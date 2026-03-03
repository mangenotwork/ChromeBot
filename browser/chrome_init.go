package browser

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"golang.org/x/sys/windows"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ChromeProcess struct {
	WindowSize string          // 窗口大小
	Proxy      string          // 代理
	UserPath   string          // 隔离环境
	Port       int             // 调试端口
	PID        int             // 浏览器进程
	NextID     int             // 自增消息id
	WSConn     *websocket.Conn // websocket连接
}

var (
	chromeInstance *ChromeProcess
	once           sync.Once
	mu             sync.RWMutex
	isInitialized  bool // 标识：是否已完成初始化
)

// GetChromeInstance 获取Chrome
func GetChromeInstance() *ChromeProcess {
	mu.RLock()
	defer mu.RUnlock()
	return chromeInstance
}

func IsChromeInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return isInitialized
}

// ChromeInit 初始化Chrome单例
func ChromeInit(windowSize, proxy, userPath string) error {

	mu.RLock()
	if isInitialized && chromeInstance != nil {
		mu.RUnlock()
		log.Printf("Chrome已初始化 | 端口：%d | PID：%d ", chromeInstance.PID, chromeInstance.Port)
		return nil
	}
	mu.RUnlock()

	var initErr error
	once.Do(func() {

		port := getAvailablePort() // 自定义函数：获取可用端口
		if port == 0 {
			initErr = errors.New("获取可用调试端口失败")
			return
		}

		// 3. 启动Chrome进程
		pid, err := startChromeProcess(windowSize, proxy, userPath, port)
		if err != nil {
			initErr = fmt.Errorf("启动Chrome进程失败：%w", err)
			return
		}

		//// 4. 连接Chrome DevTools WebSocket
		//wsConn, err := connectChromeWS(port)
		//if err != nil {
		//	// 启动失败则杀死进程，避免僵尸进程
		//	killChromeProcess(pid)
		//	initErr = fmt.Errorf("连接Chrome WS失败：%w", err)
		//	return
		//}

		// 5. 初始化单例实例
		mu.Lock()
		chromeInstance = &ChromeProcess{
			WindowSize: windowSize,
			Proxy:      proxy,
			UserPath:   userPath,
			Port:       port,
			PID:        pid,
			NextID:     1, // 初始消息ID从1开始
			//WSConn:     wsConn,
		}
		isInitialized = true // 标记：初始化完成
		mu.Unlock()

		log.Printf("Chrome始化成功 | 端口：%d | PID：%d ", port, pid)
	})

	return initErr
}

func getAvailablePort() int {
	listener, err := net.Listen("tcp", "0.0.0.0:0") // 关键：绑定0.0.0.0确保外部可访问
	if err != nil {
		log.Printf("创建监听器失败: %s", err.Error())
		return 0
	}
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}

func startChromeProcess(windowSize, proxy, userPath string, port int) (int, error) {

	chromePath, err := FindChrome()
	if err != nil {
		fmt.Printf("本机未找到Chrome浏览器，请安装后再执行")
		os.Exit(0)
	}

	args := []string{
		"--remote-debugging-port=" + strconv.Itoa(port), // 远程调试端口
		"--no-first-run",
	}

	if windowSize != "" {
		windowSize = strings.Replace(windowSize, "*", ",", -1)
		args = append(args, "--window-size="+windowSize)
	}

	if userPath != "" {
		args = append(args, "--user-data-dir="+userPath)
	}

	if proxy != "" {
		args = append(args, "--proxy-server="+proxy)
	}

	cmd := exec.Command(chromePath, args...)
	// 启动进程（不阻塞）
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	time.Sleep(2 * time.Second) // 强制休息2秒

	// 返回进程PID
	return cmd.Process.Pid, nil
}

// connectChromeWS 连接Chrome DevTools WebSocket
func connectChromeWS(port int) (*websocket.Conn, error) {
	//// Chrome DevTools的WS地址格式（需先获取目标页面的WS地址，此处简化）
	//// 实际场景需先调用 http://127.0.0.1:port/json/version 获取WS地址
	//wsURL := url.URL{
	//	Scheme: "ws",
	//	Host:   fmt.Sprintf("127.0.0.1:%d", port),
	//	Path:   "/devtools/browser/abc123", // 实际需动态获取，此处简化
	//}
	//
	//// 连接WS
	//conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

// GetNextMsgID 获取自增的消息ID（线程安全）
func (c *ChromeProcess) GetNextMsgID() int {
	mu.Lock()
	defer mu.Unlock()
	id := c.NextID
	c.NextID++
	return id
}

// Close 关闭Chrome实例（释放WS连接+杀死进程）
func (c *ChromeProcess) Close() error {

	mu.Lock()
	defer mu.Unlock()

	if !isInitialized || chromeInstance != nil {
		log.Printf("Chrome未初始化")
		return nil
	}

	// 1. 关闭WS连接
	if c.WSConn != nil {
		_ = c.WSConn.Close()
		c.WSConn = nil
	}

	log.Println("c.PID = ", c.PID)

	// 2. 杀死Chrome进程
	if c.PID != 0 {
		if err := SafeKillProcess(c.PID); err != nil {
			return err
		}
	}

	// 3. 重置单例（可选）
	chromeInstance = nil
	isInitialized = false
	once = sync.Once{} // 重置once，允许重新初始化

	log.Printf("Chrome实例已关闭 | PID：%d", c.PID)
	return nil
}

func SafeKillProcess(pid int) error {
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		// 先尝试 Windows API
		if err := killProcessByPID(pid); err == nil {
			return nil
		}

		// 再尝试 taskkill
		log.Println("exce taskkill")
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
		if err := cmd.Run(); err == nil {
			return nil
		}

		// 检查进程是否已结束
		if runing, _ := isProcessRunning(pid); !runing {
			return nil
		}

		if i < maxRetries-1 {
			time.Sleep(time.Duration(400*(i+1)) * time.Millisecond)
		}
	}

	return fmt.Errorf("failed to kill process %d after %d attempts", pid, maxRetries)
}

// 检查进程是否存在的辅助函数
func isProcessRunning(pid int) (bool, error) {
	// 尝试打开进程查询权限
	handle, err := windows.OpenProcess(
		windows.PROCESS_QUERY_INFORMATION,
		false,
		uint32(pid),
	)
	if err != nil {
		if err == windows.ERROR_INVALID_PARAMETER {
			// 进程不存在
			return false, nil
		}
		return false, err
	}
	defer windows.CloseHandle(handle)
	// 检查进程退出代码
	var exitCode uint32
	err = windows.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		return false, err
	}
	// 在Windows中，259表示进程仍在运行
	// STILL_ACTIVE 的值为 259
	return exitCode == 259, nil
}

func killProcessByPID(pid int) error {

	log.Println("exce killProcessByPID")

	handle, err := windows.OpenProcess(
		windows.PROCESS_TERMINATE,
		false,
		uint32(pid),
	)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	return windows.TerminateProcess(handle, 1)
}

func (c *ChromeProcess) GetPID() int {
	if !isInitialized || chromeInstance != nil {
		log.Printf("Chrome未初始化")
		return 0
	}
	return c.PID
}
