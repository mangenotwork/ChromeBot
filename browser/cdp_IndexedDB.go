package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  IndexedDB.clearObjectStore  -----------------------------------------------
// === 应用场景 ===
// 1. 数据清理: 清空特定对象存储中的所有数据
// 2. 测试环境准备: 在自动化测试前清空测试数据库
// 3. 用户数据重置: 允许用户清空特定类型的数据
// 4. 缓存清理: 清空IndexedDB中存储的缓存数据
// 5. 数据库迁移: 在数据迁移前清空旧的数据存储
// 6. 隐私保护: 清除用户的敏感数据

// CDPIndexedDBClearObjectStore 清空IndexedDB对象存储
func CDPIndexedDBClearObjectStore(databaseName, objectStoreName string) (string, error) {
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
		"method": "IndexedDB.clearObjectStore",
		"params": {
			"securityOrigin": "",
			"databaseName": "%s",
			"objectStoreName": "%s"
		}
	}`, reqID, databaseName, objectStoreName)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearObjectStore 请求失败: %w", err)
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
			return "", fmt.Errorf("clearObjectStore 请求超时")
		}
	}
}

/*

// 场景1: 清空用户会话数据
// 在用户登出时清空会话相关的IndexedDB存储
func clearUserSessionData() error {
	result, err := CDPIndexedDBClearObjectStore("myAppDB", "userSessions")
	if err != nil {
		log.Printf("清空用户会话数据失败: %v", err)
		return err
	}
	log.Printf("用户会话数据已清空: %s", result)
	return nil
}

// 场景2: 自动化测试前清理
// 在运行E2E测试前清理测试数据库
func setupTestEnvironment() error {
	// 清空测试数据存储
	result, err := CDPIndexedDBClearObjectStore("testDB", "testData")
	if err != nil {
		log.Printf("清理测试数据失败: %v", err)
		return err
	}
	log.Printf("测试环境已准备: %s", result)
	return nil
}


// 场景3: 清理应用缓存
// 定期清理IndexedDB中存储的临时缓存数据
func clearAppCache() error {
	// 清空缓存对象存储
	result, err := CDPIndexedDBClearObjectStore("appCacheDB", "imageCache")
	if err != nil {
		log.Printf("清理图片缓存失败: %v", err)
		return err
	}
	log.Printf("图片缓存已清理: %s", result)
	return nil
}

*/

// -----------------------------------------------  IndexedDB.deleteDatabase  -----------------------------------------------
// === 应用场景 ===
// 1. 数据库迁移: 删除旧版本的数据库以便创建新版本
// 2. 用户数据删除: 允许用户删除整个应用数据库
// 3. 隐私清理: 在用户登出或删除账号时清理所有本地数据
// 4. 测试环境清理: 在自动化测试完成后删除测试数据库
// 5. 错误恢复: 在数据库损坏时删除并重建数据库
// 6. 版本管理: 在不同应用版本间切换时清理旧数据库

// CDPIndexedDBDeleteDatabase 删除IndexedDB数据库
func CDPIndexedDBDeleteDatabase(databaseName string) (string, error) {
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
		"method": "IndexedDB.deleteDatabase",
		"params": {
			"securityOrigin": "",
			"databaseName": "%s"
		}
	}`, reqID, databaseName)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deleteDatabase 请求失败: %w", err)
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
			return "", fmt.Errorf("deleteDatabase 请求超时")
		}
	}
}

