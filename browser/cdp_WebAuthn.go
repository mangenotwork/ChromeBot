package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  WebAuthn.addCredential  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境预置凭证: 自动化测试前预先注入WebAuthn凭证，无需真实硬件密钥
// 2. 模拟用户登录: 调试WebAuthn登录流程时，模拟已注册的用户凭证
// 3. 跨设备测试: 无需物理设备即可测试WebAuthn凭证的识别与使用
// 4. 异常场景测试: 注入特定格式/异常凭证，验证系统的容错能力
// 5. 开发调试: 前端开发时快速模拟WebAuthn凭证，无需真实注册流程
// 6. 批量测试数据: 批量注入测试凭证，覆盖多用户多凭证场景

// CDPWebAuthnAddCredential 添加WebAuthn凭证到虚拟认证器
// credential: WebAuthn凭证信息，包含id、publicKey、rpId等必要字段
func CDPWebAuthnAddCredential(credential interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 序列化凭证参数
	credentialBytes, err := json.Marshal(credential)
	if err != nil {
		return "", fmt.Errorf("序列化凭证参数失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "WebAuthn.addCredential",
		"params": {
			"credential": %s
		}
	}`, reqID, string(credentialBytes))

	// 发送WebSocket请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.addCredential 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.addCredential 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：基础测试场景 - 注入标准WebAuthn凭证用于登录测试
func ExampleCDPWebAuthnAddCredential_Basic() {
	// 构造符合WebAuthn规范的测试凭证
	testCredential := map[string]interface{}{
		"credentialId":  "AQIDBAUGBwgJCgsMDQ4PEA", // base64编码的凭证ID
		"publicKey":     "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE", // 公钥
		"rpId":          "example.com", // 依赖方ID（域名）
		"userHandle":    "dXNlcl9oYW5kbGUxMjM=", // 用户句柄
		"signCount":     0, // 签名计数
		"isResidentKey": false, // 是否为常驻密钥
	}

	// 调用添加凭证方法
	resp, err := CDPWebAuthnAddCredential(testCredential)
	if err != nil {
		log.Fatalf("添加WebAuthn凭证失败: %v", err)
	}
	log.Printf("添加凭证成功，响应: %s", resp)
}

// 示例2：调试场景 - 注入常驻密钥凭证，模拟无密码登录
func ExampleCDPWebAuthnAddCredential_Debug() {
	// 模拟无密码登录的常驻密钥凭证
	residentCredential := map[string]interface{}{
		"credentialId":  "MTIzNDU2Nzg5MDEyMzQ1Njc4OTA",
		"publicKey":     "MIIBIjANBgkqhkiG9w0BAQEFAAO",
		"rpId":          "test-app.com",
		"userHandle":    "dXNlcl9pZF85ODc2",
		"signCount":     5,
		"isResidentKey": true, // 开启常驻密钥
	}

	resp, err := CDPWebAuthnAddCredential(residentCredential)
	if err != nil {
		log.Printf("调试注入凭证失败: %v", err)
		return
	}
	log.Println("调试用常驻凭证注入成功")
}

*/

// -----------------------------------------------  WebAuthn.addVirtualAuthenticator  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试：无需物理U盾/指纹设备，创建虚拟认证器完成WebAuthn注册/登录全流程测试
// 2. 前端调试：本地开发调试WebAuthn功能，快速模拟不同类型认证器（平台/跨平台）
// 3. 兼容性测试：测试不同协议（CTAP1/U2F、CTAP2）、不同特性的认证器兼容性
// 4. 异常场景模拟：创建支持/不支持常驻密钥、用户验证的虚拟认证器，测试业务容错逻辑
// 5. 无硬件测试：在没有物理WebAuthn设备的环境下完成功能验证
// 6. 批量测试：批量创建虚拟认证器，模拟多设备、多用户认证场景

// CDPWebAuthnAddVirtualAuthenticator 创建虚拟WebAuthn认证器
// options: 虚拟认证器配置（协议、传输方式、是否支持常驻密钥、用户验证等）
func CDPWebAuthnAddVirtualAuthenticator(options interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 序列化配置参数
	optBytes, err := json.Marshal(options)
	if err != nil {
		return "", fmt.Errorf("序列化虚拟认证器配置失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "WebAuthn.addVirtualAuthenticator",
		"params": {
			"options": %s
		}
	}`, reqID, string(optBytes))

	// 发送WebSocket请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.addVirtualAuthenticator 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.addVirtualAuthenticator 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：基础测试场景 - 创建标准跨平台虚拟认证器（模拟U盾/安全密钥）
func ExampleCDPWebAuthnAddVirtualAuthenticator_Basic() {
	// 标准跨平台认证器配置（CTAP2 / USB 传输）
	authenticatorOpts := map[string]interface{}{
		"protocol":          "ctap2", // 认证器协议：ctap2 / u2f
		"transport":         "usb", // 传输方式：usb / nfc / ble / internal
		"hasResidentKey":     true, // 支持常驻密钥
		"hasUserVerification": true, // 支持用户验证（指纹/PIN）
		"isUserVerified":     true, // 默认已完成用户验证
	}

	// 创建虚拟认证器
	resp, err := CDPWebAuthnAddVirtualAuthenticator(authenticatorOpts)
	if err != nil {
		log.Fatalf("创建虚拟认证器失败: %v", err)
	}
	log.Printf("虚拟认证器创建成功，响应: %s", resp)
}

// 示例2：平台认证器场景 - 模拟设备内置认证器（Windows Hello/ Touch ID）
func ExampleCDPWebAuthnAddVirtualAuthenticator_Platform() {
	// 平台认证器（设备内置，无物理连接）
	platformOpts := map[string]interface{}{
		"protocol":          "ctap2",
		"transport":         "internal", // 内置传输 = 平台认证器
		"hasResidentKey":     true,
		"hasUserVerification": true,
		"isUserVerified":     true,
		"automaticPresenceSupport": true, // 自动模拟用户在场
	}

	resp, err := CDPWebAuthnAddVirtualAuthenticator(platformOpts)
	if err != nil {
		log.Printf("创建平台虚拟认证器失败: %v", err)
		return
	}
	log.Println("平台虚拟认证器（Touch ID/Windows Hello）创建成功")
}

*/

// -----------------------------------------------  WebAuthn.clearCredentials  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境重置: 自动化测试用例执行完毕后，清空虚拟认证器中的所有WebAuthn凭证
// 2. 数据隔离: 不同测试场景/用户之间清空凭证，避免数据污染
// 3. 调试重置: 调试WebAuthn登录/注册逻辑时，快速清空旧凭证重新测试
// 4. 会话清理: 用户登出后，清除虚拟认证器内的临时凭证
// 5. 批量测试重置: 批量测试流程中，循环清空凭证保证每次测试初始状态一致
// 6. 错误恢复: 凭证异常/失效时，清空所有凭证重新注入有效数据

// CDPWebAuthnClearCredentials 清空虚拟认证器中所有已添加的WebAuthn凭证
func CDPWebAuthnClearCredentials() (string, error) {
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
		"method": "WebAuthn.clearCredentials"
	}`, reqID)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.clearCredentials 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.clearCredentials 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：自动化测试重置 - 测试用例执行后清空凭证，保证下一次测试干净环境
func ExampleCDPWebAuthnClearCredentials_TestReset() {
	// 执行测试逻辑...
	// 测试完成后清空所有WebAuthn凭证
	resp, err := CDPWebAuthnClearCredentials()
	if err != nil {
		log.Fatalf("清空WebAuthn凭证失败: %v", err)
	}
	log.Printf("清空凭证成功，测试环境已重置: %s", resp)
}

// 示例2：调试场景重置 - 调试过程中清空旧凭证，重新注入新凭证测试
func ExampleCDPWebAuthnClearCredentials_Debug() {
	// 先清空历史无效凭证
	_, err := CDPWebAuthnClearCredentials()
	if err != nil {
		log.Printf("清空历史凭证失败: %v", err)
		return
	}
	log.Println("已清空所有旧WebAuthn凭证，可以重新注入新凭证")

	// 后续可重新调用 CDPWebAuthnAddCredential 注入新凭证
}

*/

// -----------------------------------------------  WebAuthn.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境还原: 自动化测试完成后关闭WebAuthn虚拟认证器功能，恢复浏览器原生状态
// 2. 功能切换: 关闭虚拟认证后，让浏览器使用真实的物理WebAuthn设备（U盾、指纹）
// 3. 多场景测试: 在同一流程中先模拟、后真实验证，切换时禁用虚拟认证
// 4. 资源释放: 关闭虚拟认证器，释放浏览器占用的相关资源
// 5. 错误恢复: WebAuthn模拟功能异常时，禁用后重新初始化
// 6. 兼容性验证: 关闭模拟后测试真实环境下WebAuthn表现

// CDPWebAuthnDisable 禁用WebAuthn虚拟认证器功能
func CDPWebAuthnDisable() (string, error) {
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
		"method": "WebAuthn.disable"
	}`, reqID)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.disable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.disable 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：测试完成后还原浏览器环境 - 关闭虚拟WebAuthn，恢复真实设备
