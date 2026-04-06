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

// -----------------------------------------------  LayerTree.compositingReasons  -----------------------------------------------
// === 应用场景 ===
// 1. 页面性能调试: 分析页面卡顿、掉帧的图层合成原因
// 2. 渲染优化: 定位触发多余图层合成的DOM元素与样式
// 3. 自动化性能检测: 自动化测试中采集页面图层合成原因
// 4. 前端性能审计: 生成页面渲染性能报告
// 5. 复杂页面分析: 排查大型单页应用的图层爆炸问题
// 6. 动画性能优化: 分析CSS/JS动画触发图层合成的根源

// CDPLayerTreeCompositingReasons 获取指定图层的合成原因
// layerId: 图层ID（从LayerTree.layers事件中获取）
func CDPLayerTreeCompositingReasons(layerId string) (string, error) {
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
		"method": "LayerTree.compositingReasons",
		"params": {
			"layerId": "%s"
		}
	}`, reqID, layerId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.compositingReasons 请求失败: %w", err)
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
			return "", fmt.Errorf("LayerTree.compositingReasons 请求超时")
		}
	}
}

/*

// ==================== LayerTree.compositingReasons 使用示例 ====================
func ExampleCDPLayerTreeCompositingReasons() {
	// 1. 必备前提：已启用LayerTree并获取到有效图层ID
	// 通常先调用 LayerTree.enable 启用图层树，再从事件中获取 layerId
	layerID := "layer-123456" // 实际为LayerTree事件返回的图层ID

	// ========== 示例1：基础查询图层合成原因 ==========
	resp, err := CDPLayerTreeCompositingReasons(layerID)
	if err != nil {
		log.Fatalf("获取图层合成原因失败: %v", err)
	}
	log.Printf("图层合成原因查询成功: %s", resp)

	// ========== 示例2：性能调试完整流程 ==========
	// 步骤1：启用图层树
	CDPLayerTreeEnable()

	// 步骤2：获取页面图层列表（从事件/方法中获取layerId）
	// layerId := getPageLayerId()

	// 步骤3：查询该图层为什么会被合成
	reasons, _ := CDPLayerTreeCompositingReasons(layerId)
	log.Println("图层合成原因：", reasons)

	// 步骤4：关闭图层树
	CDPLayerTreeDisable()

}

*/

// -----------------------------------------------  LayerTree.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 资源释放: 停止图层树监听，释放浏览器渲染层占用的内存资源
// 2. 测试环境清理: 性能测试完成后关闭图层树采集，恢复浏览器默认状态
// 3. 流程收尾: 页面渲染分析完成后关闭监听，避免无效事件推送
// 4. 内存泄漏防护: 自动化任务结束后强制关闭，防止持续事件消耗内存
// 5. 调试流程结束: 图层渲染调试完成后关闭功能
// 6. 多任务切换: 切换不同调试/测试任务时关闭上一个任务的图层树功能

// CDPLayerTreeDisable 关闭图层树功能
func CDPLayerTreeDisable() (string, error) {
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
		"method": "LayerTree.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("LayerTree.disable 请求超时")
		}
	}
}

/*

// ==================== LayerTree.disable 使用示例 ====================
func ExampleCDPLayerTreeDisable() {
	// ========== 示例1：基础关闭图层树 ==========
	resp, err := CDPLayerTreeDisable()
	if err != nil {
		log.Fatalf("关闭图层树失败: %v", err)
	}
	log.Printf("关闭图层树成功，响应: %s", resp)

	// ========== 示例2：性能调试完整标准流程 ==========
	/*
	// 步骤1：启用图层树
	CDPLayerTreeEnable()

	// 步骤2：执行图层分析、获取合成原因等操作
	layerID := "layer-123"
	CDPLayerTreeCompositingReasons(layerID)

	// 步骤3：分析完成，关闭图层树释放资源
	closeResp, closeErr := CDPLayerTreeDisable()
	if closeErr != nil {
		log.Printf("关闭失败: %v", closeErr)
	} else {
		log.Println("图层树已关闭，资源已释放")
	}

	// ========== 示例3：自动化测试结束强制清理 ==========
	/*
	// 测试用例执行完毕，无论成功失败都关闭
	defer func() {
		_, _ = CDPLayerTreeDisable()
		log.Println("自动化测试完成，已关闭LayerTree")
	}()
	// 执行测试逻辑...
}

*/

// -----------------------------------------------  LayerTree.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 页面渲染调试：开启图层树监听，获取页面图层结构与渲染状态
// 2. 性能分析前置：调试页面卡顿、掉帧前必须先启用图层树功能
// 3. 自动化性能测试：测试开始时启用，采集页面图层渲染数据
// 4. 前端渲染优化：分析图层合成、图层爆炸问题前启用
// 5. 实时图层监控：持续接收图层更新事件，实时追踪页面渲染变化
// 6. 动画渲染分析：启用后追踪CSS/JS动画的图层创建与合成

// CDPLayerTreeEnable 启用图层树功能
func CDPLayerTreeEnable() (string, error) {
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
		"method": "LayerTree.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("LayerTree.enable 请求超时")
		}
	}
}

/*

// ==================== LayerTree.enable 使用示例 ====================
func ExampleCDPLayerTreeEnable() {
	// ========== 示例1：基础启用图层树 ==========
	resp, err := CDPLayerTreeEnable()
	if err != nil {
		log.Fatalf("启用图层树失败: %v", err)
	}
	log.Printf("启用图层树成功，响应: %s", resp)

	// ========== 示例2：性能调试标准完整流程（启用→使用→关闭） ==========
	// 步骤1：启用图层树（必须第一步）
	_, err := CDPLayerTreeEnable()
	if err != nil {
		log.Fatalf("启用失败: %v", err)
	}
	log.Println("已启用LayerTree，开始接收图层事件")

	// 步骤2：执行图层相关操作（查询图层、合成原因等）
	layerID := "layer-123456"
	reasonResp, _ := CDPLayerTreeCompositingReasons(layerID)
	log.Println("图层合成原因：", reasonResp)

	// 步骤3：完成后关闭，释放资源
	defer CDPLayerTreeDisable()
	log.Println("性能调试完成")


	// ========== 示例3：自动化测试前置初始化 ==========
	// 测试用例开始前强制启用
	func TestPagePerformance(t *testing.T) {
		// 前置：启用图层树
		_, err := CDPLayerTreeEnable()
		if err != nil {
			t.Fatalf("测试前置启用LayerTree失败: %v", err)
		}
		// 测试结束自动关闭
		defer CDPLayerTreeDisable()

		// 执行性能测试...
	}
}

*/

// -----------------------------------------------  LayerTree.loadSnapshot  -----------------------------------------------
// === 应用场景 ===
// 1. 离线渲染分析：保存页面图层后，离线加载快照进行性能调试
// 2. 自动化测试存档：测试时保存图层快照，复现问题时加载快照
// 3. 渲染问题复现：线上页面卡顿/异常，加载图层快照复现场景
// 4. 性能报告生成：加载历史图层快照，生成渲染优化报告
// 5. 多版本对比：加载不同版本页面的图层快照，对比渲染差异
// 6. 教学演示：加载预生成图层快照，演示页面渲染原理

// CDPLayerTreeLoadSnapshot 从图层树数据加载渲染快照
// layerTreeData: 完整的图层树JSON数据（从LayerTree.layers事件获取）
func CDPLayerTreeLoadSnapshot(layerTreeData string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 转义JSON数据，避免格式化冲突
	escapedData := strings.ReplaceAll(layerTreeData, `"`, `\"`)

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "LayerTree.loadSnapshot",
		"params": {
			"layerTree": "%s"
		}
	}`, reqID, escapedData)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.loadSnapshot 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
	timeout := 8 * time.Second // 快照加载稍慢，延长超时
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
			return "", fmt.Errorf("LayerTree.loadSnapshot 请求超时")
		}
	}
}

/*

// ==================== LayerTree.loadSnapshot 使用示例 ====================
func ExampleCDPLayerTreeLoadSnapshot() {
	// 1. 前提：已获取合法的图层树数据（来自LayerTree.layers事件）
	// 模拟真实图层树JSON数据（实际为事件返回的完整结构）
	layerTree := `{
		"layers": [{"layerId":"layer-1","width":1920,"height":1080}],
		"timestamp":1744000000
	}`

	// ========== 示例1：基础加载图层快照 ==========
	resp, err := CDPLayerTreeLoadSnapshot(layerTree)
	if err != nil {
		log.Fatalf("加载图层快照失败: %v", err)
	}
	log.Printf("加载图层快照成功: %s", resp)

	// ========== 示例2：离线调试完整流程 ==========
	// 步骤1：启用图层树
	CDPLayerTreeEnable()
	defer CDPLayerTreeDisable()

	// 步骤2：读取本地保存的图层数据
	savedLayerData, _ := os.ReadFile("layer_snapshot.json")

	// 步骤3：加载快照进行离线分析
	loadResp, err := CDPLayerTreeLoadSnapshot(string(savedLayerData))
	if err != nil {
		log.Fatalf("离线快照加载失败: %v", err)
	}
	log.Println("离线快照加载完成，可开始分析渲染问题")

	// ========== 示例3：自动化问题复现 ==========
	// 复现线上卡顿问题：加载上报的图层快照
	onlineLayerData := getReportLayerData()
	_, err := CDPLayerTreeLoadSnapshot(onlineLayerData)
	if err == nil {
		log.Println("线上渲染问题已成功复现")
	}
}

*/

// -----------------------------------------------  LayerTree.makeSnapshot  -----------------------------------------------
// === 应用场景 ===
// 1. 图层截图留存：对页面指定渲染图层进行截图保存，复现问题
// 2. 性能调试可视化：直观查看卡顿/异常图层的渲染效果
// 3. 自动化测试截图：测试用例中截取关键图层，生成测试报告
// 4. 页面渲染存档：保存页面核心图层快照，用于离线对比分析
// 5. 前端优化验证：优化后截图对比，验证渲染效果改善
// 6. 异常监控上报：页面图层渲染异常时，自动截图上报

// CDPLayerTreeMakeSnapshot 为指定图层创建快照
// layerId: 要截图的图层ID
// format: 图片格式，可选 jpeg/png
// quality: 图片质量 0-100，仅jpeg生效
func CDPLayerTreeMakeSnapshot(layerId string, format string, quality int) (string, error) {
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
		"method": "LayerTree.makeSnapshot",
		"params": {
			"layerId": "%s",
			"format": "%s",
			"quality": %d
		}
	}`, reqID, layerId, format, quality)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.makeSnapshot 请求失败: %w", err)
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
			return "", fmt.Errorf("LayerTree.makeSnapshot 请求超时")
		}
	}
}

