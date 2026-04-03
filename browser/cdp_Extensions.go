package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Extensions.clearStorageItems  -----------------------------------------------
// === 应用场景 ===
// 1. 数据清理: 清理扩展的存储数据
// 2. 测试重置: 在自动化测试中重置扩展状态
// 3. 隐私保护: 清理用户敏感数据
// 4. 故障恢复: 清除损坏的存储数据
// 5. 扩展更新: 更新前清理旧数据
// 6. 环境隔离: 隔离不同测试环境的数据

// CDPExtensionsClearStorageItems 清理扩展的存储项
// extensionID: 扩展ID
// storageTypes: 存储类型数组 ["appcache", "cookies", "file_systems", "indexeddb", "local_storage", "shader_cache", "websql", "service_workers", "cache_storage"]
// options: 可选参数，可以为空
func CDPExtensionsClearStorageItems(extensionID string, storageTypes []string, options map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建存储类型数组
	storageTypesJSON, err := json.Marshal(storageTypes)
	if err != nil {
		return "", fmt.Errorf("序列化存储类型失败: %w", err)
	}

	// 构建选项参数
	var optionsJSON string
	if len(options) > 0 {
		optionsBytes, err := json.Marshal(options)
		if err != nil {
			return "", fmt.Errorf("序列化选项失败: %w", err)
		}
		optionsJSON = fmt.Sprintf(`"options": %s,`, string(optionsBytes))
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Extensions.clearStorageItems",
        "params": {
            "extensionId": "%s",
            "storageTypes": %s,
            %s
            "clearSince": null
        }
    }`, reqID, extensionID, string(storageTypesJSON), optionsJSON)

	// 移除可能的多余逗号
	message = strings.ReplaceAll(message, ",,", ",")
	message = strings.ReplaceAll(message, ",\n        }", "\n        }")

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 clearStorageItems 请求失败: %w", err)
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
			return "", fmt.Errorf("clearStorageItems 请求超时")
		}
	}
}

// 辅助函数: CDPExtensionsClearStorageItemsAllTypes 清理扩展的所有存储类型
func CDPExtensionsClearStorageItemsAllTypes(extensionID string) (string, error) {
	// 所有可能的存储类型
	allStorageTypes := []string{
		"appcache",
		"cookies",
		"file_systems",
		"indexeddb",
		"local_storage",
		"shader_cache",
		"websql",
		"service_workers",
		"cache_storage",
	}

	return CDPExtensionsClearStorageItems(extensionID, allStorageTypes, nil)
}

// 辅助函数: CDPExtensionsClearStorageItemsForPrivacy 清理隐私相关的存储
func CDPExtensionsClearStorageItemsForPrivacy(extensionID string) (string, error) {
	// 隐私相关的存储类型
	privacyStorageTypes := []string{
		"cookies",
		"local_storage",
		"indexeddb",
		"websql",
	}

	return CDPExtensionsClearStorageItems(extensionID, privacyStorageTypes, nil)
}

/*

// 示例1: 测试环境数据清理
func CleanupTestEnvironment() {
    extensionID := "abcdefghijklmnopqrstuvwxyz123456"

    // 清理所有存储数据
    result, err := CDPExtensionsClearStorageItemsAllTypes(extensionID)
    if err != nil {
        log.Printf("清理存储失败: %v", err)
        return
    }

    log.Printf("测试环境清理完成: %s", result)
}

// 示例2: 用户隐私数据清理
func CleanupUserPrivacyData() {
    // 获取所有扩展
    extensions, err := CDPExtensionsGetExtensions()
    if err != nil {
        log.Printf("获取扩展列表失败: %v", err)
        return
    }

    var extList struct {
        Result []struct {
            Id   string `json:"id"`
            Name string `json:"name"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(extensions), &extList); err != nil {
        log.Printf("解析扩展列表失败: %v", err)
        return
    }

    // 为每个扩展清理隐私数据
    for _, ext := range extList.Result {
        log.Printf("正在清理扩展 [%s] 的隐私数据...", ext.Name)

        result, err := CDPExtensionsClearStorageItemsForPrivacy(ext.Id)
        if err != nil {
            log.Printf("清理扩展 %s 失败: %v", ext.Name, err)
        } else {
            log.Printf("扩展 %s 隐私数据清理完成", ext.Name)
        }
    }
}

// 示例3: 指定存储类型清理
func SelectiveStorageCleanup() {
    extensionID := "test-extension-id-123"

    // 只清理特定的存储类型
    storageTypes := []string{"local_storage", "cookies", "indexeddb"}

    result, err := CDPExtensionsClearStorageItems(extensionID, storageTypes, nil)
    if err != nil {
        log.Printf("选择性清理失败: %v", err)
        return
    }

    log.Printf("选择性清理完成: %s", result)
}

// 示例4: 自动化测试套件中的数据清理
func TestSuiteWithStorageCleanup() {
    extensionID := "test-extension-456"

    // 测试前清理存储
    log.Println("=== 测试前准备 ===")
    cleanupResult, err := CDPExtensionsClearStorageItemsAllTypes(extensionID)
    if err != nil {
        log.Printf("测试前清理失败: %v", err)
    } else {
        log.Printf("测试前清理完成: %v", cleanupResult)
    }

    // 执行测试...
    log.Println("=== 执行测试用例 ===")
    // 这里执行实际的测试代码

    // 测试后清理
    log.Println("=== 测试后清理 ===")
    postCleanup, err := CDPExtensionsClearStorageItemsAllTypes(extensionID)
    if err != nil {
        log.Printf("测试后清理失败: %v", err)
    } else {
        log.Printf("测试后清理完成: %v", postCleanup)
    }
}


*/

// -----------------------------------------------  Extensions.getExtensions  -----------------------------------------------
// === 应用场景 ===
// 1. 扩展审计: 获取当前已安装的扩展列表
// 2. 状态检查: 检查特定扩展是否已安装
// 3. 版本管理: 获取扩展的版本信息
// 4. 权限验证: 检查扩展的权限配置
// 5. 环境诊断: 诊断浏览器扩展环境
// 6. 自动化测试: 验证扩展安装状态

// CDPExtensionsGetExtensions 获取已安装的扩展列表
func CDPExtensionsGetExtensions() (string, error) {
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
        "method": "Extensions.getExtensions"
    }`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getExtensions 请求失败: %w", err)
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
			return "", fmt.Errorf("getExtensions 请求超时")
		}
	}
}

// -----------------------------------------------  Extensions.getStorageItems  -----------------------------------------------
// === 应用场景 ===
// 1. 数据审计: 审计扩展的存储数据使用情况
// 2. 数据备份: 备份扩展的存储数据
// 3. 调试诊断: 诊断扩展的数据存储问题
// 4. 迁移支持: 支持数据迁移和同步
// 5. 状态检查: 检查扩展的存储状态
// 6. 数据分析: 分析扩展的数据使用模式

// CDPExtensionsGetStorageItems 获取扩展的存储项
// extensionID: 扩展ID
// storageTypes: 存储类型数组 ["appcache", "cookies", "file_systems", "indexeddb", "local_storage", "shader_cache", "websql", "service_workers", "cache_storage"]
func CDPExtensionsGetStorageItems(extensionID string, storageTypes []string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建存储类型数组
	storageTypesJSON, err := json.Marshal(storageTypes)
	if err != nil {
		return "", fmt.Errorf("序列化存储类型失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Extensions.getStorageItems",
        "params": {
            "extensionId": "%s",
            "storageTypes": %s
        }
    }`, reqID, extensionID, string(storageTypesJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 getStorageItems 请求失败: %w", err)
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
			return "", fmt.Errorf("getStorageItems 请求超时")
		}
	}
}

