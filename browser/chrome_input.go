package browser

import (
	"ChromeBot/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

//go:embed chrome_input.js
var chromeInputJS string

func Input(xPath, text string) error {
	if !DefaultNowTab(true) {
		return nil
	}

	utils.Debug("输入内容 : ", text)

	xPath = "'" + strings.ReplaceAll(xPath, "\"", "\\\"") + "'"
	text = "'" + strings.ReplaceAll(text, "\"", "\\\"") + "'"
	js := strings.ReplaceAll(chromeInputJS, "__XPATH__", xPath)
	js = strings.ReplaceAll(js, "__INPUTTEXT__", text)

	id := GetNextMsgID()
	msg := map[string]interface{}{
		"id":     id,
		"method": "Runtime.evaluate",
		"params": map[string]interface{}{
			"expression":    js,
			"returnByValue": true,
			"awaitPromise":  true,
		},
		"sessionId": chromeInstance.NowTabSession,
	}
	err := chromeInstance.NowTabWSConn.WriteJSON(msg)
	if err != nil {
		log.Println("发送消息失败:", err)
		return fmt.Errorf("发送消息失败")
	}
	msgStr, _ := json.Marshal(msg)
	utils.Debugf("发送消息: %s", string(msgStr))

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 重要：确保计时器被清理

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				log.Println("消息队列已关闭")
				return fmt.Errorf("消息队列已关闭")
			}
			if id == respMsg.ID {
				utils.Debug("收到的消息 -> ", respMsg.Content)
				return nil
			} else {
				utils.Debug("不是自己的消息")
			}

		case <-timer.C:
			utils.Debug("6秒未收到消息")
			return fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}
}
