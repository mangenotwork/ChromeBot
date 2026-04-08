package browser

import (
	"ChromeBot/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Page.addScriptToEvaluateOnNewDocument  -----------------------------------------------
// === 应用场景 ===
// 1. 全局变量注入: 在页面加载前注入全局变量、工具函数
// 2. 反爬绕过: 提前修改页面环境、覆盖浏览器指纹特征
// 3. 初始化配置: 为单页应用提前设置全局配置、接口地址
// 4. 日志/监控: 注入全局日志、错误捕获、性能监控脚本
// 5. 自动化前置: 自动化测试前注入测试工具、mock函数
// 6. 页面劫持: 提前拦截fetch/XMLHttpRequest/console等原生方法

// CDPPageAddScriptToEvaluateOnNewDocument 在每个新文档加载前注入脚本
// 参数 source: 需要注入的JavaScript脚本字符串
// 返回值: 脚本注入ID(用于后续移除)、响应内容、错误信息
func CDPPageAddScriptToEvaluateOnNewDocument(source string) (string, string, error) {
	if !DefaultBrowserWS() {
		return "", "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.addScriptToEvaluateOnNewDocument",
		"params": {
			"source": %s
		}
	}`, reqID, strconv.Quote(source))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", "", fmt.Errorf("发送 addScriptToEvaluateOnNewDocument 请求失败: %w", err)
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
				return "", "", fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return "", content, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return "", content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				// 提取脚本ID
				if result, ok := response["result"].(map[string]interface{}); ok {
					if scriptID, ok := result["identifier"].(string); ok {
						return scriptID, content, nil
					}
				}

				return "", content, nil
			}

		case <-timer.C:
			return "", "", fmt.Errorf("addScriptToEvaluateOnNewDocument 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：注入全局变量 ===
func ExampleInjectGlobalVariable() {
	// 注入全局变量 window.ENV = "test"
	script := `window.ENV = "test";
			   window.API_HOST = "https://api.test.com";`
	scriptID, resp, err := CDPPageAddScriptToEvaluateOnNewDocument(script)
	if err != nil {
		log.Fatalf("注入全局变量失败: %v, 响应: %s", err, resp)
	}
	log.Printf("全局变量注入成功，脚本ID: %s", scriptID)
}

// === 使用场景示例代码：注入工具函数 ===
func ExampleInjectUtilFunction() {
	// 注入全局格式化时间函数
	script := `window.formatDate = function(timestamp) {
		const date = new Date(timestamp);
		return date.toLocaleString();
	};`
	scriptID, resp, err := CDPPageAddScriptToEvaluateOnNewDocument(script)
	if err != nil {
		log.Fatalf("注入工具函数失败: %v, 响应: %s", err, resp)
	}
	log.Printf("工具函数注入成功，脚本ID: %s", scriptID)
}

// === 使用场景示例代码：拦截console.log ===
func ExampleInterceptConsole() {
	// 提前拦截控制台输出
	script := `const originLog = console.log;
	console.log = function(...args) {
		originLog.apply(console, ["[注入日志]:", ...args]);
	};`
	scriptID, resp, err := CDPPageAddScriptToEvaluateOnNewDocument(script)
	if err != nil {
		log.Fatalf("拦截console失败: %v, 响应: %s", err, resp)
	}
	log.Printf("console拦截注入成功，脚本ID: %s", scriptID)
}

// === 使用场景示例代码：修改浏览器指纹(webdriver) ===
func ExampleBypassWebdriver() {
	// 绕过无头浏览器webdriver检测
	script := `Object.defineProperty(navigator, 'webdriver', {
		get: () => false
	});`
	scriptID, resp, err := CDPPageAddScriptToEvaluateOnNewDocument(script)
	if err != nil {
		log.Fatalf("修改webdriver失败: %v, 响应: %s", err, resp)
	}
	log.Printf("webdriver修改成功，脚本ID: %s", scriptID)
}

*/

// -----------------------------------------------  Page.bringToFront  -----------------------------------------------
// === 应用场景 ===
// 1. 窗口置顶: 将浏览器窗口/标签页切换到前台最顶层显示
// 2. 自动化操作: 自动化测试时确保页面可见，避免操作失败
// 3. 调试辅助: 调试时快速将目标页面切换到前台
// 4. 多标签管理: 多标签页场景下切换到指定操作标签页
// 5. 交互触发: 需要页面可见才能触发的交互/弹窗前置操作
// 6. 监控告警: 异常时自动将监控页面切换到前台提醒

// CDPPageBringToFront 将页面切换到前台置顶显示
func CDPPageBringToFront() (string, error) {
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
		"method": "Page.bringToFront"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 bringToFront 请求失败: %w", err)
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
			return "", fmt.Errorf("bringToFront 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：自动化操作前置顶页面 ===
func ExampleBringToFrontBeforeAutoTest() {
	// 自动化点击/输入操作前，确保页面在前台
	resp, err := CDPPageBringToFront()
	if err != nil {
		log.Fatalf("页面置顶失败: %v, 响应: %s", err, resp)
	}
	log.Println("页面已成功切换到前台，可执行后续自动化操作")
}

// === 使用场景示例代码：多标签页切换到目标页 ===
func ExampleSwitchTabInMultiPages() {
	// 多标签操作时，切换到需要操作的标签页
	resp, err := CDPPageBringToFront()
	if err != nil {
		log.Printf("标签页切换失败: %v", err)
		return
	}
	log.Println("已切换到目标标签页并置顶显示")
}

// === 使用场景示例代码：异常告警时自动前台展示 ===
func ExampleShowPageWhenError() {
	// 业务监控发现异常，自动将页面置顶提醒
	errOccur := true
	if errOccur {
		resp, err := CDPPageBringToFront()
		if err != nil {
			log.Printf("异常页面展示失败: %v", err)
			return
		}
		log.Println("⚠️ 异常告警：已将监控页面切换到前台")
	}
}

*/

// -----------------------------------------------  Page.captureScreenshot  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试截图: 测试用例执行失败时截取页面现场
// 2. 页面快照留存: 关键页面/数据截图永久保存
// 3. 报表生成: 生成页面可视化报表图片
// 4. 异常监控取证: 页面报错、异常时自动取证截图
// 5. 内容分享: 将页面内容转为图片用于分享
// 6. 定时截图监控: 定时截取页面监控状态变化

// CDPPageCaptureScreenshot 截取当前页面截图
// 参数 format: 图片格式 png、jpeg、webp
// 参数 quality: 图片质量 0-100 (仅jpeg/webp生效)
// 返回值: base64编码的图片数据、响应内容、错误信息
func CDPPageCaptureScreenshot(format string, quality int) (string, string, error) {
	if !DefaultBrowserWS() {
		return "", "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.captureScreenshot",
		"params": {
			"format": "%s",
			"quality": %d
		}
	}`, reqID, format, quality)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", "", fmt.Errorf("发送 captureScreenshot 请求失败: %w", err)
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
				return "", "", fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return "", content, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return "", content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				// 提取base64图片数据
				if result, ok := response["result"].(map[string]interface{}); ok {
					if data, ok := result["data"].(string); ok {
						return data, content, nil
					}
				}

				return "", content, nil
			}

		case <-timer.C:
			return "", "", fmt.Errorf("captureScreenshot 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：截取PNG高清截图并保存到本地 ===
func ExampleCapturePNGScreenshot() {
	// 截取PNG格式，质量100
	base64Data, resp, err := CDPPageCaptureScreenshot("png", 100)
	if err != nil {
		log.Fatalf("截图失败: %v, 响应: %s", err, resp)
	}

	// base64解码保存为文件
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Fatalf("base64解码失败: %v", err)
	}
	err = os.WriteFile("screenshot.png", imgData, 0644)
	if err != nil {
		log.Fatalf("保存截图失败: %v", err)
	}
	log.Println("PNG截图已保存：screenshot.png")
}

// === 使用场景示例代码：截取JPG压缩截图节省存储空间 ===
func ExampleCaptureJPGScreenshot() {
	// 截取JPEG格式，质量80（兼顾清晰度和体积）
	base64Data, resp, err := CDPPageCaptureScreenshot("jpeg", 80)
	if err != nil {
		log.Printf("截图失败: %v", err)
		return
	}

	imgData, _ := base64.StdEncoding.DecodeString(base64Data)
	_ = os.WriteFile("screenshot.jpg", imgData, 0644)
	log.Println("JPG截图已保存")
}

// === 使用场景示例代码：自动化测试失败自动取证截图 ===
func ExampleCaptureScreenshotOnTestFailed() {
	// 模拟测试失败
	testFailed := true
	if testFailed {
		base64Data, resp, err := CDPPageCaptureScreenshot("png", 90)
		if err != nil {
			log.Printf("失败截图失败: %v", err)
			return
		}

		// 以时间戳命名保存
		filename := fmt.Sprintf("test_fail_%d.png", time.Now().Unix())
		imgData, _ := base64.StdEncoding.DecodeString(base64Data)
		_ = os.WriteFile(filename, imgData, 0644)
		log.Printf("测试失败，已保存现场截图：%s", filename)
	}
}

// === 使用场景示例代码：定时截图页面监控 ===
func ExampleTimedCaptureScreenshot() {
	// 每30秒截取一次页面
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		base64Data, _, err := CDPPageCaptureScreenshot("webp", 85)
		if err != nil {
			log.Printf("定时截图失败: %v", err)
			continue
		}
		imgData, _ := base64.StdEncoding.DecodeString(base64Data)
		filename := fmt.Sprintf("monitor_%d.webp", time.Now().Unix())
		_ = os.WriteFile(filename, imgData, 0644)
		log.Printf("已保存定时监控截图：%s", filename)
	}
}

*/

// -----------------------------------------------  Page.close  -----------------------------------------------
// === 应用场景 ===
// 1. 标签页关闭: 关闭当前操作的浏览器标签页
// 2. 自动化清理: 自动化测试完成后关闭无用标签页释放资源
// 3. 多页面管理: 多标签场景下按需关闭指定页面
// 4. 任务结束回收: 爬虫/自动化任务完成后关闭页面避免内存占用
// 5. 错误重试清理: 任务失败后关闭异常页面重新初始化
// 6. 定时关闭: 超时无操作页面自动关闭