// 辅助函数: CDPExtensionsGetAllStorageItems 获取扩展的所有存储类型的数据
func CDPExtensionsGetAllStorageItems(extensionID string) (string, error) {
	// 所有可能的存储类型
	allStorageTypes := []string{
		"appcache",
		"cookies",
		"file_systems",
		"indexeddb",
		"local_storage",
		"shader_cache",
		"websql",
		"service_workers",
		"cache_storage",
	}

	return CDPExtensionsGetStorageItems(extensionID, allStorageTypes)
}

// 辅助函数: CDPExtensionsGetStorageStats 获取扩展存储统计信息
func CDPExtensionsGetStorageStats(extensionID string) (map[string]interface{}, error) {
	data, err := CDPExtensionsGetAllStorageItems(extensionID)
	if err != nil {
		return nil, fmt.Errorf("获取存储数据失败: %w", err)
	}

	var response struct {
		Result struct {
			StorageData []struct {
				Type string      `json:"type"`
				Data interface{} `json:"data"`
			} `json:"storageData"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, fmt.Errorf("解析存储数据失败: %w", err)
	}

	stats := make(map[string]interface{})
	totalItems := 0
	totalSize := 0

	for _, storage := range response.Result.StorageData {
		// 计算每种存储类型的数据量
		storageDataJSON, err := json.Marshal(storage.Data)
		if err == nil {
			size := len(storageDataJSON)
			stats[storage.Type] = map[string]interface{}{
				"size_bytes": size,
				"has_data":   size > 2, // 大于2表示不是空对象"{}"
			}
			totalSize += size
			totalItems++
		}
	}

	stats["total_storage_types"] = len(response.Result.StorageData)
	stats["total_size_bytes"] = totalSize
	stats["has_data"] = totalSize > 0

	return stats, nil
}

/*

// 示例1: 存储数据审计工具
func AuditExtensionStorage() {
    extensionID := "test-extension-123"

    // 获取所有存储数据
    storageData, err := CDPExtensionsGetAllStorageItems(extensionID)
    if err != nil {
        log.Printf("获取存储数据失败: %v", err)
        return
    }

    // 解析存储数据
    var response struct {
        Result struct {
            StorageData []struct {
                Type string      `json:"type"`
                Data interface{} `json:"data"`
            } `json:"storageData"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(storageData), &response); err != nil {
        log.Printf("解析存储数据失败: %v", err)
        return
    }

    fmt.Println("=== 扩展存储审计报告 ===")
    fmt.Printf("扩展ID: %s\n", extensionID)
    fmt.Printf("存储类型数量: %d\n\n", len(response.Result.StorageData))

    for _, storage := range response.Result.StorageData {
        dataJSON, _ := json.MarshalIndent(storage.Data, "  ", "  ")
        dataSize := len(dataJSON)

        fmt.Printf("存储类型: %s\n", storage.Type)
        fmt.Printf("数据大小: %d 字节\n", dataSize)
        fmt.Printf("数据内容:\n")

        if dataSize > 1024 { // 超过1KB的数据，只显示摘要
            var dataMap map[string]interface{}
            if err := json.Unmarshal(dataJSON, &dataMap); err == nil {
                keys := make([]string, 0, len(dataMap))
                for k := range dataMap {
                    keys = append(keys, k)
                }
                fmt.Printf("  包含 %d 个键: %v\n", len(keys), keys)
            }
        } else if dataSize > 2 { // 非空数据
            fmt.Printf("  %s\n", string(dataJSON))
        } else {
            fmt.Printf("  空数据\n")
        }
        fmt.Println()
    }
}

// 示例2: 存储数据备份工具
func BackupExtensionStorage(extensionID, backupPath string) error {
    // 获取存储数据
    storageData, err := CDPExtensionsGetAllStorageItems(extensionID)
    if err != nil {
        return fmt.Errorf("获取存储数据失败: %w", err)
    }

    // 解析数据
    var response struct {
        Result struct {
            StorageData []struct {
                Type string      `json:"type"`
                Data interface{} `json:"data"`
            } `json:"storageData"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(storageData), &response); err != nil {
        return fmt.Errorf("解析存储数据失败: %w", err)
    }

    // 创建备份结构
    backup := struct {
        ExtensionID string    `json:"extensionId"`
        BackupTime  time.Time `json:"backupTime"`
        StorageData []struct {
            Type string      `json:"type"`
            Data interface{} `json:"data"`
        } `json:"storageData"`
    }{
        ExtensionID: extensionID,
        BackupTime:  time.Now(),
        StorageData: response.Result.StorageData,
    }

    // 保存到文件
    backupJSON, err := json.MarshalIndent(backup, "", "  ")
    if err != nil {
        return fmt.Errorf("序列化备份数据失败: %w", err)
    }

    if err := os.WriteFile(backupPath, backupJSON, 0644); err != nil {
        return fmt.Errorf("保存备份文件失败: %w", err)
    }

    log.Printf("备份成功: %s (%d 字节)", backupPath, len(backupJSON))
    return nil
}

// 示例3: 存储使用情况监控
func MonitorStorageUsage(extensionID string) {
    stats, err := CDPExtensionsGetStorageStats(extensionID)
    if err != nil {
        log.Printf("获取存储统计失败: %v", err)
        return
    }

    fmt.Println("=== 存储使用情况监控 ===")
    fmt.Printf("扩展ID: %s\n", extensionID)
    fmt.Printf("监控时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

    for storageType, data := range stats {
        if dataMap, ok := data.(map[string]interface{}); ok {
            hasData, _ := dataMap["has_data"].(bool)
            size, _ := dataMap["size_bytes"].(float64)

            status := "无数据"
            if hasData {
                status = fmt.Sprintf("%.0f 字节", size)
            }

            fmt.Printf("%-20s: %s\n", storageType, status)
        }
    }

    totalSize, _ := stats["total_size_bytes"].(float64)
    fmt.Printf("\n总计: %.0f 字节 (%.2f KB)\n", totalSize, totalSize/1024)
}

// 示例4: 存储数据分析
func AnalyzeStorageData(extensionID string) {
    // 只获取特定的存储类型进行分析
    storageTypes := []string{"local_storage", "cookies", "indexeddb"}

    data, err := CDPExtensionsGetStorageItems(extensionID, storageTypes)
    if err != nil {
        log.Printf("获取存储数据失败: %v", err)
        return
    }

    var response struct {
        Result struct {
            StorageData []struct {
                Type string      `json:"type"`
                Data interface{} `json:"data"`
            } `json:"storageData"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(data), &response); err != nil {
        log.Printf("解析存储数据失败: %v", err)
        return
    }

    fmt.Println("=== 存储数据分析报告 ===")

    for _, storage := range response.Result.StorageData {
        switch storage.Type {
        case "local_storage":
            analyzeLocalStorage(storage.Data)
        case "cookies":
            analyzeCookies(storage.Data)
        case "indexeddb":
            analyzeIndexedDB(storage.Data)
        }
    }
}

func analyzeLocalStorage(data interface{}) {
    fmt.Println("\n[Local Storage 分析]")
    if dataMap, ok := data.(map[string]interface{}); ok {
        fmt.Printf("键数量: %d\n", len(dataMap))
        for key, value := range dataMap {
            valueStr := fmt.Sprintf("%v", value)
            if len(valueStr) > 50 {
                valueStr = valueStr[:50] + "..."
            }
            fmt.Printf("  %s: %s\n", key, valueStr)
        }
    }
}

func analyzeCookies(data interface{}) {
    fmt.Println("\n[Cookies 分析]")
    if cookies, ok := data.([]interface{}); ok {
        fmt.Printf("Cookie数量: %d\n", len(cookies))

        // 按域名分组统计
        domainCount := make(map[string]int)
        secureCount := 0
        httpOnlyCount := 0

        for _, cookie := range cookies {
            if cookieMap, ok := cookie.(map[string]interface{}); ok {
                domain, _ := cookieMap["domain"].(string)
                secure, _ := cookieMap["secure"].(bool)
                httpOnly, _ := cookieMap["httpOnly"].(bool)

                if domain != "" {
                    domainCount[domain]++
                }
                if secure {
                    secureCount++
                }
                if httpOnly {
                    httpOnlyCount++
                }
            }
        }

        fmt.Println("按域名统计:")
        for domain, count := range domainCount {
            fmt.Printf("  %s: %d\n", domain, count)
        }
        fmt.Printf("安全Cookie: %d\n", secureCount)
        fmt.Printf("HttpOnly Cookie: %d\n", httpOnlyCount)
    }
}

*/

// -----------------------------------------------  Extensions.loadUnpacked  -----------------------------------------------
// === 应用场景 ===
// 1. 开发调试: 加载本地未打包的扩展程序进行开发调试
// 2. 测试验证: 在自动化测试中动态加载扩展程序
// 3. 环境准备: 为测试环境预先安装必要的扩展程序
// 4. 动态加载: 运行时根据需要动态添加功能扩展
// 5. 扩展管理: 程序化管理浏览器扩展的加载和卸载
// 6. 插件验证: 验证自定义扩展程序的功能完整性

// CDPExtensionsLoadUnpacked 加载未打包的扩展程序
func CDPExtensionsLoadUnpacked(path string) (string, error) {
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
        "method": "Extensions.loadUnpacked",
        "params": {
            "path": "%s"
        }
    }`, reqID, path)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 loadUnpacked 请求失败: %w", err)
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
			return "", fmt.Errorf("loadUnpacked 请求超时")
		}
	}
}

// -----------------------------------------------  Extensions.removeStorageItems  -----------------------------------------------
// === 应用场景 ===
// 1. 数据清理: 选择性清理扩展存储中的特定数据
// 2. 隐私保护: 删除用户敏感数据
// 3. 存储管理: 管理扩展的存储空间使用
// 4. 数据重置: 重置特定类型的数据
// 5. 缓存清理: 清理扩展的缓存数据
// 6. 存储优化: 优化扩展存储结构

// CDPExtensionsRemoveStorageItems 删除扩展的存储项
// extensionID: 扩展ID
// storageTypes: 存储类型数组 ["appcache", "cookies", "file_systems", "indexeddb", "local_storage", "shader_cache", "websql", "service_workers", "cache_storage"]
// options: 可选参数，控制删除行为
func CDPExtensionsRemoveStorageItems(extensionID string, storageTypes []string, options map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建存储类型数组
	storageTypesJSON, err := json.Marshal(storageTypes)
	if err != nil {
		return "", fmt.Errorf("序列化存储类型失败: %w", err)
	}

	// 构建选项参数
	var optionsJSON string
	if len(options) > 0 {
		optionsBytes, err := json.Marshal(options)
		if err != nil {
			return "", fmt.Errorf("序列化选项失败: %w", err)
		}
		optionsJSON = fmt.Sprintf(`"options": %s,`, string(optionsBytes))
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Extensions.removeStorageItems",
        "params": {
            "extensionId": "%s",
            "storageTypes": %s,
            %s
            "clearSince": null
        }
    }`, reqID, extensionID, string(storageTypesJSON), optionsJSON)

	// 移除可能的多余逗号
	message = strings.ReplaceAll(message, ",,", ",")
	message = strings.ReplaceAll(message, ",\n        }", "\n        }")

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 removeStorageItems 请求失败: %w", err)
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
			return "", fmt.Errorf("removeStorageItems 请求超时")
		}
	}
}

