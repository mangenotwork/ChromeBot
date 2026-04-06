package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Media.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 媒体监听关闭: 停止接收所有媒体相关事件通知
// 2. 资源释放: 关闭媒体域后释放浏览器监听资源
// 3. 测试收尾: 自动化测试完成后关闭媒体事件监听
// 4. 性能优化: 无需媒体数据时关闭以减少性能消耗
// 5. 功能切换: 从媒体监听模式切换到其他功能模式
// 6. 异常恢复: 媒体监听异常时关闭并重新初始化

// CDPMediaDisable 关闭Media域，停止接收媒体相关事件
func CDPMediaDisable() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Media.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Media.disable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Media.disable 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景：自动化测试结束后关闭媒体监听
func ExampleMediaDisable() {
	// 1. 先启用媒体监听（业务逻辑）
	// resp, err := CDPMediaEnable()
	// if err != nil {
	// 	log.Fatalf("启用媒体监听失败: %v", err)
	// }

	// 2. 执行媒体相关测试逻辑...

	// 3. 测试完成后关闭媒体监听，释放资源
	resp, err := CDPMediaDisable()
	if err != nil {
		log.Printf("关闭媒体监听失败: %v, 响应内容: %s", err, resp)
		return
	}
	log.Println("成功关闭媒体域，停止接收媒体事件")
}

*/

// -----------------------------------------------  Media.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 媒体监听开启: 启动接收所有媒体相关事件通知
// 2. 媒体调试: 实时获取视频/音频播放状态、错误、日志信息
// 3. 自动化测试: 监听媒体播放异常、卡顿、加载失败等问题
// 4. 性能监控: 收集媒体播放性能数据、缓冲状态、解码信息
// 5. 业务监控: 监控页面音视频播放是否正常运行
// 6. 日志采集: 获取浏览器输出的媒体运行详细日志

// CDPMediaEnable 启用Media域，开始接收媒体相关事件
func CDPMediaEnable() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Media.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Media.enable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Media.enable 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：启动媒体监控，实时监听页面音视频状态
func ExampleMediaEnable() {
	// 1. 建立浏览器CDP连接（前置逻辑）
	// err := ConnectBrowserCDP()
	// if err != nil {
	// 	log.Fatalf("浏览器连接失败: %v", err)
	// }

	// 2. 启用媒体事件监听
	resp, err := CDPMediaEnable()
	if err != nil {
		log.Printf("启用媒体域失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功启用媒体域，开始接收视频/音频事件")

	// 后续可监听 messageQueue 中的 PlayerPropertiesChanged、PlayerErrors 等媒体事件
}
*/
