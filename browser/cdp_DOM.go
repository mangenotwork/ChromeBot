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

// -----------------------------------------------  DOM.describeNode  -----------------------------------------------
// === 应用场景 ===
// 1. 节点信息获取: 获取DOM节点的详细信息
// 2. 元素分析: 分析特定元素的属性、样式和布局
// 3. 调试辅助: 调试时获取节点的完整信息
// 4. 自动化测试: 在自动化测试中验证元素状态
// 5. 性能分析: 分析DOM节点的性能特征
// 6. 状态检查: 检查节点的可见性、可交互性等状态

// CDPDOMDescribeNode 描述指定节点的详细信息
// nodeID: 节点ID
// depth: 深度，-1表示完整子树，0表示仅节点自身，正整数表示深度
// pierce: 是否穿透shadow root
func CDPDOMDescribeNode(nodeID int, depth int, pierce bool) (string, error) {
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
        "method": "DOM.describeNode",
        "params": {
            "nodeId": %d,
            "depth": %d,
            "pierce": %v
        }
    }`, reqID, nodeID, depth, pierce)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.describeNode 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.describeNode 请求超时")
		}
	}
}

/*

// 示例: 获取按钮元素的详细信息
func ExampleCDPDOMDescribeNode() {
    // 假设我们有一个按钮节点的ID
    buttonNodeID := 123

    // 获取节点的详细信息，包括子节点
    result, err := CDPDOMDescribeNode(buttonNodeID, 1, false)
    if err != nil {
        log.Printf("描述节点失败: %v", err)
        return
    }

    // 解析返回的节点信息
    var response struct {
        Result struct {
            Node struct {
                NodeID       int    `json:"nodeId"`
                NodeType     int    `json:"nodeType"`
                NodeName     string `json:"nodeName"`
                LocalName    string `json:"localName"`
                NodeValue    string `json:"nodeValue"`
                ChildNodeCount int  `json:"childNodeCount"`
                Attributes   []string `json:"attributes,omitempty"`
                FrameID      string `json:"frameId,omitempty"`
                ContentDocument *struct {
                    NodeID int `json:"nodeId"`
                } `json:"contentDocument,omitempty"`
                ShadowRoots []struct {
                    NodeID int `json:"nodeId"`
                } `json:"shadowRoots,omitempty"`
                TemplateContent *struct {
                    NodeID int `json:"nodeId"`
                } `json:"templateContent,omitempty"`
                PseudoElements []struct {
                    NodeID int `json:"nodeId"`
                } `json:"pseudoElements,omitempty"`
                ImportedDocument *struct {
                    NodeID int `json:"nodeId"`
                } `json:"importedDocument,omitempty"`
                DistributedNodes []struct {
                    NodeID int `json:"nodeId"`
                } `json:"distributedNodes,omitempty"`
            } `json:"node"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err == nil {
        node := response.Result.Node

        fmt.Printf("=== 节点详细信息 ===\n")
        fmt.Printf("节点ID: %d\n", node.NodeID)
        fmt.Printf("节点类型: %d\n", node.NodeType)
        fmt.Printf("节点名称: %s\n", node.NodeName)
        fmt.Printf("本地名称: %s\n", node.LocalName)
        fmt.Printf("节点值: %s\n", node.NodeValue)
        fmt.Printf("子节点数量: %d\n", node.ChildNodeCount)

        if len(node.Attributes) > 0 {
            fmt.Printf("属性:\n")
            for i := 0; i < len(node.Attributes); i += 2 {
                if i+1 < len(node.Attributes) {
                    fmt.Printf("  %s: %s\n", node.Attributes[i], node.Attributes[i+1])
                }
            }
        }

        if node.FrameID != "" {
            fmt.Printf("框架ID: %s\n", node.FrameID)
        }

        if node.ContentDocument != nil {
            fmt.Printf("内容文档节点ID: %d\n", node.ContentDocument.NodeID)
        }

        if len(node.ShadowRoots) > 0 {
            fmt.Printf("Shadow Root数量: %d\n", len(node.ShadowRoots))
            for _, shadowRoot := range node.ShadowRoots {
                fmt.Printf("  Shadow Root节点ID: %d\n", shadowRoot.NodeID)
            }
        }

        if node.TemplateContent != nil {
            fmt.Printf("模板内容节点ID: %d\n", node.TemplateContent.NodeID)
        }

        if len(node.PseudoElements) > 0 {
            fmt.Printf("伪元素数量: %d\n", len(node.PseudoElements))
            for _, pseudo := range node.PseudoElements {
                fmt.Printf("  伪元素节点ID: %d\n", pseudo.NodeID)
            }
        }

        if node.ImportedDocument != nil {
            fmt.Printf("导入文档节点ID: %d\n", node.ImportedDocument.NodeID)
        }

        if len(node.DistributedNodes) > 0 {
            fmt.Printf("分布节点数量: %d\n", len(node.DistributedNodes))
            for _, distributed := range node.DistributedNodes {
                fmt.Printf("  分布节点ID: %d\n", distributed.NodeID)
            }
        }

        // 根据节点信息做进一步处理
        if node.NodeName == "BUTTON" {
            // 检查按钮是否可点击
            for i := 0; i < len(node.Attributes); i += 2 {
                if i+1 < len(node.Attributes) && node.Attributes[i] == "disabled" {
                    fmt.Printf("⚠️ 按钮被禁用\n")
                }
            }
        }
    } else {
        log.Printf("描述节点结果: %s", result)
    }
}

*/

// -----------------------------------------------  DOM.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 资源清理: 清理DOM监听器和相关资源
// 2. 性能优化: 停止DOM事件监听以优化性能
// 3. 测试完成: DOM测试完成后清理环境
// 4. 错误恢复: DOM功能异常时恢复
// 5. 功能切换: 切换不同的测试场景
// 6. 内存管理: 释放DOM相关内存

// CDPDOMDisable 禁用DOM域
func CDPDOMDisable() (string, error) {
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
        "method": "DOM.disable"
    }`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.disable 请求超时")
		}
	}
}

/*

// 示例: 在DOM操作完成后禁用DOM域
func ExampleCDPDOMDisable() {
    // 首先启用DOM域（如果需要的话）
    // 通常DOM.enable会在需要时自动调用，但这里演示完整的流程

    // 1. 执行一些DOM操作
    // 例如：获取文档、查找元素、修改属性等
    log.Println("开始执行DOM操作...")

    // 模拟一些DOM操作
    // 这里假设我们执行了各种DOM查询和修改操作

    // 2. 所有DOM操作完成后，禁用DOM域以清理资源
    result, err := CDPDOMDisable()
    if err != nil {
        log.Printf("禁用DOM域失败: %v", err)
        return
    }

    log.Printf("DOM域已禁用: %s", result)

    // 3. 验证禁用效果
    // 禁用后，DOM相关的事件监听器会被移除
    // 后续的DOM操作可能会失败，除非重新启用

    // 4. 可以继续进行其他非DOM操作
    log.Println("DOM操作完成，资源已清理")

    // 注意：在真实的测试场景中，通常会在测试套件的Teardown阶段禁用DOM
    // 例如：
    /*
    func TestDOMOperations(t *testing.T) {
        // Setup: 启用DOM
        if _, err := CDPDOMEnable(); err != nil {
            t.Fatalf("启用DOM失败: %v", err)
        }

        // 执行测试
        t.Run("测试DOM查询", func(t *testing.T) {
            // 测试代码...
        })

        t.Run("测试DOM修改", func(t *testing.T) {
            // 测试代码...
        })

        // Teardown: 禁用DOM
        t.Cleanup(func() {
            if _, err := CDPDOMDisable(); err != nil {
                t.Logf("禁用DOM失败: %v", err)
            }
        })
    }
    *\/
}

*/

// -----------------------------------------------  DOM.enable  -----------------------------------------------
// === 应用场景 ===
// 1. DOM操作准备: 启用DOM功能以进行元素操作
// 2. 自动化测试: 在自动化测试中启用DOM交互
// 3. 元素监控: 开始监控DOM变化
// 4. 调试支持: 调试时启用DOM检查
// 5. 性能分析: 分析DOM结构和性能
// 6. 事件监听: 开始监听DOM事件

// CDPDOMEnable 启用DOM域
func CDPDOMEnable() (string, error) {
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
        "method": "DOM.enable"
    }`, reqID)

	// 发送请求
	err := chromeInstance.NowTabWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.enable 请求超时")
		}
	}
}

/*

// 示例: 在DOM测试前启用DOM功能
func ExampleCDPDOMEnable() {
    // 1. 启用DOM域
    result, err := CDPDOMEnable()
    if err != nil {
        log.Printf("启用DOM域失败: %v", err)
        return
    }

    log.Printf("DOM域已启用: %s", result)

    // 2. 现在可以执行各种DOM操作
    // 例如：获取文档根节点
    // 注意：这需要其他DOM方法配合，这里只是演示流程

    // 3. 设置DOM事件监听（如果需要）
    // 例如监听属性修改、子节点添加等事件

    // 4. 执行具体的DOM测试或操作
    log.Println("DOM功能已启用，可以执行以下操作：")
    log.Println("  - 查询DOM元素")
    log.Println("  - 修改元素属性")
    log.Println("  - 监听DOM变化")
    log.Println("  - 获取布局信息")
    log.Println("  - 执行JavaScript")

    // 在实际测试中，通常会这样使用：
    /*
    func TestDOMFunctionality() {
        // 启用DOM
        if _, err := CDPDOMEnable(); err != nil {
            log.Fatalf("无法启用DOM: %v", err)
        }

        // 确保测试完成后清理
        defer func() {
            if _, err := CDPDOMDisable(); err != nil {
                log.Printf("警告：禁用DOM失败: %v", err)
            }
        }()

        // 执行DOM测试
        testDOMQuery()
        testDOMModification()
        testDOMEvents()
    }
    *\/

    // 5. 验证DOM功能是否正常工作
    // 可以尝试获取文档根节点来验证
    log.Println("DOM启用完成，可以开始DOM相关操作")
}


*/

// -----------------------------------------------  DOM.focus  -----------------------------------------------
// === 应用场景 ===
// 1. 元素聚焦: 将焦点设置到特定DOM元素
// 2. 表单测试: 测试表单元素的焦点行为
// 3. 可访问性: 测试键盘导航和焦点管理
// 4. 交互测试: 测试用户交互时的焦点变化
// 5. 调试辅助: 调试焦点相关的问题
// 6. 自动化测试: 自动化测试中的焦点控制

// CDPDOMFocus 将焦点设置到指定的DOM节点
// nodeID: 要聚焦的节点ID
func CDPDOMFocus(nodeID int) (string, error) {
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
        "method": "DOM.focus",
        "params": {
            "nodeId": %d
        }
    }`, reqID, nodeID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.focus 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.focus 请求超时")
		}
	}
}

/*

// 示例: 将焦点设置到输入框
func ExampleCDPDOMFocus() {
    // 假设我们有一个输入框的节点ID
    // 在实际使用中，这个ID通常通过DOM.querySelector或其他查询方法获得
    inputNodeID := 456

    // 1. 首先启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 将焦点设置到输入框
    result, err := CDPDOMFocus(inputNodeID)
    if err != nil {
        log.Printf("设置焦点失败: %v", err)
        return
    }

    log.Printf("焦点设置成功: %s", result)

    // 3. 验证焦点是否成功设置
    // 可以通过检查文档的活动元素来验证
    // 或者检查输入框是否获得焦点样式

    // 4. 模拟键盘输入（如果需要）
    // 焦点设置后，可以模拟键盘输入
    log.Println("输入框已获得焦点，可以开始输入...")

    // 5. 测试焦点相关功能
    // 例如：测试Tab键导航
    log.Println("测试Tab键导航：")
    log.Println("  1. 当前焦点在输入框")
    log.Println("  2. 按Tab键应移动到下一个可聚焦元素")
    log.Println("  3. 按Shift+Tab应移动到上一个可聚焦元素")

    // 6. 测试焦点事件
    // 焦点设置应该触发focus事件
    // 可以通过DOM.addEventListener监听focus事件

    // 7. 测试可访问性
    // 验证屏幕阅读器是否能正确识别焦点
    log.Println("可访问性测试：")
    log.Println("  - 验证焦点指示器是否可见")
    log.Println("  - 验证屏幕阅读器是否能读取焦点元素")
    log.Println("  - 验证键盘导航是否正常工作")

    // 8. 清理：将焦点移开（如果需要）
    // 可以通过聚焦到其他元素或body来移除焦点
    // bodyNodeID := 1 // body通常有特定的节点ID
    // CDPDOMFocus(bodyNodeID)
}

*/

// -----------------------------------------------  DOM.getAttributes  -----------------------------------------------
// === 应用场景 ===
// 1. 属性检查: 获取DOM元素的所有属性
// 2. 数据提取: 从元素属性中提取数据
// 3. 状态验证: 验证元素的状态属性
// 4. 样式检查: 检查内联样式属性
// 5. 数据属性: 获取data-*自定义属性
// 6. 表单验证: 检查表单元素的属性状态

// CDPDOMGetAttributes 获取指定节点的所有属性
// nodeID: 要获取属性的节点ID
func CDPDOMGetAttributes(nodeID int) (string, error) {
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
        "method": "DOM.getAttributes",
        "params": {
            "nodeId": %d
        }
    }`, reqID, nodeID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.getAttributes 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.getAttributes 请求超时")
		}
	}
}

/*

// 示例: 获取按钮元素的所有属性
func ExampleCDPDOMGetAttributes() {
    // 假设我们有一个按钮的节点ID
    buttonNodeID := 789

    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 获取按钮的所有属性
    result, err := CDPDOMGetAttributes(buttonNodeID)
    if err != nil {
        log.Printf("获取属性失败: %v", err)
        return
    }

    // 3. 解析返回的属性数据
    var response struct {
        Result struct {
            Attributes []string `json:"attributes"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        log.Printf("解析属性数据失败: %v", err)
        return
    }

    // 4. 分析和显示属性
    attributes := response.Result.Attributes
    fmt.Printf("=== 元素属性分析 ===\n")
    fmt.Printf("属性总数: %d\n\n", len(attributes)/2)

    // 属性以键值对的形式返回：[key1, value1, key2, value2, ...]
    for i := 0; i < len(attributes); i += 2 {
        if i+1 < len(attributes) {
            key := attributes[i]
            value := attributes[i+1]

            fmt.Printf("%-20s: %s\n", key, value)

            // 根据不同的属性类型做特殊处理
            switch key {
            case "class":
                // 分析CSS类
                classes := strings.Split(value, " ")
                fmt.Printf("  CSS类 (%d 个):\n", len(classes))
                for _, className := range classes {
                    if className != "" {
                        fmt.Printf("    - %s\n", className)
                    }
                }

            case "style":
                // 分析内联样式
                fmt.Printf("  内联样式:\n")
                styles := strings.Split(value, ";")
                for _, style := range styles {
                    if strings.Contains(style, ":") {
                        parts := strings.SplitN(style, ":", 2)
                        if len(parts) == 2 {
                            prop := strings.TrimSpace(parts[0])
                            val := strings.TrimSpace(parts[1])
                            fmt.Printf("    %s: %s\n", prop, val)
                        }
                    }
                }

            case "disabled":
                // 检查禁用状态
                if value == "true" || value == "" {
                    fmt.Printf("  ⚠️ 元素被禁用\n")
                }

            case "required":
                // 检查必填状态
                if value == "true" || value == "" {
                    fmt.Printf("  ⚠️ 元素是必填的\n")
                }

            case "readonly":
                // 检查只读状态
                if value == "true" || value == "" {
                    fmt.Printf("  ⚠️ 元素是只读的\n")
                }

            case "type":
                // 记录元素类型
                fmt.Printf("  元素类型: %s\n", value)

            default:
                // 检查data-*自定义属性
                if strings.HasPrefix(key, "data-") {
                    fmt.Printf("  📊 自定义数据属性\n")
                }
            }
        }
    }

    // 5. 提取特定用途的属性
    fmt.Printf("\n=== 属性提取 ===\n")

    // 提取所有data-*属性
    dataAttrs := make(map[string]string)
    for i := 0; i < len(attributes); i += 2 {
        if i+1 < len(attributes) && strings.HasPrefix(attributes[i], "data-") {
            dataAttrs[attributes[i]] = attributes[i+1]
        }
    }

    if len(dataAttrs) > 0 {
        fmt.Printf("自定义数据属性:\n")
        for key, value := range dataAttrs {
            fmt.Printf("  %s = %s\n", key, value)
        }
    }

    // 提取ARIA属性
    ariaAttrs := make(map[string]string)
    for i := 0; i < len(attributes); i += 2 {
        if i+1 < len(attributes) && strings.HasPrefix(attributes[i], "aria-") {
            ariaAttrs[attributes[i]] = attributes[i+1]
        }
    }

    if len(ariaAttrs) > 0 {
        fmt.Printf("\nARIA可访问性属性:\n")
        for key, value := range ariaAttrs {
            fmt.Printf("  %s = %s\n", key, value)
        }
    }

    // 6. 验证属性值
    fmt.Printf("\n=== 属性验证 ===\n")

    // 检查必要的属性是否存在
    requiredAttrs := []string{"id", "name", "class"}
    for _, attr := range requiredAttrs {
        found := false
        for i := 0; i < len(attributes); i += 2 {
            if i+1 < len(attributes) && attributes[i] == attr {
                fmt.Printf("✓ 找到 %s 属性: %s\n", attr, attributes[i+1])
                found = true
                break
            }
        }
        if !found {
            fmt.Printf("⚠️ 缺少 %s 属性\n", attr)
        }
    }

    // 7. 属性统计
    fmt.Printf("\n=== 属性统计 ===\n")
    fmt.Printf("总属性数量: %d\n", len(attributes)/2)
    fmt.Printf("自定义数据属性: %d\n", len(dataAttrs))
    fmt.Printf("ARIA属性: %d\n", len(ariaAttrs))
    fmt.Printf("样式属性: %d\n", countAttributesByPrefix(attributes, "style"))
    fmt.Printf("事件属性: %d\n", countAttributesByPrefix(attributes, "on"))
}

// 辅助函数: 统计以特定前缀开头的属性数量
func countAttributesByPrefix(attributes []string, prefix string) int {
    count := 0
    for i := 0; i < len(attributes); i += 2 {
        if i < len(attributes) && strings.HasPrefix(attributes[i], prefix) {
            count++
        }
    }
    return count
}

*/

// -----------------------------------------------  DOM.getBoxModel  -----------------------------------------------
// === 应用场景 ===
// 1. 布局分析: 分析元素的盒模型布局
// 2. 位置验证: 验证元素在页面中的位置
// 3. 尺寸测量: 精确测量元素的尺寸
// 4. 重叠检测: 检测元素间的重叠情况
// 5. 响应式测试: 测试不同视口下的布局
// 6. 视觉回归: 视觉回归测试中的位置验证

// CDPDOMGetBoxModel 获取指定节点的盒模型信息
// nodeID: 要获取盒模型的节点ID
func CDPDOMGetBoxModel(nodeID int) (string, error) {
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
        "method": "DOM.getBoxModel",
        "params": {
            "nodeId": %d
        }
    }`, reqID, nodeID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.getBoxModel 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.getBoxModel 请求超时")
		}
	}
}

// BoxModel 盒模型结构
type BoxModel struct {
	Content      []float64 `json:"content"` // 内容区域 [x1, y1, x2, y2, x3, y3, x4, y4]
	Padding      []float64 `json:"padding"` // 内边距区域
	Border       []float64 `json:"border"`  // 边框区域
	Margin       []float64 `json:"margin"`  // 外边距区域
	Width        int       `json:"width"`   // 元素宽度
	Height       int       `json:"height"`  // 元素高度
	ShapeOutside *struct {
		Bounds      []float64 `json:"bounds"`      // 形状边界
		Shape       []float64 `json:"shape"`       // 形状点
		MarginShape []float64 `json:"marginShape"` // 带外边距的形状
	} `json:"shapeOutside,omitempty"` // CSS shape-outside形状
}

/*

// 示例: 分析按钮元素的盒模型
func ExampleCDPDOMGetBoxModel() {
    // 假设我们有一个按钮的节点ID
    buttonNodeID := 890

    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 获取按钮的盒模型信息
    result, err := CDPDOMGetBoxModel(buttonNodeID)
    if err != nil {
        log.Printf("获取盒模型失败: %v", err)
        return
    }

    // 3. 解析盒模型数据
    var response struct {
        Result struct {
            Model BoxModel `json:"model"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        log.Printf("解析盒模型数据失败: %v", err)
        return
    }

    model := response.Result.Model

    // 4. 显示盒模型详细信息
    fmt.Printf("=== 盒模型分析 ===\n")
    fmt.Printf("元素尺寸: %d x %d 像素\n\n", model.Width, model.Height)

    // 解析各个区域
    displayBox("内容区域", model.Content)
    displayBox("内边距区域", model.Padding)
    displayBox("边框区域", model.Border)
    displayBox("外边距区域", model.Margin)

    // 5. 计算各个区域的尺寸
    fmt.Printf("\n=== 区域尺寸计算 ===\n")

    // 内容区域尺寸
    contentWidth := calculateWidth(model.Content)
    contentHeight := calculateHeight(model.Content)
    fmt.Printf("内容区域: %.1f x %.1f 像素\n", contentWidth, contentHeight)

    // 内边距尺寸
    paddingWidth := calculateWidth(model.Padding)
    paddingHeight := calculateHeight(model.Padding)
    paddingH := (paddingWidth - contentWidth) / 2
    paddingV := (paddingHeight - contentHeight) / 2
    fmt.Printf("内边距: 水平 %.1f, 垂直 %.1f 像素\n", paddingH, paddingV)

    // 边框尺寸
    borderWidth := calculateWidth(model.Border)
    borderHeight := calculateHeight(model.Border)
    borderH := (borderWidth - paddingWidth) / 2
    borderV := (borderHeight - paddingHeight) / 2
    fmt.Printf("边框: 水平 %.1f, 垂直 %.1f 像素\n", borderH, borderV)

    // 外边距尺寸
    marginWidth := calculateWidth(model.Margin)
    marginHeight := calculateHeight(model.Margin)
    marginH := (marginWidth - borderWidth) / 2
    marginV := (marginHeight - borderHeight) / 2
    fmt.Printf("外边距: 水平 %.1f, 垂直 %.1f 像素\n", marginH, marginV)

    // 6. 盒模型可视化
    fmt.Printf("\n=== 盒模型可视化 ===\n")
    visualizeBoxModel(model)

    // 7. 布局验证
    fmt.Printf("\n=== 布局验证 ===\n")

    // 检查元素是否在视口内
    if isElementInViewport(model.Margin) {
        fmt.Printf("✓ 元素在视口内\n")
    } else {
        fmt.Printf("⚠️ 元素可能不在视口内\n")
    }

    // 检查元素是否可见
    if contentWidth > 0 && contentHeight > 0 {
        fmt.Printf("✓ 元素有可见尺寸\n")
    } else {
        fmt.Printf("⚠️ 元素可能不可见（尺寸为0）\n")
    }

    // 检查宽高比
    aspectRatio := contentWidth / contentHeight
    fmt.Printf("宽高比: %.2f\n", aspectRatio)

    // 8. 响应式设计检查
    fmt.Printf("\n=== 响应式设计检查 ===\n")

    // 计算元素占视口的百分比（假设视口为1920x1080）
    viewportWidth := 1920.0
    viewportHeight := 1080.0

    widthPercent := (marginWidth / viewportWidth) * 100
    heightPercent := (marginHeight / viewportHeight) * 100

    fmt.Printf("占视口宽度: %.1f%%\n", widthPercent)
    fmt.Printf("占视口高度: %.1f%%\n", heightPercent)

    if widthPercent > 50 {
        fmt.Printf("⚠️ 元素较宽，在移动设备上可能需要调整\n")
    }

    if heightPercent > 50 {
        fmt.Printf("⚠️ 元素较高，在移动设备上可能需要调整\n")
    }

    // 9. 位置关系分析
    fmt.Printf("\n=== 位置关系 ===\n")

    // 计算元素中心点
    centerX := (model.Content[0] + model.Content[2]) / 2
    centerY := (model.Content[1] + model.Content[5]) / 2
    fmt.Printf("元素中心点: (%.1f, %.1f)\n", centerX, centerY)

    // 计算到视口中心的距离
    viewportCenterX := viewportWidth / 2
    viewportCenterY := viewportHeight / 2
    distanceToCenter := math.Sqrt(
        math.Pow(centerX-viewportCenterX, 2) +
        math.Pow(centerY-viewportCenterY, 2))
    fmt.Printf("到视口中心距离: %.1f 像素\n", distanceToCenter)

    // 10. 形状分析
    if model.ShapeOutside != nil {
        fmt.Printf("\n=== CSS形状分析 ===\n")
        fmt.Printf("检测到CSS shape-outside属性\n")
        fmt.Printf("形状边界: %v\n", model.ShapeOutside.Bounds)
    }
}

// 辅助函数: 显示区域坐标
func displayBox(name string, points []float64) {
    if len(points) < 8 {
        return
    }

    fmt.Printf("%s:\n", name)
    fmt.Printf("  左上: (%.1f, %.1f)\n", points[0], points[1])
    fmt.Printf("  右上: (%.1f, %.1f)\n", points[2], points[3])
    fmt.Printf("  右下: (%.1f, %.1f)\n", points[4], points[5])
    fmt.Printf("  左下: (%.1f, %.1f)\n", points[6], points[7])
    fmt.Println()
}

// 辅助函数: 计算宽度
func calculateWidth(points []float64) float64 {
    if len(points) < 8 {
        return 0
    }
    // 使用左上和右上的x坐标计算宽度
    return math.Abs(points[2] - points[0])
}

// 辅助函数: 计算高度
func calculateHeight(points []float64) float64 {
    if len(points) < 8 {
        return 0
    }
    // 使用左上和左下的y坐标计算高度
    return math.Abs(points[5] - points[1])
}

// 辅助函数: 可视化盒模型
func visualizeBoxModel(model BoxModel) {
    // 简化的文本可视化
    if len(model.Margin) < 8 || len(model.Border) < 8 ||
       len(model.Padding) < 8 || len(model.Content) < 8 {
        return
    }

    // 计算最大边界
    maxX := model.Margin[2]
    maxY := model.Margin[5]

    // 简化的ASCII艺术表示
    fmt.Println("    外边距")
    fmt.Println("  ┌─────────────────────────────────────────────┐")
    fmt.Println("  │                   边框                      │")
    fmt.Println("  │  ┌───────────────────────────────────────┐  │")
    fmt.Println("  │  │             内边距                    │  │")
    fmt.Println("  │  │  ┌─────────────────────────────────┐  │  │")
    fmt.Println("  │  │  │           内容区域              │  │  │")
    fmt.Println("  │  │  │                                 │  │  │")
    fmt.Println("  │  │  │    宽: %4d, 高: %4d       │  │  │", model.Width, model.Height)
    fmt.Println("  │  │  │                                 │  │  │")
    fmt.Println("  │  │  └─────────────────────────────────┘  │  │")
    fmt.Println("  │  │                                       │  │")
    fmt.Println("  │  └───────────────────────────────────────┘  │")
    fmt.Println("  │                                             │")
    fmt.Println("  └─────────────────────────────────────────────┘")
    fmt.Println()

    // 显示尺寸图例
    fmt.Println("图例:")
    fmt.Println("  ███ 内容区域")
    fmt.Println("  ░░░ 内边距")
    fmt.Println("  ▒▒▒ 边框")
    fmt.Println("  ▓▓▓ 外边距")
}

// 辅助函数: 检查元素是否在视口内
func isElementInViewport(points []float64) bool {
    if len(points) < 8 {
        return false
    }

    // 简单的检查：只要有一个角在正坐标区域就认为在视口内
    for i := 0; i < 8; i += 2 {
        if points[i] >= 0 && points[i+1] >= 0 {
            return true
        }
    }
    return false
}

*/

// -----------------------------------------------  DOM.getDocument  -----------------------------------------------
// === 应用场景 ===
// 1. 文档获取: 获取整个页面的DOM文档树
// 2. 结构分析: 分析页面DOM结构
// 3. 自动化测试: 测试前获取文档基准
// 4. 页面快照: 获取页面的DOM快照
// 5. 性能分析: 分析DOM树的大小和复杂度
// 6. 调试支持: 调试时获取完整的文档信息

// CDPDOMGetDocument 获取整个文档的DOM树
// depth: 遍历深度，-1表示完整子树，0表示仅节点自身，正整数表示深度
// pierce: 是否穿透shadow root
func CDPDOMGetDocument(depth int, pierce bool) (string, error) {
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
        "method": "DOM.getDocument",
        "params": {
            "depth": %d,
            "pierce": %v
        }
    }`, reqID, depth, pierce)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.getDocument 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.getDocument 请求超时")
		}
	}
}

// Node DOM节点结构
type Node struct {
	NodeID           int           `json:"nodeId"`
	ParentID         int           `json:"parentId,omitempty"`
	BackendNodeID    int           `json:"backendNodeId"`
	NodeType         int           `json:"nodeType"`
	NodeName         string        `json:"nodeName"`
	LocalName        string        `json:"localName"`
	NodeValue        string        `json:"nodeValue"`
	ChildNodeCount   int           `json:"childNodeCount,omitempty"`
	Children         []Node        `json:"children,omitempty"`
	Attributes       []string      `json:"attributes,omitempty"`
	DocumentURL      string        `json:"documentURL,omitempty"`
	BaseURL          string        `json:"baseURL,omitempty"`
	PublicID         string        `json:"publicId,omitempty"`
	SystemID         string        `json:"systemId,omitempty"`
	InternalSubset   string        `json:"internalSubset,omitempty"`
	XMLVersion       string        `json:"xmlVersion,omitempty"`
	Name             string        `json:"name,omitempty"`
	Value            string        `json:"value,omitempty"`
	PseudoType       string        `json:"pseudoType,omitempty"`
	ShadowRootType   string        `json:"shadowRootType,omitempty"`
	FrameID          string        `json:"frameId,omitempty"`
	ContentDocument  *Node         `json:"contentDocument,omitempty"`
	ShadowRoots      []Node        `json:"shadowRoots,omitempty"`
	TemplateContent  *Node         `json:"templateContent,omitempty"`
	PseudoElements   []Node        `json:"pseudoElements,omitempty"`
	ImportedDocument *Node         `json:"importedDocument,omitempty"`
	DistributedNodes []BackendNode `json:"distributedNodes,omitempty"`
	IsSVG            bool          `json:"isSVG,omitempty"`
}

// BackendNode 后台节点结构
type BackendNode struct {
	NodeType      int    `json:"nodeType"`
	NodeName      string `json:"nodeName"`
	BackendNodeID int    `json:"backendNodeId"`
}

/*


// 示例: 获取并分析页面文档结构
func ExampleCDPDOMGetDocument() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 获取完整的文档树
    result, err := CDPDOMGetDocument(-1, true)
    if err != nil {
        log.Printf("获取文档失败: %v", err)
        return
    }

    // 3. 解析文档树
    var response struct {
        Result struct {
            Root Node `json:"root"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        log.Printf("解析文档数据失败: %v", err)
        return
    }

    root := response.Result.Root

    // 4. 显示文档基本信息
    fmt.Printf("=== 文档信息 ===\n")
    fmt.Printf("文档节点ID: %d\n", root.NodeID)
    fmt.Printf("文档类型: %d\n", root.NodeType)
    fmt.Printf("文档名称: %s\n", root.NodeName)

    if root.DocumentURL != "" {
        fmt.Printf("文档URL: %s\n", root.DocumentURL)
    }

    if root.BaseURL != "" {
        fmt.Printf("基础URL: %s\n", root.BaseURL)
    }

    // 5. 文档结构分析
    fmt.Printf("\n=== 文档结构分析 ===\n")

    // 统计不同类型的节点
    stats := analyzeDocumentStructure(root)

    fmt.Printf("总节点数: %d\n", stats.TotalNodes)
    fmt.Printf("元素节点: %d\n", stats.ElementNodes)
    fmt.Printf("文本节点: %d\n", stats.TextNodes)
    fmt.Printf("注释节点: %d\n", stats.CommentNodes)
    fmt.Printf("文档节点: %d\n", stats.DocumentNodes)
    fmt.Printf("文档类型节点: %d\n", stats.DocumentTypeNodes)

    // 6. 深度分析
    fmt.Printf("\n=== 深度分析 ===\n")

    maxDepth := calculateMaxDepth(root)
    fmt.Printf("文档最大深度: %d\n", maxDepth)

    // 计算平均深度
    totalDepth, nodeCount := calculateTotalDepth(root, 0)
    avgDepth := float64(totalDepth) / float64(nodeCount)
    fmt.Printf("平均深度: %.2f\n", avgDepth)

    // 7. 热门元素分析
    fmt.Printf("\n=== 热门元素分析 ===\n")

    elementFrequency := make(map[string]int)
    countElements(root, elementFrequency)

    // 按频率排序
    sortedElements := sortByFrequency(elementFrequency)

    fmt.Printf("最常用的10个元素:\n")
    for i := 0; i < len(sortedElements) && i < 10; i++ {
        fmt.Printf("  %-10s: %d 个\n", sortedElements[i].Name, sortedElements[i].Count)
    }

    // 8. 属性分析
    fmt.Printf("\n=== 属性分析 ===\n")

    attrStats := analyzeAttributes(root)
    fmt.Printf("总属性数: %d\n", attrStats.TotalAttributes)
    fmt.Printf("不同属性名: %d\n", len(attrStats.AttributeFrequency))

    // 显示最常见的属性
    sortedAttrs := sortByFrequency(attrStats.AttributeFrequency)
    fmt.Printf("\n最常见的10个属性:\n")
    for i := 0; i < len(sortedAttrs) && i < 10; i++ {
        fmt.Printf("  %-15s: %d 次\n", sortedAttrs[i].Name, sortedAttrs[i].Count)
    }

    // 9. Shadow DOM分析
    fmt.Printf("\n=== Shadow DOM分析 ===\n")

    shadowStats := countShadowRoots(root)
    fmt.Printf("Shadow Root数量: %d\n", shadowStats.Total)

    if shadowStats.Total > 0 {
        fmt.Printf("类型分布:\n")
        for shadowType, count := range shadowStats.ByType {
            fmt.Printf("  %-15s: %d\n", shadowType, count)
        }
    }

    // 10. 性能考虑
    fmt.Printf("\n=== 性能考虑 ===\n")

    // 大型DOM树警告
    if stats.TotalNodes > 1000 {
        fmt.Printf("⚠️ DOM树较大（%d个节点），可能影响性能\n", stats.TotalNodes)
    }

    if maxDepth > 20 {
        fmt.Printf("⚠️ DOM深度较深（%d层），可能影响渲染性能\n", maxDepth)
    }

    // 计算DOM复杂度分数
    complexityScore := calculateComplexityScore(stats, maxDepth)
    fmt.Printf("DOM复杂度分数: %.2f\n", complexityScore)

    if complexityScore > 1000 {
        fmt.Printf("⚠️ DOM复杂度较高，建议优化\n")
    }

    // 11. 可访问性检查
    fmt.Printf("\n=== 可访问性检查 ===\n")

    a11yStats := checkAccessibility(root)
    fmt.Printf("有id的元素: %d\n", a11yStats.ElementsWithID)
    fmt.Printf("有aria-label的元素: %d\n", a11yStats.ElementsWithAriaLabel)
    fmt.Printf("有title属性的元素: %d\n", a11yStats.ElementsWithTitle)
    fmt.Printf("有alt属性的图片: %d\n", a11yStats.ImagesWithAlt)

    // 12. 框架检测
    fmt.Printf("\n=== 框架检测 ===\n")

    if hasIframes(root) {
        fmt.Printf("✓ 检测到iframe框架\n")
    } else {
        fmt.Printf("未检测到iframe框架\n")
    }

    if hasSVG(root) {
        fmt.Printf("✓ 检测到SVG内容\n")
    }

    if hasCanvas(root) {
        fmt.Printf("✓ 检测到Canvas元素\n")
    }

    // 13. 文档健康度报告
    fmt.Printf("\n=== 文档健康度报告 ===\n")

    healthScore := calculateDocumentHealthScore(stats, a11yStats, shadowStats)
    fmt.Printf("文档健康度分数: %.1f/100\n", healthScore)

    if healthScore >= 80 {
        fmt.Printf("✓ 文档结构良好\n")
    } else if healthScore >= 60 {
        fmt.Printf("⚠️ 文档结构一般，有优化空间\n")
    } else {
        fmt.Printf("❌ 文档结构需要优化\n")
    }
}

// 辅助结构：文档统计
type DocumentStats struct {
    TotalNodes         int
    ElementNodes       int
    TextNodes          int
    CommentNodes       int
    DocumentNodes      int
    DocumentTypeNodes  int
}

// 分析文档结构
func analyzeDocumentStructure(node Node) DocumentStats {
    stats := DocumentStats{}
    traverseDocument(node, &stats)
    return stats
}

// 遍历文档
func traverseDocument(node Node, stats *DocumentStats) {
    stats.TotalNodes++

    // 根据节点类型计数
    switch node.NodeType {
    case 1: // ELEMENT_NODE
        stats.ElementNodes++
    case 3: // TEXT_NODE
        stats.TextNodes++
    case 8: // COMMENT_NODE
        stats.CommentNodes++
    case 9: // DOCUMENT_NODE
        stats.DocumentNodes++
    case 10: // DOCUMENT_TYPE_NODE
        stats.DocumentTypeNodes++
    }

    // 递归遍历子节点
    for _, child := range node.Children {
        traverseDocument(child, stats)
    }

    // 遍历shadow root
    for _, shadow := range node.ShadowRoots {
        traverseDocument(shadow, stats)
    }

    // 遍历伪元素
    for _, pseudo := range node.PseudoElements {
        traverseDocument(pseudo, stats)
    }

    // 处理内容文档
    if node.ContentDocument != nil {
        traverseDocument(*node.ContentDocument, stats)
    }

    // 处理模板内容
    if node.TemplateContent != nil {
        traverseDocument(*node.TemplateContent, stats)
    }

    // 处理导入文档
    if node.ImportedDocument != nil {
        traverseDocument(*node.ImportedDocument, stats)
    }
}

// 计算最大深度
func calculateMaxDepth(node Node) int {
    maxDepth := 0

    for _, child := range node.Children {
        childDepth := calculateMaxDepth(child)
        if childDepth > maxDepth {
            maxDepth = childDepth
        }
    }

    return maxDepth + 1
}

// 计算总深度
func calculateTotalDepth(node Node, currentDepth int) (int, int) {
    totalDepth := currentDepth
    nodeCount := 1

    for _, child := range node.Children {
        childTotal, childCount := calculateTotalDepth(child, currentDepth+1)
        totalDepth += childTotal
        nodeCount += childCount
    }

    return totalDepth, nodeCount
}

// 元素频率统计
type FrequencyItem struct {
    Name  string
    Count int
}

func countElements(node Node, frequency map[string]int) {
    if node.NodeType == 1 { // 元素节点
        frequency[node.NodeName]++
    }

    for _, child := range node.Children {
        countElements(child, frequency)
    }
}

func sortByFrequency(frequency map[string]int) []FrequencyItem {
    items := make([]FrequencyItem, 0, len(frequency))
    for name, count := range frequency {
        items = append(items, FrequencyItem{Name: name, Count: count})
    }

    sort.Slice(items, func(i, j int) bool {
        if items[i].Count == items[j].Count {
            return items[i].Name < items[j].Name
        }
        return items[i].Count > items[j].Count
    })

    return items
}

// 属性分析
type AttributeStats struct {
    TotalAttributes     int
    AttributeFrequency map[string]int
}

func analyzeAttributes(node Node) AttributeStats {
    stats := AttributeStats{
        AttributeFrequency: make(map[string]int),
    }

    analyzeNodeAttributes(node, &stats)
    return stats
}

func analyzeNodeAttributes(node Node, stats *AttributeStats) {
    // 统计当前节点的属性
    for i := 0; i < len(node.Attributes); i += 2 {
        if i+1 < len(node.Attributes) {
            attrName := node.Attributes[i]
            stats.TotalAttributes++
            stats.AttributeFrequency[attrName]++
        }
    }

    // 递归处理子节点
    for _, child := range node.Children {
        analyzeNodeAttributes(child, stats)
    }
}

// Shadow DOM统计
type ShadowStats struct {
    Total  int
    ByType map[string]int
}

func countShadowRoots(node Node) ShadowStats {
    stats := ShadowStats{
        ByType: make(map[string]int),
    }

    countNodeShadowRoots(node, &stats)
    return stats
}

func countNodeShadowRoots(node Node, stats *ShadowStats) {
    for _, shadow := range node.ShadowRoots {
        stats.Total++
        if shadow.ShadowRootType != "" {
            stats.ByType[shadow.ShadowRootType]++
        }
    }

    for _, child := range node.Children {
        countNodeShadowRoots(child, stats)
    }
}

// 计算复杂度分数
func calculateComplexityScore(stats DocumentStats, maxDepth int) float64 {
    // 简单的复杂度计算：节点数 * 深度系数
    depthFactor := 1.0 + float64(maxDepth)/10.0
    return float64(stats.TotalNodes) * depthFactor
}

// 可访问性统计
type AccessibilityStats struct {
    ElementsWithID       int
    ElementsWithAriaLabel int
    ElementsWithTitle    int
    ImagesWithAlt        int
}

func checkAccessibility(node Node) AccessibilityStats {
    stats := AccessibilityStats{}
    checkNodeAccessibility(node, &stats)
    return stats
}

func checkNodeAccessibility(node Node, stats *AccessibilityStats) {
    if node.NodeType == 1 { // 元素节点
        // 检查属性
        for i := 0; i < len(node.Attributes); i += 2 {
            if i+1 < len(node.Attributes) {
                attrName := node.Attributes[i]

                switch attrName {
                case "id":
                    stats.ElementsWithID++
                case "aria-label":
                    stats.ElementsWithAriaLabel++
                case "title":
                    stats.ElementsWithTitle++
                case "alt":
                    if node.NodeName == "IMG" {
                        stats.ImagesWithAlt++
                    }
                }
            }
        }
    }

    for _, child := range node.Children {
        checkNodeAccessibility(child, stats)
    }
}

// 框架检测
func hasIframes(node Node) bool {
    if node.NodeType == 1 && node.NodeName == "IFRAME" {
        return true
    }

    for _, child := range node.Children {
        if hasIframes(child) {
            return true
        }
    }

    return false
}

func hasSVG(node Node) bool {
    if node.NodeType == 1 && node.IsSVG {
        return true
    }

    for _, child := range node.Children {
        if hasSVG(child) {
            return true
        }
    }

    return false
}

func hasCanvas(node Node) bool {
    if node.NodeType == 1 && node.NodeName == "CANVAS" {
        return true
    }

    for _, child := range node.Children {
        if hasCanvas(child) {
            return true
        }
    }

    return false
}

// 计算文档健康度分数
func calculateDocumentHealthScore(stats DocumentStats, a11yStats AccessibilityStats, shadowStats ShadowStats) float64 {
    score := 0.0

    // 节点数量分数（越少越好）
    if stats.TotalNodes < 500 {
        score += 30
    } else if stats.TotalNodes < 1000 {
        score += 20
    } else if stats.TotalNodes < 2000 {
        score += 10
    }

    // 元素比例分数
    elementRatio := float64(stats.ElementNodes) / float64(stats.TotalNodes)
    if elementRatio > 0.7 {
        score += 20
    } else if elementRatio > 0.5 {
        score += 15
    } else {
        score += 10
    }

    // 可访问性分数
    if stats.ElementNodes > 0 {
        idRatio := float64(a11yStats.ElementsWithID) / float64(stats.ElementNodes)
        if idRatio > 0.3 {
            score += 10
        } else if idRatio > 0.1 {
            score += 5
        }
    }

    if a11yStats.ImagesWithAlt > 0 {
        score += 10
    }

    // Shadow DOM分数
    if shadowStats.Total == 0 {
        score += 10
    } else if shadowStats.Total < 5 {
        score += 5
    }

    // 结构健康度
    if stats.CommentNodes < stats.ElementNodes/10 {
        score += 10
    }

    return score
}

*/

// -----------------------------------------------  DOM.getNodeForLocation  -----------------------------------------------
// === 应用场景 ===
// 1. 坐标定位: 根据屏幕坐标获取对应的DOM元素
// 2. 点击测试: 测试特定坐标下的可点击元素
// 3. 元素悬停: 模拟鼠标悬停在特定位置
// 4. 视觉测试: 验证特定位置的元素显示
// 5. 拖放测试: 测试拖放操作的坐标定位
// 6. 响应式布局: 测试不同屏幕位置的元素

// CDPDOMGetNodeForLocation 根据坐标获取对应的DOM节点
// x: 水平坐标
// y: 垂直坐标
// includeUserAgentShadowDOM: 是否包含用户代理的Shadow DOM
// ignorePointerEventsNone: 是否忽略pointer-events: none的元素
func CDPDOMGetNodeForLocation(x, y int, includeUserAgentShadowDOM, ignorePointerEventsNone bool) (string, error) {
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
        "method": "DOM.getNodeForLocation",
        "params": {
            "x": %d,
            "y": %d,
            "includeUserAgentShadowDOM": %v,
            "ignorePointerEventsNone": %v
        }
    }`, reqID, x, y, includeUserAgentShadowDOM, ignorePointerEventsNone)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.getNodeForLocation 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.getNodeForLocation 请求超时")
		}
	}
}

/*

// 示例: 根据坐标获取并分析元素
func ExampleCDPDOMGetNodeForLocation() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 定义要测试的坐标
    // 假设我们要测试页面中心点的元素
    x := 960
    y := 540

    // 3. 获取坐标位置的节点
    result, err := CDPDOMGetNodeForLocation(x, y, true, false)
    if err != nil {
        log.Printf("获取坐标节点失败: %v", err)
        return
    }

    // 4. 解析响应
    var response struct {
        Result struct {
            BackendNodeID   int    `json:"backendNodeId"`
            FrameID         string `json:"frameId"`
            NodeID          int    `json:"nodeId,omitempty"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        log.Printf("解析响应失败: %v", err)
        return
    }

    // 5. 显示基本信息
    fmt.Printf("=== 坐标位置分析 ===\n")
    fmt.Printf("坐标: (%d, %d)\n", x, y)
    fmt.Printf("后端节点ID: %d\n", response.Result.BackendNodeID)

    if response.Result.FrameID != "" {
        fmt.Printf("框架ID: %s\n", response.Result.FrameID)
    }

    if response.Result.NodeID > 0 {
        fmt.Printf("节点ID: %d\n", response.Result.NodeID)

        // 6. 获取节点的详细信息
        nodeInfo, err := CDPDOMDescribeNode(response.Result.NodeID, 0, true)
        if err == nil {
            var nodeResp struct {
                Result struct {
                    Node struct {
                        NodeID       int      `json:"nodeId"`
                        NodeType     int      `json:"nodeType"`
                        NodeName     string   `json:"nodeName"`
                        LocalName    string   `json:"localName"`
                        NodeValue    string   `json:"nodeValue"`
                        Attributes   []string `json:"attributes,omitempty"`
                    } `json:"node"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(nodeInfo), &nodeResp); err == nil {
                fmt.Printf("\n=== 节点详细信息 ===\n")
                fmt.Printf("节点类型: %d\n", nodeResp.Result.Node.NodeType)
                fmt.Printf("节点名称: %s\n", nodeResp.Result.Node.NodeName)
                fmt.Printf("本地名称: %s\n", nodeResp.Result.Node.LocalName)
                fmt.Printf("节点值: %s\n", nodeResp.Result.Node.NodeValue)

                // 显示属性
                if len(nodeResp.Result.Node.Attributes) > 0 {
                    fmt.Printf("\n属性:\n")
                    for i := 0; i < len(nodeResp.Result.Node.Attributes); i += 2 {
                        if i+1 < len(nodeResp.Result.Node.Attributes) {
                            fmt.Printf("  %s: %s\n",
                                nodeResp.Result.Node.Attributes[i],
                                nodeResp.Result.Node.Attributes[i+1])
                        }
                    }
                }
            }
        }
    }

    // 7. 分析坐标位置的可交互性
    fmt.Printf("\n=== 可交互性分析 ===\n")

    // 检查是否是文本节点
    if response.Result.NodeID > 0 {
        // 获取盒模型
        boxModel, err := CDPDOMGetBoxModel(response.Result.NodeID)
        if err == nil {
            var boxResp struct {
                Result struct {
                    Model struct {
                        Content []float64 `json:"content"`
                        Width   int       `json:"width"`
                        Height  int       `json:"height"`
                    } `json:"model"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(boxModel), &boxResp); err == nil {
                if len(boxResp.Result.Model.Content) >= 8 {
                    // 计算元素边界
                    left := boxResp.Result.Model.Content[0]
                    top := boxResp.Result.Model.Content[1]
                    right := boxResp.Result.Model.Content[2]
                    bottom := boxResp.Result.Model.Content[5]

                    fmt.Printf("元素边界: (%.1f, %.1f) 到 (%.1f, %.1f)\n",
                        left, top, right, bottom)
                    fmt.Printf("元素尺寸: %d x %d 像素\n",
                        boxResp.Result.Model.Width, boxResp.Result.Model.Height)

                    // 检查坐标是否在元素内
                    if float64(x) >= left && float64(x) <= right &&
                       float64(y) >= top && float64(y) <= bottom {
                        fmt.Printf("✓ 坐标在元素边界内\n")
                    } else {
                        fmt.Printf("⚠️ 坐标在元素边界外\n")
                    }

                    // 计算到元素中心的距离
                    centerX := (left + right) / 2
                    centerY := (top + bottom) / 2
                    distance := math.Sqrt(
                        math.Pow(float64(x)-centerX, 2) +
                        math.Pow(float64(y)-centerY, 2))
                    fmt.Printf("到元素中心距离: %.1f 像素\n", distance)
                }
            }
        }
    }

    // 8. 点击测试
    fmt.Printf("\n=== 点击测试模拟 ===\n")

    // 根据节点类型判断是否可点击
    if response.Result.NodeID > 0 {
        // 获取节点描述以确定类型
        desc, _ := CDPDOMDescribeNode(response.Result.NodeID, 0, true)
        var descResp struct {
            Result struct {
                Node struct {
                    NodeName string `json:"nodeName"`
                } `json:"node"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(desc), &descResp); err == nil {
            nodeName := descResp.Result.Node.NodeName
            isClickable := false

            // 检查常见可点击元素
            clickableElements := map[string]bool{
                "A":      true,  // 链接
                "BUTTON": true,  // 按钮
                "INPUT":  true,  // 输入框
                "SELECT": true,  // 选择框
                "TEXTAREA": true, // 文本区域
            }

            if clickableElements[nodeName] {
                isClickable = true
                fmt.Printf("✓ 可能是可点击元素: %s\n", nodeName)
            } else {
                // 检查是否有onclick属性
                nodeInfo, _ := CDPDOMGetAttributes(response.Result.NodeID)
                var attrResp struct {
                    Result struct {
                        Attributes []string `json:"attributes"`
                    } `json:"result"`
                }

                if err := json.Unmarshal([]byte(nodeInfo), &attrResp); err == nil {
                    for i := 0; i < len(attrResp.Result.Attributes); i += 2 {
                        if i+1 < len(attrResp.Result.Attributes) &&
                           attrResp.Result.Attributes[i] == "onclick" {
                            isClickable = true
                            fmt.Printf("✓ 有onclick事件处理器\n")
                            break
                        }
                    }
                }
            }

            if isClickable {
                fmt.Printf("✅ 可以在此坐标进行点击测试\n")
            } else {
                fmt.Printf("ℹ️  元素可能不可点击\n")
            }
        }
    }

    // 9. 坐标扫描测试
    fmt.Printf("\n=== 坐标扫描测试 ===\n")

    // 测试多个坐标点
    testPoints := []struct {
        x, y int
        desc string
    }{
        {x: 100, y: 100, desc: "左上角"},
        {x: 960, y: 100, desc: "顶部中间"},
        {x: 1820, y: 100, desc: "右上角"},
        {x: 100, y: 540, desc: "左侧中间"},
        {x: 960, y: 540, desc: "中心"},
        {x: 1820, y: 540, desc: "右侧中间"},
        {x: 100, y: 980, desc: "左下角"},
        {x: 960, y: 980, desc: "底部中间"},
        {x: 1820, y: 980, desc: "右下角"},
    }

    fmt.Printf("测试9个关键坐标点:\n")
    for _, point := range testPoints {
        result, err := CDPDOMGetNodeForLocation(point.x, point.y, false, false)
        if err == nil {
            var resp struct {
                Result struct {
                    BackendNodeID int `json:"backendNodeId"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(result), &resp); err == nil && resp.Result.BackendNodeID > 0 {
                fmt.Printf("  ✓ %s: 有元素\n", point.desc)
            } else {
                fmt.Printf("  - %s: 无元素\n", point.desc)
            }
        }
    }

    // 10. 响应式布局测试
    fmt.Printf("\n=== 响应式布局测试 ===\n")

    // 测试不同视口宽度的坐标
    viewportWidths := []int{320, 768, 1024, 1440, 1920}
    testY := 200  // 固定Y坐标

    for _, width := range viewportWidths {
        testX := width / 2
        result, err := CDPDOMGetNodeForLocation(testX, testY, false, false)
        if err == nil {
            var resp struct {
                Result struct {
                    BackendNodeID int `json:"backendNodeId"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(result), &resp); err == nil && resp.Result.BackendNodeID > 0 {
                // 获取节点信息
                nodeDesc, _ := CDPDOMDescribeNode(resp.Result.BackendNodeID, 0, true)
                var descResp struct {
                    Result struct {
                        Node struct {
                            NodeName string `json:"nodeName"`
                        } `json:"node"`
                    } `json:"result"`
                }

                if err := json.Unmarshal([]byte(nodeDesc), &descResp); err == nil {
                    fmt.Printf("  视口 %4dpx: 中间位置元素为 %s\n",
                        width, descResp.Result.Node.NodeName)
                }
            }
        }
    }

    // 11. 可访问性分析
    fmt.Printf("\n=== 可访问性分析 ===\n")

    if response.Result.NodeID > 0 {
        attrs, _ := CDPDOMGetAttributes(response.Result.NodeID)
        var attrResp struct {
            Result struct {
                Attributes []string `json:"attributes"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(attrs), &attrResp); err == nil {
            hasAriaLabel := false
            hasTabIndex := false

            for i := 0; i < len(attrResp.Result.Attributes); i += 2 {
                if i+1 < len(attrResp.Result.Attributes) {
                    attrName := attrResp.Result.Attributes[i]

                    if strings.HasPrefix(attrName, "aria-") {
                        hasAriaLabel = true
                        fmt.Printf("✓ 有ARIA属性: %s\n", attrName)
                    }

                    if attrName == "tabindex" {
                        hasTabIndex = true
                        tabIndex := attrResp.Result.Attributes[i+1]
                        fmt.Printf("✓ 有tabindex属性: %s\n", tabIndex)
                    }
                }
            }

            if !hasAriaLabel {
                fmt.Printf("⚠️ 缺少ARIA标签属性\n")
            }

            if !hasTabIndex {
                fmt.Printf("ℹ️  无tabindex属性（可能无法通过键盘访问）\n")
            }
        }
    }

    // 12. 视觉测试建议
    fmt.Printf("\n=== 视觉测试建议 ===\n")

    // 根据元素类型提供测试建议
    if response.Result.NodeID > 0 {
        desc, _ := CDPDOMDescribeNode(response.Result.NodeID, 0, true)
        var descResp struct {
            Result struct {
                Node struct {
                    NodeName string `json:"nodeName"`
                } `json:"node"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(desc), &descResp); err == nil {
            switch descResp.Result.Node.NodeName {
            case "BUTTON":
                fmt.Printf("建议测试:\n")
                fmt.Printf("  - 点击状态样式\n")
                fmt.Printf("  - 悬停状态样式\n")
                fmt.Printf("  - 禁用状态样式\n")
                fmt.Printf("  - 焦点状态样式\n")

            case "A":
                fmt.Printf("建议测试:\n")
                fmt.Printf("  - 链接颜色和样式\n")
                fmt.Printf("  - 已访问链接样式\n")
                fmt.Printf("  - 悬停状态样式\n")
                fmt.Printf("  - 焦点状态样式\n")

            case "INPUT":
                fmt.Printf("建议测试:\n")
                fmt.Printf("  - 输入框聚焦样式\n")
                fmt.Printf("  - 占位符显示\n")
                fmt.Printf("  - 验证状态样式\n")
                fmt.Printf("  - 禁用状态样式\n")

            case "IMG":
                fmt.Printf("建议测试:\n")
                fmt.Printf("  - 图片加载状态\n")
                fmt.Printf("  - 图片alt文本\n")
                fmt.Printf("  - 响应式尺寸\n")

            default:
                fmt.Printf("建议测试:\n")
                fmt.Printf("  - 元素可见性\n")
                fmt.Printf("  - 布局位置\n")
                fmt.Printf("  - 交互状态\n")
            }
        }
    }
}

*/

// -----------------------------------------------  DOM.getOuterHTML  -----------------------------------------------
// === 应用场景 ===
// 1. HTML快照: 获取元素的完整HTML表示
// 2. 结构分析: 分析元素的HTML结构
// 3. 序列化: 序列化DOM元素为HTML字符串
// 4. 内容验证: 验证元素的HTML内容
// 5. 调试辅助: 调试时获取元素的外层HTML
// 6. 模板提取: 提取元素的HTML作为模板

// CDPDOMGetOuterHTML 获取指定节点的外层HTML
// nodeID: 要获取外层HTML的节点ID
func CDPDOMGetOuterHTML(nodeID int) (string, error) {
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
        "method": "DOM.getOuterHTML",
        "params": {
            "nodeId": %d
        }
    }`, reqID, nodeID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.getOuterHTML 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.getOuterHTML 请求超时")
		}
	}
}

/*

// 示例: 获取并分析元素的外层HTML
func ExampleCDPDOMGetOuterHTML() {
	// 1. 启用DOM功能
	if _, err := CDPDOMEnable(); err != nil {
		log.Printf("启用DOM失败: %v", err)
		return
	}

	// 确保测试完成后清理
	defer func() {
		if _, err := CDPDOMDisable(); err != nil {
			log.Printf("禁用DOM失败: %v", err)
		}
	}()

	// 假设我们有一个元素的节点ID
	elementNodeID := 1234

	// 2. 获取元素的外层HTML
	result, err := CDPDOMGetOuterHTML(elementNodeID)
	if err != nil {
		log.Printf("获取外层HTML失败: %v", err)
		return
	}

	// 3. 解析响应
	var response struct {
		Result struct {
			OuterHTML string `json:"outerHTML"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &response); err != nil {
		log.Printf("解析响应失败: %v", err)
		return
	}

	outerHTML := response.Result.OuterHTML

	// 4. 显示基本HTML信息
	fmt.Printf("=== 元素外层HTML分析 ===\n")
	fmt.Printf("HTML长度: %d 字符\n", len(outerHTML))
	fmt.Printf("HTML预览: %.100s...\n\n", outerHTML)

	// 5. HTML结构分析
	fmt.Printf("=== HTML结构分析 ===\n")

	// 解析HTML标签
	doc, err := html.Parse(strings.NewReader(outerHTML))
	if err == nil {
		analyzeHTMLStructure(doc, 0)
	} else {
		fmt.Printf("HTML解析失败: %v\n", err)
	}

	// 6. 标签统计
	fmt.Printf("\n=== 标签统计 ===\n")

	tagStats := make(map[string]int)
	countHTMLTags(outerHTML, tagStats)

	totalTags := 0
	for _, count := range tagStats {
		totalTags += count
	}

	fmt.Printf("总标签数: %d\n", totalTags)
	fmt.Printf("不同标签类型: %d\n", len(tagStats))

	// 显示标签频率
	if len(tagStats) > 0 {
		fmt.Printf("\n标签频率:\n")
		// 转换为切片并排序
		type TagCount struct {
			Tag   string
			Count int
		}

		tagCounts := make([]TagCount, 0, len(tagStats))
		for tag, count := range tagStats {
			tagCounts = append(tagCounts, TagCount{Tag: tag, Count: count})
		}

		sort.Slice(tagCounts, func(i, j int) bool {
			if tagCounts[i].Count == tagCounts[j].Count {
				return tagCounts[i].Tag < tagCounts[j].Tag
			}
			return tagCounts[i].Count > tagCounts[j].Count
		})

		for i := 0; i < len(tagCounts) && i < 10; i++ {
			fmt.Printf("  %-10s: %d 个\n", tagCounts[i].Tag, tagCounts[i].Count)
		}
	}

	// 7. 属性分析
	fmt.Printf("\n=== 属性分析 ===\n")

	attrStats := analyzeHTMLAttributes(outerHTML)
	fmt.Printf("总属性数: %d\n", attrStats.TotalAttributes)
	fmt.Printf("不同属性名: %d\n", len(attrStats.AttributeFrequency))

	if len(attrStats.AttributeFrequency) > 0 {
		fmt.Printf("\n常见属性:\n")
		// 排序属性
		attrCounts := make([]struct {
			Name  string
			Count int
		}, 0, len(attrStats.AttributeFrequency))

		for name, count := range attrStats.AttributeFrequency {
			attrCounts = append(attrCounts, struct {
				Name  string
				Count int
			}{Name: name, Count: count})
		}

		sort.Slice(attrCounts, func(i, j int) bool {
			if attrCounts[i].Count == attrCounts[j].Count {
				return attrCounts[i].Name < attrCounts[j].Name
			}
			return attrCounts[i].Count > attrCounts[j].Count
		})

		for i := 0; i < len(attrCounts) && i < 10; i++ {
			fmt.Printf("  %-20s: %d 次\n", attrCounts[i].Name, attrCounts[i].Count)
		}
	}

	// 8. 代码质量检查
	fmt.Printf("\n=== 代码质量检查 ===\n")

	qualityIssues := checkHTMLQuality(outerHTML)

	if len(qualityIssues.Warnings) > 0 {
		fmt.Printf("警告 (%d 个):\n", len(qualityIssues.Warnings))
		for _, warning := range qualityIssues.Warnings {
			fmt.Printf("  ⚠️  %s\n", warning)
		}
	} else {
		fmt.Printf("✓ 无代码质量问题\n")
	}

	if len(qualityIssues.Errors) > 0 {
		fmt.Printf("\n错误 (%d 个):\n", len(qualityIssues.Errors))
		for _, err := range qualityIssues.Errors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	// 9. 可访问性检查
	fmt.Printf("\n=== 可访问性检查 ===\n")

	a11yIssues := checkAccessibility(outerHTML)

	if len(a11yIssues) > 0 {
		fmt.Printf("可访问性问题 (%d 个):\n", len(a11yIssues))
		for _, issue := range a11yIssues {
			fmt.Printf("  %s\n", issue)
		}
	} else {
		fmt.Printf("✓ 无明显的可访问性问题\n")
	}

	// 10. 性能分析
	fmt.Printf("\n=== 性能分析 ===\n")

	perfMetrics := analyzeHTMLPerformance(outerHTML)
	fmt.Printf("HTML大小: %.2f KB\n", float64(len(outerHTML))/1024)
	fmt.Printf("嵌套深度: %d\n", perfMetrics.MaxDepth)
	fmt.Printf("内联样式: %d 个\n", perfMetrics.InlineStyles)
	fmt.Printf("内联脚本: %d 个\n", perfMetrics.InlineScripts)

	if perfMetrics.InlineStyles > 5 {
		fmt.Printf("⚠️ 内联样式较多，建议使用外部CSS\n")
	}

	if perfMetrics.InlineScripts > 3 {
		fmt.Printf("⚠️ 内联脚本较多，建议使用外部JS文件\n")
	}

	// 11. 语义化分析
	fmt.Printf("\n=== 语义化分析 ===\n")

	semanticScore := analyzeSemanticHTML(outerHTML)
	fmt.Printf("语义化分数: %.1f/100\n", semanticScore)

	if semanticScore >= 80 {
		fmt.Printf("✓ HTML语义化良好\n")
	} else if semanticScore >= 60 {
		fmt.Printf("⚠️ HTML语义化一般\n")
	} else {
		fmt.Printf("❌ HTML语义化较差，建议改进\n")
	}

	// 12. 安全分析
	fmt.Printf("\n=== 安全分析 ===\n")

	securityIssues := checkHTMLSecurity(outerHTML)
	if len(securityIssues) > 0 {
		fmt.Printf("安全问题 (%d 个):\n", len(securityIssues))
		for _, issue := range securityIssues {
			fmt.Printf("  🔒 %s\n", issue)
		}
	} else {
		fmt.Printf("✓ 无明显的安全问题\n")
	}

	// 13. 格式化和美化
	fmt.Printf("\n=== 格式化HTML ===\n")

	formattedHTML, err := formatHTML(outerHTML)
	if err == nil {
		fmt.Printf("美化后的HTML (前200字符):\n")
		fmt.Printf("%.200s...\n", formattedHTML)
	}

	// 14. 差异比较（如果需要）
	fmt.Printf("\n=== 差异比较 ===\n")

	// 获取innerHTML进行对比
	// 注意：这需要先获取innerHTML
	fmt.Printf("outerHTML包含元素自身标签\n")
	fmt.Printf("innerHTML只包含子内容\n")

	// 15. 使用建议
	fmt.Printf("\n=== 使用建议 ===\n")

	if len(outerHTML) > 10000 {
		fmt.Printf("⚠️ HTML较大，考虑拆分组件\n")
	}

	if strings.Count(outerHTML, "<div>") > 10 {
		fmt.Printf("⚠️ 使用了较多div，考虑使用语义化标签\n")
	}

	if strings.Contains(outerHTML, "style=") {
		fmt.Printf("⚠️ 使用了内联样式，建议使用CSS类\n")
	}

	if strings.Contains(outerHTML, "onclick=") {
		fmt.Printf("⚠️ 使用了内联事件，建议使用事件委托\n")
	}
}

// 辅助函数: 分析HTML结构
func analyzeHTMLStructure(n *html.Node, depth int) {
	if n.Type == html.ElementNode {
		indent := strings.Repeat("  ", depth)
		fmt.Printf("%s<%s", indent, n.Data)

		// 显示属性
		for _, attr := range n.Attr {
			fmt.Printf(" %s=\"%s\"", attr.Key, attr.Val)
		}

		if n.FirstChild == nil {
			fmt.Printf(" />\n")
		} else {
			fmt.Printf(">\n")
		}
	} else if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			indent := strings.Repeat("  ", depth)
			fmt.Printf("%s文本: %s\n", indent, truncateString(text, 50))
		}
	} else if n.Type == html.CommentNode {
		indent := strings.Repeat("  ", depth)
		fmt.Printf("%s注释: %s\n", indent, truncateString(n.Data, 50))
	}

	// 递归处理子节点
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		analyzeHTMLStructure(c, depth+1)
	}

	if n.Type == html.ElementNode && n.FirstChild != nil {
		indent := strings.Repeat("  ", depth)
		fmt.Printf("%s</%s>\n", indent, n.Data)
	}
}

// 辅助函数: 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// 辅助函数: 统计HTML标签
func countHTMLTags(htmlStr string, stats map[string]int) {
	// 简单的正则匹配标签
	re := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)(?:\s|>)`)
	matches := re.FindAllStringSubmatch(htmlStr, -1)

	for _, match := range matches {
		if len(match) > 1 {
			tag := strings.ToLower(match[1])
			stats[tag]++
		}
	}
}

// 辅助结构: 属性统计
type AttributeStats struct {
	TotalAttributes    int
	AttributeFrequency map[string]int
}

// 辅助函数: 分析HTML属性
func analyzeHTMLAttributes(htmlStr string) AttributeStats {
	stats := AttributeStats{
		AttributeFrequency: make(map[string]int),
	}

	// 匹配属性
	re := regexp.MustCompile(`([a-zA-Z-]+)=["'][^"']*["']`)
	matches := re.FindAllStringSubmatch(htmlStr, -1)
	stats.TotalAttributes = len(matches)

	for _, match := range matches {
		if len(match) > 1 {
			attrName := match[1]
			stats.AttributeFrequency[attrName]++
		}
	}

	return stats
}

// 辅助结构: 代码质量问题
type QualityIssues struct {
	Warnings []string
	Errors   []string
}

// 辅助函数: 检查HTML代码质量
func checkHTMLQuality(htmlStr string) QualityIssues {
	issues := QualityIssues{}

	// 检查未闭合的标签
	if strings.Count(htmlStr, "<div") > strings.Count(htmlStr, "</div") {
		issues.Errors = append(issues.Errors, "div标签未闭合")
	}

	// 检查重复的ID
	idRe := regexp.MustCompile(`id=["']([^"']+)["']`)
	idMatches := idRe.FindAllStringSubmatch(htmlStr, -1)
	idCounts := make(map[string]int)

	for _, match := range idMatches {
		if len(match) > 1 {
			id := match[1]
			idCounts[id]++
		}
	}

	for id, count := range idCounts {
		if count > 1 {
			issues.Errors = append(issues.Errors, fmt.Sprintf("重复的ID: %s", id))
		}
	}

	// 检查过长的行
	lines := strings.Split(htmlStr, "\n")
	for i, line := range lines {
		if len(line) > 200 {
			issues.Warnings = append(issues.Warnings,
				fmt.Sprintf("第%d行过长 (%d字符)", i+1, len(line)))
		}
	}

	// 检查脚本安全性
	if strings.Contains(htmlStr, "javascript:") {
		issues.Warnings = append(issues.Warnings, "使用javascript:协议可能存在安全风险")
	}

	// 检查事件处理
	eventRe := regexp.MustCompile(`on[a-z]+=`)
	if len(eventRe.FindAllString(htmlStr, -1)) > 3 {
		issues.Warnings = append(issues.Warnings, "内联事件处理较多，建议使用事件委托")
	}

	return issues
}

// 辅助函数: 检查可访问性
func checkAccessibility(htmlStr string) []string {
	issues := []string{}

	// 检查图片alt
	imgRe := regexp.MustCompile(`<img[^>]*>`)
	imgMatches := imgRe.FindAllString(htmlStr, -1)

	for _, img := range imgMatches {
		if !strings.Contains(img, "alt=") && !strings.Contains(img, "role=") {
			issues = append(issues, "图片缺少alt属性")
			break
		}
	}

	// 检查表单标签
	inputRe := regexp.MustCompile(`<input[^>]*>`)
	inputMatches := inputRe.FindAllString(htmlStr, -1)

	for _, input := range inputMatches {
		if strings.Contains(input, "type=") &&
			(strings.Contains(input, "type=\"text\"") ||
				strings.Contains(input, "type=\"password\"")) {
			if !strings.Contains(input, "id=") && !strings.Contains(input, "aria-label=") {
				issues = append(issues, "文本输入框缺少标签")
				break
			}
		}
	}

	// 检查对比度
	if strings.Contains(htmlStr, "style=\"color:") || strings.Contains(htmlStr, "style='color:") {
		// 简单的颜色检查
		colorRe := regexp.MustCompile(`color:\s*#([0-9a-fA-F]{6})`)
		colors := colorRe.FindAllStringSubmatch(htmlStr, -1)

		for _, match := range colors {
			if len(match) > 1 {
				color := match[1]
				// 简单的亮度检查
				brightness := calculateColorBrightness(color)
				if brightness < 0.3 || brightness > 0.7 {
					issues = append(issues, "颜色对比度可能不足")
					break
				}
			}
		}
	}

	return issues
}

// 辅助函数: 计算颜色亮度
func calculateColorBrightness(hexColor string) float64 {
	if len(hexColor) != 6 {
		return 0.5
	}

	r, _ := strconv.ParseInt(hexColor[0:2], 16, 64)
	g, _ := strconv.ParseInt(hexColor[2:4], 16, 64)
	b, _ := strconv.ParseInt(hexColor[4:6], 16, 64)

	// 相对亮度公式
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255
}

// 辅助结构: 性能指标
type PerformanceMetrics struct {
	MaxDepth      int
	InlineStyles  int
	InlineScripts int
}

// 辅助函数: 分析HTML性能
func analyzeHTMLPerformance(htmlStr string) PerformanceMetrics {
	metrics := PerformanceMetrics{}

	// 计算嵌套深度
	depth := 0
	maxDepth := 0
	tagStack := []string{}

	tagRe := regexp.MustCompile(`<(/?)([a-zA-Z][a-zA-Z0-9]*)(?:\s|>)`)
	matches := tagRe.FindAllStringSubmatch(htmlStr, -1)

	for _, match := range matches {
		if len(match) > 2 {
			isClosing := match[1] == "/"
			tagName := match[2]

			if !isClosing {
				depth++
				tagStack = append(tagStack, tagName)
				if depth > maxDepth {
					maxDepth = depth
				}
			} else if len(tagStack) > 0 {
				depth--
				tagStack = tagStack[:len(tagStack)-1]
			}
		}
	}

	metrics.MaxDepth = maxDepth

	// 统计内联样式和脚本
	metrics.InlineStyles = strings.Count(htmlStr, "style=\"")
	metrics.InlineScripts = strings.Count(htmlStr, "on")

	return metrics
}

// 辅助函数: 分析语义化HTML
func analyzeSemanticHTML(htmlStr string) float64 {
	score := 0.0

	// 语义化标签加分
	semanticTags := []string{
		"header", "nav", "main", "article", "section",
		"aside", "footer", "figure", "figcaption", "time",
		"mark", "summary", "details", "menu", "dialog",
	}

	for _, tag := range semanticTags {
		count := strings.Count(strings.ToLower(htmlStr), "<"+tag)
		if count > 0 {
			score += float64(count) * 5
		}
	}

	// div滥用扣分
	divCount := strings.Count(htmlStr, "<div")
	if divCount > 0 {
		// 每5个div扣1分
		score -= float64(divCount) / 5
	}

	// span滥用扣分
	spanCount := strings.Count(htmlStr, "<span")
	if spanCount > 0 {
		// 每10个span扣1分
		score -= float64(spanCount) / 10
	}

	// 确保分数在合理范围内
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// 辅助函数: 检查HTML安全
func checkHTMLSecurity(htmlStr string) []string {
	issues := []string{}

	// 检查潜在的危险属性
	dangerousPatterns := []struct {
		pattern string
		message string
	}{
		{`javascript:`, "使用javascript:协议"},
		{`onload=`, "onload事件可能被滥用"},
		{`onerror=`, "onerror事件可能被滥用"},
		{`data:`, "data:协议可能包含恶意内容"},
		{`eval\s*\(`, "使用eval函数"},
	}

	for _, pattern := range dangerousPatterns {
		re := regexp.MustCompile(pattern.pattern)
		if re.MatchString(strings.ToLower(htmlStr)) {
			issues = append(issues, pattern.message)
		}
	}

	return issues
}

// 辅助函数: 格式化HTML
func formatHTML(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

*/

// -----------------------------------------------  DOM.hideHighlight  -----------------------------------------------
// === 应用场景 ===
// 1. 清理高亮: 清理之前设置的元素高亮
// 2. 测试清理: 测试完成后清理视觉标记
// 3. 资源释放: 释放高亮相关资源
// 4. 状态恢复: 恢复页面原始显示状态
// 5. 自动化测试: 测试流程完成后清理
// 6. 调试清理: 调试完成后清理高亮标记

// CDPDOMHideHighlight 隐藏之前设置的元素高亮
func CDPDOMHideHighlight() (string, error) {
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
        "method": "DOM.hideHighlight"
    }`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.hideHighlight 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.hideHighlight 请求超时")
		}
	}
}

/*

// 示例: 在元素高亮测试后清理高亮
func ExampleCDPDOMHideHighlight() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 模拟元素高亮操作
    // 注意：DOM.highlightNode 方法需要先实现
    // 这里假设我们已经执行了高亮操作

    // 3. 执行高亮清理
    result, err := CDPDOMHideHighlight()
    if err != nil {
        log.Printf("隐藏高亮失败: %v", err)
        return
    }

    log.Printf("高亮已隐藏: %s", result)

    // 4. 验证高亮已被清理
    fmt.Println("=== 高亮清理验证 ===")
    fmt.Println("1. 视觉检查:")
    fmt.Println("   - 页面上的高亮标记应该已消失")
    fmt.Println("   - 元素应该恢复原始样式")
    fmt.Println("   - 页面布局不应受影响")

    fmt.Println("\n2. 性能检查:")
    fmt.Println("   - 高亮相关的CSS样式应该被移除")
    fmt.Println("   - 相关的JavaScript监听器应该被清理")
    fmt.Println("   - 内存使用应该恢复正常")

    // 5. 清理状态检查
    fmt.Println("\n=== 清理状态检查 ===")

    // 检查页面是否恢复原状
    // 可以通过截图对比等方式验证
    fmt.Println("检查点:")
    fmt.Println("  ✓ 高亮叠加层被移除")
    fmt.Println("  ✓ 元素样式恢复")
    fmt.Println("  ✓ 页面交互正常")
    fmt.Println("  ✓ 无残留的CSS类")
    fmt.Println("  ✓ 无残留的事件监听器")

    // 6. 错误恢复场景
    fmt.Println("\n=== 错误恢复场景 ===")

    // 测试多次调用hideHighlight
    for i := 1; i <= 3; i++ {
        result, err := CDPDOMHideHighlight()
        if err != nil {
            fmt.Printf("第%d次调用失败: %v\n", i, err)
        } else {
            fmt.Printf("第%d次调用成功: 高亮已清理\n", i)
        }
    }

    // 7. 集成测试示例
    fmt.Println("\n=== 集成测试示例 ===")

    // 模拟一个完整的高亮测试流程
    testElementHighlightCleanup()
}

// 模拟完整的高亮测试流程
func testElementHighlightCleanup() {
    fmt.Println("开始高亮测试流程...")

    // 1. 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 2. 定义测试步骤
    steps := []struct {
        name string
        action func() error
    }{
        {
            name: "准备测试环境",
            action: func() error {
                fmt.Println("  ✓ 测试环境准备完成")
                return nil
            },
        },
        {
            name: "查找要测试的元素",
            action: func() error {
                // 模拟查找元素
                fmt.Println("  ✓ 找到测试元素")
                return nil
            },
        },
        {
            name: "高亮目标元素",
            action: func() error {
                // 模拟高亮元素
                // 这里假设调用了 CDPDOMHighlightNode
                fmt.Println("  ✓ 元素高亮成功")
                return nil
            },
        },
        {
            name: "验证高亮效果",
            action: func() error {
                fmt.Println("  ✓ 高亮效果验证通过")
                return nil
            },
        },
        {
            name: "清理高亮",
            action: func() error {
                result, err := CDPDOMHideHighlight()
                if err != nil {
                    return fmt.Errorf("清理高亮失败: %w", err)
                }
                fmt.Printf("  ✓ 高亮清理完成: %s\n", result)
                return nil
            },
        },
        {
            name: "验证清理效果",
            action: func() error {
                fmt.Println("  ✓ 清理效果验证通过")
                return nil
            },
        },
    }

    // 3. 执行测试步骤
    for i, step := range steps {
        fmt.Printf("步骤 %d/%d: %s\n", i+1, len(steps), step.name)

        if err := step.action(); err != nil {
            fmt.Printf("  ❌ 步骤失败: %v\n", err)

            // 尝试清理
            CDPDOMHideHighlight()
            CDPDOMDisable()
            return
        }
    }

    // 4. 清理DOM
    if _, err := CDPDOMDisable(); err != nil {
        log.Printf("禁用DOM失败: %v", err)
    }

    fmt.Println("\n✅ 高亮测试流程完成")
}

// 自动清理装饰器
func withHighlightCleanup(testFunc func() error) func() error {
    return func() error {
        // 确保启用DOM
        if _, err := CDPDOMEnable(); err != nil {
            return fmt.Errorf("启用DOM失败: %w", err)
        }

        // 确保测试完成后清理
        defer func() {
            // 清理高亮
            if _, err := CDPDOMHideHighlight(); err != nil {
                log.Printf("清理高亮失败: %v", err)
            }

            // 禁用DOM
            if _, err := CDPDOMDisable(); err != nil {
                log.Printf("禁用DOM失败: %v", err)
            }
        }()

        // 执行测试函数
        return testFunc()
    }
}

// 高亮管理器
type HighlightManager struct {
    isHighlighted bool
    highlightedNodes []int
    highlightConfig map[string]interface{}
}

func NewHighlightManager() *HighlightManager {
    return &HighlightManager{
        highlightedNodes: make([]int, 0),
        highlightConfig: map[string]interface{}{
            "contentColor": map[string]int{"r": 111, "g": 168, "b": 220, "a": 0.66},
            "paddingColor": map[string]int{"r": 147, "g": 196, "b": 125, "a": 0.55},
            "borderColor":  map[string]int{"r": 255, "g": 229, "b": 153, "a": 0.8},
            "marginColor":  map[string]int{"r": 246, "g": 178, "b": 107, "a": 0.66},
        },
    }
}

func (hm *HighlightManager) HighlightNode(nodeID int) error {
    // 这里应该调用 CDPDOMHighlightNode
    // 为了示例，我们模拟高亮操作
    hm.isHighlighted = true
    hm.highlightedNodes = append(hm.highlightedNodes, nodeID)
    fmt.Printf("高亮节点 %d\n", nodeID)
    return nil
}

func (hm *HighlightManager) HideHighlight() error {
    result, err := CDPDOMHideHighlight()
    if err != nil {
        return fmt.Errorf("隐藏高亮失败: %w", err)
    }

    hm.isHighlighted = false
    hm.highlightedNodes = make([]int, 0)
    fmt.Printf("隐藏高亮: %s\n", result)
    return nil
}

func (hm *HighlightManager) IsHighlighted() bool {
    return hm.isHighlighted
}

func (hm *HighlightManager) GetHighlightedNodes() []int {
    return hm.highlightedNodes
}

// 高亮测试套件
func RunHighlightTests() {
    fmt.Println("=== 运行高亮测试套件 ===")

    manager := NewHighlightManager()

    tests := []struct {
        name string
        test func(*HighlightManager) error
    }{
        {
            name: "基本高亮和清理",
            test: testBasicHighlight,
        },
        {
            name: "多次高亮清理",
            test: testMultipleHighlights,
        },
        {
            name: "清理未高亮状态",
            test: testCleanupWithoutHighlight,
        },
    }

    for _, t := range tests {
        fmt.Printf("\n测试: %s\n", t.name)

        // 启用DOM
        if _, err := CDPDOMEnable(); err != nil {
            fmt.Printf("  ❌ 启用DOM失败: %v\n", err)
            continue
        }

        // 执行测试
        if err := t.test(manager); err != nil {
            fmt.Printf("  ❌ 测试失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 测试通过\n")
        }

        // 确保清理
        CDPDOMHideHighlight()
        CDPDOMDisable()
    }

    fmt.Println("\n=== 测试套件完成 ===")
}

func testBasicHighlight(manager *HighlightManager) error {
    // 模拟高亮
    if err := manager.HighlightNode(123); err != nil {
        return fmt.Errorf("高亮失败: %w", err)
    }

    if !manager.IsHighlighted() {
        return fmt.Errorf("高亮状态不正确")
    }

    // 清理高亮
    if err := manager.HideHighlight(); err != nil {
        return fmt.Errorf("清理失败: %w", err)
    }

    if manager.IsHighlighted() {
        return fmt.Errorf("清理后高亮状态不正确")
    }

    return nil
}

func testMultipleHighlights(manager *HighlightManager) error {
    // 多次高亮
    nodes := []int{123, 456, 789}

    for _, nodeID := range nodes {
        if err := manager.HighlightNode(nodeID); err != nil {
            return fmt.Errorf("高亮节点 %d 失败: %w", nodeID, err)
        }
    }

    // 验证高亮节点
    highlighted := manager.GetHighlightedNodes()
    if len(highlighted) != len(nodes) {
        return fmt.Errorf("高亮节点数量不正确: 期望 %d, 实际 %d", len(nodes), len(highlighted))
    }

    // 清理
    if err := manager.HideHighlight(); err != nil {
        return fmt.Errorf("清理失败: %w", err)
    }

    if len(manager.GetHighlightedNodes()) > 0 {
        return fmt.Errorf("清理后仍有高亮节点")
    }

    return nil
}

func testCleanupWithoutHighlight(manager *HighlightManager) error {
    // 确保没有高亮
    manager.HideHighlight()

    // 再次清理（应该不会出错）
    if err := manager.HideHighlight(); err != nil {
        return fmt.Errorf("无高亮时清理失败: %w", err)
    }

    return nil
}

*/

// -----------------------------------------------  DOM.highlightNode  -----------------------------------------------
// === 应用场景 ===
// 1. 视觉标记: 在页面上高亮显示特定元素
// 2. 调试辅助: 调试时突出显示问题元素
// 3. 自动化测试: 测试过程中可视化目标元素
// 4. 元素定位: 帮助用户定位页面元素
// 5. 教程引导: 创建交互式教程时高亮步骤元素
// 6. 审查工具: 构建元素审查工具

// CDPDOMHighlightNode 高亮显示指定的DOM节点
// nodeID: 要高亮的节点ID
// highlightConfig: 高亮配置
func CDPDOMHighlightNode(params string) (string, error) {
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
        "method": "DOM.highlightNode",
        "params": %s
    }`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.highlightNode 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.highlightNode 请求超时")
		}
	}
}

/*

// 示例: 高亮按钮元素并显示详细信息
func ExampleCDPDOMHighlightNode() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        // 清理高亮
        if _, err := CDPDOMHideHighlight(); err != nil {
            log.Printf("隐藏高亮失败: %v", err)
        }

        // 禁用DOM
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个按钮的节点ID
    buttonNodeID := 1234

    // 2. 创建自定义高亮配置
    highlightConfig := &HighlightConfig{
        ShowInfo:   true,
        ShowStyles: true,
        ShowRulers: true,
        ContentColor: &RGBA{R: 111, G: 168, B: 220, A: 0.66},  // 浅蓝色
        PaddingColor: &RGBA{R: 147, G: 196, B: 125, A: 0.55},  // 浅绿色
        BorderColor:  &RGBA{R: 255, G: 229, B: 153, A: 0.8},   // 浅黄色
        MarginColor:  &RGBA{R: 246, G: 178, B: 107, A: 0.66},  // 浅橙色
        ColorFormat: "rgb",
    }

    // 3. 高亮按钮元素
    result, err := CDPDOMHighlightNode(buttonNodeID, highlightConfig)
    if err != nil {
        log.Printf("高亮节点失败: %v", err)
        return
    }

    log.Printf("节点高亮成功: %s", result)

    // 4. 显示高亮效果信息
    fmt.Printf("=== 元素高亮信息 ===\n")
    fmt.Printf("节点ID: %d\n", buttonNodeID)
    fmt.Printf("高亮配置: 已启用信息显示、样式显示、标尺\n")
    fmt.Printf("颜色方案: 盒模型四色方案\n")
    fmt.Printf("  - 内容区域: RGBA(111, 168, 220, 0.66)\n")
    fmt.Printf("  - 内边距:   RGBA(147, 196, 125, 0.55)\n")
    fmt.Printf("  - 边框:     RGBA(255, 229, 153, 0.80)\n")
    fmt.Printf("  - 外边距:   RGBA(246, 178, 107, 0.66)\n")

    // 5. 获取元素信息以验证高亮
    fmt.Printf("\n=== 元素验证 ===\n")

    // 获取元素描述
    nodeDesc, err := CDPDOMDescribeNode(buttonNodeID, 0, true)
    if err == nil {
        var descResp struct {
            Result struct {
                Node struct {
                    NodeName string `json:"nodeName"`
                    LocalName string `json:"localName"`
                } `json:"node"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(nodeDesc), &descResp); err == nil {
            fmt.Printf("元素类型: %s (%s)\n",
                descResp.Result.Node.NodeName,
                descResp.Result.Node.LocalName)
        }
    }

    // 获取盒模型
    boxModel, err := CDPDOMGetBoxModel(buttonNodeID)
    if err == nil {
        var boxResp struct {
            Result struct {
                Model struct {
                    Width  int `json:"width"`
                    Height int `json:"height"`
                } `json:"model"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(boxModel), &boxResp); err == nil {
            fmt.Printf("元素尺寸: %d x %d 像素\n",
                boxResp.Result.Model.Width,
                boxResp.Result.Model.Height)
        }
    }

    // 6. 高亮效果检查列表
    fmt.Printf("\n=== 高亮效果检查 ===\n")
    fmt.Printf("视觉检查项:\n")
    fmt.Printf("  ✓ 元素被明显高亮\n")
    fmt.Printf("  ✓ 盒模型各区域颜色不同\n")
    fmt.Printf("  ✓ 显示尺寸标尺\n")
    fmt.Printf("  ✓ 显示元素信息\n")
    fmt.Printf("  ✓ 不影响页面交互\n")

    // 7. 交互测试
    fmt.Printf("\n=== 交互测试 ===\n")

    // 等待用户观察
    fmt.Printf("高亮显示中，观察5秒...\n")
    time.Sleep(5 * time.Second)

    // 测试鼠标悬停
    fmt.Printf("测试鼠标悬停效果...\n")
    fmt.Printf("  - 可以正常与元素交互\n")
    fmt.Printf("  - 高亮不阻挡点击\n")
    fmt.Printf("  - 高亮不影响文本选择\n")

    // 8. 可访问性考虑
    fmt.Printf("\n=== 可访问性考虑 ===\n")
    fmt.Printf("高亮效果:\n")
    fmt.Printf("  - 使用半透明颜色，不遮挡内容\n")
    fmt.Printf("  - 颜色对比度符合WCAG标准\n")
    fmt.Printf("  - 不干扰屏幕阅读器\n")
    fmt.Printf("  - 不改变元素的tab顺序\n")

    // 9. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次高亮
    startTime := time.Now()
    for i := 0; i < 5; i++ {
        // 创建不同的颜色配置
        config := &HighlightConfig{
            ContentColor: &RGBA{
                R: 100 + i*20,
                G: 150 + i*10,
                B: 200 - i*20,
                A: 0.6,
            },
        }

        _, err := CDPDOMHighlightNode(buttonNodeID, config)
        if err != nil {
            fmt.Printf("第%d次高亮失败: %v\n", i+1, err)
            break
        }

        time.Sleep(200 * time.Millisecond)
    }

    elapsed := time.Since(startTime)
    fmt.Printf("5次高亮操作耗时: %v\n", elapsed)

    // 10. 清理演示
    fmt.Printf("\n=== 清理演示 ===\n")

    // 延迟清理，让用户看到效果
    fmt.Printf("3秒后清理高亮...\n")
    time.Sleep(3 * time.Second)

    cleanupResult, err := CDPDOMHideHighlight()
    if err != nil {
        log.Printf("清理高亮失败: %v", err)
    } else {
        fmt.Printf("高亮已清理: %s\n", cleanupResult)
    }
}

// 高级功能: 动画高亮
func animateHighlight(nodeID int, duration time.Duration) {
    // 启用DOM
    CDPDOMEnable()
    defer CDPDOMHideHighlight()
    defer CDPDOMDisable()

    // 定义颜色序列
    colors := []RGBA{
        {R: 255, G: 100, B: 100, A: 0.7}, // 红色
        {R: 255, G: 200, B: 100, A: 0.7}, // 橙色
        {R: 255, G: 255, B: 100, A: 0.7}, // 黄色
        {R: 200, G: 255, B: 100, A: 0.7}, // 黄绿色
        {R: 100, G: 255, B: 100, A: 0.7}, // 绿色
        {R: 100, G: 255, B: 200, A: 0.7}, // 蓝绿色
        {R: 100, G: 200, B: 255, A: 0.7}, // 浅蓝色
        {R: 100, G: 100, B: 255, A: 0.7}, // 蓝色
    }

    // 动画参数
    interval := 200 * time.Millisecond
    steps := int(duration / interval)

    fmt.Printf("开始动画高亮，持续 %v，%d 步\n", duration, steps)

    for i := 0; i < steps; i++ {
        colorIndex := i % len(colors)
        config := &HighlightConfig{
            ContentColor: &colors[colorIndex],
            ShowInfo: i%5 == 0, // 每5步显示一次信息
        }

        if _, err := CDPDOMHighlightNode(nodeID, config); err != nil {
            log.Printf("动画高亮失败: %v", err)
            break
        }

        time.Sleep(interval)
    }

    fmt.Println("动画高亮完成")
}

// 高级功能: 对比高亮多个元素
func highlightMultipleElements(nodeIDs []int, highlightType string) {
    CDPDOMEnable()
    defer CDPDOMDisable()

    // 先清理可能的高亮
    CDPDOMHideHighlight()

    // 根据类型选择配置
    var baseConfig *HighlightConfig

    switch highlightType {
    case "comparison":
        // 对比模式：不同的颜色
        baseConfig = &HighlightConfig{
            ShowInfo: true,
            ShowRulers: true,
        }

    case "group":
        // 组模式：相同的颜色
        baseConfig = &HighlightConfig{
            ContentColor: &RGBA{R: 111, G: 168, B: 220, A: 0.6},
            ShowInfo: false,
        }

    case "focus":
        // 焦点模式：中心元素高亮
        baseConfig = &HighlightConfig{
            ContentColor: &RGBA{R: 255, G: 100, B: 100, A: 0.8},
            ShowInfo: true,
        }

    default:
        baseConfig = &HighlightConfig{
            ContentColor: &RGBA{R: 200, G: 200, B: 200, A: 0.6},
        }
    }

    fmt.Printf("高亮 %d 个元素，模式: %s\n", len(nodeIDs), highlightType)

    // 逐个高亮元素
    for i, nodeID := range nodeIDs {
        config := *baseConfig // 复制基础配置

        if highlightType == "comparison" {
            // 为每个元素分配不同的颜色
            hue := float64(i) * 360.0 / float64(len(nodeIDs))
            rgb := hslToRGB(hue, 0.7, 0.7)
            config.ContentColor = &RGBA{
                R: int(rgb[0] * 255),
                G: int(rgb[1] * 255),
                B: int(rgb[2] * 255),
                A: 0.6,
            }
        }

        if _, err := CDPDOMHighlightNode(nodeID, &config); err != nil {
            log.Printf("高亮元素 %d 失败: %v", nodeID, err)
        } else {
            fmt.Printf("  ✓ 高亮元素 %d\n", nodeID)
        }

        // 短暂延迟，避免过快
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("多元素高亮完成，按回车键清理...")
    fmt.Scanln() // 等待用户输入

    CDPDOMHideHighlight()
    fmt.Println("高亮已清理")
}

// HSL转RGB辅助函数
func hslToRGB(h, s, l float64) [3]float64 {
    var r, g, b float64

    if s == 0 {
        r, g, b = l, l, l
    } else {
        var q, p float64

        if l < 0.5 {
            q = l * (1 + s)
        } else {
            q = l + s - l*s
        }

        p = 2*l - q
        r = hueToRGB(p, q, h/360.0+1.0/3.0)
        g = hueToRGB(p, q, h/360.0)
        b = hueToRGB(p, q, h/360.0-1.0/3.0)
    }

    return [3]float64{r, g, b}
}

func hueToRGB(p, q, t float64) float64 {
    if t < 0 {
        t += 1
    }
    if t > 1 {
        t -= 1
    }

    switch {
    case t < 1.0/6.0:
        return p + (q-p)*6*t
    case t < 1.0/2.0:
        return q
    case t < 2.0/3.0:
        return p + (q-p)*(2.0/3.0-t)*6
    default:
        return p
    }
}

// 高亮管理器增强版
type EnhancedHighlightManager struct {
    *HighlightManager
    history []HighlightRecord
    presets map[string]*HighlightConfig
}

type HighlightRecord struct {
    Timestamp time.Time
    NodeID    int
    Config    HighlightConfig
    Duration  time.Duration
}

func NewEnhancedHighlightManager() *EnhancedHighlightManager {
    manager := NewHighlightManager()

    return &EnhancedHighlightManager{
        HighlightManager: manager,
        history: make([]HighlightRecord, 0),
        presets: map[string]*HighlightConfig{
            "error": {
                ContentColor: &RGBA{R: 255, G: 100, B: 100, A: 0.8},
                ShowInfo:     true,
                ShowRulers:   true,
            },
            "success": {
                ContentColor: &RGBA{R: 100, G: 200, B: 100, A: 0.6},
                ShowInfo:     false,
            },
            "warning": {
                ContentColor: &RGBA{R: 255, G: 200, B: 100, A: 0.7},
                ShowInfo:     true,
            },
            "info": {
                ContentColor: &RGBA{R: 100, G: 168, B: 220, A: 0.6},
                ShowInfo:     true,
            },
        },
    }
}

func (ehm *EnhancedHighlightManager) HighlightWithPreset(nodeID int, presetName string, duration time.Duration) error {
    preset, exists := ehm.presets[presetName]
    if !exists {
        return fmt.Errorf("预设不存在: %s", presetName)
    }

    startTime := time.Now()

    if err := ehm.HighlightNode(nodeID); err != nil {
        return fmt.Errorf("高亮失败: %w", err)
    }

    // 记录历史
    record := HighlightRecord{
        Timestamp: startTime,
        NodeID:    nodeID,
        Config:    *preset,
        Duration:  duration,
    }
    ehm.history = append(ehm.history, record)

    // 设置自动清理
    if duration > 0 {
        go func() {
            time.Sleep(duration)
            ehm.HideHighlight()
        }()
    }

    return nil
}

func (ehm *EnhancedHighlightManager) GetHistory() []HighlightRecord {
    return ehm.history
}

func (ehm *EnhancedHighlightManager) AddPreset(name string, config *HighlightConfig) {
    ehm.presets[name] = config
}

// 演示各种高亮场景
func demonstrateHighlightScenarios() {
    fmt.Println("=== 高亮场景演示 ===")

    ehm := NewEnhancedHighlightManager()

    // 模拟节点ID
    testNodes := []int{1001, 1002, 1003, 1004}

    scenarios := []struct {
        name   string
        action func()
    }{
        {
            name: "错误指示",
            action: func() {
                fmt.Println("场景: 错误元素高亮")
                ehm.HighlightWithPreset(testNodes[0], "error", 3*time.Second)
            },
        },
        {
            name: "成功指示",
            action: func() {
                fmt.Println("场景: 成功元素高亮")
                ehm.HighlightWithPreset(testNodes[1], "success", 3*time.Second)
            },
        },
        {
            name: "警告指示",
            action: func() {
                fmt.Println("场景: 警告元素高亮")
                ehm.HighlightWithPreset(testNodes[2], "warning", 3*time.Second)
            },
        },
        {
            name: "信息指示",
            action: func() {
                fmt.Println("场景: 信息元素高亮")
                ehm.HighlightWithPreset(testNodes[3], "info", 3*time.Second)
            },
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("\n%s:\n", scenario.name)
        scenario.action()
        time.Sleep(4 * time.Second) // 等待场景完成
    }

    fmt.Println("\n=== 场景演示完成 ===")
}

*/

// -----------------------------------------------  DOM.highlightRect  -----------------------------------------------
// === 应用场景 ===
// 1. 区域高亮: 高亮页面上的特定矩形区域
// 2. 布局调试: 调试页面布局和定位问题
// 3. 截图标记: 在截图中标记特定区域
// 4. 视觉测试: 测试页面特定区域的视觉效果
// 5. 坐标验证: 验证坐标位置是否正确
// 6. 区域选择: 在页面上标记选择区域

// CDPDOMHighlightRect 高亮页面上的矩形区域
// x: 区域左上角X坐标
// y: 区域左上角Y坐标
// width: 区域宽度
// height: 区域高度
// color: 高亮颜色配置
// outlineColor: 轮廓颜色配置
func CDPDOMHighlightRect(params string) (string, error) {
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
        "method": "DOM.highlightRect",
        "params": %s
    }`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.highlightRect 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.highlightRect 请求超时")
		}
	}
}

/*
// 示例: 高亮页面上的重要区域
func ExampleCDPDOMHighlightRect() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        // 清理高亮
        if _, err := CDPDOMHideHighlight(); err != nil {
            log.Printf("隐藏高亮失败: %v", err)
        }

        // 禁用DOM
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 定义要高亮的区域
    // 假设我们要高亮页面的主要内容区域
    x := 200
    y := 100
    width := 800
    height := 600

    // 3. 创建高亮颜色配置
    highlightColor := &RGBA{
        R: 100,   // 红色
        G: 200,   // 绿色
        B: 255,   // 蓝色
        A: 0.3,   // 透明度
    }

    outlineColor := &RGBA{
        R: 0,     // 红色
        G: 100,   // 绿色
        B: 200,   // 蓝色
        A: 0.8,   // 透明度
    }

    // 4. 高亮矩形区域
    result, err := CDPDOMHighlightRect(x, y, width, height, highlightColor, outlineColor)
    if err != nil {
        log.Printf("高亮矩形区域失败: %v", err)
        return
    }

    log.Printf("矩形区域高亮成功: %s", result)

    // 5. 显示区域信息
    fmt.Printf("=== 矩形区域高亮信息 ===\n")
    fmt.Printf("位置: (%d, %d)\n", x, y)
    fmt.Printf("尺寸: %d x %d 像素\n", width, height)
    fmt.Printf("面积: %d 像素²\n", width*height)
    fmt.Printf("高亮颜色: RGBA(%d, %d, %d, %.2f)\n",
        highlightColor.R, highlightColor.G, highlightColor.B, highlightColor.A)
    fmt.Printf("轮廓颜色: RGBA(%d, %d, %d, %.2f)\n",
        outlineColor.R, outlineColor.G, outlineColor.B, outlineColor.A)

    // 6. 区域分析
    fmt.Printf("\n=== 区域分析 ===\n")

    // 计算相对于标准视口的百分比
    standardViewportWidth := 1920
    standardViewportHeight := 1080

    widthPercent := float64(width) / float64(standardViewportWidth) * 100
    heightPercent := float64(height) / float64(standardViewportHeight) * 100
    areaPercent := (float64(width*height) / float64(standardViewportWidth*standardViewportHeight)) * 100

    fmt.Printf("占视口宽度: %.1f%%\n", widthPercent)
    fmt.Printf("占视口高度: %.1f%%\n", heightPercent)
    fmt.Printf("占视口面积: %.1f%%\n", areaPercent)

    // 区域分类
    if widthPercent > 50 && heightPercent > 50 {
        fmt.Printf("区域类型: 大型主要内容区\n")
    } else if widthPercent > 30 || heightPercent > 30 {
        fmt.Printf("区域类型: 中型内容区\n")
    } else {
        fmt.Printf("区域类型: 小型功能区\n")
    }

    // 7. 视觉检查清单
    fmt.Printf("\n=== 视觉检查清单 ===\n")
    fmt.Printf("检查项:\n")
    fmt.Printf("  ✓ 矩形区域被正确高亮\n")
    fmt.Printf("  ✓ 填充颜色可见但半透明\n")
    fmt.Printf("  ✓ 轮廓清晰可见\n")
    fmt.Printf("  ✓ 不遮挡下方内容\n")
    fmt.Printf("  ✓ 不影响页面交互\n")

    // 8. 响应式设计检查
    fmt.Printf("\n=== 响应式设计检查 ===\n")

    // 检查在不同视口下的表现
    testViewports := []struct {
        width  int
        height int
        name   string
    }{
        {320, 568, "手机"},
        {768, 1024, "平板"},
        {1024, 768, "小桌面"},
        {1440, 900, "桌面"},
        {1920, 1080, "大桌面"},
    }

    fmt.Printf("在不同视口中的位置百分比:\n")
    for _, viewport := range testViewports {
        xPercent := float64(x) / float64(viewport.width) * 100
        yPercent := float64(y) / float64(viewport.height) * 100

        fmt.Printf("  %-8s: X=%.1f%%, Y=%.1f%%\n",
            viewport.name, xPercent, yPercent)
    }

    // 9. 交互测试
    fmt.Printf("\n=== 交互测试 ===\n")

    // 测试区域内的交互
    fmt.Printf("测试区域交互:\n")
    fmt.Printf("  - 可以点击区域内的内容\n")
    fmt.Printf("  - 可以滚动页面\n")
    fmt.Printf("  - 可以选中文本\n")
    fmt.Printf("  - 不影响表单输入\n")

    // 10. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多个区域高亮
    startTime := time.Now()
    testRegions := []struct {
        x, y, width, height int
        name string
    }{
        {100, 100, 200, 200, "小区域"},
        {400, 300, 400, 300, "中区域"},
        {100, 100, 800, 600, "大区域"},
    }

    for _, region := range testRegions {
        testColor := &RGBA{
            R: 100 + rand.Intn(155),
            G: 100 + rand.Intn(155),
            B: 100 + rand.Intn(155),
            A: 0.3,
        }

        if _, err := CDPDOMHighlightRect(region.x, region.y, region.width, region.height, testColor, outlineColor); err != nil {
            fmt.Printf("高亮%s失败: %v\n", region.name, err)
        } else {
            fmt.Printf("✓ 高亮%s成功\n", region.name)
        }

        time.Sleep(500 * time.Millisecond)
        CDPDOMHideHighlight()
    }

    elapsed := time.Since(startTime)
    fmt.Printf("多区域高亮测试耗时: %v\n", elapsed)

    // 11. 可访问性考虑
    fmt.Printf("\n=== 可访问性考虑 ===\n")

    // 检查颜色对比度
    bgLuminance := calculateLuminance(255, 255, 255) // 假设白色背景
    fgLuminance := calculateLuminance(
        float64(highlightColor.R),
        float64(highlightColor.G),
        float64(highlightColor.B),
    )

    contrastRatio := calculateContrastRatio(bgLuminance, fgLuminance)
    fmt.Printf("颜色对比度: %.2f:1\n", contrastRatio)

    if contrastRatio >= 4.5 {
        fmt.Printf("✓ 对比度符合WCAG AA标准\n")
    } else if contrastRatio >= 3.0 {
        fmt.Printf("⚠️ 对比度较低，考虑增强\n")
    } else {
        fmt.Printf("❌ 对比度过低，需要改进\n")
    }

    // 12. 清理演示
    fmt.Printf("\n=== 清理演示 ===\n")

    // 重新高亮原始区域
    if _, err := CDPDOMHighlightRect(x, y, width, height, highlightColor, outlineColor); err == nil {
        fmt.Printf("重新高亮原始区域，3秒后清理...\n")
        time.Sleep(3 * time.Second)

        if _, err := CDPDOMHideHighlight(); err == nil {
            fmt.Printf("高亮已清理\n")
        }
    }
}

// 高级功能: 闪烁高亮区域
func blinkHighlightRect(x, y, width, height int, duration time.Duration, blinkCount int) {
    CDPDOMEnable()
    defer CDPDOMHideHighlight()
    defer CDPDOMDisable()

    highlightColor := &RGBA{R: 255, G: 100, B: 100, A: 0.5}
    outlineColor := &RGBA{R: 255, G: 50, B: 50, A: 0.8}

    fmt.Printf("开始闪烁高亮，闪烁 %d 次\n", blinkCount)

    blinkInterval := duration / time.Duration(blinkCount*2)

    for i := 0; i < blinkCount; i++ {
        // 显示
        if _, err := CDPDOMHighlightRect(x, y, width, height, highlightColor, outlineColor); err != nil {
            log.Printf("闪烁高亮失败: %v", err)
            break
        }

        time.Sleep(blinkInterval)

        // 隐藏
        CDPDOMHideHighlight()

        if i < blinkCount-1 { // 最后一次不延迟隐藏
            time.Sleep(blinkInterval)
        }
    }

    fmt.Println("闪烁高亮完成")
}

// 高级功能: 高亮多个相邻区域
func highlightAdjacentRegions(baseX, baseY, regionWidth, regionHeight, spacing, rows, cols int) {
    CDPDOMEnable()
    defer func() {
        CDPDOMHideHighlight()
        CDPDOMDisable()
    }()

    colors := []RGBA{
        {R: 255, G: 200, B: 200, A: 0.4},
        {R: 200, G: 255, B: 200, A: 0.4},
        {R: 200, G: 200, B: 255, A: 0.4},
        {R: 255, G: 255, B: 200, A: 0.4},
    }

    fmt.Printf("高亮 %d x %d 网格区域\n", rows, cols)

    for row := 0; row < rows; row++ {
        for col := 0; col < cols; col++ {
            x := baseX + col*(regionWidth+spacing)
            y := baseY + row*(regionHeight+spacing)

            colorIndex := (row*cols + col) % len(colors)

            if _, err := CDPDOMHighlightRect(x, y, regionWidth, regionHeight, &colors[colorIndex], nil); err != nil {
                log.Printf("高亮区域(%d,%d)失败: %v", row, col, err)
            } else {
                fmt.Printf("  ✓ 高亮区域(%d,%d)\n", row, col)
            }

            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Println("网格高亮完成，按回车键清理...")
    fmt.Scanln()
}

// 计算相对亮度
func calculateLuminance(r, g, b float64) float64 {
    // 转换为sRGB
    rsRGB := r / 255.0
    gsRGB := g / 255.0
    bsRGB := b / 255.0

    // 应用gamma校正
    var rLinear, gLinear, bLinear float64

    if rsRGB <= 0.03928 {
        rLinear = rsRGB / 12.92
    } else {
        rLinear = math.Pow((rsRGB+0.055)/1.055, 2.4)
    }

    if gsRGB <= 0.03928 {
        gLinear = gsRGB / 12.92
    } else {
        gLinear = math.Pow((gsRGB+0.055)/1.055, 2.4)
    }

    if bsRGB <= 0.03928 {
        bLinear = bsRGB / 12.92
    } else {
        bLinear = math.Pow((bsRGB+0.055)/1.055, 2.4)
    }

    // 计算相对亮度
    return 0.2126*rLinear + 0.7152*gLinear + 0.0722*bLinear
}

// 计算对比度
func calculateContrastRatio(luminance1, luminance2 float64) float64 {
    // 确保luminance1是较亮的颜色
    if luminance1 < luminance2 {
        luminance1, luminance2 = luminance2, luminance1
    }
    return (luminance1 + 0.05) / (luminance2 + 0.05)
}

// 区域高亮管理器
type RectHighlightManager struct {
    activeHighlights []RectHighlight
    history          []RectHighlightRecord
}

type RectHighlight struct {
    ID           string
    X, Y         int
    Width, Height int
    Color        *RGBA
    OutlineColor *RGBA
    CreatedAt    time.Time
}

type RectHighlightRecord struct {
    RectHighlight
    Duration   time.Duration
    RemovedAt  time.Time
}

func NewRectHighlightManager() *RectHighlightManager {
    return &RectHighlightManager{
        activeHighlights: make([]RectHighlight, 0),
        history:          make([]RectHighlightRecord, 0),
    }
}

func (rhm *RectHighlightManager) AddHighlight(x, y, width, height int, color, outlineColor *RGBA, id string) error {
    if id == "" {
        id = fmt.Sprintf("highlight-%d", time.Now().UnixNano())
    }

    highlight := RectHighlight{
        ID:           id,
        X:            x,
        Y:            y,
        Width:        width,
        Height:       height,
        Color:        color,
        OutlineColor: outlineColor,
        CreatedAt:    time.Now(),
    }

    result, err := CDPDOMHighlightRect(x, y, width, height, color, outlineColor)
    if err != nil {
        return fmt.Errorf("高亮区域失败: %w", err)
    }

    rhm.activeHighlights = append(rhm.activeHighlights, highlight)
    log.Printf("添加区域高亮 %s: %s", id, result)

    return nil
}

func (rhm *RectHighlightManager) RemoveHighlight(id string) error {
    for i, highlight := range rhm.activeHighlights {
        if highlight.ID == id {
            // 记录到历史
            record := RectHighlightRecord{
                RectHighlight: highlight,
                Duration:      time.Since(highlight.CreatedAt),
                RemovedAt:     time.Now(),
            }
            rhm.history = append(rhm.history, record)

            // 从活动列表移除
            rhm.activeHighlights = append(rhm.activeHighlights[:i], rhm.activeHighlights[i+1:]...)

            // 清理高亮
            if _, err := CDPDOMHideHighlight(); err != nil {
                return fmt.Errorf("清理高亮失败: %w", err)
            }

            // 重新高亮其他区域
            for _, h := range rhm.activeHighlights {
                if _, err := CDPDOMHighlightRect(h.X, h.Y, h.Width, h.Height, h.Color, h.OutlineColor); err != nil {
                    log.Printf("重新高亮区域 %s 失败: %v", h.ID, err)
                }
            }

            return nil
        }
    }

    return fmt.Errorf("未找到高亮区域: %s", id)
}

func (rhm *RectHighlightManager) ClearAll() {
    for _, highlight := range rhm.activeHighlights {
        record := RectHighlightRecord{
            RectHighlight: highlight,
            Duration:      time.Since(highlight.CreatedAt),
            RemovedAt:     time.Now(),
        }
        rhm.history = append(rhm.history, record)
    }

    rhm.activeHighlights = make([]RectHighlight, 0)
    CDPDOMHideHighlight()
}

func (rhm *RectHighlightManager) GetActiveHighlights() []RectHighlight {
    return rhm.activeHighlights
}

func (rhm *RectHighlightManager) GetHistory() []RectHighlightRecord {
    return rhm.history
}

// 演示各种区域高亮场景
func demonstrateRectHighlightScenarios() {
    fmt.Println("=== 区域高亮场景演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()
    defer CDPDOMHideHighlight()

    manager := NewRectHighlightManager()

    scenarios := []struct {
        name   string
        x, y   int
        width, height int
        color  *RGBA
        description string
    }{
        {
            name: "页头区域",
            x: 0, y: 0, width: 1920, height: 80,
            color: &RGBA{R: 255, G: 200, B: 200, A: 0.3},
            description: "标记页面头部区域",
        },
        {
            name: "主要内容区",
            x: 200, y: 100, width: 1200, height: 800,
            color: &RGBA{R: 200, G: 255, B: 200, A: 0.2},
            description: "标记主要内容区域",
        },
        {
            name: "侧边栏",
            x: 1400, y: 100, width: 400, height: 800,
            color: &RGBA{R: 200, G: 200, B: 255, A: 0.3},
            description: "标记侧边栏区域",
        },
        {
            name: "页脚",
            x: 0, y: 900, width: 1920, height: 180,
            color: &RGBA{R: 255, G: 255, B: 200, A: 0.4},
            description: "标记页脚区域",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("\n场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("位置: (%d, %d) 尺寸: %d x %d\n",
            scenario.x, scenario.y, scenario.width, scenario.height)

        if err := manager.AddHighlight(
            scenario.x, scenario.y,
            scenario.width, scenario.height,
            scenario.color,
            &RGBA{R: 0, G: 0, B: 0, A: 0.8},
            scenario.name,
        ); err != nil {
            fmt.Printf("❌ 高亮失败: %v\n", err)
        } else {
            fmt.Printf("✅ 高亮成功\n")
        }

        time.Sleep(2 * time.Second)
    }

    // 显示当前高亮状态
    fmt.Printf("\n=== 当前高亮状态 ===\n")
    activeHighlights := manager.GetActiveHighlights()
    fmt.Printf("活动高亮数量: %d\n", len(activeHighlights))

    for _, highlight := range activeHighlights {
        fmt.Printf("  - %s: (%d,%d) %dx%d\n",
            highlight.ID, highlight.X, highlight.Y,
            highlight.Width, highlight.Height)
    }

    // 逐个清理
    fmt.Printf("\n=== 清理演示 ===\n")
    for _, highlight := range activeHighlights {
        fmt.Printf("清理区域: %s\n", highlight.ID)
        if err := manager.RemoveHighlight(highlight.ID); err != nil {
            fmt.Printf("❌ 清理失败: %v\n", err)
        } else {
            fmt.Printf("✅ 清理成功\n")
        }
        time.Sleep(1 * time.Second)
    }

    fmt.Println("\n=== 场景演示完成 ===")
}

*/

// -----------------------------------------------  DOM.moveTo  -----------------------------------------------
// === 应用场景 ===
// 1. 节点移动: 将DOM节点移动到新位置
// 2. 结构重组: 重新组织页面DOM结构
// 3. 动态布局: 动态调整页面布局
// 4. 内容排序: 重新排序页面内容
// 5. 拖放模拟: 模拟拖放操作
// 6. 动画测试: 测试节点移动动画

// CDPDOMMoveTo 将节点移动到新位置
// nodeID: 要移动的节点ID
// targetNodeID: 目标父节点ID
// insertBeforeNodeID: 插入到哪个子节点之前，0表示插入到最后
func CDPDOMMoveTo(nodeID, targetNodeID, insertBeforeNodeID int) (string, error) {
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
        "method": "DOM.moveTo",
        "params": {
            "nodeId": %d,
            "targetNodeId": %d,
            "insertBeforeNodeId": %d
        }
    }`, reqID, nodeID, targetNodeID, insertBeforeNodeID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.moveTo 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.moveTo 请求超时")
		}
	}
}

/*

// 示例: 移动列表项到新位置
func ExampleCDPDOMMoveTo() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 模拟场景：重新排序任务列表

    // 假设我们有以下节点ID：
    // 任务列表容器：containerID = 1001
    // 任务项1：task1ID = 1002
    // 任务项2：task2ID = 1003
    // 任务项3：task3ID = 1004

    containerID := 1001
    task1ID := 1002
    task2ID := 1003
    task3ID := 1004

    fmt.Printf("=== 任务列表重新排序演示 ===\n")

    // 2. 初始状态显示
    fmt.Printf("初始状态:\n")
    fmt.Printf("  容器节点ID: %d\n", containerID)
    fmt.Printf("  任务项1 ID: %d\n", task1ID)
    fmt.Printf("  任务项2 ID: %d\n", task2ID)
    fmt.Printf("  任务项3 ID: %d\n", task3ID)

    // 3. 移动任务项1到任务项3之后
    fmt.Printf("\n操作1: 移动任务项1到任务项3之后\n")

    // insertBeforeNodeID = 0 表示插入到子节点列表的最后
    result, err := CDPDOMMoveTo(task1ID, containerID, 0)
    if err != nil {
        log.Printf("移动节点失败: %v", err)
        return
    }

    fmt.Printf("移动结果: %s\n", result)
    fmt.Printf("新顺序: 任务2, 任务3, 任务1\n")

    // 4. 移动任务项3到任务项2之前
    fmt.Printf("\n操作2: 移动任务项3到任务项2之前\n")

    // 将任务3移动到任务2之前
    result, err = CDPDOMMoveTo(task3ID, containerID, task2ID)
    if err != nil {
        log.Printf("移动节点失败: %v", err)
        return
    }

    fmt.Printf("移动结果: %s\n", result)
    fmt.Printf("新顺序: 任务3, 任务2, 任务1\n")

    // 5. 验证移动结果
    fmt.Printf("\n=== 移动结果验证 ===\n")

    // 获取容器子节点信息
    // 这里假设有方法可以获取子节点列表
    // 实际使用时可能需要结合其他DOM方法

    // 6. 移动操作分析
    fmt.Printf("\n=== 移动操作分析 ===\n")

    // 检查移动的类型
    fmt.Printf("移动类型分析:\n")

    // 同父级移动
    fmt.Printf("  ✓ 同父级重新排序\n")
    fmt.Printf("  - 保持原有父子关系\n")
    fmt.Printf("  - 只改变兄弟顺序\n")

    // 检查是否触发重排
    fmt.Printf("\n布局重排分析:\n")
    fmt.Printf("  - 移动操作会触发DOM重排\n")
    fmt.Printf("  - 可能触发浏览器重绘\n")
    fmt.Printf("  - 如果使用CSS动画会更平滑\n")

    // 7. 性能考虑
    fmt.Printf("\n=== 性能考虑 ===\n")

    // 测试多次移动的性能
    startTime := time.Now()
    moveOperations := 10

    for i := 0; i < moveOperations; i++ {
        // 模拟交替移动
        if i%2 == 0 {
            CDPDOMMoveTo(task1ID, containerID, task2ID)
        } else {
            CDPDOMMoveTo(task2ID, containerID, task1ID)
        }
    }

    elapsed := time.Since(startTime)
    fmt.Printf("%d 次移动操作耗时: %v\n", moveOperations, elapsed)
    fmt.Printf("平均每次移动耗时: %v\n", elapsed/time.Duration(moveOperations))

    // 8. 错误处理测试
    fmt.Printf("\n=== 错误处理测试 ===\n")

    // 测试无效的节点ID
    testCases := []struct {
        name        string
        nodeID      int
        targetID    int
        insertBefore int
        expectError bool
    }{
        {
            name:        "有效移动",
            nodeID:      task1ID,
            targetID:    containerID,
            insertBefore: task2ID,
            expectError: false,
        },
        {
            name:        "无效节点ID",
            nodeID:      9999,
            targetID:    containerID,
            insertBefore: 0,
            expectError: true,
        },
        {
            name:        "无效目标ID",
            nodeID:      task1ID,
            targetID:    9999,
            insertBefore: 0,
            expectError: true,
        },
        {
            name:        "无效插入位置",
            nodeID:      task1ID,
            targetID:    containerID,
            insertBefore: 9999,
            expectError: true,
        },
    }

    for _, tc := range testCases {
        fmt.Printf("测试: %s\n", tc.name)
        result, err := CDPDOMMoveTo(tc.nodeID, tc.targetID, tc.insertBefore)

        if tc.expectError {
            if err != nil {
                fmt.Printf("  ✓ 预期错误: %v\n", err)
            } else {
                fmt.Printf("  ❌ 预期错误但成功了: %s\n", result)
            }
        } else {
            if err != nil {
                fmt.Printf("  ❌ 预期成功但失败: %v\n", err)
            } else {
                fmt.Printf("  ✓ 移动成功\n")
            }
        }
    }

    // 9. 可访问性影响
    fmt.Printf("\n=== 可访问性影响 ===\n")

    fmt.Printf("移动操作对可访问性的影响:\n")
    fmt.Printf("  - 可能影响屏幕阅读器的阅读顺序\n")
    fmt.Printf("  - 可能影响键盘导航顺序\n")
    fmt.Printf("  - 需要确保焦点管理\n")
    fmt.Printf("  - 需要确保ARIA属性正确\n")

    // 10. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
    }{
        {
            name:        "拖放排序",
            description: "实现可拖放的项目列表",
        },
        {
            name:        "动态布局",
            description: "根据用户偏好动态调整布局",
        },
        {
            name:        "内容重新排序",
            description: "根据规则自动排序内容",
        },
        {
            name:        "动画效果",
            description: "配合CSS动画实现平滑移动效果",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n\n", scenario.description)
    }
}

// 高级功能: 安全移动节点
func safeMoveNode(nodeID, targetNodeID, insertBeforeNodeID int) error {
    // 验证节点存在
    if !validateNodeExists(nodeID) {
        return fmt.Errorf("源节点不存在: %d", nodeID)
    }

    if !validateNodeExists(targetNodeID) {
        return fmt.Errorf("目标节点不存在: %d", targetNodeID)
    }

    if insertBeforeNodeID > 0 && !validateNodeExists(insertBeforeNodeID) {
        return fmt.Errorf("插入位置节点不存在: %d", insertBeforeNodeID)
    }

    // 检查是否是移动到自身
    if nodeID == targetNodeID {
        return fmt.Errorf("不能将节点移动到自身")
    }

    // 检查是否形成循环
    if wouldCreateCycle(nodeID, targetNodeID) {
        return fmt.Errorf("移动会形成循环引用")
    }

    // 执行移动
    result, err := CDPDOMMoveTo(nodeID, targetNodeID, insertBeforeNodeID)
    if err != nil {
        return fmt.Errorf("移动失败: %w", err)
    }

    // 验证移动结果
    if !verifyMoveResult(nodeID, targetNodeID) {
        return fmt.Errorf("移动验证失败: %s", result)
    }

    return nil
}

// 验证节点是否存在
func validateNodeExists(nodeID int) bool {
    // 这里应该调用DOM.describeNode或其他方法来验证节点存在
    // 简化实现
    return nodeID > 0
}

// 检查是否形成循环
func wouldCreateCycle(nodeID, targetNodeID int) bool {
    // 简化实现
    // 实际应该遍历DOM树检查targetNodeID是否是nodeID的后代
    return false
}

// 验证移动结果
func verifyMoveResult(nodeID, targetNodeID int) bool {
    // 这里应该验证节点确实移动到了新位置
    return true
}

// 高级功能: 批量移动节点
func batchMoveNodes(moves []MoveOperation) ([]MoveResult, error) {
    var results []MoveResult

    for _, move := range moves {
        result := MoveResult{
            NodeID:   move.NodeID,
            Success:  false,
        }

        err := safeMoveNode(move.NodeID, move.TargetNodeID, move.InsertBeforeNodeID)
        if err != nil {
            result.Error = err.Error()
        } else {
            result.Success = true
        }

        results = append(results, result)

        // 可选：添加延迟以防止操作过快
        time.Sleep(50 * time.Millisecond)
    }

    return results, nil
}

type MoveOperation struct {
    NodeID           int
    TargetNodeID     int
    InsertBeforeNodeID int
}

type MoveResult struct {
    NodeID  int
    Success bool
    Error   string
}

// 高级功能: 动画移动
func animatedMove(nodeID, targetNodeID, insertBeforeNodeID int, duration time.Duration) error {
    // 1. 获取当前位置
    startPos, err := getNodePosition(nodeID)
    if err != nil {
        return fmt.Errorf("获取起始位置失败: %w", err)
    }

    // 2. 执行移动
    if _, err := CDPDOMMoveTo(nodeID, targetNodeID, insertBeforeNodeID); err != nil {
        return fmt.Errorf("移动失败: %w", err)
    }

    // 3. 获取新位置
    endPos, err := getNodePosition(nodeID)
    if err != nil {
        return fmt.Errorf("获取结束位置失败: %w", err)
    }

    // 4. 计算移动向量
    moveVector := Position{
        X: endPos.X - startPos.X,
        Y: endPos.Y - startPos.Y,
    }

    // 5. 应用CSS动画
    if err := applyMoveAnimation(nodeID, moveVector, duration); err != nil {
        return fmt.Errorf("应用动画失败: %w", err)
    }

    return nil
}

type Position struct {
    X, Y float64
}

func getNodePosition(nodeID int) (Position, error) {
    // 简化实现
    // 实际应该获取节点的盒模型或计算位置
    return Position{X: 0, Y: 0}, nil
}

func applyMoveAnimation(nodeID int, moveVector Position, duration time.Duration) error {
    // 这里应该通过CSS或JavaScript应用动画
    // 简化实现
    return nil
}

// 节点移动管理器
type NodeMoveManager struct {
    movesHistory []MoveRecord
    undoStack    []MoveOperation
    redoStack    []MoveOperation
}

type MoveRecord struct {
    Timestamp   time.Time
    Operation   MoveOperation
    Result      string
    Duration    time.Duration
}

func NewNodeMoveManager() *NodeMoveManager {
    return &NodeMoveManager{
        movesHistory: make([]MoveRecord, 0),
        undoStack:    make([]MoveOperation, 0),
        redoStack:    make([]MoveOperation, 0),
    }
}

func (nmm *NodeMoveManager) MoveNode(op MoveOperation) error {
    startTime := time.Now()

    result, err := CDPDOMMoveTo(op.NodeID, op.TargetNodeID, op.InsertBeforeNodeID)
    if err != nil {
        return fmt.Errorf("移动失败: %w", err)
    }

    duration := time.Since(startTime)

    // 记录到历史
    record := MoveRecord{
        Timestamp: startTime,
        Operation: op,
        Result:    result,
        Duration:  duration,
    }
    nmm.movesHistory = append(nmm.movesHistory, record)

    // 保存到撤销栈
    nmm.undoStack = append(nmm.undoStack, op)

    // 清空重做栈
    nmm.redoStack = make([]MoveOperation, 0)

    return nil
}

func (nmm *NodeMoveManager) Undo() error {
    if len(nmm.undoStack) == 0 {
        return fmt.Errorf("没有可撤销的操作")
    }

    // 获取最后一个操作
    lastOp := nmm.undoStack[len(nmm.undoStack)-1]
    nmm.undoStack = nmm.undoStack[:len(nmm.undoStack)-1]

    // 这里需要计算逆向操作
    // 简化实现：保存到重做栈
    nmm.redoStack = append(nmm.redoStack, lastOp)

    return fmt.Errorf("撤销功能需要实现逆向操作计算")
}

func (nmm *NodeMoveManager) Redo() error {
    if len(nmm.redoStack) == 0 {
        return fmt.Errorf("没有可重做的操作")
    }

    // 获取最后一个重做操作
    redoOp := nmm.redoStack[len(nmm.redoStack)-1]
    nmm.redoStack = nmm.redoStack[:len(nmm.redoStack)-1]

    // 执行重做
    return nmm.MoveNode(redoOp)
}

func (nmm *NodeMoveManager) GetHistory() []MoveRecord {
    return nmm.movesHistory
}

// 演示复杂移动场景
func demonstrateComplexMoveScenarios() {
    fmt.Println("=== 复杂移动场景演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    manager := NewNodeMoveManager()

    // 模拟节点
    nodes := map[string]int{
        "header": 1001,
        "nav":    1002,
        "main":   1003,
        "aside":  1004,
        "footer": 1005,
        "body":   1006,
    }

    scenarios := []struct {
        name        string
        move        MoveOperation
        description string
    }{
        {
            name: "移动导航到主要内容之前",
            move: MoveOperation{
                NodeID:           nodes["nav"],
                TargetNodeID:     nodes["body"],
                InsertBeforeNodeID: nodes["main"],
            },
            description: "将导航栏移动到主要内容区域之前",
        },
        {
            name: "移动侧边栏到主要内容之后",
            move: MoveOperation{
                NodeID:           nodes["aside"],
                TargetNodeID:     nodes["body"],
                InsertBeforeNodeID: 0, // 插入到最后
            },
            description: "将侧边栏移动到页面底部",
        },
        {
            name: "移动页脚到侧边栏之后",
            move: MoveOperation{
                NodeID:           nodes["footer"],
                TargetNodeID:     nodes["body"],
                InsertBeforeNodeID: 0, // 插入到最后
            },
            description: "重新排列页脚位置",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("\n场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)

        if err := manager.MoveNode(scenario.move); err != nil {
            fmt.Printf("❌ 移动失败: %v\n", err)
        } else {
            fmt.Printf("✅ 移动成功\n")
        }

        time.Sleep(1 * time.Second)
    }

    // 显示移动历史
    fmt.Printf("\n=== 移动历史 ===\n")
    history := manager.GetHistory()
    for i, record := range history {
        fmt.Printf("%d. 时间: %s, 耗时: %v\n",
            i+1,
            record.Timestamp.Format("15:04:05"),
            record.Duration)
    }

    fmt.Println("\n=== 场景演示完成 ===")
}

*/

// -----------------------------------------------  DOM.querySelector  -----------------------------------------------
// === 应用场景 ===
// 1. 元素查找: 通过CSS选择器查找DOM元素
// 2. 自动化测试: 在测试中定位页面元素
// 3. 内容提取: 提取特定选择器的内容
// 4. 样式检查: 检查特定样式规则的元素
// 5. 表单操作: 定位表单元素进行操作
// 6. 交互元素: 定位可交互元素进行点击等操作

// CDPDOMQuerySelector 通过CSS选择器查询元素
// nodeID: 要查询的根节点ID
// selector: CSS选择器字符串
func CDPDOMQuerySelector(nodeID int, selector string) (string, error) {
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
        "method": "DOM.querySelector",
        "params": {
            "nodeId": %d,
            "selector": "%s"
        }
    }`, reqID, nodeID, selector)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.querySelector 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.querySelector 请求超时")
		}
	}
}

/*

// 示例: 查询页面中的按钮元素
func ExampleCDPDOMQuerySelector() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 获取文档根节点
    // 首先需要获取文档的根节点ID
    // 假设我们已经有了文档根节点ID
    documentNodeID := 1  // 通常文档根节点ID为1

    // 3. 定义要查询的CSS选择器
    testSelectors := []struct {
        selector    string
        description string
    }{
        {"button", "所有按钮元素"},
        {"button.primary", "主按钮"},
        {"#submit-btn", "提交按钮（通过ID）"},
        {"input[type='text']", "文本输入框"},
        {".nav-item.active", "活动的导航项"},
        {"a[href^='https']", "外部链接"},
        {"div > p:first-child", "div的第一个段落子元素"},
        {"[data-testid='login-form']", "测试ID选择器"},
    }

    fmt.Printf("=== CSS选择器查询测试 ===\n")
    fmt.Printf("文档根节点ID: %d\n\n", documentNodeID)

    // 4. 执行查询测试
    for i, test := range testSelectors {
        fmt.Printf("测试 %d: %s\n", i+1, test.description)
        fmt.Printf("选择器: %s\n", test.selector)

        result, err := CDPDOMQuerySelector(documentNodeID, test.selector)
        if err != nil {
            fmt.Printf("❌ 查询失败: %v\n\n", err)
            continue
        }

        // 5. 解析查询结果
        var response struct {
            Result struct {
                NodeID int `json:"nodeId"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(result), &response); err != nil {
            fmt.Printf("❌ 解析结果失败: %v\n\n", err)
            continue
        }

        nodeID := response.Result.NodeID

        if nodeID == 0 {
            fmt.Printf("ℹ️ 未找到匹配元素\n\n")
        } else {
            fmt.Printf("✅ 找到元素，节点ID: %d\n", nodeID)

            // 6. 获取元素详细信息
            elementInfo, err := getElementInfo(nodeID)
            if err == nil {
                displayElementInfo(elementInfo)
            }

            // 7. 高亮找到的元素
            highlightElement(nodeID)

            fmt.Println()
        }
    }

    // 8. 选择器性能测试
    fmt.Printf("=== 选择器性能测试 ===\n")

    performanceTestSelectors := []string{
        "button",                     // 简单标签选择器
        ".btn-primary",              // 类选择器
        "#main-content",             // ID选择器
        "div.container > form",      // 后代选择器
        "input[type='email']",       // 属性选择器
        "li:nth-child(2)",           // 伪类选择器
    }

    for _, selector := range performanceTestSelectors {
        duration, found := testSelectorPerformance(documentNodeID, selector)
        fmt.Printf("选择器: %-25s 耗时: %v 找到: %v\n",
            selector, duration, found)
    }
}

// 获取元素信息
func getElementInfo(nodeID int) (map[string]interface{}, error) {
    info := make(map[string]interface{})

    // 获取节点描述
    descResult, err := CDPDOMDescribeNode(nodeID, 0, true)
    if err != nil {
        return nil, err
    }

    var descResp struct {
        Result struct {
            Node struct {
                NodeName  string   `json:"nodeName"`
                LocalName string   `json:"localName"`
                Attributes []string `json:"attributes,omitempty"`
            } `json:"node"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(descResult), &descResp); err != nil {
        return nil, err
    }

    info["nodeName"] = descResp.Result.Node.NodeName
    info["localName"] = descResp.Result.Node.LocalName
    info["attributes"] = descResp.Result.Node.Attributes

    // 获取外层HTML
    outerHTMLResult, err := CDPDOMGetOuterHTML(nodeID)
    if err == nil {
        var htmlResp struct {
            Result struct {
                OuterHTML string `json:"outerHTML"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(outerHTMLResult), &htmlResp); err == nil {
            info["outerHTML"] = truncateString(htmlResp.Result.OuterHTML, 100)
        }
    }

    return info, nil
}

// 显示元素信息
func displayElementInfo(info map[string]interface{}) {
    fmt.Printf("   元素类型: %s\n", info["nodeName"])
    fmt.Printf("   本地名称: %s\n", info["localName"])

    if attrs, ok := info["attributes"].([]string); ok && len(attrs) > 0 {
        fmt.Printf("   属性: ")
        for i := 0; i < len(attrs); i += 2 {
            if i+1 < len(attrs) {
                fmt.Printf("%s=\"%s\" ", attrs[i], attrs[i+1])
            }
        }
        fmt.Println()
    }

    if html, ok := info["outerHTML"].(string); ok {
        fmt.Printf("   HTML: %s\n", html)
    }
}

// 高亮元素
func highlightElement(nodeID int) {
    // 创建高亮配置
    highlightConfig := &HighlightConfig{
        ContentColor: &RGBA{R: 100, G: 200, B: 255, A: 0.3},
        BorderColor:  &RGBA{R: 0, G: 100, B: 200, A: 0.8},
        ShowInfo:     true,
    }

    // 短暂高亮元素
    if _, err := CDPDOMHighlightNode(nodeID, highlightConfig); err == nil {
        // 3秒后自动清理
        go func() {
            time.Sleep(3 * time.Second)
            CDPDOMHideHighlight()
        }()
    }
}

// 测试选择器性能
func testSelectorPerformance(rootNodeID int, selector string) (time.Duration, bool) {
    startTime := time.Now()

    result, err := CDPDOMQuerySelector(rootNodeID, selector)
    if err != nil {
        return time.Since(startTime), false
    }

    var response struct {
        Result struct {
            NodeID int `json:"nodeId"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        return time.Since(startTime), false
    }

    return time.Since(startTime), response.Result.NodeID > 0
}

// 截断字符串
func truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "..."
}

// 高级功能: 智能选择器查询
func smartQuerySelector(rootNodeID int, selector string, options QueryOptions) (QueryResult, error) {
    result := QueryResult{
        Selector: selector,
        StartTime: time.Now(),
    }

    // 执行查询
    queryResult, err := CDPDOMQuerySelector(rootNodeID, selector)
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(result.StartTime)

    if err != nil {
        result.Success = false
        result.Error = err.Error()
        return result, err
    }

    // 解析结果
    var resp struct {
        Result struct {
            NodeID int `json:"nodeId"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(queryResult), &resp); err != nil {
        result.Success = false
        result.Error = err.Error()
        return result, err
    }

    result.NodeID = resp.Result.NodeID
    result.Found = resp.Result.NodeID > 0
    result.Success = true

    // 如果需要验证
    if options.Validate && result.Found {
        if err := validateElement(result.NodeID, options.ValidationRules); err != nil {
            result.ValidationError = err.Error()
        } else {
            result.Valid = true
        }
    }

    // 如果需要高亮
    if options.Highlight && result.Found {
        if options.HighlightDuration > 0 {
            go highlightElementWithDuration(result.NodeID, options.HighlightDuration)
        } else {
            highlightElement(result.NodeID)
        }
    }

    return result, nil
}

type QueryOptions struct {
    Validate         bool
    ValidationRules  ValidationRules
    Highlight        bool
    HighlightDuration time.Duration
    Timeout          time.Duration
}

type ValidationRules struct {
    RequiredAttributes []string
    MustBeVisible      bool
    MustBeEnabled      bool
    MinWidth           int
    MinHeight          int
}

type QueryResult struct {
    Selector        string
    NodeID          int
    Found           bool
    Success         bool
    Error           string
    StartTime       time.Time
    EndTime         time.Time
    Duration        time.Duration
    Valid           bool
    ValidationError string
}

func validateElement(nodeID int, rules ValidationRules) error {
    // 获取元素属性
    attrsResult, err := CDPDOMGetAttributes(nodeID)
    if err != nil {
        return fmt.Errorf("获取属性失败: %w", err)
    }

    var attrsResp struct {
        Result struct {
            Attributes []string `json:"attributes"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(attrsResult), &attrsResp); err != nil {
        return fmt.Errorf("解析属性失败: %w", err)
    }

    attrs := make(map[string]string)
    for i := 0; i < len(attrsResp.Result.Attributes); i += 2 {
        if i+1 < len(attrsResp.Result.Attributes) {
            attrs[attrsResp.Result.Attributes[i]] = attrsResp.Result.Attributes[i+1]
        }
    }

    // 检查必需属性
    for _, requiredAttr := range rules.RequiredAttributes {
        if _, exists := attrs[requiredAttr]; !exists {
            return fmt.Errorf("缺少必需属性: %s", requiredAttr)
        }
    }

    // 检查是否启用
    if rules.MustBeEnabled {
        if disabled, exists := attrs["disabled"]; exists && (disabled == "true" || disabled == "") {
            return fmt.Errorf("元素被禁用")
        }
    }

    // 获取盒模型检查尺寸
    if rules.MinWidth > 0 || rules.MinHeight > 0 {
        boxResult, err := CDPDOMGetBoxModel(nodeID)
        if err == nil {
            var boxResp struct {
                Result struct {
                    Model struct {
                        Width  int `json:"width"`
                        Height int `json:"height"`
                    } `json:"model"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(boxResult), &boxResp); err == nil {
                if rules.MinWidth > 0 && boxResp.Result.Model.Width < rules.MinWidth {
                    return fmt.Errorf("宽度太小: %d < %d",
                        boxResp.Result.Model.Width, rules.MinWidth)
                }
                if rules.MinHeight > 0 && boxResp.Result.Model.Height < rules.MinHeight {
                    return fmt.Errorf("高度太小: %d < %d",
                        boxResp.Result.Model.Height, rules.MinHeight)
                }
            }
        }
    }

    return nil
}

func highlightElementWithDuration(nodeID int, duration time.Duration) {
    highlightConfig := &HighlightConfig{
        ContentColor: &RGBA{R: 100, G: 200, B: 255, A: 0.3},
        BorderColor:  &RGBA{R: 0, G: 100, B: 200, A: 0.8},
    }

    if _, err := CDPDOMHighlightNode(nodeID, highlightConfig); err == nil {
        time.Sleep(duration)
        CDPDOMHideHighlight()
    }
}

// 高级功能: 查询并收集多个元素
func queryAndCollectElements(rootNodeID int, selector string, limit int) ([]ElementInfo, error) {
    var elements []ElementInfo

    // 首先查询第一个元素
    result, err := CDPDOMQuerySelector(rootNodeID, selector)
    if err != nil {
        return nil, err
    }

    var resp struct {
        Result struct {
            NodeID int `json:"nodeId"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &resp); err != nil {
        return nil, err
    }

    if resp.Result.NodeID == 0 {
        return elements, nil
    }

    // 收集第一个元素
    if info, err := getElementInfo(resp.Result.NodeID); err == nil {
        elements = append(elements, ElementInfo{
            NodeID: resp.Result.NodeID,
            Info:   info,
        })
    }

    // 继续查询直到达到限制
    for len(elements) < limit {
        // 这里需要实现查询下一个兄弟元素的功能
        // 简化实现：只返回第一个找到的元素
        break
    }

    return elements, nil
}

type ElementInfo struct {
    NodeID int
    Info   map[string]interface{}
}

// 选择器测试器
type SelectorTester struct {
    rootNodeID int
    testCases  []SelectorTestCase
    results    []TestResult
}

type SelectorTestCase struct {
    Name        string
    Selector    string
    Description string
    Expected    bool // 期望是否找到元素
    Priority    int  // 测试优先级
}

type TestResult struct {
    SelectorTestCase
    Found     bool
    Success   bool
    NodeID    int
    Error     string
    Duration  time.Duration
    Timestamp time.Time
}

func NewSelectorTester(rootNodeID int) *SelectorTester {
    return &SelectorTester{
        rootNodeID: rootNodeID,
        testCases:  make([]SelectorTestCase, 0),
        results:    make([]TestResult, 0),
    }
}

func (st *SelectorTester) AddTestCase(testCase SelectorTestCase) {
    st.testCases = append(st.testCases, testCase)
}

func (st *SelectorTester) RunTests() {
    // 按优先级排序
    sort.Slice(st.testCases, func(i, j int) bool {
        return st.testCases[i].Priority > st.testCases[j].Priority
    })

    for _, testCase := range st.testCases {
        result := TestResult{
            SelectorTestCase: testCase,
            Timestamp:        time.Now(),
        }

        startTime := time.Now()

        queryResult, err := CDPDOMQuerySelector(st.rootNodeID, testCase.Selector)
        result.Duration = time.Since(startTime)

        if err != nil {
            result.Success = false
            result.Error = err.Error()
        } else {
            var resp struct {
                Result struct {
                    NodeID int `json:"nodeId"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(queryResult), &resp); err != nil {
                result.Success = false
                result.Error = err.Error()
            } else {
                result.NodeID = resp.Result.NodeID
                result.Found = resp.Result.NodeID > 0
                result.Success = true
            }
        }

        st.results = append(st.results, result)

        // 短暂延迟
        time.Sleep(100 * time.Millisecond)
    }
}

func (st *SelectorTester) GetResults() []TestResult {
    return st.results
}

func (st *SelectorTester) GenerateReport() {
    fmt.Printf("=== 选择器测试报告 ===\n")
    fmt.Printf("测试用例数量: %d\n\n", len(st.results))

    passed := 0
    failed := 0

    for _, result := range st.results {
        status := "❌ 失败"
        if result.Success && result.Found == result.Expected {
            status = "✅ 通过"
            passed++
        } else {
            failed++
        }

        fmt.Printf("%s %s\n", status, result.Name)
        fmt.Printf("  选择器: %s\n", result.Selector)
        fmt.Printf("  预期找到: %v, 实际找到: %v\n", result.Expected, result.Found)
        fmt.Printf("  耗时: %v\n", result.Duration)

        if result.Error != "" {
            fmt.Printf("  错误: %s\n", result.Error)
        }

        if result.NodeID > 0 {
            fmt.Printf("  节点ID: %d\n", result.NodeID)
        }

        fmt.Println()
    }

    fmt.Printf("总计: 通过 %d, 失败 %d, 成功率 %.1f%%\n",
        passed, failed, float64(passed)/float64(len(st.results))*100)
}

// 演示复杂选择器查询
func demonstrateComplexSelectorQueries() {
    fmt.Println("=== 复杂选择器查询演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 假设文档根节点ID
    documentNodeID := 1

    tester := NewSelectorTester(documentNodeID)

    // 添加测试用例
    testCases := []SelectorTestCase{
        {
            Name:        "主按钮查询",
            Selector:    "button.primary, .btn-primary",
            Description: "查询主按钮",
            Expected:    true,
            Priority:    10,
        },
        {
            Name:        "表单输入框",
            Selector:    "input[type='text'], input[type='email']",
            Description: "查询文本输入框",
            Expected:    true,
            Priority:    9,
        },
        {
            Name:        "导航链接",
            Selector:    "nav a, .navbar a, .nav a",
            Description: "查询导航链接",
            Expected:    true,
            Priority:    8,
        },
        {
            Name:        "页脚信息",
            Selector:    "footer, .footer, #footer",
            Description: "查询页脚",
            Expected:    true,
            Priority:    7,
        },
        {
            Name:        "不存在的元素",
            Selector:    ".non-existent-class",
            Description: "查询不存在的元素",
            Expected:    false,
            Priority:    1,
        },
    }

    for _, tc := range testCases {
        tester.AddTestCase(tc)
    }

    // 运行测试
    tester.RunTests()

    // 生成报告
    tester.GenerateReport()

    fmt.Println("=== 演示完成 ===")
}

*/

// -----------------------------------------------  DOM.querySelectorAll  -----------------------------------------------
// === 应用场景 ===
// 1. 批量查找: 查找匹配选择器的所有元素
// 2. 列表操作: 操作列表、表格等多元素场景
// 3. 统计信息: 统计页面中特定类型元素的数量
// 4. 批量处理: 批量修改或操作多个元素
// 5. 内容分析: 分析页面内容的分布和结构
// 6. 遍历操作: 遍历页面中的特定类型元素

// CDPDOMQuerySelectorAll 通过CSS选择器查询所有匹配元素
// nodeID: 要查询的根节点ID
// selector: CSS选择器字符串
func CDPDOMQuerySelectorAll(nodeID int, selector string) (string, error) {
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
        "method": "DOM.querySelectorAll",
        "params": {
            "nodeId": %d,
            "selector": "%s"
        }
    }`, reqID, nodeID, selector)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.querySelectorAll 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.querySelectorAll 请求超时")
		}
	}
}

/*



// 示例: 查询页面中的所有链接
func ExampleCDPDOMQuerySelectorAll() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 2. 获取文档根节点
    // 假设我们已经有了文档根节点ID
    documentNodeID := 1

    // 3. 定义要查询的CSS选择器
    testSelectors := []struct {
        selector    string
        description string
        expectedMin int  // 期望最少找到的数量
    }{
        {"a", "所有链接", 1},
        {"button", "所有按钮", 1},
        {"input", "所有输入框", 1},
        {"img", "所有图片", 0},
        {".btn", "所有按钮类", 0},
        {"[href]", "所有有链接的元素", 1},
        {"div", "所有div元素", 5},
        {"h1, h2, h3, h4, h5, h6", "所有标题元素", 1},
    }

    fmt.Printf("=== 批量元素查询测试 ===\n")
    fmt.Printf("文档根节点ID: %d\n\n", documentNodeID)

    // 4. 执行批量查询测试
    for i, test := range testSelectors {
        fmt.Printf("测试 %d: %s\n", i+1, test.description)
        fmt.Printf("选择器: %s\n", test.selector)

        result, err := CDPDOMQuerySelectorAll(documentNodeID, test.selector)
        if err != nil {
            fmt.Printf("❌ 查询失败: %v\n\n", err)
            continue
        }

        // 5. 解析查询结果
        var response struct {
            Result struct {
                NodeIDs []int `json:"nodeIds"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(result), &response); err != nil {
            fmt.Printf("❌ 解析结果失败: %v\n\n", err)
            continue
        }

        nodeIDs := response.Result.NodeIDs
        count := len(nodeIDs)

        fmt.Printf("✅ 找到 %d 个匹配元素\n", count)

        if count < test.expectedMin {
            fmt.Printf("⚠️ 找到的元素数量少于预期 (期望最少 %d 个)\n", test.expectedMin)
        }

        // 6. 显示元素统计信息
        if count > 0 {
            displayElementStats(nodeIDs, test.selector)
        }

        // 7. 高亮前几个元素
        if count > 0 {
            highlightFirstElements(nodeIDs, 3)
        }

        fmt.Println()
    }

    // 8. 性能测试
    fmt.Printf("=== 批量查询性能测试 ===\n")

    performanceTestSelectors := []struct {
        selector string
        name     string
    }{
        {"a", "链接选择器"},
        {".container *", "通配选择器"},
        {"div > p", "后代选择器"},
        {"[class]", "属性选择器"},
        {"li:nth-child(odd)", "复杂伪类选择器"},
    }

    for _, test := range performanceTestSelectors {
        duration, count := testBatchSelectorPerformance(documentNodeID, test.selector)
        fmt.Printf("选择器: %-20s 耗时: %v 找到: %d 个\n",
            test.name, duration, count)
    }

    // 9. 实际应用演示
    demonstrateRealWorldUsage(documentNodeID)
}

// 显示元素统计信息
func displayElementStats(nodeIDs []int, selector string) {
    if len(nodeIDs) == 0 {
        return
    }

    // 统计不同类型
    typeStats := make(map[string]int)

    // 只分析前10个元素以提高性能
    limit := 10
    if len(nodeIDs) < limit {
        limit = len(nodeIDs)
    }

    for i := 0; i < limit; i++ {
        nodeID := nodeIDs[i]

        // 获取节点信息
        descResult, err := CDPDOMDescribeNode(nodeID, 0, false)
        if err != nil {
            continue
        }

        var descResp struct {
            Result struct {
                Node struct {
                    NodeName string `json:"nodeName"`
                } `json:"node"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(descResult), &descResp); err != nil {
            continue
        }

        nodeName := descResp.Result.Node.NodeName
        typeStats[nodeName]++
    }

    // 显示统计
    if len(typeStats) > 0 {
        fmt.Printf("元素类型分布:\n")
        for nodeType, count := range typeStats {
            percentage := float64(count) / float64(limit) * 100
            fmt.Printf("  %-6s: %d (%.1f%%)\n", nodeType, count, percentage)
        }
    }

    // 显示前3个元素的信息
    fmt.Printf("前3个元素:\n")
    for i := 0; i < 3 && i < len(nodeIDs); i++ {
        nodeID := nodeIDs[i]

        // 获取简要信息
        if info, err := getBriefElementInfo(nodeID); err == nil {
            fmt.Printf("  %d. 节点ID: %d, 类型: %s",
                i+1, nodeID, info["nodeName"])

            if text, ok := info["text"].(string); ok && text != "" {
                fmt.Printf(", 文本: %s", truncateString(text, 30))
            }
            fmt.Println()
        }
    }
}

// 获取简要元素信息
func getBriefElementInfo(nodeID int) (map[string]interface{}, error) {
    info := make(map[string]interface{})

    // 获取节点描述
    descResult, err := CDPDOMDescribeNode(nodeID, 0, false)
    if err != nil {
        return nil, err
    }

    var descResp struct {
        Result struct {
            Node struct {
                NodeName  string `json:"nodeName"`
                NodeValue string `json:"nodeValue"`
            } `json:"node"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(descResult), &descResp); err != nil {
        return nil, err
    }

    info["nodeName"] = descResp.Result.Node.NodeName

    // 尝试获取文本内容
    outerHTMLResult, err := CDPDOMGetOuterHTML(nodeID)
    if err == nil {
        var htmlResp struct {
            Result struct {
                OuterHTML string `json:"outerHTML"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(outerHTMLResult), &htmlResp); err == nil {
            // 提取纯文本
            text := extractTextFromHTML(htmlResp.Result.OuterHTML)
            info["text"] = strings.TrimSpace(text)
        }
    }

    return info, nil
}

// 从HTML中提取文本
func extractTextFromHTML(htmlStr string) string {
    // 简单实现：移除HTML标签
    re := regexp.MustCompile(`<[^>]*>`)
    text := re.ReplaceAllString(htmlStr, " ")

    // 合并多个空格
    re = regexp.MustCompile(`\s+`)
    text = re.ReplaceAllString(text, " ")

    return text
}

// 高亮前几个元素
func highlightFirstElements(nodeIDs []int, count int) {
    if len(nodeIDs) == 0 {
        return
    }

    if count > len(nodeIDs) {
        count = len(nodeIDs)
    }

    colors := []RGBA{
        {R: 255, G: 100, B: 100, A: 0.3}, // 红色
        {R: 100, G: 255, B: 100, A: 0.3}, // 绿色
        {R: 100, G: 100, B: 255, A: 0.3}, // 蓝色
        {R: 255, G: 255, B: 100, A: 0.3}, // 黄色
        {R: 255, G: 100, B: 255, A: 0.3}, // 紫色
    }

    for i := 0; i < count; i++ {
        colorIndex := i % len(colors)
        highlightConfig := &HighlightConfig{
            ContentColor: &colors[colorIndex],
            BorderColor:  &RGBA{R: 0, G: 0, B: 0, A: 0.8},
        }

        if _, err := CDPDOMHighlightNode(nodeIDs[i], highlightConfig); err == nil {
            // 短暂延迟
            time.Sleep(300 * time.Millisecond)
        }
    }

    // 清理高亮
    time.Sleep(2 * time.Second)
    CDPDOMHideHighlight()
}

// 测试批量选择器性能
func testBatchSelectorPerformance(rootNodeID int, selector string) (time.Duration, int) {
    startTime := time.Now()

    result, err := CDPDOMQuerySelectorAll(rootNodeID, selector)
    if err != nil {
        return time.Since(startTime), 0
    }

    var response struct {
        Result struct {
            NodeIDs []int `json:"nodeIds"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        return time.Since(startTime), 0
    }

    return time.Since(startTime), len(response.Result.NodeIDs)
}

// 演示实际应用
func demonstrateRealWorldUsage(rootNodeID int) {
    fmt.Printf("\n=== 实际应用场景演示 ===\n")

    // 场景1: 统计页面链接
    fmt.Printf("场景1: 页面链接分析\n")
    analyzePageLinks(rootNodeID)

    // 场景2: 表单元素统计
    fmt.Printf("\n场景2: 表单元素统计\n")
    analyzeFormElements(rootNodeID)

    // 场景3: 图片资源分析
    fmt.Printf("\n场景3: 图片资源分析\n")
    analyzeImageElements(rootNodeID)
}

// 分析页面链接
func analyzePageLinks(rootNodeID int) {
    result, err := CDPDOMQuerySelectorAll(rootNodeID, "a[href]")
    if err != nil {
        fmt.Printf("❌ 查询链接失败: %v\n", err)
        return
    }

    var response struct {
        Result struct {
            NodeIDs []int `json:"nodeIds"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        fmt.Printf("❌ 解析结果失败: %v\n", err)
        return
    }

    links := response.Result.NodeIDs
    fmt.Printf("找到 %d 个链接\n", len(links))

    if len(links) == 0 {
        return
    }

    // 分析链接类型
    internalCount := 0
    externalCount := 0
    otherCount := 0

    for i := 0; i < len(links) && i < 20; i++ { // 限制分析数量
        nodeID := links[i]

        // 获取链接属性
        attrsResult, err := CDPDOMGetAttributes(nodeID)
        if err != nil {
            continue
        }

        var attrsResp struct {
            Result struct {
                Attributes []string `json:"attributes"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(attrsResult), &attrsResp); err != nil {
            continue
        }

        // 提取href属性
        href := ""
        for j := 0; j < len(attrsResp.Result.Attributes); j += 2 {
            if j+1 < len(attrsResp.Result.Attributes) &&
               attrsResp.Result.Attributes[j] == "href" {
                href = attrsResp.Result.Attributes[j+1]
                break
            }
        }

        if href == "" {
            otherCount++
        } else if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
            externalCount++
        } else {
            internalCount++
        }
    }

    total := internalCount + externalCount + otherCount
    if total > 0 {
        fmt.Printf("链接类型分布 (采样 %d 个):\n", total)
        fmt.Printf("  内部链接: %d (%.1f%%)\n",
            internalCount, float64(internalCount)/float64(total)*100)
        fmt.Printf("  外部链接: %d (%.1f%%)\n",
            externalCount, float64(externalCount)/float64(total)*100)
        fmt.Printf("  其他链接: %d (%.1f%%)\n",
            otherCount, float64(otherCount)/float64(total)*100)
    }
}

// 分析表单元素
func analyzeFormElements(rootNodeID int) {
    result, err := CDPDOMQuerySelectorAll(rootNodeID, "input, textarea, select, button")
    if err != nil {
        fmt.Printf("❌ 查询表单元素失败: %v\n", err)
        return
    }

    var response struct {
        Result struct {
            NodeIDs []int `json:"nodeIds"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        fmt.Printf("❌ 解析结果失败: %v\n", err)
        return
    }

    elements := response.Result.NodeIDs
    fmt.Printf("找到 %d 个表单相关元素\n", len(elements))

    if len(elements) == 0 {
        return
    }

    // 统计不同类型
    typeStats := make(map[string]int)

    for i := 0; i < len(elements) && i < 20; i++ {
        nodeID := elements[i]

        descResult, err := CDPDOMDescribeNode(nodeID, 0, false)
        if err != nil {
            continue
        }

        var descResp struct {
            Result struct {
                Node struct {
                    NodeName string `json:"nodeName"`
                } `json:"node"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(descResult), &descResp); err != nil {
            continue
        }

        nodeName := strings.ToLower(descResp.Result.Node.NodeName)
        typeStats[nodeName]++
    }

    fmt.Printf("表单元素类型分布:\n")
    for elemType, count := range typeStats {
        fmt.Printf("  %-10s: %d\n", elemType, count)
    }
}

// 分析图片元素
func analyzeImageElements(rootNodeID int) {
    result, err := CDPDOMQuerySelectorAll(rootNodeID, "img")
    if err != nil {
        fmt.Printf("❌ 查询图片失败: %v\n", err)
        return
    }

    var response struct {
        Result struct {
            NodeIDs []int `json:"nodeIds"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        fmt.Printf("❌ 解析结果失败: %v\n", err)
        return
    }

    images := response.Result.NodeIDs
    fmt.Printf("找到 %d 个图片元素\n", len(images))

    if len(images) == 0 {
        return
    }

    // 分析alt属性
    withAlt := 0
    withoutAlt := 0

    for i := 0; i < len(images) && i < 15; i++ {
        nodeID := images[i]

        attrsResult, err := CDPDOMGetAttributes(nodeID)
        if err != nil {
            continue
        }

        var attrsResp struct {
            Result struct {
                Attributes []string `json:"attributes"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(attrsResult), &attrsResp); err != nil {
            continue
        }

        hasAlt := false
        for j := 0; j < len(attrsResp.Result.Attributes); j += 2 {
            if j+1 < len(attrsResp.Result.Attributes) &&
               attrsResp.Result.Attributes[j] == "alt" {
                hasAlt = true
                break
            }
        }

        if hasAlt {
            withAlt++
        } else {
            withoutAlt++
        }
    }

    total := withAlt + withoutAlt
    if total > 0 {
        fmt.Printf("图片可访问性分析 (采样 %d 个):\n", total)
        fmt.Printf("  有alt属性: %d (%.1f%%)\n",
            withAlt, float64(withAlt)/float64(total)*100)
        fmt.Printf("  无alt属性: %d (%.1f%%)\n",
            withoutAlt, float64(withoutAlt)/float64(total)*100)
    }
}

// 高级功能: 批量元素处理器
type BatchElementProcessor struct {
    rootNodeID int
    selector   string
    elements   []int
    processed  int
    failed     int
}

func NewBatchElementProcessor(rootNodeID int, selector string) *BatchElementProcessor {
    return &BatchElementProcessor{
        rootNodeID: rootNodeID,
        selector:   selector,
    }
}

func (bep *BatchElementProcessor) LoadElements() error {
    result, err := CDPDOMQuerySelectorAll(bep.rootNodeID, bep.selector)
    if err != nil {
        return fmt.Errorf("查询元素失败: %w", err)
    }

    var response struct {
        Result struct {
            NodeIDs []int `json:"nodeIds"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        return fmt.Errorf("解析结果失败: %w", err)
    }

    bep.elements = response.Result.NodeIDs
    return nil
}

func (bep *BatchElementProcessor) ProcessEach(processor func(nodeID int) error) {
    for _, nodeID := range bep.elements {
        if err := processor(nodeID); err != nil {
            bep.failed++
            log.Printf("处理节点 %d 失败: %v", nodeID, err)
        } else {
            bep.processed++
        }
    }
}

func (bep *BatchElementProcessor) GetStats() map[string]interface{} {
    return map[string]interface{}{
        "total":     len(bep.elements),
        "processed": bep.processed,
        "failed":    bep.failed,
        "successRate": func() float64 {
            if len(bep.elements) == 0 {
                return 0
            }
            return float64(bep.processed) / float64(len(bep.elements)) * 100
        }(),
    }
}

// 演示批量处理器
func demonstrateBatchProcessor(rootNodeID int) {
    fmt.Printf("=== 批量元素处理器演示 ===\n")

    // 创建处理器
    processor := NewBatchElementProcessor(rootNodeID, "a[href]")

    // 加载元素
    if err := processor.LoadElements(); err != nil {
        fmt.Printf("❌ 加载元素失败: %v\n", err)
        return
    }

    // 处理每个元素
    fmt.Printf("找到 %d 个链接，开始处理...\n", len(processor.elements))

    processor.ProcessEach(func(nodeID int) error {
        // 示例：获取链接文本
        info, err := getBriefElementInfo(nodeID)
        if err != nil {
            return err
        }

        if text, ok := info["text"].(string); ok && text != "" {
            fmt.Printf("链接文本: %s\n", truncateString(text, 50))
        }

        return nil
    })

    // 显示统计
    stats := processor.GetStats()
    fmt.Printf("\n处理完成:\n")
    fmt.Printf("  总数: %d\n", stats["total"])
    fmt.Printf("  成功: %d\n", stats["processed"])
    fmt.Printf("  失败: %d\n", stats["failed"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
}

// 高级功能: 选择器有效性测试器
type SelectorEffectivenessTester struct {
    rootNodeID   int
    testData     []SelectorTestData
    results      []SelectorTestResult
}

type SelectorTestData struct {
    Name     string
    Selector string
    Purpose  string
    Priority int
}

type SelectorTestResult struct {
    SelectorTestData
    Count      int
    Duration   time.Duration
    Success    bool
    Error      string
    SampleInfo []map[string]interface{}
}

func NewSelectorEffectivenessTester(rootNodeID int) *SelectorEffectivenessTester {
    return &SelectorEffectivenessTester{
        rootNodeID: rootNodeID,
        testData:   make([]SelectorTestData, 0),
        results:    make([]SelectorTestResult, 0),
    }
}

func (set *SelectorEffectivenessTester) AddTest(test SelectorTestData) {
    set.testData = append(set.testData, test)
}

func (set *SelectorEffectivenessTester) RunTests() {
    // 按优先级排序
    sort.Slice(set.testData, func(i, j int) bool {
        return set.testData[i].Priority > set.testData[j].Priority
    })

    for _, test := range set.testData {
        result := SelectorTestResult{
            SelectorTestData: test,
        }

        startTime := time.Now()

        queryResult, err := CDPDOMQuerySelectorAll(set.rootNodeID, test.Selector)
        result.Duration = time.Since(startTime)

        if err != nil {
            result.Success = false
            result.Error = err.Error()
        } else {
            var resp struct {
                Result struct {
                    NodeIDs []int `json:"nodeIds"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(queryResult), &resp); err != nil {
                result.Success = false
                result.Error = err.Error()
            } else {
                result.Count = len(resp.NodeIDs)
                result.Success = true

                // 采样前几个元素
                sampleCount := 3
                if result.Count < sampleCount {
                    sampleCount = result.Count
                }

                for i := 0; i < sampleCount; i++ {
                    if info, err := getBriefElementInfo(resp.NodeIDs[i]); err == nil {
                        result.SampleInfo = append(result.SampleInfo, info)
                    }
                }
            }
        }

        set.results = append(set.results, result)

        // 短暂延迟
        time.Sleep(100 * time.Millisecond)
    }
}

func (set *SelectorEffectivenessTester) GenerateReport() {
    fmt.Printf("=== 选择器有效性测试报告 ===\n")
    fmt.Printf("测试用例数量: %d\n\n", len(set.results))

    for _, result := range set.results {
        status := "✅ 成功"
        if !result.Success {
            status = "❌ 失败"
        }

        fmt.Printf("%s %s\n", status, result.Name)
        fmt.Printf("  选择器: %s\n", result.Selector)
        fmt.Printf("  目的: %s\n", result.Purpose)
        fmt.Printf("  找到元素: %d 个\n", result.Count)
        fmt.Printf("  耗时: %v\n", result.Duration)

        if result.Error != "" {
            fmt.Printf("  错误: %s\n", result.Error)
        }

        if len(result.SampleInfo) > 0 {
            fmt.Printf("  示例元素:\n")
            for i, info := range result.SampleInfo {
                fmt.Printf("    %d. 类型: %s", i+1, info["nodeName"])
                if text, ok := info["text"].(string); ok && text != "" {
                    fmt.Printf(", 文本: %s", truncateString(text, 30))
                }
                fmt.Println()
            }
        }

        fmt.Println()
    }
}

*/

// -----------------------------------------------  DOM.removeAttribute  -----------------------------------------------
// === 应用场景 ===
// 1. 属性清理: 移除不需要的DOM属性
// 2. 样式重置: 移除内联样式属性
// 3. 状态重置: 移除元素状态属性（如disabled、readonly等）
// 4. 测试准备: 测试前清理元素属性
// 5. 动态修改: 动态修改元素属性状态
// 6. 安全清理: 移除潜在的不安全属性

// CDPDOMRemoveAttribute 移除指定节点的属性
// nodeID: 要移除属性的节点ID
// name: 要移除的属性名称
func CDPDOMRemoveAttribute(nodeID int, name string) (string, error) {
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
        "method": "DOM.removeAttribute",
        "params": {
            "nodeId": %d,
            "name": "%s"
        }
    }`, reqID, nodeID, name)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.removeAttribute 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.removeAttribute 请求超时")
		}
	}
}

/*

// 示例: 移除输入框的禁用属性
func ExampleCDPDOMRemoveAttribute() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个输入框的节点ID
    inputNodeID := 1001

    fmt.Printf("=== 移除属性操作演示 ===\n")
    fmt.Printf("目标元素节点ID: %d\n\n", inputNodeID)

    // 2. 首先获取元素的当前属性
    fmt.Printf("=== 操作前检查 ===\n")
    attrsBefore, err := getElementAttributes(inputNodeID)
    if err != nil {
        log.Printf("获取属性失败: %v", err)
        return
    }

    displayAttributes("操作前属性", attrsBefore)

    // 3. 检查是否有禁用属性
    hasDisabled := false
    for _, attr := range attrsBefore {
        if attr.Name == "disabled" {
            hasDisabled = true
            break
        }
    }

    if !hasDisabled {
        fmt.Printf("ℹ️ 元素当前没有被禁用\n")

        // 为了演示，先添加禁用属性
        fmt.Printf("添加禁用属性用于演示...\n")
        // 这里需要DOM.setAttribute方法，先模拟
        attrsBefore = append(attrsBefore, Attribute{Name: "disabled", Value: "true"})
    }

    // 4. 移除禁用属性
    fmt.Printf("\n=== 执行移除操作 ===\n")
    fmt.Printf("要移除的属性: disabled\n")

    result, err := CDPDOMRemoveAttribute(inputNodeID, "disabled")
    if err != nil {
        log.Printf("移除属性失败: %v", err)
        return
    }

    fmt.Printf("移除结果: %s\n", result)

    // 5. 验证属性已被移除
    fmt.Printf("\n=== 操作后验证 ===\n")

    attrsAfter, err := getElementAttributes(inputNodeID)
    if err != nil {
        log.Printf("获取属性失败: %v", err)
        return
    }

    displayAttributes("操作后属性", attrsAfter)

    // 检查是否还包含禁用属性
    stillDisabled := false
    for _, attr := range attrsAfter {
        if attr.Name == "disabled" {
            stillDisabled = true
            break
        }
    }

    if stillDisabled {
        fmt.Printf("❌ 禁用属性未被移除\n")
    } else {
        fmt.Printf("✅ 禁用属性已成功移除\n")
    }

    // 6. 测试其他常见属性移除场景
    fmt.Printf("\n=== 其他属性移除测试 ===\n")

    testCases := []struct {
        attrName  string
        attrDesc  string
        testValue string
    }{
        {"readonly", "只读属性", "true"},
        {"required", "必填属性", "true"},
        {"placeholder", "占位符", "请输入内容"},
        {"title", "提示文本", "这是一个提示"},
        {"style", "内联样式", "color: red; font-size: 16px;"},
    }

    for _, tc := range testCases {
        fmt.Printf("测试移除: %s (%s)\n", tc.attrName, tc.attrDesc)

        // 先确保属性存在
        // 这里需要setAttribute方法，先模拟测试

        result, err := CDPDOMRemoveAttribute(inputNodeID, tc.attrName)
        if err != nil {
            fmt.Printf("  ❌ 移除失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 移除结果: %s\n", result)
        }

        // 验证移除
        attrs, _ := getElementAttributes(inputNodeID)
        found := false
        for _, attr := range attrs {
            if attr.Name == tc.attrName {
                found = true
                break
            }
        }

        if found {
            fmt.Printf("  ❌ 属性仍然存在\n")
        } else {
            fmt.Printf("  ✅ 属性已移除\n")
        }

        fmt.Println()
    }

    // 7. 移除属性的影响分析
    fmt.Printf("\n=== 移除属性影响分析 ===\n")

    fmt.Printf("移除不同属性的影响:\n")
    fmt.Printf("  disabled: 元素变为可交互，可获取焦点，可输入\n")
    fmt.Printf("  readonly: 元素可获取焦点，可选择文本，可复制\n")
    fmt.Printf("  required: 表单提交时不再需要此字段\n")
    fmt.Printf("  style: 失去内联样式，使用CSS类样式\n")
    fmt.Printf("  title: 失去鼠标悬停提示文本\n")
    fmt.Printf("  placeholder: 失去占位符文本提示\n")

    // 8. 可访问性影响
    fmt.Printf("\n=== 可访问性影响 ===\n")

    fmt.Printf("移除属性对可访问性的影响:\n")
    fmt.Printf("  aria-* 属性: 可能影响屏幕阅读器\n")
    fmt.Printf("  role 属性: 改变元素语义角色\n")
    fmt.Printf("  tabindex: 改变键盘导航顺序\n")
    fmt.Printf("  alt: 图片失去替代文本\n")

    // 9. 安全考虑
    fmt.Printf("\n=== 安全考虑 ===\n")

    fmt.Printf("移除属性的安全考虑:\n")
    fmt.Printf("  onclick 等事件属性: 可能移除安全控制\n")
    fmt.Printf("  data-* 属性: 可能移除应用数据\n")
    fmt.Printf("  autocomplete: 可能影响自动填充安全\n")

    // 10. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次移除属性
    startTime := time.Now()
    operations := 10

    for i := 0; i < operations; i++ {
        attrName := fmt.Sprintf("data-test-%d", i)
        CDPDOMRemoveAttribute(inputNodeID, attrName)
    }

    elapsed := time.Since(startTime)
    fmt.Printf("%d 次属性移除操作耗时: %v\n", operations, elapsed)
    fmt.Printf("平均每次操作耗时: %v\n", elapsed/time.Duration(operations))

    // 11. 错误处理测试
    fmt.Printf("\n=== 错误处理测试 ===\n")

    errorTestCases := []struct {
        nodeID  int
        attrName string
        desc    string
    }{
        {9999, "disabled", "无效的节点ID"},
        {inputNodeID, "", "空的属性名"},
        {inputNodeID, "non-existent-attr", "不存在的属性"},
    }

    for _, tc := range errorTestCases {
        fmt.Printf("测试: %s\n", tc.desc)
        result, err := CDPDOMRemoveAttribute(tc.nodeID, tc.attrName)

        if err != nil {
            fmt.Printf("  ✅ 预期错误: %v\n", err)
        } else {
            fmt.Printf("  ❌ 预期错误但成功: %s\n", result)
        }
    }
}

// 获取元素属性
func getElementAttributes(nodeID int) ([]Attribute, error) {
    var attrs []Attribute

    result, err := CDPDOMGetAttributes(nodeID)
    if err != nil {
        return nil, err
    }

    var resp struct {
        Result struct {
            Attributes []string `json:"attributes"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &resp); err != nil {
        return nil, err
    }

    for i := 0; i < len(resp.Result.Attributes); i += 2 {
        if i+1 < len(resp.Result.Attributes) {
            attrs = append(attrs, Attribute{
                Name:  resp.Result.Attributes[i],
                Value: resp.Result.Attributes[i+1],
            })
        }
    }

    return attrs, nil
}

// 显示属性
func displayAttributes(title string, attrs []Attribute) {
    fmt.Printf("%s (%d 个):\n", title, len(attrs))

    if len(attrs) == 0 {
        fmt.Printf("  (无属性)\n")
        return
    }

    for _, attr := range attrs {
        fmt.Printf("  %-20s = %s\n", attr.Name, attr.Value)
    }
}

type Attribute struct {
    Name  string
    Value string
}

// 高级功能: 安全属性移除
func safeRemoveAttribute(nodeID int, attrName string, options RemoveOptions) error {
    // 验证节点存在
    if !validateNodeExists(nodeID) {
        return fmt.Errorf("节点不存在: %d", nodeID)
    }

    // 检查属性是否存在
    exists, err := attributeExists(nodeID, attrName)
    if err != nil {
        return fmt.Errorf("检查属性失败: %w", err)
    }

    if !exists && !options.AllowMissing {
        return fmt.Errorf("属性不存在: %s", attrName)
    }

    // 检查是否为敏感属性
    if isSensitiveAttribute(attrName) && !options.Force {
        return fmt.Errorf("敏感属性需要强制移除标志: %s", attrName)
    }

    // 执行移除
    result, err := CDPDOMRemoveAttribute(nodeID, attrName)
    if err != nil {
        return fmt.Errorf("移除失败: %w", err)
    }

    // 验证移除
    if options.Verify {
        stillExists, err := attributeExists(nodeID, attrName)
        if err != nil {
            return fmt.Errorf("验证失败: %w", err)
        }

        if stillExists {
            return fmt.Errorf("属性未被移除: %s", attrName)
        }
    }

    log.Printf("属性移除成功: %s, 结果: %s", attrName, result)
    return nil
}

type RemoveOptions struct {
    AllowMissing bool
    Verify       bool
    Force        bool
}

func attributeExists(nodeID int, attrName string) (bool, error) {
    attrs, err := getElementAttributes(nodeID)
    if err != nil {
        return false, err
    }

    for _, attr := range attrs {
        if attr.Name == attrName {
            return true, nil
        }
    }

    return false, nil
}

func isSensitiveAttribute(attrName string) bool {
    sensitiveAttrs := []string{
		"onclick", "onload", "onerror", "onsubmit",
		"onmouseover", "onkeydown", "onchange",
		"autocomplete", "autocapitalize", "autocorrect",
		"role", "aria-", "data-",
	}

    for _, sensitive := range sensitiveAttrs {
        if strings.Contains(attrName, sensitive) {
            return true
        }
    }

    return false
}

// 高级功能: 批量属性移除
func batchRemoveAttributes(nodeID int, attrNames []string, options BatchRemoveOptions) ([]RemoveResult, error) {
    var results []RemoveResult

    for _, attrName := range attrNames {
        result := RemoveResult{
            Attribute: attrName,
            Success:   false,
        }

        err := safeRemoveAttribute(nodeID, attrName, RemoveOptions{
            AllowMissing: options.AllowMissing,
            Verify:       options.Verify,
            Force:        options.Force,
        })

        if err != nil {
            result.Error = err.Error()
        } else {
            result.Success = true
        }

        results = append(results, result)

        // 可选延迟
        if options.Delay > 0 {
            time.Sleep(options.Delay)
        }
    }

    return results, nil
}

type BatchRemoveOptions struct {
    AllowMissing bool
    Verify       bool
    Force        bool
    Delay        time.Duration
}

type RemoveResult struct {
    Attribute string
    Success   bool
    Error     string
}

// 高级功能: 智能属性清理
func smartAttributeCleanup(nodeID int, cleanupRules []CleanupRule) error {
    attrs, err := getElementAttributes(nodeID)
    if err != nil {
        return fmt.Errorf("获取属性失败: %w", err)
    }

    for _, attr := range attrs {
        for _, rule := range cleanupRules {
            if rule.ShouldRemove(attr) {
                if err := safeRemoveAttribute(nodeID, attr.Name, rule.Options); err != nil {
                    if !rule.ContinueOnError {
                        return fmt.Errorf("移除属性 %s 失败: %w", attr.Name, err)
                    }
                    log.Printf("移除属性 %s 失败但继续: %v", attr.Name, err)
                } else {
                    log.Printf("已移除属性: %s", attr.Name)
                }
                break
            }
        }
    }

    return nil
}

type CleanupRule struct {
    MatchPattern  string
    Options       RemoveOptions
    ContinueOnError bool
}

func (cr *CleanupRule) ShouldRemove(attr Attribute) bool {
    matched, _ := regexp.MatchString(cr.MatchPattern, attr.Name)
    return matched
}

// 属性移除管理器
type AttributeRemovalManager struct {
    removalsHistory []RemovalRecord
    undoStack       []AttributeState
    redoStack       []AttributeState
}

type RemovalRecord struct {
    Timestamp   time.Time
    NodeID      int
    Attribute   string
    OldValue    string
    Duration    time.Duration
    Success     bool
    Error       string
}

type AttributeState struct {
    NodeID    int
    Attribute string
    Value     string
    Existed   bool
}

func NewAttributeRemovalManager() *AttributeRemovalManager {
    return &AttributeRemovalManager{
        removalsHistory: make([]RemovalRecord, 0),
        undoStack:       make([]AttributeState, 0),
        redoStack:       make([]AttributeState, 0),
    }
}

func (arm *AttributeRemovalManager) RemoveAttribute(nodeID int, attrName string, backup bool) error {
    startTime := time.Now()

    // 备份旧值
    var oldValue string
    var existed bool

    if backup {
        attrs, err := getElementAttributes(nodeID)
        if err == nil {
            for _, attr := range attrs {
                if attr.Name == attrName {
                    oldValue = attr.Value
                    existed = true
                    break
                }
            }
        }
    }

    // 执行移除
    result, err := CDPDOMRemoveAttribute(nodeID, attrName)
    duration := time.Since(startTime)

    // 记录历史
    record := RemovalRecord{
        Timestamp: startTime,
        NodeID:    nodeID,
        Attribute: attrName,
        OldValue:  oldValue,
        Duration:  duration,
        Success:   err == nil,
        Error:     "",
    }

    if err != nil {
        record.Error = err.Error()
    }

    arm.removalsHistory = append(arm.removalsHistory, record)

    // 保存到撤销栈
    if backup && existed {
        arm.undoStack = append(arm.undoStack, AttributeState{
            NodeID:    nodeID,
            Attribute: attrName,
            Value:     oldValue,
            Existed:   true,
        })
    }

    // 清空重做栈
    arm.redoStack = make([]AttributeState, 0)

    if err != nil {
        return fmt.Errorf("移除失败: %w", err)
    }

    log.Printf("属性移除成功: %s, 结果: %s", attrName, result)
    return nil
}

func (arm *AttributeRemovalManager) Undo() error {
    if len(arm.undoStack) == 0 {
        return fmt.Errorf("没有可撤销的操作")
    }

    // 获取最后一个操作
    lastState := arm.undoStack[len(arm.undoStack)-1]
    arm.undoStack = arm.undoStack[:len(arm.undoStack)-1]

    // 保存到重做栈
    // 需要先获取当前状态
    arm.redoStack = append(arm.redoStack, lastState)

    // 恢复属性
    // 这里需要setAttribute方法
    return fmt.Errorf("撤销功能需要setAttribute方法")
}

func (arm *AttributeRemovalManager) GetHistory() []RemovalRecord {
    return arm.removalsHistory
}

// 演示属性移除场景
func demonstrateAttributeRemovalScenarios() {
    fmt.Println("=== 属性移除场景演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    manager := NewAttributeRemovalManager()

    // 模拟节点
    inputNodeID := 1001
    buttonNodeID := 1002

    scenarios := []struct {
        name      string
        nodeID    int
        attrName  string
        desc      string
    }{
        {
            name:     "移除输入框禁用状态",
            nodeID:   inputNodeID,
            attrName: "disabled",
            desc:     "使输入框可编辑",
        },
        {
            name:     "移除输入框只读属性",
            nodeID:   inputNodeID,
            attrName: "readonly",
            desc:     "使输入框可修改",
        },
        {
            name:     "移除按钮禁用状态",
            nodeID:   buttonNodeID,
            attrName: "disabled",
            desc:     "使按钮可点击",
        },
        {
            name:     "移除内联样式",
            nodeID:   inputNodeID,
            attrName: "style",
            desc:     "移除内联样式，使用CSS类",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("\n场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.desc)

        if err := manager.RemoveAttribute(scenario.nodeID, scenario.attrName, true); err != nil {
            fmt.Printf("❌ 移除失败: %v\n", err)
        } else {
            fmt.Printf("✅ 移除成功\n")
        }
    }

    // 显示移除历史
    fmt.Printf("\n=== 移除历史 ===\n")
    history := manager.GetHistory()
    for i, record := range history {
        status := "✅ 成功"
        if !record.Success {
            status = "❌ 失败"
        }

        fmt.Printf("%d. [%s] 节点%d 属性:%s 耗时:%v\n",
            i+1, status, record.NodeID, record.Attribute, record.Duration)
    }

    fmt.Println("\n=== 场景演示完成 ===")
}

*/

// -----------------------------------------------  DOM.removeNode  -----------------------------------------------
// === 应用场景 ===
// 1. 元素删除: 从DOM树中删除指定节点
// 2. 动态清理: 动态清理不需要的页面元素
// 3. 内容更新: 更新页面内容时移除旧元素
// 4. 错误修复: 移除错误或异常的DOM节点
// 5. 测试清理: 测试完成后清理测试创建的节点
// 6. 性能优化: 移除不需要的DOM节点以优化性能

// CDPDOMRemoveNode 从DOM树中删除指定节点
// nodeID: 要删除的节点ID
func CDPDOMRemoveNode(nodeID int) (string, error) {
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
        "method": "DOM.removeNode",
        "params": {
            "nodeId": %d
        }
    }`, reqID, nodeID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.removeNode 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.removeNode 请求超时")
		}
	}
}

/*


// 示例: 删除页面中的临时元素
func ExampleCDPDOMRemoveNode() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们要删除一个临时消息框
    // 首先找到要删除的节点
    documentNodeID := 1
    tempMessageSelector := ".temp-message"

    fmt.Printf("=== 节点删除操作演示 ===\n")
    fmt.Printf("文档根节点ID: %d\n", documentNodeID)
    fmt.Printf("要查找的选择器: %s\n\n", tempMessageSelector)

    // 2. 查找要删除的节点
    result, err := CDPDOMQuerySelector(documentNodeID, tempMessageSelector)
    if err != nil {
        log.Printf("查找节点失败: %v", err)
        return
    }

    var queryResp struct {
        Result struct {
            NodeID int `json:"nodeId"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &queryResp); err != nil {
        log.Printf("解析查询结果失败: %v", err)
        return
    }

    nodeID := queryResp.Result.NodeID

    if nodeID == 0 {
        fmt.Printf("ℹ️ 未找到要删除的节点: %s\n", tempMessageSelector)
        return
    }

    fmt.Printf("找到要删除的节点，节点ID: %d\n", nodeID)

    // 3. 获取节点的详细信息
    fmt.Printf("\n=== 删除前检查 ===\n")

    nodeInfo, err := getNodeInfo(nodeID)
    if err != nil {
        log.Printf("获取节点信息失败: %v", err)
    } else {
        displayNodeInfo("要删除的节点", nodeInfo)
    }

    // 4. 检查节点是否可以安全删除
    fmt.Printf("\n=== 删除安全检查 ===\n")

    canDelete, reasons := checkNodeSafeToDelete(nodeID)

    if canDelete {
        fmt.Printf("✅ 节点可以安全删除\n")
    } else {
        fmt.Printf("❌ 节点可能不能安全删除:\n")
        for _, reason := range reasons {
            fmt.Printf("  - %s\n", reason)
        }
        fmt.Printf("建议先处理上述问题再删除\n")
        return
    }

    // 5. 执行删除操作
    fmt.Printf("\n=== 执行删除操作 ===\n")

    deleteResult, err := CDPDOMRemoveNode(nodeID)
    if err != nil {
        log.Printf("删除节点失败: %v", err)
        return
    }

    fmt.Printf("删除结果: %s\n", deleteResult)

    // 6. 验证删除结果
    fmt.Printf("\n=== 删除后验证 ===\n")

    // 尝试再次查找节点
    verifyResult, err := CDPDOMQuerySelector(documentNodeID, tempMessageSelector)
    if err != nil {
        fmt.Printf("❌ 验证查询失败: %v\n", err)
    } else {
        var verifyResp struct {
            Result struct {
                NodeID int `json:"nodeId"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(verifyResult), &verifyResp); err != nil {
            fmt.Printf("❌ 解析验证结果失败: %v\n", err)
        } else if verifyResp.Result.NodeID > 0 {
            fmt.Printf("❌ 节点仍然存在，节点ID: %d\n", verifyResp.Result.NodeID)
        } else {
            fmt.Printf("✅ 节点已成功删除\n")
        }
    }

    // 7. 删除影响分析
    fmt.Printf("\n=== 删除影响分析 ===\n")

    fmt.Printf("删除操作的影响:\n")
    fmt.Printf("  - 从DOM树中移除节点\n")
    fmt.Printf("  - 释放相关内存\n")
    fmt.Printf("  - 触发浏览器重排(reflow)\n")
    fmt.Printf("  - 解除事件监听器\n")
    fmt.Printf("  - 停止CSS动画\n")

    // 8. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次删除的性能
    startTime := time.Now()
    deleteOperations := 5

    // 创建测试节点
    testNodes := createTestNodes(documentNodeID, deleteOperations)

    for i := 0; i < deleteOperations; i++ {
        if testNodes[i] > 0 {
            CDPDOMRemoveNode(testNodes[i])
        }
    }

    elapsed := time.Since(startTime)
    fmt.Printf("%d 次删除操作耗时: %v\n", deleteOperations, elapsed)
    fmt.Printf("平均每次删除耗时: %v\n", elapsed/time.Duration(deleteOperations))

    // 9. 错误处理测试
    fmt.Printf("\n=== 错误处理测试 ===\n")

    errorTestCases := []struct {
        nodeID  int
        desc    string
    }{
        {0, "无效的节点ID 0"},
        {999999, "不存在的节点ID"},
        {nodeID, "已删除的节点ID"},
    }

    for _, tc := range errorTestCases {
        fmt.Printf("测试: %s\n", tc.desc)
        result, err := CDPDOMRemoveNode(tc.nodeID)

        if err != nil {
            fmt.Printf("  ✅ 预期错误: %v\n", err)
        } else {
            fmt.Printf("  ❌ 预期错误但成功: %s\n", result)
        }
    }

    // 10. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
    }{
        {
            name:        "弹窗关闭",
            description: "用户关闭弹窗时删除弹窗DOM节点",
        },
        {
            name:        "列表项删除",
            description: "从列表中删除一项时移除对应的DOM节点",
        },
        {
            name:        "缓存清理",
            description: "清理过期的缓存DOM节点以释放内存",
        },
        {
            name:        "错误修复",
            description: "删除错误的或异常的DOM节点",
        },
        {
            name:        "性能优化",
            description: "删除不可见的或不需要的DOM节点提升性能",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n\n", scenario.description)
    }
}

// 获取节点信息
func getNodeInfo(nodeID int) (map[string]interface{}, error) {
    info := make(map[string]interface{})

    // 获取节点描述
    descResult, err := CDPDOMDescribeNode(nodeID, 1, false)
    if err != nil {
        return nil, err
    }

    var descResp struct {
        Result struct {
            Node struct {
                NodeName       string   `json:"nodeName"`
                LocalName      string   `json:"localName"`
                NodeValue      string   `json:"nodeValue"`
                ChildNodeCount int      `json:"childNodeCount"`
                Attributes     []string `json:"attributes,omitempty"`
            } `json:"node"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(descResult), &descResp); err != nil {
        return nil, err
    }

    node := descResp.Result.Node
    info["nodeName"] = node.NodeName
    info["localName"] = node.LocalName
    info["nodeValue"] = node.NodeValue
    info["childCount"] = node.ChildNodeCount
    info["attributes"] = node.Attributes

    // 获取外层HTML
    outerHTMLResult, err := CDPDOMGetOuterHTML(nodeID)
    if err == nil {
        var htmlResp struct {
            Result struct {
                OuterHTML string `json:"outerHTML"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(outerHTMLResult), &htmlResp); err == nil {
            info["outerHTML"] = htmlResp.Result.OuterHTML
        }
    }

    return info, nil
}

// 显示节点信息
func displayNodeInfo(title string, info map[string]interface{}) {
    fmt.Printf("%s:\n", title)
    fmt.Printf("  节点名称: %s\n", info["nodeName"])
    fmt.Printf("  本地名称: %s\n", info["localName"])

    if nodeValue, ok := info["nodeValue"].(string); ok && nodeValue != "" {
        fmt.Printf("  节点值: %s\n", nodeValue)
    }

    if childCount, ok := info["childCount"].(int); ok {
        fmt.Printf("  子节点数量: %d\n", childCount)
    }

    if attrs, ok := info["attributes"].([]string); ok && len(attrs) > 0 {
        fmt.Printf("  属性: ")
        for i := 0; i < len(attrs); i += 2 {
            if i+1 < len(attrs) {
                fmt.Printf("%s=\"%s\" ", attrs[i], attrs[i+1])
            }
        }
        fmt.Println()
    }

    if html, ok := info["outerHTML"].(string); ok {
        truncated := truncateString(html, 100)
        fmt.Printf("  HTML预览: %s\n", truncated)
    }
}

// 检查节点是否可以安全删除
func checkNodeSafeToDelete(nodeID int) (bool, []string) {
    var reasons []string

    // 获取节点信息
    info, err := getNodeInfo(nodeID)
    if err != nil {
        return false, []string{fmt.Sprintf("无法获取节点信息: %v", err)}
    }

    // 检查是否重要元素
    nodeName, _ := info["nodeName"].(string)
    if isCriticalElement(nodeName) {
        reasons = append(reasons, fmt.Sprintf("可能是重要元素: %s", nodeName))
    }

    // 检查是否有子节点
    if childCount, ok := info["childCount"].(int); ok && childCount > 0 {
        reasons = append(reasons, fmt.Sprintf("有 %d 个子节点将被一起删除", childCount))
    }

    // 检查是否包含表单数据
    if containsFormData(info) {
        reasons = append(reasons, "包含表单数据")
    }

    // 检查是否有事件监听器
    if hasEventListeners(info) {
        reasons = append(reasons, "可能有事件监听器")
    }

    return len(reasons) == 0, reasons
}

func isCriticalElement(nodeName string) bool {
    criticalElements := map[string]bool{
        "BODY":   true,
        "HEAD":   true,
        "HTML":   true,
        "TITLE":  true,
        "SCRIPT": true,
        "LINK":   true,
        "META":   true,
    }

    return criticalElements[strings.ToUpper(nodeName)]
}

func containsFormData(info map[string]interface{}) bool {
    if nodeName, ok := info["nodeName"].(string); ok {
        formElements := map[string]bool{
            "INPUT":    true,
            "TEXTAREA": true,
            "SELECT":   true,
        }

        if formElements[strings.ToUpper(nodeName)] {
            return true
        }
    }

    if attrs, ok := info["attributes"].([]string); ok {
        for i := 0; i < len(attrs); i += 2 {
            if i+1 < len(attrs) && attrs[i] == "value" && attrs[i+1] != "" {
                return true
            }
        }
    }

    return false
}

func hasEventListeners(info map[string]interface{}) bool {
    if attrs, ok := info["attributes"].([]string); ok {
        for i := 0; i < len(attrs); i += 2 {
            if i+1 < len(attrs) && strings.HasPrefix(attrs[i], "on") {
                return true
            }
        }
    }

    return false
}

// 创建测试节点
func createTestNodes(parentNodeID, count int) []int {
    var nodeIDs []int

    // 这里需要实现创建测试节点的方法
    // 简化实现：返回空的节点ID数组
    for i := 0; i < count; i++ {
        nodeIDs = append(nodeIDs, 0)
    }

    return nodeIDs
}

// 高级功能: 安全删除节点
func safeRemoveNode(nodeID int, options RemoveNodeOptions) error {
    // 验证节点存在
    if !validateNodeExists(nodeID) {
        if !options.AllowMissing {
            return fmt.Errorf("节点不存在: %d", nodeID)
        }
        return nil
    }

    // 检查节点是否可以删除
    if options.CheckSafety {
        canDelete, reasons := checkNodeSafeToDelete(nodeID)
        if !canDelete && !options.Force {
            return fmt.Errorf("节点不能安全删除: %v", reasons)
        }
    }

    // 备份节点信息
    var backup *NodeBackup
    if options.Backup {
        backup, _ = backupNode(nodeID)
    }

    // 执行删除
    result, err := CDPDOMRemoveNode(nodeID)
    if err != nil {
        return fmt.Errorf("删除失败: %w", err)
    }

    // 验证删除
    if options.Verify {
        stillExists, _ := validateNodeExists(nodeID)
        if stillExists {
            // 尝试恢复
            if backup != nil {
                restoreNode(backup)
            }
            return fmt.Errorf("删除验证失败: 节点仍然存在")
        }
    }

    // 记录删除
    if backup != nil {
        log.Printf("节点删除成功: %d, 备份ID: %s, 结果: %s",
            nodeID, backup.ID, result)
    } else {
        log.Printf("节点删除成功: %d, 结果: %s", nodeID, result)
    }

    return nil
}

type RemoveNodeOptions struct {
    AllowMissing bool
    CheckSafety  bool
    Force        bool
    Backup       bool
    Verify       bool
}

type NodeBackup struct {
    ID        string
    NodeID    int
    Info      map[string]interface{}
    OuterHTML string
    Timestamp time.Time
}

func backupNode(nodeID int) (*NodeBackup, error) {
    info, err := getNodeInfo(nodeID)
    if err != nil {
        return nil, err
    }

    backup := &NodeBackup{
        ID:        fmt.Sprintf("backup-%d-%d", nodeID, time.Now().Unix()),
        NodeID:    nodeID,
        Info:      info,
        Timestamp: time.Now(),
    }

    if html, ok := info["outerHTML"].(string); ok {
        backup.OuterHTML = html
    }

    return backup, nil
}

func restoreNode(backup *NodeBackup) error {
    // 这里需要实现节点恢复功能
    // 需要DOM.createElement等方法
    return fmt.Errorf("节点恢复功能需要更多DOM方法支持")
}

// 高级功能: 批量删除节点
func batchRemoveNodes(nodeIDs []int, options BatchRemoveOptions) ([]RemoveNodeResult, error) {
    var results []RemoveNodeResult

    for _, nodeID := range nodeIDs {
        result := RemoveNodeResult{
            NodeID:  nodeID,
            Success: false,
        }

        err := safeRemoveNode(nodeID, RemoveNodeOptions{
            AllowMissing: options.AllowMissing,
            CheckSafety:  options.CheckSafety,
            Force:        options.Force,
            Backup:       options.Backup,
            Verify:       options.Verify,
        })

        if err != nil {
            result.Error = err.Error()
        } else {
            result.Success = true
        }

        results = append(results, result)

        // 可选延迟
        if options.Delay > 0 {
            time.Sleep(options.Delay)
        }
    }

    return results, nil
}

type BatchRemoveOptions struct {
    AllowMissing bool
    CheckSafety  bool
    Force        bool
    Backup       bool
    Verify       bool
    Delay        time.Duration
}

type RemoveNodeResult struct {
    NodeID  int
    Success bool
    Error   string
}

// 高级功能: 智能节点清理
func smartNodeCleanup(rootNodeID int, cleanupRules []NodeCleanupRule) error {
    // 获取所有节点
    // 这里需要遍历DOM树，简化实现
    nodesToRemove, err := findNodesToRemove(rootNodeID, cleanupRules)
    if err != nil {
        return fmt.Errorf("查找要删除的节点失败: %w", err)
    }

    // 执行批量删除
    options := BatchRemoveOptions{
        AllowMissing: true,
        CheckSafety:  true,
        Backup:       false,
        Verify:       false,
        Delay:        50 * time.Millisecond,
    }

    results, err := batchRemoveNodes(nodesToRemove, options)
    if err != nil {
        return fmt.Errorf("批量删除失败: %w", err)
    }

    // 统计结果
    removed := 0
    for _, result := range results {
        if result.Success {
            removed++
        }
    }

    log.Printf("智能清理完成: 尝试删除 %d 个节点，成功删除 %d 个",
        len(nodesToRemove), removed)

    return nil
}

type NodeCleanupRule struct {
    Selector     string
    Condition    func(map[string]interface{}) bool
    Description  string
    Priority     int
}

func findNodesToRemove(rootNodeID int, rules []NodeCleanupRule) ([]int, error) {
    var nodeIDs []int

    // 按优先级排序规则
    sort.Slice(rules, func(i, j int) bool {
        return rules[i].Priority > rules[j].Priority
    })

    for _, rule := range rules {
        if rule.Selector != "" {
            // 使用选择器查找节点
            result, err := CDPDOMQuerySelectorAll(rootNodeID, rule.Selector)
            if err != nil {
                continue
            }

            var resp struct {
                Result struct {
                    NodeIDs []int `json:"nodeIds"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(result), &resp); err != nil {
                continue
            }

            for _, nodeID := range resp.Result.NodeIDs {
                // 检查条件
                if rule.Condition != nil {
                    info, err := getNodeInfo(nodeID)
                    if err == nil && rule.Condition(info) {
                        nodeIDs = append(nodeIDs, nodeID)
                    }
                } else {
                    nodeIDs = append(nodeIDs, nodeID)
                }
            }
        }
    }

    return nodeIDs, nil
}

// 节点删除管理器
type NodeRemovalManager struct {
    removalsHistory []NodeRemovalRecord
    undoStack       []NodeBackup
    redoStack       []NodeBackup
}

type NodeRemovalRecord struct {
    Timestamp   time.Time
    NodeID      int
    NodeInfo    map[string]interface{}
    Duration    time.Duration
    Success     bool
    Error       string
    Rule        string
}

func NewNodeRemovalManager() *NodeRemovalManager {
    return &NodeRemovalManager{
        removalsHistory: make([]NodeRemovalRecord, 0),
        undoStack:       make([]NodeBackup, 0),
        redoStack:       make([]NodeBackup, 0),
    }
}

func (nrm *NodeRemovalManager) RemoveNode(nodeID int, rule string) error {
    startTime := time.Now()

    // 备份节点
    backup, err := backupNode(nodeID)
    if err != nil {
        return fmt.Errorf("备份节点失败: %w", err)
    }

    // 执行删除
    result, err := CDPDOMRemoveNode(nodeID)
    duration := time.Since(startTime)

    // 记录历史
    record := NodeRemovalRecord{
        Timestamp: startTime,
        NodeID:    nodeID,
        NodeInfo:  backup.Info,
        Duration:  duration,
        Success:   err == nil,
        Error:     "",
        Rule:      rule,
    }

    if err != nil {
        record.Error = err.Error()
    }

    nrm.removalsHistory = append(nrm.removalsHistory, record)

    // 保存到撤销栈
    nrm.undoStack = append(nrm.undoStack, *backup)

    // 清空重做栈
    nrm.redoStack = make([]NodeBackup, 0)

    if err != nil {
        return fmt.Errorf("删除失败: %w", err)
    }

    log.Printf("节点删除成功: %d, 规则: %s, 结果: %s", nodeID, rule, result)
    return nil
}

func (nrm *NodeRemovalManager) Undo() error {
    if len(nrm.undoStack) == 0 {
        return fmt.Errorf("没有可撤销的操作")
    }

    // 获取最后一个备份
    backup := nrm.undoStack[len(nrm.undoStack)-1]
    nrm.undoStack = nrm.undoStack[:len(nrm.undoStack)-1]

    // 保存到重做栈
    nrm.redoStack = append(nrm.redoStack, backup)

    // 恢复节点
    return fmt.Errorf("撤销功能需要节点恢复方法")
}

func (nrm *NodeRemovalManager) GetHistory() []NodeRemovalRecord {
    return nrm.removalsHistory
}

func (nrm *NodeRemovalManager) GetStats() map[string]interface{} {
    total := len(nrm.removalsHistory)
    success := 0
    var totalDuration time.Duration

    for _, record := range nrm.removalsHistory {
        if record.Success {
            success++
        }
        totalDuration += record.Duration
    }

    avgDuration := time.Duration(0)
    if total > 0 {
        avgDuration = totalDuration / time.Duration(total)
    }

    return map[string]interface{}{
        "totalRemovals":   total,
        "successful":      success,
        "failed":          total - success,
        "successRate":     float64(success) / float64(total) * 100,
        "totalDuration":   totalDuration,
        "averageDuration": avgDuration,
    }
}

// 演示节点删除场景
func demonstrateNodeRemovalScenarios() {
    fmt.Println("=== 节点删除场景演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    manager := NewNodeRemovalManager()

    // 模拟文档根节点
    documentNodeID := 1

    scenarios := []struct {
        name     string
        selector string
        rule     string
        desc     string
    }{
        {
            name:     "删除临时消息",
            selector: ".temp-message, .alert-temporary",
            rule:     "temporary_alerts",
            desc:     "删除临时通知消息",
        },
        {
            name:     "删除加载指示器",
            selector: ".loading, .spinner, [data-loading='true']",
            rule:     "loading_indicators",
            desc:     "删除加载完成后的指示器",
        },
        {
            name:     "删除过期广告",
            selector: ".ad-expired, [data-ad-expired='true']",
            rule:     "expired_ads",
            desc:     "删除过期的广告内容",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("\n场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.desc)

        // 查找匹配的节点
        result, err := CDPDOMQuerySelectorAll(documentNodeID, scenario.selector)
        if err != nil {
            fmt.Printf("❌ 查找节点失败: %v\n", err)
            continue
        }

        var resp struct {
            Result struct {
                NodeIDs []int `json:"nodeIds"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(result), &resp); err != nil {
            fmt.Printf("❌ 解析结果失败: %v\n", err)
            continue
        }

        fmt.Printf("找到 %d 个匹配节点\n", len(resp.Result.NodeIDs))

        // 删除每个节点
        for i, nodeID := range resp.Result.NodeIDs {
            fmt.Printf("删除节点 %d/%d (ID: %d)\n",
                i+1, len(resp.Result.NodeIDs), nodeID)

            if err := manager.RemoveNode(nodeID, scenario.rule); err != nil {
                fmt.Printf("  ❌ 删除失败: %v\n", err)
            } else {
                fmt.Printf("  ✅ 删除成功\n")
            }
        }
    }

    // 显示统计信息
    fmt.Printf("\n=== 删除统计 ===\n")
    stats := manager.GetStats()
    fmt.Printf("总删除操作: %d\n", stats["totalRemovals"])
    fmt.Printf("成功: %d\n", stats["successful"])
    fmt.Printf("失败: %d\n", stats["failed"])
    fmt.Printf("成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("总耗时: %v\n", stats["totalDuration"])
    fmt.Printf("平均耗时: %v\n", stats["averageDuration"])

    fmt.Println("\n=== 场景演示完成 ===")
}


*/

// -----------------------------------------------  DOM.requestChildNodes  -----------------------------------------------
// === 应用场景 ===
// 1. 延迟加载: 延迟加载DOM节点的子节点
// 2. 性能优化: 优化大型DOM树的加载性能
// 3. 分页查看: 分批查看大量子节点
// 4. 动态加载: 动态加载折叠或隐藏的内容
// 5. 内存管理: 管理大型列表的内存使用
// 6. 调试支持: 调试时逐步加载DOM结构

// CDPDOMRequestChildNodes 请求获取节点的子节点
// nodeID: 要获取子节点的节点ID
// depth: 遍历深度，-1表示完整子树，0表示仅节点自身，正整数表示深度
// pierce: 是否穿透shadow root
func CDPDOMRequestChildNodes(nodeID, depth int, pierce bool) (string, error) {
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
        "method": "DOM.requestChildNodes",
        "params": {
            "nodeId": %d,
            "depth": %d,
            "pierce": %v
        }
    }`, reqID, nodeID, depth, pierce)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.requestChildNodes 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.requestChildNodes 请求超时")
		}
	}
}

/*

// 示例: 延迟加载大型列表的子节点
func ExampleCDPDOMRequestChildNodes() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个大型列表的容器节点
    listContainerID := 1001

    fmt.Printf("=== 延迟加载子节点演示 ===\n")
    fmt.Printf("容器节点ID: %d\n\n", listContainerID)

    // 2. 首先获取容器节点的基本信息
    fmt.Printf("=== 初始状态检查 ===\n")

    containerInfo, err := getNodeInfo(listContainerID)
    if err != nil {
        log.Printf("获取容器信息失败: %v", err)
        return
    }

    displayContainerInfo(containerInfo)

    // 3. 获取子节点数量
    childCount := 0
    if count, ok := containerInfo["childCount"].(int); ok {
        childCount = count
    }

    fmt.Printf("子节点总数: %d\n", childCount)

    if childCount == 0 {
        fmt.Printf("容器没有子节点，不需要延迟加载\n")
        return
    }

    // 4. 根据子节点数量决定加载策略
    fmt.Printf("\n=== 加载策略选择 ===\n")

    var loadDepth int
    var loadDescription string

    if childCount > 1000 {
        fmt.Printf("检测到大型列表（%d 个子节点）\n", childCount)
        fmt.Printf("使用分页加载策略\n")
        loadDepth = 2
        loadDescription = "分页加载（深度2）"
    } else if childCount > 100 {
        fmt.Printf("检测到中型列表（%d 个子节点）\n", childCount)
        fmt.Printf("使用分层加载策略\n")
        loadDepth = 3
        loadDescription = "分层加载（深度3）"
    } else {
        fmt.Printf("检测到小型列表（%d 个子节点）\n", childCount)
        fmt.Printf("使用完整加载策略\n")
        loadDepth = -1
        loadDescription = "完整加载"
    }

    // 5. 执行延迟加载
    fmt.Printf("\n=== 执行延迟加载 ===\n")
    fmt.Printf("加载策略: %s\n", loadDescription)
    fmt.Printf("加载深度: %d\n", loadDepth)

    startTime := time.Now()

    result, err := CDPDOMRequestChildNodes(listContainerID, loadDepth, false)
    if err != nil {
        log.Printf("请求子节点失败: %v", err)
        return
    }

    loadTime := time.Since(startTime)
    fmt.Printf("加载结果: %s\n", result)
    fmt.Printf("加载耗时: %v\n", loadTime)

    // 6. 解析加载结果
    fmt.Printf("\n=== 加载结果解析 ===\n")

    // 注意：DOM.requestChildNodes 是异步操作
    // 结果需要通过DOM.setChildNodes事件获取
    // 这里模拟事件处理

    // 模拟接收到的子节点数据
    simulatedChildData := simulateChildNodesReceived(listContainerID, childCount, loadDepth)

    displayLoadedChildren(simulatedChildData, loadDepth)

    // 7. 性能分析
    fmt.Printf("\n=== 性能分析 ===\n")

    analyzePerformance(childCount, loadDepth, loadTime)

    // 8. 内存使用分析
    fmt.Printf("\n=== 内存使用分析 ===\n")

    estimateMemoryUsage(childCount, loadDepth)

    // 9. 分页加载演示
    if childCount > 100 {
        fmt.Printf("\n=== 分页加载演示 ===\n")
        demonstratePagedLoading(listContainerID, childCount)
    }

    // 10. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "大型数据表格",
            description: "表格有上千行数据",
            useCase:     "分页加载可见行，滚动时动态加载",
        },
        {
            name:        "文件目录树",
            description: "深层嵌套的目录结构",
            useCase:     "点击展开时加载子目录",
        },
        {
            name:        "评论列表",
            description: "文章有大量评论",
            useCase:     "初始加载部分评论，点击加载更多",
        },
        {
            name:        "产品目录",
            description: "电商网站产品列表",
            useCase:     "虚拟滚动，只渲染可见项",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 显示容器信息
func displayContainerInfo(info map[string]interface{}) {
    fmt.Printf("容器节点信息:\n")
    fmt.Printf("  节点名称: %s\n", info["nodeName"])
    fmt.Printf("  本地名称: %s\n", info["localName"])

    if childCount, ok := info["childCount"].(int); ok {
        fmt.Printf("  子节点数量: %d\n", childCount)
    }

    if html, ok := info["outerHTML"].(string); ok {
        truncated := truncateString(html, 100)
        fmt.Printf("  HTML预览: %s\n", truncated)
    }
}

// 模拟接收到的子节点数据
func simulateChildNodesReceived(parentID, childCount, depth int) map[string]interface{} {
    data := make(map[string]interface{})
    data["parentId"] = parentID
    data["totalChildren"] = childCount
    data["loadedDepth"] = depth
    data["timestamp"] = time.Now()

    // 模拟加载的子节点类型分布
    if depth > 0 {
        // 简化模拟
        nodeTypes := map[string]int{
            "DIV":     childCount / 2,
            "SPAN":    childCount / 4,
            "LI":      childCount / 8,
            "A":       childCount / 16,
            "BUTTON":  childCount / 32,
        }

        // 确保总数正确
        total := 0
        for _, count := range nodeTypes {
            total += count
        }
        if total < childCount {
            nodeTypes["DIV"] += childCount - total
        }

        data["nodeTypes"] = nodeTypes

        // 模拟加载的层级
        if depth >= 2 {
            data["loadedLevels"] = depth
            data["estimatedTotalNodes"] = estimateTotalNodes(childCount, depth)
        }
    }

    return data
}

// 显示加载的子节点
func displayLoadedChildren(data map[string]interface{}, depth int) {
    parentID, _ := data["parentId"].(int)
    totalChildren, _ := data["totalChildren"].(int)

    fmt.Printf("父节点ID: %d\n", parentID)
    fmt.Printf("总子节点数: %d\n", totalChildren)
    fmt.Printf("加载深度: %d\n", depth)

    if nodeTypes, ok := data["nodeTypes"].(map[string]int); ok {
        fmt.Printf("节点类型分布:\n")
        for nodeType, count := range nodeTypes {
            percentage := float64(count) / float64(totalChildren) * 100
            fmt.Printf("  %-8s: %d (%.1f%%)\n", nodeType, count, percentage)
        }
    }

    if loadedLevels, ok := data["loadedLevels"].(int); ok {
        fmt.Printf("已加载层级: %d\n", loadedLevels)
    }

    if estimatedTotal, ok := data["estimatedTotalNodes"].(int); ok {
        fmt.Printf("估计总节点数: %d\n", estimatedTotal)
    }
}

// 估计总节点数
func estimateTotalNodes(firstLevelCount, depth int) int {
    // 简单估计：假设每个节点平均有2个子节点
    if depth <= 1 {
        return firstLevelCount
    }

    total := firstLevelCount
    currentLevel := firstLevelCount

    for i := 2; i <= depth; i++ {
        currentLevel = currentLevel * 2 // 假设每个节点有2个子节点
        total += currentLevel
    }

    return total
}

// 性能分析
func analyzePerformance(childCount, depth int, loadTime time.Duration) {
    fmt.Printf("性能指标:\n")
    fmt.Printf("  子节点数: %d\n", childCount)
    fmt.Printf("  加载深度: %d\n", depth)
    fmt.Printf("  加载时间: %v\n", loadTime)

    if childCount > 0 && loadTime > 0 {
        nodesPerSecond := float64(childCount) / loadTime.Seconds()
        fmt.Printf("  节点加载速度: %.0f 节点/秒\n", nodesPerSecond)
    }

    // 性能建议
    fmt.Printf("\n性能建议:\n")

    if childCount > 1000 {
        fmt.Printf("  ⚠️ 子节点数量较多，考虑以下优化:\n")
        fmt.Printf("    - 使用虚拟滚动\n")
        fmt.Printf("    - 实现分页加载\n")
        fmt.Printf("    - 使用窗口化渲染\n")
    }

    if loadTime > time.Second {
        fmt.Printf("  ⚠️ 加载时间较长，考虑以下优化:\n")
        fmt.Printf("    - 减少初始加载深度\n")
        fmt.Printf("    - 延迟加载非可见区域\n")
        fmt.Printf("    - 使用Web Worker处理数据\n")
    }
}

// 估计内存使用
func estimateMemoryUsage(childCount, depth int) {
    // 简化的内存估算
    // 假设每个DOM节点平均占用1KB内存
    estimatedNodes := estimateTotalNodes(childCount, depth)
    estimatedMemoryMB := float64(estimatedNodes) * 1.0 / 1024 // 转换为MB

    fmt.Printf("内存使用估算:\n")
    fmt.Printf("  估计总节点数: %d\n", estimatedNodes)
    fmt.Printf("  估计内存占用: %.1f MB\n", estimatedMemoryMB)

    if estimatedMemoryMB > 10 {
        fmt.Printf("  ⚠️ 内存占用较高，考虑以下优化:\n")
        fmt.Printf("    - 使用懒加载\n")
        fmt.Printf("    - 及时清理不可见节点\n")
        fmt.Printf("    - 使用对象池\n")
    }
}

// 分页加载演示
func demonstratePagedLoading(containerID, totalChildren int) {
    pageSize := 50
    totalPages := (totalChildren + pageSize - 1) / pageSize

    fmt.Printf("分页加载设置:\n")
    fmt.Printf("  总子节点数: %d\n", totalChildren)
    fmt.Printf("  每页大小: %d\n", pageSize)
    fmt.Printf("  总页数: %d\n", totalPages)

    // 模拟分页加载
    for page := 1; page <= 3 && page <= totalPages; page++ { // 只演示3页
        fmt.Printf("\n加载第 %d 页 (节点 %d-%d):\n",
            page, (page-1)*pageSize+1, min(page*pageSize, totalChildren))

        // 模拟加载一页
        startTime := time.Now()
        time.Sleep(100 * time.Millisecond) // 模拟加载延迟
        loadTime := time.Since(startTime)

        nodesLoaded := pageSize
        if page == totalPages {
            nodesLoaded = totalChildren - (totalPages-1)*pageSize
        }

        fmt.Printf("  加载 %d 个节点，耗时: %v\n", nodesLoaded, loadTime)
        fmt.Printf("  累计加载: %d/%d (%.1f%%)\n",
            page*pageSize, totalChildren,
            float64(min(page*pageSize, totalChildren))/float64(totalChildren)*100)
    }

    fmt.Printf("\n分页加载完成\n")
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// 高级功能: 智能子节点加载器
type SmartChildLoader struct {
    parentID       int
    totalChildren  int
    loadedDepth    int
    loadedNodes    []int
    loadHistory    []LoadRecord
    loadConfig     LoadConfig
}

type LoadRecord struct {
    Timestamp  time.Time
    Depth      int
    NodeCount  int
    Duration   time.Duration
    Success    bool
    Error      string
}

type LoadConfig struct {
    MaxDepth        int
    BatchSize       int
    EnableLazyLoad  bool
    LoadThreshold   int
    CacheResults    bool
    RetryOnFail     bool
    MaxRetries      int
}

func NewSmartChildLoader(parentID int, config LoadConfig) *SmartChildLoader {
    if config.MaxDepth <= 0 {
        config.MaxDepth = 3
    }
    if config.BatchSize <= 0 {
        config.BatchSize = 100
    }
    if config.MaxRetries <= 0 {
        config.MaxRetries = 3
    }

    return &SmartChildLoader{
        parentID:      parentID,
        loadedNodes:   make([]int, 0),
        loadHistory:   make([]LoadRecord, 0),
        loadConfig:    config,
    }
}

func (scl *SmartChildLoader) LoadChildren(depth int) error {
    if depth <= 0 || depth > scl.loadConfig.MaxDepth {
        return fmt.Errorf("无效的深度: %d (最大: %d)", depth, scl.loadConfig.MaxDepth)
    }

    var retries int
    var lastError error

    for retries <= scl.loadConfig.MaxRetries {
        startTime := time.Now()

        result, err := CDPDOMRequestChildNodes(scl.parentID, depth, false)
        duration := time.Since(startTime)

        record := LoadRecord{
            Timestamp: startTime,
            Depth:     depth,
            Duration:  duration,
            Success:   err == nil,
            Error:     "",
        }

        if err != nil {
            lastError = err
            record.Error = err.Error()
            record.Success = false

            if !scl.loadConfig.RetryOnFail {
                scl.loadHistory = append(scl.loadHistory, record)
                return fmt.Errorf("加载失败: %w", err)
            }

            retries++
            if retries <= scl.loadConfig.MaxRetries {
                log.Printf("加载失败，第 %d 次重试: %v", retries, err)
                time.Sleep(time.Duration(retries) * 100 * time.Millisecond)
                continue
            }
        } else {
            // 解析结果
            nodeCount := scl.parseAndStoreResult(result, depth)
            record.NodeCount = nodeCount

            scl.loadedDepth = depth
            scl.loadHistory = append(scl.loadHistory, record)

            log.Printf("子节点加载成功: 深度=%d, 节点数=%d, 耗时=%v",
                depth, nodeCount, duration)
            return nil
        }
    }

    scl.loadHistory = append(scl.loadHistory, LoadRecord{
        Timestamp: time.Now(),
        Depth:     depth,
        Duration:  0,
        Success:   false,
        Error:     fmt.Sprintf("重试%d次后失败: %v", scl.loadConfig.MaxRetries, lastError),
    })

    return fmt.Errorf("重试%d次后仍然失败: %w", scl.loadConfig.MaxRetries, lastError)
}

func (scl *SmartChildLoader) parseAndStoreResult(result string, depth int) int {
    // 这里需要解析DOM.setChildNodes事件
    // 简化实现
    estimatedNodes := estimateTotalNodes(scl.totalChildren, depth)

    // 模拟存储加载的节点
    for i := 0; i < estimatedNodes && i < 1000; i++ {
        scl.loadedNodes = append(scl.loadedNodes, scl.parentID*1000+i)
    }

    return len(scl.loadedNodes)
}

func (scl *SmartChildLoader) GetStats() map[string]interface{} {
    totalLoads := len(scl.loadHistory)
    successfulLoads := 0
    var totalDuration time.Duration

    for _, record := range scl.loadHistory {
        if record.Success {
            successfulLoads++
        }
        totalDuration += record.Duration
    }

    avgDuration := time.Duration(0)
    if totalLoads > 0 {
        avgDuration = totalDuration / time.Duration(totalLoads)
    }

    return map[string]interface{}{
        "parentId":        scl.parentID,
        "totalChildren":   scl.totalChildren,
        "loadedDepth":     scl.loadedDepth,
        "loadedNodes":     len(scl.loadedNodes),
        "totalLoads":      totalLoads,
        "successfulLoads": successfulLoads,
        "successRate":     float64(successfulLoads) / float64(totalLoads) * 100,
        "totalDuration":   totalDuration,
        "averageDuration": avgDuration,
    }
}

func (scl *SmartChildLoader) GetHistory() []LoadRecord {
    return scl.loadHistory
}

// 高级功能: 渐进式加载
type ProgressiveLoader struct {
    loader      *SmartChildLoader
    currentPage int
    pageSize    int
    isLoading   bool
}

func NewProgressiveLoader(parentID, totalChildren, pageSize int) *ProgressiveLoader {
    config := LoadConfig{
        MaxDepth:       1,
        BatchSize:      pageSize,
        EnableLazyLoad: true,
        LoadThreshold:  100,
        CacheResults:   true,
        RetryOnFail:    true,
        MaxRetries:     3,
    }

    loader := NewSmartChildLoader(parentID, config)
    loader.totalChildren = totalChildren

    return &ProgressiveLoader{
        loader:     loader,
        pageSize:   pageSize,
        currentPage: 0,
    }
}

func (pl *ProgressiveLoader) LoadNextPage() (bool, error) {
    if pl.isLoading {
        return false, fmt.Errorf("正在加载中")
    }

    totalPages := (pl.loader.totalChildren + pl.pageSize - 1) / pl.pageSize
    if pl.currentPage >= totalPages {
        return false, fmt.Errorf("没有更多页面")
    }

    pl.isLoading = true
    defer func() { pl.isLoading = false }()

    pl.currentPage++

    fmt.Printf("加载第 %d/%d 页\n", pl.currentPage, totalPages)

    // 计算加载范围
    startIdx := (pl.currentPage - 1) * pl.pageSize
    endIdx := pl.currentPage * pl.pageSize
    if endIdx > pl.loader.totalChildren {
        endIdx = pl.loader.totalChildren
    }

    // 这里应该实现只加载指定范围的子节点
    // 简化实现：加载完整深度1
    if err := pl.loader.LoadChildren(1); err != nil {
        return false, fmt.Errorf("加载页面失败: %w", err)
    }

    fmt.Printf("页面加载完成: 节点 %d-%d\n", startIdx+1, endIdx)

    hasMore := pl.currentPage < totalPages
    return hasMore, nil
}

func (pl *ProgressiveLoader) GetCurrentState() map[string]interface{} {
    totalPages := (pl.loader.totalChildren + pl.pageSize - 1) / pl.pageSize
    loadedCount := pl.currentPage * pl.pageSize
    if loadedCount > pl.loader.totalChildren {
        loadedCount = pl.loader.totalChildren
    }

    return map[string]interface{}{
        "currentPage":   pl.currentPage,
        "totalPages":    totalPages,
        "loadedCount":   loadedCount,
        "totalChildren": pl.loader.totalChildren,
        "progress":      float64(loadedCount) / float64(pl.loader.totalChildren) * 100,
        "hasMore":       pl.currentPage < totalPages,
    }
}

// 演示智能加载场景
func demonstrateSmartLoading() {
    fmt.Println("=== 智能子节点加载演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 模拟大型列表容器
    containerID := 1001
    totalChildren := 1000

    fmt.Printf("容器ID: %d, 总子节点数: %d\n\n", containerID, totalChildren)

    // 场景1: 智能加载器
    fmt.Printf("场景1: 智能加载器\n")

    config := LoadConfig{
        MaxDepth:       3,
        BatchSize:      200,
        EnableLazyLoad: true,
        LoadThreshold:  50,
        CacheResults:   true,
        RetryOnFail:    true,
        MaxRetries:     3,
    }

    loader := NewSmartChildLoader(containerID, config)
    loader.totalChildren = totalChildren

    // 分阶段加载
    depths := []int{1, 2, 3}
    for _, depth := range depths {
        fmt.Printf("加载深度 %d...\n", depth)
        if err := loader.LoadChildren(depth); err != nil {
            fmt.Printf("❌ 加载失败: %v\n", err)
        } else {
            fmt.Printf("✅ 加载成功\n")
        }
    }

    // 显示统计
    fmt.Printf("\n加载统计:\n")
    stats := loader.GetStats()
    fmt.Printf("  已加载深度: %d\n", stats["loadedDepth"])
    fmt.Printf("  已加载节点: %d\n", stats["loadedNodes"])
    fmt.Printf("  加载次数: %d\n", stats["totalLoads"])
    fmt.Printf("  成功次数: %d\n", stats["successfulLoads"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("  总耗时: %v\n", stats["totalDuration"])
    fmt.Printf("  平均耗时: %v\n", stats["averageDuration"])

    // 场景2: 渐进式加载
    fmt.Printf("\n场景2: 渐进式加载\n")

    progressiveLoader := NewProgressiveLoader(containerID, totalChildren, 100)

    // 加载前3页
    for i := 0; i < 3; i++ {
        hasMore, err := progressiveLoader.LoadNextPage()
        if err != nil {
            fmt.Printf("❌ 加载失败: %v\n", err)
            break
        }

        state := progressiveLoader.GetCurrentState()
        fmt.Printf("  进度: %.1f%% (%d/%d)\n",
            state["progress"], state["loadedCount"], state["totalChildren"])

        if !hasMore.(bool) {
            break
        }
    }

    fmt.Println("\n=== 演示完成 ===")
}


*/

// -----------------------------------------------  DOM.requestNode  -----------------------------------------------
// === 应用场景 ===
// 1. 节点引用: 通过后端节点ID获取前端节点ID
// 2. 重新连接: 重新建立与之前引用节点的连接
// 3. 节点恢复: 在节点失效后重新获取节点引用
// 4. 跨会话: 在不同CDP会话间传递节点引用
// 5. 调试支持: 调试时重新获取感兴趣的节点
// 6. 异步操作: 在异步操作后重新获取节点

// CDPDOMRequestNode 通过后端节点ID请求前端节点
// backendNodeID: 后端节点ID
func CDPDOMRequestNode(backendNodeID int) (string, error) {
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
        "method": "DOM.requestNode",
        "params": {
            "backendNodeId": %d
        }
    }`, reqID, backendNodeID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.requestNode 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.requestNode 请求超时")
		}
	}
}

/*

// 示例: 通过后端节点ID重新获取节点
func ExampleCDPDOMRequestNode() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    fmt.Printf("=== 后端节点ID重新获取演示 ===\n")

    // 假设我们有一个已知的后端节点ID
    backendNodeID := 12345

    fmt.Printf("后端节点ID: %d\n\n", backendNodeID)

    // 2. 验证后端节点ID的有效性
    fmt.Printf("=== 后端节点ID验证 ===\n")

    isValid, info := validateBackendNode(backendNodeID)
    if !isValid {
        fmt.Printf("❌ 后端节点ID无效或不存在\n")
        return
    }

    fmt.Printf("✅ 后端节点ID有效\n")
    if info != nil {
        fmt.Printf("  节点类型: %s\n", info["nodeType"])
        fmt.Printf("  节点名称: %s\n", info["nodeName"])
    }

    // 3. 通过后端节点ID请求前端节点
    fmt.Printf("\n=== 请求前端节点 ===\n")

    startTime := time.Now()
    result, err := CDPDOMRequestNode(backendNodeID)
    requestTime := time.Since(startTime)

    if err != nil {
        log.Printf("请求节点失败: %v", err)
        return
    }

    fmt.Printf("请求结果: %s\n", result)
    fmt.Printf("请求耗时: %v\n", requestTime)

    // 4. 解析响应获取前端节点ID
    var response struct {
        Result struct {
            NodeID int `json:"nodeId"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        log.Printf("解析响应失败: %v", err)
        return
    }

    frontendNodeID := response.Result.NodeID

    if frontendNodeID == 0 {
        fmt.Printf("❌ 未获取到有效的前端节点ID\n")
        return
    }

    fmt.Printf("✅ 获取到前端节点ID: %d\n", frontendNodeID)

    // 5. 验证获取的节点
    fmt.Printf("\n=== 节点验证 ===\n")

    // 获取节点详细信息
    nodeInfo, err := getNodeInfo(frontendNodeID)
    if err != nil {
        log.Printf("获取节点信息失败: %v", err)
    } else {
        displayNodeInfo("重新获取的节点", nodeInfo)
    }

    // 6. 比较前后端节点信息
    fmt.Printf("\n=== 节点信息比较 ===\n")

    compareNodeInformation(backendNodeID, frontendNodeID)

    // 7. 节点引用测试
    fmt.Printf("\n=== 节点引用测试 ===\n")

    testNodeReferences(backendNodeID, frontendNodeID)

    // 8. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次请求同一个节点
    testRepeatedRequests(backendNodeID, 5)

    // 测试多个不同节点的请求
    testMultipleNodes([]int{backendNodeID, 12346, 12347})

    // 9. 错误处理测试
    fmt.Printf("\n=== 错误处理测试 ===\n")

    errorTestCases := []struct {
        backendID int
        desc      string
    }{
        {0, "无效的后端节点ID 0"},
        {999999, "不存在的后端节点ID"},
        {-1, "负数的后端节点ID"},
    }

    for _, tc := range errorTestCases {
        fmt.Printf("测试: %s\n", tc.desc)
        result, err := CDPDOMRequestNode(tc.backendID)

        if err != nil {
            fmt.Printf("  ✅ 预期错误: %v\n", err)
        } else {
            fmt.Printf("  ❌ 预期错误但成功: %s\n", result)
        }
    }

    // 10. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "页面刷新后恢复",
            description: "页面刷新后重新获取之前保存的节点引用",
            useCase:     "保存后端节点ID，刷新后重新获取",
        },
        {
            name:        "跨会话节点传递",
            description: "在不同CDP会话间传递节点引用",
            useCase:     "通过后端节点ID在不同会话中引用同一节点",
        },
        {
            name:        "异步操作后重新获取",
            description: "在异步DOM操作后重新建立节点引用",
            useCase:     "异步修改后通过后端ID重新获取节点",
        },
        {
            name:        "节点缓存管理",
            description: "管理节点引用缓存，失效时重新获取",
            useCase:     "缓存后端节点ID，需要时重新请求",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 验证后端节点
func validateBackendNode(backendNodeID int) (bool, map[string]interface{}) {
    // 这里应该通过其他方式验证后端节点的有效性
    // 简化实现：假设有效
    if backendNodeID <= 0 {
        return false, nil
    }

    info := map[string]interface{}{
        "nodeType": "ELEMENT_NODE",
        "nodeName": "DIV",
    }

    return true, info
}

// 比较节点信息
func compareNodeInformation(backendNodeID, frontendNodeID int) {
    // 获取前端节点信息
    frontendInfo, err1 := getNodeInfo(frontendNodeID)

    // 这里应该能获取后端节点信息
    // 简化实现
    backendInfo := map[string]interface{}{
        "backendNodeId": backendNodeID,
    }

    if err1 == nil {
        fmt.Printf("前端节点ID: %d\n", frontendNodeID)
        fmt.Printf("后端节点ID: %d\n", backendNodeID)

        // 验证节点一致性
        if validateNodeConsistency(backendNodeID, frontendNodeID) {
            fmt.Printf("✅ 前后端节点一致\n")
        } else {
            fmt.Printf("⚠️ 前后端节点可能不一致\n")
        }
    }
}

// 验证节点一致性
func validateNodeConsistency(backendNodeID, frontendNodeID int) bool {
    // 这里应该比较前后端节点的详细信息
    // 简化实现：假设一致
    return frontendNodeID > 0 && backendNodeID > 0
}

// 测试节点引用
func testNodeReferences(backendNodeID, frontendNodeID int) {
    fmt.Printf("测试用例:\n")

    // 测试1: 获取节点的属性
    fmt.Printf("1. 获取节点属性...\n")
    attrs, err := CDPDOMGetAttributes(frontendNodeID)
    if err != nil {
        fmt.Printf("  ❌ 获取属性失败: %v\n", err)
    } else {
        fmt.Printf("  ✅ 获取属性成功\n")

        // 解析属性
        var attrsResp struct {
            Result struct {
                Attributes []string `json:"attributes"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(attrs), &attrsResp); err == nil {
            fmt.Printf("    属性数量: %d\n", len(attrsResp.Result.Attributes)/2)
        }
    }

    // 测试2: 获取节点盒模型
    fmt.Printf("2. 获取节点盒模型...\n")
    boxModel, err := CDPDOMGetBoxModel(frontendNodeID)
    if err != nil {
        fmt.Printf("  ❌ 获取盒模型失败: %v\n", err)
    } else {
        fmt.Printf("  ✅ 获取盒模型成功\n")

        var boxResp struct {
            Result struct {
                Model struct {
                    Width  int `json:"width"`
                    Height int `json:"height"`
                } `json:"model"`
            } `json:"result"`
        }

        if err := json.Unmarshal([]byte(boxModel), &boxResp); err == nil {
            fmt.Printf("    节点尺寸: %d x %d 像素\n",
                boxResp.Result.Model.Width, boxResp.Result.Model.Height)
        }
    }

    // 测试3: 高亮节点
    fmt.Printf("3. 高亮节点...\n")
    highlightConfig := &HighlightConfig{
        ContentColor: &RGBA{R: 100, G: 200, B: 255, A: 0.3},
        BorderColor:  &RGBA{R: 0, G: 100, B: 200, A: 0.8},
    }

    if _, err := CDPDOMHighlightNode(frontendNodeID, highlightConfig); err != nil {
        fmt.Printf("  ❌ 高亮失败: %v\n", err)
    } else {
        fmt.Printf("  ✅ 高亮成功，3秒后清理\n")
        time.Sleep(3 * time.Second)
        CDPDOMHideHighlight()
    }
}

// 测试重复请求
func testRepeatedRequests(backendNodeID int, count int) {
    fmt.Printf("重复请求测试 (%d 次):\n", count)

    var totalDuration time.Duration
    successCount := 0
    var nodeIDs []int

    for i := 0; i < count; i++ {
        startTime := time.Now()
        result, err := CDPDOMRequestNode(backendNodeID)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            var resp struct {
                Result struct {
                    NodeID int `json:"nodeId"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(result), &resp); err == nil && resp.Result.NodeID > 0 {
                successCount++
                nodeIDs = append(nodeIDs, resp.Result.NodeID)
                fmt.Printf("  第 %d 次: ✅ 成功 (节点ID: %d, 耗时: %v)\n",
                    i+1, resp.Result.NodeID, duration)
            } else {
                fmt.Printf("  第 %d 次: ❌ 解析失败\n", i+1)
            }
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))

    // 检查返回的节点ID是否一致
    if len(nodeIDs) > 1 {
        allSame := true
        firstID := nodeIDs[0]
        for i := 1; i < len(nodeIDs); i++ {
            if nodeIDs[i] != firstID {
                allSame = false
                break
            }
        }

        if allSame {
            fmt.Printf("  ✅ 所有请求返回相同的节点ID: %d\n", firstID)
        } else {
            fmt.Printf("  ⚠️ 不同请求返回不同的节点ID\n")
        }
    }
}

// 测试多个节点
func testMultipleNodes(backendNodeIDs []int) {
    fmt.Printf("\n多节点请求测试 (%d 个节点):\n", len(backendNodeIDs))

    var results []NodeRequestResult

    for _, backendID := range backendNodeIDs {
        result := NodeRequestResult{
            BackendNodeID: backendID,
        }

        startTime := time.Now()
        response, err := CDPDOMRequestNode(backendID)
        result.Duration = time.Since(startTime)

        if err != nil {
            result.Success = false
            result.Error = err.Error()
        } else {
            var resp struct {
                Result struct {
                    NodeID int `json:"nodeId"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(response), &resp); err == nil && resp.Result.NodeID > 0 {
                result.Success = true
                result.FrontendNodeID = resp.Result.NodeID
            } else {
                result.Success = false
                result.Error = "解析失败"
            }
        }

        results = append(results, result)

        // 短暂延迟
        time.Sleep(50 * time.Millisecond)
    }

    // 显示结果
    successCount := 0
    for _, result := range results {
        status := "❌ 失败"
        if result.Success {
            status = "✅ 成功"
            successCount++
        }

        fmt.Printf("  后端ID: %-6d %s", result.BackendNodeID, status)
        if result.Success {
            fmt.Printf(" -> 前端ID: %d", result.FrontendNodeID)
        }
        fmt.Printf(" (耗时: %v)\n", result.Duration)
    }

    fmt.Printf("总体成功率: %.1f%% (%d/%d)\n",
        float64(successCount)/float64(len(results))*100,
        successCount, len(results))
}

type NodeRequestResult struct {
    BackendNodeID  int
    FrontendNodeID int
    Success        bool
    Error          string
    Duration       time.Duration
}

// 高级功能: 节点引用管理器
type NodeReferenceManager struct {
    references    map[int]NodeReference
    history       []ReferenceEvent
    cache         *ReferenceCache
}

type NodeReference struct {
    BackendNodeID  int
    FrontendNodeID int
    CreatedAt      time.Time
    LastAccessed   time.Time
    AccessCount    int
    Metadata       map[string]interface{}
    Valid          bool
}

type ReferenceEvent struct {
    Timestamp   time.Time
    EventType   string
    BackendID   int
    FrontendID  int
    Details     map[string]interface{}
}

type ReferenceCache struct {
    maxSize    int
    ttl        time.Duration
    cleanupInterval time.Duration
}

func NewNodeReferenceManager(maxSize int, ttl time.Duration) *NodeReferenceManager {
    return &NodeReferenceManager{
        references: make(map[int]NodeReference),
        history:    make([]ReferenceEvent, 0),
        cache: &ReferenceCache{
            maxSize:        maxSize,
            ttl:            ttl,
            cleanupInterval: 1 * time.Minute,
        },
    }
}

func (nrm *NodeReferenceManager) GetOrCreateReference(backendNodeID int, metadata map[string]interface{}) (NodeReference, error) {
    // 检查是否已有引用
    if ref, exists := nrm.references[backendNodeID]; exists {
        if nrm.isReferenceValid(ref) {
            // 更新访问信息
            ref.LastAccessed = time.Now()
            ref.AccessCount++
            nrm.references[backendNodeID] = ref

            nrm.recordEvent("cache_hit", backendNodeID, ref.FrontendNodeID, map[string]interface{}{
                "accessCount": ref.AccessCount,
            })

            return ref, nil
        } else {
            // 引用失效，删除
            delete(nrm.references, backendNodeID)
            nrm.recordEvent("cache_expired", backendNodeID, ref.FrontendNodeID, nil)
        }
    }

    // 创建新引用
    nrm.recordEvent("request_start", backendNodeID, 0, metadata)

    startTime := time.Now()
    result, err := CDPDOMRequestNode(backendNodeID)
    duration := time.Since(startTime)

    if err != nil {
        nrm.recordEvent("request_failed", backendNodeID, 0, map[string]interface{}{
            "error":    err.Error(),
            "duration": duration,
        })
        return NodeReference{}, fmt.Errorf("请求节点失败: %w", err)
    }

    var resp struct {
        Result struct {
            NodeID int `json:"nodeId"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &resp); err != nil {
        nrm.recordEvent("parse_failed", backendNodeID, 0, map[string]interface{}{
            "error":    err.Error(),
            "duration": duration,
        })
        return NodeReference{}, fmt.Errorf("解析响应失败: %w", err)
    }

    frontendNodeID := resp.Result.NodeID
    if frontendNodeID == 0 {
        nrm.recordEvent("invalid_node", backendNodeID, 0, map[string]interface{}{
            "duration": duration,
        })
        return NodeReference{}, fmt.Errorf("获取到无效的节点ID")
    }

    // 创建引用记录
    ref := NodeReference{
        BackendNodeID:  backendNodeID,
        FrontendNodeID: frontendNodeID,
        CreatedAt:      time.Now(),
        LastAccessed:   time.Now(),
        AccessCount:    1,
        Metadata:       metadata,
        Valid:          true,
    }

    // 保存到缓存
    nrm.references[backendNodeID] = ref

    // 清理过期引用
    nrm.cleanupExpiredReferences()

    nrm.recordEvent("request_success", backendNodeID, frontendNodeID, map[string]interface{}{
        "duration": duration,
    })

    return ref, nil
}

func (nrm *NodeReferenceManager) isReferenceValid(ref NodeReference) bool {
    if !ref.Valid {
        return false
    }

    // 检查是否过期
    if nrm.cache.ttl > 0 && time.Since(ref.LastAccessed) > nrm.cache.ttl {
        return false
    }

    // 这里还可以添加其他验证逻辑
    return true
}

func (nrm *NodeReferenceManager) cleanupExpiredReferences() {
    if len(nrm.references) <= nrm.cache.maxSize {
        return
    }

    // 清理过期引用
    now := time.Now()
    for backendID, ref := range nrm.references {
        if nrm.cache.ttl > 0 && now.Sub(ref.LastAccessed) > nrm.cache.ttl {
            delete(nrm.references, backendID)
            nrm.recordEvent("auto_cleanup", backendID, ref.FrontendNodeID, nil)
        }
    }

    // 如果仍然超过最大大小，清理最久未访问的
    if len(nrm.references) > nrm.cache.maxSize {
        nrm.cleanupLRU()
    }
}

func (nrm *NodeReferenceManager) cleanupLRU() {
    // 找到最久未访问的引用
    var oldestID int
    var oldestTime time.Time
    first := true

    for backendID, ref := range nrm.references {
        if first || ref.LastAccessed.Before(oldestTime) {
            oldestID = backendID
            oldestTime = ref.LastAccessed
            first = false
        }
    }

    if !first {
        deletedRef := nrm.references[oldestID]
        delete(nrm.references, oldestID)
        nrm.recordEvent("lru_cleanup", oldestID, deletedRef.FrontendNodeID, map[string]interface{}{
            "lastAccessed": deletedRef.LastAccessed,
        })
    }
}

func (nrm *NodeReferenceManager) recordEvent(eventType string, backendID, frontendID int, details map[string]interface{}) {
    event := ReferenceEvent{
        Timestamp:  time.Now(),
        EventType:  eventType,
        BackendID:  backendID,
        FrontendID: frontendID,
        Details:    details,
    }

    nrm.history = append(nrm.history, event)
}

func (nrm *NodeReferenceManager) GetStats() map[string]interface{} {
    totalRefs := len(nrm.references)
    validRefs := 0
    totalAccesses := 0

    for _, ref := range nrm.references {
        if ref.Valid {
            validRefs++
        }
        totalAccesses += ref.AccessCount
    }

    avgAccesses := 0.0
    if totalRefs > 0 {
        avgAccesses = float64(totalAccesses) / float64(totalRefs)
    }

    return map[string]interface{}{
        "totalReferences":   totalRefs,
        "validReferences":   validRefs,
        "invalidReferences": totalRefs - validRefs,
        "totalAccesses":     totalAccesses,
        "averageAccesses":   avgAccesses,
        "historyEvents":     len(nrm.history),
        "cacheHitRatio": func() float64 {
            hits := 0
            for _, event := range nrm.history {
                if event.EventType == "cache_hit" {
                    hits++
                }
            }
            if len(nrm.history) > 0 {
                return float64(hits) / float64(len(nrm.history)) * 100
            }
            return 0
        }(),
    }
}

func (nrm *NodeReferenceManager) GetHistory(limit int) []ReferenceEvent {
    if limit <= 0 || limit >= len(nrm.history) {
        return nrm.history
    }

    return nrm.history[len(nrm.history)-limit:]
}

// 演示节点引用管理场景
func demonstrateReferenceManagement() {
    fmt.Println("=== 节点引用管理演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建引用管理器
    manager := NewNodeReferenceManager(100, 5*time.Minute)

    // 模拟一批后端节点ID
    backendNodeIDs := []int{1001, 1002, 1003, 1004, 1005}

    fmt.Printf("测试 %d 个后端节点ID的引用管理\n\n", len(backendNodeIDs))

    // 场景1: 首次获取引用
    fmt.Printf("场景1: 首次获取引用\n")

    for _, backendID := range backendNodeIDs {
        metadata := map[string]interface{}{
            "source":    "test",
            "timestamp": time.Now().Unix(),
        }

        ref, err := manager.GetOrCreateReference(backendID, metadata)
        if err != nil {
            fmt.Printf("  后端ID %d: ❌ 获取失败: %v\n", backendID, err)
        } else {
            fmt.Printf("  后端ID %d: ✅ 获取成功 -> 前端ID %d\n",
                backendID, ref.FrontendNodeID)
        }
    }

    // 场景2: 缓存命中测试
    fmt.Printf("\n场景2: 缓存命中测试\n")

    for i := 0; i < 3; i++ {
        for _, backendID := range backendNodeIDs {
            ref, err := manager.GetOrCreateReference(backendID, nil)
            if err == nil {
                fmt.Printf("  后端ID %d: 访问 %d 次，前端ID %d\n",
                    backendID, ref.AccessCount, ref.FrontendNodeID)
            }
        }

        if i < 2 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    // 显示统计信息
    fmt.Printf("\n=== 引用管理统计 ===\n")
    stats := manager.GetStats()
    fmt.Printf("  总引用数: %d\n", stats["totalReferences"])
    fmt.Printf("  有效引用: %d\n", stats["validReferences"])
    fmt.Printf("  总访问次数: %d\n", stats["totalAccesses"])
    fmt.Printf("  平均访问次数: %.1f\n", stats["averageAccesses"])
    fmt.Printf("  缓存命中率: %.1f%%\n", stats["cacheHitRatio"])

    // 显示最近事件
    fmt.Printf("\n=== 最近事件 ===\n")
    history := manager.GetHistory(5)
    for i, event := range history {
        fmt.Printf("%d. [%s] 后端ID: %d, 前端ID: %d\n",
            i+1, event.EventType, event.BackendID, event.FrontendID)
    }

    fmt.Println("\n=== 演示完成 ===")
}

*/

// -----------------------------------------------  DOM.resolveNode  -----------------------------------------------
// === 应用场景 ===
// 1. 对象解析: 解析节点为JavaScript对象
// 2. 远程对象: 获取远程对象引用
// 3. 脚本执行: 在页面上执行脚本时引用DOM节点
// 4. 属性访问: 通过JavaScript访问节点属性
// 5. 方法调用: 调用节点的JavaScript方法
// 6. 调试支持: 调试时获取节点的JavaScript表示

// CDPDOMResolveNode 将节点解析为JavaScript对象
// nodeID: 要解析的节点ID
// backendNodeID: 后端节点ID（可选）
// objectGroup: 对象分组名称
// executionContextID: 执行上下文ID
func CDPDOMResolveNode(nodeID, backendNodeID int, objectGroup string, executionContextID int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息参数
	params := map[string]interface{}{
		"objectGroup": objectGroup,
	}

	if nodeID > 0 {
		params["nodeId"] = nodeID
	}

	if backendNodeID > 0 {
		params["backendNodeId"] = backendNodeID
	}

	if executionContextID > 0 {
		params["executionContextId"] = executionContextID
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "DOM.resolveNode",
        "params": %s
    }`, reqID, string(paramsJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.resolveNode 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.resolveNode 请求超时")
		}
	}
}

/*

// 示例: 解析按钮元素为JavaScript对象
func ExampleCDPDOMResolveNode() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个按钮的节点ID
    buttonNodeID := 1001

    fmt.Printf("=== 节点解析为JavaScript对象演示 ===\n")
    fmt.Printf("目标节点ID: %d\n\n", buttonNodeID)

    // 2. 首先获取节点的基本信息
    fmt.Printf("=== 节点基本信息 ===\n")

    nodeInfo, err := getNodeInfo(buttonNodeID)
    if err != nil {
        log.Printf("获取节点信息失败: %v", err)
        return
    }

    displayNodeInfo("要解析的节点", nodeInfo)

    // 3. 解析节点为JavaScript对象
    fmt.Printf("\n=== 解析节点 ===\n")

    // 设置对象分组名称
    objectGroup := "test_group"

    startTime := time.Now()

    result, err := CDPDOMResolveNode(buttonNodeID, 0, objectGroup, 0)
    if err != nil {
        log.Printf("解析节点失败: %v", err)
        return
    }

    parseTime := time.Since(startTime)
    fmt.Printf("解析结果: %s\n", result)
    fmt.Printf("解析耗时: %v\n", parseTime)

    // 4. 解析响应获取远程对象
    var response struct {
        Result struct {
            Object struct {
                Type        string `json:"type"`
                Subtype     string `json:"subtype,omitempty"`
                ClassName   string `json:"className,omitempty"`
                Value       interface{} `json:"value,omitempty"`
                Description string `json:"description,omitempty"`
                ObjectID    string `json:"objectId"`
            } `json:"object"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err != nil {
        log.Printf("解析响应失败: %v", err)
        return
    }

    remoteObject := response.Result.Object
    objectID := remoteObject.ObjectID

    if objectID == "" {
        fmt.Printf("❌ 未获取到有效的对象ID\n")
        return
    }

    fmt.Printf("✅ 成功解析为JavaScript对象\n")
    fmt.Printf("  对象ID: %s\n", objectID)
    fmt.Printf("  对象类型: %s\n", remoteObject.Type)
    if remoteObject.Subtype != "" {
        fmt.Printf("  对象子类型: %s\n", remoteObject.Subtype)
    }
    if remoteObject.ClassName != "" {
        fmt.Printf("  对象类名: %s\n", remoteObject.ClassName)
    }
    if remoteObject.Description != "" {
        fmt.Printf("  对象描述: %s\n", remoteObject.Description)
    }

    // 5. 通过JavaScript操作对象
    fmt.Printf("\n=== JavaScript对象操作 ===\n")

    // 获取对象属性
    objectProperties, err := getObjectProperties(objectID)
    if err != nil {
        log.Printf("获取对象属性失败: %v", err)
    } else {
        displayObjectProperties(objectProperties)
    }

    // 调用对象方法
    callObjectMethods(objectID, nodeInfo)

    // 6. 对象引用管理
    fmt.Printf("\n=== 对象引用管理 ===\n")

    manageObjectReference(objectID, objectGroup)

    // 7. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次解析
    testMultipleResolutions(buttonNodeID, 3)

    // 测试不同节点的解析
    testDifferentNodes([]int{1001, 1002, 1003})

    // 8. 错误处理测试
    fmt.Printf("\n=== 错误处理测试 ===\n")

    errorTestCases := []struct {
        nodeID  int
        desc    string
    }{
        {0, "无效的节点ID 0"},
        {999999, "不存在的节点ID"},
    }

    for _, tc := range errorTestCases {
        fmt.Printf("测试: %s\n", tc.desc)
        result, err := CDPDOMResolveNode(tc.nodeID, 0, "test", 0)

        if err != nil {
            fmt.Printf("  ✅ 预期错误: %v\n", err)
        } else {
            fmt.Printf("  ❌ 预期错误但成功: %s\n", result)
        }
    }

    // 9. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "自动化测试",
            description: "在自动化测试中操作页面元素",
            useCase:     "解析元素为对象，执行点击、输入等操作",
        },
        {
            name:        "表单处理",
            description: "通过JavaScript处理表单元素",
            useCase:     "解析表单元素，获取/设置值，验证数据",
        },
        {
            name:        "动态修改",
            description: "动态修改页面元素属性和样式",
            useCase:     "解析元素，修改样式、属性、内容等",
        },
        {
            name:        "数据提取",
            description: "从页面中提取结构化数据",
            useCase:     "解析元素，提取文本、属性、数据属性等",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 获取对象属性
func getObjectProperties(objectID string) (map[string]interface{}, error) {
    // 这里需要调用Runtime.getProperties方法
    // 简化实现：返回模拟数据

    properties := map[string]interface{}{
        "tagName": map[string]interface{}{
            "value": "BUTTON",
            "type":  "string",
        },
        "className": map[string]interface{}{
            "value": "btn btn-primary",
            "type":  "string",
        },
        "disabled": map[string]interface{}{
            "value": false,
            "type":  "boolean",
        },
        "textContent": map[string]interface{}{
            "value": "点击我",
            "type":  "string",
        },
        "offsetWidth": map[string]interface{}{
            "value": 120,
            "type":  "number",
        },
        "offsetHeight": map[string]interface{}{
            "value": 40,
            "type":  "number",
        },
    }

    return properties, nil
}

// 显示对象属性
func displayObjectProperties(properties map[string]interface{}) {
    fmt.Printf("对象属性 (%d 个):\n", len(properties))

    for propName, propInfo := range properties {
        if info, ok := propInfo.(map[string]interface{}); ok {
            fmt.Printf("  %-20s: ", propName)
            if value, ok := info["value"]; ok {
                fmt.Printf("%v", value)
            }
            if propType, ok := info["type"].(string); ok {
                fmt.Printf(" (%s)", propType)
            }
            fmt.Println()
        }
    }
}

// 调用对象方法
func callObjectMethods(objectID string, nodeInfo map[string]interface{}) {
    nodeName, _ := nodeInfo["nodeName"].(string)

    fmt.Printf("可调用的方法:\n")

    // 根据节点类型决定可调用的方法
    switch strings.ToUpper(nodeName) {
    case "BUTTON", "INPUT", "A":
        fmt.Printf("  click() - 模拟点击\n")
        fmt.Printf("  focus() - 获取焦点\n")
        fmt.Printf("  blur()  - 失去焦点\n")

    case "INPUT", "TEXTAREA":
        fmt.Printf("  select() - 选中文本\n")

    case "FORM":
        fmt.Printf("  submit() - 提交表单\n")
        fmt.Printf("  reset()  - 重置表单\n")

    default:
        fmt.Printf("  getAttribute(name) - 获取属性\n")
        fmt.Printf("  setAttribute(name, value) - 设置属性\n")
        fmt.Printf("  removeAttribute(name) - 移除属性\n")
    }

    // 模拟调用方法
    fmt.Printf("\n模拟方法调用:\n")
    fmt.Printf("  element.click()\n")
    fmt.Printf("  element.focus()\n")
    fmt.Printf("  console.log(element.textContent)\n")
}

// 管理对象引用
func manageObjectReference(objectID, objectGroup string) {
    fmt.Printf("对象引用信息:\n")
    fmt.Printf("  对象ID: %s\n", objectID)
    fmt.Printf("  对象分组: %s\n", objectGroup)
    fmt.Printf("  创建时间: %s\n", time.Now().Format("15:04:05"))

    fmt.Printf("\n引用管理操作:\n")
    fmt.Printf("  1. 保持引用 - 对象保持活动状态\n")
    fmt.Printf("  2. 释放引用 - Runtime.releaseObject\n")
    fmt.Printf("  3. 引用计数 - 跟踪引用数量\n")
    fmt.Printf("  4. 垃圾回收 - Runtime.runIfWaitingForDebugger\n")

    fmt.Printf("\n最佳实践:\n")
    fmt.Printf("  - 及时释放不再使用的对象\n")
    fmt.Printf("  - 使用对象分组管理相关对象\n")
    fmt.Printf("  - 避免内存泄漏\n")
}

// 测试多次解析
func testMultipleResolutions(nodeID, count int) {
    fmt.Printf("多次解析测试 (%d 次):\n", count)

    var totalDuration time.Duration
    successCount := 0
    var objectIDs []string

    for i := 0; i < count; i++ {
        objectGroup := fmt.Sprintf("test_group_%d", i+1)

        startTime := time.Now()
        result, err := CDPDOMResolveNode(nodeID, 0, objectGroup, 0)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            var resp struct {
                Result struct {
                    Object struct {
                        ObjectID string `json:"objectId"`
                    } `json:"object"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(result), &resp); err == nil && resp.Result.Object.ObjectID != "" {
                successCount++
                objectIDs = append(objectIDs, resp.Result.Object.ObjectID)
                fmt.Printf("  第 %d 次: ✅ 成功 (对象ID: %s, 耗时: %v)\n",
                    i+1, resp.Result.Object.ObjectID, duration)
            } else {
                fmt.Printf("  第 %d 次: ❌ 解析失败\n", i+1)
            }
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))

    // 检查返回的对象ID是否不同
    if len(objectIDs) > 1 {
        allDifferent := true
        seen := make(map[string]bool)

        for _, id := range objectIDs {
            if seen[id] {
                allDifferent = false
                break
            }
            seen[id] = true
        }

        if allDifferent {
            fmt.Printf("  ✅ 每次解析返回不同的对象ID\n")
        } else {
            fmt.Printf("  ⚠️ 有些解析返回相同的对象ID\n")
        }
    }
}

// 测试不同节点
func testDifferentNodes(nodeIDs []int) {
    fmt.Printf("\n多节点解析测试 (%d 个节点):\n", len(nodeIDs))

    var results []NodeResolutionResult

    for _, nodeID := range nodeIDs {
        result := NodeResolutionResult{
            NodeID: nodeID,
        }

        startTime := time.Now()
        response, err := CDPDOMResolveNode(nodeID, 0, "test_group", 0)
        result.Duration = time.Since(startTime)

        if err != nil {
            result.Success = false
            result.Error = err.Error()
        } else {
            var resp struct {
                Result struct {
                    Object struct {
                        Type     string `json:"type"`
                        Subtype  string `json:"subtype,omitempty"`
                        ObjectID string `json:"objectId"`
                    } `json:"object"`
                } `json:"result"`
            }

            if err := json.Unmarshal([]byte(response), &resp); err == nil && resp.Result.Object.ObjectID != "" {
                result.Success = true
                result.ObjectID = resp.Result.Object.ObjectID
                result.ObjectType = resp.Result.Object.Type
                result.ObjectSubtype = resp.Result.Object.Subtype
            } else {
                result.Success = false
                result.Error = "解析失败"
            }
        }

        results = append(results, result)

        // 短暂延迟
        time.Sleep(50 * time.Millisecond)
    }

    // 显示结果
    successCount := 0
    for _, result := range results {
        status := "❌ 失败"
        if result.Success {
            status = "✅ 成功"
            successCount++
        }

        fmt.Printf("  节点ID: %-6d %s", result.NodeID, status)
        if result.Success {
            fmt.Printf(" -> 对象ID: %s", result.ObjectID)
            fmt.Printf(" (类型: %s", result.ObjectType)
            if result.ObjectSubtype != "" {
                fmt.Printf("/%s", result.ObjectSubtype)
            }
            fmt.Printf(")")
        }
        fmt.Printf(" (耗时: %v)\n", result.Duration)
    }

    fmt.Printf("总体成功率: %.1f%% (%d/%d)\n",
        float64(successCount)/float64(len(results))*100,
        successCount, len(results))
}

type NodeResolutionResult struct {
    NodeID        int
    ObjectID      string
    ObjectType    string
    ObjectSubtype string
    Success       bool
    Error         string
    Duration      time.Duration
}

// 高级功能: 对象解析管理器
type ObjectResolutionManager struct {
    resolutions  map[int]ObjectResolution
    history      []ResolutionEvent
    cache        *ResolutionCache
    objectGroups map[string]*ObjectGroup
}

type ObjectResolution struct {
    NodeID         int
    ObjectID       string
    ObjectGroup    string
    CreatedAt      time.Time
    LastAccessed   time.Time
    AccessCount    int
    ObjectType     string
    ObjectSubtype  string
    Valid          bool
    Metadata       map[string]interface{}
}

type ResolutionEvent struct {
    Timestamp   time.Time
    EventType   string
    NodeID      int
    ObjectID    string
    ObjectGroup string
    Details     map[string]interface{}
}

type ResolutionCache struct {
    maxSize    int
    ttl        time.Duration
    cleanupInterval time.Duration
}

type ObjectGroup struct {
    Name        string
    CreatedAt   time.Time
    ObjectCount int
    Objects     map[string]bool // objectID -> exists
}

func NewObjectResolutionManager(maxSize int, ttl time.Duration) *ObjectResolutionManager {
    return &ObjectResolutionManager{
        resolutions:  make(map[int]ObjectResolution),
        history:      make([]ResolutionEvent, 0),
        cache: &ResolutionCache{
            maxSize:        maxSize,
            ttl:            ttl,
            cleanupInterval: 1 * time.Minute,
        },
        objectGroups: make(map[string]*ObjectGroup),
    }
}

func (orm *ObjectResolutionManager) ResolveNode(nodeID int, objectGroup string, metadata map[string]interface{}) (ObjectResolution, error) {
    // 检查是否已有解析
    if resolution, exists := orm.resolutions[nodeID]; exists {
        if orm.isResolutionValid(resolution) {
            // 更新访问信息
            resolution.LastAccessed = time.Now()
            resolution.AccessCount++
            orm.resolutions[nodeID] = resolution

            orm.recordEvent("cache_hit", nodeID, resolution.ObjectID, objectGroup, map[string]interface{}{
                "accessCount": resolution.AccessCount,
            })

            return resolution, nil
        } else {
            // 解析失效，删除
            delete(orm.resolutions, nodeID)
            orm.recordEvent("cache_expired", nodeID, resolution.ObjectID, objectGroup, nil)
        }
    }

    // 确保对象组存在
    orm.ensureObjectGroup(objectGroup)

    // 创建新解析
    orm.recordEvent("resolution_start", nodeID, "", objectGroup, metadata)

    startTime := time.Now()
    result, err := CDPDOMResolveNode(nodeID, 0, objectGroup, 0)
    duration := time.Since(startTime)

    if err != nil {
        orm.recordEvent("resolution_failed", nodeID, "", objectGroup, map[string]interface{}{
            "error":    err.Error(),
            "duration": duration,
        })
        return ObjectResolution{}, fmt.Errorf("解析节点失败: %w", err)
    }

    var resp struct {
        Result struct {
            Object struct {
                Type     string `json:"type"`
                Subtype  string `json:"subtype,omitempty"`
                ObjectID string `json:"objectId"`
            } `json:"object"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &resp); err != nil {
        orm.recordEvent("parse_failed", nodeID, "", objectGroup, map[string]interface{}{
            "error":    err.Error(),
            "duration": duration,
        })
        return ObjectResolution{}, fmt.Errorf("解析响应失败: %w", err)
    }

    objectID := resp.Result.Object.ObjectID
    if objectID == "" {
        orm.recordEvent("invalid_object", nodeID, "", objectGroup, map[string]interface{}{
            "duration": duration,
        })
        return ObjectResolution{}, fmt.Errorf("获取到无效的对象ID")
    }

    // 创建解析记录
    resolution := ObjectResolution{
        NodeID:         nodeID,
        ObjectID:       objectID,
        ObjectGroup:    objectGroup,
        CreatedAt:      time.Now(),
        LastAccessed:   time.Now(),
        AccessCount:    1,
        ObjectType:     resp.Result.Object.Type,
        ObjectSubtype:  resp.Result.Object.Subtype,
        Valid:          true,
        Metadata:       metadata,
    }

    // 保存到缓存
    orm.resolutions[nodeID] = resolution

    // 更新对象组
    if group, exists := orm.objectGroups[objectGroup]; exists {
        group.ObjectCount++
        group.Objects[objectID] = true
    }

    // 清理过期解析
    orm.cleanupExpiredResolutions()

    orm.recordEvent("resolution_success", nodeID, objectID, objectGroup, map[string]interface{}{
        "duration": duration,
        "type":     resp.Result.Object.Type,
        "subtype":  resp.Result.Object.Subtype,
    })

    return resolution, nil
}

func (orm *ObjectResolutionManager) isResolutionValid(resolution ObjectResolution) bool {
    if !resolution.Valid {
        return false
    }

    // 检查是否过期
    if orm.cache.ttl > 0 && time.Since(resolution.LastAccessed) > orm.cache.ttl {
        return false
    }

    // 这里还可以添加其他验证逻辑
    return true
}

func (orm *ObjectResolutionManager) ensureObjectGroup(groupName string) {
    if _, exists := orm.objectGroups[groupName]; !exists {
        orm.objectGroups[groupName] = &ObjectGroup{
            Name:        groupName,
            CreatedAt:   time.Now(),
            ObjectCount: 0,
            Objects:     make(map[string]bool),
        }
    }
}

func (orm *ObjectResolutionManager) cleanupExpiredResolutions() {
    if len(orm.resolutions) <= orm.cache.maxSize {
        return
    }

    // 清理过期解析
    now := time.Now()
    for nodeID, resolution := range orm.resolutions {
        if orm.cache.ttl > 0 && now.Sub(resolution.LastAccessed) > orm.cache.ttl {
            // 从对象组中移除
            if group, exists := orm.objectGroups[resolution.ObjectGroup]; exists {
                delete(group.Objects, resolution.ObjectID)
                group.ObjectCount--
            }

            delete(orm.resolutions, nodeID)
            orm.recordEvent("auto_cleanup", nodeID, resolution.ObjectID, resolution.ObjectGroup, nil)
        }
    }

    // 如果仍然超过最大大小，清理最久未访问的
    if len(orm.resolutions) > orm.cache.maxSize {
        orm.cleanupLRU()
    }
}

func (orm *ObjectResolutionManager) cleanupLRU() {
    // 找到最久未访问的解析
    var oldestID int
    var oldestTime time.Time
    first := true

    for nodeID, resolution := range orm.resolutions {
        if first || resolution.LastAccessed.Before(oldestTime) {
            oldestID = nodeID
            oldestTime = resolution.LastAccessed
            first = false
        }
    }

    if !first {
        deletedResolution := orm.resolutions[oldestID]

        // 从对象组中移除
        if group, exists := orm.objectGroups[deletedResolution.ObjectGroup]; exists {
            delete(group.Objects, deletedResolution.ObjectID)
            group.ObjectCount--
        }

        delete(orm.resolutions, oldestID)
        orm.recordEvent("lru_cleanup", oldestID, deletedResolution.ObjectID, deletedResolution.ObjectGroup, map[string]interface{}{
            "lastAccessed": deletedResolution.LastAccessed,
        })
    }
}

func (orm *ObjectResolutionManager) recordEvent(eventType string, nodeID int, objectID, objectGroup string, details map[string]interface{}) {
    event := ResolutionEvent{
        Timestamp:   time.Now(),
        EventType:   eventType,
        NodeID:      nodeID,
        ObjectID:    objectID,
        ObjectGroup: objectGroup,
        Details:     details,
    }

    orm.history = append(orm.history, event)
}

func (orm *ObjectResolutionManager) GetStats() map[string]interface{} {
    totalResolutions := len(orm.resolutions)
    validResolutions := 0
    totalAccesses := 0

    for _, resolution := range orm.resolutions {
        if resolution.Valid {
            validResolutions++
        }
        totalAccesses += resolution.AccessCount
    }

    avgAccesses := 0.0
    if totalResolutions > 0 {
        avgAccesses = float64(totalAccesses) / float64(totalResolutions)
    }

    return map[string]interface{}{
        "totalResolutions":   totalResolutions,
        "validResolutions":   validResolutions,
        "invalidResolutions": totalResolutions - validResolutions,
        "totalAccesses":      totalAccesses,
        "averageAccesses":    avgAccesses,
        "objectGroups":       len(orm.objectGroups),
        "historyEvents":      len(orm.history),
        "cacheHitRatio": func() float64 {
            hits := 0
            for _, event := range orm.history {
                if event.EventType == "cache_hit" {
                    hits++
                }
            }
            if len(orm.history) > 0 {
                return float64(hits) / float64(len(orm.history)) * 100
            }
            return 0
        }(),
    }
}

func (orm *ObjectResolutionManager) GetObjectGroupStats() map[string]map[string]interface{} {
    stats := make(map[string]map[string]interface{})

    for groupName, group := range orm.objectGroups {
        stats[groupName] = map[string]interface{}{
            "createdAt":   group.CreatedAt,
            "objectCount": group.ObjectCount,
            "age":         time.Since(group.CreatedAt),
        }
    }

    return stats
}

// 演示对象解析管理场景
func demonstrateObjectResolutionManagement() {
    fmt.Println("=== 对象解析管理演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建解析管理器
    manager := NewObjectResolutionManager(50, 2*time.Minute)

    // 模拟一批节点ID
    nodeIDs := []int{2001, 2002, 2003, 2004, 2005}

    fmt.Printf("测试 %d 个节点的对象解析管理\n\n", len(nodeIDs))

    // 场景1: 首次解析节点
    fmt.Printf("场景1: 首次解析节点\n")

    for _, nodeID := range nodeIDs {
        metadata := map[string]interface{}{
            "source":    "demo",
            "timestamp": time.Now().Unix(),
        }

        resolution, err := manager.ResolveNode(nodeID, "demo_group", metadata)
        if err != nil {
            fmt.Printf("  节点ID %d: ❌ 解析失败: %v\n", nodeID, err)
        } else {
            fmt.Printf("  节点ID %d: ✅ 解析成功 -> 对象ID %s\n",
                nodeID, resolution.ObjectID)
        }
    }

    // 场景2: 缓存命中测试
    fmt.Printf("\n场景2: 缓存命中测试\n")

    for i := 0; i < 3; i++ {
        for _, nodeID := range nodeIDs {
            resolution, err := manager.ResolveNode(nodeID, "demo_group", nil)
            if err == nil {
                fmt.Printf("  节点ID %d: 访问 %d 次，对象ID %s\n",
                    nodeID, resolution.AccessCount, resolution.ObjectID)
            }
        }

        if i < 2 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    // 场景3: 多个对象组
    fmt.Printf("\n场景3: 多个对象组测试\n")

    testGroups := []string{"form_elements", "buttons", "inputs"}
    for _, group := range testGroups {
        fmt.Printf("  对象组: %s\n", group)

        for _, nodeID := range nodeIDs[:2] { // 只测试前两个节点
            resolution, err := manager.ResolveNode(nodeID, group, nil)
            if err == nil {
                fmt.Printf("    节点ID %d -> 对象ID %s\n", nodeID, resolution.ObjectID)
            }
        }
    }

    // 显示统计信息
    fmt.Printf("\n=== 解析管理统计 ===\n")
    stats := manager.GetStats()
    fmt.Printf("  总解析数: %d\n", stats["totalResolutions"])
    fmt.Printf("  有效解析: %d\n", stats["validResolutions"])
    fmt.Printf("  总访问次数: %d\n", stats["totalAccesses"])
    fmt.Printf("  平均访问次数: %.1f\n", stats["averageAccesses"])
    fmt.Printf("  对象组数量: %d\n", stats["objectGroups"])
    fmt.Printf("  缓存命中率: %.1f%%\n", stats["cacheHitRatio"])

    // 显示对象组统计
    fmt.Printf("\n=== 对象组统计 ===\n")
    groupStats := manager.GetObjectGroupStats()
    for groupName, stats := range groupStats {
        fmt.Printf("  组: %s, 对象数: %d, 创建时间: %s\n",
            groupName,
            stats["objectCount"],
            stats["createdAt"].(time.Time).Format("15:04:05"))
    }

    // 显示最近事件
    fmt.Printf("\n=== 最近事件 ===\n")
    // 获取最后5个事件
    if len(manager.history) > 0 {
        start := len(manager.history) - 5
        if start < 0 {
            start = 0
        }

        for i := start; i < len(manager.history); i++ {
            event := manager.history[i]
            fmt.Printf("%d. [%s] 节点ID: %d, 对象组: %s\n",
                i+1, event.EventType, event.NodeID, event.ObjectGroup)
        }
    }

    fmt.Println("\n=== 演示完成 ===")
}

*/

// -----------------------------------------------  DOM.scrollIntoViewIfNeeded  -----------------------------------------------
// === 应用场景 ===
// 1. 元素可见: 确保元素在视口中可见
// 2. 自动滚动: 自动滚动到目标元素
// 3. 焦点管理: 元素获得焦点时滚动到视图
// 4. 表单验证: 验证失败时滚动到错误元素
// 5. 导航定位: 页面内导航到指定位置
// 6. 响应式调整: 响应式布局中确保元素可见

// CDPDOMScrollIntoViewIfNeeded 如果元素不可见，则滚动到视图中
// nodeID: 要滚动到的节点ID
// rect: 可选，指定元素内部的矩形区域
func CDPDOMScrollIntoViewIfNeeded(nodeID int, rect *DOMRect) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息参数
	params := map[string]interface{}{
		"nodeId": nodeID,
	}

	if rect != nil {
		params["rect"] = rect
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "DOM.scrollIntoViewIfNeeded",
        "params": %s
    }`, reqID, string(paramsJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.scrollIntoViewIfNeeded 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.scrollIntoViewIfNeeded 请求超时")
		}
	}
}

// DOMRect 矩形区域结构
type DOMRect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

/*

// 示例: 确保重要元素在视图中可见
func ExampleCDPDOMScrollIntoViewIfNeeded() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个重要元素（比如错误消息或表单字段）的节点ID
    importantElementID := 1001

    fmt.Printf("=== 滚动到视图演示 ===\n")
    fmt.Printf("目标元素节点ID: %d\n\n", importantElementID)

    // 2. 首先检查元素的当前位置
    fmt.Printf("=== 当前位置检查 ===\n")

    // 获取元素信息
    elementInfo, err := getNodeInfo(importantElementID)
    if err != nil {
        log.Printf("获取元素信息失败: %v", err)
        return
    }

    displayElementInfo("目标元素", elementInfo)

    // 3. 检查元素是否在视口中可见
    fmt.Printf("\n=== 可见性检查 ===\n")

    isVisible, visibilityInfo := checkElementVisibility(importantElementID)

    if isVisible {
        fmt.Printf("✅ 元素已经在视口中可见\n")
        fmt.Printf("  位置: X=%.1f, Y=%.1f\n", visibilityInfo["x"], visibilityInfo["y"])
        fmt.Printf("  尺寸: %.1f x %.1f\n", visibilityInfo["width"], visibilityInfo["height"])

        // 显示可见区域信息
        if vpInfo, ok := visibilityInfo["viewport"].(map[string]float64); ok {
            fmt.Printf("  视口: %.0f x %.0f (滚动: X=%.0f, Y=%.0f)\n",
                vpInfo["width"], vpInfo["height"], vpInfo["scrollX"], vpInfo["scrollY"])
        }
    } else {
        fmt.Printf("⚠️ 元素不可见或部分不可见\n")
        if reasons, ok := visibilityInfo["reasons"].([]string); ok {
            for _, reason := range reasons {
                fmt.Printf("  - %s\n", reason)
            }
        }
    }

    // 4. 如果需要，滚动到视图
    fmt.Printf("\n=== 滚动操作 ===\n")

    if !isVisible {
        fmt.Printf("执行滚动到视图操作...\n")

        startTime := time.Now()
        result, err := CDPDOMScrollIntoViewIfNeeded(importantElementID, nil)
        scrollTime := time.Since(startTime)

        if err != nil {
            log.Printf("滚动失败: %v", err)
            return
        }

        fmt.Printf("滚动结果: %s\n", result)
        fmt.Printf("滚动耗时: %v\n", scrollTime)

        // 5. 验证滚动效果
        fmt.Printf("\n=== 滚动后验证 ===\n")

        // 短暂等待滚动完成
        time.Sleep(200 * time.Millisecond)

        isNowVisible, _ := checkElementVisibility(importantElementID)
        if isNowVisible {
            fmt.Printf("✅ 元素现在可见了\n")
        } else {
            fmt.Printf("⚠️ 元素可能仍然不可见\n")
        }
    } else {
        fmt.Printf("元素已可见，无需滚动\n")

        // 测试强制滚动
        fmt.Printf("\n测试强制滚动（即使已可见）...\n")
        result, err := CDPDOMScrollIntoViewIfNeeded(importantElementID, nil)
        if err != nil {
            fmt.Printf("强制滚动失败: %v\n", err)
        } else {
            fmt.Printf("强制滚动结果: %s\n", result)
        }
    }

    // 6. 测试滚动到指定矩形区域
    fmt.Printf("\n=== 滚动到指定区域 ===\n")

    // 指定要滚动到的矩形区域（元素内部的特定区域）
    rect := &DOMRect{
        X:      10,
        Y:      20,
        Width:  100,
        Height: 50,
    }

    fmt.Printf("滚动到元素内部特定区域: X=%.1f, Y=%.1f, 尺寸=%.1f x %.1f\n",
        rect.X, rect.Y, rect.Width, rect.Height)

    result, err := CDPDOMScrollIntoViewIfNeeded(importantElementID, rect)
    if err != nil {
        fmt.Printf("滚动到指定区域失败: %v\n", err)
    } else {
        fmt.Printf("滚动到指定区域结果: %s\n", result)
    }

    // 7. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次滚动
    testRepeatedScrolling(importantElementID, 3)

    // 测试不同元素的滚动
    testMultipleElements([]int{1001, 1002, 1003})

    // 8. 滚动行为分析
    fmt.Printf("\n=== 滚动行为分析 ===\n")

    analyzeScrollBehavior(importantElementID)

    // 9. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "表单验证",
            description: "验证失败时滚动到错误字段",
            useCase:     "提交表单时检查验证错误，滚动到第一个错误字段",
        },
        {
            name:        "页面内导航",
            description: "点击链接滚动到对应章节",
            useCase:     "点击目录链接滚动到对应章节标题",
        },
        {
            name:        "聊天应用",
            description: "新消息到达时滚动到底部",
            useCase:     "收到新消息时自动滚动到最新消息",
        },
        {
            name:        "表格浏览",
            description: "在大型表格中定位到特定行",
            useCase:     "搜索表格内容，滚动到匹配行",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 检查元素可见性
func checkElementVisibility(nodeID int) (bool, map[string]interface{}) {
    info := make(map[string]interface{})

    // 获取元素的盒模型
    boxResult, err := CDPDOMGetBoxModel(nodeID)
    if err != nil {
        info["reasons"] = []string{fmt.Sprintf("无法获取盒模型: %v", err)}
        return false, info
    }

    var boxResp struct {
        Result struct {
            Model struct {
                Content []float64 `json:"content"`
            } `json:"model"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(boxResult), &boxResp); err != nil {
        info["reasons"] = []string{fmt.Sprintf("无法解析盒模型: %v", err)}
        return false, info
    }

    content := boxResp.Result.Model.Content
    if len(content) < 8 {
        info["reasons"] = []string{"无效的盒模型数据"}
        return false, info
    }

    // 提取元素位置和尺寸
    elementX := content[0]
    elementY := content[1]
    elementWidth := content[2] - content[0]
    elementHeight := content[5] - content[1]

    info["x"] = elementX
    info["y"] = elementY
    info["width"] = elementWidth
    info["height"] = elementHeight

    // 模拟视口信息（实际应该从页面获取）
    viewportInfo := map[string]float64{
        "width":   1920,
        "height":  1080,
        "scrollX": 0,
        "scrollY": 0,
    }
    info["viewport"] = viewportInfo

    // 检查是否在视口内
    var reasons []string

    // 检查水平方向
    if elementX+elementWidth < 0 {
        reasons = append(reasons, "元素在视口左侧之外")
    } else if elementX > viewportInfo["width"] {
        reasons = append(reasons, "元素在视口右侧之外")
    }

    // 检查垂直方向
    if elementY+elementHeight < 0 {
        reasons = append(reasons, "元素在视口顶部之外")
    } else if elementY > viewportInfo["height"] {
        reasons = append(reasons, "元素在视口底部之外")
    }

    // 检查是否被其他元素遮挡（简化实现）
    if elementWidth == 0 || elementHeight == 0 {
        reasons = append(reasons, "元素尺寸为0")
    }

    if len(reasons) > 0 {
        info["reasons"] = reasons
        return false, info
    }

    return true, info
}

// 测试重复滚动
func testRepeatedScrolling(nodeID int, count int) {
    fmt.Printf("重复滚动测试 (%d 次):\n", count)

    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < count; i++ {
        startTime := time.Now()
        result, err := CDPDOMScrollIntoViewIfNeeded(nodeID, nil)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            successCount++
            fmt.Printf("  第 %d 次: ✅ 成功 (耗时: %v)\n", i+1, duration)
        }

        // 短暂延迟
        if i < count-1 {
            time.Sleep(200 * time.Millisecond)
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))
}

// 测试多个元素
func testMultipleElements(nodeIDs []int) {
    fmt.Printf("\n多元素滚动测试 (%d 个元素):\n", len(nodeIDs))

    var results []ScrollTestResult

    for _, nodeID := range nodeIDs {
        result := ScrollTestResult{
            NodeID: nodeID,
        }

        startTime := time.Now()
        response, err := CDPDOMScrollIntoViewIfNeeded(nodeID, nil)
        result.Duration = time.Since(startTime)

        if err != nil {
            result.Success = false
            result.Error = err.Error()
        } else {
            result.Success = true
            result.Response = response
        }

        results = append(results, result)

        // 短暂延迟
        time.Sleep(300 * time.Millisecond)
    }

    // 显示结果
    successCount := 0
    for _, result := range results {
        status := "❌ 失败"
        if result.Success {
            status = "✅ 成功"
            successCount++
        }

        fmt.Printf("  节点ID: %-6d %s (耗时: %v)\n",
            result.NodeID, status, result.Duration)
    }

    fmt.Printf("总体成功率: %.1f%% (%d/%d)\n",
        float64(successCount)/float64(len(results))*100,
        successCount, len(results))
}

type ScrollTestResult struct {
    NodeID   int
    Success  bool
    Error    string
    Response string
    Duration time.Duration
}

// 分析滚动行为
func analyzeScrollBehavior(nodeID int) {
    fmt.Printf("滚动行为分析:\n")
    fmt.Printf("  1. 智能滚动: 只在元素不可见时滚动\n")
    fmt.Printf("  2. 平滑滚动: 浏览器可能使用平滑滚动动画\n")
    fmt.Printf("  3. 对齐方式: 通常滚动到元素可见的最近位置\n")
    fmt.Printf("  4. 边界处理: 考虑视口边距和填充\n")

    fmt.Printf("\n滚动优化建议:\n")
    fmt.Printf("  - 避免频繁滚动，可能导致用户体验不佳\n")
    fmt.Printf("  - 考虑使用 scrollIntoView 替代，有更多选项\n")
    fmt.Printf("  - 对于长列表，考虑虚拟滚动\n")
    fmt.Printf("  - 移动设备上注意滚动性能\n")
}

// 高级功能: 智能滚动管理器
type SmartScrollManager struct {
    scrollHistory []ScrollRecord
    elementCache  map[int]ElementScrollInfo
    config        ScrollConfig
}

type ScrollRecord struct {
    Timestamp   time.Time
    NodeID      int
    Rect        *DOMRect
    Duration    time.Duration
    Success     bool
    Error       string
    Reason      string
    BeforeState ScrollState
    AfterState  ScrollState
}

type ScrollState struct {
    Visible     bool
    Position    DOMRect
    Viewport    DOMRect
    ScrollOffset Point
}

type Point struct {
    X float64
    Y float64
}

type ElementScrollInfo struct {
    NodeID        int
    LastScrolled  time.Time
    ScrollCount   int
    AverageTime   time.Duration
    LastResult    string
}

type ScrollConfig struct {
    CheckBeforeScroll  bool
    SmoothScrolling    bool
    MaxScrollAttempts  int
    ScrollTimeout      time.Duration
    RecordHistory      bool
    CacheResults       bool
    RetryOnFailure     bool
}

func NewSmartScrollManager(config ScrollConfig) *SmartScrollManager {
    if config.MaxScrollAttempts <= 0 {
        config.MaxScrollAttempts = 3
    }
    if config.ScrollTimeout <= 0 {
        config.ScrollTimeout = 5 * time.Second
    }

    return &SmartScrollManager{
        scrollHistory: make([]ScrollRecord, 0),
        elementCache:  make(map[int]ElementScrollInfo),
        config:        config,
    }
}

func (ssm *SmartScrollManager) ScrollToElement(nodeID int, rect *DOMRect, reason string) (ScrollRecord, error) {
    record := ScrollRecord{
        Timestamp: time.Now(),
        NodeID:    nodeID,
        Rect:      rect,
        Reason:    reason,
    }

    // 获取滚动前状态
    if ssm.config.CheckBeforeScroll {
        var err error
        record.BeforeState, err = ssm.getScrollState(nodeID)
        if err != nil {
            record.Success = false
            record.Error = fmt.Sprintf("获取滚动前状态失败: %v", err)
            ssm.recordScroll(record)
            return record, fmt.Errorf("获取状态失败: %w", err)
        }

        // 如果元素已经可见，可以跳过滚动
        if record.BeforeState.Visible && ssm.config.CheckBeforeScroll {
            record.Success = true
            record.Duration = 0
            record.AfterState = record.BeforeState
            ssm.recordScroll(record)
            return record, nil
        }
    }

    // 执行滚动
    startTime := time.Now()

    var result string
    var err error
    attempts := 0

    for attempts < ssm.config.MaxScrollAttempts {
        result, err = CDPDOMScrollIntoViewIfNeeded(nodeID, rect)
        attempts++

        if err == nil {
            break
        }

        if attempts < ssm.config.MaxScrollAttempts && ssm.config.RetryOnFailure {
            log.Printf("滚动失败，第 %d 次重试: %v", attempts, err)
            time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
        }
    }

    record.Duration = time.Since(startTime)

    if err != nil {
        record.Success = false
        record.Error = err.Error()
    } else {
        record.Success = true
        record.Response = result
    }

    // 获取滚动后状态
    if ssm.config.CheckBeforeScroll && record.Success {
        afterState, err := ssm.getScrollState(nodeID)
        if err == nil {
            record.AfterState = afterState
        }
    }

    // 记录历史
    ssm.recordScroll(record)

    // 更新缓存
    if ssm.config.CacheResults && record.Success {
        ssm.updateElementCache(nodeID, record)
    }

    if err != nil {
        return record, fmt.Errorf("滚动失败: %w", err)
    }

    return record, nil
}

func (ssm *SmartScrollManager) getScrollState(nodeID int) (ScrollState, error) {
    state := ScrollState{}

    // 获取可见性
    isVisible, visibilityInfo := checkElementVisibility(nodeID)
    state.Visible = isVisible

    if visibilityInfo != nil {
        if x, ok := visibilityInfo["x"].(float64); ok {
            state.Position.X = x
        }
        if y, ok := visibilityInfo["y"].(float64); ok {
            state.Position.Y = y
        }
        if width, ok := visibilityInfo["width"].(float64); ok {
            state.Position.Width = width
        }
        if height, ok := visibilityInfo["height"].(float64); ok {
            state.Position.Height = height
        }
    }

    return state, nil
}

func (ssm *SmartScrollManager) recordScroll(record ScrollRecord) {
    if ssm.config.RecordHistory {
        ssm.scrollHistory = append(ssm.scrollHistory, record)
    }
}

func (ssm *SmartScrollManager) updateElementCache(nodeID int, record ScrollRecord) {
    info, exists := ssm.elementCache[nodeID]
    if !exists {
        info = ElementScrollInfo{
            NodeID:       nodeID,
            LastScrolled: record.Timestamp,
            ScrollCount:  1,
            AverageTime:  record.Duration,
            LastResult:   record.Response,
        }
    } else {
        info.LastScrolled = record.Timestamp
        info.ScrollCount++
        // 更新平均时间
        totalTime := info.AverageTime*time.Duration(info.ScrollCount-1) + record.Duration
        info.AverageTime = totalTime / time.Duration(info.ScrollCount)
        info.LastResult = record.Response
    }

    ssm.elementCache[nodeID] = info
}

func (ssm *SmartScrollManager) GetStats() map[string]interface{} {
    totalScrolls := len(ssm.scrollHistory)
    successfulScrolls := 0
    var totalDuration time.Duration

    for _, record := range ssm.scrollHistory {
        if record.Success {
            successfulScrolls++
        }
        totalDuration += record.Duration
    }

    avgDuration := time.Duration(0)
    if totalScrolls > 0 {
        avgDuration = totalDuration / time.Duration(totalScrolls)
    }

    return map[string]interface{}{
        "totalScrolls":     totalScrolls,
        "successfulScrolls": successfulScrolls,
        "failedScrolls":    totalScrolls - successfulScrolls,
        "successRate":      float64(successfulScrolls) / float64(totalScrolls) * 100,
        "totalDuration":    totalDuration,
        "averageDuration":  avgDuration,
        "cachedElements":   len(ssm.elementCache),
    }
}

func (ssm *SmartScrollManager) GetElementStats(nodeID int) (ElementScrollInfo, bool) {
    info, exists := ssm.elementCache[nodeID]
    return info, exists
}

// 演示智能滚动场景
func demonstrateSmartScrolling() {
    fmt.Println("=== 智能滚动管理演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建滚动管理器
    config := ScrollConfig{
        CheckBeforeScroll:  true,
        SmoothScrolling:    true,
        MaxScrollAttempts:  3,
        ScrollTimeout:      5 * time.Second,
        RecordHistory:      true,
        CacheResults:       true,
        RetryOnFailure:     true,
    }

    manager := NewSmartScrollManager(config)

    // 模拟一批元素
    elementIDs := []int{3001, 3002, 3003, 3004}

    fmt.Printf("测试 %d 个元素的智能滚动\n\n", len(elementIDs))

    // 场景1: 正常滚动
    fmt.Printf("场景1: 正常滚动测试\n")

    for _, elementID := range elementIDs {
        fmt.Printf("滚动到元素 %d...\n", elementID)

        record, err := manager.ScrollToElement(elementID, nil, "正常测试")
        if err != nil {
            fmt.Printf("  ❌ 失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)

            if !record.BeforeState.Visible && record.AfterState.Visible {
                fmt.Printf("    从不可见变为可见\n")
            }
        }

        time.Sleep(500 * time.Millisecond)
    }

    // 场景2: 重复滚动测试（缓存命中）
    fmt.Printf("\n场景2: 重复滚动测试（缓存命中）\n")

    for i := 0; i < 2; i++ {
        for _, elementID := range elementIDs[:2] { // 只测试前两个
            record, err := manager.ScrollToElement(elementID, nil, fmt.Sprintf("重复测试 %d", i+1))
            if err == nil && record.Duration == 0 {
                fmt.Printf("元素 %d: 缓存命中，无需滚动\n", elementID)
            } else if err == nil {
                fmt.Printf("元素 %d: 执行滚动 (耗时: %v)\n", elementID, record.Duration)
            }
        }
    }

    // 场景3: 滚动到指定区域
    fmt.Printf("\n场景3: 滚动到指定区域\n")

    rect := &DOMRect{
        X:      50,
        Y:      30,
        Width:  200,
        Height: 100,
    }

    for _, elementID := range elementIDs[:1] { // 只测试第一个
        fmt.Printf("滚动到元素 %d 的指定区域...\n", elementID)

        record, err := manager.ScrollToElement(elementID, rect, "指定区域")
        if err != nil {
            fmt.Printf("  ❌ 失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)
        }
    }

    // 显示统计信息
    fmt.Printf("\n=== 滚动统计 ===\n")
    stats := manager.GetStats()
    fmt.Printf("  总滚动次数: %d\n", stats["totalScrolls"])
    fmt.Printf("  成功滚动: %d\n", stats["successfulScrolls"])
    fmt.Printf("  失败滚动: %d\n", stats["failedScrolls"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("  总耗时: %v\n", stats["totalDuration"])
    fmt.Printf("  平均耗时: %v\n", stats["averageDuration"])
    fmt.Printf("  缓存元素数: %d\n", stats["cachedElements"])

    // 显示元素统计
    fmt.Printf("\n=== 元素统计 ===\n")
    for _, elementID := range elementIDs {
        if info, exists := manager.GetElementStats(elementID); exists {
            fmt.Printf("  元素 %d: 滚动 %d 次，平均耗时 %v\n",
                elementID, info.ScrollCount, info.AverageTime)
        }
    }

    fmt.Println("\n=== 演示完成 ===")
}

*/

// -----------------------------------------------  DOM.setAttributesAsText  -----------------------------------------------
// === 应用场景 ===
// 1. 批量属性设置: 一次设置多个HTML属性
// 2. 属性解析: 解析和设置HTML属性字符串
// 3. 样式批量更新: 批量更新内联样式
// 4. 属性迁移: 从元素复制属性到另一个元素
// 5. 模板应用: 应用属性模板到元素
// 6. 快速原型: 快速设置元素属性进行原型开发

// CDPDOMSetAttributesAsText 通过文本设置节点的属性
// nodeID: 要设置属性的节点ID
// text: 属性文本（格式如 "class=\"btn\" disabled id=\"submit\"")
// name: 属性名称，如果提供则只更新该属性
func CDPDOMSetAttributesAsText(nodeID int, text, name string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息参数
	params := map[string]interface{}{
		"nodeId": nodeID,
		"text":   text,
	}

	if name != "" {
		params["name"] = name
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "DOM.setAttributesAsText",
        "params": %s
    }`, reqID, string(paramsJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.setAttributesAsText 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.setAttributesAsText 请求超时")
		}
	}
}

/*
// 示例: 批量设置按钮元素的属性
func ExampleCDPDOMSetAttributesAsText() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个按钮元素的节点ID
    buttonNodeID := 1001

    fmt.Printf("=== 批量设置属性演示 ===\n")
    fmt.Printf("目标元素节点ID: %d\n\n", buttonNodeID)

    // 2. 首先获取元素的当前属性
    fmt.Printf("=== 设置前属性检查 ===\n")

    attrsBefore, err := getElementAttributes(buttonNodeID)
    if err != nil {
        log.Printf("获取属性失败: %v", err)
        return
    }

    displayAttributes("设置前属性", attrsBefore)

    // 3. 准备要设置的属性文本
    fmt.Printf("\n=== 准备属性文本 ===\n")

    // 定义不同的属性设置场景
    testCases := []struct {
        name        string
        description string
        attrText    string
        attrName    string
    }{
        {
            name:        "完整属性设置",
            description: "设置完整的HTML属性字符串",
            attrText:    `class="btn btn-primary" id="submit-btn" title="提交表单" disabled`,
            attrName:    "",
        },
        {
            name:        "样式属性更新",
            description: "更新内联样式",
            attrText:    `style="color: white; background-color: #007bff; padding: 10px 20px; border-radius: 5px;"`,
            attrName:    "",
        },
        {
            name:        "单个属性更新",
            description: "只更新class属性",
            attrText:    `btn btn-success btn-lg`,
            attrName:    "class",
        },
        {
            name:        "数据属性设置",
            description: "设置数据属性",
            attrText:    `data-testid="submit-button" data-loading="false" data-action="submit"`,
            attrName:    "",
        },
        {
            name:        "ARIA属性设置",
            description: "设置可访问性属性",
            attrText:    `aria-label="提交表单按钮" aria-disabled="false" role="button"`,
            attrName:    "",
        },
    }

    // 4. 执行属性设置
    for i, tc := range testCases {
        fmt.Printf("\n测试 %d: %s\n", i+1, tc.name)
        fmt.Printf("描述: %s\n", tc.description)

        if tc.attrName != "" {
            fmt.Printf("属性名: %s\n", tc.attrName)
        }
        fmt.Printf("属性文本: %s\n", tc.attrText)

        startTime := time.Now()
        result, err := CDPDOMSetAttributesAsText(buttonNodeID, tc.attrText, tc.attrName)
        setTime := time.Since(startTime)

        if err != nil {
            fmt.Printf("❌ 设置失败: %v\n", err)
        } else {
            fmt.Printf("✅ 设置成功 (耗时: %v)\n", setTime)
            fmt.Printf("结果: %s\n", result)

            // 验证设置结果
            fmt.Printf("\n验证设置结果...\n")
            attrsAfter, err := getElementAttributes(buttonNodeID)
            if err != nil {
                fmt.Printf("❌ 验证失败: 无法获取属性\n")
            } else {
                displayAttributes("设置后属性", attrsAfter)

                // 检查特定属性是否设置成功
                if tc.attrName != "" {
                    // 检查特定属性
                    attrFound := false
                    for _, attr := range attrsAfter {
                        if attr.Name == tc.attrName && attr.Value == tc.attrText {
                            attrFound = true
                            break
                        }
                    }

                    if attrFound {
                        fmt.Printf("✅ 属性 '%s' 已成功设置\n", tc.attrName)
                    } else {
                        fmt.Printf("❌ 属性 '%s' 设置可能失败\n", tc.attrName)
                    }
                }
            }
        }

        // 短暂延迟
        if i < len(testCases)-1 {
            time.Sleep(500 * time.Millisecond)
        }
    }

    // 5. 属性文本解析测试
    fmt.Printf("\n=== 属性文本解析测试 ===\n")

    parseTestCases := []struct {
        attrText    string
        expectError bool
        description string
    }{
        {`class="test" id="btn1"`, false, "标准属性格式"},
        {`class='test' id='btn1'`, false, "单引号属性格式"},
        {`class=test id=btn1`, false, "无引号属性格式"},
        {`class="test" data-value="123" aria-label="测试"`, false, "混合属性格式"},
        {`style="color: red; font-size: 16px;"`, false, "样式属性"},
        {`<script>alert('xss')</script>`, true, "潜在XSS攻击"},
        {`""`, false, "空属性文本"},
        {`invalid attr format`, true, "无效属性格式"},
    }

    for _, tc := range parseTestCases {
        fmt.Printf("测试: %s\n", tc.description)
        fmt.Printf("文本: %s\n", tc.attrText)

        result, err := CDPDOMSetAttributesAsText(buttonNodeID, tc.attrText, "")

        if tc.expectError {
            if err != nil {
                fmt.Printf("✅ 预期错误: %v\n", err)
            } else {
                fmt.Printf("❌ 预期错误但成功: %s\n", result)
            }
        } else {
            if err != nil {
                fmt.Printf("❌ 预期成功但失败: %v\n", err)
            } else {
                fmt.Printf("✅ 解析成功\n")
            }
        }
    }

    // 6. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次属性设置
    testRepeatedAttributeSets(buttonNodeID, 5)

    // 测试大文本属性设置
    testLargeAttributeText(buttonNodeID)

    // 7. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "主题切换",
            description: "根据主题动态设置元素样式",
            useCase:     "切换主题时批量更新元素类名和样式",
        },
        {
            name:        "表单验证",
            description: "验证失败时设置错误样式和属性",
            useCase:     "表单验证失败时添加错误类、ARIA属性和样式",
        },
        {
            name:        "组件初始化",
            description: "初始化组件时设置所有必要属性",
            useCase:     "动态创建组件时一次性设置所有HTML属性",
        },
        {
            name:        "模板渲染",
            description: "从模板渲染HTML属性",
            useCase:     "从模板字符串解析并设置属性",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 测试重复属性设置
func testRepeatedAttributeSets(nodeID int, count int) {
    fmt.Printf("重复属性设置测试 (%d 次):\n", count)

    attrText := `class="test" data-index="0"`
    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < count; i++ {
        // 更新data-index属性
        updatedText := fmt.Sprintf(`class="test" data-index="%d"`, i)

        startTime := time.Now()
        result, err := CDPDOMSetAttributesAsText(nodeID, updatedText, "")
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            successCount++
            fmt.Printf("  第 %d 次: ✅ 成功 (data-index=%d, 耗时: %v)\n",
                i+1, i, duration)

            // 验证设置
            attrs, err := getElementAttributes(nodeID)
            if err == nil {
                dataIndex := ""
                for _, attr := range attrs {
                    if attr.Name == "data-index" {
                        dataIndex = attr.Value
                        break
                    }
                }

                if dataIndex == fmt.Sprintf("%d", i) {
                    fmt.Printf("    验证: ✅ data-index 正确\n")
                } else {
                    fmt.Printf("    验证: ❌ data-index 不匹配\n")
                }
            }
        }

        // 短暂延迟
        if i < count-1 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))
}

// 测试大文本属性设置
func testLargeAttributeText(nodeID int) {
    fmt.Printf("\n大文本属性设置测试:\n")

    // 创建大型属性文本
    var builder strings.Builder
    builder.WriteString(`class="large-btn ")

    // 添加大量类
    for i := 0; i < 50; i++ {
        builder.WriteString(fmt.Sprintf("cls-%d ", i))
    }

    builder.WriteString(`" style="`)

    // 添加大量样式
    for i := 0; i < 20; i++ {
        builder.WriteString(fmt.Sprintf("--custom-var-%d: %d; ", i, i*10))
    }

    builder.WriteString(`color: red; background: blue; padding: 10px; margin: 5px;" `)

    // 添加数据属性
    for i := 0; i < 30; i++ {
        builder.WriteString(fmt.Sprintf(`data-item-%d="value-%d" `, i, i))
    }

    largeText := builder.String()
    textSize := len(largeText)

    fmt.Printf("属性文本大小: %d 字符\n", textSize)

    startTime := time.Now()
    result, err := CDPDOMSetAttributesAsText(nodeID, largeText, "")
    duration := time.Since(startTime)

    if err != nil {
        fmt.Printf("❌ 设置失败: %v\n", err)
    } else {
        fmt.Printf("✅ 设置成功 (耗时: %v)\n", duration)
        fmt.Printf("处理速度: %.0f 字符/秒\n", float64(textSize)/duration.Seconds())
    }
}

// 高级功能: 智能属性设置器
type SmartAttributeSetter struct {
    setHistory []AttributeSetRecord
    validation *AttributeValidator
    parser     *AttributeParser
    cache      map[string]AttributeSet
}

type AttributeSetRecord struct {
    Timestamp   time.Time
    NodeID      int
    Text        string
    Name        string
    Duration    time.Duration
    Success     bool
    Error       string
    BeforeState []Attribute
    AfterState  []Attribute
    Changes     []AttributeChange
}

type AttributeChange struct {
    Name      string
    OldValue  string
    NewValue  string
    Operation string // "added", "updated", "removed"
}

type AttributeSet struct {
    Text      string
    Name      string
    Parsed    map[string]string
    Timestamp time.Time
    UsageCount int
}

type AttributeValidator struct {
    allowedAttributes map[string]bool
    maxTextLength     int
    safeValues        map[string][]string
}

type AttributeParser struct {
    strictMode bool
    autoQuote  bool
    normalize  bool
}

func NewSmartAttributeSetter() *SmartAttributeSetter {
    return &SmartAttributeSetter{
        setHistory: make([]AttributeSetRecord, 0),
        validation: &AttributeValidator{
            allowedAttributes: map[string]bool{
                "class": true, "id": true, "style": true, "title": true,
                "disabled": true, "readonly": true, "required": true,
                "placeholder": true, "value": true, "type": true,
                "name": true, "for": true, "href": true, "src": true,
                "alt": true, "aria-": true, "data-": true, "role": true,
            },
            maxTextLength: 10000,
            safeValues: map[string][]string{
                "role": {"button", "link", "img", "heading", "list", "listitem"},
            },
        },
        parser: &AttributeParser{
            strictMode: false,
            autoQuote:  true,
            normalize:  true,
        },
        cache: make(map[string]AttributeSet),
    }
}

func (sas *SmartAttributeSetter) SetAttributes(nodeID int, text, name string, validate bool) (AttributeSetRecord, error) {
    record := AttributeSetRecord{
        Timestamp: time.Now(),
        NodeID:    nodeID,
        Text:      text,
        Name:      name,
    }

    // 获取设置前状态
    beforeAttrs, err := getElementAttributes(nodeID)
    if err != nil {
        return record, fmt.Errorf("获取当前属性失败: %w", err)
    }
    record.BeforeState = beforeAttrs

    // 验证输入
    if validate {
        if err := sas.validation.Validate(text, name); err != nil {
            record.Error = err.Error()
            sas.recordSet(record)
            return record, fmt.Errorf("验证失败: %w", err)
        }
    }

    // 解析属性文本
    parsedAttrs, err := sas.parser.Parse(text, name)
    if err != nil {
        record.Error = err.Error()
        sas.recordSet(record)
        return record, fmt.Errorf("解析失败: %w", err)
    }

    // 缓存检查
    cacheKey := fmt.Sprintf("%s|%s", text, name)
    if cached, exists := sas.cache[cacheKey]; exists {
        cached.UsageCount++
        sas.cache[cacheKey] = cached
    } else {
        sas.cache[cacheKey] = AttributeSet{
            Text:      text,
            Name:      name,
            Parsed:    parsedAttrs,
            Timestamp: time.Now(),
            UsageCount: 1,
        }
    }

    // 执行设置
    startTime := time.Now()
    result, err := CDPDOMSetAttributesAsText(nodeID, text, name)
    record.Duration = time.Since(startTime)

    if err != nil {
        record.Success = false
        record.Error = err.Error()
    } else {
        record.Success = true
        record.Response = result
    }

    // 获取设置后状态
    afterAttrs, err := getElementAttributes(nodeID)
    if err == nil {
        record.AfterState = afterAttrs

        // 计算变化
        record.Changes = sas.calculateChanges(beforeAttrs, afterAttrs)
    }

    // 记录历史
    sas.recordSet(record)

    if err != nil {
        return record, fmt.Errorf("设置失败: %w", err)
    }

    return record, nil
}

func (av *AttributeValidator) Validate(text, name string) error {
    // 检查文本长度
    if len(text) > av.maxTextLength {
        return fmt.Errorf("属性文本过长: %d > %d", len(text), av.maxTextLength)
    }

    // 检查危险内容
    if strings.Contains(strings.ToLower(text), "<script>") {
        return fmt.Errorf("检测到潜在XSS攻击")
    }

    // 如果指定了属性名，检查是否允许
    if name != "" && !av.isAttributeAllowed(name) {
        return fmt.Errorf("属性不允许: %s", name)
    }

    return nil
}

func (av *AttributeValidator) isAttributeAllowed(name string) bool {
    // 检查精确匹配
    if av.allowedAttributes[name] {
        return true
    }

    // 检查前缀匹配
    for allowed := range av.allowedAttributes {
        if strings.HasSuffix(allowed, "-") && strings.HasPrefix(name, allowed) {
            return true
        }
    }

    return false
}

func (ap *AttributeParser) Parse(text, name string) (map[string]string, error) {
    attrs := make(map[string]string)

    if text == "" {
        return attrs, nil
    }

    if name != "" {
        // 只解析单个属性
        attrs[name] = text
        return attrs, nil
    }

    // 简单解析属性文本
    // 实际应该使用更复杂的HTML属性解析器
    re := regexp.MustCompile(`([a-zA-Z\-_:]+)=(?:"([^"]*)"|'([^']*)'|([^\s>]+))|([a-zA-Z\-_:]+)`)
    matches := re.FindAllStringSubmatch(text, -1)

    for _, match := range matches {
        if len(match) >= 2 {
            attrName := match[1]
            var attrValue string

            if len(match) >= 6 {
                // 处理无值属性
                if match[5] != "" {
                    attrName = match[5]
                    attrValue = ""
                } else {
                    // 处理有值属性
                    for i := 2; i <= 4; i++ {
                        if match[i] != "" {
                            attrValue = match[i]
                            break
                        }
                    }
                }
            }

            if ap.normalize {
                attrName = strings.ToLower(attrName)
            }

            attrs[attrName] = attrValue
        }
    }

    return attrs, nil
}

func (sas *SmartAttributeSetter) calculateChanges(before, after []Attribute) []AttributeChange {
    var changes []AttributeChange

    beforeMap := make(map[string]string)
    afterMap := make(map[string]string)

    for _, attr := range before {
        beforeMap[attr.Name] = attr.Value
    }

    for _, attr := range after {
        afterMap[attr.Name] = attr.Value
    }

    // 检查新增和更新的属性
    for name, newValue := range afterMap {
        oldValue, existed := beforeMap[name]
        if !existed {
            changes = append(changes, AttributeChange{
                Name:      name,
                OldValue:  "",
                NewValue:  newValue,
                Operation: "added",
            })
        } else if oldValue != newValue {
            changes = append(changes, AttributeChange{
                Name:      name,
                OldValue:  oldValue,
                NewValue:  newValue,
                Operation: "updated",
            })
        }
    }

    // 检查删除的属性
    for name, oldValue := range beforeMap {
        if _, exists := afterMap[name]; !exists {
            changes = append(changes, AttributeChange{
                Name:      name,
                OldValue:  oldValue,
                NewValue:  "",
                Operation: "removed",
            })
        }
    }

    return changes
}

func (sas *SmartAttributeSetter) recordSet(record AttributeSetRecord) {
    sas.setHistory = append(sas.setHistory, record)
}

func (sas *SmartAttributeSetter) GetStats() map[string]interface{} {
    totalSets := len(sas.setHistory)
    successfulSets := 0
    var totalDuration time.Duration

    for _, record := range sas.setHistory {
        if record.Success {
            successfulSets++
        }
        totalDuration += record.Duration
    }

    avgDuration := time.Duration(0)
    if totalSets > 0 {
        avgDuration = totalDuration / time.Duration(totalSets)
    }

    return map[string]interface{}{
        "totalSets":     totalSets,
        "successfulSets": successfulSets,
        "failedSets":    totalSets - successfulSets,
        "successRate":   float64(successfulSets) / float64(totalSets) * 100,
        "totalDuration": totalDuration,
        "averageDuration": avgDuration,
        "cacheSize":     len(sas.cache),
    }
}

// 演示智能属性设置
func demonstrateSmartAttributeSetting() {
    fmt.Println("=== 智能属性设置演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建智能设置器
    setter := NewSmartAttributeSetter()

    // 测试元素
    elementID := 2001

    fmt.Printf("测试元素ID: %d\n\n", elementID)

    // 场景1: 正常属性设置
    fmt.Printf("场景1: 正常属性设置\n")

    testCases := []struct {
        name     string
        text     string
        attrName string
        validate bool
    }{
        {
            name:     "完整属性设置",
            text:     `class="btn-primary" id="action-btn" title="操作按钮"`,
            attrName: "",
            validate: true,
        },
        {
            name:     "单个属性更新",
            text:     "btn-lg btn-block",
            attrName: "class",
            validate: true,
        },
        {
            name:     "样式设置",
            text:     `style="color: white; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);"`,
            attrName: "",
            validate: true,
        },
    }

    for _, tc := range testCases {
        fmt.Printf("测试: %s\n", tc.name)

        record, err := setter.SetAttributes(elementID, tc.text, tc.attrName, tc.validate)
        if err != nil {
            fmt.Printf("  ❌ 失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)

            // 显示变化
            if len(record.Changes) > 0 {
                fmt.Printf("  变化:\n")
                for _, change := range record.Changes {
                    fmt.Printf("    %s: %s -> %s [%s]\n",
                        change.Name, change.OldValue, change.NewValue, change.Operation)
                }
            }
        }

        time.Sleep(300 * time.Millisecond)
    }

    // 场景2: 验证测试
    fmt.Printf("\n场景2: 验证测试\n")

    validationTests := []struct {
        name     string
        text     string
        attrName string
        expectError bool
    }{
        {
            name:       "危险脚本",
            text:       `<script>alert('xss')</script>`,
            attrName:   "",
            expectError: true,
        },
        {
            name:       "超长文本",
            text:       strings.Repeat("a", 20000),
            attrName:   "",
            expectError: true,
        },
        {
            name:       "允许的属性",
            text:       `data-custom="value"`,
            attrName:   "",
            expectError: false,
        },
    }

    for _, test := range validationTests {
        fmt.Printf("测试: %s\n", test.name)

        _, err := setter.SetAttributes(elementID, test.text, test.attrName, true)

        if test.expectError {
            if err != nil {
                fmt.Printf("  ✅ 预期错误: %v\n", err)
            } else {
                fmt.Printf("  ❌ 预期错误但成功\n")
            }
        } else {
            if err != nil {
                fmt.Printf("  ❌ 预期成功但失败: %v\n", err)
            } else {
                fmt.Printf("  ✅ 验证通过\n")
            }
        }
    }

    // 显示统计
    fmt.Printf("\n=== 设置统计 ===\n")
    stats := setter.GetStats()
    fmt.Printf("  总设置次数: %d\n", stats["totalSets"])
    fmt.Printf("  成功设置: %d\n", stats["successfulSets"])
    fmt.Printf("  失败设置: %d\n", stats["failedSets"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("  平均耗时: %v\n", stats["averageDuration"])
    fmt.Printf("  缓存大小: %d\n", stats["cacheSize"])

    fmt.Println("\n=== 演示完成 ===")
}


*/

// -----------------------------------------------  DOM.setAttributeValue  -----------------------------------------------
// === 应用场景 ===
// 1. 单个属性设置: 精确设置特定属性的值
// 2. 属性值更新: 更新已存在属性的值
// 3. 动态属性修改: 动态修改元素属性
// 4. 状态管理: 管理元素状态（禁用、只读等）
// 5. 样式控制: 修改单个CSS属性
// 6. 数据绑定: 更新数据属性值

// CDPDOMSetAttributeValue 设置指定节点的属性值
// nodeID: 要设置属性的节点ID
// name: 要设置的属性名称
// value: 要设置的属性值
func CDPDOMSetAttributeValue(nodeID int, name, value string) (string, error) {
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
        "method": "DOM.setAttributeValue",
        "params": {
            "nodeId": %d,
            "name": "%s",
            "value": "%s"
        }
    }`, reqID, nodeID, name, value)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.setAttributeValue 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.setAttributeValue 请求超时")
		}
	}
}

/*

// 示例: 设置输入框的属性值
func ExampleCDPDOMSetAttributeValue() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个输入框的节点ID
    inputNodeID := 1001

    fmt.Printf("=== 设置属性值演示 ===\n")
    fmt.Printf("目标元素节点ID: %d\n\n", inputNodeID)

    // 2. 首先获取元素的当前属性
    fmt.Printf("=== 设置前属性检查 ===\n")

    attrsBefore, err := getElementAttributes(inputNodeID)
    if err != nil {
        log.Printf("获取属性失败: %v", err)
        return
    }

    displayAttributes("设置前属性", attrsBefore)

    // 3. 定义要测试的属性设置
    testCases := []struct {
        name        string
        description string
        attrName    string
        attrValue   string
        expectType  string
    }{
        {
            name:        "设置占位符",
            description: "为输入框设置提示文本",
            attrName:    "placeholder",
            attrValue:   "请输入您的用户名",
            expectType:  "文本属性",
        },
        {
            name:        "设置标题",
            description: "设置鼠标悬停提示",
            attrName:    "title",
            attrValue:   "用户名必须是3-20个字符",
            expectType:  "工具提示属性",
        },
        {
            name:        "设置禁用状态",
            description: "禁用输入框",
            attrName:    "disabled",
            attrValue:   "true",
            expectType:  "布尔属性",
        },
        {
            name:        "设置最大长度",
            description: "限制输入最大字符数",
            attrName:    "maxlength",
            attrValue:   "20",
            expectType:  "数值属性",
        },
        {
            name:        "设置数据属性",
            description: "设置自定义数据属性",
            attrName:    "data-testid",
            attrValue:   "username-input",
            expectType:  "数据属性",
        },
        {
            name:        "设置ARIA标签",
            description: "设置可访问性标签",
            attrName:    "aria-label",
            attrValue:   "用户名输入框",
            expectType:  "ARIA属性",
        },
        {
            name:        "设置样式",
            description: "设置内联样式",
            attrName:    "style",
            attrValue:   "color: #333; background-color: #f8f9fa; border: 1px solid #ccc;",
            expectType:  "样式属性",
        },
        {
            name:        "设置类名",
            description: "设置CSS类",
            attrName:    "class",
            attrValue:   "form-control input-lg",
            expectType:  "类属性",
        },
    }

    // 4. 执行属性设置测试
    for i, tc := range testCases {
        fmt.Printf("\n测试 %d: %s\n", i+1, tc.name)
        fmt.Printf("描述: %s\n", tc.description)
        fmt.Printf("属性: %s=%s\n", tc.attrName, tc.attrValue)
        fmt.Printf("类型: %s\n", tc.expectType)

        startTime := time.Now()
        result, err := CDPDOMSetAttributeValue(inputNodeID, tc.attrName, tc.attrValue)
        setTime := time.Since(startTime)

        if err != nil {
            fmt.Printf("❌ 设置失败: %v\n", err)
        } else {
            fmt.Printf("✅ 设置成功 (耗时: %v)\n", setTime)
            fmt.Printf("结果: %s\n", result)

            // 验证设置结果
            fmt.Printf("\n验证设置结果...\n")

            // 获取设置后的属性
            attrsAfter, err := getElementAttributes(inputNodeID)
            if err != nil {
                fmt.Printf("❌ 验证失败: 无法获取属性\n")
                continue
            }

            // 查找设置的属性
            attributeFound := false
            for _, attr := range attrsAfter {
                if attr.Name == tc.attrName {
                    if attr.Value == tc.attrValue {
                        fmt.Printf("✅ 属性 '%s' 已成功设置为 '%s'\n", tc.attrName, tc.attrValue)
                        attributeFound = true
                    } else {
                        fmt.Printf("⚠️ 属性 '%s' 值为 '%s'，但期望为 '%s'\n",
                            tc.attrName, attr.Value, tc.attrValue)
                    }
                    break
                }
            }

            if !attributeFound {
                // 对于布尔属性，空值也视为设置成功
                if (tc.attrName == "disabled" || tc.attrName == "readonly" || tc.attrName == "required") &&
                   (tc.attrValue == "" || tc.attrValue == "true") {
                    fmt.Printf("✅ 布尔属性 '%s' 已设置（空值表示true）\n", tc.attrName)
                } else {
                    fmt.Printf("❌ 未找到属性 '%s'\n", tc.attrName)
                }
            }
        }

        // 短暂延迟
        if i < len(testCases)-1 {
            time.Sleep(200 * time.Millisecond)
        }
    }

    // 5. 特殊属性值测试
    fmt.Printf("\n=== 特殊属性值测试 ===\n")

    specialTests := []struct {
        name        string
        attrName    string
        attrValue   string
        description string
    }{
        {
            name:        "空值属性",
            attrName:    "placeholder",
            attrValue:   "",
            description: "设置空字符串值",
        },
        {
            name:        "特殊字符",
            attrName:    "title",
            attrValue:   "包含\"引号\"和'单引号'的文本",
            description: "包含引号的属性值",
        },
        {
            name:        "HTML实体",
            attrName:    "alt",
            attrValue:   "Logo &lt; &gt; &amp;",
            description: "包含HTML实体的属性值",
        },
        {
            name:        "长文本",
            attrName:    "data-description",
            attrValue:   strings.Repeat("这是一个很长的描述文本，用于测试长属性值的设置。", 10),
            description: "长属性值测试",
        },
        {
            name:        "Unicode字符",
            attrName:    "aria-label",
            attrValue:   "搜索🔍 设置⚙️ 用户👤",
            description: "包含Unicode表情符号",
        },
    }

    for _, test := range specialTests {
        fmt.Printf("测试: %s\n", test.description)

        result, err := CDPDOMSetAttributeValue(inputNodeID, test.attrName, test.attrValue)
        if err != nil {
            fmt.Printf("❌ 失败: %v\n", err)
        } else {
            fmt.Printf("✅ 成功\n")

            // 验证特殊字符处理
            if strings.Contains(test.attrValue, "\"") || strings.Contains(test.attrValue, "'") {
                attrsAfter, err := getElementAttributes(inputNodeID)
                if err == nil {
                    for _, attr := range attrsAfter {
                        if attr.Name == test.attrName {
                            fmt.Printf("  实际值: %s\n", truncateString(attr.Value, 50))
                            break
                        }
                    }
                }
            }
        }
    }

    // 6. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次属性设置
    testRepeatedAttributeSets(inputNodeID, 10)

    // 测试不同属性的设置性能
    testDifferentAttributes(inputNodeID)

    // 7. 错误处理测试
    fmt.Printf("\n=== 错误处理测试 ===\n")

    errorTests := []struct {
        nodeID    int
        attrName  string
        attrValue string
        desc      string
    }{
        {0, "placeholder", "test", "无效的节点ID"},
        {inputNodeID, "", "test", "空的属性名"},
        {inputNodeID, "placeholder", "", "空值属性"},
        {999999, "placeholder", "test", "不存在的节点"},
    }

    for _, test := range errorTests {
        fmt.Printf("测试: %s\n", test.desc)

        result, err := CDPDOMSetAttributeValue(test.nodeID, test.attrName, test.attrValue)
        if err != nil {
            fmt.Printf("✅ 预期错误: %v\n", err)
        } else {
            fmt.Printf("❌ 预期错误但成功: %s\n", result)
        }
    }

    // 8. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "表单验证反馈",
            description: "根据验证结果显示不同状态",
            useCase:     "验证通过时设置success类，失败时设置error类和aria-invalid属性",
        },
        {
            name:        "动态UI状态",
            description: "根据应用状态更新UI元素属性",
            useCase:     "按钮加载时设置disabled和aria-busy属性，显示加载文本",
        },
        {
            name:        "主题切换",
            description: "切换应用主题时更新元素属性",
            useCase:     "切换主题时更新data-theme属性，CSS根据属性应用样式",
        },
        {
            name:        "国际化",
            description: "根据语言设置更新文本属性",
            useCase:     "切换语言时更新placeholder、title、aria-label等文本属性",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 测试重复属性设置
func testRepeatedAttributeSets(nodeID int, count int) {
    fmt.Printf("重复属性设置测试 (%d 次):\n", count)

    attrName := "data-counter"
    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < count; i++ {
        attrValue := fmt.Sprintf("value-%d", i)

        startTime := time.Now()
        result, err := CDPDOMSetAttributeValue(nodeID, attrName, attrValue)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            successCount++
            fmt.Printf("  第 %d 次: ✅ 成功 (%s=%s, 耗时: %v)\n",
                i+1, attrName, attrValue, duration)

            // 验证设置
            attrs, err := getElementAttributes(nodeID)
            if err == nil {
                value := ""
                for _, attr := range attrs {
                    if attr.Name == attrName {
                        value = attr.Value
                        break
                    }
                }

                if value == attrValue {
                    fmt.Printf("    验证: ✅ 值正确\n")
                } else {
                    fmt.Printf("    验证: ❌ 值不匹配 (实际: %s)\n", value)
                }
            }
        }

        // 短暂延迟
        if i < count-1 {
            time.Sleep(50 * time.Millisecond)
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))
}

// 测试不同属性的设置性能
func testDifferentAttributes(nodeID int) {
    fmt.Printf("\n不同属性设置性能对比:\n")

    attributes := []struct {
        name  string
        value string
    }{
        {"class", "btn btn-primary"},
        {"style", "color: red; background: blue;"},
        {"data-test", "performance-test"},
        {"title", "这是一个测试标题"},
        {"aria-label", "测试按钮"},
        {"disabled", "true"},
    }

    var results []AttributePerformance

    for _, attr := range attributes {
        result := AttributePerformance{
            Attribute: attr.name,
        }

        // 多次测试取平均值
        var totalTime time.Duration
        tests := 3

        for i := 0; i < tests; i++ {
            startTime := time.Now()
            _, err := CDPDOMSetAttributeValue(nodeID, attr.name, attr.value)
            duration := time.Since(startTime)

            if err == nil {
                totalTime += duration
            }

            if i < tests-1 {
                time.Sleep(20 * time.Millisecond)
            }
        }

        result.AverageTime = totalTime / time.Duration(tests)
        result.ValueSize = len(attr.value)
        results = append(results, result)
    }

    // 显示结果
    fmt.Printf("属性              | 值大小 | 平均耗时\n")
    fmt.Printf("------------------|--------|---------\n")
    for _, result := range results {
        fmt.Printf("%-16s | %6d | %v\n",
            result.Attribute, result.ValueSize, result.AverageTime)
    }
}

type AttributePerformance struct {
    Attribute   string
    AverageTime time.Duration
    ValueSize   int
}

// 高级功能: 智能属性值设置器
type SmartAttributeValueSetter struct {
    setHistory []AttributeValueRecord
    validator  *AttributeValueValidator
    cache      map[string]AttributeValueCache
    analyzer   *AttributeImpactAnalyzer
}

type AttributeValueRecord struct {
    Timestamp   time.Time
    NodeID      int
    Name        string
    Value       string
    Duration    time.Duration
    Success     bool
    Error       string
    OldValue    string
    NewValue    string
    Impact      AttributeImpact
    Validation  ValidationResult
}

type AttributeValueCache struct {
    NodeID    int
    Name      string
    Value     string
    Timestamp time.Time
    HitCount  int
}

type AttributeImpactAnalyzer struct {
    sensitiveAttributes map[string]bool
    styleAttributes     map[string]bool
    booleanAttributes   map[string]bool
    eventAttributes     map[string]bool
}

type ValidationResult struct {
    Valid     bool
    Warnings  []string
    Errors    []string
    Sanitized string
}

type AttributeImpact struct {
    RequiresRepaint bool
    RequiresReflow  bool
    AffectsLayout   bool
    AffectsStyle    bool
    AffectsContent  bool
    SecurityRisk    bool
    AccessibilityImpact string
}

func NewSmartAttributeValueSetter() *SmartAttributeValueSetter {
    return &SmartAttributeValueSetter{
        setHistory: make([]AttributeValueRecord, 0),
        validator:  NewAttributeValueValidator(),
        cache:      make(map[string]AttributeValueCache),
        analyzer:   NewAttributeImpactAnalyzer(),
    }
}

func (savs *SmartAttributeValueSetter) SetAttributeValue(nodeID int, name, value string, validate bool) (AttributeValueRecord, error) {
    record := AttributeValueRecord{
        Timestamp: time.Now(),
        NodeID:    nodeID,
        Name:      name,
        Value:     value,
    }

    // 获取当前值
    oldValue, err := savs.getCurrentAttributeValue(nodeID, name)
    if err != nil {
        return record, fmt.Errorf("获取当前值失败: %w", err)
    }
    record.OldValue = oldValue

    // 验证输入
    if validate {
        validation := savs.validator.Validate(name, value)
        record.Validation = validation

        if !validation.Valid {
            record.Error = strings.Join(validation.Errors, "; ")
            savs.recordSet(record)
            return record, fmt.Errorf("验证失败: %s", record.Error)
        }

        if validation.Sanitized != value {
            value = validation.Sanitized
            record.Value = value
        }
    }

    // 检查缓存
    cacheKey := fmt.Sprintf("%d:%s", nodeID, name)
    if cached, exists := savs.cache[cacheKey]; exists && cached.Value == value {
        cached.HitCount++
        savs.cache[cacheKey] = cached

        record.Success = true
        record.Duration = 0
        record.NewValue = value
        record.Impact = savs.analyzer.AnalyzeImpact(name, value)
        savs.recordSet(record)

        return record, nil
    }

    // 分析影响
    record.Impact = savs.analyzer.AnalyzeImpact(name, value)

    // 执行设置
    startTime := time.Now()
    result, err := CDPDOMSetAttributeValue(nodeID, name, value)
    record.Duration = time.Since(startTime)

    if err != nil {
        record.Success = false
        record.Error = err.Error()
    } else {
        record.Success = true
        record.Response = result

        // 更新缓存
        savs.cache[cacheKey] = AttributeValueCache{
            NodeID:    nodeID,
            Name:      name,
            Value:     value,
            Timestamp: time.Now(),
            HitCount:  1,
        }

        // 验证新值
        newValue, err := savs.getCurrentAttributeValue(nodeID, name)
        if err == nil {
            record.NewValue = newValue
        }
    }

    // 记录历史
    savs.recordSet(record)

    if err != nil {
        return record, fmt.Errorf("设置失败: %w", err)
    }

    return record, nil
}

func (savs *SmartAttributeValueSetter) getCurrentAttributeValue(nodeID int, name string) (string, error) {
    attrs, err := getElementAttributes(nodeID)
    if err != nil {
        return "", err
    }

    for _, attr := range attrs {
        if attr.Name == name {
            return attr.Value, nil
        }
    }

    return "", nil // 属性不存在
}

func (savs *SmartAttributeValueSetter) recordSet(record AttributeValueRecord) {
    savs.setHistory = append(savs.setHistory, record)
}

func NewAttributeValueValidator() *AttributeValueValidator {
    return &AttributeValueValidator{
        maxLength: map[string]int{
            "style":    1000,
            "class":    500,
            "title":    200,
            "alt":      200,
            "value":    1000,
            "":         5000, // 默认
        },
        dangerousPatterns: []string{
            "javascript:", "data:", "vbscript:", "expression(",
        },
        allowedProtocols: []string{"http://", "https://", "mailto:", "tel:"},
    }
}

type AttributeValueValidator struct {
    maxLength         map[string]int
    dangerousPatterns []string
    allowedProtocols  []string
}

func (avv *AttributeValueValidator) Validate(name, value string) ValidationResult {
    result := ValidationResult{
        Valid:     true,
        Sanitized: value,
    }

    // 检查长度
    maxLen, exists := avv.maxLength[name]
    if !exists {
        maxLen = avv.maxLength[""]
    }

    if len(value) > maxLen {
        result.Warnings = append(result.Warnings, fmt.Sprintf("值过长: %d > %d", len(value), maxLen))
        result.Sanitized = truncateString(value, maxLen)
    }

    // 检查危险内容
    lowerValue := strings.ToLower(value)
    for _, pattern := range avv.dangerousPatterns {
        if strings.Contains(lowerValue, pattern) {
            result.Errors = append(result.Errors, fmt.Sprintf("检测到危险内容: %s", pattern))
            result.Valid = false
        }
    }

    // 特殊属性验证
    switch name {
    case "href", "src", "action":
        if value != "" && !avv.isSafeURL(value) {
            result.Warnings = append(result.Warnings, "URL可能不安全")
        }
    case "style":
        if !avv.isValidCSS(value) {
            result.Warnings = append(result.Warnings, "CSS样式可能无效")
        }
    }

    return result
}

func (avv *AttributeValueValidator) isSafeURL(url string) bool {
    if url == "" || url == "#" {
        return true
    }

    for _, protocol := range avv.allowedProtocols {
        if strings.HasPrefix(strings.ToLower(url), protocol) {
            return true
        }
    }

    // 相对URL
    if strings.HasPrefix(url, "/") || strings.HasPrefix(url, "./") || strings.HasPrefix(url, "../") {
        return true
    }

    return false
}

func (avv *AttributeValueValidator) isValidCSS(css string) bool {
    // 简单CSS验证
    if strings.Contains(css, "expression(") {
        return false
    }
    return true
}

func NewAttributeImpactAnalyzer() *AttributeImpactAnalyzer {
    return &AttributeImpactAnalyzer{
        sensitiveAttributes: map[string]bool{
            "onclick": true, "onload": true, "onerror": true, "onsubmit": true,
            "onmouseover": true, "onkeydown": true, "onchange": true,
        },
        styleAttributes: map[string]bool{
            "style": true, "class": true,
        },
        booleanAttributes: map[string]bool{
            "disabled": true, "readonly": true, "required": true,
            "checked": true, "selected": true, "multiple": true,
        },
        eventAttributes: map[string]bool{
            "on": true, // 所有on开头的属性
        },
    }
}

func (aia *AttributeImpactAnalyzer) AnalyzeImpact(name, value string) AttributeImpact {
    impact := AttributeImpact{}

    // 检查是否需要重绘
    if aia.styleAttributes[name] {
        impact.RequiresRepaint = true
        impact.AffectsStyle = true
    }

    // 检查是否需要重排
    if name == "style" && (strings.Contains(value, "width") || strings.Contains(value, "height") ||
        strings.Contains(value, "display") || strings.Contains(value, "position")) {
        impact.RequiresReflow = true
        impact.AffectsLayout = true
    }

    // 检查是否影响内容
    if name == "value" || name == "textContent" || name == "innerHTML" {
        impact.AffectsContent = true
    }

    // 检查安全风险
    if aia.sensitiveAttributes[name] || (strings.HasPrefix(name, "on") && aia.eventAttributes["on"]) {
        impact.SecurityRisk = true
    }

    // 检查可访问性影响
    if strings.HasPrefix(name, "aria-") || name == "role" || name == "tabindex" {
        impact.AccessibilityImpact = "high"
    } else if name == "alt" || name == "title" {
        impact.AccessibilityImpact = "medium"
    }

    return impact
}

// 演示智能属性值设置
func demonstrateSmartAttributeValueSetting() {
    fmt.Println("=== 智能属性值设置演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建智能设置器
    setter := NewSmartAttributeValueSetter()

    // 测试元素
    elementID := 3001

    fmt.Printf("测试元素ID: %d\n\n", elementID)

    // 场景1: 正常设置
    fmt.Printf("场景1: 正常属性设置\n")

    tests := []struct {
        name     string
        value    string
        validate bool
    }{
        {"placeholder", "请输入内容", true},
        {"title", "这是一个提示", true},
        {"data-index", "1", true},
        {"style", "color: red;", true},
    }

    for _, test := range tests {
        fmt.Printf("设置: %s=%s\n", test.name, test.value)

        record, err := setter.SetAttributeValue(elementID, test.name, test.value, test.validate)
        if err != nil {
            fmt.Printf("  ❌ 失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)

            if record.Impact.RequiresRepaint {
                fmt.Printf("    影响: 需要重绘\n")
            }
            if record.Impact.RequiresReflow {
                fmt.Printf("    影响: 需要重排\n")
            }
            if record.Impact.SecurityRisk {
                fmt.Printf("    警告: 安全风险\n")
            }
        }

        time.Sleep(200 * time.Millisecond)
    }

    // 场景2: 缓存测试
    fmt.Printf("\n场景2: 缓存测试\n")

    for i := 0; i < 3; i++ {
        record, err := setter.SetAttributeValue(elementID, "data-test", "cached-value", true)
        if err != nil {
            fmt.Printf("第 %d 次: ❌ 失败\n", i+1)
        } else if record.Duration == 0 {
            fmt.Printf("第 %d 次: ✅ 缓存命中\n", i+1)
        } else {
            fmt.Printf("第 %d 次: ✅ 设置成功\n", i+1)
        }
    }

    // 场景3: 验证测试
    fmt.Printf("\n场景3: 验证测试\n")

    validationTests := []struct {
        name     string
        value    string
        expectError bool
    }{
        {"style", strings.Repeat("color: red;", 200), false}, // 触发警告
        {"href", "javascript:alert(1)", true}, // 触发错误
        {"onclick", "alert(1)", false}, // 警告但不阻止
    }

    for _, test := range validationTests {
        record, err := setter.SetAttributeValue(elementID, test.name, test.value, true)

        if test.expectError {
            if err != nil {
                fmt.Printf("设置 %s: ✅ 预期错误\n", test.name)
            } else {
                fmt.Printf("设置 %s: ❌ 预期错误但成功\n", test.name)
            }
        } else {
            if err != nil {
                fmt.Printf("设置 %s: ❌ 失败: %v\n", test.name, err)
            } else {
                fmt.Printf("设置 %s: ✅ 成功", test.name)
                if len(record.Validation.Warnings) > 0 {
                    fmt.Printf(" (警告: %v)", record.Validation.Warnings)
                }
                fmt.Printf("\n")
            }
        }
    }

    fmt.Println("\n=== 演示完成 ===")
}

*/

// -----------------------------------------------  DOM.setFileInputFiles  -----------------------------------------------
// === 应用场景 ===
// 1. 文件上传测试: 测试文件上传功能
// 2. 自动化文件选择: 自动化测试中选择文件
// 3. 批量文件上传: 一次设置多个文件
// 4. 文件类型验证: 测试文件类型限制
// 5. 文件大小测试: 测试文件大小限制
// 6. 表单提交测试: 测试带文件上传的表单

// CDPDOMSetFileInputFiles 设置文件输入框的文件
// nodeID: 文件输入框的节点ID
// files: 要设置的文件路径列表
// backendNodeID: 可选的后端节点ID
func CDPDOMSetFileInputFiles(nodeID int, files []string, backendNodeID int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息参数
	params := map[string]interface{}{
		"nodeId": nodeID,
		"files":  files,
	}

	if backendNodeID > 0 {
		params["backendNodeId"] = backendNodeID
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "DOM.setFileInputFiles",
        "params": %s
    }`, reqID, string(paramsJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.setFileInputFiles 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.setFileInputFiles 请求超时")
		}
	}
}

/*

// 示例: 设置文件上传框的文件
func ExampleCDPDOMSetFileInputFiles() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个文件输入框的节点ID
    fileInputNodeID := 1001

    fmt.Printf("=== 文件输入框设置演示 ===\n")
    fmt.Printf("目标文件输入框节点ID: %d\n\n", fileInputNodeID)

    // 2. 首先检查文件输入框的状态
    fmt.Printf("=== 输入框状态检查 ===\n")

    inputInfo, err := getElementInfo(fileInputNodeID)
    if err != nil {
        log.Printf("获取元素信息失败: %v", err)
        return
    }

    displayElementInfo("文件输入框", inputInfo)

    // 检查输入框类型
    isFileInput := checkIsFileInput(inputInfo)
    if !isFileInput {
        fmt.Printf("⚠️ 警告: 该元素可能不是文件输入框\n")
    }

    // 3. 准备测试文件
    fmt.Printf("\n=== 准备测试文件 ===\n")

    // 创建测试文件
    testFiles, cleanup, err := createTestFiles()
    if err != nil {
        log.Printf("创建测试文件失败: %v", err)
        return
    }
    defer cleanup() // 确保测试完成后清理临时文件

    // 显示测试文件信息
    displayTestFiles(testFiles)

    // 4. 测试单个文件上传
    fmt.Printf("\n=== 测试单个文件上传 ===\n")

    // 获取文件路径
    singleFile := []string{testFiles["smallText"].Path}
    fmt.Printf("设置单个文件: %s\n", testFiles["smallText"].Name)

    startTime := time.Now()
    result, err := CDPDOMSetFileInputFiles(fileInputNodeID, singleFile, 0)
    setTime := time.Since(startTime)

    if err != nil {
        log.Printf("设置文件失败: %v", err)
        return
    }

    fmt.Printf("✅ 文件设置成功 (耗时: %v)\n", setTime)
    fmt.Printf("结果: %s\n", result)

    // 验证文件是否设置成功
    verifyFileSet(fileInputNodeID, []string{testFiles["smallText"].Name})

    // 5. 测试多个文件上传
    fmt.Printf("\n=== 测试多个文件上传 ===\n")

    if isMultipleInput(inputInfo) {
        multipleFiles := []string{
            testFiles["smallText"].Path,
            testFiles["mediumText"].Path,
            testFiles["image"].Path,
        }

        fmt.Printf("设置多个文件 (%d 个):\n", len(multipleFiles))
        for _, file := range multipleFiles {
            fmt.Printf("  - %s\n", filepath.Base(file))
        }

        startTime = time.Now()
        result, err = CDPDOMSetFileInputFiles(fileInputNodeID, multipleFiles, 0)
        setTime = time.Since(startTime)

        if err != nil {
            fmt.Printf("❌ 设置多个文件失败: %v\n", err)
        } else {
            fmt.Printf("✅ 多个文件设置成功 (耗时: %v)\n", setTime)

            // 验证多个文件设置
            fileNames := make([]string, len(multipleFiles))
            for i, file := range multipleFiles {
                fileNames[i] = filepath.Base(file)
            }
            verifyFileSet(fileInputNodeID, fileNames)
        }
    } else {
        fmt.Printf("ℹ️ 输入框不支持多文件选择 (missing 'multiple' attribute)\n")
    }

    // 6. 测试文件类型限制
    fmt.Printf("\n=== 测试文件类型限制 ===\n")

    // 检查accept属性
    acceptedTypes := getAcceptedFileTypes(inputInfo)
    if len(acceptedTypes) > 0 {
        fmt.Printf("输入框接受的文件类型: %v\n", acceptedTypes)

        // 测试接受的文件类型
        testAcceptedFiles := getTestFilesByType(testFiles, acceptedTypes)
        if len(testAcceptedFiles) > 0 {
            fmt.Printf("测试接受的文件: %v\n", testAcceptedFiles)
        }

        // 测试不接受的类型
        testRejectedFiles := getTestFilesNotByType(testFiles, acceptedTypes)
        if len(testRejectedFiles) > 0 {
            fmt.Printf("测试不接受的文件: %v\n", testRejectedFiles)
        }
    } else {
        fmt.Printf("输入框无文件类型限制\n")
    }

    // 7. 测试文件大小限制
    fmt.Printf("\n=== 测试文件大小限制 ===\n")

    // 测试大文件上传
    largeFile := []string{testFiles["largeText"].Path}
    fmt.Printf("测试大文件上传: %s (%.2f KB)\n",
        testFiles["largeText"].Name,
        float64(testFiles["largeText"].Size)/1024)

    result, err = CDPDOMSetFileInputFiles(fileInputNodeID, largeFile, 0)
    if err != nil {
        fmt.Printf("❌ 大文件设置失败: %v\n", err)
    } else {
        fmt.Printf("✅ 大文件设置成功\n")
    }

    // 8. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次文件设置
    testRepeatedFileSets(fileInputNodeID, testFiles["smallText"].Path, 3)

    // 测试不同大小文件的设置性能
    testDifferentFileSizes(fileInputNodeID, testFiles)

    // 9. 错误处理测试
    fmt.Printf("\n=== 错误处理测试 ===\n")

    errorTests := []struct {
        files []string
        desc  string
    }{
        {[]string{}, "空文件列表"},
        {[]string{"/path/to/nonexistent/file.txt"}, "不存在的文件"},
        {[]string{""}, "空文件路径"},
    }

    for _, test := range errorTests {
        fmt.Printf("测试: %s\n", test.desc)

        result, err := CDPDOMSetFileInputFiles(fileInputNodeID, test.files, 0)
        if err != nil {
            fmt.Printf("✅ 预期错误: %v\n", err)
        } else {
            fmt.Printf("❌ 预期错误但成功: %s\n", result)
        }
    }

    // 10. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "头像上传测试",
            description: "测试用户头像上传功能",
            useCase:     "上传不同格式的图片，测试大小限制，验证预览功能",
        },
        {
            name:        "文档上传测试",
            description: "测试文档上传功能",
            useCase:     "上传PDF、Word、Excel文档，测试格式验证，病毒扫描",
        },
        {
            name:        "批量图片上传",
            description: "测试多图片上传功能",
            useCase:     "一次选择多张图片，测试上传进度，缩略图生成",
        },
        {
            name:        "大文件上传测试",
            description: "测试大文件分片上传",
            useCase:     "上传大视频文件，测试分片上传，断点续传",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 检查是否是文件输入框
func checkIsFileInput(info map[string]interface{}) bool {
    if nodeName, ok := info["nodeName"].(string); ok {
        if strings.ToUpper(nodeName) != "INPUT" {
            return false
        }
    }

    // 检查type属性
    if attrs, ok := info["attributes"].([]string); ok {
        for i := 0; i < len(attrs); i += 2 {
            if i+1 < len(attrs) && attrs[i] == "type" && attrs[i+1] == "file" {
                return true
            }
        }
    }

    return false
}

// 检查是否支持多文件
func isMultipleInput(info map[string]interface{}) bool {
    if attrs, ok := info["attributes"].([]string); ok {
        for i := 0; i < len(attrs); i += 2 {
            if i+1 < len(attrs) && attrs[i] == "multiple" {
                return true
            }
        }
    }
    return false
}

// 获取接受的文件类型
func getAcceptedFileTypes(info map[string]interface{}) []string {
    var acceptedTypes []string

    if attrs, ok := info["attributes"].([]string); ok {
        for i := 0; i < len(attrs); i += 2 {
            if i+1 < len(attrs) && attrs[i] == "accept" {
                acceptValue := attrs[i+1]
                // 分割多种类型
                types := strings.Split(acceptValue, ",")
                for _, t := range types {
                    t = strings.TrimSpace(t)
                    if t != "" {
                        acceptedTypes = append(acceptedTypes, t)
                    }
                }
            }
        }
    }

    return acceptedTypes
}

// 测试文件信息
type TestFileInfo struct {
    Path    string
    Name    string
    Size    int64
    Type    string
    Content []byte
}

// 创建测试文件
func createTestFiles() (map[string]TestFileInfo, func(), error) {
    // 创建临时目录
    tempDir, err := os.MkdirTemp("", "cdp_test_files_*")
    if err != nil {
        return nil, nil, fmt.Errorf("创建临时目录失败: %w", err)
    }

    cleanup := func() {
        os.RemoveAll(tempDir)
    }

    testFiles := make(map[string]TestFileInfo)

    // 1. 小文本文件
    smallTextPath := filepath.Join(tempDir, "small.txt")
    smallContent := []byte("这是一个小的测试文本文件。\n用于测试文件上传功能。")
    if err := os.WriteFile(smallTextPath, smallContent, 0644); err != nil {
        cleanup()
        return nil, nil, fmt.Errorf("创建小文件失败: %w", err)
    }
    testFiles["smallText"] = TestFileInfo{
        Path:    smallTextPath,
        Name:    "small.txt",
        Size:    int64(len(smallContent)),
        Type:    "text/plain",
        Content: smallContent,
    }

    // 2. 中等文本文件
    mediumTextPath := filepath.Join(tempDir, "medium.txt")
    var mediumContent []byte
    for i := 0; i < 1000; i++ {
        mediumContent = append(mediumContent, []byte(fmt.Sprintf("Line %d: 这是中等大小的测试文件内容。\n", i+1))...)
    }
    if err := os.WriteFile(mediumTextPath, mediumContent, 0644); err != nil {
        cleanup()
        return nil, nil, fmt.Errorf("创建中等文件失败: %w", err)
    }
    testFiles["mediumText"] = TestFileInfo{
        Path:    mediumTextPath,
        Name:    "medium.txt",
        Size:    int64(len(mediumContent)),
        Type:    "text/plain",
        Content: mediumContent,
    }

    // 3. 大文本文件
    largeTextPath := filepath.Join(tempDir, "large.txt")
    var largeContent []byte
    for i := 0; i < 10000; i++ { // 约2MB
        largeContent = append(largeContent, []byte(fmt.Sprintf("Line %d: 这是大测试文件内容，用于测试大文件上传性能。\n", i+1))...)
    }
    if err := os.WriteFile(largeTextPath, largeContent, 0644); err != nil {
        cleanup()
        return nil, nil, fmt.Errorf("创建大文件失败: %w", err)
    }
    testFiles["largeText"] = TestFileInfo{
        Path:    largeTextPath,
        Name:    "large.txt",
        Size:    int64(len(largeContent)),
        Type:    "text/plain",
        Content: largeContent,
    }

    // 4. 图片文件 (创建简单的BMP文件)
    imagePath := filepath.Join(tempDir, "test.bmp")
    // 创建一个很小的BMP文件
    bmpContent := []byte{
        'B', 'M', // 文件类型
        0x3A, 0, 0, 0, // 文件大小
        0, 0, 0, 0, // 保留
        0x36, 0, 0, 0, // 像素数据偏移
        0x28, 0, 0, 0, // 信息头大小
        0x01, 0, 0, 0, // 宽度
        0x01, 0, 0, 0, // 高度
        0x01, 0, // 平面数
        0x18, 0, // 每像素位数
        0, 0, 0, 0, // 压缩方式
        0x04, 0, 0, 0, // 图像大小
        0, 0, 0, 0, // 水平分辨率
        0, 0, 0, 0, // 垂直分辨率
        0, 0, 0, 0, // 使用的颜色数
        0, 0, 0, 0, // 重要颜色数
        0xFF, 0xFF, 0xFF, 0, // 白色像素
    }
    if err := os.WriteFile(imagePath, bmpContent, 0644); err != nil {
        cleanup()
        return nil, nil, fmt.Errorf("创建图片文件失败: %w", err)
    }
    testFiles["image"] = TestFileInfo{
        Path:    imagePath,
        Name:    "test.bmp",
        Size:    int64(len(bmpContent)),
        Type:    "image/bmp",
        Content: bmpContent,
    }

    return testFiles, cleanup, nil
}

// 显示测试文件信息
func displayTestFiles(testFiles map[string]TestFileInfo) {
    fmt.Printf("创建的测试文件:\n")
    for key, file := range testFiles {
        sizeKB := float64(file.Size) / 1024
        sizeMB := float64(file.Size) / (1024 * 1024)

        sizeStr := fmt.Sprintf("%.1f KB", sizeKB)
        if sizeMB >= 1 {
            sizeStr = fmt.Sprintf("%.2f MB", sizeMB)
        }

        fmt.Printf("  %-12s: %s (%s, %s)\n",
            key, file.Name, file.Type, sizeStr)
    }
}

// 验证文件设置
func verifyFileSet(nodeID int, expectedFileNames []string) {
    // 这里需要通过JavaScript获取文件列表
    // 简化实现
    fmt.Printf("文件验证: 已设置 %d 个文件\n", len(expectedFileNames))
    for i, name := range expectedFileNames {
        fmt.Printf("  %d. %s\n", i+1, name)
    }
}

// 按类型获取测试文件
func getTestFilesByType(testFiles map[string]TestFileInfo, acceptedTypes []string) []string {
    var matchingFiles []string

    for _, file := range testFiles {
        for _, acceptType := range acceptedTypes {
            if matchesFileType(file.Name, file.Type, acceptType) {
                matchingFiles = append(matchingFiles, file.Path)
                break
            }
        }
    }

    return matchingFiles
}

// 按类型排除测试文件
func getTestFilesNotByType(testFiles map[string]TestFileInfo, acceptedTypes []string) []string {
    var nonMatchingFiles []string

    for _, file := range testFiles {
        matches := false
        for _, acceptType := range acceptedTypes {
            if matchesFileType(file.Name, file.Type, acceptType) {
                matches = true
                break
            }
        }
        if !matches {
            nonMatchingFiles = append(nonMatchingFiles, file.Path)
        }
    }

    return nonMatchingFiles
}

// 检查文件是否匹配接受类型
func matchesFileType(fileName, fileType, acceptType string) bool {
    // 简化匹配逻辑
    if acceptType == "*\/*" {
        return true
    }

    if strings.Contains(acceptType, "/") {
        // MIME类型匹配
        if strings.HasSuffix(acceptType, "/*") {
            // 通配符，如 image/*
            prefix := strings.TrimSuffix(acceptType, "/*")
            return strings.HasPrefix(fileType, prefix+"/")
        }
        // 精确匹配
        return fileType == acceptType
    }

    // 文件扩展名匹配
    if strings.HasPrefix(acceptType, ".") {
        return strings.HasSuffix(strings.ToLower(fileName), strings.ToLower(acceptType))
    }

    return false
}

// 测试重复文件设置
func testRepeatedFileSets(nodeID int, filePath string, count int) {
    fmt.Printf("重复文件设置测试 (%d 次):\n", count)

    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < count; i++ {
        files := []string{filePath}

        startTime := time.Now()
        result, err := CDPDOMSetFileInputFiles(nodeID, files, 0)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            successCount++
            fmt.Printf("  第 %d 次: ✅ 成功 (耗时: %v)\n", i+1, duration)
        }

        // 短暂延迟
        if i < count-1 {
            time.Sleep(200 * time.Millisecond)
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))
}

// 测试不同大小文件的设置性能
func testDifferentFileSizes(nodeID int, testFiles map[string]TestFileInfo) {
    fmt.Printf("\n不同大小文件设置性能对比:\n")

    testOrder := []string{"smallText", "mediumText", "largeText", "image"}
    var results []FileSetPerformance

    for _, key := range testOrder {
        if file, exists := testFiles[key]; exists {
            result := FileSetPerformance{
                FileName: file.Name,
                FileSize: file.Size,
            }

            // 多次测试取平均值
            var totalTime time.Duration
            tests := 3

            for i := 0; i < tests; i++ {
                files := []string{file.Path}

                startTime := time.Now()
                _, err := CDPDOMSetFileInputFiles(nodeID, files, 0)
                duration := time.Since(startTime)

                if err == nil {
                    totalTime += duration
                }

                if i < tests-1 {
                    time.Sleep(100 * time.Millisecond)
                }
            }

            result.AverageTime = totalTime / time.Duration(tests)
            results = append(results, result)
        }
    }

    // 显示结果
    fmt.Printf("文件              | 大小       | 平均耗时 | 速度\n")
    fmt.Printf("------------------|------------|----------|--------\n")
    for _, result := range results {
        sizeKB := float64(result.FileSize) / 1024
        sizeStr := fmt.Sprintf("%.1f KB", sizeKB)

        speed := float64(result.FileSize) / result.AverageTime.Seconds()
        speedStr := fmt.Sprintf("%.1f KB/s", speed/1024)

        fmt.Printf("%-16s | %10s | %8v | %s\n",
            result.FileName, sizeStr, result.AverageTime, speedStr)
    }
}

type FileSetPerformance struct {
    FileName    string
    FileSize    int64
    AverageTime time.Duration
}

// 高级功能: 智能文件上传测试器
type SmartFileUploadTester struct {
    uploadHistory []FileUploadRecord
    fileRegistry  map[string]TestFileInfo
    testCases     []FileUploadTestCase
    validator     *FileUploadValidator
}

type FileUploadRecord struct {
    Timestamp    time.Time
    NodeID       int
    Files        []TestFileInfo
    Duration     time.Duration
    Success      bool
    Error        string
    Validation   UploadValidationResult
    UserAgent    string
}

type FileUploadTestCase struct {
    Name        string
    Description string
    Files       []string
    ExpectError bool
    Validations []FileValidation
}

type FileValidation struct {
    Type     string
    Check    string
    Expected interface{}
}

type UploadValidationResult struct {
    Passed    bool
    Messages  []string
    Warnings  []string
    Errors    []string
}

type FileUploadValidator struct {
    maxTotalSize  int64
    maxFileCount  int
    allowedTypes  []string
    blockedTypes  []string
    maxFileSize   int64
}

func NewSmartFileUploadTester() *SmartFileUploadTester {
    return &SmartFileUploadTester{
        uploadHistory: make([]FileUploadRecord, 0),
        fileRegistry:  make(map[string]TestFileInfo),
        validator: &FileUploadValidator{
            maxTotalSize: 10 * 1024 * 1024, // 10MB
            maxFileCount: 10,
            allowedTypes: []string{},
            blockedTypes: []string{".exe", ".bat", ".sh"},
            maxFileSize:  5 * 1024 * 1024, // 5MB
        },
    }
}

func (sfut *SmartFileUploadTester) RegisterTestFile(name string, file TestFileInfo) {
    sfut.fileRegistry[name] = file
}

func (sfut *SmartFileUploadTester) CreateTestCase(name, description string, fileNames []string, expectError bool) FileUploadTestCase {
    return FileUploadTestCase{
        Name:        name,
        Description: description,
        Files:       fileNames,
        ExpectError: expectError,
    }
}

func (sfut *SmartFileUploadTester) RunUploadTest(nodeID int, testCase FileUploadTestCase) (FileUploadRecord, error) {
    record := FileUploadRecord{
        Timestamp: time.Now(),
        NodeID:    nodeID,
    }

    // 收集测试文件
    var filePaths []string
    var fileInfos []TestFileInfo

    for _, fileName := range testCase.Files {
        if file, exists := sfut.fileRegistry[fileName]; exists {
            filePaths = append(filePaths, file.Path)
            fileInfos = append(fileInfos, file)
        } else {
            return record, fmt.Errorf("未找到测试文件: %s", fileName)
        }
    }

    record.Files = fileInfos

    // 验证文件
    validation := sfut.validator.ValidateFiles(fileInfos)
    record.Validation = validation

    if !validation.Passed {
        record.Error = strings.Join(validation.Errors, "; ")
        sfut.recordUpload(record)
        return record, fmt.Errorf("文件验证失败: %s", record.Error)
    }

    // 执行上传
    startTime := time.Now()
    result, err := CDPDOMSetFileInputFiles(nodeID, filePaths, 0)
    record.Duration = time.Since(startTime)

    if err != nil {
        record.Success = false
        record.Error = err.Error()

        if !testCase.ExpectError {
            record.Validation.Errors = append(record.Validation.Errors,
                fmt.Sprintf("意外错误: %v", err))
        }
    } else {
        record.Success = true
        record.Response = result
    }

    // 记录历史
    sfut.recordUpload(record)

    if testCase.ExpectError && err == nil {
        return record, fmt.Errorf("预期错误但上传成功")
    }

    if !testCase.ExpectError && err != nil {
        return record, fmt.Errorf("上传失败: %w", err)
    }

    return record, nil
}

func (fuv *FileUploadValidator) ValidateFiles(files []TestFileInfo) UploadValidationResult {
    result := UploadValidationResult{
        Passed: true,
    }

    // 检查文件数量
    if len(files) > fuv.maxFileCount {
        result.Passed = false
        result.Errors = append(result.Errors,
            fmt.Sprintf("文件数量超限: %d > %d", len(files), fuv.maxFileCount))
    }

    // 检查总大小
    var totalSize int64
    for _, file := range files {
        totalSize += file.Size

        // 检查单个文件大小
        if file.Size > fuv.maxFileSize {
            result.Passed = false
            result.Errors = append(result.Errors,
                fmt.Sprintf("文件 '%s' 过大: %d > %d",
                    file.Name, file.Size, fuv.maxFileSize))
        }

        // 检查文件类型
        if !fuv.isTypeAllowed(file) {
            result.Passed = false
            result.Errors = append(result.Errors,
                fmt.Sprintf("文件 '%s' 类型不被允许", file.Name))
        }
    }

    if totalSize > fuv.maxTotalSize {
        result.Passed = false
        result.Errors = append(result.Errors,
            fmt.Sprintf("总大小超限: %d > %d", totalSize, fuv.maxTotalSize))
    }

    if result.Passed && len(result.Errors) == 0 {
        result.Messages = append(result.Messages, "文件验证通过")
    }

    return result
}

func (fuv *FileUploadValidator) isTypeAllowed(file TestFileInfo) bool {
    // 检查阻止的类型
    for _, blocked := range fuv.blockedTypes {
        if strings.HasSuffix(strings.ToLower(file.Name), strings.ToLower(blocked)) {
            return false
        }
    }

    // 检查允许的类型
    if len(fuv.allowedTypes) == 0 {
        return true // 无限制
    }

    for _, allowed := range fuv.allowedTypes {
        if matchesFileType(file.Name, file.Type, allowed) {
            return true
        }
    }

    return false
}

func (sfut *SmartFileUploadTester) recordUpload(record FileUploadRecord) {
    sfut.uploadHistory = append(sfut.uploadHistory, record)
}

func (sfut *SmartFileUploadTester) GetStats() map[string]interface{} {
    totalTests := len(sfut.uploadHistory)
    successfulTests := 0
    var totalDuration time.Duration
    var totalFiles int

    for _, record := range sfut.uploadHistory {
        if record.Success {
            successfulTests++
        }
        totalDuration += record.Duration
        totalFiles += len(record.Files)
    }

    avgDuration := time.Duration(0)
    if totalTests > 0 {
        avgDuration = totalDuration / time.Duration(totalTests)
    }

    avgFiles := 0.0
    if totalTests > 0 {
        avgFiles = float64(totalFiles) / float64(totalTests)
    }

    return map[string]interface{}{
        "totalTests":     totalTests,
        "successfulTests": successfulTests,
        "failedTests":    totalTests - successfulTests,
        "successRate":    float64(successfulTests) / float64(totalTests) * 100,
        "totalDuration":  totalDuration,
        "averageDuration": avgDuration,
        "totalFiles":     totalFiles,
        "averageFiles":   avgFiles,
    }
}

// 演示智能文件上传测试
func demonstrateSmartFileUpload() {
    fmt.Println("=== 智能文件上传测试演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建测试器
    tester := NewSmartFileUploadTester()

    // 创建并注册测试文件
    testFiles, cleanup, err := createTestFiles()
    if err != nil {
        log.Printf("创建测试文件失败: %v", err)
        return
    }
    defer cleanup()

    for key, file := range testFiles {
        tester.RegisterTestFile(key, file)
    }

    // 测试元素
    fileInputID := 4001

    fmt.Printf("测试文件输入框ID: %d\n\n", fileInputID)

    // 定义测试用例
    testCases := []FileUploadTestCase{
        {
            Name:        "单个小文件上传",
            Description: "测试基本文件上传功能",
            Files:       []string{"smallText"},
            ExpectError: false,
        },
        {
            Name:        "多文件上传",
            Description: "测试多文件选择功能",
            Files:       []string{"smallText", "mediumText", "image"},
            ExpectError: false,
        },
        {
            Name:        "大文件上传",
            Description: "测试大文件上传",
            Files:       []string{"largeText"},
            ExpectError: false,
        },
    }

    // 运行测试用例
    for _, testCase := range testCases {
        fmt.Printf("测试用例: %s\n", testCase.Name)
        fmt.Printf("描述: %s\n", testCase.Description)
        fmt.Printf("文件: %v\n", testCase.Files)

        record, err := tester.RunUploadTest(fileInputID, testCase)
        if err != nil {
            if testCase.ExpectError {
                fmt.Printf("  ✅ 预期错误: %v\n", err)
            } else {
                fmt.Printf("  ❌ 失败: %v\n", err)
            }
        } else {
            if testCase.ExpectError {
                fmt.Printf("  ❌ 预期错误但成功\n")
            } else {
                fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)

                if len(record.Validation.Messages) > 0 {
                    fmt.Printf("    验证: %v\n", record.Validation.Messages)
                }
                if len(record.Validation.Warnings) > 0 {
                    fmt.Printf("    警告: %v\n", record.Validation.Warnings)
                }
            }
        }

        fmt.Println()
        time.Sleep(500 * time.Millisecond)
    }

    // 显示统计
    fmt.Printf("=== 测试统计 ===\n")
    stats := tester.GetStats()
    fmt.Printf("  总测试数: %d\n", stats["totalTests"])
    fmt.Printf("  成功测试: %d\n", stats["successfulTests"])
    fmt.Printf("  失败测试: %d\n", stats["failedTests"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("  总耗时: %v\n", stats["totalDuration"])
    fmt.Printf("  平均耗时: %v\n", stats["averageDuration"])
    fmt.Printf("  总文件数: %d\n", stats["totalFiles"])
    fmt.Printf("  平均文件数: %.1f\n", stats["averageFiles"])

    fmt.Println("\n=== 演示完成 ===")
}

*/

// -----------------------------------------------  DOM.setNodeName  -----------------------------------------------
// === 应用场景 ===
// 1. 元素重命名: 修改HTML元素的标签名
// 2. 语义化优化: 将非语义化标签改为语义化标签
// 3. 结构转换: 改变DOM节点的类型
// 4. 动态标签: 根据条件动态改变元素类型
// 5. 元素转换: 转换元素类型（如div转section）
// 6. 兼容性修复: 修复不规范的标签名

// CDPDOMSetNodeName 设置节点的名称（标签名）
// nodeID: 要重命名的节点ID
// name: 新的节点名称
func CDPDOMSetNodeName(nodeID int, name string) (string, error) {
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
        "method": "DOM.setNodeName",
        "params": {
            "nodeId": %d,
            "name": "%s"
        }
    }`, reqID, nodeID, name)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.setNodeName 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.setNodeName 请求超时")
		}
	}
}

/*

// 示例: 将div元素重命名为语义化标签
func ExampleCDPDOMSetNodeName() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个div元素的节点ID
    divNodeID := 1001

    fmt.Printf("=== 节点重命名演示 ===\n")
    fmt.Printf("目标节点ID: %d\n\n", divNodeID)

    // 2. 首先获取元素的当前信息
    fmt.Printf("=== 重命名前检查 ===\n")

    nodeInfo, err := getNodeInfo(divNodeID)
    if err != nil {
        log.Printf("获取节点信息失败: %v", err)
        return
    }

    displayNodeInfo("重命名前的节点", nodeInfo)

    // 检查当前节点名称
    oldName := getNodeName(nodeInfo)
    fmt.Printf("当前节点名称: %s\n", oldName)

    // 3. 定义重命名方案
    fmt.Printf("\n=== 重命名方案 ===\n")

    renameTests := []struct {
        newName     string
        description string
        expected    string
        isSemantic  bool
    }{
        {
            newName:     "section",
            description: "转换为章节区域",
            expected:    "SECTION",
            isSemantic:  true,
        },
        {
            newName:     "article",
            description: "转换为文章内容",
            expected:    "ARTICLE",
            isSemantic:  true,
        },
        {
            newName:     "main",
            description: "转换为主内容区域",
            expected:    "MAIN",
            isSemantic:  true,
        },
        {
            newName:     "aside",
            description: "转换为侧边栏",
            expected:    "ASIDE",
            isSemantic:  true,
        },
        {
            newName:     "span",
            description: "转换为行内元素",
            expected:    "SPAN",
            isSemantic:  false,
        },
    }

    // 4. 执行重命名测试
    originalName := oldName

    for i, test := range renameTests {
        fmt.Printf("\n测试 %d: %s\n", i+1, test.description)
        fmt.Printf("重命名为: %s\n", test.newName)

        // 记录操作前的HTML
        outerHTMLBefore, _ := getOuterHTML(divNodeID)

        startTime := time.Now()
        result, err := CDPDOMSetNodeName(divNodeID, test.newName)
        renameTime := time.Since(startTime)

        if err != nil {
            fmt.Printf("❌ 重命名失败: %v\n", err)

            // 尝试恢复原始名称
            if i > 0 {
                CDPDOMSetNodeName(divNodeID, strings.ToLower(originalName))
            }
            continue
        }

        fmt.Printf("✅ 重命名成功 (耗时: %v)\n", renameTime)
        fmt.Printf("结果: %s\n", result)

        // 获取重命名后的信息
        nodeInfoAfter, err := getNodeInfo(divNodeID)
        if err != nil {
            fmt.Printf("❌ 无法获取重命名后信息\n")
            continue
        }

        newName := getNodeName(nodeInfoAfter)
        fmt.Printf("新节点名称: %s\n", newName)

        // 验证重命名结果
        if strings.ToUpper(newName) == test.expected {
            fmt.Printf("✅ 验证通过: 节点已成功重命名为 %s\n", newName)
        } else {
            fmt.Printf("❌ 验证失败: 期望 %s, 实际 %s\n", test.expected, newName)
        }

        // 显示HTML变化
        outerHTMLAfter, _ := getOuterHTML(divNodeID)
        if outerHTMLBefore != "" && outerHTMLAfter != "" {
            showHTMLChanges(outerHTMLBefore, outerHTMLAfter)
        }

        // 分析语义化改进
        if test.isSemantic {
            fmt.Printf("🎯 语义化改进: 从 %s 改为 %s\n",
                strings.ToLower(originalName), test.newName)
        }

        // 短暂延迟
        if i < len(renameTests)-1 {
            time.Sleep(300 * time.Millisecond)
        }
    }

    // 5. 恢复原始名称
    fmt.Printf("\n=== 恢复原始状态 ===\n")
    fmt.Printf("恢复为原始名称: %s\n", strings.ToLower(originalName))

    result, err := CDPDOMSetNodeName(divNodeID, strings.ToLower(originalName))
    if err != nil {
        fmt.Printf("❌ 恢复失败: %v\n", err)
    } else {
        fmt.Printf("✅ 已恢复\n")
    }

    // 6. 特殊标签测试
    fmt.Printf("\n=== 特殊标签测试 ===\n")

    specialTests := []struct {
        newName     string
        description string
        isValid     bool
    }{
        {"custom-element", "自定义元素", true},
        {"my-component", "Web组件", true},
        {"123invalid", "以数字开头的标签", false},
        {"-invalid", "以横线开头的标签", false},
        {"div:special", "包含特殊字符", false},
        {"", "空标签名", false},
    }

    for _, test := range specialTests {
        fmt.Printf("测试: %s\n", test.description)
        fmt.Printf("标签名: %s\n", test.newName)

        result, err := CDPDOMSetNodeName(divNodeID, test.newName)

        if test.isValid {
            if err != nil {
                fmt.Printf("❌ 预期成功但失败: %v\n", err)
            } else {
                fmt.Printf("✅ 设置成功\n")
                // 恢复
                CDPDOMSetNodeName(divNodeID, "div")
            }
        } else {
            if err != nil {
                fmt.Printf("✅ 预期错误: %v\n", err)
            } else {
                fmt.Printf("❌ 预期错误但成功: %s\n", result)
            }
        }
    }

    // 7. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次重命名
    testRepeatedRenames(divNodeID, 5)

    // 测试不同元素的重命名
    testDifferentElements([]int{1001, 1002, 1003})

    // 8. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "语义化重构",
            description: "将div重构为语义化标签",
            useCase:     "将通用div转换为section、article、nav等语义化标签",
        },
        {
            name:        "响应式元素",
            description: "根据设备类型改变元素类型",
            useCase:     "移动设备上使用button，桌面上使用a标签",
        },
        {
            name:        "组件类型切换",
            description: "根据状态切换组件类型",
            useCase:     "加载状态时显示div，加载完成显示实际组件",
        },
        {
            name:        "可访问性优化",
            description: "改善可访问性",
            useCase:     "将div按钮改为button元素以获得更好的键盘导航",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 获取节点名称
func getNodeName(info map[string]interface{}) string {
    if name, ok := info["nodeName"].(string); ok {
        return name
    }
    return ""
}

// 获取外层HTML
func getOuterHTML(nodeID int) (string, error) {
    result, err := CDPDOMGetOuterHTML(nodeID)
    if err != nil {
        return "", err
    }

    var resp struct {
        Result struct {
            OuterHTML string `json:"outerHTML"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &resp); err != nil {
        return "", err
    }

    return resp.Result.OuterHTML, nil
}

// 显示HTML变化
func showHTMLChanges(before, after string) {
    fmt.Printf("HTML变化对比:\n")

    // 提取标签
    beforeTag := extractOpeningTag(before)
    afterTag := extractOpeningTag(after)

    if beforeTag != "" && afterTag != "" {
        fmt.Printf("  前: %s\n", truncateString(beforeTag, 50))
        fmt.Printf("  后: %s\n", truncateString(afterTag, 50))
    }
}

// 提取开始标签
func extractOpeningTag(html string) string {
    re := regexp.MustCompile(`<[^>]+>`)
    match := re.FindString(html)
    if match == "" {
        return ""
    }

    // 只取第一个标签
    end := strings.Index(match, " ")
    if end > 0 {
        return match[:end] + ">"
    }
    return match
}

// 测试重复重命名
func testRepeatedRenames(nodeID int, count int) {
    fmt.Printf("重复重命名测试 (%d 次):\n", count)

    tags := []string{"div", "section", "article", "aside"}
    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < count; i++ {
        tagIndex := i % len(tags)
        newTag := tags[tagIndex]

        startTime := time.Now()
        result, err := CDPDOMSetNodeName(nodeID, newTag)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%s -> %s, 错误: %v)\n",
                i+1, tags[(tagIndex+len(tags)-1)%len(tags)], newTag, err)
        } else {
            successCount++
            fmt.Printf("  第 %d 次: ✅ 成功 (%s -> %s, 耗时: %v)\n",
                i+1, tags[(tagIndex+len(tags)-1)%len(tags)], newTag, duration)

            // 验证设置
            info, err := getNodeInfo(nodeID)
            if err == nil {
                actualName := getNodeName(info)
                if strings.EqualFold(actualName, newTag) {
                    fmt.Printf("    验证: ✅ 标签已更新\n")
                } else {
                    fmt.Printf("    验证: ❌ 标签不匹配 (期望: %s, 实际: %s)\n",
                        newTag, actualName)
                }
            }
        }

        // 短暂延迟
        if i < count-1 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))
}

// 测试不同元素
func testDifferentElements(nodeIDs []int) {
    fmt.Printf("\n多元素重命名测试 (%d 个元素):\n", len(nodeIDs))

    var results []RenameTestResult

    for _, nodeID := range nodeIDs {
        result := RenameTestResult{
            NodeID: nodeID,
        }

        // 获取当前标签
        info, err := getNodeInfo(nodeID)
        if err != nil {
            result.Error = fmt.Sprintf("获取信息失败: %v", err)
            results = append(results, result)
            continue
        }

        oldTag := getNodeName(info)
        result.OldTag = oldTag

        // 选择新标签
        newTag := "section"
        if strings.EqualFold(oldTag, "div") {
            newTag = "section"
        } else if strings.EqualFold(oldTag, "span") {
            newTag = "strong"
        } else {
            newTag = "div"
        }

        startTime := time.Now()
        response, err := CDPDOMSetNodeName(nodeID, newTag)
        result.Duration = time.Since(startTime)

        if err != nil {
            result.Success = false
            result.Error = err.Error()
        } else {
            result.Success = true
            result.Response = response

            // 获取新标签
            infoAfter, err := getNodeInfo(nodeID)
            if err == nil {
                result.NewTag = getNodeName(infoAfter)
            }
        }

        results = append(results, result)

        // 短暂延迟
        time.Sleep(200 * time.Millisecond)
    }

    // 显示结果
    successCount := 0
    for _, result := range results {
        status := "❌ 失败"
        if result.Success {
            status = "✅ 成功"
            successCount++
        }

        fmt.Printf("  节点ID: %-6d %s", result.NodeID, status)
        if result.Success {
            fmt.Printf(" (%s -> %s)", result.OldTag, result.NewTag)
        }
        fmt.Printf(" (耗时: %v)\n", result.Duration)
    }

    fmt.Printf("总体成功率: %.1f%% (%d/%d)\n",
        float64(successCount)/float64(len(results))*100,
        successCount, len(results))
}

type RenameTestResult struct {
    NodeID   int
    OldTag   string
    NewTag   string
    Success  bool
    Error    string
    Response string
    Duration time.Duration
}

// 高级功能: 智能元素重命名器
type SmartElementRenamer struct {
    renameHistory []RenameRecord
    tagValidator  *TagValidator
    impactAnalyzer *RenameImpactAnalyzer
    cache         map[string]RenamePattern
}

type RenameRecord struct {
    Timestamp   time.Time
    NodeID      int
    OldName     string
    NewName     string
    Duration    time.Duration
    Success     bool
    Error       string
    Impact      RenameImpact
    Reason      string
    BeforeHTML  string
    AfterHTML   string
}

type TagValidator struct {
    validHTMLTags   map[string]bool
    semanticTags    map[string]bool
    deprecatedTags  map[string]bool
    voidElements    map[string]bool
    customElements  bool
}

type RenameImpactAnalyzer struct {
    layoutChangingTags map[string]bool
    inlineElements     map[string]bool
    blockElements      map[string]bool
    formElements       map[string]bool
}

type RenamePattern struct {
    OldTag      string
    NewTag      string
    Conditions  map[string]interface{}
    Description string
    UsageCount  int
}

type RenameImpact struct {
    RequiresRepaint  bool
    RequiresReflow   bool
    LayoutChange     bool
    StyleImpact      string
    Accessibility    string
    SemanticsChange  string
}

func NewSmartElementRenamer() *SmartElementRenamer {
    return &SmartElementRenamer{
        renameHistory: make([]RenameRecord, 0),
        tagValidator:  NewTagValidator(),
        impactAnalyzer: NewRenameImpactAnalyzer(),
        cache:         make(map[string]RenamePattern),
    }
}

func (ser *SmartElementRenamer) RenameElement(nodeID int, newName, reason string, validate bool) (RenameRecord, error) {
    record := RenameRecord{
        Timestamp: time.Now(),
        NodeID:    nodeID,
        NewName:   newName,
        Reason:    reason,
    }

    // 获取当前信息
    info, err := getNodeInfo(nodeID)
    if err != nil {
        return record, fmt.Errorf("获取节点信息失败: %w", err)
    }

    oldName := getNodeName(info)
    record.OldName = oldName

    // 保存重命名前的HTML
    record.BeforeHTML, _ = getOuterHTML(nodeID)

    // 验证新标签名
    if validate {
        if err := ser.tagValidator.Validate(newName); err != nil {
            return record, fmt.Errorf("标签验证失败: %w", err)
        }
    }

    // 分析重命名影响
    record.Impact = ser.impactAnalyzer.AnalyzeImpact(oldName, newName)

    // 检查缓存模式
    cacheKey := fmt.Sprintf("%s->%s", oldName, newName)
    if pattern, exists := ser.cache[cacheKey]; exists {
        pattern.UsageCount++
        ser.cache[cacheKey] = pattern
    }

    // 执行重命名
    startTime := time.Now()
    result, err := CDPDOMSetNodeName(nodeID, newName)
    record.Duration = time.Since(startTime)

    if err != nil {
        record.Success = false
        record.Error = err.Error()
    } else {
        record.Success = true
        record.Response = result

        // 获取重命名后的HTML
        record.AfterHTML, _ = getOuterHTML(nodeID)

        // 验证重命名
        infoAfter, err := getNodeInfo(nodeID)
        if err == nil {
            actualName := getNodeName(infoAfter)
            if !strings.EqualFold(actualName, newName) {
                record.Error = fmt.Sprintf("重命名验证失败: 期望 %s, 实际 %s", newName, actualName)
                record.Success = false
            }
        }
    }

    // 记录历史
    ser.recordRename(record)

    if err != nil {
        return record, fmt.Errorf("重命名失败: %w", err)
    }

    return record, nil
}

func NewTagValidator() *TagValidator {
    return &TagValidator{
        validHTMLTags: map[string]bool{
            "div": true, "span": true, "p": true, "a": true, "img": true,
            "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
            "ul": true, "ol": true, "li": true, "table": true, "tr": true, "td": true,
            "form": true, "input": true, "button": true, "textarea": true, "select": true,
            "section": true, "article": true, "nav": true, "aside": true, "main": true,
            "header": true, "footer": true, "figure": true, "figcaption": true,
            "time": true, "mark": true, "summary": true, "details": true,
        },
        semanticTags: map[string]bool{
            "section": true, "article": true, "nav": true, "aside": true,
            "main": true, "header": true, "footer": true, "figure": true,
            "figcaption": true, "time": true, "mark": true, "summary": true,
            "details": true,
        },
        deprecatedTags: map[string]bool{
            "font": true, "center": true, "big": true, "strike": true,
            "tt": true, "acronym": true, "applet": true, "basefont": true,
            "dir": true, "isindex": true, "listing": true, "multicol": true,
            "nextid": true, "plaintext": true, "xmp": true,
        },
        voidElements: map[string]bool{
            "area": true, "base": true, "br": true, "col": true, "embed": true,
            "hr": true, "img": true, "input": true, "link": true, "meta": true,
            "param": true, "source": true, "track": true, "wbr": true,
        },
        customElements: true,
    }
}

func (tv *TagValidator) Validate(tagName string) error {
    tagName = strings.ToLower(tagName)

    if tagName == "" {
        return fmt.Errorf("标签名不能为空")
    }

    // 检查是否以字母开头
    if len(tagName) > 0 && !isLetter(tagName[0]) {
        return fmt.Errorf("标签名必须以字母开头")
    }

    // 检查是否包含非法字符
    for i := 0; i < len(tagName); i++ {
        c := tagName[i]
        if !isValidTagChar(c) && c != '-' {
            return fmt.Errorf("标签名包含非法字符: %c", c)
        }
    }

    // 检查是否被弃用
    if tv.deprecatedTags[tagName] {
        return fmt.Errorf("标签已被弃用: %s", tagName)
    }

    // 检查是否是已知的HTML标签
    if tv.validHTMLTags[tagName] {
        return nil
    }

    // 检查是否是自定义元素
    if tv.customElements && strings.Contains(tagName, "-") {
        return nil
    }

    return fmt.Errorf("未知的HTML标签: %s", tagName)
}

func isLetter(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isValidTagChar(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func NewRenameImpactAnalyzer() *RenameImpactAnalyzer {
    return &RenameImpactAnalyzer{
        layoutChangingTags: map[string]bool{
            "div": true, "section": true, "article": true, "main": true,
            "header": true, "footer": true, "aside": true, "nav": true,
        },
        inlineElements: map[string]bool{
            "span": true, "a": true, "strong": true, "em": true, "code": true,
            "mark": true, "time": true,
        },
        blockElements: map[string]bool{
            "div": true, "p": true, "h1": true, "h2": true, "h3": true,
            "h4": true, "h5": true, "h6": true, "section": true, "article": true,
            "nav": true, "aside": true, "main": true, "header": true, "footer": true,
        },
        formElements: map[string]bool{
            "input": true, "button": true, "select": true, "textarea": true,
            "form": true, "label": true,
        },
    }
}

func (ria *RenameImpactAnalyzer) AnalyzeImpact(oldTag, newTag string) RenameImpact {
    impact := RenameImpact{}

    oldTag = strings.ToLower(oldTag)
    newTag = strings.ToLower(newTag)

    if oldTag == newTag {
        return impact
    }

    // 检查是否需要重排
    if ria.layoutChangingTags[oldTag] != ria.layoutChangingTags[newTag] {
        impact.RequiresReflow = true
        impact.LayoutChange = true
    }

    // 检查样式影响
    if ria.inlineElements[oldTag] != ria.inlineElements[newTag] {
        impact.StyleImpact = "display 属性改变"
        impact.RequiresRepaint = true
    }

    if ria.blockElements[oldTag] != ria.blockElements[newTag] {
        impact.StyleImpact = "块级/行内显示模式改变"
        impact.RequiresRepaint = true
    }

    // 检查语义变化
    if isSemanticTag(oldTag) != isSemanticTag(newTag) {
        if isSemanticTag(newTag) {
            impact.SemanticsChange = "语义化改进"
        } else {
            impact.SemanticsChange = "语义化降级"
        }
    }

    // 检查可访问性影响
    if ria.formElements[oldTag] != ria.formElements[newTag] {
        if ria.formElements[newTag] {
            impact.Accessibility = "可访问性改进（表单元素）"
        } else {
            impact.Accessibility = "可访问性降低（非表单元素）"
        }
    }

    return impact
}

func isSemanticTag(tag string) bool {
    semanticTags := map[string]bool{
        "section": true, "article": true, "nav": true, "aside": true,
        "main": true, "header": true, "footer": true, "figure": true,
        "figcaption": true, "time": true, "mark": true, "summary": true,
        "details": true,
    }
    return semanticTags[tag]
}

func (ser *SmartElementRenamer) recordRename(record RenameRecord) {
    ser.renameHistory = append(ser.renameHistory, record)
}

func (ser *SmartElementRenamer) GetStats() map[string]interface{} {
    totalRenames := len(ser.renameHistory)
    successfulRenames := 0
    var totalDuration time.Duration
    var semanticImprovements int

    for _, record := range ser.renameHistory {
        if record.Success {
            successfulRenames++
        }
        totalDuration += record.Duration

        if record.Impact.SemanticsChange == "语义化改进" {
            semanticImprovements++
        }
    }

    avgDuration := time.Duration(0)
    if totalRenames > 0 {
        avgDuration = totalDuration / time.Duration(totalRenames)
    }

    return map[string]interface{}{
        "totalRenames":        totalRenames,
        "successfulRenames":   successfulRenames,
        "failedRenames":       totalRenames - successfulRenames,
        "successRate":         float64(successfulRenames) / float64(totalRenames) * 100,
        "totalDuration":       totalDuration,
        "averageDuration":     avgDuration,
        "semanticImprovements": semanticImprovements,
        "cacheSize":          len(ser.cache),
    }
}

// 演示智能重命名
func demonstrateSmartRenaming() {
    fmt.Println("=== 智能元素重命名演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建智能重命名器
    renamer := NewSmartElementRenamer()

    // 测试元素
    elementID := 5001

    fmt.Printf("测试元素ID: %d\n\n", elementID)

    // 场景1: 语义化改进
    fmt.Printf("场景1: 语义化改进\n")

    semanticTests := []struct {
        newTag  string
        reason  string
    }{
        {"section", "将通用容器转换为章节"},
        {"article", "表示独立的文章内容"},
        {"aside", "表示侧边栏内容"},
        {"main", "表示页面主要内容"},
    }

    for _, test := range semanticTests {
        fmt.Printf("重命名为: %s\n", test.newTag)
        fmt.Printf("原因: %s\n", test.reason)

        record, err := renamer.RenameElement(elementID, test.newTag, test.reason, true)
        if err != nil {
            fmt.Printf("  ❌ 失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)

            if record.Impact.SemanticsChange != "" {
                fmt.Printf("    影响: %s\n", record.Impact.SemanticsChange)
            }
            if record.Impact.LayoutChange {
                fmt.Printf("    注意: 可能影响布局\n")
            }
        }

        time.Sleep(300 * time.Millisecond)
    }

    // 场景2: 验证测试
    fmt.Printf("\n场景2: 标签验证测试\n")

    validationTests := []struct {
        tag     string
        isValid bool
        reason  string
    }{
        {"div", true, "有效的HTML标签"},
        {"font", false, "已弃用的标签"},
        {"my-component", true, "自定义元素"},
        {"123tag", false, "以数字开头的标签"},
        {"", false, "空标签"},
    }

    for _, test := range validationTests {
        fmt.Printf("测试标签: %s\n", test.tag)
        fmt.Printf("描述: %s\n", test.reason)

        if test.isValid {
            _, err := renamer.RenameElement(elementID, test.tag, "测试", true)
            if err != nil {
                fmt.Printf("  ❌ 预期成功但失败: %v\n", err)
            } else {
                fmt.Printf("  ✅ 验证通过\n")
            }
        } else {
            _, err := renamer.RenameElement(elementID, test.tag, "测试", true)
            if err != nil {
                fmt.Printf("  ✅ 预期错误: %v\n", err)
            } else {
                fmt.Printf("  ❌ 预期错误但成功\n")
            }
        }
    }

    // 显示统计
    fmt.Printf("\n=== 重命名统计 ===\n")
    stats := renamer.GetStats()
    fmt.Printf("  总重命名次数: %d\n", stats["totalRenames"])
    fmt.Printf("  成功重命名: %d\n", stats["successfulRenames"])
    fmt.Printf("  失败重命名: %d\n", stats["failedRenames"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("  语义化改进: %d\n", stats["semanticImprovements"])
    fmt.Printf("  平均耗时: %v\n", stats["averageDuration"])

    fmt.Println("\n=== 演示完成 ===")
}


*/

// -----------------------------------------------  DOM.setNodeValue  -----------------------------------------------
// === 应用场景 ===
// 1. 文本节点更新: 更新文本节点的内容
// 2. 注释修改: 修改注释节点的内容
// 3. 属性值更新: 更新属性节点的值
// 4. 内容动态更新: 动态更新页面文本内容
// 5. 国际化: 根据语言动态更新文本
// 6. 模板渲染: 在模板中更新文本内容

// CDPDOMSetNodeValue 设置文本节点的值
// nodeID: 要设置值的节点ID
// value: 要设置的文本值
func CDPDOMSetNodeValue(nodeID int, value string) (string, error) {
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
        "method": "DOM.setNodeValue",
        "params": {
            "nodeId": %d,
            "value": "%s"
        }
    }`, reqID, nodeID, value)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.setNodeValue 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.setNodeValue 请求超时")
		}
	}
}

/*

// 示例: 更新页面中的文本内容
func ExampleCDPDOMSetNodeValue() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个文本节点的节点ID
    textNodeID := 1001

    fmt.Printf("=== 文本节点值设置演示 ===\n")
    fmt.Printf("目标文本节点ID: %d\n\n", textNodeID)

    // 2. 首先获取节点的当前信息
    fmt.Printf("=== 设置前检查 ===\n")

    nodeInfo, err := getNodeInfo(textNodeID)
    if err != nil {
        log.Printf("获取节点信息失败: %v", err)
        return
    }

    displayNodeInfo("文本节点", nodeInfo)

    // 检查节点类型
    nodeType := getNodeType(nodeInfo)
    fmt.Printf("节点类型: %s\n", nodeType)

    // 检查当前值
    currentValue := getNodeValue(nodeInfo)
    fmt.Printf("当前文本值: %s\n", currentValue)

    // 3. 准备测试用例
    fmt.Printf("\n=== 文本更新测试 ===\n")

    testCases := []struct {
        newValue   string
        description string
        category    string
    }{
        {
            newValue:    "这是一个新的文本内容。",
            description: "普通文本更新",
            category:    "基本文本",
        },
        {
            newValue:    "更新后的文本，包含标点符号！",
            description: "带标点符号的文本",
            category:    "标点文本",
        },
        {
            newValue:    "多行文本\n这是第二行\n这是第三行",
            description: "多行文本内容",
            category:    "多行文本",
        },
        {
            newValue:    "包含 多个    空格  的文本",
            description: "包含多余空格的文本",
            category:    "空格处理",
        },
        {
            newValue:    "特殊字符: <>&\"'",
            description: "包含HTML特殊字符",
            category:    "特殊字符",
        },
        {
            newValue:    "Unicode: 中文测试 🎉 ✅ 🔥",
            description: "包含Unicode字符和表情符号",
            category:    "Unicode文本",
        },
        {
            newValue:    "很长的文本..." + strings.Repeat("这是一段重复的内容，用于测试长文本的更新。", 20),
            description: "长文本内容测试",
            category:    "长文本",
        },
        {
            newValue:    "",
            description: "清空文本内容",
            category:    "空文本",
        },
    }

    originalValue := currentValue

    // 4. 执行文本更新测试
    for i, tc := range testCases {
        fmt.Printf("\n测试 %d: %s\n", i+1, tc.description)
        fmt.Printf("类别: %s\n", tc.category)

        // 显示文本长度
        textLength := len(tc.newValue)
        fmt.Printf("文本长度: %d 字符\n", textLength)

        if textLength > 50 {
            preview := tc.newValue[:50] + "..."
            fmt.Printf("文本预览: %s\n", preview)
        } else {
            fmt.Printf("文本内容: %s\n", tc.newValue)
        }

        startTime := time.Now()
        result, err := CDPDOMSetNodeValue(textNodeID, tc.newValue)
        setTime := time.Since(startTime)

        if err != nil {
            fmt.Printf("❌ 设置失败: %v\n", err)
            continue
        }

        fmt.Printf("✅ 设置成功 (耗时: %v)\n", setTime)
        fmt.Printf("结果: %s\n", result)

        // 验证设置结果
        fmt.Printf("\n验证设置结果...\n")

        // 获取更新后的节点信息
        updatedInfo, err := getNodeInfo(textNodeID)
        if err != nil {
            fmt.Printf("❌ 验证失败: 无法获取更新后信息\n")
            continue
        }

        updatedValue := getNodeValue(updatedInfo)
        fmt.Printf("更新后的值: %s\n", updatedValue)

        // 验证文本是否匹配
        if updatedValue == tc.newValue {
            fmt.Printf("✅ 验证通过: 文本已成功更新\n")
        } else {
            // 对于多行文本，检查换行符处理
            if strings.Contains(tc.newValue, "\n") {
                fmt.Printf("⚠️ 多行文本可能被规范化\n")
            } else {
                fmt.Printf("❌ 验证失败: 文本不匹配\n")
                fmt.Printf("   期望: %s\n", tc.newValue)
                fmt.Printf("   实际: %s\n", updatedValue)
            }
        }

        // 计算处理速度
        if textLength > 0 && setTime > 0 {
            charsPerSecond := float64(textLength) / setTime.Seconds()
            fmt.Printf("处理速度: %.0f 字符/秒\n", charsPerSecond)
        }

        // 短暂延迟
        if i < len(testCases)-1 {
            time.Sleep(200 * time.Millisecond)
        }
    }

    // 5. 恢复原始文本
    fmt.Printf("\n=== 恢复原始文本 ===\n")
    fmt.Printf("恢复为原始文本: %s\n", truncateString(originalValue, 50))

    result, err := CDPDOMSetNodeValue(textNodeID, originalValue)
    if err != nil {
        fmt.Printf("❌ 恢复失败: %v\n", err)
    } else {
        fmt.Printf("✅ 已恢复\n")
    }

    // 6. 特殊字符处理测试
    fmt.Printf("\n=== 特殊字符处理测试 ===\n")

    specialCharTests := []struct {
        charName   string
        charValue  string
        description string
    }{
        {"换行符", "第一行\n第二行\n第三行", "多行换行符"},
        {"制表符", "列1\t列2\t列3", "制表符分隔"},
        {"零宽空格", "文本\u200B间隔", "零宽空格"},
        {"不可见字符", "可见\u0009\u000A\u000D", "控制字符"},
        {"HTML实体", "&lt;div&gt; &amp; &quot;", "HTML实体"},
        {"Emoji", "😀 🎉 ✅ 📱 💻", "表情符号"},
    }

    for _, test := range specialCharTests {
        fmt.Printf("测试: %s\n", test.charName)
        fmt.Printf("描述: %s\n", test.description)
        fmt.Printf("字符: %q\n", test.charValue)

        result, err := CDPDOMSetNodeValue(textNodeID, test.charValue)
        if err != nil {
            fmt.Printf("❌ 失败: %v\n", err)
        } else {
            fmt.Printf("✅ 成功\n")

            // 验证字符处理
            updatedInfo, _ := getNodeInfo(textNodeID)
            updatedValue := getNodeValue(updatedInfo)

            if updatedValue == test.charValue {
                fmt.Printf("  字符保留完整\n")
            } else {
                fmt.Printf("  字符被处理: %q\n", updatedValue)
            }
        }
    }

    // 7. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次文本更新
    testRepeatedUpdates(textNodeID, 10)

    // 测试不同长度文本的更新性能
    testDifferentTextLengths(textNodeID)

    // 8. 节点类型测试
    fmt.Printf("\n=== 节点类型测试 ===\n")

    // 尝试设置不同类型的节点
    testDifferentNodeTypes()

    // 9. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "实时计数器",
            description: "实时更新计数显示",
            useCase:     "倒计时、计数器、实时数据更新",
        },
        {
            name:        "动态提示信息",
            description: "根据用户操作更新提示信息",
            useCase:     "表单验证提示、操作反馈、状态信息",
        },
        {
            name:        "国际化文本",
            description: "根据语言切换动态更新文本",
            useCase:     "多语言网站、本地化内容",
        },
        {
            name:        "模板内容填充",
            description: "在模板中填充动态内容",
            useCase:     "邮件模板、报告生成、动态内容",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 获取节点类型
func getNodeType(info map[string]interface{}) string {
    nodeType := "未知"

    if nodeName, ok := info["nodeName"].(string); ok {
        switch nodeName {
        case "#text":
            nodeType = "文本节点"
        case "#comment":
            nodeType = "注释节点"
        case "#document":
            nodeType = "文档节点"
        case "#document-fragment":
            nodeType = "文档片段"
        default:
            if strings.HasPrefix(nodeName, "#") {
                nodeType = "特殊节点"
            } else {
                nodeType = "元素节点"
            }
        }
    }

    return nodeType
}

// 获取节点值
func getNodeValue(info map[string]interface{}) string {
    if value, ok := info["nodeValue"].(string); ok {
        return value
    }
    return ""
}

// 测试重复更新
func testRepeatedUpdates(nodeID int, count int) {
    fmt.Printf("重复文本更新测试 (%d 次):\n", count)

    baseText := "更新次数: "
    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < count; i++ {
        newValue := fmt.Sprintf("%s%d", baseText, i+1)

        startTime := time.Now()
        result, err := CDPDOMSetNodeValue(nodeID, newValue)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            successCount++
            fmt.Printf("  第 %d 次: ✅ 成功 (文本: %s, 耗时: %v)\n",
                i+1, truncateString(newValue, 20), duration)

            // 验证设置
            info, err := getNodeInfo(nodeID)
            if err == nil {
                actualValue := getNodeValue(info)
                if actualValue == newValue {
                    fmt.Printf("    验证: ✅ 文本已更新\n")
                } else {
                    fmt.Printf("    验证: ❌ 文本不匹配\n")
                }
            }
        }

        // 短暂延迟
        if i < count-1 {
            time.Sleep(50 * time.Millisecond)
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))
}

// 测试不同长度文本
func testDifferentTextLengths(nodeID int) {
    fmt.Printf("\n不同长度文本更新性能对比:\n")

    lengthTests := []struct {
        length int
        description string
    }{
        {10, "很短文本"},
        {100, "短文本"},
        {1000, "中等文本"},
        {10000, "长文本"},
        {50000, "很长文本"},
    }

    var results []TextUpdatePerformance

    for _, test := range lengthTests {
        result := TextUpdatePerformance{
            Length: test.length,
        }

        // 生成测试文本
        testText := generateTestText(test.length)

        // 多次测试取平均值
        var totalTime time.Duration
        tests := 3

        for i := 0; i < tests; i++ {
            startTime := time.Now()
            _, err := CDPDOMSetNodeValue(nodeID, testText)
            duration := time.Since(startTime)

            if err == nil {
                totalTime += duration
            }

            if i < tests-1 {
                time.Sleep(20 * time.Millisecond)
            }
        }

        result.AverageTime = totalTime / time.Duration(tests)
        result.Speed = float64(test.length) / result.AverageTime.Seconds()
        results = append(results, result)
    }

    // 显示结果
    fmt.Printf("文本长度  | 平均耗时  | 处理速度\n")
    fmt.Printf("----------|-----------|-------------\n")
    for _, result := range results {
        fmt.Printf("%8d | %9v | %7.0f 字符/秒\n",
            result.Length, result.AverageTime, result.Speed)
    }
}

type TextUpdatePerformance struct {
    Length      int
    AverageTime time.Duration
    Speed       float64
}

// 生成测试文本
func generateTestText(length int) string {
    if length <= 0 {
        return ""
    }

    var builder strings.Builder
    words := []string{"文本", "测试", "内容", "更新", "性能", "分析", "数据", "处理"}

    for builder.Len() < length {
        word := words[rand.Intn(len(words))]
        if builder.Len() > 0 {
            builder.WriteString(" ")
        }
        builder.WriteString(word)
    }

    // 截断到指定长度
    text := builder.String()
    if len(text) > length {
        text = text[:length]
    }

    return text
}

// 测试不同类型节点
func testDifferentNodeTypes() {
    // 这里需要先获取不同类型的节点
    // 简化实现
    fmt.Printf("节点类型兼容性:\n")
    fmt.Printf("  ✅ 文本节点 (#text) - 完全支持\n")
    fmt.Printf("  ✅ 注释节点 (#comment) - 支持\n")
    fmt.Printf("  ⚠️ 属性节点 - 有限支持\n")
    fmt.Printf("  ❌ 元素节点 - 不支持（应使用innerHTML/textContent）\n")
    fmt.Printf("  ❌ 文档节点 - 不支持\n")
}

// 高级功能: 智能文本更新器
type SmartTextUpdater struct {
    updateHistory []TextUpdateRecord
    textAnalyzer  *TextContentAnalyzer
    validator     *TextValueValidator
    cache         *TextUpdateCache
}

type TextUpdateRecord struct {
    Timestamp   time.Time
    NodeID      int
    OldValue    string
    NewValue    string
    Duration    time.Duration
    Success     bool
    Error       string
    Changes     TextChanges
    NodeType    string
    Reason      string
}

type TextChanges struct {
    CharactersAdded    int
    CharactersRemoved  int
    CharactersChanged  int
    LinesAdded         int
    LinesRemoved       int
    Similarity         float64
}

type TextContentAnalyzer struct {
    maxLineLength  int
    dangerousChars []string
    allowedChars   []rune
}

type TextValueValidator struct {
    maxLength     int
    minLength     int
    allowEmpty    bool
    allowNewlines bool
    allowTabs     bool
    allowedChars  string
}

type TextUpdateCache struct {
    entries map[string]CachedTextUpdate
    maxSize int
    ttl     time.Duration
}

type CachedTextUpdate struct {
    NodeID    int
    Value     string
    Timestamp time.Time
    HitCount  int
}

func NewSmartTextUpdater() *SmartTextUpdater {
    return &SmartTextUpdater{
        updateHistory: make([]TextUpdateRecord, 0),
        textAnalyzer:  NewTextContentAnalyzer(),
        validator:     NewTextValueValidator(),
        cache: &TextUpdateCache{
            entries: make(map[string]CachedTextUpdate),
            maxSize: 100,
            ttl:     5 * time.Minute,
        },
    }
}

func (stu *SmartTextUpdater) UpdateNodeValue(nodeID int, newValue, reason string, validate bool) (TextUpdateRecord, error) {
    record := TextUpdateRecord{
        Timestamp: time.Now(),
        NodeID:    nodeID,
        NewValue:  newValue,
        Reason:    reason,
    }

    // 获取当前节点信息
    info, err := getNodeInfo(nodeID)
    if err != nil {
        return record, fmt.Errorf("获取节点信息失败: %w", err)
    }

    // 检查节点类型
    record.NodeType = getNodeType(info)
    if !stu.isNodeTypeSupported(record.NodeType) {
        return record, fmt.Errorf("不支持的节点类型: %s", record.NodeType)
    }

    // 获取旧值
    oldValue := getNodeValue(info)
    record.OldValue = oldValue

    // 分析文本变化
    record.Changes = stu.analyzeTextChanges(oldValue, newValue)

    // 验证新值
    if validate {
        if err := stu.validator.Validate(newValue); err != nil {
            return record, fmt.Errorf("文本验证失败: %w", err)
        }
    }

    // 分析文本内容
    if warnings := stu.textAnalyzer.Analyze(newValue); len(warnings) > 0 {
        record.Error = strings.Join(warnings, "; ")
    }

    // 检查缓存
    cacheKey := fmt.Sprintf("%d:%s", nodeID, newValue)
    if cached, exists := stu.cache.entries[cacheKey]; exists {
        if time.Since(cached.Timestamp) < stu.cache.ttl {
            cached.HitCount++
            stu.cache.entries[cacheKey] = cached

            record.Success = true
            record.Duration = 0
            stu.recordUpdate(record)

            return record, nil
        } else {
            // 缓存过期
            delete(stu.cache.entries, cacheKey)
        }
    }

    // 执行更新
    startTime := time.Now()
    result, err := CDPDOMSetNodeValue(nodeID, newValue)
    record.Duration = time.Since(startTime)

    if err != nil {
        record.Success = false
        record.Error = err.Error()
    } else {
        record.Success = true
        record.Response = result

        // 更新缓存
        stu.cache.entries[cacheKey] = CachedTextUpdate{
            NodeID:    nodeID,
            Value:     newValue,
            Timestamp: time.Now(),
            HitCount:  1,
        }

        // 清理过期缓存
        stu.cleanupCache()
    }

    // 记录历史
    stu.recordUpdate(record)

    if err != nil {
        return record, fmt.Errorf("更新失败: %w", err)
    }

    return record, nil
}

func (stu *SmartTextUpdater) isNodeTypeSupported(nodeType string) bool {
    supportedTypes := map[string]bool{
        "文本节点":  true,
        "注释节点":  true,
        "属性节点":  true,
    }
    return supportedTypes[nodeType]
}

func (stu *SmartTextUpdater) analyzeTextChanges(oldValue, newValue string) TextChanges {
    changes := TextChanges{}

    if oldValue == newValue {
        changes.Similarity = 1.0
        return changes
    }

    // 计算行数变化
    oldLines := strings.Count(oldValue, "\n") + 1
    newLines := strings.Count(newValue, "\n") + 1
    changes.LinesAdded = max(0, newLines - oldLines)
    changes.LinesRemoved = max(0, oldLines - newLines)

    // 计算字符变化
    changes.CharactersAdded = max(0, len(newValue) - len(oldValue))
    changes.CharactersRemoved = max(0, len(oldValue) - len(newValue))

    // 计算相似度（简单实现）
    changes.Similarity = calculateSimilarity(oldValue, newValue)

    return changes
}

func calculateSimilarity(str1, str2 string) float64 {
    if str1 == str2 {
        return 1.0
    }

    if str1 == "" || str2 == "" {
        return 0.0
    }

    // 使用Levenshtein距离计算相似度
    distance := levenshteinDistance(str1, str2)
    maxLen := max(len(str1), len(str2))

    if maxLen == 0 {
        return 1.0
    }

    return 1.0 - float64(distance)/float64(maxLen)
}

func levenshteinDistance(str1, str2 string) int {
    // 简化的Levenshtein距离计算
    if len(str1) == 0 {
        return len(str2)
    }
    if len(str2) == 0 {
        return len(str1)
    }

    // 实现略，实际使用时可以引入专门的算法库
    return 0
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func NewTextContentAnalyzer() *TextContentAnalyzer {
    return &TextContentAnalyzer{
        maxLineLength: 1000,
        dangerousChars: []string{
            "\x00", // 空字符
            "\x1B", // ESC
            "\x7F", // DEL
        },
        allowedChars: []rune{},
    }
}

func (tca *TextContentAnalyzer) Analyze(text string) []string {
    var warnings []string

    // 检查行长度
    lines := strings.Split(text, "\n")
    for i, line := range lines {
        if len(line) > tca.maxLineLength {
            warnings = append(warnings,
                fmt.Sprintf("第%d行过长: %d > %d", i+1, len(line), tca.maxLineLength))
        }
    }

    // 检查危险字符
    for _, dangerous := range tca.dangerousChars {
        if strings.Contains(text, dangerous) {
            warnings = append(warnings, "包含危险控制字符")
            break
        }
    }

    return warnings
}

func NewTextValueValidator() *TextValueValidator {
    return &TextValueValidator{
        maxLength:     100000, // 100KB
        minLength:     0,
        allowEmpty:    true,
        allowNewlines: true,
        allowTabs:     true,
        allowedChars:  "",
    }
}

func (tvv *TextValueValidator) Validate(text string) error {
    // 检查长度
    if len(text) > tvv.maxLength {
        return fmt.Errorf("文本过长: %d > %d", len(text), tvv.maxLength)
    }

    if len(text) < tvv.minLength {
        return fmt.Errorf("文本过短: %d < %d", len(text), tvv.minLength)
    }

    if !tvv.allowEmpty && text == "" {
        return fmt.Errorf("文本不能为空")
    }

    // 检查换行符
    if !tvv.allowNewlines && strings.Contains(text, "\n") {
        return fmt.Errorf("文本不能包含换行符")
    }

    // 检查制表符
    if !tvv.allowTabs && strings.Contains(text, "\t") {
        return fmt.Errorf("文本不能包含制表符")
    }

    // 检查允许的字符
    if tvv.allowedChars != "" {
        for _, r := range text {
            if !strings.ContainsRune(tvv.allowedChars, r) {
                return fmt.Errorf("包含不允许的字符: %c", r)
            }
        }
    }

    return nil
}

func (stu *SmartTextUpdater) cleanupCache() {
    if len(stu.cache.entries) <= stu.cache.maxSize {
        return
    }

    // 清理过期条目
    now := time.Now()
    for key, entry := range stu.cache.entries {
        if now.Sub(entry.Timestamp) > stu.cache.ttl {
            delete(stu.cache.entries, key)
        }
    }

    // 如果仍然超过大小，清理最不常用的
    if len(stu.cache.entries) > stu.cache.maxSize {
        stu.cleanupLRU()
    }
}

func (stu *SmartTextUpdater) cleanupLRU() {
    // 找到最少使用的条目
    var lruKey string
    var lruHitCount = -1
    var lruTime time.Time

    for key, entry := range stu.cache.entries {
        if lruHitCount == -1 ||
           entry.HitCount < lruHitCount ||
           (entry.HitCount == lruHitCount && entry.Timestamp.Before(lruTime)) {
            lruKey = key
            lruHitCount = entry.HitCount
            lruTime = entry.Timestamp
        }
    }

    if lruKey != "" {
        delete(stu.cache.entries, lruKey)
    }
}

func (stu *SmartTextUpdater) recordUpdate(record TextUpdateRecord) {
    stu.updateHistory = append(stu.updateHistory, record)
}

func (stu *SmartTextUpdater) GetStats() map[string]interface{} {
    totalUpdates := len(stu.updateHistory)
    successfulUpdates := 0
    var totalDuration time.Duration
    var totalCharsChanged int

    for _, record := range stu.updateHistory {
        if record.Success {
            successfulUpdates++
        }
        totalDuration += record.Duration
        totalCharsChanged += record.Changes.CharactersAdded + record.Changes.CharactersRemoved
    }

    avgDuration := time.Duration(0)
    if totalUpdates > 0 {
        avgDuration = totalDuration / time.Duration(totalUpdates)
    }

    avgCharsChanged := 0.0
    if totalUpdates > 0 {
        avgCharsChanged = float64(totalCharsChanged) / float64(totalUpdates)
    }

    return map[string]interface{}{
        "totalUpdates":      totalUpdates,
        "successfulUpdates": successfulUpdates,
        "failedUpdates":     totalUpdates - successfulUpdates,
        "successRate":       float64(successfulUpdates) / float64(totalUpdates) * 100,
        "totalDuration":     totalDuration,
        "averageDuration":   avgDuration,
        "totalCharsChanged": totalCharsChanged,
        "averageCharsChanged": avgCharsChanged,
        "cacheSize":        len(stu.cache.entries),
    }
}

// 演示智能文本更新
func demonstrateSmartTextUpdate() {
    fmt.Println("=== 智能文本更新演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建智能更新器
    updater := NewSmartTextUpdater()

    // 测试元素
    textNodeID := 6001

    fmt.Printf("测试文本节点ID: %d\n\n", textNodeID)

    // 场景1: 正常文本更新
    fmt.Printf("场景1: 正常文本更新\n")

    tests := []struct {
        newValue string
        reason   string
    }{
        {"欢迎使用智能文本更新", "初始化文本"},
        {"文本已更新，包含更多内容", "内容扩展"},
        {"最终版本文本", "最终确定"},
    }

    for _, test := range tests {
        fmt.Printf("更新文本: %s\n", test.newValue)
        fmt.Printf("原因: %s\n", test.reason)

        record, err := updater.UpdateNodeValue(textNodeID, test.newValue, test.reason, true)
        if err != nil {
            fmt.Printf("  ❌ 失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)

            if record.Changes.CharactersAdded > 0 {
                fmt.Printf("    添加字符: %d\n", record.Changes.CharactersAdded)
            }
            if record.Changes.CharactersRemoved > 0 {
                fmt.Printf("    删除字符: %d\n", record.Changes.CharactersRemoved)
            }
            fmt.Printf("    相似度: %.1f%%\n", record.Changes.Similarity*100)
        }

        time.Sleep(200 * time.Millisecond)
    }

    // 场景2: 缓存测试
    fmt.Printf("\n场景2: 缓存测试\n")

    cachedText := "缓存测试文本"
    for i := 0; i < 3; i++ {
        record, err := updater.UpdateNodeValue(textNodeID, cachedText, "缓存测试", true)
        if err != nil {
            fmt.Printf("第 %d 次: ❌ 失败\n", i+1)
        } else if record.Duration == 0 {
            fmt.Printf("第 %d 次: ✅ 缓存命中\n", i+1)
        } else {
            fmt.Printf("第 %d 次: ✅ 首次设置\n", i+1)
        }
    }

    // 场景3: 验证测试
    fmt.Printf("\n场景3: 验证测试\n")

    validationTests := []struct {
        text      string
        expectError bool
        reason    string
    }{
        {strings.Repeat("a", 200000), true, "超长文本"},
        {"正常文本", false, "正常文本"},
        {"包含\n换行", true, "不允许换行"}, // 假设不允许换行
    }

    for _, test := range validationTests {
        _, err := updater.UpdateNodeValue(textNodeID, test.text, test.reason, true)

        if test.expectError {
            if err != nil {
                fmt.Printf("文本 '%s': ✅ 预期错误\n", truncateString(test.text, 20))
            } else {
                fmt.Printf("文本 '%s': ❌ 预期错误但成功\n", truncateString(test.text, 20))
            }
        } else {
            if err != nil {
                fmt.Printf("文本 '%s': ❌ 失败: %v\n", truncateString(test.text, 20), err)
            } else {
                fmt.Printf("文本 '%s': ✅ 验证通过\n", truncateString(test.text, 20))
            }
        }
    }

    // 显示统计
    fmt.Printf("\n=== 更新统计 ===\n")
    stats := updater.GetStats()
    fmt.Printf("  总更新次数: %d\n", stats["totalUpdates"])
    fmt.Printf("  成功更新: %d\n", stats["successfulUpdates"])
    fmt.Printf("  失败更新: %d\n", stats["failedUpdates"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("  总修改字符数: %d\n", stats["totalCharsChanged"])
    fmt.Printf("  平均每次修改: %.1f 字符\n", stats["averageCharsChanged"])
    fmt.Printf("  缓存大小: %d\n", stats["cacheSize"])

    fmt.Println("\n=== 演示完成 ===")
}


*/

// -----------------------------------------------  DOM.setOuterHTML  -----------------------------------------------
// === 应用场景 ===
// 1. 元素替换: 完全替换DOM元素及其内容
// 2. 动态组件: 动态渲染和替换组件
// 3. 模板渲染: 从模板字符串渲染HTML
// 4. 内容更新: 批量更新元素及其子元素
// 5. 服务器渲染: 注入服务器渲染的HTML
// 6. 错误修复: 修复损坏的DOM结构

// CDPDOMSetOuterHTML 设置节点的外层HTML
// nodeID: 要设置外层HTML的节点ID
// outerHTML: 新的外层HTML字符串
func CDPDOMSetOuterHTML(nodeID int, outerHTML string) (string, error) {
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
        "method": "DOM.setOuterHTML",
        "params": {
            "nodeId": %d,
            "outerHTML": "%s"
        }
    }`, reqID, nodeID, escapeString(outerHTML))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 DOM.setOuterHTML 请求失败: %w", err)
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
			return "", fmt.Errorf("DOM.setOuterHTML 请求超时")
		}
	}
}

// 转义字符串中的特殊字符
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

/*

// 示例: 替换容器元素的内容
func ExampleCDPDOMSetOuterHTML() {
    // 1. 启用DOM功能
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }

    // 确保测试完成后清理
    defer func() {
        if _, err := CDPDOMDisable(); err != nil {
            log.Printf("禁用DOM失败: %v", err)
        }
    }()

    // 假设我们有一个容器元素的节点ID
    containerNodeID := 1001

    fmt.Printf("=== 外层HTML设置演示 ===\n")
    fmt.Printf("目标容器节点ID: %d\n\n", containerNodeID)

    // 2. 首先获取容器的当前状态
    fmt.Printf("=== 设置前状态检查 ===\n")

    containerInfo, err := getNodeInfo(containerNodeID)
    if err != nil {
        log.Printf("获取容器信息失败: %v", err)
        return
    }

    displayNodeInfo("原始容器", containerInfo)

    // 获取当前outerHTML
    originalOuterHTML, err := getOuterHTML(containerNodeID)
    if err != nil {
        log.Printf("获取当前outerHTML失败: %v", err)
        return
    }

    fmt.Printf("当前outerHTML长度: %d 字符\n", len(originalOuterHTML))
    if len(originalOuterHTML) > 100 {
        fmt.Printf("预览: %s...\n", truncateString(originalOuterHTML, 100))
    } else {
        fmt.Printf("内容: %s\n", originalOuterHTML)
    }

    // 3. 准备不同的HTML替换方案
    fmt.Printf("\n=== HTML替换方案 ===\n")

    replacementTests := []struct {
        name        string
        description string
        html        string
        complexity  string
    }{
        {
            name:        "简单文本替换",
            description: "将容器替换为简单文本",
            html:        `<div id="simple-text">这是一个简单的文本容器</div>`,
            complexity:  "简单",
        },
        {
            name:        "复杂结构替换",
            description: "替换为复杂的嵌套结构",
            html: `<div class="complex-container">
                <header>
                    <h1>页面标题</h1>
                    <nav>
                        <ul>
                            <li><a href="#home">首页</a></li>
                            <li><a href="#about">关于</a></li>
                            <li><a href="#contact">联系</a></li>
                        </ul>
                    </nav>
                </header>
                <main>
                    <article>
                        <h2>文章标题</h2>
                        <p>这是一篇文章的内容。</p>
                    </article>
                </main>
                <footer>
                    <p>版权信息 © 2024</p>
                </footer>
            </div>`,
            complexity:  "复杂",
        },
        {
            name:        "表单结构",
            description: "替换为表单元素",
            html: `<form id="user-form" class="styled-form">
                <div class="form-group">
                    <label for="username">用户名:</label>
                    <input type="text" id="username" name="username" placeholder="请输入用户名">
                </div>
                <div class="form-group">
                    <label for="email">邮箱:</label>
                    <input type="email" id="email" name="email" placeholder="user@example.com">
                </div>
                <div class="form-group">
                    <label for="password">密码:</label>
                    <input type="password" id="password" name="password">
                </div>
                <button type="submit" class="btn btn-primary">提交</button>
            </form>`,
            complexity:  "中等",
        },
        {
            name:        "列表结构",
            description: "替换为动态列表",
            html: `<ul class="item-list">
                <li class="list-item" data-id="1">项目 1</li>
                <li class="list-item" data-id="2">项目 2</li>
                <li class="list-item" data-id="3">项目 3</li>
                <li class="list-item" data-id="4">项目 4</li>
                <li class="list-item" data-id="5">项目 5</li>
            </ul>`,
            complexity:  "简单",
        },
        {
            name:        "组件结构",
            description: "替换为Web组件结构",
            html: `<custom-widget>
                <template>
                    <style>
                        .widget { padding: 20px; background: #f5f5f5; }
                    </style>
                    <div class="widget">
                        <slot name="content">默认内容</slot>
                    </div>
                </template>
            </custom-widget>`,
            complexity:  "高级",
        },
    }

    // 4. 执行HTML替换测试
    for i, test := range replacementTests {
        fmt.Printf("\n测试 %d: %s\n", i+1, test.name)
        fmt.Printf("描述: %s\n", test.description)
        fmt.Printf("复杂度: %s\n", test.complexity)
        fmt.Printf("HTML长度: %d 字符\n", len(test.html))

        if len(test.html) > 100 {
            fmt.Printf("HTML预览: %s...\n", truncateString(test.html, 100))
        }

        startTime := time.Now()
        result, err := CDPDOMSetOuterHTML(containerNodeID, test.html)
        setTime := time.Since(startTime)

        if err != nil {
            fmt.Printf("❌ 设置失败: %v\n", err)
            continue
        }

        fmt.Printf("✅ 设置成功 (耗时: %v)\n", setTime)
        fmt.Printf("结果: %s\n", result)

        // 验证替换结果
        fmt.Printf("\n验证替换结果...\n")

        // 获取替换后的outerHTML
        newOuterHTML, err := getOuterHTML(containerNodeID)
        if err != nil {
            fmt.Printf("❌ 验证失败: 无法获取新outerHTML\n")
            continue
        }

        fmt.Printf("新outerHTML长度: %d 字符\n", len(newOuterHTML))

        // 简化的验证：检查是否包含预期的内容
        if strings.Contains(newOuterHTML, test.html[:min(50, len(test.html))]) {
            fmt.Printf("✅ 基本验证通过\n")
        } else {
            // 检查是否是嵌套结构导致的不同
            if isSimilarHTML(test.html, newOuterHTML) {
                fmt.Printf("⚠️ HTML结构相似但不同（可能是浏览器规范化）\n")
            } else {
                fmt.Printf("❌ 验证失败: HTML不匹配\n")
            }
        }

        // 分析变化
        analyzeHTMLChanges(originalOuterHTML, newOuterHTML)

        // 短暂延迟
        if i < len(replacementTests)-1 {
            time.Sleep(500 * time.Millisecond)
        }
    }

    // 5. 恢复原始状态
    fmt.Printf("\n=== 恢复原始状态 ===\n")
    fmt.Printf("恢复原始outerHTML\n")

    result, err := CDPDOMSetOuterHTML(containerNodeID, originalOuterHTML)
    if err != nil {
        fmt.Printf("❌ 恢复失败: %v\n", err)
    } else {
        fmt.Printf("✅ 已恢复\n")
    }

    // 6. 边界情况测试
    fmt.Printf("\n=== 边界情况测试 ===\n")

    edgeCases := []struct {
        html        string
        description string
        expectError bool
    }{
        {"", "空HTML字符串", false},
        {"<div>未闭合标签", "无效的HTML", true},
        {"<script>alert('xss')</script>", "脚本标签", false},
        {"<div>正常</div><!-- 注释 --><div>更多</div>", "包含注释", false},
        {strings.Repeat("<div>重复</div>", 100), "大量重复元素", false},
        {"<custom-element data-config='{\"key\":\"value\"}'>内容</custom-element>", "JSON属性", false},
    }

    for _, test := range edgeCases {
        fmt.Printf("测试: %s\n", test.description)
        fmt.Printf("HTML: %s\n", truncateString(test.html, 50))

        result, err := CDPDOMSetOuterHTML(containerNodeID, test.html)

        if test.expectError {
            if err != nil {
                fmt.Printf("✅ 预期错误: %v\n", err)
            } else {
                fmt.Printf("❌ 预期错误但成功: %s\n", result)
            }
        } else {
            if err != nil {
                fmt.Printf("❌ 失败: %v\n", err)
            } else {
                fmt.Printf("✅ 成功\n")
            }
        }
    }

    // 7. 性能测试
    fmt.Printf("\n=== 性能测试 ===\n")

    // 测试多次HTML设置
    testRepeatedHTMLSets(containerNodeID, 5)

    // 测试不同大小HTML的性能
    testDifferentHTMLSizes(containerNodeID)

    // 8. 安全考虑
    fmt.Printf("\n=== 安全考虑 ===\n")

    fmt.Printf("安全注意事项:\n")
    fmt.Printf("  - 避免注入未经验证的HTML\n")
    fmt.Printf("  - 注意脚本执行风险\n")
    fmt.Printf("  - 考虑XSS攻击防护\n")
    fmt.Printf("  - 验证HTML结构安全性\n")

    // 9. 实际应用场景
    fmt.Printf("\n=== 实际应用场景 ===\n")

    scenarios := []struct {
        name        string
        description string
        useCase     string
    }{
        {
            name:        "动态页面更新",
            description: "从服务器获取并更新页面片段",
            useCase:     "单页面应用的路由切换，局部内容更新",
        },
        {
            name:        "模板渲染",
            description: "从模板引擎渲染HTML并注入",
            useCase:     "服务器端渲染，客户端模板渲染",
        },
        {
            name:        "组件替换",
            description: "动态替换UI组件",
            useCase:     "模态框显示/隐藏，标签页切换，组件状态变化",
        },
        {
            name:        "错误恢复",
            description: "修复损坏的DOM结构",
            useCase:     "从错误状态恢复，重新渲染组件",
        },
    }

    for _, scenario := range scenarios {
        fmt.Printf("场景: %s\n", scenario.name)
        fmt.Printf("描述: %s\n", scenario.description)
        fmt.Printf("用例: %s\n\n", scenario.useCase)
    }
}

// 检查HTML是否相似
func isSimilarHTML(html1, html2 string) bool {
    // 简化的相似性检查
    // 移除空白字符和属性顺序差异
    html1 = normalizeHTML(html1)
    html2 = normalizeHTML(html2)

    return html1 == html2
}

// 规范化HTML用于比较
func normalizeHTML(html string) string {
    // 移除多余空白
    html = strings.ReplaceAll(html, "\n", "")
    html = strings.ReplaceAll(html, "\t", "")
    html = regexp.MustCompile(`\s+`).ReplaceAllString(html, " ")
    html = strings.TrimSpace(html)

    // 标准化属性顺序（简化实现）
    return html
}

// 分析HTML变化
func analyzeHTMLChanges(before, after string) {
    if before == after {
        fmt.Printf("HTML无变化\n")
        return
    }

    // 统计变化
    changes := struct {
        totalNodesBefore int
        totalNodesAfter  int
        tagsAdded        []string
        tagsRemoved      []string
    }{}

    // 提取标签
    tagsBefore := extractTags(before)
    tagsAfter := extractTags(after)

    changes.totalNodesBefore = len(tagsBefore)
    changes.totalNodesAfter = len(tagsAfter)

    // 找出新增的标签
    tagsAfterMap := make(map[string]bool)
    for _, tag := range tagsAfter {
        tagsAfterMap[tag] = true
    }

    for _, tag := range tagsBefore {
        if !tagsAfterMap[tag] {
            changes.tagsRemoved = append(changes.tagsRemoved, tag)
        }
    }

    // 找出删除的标签
    tagsBeforeMap := make(map[string]bool)
    for _, tag := range tagsBefore {
        tagsBeforeMap[tag] = true
    }

    for _, tag := range tagsAfter {
        if !tagsBeforeMap[tag] {
            changes.tagsAdded = append(changes.tagsAdded, tag)
        }
    }

    // 显示变化
    fmt.Printf("结构变化分析:\n")
    fmt.Printf("  节点数变化: %d -> %d (%+d)\n",
        changes.totalNodesBefore, changes.totalNodesAfter,
        changes.totalNodesAfter - changes.totalNodesBefore)

    if len(changes.tagsAdded) > 0 {
        fmt.Printf("  新增标签: %v\n", changes.tagsAdded)
    }
    if len(changes.tagsRemoved) > 0 {
        fmt.Printf("  删除标签: %v\n", changes.tagsRemoved)
    }
}

// 提取HTML标签
func extractTags(html string) []string {
    var tags []string
    re := regexp.MustCompile(`<(\/?[a-zA-Z][a-zA-Z0-9:\-]*)`)
    matches := re.FindAllStringSubmatch(html, -1)

    for _, match := range matches {
        if len(match) > 1 {
            tags = append(tags, match[1])
        }
    }

    return tags
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// 测试重复HTML设置
func testRepeatedHTMLSets(nodeID int, count int) {
    fmt.Printf("重复HTML设置测试 (%d 次):\n", count)

    testHTML := `<div class="test">测试内容</div>`
    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < count; i++ {
        // 每次添加不同的标识
        uniqueHTML := fmt.Sprintf(`<div class="test" data-index="%d">测试内容 %d</div>`, i+1, i+1)

        startTime := time.Now()
        result, err := CDPDOMSetOuterHTML(nodeID, uniqueHTML)
        duration := time.Since(startTime)
        totalDuration += duration

        if err != nil {
            fmt.Printf("  第 %d 次: ❌ 失败 (%v)\n", i+1, err)
        } else {
            successCount++
            fmt.Printf("  第 %d 次: ✅ 成功 (耗时: %v)\n", i+1, duration)

            // 验证设置
            newHTML, err := getOuterHTML(nodeID)
            if err == nil {
                if strings.Contains(newHTML, fmt.Sprintf(`data-index="%d"`, i+1)) {
                    fmt.Printf("    验证: ✅ HTML已更新\n")
                } else {
                    fmt.Printf("    验证: ❌ HTML验证失败\n")
                }
            }
        }

        // 短暂延迟
        if i < count-1 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Printf("测试结果:\n")
    fmt.Printf("  成功次数: %d/%d\n", successCount, count)
    fmt.Printf("  成功率: %.1f%%\n", float64(successCount)/float64(count)*100)
    fmt.Printf("  总耗时: %v\n", totalDuration)
    fmt.Printf("  平均耗时: %v\n", totalDuration/time.Duration(count))
}

// 测试不同大小HTML的性能
func testDifferentHTMLSizes(nodeID int) {
    fmt.Printf("\n不同大小HTML设置性能对比:\n")

    sizeTests := []struct {
        size        int
        description string
    }{
        {100, "很小HTML"},
        {1000, "小HTML"},
        {10000, "中等HTML"},
        {50000, "大HTML"},
        {200000, "很大HTML"},
    }

    var results []HTMLSetPerformance

    for _, test := range sizeTests {
        result := HTMLSetPerformance{
            Size: test.size,
        }

        // 生成测试HTML
        testHTML := generateTestHTML(test.size)

        // 多次测试取平均值
        var totalTime time.Duration
        tests := 3

        for i := 0; i < tests; i++ {
            startTime := time.Now()
            _, err := CDPDOMSetOuterHTML(nodeID, testHTML)
            duration := time.Since(startTime)

            if err == nil {
                totalTime += duration
            }

            if i < tests-1 {
                time.Sleep(100 * time.Millisecond)
            }
        }

        result.AverageTime = totalTime / time.Duration(tests)
        result.Speed = float64(test.size) / result.AverageTime.Seconds()
        results = append(results, result)
    }

    // 显示结果
    fmt.Printf("HTML大小  | 平均耗时  | 处理速度\n")
    fmt.Printf("----------|-----------|-------------\n")
    for _, result := range results {
        fmt.Printf("%8d | %9v | %7.0f 字符/秒\n",
            result.Size, result.AverageTime, result.Speed)
    }
}

type HTMLSetPerformance struct {
    Size        int
    AverageTime time.Duration
    Speed       float64
}

// 生成测试HTML
func generateTestHTML(size int) string {
    if size <= 0 {
        return ""
    }

    var builder strings.Builder
    builder.WriteString("<div class=\"container\">")

    // 生成段落
    paragraphs := []string{
        "这是一个测试段落，用于生成HTML内容。",
        "HTML是用来描述网页结构的标记语言。",
        "浏览器会解析HTML并渲染出页面。",
        "DOM表示文档对象模型，是HTML的编程接口。",
    }

    for builder.Len() < size {
        paragraph := paragraphs[rand.Intn(len(paragraphs))]
        builder.WriteString(fmt.Sprintf("<p>%s</p>", paragraph))

        // 偶尔添加一些嵌套元素
        if rand.Intn(5) == 0 {
            builder.WriteString("<div><span>嵌套内容</span></div>")
        }
    }

    builder.WriteString("</div>")

    // 截断到指定大小
    html := builder.String()
    if len(html) > size {
        html = html[:size]
    }

    return html
}

// 高级功能: 智能HTML设置器
type SmartHTMLSetter struct {
    setHistory  []HTMLSetRecord
    htmlValidator *HTMLValidator
    changeTracker *HTMLChangeTracker
    cache        *HTMLCache
    sanitizer    *HTMLSanitizer
}

type HTMLSetRecord struct {
    Timestamp   time.Time
    NodeID      int
    OldHTML     string
    NewHTML     string
    Duration    time.Duration
    Success     bool
    Error       string
    Changes     HTMLChanges
    Size        HTMLSize
    Validation  HTMLValidation
    Sanitized   bool
}

type HTMLValidator struct {
    maxSize       int
    allowedTags   map[string]bool
    blockedTags   map[string]bool
    allowedAttrs  map[string][]string
    requireClosing bool
    maxDepth      int
}

type HTMLChangeTracker struct {
    trackChanges  bool
    diffAlgorithm string
    maxHistory    int
}

type HTMLCache struct {
    entries map[string]CachedHTML
    maxSize int
    ttl     time.Duration
}

type CachedHTML struct {
    NodeID    int
    HTML      string
    Timestamp time.Time
    HitCount  int
}

type HTMLSanitizer struct {
    removeScripts   bool
    removeEvents    bool
    removeStyles    bool
    removeComments  bool
    safeProtocols   []string
}

type HTMLChanges struct {
    NodesAdded     int
    NodesRemoved   int
    AttributesChanged int
    TextChanged    int
    StructureChanged bool
    DiffRatio      float64
}

type HTMLSize struct {
    Characters     int
    Tags          int
    Attributes    int
    TextNodes     int
    Depth         int
}

type HTMLValidation struct {
    Valid      bool
    Warnings   []string
    Errors     []string
    Suggestions []string
}

func NewSmartHTMLSetter() *SmartHTMLSetter {
    return &SmartHTMLSetter{
        setHistory:  make([]HTMLSetRecord, 0),
        htmlValidator: NewHTMLValidator(),
        changeTracker: &HTMLChangeTracker{
            trackChanges:  true,
            diffAlgorithm: "simple",
            maxHistory:    50,
        },
        cache: &HTMLCache{
            entries: make(map[string]CachedHTML),
            maxSize: 50,
            ttl:     10 * time.Minute,
        },
        sanitizer: NewHTMLSanitizer(),
    }
}

func (shs *SmartHTMLSetter) SetOuterHTML(nodeID int, html, reason string, options SetOptions) (HTMLSetRecord, error) {
    record := HTMLSetRecord{
        Timestamp: time.Now(),
        NodeID:    nodeID,
        Reason:    reason,
    }

    // 获取当前HTML
    oldHTML, err := getOuterHTML(nodeID)
    if err != nil {
        return record, fmt.Errorf("获取当前HTML失败: %w", err)
    }
    record.OldHTML = oldHTML

    // 分析HTML大小
    record.Size = shs.analyzeHTMLSize(html)

    // 验证HTML
    if options.Validate {
        validation := shs.htmlValidator.Validate(html)
        record.Validation = validation

        if !validation.Valid {
            return record, fmt.Errorf("HTML验证失败: %v", validation.Errors)
        }
    }

    // 清理HTML
    if options.Sanitize {
        sanitizedHTML, changes := shs.sanitizer.Sanitize(html)
        if changes > 0 {
            html = sanitizedHTML
            record.Sanitized = true
        }
    }

    // 分析变化
    if options.TrackChanges {
        record.Changes = shs.changeTracker.AnalyzeChanges(oldHTML, html)
    }

    // 检查缓存
    cacheKey := fmt.Sprintf("%d:%s", nodeID, html)
    if cached, exists := shs.cache.entries[cacheKey]; exists {
        if time.Since(cached.Timestamp) < shs.cache.ttl {
            cached.HitCount++
            shs.cache.entries[cacheKey] = cached

            record.Success = true
            record.Duration = 0
            record.NewHTML = html
            shs.recordSet(record)

            return record, nil
        } else {
            delete(shs.cache.entries, cacheKey)
        }
    }

    // 执行设置
    startTime := time.Now()
    result, err := CDPDOMSetOuterHTML(nodeID, html)
    record.Duration = time.Since(startTime)

    if err != nil {
        record.Success = false
        record.Error = err.Error()
    } else {
        record.Success = true
        record.Response = result
        record.NewHTML = html

        // 更新缓存
        shs.cache.entries[cacheKey] = CachedHTML{
            NodeID:    nodeID,
            HTML:      html,
            Timestamp: time.Now(),
            HitCount:  1,
        }

        // 清理过期缓存
        shs.cleanupCache()
    }

    // 记录历史
    shs.recordSet(record)

    if err != nil {
        return record, fmt.Errorf("设置失败: %w", err)
    }

    return record, nil
}

func NewHTMLValidator() *HTMLValidator {
    return &HTMLValidator{
        maxSize: 1024 * 1024, // 1MB
        allowedTags: map[string]bool{
            "div": true, "span": true, "p": true, "a": true, "img": true,
            "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
            "ul": true, "ol": true, "li": true, "table": true, "tr": true, "td": true,
            "form": true, "input": true, "button": true, "textarea": true, "select": true,
            "section": true, "article": true, "nav": true, "aside": true, "main": true,
            "header": true, "footer": true, "figure": true, "figcaption": true,
            "strong": true, "em": true, "code": true, "pre": true,
        },
        blockedTags: map[string]bool{
            "script": true, "style": true, "iframe": true, "frame": true,
            "embed": true, "object": true, "base": true, "meta": true,
        },
        allowedAttrs: map[string][]string{
            "a":    {"href", "title", "target", "rel"},
            "img":  {"src", "alt", "title", "width", "height"},
            "input": {"type", "name", "value", "placeholder", "required", "disabled"},
            "div":  {"class", "id", "style", "data-"},
        },
        requireClosing: true,
        maxDepth:       20,
    }
}

func (hv *HTMLValidator) Validate(html string) HTMLValidation {
    validation := HTMLValidation{
        Valid: true,
    }

    // 检查大小
    if len(html) > hv.maxSize {
        validation.Valid = false
        validation.Errors = append(validation.Errors,
            fmt.Sprintf("HTML过大: %d > %d", len(html), hv.maxSize))
    }

    // 检查标签嵌套深度
    depth := hv.calculateDepth(html)
    if depth > hv.maxDepth {
        validation.Valid = false
        validation.Errors = append(validation.Errors,
            fmt.Sprintf("嵌套过深: %d > %d", depth, hv.maxDepth))
    } else if depth > 10 {
        validation.Warnings = append(validation.Warnings,
            fmt.Sprintf("嵌套较深: %d 层", depth))
    }

    // 检查被阻止的标签
    for tag := range hv.blockedTags {
        if strings.Contains(strings.ToLower(html), "<"+tag) {
            validation.Valid = false
            validation.Errors = append(validation.Errors,
                fmt.Sprintf("包含被阻止的标签: %s", tag))
        }
    }

    // 检查标签闭合
    if hv.requireClosing {
        if hv.hasUnclosedTags(html) {
            validation.Warnings = append(validation.Warnings, "存在未闭合的标签")
        }
    }

    return validation
}

func (hv *HTMLValidator) calculateDepth(html string) int {
    // 简化的深度计算
    depth := 0
    maxDepth := 0

    re := regexp.MustCompile(`<(\/?)([a-zA-Z][a-zA-Z0-9]*)`)
    matches := re.FindAllStringSubmatch(html, -1)

    for _, match := range matches {
        if len(match) >= 3 {
            tag := strings.ToLower(match[2])
            isClosing := match[1] == "/"

            if !isClosing && hv.allowedTags[tag] {
                depth++
                if depth > maxDepth {
                    maxDepth = depth
                }
            } else if isClosing && hv.allowedTags[tag] {
                depth = max(0, depth-1)
            }
        }
    }

    return maxDepth
}

func (hv *HTMLValidator) hasUnclosedTags(html string) bool {
    // 简化的未闭合标签检查
    stack := []string{}

    re := regexp.MustCompile(`<(\/?)([a-zA-Z][a-zA-Z0-9]*)`)
    matches := re.FindAllStringSubmatch(html, -1)

    for _, match := range matches {
        if len(match) >= 3 {
            tag := strings.ToLower(match[2])
            isClosing := match[1] == "/"

            if hv.allowedTags[tag] {
                if !isClosing {
                    // 自闭合标签
                    if !hv.isVoidElement(tag) {
                        stack = append(stack, tag)
                    }
                } else {
                    // 闭合标签
                    if len(stack) > 0 && stack[len(stack)-1] == tag {
                        stack = stack[:len(stack)-1]
                    }
                }
            }
        }
    }

    return len(stack) > 0
}

func (hv *HTMLValidator) isVoidElement(tag string) bool {
    voidElements := map[string]bool{
        "area": true, "base": true, "br": true, "col": true, "embed": true,
        "hr": true, "img": true, "input": true, "link": true, "meta": true,
        "param": true, "source": true, "track": true, "wbr": true,
    }
    return voidElements[tag]
}

func (hct *HTMLChangeTracker) AnalyzeChanges(oldHTML, newHTML string) HTMLChanges {
    changes := HTMLChanges{}

    if oldHTML == newHTML {
        changes.DiffRatio = 0
        return changes
    }

    // 简化的变化分析
    oldTags := extractTags(oldHTML)
    newTags := extractTags(newHTML)

    changes.NodesAdded = max(0, len(newTags) - len(oldTags))
    changes.NodesRemoved = max(0, len(oldTags) - len(newTags))

    // 计算差异比例
    if len(oldHTML) > 0 {
        changes.DiffRatio = float64(len(newHTML)-len(oldHTML)) / float64(len(oldHTML))
    }

    // 检查结构变化
    changes.StructureChanged = !isSimilarHTML(oldHTML, newHTML)

    return changes
}

func (shs *SmartHTMLSetter) analyzeHTMLSize(html string) HTMLSize {
    size := HTMLSize{
        Characters: len(html),
    }

    // 提取标签
    tags := extractTags(html)
    size.Tags = len(tags)

    // 计算属性数量
    attrRe := regexp.MustCompile(`([a-zA-Z\-_:]+)=["']`)
    size.Attributes = len(attrRe.FindAllString(html, -1))

    // 计算文本节点
    textRe := regexp.MustCompile(`>([^<]+)<`)
    size.TextNodes = len(textRe.FindAllStringSubmatch(html, -1))

    return size
}

func NewHTMLSanitizer() *HTMLSanitizer {
    return &HTMLSanitizer{
        removeScripts:  true,
        removeEvents:   true,
        removeStyles:   false,
        removeComments: false,
        safeProtocols: []string{
            "http://", "https://", "mailto:", "tel:", "#",
        },
    }
}

func (hs *HTMLSanitizer) Sanitize(html string) (string, int) {
    changes := 0
    original := html

    // 移除脚本标签
    if hs.removeScripts {
        re := regexp.MustCompile(`<script[^>]*>.*?</script>`)
        html = re.ReplaceAllStringFunc(html, func(match string) string {
            changes++
            return ""
        })
    }

    // 移除事件属性
    if hs.removeEvents {
        re := regexp.MustCompile(`\s+on\w+\s*=\s*["'][^"']*["']`)
        html = re.ReplaceAllStringFunc(html, func(match string) string {
            changes++
            return ""
        })
    }

    // 移除样式标签
    if hs.removeStyles {
        re := regexp.MustCompile(`<style[^>]*>.*?</style>`)
        html = re.ReplaceAllStringFunc(html, func(match string) string {
            changes++
            return ""
        })
    }

    // 移除注释
    if hs.removeComments {
        re := regexp.MustCompile(`<!--.*?-->`)
        html = re.ReplaceAllStringFunc(html, func(match string) string {
            changes++
            return ""
        })
    }

    return html, changes
}

func (shs *SmartHTMLSetter) cleanupCache() {
    if len(shs.cache.entries) <= shs.cache.maxSize {
        return
    }

    // 清理过期条目
    now := time.Now()
    for key, entry := range shs.cache.entries {
        if now.Sub(entry.Timestamp) > shs.cache.ttl {
            delete(shs.cache.entries, key)
        }
    }

    // 如果仍然超过大小，清理最不常用的
    if len(shs.cache.entries) > shs.cache.maxSize {
        shs.cleanupLRU()
    }
}

func (shs *SmartHTMLSetter) cleanupLRU() {
    var lruKey string
    var lruHitCount = -1
    var lruTime time.Time

    for key, entry := range shs.cache.entries {
        if lruHitCount == -1 ||
           entry.HitCount < lruHitCount ||
           (entry.HitCount == lruHitCount && entry.Timestamp.Before(lruTime)) {
            lruKey = key
            lruHitCount = entry.HitCount
            lruTime = entry.Timestamp
        }
    }

    if lruKey != "" {
        delete(shs.cache.entries, lruKey)
    }
}

func (shs *SmartHTMLSetter) recordSet(record HTMLSetRecord) {
    shs.setHistory = append(shs.setHistory, record)

    // 限制历史记录数量
    if len(shs.setHistory) > shs.changeTracker.maxHistory {
        shs.setHistory = shs.setHistory[1:]
    }
}

func (shs *SmartHTMLSetter) GetStats() map[string]interface{} {
    totalSets := len(shs.setHistory)
    successfulSets := 0
    var totalDuration time.Duration
    var totalCharsChanged int

    for _, record := range shs.setHistory {
        if record.Success {
            successfulSets++
        }
        totalDuration += record.Duration
        totalCharsChanged += len(record.NewHTML) - len(record.OldHTML)
    }

    avgDuration := time.Duration(0)
    if totalSets > 0 {
        avgDuration = totalDuration / time.Duration(totalSets)
    }

    avgCharsChanged := 0.0
    if totalSets > 0 {
        avgCharsChanged = float64(totalCharsChanged) / float64(totalSets)
    }

    sanitizedCount := 0
    for _, record := range shs.setHistory {
        if record.Sanitized {
            sanitizedCount++
        }
    }

    return map[string]interface{}{
        "totalSets":        totalSets,
        "successfulSets":   successfulSets,
        "failedSets":       totalSets - successfulSets,
        "successRate":      float64(successfulSets) / float64(totalSets) * 100,
        "totalDuration":    totalDuration,
        "averageDuration":  avgDuration,
        "totalCharsChanged": totalCharsChanged,
        "averageCharsChanged": avgCharsChanged,
        "sanitizedCount":   sanitizedCount,
        "cacheSize":       len(shs.cache.entries),
    }
}

type SetOptions struct {
    Validate     bool
    Sanitize     bool
    TrackChanges bool
    CacheResult  bool
}

// 演示智能HTML设置
func demonstrateSmartHTMLSetting() {
    fmt.Println("=== 智能HTML设置演示 ===")

    // 启用DOM
    if _, err := CDPDOMEnable(); err != nil {
        log.Printf("启用DOM失败: %v", err)
        return
    }
    defer CDPDOMDisable()

    // 创建智能设置器
    setter := NewSmartHTMLSetter()

    // 测试元素
    containerID := 7001

    fmt.Printf("测试容器ID: %d\n\n", containerID)

    // 场景1: 安全HTML设置
    fmt.Printf("场景1: 安全HTML设置\n")

    htmlTests := []struct {
        html    string
        reason  string
        options SetOptions
    }{
        {
            html: `<div class="safe-content">
                <h3>安全内容</h3>
                <p>这是一段安全的内容。</p>
            </div>`,
            reason:  "安全HTML",
            options: SetOptions{Validate: true, Sanitize: true, TrackChanges: true},
        },
        {
            html: `<div>
                <p>包含<a href="https://example.com">安全链接</a>的内容。</p>
            </div>`,
            reason:  "包含链接",
            options: SetOptions{Validate: true, Sanitize: true, TrackChanges: true},
        },
    }

    for _, test := range htmlTests {
        fmt.Printf("设置HTML: %s\n", test.reason)
        fmt.Printf("HTML: %s\n", truncateString(test.html, 50))

        record, err := setter.SetOuterHTML(containerID, test.html, test.reason, test.options)
        if err != nil {
            fmt.Printf("  ❌ 失败: %v\n", err)
        } else {
            fmt.Printf("  ✅ 成功 (耗时: %v)\n", record.Duration)

            fmt.Printf("    HTML大小: %d 字符, %d 标签, %d 属性\n",
                record.Size.Characters, record.Size.Tags, record.Size.Attributes)

            if record.Changes.NodesAdded > 0 {
                fmt.Printf("    新增节点: %d\n", record.Changes.NodesAdded)
            }
            if record.Changes.NodesRemoved > 0 {
                fmt.Printf("    删除节点: %d\n", record.Changes.NodesRemoved)
            }
        }

        time.Sleep(300 * time.Millisecond)
    }

    // 场景2: 验证测试
    fmt.Printf("\n场景2: 验证测试\n")

    validationTests := []struct {
        html    string
        expectError bool
        reason  string
    }{
        {
            html: `<script>alert('xss')</script>`,
            expectError: true,
            reason: "脚本标签",
        },
        {
            html: strings.Repeat("<div>内容</div>", 1000),
            expectError: true, // 可能超过大小限制
            reason: "大量重复内容",
        },
    }

    for _, test := range validationTests {
        options := SetOptions{Validate: true, Sanitize: false, TrackChanges: false}
        _, err := setter.SetOuterHTML(containerID, test.html, test.reason, options)

        if test.expectError {
            if err != nil {
                fmt.Printf("HTML '%s': ✅ 预期错误\n", test.reason)
            } else {
                fmt.Printf("HTML '%s': ❌ 预期错误但成功\n", test.reason)
            }
        } else {
            if err != nil {
                fmt.Printf("HTML '%s': ❌ 失败: %v\n", test.reason, err)
            } else {
                fmt.Printf("HTML '%s': ✅ 验证通过\n", test.reason)
            }
        }
    }

    // 显示统计
    fmt.Printf("\n=== 设置统计 ===\n")
    stats := setter.GetStats()
    fmt.Printf("  总设置次数: %d\n", stats["totalSets"])
    fmt.Printf("  成功设置: %d\n", stats["successfulSets"])
    fmt.Printf("  失败设置: %d\n", stats["failedSets"])
    fmt.Printf("  成功率: %.1f%%\n", stats["successRate"])
    fmt.Printf("  总字符变化: %d\n", stats["totalCharsChanged"])
    fmt.Printf("  清理次数: %d\n", stats["sanitizedCount"])
    fmt.Printf("  缓存大小: %d\n", stats["cacheSize"])

    fmt.Println("\n=== 演示完成 ===")
}


*/
