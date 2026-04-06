package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Log.clear  -----------------------------------------------
// === 应用场景 ===
// 1. 日志环境重置: 清空旧日志，避免干扰新一轮日志采集
// 2. 自动化测试清理: 每个测试用例执行前清空日志，保证日志纯净
// 3. 调试流程初始化: 调试问题前清空历史日志，只关注最新输出
// 4. 页面切换清理: 切换页面/模块时清空上一个页面的残留日志
// 5. 错误排查准备: 定位问题前清空无关日志，精准捕获目标错误
// 6. 长时间运行程序: 定期清空日志缓存，防止内存占用过高

// CDPLogClear 清空浏览器控制台日志
func CDPLogClear() (string, error) {
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
		"method": "Log.clear"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Log.clear 请求失败: %w", err)
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
			return "", fmt.Errorf("Log.clear 请求超时")
		}
	}
}

/*

// ==================== Log.clear 使用示例 ====================
func ExampleCDPLogClear() {
	// ========== 示例1：基础清空日志 ==========
	resp, err := CDPLogClear()
	if err != nil {
		log.Fatalf("清空日志失败: %v", err)
	}
	log.Printf("清空浏览器日志成功，响应: %s", resp)

	// ========== 示例2：自动化测试标准流程 ==========
	// 每个测试用例开始前清空日志
	func TestPageLog(t *testing.T) {
		// 前置：清空历史日志
		_, err := CDPLogClear()
		if err != nil {
			t.Fatalf("测试前置清空日志失败: %v", err)
		}
		t.Log("已清空日志，开始执行测试用例")

		// 执行页面操作...
		// 捕获操作产生的新日志
	}

	// ========== 示例3：问题调试前清空 ==========
	// 调试JS报错前，清空无关日志
	debugBefore() {
		// 清空日志
		CDPLogClear()
		log.Println("已清空日志，开始复现问题并捕获日志")
		// 复现问题...
	}

	// ========== 示例4：配合Log.enable使用 ==========
	// 启用日志 → 清空 → 采集
	CDPLogEnable()
	defer CDPLogDisable()

	// 清空旧日志
	CDPLogClear()
	log.Println("开始监听全新日志")

}

*/

// -----------------------------------------------  Log.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 停止日志监听: 不再接收console、报错、警告等日志推送
// 2. 资源释放: 关闭日志模块，释放浏览器和程序占用的内存资源
// 3. 测试流程收尾: 自动化测试完成后关闭日志监听，清理环境
// 4. 性能优化: 不需要日志时关闭，减少程序消息处理压力
// 5. 多任务切换: 切换不同调试任务时，关闭上一个任务的日志监听
// 6. 程序退出清理: 程序结束时强制关闭日志，避免资源泄漏

// CDPLogDisable 关闭日志监听功能
func CDPLogDisable() (string, error) {
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
		"method": "Log.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Log.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("Log.disable 请求超时")
		}
	}
}

/*

// ==================== Log.disable 使用示例 ====================
func ExampleCDPLogDisable() {
	// ========== 示例1：基础关闭日志监听 ==========
	resp, err := CDPLogDisable()
	if err != nil {
		log.Fatalf("关闭日志功能失败: %v", err)
	}
	log.Printf("关闭日志功能成功，响应: %s", resp)

	// ========== 示例2：标准日志采集完整流程 ==========
	// 1. 启用日志
	CDPLogEnable()
	// 2. 清空旧日志
	CDPLogClear()

	// 3. 执行业务操作、采集日志...
	log.Println("正在采集页面日志...")

	// 4. 采集完成，关闭日志（推荐用defer确保执行）
	defer func() {
		_, err := CDPLogDisable()
		if err != nil {
			log.Printf("关闭日志失败: %v", err)
		} else {
			log.Println("日志功能已关闭，资源已释放")
		}
	}()

	// ========== 示例3：自动化测试用例收尾 ==========
	func TestPageErrorLog(t *testing.T) {
		// 测试前启用
		CDPLogEnable()
		// 测试后无论如何都关闭
		defer CDPLogDisable()

		// 执行测试...
		t.Log("测试完成，自动关闭日志监听")
	}
}


*/

