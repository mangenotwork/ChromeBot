package browser

import (
	"ChromeBot/utils"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func OpenUrl(url string) (string, error) {
	if !DefaultNowTab(false) {
		return "", nil
	}

	chromeInstance.NextID++
	message := fmt.Sprintf(`{
	  "id": %d,
	  "method": "Page.navigate",
	  "params": {
	      "url": "%s"
	  },
		"sessionId":"%s"
	}`, chromeInstance.NextID, utils.FixURLProtocol(url), chromeInstance.NowTabSession)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
			utils.Debug("OpenUrl 收到的消息 -> ", msg.Content)
			if chromeInstance.NextID == msg.ID {
				utils.Debug("OpenUrl 匹配到消息 -> ", msg.ID)

				select {
				case session := <-NowPageLoadEventFired:
					utils.Debug("页面已完全加载 session = ", session)
					return msg.Content, nil
				case <-time.After(6 * time.Second):
					return "", fmt.Errorf("页面加载超时")
				}

				//return msg.Content, nil

			} //else {
			//log.Println("不是自己的消息")
			//}

		case <-timer.C:
			log.Println("30秒未收到消息")
			return "", fmt.Errorf("接收消息超时; 30秒未收到消息")
		}
	}

}