/*

// ==================== LayerTree.makeSnapshot 使用示例 ====================
func ExampleCDPLayerTreeMakeSnapshot() {
	// 前提：已启用LayerTree并获取有效图层ID
	layerID := "layer-123456"

	// ========== 示例1：创建JPEG高质量图层快照 ==========
	resp, err := CDPLayerTreeMakeSnapshot(layerID, "jpeg", 90)
	if err != nil {
		log.Fatalf("创建图层快照失败: %v", err)
	}
	log.Printf("快照创建成功，响应: %s", resp)

	// ========== 示例2：创建PNG无损图层快照 ==========
	resp2, err2 := CDPLayerTreeMakeSnapshot(layerID, "png", 100)
	if err2 != nil {
		log.Fatalf("创建PNG快照失败: %v", err2)
	}
	log.Printf("PNG快照创建成功: %s", resp2)

	// ========== 示例3：性能调试完整流程 ==========
	// 1. 启用图层树
	CDPLayerTreeEnable()
	defer CDPLayerTreeDisable()

	// 2. 对异常图层创建快照
	layerId := "abnormal-layer-01"
	snapshotResp, _ := CDPLayerTreeMakeSnapshot(layerId, "jpeg", 85)
	log.Println("异常图层快照已创建：", snapshotResp)

	// 3. 后续可通过快照ID获取图片数据进行保存/上报
}

*/

