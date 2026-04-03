package browser

import (
	"ChromeBot/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Fetch.continueRequest  -----------------------------------------------
// === 应用场景 ===
// 1. 请求继续: 继续被拦截的请求
// 2. 请求修改: 修改请求后继续
// 3. 请求重放: 重放修改后的请求
// 4. 流量控制: 控制请求流量
// 5. 调试支持: 调试请求流程
// 6. 测试验证: 验证请求处理逻辑

// CDPFetchContinueRequest 继续被拦截的请求
func CDPFetchContinueRequest(requestID string, modifications *ContinueRequestModifications) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建修改参数
	var modificationsJSON string
	if modifications != nil {
		modsBytes, err := json.Marshal(modifications)
		if err != nil {
			return "", fmt.Errorf("序列化修改参数失败: %w", err)
		}
		modificationsJSON = fmt.Sprintf(`"modifications": %s,`, string(modsBytes))
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Fetch.continueRequest",
        "params": {
            "requestId": "%s",
            %s
            "interceptResponse": false
        }
    }`, reqID, requestID, modificationsJSON)

	// 移除可能的多余逗号
	message = strings.ReplaceAll(message, ",,", ",")
	message = strings.ReplaceAll(message, ",\n        }", "\n        }")

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.continueRequest 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.continueRequest 请求超时")
		}
	}
}

// ContinueRequestModifications 继续请求的修改参数
type ContinueRequestModifications struct {
	URL               string   `json:"url,omitempty"`
	Method            string   `json:"method,omitempty"`
	PostData          string   `json:"postData,omitempty"`
	Headers           []Header `json:"headers,omitempty"`
	InterceptResponse bool     `json:"interceptResponse,omitempty"`
}

// Header 请求头结构
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// -----------------------------------------------  Fetch.continueWithAuth  -----------------------------------------------
// === 应用场景 ===
// 1. 认证处理: 处理HTTP认证挑战
// 2. 自动登录: 自动提供认证凭据
// 3. 测试认证: 测试认证流程
// 4. 安全测试: 测试安全认证机制
// 5. 代理认证: 处理代理认证
// 6. 单点登录: 测试单点登录流程

// CDPFetchContinueWithAuth 继续带有认证的请求
// requestID: 请求ID
// authChallengeResponse: 认证挑战响应
func CDPFetchContinueWithAuth(requestID string, authChallengeResponse *AuthChallengeResponse) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建认证挑战响应
	authResponseBytes, err := json.Marshal(authChallengeResponse)
	if err != nil {
		return "", fmt.Errorf("序列化认证响应失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Fetch.continueWithAuth",
        "params": {
            "requestId": "%s",
            "authChallengeResponse": %s
        }
    }`, reqID, requestID, string(authResponseBytes))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.continueWithAuth 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.continueWithAuth 请求超时")
		}
	}
}

// AuthChallengeResponse 认证挑战响应结构
type AuthChallengeResponse struct {
	Response string `json:"response"` // "Default", "CancelAuth", "ProvideCredentials"
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// 辅助函数: CDPFetchContinueWithCredentials 使用凭据继续认证
func CDPFetchContinueWithCredentials(requestID, username, password string) (string, error) {
	authResponse := &AuthChallengeResponse{
		Response: "ProvideCredentials",
		Username: username,
		Password: password,
	}

	return CDPFetchContinueWithAuth(requestID, authResponse)
}

// 辅助函数: CDPFetchCancelAuth 取消认证
func CDPFetchCancelAuth(requestID string) (string, error) {
	authResponse := &AuthChallengeResponse{
		Response: "CancelAuth",
	}

	return CDPFetchContinueWithAuth(requestID, authResponse)
}

// 辅助函数: CDPFetchUseDefaultAuth 使用默认认证
func CDPFetchUseDefaultAuth(requestID string) (string, error) {
	authResponse := &AuthChallengeResponse{
		Response: "Default",
	}

	return CDPFetchContinueWithAuth(requestID, authResponse)
}

// -----------------------------------------------  Fetch.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 拦截器清理: 清理网络请求拦截器
// 2. 环境重置: 重置网络拦截状态
// 3. 资源释放: 释放拦截器相关资源
// 4. 测试结束: 网络测试完成后清理
// 5. 错误恢复: 网络拦截异常时恢复
// 6. 功能切换: 切换网络监控状态

// CDPFetchDisable 禁用Fetch拦截
func CDPFetchDisable() (string, error) {
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
        "method": "Fetch.disable"
    }`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.disable 请求超时")
		}
	}
}

// -----------------------------------------------  Fetch.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 网络监控: 启用网络请求监控
// 2. 请求拦截: 拦截和处理网络请求
// 3. 性能分析: 分析网络请求性能
// 4. 调试支持: 调试网络相关问题
// 5. 安全测试: 测试网络安全策略
// 6. 自动化测试: 自动化网络请求测试

// CDPFetchEnable 启用Fetch拦截
// patterns: 拦截模式数组
// handleAuthRequests: 是否处理认证请求
func CDPFetchEnable(patterns []FetchPattern, handleAuthRequests bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建拦截模式数组
	var patternsJSON string
	if len(patterns) > 0 {
		patternsBytes, err := json.Marshal(patterns)
		if err != nil {
			return "", fmt.Errorf("序列化拦截模式失败: %w", err)
		}
		patternsJSON = string(patternsBytes)
	} else {
		patternsJSON = "[]"
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Fetch.enable",
        "params": {
            "patterns": %s,
            "handleAuthRequests": %v
        }
    }`, reqID, patternsJSON, handleAuthRequests)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.enable 请求超时")
		}
	}
}

// FetchPattern 拦截模式结构
type FetchPattern struct {
	URLPattern   string `json:"urlPattern,omitempty"`
	ResourceType string `json:"resourceType,omitempty"` // "Document", "Stylesheet", "Image", "Media", "Font", "Script", "TextTrack", "XHR", "Fetch", "EventSource", "WebSocket", "Manifest", "SignedExchange", "Ping", "CSPViolationReport", "Other"
	RequestStage string `json:"requestStage,omitempty"` // "Request", "Response"
}

// -----------------------------------------------  Fetch.failRequest  -----------------------------------------------
// === 应用场景 ===
// 1. 错误模拟: 模拟请求失败
// 2. 故障测试: 测试故障处理
// 3. 超时测试: 测试请求超时
// 4. 网络异常: 模拟网络异常
// 5. 安全测试: 测试错误处理
// 6. 容错测试: 测试系统容错性

// CDPFetchFailRequest 使被拦截的请求失败
func CDPFetchFailRequest(requestID, errorReason string) (string, error) {
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
        "method": "Fetch.failRequest",
        "params": {
            "requestId": "%s",
            "errorReason": "%s"
        }
    }`, reqID, requestID, errorReason)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.failRequest 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.failRequest 请求超时")
		}
	}
}

// -----------------------------------------------  Fetch.fulfillRequest  -----------------------------------------------
// === 应用场景 ===
// 1. 模拟响应: 模拟请求的响应
// 2. 测试数据: 提供测试数据响应
// 3. 错误模拟: 模拟错误响应
// 4. 缓存测试: 测试缓存行为
// 5. 性能测试: 测试响应性能
// 6. 离线测试: 模拟离线场景

// CDPFetchFulfillRequest 完成被拦截的请求
func CDPFetchFulfillRequest(requestID string, response *FulfillRequestResponse) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建响应参数
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("序列化响应参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Fetch.fulfillRequest",
        "params": {
            "requestId": "%s",
            "response": %s
        }
    }`, reqID, requestID, string(responseBytes))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.fulfillRequest 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.fulfillRequest 请求超时")
		}
	}
}

// FulfillRequestResponse 完成请求的响应参数
type FulfillRequestResponse struct {
	ResponseCode          int      `json:"responseCode"`
	ResponseHeaders       []Header `json:"responseHeaders,omitempty"`
	BinaryResponseHeaders string   `json:"binaryResponseHeaders,omitempty"`
	Body                  string   `json:"body,omitempty"`
	ResponsePhrase        string   `json:"responsePhrase,omitempty"`
}

