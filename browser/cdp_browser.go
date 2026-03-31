package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// CDPBrowserClose 关闭浏览器 Browser.close
func CDPBrowserClose() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP 未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 Browser.close")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// Browser.close 方法不需要参数
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Browser.close"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 Browser.close 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)

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
				utils.Debugf("[CDP Browser.close] 收到回复 -> %s", content)
				if chromeInstance.BrowserWSConn != nil {
					chromeInstance.BrowserWSConn.Close()
					chromeInstance.BrowserWSConn = nil
				}
				return content, nil
			}
		case <-timer.C:
			if chromeInstance.BrowserWSConn != nil {
				chromeInstance.BrowserWSConn.Close()
				chromeInstance.BrowserWSConn = nil
			}
			return "", fmt.Errorf("Browser.close 请求超时")
		}
	}
}

// CDPBrowserResetPermissions 重置指定来源的浏览器权限
// 参数说明:
//   - origin: 要重置权限的来源，格式为完整的origin（如 "https://example.com:8080"）
//
// 适用场景:
// 1. 自动化测试中清理测试环境
// 2. 用户隐私保护，退出时清除权限
// 3. 开发者调试权限相关功能
// 4. 合规性要求的权限定期清理
func CDPBrowserResetPermissions(origin string) (string, error) {
	if origin == "" {
		return "", fmt.Errorf("origin 参数不能为空")
	}

	if !utils.IsValidOrigin(origin) {
		return "", fmt.Errorf("无效的 origin 格式: %s。应为完整origin，如 https://example.com", origin)
	}

	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP 功能未启用，无法调用 Browser.resetPermissions")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Browser.resetPermissions",
		"params": {
			"origin": "%s"
		}
	}`, reqID, origin)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("[ERROR] 发送 Browser.resetPermissions 请求失败 (origin: %s): %v", origin, err)
		return "", fmt.Errorf("发送重置权限请求失败: %w", err)
	}

	utils.Debugf("[CDP] 发送重置权限请求: %s", message)

	timeout := 3 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("CDP消息队列异常关闭")
			}
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				utils.Debugf("[CDP] Browser.resetPermissions 响应: %s", content)
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err == nil {
					if errorObj, exists := response["error"]; exists {
						errorMsg := fmt.Sprintf("%v", errorObj)
						log.Printf("[WARN] Browser.resetPermissions 返回错误 (origin: %s): %v", origin, errorObj)
						// 处理特定错误类型
						if strings.Contains(errorMsg, "invalid origin") {
							return "", fmt.Errorf("无效的origin: %s", origin)
						} else if strings.Contains(errorMsg, "not found") {
							return "", fmt.Errorf("指定的origin不存在或没有权限设置: %s", origin)
						}
						return "", fmt.Errorf("CDP错误: %s", errorMsg)
					}
				}
				log.Printf("[INFO] 已重置 origin '%s' 的所有权限", origin)
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("重置权限请求超时（%v），origin: %s", timeout, origin)
		}
	}
}

// CDPBrowserGetWindowForTarget 通过targetId获取对应的窗口ID
// targetId: 可以通过 Target.getTargets 获取
func CDPBrowserGetWindowForTarget(targetId string) (int, error) {
	if !DefaultBrowserWS() {
		return 0, fmt.Errorf("CDP 功能未启用")
	}

	if chromeInstance.BrowserWSConn == nil {
		return 0, fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Browser.getWindowForTarget",
		"params": {
			"targetId": "%s"
		}
	}`, reqID, targetId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("[ERROR] 发送 Browser.getWindowForTarget 失败: %v", err)
		return 0, err
	}

	timeout := 3 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return 0, fmt.Errorf("消息队列已关闭")
			}
			if reqID == respMsg.ID {
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return 0, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return 0, fmt.Errorf("CDP错误: %v", errorObj)
				}

				if resultObj, ok := response["result"].(map[string]interface{}); ok {
					if windowId, ok := resultObj["windowId"].(float64); ok {
						return int(windowId), nil
					}
				}
				return 0, fmt.Errorf("未找到 windowId")
			}
		case <-timer.C:
			return 0, fmt.Errorf("获取窗口ID超时")
		}
	}
}

