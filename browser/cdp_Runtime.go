package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Runtime.addBinding  -----------------------------------------------
// === 应用场景 ===
// 1. 网页与Go程序通信: 让前端JS调用Go后端方法，实现双向数据交互
// 2. 自动化注入API: 向页面注入自定义JS函数，供页面代码主动触发Go逻辑
// 3. 日志/数据上报: 网页JS将日志、埋点数据发送给Go程序处理
// 4. 页面操作控制: 前端通过绑定函数通知Go执行浏览器控制、文件操作等权限操作
// 5. 跨环境数据传递: 在浏览器渲染层与自动化控制层之间传递结构化数据
// 6. 自定义调试工具: 注入调试函数，让页面主动上报状态给Go调试器

// CDPRuntimeAddBinding 向页面注入绑定函数，使JS可以调用Go侧的处理逻辑
func CDPRuntimeAddBinding(bindingName string) (string, error) {
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
		"method": "Runtime.addBinding",
		"params": {
			"name": "%s"
		}
	}`, reqID, bindingName)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.addBinding 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.addBinding 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础注入绑定函数 ====================
func ExampleAddBinding_Basic() {
	// 注入名为 "sendLogToGo" 的绑定函数
	resp, err := CDPRuntimeAddBinding("sendLogToGo")
	if err != nil {
		log.Fatalf("注入绑定失败: %v", err)
	}
	log.Println("注入成功，响应:", resp)

	// 页面JS可直接调用：sendLogToGo("页面日志内容")
}

// ==================== 使用示例 2：注入多个绑定函数 ====================
func ExampleAddBinding_Multiple() {
	// 注入数据上报函数
	CDPRuntimeAddBinding("reportData")
	// 注入页面关闭通知函数
	CDPRuntimeAddBinding("onPageClosed")
	// 注入文件选择通知函数
	CDPRuntimeAddBinding("selectFile")

	log.Println("所有绑定函数注入完成")
}

// ==================== 使用示例 3：自动化测试中注入通信API ====================
func ExampleAddBinding_AutoTest() {
	// 测试场景：让页面JS主动通知Go测试结果
	_, err := CDPRuntimeAddBinding("testResultCallback")
	if err != nil {
		log.Fatal(err)
	}

	// JS调用示例：testResultCallback({caseId: 101, status: "success"})
}

*/

// -----------------------------------------------  Runtime.awaitPromise  -----------------------------------------------
// === 应用场景 ===
// 1. 异步结果获取: 等待页面Promise执行完成，获取最终结果
// 2. 自动化测试: 等待页面异步接口、异步渲染完成后再执行断言
// 3. 数据抓取: 等待页面异步加载的数据返回后提取内容
// 4. 异步操作同步化: 将页面异步Promise转为同步等待逻辑
// 5. 错误捕获: 精准捕获Promise执行失败的原因与堆栈
// 6. 前端调试: 调试页面异步逻辑，等待Promise完成查看结果

// CDPRuntimeAwaitPromise 等待页面中的Promise对象执行完成并返回结果
func CDPRuntimeAwaitPromise(promiseObjectID string, returnByValue bool) (string, error) {
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
		"method": "Runtime.awaitPromise",
		"params": {
			"promiseObjectId": "%s",
			"returnByValue": %t
		}
	}`, reqID, promiseObjectID, returnByValue)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.awaitPromise 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 30 * time.Second // Promise异步操作延长超时时间
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.awaitPromise 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：等待Promise并返回完整值 ====================
func ExampleAwaitPromise_ReturnValue() {
	// promiseObjectId来自Runtime.evaluate执行Promise返回的对象ID
	promiseObjID := "123456"

	// 等待Promise完成，直接返回结果值
	resp, err := CDPRuntimeAwaitPromise(promiseObjID, true)
	if err != nil {
		log.Fatalf("等待Promise失败: %v", err)
	}
	log.Println("Promise执行完成，结果:", resp)
}

// ==================== 使用示例 2：等待Promise仅返回对象引用 ====================
func ExampleAwaitPromise_ObjectRef() {
	// 页面大型Promise对象（如大量列表数据），只获取引用不序列化值
	promiseObjID := "789012"

	resp, err := CDPRuntimeAwaitPromise(promiseObjID, false)
	if err != nil {
		log.Fatalf("等待Promise失败: %v", err)
	}
	log.Println("Promise对象引用返回:", resp)
}

// ==================== 使用示例 3：自动化测试等待异步接口 ====================
func ExampleAwaitPromise_AutoTest() {
	// 场景：等待页面fetch接口Promise完成
	// 先通过evaluate获取Promise对象ID
	evalResp, _ := CDPRuntimeEvaluate(`fetch('/api/data').then(res=>res.json())`, true, false)

	// 解析获取promiseObjectId
	promiseObjID := "解析evalResp得到的对象ID"

	// 等待接口请求完成
	result, err := CDPRuntimeAwaitPromise(promiseObjID, true)
	if err != nil {
		log.Fatal("接口请求超时:", err)
	}
	log.Println("异步接口数据返回:", result)
}

*/

// -----------------------------------------------  Runtime.callFunctionOn  -----------------------------------------------
// === 应用场景 ===
// 1. 精准调用页面函数: 在指定DOM对象/上下文上调用页面已有的JS函数
// 2. 自动化操作DOM: 对指定元素执行点击、赋值、获取属性等操作
// 3. 跨上下文执行: 在指定的JS执行上下文、iframe中调用函数
// 4. 带参数调用: 向页面函数传递自定义参数并获取返回值
// 5. 数据获取: 调用页面方法获取组件数据、状态、计算结果
// 6. 页面逻辑控制: 主动触发页面业务逻辑、事件处理函数

// CDPRuntimeCallFunctionOn 在指定对象上调用函数
// objectID: 要调用函数的目标对象ID (DOM元素/JS对象)
// functionDeclaration: 要执行的函数声明字符串
// args: 函数参数数组JSON字符串 (格式: [{"value":"xxx"},{"objectId":"yyy"}])
// returnByValue: 是否直接返回结果值(true)还是对象引用(false)
func CDPRuntimeCallFunctionOn(objectID, functionDeclaration, args string, returnByValue bool) (string, error) {
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
		"method": "Runtime.callFunctionOn",
		"params": {
			"objectId": "%s",
			"functionDeclaration": "%s",
			"arguments": %s,
			"returnByValue": %t,
			"awaitPromise": true
		}
	}`, reqID, objectID, functionDeclaration, args, returnByValue)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.callFunctionOn 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 15 * time.Second
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.callFunctionOn 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：调用DOM元素的点击方法 ====================
func ExampleCallFunctionOn_ClickElement() {
	// 目标DOM元素的objectID (通过DOM.getDocument等API获取)
	elemObjectID := "123456"

	// 调用元素的click()函数
	args := "[]" // 无参数
	resp, err := CDPRuntimeCallFunctionOn(elemObjectID, "function() { this.click(); }", args, false)
	if err != nil {
		log.Fatalf("元素点击失败: %v", err)
	}
	log.Println("元素点击执行成功:", resp)
}

// ==================== 使用示例 2：给输入框赋值并获取值 ====================
func ExampleCallFunctionOn_InputValue() {
	inputObjectID := "789012"

	// 调用函数给input赋值并返回值
	funcDecl := `function(val) { this.value = val; return this.value; }`
	// 传递参数
	args := `[{"value":"自动化测试内容"}]`

	resp, err := CDPRuntimeCallFunctionOn(inputObjectID, funcDecl, args, true)
	if err != nil {
		log.Fatalf("输入框赋值失败: %v", err)
	}
	log.Println("输入框操作结果:", resp)
}

// ==================== 使用示例 3：调用页面自定义JS函数并传参 ====================
func ExampleCallFunctionOn_PageFunction() {
	// 页面全局对象的objectID
	globalObjectID := "987654"

	// 调用页面已定义的函数: getUserInfo(1001)
	funcDecl := `function(userId) { return window.getUserInfo(userId); }`
	args := `[{"value":1001}]`

	// 等待Promise并返回值
	resp, err := CDPRuntimeCallFunctionOn(globalObjectID, funcDecl, args, true)
	if err != nil {
		log.Fatalf("调用页面函数失败: %v", err)
	}
	log.Println("用户信息返回结果:", resp)
}

*/

// -----------------------------------------------  Runtime.compileScript  -----------------------------------------------
// === 应用场景 ===
// 1. 预编译JS脚本: 提前编译JS代码，后续执行时提升运行速度
// 2. 自动化脚本缓存: 缓存高频执行的自动化脚本，避免重复编译
// 3. 安全校验: 先编译验证JS语法正确性，再执行脚本
// 4. 批量脚本管理: 编译多个脚本，按需执行，提升执行效率
// 5. 调试预处理: 调试前预编译脚本，快速定位语法错误
// 6. 隔离执行环境: 编译脚本到指定上下文，实现隔离运行

// CDPRuntimeCompileScript 编译JS脚本但不执行，返回脚本ID
// expression: 要编译的JS脚本内容
// sourceURL: 脚本来源URL（用于调试标识）
// persistScript: 是否持久化脚本（true: 长期保存; false: 临时使用）
// executionContextID: 执行上下文ID（0表示使用默认上下文）
func CDPRuntimeCompileScript(expression, sourceURL string, persistScript bool, executionContextID int) (string, error) {
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
		"method": "Runtime.compileScript",
		"params": {
			"expression": %s,
			"sourceURL": "%s",
			"persistScript": %t,
			"executionContextId": %d
		}
	}`, reqID, jsonEscape(expression), sourceURL, persistScript, executionContextID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.compileScript 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 10 * time.Second
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.compileScript 请求超时")
		}
	}
}

