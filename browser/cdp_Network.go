package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Network.clearBrowserCache  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试: 每次测试前清除浏览器缓存，保证测试环境干净
// 2. 前端调试: 强制清除缓存，确保加载最新的静态资源（JS/CSS/图片）
// 3. 缓存问题修复: 解决因缓存导致的页面加载异常、资源不更新问题
// 4. 性能测试: 清除缓存后重新加载页面，测试首次加载性能
// 5. 自动化爬虫: 避免缓存干扰爬取结果，每次请求都获取最新页面
// 6. 页面刷新重置: 重置浏览器缓存状态，确保页面完全重新加载

// CDPNetworkClearBrowserCache 清除浏览器缓存（对应CDP方法：Network.clearBrowserCache）
func CDPNetworkClearBrowserCache() (string, error) {
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
		"method": "Network.clearBrowserCache"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearBrowserCache 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("clearBrowserCache 请求超时")
		}
	}
}

/*
// === 使用场景示例代码：自动化测试前清除缓存 ===
func ExampleClearCacheForTest() {
	// 场景：UI自动化测试开始前，清除浏览器缓存
	resp, err := CDPNetworkClearBrowserCache()
	if err != nil {
		log.Fatalf("清除缓存失败: %v", err)
	}
	log.Println("缓存清除成功，响应：", resp)

	// 后续执行测试逻辑，确保加载最新页面
}

// === 使用场景示例代码：前端调试强制刷新资源 ===
func ExampleDebugForceRefresh() {
	// 场景：调试前端页面，清除缓存后重新加载
	_, err := CDPNetworkClearBrowserCache()
	if err != nil {
		log.Printf("清除缓存异常: %v", err)
		return
	}
	log.Println("已清除浏览器缓存，可刷新页面加载最新资源")
}

*/

// -----------------------------------------------  Network.clearBrowserCookies  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试: 测试前清除Cookie，保证每次测试都是全新登录状态
// 2. 登录态重置: 强制退出登录，清除用户会话Cookie
// 3. 隐私数据清理: 自动化清理浏览器Cookie，保护隐私
// 4. 多账号测试: 切换测试账号前清除原有Cookie，避免状态冲突
// 5. 爬虫防屏蔽: 清除Cookie后重新请求，避免被服务器识别会话
// 6. 调试Cookie问题: 定位Cookie异常时重置所有Cookie状态

// CDPNetworkClearBrowserCookies 清除浏览器所有Cookie（对应CDP方法：Network.clearBrowserCookies）
func CDPNetworkClearBrowserCookies() (string, error) {
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
		"method": "Network.clearBrowserCookies"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearBrowserCookies 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("clearBrowserCookies 请求超时")
		}
	}
}

/*

// === 使用场景示例代码：自动化测试前重置登录状态 ===
func ExampleClearCookiesForTest() {
	// 场景：UI自动化测试，清除Cookie后重新登录
	resp, err := CDPNetworkClearBrowserCookies()
	if err != nil {
		log.Fatalf("清除Cookie失败: %v", err)
	}
	log.Println("Cookie清除成功，响应：", resp)

	// 后续执行登录测试，无历史登录态干扰
}

// === 使用场景示例代码：多账号测试切换 ===
func ExampleSwitchAccountByClearCookies() {
	// 场景：测试多账号登录，切换前清除Cookie
	_, err := CDPNetworkClearBrowserCookies()
	if err != nil {
		log.Printf("切换账号清除Cookie失败: %v", err)
		return
	}
	log.Println("已清除Cookie，可登录新测试账号")
}

*/

// -----------------------------------------------  Network.deleteCookies  -----------------------------------------------
// === 应用场景 ===
// 1. 精准清除会话: 只删除指定名称的Cookie，不影响其他Cookie
// 2. 登录态精准重置: 仅删除登录会话Cookie，保留其他配置Cookie
// 3. 前端调试: 测试特定Cookie失效后的页面行为
// 4. 接口测试: 手动删除指定Cookie后重新请求接口
// 5. 多用户切换: 精准删除用户标识Cookie
// 6. 爬虫精准控制: 只清除需要刷新的Cookie，避免全量清除

// CDPNetworkDeleteCookies 删除指定名称的Cookie（对应CDP方法：Network.deleteCookies）
// name: Cookie名称
// url: 可选，Cookie所属URL（优先级高于domain/path）
// domain: 可选，Cookie所属域名
// path: 可选，Cookie所属路径
func CDPNetworkDeleteCookies(name, url, domain, path string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := make(map[string]interface{})
	params["name"] = name
	if url != "" {
		params["url"] = url
	}
	if domain != "" {
		params["domain"] = domain
	}
	if path != "" {
		params["path"] = path
	}

	// 序列化为JSON
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.deleteCookies",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deleteCookies 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("deleteCookies 请求超时")
		}
	}
}

/*
// === 使用场景示例：删除指定URL下的Cookie ===
func ExampleDeleteCookieByURL() {
	// 场景：删除指定URL下的session Cookie
	resp, err := CDPNetworkDeleteCookies("session_id", "https://example.com", "", "")
	if err != nil {
		log.Fatalf("删除Cookie失败: %v", err)
	}
	log.Println("按URL删除Cookie成功:", resp)
}

// === 使用场景示例：按域名+路径删除Cookie ===
func ExampleDeleteCookieByDomainPath() {
	// 场景：精准删除指定域名和路径的Cookie
	resp, err := CDPNetworkDeleteCookies("user_token", "", ".example.com", "/api")
	if err != nil {
		log.Fatalf("删除Cookie失败: %v", err)
	}
	log.Println("按域名路径删除Cookie成功:", resp)
}
*/

// -----------------------------------------------  Network.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 停止网络监听: 不再接收Network域的事件通知，减少性能消耗
// 2. 测试流程结束: 自动化测试完成后关闭网络监控
// 3. 资源释放: 释放Network相关的监听和内存资源
// 4. 切换监控模块: 关闭网络监控后专注调试其他模块
// 5. 避免干扰: 防止网络事件干扰后续操作
// 6. 调试控制: 手动控制网络监控的开启与关闭

// CDPNetworkDisable 禁用Network域（对应CDP方法：Network.disable）
func CDPNetworkDisable() (string, error) {
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
		"method": "Network.disable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 disable 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("disable 请求超时")
		}
	}
}

/*
// === 使用场景示例：自动化测试结束后关闭网络监听 ===
func ExampleNetworkDisableAfterTest() {
	// 场景：测试完成，停止网络监控释放资源
	resp, err := CDPNetworkDisable()
	if err != nil {
		log.Fatalf("禁用Network域失败: %v", err)
	}
	log.Println("Network已禁用，停止接收网络事件: ", resp)
}

// === 使用场景示例：释放浏览器资源 ===
func ExampleNetworkDisableForResource() {
	// 场景：无需网络监控时禁用，降低浏览器资源占用
	_, err := CDPNetworkDisable()
	if err != nil {
		log.Printf("禁用Network异常: %v", err)
		return
	}
	log.Println("已禁用网络监控，资源已释放")
}

*/

// -----------------------------------------------  Network.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 启动网络监听: 开启Network域，开始接收网络请求/响应事件
// 2. 自动化测试: 测试前启用网络监控，捕获接口请求数据
// 3. 前端调试: 捕获页面资源加载、接口调用信息
// 4. 爬虫数据采集: 监听网络请求，获取接口返回数据
// 5. 性能分析: 启用后统计页面资源加载耗时
// 6. 接口调试: 实时捕获前端请求参数与响应结果

// CDPNetworkEnable 启用Network域（对应CDP方法：Network.enable）
// maxTotalBufferSize: 最大总缓冲大小（可选，传0使用默认值）
// maxResourceBufferSize: 最大资源缓冲大小（可选，传0使用默认值）
// maxPostDataSize: 最大POST数据大小（可选，传0使用默认值）
func CDPNetworkEnable(maxTotalBufferSize, maxResourceBufferSize, maxPostDataSize int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := make(map[string]interface{})
	if maxTotalBufferSize > 0 {
		params["maxTotalBufferSize"] = maxTotalBufferSize
	}
	if maxResourceBufferSize > 0 {
		params["maxResourceBufferSize"] = maxResourceBufferSize
	}
	if maxPostDataSize > 0 {
		params["maxPostDataSize"] = maxPostDataSize
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.enable",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 enable 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("enable 请求超时")
		}
	}
}