// CDPPageClose 关闭当前页面标签
func CDPPageClose() (string, error) {
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
		"method": "Page.close"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 close 请求失败: %w", err)
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
			return "", fmt.Errorf("close 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：自动化测试完成后关闭页面 ===
func ExampleCloseAfterTest() {
	// 执行测试逻辑...
	log.Println("自动化测试执行完成")

	// 关闭当前页面
	resp, err := CDPPageClose()
	if err != nil {
		log.Fatalf("关闭页面失败: %v, 响应: %s", err, resp)
	}
	log.Println("测试页面已成功关闭，资源已释放")
}

// === 使用场景示例代码：多标签页按需关闭 ===
func ExampleCloseTabInMultiPages() {
	// 多标签操作完成后关闭当前标签
	resp, err := CDPPageClose()
	if err != nil {
		log.Printf("关闭标签页失败: %v", err)
		return
	}
	log.Println("当前无用标签页已关闭")
}

// === 使用场景示例代码：任务完成后关闭页面 ===
func ExampleCloseAfterTaskFinish() {
	// 模拟数据采集/业务任务完成
	taskFinish := true
	if taskFinish {
		resp, err := CDPPageClose()
		if err != nil {
			log.Printf("任务完成关闭页面失败: %v", err)
			return
		}
		log.Println("✅ 任务已完成，页面已关闭")
	}
}

// === 使用场景示例代码：超时自动关闭页面 ===
func ExampleCloseOnTimeout() {
	// 页面操作超时（10秒）自动关闭
	time.Sleep(10 * time.Second)
	resp, err := CDPPageClose()
	if err != nil {
		log.Printf("超时关闭页面失败: %v", err)
		return
	}
	log.Println("⌛ 页面操作超时，已自动关闭")
}

*/

// -----------------------------------------------  Page.createIsolatedWorld  -----------------------------------------------
// === 应用场景 ===
// 1. 安全隔离执行：在独立沙箱环境执行JS，避免污染页面全局环境
// 2. 脚本隔离：注入的脚本与页面原生脚本隔离，防止冲突
// 3. 反爬/环境隔离：创建独立JS环境，绕过页面环境检测、防止变量覆盖
// 4. 自动化安全执行：测试脚本在隔离世界运行，不影响页面业务逻辑
// 5. 多环境执行：同时存在页面原生环境+隔离环境，互不干扰
// 6. 敏感操作执行：执行高危JS操作，避免影响页面正常运行

// CDPPageCreateIsolatedWorld 创建页面隔离环境（独立JS世界）
// 参数 frameId：目标框架ID（空则使用主框架）
// 参数 worldName：隔离世界名称（自定义标识）
// 参数 grantUniveralAccess：是否授予跨域访问权限（true=开启）
// 返回值：执行结果、响应内容、错误信息
func CDPPageCreateIsolatedWorld(frameId string, worldName string, grantUniveralAccess bool) (string, string, error) {
	if !DefaultBrowserWS() {
		return "", "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数（frameId 可选，为空则不传递）
	var params string
	if frameId == "" {
		params = fmt.Sprintf(`"worldName": "%s", "grantUniveralAccess": %t`, worldName, grantUniveralAccess)
	} else {
		params = fmt.Sprintf(`"frameId": "%s", "worldName": "%s", "grantUniveralAccess": %t`, frameId, worldName, grantUniveralAccess)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.createIsolatedWorld",
		"params": {
			%s
		}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", "", fmt.Errorf("发送 createIsolatedWorld 请求失败: %w", err)
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
				return "", "", fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return "", content, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return "", content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				// 提取执行结果
				if result, ok := response["result"].(map[string]interface{}); ok {
					resultJson, _ := json.Marshal(result)
					return string(resultJson), content, nil
				}

				return "", content, nil
			}

		case <-timer.C:
			return "", "", fmt.Errorf("createIsolatedWorld 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：创建主页面隔离环境（无跨域权限）===
func ExampleCreateNormalIsolatedWorld() {
	// 空frameId=主框架，自定义名称，无跨域权限
	result, resp, err := CDPPageCreateIsolatedWorld("", "my-isolated-world", false)
	if err != nil {
		log.Fatalf("创建隔离环境失败: %v, 响应: %s", err, resp)
	}
	log.Printf("隔离环境创建成功: %s", result)
}

// === 使用场景示例代码：创建带跨域权限的隔离环境 ===
func ExampleCreateCrossDomainWorld() {
	// 授予跨域权限，适用于跨域资源读取场景
	result, resp, err := CDPPageCreateIsolatedWorld("", "cross-domain-world", true)
	if err != nil {
		log.Printf("创建跨域隔离环境失败: %v", err)
		return
	}
	log.Printf("带跨域权限的隔离环境创建成功: %s", result)
}

// === 使用场景示例代码：子框架（iframe）创建隔离环境 ===
func ExampleCreateIsolatedWorldForIframe() {
	// 指定iframe的frameId，为子框架创建独立环境
	iframeFrameId := "iframe-123456" // 实际使用时获取真实frameId
	result, resp, err := CDPPageCreateIsolatedWorld(iframeFrameId, "iframe-isolated-world", false)
	if err != nil {
		log.Printf("iframe隔离环境创建失败: %v", err)
		return
	}
	log.Printf("iframe隔离环境创建成功: %s", result)
}

// === 使用场景示例代码：自动化测试安全隔离执行 ===
func ExampleCreateWorldForTest() {
	// 测试前创建独立环境，避免测试脚本污染页面
	_, resp, err := CDPPageCreateIsolatedWorld("", "test-safe-world", false)
	if err != nil {
		log.Fatalf("测试隔离环境创建失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 测试安全隔离环境已创建，可执行测试脚本")
}

*/

// -----------------------------------------------  Page.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 停止页面事件监听: 取消之前通过 Page.enable 开启的所有页面事件通知
// 2. 资源释放: 不再需要页面事件时关闭事件推送，节省浏览器资源
// 3. 流程控制: 自动化任务分段执行，暂停页面事件接收
// 4. 避免干扰: 防止页面事件影响后续其他模块操作
// 5. 清理环境: 任务结束后清理CDP页面事件订阅状态

// CDPPageDisable 禁用Page域事件，停止接收页面相关通知
func CDPPageDisable() (string, error) {
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
		"method": "Page.disable"
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

// === 使用场景示例代码：自动化测试结束后关闭页面事件 ===
func ExampleDisableAfterTest() {
	// 自动化测试逻辑执行完毕
	log.Println("自动化测试完成，不再需要页面事件")

	// 禁用Page事件，停止接收页面通知
	resp, err := CDPPageDisable()
	if err != nil {
		log.Fatalf("禁用Page事件失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ Page事件已禁用，停止接收页面加载/跳转等事件")
}

// === 使用场景示例代码：释放浏览器资源，减少性能消耗 ===
func ExampleDisableForResourceRelease() {
	// 长时间运行任务，中间暂停页面事件
	resp, err := CDPPageDisable()
	if err != nil {
		log.Printf("禁用Page事件失败: %v", err)
		return
	}
	log.Println("页面事件已禁用，浏览器资源消耗降低")
}

// === 使用场景示例代码：避免页面事件干扰后续操作 ===
func ExampleDisableAvoidEventInterrupt() {
	// 即将执行核心逻辑，不希望页面事件干扰
	resp, err := CDPPageDisable()
	if err != nil {
		log.Printf("禁用事件失败: %v", err)
		return
	}

	// 执行核心无干扰操作
	log.Println("已禁用页面事件，开始执行核心操作...")
}

// === 使用场景示例代码：任务结束清理CDP订阅状态 ===
func ExampleDisableOnTaskFinish() {
	// 整体任务完成，清理环境
	resp, err := CDPPageDisable()
	if err != nil {
		log.Printf("任务结束清理失败: %v", err)
		return
	}
	log.Println("🏁 任务完成，Page事件已禁用，CDP环境清理完成")
}

*/

// -----------------------------------------------  Page.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 开启页面事件监听: 启用Page域，接收页面加载、DOMContentLoaded、页面跳转等事件
// 2. 自动化测试前置: 测试开始前必须启用，才能监听页面状态完成测试
// 3. 页面监控: 监控页面加载进度、错误、重定向等生命周期事件
// 4. 调试辅助: 调试页面时实时获取页面加载、刷新、导航信息
// 5. 任务初始化: 自动化/爬虫任务启动时初始化页面事件监听
// 6. 恢复事件监听: 调用Page.disable后重新启用页面事件

// CDPPageEnable 启用Page域，开启页面相关事件监听
func CDPPageEnable() (string, error) {
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
		"method": "Page.enable"
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

// === 使用场景示例代码：自动化测试开始前启用页面事件 ===
func ExampleEnableBeforeTest() {
	// 自动化测试必须先启用Page事件
	resp, err := CDPPageEnable()
	if err != nil {
		log.Fatalf("启用Page事件失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ Page事件已启用，可监听页面加载、跳转等状态")
}

// === 使用场景示例代码：爬虫任务初始化启用页面监听 ===
func ExampleEnableForCrawlerTask() {
	// 爬虫任务启动，开启页面事件监控加载状态
	resp, err := CDPPageEnable()
	if err != nil {
		log.Printf("爬虫初始化失败: %v", err)
		return
	}
	log.Println("🕷️ 爬虫页面事件已启用，准备开始采集")
}

// === 使用场景示例代码：禁用后恢复页面事件监听 ===
func ExampleEnableAfterDisable() {
	// 先禁用
	_, _ = CDPPageDisable()
	log.Println("已临时禁用Page事件")

	// 重新启用
	resp, err := CDPPageEnable()
	if err != nil {
		log.Printf("恢复Page事件失败: %v", err)
		return
	}
	log.Println("🔄 Page事件已恢复监听")
}

// === 使用场景示例代码：页面调试实时监控 ===
func ExampleEnableForDebug() {
	// 调试页面时启用，获取加载、错误、刷新信息
	resp, err := CDPPageEnable()
	if err != nil {
		log.Fatalf("调试模式启用失败: %v", err)
	}
	log.Println("🔍 Page调试监控已启用")
}

*/

// -----------------------------------------------  Page.getAppManifest  -----------------------------------------------
// === 应用场景 ===
// 1.  PWA 应用验证：检查网页是否配置 Web App Manifest，验证 manifest.json 有效性
// 2.  元数据提取：获取 PWA 应用名称、图标、启动页、显示模式等核心配置
// 3.  安装性检测：结合 manifest 与浏览器规则，判断网页是否可被安装为桌面应用
// 4.  错误诊断：解析 manifest 加载、语法、字段合法性错误（如缺失图标、start_url 异常）
// 5.  自动化审计：爬虫/测试工具批量校验 PWA 配置完整性与合规性
// 6.  多 manifest 管理：通过 manifestId 精准获取指定 Web 应用的清单（多应用场景）

// CDPPageGetAppManifest 获取当前页面的 Web App Manifest（PWA 清单）
// 参数:
//
//	manifestId - 可选，指定 manifest ID，不匹配时返回错误；为空则获取当前页 manifest
//
// 返回:
//
//	url    - manifest 文件 URL
//	data   - manifest 原始 JSON 字符串
//	errors - 解析错误列表（无错为空）
//	err    - 调用/网络异常
func CDPPageGetAppManifest(manifestId string) (url, data string, errors []AppManifestError, err error) {
	if !DefaultBrowserWS() {
		return "", "", nil, fmt.Errorf("CDP 未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", "", nil, fmt.Errorf("WebSocket 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建请求（带可选 manifestId）
	var params string
	if manifestId != "" {
		params = fmt.Sprintf(`,"params":{"manifestId":"%s"}`, manifestId)
	}
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getAppManifest"
		%s
	}`, reqID, params)

	// 发送
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", "", nil, fmt.Errorf("发送失败: %w", err)
	}
	log.Printf("[DEBUG] Page.getAppManifest → %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", "", nil, fmt.Errorf("消息队列关闭")
			}
			if respMsg.ID != reqID {
				continue
			}

			// 解析响应
			var resp struct {
				ID     int               `json:"id"`
				Result AppManifestResult `json:"result"`
				Error  *CDPError         `json:"error"`
			}
			if err := json.Unmarshal([]byte(respMsg.Content), &resp); err != nil {
				return "", "", nil, fmt.Errorf("解析失败: %w", err)
			}
			if resp.Error != nil {
				return "", "", nil, fmt.Errorf("CDP 错误: %s", resp.Error.Message)
			}

			log.Printf("[DEBUG] getAppManifest ← url=%s errors=%d",
				resp.Result.URL, len(resp.Result.Errors))
			return resp.Result.URL, resp.Result.Data, resp.Result.Errors, nil

		case <-timer.C:
			return "", "", nil, fmt.Errorf("请求超时(%s)", timeout)
		}
	}
}

// --- 类型定义 ---
// AppManifestResult CDP 返回结构
type AppManifestResult struct {
	URL    string             `json:"url"`
	Data   string             `json:"data"`
	Errors []AppManifestError `json:"errors"`
}

// AppManifestError manifest 解析错误
type AppManifestError struct {
	Message  string `json:"message"`
	Critical int    `json:"critical"` // 1=严重错误(无法解析) 0=警告
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

// CDPError CDP 协议错误
type CDPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/*
// === 使用场景 1：基础获取（当前页 PWA 清单）===
func ExampleGetAppManifest_Basic() {
	url, data, errs, err := CDPPageGetAppManifest("")
	if err != nil {
		log.Fatalf("获取失败: %v", err)
	}
	// 输出错误
	for _, e := range errs {
		log.Printf("Manifest错误[%d:%d] %s (critical=%d)",
			e.Line, e.Column, e.Message, e.Critical)
	}
	log.Printf("清单URL: %s", url)
	log.Printf("原始JSON: %s", data)
}

// === 使用场景 2：指定 manifestId 获取 ===
func ExampleGetAppManifest_ByManifestId() {
	manifestId := "https://example.com/manifest.webmanifest"
	url, data, errs, err := CDPPageGetAppManifest(manifestId)
	if err != nil {
		log.Fatalf("获取失败: %v", err)
	}
	log.Printf("匹配清单: %s", url)
}

// === 使用场景 3：PWA 可安装性审计 ===
func ExampleGetAppManifest_PWAInstallCheck() {
	_, data, errs, err := CDPPageGetAppManifest("")
	if err != nil {
		log.Fatalf("获取失败: %v", err)
	}
	// 1. 检查严重错误
	var criticalErr bool
	for _, e := range errs {
		if e.Critical == 1 {
			criticalErr = true
			log.Printf("严重错误: %s", e.Message)
		}
	}
	if criticalErr {
		log.Println("❌ 清单无效，不可安装")
		return
	}
	// 2. 解析并校验关键字段
	var manifest struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
		StartURL  string `json:"start_url"`
		Display   string `json:"display"`
		Icons     []struct {
			Src string `json:"src"`
			Sizes string `json:"sizes"`
		} `json:"icons"`
	}
	if err := json.Unmarshal([]byte(data), &manifest); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}
	// 3. 合规判断
	valid := manifest.Name != "" && manifest.StartURL != "" &&
		(manifest.Display == "standalone" || manifest.Display == "fullscreen") &&
		len(manifest.Icons) > 0
	if valid {
		log.Printf("✅ 可安装: %s (display=%s)", manifest.Name, manifest.Display)
	} else {
		log.Println("❌ 缺失必填字段（name/start_url/display/icons）")
	}
}

*/

// -----------------------------------------------  Page.getFrameTree  -----------------------------------------------
// === 应用场景 ===
// 1. 页面结构分析: 获取完整的页面框架树，包含主页面+所有iframe嵌套结构
// 2. iframe定位: 获取所有iframe的id、url、父框架信息，用于精准操作子页面
// 3. 多框架调试: 调试嵌套页面时理清框架层级关系
// 4. 自动化操作前置: 操作iframe前必须获取frameId进行切换
// 5. 页面安全检测: 检测页面嵌套的第三方iframe，排查安全风险
// 6. 爬虫数据采集: 遍历框架树获取所有子页面内容

// CDPPageGetFrameTreeAt 获取当前页面的完整框架树（主框架+所有iframe）
func CDPPageGetFrameTreeAt() (string, error) {
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
		"method": "Page.getFrameTree"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getFrameTree 请求失败: %w", err)
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
			return "", fmt.Errorf("getFrameTree 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：获取并打印完整页面框架结构 ===
func ExampleGetFullFrameTree() {
	// 获取框架树信息
	frameTree, err := CDPPageGetFrameTree()
	if err != nil {
		log.Fatalf("获取框架树失败: %v, 响应: %s", err, frameTree)
	}
	log.Println("=== 完整页面框架树 ===")
	log.Println(frameTree)
}

// === 使用场景示例代码：提取所有iframe的id和url ===
func ExampleExtractAllIframes() {
	frameTreeResp, err := CDPPageGetFrameTree()
	if err != nil {
		log.Printf("获取框架树失败: %v", err)
		return
	}

	// 解析框架树数据
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(frameTreeResp), &result); err != nil {
		log.Printf("解析框架树失败: %v", err)
		return
	}

	// 递归遍历所有框架，提取iframe信息
	var extractIframes func(node map[string]interface{})
	extractIframes = func(node map[string]interface{}) {
		if frame, ok := node["frame"].(map[string]interface{}); ok {
			frameId, _ := frame["id"].(string)
			url, _ := frame["url"].(string)
			// 非主框架即为iframe
			if parentId, hasParent := frame["parentId"].(string); hasParent {
				log.Printf("【IFrame】ID: %s, 父框架ID: %s, URL: %s", frameId, parentId, url)
			} else {
				log.Printf("【主框架】ID: %s, URL: %s", frameId, url)
			}
		}

		// 遍历子框架
		if children, ok := node["children"].([]interface{}); ok {
			for _, child := range children {
				if childNode, ok := child.(map[string]interface{}); ok {
					extractIframes(childNode)
				}
			}
		}
	}

	// 开始解析
	if resultMap, ok := result["result"].(map[string]interface{}); ok {
		if tree, ok := resultMap["frameTree"].(map[string]interface{}); ok {
			extractIframes(tree)
		}
	}
}



// === 使用场景示例代码：自动化操作iframe前获取目标frameId ===
func ExampleGetFrameIdForIframeOperation() {
	frameTree, err := CDPPageGetFrameTree()
	if err != nil {
		log.Fatalf("获取框架树失败: %v", err)
	}

	// 模拟：根据url匹配需要操作的iframe
	targetUrl := "https://example.com/iframe-content"
	log.Printf("正在查找包含URL: %s 的iframe...", targetUrl)

	// 解析逻辑同上，找到匹配的frameId后即可进行后续操作
	log.Println("已获取目标frameId: abc123-def456")
	log.Println("✅ 可使用frameId执行注入、点击等操作")
	log.Println(frameTree)
}

// === 使用场景示例代码：页面安全检测-列出所有第三方iframe ===
func ExampleCheckThirdPartyIframes() {
	frameTree, err := CDPPageGetFrameTree()
	if err != nil {
		log.Printf("安全检测-获取框架树失败: %v", err)
		return
	}

	log.Println("=== 安全检测：所有嵌套iframe列表 ===")
	log.Println(frameTree)
	log.Println("⚠️  请检查以上iframe是否为可信第三方源")
}

*/

// -----------------------------------------------  Page.getLayoutMetrics  -----------------------------------------------
// === 应用场景 ===
// 1. 获取页面布局尺寸: 获取可视区域、内容区域、滚动条等真实布局参数
// 2. 响应式适配: 获取屏幕/页面宽高，用于自动化适配不同设备
// 3. 滚动计算: 基于内容尺寸和可视尺寸计算滚动范围、滚动位置
// 4. 截图辅助: 全屏截图时获取完整页面高度，避免截取不全
// 5. 元素定位: 结合布局参数精确定位页面元素位置
// 6. 页面渲染监控: 监控页面布局是否正常渲染，尺寸是否符合预期

// CDPPageGetLayoutMetrics 获取页面布局度量信息（尺寸、滚动、可视区域）
func CDPPageGetLayoutMetrics() (string, error) {
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
		"method": "Page.getLayoutMetrics"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getLayoutMetrics 请求失败: %w", err)
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
			return "", fmt.Errorf("getLayoutMetrics 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：获取页面核心尺寸（可视区域+内容区域）===
func ExampleGetPageSize() {
	resp, err := CDPPageGetLayoutMetrics()
	if err != nil {
		log.Fatalf("获取布局信息失败: %v, 响应: %s", err, resp)
	}

	// 解析关键尺寸
	var result map[string]interface{}
	json.Unmarshal([]byte(resp), &result)
	layout := result["result"].(map[string]interface{})

	// 可视区域（屏幕可见范围）
	visual := layout["visualViewport"].(map[string]interface{})
	// 内容区域（整个网页实际大小）
	content := layout["contentSize"].(map[string]interface{})

	log.Printf("可视区域宽度: %.0fpx", visual["width"].(float64))
	log.Printf("可视区域高度: %.0fpx", visual["height"].(float64))
	log.Printf("页面内容总宽度: %.0fpx", content["width"].(float64))
	log.Printf("页面内容总高度: %.0fpx", content["height"].(float64))
}

// === 使用场景示例代码：计算可滚动范围，用于自动化滚动 ===
func ExampleCalcScrollRange() {
	resp, _ := CDPPageGetLayoutMetrics()
	var result map[string]interface{}
	json.Unmarshal([]byte(resp), &result)
	layout := result["result"].(map[string]interface{})

	visualH := layout["visualViewport"].(map[string]interface{})["height"].(float64)
	contentH := layout["contentSize"].(map[string]interface{})["height"].(float64)
	scrollMax := contentH - visualH

	log.Printf("可视高度: %.0fpx", visualH)
	log.Printf("内容高度: %.0fpx", contentH)
	log.Printf("最大可滚动高度: %.0fpx", scrollMax)
	log.Println("✅ 可用于执行页面滚动、滑动加载操作")
}

// === 使用场景示例代码：全屏截图前获取完整页面高度 ===
func ExampleGetFullPageHeightForScreenshot() {
	resp, err := CDPPageGetLayoutMetrics()
	if err != nil {
		log.Printf("获取页面高度失败: %v", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(resp), &result)
	contentSize := result["result"].(map[string]interface{})["contentSize"].(map[string]interface{})
	fullHeight := contentSize["height"].(float64)

	log.Printf("完整页面高度: %.0fpx", fullHeight)
	log.Println("📸 使用该高度可截取完整长截图，无缺失")
}

// === 使用场景示例代码：监控页面是否正常渲染 ===
func ExampleCheckPageRendered() {
	resp, _ := CDPPageGetLayoutMetrics()
	var result map[string]interface{}
	json.Unmarshal([]byte(resp), &result)
	layout := result["result"].(map[string]interface{})
	content := layout["contentSize"].(map[string]interface{})

	width := content["width"].(float64)
	height := content["height"].(float64)

	// 正常页面宽高不可能为0
	if width > 0 && height > 0 {
		log.Println("✅ 页面布局正常渲染完成")
	} else {
		log.Println("❌ 页面未正常渲染，宽高异常")
	}
}

*/

// -----------------------------------------------  Page.getNavigationHistory  -----------------------------------------------
// === 应用场景 ===
// 1. 导航记录获取: 获取当前标签页完整的前进/后退浏览历史
// 2. 测试验证: 验证自动化操作中页面跳转是否符合预期
// 3. 页面溯源: 追溯当前页面从哪个链接跳转而来
// 4. 前进后退控制: 结合历史记录序号实现精准前进/后退导航
// 5. 去重检测: 检测是否存在重复跳转、循环跳转
// 6. 日志审计: 记录用户/自动化流程的页面访问轨迹

// CDPPageGetNavigationHistory 获取页面导航历史记录
func CDPPageGetNavigationHistory() (string, error) {
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
		"method": "Page.getNavigationHistory"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getNavigationHistory 请求失败: %w", err)
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
			return "", fmt.Errorf("getNavigationHistory 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：获取并打印完整导航历史 ===
func ExampleGetFullNavigationHistory() {
	historyResp, err := CDPPageGetNavigationHistory()
	if err != nil {
		log.Fatalf("获取导航历史失败: %v, 响应: %s", err, historyResp)
	}
	log.Println("=== 完整页面导航历史 ===")
	log.Println(historyResp)
}

// === 使用场景示例代码：解析历史记录，提取跳转轨迹 ===
func ExampleParseNavigationHistory() {
	historyResp, err := CDPPageGetNavigationHistory()
	if err != nil {
		log.Printf("获取导航历史失败: %v", err)
		return
	}

	// 解析结构
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(historyResp), &result); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	res := result["result"].(map[string]interface{})
	currentIdx := res["currentIndex"].(float64)
	entries := res["entries"].([]interface{})

	log.Printf("当前所在索引: %d", int(currentIdx))
	log.Printf("总历史记录数: %d", len(entries))

	// 遍历所有记录
	for i, item := range entries {
		entry := item.(map[string]interface{})
		id := entry["id"].(float64)
		url := entry["url"].(string)
		title := entry["title"].(string)

		// 标记当前页面
		currMark := ""
		if float64(i) == currentIdx {
			currMark = "【当前页】"
		}

		log.Printf("记录%d | ID:%v | %s URL:%s | 标题:%s", i+1, id, currMark, url, title)
	}
}

// === 使用场景示例代码：获取上一页地址（页面溯源） ===
func ExampleGetPreviousPageUrl() {
	historyResp, _ := CDPPageGetNavigationHistory()
	var result map[string]interface{}
	json.Unmarshal([]byte(historyResp), &result)

	res := result["result"].(map[string]interface{})
	currentIdx := res["currentIndex"].(float64)
	entries := res["entries"].([]interface{})

	// 存在上一页
	if currentIdx > 0 {
		prevEntry := entries[int(currentIdx)-1].(map[string]interface{})
		log.Printf("上一页地址: %s", prevEntry["url"].(string))
		log.Printf("上一页标题: %s", prevEntry["title"].(string))
	} else {
		log.Println("当前为第一页，无上一页记录")
	}
}

// === 使用场景示例代码：自动化测试验证跳转是否正确 ===
func ExampleTestNavigationHistory() {
	historyResp, _ := CDPPageGetNavigationHistory()
	var result map[string]interface{}
	json.Unmarshal([]byte(historyResp), &result)

	res := result["result"].(map[string]interface{})
	entries := res["entries"].([]interface{})

	// 预期目标地址
	expectUrl := "https://github.com"

	// 检查是否跳转到了正确页面
	last := entries[len(entries)-1].(map[string]interface{})
	actualUrl := last["url"].(string)

	if actualUrl == expectUrl {
		log.Printf("✅ 跳转验证通过: %s", actualUrl)
	} else {
		log.Printf("❌ 跳转验证失败，预期: %s，实际: %s", expectUrl, actualUrl)
	}
}

*/

// -----------------------------------------------  Page.handleJavaScriptDialog  -----------------------------------------------
// === 应用场景 ===
// 1. 自动处理弹窗: 自动化测试中自动处理alert/confirm/prompt弹窗
// 2. 弹窗拦截: 阻止页面弹窗阻塞流程，自动确认/取消
// 3. 输入弹窗响应: 对prompt输入框自动填写内容并提交
// 4. 异常弹窗处理: 页面意外弹窗时自动关闭，避免流程中断
// 5. 测试验证: 验证弹窗是否正常触发，并模拟用户操作
// 6. 爬虫自动化: 处理网站验证、提示类弹窗，继续采集流程

// CDPPageHandleJavaScriptDialog 处理页面JavaScript弹窗(alert/confirm/prompt)
// 参数 accept: true-确认/关闭弹窗，false-取消弹窗(仅confirm支持)
// 参数 promptText: prompt弹窗输入的文字(仅prompt需要)
func CDPPageHandleJavaScriptDialog(accept bool, promptText string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	var params string
	if promptText != "" {
		// 处理prompt输入弹窗
		params = fmt.Sprintf(`"accept": %t, "promptText": %s`, accept, strconv.Quote(promptText))
	} else {
		// 处理alert/confirm弹窗
		params = fmt.Sprintf(`"accept": %t`, accept)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.handleJavaScriptDialog",
		"params": {
			%s
		}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 handleJavaScriptDialog 请求失败: %w", err)
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
			return "", fmt.Errorf("handleJavaScriptDialog 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：自动关闭alert弹窗 ===
func ExampleHandleAlertDialog() {
	// 自动确认关闭alert提示框
	resp, err := CDPPageHandleJavaScriptDialog(true, "")
	if err != nil {
		log.Fatalf("处理alert弹窗失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ alert弹窗已自动关闭")
}

// === 使用场景示例代码：自动确认confirm弹窗 ===
func ExampleConfirmDialog() {
	// 自动点击确认(OK)
	resp, err := CDPPageHandleJavaScriptDialog(true, "")
	if err != nil {
		log.Printf("确认confirm失败: %v", err)
		return
	}
	log.Println("✅ confirm弹窗已自动确认")
}

// === 使用场景示例代码：自动取消confirm弹窗 ===
func ExampleCancelDialog() {
	// 自动点击取消(Cancel)
	resp, err := CDPPageHandleJavaScriptDialog(false, "")
	if err != nil {
		log.Printf("取消confirm失败: %v", err)
		return
	}
	log.Println("✅ confirm弹窗已自动取消")
}

// === 使用场景示例代码：自动填写prompt输入框并提交 ===
func ExampleHandlePromptDialog() {
	// 自动输入文字并确认
	inputText := "自动化测试输入内容"
	resp, err := CDPPageHandleJavaScriptDialog(true, inputText)
	if err != nil {
		log.Fatalf("处理prompt弹窗失败: %v, 响应: %s", err, resp)
	}
	log.Printf("✅ prompt弹窗已自动输入：%s，并确认提交", inputText)
}

*/

// -----------------------------------------------  Page.navigate  -----------------------------------------------
// === 应用场景 ===
// 1. 页面跳转：控制浏览器跳转到指定 URL
// 2. 自动化测试：打开目标页面开始执行测试流程
// 3. 爬虫初始化：导航到采集页面
// 4. 多任务切换：在同一个标签页切换不同页面
// 5. 调试导航：手动控制页面加载调试问题
// 6. 定时任务：定时打开指定页面

// CDPPageNavigate 跳转到指定 URL
func CDPPageNavigate(url string) (string, error) {
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
		"method": "Page.navigate",
		"params": {
			"url": %s
		}
	}`, reqID, strconv.Quote(url))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 navigate 请求失败: %w", err)
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
			return "", fmt.Errorf("navigate 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：基础页面跳转 ===
func ExampleNavigateBasic() {
	url := "https://www.baidu.com"
	resp, err := CDPPageNavigate(url)
	if err != nil {
		log.Fatalf("页面跳转失败: %v, 响应: %s", err, resp)
	}
	log.Printf("✅ 正在跳转到: %s", url)
}

// === 使用场景示例代码：自动化测试打开测试页面 ===
func ExampleNavigateTestPage() {
	testUrl := "https://test.example.com/login"
	resp, err := CDPPageNavigate(testUrl)
	if err != nil {
		log.Fatalf("打开测试页面失败: %v", err)
	}
	log.Println("✅ 测试页面已加载，开始执行测试")
}

// === 使用场景示例代码：爬虫导航到目标页面 ===
func ExampleNavigateCrawler() {
	target := "https://news.example.com"
	resp, err := CDPPageNavigate(target)
	if err != nil {
		log.Printf("导航失败: %v", err)
		return
	}
	log.Println("✅ 已到达采集页面，准备提取数据")
}

// === 使用场景示例代码：带 referrer 跳转（部分网站需要） ===
func ExampleNavigateWithReferrer() {
	url := "https://example.com"
	referrer := "https://google.com"

	// 扩展用法：如需带 referrer，可扩展参数
	// 这里使用完整参数格式演示
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	msg := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.navigate",
		"params": {
			"url": %s,
			"referrer": %s
		}
	}`, reqID, strconv.Quote(url), strconv.Quote(referrer))

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("带referrer跳转失败: %v", err)
		return
	}
	log.Printf("✅ 带 referrer 跳转到: %s", url)
}

*/

// -----------------------------------------------  Page.navigateToHistoryEntry  -----------------------------------------------
// === 应用场景 ===
// 1. 历史记录跳转: 通过条目ID精准跳转到浏览历史中的任意页面
// 2. 自动化前进后退: 替代浏览器前进/后退按钮，实现精准导航控制
// 3. 测试流程回溯: 测试中快速回到之前的页面状态，复现/验证流程
// 4. 页面溯源返回: 从深层页面直接返回指定层级的历史页面
// 5. 多步骤回退: 一次性回退多步，而非逐次back
// 6. 异常恢复导航: 页面出错时，直接跳转到稳定的历史记录页面

// CDPPageNavigateToHistoryEntry 通过历史条目ID跳转到对应浏览记录
// 参数 entryId: 从 Page.getNavigationHistory 获取的历史条目唯一ID
func CDPPageNavigateToHistoryEntry(entryId int64) (string, error) {
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
		"method": "Page.navigateToHistoryEntry",
		"params": {
			"entryId": %d
		}
	}`, reqID, entryId)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 navigateToHistoryEntry 请求失败: %w", err)
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
			return "", fmt.Errorf("navigateToHistoryEntry 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：获取历史ID并跳转到指定历史页面 ===
func ExampleNavigateToSpecifiedHistory() {
	// 1. 先获取导航历史
	historyResp, err := CDPPageGetNavigationHistory()
	if err != nil {
		log.Fatalf("获取导航历史失败: %v", err)
	}

	// 2. 解析历史条目，获取 entryId
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(historyResp), &result); err != nil {
		log.Fatalf("解析历史失败: %v", err)
	}

	res := result["result"].(map[string]interface{})
	entries := res["entries"].([]interface{})

	if len(entries) == 0 {
		log.Fatalf("无浏览历史")
	}

	// 3. 选择第一个历史条目进行跳转（可按需选择）
	firstEntry := entries[0].(map[string]interface{})
	entryId := int64(firstEntry["id"].(float64))
	url := firstEntry["url"].(string)

	log.Printf("准备跳转到历史记录: ID=%d, URL=%s", entryId, url)

	// 4. 执行历史跳转
	resp, err := CDPPageNavigateToHistoryEntry(entryId)
	if err != nil {
		log.Fatalf("历史跳转失败: %v, 响应: %s", err, resp)
	}

	log.Println("✅ 成功跳转到指定浏览历史")
}

// === 使用场景示例代码：自动化测试 - 回到上一个测试页面 ===
func ExampleTestBackToPreviousPage() {
	historyResp, _ := CDPPageGetNavigationHistory()
	var result map[string]interface{}
	json.Unmarshal([]byte(historyResp), &result)

	res := result["result"].(map[string]interface{})
	currentIdx := int(res["currentIndex"].(float64))
	entries := res["entries"].([]interface{})

	// 回到上一页（当前索引 - 1）
	if currentIdx > 0 {
		prevEntry := entries[currentIdx-1].(map[string]interface{})
		entryId := int64(prevEntry["id"].(float64))

		resp, err := CDPPageNavigateToHistoryEntry(entryId)
		if err != nil {
			log.Printf("返回上一页失败: %v", err)
			return
		}
		log.Println("✅ 测试已回到上一个页面")
	}
}

// === 使用场景示例代码：直接跳转到最初的页面（首页） ===
func ExampleNavigateToFirstHistory() {
	historyResp, _ := CDPPageGetNavigationHistory()
	var result map[string]interface{}
	json.Unmarshal([]byte(historyResp), &result)

	res := result["result"].(map[string]interface{})
	entries := res["entries"].([]interface{})

	// 直接跳转到第一条历史记录（通常为初始页面）
	firstEntry := entries[0].(map[string]interface{})
	entryId := int64(firstEntry["id"].(float64))

	resp, err := CDPPageNavigateToHistoryEntry(entryId)
	if err != nil {
		log.Printf("返回首页失败: %v", err)
		return
	}
	log.Println("✅ 已直接跳转到最初页面")
}

*/

// -----------------------------------------------  Page.printToPDF  -----------------------------------------------
// === 应用场景 ===
// 1. 页面转PDF: 将当前网页完整导出为PDF文件
// 2. 报表生成: 自动化生成页面数据报表、合同、发票PDF
// 3. 文章存档: 网页文章、文档、资料永久保存为PDF
// 4. 测试报告: 自动化测试结果页面导出为PDF报告
// 5. 离线分享: 将页面转为PDF方便离线查看、分享
// 6. 批量导出: 批量抓取页面并生成PDF存档

// CDPPagePrintToPDF 将页面导出为PDF
// 参数 landscape: 横向打印(true) / 纵向(false)
// 参数 printBackground: 是否打印背景颜色和图片
// 参数 paperWidth: 纸张宽度(英寸)
// 参数 paperHeight: 纸张高度(英寸)
// 返回值: base64编码的PDF数据、响应内容、错误信息
func CDPPagePrintToPDF(landscape bool, printBackground bool, paperWidth float64, paperHeight float64) (string, string, error) {
	if !DefaultBrowserWS() {
		return "", "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.printToPDF",
		"params": {
			"landscape": %t,
			"printBackground": %t,
			"paperWidth": %.2f,
			"paperHeight": %.2f
		}
	}`, reqID, landscape, printBackground, paperWidth, paperHeight)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", "", fmt.Errorf("发送 printToPDF 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（PDF生成稍慢，超时15秒）
	timeout := 15 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", "", fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应检查错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return "", content, fmt.Errorf("解析响应失败: %w", err)
				}

				if errorObj, exists := response["error"]; exists {
					return "", content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				// 提取base64数据
				if result, ok := response["result"].(map[string]interface{}); ok {
					if data, ok := result["data"].(string); ok {
						return data, content, nil
					}
				}

				return "", content, nil
			}

		case <-timer.C:
			return "", "", fmt.Errorf("printToPDF 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：标准A4 PDF导出（纵向+带背景）===
func ExamplePrintToPDFA4() {
	// A4纸张：宽8.27英寸，高11.7英寸，纵向，打印背景
	base64Data, resp, err := CDPPagePrintToPDF(false, true, 8.27, 11.7)
	if err != nil {
		log.Fatalf("生成PDF失败: %v, 响应: %s", err, resp)
	}

	// 解码保存文件
	pdfData, _ := base64.StdEncoding.DecodeString(base64Data)
	err = os.WriteFile("page_a4.pdf", pdfData, 0644)
	if err != nil {
		log.Fatalf("保存PDF失败: %v", err)
	}
	log.Println("✅ A4格式PDF已保存：page_a4.pdf")
}

// === 使用场景示例代码：横向PDF导出（适合宽表格/报表）===
func ExamplePrintToPDFLandscape() {
	// 横向A4，打印背景
	base64Data, resp, err := CDPPagePrintToPDF(true, true, 11.7, 8.27)
	if err != nil {
		log.Printf("横向PDF生成失败: %v", err)
		return
	}

	pdfData, _ := base64.StdEncoding.DecodeString(base64Data)
	_ = os.WriteFile("page_landscape.pdf", pdfData, 0644)
	log.Println("✅ 横向PDF已保存")
}

// === 使用场景示例代码：自动化测试报告导出PDF ===
func ExampleTestReportToPDF() {
	// 测试完成后导出报告
	log.Println("开始生成测试报告...")
	base64Data, _, err := CDPPagePrintToPDF(false, true, 8.27, 11.7)
	if err != nil {
		log.Fatalf("报告生成失败: %v", err)
	}

	// 按时间戳命名保存
	filename := fmt.Sprintf("test_report_%d.pdf", time.Now().Unix())
	pdfData, _ := base64.StdEncoding.DecodeString(base64Data)
	_ = os.WriteFile(filename, pdfData, 0644)
	log.Printf("✅ 测试报告已导出：%s", filename)
}

// === 使用场景示例代码：无背景简洁PDF（适合文字文档）===
func ExamplePrintToPDFWithoutBackground() {
	// 纵向A4，不打印背景（节省墨水）
	base64Data, resp, err := CDPPagePrintToPDF(false, false, 8.27, 11.7)
	if err != nil {
		log.Printf("简洁PDF生成失败: %v", err)
		return
	}

	pdfData, _ := base64.StdEncoding.DecodeString(base64Data)
	_ = os.WriteFile("page_clean.pdf", pdfData, 0644)
	log.Println("✅ 无背景简洁PDF已保存")
}

*/

// -----------------------------------------------  Page.reload  -----------------------------------------------
// === 应用场景 ===
// 1. 页面刷新重置: 刷新当前页面恢复初始状态
// 2. 测试环境重置: 自动化测试用例之间刷新页面清空状态
// 3. 数据强制更新: 忽略缓存强制刷新获取最新页面数据
// 4. 异常页面恢复: 页面报错/卡死时刷新修复
// 5. 定时刷新: 监控、报表页面定时刷新获取最新内容
// 6. 配置生效刷新: 修改设置后刷新页面使配置生效

// CDPPageReload 刷新当前页面
// 参数 ignoreCache: 是否忽略浏览器缓存（true=强制刷新，false=普通刷新）
func CDPPageReload(ignoreCache bool) (string, error) {
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
		"method": "Page.reload",
		"params": {
			"ignoreCache": %t
		}
	}`, reqID, ignoreCache)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 reload 请求失败: %w", err)
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
			return "", fmt.Errorf("reload 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：普通刷新页面（使用缓存）===
func ExampleNormalReload() {
	// 普通刷新（等价于 F5）
	resp, err := CDPPageReload(false)
	if err != nil {
		log.Fatalf("页面刷新失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面已普通刷新（F5）")
}

// === 使用场景示例代码：强制刷新忽略缓存（等价Ctrl+F5）===
func ExampleForceReloadIgnoreCache() {
	// 强制刷新，获取最新页面内容
	resp, err := CDPPageReload(true)
	if err != nil {
		log.Printf("强制刷新失败: %v", err)
		return
	}
	log.Println("✅ 页面已强制刷新，忽略缓存获取最新数据")
}

// === 使用场景示例代码：自动化测试用例之间重置页面 ===
func ExampleReloadBetweenTestCases() {
	// 一个测试用例执行完毕，刷新页面准备下一个用例
	log.Println("测试用例1执行完成，重置页面状态...")

	resp, err := CDPPageReload(true)
	if err != nil {
		log.Fatalf("测试页面重置失败: %v", err)
	}
	log.Println("✅ 页面已重置，准备执行下一个测试用例")
}

// === 使用场景示例代码：页面异常时自动修复刷新 ===
func ExampleReloadOnPageError() {
	// 模拟检测到页面异常/卡死
	pageError := true

	if pageError {
		log.Println("⚠️ 检测到页面异常，尝试刷新修复...")
		resp, err := CDPPageReload(true)
		if err != nil {
			log.Printf("修复刷新失败: %v", err)
			return
		}
		log.Println("✅ 页面已刷新修复，恢复正常")
	}
}

// === 使用场景示例代码：定时自动刷新页面 ===
func ExampleTimedReload() {
	// 每30秒强制刷新一次
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("执行定时刷新...")
		_, err := CDPPageReload(true)
		if err != nil {
			log.Printf("定时刷新失败: %v", err)
		}
	}
}

*/

// -----------------------------------------------  Page.removeScriptToEvaluateOnNewDocument  -----------------------------------------------
// === 应用场景 ===
// 1. 清理注入脚本: 移除之前通过 Page.addScriptToEvaluateOnNewDocument 注入的脚本
// 2. 环境重置: 自动化测试后清理全局注入脚本，恢复页面原始环境
// 3. 动态脚本管理: 按需启用/禁用前置注入脚本，灵活切换运行环境
// 4. 反爬策略关闭: 关闭之前注入的反爬绕过脚本，避免长期驻留
// 5. 多任务隔离: 不同任务之间清除上一个任务的注入脚本，防止冲突

// CDPPageRemoveScriptToEvaluateOnNewDocument 移除新文档加载时执行的注入脚本
// 参数 identifier: 添加脚本时返回的唯一标识ID（scriptID）
func CDPPageRemoveScriptToEvaluateOnNewDocument(identifier string) (string, error) {
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
		"method": "Page.removeScriptToEvaluateOnNewDocument",
		"params": {
			"identifier": "%s"
		}
	}`, reqID, identifier)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 removeScriptToEvaluateOnNewDocument 请求失败: %w", err)
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
			return "", fmt.Errorf("removeScriptToEvaluateOnNewDocument 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：添加脚本后立即移除清理 ===
func ExampleAddAndRemoveScript() {
	// 1. 注入脚本
	script := `window.test = "hello"`
	scriptID, resp, err := CDPPageAddScriptToEvaluateOnNewDocument(script)
	if err != nil {
		log.Fatalf("注入脚本失败: %v", err)
	}
	log.Printf("脚本注入成功，ID: %s", scriptID)

	// 2. 移除脚本（清理环境）
	resp, err = CDPPageRemoveScriptToEvaluateOnNewDocument(scriptID)
	if err != nil {
		log.Fatalf("移除脚本失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 注入脚本已成功移除")
}

// === 使用场景示例代码：自动化测试后清理注入脚本 ===
func ExampleCleanupAfterTest() {
	// 测试用的注入脚本ID（来自之前添加时保存的ID）
	testScriptID := "test-script-123"

	// 测试完成后清理，避免影响其他用例
	resp, err := CDPPageRemoveScriptToEvaluateOnNewDocument(testScriptID)
	if err != nil {
		log.Printf("测试脚本清理失败: %v", err)
		return
	}
	log.Println("✅ 测试注入脚本已清理完成")
}

// === 使用场景示例代码：动态关闭反爬注入环境 ===
func ExampleDisableCrawlerScript() {
	// 之前注入的绕过webdriver检测脚本ID
	crawlerScriptID := "webdriver-bypass-script"

	// 关闭爬虫模式，移除注入脚本
	resp, err := CDPPageRemoveScriptToEvaluateOnNewDocument(crawlerScriptID)
	if err != nil {
		log.Printf("关闭爬虫环境失败: %v", err)
		return
	}
	log.Println("✅ 爬虫注入脚本已移除，恢复原始环境")
}

*/

// -----------------------------------------------  Page.resetNavigationHistory  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境隔离: 自动化测试用例间清空历史记录，防止状态污染、避免后退/前进按钮异常
// 2. 单页应用重置: SPA路由跳转过多后重置历史，防止后退逻辑混乱
// 3. 隐私保护: 敏感操作后清除当前页面会话历史，防止信息泄露
// 4. 界面状态重置: 清除历史后禁用浏览器后退/前进，强制停留在当前页
// 5. 会话初始化: 新用户会话开始时清空历史，保证初始状态干净

// CDPPageResetNavigationHistory 重置当前页面的导航历史栈
// 效果：清空后退/前进记录，当前页成为历史唯一条目
func CDPPageResetNavigationHistory() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建命令（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.resetNavigationHistory"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 resetNavigationHistory 请求失败: %w", err)
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

				// 解析错误
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errorObj, has := response["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}
				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("resetNavigationHistory 请求超时")
		}
	}
}

/*

// === 使用场景1：自动化测试用例间清空历史 ===
func ExampleResetForTestIsolation() {
	log.Println("执行测试用例前，清空页面历史...")
	resp, err := CDPPageResetNavigationHistory()
	if err != nil {
		log.Fatalf("重置历史失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 导航历史已重置，当前页面为历史唯一条目")
}

// === 使用场景2：敏感操作后清除历史（防后退/泄露） ===
func ExampleResetAfterSensitiveOperation() {
	// 模拟支付/登录等敏感流程完成
	log.Println("敏感操作完成，清除会话历史...")
	resp, err := CDPPageResetNavigationHistory()
	if err != nil {
		log.Printf("清除历史失败: %v", err)
		return
	}
	log.Println("✅ 历史已清空，无法后退到上一页面")
}

// === 使用场景3：SPA路由过多后重置历史栈 ===
func ExampleResetSPAHistory() {
	// 单页应用多次路由后历史混乱
	log.Println("SPA路由过多，重置历史栈...")
	resp, err := CDPPageResetNavigationHistory()
	if err != nil {
		log.Printf("重置失败: %v", err)
		return
	}
	log.Println("✅ 历史已重置，后退/前进按钮恢复正常")
}

// === 使用场景4：配合getNavigationHistory验证 ===
func ExampleResetAndVerifyHistory() {
	// 先重置
	resp, err := CDPPageResetNavigationHistory()
	if err != nil {
		log.Fatalf("重置失败: %v", err)
	}
	// 再获取历史验证
	history, err := CDPPageGetNavigationHistory()
	if err != nil {
		log.Fatalf("获取历史失败: %v", err)
	}
	log.Printf("✅ 重置后历史条目数: %d, 当前索引: %d",
		len(history.Entries), history.CurrentIndex)
}

*/

// -----------------------------------------------  Page.setBypassCSP  -----------------------------------------------
// === 应用场景 ===
// 1. 注入脚本不受限：绕过页面CSP限制，成功注入自定义JS脚本
// 2. 自动化测试兼容：解决测试脚本因CSP无法运行、无法加载资源的问题
// 3. 爬虫数据获取：绕过CSP限制，正常执行页面逻辑、提取数据
// 4. 调试页面：调试时允许执行内联脚本、加载被阻止的资源
// 5. 跨域资源访问：解决因CSP导致的跨域请求、资源加载失败
// 6. 自定义工具注入：在严格CSP页面注入调试/监控工具

// CDPPageSetBypassCSP 设置是否绕过页面CSP(内容安全策略)
// 参数 bypass: true=启用绕过CSP，false=关闭绕过，使用页面原始CSP
func CDPPageSetBypassCSP(bypass bool) (string, error) {
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
		"method": "Page.setBypassCSP",
		"params": {
			"enabled": %t
		}
	}`, reqID, bypass)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setBypassCSP 请求失败: %w", err)
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
			return "", fmt.Errorf("setBypassCSP 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：页面加载前启用CSP绕过（最常用）===
func ExampleEnableBypassCSPBeforeLoad() {
	// 必须在页面导航/刷新前执行才生效
	resp, err := CDPPageSetBypassCSP(true)
	if err != nil {
		log.Fatalf("启用CSP绕过失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 已启用CSP绕过，后续页面可自由注入脚本")

	// 启用后再跳转页面，配置生效
	_, _ = CDPPageNavigate("https://example.com")
}

// === 使用场景示例代码：自动化测试解决脚本注入失败 ===
func ExampleTestBypassCSP() {
	// 测试前置：开启CSP绕过
	resp, err := CDPPageSetBypassCSP(true)
	if err != nil {
		log.Fatalf("测试环境CSP配置失败: %v", err)
	}
	log.Println("✅ 测试环境：已绕过页面CSP限制，可正常执行测试脚本")
}

// === 使用场景示例代码：爬虫绕过CSP限制采集数据 ===
func ExampleCrawlerBypassCSP() {
	// 目标网站CSP严格，无法执行采集脚本
	resp, err := CDPPageSetBypassCSP(true)
	if err != nil {
		log.Printf("爬虫CSP绕过失败: %v", err)
		return
	}
	log.Println("✅ 爬虫：已绕过CSP，可正常执行数据采集逻辑")
}

// === 使用场景示例代码：关闭CSP绕过，恢复页面安全策略 ===
func ExampleDisableBypassCSP() {
	// 任务完成后恢复原始CSP策略
	resp, err := CDPPageSetBypassCSP(false)
	if err != nil {
		log.Printf("关闭CSP绕过失败: %v", err)
		return
	}
	log.Println("✅ 已关闭CSP绕过，恢复页面原始安全策略")
}

*/

// -----------------------------------------------  Page.setDocumentContent  -----------------------------------------------
// === 应用场景 ===
// 1. 动态渲染页面：直接设置页面HTML内容，快速渲染自定义页面
// 2. 模板渲染：将后端模板直接渲染到浏览器，无需加载URL
// 3. 数据预览：将采集/生成的HTML内容实时预览
// 4. 测试用例构造：快速构造测试所需的页面结构
// 5. 离线页面生成：直接生成完整HTML页面用于截图/PDF
// 6. 修复页面：替换错误/异常页面内容为正常内容

// CDPPageSetDocumentContent 直接设置当前页面的HTML内容
// 参数 frameId：目标框架ID（空字符串使用主框架）
// 参数 html：要设置的完整HTML字符串
func CDPPageSetDocumentContent(frameId string, html string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	var params string
	if frameId == "" {
		params = fmt.Sprintf(`"html": %s`, strconv.Quote(html))
	} else {
		params = fmt.Sprintf(`"frameId": "%s", "html": %s`, frameId, strconv.Quote(html))
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setDocumentContent",
		"params": {
			%s
		}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setDocumentContent 请求失败: %w", err)
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
			return "", fmt.Errorf("setDocumentContent 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：直接渲染自定义HTML页面 ===
func ExampleSetCustomHtmlContent() {
	// 自定义完整HTML内容
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>自定义页面</title>
</head>
<body>
    <h1>CDP 直接渲染页面</h1>
    <p>这是通过 Page.setDocumentContent 设置的内容</p>
</body>
</html>`

	// 设置到主框架
	resp, err := CDPPageSetDocumentContent("", html)
	if err != nil {
		log.Fatalf("设置页面内容失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面HTML内容已成功替换渲染")
}

// === 使用场景示例代码：构造测试页面 ===
func ExampleSetTestPageContent() {
	// 构造测试用的表单页面
	testHtml := `<form id="testForm">
		<input name="username" value="testuser">
		<button type="submit">提交</button>
	</form>`

	resp, err := CDPPageSetDocumentContent("", testHtml)
	if err != nil {
		log.Fatalf("构造测试页面失败: %v", err)
	}
	log.Println("✅ 测试页面已构造完成，可执行自动化测试")
}

// === 使用场景示例代码：预览采集的HTML内容 ===
func ExamplePreviewCrawledHtml() {
	// 模拟已采集的HTML
	crawledHtml := `<div>采集的内容标题<p>采集的详细内容</p></div>`

	// 直接预览
	resp, err := CDPPageSetDocumentContent("", crawledHtml)
	if err != nil {
		log.Printf("预览HTML失败: %v", err)
		return
	}
	log.Println("✅ 采集的HTML内容已实时预览")
}

// === 使用场景示例代码：修改iframe内容 ===
func ExampleSetIframeContent() {
	// 先获取frameId（实际使用时从getFrameTree获取）
	iframeFrameId := "iframe-123"

	// 设置iframe的HTML
	iframeHtml := `<h3>自定义IFrame内容</h3>`
	resp, err := CDPPageSetDocumentContent(iframeFrameId, iframeHtml)
	if err != nil {
		log.Printf("修改iframe内容失败: %v", err)
		return
	}
	log.Println("✅ IFrame内容已成功替换")
}

*/

// -----------------------------------------------  Page.setInterceptFileChooserDialog  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化文件上传：拦截文件选择对话框，自动上传文件，无需人工点击
// 2. 测试文件上传功能：自动化测试中模拟文件选择，验证上传流程
// 3. 无头模式文件操作：无头浏览器（headless）无法弹出对话框，必须拦截处理
// 4. 批量文件上传：批量任务中自动选择文件，避免手动操作
// 5. 避免流程阻塞：防止文件选择框阻塞自动化爬虫/测试流程

// CDPPageSetInterceptFileChooserDialog 设置是否拦截文件选择对话框
// 参数 enabled：true=启用拦截（对话框不弹出，由CDP处理），false=关闭拦截（正常弹出）
func CDPPageSetInterceptFileChooserDialog(enabled bool) (string, error) {
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
		"method": "Page.setInterceptFileChooserDialog",
		"params": {
			"enabled": %t
		}
	}`, reqID, enabled)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setInterceptFileChooserDialog 请求失败: %w", err)
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
			return "", fmt.Errorf("setInterceptFileChooserDialog 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：启用文件对话框拦截，准备自动化上传 ===
func ExampleEnableInterceptFileUpload() {
	// 启用拦截：触发文件选择时不弹出窗口
	resp, err := CDPPageSetInterceptFileChooserDialog(true)
	if err != nil {
		log.Fatalf("启用文件选择框拦截失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 文件选择对话框已拦截，可自动处理文件上传")
}

// === 使用场景示例代码：关闭拦截，恢复手动选择文件 ===
func ExampleDisableInterceptFileChooser() {
	// 关闭拦截，恢复浏览器原生行为
	resp, err := CDPPageSetInterceptFileChooserDialog(false)
	if err != nil {
		log.Printf("关闭文件拦截失败: %v", err)
		return
	}
	log.Println("✅ 文件选择框已恢复正常弹出")
}

// === 使用场景示例代码：无头浏览器自动化文件上传前置 ===
func ExampleHeadlessBrowserFileUpload() {
	// 无头模式必须拦截，否则无法处理上传
	resp, err := CDPPageSetInterceptFileChooserDialog(true)
	if err != nil {
		log.Fatalf("无头模式上传配置失败: %v", err)
	}
	log.Println("✅ 无头模式：文件上传拦截已启用")
}

// === 使用场景示例代码：自动化测试前置配置 ===
func ExampleTestFileUploadSetup() {
	// 测试文件上传功能前开启拦截
	resp, err := CDPPageSetInterceptFileChooserDialog(true)
	if err != nil {
		log.Fatalf("测试环境文件上传配置失败: %v", err)
	}
	log.Println("✅ 自动化测试：文件上传拦截已开启，可执行上传测试")
}

*/

// -----------------------------------------------  Page.setLifecycleEventsEnabled  -----------------------------------------------
// === 应用场景 ===
// 1. 页面生命周期监控：监听加载、渲染、空闲等生命周期事件
// 2. 自动化等待时机：等待页面完全加载/空闲后再执行操作
// 3. 性能采集：监控加载阶段，用于性能分析
// 4. 测试稳定性：确保页面就绪后再点击/输入，避免失败
// 5. 爬虫稳爬：等待页面完全渲染再提取数据
// 6. 调试事件流：观察页面从加载到销毁的完整生命周期

// CDPPageSetLifecycleEventsEnabled 开启/关闭页面生命周期事件推送
// 参数 enabled：true=开启推送，false=关闭推送
func CDPPageSetLifecycleEventsEnabled(enabled bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setLifecycleEventsEnabled",
		"params": {
			"enabled": %t
		}
	}`, reqID, enabled)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送生命周期事件开关失败: %w", err)
	}

	log.Printf("[DEBUG] 发送CDP消息: %s", message)

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到响应: %s", content)

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
			return "", fmt.Errorf("setLifecycleEventsEnabled 请求超时")
		}
	}
}

/*
// === 使用场景示例代码：开启生命周期监控，等待页面完全就绪 ===
func ExampleEnableLifecycleAndWaitForReady() {
	// 开启事件推送
	resp, err := CDPPageSetLifecycleEventsEnabled(true)
	if err != nil {
		log.Fatalf("开启生命周期事件失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面生命周期事件已开启，可监听加载/渲染/空闲")

	// 常见可监听事件：
	// - init
	// - DOMContentLoaded
	// - load
	// - firstPaint
	// - firstContentfulPaint
	// - firstMeaningfulPaint
	// - networkAlmostIdle
	// - networkIdle
	// - lifecycleState: idle
}

// === 使用场景示例代码：自动化操作前等待页面空闲 ===
func ExampleWaitForNetworkIdleBeforeAction() {
	resp, err := CDPPageSetLifecycleEventsEnabled(true)
	if err != nil {
		log.Fatalf("开启失败: %v", err)
	}
	log.Println("✅ 已开启生命周期监听，等待 networkIdle 后执行点击/输入")
}

// === 使用场景示例代码：爬虫等待页面完全渲染 ===
func ExampleCrawlerWaitForPaint() {
	resp, _ := CDPPageSetLifecycleEventsEnabled(true)
	log.Println("🕷️ 爬虫已开启生命周期监听，等待内容绘制完成再提取数据")
}

// === 使用场景示例代码：关闭生命周期事件，减少消息流量 ===
func ExampleDisableLifecycleEvents() {
	resp, err := CDPPageSetLifecycleEventsEnabled(false)
	if err != nil {
		log.Printf("关闭事件推送失败: %v", err)
		return
	}
	log.Println("✅ 页面生命周期事件已关闭，节省性能")
}

*/

// -----------------------------------------------  Page.stopLoading  -----------------------------------------------
// === 应用场景 ===
// 1. 停止页面加载：强制终止页面正在进行的网络请求、资源加载
// 2. 超时中断：页面加载过慢时主动停止，避免阻塞流程
// 3. 非必要资源拦截：无需加载图片/广告/视频时立即停止
// 4. 爬虫提速：只需要DOM结构，加载完成立即停止
// 5. 卡死页面恢复：页面无限加载、卡死时强制终止
// 6. 流量节省：停止不必要的静态资源加载，减少带宽消耗

// CDPPageStopLoading 立即停止当前页面的所有加载行为
func CDPPageStopLoading() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建命令（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.stopLoading"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stopLoading 请求失败: %w", err)
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
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到响应: %s", content)

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
			return "", fmt.Errorf("stopLoading 请求超时")
		}
	}
}

/*
// === 使用场景示例代码：页面加载超时强制停止 ===
func ExampleStopLoadingOnTimeout() {
	// 模拟等待 3 秒后仍未加载完成，主动停止
	go func() {
		time.Sleep(3 * time.Second)
		log.Println("⚠️ 页面加载超时，强制停止...")
		resp, err := CDPPageStopLoading()
		if err != nil {
			log.Fatalf("停止加载失败: %v, 响应: %s", err, resp)
		}
		log.Println("✅ 页面已停止加载，流程继续")
	}()

	// 开始导航
	_, _ = CDPPageNavigate("https://slow-loading-site.com")
}

// === 使用场景示例代码：爬虫DOM就绪后立即停止加载 ===
func ExampleCrawlerStopAfterDOMReady() {
	log.Println("🕷️ 爬虫：DOM已获取，停止多余资源加载...")
	resp, err := CDPPageStopLoading()
	if err != nil {
		log.Printf("停止加载失败: %v", err)
		return
	}
	log.Println("✅ 已停止加载，开始提取数据")
}

// === 使用场景示例代码：页面卡死/无限加载时恢复流程 ===
func ExampleStopLoadingOnPageStuck() {
	// 检测到页面卡死
	isStuck := true
	if isStuck {
		log.Println("⚠️ 页面无限加载，执行强制终止...")
		resp, err := CDPPageStopLoading()
		if err != nil {
			log.Fatalf("修复失败: %v", err)
		}
		log.Println("✅ 页面已恢复，可继续操作")
	}
}

// === 使用场景示例代码：节省流量，不加载图片/视频 ===
func ExampleStopLoadingForTrafficSave() {
	// 页面开始加载后立即停止非必要资源
	time.Sleep(500 * time.Millisecond)
	resp, _ := CDPPageStopLoading()
	log.Println("✅ 已停止资源加载，节省流量")
}


*/

// -----------------------------------------------  Page.addCompilationCache  -----------------------------------------------
// === 应用场景 ===
// 1. 页面秒开：重复访问时 JS 瞬间执行，无解析延迟
// 2. 自动化提速：爬虫/自动化频繁打开页面，显著减少加载时间
// 3. 性能测试：模拟缓存环境，测极致性能表现
// 4. 混合应用：Electron/CEF 内嵌页复用缓存，提升流畅度

// CDPPageAddCompilationCache 注入 JS 编译缓存
func CDPPageAddCompilationCache(url string, cacheData []byte) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP 未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("WebSocket 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// data 必须是 Base64 字符串
	dataBase64 := base64.StdEncoding.EncodeToString(cacheData)

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.addCompilationCache",
		"params": {
			"url": %q,
			"data": %q
		}
	}`, reqID, url, dataBase64)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 addCompilationCache 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → addCompilationCache: url=%s", url)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP 错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("addCompilationCache 超时")
		}
	}
}

/*
// === 配套：Page.produceCompilationCache ===
// 主动生成指定脚本的编译缓存（用于获取 data）
func CDPPageProduceCompilationCache(scripts []CompilationCacheParams) (string, error) {
	// 结构同官方：scripts 数组，每个含 url
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	scriptJson, _ := json.Marshal(scripts)
	msg := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.produceCompilationCache",
		"params": {
			"scripts": %s
		}
	}`, reqID, scriptJson)
	// 发送与等待逻辑同上（省略）
	return "", nil
}

// CompilationCacheParams 生成缓存用的脚本参数
type CompilationCacheParams struct {
	URL string `json:"url"`
}

// === 用法示例：完整缓存流程 ===
func ExampleCompilationCacheFlow() {
	// 1. 先访问页面，生成缓存
	_, _ = CDPPageNavigate("https://example.com")
	// 2. 告诉 Chrome 为指定 JS 生成缓存
	_, _ = CDPPageProduceCompilationCache([]CompilationCacheParams{
		{URL: "https://example.com/app.js"},
	})
	// 3. 监听 Page.compilationCacheProduced 事件获取 data
	// 4. 下次导航前注入缓存
	_ = CDPPageAddCompilationCache(
		"https://example.com/app.js",
		yourCachedBinaryData, // 从事件中拿到的原始字节
	)
	// 5. 再导航：app.js 直接用缓存，秒开
	_, _ = CDPPageNavigate("https://example.com")
}

*/

// -----------------------------------------------  Page.captureSnapshot  -----------------------------------------------
// === 应用场景 ===
// 1. 页面完整存档：导出MHTML，包含所有资源、iframe、样式
// 2. 爬虫数据固化：保存页面原始完整状态，防止内容篡改
// 3. 自动化测试：失败时保存完整页面环境
// 4. 离线网页生成：生成可离线打开的单文件网页

// CDPPageCaptureSnapshot 捕获页面MHTML快照
func CDPPageCaptureSnapshot(format string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数（默认mhtml）
	params := "{}"
	if format != "" {
		params = fmt.Sprintf(`{"format": %q}`, format)
	}

	// 构建CDP命令
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.captureSnapshot",
		"params": %s
	}`, reqID, params)

	// 发送消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送captureSnapshot失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.captureSnapshot: %s", message)

	// 等待响应（10秒超时，MHTML较大可能稍慢）
	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					Data  string      `json:"data"`
					Error interface{} `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return respMsg.Content, fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return respMsg.Content, fmt.Errorf("CDP错误: %v", res.Error)
				}
				// 返回原始MHTML字符串
				return res.Data, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("captureSnapshot请求超时")
		}
	}
}

/*

// === 使用示例：保存页面为MHTML文件 ===
func ExampleCaptureSnapshotAndSave() {
	// 1. 导航页面
	_, _ = CDPPageNavigate("https://www.example.com")

	// 2. 捕获快照（默认mhtml）
	mhtmlData, err := CDPPageCaptureSnapshot("")
	if err != nil {
		log.Fatalf("捕获快照失败: %v", err)
	}

	// 3. 写入文件
	err = os.WriteFile("page_archive.mhtml", []byte(mhtmlData), 0644)
	if err != nil {
		log.Fatalf("保存文件失败: %v", err)
	}
	log.Println("✅ 页面已保存为MHTML：page_archive.mhtml")
}

// === 使用示例：爬虫取证保存 ===
func ExampleCrawlerSnapshotForensics() {
	// 爬取目标页后立即保存完整快照
	targetURL := "https://news.example.com/article"
	_, _ = CDPPageNavigate(targetURL)

	// 捕获快照
	mhtml, _ := CDPPageCaptureSnapshot("mhtml")
	// 保存到取证目录
	_ = os.WriteFile("./forensics/20260408_article.mhtml", []byte(mhtml), 0644)
	log.Println("🕵️ 页面已取证保存")
}

*/

// -----------------------------------------------  Page.clearCompilationCache  -----------------------------------------------
// === 应用场景 ===
// 1. 清理JS编译缓存：强制清除页面已缓存的脚本编译数据
// 2. 调试环境重置：确保JS重新解析编译，避免缓存影响调试结果
// 3. 测试环境隔离：用例之间清空缓存，保证每次运行都是全新环境
// 4. 修复异常脚本：缓存损坏导致执行异常时，清空后重新加载
// 5. 开发热重载：修改JS后清空缓存，强制加载最新版本

// CDPPageClearCompilationCache 清空当前页面的 JavaScript 编译缓存
func CDPPageClearCompilationCache() (string, error) {
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
		"method": "Page.clearCompilationCache"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearCompilationCache 请求失败: %w", err)
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
			return "", fmt.Errorf("clearCompilationCache 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：测试前清空编译缓存 ===
func ExampleClearCacheBeforeTest() {
	// 执行测试前先清空缓存，保证环境干净
	resp, err := CDPPageClearCompilationCache()
	if err != nil {
		log.Fatalf("清空编译缓存失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ JS编译缓存已清空，页面脚本将重新解析编译")
}

// === 使用场景示例代码：脚本异常修复 ===
func ExampleFixScriptByClearCache() {
	// 脚本执行异常，怀疑缓存损坏
	log.Println("检测到脚本执行异常，尝试清空编译缓存...")
	resp, err := CDPPageClearCompilationCache()
	if err != nil {
		log.Printf("清理失败: %v", err)
		return
	}
	// 刷新页面重新加载
	_, _ = CDPPageReload(true)
	log.Println("✅ 缓存已清空，页面已刷新修复")
}

// === 使用场景示例代码：配合addCompilationCache重置缓存环境 ===
func ExampleResetCompilationCache() {
	// 先清空旧缓存
	_, _ = CDPPageClearCompilationCache()
	// 再注入新的编译缓存
	// CDPPageAddCompilationCache(...)
	log.Println("✅ 编译缓存环境已重置")
}

*/

// -----------------------------------------------  Page.crash  -----------------------------------------------
// === 应用场景 ===
// 1. 崩溃测试：模拟页面崩溃，验证监控系统、异常上报、重启机制是否正常工作
// 2. 容错测试：验证服务在页面崩溃后能否自动恢复、重试、兜底
// 3. 调试崩溃恢复：测试浏览器/客户端崩溃重启流程
// 4. 稳定性测试：压力测试中故意触发崩溃，观察系统表现
// 5. 异常演练：模拟线上极端场景，提高系统健壮性

// CDPPageCrash 强制使当前页面渲染进程崩溃（测试专用，不可逆）
func CDPPageCrash() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建命令（无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.crash"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 crash 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 注意：页面崩溃后通常不会返回响应，因此设置较短超时
	timeout := 2 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列关闭（页面已崩溃）")
			}
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
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
			// 无响应通常意味着崩溃成功
			return "", fmt.Errorf("请求超时（页面已成功崩溃，无响应）")
		}
	}
}

/*

// === 使用场景示例代码：模拟页面崩溃，测试异常监控 ===
func ExampleCrashPageForTesting() {
	log.Println("⚠️  开始测试：模拟页面崩溃...")
	resp, err := CDPPageCrash()
	if err != nil {
		log.Printf("✅ 崩溃模拟成功：%v", err)
	} else {
		log.Printf("崩溃响应：%s", resp)
	}
}

// === 使用场景示例代码：崩溃后自动重启恢复 ===
func ExampleCrashAndRecover() {
	// 模拟崩溃
	_, _ = CDPPageCrash()
	log.Println("✅ 页面已崩溃，开始执行自动恢复流程...")

	// 后续可执行：重启标签页、重新连接、重新加载页面
	// browser.Reconnect()
	// CDPPageNavigate("https://example.com")
}

*/

// -----------------------------------------------  Page.generateTestReport -----------------------------------------------
// === 应用场景 ===
// 1. 生成浏览器内置测试报告：用于浏览器兼容性、功能验证、诊断报告
// 2. 自动化测试报告输出：生成标准化测试结果
// 3. 调试页面行为：触发浏览器内部测试并生成诊断数据
// 4. 质量检测：验证页面是否符合浏览器标准行为
// 5. 测试工具集成：对接浏览器测试套件生成报告

// CDPPageGenerateTestReport 触发浏览器生成测试报告
// 参数 message：自定义测试报告消息
func CDPPageGenerateTestReport(message string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建命令
	messageBody := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.generateTestReport",
		"params": {
			"message": %s
		}
	}`, reqID, strconv.Quote(message))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(messageBody))
	if err != nil {
		return "", fmt.Errorf("发送 generateTestReport 请求失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.generateTestReport: %s", messageBody)

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
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("generateTestReport 请求超时")
		}
	}
}

/*
// === 使用场景示例代码：生成测试报告 ===
func ExampleGenerateTestReport() {
	resp, err := CDPPageGenerateTestReport("自动化测试流程完成")
	if err != nil {
		log.Fatalf("生成测试报告失败: %v", err)
	}
	log.Println("✅ 测试报告已生成：", resp)
}

*/

// -----------------------------------------------  Page.getAdScriptAncestry  -----------------------------------------------
// === 应用场景 ===
// 1. 广告追踪溯源：定位广告脚本的完整调用链，明确哪个根脚本导致框架被标记为广告
// 2. 广告拦截调试：分析过滤规则匹配的根广告脚本，优化广告屏蔽策略
// 3. 页面合规审计：检测页面中广告脚本的传播路径，验证广告合规性
// 4. 性能优化：识别嵌套广告脚本，优化页面加载与执行效率
// 5. 安全分析：追踪恶意广告脚本的注入来源与调用层级

// CDPPageGetAdScriptAncestry 获取被标记为广告的框架的脚本祖先链
// 参数 frameId：目标框架的ID
// 返回 AdScriptAncestry：包含广告脚本祖先链与根过滤规则的结构体
func CDPPageGetAdScriptAncestry(frameId string) (*AdScriptAncestry, string, error) {
	if !DefaultBrowserWS() {
		return nil, "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return nil, "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建命令
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getAdScriptAncestry",
		"params": {
			"frameId": %s
		}
	}`, reqID, strconv.Quote(frameId))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return nil, "", fmt.Errorf("发送getAdScriptAncestry请求失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.getAdScriptAncestry: %s", message)

	// 等待响应
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res struct {
					Result *AdScriptAncestry `json:"result"`
					Error  *struct {
						Message string `json:"message"`
					} `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return nil, content, fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return nil, content, fmt.Errorf("CDP错误: %s", res.Error.Message)
				}
				return res.Result, content, nil
			}
		case <-timer.C:
			return nil, "", fmt.Errorf("getAdScriptAncestry请求超时")
		}
	}
}

// AdScriptId 广告脚本标识
type AdScriptId struct {
	DebuggerId string `json:"debuggerId"` // 调试器唯一ID
	ScriptId   string `json:"scriptId"`   // 脚本ID
}

// AdScriptAncestry 广告脚本祖先链结构
type AdScriptAncestry struct {
	AncestryChain            []AdScriptId `json:"ancestryChain"`            // 脚本祖先链（从当前脚本到根脚本）
	RootScriptFilterlistRule string       `json:"rootScriptFilterlistRule"` // 匹配根脚本的过滤规则
}

/*

// === 使用场景示例：溯源广告脚本 ===
func ExampleTraceAdScriptSource() {
	// 假设已知广告框架ID
	adFrameId := "ADE364F9B2C3D45E"

	ancestry, resp, err := CDPPageGetAdScriptAncestry(adFrameId)
	if err != nil {
		log.Fatalf("获取广告脚本祖先链失败: %v, 响应: %s", err, resp)
	}

	log.Println("✅ 广告脚本祖先链获取成功")
	log.Printf("🔗 祖先链层级: %d 层", len(ancestry.AncestryChain))
	for i, script := range ancestry.AncestryChain {
		log.Printf("  %d. 脚本ID: %s (调试器: %s)",
			i+1, script.ScriptId, script.DebuggerId)
	}
	if ancestry.RootScriptFilterlistRule != "" {
		log.Printf("🎯 根过滤规则: %s", ancestry.RootScriptFilterlistRule)
	}
}

*/

// -----------------------------------------------  Page.getAnnotatedPageContent -----------------------------------------------
// === 应用场景 ===
// 1. 获取带浏览器标注的页面DOM内容（无障碍、语义、节点标注）
// 2. 自动化测试：获取带语义标注的页面结构做校验
// 3. 爬虫增强：获取浏览器原生解析后的带注释页面内容
// 4. 页面可访问性（a11y）检测与提取
// 5. 页面结构审计：获取浏览器理解的“增强版DOM”

// CDPPageGetAnnotatedPageContent 获取带浏览器标注的页面内容
func CDPPageGetAnnotatedPageContent() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getAnnotatedPageContent"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getAnnotatedPageContent 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.getAnnotatedPageContent")

	timeout := 8 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					Content string `json:"content"`
					Error   any    `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return respMsg.Content, fmt.Errorf("解析失败: %w", err)
				}
				if res.Error != nil {
					return respMsg.Content, fmt.Errorf("CDP错误: %v", res.Error)
				}
				return res.Content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("getAnnotatedPageContent 超时")
		}
	}
}

/*


// === 使用示例：获取带标注的页面内容并保存 ===
func ExampleGetAnnotatedPageContent() {
	content, err := CDPPageGetAnnotatedPageContent()
	if err != nil {
		log.Fatalf("获取失败: %v", err)
	}

	// 保存为带浏览器标注的HTML
	_ = os.WriteFile("annotated_page.html", []byte(content), 0644)
	log.Println("✅ 已获取带标注页面内容，保存完成")
}

*/

// -----------------------------------------------  Page.getAppId -----------------------------------------------
// === 应用场景 ===
// 1. 获取页面对应的安装应用 ID（PWA/安装的网页应用）
// 2. 识别当前页面是否作为已安装应用运行
// 3. 自动化测试区分浏览器模式 / PWA 模式
// 4. 客户端行为统计：区分普通页面与安装应用的埋点
// 5. PWA 功能调试与验证

// CDPPageGetAppId 获取页面的应用 ID（仅安装的 PWA/网页应用有效）
func CDPPageGetAppId() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP 功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器 WebSocket 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getAppId"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getAppId 请求失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.getAppId")

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					AppId string `json:"appId"`
					Error any    `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return "", fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return "", fmt.Errorf("CDP 错误: %v", res.Error)
				}
				return res.AppId, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("getAppId 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：获取并判断是否为安装应用页面 ===
func ExampleGetAppId() {
	appId, err := CDPPageGetAppId()
	if err != nil {
		log.Fatalf("获取 appId 失败: %v", err)
	}

	if appId != "" {
		log.Printf("✅ 当前页面是已安装的应用，AppId: %s", appId)
	} else {
		log.Println("ℹ️ 当前页面是普通浏览器页面，无应用 ID")
	}
}

*/

// -----------------------------------------------  Page.getInstallabilityErrors -----------------------------------------------
// === 应用场景 ===
// 1. PWA 调试：检查页面为什么**不能被安装为 PWA**
// 2. 自动化检测：验证网站是否满足 PWA 安装条件
// 3. 发布前检查：排查 manifest、service worker、HTTPS 等错误
// 4. 兼容性测试：检测浏览器阻止安装的具体原因

// CDPPageGetInstallabilityErrors 获取PWA不可安装的错误列表
func CDPPageGetInstallabilityErrors() ([]InstallabilityError, error) {
	if !DefaultBrowserWS() {
		return nil, fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getInstallabilityErrors"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return nil, fmt.Errorf("发送getInstallabilityErrors失败: %w", err)
	}

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					Errors []InstallabilityError `json:"errors"`
					Error  any                   `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return nil, fmt.Errorf("解析失败: %w", err)
				}
				if res.Error != nil {
					return nil, fmt.Errorf("CDP错误: %v", res.Error)
				}
				return res.Errors, nil
			}
		case <-timer.C:
			return nil, fmt.Errorf("getInstallabilityErrors超时")
		}
	}
}