// 辅助函数: CDPExtensionsRemoveCookies 删除扩展的cookies
func CDPExtensionsRemoveCookies(extensionID string, domainPatterns ...string) (string, error) {
	storageTypes := []string{"cookies"}

	var options map[string]interface{}
	if len(domainPatterns) > 0 {
		options = map[string]interface{}{
			"origin": domainPatterns[0],
		}
	}

	return CDPExtensionsRemoveStorageItems(extensionID, storageTypes, options)
}

// 辅助函数: CDPExtensionsRemoveLocalStorage 删除扩展的localStorage
func CDPExtensionsRemoveLocalStorage(extensionID string, keyPatterns ...string) (string, error) {
	storageTypes := []string{"local_storage"}

	var options map[string]interface{}
	if len(keyPatterns) > 0 {
		options = map[string]interface{}{
			"key": keyPatterns[0],
		}
	}

	return CDPExtensionsRemoveStorageItems(extensionID, storageTypes, options)
}

// 辅助函数: CDPExtensionsRemoveIndexedDB 删除扩展的IndexedDB数据
func CDPExtensionsRemoveIndexedDB(extensionID string, dbNamePatterns ...string) (string, error) {
	storageTypes := []string{"indexeddb"}

	var options map[string]interface{}
	if len(dbNamePatterns) > 0 {
		options = map[string]interface{}{
			"databaseName": dbNamePatterns[0],
		}
	}

	return CDPExtensionsRemoveStorageItems(extensionID, storageTypes, options)
}

// 辅助函数: CDPExtensionsRemoveStorageDataByAge 删除指定时间前的数据
func CDPExtensionsRemoveStorageDataByAge(extensionID string, storageTypes []string, maxAge time.Duration) (string, error) {
	clearSince := time.Now().Add(-maxAge).Unix()

	options := map[string]interface{}{
		"clearSince": clearSince,
	}

	return CDPExtensionsRemoveStorageItems(extensionID, storageTypes, options)
}