/*
// === 使用场景示例：基础启用网络监听（默认配置） ===
func ExampleNetworkEnableDefault() {
	// 场景：自动化测试/前端调试，启用默认网络监听
	resp, err := CDPNetworkEnable(0, 0, 0)
	if err != nil {
		log.Fatalf("启用Network域失败: %v", err)
	}
	log.Println("Network已启用，开始捕获网络请求: ", resp)
}

// === 使用场景示例：自定义缓冲区大小启用 ===
func ExampleNetworkEnableCustomSize() {
	// 场景：捕获大接口、大表单数据，自定义缓冲区
	// 总缓冲100MB，资源缓冲50MB，POST数据10MB
	resp, err := CDPNetworkEnable(100*1024*1024, 50*1024*1024, 10*1024*1024)
	if err != nil {
		log.Fatalf("启用Network失败: %v", err)
	}
	log.Println("Network已启用(自定义缓冲): ", resp)
}

*/

// -----------------------------------------------  Network.getCookies  -----------------------------------------------
// === 应用场景 ===
// 1. 会话获取: 获取当前浏览器所有Cookie，用于校验登录态、会话信息
// 2. 自动化测试: 断言Cookie是否存在、值是否正确
// 3. 数据持久化: 保存Cookie到本地，下次启动直接复用
// 4. 接口调试: 查看请求携带的Cookie详情，排查鉴权问题
// 5. 爬虫会话保持: 获取Cookie后维持登录状态
// 6. 前端调试: 实时查看浏览器Cookie完整信息

// CDPNetworkGetCookies 获取浏览器所有Cookie（对应CDP方法：Network.getCookies）
func CDPNetworkGetCookies() (string, error) {
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
		"method": "Network.getCookies"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getCookies 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("getCookies 请求超时")
		}
	}
}

/*
// === 使用场景示例：获取所有Cookie并打印 ===
func ExampleGetAllCookies() {
	// 场景：获取当前浏览器全部Cookie列表
	resp, err := CDPNetworkGetCookies()
	if err != nil {
		log.Fatalf("获取Cookie失败: %v", err)
	}
	log.Println("获取Cookie成功，全部Cookie：\n", resp)
}

// === 使用场景示例：测试校验登录Cookie是否存在 ===
func ExampleCheckLoginCookie() {
	// 场景：自动化测试，检查是否包含登录token Cookie
	resp, err := CDPNetworkGetCookies()
	if err != nil {
		log.Printf("获取Cookie异常: %v", err)
		return
	}

	// 简单判断是否包含登录token
	if strings.Contains(resp, "token") || strings.Contains(resp, "session") {
		log.Println("登录状态有效，已找到会话Cookie")
	} else {
		log.Println("未检测到登录Cookie")
	}
}

*/

// -----------------------------------------------  Network.getRequestPostData  -----------------------------------------------
// === 应用场景 ===
// 1. 接口调试: 获取POST请求的原始请求体数据，排查参数传递问题
// 2. 自动化测试: 校验前端发送的POST参数是否符合预期
// 3. 爬虫数据采集: 捕获POST表单/JSON请求的原始数据
// 4. 问题定位: 复现接口异常时，获取真实发送的Post数据
// 5. 日志审计: 记录接口请求的原始参数用于追溯
// 6. 接口Mock: 根据真实Post数据构造Mock请求

// CDPNetworkGetRequestPostData 获取指定请求的POST数据（对应CDP方法：Network.getRequestPostData）
// requestId: 请求ID（通过Network.requestWillBeSent等事件获取）
func CDPNetworkGetRequestPostData(requestId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{"requestId": requestId}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.getRequestPostData",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getRequestPostData 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("getRequestPostData 请求超时")
		}
	}
}

/*
// === 使用场景示例：获取指定接口的POST数据 ===
func ExampleGetRequestPostData() {
	// 场景：从网络事件中获取requestId，查询该请求的POST表单/JSON数据
	requestId := "123456-request-id-7890" // 实际使用时从requestWillBeSent事件获取
	resp, err := CDPNetworkGetRequestPostData(requestId)
	if err != nil {
		log.Fatalf("获取POST数据失败: %v", err)
	}
	log.Println("获取POST数据成功:\n", resp)
}


// === 使用场景示例：自动化测试校验POST参数 ===
func ExampleTestPostData() {
	// 场景：测试提交表单接口，验证发送的参数是否正确
	requestId := "test-request-id-001"
	postDataResp, err := CDPNetworkGetRequestPostData(requestId)
	if err != nil {
		log.Printf("获取请求数据异常: %v", err)
		return
	}

	// 校验是否包含预期参数
	if strings.Contains(postDataResp, "username") && strings.Contains(postDataResp, "password") {
		log.Println("✅ 接口POST参数校验通过")
	} else {
		log.Println("❌ 接口POST参数缺失关键数据")
	}
}

*/

// -----------------------------------------------  Network.getResponseBody  -----------------------------------------------
// === 应用场景 ===
// 1. 接口数据抓取: 获取接口返回的原始响应体数据，用于解析业务数据
// 2. 自动化测试: 校验接口返回结果是否符合预期
// 3. 前端调试: 查看真实接口返回数据，排查前后端数据不一致问题
// 4. 爬虫数据提取: 捕获页面异步请求的JSON/HTML响应数据
// 5. 异常排查: 接口报错时获取真实错误信息
// 6. 数据日志: 记录接口响应数据用于问题追溯

// CDPNetworkGetResponseBody 获取指定请求的响应体数据（对应CDP方法：Network.getResponseBody）
// requestId: 请求ID（通过Network.requestWillBeSent/responseReceived事件获取）
func CDPNetworkGetResponseBody(requestId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{"requestId": requestId}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.getResponseBody",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getResponseBody 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("getResponseBody 请求超时")
		}
	}
}

/*
// === 使用场景示例：抓取接口JSON响应数据 ===
func ExampleGetApiResponseBody() {
	// 场景：获取异步接口的返回数据（如用户信息、列表数据）
	requestId := "api-request-id-123456" // 从网络事件中获取真实ID
	resp, err := CDPNetworkGetResponseBody(requestId)
	if err != nil {
		log.Fatalf("获取响应体失败: %v", err)
	}
	log.Println("接口响应数据：\n", resp)
}



// === 使用场景示例：自动化测试校验接口返回 ===
func ExampleTestApiResponse() {
	// 场景：测试接口是否返回成功状态、预期数据
	requestId := "test-api-request-id"
	respBody, err := CDPNetworkGetResponseBody(requestId)
	if err != nil {
		log.Printf("获取测试接口响应失败: %v", err)
		return
	}

	// 校验返回结果
	if strings.Contains(respBody, "\"code\":0") || strings.Contains(respBody, "success") {
		log.Println("✅ 接口返回成功，测试通过")
	} else {
		log.Println("❌ 接口返回异常，测试失败")
	}
}

*/

// -----------------------------------------------  Network.setBypassServiceWorker  -----------------------------------------------
// === 应用场景 ===
// 1. 前端调试: 跳过ServiceWorker缓存，直接加载服务器最新资源
// 2. 离线功能测试: 验证关闭SW后页面的真实网络请求行为
// 3. 资源更新测试: 确保加载最新静态资源，排除SW缓存干扰
// 4. 接口调试: 绕过SW拦截，直接请求真实后端接口
// 5. 兼容性测试: 测试不支持ServiceWorker环境下的页面表现
// 6. 问题排查: 定位是SW缓存问题还是真实服务端问题

// CDPNetworkSetBypassServiceWorker 设置是否绕过ServiceWorker（对应CDP方法：Network.setBypassServiceWorker）
// bypass: true-绕过SW  false-正常使用SW
func CDPNetworkSetBypassServiceWorker(bypass bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]bool{"bypass": bypass}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setBypassServiceWorker",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setBypassServiceWorker 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setBypassServiceWorker 请求超时")
		}
	}
}

/*
// === 使用场景示例：调试时绕过ServiceWorker ===
func ExampleBypassServiceWorkerForDebug() {
	// 场景：前端开发调试，跳过SW加载最新资源
	resp, err := CDPNetworkSetBypassServiceWorker(true)
	if err != nil {
		log.Fatalf("设置绕过SW失败: %v", err)
	}
	log.Println("已启用绕过ServiceWorker，刷新页面获取最新资源: ", resp)
}

// === 使用场景示例：恢复正常使用ServiceWorker ===
func ExampleEnableServiceWorker() {
	// 场景：调试完成，恢复SW正常缓存功能
	resp, err := CDPNetworkSetBypassServiceWorker(false)
	if err != nil {
		log.Fatalf("恢复SW失败: %v", err)
	}
	log.Println("已恢复正常使用ServiceWorker: ", resp)
}

*/

// -----------------------------------------------  Network.setCacheDisabled  -----------------------------------------------
// === 应用场景 ===
// 1. 前端调试: 禁用缓存强制加载最新资源，解决缓存导致的页面不更新问题
// 2. 自动化测试: 每次运行测试都禁用缓存，确保测试环境干净无缓存干扰
// 3. 性能测试: 测试页面首次加载性能，禁用缓存模拟新用户访问
// 4. 接口调试: 禁用缓存确保每次请求都获取最新接口数据
// 5. 爬虫采集: 禁用缓存避免获取到缓存的旧数据
// 6. 问题复现: 排查缓存相关的页面异常时使用