// -----------------------------------------------  LayerTree.profileSnapshot  -----------------------------------------------
// === 应用场景 ===
// 1. 快照性能分析：分析图层快照的渲染耗时、绘制性能瓶颈
// 2. 离线性能调试：加载历史图层快照，离线分析渲染性能问题
// 3. 性能对比测试：对比不同版本/优化前后的图层渲染性能
// 4. 自动化性能审计：自动化生成图层快照性能报告
// 5. 卡顿根因定位：通过快照性能数据定位页面掉帧、卡顿原因
// 6. 渲染优化验证：优化后采集性能数据，验证优化效果

// CDPLayerTreeProfileSnapshot 获取图层快照的性能分析数据
// snapshotId: 图层快照ID（从LayerTree.makeSnapshot获取）
// minRepeatCount: 最小重复采样次数（默认≥5）
// minDurationSeconds: 最小采样时长（默认≥0.5）
func CDPLayerTreeProfileSnapshot(snapshotId string, minRepeatCount int, minDurationSeconds float64) (string, error) {
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
		"method": "LayerTree.profileSnapshot",
		"params": {
			"snapshotId": "%s",
			"minRepeatCount": %d,
			"minDurationSeconds": %f
		}
	}`, reqID, snapshotId, minRepeatCount, minDurationSeconds)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.profileSnapshot 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（性能采样需要更长时间）
	timeout := 10 * time.Second
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
			return "", fmt.Errorf("LayerTree.profileSnapshot 请求超时")
		}
	}
}

/*

// ==================== LayerTree.profileSnapshot 使用示例 ====================
func ExampleCDPLayerTreeProfileSnapshot() {
	// 前提条件：
	// 1. 已启用 LayerTree.enable
	// 2. 已通过 makeSnapshot 创建快照并获取 snapshotId
	snapshotID := "snapshot-123456"

	// ========== 示例1：标准性能采样（推荐参数） ==========
	resp, err := CDPLayerTreeProfileSnapshot(snapshotID, 5, 0.8)
	if err != nil {
		log.Fatalf("图层快照性能分析失败: %v", err)
	}
	log.Printf("性能分析完成，结果: %s", resp)

	// ========== 示例2：高精度采样（更长时间、更多次数） ==========
	resp2, err2 := CDPLayerTreeProfileSnapshot(snapshotID, 10, 1.5)
	if err2 != nil {
		log.Fatalf("高精度性能分析失败: %v", err2)
	}
	log.Printf("高精度性能数据: %s", resp2)

	// ========== 示例3：完整性能调试流程 ==========
	// 1. 启用图层树
	CDPLayerTreeEnable()
	defer CDPLayerTreeDisable()

	// 2. 创建图层快照
	layerID := "layer-optimize-01"
	makeResp, _ := CDPLayerTreeMakeSnapshot(layerID, "jpeg", 90)
	// 解析获取 snapshotId

	// 3. 对快照进行性能采样
	profileResp, _ := CDPLayerTreeProfileSnapshot(snapshotID, 5, 0.8)
	log.Println("图层渲染性能数据：", profileResp)
}

*/

// -----------------------------------------------  LayerTree.releaseSnapshot  -----------------------------------------------
// === 应用场景 ===
// 1. 内存资源释放: 用完图层快照后立即释放，避免浏览器内存持续占用
// 2. 快照生命周期管理: 自动化测试中创建快照后规范清理
// 3. 批量快照清理: 批量生成多个图层快照后统一释放
// 4. 性能调试收尾: 渲染分析完成后清理快照资源
// 5. 错误恢复: 快照使用异常时强制释放防止句柄泄漏
// 6. 长时间运行程序: 守护进程/自动化脚本中防止内存溢出

// CDPLayerTreeReleaseSnapshot 释放指定图层快照资源
// snapshotId: 通过LayerTree.makeSnapshot / LayerTree.loadSnapshot 获取的快照ID
func CDPLayerTreeReleaseSnapshot(snapshotId string) (string, error) {
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
		"method": "LayerTree.releaseSnapshot",
		"params": {
			"snapshotId": "%s"
		}
	}`, reqID, snapshotId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.releaseSnapshot 请求失败: %w", err)
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
			return "", fmt.Errorf("LayerTree.releaseSnapshot 请求超时")
		}
	}
}

