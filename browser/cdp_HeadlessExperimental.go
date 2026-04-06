package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  HeadlessExperimental.beginFrame  -----------------------------------------------
// === 应用场景 ===
// 1. 无头模式截图/录屏：在Chrome无头模式下主动触发一帧渲染，用于生成页面截图、视频录制
// 2. 自动化渲染控制：精准控制页面渲染时机，确保页面元素完全渲染后执行后续操作
// 3. 性能测试：手动触发帧渲染，测试页面渲染性能、帧率指标
// 4. 页面状态捕获：在无头环境下获取页面最新渲染帧，保存页面可视化状态
// 5. 自动化测试：测试流程中强制渲染页面，验证UI渲染结果正确性
// 6. 服务端渲染输出：服务端无头浏览器环境下生成渲染帧，用于PDF/图片导出

// CDPHeadlessExperimentalBeginFrame 主动触发无头模式下的一帧渲染
// 参数说明：
//
//	frameTimeTicks: 帧时间戳（毫秒，可选，传0使用当前时间）
//	includeDamage: 是否包含损坏区域（可选，默认false）
func CDPHeadlessExperimentalBeginFrame(frameTimeTicks float64, includeDamage bool) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建带参数的CDP请求消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "HeadlessExperimental.beginFrame",
		"params": {
			"frameTimeTicks": %f,
			"includeDamage": %t
		}
	}`, reqID, frameTimeTicks, includeDamage)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 beginFrame 请求失败: %w", err)
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
			return "", fmt.Errorf("beginFrame 请求超时")
		}
	}
}

/*
// === 使用示例1：基础调用（默认帧渲染） ===
func ExampleBeginFrame_Base() {
	// 使用当前时间戳，不包含损坏区域
	resp, err := CDPHeadlessExperimentalBeginFrame(0, false)
	if err != nil {
		log.Fatalf("触发帧渲染失败: %v", err)
	}
	log.Printf("渲染响应: %s", resp)
}

// === 使用示例2：截图前置调用（确保页面完全渲染） ===
func ExampleBeginFrame_ForScreenshot() {
	// 截图前主动触发渲染，保证画面完整
	_, err := CDPHeadlessExperimentalBeginFrame(0, true)
	if err != nil {
		log.Printf("帧渲染失败: %v", err)
		return
	}
	// 此处可调用Page.captureScreenshot进行截图
	log.Println("帧渲染完成，准备截图...")
}

// === 使用示例3：自定义时间戳渲染 ===
func ExampleBeginFrame_CustomTime() {
	// 使用自定义时间戳（毫秒）
	customTime := float64(time.Now().UnixMilli())
	resp, err := CDPHeadlessExperimentalBeginFrame(customTime, false)
	if err != nil {
		log.Fatalf("渲染失败: %v", err)
	}
	log.Println("自定义时间戳渲染完成:", resp)
}
*/
