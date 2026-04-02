package browser

import (
	"ChromeBot/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// CDPDOMStorageClear 清除指定存储区域的所有数据
// 参数说明:
//   - storageId: 存储ID，包含安全源和存储类型信息
//
// DOMStorage.clear用于清除指定存储区域的所有数据：
// 清除localStorage或sessionStorage中的所有键值对
// 针对特定的存储ID（包括安全源、是否是localStorage）
// 不会清除其他源的存储数据
// 不会清除cookie或其他存储类型
func CDPDOMStorageClear(storageId DOMStorageId) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMStorage.clear")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMStorage.clear",
		"params": {
			"storageId": {
				"securityOrigin": "%s",
				"isLocalStorage": %t
			}
		}
	}`, reqID, storageId.SecurityOrigin, storageId.IsLocalStorage)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMStorage.clear 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
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
				fmt.Println("[CDP DOMStorage.clear] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMStorage.clear 请求超时")
		}
	}
}

// DOMStorageId DOM存储ID结构
type DOMStorageId struct {
	SecurityOrigin string `json:"securityOrigin"`
	IsLocalStorage bool   `json:"isLocalStorage"`
}

// 创建存储ID的辅助函数
func NewDOMStorageId(origin string, isLocalStorage bool) DOMStorageId {
	return DOMStorageId{
		SecurityOrigin: origin,
		IsLocalStorage: isLocalStorage,
	}
}

// LocalStorageId 创建localStorage ID
func LocalStorageId(origin string) DOMStorageId {
	return DOMStorageId{
		SecurityOrigin: origin,
		IsLocalStorage: true,
	}
}

// SessionStorageId 创建sessionStorage ID
func SessionStorageId(origin string) DOMStorageId {
	return DOMStorageId{
		SecurityOrigin: origin,
		IsLocalStorage: false,
	}
}

// ClearLocalStorage 清除指定源的localStorage
func ClearLocalStorage(origin string) {
	storageId := LocalStorageId(origin)
	fmt.Printf("正在清除 %s 的localStorage...", origin)
	response, err := CDPDOMStorageClear(storageId)
	if err != nil {
		log.Printf("清除localStorage失败: %v", err)
		return
	}
	fmt.Printf("localStorage清除成功: %s", response)
}

// ClearSessionStorage 清除指定源的sessionStorage
func ClearSessionStorage(origin string) {
	// 创建sessionStorage ID
	storageId := SessionStorageId(origin)
	fmt.Printf("正在清除 %s 的sessionStorage...", origin)
	response, err := CDPDOMStorageClear(storageId)
	if err != nil {
		log.Printf("清除sessionStorage失败: %v", err)
		return
	}
	fmt.Printf("sessionStorage清除成功: %s", response)
}

// CDPDOMStorageDisable 禁用DOMStorage域
func CDPDOMStorageDisable() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMStorage.disable")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMStorage.disable"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMStorage.disable 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
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
				fmt.Println("[CDP DOMStorage.disable] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMStorage.disable 请求超时")
		}
	}
}

// CDPDOMStorageEnable 启用DOMStorage域
func CDPDOMStorageEnable() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMStorage.enable")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMStorage.enable"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMStorage.enable 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
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
				fmt.Println("[CDP DOMStorage.enable] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMStorage.enable 请求超时")
		}
	}
}

// CDPDOMStorageGetDOMStorageItems 获取指定存储区域的所有项目
// 参数说明:
//   - storageId: 存储ID，包含安全源和存储类型信息
//
// DOMStorage.getDOMStorageItems用于获取指定存储区域的所有项目：
// 获取localStorage或sessionStorage中的所有键值对
// 针对特定的存储ID（包括安全源、是否是localStorage）
func CDPDOMStorageGetDOMStorageItems(storageId DOMStorageId) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMStorage.getDOMStorageItems")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMStorage.getDOMStorageItems",
		"params": {
			"storageId": {
				"securityOrigin": "%s",
				"isLocalStorage": %t
			}
		}
	}`, reqID, storageId.SecurityOrigin, storageId.IsLocalStorage)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMStorage.getDOMStorageItems 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
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
				fmt.Println("[CDP DOMStorage.getDOMStorageItems] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMStorage.getDOMStorageItems 请求超时")
		}
	}
}