func ExampleCDPWebAuthnDisable_TestReset() {
	// 执行完WebAuthn虚拟测试逻辑
	// 关闭虚拟认证器，还原浏览器原生状态
	resp, err := CDPWebAuthnDisable()
	if err != nil {
		log.Fatalf("禁用WebAuthn虚拟功能失败: %v", err)
	}
	log.Printf("已禁用WebAuthn虚拟认证器，浏览器恢复原生状态: %s", resp)
}

// 示例2：切换到真实物理设备 - 先禁用虚拟认证，再使用真实U盾/指纹登录
func ExampleCDPWebAuthnDisable_RealDevice() {
	// 禁用虚拟WebAuthn
	_, err := CDPWebAuthnDisable()
	if err != nil {
		log.Printf("禁用虚拟认证失败: %v", err)
		return
	}
	log.Println("已关闭虚拟WebAuthn，当前将使用真实物理认证设备")

	// 此时页面调用WebAuthn会触发真实的指纹/U盾设备，而非虚拟认证器
}

*/

// -----------------------------------------------  WebAuthn.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试前置：启动WebAuthn模拟功能，为后续虚拟认证器、凭证注入做准备
// 2. 前端开发调试：开启浏览器WebAuthn模拟，无需物理设备即可调试登录/注册流程
// 3. 测试环境初始化：自动化测试套件启动时初始化WebAuthn模拟能力
// 4. 功能恢复：禁用WebAuthn模拟后，重新启用继续进行模拟测试
// 5. 无硬件测试：在无物理U盾/指纹设备环境下，启用模拟完成功能验证
// 6. 异常恢复：模拟流程出错后，重新启用WebAuthn模块重置状态