// -----------------------------------------------  Log.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 页面日志采集：开启监听，捕获console.log/error/warn等所有日志
// 2. 前端错误监控：实时获取JS报错、渲染异常、资源加载失败日志
// 3. 自动化测试：测试过程中全程监听日志，验证是否存在异常输出
// 4. 调试问题定位：复现BUG时，开启日志捕获关键报错信息
// 5. 性能日志分析：收集控制台性能、警告类日志进行优化
// 6. 线上问题复现：还原用户场景时，开启日志捕获完整异常链

// CDPLogEnable 启用浏览器日志监听
func CDPLogEnable() (string, error) {
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
		"method": "Log.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Log.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("Log.enable 请求超时")
		}
	}
}

/*


// ==================== Log.enable 使用示例 ====================
func ExampleCDPLogEnable() {
	// ========== 示例1：基础启用日志监听 ==========
	resp, err := CDPLogEnable()
	if err != nil {
		log.Fatalf("启用日志监听失败: %v", err)
	}
	log.Printf("启用日志功能成功，响应: %s", resp)

	// ========== 示例2：标准日志监听完整流程（推荐） ==========
	// 1. 启用日志（必须第一步）
	_, err := CDPLogEnable()
	if err != nil {
		log.Fatalf("启用Log失败: %v", err)
	}
	log.Println("已启用日志监听，开始捕获浏览器输出")

	// 2. 清空历史日志，避免干扰
	CDPLogClear()

	// 3. 执行业务操作、页面操作...

	// 4. 任务完成后关闭日志（defer确保一定执行）
	defer CDPLogDisable()

	// ========== 示例3：自动化错误捕获 ==========
	// 测试用例开始时启用日志
	func TestPageConsoleError(t *testing.T) {
		// 前置：开启日志
		_, err := CDPLogEnable()
		if err != nil {
			t.Fatal(err)
		}
		// 结束自动关闭
		defer CDPLogDisable()

		// 执行页面操作...
		// 捕获console.error、JS异常等
		t.Log("正在监听页面所有日志与错误")
	}
}

*/

// -----------------------------------------------  Log.startViolationsReport  -----------------------------------------------
// === 应用场景 ===
// 1. 性能违规监控：实时捕获长任务、布局抖动、强制同步布局等性能问题
// 2. 自动化性能测试：测试过程中监听违规行为，生成性能缺陷报告
// 3. 页面卡顿排查：定位导致页面掉帧、卡顿的核心违规操作
// 4. 前端代码审计：检测不符合性能规范的JS/CSS操作
// 5. 长任务分析：捕获执行时间过长的JS任务，优化执行效率
// 6. 渲染优化验证：优化后监听违规，验证性能问题是否修复