// CDPNetworkSetCacheDisabled 设置是否禁用浏览器缓存（对应CDP方法：Network.setCacheDisabled）
// disabled: true-禁用缓存  false-启用缓存
func CDPNetworkSetCacheDisabled(disabled bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]bool{"disabled": disabled}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setCacheDisabled",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setCacheDisabled 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setCacheDisabled 请求超时")
		}
	}
}

/*
// === 使用场景示例：调试时禁用缓存 ===
func ExampleDisableCacheForDebug() {
	// 场景：前端开发调试，禁用缓存强制加载最新JS/CSS/接口数据
	resp, err := CDPNetworkSetCacheDisabled(true)
	if err != nil {
		log.Fatalf("禁用缓存失败: %v", err)
	}
	log.Println("已禁用浏览器缓存，刷新页面加载最新资源: ", resp)
}

// === 使用场景示例：恢复缓存功能 ===
func ExampleEnableCache() {
	// 场景：调试完成，恢复正常缓存机制提升加载速度
	resp, err := CDPNetworkSetCacheDisabled(false)
	if err != nil {
		log.Fatalf("恢复缓存失败: %v", err)
	}
	log.Println("已恢复浏览器缓存功能: ", resp)
}

*/

// -----------------------------------------------  Network.setCookie  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化登录: 提前设置登录Cookie，跳过登录界面直接进入系统
// 2. 会话保持: 爬虫/自动化测试中维持用户登录状态
// 3. 前端调试: 手动设置特定Cookie测试不同用户状态
// 4. 多账号切换: 设置不同账号Cookie快速切换测试身份
// 5. 接口测试: 注入鉴权Cookie请求需要登录的接口
// 6. 环境切换: 设置环境Cookie切换测试/预发布环境

// CDPNetworkSetCookie 设置单个Cookie（对应CDP方法：Network.setCookie）
// name: Cookie名称
// value: Cookie值
// url: 可选，Cookie适用URL（优先级高于domain/path）
// domain: 可选，Cookie域名
// path: 可选，Cookie路径
// secure: 可选，是否仅HTTPS传输
// httpOnly: 可选，是否禁止JS访问
// sameSite: 可选，SameSite模式（Strict/Lax/None）
// expires: 可选，过期时间（Unix时间戳，秒）
func CDPNetworkSetCookie(name, value, url, domain, path string, secure, httpOnly bool, sameSite string, expires float64) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := make(map[string]interface{})
	params["name"] = name
	params["value"] = value

	// 可选参数非空时才加入
	if url != "" {
		params["url"] = url
	}
	if domain != "" {
		params["domain"] = domain
	}
	if path != "" {
		params["path"] = path
	}
	params["secure"] = secure
	params["httpOnly"] = httpOnly
	if sameSite != "" {
		params["sameSite"] = sameSite
	}
	if expires > 0 {
		params["expires"] = expires
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setCookie",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setCookie 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setCookie 请求超时")
		}
	}
}

/*
// === 使用场景示例：基础设置登录Cookie ===
func ExampleSetLoginCookie() {
	// 场景：自动化测试直接设置登录态，跳过登录页
	resp, err := CDPNetworkSetCookie(
		"token",                // name
		"user_login_token_123", // value
		"https://example.com",  // url
		"",                    // domain
		"/",                   // path
		true,                  // secure
		true,                  // httpOnly
		"Lax",                 // sameSite
		0,                     // expires
	)
	if err != nil {
		log.Fatalf("设置登录Cookie失败: %v", err)
	}
	log.Println("✅ 登录Cookie设置成功: ", resp)
}

// === 使用场景示例：带域名+过期时间的Cookie ===
func ExampleSetDomainExpireCookie() {
	// 场景：设置主域名Cookie，7天过期
	expireTime := float64(time.Now().Add(7*24*time.Hour).Unix()) // 7天后过期
	resp, err := CDPNetworkSetCookie(
		"uid",
		"10001",
		"",
		".example.com", // 适配子域名
		"/",
		true,
		false,
		"Strict",
		expireTime,
	)
	if err != nil {
		log.Fatalf("设置域名Cookie失败: %v", err)
	}
	log.Println("✅ 域名Cookie设置成功: ", resp)
}

*/

// -----------------------------------------------  Network.setCookies  -----------------------------------------------
// === 应用场景 ===
// 1. 批量会话恢复: 一次性导入多个Cookie，快速恢复登录状态
// 2. 自动化测试: 批量设置测试环境所需的多组Cookie
// 3. 爬虫数据持久化: 保存并批量加载Cookie，维持长期登录
// 4. 多身份切换: 一键切换整套用户身份Cookie
// 5. 环境配置: 批量设置接口鉴权、环境标识等多个Cookie
// 6. 数据迁移: 从其他浏览器/客户端导入整套Cookie

// CDPNetworkSetCookies 批量设置多个Cookie（对应CDP方法：Network.setCookies）
// cookies: Cookie列表，每个Cookie必须包含name/value/url/domain/path等核心参数
func CDPNetworkSetCookies(cookies []map[string]interface{}) (string, error) {
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
		"cookies": cookies,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setCookies",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setCookies 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setCookies 请求超时")
		}
	}
}

/*
// === 使用场景示例：批量设置登录+用户信息Cookie ===
func ExampleBatchSetLoginCookies() {
	// 场景：一次性设置登录token、用户ID、会话Cookie
	cookies := []map[string]interface{}{
		{
			"name":     "token",
			"value":    "login_token_abc123",
			"url":      "https://example.com",
			"secure":   true,
			"httpOnly": true,
			"sameSite": "Lax",
		},
		{
			"name":     "uid",
			"value":    "10001",
			"domain":   ".example.com",
			"path":     "/",
			"secure":   true,
			"sameSite": "Strict",
		},
		{
			"name":     "session_id",
			"value":    "sess_xyz789",
			"url":      "https://example.com",
			"secure":   true,
			"httpOnly": true,
		},
	}

	resp, err := CDPNetworkSetCookies(cookies)
	if err != nil {
		log.Fatalf("批量设置Cookie失败: %v", err)
	}
	log.Println("✅ 批量Cookie设置成功，已恢复登录状态: ", resp)
}

// === 使用场景示例：从本地读取并批量加载Cookie ===
func ExampleLoadCookiesFromFile() {
	// 场景：读取本地保存的Cookie列表，批量导入浏览器
	// 模拟本地存储的Cookie数据
	localCookies := []map[string]interface{}{
		{
			"name":   "theme",
			"value":  "dark",
			"url":    "https://example.com",
			"secure": true,
		},
		{
			"name":   "language",
			"value":  "zh-CN",
			"domain": ".example.com",
			"path":   "/",
		},
	}

	resp, err := CDPNetworkSetCookies(localCookies)
	if err != nil {
		log.Fatalf("加载本地Cookie失败: %v", err)
	}
	log.Println("✅ 本地Cookie批量导入成功: ", resp)
}

*/

// -----------------------------------------------  Network.setExtraHTTPHeaders  -----------------------------------------------
// === 应用场景 ===
// 1. 接口鉴权：全局添加 Token、Authorization 等请求头
// 2. 跨域调试：手动添加 Origin、Referer 等头绕过浏览器限制
// 3. 接口测试：统一添加版本号、设备信息、渠道标识等
// 4. 爬虫伪装：添加 User-Agent、Accept-Language 模拟真实浏览器
// 5. 环境区分：通过请求头标识测试/预发/生产环境
// 6. 日志追踪：全局注入 traceId、requestId 便于链路追踪

// CDPNetworkSetExtraHTTPHeaders 为所有请求设置额外的HTTP请求头（对应CDP方法：Network.setExtraHTTPHeaders）
// headers: 键值对形式的请求头集合，例如 {"User-Agent": "test", "Token": "xxx"}
func CDPNetworkSetExtraHTTPHeaders(headers map[string]string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数（CDP要求 headers 字段值为字符串类型）
	params := map[string]interface{}{
		"headers": headers,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setExtraHTTPHeaders",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setExtraHTTPHeaders 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setExtraHTTPHeaders 请求超时")
		}
	}
}

