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
	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	message := fmt.Sprintf(`{
		   "id": %d,
		   "method": "SystemInfo.getFeatureState",
		   "params": {
			   "featureState": "%s"
		   }
		}`, reqID, featureName)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 SystemInfo.getFeatureState 失败:", err)
		chromeInstance.BrowserWSConn = nil
	}

	utils.Debugf("发送 CDP 消息: %s", message)

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return false, fmt.Errorf("消息队列已关闭")
			}
			if reqID == respMsg.ID {
				fmt.Println("[CDP SystemInfo.getFeatureState] 收到回复 -> ", respMsg.Content)
				result, err := gt.Json2Map(respMsg.Content)
				if err != nil {
					return false, fmt.Errorf("解析回复失败: %v", err)
				}
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
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
	   "id": %d,
	   "method": "SystemInfo.getInfo"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 SystemInfo.getInfo 失败:", err)
		chromeInstance.BrowserWSConn = nil
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)

	timeout := 6 * time.Second
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
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 SystemInfo.getProcessInfo")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "SystemInfo.getProcessInfo"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 SystemInfo.getProcessInfo 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
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
				fmt.Println("[CDP SystemInfo.getProcessInfo] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("getInfo 请求超时")
		}
	}
}