/*

// 示例1: 选择性清理特定类型的数据
func SelectiveDataCleanup() {
    extensionID := "my-extension-123"

    // 只清理cookies和localStorage
    storageTypes := []string{"cookies", "local_storage"}

    result, err := CDPExtensionsRemoveStorageItems(extensionID, storageTypes, nil)
    if err != nil {
        log.Printf("选择性清理失败: %v", err)
        return
    }

    log.Printf("选择性清理完成: %s", result)
}

// 示例2: 清理指定域名的cookies
func CleanupCookiesByDomain() {
    extensionID := "my-extension-123"

    // 清理指定域名的cookies
    result, err := CDPExtensionsRemoveCookies(extensionID, "*.example.com")
    if err != nil {
        log.Printf("清理cookies失败: %v", err)
        return
    }

    log.Printf("cookies清理完成: %s", result)
}

// 示例3: 清理旧的存储数据
func CleanupOldStorageData() {
    extensionID := "my-extension-123"

    // 清理30天前的数据
    storageTypes := []string{
        "local_storage",
        "indexeddb",
        "cookies",
    }

    result, err := CDPExtensionsRemoveStorageDataByAge(extensionID, storageTypes, 30 * 24*time.Hour)
    if err != nil {
        log.Printf("清理旧数据失败: %v", err)
        return
    }

    log.Printf("旧数据清理完成: %s", result)
}

// 示例4: 存储空间管理
type StorageCleanupManager struct {
    ExtensionID string
    MaxSizeMB   float64
    CleanupLog  []CleanupRecord
}

type CleanupRecord struct {
    Timestamp   time.Time
    Action      string
    StorageType string
    Options     map[string]interface{}
    Result      string
}

func (m *StorageCleanupManager) CheckAndCleanup() error {
    // 先获取存储统计
    stats, err := CDPExtensionsGetStorageStats(m.ExtensionID)
    if err != nil {
        return fmt.Errorf("获取存储统计失败: %w", err)
    }

    totalSizeBytes, ok := stats["total_size_bytes"].(float64)
    if !ok {
        return fmt.Errorf("无法获取总大小")
    }

    totalSizeMB := totalSizeBytes / 1024 / 1024

    if totalSizeMB <= m.MaxSizeMB {
        log.Printf("存储空间正常: %.2f MB / %.2f MB", totalSizeMB, m.MaxSizeMB)
        return nil
    }

    log.Printf("存储空间超出限制: %.2f MB / %.2f MB, 开始清理...", totalSizeMB, m.MaxSizeMB)

    // 清理缓存数据
    cleanupTypes := []string{"appcache", "cache_storage", "shader_cache"}
    result, err := CDPExtensionsRemoveStorageItems(m.ExtensionID, cleanupTypes, nil)
    if err != nil {
        return fmt.Errorf("清理缓存失败: %w", err)
    }

    m.CleanupLog = append(m.CleanupLog, CleanupRecord{
        Timestamp:   time.Now(),
        Action:      "cleanup_cache",
        StorageType: "cache",
        Result:      result,
    })

    // 重新检查大小
    stats, _ = CDPExtensionsGetStorageStats(m.ExtensionID)
    if newTotalSize, ok := stats["total_size_bytes"].(float64); ok {
        newTotalSizeMB := newTotalSize / 1024 / 1024
        log.Printf("清理后存储空间: %.2f MB / %.2f MB", newTotalSizeMB, m.MaxSizeMB)
    }

    return nil
}

// 示例5: 清理指定模式的localStorage键
func CleanupLocalStorageByPattern(extensionID, keyPattern string) error {
    // 先获取localStorage数据
    storageData, err := CDPExtensionsGetStorageItems(extensionID, []string{"local_storage"})
    if err != nil {
        return fmt.Errorf("获取localStorage数据失败: %w", err)
    }

    var response struct {
        Result struct {
            StorageData []struct {
                Type string      `json:"type"`
                Data interface{} `json:"data"`
            } `json:"storageData"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(storageData), &response); err != nil {
        return fmt.Errorf("解析localStorage数据失败: %w", err)
    }

    if len(response.Result.StorageData) == 0 {
        return fmt.Errorf("没有localStorage数据")
    }

    localStorageData := response.Result.StorageData[0].Data
    if localStorageData == nil {
        return fmt.Errorf("localStorage数据为空")
    }

    // 匹配要删除的键
    pattern := regexp.MustCompile(keyPattern)
    var keysToRemove []string

    if dataMap, ok := localStorageData.(map[string]interface{}); ok {
        for key := range dataMap {
            if pattern.MatchString(key) {
                keysToRemove = append(keysToRemove, key)
            }
        }
    }

    if len(keysToRemove) == 0 {
        log.Printf("没有匹配的键需要删除")
        return nil
    }

    log.Printf("找到 %d 个匹配的键需要删除: %v", len(keysToRemove), keysToRemove)

    // 由于CDP API不支持按键删除，我们只能清理整个localStorage
    // 或者可以记录这些键，然后在应用层面处理
    for _, key := range keysToRemove {
        log.Printf("需要删除的键: %s", key)
    }

    // 如果需要完全清理，可以调用
    result, err := CDPExtensionsRemoveLocalStorage(extensionID)
    if err != nil {
        return fmt.Errorf("清理localStorage失败: %w", err)
    }

    log.Printf("localStorage清理完成: %s", result)
    return nil
}

// 示例6: 分阶段清理策略
func PhasedCleanupStrategy(extensionID string) {
    log.Println("=== 开始分阶段清理 ===")

    // 阶段1: 清理缓存数据
    log.Println("阶段1: 清理缓存数据...")
    cacheTypes := []string{"appcache", "cache_storage", "shader_cache"}
    result1, err1 := CDPExtensionsRemoveStorageItems(extensionID, cacheTypes, nil)
    if err1 != nil {
        log.Printf("阶段1失败: %v", err1)
    } else {
        log.Printf("阶段1完成: 清理了缓存数据")
    }

    // 阶段2: 清理旧数据
    log.Println("阶段2: 清理7天前的数据...")
    oldDataTypes := []string{"cookies", "local_storage", "indexeddb"}
    result2, err2 := CDPExtensionsRemoveStorageDataByAge(extensionID, oldDataTypes, 7 * 24*time.Hour)
    if err2 != nil {
        log.Printf("阶段2失败: %v", err2)
    } else {
        log.Printf("阶段2完成: 清理了旧数据")
    }

    // 阶段3: 清理大型数据
    log.Println("阶段3: 清理特定的大型存储...")
    largeDataTypes := []string{"file_systems", "websql"}
    result3, err3 := CDPExtensionsRemoveStorageItems(extensionID, largeDataTypes, nil)
    if err3 != nil {
        log.Printf("阶段3失败: %v", err3)
    } else {
        log.Printf("阶段3完成: 清理了大型数据")
    }

    // 最终检查
    log.Println("=== 清理完成 ===")
    if err1 == nil && err2 == nil && err3 == nil {
        log.Println("所有阶段清理成功完成")
    } else {
        log.Println("清理完成，但部分阶段失败")
    }
}

// 示例7: 备份后清理
func BackupAndCleanup(extensionID, backupDir string) error {
    timestamp := time.Now().Format("20060102_150405")
    backupFile := fmt.Sprintf("%s/backup_%s_%s.json", backupDir, extensionID, timestamp)

    // 1. 备份数据
    log.Printf("开始备份扩展 %s 的数据...", extensionID)
    if err := BackupExtensionStorage(extensionID, backupFile); err != nil {
        return fmt.Errorf("备份失败: %w", err)
    }
    log.Printf("数据已备份到: %s", backupFile)

    // 2. 清理用户数据但保留配置
    log.Println("清理用户数据...")
    userDataTypes := []string{"cookies", "indexeddb", "websql"}
    result, err := CDPExtensionsRemoveStorageItems(extensionID, userDataTypes, nil)
    if err != nil {
        return fmt.Errorf("清理用户数据失败: %w", err)
    }
    log.Printf("用户数据清理完成: %s", result)

    // 3. 清理缓存
    log.Println("清理缓存数据...")
    cacheTypes := []string{"appcache", "cache_storage", "shader_cache"}
    result, err = CDPExtensionsRemoveStorageItems(extensionID, cacheTypes, nil)
    if err != nil {
        return fmt.Errorf("清理缓存失败: %w", err)
    }
    log.Printf("缓存清理完成: %s", result)

    return nil
}

// 示例8: 定时清理任务
func ScheduledCleanup(extensionID string, cleanupTypes []string, interval time.Duration, stopChan chan bool) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    log.Printf("启动定时清理，间隔: %v", interval)

    for {
        select {
        case <-ticker.C:
            log.Println("执行定时清理...")
            result, err := CDPExtensionsRemoveStorageItems(extensionID, cleanupTypes, nil)
            if err != nil {
                log.Printf("定时清理失败: %v", err)
            } else {
                log.Printf("定时清理完成: %s", result)

                // 检查清理后的存储大小
                stats, err := CDPExtensionsGetStorageStats(extensionID)
                if err == nil {
                    if totalSize, ok := stats["total_size_bytes"].(float64); ok {
                        log.Printf("当前存储大小: %.2f KB", totalSize/1024)
                    }
                }
            }

        case <-stopChan:
            log.Println("停止定时清理")
            return
        }
    }
}


*/

// -----------------------------------------------  Extensions.setStorageItems  -----------------------------------------------
// === 应用场景 ===
// 1. 数据恢复: 恢复备份的扩展数据
// 2. 数据注入: 向扩展注入测试数据
// 3. 配置预设: 预设扩展的配置数据
// 4. 数据迁移: 迁移数据到新扩展
// 5. 测试数据: 为测试设置特定数据
// 6. 状态恢复: 恢复扩展的特定状态

