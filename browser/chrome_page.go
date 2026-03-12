package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func PageEnable() error {
	if !DefaultNowTab() {
		return fmt.Errorf("浏览器未初始化")
	}

	// 1. 启用Page事件监听（必须）
	chromeInstance.NextID++
	msg := map[string]interface{}{
		"id":        chromeInstance.NextID,
		"method":    "Page.enable",
		"params":    map[string]interface{}{},
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
			if chromeInstance.NextID == respMsg.ID {
				utils.Debug("收到的消息 -> ", respMsg.Content)

				chromeInstance.NextID++
				disableFrameEvents := map[string]interface{}{
					"id":        chromeInstance.NextID,
					"method":    "Page.setLifecycleEventsEnabled", // 禁用生命周期事件
					"params":    map[string]interface{}{"enabled": false},
					"sessionId": chromeInstance.NowTabSession,
				}
				if err := chromeInstance.NowTabWSConn.WriteJSON(disableFrameEvents); err != nil {
					utils.Debug("禁用frame事件失败: ", err)
				}

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
