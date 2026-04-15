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

// -----------------------------------------------  CSS.addRule  -----------------------------------------------

// CDPCSSAddRule 向样式表中添加新的CSS规则
// 参数说明:
//   - styleSheetId: 样式表ID
//   - ruleText: 要添加的CSS规则文本
//   - location: 可选，规则插入位置
//
// 适用场景:
// 1. 动态修改页面样式
// 2. 自动化测试中注入样式规则
// 3. 调试工具中临时修改样式
// 4. 浏览器扩展中修改页面样式
//
// 参数示例:
//
//	基础插入  {
//	 "styleSheetId": "1",
//	 "rule": "div.test { background: blue; padding: 10px; }",
//	 "index": 0
//	}
//
//	覆盖页面默认样式 {
//	 "styleSheetId": "2",
//	 "rule": "* { margin: 0; padding: 0; box-sizing: border-box; }"
//	}
//
//	动态修改元素样式 {
//	 "styleSheetId": "1",
//	 "rule": "#app { display: flex; justify-content: center; }",
//	 "index": 5
//	}
func CDPCSSAddRule(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.addRule",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 addRule 请求失败: %w", err)
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
			return "", fmt.Errorf("addRule 请求超时")
		}
	}
}

// SourceRange 源范围结构
type SourceRange struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndLine     int `json:"endLine"`
	EndColumn   int `json:"endColumn"`
}

// 创建SourceRange的便捷函数
func NewSourceRange(startLine, startColumn, endLine, endColumn int) SourceRange {
	return SourceRange{
		StartLine:   startLine,
		StartColumn: startColumn,
		EndLine:     endLine,
		EndColumn:   endColumn,
	}
}

// -----------------------------------------------  CSS.collectClassNames  -----------------------------------------------

// CDPCSSColectClassNames 从指定样式表中收集所有类名
// 这是 CDP 协议的原生方法
// 参数说明:
//   - styleSheetId: 样式表ID
func CDPCSSColectClassNames(styleSheetId string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.collectClassNames",
		"params": {
			"styleSheetId": "%s"
		}
	}`, reqID, styleSheetId)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 collectClassNames 请求失败: %w", err)
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
			return "", fmt.Errorf("collectClassNames 请求超时")
		}
	}
}

// CDPCSSEnable 启用CSS域
func CDPCSSEnable() (string, error) {
	if !DefaultNowTab(false) {
		return "", nil
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.enable"
	}`, reqID)

	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 CSS.enable 失败:", err)
		return "", err
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)
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
				fmt.Println("[CDP CSS.enable] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("CSS.enable 请求超时")
		}
	}
}

// CDPCSSDisable 禁用CSS域
func CDPCSSDisable() (string, error) {
	if !DefaultNowTab(false) {
		return "", nil
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.disable"
	}`, reqID)

	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 CSS.disable 失败:", err)
		return "", err
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)
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
				fmt.Println("[CDP CSS.disable] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("CSS.disable 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.createStyleSheet  -----------------------------------------------
// === 应用场景 ===
// 1. 开发工具: 用于调试工具中创建临时的样式表
// 2. 自动化测试: 测试中动态修改页面样式
// 3. 浏览器扩展: 扩展程序需要修改页面样式时
// 4. 样式调试: 临时应用调试样式而不修改原代码
// 5. 性能测试: 测试不同样式对页面性能的影响
// 6. 主题切换: 动态切换页面主题样式

// CDPCSSCreateStyleSheet 创建一个新的样式表
// 这是 CDP 协议的原生方法
// 参数说明:
//   - frameId: 页面框架ID
//   - force: 是否强制创建新样式表
func CDPCSSCreateStyleSheet(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.createStyleSheet",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 createStyleSheet 请求失败: %w", err)
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
			return "", fmt.Errorf("createStyleSheet 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.forcePseudoState  -----------------------------------------------
// === 应用场景 ===
// 1. 开发工具: 调试伪类样式，如:hover、:focus状态
// 2. 自动化测试: 测试元素在不同伪类状态下的样式
// 3. 视觉测试: 验证伪类状态的UI表现
// 4. 交互调试: 调试复杂的交互状态
// 5. 组件库测试: 测试组件的各种状态样式
// 6. 无障碍测试: 测试焦点状态等无障碍特性

// CDPCSSForcePseudoState 强制元素应用指定的伪类状态
// 参数说明:
//   - nodeId: DOM节点ID
//   - forcedPseudoClasses: 要强制应用的伪类数组
func CDPCSSForcePseudoState(nodeId int, forcedPseudoClasses []string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	if len(forcedPseudoClasses) == 0 {
		return "", fmt.Errorf("必须指定至少一个伪类")
	}

	// 常见的CSS伪类
	var supportedPseudoClasses = []string{
		// 链接相关
		"link", "visited", "hover", "active",

		// 用户操作
		"hover", "active", "focus", "focus-visible", "focus-within",

		// 表单状态
		"enabled", "disabled", "checked", "indeterminate", "default",
		"valid", "invalid", "in-range", "out-of-range",
		"required", "optional", "read-only", "read-write",
		"placeholder-shown",

		// 结构伪类
		"first-child", "last-child", "nth-child", "nth-last-child",
		"first-of-type", "last-of-type", "nth-of-type", "nth-last-of-type",
		"only-child", "only-of-type", "empty", "root",

		// 其他
		"target", "lang", "not", "is", "where",
	}

	// 验证伪类格式
	validPseudoClasses := map[string]bool{}

	for _, v := range supportedPseudoClasses {
		validPseudoClasses[v] = true
	}

	// 检查所有伪类是否有效
	for _, pseudoClass := range forcedPseudoClasses {
		// 移除可能的前导冒号
		className := strings.TrimPrefix(pseudoClass, ":")
		if !validPseudoClasses[className] {
			return "", fmt.Errorf("无效的伪类: %s", pseudoClass)
		}
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建伪类数组JSON
	pseudoClassesJSON, err := json.Marshal(forcedPseudoClasses)
	if err != nil {
		return "", fmt.Errorf("构建伪类数组失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.forcePseudoState",
		"params": {
			"nodeId": %d,
			"forcedPseudoClasses": %s
		}
	}`, reqID, nodeId, pseudoClassesJSON)

	// 发送请求
	err = chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 forcePseudoState 请求失败: %w", err)
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
			return "", fmt.Errorf("forcePseudoState 请求超时")
		}
	}
}