/*

// 场景1: 应用数据重置
// 用户点击"重置应用数据"时删除所有本地存储
func resetApplicationData() error {
	// 删除主应用数据库
	result, err := CDPIndexedDBDeleteDatabase("myAppDB")
	if err != nil {
		log.Printf("删除应用数据库失败: %v", err)
		return err
	}
	log.Printf("应用数据库已删除: %s", result)
	return nil
}

// 场景2: 版本升级时数据库迁移
// 应用从v1升级到v2时删除旧版本数据库
func migrateDatabaseFromV1ToV2() error {
	// 删除旧版本数据库
	result, err := CDPIndexedDBDeleteDatabase("myAppDB_v1")
	if err != nil {
		log.Printf("删除旧版本数据库失败: %v", err)
		return err
	}
	log.Printf("旧版本数据库已删除: %s", result)

	// 创建新版本数据库
	// ... 创建v2数据库的逻辑
	return nil
}

// 场景3: 用户登出时清理数据
// 用户退出登录时删除所有本地存储的数据
func cleanupOnUserLogout() error {
	// 删除用户相关的数据库
	databases := []string{"userDataDB", "preferencesDB", "cacheDB"}

	for _, dbName := range databases {
		result, err := CDPIndexedDBDeleteDatabase(dbName)
		if err != nil {
			log.Printf("删除数据库 %s 失败: %v", dbName, err)
			continue
		}
		log.Printf("数据库 %s 已删除: %s", dbName, result)
	}
	return nil
}


*/

// -----------------------------------------------  IndexedDB.deleteObjectStoreEntries  -----------------------------------------------
// === 应用场景 ===
// 1. 批量数据删除: 根据键范围批量删除对象存储中的记录
// 2. 时间范围清理: 删除特定时间范围内的历史记录
// 3. 数据分区删除: 根据特定条件删除部分数据而非全部
// 4. 过期数据清理: 清理过期的缓存或会话数据
// 5. 数据归档: 在归档前删除已处理的数据
// 6. 条件性清理: 根据业务规则清理特定数据子集

