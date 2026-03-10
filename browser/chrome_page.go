package browser

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (c *ChromeProcess) PageEnable() error {
	if c.NowTabWSConn == nil {
		c.DefaultNowTab()
	}

	// 1. 启用Page事件监听（必须）
	c.NextID++
	msg := map[string]interface{}{
		"id":        c.NextID,
		"method":    "Page.enable",
		"params":    map[string]interface{}{},
		"sessionId": c.NowTabSession,
	}
	err := c.NowTabWSConn.WriteJSON(msg)
	if err != nil {
		log.Println("发送消息失败:", err)
		return fmt.Errorf("发送消息失败")
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
				return fmt.Errorf("消息队列已关闭")
			}
			if c.NextID == respMsg.ID {
				log.Println("收到的消息 -> ", respMsg.Content)
				return nil
			} else {
				log.Println("不是自己的消息")
			}

		case <-timer.C:
			log.Println("6秒未收到消息")
			return fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}

}
