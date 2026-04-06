package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Overlay.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 关闭调试高亮: 关闭DOM/图层/网格调试的可视化覆盖层
// 2. 页面恢复原始状态: 截图/录屏前清除调试标记，保证画面纯净
// 3. 测试环境清理: 自动化调试完成后清除所有覆盖层标记
// 4. 资源释放: 关闭Overlay模块，释放浏览器渲染覆盖层资源
// 5. 多任务切换: 切换调试任务时关闭上一个任务的覆盖层
// 6. 错误恢复: 覆盖层显示异常时强制关闭恢复正常

// CDPOverlayDisable 关闭覆盖层调试功能
func CDPOverlayDisable() (string, error) {
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
		"method": "Overlay.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.disable 请求超时")
		}
	}
}

/*

// ==================== Overlay.disable 使用示例 ====================
func ExampleCDPOverlayDisable() {
	// ========== 示例1：基础关闭覆盖层 ==========
	resp, err := CDPOverlayDisable()
	if err != nil {
		log.Fatalf("关闭覆盖层失败: %v", err)
	}
	log.Printf("关闭覆盖层成功，页面恢复原始状态: %s", resp)

	// ========== 示例2：调试完成后标准清理流程 ==========
	/*
	// 调试步骤：启用 → 高亮/标记 → 关闭
	CDPOverlayEnable()
	// 执行DOM高亮、图层调试、网格显示等操作...

	// 调试完成，关闭覆盖层（defer确保一定执行）
	defer func() {
		_, err := CDPOverlayDisable()
		if err != nil {
			log.Printf("关闭覆盖层失败: %v", err)
		} else {
			log.Println("覆盖层已关闭，页面已恢复正常")
		}
	}()

	// ========== 示例3：截图前清除调试标记 ==========
	// 页面截图前必须关闭覆盖层，避免截图带调试框
	func capturePageScreenshot() {
		// 先关闭所有覆盖层
		CDPOverlayDisable()

		// 执行截图逻辑...
		log.Println("覆盖层已清除，开始纯净截图")
	}

	// ========== 示例4：自动化测试用例收尾 ==========
	func TestElementDebug(t *testing.T) {
		// 测试前启用
		CDPOverlayEnable()
		// 测试后自动关闭
		defer CDPOverlayDisable()

		// 执行元素调试、高亮操作...
		t.Log("测试完成，自动清理覆盖层")
	}
}

*/

// -----------------------------------------------  Overlay.enable  -----------------------------------------------
// === 应用场景 ===
// 1. DOM调试：开启覆盖层，对页面元素进行框选、高亮、尺寸标注
// 2. 图层渲染调试：显示页面图层边界、合成信息，辅助性能分析
// 3. 自动化可视化调试：测试过程中高亮关键元素，便于排查问题
// 4. 页面布局调试：显示网格、辅助线，检查页面排版是否正常
// 5. 截图标注：截图前对目标元素进行标记高亮，增强调试效果
// 6. 前端教学演示：实时标注页面元素，直观展示DOM结构与渲染逻辑

// CDPOverlayEnable 启用页面调试覆盖层功能
func CDPOverlayEnable() (string, error) {
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
		"method": "Overlay.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.enable 请求超时")
		}
	}
}

/*

// ==================== Overlay.enable 使用示例 ====================
func ExampleCDPOverlayEnable() {
	// ========== 示例1：基础启用覆盖层 ==========
	resp, err := CDPOverlayEnable()
	if err != nil {
		log.Fatalf("启用覆盖层功能失败: %v", err)
	}
	log.Printf("启用覆盖层成功，可开始可视化调试: %s", resp)

	// ========== 示例2：标准可视化调试完整流程（推荐） ==========
	// 1. 启用覆盖层（必须第一步）
	_, err := CDPOverlayEnable()
	if err != nil {
		log.Fatalf("启用Overlay失败: %v", err)
	}
	log.Println("已启用覆盖层，可进行元素高亮、图层调试")

	// 2. 执行调试操作（高亮元素、显示网格等）
	// 例：CDPOverlayHighlightNode(节点ID)

	// 3. 调试完成后关闭，释放资源（defer确保一定执行）
	defer CDPOverlayDisable()

	// ========== 示例3：自动化元素定位调试 ==========
	// 测试查找并高亮页面按钮元素
	func TestButtonElement(t *testing.T) {
		// 前置：开启覆盖层
		_, err := CDPOverlayEnable()
		if err != nil {
			t.Fatal(err)
		}
		// 结束自动关闭
		defer CDPOverlayDisable()

		// 执行元素查找、高亮操作...
		t.Log("已启用覆盖层，正在高亮目标元素")
	}

	// ========== 示例4：配合LayerTree进行渲染调试 ==========
	// 同时启用图层树 + 覆盖层，进行深度渲染分析
	CDPLayerTreeEnable()
	CDPOverlayEnable()
	defer CDPLayerTreeDisable()
	defer CDPOverlayDisable()

	log.Println("已开启图层+覆盖层调试，可查看图层边界与结构")
}

*/

// -----------------------------------------------  Overlay.getGridHighlightObjectsForTest  -----------------------------------------------
// === 应用场景 ===
//1. CSS Grid布局自动化测试：获取网格布局数据，验证布局是否符合预期
//2. 网格渲染调试：获取Grid高亮信息，排查布局错位、渲染异常
//3. 前端布局审计：自动采集页面网格结构，生成布局合规报告
//4. 可视化渲染验证：对比网格高亮数据，确认渲染与设计一致
//5. 自动化截图辅助：获取网格边界数据，精准定位Grid元素截图
//6. 多端布局兼容测试：获取网格数据，校验不同设备下Grid渲染效果

// CDPOverlayGetGridHighlightObjectsForTest 获取页面CSS Grid高亮对象数据
func CDPOverlayGetGridHighlightObjectsForTest() (string, error) {
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
		"method": "Overlay.getGridHighlightObjectsForTest"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.getGridHighlightObjectsForTest 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
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
			return "", fmt.Errorf("Overlay.getGridHighlightObjectsForTest 请求超时")
		}
	}
}

/*

// ==================== Overlay.getGridHighlightObjectsForTest 使用示例 ====================
func ExampleCDPOverlayGetGridHighlightObjectsForTest() {
	// ========== 示例1：基础获取网格高亮数据 ==========
	resp, err := CDPOverlayGetGridHighlightObjectsForTest()
	if err != nil {
		log.Fatalf("获取Grid高亮数据失败: %v", err)
	}
	log.Printf("获取网格布局高亮对象成功: %s", resp)

	// ========== 示例2：完整网格布局调试流程 ==========
	// 1. 启用覆盖层（必须前置）
	_, err := CDPOverlayEnable()
	if err != nil {
		log.Fatalf("启用Overlay失败: %v", err)
	}
	defer CDPOverlayDisable()

	// 2. 获取页面所有CSS Grid布局数据
	gridData, _ := CDPOverlayGetGridHighlightObjectsForTest()
	log.Println("页面网格布局高亮数据：", gridData)

	// 3. 解析数据验证网格行列数、尺寸、位置
	// 可用于自动化断言：网格数量、行列布局是否符合预期

	// ========== 示例3：自动化Grid布局测试用例 ==========
	func TestGridLayout(t *testing.T) {
		// 前置启用覆盖层
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 获取网格数据
		gridHighlight, err := CDPOverlayGetGridHighlightObjectsForTest()
		if err != nil {
			t.Fatalf("获取网格数据失败: %v", err)
		}

		// 断言：页面必须包含至少1个Grid布局
		if len(gridHighlight) == 0 {
			t.Error("未检测到CSS Grid布局")
		} else {
			t.Log("Grid布局获取成功，布局验证通过")
		}
	}
}

*/

// -----------------------------------------------  Overlay.getHighlightObjectForTest  -----------------------------------------------
// === 应用场景 ===
// 1. 元素可视化测试：获取节点高亮数据，自动校验元素渲染是否正常
// 2. 自动化DOM校验：验证节点位置、尺寸、样式是否符合预期
// 3. 调试数据采集：采集元素高亮配置，用于离线渲染分析
// 4. 页面截图辅助：获取元素高亮边界，实现精准元素截图
// 5. 前端自动化验收：对比高亮数据，校验UI还原度
// 6. 异常元素排查：定位渲染错位、隐藏、异常显示的DOM节点