// CDPBrowserGetWindowBounds 获取指定浏览器窗口的边界信息
// 参数说明:
//   - windowId: 目标窗口的唯一标识符 CDPBrowserGetWindowForTarget(targetId)来获取
//
// 返回值:
//   - bounds: 窗口边界信息，包括位置、尺寸和状态
//
// 适用场景:
// 1. 自动化测试中验证窗口位置和大小
// 2. 多窗口应用管理窗口布局
// 3. 远程协助获取用户窗口位置
// 4. 用户行为分析记录窗口使用习惯
func CDPBrowserGetWindowBounds(windowId int) (map[string]interface{}, error) {
	if windowId <= 0 {
		return nil, fmt.Errorf("无效的窗口ID: %d。windowId 必须为正整数", windowId)
	}

	if !DefaultBrowserWS() {
		return nil, fmt.Errorf("CDP 功能未启用，无法调用 Browser.getWindowBounds")
	}

	if chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Browser.getWindowBounds",
		"params": {
			"windowId": %d
		}
	}`, reqID, windowId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("[ERROR] 发送 Browser.getWindowBounds 请求失败 (windowId: %d): %v", windowId, err)
		return nil, fmt.Errorf("发送获取窗口边界请求失败: %w", err)
	}

	utils.Debugf("[CDP] 发送获取窗口边界请求: %s", message)

	timeout := 3 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("CDP消息队列异常关闭")
			}
			if reqID == respMsg.ID {
				utils.Debugf("[CDP] Browser.getWindowBounds 原始响应: %s", respMsg.Content)
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					log.Printf("[ERROR] 解析响应JSON失败: %v", err)
					return nil, fmt.Errorf("解析响应数据失败: %w", err)
				}
				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("%v", errorObj)
					log.Printf("[ERROR] Browser.getWindowBounds 返回错误 (windowId: %d): %v", windowId, errorObj)
					// 处理特定错误类型
					if strings.Contains(errorMsg, "Invalid window id") {
						return nil, fmt.Errorf("无效的窗口ID: %d", windowId)
					} else if strings.Contains(errorMsg, "cannot find window") {
						return nil, fmt.Errorf("找不到窗口ID: %d", windowId)
					} else if strings.Contains(errorMsg, "permission") {
						return nil, fmt.Errorf("权限不足，无法获取窗口边界信息")
					}
					return nil, fmt.Errorf("CDP错误: %s", errorMsg)
				}
				resultObj, hasResult := response["result"]
				if !hasResult {
					return nil, fmt.Errorf("响应中缺少 result 字段")
				}
				resultMap, ok := resultObj.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("result 字段格式不正确")
				}
				boundsObj, hasBounds := resultMap["bounds"]
				if !hasBounds {
					return nil, fmt.Errorf("响应中缺少 bounds 字段")
				}
				bounds, ok := boundsObj.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("bounds 字段格式不正确")
				}
				// 验证 bounds 字段完整性
				requiredFields := []string{"left", "top", "width", "height", "windowState"}
				for _, field := range requiredFields {
					if _, exists := bounds[field]; !exists {
						log.Printf("[WARN] bounds 字段缺少 %s", field)
					}
				}
				// 添加 windowId 到返回数据中便于识别
				bounds["windowId"] = windowId
				utils.Debugf("[INFO] 获取到窗口 %d 的边界信息: %+v", windowId, bounds)
				return bounds, nil
			}
		case <-timer.C:
			return nil, fmt.Errorf("获取窗口边界请求超时（%v），windowId: %d", timeout, windowId)
		}
	}
}

// GetMainWindowID 获取主窗口ID（通常为1，但更可靠的方法）
func GetMainWindowID() (int, error) {
	// 方法1: 尝试默认值1（适用于单窗口情况）
	// 先尝试获取窗口边界，如果成功则说明windowId=1存在
	_, err := CDPBrowserGetWindowBounds(1)
	if err == nil {
		log.Println("[INFO] 使用默认窗口ID: 1")
		return 1, nil
	}
	log.Printf("[DEBUG] 窗口ID=1 获取失败: %v", err)

	// 方法2: 获取所有targets并查找主窗口
	targets, err := getAllTargets()
	if err != nil {
		return 0, fmt.Errorf("获取targets失败: %w", err)
	}

	if len(targets) == 0 {
		return 0, fmt.Errorf("未找到任何目标")
	}

	log.Printf("[DEBUG] 共找到 %d 个目标", len(targets))

	// 2.1 查找浏览器主窗口
	for _, target := range targets {
		log.Printf("[DEBUG] 检查目标: ID=%s, Type=%s, Title=%s", target.ID, target.Type, target.Title)

		// 检查是否为主窗口的常见特征
		isMainWindow := false

		// 特征1: 类型为"browser"
		if target.Type == "browser" {
			isMainWindow = true
			log.Printf("[DEBUG] 目标 %s 类型为 browser，可能是主窗口", target.ID)
		}

		// 特征2: 标题包含浏览器特定内容
		browserTitles := []string{"New Tab", "新标签页", "about:blank", "chrome://"}
		for _, title := range browserTitles {
			if strings.Contains(target.Title, title) {
				isMainWindow = true
				log.Printf("[DEBUG] 目标 %s 标题包含 '%s'，可能是主窗口", target.ID, title)
				break
			}
		}

		// 特征3: URL为空或特定模式
		if target.URL == "" || strings.HasPrefix(target.URL, "chrome://") ||
			strings.HasPrefix(target.URL, "about:") {
			isMainWindow = true
			log.Printf("[DEBUG] 目标 %s URL为 '%s'，可能是主窗口", target.ID, target.URL)
		}

		// 特征4: 没有opener（顶级窗口）
		if target.OpenerID == "" {
			isMainWindow = true
			log.Printf("[DEBUG] 目标 %s 没有opener，可能是主窗口", target.ID)
		}

		if isMainWindow {
			windowId, err := CDPBrowserGetWindowForTarget(target.ID)
			if err == nil {
				log.Printf("[INFO] 找到主窗口: ID=%d (来自target: %s)", windowId, target.ID)
				return windowId, nil
			}
			log.Printf("[DEBUG] 目标 %s 无法获取窗口ID: %v", target.ID, err)
		}
	}

	// 2.2 如果没有明确的主窗口特征，尝试按类型筛选
	priorityOrder := []string{"page", "browser", "background_page", "service_worker", "other"}

	for _, targetType := range priorityOrder {
		for _, target := range targets {
			if target.Type == targetType {
				windowId, err := CDPBrowserGetWindowForTarget(target.ID)
				if err == nil {
					log.Printf("[INFO] 按类型顺序找到窗口: ID=%d (类型: %s)", windowId, targetType)
					return windowId, nil
				}
			}
		}
	}

	// 方法3: 尝试所有targets
	log.Println("[DEBUG] 尝试所有targets获取窗口ID...")
	for _, target := range targets {
		windowId, err := CDPBrowserGetWindowForTarget(target.ID)
		if err == nil {
			log.Printf("[INFO] 找到第一个可用窗口: ID=%d (target: %s)", windowId, target.ID)
			return windowId, nil
		}
		log.Printf("[DEBUG] 目标 %s 无法获取窗口ID: %v", target.ID, err)
	}

	return 0, fmt.Errorf("无法找到任何有效窗口，已尝试 %d 个目标", len(targets))
}

// TargetInfo 目标信息结构
type TargetInfo struct {
	ID        string `json:"targetId"`
	Type      string `json:"type"` // 可能的值: "page", "background_page", "service_worker", "shared_worker", "browser", "other"
	Title     string `json:"title"`
	URL       string `json:"url"`
	Attached  bool   `json:"attached"`
	OpenerID  string `json:"openerId,omitempty"`
	BrowserID string `json:"browserContextId,omitempty"`
	Subtype   string `json:"subtype,omitempty"` // 如: "iframe", "prerender" 等
	CanAccess bool   `json:"canAccessOpener"`
}

// getAllTargets 获取所有目标（targets）
func getAllTargets() ([]TargetInfo, error) {
	if !DefaultBrowserWS() {
		return nil, fmt.Errorf("CDP 功能未启用")
	}

	if chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 Target.getTargets 请求
	// 该方法可以获取所有可用的目标列表
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Target.getTargets"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("[ERROR] 发送 Target.getTargets 失败: %v", err)
		return nil, fmt.Errorf("发送获取目标列表请求失败: %w", err)
	}

	utils.Debugf("[CDP] 发送获取目标列表请求: %s", message)

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				return parseTargetsResponse(respMsg.Content)
			}

		case <-timer.C:
			return nil, fmt.Errorf("获取目标列表超时")
		}
	}
}

// parseTargetsResponse 解析 Target.getTargets 的响应
func parseTargetsResponse(responseContent string) ([]TargetInfo, error) {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseContent), &response); err != nil {
		log.Printf("[ERROR] 解析响应JSON失败: %v", err)
		return nil, fmt.Errorf("解析目标列表响应失败: %w", err)
	}

	// 检查错误
	if errorObj, exists := response["error"]; exists {
		errorMsg := fmt.Sprintf("%v", errorObj)
		log.Printf("[ERROR] Target.getTargets 返回错误: %v", errorObj)

		if strings.Contains(errorMsg, "permission") {
			return nil, fmt.Errorf("权限不足，无法获取目标列表")
		}
		return nil, fmt.Errorf("CDP错误: %s", errorMsg)
	}

	// 提取结果
	resultObj, hasResult := response["result"]
	if !hasResult {
		return nil, fmt.Errorf("响应中缺少 result 字段")
	}

	resultMap, ok := resultObj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("result 字段格式不正确")
	}

	// 提取 targetInfos
	targetInfosObj, hasTargetInfos := resultMap["targetInfos"]
	if !hasTargetInfos {
		return nil, fmt.Errorf("响应中缺少 targetInfos 字段")
	}

	targetInfosArray, ok := targetInfosObj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("targetInfos 字段格式不正确")
	}

	// 转换为 TargetInfo 结构
	var targets []TargetInfo
	for _, targetInfoObj := range targetInfosArray {
		targetInfoMap, ok := targetInfoObj.(map[string]interface{})
		if !ok {
			continue
		}

		target := TargetInfo{
			ID:       getString(targetInfoMap, "targetId"),
			Type:     getString(targetInfoMap, "type"),
			Title:    getString(targetInfoMap, "title"),
			URL:      getString(targetInfoMap, "url"),
			Attached: getBool(targetInfoMap, "attached"),
		}

		// 可选字段
		if openerID, ok := targetInfoMap["openerId"].(string); ok && openerID != "" {
			target.OpenerID = openerID
		}

		if browserID, ok := targetInfoMap["browserContextId"].(string); ok && browserID != "" {
			target.BrowserID = browserID
		}

		if subtype, ok := targetInfoMap["subtype"].(string); ok && subtype != "" {
			target.Subtype = subtype
		}

		if canAccess, ok := targetInfoMap["canAccessOpener"].(bool); ok {
			target.CanAccess = canAccess
		}

		targets = append(targets, target)
	}

	utils.Debugf("[INFO] 获取到 %d 个目标", len(targets))
	return targets, nil
}

// getString 从map中安全获取字符串
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getBool 从map中安全获取布尔值
func getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// GetCurrentWindowInfo 获取当前活动窗口的信息
// 这是最完整的实现，集成了所有必要的功能
func GetCurrentWindowInfo() (*WindowInfo, error) {
	// 1. 获取所有targets
	targets, err := getAllTargets()
	if err != nil {
		return nil, fmt.Errorf("获取targets失败: %w", err)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("没有找到任何target")
	}

	log.Printf("[INFO] 开始查找活动窗口，共 %d 个target", len(targets))

	// 2. 查找活动target
	activeTarget, err := findActiveTarget(targets)
	if err != nil {
		return nil, fmt.Errorf("查找活动target失败: %w", err)
	}

	log.Printf("[INFO] 找到活动target: ID=%s, Title=%s", activeTarget.ID, activeTarget.Title)

	// 3. 获取窗口ID
	windowId, err := CDPBrowserGetWindowForTarget(activeTarget.ID)
	if err != nil {
		return nil, fmt.Errorf("获取窗口ID失败: %w", err)
	}

	// 4. 获取窗口边界信息
	var bounds map[string]interface{}
	bounds, err = CDPBrowserGetWindowBounds(windowId)
	if err != nil {
		log.Printf("[WARN] 获取窗口边界失败: %v，使用默认边界", err)
		// 使用默认边界
		bounds = map[string]interface{}{
			"left":        0,
			"top":         0,
			"width":       1024,
			"height":      768,
			"windowState": "normal",
		}
	}

	// 5. 获取页面详细信息
	pageDetails, err := getPageDetails(activeTarget.ID)
	if err != nil {
		log.Printf("[WARN] 获取页面详情失败: %v", err)
		pageDetails = map[string]interface{}{
			"hasFocus":   false,
			"isVisible":  false,
			"readyState": "unknown",
		}
	}

	// 6. 构建完整的窗口信息
	windowInfo := &WindowInfo{
		WindowID:    windowId,
		TargetID:    activeTarget.ID,
		Type:        activeTarget.Type,
		Title:       activeTarget.Title,
		URL:         activeTarget.URL,
		Bounds:      bounds,
		PageDetails: pageDetails,
		Attached:    activeTarget.Attached,
		OpenerID:    activeTarget.OpenerID,
		BrowserID:   activeTarget.BrowserID,
		Subtype:     activeTarget.Subtype,
		CanAccess:   activeTarget.CanAccess,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	log.Printf("[INFO] 活动窗口信息获取完成: WindowID=%d, Title=%s", windowId, activeTarget.Title)
	return windowInfo, nil
}

// WindowInfo 扩展的窗口信息结构
type WindowInfo struct {
	WindowID    int                    `json:"windowId"`
	TargetID    string                 `json:"targetId"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	URL         string                 `json:"url"`
	Bounds      map[string]interface{} `json:"bounds"`
	PageDetails map[string]interface{} `json:"pageDetails"`
	Attached    bool                   `json:"attached"`
	OpenerID    string                 `json:"openerId,omitempty"`
	BrowserID   string                 `json:"browserContextId,omitempty"`
	Subtype     string                 `json:"subtype,omitempty"`
	CanAccess   bool                   `json:"canAccess"`
	IsActive    bool                   `json:"isActive"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// findActiveTarget 查找活动target的完整实现
func findActiveTarget(targets []TargetInfo) (*TargetInfo, error) {
	// 优先级1: 通过JavaScript精确查找可见页面
	if target := findVisiblePageByJS(targets); target != nil {
		return target, nil
	}

	// 优先级2: 查找有实际内容的页面
	if target := findPageWithContent(targets); target != nil {
		return target, nil
	}

	// 优先级3: 查找已附加的页面
	if target := findAttachedPage(targets); target != nil {
		return target, nil
	}

	// 优先级4: 查找浏览器窗口
	if target := findBrowserWindow(targets); target != nil {
		return target, nil
	}

	// 优先级5: 返回第一个页面
	if target := findFirstPage(targets); target != nil {
		return target, nil
	}

	// 优先级6: 返回任意target
	if len(targets) > 0 {
		return &targets[0], nil
	}

	return nil, fmt.Errorf("无法找到活动target")
}

// findVisiblePageByJS 通过JavaScript查找可见页面
func findVisiblePageByJS(targets []TargetInfo) *TargetInfo {
	for _, target := range targets {
		if target.Type != "page" {
			continue
		}

		visible, err := checkVisibility(target.ID)
		if err != nil {
			log.Printf("[DEBUG] 检查页面 %s 可见性失败: %v", target.ID, err)
			continue
		}

		if visible {
			log.Printf("[DEBUG] 找到可见页面: %s", target.ID)
			return &target
		}
	}
	return nil
}

// checkVisibility 检查页面是否可见
func checkVisibility(targetID string) (bool, error) {
	// 尝试附加到target
	session, err := attachToTarget()
	if err != nil {
		return false, err
	}
	defer detachFromTarget(session)

	// 执行检查可见性的JavaScript
	result, err := evaluateJavaScriptSimple(session, `
		(function() {
			try {
				// 检查页面是否可见
				return document.visibilityState === 'visible' && 
					   !document.hidden && 
					   document.hasFocus();
			} catch(e) {
				return false;
			}
		})()
	`)

	if err != nil {
		return false, err
	}

	if visible, ok := result.(bool); ok {
		return visible, nil
	}

	return false, nil
}

// evaluateJavaScriptSimple 简化的JavaScript执行
func evaluateJavaScriptSimple(sessionId, script string) (interface{}, error) {
	if !DefaultBrowserWS() || chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("CDP连接不可用")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	escapedScript := strings.ReplaceAll(script, `"`, `\"`)
	escapedScript = strings.ReplaceAll(escapedScript, "\n", "\\n")

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Runtime.evaluate",
		"params": {
			"expression": "%s",
			"returnByValue": true,
			"awaitPromise": true
		},
		"sessionId": "%s"
	}`, reqID, escapedScript, sessionId)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return nil, fmt.Errorf("发送JS执行请求失败: %w", err)
	}

	timeout := 3 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列关闭")
			}

			if reqID == respMsg.ID {
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return nil, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return nil, fmt.Errorf("JS执行错误: %v", errorObj)
				}

				if resultObj, ok := response["result"].(map[string]interface{}); ok {
					if resultValue, ok := resultObj["result"]; ok {
						if valueMap, ok := resultValue.(map[string]interface{}); ok {
							return valueMap["value"], nil
						}
					}
				}
				return nil, fmt.Errorf("响应格式错误")
			}
		case <-timer.C:
			return nil, fmt.Errorf("执行超时")
		}
	}
}

// findPageWithContent 查找有实际内容的页面
func findPageWithContent(targets []TargetInfo) *TargetInfo {
	for _, target := range targets {
		if target.Type != "page" {
			continue
		}

		// 检查是否有实际内容
		if target.URL != "" &&
			!strings.HasPrefix(target.URL, "chrome://") &&
			!strings.HasPrefix(target.URL, "about:") &&
			target.URL != "about:blank" &&
			target.Title != "" &&
			target.Title != "New Tab" &&
			target.Title != "新标签页" {

			log.Printf("[DEBUG] 找到有内容页面: %s (URL: %s)", target.ID, target.URL)
			return &target
		}
	}
	return nil
}

// findAttachedPage 查找已附加的页面
func findAttachedPage(targets []TargetInfo) *TargetInfo {
	for _, target := range targets {
		if target.Type == "page" && target.Attached {
			log.Printf("[DEBUG] 找到已附加页面: %s", target.ID)
			return &target
		}
	}
	return nil
}

// findBrowserWindow 查找浏览器窗口
func findBrowserWindow(targets []TargetInfo) *TargetInfo {
	for _, target := range targets {
		if target.Type == "browser" {
			log.Printf("[DEBUG] 找到浏览器窗口: %s", target.ID)
			return &target
		}
	}
	return nil
}

// findFirstPage 查找第一个页面
func findFirstPage(targets []TargetInfo) *TargetInfo {
	for _, target := range targets {
		if target.Type == "page" {
			log.Printf("[DEBUG] 找到第一个页面: %s", target.ID)
			return &target
		}
	}
	return nil
}

// getPageDetails 获取页面详细信息
func getPageDetails(targetID string) (map[string]interface{}, error) {
	session, err := attachToTarget()
	if err != nil {
		return nil, fmt.Errorf("附加到target失败: %w", err)
	}
	defer detachFromTarget(session)

	result, err := evaluateJavaScriptSimple(session, `
		(function() {
			try {
				return {
					// 页面可见性
					visibilityState: document.visibilityState,
					hidden: document.hidden,
					hasFocus: document.hasFocus(),
					
					// 活动元素
					activeElement: document.activeElement ? {
						tagName: document.activeElement.tagName,
						id: document.activeElement.id,
						className: document.activeElement.className
					} : null,
					
					// 页面状态
					readyState: document.readyState,
					title: document.title,
					url: window.location.href,
					
					// 视口信息
					viewportWidth: window.innerWidth,
					viewportHeight: window.innerHeight,
					devicePixelRatio: window.devicePixelRatio,
					
					// 用户交互状态
					hasUserActivation: navigator.userActivation && navigator.userActivation.hasBeenActive,
					
					// 时间戳
					timestamp: Date.now()
				};
			} catch(e) {
				return { error: e.message };
			}
		})()
	`)

	if err != nil {
		return nil, fmt.Errorf("获取页面详情失败: %w", err)
	}

	if details, ok := result.(map[string]interface{}); ok {
		return details, nil
	}

	return map[string]interface{}{"error": "解析页面详情失败"}, nil
}

// detachFromTarget 从target分离
func detachFromTarget(sessionId string) error {
	if !DefaultBrowserWS() || chromeInstance.BrowserWSConn == nil {
		return fmt.Errorf("CDP连接不可用")
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
		return fmt.Errorf("发送分离请求失败: %w", err)
	}

	// 不等待响应，直接返回
	go func() {
		timeout := 2 * time.Second
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		for {
			select {
			case respMsg, ok := <-messageQueue:
				if !ok {
					return
				}

				if reqID == respMsg.ID {
					log.Printf("[DEBUG] 从session %s 分离完成", sessionId)
					return
				}
			case <-timer.C:
				log.Printf("[WARN] 分离session %s 超时", sessionId)
				return
			}
		}
	}()

	return nil
}

// CDPBrowserSetContentsSize 设置浏览器内容区域尺寸
// 通过 Browser.setWindowBounds 实现，可以精确控制窗口大小从而控制内容区域
/*

1.  设置窗口1的内容尺寸为800x600
response, err := CDPBrowserSetContentsSize(1, 800, 600)

2. 包括浏览器边框
response, err := CDPBrowserSetContentsSize(1, 800, 600,
	WithIncludeChrome(true),  // 包括浏览器边框
		WithKeepPosition(true),   // 保持当前位置
	)

3. 设置到指定位置
response, err := CDPBrowserSetContentsSize(1, 1024, 768,
		WithPosition(100, 100),  // 设置到屏幕位置(100, 100)
		WithIncludeChrome(true), // 包括浏览器边框
	)

4. 最大化窗口
response, err := CDPBrowserSetContentsSize(1, 1920, 1080,
		WithWindowState("maximized"),  // 最大化窗口
	)

*/
func CDPBrowserSetContentsSize(windowId int, width, height int, options ...SetSizeOption) (string, error) {
	// 1. 参数验证
	if windowId <= 0 {
		return "", fmt.Errorf("无效的窗口ID: %d", windowId)
	}

	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("宽度和高度必须是正数: %dx%d", width, height)
	}

	if width > 10000 || height > 10000 {
		return "", fmt.Errorf("尺寸超出合理范围: %dx%d", width, height)
	}

	// 2. 检查CDP连接
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}

	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 3. 解析选项
	config := &SizeConfig{
		WindowID:      windowId,
		Width:         width,
		Height:        height,
		KeepPosition:  true,  // 默认保持当前位置
		IncludeChrome: false, // 默认不包括浏览器边框
	}

	for _, option := range options {
		option(config)
	}

	// 4. 获取当前窗口边界
	currentBounds, err := CDPBrowserGetWindowBounds(windowId)
	if err != nil {
		return "", fmt.Errorf("获取当前窗口边界失败: %w", err)
	}

	// 5. 计算窗口尺寸
	windowWidth, windowHeight := calculateWindowSize(width, height, config.IncludeChrome)

	// 6. 构建新的边界设置
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := buildSetWindowBoundsMessage(reqID, windowId, windowWidth, windowHeight, config, currentBounds)

	// 7. 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("[ERROR] 发送 Browser.setWindowBounds 失败 (windowId: %d): %v", windowId, err)
		return "", fmt.Errorf("发送设置窗口边界请求失败: %w", err)
	}

	utils.Debugf("[CDP] 发送设置窗口边界请求: %s", message)

	// 8. 等待响应
	return waitForSetWindowBoundsResponse(reqID, windowId, width, height)
}

// SizeConfig 尺寸配置
type SizeConfig struct {
	WindowID      int    `json:"windowId"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	KeepPosition  bool   `json:"keepPosition"`   // 是否保持当前位置
	IncludeChrome bool   `json:"includeChrome"`  // 是否包括浏览器边框
	WindowState   string `json:"windowState"`    // 窗口状态
	Left          *int   `json:"left,omitempty"` // 指定X坐标
	Top           *int   `json:"top,omitempty"`  // 指定Y坐标
}