// -----------------------------------------------  Fetch.getResponseBody  -----------------------------------------------
// === 应用场景 ===
// 1. 响应体获取: 获取被拦截请求的响应体
// 2. 数据验证: 验证API响应数据
// 3. 性能分析: 分析响应体大小和内容
// 4. 调试辅助: 调试网络请求响应
// 5. 安全审计: 审计敏感数据泄露
// 6. 数据备份: 备份网络响应数据

// CDPFetchGetResponseBody 获取请求的响应体
// requestID: 请求ID
func CDPFetchGetResponseBody(requestID string) (string, error) {
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
        "method": "Fetch.getResponseBody",
        "params": {
            "requestId": "%s"
        }
    }`, reqID, requestID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.getResponseBody 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.getResponseBody 请求超时")
		}
	}
}

// ResponseBodyInfo 响应体信息结构
type ResponseBodyInfo struct {
	Body          string `json:"body"`
	Base64Encoded bool   `json:"base64Encoded"`
	Size          int    `json:"size"`
	ContentType   string `json:"contentType,omitempty"`
}

// 辅助函数: CDPFetchGetResponseBodyParsed 获取并解析响应体
func CDPFetchGetResponseBodyParsed(requestID string) (*ResponseBodyInfo, error) {
	result, err := CDPFetchGetResponseBody(requestID)
	if err != nil {
		return nil, fmt.Errorf("获取响应体失败: %w", err)
	}

	var response struct {
		Result struct {
			Body          string `json:"body"`
			Base64Encoded bool   `json:"base64Encoded"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return nil, fmt.Errorf("解析响应体失败: %w", err)
	}

	info := &ResponseBodyInfo{
		Body:          response.Result.Body,
		Base64Encoded: response.Result.Base64Encoded,
		Size:          len(response.Result.Body),
	}

	return info, nil
}

// 辅助函数: CDPFetchGetResponseBodyAsJSON 获取响应体并解析为JSON
func CDPFetchGetResponseBodyAsJSON(requestID string) (map[string]interface{}, error) {
	info, err := CDPFetchGetResponseBodyParsed(requestID)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	// 如果是Base64编码，先解码
	body := info.Body
	if info.Base64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, fmt.Errorf("Base64解码失败: %w", err)
		}
		body = string(decoded)
	}

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return data, nil
}

// 辅助函数: CDPFetchGetResponseBodyAsString 获取响应体为字符串
func CDPFetchGetResponseBodyAsString(requestID string) (string, error) {
	info, err := CDPFetchGetResponseBodyParsed(requestID)
	if err != nil {
		return "", err
	}

	if info.Base64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(info.Body)
		if err != nil {
			return "", fmt.Errorf("Base64解码失败: %w", err)
		}
		return string(decoded), nil
	}

	return info.Body, nil
}

/*

// 示例1: 获取并分析API响应
func AnalyzeAPIResponse(requestID string) error {
    log.Printf("分析API响应，请求ID: %s", requestID)

    // 获取响应体
    info, err := CDPFetchGetResponseBodyParsed(requestID)
    if err != nil {
        return fmt.Errorf("获取响应体失败: %w", err)
    }

    fmt.Println("=== API响应分析 ===")
    fmt.Printf("请求ID: %s\n", requestID)
    fmt.Printf("响应体大小: %d 字节\n", info.Size)
    fmt.Printf("Base64编码: %v\n", info.Base64Encoded)

    // 如果是Base64编码，解码
    var bodyStr string
    if info.Base64Encoded {
        decoded, err := base64.StdEncoding.DecodeString(info.Body)
        if err != nil {
            return fmt.Errorf("Base64解码失败: %w", err)
        }
        bodyStr = string(decoded)
    } else {
        bodyStr = info.Body
    }

    // 尝试解析为JSON
    var jsonData interface{}
    if err := json.Unmarshal([]byte(bodyStr), &jsonData); err == nil {
        fmt.Println("格式: JSON")

        // 如果是JSON对象，显示键
        if obj, ok := jsonData.(map[string]interface{}); ok {
            fmt.Printf("JSON键数量: %d\n", len(obj))
            fmt.Println("JSON键:")
            for key := range obj {
                fmt.Printf("  - %s\n", key)
            }
        }

        // 如果是JSON数组
        if arr, ok := jsonData.([]interface{}); ok {
            fmt.Printf("JSON数组长度: %d\n", len(arr))
        }
    } else if strings.Contains(strings.ToLower(info.Body), "<html") {
        fmt.Println("格式: HTML")
        fmt.Printf("HTML大小: %d 字符\n", len(bodyStr))
    } else if strings.Contains(strings.ToLower(info.Body), "<?xml") {
        fmt.Println("格式: XML")
    } else {
        fmt.Println("格式: 纯文本/未知")
    }

    // 显示前200个字符
    preview := bodyStr
    if len(bodyStr) > 200 {
        preview = bodyStr[:200] + "..."
    }
    fmt.Printf("\n响应预览:\n%s\n", preview)

    return nil
}

// 示例2: 响应体验证器
type ResponseValidator struct {
    Rules   []ValidationRule
    Results []ValidationResult
    Mutex   sync.RWMutex
}

type ValidationRule struct {
    ID          string
    Name        string
    Type        string // "size", "json", "regex", "contains", "status"
    Condition   interface{}
    Description string
    Severity    string // "info", "warning", "error"
}

type ValidationResult struct {
    RequestID   string
    RuleID      string
    Timestamp   time.Time
    Passed      bool
    Message     string
    Details     interface{}
    Severity    string
}

func NewResponseValidator() *ResponseValidator {
    return &ResponseValidator{
        Rules:   make([]ValidationRule, 0),
        Results: make([]ValidationResult, 0),
    }
}

func (v *ResponseValidator) AddRule(rule ValidationRule) {
    v.Mutex.Lock()
    defer v.Mutex.Unlock()

    v.Rules = append(v.Rules, rule)
    log.Printf("添加验证规则: %s", rule.Name)
}

func (v *ResponseValidator) ValidateResponse(requestID string) error {
    log.Printf("验证响应，请求ID: %s", requestID)

    // 获取响应体
    info, err := CDPFetchGetResponseBodyParsed(requestID)
    if err != nil {
        return fmt.Errorf("获取响应体失败: %w", err)
    }

    // 解码响应体
    var bodyStr string
    if info.Base64Encoded {
        decoded, err := base64.StdEncoding.DecodeString(info.Body)
        if err != nil {
            return fmt.Errorf("Base64解码失败: %w", err)
        }
        bodyStr = string(decoded)
    } else {
        bodyStr = info.Body
    }

    v.Mutex.Lock()
    defer v.Mutex.Unlock()

    // 应用所有规则
    for _, rule := range v.Rules {
        result := ValidationResult{
            RequestID: requestID,
            RuleID:    rule.ID,
            Timestamp: time.Now(),
            Severity:  rule.Severity,
        }

        passed, message, details := v.applyRule(rule, bodyStr, info.Size)
        result.Passed = passed
        result.Message = message
        result.Details = details

        v.Results = append(v.Results, result)

        status := "✓ 通过"
        if !passed {
            status = "✗ 失败"
        }

        log.Printf("%s [%s] %s: %s", status, rule.Severity, rule.Name, message)
    }

    return nil
}

func (v *ResponseValidator) applyRule(rule ValidationRule, body string, size int) (bool, string, interface{}) {
    switch rule.Type {
    case "size":
        maxSize, ok := rule.Condition.(float64)
        if !ok {
            return false, "无效的大小条件", nil
        }

        if size > int(maxSize) {
            return false,
                fmt.Sprintf("响应体大小 %d 超过限制 %d", size, int(maxSize)),
                map[string]interface{}{
                    "actual":   size,
                    "expected": int(maxSize),
                }
        }
        return true,
            fmt.Sprintf("响应体大小 %d 符合要求", size),
            map[string]interface{}{"size": size}

    case "json":
        var data interface{}
        if err := json.Unmarshal([]byte(body), &data); err != nil {
            return false,
                "响应不是有效的JSON格式",
                map[string]interface{}{"error": err.Error()}
        }
        return true, "响应是有效的JSON格式", nil

    case "regex":
        pattern, ok := rule.Condition.(string)
        if !ok {
            return false, "无效的正则表达式", nil
        }

        matched, err := regexp.MatchString(pattern, body)
        if err != nil {
            return false,
                fmt.Sprintf("正则表达式错误: %v", err),
                nil
        }

        if !matched {
            return false,
                "响应不匹配正则表达式",
                map[string]interface{}{"pattern": pattern}
        }
        return true, "响应匹配正则表达式", nil

    case "contains":
        text, ok := rule.Condition.(string)
        if !ok {
            return false, "无效的包含条件", nil
        }

        if !strings.Contains(body, text) {
            return false,
                fmt.Sprintf("响应不包含文本: %s", text),
                nil
        }
        return true,
            fmt.Sprintf("响应包含文本: %s", text),
            nil

    default:
        return false, "未知的规则类型", nil
    }
}

func (v *ResponseValidator) GetStats() map[string]interface{} {
    v.Mutex.RLock()
    defer v.Mutex.RUnlock()

    stats := make(map[string]interface{})
    stats["rules_count"] = len(v.Rules)
    stats["results_count"] = len(v.Results)

    if len(v.Results) > 0 {
        passed := 0
        for _, result := range v.Results {
            if result.Passed {
                passed++
            }
        }
        stats["passed_count"] = passed
        stats["failed_count"] = len(v.Results) - passed
        stats["pass_rate"] = float64(passed) / float64(len(v.Results)) * 100
    }

    return stats
}

func (v *ResponseValidator) GenerateReport() {
    stats := v.GetStats()

    fmt.Println("=== 响应验证报告 ===")
    fmt.Printf("验证规则数量: %d\n", stats["rules_count"])
    fmt.Printf("验证结果数量: %d\n", stats["results_count"])

    if passRate, ok := stats["pass_rate"].(float64); ok {
        fmt.Printf("通过率: %.1f%%\n", passRate)
    }

    v.Mutex.RLock()
    defer v.Mutex.RUnlock()

    if len(v.Results) > 0 {
        fmt.Println("\n最近验证结果:")
        startIdx := len(v.Results) - 5
        if startIdx < 0 {
            startIdx = 0
        }

        for i := startIdx; i < len(v.Results); i++ {
            result := v.Results[i]
            status := "✓ 通过"
            if !result.Passed {
                status = "✗ 失败"
            }

            fmt.Printf("[%s] %s %s: %s\n",
                result.Timestamp.Format("15:04:05"),
                status,
                result.RuleID,
                result.Message)
        }
    }
}

*/

