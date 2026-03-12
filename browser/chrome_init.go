package browser

import (
	"ChromeBot/internal/host"
	"ChromeBot/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	gt "github.com/mangenotwork/gathertool"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
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
	WindowSize     string          // 窗口大小
	Proxy          string          // 代理
	UserPath       string          // 隔离环境
	Port           int             // 调试端口
	PID            int             // 浏览器进程
	NextID         int             // 自增消息id
	NowTab         string          // 当前操作的tab
	NowTabWSConn   *websocket.Conn // 当前操作的tab的websocket连接
	NowTabTargetId string          // 当前操作的tab的TargetId
	NowTabWSUrl    string          // 当前操作的tab的WSUrl
	NowTabSession  string          // 当前操作的tab的Session
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

// ChromeInit 初始化Chrome单例
func ChromeInit(windowSize, proxy, userPath string, isNew bool) error {

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

		mu.Lock()
		defer mu.Unlock()

		port := getAvailablePort() // 自定义函数：获取可用端口
		if port == 0 {
			initErr = errors.New("获取可用调试端口失败")
			return
		}

		chromePath, err := FindChrome()
		if err != nil {
			fmt.Printf("本机未找到Chrome浏览器，请安装后再执行")
			os.Exit(0)
		}

		utils.Debug("chromePath = ", chromePath)

		// 获取可执行文件的完整路径
		wd, _ := os.Getwd()

		// userPath 与 isNew 用时在时，优先使用 userPath
		if userPath == "" && isNew {
			fmt.Println("新建chrome隔离环境")
			n, _ := countDirectSubDirs(fmt.Sprintf("%s\\profiles\\", wd), false)
			userPath = fmt.Sprintf("%s\\profiles\\%d", wd, n)
		} else if userPath == "" && !isNew {
			userPath = fmt.Sprintf("%s\\profiles\\default", wd) // 默认
			if HasLocalRecord(userPath) {
				fmt.Printf("当前谷歌浏览器工作目录：%s 已经在运行，是否新创建一个工作目录 \n", userPath)
				isRun, _ := host.SystemConfirmBox("确认操作", fmt.Sprintf("当前谷歌浏览器工作目录：%s 已经在运行，是否新创建一个工作目录?", userPath))
				if isRun {
					n, _ := countDirectSubDirs(fmt.Sprintf("%s\\profiles\\", wd), false)
					userPath = fmt.Sprintf("%s\\profiles\\%d", wd, n)
				} else {
					fmt.Printf("当前谷歌浏览器工作目录:%s 正在被其他任务执行, 该脚本终止 \n", userPath)
					os.Exit(0)
				}
			}
		}

		utils.Debug("userPath = ", userPath)
		fmt.Printf("当前谷歌浏览器工作目录：%s\n", userPath)

		// 启动Chrome进程
		pid, err := startChromeProcess(chromePath, windowSize, proxy, userPath, port)
		if err != nil {
			initErr = fmt.Errorf("启动Chrome进程失败：%w", err)
			return
		}

		AddLocalRecord(userPath, pid)

		chromeInstance = &ChromeProcess{
			WindowSize: windowSize,
			Proxy:      proxy,
			UserPath:   userPath,
			Port:       port,
			PID:        pid,
			NextID:     1, // 初始消息ID从1开始
		}
		isInitialized = true // 标记：初始化完成

		utils.Debugf("Chrome始化成功 | 端口：%d | PID：%d ", port, pid)

		time.Sleep(1 * time.Second)
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

func startChromeProcess(chromePath, windowSize, proxy, userPath string, port int) (int, error) {

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

var ConnTabDone = make(chan struct{})

// 定义全局消息通道
type mess struct {
	ID      int
	Content string
}

var messageQueue = make(chan mess, 100) // 缓冲队列

func (c *ChromeProcess) DefaultNowTab() {
	if c.NowTabWSConn != nil {
		log.Println("当前已经选中了tab")
		return
	}
	targetId, webSocketDebuggerUrl, err := c.GetFirstTabWs()
	if err != nil {
		log.Println("[Chrome] 初始化失败 err = ", err)
	}
	c.NowTabTargetId = targetId
	c.NowTabWSUrl = webSocketDebuggerUrl

	// 默认第一个Tab,并连接Chrome DevTools WebSocket
	wsConn, err := c.ConnTab()
	if err != nil {
		fmt.Println("[Chrome] 默认连接第一个Tab出现错误, err : ", err)
	}
	c.NowTabWSConn = wsConn

	session, err := c.GetSession()
	if err != nil {
		fmt.Println("[Chrome] 默认连接第一个Tab创建session出现错误, err : ", err)
	}
	utils.Debug("获取到 session = ", session)
	c.NowTabSession = session

	// 启动页面监听
	err = c.PageEnable()
	if err != nil {
		log.Println("页面加载失败")
	}
}

func (c *ChromeProcess) ConnTab() (*websocket.Conn, error) {
	if c.NowTabWSUrl == "" {
		return nil, fmt.Errorf("url is null")
	}
	// 使用Dialer建立连接
	utils.Debug("建立连接: ", c.NowTabWSUrl)
	conn, _, err := websocket.DefaultDialer.Dial(c.NowTabWSUrl, nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	// 启动一个goroutine来接收服务器消息

	go func() {

		nowDone := make(chan struct{}, 1)
		defer close(nowDone)

		for {
			select {
			case <-ConnTabDone:
				fmt.Println("[Chrome] ws 连接收到结束....")
				_ = conn.Close()
				return
			case <-nowDone:
				fmt.Println("[Chrome] ws 连接断开执行结束....")
				return
			default:
				if conn == nil {
					fmt.Println("[Chrome] ws 连接失败")
					break
				}
				_, message, err := conn.ReadMessage()
				defer func() {
					if r := recover(); r != nil {
						// 转换为错误返回，而不是 panic
						err = fmt.Errorf("panic in read: %v", r)
						// 记录日志但不 panic
						log.Printf("[SafeWebSocket] Recovered panic in SafeRead: %v", r)
					}
				}()

				if err != nil {
					if err != io.EOF && !strings.Contains(err.Error(), "unexpected EOF") {
						//log.Println("接收消息失败:", err)
						time.Sleep(1 * time.Second) // 避免太快阻塞了
						continue

					} else {
						log.Println("控制谷歌似乎断开了 p = ", c.Port, " ,err = ", err)

						// 5秒后处理，避免马上操作
						time.Sleep(5 * time.Second)

						//FailConn <- 0
						//time.Sleep(1 * time.Second)

						nowDone <- struct{}{}
						return
					}

				}

				msgDebug := string(message)
				if len(msgDebug) > 4000 {
					msgDebug = msgDebug[0:4000] + " --> 太多了省略 ..."
				}

				utils.Debug("=====> 收到服务器回复: %s", msgDebug)

				// getRequestImg(string(message))  // 监听到图片资源

				result, err := gt.Json2Map(string(message))
				if err != nil {
					gt.Error("回复内容解析错误")
				} else {

					// 提取sessionId
					sessionId, sessionIdOK := result["sessionId"].(string)
					if !sessionIdOK {
						sessionId = ""
					}

					// 监听页面加载事件
					method, methodOK := result["method"].(string)
					if methodOK {
						// 关键修改5：优化通道发送逻辑，避免阻塞
						var sendFlag bool
						if method == "Page.loadEventFired" {
							sendFlag = true
						}

						if sendFlag && sessionId != "" {
							// 使用select+default，避免NowPageLoadEventFired无缓冲时阻塞
							select {
							case NowPageLoadEventFired <- sessionId:
								utils.Debug("发送页面加载事件，sessionId: %s", sessionId)
							default:
								utils.Debug("NowPageLoadEventFired通道阻塞，跳过发送: %s", sessionId)
							}
						}
					}

					id, ok := result["id"].(float64)
					//gt.Info(id, ok)
					if !ok {
						//gt.Error("回复消息没有id")
					} else {
						messageQueue <- mess{
							ID:      int(id),
							Content: string(message),
						}
					}

				}

			}
		}
	}()
	return conn, nil
}

func (c *ChromeProcess) CloseNowTabConn() {
	ConnTabDone <- struct{}{}
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

	//mu.Lock()
	//defer mu.Unlock()

	if !isInitialized || chromeInstance == nil {
		fmt.Println("[Chrome]未初始化")
		return nil
	}

	isRun, _ := isProcessRunning(c.PID)
	if !isRun {
		fmt.Println("[Chrome]未初始化")
		return nil
	}
	utils.Debug("关闭WS连接")
	// 关闭WS连接
	go c.CloseNowTabConn()

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

func (c *ChromeProcess) GetFirstTabWs() (string, string, error) {
	utils.Debug("c.UserPath = ", c.UserPath)
	isNew := !utils.PathExists(c.UserPath)
	utils.Debug("isNew = ", isNew)

	tabUrl := fmt.Sprintf("http://127.0.0.1:%d/json/list", c.Port)
	utils.Debug("tabUrl = ", tabUrl)

	var targetId = ""
	var webSocketDebuggerUrl = ""
	var e2r gt.Err2Retry = true
	ctx, err := gt.Get(tabUrl, gt.RetryTimes(5), e2r, gt.ReqTimeOutMs(5000))
	if err != nil {
		gt.Error(err)
		return "", "", err
	}
	utils.Debug("json/list = ", ctx.RespBodyString())

	dataArr := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(ctx.RespBodyString()), &dataArr)
	if err != nil {
		return "", "", err
	}

	utils.Debug("rList = ", dataArr)

	if isNew && len(dataArr) > 1 {
		dataMap := gt.Any2Map(dataArr[1])
		targetId = dataMap["id"].(string)
		webSocketDebuggerUrl = dataMap["webSocketDebuggerUrl"].(string)
		c.NowTab = dataMap["title"].(string)
	} else if len(dataArr) > 0 {
		dataMap := gt.Any2Map(dataArr[0])
		targetId = dataMap["id"].(string)
		webSocketDebuggerUrl = dataMap["webSocketDebuggerUrl"].(string)
		c.NowTab = dataMap["title"].(string)
	} else {
		return "", "", err
	}

	utils.Debugf("ws url %s", webSocketDebuggerUrl)
	return targetId, webSocketDebuggerUrl, nil
}

func (c *ChromeProcess) GetSession() (string, error) {
	_, err := c.activateTarget()
	if err != nil {
		log.Println("[Chrome]执行activateTarget遇到错误， err = ", err.Error())
	}
	return c.attachToTarget()
}

func (c *ChromeProcess) attachToTarget() (string, error) {
	c.NextID++
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "Target.attachToTarget",
	   "params": {
	       "targetId": "%s",
        	"flatten": true
	   }
	}`, c.NextID, c.NowTabTargetId)
	err := c.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("[Chrome]发送消息失败:", err)
		return "", fmt.Errorf("发送消息失败")
	}
	utils.Debugf("发送消息: %s", message)

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 重要：确保计时器被清理

	for {
		select {
		case msg, ok := <-messageQueue:
			if !ok {
				gt.Info("消息队列已关闭")
				return "", fmt.Errorf("消息队列已关闭")
			}
			utils.Debug("收到的消息 -> ", msg.Content)
			if c.NextID == msg.ID {
				result, err := gt.Json2Map(msg.Content)
				if err != nil {
					gt.Error("回复内容解析错误")
				} else {
					resultData, resultDataOK := result["result"]
					if resultDataOK {
						resultDataMap, hasMap := resultData.(map[string]any)
						if hasMap {
							sessionIdData, sessionIdHas := resultDataMap["sessionId"]
							if sessionIdHas {
								return sessionIdData.(string), nil
							}
						}
					}
				}

				return msg.Content, nil
			} else {
				gt.Info("不是自己的消息")
			}

		case <-timer.C:
			gt.Info("6秒未收到消息")
			return "", fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}

}

func (c *ChromeProcess) activateTarget() (string, error) {
	c.NextID++
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "Target.activateTarget",
	   "params": {
	       "targetId": "%s",
        	"flatten": true
	   }
	}`, c.NextID, c.NowTabTargetId)
	err := c.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("[Chrome]发送消息失败:", err)
		return "", fmt.Errorf("发送消息失败")
	}
	utils.Debugf("发送消息: %s", message)

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 重要：确保计时器被清理

	for {
		select {
		case msg, ok := <-messageQueue:
			if !ok {
				utils.Debug("消息队列已关闭")
				return "", fmt.Errorf("消息队列已关闭")
			}
			utils.Debug("收到的消息 -> ", msg.Content)
			if c.NextID == msg.ID {
				result, err := gt.Json2Map(msg.Content)
				if err != nil {
					utils.Debug("回复内容解析错误")
				} else {
					resultData, resultDataOK := result["result"]
					if resultDataOK {
						resultDataMap, hasMap := resultData.(map[string]any)
						if hasMap {
							sessionIdData, sessionIdHas := resultDataMap["sessionId"]
							if sessionIdHas {
								return sessionIdData.(string), nil
							}
						}
					}
				}

				return msg.Content, nil
			} else {
				utils.Debug("不是自己的消息")
			}

		case <-timer.C:
			utils.Debug("6秒未收到消息")
			return "", fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}
}

func (c *ChromeProcess) OpenUrl(url string) (string, error) {

	if c.NowTabWSConn == nil {
		c.DefaultNowTab()
	}

	c.NextID++
	message := fmt.Sprintf(`{
	  "id": %d,
	  "method": "Page.navigate",
	  "params": {
	      "url": "%s"
	  },
		"sessionId":"%s"
	}`, c.NextID, utils.FixURLProtocol(url), c.NowTabSession)
	err := c.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送消息失败:", err)
		return "", fmt.Errorf("发送消息失败")
	}
	utils.Debugf("发送消息: %s", message)

	timeout := 30 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 重要：确保计时器被清理

	for {
		select {
		case msg, ok := <-messageQueue:
			if !ok {
				log.Println("消息队列已关闭")
				return "", fmt.Errorf("消息队列已关闭")
			}
			utils.Debug("收到的消息 -> ", msg.Content)
			if c.NextID == msg.ID {

				select {
				case session := <-NowPageLoadEventFired:
					utils.Debug("页面已完全加载 session = ", session)
					return msg.Content, nil
				case <-time.After(6 * time.Second):
					return "", fmt.Errorf("页面加载超时")
				}

				//return msg.Content, nil

			} else {
				log.Println("不是自己的消息")
			}

		case <-timer.C:
			log.Println("30秒未收到消息")
			return "", fmt.Errorf("接收消息超时; 30秒未收到消息")
		}
	}

}
