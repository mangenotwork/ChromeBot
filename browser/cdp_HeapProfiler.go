package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  HeapProfiler.addInspectedHeapObject  -----------------------------------------------
// === 应用场景 ===
// 1. 内存调试：将指定堆对象加入监控列表，用于针对性分析内存占用
// 2. 泄漏检测：跟踪关键对象，检测是否发生意外内存泄漏
// 3. 性能分析：对核心业务对象进行堆内存专项监控
// 4. 自动化诊断：在自动化测试中监控指定对象的内存生命周期
// 5. 问题定位：快速定位特定JS对象的内存引用与占用情况
// 6. 长期监控：持续监控高频创建/销毁对象，优化内存使用

// CDPHeapProfilerAddInspectedHeapObject 将指定堆对象添加到受检查列表
// 参数：heapObjectId - 堆对象ID
func CDPHeapProfilerAddInspectedHeapObject(heapObjectId string) (string, error) {
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
		"method": "HeapProfiler.addInspectedHeapObject",
		"params": {
			"heapObjectId": "%s"
		}
	}`, reqID, heapObjectId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 addInspectedHeapObject 请求失败: %w", err)
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
			return "", fmt.Errorf("addInspectedHeapObject 请求超时")
		}
	}
}

/*
// === 使用示例1：基础监控堆对象 ===
func ExampleAddInspectedHeapObject_Base() {
	// 传入已知的堆对象ID
	objectId := "heap://12345678"
	resp, err := CDPHeapProfilerAddInspectedHeapObject(objectId)
	if err != nil {
		log.Fatalf("添加监控堆对象失败: %v", err)
	}
	log.Printf("添加成功: %s", resp)
}

// === 使用示例2：内存泄漏检测流程 ===
func ExampleAddInspectedHeapObject_LeakCheck() {
	// 业务关键对象ID
	targetObj := "heap://business-core-object"
	_, err := CDPHeapProfilerAddInspectedHeapObject(targetObj)
	if err != nil {
		log.Printf("添加监控失败: %v", err)
		return
	}
	log.Println("已加入堆监控，可开始检测内存泄漏")
}

// === 使用示例3：自动化内存诊断 ===
func ExampleAddInspectedHeapObject_AutoTest() {
	// 自动化测试中监控高频对象
	testObjectId := "heap://test-auto-obj"
	resp, err := CDPHeapProfilerAddInspectedHeapObject(testObjectId)
	if err != nil {
		log.Fatalf("自动化监控失败: %v", err)
	}
	log.Println("自动化堆对象监控已启动", resp)
}
*/

// -----------------------------------------------  HeapProfiler.collectGarbage  -----------------------------------------------
// === 应用场景 ===
// 1. 内存测试前清理：在内存分析、堆快照前手动触发GC，排除垃圾数据干扰
// 2. 自动化测试：测试流程中强制回收内存，保证每次测试环境内存干净
// 3. 性能基准测试：获取准确内存占用基线，避免垃圾内存影响测试结果
// 4. 服务端渲染优化：长时间运行无头浏览器时定期回收内存
// 5. 内存泄漏排查：GC后仍不释放的对象可判定为泄漏对象
// 6. 页面资源释放：页面切换/关闭后强制GC，释放DOM、JS对象内存

// CDPHeapProfilerCollectGarbage 手动触发浏览器垃圾回收
func CDPHeapProfilerCollectGarbage() (string, error) {
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
		"method": "HeapProfiler.collectGarbage"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 collectGarbage 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
	timeout := 10 * time.Second // GC可能稍慢，延长超时
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
			return "", fmt.Errorf("collectGarbage 请求超时")
		}
	}
}

/*
// === 使用示例1：内存快照前触发GC ===
func ExampleCollectGarbage_BeforeSnapshot() {
	// 先GC清理垃圾
	_, err := CDPHeapProfilerCollectGarbage()
	if err != nil {
		log.Fatalf("GC触发失败: %v", err)
	}
	log.Println("GC完成，可开始采集堆快照")
	// 后续调用 HeapProfiler.takeHeapSnapshot
}

// === 使用示例2：自动化测试环境清理 ===
func ExampleCollectGarbage_AutoTest() {
	// 每个测试用例执行前GC，保证环境干净
	resp, err := CDPHeapProfilerCollectGarbage()
	if err != nil {
		log.Printf("测试前GC失败: %v", err)
		return
	}
	log.Println("测试环境GC完成：", resp)
}

// === 使用示例3：无头浏览器定期内存优化 ===
func ExampleCollectGarbage_Periodic() {
	// 定时任务：每5分钟触发一次GC
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		_, err := CDPHeapProfilerCollectGarbage()
		if err != nil {
			log.Printf("定期GC失败: %v", err)
			continue
		}
		log.Println("定期GC执行完成，内存已优化")
	}
}
*/

// -----------------------------------------------  HeapProfiler.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 停止堆分析：结束堆内存监控、停止所有HeapProfiler相关事件
// 2. 资源释放：关闭堆分析器，释放浏览器占用的CPU与内存资源
// 3. 测试流程结束：自动化测试完成后关闭堆分析，恢复浏览器默认状态
// 4. 性能模式切换：从内存调试模式切换回正常运行模式
// 5. 内存泄漏检测结束：完成泄漏排查后关闭堆分析器
// 6. 服务端资源回收：无头浏览器任务结束后关闭不必要的分析功能

// CDPHeapProfilerDisable 禁用堆分析器，停止所有堆监控事件
func CDPHeapProfilerDisable() (string, error) {
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
		"method": "HeapProfiler.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 disable 请求失败: %w", err)
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
			return "", fmt.Errorf("disable 请求超时")
		}
	}
}

/*
// === 使用示例1：堆分析完成后关闭 ===
func ExampleHeapProfilerDisable_Finish() {
	// 堆快照、内存分析完成后关闭
	resp, err := CDPHeapProfilerDisable()
	if err != nil {
		log.Fatalf("关闭堆分析器失败: %v", err)
	}
	log.Println("堆分析器已关闭，资源已释放:", resp)
}

// === 使用示例2：自动化测试结束清理 ===
func ExampleHeapProfilerDisable_TestClean() {
	// 测试用例执行完毕，关闭堆监控
	_, err := CDPHeapProfilerDisable()
	if err != nil {
		log.Printf("清理堆分析器失败: %v", err)
		return
	}
	log.Println("测试完成：HeapProfiler 已禁用")
}

// === 使用示例3：切换浏览器运行模式 ===
func ExampleHeapProfilerDisable_SwitchMode() {
	// 从调试模式切换为高性能模式
	resp, err := CDPHeapProfilerDisable()
	if err != nil {
		log.Fatalf("切换模式失败: %v", err)
	}
	log.Println("已关闭内存分析，浏览器恢复正常模式:", resp)
}
*/

// -----------------------------------------------  HeapProfiler.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 内存监控启动：开启堆分析器，开始监听堆内存相关事件
// 2. 内存泄漏检测：启动堆分析，准备追踪对象引用与内存分配
// 3. 性能调试：开启堆分析功能，用于页面性能与内存问题排查
// 4. 自动化测试：测试前启动堆分析器，收集内存数据
// 5. 堆快照采集：采集堆快照前必须先启用HeapProfiler
// 6. 实时内存跟踪：开启后可接收堆对象分配、销毁等实时事件

// CDPHeapProfilerEnable 启用堆分析器
func CDPHeapProfilerEnable() (string, error) {
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
		"method": "HeapProfiler.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 enable 请求失败: %w", err)
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
			return "", fmt.Errorf("enable 请求超时")
		}
	}
}

/*
// === 使用示例1：启动内存分析流程 ===
func ExampleHeapProfilerEnable_Start() {
	// 开始内存分析前必须先enable
	resp, err := CDPHeapProfilerEnable()
	if err != nil {
		log.Fatalf("启用堆分析器失败: %v", err)
	}
	log.Println("HeapProfiler 已启用，可以开始内存分析:", resp)
}

// === 使用示例2：内存泄漏检测启动 ===
func ExampleHeapProfilerEnable_LeakDetection() {
	// 启动泄漏监控
	_, err := CDPHeapProfilerEnable()
	if err != nil {
		log.Printf("启动内存监控失败: %v", err)
		return
	}
	log.Println("已启动堆分析，开始检测内存泄漏")
}

// === 使用示例3：自动化测试前置启动 ===
func ExampleHeapProfilerEnable_AutoTest() {
	// 测试用例执行前开启堆分析
	resp, err := CDPHeapProfilerEnable()
	if err != nil {
		log.Fatalf("测试前置开启堆分析失败: %v", err)
	}
	log.Println("自动化测试：HeapProfiler 已就绪", resp)
}
*/

// -----------------------------------------------  HeapProfiler.getHeapObjectId  -----------------------------------------------
// === 应用场景 ===
// 1. 对象追踪：通过运行时对象ID获取堆对象ID，用于内存分析
// 2. 内存调试：关联Runtime对象与Heap对象，进行完整引用链排查
// 3. 泄漏定位：获取对象堆ID后，监控其内存回收状态
// 4. 自动化分析：在自动化测试中获取对象唯一堆标识
// 5. 堆快照辅助：为指定对象生成堆ID，用于快照中精准查找
// 6. 性能监控：跟踪核心业务对象的堆内存分配

// CDPHeapProfilerGetHeapObjectId 根据运行时对象ID获取堆对象ID
// 参数：objectId - Runtime.getProperties 等返回的运行时对象ID
func CDPHeapProfilerGetHeapObjectId(objectId string) (string, error) {
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
		"method": "HeapProfiler.getHeapObjectId",
		"params": {
			"objectId": "%s"
		}
	}`, reqID, objectId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getHeapObjectId 请求失败: %w", err)
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
			return "", fmt.Errorf("getHeapObjectId 请求超时")
		}
	}
}

