package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// CDPTargetActivateTarget 激活指定的目标
func CDPTargetActivateTarget(targetId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.activateTarget")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.activateTarget",
		"params": {
			"targetId": "%s"
		}
	}`, reqID, targetId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.activateTarget 失败:", err)
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
				fmt.Println("[CDP Target.activateTarget] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("activateTarget 请求超时")
		}
	}
}

// CDPTargetAttachToTarget 附加到指定的目标
func CDPTargetAttachToTarget(targetId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.attachToTarget")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.attachToTarget",
		"params": {
			"targetId": "%s"
		}
	}`, reqID, targetId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.attachToTarget 失败:", err)
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
				fmt.Println("[CDP Target.attachToTarget] 收到回复 -> ", content)

				// 解析响应，提取sessionId
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return "", fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查错误
				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("%v", errorObj)
					return "", fmt.Errorf("CDP错误: %s", errorMsg)
				}

				// 提取sessionId
				if result, ok := response["result"].(map[string]interface{}); ok {
					if sessionId, ok := result["sessionId"].(string); ok {
						return sessionId, nil
					}
				}

				return "", fmt.Errorf("响应中未找到sessionId")
			}
		case <-timer.C:
			return "", fmt.Errorf("attachToTarget 请求超时")
		}
	}
}

// CDPTargetCloseTarget 关闭指定的目标
func CDPTargetCloseTarget(targetId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.closeTarget")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.closeTarget",
		"params": {
			"targetId": "%s"
		}
	}`, reqID, targetId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.closeTarget 失败:", err)
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
				fmt.Println("[CDP Target.closeTarget] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("closeTarget 请求超时")
		}
	}
}

// CDPTargetCreateBrowserContext 创建一个新的浏览器上下文
// 返回值:
//   - browserContextId: 新创建的浏览器上下文ID
//   - error: 创建过程中发生的错误
//
// Target.createBrowserContext用于创建一个新的浏览器上下文，这相当于创建一个独立的浏览器会话环境：
// 创建隔离的浏览器会话, 独立的cookie、localStorage、sessionStorage存储, 独立的缓存和网络数据, 独立的扩展程序环境, 独立的权限设置
func CDPTargetCreateBrowserContext() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.createBrowserContext")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.createBrowserContext"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.createBrowserContext 失败:", err)
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
				fmt.Println("[CDP Target.createBrowserContext] 收到回复 -> ", content)

				// 解析响应，提取browserContextId
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return "", fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查错误
				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("%v", errorObj)
					return "", fmt.Errorf("CDP错误: %s", errorMsg)
				}

				// 提取browserContextId
				if result, ok := response["result"].(map[string]interface{}); ok {
					if browserContextId, ok := result["browserContextId"].(string); ok {
						return browserContextId, nil
					}
				}

				return "", fmt.Errorf("响应中未找到browserContextId")
			}
		case <-timer.C:
			return "", fmt.Errorf("createBrowserContext 请求超时")
		}
	}
}

// CDPTargetCreateTarget 创建一个新的目标（页面/标签页）
// 参数说明:
//   - url: 新目标的初始URL
//
// 可选参数:
//   - width: 页面宽度
//   - height: 页面高度
//   - browserContextId: 浏览器上下文ID
//
// 返回值:
//   - targetId: 新创建的目标ID
//   - error: 创建过程中发生的错误
func CDPTargetCreateTarget(url string, options ...CreateTargetOption) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.createTarget")
	}

	// 默认配置
	config := &CreateTargetConfig{
		URL: url,
		//Width:                   800,
		//Height:                  600,
		EnableBeginFrameControl: false, // 可选，是否启用帧控制
		NewWindow:               false, // 可选，是否在新窗口中打开（默认 true）
		Background:              false, // 可选，是否在后台打开
	}

	// 应用选项
	for _, option := range options {
		option(config)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.createTarget",
		"params": {
			"url": "%s",
			"enableBeginFrameControl": %t,
			"newWindow": %t,
			"background": %t`,
		reqID, config.URL, config.EnableBeginFrameControl, config.NewWindow, config.Background)

	// 可选参数
	if config.BrowserContextId != "" {
		message += fmt.Sprintf(`, "browserContextId": "%s"`, config.BrowserContextId)
	}

	message += `}}`

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.createTarget 失败:", err)
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
				fmt.Println("[CDP Target.createTarget] 收到回复 -> ", content)

				// 解析响应，提取targetId
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return "", fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查错误
				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("%v", errorObj)
					return "", fmt.Errorf("CDP错误: %s", errorMsg)
				}

				// 提取targetId
				if result, ok := response["result"].(map[string]interface{}); ok {
					if targetId, ok := result["targetId"].(string); ok {
						return targetId, nil
					}
				}

				return "", fmt.Errorf("响应中未找到targetId")
			}
		case <-timer.C:
			return "", fmt.Errorf("createTarget 请求超时")
		}
	}
}

