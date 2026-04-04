package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  DOMDebugger.getEventListeners  -----------------------------------------------
// === 应用场景 ===
// 1. 事件监听器检查: 检查特定DOM节点上已注册的事件监听器
// 2. 内存泄漏检测: 识别可能造成内存泄漏的事件监听器
// 3. 代码调试: 调试复杂的事件处理逻辑
// 4. 性能分析: 分析事件监听器对性能的影响
// 5. 兼容性检查: 检查不同浏览器下事件监听器的差异
// 6. 自动化测试: 验证事件监听器是否正确注册

// CDPDOMDebuggerGetEventListeners 获取指定节点的事件监听器
func CDPDOMDebuggerGetEventListeners(objectID string) (string, error) {
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
        "method": "DOMDebugger.getEventListeners",
        "params": {
            "objectId": "%s"
        }
    }`, reqID, objectID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getEventListeners 请求失败: %w", err)
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
			return "", fmt.Errorf("getEventListeners 请求超时")
		}
	}
}

/*

// === 使用示例: 检查按钮的事件监听器 ===
func exampleGetEventListeners() {
    // 假设我们已经获取了一个按钮的DOM对象ID
    buttonObjectID := "{\"injectedScriptId\":1,\"id\":1}"

    // 获取该按钮上的所有事件监听器
    result, err := CDPDOMDebuggerGetEventListeners(buttonObjectID)
    if err != nil {
        log.Printf("获取事件监听器失败: %v", err)
        return
    }

    log.Printf("按钮事件监听器信息: %s", result)

    // 解析返回的JSON，通常包含如下信息：
    // {
    //     "listeners": [{
    //         "type": "click",
    //         "useCapture": false,
    //         "passive": false,
    //         "once": false,
    //         "handler": {
    //             "className": "Function",
    //             "description": "function() {...}"
    //         }
    //     }]
    // }
}

*/

// -----------------------------------------------  DOMDebugger.removeDOMBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 断点清理: 在调试完成后移除DOM断点
// 2. 动态断点管理: 在运行时根据条件移除断点
// 3. 内存优化: 移除不再需要的DOM断点以释放资源
// 4. 条件调试: 在某些条件下移除断点避免干扰
// 5. 自动化测试: 在测试完成后清理测试断点
// 6. 调试流程控制: 在特定步骤移除断点继续执行

// CDPDOMDebuggerRemoveDOMBreakpoint 移除DOM断点
func CDPDOMDebuggerRemoveDOMBreakpoint(nodeID int, breakpointType string) (string, error) {
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
        "method": "DOMDebugger.removeDOMBreakpoint",
        "params": {
            "nodeId": %d,
            "type": "%s"
        }
    }`, reqID, nodeID, breakpointType)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 removeDOMBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("removeDOMBreakpoint 请求超时")
		}
	}
}

/*

// === 使用示例: 移除元素子节点变更断点 ===
func exampleRemoveDOMBreakpoint() {
    // 假设我们之前在一个节点上设置了断点
    // nodeID 是通过 DOM.getNodeId 或其他方法获取的
    targetNodeID := 45

    // 移除该节点上的子节点变更断点
    // breakpointType 可以是: "subtree-modified", "attribute-modified", "node-removed"
    result, err := CDPDOMDebuggerRemoveDOMBreakpoint(targetNodeID, "subtree-modified")
    if err != nil {
        log.Printf("移除DOM断点失败: %v", err)
        return
    }

    log.Printf("DOM断点移除成功: %s", result)
    // 响应示例: {} (空对象表示成功)

    // 可以在调试完成后清理所有断点类型
    breakpointTypes := []string{"subtree-modified", "attribute-modified", "node-removed"}
    for _, breakpointType := range breakpointTypes {
        if _, err := CDPDOMDebuggerRemoveDOMBreakpoint(targetNodeID, breakpointType); err != nil {
            log.Printf("移除断点类型 %s 失败: %v", breakpointType, err)
        }
    }
}

*/

