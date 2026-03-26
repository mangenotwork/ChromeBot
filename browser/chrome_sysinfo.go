package browser

import (
	"ChromeBot/utils"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	gt "github.com/mangenotwork/gathertool"
)

// https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo/

// CDPSystemInfoGetFeatureState 获取feature的信息
// featureName :
//
//	gpu_acceleration   GPU 加速
//	vulkan   Vulkan 渲染
//	direct3d11  D3D11
//	canvas_oop_rasterization 画布离屏渲染
//	video_acceleration   视频硬件加速
//	webgl    WebGL
//	webgl2   WebGL2
//	webgpu    WebGPU
func CDPSystemInfoGetFeatureState(featureName string) (bool, error) {
	if !DefaultBrowserWS() {
		return false, nil
	}

	// 自增ID（和你现有代码一致）
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 CDP 请求：SystemInfo.getFeatureState
	message := fmt.Sprintf(`{
		   "id": %d,
		   "method": "SystemInfo.getFeatureState",
		   "params": {
			   "featureState": "%s"
		   }
		}`, reqID, featureName)

	// 发送 WebSocket 消息（完全沿用你的写法）
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 SystemInfo.getFeatureState 失败:", err)
		chromeInstance.BrowserWSConn = nil
	}

	utils.Debugf("发送 CDP 消息: %s", message)

	// 超时机制（和你现有代码一致）
	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 等待回复
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return false, fmt.Errorf("消息队列已关闭")
			}

			// 匹配自己的请求ID
			if reqID == respMsg.ID {
				fmt.Println("[CDP SystemInfo.getFeatureState] 收到回复 -> ", respMsg.Content)

				// 解析JSON
				result, err := gt.Json2Map(respMsg.Content)
				if err != nil {
					return false, fmt.Errorf("解析回复失败: %v", err)
				}

				// 提取 result -> enabled
				if resultData, ok := result["result"].(map[string]any); ok {
					if enabled, ok := resultData["enabled"].(bool); ok {
						fmt.Printf("[CDP SystemInfo.getFeatureState]功能 [%s] 状态: %t", featureName, enabled)
						return enabled, nil
					}
				}

				return false, fmt.Errorf("未获取到功能状态")
			}

		case <-timer.C:
			return false, fmt.Errorf("getFeatureState 请求超时")
		}
	}
}

// CDPSystemInfoGetInfo 获取Chrome完整系统信息（GPU、CPU、型号、驱动、版本等）
func CDPSystemInfoGetInfo() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}

	// 自增ID（和你现有代码完全一致）
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建 CDP 请求：SystemInfo.getInfo（该方法 无参数！）
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "SystemInfo.getInfo"
	}`, reqID)

	// 发送 WebSocket 消息（完全沿用你的写法）
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 SystemInfo.getInfo 失败:", err)
		chromeInstance.BrowserWSConn = nil
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)

	// 超时机制（和你现有代码一致）
	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 等待回复
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配自己的请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				fmt.Println("[CDP SystemInfo.getInfo] 收到回复 -> ", content)

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("getInfo 请求超时")
		}
	}
}

// CDPSystemInfoGetProcessInfo 获取Chrome所有进程信息（PID/类型/内存/CPU等）
func CDPSystemInfoGetProcessInfo() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}

	// 关键：必须用浏览器连接
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 SystemInfo.getProcessInfo")
	}

	// 请求ID自增
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建CDP请求（该接口无参数）
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "SystemInfo.getProcessInfo"
	}`, reqID)

	// 发送（Browser WS）
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 SystemInfo.getProcessInfo 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 等待回复
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}

			// 匹配自己的请求ID
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				fmt.Println("[CDP SystemInfo.getProcessInfo] 收到回复 -> ", content)

				return content, nil
			}

		case <-timer.C:
			return "", fmt.Errorf("getInfo 请求超时")
		}
	}
}