// jsonEscape JSON转义工具函数，防止脚本字符串破坏JSON格式
func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

/*

// ==================== 使用示例 1：基础编译临时脚本 ====================
func ExampleCompileScript_Basic() {
	// 要编译的JS脚本
	script := `console.log("预编译脚本执行"); return 100 + 200;`

	// 编译临时脚本（不持久化），使用默认上下文
	resp, err := CDPRuntimeCompileScript(script, "temp://testScript.js", false, 0)
	if err != nil {
		log.Fatalf("脚本编译失败: %v", err)
	}
	log.Println("脚本编译成功，返回脚本ID:", resp)
}

// ==================== 使用示例 2：编译持久化自动化脚本 ====================
func ExampleCompileScript_Persist() {
	// 自动化业务脚本
	autoScript := `
		function autoRun() {
			document.querySelector("button").click();
			return document.title;
		}
		autoRun();
	`

	// 编译并持久化脚本，长期复用
	resp, err := CDPRuntimeCompileScript(autoScript, "auto://task.js", true, 0)
	if err != nil {
		log.Fatalf("持久化脚本编译失败: %v", err)
	}
	log.Println("持久化脚本编译成功:", resp)
}

// ==================== 使用示例 3：指定执行上下文编译脚本 ====================
func ExampleCompileScript_Context() {
	// 自定义脚本
	script := `return "在指定上下文执行";`
	// 指定上下文ID（来自Runtime.createExecutionContext创建）
	contextID := 1001

	resp, err := CDPRuntimeCompileScript(script, "context://custom.js", false, contextID)
	if err != nil {
		log.Fatalf("上下文脚本编译失败: %v", err)
	}
	log.Println("指定上下文编译成功:", resp)
}

*/

// -----------------------------------------------  Runtime.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 关闭Runtime域事件监听: 停止接收执行上下文创建、销毁、绑定调用等事件
// 2. 资源释放: 关闭Runtime域后释放浏览器相关监听资源，降低性能消耗
// 3. 流程收尾: 自动化测试/调试完成后，关闭Runtime功能
// 4. 安全隔离: 关闭后禁止JS执行、函数调用等Runtime相关操作
// 5. 状态重置: 重置Runtime域状态，为后续重新启用做准备
// 6. 多流程切换: 切换不同CDP域操作时，关闭不需要的Runtime域

