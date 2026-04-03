package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  CacheStorage.deleteCache  -----------------------------------------------
// === 应用场景 ===
// 1. **缓存清理测试**：在自动化测试中清除特定缓存，验证页面加载时是否重新获取最新资源
// 2. **用户数据管理**：浏览器扩展开发时，允许用户手动清除指定缓存（如隐私数据缓存）
// 3. **性能优化验证**：分析缓存使用情况后，清理无效或过大的缓存条目
// 4. **服务端更新同步**：当服务端资源更新时，强制清除客户端缓存获取最新内容
// 5. **调试辅助**：开发阶段手动清除缓存，避免因缓存导致无法复现的问题
// 6. **安全审计**：清除可能包含敏感信息的缓存，防止数据泄露风险

// CDPCacheStorageDeleteCache 删除指定缓存
func CDPCacheStorageDeleteCache(cacheID string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的请求消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "CacheStorage.deleteCache",
        "params": {
            "cacheId": "%s"
        }
    }`, reqID, cacheID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deleteCache 请求失败: %w", err)
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
			return "", fmt.Errorf("deleteCache 请求超时")
		}
	}
}

/*

示例1：测试场景 - 清除缓存后验证资源加载

func TestCacheUpdate() {
    // 1. 先获取所有缓存名称
    cacheNames, err := CDPCacheStorageRequestCacheNames()
    if err != nil {
        log.Fatalf("获取缓存失败: %v", err)
    }

    // 2. 假设我们有一个已知的缓存ID（实际可通过缓存名称映射获取）
    targetCacheID := "test-cache-123"

    // 3. 删除指定缓存
    result, err := CDPCacheStorageDeleteCache(targetCacheID)
    if err != nil {
        log.Fatalf("删除缓存失败: %v", err)
    }
    log.Printf("删除结果: %s", result)

    // 4. 重新加载页面验证是否获取新资源
    // ...（此处可添加页面导航和资源验证逻辑）
}

示例2：安全场景 - 清除敏感缓存
func ClearSensitiveCache() {
    // 1. 获取所有缓存名称
    cacheNames, err := CDPCacheStorageRequestCacheNames()
    if err != nil {
        log.Fatalf("获取缓存失败: %v", err)
    }

    // 2. 遍历并删除可能包含敏感信息的缓存
    for _, cache := range cacheNames.Caches {
        if strings.Contains(cache.CacheName, "auth-") { // 假设以auth-开头的缓存包含认证信息
            _, err := CDPCacheStorageDeleteCache(cache.CacheId)
            if err != nil {
                log.Printf("警告: 删除缓存 %s 失败: %v", cache.CacheId, err)
                continue
            }
            log.Printf("已清除敏感缓存: %s", cache.CacheName)
        }
    }
}


示例3：性能优化 - 清理过期缓存
func CleanupExpiredCaches() {
    // 1. 获取所有缓存名称
    cacheNames, err := CDPCacheStorageRequestCacheNames()
    if err != nil {
        log.Fatalf("获取缓存失败: %v", err)
    }

    // 2. 定义清理阈值（例如超过100MB的缓存）
    sizeThreshold := int64(100 * 1024 * 1024) // 100MB

    // 3. 遍历并删除大缓存
    for _, cache := range cacheNames.Caches {
        stats, err := CDPCacheStorageRequestCachedResponse(cache.CacheId, "/*") // 获取缓存统计
        if err != nil {
            log.Printf("获取缓存统计失败: %v", err)
            continue
        }

        // 假设统计中有size字段（实际需根据协议返回结构调整）
        if stats.BodySize > sizeThreshold {
            _, err := CDPCacheStorageDeleteCache(cache.CacheId)
            if err != nil {
                log.Printf("删除大缓存失败: %v", err)
                continue
            }
            log.Printf("已清理大缓存: %s (%.2fMB)",
                cache.CacheName,
                float64(stats.BodySize)/1024/1024)
        }
    }
}


*/

// -----------------------------------------------  CacheStorage.deleteEntry  -----------------------------------------------
// === 应用场景 ===
// 1. 精确缓存清理：清除特定URL的缓存条目（如API响应）
// 2. 缓存更新验证：在服务端更新后强制客户端获取新内容
// 3. 缓存污染修复：删除被错误缓存的无效响应
// 4. 性能优化：清理长期未访问的缓存条目
// 5. 安全审计：清除可能包含敏感信息的特定缓存条目
// 6. 调试辅助：验证特定缓存条目是否被正确存储

// CDPCacheStorageDeleteEntry 删除指定缓存条目
func CDPCacheStorageDeleteEntry(cacheID, request string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的请求消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "CacheStorage.deleteEntry",
        "params": {
            "cacheId": "%s",
            "request": "%s"
        }
    }`, reqID, cacheID, request)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 deleteEntry 请求失败: %w", err)
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
			return "", fmt.Errorf("deleteEntry 请求超时")
		}
	}
}