// -----------------------------------------------  DOMDebugger.removeEventListenerBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 事件调试清理: 移除事件监听器的调试断点
// 2. 性能优化: 移除不再需要的事件断点以提升性能
// 3. 条件调试: 在特定条件下移除事件断点
// 4. 内存管理: 释放事件断点占用的资源
// 5. 测试清理: 在自动化测试完成后清理事件断点
// 6. 调试流程控制: 控制事件断点的生命周期

// CDPDOMDebuggerRemoveEventListenerBreakpoint 移除事件监听器断点
func CDPDOMDebuggerRemoveEventListenerBreakpoint(eventName string) (string, error) {
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
        "method": "DOMDebugger.removeEventListenerBreakpoint",
        "params": {
            "eventName": "%s"
        }
    }`, reqID, eventName)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 removeEventListenerBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("removeEventListenerBreakpoint 请求超时")
		}
	}
}

/*

// === 使用示例: 移除click事件断点 ===
func exampleRemoveEventListenerBreakpoint() {
    // 假设我们之前在调试页面中设置了click事件的断点
    // 现在调试完成，需要移除这个事件断点

    // 移除click事件监听器断点
    result, err := CDPDOMDebuggerRemoveEventListenerBreakpoint("click")
    if err != nil {
        log.Printf("移除事件监听器断点失败: %v", err)
        return
    }

    log.Printf("click事件监听器断点移除成功: %s", result)
    // 响应示例: {} (空对象表示成功)

    // 可以批量移除多个事件断点
    eventsToRemove := []string{
        "click", "mouseover", "keydown", "submit", "change",
    }

    for _, eventName := range eventsToRemove {
        if _, err := CDPDOMDebuggerRemoveEventListenerBreakpoint(eventName); err != nil {
            log.Printf("移除事件断点 %s 失败: %v", eventName, err)
        } else {
            log.Printf("已移除事件断点: %s", eventName)
        }
    }

    // 也可以在特定条件下移除断点
    // 例如：只在事件触发次数达到一定阈值后移除
    clickCounter := 0
    // 模拟条件判断
    if clickCounter > 10 {
        if _, err := CDPDOMDebuggerRemoveEventListenerBreakpoint("click"); err == nil {
            log.Printf("click事件触发超过10次，已移除断点")
        }
    }
}

*/

// -----------------------------------------------  DOMDebugger.removeXHRBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 网络请求调试清理: 移除XHR/Fetch请求的调试断点
// 2. 性能优化: 移除不再需要的网络请求断点
// 3. 条件调试: 在特定条件下移除网络请求断点
// 4. 内存管理: 释放网络请求断点占用的资源
// 5. 测试清理: 在自动化测试完成后清理网络请求断点
// 6. 调试流程控制: 控制网络请求断点的生命周期

// CDPDOMDebuggerRemoveXHRBreakpoint 移除XHR断点
func CDPDOMDebuggerRemoveXHRBreakpoint(url string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	params := fmt.Sprintf(`"id": %d, "method": "DOMDebugger.removeXHRBreakpoint"`, reqID)
	if url != "" {
		params += fmt.Sprintf(`, "params": { "url": "%s" }`, url)
	} else {
		params += `, "params": {}`
	}

	message := fmt.Sprintf(`{ %s }`, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 removeXHRBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("removeXHRBreakpoint 请求超时")
		}
	}
}

/*

// === 使用示例1: 移除特定URL的XHR断点 ===
func exampleRemoveSpecificXHRBreakpoint() {
    // 假设我们之前在调试API请求时设置了断点
    // 现在需要移除特定API端点的XHR断点

    apiURL := "https://api.example.com/users"
    result, err := CDPDOMDebuggerRemoveXHRBreakpoint(apiURL)
    if err != nil {
        log.Printf("移除XHR断点失败: %v", err)
        return
    }

    log.Printf("已移除 %s 的XHR断点: %s", apiURL, result)
    // 响应示例: {} (空对象表示成功)
}

// === 使用示例2: 移除所有XHR断点 ===
func exampleRemoveAllXHRBreakpoints() {
    // 不传递URL参数可以移除所有XHR断点
    result, err := CDPDOMDebuggerRemoveXHRBreakpoint("")
    if err != nil {
        log.Printf("移除所有XHR断点失败: %v", err)
        return
    }

    log.Printf("已移除所有XHR断点: %s", result)
    // 响应示例: {} (空对象表示成功)
}

// === 使用示例3: 批量移除多个API断点 ===
func exampleRemoveMultipleXHRBreakpoints() {
    // 在测试或调试完成后，清理多个API断点
    apiEndpoints := []string{
        "https://api.example.com/users",
        "https://api.example.com/orders",
        "https://api.example.com/products",
        "/api/auth/login",  // 相对路径也可以
    }

    for _, endpoint := range apiEndpoints {
        if _, err := CDPDOMDebuggerRemoveXHRBreakpoint(endpoint); err != nil {
            log.Printf("移除端点 %s 的断点失败: %v", endpoint, err)
        } else {
            log.Printf("已移除端点断点: %s", endpoint)
        }
    }

    // 最后清理所有剩余的断点
    if _, err := CDPDOMDebuggerRemoveXHRBreakpoint(""); err == nil {
        log.Printf("已清理所有剩余的XHR断点")
    }
}

*/

// -----------------------------------------------  DOMDebugger.setDOMBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. DOM变化调试: 在DOM节点发生变化时触发断点
// 2. 属性变更监控: 监控特定元素属性的变化
// 3. 节点删除检测: 检测DOM节点被删除的情况
// 4. 子树变更追踪: 追踪DOM子树的结构变化
// 5. 内存泄漏分析: 分析可能导致内存泄漏的DOM操作
// 6. 性能优化: 定位频繁的DOM操作导致的性能问题

// CDPDOMDebuggerSetDOMBreakpoint 设置DOM断点
func CDPDOMDebuggerSetDOMBreakpoint(nodeID int, breakpointType string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 验证断点类型
	validTypes := map[string]bool{
		"subtree-modified":   true,
		"attribute-modified": true,
		"node-removed":       true,
	}
	if !validTypes[breakpointType] {
		return "", fmt.Errorf("无效的断点类型: %s，可选值: subtree-modified, attribute-modified, node-removed", breakpointType)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "DOMDebugger.setDOMBreakpoint",
        "params": {
            "nodeId": %d,
            "type": "%s"
        }
    }`, reqID, nodeID, breakpointType)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setDOMBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("setDOMBreakpoint 请求超时")
		}
	}
}