/*

// ==================== LayerTree.releaseSnapshot 使用示例 ====================
func ExampleCDPLayerTreeReleaseSnapshot() {
	// 前提：已创建图层快照并获取 snapshotId
	snapshotID := "snapshot-123456"

	// ========== 示例1：基础释放快照 ==========
	resp, err := CDPLayerTreeReleaseSnapshot(snapshotID)
	if err != nil {
		log.Fatalf("释放图层快照失败: %v", err)
	}
	log.Printf("图层快照释放成功，响应: %s", resp)

	// ========== 示例2：标准完整流程（创建→使用→释放） ==========
	// 1. 启用图层树
	CDPLayerTreeEnable()
	defer CDPLayerTreeDisable()

	// 2. 创建快照
	layerID := "layer-001"
	_, _ = CDPLayerTreeMakeSnapshot(layerID, "jpeg", 90)

	// 3. 性能分析
	_, _ = CDPLayerTreeProfileSnapshot(snapshotID, 5, 0.5)

	// 4. 用完立即释放（核心步骤）
	releaseResp, releaseErr := CDPLayerTreeReleaseSnapshot(snapshotID)
	if releaseErr != nil {
		log.Printf("释放失败: %v", releaseErr)
	} else {
		log.Println("快照已释放，内存已回收")
	}

	// ========== 示例3：defer 确保一定释放（自动化推荐） ==========
	// 函数退出时自动释放，防止遗漏
	func snapshotAnalysis() {
		snapshotID := "snapshot-auto-release"
		defer CDPLayerTreeReleaseSnapshot(snapshotID)

		// 执行快照分析逻辑...
		log.Println("开始分析，结束后自动释放快照")
	}
}

*/

