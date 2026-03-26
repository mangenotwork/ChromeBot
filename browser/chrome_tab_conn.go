package browser

import (
	"ChromeBot/internal/host"
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	gt "github.com/mangenotwork/gathertool"
)

// DefaultNowTab 默认当前交互Tab
// isOP 是否是操作， 点击，输入，截图
func DefaultNowTab(isOP bool) bool {
	if chromeInstance == nil {
		fmt.Println("[Chrome]未初始化浏览器进程,请执行chrome init命令进行初始化")
		return false
	}

	if isOP {
		fmt.Println("CheckTab ....")
		CheckTab() // 检查当前操作是tab是否被意外丢失
	}

	if chromeInstance.NowTabWSConn != nil {
		fmt.Println("[Chrome]获取浏览器tab页聚焦")
		return true
	}

	windowSize := chromeInstance.WindowSize
	proxy := chromeInstance.Proxy
	userPath := chromeInstance.UserPath
	device := chromeInstance.Device
	isNew := chromeInstance.IsNew
	retryTimes := 4
	firstTabWsOK := false

	targetId, webSocketDebuggerUrl, err := GetFirstTabWs()
	if err != nil {
		log.Println("[Chrome] 初始化失败 err = ", err)
		for i := 0; i < retryTimes; i++ {
			_ = Close()
			ChromeInit(windowSize, proxy, userPath, device, isNew)
			var newErr error
			targetId, webSocketDebuggerUrl, newErr = GetFirstTabWs()
			if newErr == nil {
				firstTabWsOK = true
				break
			}
		}
		if !firstTabWsOK {
			_, _ = host.ErrorTipBox("浏览器初始化失败！")
			os.Exit(0)
		}
	}
	chromeInstance.NowTabTargetId = targetId
	chromeInstance.NowTabWSUrl = webSocketDebuggerUrl

	// 默认第一个Tab,并连接Chrome DevTools WebSocket
	wsConn, err := ConnTab()
	if err != nil {
		fmt.Println("[Chrome] 默认连接第一个Tab出现错误, err : ", err)
	}
	chromeInstance.NowTabWSConn = wsConn

	session, err := GetSession()
	if err != nil {
		fmt.Println("[Chrome] 默认连接第一个Tab创建session出现错误, err : ", err)
	}
	utils.Debug("获取到 session = ", session)
	chromeInstance.NowTabSession = session

	// 启动页面监听
	err = PageEnable()
	if err != nil {
		log.Println("页面加载失败")
	}
	return true
}

func ConnTab() (*websocket.Conn, error) {
	if chromeInstance.NowTabWSUrl == "" {
		return nil, fmt.Errorf("url is null")
	}
	// 使用Dialer建立连接
	utils.Debug("建立连接: ", chromeInstance.NowTabWSUrl)
	conn, _, err := websocket.DefaultDialer.Dial(chromeInstance.NowTabWSUrl, nil)
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
			case <-ConnTabDone:
				fmt.Println("[Chrome] ws 连接收到结束....")

				if chromeInstance.CloseState {
					_ = conn.Close()
					return
				}

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

func CloseNowTabConn() {
	ConnTabDone <- struct{}{}
}

func GetFirstTabWs() (string, string, error) {
	utils.Debug("c.UserPath = ", chromeInstance.UserPath)
	isNew := !utils.PathExists(chromeInstance.UserPath)
	utils.Debug("isNew = ", isNew)
	fmt.Println("当前进程port = ", chromeInstance.Port)
	tabUrl := fmt.Sprintf("http://127.0.0.1:%d/json/list", chromeInstance.Port)
	utils.Debug("tabUrl = ", tabUrl)
	fmt.Println("tabUrl = ", tabUrl)

	var targetId = ""
	var webSocketDebuggerUrl = ""
	var e2r gt.Err2Retry = true
	ctx, err := gt.Get(tabUrl, gt.RetryTimes(2), e2r, gt.ReqTimeOutMs(2000))
	if err != nil {
		fmt.Println("请求失败: ", err)
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
		chromeInstance.NowTab = dataMap["title"].(string)
	} else if len(dataArr) > 0 {
		dataMap := gt.Any2Map(dataArr[0])
		targetId = dataMap["id"].(string)
		webSocketDebuggerUrl = dataMap["webSocketDebuggerUrl"].(string)
		chromeInstance.NowTab = dataMap["title"].(string)
	} else {
		return "", "", err
	}

	utils.Debugf("ws url %s", webSocketDebuggerUrl)
	return targetId, webSocketDebuggerUrl, nil
}

func GetSession() (string, error) {
	_, err := activateTarget()
	if err != nil {
		log.Println("[Chrome]执行activateTarget遇到错误， err = ", err.Error())
	}
	return attachToTarget()
}

func attachToTarget() (string, error) {
	chromeInstance.NextID++
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "Target.attachToTarget",
	   "params": {
	       "targetId": "%s",
        	"flatten": true
	   }
	}`, chromeInstance.NextID, chromeInstance.NowTabTargetId)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
			if chromeInstance.NextID == msg.ID {
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

func activateTarget() (string, error) {

	if chromeInstance.NowTabWSConn == nil {
		// todo
	}

	chromeInstance.NextID++
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "Target.activateTarget",
	   "params": {
	       "targetId": "%s",
        	"flatten": true
	   }
	}`, chromeInstance.NextID, chromeInstance.NowTabTargetId)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
			if chromeInstance.NextID == msg.ID {
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
