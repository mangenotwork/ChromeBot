package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Inspector.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 关闭调试代理：关闭Inspector调试代理，停止接收调试事件
// 2. 资源释放：释放调试相关的连接、监听资源，降低浏览器开销
// 3. 自动化测试清理：测试流程结束后关闭调试功能，恢复浏览器原始状态
// 4. 连接关闭：主动断开调试连接，避免无效连接占用
// 5. 无头浏览器优化：长时间运行无头Chrome时关闭调试减少资源占用
// 6. 调试流程收尾：完成调试任务后优雅关闭Inspector模块

// CDPInspectorDisable 禁用Inspector调试代理
func CDPInspectorDisable() (string, error) {
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
		"method": "Inspector.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 disable 请求失败: %w", err)
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
			return "", fmt.Errorf("Inspector.disable 请求超时")
		}
	}
}

/*
// === 使用示例1：调试结束后关闭Inspector ===
func ExampleInspectorDisable_FinishDebug() {
	// 调试任务完成后关闭
	resp, err := CDPInspectorDisable()
	if err != nil {
		log.Fatalf("关闭Inspector失败: %v", err)
	}
	log.Println("Inspector已关闭，调试代理已停止:", resp)
}

// === 使用示例2：自动化测试清理 ===
func ExampleInspectorDisable_TestClean() {
	// 测试用例执行完毕，释放调试资源
	_, err := CDPInspectorDisable()
	if err != nil {
		log.Printf("清理Inspector失败: %v", err)
		return
	}
	log.Println("测试完成：Inspector调试已禁用")
}

// === 使用示例3：无头浏览器资源优化 ===
func ExampleInspectorDisable_Headless() {
	// 长期运行无头浏览器，关闭调试减少开销
	resp, err := CDPInspectorDisable()
	if err != nil {
		log.Fatalf("关闭调试代理失败: %v", err)
	}
	log.Println("已关闭Inspector，浏览器资源占用已降低:", resp)
}
*/

// -----------------------------------------------  Inspector.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 启动调试代理：开启Inspector调试代理，建立调试通信通道
// 2. 调试会话初始化：在开始调试前启用Inspector模块
// 3. 自动化测试前置：测试启动时开启调试功能，接收调试事件
// 4. 浏览器状态监控：启用后监听调试相关事件（断开、重启等）
// 5. 无头浏览器调试：无头模式下启用调试，支持远程控制与事件监听
// 6. 调试流程启动：作为CDP调试的基础入口方法，优先调用

// CDPInspectorEnable 启用Inspector调试代理
func CDPInspectorEnable() (string, error) {
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
		"method": "Inspector.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 enable 请求失败: %w", err)
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
			return "", fmt.Errorf("Inspector.enable 请求超时")
		}
	}
}

/*
// === 使用示例1：调试启动初始化 ===
func ExampleInspectorEnable_Init() {
	// 调试开始前必须启用
	resp, err := CDPInspectorEnable()
	if err != nil {
		log.Fatalf("启用Inspector失败: %v", err)
	}
	log.Println("Inspector已启用，调试通道已建立:", resp)
}

// === 使用示例2：自动化测试启动 ===
func ExampleInspectorEnable_TestStart() {
	// 测试任务启动时开启调试
	_, err := CDPInspectorEnable()
	if err != nil {
		log.Printf("测试前置启用Inspector失败: %v", err)
		return
	}
	log.Println("自动化测试：Inspector调试已就绪")
}

// === 使用示例3：无头浏览器调试初始化 ===
func ExampleInspectorEnable_Headless() {
	// 无头Chrome初始化调试能力
	resp, err := CDPInspectorEnable()
	if err != nil {
		log.Fatalf("无头浏览器启用Inspector失败: %v", err)
	}
	log.Println("无头Chrome调试已启用:", resp)
}

*/
