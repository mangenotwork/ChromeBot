package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// ----------------------------------------------- BackgroundService.clearEvents  -----------------------------------------------
// === 应用场景 ===
// 1. 数据清理: 清理已记录的后台服务事件
// 2. 环境重置: 重置测试环境状态
// 3. 内存管理: 清理内存中的事件数据
// 4. 测试隔离: 隔离不同测试的数据
// 5. 调试清理: 清理调试过程中的事件
// 6. 性能优化: 优化内存使用

// CDPBackgroundServiceClearEvents 清理后台服务事件
// service: 服务类型
func CDPBackgroundServiceClearEvents(service string) (string, error) {
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
        "method": "BackgroundService.clearEvents",
        "params": {
            "service": "%s"
        }
    }`, reqID, service)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 BackgroundService.clearEvents 请求失败: %w", err)
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
			return "", fmt.Errorf("BackgroundService.clearEvents 请求超时")
		}
	}
}

// ----------------------------------------------- BackgroundService.setRecording  -----------------------------------------------
// === 应用场景 ===
// 1. 录制开始: 开始录制后台服务事件
// 2. 录制停止: 停止录制后台服务事件
// 3. 数据收集: 收集后台服务调试数据
// 4. 故障排查: 排查后台服务问题
// 5. 性能分析: 分析后台服务性能
// 6. 质量保证: 保证后台服务质量

// CDPBackgroundServiceSetRecording 设置后台服务录制状态
// shouldRecord: 是否录制
// service: 服务类型
func CDPBackgroundServiceSetRecording(shouldRecord bool, service string) (string, error) {
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
        "method": "BackgroundService.setRecording",
        "params": {
            "shouldRecord": %v,
            "service": "%s"
        }
    }`, reqID, shouldRecord, service)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 BackgroundService.setRecording 请求失败: %w", err)
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
			return "", fmt.Errorf("BackgroundService.setRecording 请求超时")
		}
	}
}

// ----------------------------------------------- BackgroundService.startObserving  -----------------------------------------------
// === 应用场景 ===
// 1. 监控启动: 开始监控后台服务活动
// 2. 性能监控: 监控后台服务的性能表现
// 3. 资源监控: 监控后台服务的资源使用情况
// 4. 调试辅助: 调试后台服务相关问题
// 5. 测试支持: 测试后台服务的行为
// 6. 安全审计: 审计后台服务的活动

// CDPBackgroundServiceStartObserving 开始观察后台服务
// service: 要观察的服务类型 ["backgroundFetch", "backgroundSync", "pushMessaging", "notifications", "paymentHandler", "periodicBackgroundSync"]
func CDPBackgroundServiceStartObserving(service string) (string, error) {
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
        "method": "BackgroundService.startObserving",
        "params": {
            "service": "%s"
        }
    }`, reqID, service)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 BackgroundService.startObserving 请求失败: %w", err)
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
			return "", fmt.Errorf("BackgroundService.startObserving 请求超时")
		}
	}
}

// ----------------------------------------------- BackgroundService.stopObserving  -----------------------------------------------
// === 应用场景 ===
// 1. 监控停止: 停止监控后台服务活动
// 2. 资源释放: 释放监控资源
// 3. 环境清理: 清理测试环境
// 4. 调试结束: 结束调试会话
// 5. 测试完成: 测试完成后停止监控
// 6. 性能优化: 优化运行时性能

// CDPBackgroundServiceStopObserving 停止观察后台服务
// service: 要停止观察的服务类型
func CDPBackgroundServiceStopObserving(service string) (string, error) {
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
        "method": "BackgroundService.stopObserving",
        "params": {
            "service": "%s"
        }
    }`, reqID, service)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 BackgroundService.stopObserving 请求失败: %w", err)
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
			return "", fmt.Errorf("BackgroundService.stopObserving 请求超时")
		}
	}
}