// CDPWebAuthnEnable 启用WebAuthn虚拟认证器功能
func CDPWebAuthnEnable() (string, error) {
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
		"method": "WebAuthn.enable"
	}`, reqID)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.enable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.enable 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：测试初始化 - 启动自动化测试前启用WebAuthn模拟功能
func ExampleCDPWebAuthnEnable_TestInit() {
	// 初始化WebAuthn模拟模块（必须先执行enable才能使用其他WebAuthn方法）
	resp, err := CDPWebAuthnEnable()
	if err != nil {
		log.Fatalf("启用WebAuthn虚拟功能失败: %v", err)
	}
	log.Printf("WebAuthn模拟功能已启用，可创建虚拟认证器: %s", resp)

	// 后续可调用：CDPWebAuthnAddVirtualAuthenticator、CDPWebAuthnAddCredential等
}

// 示例2：调试恢复 - 禁用后重新启用WebAuthn模拟，继续调试登录流程
func ExampleCDPWebAuthnEnable_DebugRestore() {
	// 重新启用WebAuthn模拟
	_, err := CDPWebAuthnEnable()
	if err != nil {
		log.Printf("重新启用WebAuthn失败: %v", err)
		return
	}
	log.Println("WebAuthn模拟已重新启用，可继续调试注册/登录流程")

	// 调试场景：可立即创建虚拟认证器注入测试凭证
}

*/

// -----------------------------------------------  WebAuthn.getCredential  -----------------------------------------------
// === 应用场景 ===
// 1. 测试验证：获取已注入的WebAuthn凭证，校验凭证是否正确添加
// 2. 调试排查：获取凭证详情，排查登录/注册失败时的凭证配置问题
// 3. 数据校验：自动化测试中断言凭证属性（rpId、签名计数、公钥等）是否符合预期
// 4. 日志记录：获取凭证信息并记录，用于测试报告与问题追踪
// 5. 流程校验：确认虚拟认证器中存在指定凭证，再执行后续认证流程
// 6. 异常排查：获取凭证信息，分析认证失败是否由凭证不匹配导致

// CDPWebAuthnGetCredential 根据凭证ID获取WebAuthn凭证详情
// credentialId: base64url编码的凭证ID
func CDPWebAuthnGetCredential(credentialId string) (string, error) {
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
		"method": "WebAuthn.getCredential",
		"params": {
			"credentialId": "%s"
		}
	}`, reqID, credentialId)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.getCredential 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.getCredential 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：测试校验 - 获取凭证并验证是否正确注入