// CDPOverlayGetHighlightObjectForTest 获取指定节点的高亮调试对象
// nodeId: DOM节点ID（从DOM.getDocument等方法获取）
func CDPOverlayGetHighlightObjectForTest(nodeId int) (string, error) {
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
		"method": "Overlay.getHighlightObjectForTest",
		"params": {
			"nodeId": %d
		}
	}`, reqID, nodeId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.getHighlightObjectForTest 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
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
			return "", fmt.Errorf("Overlay.getHighlightObjectForTest 请求超时")
		}
	}
}

/*

// ==================== Overlay.getHighlightObjectForTest 使用示例 ====================
func ExampleCDPOverlayGetHighlightObjectForTest() {
	// 前提：已获取有效DOM节点ID（nodeId）
	nodeID := 123

	// ========== 示例1：基础获取节点高亮数据 ==========
	resp, err := CDPOverlayGetHighlightObjectForTest(nodeID)
	if err != nil {
		log.Fatalf("获取节点高亮对象失败: %v", err)
	}
	log.Printf("获取元素高亮配置成功: %s", resp)

	// ========== 示例2：完整元素调试流程 ==========
	// 1. 启用覆盖层（必须前置）
	_, err := CDPOverlayEnable()
	if err != nil {
		log.Fatalf("启用Overlay失败: %v", err)
	}
	defer CDPOverlayDisable()

	// 2. 获取目标节点高亮信息（尺寸、位置、颜色、边界）
	nodeId := 128
	highlightData, _ := CDPOverlayGetHighlightObjectForTest(nodeId)
	log.Println("元素高亮信息：", highlightData)

	// 3. 解析数据用于自动化校验

	// ========== 示例3：自动化元素渲染测试用例 ==========
	func TestElementRender(t *testing.T) {
		// 前置启用覆盖层
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 获取按钮节点高亮数据
		nodeId := 130
		highlight, err := CDPOverlayGetHighlightObjectForTest(nodeId)
		if err != nil {
			t.Fatalf("获取高亮数据失败: %v", err)
		}

		// 断言：元素必须正常渲染（存在高亮数据）
		if highlight == "" {
			t.Error("元素渲染异常，未获取到高亮信息")
		} else {
			t.Log("元素渲染正常，高亮数据获取成功")
		}
	}
}

*/

// -----------------------------------------------  Overlay.getSourceOrderHighlightObjectForTest  -----------------------------------------------
// === 应用场景 ===
// 1. Flex/Grid顺序调试：查看元素源码顺序 vs 渲染顺序是否一致
// 2. 布局错乱排查：定位order属性导致的元素显示错位问题
// 3. 自动化顺序校验：验证页面元素渲染顺序符合设计规范
// 4. 可视化教学演示：直观展示DOM顺序与视觉顺序差异
// 5. 响应式布局检查：不同屏幕下order变化的渲染验证
// 6. 前端自动化验收：校验元素顺序是否符合需求文档

// CDPOverlayGetSourceOrderHighlightObjectForTest 获取节点源码顺序高亮对象
// nodeId: DOM节点ID
func CDPOverlayGetSourceOrderHighlightObjectForTest(nodeId int) (string, error) {
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
		"method": "Overlay.getSourceOrderHighlightObjectForTest",
		"params": {
			"nodeId": %d
		}
	}`, reqID, nodeId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.getSourceOrderHighlightObjectForTest 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
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
			return "", fmt.Errorf("Overlay.getSourceOrderHighlightObjectForTest 请求超时")
		}
	}
}

/*

// ==================== Overlay.getSourceOrderHighlightObjectForTest 使用示例 ====================
func ExampleCDPOverlayGetSourceOrderHighlightObjectForTest() {
	// 前提：已获取DOM节点ID
	nodeID := 125

	// ========== 示例1：获取节点源码顺序高亮数据 ==========
	resp, err := CDPOverlayGetSourceOrderHighlightObjectForTest(nodeID)
	if err != nil {
		log.Fatalf("获取源码顺序高亮对象失败: %v", err)
	}
	log.Printf("获取成功: %s", resp)

	// ========== 示例2：Flex/Grid顺序调试完整流程 ==========
	// 1. 启用覆盖层（必须）
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 获取目标节点的源码顺序高亮信息
	nodeId := 128
	highlight, _ := CDPOverlayGetSourceOrderHighlightObjectForTest(nodeId)
	log.Println("源码渲染顺序高亮数据：", highlight)

	// 可用于排查：order属性是否导致顺序错乱

	// ========== 示例3：自动化元素顺序测试 ==========
	func TestElementOrder(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 获取节点顺序高亮数据
		nodeId := 132
		data, err := CDPOverlayGetSourceOrderHighlightObjectForTest(nodeId)
		if err != nil {
			t.Fatal(err)
		}

		// 断言：元素顺序必须符合预期
		if data == "" {
			t.Error("元素顺序渲染异常")
		}
		t.Log("元素源码顺序渲染正常")
	}
}

*/

// -----------------------------------------------  Overlay.hideHighlight  -----------------------------------------------
// === 应用场景 ===
// 1. 快速清除高亮：调试完元素后隐藏标注，不关闭Overlay功能
// 2. 截图前清理：页面截图前临时清除高亮，保证画面干净
// 3. 多元素调试：切换调试不同元素前，清除上一个高亮
// 4. 自动化演示：演示过程中按需隐藏/显示调试标记
// 5. 录屏净化：录屏时临时隐藏高亮，不影响画面观感
// 6. 调试流程切换：从元素调试切换到布局调试时清空标记

// CDPOverlayHideHighlight 隐藏当前所有覆盖层高亮
func CDPOverlayHideHighlight() (string, error) {
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
		"method": "Overlay.hideHighlight"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.hideHighlight 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.hideHighlight 请求超时")
		}
	}
}

/*

// ==================== Overlay.hideHighlight 使用示例 ====================
func ExampleCDPOverlayHideHighlight() {
	// ========== 示例1：基础隐藏高亮 ==========
	resp, err := CDPOverlayHideHighlight()
	if err != nil {
		log.Fatalf("隐藏高亮失败: %v", err)
	}
	log.Printf("隐藏页面高亮成功，Overlay仍保持启用: %s", resp)

	// ========== 示例2：截图前快速清理高亮（推荐） ==========
	/*
	// 启用覆盖层（保持开启）
	CDPOverlayEnable()

	// 1. 高亮目标元素
	// CDPOverlayHighlightNode(123)

	// 2. 截图前隐藏所有高亮
	CDPOverlayHideHighlight()
	log.Println("已清除高亮，开始纯净截图")

	// 3. 截图完成后可继续高亮其他元素，无需重新enable

	// ========== 示例3：多元素调试切换 ==========
	// 调试第一个元素
	// CDPOverlayHighlightNode(node1)

	// 切换调试前清除
	CDPOverlayHideHighlight()

	// 调试第二个元素
	// CDPOverlayHighlightNode(node2)
	log.Println("切换元素调试，清除历史标记完成")

	// ========== 示例4：自动化测试用例 ==========
	func TestElementDebug(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 高亮 → 验证 → 隐藏
		// CDPOverlayHighlightNode(nodeId)
		// 执行验证逻辑

		CDPOverlayHideHighlight()
		t.Log("验证完成，已隐藏高亮标记")
	}
}

*/

// -----------------------------------------------  Overlay.highlightNode  -----------------------------------------------
// === 应用场景 ===
// 1. 元素可视化调试：给目标DOM节点添加高亮边框，快速定位元素
// 2. 自动化测试验证：高亮目标元素，确认元素是否正确加载
// 3. 页面元素检查：调试页面布局、错位、隐藏元素问题
// 4. 教学演示：实时标记页面元素，直观展示DOM结构
// 5. 自动化截图标注：截图前高亮元素，让报告更清晰
// 6. 前端问题排查：定位点击失效、渲染异常的元素

