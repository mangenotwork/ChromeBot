package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Memory.forciblyPurgeJavaScriptMemory  -----------------------------------------------
// === 应用场景 ===
// 1. 内存泄漏测试: 强制清理JS内存后重新检测内存占用
// 2. 性能压测: 长时间测试后主动释放内存，避免浏览器崩溃
// 3. 自动化测试: 测试用例执行完毕后清理内存，保证用例隔离性
// 4. 页面切换优化: 单页应用路由切换后强制回收无用内存
// 5. 调试辅助: 排查内存问题时手动触发垃圾回收
// 6. 资源紧张场景: 低配置设备上主动释放JS内存提升运行流畅度

// CDPMemoryForciblyPurgeJavaScriptMemory 强制清理JavaScript内存，触发浏览器垃圾回收
func CDPMemoryForciblyPurgeJavaScriptMemory() (string, error) {
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
		"method": "Memory.forciblyPurgeJavaScriptMemory"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.forciblyPurgeJavaScriptMemory 请求失败: %w", err)
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
			return "", fmt.Errorf("forciblyPurgeJavaScriptMemory 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：自动化测试后强制清理JS内存，防止内存累积影响后续用例
func ExampleMemoryForciblyPurgeJavaScriptMemory() {
	// 1. 执行页面业务逻辑/测试用例
	// DoPageTestLogic()

	// 2. 强制清理JavaScript内存，触发垃圾回收
	resp, err := CDPMemoryForciblyPurgeJavaScriptMemory()
	if err != nil {
		log.Printf("强制清理JS内存失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功强制清理JavaScript内存，浏览器已执行垃圾回收")
}
*/

// -----------------------------------------------  Memory.getAllTimeSamplingProfile  -----------------------------------------------
// === 应用场景 ===
// 1. 内存泄漏分析: 获取全时段内存采样数据，定位长期增长的内存占用
// 2. 性能诊断: 分析页面从加载到当前的完整内存使用趋势
// 3. 自动化测试: 测试结束后导出全量内存profile，用于离线分析
// 4. 长期运行监控: 收集长时间运行页面的内存累积数据
// 5. 前端优化: 基于全时段采样数据优化JS内存占用和垃圾回收
// 6. 问题复现: 保存内存profile文件，用于复现和排查内存异常问题

// CDPMemoryGetAllTimeSamplingProfile 获取浏览器从启动至今的全时段内存采样配置文件
func CDPMemoryGetAllTimeSamplingProfile() (string, error) {
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
		"method": "Memory.getAllTimeSamplingProfile"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.getAllTimeSamplingProfile 请求失败: %w", err)
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
			return "", fmt.Errorf("getAllTimeSamplingProfile 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：自动化测试完成后，获取全时段内存采样数据并保存用于泄漏分析
func ExampleMemoryGetAllTimeSamplingProfile() {
	// 1. 执行页面业务测试流程
	// RunPageTestCases()

	// 2. 获取全时段内存采样profile
	profileResp, err := CDPMemoryGetAllTimeSamplingProfile()
	if err != nil {
		log.Printf("获取全时段内存采样数据失败: %v, 响应: %s", err, profileResp)
		return
	}

	// 3. 保存内存数据到文件，用于离线分析
	err = os.WriteFile("memory_all_time_profile.json", []byte(profileResp), 0644)
	if err != nil {
		log.Printf("保存内存profile文件失败: %v", err)
		return
	}

	log.Println("成功获取并保存全时段内存采样配置文件：memory_all_time_profile.json")
}
*/

// -----------------------------------------------  Memory.getBrowserSamplingProfile  -----------------------------------------------
// === 应用场景 ===
// 1. 浏览器进程内存监控: 获取浏览器主进程的内存采样数据
// 2. 多标签内存分析: 监控整个浏览器而非单个页面的内存占用
// 3. 浏览器性能诊断: 排查浏览器整体内存泄漏问题
// 4. 自动化测试: 测试后采集浏览器全局内存数据
// 5. 资源监控: 长期监控浏览器进程内存消耗趋势
// 6. 调试辅助: 定位浏览器内核级别的内存异常

// CDPMemoryGetBrowserSamplingProfile 获取浏览器进程的内存采样配置文件
func CDPMemoryGetBrowserSamplingProfile() (string, error) {
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
		"method": "Memory.getBrowserSamplingProfile"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.getBrowserSamplingProfile 请求失败: %w", err)
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
			return "", fmt.Errorf("getBrowserSamplingProfile 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：监控浏览器整体进程内存，采集数据并保存用于诊断
func ExampleMemoryGetBrowserSamplingProfile() {
	// 1. 执行多标签/多页面测试操作
	// OpenMultiPagesTest()

	// 2. 获取浏览器进程级别的内存采样数据
	resp, err := CDPMemoryGetBrowserSamplingProfile()
	if err != nil {
		log.Printf("获取浏览器内存采样失败: %v, 响应: %s", err, resp)
		return
	}

	// 3. 保存浏览器内存数据到文件
	err = os.WriteFile("browser_memory_profile.json", []byte(resp), 0644)
	if err != nil {
		log.Printf("保存浏览器内存profile失败: %v", err)
		return
	}

	log.Println("成功获取并保存浏览器进程内存采样文件")
}
*/

// -----------------------------------------------  Memory.getDOMCounters  -----------------------------------------------
// === 应用场景 ===
// 1. DOM节点监控: 实时获取页面DOM节点数量、事件监听器数量
// 2. 内存泄漏排查: 检测DOM节点未释放、监听器泄漏问题
// 3. 页面性能优化: 监控DOM数量过高导致的渲染性能下降
// 4. 自动化测试: 测试用例执行后检查DOM计数器是否异常增长
// 5. 单页应用监控: 路由切换后验证DOM是否正确清理回收
// 6. 调试辅助: 快速定位页面DOM相关的内存占用问题

// CDPMemoryGetDOMCounters 获取页面DOM节点计数器信息
func CDPMemoryGetDOMCounters() (string, error) {
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
		"method": "Memory.getDOMCounters"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.getDOMCounters 请求失败: %w", err)
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
			return "", fmt.Errorf("getDOMCounters 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：页面路由切换后，获取DOM计数器检测是否存在泄漏
func ExampleMemoryGetDOMCounters() {
	// 1. 模拟页面路由切换操作
	// SwitchPageRoute()

	// 2. 获取DOM节点、事件监听器等计数器数据
	resp, err := CDPMemoryGetDOMCounters()
	if err != nil {
		log.Printf("获取DOM计数器失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功获取DOM计数器信息：", resp)

	// 响应包含字段：nodes(DOM节点数), eventListeners(事件监听器数), jsEventListeners(JS监听器数)
}
*/

// -----------------------------------------------  Memory.getDOMCountersForLeakDetection  -----------------------------------------------
// === 应用场景 ===
// 1. 内存泄漏专项检测：专门用于检测DOM节点、事件监听器泄漏
// 2. 自动化泄漏测试：测试流程中精准采集泄漏检测专用DOM计数器
// 3. 单页应用路由泄漏：路由切换后对比计数器判断DOM未释放
// 4. 组件卸载检测：前端组件销毁后检查是否残留DOM/监听器
// 5. 长期运行监控：持续采集专用指标预警内存泄漏风险
// 6. 调试定位：快速定位导致泄漏的DOM节点与监听器来源

// CDPMemoryGetDOMCountersForLeakDetection 获取用于泄漏检测的DOM计数器信息
func CDPMemoryGetDOMCountersForLeakDetection() (string, error) {
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
		"method": "Memory.getDOMCountersForLeakDetection"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.getDOMCountersForLeakDetection 请求失败: %w", err)
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
			return "", fmt.Errorf("getDOMCountersForLeakDetection 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：自动化检测DOM泄漏，对比组件加载前后计数器判断是否泄漏
func ExampleMemoryGetDOMCountersForLeakDetection() {
	// 1. 组件加载前获取泄漏检测计数器
	before, err := CDPMemoryGetDOMCountersForLeakDetection()
	if err != nil {
		log.Printf("获取前置DOM泄漏计数器失败：%v", err)
		return
	}

	// 2. 执行组件加载、交互、卸载操作
	// LoadComponent()
	// ComponentInteraction()
	// UnloadComponent()

	// 3. 组件卸载后再次获取专用计数器
	after, err := CDPMemoryGetDOMCountersForLeakDetection()
	if err != nil {
		log.Printf("获取后置DOM泄漏计数器失败：%v", err)
		return
	}

	log.Println("==== DOM泄漏检测报告 ====")
	log.Println("操作前计数器：", before)
	log.Println("操作后计数器：", after)
	log.Println("若节点数/监听器数持续增长，说明存在DOM泄漏！")
}
*/

// -----------------------------------------------  Memory.getSamplingProfile  -----------------------------------------------
// === 应用场景 ===
// 1. 实时内存分析: 获取当前渲染进程的内存采样数据
// 2. 页面内存诊断: 定位单页面JS内存占用过高的问题
// 3. 自动化测试: 采集页面内存profile用于性能回归检测
// 4. 组件性能监控: 检测单个前端组件的内存消耗情况
// 5. 内存泄漏排查: 对比多次采样数据发现内存持续增长
// 6. 性能优化: 基于内存采样数据优化JS对象与内存使用

// CDPMemoryGetSamplingProfile 获取当前渲染进程的内存采样配置文件
func CDPMemoryGetSamplingProfile() (string, error) {
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
		"method": "Memory.getSamplingProfile"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.getSamplingProfile 请求失败: %w", err)
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
			return "", fmt.Errorf("getSamplingProfile 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：页面交互后采集内存采样数据，保存到文件用于性能分析
func ExampleMemoryGetSamplingProfile() {
	// 1. 执行页面交互逻辑（如点击、加载数据等）
	// PerformPageInteraction()

	// 2. 获取当前页面渲染进程内存采样数据
	resp, err := CDPMemoryGetSamplingProfile()
	if err != nil {
		log.Printf("获取页面内存采样失败: %v, 响应: %s", err, resp)
		return
	}

	// 3. 保存内存采样文件，可在Chrome DevTools中导入分析
	err = os.WriteFile("page_sampling_profile.json", []byte(resp), 0644)
	if err != nil {
		log.Printf("保存内存采样文件失败: %v", err)
		return
	}

	log.Println("成功获取页面内存采样文件，可在Chrome DevTools Memory面板导入分析")
}
*/

// -----------------------------------------------  Memory.prepareForLeakDetection  -----------------------------------------------
// === 应用场景 ===
// 1. 泄漏检测前置准备: 内存泄漏测试前执行环境初始化
// 2. 自动化测试隔离: 确保测试用例间内存环境干净独立
// 3. DOM泄漏检测: 检测前清理全局DOM引用，提升检测准确性
// 4. 单页应用测试: 路由切换前重置内存检测环境
// 5. 组件测试: 组件加载前准备检测环境，避免干扰
// 6. 精准泄漏定位: 排除历史数据干扰，获取真实泄漏数据

// CDPMemoryPrepareForLeakDetection 为内存泄漏检测准备环境（强制GC+清理引用）
func CDPMemoryPrepareForLeakDetection() (string, error) {
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
		"method": "Memory.prepareForLeakDetection"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.prepareForLeakDetection 请求失败: %w", err)
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
			return "", fmt.Errorf("prepareForLeakDetection 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景：自动化内存泄漏检测流程 - 准备环境 → 执行操作 → 检测泄漏
func ExampleMemoryPrepareForLeakDetection() {
	// 1. 【关键】泄漏检测前准备环境：强制垃圾回收、清理引用、排除干扰
	resp, err := CDPMemoryPrepareForLeakDetection()
	if err != nil {
		log.Printf("泄漏检测环境准备失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功准备泄漏检测环境，可开始采集DOM/内存数据")

	// 2. 执行测试操作（加载/交互/卸载组件）
	// LoadTestComponent()
	// UnloadTestComponent()

	// 3. 获取专用计数器进行泄漏检测
	// counters, _ := CDPMemoryGetDOMCountersForLeakDetection()
	// log.Println("泄漏检测DOM计数器:", counters)
}
*/

// -----------------------------------------------  Memory.setPressureNotificationsSuppressed  -----------------------------------------------
// === 应用场景 ===
// 1. 压力测试屏蔽: 内存压力测试时屏蔽系统通知，避免干扰
// 2. 自动化测试稳定: 抑制内存压力提醒，保证测试流程不中断
// 3. 性能监控静默: 后台采集内存数据时不触发浏览器通知
// 4. 调试环境优化: 调试内存问题时关闭压力提示弹窗
// 5. 长时间运行防护: 避免频繁内存压力通知影响程序稳定性
// 6. 沉浸式测试: 全量压测时屏蔽所有非必要内存提醒

// CDPMemorySetPressureNotificationsSuppressed 设置是否屏蔽内存压力通知
// 参数：suppressed - true 屏蔽通知，false 开启通知
func CDPMemorySetPressureNotificationsSuppressed(suppressed bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Memory.setPressureNotificationsSuppressed",
		"params": {
			"suppressed": %t
		}
	}`, reqID, suppressed)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.setPressureNotificationsSuppressed 请求失败: %w", err)
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
			return "", fmt.Errorf("setPressureNotificationsSuppressed 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：内存压力测试开始前屏蔽通知，测试结束后恢复
func ExampleMemorySetPressureNotificationsSuppressed() {
	// 1. 屏蔽内存压力通知，避免测试过程中弹出提醒
	resp, err := CDPMemorySetPressureNotificationsSuppressed(true)
	if err != nil {
		log.Printf("屏蔽内存压力通知失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("已成功屏蔽内存压力通知，开始执行内存压力测试...")

	// 2. 执行内存压力测试逻辑（如大量创建DOM、大对象等）
	// RunMemoryPressureTest()

	// 3. 测试完成，恢复内存压力通知（可选）
	// restoreResp, restoreErr := CDPMemorySetPressureNotificationsSuppressed(false)
	// if restoreErr != nil {
	// 	log.Printf("恢复内存通知失败: %v", restoreErr)
	// 	return
	// }
	// log.Println("内存压力测试完成，已恢复通知")
}
*/

// -----------------------------------------------  Memory.simulatePressureNotification  -----------------------------------------------
// === 应用场景 ===
// 1. 内存压力测试：模拟不同等级内存压力，测试页面/浏览器应对策略
// 2. 稳定性测试：验证内存紧张时程序是否崩溃、卡顿
// 3. 前端优化验证：测试内存告警时代码是否正常释放资源
// 4. 自动化压测：模拟内存压力触发浏览器回收机制
// 5. 低内存适配：模拟移动端/低配设备内存不足场景
// 6. 告警逻辑测试：测试内存压力通知相关业务逻辑

// CDPMemorySimulatePressureNotification 模拟内存压力通知
// 参数：level - 内存压力等级，可选值：moderate(中等)、critical(严重)、none(无压力)
func CDPMemorySimulatePressureNotification(level string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	// 校验参数合法性
	validLevels := map[string]bool{"moderate": true, "critical": true, "none": true}
	if !validLevels[level] {
		return "", fmt.Errorf("无效的内存压力等级: %s，支持: moderate, critical, none", level)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Memory.simulatePressureNotification",
		"params": {
			"level": "%s"
		}
	}`, reqID, level)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.simulatePressureNotification 请求失败: %w", err)
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
			return "", fmt.Errorf("simulatePressureNotification 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：模拟严重内存压力，测试页面是否能正常处理内存告警
func ExampleMemorySimulatePressureNotification() {
	// 1. 模拟【严重】内存压力，触发浏览器内存告警
	resp, err := CDPMemorySimulatePressureNotification("critical")
	if err != nil {
		log.Printf("模拟严重内存压力失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功模拟【严重级别】内存压力通知，页面应触发内存优化/回收逻辑")

	// 可选：模拟中等压力
	// resp, err = CDPMemorySimulatePressureNotification("moderate")

	// 可选：取消模拟，恢复正常内存状态
	// resp, err = CDPMemorySimulatePressureNotification("none")
}
*/

// -----------------------------------------------  Memory.startSampling  -----------------------------------------------
// === 应用场景 ===
// 1. 内存采样监控：启动对JS堆内存的持续采样
// 2. 性能测试前置：测试开始前启动采样，记录内存变化轨迹
// 3. 内存泄漏追踪：长期采样定位缓慢增长的内存泄漏点
// 4. 页面加载分析：记录页面从加载到运行的完整内存曲线
// 5. 自动化性能采集：自动化用例中开启采样做性能基线
// 6. 交互内存分析：用户交互过程中实时采集内存分配数据

// CDPMemoryStartSampling 启动内存采样
// 参数：samplingInterval - 采样间隔（单位：字节，默认126976）
func CDPMemoryStartSampling(samplingInterval int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Memory.startSampling",
		"params": {
			"samplingInterval": %d
		}
	}`, reqID, samplingInterval)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.startSampling 请求失败: %w", err)
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
			return "", fmt.Errorf("startSampling 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：启动内存采样，执行页面操作后获取采样数据做性能分析
func ExampleMemoryStartSampling() {
	// 1. 启动内存采样，使用默认间隔 126976 字节
	interval := 126976
	resp, err := CDPMemoryStartSampling(interval)
	if err != nil {
		log.Printf("启动内存采样失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功启动内存采样，开始记录页面内存分配数据")

	// 2. 执行页面业务操作（加载、交互、路由切换等）
	// RunPageBusinessLogic()

	// 3. 后续可调用 CDPMemoryStopSampling() 停止采样并获取数据
}
*/

// -----------------------------------------------  Memory.stopSampling  -----------------------------------------------
// === 应用场景 ===
// 1. 内存采样结束：停止内存采样并获取最终采样数据
// 2. 性能测试收尾：测试流程完成后终止采样，避免资源占用
// 3. 数据采集完成：获取完整内存采样profile用于离线分析
// 4. 自动化测试闭环：与startSampling配对完成性能采集
// 5. 内存诊断结束：停止监控释放浏览器采样资源
// 6. 数据保存：采样停止后导出数据文件用于调试分析

// CDPMemoryStopSampling 停止内存采样并返回采样数据
func CDPMemoryStopSampling() (string, error) {
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
		"method": "Memory.stopSampling"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Memory.stopSampling 请求失败: %w", err)
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
			return "", fmt.Errorf("stopSampling 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：完整内存采样流程 - 启动 → 操作 → 停止 → 保存数据
func ExampleMemoryStopSampling() {
	// 1. 先启动内存采样
	// err := CDPMemoryStartSampling(126976)
	// if err != nil {
	// 	log.Fatalf("启动采样失败: %v", err)
	// }

	// 2. 执行页面交互/测试逻辑
	// RunPageTest()

	// 3. 停止采样并获取采样数据
	resp, err := CDPMemoryStopSampling()
	if err != nil {
		log.Printf("停止内存采样失败: %v, 响应: %s", err, resp)
		return
	}

	// 4. 保存采样数据到文件
	err = os.WriteFile("memory_sampling_result.json", []byte(resp), 0644)
	if err != nil {
		log.Printf("保存采样数据失败: %v", err)
		return
	}

	log.Println("成功停止内存采样，数据已保存：memory_sampling_result.json")
}
*/