func ExampleCDPWebAuthnGetCredential_TestVerify() {
	// 已注入的凭证ID
	testCredentialId := "AQIDBAUGBwgJCgsMDQ4PEA"

	// 获取凭证详情
	resp, err := CDPWebAuthnGetCredential(testCredentialId)
	if err != nil {
		log.Fatalf("获取WebAuthn凭证失败: %v", err)
	}
	log.Printf("获取凭证成功，详情: %s", resp)
}

// 示例2：调试排查 - 登录失败时获取凭证检查配置
func ExampleCDPWebAuthnGetCredential_Debug() {
	// 登录流程异常，获取凭证排查rpId、userHandle等配置
	credId := "MTIzNDU2Nzg5MDEyMzQ1Njc4OTA"
	resp, err := CDPWebAuthnGetCredential(credId)
	if err != nil {
		log.Printf("获取凭证失败: %v", err)
		return
	}
	log.Println("凭证信息获取成功，可核对域名、公钥、用户句柄等参数")
}

*/

// -----------------------------------------------  WebAuthn.getCredentials  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境检查：获取虚拟认证器中所有WebAuthn凭证，校验注入的凭证数量与内容是否正确
// 2. 调试问题定位：登录/注册异常时，查看全部凭证列表排查缺失、重复、配置错误问题
// 3. 自动化断言：测试用例中断言凭证总数、rpId、用户句柄等是否符合预期结果
// 4. 测试报告生成：获取所有凭证信息记录到测试报告，便于追踪测试数据
// 5. 环境清理确认：清空凭证后调用，验证是否已彻底清空
// 6. 多凭证场景验证：批量注入后，获取列表确认所有凭证都已生效

