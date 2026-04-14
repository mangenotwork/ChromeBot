package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  WebAudio.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 音频调试终止: 调试Web Audio API问题时，禁用音频上下文停止所有音频处理
// 2. 性能优化: 页面无需音频时，禁用WebAudio释放音频相关资源
// 3. 自动化测试: 测试流程中关闭音频功能，避免音频干扰测试结果
// 4. 浏览器状态重置: 切换页面/测试场景时，重置WebAudio为禁用状态
// 5. 内存释放: 长时间运行页面后，禁用WebAudio回收音频内存
// 6. 静音控制: 快速实现页面全局静音，停止所有WebAudio生成的音频

// CDPWebAudioDisable 禁用WebAudio功能
func CDPWebAudioDisable() (string, error) {
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
		"method": "WebAudio.disable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAudio.disable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 设置响应超时时间
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

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应JSON
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP响应错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAudio.disable 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1：调试时禁用WebAudio
func ExampleDisableWebAudioForDebug() {
	resp, err := CDPWebAudioDisable()
	if err != nil {
		log.Fatalf("禁用WebAudio失败: %v", err)
	}
	log.Printf("禁用WebAudio成功，响应: %s", resp)
}

// 示例2：页面性能优化，释放音频资源
func ExampleDisableWebAudioForPerformance() {
	// 业务逻辑：页面进入后台，无需音频播放
	resp, err := CDPWebAudioDisable()
	if err != nil {
		log.Printf("禁用WebAudio异常: %v", err)
		return
	}
	log.Println("已禁用WebAudio，释放音频资源")
}

// 示例3：自动化测试中静音页面
func ExampleDisableWebAudioInTest() {
	// 测试前置：禁用所有音频避免干扰
	_, err := CDPWebAudioDisable()
	if err != nil {
		panic("测试初始化失败：无法禁用WebAudio")
	}
	// 执行后续测试逻辑
}

*/

// -----------------------------------------------  WebAudio.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 音频调试启用: 调试Web Audio API时，开启音频上下文监听
// 2. 页面音频恢复: 禁用WebAudio后，重新启用音频播放功能
// 3. 自动化测试: 测试音频相关功能前，启用WebAudio模块
// 4. 前台唤醒: 页面从后台切回前台，恢复WebAudio服务
// 5. 功能初始化: 页面加载完成后，初始化WebAudio监听
// 6. 音频事件监听: 启用后才能接收WebAudio相关事件回调

// CDPWebAudioEnable 启用WebAudio功能
func CDPWebAudioEnable() (string, error) {
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
		"method": "WebAudio.enable"
	}`, reqID)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAudio.enable 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 设置响应超时时间
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

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应JSON
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP响应错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAudio.enable 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1：页面初始化时启用WebAudio
func ExampleEnableWebAudioOnPageInit() {
	resp, err := CDPWebAudioEnable()
	if err != nil {
		log.Fatalf("启用WebAudio失败: %v", err)
	}
	log.Printf("启用WebAudio成功，响应: %s", resp)
}

// 示例2：禁用后恢复WebAudio功能
func ExampleEnableWebAudioAfterDisable() {
	// 先禁用
	CDPWebAudioDisable()
	// 业务需要时重新启用
	resp, err := CDPWebAudioEnable()
	if err != nil {
		log.Printf("恢复WebAudio失败: %v", err)
		return
	}
	log.Println("已恢复WebAudio功能")
}

// 示例3：自动化测试前置启用WebAudio
func ExampleEnableWebAudioInTest() {
	// 测试音频功能前必须启用
	_, err := CDPWebAudioEnable()
	if err != nil {
		panic("测试初始化失败：无法启用WebAudio")
	}
	// 执行音频相关测试逻辑
}

*/

// -----------------------------------------------  WebAudio.getRealtimeData  -----------------------------------------------
// === 应用场景 ===
// 1. 音频实时监控: 获取WebAudio上下文的实时数据，监控音频播放状态
// 2. 可视化调试: 为音频可视化（频谱、波形）提供实时数据
// 3. 性能分析: 分析WebAudio的实时性能指标，排查卡顿问题
// 4. 自动化测试: 验证音频是否正常播放，获取实时数据断言
// 5. 调试辅助: 调试Web Audio API时，查看实时音频参数
// 6. 音频质量检测: 实时检测音频输出的质量和数据指标

// CDPWebAudioGetRealtimeData 获取WebAudio实时数据
// 参数 contextId: WebAudio上下文ID
func CDPWebAudioGetRealtimeData(contextId string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求消息（带参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "WebAudio.getRealtimeData",
		"params": {
			"contextId": "%s"
		}
	}`, reqID, contextId)

	// 发送WebSocket消息
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 WebAudio.getRealtimeData 请求失败: %w", err)
	}

	log.Printf("[DEBUG] 发送 CDP 消息: %s", message)

	// 设置响应超时时间
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

			// 匹配请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				log.Printf("[DEBUG] 收到回复: %s", content)

				// 解析响应JSON
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(respMsg.Content), &response); err != nil {
					return content, fmt.Errorf("解析响应失败: %w", err)
				}

				// 检查CDP响应错误
				if errorObj, exists := response["error"]; exists {
					return content, fmt.Errorf("CDP错误: %v", errorObj)
				}

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("WebAudio.getRealtimeData 请求超时")
		}
	}
}

/*


// === 使用示例代码 ===
// 示例1：获取指定音频上下文的实时数据
func ExampleGetWebAudioRealtimeData() {
	// 音频上下文ID（从WebAudio.contextCreated事件获取）
	audioContextID := "your-audio-context-id"

	resp, err := CDPWebAudioGetRealtimeData(audioContextID)
	if err != nil {
		log.Fatalf("获取WebAudio实时数据失败: %v", err)
	}
	log.Printf("WebAudio实时数据: %s", resp)
}

// 示例2：音频可视化数据获取
func ExampleGetRealtimeDataForVisualization() {
	ctxID := "audio-ctx-123"
	// 循环获取实时数据用于可视化展示
	for i := 0; i < 10; i++ {
		data, err := CDPWebAudioGetRealtimeData(ctxID)
		if err != nil {
			log.Printf("获取实时数据失败: %v", err)
			continue
		}
		log.Printf("第%d次音频实时数据: %s", i+1, data)
		time.Sleep(500 * time.Millisecond)
	}
}

// 示例3：自动化测试验证音频播放
func ExampleTestAudioPlayWithRealtimeData() {
	ctxID := "test-context-id"
	// 启用WebAudio
	CDPWebAudioEnable()

	// 获取实时数据判断音频是否正常
	data, err := CDPWebAudioGetRealtimeData(ctxID)
	if err != nil {
		panic("音频测试失败：无法获取实时数据")
	}
	log.Println("音频播放正常，实时数据：", data)
}

*/