// CDPExtensionsSetStorageItems 设置扩展的存储项
// extensionID: 扩展ID
// storageItems: 存储项数组，每个项包含存储类型和数据
func CDPExtensionsSetStorageItems(extensionID string, storageItems []StorageItem) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建存储项数组
	storageItemsJSON, err := json.Marshal(storageItems)
	if err != nil {
		return "", fmt.Errorf("序列化存储项失败: %w", err)
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Extensions.setStorageItems",
        "params": {
            "extensionId": "%s",
            "storageItems": %s
        }
    }`, reqID, extensionID, string(storageItemsJSON))

	// 发送请求
	err = chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setStorageItems 请求失败: %w", err)
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
			return "", fmt.Errorf("setStorageItems 请求超时")
		}
	}
}

// StorageItem 存储项结构
type StorageItem struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// 辅助函数: CDPExtensionsSetLocalStorage 设置扩展的localStorage
func CDPExtensionsSetLocalStorage(extensionID string, data map[string]interface{}) (string, error) {
	storageItem := StorageItem{
		Type: "local_storage",
		Data: data,
	}

	return CDPExtensionsSetStorageItems(extensionID, []StorageItem{storageItem})
}

// 辅助函数: CDPExtensionsSetCookies 设置扩展的cookies
func CDPExtensionsSetCookies(extensionID string, cookies []Cookie) (string, error) {
	storageItem := StorageItem{
		Type: "cookies",
		Data: cookies,
	}

	return CDPExtensionsSetStorageItems(extensionID, []StorageItem{storageItem})
}

// Cookie 结构定义
type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain,omitempty"`
	Path     string  `json:"path,omitempty"`
	Expires  float64 `json:"expires,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	HttpOnly bool    `json:"httpOnly,omitempty"`
	SameSite string  `json:"sameSite,omitempty"` // "Strict" | "Lax" | "None"
}

// 辅助函数: CDPExtensionsRestoreFromBackup 从备份恢复数据
func CDPExtensionsRestoreFromBackup(extensionID, backupFile string) (string, error) {
	// 读取备份文件
	backupData, err := os.ReadFile(backupFile)
	if err != nil {
		return "", fmt.Errorf("读取备份文件失败: %w", err)
	}

	var backup struct {
		ExtensionID string        `json:"extensionId"`
		BackupTime  time.Time     `json:"backupTime"`
		StorageData []StorageItem `json:"storageData"`
	}

	if err := json.Unmarshal(backupData, &backup); err != nil {
		return "", fmt.Errorf("解析备份数据失败: %w", err)
	}

	// 验证扩展ID
	if backup.ExtensionID != "" && backup.ExtensionID != extensionID {
		log.Printf("警告: 备份文件是为扩展 %s 创建的，当前扩展是 %s",
			backup.ExtensionID, extensionID)
	}

	// 恢复数据
	return CDPExtensionsSetStorageItems(extensionID, backup.StorageData)
}

// 辅助函数: CDPExtensionsSetIndexedDB 设置扩展的IndexedDB数据
func CDPExtensionsSetIndexedDB(extensionID string, databases []IndexedDBDatabase) (string, error) {
	storageItem := StorageItem{
		Type: "indexeddb",
		Data: databases,
	}

	return CDPExtensionsSetStorageItems(extensionID, []StorageItem{storageItem})
}

// IndexedDBDatabase IndexedDB数据库结构
type IndexedDBDatabase struct {
	Name    string           `json:"name"`
	Version int              `json:"version"`
	Stores  []IndexedDBStore `json:"stores,omitempty"`
}

// IndexedDBStore IndexedDB存储结构
type IndexedDBStore struct {
	Name          string        `json:"name"`
	KeyPath       string        `json:"keyPath,omitempty"`
	AutoIncrement bool          `json:"autoIncrement,omitempty"`
	Records       []interface{} `json:"records,omitempty"`
}

/*

// 示例1: 恢复备份数据
func RestoreExtensionData() {
    extensionID := "my-extension-123"
    backupFile := "./backups/backup_my-extension-123_20240101_120000.json"

    result, err := CDPExtensionsRestoreFromBackup(extensionID, backupFile)
    if err != nil {
        log.Printf("恢复数据失败: %v", err)
        return
    }

    log.Printf("数据恢复完成: %s", result)
}

// 示例2: 注入测试数据
func InjectTestData() {
    extensionID := "test-extension-456"

    // 创建测试数据
    localStorageData := map[string]interface{}{
        "user_preferences": map[string]interface{}{
            "theme": "dark",
            "language": "zh-CN",
            "notifications": true,
        },
        "test_settings": map[string]interface{}{
            "debug_mode": true,
            "log_level": "verbose",
        },
    }

    // 设置localStorage
    result, err := CDPExtensionsSetLocalStorage(extensionID, localStorageData)
    if err != nil {
        log.Printf("设置localStorage失败: %v", err)
        return
    }

    log.Printf("测试数据注入完成: %s", result)
}

// 示例3: 设置cookies用于认证测试
func SetupAuthCookies() {
    extensionID := "web-app-extension"

    // 创建认证cookies
    cookies := []Cookie{
        {
            Name:     "session_id",
            Value:    "abc123def456",
            Domain:   ".example.com",
            Path:     "/",
            Expires:  float64(time.Now().Add(24 * time.Hour).Unix()),
            Secure:   true,
            HttpOnly: true,
            SameSite: "Strict",
        },
        {
            Name:     "auth_token",
            Value:    "xyz789uvw012",
            Domain:   ".example.com",
            Path:     "/",
            Secure:   true,
            SameSite: "Lax",
        },
        {
            Name:    "preferences",
            Value:   "theme=dark&lang=zh",
            Domain:  ".example.com",
            Path:    "/",
            Expires: float64(time.Now().Add(365 * 24 * time.Hour).Unix()),
        },
    }

    result, err := CDPExtensionsSetCookies(extensionID, cookies)
    if err != nil {
        log.Printf("设置cookies失败: %v", err)
        return
    }

    log.Printf("认证cookies设置完成: %s", result)
}

// 示例4: 创建复杂的IndexedDB测试数据
func SetupIndexedDBTestData() {
    extensionID := "data-extension-789"

    // 创建IndexedDB测试数据
    databases := []IndexedDBDatabase{
        {
            Name:    "UserDatabase",
            Version: 1,
            Stores: []IndexedDBStore{
                {
                    Name:          "users",
                    KeyPath:       "id",
                    AutoIncrement: true,
                    Records: []interface{}{
                        map[string]interface{}{
                            "id":    1,
                            "name":  "张三",
                            "email": "zhangsan@example.com",
                            "age":   30,
                        },
                        map[string]interface{}{
                            "id":    2,
                            "name":  "李四",
                            "email": "lisi@example.com",
                            "age":   25,
                        },
                        map[string]interface{}{
                            "id":    3,
                            "name":  "王五",
                            "email": "wangwu@example.com",
                            "age":   28,
                        },
                    },
                },
                {
                    Name:          "products",
                    KeyPath:       "sku",
                    AutoIncrement: false,
                    Records: []interface{}{
                        map[string]interface{}{
                            "sku":    "P001",
                            "name":   "笔记本电脑",
                            "price":  5999.99,
                            "stock":  50,
                        },
                        map[string]interface{}{
                            "sku":    "P002",
                            "name":   "智能手机",
                            "price":  2999.99,
                            "stock":  100,
                        },
                    },
                },
            },
        },
        {
            Name:    "LogDatabase",
            Version: 1,
            Stores: []IndexedDBStore{
                {
                    Name:          "logs",
                    KeyPath:       "timestamp",
                    AutoIncrement: false,
                    Records: []interface{}{
                        map[string]interface{}{
                            "timestamp": time.Now().Add(-1 * time.Hour).Unix(),
                            "level":     "INFO",
                            "message":   "应用启动成功",
                        },
                        map[string]interface{}{
                            "timestamp": time.Now().Add(-30 * time.Minute).Unix(),
                            "level":     "WARN",
                            "message":   "磁盘空间不足",
                        },
                        map[string]interface{}{
                            "timestamp": time.Now().Unix(),
                            "level":     "INFO",
                            "message":   "数据加载完成",
                        },
                    },
                },
            },
        },
    }

    result, err := CDPExtensionsSetIndexedDB(extensionID, databases)
    if err != nil {
        log.Printf("设置IndexedDB失败: %v", err)
        return
    }

    log.Printf("IndexedDB测试数据设置完成: %s", result)
}

// 示例5: 完整的数据迁移流程
func MigrateExtensionData(sourceExtensionID, targetExtensionID string) error {
    log.Printf("开始数据迁移: %s -> %s", sourceExtensionID, targetExtensionID)

    // 1. 从源扩展获取数据
    log.Println("步骤1: 获取源扩展数据...")
    sourceData, err := CDPExtensionsGetAllStorageItems(sourceExtensionID)
    if err != nil {
        return fmt.Errorf("获取源数据失败: %w", err)
    }

    // 解析源数据
    var sourceResponse struct {
        Result struct {
            StorageData []StorageItem `json:"storageData"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(sourceData), &sourceResponse); err != nil {
        return fmt.Errorf("解析源数据失败: %w", err)
    }

    log.Printf("获取到 %d 种存储类型的数据", len(sourceResponse.Result.StorageData))

    // 2. 转换数据（如果需要）
    log.Println("步骤2: 转换数据...")
    var targetStorageItems []StorageItem
    for _, item := range sourceResponse.Result.StorageData {
        // 这里可以添加数据转换逻辑
        // 例如：修改特定字段、过滤数据等
        targetStorageItems = append(targetStorageItems, item)
    }

    // 3. 设置到目标扩展
    log.Println("步骤3: 设置目标扩展数据...")
    result, err := CDPExtensionsSetStorageItems(targetExtensionID, targetStorageItems)
    if err != nil {
        return fmt.Errorf("设置目标数据失败: %w", err)
    }

    log.Printf("数据迁移完成: %s", result)
    return nil
}

