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

/*

// 示例1: 设备模拟清理和恢复
func exampleDeviceSimulationCleanup() {
	// === 应用场景描述 ===
	// 场景: 设备模拟清理和恢复
	// 用途: 在设备模拟测试完成后恢复原始设备状态
	// 优势: 确保后续测试不受之前设备模拟的影响
	// 典型工作流: 设置设备模拟 -> 执行测试 -> 清理设备覆盖 -> 验证恢复

	log.Println("设备模拟清理和恢复示例...")

	// 定义设备模拟清理场景
	deviceCleanupScenarios := []struct {
		name        string
		description string
		deviceType  string
		simulationConfig map[string]interface{}
		cleanupReason string
		expectedRestoration []string
	}{
		{
			name:        "移动设备模拟清理",
			description: "清理移动设备模拟，恢复桌面环境",
			deviceType:  "iPhone 12",
			simulationConfig: map[string]interface{}{
				"width": 390,
				"height": 844,
				"deviceScaleFactor": 3.0,
				"mobile": true,
			},
			cleanupReason: "移动端测试完成，返回桌面环境",
			expectedRestoration: []string{
				"视口尺寸恢复",
				"设备像素比恢复",
				"移动设备标识清除",
				"触摸支持状态恢复",
			},
		},
		{
			name:        "平板设备模拟清理",
			description: "清理平板设备模拟，恢复原始状态",
			deviceType:  "iPad Pro",
			simulationConfig: map[string]interface{}{
				"width": 1024,
				"height": 1366,
				"deviceScaleFactor": 2.0,
				"mobile": true,
			},
			cleanupReason: "平板端响应式测试完成",
			expectedRestoration: []string{
				"屏幕分辨率恢复",
				"视口方向重置",
				"设备类型标识清除",
			},
		},
		{
			name:        "高清屏幕模拟清理",
			description: "清理高DPI设备模拟",
			deviceType:  "4K显示器",
			simulationConfig: map[string]interface{}{
				"width": 3840,
				"height": 2160,
				"deviceScaleFactor": 2.0,
				"mobile": false,
			},
			cleanupReason: "高DPI测试完成，返回标准屏幕",
			expectedRestoration: []string{
				"像素密度恢复",
				"屏幕尺寸重置",
				"视网膜显示标志清除",
			},
		},
		{
			name:        "小屏幕设备清理",
			description: "清理小屏幕设备模拟",
			deviceType:  "功能手机",
			simulationConfig: map[string]interface{}{
				"width": 320,
				"height": 480,
				"deviceScaleFactor": 1.0,
				"mobile": true,
			},
			cleanupReason: "移动优先设计验证完成",
			expectedRestoration: []string{
				"窄屏适配清除",
				"移动布局恢复",
				"触摸目标大小重置",
			},
		},
		{
			name:        "自定义分辨率清理",
			description: "清理自定义分辨率模拟",
			deviceType:  "自定义分辨率",
			simulationConfig: map[string]interface{}{
				"width": 1920,
				"height": 1080,
				"deviceScaleFactor": 1.5,
				"mobile": false,
			},
			cleanupReason: "特定分辨率测试完成",
			expectedRestoration: []string{
				"宽高比重置",
				"自定义尺寸清除",
				"比例因子恢复",
			},
		},
	}

	// 测试每种清理场景
	for i, scenario := range deviceCleanupScenarios {
		log.Printf("\n=== 设备清理场景 %d/%d: %s ===", i+1, len(deviceCleanupScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("设备类型: %s", scenario.deviceType)
		log.Printf("清理原因: %s", scenario.cleanupReason)
		log.Printf("模拟配置: %v", scenario.simulationConfig)
		log.Printf("预期恢复项:")
		for _, item := range scenario.expectedRestoration {
			log.Printf("  - %s", item)
		}

		// 模拟设备指标覆盖状态
		simulatedState := simulateDeviceMetricsOverride(scenario.deviceType, scenario.simulationConfig)
		log.Printf("当前设备覆盖状态:")
		for key, value := range simulatedState {
			log.Printf("  %s: %v", key, value)
		}

		// 执行设备指标覆盖清理
		log.Printf("执行设备指标覆盖清理...")
		response, err := CDPEmulationClearDeviceMetricsOverride()
		if err != nil {
			log.Printf("清理设备指标覆盖失败: %v", err)
			continue
		}

		log.Printf("✅ 设备指标覆盖清理成功: %s", response)

		// 模拟清理后的状态
		restoredState := simulateDeviceMetricsRestoration(scenario.expectedRestoration)

		// 分析清理效果
		analyzeDeviceCleanupEffect(simulatedState, restoredState, scenario.expectedRestoration)

		// 提供使用建议
		provideDeviceCleanupAdvice(scenario.deviceType, scenario.cleanupReason)

		// 短暂延迟
		time.Sleep(100 * time.Millisecond)
	}
}

func simulateDeviceMetricsOverride(deviceType string, config map[string]interface{}) map[string]interface{} {
	state := make(map[string]interface{})

	// 模拟设备指标覆盖状态
	state["deviceType"] = deviceType
	state["isOverridden"] = true
	state["overrideTime"] = time.Now().Add(-2 * time.Minute).Format("15:04:05")

	// 添加配置项
	for key, value := range config {
		state[key] = value
	}

	// 添加模拟的影响
	state["affectedAPIs"] = []string{
		"window.innerWidth",
		"window.innerHeight",
		"window.screen.width",
		"window.screen.height",
		"window.devicePixelRatio",
		"CSS媒体查询",
		"响应式布局",
	}

	state["userAgentModified"] = false
	state["touchEmulation"] = config["mobile"].(bool)

	return state
}

func simulateDeviceMetricsRestoration(expectedRestoration []string) map[string]interface{} {
	state := make(map[string]interface{})

	// 模拟恢复后的状态
	state["isOverridden"] = false
	state["restorationTime"] = time.Now().Format("15:04:05")
	state["restorationComplete"] = true

	// 模拟恢复的项
	restoredItems := make(map[string]bool)
	for _, item := range expectedRestoration {
		restoredItems[item] = true
	}

	state["restoredItems"] = restoredItems

	// 模拟恢复后的设备值
	state["width"] = 1920
	state["height"] = 1080
	state["deviceScaleFactor"] = 1.0
	state["mobile"] = false
	state["touchEnabled"] = false

	// 模拟恢复的影响
	state["originalAPIsRestored"] = []string{
		"window.innerWidth",
		"window.innerHeight",
		"window.screen.width",
		"window.screen.height",
		"window.devicePixelRatio",
		"navigator.userAgent",
		"CSS媒体查询响应",
	}

	return state
}

func analyzeDeviceCleanupEffect(originalState, restoredState map[string]interface{}, expectedRestoration []string) {
	log.Printf("设备清理效果分析:")

	wasOverridden, _ := originalState["isOverridden"].(bool)
	isOverridden, _ := restoredState["isOverridden"].(bool)

	log.Printf("  覆盖状态: %v -> %v", wasOverridden, isOverridden)
	log.Printf("  清理完成: %v", !isOverridden)

	if wasOverridden && !isOverridden {
		log.Printf("  ✅ 成功清除设备指标覆盖")
	} else if !wasOverridden && !isOverridden {
		log.Printf("  ℹ 设备指标未覆盖，无需清理")
	} else {
		log.Printf("  ❌ 设备指标覆盖状态异常")
	}

	// 检查恢复项
	if restoredItems, ok := restoredState["restoredItems"].(map[string]bool); ok {
		log.Printf("  恢复项验证:")
		for _, expected := range expectedRestoration {
			if restoredItems[expected] {
				log.Printf("    ✅ %s: 已恢复", expected)
			} else {
				log.Printf("    ❌ %s: 未恢复", expected)
			}
		}
	}

	// 影响分析
	if originalAPIs, ok := originalState["affectedAPIs"].([]string); ok {
		log.Printf("  影响的API数量: %d", len(originalAPIs))
	}

	if restoredAPIs, ok := restoredState["originalAPIsRestored"].([]string); ok {
		log.Printf("  恢复的API数量: %d", len(restoredAPIs))
	}

	// 设备值对比
	log.Printf("  设备值对比:")
	originalWidth, _ := originalState["width"].(int)
	restoredWidth, _ := restoredState["width"].(int)
	log.Printf("    宽度: %d -> %d", originalWidth, restoredWidth)

	originalHeight, _ := originalState["height"].(int)
	restoredHeight, _ := restoredState["height"].(int)
	log.Printf("    高度: %d -> %d", originalHeight, restoredHeight)

	originalScale, _ := originalState["deviceScaleFactor"].(float64)
	restoredScale, _ := restoredState["deviceScaleFactor"].(float64)
	log.Printf("    设备比例因子: %.1f -> %.1f", originalScale, restoredScale)
}

func provideDeviceCleanupAdvice(deviceType, cleanupReason string) {
	log.Printf("使用建议:")

	switch cleanupReason {
	case "移动端测试完成，返回桌面环境":
		log.Printf("  ⚡ 移动端测试完成后清理是必要的")
		log.Printf("    建议: 在移动端和桌面端测试之间切换时执行清理")

	case "平板端响应式测试完成":
		log.Printf("  ⚡ 平板设备模拟清理确保后续测试准确性")
		log.Printf("    建议: 在不同设备类型测试后都执行清理")

	case "高DPI测试完成，返回标准屏幕":
		log.Printf("  ⚡ 高DPI模拟会影响字体渲染和图片显示")
		log.Printf("    建议: 在高DPI测试后立即清理")

	case "移动优先设计验证完成":
		log.Printf("  ⚡ 小屏幕设备模拟清理避免布局问题")
		log.Printf("    建议: 在响应式设计验证后清理")

	case "特定分辨率测试完成":
		log.Printf("  ⚡ 自定义分辨率清理确保测试环境干净")
		log.Printf("    建议: 在分辨率相关的测试套件中自动清理")
	}

	// 通用建议
	log.Printf("  通用建议:")
	log.Printf("    - 在测试套件开始时和结束时都清理设备覆盖")
	log.Printf("    - 在不同设备模拟之间切换时执行清理")
	log.Printf("    - 自动化测试中确保清理步骤")
	log.Printf("    - 错误处理中包括设备覆盖清理")
}

*/

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