// CDPDOMStorageRemoveDOMStorageItem 删除指定存储区域的特定项目
// 参数说明:
//   - storageId: 存储ID，包含安全源和存储类型信息
//   - key: 要删除的键
//
// DOMStorage.removeDOMStorageItem用于删除指定存储区域中的特定项目：
// 从localStorage或sessionStorage中删除指定的键值对
// 针对特定的存储ID（包括安全源、是否是localStorage）
// 只删除指定的键，不影响其他数据
// 如果键不存在，操作也会成功
// 会触发DOMStorage存储变化事件
func CDPDOMStorageRemoveDOMStorageItem(storageId DOMStorageId, key string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMStorage.removeDOMStorageItem")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMStorage.removeDOMStorageItem",
		"params": {
			"storageId": {
				"securityOrigin": "%s",
				"isLocalStorage": %t
			},
			"key": "%s"
		}
	}`, reqID, storageId.SecurityOrigin, storageId.IsLocalStorage, key)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMStorage.removeDOMStorageItem 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
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
				fmt.Println("[CDP DOMStorage.removeDOMStorageItem] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMStorage.removeDOMStorageItem 请求超时")
		}
	}
}

// CDPDOMStorageSetDOMStorageItem 在指定存储区域中设置项目
// 参数说明:
//   - storageId: 存储ID，包含安全源和存储类型信息
//   - key: 要设置的键
//   - value: 要设置的值
//
// DOMStorage.setDOMStorageItem用于在指定存储区域中设置项目：
// 在localStorage或sessionStorage中设置键值对
// 如果键已存在，则更新其值
// 针对特定的存储ID（包括安全源、是否是localStorage）
// 会触发DOMStorage存储变化事件
// 遵循同源策略的安全限制
func CDPDOMStorageSetDOMStorageItem(storageId DOMStorageId, key, value string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMStorage.setDOMStorageItem")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 转义特殊字符
	escapedKey := strings.ReplaceAll(key, `"`, `\"`)
	escapedValue := strings.ReplaceAll(value, `"`, `\"`)

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMStorage.setDOMStorageItem",
		"params": {
			"storageId": {
				"securityOrigin": "%s",
				"isLocalStorage": %t
			},
			"key": "%s",
			"value": "%s"
		}
	}`, reqID, storageId.SecurityOrigin, storageId.IsLocalStorage, escapedKey, escapedValue)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMStorage.setDOMStorageItem 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
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
				fmt.Println("[CDP DOMStorage.setDOMStorageItem] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMStorage.setDOMStorageItem 请求超时")
		}
	}
}

/*

将下面示例写成cbs脚本的示例

// 示例2: 设置测试数据
func exampleSetTestData() {
	origin := "https://test.example.com"

	log.Printf("为 %s 设置测试数据...", origin)
	storageId := LocalStorageId(origin)

	// 测试数据
	testData := map[string]string{
		"test_user_id":      "user_12345",
		"test_session_id":   "session_abc123def456",
		"test_timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
		"test_config":       `{"env":"testing","debug":true,"logLevel":"debug"}`,
		"test_items":        `["item1","item2","item3"]`,
		"test_counter":      "0",
		"test_flag_enabled": "true",
		"test_api_endpoint": "https://api.test.example.com",
		"test_version":      "1.0.0",
		"test_mode":         "automation",
	}

	successCount := 0
	for key, value := range testData {
		log.Printf("设置: %s = %s", key, value)

		_, err := CDPDOMStorageSetDOMStorageItem(storageId, key, value)
		if err != nil {
			log.Printf("  设置失败: %v", err)
		} else {
			successCount++
		}

		time.Sleep(50 * time.Millisecond) // 避免过快
	}

	log.Printf("测试数据设置完成: %d/%d 成功", successCount, len(testData))
}


// 示例4: 模拟用户会话
func exampleSimulateUserSession() {
	origin := "https://app.example.com"

	log.Printf("为 %s 模拟用户会话...", origin)
	storageId := LocalStorageId(origin)

	// 模拟用户登录
	log.Println("模拟用户登录...")
	sessionData := map[string]string{
		"user_id":       "user_789012",
		"username":      "john_doe",
		"email":         "john@example.com",
		"auth_token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"session_id":    "session_xyz789",
		"login_time":    fmt.Sprintf("%d", time.Now().Unix()),
		"permissions":   `["read","write","delete"]`,
		"preferences":   `{"theme":"light","notifications":true,"timezone":"UTC"}`,
		"last_active":   fmt.Sprintf("%d", time.Now().Unix()),
		"user_role":     "admin",
	}

	// 设置会话数据
	for key, value := range sessionData {
		_, err := CDPDOMStorageSetDOMStorageItem(storageId, key, value)
		if err != nil {
			log.Printf("设置 %s 失败: %v", key, err)
		}
		time.Sleep(20 * time.Millisecond)
	}

	log.Println("用户会话已设置")

	// 模拟用户操作
	log.Println("模拟用户操作...")
	time.Sleep(2 * time.Second)

	// 更新最后活跃时间
	newLastActive := fmt.Sprintf("%d", time.Now().Unix())
	_, err := CDPDOMStorageSetDOMStorageItem(storageId, "last_active", newLastActive)
	if err != nil {
		log.Printf("更新最后活跃时间失败: %v", err)
	} else {
		log.Printf("最后活跃时间已更新: %s", newLastActive)
	}

	// 模拟添加用户偏好
	newPreferences := `{"theme":"dark","notifications":true,"timezone":"UTC","fontSize":"medium"}`
	_, err = CDPDOMStorageSetDOMStorageItem(storageId, "preferences", newPreferences)
	if err != nil {
		log.Printf("更新偏好失败: %v", err)
	} else {
		log.Println("用户偏好已更新")
	}

	log.Println("用户会话模拟完成")
}


*/