// 示例6: 为自动化测试设置预设状态
func SetupTestEnvironment(extensionID string, testCase string) error {
    log.Printf("为测试用例 '%s' 设置环境", testCase)

    switch testCase {
    case "auth_required":
        // 设置认证状态
        localStorageData := map[string]interface{}{
            "auth_state": map[string]interface{}{
                "isLoggedIn": true,
                "userId":     "user123",
                "userRole":   "admin",
            },
        }

        cookies := []Cookie{
            {
                Name:   "auth_token",
                Value:  "test_token_123",
                Domain: ".test.com",
                Secure: true,
            },
        }

        // 设置localStorage
        if _, err := CDPExtensionsSetLocalStorage(extensionID, localStorageData); err != nil {
            return err
        }

        // 设置cookies
        if _, err := CDPExtensionsSetCookies(extensionID, cookies); err != nil {
            return err
        }

    case "data_loaded":
        // 设置加载了数据的场景
        localStorageData := map[string]interface{}{
            "data_state": map[string]interface{}{
                "lastSyncTime": time.Now().Unix(),
                "totalItems":   150,
                "lastPage":     5,
            },
        }

        indexedDBData := []IndexedDBDatabase{
            {
                Name:    "AppData",
                Version: 1,
                Stores: []IndexedDBStore{
                    {
                        Name:    "items",
                        KeyPath: "id",
                        Records: []interface{}{
                            map[string]interface{}{"id": 1, "name": "测试项目1", "status": "active"},
                            map[string]interface{}{"id": 2, "name": "测试项目2", "status": "completed"},
                            map[string]interface{}{"id": 3, "name": "测试项目3", "status": "pending"},
                        },
                    },
                },
            },
        }

        storageItems := []StorageItem{
            {Type: "local_storage", Data: localStorageData},
            {Type: "indexeddb", Data: indexedDBData},
        }

        if _, err := CDPExtensionsSetStorageItems(extensionID, storageItems); err != nil {
            return err
        }

    case "empty_state":
        // 设置为空状态
        storageItems := []StorageItem{
            {Type: "local_storage", Data: map[string]interface{}{}},
            {Type: "cookies", Data: []Cookie{}},
        }

        if _, err := CDPExtensionsSetStorageItems(extensionID, storageItems); err != nil {
            return err
        }

    default:
        return fmt.Errorf("未知的测试用例: %s", testCase)
    }

    log.Println("测试环境设置完成")
    return nil
}