// CDPOverlayHighlightNode 高亮指定DOM节点
// nodeId: 要高亮的节点ID
// color: 高亮颜色 (如 "rgba(255,0,0,0.5)")
func CDPOverlayHighlightNode(nodeId int, color string) (string, error) {
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
		"method": "Overlay.highlightNode",
		"params": {
			"nodeId": %d,
			"highlightConfig": {
				"showInfo": true,
				"showRulers": false,
				"borderColor": %s,
				"contentColor": %s
			}
		}
	}`, reqID, nodeId, color, color)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.highlightNode 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.highlightNode 请求超时")
		}
	}
}

/*

// ==================== Overlay.highlightNode 使用示例 ====================
func ExampleCDPOverlayHighlightNode() {
	// 前提：已启用Overlay，已获取有效nodeId
	nodeID := 128
	redColor := `{"r":255,"g":0,"b":0,"a":0.6}`

	// ========== 示例1：基础高亮节点 ==========
	resp, err := CDPOverlayHighlightNode(nodeID, redColor)
	if err != nil {
		log.Fatalf("高亮节点失败: %v", err)
	}
	log.Printf("高亮元素成功，页面已显示红色边框: %s", resp)

	// ========== 示例2：完整调试流程（启用→高亮→隐藏→关闭） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 高亮目标元素
	nodeId := 130
	color := `{"r":0,"g":200,"b":255,"a":0.6}`
	CDPOverlayHighlightNode(nodeId, color)
	log.Println("元素已高亮，可查看调试")

	// 3. 停留1秒后隐藏高亮
	time.Sleep(1 * time.Second)
	CDPOverlayHideHighlight()
	log.Println("已隐藏高亮，调试完成")

	// ========== 示例3：自动化测试用例 ==========
	func TestButtonElement(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 高亮按钮元素
		nodeId := 135
		color := `{"r":50,"g":255,"b":50,"a":0.7}`
		_, err := CDPOverlayHighlightNode(nodeId, color)
		if err != nil {
			t.Fatalf("元素高亮失败，元素不存在: %v", err)
		}

		// 验证高亮成功
		t.Log("目标按钮元素已正常高亮，测试通过")
	}
}

*/

// -----------------------------------------------  Overlay.highlightQuad  -----------------------------------------------
// === 应用场景 ===
// 1. 区域渲染调试：高亮页面任意坐标区域，验证渲染范围
// 2. 坐标可视化：将计算出的坐标四边形显示在页面上，直观调试
// 3. 自动化区域校验：确认元素/视图区域是否符合预期坐标
// 4. 截图区域验证：截图前高亮截选区域，确认范围正确
// 5. 自定义框选调试：非元素区域、画布区域高亮调试
// 6. 图层/布局边界验证：可视化边界坐标是否正确

// CDPOverlayHighlightQuad 高亮自定义四边形区域
// quad: 四边形坐标数组，格式 [x1,y1,x2,y2,x3,y3,x4,y4]
// color: 高亮颜色 RGBA 对象
func CDPOverlayHighlightQuad(quad string, color string) (string, error) {
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
		"method": "Overlay.highlightQuad",
		"params": {
			"quad": %s,
			"color": %s,
			"fillColor": {"r":0,"g":0,"b":0,"a":0}
		}
	}`, reqID, quad, color)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.highlightQuad 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.highlightQuad 请求超时")
		}
	}
}

/*

// ==================== Overlay.highlightQuad 使用示例 ====================
func ExampleCDPOverlayHighlightQuad() {
	// ========== 示例1：高亮矩形区域（标准四边形） ==========
	// 坐标格式：[x1,y1, x2,y2, x3,y3, x4,y4]
	quad := `[100,100, 400,100, 400,300, 100,300]`
	// 蓝色半透明
	color := `{"r":0,"g":100,"b":255,"a":0.7}`

	resp, err := CDPOverlayHighlightQuad(quad, color)
	if err != nil {
		log.Fatalf("高亮四边形失败: %v", err)
	}
	log.Printf("自定义区域高亮成功: %s", resp)

	// ========== 示例2：完整调试流程（启用→高亮→隐藏） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 高亮目标区域
	area := `[50,50, 500,50, 500,400, 50,400]`
	blue := `{"r":0,"g":120,"b":255,"a":0.6}`
	CDPOverlayHighlightQuad(area, blue)
	log.Println("区域已高亮")

	// 3. 调试完成后隐藏
	time.Sleep(2 * time.Second)
	CDPOverlayHideHighlight()

	// ========== 示例3：截图区域验证 ==========
	func TestScreenshotArea(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 高亮即将截图的区域
		area := `[0,0, 800,600]`
		CDPOverlayHighlightQuad(area, `{"r":255,"g":50,"b":50,"a":0.5}`)

		t.Log("截图区域已可视化验证")
	}
}


*/

// -----------------------------------------------  Overlay.highlightRect  -----------------------------------------------
// === 应用场景 ===
// 1. 区域可视化调试：高亮页面任意矩形区域，验证坐标/尺寸是否正确
// 2. 截图预览：截图前高亮目标区域，确认范围无误
// 3. 元素边界校验：不依赖DOM节点，直接用坐标高亮区域
// 4. 自动化布局测试：验证页面区块位置、大小符合预期
// 5. 自定义标注：给页面特定区块做临时标记、教学演示
// 6. 渲染异常排查：可视化查看区域是否被遮挡、错位

// CDPOverlayHighlightRect 高亮矩形区域
// x: 左上角X坐标
// y: 左上角Y坐标
// width: 宽度
// height: 高度
// color: 边框颜色 RGBA 字符串
func CDPOverlayHighlightRect(x int, y int, width int, height int, color string) (string, error) {
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
		"method": "Overlay.highlightRect",
		"params": {
			"x": %d,
			"y": %d,
			"width": %d,
			"height": %d,
			"color": %s,
			"fillColor": {"r":0,"g":0,"b":0,"a":0}
		}
	}`, reqID, x, y, width, height, color)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.highlightRect 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.highlightRect 请求超时")
		}
	}
}

/*

// ==================== Overlay.highlightRect 使用示例 ====================
func ExampleCDPOverlayHighlightRect() {
	// ========== 示例1：基础高亮矩形区域 ==========
	color := `{"r":0,"g":160,"b":255,"a":0.7}` // 蓝色半透明
	resp, err := CDPOverlayHighlightRect(100, 100, 400, 200, color)
	if err != nil {
		log.Fatalf("高亮矩形失败: %v", err)
	}
	log.Printf("矩形高亮成功: %s", resp)

	// ========== 示例2：完整调试流程（启用→高亮→隐藏→关闭） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 高亮矩形区域
	redColor := `{"r":255,"g":50,"b":50,"a":0.6}`
	CDPOverlayHighlightRect(50, 50, 600, 300, redColor)
	log.Println("矩形区域已高亮")

	// 3. 2秒后隐藏高亮
	time.Sleep(2 * time.Second)
	CDPOverlayHideHighlight()
	log.Println("已清除矩形高亮")


	// ========== 示例3：截图前预览区域 ==========
	func TestScreenshotArea(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 高亮截图范围
		CDPOverlayHighlightRect(0, 0, 1200, 800, `{"r":50,"g":255,"b":50,"a":0.5}`)
		t.Log("截图区域已预览验证")
	}
}

*/

// -----------------------------------------------  Overlay.highlightSourceOrder  -----------------------------------------------
// === 应用场景 ===
// 1. Flex/Grid 顺序调试：验证子元素源码顺序 vs 视觉显示顺序
// 2. 浮动/定位错乱排查：定位/浮动导致元素视觉顺序与HTML顺序不符
// 3. 无障碍/屏幕阅读器校验：确保源码顺序符合语义与阅读逻辑
// 4. 组件嵌套顺序验证：复杂组件中子元素渲染顺序是否符合预期
// 5. 动态插入节点调试：JS动态添加元素后，确认源码顺序是否正确
// 6. 教学演示：直观展示HTML结构与视觉呈现的顺序关系

// CDPOverlayHighlightSourceOrder 高亮节点子元素的源码顺序
// nodeID: 目标DOM节点ID（其子元素将被标注顺序）
// color: 顺序数字与边框颜色（RGBA）
func CDPOverlayHighlightSourceOrder(nodeID int, color string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建SourceOrderConfig（样式配置）
	sourceOrderConfig := fmt.Sprintf(`{
		"color": %s,
		"fillColor": {"r":0,"g":0,"b":0,"a":0},
		"showText": true
	}`, color)

	// 构建CDP消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Overlay.highlightSourceOrder",
		"params": {
			"nodeId": %d,
			"sourceOrderConfig": %s
		}
	}`, reqID, nodeID, sourceOrderConfig)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.highlightSourceOrder 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.highlightSourceOrder 请求超时")
		}
	}
}