/*

// === 使用示例1: 监控元素属性变化 ===
func exampleMonitorAttributeChanges() {
    // 假设我们有一个input元素，想要监控它的属性变化
    // 先通过DOM.getNodeId获取元素的nodeId
    inputNodeID := 42

    // 设置属性变化断点
    result, err := CDPDOMDebuggerSetDOMBreakpoint(inputNodeID, "attribute-modified")
    if err != nil {
        log.Printf("设置DOM属性断点失败: %v", err)
        return
    }

    log.Printf("DOM属性变化断点设置成功: %s", result)
    // 响应示例: {} (空对象表示成功)

    // 当元素的属性（如value、disabled、class等）发生变化时
    // 会触发调试器暂停，可以在DevTools中查看具体的变化
}

// === 使用示例2: 监控节点被删除 ===
func exampleMonitorNodeRemoval() {
    // 监控一个重要的DOM节点是否会被意外删除
    importantNodeID := 73

    // 设置节点删除断点
    result, err := CDPDOMDebuggerSetDOMBreakpoint(importantNodeID, "node-removed")
    if err != nil {
        log.Printf("设置DOM节点删除断点失败: %v", err)
        return
    }

    log.Printf("DOM节点删除断点设置成功: %s", result)

    // 当这个节点被从DOM树中移除时
    // 调试器会暂停执行，可以查看调用栈和状态
}

// === 使用示例3: 监控子树结构变化 ===
func exampleMonitorSubtreeChanges() {
    // 监控一个容器的整个子树结构变化
    containerNodeID := 89

    // 设置子树修改断点
    result, err := CDPDOMDebuggerSetDOMBreakpoint(containerNodeID, "subtree-modified")
    if err != nil {
        log.Printf("设置DOM子树断点失败: %v", err)
        return
    }

    log.Printf("DOM子树修改断点设置成功: %s", result)

    // 当容器的子节点被添加、删除或重新排列时
    // 调试器会暂停，可以检查DOM操作的原因
}

*/