/*

// === 使用场景示例：全局添加接口鉴权 Token ===
func ExampleSetAuthHeader() {
	// 场景：所有请求自动带上登录token，无需每个接口手动添加
	headers := map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"Source":        "automation-test",
	}

	resp, err := CDPNetworkSetExtraHTTPHeaders(headers)
	if err != nil {
		log.Fatalf("设置全局请求头失败: %v", err)
	}
	log.Println("✅ 全局鉴权请求头设置成功: ", resp)
}

// === 使用场景示例：伪装浏览器请求 + 自定义追踪ID ===
func ExampleSetMockAndTraceHeader() {
	// 场景：爬虫/测试时伪装UA，并添加追踪ID
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0",
		"Referer":    "https://example.com",
		"Trace-Id":   "trace_0000123456",
	}

	resp, err := CDPNetworkSetExtraHTTPHeaders(headers)
	if err != nil {
		log.Fatalf("设置伪装请求头失败: %v", err)
	}
	log.Println("✅ 全局伪装+追踪请求头设置成功: ", resp)
}

*/

// -----------------------------------------------  Network.setUserAgentOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 设备模拟: 模拟手机/平板/PC浏览器访问页面，测试响应式展示
// 2. 爬虫伪装: 伪装成真实浏览器，避免被服务器识别为爬虫
// 3. 兼容性测试: 模拟低版本浏览器、不同内核设备测试页面兼容性
// 4. 接口调试: 服务端根据UA区分设备时，指定UA调试对应逻辑
// 5. 自动化测试: 固定UA，避免环境差异导致测试结果不一致
// 6. 跨平台适配: 模拟微信/支付宝/小程序内置浏览器UA

// CDPNetworkSetUserAgentOverride 设置用户代理UA（对应CDP方法：Network.setUserAgentOverride）
// userAgent: 完整的User-Agent字符串
func CDPNetworkSetUserAgentOverride(userAgent string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{"userAgent": userAgent}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setUserAgentOverride",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setUserAgentOverride 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setUserAgentOverride 请求超时")
		}
	}
}

/*
// === 使用场景示例：模拟iPhone手机浏览器UA ===
func ExampleSetIPhoneUserAgent() {
	// 场景：模拟iPhone访问，测试移动端页面适配
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
	resp, err := CDPNetworkSetUserAgentOverride(ua)
	if err != nil {
		log.Fatalf("设置手机UA失败: %v", err)
	}
	log.Println("✅ 已模拟iPhone浏览器UA: ", resp)
}

// === 使用场景示例：模拟Windows PC Chrome浏览器 ===
func ExampleSetPCChromeUserAgent() {
	// 场景：固定PC端UA，用于自动化测试/爬虫伪装
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"
	resp, err := CDPNetworkSetUserAgentOverride(ua)
	if err != nil {
		log.Fatalf("设置PC浏览器UA失败: %v", err)
	}
	log.Println("✅ 已模拟PC端Chrome浏览器UA: ", resp)
}

*/

// -----------------------------------------------  Network.clearAcceptedEncodingsOverride  -----------------------------------------------
// === 应用场景 ===
// 1. 编码模拟恢复: 清除自定义的Accept-Encoding覆盖，恢复浏览器默认编码设置
// 2. 测试环境清理: 自动化测试完成后清理自定义编码配置，避免影响后续用例
// 3. 网络调试重置: 调试gzip/br/deflate压缩问题后，恢复默认编码行为
// 4. 多场景切换: 切换不同网络压缩测试场景时，重置编码状态
// 5. 自动化流程清理: 自动化测试/爬虫流程结束后清理网络编码配置
// 6. 错误恢复: 编码模拟异常时快速恢复浏览器默认状态

// CDPNetworkClearAcceptedEncodingsOverride 清除Accept-Encoding编码覆盖（对应CDP方法：Network.clearAcceptedEncodingsOverride）
func CDPNetworkClearAcceptedEncodingsOverride() (string, error) {
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
		"method": "Network.clearAcceptedEncodingsOverride"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearAcceptedEncodingsOverride 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("clearAcceptedEncodingsOverride 请求超时")
		}
	}
}

/*
// === 使用场景示例：自动化测试后清除编码覆盖 ===
func ExampleClearAcceptedEncodingsForTest() {
	// 场景：压缩编码测试完成，恢复默认Accept-Encoding
	resp, err := CDPNetworkClearAcceptedEncodingsOverride()
	if err != nil {
		log.Fatalf("清除编码覆盖失败: %v", err)
	}
	log.Println("✅ 已清除Accept-Encoding覆盖，恢复默认编码: ", resp)
}

// === 使用场景示例：网络调试完成后重置编码 ===
func ExampleClearEncodingsAfterDebug() {
	// 场景：gzip/br压缩调试结束，还原浏览器默认网络编码行为
	_, err := CDPNetworkClearAcceptedEncodingsOverride()
	if err != nil {
		log.Printf("重置编码异常: %v", err)
		return
	}
	log.Println("✅ 网络编码已重置为浏览器默认状态")
}
*/

// -----------------------------------------------  Network.configureDurableMessages  -----------------------------------------------
// === 应用场景 ===
// 1. 网络日志持久化：配置网络消息持久化，长时间保留请求/响应记录
// 2. 大型页面调试：避免网络日志被自动清理，完整捕获复杂页面的网络流
// 3. 自动化故障排查：测试崩溃/异常时，保留完整网络日志用于事后分析
// 4. 性能录制：持久化存储网络请求，用于离线性能分析
// 5. 合规审计：长期保存网络交互数据，满足审计追溯需求
// 6. 长时间运行监控：持续运行的爬虫/自动化程序，持久化网络日志

// CDPNetworkConfigureDurableMessages 配置持久化网络消息（对应CDP方法：Network.configureDurableMessages）
// enabled: 是否启用持久化消息
// maxTotalSize: 最大总存储大小（字节，0为使用默认值）
func CDPNetworkConfigureDurableMessages(enabled bool, maxTotalSize int64) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := make(map[string]interface{})
	params["enabled"] = enabled
	if maxTotalSize > 0 {
		params["maxTotalSize"] = maxTotalSize
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.configureDurableMessages",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 configureDurableMessages 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("configureDurableMessages 请求超时")
		}
	}
}

/*
// === 使用场景示例：启用持久化网络消息（默认大小） ===
func ExampleEnableDurableMessages() {
	// 场景：调试复杂页面，启用网络日志持久化防止丢失
	resp, err := CDPNetworkConfigureDurableMessages(true, 0)
	if err != nil {
		log.Fatalf("启用持久化网络消息失败: %v", err)
	}
	log.Println("✅ 已启用网络消息持久化（默认大小）: ", resp)
}

// === 使用场景示例：启用并配置最大存储容量 ===
func ExampleEnableDurableMessagesWithSize() {
	// 场景：长时间监控，配置100MB持久化存储空间
	const maxSize = 100 * 1024 * 1024 // 100MB
	resp, err := CDPNetworkConfigureDurableMessages(true, maxSize)
	if err != nil {
		log.Fatalf("配置持久化网络消息失败: %v", err)
	}
	log.Println("✅ 已启用持久化网络消息，最大容量100MB: ", resp)
}

*/

// -----------------------------------------------  Network.emulateNetworkConditionsByRule  -----------------------------------------------
// === 应用场景 ===
// 1. 精细化弱网测试：按请求域名/URL规则单独限速（如仅对API接口限流）
// 2. 前端资源调试：单独模拟图片/JS/CSS资源的网络延迟，排查加载问题
// 3. 接口降级测试：针对特定后端服务模拟超时、丢包，测试容错能力
// 4. 多环境模拟：不同请求规则应用不同网络条件，更真实模拟复杂网络
// 5. 精准性能分析：只对核心业务请求限速，不影响其他资源
// 6. 自动化专项测试：针对特定请求规则做弱网/断网专项测试

// CDPNetworkEmulateNetworkConditionsByRule 按请求规则模拟网络条件（对应CDP方法：Network.emulateNetworkConditionsByRule）
// rule: 匹配规则，支持 urlPattern / domain 等匹配方式
// download: 下载速度 (bytes/s，0 不限速)
// upload: 上传速度 (bytes/s，0 不限速)
// latency: 延迟 (ms)
// packetLoss: 丢包率 (0-100)
func CDPNetworkEmulateNetworkConditionsByRule(rule map[string]interface{}, download, upload int64, latency float64, packetLoss float64) (string, error) {
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
		"rule":       rule,
		"download":   download,
		"upload":     upload,
		"latency":    latency,
		"packetLoss": packetLoss,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.emulateNetworkConditionsByRule",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 emulateNetworkConditionsByRule 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("emulateNetworkConditionsByRule 请求超时")
		}
	}
}

