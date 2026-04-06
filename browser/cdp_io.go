package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  IO.close  -----------------------------------------------
// === 应用场景 ===
// 1. 资源释放: 读取完文件/流数据后立即关闭句柄，避免浏览器资源泄漏
// 2. 自动化清理: 自动化测试中读取页面资源后，规范关闭IO句柄
// 3. 错误处理: 读取流失败时，主动关闭句柄防止资源占用
// 4. 批量读取: 批量读取多个文件流后，逐个关闭句柄
// 5. 调试辅助: 调试IO流操作时，确保句柄正常关闭
// 6. 页面资源回收: 读取完页面网络流、日志流后回收句柄

// CDPIOClose 关闭IO流句柄
// handle: IO.open返回的流句柄
func CDPIOClose(handle string) (string, error) {
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
		"method": "IO.close",
		"params": {
			"handle": "%s"
		}
	}`, reqID, handle)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 IO.close 请求失败: %w", err)
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
			return "", fmt.Errorf("IO.close 请求超时")
		}
	}
}

/*

// ==================== IO.close 使用示例 ====================
func ExampleCDPIOClose() {
	// 1. 假设已通过 IO.open 获取到流句柄
	streamHandle := "123456" // 实际为IO.open返回的handle值

	// 2. 读取完流数据后，关闭句柄
	resp, err := CDPIOClose(streamHandle)
	if err != nil {
		log.Fatalf("关闭IO句柄失败: %v", err)
	}

	log.Printf("关闭IO句柄成功，响应: %s", resp)

	// 3. 错误处理示例：读取失败时强制关闭

		readErr := CDPIORead(streamHandle)
		if readErr != nil {
			// 读取失败，必须关闭句柄
			CDPIOClose(streamHandle)
			log.Fatalf("读取失败，已关闭句柄: %v", readErr)
		}

}

*/

// -----------------------------------------------  IO.read  -----------------------------------------------
// === 应用场景 ===
// 1. 页面资源读取: 读取页面加载的JS/CSS/图片等二进制/文本资源
// 2. 日志流获取: 读取浏览器控制台输出、网络日志流数据
// 3. 文件流读取: 读取通过IO.open打开的本地/远程文件流
// 4. 自动化数据采集: 自动化测试中抓取页面资源原始内容
// 5. 调试分析: 调试页面资源加载异常时，查看原始资源数据
// 6. 大文件分片读取: 大体积资源分多次读取，避免内存溢出

// CDPIORead 从IO流句柄读取数据
// handle: IO.open获取的流句柄
// size: 可选，要读取的字节大小，传0则使用默认值
// offset: 可选，读取偏移量，传0则从开头读取
func CDPIORead(handle string, size int, offset int) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数（size/offset为可选参数）
	params := ""
	if size > 0 && offset > 0 {
		params = fmt.Sprintf(`"handle": "%s", "size": %d, "offset": %d`, handle, size, offset)
	} else if size > 0 {
		params = fmt.Sprintf(`"handle": "%s", "size": %d`, handle, size)
	} else if offset > 0 {
		params = fmt.Sprintf(`"handle": "%s", "offset": %d`, handle, offset)
	} else {
		params = fmt.Sprintf(`"handle": "%s"`, handle)
	}

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "IO.read",
		"params": {%s}
	}`, reqID, params)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 IO.read 请求失败: %w", err)
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
			return "", fmt.Errorf("IO.read 请求超时")
		}
	}
}

/*

// ==================== IO.read 使用示例 ====================
func ExampleCDPIORead() {
	// 1. 必备前提：已通过 IO.open 获取到有效流句柄
	streamHandle := "123456" // 实际为IO.open返回的handle

	// ========== 示例1：基础读取（读取全部数据） ==========
	resp, err := CDPIORead(streamHandle, 0, 0)
	if err != nil {
		log.Fatalf("IO流读取失败: %v", err)
	}
	log.Printf("IO流基础读取成功，响应: %s", resp)

	// ========== 示例2：指定大小读取（限制1024字节） ==========
	resp2, err2 := CDPIORead(streamHandle, 1024, 0)
	if err2 != nil {
		log.Fatalf("指定大小IO读取失败: %v", err2)
	}
	log.Printf("指定大小读取成功: %s", resp2)

	// ========== 示例3：偏移量读取（从第512字节开始读） ==========
	resp3, err3 := CDPIORead(streamHandle, 512, 512)
	if err3 != nil {
		log.Fatalf("偏移量IO读取失败: %v", err3)
	}
	log.Printf("偏移量读取成功: %s", resp3)

	// ========== 示例4：读取+关闭标准流程 ==========

	// 读取完成后必须关闭句柄，释放资源
	resp4, _ := CDPIORead(streamHandle, 0, 0)
	_, closeErr := CDPIOClose(streamHandle)
	if closeErr != nil {
		log.Printf("关闭句柄失败: %v", closeErr)
	}

}

*/