/*

// ==================== Overlay.highlightSourceOrder 使用示例 ====================
func ExampleCDPOverlayHighlightSourceOrder() {
	// ========== 前提：先获取目标容器节点ID ==========
	// 示例：通过DOM.querySelector获取容器节点ID
	containerNodeID, err := CDPDOMGetNodeIDBySelector("body > div.container")
	if err != nil {
		log.Fatalf("获取容器节点ID失败: %v", err)
	}

	// ========== 示例1：基础高亮源码顺序 ==========
	// 蓝色半透明
	color := `{"r":0,"g":120,"b":255,"a":0.8}`
	resp, err := CDPOverlayHighlightSourceOrder(containerNodeID, color)
	if err != nil {
		log.Fatalf("高亮源码顺序失败: %v", err)
	}
	log.Printf("源码顺序高亮成功: %s", resp)

	// ========== 示例2：完整调试流程（启用→高亮→隐藏） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 获取Flex容器节点ID
	flexNodeID, _ := CDPDOMGetNodeIDBySelector(".flex-container")

	// 3. 高亮Flex子元素源码顺序（红色标注）
	redColor := `{"r":255,"g":50,"b":50,"a":0.7}`
	CDPOverlayHighlightSourceOrder(flexNodeID, redColor)
	log.Println("Flex子元素源码顺序已标注")

	// 4. 3秒后清除高亮
	time.Sleep(3 * time.Second)
	CDPOverlayHideHighlight()
	log.Println("已清除源码顺序高亮")

	// ========== 示例3：自动化顺序校验测试 ==========
	func TestFlexSourceOrder(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 获取Grid容器
		gridNodeID, err := CDPDOMGetNodeIDBySelector(".grid-layout")
		if err != nil {
			t.Fatal(err)
		}

		// 高亮顺序
		CDPOverlayHighlightSourceOrder(gridNodeID, `{"r":50,"g":255,"b":50,"a":0.6}`)
		t.Log("Grid子元素源码顺序已可视化，请校验视觉顺序是否符合预期")
	}
}

*/

// -----------------------------------------------  Overlay.setInspectMode  -----------------------------------------------
// === 应用场景 ===
// 1. 自定义调试工具：实现类似Chrome开发者工具的“选择元素”功能
// 2. 自动化元素拾取：鼠标悬浮自动获取节点信息、高亮元素
// 3. 可视化配置工具：让用户手动点击选择页面配置区域
// 4. 教学演示：实时展示鼠标拾取的DOM节点结构
// 5. 页面分析工具：自动捕获鼠标指向的元素信息
// 6. 错误排查：快速定位页面点击/渲染异常的目标元素