// CDPIndexedDBDeleteObjectStoreEntries 删除对象存储中的条目
func CDPIndexedDBDeleteObjectStoreEntries(databaseName, objectStoreName string, keyRange KeyRange) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 将KeyRange转换为JSON
	keyRangeJSON, err := json.Marshal(keyRange)
	if err != nil {
		return "", fmt.Errorf("序列化keyRange失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "IndexedDB.deleteObjectStoreEntries",
		"params": {
			"securityOrigin": "",
			"databaseName": "%s",
			"objectStoreName": "%s",
			"keyRange": %s
		}
	}`, reqID, databaseName, objectStoreName, string(keyRangeJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deleteObjectStoreEntries 请求失败: %w", err)
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
			return "", fmt.Errorf("deleteObjectStoreEntries 请求超时")
		}
	}
}

// KeyRange 表示IndexedDB的键范围
type KeyRange struct {
	Lower     interface{} `json:"lower,omitempty"`
	Upper     interface{} `json:"upper,omitempty"`
	LowerOpen bool        `json:"lowerOpen,omitempty"`
	UpperOpen bool        `json:"upperOpen,omitempty"`
}

/*

// 场景1: 删除过期的会话数据
// 清理一周前的用户会话记录
func cleanupExpiredSessions() error {
	// 计算一周前的时间戳
	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour).UnixNano() / 1e6 // 转换为毫秒

	// 创建键范围：删除所有timestamp小于一周前的记录
	keyRange := KeyRange{
		Upper:     oneWeekAgo,
		UpperOpen: true, // 不包含上限
	}

	result, err := CDPIndexedDBDeleteObjectStoreEntries("sessionDB", "userSessions", keyRange)
	if err != nil {
		log.Printf("清理过期会话失败: %v", err)
		return err
	}
	log.Printf("已清理过期会话: %s", result)
	return nil
}

// 场景2: 批量删除特定用户的数据
// 删除ID在1000-2000范围内的用户数据
func deleteUsersInRange() error {
	keyRange := KeyRange{
		Lower:     1000,
		Upper:     2000,
		LowerOpen: false, // 包含下限
		UpperOpen: false, // 包含上限
	}

	result, err := CDPIndexedDBDeleteObjectStoreEntries("userDB", "profiles", keyRange)
	if err != nil {
		log.Printf("批量删除用户失败: %v", err)
		return err
	}
	log.Printf("已批量删除用户: %s", result)
	return nil
}

// 场景3: 清理特定前缀的缓存数据
// 删除以"temp_"开头的临时缓存
func cleanupTempCache() error {
	keyRange := KeyRange{
		Lower:     "temp_",
		Upper:     "temp_" + string([]byte{0xFF}), // 确保包含所有以temp_开头的键
		LowerOpen: false,
		UpperOpen: true,
	}

	result, err := CDPIndexedDBDeleteObjectStoreEntries("cacheDB", "resources", keyRange)
	if err != nil {
		log.Printf("清理临时缓存失败: %v", err)
		return err
	}
	log.Printf("已清理临时缓存: %s", result)
	return nil
}

*/

// -----------------------------------------------  IndexedDB.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 性能优化: 禁用IndexedDB以提升页面性能
// 2. 调试工具关闭: 在调试完成后关闭IndexedDB观察功能
// 3. 资源清理: 在页面卸载前释放IndexedDB相关资源
// 4. 测试环境控制: 控制测试环境中IndexedDB的可用性
// 5. 功能切换: 根据需要动态禁用IndexedDB功能
// 6. 错误处理: 在IndexedDB发生错误时临时禁用

// CDPIndexedDBDisable 禁用IndexedDB域
func CDPIndexedDBDisable() (string, error) {
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
		"method": "IndexedDB.disable"
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

// 场景1: 性能敏感页面优化
// 在性能要求高的页面禁用IndexedDB以减少开销
func optimizePagePerformance() error {
	// 禁用IndexedDB以减少性能开销
	result, err := CDPIndexedDBDisable()
	if err != nil {
		log.Printf("禁用IndexedDB失败: %v", err)
		return err
	}
	log.Printf("IndexedDB已禁用，性能优化完成: %s", result)
	return nil
}


// 场景2: 调试完成后清理
// 在IndexedDB调试会话结束后禁用相关功能
func endIndexedDBDebugSession() error {
	// 先完成调试操作...

	// 然后禁用IndexedDB域
	result, err := CDPIndexedDBDisable()
	if err != nil {
		log.Printf("结束调试会话失败: %v", err)
		return err
	}
	log.Printf("调试会话已结束，IndexedDB已禁用: %s", result)
	return nil
}

// 场景3: 页面卸载前资源清理
// 在页面即将卸载时禁用IndexedDB相关功能
func cleanupBeforePageUnload() error {
	// 禁用IndexedDB以释放资源
	result, err := CDPIndexedDBDisable()
	if err != nil {
		log.Printf("页面卸载前清理失败: %v", err)
		return err
	}
	log.Printf("IndexedDB资源已释放: %s", result)
	return nil
}

*/

// -----------------------------------------------  IndexedDB.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 调试工具启动: 启用IndexedDB调试功能
// 2. 数据库操作前准备: 在执行IndexedDB操作前启用域
// 3. 功能恢复: 在禁用后重新启用IndexedDB功能
// 4. 自动化测试: 在测试前启用IndexedDB支持
// 5. 监控启用: 启用IndexedDB事件监控
// 6. 开发工具集成: 在开发工具中启用IndexedDB支持

// CDPIndexedDBEnable 启用IndexedDB域
func CDPIndexedDBEnable() (string, error) {
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
		"method": "IndexedDB.enable"
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

// 场景1: 开启IndexedDB调试功能
// 在开发调试时启用IndexedDB域以便进行调试
func startIndexedDBDebugging() error {
	// 启用IndexedDB域
	result, err := CDPIndexedDBEnable()
	if err != nil {
		log.Printf("启用IndexedDB调试失败: %v", err)
		return err
	}
	log.Printf("IndexedDB调试已启用: %s", result)

	// 现在可以执行其他IndexedDB调试操作
	return nil
}

// 场景2: 自动化测试初始化
// 在自动化测试开始时启用IndexedDB支持
func setupIndexedDBForTesting() error {
	// 启用IndexedDB域
	result, err := CDPIndexedDBEnable()
	if err != nil {
		log.Printf("测试环境IndexedDB初始化失败: %v", err)
		return err
	}
	log.Printf("IndexedDB测试环境已就绪: %s", result)

	// 可以继续执行测试相关的IndexedDB操作
	return nil
}

// 场景3: 功能恢复和重新启用
// 在临时禁用后重新启用IndexedDB功能
func restoreIndexedDBFunctionality() error {
	// 重新启用IndexedDB域
	result, err := CDPIndexedDBEnable()
	if err != nil {
		log.Printf("重新启用IndexedDB失败: %v", err)
		return err
	}
	log.Printf("IndexedDB功能已恢复: %s", result)

	// 现在可以正常使用IndexedDB功能
	return nil
}

*/

// -----------------------------------------------  IndexedDB.getMetadata  -----------------------------------------------
// === 应用场景 ===
// 1. 数据库状态检查: 获取数据库的元数据信息
// 2. 空间使用分析: 分析数据库大小和存储使用情况
// 3. 性能监控: 监控IndexedDB的存储容量
// 4. 容量规划: 根据数据库大小进行存储规划
// 5. 清理决策: 根据存储使用情况决定是否清理数据
// 6. 调试诊断: 获取元数据用于调试数据库问题

// CDPIndexedDBGetMetadata 获取IndexedDB数据库的元数据
func CDPIndexedDBGetMetadata(databaseName string) (string, error) {
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
		"method": "IndexedDB.getMetadata",
		"params": {
			"securityOrigin": "",
			"databaseName": "%s"
		}
	}`, reqID, databaseName)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getMetadata 请求失败: %w", err)
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
			return "", fmt.Errorf("getMetadata 请求超时")
		}
	}
}