/*
// === 使用示例1：基础获取堆对象ID ===
func ExampleGetHeapObjectId_Base() {
	// 传入从 Runtime 接口获取的 objectId
	runtimeObjId := "{\"injectedScriptId\":1,\"id\":2}"
	resp, err := CDPHeapProfilerGetHeapObjectId(runtimeObjId)
	if err != nil {
		log.Fatalf("获取堆对象ID失败: %v", err)
	}
	log.Printf("堆对象ID获取成功: %s", resp)
}

// === 使用示例2：获取ID后加入堆监控 ===
func ExampleGetHeapObjectId_AddInspected() {
	// 1. 获取运行时对象ID
	objId := "{\"injectedScriptId\":1,\"id\":100}"

	// 2. 获取堆ID
	resp, err := CDPHeapProfilerGetHeapObjectId(objId)
	if err != nil {
		log.Printf("获取堆ID失败: %v", err)
		return
	}
	log.Println("获取HeapObjectId成功:", resp)

	// 3. 加入监控列表（配合 addInspectedHeapObject 使用）
	// 解析 resp 获取 heapObjectId 后调用
	// CDPHeapProfilerAddInspectedHeapObject(heapObjectId)
}

// === 使用示例3：内存泄漏排查流程 ===
func ExampleGetHeapObjectId_LeakDebug() {
	// 目标对象ID
	targetObj := "{\"injectedScriptId\":1,\"id\":5}"
	_, err := CDPHeapProfilerGetHeapObjectId(targetObj)
	if err != nil {
		log.Fatalf("泄漏排查失败: %v", err)
	}
	log.Println("已获取目标对象堆ID，可进行引用链分析")
}
*/