// 示例7: 批量设置多个存储类型的数据
func BulkSetStorageData() {
    extensionID := "complex-extension-999"

    // 准备多种存储类型的数据
    storageItems := []StorageItem{
        {
            Type: "local_storage",
            Data: map[string]interface{}{
                "app_config": map[string]interface{}{
                    "version": "1.0.0",
                    "theme":   "blue",
                    "locale":  "zh_CN",
                },
                "user_data": map[string]interface{}{
                    "username": "testuser",
                    "preferences": map[string]interface{}{
                        "auto_save": true,
                        "notify":    false,
                    },
                },
            },
        },
        {
            Type: "cookies",
            Data: []Cookie{
                {Name: "session", Value: "abc123", Domain: ".example.com"},
                {Name: "pref", Value: "theme=dark", Domain: ".example.com"},
            },
        },
        {
            Type: "indexeddb",
            Data: []IndexedDBDatabase{
                {
                    Name:    "CacheDB",
                    Version: 1,
                    Stores: []IndexedDBStore{
                        {
                            Name:    "cache",
                            KeyPath: "key",
                            Records: []interface{}{
                                map[string]interface{}{
                                    "key":   "homepage",
                                    "value": "<html>...</html>",
                                    "timestamp": time.Now().Unix(),
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    result, err := CDPExtensionsSetStorageItems(extensionID, storageItems)
    if err != nil {
        log.Printf("批量设置存储数据失败: %v", err)
        return
    }

    log.Printf("批量数据设置完成: %s", result)
}

// 示例8: 实时数据同步工具
type DataSyncManager struct {
    ExtensionID    string
    SyncInterval   time.Duration
    LastSyncTime   time.Time
    DataTemplate   []StorageItem
    IsSyncing      bool
}

func NewDataSyncManager(extensionID string, interval time.Duration) *DataSyncManager {
    return &DataSyncManager{
        ExtensionID:  extensionID,
        SyncInterval: interval,
        DataTemplate: make([]StorageItem, 0),
    }
}

func (m *DataSyncManager) AddDataTemplate(storageType string, dataGenerator func() interface{}) {
    m.DataTemplate = append(m.DataTemplate, StorageItem{
        Type: storageType,
        Data: dataGenerator(),
    })
}

func (m *DataSyncManager) StartSync() {
    ticker := time.NewTicker(m.SyncInterval)

    go func() {
        for range ticker.C {
            if m.IsSyncing {
                continue
            }

            m.IsSyncing = true
            m.syncData()
            m.IsSyncing = false
        }
    }()
}

func (m *DataSyncManager) syncData() {
    // 更新数据模板中的动态数据
    var storageItems []StorageItem
    for _, item := range m.DataTemplate {
        // 如果是函数类型，执行函数获取新数据
        if dataFunc, ok := item.Data.(func() interface{}); ok {
            storageItems = append(storageItems, StorageItem{
                Type: item.Type,
                Data: dataFunc(),
            })
        } else {
            storageItems = append(storageItems, item)
        }
    }

    // 设置数据
    result, err := CDPExtensionsSetStorageItems(m.ExtensionID, storageItems)
    if err != nil {
        log.Printf("数据同步失败: %v", err)
        return
    }

    m.LastSyncTime = time.Now()
    log.Printf("数据同步完成: %s", result)
}

func main() {
    extensionID := "sync-extension-001"

    fmt.Println("=== 示例1: 恢复备份数据 ===")
    RestoreExtensionData()

    fmt.Println("\n=== 示例2: 注入测试数据 ===")
    InjectTestData()

    fmt.Println("\n=== 示例3: 设置认证cookies ===")
    SetupAuthCookies()

    fmt.Println("\n=== 示例4: 创建IndexedDB测试数据 ===")
    SetupIndexedDBTestData()

    fmt.Println("\n=== 示例5: 数据迁移 ===")
    if err := MigrateExtensionData("old-extension", "new-extension"); err != nil {
        log.Printf("数据迁移失败: %v", err)
    }

    fmt.Println("\n=== 示例6: 设置测试环境 ===")
    if err := SetupTestEnvironment(extensionID, "auth_required"); err != nil {
        log.Printf("设置测试环境失败: %v", err)
    }

    fmt.Println("\n=== 示例7: 批量设置存储数据 ===")
    BulkSetStorageData()

    fmt.Println("\n=== 示例8: 实时数据同步 ===")
    syncManager := NewDataSyncManager(extensionID, 1*time.Minute)

    // 添加数据模板
    syncManager.AddDataTemplate("local_storage", func() interface{} {
        return map[string]interface{}{
            "last_updated": time.Now().Unix(),
            "sync_count":   time.Now().Unix() % 100,
        }
    })

    syncManager.AddDataTemplate("cookies", func() interface{} {
        return []Cookie{
            {
                Name:   "timestamp",
                Value:  fmt.Sprintf("%d", time.Now().Unix()),
                Domain: ".example.com",
            },
        }
    })

    syncManager.StartSync()
    time.Sleep(2 * time.Minute) // 让同步运行一段时间
}

*/

// -----------------------------------------------  Extensions.triggerAction  -----------------------------------------------
// === 应用场景 ===
// 1. 扩展控制: 程序化触发扩展的特定动作
// 2. 自动化测试: 自动化测试扩展的功能
// 3. 远程控制: 远程控制扩展的行为
// 4. 功能触发: 触发扩展的特定功能
// 5. 集成测试: 集成测试中模拟用户操作
// 6. 批处理操作: 批量触发多个扩展动作

// CDPExtensionsTriggerAction 触发扩展动作
// extensionID: 扩展ID
// actionName: 动作名称
// parameters: 动作参数
func CDPExtensionsTriggerAction(extensionID, actionName string, parameters map[string]interface{}) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	var parametersJSON string
	if len(parameters) > 0 {
		paramsBytes, err := json.Marshal(parameters)
		if err != nil {
			return "", fmt.Errorf("序列化参数失败: %w", err)
		}
		parametersJSON = fmt.Sprintf(`"parameters": %s,`, string(paramsBytes))
	}

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "Extensions.triggerAction",
        "params": {
            "extensionId": "%s",
            "actionName": "%s",
            %s
            "timeout": 30000
        }
    }`, reqID, extensionID, actionName, parametersJSON)

	// 移除可能的多余逗号
	message = strings.ReplaceAll(message, ",,", ",")
	message = strings.ReplaceAll(message, ",\n        }", "\n        }")

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 triggerAction 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 等待响应
	timeout := 30 * time.Second // 与CDP中的timeout参数保持一致
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
			return "", fmt.Errorf("triggerAction 请求超时")
		}
	}
}

// 辅助函数: CDPExtensionsTriggerActionWithTimeout 自定义超时的触发动作
func CDPExtensionsTriggerActionWithTimeout(extensionID, actionName string, parameters map[string]interface{}, timeout time.Duration) (string, error) {
	if parameters == nil {
		parameters = make(map[string]interface{})
	}
	parameters["timeout"] = timeout.Milliseconds()

	return CDPExtensionsTriggerAction(extensionID, actionName, parameters)
}

// 辅助函数: CDPExtensionsTriggerPageAction 触发页面动作
func CDPExtensionsTriggerPageAction(extensionID, actionName string, tabID int, url string) (string, error) {
	parameters := map[string]interface{}{
		"tabId": tabID,
		"url":   url,
	}

	return CDPExtensionsTriggerAction(extensionID, actionName, parameters)
}

// 辅助函数: CDPExtensionsTriggerBrowserAction 触发浏览器动作
func CDPExtensionsTriggerBrowserAction(extensionID, actionName string) (string, error) {
	return CDPExtensionsTriggerAction(extensionID, actionName, nil)
}

// 辅助函数: CDPExtensionsTriggerActionWithCallback 带回调的触发动作
func CDPExtensionsTriggerActionWithCallback(extensionID, actionName string, parameters map[string]interface{}, callback func(string, error)) {
	go func() {
		result, err := CDPExtensionsTriggerAction(extensionID, actionName, parameters)
		callback(result, err)
	}()
}

/*

// 示例1: 触发广告拦截扩展的更新
func TriggerAdBlockUpdate() {
    extensionID := "adblock-extension-id"

    result, err := CDPExtensionsTriggerAction(extensionID, "updateFilters", map[string]interface{}{
        "force": true,
    })
    if err != nil {
        log.Printf("触发广告拦截更新失败: %v", err)
        return
    }

    log.Printf("广告拦截更新已触发: %s", result)
}

// 示例2: 触发密码管理器的保存操作
func TriggerPasswordSave() {
    extensionID := "password-manager-extension"

    result, err := CDPExtensionsTriggerAction(extensionID, "savePassword", map[string]interface{}{
        "url":      "https://example.com/login",
        "username": "user@example.com",
        "password": "secure_password_123",
        "title":    "Example Site",
    })
    if err != nil {
        log.Printf("触发密码保存失败: %v", err)
        return
    }

    log.Printf("密码保存已触发: %s", result)
}

// 示例3: 触发翻译扩展的翻译操作
func TriggerTranslation() {
    extensionID := "translate-extension-id"

    result, err := CDPExtensionsTriggerAction(extensionID, "translatePage", map[string]interface{}{
        "sourceLang": "en",
        "targetLang": "zh-CN",
        "tabId":      123,
    })
    if err != nil {
        log.Printf("触发翻译失败: %v", err)
        return
    }

    log.Printf("页面翻译已触发: %s", result)
}

// 示例4: 触发截图扩展的截图操作
func TriggerScreenshot() {
    extensionID := "screenshot-extension-id"

    result, err := CDPExtensionsTriggerAction(extensionID, "captureVisibleTab", map[string]interface{}{
        "format":   "png",
        "quality":  90,
        "tabId":    123,
        "callback": "saveScreenshot",
    })
    if err != nil {
        log.Printf("触发截图失败: %v", err)
        return
    }

    // 解析响应获取截图数据
    var response struct {
        Result struct {
            Data   string `json:"data"`
            Format string `json:"format"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err == nil {
        log.Printf("截图完成，格式: %s，数据长度: %d",
            response.Result.Format, len(response.Result.Data))
    }
}

// 示例5: 自动化测试套件
type ExtensionTestSuite struct {
    ExtensionID string
    TestCases   []TestCase
    Results     []TestResult
}

type TestCase struct {
    Name       string
    ActionName string
    Parameters map[string]interface{}
    Expected   interface{}
    Timeout    time.Duration
}

type TestResult struct {
    TestName   string
    Success    bool
    Actual     interface{}
    Error      string
    Duration   time.Duration
}

func (suite *ExtensionTestSuite) Run() {
    log.Printf("开始运行扩展测试套件，扩展ID: %s", suite.ExtensionID)

    for _, testCase := range suite.TestCases {
        log.Printf("运行测试用例: %s", testCase.Name)

        startTime := time.Now()
        result, err := CDPExtensionsTriggerActionWithTimeout(
            suite.ExtensionID,
            testCase.ActionName,
            testCase.Parameters,
            testCase.Timeout,
        )
        duration := time.Since(startTime)

        testResult := TestResult{
            TestName: testCase.Name,
            Duration: duration,
        }

        if err != nil {
            testResult.Success = false
            testResult.Error = err.Error()
            log.Printf("测试用例 %s 失败: %v", testCase.Name, err)
        } else {
            // 验证结果
            testResult.Success = suite.validateResult(result, testCase.Expected)
            testResult.Actual = result
            if testResult.Success {
                log.Printf("测试用例 %s 成功，耗时: %v", testCase.Name, duration)
            } else {
                log.Printf("测试用例 %s 验证失败", testCase.Name)
            }
        }

        suite.Results = append(suite.Results, testResult)
    }

    suite.GenerateReport()
}

func (suite *ExtensionTestSuite) validateResult(actual string, expected interface{}) bool {
    // 这里可以实现具体的验证逻辑
    // 简单示例：检查是否包含期望的关键字
    if expectedStr, ok := expected.(string); ok {
        return actual == expectedStr
    }
    return true
}

func (suite *ExtensionTestSuite) GenerateReport() {
    fmt.Println("=== 扩展测试报告 ===")
    fmt.Printf("扩展ID: %s\n", suite.ExtensionID)
    fmt.Printf("测试时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

    total := len(suite.Results)
    passed := 0

    for _, result := range suite.Results {
        status := "✓ 通过"
        if !result.Success {
            status = "✗ 失败"
        } else {
            passed++
        }

        fmt.Printf("%s %-30s 耗时: %v\n",
            status, result.TestName, result.Duration)

        if result.Error != "" {
            fmt.Printf("   错误: %s\n", result.Error)
        }
    }

    fmt.Printf("\n总计: %d/%d 通过 (%.1f%%)\n",
        passed, total, float64(passed)/float64(total)*100)
}

// 示例6: 批量触发多个扩展动作
func BatchTriggerActions() {
    // 定义要触发的动作列表
    actions := []struct {
        ExtensionID string
        ActionName  string
        Parameters  map[string]interface{}
    }{
        {
            ExtensionID: "extension-1",
            ActionName:  "refresh",
            Parameters:  map[string]interface{}{"force": true},
        },
        {
            ExtensionID: "extension-2",
            ActionName:  "sync",
            Parameters:  map[string]interface{}{"type": "full"},
        },
        {
            ExtensionID: "extension-3",
            ActionName:  "cleanup",
            Parameters:  nil,
        },
    }

    var wg sync.WaitGroup
    results := make(chan string, len(actions))
    errors := make(chan error, len(actions))

    for _, action := range actions {
        wg.Add(1)
        go func(extID, actName string, params map[string]interface{}) {
            defer wg.Done()

            result, err := CDPExtensionsTriggerAction(extID, actName, params)
            if err != nil {
                errors <- fmt.Errorf("扩展 %s 动作 %s 失败: %w",
                    extID, actName, err)
            } else {
                results <- fmt.Sprintf("扩展 %s 动作 %s 成功: %s",
                    extID, actName, result)
            }
        }(action.ExtensionID, action.ActionName, action.Parameters)
    }

    wg.Wait()
    close(results)
    close(errors)

    // 输出结果
    fmt.Println("=== 批量动作执行结果 ===")
    for result := range results {
        fmt.Println(result)
    }

    for err := range errors {
        fmt.Printf("错误: %v\n", err)
    }
}

// 示例7: 动作链执行
type ActionChain struct {
    ExtensionID string
    Actions     []ChainAction
    Context     map[string]interface{}
}

type ChainAction struct {
    Name       string
    Parameters map[string]interface{}
    Condition  func(map[string]interface{}) bool
    OnSuccess  func(string, map[string]interface{})
    OnError    func(error, map[string]interface{})
}

func (chain *ActionChain) Execute() error {
    for _, action := range chain.Actions {
        // 检查执行条件
        if action.Condition != nil && !action.Condition(chain.Context) {
            log.Printf("跳过动作: %s，条件不满足", action.Name)
            continue
        }

        log.Printf("执行动作: %s", action.Name)
        result, err := CDPExtensionsTriggerAction(chain.ExtensionID, action.Name, action.Parameters)

        if err != nil {
            log.Printf("动作 %s 失败: %v", action.Name, err)
            if action.OnError != nil {
                action.OnError(err, chain.Context)
            }
            return err
        }

        log.Printf("动作 %s 成功: %s", action.Name, result)
        if action.OnSuccess != nil {
            action.OnSuccess(result, chain.Context)
        }

        // 更新上下文
        chain.Context["last_action"] = action.Name
        chain.Context["last_result"] = result
    }

    return nil
}

// 示例8: 监控扩展动作的执行
type ActionMonitor struct {
    ExtensionID    string
    ActionPatterns []string
    EventChannel   chan ActionEvent
    IsMonitoring   bool
}

type ActionEvent struct {
    Timestamp  time.Time
    ActionName string
    Parameters map[string]interface{}
    Result     string
    Error      string
    Duration   time.Duration
}

func (m *ActionMonitor) Start() {
    m.IsMonitoring = true
    m.EventChannel = make(chan ActionEvent, 100)

    go m.monitorLoop()
}

func (m *ActionMonitor) Stop() {
    m.IsMonitoring = false
    if m.EventChannel != nil {
        close(m.EventChannel)
    }
}

func (m *ActionMonitor) monitorLoop() {
    // 这里可以实现具体的监控逻辑
    // 例如：定期检查扩展状态，记录动作执行等
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for m.IsMonitoring {
        select {
        case <-ticker.C:
            // 模拟监控事件
            event := ActionEvent{
                Timestamp:  time.Now(),
                ActionName: "heartbeat",
                Duration:   0,
            }
            m.EventChannel <- event
        }
    }
}

func (m *ActionMonitor) TriggerAndMonitor(actionName string, parameters map[string]interface{}) (ActionEvent, error) {
    startTime := time.Now()
    result, err := CDPExtensionsTriggerAction(m.ExtensionID, actionName, parameters)
    duration := time.Since(startTime)

    event := ActionEvent{
        Timestamp:  startTime,
        ActionName: actionName,
        Parameters: parameters,
        Result:     result,
        Error:      "",
        Duration:   duration,
    }

    if err != nil {
        event.Error = err.Error()
    }

    m.EventChannel <- event
    return event, err
}



*/

// -----------------------------------------------  Extensions.uninstall  -----------------------------------------------
// === 应用场景 ===
// 1. 扩展清理: 清理测试过程中安装的扩展程序
// 2. 环境重置: 恢复浏览器到干净的初始状态
// 3. 扩展管理: 程序化卸载不需要的扩展
// 4. 故障恢复: 卸载有问题的扩展程序
// 5. 版本更新: 卸载旧版本扩展准备安装新版本
// 6. 权限管理: 移除有安全风险的扩展程序

// CDPExtensionsUninstall 卸载扩展程序
func CDPExtensionsUninstall(extensionID string) (string, error) {
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
        "method": "Extensions.uninstall",
        "params": {
            "extensionId": "%s"
        }
    }`, reqID, extensionID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 uninstall 请求失败: %w", err)
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
			return "", fmt.Errorf("uninstall 请求超时")
		}
	}
}