/*
// === 使用场景示例：对指定API接口域名进行弱网模拟 ===
func ExampleEmulateNetworkByDomain() {
	// 场景：仅对 api.example.com 接口限速，不影响页面资源
	rule := map[string]interface{}{
		"type":        "Domain",
		"value":       "api.example.com",
		"isNegative": false,
	}

	// 下载 512KB/s，上传 128KB/s，延迟 500ms，丢包 5%
	resp, err := CDPNetworkEmulateNetworkConditionsByRule(
		rule,
		512*1024,
		128*1024,
		500,
		5,
	)
	if err != nil {
		log.Fatalf("按域名模拟弱网失败: %v", err)
	}
	log.Println("✅ 已对接口域名应用弱网规则: ", resp)
}

// === 使用场景示例：对URL路径进行模糊匹配限速 ===
func ExampleEmulateNetworkByURLPattern() {
	// 场景：对所有 /api/* 接口应用高延迟模拟
	rule := map[string]interface{}{
		"type":        "UrlPattern",
		"value":       "https://*\/api/*",
		"isNegative": false,
	}

	// 无带宽限制，仅模拟 1500ms 高延迟
	resp, err := CDPNetworkEmulateNetworkConditionsByRule(
		rule,
		0,
		0,
		1500,
		0,
	)
	if err != nil {
		log.Fatalf("按URL路径模拟弱网失败: %v", err)
	}
	log.Println("✅ 已对接口路径应用高延迟规则: ", resp)
}

*/

// -----------------------------------------------  Network.enableDeviceBoundSessions  -----------------------------------------------
// === 应用场景 ===
// 1. 安全会话测试：启用设备绑定会话，测试账号安全锁定机制
// 2. 登录态防护：模拟浏览器设备绑定会话，防止会话劫持
// 3. 金融/电商测试：验证高安全场景下的设备绑定会话功能
// 4. 合规测试：满足安全合规要求，启用设备级会话绑定
// 5. 多设备登录控制：测试单设备登录、禁止跨设备共享会话
// 6. 安全调试：调试服务端设备绑定会话的校验逻辑

// CDPNetworkEnableDeviceBoundSessions 启用设备绑定会话（对应CDP方法：Network.enableDeviceBoundSessions）
func CDPNetworkEnableDeviceBoundSessions() (string, error) {
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
		"method": "Network.enableDeviceBoundSessions"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 enableDeviceBoundSessions 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("enableDeviceBoundSessions 请求超时")
		}
	}
}

/*
// === 使用场景示例：安全测试启用设备绑定会话 ===
func ExampleEnableDeviceBoundSessionsForSecurity() {
	// 场景：金融/电商安全测试，启用设备绑定会话防止劫持
	resp, err := CDPNetworkEnableDeviceBoundSessions()
	if err != nil {
		log.Fatalf("启用设备绑定会话失败: %v", err)
	}
	log.Println("✅ 已启用设备绑定会话，会话安全加固完成: ", resp)
}

// === 使用场景示例：自动化测试验证登录态防护 ===
func ExampleTestDeviceBoundSession() {
	// 场景：自动化测试单设备登录限制
	_, err := CDPNetworkEnableDeviceBoundSessions()
	if err != nil {
		log.Printf("启用设备绑定会话异常: %v", err)
		return
	}
	log.Println("✅ 已启用设备绑定，开始测试跨设备登录拦截逻辑")
}


*/

// -----------------------------------------------  Network.enableReportingApi  -----------------------------------------------
// === 应用场景 ===
// 1. 前端监控调试：启用 Reporting API 捕获网络错误、安全策略、崩溃等报告
// 2. 自动化测试：监听页面上报的诊断报告，验证监控系统是否正常工作
// 3. 安全策略调试：调试 CSP、COOP 等安全策略触发的报告
// 4. 性能监控：捕获页面性能相关的上报数据
// 5. 问题排查：线上问题复现时，收集浏览器自动生成的异常报告
// 6. 合规审计：记录浏览器上报的各类报告用于追溯分析

// CDPNetworkEnableReportingApi 启用 Reporting API（对应CDP方法：Network.enableReportingApi）
func CDPNetworkEnableReportingApi() (string, error) {
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
		"method": "Network.enableReportingApi"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 enableReportingApi 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("enableReportingApi 请求超时")
		}
	}
}

/*
// === 使用场景示例：调试时启用Reporting API捕获异常报告 ===
func ExampleEnableReportingApiForDebug() {
	// 场景：前端调试，捕获CSP、网络异常、安全策略报告
	resp, err := CDPNetworkEnableReportingApi()
	if err != nil {
		log.Fatalf("启用Reporting API失败: %v", err)
	}
	log.Println("✅ 已启用Reporting API，开始捕获浏览器报告: ", resp)
}

// === 使用场景示例：自动化测试监听页面上报报告 ===
func ExampleEnableReportingApiForTest() {
	// 场景：自动化测试，监听页面是否正常上报监控报告
	_, err := CDPNetworkEnableReportingApi()
	if err != nil {
		log.Printf("启用Reporting API异常: %v", err)
		return
	}
	log.Println("✅ 已启用Reporting API，准备接收页面上报数据")
}


*/

// -----------------------------------------------  Network.fetchSchemefulSite  -----------------------------------------------
// === 应用场景 ===
// 1. 同源策略调试：获取请求对应的同源站点（RFC 6454），排查跨域问题
// 2. 安全测试：验证站点同源策略、Cookie 隔离、存储隔离是否符合预期
// 3. 跨域权限校验：判断请求是否属于同一站点，辅助权限控制测试
// 4. 浏览器隔离机制调试：分析 Fetch / Network 站点隔离规则
// 5. 自动化测试：校验请求的同源站点是否符合预期，确保跨域配置正确
// 6. 第三方资源分析：识别请求是否属于第三方站点，用于安全审计

// CDPNetworkFetchSchemefulSite 获取请求的同源站点（scheme://host:port，对应CDP方法：Network.fetchSchemefulSite）
// requestId：网络请求ID（从 requestWillBeSent / responseReceived 事件获取）
func CDPNetworkFetchSchemefulSite(requestId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{"requestId": requestId}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.fetchSchemefulSite",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 fetchSchemefulSite 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("fetchSchemefulSite 请求超时")
		}
	}
}

/*
// === 使用场景示例：获取接口请求的同源站点 ===
func ExampleFetchSchemefulSiteByRequestId() {
	// 场景：查询某个接口请求对应的同源站点，排查跨域/Cookie问题
	requestId := "api-request-123456" // 从网络事件中获取真实ID
	resp, err := CDPNetworkFetchSchemefulSite(requestId)
	if err != nil {
		log.Fatalf("获取同源站点失败: %v", err)
	}
	log.Println("✅ 获取请求同源站点成功：\n", resp)
}

// === 使用场景示例：自动化测试校验同源策略 ===
func ExampleTestSameSiteBySchemefulSite() {
	// 场景：测试接口是否与页面同源，验证跨域配置是否正确
	requestId := "test-request-id-001"
	siteResp, err := CDPNetworkFetchSchemefulSite(requestId)
	if err != nil {
		log.Printf("获取站点信息异常: %v", err)
		return
	}

	// 判断是否为预期同源站点
	if strings.Contains(siteResp, "https://example.com") {
		log.Println("✅ 请求同源策略校验通过")
	} else {
		log.Println("⚠️ 请求属于第三方/跨域站点")
	}
}


*/

// -----------------------------------------------  Network.getCertificate  -----------------------------------------------
// === 应用场景 ===
// 1. HTTPS调试：获取请求的SSL证书详情，排查证书过期、不匹配问题
// 2. 安全审计：校验服务器证书是否合法、有效、受信任
// 3. 接口测试：验证HTTPS接口证书配置是否正确
// 4. 爬虫适配：处理自定义证书、自签名证书的校验场景
// 5. 性能监控：检查SSL握手、证书链是否正常
// 6. 合规检查：确保证书算法、有效期符合安全规范

// CDPNetworkGetCertificate 获取指定请求的SSL证书信息（对应CDP方法：Network.getCertificate）
// requestId: 网络请求ID（从 requestWillBeSent / responseReceived 事件获取）
func CDPNetworkGetCertificate(requestId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{"requestId": requestId}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.getCertificate",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getCertificate 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("getCertificate 请求超时")
		}
	}
}

/*
// === 使用场景示例：获取HTTPS接口证书详情 ===
func ExampleGetApiCertificate() {
	// 场景：调试HTTPS接口，获取证书有效期、颁发者、指纹等信息
	requestId := "https-request-id-123456" // 从网络事件获取真实ID
	resp, err := CDPNetworkGetCertificate(requestId)
	if err != nil {
		log.Fatalf("获取证书失败: %v", err)
	}
	log.Println("✅ 获取SSL证书成功：\n", resp)
}

// === 使用场景示例：自动化测试校验证书有效期 ===
func ExampleCheckCertificateValidity() {
	// 场景：测试证书是否过期，确保接口HTTPS安全
	requestId := "test-https-request-id"
	certResp, err := CDPNetworkGetCertificate(requestId)
	if err != nil {
		log.Printf("获取证书异常: %v", err)
		return
	}

	// 简单校验证书是否包含有效期字段
	if strings.Contains(certResp, "validFrom") && strings.Contains(certResp, "validTo") {
		log.Println("✅ 证书包含有效期限，配置正常")
	} else {
		log.Println("⚠️ 证书信息不完整")
	}
}

*/