// -----------------------------------------------  HeapProfiler.getObjectByHeapObjectId  -----------------------------------------------
// === 应用场景 ===
// 1. 堆对象还原：通过堆对象ID获取原始JS对象，用于内存分析
// 2. 泄漏对象查看：查看内存泄漏的堆对象具体内容与属性
// 3. 运行时关联：将堆分析数据还原为可操作的Runtime对象
// 4. 内存快照验证：验证快照中的对象是否为业务关键对象
// 5. 自动化调试：自动获取泄漏对象详情，生成调试报告
// 6. 引用链分析：获取对象后分析其属性、子对象与引用关系

// CDPHeapProfilerGetObjectByHeapObjectId 根据堆对象ID获取运行时对象
// 参数：heapObjectId - 堆对象ID
func CDPHeapProfilerGetObjectByHeapObjectId(heapObjectId string) (string, error) {
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
		"method": "HeapProfiler.getObjectByHeapObjectId",
		"params": {
			"heapObjectId": "%s"
		}
	}`, reqID, heapObjectId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getObjectByHeapObjectId 请求失败: %w", err)
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
			return "", fmt.Errorf("getObjectByHeapObjectId 请求超时")
		}
	}
}

/*
// === 使用示例1：基础获取堆对象详情 ===
func ExampleGetObjectByHeapObjectId_Base() {
	// 传入已知堆ID
	heapId := "heap://123456"
	resp, err := CDPHeapProfilerGetObjectByHeapObjectId(heapId)
	if err != nil {
		log.Fatalf("获取堆对象失败: %v", err)
	}
	log.Printf("堆对象详情: %s", resp)
}

// === 使用示例2：内存泄漏对象分析 ===
func ExampleGetObjectByHeapObjectId_Leak() {
	// 泄漏对象堆ID
	leakHeapId := "heap://leak-obj-789"
	resp, err := CDPHeapProfilerGetObjectByHeapObjectId(leakHeapId)
	if err != nil {
		log.Printf("分析泄漏对象失败: %v", err)
		return
	}
	log.Println("泄漏对象内容获取成功:", resp)
}

// === 使用示例3：配合getHeapObjectId完整流程 ===
func ExampleGetObjectByHeapObjectId_Full() {
	// 1. 获取运行时对象堆ID
	runtimeObjId := "{\"injectedScriptId\":1,\"id\":2}"
	heapResp, _ := CDPHeapProfilerGetHeapObjectId(runtimeObjId)

	// 2. 解析堆ID
	var heapResult struct {
		Result struct {
			HeapObjectId string `json:"heapObjectId"`
		} `json:"result"`
	}
	json.Unmarshal([]byte(heapResp), &heapResult)
	heapId := heapResult.Result.HeapObjectId

	// 3. 获取完整对象
	resp, err := CDPHeapProfilerGetObjectByHeapObjectId(heapId)
	if err != nil {
		log.Fatalf("获取对象失败: %v", err)
	}
	log.Println("对象获取成功:", resp)
}

*/

// -----------------------------------------------  HeapProfiler.getSamplingProfile  -----------------------------------------------
// === 应用场景 ===
// 1. 内存采样分析：获取当前堆内存采样数据，分析内存分配热点
// 2. 性能瓶颈定位：通过采样数据找出占用内存最高的函数/代码块
// 3. 轻量级内存监控：非侵入式采样，不影响程序运行性能
// 4. 自动化性能测试：生成内存采样报告，用于性能回归检测
// 5. 线上问题排查：生产环境安全获取内存数据，定位内存异常
// 6. 长期内存趋势分析：定期采样，观察内存分配变化趋势

// CDPHeapProfilerGetSamplingProfile 获取堆内存采样分析报告
func CDPHeapProfilerGetSamplingProfile() (string, error) {
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
		"method": "HeapProfiler.getSamplingProfile"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getSamplingProfile 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
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
			return "", fmt.Errorf("getSamplingProfile 请求超时")
		}
	}
}

/*
// === 使用示例1：基础获取内存采样报告 ===
func ExampleGetSamplingProfile_Base() {
	// 直接获取当前堆采样数据
	resp, err := CDPHeapProfilerGetSamplingProfile()
	if err != nil {
		log.Fatalf("获取堆采样数据失败: %v", err)
	}
	log.Printf("堆采样报告: %s", resp)
}

// === 使用示例2：性能测试后采集内存数据 ===
func ExampleGetSamplingProfile_PerfTest() {
	// 执行业务逻辑后采集内存采样
	log.Println("执行压力测试...")
	// 测试代码...

	// 采集分析
	resp, err := CDPHeapProfilerGetSamplingProfile()
	if err != nil {
		log.Printf("测试后采样失败: %v", err)
		return
	}
	log.Println("性能测试内存采样完成:", resp)
}

// === 使用示例3：定期采集内存分析趋势 ===
func ExampleGetSamplingProfile_Periodic() {
	// 每30秒采集一次内存数据
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		resp, err := CDPHeapProfilerGetSamplingProfile()
		if err != nil {
			log.Printf("定时采样失败: %v", err)
			continue
		}
		log.Println("内存采样数据已获取")
		// 可保存到文件/数据库做趋势分析
	}
}
*/

// -----------------------------------------------  HeapProfiler.startSampling  -----------------------------------------------
// === 应用场景 ===
// 1. 内存采样监控：启动堆内存采样，开始记录内存分配情况
// 2. 性能测试前置：在性能/压力测试前启动采样，追踪内存分配热点
// 3. 内存泄漏排查：轻量级采样，长期运行检测内存增长趋势
// 4. 自动化测试：测试流程中启动采样，获取函数级内存分配数据
// 5. 线上性能分析：低开销采样，不影响页面正常运行
// 6. 代码内存优化：定位高频内存分配的代码位置

// CDPHeapProfilerStartSampling 启动堆内存采样
// 参数：samplingInterval - 采样间隔（单位：字节，默认 1048576 字节=1MB，传 0 使用默认值）
func CDPHeapProfilerStartSampling(samplingInterval int) (string, error) {
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
		"method": "HeapProfiler.startSampling",
		"params": {
			"samplingInterval": %d
		}
	}`, reqID, samplingInterval)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 startSampling 请求失败: %w", err)
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
// === 使用示例1：默认采样间隔启动 ===
func ExampleStartSampling_Default() {
	// 使用默认 1MB 采样间隔
	resp, err := CDPHeapProfilerStartSampling(0)
	if err != nil {
		log.Fatalf("启动堆采样失败: %v", err)
	}
	log.Println("堆内存采样已启动:", resp)
}

