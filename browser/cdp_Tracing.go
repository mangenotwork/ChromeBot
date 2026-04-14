package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Tracing.end  -----------------------------------------------
// === 应用场景 ===
// 1. 性能追踪结束: 完成性能数据采集后停止追踪
// 2. 自动化测试收尾: 自动化测试流程中结束性能追踪
// 3. 调试完成: 问题调试结束后关闭追踪功能
// 4. 资源释放: 停止追踪以释放浏览器性能消耗资源
// 5. 数据采集完成: 性能数据采集完毕后终止追踪会话
// 6. 定时追踪结束: 达到预设时间后自动结束追踪

// CDPTracingEnd 停止浏览器性能追踪
func CDPTracingEnd() (string, error) {
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
		"method": "Tracing.end"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tracing.end 请求失败: %w", err)
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
			return "", fmt.Errorf("Tracing.end 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1: 基础结束追踪并处理结果
func ExampleCDPTracingEnd_Basic() {
	// 结束追踪
	result, err := CDPTracingEnd()
	if err != nil {
		log.Fatalf("结束追踪失败: %v", err)
	}
	log.Printf("追踪已停止，响应结果: %s", result)
}

// 示例2: 自动化测试中结束追踪+清理
func ExampleCDPTracingEnd_AutoTest() {
	// 假设已执行 Tracing.start 开启追踪
	// ... 执行测试操作 ...

	// 结束追踪
	resp, err := CDPTracingEnd()
	if err != nil {
		log.Printf("结束追踪异常: %v", err)
		return
	}
	log.Printf("测试完成，追踪已终止: %s", resp)
}

// 示例3: 定时自动结束追踪
func ExampleCDPTracingEnd_Timed() {
	// 开启追踪后等待3秒自动结束
	time.Sleep(3 * time.Second)

	result, err := CDPTracingEnd()
	if err != nil {
		log.Printf("定时结束追踪失败: %v", err)
		return
	}
	log.Printf("定时任务：追踪已停止: %s", result)
}

*/

// -----------------------------------------------  Tracing.start  -----------------------------------------------
// === 应用场景 ===
// 1. 页面性能采集: 启动浏览器性能追踪，采集页面加载、渲染、JS执行等性能数据
// 2. 自动化性能测试: 自动化测试流程中开启性能追踪，监控页面性能瓶颈
// 3. 卡顿问题调试: 调试页面卡顿、延迟、无响应等问题时采集追踪数据
// 4. 内存泄漏分析: 开启追踪采集内存分配、GC相关数据，分析内存泄漏问题
// 5. 网络性能监控: 结合网络追踪，采集网络请求、响应的完整性能数据
// 6. 前端性能优化: 优化前端代码前采集基准数据，优化后对比性能提升效果

// CDPTracingStart 启动浏览器性能追踪
// 参数 traceConfig: 追踪配置JSON字符串，参考Chrome DevTools Protocol Tracing域配置
func CDPTracingStart(traceConfig string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息，支持自定义追踪配置
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Tracing.start",
		"params": %s
	}`, reqID, traceConfig)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tracing.start 请求失败: %w", err)
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
			return "", fmt.Errorf("Tracing.start 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1: 基础性能追踪（采集所有性能数据，标准配置）
func ExampleCDPTracingStart_Basic() {
	// 标准追踪配置：采集所有性能数据
	traceConfig := `{
		"categories": "*,devtools.timeline",
		"transferMode": "ReturnAsStream"
	}`

	// 启动追踪
	result, err := CDPTracingStart(traceConfig)
	if err != nil {
		log.Fatalf("启动性能追踪失败: %v", err)
	}
	log.Printf("追踪已启动，响应结果: %s", result)
}

// 示例2: 专注页面渲染/JS执行追踪（轻量化配置）
func ExampleCDPTracingStart_RenderJS() {
	// 仅追踪渲染、JS执行、事件循环相关数据
	traceConfig := `{
		"categories": "blink,blink_gc,v8,devtools.timeline",
		"transferMode": "ReturnAsStream"
	}`

	resp, err := CDPTracingStart(traceConfig)
	if err != nil {
		log.Printf("启动渲染追踪失败: %v", err)
		return
	}
	log.Printf("渲染&JS追踪已启动: %s", resp)
}

// 示例3: 网络+性能综合追踪
func ExampleCDPTracingStart_Network() {
	// 包含网络请求、网络栈、页面性能的综合追踪配置
	traceConfig := `{
		"categories": "netlog,blink,devtools.timeline,devtools.network",
		"transferMode": "ReturnAsStream"
	}`

	result, err := CDPTracingStart(traceConfig)
	if err != nil {
		log.Printf("启动网络性能追踪失败: %v", err)
		return
	}
	log.Printf("网络+性能综合追踪已启动: %s", result)
}

*/

// -----------------------------------------------  Tracing.getCategories  -----------------------------------------------
// === 应用场景 ===
// 1. 追踪配置前置检查: 启动追踪前获取可用分类，验证配置合法性
// 2. 动态生成追踪配置: 根据支持的分类自动构建性能追踪参数
// 3. 调试追踪配置: 排查因分类名称错误导致的追踪启动失败
// 4. 跨浏览器兼容: 不同Chrome版本兼容，获取当前环境支持的追踪分类
// 5. 性能工具集成: 性能监控工具中展示可选的追踪分类列表
// 6. 自动化配置校验: 自动化测试中校验追踪配置是否有效

// CDPTracingGetCategories 获取浏览器支持的性能追踪分类列表
func CDPTracingGetCategories() (string, error) {
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
		"method": "Tracing.getCategories"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tracing.getCategories 请求失败: %w", err)
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
			return "", fmt.Errorf("Tracing.getCategories 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1: 基础获取所有追踪分类并打印
func ExampleCDPTracingGetCategories_Basic() {
	// 获取支持的追踪分类
	categories, err := CDPTracingGetCategories()
	if err != nil {
		log.Fatalf("获取追踪分类失败: %v", err)
	}
	log.Printf("浏览器支持的追踪分类: %s", categories)
}

// 示例2: 启动追踪前校验+获取分类（推荐生产使用）
func ExampleCDPTracingGetCategories_PreCheck() {
	// 先获取支持的分类，确保配置有效
	_, err := CDPTracingGetCategories()
	if err != nil {
		log.Printf("获取分类失败，无法启动追踪: %v", err)
		return
	}

	// 分类获取成功，继续启动追踪
	log.Println("追踪分类校验通过，准备启动Tracing.start...")
}

// 示例3: 自动化测试中动态获取分类配置
func ExampleCDPTracingGetCategories_AutoTest() {
	// 自动化测试流程：获取分类 -> 启动追踪 -> 执行测试
	catResp, err := CDPTracingGetCategories()
	if err != nil {
		log.Fatalf("自动化测试：获取追踪分类失败: %v", err)
	}
	log.Printf("自动化测试：可用分类获取成功: %s", catResp)

	// 后续可基于返回的categories动态构建traceConfig启动追踪
}

*/

// -----------------------------------------------  Tracing.getTrackEventDescriptors  -----------------------------------------------
// === 应用场景 ===
// 1. 高级追踪配置：获取事件描述符，用于精准定义需要追踪的事件类型
// 2. 自定义性能分析：基于描述符筛选特定事件，实现精细化性能数据采集
// 3. 调试追踪规则：验证和调试追踪事件过滤规则是否生效
// 4. 性能工具开发：构建专业性能分析工具时获取标准事件描述符
// 5. 精准问题定位：只采集特定类型事件，减少无效数据，快速定位问题
// 6. 跨版本兼容：确保不同Chrome版本追踪事件类型的兼容性

// CDPTracingGetTrackEventDescriptors 获取浏览器追踪事件描述符
func CDPTracingGetTrackEventDescriptors() (string, error) {
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
		"method": "Tracing.getTrackEventDescriptors"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tracing.getTrackEventDescriptors 请求失败: %w", err)
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
			return "", fmt.Errorf("Tracing.getTrackEventDescriptors 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1: 基础获取事件描述符
func ExampleCDPTracingGetTrackEventDescriptors_Basic() {
	// 获取追踪事件描述符列表
	descriptors, err := CDPTracingGetTrackEventDescriptors()
	if err != nil {
		log.Fatalf("获取追踪事件描述符失败: %v", err)
	}
	log.Printf("追踪事件描述符: %s", descriptors)
}

// 示例2: 高级性能分析前获取描述符
func ExampleCDPTracingGetTrackEventDescriptors_Analysis() {
	// 先获取事件描述符，用于自定义精准追踪配置
	resp, err := CDPTracingGetTrackEventDescriptors()
	if err != nil {
		log.Printf("获取事件描述符异常: %v", err)
		return
	}
	log.Printf("已获取事件描述符，可用于自定义追踪规则: %s", resp)
}

// 示例3: 自动化性能测试流程
func ExampleCDPTracingGetTrackEventDescriptors_AutoTest() {
	// 步骤1：获取事件描述符
	_, err := CDPTracingGetTrackEventDescriptors()
	if err != nil {
		log.Fatalf("自动化测试：获取事件描述符失败: %v", err)
	}

	// 步骤2：基于描述符启动精准追踪
	log.Println("自动化测试：事件描述符获取完成，开始启动自定义追踪...")

	// 步骤3：执行测试...
}

*/

// -----------------------------------------------  Tracing.recordClockSyncMarker  -----------------------------------------------
// === 应用场景 ===
// 1. 跨进程时钟同步：在浏览器与服务端/客户端之间记录时间同步标记，统一时间轴
// 2. 分布式性能分析：多进程、多设备协作时对齐追踪数据的时间戳
// 3. 精准时序分析：确保性能追踪事件与外部日志/数据的时间完全一致
// 4. 自动化测试校准：测试流程中标记同步点，精准分析操作耗时
// 5. 混合应用调试：Native + Web 混合调试时统一时钟，定位时序问题
// 6. 性能报告生成：生成标准化性能报告时校准时间基准

// CDPTracingRecordClockSyncMarker 记录时钟同步标记
// 参数 syncId：自定义同步标识符，用于关联外部系统的同步标记
func CDPTracingRecordClockSyncMarker(syncId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息，传入同步ID参数
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Tracing.recordClockSyncMarker",
		"params": {
			"syncId": "%s"
		}
	}`, reqID, syncId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tracing.recordClockSyncMarker 请求失败: %w", err)
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
			return "", fmt.Errorf("Tracing.recordClockSyncMarker 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1: 基础时钟同步标记（简单同步ID）
func ExampleCDPTracingRecordClockSyncMarker_Basic() {
	// 记录基础同步标记，使用固定ID
	result, err := CDPTracingRecordClockSyncMarker("sync_1001")
	if err != nil {
		log.Fatalf("记录时钟同步标记失败: %v", err)
	}
	log.Printf("同步标记已记录，响应: %s", result)
}

// 示例2: 动态生成同步ID（推荐用于自动化场景）
func ExampleCDPTracingRecordClockSyncMarker_Dynamic() {
	// 生成唯一同步ID
	uuid := fmt.Sprintf("sync_%d", time.Now().UnixMilli())

	// 记录同步标记
	resp, err := CDPTracingRecordClockSyncMarker(uuid)
	if err != nil {
		log.Printf("记录动态同步标记失败: %v", err)
		return
	}
	log.Printf("动态同步标记记录成功: %s, ID: %s", resp, uuid)
}

// 示例3: 性能追踪流程中插入时钟同步（完整流程）
func ExampleCDPTracingRecordClockSyncMarker_FullTrace() {
	// 1. 启动性能追踪
	traceConfig := `{"categories": "*,devtools.timeline","transferMode": "ReturnAsStream"}`
	CDPTracingStart(traceConfig)

	// 2. 记录时钟同步点（关键：对齐外部时间）
	syncID := "test_sync_001"
	_, err := CDPTracingRecordClockSyncMarker(syncID)
	if err != nil {
		log.Fatalf("同步失败: %v", err)
	}

	// 3. 执行业务操作...
	time.Sleep(1 * time.Second)

	// 4. 结束追踪
	CDPTracingEnd()
	log.Println("完整追踪流程完成，已插入时钟同步标记")
}

*/

// -----------------------------------------------  Tracing.requestMemoryDump  -----------------------------------------------
// === 应用场景 ===
// 1. 内存泄漏检测：主动触发内存转储，分析页面内存占用情况
// 2. 自动化性能测试：测试流程中定点采集内存快照，监控内存变化趋势
// 3. 页面卡顿调试：页面卡顿、崩溃时采集内存数据，定位异常占用
// 4. 资源优化分析：采集JS堆、DOM节点、渲染节点等内存数据，优化资源占用
// 5. 长时间运行监控：长时间运行的Web应用定时采集内存，防止内存溢出
// 6. 混合应用内存分析：Native+Web混合应用，采集WebView内存数据

// CDPTracingRequestMemoryDump 主动请求浏览器执行内存转储
// 参数 dumpMode：内存转储模式，支持：light/heap/verbose
func CDPTracingRequestMemoryDump(dumpMode string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息，传入内存转储配置
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Tracing.requestMemoryDump",
		"params": {
			"dumpMode": "%s"
		}
	}`, reqID, dumpMode)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Tracing.requestMemoryDump 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
	timeout := 8 * time.Second // 内存转储耗时稍长，超时时间延长
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
			return "", fmt.Errorf("Tracing.requestMemoryDump 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1: 轻量级内存转储（快速采集，性能消耗低）
func ExampleCDPTracingRequestMemoryDump_Light() {
	// light模式：快速采集基础内存数据
	result, err := CDPTracingRequestMemoryDump("light")
	if err != nil {
		log.Fatalf("轻量级内存转储失败: %v", err)
	}
	log.Printf("轻量级内存转储完成: %s", result)
}

// 示例2: 完整堆内存转储（深度分析内存泄漏）
func ExampleCDPTracingRequestMemoryDump_Heap() {
	// heap模式：采集完整JS堆内存，用于泄漏分析
	resp, err := CDPTracingRequestMemoryDump("heap")
	if err != nil {
		log.Printf("堆内存转储失败: %v", err)
		return
	}
	log.Printf("堆内存转储完成，可用于泄漏分析: %s", resp)
}

// 示例3: 性能追踪+内存转储组合流程（完整内存分析）
func ExampleCDPTracingRequestMemoryDump_Full() {
	// 1. 启动性能追踪
	traceConfig := `{"categories": "*,devtools.memory","transferMode": "ReturnAsStream"}`
	CDPTracingStart(traceConfig)

	// 2. 执行业务操作
	time.Sleep(2 * time.Second)

	// 3. 主动触发详细内存转储
	_, err := CDPTracingRequestMemoryDump("verbose")
	if err != nil {
		log.Fatalf("详细内存转储失败: %v", err)
	}
	log.Println("已触发详细内存转储")

	// 4. 结束追踪
	CDPTracingEnd()
	log.Println("内存分析流程完成")
}

*/