// -----------------------------------------------  Fetch.takeResponseBodyAsStream  -----------------------------------------------
// === 应用场景 ===
// 1. 流式处理: 处理大型响应体，避免内存溢出
// 2. 实时处理: 实时处理流式数据
// 3. 大文件下载: 下载大文件时流式处理
// 4. 视频流处理: 处理视频或音频流
// 5. 数据管道: 构建数据流处理管道
// 6. 性能优化: 优化大响应体的处理性能

// CDPFetchTakeResponseBodyAsStream 以流的方式获取响应体
// requestID: 请求ID
func CDPFetchTakeResponseBodyAsStream(requestID string) (string, error) {
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
        "method": "Fetch.takeResponseBodyAsStream",
        "params": {
            "requestId": "%s"
        }
    }`, reqID, requestID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Fetch.takeResponseBodyAsStream 请求失败: %w", err)
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
			return "", fmt.Errorf("Fetch.takeResponseBodyAsStream 请求超时")
		}
	}
}

// StreamResult 流处理结果
type StreamResult struct {
	StreamID string `json:"stream"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// 流处理器接口
type StreamProcessor interface {
	ProcessChunk(data []byte) error
	OnComplete() error
	OnError(err error)
	GetStats() map[string]interface{}
}

// 基础流处理器
type BaseStreamProcessor struct {
	StreamID       string
	TotalBytes     int64
	ChunkCount     int
	StartTime      time.Time
	EndTime        time.Time
	IsComplete     bool
	HasError       bool
	ErrorMsg       string
	ContentType    string
	LastChunkTime  time.Time
	ProcessingTime time.Duration
}

// 创建流式响应处理器
func NewStreamResponseProcessor(streamID string) *StreamResponseProcessor {
	return &StreamResponseProcessor{
		BaseStreamProcessor: BaseStreamProcessor{
			StreamID:    streamID,
			StartTime:   time.Now(),
			ContentType: "unknown",
		},
		Data:              make([]byte, 0),
		ChunkCallbacks:    make([]func([]byte) error, 0),
		CompleteCallbacks: make([]func() error, 0),
		ErrorCallbacks:    make([]func(error), 0),
	}
}

// 流式响应处理器
type StreamResponseProcessor struct {
	BaseStreamProcessor
	Data              []byte
	ChunkCallbacks    []func([]byte) error
	CompleteCallbacks []func() error
	ErrorCallbacks    []func(error)
	Mutex             sync.RWMutex
}

// 处理数据块
func (p *StreamResponseProcessor) ProcessChunk(data []byte) error {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	start := time.Now()

	p.ChunkCount++
	p.TotalBytes += int64(len(data))
	p.LastChunkTime = time.Now()

	// 保存数据
	p.Data = append(p.Data, data...)

	// 调用块回调
	for _, callback := range p.ChunkCallbacks {
		if err := callback(data); err != nil {
			p.HasError = true
			p.ErrorMsg = err.Error()
			p.ProcessingTime += time.Since(start)
			return err
		}
	}

	p.ProcessingTime += time.Since(start)
	return nil
}

// 完成处理
func (p *StreamResponseProcessor) OnComplete() error {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.IsComplete = true
	p.EndTime = time.Now()

	// 调用完成回调
	for _, callback := range p.CompleteCallbacks {
		if err := callback(); err != nil {
			p.HasError = true
			p.ErrorMsg = err.Error()
			return err
		}
	}

	return nil
}

// 错误处理
func (p *StreamResponseProcessor) OnError(err error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.HasError = true
	p.ErrorMsg = err.Error()
	p.EndTime = time.Now()

	// 调用错误回调
	for _, callback := range p.ErrorCallbacks {
		callback(err)
	}
}

// 注册块处理回调
func (p *StreamResponseProcessor) RegisterChunkCallback(callback func([]byte) error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.ChunkCallbacks = append(p.ChunkCallbacks, callback)
}

// 注册完成回调
func (p *StreamResponseProcessor) RegisterCompleteCallback(callback func() error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.CompleteCallbacks = append(p.CompleteCallbacks, callback)
}

// 注册错误回调
func (p *StreamResponseProcessor) RegisterErrorCallback(callback func(error)) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.ErrorCallbacks = append(p.ErrorCallbacks, callback)
}

// 获取统计数据
func (p *BaseStreamProcessor) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["stream_id"] = p.StreamID
	stats["total_bytes"] = p.TotalBytes
	stats["chunk_count"] = p.ChunkCount
	stats["start_time"] = p.StartTime
	stats["content_type"] = p.ContentType
	stats["is_complete"] = p.IsComplete
	stats["has_error"] = p.HasError
	stats["error"] = p.ErrorMsg

	if p.StartTime.IsZero() {
		stats["duration"] = 0
	} else if p.EndTime.IsZero() {
		stats["duration"] = time.Since(p.StartTime)
	} else {
		stats["duration"] = p.EndTime.Sub(p.StartTime)
	}

	if p.TotalBytes > 0 && !p.StartTime.IsZero() {
		if p.EndTime.IsZero() {
			stats["throughput_bps"] = float64(p.TotalBytes) / time.Since(p.StartTime).Seconds()
		} else {
			stats["throughput_bps"] = float64(p.TotalBytes) / p.EndTime.Sub(p.StartTime).Seconds()
		}
	}

	if p.ChunkCount > 0 {
		stats["avg_chunk_size"] = p.TotalBytes / int64(p.ChunkCount)
	}

	return stats
}