/*

示例1：清除特定API缓存
func ClearApiCache() {
    // 1. 获取目标缓存ID
    cacheID := "api-cache-v1"

    // 2. 构建要删除的请求URL（需要URL编码）
    requestURL := "https%3A%2F%2Fapi.example.com%2Fusers%2F123"

    // 3. 执行删除操作
    result, err := CDPCacheStorageDeleteEntry(cacheID, requestURL)
    if err != nil {
        log.Printf("删除失败: %v", err)
        return
    }

    log.Printf("成功删除缓存条目: %s", result)

    // 4. 验证结果（重新请求API验证是否获取最新数据）
    // ...
}


示例2：批量清理过期缓存
func CleanExpiredCacheEntries(cacheID string) {
    // 1. 获取缓存中的所有条目
    entries, err := CDPCacheStorageRequestEntries(cacheID, 0, 100) // 假设有分页获取方法
    if err != nil {
        log.Fatalf("获取缓存条目失败: %v", err)
    }

    // 2. 定义过期时间（24小时前）
    expireTime := time.Now().Add(-24 * time.Hour)

    // 3. 遍历并删除过期条目
    for _, entry := range entries {
        // 解析entry.response中的时间戳
        timestamp, err := parseCacheTimestamp(entry.Response)
        if err != nil {
            continue // 跳过无法解析的条目
        }

        // 4. 检查是否过期
        if timestamp.Before(expireTime) {
            _, err := CDPCacheStorageDeleteEntry(cacheID, entry.Request)
            if err != nil {
                log.Printf("删除条目失败: %v", err)
                continue
            }
            log.Printf("已清理过期缓存: %s", entry.Request)
        }
    }
}


示例3：调试缓存问题
func DebugCacheIssue() {
    // 1. 定位问题缓存
    cacheID := "debug-cache-1"
    problemURL := "https://example.com/buggy-page"

    // 2. 删除特定缓存条目
    result, err := CDPCacheStorageDeleteEntry(cacheID, problemURL)
    if err != nil {
        log.Fatalf("删除失败: %v", err)
    }

    // 3. 重新加载页面验证问题是否解决
    if err := ReloadPage(); err != nil {
        log.Fatalf("页面重载失败: %v", err)
    }

    // 4. 验证问题是否解决
    if IsPageFixed() {
        log.Println("缓存问题已确认")
    }
}


*/

// -----------------------------------------------  CacheStorage.requestCachedResponse  -----------------------------------------------
// === 应用场景 ===
// 1. 缓存内容验证：检查特定缓存条目的实际响应内容
// 2. 缓存策略调试：分析缓存响应头确认缓存策略是否生效
// 3. 性能分析：获取缓存响应体大小和头部信息用于性能评估
// 4. 安全审计：检查缓存内容是否包含敏感信息
// 5. 问题诊断：验证缓存是否包含预期的响应数据
// 6. 自动化测试：在测试中验证缓存内容是否符合预期