// 组件库测试 - 测试按钮的各种状态
// === 应用场景描述 ===
// 场景: UI组件库的交互状态测试
// 用途: 自动化测试按钮在各种交互状态下的样式
// 优势: 确保组件在所有状态下的视觉一致性
// 典型工作流: 遍历所有状态 -> 应用伪类 -> 验证样式 -> 生成报告
func ComponentLibraryTest(buttonNodeId int) {
	testResults := make(map[string]bool)
	// 定义要测试的所有状态
	testStates := []struct {
		name          string
		pseudoClasses []string
		description   string
	}{
		{"默认状态", []string{}, "按钮的默认外观"},
		{"悬停状态", []string{"hover"}, "鼠标悬停时的样式"},
		{"激活状态", []string{"active"}, "按钮被点击时的样式"},
		{"焦点状态", []string{"focus"}, "获得键盘焦点时的样式"},
		{"禁用状态", []string{"disabled"}, "禁用按钮的样式"},
		{"悬停+焦点", []string{"hover", "focus"}, "同时悬停和获得焦点"},
		{"激活+焦点", []string{"active", "focus"}, "点击时获得焦点"},
	}

	log.Println("开始按钮组件状态测试...")

	for _, test := range testStates {
		log.Printf("测试: %s (%s)", test.name, test.description)
		_, err := CDPCSSForcePseudoState(buttonNodeId, test.pseudoClasses)
		if err != nil {
			log.Printf("  失败: %v", err)
			testResults[test.name] = false
		} else {
			log.Printf("  成功: 状态已应用")
			testResults[test.name] = true
			// 在这里可以进行样式验证或截图
			// 延迟一下，让状态生效
			time.Sleep(200 * time.Millisecond)
		}
	}

	// 清理状态
	CDPCSSForcePseudoState(buttonNodeId, []string{})

	// 生成测试报告
	log.Println("\n=== 组件状态测试报告 ===")
	successCount := 0
	for testName, success := range testResults {
		status := "✗ 失败"
		if success {
			status = "✓ 通过"
			successCount++
		}
		log.Printf("%s: %s", testName, status)
	}

	log.Printf("测试完成: %d/%d 通过", successCount, len(testResults))
}

// -----------------------------------------------  CSS.forceStartingStyle  -----------------------------------------------
// === 应用场景 ===
// 1. CSS过渡和动画调试: 调试CSS @starting-style规则的动画起始状态
// 2. 组件库开发: 测试组件在应用starting-style时的表现
// 3. 视觉测试: 验证元素在起始样式状态下的UI表现
// 4. 动画开发: 调试复杂动画的起始状态
// 5. 性能优化: 测试starting-style对性能的影响
// 6. 浏览器兼容性测试: 测试不同浏览器对@starting-style的支持