// InstallabilityError PWA不可安装错误结构
type InstallabilityError struct {
	ErrorId        string                   `json:"errorId"`
	ErrorArguments []InstallabilityErrorArg `json:"errorArguments"`
}

type InstallabilityErrorArg struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

/*

// === 使用场景示例代码：检查PWA是否可安装，并输出错误 ===
func ExampleCheckPwaInstallable() {
	errors, err := CDPPageGetInstallabilityErrors()
	if err != nil {
		log.Fatalf("获取安装错误失败: %v", err)
	}

	if len(errors) == 0 {
		log.Println("✅ 当前页面满足PWA安装条件")
		return
	}

	log.Printf("⚠️ 当前页面不可安装，共 %d 个错误：", len(errors))
	for _, e := range errors {
		log.Printf("  - 错误ID: %s", e.ErrorId)
		for _, arg := range e.ErrorArguments {
			log.Printf("    → %s: %s", arg.Name, arg.Value)
		}
	}
}

*/

// -----------------------------------------------  Page.getOriginTrials  -----------------------------------------------
// === 应用场景 ===
// 1. 查看当前页面启用的**Origin Trial**实验特性
// 2. 调试 Web 试验性功能是否正常生效
// 3. 验证第三方 token 与域名、特性匹配情况
// 4. 前端特性检测、自动化兼容性测试
// 5. 排查试验 API 不可用原因

