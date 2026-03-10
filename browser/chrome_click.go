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

func (c *ChromeProcess) Click(xPath string) error {
	if c.NowTabWSConn == nil {
		c.DefaultNowTab()
	}

	xPath = "'" + strings.ReplaceAll(xPath, "\"", "\\\"") + "'"
	js := strings.ReplaceAll(chromeClickJS, "__XPATH__", xPath)

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
			if c.NextID == respMsg.ID {

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

// todo 下面都是点击参考的代码

//func MouseEventsBtn(conn *websocket.Conn, sessionId string, id int, xpath string) (string, error) {
//	js := fmt.Sprintf(`
//				(function() {
//					const xpath = %s;
//					const element = document.evaluate(xpath, document, null,
//						XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
//
//					if (!element) return {success: false};
//
//					const rect = element.getBoundingClientRect();
//					const x = rect.left + rect.width / 2;
//					const y = rect.top + rect.height / 2;
//
//					// 触发完整的鼠标事件序列
//					const events = [
//						['mousedown', {button: 0, clientX: x, clientY: y}],
//						['mouseup', {button: 0, clientX: x, clientY: y}],
//						['click', {button: 0, clientX: x, clientY: y}]
//					];
//
//					events.forEach(([type, options]) => {
//						const event = new MouseEvent(type, {
//							view: window,
//							bubbles: true,
//							cancelable: true,
//							button: options.button,
//							buttons: 1,
//							clientX: options.clientX,
//							clientY: options.clientY
//						});
//						element.dispatchEvent(event);
//					});
//
//					return {success: true};
//				})()
//			`, jsonEscapeString(xpath))
//
//	msg := map[string]interface{}{
//		"id":     id,
//		"method": "Runtime.evaluate",
//		"params": map[string]interface{}{
//			"expression":    js,
//			"returnByValue": true,
//			"awaitPromise":  true,
//		},
//		"sessionId": sessionId,
//	}
//	err := conn.WriteJSON(msg)
//	if err != nil {
//		gt.Error("发送消息失败:", err)
//		return "", fmt.Errorf("发送消息失败")
//	}
//	//msgStr, _ := json.Marshal(msg)
//	//log.Printf("发送消息: %s", string(msgStr))
//
//	timeout := 6 * time.Second
//	timer := time.NewTimer(timeout)
//	defer timer.Stop() // 重要：确保计时器被清理
//
//	for {
//		select {
//		case msg, ok := <-messageQueue:
//			if !ok {
//				gt.Info("消息队列已关闭")
//				return "", fmt.Errorf("消息队列已关闭")
//			}
//			//gt.Info("收到的消息 -> ", msg.Content)
//			if id == msg.ID {
//				//gt.Info("是自己的消息")
//				return msg.Content, nil
//			} else {
//				gt.Info("不是自己的消息")
//			}
//
//		case <-timer.C:
//			gt.Info("6秒未收到消息")
//			return "", fmt.Errorf("接收消息超时; 6秒未收到消息")
//		}
//	}
//
//}
//
//func ReactDispatchEventBtn(conn *websocket.Conn, sessionId string, id int, xpath string) (string, error) {
//	js := fmt.Sprintf(`
//				(function() {
//					const xpath = %s;
//					const element = document.evaluate(xpath, document, null,
//						XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
//
//					if (!element) return {success: false};
//
//					// 创建并触发自定义事件
//					const event = new Event('click', {
//						bubbles: true,
//						cancelable: true
//					});
//
//					// 添加一些框架可能需要的属性
//					event._synthetic = true;
//					event._reactName = 'onClick';
//					event.nativeEvent = new MouseEvent('click', { bubbles: true });
//
//					element.dispatchEvent(event);
//
//					// 也触发原生的click
//					element.click();
//
//					return {success: true};
//				})()
//			`, jsonEscapeString(xpath))
//
//	msg := map[string]interface{}{
//		"id":     id,
//		"method": "Runtime.evaluate",
//		"params": map[string]interface{}{
//			"expression":    js,
//			"returnByValue": true,
//			"awaitPromise":  true,
//		},
//		"sessionId": sessionId,
//	}
//	err := conn.WriteJSON(msg)
//	if err != nil {
//		gt.Error("发送消息失败:", err)
//		return "", fmt.Errorf("发送消息失败")
//	}
//	//msgStr, _ := json.Marshal(msg)
//	//log.Printf("发送消息: %s", string(msgStr))
//
//	timeout := 6 * time.Second
//	timer := time.NewTimer(timeout)
//	defer timer.Stop() // 重要：确保计时器被清理
//
//	for {
//		select {
//		case msg, ok := <-messageQueue:
//			if !ok {
//				gt.Info("消息队列已关闭")
//				return "", fmt.Errorf("消息队列已关闭")
//			}
//			//gt.Info("收到的消息 -> ", msg.Content)
//			if id == msg.ID {
//				//gt.Info("是自己的消息")
//				return msg.Content, nil
//			} else {
//				gt.Info("不是自己的消息")
//			}
//
//		case <-timer.C:
//			gt.Info("6秒未收到消息")
//			return "", fmt.Errorf("接收消息超时; 6秒未收到消息")
//		}
//	}
//
//}
//
//func FocusAndEnterBtn(conn *websocket.Conn, sessionId string, id int, xpath string) (string, error) {
//	js := fmt.Sprintf(`
//				(function() {
//					const xpath = %s;
//					const element = document.evaluate(xpath, document, null,
//						XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
//
//					if (!element) return {success: false};
//
//					// 聚焦元素
//					element.focus();
//
//					// 触发focus事件
//					element.dispatchEvent(new Event('focus', { bubbles: true }));
//
//					// 模拟Enter键
//					const enterEvent = new KeyboardEvent('keydown', {
//						key: 'Enter',
//						code: 'Enter',
//						keyCode: 13,
//						charCode: 13,
//						which: 13,
//						bubbles: true,
//						cancelable: true
//					});
//
//					element.dispatchEvent(enterEvent);
//
//					// 触发keyup
//					const keyupEvent = new KeyboardEvent('keyup', {
//						key: 'Enter',
//						code: 'Enter',
//						keyCode: 13,
//						charCode: 13,
//						which: 13,
//						bubbles: true,
//						cancelable: true
//					});
//
//					element.dispatchEvent(keyupEvent);
//
//					// 最后点击一次
//					element.click();
//
//					return {success: true};
//				})()
//			`, jsonEscapeString(xpath))
//
//	msg := map[string]interface{}{
//		"id":     id,
//		"method": "Runtime.evaluate",
//		"params": map[string]interface{}{
//			"expression":    js,
//			"returnByValue": true,
//			"awaitPromise":  true,
//		},
//		"sessionId": sessionId,
//	}
//	err := conn.WriteJSON(msg)
//	if err != nil {
//		gt.Error("发送消息失败:", err)
//		return "", fmt.Errorf("发送消息失败")
//	}
//	//msgStr, _ := json.Marshal(msg)
//	//log.Printf("发送消息: %s", string(msgStr))
//
//	timeout := 6 * time.Second
//	timer := time.NewTimer(timeout)
//	defer timer.Stop() // 重要：确保计时器被清理
//
//	for {
//		select {
//		case msg, ok := <-messageQueue:
//			if !ok {
//				gt.Info("消息队列已关闭")
//				return "", fmt.Errorf("消息队列已关闭")
//			}
//			//gt.Info("收到的消息 -> ", msg.Content)
//			if id == msg.ID {
//				//gt.Info("是自己的消息")
//				return msg.Content, nil
//			} else {
//				gt.Info("不是自己的消息")
//			}
//
//		case <-timer.C:
//			gt.Info("6秒未收到消息")
//			return "", fmt.Errorf("接收消息超时; 6秒未收到消息")
//		}
//	}
//
//}

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