/*

// 示例1: 地理位置模拟清理和恢复
func exampleGeolocationSimulationCleanup() {
	// === 应用场景描述 ===
	// 场景: 地理位置模拟清理和恢复
	// 用途: 在地理位置模拟测试完成后恢复原始定位
	// 优势: 确保后续功能不受模拟地理位置的影响
	// 典型工作流: 设置位置模拟 -> 执行测试 -> 清理位置覆盖 -> 验证恢复

	log.Println("地理位置模拟清理和恢复示例...")

	// 定义地理位置模拟清理场景
	geolocationCleanupScenarios := []struct {
		name        string
		description string
		location    string
		simulationConfig map[string]interface{}
		cleanupReason string
		expectedRestoration []string
	}{
		{
			name:        "美国硅谷位置清理",
			description: "清理美国硅谷的地理位置模拟",
			location:    "硅谷, 美国",
			simulationConfig: map[string]interface{}{
				"latitude": 37.7749,
				"longitude": -122.4194,
				"accuracy": 50,
				"altitude": 20,
				"altitudeAccuracy": 10,
				"heading": 45.0,
				"speed": 0.5,
			},
			cleanupReason: "地理位置相关的功能测试完成",
			expectedRestoration: []string{
				"GPS坐标恢复",
				"定位精度重置",
				"海拔高度清除",
				"行进方向恢复",
				"速度信息清除",
			},
		},
		{
			name:        "中国北京位置清理",
			description: "清理中国北京的地理位置模拟",
			location:    "北京, 中国",
			simulationConfig: map[string]interface{}{
				"latitude": 39.9042,
				"longitude": 116.4074,
				"accuracy": 100,
				"altitude": 50,
				"altitudeAccuracy": 20,
				"heading": 90.0,
				"speed": 0.0,
			},
			cleanupReason: "时区相关的功能验证完成",
			expectedRestoration: []string{
				"地理位置API恢复",
				"时区信息重置",
				"本地化功能恢复",
				"地图服务重置",
			},
		},
		{
			name:        "伦敦位置清理",
			description: "清理英国伦敦的地理位置模拟",
			location:    "伦敦, 英国",
			simulationConfig: map[string]interface{}{
				"latitude": 51.5074,
				"longitude": -0.1278,
				"accuracy": 25,
				"altitude": 15,
				"altitudeAccuracy": 5,
				"heading": 180.0,
				"speed": 1.2,
			},
			cleanupReason: "多语言和本地化测试完成",
			expectedRestoration: []string{
				"定位权限恢复",
				"语言设置重置",
				"货币格式恢复",
				"日期格式重置",
			},
		},
		{
			name:        "东京位置清理",
			description: "清理日本东京的地理位置模拟",
			location:    "东京, 日本",
			simulationConfig: map[string]interface{}{
				"latitude": 35.6762,
				"longitude": 139.6503,
				"accuracy": 30,
				"altitude": 10,
				"altitudeAccuracy": 8,
				"heading": 270.0,
				"speed": 0.8,
			},
			cleanupReason: "RTL布局和东方语言测试完成",
			expectedRestoration: []string{
				"文本方向恢复",
				"字体渲染重置",
				"输入法支持恢复",
				"UI布局重置",
			},
		},
		{
			name:        "悉尼位置清理",
			description: "清理澳大利亚悉尼的地理位置模拟",
			location:    "悉尼, 澳大利亚",
			simulationConfig: map[string]interface{}{
				"latitude": -33.8688,
				"longitude": 151.2093,
				"accuracy": 80,
				"altitude": 5,
				"altitudeAccuracy": 3,
				"heading": 0.0,
				"speed": 0.0,
			},
			cleanupReason: "南半球特定功能测试完成",
			expectedRestoration: []string{
				"气候相关功能恢复",
				"季节性功能重置",
				"时区计算恢复",
				"地图投影重置",
			},
		},
	}

	// 测试每种清理场景
	for i, scenario := range geolocationCleanupScenarios {
		log.Printf("\n=== 地理位置清理场景 %d/%d: %s ===", i+1, len(geolocationCleanupScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("位置: %s", scenario.location)
		log.Printf("清理原因: %s", scenario.cleanupReason)
		log.Printf("模拟配置: %v", scenario.simulationConfig)
		log.Printf("预期恢复项:")
		for _, item := range scenario.expectedRestoration {
			log.Printf("  - %s", item)
		}

		// 模拟地理位置覆盖状态
		simulatedState := simulateGeolocationOverride(scenario.location, scenario.simulationConfig)
		log.Printf("当前地理位置覆盖状态:")
		for key, value := range simulatedState {
			log.Printf("  %s: %v", key, value)
		}

		// 执行地理位置覆盖清理
		log.Printf("执行地理位置覆盖清理...")
		response, err := CDPEmulationClearGeolocationOverride()
		if err != nil {
			log.Printf("清理地理位置覆盖失败: %v", err)
			continue
		}

		log.Printf("✅ 地理位置覆盖清理成功: %s", response)

		// 模拟清理后的状态
		restoredState := simulateGeolocationRestoration(scenario.expectedRestoration)

		// 分析清理效果
		analyzeGeolocationCleanupEffect(simulatedState, restoredState, scenario.expectedRestoration)

		// 提供使用建议
		provideGeolocationCleanupAdvice(scenario.location, scenario.cleanupReason)

		// 短暂延迟
		time.Sleep(100 * time.Millisecond)
	}
}

func simulateGeolocationOverride(location string, config map[string]interface{}) map[string]interface{} {
	state := make(map[string]interface{})

	// 模拟地理位置覆盖状态
	state["location"] = location
	state["isOverridden"] = true
	state["overrideTime"] = time.Now().Add(-3 * time.Minute).Format("15:04:05")

	// 添加配置项
	for key, value := range config {
		state[key] = value
	}

	// 添加模拟的影响
	state["affectedAPIs"] = []string{
		"navigator.geolocation",
		"Geolocation.getCurrentPosition()",
		"Geolocation.watchPosition()",
		"位置相关的Web API",
		"地图服务集成",
	}

	state["affectedFeatures"] = []string{
		"基于位置的内容",
		"本地化服务",
		"天气信息",
		"交通信息",
		"位置相关的广告",
	}

	state["privacyImplications"] = []string{
		"位置信息伪造",
		"隐私数据泄露风险",
		"位置追踪干扰",
		"地理围栏测试",
	}

	return state
}

func simulateGeolocationRestoration(expectedRestoration []string) map[string]interface{} {
	state := make(map[string]interface{})

	// 模拟恢复后的状态
	state["isOverridden"] = false
	state["restorationTime"] = time.Now().Format("15:04:05")
	state["restorationComplete"] = true

	// 模拟恢复的项
	restoredItems := make(map[string]bool)
	for _, item := range expectedRestoration {
		restoredItems[item] = true
	}

	state["restoredItems"] = restoredItems

	// 模拟恢复后的默认值
	state["latitude"] = 0.0
	state["longitude"] = 0.0
	state["accuracy"] = 0.0
	state["hasPosition"] = false

	// 模拟恢复的影响
	state["originalAPIsRestored"] = []string{
		"navigator.geolocation.getCurrentPosition",
		"Geolocation API 权限",
		"真实位置检测",
		"浏览器原生定位",
	}

	state["privacyRestored"] = []string{
		"真实位置隐私",
		"位置权限控制",
		"地理围栏重置",
		"追踪保护恢复",
	}

	return state
}

func analyzeGeolocationCleanupEffect(originalState, restoredState map[string]interface{}, expectedRestoration []string) {
	log.Printf("地理位置清理效果分析:")

	wasOverridden, _ := originalState["isOverridden"].(bool)
	isOverridden, _ := restoredState["isOverridden"].(bool)

	log.Printf("  覆盖状态: %v -> %v", wasOverridden, isOverridden)
	log.Printf("  清理完成: %v", !isOverridden)

	if wasOverridden && !isOverridden {
		log.Printf("  ✅ 成功清除地理位置覆盖")
	} else if !wasOverridden && !isOverridden {
		log.Printf("  ℹ 地理位置未覆盖，无需清理")
	} else {
		log.Printf("  ❌ 地理位置覆盖状态异常")
	}

	// 检查恢复项
	if restoredItems, ok := restoredState["restoredItems"].(map[string]bool); ok {
		log.Printf("  恢复项验证:")
		for _, expected := range expectedRestoration {
			if restoredItems[expected] {
				log.Printf("    ✅ %s: 已恢复", expected)
			} else {
				log.Printf("    ❌ %s: 未恢复", expected)
			}
		}
	}

	// 位置值对比
	log.Printf("  位置值对比:")
	originalLat, _ := originalState["latitude"].(float64)
	restoredLat, _ := restoredState["latitude"].(float64)
	log.Printf("    纬度: %.4f -> %.4f", originalLat, restoredLat)

	originalLon, _ := originalState["longitude"].(float64)
	restoredLon, _ := restoredState["longitude"].(float64)
	log.Printf("    经度: %.4f -> %.4f", originalLon, restoredLon)

	originalAcc, _ := originalState["accuracy"].(float64)
	restoredAcc, _ := restoredState["accuracy"].(float64)
	log.Printf("    精度: %.1f米 -> %.1f米", originalAcc, restoredAcc)

	// 隐私影响分析
	if originalPrivacy, ok := originalState["privacyImplications"].([]string); ok {
		log.Printf("  原始隐私影响: %d 项", len(originalPrivacy))
	}

	if restoredPrivacy, ok := restoredState["privacyRestored"].([]string); ok {
		log.Printf("  隐私恢复项: %d 项", len(restoredPrivacy))
	}
}

func provideGeolocationCleanupAdvice(location, cleanupReason string) {
	log.Printf("使用建议:")

	switch cleanupReason {
	case "地理位置相关的功能测试完成":
		log.Printf("  🌍 位置功能测试后清理是必要的")
		log.Printf("    建议: 在不同地区功能测试之间执行清理")

	case "时区相关的功能验证完成":
		log.Printf("  🌍 时区功能测试后清理确保准确性")
		log.Printf("    建议: 在跨时区测试后立即清理")

	case "多语言和本地化测试完成":
		log.Printf("  🌍 本地化测试清理避免语言设置冲突")
		log.Printf("    建议: 在本地化测试套件中自动清理")

	case "RTL布局和东方语言测试完成":
		log.Printf("  🌍 地区特定UI测试后清理重置布局")
		log.Printf("    建议: 在UI国际化测试后执行清理")

	case "南半球特定功能测试完成":
		log.Printf("  🌍 半球特定功能测试清理确保功能正常")
		log.Printf("    建议: 在气候和季节相关测试后清理")
	}

	// 通用建议
	log.Printf("  通用建议:")
	log.Printf("    - 地理位置模拟可能影响隐私，测试后必须清理")
	log.Printf("    - 自动化测试中确保地理位置状态清理")
	log.Printf("    - 在切换测试地区时执行清理")
	log.Printf("    - 监控地理位置API的使用情况")
}

*/

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