// CDPPageGetOriginTrials 获取页面当前启用的 Origin Trial 列表
func CDPPageGetOriginTrials() ([]OriginTrial, error) {
	if !DefaultBrowserWS() {
		return nil, fmt.Errorf("CDP 功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("浏览器 WebSocket 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	msg := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getOriginTrials"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return nil, fmt.Errorf("发送 getOriginTrials 失败: %w", err)
	}

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					OriginTrials []OriginTrial `json:"originTrials"`
					Error        interface{}   `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return nil, fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return nil, fmt.Errorf("CDP 错误: %v", res.Error)
				}
				return res.OriginTrials, nil
			}

		case <-timer.C:
			return nil, fmt.Errorf("getOriginTrials 请求超时")
		}
	}
}

// OriginTrial 对应浏览器 Origin Trial 结构
type OriginTrial struct {
	TrialName  string `json:"trialName"`
	Status     string `json:"status"`
	ExpiryTime int64  `json:"expiryTime,omitempty"`
}

// -----------------------------------------------  Page.getPermissionsPolicyState  -----------------------------------------------
// === 应用场景 ===
// 1. 读取当前页面（或指定frame）的 **Permissions-Policy (原Feature-Policy)** 完整状态
// 2. 自动化审计：检查 camera、geolocation、microphone、autoplay、sync-xhr 等功能是否被允许/禁止
// 3. 多iframe场景：验证子帧权限是否符合预期（allow 属性、HTTP 头）
// 4. 安全测试：排查第三方资源权限泄露、权限配置错误
// 5. 调试：为什么某 API 报 NotAllowedError

// CDPPageGetPermissionsPolicyState 获取页面/frame 的权限策略状态
func CDPPageGetPermissionsPolicyState(frameID string) ([]PermissionsPolicyFeature, error) {
	if !DefaultBrowserWS() {
		return nil, fmt.Errorf("CDP 功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("浏览器 WebSocket 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 可选参数 frameId：指定 frame，不填则为当前主 frame
	var params string
	if frameID != "" {
		params = fmt.Sprintf(`, "params": {"frameId": "%s"}`, frameID)
	}

	msg := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getPermissionsPolicyState"
		%s
	}`, reqID, params)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return nil, fmt.Errorf("发送 getPermissionsPolicyState 失败: %w", err)
	}

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					Features []PermissionsPolicyFeature `json:"features"`
					Error    interface{}                `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return nil, fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return nil, fmt.Errorf("CDP 错误: %v", res.Error)
				}
				return res.Features, nil
			}

		case <-timer.C:
			return nil, fmt.Errorf("getPermissionsPolicyState 请求超时")
		}
	}
}

// PermissionsPolicyFeature 单个权限策略特性的完整结构
type PermissionsPolicyFeature struct {
	FeatureName      string   `json:"featureName"`      // 如: camera, geolocation, microphone, autoplay, clipboard-read
	AllowedByDefault bool     `json:"allowedByDefault"` // 浏览器默认是否允许
	AllowedOrigins   []string `json:"allowedOrigins"`   // 显式允许的源（self, *, https://a.com）
	BlockedByOrigins []string `json:"blockedByOrigins"` // 显式禁止的源
	IsEffective      bool     `json:"isEffective"`      // 当前是否生效（被策略/iframe 控制）
	IframePolicy     string   `json:"iframePolicy"`     // iframe allow 属性中的规则
}

/*

func ExamplePrintPermissionsPolicy() {
	// 获取主 frame 权限
	features, err := CDPPageGetPermissionsPolicyState("")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("=== Permissions Policy (共 %d 项) ===\n", len(features))
	for _, f := range features {
		status := "❌ 禁止"
		if f.IsEffective && len(f.AllowedOrigins) > 0 {
			status = "✅ 允许"
		}
		fmt.Printf("%-20s %s | 默认:%t | 允许:%v | iframe:%s\n",
			f.FeatureName,
			status,
			f.AllowedByDefault,
			f.AllowedOrigins,
			f.IframePolicy,
		)
	}
}

*/

// -----------------------------------------------  Page.getResourceContent -----------------------------------------------
// === 应用场景 ===
// 1. 获取页面已加载的任意资源：JS、CSS、图片、字体、XHR、fetch 等
// 2. 爬虫资源提取：不重新请求，直接从浏览器内存读取资源内容
// 3. 调试：查看当前页面实际加载的资源原文
// 4. 对比资源：检查缓存/网络返回的资源是否一致
// 5. 提取图片/字体二进制（base64）

// CDPPageGetResourceContent 获取指定资源的内容
// 参数：
//   - frameId: 框架ID（可选，传空使用主框架）
//   - url: 资源URL（必须与Page.getResources中返回的完全一致）
//
// 返回：
//   - content: 资源文本内容 或 base64 编码的二进制
//   - base64: 是否为base64编码
func CDPPageGetResourceContent(frameId, url string) (content string, base64 bool, err error) {
	if !DefaultBrowserWS() {
		return "", false, fmt.Errorf("CDP 功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", false, fmt.Errorf("浏览器 WebSocket 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]interface{}{
		"url": url,
	}
	if frameId != "" {
		params["frameId"] = frameId
	}
	paramBytes, _ := json.Marshal(params)

	msg := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getResourceContent",
		"params": %s
	}`, reqID, paramBytes)

	// 发送
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return "", false, fmt.Errorf("发送请求失败: %w", err)
	}

	// 等待响应
	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", false, fmt.Errorf("消息队列关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					Content string      `json:"content"`
					Base64  bool        `json:"base64"`
					Error   interface{} `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return "", false, fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return "", false, fmt.Errorf("CDP 错误: %v", res.Error)
				}
				return res.Content, res.Base64, nil
			}

		case <-timer.C:
			return "", false, fmt.Errorf("getResourceContent 请求超时")
		}
	}
}

/*


// ------------------------------ 用法示例 ------------------------------
// 1. 获取 JS/CSS 文本内容
func ExampleGetResourceJS() {
	url := "https://example.com/script.js"
	content, isBase64, err := CDPPageGetResourceContent("", url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("✅ JS 内容：", content[:50])
}

// 2. 获取图片（二进制 base64）
func ExampleGetResourceImage() {
	url := "https://example.com/logo.png"
	data, isBase64, err := CDPPageGetResourceContent("", url)
	if err != nil {
		log.Fatal(err)
	}
	if isBase64 {
		// 解码保存图片
		imgBytes, _ := base64.StdEncoding.DecodeString(data)
		_ = os.WriteFile("logo.png", imgBytes, 0644)
		log.Println("✅ 图片已保存")
	}
}

*/

// -----------------------------------------------  Page.getResourceTree -----------------------------------------------
// === 应用场景 ===
// 1. 获取页面完整资源树（主框架 + 所有 iframe + 所有加载资源）
// 2. 爬虫：获取页面所有资源（JS/CSS/IMG/FONT/XHR/ Fetch 等）
// 3. 调试：查看页面资源加载层级关系
// 4. 性能审计：分析资源依赖、框架结构、资源类型分布
// 5. 配合 getResourceContent 实现完整资源抓取

// CDPPageGetResourceTree 获取当前页面的完整资源树
func CDPPageGetResourceTree() (*ResourceTree, error) {
	if !DefaultBrowserWS() {
		return nil, fmt.Errorf("CDP 功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("浏览器 WebSocket 未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	msg := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.getResourceTree"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return nil, fmt.Errorf("发送 getResourceTree 失败: %w", err)
	}

	timeout := 8 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					FrameTree *ResourceTree `json:"frameTree"`
					Error     interface{}   `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return nil, fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return nil, fmt.Errorf("CDP 错误: %v", res.Error)
				}
				return res.FrameTree, nil
			}

		case <-timer.C:
			return nil, fmt.Errorf("getResourceTree 请求超时")
		}
	}
}

