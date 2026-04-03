package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  FileSystem.getDirectory  -----------------------------------------------
// === 应用场景 ===
// 1. 目录获取: 获取指定目录信息
// 2. 路径解析: 解析目录路径
// 3. 目录浏览: 浏览目录内容
// 4. 结构分析: 分析目录结构
// 5. 权限检查: 检查目录访问权限
// 6. 元数据获取: 获取目录元数据

// CDPFileSystemGetDirectory 获取目录
// fileSystemID: 文件系统ID
// path: 目录路径
func CDPFileSystemGetDirectory(fileSystemID, path string) (string, error) {
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
        "method": "FileSystem.getDirectory",
        "params": {
            "fileSystemId": "%s",
            "path": "%s"
        }
    }`, reqID, fileSystemID, path)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FileSystem.getDirectory 请求失败: %w", err)
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
			return "", fmt.Errorf("FileSystem.getDirectory 请求超时")
		}
	}
}

/*


// 示例: 获取特定目录
func ExampleCDPFileSystemGetDirectory() {
	// 启用文件系统
	if _, err := CDPFileSystemEnable(); err != nil {
		log.Printf("启用文件系统失败: %v", err)
		return
	}

	// 假设有文件系统ID
	fileSystemID := "filesystem:https://example.com/persistent/0"
	directoryPath := "/documents"

	result, err := CDPFileSystemGetDirectory(fileSystemID, directoryPath)
	if err != nil {
		log.Printf("获取目录失败: %v", err)
		return
	}

	// 解析目录信息
	var resp struct {
		Result struct {
			Name        string `json:"name"`
			Path        string `json:"path"`
			IsFile      bool   `json:"isFile"`
			IsDirectory bool   `json:"isDirectory"`
			Size        int64  `json:"size,omitempty"`
			Modified    int64  `json:"modified,omitempty"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(result), &resp); err == nil {
		if resp.Result.IsDirectory {
			log.Printf("获取到目录: 名称=%s, 路径=%s",
				resp.Result.Name, resp.Result.Path)

			// 可以进一步操作目录，如列出内容
		} else {
			log.Printf("路径不是目录: %s", directoryPath)
		}
	} else {
		log.Printf("获取目录信息: %s", result)
	}
}

*/