// CDPRuntimeDisable 禁用Runtime域，停止事件通知并释放相关资源
func CDPRuntimeDisable() (string, error) {
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
		"method": "Runtime.disable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.disable 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.disable 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础关闭Runtime域 ====================
func ExampleDisable_Basic() {
	// 关闭Runtime功能，停止所有事件监听
	resp, err := CDPRuntimeDisable()
	if err != nil {
		log.Fatalf("禁用Runtime域失败: %v", err)
	}
	log.Println("成功禁用Runtime域，响应:", resp)
}

// ==================== 使用示例 2：自动化测试完成后清理关闭 ====================
func ExampleDisable_AutoTest() {
	// 1. 执行测试逻辑
	log.Println("执行自动化测试...")
	// ... 测试代码 ...

	// 2. 测试完成后禁用Runtime域释放资源
	_, err := CDPRuntimeDisable()
	if err != nil {
		log.Fatalf("测试后清理Runtime失败: %v", err)
	}
	log.Println("测试完成，Runtime域已禁用")
}

// ==================== 使用示例 3：先关闭再重新启用Runtime ====================
func ExampleDisable_Reset() {
	// 关闭Runtime域重置状态
	CDPRuntimeDisable()
	log.Println("Runtime域已关闭，准备重置...")

	// 等待1秒后重新启用
	time.Sleep(1 * time.Second)
	// 对应启用方法：CDPRuntimeEnable()
	// CDPRuntimeEnable()
	log.Println("Runtime域已重置并重新启用")
}

*/

// -----------------------------------------------  Runtime.discardConsoleEntries  -----------------------------------------------
// === 应用场景 ===
// 1. 控制台日志清理: 清空浏览器控制台所有历史日志，避免日志堆积
// 2. 测试环境重置: 单条测试用例执行完毕后清空日志，避免干扰下一条用例
// 3. 性能优化: 长时间运行自动化任务时，定期清理控制台释放内存
// 4. 调试隔离: 调试特定功能前清空既有日志，只关注新产生的日志
// 5. 日志数据重置: 重新开始捕获控制台日志前的初始化操作
// 6. 界面整洁: 清空冗余日志，让控制台只展示关键运行信息

// CDPRuntimeDiscardConsoleEntries 清空浏览器控制台所有日志条目
func CDPRuntimeDiscardConsoleEntries() (string, error) {
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
		"method": "Runtime.discardConsoleEntries"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.discardConsoleEntries 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.discardConsoleEntries 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础清空控制台日志 ====================
func ExampleDiscardConsoleEntries_Basic() {
	// 一键清空所有控制台日志
	resp, err := CDPRuntimeDiscardConsoleEntries()
	if err != nil {
		log.Fatalf("清空控制台日志失败: %v", err)
	}
	log.Println("控制台已清空，响应:", resp)
}

// ==================== 使用示例 2：自动化测试用例间清理日志 ====================
func ExampleDiscardConsoleEntries_AutoTest() {
	// 执行第一个测试用例
	log.Println("执行测试用例1...")
	// 测试逻辑...

	// 清空控制台，避免日志影响下一个用例
	_, err := CDPRuntimeDiscardConsoleEntries()
	if err != nil {
		log.Fatalf("用例1清理日志失败: %v", err)
	}

	// 执行第二个测试用例
	log.Println("执行测试用例2...")
	// 测试逻辑...
	log.Println("所有测试用例执行完成，控制台无冗余日志")
}

// ==================== 使用示例 3：定时清理控制台（长时间运行任务） ====================
func ExampleDiscardConsoleEntries_Timer() {
	// 模拟长时间运行的自动化任务，每30秒清空一次控制台防止内存溢出
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := CDPRuntimeDiscardConsoleEntries()
		if err != nil {
			log.Printf("定时清空日志失败: %v", err)
			continue
		}
		log.Println("已定时清空控制台日志")
	}
}

*/

// -----------------------------------------------  Runtime.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 启用Runtime域: 开启JS执行、控制台、绑定调用等核心功能
// 2. 监听控制台事件: 开始接收console.log、warn、error等日志事件
// 3. 自动化初始化: 启动自动化测试/爬虫前必须启用Runtime域
// 4. 调试环境准备: 调试JS前启用Runtime，获取执行上下文与日志
// 5. 恢复功能: 执行Runtime.disable后重新启用Runtime功能
// 6. 上下文监听: 监听页面JS执行上下文创建、销毁事件

// CDPRuntimeEnable 启用Runtime域，开启JS执行与事件监听
func CDPRuntimeEnable() (string, error) {
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
		"method": "Runtime.enable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.enable 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.enable 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础启用Runtime域 ====================
func ExampleEnable_Basic() {
	// 初始化必须先启用Runtime域
	resp, err := CDPRuntimeEnable()
	if err != nil {
		log.Fatalf("启用Runtime域失败: %v", err)
	}
	log.Println("Runtime域已成功启用，响应:", resp)
}

// ==================== 使用示例 2：自动化测试初始化流程 ====================
func ExampleEnable_AutoTest() {
	// 自动化测试初始化步骤
	log.Println("开始初始化自动化环境...")

	// 1. 先启用Runtime域（核心前置依赖）
	_, err := CDPRuntimeEnable()
	if err != nil {
		log.Fatalf("初始化Runtime失败: %v", err)
	}

	// 2. 清空历史控制台日志
	CDPRuntimeDiscardConsoleEntries()

	// 3. 注入页面绑定函数
	CDPRuntimeAddBinding("pageCallback")

	log.Println("自动化Runtime环境初始化完成")
}

// ==================== 使用示例 3：关闭后重新启用Runtime ====================
func ExampleEnable_Reset() {
	// 先关闭Runtime域
	CDPRuntimeDisable()
	log.Println("Runtime域已关闭")

	// 业务逻辑处理...
	time.Sleep(1 * time.Second)

	// 重新启用Runtime恢复功能
	resp, err := CDPRuntimeEnable()
	if err != nil {
		log.Fatalf("重新启用Runtime失败: %v", err)
	}
	log.Println("Runtime域已重新启用，功能恢复正常")
}

*/

// -----------------------------------------------  Runtime.evaluate  -----------------------------------------------
// === 应用场景 ===
// 1. 执行JS代码: 在页面中直接执行任意JavaScript代码
// 2. 获取页面数据: 获取页面标题、URL、DOM内容、全局变量等数据
// 3. 自动化操作: 执行点击、输入、表单提交等页面操作
// 4. 调试诊断: 执行调试代码，检查页面状态、变量值
// 5. 异步代码执行: 执行async/await、Promise等异步JS逻辑
// 6. 页面环境修改: 修改页面全局变量、样式、属性等环境信息

// CDPRuntimeEvaluate 在页面执行JS表达式并返回结果
// expression: 要执行的JS表达式/代码
// returnByValue: 是否直接返回结果值(true)还是返回对象引用(false)
// awaitPromise: 是否自动等待Promise完成
func CDPRuntimeEvaluate(expression string, returnByValue bool, awaitPromise bool) (string, error) {
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
		"method": "Runtime.evaluate",
		"params": {
			"expression": %s,
			"returnByValue": %t,
			"awaitPromise": %t,
			"contextId": 0
		}
	}`, reqID, jsonEscape(expression), returnByValue, awaitPromise)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.evaluate 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 15 * time.Second
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.evaluate 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：执行简单JS获取页面标题 ====================
func ExampleEvaluate_Basic() {
	// 获取页面document.title
	resp, err := CDPRuntimeEvaluate(`document.title`, true, false)
	if err != nil {
		log.Fatalf("执行JS失败: %v", err)
	}
	log.Println("页面标题:", resp)
}

// ==================== 使用示例 2：执行异步Promise代码（等待结果） ====================
func ExampleEvaluate_AwaitPromise() {
	// 执行异步fetch请求，自动等待Promise完成
	js := `fetch('https://httpbin.org/get').then(res => res.json())`
	resp, err := CDPRuntimeEvaluate(js, true, true)
	if err != nil {
		log.Fatalf("执行异步JS失败: %v", err)
	}
	log.Println("异步接口返回结果:", resp)
}

// ==================== 使用示例 3：执行DOM操作（点击按钮） ====================
func ExampleEvaluate_DOM() {
	// 执行JS点击页面按钮
	js := `document.querySelector('#submit-btn').click(); '点击成功'`
	resp, err := CDPRuntimeEvaluate(js, true, false)
	if err != nil {
		log.Fatalf("执行DOM操作失败: %v", err)
	}
	log.Println("DOM操作结果:", resp)
}

// ==================== 使用示例 4：获取页面全局变量 ====================
func ExampleEvaluate_GlobalVar() {
	// 获取页面全局变量 userInfo
	resp, err := CDPRuntimeEvaluate(`window.userInfo`, true, false)
	if err != nil {
		log.Fatalf("获取全局变量失败: %v", err)
	}
	log.Println("页面用户信息:", resp)
}

*/

// -----------------------------------------------  Runtime.getProperties  -----------------------------------------------
// === 应用场景 ===
// 1. 对象属性解析: 获取JS对象/ DOM元素的所有属性、方法、自有属性
// 2. 调试数据提取: 调试时获取页面变量、组件实例的完整属性结构
// 3. 自动化数据抓取: 提取DOM元素、JS对象的完整状态与数据
// 4. 原型链分析: 获取对象的原型链属性，分析继承关系
// 5. 函数参数解析: 获取函数对象的参数、属性、闭包信息
// 6. 元素属性提取: 获取DOM元素的class、id、样式、属性等完整信息

// CDPRuntimeGetProperties 获取指定对象的所有属性
// objectID: 目标对象ID (来自evaluate、DOM等API返回的objectId)
// ownProperties: 是否只获取自身属性(true: 不包含原型链; false: 包含原型链)
// accessorPropertiesOnly: 是否只获取访问器属性(getter/setter)
// generatePreview: 是否生成属性预览数据
func CDPRuntimeGetProperties(objectID string, ownProperties, accessorPropertiesOnly, generatePreview bool) (string, error) {
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
		"method": "Runtime.getProperties",
		"params": {
			"objectId": "%s",
			"ownProperties": %t,
			"accessorPropertiesOnly": %t,
			"generatePreview": %t
		}
	}`, reqID, objectID, ownProperties, accessorPropertiesOnly, generatePreview)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.getProperties 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 10 * time.Second
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.getProperties 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：获取DOM元素自身属性 ====================
func ExampleGetProperties_DOMElement() {
	// DOM元素objectID (通过Runtime.evaluate获取)
	elemObjID := "123456"

	// 只获取元素自身属性，不包含原型链，生成预览
	resp, err := CDPRuntimeGetProperties(elemObjID, true, false, true)
	if err != nil {
		log.Fatalf("获取DOM元素属性失败: %v", err)
	}
	log.Println("DOM元素完整属性:", resp)
}

// ==================== 使用示例 2：获取JS对象所有属性(含原型链) ====================
func ExampleGetProperties_Object() {
	// JS对象objectID (如window.user、组件实例)
	userObjID := "789012"

	// 获取对象全部属性(含原型)，不限制属性类型，生成预览
	resp, err := CDPRuntimeGetProperties(userObjID, false, false, true)
	if err != nil {
		log.Fatalf("获取JS对象属性失败: %v", err)
	}
	log.Println("JS对象完整属性(含原型链):", resp)
}

// ==================== 使用示例 3：调试获取函数对象属性 ====================
func ExampleGetProperties_Function() {
	// 函数对象objectID
	funcObjID := "345678"

	// 只获取函数自身属性，用于调试分析
	resp, err := CDPRuntimeGetProperties(funcObjID, true, false, true)
	if err != nil {
		log.Fatalf("获取函数属性失败: %v", err)
	}
	log.Println("函数对象属性:", resp)
}

// ==================== 使用示例 4：仅获取对象的访问器属性(getter/setter) ====================
func ExampleGetProperties_Accessor() {
	objID := "901234"

	// 只获取getter/setter访问器属性
	resp, err := CDPRuntimeGetProperties(objID, true, true, false)
	if err != nil {
		log.Fatalf("获取访问器属性失败: %v", err)
	}
	log.Println("对象访问器属性:", resp)
}

*/

// -----------------------------------------------  Runtime.globalLexicalScopeNames  -----------------------------------------------
// === 应用场景 ===
// 1. 全局变量扫描: 获取页面顶层通过 let/const 声明的全局变量名称
// 2. 页面变量审计: 审计页面全局作用域污染情况，排查冗余全局变量
// 3. 调试分析: 分析页面顶层词法作用域变量，定位全局变量问题
// 4. 自动化检测: 检测页面是否声明了非法/敏感全局词法变量
// 5. 框架变量识别: 提取前端框架(React/Vue/Angular)挂载的全局词法变量
// 6. 代码规范检查: 检查页面是否存在过多不规范的全局let/const声明

// CDPRuntimeGlobalLexicalScopeNames 获取全局词法作用域中的变量名（let/const 声明）
// executionContextID: 执行上下文ID（0表示默认主上下文）
func CDPRuntimeGlobalLexicalScopeNames(executionContextID int) (string, error) {
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
		"method": "Runtime.globalLexicalScopeNames",
		"params": {
			"executionContextId": %d
		}
	}`, reqID, executionContextID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.globalLexicalScopeNames 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.globalLexicalScopeNames 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：获取默认上下文全局词法变量 ====================
func ExampleGlobalLexicalScopeNames_DefaultContext() {
	// 获取默认执行上下文（主页面）的 let/const 全局变量
	resp, err := CDPRuntimeGlobalLexicalScopeNames(0)
	if err != nil {
		log.Fatalf("获取全局词法变量失败: %v", err)
	}
	log.Println("页面全局词法作用域变量(let/const):", resp)
}

// ==================== 使用示例 2：指定执行上下文获取词法变量 ====================
func ExampleGlobalLexicalScopeNames_CustomContext() {
	// 指定执行上下文ID（来自Runtime.createExecutionContext创建）
	contextID := 1001

	resp, err := CDPRuntimeGlobalLexicalScopeNames(contextID)
	if err != nil {
		log.Fatalf("获取指定上下文词法变量失败: %v", err)
	}
	log.Println("自定义上下文词法变量:", resp)
}

// ==================== 使用示例 3：自动化审计页面全局变量规范 ====================
func ExampleGlobalLexicalScopeNames_Audit() {
	// 审计页面是否存在过多全局词法变量
	resp, err := CDPRuntimeGlobalLexicalScopeNames(0)
	if err != nil {
		log.Fatalf("变量审计失败: %v", err)
	}

	// 解析结果并检查变量数量
	log.Println("页面全局词法变量扫描完成:", resp)
	// 可进一步判断是否符合代码规范
}

*/

// -----------------------------------------------  Runtime.queryObjects  -----------------------------------------------
// === 应用场景 ===
// 1. 全局实例检索: 获取页面中所有由指定构造函数（如Array、Map、自定义类）创建的对象实例
// 2. 内存分析: 查找页面中所有某类对象，辅助内存泄漏排查
// 3. 数据抓取: 批量获取页面中所有数组、Map、Set等数据结构
// 4. 框架实例扫描: 查找Vue/React组件实例、自定义类的所有全局实例
// 5. 调试溯源: 定位某类对象的所有创建实例，分析数据来源
// 6. 状态收集: 收集页面中所有某类状态对象，用于调试与监控

// CDPRuntimeQueryObjects 根据原型对象ID，查询所有继承该原型的对象
// prototypeObjectID: 原型对象的ID（通过evaluate获取构造函数的prototype）
func CDPRuntimeQueryObjects(prototypeObjectID string) (string, error) {
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
		"method": "Runtime.queryObjects",
		"params": {
			"prototypeObjectId": "%s"
		}
	}`, reqID, prototypeObjectID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.queryObjects 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 10 * time.Second
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.queryObjects 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：获取页面所有 Array 数组实例 ====================
func ExampleQueryObjects_Array() {
	// 1. 先获取 Array.prototype 的对象ID
	evalResp, _ := CDPRuntimeEvaluate(`Array.prototype`, false, false)
	// 解析 evalResp 得到 prototypeObjectID
	arrayProtoID := "123456"

	// 2. 查询页面中所有数组对象
	resp, err := CDPRuntimeQueryObjects(arrayProtoID)
	if err != nil {
		log.Fatalf("查询所有数组失败: %v", err)
	}
	log.Println("页面所有数组实例:", resp)
}

// ==================== 使用示例 2：获取页面所有 Map 实例 ====================
func ExampleQueryObjects_Map() {
	// 获取 Map.prototype 对象ID
	mapProtoID := "789012"

	// 查询所有 Map 实例
	resp, err := CDPRuntimeQueryObjects(mapProtoID)
	if err != nil {
		log.Fatalf("查询所有Map失败: %v", err)
	}
	log.Println("页面所有Map实例:", resp)
}

// ==================== 使用示例 3：获取页面所有自定义类实例 ====================
func ExampleQueryObjects_CustomClass() {
	// 页面有自定义类：class User {}
	// 获取 User.prototype 的对象ID
	userProtoID := "345678"

	// 查询页面所有 User 类实例
	resp, err := CDPRuntimeQueryObjects(userProtoID)
	if err != nil {
		log.Fatalf("查询User实例失败: %v", err)
	}
	log.Println("页面所有User类实例:", resp)
}

// ==================== 使用示例 4：内存泄漏排查 - 检索指定对象 ====================
func ExampleQueryObjects_MemoryDebug() {
	// 排查某个类是否产生多余实例
	protoID := "901234"

	resp, err := CDPRuntimeQueryObjects(protoID)
	if err != nil {
		log.Fatalf("内存排查查询失败: %v", err)
	}
	log.Println("目标对象所有实例（用于内存分析）:", resp)
}

*/

// -----------------------------------------------  Runtime.releaseObject  -----------------------------------------------
// === 应用场景 ===
// 1. 手动释放对象引用: 释放通过CDP获取的JS对象引用，避免内存泄漏
// 2. 自动化资源清理: 执行完对象操作后及时释放，长时间运行不卡顿
// 3. 内存优化: 爬虫/自动化工具高频操作对象时，主动回收内存
// 4. 引用计数管理: 控制浏览器端对象引用计数，防止堆积
// 5. 长任务维护: 长时间运行自动化流程，定期释放无用对象
// 6. 调试后清理: 调试获取对象属性、结构后，释放引用资源

// CDPRuntimeReleaseObject 释放指定的JS对象引用（objectId）
// objectID: 要释放的对象ID（来自evaluate、getProperties等返回）
func CDPRuntimeReleaseObject(objectID string) (string, error) {
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
		"method": "Runtime.releaseObject",
		"params": {
			"objectId": "%s"
		}
	}`, reqID, objectID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.releaseObject 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.releaseObject 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：获取属性后立即释放对象 ====================
func ExampleReleaseObject_Basic() {
	// 1. 获取DOM元素对象ID
	elemObjID := "123456"

	// 2. 获取元素属性
	CDPRuntimeGetProperties(elemObjID, true, false, true)

	// 3. 立即释放对象引用，防止内存泄漏
	resp, err := CDPRuntimeReleaseObject(elemObjID)
	if err != nil {
		log.Fatalf("释放对象失败: %v", err)
	}
	log.Println("对象已成功释放，内存回收:", resp)
}

// ==================== 使用示例 2：自动化循环操作后批量释放 ====================
func ExampleReleaseObject_Loop() {
	// 模拟批量抓取页面元素
	objIDs := []string{"obj-1", "obj-2", "obj-3", "obj-4"}

	for _, id := range objIDs {
		// 操作对象...
		log.Println("处理对象:", id)

		// 操作完成立即释放
		CDPRuntimeReleaseObject(id)
	}

	log.Println("批量对象处理完成，全部引用已释放，无内存泄漏")
}

// ==================== 使用示例 3：Promise对象使用后释放 ====================
func ExampleReleaseObject_Promise() {
	// Promise对象ID
	promiseObjID := "789012"

	// 等待Promise执行完成
	CDPRuntimeAwaitPromise(promiseObjID, true)

	// 释放Promise对象引用
	resp, err := CDPRuntimeReleaseObject(promiseObjID)
	if err != nil {
		log.Fatalf("释放Promise对象失败: %v", err)
	}
	log.Println("Promise对象已释放:", resp)
}

// ==================== 使用示例 4：长时间爬虫任务内存优化 ====================
func ExampleReleaseObject_Optimize() {
	// 抓取数据后必须释放
	dataObjID := "901234"

	// 提取数据...
	log.Println("提取页面数据完成")

	// 释放对象，优化内存
	CDPRuntimeReleaseObject(dataObjID)
	log.Println("数据对象已释放，内存占用降低")
}

*/

// -----------------------------------------------  Runtime.releaseObjectGroup  -----------------------------------------------
// === 应用场景 ===
// 1. 批量释放对象引用: 一次性释放同一分组下的所有JS对象引用，高效清理内存
// 2. 长时运行优化: 自动化/爬虫长时间运行时，定期批量释放对象，防止浏览器内存溢出
// 3. 任务隔离清理: 每个自动化任务使用独立对象组，任务结束后一键释放该组所有对象
// 4. 资源批量回收: 替代多次调用releaseObject，批量回收提升性能
// 5. 页面切换清理: 切换页面/任务时，释放上一个任务的所有对象引用
// 6. 调试会话清理: 单次调试结束后，释放该会话创建的所有对象，不影响其他会话

// CDPRuntimeReleaseObjectGroup 批量释放指定对象组下的所有JS对象引用
// objectGroup: 对象组名称（通过evaluate、callFunctionOn指定的group名）
func CDPRuntimeReleaseObjectGroup(objectGroup string) (string, error) {
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
		"method": "Runtime.releaseObjectGroup",
		"params": {
			"objectGroup": "%s"
		}
	}`, reqID, objectGroup)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.releaseObjectGroup 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.releaseObjectGroup 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础批量释放对象组 ====================