// CDPCSSForceStartingStyle 强制元素应用起始样式状态
// 参数说明:
//   - nodeId: DOM节点ID
//   - forced: 是否启用起始样式状态
func CDPCSSForceStartingStyle(nodeId int, forced bool) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.forceStartingStyle",
		"params": {
			"nodeId": %d,
			"forced": %t
		}
	}`, reqID, nodeId, forced)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 forceStartingStyle 请求失败: %w", err)
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
			return "", fmt.Errorf("forceStartingStyle 请求超时")
		}
	}
}

// 示例: 组件测试 - 测试Modal弹窗的起始动画
// === 应用场景描述 ===
// 场景: Modal弹窗组件的起始动画测试
// 用途: 测试Modal在显示时的起始动画样式
// 优势: 自动化验证Modal的显示动画起始状态
// 典型工作流: 触发Modal显示 -> 启用起始样式 -> 截图验证 -> 测试动画
func exampleModalComponentTest(modalElementId int) {

	log.Println("测试Modal弹窗的起始动画样式...")
	// 定义测试用例
	testCases := []struct {
		name   string
		forced bool
		desc   string
	}{
		{"启用起始样式", true, "测试Modal在起始样式状态下的表现"},
		{"禁用起始样式", false, "测试Modal在正常状态下的表现"},
		{"重新启用起始样式", true, "再次测试起始样式状态"},
	}

	for i, testCase := range testCases {
		log.Printf("测试用例 %d/%d: %s", i+1, len(testCases), testCase.name)
		log.Printf("  描述: %s", testCase.desc)
		_, err := CDPCSSForceStartingStyle(modalElementId, testCase.forced)
		if err != nil {
			log.Printf("  失败: %v", err)
		} else {
			log.Printf("  成功: 状态已%s", map[bool]string{true: "启用", false: "禁用"}[testCase.forced])

			// 在这里可以进行截图或样式验证
			// screenshotName := fmt.Sprintf("modal-starting-style-%d-%t.png", i+1, testCase.forced)
			// takeScreenshot(screenshotName)

			// 验证样式
			// verifyStartingStyles(modalElementId, testCase.forced)

			// 给一些时间观察效果
			time.Sleep(1 * time.Second)
		}
	}
	// 最终清理
	log.Println("清理起始样式状态...")
	CDPCSSForceStartingStyle(modalElementId, false)
	log.Println("Modal起始样式测试完成")
}

// -----------------------------------------------  CSS.getBackgroundColors  -----------------------------------------------
// === 应用场景 ===
// 1. 无障碍测试: 验证文本与背景颜色的对比度是否符合WCAG标准
// 2. 视觉测试: 自动化验证UI元素的背景颜色
// 3. 主题系统测试: 验证主题切换时的背景颜色变化
// 4. 渐变分析: 分析CSS渐变背景的颜色分布
// 5. 颜色对比度计算: 计算文本可读性的颜色对比度
// 6. 设计系统验证: 验证设计系统中的颜色使用一致性

// CDPCSSGetBackgroundColors 获取元素背后的背景颜色范围
// 参数说明:
//   - nodeId: DOM节点ID
func CDPCSSGetBackgroundColors(nodeId int) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getBackgroundColors",
		"params": {
			"nodeId": %d
		}
	}`, reqID, nodeId)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getBackgroundColors 请求失败: %w", err)
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
			return "", fmt.Errorf("getBackgroundColors 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.getComputedStyleForNode  -----------------------------------------------
// === 应用场景 ===
// 1. 样式调试工具: 开发者工具中查看元素最终的计算样式
// 2. 自动化测试: 验证元素在特定状态下的实际样式值
// 3. 视觉回归测试: 比较不同版本间样式的变化
// 4. 响应式设计测试: 验证在不同视口下的样式计算
// 5. 样式继承分析: 分析CSS属性的继承链和计算过程
// 6. 性能分析: 监控重排和重绘相关的样式计算

// CDPCSSGetComputedStyleForNode 获取指定节点的计算样式
// 参数说明:
//   - nodeId: DOM节点ID
func CDPCSSGetComputedStyleForNode(nodeId int) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getComputedStyleForNode",
		"params": {
			"nodeId": %d
		}
	}`, reqID, nodeId)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getComputedStyleForNode 请求失败: %w", err)
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
			return "", fmt.Errorf("getComputedStyleForNode 请求超时")
		}
	}
}

// CSSComputedStyleProperty CSS计算样式属性
// 基于文档中的CSSComputedStyleProperty类型定义
type CSSComputedStyleProperty struct {
	Name  string `json:"name"`  // 计算样式属性名
	Value string `json:"value"` // 计算样式属性值
}

// ComputedStyleResult 计算样式结果
type ComputedStyleResult struct {
	ComputedStyle []CSSComputedStyleProperty `json:"computedStyle"`         // 计算样式数组
	ExtraFields   interface{}                `json:"extraFields,omitempty"` // 额外字段（实验性）
}

// ParseComputedStyle 解析计算样式响应
func ParseComputedStyle(response string) (*ComputedStyleResult, error) {
	var data struct {
		Result *ComputedStyleResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

// 示例: 自动化测试 - 样式验证
// === 应用场景描述 ===
// 场景: 自动化测试中的样式验证
// 用途: 验证元素在实际渲染后的样式值是否符合预期
// 优势: 自动化检查视觉一致性，防止样式回归
// 典型工作流: 定义预期样式 -> 获取计算样式 -> 比较差异 -> 生成报告
func exampleAutomationStyleVerification(buttonId int) {
	// 预期的样式规范
	expectedStyles := map[string]string{
		"padding":          "10px 20px",
		"border-radius":    "4px",
		"font-weight":      "500",
		"cursor":           "pointer",
		"user-select":      "none",
		"background-color": "rgb(0, 123, 255)",
		"color":            "rgb(255, 255, 255)",
	}

	log.Println("验证按钮样式符合规范...")

	// 获取计算样式
	response, err := CDPCSSGetComputedStyleForNode(buttonId)
	if err != nil {
		log.Printf("获取计算样式失败: %v", err)
		return
	}

	result, err := ParseComputedStyle(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	// 验证样式
	verificationResults := make(map[string]bool)

	for property, expectedValue := range expectedStyles {
		found := false
		for _, computedProp := range result.ComputedStyle {
			if computedProp.Name == property {
				found = true
				matches := computedProp.Value == expectedValue
				verificationResults[property] = matches

				if matches {
					log.Printf("  ✓ %s: %s (符合预期)", property, computedProp.Value)
				} else {
					log.Printf("  ✗ %s: %s (预期: %s)", property, computedProp.Value, expectedValue)
				}
				break
			}
		}

		if !found {
			log.Printf("  ✗ %s: 未找到该属性", property)
			verificationResults[property] = false
		}
	}

	// 生成验证报告
	passCount := 0
	for _, passed := range verificationResults {
		if passed {
			passCount++
		}
	}

	log.Printf("\n样式验证结果: %d/%d 通过", passCount, len(expectedStyles))
}

// 示例: 响应式设计测试
// === 应用场景描述 ===
// 场景: 响应式设计的断点测试
// 用途: 在不同视口大小下验证元素的计算样式变化
// 优势: 自动化测试响应式布局的断点切换
// 典型工作流: 设置视口大小 -> 获取计算样式 -> 验证响应式规则 -> 测试断点
func exampleResponsiveDesignTest(containerId int) {
	// 定义测试的视口大小
	viewportTests := []struct {
		name   string
		width  int
		height int
		expect struct {
			display       string
			flexDirection string
			maxWidth      string
		}
	}{
		{
			name:   "移动端 (<768px)",
			width:  375,
			height: 667,
			expect: struct {
				display       string
				flexDirection string
				maxWidth      string
			}{
				display:       "block",
				flexDirection: "column",
				maxWidth:      "100%",
			},
		},
		{
			name:   "平板端 (768px-1024px)",
			width:  768,
			height: 1024,
			expect: struct {
				display       string
				flexDirection string
				maxWidth      string
			}{
				display:       "flex",
				flexDirection: "row",
				maxWidth:      "720px",
			},
		},
		{
			name:   "桌面端 (>1024px)",
			width:  1200,
			height: 800,
			expect: struct {
				display       string
				flexDirection string
				maxWidth      string
			}{
				display:       "flex",
				flexDirection: "row",
				maxWidth:      "1140px",
			},
		},
	}

	log.Println("响应式设计断点测试...")

	for _, test := range viewportTests {
		log.Printf("\n测试视口: %s (%dx%d)", test.name, test.width, test.height)

		// 这里需要设置视口大小
		// CDPPageSetDeviceMetricsOverride(test.width, test.height)

		// 等待布局更新
		time.Sleep(500 * time.Millisecond)

		// 获取计算样式
		response, err := CDPCSSGetComputedStyleForNode(containerId)
		if err != nil {
			log.Printf("  获取计算样式失败: %v", err)
			continue
		}

		result, err := ParseComputedStyle(response)
		if err != nil {
			log.Printf("  解析结果失败: %v", err)
			continue
		}

		// 验证响应式样式
		verifyResponsiveProperty(result, "display", test.expect.display)
		verifyResponsiveProperty(result, "flex-direction", test.expect.flexDirection)
		verifyResponsiveProperty(result, "max-width", test.expect.maxWidth)
	}
}

func verifyResponsiveProperty(result *ComputedStyleResult, property, expected string) {
	for _, prop := range result.ComputedStyle {
		if prop.Name == property {
			if prop.Value == expected {
				log.Printf("  ✓ %s: %s", property, prop.Value)
			} else {
				log.Printf("  ✗ %s: %s (预期: %s)", property, prop.Value, expected)
			}
			return
		}
	}
	log.Printf("  ✗ %s: 未找到", property)
}

// -----------------------------------------------  CSS.getInlineStylesForNode  -----------------------------------------------
// === 应用场景 ===
// 1. 样式优先级分析: 分析内联样式在CSS优先级中的权重
// 2. 样式来源调试: 识别样式是来自内联属性还是外部样式表
// 3. DOM操作监控: 监控通过JavaScript动态修改的内联样式
// 4. 样式覆盖分析: 分析内联样式如何覆盖其他来源的样式
// 5. 无障碍测试: 验证内联样式对无障碍性的影响
// 6. 性能优化: 识别可能导致性能问题的内联样式
// 这个方法是构建Web性能分析工具、代码质量检查工具和样式调试工具的重要基础。

// CDPCSSGetInlineStylesForNode 获取指定节点的内联样式
// 参数说明:
//   - nodeId: DOM节点ID
func CDPCSSGetInlineStylesForNode(nodeId int) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getInlineStylesForNode",
		"params": {
			"nodeId": %d
		}
	}`, reqID, nodeId)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getInlineStylesForNode 请求失败: %w", err)
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
			return "", fmt.Errorf("getInlineStylesForNode 请求超时")
		}
	}
}

