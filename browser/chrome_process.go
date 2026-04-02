package browser

import (
	"ChromeBot/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type ChromeProcess struct {
	WindowSize           string          // 窗口大小
	Proxy                string          // 代理
	UserPath             string          // 隔离环境
	Device               string          // 设备
	Port                 int             // 调试端口
	PID                  int             // 浏览器进程
	NextID               int             // 自增消息id
	NowTab               string          // 当前操作的tab
	NowTabWSConn         *websocket.Conn // 当前操作的tab的websocket连接
	NowTabTargetId       string          // 当前操作的tab的TargetId
	NowTabWSUrl          string          // 当前操作的tab的WSUrl
	NowTabSession        string          // 当前操作的tab的Session
	IsNew                bool            // 是否是新隔离环境
	CloseState           bool            // 关闭状态
	WebSocketDebuggerUrl string          // 浏览器的debugger调试url
	BrowserWSConn        *websocket.Conn // 当前浏览器的debugger调试WS连接
}

var (
	chromeInstance *ChromeProcess
	once           sync.Once
	mu             sync.RWMutex
	isInitialized  bool // 标识：是否已完成初始化
)

// GetChromeInstance 获取Chrome
func GetChromeInstance() *ChromeProcess {
	//mu.RLock()
	//defer mu.RUnlock()
	return chromeInstance
}

func GetNowTabSession() string {
	return chromeInstance.NowTabSession
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

func startChromeProcess(chromePath, windowSize, proxy, userPath, device string, port int) (int, error) {

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

	if deviceData, ok := chromeDevice[device]; ok {
		fmt.Println("设置设备:", device)
		args = append(args, deviceData.userAgent)
		args = append(args, deviceData.windowSize)
	}

	cmd := exec.Command(chromePath, args...)
	// 启动进程（不阻塞）
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	pid := cmd.Process.Pid
	fmt.Println("[Chrome]浏览器进程 PID: ", pid)

	for i := 0; i < 40; i++ {
		if ok, _ := isProcessRunning(pid); ok {
			return pid, nil
		}
		time.Sleep(40 * time.Millisecond)
	}

	// 最后一次检查
	if ok, _ := isProcessRunning(pid); ok {
		return pid, nil
	} else {
		return 0, fmt.Errorf("[Chrome]未找到进程")
	}

}

// Close 关闭Chrome实例（释放WS连接+杀死进程）
func Close() error {

	//mu.Lock()
	//defer mu.Unlock()

	if !isInitialized || chromeInstance == nil {
		fmt.Println("[Chrome]未初始化")
		return nil
	}

	isRun, _ := isProcessRunning(chromeInstance.PID)
	if !isRun {
		fmt.Println("[Chrome]未初始化")
		return nil
	}

	chromeInstance.CloseState = true

	utils.Debug("关闭WS连接")
	// 关闭WS连接
	go CloseNowTabConn()
	utils.Debug("c.PID = ", chromeInstance.PID)

	if chromeInstance.PID != 0 {
		if err := SafeKillProcess(chromeInstance.PID); err != nil {
			utils.Debug("[ERR]关闭进程错误:", err.Error())
			return err
		}
	}

	fmt.Printf("[Chrome]浏览器进程已关闭 | PID：%d \n", chromeInstance.PID)

	chromeInstance = nil
	isInitialized = false
	once = sync.Once{}
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
	log.Println("exitCode = ", exitCode)
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

func GetPID() int {
	if !isInitialized || chromeInstance != nil {
		fmt.Println("[Chrome]未初始化")
		return 0
	}
	return chromeInstance.PID
}