// CDPCacheStorageRequestCachedResponse 获取指定缓存条目的响应内容
func CDPCacheStorageRequestCachedResponse(cacheID, request string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的请求消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "CacheStorage.requestCachedResponse",
        "params": {
            "cacheId": "%s",
            "requestURL": "%s"
        }
    }`, reqID, cacheID, request)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 requestCachedResponse 请求失败: %w", err)
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

				// 解析响应结构
				var response struct {
					Result struct {
						Response struct {
							Status     int         `json:"status"`
							StatusText string      `json:"statusText"`
							Headers    interface{} `json:"headers"`
							Content    struct {
								BodySize int         `json:"bodySize"`
								Content  interface{} `json:"content"`
							} `json:"content"`
						} `json:"response"`
					} `json:"result"`
					Error interface{} `json:"error"`
				}

				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				if response.Error != nil {
					return content, fmt.Errorf("CDP错误: %v", response.Error)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("requestCachedResponse 请求超时")
		}
	}
}

/*

示例1：验证API响应缓存
func VerifyApiCache() {
    cacheID := "api-cache-v1"
    targetURL := "https%3A%2F%2Fapi.example.com%2Fusers%2F123"

    // 获取缓存响应
    response, err := CDPCacheStorageRequestCachedResponse(cacheID, targetURL)
    if err != nil {
        log.Fatalf("获取缓存失败: %v", err)
    }

    // 解析响应内容
    var cachedResponse struct {
        Status     int
        StatusText string
        Headers    map[string]string
        Body       string
    }
    if err := json.Unmarshal([]byte(response), &cachedResponse); err != nil {
        log.Fatalf("解析缓存响应失败: %v", err)
    }

    // 验证状态码
    if cachedResponse.Status != 200 {
        log.Fatalf("非预期状态码: %d", cachedResponse.Status)
    }

    // 验证响应体内容
    if !strings.Contains(cachedResponse.Body, "user-123") {
        log.Fatalf("响应体缺少预期内容")
    }

    log.Println("缓存验证通过")
}


示例2：检查缓存策略生效情况
func CheckCacheHeaders() {
    cacheID := "static-assets-cache"
    targetURL := "https%3A%2F%2Fexample.com%2Fstyles.css"

    response, err := CDPCacheStorageRequestCachedResponse(cacheID, targetURL)
    if err != nil {
        log.Fatalf("获取缓存失败: %v", err)
    }

    // 解析响应头
    var headers map[string]interface{}
    if err := json.Unmarshal([]byte(response), &headers); err != nil {
        log.Fatalf("解析响应头失败: %v", err)
    }

    // 检查缓存控制头
    cacheControl, exists := headers["cache-control"]
    if !exists {
        log.Fatalf("缺少Cache-Control头")
    }

    // 验证缓存策略
    if cacheControl != "public, max-age=31536000" {
        log.Fatalf("非预期的缓存策略: %s", cacheControl)
    }

    log.Println("缓存策略验证通过")
}


示例3：安全审计缓存内容
func AuditCacheForSensitiveData() {
    cacheID := "user-data-cache"
    targetURL := "https%3A%2F%2Fexample.com%2Fprofile"

    response, err := CDPCacheStorageRequestCachedResponse(cacheID, targetURL)
    if err != nil {
        log.Printf("警告: 获取缓存失败: %v", err)
        return
    }

    // 检查响应体是否包含敏感信息
    if containsSensitiveData(response) {
        log.Printf("警告: 缓存包含敏感信息 - %s", targetURL)
        // 执行清理操作
        _, err := CDPCacheStorageDeleteEntry(cacheID, targetURL)
        if err != nil {
            log.Printf("清理失败: %v", err)
        }
    }
}

func containsSensitiveData(response string) bool {
    // 实现敏感信息检测逻辑（如检测密码、令牌等）
    return strings.Contains(response, "password") ||
           strings.Contains(response, "auth_token")
}



*/

// -----------------------------------------------  CacheStorage.requestCacheNames  -----------------------------------------------
// === 应用场景 ===
// 1. 缓存清单管理：获取所有缓存名称进行批量管理
// 2. 自动化测试：验证缓存创建/删除后的状态变化
// 3. 性能分析：统计缓存数量评估存储使用情况
// 4. 调试辅助：开发阶段查看当前页面缓存状态
// 5. 安全审计：检测是否存在异常缓存实例
// 6. 缓存策略验证：确认缓存是否按预期创建

// CDPCacheStorageRequestCacheNames 获取所有缓存名称
func CDPCacheStorageRequestCacheNames() (struct {
	CacheIds []string `json:"cacheIds"`
}, error) {
	if !DefaultBrowserWS() {
		return struct {
			CacheIds []string `json:"cacheIds"`
		}{}, fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return struct {
			CacheIds []string `json:"cacheIds"`
		}{}, fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "CacheStorage.requestCacheNames"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return struct {
			CacheIds []string `json:"cacheIds"`
		}{}, fmt.Errorf("发送 requestCacheNames 请求失败: %w", err)
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
				return struct {
					CacheIds []string `json:"cacheIds"`
				}{}, fmt.Errorf("消息队列已关闭")
			}

			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应结构
				var response struct {
					Result struct {
						CacheIds []string `json:"cacheIds"`
					} `json:"result"`
					Error interface{} `json:"error"`
				}

				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return struct {
						CacheIds []string `json:"cacheIds"`
					}{}, fmt.Errorf("解析响应失败: %w", err)
				}

				if response.Error != nil {
					return struct {
						CacheIds []string `json:"cacheIds"`
					}{}, fmt.Errorf("CDP错误: %v", response.Error)
				}

				return response.Result, nil
			}

		case <-timer.C:
			return struct {
				CacheIds []string `json:"cacheIds"`
			}{}, fmt.Errorf("requestCacheNames 请求超时")
		}
	}
}

/*

示例1：列出所有缓存名称
func ListAllCaches() {
    // 获取所有缓存名称
    cacheNames, err := CDPCacheStorageRequestCacheNames()
    if err != nil {
        log.Fatalf("获取缓存名称失败: %v", err)
    }

    // 打印结果
    log.Printf("发现 %d 个缓存实例:", len(cacheNames.CacheIds))
    for i, cacheID := range cacheNames.CacheIds {
        log.Printf("%d. %s", i+1, cacheID)
    }

    // 后续可进行批量操作（如清理、分析等）
}



示例2：自动化测试缓存创建
func TestCacheCreation() {
    // 测试前获取缓存列表
    before, err := CDPCacheStorageRequestCacheNames()
    if err != nil {
        log.Fatalf("初始缓存获取失败: %v", err)
    }

    // 执行创建缓存的操作（如加载特定页面）
    LoadTestPage()

    // 测试后获取缓存列表
    after, err := CDPCacheStorageRequestCacheNames()
    if err != nil {
        log.Fatalf("测试后缓存获取失败: %v", err)
    }

    // 验证是否新增了缓存
    newCaches := make(map[string]bool)
    for _, id := range after.CacheIds {
        newCaches[id] = true
    }
    for _, id := range before.CacheIds {
        delete(newCaches, id)
    }

    if len(newCaches) == 0 {
        log.Fatalf("未检测到新缓存创建")
    }

    log.Printf("成功创建 %d 个新缓存", len(newCaches))
}



示例3：缓存使用情况分析
func AnalyzeCacheUsage() {
    // 获取所有缓存
    caches, err := CDPCacheStorageRequestCacheNames()
    if err != nil {
        log.Fatalf("获取缓存失败: %v", err)
    }

    // 统计缓存数量和大小
    totalCaches := len(caches.CacheIds)
    totalSize := int64(0)

    // 遍历获取每个缓存的大小
    for _, cacheID := range caches.CacheIds {
        stats, err := CDPCacheStorageRequestCachedResponse(cacheID, "/*")
        if err != nil {
            log.Printf("获取缓存统计失败: %v", err)
            continue
        }
        // 假设响应中包含bodySize字段
        if size, exists := stats["bodySize"]; exists {
            totalSize += int64(size.(float64))
        }
    }

    log.Printf("缓存统计: %d 个缓存实例, 总大小 %.2fMB",
        totalCaches,
        float64(totalSize)/1024/1024)
}

*/

// -----------------------------------------------  CacheStorage.requestEntries  -----------------------------------------------
// === 应用场景 ===
// 1. 缓存批量管理：获取指定缓存的所有条目进行批量操作
// 2. 分页加载：支持分页获取缓存条目，避免一次性加载过多数据
// 3. 路径过滤：根据URL路径过滤缓存条目
// 4. 缓存分析：统计缓存使用情况，如最常访问的资源
// 5. 自动化测试：验证缓存条目是否符合预期

// CDPCacheStorageRequestEntries 获取缓存条目列表
func CDPCacheStorageRequestEntries(cacheID string, skipCount, pageSize int, pathFilter string) (RequestEntriesResult, error) {
	if !DefaultBrowserWS() {
		return RequestEntriesResult{}, fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return RequestEntriesResult{}, fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建请求参数
	params := map[string]interface{}{
		"cacheId": cacheID,
	}
	if skipCount > 0 {
		params["skipCount"] = skipCount
	}
	if pageSize > 0 {
		params["pageSize"] = pageSize
	}
	if pathFilter != "" {
		params["pathFilter"] = pathFilter
	}

	// 序列化参数
	paramsBytes, _ := json.Marshal(params)

	// 构建请求消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "CacheStorage.requestEntries",
        "params": %s
    }`, reqID, paramsBytes)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return RequestEntriesResult{}, fmt.Errorf("发送请求失败: %w", err)
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
				return RequestEntriesResult{}, fmt.Errorf("消息队列已关闭")
			}
			if respMsg.ID == reqID {
				// 解析响应
				var response struct {
					Result RequestEntriesResult `json:"result"`
					Error  interface{}          `json:"error"`
				}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return RequestEntriesResult{}, fmt.Errorf("解析响应失败: %w", err)
				}
				if response.Error != nil {
					return RequestEntriesResult{}, fmt.Errorf("CDP错误: %v", response.Error)
				}
				return response.Result, nil
			}
		case <-timer.C:
			return RequestEntriesResult{}, fmt.Errorf("请求超时")
		}
	}
}

