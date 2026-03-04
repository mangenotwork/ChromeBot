package browser

import (
	"ChromeBot/utils"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
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

// ChromeInit 初始化Chrome单例
func ChromeInit(windowSize, proxy, userPath string) error {

	if isInitialized && chromeInstance != nil {
		isRun, _ := isProcessRunning(chromeInstance.PID)
		if isRun {
			utils.Debugf("Chrome已初始化 | 端口：%d | PID：%d ", chromeInstance.Port, chromeInstance.PID)
			fmt.Println("[Chrome]已初始化")
			return nil
		} else {
			chromeInstance = nil
			isInitialized = false
			once = sync.Once{} // 重置once，允许重新初始化
		}
	}

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

		utils.Debugf("Chrome始化成功 | 端口：%d | PID：%d ", port, pid)
	})

	return initErr
}

func getAvailablePort() int {
	listener, err := net.Listen("tcp", "0.0.0.0:0") // 关键：绑定0.0.0.0确保外部可访问
	if err != nil {
		utils.Debugf("创建监听器失败: %s", err.Error())
		return 0
	}
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}

func startChromeProcess(windowSize, proxy, userPath string, port int) (int, error) {

	userPath = fmt.Sprintf("./ChromeBot/profiles/default") // 谷歌目录下  \Google\Chrome\Application\

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

	pid := cmd.Process.Pid

	for i := 0; i < 40; i++ {
		if ok, _ := isProcessRunning(pid); ok {
			break
		}
		time.Sleep(40 * time.Millisecond)
	}
	return pid, nil
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

	if !isInitialized || chromeInstance == nil {
		fmt.Println("[Chrome]未初始化")
		return nil
	}

	isRun, _ := isProcessRunning(c.PID)
	if !isRun {
		fmt.Println("[Chrome]未初始化")
		return nil
	}

	// 1. 关闭WS连接
	if c.WSConn != nil {
		_ = c.WSConn.Close()
		c.WSConn = nil
	}

	utils.Debug("c.PID = ", c.PID)

	// 2. 杀死Chrome进程
	if c.PID != 0 {
		if err := SafeKillProcess(c.PID); err != nil {
			utils.Debug("[ERR]关闭进程错误:", err.Error())
			return err
		}
	}

	// 3. 重置单例（可选）
	chromeInstance = nil
	isInitialized = false
	once = sync.Once{} // 重置once，允许重新初始化

	fmt.Printf("[Chrome]浏览器进程已关闭 | PID：%d \n", c.PID)
	return nil
}

func SafeKillProcess(pid int) error {
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		// 先尝试 Windows API
		if err := killProcessByPID(pid); err != nil {
			utils.Debug("Windows API err : ", err.Error())
		}

		// 再尝试 taskkill
		utils.Debug("exce taskkill")
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
		if err := cmd.Run(); err != nil {
			utils.Debug("taskkill执行失败")
		}

		output, err := cmd.CombinedOutput() // 同时捕获stdout/stderr
		if err != nil {
			gbkOutput, decodeErr := gbkToUtf8(output)
			if decodeErr != nil {
				// 解码失败则用原始字符串（避免二次错误）
				gbkOutput = strings.TrimSpace(string(output))
			}
			utils.Debugf("taskkill执行失败 | PID：%d | 退出码：%v | 错误详情：%s", pid, err, gbkOutput)
		}

		if i < maxRetries-1 {
			time.Sleep(time.Duration(400*(i+1)) * time.Millisecond)
		}

		// 检查进程是否已结束
		isRun, _ := isProcessRunning(pid)
		utils.Debug("isRun = ", isRun)
		if isRun {
			continue
		} else {
			return nil
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

// gbkToUtf8 将GBK编码的字节数组转为UTF-8字符串（核心解码函数）
func gbkToUtf8(gbkBytes []byte) (string, error) {
	// 创建GBK转UTF-8的转换器
	reader := transform.NewReader(strings.NewReader(string(gbkBytes)), simplifiedchinese.GBK.NewDecoder())
	// 读取转换后的字节
	utf8Bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(utf8Bytes)), nil
}

func killProcessByPID(pid int) error {

	utils.Debug("exce killProcessByPID")

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
		fmt.Println("[Chrome]未初始化")
		return 0
	}
	return c.PID
}