// CSSProperty CSS属性（基于文档中的CSSProperty类型定义）
type CSSProperty struct {
	Name      string       `json:"name"`                // 属性名
	Value     string       `json:"value"`               // 属性值
	Important bool         `json:"important,omitempty"` // 是否包含!important
	Implicit  bool         `json:"implicit,omitempty"`  // 是否是隐式属性
	Text      string       `json:"text,omitempty"`      // 完整的属性文本
	ParsedOk  bool         `json:"parsedOk,omitempty"`  // 是否解析成功
	Disabled  bool         `json:"disabled,omitempty"`  // 是否被禁用
	Range     *SourceRange `json:"range,omitempty"`     // 属性范围
}

// CSSStyle CSS样式对象（基于文档中的CSSStyle类型定义）
type CSSStyle struct {
	StyleSheetID     string           `json:"styleSheetId,omitempty"`     // 样式表ID
	CSSProperties    []CSSProperty    `json:"cssProperties"`              // CSS属性数组
	ShorthandEntries []ShorthandEntry `json:"shorthandEntries,omitempty"` // 简写属性条目
	CssText          string           `json:"cssText,omitempty"`          // 样式声明文本
	Range            *SourceRange     `json:"range,omitempty"`            // 样式范围
}

// ShorthandEntry 简写属性条目
type ShorthandEntry struct {
	Name      string `json:"name"`                // 简写属性名
	Value     string `json:"value"`               // 简写属性值
	Important bool   `json:"important,omitempty"` // 是否包含!important
}

