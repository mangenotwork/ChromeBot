package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Storage.clearCookies  -----------------------------------------------
// === 应用场景 ===
// 1. 登录态清除: 清除浏览器Cookie，强制用户重新登录
// 2. 自动化测试: 测试用例执行前清空Cookie，保证测试环境干净
// 3. 隐私保护: 退出账号时清除所有会话Cookie
// 4. 调试缓存: 调试Cookie相关问题时快速清空所有Cookie
// 5. 多账号测试: 切换测试账号前清除原有Cookie避免冲突
// 6. 页面重置: 重置网站状态，清除所有Cookie相关的用户配置

// CDPStorageClearCookies 清除浏览器所有Cookie
func CDPStorageClearCookies() (string, error) {
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
		"method": "Storage.clearCookies"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearCookies 请求失败: %w", err)
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
			return "", fmt.Errorf("clearCookies 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：自动化测试前清空Cookie，保证测试环境纯净
func ExampleTestBeforeClearCookies() {
	resp, err := CDPStorageClearCookies()
	if err != nil {
		log.Fatalf("清空Cookie失败: %v", err)
	}
	log.Printf("清空Cookie成功，响应: %s", resp)
}

// 场景2：用户退出登录时清除所有Cookie
func ExampleUserLogoutClearCookies() {
	resp, err := CDPStorageClearCookies()
	if err != nil {
		log.Printf("退出时清除Cookie失败: %v", err)
		return
	}
	log.Println("已清除所有Cookie，用户安全退出")
}

*/

// -----------------------------------------------  Storage.clearDataForOrigin  -----------------------------------------------
// === 应用场景 ===
// 1. 单站点数据清空: 只清除指定域名的存储数据，不影响其他网站
// 2. 自动化测试: 针对测试域名清理缓存、LocalStorage、Cookie等数据
// 3. 隐私清理: 退出指定网站时清除该站点的所有本地存储
// 4. 调试存储问题: 精准清理指定域名的存储数据排查问题
// 5. 账号切换: 切换同一网站不同账号前清空该域名所有存储数据
// 6. 页面强制重置: 重置指定网站的本地状态，等同于清除站点数据

// CDPStorageClearDataForOrigin 清除指定源(域名)的存储数据
// origin: 要清除数据的源地址，例如 "https://www.example.com"
// storageTypes: 要清除的存储类型，例如 "all" 清除所有，"cookies,local_storage" 清除指定类型
func CDPStorageClearDataForOrigin(origin string, storageTypes string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息，带参数
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Storage.clearDataForOrigin",
		"params": {
			"origin": "%s",
			"storageTypes": "%s"
		}
	}`, reqID, origin, storageTypes)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearDataForOrigin 请求失败: %w", err)
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
			return "", fmt.Errorf("clearDataForOrigin 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：清除指定网站所有存储数据（最常用）
func ExampleClearAllDataForOrigin() {
	origin := "https://github.com"
	resp, err := CDPStorageClearDataForOrigin(origin, "all")
	if err != nil {
		log.Fatalf("清除 %s 所有数据失败: %v", origin, err)
	}
	log.Printf("成功清除 %s 所有存储数据，响应: %s", origin, resp)
}

// 场景2：只清除指定网站的Cookie和LocalStorage
func ExampleClearSpecificStorageForOrigin() {
	origin := "https://www.google.com"
	// 可选类型：cookies, local_storage, session_storage, web_sql, indexed_db, cache_storage, all
	storage := "cookies,local_storage"
	resp, err := CDPStorageClearDataForOrigin(origin, storage)
	if err != nil {
		log.Printf("清除 %s 数据失败: %v", origin, err)
		return
	}
	log.Printf("成功清除 %s 的Cookie和本地存储", origin)
}

// 场景3：自动化测试前清理测试站点数据
func ExampleTestClearOriginData() {
	testOrigin := "https://test.example.com"
	resp, err := CDPStorageClearDataForOrigin(testOrigin, "all")
	if err != nil {
		log.Fatalf("测试环境清理失败: %v", err)
	}
	log.Println("测试站点数据已清空，可以开始执行测试用例")
}

*/

// -----------------------------------------------  Storage.clearDataForStorageKey  -----------------------------------------------
// === 应用场景 ===
// 1. 精准存储数据清除：基于StorageKey清除数据，支持分区存储场景下的精准清理
// 2. 第三方站点数据清理：清除第三方嵌入站点的存储数据，不影响主站
// 3. 自动化测试：针对带存储分区的测试站点进行数据隔离清理
// 4. 隐私数据保护：精准删除特定存储分区下的用户数据
// 5. 调试存储问题：定位分区存储异常时，精准清除对应存储键数据
// 6. 多租户/多实例隔离：清除不同存储分区的独立数据，互不干扰

// CDPStorageClearDataForStorageKey 清除指定StorageKey的存储数据
// storageKey: 存储键，用于标识分区存储的唯一键
// storageTypes: 要清除的存储类型，例如 "all" 清除所有，"cookies,local_storage" 清除指定类型
func CDPStorageClearDataForStorageKey(storageKey string, storageTypes string) (string, error) {
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
		"method": "Storage.clearDataForStorageKey",
		"params": {
			"storageKey": "%s",
			"storageTypes": "%s"
		}
	}`, reqID, storageKey, storageTypes)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearDataForStorageKey 请求失败: %w", err)
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
			return "", fmt.Errorf("clearDataForStorageKey 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：清除指定StorageKey的所有存储数据
func ExampleClearAllDataForStorageKey() {
	// 示例StorageKey（实际使用时通过CDP获取真实存储键）
	storageKey := "https://example.com^180daypartition"
	resp, err := CDPStorageClearDataForStorageKey(storageKey, "all")
	if err != nil {
		log.Fatalf("清除 StorageKey 所有数据失败: %v", err)
	}
	log.Printf("成功清除 StorageKey[%s] 所有存储数据", storageKey)
}

// 场景2：仅清除指定StorageKey的Cookie和会话存储
func ExampleClearSpecificStorageForKey() {
	storageKey := "https://sub.example.com^partition"
	storageTypes := "cookies,session_storage"
	resp, err := CDPStorageClearDataForStorageKey(storageKey, storageTypes)
	if err != nil {
		log.Printf("清除指定存储数据失败: %v", err)
		return
	}
	log.Printf("成功清理 StorageKey 分区数据")
}

// 场景3：自动化测试 - 分区存储环境数据重置
func ExampleTestResetPartitionStorage() {
	testStorageKey := "https://test.app^testpartition"
	resp, err := CDPStorageClearDataForStorageKey(testStorageKey, "all")
	if err != nil {
		log.Fatalf("测试分区存储重置失败: %v", err)
	}
	log.Println("测试存储分区已清空，可执行隔离测试")
}

*/

// -----------------------------------------------  Storage.getCookies  -----------------------------------------------
// === 应用场景 ===
// 1. Cookie获取: 获取浏览器当前所有Cookie，用于分析登录态、会话信息
// 2. 自动化测试: 验证Cookie是否正确设置、过期、携带指定字段
// 3. 调试分析: 排查Cookie失效、跨域、过期时间等问题
// 4. 数据采集: 采集会话Cookie用于接口请求、爬虫鉴权
// 5. 登录校验: 检查关键登录Cookie是否存在，判断登录状态
// 6. 安全审计: 审计Cookie的Secure、HttpOnly、SameSite等安全属性

// CDPStorageGetCookies 获取当前浏览器所有Cookie
func CDPStorageGetCookies() (string, error) {
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
		"method": "Storage.getCookies"
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


// === 使用场景示例代码 ===
// 场景1：获取所有Cookie并打印完整信息
func ExampleGetAllCookies() {
	cookiesResp, err := CDPStorageGetCookies()
	if err != nil {
		log.Fatalf("获取Cookie失败: %v", err)
	}
	log.Printf("当前所有Cookie信息:\n%s", cookiesResp)
}

// 场景2：自动化测试 - 验证登录Cookie是否存在
func ExampleTestCheckLoginCookie() {
	cookiesResp, err := CDPStorageGetCookies()
	if err != nil {
		log.Printf("获取Cookie失败: %v", err)
		return
	}

	// 检查是否包含登录态Cookie（示例：token）
	if strings.Contains(cookiesResp, "token") {
		log.Println("✅ 登录Cookie存在，用户已登录")
	} else {
		log.Println("❌ 登录Cookie不存在，用户未登录")
	}
}

// 场景3：采集Cookie用于接口请求鉴权
func ExampleCollectCookieForRequest() {
	cookiesResp, err := CDPStorageGetCookies()
	if err != nil {
		log.Fatalf("采集Cookie失败: %v", err)
	}
	log.Println("成功采集浏览器Cookie，可用于接口请求头")
	// 后续可解析cookiesResp拼接成Cookie字符串
}

*/

// -----------------------------------------------  Storage.getUsageAndQuota  -----------------------------------------------
// === 应用场景 ===
// 1. 存储容量监控: 查询浏览器为当前域名分配的存储配额和已使用空间
// 2. 自动化测试: 验证站点存储是否超出限制、排查存储溢出问题
// 3. 前端调试: 定位 LocalStorage/IndexedDB 存储爆满导致的功能异常
// 4. 存储优化: 分析各类型存储占用，针对性清理释放空间
// 5. 用户提示: 监控存储使用率，接近上限时提示用户清理数据
// 6. 多站点对比: 对比不同域名的存储占用与配额分配情况

// CDPStorageGetUsageAndQuota 查询指定源的存储使用情况与配额限制
// origin: 要查询的站点源地址，例如 "https://www.example.com"
func CDPStorageGetUsageAndQuota(origin string) (string, error) {
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
		"method": "Storage.getUsageAndQuota",
		"params": {
			"origin": "%s"
		}
	}`, reqID, origin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getUsageAndQuota 请求失败: %w", err)
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
			return "", fmt.Errorf("getUsageAndQuota 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：查询指定站点的总存储配额和已使用空间
func ExampleGetStorageQuota() {
	origin := "https://github.com"
	resp, err := CDPStorageGetUsageAndQuota(origin)
	if err != nil {
		log.Fatalf("查询 %s 存储配额失败: %v", origin, err)
	}
	log.Printf("=== %s 存储信息 ===", origin)
	log.Printf("配额与使用详情：\n%s", resp)
}

// 场景2：自动化测试 - 校验站点存储是否超过警戒阈值
func ExampleTestCheckStorageUsage() {
	origin := "https://test.example.com"
	resp, err := CDPStorageGetUsageAndQuota(origin)
	if err != nil {
		log.Printf("获取存储信息失败: %v", err)
		return
	}

	// 解析判断存储使用率（示例：超过80%告警）
	if strings.Contains(resp, "\"usage\":") && strings.Contains(resp, "\"quota\":") {
		log.Println("✅ 成功获取存储数据，可进行使用率计算")
		// 可自行解析JSON提取 usage 和 quota 数值计算百分比
	} else {
		log.Println("❌ 未获取到有效存储数据")
	}
}

// 场景3：调试存储溢出问题
func ExampleDebugStorageFull() {
	origin := "https://app.example.com"
	resp, err := CDPStorageGetUsageAndQuota(origin)
	if err != nil {
		log.Fatalf("调试查询失败: %v", err)
	}
	log.Println("存储调试信息：", resp)
	// 返回结果包含 usage(已使用)、quota(总配额)、各类型存储细分占用
}

*/

// -----------------------------------------------  Storage.setCookies  -----------------------------------------------
// === 应用场景 ===
// 1. 登录态注入: 提前设置登录Cookie，实现免登录直接进入目标页面
// 2. 自动化测试: 自定义测试Cookie，模拟不同用户状态、权限场景
// 3. 会话保持: 注入有效会话Cookie，维持浏览器登录状态
// 4. 调试Cookie: 手动设置Cookie属性，调试Secure、HttpOnly、SameSite等特性
// 5. 多账号切换: 快速设置不同账号的Cookie，实现账号无缝切换
// 6. 接口鉴权: 设置必要Cookie，保证页面请求接口能正常鉴权通过

// CDPStorageSetCookies 设置浏览器Cookie（支持批量/单个）
// cookiesJson: Cookie数组的JSON字符串，格式参考CDP规范，包含name、value、domain、path等字段
func CDPStorageSetCookies(cookiesJson string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息，直接传入拼接好的cookies数组
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "Storage.setCookies",
		"params": {
			"cookies": %s
		}
	}`, reqID, cookiesJson)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
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


// === 使用场景示例代码 ===
// 场景1：设置单个登录Cookie（最常用：免登录访问）
func ExampleSetSingleLoginCookie() {
	// 单个Cookie标准JSON格式
	cookie := `[{
		"name": "session_id",
		"value": "abc123def456",
		"domain": ".example.com",
		"path": "/",
		"secure": true,
		"httpOnly": true,
		"sameSite": "Lax"
	}]`

	resp, err := CDPStorageSetCookies(cookie)
	if err != nil {
		log.Fatalf("设置登录Cookie失败: %v", err)
	}
	log.Println("✅ 成功注入登录Cookie，可直接访问目标页面")
}

// 场景2：批量设置多个Cookie（完整会话状态）
func ExampleSetMultipleCookies() {
	// 批量Cookie数组
	cookies := `[{
		"name": "user_id",
		"value": "10001",
		"domain": "app.example.com",
		"path": "/"
	}, {
		"name": "theme",
		"value": "dark",
		"domain": "app.example.com",
		"path": "/"
	}]`

	resp, err := CDPStorageSetCookies(cookies)
	if err != nil {
		log.Printf("批量设置Cookie失败: %v", err)
		return
	}
	log.Println("✅ 批量Cookie设置完成，用户状态已配置")
}

// 场景3：自动化测试 - 设置测试专用Cookie
func ExampleTestSetCustomCookie() {
	testCookie := `[{
		"name": "test_env",
		"value": "uat",
		"domain": "test.example.com",
		"secure": false
	}]`

	resp, err := CDPStorageSetCookies(testCookie)
	if err != nil {
		log.Fatalf("测试Cookie设置失败: %v", err)
	}
	log.Println("✅ 测试环境Cookie已设置，开始执行测试用例")
}

*/

// -----------------------------------------------  Storage.setProtectedAudienceKAnonymity  -----------------------------------------------
// === 应用场景 ===
// 1. 广告测试环境：强制设置 Protected Audience API 的 K-Anonymity 状态，用于测试广告投放逻辑
// 2. 隐私沙箱调试：模拟 K 匿名检查结果，验证广告竞价、过滤逻辑
// 3. 自动化测试：绕过真实 K 匿名计算，稳定复现广告相关测试用例
// 4. 开发环境配置：本地开发时强制指定广告受众的匿名状态
// 5. 兼容性验证：测试不同 K-Anonymity 值下的广告行为表现
// 6. 合规测试：验证广告系统在隐私合规约束下的运行逻辑

// CDPStorageSetProtectedAudienceKAnonymity 强制设置 Protected Audience API 的 K-Anonymity 状态
// name: 标识符（如广告 buyer/bidding/signals）
// value: K-Anonymity 结果（bool 类型）
func CDPStorageSetProtectedAudienceKAnonymity(name string, value bool) (string, error) {
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
		"method": "Storage.setProtectedAudienceKAnonymity",
		"params": {
			"name": "%s",
			"value": %t
		}
	}`, reqID, name, value)

	// 发送 WebSocket 消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setProtectedAudienceKAnonymity 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应超时控制
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 循环等待对应 ID 的响应
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配响应 ID
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
			return "", fmt.Errorf("setProtectedAudienceKAnonymity 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：广告测试 - 强制开启 Protected Audience K-Anonymity 允许状态
func ExampleSetKAnonymityAllowed() {
	// 强制标记为满足 K 匿名要求
	name := "protected_audience_k_anonymity"
	value := true
	resp, err := CDPStorageSetProtectedAudienceKAnonymity(name, value)
	if err != nil {
		log.Fatalf("设置 K-Anonymity 允许失败: %v", err)
	}
	log.Println("✅ 已强制启用 K-Anonymity 允许状态，广告可正常参与竞价")
}

// 场景2：隐私测试 - 强制禁用 K-Anonymity，验证广告过滤逻辑
func ExampleSetKAnonymityDenied() {
	// 强制标记为不满足匿名要求
	name := "protected_audience_k_anonymity"
	value := false
	resp, err := CDPStorageSetProtectedAudienceKAnonymity(name, value)
	if err != nil {
		log.Fatalf("设置 K-Anonymity 禁用失败: %v", err)
	}
	log.Println("✅ 已强制禁用 K-Anonymity，广告将被过滤不参与竞价")
}

// 场景3：自动化广告测试用例前置配置
func ExampleTestSetupKAnonymity() {
	// 测试环境固定设置为允许，保证测试用例稳定执行
	resp, err := CDPStorageSetProtectedAudienceKAnonymity("bidding_signal", true)
	if err != nil {
		log.Fatalf("测试前置配置 K-Anonymity 失败: %v", err)
	}
	log.Println("✅ 广告测试环境 K-Anonymity 配置完成")
}

*/

// -----------------------------------------------  Storage.trackCacheStorageForOrigin  -----------------------------------------------
// === 应用场景 ===
// 1. 缓存监控：开启对指定源 CacheStorage 的实时跟踪，监听缓存创建/更新/删除
// 2. PWA 调试：调试渐进式 Web 应用的离线缓存策略，监控缓存变更
// 3. 性能分析：跟踪静态资源缓存行为，分析缓存命中与加载性能
// 4. 自动化测试：监听测试站点缓存操作，验证缓存逻辑是否符合预期
// 5. 资源管理：监控第三方站点 CacheStorage 占用与变更
// 6. 问题排查：定位缓存未更新、缓存污染、离线加载异常等问题

// CDPStorageTrackCacheStorageForOrigin 开启跟踪指定源的 CacheStorage 存储
// origin: 需要跟踪缓存的站点源地址，例如 "https://www.example.com"
func CDPStorageTrackCacheStorageForOrigin(origin string) (string, error) {
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
		"method": "Storage.trackCacheStorageForOrigin",
		"params": {
			"origin": "%s"
		}
	}`, reqID, origin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 trackCacheStorageForOrigin 请求失败: %w", err)
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
			return "", fmt.Errorf("trackCacheStorageForOrigin 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启PWA站点的CacheStorage跟踪，调试离线缓存
func ExampleTrackPwaCacheStorage() {
	origin := "https://pwa.example.com"
	resp, err := CDPStorageTrackCacheStorageForOrigin(origin)
	if err != nil {
		log.Fatalf("开启PWA缓存跟踪失败: %v", err)
	}
	log.Printf("✅ 成功开启 [%s] CacheStorage 实时跟踪，可监听缓存变更事件", origin)
}

// 场景2：自动化测试 - 前置开启缓存跟踪，验证缓存写入逻辑
func ExampleTestSetupCacheTracking() {
	testOrigin := "https://test.app.com"
	resp, err := CDPStorageTrackCacheStorageForOrigin(testOrigin)
	if err != nil {
		log.Fatalf("测试缓存跟踪开启失败: %v", err)
	}
	log.Println("✅ 测试站点缓存跟踪已开启，准备执行缓存相关测试用例")
}

// 场景3：调试生产站点静态资源缓存异常
func ExampleDebugProductionCache() {
	origin := "https://www.example.com"
	resp, err := CDPStorageTrackCacheStorageForOrigin(origin)
	if err != nil {
		log.Printf("开启缓存跟踪失败: %v", err)
		return
	}
	log.Println("✅ 已开启生产环境缓存跟踪，可实时监控缓存更新/删除行为")
}

*/

// -----------------------------------------------  Storage.trackCacheStorageForStorageKey  -----------------------------------------------
// === 应用场景 ===
// 1. 分区缓存监控：基于StorageKey开启分区存储下的CacheStorage跟踪，适配浏览器存储隔离
// 2. 第三方站点缓存调试：跟踪第三方嵌入站点的独立缓存分区，不影响主站缓存
// 3. PWA多实例监控：监控同一域名下不同分区、不同实例的PWA离线缓存
// 4. 自动化测试：针对带存储分区的测试环境做缓存行为隔离监控
// 5. 缓存问题定位：精准定位分区存储下的缓存未更新、污染、加载异常
// 6. 多租户缓存管理：管理不同存储分区下的独立缓存，互不干扰

// CDPStorageTrackCacheStorageForStorageKey 开启跟踪指定StorageKey的CacheStorage存储
// storageKey: 存储键，用于标识分区存储的唯一键
func CDPStorageTrackCacheStorageForStorageKey(storageKey string) (string, error) {
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
		"method": "Storage.trackCacheStorageForStorageKey",
		"params": {
			"storageKey": "%s"
		}
	}`, reqID, storageKey)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 trackCacheStorageForStorageKey 请求失败: %w", err)
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
			return "", fmt.Errorf("trackCacheStorageForStorageKey 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启指定StorageKey分区缓存跟踪，调试分区存储PWA
func ExampleTrackPartitionCacheByStorageKey() {
	storageKey := "https://pwa.example.com^partition1"
	resp, err := CDPStorageTrackCacheStorageForStorageKey(storageKey)
	if err != nil {
		log.Fatalf("开启分区缓存跟踪失败: %v", err)
	}
	log.Printf("✅ 成功开启 StorageKey[%s] CacheStorage 跟踪", storageKey)
}

// 场景2：自动化测试 - 分区缓存环境监控
func ExampleTestSetupPartitionCacheTracking() {
	testStorageKey := "https://test.app^test-partition"
	resp, err := CDPStorageTrackCacheStorageForStorageKey(testStorageKey)
	if err != nil {
		log.Fatalf("测试分区缓存监控开启失败: %v", err)
	}
	log.Println("✅ 测试分区缓存跟踪已启动，可监控缓存变更事件")
}

// 场景3：第三方站点分区缓存调试
func ExampleDebugThirdPartyCache() {
	storageKey := "https://third-party.com^embed-partition"
	resp, err := CDPStorageTrackCacheStorageForStorageKey(storageKey)
	if err != nil {
		log.Printf("第三方缓存跟踪失败: %v", err)
		return
	}
	log.Println("✅ 第三方分区缓存跟踪已开启，可排查离线加载问题")
}

*/

// -----------------------------------------------  Storage.trackIndexedDBForOrigin  -----------------------------------------------
// === 应用场景 ===
// 1. IndexedDB 实时监控：开启对指定源 IndexedDB 数据库的跟踪，监听创建/更新/删除操作
// 2. PWA 离线数据调试：调试渐进式 Web 应用的本地数据库存储行为
// 3. 前端数据存储分析：跟踪 IndexedDB 数据读写逻辑，分析数据持久化流程
// 4. 自动化测试：监听测试站点 IndexedDB 操作，验证数据存储逻辑是否符合预期
// 5. 数据问题排查：定位 IndexedDB 数据丢失、写入失败、查询异常等问题
// 6. 存储性能监控：跟踪 IndexedDB 操作耗时与数据量，优化存储性能

// CDPStorageTrackIndexedDBForOrigin 开启跟踪指定源的 IndexedDB 存储
// origin: 需要跟踪 IndexedDB 的站点源地址，例如 "https://www.example.com"
func CDPStorageTrackIndexedDBForOrigin(origin string) (string, error) {
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
		"method": "Storage.trackIndexedDBForOrigin",
		"params": {
			"origin": "%s"
		}
	}`, reqID, origin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 trackIndexedDBForOrigin 请求失败: %w", err)
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
			return "", fmt.Errorf("trackIndexedDBForOrigin 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启PWA站点的IndexedDB跟踪，调试离线数据存储
func ExampleTrackPwaIndexedDB() {
	origin := "https://pwa.example.com"
	resp, err := CDPStorageTrackIndexedDBForOrigin(origin)
	if err != nil {
		log.Fatalf("开启PWA IndexedDB跟踪失败: %v", err)
	}
	log.Printf("✅ 成功开启 [%s] IndexedDB 实时跟踪，可监听数据库变更事件", origin)
}

// 场景2：自动化测试 - 前置开启数据库跟踪，验证数据写入逻辑
func ExampleTestSetupIndexedDBTracking() {
	testOrigin := "https://test.app.com"
	resp, err := CDPStorageTrackIndexedDBForOrigin(testOrigin)
	if err != nil {
		log.Fatalf("测试数据库跟踪开启失败: %v", err)
	}
	log.Println("✅ 测试站点数据库跟踪已开启，准备执行数据存储测试用例")
}

// 场景3：调试生产站点IndexedDB数据异常问题
func ExampleDebugProductionIndexedDB() {
	origin := "https://www.example.com"
	resp, err := CDPStorageTrackIndexedDBForOrigin(origin)
	if err != nil {
		log.Printf("开启数据库跟踪失败: %v", err)
		return
	}
	log.Println("✅ 已开启生产环境IndexedDB跟踪，可实时监控数据增删改查行为")
}

*/

// -----------------------------------------------  Storage.trackIndexedDBForStorageKey  -----------------------------------------------
// === 应用场景 ===
// 1. 分区IndexedDB监控：基于StorageKey开启分区存储下的IndexedDB跟踪，适配浏览器存储隔离
// 2. 第三方站点数据库调试：跟踪第三方嵌入站点的独立IndexedDB分区，不影响主站
// 3. PWA多实例数据监控：监控同一域名下不同分区、不同实例的PWA离线数据库
// 4. 自动化测试：针对带存储分区的测试环境做IndexedDB行为隔离监控
// 5. 数据问题定位：精准定位分区存储下的IndexedDB数据丢失、写入失败、查询异常
// 6. 多租户数据管理：管理不同存储分区下的独立数据库，互不干扰

// CDPStorageTrackIndexedDBForStorageKey 开启跟踪指定StorageKey的IndexedDB存储
// storageKey: 存储键，用于标识分区存储的唯一键
func CDPStorageTrackIndexedDBForStorageKey(storageKey string) (string, error) {
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
		"method": "Storage.trackIndexedDBForStorageKey",
		"params": {
			"storageKey": "%s"
		}
	}`, reqID, storageKey)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 trackIndexedDBForStorageKey 请求失败: %w", err)
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
			return "", fmt.Errorf("trackIndexedDBForStorageKey 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启指定StorageKey分区IndexedDB跟踪，调试分区存储PWA
func ExampleTrackPartitionIndexedDBByStorageKey() {
	storageKey := "https://pwa.example.com^partition1"
	resp, err := CDPStorageTrackIndexedDBForStorageKey(storageKey)
	if err != nil {
		log.Fatalf("开启分区IndexedDB跟踪失败: %v", err)
	}
	log.Printf("✅ 成功开启 StorageKey[%s] IndexedDB 跟踪", storageKey)
}

// 场景2：自动化测试 - 分区数据库环境监控
func ExampleTestSetupPartitionIndexedDBTracking() {
	testStorageKey := "https://test.app^test-partition"
	resp, err := CDPStorageTrackIndexedDBForStorageKey(testStorageKey)
	if err != nil {
		log.Fatalf("测试分区数据库监控开启失败: %v", err)
	}
	log.Println("✅ 测试分区IndexedDB跟踪已启动，可监控数据库变更事件")
}

// 场景3：第三方站点分区IndexedDB调试
func ExampleDebugThirdPartyIndexedDB() {
	storageKey := "https://third-party.com^embed-partition"
	resp, err := CDPStorageTrackIndexedDBForStorageKey(storageKey)
	if err != nil {
		log.Printf("第三方IndexedDB跟踪失败: %v", err)
		return
	}
	log.Println("✅ 第三方分区数据库跟踪已开启，可排查数据异常问题")
}

*/

// -----------------------------------------------  Storage.untrackCacheStorageForOrigin  -----------------------------------------------
// === 应用场景 ===
// 1. 关闭缓存监控：停止对指定源 CacheStorage 的跟踪，释放资源
// 2. 测试环境清理：自动化测试完成后，关闭缓存跟踪避免干扰
// 3. 按需跟踪：完成调试后关闭跟踪，减少浏览器性能消耗
// 4. 多站点切换：切换监控站点时，关闭上一个站点的缓存跟踪
// 5. 资源释放：长时间运行后关闭无用跟踪，降低内存占用
// 6. 流程收尾：PWA调试完成后，正常关闭缓存监听流程

// CDPStorageUntrackCacheStorageForOrigin 停止跟踪指定源的 CacheStorage 存储
// origin: 需要停止跟踪的站点源地址，例如 "https://www.example.com"
func CDPStorageUntrackCacheStorageForOrigin(origin string) (string, error) {
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
		"method": "Storage.untrackCacheStorageForOrigin",
		"params": {
			"origin": "%s"
		}
	}`, reqID, origin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 untrackCacheStorageForOrigin 请求失败: %w", err)
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
			return "", fmt.Errorf("untrackCacheStorageForOrigin 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：PWA调试完成后关闭CacheStorage跟踪
func ExampleUntrackPwaCacheStorage() {
	origin := "https://pwa.example.com"
	resp, err := CDPStorageUntrackCacheStorageForOrigin(origin)
	if err != nil {
		log.Fatalf("关闭PWA缓存跟踪失败: %v", err)
	}
	log.Printf("✅ 成功停止 [%s] CacheStorage 跟踪", origin)
}

// 场景2：自动化测试用例执行完毕，清理缓存跟踪
func ExampleTestCleanupCacheTracking() {
	testOrigin := "https://test.app.com"
	resp, err := CDPStorageUntrackCacheStorageForOrigin(testOrigin)
	if err != nil {
		log.Printf("测试缓存跟踪关闭失败: %v", err)
		return
	}
	log.Println("✅ 测试站点缓存跟踪已关闭，环境清理完成")
}

// 场景3：调试完成释放资源，关闭生产站点缓存监控
func ExampleDebugReleaseCacheResource() {
	origin := "https://www.example.com"
	resp, err := CDPStorageUntrackCacheStorageForOrigin(origin)
	if err != nil {
		log.Printf("关闭缓存监控失败: %v", err)
		return
	}
	log.Println("✅ 缓存监控已关闭，浏览器资源已释放")
}

*/

// -----------------------------------------------  Storage.untrackCacheStorageForStorageKey  -----------------------------------------------
// === 应用场景 ===
// 1. 关闭分区缓存监控：停止基于StorageKey的分区CacheStorage跟踪，释放资源
// 2. 测试环境清理：自动化分区测试完成后，关闭缓存跟踪避免干扰与内存泄漏
// 3. 按需跟踪：分区存储调试完成后关闭跟踪，降低浏览器性能消耗
// 4. 多分区切换：切换监控存储分区时，关闭上一个分区的缓存跟踪
// 5. 资源释放：长时间运行后关闭无用的分区跟踪，减少资源占用
// 6. 流程收尾：第三方站点/分区PWA调试完成后，正常关闭缓存监听

// CDPStorageUntrackCacheStorageForStorageKey 停止跟踪指定StorageKey的CacheStorage存储
// storageKey: 存储键，用于标识分区存储的唯一键
func CDPStorageUntrackCacheStorageForStorageKey(storageKey string) (string, error) {
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
		"method": "Storage.untrackCacheStorageForStorageKey",
		"params": {
			"storageKey": "%s"
		}
	}`, reqID, storageKey)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 untrackCacheStorageForStorageKey 请求失败: %w", err)
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
			return "", fmt.Errorf("untrackCacheStorageForStorageKey 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：分区PWA调试完成后，关闭指定StorageKey缓存跟踪
func ExampleUntrackPartitionCacheByStorageKey() {
	storageKey := "https://pwa.example.com^partition1"
	resp, err := CDPStorageUntrackCacheStorageForStorageKey(storageKey)
	if err != nil {
		log.Fatalf("关闭分区缓存跟踪失败: %v", err)
	}
	log.Printf("✅ 成功停止 StorageKey[%s] CacheStorage 跟踪", storageKey)
}

// 场景2：自动化分区测试完成，清理缓存跟踪资源
func ExampleTestCleanupPartitionCacheTracking() {
	testStorageKey := "https://test.app^test-partition"
	resp, err := CDPStorageUntrackCacheStorageForStorageKey(testStorageKey)
	if err != nil {
		log.Printf("测试分区缓存跟踪关闭失败: %v", err)
		return
	}
	log.Println("✅ 测试分区缓存跟踪已关闭，资源释放完成")
}

// 场景3：第三方站点调试结束，关闭分区缓存监控
func ExampleUntrackThirdPartyCache() {
	storageKey := "https://third-party.com^embed-partition"
	resp, err := CDPStorageUntrackCacheStorageForStorageKey(storageKey)
	if err != nil {
		log.Printf("第三方分区缓存跟踪关闭失败: %v", err)
		return
	}
	log.Println("✅ 第三方站点缓存监控已关闭")
}

*/

// -----------------------------------------------  Storage.untrackIndexedDBForOrigin  -----------------------------------------------
// === 应用场景 ===
// 1. 关闭IndexedDB监控：停止对指定源IndexedDB数据库的跟踪，释放浏览器监听资源
// 2. 测试环境清理：自动化测试完成后关闭IndexedDB跟踪，避免资源占用
// 3. 按需跟踪：调试完成后关闭跟踪，降低性能消耗
// 4. 多站点切换：切换监控站点时，关闭上一个站点的IndexedDB跟踪
// 5. 内存优化：长时间运行后关闭无用跟踪，减少内存泄漏风险
// 6. 流程收尾：PWA/离线应用调试完毕，正常结束数据库监听流程

// CDPStorageUntrackIndexedDBForOrigin 停止跟踪指定源的IndexedDB存储
// origin: 需要停止跟踪的站点源地址，例如 "https://www.example.com"
func CDPStorageUntrackIndexedDBForOrigin(origin string) (string, error) {
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
		"method": "Storage.untrackIndexedDBForOrigin",
		"params": {
			"origin": "%s"
		}
	}`, reqID, origin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 untrackIndexedDBForOrigin 请求失败: %w", err)
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
			return "", fmt.Errorf("untrackIndexedDBForOrigin 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：PWA调试完成后关闭IndexedDB跟踪，释放资源
func ExampleUntrackPwaIndexedDB() {
	origin := "https://pwa.example.com"
	resp, err := CDPStorageUntrackIndexedDBForOrigin(origin)
	if err != nil {
		log.Fatalf("关闭PWA IndexedDB跟踪失败: %v", err)
	}
	log.Printf("✅ 成功停止 [%s] IndexedDB 跟踪", origin)
}

// 场景2：自动化测试用例执行完毕，清理IndexedDB跟踪
func ExampleTestCleanupIndexedDBTracking() {
	testOrigin := "https://test.app.com"
	resp, err := CDPStorageUntrackIndexedDBForOrigin(testOrigin)
	if err != nil {
		log.Printf("测试IndexedDB跟踪关闭失败: %v", err)
		return
	}
	log.Println("✅ 测试站点IndexedDB跟踪已关闭，环境清理完成")
}

// 场景3：调试结束释放资源，关闭生产站点IndexedDB监控
func ExampleDebugReleaseIndexedDBResource() {
	origin := "https://www.example.com"
	resp, err := CDPStorageUntrackIndexedDBForOrigin(origin)
	if err != nil {
		log.Printf("关闭IndexedDB监控失败: %v", err)
		return
	}
	log.Println("✅ IndexedDB监控已关闭，浏览器资源已释放")
}

*/

// -----------------------------------------------  Storage.untrackIndexedDBForStorageKey  -----------------------------------------------
// === 应用场景 ===
// 1. 关闭分区IndexedDB监控：停止基于StorageKey的分区数据库跟踪，释放浏览器资源
// 2. 测试环境清理：自动化分区测试完成后，关闭IndexedDB跟踪避免内存泄漏
// 3. 按需跟踪：分区存储调试完成后关闭跟踪，降低浏览器性能消耗
// 4. 多分区切换：切换监控存储分区时，关闭上一个分区的IndexedDB跟踪
// 5. 资源释放：长时间运行后关闭无用的分区跟踪，减少内存与CPU占用
// 6. 流程收尾：第三方站点/分区PWA调试完成后，正常关闭数据库监听

// CDPStorageUntrackIndexedDBForStorageKey 停止跟踪指定StorageKey的IndexedDB存储
// storageKey: 存储键，用于标识分区存储的唯一键
func CDPStorageUntrackIndexedDBForStorageKey(storageKey string) (string, error) {
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
		"method": "Storage.untrackIndexedDBForStorageKey",
		"params": {
			"storageKey": "%s"
		}
	}`, reqID, storageKey)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 untrackIndexedDBForStorageKey 请求失败: %w", err)
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
			return "", fmt.Errorf("untrackIndexedDBForStorageKey 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：分区PWA调试完成后，关闭指定StorageKey的IndexedDB跟踪
func ExampleUntrackPartitionIndexedDBByStorageKey() {
	storageKey := "https://pwa.example.com^partition1"
	resp, err := CDPStorageUntrackIndexedDBForStorageKey(storageKey)
	if err != nil {
		log.Fatalf("关闭分区IndexedDB跟踪失败: %v", err)
	}
	log.Printf("✅ 成功停止 StorageKey[%s] IndexedDB 跟踪", storageKey)
}

// 场景2：自动化分区测试结束，清理IndexedDB跟踪资源
func ExampleTestCleanupPartitionIndexedDBTracking() {
	testStorageKey := "https://test.app^test-partition"
	resp, err := CDPStorageUntrackIndexedDBForStorageKey(testStorageKey)
	if err != nil {
		log.Printf("测试分区IndexedDB跟踪关闭失败: %v", err)
		return
	}
	log.Println("✅ 测试分区IndexedDB跟踪已关闭，资源释放完成")
}

// 场景3：第三方站点调试完毕，关闭分区数据库监控
func ExampleUntrackThirdPartyIndexedDB() {
	storageKey := "https://third-party.com^embed-partition"
	resp, err := CDPStorageUntrackIndexedDBForStorageKey(storageKey)
	if err != nil {
		log.Printf("第三方分区IndexedDB跟踪关闭失败: %v", err)
		return
	}
	log.Println("✅ 第三方站点IndexedDB监控已关闭")
}

*/

// -----------------------------------------------  Storage.clearSharedStorageEntries  -----------------------------------------------
// === 应用场景 ===
// 1. 共享存储数据清空：清除指定源的 Shared Storage 所有数据
// 2. 隐私沙箱调试：调试 Shared Storage 隐私存储功能，重置数据环境
// 3. 自动化测试：测试用例执行前清空共享存储，保证环境纯净无干扰
// 4. 数据隔离：多账号/多场景测试时清除上一轮的共享存储数据
// 5. 存储重置：修复 Shared Storage 异常、脏数据问题
// 6. 开发环境清理：开发过程中快速重置站点隐私存储状态

// CDPStorageClearSharedStorageEntries 清除指定源的 Shared Storage 条目
// origin: 要清除共享存储的站点源地址，例如 "https://www.example.com"
func CDPStorageClearSharedStorageEntries(origin string) (string, error) {
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
		"method": "Storage.clearSharedStorageEntries",
		"params": {
			"ownerOrigin": "%s"
		}
	}`, reqID, origin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearSharedStorageEntries 请求失败: %w", err)
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
			return "", fmt.Errorf("clearSharedStorageEntries 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：清除指定网站的所有 Shared Storage 数据
func ExampleClearSharedStorage() {
	origin := "https://example.com"
	resp, err := CDPStorageClearSharedStorageEntries(origin)
	if err != nil {
		log.Fatalf("清除 Shared Storage 失败: %v", err)
	}
	log.Printf("✅ 成功清除 [%s] 的所有 Shared Storage 条目", origin)
}

// 场景2：自动化测试前置 - 清空 Shared Storage 保证测试环境干净
func ExampleTestSetupClearSharedStorage() {
	testOrigin := "https://test.example.com"
	resp, err := CDPStorageClearSharedStorageEntries(testOrigin)
	if err != nil {
		log.Fatalf("测试环境清理 Shared Storage 失败: %v", err)
	}
	log.Println("✅ 测试站点 Shared Storage 已清空，可开始执行测试")
}

// 场景3：调试隐私沙箱功能，重置 Shared Storage 状态
func ExampleDebugResetSharedStorage() {
	origin := "https://privacy-sandbox.example.com"
	resp, err := CDPStorageClearSharedStorageEntries(origin)
	if err != nil {
		log.Printf("重置 Shared Storage 失败: %v", err)
		return
	}
	log.Println("✅ Shared Storage 已重置，可重新调试隐私存储功能")
}

*/

// -----------------------------------------------  Storage.clearTrustTokens  -----------------------------------------------
// === 应用场景 ===
// 1. 信任令牌清除: 清除浏览器所有Trust Tokens，重置身份验证状态
// 2. 隐私沙箱调试: 调试Trust Token隐私API时快速清空令牌数据
// 3. 自动化测试: 测试用例执行前清空Trust Tokens，保证测试环境独立
// 4. 身份状态重置: 清除网站信任凭证，强制重新进行身份验证
// 5. 安全测试: 验证无信任令牌时的网站行为与权限控制
// 6. 错误恢复: 信任令牌失效/异常时清空并重新获取

// CDPStorageClearTrustTokens 清除浏览器所有Trust Tokens信任令牌
func CDPStorageClearTrustTokens() (string, error) {
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
		"method": "Storage.clearTrustTokens"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearTrustTokens 请求失败: %w", err)
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
			return "", fmt.Errorf("clearTrustTokens 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：清空所有Trust Tokens，重置浏览器信任状态
func ExampleClearAllTrustTokens() {
	resp, err := CDPStorageClearTrustTokens()
	if err != nil {
		log.Fatalf("清除Trust Tokens失败: %v", err)
	}
	log.Println("✅ 成功清除浏览器所有Trust Tokens信任令牌")
}

// 场景2：自动化测试前置 - 清空信任令牌保证测试环境纯净
func ExampleTestSetupClearTrustTokens() {
	resp, err := CDPStorageClearTrustTokens()
	if err != nil {
		log.Fatalf("测试环境清理Trust Tokens失败: %v", err)
	}
	log.Println("✅ 信任令牌已清空，可执行无身份凭证测试用例")
}

// 场景3：调试隐私沙箱Trust Token功能，重置令牌状态
func ExampleDebugResetTrustTokens() {
	resp, err := CDPStorageClearTrustTokens()
	if err != nil {
		log.Printf("重置Trust Tokens失败: %v", err)
		return
	}
	log.Println("✅ Trust Tokens已重置，可重新进行身份签发与验证")
}

*/

// -----------------------------------------------  Storage.deleteSharedStorageEntry  -----------------------------------------------
// === 应用场景 ===
// 1. 精准删除共享存储条目：只删除指定key的Shared Storage数据，不影响其他键值
// 2. 隐私沙箱调试：调试Shared Storage时删除特定测试键值
// 3. 自动化测试：测试过程中删除指定条目，验证业务逻辑容错性
// 4. 数据清理：清理过期/无用的单个共享存储键值
// 5. 状态重置：重置某个功能对应的共享存储状态
// 6. 多场景隔离：测试不同场景时单独删除对应key避免干扰

// CDPStorageDeleteSharedStorageEntry 删除指定源的单个Shared Storage键值对
// origin: 站点源地址，例如 "https://www.example.com"
// key: 要删除的共享存储键名
func CDPStorageDeleteSharedStorageEntry(origin string, key string) (string, error) {
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
		"method": "Storage.deleteSharedStorageEntry",
		"params": {
			"ownerOrigin": "%s",
			"key": "%s"
		}
	}`, reqID, origin, key)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deleteSharedStorageEntry 请求失败: %w", err)
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
			return "", fmt.Errorf("deleteSharedStorageEntry 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：删除指定网站的单个Shared Storage键
func ExampleDeleteSingleSharedStorageEntry() {
	origin := "https://example.com"
	key := "user_preference"
	resp, err := CDPStorageDeleteSharedStorageEntry(origin, key)
	if err != nil {
		log.Fatalf("删除 Shared Storage 键失败: %v", err)
	}
	log.Printf("✅ 成功删除 [%s] 中的键: %s", origin, key)
}

// 场景2：自动化测试 - 删除测试用的共享存储条目
func ExampleTestDeleteSharedStorageKey() {
	testOrigin := "https://test.example.com"
	testKey := "test_flag"
	resp, err := CDPStorageDeleteSharedStorageEntry(testOrigin, testKey)
	if err != nil {
		log.Printf("测试删除 Shared Storage 失败: %v", err)
		return
	}
	log.Println("✅ 测试键已删除，可继续执行后续测试逻辑")
}

// 场景3：调试隐私沙箱，清理无效共享存储键
func ExampleDebugCleanInvalidSharedStorage() {
	origin := "https://privacy-sandbox.app"
	invalidKey := "deprecated_key"
	resp, err := CDPStorageDeleteSharedStorageEntry(origin, invalidKey)
	if err != nil {
		log.Printf("清理无效键失败: %v", err)
		return
	}
	log.Println("✅ 无效共享存储条目已删除")
}

*/

// -----------------------------------------------  Storage.deleteStorageBucket  -----------------------------------------------
// === 应用场景 ===
// 1. 存储桶删除：直接删除指定的存储桶（Storage Bucket），释放所有关联数据
// 2. 分区存储清理：删除浏览器分桶存储的独立数据分区，彻底释放空间
// 3. 自动化测试：测试完成后删除测试专用存储桶，保证环境干净
// 4. 调试存储异常：删除异常存储桶，解决数据损坏、无法访问问题
// 5. 隐私数据清理：彻底删除用户隐私相关的分桶存储数据
// 6. 开发环境重置：重置存储桶状态，重新测试存储分配逻辑

// CDPStorageDeleteStorageBucket 删除指定的存储桶
// storageBucketId: 要删除的存储桶唯一ID
func CDPStorageDeleteStorageBucket(storageBucketId string) (string, error) {
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
		"method": "Storage.deleteStorageBucket",
		"params": {
			"storageBucketId": "%s"
		}
	}`, reqID, storageBucketId)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deleteStorageBucket 请求失败: %w", err)
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
			return "", fmt.Errorf("deleteStorageBucket 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：删除指定ID的存储桶，彻底清理数据
func ExampleDeleteStorageBucket() {
	bucketID := "example-bucket-id-123456"
	resp, err := CDPStorageDeleteStorageBucket(bucketID)
	if err != nil {
		log.Fatalf("删除存储桶失败: %v", err)
	}
	log.Printf("✅ 成功删除存储桶: %s", bucketID)
}

// 场景2：自动化测试 - 清理测试创建的存储桶
func ExampleTestCleanupStorageBucket() {
	testBucketID := "test-bucket-789"
	resp, err := CDPStorageDeleteStorageBucket(testBucketID)
	if err != nil {
		log.Printf("测试存储桶删除失败: %v", err)
		return
	}
	log.Println("✅ 测试存储桶已删除，测试环境清理完成")
}

// 场景3：调试存储桶异常，删除损坏的存储桶
func ExampleDebugDeleteCorruptedBucket() {
	corruptedID := "corrupted-bucket-001"
	resp, err := CDPStorageDeleteStorageBucket(corruptedID)
	if err != nil {
		log.Fatalf("删除异常存储桶失败: %v", err)
	}
	log.Println("✅ 异常存储桶已删除，可重新创建正常存储桶")
}

*/

// -----------------------------------------------  Storage.getAffectedUrlsForThirdPartyCookieMetadata  -----------------------------------------------
// === 应用场景 ===
// 1. 第三方Cookie审计：查询哪些URL会受第三方Cookie元数据影响
// 2. 隐私合规检测：验证第三方Cookie政策对站点的影响范围
// 3. 广告/统计调试：排查第三方Cookie被拦截导致的功能异常
// 4. 自动化测试：验证第三方Cookie策略生效范围
// 5. 站点兼容性分析：检查第三方Cookie依赖的URL清单
// 6. 安全审计：识别存在第三方Cookie风险的URL路径

// CDPStorageGetAffectedUrlsForThirdPartyCookieMetadata 获取受第三方Cookie元数据影响的URL列表
// cookieMetadataJson: 第三方Cookie元数据的JSON字符串
func CDPStorageGetAffectedUrlsForThirdPartyCookieMetadata(cookieMetadataJson string) (string, error) {
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
		"method": "Storage.getAffectedUrlsForThirdPartyCookieMetadata",
		"params": {
			"cookieMetadata": %s
		}
	}`, reqID, cookieMetadataJson)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getAffectedUrlsForThirdPartyCookieMetadata 请求失败: %w", err)
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
			return "", fmt.Errorf("getAffectedUrlsForThirdPartyCookieMetadata 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：查询指定第三方Cookie元数据影响的URL（基础用法）
func ExampleGetAffectedUrlsBasic() {
	// 第三方Cookie元数据（name/domain/path等）
	metadata := `{
		"name": "third_party_tracker",
		"domain": ".analytics.com",
		"path": "/"
	}`

	resp, err := CDPStorageGetAffectedUrlsForThirdPartyCookieMetadata(metadata)
	if err != nil {
		log.Fatalf("查询受影响URL失败: %v", err)
	}
	log.Printf("✅ 受第三方Cookie影响的URL列表: %s", resp)
}

// 场景2：隐私合规审计 - 检测第三方Cookie影响范围
func ExamplePrivacyAuditAffectedUrls() {
	metadata := `{
		"name": "ad_cookie",
		"domain": ".adsnetwork.com",
		"secure": true
	}`

	resp, err := CDPStorageGetAffectedUrlsForThirdPartyCookieMetadata(metadata)
	if err != nil {
		log.Printf("合规审计查询失败: %v", err)
		return
	}
	log.Println("✅ 第三方Cookie合规影响范围查询完成")
}

// 场景3：自动化测试 - 验证第三方Cookie策略影响URL
func ExampleTestAffectedUrls() {
	testMetadata := `{
		"name": "test_partner_cookie",
		"domain": ".test-partner.com"
	}`

	resp, err := CDPStorageGetAffectedUrlsForThirdPartyCookieMetadata(testMetadata)
	if err != nil {
		log.Fatalf("测试查询失败: %v", err)
	}
	log.Println("✅ 测试用第三方Cookie影响URL获取成功")
}

*/

// -----------------------------------------------  Storage.getInterestGroupDetails  -----------------------------------------------
// === 应用场景 ===
// 1. 广告兴趣组调试：获取 Protected Audience API 兴趣组的完整配置信息
// 2. 隐私沙箱测试：查询浏览器本地存储的广告兴趣组数据
// 3. 自动化测试：验证兴趣组是否正确加入、更新、配置生效
// 4. 广告投放分析：检查兴趣组出价、广告、元数据等核心配置
// 5. 开发环境调试：本地调试广告联盟、需求方平台（DSP）逻辑
// 6. 数据合规审计：查看兴趣组存储的用户数据与权限配置

// CDPStorageGetInterestGroupDetails 获取指定广告兴趣组的详细信息
// ownerOrigin: 兴趣组所属源站，例如 "https://www.example.com"
// name: 兴趣组名称
func CDPStorageGetInterestGroupDetails(ownerOrigin string, name string) (string, error) {
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
		"method": "Storage.getInterestGroupDetails",
		"params": {
			"ownerOrigin": "%s",
			"name": "%s"
		}
	}`, reqID, ownerOrigin, name)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getInterestGroupDetails 请求失败: %w", err)
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
			return "", fmt.Errorf("getInterestGroupDetails 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：获取指定广告兴趣组完整详情（基础用法）
func ExampleGetInterestGroupBasic() {
	origin := "https://advertiser.example.com"
	name := "shoes_interest_group"

	resp, err := CDPStorageGetInterestGroupDetails(origin, name)
	if err != nil {
		log.Fatalf("获取兴趣组信息失败: %v", err)
	}
	log.Printf("✅ 兴趣组详情:\n%s", resp)
}

// 场景2：自动化测试 - 验证兴趣组是否正确创建
func ExampleTestVerifyInterestGroup() {
	testOrigin := "https://test-dsp.example.com"
	testName := "test_electronics"

	resp, err := CDPStorageGetInterestGroupDetails(testOrigin, testName)
	if err != nil {
		log.Fatalf("测试兴趣组查询失败: %v", err)
	}

	if strings.Contains(resp, "biddingLogicUrl") {
		log.Println("✅ 兴趣组配置正常，可参与广告竞价")
	} else {
		log.Println("❌ 兴趣组信息不完整，配置异常")
	}
}

// 场景3：隐私沙箱广告调试，查看兴趣组出价与广告配置
func ExampleDebugProtectedAudience() {
	origin := "https://privacy-sandbox.example.com"
	name := "campaign_2025"

	resp, err := CDPStorageGetInterestGroupDetails(origin, name)
	if err != nil {
		log.Printf("广告调试查询失败: %v", err)
		return
	}
	log.Println("✅ 广告兴趣组配置已获取，可调试竞价逻辑")
}

*/

// -----------------------------------------------  Storage.getRelatedWebsiteSets  -----------------------------------------------
// === 应用场景 ===
// 1. 关联网站集查询：获取浏览器当前配置的 Related Website Sets（相关网站集合）
// 2. 第三方Cookie调试：验证关联网站是否可在隐私限制下共享Cookie
// 3. 跨站身份验证调试：检查集团旗下多网站是否属于同一关联集合
// 4. 自动化测试：验证关联网站集配置是否正确加载
// 5. 隐私合规检测：确认跨站数据共享范围符合规范
// 6. 站点权限排查：定位跨站存储、登录态失效问题

// CDPStorageGetRelatedWebsiteSets 获取浏览器当前的 Related Website Sets
func CDPStorageGetRelatedWebsiteSets() (string, error) {
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
		"method": "Storage.getRelatedWebsiteSets"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getRelatedWebsiteSets 请求失败: %w", err)
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
			return "", fmt.Errorf("getRelatedWebsiteSets 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：获取浏览器全部关联网站集配置
func ExampleGetRelatedWebsiteSets() {
	resp, err := CDPStorageGetRelatedWebsiteSets()
	if err != nil {
		log.Fatalf("获取关联网站集失败: %v", err)
	}
	log.Printf("✅ 关联网站集配置:\n%s", resp)
}

// 场景2：自动化测试 - 验证关联网站集是否正确加载
func ExampleTestVerifyRelatedWebsiteSets() {
	resp, err := CDPStorageGetRelatedWebsiteSets()
	if err != nil {
		log.Fatalf("测试查询失败: %v", err)
	}

	// 检查是否包含预期的关联网站
	if strings.Contains(resp, "example.com") && strings.Contains(resp, "example.co.jp") {
		log.Println("✅ 关联网站集配置正确，跨站权限正常")
	} else {
		log.Println("❌ 未找到预期的关联网站配置")
	}
}

// 场景3：调试第三方Cookie跨站失效问题
func ExampleDebugThirdPartyCookieByRWS() {
	resp, err := CDPStorageGetRelatedWebsiteSets()
	if err != nil {
		log.Printf("调试查询失败: %v", err)
		return
	}
	log.Println("✅ 已获取关联网站集，可排查跨站Cookie共享权限")
}

*/

// -----------------------------------------------  Storage.getSharedStorageEntries  -----------------------------------------------
// === 应用场景 ===
// 1. 共享存储读取：获取指定源的 Shared Storage 所有键值对数据
// 2. 隐私沙箱调试：调试 Shared Storage 隐私存储功能，查看存储内容
// 3. 自动化测试：验证共享存储数据是否正确写入、读取
// 4. 数据审计：查看站点隐私存储的具体内容与配置
// 5. 功能调试：排查基于 Shared Storage 的业务逻辑异常
// 6. 开发验证：确认共享存储数据符合预期

// CDPStorageGetSharedStorageEntries 获取指定源的 Shared Storage 全部条目
// ownerOrigin: 共享存储所属源地址，例如 "https://www.example.com"
func CDPStorageGetSharedStorageEntries(ownerOrigin string) (string, error) {
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
		"method": "Storage.getSharedStorageEntries",
		params": {
			"ownerOrigin": "%s"
		}
	}`, reqID, ownerOrigin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getSharedStorageEntries 请求失败: %w", err)
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
			return "", fmt.Errorf("getSharedStorageEntries 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：获取指定网站的所有 Shared Storage 数据
func ExampleGetSharedStorageEntries() {
	origin := "https://example.com"
	resp, err := CDPStorageGetSharedStorageEntries(origin)
	if err != nil {
		log.Fatalf("获取 Shared Storage 失败: %v", err)
	}
	log.Printf("✅ [%s] 的共享存储数据:\n%s", origin, resp)
}

// 场景2：自动化测试 - 验证共享存储数据是否正确
func ExampleTestCheckSharedStorageData() {
	testOrigin := "https://test.example.com"
	resp, err := CDPStorageGetSharedStorageEntries(testOrigin)
	if err != nil {
		log.Fatalf("测试获取共享存储失败: %v", err)
	}

	// 检查是否包含预期键
	if strings.Contains(resp, "user_config") {
		log.Println("✅ 共享存储数据正常，包含预期键值")
	} else {
		log.Println("❌ 共享存储未找到预期数据")
	}
}

// 场景3：调试隐私沙箱，查看共享存储内容
func ExampleDebugSharedStorageContent() {
	origin := "https://privacy-sandbox.app"
	resp, err := CDPStorageGetSharedStorageEntries(origin)
	if err != nil {
		log.Printf("调试获取共享存储失败: %v", err)
		return
	}
	log.Println("✅ 共享存储内容获取成功，可进行调试分析")
}

*/

// -----------------------------------------------  Storage.getSharedStorageMetadata  -----------------------------------------------
// === 应用场景 ===
// 1. 共享存储元数据查询：获取指定源 Shared Storage 的元信息（非键值数据）
// 2. 隐私沙箱调试：查看 Shared Storage 存储大小、创建时间、剩余容量等信息
// 3. 自动化测试：验证共享存储是否正常初始化、配额使用情况
// 4. 存储监控：监控站点共享存储占用，判断是否需要清理
// 5. 合规审计：查看共享存储基础信息，确认隐私合规状态
// 6. 问题排查：定位 Shared Storage 初始化失败、写入受限问题

// CDPStorageGetSharedStorageMetadata 获取指定源的 Shared Storage 元数据
// ownerOrigin: 共享存储所属源地址，例如 "https://www.example.com"
func CDPStorageGetSharedStorageMetadata(ownerOrigin string) (string, error) {
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
		"method": "Storage.getSharedStorageMetadata",
		"params": {
			"ownerOrigin": "%s"
		}
	}`, reqID, ownerOrigin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getSharedStorageMetadata 请求失败: %w", err)
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
			return "", fmt.Errorf("getSharedStorageMetadata 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：获取指定网站 Shared Storage 元数据（基础用法）
func ExampleGetSharedStorageMetadata() {
	origin := "https://example.com"
	resp, err := CDPStorageGetSharedStorageMetadata(origin)
	if err != nil {
		log.Fatalf("获取共享存储元数据失败: %v", err)
	}
	log.Printf("✅ [%s] 共享存储元数据:\n%s", origin, resp)
}

// 场景2：自动化测试 - 验证共享存储元数据是否合法
func ExampleTestCheckSharedStorageMeta() {
	testOrigin := "https://test.example.com"
	resp, err := CDPStorageGetSharedStorageMetadata(testOrigin)
	if err != nil {
		log.Fatalf("测试获取元数据失败: %v", err)
	}

	if strings.Contains(resp, "bytesUsed") && strings.Contains(resp, "creationTime") {
		log.Println("✅ 共享存储元数据完整，状态正常")
	} else {
		log.Println("❌ 共享存储元数据异常")
	}
}

// 场景3：调试隐私沙箱，监控共享存储占用
func ExampleDebugSharedStorageMeta() {
	origin := "https://privacy-sandbox.app"
	resp, err := CDPStorageGetSharedStorageMetadata(origin)
	if err != nil {
		log.Printf("调试获取元数据失败: %v", err)
		return
	}
	log.Println("✅ 共享存储元数据获取成功，可分析存储状态")
}

*/

// -----------------------------------------------  Storage.getStorageKey -----------------------------------------------
// === 应用场景 ===
// 1. 获取存储键：根据 origin、topFrameOrigin 等信息生成标准 StorageKey
// 2. 分区存储调试：用于调试浏览器存储分区、第三方隔离机制
// 3. CacheStorage / IndexedDB 定位：精准定位分区存储的归属
// 4. 自动化测试：生成稳定的 StorageKey 用于跟踪/删除缓存/数据库
// 5. 跨站存储分析：验证第一方/第三方存储分区是否正确隔离
// 6. 隐私沙箱调试：确认存储分区规则符合最新 Chrome 规范

// CDPStorageGetStorageKey 根据源信息获取浏览器标准 StorageKey
// origin: 资源源地址
// topFrameOrigin: 顶层框架源地址（可为空）
func CDPStorageGetStorageKey(origin string, topFrameOrigin string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息
	var message string
	if topFrameOrigin == "" {
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "Storage.getStorageKey",
			"params": {
				"origin": "%s"
			}
		}`, reqID, origin)
	} else {
		message = fmt.Sprintf(`{
			"id": %d,
			"method": "Storage.getStorageKey",
			"params": {
				"origin": "%s",
				"topFrameOrigin": "%s"
			}
		}`, reqID, origin, topFrameOrigin)
	}

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getStorageKey 请求失败: %w", err)
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
			return "", fmt.Errorf("getStorageKey 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：获取第一方页面 StorageKey（无顶层框架区分）
func ExampleGetStorageKeyFirstParty() {
	origin := "https://www.example.com"
	resp, err := CDPStorageGetStorageKey(origin, "")
	if err != nil {
		log.Fatalf("获取 StorageKey 失败: %v", err)
	}
	log.Printf("✅ 第一方 StorageKey: %s", resp)
}

// 场景2：获取第三方嵌入资源 StorageKey（带顶层框架）
func ExampleGetStorageKeyThirdParty() {
	topOrigin := "https://parent-site.com"
	origin := "https://third-party.com"
	resp, err := CDPStorageGetStorageKey(origin, topOrigin)
	if err != nil {
		log.Fatalf("获取第三方 StorageKey 失败: %v", err)
	}
	log.Printf("✅ 第三方分区 StorageKey: %s", resp)
}

// 场景3：自动化测试 - 生成 StorageKey 用于跟踪缓存/IndexedDB
func ExampleTestSetupStorageKey() {
	origin := "https://test.app"
	resp, err := CDPStorageGetStorageKey(origin, "")
	if err != nil {
		log.Fatalf("测试生成 StorageKey 失败: %v", err)
	}
	log.Println("✅ 测试用 StorageKey 生成完成，可用于缓存/数据库跟踪")
}

*/

// -----------------------------------------------  Storage.getTrustTokens  -----------------------------------------------
// === 应用场景 ===
// 1. 信任令牌查询：获取浏览器当前存储的所有 Trust Tokens 列表
// 2. 隐私沙箱调试：查看 Trust Token 签发、状态、有效期等信息
// 3. 自动化测试：验证令牌是否正确签发与存储
// 4. 身份凭证审计：检查第三方签发的信任令牌详情
// 5. 问题排查：定位令牌缺失、失效导致的验证失败
// 6. 合规检查：查看浏览器存储的隐私敏感身份令牌

// CDPStorageGetTrustTokens 获取浏览器中存储的所有 Trust Tokens
func CDPStorageGetTrustTokens() (string, error) {
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
		"method": "Storage.getTrustTokens"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getTrustTokens 请求失败: %w", err)
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
			return "", fmt.Errorf("getTrustTokens 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：获取浏览器所有 Trust Tokens 信息
func ExampleGetTrustTokens() {
	resp, err := CDPStorageGetTrustTokens()
	if err != nil {
		log.Fatalf("获取 Trust Tokens 失败: %v", err)
	}
	log.Printf("✅ 浏览器 Trust Tokens 信息:\n%s", resp)
}

// 场景2：自动化测试 - 验证令牌是否成功签发
func ExampleTestCheckTrustTokensIssued() {
	resp, err := CDPStorageGetTrustTokens()
	if err != {
		log.Fatalf("测试获取令牌失败: %v", err)
	}

	if strings.Contains(resp, "issuer") && strings.Contains
(resp, "token") {
		log.Println("✅ Trust Token 签发成功，已正常存储")
	} else {
		log.Println("❌ 未获取到有效 Trust Token")
	}
}

// 场景3：调试隐私沙箱身份验证
func ExampleDebugTrustTokens() {
	resp, err := CDPStorageGetTrustTokens()
	if err != nil {
		log.Printf("调试获取令牌失败: %v", err)
		return
	}
	log.Println("✅ Trust Token 信息已获取，可进行身份验证调试")
}

*/

// -----------------------------------------------  Storage.overrideQuotaForOrigin  -----------------------------------------------
// === 应用场景 ===
// 1. 前端存储调试：临时突破存储配额限制，测试大容量数据写入
// 2. PWA 开发调试：给离线应用分配更大缓存空间，验证极限存储场景
// 3. 自动化测试：固定配额大小，确保测试环境一致性
// 4. 性能压测：模拟不同配额大小，测试应用降级/提示逻辑
// 5. 问题复现：快速触发“存储已满”异常，验证错误处理
// 6. 开发环境优化：避免本地调试时频繁出现配额不足

// CDPStorageOverrideQuotaForOrigin 为指定源覆盖存储配额大小
// origin: 站点源地址
// quotaSize: 配额大小（字节），如 1024*1024*100 = 100MB
func CDPStorageOverrideQuotaForOrigin(origin string, quotaSize int64) (string, error) {
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
		"method": "Storage.overrideQuotaForOrigin",
		"params": {
			"origin": "%s",
			"quotaSize": %d
		}
	}`, reqID, origin, quotaSize)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 overrideQuotaForOrigin 请求失败: %w", err)
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
			return "", fmt.Errorf("overrideQuotaForOrigin 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：给测试网站设置 100MB 存储配额，方便调试大容量写入
func ExampleOverrideQuotaNormal() {
	origin := "https://localhost:3000"
	quota := int64(100 * 1024 * 1024) // 100MB
	resp, err := CDPStorageOverrideQuotaForOrigin(origin, quota)
	if err != nil {
		log.Fatalf("设置配额失败: %v", err)
	}
	log.Printf("✅ 成功为 [%s] 设置配额: 100MB", origin)
}

// 场景2：自动化测试 - 设置极小配额(1KB)快速触发存储已满异常
func ExampleTestTriggerQuotaFull() {
	origin := "https://test.app"
	quota := int64(1024) // 1KB
	resp, err := CDPStorageOverrideQuotaForOrigin(origin, quota)
	if err != nil {
		log.Fatalf("测试配额设置失败: %v", err)
	}
	log.Println("✅ 极小配额已设置，可快速测试存储已满逻辑")
}

// 场景3：PWA调试 - 分配大空间用于离线资源缓存
func ExampleDebugPwaStorageQuota() {
	origin := "https://pwa.example.com"
	quota := int64(500 * 1024 * 1024) // 500MB
	resp, err := CDPStorageOverrideQuotaForOrigin(origin, quota)
	if err != nil {
		log.Printf("PWA配额设置失败: %v", err)
		return
	}
	log.Println("✅ PWA存储配额已提升，离线缓存调试更顺畅")
}

*/

// -----------------------------------------------  Storage.resetSharedStorageBudget  -----------------------------------------------
// === 应用场景 ===
// 1. 共享存储预算重置：将指定源的 Shared Storage 操作预算恢复为初始值
// 2. 隐私沙箱调试：调试 Shared Storage 频率限制、预算耗尽场景
// 3. 自动化测试：每次测试前重置预算，确保测试环境一致
// 4. 开发环境：避免调试时频繁触发预算不足限制
// 5. 功能验证：测试预算耗尽后的恢复逻辑
// 6. 问题排查：解决 Shared Storage 写入被限制的问题

// CDPStorageResetSharedStorageBudget 重置指定源的 Shared Storage 操作预算
// ownerOrigin: 共享存储所属源地址
func CDPStorageResetSharedStorageBudget(ownerOrigin string) (string, error) {
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
		"method": "Storage.resetSharedStorageBudget",
		"params": {
			"ownerOrigin": "%s"
		}
	}`, reqID, ownerOrigin)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 resetSharedStorageBudget 请求失败: %w", err)
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
			return "", fmt.Errorf("resetSharedStorageBudget 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：重置指定网站的 Shared Storage 预算
func ExampleResetSharedStorageBudget() {
	origin := "https://example.com"
	resp, err := CDPStorageResetSharedStorageBudget(origin)
	if err != nil {
		log.Fatalf("重置 Shared Storage 预算失败: %v", err)
	}
	log.Printf("✅ 成功重置 [%s] 的共享存储预算", origin)
}

// 场景2：自动化测试 - 每次用例前重置预算
func ExampleTestResetSharedStorageBudget() {
	testOrigin := "https://test.example.com"
	resp, err := CDPStorageResetSharedStorageBudget(testOrigin)
	if err != nil {
		log.Fatalf("测试预算重置失败: %v", err)
	}
	log.Println("✅ 测试用 Shared Storage 预算已重置，可正常执行")
}

// 场景3：调试预算耗尽问题，恢复写入权限
func ExampleDebugResetSharedStorageBudget() {
	origin := "https://privacy-sandbox.app"
	resp, err := CDPStorageResetSharedStorageBudget(origin)
	if err != nil {
		log.Printf("调试预算重置失败: %v", err)
		return
	}
	log.Println("✅ Shared Storage 预算已恢复，可继续调试写入逻辑")
}

*/

// -----------------------------------------------  Storage.runBounceTrackingMitigations  -----------------------------------------------
// === 应用场景 ===
// 1. 反弹追踪防护：手动触发浏览器反弹追踪 mitigation 清理机制
// 2. 隐私防护调试：测试 bounce tracking 防护是否正常工作
// 3. 自动化测试：验证反弹追踪后存储、Cookie 清理逻辑
// 4. 隐私沙箱验证：测试浏览器反追踪功能有效性
// 5. 开发调试：手动触发清理，复现隐私保护行为
// 6. 安全测试：验证站点跳转追踪数据是否被正确清除

// CDPStorageRunBounceTrackingMitigations 手动触发浏览器反弹追踪缓解措施
func CDPStorageRunBounceTrackingMitigations() (string, error) {
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
		"method": "Storage.runBounceTrackingMitigations"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 runBounceTrackingMitigations 请求失败: %w", err)
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
			return "", fmt.Errorf("runBounceTrackingMitigations 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：手动触发反弹追踪防护清理
func ExampleRunBounceTrackingMitigations() {
	resp, err := CDPStorageRunBounceTrackingMitigations()
	if err != nil {
		log.Fatalf("触发反弹追踪防护失败: %v", err)
	}
	log.Printf("✅ 成功触发反弹追踪缓解措施:\n%s", resp)
}

// 场景2：自动化测试 - 验证反追踪机制是否生效
func ExampleTestBounceTrackingMitigations() {
	resp, err := CDPStorageRunBounceTrackingMitigations()
	if err != nil {
		log.Fatalf("测试触发反追踪失败: %v", err)
	}
	log.Println("✅ 反弹追踪防护已执行，可验证清理结果")
}

// 场景3：隐私保护调试，手动清理追踪状态
func ExampleDebugBounceTrackingProtection() {
	resp, err := CDPStorageRunBounceTrackingMitigations()
	if err != nil {
		log.Printf("调试反追踪功能失败: %v", err)
		return
	}
	log.Println("✅ 反弹追踪防护已触发，可检查存储/Cookie清理情况")
}

*/

// -----------------------------------------------  Storage.sendPendingAttributionReports  -----------------------------------------------
// === 应用场景 ===
// 1. 转化归因调试：立即发送浏览器中等待的归因报告（Attribution Reporting）
// 2. 广告转化测试：加速验证广告点击/转化上报流程
// 3. 自动化测试：确保归因报告立即发送，不等待浏览器默认延迟
// 4. 隐私沙箱调试：调试 Attribution Reporting API 数据上报逻辑
// 5. 开发联调：快速触发上报，对接后端转化接收接口
// 6. 问题排查：定位报告未发送、延迟上报异常问题

// CDPStorageSendPendingAttributionReports 立即发送所有挂起的归因报告
func CDPStorageSendPendingAttributionReports() (string, error) {
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
		"method": "Storage.sendPendingAttributionReports"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 sendPendingAttributionReports 请求失败: %w", err)
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
			return "", fmt.Errorf("sendPendingAttributionReports 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：立即发送所有待上报的归因转化报告
func ExampleSendPendingAttributionReports() {
	resp, err := CDPStorageSendPendingAttributionReports()
	if err != nil {
		log.Fatalf("触发归因报告上报失败: %v", err)
	}
	log.Printf("✅ 成功触发挂起归因报告发送:\n%s", resp)
}

// 场景2：自动化测试 - 广告转化归因验证
func ExampleTestAttributionReporting() {
	resp, err := CDPStorageSendPendingAttributionReports()
	if err != nil {
		log.Fatalf("测试归因报告失败: %v", err)
	}
	log.Println("✅ 归因报告已立即发送，可验证后端接收")
}

// 场景3：调试隐私沙箱归因API
func ExampleDebugAttributionReporting() {
	resp, err := CDPStorageSendPendingAttributionReports()
	if err != nil {
		log.Printf("调试归因上报失败: %v", err)
		return
	}
	log.Println("✅ 挂起的归因报告已触发发送")
}

*/

// -----------------------------------------------  Storage.setAttributionReportingLocalTestingMode  -----------------------------------------------
// === 应用场景 ===
// 1. 归因报告本地测试：开启 Attribution Reporting 本地调试模式
// 2. 本地开发联调：允许在 localhost 环境下正常测试归因转化上报
// 3. 自动化测试：禁用浏览器安全限制，稳定验证归因API流程
// 4. 隐私沙箱调试：本地无HTTPS环境调试广告转化归因
// 5. 问题快速复现：跳过浏览器安全校验，专注逻辑验证
// 6. 接口对接测试：本地后端直接接收归因报告

// CDPStorageSetAttributionReportingLocalTestingMode 设置归因报告本地测试模式
// enabled: true 开启本地测试模式，false 关闭
func CDPStorageSetAttributionReportingLocalTestingMode(enabled bool) (string, error) {
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
		"method": "Storage.setAttributionReportingLocalTestingMode",
		"params": {
			"enabled": %t
		}
	}`, reqID, enabled)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setAttributionReportingLocalTestingMode 请求失败: %w", err)
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
			return "", fmt.Errorf("setAttributionReportingTracking 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启归因报告本地测试模式（localhost调试专用）
func ExampleEnableAttributionLocalTesting() {
	resp, err := CDPStorageSetAttributionReportingLocalTestingMode(true)
	if err != nil {
		log.Fatalf("开启本地测试模式失败: %v", err)
	}
	log.Println("✅ 已开启归因报告本地测试模式，可以 localhost 调试")
}

// 场景2：自动化测试启用本地模式
func ExampleTestAttributionLocalMode() {
	resp, err := CDPStorageSetAttributionReportingLocalTestingMode(true)
	if err != nil {
		log.Fatalf("测试设置失败: %v", err)
	}
	log.Println("✅ 自动化测试：归因本地测试模式已启用")
}

// 场景3：关闭本地测试模式，恢复浏览器默认行为
func ExampleDisableAttributionLocalTesting() {
	resp, err := CDPStorageSetAttributionReportingLocalTestingMode(false)
	if err != nil {
		log.Printf("关闭本地测试模式失败: %v", err)
		return
	}
	log.Println("✅ 已关闭归因报告本地测试模式")
}

*/

// -----------------------------------------------  Storage.setAttributionReportingTracking  -----------------------------------------------
// === 应用场景 ===
// 1. 归因报告追踪控制：启用/禁用 Attribution Reporting API 的事件追踪
// 2. 广告转化调试：控制是否捕获归因转化与点击数据
// 3. 自动化测试：按需开关归因追踪，避免测试数据污染
// 4. 隐私沙箱调试：验证归因功能启用/禁用后的行为差异
// 5. 合规测试：验证关闭归因后是否停止上报转化数据
// 6. 性能调试：关闭归因减少浏览器后台请求干扰

// CDPStorageSetAttributionReportingTracking 设置归因报告追踪开关
// enable: true 开启追踪，false 关闭追踪
func CDPStorageSetAttributionReportingTracking(enable bool) (string, error) {
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
		"method": "Storage.setAttributionReportingTracking",
		"params": {
			"enable": %t
		}
	}`, reqID, enable)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setAttributionReportingTracking 请求失败: %w", err)
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
			return "", fmt.Errorf("setAttributionReportingTracking 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启归因报告追踪（默认启用）
func ExampleEnableAttributionTracking() {
	resp, err := CDPStorageSetAttributionReportingTracking(true)
	if err != nil {
		log.Fatalf("开启归因追踪失败: %v", err)
	}
	log.Println("✅ 归因报告追踪已开启")
}

// 场景2：关闭归因报告追踪（测试无归因环境）
func ExampleDisableAttributionTracking() {
	resp, err := CDPStorageSetAttributionReportingTracking(false)
	if err != nil {
		log.Fatalf("关闭归因追踪失败: %v", err)
	}
	log.Println("✅ 归因报告追踪已关闭，浏览器将不再捕获归因数据")
}

// 场景3：自动化测试 - 测试前后控制归因追踪状态
func ExampleTestAttributionTrackingControl() {
	// 测试前开启
	CDPStorageSetAttributionReportingTracking(true)
	// 执行测试逻辑...
	log.Println("✅ 测试中：归因追踪已启用")

	// 测试后关闭
	CDPStorageSetAttributionReportingTracking(false)
	log.Println("✅ 测试完成：归因追踪已关闭")
}

*/

// -----------------------------------------------  Storage.setInterestGroupAuctionTracking  -----------------------------------------------
// === 应用场景 ===
// 1. 广告竞价追踪控制：启用/禁用 Protected Audience API 兴趣组竞价追踪
// 2. 隐私沙箱调试：控制竞价行为日志与数据上报开关
// 3. 自动化测试：避免测试过程中产生广告竞价追踪数据
// 4. 竞价流程调试：开关追踪以验证竞价逻辑是否受影响
// 5. 合规测试：验证关闭追踪后是否停止记录竞价行为
// 6. 性能调试：减少后台追踪开销，专注核心逻辑调试

// CDPStorageSetInterestGroupAuctionTracking 设置兴趣组广告竞价追踪开关
// enable: true 开启追踪，false 关闭追踪
func CDPStorageSetInterestGroupAuctionTracking(enable bool) (string, error) {
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
		"method": "Storage.setInterestGroupAuctionTracking",
		"params": {
			"enable": %t
		}
	}`, reqID, enable)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setInterestGroupAuctionTracking 请求失败: %w", err)
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
			return "", fmt.Errorf("setInterestGroupAuctionTracking 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启兴趣组广告竞价追踪
func ExampleEnableInterestGroupAuctionTracking() {
	resp, err := CDPStorageSetInterestGroupAuctionTracking(true)
	if err != nil {
		log.Fatalf("开启竞价追踪失败: %v", err)
	}
	log.Println("✅ 兴趣组广告竞价追踪已开启")
}

// 场景2：关闭兴趣组广告竞价追踪（调试无追踪环境）
func ExampleDisableInterestGroupAuctionTracking() {
	resp, err := CDPStorageSetInterestGroupAuctionTracking(false)
	if err != nil {
		log.Fatalf("关闭竞价追踪失败: %v", err)
	}
	log.Println("✅ 兴趣组广告竞价追踪已关闭")
}

// 场景3：自动化测试 - 测试前后控制竞价追踪
func ExampleTestAuctionTrackingControl() {
	// 测试前开启
	CDPStorageSetInterestGroupAuctionTracking(true)
	log.Println("✅ 测试中：竞价追踪已启用")

	// 测试完成后关闭，避免数据污染
	CDPStorageSetInterestGroupAuctionTracking(false)
	log.Println("✅ 测试完成：竞价追踪已关闭")
}

*/

// -----------------------------------------------  Storage.setSharedStorageEntry  -----------------------------------------------
// === 应用场景 ===
// 1. 共享存储写入：向指定源的 Shared Storage 写入/更新键值对
// 2. 隐私沙箱调试：手动设置 Shared Storage 数据用于测试
// 3. 自动化测试：预置测试数据，验证业务逻辑
// 4. 开发联调：快速注入配置项，无需前端页面操作
// 5. 状态模拟：模拟用户偏好、实验标记等存储状态
// 6. 数据修复：手动修复错误的共享存储条目

// CDPStorageSetSharedStorageEntry 向 Shared Storage 写入键值对
// ownerOrigin: 站点源地址
// key: 存储键名
// value: 存储值
func CDPStorageSetSharedStorageEntry(ownerOrigin string, key string, value string) (string, error) {
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
		"method": "Storage.setSharedStorageEntry",
		"params": {
			"ownerOrigin": "%s",
			"key": "%s",
			"value": "%s"
		}
	}`, reqID, ownerOrigin, key, value)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setSharedStorageEntry 请求失败: %w", err)
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
			return "", fmt.Errorf("setInterestGroupAuctionTracking 请求超时")
		}
	}
}

/*

// === 使用场景示例代码 ===
// 场景1：向指定网站写入 Shared Storage 键值对
func ExampleSetSharedStorageEntry() {
	origin := "https://example.com"
	key := "user_theme"
	value := "dark"

	resp, err := CDPStorageSetSharedStorageEntry(origin, key, value)
	if err != nil {
		log.Fatalf("写入共享存储失败: %v", err)
	}
	log.Printf("✅ 成功写入 [%s] %s=%s", origin, key, value)
}

// 场景2：自动化测试 - 预置测试用共享存储数据
func ExampleTestSetupSharedStorage() {
	testOrigin := "https://test.example.com"
	resp, err := CDPStorageSetSharedStorageEntry(testOrigin, "test_flag", "active")
	if err != nil {
		log.Fatalf("测试数据写入失败: %v", err)
	}
	log.Println("✅ 测试用 Shared Storage 数据已预置")
}

// 场景3：调试隐私沙箱，手动设置存储状态
func ExampleDebugSetSharedStorage() {
	origin := "https://privacy-sandbox.app"
	resp, err := CDPStorageSetSharedStorageEntry(origin, "experiment", "enabled")
	if err != nil {
		log.Printf("调试设置失败: %v", err)
		return
	}
	log.Println("✅ Shared Storage 状态已设置，可调试功能")
}

*/

// -----------------------------------------------  Storage.setSharedStorageTracking  -----------------------------------------------
// === 应用场景 ===
// 1. 共享存储追踪控制：全局启用/禁用 Shared Storage 追踪与记录
// 2. 隐私沙箱调试：控制共享存储读写行为的日志与上报
// 3. 自动化测试：避免测试产生额外存储追踪数据
// 4. 合规测试：验证关闭后是否停止存储用户隐私数据
// 5. 本地开发调试：减少干扰，专注核心逻辑调试
// 6. 性能优化：关闭追踪减少浏览器后台开销

// CDPStorageSetSharedStorageTracking 设置 Shared Storage 全局追踪开关
// enable: true 开启追踪，false 关闭追踪
func CDPStorageSetSharedStorageTracking(enable bool) (string, error) {
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
		"method": "Storage.setSharedStorageTracking",
		"params": {
			"enable": %t
		}
	}`, reqID, enable)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setSharedStorageTracking 请求失败: %w", err)
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
			return "", fmt.Errorf("setSharedStorageTracking 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启 Shared Storage 追踪
func ExampleEnableSharedStorageTracking() {
	resp, err := CDPStorageSetSharedStorageTracking(true)
	if err != nil {
		log.Fatalf("开启共享存储追踪失败: %v", err)
	}
	log.Println("✅ Shared Storage 追踪已开启")
}

// 场景2：关闭 Shared Storage 追踪（隐私/调试模式）
func ExampleDisableSharedStorageTracking() {
	resp, err := CDPStorageSetSharedStorageTracking(false)
	if err != nil {
		log.Fatalf("关闭共享存储追踪失败: %v", err)
	}
	log.Println("✅ Shared Storage 追踪已关闭")
}

// 场景3：自动化测试 - 测试前后控制追踪状态
func ExampleTestSharedStorageTrackingControl() {
	// 测试前开启
	CDPStorageSetSharedStorageTracking(true)
	log.Println("✅ 测试中：Shared Storage 追踪已启用")

	// 测试后关闭
	CDPStorageSetSharedStorageTracking(false)
	log.Println("✅ 测试完成：Shared Storage 追踪已关闭")
}

*/

// -----------------------------------------------  Storage.setStorageBucketTracking  -----------------------------------------------
// === 应用场景 ===
// 1. 存储桶追踪控制：全局启用/禁用 Storage Bucket 存储行为追踪
// 2. 分区存储调试：监控分桶存储的创建、读写、删除行为
// 3. 自动化测试：避免测试过程产生存储桶追踪数据
// 4. 隐私合规测试：验证关闭后是否停止记录存储桶操作
// 5. 存储隔离调试：开关追踪以验证分区存储逻辑
// 6. 本地开发：减少后台数据干扰，专注核心调试

// CDPStorageSetStorageBucketTracking 设置存储桶操作追踪开关
// enable: true 开启追踪，false 关闭追踪
func CDPStorageSetStorageBucketTracking(enable bool) (string, error) {
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
		"method": "Storage.setStorageBucketTracking",
		"params": {
			"enable": %t
		}
	}`, reqID, enable)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setStorageBucketTracking 请求失败: %w", err)
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
			return "", fmt.Errorf("setStorageBucketTracking 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：开启存储桶操作追踪
func ExampleEnableStorageBucketTracking() {
	resp, err := CDPStorageSetStorageBucketTracking(true)
	if err != nil {
		log.Fatalf("开启存储桶追踪失败: %v", err)
	}
	log.Println("✅ 存储桶追踪已开启")
}

// 场景2：关闭存储桶操作追踪
func ExampleDisableStorageBucketTracking() {
	resp, err := CDPStorageSetStorageBucketTracking(false)
	if err != nil {
		log.Fatalf("关闭存储桶追踪失败: %v", err)
	}
	log.Println("✅ 存储桶追踪已关闭")
}

// 场景3：自动化测试 - 测试前后控制追踪状态
func ExampleTestStorageBucketTrackingControl() {
	// 测试前开启
	CDPStorageSetStorageBucketTracking(true)
	log.Println("✅ 测试中：存储桶追踪已启用")

	// 测试后关闭
	CDPStorageSetStorageBucketTracking(false)
	log.Println("✅ 测试完成：存储桶追踪已关闭")
}

*/
