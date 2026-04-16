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
func CDPDebuggerEvaluateOnCallFrame(params string) (string, error) {
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
		"method": "Debugger.evaluateOnCallFrame",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
func CDPDebuggerGetPossibleBreakpoints(params string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.getPossibleBreakpoints",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// 示例: 代码审计和安全搜索
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
		name      string
		scriptId  string
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
func CDPDebuggerSetBreakpoint(params string) (string, error) {
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
		"method": "Debugger.setBreakpoint",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
func CDPDebuggerSetBreakpointByUrl(params string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setBreakpointByUrl",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
func CDPDebuggerSetInstrumentationBreakpoint(params string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setInstrumentationBreakpoint",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setInstrumentationBreakpoint 请求失败: %w", err)
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
func CDPDebuggerSetScriptSource(params string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setScriptSource",
		"params": %s
	}`, reqID, params)
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
func CDPDebuggerSetVariableValue(params string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Debugger.setVariableValue",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
func CDPDebuggerStepInto(params string) (string, error) {
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
			"method": "Debugger.stepInto",
			"params": %s
		}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