func ExampleReleaseObjectGroup_Basic() {
	// 定义对象组名：task_group_1001
	groupName := "task_group_1001"

	// 一键释放该组下所有对象引用
	resp, err := CDPRuntimeReleaseObjectGroup(groupName)
	if err != nil {
		log.Fatalf("批量释放对象组失败: %v", err)
	}
	log.Println("对象组所有引用已释放，内存回收完成:", resp)
}

// ==================== 使用示例 2：自动化任务完成后清理 ====================
func ExampleReleaseObjectGroup_AutoTask() {
	// 1. 执行自动化任务，使用专属对象组
	taskGroup := "login_task_group"
	log.Println("执行登录任务，对象组:", taskGroup)
	// 执行JS、获取对象时指定 objectGroup: taskGroup

	// 2. 任务完成，批量释放该任务所有对象（无需逐个释放）
	resp, err := CDPRuntimeReleaseObjectGroup(taskGroup)
	if err != nil {
		log.Fatalf("任务清理失败: %v", err)
	}
	log.Println("登录任务完成，所有对象已批量释放:", resp)
}

// ==================== 使用示例 3：定时批量释放（长时间爬虫优化） ====================
func ExampleReleaseObjectGroup_Timer() {
	// 爬虫核心对象组
	crawlGroup := "crawl_data_group"

	// 每30秒批量释放一次，防止内存堆积
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := CDPRuntimeReleaseObjectGroup(crawlGroup)
		if err != nil {
			log.Printf("定时释放失败: %v", err)
			continue
		}
		log.Println("爬虫对象组已定时清理，内存优化完成")
	}
}