// 定义缓存条目结构体
type CacheDataEntry struct {
	RequestURL      string                 `json:"requestURL"`
	RequestHeaders  map[string]interface{} `json:"requestHeaders"`
	ResponseHeaders map[string]interface{} `json:"responseHeaders"`
	Timestamp       float64                `json:"timestamp"`
}

// 定义请求结果结构体
type RequestEntriesResult struct {
	CacheDataEntries []CacheDataEntry `json:"cacheDataEntries"`
	ReturnCount      int              `json:"returnCount"`
}

/*

示例1：获取所有缓存条目
cacheID := "test-cache-123"
result, err := CDPCacheStorageRequestEntries(cacheID, 0, 0, "")
if err != nil {
    log.Printf("获取缓存条目失败: %v", err)
    return
}
log.Printf("获取到 %d 个缓存条目", len(result.CacheDataEntries))



示例2：分页获取缓存条目
cacheID := "large-cache-456"
pageSize := 50
skipCount := 0
total := 0

for {
    result, err := CDPCacheStorageRequestEntries(cacheID, skipCount, pageSize, "")
    if err != nil {
        log.Printf("分页获取失败: %v", err)
        break
    }
    total += len(result.CacheDataEntries)
    log.Printf("第 %d 页: %d 条", skipCount/pageSize+1, len(result.CacheDataEntries))
    if len(result.CacheDataEntries) < pageSize {
        break
    }
    skipCount += pageSize
}
log.Printf("总共获取 %d 条", total)


示例3：根据路径过滤缓存条目
cacheID := "api-cache-789"
pathFilter := "/api/users"
result, err := CDPCacheStorageRequestEntries(cacheID, 0, 0, pathFilter)
if err != nil {
    log.Printf("过滤获取失败: %v", err)
    return
}
log.Printf("过滤后获取 %d 条", len(result.CacheDataEntries))
for _, entry := range result.CacheDataEntries {
    if strings.Contains(entry.RequestURL, pathFilter) {
        log.Printf("匹配条目: %s", entry.RequestURL)
    }
}


*/