// === 使用示例2：高精度采样（性能测试专用） ===
func ExampleStartSampling_HighPrecision() {
	// 高精度：间隔 4096 字节，适合详细分析
	resp, err := CDPHeapProfilerStartSampling(4096)
	if err != nil {
		log.Fatalf("启动高精度采样失败: %v", err)
	}
	log.Println("高精度内存采样已启动:", resp)
}

// === 使用示例3：性能测试流程启动 ===
func ExampleStartSampling_PerfTest() {
	// 测试前启动采样
	_, err := CDPHeapProfilerStartSampling(102400)
	if err != nil {
		log.Printf("测试采样启动失败: %v", err)
		return
	}
	log.Println("性能测试：内存采样已开始，可执行业务逻辑")
	// 执行业务代码 → 调用 getSamplingProfile 获取结果 → stopSampling
}
*/

// -----------------------------------------------  HeapProfiler.startTrackingHeapObjects  -----------------------------------------------
// === 应用场景 ===
// 1. 内存泄漏精准检测：启动实时堆对象跟踪，捕获对象分配与回收轨迹
// 2. 实时内存监控：持续追踪堆内存变化，实时获取内存增长详情
// 3. 泄漏对象定位：跟踪未回收对象，快速定位泄漏源
// 4. 自动化内存诊断：测试流程中启动跟踪，自动分析内存异常
// 5. 页面内存生命周期监控：全程监控页面加载到运行的内存状态
// 6. 重度页面优化：针对复杂页面跟踪对象引用，优化内存占用

