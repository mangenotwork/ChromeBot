package browser

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"html"
	"log"
	"strconv"
	"time"
)

//go:embed chrome_html.js
var chromeHtmlJS string

// GetHtml 获取页面的html
func (c *ChromeProcess) GetHtml() (string, error) {
	if c.NowTabWSConn == nil {
		c.DefaultNowTab()
	}

	c.NextID++
	msg := map[string]interface{}{
		"id":     c.NextID,
		"method": "Runtime.evaluate",
		"params": map[string]interface{}{
			"expression":    chromeHtmlJS,
			"returnByValue": true,
			"awaitPromise":  true,
		},
		"sessionId": c.NowTabSession,
	}
	err := c.NowTabWSConn.WriteJSON(msg)
	if err != nil {
		gt.Error("发送消息失败:", err)
		return "", fmt.Errorf("发送消息失败")
	}
	msgStr, _ := json.Marshal(msg)
	log.Printf("发送消息: %s", string(msgStr))

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

			//content, err := decodeUnicodeInHTML(msg.Content)
			//if err != nil {
			//	gt.Error("编码处理失败, err  = ", err)
			//	return "", fmt.Errorf("编码处理失败")
			//}
			//log.Println("收到的消息 -> ", content)
			if c.NextID == msg.ID {
				result, err := gt.Json2Map(msg.Content)
				if err != nil {
					log.Println("回复内容解析错误")
				} else {
					resultData, resultDataOK := result["result"]
					if resultDataOK {
						resultDataMap, hasMap := resultData.(map[string]any)
						if hasMap {
							resultBody, resultBodyOK := resultDataMap["result"]
							if resultBodyOK {
								resultBodyMap, resultBodyMapOK := resultBody.(map[string]any)
								if resultBodyMapOK {
									value, valueOK := resultBodyMap["value"].(map[string]any)
									if valueOK {
										htmlBody, htmlBodyOK := value["html"]
										if htmlBodyOK {
											htmlBodyStr, err := decodeUnicodeInHTML(htmlBody.(string))
											if err != nil {
												htmlBodyStr = htmlBody.(string)
											}
											return htmlBodyStr, nil
										} else {
											log.Println("未获取到页面html = ", value["error"])
											return "", fmt.Errorf("[Chrome] 未获取到页面html,err: %s", value["error"].(string))
										}
									}
								}
							}
						}
					}
					return "", fmt.Errorf("[Chrome] 未获取到页面html")
				}

			} else {
				log.Println("不是自己的消息")
			}

		case <-timer.C:
			log.Println("6秒未收到消息")
			return "", fmt.Errorf("接收消息超时; 6秒未收到消息")
		}
	}

}

// 从HTML中提取并转换\u编码
func decodeUnicodeInHTML(htmlStr string) (string, error) {
	return decodeHTMLWithUnicode(htmlStr)
}

// 解码HTML中的Unicode转义
func decodeHTMLWithUnicode(htmlStr string) (string, error) {
	// 先解码HTML实体（如 &lt; 等）
	decodedHTML := html.UnescapeString(htmlStr)

	// 提取并转换所有的\u编码
	result, err := decodeAllUnicodeEscapes(decodedHTML)
	if err != nil {
		return "", err
	}

	return result, nil
}

// 解码字符串中所有的Unicode转义序列
func decodeAllUnicodeEscapes(s string) (string, error) {
	var buf bytes.Buffer
	i := 0

	for i < len(s) {
		// 检查是否是\u转义序列
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == 'u' && i+6 <= len(s) {
			// 解析\uXXXX
			hexStr := s[i+2 : i+6]
			runeValue, err := strconv.ParseInt(hexStr, 16, 32)
			if err != nil {
				// 如果解析失败，保留原样
				buf.WriteString(s[i : i+6])
				i += 6
				continue
			}

			// 写入解码后的字符
			buf.WriteRune(rune(runeValue))
			i += 6
		} else if s[i] == '\\' && i+1 < len(s) && s[i+1] == 'U' && i+10 <= len(s) {
			// 解析\UXXXXXXXX (8位)
			hexStr := s[i+2 : i+10]
			runeValue, err := strconv.ParseInt(hexStr, 16, 32)
			if err != nil {
				buf.WriteString(s[i : i+10])
				i += 10
				continue
			}

			buf.WriteRune(rune(runeValue))
			i += 10
		} else {
			// 普通字符
			buf.WriteByte(s[i])
			i++
		}
	}

	return buf.String(), nil
}