// CDPWebAuthnGetCredentials 获取虚拟认证器中所有的WebAuthn凭证列表
func CDPWebAuthnGetCredentials() (string, error) {
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
		"method": "WebAuthn.getCredentials"
	}`, reqID)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.getCredentials 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.getCredentials 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：测试断言 - 获取所有凭证并校验数量是否正确
func ExampleCDPWebAuthnGetCredentials_TestVerify() {
	// 获取虚拟认证器中全部凭证
	resp, err := CDPWebAuthnGetCredentials()
	if err != nil {
		log.Fatalf("获取所有WebAuthn凭证失败: %v", err)
	}
	log.Printf("获取凭证列表成功: %s", resp)
}

// 示例2：调试环境检查 - 清空凭证后确认环境干净
func ExampleCDPWebAuthnGetCredentials_DebugCheck() {
	// 先清空所有凭证
	_, _ = CDPWebAuthnClearCredentials()

	// 获取列表确认已清空
	resp, err := CDPWebAuthnGetCredentials()
	if err != nil {
		log.Printf("检查凭证失败: %v", err)
		return
	}
	log.Println("当前虚拟认证器凭证列表：", resp)
	log.Println("可用于确认测试环境是否已重置干净")
}

*/

// -----------------------------------------------  WebAuthn.removeCredential  -----------------------------------------------
// === 应用场景 ===
// 1. 单凭证清理：测试中仅删除失效/错误的单个凭证，保留其他有效凭证
// 2. 测试数据维护：动态移除不需要的测试凭证，保持虚拟认证器数据整洁
// 3. 异常凭证处理：认证失败后，删除异常凭证并重新注入正确凭证
// 4. 多用户测试：切换测试用户时，移除上一个用户的专属凭证
// 5. 调试验证：删除指定凭证后，验证业务是否正确处理“凭证不存在”场景
// 6. 测试流程控制：按步骤删除指定凭证，验证渐进式认证流程

// CDPWebAuthnRemoveCredential 根据凭证ID删除单个WebAuthn凭证
// credentialId: base64url 编码的要删除的凭证ID
func CDPWebAuthnRemoveCredential(credentialId string) (string, error) {
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
		"method": "WebAuthn.removeCredential",
		"params": {
			"credentialId": "%s"
		}
	}`, reqID, credentialId)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.removeCredential 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.removeCredential 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：测试维护 - 删除失效的测试凭证，保留有效凭证
func ExampleCDPWebAuthnRemoveCredential_TestClean() {
	// 需要删除的失效凭证ID
	invalidCredentialId := "AQIDBAUGBwgJCgsMDQ4PEA"

	// 删除单个指定凭证
	resp, err := CDPWebAuthnRemoveCredential(invalidCredentialId)
	if err != nil {
		log.Fatalf("删除WebAuthn凭证失败: %v", err)
	}
	log.Printf("成功删除失效凭证，响应: %s", resp)
}

// 示例2：调试验证 - 删除凭证后测试业务“凭证不存在”的异常处理
func ExampleCDPWebAuthnRemoveCredential_Debug() {
	// 要移除的测试凭证ID
	credId := "MTIzNDU2Nzg5MDEyMzQ1Njc4OTA"

	// 执行删除
	_, err := CDPWebAuthnRemoveCredential(credId)
	if err != nil {
		log.Printf("删除凭证失败: %v", err)
		return
	}
	log.Println("指定凭证已删除，可测试页面访问时是否正确提示无可用凭证")
}

*/

// -----------------------------------------------  WebAuthn.removeVirtualAuthenticator  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境清理: 测试完成后删除虚拟认证器，释放浏览器资源
// 2. 多认证器切换: 移除旧的虚拟认证器，创建新的不同配置的认证器
// 3. 资源回收: 长时间运行自动化测试时，及时销毁不再使用的虚拟认证器
// 4. 状态重置: 虚拟认证器异常时，删除重建恢复正常状态
// 5. 场景隔离: 不同测试用例/用户之间独立认证器环境，避免数据干扰
// 6. 调试收尾: 调试结束后删除虚拟认证器，恢复浏览器原生WebAuthn行为

