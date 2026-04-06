package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Preload.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 预加载功能关闭: 停止Preload域的所有预加载、预连接、预解析监控
// 2. 测试环境清理: 自动化测试后关闭预加载监听，释放资源
// 3. 网络性能调试: 关闭预加载后对比页面加载性能差异
// 4. 资源监控停止: 不再接收预加载相关事件与状态通知
// 5. 异常恢复: 预加载监控异常时关闭并重置状态
// 6. 功能切换: 从预加载调试模式切换回正常运行模式

// CDPPreloadDisable 关闭Preload域，停止预加载相关事件监听
func CDPPreloadDisable() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Preload.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Preload.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("Preload.disable 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：预加载性能测试完成后关闭Preload域，释放浏览器资源
func ExamplePreloadDisable() {
	// 1. 执行预加载相关测试逻辑
	// RunPreloadTestCases()

	// 2. 关闭Preload域，停止监听
	resp, err := CDPPreloadDisable()
	if err != nil {
		log.Printf("关闭Preload域失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功关闭Preload域，已停止所有预加载事件监听")
}
*/

// -----------------------------------------------  Preload.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 预加载监控开启: 启动监听预加载、预连接、预解析、Speculative加载事件
// 2. 页面性能优化: 分析预加载策略是否生效、资源是否提前加载
// 3. 自动化测试: 验证页面预加载规则执行情况与加载效率
// 4. 网络调试: 追踪浏览器预加载行为、失败原因、加载状态
// 5. 前端资源优化: 基于预加载事件优化资源加载顺序
// 6. 调试辅助: 实时查看预加载触发流程与资源请求状态

// CDPPreloadEnable 启用Preload域，开始接收预加载相关事件
func CDPPreloadEnable() (string, error) {
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
		"method": "Preload.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Preload.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("Preload.enable 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：启动预加载监控，采集页面预加载、预连接行为用于性能分析
func ExamplePreloadEnable() {
	// 1. 建立浏览器CDP连接（前置逻辑）
	// err := ConnectBrowserCDP()
	// if err != nil {
	// 	log.Fatalf("浏览器连接失败: %v", err)
	// }

	// 2. 启用Preload域，开始监听预加载事件
	resp, err := CDPPreloadEnable()
	if err != nil {
		log.Printf("启用Preload域失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功启用Preload域，开始监听预加载、预解析、预连接事件")

	// 后续可监听 messageQueue 中的预加载相关事件，分析资源提前加载效果
}
*/