// -----------------------------------------------  Network.getResponseBodyForInterception  -----------------------------------------------
// === 应用场景 ===
// 1. 请求拦截调试：获取被拦截请求的原始响应体，用于查看/修改返回数据
// 2. 接口数据篡改：抓取拦截后的响应数据，进行修改后再返回给页面
// 3. 爬虫数据提取：拦截敏感接口并直接获取响应体，无需二次解析
// 4. 自动化Mock：获取真实响应后构造Mock数据
// 5. 加密数据解密：拦截加密接口响应，获取原始密文用于调试
// 6. 故障定位：查看被拦截请求的真实服务端返回数据

// CDPNetworkGetResponseBodyForInterception 获取拦截请求的响应体（对应CDP方法：Network.getResponseBodyForInterception）
// interceptionId: 请求拦截ID（从Network.requestIntercepted事件获取）
func CDPNetworkGetResponseBodyForInterception(interceptionId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{"interceptionId": interceptionId}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.getResponseBodyForInterception",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getResponseBodyForInterception 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("getResponseBodyForInterception 请求超时")
		}
	}
}

/*
// === 使用场景示例：获取拦截接口的原始响应数据 ===
func ExampleGetInterceptedResponseBody() {
	// 场景：从requestIntercepted事件获取拦截ID，读取原始响应体
	interceptionId := "intercept-123456-api" // 从拦截事件中获取
	resp, err := CDPNetworkGetResponseBodyForInterception(interceptionId)
	if err != nil {
		log.Fatalf("获取拦截响应体失败: %v", err)
	}
	log.Println("✅ 获取拦截请求响应体成功:\n", resp)
}

// === 使用场景示例：拦截后修改数据并返回 ===
func ExampleModifyInterceptedResponse() {
	// 场景：获取拦截响应 → 修改数据 → 继续返回页面
	interceptionId := "intercept-login-api"
	respBody, err := CDPNetworkGetResponseBodyForInterception(interceptionId)
	if err != nil {
		log.Printf("获取拦截响应失败: %v", err)
		return
	}

	// 此处可对respBody进行解析、修改、加密等操作
	log.Println("✅ 获取原始响应:", respBody)
	log.Println("✅ 可修改响应数据后调用continueInterceptedRequest继续请求")
}

*/

// -----------------------------------------------  Network.getSecurityIsolationStatus  -----------------------------------------------
// === 应用场景 ===
// 1. 安全隔离调试：查询页面/Frame的COOP/COEP安全隔离状态，排查跨源隔离配置问题
// 2. 跨源通信审计：验证页面是否正确启用Cross-Origin-Opener-Policy/Embedder-Policy
// 3. 浏览器兼容性测试：检查不同站点的安全隔离标头配置是否符合规范
// 4. 安全漏洞检测：发现未正确隔离的页面，防止Spectre/Meltdown类攻击
// 5. 自动化安全校验：集成测试中校验COOP/COEP标头是否正确生效

// CDPNetworkGetSecurityIsolationStatus 获取安全隔离状态（对应CDP方法：Network.getSecurityIsolationStatus）
// frameId: 可选，页面Frame ID；不提供则返回当前目标默认Frame的状态
func CDPNetworkGetSecurityIsolationStatus(frameId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数（可选frameId）
	params := make(map[string]interface{})
	if frameId != "" {
		params["frameId"] = frameId
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.getSecurityIsolationStatus",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getSecurityIsolationStatus 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("getSecurityIsolationStatus 请求超时")
		}
	}
}

/*
// === 使用场景示例：查询当前页面主Frame安全隔离状态 ===
func ExampleCheckMainFrameSecurityIsolation() {
	// 场景：检查页面COOP/COEP是否配置正确
	status, err := CDPNetworkGetSecurityIsolationStatus("") // 空frameId查主Frame
	if err != nil {
		log.Fatalf("获取安全隔离状态失败: %v", err)
	}
	log.Println("✅ 当前页面安全隔离状态:\n", status)
}

// === 使用场景示例：校验iframe是否满足跨源隔离要求 ===
func ExampleCheckIframeSecurityIsolation() {
	// 场景：测试iframe的隔离策略，确保符合安全规范
	iframeId := "iframe-123456" // 从Page.getFrameTree获取真实frameId
	status, err := CDPNetworkGetSecurityIsolationStatus(iframeId)
	if err != nil {
		log.Printf("获取iframe隔离状态失败: %v", err)
		return
	}

	// 解析并校验COOP/COEP字段
	if strings.Contains(status, "coop") && strings.Contains(status, "coep") {
		log.Println("✅ iframe 安全隔离配置完整")
	} else {
		log.Println("⚠️ iframe 缺少COOP/COEP安全标头")
	}
}

*/

// -----------------------------------------------  Network.loadNetworkResource  -----------------------------------------------
// === 应用场景 ===
// 1. 浏览器内资源加载：让浏览器直接发起网络请求加载资源，无需页面JS触发
// 2. 静态资源预加载：提前加载JS/CSS/图片等静态资源，提升页面渲染速度
// 3. 资源完整性校验：加载资源并校验哈希值，确保文件未被篡改
// 4. 离线资源缓存：主动加载并缓存资源，用于离线可用场景
// 5. 接口直连调试：直接请求后端接口，获取原始响应数据
// 6. 自动化资源采集：批量加载网络资源用于分析、备份、测试

// CDPNetworkLoadNetworkResource 浏览器主动加载网络资源（对应CDP方法：Network.loadNetworkResource）
// url: 资源完整URL
// options: 加载选项，包含timeout、includeCredentials、disableCache等
func CDPNetworkLoadNetworkResource(url string, options map[string]interface{}) (string, error) {
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
		"url":     url,
		"options": options,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.loadNetworkResource",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 loadNetworkResource 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("loadNetworkResource 请求超时")
		}
	}
}

/*
// === 使用场景示例：基础加载网络图片资源 ===
func ExampleLoadImageResource() {
	// 场景：主动加载图片资源，带10秒超时
	url := "https://example.com/logo.png"
	options := map[string]interface{}{
		"timeout":           10.0, // 10秒超时
		"includeCredentials": false,
		"disableCache":       false,
	}

	resp, err := CDPNetworkLoadNetworkResource(url, options)
	if err != nil {
		log.Fatalf("加载图片资源失败: %v", err)
	}
	log.Println("✅ 图片资源加载成功:\n", resp)
}

// === 使用场景示例：加载JS文件并禁用缓存 ===
func ExampleLoadScriptWithNoCache() {
	// 场景：强制加载最新JS文件，不走缓存，带证书校验
	url := "https://example.com/app.js"
	options := map[string]interface{}{
		"timeout":           15.0,
		"includeCredentials": true,
		"disableCache":       true, // 禁用缓存
	}

	resp, err := CDPNetworkLoadNetworkResource(url, options)
	if err != nil {
		log.Fatalf("加载JS资源失败: %v", err)
	}
	log.Println("✅ JS资源加载成功（无缓存）:\n", resp)
}

*/

// -----------------------------------------------  Network.overrideNetworkState  -----------------------------------------------
// === 应用场景 ===
// 1. 全局弱网模拟：统一模拟网络延迟、上传/下载限速、丢包，测试页面加载容错
// 2. 离线模式测试：模拟断网状态，验证页面离线展示、异常处理逻辑
// 3. 网络环境适配：模拟3G/4G/5G/弱网等不同网络条件
// 4. 性能压测：限制带宽测试页面在低速网络下的渲染、卡顿情况
// 5. 接口超时测试：高延迟+丢包，验证接口超时重试机制
// 6. 自动化专项测试：统一网络环境，保证测试结果一致性

// CDPNetworkOverrideNetworkState 全局覆盖网络状态（对应CDP方法：Network.overrideNetworkState）
// offline: 是否模拟离线 true-离线 false-在线
// latency: 网络延迟（毫秒）
// downloadThroughput: 下载速度（bytes/s，-1不限制）
// uploadThroughput: 上传速度（bytes/s，-1不限制）
// connectionType: 连接类型（wifi/4g/3g/2g/none）
func CDPNetworkOverrideNetworkState(
	offline bool,
	latency float64,
	downloadThroughput,
	uploadThroughput int64,
	connectionType string,
) (string, error) {
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
		"offline":            offline,
		"latency":            latency,
		"downloadThroughput": downloadThroughput,
		"uploadThroughput":   uploadThroughput,
		"connectionType":     connectionType,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.overrideNetworkState",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 overrideNetworkState 请求失败: %w", err)
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
				log.Println("[DEBUG] 收到回复: ", content)

				// 解析响应
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
			return "", fmt.Errorf("overrideNetworkState 请求超时")
		}
	}
}

