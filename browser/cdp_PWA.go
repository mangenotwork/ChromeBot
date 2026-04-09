package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  PWA.changeAppUserSettings  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 自动化测试：动态修改 PWA 应用用户设置（通知、安装状态等）
// 2. 调试 PWA 权限：测试不同用户设置下 PWA 的行为表现
// 3. 自动化配置：批量为 PWA 配置用户偏好设置
// 4. 测试环境模拟：模拟用户开启/关闭 PWA 相关权限
// 5. 端到端测试：验证 PWA 在不同用户设置下的兼容性
// 6. 开发调试：快速切换 PWA 用户配置，无需手动操作浏览器

// CDPPWAChangeAppUserSettings 修改 PWA 应用用户设置
// 参数：settings - PWA 用户设置结构体（遵循 CDP 协议定义）
func CDPPWAChangeAppUserSettings(settings map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 序列化参数
	paramsBytes, err := json.Marshal(settings)
	if err != nil {
		return "", fmt.Errorf("序列化 PWA 用户设置参数失败: %w", err)
	}

	// 构建 CDP 请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "PWA.changeAppUserSettings",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送 WebSocket 消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 changeAppUserSettings 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查 CDP 错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("changeAppUserSettings 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础用法 - 启用 PWA 通知权限 ====================
func ExampleCDPPWAChangeAppUserSettings_EnableNotifications() {
	// 构建 PWA 用户设置（遵循 CDP PWA 协议参数）
	settings := map[string]interface{}{
		"notificationsEnabled": true,  // 启用通知
		"appInstalled":         false, // 标记应用未安装
	}

	// 调用函数
	resp, err := CDPPWAChangeAppUserSettings(settings)
	if err != nil {
		log.Fatalf("修改 PWA 用户设置失败: %v", err)
	}
	log.Printf("修改成功，响应: %s", resp)
}

// ==================== 使用示例 2：完整配置 - 模拟 PWA 已安装并开启所有权限 ====================
func ExampleCDPPWAChangeAppUserSettings_FullConfig() {
	settings := map[string]interface{}{
		"notificationsEnabled": true,
		"contentScriptsEnabled": true,
		"autoLaunchEnabled":     true,
		"appInstalled":          true,
	}

	resp, err := CDPPWAChangeAppUserSettings(settings)
	if err != nil {
		log.Printf("修改失败: %v", err)
		return
	}
	log.Println("PWA 用户设置已完整配置")
}

*/

// -----------------------------------------------  PWA.getOsAppState  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 状态检测：获取当前操作系统中 PWA 应用的安装/运行状态
// 2. 自动化测试：验证 PWA 是否已成功安装到操作系统
// 3. 调试诊断：排查 PWA 安装、启动、系统集成相关问题
// 4. 跨平台校验：获取 Windows/macOS/Linux 下的 PWA 系统状态
// 5. 测试前置检查：在执行 PWA 测试前确认系统状态是否符合预期
// 6. 自动化监控：持续监测 PWA 与操作系统的集成状态

// CDPPWAGetOsAppState 获取 PWA 应用在操作系统中的状态
func CDPPWAGetOsAppState() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 CDP 请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "PWA.getOsAppState"
	}`, reqID)

	// 发送 WebSocket 消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getOsAppState 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查 CDP 错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("getOsAppState 请求超时")
		}
	}
}

/*


// ==================== 使用示例 1：基础用法 - 获取并打印 PWA 操作系统状态 ====================
func ExampleCDPPWAGetOsAppState_Basic() {
	// 调用函数获取 PWA 系统状态
	resp, err := CDPPWAGetOsAppState()
	if err != nil {
		log.Fatalf("获取 PWA 操作系统状态失败: %v", err)
	}

	// 打印完整状态信息
	log.Printf("PWA 操作系统状态: \n%s", resp)
}

// ==================== 使用示例 2：测试场景 - 检查 PWA 是否已安装到系统 ====================
func ExampleCDPPWAGetOsAppState_CheckInstalled() {
	resp, err := CDPPWAGetOsAppState()
	if err != nil {
		log.Printf("获取状态失败: %v", err)
		return
	}

	// 解析响应判断安装状态
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err == nil {
		if result["result"] != nil {
			osState := result["result"].(map[string]interface{})
			isInstalled := osState["isAppInstalled"].(bool)
			log.Printf("PWA 是否已安装到系统: %t", isInstalled)
		}
	}
}

*/

// -----------------------------------------------  PWA.install  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 自动化安装：无需手动点击浏览器安装按钮，自动完成 PWA 安装
// 2. 端到端测试：自动化测试 PWA 安装流程、安装后功能
// 3. 批量部署：批量为多个浏览器安装指定 PWA 应用
// 4. 调试安装问题：模拟用户安装行为，排查 PWA 安装失败、异常问题
// 5. 演示环境：快速安装 PWA 用于演示、功能展示
// 6. 自动化验收：验证页面是否满足 PWA 安装条件

// CDPPWAInstall 触发 PWA 应用安装
func CDPPWAInstall() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 CDP 请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "PWA.install"
	}`, reqID)

	// 发送 WebSocket 消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 PWA.install 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查 CDP 错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("PWA.install 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础用法 - 自动安装当前页面 PWA ====================
