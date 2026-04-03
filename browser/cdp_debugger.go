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

// -----------------------------------------------  Debugger.continueToLocation  -----------------------------------------------
// === 应用场景 ===
// 1. 精确调试: 在代码执行到特定位置时暂停
// 2. 条件断点替代: 替代复杂的条件断点逻辑
// 3. 循环调试: 调试特定循环的特定迭代
// 4. 异步代码调试: 调试异步操作的完成点
// 5. 事件处理调试: 在特定事件处理程序中暂停
// 6. 性能分析: 分析特定代码路径的执行

// CDPDebuggerContinueToLocation 继续执行直到到达特定位置
// 参数说明:
//   - scriptId: 脚本ID
//   - lineNumber: 行号（0-based）
//   - columnNumber: 列号（0-based）
func CDPDebuggerContinueToLocation(scriptId string, lineNumber, columnNumber int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if scriptId == "" {
		return "", fmt.Errorf("脚本ID不能为空")
	}
	if lineNumber < 0 {
		return "", fmt.Errorf("行号必须是非负整数")
	}
	if columnNumber < 0 {
		return "", fmt.Errorf("列号必须是非负整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建位置对象
	location := map[string]interface{}{
		"scriptId":     scriptId,
		"lineNumber":   lineNumber,
		"columnNumber": columnNumber,
	}

	locationJSON, err := json.Marshal(location)
	if err != nil {
		return "", fmt.Errorf("序列化位置对象失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.continueToLocation",
		"params": {
			"location": %s
		}
	}`, reqID, locationJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 continueToLocation 请求失败: %w", err)
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
			return "", fmt.Errorf("continueToLocation 请求超时")
		}
	}
}

/*

// 示例1: 精确位置调试
func examplePreciseLocationDebugging() {
	// === 应用场景描述 ===
	// 场景: 精确位置调试
	// 用途: 在代码执行到特定行和列时暂停
	// 优势: 可以精确定位到特定的表达式或语句
	// 典型工作流: 设置目标位置 -> 继续执行 -> 等待暂停 -> 调试分析

	log.Println("精确位置调试示例...")

	// 假设我们要调试的脚本ID
	scriptId := "script-1234567890"

	// 定义要调试的目标位置
	debugLocations := []struct {
		line   int
		column int
		description string
	}{
		{10, 0, "函数开始位置"},
		{15, 4, "条件语句内部"},
		{20, 8, "循环体内部"},
		{25, 12, "函数调用前"},
		{30, 0, "错误处理代码"},
	}

	for i, loc := range debugLocations {
		log.Printf("调试位置 %d/%d: %s", i+1, len(debugLocations), loc.description)
		log.Printf("  目标位置: 脚本 %s, 行 %d, 列 %d", scriptId, loc.line, loc.column)

		// 继续执行到目标位置
		response, err := CDPDebuggerContinueToLocation(scriptId, loc.line, loc.column)
		if err != nil {
			log.Printf("  继续执行失败: %v", err)
			continue
		}

		log.Printf("  继续执行成功，等待到达目标位置...")

		// 这里应该监听 Debugger.paused 事件
		// 当执行到达目标位置时，会触发 paused 事件

		// 模拟等待到达目标位置的时间
		time.Sleep(1 * time.Second)

		log.Printf("  已到达目标位置")
	}
}

// 示例2: 条件执行路径调试
func exampleConditionalExecutionPathDebugging() {
	// === 应用场景描述 ===
	// 场景: 条件执行路径调试
	// 用途: 调试特定条件下的代码执行路径
	// 优势: 可以测试不同分支的代码执行
	// 典型工作流: 设置条件断点位置 -> 触发条件 -> 继续执行 -> 调试特定路径

	log.Println("条件执行路径调试示例...")

	// 定义条件断点位置
	conditionalBreakpoints := []struct {
		scriptId    string
		line        int
		column      int
		condition   string
		testData    string
	}{
		{
			scriptId:  "app-utils.js",
			line:      42,
			column:    0,
			condition: "用户登录状态检查",
			testData:  "user.isLoggedIn === true",
		},
		{
			scriptId:  "form-validation.js",
			line:      78,
			column:    4,
			condition: "表单验证失败分支",
			testData:  "validationResult.errors.length > 0",
		},
		{
			scriptId:  "api-handler.js",
			line:      105,
			column:    8,
			condition: "API请求成功处理",
			testData:  "response.status === 200",
		},
		{
			scriptId:  "error-handler.js",
			line:      33,
			column:    0,
			condition: "异常捕获处理",
			testData:  "error.code === 'NETWORK_ERROR'",
		},
		{
			scriptId:  "pagination.js",
			line:      56,
			column:    12,
			condition: "最后一页判断",
			testData:  "currentPage >= totalPages",
		},
	}

	// 模拟不同的测试场景
	testScenarios := []struct {
		name     string
		setup    func()
		cleanup  func()
	}{
		{
			name: "用户已登录场景",
			setup: func() {
				log.Println("设置用户已登录状态...")
				// 这里可以设置测试数据
			},
			cleanup: func() {
				log.Println("清理用户登录状态...")
			},
		},
		{
			name: "表单验证失败场景",
			setup: func() {
				log.Println("设置表单无效数据...")
			},
			cleanup: func() {
				log.Println("清理表单数据...")
			},
		},
		{
			name: "API成功响应场景",
			setup: func() {
				log.Println("模拟API成功响应...")
			},
			cleanup: func() {
				log.Println("清理API响应模拟...")
			},
		},
	}

	for scenarioIndex, scenario := range testScenarios {
		log.Printf("\n=== 测试场景 %d/%d: %s ===",
			scenarioIndex+1, len(testScenarios), scenario.name)

		// 设置测试场景
		scenario.setup()

		// 为每个条件断点位置继续执行
		for bpIndex, bp := range conditionalBreakpoints {
			if bpIndex < scenarioIndex+2 { // 每个场景测试2个断点
				log.Printf("调试断点 %d: %s", bpIndex+1, bp.condition)
				log.Printf("  位置: %s:%d:%d", bp.scriptId, bp.line, bp.column)
				log.Printf("  测试数据: %s", bp.testData)

				// 继续执行到目标位置
				response, err := CDPDebuggerContinueToLocation(bp.scriptId, bp.line, bp.column)
				if err != nil {
					log.Printf("  继续执行失败: %v", err)
				} else {
					log.Printf("  继续执行成功，等待到达断点...")

					// 这里可以检查是否真的到达了预期位置
					// 通过监听 Debugger.paused 事件来判断

					time.Sleep(500 * time.Millisecond)
					log.Printf("  已到达条件断点")
				}
			}
		}

		// 清理测试场景
		scenario.cleanup()

		// 短暂延迟
		time.Sleep(300 * time.Millisecond)
	}
}

// 示例3: 异步代码调试
func exampleAsyncCodeDebugging() {
	// === 应用场景描述 ===
	// 场景: 异步代码调试
	// 用途: 调试异步操作完成后的代码执行
	// 优势: 可以调试Promise、async/await、setTimeout等异步代码
	// 典型工作流: 设置异步回调位置 -> 触发异步操作 -> 继续执行 -> 调试回调代码

	log.Println("异步代码调试示例...")

	// 定义常见的异步代码位置
	asyncCodeLocations := []struct {
		name        string
		scriptId    string
		line        int
		column      int
		asyncType   string
		description string
	}{
		{
			name:      "Promise.then回调",
			scriptId:  "api-service.js",
			line:      88,
			column:    4,
			asyncType: "Promise",
			description: "API响应处理回调",
		},
		{
			name:      "async函数await后",
			scriptId:  "data-loader.js",
			line:      52,
			column:    8,
			asyncType: "async/await",
			description: "异步数据加载完成后的处理",
		},
		{
			name:      "setTimeout回调",
			scriptId:  "animation-controller.js",
			line:      120,
			column:    0,
			asyncType: "setTimeout",
			description: "动画延迟执行回调",
		},
		{
			name:      "事件监听器",
			scriptId:  "event-handler.js",
			line:      45,
			column:    12,
			asyncType: "Event",
			description: "DOM事件处理回调",
		},
		{
			name:      "fetch响应处理",
			scriptId:  "http-client.js",
			line:      76,
			column:    16,
			asyncType: "fetch",
			description: "网络请求响应处理",
		},
		{
			name:      "微任务队列",
			scriptId:  "queue-processor.js",
			line:      33,
			column:    0,
			asyncType: "microtask",
			description: "微任务队列处理",
		},
	}

	// 模拟不同的异步操作
	log.Println("开始异步操作调试...")

	for i, location := range asyncCodeLocations {
		log.Printf("\n异步调试 %d/%d: %s", i+1, len(asyncCodeLocations), location.name)
		log.Printf("  异步类型: %s", location.asyncType)
		log.Printf("  描述: %s", location.description)
		log.Printf("  目标位置: %s:%d:%d", location.scriptId, location.line, location.column)

		// 模拟触发异步操作
		log.Printf("  触发异步操作...")

		// 继续执行到异步回调位置
		response, err := CDPDebuggerContinueToLocation(location.scriptId, location.line, location.column)
		if err != nil {
			log.Printf("  继续执行失败: %v", err)
			continue
		}

		log.Printf("  继续执行成功，等待异步操作完成...")

		// 模拟异步操作完成的时间
		waitTime := 300 + i*100 // 递增的等待时间
		time.Sleep(time.Duration(waitTime) * time.Millisecond)

		log.Printf("  异步操作完成，已到达回调位置")

		// 这里可以添加调试逻辑，比如检查变量状态、调用栈等
	}
}

*/

// -----------------------------------------------  Debugger.disable  -----------------------------------------------

// CDPDebuggerDisable 禁用调试器
func CDPDebuggerDisable() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.disable"
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
			return "", fmt.Errorf("disable 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.enable  -----------------------------------------------

// CDPDebuggerEnable 启用调试器
// 参数说明:
//   - maxScriptsCacheSize: 收集脚本的最大缓存大小（字节），0表示无限制
func CDPDebuggerEnable(maxScriptsCacheSize int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	params := ""
	if maxScriptsCacheSize > 0 {
		params = fmt.Sprintf(`"maxScriptsCacheSize": %d`, maxScriptsCacheSize)
	}

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.enable"%s
	}`, reqID, func() string {
		if params != "" {
			return fmt.Sprintf(`, "params": {%s}`, params)
		}
		return ""
	}())

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
			return "", fmt.Errorf("enable 请求超时")
		}
	}
}

/*

// 示例1: 完整的调试会话生命周期管理
func exampleDebugSessionLifecycle() {
	// === 应用场景描述 ===
	// 场景: 完整的调试会话生命周期管理
	// 用途: 管理从启用到禁用的完整调试流程
	// 优势: 确保调试资源的正确初始化和清理
	// 典型工作流: 启用调试器 -> 执行调试操作 -> 禁用调试器 -> 清理资源

	log.Println("调试会话生命周期管理示例...")

	// 1. 启用调试器
	log.Println("步骤1: 启用调试器")
	enableResponse, err := CDPDebuggerEnable(1024 * 1024) // 1MB缓存
	if err != nil {
		log.Printf("启用调试器失败: %v", err)
		return
	}

	// 解析返回的调试器ID
	type EnableResult struct {
		DebuggerID string `json:"debuggerId"`
	}

	var data struct {
		Result EnableResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(enableResponse), &data); err == nil {
		log.Printf("调试器已启用，ID: %s", data.Result.DebuggerID)
	} else {
		log.Printf("启用调试器成功，但解析ID失败: %v", err)
	}

	// 2. 模拟调试操作
	log.Println("步骤2: 执行调试操作")

	// 这里可以执行各种调试操作，如设置断点、单步调试等
	debugOperations := []struct {
		name        string
		operation   func() error
		description string
	}{
		{
			name:        "设置断点",
			description: "在关键位置设置断点",
			operation: func() error {
				// 模拟设置断点操作
				log.Println("  设置断点操作...")
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		},
		{
			name:        "暂停执行",
			description: "暂停JavaScript执行",
			operation: func() error {
				// 模拟暂停操作
				log.Println("  暂停执行操作...")
				time.Sleep(150 * time.Millisecond)
				return nil
			},
		},
		{
			name:        "检查变量",
			description: "检查当前作用域的变量",
			operation: func() error {
				log.Println("  检查变量操作...")
				time.Sleep(100 * time.Millisecond)
				return nil
			},
		},
		{
			name:        "单步调试",
			description: "逐行执行代码",
			operation: func() error {
				log.Println("  单步调试操作...")
				time.Sleep(250 * time.Millisecond)
				return nil
			},
		},
	}

	for i, op := range debugOperations {
		log.Printf("调试操作 %d/%d: %s", i+1, len(debugOperations), op.name)
		log.Printf("  描述: %s", op.description)

		if err := op.operation(); err != nil {
			log.Printf("  操作失败: %v", err)
		} else {
			log.Printf("  操作成功")
		}
	}

	// 3. 禁用调试器
	log.Println("步骤3: 禁用调试器")
	disableResponse, err := CDPDebuggerDisable()
	if err != nil {
		log.Printf("禁用调试器失败: %v", err)
	} else {
		log.Printf("调试器已禁用: %s", disableResponse)
	}

	log.Println("调试会话生命周期管理完成")
}

// 示例2: 错误恢复和资源清理
func exampleErrorRecoveryAndCleanup() {
	// === 应用场景描述 ===
	// 场景: 错误恢复和资源清理
	// 用途: 在调试过程中发生错误时确保资源正确清理
	// 优势: 防止资源泄漏，确保系统稳定性
	// 典型工作流: 尝试调试 -> 捕获错误 -> 清理资源 -> 恢复状态

	log.Println("错误恢复和资源清理示例...")

	// 模拟可能的错误场景
	errorScenarios := []struct {
		name        string
		description string
		simulateError func() error
		recovery    func() error
	}{
		{
			name:        "调试器初始化失败",
			description: "启用调试器时发生错误",
			simulateError: func() error {
				log.Println("  模拟调试器初始化失败...")
				return fmt.Errorf("无法连接到调试接口")
			},
			recovery: func() error {
				log.Println("  执行恢复操作: 清理部分初始化的资源")
				// 即使启用失败，也尝试禁用以确保状态一致
				_, _ = CDPDebuggerDisable()
				return nil
			},
		},
		{
			name:        "调试操作超时",
			description: "调试操作执行时间过长",
			simulateError: func() error {
				log.Println("  模拟调试操作超时...")
				return fmt.Errorf("操作执行超时")
			},
			recovery: func() error {
				log.Println("  执行恢复操作: 强制停止当前调试操作并清理")
				_, err := CDPDebuggerDisable()
				time.Sleep(100 * time.Millisecond)
				// 重新启用调试器
				_, _ = CDPDebuggerEnable(0)
				return err
			},
		},
		{
			name:        "内存资源不足",
			description: "调试器缓存超出限制",
			simulateError: func() error {
				log.Println("  模拟内存资源不足错误...")
				return fmt.Errorf("脚本缓存超出限制")
			},
			recovery: func() error {
				log.Println("  执行恢复操作: 禁用并重新启用调试器，使用更小缓存")
				_, err := CDPDebuggerDisable()
				if err != nil {
					return err
				}
				// 使用更小的缓存重新启用
				_, err = CDPDebuggerEnable(512 * 1024) // 512KB
				return err
			},
		},
		{
			name:        "网络连接中断",
			description: "与调试目标的连接断开",
			simulateError: func() error {
				log.Println("  模拟网络连接中断...")
				return fmt.Errorf("调试连接已断开")
			},
			recovery: func() error {
				log.Println("  执行恢复操作: 清理连接资源")
				// 尝试禁用，但可能会因为连接已断开而失败
				_, _ = CDPDebuggerDisable()
				log.Println("  连接资源已清理")
				return nil
			},
		},
	}

	for scenarioIndex, scenario := range errorScenarios {
		log.Printf("\n=== 错误场景 %d/%d: %s ===",
			scenarioIndex+1, len(errorScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)

		// 先启用调试器
		log.Println("启用调试器...")
		_, enableErr := CDPDebuggerEnable(1024 * 1024)
		if enableErr != nil {
			log.Printf("启用调试器失败: %v", enableErr)
			continue
		}

		// 模拟错误
		log.Println("执行调试操作...")
		err := scenario.simulateError()
		if err != nil {
			log.Printf("发生错误: %v", err)
		}

		// 执行恢复操作
		log.Println("执行错误恢复...")
		if recoveryErr := scenario.recovery(); recoveryErr != nil {
			log.Printf("恢复失败: %v", recoveryErr)

			// 强制清理
			log.Println("执行强制清理...")
			_, _ = CDPDebuggerDisable()
		} else {
			log.Println("恢复成功")
		}

		// 确保调试器被禁用
		log.Println("确保调试器被禁用...")
		_, _ = CDPDebuggerDisable()

		// 短暂延迟
		time.Sleep(200 * time.Millisecond)
	}

	log.Println("错误恢复和资源清理示例完成")
}

// 示例3: 生产环境调试器管理
func exampleProductionEnvironmentDebuggerManagement() {
	// === 应用场景描述 ===
	// 场景: 生产环境调试器管理
	// 用途: 在生产环境中安全地启用和禁用调试器
	// 优势: 确保生产环境的性能和安全性
	// 典型工作流: 按需启用 -> 快速调试 -> 立即禁用 -> 性能监控

	log.Println("生产环境调试器管理示例...")

	// 生产环境管理策略
	productionStrategies := []struct {
		name        string
		description string
		maxCache    int
		timeout     time.Duration
		shouldCleanup bool
	}{
		{
			name:        "最小化调试",
			description: "使用最小缓存，最短超时时间",
			maxCache:    256 * 1024, // 256KB
			timeout:     1 * time.Second,
			shouldCleanup: true,
		},
		{
			name:        "性能优先",
			description: "平衡调试能力和性能影响",
			maxCache:    512 * 1024, // 512KB
			timeout:     3 * time.Second,
			shouldCleanup: true,
		},
		{
			name:        "调试优先",
			description: "提供更好的调试体验，适当牺牲性能",
			maxCache:    1024 * 1024, // 1MB
			timeout:     5 * time.Second,
			shouldCleanup: false, // 可能需要保持状态
		},
		{
			name:        "紧急调试",
			description: "用于生产环境紧急问题诊断",
			maxCache:    2048 * 1024, // 2MB
			timeout:     10 * time.Second,
			shouldCleanup: true,
		},
	}

	for strategyIndex, strategy := range productionStrategies {
		log.Printf("\n=== 生产环境策略 %d/%d: %s ===",
			strategyIndex+1, len(productionStrategies), strategy.name)
		log.Printf("描述: %s", strategy.description)
		log.Printf("配置: 缓存=%dKB, 超时=%v", strategy.maxCache/1024, strategy.timeout)

		// 启用调试器
		log.Printf("启用调试器 (缓存: %dKB)...", strategy.maxCache/1024)
		startTime := time.Now()

		enableResponse, enableErr := CDPDebuggerEnable(strategy.maxCache)
		enableDuration := time.Since(startTime)

		if enableErr != nil {
			log.Printf("启用失败: %v (耗时: %v)", enableErr, enableDuration)
			continue
		}

		log.Printf("启用成功 (耗时: %v)", enableDuration)

		// 模拟生产环境调试操作
		log.Println("执行生产环境调试操作...")
		debugStartTime := time.Now()

		// 这里执行关键的调试操作
		criticalDebugOperations := []string{
			"检查当前执行栈",
			"查看关键变量状态",
			"分析内存使用情况",
			"检查网络请求状态",
		}

		for i, op := range criticalDebugOperations {
			log.Printf("  调试操作 %d: %s", i+1, op)
			time.Sleep(100 * time.Millisecond)
		}

		debugDuration := time.Since(debugStartTime)
		log.Printf("调试操作完成 (耗时: %v)", debugDuration)

		// 检查是否超时
		if debugDuration > strategy.timeout {
			log.Printf("⚠ 调试操作超时! 配置超时: %v, 实际耗时: %v",
				strategy.timeout, debugDuration)

			// 立即清理
			log.Println("立即清理调试器...")
			_, _ = CDPDebuggerDisable()
		} else if strategy.shouldCleanup {
			// 正常清理
			log.Println("清理调试器...")
			disableStartTime := time.Now()
			_, disableErr := CDPDebuggerDisable()
			disableDuration := time.Since(disableStartTime)

			if disableErr != nil {
				log.Printf("清理失败: %v (耗时: %v)", disableErr, disableDuration)
			} else {
				log.Printf("清理成功 (耗时: %v)", disableDuration)
			}
		} else {
			log.Println("保持调试器启用状态（根据策略配置）")
		}

		// 生成性能报告
		totalDuration := time.Since(startTime)
		log.Printf("总耗时: %v", totalDuration)

		// 性能评估
		log.Println("性能评估:")
		if totalDuration > 5*time.Second {
			log.Printf("  ❌ 性能不佳: 总耗时过长 (%v)", totalDuration)
		} else if totalDuration > 2*time.Second {
			log.Printf("  ⚠ 性能中等: 总耗时适中 (%v)", totalDuration)
		} else {
			log.Printf("  ✅ 性能良好: 总耗时较短 (%v)", totalDuration)
		}
	}

	log.Println("生产环境调试器管理示例完成")
}

*/

// -----------------------------------------------  Debugger.evaluateOnCallFrame  -----------------------------------------------

// === 应用场景 ===
// 1. 变量检查: 在断点处检查变量的值和状态
// 2. 表达式测试: 测试代码片段在当前上下文中的行为
// 3. 调试工具: 构建交互式调试工具
// 4. 条件求值: 在断点条件下求值复杂表达式
// 5. 运行时修改: 临时修改变量值以测试不同场景
// 6. 性能分析: 分析表达式在特定上下文中的性能

// CDPDebuggerEvaluateOnCallFrame 在指定调用帧上求值表达式
// 参数说明:
//   - callFrameId: 调用帧ID
//   - expression: 要求值的JavaScript表达式
//   - options: 求值选项配置
func CDPDebuggerEvaluateOnCallFrame(callFrameId string, expression string, options *EvaluateOptions) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if callFrameId == "" {
		return "", fmt.Errorf("调用帧ID不能为空")
	}
	if expression == "" {
		return "", fmt.Errorf("表达式不能为空")
	}

	// 使用默认选项
	if options == nil {
		options = &EvaluateOptions{
			ObjectGroup:           "debugger-eval",
			IncludeCommandLineAPI: false,
			Silent:                false,
			ReturnByValue:         false,
			GeneratePreview:       true,
			ThrowOnSideEffect:     true,
		}
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 转义表达式中的特殊字符
	escapedExpression := strings.ReplaceAll(expression, `"`, `\"`)
	escapedExpression = strings.ReplaceAll(escapedExpression, "\n", "\\n")
	escapedExpression = strings.ReplaceAll(escapedExpression, "\t", "\\t")

	// 构建参数对象
	params := map[string]interface{}{
		"callFrameId": callFrameId,
		"expression":  escapedExpression,
	}

	// 添加可选参数
	if options.ObjectGroup != "" {
		params["objectGroup"] = options.ObjectGroup
	}
	if options.IncludeCommandLineAPI {
		params["includeCommandLineAPI"] = true
	}
	if options.Silent {
		params["silent"] = true
	}
	if options.ReturnByValue {
		params["returnByValue"] = true
	}
	if options.GeneratePreview {
		params["generatePreview"] = true
	}
	if options.ThrowOnSideEffect {
		params["throwOnSideEffect"] = true
	}
	if options.Timeout > 0 {
		params["timeout"] = options.Timeout.Milliseconds()
	}

	// 序列化参数
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.evaluateOnCallFrame",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 evaluateOnCallFrame 请求失败: %w", err)
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
			return "", fmt.Errorf("evaluateOnCallFrame 请求超时")
		}
	}
}