// CDPHeapProfilerStartTrackingHeapObjects 启动堆对象实时跟踪
// 参数：trackAllocations - 是否跟踪对象分配（true=开启详细分配跟踪，false=基础跟踪）
func CDPHeapProfilerStartTrackingHeapObjects(trackAllocations bool) (string, error) {
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
		"method": "HeapProfiler.startTrackingHeapObjects",
		"params": {
			"trackAllocations": %t
		}
	}`, reqID, trackAllocations)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 startTrackingHeapObjects 请求失败: %w", err)
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
			return "", fmt.Errorf("startTrackingHeapObjects 请求超时")
		}
	}
}

/*
// === 使用示例1：基础堆对象跟踪（低性能消耗） ===
func ExampleStartTrackingHeapObjects_Base() {
	// 关闭分配跟踪，性能消耗低
	resp, err := CDPHeapProfilerStartTrackingHeapObjects(false)
	if err != nil {
		log.Fatalf("启动堆对象跟踪失败: %v", err)
	}
	log.Println("堆对象基础跟踪已启动:", resp)
}

// === 使用示例2：精准内存泄漏检测（完整跟踪） ===
func ExampleStartTrackingHeapObjects_LeakDetection() {
	// 开启分配跟踪，精准检测泄漏
	resp, err := CDPHeapProfilerStartTrackingHeapObjects(true)
	if err != nil {
		log.Fatalf("启动泄漏检测模式失败: %v", err)
	}
	log.Println("内存泄漏精准检测已启动，可记录完整对象分配:", resp)
}

// === 使用示例3：自动化测试内存监控 ===
func ExampleStartTrackingHeapObjects_AutoTest() {
	// 测试前启动完整跟踪
	_, err := CDPHeapProfilerStartTrackingHeapObjects(true)
	if err != nil {
		log.Printf("测试内存跟踪启动失败: %v", err)
		return
	}
	log.Println("自动化测试：堆对象跟踪已开启，运行测试用例...")
	// 测试完成后调用 stopTrackingHeapObjects + takeHeapSnapshot
}

*/

// -----------------------------------------------  HeapProfiler.stopSampling  -----------------------------------------------
// === 应用场景 ===
// 1. 结束内存采样：停止堆内存采样，释放采样相关资源
// 2. 性能测试收尾：性能测试完成后停止采样，避免资源占用
// 3. 内存分析流程结束：完成内存分配热点分析后关闭采样
// 4. 自动化测试清理：测试用例结束后清理采样状态
// 5. 性能恢复：停止采样后浏览器恢复正常运行性能
// 6. 报告生成：停止采样后获取最终完整的内存分析报告

// CDPHeapProfilerStopSampling 停止堆内存采样
func CDPHeapProfilerStopSampling() (string, error) {
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
		"method": "HeapProfiler.stopSampling"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stopSampling 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
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
			return "", fmt.Errorf("stopSampling 请求超时")
		}
	}
}

/*
// === 使用示例1：标准采样结束流程 ===
func ExampleStopSampling_Normal() {
	// 先获取采样报告
	_, _ = CDPHeapProfilerGetSamplingProfile()

	// 再停止采样
	resp, err := CDPHeapProfilerStopSampling()
	if err != nil {
		log.Fatalf("停止内存采样失败: %v", err)
	}
	log.Println("内存采样已停止:", resp)
}

// === 使用示例2：性能测试完成清理 ===
func ExampleStopSampling_PerfTest() {
	// 测试业务逻辑执行完毕
	log.Println("性能测试结束，停止内存采样...")

	resp, err := CDPHeapProfilerStopSampling()
	if err != nil {
		log.Printf("测试后停止采样失败: %v", err)
		return
	}
	log.Println("采样已停止，资源已释放:", resp)
}

// === 使用示例3：完整内存分析流程 ===
func ExampleStopSampling_FullProcess() {
	// 1. 启动采样
	CDPHeapProfilerStartSampling(0)
	// 2. 执行业务操作
	log.Println("运行测试代码...")
	// 3. 获取报告
	CDPHeapProfilerGetSamplingProfile()
	// 4. 停止采样
	resp, _ := CDPHeapProfilerStopSampling()
	log.Println("完整内存分析流程完成，采样已停止:", resp)
}
*/

// -----------------------------------------------  HeapProfiler.stopTrackingHeapObjects  -----------------------------------------------
// === 应用场景 ===
// 1. 停止堆对象跟踪：结束堆对象实时监控，释放浏览器分析资源
// 2. 内存泄漏检测结束：完成泄漏排查后关闭跟踪，恢复浏览器性能
// 3. 自动化测试清理：测试用例执行完毕后清理内存跟踪状态
// 4. 长时间监控收尾：长期内存监控结束，停止数据采集
// 5. 浏览器资源释放：关闭高消耗的堆跟踪功能，避免CPU/内存占用
// 6. 内存报告生成：停止跟踪后生成完整的对象分配/回收报告

// CDPHeapProfilerStopTrackingHeapObjects 停止堆对象跟踪
// 参数：reportProgress - 是否在停止时上报进度信息（true=上报，false=直接停止）
func CDPHeapProfilerStopTrackingHeapObjects(reportProgress bool) (string, error) {
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
		"method": "HeapProfiler.stopTrackingHeapObjects",
		"params": {
			"reportProgress": %t
		}
	}`, reqID, reportProgress)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stopTrackingHeapObjects 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
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
			return "", fmt.Errorf("stopTrackingHeapObjects 请求超时")
		}
	}
}