// CDPOverlaySetInspectMode 设置元素检查模式
// mode: 检查模式（searchForNode/searchForUAShadowDOM/none）
// highlightConfig: 高亮配置（RGBA颜色JSON）
func CDPOverlaySetInspectMode(mode string, highlightConfig string) (string, error) {
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
		"method": "Overlay.setInspectMode",
		params: {
			"mode": "%s",
			"highlightConfig": %s
		}
	}`, reqID, mode, highlightConfig)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setInspectMode 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setInspectMode 请求超时")
		}
	}
}

/*


// ==================== Overlay.setInspectMode 使用示例 ====================
func ExampleCDPOverlaySetInspectMode() {
	// ========== 示例1：开启元素拾取模式（标准蓝色高亮） ==========
	// 模式：searchForNode = 拾取普通DOM节点
	mode := "searchForNode"
	// 蓝色高亮配置
	highlightConfig := `{
		"borderColor": {"r":0,"g":120,"b":255,"a":0.8},
		"contentColor": {"r":0,"g":0,"b":0,"a":0},
		"showInfo": true
	}`

	resp, err := CDPOverlaySetInspectMode(mode, highlightConfig)
	if err != nil {
		log.Fatalf("开启检查模式失败: %v", err)
	}
	log.Printf("已开启元素拾取模式，鼠标悬浮自动高亮: %s", resp)

	// ========== 示例2：关闭元素拾取模式 ==========
	resp, err := CDPOverlaySetInspectMode("none", "{}")
	if err != nil {
		log.Fatalf("关闭检查模式失败: %v", err)
	}
	log.Println("已关闭元素拾取，恢复正常鼠标")
	// ========== 示例3：完整自定义调试工具流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启元素拾取（红色高亮）
	mode := "searchForNode"
	redConfig := `{
		"borderColor": {"r":255,"g":50,"b":50,"a":0.7},
		"showInfo": true,
		"showRulers": true
	}`
	CDPOverlaySetInspectMode(mode, redConfig)
	log.Println("自定义元素拾取已启动，开始选择元素")

	// 3. 使用完成后关闭拾取
	// CDPOverlaySetInspectMode("none", "{}")

	// ========== 示例4：自动化拾取元素测试 ==========
	func TestElementPicker(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 启动拾取
		CDPOverlaySetInspectMode("searchForNode", `{"borderColor":{"r":50,"g":255,"b":50,"a":0.6}}`)
		t.Log("元素拾取模式已启动，可手动选择页面元素")

		// 测试完成关闭
		defer CDPOverlaySetInspectMode("none", "{}")
	}
}

*/

// -----------------------------------------------  Overlay.setPausedInDebuggerMessage  -----------------------------------------------
// === 应用场景 ===
// 1. 调试器可视化：JS断点暂停时，页面显示清晰暂停提示
// 2. 自动化调试：调试过程中给页面添加暂停状态标记
// 3. 教学演示：直观展示代码执行已暂停在断点处
// 4. 错误排查：明确告知用户当前脚本已暂停运行
// 5. 自定义调试工具：实现专业调试器的暂停状态UI
// 6. 录屏演示：让观看者清晰看到调试暂停节点

// CDPOverlaySetPausedInDebuggerMessage 设置调试暂停提示信息
// message: 要显示的文字（传nil/空则隐藏提示）
func CDPOverlaySetPausedInDebuggerMessage(message string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（message为nil时隐藏提示）
	var messageParam string
	if message == "" {
		messageParam = "null"
	} else {
		messageParam = fmt.Sprintf(`"%s"`, message)
	}

	messageBody := fmt.Sprintf(`{
		"id": %d,
		"method": "Overlay.setPausedInDebuggerMessage",
		"params": {
			"message": %s
		}
	}`, reqID, messageParam)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(messageBody))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setPausedInDebuggerMessage 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", messageBody)

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
			return "", fmt.Errorf("Overlay.setPausedInDebuggerMessage 请求超时")
		}
	}
}

/*

// ==================== Overlay.setPausedInDebuggerMessage 使用示例 ====================
func ExampleCDPOverlaySetPausedInDebuggerMessage() {
	// ========== 示例1：显示调试暂停提示 ==========
	resp, err := CDPOverlaySetPausedInDebuggerMessage("已暂停执行 - 等待调试")
	if err != nil {
		log.Fatalf("显示暂停提示失败: %v", err)
	}
	log.Printf("已显示调试暂停覆盖层: %s", resp)

	// ========== 示例2：隐藏暂停提示 ==========
	resp, err := CDPOverlaySetPausedInDebuggerMessage("")
	if err != nil {
		log.Fatalf("隐藏暂停提示失败: %v", err)
	}
	log.Println("已清除调试暂停提示")

	// ========== 示例3：完整调试流程（启用→显示→隐藏） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. JS断点触发时显示提示
	CDPOverlaySetPausedInDebuggerMessage("代码已暂停｜点击继续执行")
	log.Println("调试暂停，页面已显示提示")

	// 3. 调试结束后清除提示
	// time.Sleep(3 * time.Second)
	CDPOverlaySetPausedInDebuggerMessage("")
	log.Println("已恢复执行，清除暂停提示")

	// ========== 示例4：自动化调试测试 ==========
	func TestDebugPause(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 模拟断点暂停
		CDPOverlaySetPausedInDebuggerMessage("自动化测试：脚本已暂停")
		t.Log("暂停提示已显示")

		// 测试结束关闭
		defer CDPOverlaySetPausedInDebuggerMessage("")
	}
}

*/

// -----------------------------------------------  Overlay.setShowAdHighlights  -----------------------------------------------
// === 应用场景 ===
// 1. 广告区域检测：自动识别并高亮页面所有广告位，便于分析页面结构
// 2. 广告屏蔽验证：验证广告屏蔽插件/规则是否生效，查看未被屏蔽的广告
// 3. 页面净化测试：检查页面纯净度，统计广告区域占比
// 4. 前端合规检查：确认页面广告展示是否符合规范
// 5. 自动化爬虫辅助：定位广告区域，避免采集广告内容
// 6. 教学演示：直观展示页面中广告与正常内容的区分

// CDPOverlaySetShowAdHighlights 设置是否高亮广告区域
// show: true 开启广告高亮，false 关闭广告高亮
func CDPOverlaySetShowAdHighlights(show bool) (string, error) {
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
		"method": "Overlay.setShowAdHighlights",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowAdHighlights 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowAdHighlights 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowAdHighlights 使用示例 ====================
func ExampleCDPOverlaySetShowAdHighlights() {
	// ========== 示例1：开启广告高亮 ==========
	resp, err := CDPOverlaySetShowAdHighlights(true)
	if err != nil {
		log.Fatalf("开启广告高亮失败: %v", err)
	}
	log.Printf("已开启广告区域高亮，页面广告已标记: %s", resp)

	// ========== 示例2：关闭广告高亮 ==========
	resp, err := CDPOverlaySetShowAdHighlights(false)
	if err != nil {
		log.Fatalf("关闭广告高亮失败: %v", err)
	}
	log.Println("已关闭广告区域标记，恢复正常显示")

	// ========== 示例3：完整广告检测流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启广告高亮，检测页面广告
	_, err := CDPOverlaySetShowAdHighlights(true)
	if err != nil {
		log.Fatalf("广告检测失败: %v", err)
	}
	log.Println("广告区域已高亮，可查看页面广告分布")

	// 3. 检测完成后关闭
	// time.Sleep(3 * time.Second)
	// CDPOverlaySetShowAdHighlights(false)

	// ========== 示例4：广告屏蔽效果测试用例 ==========
	func TestAdBlockEffect(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启广告高亮
		CDPOverlaySetShowAdHighlights(true)
		// 断言：页面应无广告高亮区域（广告已被屏蔽）
		t.Log("已开启广告检测，验证广告屏蔽效果")

		// 测试完成关闭
		defer CDPOverlaySetShowAdHighlights(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowContainerQueryOverlays  -----------------------------------------------
// === 应用场景 ===
// 1. 容器查询调试：可视化查看哪些元素开启了 container-type / container-name
// 2. 响应式布局校验：确认容器查询是否正确应用、生效
// 3. 自动化布局测试：验证容器查询渲染是否符合预期
// 4. 教学演示：直观展示 CSS Container Queries 工作机制
// 5. 复杂页面调试：快速定位所有使用容器查询的父容器
// 6. UI 组件库调试：检查组件容器查询配置是否正确

// CDPOverlaySetShowContainerQueryOverlays 设置是否显示容器查询覆盖层
// show: true 显示容器查询覆盖层，false 关闭
func CDPOverlaySetShowContainerQueryOverlays(show bool) (string, error) {
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
		"method": "Overlay.setShowContainerQueryOverlays",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowContainerQueryOverlays 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowContainerQueryOverlays 请求超时")
		}
	}
}

/*
// ==================== Overlay.setShowContainerQueryOverlays 使用示例 ====================
func ExampleCDPOverlaySetShowContainerQueryOverlays() {
	// ========== 示例1：开启容器查询覆盖层 ==========
	resp, err := CDPOverlaySetShowContainerQueryOverlays(true)
	if err != nil {
		log.Fatalf("开启容器查询覆盖层失败: %v", err)
	}
	log.Printf("已开启容器查询可视化，所有容器已标注: %s", resp)

	// ========== 示例2：关闭容器查询覆盖层 ==========
	resp, err := CDPOverlaySetShowContainerQueryOverlays(false)
	if err != nil {
		log.Fatalf("关闭容器查询覆盖层失败: %v", err)
	}
	log.Println("已关闭容器查询可视化")

	// ========== 示例3：完整容器查询调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启容器查询可视化
	_, err := CDPOverlaySetShowContainerQueryOverlays(true)
	if err != nil {
		log.Fatalf("启动容器查询调试失败: %v", err)
	}
	log.Println("容器查询调试已启动，可查看所有容器标记")

	// 3. 调试完成后关闭
	// CDPOverlaySetShowContainerQueryOverlays(false)

	// ========== 示例4：自动化容器查询测试用例 ==========
	func TestContainerQuery(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启可视化
		CDPOverlaySetShowContainerQueryOverlays(true)
		t.Log("容器查询覆盖层已启用，验证布局是否正确")

		// 测试结束关闭
		defer CDPOverlaySetShowContainerQueryOverlays(false)
	}
}


*/

// -----------------------------------------------  Overlay.setShowDebugBorders  -----------------------------------------------
// === 应用场景 ===
// 1. 渲染性能调试：查看图层、复合层、重绘区域边界
// 2. 图层结构分析：识别页面会被GPU分层的元素
// 3. 滚动区域调试：可视化滚动容器与滚动裁剪边界
// 4. 渲染异常排查：定位闪烁、卡顿、图层错乱问题
// 5. 自动化性能审计：验证页面是否存在过多图层创建
// 6. CSS transform/opacity 调试：查看硬件加速区域

// CDPOverlaySetShowDebugBorders 设置是否显示渲染调试边框
// show: true 显示调试边框，false 关闭调试边框
func CDPOverlaySetShowDebugBorders(show bool) (string, error) {
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
		"method": "Overlay.setShowDebugBorders",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowDebugBorders 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowDebugBorders 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowDebugBorders 使用示例 ====================
func ExampleCDPOverlaySetShowDebugBorders() {
	// ========== 示例1：开启渲染调试边框 ==========
	resp, err := CDPOverlaySetShowDebugBorders(true)
	if err != nil {
		log.Fatalf("开启调试边框失败: %v", err)
	}
	log.Printf("已开启渲染调试边框: %s", resp)

	// ========== 示例2：关闭调试边框 ==========
	resp, err := CDPOverlaySetShowDebugBorders(false)
	if err != nil {
		log.Fatalf("关闭调试边框失败: %v", err)
	}
	log.Println("已关闭渲染调试边框")

	// ========== 示例3：完整性能渲染调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启渲染边界调试
	_, err := CDPOverlaySetShowDebugBorders(true)
	if err != nil {
		log.Fatalf("启动渲染调试失败: %v", err)
	}
	log.Println("已显示图层、复合层、裁剪区域边界")

	// 3. 执行页面滚动、动画等操作观察渲染
	// 4. 调试完成关闭
	// CDPOverlaySetShowDebugBorders(false)

	// ========== 示例4：自动化性能审计测试用例 ==========
	func TestRenderLayers(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启图层调试边框
		CDPOverlaySetShowDebugBorders(true)
		t.Log("渲染调试边框已启用，检查页面图层数量是否合理")

		// 测试结束关闭
		defer CDPOverlaySetShowDebugBorders(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowFlexOverlays  -----------------------------------------------
// === 应用场景 ===
// 1. Flex 布局调试：可视化查看主轴、交叉轴、对齐方式、gap 间距
// 2. 布局错乱排查：快速定位 justify/align/align-self 失效问题
// 3. 响应式布局验证：不同屏幕下 Flex 排列是否符合预期
// 4. 自动化布局测试：验证 Flex 容器渲染是否正确
// 5. 教学演示：直观展示 Flex 布局工作原理
// 6. 复杂页面调试：快速找到所有 Flex 容器

// CDPOverlaySetShowFlexOverlays 设置是否显示 Flex 布局覆盖层
// show: true 显示，false 关闭
func CDPOverlaySetShowFlexOverlays(show bool) (string, error) {
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
		"method": "Overlay.setShowFlexOverlays",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowFlexOverlays 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowFlexOverlays 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowFlexOverlays 使用示例 ====================
func ExampleCDPOverlaySetShowFlexOverlays() {
	// ========== 示例1：开启 Flex 布局覆盖层 ==========
	resp, err := CDPOverlaySetShowFlexOverlays(true)
	if err != nil {
		log.Fatalf("开启 Flex 覆盖层失败: %v", err)
	}
	log.Printf("已开启 Flex 布局可视化调试: %s", resp)

	// ========== 示例2：关闭 Flex 布局覆盖层 ==========
	resp, err := CDPOverlaySetShowFlexOverlays(false)
	if err != nil {
		log.Fatalf("关闭 Flex 覆盖层失败: %v", err)
	}
	log.Println("已关闭 Flex 布局覆盖层")

	// ========== 示例3：完整 Flex 调试流程（推荐） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启 Flex 可视化
	_, err := CDPOverlaySetShowFlexOverlays(true)
	if err != nil {
		log.Fatalf("启动 Flex 调试失败: %v", err)
	}
	log.Println("已显示所有 Flex 容器、轴、间距、对齐信息")

	// 3. 调试完成后关闭
	// CDPOverlaySetShowFlexOverlays(false)

	// ========== 示例4：自动化 Flex 布局测试用例 ==========
	func TestFlexLayout(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启覆盖层
		CDPOverlaySetShowFlexOverlays(true)
		t.Log("Flex 布局调试已启用，验证排列与对齐")

		// 测试结束自动关闭
		defer CDPOverlaySetShowFlexOverlays(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowFPSCounter  -----------------------------------------------
// === 应用场景 ===
// 1. 页面流畅度监控：实时查看 FPS 判断页面是否卡顿
// 2. 动画性能调试：检测 CSS/JS 动画是否掉帧
// 3. 滚动性能分析：检查滚动时帧率是否稳定在 60fps
// 4. 渲染性能优化：定位长任务、频繁重绘导致的性能问题
// 5. 自动化性能测试：记录帧率数据生成性能报告
// 6. 游戏/H5 性能验收：确保交互场景帧率达标

// CDPOverlaySetShowFPSCounter 设置是否显示帧率计数器
// show: true 显示FPS计数器，false 隐藏
func CDPOverlaySetShowFPSCounter(show bool) (string, error) {
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
		"method": "Overlay.setShowFPSCounter",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowFPSCounter 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowFPSCounter 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowFPSCounter 使用示例 ====================
func ExampleCDPOverlaySetShowFPSCounter() {
	// ========== 示例1：显示FPS性能计数器 ==========
	resp, err := CDPOverlaySetShowFPSCounter(true)
	if err != nil {
		log.Fatalf("开启FPS计数器失败: %v", err)
	}
	log.Printf("已显示FPS计数器，实时监控页面流畅度: %s", resp)

	// ========== 示例2：隐藏FPS计数器 ==========
	resp, err := CDPOverlaySetShowFPSCounter(false)
	if err != nil {
		log.Fatalf("关闭FPS计数器失败: %v", err)
	}
	log.Println("已隐藏FPS性能计数器")

	// ========== 示例3：完整性能监控流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启FPS监控
	_, err := CDPOverlaySetShowFPSCounter(true)
	if err != nil {
		log.Fatalf("启动FPS监控失败: %v", err)
	}
	log.Println("FPS监控已启动，可测试滚动、动画等场景")

	// 3. 测试完成关闭
	// CDPOverlaySetShowFPSCounter(false)

	// ========== 示例4：自动化流畅度测试用例 ==========
	func TestPageFPS(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启帧率监控
		CDPOverlaySetShowFPSCounter(true)
		t.Log("FPS计数器已显示，验证页面动画/滚动流畅度")

		// 测试结束关闭
		defer CDPOverlaySetShowFPSCounter(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowGridOverlays  -----------------------------------------------
// === 应用场景 ===
// 1. Grid 布局调试：可视化网格线、行列、gap、区域名称
// 2. 布局错乱排查：定位 grid-template-areas / 行列尺寸问题
// 3. 响应式布局验证：不同屏幕下网格自适应是否符合预期
// 4. 自动化布局测试：验证 Grid 容器渲染是否正确
// 5. 教学演示：直观展示 CSS Grid 布局工作原理
// 6. 复杂页面调试：快速找到所有 Grid 容器

// CDPOverlaySetShowGridOverlays 设置是否显示 Grid 布局覆盖层
// show: true 显示，false 关闭
func CDPOverlaySetShowGridOverlays(show bool) (string, error) {
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
		"method": "Overlay.setShowGridOverlays",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowGridOverlays 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowGridOverlays 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowGridOverlays 使用示例 ====================
func ExampleCDPOverlaySetShowGridOverlays() {
	// ========== 示例1：开启 Grid 布局覆盖层 ==========
	resp, err := CDPOverlaySetShowGridOverlays(true)
	if err != nil {
		log.Fatalf("开启 Grid 覆盖层失败: %v", err)
	}
	log.Printf("已开启 Grid 布局可视化调试: %s", resp)

	// ========== 示例2：关闭 Grid 布局覆盖层 ==========
	resp, err := CDPOverlaySetShowGridOverlays(false)
	if err != nil {
		log.Fatalf("关闭 Grid 覆盖层失败: %v", err)
	}
	log.Println("已关闭 Grid 布局覆盖层")

	// ========== 示例3：完整 Grid 调试流程（推荐） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启 Grid 可视化
	_, err := CDPOverlaySetShowGridOverlays(true)
	if err != nil {
		log.Fatalf("启动 Grid 调试失败: %v", err)
	}
	log.Println("已显示所有 Grid 网格线、行列、间距、区域信息")

	// 3. 调试完成后关闭
	// CDPOverlaySetShowGridOverlays(false)

	// ========== 示例4：自动化 Grid 布局测试用例 ==========
	func TestGridLayout(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启覆盖层
		CDPOverlaySetShowGridOverlays(true)
		t.Log("Grid 布局调试已启用，验证网格、区域、对齐")

		// 测试结束自动关闭
		defer CDPOverlaySetShowGridOverlays(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowHinge  -----------------------------------------------
// === 应用场景 ===
// 1. 元素销毁动画调试：预览元素合页脱落效果
// 2. 层叠与动画测试：验证元素动画执行状态
// 3. 调试效果演示：展示元素移除时的视觉动画
// 4. 自动化UI测试：验证元素动画是否正常触发
// 5. 教学演示：直观展示元素退场动画机制
// 6. 渲染异常排查：检查动画图层是否正常渲染

// CDPOverlaySetShowHinge 设置是否显示元素合页脱落效果
// show: true 开启hinge效果，false 关闭
func CDPOverlaySetShowHinge(show bool) (string, error) {
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
		"method": "Overlay.setShowHinge",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowHinge 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowHinge 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowHinge 使用示例 ====================
func ExampleCDPOverlaySetShowHinge() {
	// ========== 示例1：开启合页脱落效果 ==========
	resp, err := CDPOverlaySetShowHinge(true)
	if err != nil {
		log.Fatalf("开启hinge效果失败: %v", err)
	}
	log.Printf("已开启元素合页脱落调试效果: %s", resp)

	// ========== 示例2：关闭合页效果 ==========
	resp, err := CDPOverlaySetShowHinge(false)
	if err != nil {
		log.Fatalf("关闭hinge效果失败: %v", err)
	}
	log.Println("已关闭元素合页脱落动画效果")

	// ========== 示例3：完整调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启hinge效果
	_, err := CDPOverlaySetShowHinge(true)
	if err != nil {
		log.Fatalf("启动hinge调试失败: %v", err)
	}
	log.Println("元素合页脱落效果已启用")

	// 3. 调试完成关闭
	// CDPOverlaySetShowHinge(false)

	// ========== 示例4：自动化动画测试用例 ==========
	func TestElementHingeAnimation(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		CDPOverlaySetShowHinge(true)
		t.Log("合页脱落效果已启用，验证动画渲染")

		defer CDPOverlaySetShowHinge(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowInspectedElementAnchor  -----------------------------------------------
// === 应用场景 ===
// 1. 元素调试定位：给当前审查元素显示锚点，快速识别目标节点
// 2. 多元素调试区分：避免多个高亮混淆，锚点标记当前审查对象
// 3. 自定义调试工具：实现类似 DevTools 的选中锚点效果
// 4. 自动化元素验证：明确标识当前正在检测的元素
// 5. 教学演示：直观指出当前被审查的目标元素
// 6. 复杂页面排查：在多层嵌套结构中快速定位选中元素

// CDPOverlaySetShowInspectedElementAnchor 设置是否显示被审查元素锚点
// show: true 显示锚点标记，false 隐藏
func CDPOverlaySetShowInspectedElementAnchor(show bool) (string, error) {
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
		"method": "Overlay.setShowInspectedElementAnchor",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowInspectedElementAnchor 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowInspectedElementAnchor 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowInspectedElementAnchor 使用示例 ====================
func ExampleCDPOverlaySetShowInspectedElementAnchor() {
	// ========== 示例1：显示被审查元素锚点 ==========
	resp, err := CDPOverlaySetShowInspectedElementAnchor(true)
	if err != nil {
		log.Fatalf("显示锚点失败: %v", err)
	}
	log.Printf("已显示当前审查元素锚点: %s", resp)

	// ========== 示例2：隐藏审查元素锚点 ==========
	resp, err := CDPOverlaySetShowInspectedElementAnchor(false)
	if err != nil {
		log.Fatalf("隐藏锚点失败: %v", err)
	}
	log.Println("已隐藏审查元素锚点")


	// ========== 示例3：完整元素调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启锚点标记
	_, err := CDPOverlaySetShowInspectedElementAnchor(true)
	if err != nil {
		log.Fatalf("启动锚点显示失败: %v", err)
	}
	log.Println("已为当前审查元素显示锚点标记")

	// 3. 调试完成关闭
	// CDPOverlaySetShowInspectedElementAnchor(false)

	// ========== 示例4：自动化元素定位测试用例 ==========
	func TestInspectedElement(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启锚点标识
		CDPOverlaySetShowInspectedElementAnchor(true)
		t.Log("已显示当前审查元素锚点，定位更清晰")

		// 测试结束关闭
		defer CDPOverlaySetShowInspectedElementAnchor(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowIsolatedElements  -----------------------------------------------
// === 应用场景 ===
// 1. z-index 层级错乱调试：可视化查看哪些元素是独立堆叠上下文
// 2. 遮挡/覆盖问题排查：定位导致层级异常的隔离元素
// 3. 堆叠上下文分析：理解页面层级结构
// 4. 复杂弹窗/浮层调试：解决层级不生效问题
// 5. CSS 层级优化：识别不必要的隔离层
// 6. 自动化层级验证：校验页面层级结构是否符合预期

// CDPOverlaySetShowIsolatedElements 设置是否高亮隔离/堆叠上下文元素
// show: true 显示高亮，false 关闭
func CDPOverlaySetShowIsolatedElements(show bool) (string, error) {
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
		"method": "Overlay.setShowIsolatedElements",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowIsolatedElements 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowIsolatedElements 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowIsolatedElements 使用示例 ====================
func ExampleCDPOverlaySetShowIsolatedElements() {
	// ========== 示例1：开启隔离元素高亮 ==========
	resp, err := CDPOverlaySetShowIsolatedElements(true)
	if err != nil {
		log.Fatalf("开启隔离元素高亮失败: %v", err)
	}
	log.Printf("已高亮所有隔离/堆叠上下文元素: %s", resp)

	// ========== 示例2：关闭隔离元素高亮 ==========
	resp, err := CDPOverlaySetShowIsolatedElements(false)
	if err != nil {
		log.Fatalf("关闭隔离元素高亮失败: %v", err)
	}
	log.Println("已关闭隔离元素覆盖层")

	// ========== 示例3：完整层级调试流程（推荐） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启隔离元素可视化
	_, err := CDPOverlaySetShowIsolatedElements(true)
	if err != nil {
		log.Fatalf("启动层级调试失败: %v", err)
	}
	log.Println("已显示所有堆叠上下文、隔离元素")

	// 3. 调试 z-index 遮挡问题
	// 4. 完成后关闭
	// CDPOverlaySetShowIsolatedElements(false)

	// ========== 示例4：自动化层级问题测试用例 ==========
	func TestZIndexIsolation(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启隔离检测
		CDPOverlaySetShowIsolatedElements(true)
		t.Log("已开启层级隔离检测，检查z-index遮挡问题")

		// 测试结束关闭
		defer CDPOverlaySetShowIsolatedElements(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowLayoutShiftRegions  -----------------------------------------------
// === 应用场景 ===
// 1. CLS 布局偏移调试：直观看到页面哪里发生了抖动
// 2. 图片未设尺寸导致偏移排查
// 3. 动态插入内容导致布局抖动定位
// 4. 字体加载闪烁导致布局偏移分析
// 5. 自动化页面稳定性测试
// 6. 前端体验优化：降低布局偏移提升用户体验

// CDPOverlaySetShowLayoutShiftRegions 设置是否显示布局偏移区域
// show: true 显示布局偏移高亮，false 关闭
func CDPOverlaySetShowLayoutShiftRegions(show bool) (string, error) {
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
		"method": "Overlay.setShowLayoutShiftRegions",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowLayoutShiftRegions 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowLayoutShiftRegions 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowLayoutShiftRegions 使用示例 ====================
func ExampleCDPOverlaySetShowLayoutShiftRegions() {
	// ========== 示例1：开启布局偏移区域高亮 ==========
	resp, err := CDPOverlaySetShowLayoutShiftRegions(true)
	if err != nil {
		log.Fatalf("开启布局偏移高亮失败: %v", err)
	}
	log.Printf("已显示布局抖动区域: %s", resp)

	// ========== 示例2：关闭布局偏移高亮 ==========
	resp, err := CDPOverlaySetShowLayoutShiftRegions(false)
	if err != nil {
		log.Fatalf("关闭布局偏移高亮失败: %v", err)
	}
	log.Println("已关闭布局偏移覆盖层")

	// ========== 示例3：完整 CLS 调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启布局偏移监控
	_, err := CDPOverlaySetShowLayoutShiftRegions(true)
	if err != nil {
		log.Fatalf("启动 CLS 调试失败: %v", err)
	}
	log.Println("已可视化布局抖动，可滚动/交互观察偏移")

	// 3. 调试完成关闭
	// CDPOverlaySetShowLayoutShiftRegions(false)

	// ========== 示例4：自动化页面稳定性测试 ==========
	func TestCLS(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		CDPOverlaySetShowLayoutShiftRegions(true)
		t.Log("布局偏移检测已启用，验证页面稳定性")

		defer CDPOverlaySetShowLayoutShiftRegions(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowPaintRects  -----------------------------------------------
// === 应用场景 ===
// 1. 页面重绘优化：查看哪些区域频繁重绘、导致卡顿
// 2. 动画性能调试：检测 CSS/JS 动画是否触发大面积重绘
// 3. 滚动性能分析：滚动时是否触发不必要的全页面重绘
// 4. 渲染异常排查：定位闪烁、频繁刷新的区域
// 5. 自动化性能审计：验证页面重绘区域是否合理
// 6. 长列表/虚拟列表优化：检查是否只重绘可见区域

// CDPOverlaySetShowPaintRects 设置是否显示绘制矩形
// show: true 显示重绘区域，false 关闭
func CDPOverlaySetShowPaintRects(show bool) (string, error) {
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
		"method": "Overlay.setShowPaintRects",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowPaintRects 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowPaintRects 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowPaintRects 使用示例 ====================
func ExampleCDPOverlaySetShowPaintRects() {
	// ========== 示例1：开启重绘区域高亮（绿色矩形） ==========
	resp, err := CDPOverlaySetShowPaintRects(true)
	if err != nil {
		log.Fatalf("开启重绘检测失败: %v", err)
	}
	log.Printf("已显示页面重绘区域: %s", resp)

	// ========== 示例2：关闭重绘高亮 ==========
	resp, err := CDPOverlaySetShowPaintRects(false)
	if err != nil {
		log.Fatalf("关闭重绘检测失败: %v", err)
	}
	log.Println("已关闭页面重绘区域标记")

	// ========== 示例3：完整性能调试流程（推荐） ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启重绘可视化
	_, err := CDPOverlaySetShowPaintRects(true)
	if err != nil {
		log.Fatalf("启动重绘调试失败: %v", err)
	}
	log.Println("已开启重绘监控，滚动/交互查看绿色重绘区域")

	// 3. 调试完成关闭
	// CDPOverlaySetShowPaintRects(false)

	// ========== 示例4：自动化重绘性能测试用例 ==========
	func TestPagePaintPerformance(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启重绘检测
		CDPOverlaySetShowPaintRects(true)
		t.Log("重绘区域已显示，验证页面渲染性能")

		// 测试结束关闭
		defer CDPOverlaySetShowPaintRects(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowScrollBottleneckRects  -----------------------------------------------
// === 应用场景 ===
// 1. 滚动卡顿问题定位：找出导致滚动不流畅的瓶颈元素
// 2. 滚动性能优化：识别触发频繁布局/重绘的区域
// 3. 长列表滚动调试：定位滚动性能瓶颈
// 4. 固定定位/transform 滚动问题排查
// 5. 自动化滚动流畅度测试
// 6. 页面滚动性能审计

// CDPOverlaySetShowScrollBottleneckRects 设置是否显示滚动瓶颈区域
// show: true 显示瓶颈区域，false 关闭
func CDPOverlaySetShowScrollBottleneckRects(show bool) (string, error) {
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
		"method": "Overlay.setShowScrollBottleneckRects",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowScrollBottleneckRects 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowScrollBottleneckRects 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowScrollBottleneckRects 使用示例 ====================
func ExampleCDPOverlaySetShowScrollBottleneckRects() {
	// ========== 示例1：显示滚动瓶颈区域 ==========
	resp, err := CDPOverlaySetShowScrollBottleneckRects(true)
	if err != nil {
		log.Fatalf("开启滚动瓶颈高亮失败: %v", err)
	}
	log.Printf("已显示滚动瓶颈区域标记: %s", resp)

	// ========== 示例2：关闭滚动瓶颈区域 ==========
	resp, err := CDPOverlaySetShowScrollBottleneckRects(false)
	if err != nil {
		log.Fatalf("关闭滚动瓶颈高亮失败: %v", err)
	}
	log.Println("已关闭滚动瓶颈区域覆盖层")

	// ========== 示例3：完整滚动性能调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启滚动瓶颈检测
	_, err := CDPOverlaySetShowScrollBottleneckRects(true)
	if err != nil {
		log.Fatalf("启动滚动性能调试失败: %v", err)
	}
	log.Println("已显示滚动瓶颈，滚动页面查看卡顿区域")

	// 3. 调试完成关闭
	// CDPOverlaySetShowScrollBottleneckRects(false)

	// ========== 示例4：自动化滚动流畅度测试用例 ==========
	func TestScrollPerformance(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		CDPOverlaySetShowScrollBottleneckRects(true)
		t.Log("滚动瓶颈检测已启用，验证滚动流畅度")

		defer CDPOverlaySetShowScrollBottleneckRects(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowScrollSnapOverlays  -----------------------------------------------
// === 应用场景 ===
// 1. Scroll Snap 滚动吸附效果调试
// 2. 轮播图 / 横向滚动列表吸附位置校验
// 3. 滚动对齐点、吸附区域可视化
// 4. scroll-snap-type / scroll-snap-align 效果验证
// 5. 复杂滚动容器交互调试
// 6. 自动化滚动吸附规则测试

// CDPOverlaySetShowScrollSnapOverlays 设置是否显示滚动吸附调试覆盖层
// show: true 显示，false 关闭
func CDPOverlaySetShowScrollSnapOverlays(show bool) (string, error) {
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
		"method": "Overlay.setShowScrollSnapOverlays",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowScrollSnapOverlays 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowScrollSnapOverlays 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowScrollSnapOverlays 使用示例 ====================
func ExampleCDPOverlaySetShowScrollSnapOverlays() {
	// ========== 示例1：开启滚动吸附覆盖层 ==========
	resp, err := CDPOverlaySetShowScrollSnapOverlays(true)
	if err != nil {
		log.Fatalf("开启滚动吸附调试失败: %v", err)
	}
	log.Printf("已显示滚动吸附对齐区域: %s", resp)

	// ========== 示例2：关闭滚动吸附覆盖层 ==========
	resp, err := CDPOverlaySetShowScrollSnapOverlays(false)
	if err != nil {
		log.Fatalf("关闭滚动吸附调试失败: %v", err)
	}
	log.Println("已关闭滚动吸附覆盖层")

	// ========== 示例3：完整滚动吸附调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启 Scroll Snap 可视化
	_, err := CDPOverlaySetShowScrollSnapOverlays(true)
	if err != nil {
		log.Fatalf("启动滚动吸附调试失败: %v", err)
	}
	log.Println("已显示滚动吸附位置、对齐点、吸附区域")

	// 3. 调试完成关闭
	// CDPOverlaySetShowScrollSnapOverlays(false)

	// ========== 示例4：轮播图滚动吸附测试用例 ==========
	func TestCarouselScrollSnap(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		CDPOverlaySetShowScrollSnapOverlays(true)
		t.Log("滚动吸附调试已启用，验证轮播吸附是否正确")

		defer CDPOverlaySetShowScrollSnapOverlays(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowViewportSizeOnResize  -----------------------------------------------
// === 应用场景 ===
// 1. 响应式布局调试：实时查看窗口宽高变化
// 2. 媒体查询断点验证：确认哪些宽度触发样式变化
// 3. 移动端/平板适配测试
// 4. 页面自适应效果验证
// 5. 教学演示：直观展示视口尺寸变化
// 6. 自动化窗口大小测试

// CDPOverlaySetShowViewportSizeOnResize 设置调整窗口大小时是否显示视口尺寸
// show: true 显示尺寸，false 关闭
func CDPOverlaySetShowViewportSizeOnResize(show bool) (string, error) {
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
		"method": "Overlay.setShowViewportSizeOnResize",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowViewportSizeOnResize 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowViewportSizeOnResize 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowViewportSizeOnResize 使用示例 ====================
func ExampleCDPOverlaySetShowViewportSizeOnResize() {
	// ========== 示例1：开启窗口大小变化显示视口尺寸 ==========
	resp, err := CDPOverlaySetShowViewportSizeOnResize(true)
	if err != nil {
		log.Fatalf("开启视口尺寸显示失败: %v", err)
	}
	log.Printf("已启用：拖动窗口时实时显示宽高: %s", resp)

	// ========== 示例2：关闭视口尺寸显示 ==========
	resp, err := CDPOverlaySetShowViewportSizeOnResize(false)
	if err != nil {
		log.Fatalf("关闭视口尺寸显示失败: %v", err)
	}
	log.Println("已关闭窗口大小变化提示")

	// ========== 示例3：完整响应式调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启实时尺寸显示
	_, err := CDPOverlaySetShowViewportSizeOnResize(true)
	if err != nil {
		log.Fatalf("启动响应式调试失败: %v", err)
	}
	log.Println("拖动浏览器窗口，右上角会实时显示 width × height")

	// 3. 测试完成关闭
	// CDPOverlaySetShowViewportSizeOnResize(false)

	// ========== 示例4：自动化响应式断点测试 ==========
	func TestResponsiveBreakpoints(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		CDPOverlaySetShowViewportSizeOnResize(true)
		t.Log("已开启视口尺寸实时显示，验证媒体查询断点")

		defer CDPOverlaySetShowViewportSizeOnResize(false)
	}
}

*/

// -----------------------------------------------  Overlay.setShowWindowControlsOverlay  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 窗口控件区域调试
// 2. 桌面端安装应用标题栏适配
// 3. 避免内容被窗口按钮遮挡
// 4. 验证 title bar 区域样式与占位
// 5. PWA 安装体验优化
// 6. 自动化窗口控件区域布局测试

// CDPOverlaySetShowWindowControlsOverlay 设置是否显示窗口控件区域覆盖层
// show: true 显示，false 关闭
func CDPOverlaySetShowWindowControlsOverlay(show bool) (string, error) {
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
		"method": "Overlay.setShowWindowControlsOverlay",
		"params": {
			"show": %t
		}
	}`, reqID, show)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Overlay.setShowWindowControlsOverlay 请求失败: %w", err)
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
			return "", fmt.Errorf("Overlay.setShowWindowControlsOverlay 请求超时")
		}
	}
}