// ==================== 使用示例 4：页面切换后释放上一页对象 ====================
func ExampleReleaseObjectGroup_PageSwitch() {
	// 上一个页面的对象组
	lastPageGroup := "previous_page_group"

	// 切换页面前，释放上一页所有对象引用
	resp, err := CDPRuntimeReleaseObjectGroup(lastPageGroup)
	if err != nil {
		log.Fatalf("释放上页资源失败: %v", err)
	}
	log.Println("上一页面所有对象已释放，准备切换新页面")
}

*/

// -----------------------------------------------  Runtime.removeBinding  -----------------------------------------------
// === 应用场景 ===
// 1. 解绑通信函数: 移除通过addBinding注入的JS-Go通信函数
// 2. 测试环境清理: 自动化测试完成后清理页面注入的绑定函数
// 3. 安全隔离: 移除敏感操作绑定，防止页面后续非法调用
// 4. 动态更新: 替换绑定函数前先移除旧绑定，避免冲突
// 5. 资源释放: 清理不再使用的绑定，释放监听资源
// 6. 多任务切换: 切换任务场景时移除上一个任务的绑定函数

// CDPRuntimeRemoveBinding 移除页面注入的绑定函数
// bindingName: 要移除的绑定函数名称
func CDPRuntimeRemoveBinding(bindingName string) (string, error) {
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
		"method": "Runtime.removeBinding",
		"params": {
			"name": "%s"
		}
	}`, reqID, bindingName)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.removeBinding 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.removeBinding 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础移除绑定函数 ====================
func ExampleRemoveBinding_Basic() {
	// 移除名为 sendLogToGo 的绑定函数
	resp, err := CDPRuntimeRemoveBinding("sendLogToGo")
	if err != nil {
		log.Fatalf("移除绑定函数失败: %v", err)
	}
	log.Println("绑定函数已成功移除，JS无法再调用:", resp)
}

// ==================== 使用示例 2：自动化测试后清理所有绑定 ====================
func ExampleRemoveBinding_AutoTest() {
	// 测试前注入绑定
	CDPRuntimeAddBinding("testCallback")
	CDPRuntimeAddBinding("reportResult")

	log.Println("执行自动化测试逻辑...")

	// 测试完成后批量移除所有绑定
	CDPRuntimeRemoveBinding("testCallback")
	CDPRuntimeRemoveBinding("reportResult")

	log.Println("测试完成，所有绑定函数已清理，页面环境恢复")
}

// ==================== 使用示例 3：动态替换绑定函数 ====================
func ExampleRemoveBinding_Replace() {
	// 旧绑定函数
	oldBinding := "oldApi"
	// 先移除旧绑定
	CDPRuntimeRemoveBinding(oldBinding)
	log.Println("旧绑定已移除")

	// 注入新的绑定函数
	CDPRuntimeAddBinding("newApi")
	log.Println("新绑定已注入，函数替换完成")
}

// ==================== 使用示例 4：安全场景移除敏感绑定 ====================
func ExampleRemoveBinding_Security() {
	// 敏感操作绑定：文件操作、系统命令等
	sensitiveBinding := "doSystemAction"

	// 完成敏感操作后立即移除，防止恶意调用
	resp, err := CDPRuntimeRemoveBinding(sensitiveBinding)
	if err != nil {
		log.Fatalf("移除敏感绑定失败: %v", err)
	}
	log.Println("敏感操作绑定已移除，安全风险降低")
}

*/

// -----------------------------------------------  Runtime.runIfWaitingForDebugger  -----------------------------------------------
// === 应用场景 ===
// 1. 调试控制: 当页面启动后暂停等待调试器时，发送命令让页面继续执行
// 2. 自动化调试流程: 自动化完成调试初始化后，恢复页面运行
// 3. 断点调试同步: 调试工具准备就绪后，通知浏览器开始执行JS
// 4. 调试生命周期管理: 控制页面在调试模式下的启动与运行时机
// 5. 无头浏览器调试: 无界面浏览器调试时，手动控制脚本运行
// 6. 调试环境就绪: 调试器连接完成后放行页面JS执行

// CDPRuntimeRunIfWaitingForDebugger 通知浏览器停止等待调试器，继续执行页面脚本
func CDPRuntimeRunIfWaitingForDebugger() (string, error) {
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
		"method": "Runtime.runIfWaitingForDebugger"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.runIfWaitingForDebugger 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.runIfWaitingForDebugger 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础调试恢复执行 ====================
func ExampleRunIfWaitingForDebugger_Basic() {
	// 浏览器启动时配置了 --js-flags="--wait-for-debugger"
	// 调试器已连接，通知页面继续执行JS
	resp, err := CDPRuntimeRunIfWaitingForDebugger()
	if err != nil {
		log.Fatalf("通知页面继续执行失败: %v", err)
	}
	log.Println("调试器就绪，页面已恢复执行:", resp)
}

// ==================== 使用示例 2：自动化调试流程 ====================
func ExampleRunIfWaitingForDebugger_AutoDebug() {
	log.Println("1. 启动浏览器（等待调试器模式）")
	log.Println("2. 成功连接CDP调试接口")
	log.Println("3. 启用Runtime域")
	CDPRuntimeEnable()

	// 4. 所有调试准备完成，放行页面执行
	resp, err := CDPRuntimeRunIfWaitingForDebugger()
	if err != nil {
		log.Fatalf("自动化调试启动失败: %v", err)
	}
	log.Println("自动化调试流程完成，页面正常运行:", resp)
}

// ==================== 使用示例 3：无头浏览器调试控制 ====================
func ExampleRunIfWaitingForDebugger_Headless() {
	// 无头浏览器环境，页面等待调试器
	log.Println("无头浏览器已暂停，等待调试指令...")

	// 执行恢复运行命令
	_, err := CDPRuntimeRunIfWaitingForDebugger()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("无头浏览器恢复执行，业务脚本开始运行")
}

*/

// -----------------------------------------------  Runtime.runScript  -----------------------------------------------
// === 应用场景 ===
// 1. 执行预编译脚本: 运行通过compileScript提前编译好的脚本，提升执行效率
// 2. 隔离脚本执行: 在指定执行上下文运行脚本，实现环境隔离
// 3. 复用编译结果: 多次执行同一编译好的脚本，避免重复编译
// 4. 精准执行控制: 精确控制脚本执行的上下文、时机、返回方式
// 5. 模块化执行: 按模块执行编译后的JS脚本，实现分步执行
// 6. 调试执行: 执行指定脚本并调试，支持等待Promise、异常捕获

// CDPRuntimeRunScript 执行预编译的脚本（通过scriptId）
// scriptID: 预编译脚本ID（来自compileScript返回）
// executionContextID: 执行上下文ID（0=默认上下文）
// returnByValue: 是否直接返回结果值
// awaitPromise: 是否自动等待Promise完成
func CDPRuntimeRunScript(scriptID string, executionContextID int, returnByValue bool, awaitPromise bool) (string, error) {
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
		"method": "Runtime.runScript",
		"params": {
			"scriptId": "%s",
			"executionContextId": %d,
			"returnByValue": %t,
			"awaitPromise": %t
		}
	}`, reqID, scriptID, executionContextID, returnByValue, awaitPromise)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.runScript 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 15 * time.Second
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.runScript 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：执行预编译脚本（基础） ====================
func ExampleRunScript_Basic() {
	// 脚本ID来自 CDPRuntimeCompileScript 返回
	scriptID := "123456"

	// 在默认上下文执行，返回结果值，不等待Promise
	resp, err := CDPRuntimeRunScript(scriptID, 0, true, false)
	if err != nil {
		log.Fatalf("执行预编译脚本失败: %v", err)
	}
	log.Println("脚本执行成功，结果:", resp)
}