// CreateTargetConfig 创建目标的配置
type CreateTargetConfig struct {
	URL                     string
	Width                   int
	Height                  int
	BrowserContextId        string
	EnableBeginFrameControl bool
	NewWindow               bool
	Background              bool
}

// CreateTargetOption 配置选项
type CreateTargetOption func(*CreateTargetConfig)

// WithWidth 设置页面宽度
func WithWidth(width int) CreateTargetOption {
	return func(c *CreateTargetConfig) {
		c.Width = width
	}
}

// WithHeight 设置页面高度
func WithHeight(height int) CreateTargetOption {
	return func(c *CreateTargetConfig) {
		c.Height = height
	}
}

// WithBrowserContext 设置浏览器上下文
func WithBrowserContext(contextId string) CreateTargetOption {
	return func(c *CreateTargetConfig) {
		c.BrowserContextId = contextId
	}
}

// WithFrameControl 启用帧控制
func WithFrameControl(enable bool) CreateTargetOption {
	return func(c *CreateTargetConfig) {
		c.EnableBeginFrameControl = enable
	}
}

// InCurrentWindow 在当前窗口打开
func InCurrentWindow() CreateTargetOption {
	return func(c *CreateTargetConfig) {
		c.NewWindow = false
	}
}

// OpenInBackground 在后台打开
func OpenInBackground() CreateTargetOption {
	return func(c *CreateTargetConfig) {
		c.Background = true
	}
}

// CDPTargetDetachFromTarget 从目标分离会话
func CDPTargetDetachFromTarget(sessionId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.detachFromTarget")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.detachFromTarget",
		"params": {
			"sessionId": "%s"
		}
	}`, reqID, sessionId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.detachFromTarget 失败:", err)
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
				fmt.Println("[CDP Target.detachFromTarget] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("detachFromTarget 请求超时")
		}
	}
}

/*
// 示例2: 安全的临时操作
func exampleSafeOperation() {
	targetId := "ABC123DEF456"

	// 使用defer确保总是分离
	sessionId, err := CDPTargetAttachToTarget(targetId)
	if err != nil {
		log.Printf("附加失败: %v", err)
		return
	}

	// 确保在函数退出时分离
	defer func() {
		if _, err := CDPTargetDetachFromTarget(sessionId); err != nil {
			log.Printf("分离失败: %v", err)
		} else {
			log.Printf("会话 %s 已安全分离", sessionId)
		}
	}()

	log.Printf("已附加，sessionId: %s", sessionId)

	// 执行一些操作
	// 如果这里发生panic，defer仍然会执行分离

	// 操作1
	// result1, err := doOperation1(sessionId)

	// 操作2
	// result2, err := doOperation2(sessionId)

	// ... 更多操作
}
*/

// CDPTargetDisposeBrowserContext 销毁指定的浏览器上下文
// 参数说明:
//   - browserContextId: 要销毁的浏览器上下文ID
//
// Target.disposeBrowserContext用于销毁指定的浏览器上下文，这会清理该上下文中的所有资源：
// 关闭该上下文中的所有目标（标签页、窗口）, 清理上下文的cookies、localStorage、sessionStorage, 释放上下文占用的所有内存和资源
// 清理缓存、网络数据等, 销毁扩展程序环境, 不可逆操作，上下文被永久删除
func CDPTargetDisposeBrowserContext(browserContextId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.disposeBrowserContext")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.disposeBrowserContext",
		"params": {
			"browserContextId": "%s"
		}
	}`, reqID, browserContextId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.disposeBrowserContext 失败:", err)
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
				fmt.Println("[CDP Target.disposeBrowserContext] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("disposeBrowserContext 请求超时")
		}
	}
}

// CDPTargetGetBrowserContexts 获取所有浏览器上下文
func CDPTargetGetBrowserContexts() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.getBrowserContexts")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.getBrowserContexts"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.getBrowserContexts 失败:", err)
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
				fmt.Println("[CDP Target.getBrowserContexts] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("getBrowserContexts 请求超时")
		}
	}
}

// CDPTargetGetTargets 获取所有可用的目标
func CDPTargetGetTargets() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.getTargets")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.getTargets"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.getTargets 失败:", err)
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
				fmt.Println("[CDP Target.getTargets] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("getTargets 请求超时")
		}
	}
}

// CDPTargetGetTargetInfo 获取指定目标的详细信息
func CDPTargetGetTargetInfo(targetId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Target.getTargetInfo")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.getTargetInfo",
		"params": {
			"targetId": "%s"
		}
	}`, reqID, targetId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Target.getTargetInfo 失败:", err)
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
				fmt.Println("[CDP Target.getTargetInfo] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("getTargetInfo 请求超时")
		}
	}
}