/*
// === 使用示例1：标准停止跟踪（带进度上报） ===
func ExampleStopTrackingHeapObjects_Default() {
	// 停止并上报进度，适合调试、完整分析场景
	resp, err := CDPHeapProfilerStopTrackingHeapObjects(true)
	if err != nil {
		log.Fatalf("停止堆对象跟踪失败: %v", err)
	}
	log.Println("堆对象跟踪已停止，进度已上报:", resp)
}

// === 使用示例2：快速停止（无进度上报） ===
func ExampleStopTrackingHeapObjects_Fast() {
	// 直接停止，不等待进度上报，适合自动化、快速清理场景
	resp, err := CDPHeapProfilerStopTrackingHeapObjects(false)
	if err != nil {
		log.Fatalf("快速停止跟踪失败: %v", err)
	}
	log.Println("堆对象跟踪已快速停止:", resp)
}

// === 使用示例3：内存泄漏检测完整流程 ===
func ExampleStopTrackingHeapObjects_Full() {
	// 1. 启动跟踪
	CDPHeapProfilerStartTrackingHeapObjects(true)
	// 2. 执行业务逻辑、页面操作
	log.Println("正在检测内存泄漏...")
	// 3. 停止跟踪并生成报告
	resp, err := CDPHeapProfilerStopTrackingHeapObjects(true)
	if err != nil {
		log.Printf("停止失败: %v", err)
		return
	}
	log.Println("内存泄漏检测完成，跟踪已停止:", resp)
}
*/