// InlineStyleResult 内联样式结果
type InlineStyleResult struct {
	InlineStyle     *CSSStyle `json:"inlineStyle"`     // 内联样式（style属性）
	AttributesStyle *CSSStyle `json:"attributesStyle"` // 属性样式（DOM属性）
}

// -----------------------------------------------  CSS.getMatchedStylesForNode  -----------------------------------------------
// === 应用场景 ===
// 1. 完整的样式调试: 开发者工具中查看元素的所有样式来源和匹配规则
// 2. CSS继承链分析: 分析样式如何从父元素继承到当前元素
// 3. 伪元素样式分析: 获取::before、::after等伪元素的匹配样式
// 4. 动画规则匹配: 分析应用于元素的CSS动画和关键帧规则
// 5. 样式优先级计算: 计算CSS选择器特异性和样式优先级
// 6. CSS规则审计: 审计CSS规则的使用情况和覆盖关系

// CDPCSSGetMatchedStylesForNode 获取指定节点的匹配样式
// 参数说明:
//   - nodeId: DOM节点ID
func CDPCSSGetMatchedStylesForNode(nodeId int) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getMatchedStylesForNode",
		"params": {
			"nodeId": %d
		}
	}`, reqID, nodeId)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getMatchedStylesForNode 请求失败: %w", err)
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
			return "", fmt.Errorf("getMatchedStylesForNode 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.getMediaQueries  -----------------------------------------------
// === 应用场景 ===
// 1. 响应式设计调试: 获取页面中定义的所有媒体查询及其状态
// 2. 断点分析: 分析网站的响应式断点设计和激活状态
// 3. 媒体查询审计: 检查媒体查询的使用情况和优化空间
// 4. 性能分析: 分析媒体查询对页面渲染性能的影响
// 5. 兼容性测试: 测试不同设备下媒体查询的激活状态
// 6. 设计系统验证: 验证设计系统中的响应式断点一致性

// CDPCSSGetMediaQueries 获取所有媒体查询
func CDPCSSGetMediaQueries() (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getMediaQueries"
	}`, reqID)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getMediaQueries 请求失败: %w", err)
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
			return "", fmt.Errorf("getMediaQueries 请求超时")
		}
	}
}

// MediaQueryExpression 媒体查询表达式
type MediaQueryExpression struct {
	Value          float64      `json:"value"`                    // 媒体查询表达式值
	Unit           string       `json:"unit"`                     // 媒体查询表达式单位
	Feature        string       `json:"feature"`                  // 媒体查询表达式特性
	ValueRange     *SourceRange `json:"valueRange,omitempty"`     // 值范围
	ComputedLength float64      `json:"computedLength,omitempty"` // 计算长度
}

// MediaQuery 媒体查询
type MediaQuery struct {
	Expressions []MediaQueryExpression `json:"expressions"` // 媒体查询表达式数组
	Active      bool                   `json:"active"`      // 媒体查询条件是否满足
}

// CSSMedia CSS媒体规则
type CSSMedia struct {
	Text         string       `json:"text"`            // 媒体查询文本
	Source       string       `json:"source"`          // 媒体查询来源
	SourceURL    string       `json:"sourceURL"`       // 源URL
	Range        *SourceRange `json:"range,omitempty"` // 范围
	StyleSheetID string       `json:"styleSheetId"`    // 样式表ID
	MediaList    []MediaQuery `json:"mediaList"`       // 媒体查询列表
}

// MediaQueriesResult 媒体查询结果
type MediaQueriesResult struct {
	Medias []CSSMedia `json:"medias"` // 媒体查询数组
}