// CDPWebAuthnRemoveVirtualAuthenticator 删除指定的虚拟WebAuthn认证器
// authenticatorId: 要删除的虚拟认证器ID（通过addVirtualAuthenticator返回）
func CDPWebAuthnRemoveVirtualAuthenticator(authenticatorId string) (string, error) {
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
		"method": "WebAuthn.removeVirtualAuthenticator",
		"params": {
			"authenticatorId": "%s"
		}
	}`, reqID, authenticatorId)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.removeVirtualAuthenticator 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.removeVirtualAuthenticator 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：测试完成清理 - 自动化测试结束后删除虚拟认证器
func ExampleCDPWebAuthnRemoveVirtualAuthenticator_TestClean() {
	// 之前创建的虚拟认证器ID
	testAuthId := "auth_123456"

	// 删除虚拟认证器
	resp, err := CDPWebAuthnRemoveVirtualAuthenticator(testAuthId)
	if err != nil {
		log.Fatalf("删除虚拟认证器失败: %v", err)
	}
	log.Printf("成功删除虚拟认证器，测试环境已清理: %s", resp)
}

// 示例2：重建认证器 - 认证器异常时删除后重新创建
func ExampleCDPWebAuthnRemoveVirtualAuthenticator_Recreate() {
	// 异常的旧认证器ID
	oldAuthId := "auth_old_789"

	// 先删除异常认证器
	_, err := CDPWebAuthnRemoveVirtualAuthenticator(oldAuthId)
	if err != nil {
		log.Printf("删除旧认证器失败: %v", err)
		return
	}
	log.Println("已删除异常虚拟认证器，准备重新创建...")

	// 后续可重新调用 CDPWebAuthnAddVirtualAuthenticator 创建新认证器
}

*/

// -----------------------------------------------  WebAuthn.setAutomaticPresenceSimulation  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试：自动模拟用户在场验证，无需手动确认，实现无值守测试
// 2. 无交互认证：关闭弹窗/用户确认，让WebAuthn流程直接通过，提升测试效率
// 3. 持续集成：CI/CD环境中必须开启自动模拟，避免流程阻塞等待人工操作
// 4. 压力测试：批量并发认证时，自动模拟用户在场，保证流程顺畅
// 5. 调试跳过验证：调试时跳过指纹/PIN/触摸验证，快速进入核心逻辑
// 6. 场景切换：同一认证器在手动测试与自动化测试之间快速切换验证模式

// CDPWebAuthnSetAutomaticPresenceSimulation 为虚拟认证器设置是否自动模拟用户在场验证
// authenticatorId: 虚拟认证器ID
// enabled: true=自动模拟用户在场，false=需要用户手动验证
func CDPWebAuthnSetAutomaticPresenceSimulation(authenticatorId string, enabled bool) (string, error) {
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
		"method": "WebAuthn.setAutomaticPresenceSimulation",
		"params": {
			"authenticatorId": "%s",
			"enabled": %t
		}
	}`, reqID, authenticatorId, enabled)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.setAutomaticPresenceSimulation 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.setAutomaticPresenceSimulation 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：自动化测试/CI环境 - 开启自动模拟用户在场，无值守运行
func ExampleCDPWebAuthnSetAutomaticPresenceSimulation_AutoTest() {
	// 已创建的虚拟认证器ID
	authId := "auth_123456"

	// 开启自动用户在场模拟（无需触摸/指纹确认）
	resp, err := CDPWebAuthnSetAutomaticPresenceSimulation(authId, true)
	if err != nil {
		log.Fatalf("设置自动在场模拟失败: %v", err)
	}
	log.Printf("已开启自动用户在场模拟，自动化流程可无阻塞执行: %s", resp)
}

// 示例2：调试/手动测试 - 关闭自动模拟，恢复需要手动验证的真实场景
func ExampleCDPWebAuthnSetAutomaticPresenceSimulation_Manual() {
	authId := "auth_7890"

	// 关闭自动模拟，恢复真实场景：需要用户手动确认
	resp, err := CDPWebAuthnSetAutomaticPresenceSimulation(authId, false)
	if err != nil {
		log.Printf("关闭自动模拟失败: %v", err)
		return
	}
	log.Println("已关闭自动在场模拟，WebAuthn将等待手动验证（贴近真实用户体验）")
}