// ------------------------------ 核心结构体定义 ------------------------------
type ResourceTree struct {
	Frame       *Frame          `json:"frame"`
	Resources   []*Resource     `json:"resources"`
	ChildFrames []*ResourceTree `json:"childFrames,omitempty"`
}

type Frame struct {
	Id             string `json:"id"`
	ParentId       string `json:"parentId,omitempty"`
	Name           string `json:"name,omitempty"`
	Url            string `json:"url"`
	SecurityOrigin string `json:"securityOrigin"`
	MimeType       string `json:"mimeType"`
}

type Resource struct {
	Url         string `json:"url"`
	Type        string `json:"type"`
	MimeType    string `json:"mimeType"`
	ContentSize int    `json:"contentSize,omitempty"`
	Failed      bool   `json:"failed,omitempty"`
	Canceled    bool   `json:"canceled,omitempty"`
}

/*

// ------------------------------ 使用示例 ------------------------------
func ExampleGetResourceTree() {
	// 获取完整资源树
	tree, err := CDPPageGetResourceTree()
	if err != nil {
		log.Fatalf("获取资源树失败: %v", err)
	}

	// 打印主框架信息
	log.Printf("✅ 主框架 URL: %s", tree.Frame.Url)
	log.Printf("✅ 主框架资源总数: %d", len(tree.Resources))

	// 遍历所有资源
	for _, res := range tree.Resources {
		log.Printf("资源: %s | 类型: %s | MIME: %s",
			res.Url, res.Type, res.MimeType)
	}

	// 遍历子框架（iframe）
	for _, child := range tree.ChildFrames {
		log.Printf("=== IFrame: %s ===", child.Frame.Url)
	}
}

*/

