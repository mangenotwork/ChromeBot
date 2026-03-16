package browser

import (
	"ChromeBot/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"strings"
	"time"
)

//go:embed chrome_click.js
var chromeClickJS string

func Click(xPath string) error {
	if !DefaultNowTab(true) {
		return fmt.Errorf("浏览器未初始化")
	}

	xPath = "'" + strings.ReplaceAll(xPath, "\"", "\\\"") + "'"
	js := strings.ReplaceAll(chromeClickJS, "__XPATH__", xPath)

	chromeInstance.NextID++
	msg := map[string]interface{}{
		"id":     chromeInstance.NextID,
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
			utils.Debug("收到的消息 -> ", utils.UnescapeUnicode(respMsg.Content))
			if chromeInstance.NextID == respMsg.ID {

				resultValue, err := gt.JsonFind(respMsg.Content, "result/result/subtype")
				if err != nil {
					log.Println(err)
				}
				if resultValue == "error" {
					log.Println("点击出现了错误: ", respMsg.Content)
					return fmt.Errorf("点击出现了错误: %s", respMsg.Content)
				}

				resultValue, err = gt.JsonFind(respMsg.Content, "result/result/value")
				if err != nil {
					log.Println(err)
				}
				utils.Debug("执行点击操作: ", resultValue)

				select {
				case session := <-NowPageLoadEventFired:
					utils.Debug("点击后页面已完全加载 session = ", session)
					return nil
				case <-time.After(6 * time.Second):
					return nil
				}

			} else {
				utils.Debug("不是自己的消息")
			}

		case <-timer.C:
			utils.Debug("6秒未收到消息")
			return fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}
}

// todo  如果是知道坐标可以使用 cdp来替代js点击
/*
// CDP触发鼠标点击
func (c *ChromeProcess) CDPClick(x, y int) error {
    c.NextID++
    msg := map[string]interface{}{
        "id":     c.NextID,
        "method": "Input.dispatchMouseEvent",
        "params": map[string]interface{}{
            "type": "mousePressed",
            "x":    x,
            "y":    y,
            "button": "left",
            "clickCount": 1,
        },
        "sessionId": c.NowTabSession,
    }
    // 发送请求...
}
....
*/