/*

// 示例1: 空闲状态模拟清理和恢复
func exampleIdleStateSimulationCleanup() {
	// === 应用场景描述 ===
	// 场景: 空闲状态模拟清理和恢复
	// 用途: 在空闲状态模拟测试完成后恢复原始空闲检测
	// 优势: 确保后续功能不受模拟空闲状态的影响
	// 典型工作流: 设置空闲状态模拟 -> 执行测试 -> 清理空闲覆盖 -> 验证恢复

	log.Println("空闲状态模拟清理和恢复示例...")

	// 定义空闲状态模拟清理场景
	idleCleanupScenarios := []struct {
		name        string
		description string
		idleState   string
		simulationConfig map[string]interface{}
		cleanupReason string
		expectedRestoration []string
	}{
		{
			name:        "活动状态清理",
			description: "清理模拟的用户活动状态",
			idleState:   "活跃状态",
			simulationConfig: map[string]interface{}{
				"isUserActive": true,
				"isScreenUnlocked": true,
				"idleTime": 0,
			},
			cleanupReason: "用户活跃性测试完成",
			expectedRestoration: []string{
				"用户活动检测恢复",
				"屏幕锁定状态重置",
				"空闲计时器恢复",
				"活动监听器重置",
			},
		},
		{
			name:        "空闲状态清理",
			description: "清理模拟的用户空闲状态",
			idleState:   "空闲状态",
			simulationConfig: map[string]interface{}{
				"isUserActive": false,
				"isScreenUnlocked": true,
				"idleTime": 300, // 5分钟
			},
			cleanupReason: "自动锁定功能测试完成",
			expectedRestoration: []string{
				"空闲检测恢复",
				"自动锁定重置",
				"节能模式恢复",
				"通知抑制重置",
			},
		},
		{
			name:        "锁定状态清理",
			description: "清理模拟的屏幕锁定状态",
			idleState:   "锁定状态",
			simulationConfig: map[string]interface{}{
				"isUserActive": false,
				"isScreenUnlocked": false,
				"idleTime": 600, // 10分钟
			},
			cleanupReason: "屏幕保护测试完成",
			expectedRestoration: []string{
				"屏幕状态检测恢复",
				"保护程序重置",
				"密码保护恢复",
				"会话管理重置",
			},
		},
		{
			name:        "深度空闲清理",
			description: "清理模拟的深度空闲状态",
			idleState:   "深度空闲",
			simulationConfig: map[string]interface{}{
				"isUserActive": false,
				"isScreenUnlocked": true,
				"idleTime": 1800, // 30分钟
			},
			cleanupReason: "深度节能功能测试完成",
			expectedRestoration: []string{
				"深度空闲检测恢复",
				"后台任务管理重置",
				"资源限制恢复",
				"网络连接重置",
			},
		},
		{
			name:        "混合状态清理",
			description: "清理复杂的混合空闲状态",
			idleState:   "混合状态",
			simulationConfig: map[string]interface{}{
				"isUserActive": true,  // 用户活动
				"isScreenUnlocked": false, // 屏幕锁定
				"idleTime": 120, // 2分钟
			},
			cleanupReason: "异常状态处理测试完成",
			expectedRestoration: []string{
				"状态机恢复",
				"异常处理重置",
				"状态同步恢复",
				"冲突解决重置",
			},
		},
	}

	// 测试每种清理场景
	for i, scenario := range idleCleanupScenarios {
		log.Printf("\n=== 空闲状态清理场景 %d/%d: %s ===", i+1, len(idleCleanupScenarios), scenario.name)
		log.Printf("描述: %s", scenario.description)
		log.Printf("空闲状态: %s", scenario.idleState)
		log.Printf("清理原因: %s", scenario.cleanupReason)
		log.Printf("模拟配置: %v", scenario.simulationConfig)
		log.Printf("预期恢复项:")
		for _, item := range scenario.expectedRestoration {
			log.Printf("  - %s", item)
		}

		// 模拟空闲状态覆盖
		simulatedState := simulateIdleStateOverride(scenario.idleState, scenario.simulationConfig)
		log.Printf("当前空闲状态覆盖:")
		for key, value := range simulatedState {
			log.Printf("  %s: %v", key, value)
		}

		// 执行空闲状态覆盖清理
		log.Printf("执行空闲状态覆盖清理...")
		response, err := CDPEmulationClearIdleOverride()
		if err != nil {
			log.Printf("清理空闲状态覆盖失败: %v", err)
			continue
		}

		log.Printf("✅ 空闲状态覆盖清理成功: %s", response)

		// 模拟清理后的状态
		restoredState := simulateIdleStateRestoration(scenario.expectedRestoration)

		// 分析清理效果
		analyzeIdleStateCleanupEffect(simulatedState, restoredState, scenario.expectedRestoration)

		// 提供使用建议
		provideIdleStateCleanupAdvice(scenario.idleState, scenario.cleanupReason)

		// 短暂延迟
		time.Sleep(100 * time.Millisecond)
	}
}

func simulateIdleStateOverride(idleState string, config map[string]interface{}) map[string]interface{} {
	state := make(map[string]interface{})

	// 模拟空闲状态覆盖
	state["idleState"] = idleState
	state["isOverridden"] = true
	state["overrideTime"] = time.Now().Add(-4 * time.Minute).Format("15:04:05")

	// 添加配置项
	for key, value := range config {
		state[key] = value
	}

	// 添加模拟的影响
	state["affectedAPIs"] = []string{
		"navigator.getBattery",
		"Idle Detection API",
		"Page Visibility API",
		"Screen Wake Lock API",
		"Power Saving API",
	}

	state["affectedBehaviors"] = []string{
		"自动锁定",
		"屏幕保护",
		"通知抑制",
		"后台任务限制",
		"网络节流",
	}

	state["powerImplications"] = []string{
		"电池消耗模拟",
		"电源管理干扰",
		"性能限制影响",
		"资源分配异常",
	}

	// 计算空闲时间
	if idleTime, ok := config["idleTime"].(int); ok {
		state["idleDuration"] = fmt.Sprintf("%d秒", idleTime)
		state["idleLevel"] = calculateIdleLevel(idleTime)
	}

	return state
}

func calculateIdleLevel(idleTime int) string {
	switch {
	case idleTime == 0:
		return "活跃"
	case idleTime <= 60:
		return "轻微空闲"
	case idleTime <= 300:
		return "中度空闲"
	case idleTime <= 900:
		return "高度空闲"
	default:
		return "深度空闲"
	}
}

func simulateIdleStateRestoration(expectedRestoration []string) map[string]interface{} {
	state := make(map[string]interface{})

	// 模拟恢复后的状态
	state["isOverridden"] = false
	state["restorationTime"] = time.Now().Format("15:04:05")
	state["restorationComplete"] = true

	// 模拟恢复的项
	restoredItems := make(map[string]bool)
	for _, item := range expectedRestoration {
		restoredItems[item] = true
	}

	state["restoredItems"] = restoredItems

	// 模拟恢复后的默认值
	state["isUserActive"] = true
	state["isScreenUnlocked"] = true
	state["idleTime"] = 0
	state["idleLevel"] = "活跃"

	// 模拟恢复的影响
	state["originalAPIsRestored"] = []string{
		"真实的空闲检测",
		"原生电源管理",
		"实际用户活动",
		"屏幕状态检测",
	}

	state["powerManagementRestored"] = []string{
		"正常电池消耗",
		"真实电源管理",
		"准确的性能测量",
		"正确的资源分配",
	}

	return state
}

func analyzeIdleStateCleanupEffect(originalState, restoredState map[string]interface{}, expectedRestoration []string) {
	log.Printf("空闲状态清理效果分析:")

	wasOverridden, _ := originalState["isOverridden"].(bool)
	isOverridden, _ := restoredState["isOverridden"].(bool)

	log.Printf("  覆盖状态: %v -> %v", wasOverridden, isOverridden)
	log.Printf("  清理完成: %v", !isOverridden)

	if wasOverridden && !isOverridden {
		log.Printf("  ✅ 成功清除空闲状态覆盖")
	} else if !wasOverridden && !isOverridden {
		log.Printf("  ℹ 空闲状态未覆盖，无需清理")
	} else {
		log.Printf("  ❌ 空闲状态覆盖状态异常")
	}

	// 检查恢复项
	if restoredItems, ok := restoredState["restoredItems"].(map[string]bool); ok {
		log.Printf("  恢复项验证:")
		for _, expected := range expectedRestoration {
			if restoredItems[expected] {
				log.Printf("    ✅ %s: 已恢复", expected)
			} else {
				log.Printf("    ❌ %s: 未恢复", expected)
			}
		}
	}

	// 状态值对比
	log.Printf("  状态值对比:")
	originalActive, _ := originalState["isUserActive"].(bool)
	restoredActive, _ := restoredState["isUserActive"].(bool)
	log.Printf("    用户活动: %v -> %v", originalActive, restoredActive)

	originalUnlocked, _ := originalState["isScreenUnlocked"].(bool)
	restoredUnlocked, _ := restoredState["isScreenUnlocked"].(bool)
	log.Printf("    屏幕锁定: %v -> %v", originalUnlocked, restoredUnlocked)

	originalIdleTime, _ := originalState["idleTime"].(int)
	restoredIdleTime, _ := restoredState["idleTime"].(int)
	log.Printf("    空闲时间: %d秒 -> %d秒", originalIdleTime, restoredIdleTime)

	originalLevel, _ := originalState["idleLevel"].(string)
	restoredLevel, _ := restoredState["idleLevel"].(string)
	log.Printf("    空闲等级: %s -> %s", originalLevel, restoredLevel)

	// 电源管理影响分析
	if originalPower, ok := originalState["powerImplications"].([]string); ok {
		log.Printf("  原始电源影响: %d 项", len(originalPower))
	}

	if restoredPower, ok := restoredState["powerManagementRestored"].([]string); ok {
		log.Printf("  电源管理恢复项: %d 项", len(restoredPower))
	}
}

func provideIdleStateCleanupAdvice(idleState, cleanupReason string) {
	log.Printf("使用建议:")

	switch cleanupReason {
	case "用户活跃性测试完成":
		log.Printf("  🔋 活跃状态测试后清理确保准确性")
		log.Printf("    建议: 在用户交互测试后立即清理")

	case "自动锁定功能测试完成":
		log.Printf("  🔋 锁定功能测试清理避免误锁")
		log.Printf("    建议: 在自动锁定测试套件中自动清理")

	case "屏幕保护测试完成":
		log.Printf("  🔋 屏幕保护测试清理恢复显示")
		log.Printf("    建议: 在显示相关测试后执行清理")

	case "深度节能功能测试完成":
		log.Printf("  🔋 深度节能测试清理恢复性能")
		log.Printf("    建议: 在性能测试前确保清理完成")

	case "异常状态处理测试完成":
		log.Printf("  🔋 异常状态测试清理重置状态机")
		log.Printf("    建议: 在状态机测试后立即清理")
	}

	// 通用建议
	log.Printf("  通用建议:")
	log.Printf("    - 空闲状态模拟可能影响性能测试，测试后必须清理")
	log.Printf("    - 自动化测试中确保空闲状态清理")
	log.Printf("    - 在切换测试场景时执行清理")
	log.Printf("    - 监控电池和电源管理API的使用情况")
}

*/

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