// -----------------------------------------------  Page.produceCompilationCache  -----------------------------------------------
// === 应用场景 ===
// 1. 为指定 JavaScript 脚本生成 V8 编译缓存，提升重复加载执行速度
// 2. 自动化/爬虫性能优化：预缓存公共库，减少重复解析开销
// 3. 前端性能测试：模拟编译缓存环境，验证页面加载速度

// CDPPageProduceCompilationCache 触发指定脚本的编译缓存生成
func CDPPageProduceCompilationCache(scripts []ProduceCompilationCacheScript) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	scriptBytes, _ := json.Marshal(scripts)
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.produceCompilationCache",
		"params": {
			"scripts": %s
		}
	}`, reqID, scriptBytes)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 produceCompilationCache 请求失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.produceCompilationCache")

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
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("produceCompilationCache 请求超时")
		}
	}
}

// ProduceCompilationCacheScript 生成编译缓存的脚本参数
type ProduceCompilationCacheScript struct {
	URL string `json:"url"`
}

/*


// === 使用场景示例代码：为指定 JS 生成编译缓存 ===
func ExampleProduceCompilationCache() {
	scripts := []ProduceCompilationCacheScript{
		{URL: "https://example.com/main.js"},
		{URL: "https://example.com/libs/jquery.min.js"},
	}

	resp, err := CDPPageProduceCompilationCache(scripts)
	if err != nil {
		log.Fatalf("生成编译缓存失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 已触发编译缓存生成，等待 Page.compilationCacheProduced 事件")
}

*/

// -----------------------------------------------  Page.screencastFrameAck  -----------------------------------------------
// === 应用场景 ===
// 1. 确认接收屏幕录制帧，控制录屏流帧率
// 2. 启动 Page.startScreencast 后必须使用 ack 同步帧
// 3. 避免录屏帧堆积，稳定录屏性能
// 4. 远程桌面、实时预览、屏幕监控必备

// CDPPageScreencastFrameAck 确认接收 screencast 帧
// 参数 frameNumber：接收到的帧编号（从事件 Page.screencastFrame 中获取）
func CDPPageScreencastFrameAck(frameNumber int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.screencastFrameAck",
		"params": {
			"frameNumber": %d
		}
	}`, reqID, frameNumber)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 screencastFrameAck 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.screencastFrameAck: frameNumber=%d", frameNumber)

	timeout := 3 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("screencastFrameAck 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：录屏帧接收 + ACK 确认 ---------------------------
func ExampleHandleScreencastFrame(frameNumber int) {
	// 处理帧数据...

	// 必须发送 ACK，浏览器才会继续发送下一帧
	_, err := CDPPageScreencastFrameAck(frameNumber)
	if err != nil {
		log.Printf("ACK 确认失败: %v", err)
		return
	}
	log.Println("✅ 帧已确认:", frameNumber)
}

*/

// -----------------------------------------------  Page.searchInResource  -----------------------------------------------
// === 应用场景 ===
// 1. 在指定页面资源（JS/CSS/HTML等）中搜索关键词
// 2. 爬虫：快速定位资源中的关键内容、链接、配置
// 3. 调试：查找脚本中的变量、函数、字符串位置
// 4. 安全审计：检测敏感信息、接口地址、密钥等

// CDPPageSearchInResource 在指定资源中搜索文本
// 参数：
//
//	frameId - 框架ID（为空则使用主框架）
//	url - 资源URL
//	query - 搜索关键词
//	caseSensitive - 是否区分大小写
//	isRegex - 是否为正则表达式搜索
//
// 返回：搜索结果列表
func CDPPageSearchInResource(frameId, url, query string, caseSensitive, isRegex bool) ([]SearchMatch, error) {
	if !DefaultBrowserWS() {
		return nil, fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return nil, fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]interface{}{
		"url":           url,
		"query":         query,
		"caseSensitive": caseSensitive,
		"isRegex":       isRegex,
	}
	if frameId != "" {
		params["frameId"] = frameId
	}
	paramBytes, _ := json.Marshal(params)

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.searchInResource",
		"params": %s
	}`, reqID, paramBytes)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return nil, fmt.Errorf("发送 searchInResource 请求失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.searchInResource: url=%s, query=%s", url, query)

	// 等待响应
	timeout := 8 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				var res struct {
					Result []SearchMatch `json:"result"`
					Error  interface{}   `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return nil, fmt.Errorf("解析响应失败: %w", err)
				}
				if res.Error != nil {
					return nil, fmt.Errorf("CDP错误: %v", res.Error)
				}
				return res.Result, nil
			}
		case <-timer.C:
			return nil, fmt.Errorf("searchInResource 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：在JS脚本中搜索关键词 ===
func ExampleSearchInResource() {
	// 在指定JS中搜索 "apiEndpoint"
	matches, err := CDPPageSearchInResource(
		"",
		"https://example.com/app.js",
		"apiEndpoint",
		false,
		false,
	)
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}

	log.Printf("✅ 搜索到 %d 个匹配项", len(matches))
	for _, match := range matches {
		log.Printf("行 %d: %s", match.LineNumber+1, match.LineContent)
	}
}