// -----------------------------------------------  DOMDebugger.setEventListenerBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 事件触发调试: 在特定类型事件触发时暂停执行
// 2. 事件处理分析: 分析事件处理函数的执行流程
// 3. 事件传播追踪: 追踪事件的捕获和冒泡阶段
// 4. 事件委托调试: 调试事件委托模式的实现
// 5. 第三方库事件分析: 分析第三方库的事件处理逻辑
// 6. 复杂交互调试: 调试复杂的用户交互事件链

// CDPDOMDebuggerSetEventListenerBreakpoint 设置事件监听器断点
func CDPDOMDebuggerSetEventListenerBreakpoint(eventName string) (string, error) {
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
        "method": "DOMDebugger.setEventListenerBreakpoint",
        "params": {
            "eventName": "%s"
        }
    }`, reqID, eventName)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setEventListenerBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("setEventListenerBreakpoint 请求超时")
		}
	}
}

/*

// === 使用示例1: 调试点击事件 ===
func exampleDebugClickEvents() {
    // 设置click事件的全局断点
    result, err := CDPDOMDebuggerSetEventListenerBreakpoint("click")
    if err != nil {
        log.Printf("设置click事件断点失败: %v", err)
        return
    }

    log.Printf("click事件断点设置成功: %s", result)
    // 响应示例: {} (空对象表示成功)

    // 当页面中任何元素触发click事件时
    // 调试器会在事件处理函数执行前暂停
    // 可以在DevTools中查看事件目标、事件对象和调用栈
}

// === 使用示例2: 调试表单事件 ===
func exampleDebugFormEvents() {
    // 设置表单相关事件的断点
    formEvents := []string{"submit", "change", "input", "focus", "blur"}

    for _, eventName := range formEvents {
        result, err := CDPDOMDebuggerSetEventListenerBreakpoint(eventName)
        if err != nil {
            log.Printf("设置 %s 事件断点失败: %v", eventName, err)
        } else {
            log.Printf("已设置 %s 事件断点", eventName)
        }
    }

    // 这样可以全面调试表单的交互流程
    // 当用户填写表单、提交表单时，调试器会在相应事件处暂停
}

// === 使用示例3: 调试键盘和鼠标事件 ===
func exampleDebugKeyboardMouseEvents() {
    // 设置键盘和鼠标事件的断点
    keyEvents := []string{"keydown", "keyup", "keypress"}
    mouseEvents := []string{"mousedown", "mouseup", "mousemove", "mouseover", "mouseout"}

    allEvents := append(keyEvents, mouseEvents...)

    for _, eventName := range allEvents {
        if _, err := CDPDOMDebuggerSetEventListenerBreakpoint(eventName); err == nil {
            log.Printf("已设置键盘/鼠标事件断点: %s", eventName)
        }
    }

    // 这样可以调试复杂的用户交互
    // 例如：键盘快捷键、拖拽操作、鼠标悬停效果等
}

// === 使用示例4: 调试自定义事件 ===
func exampleDebugCustomEvents() {
    // 设置自定义事件的断点
    customEvents := []string{"customEvent1", "customEvent2", "myapp:event"}

    for _, eventName := range customEvents {
        if _, err := CDPDOMDebuggerSetEventListenerBreakpoint(eventName); err == nil {
            log.Printf("已设置自定义事件断点: %s", eventName)
        } else {
            log.Printf("设置自定义事件断点 %s 失败: %v", eventName, err)
        }
    }

    // 这样可以调试应用中自定义的事件系统
    // 例如：Vue/React等框架的自定义事件
}

*/

// -----------------------------------------------  DOMDebugger.setXHRBreakpoint  -----------------------------------------------
// === 应用场景 ===
// 1. 网络请求调试: 在特定XHR请求发生时暂停执行
// 2. API调用追踪: 追踪特定API端点的调用情况
// 3. 请求参数分析: 分析网络请求的参数和响应
// 4. 错误调试: 调试网络请求失败的问题
// 5. 性能分析: 分析网络请求的性能瓶颈
// 6. 安全分析: 监控敏感API的调用情况

// CDPDOMDebuggerSetXHRBreakpoint 设置XHR断点
func CDPDOMDebuggerSetXHRBreakpoint(urlPattern string) (string, error) {
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
        "method": "DOMDebugger.setXHRBreakpoint",
        "params": {
            "url": "%s"
        }
    }`, reqID, urlPattern)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setXHRBreakpoint 请求失败: %w", err)
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
			return "", fmt.Errorf("setXHRBreakpoint 请求超时")
		}
	}
}