// ==================== 使用示例 2：执行异步脚本（自动等待Promise） ====================
func ExampleRunScript_AwaitPromise() {
	// 包含异步逻辑的预编译脚本ID
	scriptID := "789012"

	// 执行异步脚本，自动等待Promise完成并返回值
	resp, err := CDPRuntimeRunScript(scriptID, 0, true, true)
	if err != nil {
		log.Fatalf("执行异步脚本失败: %v", err)
	}
	log.Println("异步脚本执行完成，结果:", resp)
}

// ==================== 使用示例 3：在指定执行上下文执行脚本 ====================
func ExampleRunScript_CustomContext() {
	scriptID := "345678"
	// 指定自定义执行上下文ID（来自createExecutionContext）
	contextID := 1001

	// 在隔离上下文执行脚本
	resp, err := CDPRuntimeRunScript(scriptID, contextID, true, false)
	if err != nil {
		log.Fatalf("指定上下文执行脚本失败: %v", err)
	}
	log.Println("隔离环境执行成功:", resp)
}

// ==================== 使用示例 4：多次复用执行同一编译脚本 ====================
func ExampleRunScript_Reuse() {
	// 一次编译，多次执行
	scriptID := "901234"

	// 第1次执行
	CDPRuntimeRunScript(scriptID, 0, true, false)
	// 第2次执行
	CDPRuntimeRunScript(scriptID, 0, true, false)
	// 第3次执行
	CDPRuntimeRunScript(scriptID, 0, true, false)

	log.Println("脚本复用执行完成，无需重复编译，效率更高")
}

*/

// -----------------------------------------------  Runtime.setAsyncCallStackDepth  -----------------------------------------------
// === 应用场景 ===
// 1. 异步调试增强: 开启Promise/async/await的完整异步调用栈追踪
// 2. 错误定位: 异步代码报错时，显示完整调用链，快速定位根源
// 3. 性能分析: 分析异步任务执行链路，排查性能瓶颈
// 4. 复杂业务调试: 多层嵌套异步逻辑下，还原完整执行路径
// 5. 自动化问题排查: 自动化测试中异步操作失败，追踪调用栈
// 6. 内存泄漏追踪: 异步任务导致的内存问题，通过调用栈溯源

// CDPRuntimeSetAsyncCallStackDepth 设置异步调用栈深度，启用异步栈追踪
// depth: 异步调用栈最大深度（建议32/64，0=关闭追踪）
func CDPRuntimeSetAsyncCallStackDepth(depth int) (string, error) {
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
		"method": "Runtime.setAsyncCallStackDepth",
		"params": {
			"maxDepth": %d
		}
	}`, reqID, depth)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.setAsyncCallStackDepth 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.setAsyncCallStackDepth 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：开启标准异步栈追踪（推荐） ====================
func ExampleSetAsyncCallStackDepth_Enable() {
	// 设置最大深度32，满足绝大多数复杂异步调试
	resp, err := CDPRuntimeSetAsyncCallStackDepth(32)
	if err != nil {
		log.Fatalf("开启异步调用栈失败: %v", err)
	}
	log.Println("已启用异步调用栈追踪，深度32:", resp)
}

// ==================== 使用示例 2：关闭异步栈追踪（释放性能） ====================
func ExampleSetAsyncCallStackDepth_Disable() {
	// 调试完成后设置为0，关闭追踪，降低浏览器内存消耗
	resp, err := CDPRuntimeSetAsyncCallStackDepth(0)
	if err != nil {
		log.Fatalf("关闭异步栈失败: %v", err)
	}
	log.Println("已关闭异步调用栈追踪，释放浏览器资源:", resp)
}

// ==================== 使用示例 3：调试深层嵌套异步代码 ====================
func ExampleSetAsyncCallStackDepth_DeepDebug() {
	// 多层嵌套Promise/async逻辑，设置更大深度64
	resp, err := CDPRuntimeSetAsyncCallStackDepth(64)
	if err != nil {
		log.Fatalf("设置深层异步栈失败: %v", err)
	}
	log.Println("已启用深层异步调用栈（64层），准备调试复杂异步业务:", resp)

	// 执行异步JS代码，报错时可看到完整调用链
	CDPRuntimeEvaluate(`async function run(){ await Promise.resolve().then(()=>{throw new Error("异步测试错误")}) } run()`, true, true)
}

// ==================== 使用示例 4：自动化调试初始化 ====================
func ExampleSetAsyncCallStackDepth_AutoTest() {
	// 自动化测试初始化：启用异步栈，方便排查异步失败
	log.Println("初始化自动化测试环境...")
	CDPRuntimeEnable()
	// 开启异步栈追踪，深度32
	CDPRuntimeSetAsyncCallStackDepth(32)

	log.Println("自动化环境初始化完成，异步问题可完整追踪调用栈")
}

*/

// -----------------------------------------------  Runtime.getExceptionDetails  -----------------------------------------------
// === 应用场景 ===
// 1. 异常详情获取: 根据异常ID获取JS执行报错的完整堆栈、信息、位置
// 2. 自动化错误捕获: 测试/爬虫执行JS出错时，提取详细错误定位问题
// 3. 调试精准定位: 获取错误行号、列号、堆栈、源码片段，快速修复BUG
// 4. 异步错误解析: 解析Promise/async抛出的异常完整信息
// 5. 错误日志上报: 收集页面JS异常详情，用于监控与上报
// 6. 崩溃分析: 页面脚本崩溃时，提取崩溃上下文与堆栈信息

// CDPRuntimeGetExceptionDetails 根据异常ID获取JS异常详细信息
// exceptionID: 异常唯一ID（来自异常事件、执行失败返回的exceptionId）
func CDPRuntimeGetExceptionDetails(exceptionID int) (string, error) {
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
		"method": "Runtime.getExceptionDetails",
		"params": {
			"exceptionId": %d
		}
	}`, reqID, exceptionID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.getExceptionDetails 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.getExceptionDetails 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础获取异常详情 ====================