type EvaluateOptions struct {
	ObjectGroup           string        // 对象组名称
	IncludeCommandLineAPI bool          // 是否包含命令行API
	Silent                bool          // 是否静默模式
	ReturnByValue         bool          // 是否按值返回
	GeneratePreview       bool          // 是否生成预览
	ThrowOnSideEffect     bool          // 有副作用时是否抛出异常
	Timeout               time.Duration // 超时时间
}

// EvaluateResult 求值结果
type EvaluateResult struct {
	Result           *RuntimeRemoteObject     `json:"result,omitempty"`           // 求值结果
	ExceptionDetails *RuntimeExceptionDetails `json:"exceptionDetails,omitempty"` // 异常详情
}

// RuntimeRemoteObject 运行时远程对象
// 表示在远程调试器中执行的JavaScript对象
type RuntimeRemoteObject struct {
	Type          string               `json:"type"`                    // 对象类型: "object", "function", "undefined", "string", "number", "boolean", "symbol", "bigint"
	Subtype       string               `json:"subtype,omitempty"`       // 子类型: "array", "null", "node", "regexp", "date", "map", "set", "weakmap", "weakset", "iterator", "generator", "error", "proxy", "promise", "typedarray", "arraybuffer", "dataview"
	ClassName     string               `json:"className,omitempty"`     // 对象类名
	Value         interface{}          `json:"value,omitempty"`         // 原始值（当type是原始类型时）
	Description   string               `json:"description,omitempty"`   // 对象描述
	ObjectID      string               `json:"objectId,omitempty"`      // 对象标识符
	Preview       *ObjectPreview       `json:"preview,omitempty"`       // 对象预览
	CustomPreview *CustomObjectPreview `json:"customPreview,omitempty"` // 自定义预览
}

// ObjectPreview 对象预览
type ObjectPreview struct {
	Type        string            `json:"type"`                  // 对象类型
	Subtype     string            `json:"subtype,omitempty"`     // 子类型
	Description string            `json:"description,omitempty"` // 描述
	Overflow    bool              `json:"overflow"`              // 是否溢出
	Properties  []PropertyPreview `json:"properties"`            // 属性预览
	Entries     []EntryPreview    `json:"entries,omitempty"`     // 条目预览
}

// PropertyPreview 属性预览
type PropertyPreview struct {
	Name         string         `json:"name"`                   // 属性名
	Type         string         `json:"type"`                   // 类型
	Value        string         `json:"value,omitempty"`        // 值
	ValuePreview *ObjectPreview `json:"valuePreview,omitempty"` // 值预览
	Subtype      string         `json:"subtype,omitempty"`      // 子类型
}

// EntryPreview 条目预览
type EntryPreview struct {
	Key   *ObjectPreview `json:"key,omitempty"`   // 键预览
	Value *ObjectPreview `json:"value,omitempty"` // 值预览
}

// CustomObjectPreview 自定义对象预览
type CustomObjectPreview struct {
	Header       string `json:"header"`                 // 头部
	BodyGetterID string `json:"bodyGetterId,omitempty"` // 主体获取器ID
}

// RuntimeExceptionDetails 运行时异常详情
type RuntimeExceptionDetails struct {
	ExceptionID        int                  `json:"exceptionId"`                  // 异常ID
	Text               string               `json:"text"`                         // 异常文本
	LineNumber         int                  `json:"lineNumber,omitempty"`         // 行号
	ColumnNumber       int                  `json:"columnNumber,omitempty"`       // 列号
	ScriptID           string               `json:"scriptId,omitempty"`           // 脚本ID
	URL                string               `json:"url,omitempty"`                // URL
	StackTrace         *StackTrace          `json:"stackTrace,omitempty"`         // 堆栈跟踪
	Exception          *RuntimeRemoteObject `json:"exception,omitempty"`          // 异常对象
	ExecutionContextID int                  `json:"executionContextId,omitempty"` // 执行上下文ID
	ExceptionMetaData  map[string]string    `json:"exceptionMetaData,omitempty"`  // 异常元数据
}

// StackTrace 堆栈跟踪
type StackTrace struct {
	Description          string      `json:"description,omitempty"`          // 描述
	CallFrames           []CallFrame `json:"callFrames"`                     // 调用帧
	Parent               *StackTrace `json:"parent,omitempty"`               // 父堆栈
	PromiseCreationFrame *CallFrame  `json:"promiseCreationFrame,omitempty"` // Promise创建帧
}

// CallFrame 调用帧
type CallFrame struct {
	FunctionName string `json:"functionName"`  // 函数名
	ScriptID     string `json:"scriptId"`      // 脚本ID
	URL          string `json:"url,omitempty"` // URL
	LineNumber   int    `json:"lineNumber"`    // 行号
	ColumnNumber int    `json:"columnNumber"`  // 列号
}

