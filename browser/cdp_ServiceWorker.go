package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  ServiceWorker.deliverPushMessage  -----------------------------------------------
// === 应用场景 ===
// 1. 推送测试: 自动化测试Web推送通知功能，无需真实后端推送服务
// 2. 调试验证: 本地调试ServiceWorker接收推送消息的逻辑是否正常执行
// 3. 离线测试: 无网络环境下测试推送消息的处理流程
// 4. 功能开发: 前端开发推送功能时快速模拟推送触发
// 5. 异常测试: 模拟不同格式/内容的推送消息测试兼容性
// 6. 自动化验收: 端到端测试中触发推送验证业务流程

// CDPServiceWorkerDeliverPushMessage 向ServiceWorker投递推送消息
// versionId: ServiceWorker版本ID
// data: 推送消息数据（字符串格式）
// timeToLive: 推送消息存活时间（秒）
func CDPServiceWorkerDeliverPushMessage(versionId string, data string, timeToLive int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"versionId": "%s", "data": "%s", "timeToLive": %d`, versionId, data, timeToLive)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.deliverPushMessage",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deliverPushMessage 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("deliverPushMessage 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：基础推送测试 - 向指定ServiceWorker发送文本推送消息
func ExampleDeliverPushMessage_Basic() {
	// 替换为你的ServiceWorker版本ID（可通过ServiceWorker.getRegistration获取）
	swVersionID := "SWVersionID_123456"
	// 推送消息内容
	pushData := "{'title':'测试通知','body':'这是一条测试推送消息'}"
	// 存活时间：3600秒
	ttl := 3600

	resp, err := CDPServiceWorkerDeliverPushMessage(swVersionID, pushData, ttl)
	if err != nil {
		log.Fatalf("投递推送消息失败: %v", err)
	}
	log.Printf("投递成功: %s", resp)
}

// 场景2：调试推送逻辑 - 空数据测试异常处理
func ExampleDeliverPushMessage_Debug() {
	swVersionID := "SWVersionID_123456"
	// 空数据测试
	pushData := ""
	ttl := 60

	resp, err := CDPServiceWorkerDeliverPushMessage(swVersionID, pushData, ttl)
	if err != nil {
		log.Printf("预期错误: %v", err)
	} else {
		log.Printf("响应: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 测试隔离: 禁用ServiceWorker避免缓存/拦截影响测试结果
// 2. 调试排查: 临时关闭ServiceWorker定位页面加载问题
// 3. 环境重置: 自动化测试后禁用ServiceWorker恢复初始状态
// 4. 性能调试: 关闭后排查ServiceWorker导致的性能问题
// 5. 兼容性测试: 对比开启/禁用状态的页面行为差异
// 6. 清理拦截: 停止ServiceWorker的网络请求拦截功能

// CDPServiceWorkerDisable 禁用ServiceWorker功能
func CDPServiceWorkerDisable() (string, error) {
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
		"method": "ServiceWorker.disable"
	}`, reqID)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 disable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 解析响应数据
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
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


// === 使用场景示例代码 ===
// 场景1：自动化测试前禁用ServiceWorker
func ExampleServiceWorkerDisable_Test() {
	resp, err := CDPServiceWorkerDisable()
	if err != nil {
		log.Fatalf("禁用ServiceWorker失败: %v", err)
	}
	log.Printf("禁用成功: %s", resp)
}

// 场景2：调试页面时临时禁用
func ExampleServiceWorkerDisable_Debug() {
	resp, err := CDPServiceWorkerDisable()
	if err != nil {
		log.Printf("禁用失败: %v", err)
		return
	}
	log.Println("已禁用ServiceWorker，可正常调试页面")
}

*/

// -----------------------------------------------  ServiceWorker.dispatchPeriodicSyncEvent  -----------------------------------------------
// === 应用场景 ===
// 1. 定时同步测试: 模拟ServiceWorker周期性后台同步事件
// 2. 离线数据同步: 测试定时同步本地数据到服务器的逻辑
// 3. 自动化验证: 验证定时同步功能是否正常触发执行
// 4. 调试同步流程: 本地调试Periodic Sync事件处理逻辑
// 5. 边界测试: 模拟频繁/异常定时同步触发测试稳定性
// 6. 业务流程测试: 测试依赖定时同步的业务功能