*/

// -----------------------------------------------  WebAuthn.setCredentialProperties  -----------------------------------------------
// === 应用场景 ===
// 1. 凭证状态修改：动态修改已注入WebAuthn凭证的属性，如更新常驻密钥状态
// 2. 测试场景模拟：修改凭证属性，测试不同凭证状态下的业务逻辑处理
// 3. 兼容性验证：修改用户验证、密钥类型等属性，验证系统兼容性
// 4. 异常流程测试：修改凭证属性为异常值，验证系统容错与错误处理能力
// 5. 无重新注入调试：无需删除重建凭证，直接修改属性快速调试
// 6. 多状态覆盖测试：一套凭证修改不同属性，覆盖多场景测试

// CDPWebAuthnSetCredentialProperties 设置WebAuthn凭证的属性
// credentialId: base64url编码的凭证ID
// properties: 凭证属性配置（可修改isResidentKey等）
func CDPWebAuthnSetCredentialProperties(credentialId string, properties interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 序列化属性参数
	propBytes, err := json.Marshal(properties)
	if err != nil {
		return "", fmt.Errorf("序列化凭证属性参数失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "WebAuthn.setCredentialProperties",
		"params": {
			"credentialId": "%s",
			"properties": %s
		}
	}`, reqID, credentialId, string(propBytes))

	// 发送WebSocket请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.setCredentialProperties 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.setCredentialProperties 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：基础修改 - 将普通凭证修改为常驻密钥凭证（无密码登录）
func ExampleCDPWebAuthnSetCredentialProperties_Basic() {
	// 目标凭证ID
	credId := "AQIDBAUGBwgJCgsMDQ4PEA"

	// 设置凭证属性：开启常驻密钥
	props := map[string]interface{}{
		"isResidentKey": true,
	}

	// 执行修改
	resp, err := CDPWebAuthnSetCredentialProperties(credId, props)
	if err != nil {
		log.Fatalf("设置凭证属性失败: %v", err)
	}
	log.Printf("凭证属性修改成功，响应: %s", resp)
}

// 示例2：测试场景 - 关闭常驻密钥，验证业务回退逻辑
func ExampleCDPWebAuthnSetCredentialProperties_Test() {
	credId := "MTIzNDU2Nzg5MDEyMzQ1Njc4OTA"

	// 关闭常驻密钥属性
	props := map[string]interface{}{
		"isResidentKey": false,
	}

	resp, err := CDPWebAuthnSetCredentialProperties(credId, props)
	if err != nil {
		log.Printf("修改凭证属性失败: %v", err)
		return
	}
	log.Println("已关闭凭证常驻密钥属性，可验证非常驻密钥登录流程")
}

*/

// -----------------------------------------------  WebAuthn.setResponseOverrideBits  -----------------------------------------------
// === 应用场景 ===
// 1. 异常测试：模拟WebAuthn认证器返回错误状态码，验证系统错误处理逻辑
// 2. 兼容性测试：覆盖不同错误响应场景，确保系统稳定不崩溃
// 3. 流程阻断测试：模拟认证中断、用户取消、设备异常等真实异常情况
// 4. 错误码覆盖：测试业务对各类CTAP错误码的识别与提示
// 5. 安全测试：验证系统对异常认证响应的防护与日志记录能力
// 6. 自动化异常用例：无需手动操作，自动模拟各类错误响应

// CDPWebAuthnSetResponseOverrideBits 为虚拟认证器设置WebAuthn响应覆盖位（模拟错误/异常响应）
// authenticatorId: 虚拟认证器ID
// bits: 覆盖配置，可设置errorCode等CTAP错误码
func CDPWebAuthnSetResponseOverrideBits(authenticatorId string, bits interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 序列化覆盖参数
	bitsBytes, err := json.Marshal(bits)
	if err != nil {
		return "", fmt.Errorf("序列化覆盖参数失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "WebAuthn.setResponseOverrideBits",
		"params": {
			"authenticatorId": "%s",
			"bits": %s
		}
	}`, reqID, authenticatorId, string(bitsBytes))

	// 发送WebSocket请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.setResponseOverrideBits 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.setResponseOverrideBits 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：模拟用户取消认证 - 模拟用户点击取消，验证系统提示逻辑
func ExampleCDPWebAuthnSetResponseOverrideBits_UserCancel() {
	authId := "auth_123456"

	// CTAP 错误码：0x27 = 操作被用户取消
	bits := map[string]interface{}{
		"errorCode": 0x27,
	}

	resp, err := CDPWebAuthnSetResponseOverrideBits(authId, bits)
	if err != nil {
		log.Fatalf("设置响应覆盖失败: %v", err)
	}
	log.Printf("已模拟用户取消认证，后续认证将返回取消错误: %s", resp)
}

// 示例2：模拟认证器超时/设备异常 - 测试系统异常处理与重连逻辑
func ExampleCDPWebAuthnSetResponseOverrideBits_DeviceError() {
	authId := "auth_7890"

	// CTAP 错误码：0x03 = 认证器超时/设备异常
	bits := map[string]interface{}{
		"errorCode": 0x03,
	}

	resp, err := CDPWebAuthnSetResponseOverrideBits(authId, bits)
	if err != nil {
		log.Printf("模拟设备异常失败: %v", err)
		return
	}
	log.Println("已模拟认证器设备异常，可测试系统错误提示与重试逻辑")
}


*/

// -----------------------------------------------  WebAuthn.setUserVerified  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试：模拟用户已完成验证（指纹/PIN），无需手动操作即可通过强认证
// 2. CI/CD 无值守测试：避免流程阻塞等待用户验证，保证自动化流程顺畅运行
// 3. 调试用户验证流程：快速切换用户验证状态，测试不同验证结果的业务处理
// 4. 权限测试：模拟已验证/未验证状态，测试系统权限控制与鉴权逻辑
// 5. 异常场景测试：设置未验证状态，测试系统对验证失败的容错与提示
// 6. 快速流程跳过：调试时跳过真实用户验证步骤，直接进入核心业务逻辑

// CDPWebAuthnSetUserVerified 设置虚拟认证器的用户验证状态
// authenticatorId: 虚拟认证器ID
// isUserVerified: true=用户已验证，false=用户未验证
func CDPWebAuthnSetUserVerified(authenticatorId string, isUserVerified bool) (string, error) {
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
		"method": "WebAuthn.setUserVerified",
		"params": {
			"authenticatorId": "%s",
			"isUserVerified": %t
		}
	}`, reqID, authenticatorId, isUserVerified)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAuthn.setUserVerified 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 监听响应消息队列
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID，获取对应响应
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应并检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAuthn.setUserVerified 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 示例1：自动化测试 - 模拟用户已验证，无值守通过WebAuthn强认证流程
func ExampleCDPWebAuthnSetUserVerified_AutoTest() {
	// 已创建的虚拟认证器ID
	authId := "auth_123456"

	// 设置为用户已验证（跳过指纹/PIN输入）
	resp, err := CDPWebAuthnSetUserVerified(authId, true)
	if err != nil {
		log.Fatalf("设置用户验证状态失败: %v", err)
	}
	log.Printf("已模拟用户验证通过，可直接完成认证流程: %s", resp)
}

// 示例2：异常测试 - 设置未验证状态，测试系统验证失败处理逻辑
func ExampleCDPWebAuthnSetUserVerified_ErrorTest() {
	authId := "auth_7890"

	// 设置用户未验证，模拟验证失败/取消
	resp, err := CDPWebAuthnSetUserVerified(authId, false)
	if err != nil {
		log.Printf("设置未验证状态失败: %v", err)
		return
	}
	log.Println("已模拟用户未验证，后续认证将返回验证失败，可测试错误提示")
}

*/