func ExampleGetExceptionDetails_Basic() {
	// 异常ID来自执行JS失败时返回的 exceptionId
	exceptionID := 1001

	// 获取完整异常详情（堆栈、行号、错误信息）
	resp, err := CDPRuntimeGetExceptionDetails(exceptionID)
	if err != nil {
		log.Fatalf("获取异常详情失败: %v", err)
	}
	log.Println("JS异常完整信息:", resp)
}

// ==================== 使用示例 2：执行JS出错后捕获异常 ====================
func ExampleGetExceptionDetails_ExecuteError() {
	// 执行一段会报错的JS代码
	// 执行后会得到 exceptionId = 1002
	jsCode := `document.querySelector(null).click()`
	_, err := CDPRuntimeEvaluate(jsCode, true, false)
	if err != nil {
		log.Println("执行JS失败，开始获取异常详情...")
	}

	// 获取异常详情
	resp, _ := CDPRuntimeGetExceptionDetails(1002)
	log.Println("异常定位信息:", resp)
}

// ==================== 使用示例 3：自动化测试异常上报 ====================
func ExampleGetExceptionDetails_AutoTest() {
	// 测试用例执行异常
	exceptionID := 2001

	// 获取异常并格式化上报
	details, err := CDPRuntimeGetExceptionDetails(exceptionID)
	if err != nil {
		log.Fatalf("测试异常获取失败: %v", err)
	}

	log.Printf("[测试异常] 用例执行失败，异常详情:\n%s", details)
	// 可上传到错误监控平台
}

// ==================== 使用示例 4：异步Promise异常解析 ====================
func ExampleGetExceptionDetails_Async() {
	// 异步代码抛出异常ID
	exceptionID := 3001

	// 获取异步异常完整堆栈
	resp, err := CDPRuntimeGetExceptionDetails(exceptionID)
	if err != nil {
		log.Fatalf("异步异常解析失败: %v", err)
	}
	log.Println("异步JS异常详情:", resp)
}

*/

// -----------------------------------------------  Runtime.getHeapUsage  -----------------------------------------------
// === 应用场景 ===
// 1. 内存监控: 实时监控浏览器V8引擎的堆内存使用情况
// 2. 内存泄漏排查: 定时获取内存数据，分析是否存在持续上涨
// 3. 性能优化: 根据内存占用判断是否需要释放对象、触发GC
// 4. 自动化健康检查: 长时间运行自动化任务时，检测内存是否超标
// 5. 资源告警: 内存占用过高时触发告警，防止浏览器崩溃
// 6. 调试分析: 对比操作前后的内存变化，定位内存占用大的逻辑

// CDPRuntimeGetHeapUsage 获取V8堆内存使用情况
func CDPRuntimeGetHeapUsage() (string, error) {
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
		"method": "Runtime.getHeapUsage"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.getHeapUsage 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.getHeapUsage 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：单次获取内存使用 ====================
func ExampleGetHeapUsage_Basic() {
	// 获取当前堆内存使用：usedSize(已使用) / totalSize(总分配)
	resp, err := CDPRuntimeGetHeapUsage()
	if err != nil {
		log.Fatalf("获取堆内存失败: %v", err)
	}
	log.Println("堆内存使用情况:", resp)
}

// ==================== 使用示例 2：定时监控内存（排查泄漏） ====================
func ExampleGetHeapUsage_Monitor() {
	// 每5秒获取一次内存，监控是否持续上涨
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, _ := CDPRuntimeGetHeapUsage()
		log.Println("[内存监控] 当前堆使用:", resp)
	}
}

// ==================== 使用示例 3：自动化任务内存阈值告警 ====================
func ExampleGetHeapUsage_Threshold() {
	resp, err := CDPRuntimeGetHeapUsage()
	if err != nil {
		log.Fatal(err)
	}

	// 解析返回结果，判断内存是否超过阈值
	log.Println("自动化任务内存检查完成:", resp)

	// 示例逻辑：
	// if used > 800MB {
	//     执行释放对象、触发GC、告警
	// }
}

// ==================== 使用示例 4：操作前后内存对比 ====================
func ExampleGetHeapUsage_Compare() {
	// 操作前内存
	before, _ := CDPRuntimeGetHeapUsage()
	log.Println("操作前内存:", before)

	// 执行页面操作（抓取数据、DOM操作）
	log.Println("执行页面业务操作...")

	// 操作后内存
	after, _ := CDPRuntimeGetHeapUsage()
	log.Println("操作后内存:", after)

	log.Println("内存变化对比完成，可用于定位泄漏点")
}

*/

// -----------------------------------------------  Runtime.getIsolateId  -----------------------------------------------
// === 应用场景 ===
// 1. 隔离环境标识: 获取V8引擎Isolate唯一ID，区分不同JS运行隔离环境
// 2. 多实例管理: 多标签页/多渲染进程场景下，标识独立V8实例
// 3. 调试追踪: 调试时定位异常、日志所属的具体V8隔离环境
// 4. 内存隔离区分: 明确堆内存、对象归属的V8隔离实例
// 5. 自动化隔离: 多任务自动化时，区分任务所属的JS运行环境
// 6. 日志溯源: 给日志、异常添加Isolate标识，方便问题定位