*/

// -----------------------------------------------  Page.setAdBlockingEnabled  -----------------------------------------------
// === 应用场景 ===
// 1. 启用/禁用浏览器内置广告拦截功能
// 2. 爬虫/自动化：屏蔽广告，提升页面加载速度
// 3. 测试广告拦截效果、页面纯净度
// 4. 避免广告弹窗干扰自动化操作

// CDPPageSetAdBlockingEnabled 启用或禁用页面广告拦截
// 参数：enabled - true=开启广告拦截，false=关闭广告拦截
func CDPPageSetAdBlockingEnabled(enabled bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setAdBlockingEnabled",
		"params": {
			"enabled": %t
		}
	}`, reqID, enabled)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setAdBlockingEnabled 请求失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.setAdBlockingEnabled: enabled=%t", enabled)

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("setAdBlockingEnabled 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：开启广告拦截，净化页面 ===
func ExampleSetAdBlockingEnabled() {
	// 开启广告拦截
	resp, err := CDPPageSetAdBlockingEnabled(true)
	if err != nil {
		log.Fatalf("开启广告拦截失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 浏览器广告拦截已启用")
}

*/

// -----------------------------------------------  Page.setFontFamilies  -----------------------------------------------
// === 应用场景 ===
// 1. 强制覆盖页面默认字体家族
// 2. 自动化截图/测试时统一字体渲染
// 3. 解决跨平台字体差异、乱码、显示不一致
// 4. 无障碍适配、自定义字体替换

// CDPPageSetFontFamilies 设置页面默认字体
func CDPPageSetFontFamilies(fonts FontFamilies) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	paramBytes, _ := json.Marshal(fonts)
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setFontFamilies",
		"params": %s
	}`, reqID, paramBytes)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setFontFamilies 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.setFontFamilies")

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("setFontFamilies 请求超时")
		}
	}
}

