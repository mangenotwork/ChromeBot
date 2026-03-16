package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"time"
)

func getAllTabData() (map[string]string, error) {
	res := make(map[string]string)
	tabUrl := fmt.Sprintf("http://127.0.0.1:%d/json/list", chromeInstance.Port)
	utils.Debug("tabUrl = ", tabUrl)

	var e2r gt.Err2Retry = true
	ctx, err := gt.Get(tabUrl, gt.RetryTimes(5), e2r, gt.ReqTimeOutMs(5000))
	if err != nil {
		log.Println(err)
		return res, err
	}
	utils.Debug("json/list = ", ctx.RespBodyString())

	dataArr := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(ctx.RespBodyString()), &dataArr)
	if err != nil {
		return res, err
	}

	log.Println("[Chrome]table List: ")
	for _, v := range dataArr {

		if v["url"].(string) == "chrome://omnibox-popup.top-chrome/" {
			// 跳过 Omnibox Popup
			continue
		}

		if v["type"].(string) == "page" {
			res[v["id"].(string)] = v["title"].(string)
			log.Printf("id : %s  title: %s", v["id"].(string), v["title"].(string))
		}
	}

	return res, nil
}

// GetAllTab 查看所有的页签
func GetAllTab() (map[string]string, error) {
	if !DefaultNowTab(false) {
		return map[string]string{}, nil
	}
	return getAllTabData()
}

// NewTab 新建标签页
func NewTab() (string, error) {
R:

	if !DefaultNowTab(false) {
		return "", nil
	}

	chromeInstance.NextID++
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "Target.createTarget",
	   "params": {
	       "url": "chrome://newtab/"
	   }
	}`, chromeInstance.NextID)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送消息失败:", err)
		chromeInstance.NowTabWSConn = nil
		goto R
		//return "", fmt.Errorf("发送消息失败")
	}
	utils.Debugf("发送消息: %s", message)

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 重要：确保计时器被清理

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				log.Println("消息队列已关闭")
				return "", fmt.Errorf("消息队列已关闭")
			}
			utils.Debug("收到的消息 -> ", respMsg.Content)
			if chromeInstance.NextID == respMsg.ID {
				result, err := gt.Json2Map(respMsg.Content)
				if err != nil {
					log.Println("回复内容解析错误")
				} else {
					utils.Debug("Target.createTarget 回复消息 : ", result)

					resultData, resultDataOK := result["result"]
					if resultDataOK {
						resultDataMap, hasMap := resultData.(map[string]any)
						if hasMap {
							targetId, targetIdHas := resultDataMap["targetId"]
							if targetIdHas {
								SelectTab(targetId.(string))
								return targetId.(string), nil
							}
						}
					}
					return "", fmt.Errorf("开启新标签页失败,未找到targetId")
				}

				return respMsg.Content, nil
			} else {
				utils.Debug("不是自己的消息")
			}

		case <-timer.C:
			log.Println("6秒未收到消息")
			return "", fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}
}

// SelectTab 切换Tab
func SelectTab(targetId string) {

	if !DefaultNowTab(false) {
		return
	}

	log.Println("[Chrome]切换Tab targetId = ", targetId)

	tabUrl := fmt.Sprintf("http://127.0.0.1:%d/json/list", chromeInstance.Port)
	utils.Debug("tabUrl = ", tabUrl)

	var e2r gt.Err2Retry = true
	ctx, err := gt.Get(tabUrl, gt.RetryTimes(5), e2r, gt.ReqTimeOutMs(5000))
	if err != nil {
		log.Println(err)
		return
	}
	utils.Debug("json/list = ", ctx.RespBodyString())
	log.Println("json/list = ", ctx.RespBodyString())

	dataArr := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(ctx.RespBodyString()), &dataArr)
	if err != nil {
		return
	}

	has := false

	for _, v := range dataArr {
		if targetId == v["id"].(string) {
			has = true
			chromeInstance.NowTabTargetId = targetId
			chromeInstance.NowTabWSUrl = v["webSocketDebuggerUrl"].(string)
			chromeInstance.NowTab = v["title"].(string)

			// go func() { ConnTabDone <- struct{}{} }()

			// 默认第一个Tab,并连接Chrome DevTools WebSocket
			wsConn, err := ConnTab()
			if err != nil {
				fmt.Println("[Chrome] 默认连接第一个Tab出现错误, err : ", err)
			}
			chromeInstance.NowTabWSConn = wsConn

			session, err := GetSession()
			if err != nil {
				fmt.Println("[Chrome] 连接Tab的session出现错误, err : ", err)
			}
			log.Println("获取到 session = ", session)
			chromeInstance.NowTabSession = session

			// 启动页面监听
			err = PageEnable()
			if err != nil {
				log.Println("页面加载失败")
			}

		}

	}

	if !has {
		fmt.Println("[Chrome]未匹配到targetId = ", targetId)
	}
}

// NowTabInfo 当前标签页的信息
func NowTabInfo() {
	if !DefaultNowTab(false) {
		return
	}
	fmt.Println("[Chrome] tab id : ", chromeInstance.NowTabTargetId)
	fmt.Println("[Chrome] tab title : ", chromeInstance.NowTab)
	fmt.Println("[Chrome] tab session : ", chromeInstance.NowTabSession)
}

// NowTabClose 关闭当前标签页
func NowTabClose() {
	if !DefaultNowTab(false) {
		return
	}

	tabList, _ := GetAllTab()
	utils.Debug("tabList len = ", len(tabList))
	if len(tabList) <= 1 {
		log.Println("[Chrome]当前页签小于等于1个，不允许被关闭")
		return
	}

	chromeInstance.NextID++
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "Target.closeTarget",
	   "params": {
	      "targetId": "%s"
	   }
	}`, chromeInstance.NextID, chromeInstance.NowTabTargetId)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送消息失败:", err)
		return
	}
	utils.Debugf("发送消息: %s", message)

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 重要：确保计时器被清理

	for {
		select {
		case msg, ok := <-messageQueue:
			if !ok {
				log.Println("消息队列已关闭")
				return
			}
			utils.Debug("收到的消息 -> ", msg.Content)
			if chromeInstance.NextID == msg.ID {
				result, err := gt.Json2Map(msg.Content)
				if err != nil {
					log.Println("回复内容解析错误")
				} else {
					utils.Debug("Target.closeTarget 回复消息 : ", result)

					resultData, resultDataOK := result["result"]
					if resultDataOK {
						resultDataMap, hasMap := resultData.(map[string]any)
						if hasMap {
							success, successHas := resultDataMap["success"]
							if successHas && success.(bool) {

								log.Println("Target.closeTarget 关闭成功, 切换tab")
								for tableItemKey, _ := range tabList {
									if tableItemKey != chromeInstance.NowTabTargetId {
										SelectTab(tableItemKey)
									}
								}

							}
						}
					}

					return
				}

				return
			} else {
				log.Println("不是自己的消息")
			}

		case <-timer.C:
			log.Println("6秒未收到消息")
			return
		}
	}
}