// CDPServiceWorkerDispatchPeriodicSyncEvent 触发ServiceWorker周期性同步事件
// origin: 注册ServiceWorker的源地址（如：https://example.com）
// registrationId: ServiceWorker注册ID
// tag: 定时同步的标签标识
func CDPServiceWorkerDispatchPeriodicSyncEvent(origin string, registrationId string, tag string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"origin": "%s", "registrationId": "%s", "tag": "%s"`, origin, registrationId, tag)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.dispatchPeriodicSyncEvent",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 dispatchPeriodicSyncEvent 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("dispatchPeriodicSyncEvent 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：基础定时同步测试 - 模拟触发标准周期性同步事件
func ExampleDispatchPeriodicSyncEvent_Basic() {
	// 站点源地址
	origin := "https://localhost:8080"
	// ServiceWorker注册ID（通过getRegistration获取）
	regID := "SW_REG_ID_123"
	// 定时同步标签
	syncTag := "data-sync-tag"

	resp, err := CDPServiceWorkerDispatchPeriodicSyncEvent(origin, regID, syncTag)
	if err != nil {
		log.Fatalf("触发定时同步失败: %v", err)
	}
	log.Printf("触发同步成功: %s", resp)
}

// 场景2：调试定时同步逻辑 - 测试自定义标签同步
func ExampleDispatchPeriodicSyncEvent_Debug() {
	origin := "https://localhost:8080"
	regID := "SW_REG_ID_123"
	// 自定义调试标签
	syncTag := "debug-sync-test"

	resp, err := CDPServiceWorkerDispatchPeriodicSyncEvent(origin, regID, syncTag)
	if err != nil {
		log.Printf("调试触发失败: %v", err)
	} else {
		log.Printf("调试同步事件已触发: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.dispatchSyncEvent  -----------------------------------------------
// === 应用场景 ===
// 1. 离线同步测试: 模拟ServiceWorker后台同步事件，测试离线数据上传逻辑
// 2. 自动化测试: 自动化验证同步功能是否正常触发
// 3. 本地调试: 无需等待浏览器自动触发，手动调试Sync事件处理流程
// 4. 异常测试: 模拟重复触发、异常标签测试同步稳定性
// 5. 业务验证: 测试依赖后台同步的业务流程
// 6. 缓存更新: 手动触发同步更新本地缓存数据

// CDPServiceWorkerDispatchSyncEvent 触发ServiceWorker同步事件
// origin: 注册ServiceWorker的源地址（例如：https://example.com）
// registrationId: ServiceWorker注册ID
// tag: 同步事件的标签标识
func CDPServiceWorkerDispatchSyncEvent(origin string, registrationId string, tag string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"origin": "%s", "registrationId": "%s", "tag": "%s"`, origin, registrationId, tag)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.dispatchSyncEvent",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 dispatchSyncEvent 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("dispatchSyncEvent 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：基础同步测试 - 模拟触发标准同步事件上传离线数据
func ExampleDispatchSyncEvent_Basic() {
	origin := "https://localhost:8080"
	regID := "SW_REG_123456"
	syncTag := "upload-offline-data"

	resp, err := CDPServiceWorkerDispatchSyncEvent(origin, regID, syncTag)
	if err != nil {
		log.Fatalf("触发同步事件失败: %v", err)
	}
	log.Printf("触发同步成功: %s", resp)
}

