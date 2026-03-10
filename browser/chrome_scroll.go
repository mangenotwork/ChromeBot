package browser

import (
	"ChromeBot/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"strconv"
	"strings"
	"time"
)

//go:embed chrome_scroll_pixel.js
var chromeScrollPixelJS string

//go:embed chrome_scroll_element.js
var chromeScrollElementJS string

// ScrollPixelParam 按像素滚动的参数
type ScrollPixelParam struct {
	X        int  `json:"x"`        // 水平滚动到的位置（像素）
	Y        int  `json:"y"`        // 垂直滚动到的位置（像素）
	IsSmooth bool `json:"isSmooth"` // 是否平滑滚动
}

// ScrollElementParam 按元素滚动的参数
type ScrollElementParam struct {
	XPath    string `json:"xpath"`    // 目标元素的XPath
	IsSmooth bool   `json:"isSmooth"` // 是否平滑滚动
}

// ScrollByPixel 按像素滚动
func (c *ChromeProcess) ScrollByPixel(x, y int) error {
	jsPixel := strings.ReplaceAll(chromeScrollPixelJS, "__SCROLL_X__", strconv.Itoa(x))
	jsPixel = strings.ReplaceAll(jsPixel, "__SCROLL_Y__", strconv.Itoa(y))
	res, err := c.scroll(jsPixel)
	log.Printf("[Chrome]滚动结果: %v", res)
	return err
}

// ScrollToElement 按元素滚动
func (c *ChromeProcess) ScrollToElement(xPath string) error {
	xPath = "'" + strings.ReplaceAll(xPath, "\"", "\\\"") + "'"
	jsElement := strings.ReplaceAll(chromeScrollElementJS, "__SCROLL_XPATH__", xPath)
	jsElement = strings.ReplaceAll(jsElement, "__SCROLL_IS_SMOOTH__", strconv.FormatBool(true))
	res, err := c.scroll(jsElement)
	log.Printf("[Chrome]滚动结果: %v", res)
	return err
}

// ScrollResult 滚动结果（对应JS返回的结构）
type ScrollResult struct {
	Success bool   `json:"success"` // 是否成功
	Error   string `json:"error"`   // 错误类型（如参数错误/元素未找到）
	Message string `json:"message"` // 详细错误/成功信息
	Stack   string `json:"stack"`   // 可选：异常堆栈
}

func (c *ChromeProcess) scroll(js string) (*ScrollResult, error) {
	if c.NowTabWSConn == nil {
		c.DefaultNowTab()
	}

	c.NextID++
	msg := map[string]interface{}{
		"id":     c.NextID,
		"method": "Runtime.evaluate",
		"params": map[string]interface{}{
			"expression":    js,
			"returnByValue": true, // 必须：返回完整的对象结构
			"awaitPromise":  false,
		},
		"sessionId": c.NowTabSession,
	}

	err := c.NowTabWSConn.WriteJSON(msg)
	if err != nil {
		return nil, fmt.Errorf("发送滚动消息失败: %w", err)
	}
	msgStr, _ := json.Marshal(msg)
	utils.Debugf("发送滚动消息: %s", string(msgStr))

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列已关闭")
			}
			utils.Debug("收到滚动响应 -> ", respMsg.Content)

			// 匹配当前请求的响应
			if c.NextID == respMsg.ID {

				resultJson, err := gt.JsonFind(respMsg.Content, "/result/result/value")
				if err != nil {
					return nil, fmt.Errorf("解析滚动结果失败: %w", err)
				}

				resultBytes, err := json.Marshal(resultJson)
				if err != nil {
					return nil, fmt.Errorf("转换结果为JSON失败: %w", err)
				}
				var scrollResult ScrollResult
				err = json.Unmarshal(resultBytes, &scrollResult)
				if err != nil {
					return nil, fmt.Errorf("反序列化滚动结果失败: %w", err)
				}

				if !scrollResult.Success {
					return &scrollResult, fmt.Errorf("滚动执行失败: %s - %s", scrollResult.Error, scrollResult.Message)
				}
				return &scrollResult, nil
			} else {
				log.Println("不是当前滚动请求的响应，忽略")
			}

		case <-timer.C:
			return nil, fmt.Errorf("滚动请求超时（6秒）")
		}
	}
}
