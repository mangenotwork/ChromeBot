package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Profiler.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 停止性能采集: 停止浏览器的JS执行性能分析器，不再收集性能数据
// 2. 资源释放: 关闭性能分析器，释放浏览器占用的内存和CPU资源
// 3. 测试流程收尾: 自动化测试中性能采集完成后的收尾操作
// 4. 调试结束: 代码性能调试完成后关闭分析器
// 5. 性能恢复: 让浏览器JS执行恢复到正常无监控状态
// 6. 多场景切换: 切换不同调试任务时关闭上一个性能采集任务

// CDPProfilerDisable 禁用性能分析器，停止收集JS执行性能数据
func CDPProfilerDisable() (string, error) {
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
		"method": "Profiler.disable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.disable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.disable 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化测试中性能采集完成后关闭分析器
func ExampleProfilerDisable_AutoTest() {
	// 1. 先启用性能采集
	// CDPProfilerEnable()

	// 2. 执行业务逻辑/测试用例
	// RunBusinessLogic()

	// 3. 采集完成后禁用分析器
	resp, err := CDPProfilerDisable()
	if err != nil {
		log.Fatalf("禁用性能分析器失败: %v", err)
	}
	log.Printf("禁用成功，响应: %s", resp)
}

// 场景2：调试结束后释放资源，恢复浏览器正常状态
func ExampleProfilerDisable_DebugEnd() {
	// 调试完成后关闭性能采集
	if _, err := CDPProfilerDisable(); err != nil {
		log.Printf("关闭性能分析器异常: %v", err)
	} else {
		log.Println("性能分析器已关闭，资源已释放")
	}
}

*/

// -----------------------------------------------  Profiler.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 启动性能采集: 开启浏览器JS执行性能分析器，开始收集性能数据
// 2. 自动化测试前置: 自动化测试前启用性能监控，为后续采集做准备
// 3. 代码调试启动: 调试JS性能问题时，开启性能分析功能
// 4. 性能监控初始化: 页面性能监控流程的初始化步骤
// 5. 多轮性能测试: 每次性能测试前重新启用分析器
// 6. 前端性能诊断: 诊断页面卡顿、JS执行耗时前启用采集

// CDPProfilerEnable 启用性能分析器，开始收集JavaScript执行性能数据
func CDPProfilerEnable() (string, error) {
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
		"method": "Profiler.enable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.enable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.enable 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化测试前启用性能采集
func ExampleProfilerEnable_AutoTestPrepare() {
	// 打开目标页面
	// CDPPageNavigate("https://your-test-page.com")

	// 启用性能分析器
	resp, err := CDPProfilerEnable()
	if err != nil {
		log.Fatalf("启用性能分析器失败: %v", err)
	}
	log.Printf("性能分析器启用成功，响应: %s", resp)

	// 后续执行性能采集、业务测试逻辑
}

// 场景2：JS性能调试前启动分析器
func ExampleProfilerEnable_DebugStart() {
	// 诊断页面JS性能问题前启用
	if _, err := CDPProfilerEnable(); err != nil {
		log.Printf("启动性能分析失败: %v", err)
		return
	}
	log.Println("已启动JS性能采集，可开始调试性能问题")

	// 执行需要调试的JS操作
}

*/

// -----------------------------------------------  Profiler.getBestEffortCoverage  -----------------------------------------------
// === 应用场景 ===
// 1. JS代码覆盖率采集：获取页面JS代码的执行覆盖数据（尽力而为模式）
// 2. 前端测试覆盖率：自动化测试中收集测试用例对JS代码的覆盖情况
// 3. 无用代码检测：分析页面未执行的JS代码，用于代码精简和优化
// 4. 性能优化分析：结合覆盖率数据定位未使用代码，减少资源加载
// 5. 持续集成监控：CI/CD流程中自动采集前端代码覆盖率报告
// 6. 调试辅助：定位哪些JS逻辑未被执行，辅助问题排查

// CDPProfilerGetBestEffortCoverage 获取JS代码尽力而为模式的覆盖率数据
func CDPProfilerGetBestEffortCoverage() (string, error) {
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
		"method": "Profiler.getBestEffortCoverage"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.getBestEffortCoverage 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.getBestEffortCoverage 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化测试完成后采集JS代码覆盖率
func ExampleProfilerGetBestEffortCoverage_AutoTest() {
	// 1. 启用性能分析器
	_, err := CDPProfilerEnable()
	if err != nil {
		log.Fatalf("启用分析器失败: %v", err)
	}

	// 2. 执行自动化测试业务逻辑
	// RunTestCases()

	// 3. 采集尽力而为的覆盖率数据
	coverage, err := CDPProfilerGetBestEffortCoverage()
	if err != nil {
		log.Fatalf("采集覆盖率失败: %v", err)
	}
	log.Printf("JS代码覆盖率数据: %s", coverage)

	// 4. 关闭分析器
	_, _ = CDPProfilerDisable()
}

// 场景2：检测页面未使用的JS代码（无用代码分析）
func ExampleProfilerGetBestEffortCoverage_CodeOptimize() {
	// 页面加载完成后采集覆盖率
	if _, err := CDPProfilerEnable(); err != nil {
		log.Printf("启用失败: %v", err)
		return
	}

	// 模拟用户操作触发JS执行
	// SimulateUserAction()

	// 获取覆盖率结果用于分析死代码
	result, err := CDPProfilerGetBestEffortCoverage()
	if err != nil {
		log.Printf("获取覆盖率失败: %v", err)
		return
	}
	// 可将结果保存为文件用于分析
	// utils.SaveToFile("coverage.json", result)
	log.Println("已获取代码覆盖率，可用于分析未使用代码")
	_, _ = CDPProfilerDisable()
}

*/

// -----------------------------------------------  Profiler.setSamplingInterval  -----------------------------------------------
// === 应用场景 ===
// 1. 性能采集精度调整：自定义JS性能分析器的采样间隔，平衡精度与性能开销
// 2. 高精度调试：设置短间隔（如100μs）用于精细分析JS执行瓶颈
// 3. 低开销监控：设置长间隔（如10000μs）减少性能采集对业务的影响
// 4. 自动化测试配置：根据测试场景动态配置采样频率
// 5. 长时间监控：长间隔用于持续性能监控，避免内存占用过高
// 6. 问题复现调试：针对复杂性能问题调整采样精度，精准定位耗时逻辑

// CDPProfilerSetSamplingInterval 设置JS性能分析器的采样间隔
// 参数 interval：采样间隔，单位**微秒(μs)**，推荐范围：100(高精度) ~ 10000(低开销)
func CDPProfilerSetSamplingInterval(interval int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Profiler.setSamplingInterval",
		"params": {
			"interval": %d
		}
	}`, reqID, interval)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.setSamplingInterval 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.setSamplingInterval 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：高精度性能调试（设置100微秒间隔，采集精细数据）
func ExampleProfilerSetSamplingInterval_HighPrecision() {
	// 启用性能分析器
	_, err := CDPProfilerEnable()
	if err != nil {
		log.Fatalf("启用分析器失败: %v", err)
	}

	// 设置高精度采样间隔：100μs，适合定位细微性能问题
	resp, err := CDPProfilerSetSamplingInterval(100)
	if err != nil {
		log.Fatalf("设置采样间隔失败: %v", err)
	}
	log.Printf("高精度采样模式已启用，响应: %s", resp)

	// 执行需要调试的JS逻辑
	// RunDebugJS()
}

// 场景2：低开销长期监控（设置10000微秒间隔，减少性能损耗）
func ExampleProfilerSetSamplingInterval_LowOverhead() {
	// 启用后设置低开销采样间隔：10000μs，适合长时间运行监控
	if _, err := CDPProfilerEnable(); err != nil {
		log.Printf("启用失败: %v", err)
		return
	}

	// 低频率采样，对业务性能影响小
	_, err := CDPProfilerSetSamplingInterval(10000)
	if err != nil {
		log.Printf("设置低开销模式失败: %v", err)
		return
	}
	log.Println("已切换至低开销性能监控模式")
}

*/

// -----------------------------------------------  Profiler.start  -----------------------------------------------
// === 应用场景 ===
// 1. 启动JS性能采样：开始采集JavaScript执行的CPU性能采样数据
// 2. 自动化性能测试：测试用例执行时启动性能采集，记录函数执行耗时
// 3. 页面卡顿分析：页面交互时启动采集，定位导致卡顿的JS函数
// 4. 接口响应优化：API请求前后启动采样，分析请求处理中的JS耗时逻辑
// 5. 前端性能诊断：用户操作触发前启动，捕获真实场景的性能数据
// 6. 持续性能监控：业务流程执行时启动，生成性能报告

// CDPProfilerStart 启动JavaScript CPU性能采样分析
func CDPProfilerStart() (string, error) {
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
		"method": "Profiler.start"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.start 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.start 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化性能测试 - 执行业务逻辑前启动CPU采样
func ExampleProfilerStart_AutoPerformanceTest() {
	// 1. 启用性能分析器
	_, err := CDPProfilerEnable()
	if err != nil {
		log.Fatalf("启用性能分析器失败: %v", err)
	}

	// 2. 配置采样间隔（可选）
	_, _ = CDPProfilerSetSamplingInterval(500)

	// 3. 启动性能采样
	resp, err := CDPProfilerStart()
	if err != nil {
		log.Fatalf("启动性能采样失败: %v", err)
	}
	log.Printf("启动CPU性能采样成功: %s", resp)

	// 4. 执行需要测试的业务逻辑
	// RunBusinessTestLogic()
}

// 场景2：页面交互卡顿调试 - 按钮点击/页面切换前启动采集
func ExampleProfilerStart_PageLagDebug() {
	// 调试页面卡顿问题
	if _, err := CDPProfilerEnable(); err != nil {
		log.Printf("启用分析器失败: %v", err)
		return
	}

	// 启动采样，捕获卡顿时刻的JS执行栈
	if _, err := CDPProfilerStart(); err != nil {
		log.Printf("启动采样失败: %v", err)
		return
	}
	log.Println("已启动CPU采样，即将触发交互操作...")

	// 触发可能导致卡顿的操作
	// TriggerPageInteraction()
}

*/

// -----------------------------------------------  Profiler.startPreciseCoverage  -----------------------------------------------
// === 应用场景 ===
// 1. 精准JS代码覆盖率采集：启动高精度、精确的代码覆盖率统计（无丢数据、高准确度）
// 2. 正式测试覆盖率报告：自动化测试生成精准覆盖率报表，用于质量门禁
// 3. 生产环境轻量监控：支持低开销精确覆盖，适合线上非侵入式采集
// 4. 死代码精准清理：基于精确覆盖数据，安全删除未执行JS代码
// 5. 单测/集成测覆盖：单元测试、E2E测试中获取100%可靠覆盖率数据
// 6. 合规/质量审计：需要精准代码覆盖数据的合规检查

// CDPProfilerStartPreciseCoverage 启动精确JS代码覆盖率采集
// 参数 enableCallCount：是否启用函数调用次数统计
// 参数 enableDetailedCoverage：是否启用详细代码块级覆盖率（行/语句级）
func CDPProfilerStartPreciseCoverage(enableCallCount bool, enableDetailedCoverage bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的精确覆盖率启动请求
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Profiler.startPreciseCoverage",
		"params": {
			"enableCallCount": %t,
			"enableDetailedCoverage": %t
		}
	}`, reqID, enableCallCount, enableDetailedCoverage)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.startPreciseCoverage 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.startPreciseCoverage 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化测试生成高精度覆盖率报告（开启调用计数+详细覆盖）
func ExampleProfilerStartPreciseCoverage_AutoTestReport() {
	// 启用分析器
	_, err := CDPProfilerEnable()
	if err != nil {
		log.Fatalf("启用分析器失败: %v", err)
	}

	// 启动精确覆盖率：开启调用次数 + 开启详细语句覆盖
	resp, err := CDPProfilerStartPreciseCoverage(true, true)
	if err != nil {
		log.Fatalf("启动精确覆盖率失败: %v", err)
	}
	log.Printf("启动精确覆盖率成功: %s", resp)

	// 执行测试用例
	// RunAllTestCases()
}

// 场景2：轻量生产环境覆盖采集（关闭详细数据，降低性能开销）
func ExampleProfilerStartPreciseCoverage_ProductionMonitor() {
	_, _ = CDPProfilerEnable()

	// 生产环境：仅开启调用统计，不采集详细覆盖，低开销
	resp, err := CDPProfilerStartPreciseCoverage(true, false)
	if err != nil {
		log.Printf("启动失败: %v", err)
		return
	}
	log.Println("已启动生产环境低开销精确覆盖率采集")
}

*/

// -----------------------------------------------  Profiler.stop  -----------------------------------------------
// === 应用场景 ===
// 1. 停止性能采样：停止已启动的JS CPU性能采样，结束数据采集
// 2. 测试流程收尾：自动化性能测试完成后停止采样，准备获取结果
// 3. 性能数据固化：停止采集后确保性能数据完整，避免数据丢失
// 4. 调试流程结束：性能瓶颈定位完成后停止监控
// 5. 资源释放：停止采样后释放浏览器性能分析相关资源
// 6. 报告生成前置：停止采样后才能获取完整的性能报告数据

// CDPProfilerStop 停止JavaScript CPU性能采样分析
func CDPProfilerStop() (string, error) {
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
		"method": "Profiler.stop"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.stop 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.stop 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化性能测试 - 业务逻辑执行完成后停止采样
func ExampleProfilerStop_AutoPerformanceTest() {
	// 1. 启用分析器、设置采样间隔、启动采样
	_, _ = CDPProfilerEnable()
	_, _ = CDPProfilerSetSamplingInterval(500)
	_, _ = CDPProfilerStart()

	// 2. 执行测试业务逻辑
	// RunPerformanceTestLogic()

	// 3. 停止性能采样
	resp, err := CDPProfilerStop()
	if err != nil {
		log.Fatalf("停止性能采样失败: %v", err)
	}
	log.Printf("停止CPU采样成功，响应: %s", resp)

	// 4. 后续可获取性能报告、关闭分析器
}

// 场景2：页面卡顿调试 - 捕获到卡顿问题后停止采样分析数据
func ExampleProfilerStop_PageLagDebug() {
	// 已启动采样并触发卡顿操作
	// ...

	// 问题复现后立即停止采样，保存数据
	if _, err := CDPProfilerStop(); err != nil {
		log.Printf("停止采样失败: %v", err)
		return
	}
	log.Println("已停止CPU采样，可分析性能数据定位卡顿原因")
}

*/

// -----------------------------------------------  Profiler.stopPreciseCoverage  -----------------------------------------------
// === 应用场景 ===
// 1. 停止精准覆盖率采集：结束高精度JS代码覆盖率统计
// 2. 测试流程收尾：自动化测试完成精确覆盖采集后停止
// 3. 资源释放：停止精确覆盖监控，释放浏览器分析资源
// 4. 报告生成前置：停止采集后获取完整精准覆盖率报告
// 5. 多任务切换：切换不同调试/测试任务时关闭当前覆盖采集
// 6. 生产监控结束：线上轻量覆盖监控完成后安全停止

// CDPProfilerStopPreciseCoverage 停止精确JavaScript代码覆盖率采集
func CDPProfilerStopPreciseCoverage() (string, error) {
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
		"method": "Profiler.stopPreciseCoverage"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.stopPreciseCoverage 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.stopPreciseCoverage 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化测试完成后停止精确覆盖率采集并生成报告
func ExampleProfilerStopPreciseCoverage_AutoTestReport() {
	// 前置：启用分析器 + 启动精确覆盖率
	_, _ = CDPProfilerEnable()
	_, _ = CDPProfilerStartPreciseCoverage(true, true)

	// 执行完整测试用例
	// RunFullTestSuit()

	// 停止精确覆盖率采集
	resp, err := CDPProfilerStopPreciseCoverage()
	if err != nil {
		log.Fatalf("停止精确覆盖率失败: %v", err)
	}
	log.Printf("停止精确覆盖率成功: %s", resp)

	// 后续获取覆盖率数据 & 生成报告
	// coverage, _ := CDPProfilerGetBestEffortCoverage()
	// GenerateCoverageReport(coverage)
}

// 场景2：生产环境监控完成后停止采集，释放资源
func ExampleProfilerStopPreciseCoverage_ProductionClean() {
	// 线上低开销监控结束
	log.Println("准备停止生产环境精确覆盖率监控...")

	if _, err := CDPProfilerStopPreciseCoverage(); err != nil {
		log.Printf("停止失败: %v", err)
		return
	}

	// 关闭分析器，完全释放资源
	_, _ = CDPProfilerDisable()
	log.Println("已停止精确覆盖率采集，浏览器资源已释放")
}

*/

// -----------------------------------------------  Profiler.takePreciseCoverage  -----------------------------------------------
// === 应用场景 ===
// 1. 获取精准覆盖率数据：获取已启动的精确代码覆盖率采集结果
// 2. 自动化测试报告：测试过程中/结束后提取精确覆盖数据生成报告
// 3. 实时覆盖率监控：动态获取当前JS代码执行覆盖情况
// 4. 多阶段覆盖采集：分步骤采集不同业务场景的精确覆盖数据
// 5. 死代码检测分析：基于精准覆盖结果识别未执行代码
// 6. 质量门禁校验：CI/流程中提取覆盖率数据进行质量校验

// CDPProfilerTakePreciseCoverage 获取精确的JavaScript代码覆盖率数据
func CDPProfilerTakePreciseCoverage() (string, error) {
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
		"method": "Profiler.takePreciseCoverage"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Profiler.takePreciseCoverage 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应ID的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应判断是否出错
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("Profiler.takePreciseCoverage 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化测试中获取精确覆盖率生成测试报告
func ExampleProfilerTakePreciseCoverage_AutoTestReport() {
	// 1. 启用分析器
	_, err := CDPProfilerEnable()
	if err != nil {
		log.Fatalf("启用分析器失败: %v", err)
	}

	// 2. 启动精确覆盖率采集
	_, err = CDPProfilerStartPreciseCoverage(true, true)
	if err != nil {
		log.Fatalf("启动精确覆盖率失败: %v", err)
	}

	// 3. 执行测试用例
	// RunTestCases()

	// 4. 获取精准覆盖率数据
	coverage, err := CDPProfilerTakePreciseCoverage()
	if err != nil {
		log.Fatalf("获取精确覆盖率失败: %v", err)
	}
	log.Printf("精确覆盖率数据: %s", coverage)

	// 5. 停止采集并关闭分析器
	_, _ = CDPProfilerStopPreciseCoverage()
	_, _ = CDPProfilerDisable()
}

// 场景2：页面操作后实时获取覆盖率分析未执行代码
func ExampleProfilerTakePreciseCoverage_RealTimeAnalysis() {
	// 启动精确覆盖后执行用户交互操作
	// SimulateUserBehavior()

	// 实时获取当前覆盖情况用于分析死代码
	result, err := CDPProfilerTakePreciseCoverage()
	if err != nil {
		log.Printf("获取覆盖率失败: %v", err)
		return
	}

	// 保存数据用于前端优化分析
	// utils.SaveToFile("precise_coverage.json", result)
	log.Println("已获取实时精确覆盖率数据，可用于未执行代码检测")
}

*/