func ExampleCDPPWAInstall_Basic() {
	// 调用函数触发 PWA 安装
	resp, err := CDPPWAInstall()
	if err != nil {
		log.Fatalf("PWA 安装失败: %v", err)
	}

	log.Printf("PWA 安装请求发送成功，响应: %s", resp)
}

// ==================== 使用示例 2：自动化测试 - 安装后校验安装状态 ====================
func ExampleCDPPWAInstall_InstallAndCheck() {
	// 1. 执行安装
	_, err := CDPPWAInstall()
	if err != nil {
		log.Printf("安装失败: %v", err)
		return
	}

	// 等待安装完成
	time.Sleep(2 * time.Second)

	// 2. 获取系统状态验证是否安装成功
	resp, err := CDPPWAGetOsAppState()
	if err != nil {
		log.Printf("获取状态失败: %v", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err == nil {
		if res, ok := result["result"].(map[string]interface{}); ok {
			installed := res["isAppInstalled"].(bool)
			if installed {
				log.Println("✅ PWA 安装成功并已验证")
			} else {
				log.Println("❌ PWA 安装未完成")
			}
		}
	}
}

*/

// -----------------------------------------------  PWA.launch  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 自动化启动：无需手动操作，自动启动已安装的 PWA 应用
// 2. 端到端测试：自动化测试 PWA 启动流程、启动后功能
// 3. 自动化验收：验证 PWA 能否正常从系统中启动
// 4. 批量启动：批量启动多个已安装的 PWA 应用
// 5. 调试启动问题：模拟用户启动行为，排查 PWA 启动失败、闪退问题
// 6. 演示环境：自动启动 PWA 用于功能演示、效果展示

// CDPPWALaunch 启动已安装的 PWA 应用
func CDPPWALaunch() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 CDP 请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "PWA.launch"
	}`, reqID)

	// 发送 WebSocket 消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 PWA.launch 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查 CDP 错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("PWA.launch 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础用法 - 直接启动已安装的 PWA ====================
func ExampleCDPPWALaunch_Basic() {
	// 调用函数启动 PWA
	resp, err := CDPPWALaunch()
	if err != nil {
		log.Fatalf("PWA 启动失败: %v", err)
	}

	log.Printf("PWA 启动成功，响应: %s", resp)
}

// ==================== 使用示例 2：完整自动化流程 - 安装 → 等待 → 启动 PWA ====================
func ExampleCDPPWALaunch_InstallAndLaunch() {
	// 1. 安装 PWA
	_, err := CDPPWAInstall()
	if err != nil {
		log.Printf("PWA 安装失败: %v", err)
		return
	}
	log.Println("✅ PWA 安装完成")

	// 等待系统完成安装注册
	time.Sleep(3 * time.Second)

	// 2. 启动 PWA
	resp, err := CDPPWALaunch()
	if err != nil {
		log.Fatalf("PWA 启动失败: %v", err)
	}

	log.Printf("✅ PWA 自动化启动成功，响应: %s", resp)
}

*/

// -----------------------------------------------  PWA.launchFilesInApp  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 文件关联测试：自动化测试 PWA 打开关联文件类型的功能
// 2. 端到端测试：验证 PWA 处理本地文件、系统文件调用的逻辑
// 3. 开发调试：快速调试 PWA 文件启动、文件接收相关功能
// 4. 自动化场景：通过系统调用方式让 PWA 打开指定文件
// 5. 兼容性测试：跨平台测试 PWA 文件关联功能
// 6. 演示场景：自动用 PWA 打开指定文件展示功能

// CDPPWALaunchFilesInApp 启动 PWA 并在应用中打开指定文件
// 参数：filePaths - 要在 PWA 中打开的本地文件路径列表
func CDPPWALaunchFilesInApp(filePaths []string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]interface{}{
		"filePaths": filePaths,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化文件路径参数失败: %w", err)
	}

	// 构建 CDP 请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "PWA.launchFilesInApp",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送 WebSocket 消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 PWA.launchFilesInApp 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查 CDP 错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("PWA.launchFilesInApp 请求超时")
		}
	}
}

/*

// ==================== 使用示例 1：基础用法 - 打开单个文件 ====================
func ExampleCDPPWALaunchFilesInApp_SingleFile() {
	// 单个文件路径
	files := []string{
		"C:\\test\\document.pdf", // Windows 路径示例
		// "/home/user/test.json", // Linux 路径示例
		// "/Users/user/test.txt", // macOS 路径示例
	}

	// 启动 PWA 并打开文件
	resp, err := CDPPWALaunchFilesInApp(files)
	if err != nil {
		log.Fatalf("PWA 打开文件失败: %v", err)
	}

	log.Printf("PWA 打开文件成功，响应: %s", resp)
}

// ==================== 使用示例 2：高级用法 - 批量打开多个文件 ====================
func ExampleCDPPWALaunchFilesInApp_MultiFiles() {
	// 批量文件路径
	files := []string{
		"C:\\files\\data1.json",
		"C:\\files\\report.pdf",
		"C:\\files\\image.png",
	}

	// 启动 PWA 批量打开文件
	resp, err := CDPPWALaunchFilesInApp(files)
	if err != nil {
		log.Printf("批量打开文件失败: %v", err)
		return
	}

	log.Println("✅ PWA 已批量打开所有指定文件")
}

*/

// -----------------------------------------------  PWA.openCurrentPageInApp  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 体验切换：将浏览器中打开的当前页面，切换到 PWA 应用窗口中打开
// 2. 自动化测试：验证 PWA 应用内打开网页的功能完整性
// 3. 调试场景：快速在 PWA 窗口中调试当前网页
// 4. 流程自动化：自动化完成浏览器到 PWA 应用的页面跳转
// 5. 用户体验测试：测试从网页无缝切换到 PWA 应用的交互流程
// 6. 演示场景：展示网页与 PWA 应用之间的页面迁移能力

// CDPPWAOpenCurrentPageInApp 将浏览器当前页面在 PWA 应用中打开
func CDPPWAOpenCurrentPageInApp() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 CDP 请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "PWA.openCurrentPageInApp"
	}`, reqID)

	// 发送 WebSocket 消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 PWA.openCurrentPageInApp 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查 CDP 错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("PWA.openCurrentPageInApp 请求超时")
		}
	}
}