/*

// 场景1: 检查数据库存储使用情况
// 监控数据库大小，避免超过存储限制
func checkDatabaseUsage() error {
	result, err := CDPIndexedDBGetMetadata("myAppDB")
	if err != nil {
		log.Printf("获取数据库元数据失败: %v", err)
		return err
	}

	// 解析元数据
	var metadata struct {
		Result struct {
			EntriesCount    int     `json:"entriesCount"`
			KeyGeneratorValue float64 `json:"keyGeneratorValue"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &metadata); err != nil {
		log.Printf("解析元数据失败: %v", err)
		return err
	}

	log.Printf("数据库使用情况: 记录数=%d, 键生成器值=%.0f",
		metadata.Result.EntriesCount, metadata.Result.KeyGeneratorValue)

	// 根据使用情况做决策
	if metadata.Result.EntriesCount > 10000 {
		log.Println("数据库记录过多，建议清理")
	}

	return nil
}

// 场景2: 容量预警和清理决策
// 根据数据库大小决定是否进行清理
func monitorAndCleanupIfNeeded() error {
	result, err := CDPIndexedDBGetMetadata("userDataDB")
	if err != nil {
		log.Printf("监控数据库失败: %v", err)
		return err
	}

	// 解析响应获取实际数据大小
	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &respData); err != nil {
		return err
	}

	if resultObj, ok := respData["result"].(map[string]interface{}); ok {
		entriesCount := int(resultObj["entriesCount"].(float64))
		keyGeneratorValue := resultObj["keyGeneratorValue"].(float64)

		log.Printf("数据库统计: 总记录数=%d, 键生成器值=%.0f",
			entriesCount, keyGeneratorValue)

		// 如果记录数过多，触发清理
		if entriesCount > 5000 {
			log.Println("数据库记录过多，触发自动清理")
			// 执行清理逻辑
			return cleanupOldRecords()
		}
	}

	return nil
}

// 场景3: 调试数据库状态
// 在调试时获取数据库的元数据信息
func debugDatabaseState() error {
	result, err := CDPIndexedDBGetMetadata("cacheDB")
	if err != nil {
		log.Printf("调试数据库失败: %v", err)
		return err
	}

	// 输出详细的元数据信息用于调试
	fmt.Printf("数据库元数据详情:\n%s\n", result)

	// 可以进一步解析和分析
	var data map[string]interface{}
	json.Unmarshal([]byte(result), &data)

	// 将元数据保存到日志文件
	logData, _ := json.MarshalIndent(data, "", "  ")
	log.Printf("数据库调试信息: %s", string(logData))

	return nil
}

*/

// -----------------------------------------------  IndexedDB.requestData  -----------------------------------------------
// === 应用场景 ===
// 1. 数据查询: 从IndexedDB中读取特定数据记录
// 2. 分页加载: 分批加载大量数据以提高性能
// 3. 数据导出: 导出IndexedDB中的数据到其他格式
// 4. 数据验证: 验证存储的数据是否正确
// 5. 数据迁移: 在数据库迁移时读取旧数据
// 6. 数据恢复: 从IndexedDB中恢复用户数据

