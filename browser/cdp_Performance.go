package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Performance.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 性能监控关闭: 停止接收性能相关事件通知
// 2. 资源释放: 关闭性能域后释放浏览器监听与采集资源
// 3. 测试收尾: 自动化性能测试完成后关闭监控
// 4. 性能优化: 无需性能数据时关闭以减少浏览器性能消耗
// 5. 功能切换: 从性能监控模式切换到其他功能模式
// 6. 异常恢复: 性能监控异常时关闭并重新初始化

// CDPPerformanceDisable 关闭Performance域，停止接收性能相关事件
func CDPPerformanceDisable() (string, error) {
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
		"method": "Performance.disable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Performance.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("Performance.disable 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：性能测试结束后关闭性能监控，释放浏览器资源
func ExamplePerformanceDisable() {
	// 1. 先启用性能监控（业务逻辑）
	// resp, err := CDPPerformanceEnable()
	// if err != nil {
	// 	log.Fatalf("启用性能监控失败: %v", err)
	// }

	// 2. 执行性能测试、采集指标等逻辑...

	// 3. 测试完成后关闭性能监控
	resp, err := CDPPerformanceDisable()
	if err != nil {
		log.Printf("关闭性能监控失败: %v, 响应内容: %s", err, resp)
		return
	}
	log.Println("成功关闭性能域，停止接收性能事件")
}
*/

// -----------------------------------------------  Performance.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 性能监控开启: 启动接收页面性能相关事件和指标
// 2. 页面性能采集: 收集加载、渲染、脚本执行等性能数据
// 3. 自动化性能测试: 监控页面卡顿、帧率、耗时等指标
// 4. 前端性能诊断: 定位页面加载慢、交互卡顿的问题
// 5. 性能基线建立: 采集正常性能数据作为对比基准
// 6. 实时性能调试: 调试页面运行时的性能损耗情况

// CDPPerformanceEnable 启用Performance域，开始接收性能相关事件
func CDPPerformanceEnable() (string, error) {
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
		"method": "Performance.enable"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Performance.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("Performance.enable 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：启动页面性能监控，准备采集页面加载与运行时性能数据
func ExamplePerformanceEnable() {
	// 1. 建立浏览器CDP连接（前置逻辑）
	// err := ConnectBrowserCDP()
	// if err != nil {
	// 	log.Fatalf("浏览器连接失败: %v", err)
	// }

	// 2. 启用性能事件监听
	resp, err := CDPPerformanceEnable()
	if err != nil {
		log.Printf("启用性能域失败: %v, 响应: %s", err, resp)
		return
	}
	log.Println("成功启用性能域，开始接收页面性能事件")

	// 后续可监听 messageQueue 中的性能事件，如 metrics、timeline 事件
}
*/

// -----------------------------------------------  Performance.getMetrics  -----------------------------------------------
// === 应用场景 ===
// 1. 实时性能采集: 获取当前页面核心性能指标（帧率、CPU、内存等）
// 2. 自动化性能测试: 采集关键性能指标用于性能回归校验
// 3. 页面健康监控: 实时监控页面流畅度、资源占用情况
// 4. 性能诊断: 快速定位页面卡顿、高CPU占用问题
// 5. 性能对比: 不同版本/操作后的性能指标差异对比
// 6. 监控告警: 采集指标并判断是否超出阈值触发告警

// CDPPerformanceGetMetrics 获取当前页面的性能指标集合
func CDPPerformanceGetMetrics() (string, error) {
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
		"method": "Performance.getMetrics"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 Performance.getMetrics 请求失败: %w", err)
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
			return "", fmt.Errorf("Performance.getMetrics 请求超时")
		}
	}
}

/*
// === 使用场景示例代码 ===
// 场景：页面交互完成后采集性能指标，检测是否存在性能瓶颈
func ExamplePerformanceGetMetrics() {
	// 1. 确保已启用性能域
	// _, err := CDPPerformanceEnable()
	// if err != nil {
	// 	log.Fatalf("启用性能域失败: %v", err)
	// }

	// 2. 执行页面交互操作（点击、滚动、加载数据等）
	// ExecutePageInteraction()

	// 3. 获取实时性能指标
	resp, err := CDPPerformanceGetMetrics()
	if err != nil {
		log.Printf("获取性能指标失败: %v, 响应: %s", err, resp)
		return
	}

	log.Println("成功获取页面实时性能指标：")
	log.Println(resp)
	// 包含常用指标：Timestamp、TickCount、JSHeapUsedSize、FramesPerSecond、ProcessCPUUsage等
}
*/
