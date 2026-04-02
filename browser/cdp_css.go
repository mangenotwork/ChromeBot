package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
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
func CDPCSSAddRule(styleSheetId string, ruleText string, location ...SourceRange) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}
	if ruleText == "" {
		return "", fmt.Errorf("规则文本不能为空")
	}

	// 转义特殊字符
	escapedRuleText := strings.ReplaceAll(ruleText, `"`, `\"`)
	escapedRuleText = strings.ReplaceAll(escapedRuleText, "\n", "\\n")

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	var message string
	if len(location) > 0 {
		// 包含位置参数
		loc := location[0]
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "CSS.addRule",
			"params": {
				"styleSheetId": "%s",
				"ruleText": "%s",
				"location": {
					"startLine": %d,
					"startColumn": %d,
					"endLine": %d,
					"endColumn": %d
				}
			}
		}`, reqID, styleSheetId, escapedRuleText,
			loc.StartLine, loc.StartColumn, loc.EndLine, loc.EndColumn)
	} else {
		// 不包含位置参数
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "CSS.addRule",
			"params": {
				"styleSheetId": "%s",
				"ruleText": "%s"
			}
		}`, reqID, styleSheetId, escapedRuleText)
	}

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

/*

示例

// 示例1: 向样式表添加基本规则
func exampleAddBasicRule() {
	styleSheetId := "style-sheet-1"
	ruleText := ".my-class { color: red; font-size: 16px; }"

	log.Printf("向样式表 %s 添加规则: %s", styleSheetId, ruleText)

	response, err := CDPCSSAddRule(styleSheetId, ruleText)
	if err != nil {
		log.Printf("添加规则失败: %v", err)
		return
	}

	log.Printf("添加成功: %s", response)

	// 解析返回的规则信息
	var data struct {
		Result struct {
			Rule struct {
				SelectorText string `json:"selectorText"`
				Style        struct {
					CSSProperties []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"cssProperties"`
				} `json:"style"`
			} `json:"rule"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err == nil {
		log.Printf("新规则选择器: %s", data.Result.Rule.SelectorText)
		log.Printf("样式属性数: %d", len(data.Result.Rule.Style.CSSProperties))
	}
}

// 示例2: 在指定位置添加规则
func exampleAddRuleWithLocation() {
	styleSheetId := "style-sheet-1"
	ruleText := ".highlight { background-color: yellow; }"

	// 在样式表的第5行第1列插入规则
	location := NewSourceRange(5, 1, 5, 1)

	log.Printf("在位置 %d:%d 添加规则", location.StartLine, location.StartColumn)

	response, err := CDPCSSAddRule(styleSheetId, ruleText, location)
	if err != nil {
		log.Printf("添加规则失败: %v", err)
		return
	}

	log.Printf("添加成功: %s", response)
}



// 完整的添加规则流程
func completeAddRuleWorkflow() {
	// 1. 启用CSS域
	_, err := CDPCSSEnable()
	if err != nil {
		log.Printf("启用CSS域失败: %v", err)
		return
	}
	defer CDPCSSDisable()

	// 2. 获取样式表列表
	styleSheets, err := getAllStyleSheets()
	if err != nil {
		log.Printf("获取样式表失败: %v", err)
		return
	}

	if len(styleSheets) == 0 {
		log.Printf("未找到样式表")
		return
	}

	// 3. 选择第一个样式表
	styleSheet := styleSheets[0]
	log.Printf("使用样式表: %s (%s)", styleSheet.ID, styleSheet.SourceURL)

	// 4. 添加规则
	ruleText := `
	.my-dynamic-class {
		color: blue;
		font-weight: bold;
		padding: 10px;
	}`

		response, err := CDPCSSAddRule(styleSheet.ID, ruleText)
		if err != nil {
			log.Printf("添加规则失败: %v", err)
			return
		}

		log.Printf("规则添加成功: %s", response)
	}

*/

// -----------------------------------------------  CSS.collectClassNames  -----------------------------------------------