/*

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
		name           string
		description    string
		throttlingRate float64
		deviceCategory string
		performanceImpact string
		testCases      []string
		expectedMetrics map[string]interface{}
	}{
		{
			name:           "高性能桌面",
			description:    "模拟高性能桌面CPU，无限制",
			throttlingRate: 1.0,
			deviceCategory: "高端桌面电脑",
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
			name:           "中端移动设备",
			description:    "模拟中端移动设备CPU性能",
			throttlingRate: 2.5,
			deviceCategory: "中端智能手机",
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
			name:           "低端设备",
			description:    "模拟低端设备CPU性能",
			throttlingRate: 5.0,
			deviceCategory: "低端智能手机/平板",
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
			name:           "极低性能设备",
			description:    "模拟性能极低的设备",
			throttlingRate: 10.0,
			deviceCategory: "老旧设备/物联网设备",
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
			name:           "网络节流设备",
			description:    "模拟网络受限设备的CPU性能",
			throttlingRate: 3.0,
			deviceCategory: "网络受限移动设备",
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
			name:           "节电模式",
			description:    "模拟节电模式下的CPU性能",
			throttlingRate: 4.0,
			deviceCategory: "节电模式设备",
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
			name:           "CPU瓶颈测试",
			description:    "测试CPU密集型应用瓶颈",
			throttlingRate: 8.0,
			deviceCategory: "CPU密集型场景",
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
	name           string
	description    string
	throttlingRate float64
	deviceCategory string
	performanceImpact string
	testCases      []string
	expectedMetrics map[string]interface{}
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
	name           string
	description    string
	throttlingRate float64
	deviceCategory string
	performanceImpact string
	testCases      []string
	expectedMetrics map[string]interface{}
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


*/

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
func CDPEmulationSetDefaultBackgroundColorOverride(color *RGBA) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	var message string
	if color == nil {
		// 清除背景颜色覆盖
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "Emulation.setDefaultBackgroundColorOverride",
			"params": {}
		}`, reqID)
		log.Printf("[DEBUG] 发送清除默认背景颜色覆盖请求")
	} else {
		// 验证颜色值
		if color.R < 0 || color.R > 255 ||
			color.G < 0 || color.G > 255 ||
			color.B < 0 || color.B > 255 {
			return "", fmt.Errorf("颜色分量值必须在0-255范围内: R=%d, G=%d, B=%d",
				color.R, color.G, color.B)
		}

		if color.A < 0 || color.A > 1 {
			return "", fmt.Errorf("alpha值必须在0-1范围内: A=%.2f", color.A)
		}

		// 设置背景颜色覆盖
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "Emulation.setDefaultBackgroundColorOverride",
			"params": {
				"color": {
					"r": %d,
					"g": %d,
					"b": %d,
					"a": %f
				}
			}
		}`, reqID, color.R, color.G, color.B, color.A)

		log.Printf("[DEBUG] 发送设置默认背景颜色覆盖请求: R=%d, G=%d, B=%d, A=%.2f",
			color.R, color.G, color.B, color.A)
	}

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

// RGBA 颜色结构体
type RGBA struct {
	R int     `json:"r"` // 红色分量 (0-255)
	G int     `json:"g"` // 绿色分量 (0-255)
	B int     `json:"b"` // 蓝色分量 (0-255)
	A float64 `json:"a"` // alpha通道 (0-1)
}

// 预定义的颜色常量
var (
	// 基本颜色
	ColorWhite     = &RGBA{R: 255, G: 255, B: 255, A: 1.0}
	ColorBlack     = &RGBA{R: 0, G: 0, B: 0, A: 1.0}
	ColorRed       = &RGBA{R: 255, G: 0, B: 0, A: 1.0}
	ColorGreen     = &RGBA{R: 0, G: 255, B: 0, A: 1.0}
	ColorBlue      = &RGBA{R: 0, G: 0, B: 255, A: 1.0}
	ColorYellow    = &RGBA{R: 255, G: 255, B: 0, A: 1.0}
	ColorCyan      = &RGBA{R: 0, G: 255, B: 255, A: 1.0}
	ColorMagenta   = &RGBA{R: 255, G: 0, B: 255, A: 1.0}
	ColorGray      = &RGBA{R: 128, G: 128, B: 128, A: 1.0}
	ColorDarkGray  = &RGBA{R: 64, G: 64, B: 64, A: 1.0}
	ColorLightGray = &RGBA{R: 192, G: 192, B: 192, A: 1.0}

	// 透明颜色
	ColorTransparent          = &RGBA{R: 0, G: 0, B: 0, A: 0.0}
	ColorSemiTransparentWhite = &RGBA{R: 255, G: 255, B: 255, A: 0.5}
	ColorSemiTransparentBlack = &RGBA{R: 0, G: 0, B: 0, A: 0.5}

	// 打印颜色
	ColorPrintWhite = &RGBA{R: 255, G: 255, B: 255, A: 1.0} // 打印白色
	ColorPrintCream = &RGBA{R: 255, G: 253, B: 240, A: 1.0} // 打印米色
	ColorPrintBlue  = &RGBA{R: 240, G: 248, B: 255, A: 1.0} // 打印淡蓝

	// 主题颜色
	ColorDarkTheme  = &RGBA{R: 18, G: 18, B: 18, A: 1.0}    // 深色主题
	ColorLightTheme = &RGBA{R: 255, G: 255, B: 255, A: 1.0} // 浅色主题
	ColorSepiaTheme = &RGBA{R: 251, G: 240, B: 217, A: 1.0} // 棕褐色主题

	// 无障碍颜色
	ColorHighContrastWhite  = &RGBA{R: 255, G: 255, B: 255, A: 1.0}
	ColorHighContrastBlack  = &RGBA{R: 0, G: 0, B: 0, A: 1.0}
	ColorHighContrastYellow = &RGBA{R: 255, G: 255, B: 0, A: 1.0}
	ColorHighContrastBlue   = &RGBA{R: 0, G: 0, B: 255, A: 1.0}

	// 网页安全色
	ColorWebSafeBlack = &RGBA{R: 0, G: 0, B: 0, A: 1.0}
	ColorWebSafeWhite = &RGBA{R: 255, G: 255, B: 255, A: 1.0}
	ColorWebSafeGray  = &RGBA{R: 153, G: 153, B: 153, A: 1.0}
	ColorWebSafeRed   = &RGBA{R: 255, G: 51, B: 51, A: 1.0}
	ColorWebSafeGreen = &RGBA{R: 51, G: 204, B: 51, A: 1.0}
	ColorWebSafeBlue  = &RGBA{R: 51, G: 102, B: 255, A: 1.0}

	// 品牌颜色
	ColorFacebookBlue = &RGBA{R: 24, G: 119, B: 242, A: 1.0}
	ColorTwitterBlue  = &RGBA{R: 29, G: 161, B: 242, A: 1.0}
	ColorGoogleRed    = &RGBA{R: 234, G: 67, B: 53, A: 1.0}
	ColorGoogleGreen  = &RGBA{R: 52, G: 168, B: 83, A: 1.0}
	ColorGoogleBlue   = &RGBA{R: 66, G: 133, B: 244, A: 1.0}
	ColorGoogleYellow = &RGBA{R: 251, G: 188, B: 5, A: 1.0}
)