// SetSizeOption 配置选项
type SetSizeOption func(*SizeConfig)

// WithKeepPosition 设置是否保持当前位置
func WithKeepPosition(keep bool) SetSizeOption {
	return func(c *SizeConfig) {
		c.KeepPosition = keep
	}
}

// WithIncludeChrome 设置是否包括浏览器边框
func WithIncludeChrome(include bool) SetSizeOption {
	return func(c *SizeConfig) {
		c.IncludeChrome = include
	}
}

// WithWindowState 设置窗口状态
func WithWindowState(state string) SetSizeOption {
	return func(c *SizeConfig) {
		c.WindowState = state
	}
}

// WithPosition 设置窗口位置
func WithPosition(left, top int) SetSizeOption {
	return func(c *SizeConfig) {
		c.Left = &left
		c.Top = &top
		c.KeepPosition = false
	}
}

// calculateWindowSize 计算窗口尺寸
func calculateWindowSize(contentWidth, contentHeight int, includeChrome bool) (int, int) {
	if includeChrome {
		// 包括浏览器边框
		// 这些值是估算的，实际值可能因操作系统和浏览器主题而异
		const chromeWidth = 16   // 边框+滚动条
		const chromeHeight = 100 // 标题栏+标签栏+边框+状态栏
		return contentWidth + chromeWidth, contentHeight + chromeHeight
	}

	// 不包括浏览器边框，直接使用内容尺寸
	return contentWidth, contentHeight
}