// CDPCSSColectClassNames 从指定样式表中收集所有类名
// 这是 CDP 协议的原生方法
// 参数说明:
//   - styleSheetId: 样式表ID
func CDPCSSColectClassNames(styleSheetId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

/*

示例

// 示例1: 收集样式表中的类名
func exampleCollectClassNames() {
	// 先获取样式表ID
	styleSheets, err := getAllStyleSheets()
	if err != nil {
		log.Printf("获取样式表失败: %v", err)
		return
	}

	if len(styleSheets) == 0 {
		log.Printf("未找到样式表")
		return
	}

	// 使用第一个样式表
	styleSheetId := styleSheets[0].ID

	log.Printf("收集样式表 %s 的类名...", styleSheetId)

	response, err := CDPCSSColectClassNames(styleSheetId)
	if err != nil {
		log.Printf("收集类名失败: %v", err)
		return
	}

	// 解析响应
	var data struct {
		Result struct {
			ClassNames []string `json:"classNames"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		log.Printf("解析响应失败: %v", err)
		return
	}

	log.Printf("收集到 %d 个类名:", len(data.Result.ClassNames))

	for i, className := range data.Result.ClassNames {
		if i < 20 { // 只显示前20个
			log.Printf("  [%d] %s", i+1, className)
		} else {
			log.Printf("  ... 还有 %d 个类名", len(data.Result.ClassNames)-20)
			break
		}
	}
}

// 示例2: 分析多个样式表的类名
func exampleAnalyzeMultipleStyleSheets() {
	styleSheets, err := getAllStyleSheets()
	if err != nil {
		log.Printf("获取样式表失败: %v", err)
		return
	}

	log.Printf("分析 %d 个样式表的类名...", len(styleSheets))

	totalClassNames := 0
	allClassNames := make(map[string]int) // 类名 -> 出现次数

	for i, sheet := range styleSheets {
		log.Printf("样式表 %d/%d: %s", i+1, len(styleSheets), sheet.SourceURL)

		response, err := CDPCSSColectClassNames(sheet.ID)
		if err != nil {
			log.Printf("  收集失败: %v", err)
			continue
		}

		var data struct {
			Result struct {
				ClassNames []string `json:"classNames"`
			} `json:"result"`
		}

		if err := json.Unmarshal([]byte(response), &data); err != nil {
			log.Printf("  解析失败: %v", err)
			continue
		}

		classCount := len(data.Result.ClassNames)
		totalClassNames += classCount

		for _, className := range data.Result.ClassNames {
			allClassNames[className]++
		}

		log.Printf("  类名数: %d", classCount)

		time.Sleep(50 * time.Millisecond) // 避免请求过快
	}

	log.Printf("\n分析结果:")
	log.Printf("  样式表总数: %d", len(styleSheets))
	log.Printf("  总类名数: %d", totalClassNames)
	log.Printf("  唯一类名数: %d", len(allClassNames))

	// 查找重复的类名
	var duplicateClasses []string
	for className, count := range allClassNames {
		if count > 1 {
			duplicateClasses = append(duplicateClasses, fmt.Sprintf("%s (%d次)", className, count))
		}
	}

	if len(duplicateClasses) > 0 {
		log.Printf("  重复的类名: %d 个", len(duplicateClasses))
		for i, dup := range duplicateClasses {
			if i < 10 { // 只显示前10个
				log.Printf("    - %s", dup)
			} else {
				log.Printf("    ... 还有 %d 个", len(duplicateClasses)-10)
				break
			}
		}
	}
}

*/

// 获取所有样式表
func getAllStyleSheets() ([]StyleSheetInfo, error) {
	// 先启用CSS域
	_, err := CDPCSSEnable()
	if err != nil {
		return nil, fmt.Errorf("启用CSS域失败: %w", err)
	}

	// 获取样式表
	response, err := CDPCSSGetStyleSheets()
	if err != nil {
		return nil, fmt.Errorf("获取样式表失败: %w", err)
	}

	// 解析响应
	var data struct {
		Result struct {
			Headers []struct {
				StyleSheetID string `json:"styleSheetId"`
				SourceURL    string `json:"sourceURL"`
				Title        string `json:"title"`
				Disabled     bool   `json:"disabled"`
			} `json:"headers"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	var styleSheets []StyleSheetInfo
	for _, header := range data.Result.Headers {
		styleSheets = append(styleSheets, StyleSheetInfo{
			ID:        header.StyleSheetID,
			SourceURL: header.SourceURL,
			Title:     header.Title,
			Disabled:  header.Disabled,
		})
	}

	return styleSheets, nil
}

// StyleSheetInfo 样式表信息
type StyleSheetInfo struct {
	ID        string
	SourceURL string
	Title     string
	Disabled  bool
}

// CDPCSSGetStyleSheets 获取页面所有样式表信息
// 这是 CDP 协议的方法
// 返回值:
//   - 包含样式表信息的JSON字符串
//   - error: 获取过程中发生的错误
func CDPCSSGetStyleSheets() (string, error) {
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
		"method": "CSS.getStyleSheets"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getStyleSheets 请求失败: %w", err)
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
						return content, fmt.Errorf("CSS.getStyleSheets 方法在当前Chrome版本中可能不存在")
					}
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("getStyleSheets 请求超时")
		}
	}
}

/*

示例

// 示例1: 获取并显示所有样式表
func exampleGetAllStyleSheets() {
	// 先启用CSS域
	_, err := CDPCSSEnable()
	if err != nil {
		log.Printf("启用CSS域失败: %v", err)
		return
	}
	defer CDPCSSDisable()

	// 获取样式表
	response, err := CDPCSSGetStyleSheets()
	if err != nil {
		log.Printf("获取样式表失败: %v", err)
		return
	}

	// 解析响应
	var data struct {
		Result struct {
			Headers []struct {
				StyleSheetID string `json:"styleSheetId"`
				FrameID      string `json:"frameId"`
				SourceURL    string `json:"sourceURL"`
				Title        string `json:"title"`
				Disabled     bool   `json:"disabled"`
				IsInline     bool   `json:"isInline"`
				StartLine    int    `json:"startLine"`
				StartColumn  int    `json:"startColumn"`
				Length       int    `json:"length"`
				Origin       string `json:"origin"`
			} `json:"headers"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		log.Printf("解析响应失败: %v", err)
		return
	}

	log.Printf("找到 %d 个样式表:", len(data.Result.Headers))

	for i, sheet := range data.Result.Headers {
		log.Printf("样式表 %d:", i+1)
		log.Printf("  ID: %s", sheet.StyleSheetID)
		log.Printf("  来源: %s", sheet.SourceURL)
		log.Printf("  标题: %s", sheet.Title)
		log.Printf("  是否内联: %v", sheet.IsInline)
		log.Printf("  是否禁用: %v", sheet.Disabled)
		log.Printf("  起始位置: 行%d, 列%d", sheet.StartLine, sheet.StartColumn)
		log.Printf("  长度: %d 字符", sheet.Length)
		log.Printf("  来源类型: %s", sheet.Origin)
		fmt.Println()
	}
}

// 获取所有样式表的辅助函数
func getAllStyleSheets() ([]StyleSheetInfo, error) {
	// 先启用CSS域
	_, err := CDPCSSEnable()
	if err != nil {
		return nil, fmt.Errorf("启用CSS域失败: %w", err)
	}

	// 获取样式表
	response, err := CDPCSSGetStyleSheets()
	if err != nil {
		return nil, fmt.Errorf("获取样式表失败: %w", err)
	}

	// 解析响应
	var data struct {
		Result struct {
			Headers []struct {
				StyleSheetID string `json:"styleSheetId"`
				SourceURL    string `json:"sourceURL"`
				Title        string `json:"title"`
				Disabled     bool   `json:"disabled"`
				IsInline     bool   `json:"isInline"`
				Origin       string `json:"origin"`
			} `json:"headers"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	var styleSheets []StyleSheetInfo
	for _, header := range data.Result.Headers {
		styleSheets = append(styleSheets, StyleSheetInfo{
			ID:        header.StyleSheetID,
			SourceURL: header.SourceURL,
			Title:     header.Title,
			Disabled:  header.Disabled,
			IsInline:  header.IsInline,
			Origin:    header.Origin,
		})
	}

	return styleSheets, nil
}

*/

// CDPCSSEnable 启用CSS域
func CDPCSSEnable() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.enable"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.disable"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// CDPCSSCreateStyleSheet 创建一个新的"via-inspector"样式表
// 这是 CDP 协议的原生方法
// 参数说明:
//   - frameId: 页面框架ID
//   - force: 是否强制创建新样式表
func CDPCSSCreateStyleSheet(frameId string, force bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if frameId == "" {
		return "", fmt.Errorf("框架ID不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.createStyleSheet",
		"params": {
			"frameId": "%s",
			"force": %t
		}
	}`, reqID, frameId, force)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// 获取当前页面框架ID
func getCurrentFrameId() (string, error) {
	// 通过Page域获取框架信息
	response, err := CDPPageGetFrameTree()
	if err != nil {
		return "", fmt.Errorf("获取框架树失败: %w", err)
	}

	// 解析响应
	var data struct {
		Result struct {
			FrameTree struct {
				Frame struct {
					ID string `json:"id"`
				} `json:"frame"`
			} `json:"frameTree"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return "", fmt.Errorf("解析框架树失败: %w", err)
	}

	if data.Result.FrameTree.Frame.ID == "" {
		return "", fmt.Errorf("未找到框架ID")
	}

	return data.Result.FrameTree.Frame.ID, nil
}

// CDPPageGetFrameTree 获取页面框架树
func CDPPageGetFrameTree() (string, error) {
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
		"method": "Page.getFrameTree"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getFrameTree 请求失败: %w", err)
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
			return "", fmt.Errorf("getFrameTree 请求超时")
		}
	}
}

/*

// 示例1: 创建新的"via-inspector"样式表
func exampleCreateStyleSheet() {
	// 先获取当前页面的框架ID
	// 假设我们已经有了frameId
	frameId := "frame-id-123"

	log.Printf("在框架 %s 中创建样式表...", frameId)

	// 强制创建新样式表
	response, err := CDPCSSCreateStyleSheet(frameId, true)
	if err != nil {
		log.Printf("创建样式表失败: %v", err)
		return
	}

	// 解析响应获取样式表ID
	var data struct {
		Result struct {
			StyleSheetID string `json:"styleSheetId"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		log.Printf("解析响应失败: %v", err)
		return
	}

	log.Printf("样式表创建成功，ID: %s", data.Result.StyleSheetID)

	// 可以使用这个样式表ID进行其他操作
	// 例如：向样式表添加规则
	ruleText := `
.via-inspector-element {
	color: red;
	border: 2px dashed #ccc;
	background-color: #f0f0f0;
	padding: 10px;
	margin: 5px;
}
	`

	// 向新创建的样式表添加规则
	addRuleResponse, err := CDPCSSAddRule(data.Result.StyleSheetID, ruleText)
	if err != nil {
		log.Printf("添加规则失败: %v", err)
	} else {
		log.Printf("规则添加成功: %s", addRuleResponse)
	}
}

// 示例2: 获取或创建样式表
func exampleGetOrCreateStyleSheet(frameId string) (string, error) {
	// 先尝试获取已存在的样式表
	styleSheets, err := getAllStyleSheets()
	if err != nil {
		return "", fmt.Errorf("获取样式表失败: %w", err)
	}

	// 查找已经存在的"via-inspector"样式表
	for _, sheet := range styleSheets {
		// 检查是否是"via-inspector"样式表
		// 可以通过检查来源或其他特征
		if sheet.Origin == "inspector" || strings.Contains(sheet.Title, "via-inspector") {
			log.Printf("找到已存在的via-inspector样式表: %s", sheet.ID)
			return sheet.ID, nil
		}
	}

	// 如果没有找到，创建新的样式表
	log.Printf("未找到via-inspector样式表，创建新的...")
	response, err := CDPCSSCreateStyleSheet(frameId, false)
	if err != nil {
		return "", fmt.Errorf("创建样式表失败: %w", err)
	}

	var data struct {
		Result struct {
			StyleSheetID string `json:"styleSheetId"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	log.Printf("新样式表创建成功，ID: %s", data.Result.StyleSheetID)
	return data.Result.StyleSheetID, nil
}


// 完整的创建和使用样式表流程
func completeCreateStyleSheetWorkflow() {
	// 1. 启用CSS域
	_, err := CDPCSSEnable()
	if err != nil {
		log.Printf("启用CSS域失败: %v", err)
		return
	}
	defer CDPCSSDisable()

	// 2. 获取当前框架ID
	frameId, err := getCurrentFrameId()
	if err != nil {
		log.Printf("获取框架ID失败: %v", err)
		// 如果无法获取框架ID，可以尝试使用默认值
		frameId = "main" // 默认值
		log.Printf("使用默认框架ID: %s", frameId)
	}

	// 3. 创建或获取样式表
	styleSheetId, err := exampleGetOrCreateStyleSheet(frameId)
	if err != nil {
		log.Printf("获取样式表失败: %v", err)
		return
	}

	log.Printf("样式表ID: %s", styleSheetId)

	// 4. 向样式表添加多个规则
	rules := []struct {
		name     string
		ruleText string
	}{
		{
			name:     "调试高亮",
			ruleText: ".debug-highlight { background-color: yellow !important; border: 3px solid red !important; }",
		},
		{
			name:     "隐藏元素",
			ruleText: ".hide-debug { display: none !important; }",
		},
		{
			name:     "显示网格",
			ruleText: ".show-grid { outline: 1px dashed #999; }",
		},
		{
			name:     "测量尺寸",
			ruleText: ".measure-size { position: relative; } .measure-size::after { content: attr(data-size); position: absolute; top: 0; right: 0; background: #333; color: white; padding: 2px 4px; font-size: 12px; }",
		},
	}

	for i, rule := range rules {
		log.Printf("添加规则 %d/%d: %s", i+1, len(rules), rule.name)

		response, err := CDPCSSAddRule(styleSheetId, rule.ruleText)
		if err != nil {
			log.Printf("  添加失败: %v", err)
		} else {
			log.Printf("  添加成功")
		}

		time.Sleep(50 * time.Millisecond)
	}

	log.Printf("样式表创建和规则添加完成")

	// 5. 可选：验证样式表内容
	log.Printf("验证样式表内容...")
	verifyStyleSheet(styleSheetId)
}

// 验证样式表内容
func verifyStyleSheet(styleSheetId string) {
	// 获取样式表文本
	response, err := CDPCSSGetStyleSheetText(styleSheetId)
	if err != nil {
		log.Printf("获取样式表文本失败: %v", err)
		return
	}

	var data struct {
		Result struct {
			Text string `json:"text"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		log.Printf("解析响应失败: %v", err)
		return
	}

	log.Printf("样式表内容:")
	log.Printf("长度: %d 字符", len(data.Result.Text))
	if len(data.Result.Text) > 0 {
		log.Printf("预览: %s", data.Result.Text[:min(100, len(data.Result.Text))])
		if len(data.Result.Text) > 100 {
			log.Printf("... (内容截断)")
		}
	}
}


*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}

	if len(forcedPseudoClasses) == 0 {
		return "", fmt.Errorf("必须指定至少一个伪类")
	}

	// 验证伪类格式
	validPseudoClasses := map[string]bool{
		"hover":             true,
		"focus":             true,
		"focus-visible":     true,
		"focus-within":      true,
		"active":            true,
		"visited":           true,
		"link":              true,
		"target":            true,
		"enabled":           true,
		"disabled":          true,
		"checked":           true,
		"indeterminate":     true,
		"default":           true,
		"in-range":          true,
		"out-of-range":      true,
		"placeholder-shown": true,
		"valid":             true,
		"invalid":           true,
		"required":          true,
		"optional":          true,
		"read-only":         true,
		"read-write":        true,
		"first-child":       true,
		"last-child":        true,
		"nth-child":         true,
		"first-of-type":     true,
		"last-of-type":      true,
		"nth-of-type":       true,
		"only-child":        true,
		"only-of-type":      true,
		"empty":             true,
		"root":              true,
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
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

/*

// 示例2: 自动化测试 - 测试焦点状态
func exampleAutomationFocusTest() {
	// === 应用场景描述 ===
	// 场景: 表单元素的焦点状态测试
	// 用途: 测试输入框在不同焦点状态下的样式
	// 优势: 自动化验证焦点样式，无需人工操作
	// 典型工作流: 选择输入框 -> 启用焦点状态 -> 截图验证 -> 清理状态

	inputNodeId := 37

	// 测试普通焦点状态
	log.Println("测试输入框焦点状态...")
	focusResponse, err := CDPCSSForcePseudoState(inputNodeId, []string{"focus"})
	if err != nil {
		log.Printf("启用焦点状态失败: %v", err)
	} else {
		log.Printf("焦点状态已启用: %s", focusResponse)

		// 在这里可以进行截图或样式验证
		// takeScreenshot("input-focus-state.png")

		// 验证焦点状态下的样式
		// verifyFocusStyles(inputNodeId)
	}

	// 测试焦点可见状态
	log.Println("测试焦点可见状态...")
	focusVisibleResponse, err := CDPCSSForcePseudoState(inputNodeId, []string{"focus-visible"})
	if err != nil {
		log.Printf("启用焦点可见状态失败: %v", err)
	} else {
		log.Printf("焦点可见状态已启用: %s", focusVisibleResponse)
	}

	// 组合多个伪类
	log.Println("测试组合伪类状态...")
	combinedResponse, err := CDPCSSForcePseudoState(inputNodeId, []string{"focus", "valid"})
	if err != nil {
		log.Printf("启用组合状态失败: %v", err)
	} else {
		log.Printf("组合状态已启用: %s", combinedResponse)
	}

	// 清理所有状态
	log.Println("清理所有伪类状态...")
	_, err = CDPCSSForcePseudoState(inputNodeId, []string{})
	if err != nil {
		log.Printf("清理状态失败: %v", err)
	}
}

// 示例3: 组件库测试 - 测试按钮的各种状态
func exampleComponentLibraryTest() {
	// === 应用场景描述 ===
	// 场景: UI组件库的交互状态测试
	// 用途: 自动化测试按钮在各种交互状态下的样式
	// 优势: 确保组件在所有状态下的视觉一致性
	// 典型工作流: 遍历所有状态 -> 应用伪类 -> 验证样式 -> 生成报告

	buttonNodeId := 42
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

		response, err := CDPCSSForcePseudoState(buttonNodeId, test.pseudoClasses)
		if err != nil {
			log.Printf("  失败: %v", err)
			testResults[test.name] = false
		} else {
			log.Printf("  成功: 状态已应用")
			testResults[test.name] = true

			// 在这里可以进行样式验证或截图
			// verifyButtonStyles(buttonNodeId, test.name)

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

*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

/*

示例

// 示例2: 组件测试 - 测试Modal弹窗的起始动画
func exampleModalComponentTest() {
	// === 应用场景描述 ===
	// 场景: Modal弹窗组件的起始动画测试
	// 用途: 测试Modal在显示时的起始动画样式
	// 优势: 自动化验证Modal的显示动画起始状态
	// 典型工作流: 触发Modal显示 -> 启用起始样式 -> 截图验证 -> 测试动画

	modalElementId := 38

	log.Println("测试Modal弹窗的起始动画样式...")

	// 定义测试用例
	testCases := []struct {
		name  string
		forced bool
		desc  string
	}{
		{"启用起始样式", true, "测试Modal在起始样式状态下的表现"},
		{"禁用起始样式", false, "测试Modal在正常状态下的表现"},
		{"重新启用起始样式", true, "再次测试起始样式状态"},
	}

	for i, testCase := range testCases {
		log.Printf("测试用例 %d/%d: %s", i+1, len(testCases), testCase.name)
		log.Printf("  描述: %s", testCase.desc)

		response, err := CDPCSSForceStartingStyle(modalElementId, testCase.forced)
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


*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// BackgroundColorsResult 背景颜色结果
type BackgroundColorsResult struct {
	BackgroundColors   []string `json:"backgroundColors"`   // 背景颜色数组
	ComputedFontSize   string   `json:"computedFontSize"`   // 计算后的字体大小
	ComputedFontWeight string   `json:"computedFontWeight"` // 计算后的字体粗细
}

// ParseBackgroundColors 解析背景颜色响应
func ParseBackgroundColors(response string) (*BackgroundColorsResult, error) {
	var data struct {
		Result *BackgroundColorsResult `json:"result"`
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

// 示例2: 设计系统验证 - 验证按钮背景颜色
func exampleDesignSystemButtonTest() {
	// === 应用场景描述 ===
	// 场景: 设计系统中按钮组件的颜色验证
	// 用途: 验证按钮在不同状态下的背景颜色是否符合设计规范
	// 优势: 自动化验证设计一致性，防止视觉回归
	// 典型工作流: 选择按钮 -> 获取背景颜色 -> 验证设计token -> 生成报告

	buttonElementId := 45

	log.Println("验证按钮背景颜色设计...")

	// 定义设计规范
	designSpecs := map[string]string{
		"primary-button":   "#007bff",
		"secondary-button": "#6c757d",
		"success-button":   "#28a745",
		"danger-button":    "#dc3545",
	}

	// 获取按钮背景颜色
	response, err := CDPCSSGetBackgroundColors(buttonElementId)
	if err != nil {
		log.Printf("获取背景颜色失败: %v", err)
		return
	}

	result, err := ParseBackgroundColors(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	// 分析背景颜色
	if len(result.BackgroundColors) == 0 {
		log.Println("警告: 按钮没有检测到背景颜色")
		return
	}

	actualColor := result.BackgroundColors[0]
	log.Printf("按钮实际背景颜色: %s", actualColor)

	// 查找最接近的设计规范颜色
	var closestSpec string
	var closestDiff float64 = math.MaxFloat64

	for specName, specColor := range designSpecs {
		diff := colorDifference(actualColor, specColor)
		if diff < closestDiff {
			closestDiff = diff
			closestSpec = specName
		}
		log.Printf("  与%s对比差异: %.2f", specName, diff)
	}

	if closestDiff < 0.1 { // 阈值
		log.Printf("✓ 符合设计规范: %s (差异: %.2f)", closestSpec, closestDiff)
	} else {
		log.Printf("⚠ 与设计规范差异较大，最接近: %s (差异: %.2f)", closestSpec, closestDiff)
	}
}


*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

/*

// 示例1: 开发者工具 - 样式检查器
func exampleDevToolsStyleInspector() {
	// === 应用场景描述 ===
	// 场景: 开发者工具中的样式检查器
	// 用途: 获取元素的所有计算样式，用于样式调试和分析
	// 优势: 查看CSS属性的最终计算值，包括继承、层叠和计算后的单位
	// 典型工作流: 选择元素 -> 获取计算样式 -> 分析样式来源 -> 调试问题

	elementId := 25

	log.Println("获取元素计算样式...")

	// 获取计算样式
	response, err := CDPCSSGetComputedStyleForNode(elementId)
	if err != nil {
		log.Printf("获取计算样式失败: %v", err)
		return
	}

	// 解析结果
	result, err := ParseComputedStyle(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	log.Printf("计算样式分析:")
	log.Printf("  样式属性总数: %d", len(result.ComputedStyle))

	// 分类显示常见样式属性
	commonProperties := map[string]string{
		"display":           "显示",
		"position":          "定位",
		"width":             "宽度",
		"height":            "高度",
		"margin":            "外边距",
		"padding":           "内边距",
		"border":            "边框",
		"color":             "颜色",
		"background-color":  "背景色",
		"font-size":         "字体大小",
		"font-weight":       "字体粗细",
		"line-height":       "行高",
		"text-align":        "文本对齐",
		"visibility":        "可见性",
		"opacity":           "透明度",
		"z-index":           "堆叠顺序",
		"box-sizing":        "盒模型",
		"flex":              "弹性盒子",
		"grid":              "网格",
	}

	log.Println("  关键样式属性:")
	for _, prop := range result.ComputedStyle {
		if desc, exists := commonProperties[prop.Name]; exists {
			log.Printf("    %s (%s): %s", prop.Name, desc, prop.Value)
		}
	}

	// 查找特定样式
	findAndLogProperty(result, "margin")
	findAndLogProperty(result, "padding")
	findAndLogProperty(result, "font-size")
}

// 查找并记录属性
func findAndLogProperty(result *ComputedStyleResult, propertyName string) {
	for _, prop := range result.ComputedStyle {
		if prop.Name == propertyName {
			log.Printf("  %s: %s", propertyName, prop.Value)
			return
		}
	}
	log.Printf("  %s: 未设置", propertyName)
}

// 示例2: 自动化测试 - 样式验证
func exampleAutomationStyleVerification() {
	// === 应用场景描述 ===
	// 场景: 自动化测试中的样式验证
	// 用途: 验证元素在实际渲染后的样式值是否符合预期
	// 优势: 自动化检查视觉一致性，防止样式回归
	// 典型工作流: 定义预期样式 -> 获取计算样式 -> 比较差异 -> 生成报告

	buttonId := 42

	// 预期的样式规范
	expectedStyles := map[string]string{
		"padding":         "10px 20px",
		"border-radius":   "4px",
		"font-weight":     "500",
		"cursor":          "pointer",
		"user-select":     "none",
		"background-color": "rgb(0, 123, 255)",
		"color":           "rgb(255, 255, 255)",
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

// 示例3: 响应式设计测试
func exampleResponsiveDesignTest() {
	// === 应用场景描述 ===
	// 场景: 响应式设计的断点测试
	// 用途: 在不同视口大小下验证元素的计算样式变化
	// 优势: 自动化测试响应式布局的断点切换
	// 典型工作流: 设置视口大小 -> 获取计算样式 -> 验证响应式规则 -> 测试断点

	containerId := 33

	// 定义测试的视口大小
	viewportTests := []struct {
		name   string
		width  int
		height int
		expect struct {
			display     string
			flexDirection string
			maxWidth    string
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

*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// ParseInlineStyles 解析内联样式响应
func ParseInlineStyles(response string) (*InlineStyleResult, error) {
	var data struct {
		Result *InlineStyleResult `json:"result"`
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

示例

// 示例1: 开发者工具 - 内联样式分析器
func exampleDevToolsInlineStyleAnalyzer() {
	// === 应用场景描述 ===
	// 场景: 开发者工具中的内联样式分析
	// 用途: 分析元素的内联样式，区分显式style属性和DOM属性样式
	// 优势: 清晰展示样式来源，帮助理解样式优先级和覆盖关系
	// 典型工作流: 选择元素 -> 获取内联样式 -> 分析来源 -> 调试优先级

	elementId := 25

	log.Println("分析元素内联样式...")

	// 获取内联样式
	response, err := CDPCSSGetInlineStylesForNode(elementId)
	if err != nil {
		log.Printf("获取内联样式失败: %v", err)
		return
	}

	// 解析结果
	result, err := ParseInlineStyles(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	log.Printf("内联样式分析:")

	// 分析显式内联样式（style属性）
	log.Println("1. 显式内联样式 (style属性):")
	if result.InlineStyle != nil && len(result.InlineStyle.CSSProperties) > 0 {
		log.Printf("   属性数量: %d", len(result.InlineStyle.CSSProperties))
		for i, prop := range result.InlineStyle.CSSProperties {
			if i < 10 { // 只显示前10个
				important := ""
				if prop.Important {
					important = " !important"
				}
				log.Printf("   - %s: %s%s", prop.Name, prop.Value, important)
			}
		}
		if len(result.InlineStyle.CSSProperties) > 10 {
			log.Printf("   ... 还有 %d 个属性", len(result.InlineStyle.CSSProperties)-10)
		}
	} else {
		log.Println("   无显式内联样式")
	}

	// 分析属性样式（DOM属性）
	log.Println("2. 属性样式 (DOM属性):")
	if result.AttributesStyle != nil && len(result.AttributesStyle.CSSProperties) > 0 {
		log.Printf("   属性数量: %d", len(result.AttributesStyle.CSSProperties))
		for i, prop := range result.AttributesStyle.CSSProperties {
			if i < 10 { // 只显示前10个
				implicit := ""
				if prop.Implicit {
					implicit = " (隐式)"
				}
				log.Printf("   - %s: %s%s", prop.Name, prop.Value, implicit)
			}
		}
		if len(result.AttributesStyle.CSSProperties) > 10 {
			log.Printf("   ... 还有 %d 个属性", len(result.AttributesStyle.CSSProperties)-10)
		}
	} else {
		log.Println("   无属性样式")
	}
}

*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// RuleMatch 规则匹配（基于文档中的RuleMatch类型定义）
type RuleMatch struct {
	Rule              *CSSRule `json:"rule"`              // 匹配的CSS规则
	MatchingSelectors []int    `json:"matchingSelectors"` // 匹配的选择器索引
}

// CSSRule CSS规则表示
// 基于文档中的CSSRule类型定义
type CSSRule struct {
	StyleSheetID          string              `json:"styleSheetId,omitempty"`          // 样式表标识符
	SelectorList          *SelectorList       `json:"selectorList,omitempty"`          // 规则选择器数据
	NestingSelectors      []string            `json:"nestingSelectors,omitempty"`      // 祖先样式规则的选择器数组
	Origin                string              `json:"origin,omitempty"`                // 父样式表的来源
	Style                 *CSSStyle           `json:"style,omitempty"`                 // 关联的样式声明
	OriginTreeScopeNodeID int                 `json:"originTreeScopeNodeId,omitempty"` // 构成此规则来源树范围的DOM节点的后端节点ID
	Media                 []CSSMedia          `json:"media,omitempty"`                 // 媒体查询列表数组
	ContainerQueries      []CSSContainerQuery `json:"containerQueries,omitempty"`      // 容器查询列表数组
	Supports              []CSSSupports       `json:"supports,omitempty"`              // @supports CSS at-rule数组
	Layers                []CSSLayer          `json:"layers,omitempty"`                // 级联层数组
	Scopes                []CSSScope          `json:"scopes,omitempty"`                // @scope CSS at-rule数组
	RuleTypes             []string            `json:"ruleTypes,omitempty"`             // 祖先CSSRule类型的数组
	StartingStyles        []CSSStartingStyle  `json:"startingStyles,omitempty"`        // @starting-style CSS at-rule数组
	Navigations           []CSSNavigation     `json:"navigations,omitempty"`           // @navigation CSS at-rule数组
}

// SelectorList 选择器列表
type SelectorList struct {
	Selectors []Value `json:"selectors"` // 列表中的选择器
	Text      string  `json:"text"`      // 规则选择器文本
}

// Value 值（CSS简单选择器的数据）
type Value struct {
	Text        string       `json:"text"`                  // 值文本
	Range       *SourceRange `json:"range,omitempty"`       // 底层资源中的值范围
	Specificity *Specificity `json:"specificity,omitempty"` // 选择器的特异性
}

// Specificity 特异性
type Specificity struct {
	A int `json:"a"` // a组件，表示ID选择器的数量
	B int `json:"b"` // b组件，表示类选择器、属性选择器和伪类的数量
	C int `json:"c"` // c组件，表示类型选择器和伪元素的数量
}

// CSSContainerQuery CSS容器查询规则描述符
type CSSContainerQuery struct {
	Text               string       `json:"text"`                         // 容器查询文本
	Range              *SourceRange `json:"range,omitempty"`              // 关联的规则标题范围
	StyleSheetID       string       `json:"styleSheetId"`                 // 包含此对象的样式表标识符
	Name               string       `json:"name,omitempty"`               // 容器的可选名称
	PhysicalAxes       string       `json:"physicalAxes,omitempty"`       // 容器的物理轴
	LogicalAxes        string       `json:"logicalAxes,omitempty"`        // 容器的逻辑轴
	QueriesScrollState bool         `json:"queriesScrollState,omitempty"` // 查询是否包含scroll-state()查询
	QueriesAnchored    bool         `json:"queriesAnchored,omitempty"`    // 查询是否包含anchored()查询
}

// CSSSupports CSS @supports规则描述符
type CSSSupports struct {
	Text         string       `json:"text"`            // supports规则文本
	Active       bool         `json:"active"`          // 是否满足supports条件
	Range        *SourceRange `json:"range,omitempty"` // 关联的规则标题范围
	StyleSheetID string       `json:"styleSheetId"`    // 包含此对象的样式表标识符
}

// CSSLayer CSS @layer规则描述符
type CSSLayer struct {
	Text         string       `json:"text"`            // 层名称
	Range        *SourceRange `json:"range,omitempty"` // 关联的规则标题范围
	StyleSheetID string       `json:"styleSheetId"`    // 包含此对象的样式表标识符
}

// CSSScope CSS @scope规则描述符
type CSSScope struct {
	Text         string       `json:"text"`            // 范围规则文本
	Range        *SourceRange `json:"range,omitempty"` // 关联的规则标题范围
	StyleSheetID string       `json:"styleSheetId"`    // 包含此对象的样式表标识符
}

// CSSStartingStyle CSS @starting-style规则描述符
type CSSStartingStyle struct {
	Range        *SourceRange `json:"range,omitempty"` // 关联的规则标题范围
	StyleSheetID string       `json:"styleSheetId"`    // 包含此对象的样式表标识符
}

// CSSNavigation CSS @navigation规则描述符
type CSSNavigation struct {
	Text         string       `json:"text"`            // 导航规则文本
	Active       bool         `json:"active"`          // 是否满足导航条件
	Range        *SourceRange `json:"range,omitempty"` // 关联的规则标题范围
	StyleSheetID string       `json:"styleSheetId"`    // 包含此对象的样式表标识符
}

// StyleSheetOrigin 样式表来源
// 基于文档中的StyleSheetOrigin类型定义
const (
	StyleSheetOriginInjected  = "injected"   // 通过扩展注入的样式表
	StyleSheetOriginUserAgent = "user-agent" // 用户代理样式表
	StyleSheetOriginInspector = "inspector"  // 检查器创建的样式表
	StyleSheetOriginRegular   = "regular"    // 常规样式表
)

// PseudoElementMatches 伪元素匹配
type PseudoElementMatches struct {
	PseudoType       string      `json:"pseudoType"`                 // 伪元素类型
	PseudoIdentifier string      `json:"pseudoIdentifier,omitempty"` // 伪元素标识符
	Matches          []RuleMatch `json:"matches"`                    // 匹配的规则
}

// InheritedStyleEntry 继承样式条目
type InheritedStyleEntry struct {
	InlineStyle     *CSSStyle   `json:"inlineStyle,omitempty"` // 内联样式
	MatchedCSSRules []RuleMatch `json:"matchedCSSRules"`       // 匹配的CSS规则
}

// InheritedPseudoElementMatches 继承伪元素匹配
type InheritedPseudoElementMatches struct {
	PseudoElements []PseudoElementMatches `json:"pseudoElements"` // 伪元素匹配
}

// CSSKeyframesRule CSS关键帧规则
type CSSKeyframesRule struct {
	AnimationName interface{}       `json:"animationName"` // 动画名称
	Keyframes     []CSSKeyframeRule `json:"keyframes"`     // 关键帧规则
}

// CSSKeyframeRule CSS关键帧规则表示
// 基于文档中的CSSKeyframeRule类型定义
type CSSKeyframeRule struct {
	StyleSheetID string      `json:"styleSheetId,omitempty"` // 样式表标识符
	Origin       string      `json:"origin,omitempty"`       // 父样式表的来源
	KeyText      interface{} `json:"keyText"`                // 关联的关键文本
	Style        *CSSStyle   `json:"style"`                  // 关联的样式声明
}

// CSSPositionTryRule CSS @position-try规则
type CSSPositionTryRule struct {
	Name         interface{} `json:"name"`         // 名称
	StyleSheetID string      `json:"styleSheetId"` // 样式表ID
	Origin       string      `json:"origin"`       // 来源
	Style        *CSSStyle   `json:"style"`        // 样式
	Active       bool        `json:"active"`       // 是否激活
}

// CSSPropertyRule CSS @property规则
type CSSPropertyRule struct {
	StyleSheetID string      `json:"styleSheetId"` // 样式表ID
	Origin       string      `json:"origin"`       // 来源
	PropertyName interface{} `json:"propertyName"` // 属性名
	Style        *CSSStyle   `json:"style"`        // 样式
}

// CSSPropertyRegistration CSS属性注册
type CSSPropertyRegistration struct {
	PropertyName string      `json:"propertyName"` // 属性名
	InitialValue interface{} `json:"initialValue"` // 初始值
	Inherits     bool        `json:"inherits"`     // 是否继承
	Syntax       string      `json:"syntax"`       // 语法
}

// CSSAtRule CSS @规则
type CSSAtRule struct {
	Type         string      `json:"type"`                 // 规则类型
	Subsection   string      `json:"subsection,omitempty"` // 子节
	Name         interface{} `json:"name,omitempty"`       // 名称
	StyleSheetID string      `json:"styleSheetId"`         // 样式表ID
	Origin       string      `json:"origin"`               // 来源
	Style        *CSSStyle   `json:"style,omitempty"`      // 样式
}

// CSSFunctionRule CSS @function规则
type CSSFunctionRule struct {
	Name         interface{}            `json:"name"`         // 函数名
	StyleSheetID string                 `json:"styleSheetId"` // 样式表ID
	Origin       string                 `json:"origin"`       // 来源
	Parameters   []CSSFunctionParameter `json:"parameters"`   // 参数
	Children     []CSSFunctionNode      `json:"children"`     // 子节点
}

// CSSFunctionParameter CSS函数参数表示
// 基于文档中的CSSFunctionParameter类型定义
type CSSFunctionParameter struct {
	Name string `json:"name"` // 参数名称
	Type string `json:"type"` // 参数类型
}

// CSSFunctionNode CSS函数节点表示
// 基于文档中的CSSFunctionNode类型定义
type CSSFunctionNode struct {
	Condition *CSSFunctionConditionNode `json:"condition,omitempty"` // 条件块
	Style     *CSSStyle                 `json:"style,omitempty"`     // 节点设置的值
}

// CSSFunctionConditionNode CSS函数条件块表示
// 基于文档中的CSSFunctionConditionNode类型定义
type CSSFunctionConditionNode struct {
	Media            *CSSMedia          `json:"media,omitempty"`            // 媒体查询条件块
	ContainerQueries *CSSContainerQuery `json:"containerQueries,omitempty"` // 容器查询条件块
	Supports         *CSSSupports       `json:"supports,omitempty"`         // @supports CSS at-rule条件
	Navigation       *CSSNavigation     `json:"navigation,omitempty"`       // @navigation条件
	Children         []CSSFunctionNode  `json:"children"`                   // 块体
	ConditionText    string             `json:"conditionText"`              // 条件文本
}

// MatchedStylesResult 匹配样式结果
type MatchedStylesResult struct {
	InlineStyle                 *CSSStyle                       `json:"inlineStyle"`                           // 内联样式
	AttributesStyle             *CSSStyle                       `json:"attributesStyle"`                       // 属性样式
	MatchedCSSRules             []RuleMatch                     `json:"matchedCSSRules"`                       // 匹配的CSS规则
	PseudoElements              []PseudoElementMatches          `json:"pseudoElements"`                        // 伪元素匹配
	Inherited                   []InheritedStyleEntry           `json:"inherited"`                             // 继承样式链
	InheritedPseudoElements     []InheritedPseudoElementMatches `json:"inheritedPseudoElements"`               // 继承伪元素
	CSSKeyframesRules           []CSSKeyframesRule              `json:"cssKeyframesRules"`                     // CSS关键帧规则
	CSSPositionTryRules         []CSSPositionTryRule            `json:"cssPositionTryRules"`                   // CSS @position-try规则
	ActivePositionFallbackIndex *int                            `json:"activePositionFallbackIndex,omitempty"` // 活动位置回退索引
	CSSPropertyRules            []CSSPropertyRule               `json:"cssPropertyRules"`                      // CSS @property规则
	CSSPropertyRegistrations    []CSSPropertyRegistration       `json:"cssPropertyRegistrations"`              // CSS属性注册
	CSSAtRules                  []CSSAtRule                     `json:"cssAtRules"`                            // CSS @规则
	ParentLayoutNodeId          *int                            `json:"parentLayoutNodeId,omitempty"`          // 父布局节点ID
	CSSFunctionRules            []CSSFunctionRule               `json:"cssFunctionRules"`                      // CSS @function规则
}

// ParseMatchedStyles 解析匹配样式响应
func ParseMatchedStyles(response string) (*MatchedStylesResult, error) {
	var data struct {
		Result *MatchedStylesResult `json:"result"`
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

// 示例1: 完整的样式来源分析
func exampleCompleteStyleSourceAnalysis() {
	// === 应用场景描述 ===
	// 场景: 开发者工具中的完整样式分析面板
	// 用途: 显示元素所有样式来源，包括继承、伪元素、动画等
	// 优势: 提供最全面的样式调试信息，帮助理解样式层叠和优先级
	// 典型工作流: 选择元素 -> 获取匹配样式 -> 分析所有来源 -> 调试优先级

	elementId := 25

	log.Println("执行完整的样式来源分析...")

	// 获取匹配样式
	response, err := CDPCSSGetMatchedStylesForNode(elementId)
	if err != nil {
		log.Printf("获取匹配样式失败: %v", err)
		return
	}

	// 解析结果
	result, err := ParseMatchedStyles(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	log.Printf("=== 完整的样式来源分析 ===")
	log.Printf("元素ID: %d", elementId)

	// 1. 内联样式
	log.Printf("\n1. 内联样式:")
	if result.InlineStyle != nil && len(result.InlineStyle.CSSProperties) > 0 {
		log.Printf("   属性数量: %d", len(result.InlineStyle.CSSProperties))
		for i, prop := range result.InlineStyle.CSSProperties {
			if i < 5 {
				log.Printf("   - %s: %s", prop.Name, prop.Value)
			}
		}
		if len(result.InlineStyle.CSSProperties) > 5 {
			log.Printf("   ... 还有 %d 个属性", len(result.InlineStyle.CSSProperties)-5)
		}
	} else {
		log.Printf("   无内联样式")
	}

	// 2. 匹配的CSS规则
	log.Printf("\n2. 匹配的CSS规则:")
	log.Printf("   规则数量: %d", len(result.MatchedCSSRules))
	for i, ruleMatch := range result.MatchedCSSRules {
		if i < 3 && ruleMatch.Rule != nil && ruleMatch.Rule.SelectorList != nil {
			log.Printf("   规则[%d]: %s", i+1, ruleMatch.Rule.SelectorList.Text)
		}
	}
	if len(result.MatchedCSSRules) > 3 {
		log.Printf("   ... 还有 %d 个规则", len(result.MatchedCSSRules)-3)
	}

	// 3. 伪元素样式
	log.Printf("\n3. 伪元素样式:")
	if len(result.PseudoElements) > 0 {
		log.Printf("   伪元素数量: %d", len(result.PseudoElements))
		for _, pseudo := range result.PseudoElements {
			log.Printf("   - %s: %d 个匹配规则", pseudo.PseudoType, len(pseudo.Matches))
		}
	} else {
		log.Printf("   无伪元素样式")
	}

	// 4. 继承样式链
	log.Printf("\n4. 继承样式链:")
	if len(result.Inherited) > 0 {
		log.Printf("   继承层级: %d", len(result.Inherited))
		for i, inherited := range result.Inherited {
			log.Printf("   层级[%d]: %d 个匹配规则", i+1, len(inherited.MatchedCSSRules))
		}
	} else {
		log.Printf("   无继承样式")
	}

	// 5. CSS动画规则
	log.Printf("\n5. CSS动画规则:")
	if len(result.CSSKeyframesRules) > 0 {
		log.Printf("   动画数量: %d", len(result.CSSKeyframesRules))
		for _, animation := range result.CSSKeyframesRules {
			log.Printf("   - 动画: %v", animation.AnimationName)
		}
	} else {
		log.Printf("   无CSS动画")
	}

	// 6. CSS @规则
	log.Printf("\n6. CSS @规则:")
	if len(result.CSSAtRules) > 0 {
		log.Printf("   @规则数量: %d", len(result.CSSAtRules))
		for _, atRule := range result.CSSAtRules {
			log.Printf("   - @%s", atRule.Type)
		}
	} else {
		log.Printf("   无CSS @规则")
	}
}


// 示例4: 网站样式审计完整工作流
func exampleWebsiteStyleAuditWorkflow() {
	// === 应用场景描述 ===
	// 场景: 网站样式质量和性能的完整审计
	// 用途: 分析网站中关键元素的样式实现，评估样式质量和性能
	// 优势: 系统化评估样式实现，提供优化建议
	// 典型工作流: 选择关键元素 -> 分析匹配样式 -> 评估质量 -> 生成报告

	log.Println("=== 网站样式审计工作流 ===")

	// 1. 启用CSS域
	log.Println("步骤1: 启用CSS域")
	_, err := CDPCSSEnable()
	if err != nil {
		log.Printf("启用CSS域失败: %v", err)
		return
	}
	defer CDPCSSDisable()

	// 2. 选择关键审计元素
	log.Println("步骤2: 选择关键审计元素")
	// 通常包括：按钮、表单元素、导航、主要内容区域等
	auditElements := []struct {
		id   int
		name string
	}{
		{25, "主导航"},
		{42, "主按钮"},
		{33, "主要内容区域"},
		{67, "表单输入框"},
		{89, "页脚"},
	}

	// 3. 分析每个元素的匹配样式
	log.Println("步骤3: 分析元素匹配样式")
	analyzers := make([]*MatchedStylesAnalyzer, 0, len(auditElements))

	for _, element := range auditElements {
		log.Printf("分析元素: %s (ID: %d)", element.name, element.id)

		analyzer := NewMatchedStylesAnalyzer(element.id)
		if err := analyzer.Analyze(); err != nil {
			log.Printf("  分析失败: %v", err)
			continue
		}

		analyzers = append(analyzers, analyzer)
		log.Printf("  分析完成")
	}

	// 4. 生成审计报告
	log.Println("步骤4: 生成样式审计报告")
	generateStyleAuditReport(auditElements, analyzers)

	log.Println("网站样式审计工作流完成")
}

// 生成样式审计报告
func generateStyleAuditReport(elements []struct{id int; name string}, analyzers []*MatchedStylesAnalyzer) {
	log.Println("\n=== 网站样式审计报告 ===")

	totalElements := len(analyzers)
	elementsWithInlineStyles := 0
	totalMatchedRules := 0
	elementsWithImportant := 0
	elementsWithAnimations := 0

	for i, analyzer := range analyzers {
		element := elements[i]

		log.Printf("\n[%s] (ID: %d)", element.name, element.id)
		log.Printf("  - 匹配规则数: %d", analyzer.GetMatchedRulesCount())
		log.Printf("  - 内联样式属性: %d", len(analyzer.GetInlineStyleProperties()))
		log.Printf("  - 伪元素: %d", len(analyzer.GetPseudoElements()))
		log.Printf("  - 继承深度: %d", analyzer.GetInheritanceDepth())
		log.Printf("  - 动画数量: %d", len(analyzer.GetAnimations()))

		// 统计
		if len(analyzer.GetInlineStyleProperties()) > 0 {
			elementsWithInlineStyles++
		}
		totalMatchedRules += analyzer.GetMatchedRulesCount()
		if analyzer.HasImportantDeclarations() {
			elementsWithImportant++
		}
		if len(analyzer.GetAnimations()) > 0 {
			elementsWithAnimations++
		}
	}

	// 生成总结报告
	log.Println("\n=== 审计总结 ===")
	log.Printf("分析元素总数: %d", totalElements)
	log.Printf("使用内联样式的元素: %d (%.1f%%)",
		elementsWithInlineStyles,
		float64(elementsWithInlineStyles)/float64(totalElements)*100)
	log.Printf("包含!important的元素: %d (%.1f%%)",
		elementsWithImportant,
		float64(elementsWithImportant)/float64(totalElements)*100)
	log.Printf("使用动画的元素: %d (%.1f%%)",
		elementsWithAnimations,
		float64(elementsWithAnimations)/float64(totalElements)*100)
	log.Printf("平均每个元素匹配规则数: %.1f", float64(totalMatchedRules)/float64(totalElements))

	// 优化建议
	log.Println("\n=== 优化建议 ===")
	if elementsWithInlineStyles > totalElements/3 {
		log.Println("⚠ 警告: 较多元素使用内联样式，建议:")
		log.Println("  - 减少内联样式的使用")
		log.Println("  - 将公共样式提取到CSS类")
	}

	if elementsWithImportant > 0 {
		log.Println("⚠ 警告: 发现!important声明，建议:")
		log.Println("  - 尽量避免使用!important")
		log.Println("  - 通过提高选择器特异性来替代!important")
	}

	avgRules := float64(totalMatchedRules) / float64(totalElements)
	if avgRules > 20 {
		log.Println("⚠ 警告: 样式规则较多，建议:")
		log.Println("  - 合并相似的样式规则")
		log.Println("  - 减少选择器复杂性")
		log.Println("  - 考虑使用CSS-in-JS或CSS模块")
	}
}

*/

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
		"method": "CSS.getMediaQueries"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

/*

// 示例1: 响应式设计调试工具
func exampleResponsiveDesignDebugger() {
	// === 应用场景描述 ===
	// 场景: 响应式设计调试工具
	// 用途: 获取页面中所有媒体查询，分析响应式断点的激活状态
	// 优势: 实时查看媒体查询状态，调试响应式布局
	// 典型工作流: 获取媒体查询 -> 分析激活状态 -> 调整视口 -> 验证变化

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

*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// PlatformFontUsage 平台字体使用统计
// 基于文档中的PlatformFontUsage类型定义
type PlatformFontUsage struct {
	FamilyName     string  `json:"familyName"`     // 平台报告的字体族名
	PostScriptName string  `json:"postScriptName"` // 平台报告的PostScript名称
	IsCustomFont   bool    `json:"isCustomFont"`   // 是否是下载的或本地解析的自定义字体
	GlyphCount     float64 `json:"glyphCount"`     // 使用此字体渲染的字形数量
}

// PlatformFontsResult 平台字体结果
type PlatformFontsResult struct {
	Fonts []PlatformFontUsage `json:"fonts"` // 平台字体使用数组
}

// ParsePlatformFonts 解析平台字体响应
func ParsePlatformFonts(response string) (*PlatformFontsResult, error) {
	var data struct {
		Result *PlatformFontsResult `json:"result"`
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

示例
// 示例1: 字体使用分析和优化
func exampleFontUsageAnalysis() {
	// === 应用场景描述 ===
	// 场景: 字体使用分析和性能优化
	// 用途: 分析元素中实际使用的字体，识别性能瓶颈
	// 优势: 了解实际字体渲染情况，优化字体加载策略
	// 典型工作流: 选择元素 -> 获取平台字体 -> 分析使用情况 -> 优化建议

	textElementId := 52

	log.Println("分析元素字体使用情况...")

	// 获取平台字体
	response, err := CDPCSSGetPlatformFontsForNode(textElementId)
	if err != nil {
		log.Printf("获取平台字体失败: %v", err)
		return
	}

	// 解析结果
	result, err := ParsePlatformFonts(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	log.Printf("=== 字体使用分析 ===")
	log.Printf("元素ID: %d", textElementId)
	log.Printf("使用的字体数量: %d", len(result.Fonts))

	totalGlyphs := 0.0
	customFontCount := 0
	systemFontCount := 0

	for i, font := range result.Fonts {
		log.Printf("\n字体[%d]:", i+1)
		log.Printf("  字体族名: %s", font.FamilyName)

		if font.PostScriptName != "" {
			log.Printf("  PostScript名称: %s", font.PostScriptName)
		}

		fontType := "系统字体"
		if font.IsCustomFont {
			fontType = "自定义字体"
			customFontCount++
		} else {
			systemFontCount++
		}
		log.Printf("  字体类型: %s", fontType)

		log.Printf("  渲染字形数: %.0f", font.GlyphCount)
		totalGlyphs += font.GlyphCount

		// 计算使用比例
		if totalGlyphs > 0 {
			percentage := (font.GlyphCount / totalGlyphs) * 100
			log.Printf("  使用比例: %.1f%%", percentage)
		}
	}

	// 生成统计报告
	log.Printf("\n=== 字体使用统计 ===")
	log.Printf("总字形数: %.0f", totalGlyphs)
	log.Printf("自定义字体数量: %d", customFontCount)
	log.Printf("系统字体数量: %d", systemFontCount)

	if customFontCount > 0 {
		log.Println("⚠ 发现自定义字体，建议:")
		log.Println("  - 检查自定义字体的加载性能")
		log.Println("  - 考虑使用font-display优化加载行为")
		log.Println("  - 实现字体回退策略")
	}
}


*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// StyleSheetTextResult 样式表文本结果
type StyleSheetTextResult struct {
	Text string `json:"text"` // 样式表文本
}

// ParseStyleSheetText 解析样式表文本响应
func ParseStyleSheetText(response string) (*StyleSheetTextResult, error) {
	var data struct {
		Result *StyleSheetTextResult `json:"result"`
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

示例

// 示例1: 样式表内容分析和调试
func exampleStyleSheetAnalysis() {
	// === 应用场景描述 ===
	// 场景: 样式表内容分析和调试
	// 用途: 获取样式表的完整内容，进行分析和调试
	// 优势: 可以查看原始CSS代码，便于调试样式问题
	// 典型工作流: 获取样式表ID -> 获取样式表文本 -> 分析内容 -> 调试问题

	styleSheetId := "style-sheet-1"

	log.Println("获取样式表内容进行分析...")

	// 获取样式表文本
	response, err := CDPCSSGetStyleSheetText(styleSheetId)
	if err != nil {
		log.Printf("获取样式表文本失败: %v", err)
		return
	}

	// 解析结果
	result, err := ParseStyleSheetText(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	log.Printf("=== 样式表分析 ===")
	log.Printf("样式表ID: %s", styleSheetId)

	// 分析样式表内容
	cssText := result.Text
	log.Printf("样式表长度: %d 字符", len(cssText))
	log.Printf("样式表行数: %d", countLines(cssText))

	// 显示样式表前几行
	log.Println("\n样式表预览:")
	lines := strings.SplitN(cssText, "\n", 10)
	for i, line := range lines {
		if i < 9 { // 显示前9行
			log.Printf("  [%d] %s", i+1, line)
		}
	}
	if len(cssText) > 0 {
		log.Printf("  ... (内容截断)")
	}

	// 分析样式表结构
	log.Println("\n=== 样式表结构分析 ===")
	analyzeStyleSheetStructure(cssText)
}

// 统计行数
func countLines(text string) int {
	return strings.Count(text, "\n") + 1
}

// 分析样式表结构
func analyzeStyleSheetStructure(cssText string) {
	// 统计规则类型
	ruleCounts := make(map[string]int)

	// 简单的正则匹配
	patterns := map[string]*regexp.Regexp{
		"class_selector":   regexp.MustCompile(`\.[a-zA-Z_][a-zA-Z0-9_-]*\s*\{`),
		"id_selector":      regexp.MustCompile(`#[a-zA-Z_][a-zA-Z0-9_-]*\s*\{`),
		"element_selector": regexp.MustCompile(`^[a-zA-Z_]+\s*\{`),
		"at_rule":          regexp.MustCompile(`@(media|keyframes|import|font-face|supports|page|charset)`),
		"pseudo_class":     regexp.MustCompile(`:[a-zA-Z-]+\s*\{`),
		"pseudo_element":   regexp.MustCompile(`::[a-zA-Z-]+\s*\{`),
	}

	for ruleType, pattern := range patterns {
		matches := pattern.FindAllString(cssText, -1)
		ruleCounts[ruleType] = len(matches)
	}

	// 显示统计结果
	log.Println("规则类型统计:")
	log.Printf("  类选择器: %d", ruleCounts["class_selector"])
	log.Printf("  ID选择器: %d", ruleCounts["id_selector"])
	log.Printf("  元素选择器: %d", ruleCounts["element_selector"])
	log.Printf("  @规则: %d", ruleCounts["at_rule"])
	log.Printf("  伪类: %d", ruleCounts["pseudo_class"])
	log.Printf("  伪元素: %d", ruleCounts["pseudo_element"])

	// 统计注释
	commentPattern := regexp.MustCompile(`/\*[\s\S]*?\*\/`)
	comments := commentPattern.FindAllString(cssText, -1)
	log.Printf("  注释块: %d", len(comments))

	// 统计属性
	propertyPattern := regexp.MustCompile(`[a-zA-Z-]+\s*:`)
	properties := propertyPattern.FindAllString(cssText, -1)
	log.Printf("  属性声明: %d", len(properties))

	// 计算注释率
	totalLines := countLines(cssText)
	commentLines := 0
	for _, comment := range comments {
		commentLines += countLines(comment)
	}

	if totalLines > 0 {
		commentRate := float64(commentLines) / float64(totalLines) * 100
		log.Printf("  注释率: %.1f%%", commentRate)
	}
}

// 示例2: 样式表对比和差异分析
func exampleStyleSheetComparison() {
	// === 应用场景描述 ===
	// 场景: 样式表对比和版本差异分析
	// 用途: 对比不同版本或不同环境的样式表差异
	// 优势: 快速识别样式变化，便于版本控制和回归测试
	// 典型工作流: 获取两个样式表 -> 对比内容 -> 识别差异 -> 分析影响

	// 假设我们有两个样式表ID
	styleSheetId1 := "style-sheet-v1"
	styleSheetId2 := "style-sheet-v2"

	log.Println("对比两个样式表版本...")

	// 获取第一个样式表
	response1, err := CDPCSSGetStyleSheetText(styleSheetId1)
	if err != nil {
		log.Printf("获取第一个样式表失败: %v", err)
		return
	}

	result1, err := ParseStyleSheetText(response1)
	if err != nil {
		log.Printf("解析第一个样式表失败: %v", err)
		return
	}

	// 获取第二个样式表
	response2, err := CDPCSSGetStyleSheetText(styleSheetId2)
	if err != nil {
		log.Printf("获取第二个样式表失败: %v", err)
		return
	}

	result2, err := ParseStyleSheetText(response2)
	if err != nil {
		log.Printf("解析第二个样式表失败: %v", err)
		return
	}

	log.Printf("=== 样式表对比分析 ===")
	log.Printf("样式表1: %s (长度: %d)", styleSheetId1, len(result1.Text))
	log.Printf("样式表2: %s (长度: %d)", styleSheetId2, len(result2.Text))

	// 简单对比
	if result1.Text == result2.Text {
		log.Println("✓ 两个样式表完全相同")
		return
	}

	// 计算差异
	diff := calculateStyleSheetDiff(result1.Text, result2.Text)
	log.Printf("差异分析:")
	log.Printf("  总行数差异: %d 行", diff.lineDiff)
	log.Printf("  新增规则: %d 个", diff.addedRules)
	log.Printf("  删除规则: %d 个", diff.removedRules)
	log.Printf("  修改规则: %d 个", diff.modifiedRules)

	// 显示主要差异
	if len(diff.majorChanges) > 0 {
		log.Println("  主要变化:")
		for i, change := range diff.majorChanges {
			if i < 5 { // 只显示前5个
				log.Printf("    - %s", change)
			}
		}
		if len(diff.majorChanges) > 5 {
			log.Printf("    ... 还有 %d 个变化", len(diff.majorChanges)-5)
		}
	}
}

// 样式表差异结果
type StyleSheetDiff struct {
	lineDiff      int
	addedRules    int
	removedRules  int
	modifiedRules int
	majorChanges  []string
}

// 计算样式表差异
func calculateStyleSheetDiff(text1, text2 string) StyleSheetDiff {
	diff := StyleSheetDiff{}

	// 行数差异
	lines1 := strings.Split(text1, "\n")
	lines2 := strings.Split(text2, "\n")
	diff.lineDiff = len(lines2) - len(lines1)

	// 简化的差异检测
	// 提取规则
	rules1 := extractCSSRules(text1)
	rules2 := extractCSSRules(text2)

	// 比较规则
	ruleSet1 := make(map[string]bool)
	ruleSet2 := make(map[string]bool)

	for _, rule := range rules1 {
		ruleSet1[rule] = true
	}
	for _, rule := range rules2 {
		ruleSet2[rule] = true
	}

	// 计算新增和删除的规则
	for rule := range ruleSet2 {
		if !ruleSet1[rule] {
			diff.addedRules++
		}
	}
	for rule := range ruleSet1 {
		if !ruleSet2[rule] {
			diff.removedRules++
		}
	}

	// 识别主要变化
	diff.majorChanges = identifyMajorChanges(text1, text2)

	return diff
}

// 提取CSS规则
func extractCSSRules(cssText string) []string {
	// 简化的规则提取
	rulePattern := regexp.MustCompile(`([.#]?[a-zA-Z_-][^{]*)\s*\{`)
	matches := rulePattern.FindAllStringSubmatch(cssText, -1)

	var rules []string
	for _, match := range matches {
		if len(match) > 1 {
			rules = append(rules, strings.TrimSpace(match[1]))
		}
	}
	return rules
}

// 识别主要变化
func identifyMajorChanges(text1, text2 string) []string {
	var changes []string

	// 检测媒体查询变化
	mediaPattern := regexp.MustCompile(`@media[^{]+\{`)
	media1 := mediaPattern.FindAllString(text1, -1)
	media2 := mediaPattern.FindAllString(text2, -1)

	if len(media1) != len(media2) {
		changes = append(changes, fmt.Sprintf("媒体查询数量变化: %d -> %d", len(media1), len(media2)))
	}

	// 检测关键帧动画变化
	keyframesPattern := regexp.MustCompile(`@keyframes[^{]+\{`)
	keyframes1 := keyframesPattern.FindAllString(text1, -1)
	keyframes2 := keyframesPattern.FindAllString(text2, -1)

	if len(keyframes1) != len(keyframes2) {
		changes = append(changes, fmt.Sprintf("关键帧动画数量变化: %d -> %d", len(keyframes1), len(keyframes2)))
	}

	return changes
}


*/

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
func CDPCSSSetEffectivePropertyValueForNode(nodeId int, propertyName, value string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if nodeId <= 0 {
		return "", fmt.Errorf("节点ID必须是正整数")
	}
	if propertyName == "" {
		return "", fmt.Errorf("属性名不能为空")
	}
	if value == "" {
		return "", fmt.Errorf("属性值不能为空")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 转义特殊字符
	escapedPropertyName := strings.ReplaceAll(propertyName, `"`, `\"`)
	escapedValue := strings.ReplaceAll(value, `"`, `\"`)
	escapedValue = strings.ReplaceAll(escapedValue, "\n", "\\n")

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setEffectivePropertyValueForNode",
		"params": {
			"nodeId": %d,
			"propertyName": "%s",
			"value": "%s"
		}
	}`, reqID, nodeId, escapedPropertyName, escapedValue)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

/*

示例

// 示例1: 实时样式调试工具
func exampleLiveStyleDebugger() {
	// === 应用场景描述 ===
	// 场景: 实时样式调试工具
	// 用途: 在开发者工具中实时修改元素的CSS属性值
	// 优势: 无需修改源代码即可预览样式变化，提高调试效率
	// 典型工作流: 选择元素 -> 修改属性值 -> 实时预览 -> 应用或撤销

	elementId := 25

	log.Println("实时样式调试 - 修改元素样式属性...")

	// 定义要测试的样式修改
	styleTests := []struct {
		propertyName string
		value        string
		description  string
	}{
		{"color", "red", "修改文本颜色为红色"},
		{"font-size", "20px", "增大字体大小"},
		{"background-color", "yellow", "修改背景色为黄色"},
		{"border", "2px solid blue", "添加蓝色边框"},
		{"padding", "10px", "增加内边距"},
		{"margin", "20px", "增加外边距"},
		{"opacity", "0.7", "设置透明度"},
		{"display", "flex", "修改为flex布局"},
		{"justify-content", "center", "水平居中"},
		{"align-items", "center", "垂直居中"},
	}

	for i, test := range styleTests {
		log.Printf("测试 %d/%d: %s", i+1, len(styleTests), test.description)
		log.Printf("  修改: %s: %s", test.propertyName, test.value)

		// 设置属性值
		response, err := CDPCSSSetEffectivePropertyValueForNode(elementId, test.propertyName, test.value)
		if err != nil {
			log.Printf("  修改失败: %v", err)
		} else {
			log.Printf("  修改成功")
		}

		// 等待一段时间查看效果
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("所有样式修改测试完成")
}

// 示例2: 响应式设计预览工具
func exampleResponsiveDesignPreview() {
	// === 应用场景描述 ===
	// 场景: 响应式设计预览和调试工具
	// 用途: 在不同断点下预览和调试元素的样式
	// 优势: 实时调整响应式样式，优化多设备体验
	// 典型工作流: 设置视口大小 -> 调整样式 -> 测试断点 -> 优化响应式设计

	containerId := 33

	log.Println("响应式设计预览 - 调整容器样式...")

	// 定义不同视口下的样式调整
	viewportStyles := []struct {
		viewport     string
		width        int
		adjustments  []struct {
			property string
			value    string
		}
	}{
		{
			viewport: "手机端 (<768px)",
			width:    375,
			adjustments: []struct {
				property string
				value    string
			}{
				{"padding", "10px"},
				{"font-size", "14px"},
				{"flex-direction", "column"},
			},
		},
		{
			viewport: "平板端 (768px-1024px)",
			width:    768,
			adjustments: []struct {
				property string
				value    string
			}{
				{"padding", "20px"},
				{"font-size", "16px"},
				{"flex-direction", "row"},
			},
		},
		{
			viewport: "桌面端 (>1024px)",
			width:    1200,
			adjustments: []struct {
				property string
				value    string
			}{
				{"padding", "30px"},
				{"font-size", "18px"},
				{"max-width", "1200px"},
				{"margin", "0 auto"},
			},
		},
	}

	for _, viewport := range viewportStyles {
		log.Printf("\n视口: %s (%dpx)", viewport.viewport, viewport.width)

		// 这里可以设置视口大小
		// CDPPageSetDeviceMetricsOverride(viewport.width, 800)
		// time.Sleep(200 * time.Millisecond)

		for i, adjustment := range viewport.adjustments {
			log.Printf("  调整[%d]: %s: %s", i+1, adjustment.property, adjustment.value)

			response, err := CDPCSSSetEffectivePropertyValueForNode(containerId, adjustment.property, adjustment.value)
			if err != nil {
				log.Printf("    失败: %v", err)
			} else {
				log.Printf("    成功")
			}

			time.Sleep(200 * time.Millisecond)
		}
	}

	log.Println("响应式设计预览完成")
}

// 示例3: 设计系统样式验证
func exampleDesignSystemValidation() {
	// === 应用场景描述 ===
	// 场景: 设计系统的样式规范验证
	// 用途: 验证元素样式是否符合设计系统规范
	// 优势: 确保UI组件的一致性和设计规范符合性
	// 典型工作流: 选择组件 -> 检查样式 -> 调整不符合项 -> 验证规范

	buttonId := 42

	log.Println("设计系统验证 - 检查按钮样式规范...")

	// 定义设计系统规范
	designSystemSpecs := map[string]string{
		"padding":         "12px 24px",
		"border-radius":   "4px",
		"font-weight":     "500",
		"font-size":       "14px",
		"line-height":     "1.5",
		"border":          "none",
		"cursor":         "pointer",
		"transition":      "all 0.2s ease",
	}

	log.Println("检查并调整按钮样式以符合设计系统规范...")

	// 应用设计系统规范
	adjustmentsMade := 0
	totalProperties := len(designSystemSpecs)

	for property, expectedValue := range designSystemSpecs {
		log.Printf("设置 %s: %s", property, expectedValue)

		response, err := CDPCSSSetEffectivePropertyValueForNode(buttonId, property, expectedValue)
		if err != nil {
			log.Printf("  设置失败: %v", err)
		} else {
			log.Printf("  设置成功")
			adjustmentsMade++
		}

		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("\n=== 设计系统验证结果 ===")
	log.Printf("总属性数: %d", totalProperties)
	log.Printf("成功设置: %d", adjustmentsMade)

	if adjustmentsMade == totalProperties {
		log.Println("✓ 所有样式属性已成功应用")
	} else {
		log.Printf("⚠ 部分属性设置失败: %d/%d", totalProperties-adjustmentsMade, totalProperties)
	}
}

*/

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
func CDPCSSSetKeyframeKey(styleSheetId string, rangeObj SourceRange, keyText string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}
	if keyText == "" {
		return "", fmt.Errorf("关键文本不能为空")
	}
	// 验证keyText格式，应为有效的关键帧位置
	if !isValidKeyframeKey(keyText) {
		return "", fmt.Errorf("无效的关键帧位置: %s，应为百分比(如'0 %%')或关键词(如'from', 'to')", keyText)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建范围JSON
	rangeJSON, err := json.Marshal(rangeObj)
	if err != nil {
		return "", fmt.Errorf("序列化范围失败: %w", err)
	}

	// 转义特殊字符
	escapedKeyText := strings.ReplaceAll(keyText, `"`, `\"`)

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setKeyframeKey",
		"params": {
			"styleSheetId": "%s",
			"range": %s,
			"keyText": "%s"
		}
	}`, reqID, styleSheetId, rangeJSON, escapedKeyText)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// 验证关键帧位置格式
func isValidKeyframeKey(keyText string) bool {
	keyText = strings.TrimSpace(keyText)

	// 检查是否为有效百分比
	if strings.HasSuffix(keyText, "%") {
		percentStr := strings.TrimSuffix(keyText, "%")
		if percent, err := strconv.ParseFloat(percentStr, 64); err == nil {
			return percent >= 0 && percent <= 100
		}
		return false
	}

	// 检查是否为关键字
	validKeywords := []string{"from", "to"}
	for _, keyword := range validKeywords {
		if strings.EqualFold(keyText, keyword) {
			return true
		}
	}

	return false
}

/*

示例

// 示例1: CSS动画编辑器
func exampleCSSAnimationEditor() {
	// === 应用场景描述 ===
	// 场景: CSS动画编辑器工具
	// 用途: 可视化的CSS关键帧动画编辑
	// 优势: 实时修改动画关键帧，预览效果
	// 典型工作流: 选择动画 -> 修改关键帧位置 -> 预览效果 -> 导出CSS

	styleSheetId := "style-sheet-1"

	log.Println("CSS动画编辑器 - 修改关键帧位置...")

	// 假设我们要修改一个已有的关键帧动画
	// 定义关键帧动画范围
	animationRange := SourceRange{
		StartLine:   10,
		StartColumn: 5,
		EndLine:     15,
		EndColumn:   20,
	}

	// 定义关键帧修改测试
	keyframeTests := []struct {
		oldKey     string
		newKey     string
		description string
	}{
		{"0%", "10%", "将起始关键帧从0%调整到10%"},
		{"50%", "40%", "将中间关键帧从50%调整到40%"},
		{"100%", "90%", "将结束关键帧从100%调整到90%"},
		{"from", "0%", "将'from'关键字改为'0%'"},
		{"to", "100%", "将'to'关键字改为'100%'"},
		{"25%", "30%", "将25%关键帧调整到30%"},
	}

	for i, test := range keyframeTests {
		log.Printf("测试 %d/%d: %s", i+1, len(keyframeTests), test.description)
		log.Printf("  修改: %s -> %s", test.oldKey, test.newKey)

		// 设置新的关键帧位置
		response, err := CDPCSSSetKeyframeKey(styleSheetId, animationRange, test.newKey)
		if err != nil {
			log.Printf("  修改失败: %v", err)
		} else {
			log.Printf("  修改成功")

			// 解析返回的关键文本
			var data struct {
				Result struct {
					KeyText interface{} `json:"keyText"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil {
				log.Printf("  返回的关键文本: %v", data.Result.KeyText)
			}
		}

		// 等待一段时间查看效果
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("关键帧修改测试完成")
}

*/

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
func CDPCSSSetMediaText(styleSheetId string, rangeObj SourceRange, text string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}
	if text == "" {
		return "", fmt.Errorf("媒体查询文本不能为空")
	}
	// 验证媒体查询文本格式
	if !isValidMediaQuery(text) {
		return "", fmt.Errorf("无效的媒体查询文本: %s", text)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建范围JSON
	rangeJSON, err := json.Marshal(rangeObj)
	if err != nil {
		return "", fmt.Errorf("序列化范围失败: %w", err)
	}

	// 转义特殊字符
	escapedText := strings.ReplaceAll(text, `"`, `\"`)
	escapedText = strings.ReplaceAll(escapedText, "\n", "\\n")

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setMediaText",
		"params": {
			"styleSheetId": "%s",
			"range": %s,
			"text": "%s"
		}
	}`, reqID, styleSheetId, rangeJSON, escapedText)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setMediaText 请求失败: %w", err)
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
			return "", fmt.Errorf("setMediaText 请求超时")
		}
	}
}

// 验证媒体查询格式
func isValidMediaQuery(text string) bool {
	text = strings.TrimSpace(text)

	// 基本格式检查
	if text == "" {
		return false
	}

	// 检查是否以@media开头
	if strings.HasPrefix(strings.ToLower(text), "@media") {
		return true
	}

	// 检查是否包含媒体特性
	mediaFeatures := []string{
		"min-width", "max-width", "width",
		"min-height", "max-height", "height",
		"orientation", "aspect-ratio", "resolution",
		"color", "monochrome", "grid", "hover", "pointer",
	}

	for _, feature := range mediaFeatures {
		if strings.Contains(text, feature) {
			return true
		}
	}

	// 检查逻辑运算符
	logicOperators := []string{"and", "or", "not", "only"}
	for _, op := range logicOperators {
		if strings.Contains(strings.ToLower(text), " "+op+" ") {
			return true
		}
	}

	return false
}

/*

示例

// 示例1: 响应式设计编辑器
func exampleResponsiveDesignEditor() {
	// === 应用场景描述 ===
	// 场景: 响应式设计编辑器工具
	// 用途: 实时编辑和调整媒体查询条件
	// 优势: 即时预览不同断点下的布局变化
	// 典型工作流: 选择媒体查询 -> 编辑条件 -> 预览效果 -> 保存修改

	styleSheetId := "style-sheet-1"

	log.Println("响应式设计编辑器 - 修改媒体查询...")

	// 定义要修改的媒体查询范围
	mediaRange := SourceRange{
		StartLine:   15,
		StartColumn: 1,
		EndLine:     18,
		EndColumn:   20,
	}

	// 定义媒体查询修改测试
	mediaTests := []struct {
		oldMedia   string
		newMedia   string
		description string
	}{
		{
			"(max-width: 768px)",
			"(max-width: 1024px)",
			"将移动端断点从768px调整到1024px",
		},
		{
			"(min-width: 769px) and (max-width: 1024px)",
			"(min-width: 1025px) and (max-width: 1440px)",
			"将平板端断点范围调整",
		},
		{
			"screen and (max-width: 768px)",
			"screen and (max-width: 768px) and (orientation: portrait)",
			"为移动端添加竖屏方向限制",
		},
		{
			"(min-width: 1200px)",
			"(min-width: 1200px) and (prefers-reduced-motion: no-preference)",
			"为桌面端添加减少动画偏好检测",
		},
		{
			"print",
			"print, (max-width: 480px)",
			"扩展打印样式到小屏幕设备",
		},
	}

	for i, test := range mediaTests {
		log.Printf("测试 %d/%d: %s", i+1, len(mediaTests), test.description)
		log.Printf("  修改: %s -> %s", test.oldMedia, test.newMedia)

		// 设置新的媒体查询文本
		response, err := CDPCSSSetMediaText(styleSheetId, mediaRange, test.newMedia)
		if err != nil {
			log.Printf("  修改失败: %v", err)
		} else {
			log.Printf("  修改成功")

			// 解析返回的媒体规则
			var data struct {
				Result struct {
					Media *CSSMedia `json:"media"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.Media != nil {
				log.Printf("  返回的媒体查询: %s", data.Result.Media.Text)
				log.Printf("  激活状态: %v", isMediaActive(data.Result.Media))
			}
		}

		// 等待一段时间查看效果
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("媒体查询修改测试完成")
}

// 检查媒体查询是否激活
func isMediaActive(media *CSSMedia) bool {
	if media == nil || len(media.MediaList) == 0 {
		return false
	}

	// 检查媒体查询列表中的所有条件
	for _, query := range media.MediaList {
		if !query.Active {
			return false
		}
	}
	return true
}

// 示例2: 跨设备测试工具
func exampleCrossDeviceTesting() {
	// === 应用场景描述 ===
	// 场景: 跨设备响应式测试工具
	// 用途: 测试和调整不同设备下的媒体查询
	// 优势: 确保网站在各种设备上都有良好的表现
	// 典型工作流: 选择设备 -> 调整媒体查询 -> 测试兼容性 -> 优化断点

	styleSheetId := "style-sheet-2"
	mediaRange := SourceRange{
		StartLine:   22,
		StartColumn: 3,
		EndLine:     25,
		EndColumn:   15,
	}

	log.Println("跨设备测试 - 优化媒体查询...")

	// 定义设备特定的媒体查询优化
	deviceOptimizations := []struct {
		device    string
		viewport  struct{ width, height int }
		oldMedia  string
		newMedia  string
		reason    string
	}{
		{
			device:   "iPhone SE",
			viewport: struct{ width, height int }{375, 667},
			oldMedia: "(max-width: 375px)",
			newMedia: "(max-width: 414px)",
			reason:   "覆盖更多小屏设备",
		},
		{
			device:   "iPad",
			viewport: struct{ width, height int }{768, 1024},
			oldMedia: "(min-width: 768px) and (max-width: 1024px)",
			newMedia: "(min-width: 768px) and (max-width: 1024px) and (orientation: landscape)",
			reason:   "针对横屏iPad优化",
		},
		{
			device:   "Desktop HD",
			viewport: struct{ width, height int }{1920, 1080},
			oldMedia: "(min-width: 1200px)",
			newMedia: "(min-width: 1200px) and (min-resolution: 2dppx)",
			reason:   "针对高分辨率屏幕优化",
		},
		{
			device:   "Large Tablet",
			viewport: struct{ width, height int }{1024, 1366},
			oldMedia: "(min-width: 1025px) and (max-width: 1366px)",
			newMedia: "(min-width: 1024px) and (max-width: 1366px)",
			reason:   "包含1024px平板",
		},
	}

	for _, optimization := range deviceOptimizations {
		log.Printf("\n设备: %s (%dx%d)",
			optimization.device, optimization.viewport.width, optimization.viewport.height)
		log.Printf("优化原因: %s", optimization.reason)
		log.Printf("修改: %s -> %s", optimization.oldMedia, optimization.newMedia)

		// 设置新的媒体查询
		response, err := CDPCSSSetMediaText(styleSheetId, mediaRange, optimization.newMedia)
		if err != nil {
			log.Printf("  优化失败: %v", err)
		} else {
			log.Printf("  优化成功")

			// 解析返回结果
			var data struct {
				Result struct {
					Media *CSSMedia `json:"media"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.Media != nil {
				log.Printf("  当前媒体查询: %s", data.Result.Media.Text)
			}
		}

		// 模拟设备视口
		// CDPPageSetDeviceMetricsOverride(optimization.viewport.width, optimization.viewport.height)
		time.Sleep(1 * time.Second)
		log.Println("  设备预览完成")
	}

	log.Println("跨设备测试完成")
}

// 示例3: 媒体查询性能分析器
func exampleMediaQueryPerformanceAnalyzer() {
	// === 应用场景描述 ===
	// 场景: 媒体查询性能分析和优化工具
	// 用途: 分析复杂媒体查询的性能影响并优化
	// 优势: 识别性能瓶颈，优化复杂媒体查询
	// 典型工作流: 分析媒体查询 -> 识别复杂度 -> 优化重构 -> 测试性能

	styleSheetId := "style-sheet-3"

	log.Println("媒体查询性能分析 - 优化复杂查询...")

	// 定义性能优化案例
	performanceCases := []struct {
		caseName   string
		original   string
		optimized  string
		complexity string
		improvement string
	}{
		{
			caseName:  "简化和合并条件",
			original:  "(min-width: 768px) and (max-width: 1024px) and (orientation: landscape) and (min-resolution: 1.5dppx)",
			optimized: "(min-width: 768px) and (max-width: 1024px) and (orientation: landscape)",
			complexity: "4个条件 -> 3个条件",
			improvement: "减少分辨率检测，提高性能",
		},
		{
			caseName:  "移除冗余条件",
			original:  "screen and (max-width: 768px) and (max-device-width: 768px)",
			optimized: "screen and (max-width: 768px)",
			complexity: "3个条件 -> 2个条件",
			improvement: "移除冗余的device-width检测",
		},
		{
			caseName:  "优化逻辑结构",
			original:  "(min-width: 320px) and (max-width: 767px) or (min-width: 1200px)",
			optimized: "(min-width: 320px) and (max-width: 767px), (min-width: 1200px)",
			complexity: "逻辑或 -> 逗号分隔",
			improvement: "使用逗号分隔替代逻辑或，提高可读性",
		},
		{
			caseName:  "简化特性检测",
			original:  "(hover: hover) and (pointer: fine) and (any-hover: hover)",
			optimized: "(hover: hover) and (pointer: fine)",
			complexity: "3个特性 -> 2个特性",
			improvement: "移除冗余的any-hover检测",
		},
	}

	for i, perfCase := range performanceCases {
		log.Printf("\n性能案例 %d/%d: %s", i+1, len(performanceCases), perfCase.caseName)
		log.Printf("  复杂度: %s", perfCase.complexity)
		log.Printf("  优化改进: %s", perfCase.improvement)

		// 这里需要获取媒体查询的范围
		mediaRange := SourceRange{
			StartLine:   10 + i*5,
			StartColumn: 1,
			EndLine:     12 + i*5,
			EndColumn:   50,
		}

		log.Printf("  应用优化: %s", perfCase.optimized)

		response, err := CDPCSSSetMediaText(styleSheetId, mediaRange, perfCase.optimized)
		if err != nil {
			log.Printf("  优化失败: %v", err)
		} else {
			log.Printf("  优化成功")

			// 分析优化结果
			var data struct {
				Result struct {
					Media *CSSMedia `json:"media"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.Media != nil {
				log.Printf("  优化后媒体查询: %s", data.Result.Media.Text)

				// 计算条件数量
				conditionCount := 0
				for _, query := range data.Result.Media.MediaList {
					conditionCount += len(query.Expressions)
				}
				log.Printf("  条件数量: %d", conditionCount)
			}
		}

		time.Sleep(300 * time.Millisecond)
	}

	log.Println("媒体查询性能分析完成")
}


*/

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
func CDPCSSSetPropertyRulePropertyName(styleSheetId string, rangeObj SourceRange, propertyName string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}
	if propertyName == "" {
		return "", fmt.Errorf("属性名不能为空")
	}
	// 验证属性名格式
	if !isValidCSSPropertyName(propertyName) {
		return "", fmt.Errorf("无效的CSS属性名: %s", propertyName)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建范围JSON
	rangeJSON, err := json.Marshal(rangeObj)
	if err != nil {
		return "", fmt.Errorf("序列化范围失败: %w", err)
	}

	// 转义特殊字符
	escapedPropertyName := strings.ReplaceAll(propertyName, `"`, `\"`)

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setPropertyRulePropertyName",
		"params": {
			"styleSheetId": "%s",
			"range": %s,
			"propertyName": "%s"
		}
	}`, reqID, styleSheetId, rangeJSON, escapedPropertyName)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// 验证CSS属性名格式
func isValidCSSPropertyName(name string) bool {
	name = strings.TrimSpace(name)

	if name == "" {
		return false
	}

	// 自定义属性通常以--开头
	if strings.HasPrefix(name, "--") {
		// 检查--后的部分
		rest := strings.TrimPrefix(name, "--")
		return isValidCustomPropertyName(rest)
	}

	// 标准属性名检查
	// 简单的格式检查
	cssPropertyPattern := regexp.MustCompile(`^[a-zA-Z_-][a-zA-Z0-9_-]*$`)
	return cssPropertyPattern.MatchString(name)
}

// 验证自定义属性名
func isValidCustomPropertyName(name string) bool {
	customPropPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)
	return customPropPattern.MatchString(name)
}

/*

示例

// 示例: CSS自定义属性重命名工具
func exampleCSSPropertyRenamingTool() {
	// === 应用场景描述 ===
	// 场景: CSS自定义属性重命名工具
	// 用途: 安全地重命名CSS @property规则
	// 优势: 批量重命名自定义属性，避免手动修改错误
	// 典型工作流: 选择@property规则 -> 重命名 -> 更新使用处 -> 验证

	styleSheetId := "style-sheet-1"

	// 定义要修改的@property规则范围
	propertyRange := SourceRange{
		StartLine:   25,
		StartColumn: 1,
		EndLine:     30,
		EndColumn:   20,
	}

	log.Println("CSS自定义属性重命名工具...")

	// 定义重命名映射
	renames := []struct {
		oldName string
		newName string
		reason  string
	}{
		{
			oldName: "--primary-color",
			newName: "--brand-primary",
			reason:  "统一设计系统命名",
		},
		{
			oldName: "--secondary-color",
			newName: "--brand-secondary",
			reason:  "统一设计系统命名",
		},
		{
			oldName: "--spacing-small",
			newName: "--space-sm",
			reason:  "使用更简洁的命名",
		},
		{
			oldName: "--animation-duration",
			newName: "--duration-normal",
			reason:  "更语义化的命名",
		},
	}

	for i, rename := range renames {
		log.Printf("重命名 %d/%d: %s", i+1, len(renames), rename.reason)
		log.Printf("  修改: %s -> %s", rename.oldName, rename.newName)

		// 注意: 这个方法的实际可用性需要验证
		response, err := CDPCSSSetPropertyRulePropertyName(styleSheetId, propertyRange, rename.newName)
		if err != nil {
			log.Printf("  重命名失败: %v", err)
			// 检查是否是方法不存在
			if strings.Contains(err.Error(), "可能不存在") {
				log.Println("  注意: 这个方法在当前Chrome版本中可能不可用")
				break
			}
		} else {
			log.Printf("  重命名成功")

			// 解析返回结果
			var data struct {
				Result struct {
					PropertyRule *CSSPropertyRule `json:"propertyRule"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.PropertyRule != nil {
				log.Printf("  返回的属性规则:")
				if data.Result.PropertyRule.PropertyName != nil {
					log.Printf("    属性名: %v", data.Result.PropertyRule.PropertyName)
				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Println("CSS属性重命名测试完成")
}

*/

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
func CDPCSSSetRuleSelector(styleSheetId string, rangeObj SourceRange, selector string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}
	if selector == "" {
		return "", fmt.Errorf("选择器不能为空")
	}
	// 验证选择器格式
	if !isValidCSSSelector(selector) {
		return "", fmt.Errorf("无效的CSS选择器: %s", selector)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建范围JSON
	rangeJSON, err := json.Marshal(rangeObj)
	if err != nil {
		return "", fmt.Errorf("序列化范围失败: %w", err)
	}

	// 转义特殊字符
	escapedSelector := strings.ReplaceAll(selector, `"`, `\"`)
	escapedSelector = strings.ReplaceAll(escapedSelector, "\n", "\\n")

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setRuleSelector",
		"params": {
			"styleSheetId": "%s",
			"range": %s,
			"selector": "%s"
		}
	}`, reqID, styleSheetId, rangeJSON, escapedSelector)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setRuleSelector 请求失败: %w", err)
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
			return "", fmt.Errorf("setRuleSelector 请求超时")
		}
	}
}

// 验证CSS选择器格式
func isValidCSSSelector(selector string) bool {
	selector = strings.TrimSpace(selector)

	if selector == "" {
		return false
	}

	// 基本格式检查
	// 不允许以数字开头
	if len(selector) > 0 && selector[0] >= '0' && selector[0] <= '9' {
		return false
	}

	// 检查常见的选择器格式
	validPatterns := []*regexp.Regexp{
		// 类选择器: .class
		regexp.MustCompile(`^\.[a-zA-Z_][a-zA-Z0-9_-]*$`),
		// ID选择器: #id
		regexp.MustCompile(`^#[a-zA-Z_][a-zA-Z0-9_-]*$`),
		// 元素选择器: element
		regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`),
		// 属性选择器: [attr]
		regexp.MustCompile(`^\[[a-zA-Z_][a-zA-Z0-9_-]*(?:[~|^$*]?=[^\]]+)?\]$`),
		// 伪类选择器: :pseudo
		regexp.MustCompile(`^:[a-zA-Z_-]+(?:\([^)]*\))?$`),
		// 伪元素选择器: ::pseudo
		regexp.MustCompile(`^::[a-zA-Z_-]+$`),
		// 组合选择器: selector1 selector2
		regexp.MustCompile(`^[a-zA-Z_.#:\[].*$`),
	}

	for _, pattern := range validPatterns {
		if pattern.MatchString(selector) {
			return true
		}
	}

	return false
}

/*

示例

// 示例1: CSS重构和重命名工具
func exampleCSSRefactoringTool() {
	// === 应用场景描述 ===
	// 场景: CSS重构和重命名工具
	// 用途: 安全地重命名CSS选择器，避免破坏现有样式
	// 优势: 批量更新选择器，确保样式一致性
	// 典型工作流: 分析现有选择器 -> 重命名 -> 验证影响 -> 更新HTML

	styleSheetId := "style-sheet-1"

	log.Println("CSS重构工具 - 重命名选择器...")

	// 定义要修改的CSS规则范围
	ruleRange := SourceRange{
		StartLine:   10,
		StartColumn: 1,
		EndLine:     12,
		EndColumn:   20,
	}

	// 定义选择器重命名映射
	selectorRenames := []struct {
		oldSelector string
		newSelector string
		description string
	}{
		{
			oldSelector: ".btn-primary",
			newSelector: ".button-primary",
			description: "统一按钮类名命名规范",
		},
		{
			oldSelector: "#header",
			newSelector: "#main-header",
			description: "更语义化的ID命名",
		},
		{
			oldSelector: ".nav-item",
			newSelector: ".navigation-item",
			description: "扩展缩写，提高可读性",
		},
		{
			oldSelector: "div.container",
			newSelector: ".container",
			description: "移除冗余的元素选择器",
		},
		{
			oldSelector: "a:hover",
			newSelector: "a:hover, a:focus",
			description: "添加焦点状态支持",
		},
		{
			oldSelector: ".card .title",
			newSelector: ".card-title",
			description: "简化嵌套选择器",
		},
	}

	for i, rename := range selectorRenames {
		log.Printf("重命名 %d/%d: %s", i+1, len(selectorRenames), rename.description)
		log.Printf("  修改: %s -> %s", rename.oldSelector, rename.newSelector)

		// 设置新的选择器
		response, err := CDPCSSSetRuleSelector(styleSheetId, ruleRange, rename.newSelector)
		if err != nil {
			log.Printf("  修改失败: %v", err)
		} else {
			log.Printf("  修改成功")

			// 解析返回的选择器列表
			var data struct {
				Result struct {
					SelectorList *SelectorList `json:"selectorList"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.SelectorList != nil {
				log.Printf("  返回的选择器: %s", data.Result.SelectorList.Text)
				log.Printf("  选择器数量: %d", len(data.Result.SelectorList.Selectors))
			}
		}

		time.Sleep(300 * time.Millisecond)
	}

	log.Println("CSS选择器重命名完成")
}

// 示例2: 设计系统迁移工具
func exampleDesignSystemMigration() {
	// === 应用场景描述 ===
	// 场景: 设计系统迁移工具
	// 用途: 从旧设计系统迁移到新设计系统的CSS选择器
	// 优势: 自动化迁移，减少手动错误
	// 典型工作流: 定义映射规则 -> 批量迁移 -> 验证结果 -> 清理旧代码

	styleSheetId := "style-sheet-2"

	log.Println("设计系统迁移 - 更新CSS选择器...")

	// 定义设计系统迁移映射
	migrationMap := []struct {
		oldSystem  string
		newSystem  string
		category   string
	}{
		// 间距系统
		{".spacing-sm", ".space-xs", "spacing"},
		{".spacing-md", ".space-sm", "spacing"},
		{".spacing-lg", ".space-md", "spacing"},
		{".spacing-xl", ".space-lg", "spacing"},

		// 颜色系统
		{".color-primary", ".primary", "color"},
		{".color-secondary", ".secondary", "color"},
		{".color-success", ".success", "color"},
		{".color-danger", ".danger", "color"},
		{".color-warning", ".warning", "color"},
		{".color-info", ".info", "color"},

		// 排版系统
		{".text-xs", ".text-sm", "typography"},
		{".text-sm", ".text-base", "typography"},
		{".text-md", ".text-lg", "typography"},
		{".text-lg", ".text-xl", "typography"},
		{".text-xl", ".text-2xl", "typography"},

		// 布局系统
		{".flex-row", ".row", "layout"},
		{".flex-col", ".col", "layout"},
		{".justify-start", ".justify-start", "layout"},
		{".justify-center", ".justify-center", "layout"},
		{".justify-end", ".justify-end", "layout"},
		{".align-start", ".items-start", "layout"},
		{".align-center", ".items-center", "layout"},
		{".align-end", ".items-end", "layout"},

		// 组件系统
		{".btn", ".button", "component"},
		{".card", ".surface", "component"},
		{".modal", ".dialog", "component"},
		{".dropdown", ".menu", "component"},
		{".badge", ".chip", "component"},
	}

	// 按类别统计
	migrationStats := make(map[string]int)

	for i, migration := range migrationMap {
		// 定义规则范围（这里简化为统一范围，实际应该根据规则位置调整）
		ruleRange := SourceRange{
			StartLine:   5 + i*2,
			StartColumn: 1,
			EndLine:     7 + i*2,
			EndColumn:   30,
		}

		log.Printf("迁移 %d/%d: [%s] %s -> %s",
			i+1, len(migrationMap), migration.category, migration.oldSystem, migration.newSystem)

		response, err := CDPCSSSetRuleSelector(styleSheetId, ruleRange, migration.newSystem)
		if err != nil {
			log.Printf("  迁移失败: %v", err)
		} else {
			log.Printf("  迁移成功")
			migrationStats[migration.category]++
		}

		// 控制迁移速度
		time.Sleep(100 * time.Millisecond)
	}

	// 生成迁移报告
	log.Println("\n=== 设计系统迁移报告 ===")
	totalMigrations := 0
	for category, count := range migrationStats {
		log.Printf("  %s: %d 个选择器", category, count)
		totalMigrations += count
	}
	log.Printf("  总计: %d/%d 个选择器迁移成功", totalMigrations, len(migrationMap))
}

// 示例3: CSS特异性优化工具
func exampleCSSSpecificityOptimizer() {
	// === 应用场景描述 ===
	// 场景: CSS特异性优化工具
	// 用途: 优化CSS选择器的特异性，提高样式性能和可维护性
	// 优势: 降低选择器复杂性，提高渲染性能
	// 典型工作流: 分析选择器特异性 -> 优化重构 -> 测试影响 -> 应用优化

	styleSheetId := "style-sheet-3"

	log.Println("CSS特异性优化 - 重构选择器...")

	// 定义优化策略
	optimizationStrategies := []struct {
		strategy  string
		examples []struct {
			before string
			after  string
			benefit string
		}
	}{
		{
			strategy: "降低ID选择器的特异性",
			examples: []struct {
				before string
				after  string
				benefit string
			}{
				{
					before: "#content .article h1",
					after:  ".content-article h1",
					benefit: "用类选择器替代ID选择器",
				},
				{
					before: "#sidebar ul.menu li a",
					after:  ".sidebar-menu-link",
					benefit: "简化深层嵌套选择器",
				},
			},
		},
		{
			strategy: "简化过度限定的选择器",
			examples: []struct {
				before string
				after  string
				benefit string
			}{
				{
					before: "div.container ul.list li.item a.link",
					after:  ".container .list-link",
					benefit: "移除冗余元素限定符",
				},
				{
					before: "body.page-home header.site-header nav.main-nav",
					after:  ".main-nav",
					benefit: "移除不必要的祖先限定",
				},
			},
		},
		{
			strategy: "合并相似选择器",
			examples: []struct {
				before string
				after  string
				benefit string
			}{
				{
					before: ".btn-primary, .btn-secondary, .btn-success",
					after:  "[class^='btn-']",
					benefit: "使用属性选择器合并",
				},
				{
					before: ".icon-home, .icon-user, .icon-settings",
					after:  "[class*='icon-']",
					benefit: "通配符合并相关类",
				},
			},
		},
		{
			strategy: "优化伪类选择器",
			examples: []struct {
				before string
				after  string
				benefit string
			}{
				{
					before: "a:link, a:visited, a:hover, a:active",
					after:  "a",
					benefit: "简化链接状态选择器",
				},
				{
					before: "input[type='text']:focus, input[type='email']:focus",
					after:  "input:focus",
					benefit: "合并相似焦点状态",
				},
			},
		},
	}

	exampleCount := 0
	for strategyIndex, strategy := range optimizationStrategies {
		log.Printf("\n优化策略: %s", strategy.strategy)

		for exampleIndex, example := range strategy.examples {
			exampleCount++
			ruleRange := SourceRange{
				StartLine:   10 + exampleCount*3,
				StartColumn: 1,
				EndLine:     12 + exampleCount*3,
				EndColumn:   50,
			}

			log.Printf("  示例 %d: %s", exampleIndex+1, example.benefit)
			log.Printf("    优化前: %s", example.before)
			log.Printf("    优化后: %s", example.after)

			// 计算特异性差异
			beforeSpec := calculateSelectorSpecificity(example.before)
			afterSpec := calculateSelectorSpecificity(example.after)

			log.Printf("    特异性: %v -> %v", beforeSpec, afterSpec)

			response, err := CDPCSSSetRuleSelector(styleSheetId, ruleRange, example.after)
			if err != nil {
				log.Printf("    优化失败: %v", err)
			} else {
				log.Printf("    优化成功")

				// 解析返回结果
				var data struct {
					Result struct {
						SelectorList *SelectorList `json:"selectorList"`
					} `json:"result"`
				}

				if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.SelectorList != nil {
					log.Printf("    实际应用的选择器: %s", data.Result.SelectorList.Text)
				}
			}

			time.Sleep(200 * time.Millisecond)
		}
	}

	log.Println("CSS特异性优化完成")
}

*/

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
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
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
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// StyleSheetTextSetResult 样式表文本设置结果
type StyleSheetTextSetResult struct {
	SourceMapURL string `json:"sourceMapURL,omitempty"` // 源映射URL
}

// ParseStyleSheetTextSet 解析样式表文本设置响应
func ParseStyleSheetTextSet(response string) (*StyleSheetTextSetResult, error) {
	var data struct {
		Result *StyleSheetTextSetResult `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("响应中没有结果")
	}

	return data.Result, nil
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
func CDPCSSSetStyleTexts(edits []StyleDeclarationEdit) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if len(edits) == 0 {
		return "", fmt.Errorf("编辑数组不能为空")
	}

	// 验证每个编辑
	for i, edit := range edits {
		if edit.StyleSheetID == "" {
			return "", fmt.Errorf("编辑[%d]的样式表ID不能为空", i)
		}
		if edit.Text == "" {
			return "", fmt.Errorf("编辑[%d]的样式文本不能为空", i)
		}
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建编辑数组JSON
	editsJSON, err := json.Marshal(edits)
	if err != nil {
		return "", fmt.Errorf("序列化编辑数组失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setStyleTexts",
		"params": {
			"edits": %s
		}
	}`, reqID, editsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// StyleDeclarationEdit 样式声明编辑
type StyleDeclarationEdit struct {
	StyleSheetID string      `json:"styleSheetId"` // 样式表ID
	Range        SourceRange `json:"range"`        // 样式文本在包含样式表中的范围
	Text         string      `json:"text"`         // 新样式文本
}

// SetStyleTextsResult 设置样式文本结果
type SetStyleTextsResult struct {
	Styles []CSSStyle `json:"styles"` // 修改后的样式数组
}

// ParseSetStyleTexts 解析设置样式文本响应
func ParseSetStyleTexts(response string) (*SetStyleTextsResult, error) {
	var data struct {
		Result *SetStyleTextsResult `json:"result"`
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

示例

// 示例1: 批量样式更新工具
func exampleBatchStyleUpdater() {
	// === 应用场景描述 ===
	// 场景: 批量样式更新工具
	// 用途: 批量更新多个样式声明，确保原子性操作
	// 优势: 减少重排重绘，提高性能，确保一致性
	// 典型工作流: 收集样式修改 -> 批量应用 -> 验证结果 -> 错误处理

	log.Println("批量样式更新工具 - 原子性批量编辑...")

	// 定义要批量编辑的样式
	styleEdits := []StyleDeclarationEdit{
		{
			StyleSheetID: "style-sheet-1",
			Range: SourceRange{
				StartLine:   10,
				StartColumn: 1,
				EndLine:     12,
				EndColumn:   30,
			},
			Text: "color: #333; font-size: 16px; line-height: 1.5;",
		},
		{
			StyleSheetID: "style-sheet-1",
			Range: SourceRange{
				StartLine:   15,
				StartColumn: 1,
				EndLine:     17,
				EndColumn:   25,
			},
			Text: "background: #f5f5f5; padding: 20px; border-radius: 8px;",
		},
		{
			StyleSheetID: "style-sheet-1",
			Range: SourceRange{
				StartLine:   20,
				StartColumn: 1,
				EndLine:     22,
				EndColumn:   20,
			},
			Text: "margin: 0 auto; max-width: 1200px;",
		},
		{
			StyleSheetID: "style-sheet-2",
			Range: SourceRange{
				StartLine:   5,
				StartColumn: 1,
				EndLine:     7,
				EndColumn:   15,
			},
			Text: "display: flex; justify-content: center; align-items: center;",
		},
		{
			StyleSheetID: "style-sheet-2",
			Range: SourceRange{
				StartLine:   25,
				StartColumn: 1,
				EndLine:     28,
				EndColumn:   40,
			},
			Text: "transition: all 0.3s ease; box-shadow: 0 2px 8px rgba(0,0,0,0.1);",
		},
	}

	log.Printf("准备批量应用 %d 个样式编辑...", len(styleEdits))

	// 批量设置样式文本
	response, err := CDPCSSSetStyleTexts(styleEdits)
	if err != nil {
		log.Printf("批量样式编辑失败: %v", err)
		return
	}

	log.Printf("批量样式编辑成功")

	// 解析返回结果
	result, err := ParseSetStyleTexts(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	log.Printf("返回的样式数量: %d", len(result.Styles))

	for i, style := range result.Styles {
		log.Printf("样式[%d]:", i+1)
		log.Printf("  样式表ID: %s", style.StyleSheetID)
		log.Printf("  属性数量: %d", len(style.CSSProperties))

		if len(style.CSSProperties) > 0 {
			log.Printf("  前3个属性:")
			for j, prop := range style.CSSProperties {
				if j < 3 {
					log.Printf("    %s: %s", prop.Name, prop.Value)
				}
			}
		}
	}
}

// 示例2: 设计系统主题切换
func exampleDesignSystemThemeSwitcher() {
	// === 应用场景描述 ===
	// 场景: 设计系统主题切换工具
	// 用途: 批量切换设计系统中的主题样式
	// 优势: 原子性主题切换，避免中间状态
	// 典型工作流: 准备主题配置 -> 批量应用 -> 切换动画 -> 持久化

	log.Println("设计系统主题切换 - 批量更新主题变量...")

	// 定义主题配置
	themes := []struct {
		name   string
		config map[string]string
	}{
		{
			name: "浅色主题",
			config: map[string]string{
				"--primary-color": "#007bff",
				"--secondary-color": "#6c757d",
				"--background-color": "#ffffff",
				"--text-color": "#333333",
				"--border-color": "#dee2e6",
				"--success-color": "#28a745",
				"--warning-color": "#ffc107",
				"--danger-color": "#dc3545",
				"--shadow-color": "rgba(0,0,0,0.1)",
			},
		},
		{
			name: "深色主题",
			config: map[string]string{
				"--primary-color": "#bb86fc",
				"--secondary-color": "#03dac6",
				"--background-color": "#121212",
				"--text-color": "#ffffff",
				"--border-color": "#2c2c2c",
				"--success-color": "#03dac6",
				"--warning-color": "#ffb74d",
				"--danger-color": "#cf6679",
				"--shadow-color": "rgba(0,0,0,0.3)",
			},
		},
		{
			name: "高对比度主题",
			config: map[string]string{
				"--primary-color": "#0056b3",
				"--secondary-color": "#495057",
				"--background-color": "#000000",
				"--text-color": "#ffffff",
				"--border-color": "#ffffff",
				"--success-color": "#155724",
				"--warning-color": "#856404",
				"--danger-color": "#721c24",
				"--shadow-color": "rgba(255,255,255,0.2)",
			},
		},
	}

	// 模拟主题切换循环
	for themeIndex, theme := range themes {
		log.Printf("切换到主题: %s", theme.name)

		// 准备样式编辑
		var styleEdits []StyleDeclarationEdit

		// 假设主题变量定义在特定位置
		// 这里创建对应的编辑操作
		for varName, varValue := range theme.config {
			styleEdits = append(styleEdits, StyleDeclarationEdit{
				StyleSheetID: "theme-variables",
				Range: SourceRange{
					StartLine:   5 + len(styleEdits),
					StartColumn: 1,
					EndLine:     6 + len(styleEdits),
					EndColumn:   50,
				},
				Text: fmt.Sprintf("%s: %s;", varName, varValue),
			})
		}

		log.Printf("  应用 %d 个主题变量...", len(styleEdits))

		// 批量应用主题变量
		if len(styleEdits) > 0 {
			response, err := CDPCSSSetStyleTexts(styleEdits)
			if err != nil {
				log.Printf("  主题切换失败: %v", err)
				continue
			}

			log.Printf("  主题切换成功")

			// 解析返回结果
			result, err := ParseSetStyleTexts(response)
			if err == nil {
				log.Printf("  更新的样式数量: %d", len(result.Styles))
			}
		}

		// 主题切换延迟
		if themeIndex < len(themes)-1 {
			log.Println("  等待2秒切换到下一个主题...")
			time.Sleep(2 * time.Second)
		}
	}

	log.Println("主题切换演示完成")
}

*/

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
		"method": "CSS.startRuleUsageTracking"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
		"method": "CSS.stopRuleUsageTracking"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
		"method": "CSS.takeCoverageDelta"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// CoverageDeltaResult 覆盖率增量结果
type CoverageDeltaResult struct {
	Timestamp     float64     `json:"timestamp"`               // 时间戳
	Coverage      []RuleUsage `json:"coverage"`                // 覆盖率数据
	IsDelta       bool        `json:"isDelta,omitempty"`       // 是否为增量数据
	SinceLastCall bool        `json:"sinceLastCall,omitempty"` // 是否自上次调用以来
}

// RuleUsage CSS规则使用情况
type RuleUsage struct {
	StyleSheetID string  `json:"styleSheetId"`           // 样式表ID
	StartOffset  int     `json:"startOffset"`            // 规则在样式表中的起始偏移
	EndOffset    int     `json:"endOffset"`              // 规则在样式表中的结束偏移
	Used         bool    `json:"used"`                   // 是否被使用
	FirstUsed    bool    `json:"firstUsed,omitempty"`    // 是否是首次使用
	LastUsedTime float64 `json:"lastUsedTime,omitempty"` // 最后使用时间
}

// ParseCoverageDelta 解析覆盖率增量响应
func ParseCoverageDelta(response string) (*CoverageDeltaResult, error) {
	var data struct {
		Result *CoverageDeltaResult `json:"result"`
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

示例

// 示例1: 实时CSS覆盖率监控工具
func exampleRealTimeCSSCoverageMonitor() {
	// === 应用场景描述 ===
	// 场景: 实时CSS覆盖率监控工具
	// 用途: 实时监控用户交互过程中的CSS规则使用变化
	// 优势: 实时反馈，精确分析CSS在用户旅程中的使用
	// 典型工作流: 开始跟踪 -> 实时获取增量 -> 分析变化 -> 生成报告

	log.Println("实时CSS覆盖率监控工具...")

	// 1. 启用CSS域
	log.Println("步骤1: 启用CSS域")
	_, err := CDPCSSEnable()
	if err != nil {
		log.Printf("启用CSS域失败: %v", err)
		return
	}
	defer CDPCSSDisable()

	// 2. 开始规则使用跟踪
	log.Println("步骤2: 开始CSS规则使用跟踪")
	_, err = CDPCSSStartRuleUsageTracking()
	if err != nil {
		log.Printf("开始跟踪失败: %v", err)
		return
	}

	// 3. 实时监控循环
	log.Println("步骤3: 开始实时监控循环")
	log.Println("按Ctrl+C停止监控...")

	monitorDuration := 30 * time.Second
	interval := 2 * time.Second
	startTime := time.Now()
	iteration := 0

	var allCoverageData []*CoverageDeltaResult

	for time.Since(startTime) < monitorDuration {
		iteration++
		currentTime := time.Now()

		log.Printf("\n监控迭代 %d (已运行: %v)", iteration, time.Since(startTime))

		// 获取覆盖率增量数据
		response, err := CDPCSSTakeCoverageDelta()
		if err != nil {
			log.Printf("获取覆盖率增量失败: %v", err)
			break
		}

		// 解析数据
		result, err := ParseCoverageDelta(response)
		if err != nil {
			log.Printf("解析覆盖率数据失败: %v", err)
			break
		}

		allCoverageData = append(allCoverageData, result)

		// 分析本次增量
		analyzeCoverageDelta(result, iteration, currentTime)

		// 等待下一个间隔
		time.Sleep(interval)
	}

	// 4. 停止跟踪
	log.Println("步骤4: 停止CSS规则使用跟踪")
	stopResponse, err := CDPCSSStopRuleUsageTracking()
	if err != nil {
		log.Printf("停止跟踪失败: %v", err)
	} else {
		log.Printf("停止跟踪响应: %s", stopResponse)
	}

	// 5. 生成综合报告
	log.Println("步骤5: 生成实时监控综合报告")
	generateRealTimeMonitorReport(allCoverageData, monitorDuration, interval)
}

// 分析覆盖率增量
func analyzeCoverageDelta(result *CoverageDeltaResult, iteration int, timestamp time.Time) {
	totalRules := len(result.Coverage)
	newlyUsed := 0
	stillUnused := 0

	for _, rule := range result.Coverage {
		if rule.Used {
			if rule.FirstUsed {
				newlyUsed++
			}
		} else {
			stillUnused++
		}
	}

	log.Printf("  时间: %s", timestamp.Format("15:04:05"))
	log.Printf("  总规则数: %d", totalRules)
	log.Printf("  新增使用规则: %d", newlyUsed)

	if totalRules > 0 {
		usedCount := totalRules - stillUnused
		coverageRate := float64(usedCount) / float64(totalRules) * 100
		log.Printf("  当前覆盖率: %.1f%% (%d/%d)", coverageRate, usedCount, totalRules)

		// 实时建议
		if newlyUsed > 0 {
			log.Printf("  📈 发现 %d 个新使用的CSS规则", newlyUsed)
		}

		if coverageRate < 30 && iteration > 3 {
			log.Printf("  ⚠ 覆盖率较低，考虑优化CSS加载策略")
		}
	}
}

// 生成实时监控报告
func generateRealTimeMonitorReport(data []*CoverageDeltaResult, duration time.Duration, interval time.Duration) {
	log.Println("\n=== 实时CSS覆盖率监控报告 ===")
	log.Printf("监控总时长: %v", duration)
	log.Printf("监控间隔: %v", interval)
	log.Printf("数据采集次数: %d", len(data))

	if len(data) == 0 {
		log.Println("没有收集到有效数据")
		return
	}

	// 分析覆盖率趋势
	firstSnapshot := data[0]
	lastSnapshot := data[len(data)-1]

	firstUsed := 0
	lastUsed := 0

	for _, rule := range firstSnapshot.Coverage {
		if rule.Used {
			firstUsed++
		}
	}

	for _, rule := range lastSnapshot.Coverage {
		if rule.Used {
			lastUsed++
		}
	}

	totalRules := len(firstSnapshot.Coverage)
	if totalRules > 0 {
		startCoverage := float64(firstUsed) / float64(totalRules) * 100
		endCoverage := float64(lastUsed) / float64(totalRules) * 100
		coverageIncrease := endCoverage - startCoverage

		log.Printf("\n覆盖率趋势分析:")
		log.Printf("  开始覆盖率: %.1f%%", startCoverage)
		log.Printf("  结束覆盖率: %.1f%%", endCoverage)
		log.Printf("  覆盖率增长: +%.1f%%", coverageIncrease)

		// 分析覆盖率增长模式
		log.Printf("\n覆盖率增长模式:")
		if coverageIncrease > 20 {
			log.Println("  📈 快速增长: 用户进行了大量交互")
		} else if coverageIncrease > 5 {
			log.Println("  📈 稳定增长: 正常的用户交互")
		} else {
			log.Println("  📉 缓慢增长: 有限的用户交互")
		}
	}

	// 识别关键发现
	log.Println("\n关键发现:")

	// 分析哪些规则是逐步使用的
	if len(data) >= 3 {
		earlyUsed := make(map[string]bool)  // 早期使用的规则
		lateUsed := make(map[string]bool)   // 后期使用的规则

		// 检查前1/3时间段
		earlyThreshold := len(data) / 3
		for i := 0; i < earlyThreshold && i < len(data); i++ {
			for _, rule := range data[i].Coverage {
				if rule.Used {
					key := fmt.Sprintf("%s:%d-%d", rule.StyleSheetID, rule.StartOffset, rule.EndOffset)
					earlyUsed[key] = true
				}
			}
		}

		// 检查最后1/3时间段
		lateStart := len(data) * 2 / 3
		for i := lateStart; i < len(data); i++ {
			for _, rule := range data[i].Coverage {
				if rule.Used {
					key := fmt.Sprintf("%s:%d-%d", rule.StyleSheetID, rule.StartOffset, rule.EndOffset)
					if !earlyUsed[key] {
						lateUsed[key] = true
					}
				}
			}
		}

		log.Printf("  早期使用规则: %d 个", len(earlyUsed))
		log.Printf("  后期使用规则: %d 个", len(lateUsed))

		if len(lateUsed) > len(earlyUsed) {
			log.Println("  ⚠ 发现: 大量CSS在后期才被使用，考虑延迟加载")
		}
	}

	// 优化建议
	log.Println("\n优化建议:")
	log.Println("  1. 基于使用时机优化CSS加载策略")
	log.Println("  2. 对后期使用的CSS实现懒加载")
	log.Println("  3. 将关键CSS内联到HTML中")
	log.Println("  4. 移除始终未使用的CSS规则")
	log.Println("  5. 监控不同用户场景的CSS使用模式")
}


*/

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
		"method": "CSS.getEnvironmentVariables"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// EnvironmentVariable CSS环境变量
type EnvironmentVariable struct {
	Name  string `json:"name"`  // 环境变量名称
	Value string `json:"value"` // 环境变量值
	Type  string `json:"type"`  // 变量类型："custom" 或 "builtin"
}

// EnvironmentVariablesResult 环境变量结果
type EnvironmentVariablesResult struct {
	Variables []EnvironmentVariable `json:"variables"` // 环境变量数组
}

// ParseEnvironmentVariables 解析环境变量响应
func ParseEnvironmentVariables(response string) (*EnvironmentVariablesResult, error) {
	var data struct {
		Result *EnvironmentVariablesResult `json:"result"`
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

示例

// 示例1: CSS环境变量分析工具
func exampleCSSEnvironmentVariableAnalyzer() {
	// === 应用场景描述 ===
	// 场景: CSS环境变量分析工具
	// 用途: 分析页面中定义的所有CSS环境变量
	// 优势: 全面了解CSS环境变量的使用情况
	// 典型工作流: 获取环境变量 -> 分类分析 -> 识别问题 -> 优化建议

	log.Println("CSS环境变量分析工具...")

	// 1. 启用CSS域
	log.Println("步骤1: 启用CSS域")
	_, err := CDPCSSEnable()
	if err != nil {
		log.Printf("启用CSS域失败: %v", err)
		return
	}
	defer CDPCSSDisable()

	// 2. 获取环境变量
	log.Println("步骤2: 获取CSS环境变量")
	response, err := CDPCSSGetEnvironmentVariables()
	if err != nil {
		log.Printf("获取环境变量失败: %v", err)
		return
	}

	// 3. 解析结果
	log.Println("步骤3: 解析环境变量数据")
	result, err := ParseEnvironmentVariables(response)
	if err != nil {
		log.Printf("解析结果失败: %v", err)
		return
	}

	// 4. 分析环境变量
	log.Println("步骤4: 分析环境变量")
	generateEnvironmentVariableReport(result)
}

// 生成环境变量报告
func generateEnvironmentVariableReport(result *EnvironmentVariablesResult) {
	log.Println("\n=== CSS环境变量分析报告 ===")
	log.Printf("总环境变量数量: %d", len(result.Variables))

	if len(result.Variables) == 0 {
		log.Println("未发现CSS环境变量")
		return
	}

	// 按类型统计
	customCount := 0
	builtinCount := 0
	customVars := make([]EnvironmentVariable, 0)
	builtinVars := make([]EnvironmentVariable, 0)

	for _, variable := range result.Variables {
		if variable.Type == "custom" {
			customCount++
			customVars = append(customVars, variable)
		} else if variable.Type == "builtin" {
			builtinCount++
			builtinVars = append(builtinVars, variable)
		}
	}

	log.Printf("自定义变量: %d 个", customCount)
	log.Printf("内置变量: %d 个", builtinCount)

	// 分析自定义变量
	if customCount > 0 {
		log.Println("\n=== 自定义环境变量分析 ===")

		// 按前缀分组
		prefixGroups := make(map[string]int)
		valueLengthStats := make(map[string]int) // 值长度统计

		for _, variable := range customVars {
			// 提取变量前缀（最后一个-之前的部分）
			parts := strings.Split(variable.Name, "-")
			if len(parts) > 1 {
				prefix := parts[0]
				prefixGroups[prefix]++
			} else {
				prefixGroups["other"]++
			}

			// 统计值长度
			length := len(variable.Value)
			rangeKey := ""
			switch {
			case length <= 5:
				rangeKey = "0-5"
			case length <= 10:
				rangeKey = "6-10"
			case length <= 20:
				rangeKey = "11-20"
			case length <= 50:
				rangeKey = "21-50"
			default:
				rangeKey = "50+"
			}
			valueLengthStats[rangeKey]++
		}

		log.Println("变量前缀分组:")
		for prefix, count := range prefixGroups {
			log.Printf("  %s: %d 个", prefix, count)
		}

		log.Println("\n变量值长度分布:")
		for rangeKey, count := range valueLengthStats {
			log.Printf("  %s字符: %d 个", rangeKey, count)
		}

		// 显示自定义变量详情
		log.Println("\n自定义变量详情 (前10个):")
		for i, variable := range customVars {
			if i < 10 {
				log.Printf("  [%d] %s = %s", i+1, variable.Name, variable.Value)
			}
		}
		if len(customVars) > 10 {
			log.Printf("  ... 还有 %d 个自定义变量", len(customVars)-10)
		}
	}

	// 分析内置变量
	if builtinCount > 0 {
		log.Println("\n=== 内置环境变量分析 ===")

		// 按常见类型分组
		builtinCategories := make(map[string]int)

		for _, variable := range builtinVars {
			category := categorizeBuiltinVariable(variable.Name)
			builtinCategories[category]++
		}

		log.Println("内置变量分类:")
		for category, count := range builtinCategories {
			log.Printf("  %s: %d 个", category, count)
		}

		// 显示内置变量详情
		log.Println("\n内置变量详情 (前5个):")
		for i, variable := range builtinVars {
			if i < 5 {
				log.Printf("  [%d] %s = %s", i+1, variable.Name, variable.Value)
			}
		}
		if len(builtinVars) > 5 {
			log.Printf("  ... 还有 %d 个内置变量", len(builtinVars)-5)
		}
	}

	// 优化建议
	log.Println("\n=== 优化建议 ===")
	if customCount > 20 {
		log.Println("  ⚠ 自定义变量较多，建议:")
		log.Println("    - 考虑变量命名规范化")
		log.Println("    - 合并相似功能的变量")
		log.Println("    - 建立变量命名约定")
	}

	if customCount == 0 {
		log.Println("  ⚠ 未使用自定义环境变量，建议:")
		log.Println("    - 考虑使用CSS环境变量提高样式可维护性")
		log.Println("    - 为设计系统定义基础变量")
		log.Println("    - 实现主题切换功能")
	}
}

// 分类内置变量
func categorizeBuiltinVariable(name string) string {
	// 常见的CSS环境变量分类
	switch {
	case strings.Contains(name, "safe-area-inset"):
		return "安全区域"
	case strings.Contains(name, "viewport"):
		return "视口相关"
	case strings.Contains(name, "color"):
		return "颜色相关"
	case strings.Contains(name, "size") || strings.Contains(name, "width") || strings.Contains(name, "height"):
		return "尺寸相关"
	case strings.Contains(name, "spacing") || strings.Contains(name, "margin") || strings.Contains(name, "padding"):
		return "间距相关"
	case strings.Contains(name, "font") || strings.Contains(name, "text"):
		return "文字相关"
	default:
		return "其他"
	}
}


*/

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
func CDPCSSSetContainerQueryText(styleSheetId string, rangeObj SourceRange, text string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if styleSheetId == "" {
		return "", fmt.Errorf("样式表ID不能为空")
	}
	if text == "" {
		return "", fmt.Errorf("容器查询文本不能为空")
	}
	// 验证容器查询文本格式
	if !isValidContainerQuery(text) {
		return "", fmt.Errorf("无效的容器查询文本: %s", text)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建范围JSON
	rangeJSON, err := json.Marshal(rangeObj)
	if err != nil {
		return "", fmt.Errorf("序列化范围失败: %w", err)
	}

	// 转义特殊字符
	escapedText := strings.ReplaceAll(text, `"`, `\"`)
	escapedText = strings.ReplaceAll(escapedText, "\n", "\\n")

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CSS.setContainerQueryText",
		"params": {
			"styleSheetId": "%s",
			"range": %s,
			"text": "%s"
		}
	}`, reqID, styleSheetId, rangeJSON, escapedText)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

// 验证容器查询格式
func isValidContainerQuery(text string) bool {
	text = strings.TrimSpace(text)

	if text == "" {
		return false
	}

	// 基本格式检查
	// 检查是否包含@container
	if !strings.Contains(text, "@container") {
		return false
	}

	// 检查是否包含大括号
	if !strings.Contains(text, "{") {
		return false
	}

	// 常见的容器查询条件
	containerFeatures := []string{
		"width", "height", "inline-size", "block-size",
		"aspect-ratio", "orientation", "style", "style-type",
		"scroll-state", "anchored", "container-type", "container-name",
	}

	// 检查是否包含有效的容器特性
	hasValidFeature := false
	for _, feature := range containerFeatures {
		if strings.Contains(text, feature) {
			hasValidFeature = true
			break
		}
	}

	return hasValidFeature
}

// SetContainerQueryTextResult 设置容器查询文本结果
type SetContainerQueryTextResult struct {
	ContainerQuery CSSContainerQuery `json:"containerQuery"` // 修改后的容器查询规则
}

// ParseSetContainerQueryText 解析设置容器查询文本响应
func ParseSetContainerQueryText(response string) (*SetContainerQueryTextResult, error) {
	var data struct {
		Result *SetContainerQueryTextResult `json:"result"`
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

示例

// 示例1: 容器查询实时编辑器
func exampleContainerQueryEditor() {
	// === 应用场景描述 ===
	// 场景: 容器查询实时编辑器
	// 用途: 实时编辑和预览CSS容器查询效果
	// 优势: 所见即所得，快速迭代容器查询设计
	// 典型工作流: 选择容器查询 -> 编辑条件 -> 实时预览 -> 优化调整

	styleSheetId := "style-sheet-1"

	log.Println("容器查询实时编辑器...")

	// 定义要修改的容器查询范围
	containerQueryRange := SourceRange{
		StartLine:   12,
		StartColumn: 1,
		EndLine:     15,
		EndColumn:   30,
	}

	// 定义容器查询编辑测试
	containerQueryTests := []struct {
		oldQuery    string
		newQuery    string
		description string
		useCase     string
	}{
		{
			oldQuery: "@container (min-width: 300px) { ... }",
			newQuery: "@container (min-width: 400px) { ... }",
			description: "增加容器断点宽度",
			useCase: "调整卡片组件的响应式断点",
		},
		{
			oldQuery: "@container card (min-width: 400px) { ... }",
			newQuery: "@container card (min-width: 400px) and (orientation: landscape) { ... }",
			description: "添加方向条件",
			useCase: "优化卡片在横屏设备上的表现",
		},
		{
			oldQuery: "@container (width > 300px) { ... }",
			newQuery: "@container (inline-size > 300px) { ... }",
			description: "使用逻辑属性",
			useCase: "改进RTL（从右到左）语言支持",
		},
		{
			oldQuery: "@container (min-height: 200px) { ... }",
			newQuery: "@container (min-height: 200px) and style(--theme: dark) { ... }",
			description: "添加样式查询",
			useCase: "支持暗色主题下的容器样式",
		},
		{
			oldQuery: "@container sidebar (max-width: 250px) { ... }",
			newQuery: "@container sidebar (max-width: 250px) and scroll-state(active) { ... }",
			description: "添加滚动状态查询",
			useCase: "优化滚动时的侧边栏表现",
		},
		{
			oldQuery: "@container (aspect-ratio > 1/1) { ... }",
			newQuery: "@container (aspect-ratio > 1/1) and anchored { ... }",
			description: "添加锚定查询",
			useCase: "处理固定位置容器的响应式",
		},
	}

	for i, test := range containerQueryTests {
		log.Printf("测试 %d/%d: %s", i+1, len(containerQueryTests), test.description)
		log.Printf("  使用场景: %s", test.useCase)
		log.Printf("  修改: %s -> %s", test.oldQuery, test.newQuery)

		// 设置新的容器查询文本
		response, err := CDPCSSSetContainerQueryText(styleSheetId, containerQueryRange, test.newQuery)
		if err != nil {
			log.Printf("  修改失败: %v", err)
		} else {
			log.Printf("  修改成功")

			// 解析返回的容器查询规则
			var data struct {
				Result struct {
					ContainerQuery *CSSContainerQuery `json:"containerQuery"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.ContainerQuery != nil {
				query := data.Result.ContainerQuery
				log.Printf("  返回的容器查询:")
				log.Printf("    文本: %s", query.Text)
				if query.Name != "" {
					log.Printf("    名称: %s", query.Name)
				}
				if query.PhysicalAxes != "" {
					log.Printf("    物理轴: %s", query.PhysicalAxes)
				}
				if query.LogicalAxes != "" {
					log.Printf("    逻辑轴: %s", query.LogicalAxes)
				}
				if query.QueriesScrollState {
					log.Printf("    包含滚动状态查询: 是")
				}
				if query.QueriesAnchored {
					log.Printf("    包含锚定查询: 是")
				}
			}
		}

		// 等待一段时间查看效果
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("容器查询编辑测试完成")
}

// 示例2: 组件级响应式设计工具
func exampleComponentResponsiveDesignTool() {
	// === 应用场景描述 ===
	// 场景: 组件级响应式设计工具
	// 用途: 为不同组件设计和调整容器查询
	// 优势: 组件级别的细粒度响应式控制
	// 典型工作流: 选择组件 -> 设计容器查询 -> 测试断点 -> 优化性能

	styleSheetId := "component-styles"

	log.Println("组件级响应式设计工具...")

	// 定义常见组件的容器查询配置
	componentConfigs := []struct {
		component   string
		container   string
		queries     []struct {
			name        string
			description string
			query       string
		}
	}{
		{
			component: "卡片组件 (Card)",
			container: "card-container",
			queries: []struct {
				name        string
				description string
				query       string
			}{
				{
					name:        "紧凑模式",
					description: "小空间显示",
					query:       "@container card-container (max-width: 300px) { .card { padding: 0.5rem; font-size: 0.875rem; } }",
				},
				{
					name:        "标准模式",
					description: "中等空间显示",
					query:       "@container card-container (min-width: 300px) and (max-width: 600px) { .card { padding: 1rem; font-size: 1rem; } }",
				},
				{
					name:        "展开模式",
					description: "大空间显示",
					query:       "@container card-container (min-width: 600px) { .card { padding: 1.5rem; font-size: 1.125rem; } }",
				},
			},
		},
		{
			component: "导航栏 (Navbar)",
			container: "navbar-container",
			queries: []struct {
				name        string
				description string
				query       string
			}{
				{
					name:        "移动端",
					description: "小屏幕导航",
					query:       "@container navbar-container (max-width: 768px) { .navbar { flex-direction: column; } }",
				},
				{
					name:        "桌面端",
					description: "大屏幕导航",
					query:       "@container navbar-container (min-width: 768px) { .navbar { flex-direction: row; } }",
				},
			},
		},
		{
			component: "模态框 (Modal)",
			container: "modal-container",
			queries: []struct {
				name        string
				description string
				query       string
			}{
				{
					name:        "小弹窗",
					description: "小屏幕弹窗",
					query:       "@container modal-container (max-width: 400px) { .modal { width: 90vw; } }",
				},
				{
					name:        "中等弹窗",
					description: "中等屏幕弹窗",
					query:       "@container modal-container (min-width: 400px) and (max-width: 800px) { .modal { width: 70vw; } }",
				},
				{
					name:        "大弹窗",
					description: "大屏幕弹窗",
					query:       "@container modal-container (min-width: 800px) { .modal { width: 50vw; } }",
				},
			},
		},
		{
			component: "表单 (Form)",
			container: "form-container",
			queries: []struct {
				name        string
				description string
				query       string
			}{
				{
					name:        "单列布局",
					description: "小屏幕表单",
					query:       "@container form-container (max-width: 500px) { .form-group { flex-direction: column; } }",
				},
				{
					name:        "双列布局",
					description: "中等屏幕表单",
					query:       "@container form-container (min-width: 500px) and (max-width: 900px) { .form-group { flex-direction: row; flex-wrap: wrap; } }",
				},
				{
					name:        "网格布局",
					description: "大屏幕表单",
					query:       "@container form-container (min-width: 900px) { .form { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; } }",
				},
			},
		},
	}

	for componentIndex, component := range componentConfigs {
		log.Printf("\n=== 配置组件: %s ===", component.component)
		log.Printf("容器名称: %s", component.container)

		// 为每个查询创建范围
		baseLine := 20 + componentIndex*20

		for queryIndex, query := range component.queries {
			lineOffset := baseLine + queryIndex*4
			queryRange := SourceRange{
				StartLine:   lineOffset,
				StartColumn: 1,
				EndLine:     lineOffset + 3,
				EndColumn:   80,
			}

			log.Printf("\n查询 %d/%d: %s", queryIndex+1, len(component.queries), query.name)
			log.Printf("描述: %s", query.description)
			log.Printf("应用查询: %s", query.query)

			// 设置容器查询
			response, err := CDPCSSSetContainerQueryText(styleSheetId, queryRange, query.query)
			if err != nil {
				log.Printf("  配置失败: %v", err)
				continue
			}

			log.Printf("  配置成功")

			// 解析返回结果
			var data struct {
				Result struct {
					ContainerQuery *CSSContainerQuery `json:"containerQuery"`
				} `json:"result"`
			}

			if err := json.Unmarshal([]byte(response), &data); err == nil && data.Result.ContainerQuery != nil {
				cq := data.Result.ContainerQuery
				log.Printf("  容器查询详情:")
				log.Printf("    容器名称: %s", cq.Name)
				log.Printf("    查询文本: %s", cq.Text)

				// 验证查询特性
				if cq.QueriesScrollState {
					log.Printf("    ⚡ 包含滚动状态查询")
				}
				if cq.QueriesAnchored {
					log.Printf("    ⚡ 包含锚定查询")
				}
			}

			// 短暂延迟，模拟用户查看效果
			time.Sleep(300 * time.Millisecond)
		}
	}

	log.Println("\n组件级响应式设计配置完成")
}


*/