// -----------------------------------------------  LayerTree.replaySnapshot  -----------------------------------------------
// === 应用场景 ===
// 1. 问题复现：回放历史图层快照，精准复现页面渲染异常、卡顿、错位问题
// 2. 离线调试：无需访问原页面，直接回放快照进行渲染分析
// 3. 自动化测试：回放快照验证页面渲染是否符合预期
// 4. 性能对比：回放优化前后的快照，对比渲染效果与性能
// 5. 演示教学：回放快照展示页面图层渲染原理与问题场景
// 6. 线上故障还原：回放用户上报的快照，还原线上真实渲染状态

// CDPLayerTreeReplaySnapshot 回放图层快照
// snapshotId: 图层快照ID（从makeSnapshot/loadSnapshot获取）
// fromStep: 可选，开始回放的步骤（从第几帧开始）
// toStep: 可选，结束回放的步骤（到第几帧结束）
// scale: 可选，回放渲染缩放比例（默认1.0）
func CDPLayerTreeReplaySnapshot(snapshotId string, fromStep int, toStep int, scale float64) (string, error) {
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
		"method": "LayerTree.replaySnapshot",
		"params": {
			"snapshotId": "%s",
			"fromStep": %d,
			"toStep": %d,
			"scale": %f
		}
	}`, reqID, snapshotId, fromStep, toStep, scale)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.replaySnapshot 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
	timeout := 8 * time.Second
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
			return "", fmt.Errorf("LayerTree.replaySnapshot 请求超时")
		}
	}
}

/*

// ==================== LayerTree.replaySnapshot 使用示例 ====================
func ExampleCDPLayerTreeReplaySnapshot() {
	// 前提：已创建/加载图层快照并获取有效 snapshotId
	snapshotID := "snapshot-123456"

	// ========== 示例1：完整回放快照（默认参数：从头至尾，1倍缩放） ==========
	resp, err := CDPLayerTreeReplaySnapshot(snapshotID, 0, -1, 1.0)
	if err != nil {
		log.Fatalf("图层快照回放失败: %v", err)
	}
	log.Printf("图层快照完整回放成功: %s", resp)

	// ========== 示例2：指定范围片段回放（从第5帧到第20帧，1.5倍缩放） ==========
	resp2, err2 := CDPLayerTreeReplaySnapshot(snapshotID, 5, 20, 1.5)
	if err2 != nil {
		log.Fatalf("片段回放失败: %v", err2)
	}
	log.Printf("快照片段回放成功: %s", resp2)

	// ========== 示例3：离线问题复现完整流程 ==========
	// 1. 启用图层树
	CDPLayerTreeEnable()
	defer CDPLayerTreeDisable()

	// 2. 加载/创建快照
	// loadResp, _ := CDPLayerTreeLoadSnapshot(savedLayerData)

	// 3. 回放快照复现线上渲染问题
	replayResp, _ := CDPLayerTreeReplaySnapshot(snapshotID, 0, -1, 1.0)
	log.Println("线上渲染问题已通过快照回放复现")

	// 4. 分析完成后释放快照
	defer CDPLayerTreeReleaseSnapshot(snapshotID)
}

*/

// -----------------------------------------------  LayerTree.snapshotCommandLog  -----------------------------------------------
// === 应用场景 ===
// 1. 深度渲染调试: 获取图层快照的绘制指令，分析渲染执行流程
// 2. 渲染异常定位: 排查页面错位、黑屏、闪烁的底层绘制原因
// 3. 性能瓶颈分析: 查看绘制指令数量与复杂度，定位慢渲染根源
// 4. 自动化渲染审计: 采集绘制日志生成页面渲染合规报告
// 5. 前端渲染优化: 根据绘制指令优化DOM结构与CSS样式
// 6. 离线问题复盘: 加载快照日志，离线分析线上渲染故障

// CDPLayerTreeSnapshotCommandLog 获取图层快照的绘制命令日志
// snapshotId: 通过LayerTree.makeSnapshot获取的图层快照ID
func CDPLayerTreeSnapshotCommandLog(snapshotId string) (string, error) {
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
		"method": "LayerTree.snapshotCommandLog",
		"params": {
			"snapshotId": "%s"
		}
	}`, reqID, snapshotId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 LayerTree.snapshotCommandLog 请求失败: %w", err)
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
			return "", fmt.Errorf("LayerTree.snapshotCommandLog 请求超时")
		}
	}
}