// HexToRGBA 从十六进制颜色字符串创建RGBA对象
func HexToRGBA(hex string, alpha float64) (*RGBA, error) {
	// 移除可能的#前缀
	hex = strings.TrimPrefix(hex, "#")

	// 验证长度
	if len(hex) != 6 && len(hex) != 3 && len(hex) != 8 {
		return nil, fmt.Errorf("无效的十六进制颜色: %s (长度必须是3,6或8)", hex)
	}

	var r, g, b int

	// 处理3位十六进制 (#RGB)
	if len(hex) == 3 {
		r1, err := strconv.ParseInt(hex[0:1], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析红色分量失败: %w", err)
		}
		g1, err := strconv.ParseInt(hex[1:2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析绿色分量失败: %w", err)
		}
		b1, err := strconv.ParseInt(hex[2:3], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析蓝色分量失败: %w", err)
		}

		// 扩展3位到6位
		r = int(r1*16 + r1)
		g = int(g1*16 + g1)
		b = int(b1*16 + b1)
	} else if len(hex) == 6 {
		// 处理6位十六进制 (#RRGGBB)
		r64, err := strconv.ParseInt(hex[0:2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析红色分量失败: %w", err)
		}
		g64, err := strconv.ParseInt(hex[2:4], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析绿色分量失败: %w", err)
		}
		b64, err := strconv.ParseInt(hex[4:6], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析蓝色分量失败: %w", err)
		}

		r = int(r64)
		g = int(g64)
		b = int(b64)
	} else if len(hex) == 8 {
		// 处理8位十六进制 (#RRGGBBAA) - 忽略alpha，使用参数
		r64, err := strconv.ParseInt(hex[0:2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析红色分量失败: %w", err)
		}
		g64, err := strconv.ParseInt(hex[2:4], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析绿色分量失败: %w", err)
		}
		b64, err := strconv.ParseInt(hex[4:6], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析蓝色分量失败: %w", err)
		}

		r = int(r64)
		g = int(g64)
		b = int(b64)
	}

	return &RGBA{R: r, G: g, B: b, A: alpha}, nil
}

// ColorFromName 从颜色名称创建RGBA对象
func ColorFromName(name string) (*RGBA, error) {
	switch strings.ToLower(name) {
	case "white":
		return ColorWhite, nil
	case "black":
		return ColorBlack, nil
	case "red":
		return ColorRed, nil
	case "green":
		return ColorGreen, nil
	case "blue":
		return ColorBlue, nil
	case "yellow":
		return ColorYellow, nil
	case "cyan":
		return ColorCyan, nil
	case "magenta":
		return ColorMagenta, nil
	case "gray", "grey":
		return ColorGray, nil
	case "darkgray", "darkgrey":
		return ColorDarkGray, nil
	case "lightgray", "lightgrey":
		return ColorLightGray, nil
	case "transparent":
		return ColorTransparent, nil
	case "darktheme":
		return ColorDarkTheme, nil
	case "lighttheme":
		return ColorLightTheme, nil
	case "sepiatheme":
		return ColorSepiaTheme, nil
	case "printwhite":
		return ColorPrintWhite, nil
	case "printcream":
		return ColorPrintCream, nil
	case "printblue":
		return ColorPrintBlue, nil
	default:
		return nil, fmt.Errorf("未知的颜色名称: %s", name)
	}
}

/*

// 示例1: 基本背景颜色测试
func exampleBasicBackgroundColorTesting() {
	// === 应用场景描述 ===
	// 场景: 基本背景颜色测试
	// 用途: 测试不同背景颜色对页面显示的影响
	// 优势: 统一页面背景颜色，便于视觉测试和对比
	// 典型工作流: 设置背景色 -> 截图 -> 分析视觉效果 -> 清除设置

	log.Println("基本背景颜色测试示例...")

	// 定义测试颜色列表
	testColors := []struct {
		name        string
		color       *RGBA
		description string
		useCases    []string
	}{
		{
			name:        "纯白背景",
			color:       ColorWhite,
			description: "标准白色背景，适用于大多数打印场景",
			useCases:    []string{"打印预览", "文档查看", "屏幕截图"},
		},
		{
			name:        "纯黑背景",
			color:       ColorBlack,
			description: "深色模式测试，验证暗色主题兼容性",
			useCases:    []string{"暗黑模式", "夜间模式", "高对比度"},
		},
		{
			name:        "米色背景",
			color:       ColorPrintCream,
			description: "护眼背景色，减少蓝光对眼睛的刺激",
			useCases:    []string{"长时间阅读", "护眼模式", "纸质效果"},
		},
		{
			name:        "淡蓝色背景",
			color:       ColorPrintBlue,
			description: "清新淡蓝色背景，提升页面可读性",
			useCases:    []string{"阅读模式", "文档浏览", "内容展示"},
		},
		{
			name:        "透明背景",
			color:       ColorTransparent,
			description: "透明背景，用于测试透明效果",
			useCases:    []string{"图层测试", "混合模式", "透明效果"},
		},
		{
			name:        "半透明白色",
			color:       ColorSemiTransparentWhite,
			description: "半透明背景，创建毛玻璃效果",
			useCases:    []string{"磨砂玻璃效果", "模态框", "遮罩层"},
		},
		{
			name:        "深色主题",
			color:       ColorDarkTheme,
			description: "深色主题背景，适用于暗色UI",
			useCases:    []string{"深色主题", "夜间模式", "低亮度环境"},
		},
		{
			name:        "棕褐色主题",
			color:       ColorSepiaTheme,
			description: "复古棕褐色背景，模仿纸张效果",
			useCases:    []string{"阅读模式", "复古主题", "纸张效果"},
		},
	}

	// 测试每种背景颜色
	for i, testColor := range testColors {
		log.Printf("\n=== 背景颜色测试 %d/%d: %s ===", i+1, len(testColors), testColor.name)
		log.Printf("描述: %s", testColor.description)
		log.Printf("颜色值: R=%d, G=%d, B=%d, A=%.2f",
			testColor.color.R, testColor.color.G, testColor.color.B, testColor.color.A)
		log.Printf("应用场景:")
		for _, useCase := range testColor.useCases {
			log.Printf("  - %s", useCase)
		}

		// 设置背景颜色
		log.Printf("设置背景颜色 '%s'...", testColor.name)
		response, err := CDPEmulationSetDefaultBackgroundColorOverride(testColor.color)
		if err != nil {
			log.Printf("设置背景颜色失败: %v", err)
			continue
		}

		log.Printf("✅ 背景颜色设置成功: %s", response)

		// 执行视觉测试
		visualTestResult := simulateVisualTesting(testColor.color)

		// 分析测试结果
		analyzeVisualTestResults(testColor, visualTestResult)

		// 短暂等待以便观察
		time.Sleep(500 * time.Millisecond)

		// 清除背景颜色
		log.Printf("清除背景颜色覆盖...")
		clearResponse, err := CDPEmulationSetDefaultBackgroundColorOverride(nil)
		if err != nil {
			log.Printf("清除背景颜色失败: %v", err)
		} else {
			log.Printf("✅ 背景颜色已清除: %s", clearResponse)
		}

		// 短暂延迟
		time.Sleep(200 * time.Millisecond)
	}
}

func simulateVisualTesting(color *RGBA) map[string]interface{} {
	metrics := make(map[string]interface{})

	// 计算亮度 (简单亮度计算)
	brightness := 0.299*float64(color.R) + 0.587*float64(color.G) + 0.114*float64(color.B)
	metrics["brightness"] = brightness

	// 计算对比度等级
	if brightness > 200 {
		metrics["contrastLevel"] = "亮色"
		metrics["textColorRecommend"] = "深色文字 (#000000 或 #333333)"
	} else if brightness > 128 {
		metrics["contrastLevel"] = "中等"
		metrics["textColorRecommend"] = "深色文字 (#000000 或 #222222)"
	} else if brightness > 64 {
		metrics["contrastLevel"] = "暗色"
		metrics["textColorRecommend"] = "浅色文字 (#FFFFFF 或 #F5F5F5)"
	} else {
		metrics["contrastLevel"] = "深色"
		metrics["textColorRecommend"] = "浅色文字 (#FFFFFF 或 #E0E0E0)"
	}

	// 透明度影响
	if color.A < 1.0 {
		metrics["transparencyEffect"] = fmt.Sprintf("半透明 (%.0f%%)", color.A*100)
		metrics["transparencyNote"] = "注意: 半透明背景会显示下层内容"
	} else {
		metrics["transparencyEffect"] = "不透明"
		metrics["transparencyNote"] = "完全覆盖底层内容"
	}

	// 颜色可读性评分
	readabilityScore := calculateReadabilityScore(color)
	metrics["readabilityScore"] = fmt.Sprintf("%.0f/100", readabilityScore)

	if readabilityScore >= 90 {
		metrics["readabilityLevel"] = "优秀"
	} else if readabilityScore >= 70 {
		metrics["readabilityLevel"] = "良好"
	} else if readabilityScore >= 50 {
		metrics["readabilityLevel"] = "一般"
	} else {
		metrics["readabilityLevel"] = "较差"
	}

	// 视觉舒适度评估
	comfortScore := calculateVisualComfortScore(color)
	metrics["comfortScore"] = fmt.Sprintf("%.0f/100", comfortScore)

	if comfortScore >= 90 {
		metrics["comfortLevel"] = "非常舒适"
	} else if comfortScore >= 70 {
		metrics["comfortLevel"] = "舒适"
	} else if comfortScore >= 50 {
		metrics["comfortLevel"] = "一般"
	} else {
		metrics["comfortLevel"] = "不舒适"
	}

	// 打印友好度评估
	printFriendliness := calculatePrintFriendliness(color)
	metrics["printFriendliness"] = fmt.Sprintf("%.0f/100", printFriendliness)

	if printFriendliness >= 90 {
		metrics["printLevel"] = "非常适合打印"
	} else if printFriendliness >= 70 {
		metrics["printLevel"] = "适合打印"
	} else if printFriendliness >= 50 {
		metrics["printLevel"] = "可打印"
	} else {
		metrics["printLevel"] = "不适合打印"
	}

	// 无障碍评估
	accessibilityScore := calculateAccessibilityScore(color)
	metrics["accessibilityScore"] = fmt.Sprintf("%.0f/100", accessibilityScore)

	if accessibilityScore >= 90 {
		metrics["accessibilityLevel"] = "无障碍友好"
	} else if accessibilityScore >= 70 {
		metrics["accessibilityLevel"] = "基本无障碍"
	} else if accessibilityScore >= 50 {
		metrics["accessibilityLevel"] = "需要改进"
	} else {
		metrics["accessibilityLevel"] = "无障碍问题"
	}

	return metrics
}

func calculateReadabilityScore(color *RGBA) float64 {
	// 基于亮度和对比度的可读性评分
	brightness := 0.299*float64(color.R) + 0.587*float64(color.G) + 0.114*float64(color.B)

	// 理想亮度范围 150-200
	readability := 100.0 - math.Abs(brightness-175)/175 * 100

	// 透明度惩罚
	if color.A < 0.9 {
		readability *= 0.8
	}

	// 确保在0-100范围内
	return math.Max(0, math.Min(100, readability))
}

func calculateVisualComfortScore(color *RGBA) float64 {
	// 视觉舒适度评分
	// 基于RGB分量的平衡
	r := float64(color.R)
	g := float64(color.G)
	b := float64(color.B)

	// 计算颜色平衡
	avg := (r + g + b) / 3
	balance := 100.0 - (math.Abs(r-avg) + math.Abs(g-avg) + math.Abs(b-avg)) / 3 / 255 * 100

	// 亮度舒适度
	brightness := 0.299*r + 0.587*g + 0.114*b
	brightnessScore := 100.0 - math.Abs(brightness-180)/180 * 100

	// 综合评分
	comfort := (balance*0.6 + brightnessScore*0.4)

	// 透明度影响
	if color.A < 1.0 {
		comfort *= 0.9
	}

	return math.Max(0, math.Min(100, comfort))
}

func calculatePrintFriendliness(color *RGBA) float64 {
	// 打印友好度评分
	// 基于亮度、对比度和颜色深度

	// 亮度适合打印 (中等亮度最佳)
	brightness := 0.299*float64(color.R) + 0.587*float64(color.G) + 0.114*float64(color.B)
	brightnessScore := 100.0 - math.Abs(brightness-200)/200 * 100

	// 墨水使用评分 (较浅颜色节省墨水)
	inkUsageScore := 100.0 - (float64(color.R) + float64(color.G) + float64(color.B)) / 3 / 255 * 100

	// 对比度评分
	var contrastScore float64
	if brightness > 200 {
		contrastScore = 90.0
	} else if brightness > 150 {
		contrastScore = 80.0
	} else if brightness > 100 {
		contrastScore = 60.0
	} else {
		contrastScore = 40.0
	}

	// 综合评分
	printScore := brightnessScore*0.4 + inkUsageScore*0.3 + contrastScore*0.3

	return math.Max(0, math.Min(100, printScore))
}

func calculateAccessibilityScore(color *RGBA) float64 {
	// 无障碍评分
	// 基于WCAG对比度标准

	brightness := 0.299*float64(color.R) + 0.587*float64(color.G) + 0.114*float64(color.B)

	// 基本对比度
	var contrastScore float64
	if brightness > 240 {
		// 非常亮的背景
		contrastScore = 95.0
	} else if brightness > 180 {
		// 亮背景
		contrastScore = 85.0
	} else if brightness > 120 {
		// 中等背景
		contrastScore = 70.0
	} else if brightness > 60 {
		// 暗背景
		contrastScore = 80.0
	} else {
		// 非常暗的背景
		contrastScore = 90.0
	}

	// 透明度无障碍
	transparencyScore := 100.0
	if color.A < 0.3 {
		transparencyScore = 60.0
	} else if color.A < 0.7 {
		transparencyScore = 80.0
	}

	// 颜色可区分性
	// 检查是否为纯色（非灰色）
	if math.Abs(float64(color.R)-float64(color.G)) < 30 &&
	   math.Abs(float64(color.G)-float64(color.B)) < 30 {
		// 接近灰色，对色盲友好
		colorScore := 100.0
	} else {
		// 彩色，检查对比度
		colorScore := 85.0
	}

	// 综合评分
	accessibility := contrastScore*0.5 + transparencyScore*0.3 + 100.0 * 0.2

	return math.Max(0, math.Min(100, accessibility))
}

func analyzeVisualTestResults(testColor struct {
	name        string
	color       *RGBA
	description string
	useCases    []string
}, metrics map[string]interface{}) {
	log.Printf("视觉测试结果分析:")

	log.Printf("  亮度: %.0f (等级: %v)",
		metrics["brightness"], metrics["contrastLevel"])

	log.Printf("  文字颜色建议: %v", metrics["textColorRecommend"])

	log.Printf("  透明度: %v", metrics["transparencyEffect"])
	if note, ok := metrics["transparencyNote"].(string); ok && note != "" {
		log.Printf("    注意: %s", note)
	}

	log.Printf("  可读性评分: %v (%v)",
		metrics["readabilityScore"], metrics["readabilityLevel"])

	log.Printf("  视觉舒适度: %v (%v)",
		metrics["comfortScore"], metrics["comfortLevel"])

	log.Printf("  打印友好度: %v (%v)",
		metrics["printFriendliness"], metrics["printLevel"])

	log.Printf("  无障碍评分: %v (%v)",
		metrics["accessibilityScore"], metrics["accessibilityLevel"])

	// 应用场景评估
	log.Printf("  适用场景评估:")
	for _, useCase := range testColor.useCases {
		suitability := evaluateColorSuitability(testColor.color, useCase)
		log.Printf("    - %s: %s", useCase, suitability)
	}

	// 优化建议
	log.Printf("  优化建议:")
	readabilityLevel, _ := metrics["readabilityLevel"].(string)
	comfortLevel, _ := metrics["comfortLevel"].(string)

	if readabilityLevel == "较差" || comfortLevel == "不舒适" {
		log.Printf("    ⚠ 考虑调整颜色以提高可读性和舒适度")
		if testColor.color.A < 0.5 {
			log.Printf("    💡 提高透明度以增强可读性")
		}
		if readabilityLevel == "较差" {
			log.Printf("    💡 建议亮度调整到150-200范围")
		}
	}

	if printScore, ok := metrics["printFriendliness"].(string); ok {
		if strings.Contains(printScore, "不适合打印") {
			log.Printf("    🖨️ 此颜色可能不适合打印，考虑使用更浅的颜色")
		}
	}

	if accessibilityLevel, ok := metrics["accessibilityLevel"].(string); ok {
		if strings.Contains(accessibilityLevel, "问题") {
			log.Printf("    ♿ 无障碍性需要改进，请参考WCAG标准")
		}
	}
}

func evaluateColorSuitability(color *RGBA, useCase string) string {
	brightness := 0.299*float64(color.R) + 0.587*float64(color.G) + 0.114*float64(color.B)
	alpha := color.A

	switch useCase {
	case "打印预览":
		if brightness > 200 && alpha >= 0.9 {
			return "✅ 非常适合"
		} else if brightness > 150 && alpha >= 0.9 {
			return "✅ 适合"
		} else {
			return "⚠ 可能需要调整"
		}

	case "暗黑模式", "夜间模式":
		if brightness < 100 && alpha >= 0.9 {
			return "✅ 非常适合"
		} else if brightness < 150 && alpha >= 0.9 {
			return "✅ 适合"
		} else {
			return "⚠ 可能需要调整"
		}

	case "阅读模式", "长时间阅读":
		if brightness > 180 && brightness < 220 && alpha >= 0.9 {
			return "✅ 非常适合"
		} else if brightness > 150 && brightness < 240 && alpha >= 0.8 {
			return "✅ 适合"
		} else {
			return "⚠ 可能需要调整"
		}

	case "高对比度":
		if (brightness < 50 || brightness > 230) && alpha >= 0.9 {
			return "✅ 非常适合"
		} else if (brightness < 100 || brightness > 200) && alpha >= 0.9 {
			return "✅ 适合"
		} else {
			return "⚠ 可能需要调整"
		}

	case "透明效果", "图层测试":
		if alpha < 1.0 {
			return "✅ 非常适合"
		} else {
			return "⚠ 不是透明效果"
		}

	case "磨砂玻璃效果", "模态框":
		if alpha < 1.0 && alpha > 0.3 {
			return "✅ 非常适合"
		} else {
			return "⚠ 透明度可能不合适"
		}

	default:
		return "🔍 需进一步测试"
	}
}

*/

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
func CDPEmulationSetDeviceMetricsOverride(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 必需参数验证
	width, hasWidth := params["width"].(int)
	height, hasHeight := params["height"].(int)
	deviceScaleFactor, hasScaleFactor := params["deviceScaleFactor"].(float64)
	_, hasMobile := params["mobile"].(bool)

	if !hasWidth || !hasHeight || !hasScaleFactor || !hasMobile {
		return "", fmt.Errorf("缺少必需参数: width, height, deviceScaleFactor, mobile")
	}

	// 验证参数范围
	if width < 0 || width > 10000000 {
		return "", fmt.Errorf("宽度必须在0-10000000范围内: %d", width)
	}
	if height < 0 || height > 10000000 {
		return "", fmt.Errorf("高度必须在0-10000000范围内: %d", height)
	}
	if deviceScaleFactor < 0 {
		return "", fmt.Errorf("设备比例因子不能为负: %f", deviceScaleFactor)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setDeviceMetricsOverride",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setDeviceMetricsOverride 请求失败: %w", err)
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
func CDPEmulationSetEmulatedMedia(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setEmulatedMedia",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setEmulatedMedia 请求失败: %w", err)
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

				log.Printf("✅ 媒体模拟设置成功")
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setEmulatedMedia 请求超时")
		}
	}
}

// MediaFeature 媒体特性
type MediaFeature struct {
	Name  string `json:"name"`  // 媒体特性名称
	Value string `json:"value"` // 媒体特性值
}

// MediaFeatureType 媒体特性类型常量
const (
	// 颜色相关
	ColorSchemeDark   = "(prefers-color-scheme: dark)"
	ColorSchemeLight  = "(prefers-color-scheme: light)"
	ColorSchemeNoPref = "(prefers-color-scheme: no-preference)"

	// 对比度相关
	ContrastMore   = "(prefers-contrast: more)"
	ContrastLess   = "(prefers-contrast: less)"
	ContrastNoPref = "(prefers-contrast: no-preference)"

	// 运动相关
	ReducedMotionReduce = "(prefers-reduced-motion: reduce)"
	ReducedMotionNoPref = "(prefers-reduced-motion: no-preference)"

	// 透明度相关
	ReducedTransparencyReduce = "(prefers-reduced-transparency: reduce)"
	ReducedTransparencyNoPref = "(prefers-reduced-transparency: no-preference)"

	// 数据保存相关
	DataSaveOn  = "(prefers-reduced-data: reduce)"
	DataSaveOff = "(prefers-reduced-data: no-preference)"

	// 显示更新频率
	UpdateFrequencyFast = "(update: fast)"
	UpdateFrequencySlow = "(update: slow)"

	// 色彩空间
	ColorGamutSRGB    = "(color-gamut: srgb)"
	ColorGamutP3      = "(color-gamut: p3)"
	ColorGamutRec2020 = "(color-gamut: rec2020)"

	// 屏幕分辨率
	Resolution1x = "(resolution: 96dpi)"
	Resolution2x = "(resolution: 192dpi)"

	// 设备像素比
	DevicePixelRatio1 = "(-webkit-device-pixel-ratio: 1)"
	DevicePixelRatio2 = "(-webkit-device-pixel-ratio: 2)"
	DevicePixelRatio3 = "(-webkit-device-pixel-ratio: 3)"

	// 设备方向
	OrientationPortrait  = "(orientation: portrait)"
	OrientationLandscape = "(orientation: landscape)"
)

// MediaType 媒体类型常量
const (
	MediaTypeScreen     = "screen"     // 屏幕
	MediaTypePrint      = "print"      // 打印
	MediaTypeSpeech     = "speech"     // 语音合成
	MediaTypeAll        = "all"        // 所有设备
	MediaTypeBraille    = "braille"    // 盲文设备
	MediaTypeEmbossed   = "embossed"   // 盲文打印
	MediaTypeHandheld   = "handheld"   // 手持设备
	MediaTypeProjection = "projection" // 投影
	MediaTypeTTY        = "tty"        // 电传打字机
	MediaTypeTV         = "tv"         // 电视
	MediaTypeNone       = ""           // 禁用媒体模拟
)

// MediaSimulationConfig 媒体模拟配置
type MediaSimulationConfig struct {
	Media       string          `json:"media,omitempty"`       // 媒体类型
	Features    []*MediaFeature `json:"features,omitempty"`    // 媒体特性列表
	Description string          `json:"description,omitempty"` // 配置描述
	Scenario    string          `json:"scenario,omitempty"`    // 测试场景
}

// 预定义的媒体模拟配置
var (
	// 暗色模式
	DarkModeConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "prefers-color-scheme", Value: "dark"},
		},
		Description: "暗色模式（深色主题）",
		Scenario:    "测试暗色主题适配",
	}

	// 亮色模式
	LightModeConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "prefers-color-scheme", Value: "light"},
		},
		Description: "亮色模式（浅色主题）",
		Scenario:    "测试亮色主题适配",
	}

	// 打印模式
	PrintConfig = &MediaSimulationConfig{
		Media:       MediaTypePrint,
		Description: "打印样式",
		Scenario:    "测试打印页面样式",
	}

	// 减少动画模式
	ReducedMotionConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "prefers-reduced-motion", Value: "reduce"},
		},
		Description: "减少动画模式",
		Scenario:    "测试动画减少模式下的用户体验",
	}

	// 高对比度模式
	HighContrastConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "prefers-contrast", Value: "more"},
		},
		Description: "高对比度模式",
		Scenario:    "测试视觉障碍用户的高对比度适配",
	}

	// 低数据模式
	ReducedDataConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "prefers-reduced-data", Value: "reduce"},
		},
		Description: "低数据模式",
		Scenario:    "测试网络条件差的环境下数据优化",
	}

	// 语音合成模式
	SpeechConfig = &MediaSimulationConfig{
		Media:       MediaTypeSpeech,
		Description: "语音合成设备",
		Scenario:    "测试屏幕阅读器和语音合成器的可访问性",
	}

	// 便携设备模式
	HandheldConfig = &MediaSimulationConfig{
		Media:       MediaTypeHandheld,
		Description: "便携设备（手机、平板等）",
		Scenario:    "测试移动设备优化",
	}

	// 电视模式
	TVConfig = &MediaSimulationConfig{
		Media:       MediaTypeTV,
		Description: "电视设备",
		Scenario:    "测试电视大屏界面和遥控器交互",
	}

	// 多特性组合：深色模式+高对比度+减少动画
	AccessibilityConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "prefers-color-scheme", Value: "dark"},
			{Name: "prefers-contrast", Value: "more"},
			{Name: "prefers-reduced-motion", Value: "reduce"},
			{Name: "prefers-reduced-transparency", Value: "reduce"},
		},
		Description: "完整无障碍访问配置",
		Scenario:    "测试多个无障碍特性的组合效果",
	}

	// 色彩空间测试
	WideGamutConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "color-gamut", Value: "p3"},
		},
		Description: "广色域屏幕（P3）",
		Scenario:    "测试广色域设备的色彩表现",
	}

	// 高DPI屏幕
	HighDPIConfig = &MediaSimulationConfig{
		Media: MediaTypeScreen,
		Features: []*MediaFeature{
			{Name: "resolution", Value: "192dpi"},
			{Name: "-webkit-device-pixel-ratio", Value: "2"},
		},
		Description: "高DPI屏幕（视网膜屏）",
		Scenario:    "测试高分辨率设备的图片和字体渲染",
	}
)