// FontFamilies 字体配置结构（全部为可选字段）
type FontFamilies struct {
	Standard   string `json:"standard,omitempty"`
	Serif      string `json:"serif,omitempty"`
	SansSerif  string `json:"sansSerif,omitempty"`
	Monospace  string `json:"monospace,omitempty"`
	Cursive    string `json:"cursive,omitempty"`
	Fantasy    string `json:"fantasy,omitempty"`
	Pictograph string `json:"pictograph,omitempty"`
}

/*

// === 使用示例：统一设置页面字体为微软雅黑 + 等线 ===
func ExampleSetFontFamilies() {
	fonts := FontFamilies{
		Standard:  "Microsoft YaHei, sans-serif",
		SansSerif: "Microsoft YaHei, sans-serif",
		Monospace: "Consolas, monospace",
	}

	resp, err := CDPPageSetFontFamilies(fonts)
	if err != nil {
		log.Fatalf("设置字体失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面默认字体已替换")
}

*/

// -----------------------------------------------  Page.setFontSizes  -----------------------------------------------
// === 应用场景 ===
// 1. 强制设置页面默认字体尺寸
// 2. 自动化截图、UI测试时统一渲染效果
// 3. 解决不同平台字体大小不一致问题
// 4. 无障碍适配、调整页面基础字号

// CDPPageSetFontSizes 设置页面字体大小
func CDPPageSetFontSizes(fontSizes FontSizes) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	paramBytes, _ := json.Marshal(fontSizes)
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setFontSizes",
		"params": %s
	}`, reqID, paramBytes)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setFontSizes 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.setFontSizes")

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("setFontSizes 请求超时")
		}
	}
}

// FontSizes 字体大小配置结构（全部为可选字段）
type FontSizes struct {
	Standard int `json:"standard,omitempty"` // 标准字体大小
	Fixed    int `json:"fixed,omitempty"`    // 等宽字体大小
}

/*

// === 使用场景示例代码：设置页面标准字体和等宽字体大小 ===
func ExampleSetFontSizes() {
	fontSizes := FontSizes{
		Standard: 16,
		Fixed:    14,
	}

	resp, err := CDPPageSetFontSizes(fontSizes)
	if err != nil {
		log.Fatalf("设置字体大小失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面字体大小已设置")
}

*/

// -----------------------------------------------  Page.setPrerenderingAllowed  -----------------------------------------------
// === 应用场景 ===
// 1. 允许/禁止页面预渲染（Prerendering）
// 2. 自动化测试：控制预渲染行为，避免干扰测试流程
// 3. 性能调试：开启/关闭预渲染对比加载速度
// 4. 资源控制：禁止预渲染节省带宽与内存

// CDPPageSetPrerenderingAllowed 允许或禁止页面预渲染
func CDPPageSetPrerenderingAllowed(allowed bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setPrerenderingAllowed",
		"params": {
			"allowed": %t
		}
	}`, reqID, allowed)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setPrerenderingAllowed 请求失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.setPrerenderingAllowed: allowed=%t", allowed)

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("setPrerenderingAllowed 请求超时")
		}
	}
}

/*


// === 使用场景示例代码：禁止页面预渲染 ===
func ExampleSetPrerenderingAllowed() {
	resp, err := CDPPageSetPrerenderingAllowed(false)
	if err != nil {
		log.Fatalf("禁止预渲染失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面预渲染已被禁止")
}

*/

// -----------------------------------------------  Page.setRPHRegistrationMode  -----------------------------------------------
// === 应用场景 ===
// RPH = registerProtocolHandler (网页协议注册处理)
// 1. 自动处理网页调用 navigator.registerProtocolHandler() 时的授权弹窗
// 2. 自动化测试中避免手动确认协议注册，稳定执行流程
// 3. 控制浏览器是否自动允许/拒绝网页注册自定义协议（如 mailto、web+自定义协议）
// 4. 兼容 WebDriver 标准的 RPH 自动化模式

// RPHRegistrationMode 协议注册授权模式
type RPHRegistrationMode string

const (
	RPHModeAutoAccept RPHRegistrationMode = "autoAccept" // 自动同意所有协议注册
	RPHModeAutoReject RPHRegistrationMode = "autoReject" // 自动拒绝所有协议注册
	RPHModeNone       RPHRegistrationMode = "none"       // 默认行为，弹出用户确认框
)

// CDPPageSetRPHRegistrationMode 设置网页协议注册（RPH）的自动处理模式
func CDPPageSetRPHRegistrationMode(mode RPHRegistrationMode) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setRPHRegistrationMode",
		"params": {
			"mode": "%s"
		}
	}`, reqID, mode)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setRPHRegistrationMode 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.setRPHRegistrationMode: mode=%s", mode)

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("setRPHRegistrationMode 请求超时")
		}
	}
}

/*

// === 使用场景示例：自动化测试自动允许协议注册 ===
func ExampleSetRPHRegistrationMode() {
	// 自动允许网页注册协议（避免弹窗）
	resp, err := CDPPageSetRPHRegistrationMode(RPHModeAutoAccept)
	if err != nil {
		log.Fatalf("设置RPH模式失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ RPH模式已设为自动同意")

	// 恢复默认（弹出确认框）
	// resp, err := CDPPageSetRPHRegistrationMode(RPHModeNone)
}

*/

// -----------------------------------------------  Page.setSPCTransactionMode  -----------------------------------------------
// === 应用场景 ===
// SPC = Secure Payment Confirmation (安全支付确认)
// 1. 控制网页SPC安全支付API的自动授权行为
// 2. 支付自动化测试：自动允许/拒绝支付确认弹窗
// 3. 金融自动化：稳定处理支付流程，避免人工干预
// 4. 兼容W3C SPC标准：控制navigator.securePaymentConfirmation调用

// SPCTransactionMode SPC交易授权模式
type SPCTransactionMode string

const (
	SPCModeAutoAllow SPCTransactionMode = "autoAllow" // 自动允许所有SPC支付确认
	SPCModeAutoDeny  SPCTransactionMode = "autoDeny"  // 自动拒绝所有SPC支付确认
	SPCModeDefault   SPCTransactionMode = "default"   // 默认行为，弹出用户确认框
)

// CDPPageSetSPCTransactionMode 设置安全支付确认(SPC)的交易模式
func CDPPageSetSPCTransactionMode(mode SPCTransactionMode) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setSPCTransactionMode",
		"params": {
			"mode": "%s"
		}
	}`, reqID, mode)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setSPCTransactionMode 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.setSPCTransactionMode: mode=%s", mode)

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("setSPCTransactionMode 请求超时")
		}
	}
}

/*


// === 使用场景示例：支付自动化自动允许SPC确认 ===
func ExampleSetSPCTransactionMode() {
	// 自动允许安全支付确认（无弹窗）
	resp, err := CDPPageSetSPCTransactionMode(SPCModeAutoAllow)
	if err != nil {
		log.Fatalf("设置SPC模式失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ SPC模式已设为自动允许")

	// 自动拒绝支付确认
	// resp, err := CDPPageSetSPCTransactionMode(SPCModeAutoDeny)

	// 恢复默认（弹出确认框）
	// resp, err := CDPPageSetSPCTransactionMode(SPCModeDefault)
}

*/

// -----------------------------------------------  Page.setWebLifecycleState  -----------------------------------------------
// === 应用场景 ===
// 1. 手动设置页面 Web 生命周期状态（激活/冻结/休眠）
// 2. 自动化测试：验证页面在 freeze、resume 等状态下的行为
// 3. 性能/内存优化：模拟页面被浏览器挂起、冻结
// 4. 调试页面在后台、休眠时的逻辑

// WebLifecycleState 网页生命周期状态
type WebLifecycleState string

const (
	WebLifecycleActive    WebLifecycleState = "active"    // 正常激活状态
	WebLifecycleFrozen    WebLifecycleState = "frozen"    // 冻结（CPU/定时器暂停）
	WebLifecycleDiscarded WebLifecycleState = "discarded" // 丢弃/休眠
)

// CDPPageSetWebLifecycleState 手动设置页面生命周期状态
func CDPPageSetWebLifecycleState(state WebLifecycleState) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.setWebLifecycleState",
		"params": {
			"state": "%s"
		}
	}`, reqID, state)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setWebLifecycleState 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.setWebLifecycleState: state=%s", state)

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("setWebLifecycleState 请求超时")
		}
	}
}

/*

// === 使用场景示例：冻结页面 → 恢复激活 ===
func ExampleSetWebLifecycleState() {
	// 冻结页面（模拟后台挂起）
	resp, err := CDPPageSetWebLifecycleState(WebLifecycleFrozen)
	if err != nil {
		log.Fatalf("设置生命周期失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面已冻结")

	// 恢复正常
	// resp, err = CDPPageSetWebLifecycleState(WebLifecycleActive)
}

*/

// -----------------------------------------------  Page.startScreencast  -----------------------------------------------
// === 应用场景 ===
// 1. 启动页面实时屏幕录制流（无头浏览器也支持）
// 2. 远程桌面、实时预览、直播推流、操作回放
// 3. 自动化监控、可视化调试
// 4. 搭配 screencastFrameAck 实现稳定帧率

// CDPPageStartScreencast 启动页面屏幕广播
// 参数：
//
//	format - 图片格式：jpeg / png
//	quality - 图片质量 0-100
//	maxWidth - 最大宽度（可选，0=不限制）
//	maxHeight - 最大高度（可选，0=不限制）
func CDPPageStartScreencast(format string, quality int, maxWidth, maxHeight int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	params := map[string]interface{}{
		"format":    format,
		"quality":   quality,
		"maxWidth":  maxWidth,
		"maxHeight": maxHeight,
	}
	paramBytes, _ := json.Marshal(params)

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.startScreencast",
		"params": %s
	}`, reqID, paramBytes)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 startScreencast 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.startScreencast")

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("startScreencast 请求超时")
		}
	}
}

/*

// === 使用场景示例：启动高质量屏幕广播 ===
func ExampleStartScreencast() {
	// 格式：jpeg，质量80，尺寸不限
	resp, err := CDPPageStartScreencast("jpeg", 80, 0, 0)
	if err != nil {
		log.Fatalf("启动录屏失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面屏幕广播已启动")
}

*/

// -----------------------------------------------  Page.stopScreencast  -----------------------------------------------
// === 应用场景 ===
// 1. 停止页面实时屏幕广播
// 2. 结束录屏、释放资源
// 3. 切换页面/任务时关闭 screencast 流

// CDPPageStopScreencast 停止页面屏幕广播
func CDPPageStopScreencast() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.stopScreencast"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stopScreencast 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.stopScreencast")

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("stopScreencast 请求超时")
		}
	}
}

/*

// === 使用场景示例：停止页面屏幕广播 ===
func ExampleStopScreencast() {
	resp, err := CDPPageStopScreencast()
	if err != nil {
		log.Fatalf("停止录屏失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面屏幕广播已停止")
}

*/

// -----------------------------------------------  Page.waitForDebugger  -----------------------------------------------
// === 应用场景 ===
// 1. 让页面加载后**暂停执行**，等待调试器连接
// 2. 自动化调试：在页面初始化前断点
// 3. 调试早期 JS 执行、页面加载生命周期

// CDPPageWaitForDebugger 使页面暂停等待调试器连接
func CDPPageWaitForDebugger() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket未连接")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Page.waitForDebugger"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 waitForDebugger 失败: %w", err)
	}

	log.Printf("[DEBUG] CDP → Page.waitForDebugger")

	timeout := 15 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				var res map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &res); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}
				if errObj, has := res["error"]; has {
					return content, fmt.Errorf("CDP错误: %v", errObj)
				}
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("waitForDebugger 请求超时")
		}
	}
}

/*

// === 使用场景示例：页面加载后暂停，等待调试器 ===
func ExampleWaitForDebugger() {
	resp, err := CDPPageWaitForDebugger()
	if err != nil {
		log.Fatalf("等待调试器失败: %v, 响应: %s", err, resp)
	}
	log.Println("✅ 页面已暂停，等待调试器连接...")
}

*/