// -----------------------------------------------  HeapProfiler.takeHeapSnapshot  -----------------------------------------------
// === 应用场景 ===
// 1. 内存快照采集：生成完整的堆内存快照，用于离线分析内存占用
// 2. 内存泄漏排查：通过快照对比找出未释放的对象，定位泄漏根源
// 3. 页面性能诊断：分析页面DOM、JS对象、闭包等内存占用情况
// 4. 自动化测试：测试流程中采集快照，自动检测内存异常
// 5. 问题复现与存档：保存内存现场，用于后续深度调试
// 6. 服务端无头浏览器：定期生成快照监控长期运行内存状态

// CDPHeapProfilerTakeHeapSnapshot 生成堆内存快照
// 参数：reportProgress - 是否开启进度上报（true=实时返回采集进度，false=静默采集）
func CDPHeapProfilerTakeHeapSnapshot(reportProgress bool) (string, error) {
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
		"method": "HeapProfiler.takeHeapSnapshot",
		"params": {
			"reportProgress": %t
		}
	}`, reqID, reportProgress)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 takeHeapSnapshot 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（快照生成较慢，超时设为20秒）
	timeout := 20 * time.Second
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
			return "", fmt.Errorf("takeHeapSnapshot 请求超时")
		}
	}
}

/*
// === 使用示例1：标准采集快照（带进度） ===
func ExampleTakeHeapSnapshot_WithProgress() {
	// 开启进度上报，可监听采集百分比
	resp, err := CDPHeapProfilerTakeHeapSnapshot(true)
	if err != nil {
		log.Fatalf("堆快照生成失败: %v", err)
	}
	log.Println("堆快照生成完成:", resp)
}

// === 使用示例2：静默采集（无进度） ===
func ExampleTakeHeapSnapshot_Silent() {
	// 不监听进度，直接生成快照
	resp, err := CDPHeapProfilerTakeHeapSnapshot(false)
	if err != nil {
		log.Fatalf("静默快照生成失败: %v", err)
	}
	log.Println("静默堆快照已生成:", resp)
}

// === 使用示例3：GC后采集精准快照 ===
func ExampleTakeHeapSnapshot_AfterGC() {
	// 1. 先手动GC清理垃圾
	_, _ = CDPHeapProfilerCollectGarbage()
	log.Println("GC完成，开始生成堆快照...")

	// 2. 生成快照（带进度）
	resp, err := CDPHeapProfilerTakeHeapSnapshot(true)
	if err != nil {
		log.Fatalf("GC后快照生成失败: %v", err)
	}
	log.Println("GC后精准堆快照已生成:", resp)
}
*/
