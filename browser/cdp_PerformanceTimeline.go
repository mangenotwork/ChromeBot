package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  PerformanceTimeline.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 性能时间线监控开启: 启动接收页面完整性能时间线事件
// 2. 加载性能分析: 采集页面从请求到渲染的全链路时间线数据
// 3. 自动化性能测试: 监听性能时间线事件用于性能瓶颈定位
// 4. 前端渲染优化: 分析脚本、样式、渲染、绘制阶段耗时
// 5. 页面卡顿诊断: 追踪长任务、重排重绘、长时间阻塞事件
// 6. 性能日志采集: 收集浏览器性能时间线详细事件用于离线分析

// CDPPerformanceTimelineEnable 启用PerformanceTimeline域，开始接收性能时间线事件
func CDPPerformanceTimelineEnable() (string, error) {
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
		"method": "PerformanceTimeline.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 PerformanceTimeline.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("PerformanceTimeline.enable 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：启动性能时间线监控，采集页面加载与渲染全流程事件
func ExamplePerformanceTimelineEnable() {
	// 1. 建立浏览器CDP连接（前置逻辑）
	// err := ConnectBrowserCDP()
	// if err != nil {
	// 	log.Fatalf("浏览器连接失败: %v", err)
	// }

	// 2. 启用性能时间线事件监听
	resp, err := CDPPerformanceTimelineEnable()
	if err != nil {
		log.Printf("启用性能时间线域失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功启用性能时间线域，开始接收页面加载、渲染、长任务等性能事件")

	// 后续可监听 messageQueue 中的 timelineEvent 事件，获取详细性能时间线数据
}
*/