// CDPLogStartViolationsReport 启动性能违规报告监听
// config：违规配置数组，包含违规类型与阈值
func CDPLogStartViolationsReport(config string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息（config直接传入JSON格式配置）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Log.startViolationsReport",
		"params": {
			"config": %s
		}
	}`, reqID, config)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Log.startViolationsReport 请求失败: %w", err)
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
			return "", fmt.Errorf("Log.startViolationsReport 请求超时")
		}
	}
}

/*

// ==================== Log.startViolationsReport 使用示例 ====================
func ExampleCDPLogStartViolationsReport() {
	// ========== 示例1：基础启动 - 监听长任务+强制同步布局 ==========
	// 标准违规配置JSON：监听长任务(>200ms)、强制同步布局
	violationConfig := `[
		{"name":"longTask","threshold":200},
		{"name":"forcedLayout","threshold":10}
	]`

	resp, err := CDPLogStartViolationsReport(violationConfig)
	if err != nil {
		log.Fatalf("启动违规报告失败: %v", err)
	}
	log.Printf("启动性能违规监控成功，响应: %s", resp)

	// ========== 示例2：完整性能监控流程（推荐） ==========
	// 1. 启用日志模块
	CDPLogEnable()
	defer CDPLogDisable()

	// 2. 配置全量违规检测
	fullConfig := `[
		{"name":"longTask","threshold":150},
		{"name":"forcedLayout","threshold":5},
		{"name":"styleRecalculation","threshold":50},
		{"name":"longLayout","threshold":20}
	]`

	// 3. 启动违规报告
	_, err := CDPLogStartViolationsReport(fullConfig)
	if err != nil {
		log.Fatalf("性能监控启动失败: %v", err)
	}
	log.Println("已启动全量性能违规监控，开始捕获页面问题")

	// 4. 执行页面操作，自动捕获违规日志

	// ========== 示例3：自动化测试性能审计 ==========
	func TestPagePerformance(t *testing.T) {
		// 前置启用日志
		CDPLogEnable()
		defer CDPLogDisable()

		// 启动违规监控
		config := `[{"name":"longTask","threshold":200}]`
		CDPLogStartViolationsReport(config)

		// 执行业务操作...
		// 断言：无长任务违规
		t.Log("性能审计中：禁止执行超过200ms的长任务")
	}
}

*/

// -----------------------------------------------  Log.stopViolationsReport  -----------------------------------------------
// === 应用场景 ===
// 1. 停止性能监控: 结束违规行为监听，不再捕获长任务、布局抖动等问题
// 2. 资源释放: 关闭违规监控，释放浏览器性能分析占用的内存资源
// 3. 测试流程收尾: 自动化性能测试完成后停止监控，清理环境
// 4. 多任务切换: 切换不同测试/调试任务时，停止上一个性能监控任务
// 5. 性能报告生成: 采集完成后停止，生成最终性能违规报告
// 6. 程序退出清理: 程序结束时强制停止，防止资源泄漏与无效监听

// CDPLogStopViolationsReport 停止性能违规报告监听
func CDPLogStopViolationsReport() (string, error) {
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
		"method": "Log.stopViolationsReport"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Log.stopViolationsReport 请求失败: %w", err)
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
			return "", fmt.Errorf("Log.stopViolationsReport 请求超时")
		}
	}
}

/*


// ==================== Log.stopViolationsReport 使用示例 ====================
func ExampleCDPLogStopViolationsReport() {
	// ========== 示例1：基础停止违规监控 ==========
	resp, err := CDPLogStopViolationsReport()
	if err != nil {
		log.Fatalf("停止性能违规报告失败: %v", err)
	}
	log.Printf("停止性能违规监控成功，响应: %s", resp)

	// ========== 示例2：完整性能审计标准流程（推荐） ==========
	// 1. 启用日志
	CDPLogEnable()
	defer CDPLogDisable()

	// 2. 配置违规规则
	config := `[{"name":"longTask","threshold":200},{"name":"forcedLayout","threshold":10}]`

	// 3. 启动监控
	CDPLogStartViolationsReport(config)
	log.Println("已启动性能违规监控，开始采集...")

	// 4. 执行页面操作、性能测试...

	// 5. 采集完成，停止监控（defer确保一定执行）
	defer func() {
		_, stopErr := CDPLogStopViolationsReport()
		if stopErr != nil {
			log.Printf("停止监控失败: %v", stopErr)
		} else {
			log.Println("性能违规监控已停止，资源已释放")
		}
	}()

	// ========== 示例3：自动化测试收尾 ==========
	func TestPerformanceViolation(t *testing.T) {
		// 前置
		CDPLogEnable()
		CDPLogStartViolationsReport(`[{"name":"longTask","threshold":150}]`)

		// 测试结束自动停止
		defer CDPLogStopViolationsReport()
		defer CDPLogDisable()

		// 执行测试逻辑...
		t.Log("测试中，违规监控已自动启停")
	}
}

*/
