package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Tethering.bind  -----------------------------------------------
// === 应用场景 ===
// 1. 绑定端口: 将设备指定端口绑定到浏览器，实现网络 tethering 通信
// 2. 远程调试: 绑定调试端口，用于远程设备调试
// 3. 网络代理: 绑定代理端口，实现设备网络通过浏览器转发
// 4. 自动化测试: 自动化流程中绑定设备通信端口
// 5. 设备连接: 建立浏览器与外部设备的端口通信通道
// 6. 调试环境初始化: 启动调试前先绑定所需端口

// CDPTetheringBind 绑定设备端口到浏览器，参数 port: 需要绑定的端口号
func CDPTetheringBind(port int) (string, error) {
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
		"method": "Tethering.bind",
		"params": {
			"port": %d
		}
	}`, reqID, port)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tethering.bind 请求失败: %w", err)
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
			return "", fmt.Errorf("Tethering.bind 请求超时")
		}
	}
}

/*


// === 使用示例 ===
// 示例1：绑定设备 9222 调试端口（最常用）
func ExampleCDPTetheringBind_DebugPort() {
	resp, err := CDPTetheringBind(9222)
	if err != nil {
		log.Fatalf("绑定调试端口失败: %v", err)
	}
	log.Printf("绑定成功: %s", resp)
}

// 示例2：绑定设备 8080 代理端口
func ExampleCDPTetheringBind_ProxyPort() {
	resp, err := CDPTetheringBind(8080)
	if err != nil {
		log.Fatalf("绑定代理端口失败: %v", err)
	}
	log.Printf("绑定成功: %s", resp)
}

// 示例3：自动化测试中绑定设备通信端口
func ExampleCDPTetheringBind_AutoTest() {
	// 测试前初始化绑定
	resp, err := CDPTetheringBind(10086)
	if err != nil {
		log.Printf("绑定测试端口失败: %v", err)
		return
	}
	log.Println("自动化测试端口绑定完成")
	// 后续执行测试逻辑
}

*/

// -----------------------------------------------  Tethering.unbind  -----------------------------------------------
// === 应用场景 ===
// 1. 解绑端口: 解除设备指定端口与浏览器的绑定关系
// 2. 资源释放: 测试/调试完成后释放占用的端口资源
// 3. 连接关闭: 主动关闭浏览器与外部设备的端口通信通道
// 4. 环境重置: 自动化测试后重置 tethering 连接状态
// 5. 错误处理: 端口通信异常时解绑并重新绑定
// 6. 程序退出: 程序关闭前优雅解绑所有绑定端口

// CDPTetheringUnbind 解除设备端口与浏览器的绑定，参数 port: 需要解绑的端口号
func CDPTetheringUnbind(port int) (string, error) {
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
		"method": "Tethering.unbind",
		"params": {
			"port": %d
		}
	}`, reqID, port)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tethering.unbind 请求失败: %w", err)
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
			return "", fmt.Errorf("Tethering.unbind 请求超时")
		}
	}
}

/*


// === 使用示例 ===
// 示例1：解绑设备 9222 调试端口
func ExampleCDPTetheringUnbind_DebugPort() {
	resp, err := CDPTetheringUnbind(9222)
	if err != nil {
		log.Fatalf("解绑调试端口失败: %v", err)
	}
	log.Printf("解绑成功: %s", resp)
}

// 示例2：解绑设备 8080 代理端口
func ExampleCDPTetheringUnbind_ProxyPort() {
	resp, err := CDPTetheringUnbind(8080)
	if err != nil {
		log.Fatalf("解绑代理端口失败: %v", err)
	}
	log.Printf("解绑成功: %s", resp)
}

// 示例3：自动化测试完成后解绑端口并释放资源
func ExampleCDPTetheringUnbind_AutoTestCleanup() {
	// 测试执行完成后解绑
	resp, err := CDPTetheringUnbind(10086)
	if err != nil {
		log.Printf("解绑测试端口失败: %v", err)
		return
	}
	log.Println("自动化测试端口解绑完成，资源已释放")
}

*/