// 场景2：调试同步逻辑 - 自定义标签测试
func ExampleDispatchSyncEvent_Debug() {
	origin := "https://localhost:8080"
	regID := "SW_REG_123456"
	syncTag := "debug-sync-test"

	resp, err := CDPServiceWorkerDispatchSyncEvent(origin, regID, syncTag)
	if err != nil {
		log.Printf("调试触发失败: %v", err)
	} else {
		log.Printf("调试同步事件已触发: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 测试初始化: 自动化测试开始前启用ServiceWorker功能
// 2. 调试恢复: 调试完成后重新启用ServiceWorker恢复正常功能
// 3. 环境准备: 测试PWA应用前初始化ServiceWorker环境
// 4. 功能切换: 从禁用状态切换回启用ServiceWorker
// 5. 缓存启动: 开启ServiceWorker的网络缓存与拦截功能
// 6. 自动化流程: 端到端测试中启动ServiceWorker服务

// CDPServiceWorkerEnable 启用ServiceWorker功能
func CDPServiceWorkerEnable() (string, error) {
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
		"method": "ServiceWorker.enable"
	}`, reqID)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 enable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 解析响应数据
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
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


// === 使用场景示例代码 ===
// 场景1：自动化测试开始前启用ServiceWorker
func ExampleServiceWorkerEnable_Test() {
	resp, err := CDPServiceWorkerEnable()
	if err != nil {
		log.Fatalf("启用ServiceWorker失败: %v", err)
	}
	log.Printf("启用成功: %s", resp)
}

// 场景2：调试完成后恢复ServiceWorker功能
func ExampleServiceWorkerEnable_Debug() {
	resp, err := CDPServiceWorkerEnable()
	if err != nil {
		log.Printf("启用失败: %v", err)
		return
	}
	log.Println("已启用ServiceWorker，恢复正常缓存与功能")
}

*/

// -----------------------------------------------  ServiceWorker.setForceUpdateOnPageLoad  -----------------------------------------------
// === 应用场景 ===
// 1. 自动化测试: 强制页面加载时更新ServiceWorker，确保使用最新脚本
// 2. 调试开发: 本地开发时避免SW缓存，每次加载自动更新
// 3. 版本发布: 测试新版本ServiceWorker强制更新逻辑
// 4. 缓存清理: 确保页面加载时跳过缓存直接更新SW
// 5. 回归测试: 验证ServiceWorker热更新功能
// 6. 环境重置: 测试环境强制刷新SW状态

// CDPServiceWorkerSetForceUpdateOnPageLoad 设置页面加载时强制更新ServiceWorker
// forceUpdate: true=开启强制更新，false=关闭强制更新
func CDPServiceWorkerSetForceUpdateOnPageLoad(forceUpdate bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"forceUpdateOnPageLoad": %t`, forceUpdate)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.setForceUpdateOnPageLoad",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setForceUpdateOnPageLoad 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("setForceUpdateOnPageLoad 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：自动化测试开启强制更新
func ExampleSetForceUpdateOnPageLoad_Enable() {
	// 强制页面加载时更新ServiceWorker
	resp, err := CDPServiceWorkerSetForceUpdateOnPageLoad(true)
	if err != nil {
		log.Fatalf("设置强制更新失败: %v", err)
	}
	log.Printf("设置成功: %s", resp)
}

// 场景2：关闭强制更新，恢复默认行为
func ExampleSetForceUpdateOnPageLoad_Disable() {
	// 关闭强制更新，使用浏览器默认SW更新策略
	resp, err := CDPServiceWorkerSetForceUpdateOnPageLoad(false)
	if err != nil {
		log.Printf("关闭强制更新失败: %v", err)
	} else {
		log.Printf("已关闭强制更新: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.skipWaiting  -----------------------------------------------
// === 应用场景 ===
// 1. PWA版本更新: 强制新ServiceWorker立即激活，无需等待页面刷新
// 2. 自动化测试: 测试SW更新激活流程，确保新版本快速生效
// 3. 调试开发: 本地调试时快速激活新ServiceWorker，验证更新逻辑
// 4. 缓存更新: 强制激活新SW以更新离线缓存和资源
// 5. 版本切换: 快速切换ServiceWorker版本，验证功能兼容性
// 6. 测试环境: 自动化测试中确保SW立即激活，不等待页面关闭

// CDPServiceWorkerSkipWaiting 强制ServiceWorker跳过等待状态并激活
// versionId: ServiceWorker版本ID
func CDPServiceWorkerSkipWaiting(versionId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"versionId": "%s"`, versionId)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.skipWaiting",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 skipWaiting 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("skipWaiting 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：PWA更新测试 - 强制新ServiceWorker立即激活
func ExampleSkipWaiting_Activate() {
	// 替换为你的ServiceWorker版本ID（通过getRegistration获取）
	swVersionID := "SW_VERSION_123456"

	resp, err := CDPServiceWorkerSkipWaiting(swVersionID)
	if err != nil {
		log.Fatalf("强制激活ServiceWorker失败: %v", err)
	}
	log.Printf("强制激活成功: %s", resp)
}

// 场景2：调试SW更新逻辑 - 本地开发快速激活新版本
func ExampleSkipWaiting_Debug() {
	swVersionID := "SW_VERSION_123456"

	resp, err := CDPServiceWorkerSkipWaiting(swVersionID)
	if err != nil {
		log.Printf("调试激活失败: %v", err)
	} else {
		log.Printf("新ServiceWorker已激活: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.startWorker  -----------------------------------------------
// === 应用场景 ===
// 1. 服务启动测试: 手动启动ServiceWorker，测试初始化逻辑
// 2. 崩溃恢复测试: 模拟SW崩溃后重新启动，验证恢复能力
// 3. 调试启动流程: 本地调试ServiceWorker的启动、安装流程
// 4. 自动化初始化: 测试流程中主动启动SW，确保服务就绪
// 5. 离线功能启动: 手动触发SW启动，保障离线缓存功能生效
// 6. 边界场景测试: 测试重复启动、异常启动的稳定性

// CDPServiceWorkerStartWorker 手动启动指定的ServiceWorker
// versionId: ServiceWorker 版本ID（通过getRegistration获取）
func CDPServiceWorkerStartWorker(versionId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"versionId": "%s"`, versionId)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.startWorker",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 startWorker 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("startWorker 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：自动化测试 - 主动启动ServiceWorker保证服务可用
func ExampleStartWorker_Basic() {
	// 替换为真实的ServiceWorker版本ID
	swVersionID := "SW_VERSION_98765"

	resp, err := CDPServiceWorkerStartWorker(swVersionID)
	if err != nil {
		log.Fatalf("启动ServiceWorker失败: %v", err)
	}
	log.Printf("启动成功: %s", resp)
}

// 场景2：调试崩溃恢复 - 模拟崩溃后重新启动
func ExampleStartWorker_DebugRecovery() {
	swVersionID := "SW_VERSION_98765"

	resp, err := CDPServiceWorkerStartWorker(swVersionID)
	if err != nil {
		log.Printf("重启失败: %v", err)
	} else {
		log.Printf("已重新启动ServiceWorker: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.stopAllWorkers  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境清理: 自动化测试完成后停止所有ServiceWorker，重置环境
// 2. 调试问题定位: 停止所有SW进程，排查页面加载/缓存问题
// 3. 资源释放: 释放ServiceWorker占用的浏览器资源
// 4. 安全关闭: 批量终止SW服务，避免后台持续运行
// 5. 状态重置: 重置所有ServiceWorker运行状态，恢复初始环境
// 6. 崩溃恢复: 异常情况下批量关闭SW，防止连锁错误

// CDPServiceWorkerStopAllWorkers 停止所有运行中的ServiceWorker
func CDPServiceWorkerStopAllWorkers() (string, error) {
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
		"method": "ServiceWorker.stopAllWorkers"
	}`, reqID)

	// 发送WebSocket请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stopAllWorkers 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 解析响应数据
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP协议错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("stopAllWorkers 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：自动化测试结束后清理所有ServiceWorker
func ExampleStopAllWorkers_TestCleanup() {
	resp, err := CDPServiceWorkerStopAllWorkers()
	if err != nil {
		log.Fatalf("停止所有ServiceWorker失败: %v", err)
	}
	log.Printf("停止成功: %s", resp)
}

// 场景2：调试时批量关闭ServiceWorker排查问题
func ExampleStopAllWorkers_Debug() {
	resp, err := CDPServiceWorkerStopAllWorkers()
	if err != nil {
		log.Printf("关闭失败: %v", err)
		return
	}
	log.Println("已关闭所有ServiceWorker，可进行页面调试")
}

*/

// -----------------------------------------------  ServiceWorker.stopWorker  -----------------------------------------------
// === 应用场景 ===
// 1. 精准测试: 停止指定版本的ServiceWorker，不影响其他服务
// 2. 问题定位: 单独关闭异常SW，排查单个服务故障
// 3. 资源管理: 释放单个SW占用的浏览器资源
// 4. 版本切换: 停止旧版SW，为新版激活做准备
// 5. 调试隔离: 单独控制某个SW的运行状态
// 6. 崩溃恢复: 关闭异常SW后重新启动

// CDPServiceWorkerStopWorker 停止指定的ServiceWorker
// versionId: ServiceWorker版本ID
func CDPServiceWorkerStopWorker(versionId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"versionId": "%s"`, versionId)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.stopWorker",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 stopWorker 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("stopWorker 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：停止指定版本的ServiceWorker
func ExampleStopWorker_Basic() {
	// 替换为你的ServiceWorker版本ID
	swVersionID := "SW_VERSION_123456"

	resp, err := CDPServiceWorkerStopWorker(swVersionID)
	if err != nil {
		log.Fatalf("停止ServiceWorker失败: %v", err)
	}
	log.Printf("停止成功: %s", resp)
}

// 场景2：调试异常ServiceWorker - 单独关闭故障服务
func ExampleStopWorker_Debug() {
	swVersionID := "SW_VERSION_123456"

	resp, err := CDPServiceWorkerStopWorker(swVersionID)
	if err != nil {
		log.Printf("关闭失败: %v", err)
	} else {
		log.Printf("已关闭异常ServiceWorker: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.unregister  -----------------------------------------------
// === 应用场景 ===
// 1. 测试环境清理: 自动化测试后注销ServiceWorker注册，还原环境
// 2. 调试问题定位: 注销SW后重新注册，解决缓存/更新异常
// 3. 功能重置: 完全清除ServiceWorker，恢复普通网页状态
// 4. 版本切换: 注销旧SW注册，强制重新安装新版本
// 5. 安全清理: 注销无用SW，释放浏览器存储与缓存资源
// 6. 离线功能关闭: 关闭PWA离线功能，停止SW服务

// CDPServiceWorkerUnregister 注销指定的ServiceWorker注册
// scopeURL: ServiceWorker注册的作用域URL（如：https://localhost:8080/）
func CDPServiceWorkerUnregister(scopeURL string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"scopeURL": "%s"`, scopeURL)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.unregister",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 unregister 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("unregister 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：自动化测试后注销ServiceWorker，清理环境
func ExampleUnregister_TestCleanup() {
	// 替换为你的SW作用域URL
	scopeURL := "https://localhost:8080/"

	resp, err := CDPServiceWorkerUnregister(scopeURL)
	if err != nil {
		log.Fatalf("注销ServiceWorker失败: %v", err)
	}
	log.Printf("注销成功: %s", resp)
}

// 场景2：调试时注销SW解决缓存异常
func ExampleUnregister_Debug() {
	scopeURL := "https://localhost:8080/"

	resp, err := CDPServiceWorkerUnregister(scopeURL)
	if err != nil {
		log.Printf("注销失败: %v", err)
	} else {
		log.Printf("已注销ServiceWorker，可重新注册: %s", resp)
	}
}

*/

// -----------------------------------------------  ServiceWorker.updateRegistration  -----------------------------------------------
// === 应用场景 ===
// 1. 强制更新测试: 手动触发ServiceWorker注册更新，测试新版本部署
// 2. 自动化验证: 验证SW更新流程是否正常执行
// 3. 调试更新逻辑: 本地调试时主动触发更新，无需等待浏览器默认检查周期
// 4. 缓存刷新: 强制更新SW以刷新静态资源缓存
// 5. 版本回归测试: 测试从旧版本升级到新版本的完整流程
// 6. PWA发布验证: 验证生产环境SW热更新功能

// CDPServiceWorkerUpdateRegistration 手动更新指定的ServiceWorker注册
// scopeURL: ServiceWorker注册的作用域URL
func CDPServiceWorkerUpdateRegistration(scopeURL string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求参数
	params := fmt.Sprintf(`"scopeURL": "%s"`, scopeURL)
	// 构建完整消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "ServiceWorker.updateRegistration",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 updateRegistration 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应，超时5秒
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

				// 检查CDP错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("updateRegistration 请求超时")
		}
	}
}

/*


// === 使用场景示例代码 ===
// 场景1：自动化测试 - 强制触发ServiceWorker更新检查
func ExampleUpdateRegistration_ForceUpdate() {
	// 替换为你的ServiceWorker作用域URL
	scopeURL := "https://localhost:8080/"

	resp, err := CDPServiceWorkerUpdateRegistration(scopeURL)
	if err != nil {
		log.Fatalf("触发ServiceWorker更新失败: %v", err)
	}
	log.Printf("更新触发成功: %s", resp)
}

// 场景2：调试更新流程 - 本地开发手动刷新SW版本
func ExampleUpdateRegistration_Debug() {
	scopeURL := "https://localhost:8080/"

	resp, err := CDPServiceWorkerUpdateRegistration(scopeURL)
	if err != nil {
		log.Printf("调试更新失败: %v", err)
	} else {
		log.Printf("已手动触发ServiceWorker更新: %s", resp)
	}
}


*/