// ParseMediaQueries 解析媒体查询响应
func ParseMediaQueries(response string) (*MediaQueriesResult, error) {
	var data struct {
		Result *MediaQueriesResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

// 示例: 响应式设计调试工具 分析媒体查询
// === 应用场景描述 ===
// 场景: 响应式设计调试工具
// 用途: 获取页面中所有媒体查询，分析响应式断点的激活状态
// 优势: 实时查看媒体查询状态，调试响应式布局
// 典型工作流: 获取媒体查询 -> 分析激活状态 -> 调整视口 -> 验证变化
func ResponsiveDesignDebugger() {
	log.Println("响应式设计调试 - 分析媒体查询...")

	// 获取所有媒体查询
	response, err := CDPCSSGetMediaQueries()
	if err != nil {
		log.Printf("获取媒体查询失败: %v", err)
		return
	}

	// 解析结果
	result, err := ParseMediaQueries(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	log.Printf("=== 媒体查询分析 ===")
	log.Printf("发现 %d 个媒体查询定义", len(result.Medias))

	// 统计激活状态
	activeCount := 0
	inactiveCount := 0

	for i, media := range result.Medias {
		log.Printf("\n媒体查询[%d]:", i+1)
		log.Printf("  查询文本: %s", media.Text)
		log.Printf("  来源: %s", getMediaSourceDescription(media.Source))

		if media.SourceURL != "" {
			log.Printf("  源URL: %s", media.SourceURL)
		}

		// 分析媒体查询列表
		if len(media.MediaList) > 0 {
			for j, query := range media.MediaList {
				status := "❌ 未激活"
				if query.Active {
					status = "✅ 已激活"
					activeCount++
				} else {
					inactiveCount++
				}

				log.Printf("  条件[%d]: %s", j+1, status)

				// 显示表达式详情
				for k, expr := range query.Expressions {
					log.Printf("    表达式[%d]: %s %s %v%s",
						k+1, expr.Feature, getComparisonOperator(expr), expr.Value, expr.Unit)
				}
			}
		}
	}

	log.Printf("\n=== 统计摘要 ===")
	log.Printf("总媒体查询定义: %d", len(result.Medias))
	log.Printf("激活的条件: %d", activeCount)
	log.Printf("未激活的条件: %d", inactiveCount)
}

// 获取媒体查询来源描述
func getMediaSourceDescription(source string) string {
	switch source {
	case "mediaRule":
		return "@media规则"
	case "importRule":
		return "@import规则"
	case "linkedSheet":
		return "链接样式表"
	case "inlineSheet":
		return "内联样式表"
	default:
		return source
	}
}

// 获取比较操作符（简化处理）
func getComparisonOperator(expr MediaQueryExpression) string {
	// 文档未详述此点，但基于我所掌握的知识
	// 媒体查询表达式通常包含min-width、max-width等
	if strings.Contains(expr.Feature, "min-") {
		return ">="
	} else if strings.Contains(expr.Feature, "max-") {
		return "<="
	}
	return "="
}

// -----------------------------------------------  CSS.getPlatformFontsForNode  -----------------------------------------------
// === 应用场景 ===
// 1. 字体性能分析: 分析页面中实际使用的字体及其性能影响
// 2. 字体优化: 识别字体使用情况，优化字体加载策略
// 3. 可访问性检查: 检查字体使用是否符合可访问性要求
// 4. 设计系统验证: 验证设计系统中字体的实际使用情况
// 5. 字体调试: 调试字体渲染问题和字体回退机制
// 6. 性能监控: 监控自定义字体和系统字体的使用性能

// CDPCSSGetPlatformFontsForNode 获取节点使用的平台字体信息
// 参数说明:
//   - nodeId: DOM节点ID
func CDPCSSGetPlatformFontsForNode(nodeId int) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getPlatformFontsForNode",
		"params": {
			"nodeId": %d
		}
	}`, reqID, nodeId)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getPlatformFontsForNode 请求失败: %w", err)
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
			return "", fmt.Errorf("getPlatformFontsForNode 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.getStyleSheetText  -----------------------------------------------
// === 应用场景 ===
// 1. 样式表调试: 获取完整样式表内容用于调试和分析
// 2. 样式对比: 比较不同版本或不同环境下的样式表差异
// 3. 样式提取: 从页面中提取特定样式表的完整内容
// 4. 代码审计: 审计样式表代码质量和规范符合性
// 5. 样式分析: 分析样式表的结构、复杂度和优化空间
// 6. 样式备份: 备份页面中使用的样式表内容

// CDPCSSGetStyleSheetText 获取样式表的文本内容
// 参数说明:
//   - styleSheetId: 样式表ID
func CDPCSSGetStyleSheetText(styleSheetId string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getStyleSheetText",
		"params": {
			"styleSheetId": "%s"
		}
	}`, reqID, styleSheetId)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getStyleSheetText 请求失败: %w", err)
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
			return "", fmt.Errorf("getStyleSheetText 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setEffectivePropertyValueForNode  -----------------------------------------------
// === 应用场景 ===
// 1. 实时样式调试: 开发工具中实时修改元素的CSS属性值
// 2. 样式预览: 在不修改源码的情况下预览样式变化
// 3. 交互式设计: 允许用户交互式调整页面元素的样式
// 4. 样式测试: 测试不同属性值对元素表现的影响
// 5. 调试覆盖: 调试样式覆盖和优先级问题
// 6. 教学演示: 演示CSS属性和值如何影响元素表现

// CDPCSSSetEffectivePropertyValueForNode 设置节点有效属性的值
// 参数说明:
//   - nodeId: DOM节点ID
//   - propertyName: 要设置的属性名
//   - value: 要设置的属性值
func CDPCSSSetEffectivePropertyValueForNode(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setEffectivePropertyValueForNode",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setEffectivePropertyValueForNode 请求失败: %w", err)
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
			return "", fmt.Errorf("setEffectivePropertyValueForNode 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setKeyframeKey  -----------------------------------------------
// === 应用场景 ===
// 1. CSS动画编辑器: 可视化修改关键帧动画的时间点
// 2. 动画调试工具: 实时调整动画的关键帧位置
// 3. 动画原型设计: 快速迭代动画效果
// 4. 性能优化: 调整关键帧位置优化动画性能
// 5. 教学演示: 演示关键帧动画的工作原理
// 6. 自动化测试: 自动化修改和测试动画效果

// CDPCSSSetKeyframeKey 修改关键帧规则的关键文本
// 参数说明:
//   - styleSheetId: 样式表ID
//   - range: 源范围，指定要修改的关键帧位置
//   - keyText: 新的关键文本
func CDPCSSSetKeyframeKey(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setKeyframeKey",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setKeyframeKey 请求失败: %w", err)
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
			return "", fmt.Errorf("setKeyframeKey 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setMediaText  -----------------------------------------------
// === 应用场景 ===
// 1. 响应式设计工具: 实时修改媒体查询条件
// 2. 断点调试器: 调试和调整响应式断点
// 3. 设计系统编辑器: 编辑设计系统的响应式规则
// 4. 跨设备测试: 调整媒体查询测试不同设备
// 5. 性能优化: 优化复杂媒体查询的性能
// 6. 教学演示: 演示媒体查询的工作原理

// CDPCSSSetMediaText 修改媒体查询规则文本
// 参数说明:
//   - styleSheetId: 样式表ID
//   - range: 源范围，指定要修改的媒体查询位置
//   - text: 新的媒体查询文本
func CDPCSSSetMediaText(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setMediaText",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setMediaText 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

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
			return "", fmt.Errorf("setMediaText 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setPropertyRulePropertyName  -----------------------------------------------
// === 应用场景 ===
// 1. CSS Houdini工具: 编辑和调试CSS自定义属性
// 2. 设计系统工具: 管理设计系统的自定义属性
// 3. 样式重构工具: 批量重命名CSS自定义属性
// 4. 代码重构: 安全地重命名CSS属性规则
// 5. 教学演示: 演示CSS自定义属性的使用

// CDPCSSSetPropertyRulePropertyName 设置CSS @property规则的属性名
// 注意: 这个方法在提供的文档链接中未明确找到
// 基于方法名推测的功能
// 参数说明:
//   - styleSheetId: 样式表ID
//   - range: 源范围，指定要修改的@property规则
//   - propertyName: 新的属性名称
func CDPCSSSetPropertyRulePropertyName(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setPropertyRulePropertyName",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setPropertyRulePropertyName 请求失败: %w", err)
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
					// 检查是否是方法不存在错误
					errorCode, _ := errorObj.(map[string]interface{})["code"].(float64)
					if errorCode == -32601 { // Method not found
						return content, fmt.Errorf("CSS.setPropertyRulePropertyName 方法在当前Chrome版本中可能不存在")
					}
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setPropertyRulePropertyName 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setRuleSelector  -----------------------------------------------
// === 应用场景 ===
// 1. CSS重构工具: 安全地重命名CSS选择器
// 2. 样式调试器: 实时修改选择器测试样式匹配
// 3. 设计系统迁移: 批量更新选择器以适配新设计系统
// 4. 代码重构: 重命名类名、ID等选择器
// 5. 自动化测试: 测试不同选择器的样式应用
// 6. 教学演示: 演示CSS选择器的工作原理

// CDPCSSSetRuleSelector 修改CSS规则的选择器
// 参数说明:
//   - styleSheetId: 样式表ID
//   - range: 源范围，指定要修改的CSS规则
//   - selector: 新的选择器文本
func CDPCSSSetRuleSelector(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setRuleSelector",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setRuleSelector 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

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
			return "", fmt.Errorf("setRuleSelector 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setStyleSheetText  -----------------------------------------------
// === 应用场景 ===
// 1. 样式表编辑器: 完整的CSS编辑器和实时预览
// 2. 主题切换系统: 动态切换整个样式表实现主题变更
// 3. 代码热重载: 开发环境下实时更新样式表
// 4. 样式表优化: 应用CSS压缩、优化后的结果
// 5. 动态样式生成: 运行时生成和设置CSS样式
// 6. 样式表版本控制: 动态切换不同版本的样式表

// CDPCSSSetStyleSheetText 设置样式表的文本内容
// 参数说明:
//   - styleSheetId: 样式表ID
//   - text: 新的样式表文本
func CDPCSSSetStyleSheetText(styleSheetId, text string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}
	if text == "" {
		return "", fmt.Errorf("样式表文本不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 转义特殊字符
	escapedText := strings.ReplaceAll(text, `"`, `\"`)
	escapedText = strings.ReplaceAll(escapedText, "\n", "\\n")
	escapedText = strings.ReplaceAll(escapedText, "\t", "\\t")

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setStyleSheetText",
		"params": {
			"styleSheetId": "%s",
			"text": "%s"
		}
	}`, reqID, styleSheetId, escapedText)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setStyleSheetText 请求失败: %w", err)
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
			return "", fmt.Errorf("setStyleSheetText 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setStyleTexts  -----------------------------------------------
// === 应用场景 ===
// 1. 批量样式编辑: 一次性应用多个样式修改
// 2. 样式原子操作: 确保多个样式修改的原子性
// 3. 样式重构工具: 批量更新相关样式的值
// 4. 设计系统同步: 同步更新设计系统中的多个样式
// 5. 性能优化: 批量修改减少重排重绘次数
// 6. 样式版本管理: 批量应用样式补丁或回滚

// CDPCSSSetStyleTexts 批量设置样式文本
// 参数说明:
//   - edits: 样式声明编辑数组
func CDPCSSSetStyleTexts(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setStyleTexts",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setStyleTexts 请求失败: %w", err)
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
			return "", fmt.Errorf("setStyleTexts 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.startRuleUsageTracking  -----------------------------------------------
// === 应用场景 ===
// 1. 性能分析工具: 跟踪CSS规则的实际使用情况
// 2. 代码覆盖率分析: 分析CSS代码的实际覆盖率
// 3. 未使用CSS检测: 识别和清理未使用的CSS规则
// 4. 优化建议工具: 基于实际使用情况提供CSS优化建议
// 5. 性能监控: 监控CSS规则在页面生命周期中的使用
// 6. 开发工具: 集成到开发者工具中提供CSS使用分析

// CDPCSSStartRuleUsageTracking 开始CSS规则使用跟踪
func CDPCSSStartRuleUsageTracking() (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.startRuleUsageTracking"
	}`, reqID)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 startRuleUsageTracking 请求失败: %w", err)
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
			return "", fmt.Errorf("startRuleUsageTracking 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.stopRuleUsageTracking  -----------------------------------------------
// === 应用场景 ===
// 1. CSS覆盖率分析工具: 停止跟踪并获取最终的规则使用数据
// 2. 性能分析工具: 完成性能分析周期，收集分析结果
// 3. 代码质量工具: 收集CSS代码使用情况的最终数据
// 4. 开发工具集成: 在开发者工具中停止CSS跟踪
// 5. 自动化测试: 在测试完成后收集CSS使用数据
// 6. 性能监控: 完成监控周期，生成报告

// CDPCSSStopRuleUsageTracking 停止CSS规则使用跟踪
func CDPCSSStopRuleUsageTracking() (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.stopRuleUsageTracking"
	}`, reqID)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stopRuleUsageTracking 请求失败: %w", err)
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
			return "", fmt.Errorf("stopRuleUsageTracking 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.takeCoverageDelta  -----------------------------------------------
// === 应用场景 ===
// 1. 增量覆盖率监控: 持续监控CSS规则使用情况的变化
// 2. 实时性能分析: 实时分析用户交互过程中的CSS使用变化
// 3. 交互式调试工具: 在开发者工具中实时显示CSS覆盖率变化
// 4. 自动化测试集成: 在测试过程中持续监控CSS覆盖率
// 5. 渐进式分析: 分阶段分析CSS规则使用情况
// 6. 性能基准测试: 比较不同时间点的CSS使用情况

// CDPCSSTakeCoverageDelta 获取CSS覆盖率增量数据
func CDPCSSTakeCoverageDelta() (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.takeCoverageDelta"
	}`, reqID)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 takeCoverageDelta 请求失败: %w", err)
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
			return "", fmt.Errorf("takeCoverageDelta 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.getEnvironmentVariables  -----------------------------------------------
// === 应用场景 ===
// 1. CSS环境调试工具: 查看和分析CSS中定义的环境变量
// 2. 设计系统工具: 监控和验证设计系统的环境变量配置
// 3. 响应式设计调试: 分析响应式设计中的环境变量使用
// 4. 样式变量管理: 管理CSS自定义属性的环境变量
// 5. 跨平台样式适配: 分析不同平台下的环境变量差异
// 6. 开发工具集成: 在开发者工具中显示CSS环境变量

// CDPCSSGetEnvironmentVariables 获取CSS环境变量
func CDPCSSGetEnvironmentVariables() (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.getEnvironmentVariables"
	}`, reqID)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getEnvironmentVariables 请求失败: %w", err)
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
			return "", fmt.Errorf("getEnvironmentVariables 请求超时")
		}
	}
}

// -----------------------------------------------  CSS.setContainerQueryText  -----------------------------------------------
// === 应用场景 ===
// 1. 容器查询编辑器: 实时编辑和调试CSS容器查询
// 2. 响应式组件工具: 调整组件级的响应式断点
// 3. 设计系统适配器: 编辑容器查询以适配不同组件
// 4. 性能优化工具: 优化容器查询的性能
// 5. 教学演示工具: 演示容器查询的工作原理
// 6. 跨组件调试: 调试多个组件间的容器查询交互

// CDPCSSSetContainerQueryText 修改容器查询文本
// 参数说明:
//   - styleSheetId: 样式表ID
//   - range: 源范围，指定要修改的容器查询位置
//   - text: 新的容器查询文本
func CDPCSSSetContainerQueryText(params string) (string, error) {
	if !DefaultNowTab(false) {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.NowTabWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setContainerQueryText",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setContainerQueryText 请求失败: %w", err)
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
			return "", fmt.Errorf("setContainerQueryText 请求超时")
		}
	}
}