/*


// ==================== 使用示例 1：基础用法 - 浏览器页面切换到 PWA 应用打开 ====================
func ExampleCDPPWAOpenCurrentPageInApp_Basic() {
	// 调用函数将当前页面在 PWA 中打开
	resp, err := CDPPWAOpenCurrentPageInApp()
	if err != nil {
		log.Fatalf("页面切换到 PWA 打开失败: %v", err)
	}

	log.Printf("页面已成功在 PWA 应用中打开，响应: %s", resp)
}

// ==================== 使用示例 2：自动化测试场景 - 启动 PWA 后打开当前页面 ====================
func ExampleCDPPWAOpenCurrentPageInApp_AutoTest() {
	// 1. 先启动 PWA 确保应用已激活
	_, err := CDPPWALaunch()
	if err != nil {
		log.Printf("PWA 启动失败: %v", err)
		return
	}
	time.Sleep(1 * time.Second)

	// 2. 将浏览器当前页面切换到 PWA 应用中打开
	resp, err := CDPPWAOpenCurrentPageInApp()
	if err != nil {
		log.Fatalf("切换页面失败: %v", err)
	}

	log.Println("✅ 自动化完成：PWA启动 + 当前页面在应用中打开成功")
}

*/

// -----------------------------------------------  PWA.uninstall  -----------------------------------------------
// === 应用场景 ===
// 1. PWA 自动化卸载：无需手动操作，自动卸载已安装的 PWA 应用
// 2. 端到端测试：自动化测试 PWA 卸载流程、卸载后状态清理
// 3. 测试环境重置：在自动化测试后卸载 PWA，恢复初始测试环境
// 4. 批量清理：批量卸载多个测试用的 PWA 应用
// 5. 调试卸载问题：模拟用户卸载行为，排查 PWA 卸载失败、残留问题
// 6. 自动化验收：验证 PWA 能否正常从操作系统中卸载

// CDPPWAUninstall 卸载已安装的 PWA 应用
func CDPPWAUninstall() (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 CDP 请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "PWA.uninstall"
	}`, reqID)

	// 发送 WebSocket 消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 PWA.uninstall 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应（5秒超时）
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查 CDP 错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("PWA.uninstall 请求超时")
		}
	}
}

/*


// ==================== 使用示例 1：基础用法 - 直接卸载 PWA 应用 ====================
func ExampleCDPPWAUninstall_Basic() {
	// 调用函数卸载 PWA
	resp, err := CDPPWAUninstall()
	if err != nil {
		log.Fatalf("PWA 卸载失败: %v", err)
	}

	log.Printf("PWA 卸载请求发送成功，响应: %s", resp)
}

// ==================== 使用示例 2：完整自动化流程 - 卸载后校验状态 ====================
func ExampleCDPPWAUninstall_UninstallAndCheck() {
	// 1. 执行卸载
	_, err := CDPPWAUninstall()
	if err != nil {
		log.Printf("卸载失败: %v", err)
		return
	}
	log.Println("✅ PWA 卸载请求已发送")

	// 等待系统完成卸载
	time.Sleep(2 * time.Second)

	// 2. 获取系统状态验证是否卸载成功
	resp, err := CDPPWAGetOsAppState()
	if err != nil {
		log.Printf("获取状态失败: %v", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err == nil {
		if res, ok := result["result"].(map[string]interface{}); ok {
			installed := res["isAppInstalled"].(bool)
			if !installed {
				log.Println("✅ PWA 已成功卸载并验证")
			} else {
				log.Println("❌ PWA 仍处于安装状态，卸载失败")
			}
		}
	}
}


*/
