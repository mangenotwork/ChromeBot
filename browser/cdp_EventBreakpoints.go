package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  EventBreakpoints.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 调试结束清理: 完成调试会话后清理所有事件断点
// 2. 错误恢复: 在调试过程中发生错误时恢复浏览器正常运行状态
// 3. 自动化测试: 在自动化测试结束后清理测试设置的事件断点
// 4. 性能优化: 清除事件断点以减少性能开销，特别是在长时间运行的页面中
// 5. 页面刷新: 在页面刷新前清理断点，避免断点残留
// 6. 多页面调试: 切换到不同页面时清理当前页面的断点设置
// 7. 安全恢复: 在安全敏感的调试环境中确保不留下调试痕迹
// 8. 内存管理: 释放事件断点占用的资源，防止内存泄漏
// 9. 断点重置: 在需要重置调试状态时清除所有事件断点
// 10. 批量操作: 一次性地移除通过 EventBreakpoints 设置的所有断点，而不是逐个移除

// EventBreakpointsDisable 移除所有通过 EventBreakpoints 设置的事件断点
// EventBreakpoints.disable
// 参数: 无
// 返回值: CDP响应内容字符串和错误信息
func EventBreakpointsDisable() (string, error) {
	// 检查WebSocket连接状态
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 生成请求ID
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息
	// EventBreakpoints.disable 方法不需要参数
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "EventBreakpoints.disable"
	}`, reqID)

	log.Printf("[DEBUG] EventBreakpointsDisable: 发送请求 ID=%d", reqID)
	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 EventBreakpoints.disable 请求失败: %w", err)
	}

	// 等待响应，设置5秒超时
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听消息队列获取响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 检查是否是对应请求ID的响应
			if reqID == respMsg.ID {
				// 格式化响应内容
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] EventBreakpointsDisable: 收到响应 ID=%d", reqID)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("CDP错误: %v", errorObj)
					return content, fmt.Errorf(errorMsg)
				}

				// 记录成功日志
				log.Printf("[INFO] EventBreakpoints.disable 成功: 已移除所有事件断点")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("EventBreakpoints.disable 请求超时 (%v)", timeout)
		}
	}
}

// -----------------------------------------------  EventBreakpoints.removeInstrumentationBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 选择性断点清理: 移除特定的不再需要的事件断点，而保留其他断点
// 2. 事件断点管理: 动态调整事件断点，在调试过程中移除已分析的特定事件
// 3. 性能调优: 移除对性能影响较大的特定事件断点
// 4. 调试流程控制: 在不同调试阶段移除不再需要的事件断点
// 5. 事件过滤: 从多个事件断点中移除特定的一个，实现更精细的控制
// 6. 条件断点移除: 当特定条件满足时移除对应的事件断点
// 7. 资源释放: 释放不再需要的特定事件断点占用的资源
// 8. 断点生命周期管理: 实现断点的精确创建和移除

// EventBreakpointsRemoveInstrumentationBreakpoint 移除特定的原生事件断点
// EventBreakpoints.removeInstrumentationBreakpoint
// 参数: eventName - 要移除的instrumentation事件名称
// 返回值: CDP响应内容字符串和错误信息
func EventBreakpointsRemoveInstrumentationBreakpoint(eventName string) (string, error) {
	// 验证参数
	if eventName == "" {
		return "", fmt.Errorf("eventName 参数不能为空")
	}

	// 检查WebSocket连接状态
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 生成请求ID
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "EventBreakpoints.removeInstrumentationBreakpoint",
		"params": {
			"eventName": "%s"
		}
	}`, reqID, eventName)

	log.Printf("[DEBUG] EventBreakpointsRemoveInstrumentationBreakpoint: 发送请求 ID=%d, eventName=%s", reqID, eventName)
	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 EventBreakpoints.removeInstrumentationBreakpoint 请求失败: %w", err)
	}

	// 等待响应，设置5秒超时
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听消息队列获取响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 检查是否是对应请求ID的响应
			if reqID == respMsg.ID {
				// 格式化响应内容
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] EventBreakpointsRemoveInstrumentationBreakpoint: 收到响应 ID=%d", reqID)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("CDP错误: %v", errorObj)
					return content, fmt.Errorf(errorMsg)
				}

				// 记录成功日志
				log.Printf("[INFO] EventBreakpoints.removeInstrumentationBreakpoint 成功: 已移除事件断点 '%s'", eventName)
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("EventBreakpoints.removeInstrumentationBreakpoint 请求超时 (%v)", timeout)
		}
	}
}

// -----------------------------------------------  EventBreakpoints.setInstrumentationBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 事件监听调试: 在特定JavaScript事件上设置断点，如点击、键盘、网络事件等
// 2. 性能分析: 在关键原生事件上设置断点，分析性能瓶颈
// 3. 事件流调试: 调试复杂的事件流和事件传播过程
// 4. 第三方库分析: 分析第三方库的事件处理逻辑
// 5. 事件处理顺序: 调试事件处理的顺序和优先级
// 6. 事件冒泡/捕获: 分析事件冒泡和捕获阶段的行为
// 7. 自定义事件: 调试自定义事件的触发和处理
// 8. 事件委托: 分析事件委托模式下的行为
// 9. 异步事件: 调试异步触发的事件
// 10. 内存泄漏检测: 在事件监听器上设置断点，检测可能的内存泄漏

// EventBreakpointsSetInstrumentationBreakpoint 在特定的原生事件上设置断点
// EventBreakpoints.setInstrumentationBreakpoint
// 参数: eventName - 要设置断点的instrumentation事件名称
// 返回值: CDP响应内容字符串和错误信息
func EventBreakpointsSetInstrumentationBreakpoint(eventName string) (string, error) {
	// 验证参数
	if eventName == "" {
		return "", fmt.Errorf("eventName 参数不能为空")
	}

	// 检查WebSocket连接状态
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 生成请求ID
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "EventBreakpoints.setInstrumentationBreakpoint",
		"params": {
			"eventName": "%s"
		}
	}`, reqID, eventName)

	log.Printf("[DEBUG] EventBreakpointsSetInstrumentationBreakpoint: 发送请求 ID=%d, eventName=%s", reqID, eventName)
	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 EventBreakpoints.setInstrumentationBreakpoint 请求失败: %w", err)
	}

	// 等待响应，设置5秒超时
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听消息队列获取响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 检查是否是对应请求ID的响应
			if reqID == respMsg.ID {
				// 格式化响应内容
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] EventBreakpointsSetInstrumentationBreakpoint: 收到响应 ID=%d", reqID)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("CDP错误: %v", errorObj)
					return content, fmt.Errorf(errorMsg)
				}

				// 记录成功日志
				log.Printf("[INFO] EventBreakpoints.setInstrumentationBreakpoint 成功: 已设置事件断点 '%s'", eventName)
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("EventBreakpoints.setInstrumentationBreakpoint 请求超时 (%v)", timeout)
		}
	}
}
