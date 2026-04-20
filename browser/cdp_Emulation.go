package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Emulation.clearDeviceMetricsOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 设备模拟恢复: 恢复设备指标为原始状态
// 2. 测试环境清理: 在测试完成后清理设备模拟配置
// 3. 设备切换: 切换不同的设备模拟配置
// 4. 错误恢复: 在模拟出现问题时恢复默认状态
// 5. 自动化测试: 在自动化测试流程中清理设备状态
// 6. 调试辅助: 在调试设备相关问题时恢复原始状态

// CDPEmulationClearDeviceMetricsOverride 清除设备指标覆盖
func CDPEmulationClearDeviceMetricsOverride() (string, error) {
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
		"method": "Emulation.clearDeviceMetricsOverride"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearDeviceMetricsOverride 请求失败: %w", err)
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
			return "", fmt.Errorf("clearDeviceMetricsOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.clearGeolocationOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 地理位置模拟恢复: 恢复原始地理位置获取
// 2. 测试环境清理: 在位置相关的测试完成后清理模拟
// 3. 地理位置切换: 切换不同的地理位置模拟配置
// 4. 错误恢复: 在地理位置模拟出现问题时恢复默认
// 5. 自动化测试: 在自动化测试流程中清理地理位置状态
// 6. 隐私保护: 清理地理位置模拟以保护用户隐私

// CDPEmulationClearGeolocationOverride 清除地理位置覆盖
func CDPEmulationClearGeolocationOverride() (string, error) {
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
		"method": "Emulation.clearGeolocationOverride"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearGeolocationOverride 请求失败: %w", err)
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
			return "", fmt.Errorf("clearGeolocationOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.clearIdleOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 空闲状态模拟恢复: 恢复原始空闲状态检测
// 2. 测试环境清理: 在空闲状态相关的测试完成后清理模拟
// 3. 状态切换: 切换不同的空闲状态模拟配置
// 4. 错误恢复: 在空闲状态模拟出现问题时恢复默认
// 5. 自动化测试: 在自动化测试流程中清理空闲状态
// 6. 性能测试: 清理空闲状态覆盖以确保准确性能测量

// CDPEmulationClearIdleOverride 清除空闲状态覆盖
func CDPEmulationClearIdleOverride() (string, error) {
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
		"method": "Emulation.clearIdleOverride"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearIdleOverride 请求失败: %w", err)
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
			return "", fmt.Errorf("clearIdleOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setCPUThrottlingRate  -----------------------------------------------
// === 应用场景 ===
// 1. 性能测试: 模拟低性能设备的渲染行为
// 2. 性能基准: 在不同CPU限制下测试页面性能
// 3. 资源限制测试: 测试页面在资源受限环境下的表现
// 4. 性能回归测试: 确保性能优化在不同硬件上有效
// 5. 加载速度测试: 测试慢速CPU下的页面加载性能
// 6. 动画性能测试: 测试CPU限制下的动画和滚动性能

// CDPEmulationSetCPUThrottlingRate 设置CPU限制率
// 参数:
//   - rate: 限制率作为减速因子（1表示不限制，2表示2倍减速，等等）
func CDPEmulationSetCPUThrottlingRate(rate float64) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 参数验证
	if rate < 0.1 {
		return "", fmt.Errorf("CPU限制率不能小于0.1: %f", rate)
	}
	if rate > 100 {
		return "", fmt.Errorf("CPU限制率不能大于100: %f", rate)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setCPUThrottlingRate",
		"params": {
			"rate": %f
		}
	}`, reqID, rate)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setCPUThrottlingRate 请求失败: %w", err)
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
			return "", fmt.Errorf("setCPUThrottlingRate 请求超时")
		}
	}
}

// 示例1: CPU性能模拟和基准测试
func exampleCPUPerformanceSimulationAndBenchmark() {
	// === 应用场景描述 ===
	// 场景: CPU性能模拟和基准测试
	// 用途: 在不同CPU性能限制下测试页面性能表现
	// 优势: 模拟不同性能等级的CPU设备
	// 典型工作流: 设置CPU限制率 -> 执行性能测试 -> 收集性能指标 -> 清理设置

	log.Println("CPU性能模拟和基准测试示例...")

	// 定义CPU性能测试场景
	cpuPerformanceScenarios := []struct {
		name              string
		description       string
		throttlingRate    float64
		deviceCategory    string
		performanceImpact string
		testCases         []string
		expectedMetrics   map[string]interface{}
	}{
		{
			name:              "高性能桌面",
			description:       "模拟高性能桌面CPU，无限制",
			throttlingRate:    1.0,
			deviceCategory:    "高端桌面电脑",
			performanceImpact: "无性能限制，最大性能",
			testCases: []string{
				"页面加载性能",
				"JavaScript执行速度",
				"CSS动画帧率",
				"DOM操作性能",
			},
			expectedMetrics: map[string]interface{}{
				"frameRate": "60+ fps",
				"loadTime":  "<2秒",
				"jsTime":    "<500ms",
				"cpuUsage":  "高",
			},
		},
		{
			name:              "中端移动设备",
			description:       "模拟中端移动设备CPU性能",
			throttlingRate:    2.5,
			deviceCategory:    "中端智能手机",
			performanceImpact: "中等性能限制，模拟移动设备",
			testCases: []string{
				"移动端页面加载",
				"触摸响应速度",
				"移动端动画性能",
				"节电模式性能",
			},
			expectedMetrics: map[string]interface{}{
				"frameRate": "30-60 fps",
				"loadTime":  "2-5秒",
				"jsTime":    "500-1500ms",
				"cpuUsage":  "中高",
			},
		},
		{
			name:              "低端设备",
			description:       "模拟低端设备CPU性能",
			throttlingRate:    5.0,
			deviceCategory:    "低端智能手机/平板",
			performanceImpact: "较高性能限制，模拟低端设备",
			testCases: []string{
				"低端设备页面渲染",
				"JavaScript密集型任务",
				"复杂动画性能",
				"内存受限性能",
			},
			expectedMetrics: map[string]interface{}{
				"frameRate": "15-30 fps",
				"loadTime":  "5-10秒",
				"jsTime":    "1500-3000ms",
				"cpuUsage":  "中",
			},
		},
		{
			name:              "极低性能设备",
			description:       "模拟性能极低的设备",
			throttlingRate:    10.0,
			deviceCategory:    "老旧设备/物联网设备",
			performanceImpact: "极高性能限制，模拟老旧设备",
			testCases: []string{
				"极限性能测试",
				"降级用户体验检测",
				"渐进增强验证",
				"兼容性测试",
			},
			expectedMetrics: map[string]interface{}{
				"frameRate": "<15 fps",
				"loadTime":  ">10秒",
				"jsTime":    ">3000ms",
				"cpuUsage":  "低",
			},
		},
		{
			name:              "网络节流设备",
			description:       "模拟网络受限设备的CPU性能",
			throttlingRate:    3.0,
			deviceCategory:    "网络受限移动设备",
			performanceImpact: "网络相关CPU限制",
			testCases: []string{
				"网络加载性能",
				"解析阻塞测试",
				"资源加载顺序",
				"懒加载性能",
			},
			expectedMetrics: map[string]interface{}{
				"frameRate": "20-40 fps",
				"loadTime":  "3-7秒",
				"jsTime":    "1000-2000ms",
				"cpuUsage":  "中等",
			},
		},
		{
			name:              "节电模式",
			description:       "模拟节电模式下的CPU性能",
			throttlingRate:    4.0,
			deviceCategory:    "节电模式设备",
			performanceImpact: "节电模式性能限制",
			testCases: []string{
				"节电模式页面渲染",
				"后台任务处理",
				"动画优化测试",
				"节能策略验证",
			},
			expectedMetrics: map[string]interface{}{
				"frameRate": "20-30 fps",
				"loadTime":  "4-8秒",
				"jsTime":    "1200-2500ms",
				"cpuUsage":  "中低",
			},
		},
		{
			name:              "CPU瓶颈测试",
			description:       "测试CPU密集型应用瓶颈",
			throttlingRate:    8.0,
			deviceCategory:    "CPU密集型场景",
			performanceImpact: "极端CPU限制，测试瓶颈",
			testCases: []string{
				"复杂计算任务",
				"大数据处理性能",
				"实时分析能力",
				"并发处理测试",
			},
			expectedMetrics: map[string]interface{}{
				"frameRate": "<10 fps",
				"loadTime":  ">8秒",
				"jsTime":    ">2000ms",
				"cpuUsage":  "极高",
			},
		},
	}

	// 测试每种CPU性能场景
	for i, scenario := range cpuPerformanceScenarios {
		log.Printf("\n=== CPU性能测试场景 %d/%d: %s ===", i+1, len(cpuPerformanceScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("CPU限制率: %.1fx (1=无限制)", scenario.throttlingRate)
		log.Printf("设备类别: %s", scenario.deviceCategory)
		log.Printf("性能影响: %s", scenario.performanceImpact)
		log.Printf("测试用例:")
		for _, testCase := range scenario.testCases {
			log.Printf("  - %s", testCase)
		}
		log.Printf("预期指标:")
		for key, value := range scenario.expectedMetrics {
			log.Printf("  %s: %v", key, value)
		}

		// 检查当前CPU限制
		currentRate := getCurrentCPUThrottlingRate()
		log.Printf("当前CPU限制率: %.1fx", currentRate)

		// 执行CPU限制率设置
		log.Printf("设置CPU限制率 %.1fx...", scenario.throttlingRate)
		response, err := CDPEmulationSetCPUThrottlingRate(scenario.throttlingRate)
		if err != nil {
			log.Printf("设置CPU限制率失败: %v", err)
			continue
		}

		log.Printf("✅ CPU限制率设置成功: %s", response)

		// 模拟性能测试
		performanceMetrics := simulatePerformanceTesting(scenario.throttlingRate, scenario.testCases)

		// 分析测试结果
		analyzePerformanceResults(scenario, performanceMetrics)

		// 提供性能优化建议
		providePerformanceRecommendations(scenario, performanceMetrics)

		// 恢复CPU性能
		log.Printf("恢复CPU性能(设置限制率1.0)...")
		clearResponse, err := CDPEmulationSetCPUThrottlingRate(1.0)
		if err != nil {
			log.Printf("恢复CPU性能失败: %v", err)
		} else {
			log.Printf("✅ CPU性能已恢复: %s", clearResponse)
		}

		// 短暂延迟
		time.Sleep(200 * time.Millisecond)
	}
}

func getCurrentCPUThrottlingRate() float64 {
	// 这里应该实现获取当前CPU限制率的逻辑
	// 由于没有对应的CDP方法，我们模拟返回
	return 1.0
}

func simulatePerformanceTesting(throttlingRate float64, testCases []string) map[string]interface{} {
	metrics := make(map[string]interface{})

	// 模拟计算各种性能指标
	baseLoadTime := 1.5
	baseFrameRate := 60.0
	baseJSTime := 300.0

	// 根据限制率计算性能指标
	loadTimeMultiplier := 1.0 + (throttlingRate-1)*0.3
	frameRateMultiplier := 1.0 / (1.0 + (throttlingRate-1)*0.1)
	jsTimeMultiplier := 1.0 + (throttlingRate-1)*0.5

	metrics["throttlingRate"] = throttlingRate
	metrics["loadTime"] = fmt.Sprintf("%.2f秒", baseLoadTime*loadTimeMultiplier)
	metrics["frameRate"] = fmt.Sprintf("%.1f fps", baseFrameRate*frameRateMultiplier)
	metrics["jsExecutionTime"] = fmt.Sprintf("%.0fms", baseJSTime*jsTimeMultiplier)
	metrics["cpuUtilization"] = fmt.Sprintf("%.0f%%", math.Min(100, 20+throttlingRate*5))

	// 根据测试用例添加额外指标
	additionalMetrics := make(map[string]string)

	for _, testCase := range testCases {
		switch {
		case strings.Contains(testCase, "页面加载"):
			additionalMetrics["domReadyTime"] = fmt.Sprintf("%.2f秒", baseLoadTime*loadTimeMultiplier*0.7)
			additionalMetrics["firstContentfulPaint"] = fmt.Sprintf("%.2f秒", baseLoadTime*loadTimeMultiplier*0.5)

		case strings.Contains(testCase, "JavaScript"):
			additionalMetrics["jsHeapUsed"] = fmt.Sprintf("%.1fMB", 20.0+throttlingRate*2)
			additionalMetrics["jsExecutionScore"] = fmt.Sprintf("%.0f/100", 100.0/(1.0+throttlingRate*0.1))

		case strings.Contains(testCase, "CSS动画") || strings.Contains(testCase, "动画性能"):
			additionalMetrics["animationFrameTime"] = fmt.Sprintf("%.0fms", 16.0*throttlingRate)
			additionalMetrics["animationSmoothness"] = fmt.Sprintf("%.0f%%", 100.0/(1.0+throttlingRate*0.2))

		case strings.Contains(testCase, "DOM操作"):
			additionalMetrics["domManipulationTime"] = fmt.Sprintf("%.0fms", 50.0*throttlingRate)
			additionalMetrics["domMutationScore"] = fmt.Sprintf("%.0f/100", 100.0/(1.0+throttlingRate*0.15))

		case strings.Contains(testCase, "触摸响应"):
			additionalMetrics["touchLatency"] = fmt.Sprintf("%.0fms", 100.0+throttlingRate*20)
			additionalMetrics["touchResponsiveness"] = fmt.Sprintf("%.0f%%", 100.0/(1.0+throttlingRate*0.25))

		case strings.Contains(testCase, "内存受限"):
			additionalMetrics["memoryPressure"] = fmt.Sprintf("%.0f%%", 30.0+throttlingRate*5)
			additionalMetrics["gcFrequency"] = fmt.Sprintf("%.1f次/秒", 0.5+throttlingRate*0.1)

		case strings.Contains(testCase, "复杂计算"):
			additionalMetrics["computationTime"] = fmt.Sprintf("%.2f秒", throttlingRate*0.5)
			additionalMetrics["computationScore"] = fmt.Sprintf("%.0f/100", 100.0/(1.0+throttlingRate*0.3))

		case strings.Contains(testCase, "网络加载"):
			additionalMetrics["networkTimeRatio"] = fmt.Sprintf("%.2f", throttlingRate*0.7)
			additionalMetrics["parseBlockingTime"] = fmt.Sprintf("%.0fms", 200.0*throttlingRate)
		}
	}

	metrics["additionalMetrics"] = additionalMetrics

	// 计算性能评分
	performanceScore := 100.0 / (1.0 + (throttlingRate-1)*0.2)
	metrics["performanceScore"] = fmt.Sprintf("%.0f/100", performanceScore)

	if performanceScore >= 80 {
		metrics["performanceLevel"] = "优秀"
	} else if performanceScore >= 60 {
		metrics["performanceLevel"] = "良好"
	} else if performanceScore >= 40 {
		metrics["performanceLevel"] = "一般"
	} else if performanceScore >= 20 {
		metrics["performanceLevel"] = "较差"
	} else {
		metrics["performanceLevel"] = "极差"
	}

	return metrics
}

func analyzePerformanceResults(scenario struct {
	name              string
	description       string
	throttlingRate    float64
	deviceCategory    string
	performanceImpact string
	testCases         []string
	expectedMetrics   map[string]interface{}
}, actualMetrics map[string]interface{}) {
	log.Printf("性能测试结果分析:")

	// 基本指标对比
	log.Printf("  基本指标:")
	for key, expectedValue := range scenario.expectedMetrics {
		actualValue, exists := actualMetrics[key]
		if exists {
			log.Printf("    %s: 预期=%v, 实际=%v", key, expectedValue, actualValue)
		}
	}

	// 性能评分
	if score, ok := actualMetrics["performanceScore"].(string); ok {
		log.Printf("  性能评分: %s", score)
	}

	if level, ok := actualMetrics["performanceLevel"].(string); ok {
		log.Printf("  性能等级: %s", level)
	}

	// CPU利用率
	if cpuUsage, ok := actualMetrics["cpuUtilization"].(string); ok {
		log.Printf("  CPU利用率: %s", cpuUsage)
	}

	// 额外指标
	if additionalMetrics, ok := actualMetrics["additionalMetrics"].(map[string]string); ok {
		log.Printf("  详细指标:")
		for key, value := range additionalMetrics {
			log.Printf("    %s: %s", key, value)
		}
	}

	// 性能影响分析
	log.Printf("  性能影响分析:")
	throttlingRate := scenario.throttlingRate

	if throttlingRate <= 1.5 {
		log.Printf("    ⚡ 性能影响: 轻微")
		log.Printf("      适合: 性能基准测试，高端设备模拟")
	} else if throttlingRate <= 3.0 {
		log.Printf("    ⚡ 性能影响: 中等")
		log.Printf("      适合: 移动设备测试，真实用户场景")
	} else if throttlingRate <= 6.0 {
		log.Printf("    ⚡ 性能影响: 显著")
		log.Printf("      适合: 低端设备测试，性能优化验证")
	} else {
		log.Printf("    ⚡ 性能影响: 严重")
		log.Printf("      适合: 极限测试，兼容性验证")
	}

	// 设备匹配度评估
	log.Printf("  设备匹配度评估:")
	rateMatch := 1.0 / (1.0 + math.Abs(throttlingRate-expectedRateForDevice(scenario.deviceCategory)))
	log.Printf("    与设备类别'%s'的匹配度: %.0f%%",
		scenario.deviceCategory, rateMatch*100)

	if rateMatch >= 0.9 {
		log.Printf("    ✅ 高匹配度: 适合%s的性能测试", scenario.deviceCategory)
	} else if rateMatch >= 0.7 {
		log.Printf("    ⚠ 中等匹配度: 可用于%s的近似测试", scenario.deviceCategory)
	} else {
		log.Printf("    ❌ 低匹配度: 考虑调整限制率以更准确模拟%s", scenario.deviceCategory)
	}
}

func expectedRateForDevice(deviceCategory string) float64 {
	switch deviceCategory {
	case "高端桌面电脑":
		return 1.0
	case "中端智能手机":
		return 2.5
	case "低端智能手机/平板":
		return 5.0
	case "老旧设备/物联网设备":
		return 10.0
	case "网络受限移动设备":
		return 3.0
	case "节电模式设备":
		return 4.0
	case "CPU密集型场景":
		return 8.0
	default:
		return 2.0
	}
}

func providePerformanceRecommendations(scenario struct {
	name              string
	description       string
	throttlingRate    float64
	deviceCategory    string
	performanceImpact string
	testCases         []string
	expectedMetrics   map[string]interface{}
}, metrics map[string]interface{}) {
	log.Printf("性能优化建议:")

	throttlingRate := scenario.throttlingRate
	performanceLevel, _ := metrics["performanceLevel"].(string)

	// 通用建议
	log.Printf("  通用优化建议:")
	log.Printf("    - 测试不同限制率下的性能表现")
	log.Printf("    - 监控关键性能指标的变化")
	log.Printf("    - 验证降级体验的可接受性")
	log.Printf("    - 记录性能回归数据")

	// 针对特定限制率的建议
	log.Printf("  %s专用建议:", scenario.deviceCategory)
	switch {
	case throttlingRate <= 1.5:
		log.Printf("    - 重点关注首屏渲染时间")
		log.Printf("    - 优化JavaScript执行效率")
		log.Printf("    - 确保高帧率动画流畅")
		log.Printf("    - 测试复杂交互的响应速度")

	case throttlingRate <= 3.0:
		log.Printf("    - 优化移动端页面加载")
		log.Printf("    - 减少主线程阻塞时间")
		log.Printf("    - 实现适当的懒加载策略")
		log.Printf("    - 优化触摸事件响应")

	case throttlingRate <= 6.0:
		log.Printf("    - 简化JavaScript代码")
		log.Printf("    - 优化CSS选择器复杂度")
		log.Printf("    - 减少不必要的重绘重排")
		log.Printf("    - 实现渐进增强")

	default:
		log.Printf("    - 提供简化版界面")
		log.Printf("    - 显著减少JavaScript使用")
		log.Printf("    - 避免复杂CSS效果")
		log.Printf("    - 确保基本功能可用")
	}

	// 针对性能等级的建议
	log.Printf("  针对性能等级'%s'的建议:", performanceLevel)
	switch performanceLevel {
	case "优秀":
		log.Printf("    - 保持当前优化水平")
		log.Printf("    - 监控性能回归")
		log.Printf("    - 探索进一步优化空间")

	case "良好":
		log.Printf("    - 优化关键性能路径")
		log.Printf("    - 减少主线程阻塞")
		log.Printf("    - 优化资源加载")

	case "一般":
		log.Printf("    - 重点优化核心功能")
		log.Printf("    - 显著减少JS执行时间")
		log.Printf("    - 考虑使用Web Worker")

	case "较差":
		log.Printf("    - 重构关键性能代码")
		log.Printf("    - 简化页面结构")
		log.Printf("    - 实现显著性能优化")

	case "极差":
		log.Printf("    - 需要重大性能重构")
		log.Printf("    - 考虑完全重写核心模块")
		log.Printf("    - 大幅减少功能复杂度")
	}

	// 测试建议
	log.Printf("  测试策略建议:")
	log.Printf("    - 在%.1fx限制率下运行完整测试套件", throttlingRate)
	log.Printf("    - 监控内存使用和GC频率")
	log.Printf("    - 测试页面加载过程中的用户感知性能")
	log.Printf("    - 验证交互响应时间")

	// 监控建议
	log.Printf("  监控建议:")
	log.Printf("    - 监控实际用户的设备性能分布")
	log.Printf("    - 建立性能基准线")
	log.Printf("    - 设置性能告警阈值")
	log.Printf("    - 定期进行性能回归测试")
}

// -----------------------------------------------  Emulation.setDefaultBackgroundColorOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 打印预览测试: 模拟白色背景用于打印预览
// 2. 暗黑模式测试: 覆盖默认背景测试暗黑模式
// 3. 主题测试: 测试不同主题背景色的兼容性
// 4. 截图测试: 统一背景颜色以确保截图一致性
// 5. 对比度测试: 测试文本在不同背景色上的可读性
// 6. 无障碍测试: 验证背景色是否符合无障碍标准

// CDPEmulationSetDefaultBackgroundColorOverride 设置或清除默认背景颜色覆盖
// 参数:
//   - color: 可选的RGBA颜色对象，包含r,g,b,a四个字段
//     如果color为nil，则清除现有的背景颜色覆盖
func CDPEmulationSetDefaultBackgroundColorOverride(params string) (string, error) {
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
		"method": "Emulation.setDefaultBackgroundColorOverride",
		"params": %s
	}`, reqID, params)
	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setDefaultBackgroundColorOverride 请求失败: %w", err)
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

				log.Printf("✅ 默认背景颜色覆盖设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setDefaultBackgroundColorOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setDeviceMetricsOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 移动设备模拟: 模拟各种移动设备的屏幕尺寸
// 2. 响应式设计测试: 测试不同屏幕尺寸下的页面布局
// 3. 高DPI设备测试: 测试视网膜显示屏等高DPI设备
// 4. 可折叠设备测试: 模拟可折叠设备的姿态和显示
// 5. 多屏设备测试: 测试多段屏幕设备
// 6. 视口测试: 测试不同视口设置下的页面行为

// CDPEmulationSetDeviceMetricsOverride 设置设备指标覆盖
// 参数说明:
//   - width: 覆盖的宽度值（像素），0-10000000，0表示禁用
//   - height: 覆盖的高度值（像素），0-10000000，0表示禁用
//   - deviceScaleFactor: 设备比例因子，0表示禁用
//   - mobile: 是否模拟移动设备
func CDPEmulationSetDeviceMetricsOverride(params string) (string, error) {
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
		"method": "Emulation.setDeviceMetricsOverride",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setDeviceMetricsOverride 请求失败: %w", err)
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

				log.Printf("✅ 设备指标覆盖设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setDeviceMetricsOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setEmulatedMedia  -----------------------------------------------
// === 应用场景 ===
// 1. 打印样式测试: 模拟"print"媒体类型，测试打印样式
// 2. 屏幕阅读器测试: 模拟"speech"媒体类型，测试无障碍访问
// 3. 媒体特性测试: 测试不同的屏幕宽度、高度、方向、颜色模式等
// 4. 设备适配测试: 测试不同设备特性下的样式适配
// 5. 响应式设计: 测试媒体查询的响应

// CDPEmulationSetEmulatedMedia 设置模拟媒体
// 参数说明:
//   - media: 可选字符串，要模拟的媒体类型（如"screen", "print", "speech"等），空字符串表示禁用
//   - features: 可选数组，要模拟的媒体特性
func CDPEmulationSetEmulatedMedia(params string) (string, error) {
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
		"method": "Emulation.setEmulatedMedia",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setEmulatedMedia 请求失败: %w", err)
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

				log.Printf("✅ 媒体模拟设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setEmulatedMedia 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setEmulatedOSTextScale  -----------------------------------------------
// === 应用场景 ===
// 1. 无障碍访问测试: 模拟视力障碍用户使用的文本缩放
// 2. 高DPI适配: 测试在高DPI设备上的文本渲染
// 3. 老年人友好性: 测试大字号模式下的布局
// 4. 多语言支持: 测试不同语言在缩放下的显示效果
// 5. 响应式设计: 测试文本缩放对布局的影响

// CDPEmulationSetEmulatedOSTextScale 设置模拟的操作系统文本缩放
// 参数说明:
//   - textScaleFactor: 文本缩放因子，通常为1.0、1.25、1.5、1.75、2.0等
func CDPEmulationSetEmulatedOSTextScale(params string) (string, error) {
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
		"method": "Emulation.setEmulatedOSTextScale",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setEmulatedOSTextScale 请求失败: %w", err)
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

				log.Println("✅ 操作系统文本缩放设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setEmulatedOSTextScale 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setEmulatedVisionDeficiency  -----------------------------------------------
// === 应用场景 ===
// 1. 颜色盲模拟: 测试网站对不同类型色盲用户的兼容性
// 2. 视力障碍模拟: 测试视力模糊或对比度降低情况下的可读性
// 3. 无障碍设计: 验证网站是否符合WCAG无障碍标准
// 4. 用户体验: 确保所有用户都能获得良好的浏览体验

// CDPEmulationSetEmulatedVisionDeficiency 设置模拟的视觉缺陷
// 参数说明:
//   - type: 视觉缺陷类型，支持以下值:
//   - "none": 无视觉缺陷（默认）
//   - "achromatopsia": 全色盲（看不到任何颜色）
//   - "deuteranopia": 绿色盲
//   - "protanopia": 红色盲
//   - "tritanopia": 蓝色盲
//   - "blurredVision": 视力模糊
//   - "reducedContrast": 对比度降低
func CDPEmulationSetEmulatedVisionDeficiency(params string) (string, error) {
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
		"method": "Emulation.setEmulatedVisionDeficiency",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setEmulatedVisionDeficiency 请求失败: %w", err)
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

				log.Printf("✅ 视觉缺陷模拟设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setEmulatedVisionDeficiency 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setGeolocationOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 位置服务测试: 测试基于位置的服务和功能
// 2. 地图应用开发: 测试地图应用在不同位置的行为
// 3. 地理围栏测试: 测试进入/离开特定区域的行为
// 4. 本地化内容: 测试基于位置的内容显示
// 5. 位置权限测试: 测试地理位置权限流程

// CDPEmulationSetGeolocationOverride 设置地理位置覆盖
// 参数说明:
//   - latitude: 纬度，范围 -90 到 90
//   - longitude: 经度，范围 -180 到 180
//   - accuracy: 精度（以米为单位），可选，默认值为0
//
// 其他可选参数:
//   - altitude: 海拔高度（以米为单位），可选
//   - altitudeAccuracy: 海拔精度（以米为单位），可选
//   - heading: 方向（以度为单位，0-360），可选
//   - speed: 速度（以米/秒为单位），可选
func CDPEmulationSetGeolocationOverride(params string) (string, error) {
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
		"method": "Emulation.setGeolocationOverride",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setGeolocationOverride 请求失败: %w", err)
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

				log.Printf("✅ 地理位置覆盖设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setGeolocationOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setIdleOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 节能模式测试: 测试网页在节能模式下的行为
// 2. 屏幕保护程序: 测试屏幕锁定/解锁时的网页行为
// 3. 页面生命周期: 测试Page Lifecycle API的各种状态
// 4. 资源管理: 测试网页在用户不活跃时的资源使用
// 5. 后台任务: 测试网页在后台运行时的行为

// CDPEmulationSetIdleOverride 设置空闲状态覆盖
// 参数说明:
//   - isUserActive: 用户是否活跃，布尔值
//   - isScreenUnlocked: 屏幕是否解锁，布尔值
//
// 可选参数:
//   - idleTime: 空闲时间（以秒为单位），可选
func CDPEmulationSetIdleOverride(params string) (string, error) {
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
		"method": "Emulation.setIdleOverride",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setIdleOverride 请求失败: %w", err)
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

				log.Printf("✅ 空闲状态覆盖设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setIdleOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setScriptExecutionDisabled  -----------------------------------------------
// === 应用场景 ===
// 1. 无脚本测试: 测试网页在没有JavaScript时的基本功能
// 2. 性能基准测试: 测量纯HTML/CSS渲染性能
// 3. 可访问性测试: 确保网站在JS禁用时仍可访问
// 4. 安全分析: 分析XSS等安全漏洞
// 5. 渐进增强测试: 验证网站是否遵循渐进增强原则
// 6. 搜索引擎优化测试: 确保搜索引擎可抓取内容
// 7. 网络条件模拟: 测试慢速网络下的回退方案
// 8. 浏览器兼容性测试: 测试不支持JS的老旧浏览器

// CDPEmulationSetScriptExecutionDisabled 设置脚本执行禁用状态
// 参数说明:
//   - value: 布尔值，true表示禁用脚本执行，false表示启用脚本执行
func CDPEmulationSetScriptExecutionDisabled(params string) (string, error) {
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
		"method": "Emulation.setScriptExecutionDisabled",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setScriptExecutionDisabled 请求失败: %w", err)
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
			return "", fmt.Errorf("setScriptExecutionDisabled 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setTimezoneOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 时区相关功能测试: 测试日期时间处理、时区转换
// 2. 国际化测试: 测试不同时区下的应用行为
// 3. 夏令时测试: 测试夏令时转换
// 4. 日历应用测试: 测试日历事件的时间显示
// 5. 会议调度测试: 测试跨时区会议调度
// 6. 时间敏感功能: 测试限时活动、定时任务
// 7. 历史日期测试: 测试历史时区规则

// CDPEmulationSetTimezoneOverride 设置时区覆盖
// 参数说明:
//   - timezoneId: 时区ID字符串，遵循IANA时区数据库格式
//
// 常见时区ID示例:
//   - "America/New_York" (美国东部时间，EST/EDT)
//   - "America/Chicago" (美国中部时间，CST/CDT)
//   - "America/Denver" (美国山地时间，MST/MDT)
//   - "America/Los_Angeles" (美国太平洋时间，PST/PDT)
//   - "Europe/London" (格林尼治标准时间，GMT/BST)
//   - "Europe/Paris" (欧洲中部时间，CET/CEST)
//   - "Asia/Shanghai" (中国标准时间，CST，UTC+8)
//   - "Asia/Tokyo" (日本标准时间，JST，UTC+9)
//   - "Australia/Sydney" (澳大利亚东部时间，AEST/AEDT)
func CDPEmulationSetTimezoneOverride(params string) (string, error) {
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
		"method": "Emulation.setTimezoneOverride",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setTimezoneOverride 请求失败: %w", err)
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
			return "", fmt.Errorf("setTimezoneOverride 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setTouchEmulationEnabled  -----------------------------------------------
// === 应用场景 ===
// 1. 移动端网页测试: 测试响应式网页的触摸交互
// 2. 触摸事件测试: 测试触摸手势、滑动、缩放等功能
// 3. 多指触控测试: 测试多点触控功能
// 4. 移动端兼容性测试: 确保网页在触摸设备上正常工作
// 5. 渐进增强测试: 测试触摸和鼠标事件兼容性
// 6. 手势识别测试: 测试滑动、长按、双击等手势
// 7. 移动端性能测试: 测试触摸事件的性能影响
// 8. 跨设备测试: 确保桌面和移动端体验一致

// CDPEmulationSetTouchEmulationEnabled 设置触摸模拟启用状态
// 参数说明:
//   - enabled: 布尔值，true表示启用触摸模拟，false表示禁用触摸模拟
//   - configuration: 可选参数，触摸模拟配置对象
//   - maxTouchPoints: 最大触摸点数，默认值为1
//   - configurationType: 配置类型，如"mobile"、"desktop"等
func CDPEmulationSetTouchEmulationEnabled(params string) (string, error) {
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
		"method": "Emulation.setTouchEmulationEnabled",
		"params": %s
	}`, reqID, params)
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setTouchEmulationEnabled 请求失败: %w", err)
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
			return "", fmt.Errorf("setTouchEmulationEnabled 请求超时")
		}
	}
}

// -----------------------------------------------  Emulation.setUserAgentOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 浏览器兼容性测试: 测试网站在不同浏览器下的表现
// 2. 设备模拟测试: 模拟移动端、平板、桌面设备
// 3. 爬虫检测绕过: 模拟不同浏览器避免被检测为爬虫
// 4. 响应式设计测试: 测试网站在不同设备上的布局
// 5. 功能特性检测: 测试浏览器特定功能的兼容性
// 6. 服务端渲染测试: 测试不同用户代理的服务器响应
// 7. 地理位置测试: 测试基于用户代理的内容本地化
// 8. 安全性测试: 测试用户代理伪造的防护

// CDPEmulationSetUserAgentOverride 设置用户代理覆盖
// 参数说明:
//   - userAgent: 字符串，要设置的用户代理字符串
//   - acceptLanguage: 可选，字符串，设置Accept-Language头部
//   - platform: 可选，字符串，设置平台信息
//   - userAgentMetadata: 可选，对象，用户代理元数据
func CDPEmulationSetUserAgentOverride(params string) (string, error) {
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
		"method": "Emulation.setUserAgentOverride",
		"params": %s
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setUserAgentOverride 请求失败: %w", err)
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
			return "", fmt.Errorf("setUserAgentOverride 请求超时")
		}
	}
}
