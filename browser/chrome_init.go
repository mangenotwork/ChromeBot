package browser

import (
	"ChromeBot/internal/host"
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	gt "github.com/mangenotwork/gathertool"
)

// ChromeInit 初始化Chrome单例
func ChromeInit(windowSize, proxy, userPath, device string, isNew bool) {

	if isInitialized && chromeInstance != nil {
		isRun, _ := isProcessRunning(chromeInstance.PID)
		if isRun {
			utils.Debugf("Chrome已初始化 | 端口：%d | PID：%d ", chromeInstance.Port, chromeInstance.PID)
			fmt.Println("[Chrome]已初始化")
			return
		} else {
			chromeInstance = nil
			isInitialized = false
			once = sync.Once{} // 重置once，允许重新初始化
		}
	}

	once.Do(func() {

		mu.Lock()
		defer mu.Unlock()

		port := getAvailablePort() // 自定义函数：获取可用端口
		if port == 0 {
			fmt.Printf("本机未获取到可用端口!!!!")
			os.Exit(0)
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
		pid, err := startChromeProcess(chromePath, windowSize, proxy, userPath, device, port)
		if err != nil {
			fmt.Printf("启动Chrome进程失败, err = %s", err.Error())
			os.Exit(0)
		}

		AddLocalRecord(userPath, pid)

		chromeInstance = &ChromeProcess{
			WindowSize: windowSize,
			Proxy:      proxy,
			UserPath:   userPath,
			Port:       port,
			PID:        pid,
			NextID:     1, // 初始消息ID从1开始
			IsNew:      isNew,
			CloseState: false,
		}
		isInitialized = true // 标记：初始化完成

		utils.Debugf("Chrome始化成功 | 端口：%d | PID：%d ", port, pid)
		fmt.Printf("Chrome始化成功 | 端口：%d | PID：%d \n", port, pid)

		if utils.RunMode == "Script" { // 脚本模式下在启动进程后增加两秒，等待系统处理进程
			time.Sleep(2 * time.Second)
		}

		time.Sleep(1 * time.Second)
	})
}

func DefaultBrowserWS() bool {
	if chromeInstance == nil {
		fmt.Println("[Chrome]未初始化浏览器进程,请执行chrome init命令进行初始化")
		return false
	}

	url := fmt.Sprintf("http://127.0.0.1:%d/json/version", chromeInstance.Port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取浏览器debug url失败， err : ", err.Error())
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		BrowserWebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("解析浏览器debug url失败， err : ", err.Error())
		return false
	}

	chromeInstance.WebSocketDebuggerUrl = result.BrowserWebSocketDebuggerUrl
	chromeInstance.BrowserWSConn, err = ConnBrowserWS(result.BrowserWebSocketDebuggerUrl)
	if err != nil {
		fmt.Println("连接浏览器debug url失败， err : ", err.Error())
	}

	return true
}

func ConnBrowserWS(wsUrl string) (*websocket.Conn, error) {
	if wsUrl == "" {
		return nil, fmt.Errorf("url is null")
	}
	// 使用Dialer建立连接
	utils.Debug("建立连接: ", wsUrl)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	// 启动一个goroutine来接收服务器消息

	go func() {

		nowDone := make(chan struct{}, 1)
		defer close(nowDone)

		defer func() {
			if r := recover(); r != nil {
				// 转换为错误返回，而不是 panic
				err = fmt.Errorf("panic in read: %v", r)
				// 记录日志但不 panic
				log.Printf("[SafeWebSocket] Recovered panic in SafeRead: %v", r)
			}
		}()

		for {
			select {
			// case <-ConnTabDone:
			// 	fmt.Println("[Chrome] ws 连接收到结束....")

			// 	if chromeInstance.CloseState {
			// 		_ = conn.Close()
			// 		return
			// 	}

			case <-nowDone:
				fmt.Println("[Chrome] ws 连接断开执行结束....")
				return
			default:
				if conn == nil {
					fmt.Println("[Chrome] ws 连接失败")
					break
				}
				_, message, err := conn.ReadMessage()

				if err != nil {
					if err != io.EOF && !strings.Contains(err.Error(), "unexpected EOF") {
						//log.Println("接收消息失败:", err)
						time.Sleep(1 * time.Second) // 避免太快阻塞了
						continue

					} else {
						log.Println("控制谷歌似乎断开了 p = ", chromeInstance.Port, " ,err = ", err)
						chromeInstance.NowTabWSConn = nil

						// 检查4次
						for i := 0; i < 4; i++ {
							time.Sleep(2 * time.Second)
							// 检查是否进程被关闭
							isRun, err := isProcessRunning(chromeInstance.PID)
							if err != nil {
								fmt.Println("控制谷歌似乎断开了,检查进程错误, err = ", err.Error())
							}
							fmt.Println("控制谷歌似乎断开了 pid = ", chromeInstance.PID, " | isRun = ", isRun)
							if !isRun {
								fmt.Println("[Chrome]浏览器进程被关闭了,请重新初始化！")
								chromeInstance = nil
								isInitialized = false
								once = sync.Once{}
								break
							}
						}

						time.Sleep(1 * time.Second)
						nowDone <- struct{}{}
						continue
					}

				}

				msgDebug := string(message)
				if len(msgDebug) > 4000 {
					msgDebug = msgDebug[0:4000] + " --> 太多了省略 ..."
				}

				utils.Debugf("=====> 收到服务器回复: %s", msgDebug)

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
								utils.Debugf("发送页面加载事件，sessionId: %s", sessionId)
							default:
								utils.Debugf("NowPageLoadEventFired通道阻塞，跳过发送: %s", sessionId)
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