// ParseEvaluateOnCallFrame 解析求值响应
func ParseEvaluateOnCallFrame(response string) (*EvaluateResult, error) {
	var data struct {
		Result *EvaluateResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

/*

// 示例1: 调用帧变量检查器
func exampleCallFrameVariableInspector() {
	// === 应用场景描述 ===
	// 场景: 调用帧变量检查器
	// 用途: 在调试时检查特定调用帧中的变量
	// 优势: 可以查看当前执行上下文中的变量状态
	// 典型工作流: 暂停执行 -> 选择调用帧 -> 求值表达式 -> 分析结果

	log.Println("调用帧变量检查器示例...")

	// 模拟不同的调用帧ID
	callFrameIds := []string{
		"call-frame-123",
		"call-frame-456",
		"call-frame-789",
		"call-frame-101",
		"call-frame-112",
	}

	// 定义要检查的变量和表达式
	inspectionExpressions := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "局部变量",
			expression:  "localVar",
			description: "检查局部变量的值",
		},
		{
			name:        "函数参数",
			expression:  "arguments.length",
			description: "检查函数参数数量",
		},
		{
			name:        "对象属性",
			expression:  "this.propertyName",
			description: "检查this对象的属性",
		},
		{
			name:        "闭包变量",
			expression:  "closureVar",
			description: "检查闭包中的变量",
		},
		{
			name:        "全局变量",
			expression:  "window.location.href",
			description: "检查全局变量",
		},
		{
			name:        "计算表达式",
			expression:  "2 + 2 * 3",
			description: "测试数学表达式求值",
		},
		{
			name:        "类型检查",
			expression:  "typeof variable",
			description: "检查变量类型",
		},
		{
			name:        "函数调用",
			expression:  "JSON.stringify(obj, null, 2)",
			description: "测试函数调用",
		},
	}

	// 测试不同调用帧
	for frameIndex, callFrameId := range callFrameIds {
		log.Printf("\n=== 检查调用帧 %d/%d: %s ===",
			frameIndex+1, len(callFrameIds), callFrameId)

		// 为每个调用帧检查几个表达式
		expressionsToTest := inspectionExpressions[frameIndex*2:]
		if len(expressionsToTest) > 3 {
			expressionsToTest = expressionsToTest[:3]
		}

		for exprIndex, expr := range expressionsToTest {
			log.Printf("求值表达式 %d: %s", exprIndex+1, expr.name)
			log.Printf("描述: %s", expr.description)
			log.Printf("表达式: %s", expr.expression)

			// 配置求值选项
			options := &EvaluateOptions{
				ObjectGroup:       fmt.Sprintf("frame-%d-inspect", frameIndex),
				GeneratePreview:   true,
				ReturnByValue:     true,
				ThrowOnSideEffect: true,
				Timeout:          2 * time.Second,
			}

			// 在调用帧上求值表达式
			response, err := CDPDebuggerEvaluateOnCallFrame(callFrameId, expr.expression, options)
			if err != nil {
				log.Printf("求值失败: %v", err)
				continue
			}

			// 解析结果
			result, err := ParseEvaluateOnCallFrame(response)
			if err != nil {
				log.Printf("解析结果失败: %v", err)
				continue
			}

			// 显示结果
			displayEvaluationResult(result, exprIndex+1)

			// 短暂延迟
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// 显示求值结果
func displayEvaluationResult(result *EvaluateResult, index int) {
	log.Printf("  [%d] 求值结果:", index)

	if result.ExceptionDetails != nil {
		log.Printf("    ❌ 异常: %v", result.ExceptionDetails)
		return
	}

	if result.Result != nil {
		log.Printf("    ✅ 结果类型: %s", result.Result.Type)
		if result.Result.Value != nil {
			log.Printf("      值: %v", result.Result.Value)
		}
		if result.Result.Description != "" {
			log.Printf("      描述: %s", result.Result.Description)
		}
		if result.Result.Preview != nil {
			log.Printf("      预览: %v", result.Result.Preview)
		}
	}
}

// 示例2: 调试条件表达式测试
func exampleDebugConditionExpressionTest() {
	// === 应用场景描述 ===
	// 场景: 调试条件表达式测试
	// 用途: 测试断点条件表达式的有效性
	// 优势: 确保断点条件在特定上下文中正确工作
	// 典型工作流: 设置条件断点 -> 测试表达式 -> 验证结果 -> 调整条件

	log.Println("调试条件表达式测试示例...")

	// 定义要测试的条件表达式
	conditionTests := []struct {
		name        string
		expression  string
		expected    interface{}
		description string
		testData    map[string]interface{}
	}{
		{
			name:        "简单布尔条件",
			expression:  "x > 10",
			expected:    true,
			description: "测试大于比较",
			testData:    map[string]interface{}{"x": 15},
		},
		{
			name:        "字符串比较",
			expression:  "name === 'test'",
			expected:    true,
			description: "测试字符串相等",
			testData:    map[string]interface{}{"name": "test"},
		},
		{
			name:        "数组长度检查",
			expression:  "items.length > 0",
			expected:    true,
			description: "测试数组非空",
			testData:    map[string]interface{}{"items": []interface{}{1, 2, 3}},
		},
		{
			name:        "对象属性存在性",
			expression:  "'property' in obj && obj.property !== undefined",
			expected:    true,
			description: "测试对象属性存在且非undefined",
			testData:    map[string]interface{}{"obj": map[string]interface{}{"property": "value"}},
		},
		{
			name:        "复杂逻辑表达式",
			expression:  "(a && b) || (!c && d)",
			expected:    true,
			description: "测试复杂逻辑组合",
			testData:    map[string]interface{}{"a": true, "b": false, "c": true, "d": true},
		},
		{
			name:        "函数调用结果",
			expression:  "isValid(input) && !hasError()",
			expected:    true,
			description: "测试函数调用和否定",
			testData:    map[string]interface{}{"input": "valid", "isValid": true, "hasError": false},
		},
		{
			name:        "类型检查组合",
			expression:  "typeof value === 'string' && value.length > 0",
			expected:    true,
			description: "测试类型和属性组合",
			testData:    map[string]interface{}{"value": "hello"},
		},
		{
			name:        "正则表达式匹配",
			expression:  "/^test/.test(str)",
			expected:    true,
			description: "测试正则表达式匹配",
			testData:    map[string]interface{}{"str": "test123"},
		},
	}

	// 模拟调用帧ID
	callFrameId := "debug-frame-001"

	// 测试所有条件表达式
	passedTests := 0
	failedTests := 0
	errorTests := 0

	for testIndex, test := range conditionTests {
		log.Printf("\n测试 %d/%d: %s", testIndex+1, len(conditionTests), test.name)
		log.Printf("描述: %s", test.description)
		log.Printf("表达式: %s", test.expression)
		log.Printf("测试数据: %v", test.testData)
		log.Printf("预期结果: %v", test.expected)

		// 准备测试环境
		// 这里可以设置测试数据到调用帧上下文

		// 创建包含测试数据的表达式
		setupCode := ""
		for key, value := range test.testData {
			setupCode += fmt.Sprintf("var %s = %v; ", key, value)
		}

		fullExpression := setupCode + "(" + test.expression + ")"

		// 配置求值选项
		options := &EvaluateOptions{
			ObjectGroup:       "condition-test",
			GeneratePreview:   false,
			ReturnByValue:     true,
			Silent:           true, // 静默模式，不报告异常
			Timeout:          1 * time.Second,
		}

		// 求值条件表达式
		response, err := CDPDebuggerEvaluateOnCallFrame(callFrameId, fullExpression, options)
		if err != nil {
			log.Printf("❌ 求值失败: %v", err)
			errorTests++
			continue
		}

		// 解析结果
		result, err := ParseEvaluateOnCallFrame(response)
		if err != nil {
			log.Printf("❌ 解析结果失败: %v", err)
			errorTests++
			continue
		}

		// 检查异常
		if result.ExceptionDetails != nil {
			log.Printf("❌ 表达式执行异常: %v", result.ExceptionDetails)
			failedTests++
			continue
		}

		// 检查结果
		if result.Result != nil && result.Result.Value == test.expected {
			log.Printf("✅ 测试通过: 结果符合预期 (%v)", result.Result.Value)
			passedTests++
		} else {
			actualValue := "无结果"
			if result.Result != nil && result.Result.Value != nil {
				actualValue = fmt.Sprintf("%v", result.Result.Value)
			}
			log.Printf("❌ 测试失败: 预期 %v, 实际 %s", test.expected, actualValue)
			failedTests++
		}

		// 短暂延迟
		time.Sleep(50 * time.Millisecond)
	}

	// 生成测试报告
	log.Println("\n=== 条件表达式测试报告 ===")
	log.Printf("总测试数: %d", len(conditionTests))
	log.Printf("通过: %d", passedTests)
	log.Printf("失败: %d", failedTests)
	log.Printf("错误: %d", errorTests)

	successRate := float64(passedTests) / float64(len(conditionTests)) * 100
	log.Printf("成功率: %.1f%%", successRate)

	if successRate >= 90 {
		log.Println("测试评级: ✅ 优秀")
	} else if successRate >= 70 {
		log.Println("测试评级: ⚠ 良好")
	} else if successRate >= 50 {
		log.Println("测试评级: ⚠ 一般")
	} else {
		log.Println("测试评级: ❌ 需改进")
	}
}

// 示例3: 运行时变量修改调试
func exampleRuntimeVariableModification() {
	// === 应用场景描述 ===
	// 场景: 运行时变量修改调试
	// 用途: 在调试时修改变量值以测试不同场景
	// 优势: 无需重新启动应用即可测试不同数据状态
	// 典型工作流: 暂停执行 -> 修改变量 -> 继续执行 -> 观察效果

	log.Println("运行时变量修改调试示例...")

	// 模拟调用帧ID
	callFrameId := "modify-frame-001"

	// 定义变量修改场景
	modificationScenarios := []struct {
		name        string
		variable    string
		modifications []struct {
			expression  string
			description string
			expected    interface{}
		}
	}{
		{
			name:     "计数器变量",
			variable: "counter",
			modifications: []struct {
				expression  string
				description string
				expected    interface{}
			}{
				{
					expression:  "counter = 0",
					description: "重置计数器",
					expected:    0,
				},
				{
					expression:  "counter = counter + 1",
					description: "递增计数器",
					expected:    1,
				},
				{
					expression:  "counter *= 2",
					description: "加倍计数器",
					expected:    2,
				},
				{
					expression:  "counter = 100",
					description: "设置最大值",
					expected:    100,
				},
			},
		},
		{
			name:     "状态标志",
			variable: "isEnabled",
			modifications: []struct {
				expression  string
				description string
				expected    interface{}
			}{
				{
					expression:  "isEnabled = true",
					description: "启用标志",
					expected:    true,
				},
				{
					expression:  "isEnabled = false",
					description: "禁用标志",
					expected:    false,
				},
				{
					expression:  "isEnabled = !isEnabled",
					description: "切换标志",
					expected:    true,
				},
				{
					expression:  "isEnabled = null",
					description: "清空标志",
					expected:    nil,
				},
			},
		},
		{
			name:     "数组操作",
			variable: "items",
			modifications: []struct {
				expression  string
				description string
				expected    interface{}
			}{
				{
					expression:  "items = []",
					description: "清空数组",
					expected:    []interface{}{},
				},
				{
					expression:  "items.push('first')",
					description: "添加元素",
					expected:    1,
				},
				{
					expression:  "items.push('second', 'third')",
					description: "添加多个元素",
					expected:    3,
				},
				{
					expression:  "items.pop()",
					description: "移除最后一个元素",
					expected:    "third",
				},
				{
					expression:  "items.length",
					description: "检查数组长度",
					expected:    2,
				},
			},
		},
		{
			name:     "对象属性",
			variable: "user",
			modifications: []struct {
				expression  string
				description string
				expected    interface{}
			}{
				{
					expression:  "user = {name: 'John', age: 30}",
					description: "创建用户对象",
					expected:    map[string]interface{}{"name": "John", "age": float64(30)},
				},
				{
					expression:  "user.age = 31",
					description: "修改年龄",
					expected:    float64(31),
				},
				{
					expression:  "user.active = true",
					description: "添加新属性",
					expected:    true,
				},
				{
					expression:  "delete user.active",
					description: "删除属性",
					expected:    true,
				},
				{
					expression:  "JSON.stringify(user)",
					description: "查看完整对象",
					expected:    `{"name":"John","age":31}`,
				},
			},
		},
	}

	// 执行变量修改场景
	for scenarioIndex, scenario := range modificationScenarios {
		log.Printf("\n=== 变量修改场景 %d/%d: %s ===",
			scenarioIndex+1, len(modificationScenarios), scenario.name)

		// 初始化变量
		log.Printf("初始化变量: %s", scenario.variable)

		// 执行每个修改
		for modifyIndex, modify := range scenario.modifications {
			log.Printf("修改 %d/%d: %s", modifyIndex+1, len(scenario.modifications), modify.description)
			log.Printf("表达式: %s", modify.expression)

			// 配置求值选项
			options := &EvaluateOptions{
				ObjectGroup:       "modify-group",
				GeneratePreview:   true,
				ReturnByValue:     true,
				Silent:           false,
				ThrowOnSideEffect: false, // 允许副作用
				Timeout:          1 * time.Second,
			}

			// 执行修改表达式
			response, err := CDPDebuggerEvaluateOnCallFrame(callFrameId, modify.expression, options)
			if err != nil {
				log.Printf("❌ 修改失败: %v", err)
				continue
			}

			// 解析结果
			result, err := ParseEvaluateOnCallFrame(response)
			if err != nil {
				log.Printf("❌ 解析结果失败: %v", err)
				continue
			}

			// 检查异常
			if result.ExceptionDetails != nil {
				log.Printf("⚠ 修改有异常: %v", result.ExceptionDetails)
			}

			// 显示修改结果
			if result.Result != nil {
				log.Printf("✅ 修改结果: %v", result.Result.Value)

				// 验证结果
				if result.Result.Value == modify.expected {
					log.Printf("   结果符合预期")
				} else {
					log.Printf("⚠ 结果与预期不符: 预期 %v", modify.expected)
				}
			}

			// 检查变量当前值
			checkOptions := &EvaluateOptions{
				ObjectGroup:   "check-group",
				ReturnByValue: true,
				Silent:       true,
				Timeout:      500 * time.Millisecond,
			}

			checkResponse, err := CDPDebuggerEvaluateOnCallFrame(callFrameId, scenario.variable, checkOptions)
			if err == nil {
				checkResult, err := ParseEvaluateOnCallFrame(checkResponse)
				if err == nil && checkResult.Result != nil {
					log.Printf("   当前 %s = %v", scenario.variable, checkResult.Result.Value)
				}
			}

			// 短暂延迟
			time.Sleep(100 * time.Millisecond)
		}
	}

	log.Println("\n运行时变量修改调试完成")
}

*/

// -----------------------------------------------  Debugger.getPossibleBreakpoints  -----------------------------------------------
// === 应用场景 ===
// 1. 断点位置分析: 分析代码中所有可以设置断点的位置
// 2. 函数调试辅助: 查找函数内部的所有可断点位置
// 3. 代码覆盖率分析: 确定代码中所有可能的执行点
// 4. 调试器优化: 帮助调试器更智能地设置断点
// 5. 自动断点设置: 自动在特定位置设置断点
// 6. 代码分析工具: 分析代码的结构和执行流程

// CDPDebuggerGetPossibleBreakpoints 获取可能的断点位置
// 参数说明:
//   - startScriptId: 起始位置脚本ID
//   - startLine: 起始行号（0-based）
//   - startColumn: 起始列号（0-based）
//   - endScriptId: 结束位置脚本ID（可选）
//   - endLine: 结束行号（可选）
//   - endColumn: 结束列号（可选）
//   - restrictToFunction: 是否限制在同一函数内
func CDPDebuggerGetPossibleBreakpoints(startScriptId string, startLine, startColumn int,
	endScriptId string, endLine, endColumn int, restrictToFunction bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if startScriptId == "" {
		return "", fmt.Errorf("起始脚本ID不能为空")
	}
	if startLine < 0 {
		return "", fmt.Errorf("起始行号必须是非负整数")
	}
	if startColumn < 0 {
		return "", fmt.Errorf("起始列号必须是非负整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建起始位置
	startLocation := map[string]interface{}{
		"scriptId":     startScriptId,
		"lineNumber":   startLine,
		"columnNumber": startColumn,
	}

	// 构建参数对象
	params := map[string]interface{}{
		"start":              startLocation,
		"restrictToFunction": restrictToFunction,
	}

	// 如果有结束位置，添加到参数
	if endScriptId != "" && endLine >= 0 && endColumn >= 0 {
		endLocation := map[string]interface{}{
			"scriptId":     endScriptId,
			"lineNumber":   endLine,
			"columnNumber": endColumn,
		}
		params["end"] = endLocation
	}

	// 序列化参数
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.getPossibleBreakpoints",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getPossibleBreakpoints 请求失败: %w", err)
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
			return "", fmt.Errorf("getPossibleBreakpoints 请求超时")
		}
	}
}

// GetPossibleBreakpointsResult 获取可能断点位置结果
type GetPossibleBreakpointsResult struct {
	Locations []BreakLocation `json:"locations"` // 可能的断点位置列表
}

// BreakLocation 断点位置
type BreakLocation struct {
	ScriptID     string `json:"scriptId"`       // 脚本ID
	LineNumber   int    `json:"lineNumber"`     // 行号（0-based）
	ColumnNumber int    `json:"columnNumber"`   // 列号（0-based）
	Type         string `json:"type,omitempty"` // 类型: "debuggerStatement", "call", "return"
}

// ParseGetPossibleBreakpoints 解析获取可能断点位置响应
func ParseGetPossibleBreakpoints(response string) (*GetPossibleBreakpointsResult, error) {
	var data struct {
		Result *GetPossibleBreakpointsResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

/*

// 示例1: 函数内断点位置分析
func exampleFunctionBreakpointAnalysis() {
	// === 应用场景描述 ===
	// 场景: 函数内断点位置分析
	// 用途: 分析函数内部所有可以设置断点的位置
	// 优势: 帮助开发者了解函数的执行路径和调试点
	// 典型工作流: 选择函数 -> 获取可断点位置 -> 分析分布 -> 设置断点

	log.Println("函数内断点位置分析示例...")

	// 模拟的函数信息
	functionSamples := []struct {
		name        string
		scriptId    string
		startLine   int
		startColumn int
		endLine     int
		endColumn   int
		description string
	}{
		{
			name:        "calculateSum",
			scriptId:    "app-math.js",
			startLine:   10,
			startColumn: 0,
			endLine:     25,
			endColumn:   0,
			description: "计算数组和的简单函数",
		},
		{
			name:        "validateInput",
			scriptId:    "form-utils.js",
			startLine:   45,
			startColumn: 4,
			endLine:     65,
			endColumn:   2,
			description: "输入验证函数，包含多个条件分支",
		},
		{
			name:        "processData",
			scriptId:    "data-processor.js",
			startLine:   88,
			startColumn: 8,
			endLine:     120,
			endColumn:   0,
			description: "复杂数据处理函数，包含循环和错误处理",
		},
		{
			name:        "handleResponse",
			scriptId:    "api-handler.js",
			startLine:   33,
			startColumn: 0,
			endLine:     50,
			endColumn:   15,
			description: "API响应处理函数，包含异步操作",
		},
		{
			name:        "formatDate",
			scriptId:    "date-utils.js",
			startLine:   12,
			startColumn: 4,
			endLine:     30,
			endColumn:   8,
			description: "日期格式化工具函数",
		},
	}

	for i, function := range functionSamples {
		log.Printf("\n=== 函数分析 %d/%d: %s ===", i+1, len(functionSamples), function.name)
		log.Printf("描述: %s", function.description)
		log.Printf("位置: %s:%d:%d - %d:%d",
			function.scriptId, function.startLine, function.startColumn,
			function.endLine, function.endColumn)

		// 获取函数内可能的断点位置
		response, err := CDPDebuggerGetPossibleBreakpoints(
			function.scriptId, function.startLine, function.startColumn,
			function.scriptId, function.endLine, function.endColumn,
			true, // 限制在同一函数内
		)

		if err != nil {
			log.Printf("获取断点位置失败: %v", err)
			continue
		}

		// 解析结果
		result, err := ParseGetPossibleBreakpoints(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 分析断点位置
		analyzeBreakpointLocations(result, function.name)
	}
}

// 分析断点位置
func analyzeBreakpointLocations(result *GetPossibleBreakpointsResult, functionName string) {
	if len(result.Locations) == 0 {
		log.Printf("  未找到可断点位置")
		return
	}

	log.Printf("  找到 %d 个可断点位置:", len(result.Locations))

	// 按类型统计
	typeStats := make(map[string]int)
	lineStats := make(map[int]int)

	for _, location := range result.Locations {
		typeStats[location.Type]++
		lineStats[location.LineNumber]++
	}

	log.Printf("  类型分布:")
	for t, count := range typeStats {
		typeName := t
		if typeName == "" {
			typeName = "常规"
		}
		log.Printf("    %s: %d 个", typeName, count)
	}

	// 行数分布
	log.Printf("  行数分布:")
	lines := make([]int, 0, len(lineStats))
	for line := range lineStats {
		lines = append(lines, line)
	}
	sort.Ints(lines)

	// 显示密度
	minLine := lines[0]
	maxLine := lines[len(lines)-1]
	lineRange := maxLine - minLine + 1

	log.Printf("    行号范围: %d - %d (%d 行)", minLine, maxLine, lineRange)
	log.Printf("    有断点的行: %d 行", len(lines))

	// 计算密度
	density := float64(len(lines)) / float64(lineRange) * 100
	log.Printf("    行覆盖密度: %.1f%%", density)

	// 显示热点行（有多个断点的行）
	hotLines := make([]int, 0)
	for line, count := range lineStats {
		if count > 1 {
			hotLines = append(hotLines, line)
		}
	}

	if len(hotLines) > 0 {
		sort.Ints(hotLines)
		log.Printf("    热点行 (多个断点): %v", hotLines)
	}

	// 显示前5个断点位置
	log.Printf("    前5个断点位置:")
	for i := 0; i < 5 && i < len(result.Locations); i++ {
		loc := result.Locations[i]
		log.Printf("      [%d] 行 %d, 列 %d",
			loc.LineNumber, loc.ColumnNumber)
	}
}

// 示例2: 代码覆盖率断点分析
func exampleCodeCoverageBreakpointAnalysis() {
	// === 应用场景描述 ===
	// 场景: 代码覆盖率断点分析
	// 用途: 分析整个脚本中所有可能的执行点
	// 优势: 帮助进行代码覆盖率测试和调试
	// 典型工作流: 选择脚本 -> 分析所有可断点位置 -> 设计测试用例 -> 评估覆盖率

	log.Println("代码覆盖率断点分析示例...")

	// 模拟的脚本信息
	scriptSamples := []struct {
		name        string
		scriptId    string
		totalLines  int
		description string
	}{
		{
			name:        "用户认证模块",
			scriptId:    "auth-module.js",
			totalLines:  150,
			description: "用户登录、注册、权限验证",
		},
		{
			name:        "数据验证器",
			scriptId:    "data-validator.js",
			totalLines:  120,
			description: "各种数据格式验证",
		},
		{
			name:        "API客户端",
			scriptId:    "api-client.js",
			totalLines:  200,
			description: "HTTP请求封装和错误处理",
		},
		{
			name:        "工具函数库",
			scriptId:    "utility-library.js",
			totalLines:  80,
			description: "常用工具函数集合",
		},
		{
			name:        "状态管理器",
			scriptId:    "state-manager.js",
			totalLines:  180,
			description: "应用状态管理和响应式更新",
		},
	}

	// 分析每个脚本
	for i, script := range scriptSamples {
		log.Printf("\n=== 脚本分析 %d/%d: %s ===", i+1, len(scriptSamples), script.name)
		log.Printf("描述: %s", script.description)
		log.Printf("总行数: %d 行", script.totalLines)

		// 分析整个脚本的断点位置
		startLine := 0
		startColumn := 0

		log.Printf("分析脚本断点位置...")

		response, err := CDPDebuggerGetPossibleBreakpoints(
			script.scriptId, startLine, startColumn,
			"", -1, -1, // 不指定结束位置，分析整个脚本
			false, // 不限制在函数内
		)

		if err != nil {
			log.Printf("分析失败: %v", err)
			continue
		}

		// 解析结果
		result, err := ParseGetPossibleBreakpoints(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 生成覆盖率分析报告
		generateCoverageAnalysisReport(result, script.name, script.totalLines)
	}
}

func generateCoverageAnalysisReport(result *GetPossibleBreakpointsResult, scriptName string, totalLines int) {
	if len(result.Locations) == 0 {
		log.Printf("  未找到可断点位置")
		return
	}

	totalBreakpoints := len(result.Locations)

	// 收集唯一的行号
	uniqueLines := make(map[int]bool)
	for _, location := range result.Locations {
		uniqueLines[location.LineNumber] = true
	}

	// 计算覆盖率
	coveredLines := len(uniqueLines)
	coverageRate := float64(coveredLines) / float64(totalLines) * 100

	log.Printf("  断点分析结果:")
	log.Printf("    总可断点位置: %d 个", totalBreakpoints)
	log.Printf("    有断点的行数: %d 行", coveredLines)
	log.Printf("    总代码行数: %d 行", totalLines)
	log.Printf("    理论覆盖率: %.1f%%", coverageRate)

	// 断点密度分析
	density := float64(totalBreakpoints) / float64(totalLines)
	log.Printf("    断点密度: %.2f 个/行", density)

	// 类型分析
	typeAnalysis := make(map[string]int)
	for _, location := range result.Locations {
		typeName := location.Type
		if typeName == "" {
			typeName = "regular"
		}
		typeAnalysis[typeName]++
	}

	log.Printf("    断点类型分布:")
	for t, count := range typeAnalysis {
		percentage := float64(count) / float64(totalBreakpoints) * 100
		typeDesc := t
		switch t {
		case "debuggerStatement":
			typeDesc = "debugger语句"
		case "call":
			typeDesc = "函数调用"
		case "return":
			typeDesc = "返回语句"
		case "regular":
			typeDesc = "常规语句"
		}
		log.Printf("      %s: %d 个 (%.1f%%)", typeDesc, count, percentage)
	}

	// 行号范围分析
	lines := make([]int, 0, len(uniqueLines))
	for line := range uniqueLines {
		lines = append(lines, line)
	}
	sort.Ints(lines)

	if len(lines) > 0 {
		firstLine := lines[0]
		lastLine := lines[len(lines)-1]
		actualRange := lastLine - firstLine + 1

		log.Printf("    行号范围: %d - %d", firstLine, lastLine)
		log.Printf("    实际代码范围: %d 行", actualRange)

		// 连续性分析
		continuousBlocks := 0
		currentBlock := 0

		for i := 0; i < len(lines); i++ {
			if i == 0 || lines[i] != lines[i-1]+1 {
				continuousBlocks++
				currentBlock++
			}
		}

		log.Printf("    连续性分析:")
		log.Printf("      连续代码块: %d 个", continuousBlocks)

		// 计算平均块大小
		if continuousBlocks > 0 {
			avgBlockSize := float64(coveredLines) / float64(continuousBlocks)
			log.Printf("      平均块大小: %.1f 行", avgBlockSize)
		}
	}

	// 覆盖率评级
	log.Printf("    覆盖率评级:")
	if coverageRate >= 90 {
		log.Printf("      ✅ 优秀: 代码可执行点覆盖很高")
	} else if coverageRate >= 70 {
		log.Printf("      ⚠ 良好: 代码可执行点覆盖适中")
	} else if coverageRate >= 50 {
		log.Printf("      ⚠ 一般: 代码可执行点覆盖较低")
	} else {
		log.Printf("      ❌ 差: 代码可执行点覆盖不足")
		log.Printf("      建议: 代码可能存在大量不可执行代码或注释")
	}
}

// 示例3: 智能断点推荐系统
func exampleSmartBreakpointRecommender() {
	// === 应用场景描述 ===
	// 场景: 智能断点推荐系统
	// 用途: 根据代码结构和执行模式推荐最佳断点位置
	// 优势: 提高调试效率，帮助新手快速找到关键调试点
	// 典型工作流: 分析代码 -> 识别关键位置 -> 推荐断点 -> 自动设置

	log.Println("智能断点推荐系统示例...")

	// 模拟的代码文件
	codeFiles := []struct {
		name        string
		scriptId    string
		codeType    string
		complexity  string
	}{
		{
			name:       "错误处理函数",
			scriptId:  "error-handler.js",
			codeType:  "error",
			complexity: "medium",
		},
		{
			name:       "数据转换器",
			scriptId:  "data-transformer.js",
			codeType:  "data",
			complexity: "high",
		},
		{
			name:       "UI渲染器",
			scriptId:  "ui-renderer.js",
			codeType:  "ui",
			complexity: "high",
		},
		{
			name:       "配置加载器",
			scriptId:  "config-loader.js",
			codeType:  "config",
			complexity: "low",
		},
		{
			name:       "缓存管理器",
			scriptId:  "cache-manager.js",
			codeType:  "cache",
			complexity: "medium",
		},
	}

	for i, file := range codeFiles {
		log.Printf("\n=== 文件分析 %d/%d: %s ===", i+1, len(codeFiles), file.name)
		log.Printf("类型: %s, 复杂度: %s", file.codeType, file.complexity)

		// 分析文件的所有可能断点
		response, err := CDPDebuggerGetPossibleBreakpoints(
			file.scriptId, 0, 0,
			"", -1, -1,
			false,
		)

		if err != nil {
			log.Printf("分析失败: %v", err)
			continue
		}

		result, err := ParseGetPossibleBreakpoints(response)
		if err != nil {
			log.Printf("解析失败: %v", err)
			continue
		}

		// 根据代码类型推荐断点
		recommendBreakpoints(result, file.codeType, file.complexity)
	}
}

func recommendBreakpoints(result *GetPossibleBreakpointsResult, codeType, complexity string) {
	if len(result.Locations) == 0 {
		log.Printf("  没有可推荐的断点位置")
		return
	}

	log.Printf("  断点位置总数: %d", len(result.Locations))

	// 根据代码类型和复杂度选择推荐策略
	maxRecommendations := 5
	if complexity == "high" {
		maxRecommendations = 8
	} else if complexity == "low" {
		maxRecommendations = 3
	}

	// 分类收集位置
	locationsByType := make(map[string][]BreakLocation)
	for _, loc := range result.Locations {
		locType := loc.Type
		if locType == "" {
			locType = "regular"
		}
		locationsByType[locType] = append(locationsByType[locType], loc)
	}

	// 根据代码类型推荐不同类型的断点
	log.Printf("  推荐断点位置:")

	switch codeType {
	case "error":
		// 错误处理代码：推荐函数调用和返回位置
		recommendErrorBreakpoints(locationsByType, maxRecommendations)
	case "data":
		// 数据处理代码：推荐关键计算和转换位置
		recommendDataBreakpoints(locationsByType, maxRecommendations)
	case "ui":
		// UI代码：推荐状态变更和渲染位置
		recommendUIBreakpoints(locationsByType, maxRecommendations)
	case "config":
		// 配置代码：推荐加载和验证位置
		recommendConfigBreakpoints(locationsByType, maxRecommendations)
	case "cache":
		// 缓存代码：推荐读写和失效位置
		recommendCacheBreakpoints(locationsByType, maxRecommendations)
	default:
		// 默认推荐：混合类型
		recommendDefaultBreakpoints(locationsByType, maxRecommendations)
	}
}

func recommendErrorBreakpoints(locationsByType map[string][]BreakLocation, maxCount int) {
	recommendations := make([]BreakLocation, 0)

	// 优先推荐调用位置（函数调用点）
	if calls, ok := locationsByType["call"]; ok && len(calls) > 0 {
		count := min(2, len(calls))
		recommendations = append(recommendations, calls[:count]...)
	}

	// 推荐返回位置（错误返回点）
	if returns, ok := locationsByType["return"]; ok && len(returns) > 0 {
		count := min(1, len(returns))
		recommendations = append(recommendations, returns[:count]...)
	}

	// 补充常规位置
	if regulars, ok := locationsByType["regular"]; ok && len(regulars) > 0 {
		needed := maxCount - len(recommendations)
		if needed > 0 && len(regulars) > 0 {
			count := min(needed, len(regulars))
			// 选择间隔分布的位置
			step := len(regulars) / count
			if step < 1 {
				step = 1
			}

			for i := 0; i < count && i*step < len(regulars); i++ {
				idx := i * step
				recommendations = append(recommendations, regulars[idx])
			}
		}
	}

	displayRecommendations(recommendations, "错误处理")
}

func recommendDataBreakpoints(locationsByType map[string][]BreakLocation, maxCount int) {
	// 数据代码：推荐所有类型的混合位置
	displayRecommendations(selectMixedLocations(locationsByType, maxCount), "数据处理")
}

func recommendUIBreakpoints(locationsByType map[string][]BreakLocation, maxCount int) {
	// UI代码：优先推荐调用和返回位置
	displayRecommendations(selectMixedLocations(locationsByType, maxCount), "UI渲染")
}

func recommendConfigBreakpoints(locationsByType map[string][]BreakLocation, maxCount int) {
	// 配置代码：推荐较少的断点
	displayRecommendations(selectMixedLocations(locationsByType, min(3, maxCount)), "配置管理")
}

func recommendCacheBreakpoints(locationsByType map[string][]BreakLocation, maxCount int) {
	// 缓存代码：推荐调用位置
	displayRecommendations(selectMixedLocations(locationsByType, maxCount), "缓存管理")
}

func recommendDefaultBreakpoints(locationsByType map[string][]BreakLocation, maxCount int) {
	displayRecommendations(selectMixedLocations(locationsByType, maxCount), "通用")
}

func selectMixedLocations(locationsByType map[string][]BreakLocation, maxCount int) []BreakLocation {
	recommendations := make([]BreakLocation, 0)

	// 按类型优先级选择
	typeOrder := []string{"call", "return", "debuggerStatement", "regular"}

	for _, t := range typeOrder {
		if locs, ok := locationsByType[t]; ok && len(locs) > 0 {
			count := min(1, len(locs))
			if len(recommendations)+count <= maxCount {
				recommendations = append(recommendations, locs[:count]...)
			}
		}
	}

	// 如果还不够，从所有位置中补充
	if len(recommendations) < maxCount {
		allLocs := make([]BreakLocation, 0)
		for _, locs := range locationsByType {
			allLocs = append(allLocs, locs...)
		}

		// 按行号排序
		sort.Slice(allLocs, func(i, j int) bool {
			return allLocs[i].LineNumber < allLocs[j].LineNumber
		})

		needed := maxCount - len(recommendations)
		if needed > 0 && len(allLocs) > 0 {
			// 选择均匀分布的位置
			step := len(allLocs) / needed
			if step < 1 {
				step = 1
			}

			for i := 0; i < needed && i*step < len(allLocs); i++ {
				idx := i * step
				recommendations = append(recommendations, allLocs[idx])
			}
		}
	}

	return recommendations
}

func displayRecommendations(locations []BreakLocation, codeType string) {
	if len(locations) == 0 {
		log.Printf("    无推荐断点")
		return
	}

	log.Printf("    %s代码推荐断点:", codeType)
	for i, loc := range locations {
		typeDesc := loc.Type
		if typeDesc == "" {
			typeDesc = "常规"
		}
		log.Printf("    [%d] 行 %d, 列 %d (%s)",
			i+1, loc.LineNumber, loc.ColumnNumber, typeDesc)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

*/

// -----------------------------------------------  Debugger.getScriptSource  -----------------------------------------------
// === 应用场景 ===
// 1. 源码查看器: 实时查看和分析JavaScript脚本源码
// 2. 源码调试工具: 在调试时获取当前执行脚本的源码
// 3. 代码分析: 分析已加载脚本的代码结构和内容
// 4. 源码比较: 比较不同版本的脚本源码差异
// 5. 源码备份: 备份运行时加载的脚本源码
// 6. 代码审计: 审计第三方或动态加载的脚本

// CDPDebuggerGetScriptSource 获取脚本源码
// 参数说明:
//   - scriptId: 脚本ID
func CDPDebuggerGetScriptSource(scriptId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if scriptId == "" {
		return "", fmt.Errorf("脚本ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.getScriptSource",
		"params": {
			"scriptId": "%s"
		}
	}`, reqID, scriptId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getScriptSource 请求失败: %w", err)
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
			return "", fmt.Errorf("getScriptSource 请求超时")
		}
	}
}

// GetScriptSourceResult 获取脚本源码结果
type GetScriptSourceResult struct {
	ScriptSource string `json:"scriptSource"`       // 脚本源码
	Bytecode     string `json:"bytecode,omitempty"` // Wasm字节码（base64编码）
}

// ParseGetScriptSource 解析获取脚本源码响应
func ParseGetScriptSource(response string) (*GetScriptSourceResult, error) {
	var data struct {
		Result *GetScriptSourceResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

/*

// 示例1: 脚本源码查看和分析工具
func exampleScriptSourceViewer() {
	// === 应用场景描述 ===
	// 场景: 脚本源码查看和分析工具
	// 用途: 查看和分析已加载JavaScript脚本的源码
	// 优势: 实时获取源码，支持代码分析和调试
	// 典型工作流: 选择脚本 -> 获取源码 -> 分析结构 -> 调试定位

	log.Println("脚本源码查看和分析工具...")

	// 模拟已加载的脚本
	loadedScripts := []struct {
		scriptId    string
		name        string
		description string
		expectedType string
	}{
		{
			scriptId:    "script-1234567890",
			name:        "主应用脚本",
			description: "应用的主要业务逻辑",
			expectedType: "javascript",
		},
		{
			scriptId:    "script-9876543210",
			name:        "工具库脚本",
			description: "工具函数和工具类",
			expectedType: "javascript",
		},
		{
			scriptId:    "wasm-module-001",
			name:        "WebAssembly模块",
			description: "高性能计算模块",
			expectedType: "wasm",
		},
		{
			scriptId:    "vendor-lib-001",
			name:        "第三方库",
			description: "第三方JavaScript库",
			expectedType: "javascript",
		},
		{
			scriptId:    "dynamic-script-001",
			name:        "动态加载脚本",
			description: "运行时动态加载的脚本",
			expectedType: "javascript",
		},
		{
			scriptId:    "inline-script-001",
			name:        "内联脚本",
			description: "HTML中内嵌的脚本",
			expectedType: "javascript",
		},
	}

	for i, script := range loadedScripts {
		log.Printf("\n=== 脚本分析 %d/%d: %s ===", i+1, len(loadedScripts), script.name)
		log.Printf("描述: %s", script.description)
		log.Printf("脚本ID: %s", script.scriptId)
		log.Printf("预期类型: %s", script.expectedType)

		// 获取脚本源码
		log.Println("正在获取脚本源码...")

		startTime := time.Now()
		response, err := CDPDebuggerGetScriptSource(script.scriptId)
		if err != nil {
			log.Printf("获取源码失败: %v", err)
			continue
		}

		getTime := time.Since(startTime)

		// 解析结果
		result, err := ParseGetScriptSource(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 分析源码
		analyzeScriptSource(result, script.name, script.expectedType, getTime)
	}
}

func analyzeScriptSource(result *GetScriptSourceResult, scriptName, expectedType string, fetchTime time.Duration) {
	log.Printf("获取时间: %v", fetchTime)

	// 检查脚本类型
	if result.Bytecode != "" && expectedType == "wasm" {
		log.Printf("✅ 识别为WebAssembly模块")
		log.Printf("  字节码大小: %d 字节 (base64)", len(result.Bytecode))

		// 如果是Wasm，尝试解码base64获取实际大小
		if decoded, err := base64.StdEncoding.DecodeString(result.Bytecode); err == nil {
			log.Printf("  解码后大小: %d 字节", len(decoded))

			// 简单的Wasm模块分析
			if len(decoded) > 0 {
				log.Printf("  Wasm魔数: %x", decoded[:4])

				// 检查标准Wasm魔数 "\0asm"
				if len(decoded) >= 4 && string(decoded[:4]) == "\x00asm" {
					log.Printf("  ✅ 有效的Wasm模块")
					version := binary.LittleEndian.Uint32(decoded[4:8])
					log.Printf("  Wasm版本: %d", version)
				} else {
					log.Printf("  ⚠ 非标准Wasm格式")
				}
			}
		}
	} else if result.ScriptSource != "" {
		log.Printf("✅ 识别为JavaScript脚本")
		log.Printf("  源码大小: %d 字符", len(result.ScriptSource))

		// 基本源码分析
		lines := strings.Split(result.ScriptSource, "\n")
		log.Printf("  行数: %d 行", len(lines))

		// 计算代码行数（排除空行和注释）
		codeLines := 0
		commentLines := 0
		emptyLines := 0

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				emptyLines++
			} else if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
				commentLines++
			} else {
				codeLines++
			}
		}

		log.Printf("  代码行: %d", codeLines)
		log.Printf("  注释行: %d", commentLines)
		log.Printf("  空行: %d", emptyLines)

		// 计算密度
		totalLines := len(lines)
		if totalLines > 0 {
			codeDensity := float64(codeLines) / float64(totalLines) * 100
			commentDensity := float64(commentLines) / float64(totalLines) * 100

			log.Printf("  代码密度: %.1f%%", codeDensity)
			log.Printf("  注释密度: %.1f%%", commentDensity)
		}

		// 源码预览
		log.Printf("  源码预览 (前10行):")
		for i := 0; i < 10 && i < len(lines); i++ {
			lineNum := i + 1
			linePreview := lines[i]
			if len(linePreview) > 60 {
				linePreview = linePreview[:60] + "..."
			}
			log.Printf("    %4d: %s", lineNum, linePreview)
		}

		// 简单语法分析
		log.Printf("  语法分析:")
		keywords := []string{
			"function", "const", "let", "var", "if", "for", "while",
			"return", "class", "import", "export", "async", "await",
		}

		keywordCounts := make(map[string]int)
		sourceLower := strings.ToLower(result.ScriptSource)

		for _, keyword := range keywords {
			count := strings.Count(sourceLower, keyword+" ")
			if count > 0 {
				keywordCounts[keyword] = count
			}
		}

		if len(keywordCounts) > 0 {
			log.Printf("    关键字统计:")
			for keyword, count := range keywordCounts {
				log.Printf("      %s: %d", keyword, count)
			}
		}
	} else {
		log.Printf("⚠ 脚本为空或无内容")
	}

	// 性能分析
	log.Printf("  获取性能: %v", fetchTime)
	if fetchTime > 1*time.Second {
		log.Printf("  ⚠ 获取时间较长")
	} else if fetchTime > 100*time.Millisecond {
		log.Printf("  ⚠ 获取时间适中")
	} else {
		log.Printf("  ✅ 获取性能良好")
	}
}

// 示例2: 源码比较和差异分析
func exampleSourceCodeComparison() {
	// === 应用场景描述 ===
	// 场景: 源码比较和差异分析
	// 用途: 比较不同版本的脚本源码，分析变化
	// 优势: 帮助理解代码演进，检测意外修改
	// 典型工作流: 获取多个版本 -> 比较差异 -> 分析变化 -> 生成报告

	log.Println("源码比较和差异分析示例...")

	// 模拟不同版本的脚本
	scriptVersions := []struct {
		version     string
		scriptId    string
		description string
		changes     string
	}{
		{
			version:     "v1.0.0",
			scriptId:    "script-v1",
			description: "初始版本",
			changes:     "基础功能实现",
		},
		{
			version:     "v1.1.0",
			scriptId:    "script-v2",
			description: "功能增强版",
			changes:     "添加新功能，优化性能",
		},
		{
			version:     "v1.2.0",
			scriptId:    "script-v3",
			description: "修复版本",
			changes:     "修复bug，改进错误处理",
		},
		{
			version:     "v2.0.0",
			scriptId:    "script-v4",
			description: "重大更新",
			changes:     "重构架构，API变更",
		},
	}

	// 获取所有版本的源码
	var versions []*ScriptVersion

	for _, ver := range scriptVersions {
		log.Printf("获取版本: %s (%s)", ver.version, ver.description)

		response, err := CDPDebuggerGetScriptSource(ver.scriptId)
		if err != nil {
			log.Printf("获取失败: %v", err)
			continue
		}

		result, err := ParseGetScriptSource(response)
		if err != nil {
			log.Printf("解析失败: %v", err)
			continue
		}

		version := &ScriptVersion{
			Version:     ver.version,
			Description: ver.description,
			Changes:     ver.changes,
			Source:      result.ScriptSource,
			Bytecode:    result.Bytecode,
		}

		versions = append(versions, version)
		log.Printf("✅ 获取成功 (%d 字符)", len(version.Source))
	}

	// 比较不同版本
	if len(versions) > 1 {
		log.Println("\n=== 版本比较分析 ===")

		for i := 0; i < len(versions)-1; i++ {
			version1 := versions[i]
			version2 := versions[i+1]

			log.Printf("\n比较: %s → %s", version1.Version, version2.Version)
			log.Printf("描述: %s", version2.Changes)

			compareScriptVersions(version1, version2)
		}
	}
}

// ScriptVersion 脚本版本
type ScriptVersion struct {
	Version     string
	Description string
	Changes     string
	Source      string
	Bytecode    string
}

func compareScriptVersions(v1, v2 *ScriptVersion) {
	// 基本统计比较
	size1 := len(v1.Source)
	size2 := len(v2.Source)
	sizeDiff := size2 - size1

	log.Printf("  大小变化: %d → %d 字符 (%+d)", size1, size2, sizeDiff)

	// 行数比较
	lines1 := strings.Split(v1.Source, "\n")
	lines2 := strings.Split(v2.Source, "\n")

	lineDiff := len(lines2) - len(lines1)
	log.Printf("  行数变化: %d → %d 行 (%+d)", len(lines1), len(lines2), lineDiff)

	// 分析行数变化类型
	if lineDiff > 0 {
		log.Printf("  新增大约 %d 行代码", lineDiff)
	} else if lineDiff < 0 {
		log.Printf("  删除大约 %d 行代码", -lineDiff)
	} else {
		log.Printf("  行数无变化")
	}

	// 简单内容比较
	if v1.Source == v2.Source {
		log.Printf("  ✅ 源码完全一致")
	} else {
		// 使用简单算法比较差异
		diff := calculateSimpleDiff(v1.Source, v2.Source)

		log.Printf("  变化分析:")
		log.Printf("    新增内容: %d 字符", diff.added)
		log.Printf("    删除内容: %d 字符", diff.removed)
		log.Printf("    修改内容: %d 字符", diff.modified)

		// 计算相似度
		similarity := calculateSimilarity(v1.Source, v2.Source)
		log.Printf("    相似度: %.1f%%", similarity*100)

		// 变化程度评估
		if similarity > 0.9 {
			log.Printf("  ✅ 微小变化")
		} else if similarity > 0.7 {
			log.Printf("  ⚠ 中等变化")
		} else if similarity > 0.5 {
			log.Printf("  ⚠ 较大变化")
		} else {
			log.Printf("  🔥 重大变化")
		}
	}
}

// 简单差异统计
type diffStats struct {
	added    int
	removed  int
	modified int
}

func calculateSimpleDiff(source1, source2 string) diffStats {
	// 简单的基于行的差异计算
	lines1 := strings.Split(source1, "\n")
	lines2 := strings.Split(source2, "\n")

	var stats diffStats

	// 创建行到索引的映射
	lineMap1 := make(map[string]int)
	for _, line := range lines1 {
		lineMap1[line]++
	}

	lineMap2 := make(map[string]int)
	for _, line := range lines2 {
		lineMap2[line]++
	}

	// 统计新增的行
	for line, count2 := range lineMap2 {
		count1 := lineMap1[line]
		if count2 > count1 {
			stats.added += (count2 - count1) * len(line)
		}
	}

	// 统计删除的行
	for line, count1 := range lineMap1 {
		count2 := lineMap2[line]
		if count1 > count2 {
			stats.removed += (count1 - count2) * len(line)
		}
	}

	return stats
}

func calculateSimilarity(source1, source2 string) float64 {
	// 简单的基于字符的相似度计算
	if source1 == source2 {
		return 1.0
	}

	if len(source1) == 0 || len(source2) == 0 {
		return 0.0
	}

	// 使用编辑距离的简化版本
	distance := levenshteinDistance(source1, source2)
	maxLen := max(len(source1), len(source2))

	return 1.0 - float64(distance)/float64(maxLen)
}

func levenshteinDistance(a, b string) int {
	// 简化的编辑距离计算
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// 使用较小的窗口进行比较
	window := 100
	if len(a) > window || len(b) > window {
		// 对于大文本，使用采样比较
		sampleA := a
		if len(a) > window {
			sampleA = a[:window]
		}

		sampleB := b
		if len(b) > window {
			sampleB = b[:window]
		}

		return simpleLevenshtein(sampleA, sampleB)
	}

	return simpleLevenshtein(a, b)
}

func simpleLevenshtein(a, b string) int {
	// 简单的Levenshtein距离实现
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
		d[i][0] = i
	}

	for j := range d[0] {
		d[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			d[i][j] = min(
				d[i-1][j]+1,   // 删除
				d[i][j-1]+1,   // 插入
				d[i-1][j-1]+cost, // 替换
			)
		}
	}

	return d[len(a)][len(b)]
}

func min(nums ...int) int {
	minNum := nums[0]
	for _, num := range nums[1:] {
		if num < minNum {
			minNum = num
		}
	}
	return minNum
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 示例3: 源码审计和安全分析
func exampleSourceCodeAudit() {
	// === 应用场景描述 ===
	// 场景: 源码审计和安全分析
	// 用途: 检查脚本源码中的安全问题
	// 优势: 及时发现潜在的安全漏洞
	// 典型工作流: 获取源码 -> 安全分析 -> 检测问题 -> 生成报告

	log.Println("源码审计和安全分析示例...")

	// 模拟需要审计的脚本
	scriptsToAudit := []struct {
		name        string
		scriptId    string
		category    string
		riskLevel   string
	}{
		{
			name:      "用户输入处理器",
			scriptId:  "user-input-handler.js",
			category:  "输入验证",
			riskLevel: "高",
		},
		{
			name:      "API调用封装",
			scriptId:  "api-wrapper.js",
			category:  "网络请求",
			riskLevel: "中",
		},
		{
			name:      "数据存储模块",
			scriptId:  "data-storage.js",
			category:  "数据安全",
			riskLevel: "高",
		},
		{
			name:      "第三方集成",
			scriptId:  "third-party-integration.js",
			category:  "外部依赖",
			riskLevel: "中",
		},
		{
			name:      "认证授权模块",
			scriptId:  "auth-module.js",
			category:  "访问控制",
			riskLevel: "极高",
		},
	}

	for i, script := range scriptsToAudit {
		log.Printf("\n=== 安全审计 %d/%d: %s ===", i+1, len(scriptsToAudit), script.name)
		log.Printf("类别: %s", script.category)
		log.Printf("风险等级: %s", script.riskLevel)

		// 获取源码
		response, err := CDPDebuggerGetScriptSource(script.scriptId)
		if err != nil {
			log.Printf("获取源码失败: %v", err)
			continue
		}

		result, err := ParseGetScriptSource(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 执行安全审计
		performSecurityAudit(result, script.name, script.category, script.riskLevel)
	}
}

func performSecurityAudit(result *GetScriptSourceResult, scriptName, category, riskLevel string) {
	if result.ScriptSource == "" {
		log.Printf("⚠ 脚本无源码，无法进行安全审计")
		return
	}

	log.Printf("源码大小: %d 字符", len(result.ScriptSource))

	// 安全检测规则
	securityRules := []struct {
		name        string
		pattern     string
		severity    string
		description string
	}{
		// XSS相关
		{
			name:        "innerHTML直接赋值",
			pattern:     `\.innerHTML\s*=`,
			severity:    "高",
			description: "直接innerHTML赋值可能导致XSS",
		},
		{
			name:        "eval函数调用",
			pattern:     `eval\(`,
			severity:    "高",
			description: "eval可能执行恶意代码",
		},
		{
			name:        "document.write",
			pattern:     `document\.write\(`,
			severity:    "中",
			description: "document.write可能被滥用",
		},

		// 数据安全
		{
			name:        "localStorage敏感数据",
			pattern:     `localStorage\.(setItem|getItem)\s*\(\s*['"][^'"]*(password|token|secret|key)`,
			severity:    "高",
			description: "敏感数据不应存储在localStorage",
		},

		// API安全
		{
			name:        "fetch无CORS检查",
			pattern:     `fetch\([^)]*\)[^{]*\.then`,
			severity:    "中",
			description: "fetch调用应检查CORS和错误",
		},

		// 代码质量
		{
			name:        "调试代码",
			pattern:     `console\.(log|debug|info|warn|error)`,
			severity:    "低",
			description: "生产环境应移除调试代码",
		},
		{
			name:        "TODO/FIXME注释",
			pattern:     `\/\/\s*(TODO|FIXME|HACK|XXX)`,
			severity:    "低",
			description: "未完成的代码标记",
		},

		// 密码学
		{
			name:        "弱随机数",
			pattern:     `Math\.random\(\)`,
			severity:    "中",
			description: "Math.random()不适合安全用途",
		},
		{
			name:        "自定义加密",
			pattern:     `function\s+\w+\s*\([^)]*\)[^{]*{[^}]*[+/%&|^][^}]*=`,
			severity:    "高",
			description: "自定义加密算法通常不安全",
		},

		// 输入验证
		{
			name:        "缺少输入验证",
			pattern:     `function\s+\w+\s*\((\w+)(,\s*\w+)*\)`,
			severity:    "中",
			description: "函数参数缺少验证检查",
		},
	}

	// 执行检测
	var findings []SecurityFinding

	for _, rule := range securityRules {
		re := regexp.MustCompile(rule.pattern)
		matches := re.FindAllStringIndex(result.ScriptSource, -1)

		if len(matches) > 0 {
			finding := SecurityFinding{
				RuleName:    rule.name,
				Severity:    rule.severity,
				Description: rule.description,
				Count:       len(matches),
				Lines:       make([]int, 0),
			}

			// 记录行号
			lines := strings.Split(result.ScriptSource, "\n")
			lineNum := 1
			for _, line := range lines {
				if re.MatchString(line) {
					finding.Lines = append(finding.Lines, lineNum)
				}
				lineNum++
			}

			findings = append(findings, finding)
		}
	}

	// 显示审计结果
	log.Printf("安全审计结果:")

	if len(findings) == 0 {
		log.Printf("  ✅ 未发现明显的安全问题")
	} else {
		log.Printf("  ⚠ 发现 %d 个安全问题:", len(findings))

		// 按严重程度分组
		severityCounts := make(map[string]int)
		for _, finding := range findings {
			severityCounts[finding.Severity]++
		}

		log.Printf("  严重程度分布:")
		for severity, count := range severityCounts {
			log.Printf("    %s: %d 个", severity, count)
		}

		// 显示严重问题
		log.Printf("  详细问题:")
		for _, finding := range findings {
			if finding.Severity == "高" || finding.Severity == "极高" {
				log.Printf("    ❌ [%s] %s", finding.Severity, finding.RuleName)
				log.Printf("        描述: %s", finding.Description)
				log.Printf("        次数: %d 次", finding.Count)
				if len(finding.Lines) > 0 {
					lineStr := ""
					for i, line := range finding.Lines {
						if i < 3 { // 只显示前3个行号
							if i > 0 {
								lineStr += ", "
							}
							lineStr += fmt.Sprintf("%d", line)
						}
					}
					if len(finding.Lines) > 3 {
						lineStr += fmt.Sprintf(", ... 等 %d 处", len(finding.Lines))
					}
					log.Printf("        行号: %s", lineStr)
				}
			}
		}

		// 安全评级
		if severityCounts["极高"] > 0 {
			log.Printf("  🔥 安全评级: 极差")
		} else if severityCounts["高"] > 2 {
			log.Printf("  ❌ 安全评级: 差")
		} else if severityCounts["高"] > 0 || severityCounts["中"] > 3 {
			log.Printf("  ⚠ 安全评级: 中")
		} else {
			log.Printf("  ⚠ 安全评级: 良")
		}
	}
}

// SecurityFinding 安全发现
type SecurityFinding struct {
	RuleName    string
	Severity    string
	Description string
	Count       int
	Lines       []int
}

*/

// -----------------------------------------------  Debugger.pause  -----------------------------------------------
// === 应用场景 ===
// 1. 紧急暂停: 在代码执行过程中紧急暂停调试
// 2. 断点替代: 手动触发暂停，替代断点功能
// 3. 执行控制: 控制代码执行流程
// 4. 状态检查: 暂停执行检查当前状态
// 5. 异步调试: 暂停异步代码执行
// 6. 动态分析: 在运行时动态暂停分析

// CDPDebuggerPause 暂停调试器
func CDPDebuggerPause() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.pause"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 pause 请求失败: %w", err)
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
			return "", fmt.Errorf("pause 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.restartFrame  -----------------------------------------------
// === 应用场景 ===
// 1. 函数重试: 重新执行特定函数调用
// 2. 调试循环: 重新执行特定迭代进行调试
// 3. 错误重现: 重现特定错误场景
// 4. 状态重置: 重置函数调用状态
// 5. 性能测试: 重新执行函数测试性能
// 6. 条件测试: 在不同条件下重新执行函数

// CDPDebuggerRestartFrame 重新启动调用帧
// 参数说明:
//   - callFrameId: 调用帧ID
func CDPDebuggerRestartFrame(callFrameId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if callFrameId == "" {
		return "", fmt.Errorf("调用帧ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.restartFrame",
		"params": {
			"callFrameId": "%s",
			"mode": "StepInto"
		}
	}`, reqID, callFrameId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 restartFrame 请求失败: %w", err)
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
			return "", fmt.Errorf("restartFrame 请求超时")
		}
	}
}

// RestartFrameResult 重启帧结果
type RestartFrameResult struct {
	CallFrames        []CallFrame   `json:"callFrames"`                  // 新的调用栈（已弃用，总是为空）
	AsyncStackTrace   *StackTrace   `json:"asyncStackTrace,omitempty"`   // 异步调用栈
	AsyncStackTraceID *StackTraceID `json:"asyncStackTraceId,omitempty"` // 异步调用栈ID
}

// StackTraceID 调用栈跟踪ID
type StackTraceID struct {
	ID         string `json:"id"`                   // 异步调用栈的唯一标识符
	DebuggerID string `json:"debuggerId,omitempty"` // 关联的调试器ID
}

// Scope 作用域
type Scope struct {
	Type   string      `json:"type"`            // 作用域类型
	Object interface{} `json:"object"`          // 作用域对象
	Start  *Location   `json:"start,omitempty"` // 起始位置
	End    *Location   `json:"end,omitempty"`   // 结束位置
}

// ParseRestartFrame 解析重启帧响应
func ParseRestartFrame(response string) (*RestartFrameResult, error) {
	var data struct {
		Result *RestartFrameResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

/*

// 示例1: 函数重试调试
func exampleFunctionRetryDebugging() {
	// === 应用场景描述 ===
	// 场景: 函数重试调试
	// 用途: 重新执行函数以调试特定执行路径
	// 优势: 无需重新启动程序即可重试函数执行
	// 典型工作流: 暂停在函数内部 -> 重启帧 -> 重新执行 -> 观察结果

	log.Println("函数重试调试示例...")

	// 模拟可重试的函数调用
	retryFunctions := []struct {
		name        string
		callFrameId string
		description string
		retryCount  int
		testCases   []struct {
			input  string
			output string
		}
	}{
		{
			name:        "数据处理函数",
			callFrameId: "process-data-frame-001",
			description: "处理用户输入数据",
			retryCount:  3,
			testCases: []struct {
				input  string
				output string
			}{
				{"正常输入", "处理成功"},
				{"边界值输入", "处理成功"},
				{"异常输入", "处理失败"},
			},
		},
		{
			name:        "验证函数",
			callFrameId: "validation-frame-001",
			description: "验证表单数据",
			retryCount:  4,
			testCases: []struct {
				input  string
				output string
			}{
				{"有效数据", "验证通过"},
				{"无效格式", "验证失败"},
				{"空数据", "验证失败"},
				{"特殊字符", "验证通过"},
			},
		},
		{
			name:        "计算函数",
			callFrameId: "calculation-frame-001",
			description: "执行复杂计算",
			retryCount:  2,
			testCases: []struct {
				input  string
				output string
			}{
				{"标准输入", "计算结果"},
				{"大数值输入", "溢出处理"},
			},
		},
		{
			name:        "网络请求函数",
			callFrameId: "network-frame-001",
			description: "发送网络请求",
			retryCount:  3,
			testCases: []struct {
				input  string
				output string
			}{
				{"成功请求", "响应成功"},
				{"超时请求", "响应超时"},
				{"错误请求", "响应错误"},
			},
		},
		{
			name:        "状态更新函数",
			callFrameId: "state-update-frame-001",
			description: "更新应用状态",
			retryCount:  3,
			testCases: []struct {
				input  string
				output string
			}{
				{"初始状态", "状态更新"},
				{"中间状态", "状态更新"},
				{"最终状态", "状态完成"},
			},
		},
	}

	// 测试每个函数的重试
	for i, function := range retryFunctions {
		log.Printf("\n=== 函数重试 %d/%d: %s ===", i+1, len(retryFunctions), function.name)
		log.Printf("描述: %s", function.description)
		log.Printf("重试次数: %d", function.retryCount)

		// 执行多次重试
		for retryIndex := 0; retryIndex < function.retryCount; retryIndex++ {
			log.Printf("重试 %d/%d:", retryIndex+1, function.retryCount)

			// 准备测试用例
			testCaseIndex := retryIndex % len(function.testCases)
			testCase := function.testCases[testCaseIndex]

			log.Printf("  测试用例: 输入='%s', 预期输出='%s'",
				testCase.input, testCase.output)

			// 模拟函数执行前的状态
			log.Printf("  模拟函数执行前状态...")
			preExecutionState := simulateFunctionState(function.name, testCase.input)
			log.Printf("  执行前状态: %v", preExecutionState)

			// 重启函数调用帧
			log.Printf("  重启调用帧...")
			response, err := CDPDebuggerRestartFrame(function.callFrameId)
			if err != nil {
				log.Printf("  重启失败: %v", err)
				continue
			}

			log.Printf("  重启成功")

			// 解析重启结果
			result, err := ParseRestartFrame(response)
			if err != nil {
				log.Printf("  解析结果失败: %v", err)
			} else {
				// 分析异步调用栈
				if result.AsyncStackTrace != nil {
					log.Printf("  异步调用栈信息:")
					if result.AsyncStackTrace.Description != "" {
						log.Printf("    描述: %s", result.AsyncStackTrace.Description)
					}
					if len(result.AsyncStackTrace.CallFrames) > 0 {
						log.Printf("    调用帧数: %d", len(result.AsyncStackTrace.CallFrames))
					}
				}
			}

			// 模拟函数执行后的状态
			log.Printf("  模拟函数执行后状态...")
			postExecutionState := simulateFunctionState(function.name, testCase.output)
			log.Printf("  执行后状态: %v", postExecutionState)

			// 比较状态变化
			log.Printf("  状态变化分析:")
			analyzeStateChange(preExecutionState, postExecutionState, testCase.output)

			// 短暂延迟
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func simulateFunctionState(functionName, input string) map[string]interface{} {
	// 模拟函数状态
	state := make(map[string]interface{})

	switch functionName {
	case "数据处理函数":
		state["inputData"] = input
		state["processingStage"] = "数据清洗"
		state["errorCount"] = 0
		if input == "异常输入" {
			state["errorCount"] = 1
			state["hasError"] = true
		}
	case "验证函数":
		state["validationInput"] = input
		state["isValid"] = input == "有效数据" || input == "特殊字符"
		state["validationRules"] = []string{"格式检查", "长度检查", "内容检查"}
	case "计算函数":
		state["calculationInput"] = input
		state["result"] = "计算结果"
		if input == "大数值输入" {
			state["hasOverflow"] = true
		}
	case "网络请求函数":
		state["requestInput"] = input
		state["requestStatus"] = "进行中"
		if input == "成功请求" {
			state["requestStatus"] = "完成"
		} else if input == "超时请求" {
			state["hasTimeout"] = true
		}
	case "状态更新函数":
		state["stateInput"] = input
		state["currentState"] = "更新中"
		if input == "最终状态" {
			state["currentState"] = "完成"
		}
	}

	return state
}

func analyzeStateChange(preState, postState map[string]interface{}, expectedOutput string) {
	// 分析状态变化
	changes := 0

	for key, preValue := range preState {
		if postValue, exists := postState[key]; exists {
			if fmt.Sprintf("%v", preValue) != fmt.Sprintf("%v", postValue) {
				log.Printf("    - %s: %v → %v", key, preValue, postValue)
				changes++
			}
		}
	}

	log.Printf("    总变化数: %d", changes)

	// 检查预期输出
	if output, ok := postState["result"]; ok {
		if fmt.Sprintf("%v", output) == expectedOutput {
			log.Printf("    ✅ 输出符合预期")
		} else {
			log.Printf("    ❌ 输出不符合预期: 实际 %v, 预期 %s", output, expectedOutput)
		}
	} else if status, ok := postState["currentState"]; ok {
		if fmt.Sprintf("%v", status) == expectedOutput {
			log.Printf("    ✅ 状态符合预期")
		} else {
			log.Printf("    ❌ 状态不符合预期")
		}
	}
}


*/

// -----------------------------------------------  Debugger.resume  -----------------------------------------------
// === 应用场景 ===
// 1. 继续执行: 在暂停后恢复代码执行
// 2. 调试控制: 控制调试流程的继续
// 3. 错误恢复: 在错误处理后恢复执行
// 4. 条件继续: 在满足条件后继续执行
// 5. 批量调试: 在批量调试中控制执行流程
// 6. 自动测试: 在自动化测试中控制执行

// CDPDebuggerResume 恢复JavaScript执行
// 参数说明:
//   - terminateOnResume: 是否在恢复时终止执行
func CDPDebuggerResume(terminateOnResume bool) (string, error) {
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
		"method": "Debugger.resume",
		"params": {
			"terminateOnResume": %v
		}
	}`, reqID, terminateOnResume)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 resume 请求失败: %w", err)
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
			return "", fmt.Errorf("resume 请求超时")
		}
	}
}

/*

// 示例1: 基本调试流程控制
func exampleBasicDebugFlowControl() {
	// === 应用场景描述 ===
	// 场景: 基本调试流程控制
	// 用途: 控制调试的暂停和恢复流程
	// 优势: 实现标准的调试流程控制
	// 典型工作流: 启用调试 -> 暂停执行 -> 检查状态 -> 恢复执行

	log.Println("基本调试流程控制示例...")

	// 模拟调试流程步骤
	debugSteps := []struct {
		name        string
		description string
		action      func() error
		delay       time.Duration
	}{
		{
			name:        "启用调试器",
			description: "启用调试功能",
			action:      func() error {
				log.Println("  模拟启用调试器...")
				return nil
			},
			delay: 100 * time.Millisecond,
		},
		{
			name:        "设置断点",
			description: "在关键位置设置断点",
			action: func() error {
				log.Println("  模拟设置断点...")
				return nil
			},
			delay: 150 * time.Millisecond,
		},
		{
			name:        "暂停执行",
			description: "暂停JavaScript执行",
			action: func() error {
				log.Println("  模拟暂停执行...")
				return nil
			},
			delay: 200 * time.Millisecond,
		},
		{
			name:        "检查变量",
			description: "检查当前执行状态的变量",
			action: func() error {
				log.Println("  模拟检查变量...")
				return nil
			},
			delay: 250 * time.Millisecond,
		},
		{
			name:        "恢复执行",
			description: "恢复JavaScript执行",
			action: func() error {
				log.Println("  执行恢复...")
				_, err := CDPDebuggerResume(false) // 正常恢复
				return err
			},
			delay: 300 * time.Millisecond,
		},
		{
			name:        "验证结果",
			description: "验证执行结果",
			action: func() error {
				log.Println("  模拟验证结果...")
				return nil
			},
			delay: 150 * time.Millisecond,
		},
	}

	// 执行调试流程
	for i, step := range debugSteps {
		log.Printf("步骤 %d/%d: %s", i+1, len(debugSteps), step.name)
		log.Printf("描述: %s", step.description)

		// 执行步骤动作
		if err := step.action(); err != nil {
			log.Printf("步骤失败: %v", err)
			continue
		}

		log.Printf("步骤成功")

		// 步骤间延迟
		time.Sleep(step.delay)
	}

	log.Println("基本调试流程控制完成")
}


*/

// -----------------------------------------------  Debugger.searchInContent  -----------------------------------------------
// === 应用场景 ===
// 1. 代码搜索: 在脚本中搜索特定文本
// 2. 调试辅助: 查找特定的函数调用或变量
// 3. 代码分析: 分析代码中的模式匹配
// 4. 错误定位: 定位错误消息在代码中的位置
// 5. 重构辅助: 查找需要重构的代码模式
// 6. 安全审计: 搜索潜在的安全问题代码

// CDPDebuggerSearchInContent 在脚本内容中搜索
// 参数说明:
//   - scriptId: 脚本ID
//   - query: 要搜索的字符串
//   - caseSensitive: 是否区分大小写
//   - isRegex: 是否使用正则表达式
func CDPDebuggerSearchInContent(scriptId, query string, caseSensitive, isRegex bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if scriptId == "" {
		return "", fmt.Errorf("脚本ID不能为空")
	}
	if query == "" {
		return "", fmt.Errorf("搜索查询不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 转义查询字符串中的特殊字符
	escapedQuery := strings.ReplaceAll(query, `"`, `\"`)
	escapedQuery = strings.ReplaceAll(escapedQuery, "\n", "\\n")
	escapedQuery = strings.ReplaceAll(escapedQuery, "\t", "\\t")

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.searchInContent",
		"params": {
			"scriptId": "%s",
			"query": "%s",
			"caseSensitive": %v,
			"isRegex": %v
		}
	}`, reqID, scriptId, escapedQuery, caseSensitive, isRegex)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 searchInContent 请求失败: %w", err)
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
			return "", fmt.Errorf("searchInContent 请求超时")
		}
	}
}

// SearchInContentResult 搜索内容结果
type SearchInContentResult struct {
	Result []SearchMatch `json:"result"` // 搜索匹配列表
}

// SearchMatch 搜索匹配
type SearchMatch struct {
	LineNumber  int    `json:"lineNumber"`  // 行号
	LineContent string `json:"lineContent"` // 包含匹配的行内容
}

// ParseSearchInContent 解析搜索内容响应
func ParseSearchInContent(response string) (*SearchInContentResult, error) {
	var data struct {
		Result *SearchInContentResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

/*

// 示例1: 基本代码搜索功能
func exampleBasicCodeSearch() {
	// === 应用场景描述 ===
	// 场景: 基本代码搜索功能
	// 用途: 在JavaScript脚本中搜索特定文本
	// 优势: 快速定位代码中的特定模式
	// 典型工作流: 选择脚本 -> 指定搜索词 -> 执行搜索 -> 查看结果

	log.Println("基本代码搜索功能示例...")

	// 模拟不同的搜索场景
	searchScenarios := []struct {
		name        string
		scriptId    string
		query       string
		description string
		caseSensitive bool
		isRegex     bool
	}{
		{
			name:        "函数名搜索",
			scriptId:    "app-main.js",
			query:       "function processData",
			description: "搜索特定的函数定义",
			caseSensitive: false,
			isRegex:     false,
		},
		{
			name:        "变量声明搜索",
			scriptId:    "utils.js",
			query:       "const\\s+\\w+\\s*=",
			description: "使用正则表达式搜索const变量声明",
			caseSensitive: false,
			isRegex:     true,
		},
		{
			name:        "错误处理搜索",
			scriptId:    "error-handler.js",
			query:       "catch",
			description: "搜索错误处理代码",
			caseSensitive: false,
			isRegex:     false,
		},
		{
			name:        "API调用搜索",
			scriptId:    "api-client.js",
			query:       "fetch\\(",
			description: "使用正则搜索fetch调用",
			caseSensitive: false,
			isRegex:     true,
		},
		{
			name:        "TODO注释搜索",
			scriptId:    "all-scripts",
			query:       "TODO|FIXME|HACK",
			description: "搜索开发注释",
			caseSensitive: false,
			isRegex:     true,
		},
		{
			name:        "精确变量名搜索",
			scriptId:    "data-service.js",
			query:       "userData",
			description: "区分大小写搜索特定变量名",
			caseSensitive: true,
			isRegex:     false,
		},
		{
			name:        "导入语句搜索",
			scriptId:    "module-imports.js",
			query:       "import.*from",
			description: "搜索ES6导入语句",
			caseSensitive: false,
			isRegex:     true,
		},
		{
			name:        "控制台日志搜索",
			scriptId:    "debug-utils.js",
			query:       "console\\.(log|debug|info|warn|error)",
			description: "搜索所有控制台输出",
			caseSensitive: false,
			isRegex:     true,
		},
	}

	// 执行各种搜索
	for i, scenario := range searchScenarios {
		log.Printf("\n=== 搜索场景 %d/%d: %s ===", i+1, len(searchScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("脚本: %s", scenario.scriptId)
		log.Printf("查询: %s", scenario.query)
		log.Printf("区分大小写: %v, 正则表达式: %v", scenario.caseSensitive, scenario.isRegex)

		// 执行搜索
		response, err := CDPDebuggerSearchInContent(
			scenario.scriptId,
			scenario.query,
			scenario.caseSensitive,
			scenario.isRegex,
		)

		if err != nil {
			log.Printf("搜索失败: %v", err)
			continue
		}

		// 解析结果
		result, err := ParseSearchInContent(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 显示搜索结果
		displaySearchResults(result, scenario.query)
	}
}

func displaySearchResults(result *SearchInContentResult, query string) {
	if len(result.Result) == 0 {
		log.Printf("  未找到匹配 '%s' 的结果", query)
		return
	}

	log.Printf("  找到 %d 个匹配:", len(result.Result))

	// 按行号排序
	matches := result.Result
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].LineNumber < matches[j].LineNumber
	})

	// 显示匹配结果
	showCount := min(5, len(matches))
	for i := 0; i < showCount; i++ {
		match := matches[i]
		lineContent := strings.TrimSpace(match.LineContent)
		if len(lineContent) > 60 {
			lineContent = lineContent[:60] + "..."
		}
		log.Printf("    [%d] 行 %d: %s", i+1, match.LineNumber, lineContent)
	}

	if len(matches) > showCount {
		log.Printf("    ... 还有 %d 个匹配", len(matches)-showCount)
	}

	// 统计分析
	analyzeSearchResults(matches, query)
}

func analyzeSearchResults(matches []SearchMatch, query string) {
	// 计算行号分布
	minLine := math.MaxInt32
	maxLine := 0
	lineNumbers := make([]int, 0, len(matches))

	for _, match := range matches {
		lineNumbers = append(lineNumbers, match.LineNumber)
		if match.LineNumber < minLine {
			minLine = match.LineNumber
		}
		if match.LineNumber > maxLine {
			maxLine = match.LineNumber
		}
	}

	// 分析分布
	if len(matches) > 1 {
		lineRange := maxLine - minLine
		density := float64(len(matches)) / float64(lineRange+1) * 100

		log.Printf("  行号范围: %d - %d (%d 行)", minLine, maxLine, lineRange+1)
		log.Printf("  匹配密度: %.1f%%", density)

		// 密度评估
		if density > 50 {
			log.Printf("  ⚠ 匹配非常密集")
		} else if density > 20 {
			log.Printf("  ⚠ 匹配密度适中")
		} else {
			log.Printf("  ✅ 匹配分布稀疏")
		}
	}
}

// 示例2: 代码审计和安全搜索
func exampleCodeAuditAndSecuritySearch() {
	// === 应用场景描述 ===
	// 场景: 代码审计和安全搜索
	// 用途: 搜索代码中的安全问题和漏洞模式
	// 优势: 自动化检测潜在的安全风险
	// 典型工作流: 定义安全模式 -> 执行搜索 -> 分析结果 -> 生成报告

	log.Println("代码审计和安全搜索示例...")

	// 定义安全审计规则
	securityAuditRules := []struct {
		name        string
		pattern     string
		description string
		severity    string
		category    string
		isRegex     bool
	}{
		// XSS相关
		{
			name:        "innerHTML直接赋值",
			pattern:     `\.innerHTML\s*=`,
			description: "直接innerHTML赋值可能导致XSS漏洞",
			severity:    "高危",
			category:    "XSS",
			isRegex:     true,
		},
		{
			name:        "eval函数调用",
			pattern:     `eval\(`,
			description: "eval可能执行恶意代码",
			severity:    "高危",
			category:    "代码注入",
			isRegex:     true,
		},
		{
			name:        "document.write使用",
			pattern:     `document\.write\(`,
			description: "document.write可能被滥用",
			severity:    "中危",
			category:    "XSS",
			isRegex:     true,
		},

		// 敏感数据处理
		{
			name:        "localStorage敏感数据",
			pattern:     `localStorage\.(setItem|getItem)\s*\(\s*['"][^'"]*(password|token|secret|key)`,
			description: "敏感数据不应存储在localStorage",
			severity:    "高危",
			category:    "数据安全",
			isRegex:     true,
		},

		// 网络请求安全
		{
			name:        "fetch无CORS检查",
			pattern:     `fetch\([^)]*\)[^{]*\.then`,
			description: "fetch调用应检查CORS和错误",
			severity:    "中危",
			category:    "网络安全",
			isRegex:     true,
		},

		// 密码学安全
		{
			name:        "弱随机数生成",
			pattern:     `Math\.random\(\)`,
			description: "Math.random()不适合安全用途",
			severity:    "中危",
			category:    "密码学",
			isRegex:     true,
		},

		// 输入验证
		{
			name:        "缺少输入验证",
			pattern:     `function\s+\w+\s*\([^)]*\)[^{]*\{`,
			description: "函数参数缺少验证检查",
			severity:    "低危",
			category:    "输入验证",
			isRegex:     true,
		},

		// 调试信息
		{
			name:        "生产环境调试代码",
			pattern:     `console\.(log|debug|info|warn|error)`,
			description: "生产环境应移除调试代码",
			severity:    "低危",
			category:    "代码质量",
			isRegex:     true,
		},

		// 硬编码凭证
		{
			name:        "硬编码API密钥",
			pattern:     `['"](?:api[_-]?key|secret|token)['"]\s*:\s*['"][^'"]{10,}['"]`,
			description: "发现硬编码的API密钥或令牌",
			severity:    "高危",
			category:    "凭证安全",
			isRegex:     true,
		},
	}

	// 模拟要审计的脚本
	scriptsToAudit := []struct {
		name     string
		scriptId string
		riskLevel string
	}{
		{"用户认证模块", "auth-module.js", "高危"},
		{"API客户端", "api-client.js", "中危"},
		{"数据处理器", "data-processor.js", "中危"},
		{"工具函数库", "utility-library.js", "低危"},
		{"第三方集成", "third-party-integration.js", "高危"},
	}

	// 执行安全审计
	auditResults := make(map[string][]SecurityFinding)

	for _, script := range scriptsToAudit {
		log.Printf("\n=== 安全审计: %s ===", script.name)
		log.Printf("脚本ID: %s", script.scriptId)
		log.Printf("风险等级: %s", script.riskLevel)

		var scriptFindings []SecurityFinding

		// 对每个安全规则执行搜索
		for _, rule := range securityAuditRules {
			log.Printf("检查规则: %s", rule.name)

			response, err := CDPDebuggerSearchInContent(
				script.scriptId,
				rule.pattern,
				false, // 不区分大小写
				rule.isRegex,
			)

			if err != nil {
				log.Printf("  搜索失败: %v", err)
				continue
			}

			result, err := ParseSearchInContent(response)
			if err != nil {
				log.Printf("  解析失败: %v", err)
				continue
			}

			// 记录发现的问题
			if len(result.Result) > 0 {
				finding := SecurityFinding{
					RuleName:    rule.name,
					Pattern:     rule.pattern,
					Description: rule.description,
					Severity:    rule.severity,
					Category:    rule.category,
					MatchCount:  len(result.Result),
					Matches:     result.Result,
				}

				scriptFindings = append(scriptFindings, finding)

				log.Printf("  ❌ 发现 %d 个问题", len(result.Result))
			} else {
				log.Printf("  ✅ 未发现问题")
			}
		}

		auditResults[script.name] = scriptFindings
	}

	// 生成安全审计报告
	generateSecurityAuditReport(auditResults)
}

// SecurityFinding 安全发现
type SecurityFinding struct {
	RuleName    string
	Pattern     string
	Description string
	Severity    string
	Category    string
	MatchCount  int
	Matches     []SearchMatch
}

func generateSecurityAuditReport(results map[string][]SecurityFinding) {
	log.Println("\n=== 安全审计报告 ===")

	totalFindings := 0
	findingsBySeverity := make(map[string]int)
	findingsByCategory := make(map[string]int)

	// 统计总体情况
	for scriptName, findings := range results {
		log.Printf("\n脚本: %s", scriptName)
		log.Printf("发现的问题数: %d", len(findings))

		if len(findings) > 0 {
			log.Printf("详细问题:")
			for _, finding := range findings {
				totalFindings++
				findingsBySeverity[finding.Severity]++
				findingsByCategory[finding.Category]++

				log.Printf("  [%s] %s", finding.Severity, finding.RuleName)
				log.Printf("    描述: %s", finding.Description)
				log.Printf("    匹配数: %d", finding.MatchCount)

				// 显示前几个匹配位置
				showCount := min(2, len(finding.Matches))
				for i := 0; i < showCount; i++ {
					match := finding.Matches[i]
					linePreview := strings.TrimSpace(match.LineContent)
					if len(linePreview) > 50 {
						linePreview = linePreview[:50] + "..."
					}
					log.Printf("    匹配 %d: 行 %d - %s", i+1, match.LineNumber, linePreview)
				}
			}
		} else {
			log.Printf("  ✅ 未发现安全问题")
		}
	}

	// 总体统计
	log.Println("\n=== 总体统计 ===")
	log.Printf("总脚本数: %d", len(results))
	log.Printf("总问题数: %d", totalFindings)

	// 按严重程度统计
	log.Println("\n按严重程度统计:")
	for severity, count := range findingsBySeverity {
		percentage := 0.0
		if totalFindings > 0 {
			percentage = float64(count) / float64(totalFindings) * 100
		}
		log.Printf("  %s: %d (%.1f%%)", severity, count, percentage)
	}

	// 按类别统计
	log.Println("\n按类别统计:")
	for category, count := range findingsByCategory {
		percentage := 0.0
		if totalFindings > 0 {
			percentage = float64(count) / float64(totalFindings) * 100
		}
		log.Printf("  %s: %d (%.1f%%)", category, count, percentage)
	}

	// 安全评级
	log.Println("\n安全评级:")
	if findingsBySeverity["高危"] > 0 {
		log.Println("  🔥 安全评级: 差 (存在高危问题)")
		log.Println("    建议: 立即修复高危问题")
	} else if findingsBySeverity["中危"] > 3 {
		log.Println("  ⚠ 安全评级: 中 (存在多个中危问题)")
		log.Println("    建议: 尽快修复中危问题")
	} else if totalFindings > 0 {
		log.Println("  ⚠ 安全评级: 良 (存在少量低危问题)")
		log.Println("    建议: 计划修复低危问题")
	} else {
		log.Println("  ✅ 安全评级: 优 (未发现问题)")
	}
}


*/

// -----------------------------------------------  Debugger.setAsyncCallStackDepth  -----------------------------------------------
// === 应用场景 ===
// 1. 异步代码调试: 控制异步调用栈的收集深度
// 2. 性能优化: 限制异步调用栈跟踪以减少内存使用
// 3. 复杂异步调试: 调试复杂的Promise链和async/await代码
// 4. 错误追踪: 追踪异步错误的完整调用栈
// 5. 性能分析: 分析异步代码的执行路径
// 6. 内存管理: 控制调试器内存使用

// CDPDebuggerSetAsyncCallStackDepth 设置异步调用栈深度
// 参数说明:
//   - maxDepth: 异步调用栈的最大深度，0表示禁用异步调用栈跟踪
func CDPDebuggerSetAsyncCallStackDepth(maxDepth int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if maxDepth < 0 {
		return "", fmt.Errorf("异步调用栈深度不能为负数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setAsyncCallStackDepth",
		"params": {
			"maxDepth": %d
		}
	}`, reqID, maxDepth)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setAsyncCallStackDepth 请求失败: %w", err)
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
			return "", fmt.Errorf("setAsyncCallStackDepth 请求超时")
		}
	}
}

/*


// 示例1: 异步调用栈深度配置管理
func exampleAsyncCallStackDepthConfiguration() {
	// === 应用场景描述 ===
	// 场景: 异步调用栈深度配置管理
	// 用途: 根据调试需求配置不同的异步调用栈深度
	// 优势: 平衡调试信息和性能开销
	// 典型工作流: 分析调试需求 -> 选择适当深度 -> 配置跟踪 -> 调试异步代码

	log.Println("异步调用栈深度配置管理示例...")

	// 定义不同的调试场景和推荐的深度配置
	debuggingScenarios := []struct {
		name        string
		description string
		recommendedDepth int
		useCase     string
		performanceImpact string
	}{
		{
			name:        "简单异步调试",
			description: "调试简单的Promise或setTimeout",
			recommendedDepth: 3,
			useCase:     "单个异步操作调试",
			performanceImpact: "低",
		},
		{
			name:        "复杂Promise链",
			description: "调试多层Promise.then()链",
			recommendedDepth: 5,
			useCase:     "复杂的异步流程调试",
			performanceImpact: "中",
		},
		{
			name:        "async/await深度调试",
			description: "调试深层嵌套的async/await函数",
			recommendedDepth: 7,
			useCase:     "复杂的异步控制流分析",
			performanceImpact: "中高",
		},
		{
			name:        "性能关键调试",
			description: "在性能敏感场景中调试",
			recommendedDepth: 1,
			useCase:     "生产环境问题诊断",
			performanceImpact: "最低",
		},
		{
			name:        "完整异步追踪",
			description: "需要完整的异步调用历史",
			recommendedDepth: 10,
			useCase:     "复杂的竞态条件分析",
			performanceImpact: "高",
		},
		{
			name:        "禁用异步追踪",
			description: "完全禁用异步调用栈收集",
			recommendedDepth: 0,
			useCase:     "性能测试或内存优化",
			performanceImpact: "无",
		},
		{
			name:        "深度递归异步",
			description: "调试递归的异步操作",
			recommendedDepth: 8,
			useCase:     "递归异步算法分析",
			performanceImpact: "高",
		},
		{
			name:        "事件驱动调试",
			description: "调试事件监听器和回调",
			recommendedDepth: 4,
			useCase:     "事件驱动架构分析",
			performanceImpact: "中",
		},
	}

	// 测试不同的深度配置
	for i, scenario := range debuggingScenarios {
		log.Printf("\n=== 配置场景 %d/%d: %s ===", i+1, len(debuggingScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("使用场景: %s", scenario.useCase)
		log.Printf("推荐深度: %d", scenario.recommendedDepth)
		log.Printf("性能影响: %s", scenario.performanceImpact)

		// 设置异步调用栈深度
		log.Printf("设置异步调用栈深度为 %d...", scenario.recommendedDepth)

		response, err := CDPDebuggerSetAsyncCallStackDepth(scenario.recommendedDepth)
		if err != nil {
			log.Printf("设置失败: %v", err)
			continue
		}

		log.Printf("设置成功: %s", response)

		// 分析配置效果
		analyzeAsyncDepthConfiguration(scenario.recommendedDepth, scenario.performanceImpact)

		// 短暂延迟
		time.Sleep(150 * time.Millisecond)
	}
}

func analyzeAsyncDepthConfiguration(depth int, performanceImpact string) {
	log.Printf("配置分析:")
	log.Printf("  深度: %d", depth)

	switch depth {
	case 0:
		log.Printf("  ✅ 异步调用栈跟踪已禁用")
		log.Printf("    优势: 最小性能开销，最大内存节省")
		log.Printf("    缺点: 无法获取异步调用栈信息")

	case 1:
		log.Printf("  ⚠ 最小跟踪深度")
		log.Printf("    优势: 低性能开销，基本异步信息")
		log.Printf("    缺点: 只能看到最近的异步操作")

	case 2, 3:
		log.Printf("  ⚠ 适中跟踪深度")
		log.Printf("    优势: 平衡性能和信息量")
		log.Printf("    缺点: 可能无法追踪复杂的异步链")

	case 4, 5, 6:
		log.Printf("  ⚠ 深度跟踪")
		log.Printf("    优势: 详细的异步调用信息")
		log.Printf("    缺点: 中等性能开销")

	default:
		log.Printf("  🔥 深度跟踪 (>=7)")
		log.Printf("    优势: 完整的异步调用历史")
		log.Printf("    缺点: 高内存使用和性能开销")
	}

	// 性能建议
	log.Printf("  性能建议:")
	switch performanceImpact {
	case "低":
		log.Printf("    ✅ 适合持续启用")
	case "中":
		log.Printf("    ⚠ 调试期间启用")
	case "中高", "高":
		log.Printf("    ❌ 仅在必要时启用，调试后及时恢复")
	case "无":
		log.Printf("    ✅ 无性能影响")
	}
}

// 示例2: 异步代码性能分析与深度优化
func exampleAsyncCodePerformanceAnalysis() {
	// === 应用场景描述 ===
	// 场景: 异步代码性能分析与深度优化
	// 用途: 分析不同异步调用栈深度对性能的影响
	// 优势: 找到性能与调试信息的最佳平衡点
	// 典型工作流: 测试不同深度 -> 测量性能 -> 分析内存 -> 优化配置

	log.Println("异步代码性能分析与深度优化示例...")

	// 定义性能测试用例
	performanceTests := []struct {
		name         string
		asyncPattern string
		depthTests   []int
		iterations   int
	}{
		{
			name:         "简单Promise链",
			asyncPattern: "Promise.then().then().then()",
			depthTests:   []int{0, 1, 3, 5},
			iterations:   1000,
		},
		{
			name:         "async/await嵌套",
			asyncPattern: "async function with 3 levels of await",
			depthTests:   []int{0, 2, 4, 6},
			iterations:   500,
		},
		{
			name:         "事件监听器链",
			asyncPattern: "Event listeners with callbacks",
			depthTests:   []int{0, 1, 3, 5},
			iterations:   800,
		},
		{
			name:         "混合异步模式",
			asyncPattern: "Promise + setTimeout + async/await",
			depthTests:   []int{0, 3, 6, 9},
			iterations:   300,
		},
		{
			name:         "递归异步操作",
			asyncPattern: "Recursive async functions",
			depthTests:   []int{0, 2, 5, 8},
			iterations:   200,
		},
	}

	// 执行性能测试
	for testIndex, test := range performanceTests {
		log.Printf("\n=== 性能测试 %d/%d: %s ===",
			testIndex+1, len(performanceTests), test.name)
		log.Printf("异步模式: %s", test.asyncPattern)
		log.Printf("迭代次数: %d", test.iterations)

		var testResults []PerformanceResult

		// 测试不同的深度
		for _, depth := range test.depthTests {
			log.Printf("\n测试深度: %d", depth)

			// 设置异步调用栈深度
			_, err := CDPDebuggerSetAsyncCallStackDepth(depth)
			if err != nil {
				log.Printf("设置深度失败: %v", err)
				continue
			}

			// 执行性能测试
			result := runAsyncPerformanceTest(test.name, depth, test.iterations)
			testResults = append(testResults, result)

			// 显示即时结果
			displayPerformanceResult(result)
		}

		// 分析最佳深度
		analyzeOptimalDepth(testResults, test.name)
	}
}

// PerformanceResult 性能结果
type PerformanceResult struct {
	TestName    string
	Depth       int
	Iterations  int
	TotalTime   time.Duration
	AverageTime time.Duration
	MemoryUsage string
	SuccessRate float64
	Metadata    map[string]interface{}
}

func runAsyncPerformanceTest(testName string, depth, iterations int) PerformanceResult {
	startTime := time.Now()

	// 模拟异步操作执行
	successCount := 0
	totalOperations := iterations

	for i := 0; i < iterations; i++ {
		// 模拟异步操作
		opTime := simulateAsyncOperation(depth, i)

		// 记录成功
		if opTime > 0 {
			successCount++
		}

		// 短暂延迟，模拟实际执行
		time.Sleep(time.Microsecond * 10)
	}

	totalTime := time.Since(startTime)
	averageTime := totalTime / time.Duration(iterations)
	successRate := float64(successCount) / float64(totalOperations) * 100

	// 模拟内存使用
	memoryUsage := "低"
	if depth >= 5 {
		memoryUsage = "中"
	}
	if depth >= 8 {
		memoryUsage = "高"
	}
	if depth == 0 {
		memoryUsage = "最低"
	}

	return PerformanceResult{
		TestName:    testName,
		Depth:       depth,
		Iterations:  iterations,
		TotalTime:   totalTime,
		AverageTime: averageTime,
		MemoryUsage: memoryUsage,
		SuccessRate: successRate,
		Metadata: map[string]interface{}{
			"testCompleted": true,
			"timestamp":     time.Now().UnixNano(),
		},
	}
}

func simulateAsyncOperation(depth, iteration int) time.Duration {
	// 模拟异步操作的执行时间
	// 深度越大，模拟的操作越复杂
	baseTime := 100 * time.Microsecond
	complexityFactor := time.Duration(depth) * 5 * time.Microsecond
	randomFactor := time.Duration(rand.Intn(50)) * time.Microsecond

	return baseTime + complexityFactor + randomFactor
}

func displayPerformanceResult(result PerformanceResult) {
	log.Printf("  性能结果:")
	log.Printf("    总时间: %v", result.TotalTime)
	log.Printf("    平均时间: %v", result.AverageTime)
	log.Printf("    内存使用: %s", result.MemoryUsage)
	log.Printf("    成功率: %.1f%%", result.SuccessRate)

	// 性能评级
	log.Printf("    性能评级:")
	avgMicroseconds := result.AverageTime.Microseconds()

	switch {
	case avgMicroseconds < 150:
		log.Printf("      ✅ 优秀 (<150µs)")
	case avgMicroseconds < 300:
		log.Printf("      ⚠ 良好 (150-300µs)")
	case avgMicroseconds < 500:
		log.Printf("      ⚠ 一般 (300-500µs)")
	default:
		log.Printf("      ❌ 较差 (>500µs)")
	}
}

func analyzeOptimalDepth(results []PerformanceResult, testName string) {
	if len(results) == 0 {
		return
	}

	log.Printf("\n深度优化分析 (%s):", testName)

	// 找到最佳性能的深度
	bestDepth := 0
	bestTime := time.Duration(math.MaxInt64)
	bestBalanceDepth := 0
	bestBalanceScore := 0.0

	for _, result := range results {
		// 检查最短时间
		if result.AverageTime < bestTime {
			bestTime = result.AverageTime
			bestDepth = result.Depth
		}

		// 计算平衡分数（考虑深度和时间）
		// 深度越大，分数越低（因为需要更多资源）
		// 时间越短，分数越高
		depthFactor := 1.0 / float64(result.Depth+1) // 避免除零
		timeFactor := 1.0 / result.AverageTime.Seconds()
		balanceScore := depthFactor * timeFactor * 1000

		if balanceScore > bestBalanceScore {
			bestBalanceScore = balanceScore
			bestBalanceDepth = result.Depth
		}
	}

	log.Printf("  最快深度: %d (%v/op)", bestDepth, bestTime)
	log.Printf("  最佳平衡深度: %d (分数: %.2f)", bestBalanceDepth, bestBalanceScore)

	// 推荐深度
	log.Printf("  推荐配置:")
	if bestDepth == 0 {
		log.Printf("    ✅ 深度 0: 最大性能，无异步栈信息")
	} else if bestDepth <= 3 {
		log.Printf("    ⚠ 深度 %d: 高性能，基本异步信息", bestDepth)
	} else if bestDepth <= 6 {
		log.Printf("    ⚠ 深度 %d: 平衡性能和信息", bestDepth)
	} else {
		log.Printf("    ❌ 深度 %d: 高开销，完整异步信息", bestDepth)
	}
}


*/

// -----------------------------------------------  Debugger.setBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 精确断点: 在代码的特定位置设置断点
// 2. 条件调试: 设置条件断点，仅在特定条件下暂停
// 3. 调试自动化: 自动化设置断点进行测试
// 4. 动态调试: 在运行时动态添加断点
// 5. 错误重现: 在错误发生位置设置断点
// 6. 性能分析: 在关键路径设置断点分析性能

// CDPDebuggerSetBreakpoint 设置断点
// 参数说明:
//   - scriptId: 脚本ID
//   - lineNumber: 行号（0-based）
//   - columnNumber: 列号（0-based）
//   - condition: 断点条件表达式（可选）
func CDPDebuggerSetBreakpoint(scriptId string, lineNumber, columnNumber int, condition string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if scriptId == "" {
		return "", fmt.Errorf("脚本ID不能为空")
	}
	if lineNumber < 0 {
		return "", fmt.Errorf("行号必须是非负整数")
	}
	if columnNumber < 0 {
		return "", fmt.Errorf("列号必须是非负整数")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建位置对象
	location := map[string]interface{}{
		"scriptId":     scriptId,
		"lineNumber":   lineNumber,
		"columnNumber": columnNumber,
	}

	locationJSON, err := json.Marshal(location)
	if err != nil {
		return "", fmt.Errorf("序列化位置对象失败: %w", err)
	}

	// 构建参数
	params := fmt.Sprintf(`"location": %s`, locationJSON)

	if condition != "" {
		// 转义条件中的特殊字符
		escapedCondition := strings.ReplaceAll(condition, `"`, `\"`)
		escapedCondition = strings.ReplaceAll(escapedCondition, "\n", "\\n")
		escapedCondition = strings.ReplaceAll(escapedCondition, "\t", "\\t")
		params += fmt.Sprintf(`, "condition": "%s"`, escapedCondition)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setBreakpoint",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("setBreakpoint 请求超时")
		}
	}
}

// SetBreakpointResult 设置断点结果
type SetBreakpointResult struct {
	BreakpointID   string    `json:"breakpointId"`   // 断点ID
	ActualLocation *Location `json:"actualLocation"` // 实际位置
}

// Location 位置
type Location struct {
	ScriptID     string `json:"scriptId"`     // 脚本ID
	LineNumber   int    `json:"lineNumber"`   // 行号
	ColumnNumber int    `json:"columnNumber"` // 列号
}

// ParseSetBreakpoint 解析设置断点响应
func ParseSetBreakpoint(response string) (*SetBreakpointResult, error) {
	var data struct {
		Result *SetBreakpointResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

/*

// 示例1: 基本断点设置
func exampleBasicBreakpointSetup() {
	// === 应用场景描述 ===
	// 场景: 基本断点设置
	// 用途: 在代码的关键位置设置断点
	// 优势: 可以精确控制调试暂停的位置
	// 典型工作流: 选择位置 -> 设置断点 -> 触发执行 -> 调试分析

	log.Println("基本断点设置示例...")

	// 定义要设置断点的关键位置
	breakpointLocations := []struct {
		name        string
		scriptId    string
		lineNumber  int
		columnNumber int
		description string
		importance  string
	}{
		{
			name:        "函数入口",
			scriptId:    "app-main.js",
			lineNumber:  42,
			columnNumber: 0,
			description: "主处理函数的入口点",
			importance:  "高",
		},
		{
			name:        "错误处理",
			scriptId:    "error-handler.js",
			lineNumber:  18,
			columnNumber: 4,
			description: "错误处理函数的开始",
			importance:  "高",
		},
		{
			name:        "数据验证",
			scriptId:    "validation.js",
			lineNumber:  33,
			columnNumber: 8,
			description: "数据验证逻辑的关键检查点",
			importance:  "中",
		},
		{
			name:        "循环开始",
			scriptId:    "data-processor.js",
			lineNumber:  67,
			columnNumber: 0,
			description: "主循环的起始位置",
			importance:  "中",
		},
		{
			name:        "API调用前",
			scriptId:    "api-client.js",
			lineNumber:  89,
			columnNumber: 12,
			description: "网络请求发送前的位置",
			importance:  "高",
		},
		{
			name:        "状态更新",
			scriptId:    "state-manager.js",
			lineNumber:  51,
			columnNumber: 4,
			description: "应用状态更新的关键点",
			importance:  "中",
		},
		{
			name:        "计算结果",
			scriptId:    "calculator.js",
			lineNumber:  24,
			columnNumber: 16,
			description: "重要计算结果的赋值位置",
			importance:  "低",
		},
		{
			name:        "资源清理",
			scriptId:    "cleanup.js",
			lineNumber:  12,
			columnNumber: 0,
			description: "资源清理操作的开始",
			importance:  "中",
		},
	}

	// 设置断点
	for i, location := range breakpointLocations {
		log.Printf("\n=== 设置断点 %d/%d: %s ===", i+1, len(breakpointLocations), location.name)
		log.Printf("描述: %s", location.description)
		log.Printf("重要性: %s", location.importance)
		log.Printf("位置: %s:%d:%d", location.scriptId, location.lineNumber, location.columnNumber)

		// 设置断点
		response, err := CDPDebuggerSetBreakpoint(
			location.scriptId,
			location.lineNumber,
			location.columnNumber,
			"", // 无条件
		)

		if err != nil {
			log.Printf("设置断点失败: %v", err)
			continue
		}

		// 解析结果
		result, err := ParseSetBreakpoint(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 显示断点信息
		displayBreakpointInfo(result, location.name)
	}
}

func displayBreakpointInfo(result *SetBreakpointResult, breakpointName string) {
	log.Printf("断点设置成功:")
	log.Printf("  断点ID: %s", result.BreakpointID)

	if result.ActualLocation != nil {
		loc := result.ActualLocation
		log.Printf("  实际位置: %s:%d:%d",
			loc.ScriptID, loc.LineNumber, loc.ColumnNumber)

		// 检查位置是否匹配
		if loc.ScriptID != "" && loc.LineNumber >= 0 && loc.ColumnNumber >= 0 {
			log.Printf("  ✅ 断点位置有效")
		}
	}

	// 断点状态
	log.Printf("  状态: 已激活")
	log.Printf("  名称: %s", breakpointName)
}

// 示例2: 条件断点设置
func exampleConditionalBreakpointSetup() {
	// === 应用场景描述 ===
	// 场景: 条件断点设置
	// 用途: 设置仅在特定条件下触发的断点
	// 优势: 可以过滤不相关的执行，提高调试效率
	// 典型工作流: 分析条件 -> 设置条件断点 -> 触发条件 -> 调试特定场景

	log.Println("条件断点设置示例...")

	// 定义条件断点场景
	conditionalBreakpoints := []struct {
		name        string
		scriptId    string
		lineNumber  int
		columnNumber int
		condition   string
		description string
		useCase     string
	}{
		{
			name:        "特定用户调试",
			scriptId:    "user-service.js",
			lineNumber:  45,
			columnNumber: 0,
			condition:   "user.id === 'debug-user-001'",
			description: "仅针对特定测试用户触发断点",
			useCase:     "用户特定问题调试",
		},
		{
			name:        "错误值检查",
			scriptId:    "data-validator.js",
			lineNumber:  28,
			columnNumber: 4,
			condition:   "value === null || value === undefined",
			description: "当值为null或undefined时触发断点",
			useCase:     "空值问题调试",
		},
		{
			name:        "性能阈值",
			scriptId:    "performance-monitor.js",
			lineNumber:  63,
			columnNumber: 8,
			condition:   "executionTime > 1000", // 1秒
			description: "执行时间超过阈值时触发断点",
			useCase:     "性能问题调试",
		},
		{
			name:        "特定状态",
			scriptId:    "state-machine.js",
			lineNumber:  37,
			columnNumber: 12,
			condition:   "currentState === 'ERROR' || currentState === 'FAILED'",
			description: "进入错误或失败状态时触发断点",
			useCase:     "状态机错误调试",
		},
		{
			name:        "数据边界",
			scriptId:    "array-processor.js",
			lineNumber:  52,
			columnNumber: 0,
			condition:   "index === 0 || index === array.length - 1",
			description: "处理数组第一个或最后一个元素时触发",
			useCase:     "边界条件调试",
		},
		{
			name:        "资源限制",
			scriptId:    "resource-manager.js",
			lineNumber:  19,
			columnNumber: 4,
			condition:   "memoryUsage > 100 * 1024 * 1024", // 100MB
			description: "内存使用超过限制时触发断点",
			useCase:     "内存泄漏调试",
		},
		{
			name:        "特定输入模式",
			scriptId:    "input-processor.js",
			lineNumber:  41,
			columnNumber: 8,
			condition:   "input.includes('test') || input.includes('debug')",
			description: "输入包含特定关键词时触发断点",
			useCase:     "输入处理调试",
		},
		{
			name:        "网络错误",
			scriptId:    "api-handler.js",
			lineNumber:  76,
			columnNumber: 0,
			condition:   "response.status >= 400",
			description: "HTTP响应错误时触发断点",
			useCase:     "网络错误调试",
		},
	}

	// 设置条件断点
	for i, bp := range conditionalBreakpoints {
		log.Printf("\n=== 条件断点 %d/%d: %s ===", i+1, len(conditionalBreakpoints), bp.name)
		log.Printf("描述: %s", bp.description)
		log.Printf("使用场景: %s", bp.useCase)
		log.Printf("位置: %s:%d:%d", bp.scriptId, bp.lineNumber, bp.columnNumber)
		log.Printf("条件: %s", bp.condition)

		// 设置条件断点
		response, err := CDPDebuggerSetBreakpoint(
			bp.scriptId,
			bp.lineNumber,
			bp.columnNumber,
			bp.condition,
		)

		if err != nil {
			log.Printf("设置条件断点失败: %v", err)
			continue
		}

		// 解析结果
		result, err := ParseSetBreakpoint(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 显示条件断点信息
		displayConditionalBreakpointInfo(result, bp.name, bp.condition)

		// 分析条件有效性
		analyzeBreakpointCondition(bp.condition, bp.useCase)
	}
}

func displayConditionalBreakpointInfo(result *SetBreakpointResult, breakpointName, condition string) {
	log.Printf("条件断点设置成功:")
	log.Printf("  断点ID: %s", result.BreakpointID)
	log.Printf("  条件: %s", condition)

	if result.ActualLocation != nil {
		loc := result.ActualLocation
		log.Printf("  实际位置: %s:%d:%d",
			loc.ScriptID, loc.LineNumber, loc.ColumnNumber)
	}

	// 条件断点特性
	log.Printf("  类型: 条件断点")
	log.Printf("  名称: %s", breakpointName)
	log.Printf("  触发条件: 当表达式为true时暂停")
}

func analyzeBreakpointCondition(condition, useCase string) {
	log.Printf("条件分析:")

	// 分析条件复杂度
	complexity := "简单"
	if strings.Contains(condition, "&&") || strings.Contains(condition, "||") {
		complexity = "中等"
	}
	if strings.Contains(condition, ">") || strings.Contains(condition, "<") ||
	   strings.Contains(condition, "===") || strings.Contains(condition, "!==") {
		complexity = "中等"
	}
	if strings.Count(condition, "&&") + strings.Count(condition, "||") > 2 {
		complexity = "复杂"
	}

	log.Printf("  复杂度: %s", complexity)

	// 评估条件
	switch useCase {
	case "用户特定问题调试":
		log.Printf("  ⚠ 用户特定调试: 需要特定测试数据")
		log.Printf("    建议: 准备测试用户数据")

	case "空值问题调试":
		log.Printf("  ✅ 空值检查: 常见调试场景")
		log.Printf("    建议: 可以长期启用")

	case "性能问题调试":
		log.Printf("  ⚠ 性能阈值: 可能频繁触发")
		log.Printf("    建议: 适当调整阈值")

	case "状态机错误调试":
		log.Printf("  ✅ 状态检查: 有效的错误捕获")
		log.Printf("    建议: 结合日志使用")

	case "边界条件调试":
		log.Printf("  ⚠ 边界条件: 针对特定情况")
		log.Printf("    建议: 在边界测试时启用")

	case "内存泄漏调试":
		log.Printf("  🔥 内存检查: 重要但可能影响性能")
		log.Printf("    建议: 在内存分析时临时启用")
	}

	// 性能影响
	if complexity == "复杂" {
		log.Printf("  ⚠ 复杂条件可能影响执行性能")
	}
}

// 示例3: 断点管理和状态跟踪
func exampleBreakpointManagementAndTracking() {
	// === 应用场景描述 ===
	// 场景: 断点管理和状态跟踪
	// 用途: 管理多个断点并跟踪其状态
	// 优势: 可以统一管理调试断点，分析断点效果
	// 典型工作流: 批量设置断点 -> 跟踪触发 -> 统计效果 -> 优化配置

	log.Println("断点管理和状态跟踪示例...")

	// 模拟断点管理场景
	managementScenarios := []struct {
		name        string
		description string
		breakpoints []struct {
			id          string
			name        string
			scriptId    string
			line        int
			column      int
			condition   string
			priority    string
		}
		managementStrategy string
	}{
		{
			name:        "函数调用追踪",
			description: "追踪关键函数的调用链",
			managementStrategy: "按调用顺序设置断点",
			breakpoints: []struct {
				id        string
				name      string
				scriptId  string
				line      int
				column    int
				condition string
				priority  string
			}{
				{"bp-func-entry", "函数入口", "workflow.js", 15, 0, "", "高"},
				{"bp-data-prep", "数据准备", "workflow.js", 28, 4, "", "中"},
				{"bp-process", "处理逻辑", "workflow.js", 42, 8, "", "高"},
				{"bp-result", "结果处理", "workflow.js", 57, 0, "", "中"},
				{"bp-cleanup", "清理操作", "workflow.js", 73, 4, "", "低"},
			},
		},
		{
			name:        "错误处理路径",
			description: "覆盖所有错误处理路径",
			managementStrategy: "在错误处理点设置断点",
			breakpoints: []struct {
				id        string
				name      string
				scriptId  string
				line      int
				column    int
				condition string
				priority  string
			}{
				{"bp-error-entry", "错误入口", "error-handler.js", 10, 0, "error !== null", "高"},
				{"bp-validation", "验证错误", "error-handler.js", 25, 4, "error.type === 'VALIDATION'", "中"},
				{"bp-network", "网络错误", "error-handler.js", 39, 8, "error.type === 'NETWORK'", "中"},
				{"bp-timeout", "超时错误", "error-handler.js", 52, 0, "error.type === 'TIMEOUT'", "低"},
				{"bp-final", "最终处理", "error-handler.js", 68, 4, "", "高"},
			},
		},
		{
			name:        "性能热点分析",
			description: "在性能关键路径设置断点",
			managementStrategy: "结合性能分析设置断点",
			breakpoints: []struct {
				id        string
				name      string
				scriptId  string
				line      int
				column    int
				condition string
				priority  string
			}{
				{"bp-loop-start", "循环开始", "perf-critical.js", 33, 0, "i === 0", "中"},
				{"bp-heavy-op", "重操作", "perf-critical.js", 47, 8, "data.length > 1000", "高"},
				{"bp-memory", "内存操作", "perf-critical.js", 62, 4, "memoryAllocated > 1024 * 1024", "高"},
				{"bp-io", "IO操作", "perf-critical.js", 78, 0, "", "中"},
				{"bp-end", "结束点", "perf-critical.js", 91, 4, "", "低"},
			},
		},
		{
			name:        "状态转换追踪",
			description: "追踪状态机的状态转换",
			managementStrategy: "在状态转换点设置断点",
			breakpoints: []struct {
				id        string
				name      string
				scriptId  string
				line      int
				column    int
				condition string
				priority  string
			}{
				{"bp-init", "初始状态", "state-machine.js", 12, 0, "fromState === 'INIT'", "高"},
				{"bp-processing", "处理中", "state-machine.js", 28, 4, "toState === 'PROCESSING'", "中"},
				{"bp-error", "错误状态", "state-machine.js", 45, 8, "toState === 'ERROR'", "高"},
				{"bp-done", "完成状态", "state-machine.js", 61, 0, "toState === 'DONE'", "中"},
				{"bp-reset", "重置状态", "state-machine.js", 77, 4, "toState === 'RESET'", "低"},
			},
		},
	}

	// 管理每个场景的断点
	for scenarioIndex, scenario := range managementScenarios {
		log.Printf("\n=== 断点管理场景 %d/%d: %s ===",
			scenarioIndex+1, len(managementScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("管理策略: %s", scenario.managementStrategy)

		// 设置该场景的所有断点
		var breakpointResults []BreakpointInfo

		for bpIndex, bp := range scenario.breakpoints {
			log.Printf("\n设置断点 %d/%d: %s",
				bpIndex+1, len(scenario.breakpoints), bp.name)
			log.Printf("优先级: %s", bp.priority)
			log.Printf("位置: %s:%d:%d", bp.scriptId, bp.line, bp.column)

			if bp.condition != "" {
				log.Printf("条件: %s", bp.condition)
			}

			// 设置断点
			response, err := CDPDebuggerSetBreakpoint(
				bp.scriptId,
				bp.line,
				bp.column,
				bp.condition,
			)

			var breakpointInfo BreakpointInfo

			if err != nil {
				log.Printf("设置断点失败: %v", err)
				breakpointInfo = BreakpointInfo{
					ID:       bp.id,
					Name:     bp.name,
					Success:  false,
					Error:    err.Error(),
					Priority: bp.priority,
				}
			} else {
				result, err := ParseSetBreakpoint(response)
				if err != nil {
					log.Printf("解析结果失败: %v", err)
					breakpointInfo = BreakpointInfo{
						ID:       bp.id,
						Name:     bp.name,
						Success:  false,
						Error:    err.Error(),
						Priority: bp.priority,
					}
				} else {
					log.Printf("断点设置成功: %s", result.BreakpointID)

					breakpointInfo = BreakpointInfo{
						ID:           result.BreakpointID,
						Name:         bp.name,
						Success:      true,
						ActualLine:   result.ActualLocation.LineNumber,
						ActualColumn: result.ActualLocation.ColumnNumber,
						Priority:     bp.priority,
					}
				}
			}

			breakpointResults = append(breakpointResults, breakpointInfo)

			// 短暂延迟
			time.Sleep(100 * time.Millisecond)
		}

		// 分析断点设置结果
		analyzeBreakpointSetupResults(breakpointResults, scenario.name)

		// 生成管理报告
		generateBreakpointManagementReport(breakpointResults, scenario)
	}
}

// BreakpointInfo 断点信息
type BreakpointInfo struct {
	ID           string
	Name         string
	Success      bool
	ActualLine   int
	ActualColumn int
	Error        string
	Priority     string
	HitCount     int
	LastHitTime  time.Time
}

func analyzeBreakpointSetupResults(results []BreakpointInfo, scenarioName string) {
	log.Printf("\n断点设置结果分析 (%s):", scenarioName)

	totalBreakpoints := len(results)
	successfulBreakpoints := 0
	failedBreakpoints := 0

	priorityCounts := make(map[string]int)
	successByPriority := make(map[string]int)

	for _, bp := range results {
		if bp.Success {
			successfulBreakpoints++
			successByPriority[bp.Priority]++
		} else {
			failedBreakpoints++
		}
		priorityCounts[bp.Priority]++
	}

	log.Printf("  总数: %d", totalBreakpoints)
	log.Printf("  成功: %d (%.1f%%)", successfulBreakpoints,
		float64(successfulBreakpoints)/float64(totalBreakpoints)*100)
	log.Printf("  失败: %d (%.1f%%)", failedBreakpoints,
		float64(failedBreakpoints)/float64(totalBreakpoints)*100)

	log.Printf("  按优先级统计:")
	for priority, count := range priorityCounts {
		successCount := successByPriority[priority]
		successRate := 0.0
		if count > 0 {
			successRate = float64(successCount) / float64(count) * 100
		}
		log.Printf("    %s: %d/%d (%.1f%%)", priority, successCount, count, successRate)
	}

	// 成功率评估
	successRate := float64(successfulBreakpoints) / float64(totalBreakpoints) * 100
	if successRate >= 90 {
		log.Printf("  ✅ 断点设置成功率优秀")
	} else if successRate >= 70 {
		log.Printf("  ⚠ 断点设置成功率良好")
	} else if successRate >= 50 {
		log.Printf("  ⚠ 断点设置成功率一般")
	} else {
		log.Printf("  ❌ 断点设置成功率较差")
	}
}

func generateBreakpointManagementReport(results []BreakpointInfo, scenario struct {
	name        string
	description string
	breakpoints []struct {
		id        string
		name      string
		scriptId  string
		line      int
		column    int
		condition string
		priority  string
	}
	managementStrategy string
}) {
	log.Printf("\n=== 断点管理报告: %s ===", scenario.name)

	// 总体统计
	successCount := 0
	highPrioritySuccess := 0
	highPriorityTotal := 0

	for _, bp := range results {
		if bp.Success {
			successCount++
		}
		if bp.Priority == "高" {
			highPriorityTotal++
			if bp.Success {
				highPrioritySuccess++
			}
		}
	}

	log.Printf("总体统计:")
	log.Printf("  场景: %s", scenario.name)
	log.Printf("  描述: %s", scenario.description)
	log.Printf("  策略: %s", scenario.managementStrategy)
	log.Printf("  总断点数: %d", len(results))
	log.Printf("  成功数: %d", successCount)

	if highPriorityTotal > 0 {
		highPriorityRate := float64(highPrioritySuccess) / float64(highPriorityTotal) * 100
		log.Printf("  高优先级成功率: %.1f%% (%d/%d)", highPriorityRate,
			highPrioritySuccess, highPriorityTotal)
	}

	// 详细结果
	log.Printf("\n详细结果:")
	for i, bp := range results {
		status := "✅ 成功"
		if !bp.Success {
			status = "❌ 失败"
		}

		log.Printf("%d. [%s] %s (优先级: %s)", i+1, status, bp.Name, bp.Priority)

		if bp.Success {
			log.Printf("    断点ID: %s", bp.ID)
			if bp.ActualLine >= 0 {
				log.Printf("    实际位置: 行 %d, 列 %d", bp.ActualLine, bp.ActualColumn)
			}
		} else {
			log.Printf("    错误: %s", bp.Error)
		}
	}

	// 建议
	log.Printf("\n建议:")
	if successCount == len(results) {
		log.Printf("  ✅ 所有断点设置成功")
		log.Printf("    建议: 可以开始调试")
	} else if highPrioritySuccess == highPriorityTotal {
		log.Printf("  ⚠ 高优先级断点全部成功")
		log.Printf("    建议: 可以开始调试，但需注意低优先级断点")
	} else {
		log.Printf("  ❌ 有高优先级断点失败")
		log.Printf("    建议: 先解决高优先级断点问题")
	}
}

*/

// -----------------------------------------------  Debugger.setBreakpointByUrl  -----------------------------------------------
// === 应用场景 ===
// 1. URL模式断点: 在匹配特定URL模式的脚本上设置断点
// 2. 批量断点: 在多个脚本的同位置设置断点
// 3. 动态脚本: 在动态加载的脚本上设置断点
// 4. 生产调试: 在生产环境特定文件上设置断点
// 5. 正则匹配: 使用正则表达式匹配多个URL
// 6. 持久断点: 设置跨页面重载有效的断点

// CDPDebuggerSetBreakpointByUrl 通过URL设置断点
// 参数说明:
//   - lineNumber: 行号
//   - url: URL字符串（与urlRegex二选一）
//   - urlRegex: URL正则表达式（与url二选一）
//   - scriptHash: 脚本哈希（可选）
//   - condition: 断点条件表达式（可选）
//   - columnNumber: 列号（可选）
func CDPDebuggerSetBreakpointByUrl(lineNumber int, url, urlRegex, scriptHash, condition string, columnNumber int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if lineNumber < 0 {
		return "", fmt.Errorf("行号必须是非负整数")
	}
	if url == "" && urlRegex == "" {
		return "", fmt.Errorf("必须提供url或urlRegex参数")
	}
	if columnNumber < 0 {
		columnNumber = 0 // 默认为0
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := make(map[string]interface{})
	params["lineNumber"] = lineNumber

	if url != "" {
		params["url"] = url
	}
	if urlRegex != "" {
		params["urlRegex"] = urlRegex
	}
	if scriptHash != "" {
		params["scriptHash"] = scriptHash
	}
	if condition != "" {
		// 转义条件中的特殊字符
		escapedCondition := strings.ReplaceAll(condition, `"`, `\"`)
		escapedCondition = strings.ReplaceAll(escapedCondition, "\n", "\\n")
		escapedCondition = strings.ReplaceAll(escapedCondition, "\t", "\\t")
		params["condition"] = escapedCondition
	}
	if columnNumber > 0 {
		params["columnNumber"] = columnNumber
	}

	// 序列化参数
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setBreakpointByUrl",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setBreakpointByUrl 请求失败: %w", err)
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
			return "", fmt.Errorf("setBreakpointByUrl 请求超时")
		}
	}
}

// SetBreakpointByUrlResult 通过URL设置断点结果
type SetBreakpointByUrlResult struct {
	BreakpointID string     `json:"breakpointId"` // 断点ID
	Locations    []Location `json:"locations"`    // 断点解析到的位置列表
}

// ParseSetBreakpointByUrl 解析通过URL设置断点响应
func ParseSetBreakpointByUrl(response string) (*SetBreakpointByUrlResult, error) {
	var data struct {
		Result *SetBreakpointByUrlResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
}

/*

// 示例1: 基于URL模式的基本断点设置
func exampleBasicUrlPatternBreakpoint() {
	// === 应用场景描述 ===
	// 场景: 基于URL模式的基本断点设置
	// 用途: 在特定URL的脚本上设置断点
	// 优势: 可以在加载特定文件时自动设置断点
	// 典型工作流: 指定URL模式 -> 设置断点 -> 脚本加载 -> 自动断点生效

	log.Println("基于URL模式的基本断点设置示例...")

	// 定义URL断点场景
	urlBreakpointScenarios := []struct {
		name        string
		description string
		lineNumber  int
		urlPattern  string
		isRegex     bool
		condition   string
		columnNumber int
		useCase     string
	}{
		{
			name:        "特定脚本文件",
			description: "在特定脚本文件的特定行设置断点",
			lineNumber:  42,
			urlPattern:  "https://example.com/js/app.js",
			isRegex:     false,
			condition:   "",
			columnNumber: 0,
			useCase:     "调试特定文件的特定问题",
		},
		{
			name:        "所有JS文件第10行",
			description: "在所有JavaScript文件的第10行设置断点",
			lineNumber:  10,
			urlPattern:  ".*\\.js$",
			isRegex:     true,
			condition:   "",
			columnNumber: 0,
			useCase:     "分析所有JS文件的初始化代码",
		},
		{
			name:        "API模块入口点",
			description: "在API模块文件的入口函数设置断点",
			lineNumber:  15,
			urlPattern:  ".*api.*\\.js$",
			isRegex:     true,
			condition:   "method === 'GET'",
			columnNumber: 4,
			useCase:     "调试API调用流程",
		},
		{
			name:        "错误处理文件",
			description: "在错误处理文件的错误捕获点设置条件断点",
			lineNumber:  28,
			urlPattern:  ".*error.*\\.js$",
			isRegex:     true,
			condition:   "error.code === 500",
			columnNumber: 8,
			useCase:     "捕获特定错误条件",
		},
		{
			name:        "第三方库调试",
			description: "在第三方库的特定版本设置断点",
			lineNumber:  5,
			urlPattern:  ".*jquery-3\\.6\\.0\\.min\\.js$",
			isRegex:     true,
			condition:   "",
			columnNumber: 0,
			useCase:     "调试第三方库问题",
		},
		{
			name:        "动态加载模块",
			description: "在动态加载的模块文件设置断点",
			lineNumber:  1,
			urlPattern:  ".*module-.*\\.js$",
			isRegex:     true,
			condition:   "initComplete === false",
			columnNumber: 0,
			useCase:     "调试动态模块加载",
		},
		{
			name:        "内联脚本",
			description: "在内联脚本设置断点",
			lineNumber:  0,
			urlPattern:  "about:blank",
			isRegex:     false,
			condition:   "",
			columnNumber: 0,
			useCase:     "调试内联脚本代码",
		},
		{
			name:        "CDN资源",
			description: "在CDN托管的资源上设置断点",
			lineNumber:  100,
			urlPattern:  "https://cdn\\.example\\.com/.*\\.js",
			isRegex:     true,
			condition:   "debug === true",
			columnNumber: 12,
			useCase:     "调试CDN托管的脚本",
		},
	}

	// 设置URL模式断点
	for i, scenario := range urlBreakpointScenarios {
		log.Printf("\n=== URL断点场景 %d/%d: %s ===", i+1, len(urlBreakpointScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("使用场景: %s", scenario.useCase)
		log.Printf("位置: 行 %d, 列 %d", scenario.lineNumber, scenario.columnNumber)
		log.Printf("URL模式: %s (正则: %v)", scenario.urlPattern, scenario.isRegex)
		if scenario.condition != "" {
			log.Printf("条件: %s", scenario.condition)
		}

		var response string
		var err error

		// 根据是否是正则表达式调用不同的参数
		if scenario.isRegex {
			response, err = CDPDebuggerSetBreakpointByUrl(
				scenario.lineNumber,
				"", // url参数为空
				scenario.urlPattern, // 使用urlRegex
				"", // scriptHash
				scenario.condition,
				scenario.columnNumber,
			)
		} else {
			response, err = CDPDebuggerSetBreakpointByUrl(
				scenario.lineNumber,
				scenario.urlPattern, // 使用url
				"", // urlRegex参数为空
				"", // scriptHash
				scenario.condition,
				scenario.columnNumber,
			)
		}

		if err != nil {
			log.Printf("设置URL断点失败: %v", err)
			continue
		}

		// 解析结果
		result, err := ParseSetBreakpointByUrl(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 显示断点信息
		displayUrlBreakpointInfo(result, scenario.name, scenario.urlPattern, scenario.isRegex)
	}
}

func displayUrlBreakpointInfo(result *SetBreakpointByUrlResult, breakpointName, pattern string, isRegex bool) {
	log.Printf("URL断点设置成功:")
	log.Printf("  断点ID: %s", result.BreakpointID)

	patternType := "URL"
	if isRegex {
		patternType = "URL正则表达式"
	}
	log.Printf("  匹配模式: %s (%s)", pattern, patternType)

	if len(result.Locations) > 0 {
		log.Printf("  匹配位置数: %d", len(result.Locations))

		// 显示前几个匹配位置
		showCount := min(3, len(result.Locations))
		for i := 0; i < showCount; i++ {
			loc := result.Locations[i]
			log.Printf("  位置 %d: %s:%d:%d",
				i+1, loc.ScriptID, loc.LineNumber, loc.ColumnNumber)
		}

		if len(result.Locations) > showCount {
			log.Printf("  ... 还有 %d 个匹配位置", len(result.Locations)-showCount)
		}

		// 分析匹配结果
		analyzeUrlBreakpointMatches(result.Locations, pattern, isRegex)
	} else {
		log.Printf("  ⚠ 当前无匹配位置")
		log.Printf("    注意: 断点将在匹配的脚本加载时生效")
	}
}

func analyzeUrlBreakpointMatches(locations []Location, pattern string, isRegex bool) {
	log.Printf("  匹配分析:")

	if len(locations) == 0 {
		log.Printf("    ⚠ 暂无已加载的匹配脚本")
		log.Printf("      断点将在匹配脚本加载时自动设置")
		return
	}

	// 统计不同脚本
	scriptSet := make(map[string]bool)
	for _, loc := range locations {
		scriptSet[loc.ScriptID] = true
	}

	log.Printf("    匹配脚本数: %d", len(scriptSet))
	log.Printf("    总匹配位置数: %d", len(locations))

	// 行号分布
	lineNumbers := make([]int, len(locations))
	for i, loc := range locations {
		lineNumbers[i] = loc.LineNumber
	}

	sort.Ints(lineNumbers)
	minLine := lineNumbers[0]
	maxLine := lineNumbers[len(lineNumbers)-1]

	log.Printf("    行号范围: %d - %d", minLine, maxLine)

	// 匹配效果评估
	matchRatio := float64(len(locations)) / float64(len(scriptSet))
	log.Printf("    平均每脚本匹配位置: %.1f", matchRatio)

	if matchRatio > 5 {
		log.Printf("    ⚠ 匹配较密集，可能影响性能")
	} else if matchRatio > 1 {
		log.Printf("    ⚠ 匹配适中")
	} else {
		log.Printf("    ✅ 匹配精确")
	}
}

// 示例2: 正则表达式URL模式匹配
func exampleRegexUrlPatternMatching() {
	// === 应用场景描述 ===
	// 场景: 正则表达式URL模式匹配
	// 用途: 使用正则表达式匹配多个URL设置断点
	// 优势: 可以批量在相关文件上设置断点
	// 典型工作流: 设计正则模式 -> 测试匹配 -> 设置断点 -> 验证效果

	log.Println("正则表达式URL模式匹配示例...")

	// 定义正则表达式断点模式
	regexPatterns := []struct {
		name        string
		regex       string
		description string
		lineNumber  int
		columnNumber int
		condition   string
		expectedMatches []string
		complexity  string
	}{
		{
			name:        "版本化文件匹配",
			regex:       `.*\/app-v\d+\.\d+\.\d+\.js$`,
			description: "匹配版本化的应用脚本文件",
			lineNumber:  10,
			columnNumber: 0,
			condition:   "",
			expectedMatches: []string{
				"https://example.com/js/app-v1.0.0.js",
				"https://example.com/js/app-v1.2.3.js",
				"https://example.com/js/app-v2.0.1.js",
			},
			complexity:  "低",
		},
		{
			name:        "模块化组件",
			regex:       `.*\/components/.*\.js$`,
			description: "匹配所有组件目录下的JS文件",
			lineNumber:  5,
			columnNumber: 0,
			condition:   "typeof component !== 'undefined'",
			expectedMatches: []string{
				"https://example.com/components/button.js",
				"https://example.com/components/modal.js",
				"https://example.com/components/form.js",
			},
			complexity:  "中",
		},
		{
			name:        "API端点调试",
			regex:       `.*\/api/(v1|v2)/.*\.js$`,
			description: "匹配API v1和v2版本的所有脚本",
			lineNumber:  20,
			columnNumber: 4,
			condition:   "requestCount > 0",
			expectedMatches: []string{
				"https://api.example.com/v1/users.js",
				"https://api.example.com/v2/products.js",
				"https://api.example.com/v1/orders.js",
			},
			complexity:  "中",
		},
		{
			name:        "第三方库特定版本",
			regex:       `.*\/(jquery|lodash|moment)-\d+\.\d+\.\d+\.min\.js$`,
			description: "匹配常见第三方库的版本化压缩文件",
			lineNumber:  1,
			columnNumber: 0,
			condition:   "",
			expectedMatches: []string{
				"https://cdn.example.com/jquery-3.6.0.min.js",
				"https://cdn.example.com/lodash-4.17.21.min.js",
				"https://cdn.example.com/moment-2.29.1.min.js",
			},
			complexity:  "高",
		},
		{
			name:        "开发环境文件",
			regex:       `.*\.dev\.js$|.*\/dev/.*\.js$`,
			description: "匹配开发环境的JS文件",
			lineNumber:  15,
			columnNumber: 8,
			condition:   "process.env.NODE_ENV === 'development'",
			expectedMatches: []string{
				"https://dev.example.com/app.dev.js",
				"https://example.com/dev/utils.js",
				"https://example.com/src/dev/helpers.js",
			},
			complexity:  "中",
		},
		{
			name:        "带哈希的资源文件",
			regex:       `.*\.[a-f0-9]{8,}\.js$`,
			description: "匹配带内容哈希的资源文件",
			lineNumber:  0,
			columnNumber: 0,
			condition:   "",
			expectedMatches: []string{
				"https://example.com/app.abc12345.js",
				"https://example.com/vendor.def67890.js",
				"https://example.com/runtime.12345678.js",
			},
			complexity:  "高",
		},
		{
			name:        "特定域名文件",
			regex:       `https://(staging|test)\.example\.com/.*\.js$`,
			description: "匹配测试和预发布环境的文件",
			lineNumber:  8,
			columnNumber: 0,
			condition:   "window.location.hostname.includes('test')",
			expectedMatches: []string{
				"https://staging.example.com/app.js",
				"https://test.example.com/utils.js",
				"https://staging.example.com/components.js",
			},
			complexity:  "中",
		},
		{
			name:        "动态生成脚本",
			regex:       `.*\?.*\.js$`,
			description: "匹配带查询参数的动态脚本",
			lineNumber:  3,
			columnNumber: 0,
			condition:   "typeof dynamic === 'true'",
			expectedMatches: []string{
				"https://example.com/script.js?version=1.0",
				"https://example.com/bundle.js?v=abc123",
				"https://example.com/module.js?timestamp=1234567890",
			},
			complexity:  "低",
		},
	}

	// 测试每个正则模式
	for i, pattern := range regexPatterns {
		log.Printf("\n=== 正则模式 %d/%d: %s ===", i+1, len(regexPatterns), pattern.name)
		log.Printf("描述: %s", pattern.description)
		log.Printf("正则表达式: %s", pattern.regex)
		log.Printf("复杂度: %s", pattern.complexity)
		log.Printf("期望匹配: %v", pattern.expectedMatches)
		log.Printf("位置: 行 %d, 列 %d", pattern.lineNumber, pattern.columnNumber)
		if pattern.condition != "" {
			log.Printf("条件: %s", pattern.condition)
		}

		// 分析正则表达式
		analyzeRegexPattern(pattern.regex, pattern.complexity)

		// 设置正则模式断点
		response, err := CDPDebuggerSetBreakpointByUrl(
			pattern.lineNumber,
			"", // url参数为空
			pattern.regex, // 使用urlRegex
			"", // scriptHash
			pattern.condition,
			pattern.columnNumber,
		)

		if err != nil {
			log.Printf("设置正则断点失败: %v", err)
			continue
		}

		// 解析结果
		result, err := ParseSetBreakpointByUrl(response)
		if err != nil {
			log.Printf("解析结果失败: %v", err)
			continue
		}

		// 分析匹配结果
		analyzeRegexBreakpointResult(result, pattern.expectedMatches, pattern.regex)

		// 短暂延迟
		time.Sleep(100 * time.Millisecond)
	}
}

func analyzeRegexPattern(regex, complexity string) {
	log.Printf("正则表达式分析:")

	// 检查特殊字符
	specialChars := []string{".", "*", "+", "?", "|", "[", "]", "(", ")", "{", "}", "^", "$", "\\"}
	charCounts := make(map[string]int)

	for _, char := range specialChars {
		charCounts[char] = strings.Count(regex, char)
	}

	log.Printf("  特殊字符统计:")
	totalSpecialChars := 0
	for char, count := range charCounts {
		if count > 0 {
			log.Printf("    %s: %d", char, count)
			totalSpecialChars += count
		}
	}

	log.Printf("  总特殊字符数: %d", totalSpecialChars)

	// 复杂度验证
	actualComplexity := "低"
	if totalSpecialChars > 10 {
		actualComplexity = "高"
	} else if totalSpecialChars > 5 {
		actualComplexity = "中"
	}

	if actualComplexity != complexity {
		log.Printf("  ⚠ 复杂度评估不一致: 预期 %s, 实际 %s",
			complexity, actualComplexity)
	} else {
		log.Printf("  ✅ 复杂度评估一致: %s", complexity)
	}

	// 性能提示
	if strings.Contains(regex, ".*") && strings.Count(regex, ".*") > 3 {
		log.Printf("  ⚠ 包含多个'.*'，可能影响匹配性能")
	}
	if strings.Contains(regex, "|") && strings.Count(regex, "|") > 2 {
		log.Printf("  ⚠ 包含多个'|'选择符，可能影响匹配性能")
	}
}

func analyzeRegexBreakpointResult(result *SetBreakpointByUrlResult, expectedMatches []string, regex string) {
	log.Printf("正则断点结果分析:")
	log.Printf("  断点ID: %s", result.BreakpointID)
	log.Printf("  匹配位置数: %d", len(result.Locations))

	// 分组统计
	locationsByScript := make(map[string][]Location)
	for _, loc := range result.Locations {
		locationsByScript[loc.ScriptID] = append(locationsByScript[loc.ScriptID], loc)
	}

	log.Printf("  匹配脚本数: %d", len(locationsByScript))

	// 显示匹配统计
	if len(locationsByScript) > 0 {
		log.Printf("  各脚本匹配位置数:")
		count := 0
		for scriptID, locs := range locationsByScript {
			if count < 3 { // 只显示前3个
				log.Printf("    %s: %d 个位置",
					truncateString(scriptID, 30), len(locs))
				count++
			}
		}
		if len(locationsByScript) > 3 {
			log.Printf("    ... 还有 %d 个脚本", len(locationsByScript)-3)
		}
	}

	// 与期望匹配对比
	if len(expectedMatches) > 0 {
		log.Printf("  期望匹配数: %d", len(expectedMatches))

		// 简单模拟匹配检查
		matchedCount := 0
		for _, expected := range expectedMatches {
			// 简化检查，实际应该使用正则引擎
			if strings.Contains(regex, "/") && strings.Contains(expected, "/") {
				matchedCount++
			}
		}

		matchRate := float64(matchedCount) / float64(len(expectedMatches)) * 100
		log.Printf("  模拟匹配率: %.1f%% (%d/%d)",
			matchRate, matchedCount, len(expectedMatches))
	}

	// 效果评估
	if len(result.Locations) == 0 {
		log.Printf("  ⚠ 暂无匹配，但断点已设置")
		log.Printf("    新加载的匹配脚本将自动设置断点")
	} else if len(result.Locations) > 20 {
		log.Printf("  🔥 匹配位置较多，可能影响性能")
		log.Printf("    建议: 优化正则表达式或减少匹配范围")
	} else if len(result.Locations) > 5 {
		log.Printf("  ⚠ 匹配位置适中")
		log.Printf("    适合批量调试相关文件")
	} else {
		log.Printf("  ✅ 匹配精确")
		log.Printf("    适合特定文件的精确调试")
	}
}


*/

// -----------------------------------------------  Debugger.setBreakpointsActive  -----------------------------------------------
// === 应用场景 ===
// 1. 全局断点控制: 一次性启用或禁用所有断点
// 2. 调试流程控制: 在调试过程中临时禁用断点
// 3. 性能测试: 在性能测试时禁用断点避免干扰
// 4. 条件调试: 根据需要动态切换断点状态
// 5. 批量操作: 批量管理多个断点的激活状态
// 6. 自动化测试: 在自动化测试中控制断点行为

// CDPDebuggerSetBreakpointsActive 设置断点激活状态
// 参数说明:
//   - active: 是否激活所有断点
func CDPDebuggerSetBreakpointsActive(active bool) (string, error) {
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
		"method": "Debugger.setBreakpointsActive",
		"params": {
			"active": %v
		}
	}`, reqID, active)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setBreakpointsActive 请求失败: %w", err)
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
			return "", fmt.Errorf("setBreakpointsActive 请求超时")
		}
	}
}

/*

// 示例1: 基本断点状态控制
func exampleBasicBreakpointStateControl() {
	// === 应用场景描述 ===
	// 场景: 基本断点状态控制
	// 用途: 控制所有断点的全局激活状态
	// 优势: 可以一键启用或禁用所有调试断点
	// 典型工作流: 设置断点 -> 控制状态 -> 执行测试 -> 分析结果

	log.Println("基本断点状态控制示例...")

	// 定义断点状态控制场景
	stateControlScenarios := []struct {
		name        string
		description string
		active      bool
		useCase     string
		expectedEffect string
	}{
		{
			name:        "启用所有断点",
			description: "激活页面上所有已设置的断点",
			active:      true,
			useCase:     "开始调试会话",
			expectedEffect: "所有断点变为活动状态，会在命中时暂停执行",
		},
		{
			name:        "禁用所有断点",
			description: "停用页面上所有已设置的断点",
			active:      false,
			useCase:     "临时跳过调试检查",
			expectedEffect: "所有断点变为非活动状态，不会暂停执行",
		},
		{
			name:        "批量调试切换",
			description: "在批量测试中切换断点状态",
			active:      true,
			useCase:     "自动化测试流程",
			expectedEffect: "在测试阶段启用断点，在性能测试阶段禁用",
		},
		{
			name:        "条件调试控制",
			description: "根据条件动态控制断点",
			active:      false,
			useCase:     "条件调试流程",
			expectedEffect: "在特定条件下启用断点，其他条件下禁用",
		},
		{
			name:        "性能测试准备",
			description: "在性能测试前禁用断点",
			active:      false,
			useCase:     "性能基准测试",
			expectedEffect: "消除断点对性能测量的干扰",
		},
		{
			name:        "生产环境调试",
			description: "在生产环境谨慎启用断点",
			active:      true,
			useCase:     "生产问题诊断",
			expectedEffect: "在受控条件下启用断点进行问题诊断",
		},
		{
			name:        "教学演示控制",
			description: "在教学演示中控制断点状态",
			active:      true,
			useCase:     "代码教学演示",
			expectedEffect: "在讲解时启用断点，在自由练习时禁用",
		},
		{
			name:        "渐进式调试",
			description: "渐进式启用断点进行调试",
			active:      true,
			useCase:     "复杂问题调试",
			expectedEffect: "逐步启用断点，缩小问题范围",
		},
	}

	// 测试不同的状态控制场景
	for i, scenario := range stateControlScenarios {
		log.Printf("\n=== 状态控制场景 %d/%d: %s ===", i+1, len(stateControlScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("使用场景: %s", scenario.useCase)
		log.Printf("激活状态: %v", scenario.active)
		log.Printf("预期效果: %s", scenario.expectedEffect)

		// 模拟当前断点状态
		currentState := simulateCurrentBreakpointState(scenario.active)
		log.Printf("当前断点状态:")
		for key, value := range currentState {
			log.Printf("  %s: %v", key, value)
		}

		// 设置断点激活状态
		log.Printf("设置断点激活状态为 %v...", scenario.active)

		response, err := CDPDebuggerSetBreakpointsActive(scenario.active)
		if err != nil {
			log.Printf("设置断点状态失败: %v", err)
			continue
		}

		log.Printf("设置成功: %s", response)

		// 模拟状态变化后的效果
		simulateBreakpointStateEffect(scenario.active, scenario.expectedEffect)

		// 分析状态变化
		analyzeBreakpointStateChange(scenario.active, currentState)

		// 短暂延迟
		time.Sleep(100 * time.Millisecond)
	}
}

func simulateCurrentBreakpointState(targetActive bool) map[string]interface{} {
	state := make(map[string]interface{})

	// 模拟断点信息
	state["totalBreakpoints"] = 8
	state["conditionalBreakpoints"] = 3
	state["activeBeforeChange"] = !targetActive // 假设之前状态相反
	state["lastStateChange"] = time.Now().Add(-5 * time.Minute).Format("15:04:05")

	// 模拟断点分布
	state["breakpointsByFile"] = map[string]int{
		"app.js": 3,
		"utils.js": 2,
		"api.js": 2,
		"ui.js": 1,
	}

	// 模拟性能影响
	if targetActive {
		state["estimatedPerformanceImpact"] = "中 (断点启用)"
	} else {
		state["estimatedPerformanceImpact"] = "低 (断点禁用)"
	}

	return state
}

func simulateBreakpointStateEffect(active bool, expectedEffect string) {
	log.Printf("状态变化效果模拟:")

	if active {
		log.Printf("  ✅ 所有断点已激活")
		log.Printf("    效果: %s", expectedEffect)

		// 模拟激活后的行为
		log.Printf("  模拟断点行为:")
		log.Printf("    - 代码执行会在断点处暂停")
		log.Printf("    - 可以检查变量状态")
		log.Printf("    - 可以单步执行代码")
		log.Printf("    - 条件断点会在条件满足时触发")
	} else {
		log.Printf("  ⚠ 所有断点已禁用")
		log.Printf("    效果: %s", expectedEffect)

		// 模拟禁用后的行为
		log.Printf("  模拟断点行为:")
		log.Printf("    - 代码执行不会在断点处暂停")
		log.Printf("    - 断点仍然存在但不会触发")
		log.Printf("    - 代码执行流畅无中断")
		log.Printf("    - 适合性能测试和正常流程")
	}

	// 性能影响
	if active {
		log.Printf("  ⚠ 性能影响: 启用断点会增加执行开销")
		log.Printf("    建议: 仅在调试时启用")
	} else {
		log.Printf("  ✅ 性能影响: 禁用断点最小化执行开销")
		log.Printf("    建议: 在生产环境和性能测试时禁用")
	}
}

func analyzeBreakpointStateChange(newActive bool, previousState map[string]interface{}) {
	log.Printf("状态变化分析:")

	oldActive, _ := previousState["activeBeforeChange"].(bool)
	totalBreakpoints, _ := previousState["totalBreakpoints"].(int)
	conditionalBreakpoints, _ := previousState["conditionalBreakpoints"].(int)

	log.Printf("  状态变化: %v -> %v", oldActive, newActive)
	log.Printf("  总断点数: %d", totalBreakpoints)
	log.Printf("  条件断点数: %d", conditionalBreakpoints)

	if newActive && !oldActive {
		log.Printf("  ⚡ 状态: 从禁用切换到启用")
		log.Printf("    影响: %d 个断点变为活动状态", totalBreakpoints)
		if conditionalBreakpoints > 0 {
			log.Printf("    其中 %d 个条件断点会在条件满足时触发", conditionalBreakpoints)
		}
	} else if !newActive && oldActive {
		log.Printf("  ⚡ 状态: 从启用切换到禁用")
		log.Printf("    影响: %d 个断点变为非活动状态", totalBreakpoints)
		log.Printf("    代码执行将不会在断点处暂停")
	} else {
		log.Printf("  ℹ 状态: 保持不变")
		log.Printf("    断点状态未发生变化")
	}

	// 建议
	log.Printf("  建议:")
	if newActive {
		log.Printf("    - 现在可以开始调试")
		log.Printf("    - 条件断点会在条件满足时暂停")
		log.Printf("    - 注意性能影响")
	} else {
		log.Printf("    - 适合运行测试和性能基准")
		log.Printf("    - 代码执行不会被打断")
		log.Printf("    - 断点配置仍然保留")
	}
}


*/

// -----------------------------------------------  Debugger.setInstrumentationBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 脚本执行跟踪: 在脚本开始执行时暂停调试
// 2. 源映射调试: 在带源映射的脚本执行前暂停
// 3. 动态脚本分析: 分析动态加载的脚本
// 4. 调试器初始化: 在脚本执行前设置调试环境
// 5. 性能分析: 跟踪脚本执行时机
// 6. 安全检测: 监控可疑脚本执行

// CDPDebuggerSetInstrumentationBreakpoint 设置检测断点
// 参数说明:
//   - instrumentation: 检测类型，可选值：
//   - "beforeScriptExecution": 在脚本执行前
//   - "beforeScriptWithSourceMapExecution": 在带有源映射的脚本执行前
func CDPDebuggerSetInstrumentationBreakpoint(instrumentation string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	validInstrumentations := map[string]bool{
		"beforeScriptExecution":              true,
		"beforeScriptWithSourceMapExecution": true,
	}

	if !validInstrumentations[instrumentation] {
		return "", fmt.Errorf("无效的检测类型: %s，有效值为: beforeScriptExecution, beforeScriptWithSourceMapExecution", instrumentation)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setInstrumentationBreakpoint",
		"params": {
			"instrumentation": "%s"
		}
	}`, reqID, instrumentation)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setInstrumentationBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("setInstrumentationBreakpoint 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.setPauseOnExceptions  -----------------------------------------------
// === 应用场景 ===
// 1. 异常调试: 在异常发生时暂停调试，检查异常上下文
// 2. 错误处理: 调试catch块中的错误处理逻辑
// 3. 异常监控: 监控应用程序中的所有异常
// 4. 调试配置: 根据调试需求配置异常暂停行为
// 5. 自动化测试: 在测试中自动暂停在异常处
// 6. 性能分析: 分析异常对性能的影响

// CDPDebuggerSetPauseOnExceptions 设置异常暂停状态
// 参数说明:
//   - state: 异常暂停状态，可选值：
//   - "none": 不在任何异常上暂停
//   - "caught": 在捕获的异常上暂停
//   - "uncaught": 在未捕获的异常上暂停
//   - "all": 在所有异常上暂停
func CDPDebuggerSetPauseOnExceptions(state string) error {
	if !DefaultBrowserWS() {
		return fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	validStates := map[string]bool{
		"none":     true,
		"caught":   true,
		"uncaught": true,
		"all":      true,
	}

	if !validStates[state] {
		return fmt.Errorf("无效的异常暂停状态: %s，有效值为: none, caught, uncaught, all", state)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setPauseOnExceptions",
		"params": {
			"state": "%s"
		}
	}`, reqID, state)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("发送 setPauseOnExceptions 请求失败: %w", err)
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
				return fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return fmt.Errorf("CDP错误: %v", errorObj)
				}

				return nil
			}

		case <-timer.C:
			return fmt.Errorf("setPauseOnExceptions 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.setScriptSource  -----------------------------------------------
// === 应用场景 ===
// 1. 实时代码编辑: 在调试过程中修改代码
// 2. 热修复测试: 测试代码修改而不实际应用
// 3. 代码实验: 尝试不同的代码实现
// 4. 调试增强: 修改代码添加调试信息
// 5. 自动化修复: 自动修复发现的问题
// 6. 代码热重载: 运行时更新代码逻辑

// CDPDebuggerSetScriptSource 设置脚本源
// 参数说明:
//   - scriptId: 脚本ID
//   - scriptSource: 脚本新内容
//   - dryRun: 是否干运行（不实际应用更改）
func CDPDebuggerSetScriptSource(scriptId, scriptSource string, dryRun bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if scriptId == "" {
		return "", fmt.Errorf("脚本ID不能为空")
	}
	if scriptSource == "" {
		return "", fmt.Errorf("脚本内容不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 转义脚本源中的特殊字符
	escapedScriptSource := strings.ReplaceAll(scriptSource, `"`, `\"`)
	escapedScriptSource = strings.ReplaceAll(escapedScriptSource, "\n", "\\n")
	escapedScriptSource = strings.ReplaceAll(escapedScriptSource, "\t", "\\t")
	escapedScriptSource = strings.ReplaceAll(escapedScriptSource, "\r", "\\r")

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setScriptSource",
		"params": {
			"scriptId": "%s",
			"scriptSource": "%s",
			"dryRun": %v
		}
	}`, reqID, scriptId, escapedScriptSource, dryRun)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setScriptSource 请求失败: %w", err)
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
			return "", fmt.Errorf("setScriptSource 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.setSkipAllPauses  -----------------------------------------------
// === 应用场景 ===
// 1. 批量调试控制: 在批量测试中临时跳过所有暂停
// 2. 性能测试: 在性能测试时避免调试暂停干扰
// 3. 自动化执行: 在自动化脚本执行时跳过暂停
// 4. 快速运行: 需要快速执行代码而不被调试器中断
// 5. 生产调试: 在生产环境谨慎控制调试行为
// 6. 用户体验: 在用户可见的场景跳过调试中断

// CDPDebuggerSetSkipAllPauses 设置跳过所有暂停
// 参数说明:
//   - skip: 是否跳过所有暂停
func CDPDebuggerSetSkipAllPauses(skip bool) (string, error) {
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
		"method": "Debugger.setSkipAllPauses",
		"params": {
			"skip": %v
		}
	}`, reqID, skip)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setSkipAllPauses 请求失败: %w", err)
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
			return "", fmt.Errorf("setSkipAllPauses 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.setVariableValue  -----------------------------------------------
// === 应用场景 ===
// 1. 动态变量修改: 在调试过程中修改变量的值
// 2. 测试用例调试: 修改变量值测试不同代码路径
// 3. 错误重现: 修改变量值重现特定的错误场景
// 4. 条件测试: 修改变量值测试条件分支
// 5. 状态注入: 注入特定的状态值进行调试
// 6. 教学演示: 实时展示变量值变化的效果

// CDPDebuggerSetVariableValue 设置变量值
// 参数说明:
//   - scopeNumber: 作用域编号（0-based）
//   - variableName: 变量名
//   - newValue: 新的变量值（可以是任意JSON可序列化的值）
//   - callFrameId: 调用帧ID
func CDPDebuggerSetVariableValue(scopeNumber int, variableName string, newValue interface{}, callFrameId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if scopeNumber < 0 {
		return "", fmt.Errorf("作用域编号不能为负数")
	}
	if variableName == "" {
		return "", fmt.Errorf("变量名不能为空")
	}
	if callFrameId == "" {
		return "", fmt.Errorf("调用帧ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建newValue参数
	var newValueParam interface{}

	// 尝试将newValue转换为适合JSON序列化的格式
	switch v := newValue.(type) {
	case nil:
		newValueParam = map[string]interface{}{"value": nil}
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		// 基本类型，使用value字段
		newValueParam = map[string]interface{}{"value": v}
	case map[string]interface{}, []interface{}:
		// 对象或数组，使用value字段
		newValueParam = map[string]interface{}{"value": v}
	default:
		// 其他类型尝试序列化
		newValueParam = map[string]interface{}{"value": v}
	}

	// 序列化参数
	params := map[string]interface{}{
		"scopeNumber":  scopeNumber,
		"variableName": variableName,
		"newValue":     newValueParam,
		"callFrameId":  callFrameId,
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setVariableValue",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setVariableValue 请求失败: %w", err)
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
			return "", fmt.Errorf("setVariableValue 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.stepInto  -----------------------------------------------
// === 应用场景 ===
// 1. 函数调用调试: 进入函数内部进行逐行调试
// 2. 嵌套函数分析: 调试嵌套的函数调用
// 3. 方法调用跟踪: 跟踪对象方法的执行
// 4. 回调函数调试: 进入回调函数内部调试
// 5. 异步函数调试: 配合breakOnAsyncCall调试异步函数
// 6. 递归调用分析: 分析递归函数的每一层调用

// CDPDebuggerStepInto 单步进入函数
// 参数说明:
//   - breakOnAsyncCall: 控制是否在异步调用时中断
func CDPDebuggerStepInto(breakOnAsyncCall map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := make(map[string]interface{})
	if len(breakOnAsyncCall) > 0 {
		params["breakOnAsyncCall"] = breakOnAsyncCall
	}

	var paramsJSON []byte
	var err error
	if len(params) > 0 {
		paramsJSON, err = json.Marshal(params)
		if err != nil {
			return "", fmt.Errorf("序列化参数失败: %w", err)
		}
	}

	// 构建消息
	message := ""
	if len(params) > 0 {
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "Debugger.stepInto",
			"params": %s
		}`, reqID, string(paramsJSON))
	} else {
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "Debugger.stepInto"
		}`, reqID)
	}

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stepInto 请求失败: %w", err)
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
			return "", fmt.Errorf("stepInto 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.stepOut  -----------------------------------------------
// === 应用场景 ===
// 1. 快速退出函数: 当函数内部调试完成后快速跳出
// 2. 跳过复杂逻辑: 跳过函数内部的复杂代码段
// 3. 调试流程控制: 控制调试流程，快速返回到调用点
// 4. 错误排查: 快速退出当前函数，检查上层调用栈
// 5. 性能优化: 跳过不需要深入调试的函数
// 6. 教学演示: 快速展示函数调用返回的结果

// CDPDebuggerStepOut 单步跳出函数
func CDPDebuggerStepOut() (string, error) {
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
		"method": "Debugger.stepOut"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stepOut 请求失败: %w", err)
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
			return "", fmt.Errorf("stepOut 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.stepOver  -----------------------------------------------
// === 应用场景 ===
// 1. 逐行调试: 在当前作用域内逐行执行代码
// 2. 跳过函数调用: 跳过不需要进入的函数调用
// 3. 循环调试: 在循环中逐次执行，不进入内部函数
// 4. 条件语句调试: 调试条件分支的执行流程
// 5. 表达式求值: 观察复杂表达式的逐步求值过程
// 6. 性能关键代码: 调试性能敏感代码而不进入内部函数

// CDPDebuggerStepOver 单步跳过语句
func CDPDebuggerStepOver() (string, error) {
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
		"method": "Debugger.stepOver"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stepOver 请求失败: %w", err)
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
			return "", fmt.Errorf("stepOver 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.disassembleWasmModule  -----------------------------------------------
// === 应用场景 ===
// 1. Wasm代码分析: 分析WebAssembly模块的字节码
// 2. 性能优化: 优化Wasm模块的性能
// 3. 安全审计: 检查Wasm模块的安全性
// 4. 调试支持: 为Wasm模块提供调试信息
// 5. 教育研究: 学习WebAssembly的内部工作原理
// 6. 逆向工程: 分析第三方Wasm模块

// CDPDebuggerDisassembleWasmModule 反汇编Wasm模块
// 参数说明:
//   - scriptId: 脚本ID
func CDPDebuggerDisassembleWasmModule(scriptId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if scriptId == "" {
		return "", fmt.Errorf("脚本ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.disassembleWasmModule",
		"params": {
			"scriptId": "%s"
		}
	}`, reqID, scriptId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 disassembleWasmModule 请求失败: %w", err)
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
			return "", fmt.Errorf("disassembleWasmModule 请求超时")
		}
	}
}

// -----------------------------------------------  Debugger.getStackTrace  -----------------------------------------------
// === 应用场景 ===
// 1. 异步调试: 获取异步操作的调用栈跟踪
// 2. 错误分析: 分析异常发生时的完整调用栈
// 3. 性能分析: 分析函数调用关系和执行路径
// 4. 监控诊断: 获取特定操作的完整调用历史
// 5. 教学演示: 展示调用栈的结构和层次
// 6. 代码审计: 分析代码执行流程

// CDPDebuggerGetStackTrace 获取调用栈跟踪
// 参数说明:
//   - stackTraceId: 调用栈ID
func CDPDebuggerGetStackTrace(stackTraceId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if stackTraceId == "" {
		return "", fmt.Errorf("调用栈ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.getStackTrace",
		"params": {
			"stackTraceId": "%s"
		}
	}`, reqID, stackTraceId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getStackTrace 请求失败: %w", err)
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
			return "", fmt.Errorf("getStackTrace 请求超时")
		}
	}
}