/*

// === 使用场景示例：模拟弱网 3G 环境（延迟+限速） ===
func ExampleOverride3GNetwork() {
	// 场景：模拟3G网络：延迟100ms、下载768KB/s、上传256KB/s
	resp, err := CDPNetworkOverrideNetworkState(
		false,       // 在线
		100,         // 延迟 100ms
		768*1024,    // 下载 768KB/s
		256*1024,    // 上传 256KB/s
		"3g",        // 连接类型
	)
	if err != nil {
		log.Fatalf("模拟3G网络失败: %v", err)
	}
	log.Println("✅ 已全局模拟 3G 弱网环境", resp)
}

// === 使用场景示例：模拟完全离线状态 ===
func ExampleOverrideOfflineNetwork() {
	// 场景：测试页面离线展示、断网异常处理
	resp, err := CDPNetworkOverrideNetworkState(
		true,   // 离线
		0,
		-1,
		-1,
		"none",
	)
	if err != nil {
		log.Fatalf("模拟离线失败: %v", err)
	}
	log.Println("✅ 已全局模拟 离线断网状态", resp)
}

*/

// -----------------------------------------------  Network.replayXHR  -----------------------------------------------
// === 应用场景 ===
// 1. 接口重放调试：重复发送XHR请求，快速复现接口问题、调试返回数据
// 2. 自动化测试：无需操作页面，直接重放已捕获的接口请求
// 3. 问题复现：线上接口异常时，重放请求验证服务端是否修复
// 4. 接口参数校验：重放请求检查参数、Header、Body是否正确
// 5. 性能压测：短时间内多次重放接口，模拟高并发请求
// 6. 调试缓存：重放请求验证接口缓存策略是否生效

// CDPNetworkReplayXHR 重放指定的XHR请求（对应CDP方法：Network.replayXHR）
// requestId: XHR请求ID（从Network.requestWillBeSent事件获取）
func CDPNetworkReplayXHR(requestId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{"requestId": requestId}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.replayXHR",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 replayXHR 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("replayXHR 请求超时")
		}
	}
}

/*
// === 使用场景示例：重放登录接口XHR请求 ===
func ExampleReplayLoginXHR() {
	// 场景：调试登录接口，无需重新登录，直接重放请求
	requestId := "xhr-login-request-123456" // 从网络事件中获取真实ID
	resp, err := CDPNetworkReplayXHR(requestId)
	if err != nil {
		log.Fatalf("重放登录接口失败: %v", err)
	}
	log.Println("✅ 登录接口XHR重放成功: ", resp)
}

// === 使用场景示例：自动化测试循环重放接口 ===
func ExampleReplayApiForTest() {
	// 场景：自动化测试，连续重放接口验证稳定性
	requestId := "xhr-api-list-request"
	// 重放3次
	for i := 0; i < 3; i++ {
		_, err := CDPNetworkReplayXHR(requestId)
		if err != nil {
			log.Printf("第%d次重放失败: %v", i+1, err)
			continue
		}
		log.Printf("✅ 第%d次接口重放成功", i+1)
	}
}

*/

// -----------------------------------------------  Network.searchInResponseBody  -----------------------------------------------
// === 应用场景 ===
// 1. 接口数据检索：在接口响应体中快速搜索关键字，无需手动解析
// 2. 问题定位：排查接口返回是否包含特定错误码、提示语、字段
// 3. 自动化校验：验证接口响应是否包含预期业务数据
// 4. 日志分析：批量检索网络响应中的关键信息
// 5. 数据提取：精准查找响应中的目标字符串，用于后续处理
// 6. 调试效率提升：避免打印大体积JSON，直接搜索定位内容

// CDPNetworkSearchInResponseBody 在指定请求的响应体中搜索关键字（对应CDP方法：Network.searchInResponseBody）
// requestId: 请求ID
// searchString: 要搜索的关键字/字符串
func CDPNetworkSearchInResponseBody(requestId, searchString string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{
		"requestId":    requestId,
		"searchString": searchString,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.searchInResponseBody",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 searchInResponseBody 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("searchInResponseBody 请求超时")
		}
	}
}

/*
// === 使用场景示例：在接口响应中搜索用户ID ===
func ExampleSearchUserIdInResponse() {
	// 场景：检查用户列表接口是否包含指定用户ID
	requestId := "api-user-list-123"
	keyword := "user_id_10001"

	resp, err := CDPNetworkSearchInResponseBody(requestId, keyword)
	if err != nil {
		log.Fatalf("响应体搜索失败: %v", err)
	}
	log.Println("✅ 搜索结果：\n", resp)
}

// === 使用场景示例：自动化测试检查错误提示 ===
func ExampleSearchErrorMsgInResponse() {
	// 场景：验证接口是否返回预期错误信息
	requestId := "api-login-fail-request"
	errorMsg := "密码错误"

	result, err := CDPNetworkSearchInResponseBody(requestId, errorMsg)
	if err != nil {
		log.Printf("搜索异常: %v", err)
		return
	}

	if strings.Contains(result, "\"line\":") {
		log.Println("✅ 找到目标错误信息，测试通过")
	} else {
		log.Println("❌ 未找到错误信息，测试失败")
	}
}

*/

// -----------------------------------------------  Network.setAcceptedEncodings  -----------------------------------------------
// === 应用场景 ===
// 1. 压缩调试：强制指定 gzip / br / deflate 等编码，测试服务端压缩适配
// 2. 接口调试：只允许明文（identity），方便查看未压缩原始响应体
// 3. 弱网优化：强制开启压缩，验证传输体积优化效果
// 4. 兼容性测试：模拟不支持 br 压缩的旧浏览器，验证降级逻辑
// 5. 性能测试：对比不同压缩算法的传输速度、体积差异
// 6. 爬虫适配：精准控制 Accept-Encoding，避免乱码、解析失败

// CDPNetworkSetAcceptedEncodings 设置接受的压缩编码（对应CDP方法：Network.setAcceptedEncodings）
// encodings: 编码列表，例如 []string{"gzip", "br", "deflate", "identity"}
func CDPNetworkSetAcceptedEncodings(encodings []string) (string, error) {
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
		"encodings": encodings,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setAcceptedEncodings",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setAcceptedEncodings 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setAcceptedEncodings 请求超时")
		}
	}
}

/*
// === 使用场景示例：只启用 gzip + br 压缩 ===
func ExampleSetGzipBrEncodings() {
	// 场景：标准现代浏览器编码配置
	encodings := []string{"gzip", "br"}
	resp, err := CDPNetworkSetAcceptedEncodings(encodings)
	if err != nil {
		log.Fatalf("设置编码失败: %v", err)
	}
	log.Println("✅ 已设置接受编码: gzip, br", resp)
}

// === 使用场景示例：禁用所有压缩（纯文本调试） ===
func ExampleSetIdentityOnlyEncoding() {
	// 场景：查看原始明文响应体，不接受任何压缩
	encodings := []string{"identity"}
	resp, err := CDPNetworkSetAcceptedEncodings(encodings)
	if err != nil {
		log.Fatalf("设置无压缩编码失败: %v", err)
	}
	log.Println("✅ 已禁用压缩，仅接受明文", resp)
}

*/

// -----------------------------------------------  Network.setAttachDebugStack  -----------------------------------------------
// === 应用场景 ===
// 1. 网络请求溯源：给每个网络请求附加调用栈信息，快速定位代码发起位置
// 2. 前端问题排查：定位哪个组件/函数触发了异常/重复接口请求
// 3. 性能分析：查看请求发起链路，优化冗余请求
// 4. 自动化调试：精准追踪请求来源，提升调试效率
// 5. 第三方库分析：识别SDK/插件发起的网络请求
// 6. 代码规范检查：检测不规范的请求调用位置

// CDPNetworkSetAttachDebugStack 设置是否为网络请求附加调试调用栈
// enabled: true-开启调用栈追踪  false-关闭
func CDPNetworkSetAttachDebugStack(enabled bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]bool{"enabled": enabled}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setAttachDebugStack",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setAttachDebugStack 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setAttachDebugStack 请求超时")
		}
	}
}

