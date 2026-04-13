package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Security.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 安全检测关闭: 停止浏览器的安全相关事件监听和检测
// 2. 自动化测试收尾: 测试完成后关闭安全模块，释放资源
// 3. 调试环境恢复: 关闭安全监控，恢复浏览器默认安全状态
// 4. 无痕/临时调试: 临时关闭安全检查后恢复默认状态
// 5. 性能优化: 关闭不必要的安全监听减少浏览器开销

// CDPSecurityDisable 关闭Security域的功能，停止安全相关检测和事件
func CDPSecurityDisable() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Security.disable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Security.disable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Security.disable 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景：自动化测试完成后关闭安全模块，清理环境
func ExampleCDPSecurityDisable() {
	// 假设之前已经启用了Security模块监听证书/安全事件
	resp, err := CDPSecurityDisable()
	if err != nil {
		log.Fatalf("关闭Security失败: %v", err)
	}
	log.Printf("Security已关闭，响应: %s", resp)

	// 后续可执行浏览器关闭、其他模块重置等操作
}

// 场景：调试时临时关闭安全检测，之后恢复默认
func ExampleCDPSecurityDisable_Debug() {
	// 执行临时安全相关操作后关闭模块
	_, err := CDPSecurityDisable()
	if err != nil {
		log.Printf("关闭安全模块失败: %v", err)
		return
	}
	log.Println("已关闭浏览器安全检测，恢复默认状态")
}

*/

// -----------------------------------------------  Security.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 安全监控启动: 开启浏览器安全相关事件监听与检测
// 2. 证书验证监听: 开始接收证书、安全相关的事件通知
// 3. 自动化测试前置: 测试安全相关功能前启用Security模块
// 4. 网页安全调试: 调试HTTPS、证书错误时启用安全检测
// 5. 安全状态检测: 实时检测网页安全状态与证书信息

// CDPSecurityEnable 启用Security域功能，开启安全相关检测和事件监听
func CDPSecurityEnable() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Security.enable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Security.enable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Security.enable 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景：自动化测试开始前，启用安全模块监听证书与安全事件
func ExampleCDPSecurityEnable() {
	// 启动安全检测，准备监听证书、安全状态变更
	resp, err := CDPSecurityEnable()
	if err != nil {
		log.Fatalf("启用Security模块失败: %v", err)
	}
	log.Printf("Security已成功启用，响应: %s", resp)
}

// 场景：调试HTTPS网站时，启用安全模块获取证书信息
func ExampleCDPSecurityEnable_DebugHTTPS() {
	// 开启安全监控，用于捕获证书错误、安全状态
	_, err := CDPSecurityEnable()
	if err != nil {
		log.Printf("启用安全监控失败: %v", err)
		return
	}
	log.Println("已启用浏览器安全监控，可监听证书与安全状态")
}

*/

// -----------------------------------------------  Security.setIgnoreCertificateErrors  -----------------------------------------------
// === 应用场景 ===
// 1. 本地HTTPS调试：忽略自签名证书错误，正常访问本地开发环境
// 2. 自动化测试：跳过无效证书拦截，确保测试用例正常执行
// 3. 内网服务访问：忽略内网自签证书、过期证书提示
// 4. 抓包调试：配合Fiddler/Charles抓包时，忽略代理证书错误
// 5. 临时测试环境：测试环境证书无效时，临时关闭证书验证

// CDPSecuritySetIgnoreCertificateErrors 设置是否忽略证书错误
// 参数：ignore - true 忽略证书错误，false 启用正常证书验证
func CDPSecuritySetIgnoreCertificateErrors(ignore bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Security.setIgnoreCertificateErrors",
		"params": {
			"ignore": %t
		}
	}`, reqID, ignore)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Security.setIgnoreCertificateErrors 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Security.setIgnoreCertificateErrors 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景：本地开发调试HTTPS，忽略自签名证书错误
func ExampleCDPSecuritySetIgnoreCertificateErrors_LocalDebug() {
	// 忽略证书错误，访问本地自签HTTPS服务
	resp, err := CDPSecuritySetIgnoreCertificateErrors(true)
	if err != nil {
		log.Fatalf("设置忽略证书错误失败: %v", err)
	}
	log.Printf("已忽略证书错误，响应: %s", resp)
}

// 场景：测试完成后恢复正常证书验证
func ExampleCDPSecuritySetIgnoreCertificateErrors_Reset() {
	// 恢复浏览器默认严格证书验证
	resp, err := CDPSecuritySetIgnoreCertificateErrors(false)
	if err != nil {
		log.Printf("恢复证书验证失败: %v", err)
		return
	}
	log.Println("已恢复正常证书验证")
}

// 场景：自动化测试中跳过证书拦截
func ExampleCDPSecuritySetIgnoreCertificateErrors_AutoTest() {
	// 测试前置：开启忽略证书
	_, err := CDPSecuritySetIgnoreCertificateErrors(true)
	if err != nil {
		log.Fatalf("测试前置设置失败: %v", err)
	}
	// 执行测试用例...
	log.Println("自动化测试：已忽略证书错误，开始执行测试")
}

*/