// 辅助函数: 处理流式响应
func ProcessStreamResponse(requestID string, processor StreamProcessor) error {
	log.Printf("开始流式处理响应，请求ID: %s", requestID)

	// 获取流ID
	result, err := CDPFetchTakeResponseBodyAsStream(requestID)
	if err != nil {
		return fmt.Errorf("获取响应流失败: %w", err)
	}

	var streamResp struct {
		Result struct {
			Stream string `json:"stream"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &streamResp); err != nil {
		return fmt.Errorf("解析流响应失败: %w", err)
	}

	streamID := streamResp.Result.Stream
	log.Printf("获取到流ID: %s", streamID)

	// 开始监听流数据
	return listenToStream(streamID, processor)
}

// 监听流数据
func listenToStream(streamID string, processor StreamProcessor) error {
	// 这里需要实现流数据监听逻辑
	// 实际上，CDP会通过IO.read等事件发送流数据

	log.Printf("开始监听流: %s", streamID)
	// 返回一个占位错误，实际实现需要集成CDP事件系统
	return fmt.Errorf("流监听需要集成CDP事件系统")
}

/*

// 示例1: 文件下载处理器
type FileDownloadProcessor struct {
    BaseStreamProcessor
    FilePath      string
    File          *os.File
    FileSize      int64
    Checksum      string
    Hash          string
    WriteMutex    sync.Mutex
}

func NewFileDownloadProcessor(streamID, filePath string) *FileDownloadProcessor {
    return &FileDownloadProcessor{
        BaseStreamProcessor: BaseStreamProcessor{
            StreamID:  streamID,
            StartTime: time.Now(),
        },
        FilePath: filePath,
    }
}

func (p *FileDownloadProcessor) ProcessChunk(data []byte) error {
    p.WriteMutex.Lock()
    defer p.WriteMutex.Unlock()

    if p.File == nil {
        // 创建文件
        file, err := os.Create(p.FilePath)
        if err != nil {
            return fmt.Errorf("创建文件失败: %w", err)
        }
        p.File = file
    }

    // 写入文件
    n, err := p.File.Write(data)
    if err != nil {
        return fmt.Errorf("写入文件失败: %w", err)
    }

    p.ChunkCount++
    p.TotalBytes += int64(n)
    p.FileSize += int64(n)
    p.LastChunkTime = time.Now()

    // 更新校验和
    p.updateChecksum(data)

    return nil
}

func (p *FileDownloadProcessor) updateChecksum(data []byte) {
    // 简单的CRC32校验
    if p.Checksum == "" {
        p.Checksum = fmt.Sprintf("%08x", crc32.ChecksumIEEE(data))
    } else {
        // 合并校验和
        current, _ := strconv.ParseUint(p.Checksum, 16, 32)
        newChecksum := crc32.ChecksumIEEE(data)
        combined := uint32(current) ^ newChecksum
        p.Checksum = fmt.Sprintf("%08x", combined)
    }

    // 更新MD5哈希
    hash := md5.New()
    if p.Hash != "" {
        prevHash, _ := hex.DecodeString(p.Hash)
        hash.Write(prevHash)
    }
    hash.Write(data)
    p.Hash = hex.EncodeToString(hash.Sum(nil))
}

func (p *FileDownloadProcessor) OnComplete() error {
    p.WriteMutex.Lock()
    defer p.WriteMutex.Unlock()

    p.IsComplete = true
    p.EndTime = time.Now()

    if p.File != nil {
        if err := p.File.Close(); err != nil {
            return fmt.Errorf("关闭文件失败: %w", err)
        }
    }

    log.Printf("文件下载完成: %s, 大小: %d 字节", p.FilePath, p.FileSize)
    return nil
}

func (p *FileDownloadProcessor) OnError(err error) {
    p.WriteMutex.Lock()
    defer p.WriteMutex.Unlock()

    p.HasError = true
    p.ErrorMsg = err.Error()
    p.EndTime = time.Now()

    if p.File != nil {
        p.File.Close()
        // 删除不完整的文件
        os.Remove(p.FilePath)
    }

    log.Printf("文件下载失败: %v", err)
}

func (p *FileDownloadProcessor) GetStats() map[string]interface{} {
    stats := p.BaseStreamProcessor.GetStats()
    stats["file_path"] = p.FilePath
    stats["file_size"] = p.FileSize
    stats["checksum"] = p.Checksum
    stats["hash"] = p.Hash
    stats["download_speed"] = 0.0

    if p.TotalBytes > 0 && !p.StartTime.IsZero() {
        var duration time.Duration
        if p.EndTime.IsZero() {
            duration = time.Since(p.StartTime)
        } else {
            duration = p.EndTime.Sub(p.StartTime)
        }

        if duration.Seconds() > 0 {
            stats["download_speed"] = float64(p.TotalBytes) / duration.Seconds() / 1024 / 1024 // MB/s
        }
    }

    return stats
}

// 示例2: JSON流式解析器
type JSONStreamParser struct {
    BaseStreamProcessor
    Buffer        []byte
    Objects       []map[string]interface{}
    ObjectCount   int
    ParseErrors   int
    BufferMutex   sync.RWMutex
    ParserConfig  ParserConfig
}

type ParserConfig struct {
    BufferSize    int
    MaxObjects    int
    ValidateJSON  bool
    ExtractFields []string
    OnObject      func(map[string]interface{}) error
}

func NewJSONStreamParser(streamID string, config ParserConfig) *JSONStreamParser {
    if config.BufferSize == 0 {
        config.BufferSize = 1024 * 1024 // 1MB
    }

    return &JSONStreamParser{
        BaseStreamProcessor: BaseStreamProcessor{
            StreamID:  streamID,
            StartTime: time.Now(),
            ContentType: "application/json",
        },
        Buffer:      make([]byte, 0, config.BufferSize),
        Objects:     make([]map[string]interface{}, 0),
        ParserConfig: config,
    }
}

func (p *JSONStreamParser) ProcessChunk(data []byte) error {
    p.BufferMutex.Lock()
    defer p.BufferMutex.Unlock()

    p.ChunkCount++
    p.TotalBytes += int64(len(data))
    p.LastChunkTime = time.Now()

    // 添加到缓冲区
    p.Buffer = append(p.Buffer, data...)

    // 尝试解析完整的JSON对象
    if err := p.tryParseObjects(); err != nil {
        p.ParseErrors++
        return fmt.Errorf("解析JSON失败: %w", err)
    }

    return nil
}

func (p *JSONStreamParser) tryParseObjects() error {
    // 寻找JSON对象的开始和结束
    start := 0
    depth := 0
    inString := false
    escapeNext := false

    for i, b := range p.Buffer {
        if escapeNext {
            escapeNext = false
            continue
        }

        switch b {
        case '\\':
            escapeNext = true
        case '"':
            if !inString {
                inString = true
            } else {
                inString = false
            }
        case '{':
            if !inString && depth == 0 {
                start = i
            }
            if !inString {
                depth++
            }
        case '}':
            if !inString {
                depth--
                if depth == 0 {
                    // 找到完整的JSON对象
                    jsonStr := p.Buffer[start : i+1]
                    if err := p.parseSingleObject(jsonStr); err != nil {
                        return err
                    }

                    // 清理已处理的数据
                    p.Buffer = p.Buffer[i+1:]
                    return p.tryParseObjects() // 递归处理剩余数据
                }
            }
        }
    }

    return nil
}

func (p *JSONStreamParser) parseSingleObject(data []byte) error {
    if p.ParserConfig.ValidateJSON {
        if !json.Valid(data) {
            return fmt.Errorf("无效的JSON: %s", string(data))
        }
    }

    var obj map[string]interface{}
    if err := json.Unmarshal(data, &obj); err != nil {
        return fmt.Errorf("JSON解析失败: %w", err)
    }

    p.ObjectCount++

    // 提取指定字段
    if len(p.ParserConfig.ExtractFields) > 0 {
        extracted := make(map[string]interface{})
        for _, field := range p.ParserConfig.ExtractFields {
            if val, exists := obj[field]; exists {
                extracted[field] = val
            }
        }
        p.Objects = append(p.Objects, extracted)
    } else {
        p.Objects = append(p.Objects, obj)
    }

    // 调用回调
    if p.ParserConfig.OnObject != nil {
        if err := p.ParserConfig.OnObject(obj); err != nil {
            return fmt.Errorf("对象回调失败: %w", err)
        }
    }

    // 检查对象数量限制
    if p.ParserConfig.MaxObjects > 0 && len(p.Objects) >= p.ParserConfig.MaxObjects {
        return fmt.Errorf("达到最大对象数量限制: %d", p.ParserConfig.MaxObjects)
    }

    return nil
}

func (p *JSONStreamParser) OnComplete() error {
    p.BufferMutex.Lock()
    defer p.BufferMutex.Unlock()

    p.IsComplete = true
    p.EndTime = time.Now()

    // 尝试解析缓冲区中剩余的数据
    if len(p.Buffer) > 0 {
        if err := p.tryParseObjects(); err != nil {
            p.ParseErrors++
            log.Printf("解析剩余数据失败: %v", err)
        }
    }

    log.Printf("JSON流解析完成: 解析了 %d 个对象", p.ObjectCount)
    return nil
}

func (p *JSONStreamParser) OnError(err error) {
    p.BufferMutex.Lock()
    defer p.BufferMutex.Unlock()

    p.HasError = true
    p.ErrorMsg = err.Error()
    p.EndTime = time.Now()

    log.Printf("JSON流解析错误: %v", err)
}

func (p *JSONStreamParser) GetStats() map[string]interface{} {
    stats := p.BaseStreamProcessor.GetStats()
    stats["object_count"] = p.ObjectCount
    stats["parse_errors"] = p.ParseErrors
    stats["buffer_size"] = len(p.Buffer)
    stats["objects_per_second"] = 0.0

    if p.ObjectCount > 0 && !p.StartTime.IsZero() {
        var duration time.Duration
        if p.EndTime.IsZero() {
            duration = time.Since(p.StartTime)
        } else {
            duration = p.EndTime.Sub(p.StartTime)
        }

        if duration.Seconds() > 0 {
            stats["objects_per_second"] = float64(p.ObjectCount) / duration.Seconds()
        }
    }

    return stats
}

// 示例3: 实时数据分析器
type RealtimeDataAnalyzer struct {
    BaseStreamProcessor
    Metrics       map[string]float64
    Counters      map[string]int64
    Samples       []DataSample
    WindowSize    int
    AnalysisFunc  func([]byte) (map[string]float64, error)
    MetricsMutex  sync.RWMutex
}

type DataSample struct {
    Timestamp time.Time
    Metrics   map[string]float64
    RawSize   int
}

func NewRealtimeDataAnalyzer(streamID string, windowSize int, analysisFunc func([]byte) (map[string]float64, error)) *RealtimeDataAnalyzer {
    return &RealtimeDataAnalyzer{
        BaseStreamProcessor: BaseStreamProcessor{
            StreamID:  streamID,
            StartTime: time.Now(),
        },
        Metrics:      make(map[string]float64),
        Counters:     make(map[string]int64),
        Samples:      make([]DataSample, 0, windowSize),
        WindowSize:   windowSize,
        AnalysisFunc: analysisFunc,
    }
}

func (a *RealtimeDataAnalyzer) ProcessChunk(data []byte) error {
    a.MetricsMutex.Lock()
    defer a.MetricsMutex.Unlock()

    a.ChunkCount++
    a.TotalBytes += int64(len(data))
    a.LastChunkTime = time.Now()

    // 分析数据
    metrics, err := a.AnalysisFunc(data)
    if err != nil {
        return fmt.Errorf("数据分析失败: %w", err)
    }

    // 记录样本
    sample := DataSample{
        Timestamp: time.Now(),
        Metrics:   metrics,
        RawSize:   len(data),
    }

    a.Samples = append(a.Samples, sample)

    // 维护窗口大小
    if len(a.Samples) > a.WindowSize {
        a.Samples = a.Samples[1:]
    }

    // 更新累计指标
    for key, value := range metrics {
        a.Metrics[key] = value
        a.Counters[key]++
    }

    return nil
}

func (a *RealtimeDataAnalyzer) OnComplete() error {
    a.MetricsMutex.Lock()
    defer a.MetricsMutex.Unlock()

    a.IsComplete = true
    a.EndTime = time.Now()

    log.Printf("实时数据分析完成: 处理了 %d 个样本", len(a.Samples))
    return nil
}

func (a *RealtimeDataAnalyzer) OnError(err error) {
    a.MetricsMutex.Lock()
    defer a.MetricsMutex.Unlock()

    a.HasError = true
    a.ErrorMsg = err.Error()
    a.EndTime = time.Now()

    log.Printf("实时数据分析错误: %v", err)
}

func (a *RealtimeDataAnalyzer) GetCurrentMetrics() map[string]interface{} {
    a.MetricsMutex.RLock()
    defer a.MetricsMutex.RUnlock()

    result := make(map[string]interface{})

    // 计算窗口统计
    if len(a.Samples) > 0 {
        // 计算平均值
        sums := make(map[string]float64)
        counts := make(map[string]int)

        for _, sample := range a.Samples {
            for key, value := range sample.Metrics {
                sums[key] += value
                counts[key]++
            }
        }

        for key, sum := range sums {
            if count, ok := counts[key]; ok && count > 0 {
                result[fmt.Sprintf("%s_avg", key)] = sum / float64(count)
            }
        }

        // 最新值
        latest := a.Samples[len(a.Samples)-1]
        for key, value := range latest.Metrics {
            result[fmt.Sprintf("%s_latest", key)] = value
        }

        // 变化率
        if len(a.Samples) >= 2 {
            oldest := a.Samples[0]
            for key, latestValue := range latest.Metrics {
                if oldValue, ok := oldest.Metrics[key]; ok {
                    result[fmt.Sprintf("%s_change", key)] = latestValue - oldValue
                }
            }
        }
    }

    // 累计统计
    result["total_samples"] = len(a.Samples)
    result["total_bytes"] = a.TotalBytes

    return result
}

func (a *RealtimeDataAnalyzer) GetStats() map[string]interface{} {
    stats := a.BaseStreamProcessor.GetStats()

    // 添加分析统计
    currentMetrics := a.GetCurrentMetrics()
    for key, value := range currentMetrics {
        stats[key] = value
    }

    stats["sample_count"] = len(a.Samples)
    stats["metric_count"] = len(a.Metrics)
    stats["window_size"] = a.WindowSize

    return stats
}

// 示例4: 多路流处理器
type MultiplexStreamProcessor struct {
    BaseStreamProcessor
    Processors    []StreamProcessor
    Router        func([]byte) (int, error)
    ChunkRouter   func([]byte) []int
    Stats         []map[string]interface{}
    Mutex         sync.RWMutex
}

func NewMultiplexStreamProcessor(streamID string) *MultiplexStreamProcessor {
    return &MultiplexStreamProcessor{
        BaseStreamProcessor: BaseStreamProcessor{
            StreamID:  streamID,
            StartTime: time.Now(),
        },
        Processors: make([]StreamProcessor, 0),
        Stats:      make([]map[string]interface{}, 0),
    }
}

func (m *MultiplexStreamProcessor) AddProcessor(processor StreamProcessor) int {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    m.Processors = append(m.Processors, processor)
    m.Stats = append(m.Stats, make(map[string]interface{}))

    return len(m.Processors) - 1
}

func (m *MultiplexStreamProcessor) ProcessChunk(data []byte) error {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    m.ChunkCount++
    m.TotalBytes += int64(len(data))
    m.LastChunkTime = time.Now()

    var errors []error

    if m.Router != nil {
        // 使用路由器选择处理器
        index, err := m.Router(data)
        if err != nil {
            return fmt.Errorf("路由失败: %w", err)
        }

        if index >= 0 && index < len(m.Processors) {
            if err := m.Processors[index].ProcessChunk(data); err != nil {
                errors = append(errors, fmt.Errorf("处理器 %d 失败: %w", index, err))
            }
        }
    } else if m.ChunkRouter != nil {
        // 使用块路由器选择多个处理器
        indices := m.ChunkRouter(data)
        for _, index := range indices {
            if index >= 0 && index < len(m.Processors) {
                if err := m.Processors[index].ProcessChunk(data); err != nil {
                    errors = append(errors, fmt.Errorf("处理器 %d 失败: %w", index, err))
                }
            }
        }
    } else {
        // 发送给所有处理器
        for i, processor := range m.Processors {
            if err := processor.ProcessChunk(data); err != nil {
                errors = append(errors, fmt.Errorf("处理器 %d 失败: %w", i, err))
            }
        }
    }

    // 更新统计
    for i, processor := range m.Processors {
        m.Stats[i] = processor.GetStats()
    }

    if len(errors) > 0 {
        return fmt.Errorf("处理器错误: %v", errors)
    }

    return nil
}

func (m *MultiplexStreamProcessor) OnComplete() error {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    m.IsComplete = true
    m.EndTime = time.Now()

    var errors []error

    for i, processor := range m.Processors {
        if err := processor.OnComplete(); err != nil {
            errors = append(errors, fmt.Errorf("处理器 %d 完成失败: %w", i, err))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("完成错误: %v", errors)
    }

    log.Printf("多路流处理器完成: 共 %d 个处理器", len(m.Processors))
    return nil
}

func (m *MultiplexStreamProcessor) OnError(err error) {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    m.HasError = true
    m.ErrorMsg = err.Error()
    m.EndTime = time.Now()

    for i, processor := range m.Processors {
        processor.OnError(fmt.Errorf("多路处理器错误: %w", err))
        m.Stats[i] = processor.GetStats()
    }

    log.Printf("多路流处理器错误: %v", err)
}

func (m *MultiplexStreamProcessor) GetStats() map[string]interface{} {
    stats := m.BaseStreamProcessor.GetStats()

    m.Mutex.RLock()
    defer m.Mutex.RUnlock()

    stats["processor_count"] = len(m.Processors)
    stats["processor_stats"] = m.Stats

    // 计算总体统计
    var totalBytes int64
    var totalChunks int

    for _, procStats := range m.Stats {
        if bytes, ok := procStats["total_bytes"].(int64); ok {
            totalBytes += bytes
        }
        if chunks, ok := procStats["chunk_count"].(int); ok {
            totalChunks += chunks
        }
    }

    stats["total_processed_bytes"] = totalBytes
    stats["total_processed_chunks"] = totalChunks

    return stats
}

// 示例5: 流监控器
type StreamMonitor struct {
    Streams       map[string]*StreamInfo
    Alerts        []StreamAlert
    EventChannel  chan StreamEvent
    IsMonitoring  bool
    Mutex         sync.RWMutex
}

type StreamInfo struct {
    StreamID      string
    Processor     StreamProcessor
    StartTime     time.Time
    LastActivity  time.Time
    Stats         map[string]interface{}
    Metadata      map[string]interface{}
    IsActive      bool
    ErrorCount    int
}

type StreamAlert struct {
    ID          string
    Condition   func(*StreamInfo) bool
    Action      func(*StreamInfo, string)
    Description string
    Severity    string
    Cooldown    time.Duration
    LastAlert   time.Time
}

type StreamEvent struct {
    Type      string
    StreamID  string
    Timestamp time.Time
    Data      interface{}
    Message   string
}

func NewStreamMonitor() *StreamMonitor {
    return &StreamMonitor{
        Streams:      make(map[string]*StreamInfo),
        Alerts:       make([]StreamAlert, 0),
        EventChannel: make(chan StreamEvent, 100),
    }
}

func (m *StreamMonitor) RegisterStream(streamID string, processor StreamProcessor, metadata map[string]interface{}) error {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    if _, exists := m.Streams[streamID]; exists {
        return fmt.Errorf("流已存在: %s", streamID)
    }

    info := &StreamInfo{
        StreamID:     streamID,
        Processor:    processor,
        StartTime:    time.Now(),
        LastActivity: time.Now(),
        Stats:        make(map[string]interface{}),
        Metadata:     metadata,
        IsActive:     true,
    }

    m.Streams[streamID] = info

    m.EventChannel <- StreamEvent{
        Type:      "stream_registered",
        StreamID:  streamID,
        Timestamp: time.Now(),
        Data:      metadata,
        Message:   "流已注册",
    }

    log.Printf("流已注册: %s", streamID)
    return nil
}

func (m *StreamMonitor) UnregisterStream(streamID string) error {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    info, exists := m.Streams[streamID]
    if !exists {
        return fmt.Errorf("流不存在: %s", streamID)
    }

    info.IsActive = false
    delete(m.Streams, streamID)

    m.EventChannel <- StreamEvent{
        Type:      "stream_unregistered",
        StreamID:  streamID,
        Timestamp: time.Now(),
        Data:      info.Stats,
        Message:   "流已注销",
    }

    log.Printf("流已注销: %s", streamID)
    return nil
}

func (m *StreamMonitor) UpdateStreamActivity(streamID string) {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    if info, exists := m.Streams[streamID]; exists {
        info.LastActivity = time.Now()
        info.Stats = info.Processor.GetStats()

        // 检查告警
        m.checkAlerts(info)
    }
}

func (m *StreamMonitor) AddAlert(alert StreamAlert) {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    m.Alerts = append(m.Alerts, alert)
    log.Printf("添加流告警: %s", alert.Description)
}

func (m *StreamMonitor) checkAlerts(info *StreamInfo) {
    for i := range m.Alerts {
        alert := &m.Alerts[i]

        if alert.Condition(info) {
            // 检查冷却时间
            if time.Since(alert.LastAlert) > alert.Cooldown {
                alert.LastAlert = time.Now()

                // 触发告警
                if alert.Action != nil {
                    go alert.Action(info, alert.Description)
                }

                m.EventChannel <- StreamEvent{
                    Type:      "alert_triggered",
                    StreamID:  info.StreamID,
                    Timestamp: time.Now(),
                    Data:      info.Stats,
                    Message:   alert.Description,
                }
            }
        }
    }
}

func (m *StreamMonitor) StartMonitoring() {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    if m.IsMonitoring {
        return
    }

    m.IsMonitoring = true
    go m.processEvents()
    go m.healthCheck()

    log.Println("流监控器已启动")
}

func (m *StreamMonitor) StopMonitoring() {
    m.Mutex.Lock()
    defer m.Mutex.Unlock()

    if !m.IsMonitoring {
        return
    }

    m.IsMonitoring = false
    close(m.EventChannel)

    log.Println("流监控器已停止")
}

func (m *StreamMonitor) processEvents() {
    for event := range m.EventChannel {
        log.Printf("[流监控] %s [%s] %s",
            event.Timestamp.Format("15:04:05"),
            event.Type,
            event.Message)
    }
}

func (m *StreamMonitor) healthCheck() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for m.IsMonitoring {
        select {
        case <-ticker.C:
            m.checkStreamHealth()
        }
    }
}

func (m *StreamMonitor) checkStreamHealth() {
    m.Mutex.RLock()
    defer m.Mutex.RUnlock()

    now := time.Now()

    for streamID, info := range m.Streams {
        if info.IsActive {
            // 检查是否超时
            if now.Sub(info.LastActivity) > 5*time.Minute {
                m.EventChannel <- StreamEvent{
                    Type:      "stream_timeout",
                    StreamID:  streamID,
                    Timestamp: now,
                    Data:      info.Stats,
                    Message:   "流已超时",
                }

                info.IsActive = false
            }
        }
    }
}

func (m *StreamMonitor) GetStats() map[string]interface{} {
    m.Mutex.RLock()
    defer m.Mutex.RUnlock()

    stats := make(map[string]interface{})
    stats["is_monitoring"] = m.IsMonitoring
    stats["active_streams"] = 0
    stats["total_streams"] = len(m.Streams)
    stats["alert_count"] = len(m.Alerts)

    var totalBytes int64
    var totalChunks int

    for _, info := range m.Streams {
        if info.IsActive {
            stats["active_streams"] = stats["active_streams"].(int) + 1
        }

        if procStats, ok := info.Stats["total_bytes"].(int64); ok {
            totalBytes += procStats
        }
        if procStats, ok := info.Stats["chunk_count"].(int); ok {
            totalChunks += procStats
        }
    }

    stats["total_processed_bytes"] = totalBytes
    stats["total_processed_chunks"] = totalChunks

    return stats
}

func (m *StreamMonitor) GenerateReport() {
    stats := m.GetStats()

    fmt.Println("=== 流监控报告 ===")
    fmt.Printf("监控状态: %v\n", stats["is_monitoring"])
    fmt.Printf("总流数量: %d\n", stats["total_streams"])
    fmt.Printf("活动流数量: %d\n", stats["active_streams"])
    fmt.Printf("告警数量: %d\n", stats["alert_count"])
    fmt.Printf("处理总字节数: %d\n", stats["total_processed_bytes"])
    fmt.Printf("处理总块数: %d\n", stats["total_processed_chunks"])

    m.Mutex.RLock()
    defer m.Mutex.RUnlock()

    if len(m.Streams) > 0 {
        fmt.Println("\n流状态:")
        for streamID, info := range m.Streams {
            status := "活跃"
            if !info.IsActive {
                status = "非活跃"
            }

            fmt.Printf("  %-20s [%s] 最后活动: %v\n",
                streamID,
                status,
                time.Since(info.LastActivity))
        }
    }
}


*/

/*


综合性的测试

// 示例2: 模拟API响应
func SimulateAPIResponse() {
    log.Println("=== 模拟API响应 ===")

    // 启用Fetch拦截
    patterns := []FetchPattern{
        {
            URLPattern:   "*\/api/*", // 拦截API请求
            ResourceType: "XHR",
            RequestStage: "Request",
        },
    }

    if _, err := CDPFetchEnable(patterns, false); err != nil {
        log.Printf("启用Fetch拦截失败: %v", err)
        return
    }
    defer CDPFetchDisable()

    // 模拟请求被拦截
    requestID := "api_request_001"

    // 创建模拟响应
    mockResponse := &FulfillRequestResponse{
        ResponseCode: 200,
        ResponseHeaders: []Header{
            {Name: "Content-Type", Value: "application/json"},
            {Name: "Cache-Control", Value: "no-cache"},
        },
        Body: `{"status": "success", "data": {"id": 1, "name": "Test User"}}`,
        ResponsePhrase: "OK",
    }

    // 完成请求
    result, err := CDPFetchFulfillRequest(requestID, mockResponse)
    if err != nil {
        log.Printf("模拟响应失败: %v", err)
    } else {
        log.Printf("模拟响应成功: %s", result)
    }
}

// 示例3: 模拟请求失败
func SimulateRequestFailure() {
    log.Println("=== 模拟请求失败 ===")

    if _, err := CDPFetchEnable([]FetchPattern{{URLPattern: "*"}}, false); err != nil {
        log.Printf("启用Fetch拦截失败: %v", err)
        return
    }
    defer CDPFetchDisable()

    requestID := "failing_request_001"

    // 模拟网络错误
    result, err := CDPFetchFailRequest(requestID, "ConnectionFailed")
    if err != nil {
        log.Printf("模拟失败失败: %v", err)
    } else {
        log.Printf("模拟失败成功: %s", result)
    }

    // 模拟超时
    requestID = "timeout_request_001"
    result, err = CDPFetchFailRequest(requestID, "Timeout")
    if err != nil {
        log.Printf("模拟超时失败: %v", err)
    } else {
        log.Printf("模拟超时成功: %s", result)
    }
}


// 示例5: 批量请求处理
func BatchRequestProcessing() {
    log.Println("=== 批量请求处理 ===")

    if _, err := CDPFetchEnable([]FetchPattern{{URLPattern: "*"}}, false); err != nil {
        log.Printf("启用Fetch拦截失败: %v", err)
        return
    }
    defer CDPFetchDisable()

    // 模拟多个请求
    requests := []struct {
        ID       string
        URL      string
        Method   string
    }{
        {"req_001", "https://api.example.com/users", "GET"},
        {"req_002", "https://api.example.com/products", "GET"},
        {"req_003", "https://api.example.com/orders", "POST"},
    }

    var wg sync.WaitGroup
    results := make(chan string, len(requests))

    for _, req := range requests {
        wg.Add(1)
        go func(requestID, url, method string) {
            defer wg.Done()

            // 修改请求
            modifications := &ContinueRequestModifications{
                Method: method,
                Headers: []Header{
                    {Name: "X-Request-ID", Value: requestID},
                    {Name: "Origin", Value: "https://example.com"},
                },
            }

            result, err := CDPFetchContinueRequest(requestID, modifications)
            if err != nil {
                results <- fmt.Sprintf("请求 %s 失败: %v", requestID, err)
            } else {
                results <- fmt.Sprintf("请求 %s 成功: %s", requestID, result)
            }
        }(req.ID, req.URL, req.Method)
    }

    wg.Wait()
    close(results)

    for result := range results {
        log.Println(result)
    }
}

// 示例6: 网络请求模拟器
type NetworkRequestSimulator struct {
    BaseURL    string
    Mocks      map[string]*MockResponse
    Patterns   []FetchPattern
    IsEnabled  bool
    RequestLog []RequestLogEntry
    Mutex      sync.RWMutex
}

type MockResponse struct {
    StatusCode  int
    Headers     []Header
    Body        interface{}
    Delay       time.Duration
    FailReason  string
    ContentType string
}

type RequestLogEntry struct {
    ID        string
    URL       string
    Method    string
    Timestamp time.Time
    Status    string
    Duration  time.Duration
    Error     string
}

func NewNetworkRequestSimulator(baseURL string) *NetworkRequestSimulator {
    return &NetworkRequestSimulator{
        BaseURL:    baseURL,
        Mocks:      make(map[string]*MockResponse),
        RequestLog: make([]RequestLogEntry, 0),
    }
}

func (sim *NetworkRequestSimulator) AddMock(endpoint string, mock *MockResponse) {
    sim.Mutex.Lock()
    defer sim.Mutex.Unlock()

    pattern := FetchPattern{
        URLPattern:   sim.BaseURL + endpoint,
        RequestStage: "Request",
    }
    sim.Patterns = append(sim.Patterns, pattern)
    sim.Mocks[endpoint] = mock

    log.Printf("添加Mock: %s%s -> %d", sim.BaseURL, endpoint, mock.StatusCode)
}

func (sim *NetworkRequestSimulator) Start() error {
    sim.Mutex.Lock()
    defer sim.Mutex.Unlock()

    if sim.IsEnabled {
        return fmt.Errorf("模拟器已启用")
    }

    if _, err := CDPFetchEnable(sim.Patterns, false); err != nil {
        return fmt.Errorf("启用Fetch拦截失败: %w", err)
    }

    sim.IsEnabled = true
    log.Printf("网络请求模拟器已启动，Mock数量: %d", len(sim.Mocks))
    return nil
}

func (sim *NetworkRequestSimulator) Stop() {
    sim.Mutex.Lock()
    defer sim.Mutex.Unlock()

    if !sim.IsEnabled {
        return
    }

    CDPFetchDisable()
    sim.IsEnabled = false
    log.Printf("网络请求模拟器已停止，请求日志数量: %d", len(sim.RequestLog))
}

func (sim *NetworkRequestSimulator) HandleRequest(requestID, url, method string) {
    sim.Mutex.Lock()
    defer sim.Mutex.Unlock()

    startTime := time.Now()
    logEntry := RequestLogEntry{
        ID:        requestID,
        URL:       url,
        Method:    method,
        Timestamp: startTime,
    }

    // 查找匹配的Mock
    var mock *MockResponse
    var endpoint string

    for ep, m := range sim.Mocks {
        if strings.Contains(url, ep) {
            mock = m
            endpoint = ep
            break
        }
    }

    if mock == nil {
        // 如果没有匹配的Mock，继续原始请求
        go func() {
            if _, err := CDPFetchContinueRequest(requestID, nil); err != nil {
                logEntry.Error = err.Error()
                logEntry.Status = "error"
            } else {
                logEntry.Status = "passed-through"
            }
            logEntry.Duration = time.Since(startTime)
            sim.RequestLog = append(sim.RequestLog, logEntry)
        }()
        return
    }

    // 应用延迟
    if mock.Delay > 0 {
        time.Sleep(mock.Delay)
    }

    // 处理失败情况
    if mock.FailReason != "" {
        go func() {
            if _, err := CDPFetchFailRequest(requestID, mock.FailReason); err != nil {
                logEntry.Error = err.Error()
                logEntry.Status = "mock-fail-error"
            } else {
                logEntry.Status = "mocked-failure"
            }
            logEntry.Duration = time.Since(startTime)
            sim.RequestLog = append(sim.RequestLog, logEntry)
        }()
        return
    }

    // 构建响应
    body, err := json.Marshal(mock.Body)
    if err != nil {
        logEntry.Error = fmt.Sprintf("序列化响应体失败: %v", err)
        logEntry.Status = "error"
        logEntry.Duration = time.Since(startTime)
        sim.RequestLog = append(sim.RequestLog, logEntry)
        return
    }

    contentType := mock.ContentType
    if contentType == "" {
        contentType = "application/json"
    }

    response := &FulfillRequestResponse{
        ResponseCode: mock.StatusCode,
        ResponseHeaders: []Header{
            {Name: "Content-Type", Value: contentType},
            {Name: "X-Mocked-By", Value: "NetworkRequestSimulator"},
            {Name: "X-Mock-Endpoint", Value: endpoint},
        },
        Body: string(body),
    }

    go func() {
        if _, err := CDPFetchFulfillRequest(requestID, response); err != nil {
            logEntry.Error = err.Error()
            logEntry.Status = "mock-error"
        } else {
            logEntry.Status = "mocked"
        }
        logEntry.Duration = time.Since(startTime)
        sim.RequestLog = append(sim.RequestLog, logEntry)
    }()
}

func (sim *NetworkRequestSimulator) GetStats() map[string]interface{} {
    sim.Mutex.RLock()
    defer sim.Mutex.RUnlock()

    stats := make(map[string]interface{})
    stats["is_enabled"] = sim.IsEnabled
    stats["mock_count"] = len(sim.Mocks)
    stats["request_log_count"] = len(sim.RequestLog)
    stats["patterns_count"] = len(sim.Patterns)

    if len(sim.RequestLog) > 0 {
        stats["first_request"] = sim.RequestLog[0].Timestamp
        stats["last_request"] = sim.RequestLog[len(sim.RequestLog)-1].Timestamp
        stats["success_count"] = countByStatus(sim.RequestLog, "mocked")
        stats["failure_count"] = countByStatus(sim.RequestLog, "mock-fail-error")
        stats["passthrough_count"] = countByStatus(sim.RequestLog, "passed-through")
    }

    return stats
}

func (sim *NetworkRequestSimulator) ClearLogs() {
    sim.Mutex.Lock()
    defer sim.Mutex.Unlock()

    sim.RequestLog = make([]RequestLogEntry, 0)
    log.Println("请求日志已清理")
}

func (sim *NetworkRequestSimulator) GenerateReport() {
    sim.Mutex.RLock()
    defer sim.Mutex.RUnlock()

    fmt.Println("=== 网络请求模拟器报告 ===")
    fmt.Printf("模拟器状态: %v\n", sim.IsEnabled)
    fmt.Printf("Mock数量: %d\n", len(sim.Mocks))
    fmt.Printf("请求日志数量: %d\n\n", len(sim.RequestLog))

    if len(sim.RequestLog) == 0 {
        fmt.Println("无请求记录")
        return
    }

    // Mock使用统计
    mockUsage := make(map[string]int)
    for _, entry := range sim.RequestLog {
        for endpoint := range sim.Mocks {
            if strings.Contains(entry.URL, endpoint) {
                mockUsage[endpoint]++
            }
        }
    }

    fmt.Println("Mock使用统计:")
    for endpoint, count := range mockUsage {
        fmt.Printf("  %-30s: %d 次\n", endpoint, count)
    }

    // 状态统计
    statusCounts := make(map[string]int)
    for _, entry := range sim.RequestLog {
        statusCounts[entry.Status]++
    }

    fmt.Println("\n请求状态统计:")
    for status, count := range statusCounts {
        fmt.Printf("  %-20s: %d 次\n", status, count)
    }

    // 最近请求
    fmt.Println("\n最近请求 (最多5个):")
    startIdx := len(sim.RequestLog) - 5
    if startIdx < 0 {
        startIdx = 0
    }

    for i := startIdx; i < len(sim.RequestLog); i++ {
        entry := sim.RequestLog[i]
        fmt.Printf("%s [%s] %s %-8s 耗时: %v\n",
            entry.Timestamp.Format("15:04:05.000"),
            entry.ID,
            entry.Method,
            entry.Status,
            entry.Duration)

        if entry.Error != "" {
            fmt.Printf("      错误: %s\n", entry.Error)
        }
    }
}

func countByStatus(logs []RequestLogEntry, status string) int {
    count := 0
    for _, entry := range logs {
        if entry.Status == status {
            count++
        }
    }
    return count
}

// 示例7: API响应时间测试
func TestAPIResponseTimes() {
    log.Println("=== API响应时间测试 ===")

    if _, err := CDPFetchEnable([]FetchPattern{{URLPattern: "*\/api/*"}}, false); err != nil {
        log.Printf("启用Fetch拦截失败: %v", err)
        return
    }
    defer CDPFetchDisable()

    // 测试不同的延迟
    delays := []time.Duration{
        100 * time.Millisecond,
        500 * time.Millisecond,
        1 * time.Second,
        2 * time.Second,
        5 * time.Second,
    }

    for i, delay := range delays {
        requestID := fmt.Sprintf("delay_test_%d", i+1)

        // 创建有延迟的响应
        response := &FulfillRequestResponse{
            ResponseCode: 200,
            ResponseHeaders: []Header{
                {Name: "Content-Type", Value: "application/json"},
                {Name: "X-Test-Delay", Value: delay.String()},
            },
            Body: fmt.Sprintf(`{"test": "delay_%d", "delay_ms": %d}`, i+1, delay.Milliseconds()),
        }

        startTime := time.Now()
        if _, err := CDPFetchFulfillRequest(requestID, response); err != nil {
            log.Printf("延迟测试 %v 失败: %v", delay, err)
        } else {
            actualDelay := time.Since(startTime)
            log.Printf("延迟测试 %v 完成，实际耗时: %v", delay, actualDelay)
        }

        time.Sleep(500 * time.Millisecond)
    }
}


*/