/*
// === 使用场景示例：开启网络请求调用栈追踪 ===
func ExampleEnableAttachDebugStack() {
	// 场景：调试接口来源，开启调用栈追踪
	resp, err := CDPNetworkSetAttachDebugStack(true)
	if err != nil {
		log.Fatalf("开启调用栈失败: %v", err)
	}
	log.Println("✅ 已开启网络请求调试调用栈，可查看请求发起代码位置", resp)
}

// === 使用场景示例：关闭调用栈追踪（提升性能） ===
func ExampleDisableAttachDebugStack() {
	// 场景：调试完成，关闭调用栈减少性能消耗
	resp, err := CDPNetworkSetAttachDebugStack(false)
	if err != nil {
		log.Fatalf("关闭调用栈失败: %v", err)
	}
	log.Println("✅ 已关闭网络请求调试调用栈，恢复正常性能", resp)
}

*/

// -----------------------------------------------  Network.setBlockedURLs  -----------------------------------------------
// === 应用场景 ===
// 1. 广告/统计屏蔽：屏蔽广告、埋点、统计接口，加速页面调试
// 2. 前端异常测试：屏蔽指定资源，测试页面降级、容错、占位逻辑
// 3. 接口隔离：屏蔽第三方接口，避免影响核心功能测试
// 4. 性能调试：屏蔽大文件、图片、非核心请求，专注核心业务
// 5. 安全测试：屏蔽敏感接口，验证权限与拦截有效性
// 6. 自动化用例：稳定测试环境，屏蔽随机/第三方波动资源

// CDPNetworkSetBlockedURLs 屏蔽指定URL规则的网络请求（对应CDP方法：Network.setBlockedURLs）
// urls: 屏蔽规则列表，支持 * 通配符
func CDPNetworkSetBlockedURLs(urls []string) (string, error) {
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
		"urls": urls,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setBlockedURLs",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setBlockedURLs 请求失败: %w", err)
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

				// 解析响应
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
			return "", fmt.Errorf("setBlockedURLs 请求超时")
		}
	}
}

/*

// === 使用场景示例：屏蔽广告、统计、图片资源 ===
func ExampleBlockCommonResources() {
	// 支持 * 通配符
	urls := []string{
		"*://*.baidu.com/*",
		"*://*.google-analytics.com/*",
		"*.png*",
		"*.jpg*",
		"*\/ad/*",
	}

	resp, err := CDPNetworkSetBlockedURLs(urls)
	if err != nil {
		log.Fatalf("设置屏蔽URL失败: %v", err)
	}
	log.Println("✅ 已屏蔽广告/统计/图片资源", resp)
}

// === 使用场景示例：清空屏蔽规则（恢复所有请求） ===
func ExampleClearBlockedURLs() {
	// 传空数组即可清空屏蔽
	resp, err := CDPNetworkSetBlockedURLs([]string{})
	if err != nil {
		log.Fatalf("清空屏蔽规则失败: %v", err)
	}
	log.Println("✅ 已清空所有URL屏蔽规则", resp)
}

*/

// -----------------------------------------------  Network.setCookieControls  -----------------------------------------------
// === 应用场景 ===
// 1. 全局禁用Cookie写入：只允许读/发Cookie，禁止任何网站设置/修改Cookie
// 2. 测试隐私模式：验证页面在无法写入Cookie时的行为（如登录态、会话保持）
// 3. 自动化防污染：防止测试过程中Cookie被意外修改，保证用例稳定
// 4. 广告/追踪屏蔽：阻止第三方SDK设置追踪Cookie，保护隐私

// CDPNetworkSetCookieControls 控制浏览器是否允许写入Cookie（实验性API）
// enableWrite: true-允许写入（默认） false-禁止写入（冻结当前Cookie）
func CDPNetworkSetCookieControls(enableWrite bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数：实验性字段，控制cookie写入
	params := map[string]bool{
		"allowWrite": enableWrite, // 核心：是否允许写入Cookie
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.setCookieControls",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setCookieControls 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 5秒超时等待响应
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
			return "", fmt.Errorf("setCookieControls 请求超时")
		}
	}
}

/*
// === 场景1：冻结Cookie（禁止写入） ===
func ExampleFreezeCookies() {
	// 效果：页面只能读取/发送已有Cookie，无法set/modify/delete
	resp, err := CDPNetworkSetCookieControls(false)
	if err != nil {
		log.Fatalf("禁用Cookie写入失败: %v", err)
	}
	log.Println("✅ 已冻结Cookie：仅允许读取/发送，禁止任何写入/修改", resp)
}

// === 场景2：恢复正常（允许写入） ===
func ExampleAllowCookieWrite() {
	resp, err := CDPNetworkSetCookieControls(true)
	if err != nil {
		log.Fatalf("启用Cookie写入失败: %v", err)
	}
	log.Println("✅ 已恢复Cookie正常读写：可自由设置/修改/删除", resp)
}

*/

// -----------------------------------------------  Network.streamResourceContent  -----------------------------------------------
// === 应用场景 ===
// 1. 大资源分段获取：获取图片/JS/CSS/接口等大体积资源，避免一次性加载卡顿
// 2. 流式读取响应：实时接收响应数据，适用于视频、日志、大文件导出
// 3. 调试原始数据：获取未修改、未解码的原生响应流
// 4. 资源抓取保存：分段接收并保存网络资源到本地
// 5. 避免内存溢出：处理GB级别的响应体，防止占用过大内存

// CDPNetworkStreamResourceContent 流式获取网络资源内容（对应CDP方法：Network.streamResourceContent）
// requestId: 要获取的网络请求ID
func CDPNetworkStreamResourceContent(requestId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{
		"requestId": requestId,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.streamResourceContent",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 streamResourceContent 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
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
			return "", fmt.Errorf("streamResourceContent 请求超时")
		}
	}
}

/*
// === 使用场景示例：流式获取大图片/大JS资源 ===
func ExampleStreamImageResource() {
	// 场景：获取大体积资源，避免内存占用过高
	requestId := "resource-large-image-123"
	resp, err := CDPNetworkStreamResourceContent(requestId)
	if err != nil {
		log.Fatalf("流式获取资源失败: %v", err)
	}
	log.Println("✅ 已启动资源流式获取:", resp)
}

// === 使用场景示例：流式获取接口响应并分段处理 ===
func ExampleStreamApiResponse() {
	// 场景：接口返回大量数据，流式接收避免解析崩溃
	requestId := "api-export-large-data-request"
	resp, err := CDPNetworkStreamResourceContent(requestId)
	if err != nil {
		log.Printf("流式获取接口响应失败: %v", err)
		return
	}
	log.Println("✅ 接口响应流式接收中:", resp)
}

*/

// -----------------------------------------------  Network.takeResponseBodyForInterceptionAsStream  -----------------------------------------------
// === 应用场景 ===
// 1. 拦截大资源流式处理：对请求拦截后的超大响应体进行流式读取，避免内存溢出
// 2. 视频/图片拦截修改：边接收拦截流边处理，无需等待完整下载
// 3. 接口加密解密：流式接收加密响应，实时解密，不占用大量内存
// 4. 数据抓取与保存：分段接收拦截响应，直接写入文件
// 5. 避免卡顿：处理GB级别的响应数据，保持程序稳定

// CDPNetworkTakeResponseBodyForInterceptionAsStream 将拦截的响应体转为流获取（对应CDP方法：Network.takeResponseBodyForInterceptionAsStream）
// interceptionId: 请求拦截ID（从Network.requestIntercepted事件获取）
func CDPNetworkTakeResponseBodyForInterceptionAsStream(interceptionId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	params := map[string]string{
		"interceptionId": interceptionId,
	}

	// 参数序列化
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("参数序列化失败: %w", err)
	}

	// 构建CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Network.takeResponseBodyForInterceptionAsStream",
		"params": %s
	}`, reqID, string(paramsBytes))

	// 发送WebSocket消息
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 takeResponseBodyForInterceptionAsStream 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
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
			return "", fmt.Errorf("takeResponseBodyForInterceptionAsStream 请求超时")
		}
	}
}

/*
// === 使用场景示例：流式获取拦截的大文件/视频响应 ===
func ExampleStreamInterceptedVideo() {
	// 场景：拦截视频/大文件下载，转为流读取，不占内存
	interceptionId := "intercept-video-stream-123456"
	resp, err := CDPNetworkTakeResponseBodyForInterceptionAsStream(interceptionId)
	if err != nil {
		log.Fatalf("流式获取拦截响应失败: %v", err)
	}
	log.Println("✅ 拦截响应已转为流式，开始分段接收数据:", resp)
}

// === 使用场景示例：流式获取拦截的加密接口响应 ===
func ExampleStreamInterceptedEncryptedApi() {
	// 场景：拦截加密接口，流式接收实时解密
	interceptionId := "intercept-encrypted-api"
	resp, err := CDPNetworkTakeResponseBodyForInterceptionAsStream(interceptionId)
	if err != nil {
		log.Printf("流式获取拦截响应失败: %v", err)
		return
	}
	log.Println("✅ 加密接口响应流式接收中:", resp)
}
*/