// CDPRuntimeGetIsolateId 获取当前V8引擎Isolate（隔离实例）的唯一ID
func CDPRuntimeGetIsolateId() (string, error) {
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
		"method": "Runtime.getIsolateId"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.getIsolateId 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.getIsolateId 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础获取V8隔离实例ID ====================
func ExampleGetIsolateId_Basic() {
	// 获取当前页面V8 Isolate唯一标识
	resp, err := CDPRuntimeGetIsolateId()
	if err != nil {
		log.Fatalf("获取Isolate ID失败: %v", err)
	}
	log.Println("当前V8隔离实例ID:", resp)
}

// ==================== 使用示例 2：多页面/多实例标识区分 ====================
func ExampleGetIsolateId_MultiPage() {
	// 场景：打开多个页面，分别获取Isolate ID进行区分
	page1Isolate, _ := CDPRuntimeGetIsolateId()
	log.Println("页面1 Isolate ID:", page1Isolate)

	// 切换页面后...
	page2Isolate, _ := CDPRuntimeGetIsolateId()
	log.Println("页面2 Isolate ID:", page2Isolate)

	// 通过ID区分不同页面的JS运行环境
}

// ==================== 使用示例 3：调试日志添加隔离环境标识 ====================
func ExampleGetIsolateId_LogTrace() {
	// 获取ID并添加到日志，方便问题溯源
	isolateID, _ := CDPRuntimeGetIsolateId()
	log.Printf("[Isolate:%s] 执行自动化任务", isolateID)

	// 执行JS操作
	CDPRuntimeEvaluate(`console.log("任务执行")`, true, false)
}

// ==================== 使用示例 4：内存/异常归属定位 ====================
func ExampleGetIsolateId_Debug() {
	// 先获取隔离ID
	isolateID, _ := CDPRuntimeGetIsolateId()
	// 获取对应Isolate的内存使用
	heap, _ := CDPRuntimeGetHeapUsage()

	log.Printf("归属隔离实例[%s] 内存使用: %s", isolateID, heap)
}

*/

// -----------------------------------------------  Runtime.setCustomObjectFormatterEnabled  -----------------------------------------------
// === 应用场景 ===
// 1. 调试对象格式化：启用/禁用页面自定义的对象格式化器（如React、Vue组件自定义展示）
// 2. 开发者体验优化：调试时让复杂对象在控制台/调试器中显示更友好的结构
// 3. 框架调试适配：支持前端框架自定义对象格式化，提升调试效率
// 4. 调试一致性控制：强制关闭自定义格式化，保证调试时看到原生对象结构
// 5. 自动化调试配置：初始化调试环境时统一设置格式化行为
// 6. 排查格式化冲突：关闭自定义格式化，避免对象展示异常

// CDPRuntimeSetCustomObjectFormatterEnabled 启用或禁用页面自定义对象格式化器
// enabled: true=启用自定义格式化；false=禁用，使用浏览器原生格式化
func CDPRuntimeSetCustomObjectFormatterEnabled(enabled bool) (string, error) {
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
		"method": "Runtime.setCustomObjectFormatterEnabled",
		"params": {
			"enabled": %t
		}
	}`, reqID, enabled)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.setCustomObjectFormatterEnabled 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.setCustomObjectFormatterEnabled 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：启用自定义对象格式化（框架调试推荐） ====================
func ExampleSetCustomObjectFormatterEnabled_Enable() {
	// 启用页面自定义格式化（React/Vue组件会显示友好结构）
	resp, err := CDPRuntimeSetCustomObjectFormatterEnabled(true)
	if err != nil {
		log.Fatalf("启用自定义对象格式化失败: %v", err)
	}
	log.Println("已启用自定义对象格式化，调试体验优化:", resp)
}

// ==================== 使用示例 2：禁用自定义格式化（查看原生对象） ====================
func ExampleSetCustomObjectFormatterEnabled_Disable() {
	// 禁用，强制使用浏览器原生对象展示，用于底层调试
	resp, err := CDPRuntimeSetCustomObjectFormatterEnabled(false)
	if err != nil {
		log.Fatalf("禁用自定义对象格式化失败: %v", err)
	}
	log.Println("已禁用自定义格式化，显示对象原生结构:", resp)
}

// ==================== 使用示例 3：自动化调试环境初始化 ====================
func ExampleSetCustomObjectFormatterEnabled_Init() {
	// 初始化调试环境
	CDPRuntimeEnable()

	// 启用友好的对象格式化
	CDPRuntimeSetCustomObjectFormatterEnabled(true)

	log.Println("调试环境初始化完成：自定义格式化已启用")
}

// ==================== 使用示例 4：排查对象展示异常 ====================
func ExampleSetCustomObjectFormatterEnabled_Debug() {
	// 发现控制台对象展示异常，先关闭自定义格式化
	CDPRuntimeSetCustomObjectFormatterEnabled(false)
	log.Println("已关闭自定义格式化，对象展示恢复原生模式")

	// 执行调试获取对象属性
	// CDPRuntimeGetProperties(objectId, ...)
}

*/

// -----------------------------------------------  Runtime.setMaxCallStackSizeToCapture  -----------------------------------------------
// === 应用场景 ===
// 1. 调试栈深度控制: 设置捕获调用栈的最大深度，平衡调试信息与性能
// 2. 复杂栈调试: 深层递归/调用时，捕获完整栈信息，避免丢失关键路径
// 3. 性能优化: 生产环境减小栈深度，降低内存与CPU消耗
// 4. 错误日志优化: 控制异常上报时的栈长度，精简日志大小
// 5. 递归问题排查: 递归函数出错时，获取足够深的调用栈定位死循环
// 6. 调试精度调节: 根据场景自由调整栈捕获深度

// CDPRuntimeSetMaxCallStackSizeToCapture 设置要捕获的最大调用栈深度
// size: 调用栈最大捕获深度（建议：调试32-128，生产10-32）
func CDPRuntimeSetMaxCallStackSizeToCapture(size int) (string, error) {
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
		"method": "Runtime.setMaxCallStackSizeToCapture",
		"params": {
			"size": %d
		}
	}`, reqID, size)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.setMaxCallStackSizeToCapture 请求失败: %w", err)
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.setMaxCallStackSizeToCapture 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：调试模式设置大深度栈（推荐） ====================
func ExampleSetMaxCallStackSizeToCapture_Debug() {
	// 调试递归/深层调用，设置栈深度64
	resp, err := CDPRuntimeSetMaxCallStackSizeToCapture(64)
	if err != nil {
		log.Fatalf("设置调用栈深度失败: %v", err)
	}
	log.Println("已设置调用栈最大捕获深度：64层，调试信息完整", resp)
}

// ==================== 使用示例 2：生产环境精简栈深度（优化性能） ====================
func ExampleSetMaxCallStackSizeToCapture_Prod() {
	// 生产环境减小栈深度，降低性能开销
	resp, err := CDPRuntimeSetMaxCallStackSizeToCapture(16)
	if err != nil {
		log.Fatalf("设置精简栈深度失败: %v", err)
	}
	log.Println("已设置生产环境栈深度：16层，性能最优", resp)
}

// ==================== 使用示例 3：排查递归死循环问题 ====================
func ExampleSetMaxCallStackSizeToCapture_Recursive() {
	// 排查深层递归，设置最大深度128
	resp, err := CDPRuntimeSetMaxCallStackSizeToCapture(128)
	if err != nil {
		log.Fatalf("设置递归栈深度失败: %v", err)
	}
	log.Println("已启用最大栈深度128层，准备排查递归问题", resp)
}

// ==================== 使用示例 4：自动化调试初始化 ====================
func ExampleSetMaxCallStackSizeToCapture_AutoTest() {
	// 初始化测试环境
	CDPRuntimeEnable()
	// 设置标准调试栈深度32
	CDPRuntimeSetMaxCallStackSizeToCapture(32)
	log.Println("自动化环境初始化完成：调用栈深度32层")
}

*/

// -----------------------------------------------  Runtime.terminateExecution  -----------------------------------------------
// === 应用场景 ===
// 1. 强制终止JS执行：立即停止当前正在运行的耗时/死循环/卡死JS代码
// 2. 超时保护：自动化执行JS超时后，强制终止防止浏览器卡死
// 3. 失控脚本终止：页面脚本死循环、递归溢出时紧急停止
// 4. 任务中断：主动中断正在执行的脚本，切换新任务
// 5. 稳定性保障：长时间运行自动化时，防止恶意/异常脚本挂起进程
// 6. 调试中断：调试过程中手动中断卡死的JS执行

// CDPRuntimeTerminateExecution 立即终止当前正在执行的JavaScript代码
func CDPRuntimeTerminateExecution() (string, error) {
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
		"method": "Runtime.terminateExecution"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Runtime.terminateExecution 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 3 * time.Second // 终止命令响应极快
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

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Runtime.terminateExecution 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础终止卡死JS ====================
func ExampleTerminateExecution_Basic() {
	// 立即终止当前执行的死循环/耗时JS
	resp, err := CDPRuntimeTerminateExecution()
	if err != nil {
		log.Fatalf("终止JS执行失败: %v", err)
	}
	log.Println("JS已成功强制终止，浏览器恢复响应:", resp)
}

// ==================== 使用示例 2：超时自动终止（保护机制） ====================
func ExampleTerminateExecution_Timeout() {
	// 执行耗时JS，5秒超时自动终止
	go func() {
		time.Sleep(5 * time.Second)
		// 超时强制终止
		_, err := CDPRuntimeTerminateExecution()
		if err == nil {
			log.Println("【超时保护】JS执行超时，已自动终止")
		}
	}()

	// 执行可能卡死的JS
	log.Println("开始执行可能超时的JS代码...")
	_, _ = CDPRuntimeEvaluate(`while(true){}`, true, false)
}

// ==================== 使用示例 3：自动化任务紧急中断 ====================
func ExampleTerminateExecution_AutoTest() {
	// 自动化任务异常，需要立即停止所有JS
	log.Println("检测到任务异常，紧急终止脚本执行...")

	resp, err := CDPRuntimeTerminateExecution()
	if err != nil {
		log.Fatalf("任务中断失败: %v", err)
	}

	log.Println("自动化任务已安全中断，浏览器状态正常:", resp)
}

// ==================== 使用示例 4：终止死循环/递归溢出 ====================
func ExampleTerminateExecution_DeadLoop() {
	// 页面陷入死循环，无响应
	log.Println("检测到JS死循环，启动终止...")

	_, err := CDPRuntimeTerminateExecution()
	if err != nil {
		log.Fatalf("死循环终止失败: %v", err)
	}

	log.Println("死循环已终止，浏览器恢复正常")
}

*/