/*

// ==================== Overlay.setShowWindowControlsOverlay 使用示例 ====================
func ExampleCDPOverlaySetShowWindowControlsOverlay() {
	// ========== 示例1：显示窗口控件区域覆盖层 ==========
	resp, err := CDPOverlaySetShowWindowControlsOverlay(true)
	if err != nil {
		log.Fatalf("开启窗口控件区域高亮失败: %v", err)
	}
	log.Printf("已显示窗口控件区域覆盖层: %s", resp)

	// ========== 示例2：关闭窗口控件区域覆盖层 ==========
	resp, err := CDPOverlaySetShowWindowControlsOverlay(false)
	if err != nil {
		log.Fatalf("关闭窗口控件区域高亮失败: %v", err)
	}
	log.Println("已关闭窗口控件区域覆盖层")

	// ========== 示例3：完整 PWA 标题栏调试流程 ==========
	// 1. 启用覆盖层
	CDPOverlayEnable()
	defer CDPOverlayDisable()

	// 2. 开启窗口控件区域可视化
	_, err := CDPOverlaySetShowWindowControlsOverlay(true)
	if err != nil {
		log.Fatalf("启动窗口控件调试失败: %v", err)
	}
	log.Println("已可视化窗口按钮区域，避免内容遮挡")

	// 3. 调试完成关闭
	// CDPOverlaySetShowWindowControlsOverlay(false)

	// ========== 示例4：PWA 安装适配测试用例 ==========
	func TestPWAWindowControls(t *testing.T) {
		CDPOverlayEnable()
		defer CDPOverlayDisable()

		// 开启控件区域检测
		CDPOverlaySetShowWindowControlsOverlay(true)
		t.Log("窗口控件区域已高亮，验证PWA标题栏适配")

		// 测试结束关闭
		defer CDPOverlaySetShowWindowControlsOverlay(false)
	}
}

*/