// CDPIndexedDBRequestData 请求IndexedDB数据
func CDPIndexedDBRequestData(databaseName, objectStoreName string, indexName string, skipCount, pageSize int, keyRange *KeyRange) (string, error) {
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
		"securityOrigin":  "",
		"databaseName":    databaseName,
		"objectStoreName": objectStoreName,
		"skipCount":       skipCount,
		"pageSize":        pageSize,
	}

	if indexName != "" {
		params["indexName"] = indexName
	}

	if keyRange != nil {
		params["keyRange"] = keyRange
	}

	// 序列化参数
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "IndexedDB.requestData",
		"params": %s
	}`, reqID, string(paramsJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 requestData 请求失败: %w", err)
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
			return "", fmt.Errorf("requestData 请求超时")
		}
	}
}

/*

// 场景1: 分页查询用户数据
// 分页加载用户列表，每页20条记录
func loadUsersPage(page int) error {
	skipCount := (page - 1) * 20
	pageSize := 20

	result, err := CDPIndexedDBRequestData("userDB", "users", "", skipCount, pageSize, nil)
	if err != nil {
		log.Printf("加载用户数据失败: %v", err)
		return err
	}

	// 解析返回的数据
	var response struct {
		Result struct {
			ObjectStoreDataEntries []struct {
				Key   interface{} `json:"key"`
				Value interface{} `json:"value"`
			} `json:"objectStoreDataEntries"`
			HasMore bool `json:"hasMore"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &response); err != nil {
		log.Printf("解析用户数据失败: %v", err)
		return err
	}

	log.Printf("加载了第%d页用户数据，共%d条，还有更多数据: %v",
		page, len(response.Result.ObjectStoreDataEntries), response.Result.HasMore)
	return nil
}

// 场景2: 按时间范围查询日志
// 查询特定时间范围内的日志记录
func queryLogsByTimeRange(startTime, endTime int64) error {
	// 创建时间范围的键范围
	keyRange := &KeyRange{
		Lower:     startTime,
		Upper:     endTime,
		LowerOpen: false,
		UpperOpen: false,
	}

	// 不跳过记录，查询最多100条
	result, err := CDPIndexedDBRequestData("appLogsDB", "logs", "timestamp", 0, 100, keyRange)
	if err != nil {
		log.Printf("查询日志失败: %v", err)
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &respData); err != nil {
		return err
	}

	if resultObj, ok := respData["result"].(map[string]interface{}); ok {
		entries := resultObj["objectStoreDataEntries"].([]interface{})
		log.Printf("找到%d条日志记录", len(entries))
	}

	return nil
}

// 场景3: 通过索引查询数据
// 通过email索引查询用户信息
func findUserByEmail(email string) error {
	// 使用索引查询特定email的用户
	keyRange := &KeyRange{
		Lower:     email,
		Upper:     email,
		LowerOpen: false,
		UpperOpen: false,
	}

	result, err := CDPIndexedDBRequestData("userDB", "profiles", "email", 0, 1, keyRange)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &respData); err != nil {
		return err
	}

	if resultObj, ok := respData["result"].(map[string]interface{}); ok {
		entries := resultObj["objectStoreDataEntries"].([]interface{})
		if len(entries) > 0 {
			entry := entries[0].(map[string]interface{})
			log.Printf("找到用户: %v", entry)
		} else {
			log.Println("未找到用户")
		}
	}

	return nil
}

*/

// -----------------------------------------------  IndexedDB.requestDatabase  -----------------------------------------------
// === 应用场景 ===
// 1. 数据库结构分析: 获取数据库的完整结构信息
// 2. 调试数据库架构: 在开发过程中查看数据库设计
// 3. 迁移规划: 在数据库迁移前分析现有结构
// 4. 兼容性检查: 检查数据库结构与应用版本的兼容性
// 5. 文档生成: 自动生成数据库结构文档
// 6. 数据恢复: 在数据恢复前了解数据库结构