/*

// ==================== LayerTree.snapshotCommandLog 使用示例 ====================
func ExampleCDPLayerTreeSnapshotCommandLog() {
	// 前提：已启用LayerTree，并通过makeSnapshot创建快照获取snapshotId
	snapshotID := "snapshot-123456"

	// ========== 示例1：基础获取快照绘制命令日志 ==========
	resp, err := CDPLayerTreeSnapshotCommandLog(snapshotID)
	if err != nil {
		log.Fatalf("获取快照绘制日志失败: %v", err)
	}
	log.Printf("绘制命令日志获取成功: %s", resp)

	// ========== 示例2：渲染异常深度分析完整流程 ==========
	// 1. 启用图层树
	CDPLayerTreeEnable()
	defer CDPLayerTreeDisable()

	// 2. 对异常图层创建快照
	layerID := "layer-abnormal-render"
	_, _ = CDPLayerTreeMakeSnapshot(layerID, "png", 100)

	// 3. 获取绘制指令日志，定位渲染问题
	logResp, _ := CDPLayerTreeSnapshotCommandLog(snapshotID)
	log.Println("图层绘制指令日志：", logResp)

	// 4. 分析完成释放快照资源
	defer CDPLayerTreeReleaseSnapshot(snapshotID)

	// ========== 示例3：自动化性能采集 ==========
	// 自动化采集页面核心图层绘制日志，生成性能报告
	func captureRenderLog(snapshotId string) {
		cmdLog, err := CDPLayerTreeSnapshotCommandLog(snapshotId)
		if err == nil {
			saveRenderReport(cmdLog) // 保存日志到文件/数据库
		}
	}
}


*/