/*

// === 使用示例1: 调试特定API端点 ===
func exampleDebugSpecificAPI() {
    // 设置对特定API端点的XHR断点
    apiURL := "https://api.example.com/users"
    result, err := CDPDOMDebuggerSetXHRBreakpoint(apiURL)
    if err != nil {
        log.Printf("设置API断点失败: %v", err)
        return
    }

    log.Printf("API断点设置成功: %s", result)
    // 响应示例: {} (空对象表示成功)

    // 当页面中发起对 https://api.example.com/users 的XHR或Fetch请求时
    // 调试器会在请求发送前暂停
    // 可以检查请求参数、请求头等信息
}

// === 使用示例2: 使用通配符匹配多个API ===
func exampleDebugMultipleAPIsWithWildcard() {
    // 使用通配符匹配一组相关的API
    // 匹配所有用户相关的API
    result, err := CDPDOMDebuggerSetXHRBreakpoint("https://api.example.com/users/*")
    if err != nil {
        log.Printf("设置通配符断点失败: %v", err)
        return
    }

    log.Printf("通配符API断点设置成功: %s", result)
    // 这会匹配:
    // - https://api.example.com/users/list
    // - https://api.example.com/users/create
    // - https://api.example.com/users/123/update
    // 等所有以users开头的API
}

// === 使用示例3: 调试认证相关的请求 ===
func exampleDebugAuthRequests() {
    // 设置对认证相关API的断点
    authAPIs := []string{
        "/api/auth/login",
        "/api/auth/logout",
        "/api/auth/token",
        "/api/auth/refresh",
    }

    for _, api := range authAPIs {
        result, err := CDPDOMDebuggerSetXHRBreakpoint(api)
        if err != nil {
            log.Printf("设置认证API断点 %s 失败: %v", api, err)
        } else {
            log.Printf("已设置认证API断点: %s", api)
        }
    }

    // 这样可以调试登录、登出、token刷新等认证流程
    // 在请求发送前暂停，可以检查认证参数是否正确
}

// === 使用示例4: 调试上传下载接口 ===
func exampleDebugFileUploadDownload() {
    // 设置文件上传下载API的断点
    fileAPIs := []string{
        "/api/upload",           // 上传接口
        "/api/download/*",       // 下载接口（使用通配符）
        "/api/attachment/",      // 附件接口
    }

    for _, api := range fileAPIs {
        if _, err := CDPDOMDebuggerSetXHRBreakpoint(api); err == nil {
            log.Printf("已设置文件接口断点: %s", api)
        }
    }

    // 这样可以调试文件上传下载的逻辑
    // 包括FormData的处理、文件分片等
}

// === 使用示例5: 调试第三方服务调用 ===
func exampleDebugThirdPartyServices() {
    // 设置对第三方服务的断点
    thirdPartyServices := []string{
        "https://maps.googleapis.com/",   // Google地图API
        "https://analytics.google.com/",  // Google分析
        "https://api.stripe.com/",        // Stripe支付
    }

    for _, service := range thirdPartyServices {
        if _, err := CDPDOMDebuggerSetXHRBreakpoint(service); err == nil {
            log.Printf("已设置第三方服务断点: %s", service)
        }
    }

    // 这样可以调试与第三方服务的集成
    // 了解何时发起请求、请求参数是否正确
}

*/