// CDPIndexedDBRequestDatabase 请求数据库信息
func CDPIndexedDBRequestDatabase(databaseName string) (string, error) {
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
		"method": "IndexedDB.requestDatabase",
		"params": {
			"securityOrigin": "",
			"databaseName": "%s"
		}
	}`, reqID, databaseName)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 requestDatabase 请求失败: %w", err)
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
			return "", fmt.Errorf("requestDatabase 请求超时")
		}
	}
}

/*

// 场景1: 分析数据库结构
// 获取并分析数据库的完整结构信息
func analyzeDatabaseStructure() error {
	result, err := CDPIndexedDBRequestDatabase("myAppDB")
	if err != nil {
		log.Printf("获取数据库结构失败: %v", err)
		return err
	}

	// 解析数据库结构
	var dbInfo struct {
		Result struct {
			DatabaseWithObjectStores struct {
				Name          string `json:"name"`
				Version       int    `json:"version"`
				ObjectStores []struct {
					Name        string `json:"name"`
					KeyPath     string `json:"keyPath"`
					AutoIncrement bool `json:"autoIncrement"`
				} `json:"objectStores"`
			} `json:"databaseWithObjectStores"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &dbInfo); err != nil {
		log.Printf("解析数据库结构失败: %v", err)
		return err
	}

	db := dbInfo.Result.DatabaseWithObjectStores
	log.Printf("数据库: %s (版本: %d)", db.Name, db.Version)
	for _, store := range db.ObjectStores {
		log.Printf("  对象存储: %s (主键: %s, 自增: %v)",
			store.Name, store.KeyPath, store.AutoIncrement)
	}

	return nil
}

// 场景2: 检查数据库兼容性
// 在应用启动时检查数据库结构是否兼容
func checkDatabaseCompatibility() error {
	result, err := CDPIndexedDBRequestDatabase("appDataDB")
	if err != nil {
		log.Printf("检查数据库兼容性失败: %v", err)
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &respData); err != nil {
		return err
	}

	if resultObj, ok := respData["result"].(map[string]interface{}); ok {
		if dbObj, ok := resultObj["databaseWithObjectStores"].(map[string]interface{}); ok {
			dbVersion := int(dbObj["version"].(float64))
			dbName := dbObj["name"].(string)

			log.Printf("数据库: %s, 版本: %d", dbName, dbVersion)

			// 检查版本兼容性
			if dbVersion < 2 {
				log.Println("警告: 数据库版本过低，建议升级")
				return upgradeDatabase(dbName)
			}
		}
	}

	return nil
}

// 场景3: 生成数据库文档
// 为数据库结构生成文档
func generateDatabaseDocumentation() error {
	result, err := CDPIndexedDBRequestDatabase("userManagementDB")
	if err != nil {
		log.Printf("获取数据库信息失败: %v", err)
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &respData); err != nil {
		return err
	}

	// 生成Markdown格式的文档
	var docBuilder strings.Builder
	docBuilder.WriteString("# 数据库文档\n\n")

	if resultObj, ok := respData["result"].(map[string]interface{}); ok {
		if dbObj, ok := resultObj["databaseWithObjectStores"].(map[string]interface{}); ok {
			dbName := dbObj["name"].(string)
			dbVersion := int(dbObj["version"].(float64))

			docBuilder.WriteString(fmt.Sprintf("## 数据库: %s\n", dbName))
			docBuilder.WriteString(fmt.Sprintf("- 版本: %d\n\n", dbVersion))

			if stores, ok := dbObj["objectStores"].([]interface{}); ok {
				docBuilder.WriteString("## 对象存储\n\n")
				for _, store := range stores {
					storeMap := store.(map[string]interface{})
					storeName := storeMap["name"].(string)
					keyPath := storeMap["keyPath"].(string)

					docBuilder.WriteString(fmt.Sprintf("### %s\n", storeName))
					docBuilder.WriteString(fmt.Sprintf("- 主键: %s\n", keyPath))
					docBuilder.WriteString("\n")
				}
			}
		}
	}

	// 保存文档
	docContent := docBuilder.String()
	log.Printf("数据库文档已生成:\n%s", docContent)

	return nil
}


*/

