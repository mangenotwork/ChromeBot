package browser

import (
	_ "embed"
	"encoding/json"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"strings"
	"time"
)

//go:embed chrome_check.js
var chromeCheckJS string

func (c *ChromeProcess) Check(xPath string) (bool, error) {
	if c.NowTabWSConn == nil {
		c.DefaultNowTab()
	}

	xPath = "\"" + strings.ReplaceAll(xPath, "\"", "\\\"") + "\""
	js := strings.ReplaceAll(chromeCheckJS, "__XPATH__", xPath)

	c.NextID++
	msg := map[string]interface{}{
		"id":     c.NextID,
		"method": "Runtime.evaluate",
		"params": map[string]interface{}{
			"expression":    js,
			"returnByValue": true,
			"awaitPromise":  true,
		},
		"sessionId": c.NowTabSession,
	}
	err := c.NowTabWSConn.WriteJSON(msg)
	if err != nil {
		gt.Error("发送消息失败:", err)
		return false, fmt.Errorf("发送消息失败")
	}
	msgStr, _ := json.Marshal(msg)
	log.Printf("发送消息: %s", string(msgStr))

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 重要：确保计时器被清理

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				log.Println("消息队列已关闭")
				return false, fmt.Errorf("消息队列已关闭")
			}
			log.Println("收到的消息 -> ", respMsg.Content)
			if c.NextID == respMsg.ID {
				resultValue, err := gt.JsonFind(respMsg.Content, "/result/result/value")
				if err != nil {
					log.Println(err)
				}
				if gt.Any2String(resultValue) == "true" {
					return true, nil
				}
				return false, nil
			} else {
				log.Println("不是自己的消息")
			}

		case <-timer.C:
			log.Println("6秒未收到消息")
			return false, fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}
}