/*

// 示例1: 基本媒体模拟
func exampleBasicMediaSimulation() {
	log.Println("基本媒体模拟示例...")

	// 定义媒体测试场景
	mediaTests := []struct {
		name         string
		mediaConfig  *MediaSimulationConfig
		cssExample   string
		jsDetection  string
		testScenarios []string
	}{
		{
			name:        "暗色模式模拟",
			mediaConfig: DarkModeConfig,
			cssExample: `@media (prefers-color-scheme: dark) {
  body { background-color: #121212; color: #ffffff; }
}`,
			jsDetection: `if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
  console.log('设备使用暗色模式');
}`,
			testScenarios: []string{
				"暗色主题切换",
				"颜色对比度检查",
				"夜间模式优化",
				"深色背景上的文字可读性",
			},
		},
		{
			name:        "打印样式模拟",
			mediaConfig: PrintConfig,
			cssExample: `@media print {
  .no-print { display: none; }
  body { color: #000; }
}`,
			jsDetection: `if (window.matchMedia('print').matches) {
  console.log('打印模式激活');
}`,
			testScenarios: []string{
				"打印页面布局",
				"隐藏不必要元素",
				"打印链接处理",
				"页眉页脚设置",
			},
		},
		{
			name:        "减少动画模式",
			mediaConfig: ReducedMotionConfig,
			cssExample: `@media (prefers-reduced-motion: reduce) {
  * { animation-duration: 0.01ms !important; }
  * { animation-iteration-count: 1 !important; }
}`,
			jsDetection: `if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
  console.log('用户偏好减少动画');
}`,
			testScenarios: []string{
				"动画减少优化",
				"运动敏感用户测试",
				"替代动画方案",
				"性能优化",
			},
		},
		{
			name:        "高对比度模式",
			mediaConfig: HighContrastConfig,
			cssExample: `@media (prefers-contrast: more) {
  .low-contrast { border: 2px solid; }
  a { text-decoration: underline; }
}`,
			jsDetection: `if (window.matchMedia('(prefers-contrast: more)').matches) {
  console.log('高对比度模式激活');
}`,
			testScenarios: []string{
				"视觉障碍用户支持",
				"颜色对比度检查",
				"可读性优化",
				"WCAG合规性测试",
			},
		},
		{
			name:        "完整无障碍配置",
			mediaConfig: AccessibilityConfig,
			cssExample:
			@media (prefers-color-scheme: dark) and
			(prefers-contrast: more) and
			(prefers-reduced-motion: reduce) {

	 }`,
				 jsDetection: `const isAccessibility =
	   window.matchMedia('(prefers-color-scheme: dark)').matches &&
	   window.matchMedia('(prefers-contrast: more)').matches &&
	   window.matchMedia('(prefers-reduced-motion: reduce)').matches;`,
				 testScenarios: []string{
					 "综合无障碍测试",
					 "多特性组合优化",
					 "可访问性审计",
					 "用户体验测试",
				 },
			 },
			 {
				 name:        "语音合成设备",
				 mediaConfig: SpeechConfig,
				 cssExample: `@media speech {
	   .speech-only { display: block; }
	   .visual-only { display: none; }
	 }`,
				 jsDetection: `if (window.matchMedia('speech').matches) {
	   console.log('语音合成设备');
	 }`,
				 testScenarios: []string{
					 "屏幕阅读器兼容性",
					 "语音导航测试",
					 "可访问性语音优化",
					 "听觉反馈测试",
				 },
			 },
			 {
				 name:        "便携设备模式",
				 mediaConfig: HandheldConfig,
				 cssExample: `@media handheld {

	   font-size: 16px;
	 }`,
				 jsDetection: `if (window.matchMedia('handheld').matches) {
	   console.log('手持设备检测');
	 }`,
				 testScenarios: []string{
					 "移动优先设计测试",
					 "触摸交互优化",
					 "移动端性能测试",
					 "小屏设备适配",
				 },
			 },
			 {
				 name:        "电视设备模式",
				 mediaConfig: TVConfig,
				 cssExample: `@media tv {

	   .tv-optimized { font-size: 24px; }
	 }`,
				 jsDetection: `if (window.matchMedia('tv').matches) {
	   console.log('电视设备检测');
	 }`,
				 testScenarios: []string{
					 "大屏界面设计",
					 "遥控器导航测试",
					 "10英尺体验优化",
					 "电视应用测试",
				 },
			 },
			 {
				 name:        "低数据模式",
				 mediaConfig: ReducedDataConfig,
				 cssExample: `@media (prefers-reduced-data: reduce) {
	   .high-res-image { display: none; }
	   .low-res-image { display: block; }
	 }`,
				 jsDetection: `if (window.matchMedia('(prefers-reduced-data: reduce)').matches) {
	   console.log('低数据模式激活');
	 }`,
				 testScenarios: []string{
					 "数据优化测试",
					 "网络条件差的环境",
					 "图片懒加载优化",
					 "资源加载策略",
				 },
			 },
			 {
				 name:        "高DPI屏幕",
				 mediaConfig: HighDPIConfig,
				 cssExample: `@media (-webkit-min-device-pixel-ratio: 2),
			(min-resolution: 192dpi) {
	   .logo { background-image: url('logo@2x.png'); }
	 }`,
				 jsDetection: `const isRetina =
	   window.devicePixelRatio >= 2 ||
	   window.matchMedia('(min-resolution: 192dpi)').matches;`,
				 testScenarios: []string{
					 "视网膜显示屏优化",
					 "高分辨率图片",
					 "字体抗锯齿",
					 "高DPI适配",
				 },
			 },
		 }

		 // 执行媒体模拟测试
		 for i, test := range mediaTests {
			 log.Printf("\n=== 媒体模拟测试 %d/%d: %s ===", i+1, len(mediaTests), test.name)
			 log.Printf("描述: %s", test.mediaConfig.Description)
			 log.Printf("测试场景: %s", test.mediaConfig.Scenario)

			 log.Printf("媒体类型: %s", test.mediaConfig.Media)
			 if len(test.mediaConfig.Features) > 0 {
				 log.Printf("媒体特性:")
				 for _, feature := range test.mediaConfig.Features {
					 log.Printf("  - %s: %s", feature.Name, feature.Value)
				 }
			 }

			 log.Printf("CSS示例:")
			 log.Printf("  %s", test.cssExample)

			 log.Printf("JavaScript检测:")
			 log.Printf("  %s", test.jsDetection)

			 log.Printf("测试场景:")
			 for _, scenario := range test.testScenarios {
				 log.Printf("  - %s", scenario)
			 }

			 // 构建参数
			 params := make(map[string]interface{})

			 // 设置媒体类型
			 params["media"] = test.mediaConfig.Media

			 // 设置媒体特性
			 if len(test.mediaConfig.Features) > 0 {
				 features := make([]map[string]interface{}, len(test.mediaConfig.Features))
				 for i, feature := range test.mediaConfig.Features {
					 features[i] = map[string]interface{}{
						 "name":  feature.Name,
						 "value": feature.Value,
					 }
				 }
				 params["features"] = features
			 }

			 log.Printf("设置媒体模拟...")
			 response, err := CDPEmulationSetEmulatedMedia(params)
			 if err != nil {
				 log.Printf("❌ 媒体模拟失败: %v", err)
				 continue
			 }

			 log.Printf("✅ 媒体模拟成功")

			 // 执行验证测试
			 log.Printf("执行验证测试...")
			 testResults := validateMediaSimulation(test.mediaConfig, test.testScenarios)

			 // 提供优化建议
			 provideMediaOptimizationAdvice(test.mediaConfig, testResults)

			 // 清除媒体模拟
			 log.Printf("清除媒体模拟...")
			 clearParams := map[string]interface{}{
				 "media": "",
			 }
			 clearResponse, err := CDPEmulationSetEmulatedMedia(clearParams)
			 if err != nil {
				 log.Printf("⚠ 清除媒体模拟失败: %v", err)
			 } else {
				 log.Printf("✅ 媒体模拟已清除")
			 }

			 // 短暂延迟
			 time.Sleep(300 * time.Millisecond)
		 }
	 }


*/

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
func CDPEmulationSetEmulatedOSTextScale(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 验证参数
	textScaleFactor, ok := params["textScaleFactor"].(float64)
	if !ok {
		return "", fmt.Errorf("缺少必需的textScaleFactor参数")
	}

	// 验证缩放因子的合理性
	if textScaleFactor < 0.5 || textScaleFactor > 5.0 {
		return "", fmt.Errorf("文本缩放因子必须在0.5-5.0范围内: %f", textScaleFactor)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setEmulatedOSTextScale",
		"params": {
			"textScaleFactor": %f
		}
	}`, reqID, textScaleFactor)

	// 发送请求
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

				log.Printf("✅ 操作系统文本缩放设置成功: %.2fx", textScaleFactor)
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
func CDPEmulationSetEmulatedVisionDeficiency(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setEmulatedVisionDeficiency",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
func CDPEmulationSetGeolocationOverride(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setGeolocationOverride",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
func CDPEmulationSetIdleOverride(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setIdleOverride",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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
func CDPEmulationSetScriptExecutionDisabled(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 验证参数
	value, ok := params["value"].(bool)
	if !ok {
		return "", fmt.Errorf("缺少必需的value参数（应为布尔值）")
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setScriptExecutionDisabled",
		"params": {
			"value": %t
		}
	}`, reqID, value)

	// 发送请求
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

				status := "禁用"
				if !value {
					status = "启用"
				}
				log.Printf("✅ 脚本执行已%s", status)
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
func CDPEmulationSetTimezoneOverride(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 验证参数
	timezoneId, ok := params["timezoneId"].(string)
	if !ok {
		return "", fmt.Errorf("缺少必需的timezoneId参数（应为字符串）")
	}

	// 验证时区ID格式
	if !isValidTimezone(timezoneId) {
		return "", fmt.Errorf("无效的时区ID: %s，请使用IANA时区数据库格式", timezoneId)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setTimezoneOverride",
		"params": {
			"timezoneId": "%s"
		}
	}`, reqID, timezoneId)

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

				log.Printf("✅ 时区已覆盖为: %s", timezoneId)
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setTimezoneOverride 请求超时")
		}
	}
}

// isValidTimezone 验证时区ID是否有效
// 检查时区是否符合IANA时区数据库格式
// 支持格式: 地区/城市 或 Etc/GMT±[偏移量]
func isValidTimezone(timezoneID string) bool {
	if timezoneID == "" {
		return false
	}

	// 移除可能的空格
	timezoneID = strings.TrimSpace(timezoneID)

	// 基本格式检查
	// IANA时区格式通常是: 地区/城市
	// 或者特殊的 Etc/GMT±[偏移量]
	parts := strings.Split(timezoneID, "/")
	if len(parts) < 1 || len(parts) > 2 {
		return false
	}

	// 检查常见的有效时区模式
	validPatterns := []*regexp.Regexp{
		// 1. 标准地区/城市格式: America/New_York, Europe/London
		regexp.MustCompile(`^[A-Za-z_]+/[A-Za-z_]+(?:/[A-Za-z_]+)?$`),
		// 2. Etc/GMT格式: Etc/GMT, Etc/GMT+1, Etc/GMT-1
		regexp.MustCompile(`^Etc/GMT(?:[+-]\d{1,2})?$`),
		// 3. GMT格式: GMT, GMT+1, GMT-1
		regexp.MustCompile(`^GMT(?:[+-]\d{1,2})?$`),
		// 4. UTC格式: UTC, UTC+1, UTC-1
		regexp.MustCompile(`^UTC(?:[+-]\d{1,2})?$`),
		// 5. 时区缩写: EST, PST, CET
		regexp.MustCompile(`^[A-Z]{3,4}$`),
	}

	// 检查是否匹配任何有效模式
	matchesAnyPattern := false
	for _, pattern := range validPatterns {
		if pattern.MatchString(timezoneID) {
			matchesAnyPattern = true
			break
		}
	}

	if !matchesAnyPattern {
		return false
	}

	// 尝试加载时区来验证
	_, err := time.LoadLocation(timezoneID)
	return err == nil
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
func CDPEmulationSetTouchEmulationEnabled(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 验证必填参数
	enabled, ok := params["enabled"].(bool)
	if !ok {
		return "", fmt.Errorf("缺少必需的enabled参数（应为布尔值）")
	}

	// 构建参数
	messageParams := map[string]interface{}{
		"enabled": enabled,
	}

	// 添加可选参数
	if configuration, ok := params["configuration"]; ok {
		messageParams["configuration"] = configuration
	}
	if maxTouchPoints, ok := params["maxTouchPoints"]; ok {
		messageParams["maxTouchPoints"] = maxTouchPoints
	}

	// 构建消息
	paramsJSON, err := json.Marshal(messageParams)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setTouchEmulationEnabled",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setTouchEmulationEnabled 请求失败: %w", err)
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

				status := "启用"
				if !enabled {
					status = "禁用"
				}
				log.Printf("✅ 触摸模拟已%s", status)
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
func CDPEmulationSetUserAgentOverride(params map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 验证必填参数
	userAgent, ok := params["userAgent"].(string)
	if !ok {
		return "", fmt.Errorf("缺少必需的userAgent参数（应为字符串）")
	}

	// 构建参数
	messageParams := map[string]interface{}{
		"userAgent": userAgent,
	}

	// 添加可选参数
	if acceptLanguage, ok := params["acceptLanguage"]; ok {
		messageParams["acceptLanguage"] = acceptLanguage
	}
	if platform, ok := params["platform"]; ok {
		messageParams["platform"] = platform
	}
	if metadata, ok := params["userAgentMetadata"]; ok {
		messageParams["userAgentMetadata"] = metadata
	}

	// 构建消息
	paramsJSON, err := json.Marshal(messageParams)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Emulation.setUserAgentOverride",
		"params": %s
	}`, reqID, paramsJSON)

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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

				log.Printf("✅ 用户代理已覆盖为: %s", userAgent)
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setUserAgentOverride 请求超时")
		}
	}
}