// -----------------------------------------------  IndexedDB.requestDatabaseNames  -----------------------------------------------
// === 应用场景 ===
// 1. 数据库枚举: 列出当前页面中的所有IndexedDB数据库
// 2. 数据库管理: 在数据库管理工具中显示可用数据库
// 3. 数据清理: 查找并清理所有数据库
// 4. 调试辅助: 在调试时查看当前有哪些数据库
// 5. 迁移检查: 检查数据库迁移时的所有相关数据库
// 6. 安全审计: 审计页面中存储的所有数据库

// CDPIndexedDBRequestDatabaseNames 请求所有数据库名称
func CDPIndexedDBRequestDatabaseNames() (string, error) {
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
		"method": "IndexedDB.requestDatabaseNames",
		"params": {
			"securityOrigin": ""
		}
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 requestDatabaseNames 请求失败: %w", err)
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
			return "", fmt.Errorf("requestDatabaseNames 请求超时")
		}
	}
}

/*

// 场景1: 枚举所有数据库用于管理界面
// 在数据库管理工具中显示所有数据库
func listAllDatabases() ([]string, error) {
	result, err := CDPIndexedDBRequestDatabaseNames()
	if err != nil {
		log.Printf("获取数据库列表失败: %v", err)
		return nil, err
	}

	// 解析响应获取数据库名称列表
	var response struct {
		Result struct {
			DatabaseNames []string `json:"databaseNames"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &response); err != nil {
		log.Printf("解析数据库列表失败: %v", err)
		return nil, err
	}

	log.Printf("找到 %d 个数据库: %v",
		len(response.Result.DatabaseNames), response.Result.DatabaseNames)

	return response.Result.DatabaseNames, nil
}

// 场景2: 批量清理所有数据库
// 在应用卸载或重置时清理所有本地数据库
func cleanupAllDatabases() error {
	result, err := CDPIndexedDBRequestDatabaseNames()
	if err != nil {
		log.Printf("获取数据库列表失败: %v", err)
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &respData); err != nil {
		return err
	}

	if resultObj, ok := respData["result"].(map[string]interface{}); ok {
		if dbNames, ok := resultObj["databaseNames"].([]interface{}); ok {
			log.Printf("开始清理 %d 个数据库", len(dbNames))

			for _, dbName := range dbNames {
				dbNameStr := dbName.(string)
				log.Printf("正在删除数据库: %s", dbNameStr)

				// 删除数据库
				deleteResult, err := CDPIndexedDBDeleteDatabase(dbNameStr)
				if err != nil {
					log.Printf("删除数据库 %s 失败: %v", dbNameStr, err)
				} else {
					log.Printf("数据库 %s 已删除: %s", dbNameStr, deleteResult)
				}
			}
		}
	}

	return nil
}

// 场景3: 数据库迁移前的检查
// 在数据库迁移前检查所有需要处理的数据库
func checkDatabasesBeforeMigration() error {
	result, err := CDPIndexedDBRequestDatabaseNames()
	if err != nil {
		log.Printf("检查数据库失败: %v", err)
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &respData); err != nil {
		return err
	}

	if resultObj, ok := respData["result"].(map[string]interface{}); ok {
		if dbNames, ok := resultObj["databaseNames"].([]interface{}); ok {
			var oldVersionDBs []string
			var newVersionDBs []string

			for _, dbName := range dbNames {
				dbNameStr := dbName.(string)

				// 根据命名模式分类数据库
				if strings.HasSuffix(dbNameStr, "_v1") {
					oldVersionDBs = append(oldVersionDBs, dbNameStr)
				} else if strings.HasSuffix(dbNameStr, "_v2") {
					newVersionDBs = append(newVersionDBs, dbNameStr)
				}
			}

			log.Printf("迁移前检查: 旧版本数据库: %v", oldVersionDBs)
			log.Printf("迁移前检查: 新版本数据库: %v", newVersionDBs)

			if len(oldVersionDBs) > 0 {
				log.Println("发现需要迁移的旧版本数据库")
				// 触发迁移逻辑
				return migrateOldDatabases(oldVersionDBs)
			}
		}
	}

	return nil
}

*/