// buildSetWindowBoundsMessage 构建设置窗口边界的消息
func buildSetWindowBoundsMessage(reqID, windowId, width, height int, config *SizeConfig, currentBounds map[string]interface{}) string {
	// 基础消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Browser.setWindowBounds",
		"params": {
			"windowId": %d,
			"bounds": {
				"width": %d,
				"height": %d`, reqID, windowId, width, height)

	// 设置窗口状态
	if config.WindowState != "" {
		message += fmt.Sprintf(`, "windowState": "%s"`, config.WindowState)
	} else {
		message += `, "windowState": "normal"`
	}

	// 设置位置
	left, top := calculatePosition(config, currentBounds)
	if left != nil {
		message += fmt.Sprintf(`, "left": %d`, *left)
	}
	if top != nil {
		message += fmt.Sprintf(`, "top": %d`, *top)
	}

	// 关闭消息
	message += `}}}`

	return message
}

// calculatePosition 计算窗口位置
func calculatePosition(config *SizeConfig, currentBounds map[string]interface{}) (*int, *int) {
	// 如果指定了位置，使用指定位置
	if config.Left != nil && config.Top != nil {
		return config.Left, config.Top
	}

	// 如果保持当前位置，使用当前位置
	if config.KeepPosition {
		if left, ok := currentBounds["left"].(float64); ok {
			if top, ok := currentBounds["top"].(float64); ok {
				leftInt := int(left)
				topInt := int(top)
				return &leftInt, &topInt
			}
		}
	}

	// 否则由系统决定位置
	return nil, nil
}

// waitForSetWindowBoundsResponse 等待设置窗口边界的响应
func waitForSetWindowBoundsResponse(reqID, windowId, width, height int) (string, error) {
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
				utils.Debugf("[CDP] Browser.setWindowBounds 收到回复 -> %s", content)

				// 检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					errorMsg := fmt.Sprintf("%v", errorObj)

					// 处理特定错误
					if strings.Contains(errorMsg, "Invalid window id") {
						return content, fmt.Errorf("无效的窗口ID: %d", windowId)
					} else if strings.Contains(errorMsg, "cannot find window") {
						return content, fmt.Errorf("找不到窗口: %d", windowId)
					} else if strings.Contains(errorMsg, "permission") {
						return content, fmt.Errorf("权限不足，无法设置窗口边界")
					}

					return content, fmt.Errorf("CDP错误: %s", errorMsg)
				}

				// 验证设置是否成功
				log.Printf("[INFO] 窗口 %d 内容尺寸已设置为: %dx%d", windowId, width, height)

				// 返回成功响应
				successResponse := map[string]interface{}{
					"success":       true,
					"windowId":      windowId,
					"contentWidth":  width,
					"contentHeight": height,
					"timestamp":     time.Now().Unix(),
				}

				successJSON, _ := json.Marshal(successResponse)
				return string(successJSON), nil
			}

		case <-timer.C:
			return "", fmt.Errorf("设置窗口边界请求超时 (%v), windowId: %d", timeout, windowId)
		}
	}
}
