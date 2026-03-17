package browser

import (
	"ChromeBot/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ScreenshotResult 截图结果结构体
type ScreenshotResult struct {
	Success bool `json:"success"` // 是否成功
	//Base64Str string `json:"base64Str"` // 截图base64字符串
	Width    int    `json:"width"`    // 截图宽度
	Height   int    `json:"height"`   // 截图高度
	FilePath string `json:"filePath"` // 保存的文件路径（可选）
	Error    string `json:"error"`    // 错误信息
}

func CaptureFullPageScreenshot(outputPath string) (*ScreenshotResult, error) {
	result := &ScreenshotResult{Success: false}

	if !DefaultNowTab(true) {
		return result, nil
	}

	// 步骤1：先获取页面完整尺寸（宽高）
	pageSize, err := getPageFullSize()
	if err != nil {
		result.Error = fmt.Sprintf("获取页面尺寸失败: %v", err)
		return result, fmt.Errorf("获取页面尺寸失败: %w", err)
	}
	result.Width = gt.Any2Int(pageSize["width"])
	result.Height = gt.Any2Int(pageSize["height"])
	utils.Debugf("页面完整尺寸：宽=%dpx，高=%dpx", result.Width, result.Height)

	// 步骤2：构造 Page.captureScreenshot 请求（关键参数）
	chromeInstance.NextID++
	screenshotMsg := map[string]interface{}{
		"id":     chromeInstance.NextID,
		"method": "Page.captureScreenshot",
		"params": map[string]interface{}{
			"format":                "png", // 截图格式（png/jpeg/webp）
			"quality":               100,   // 质量（仅jpeg/webp有效）
			"captureBeyondViewport": true,  // 核心：捕获视口外内容（全屏）
			"fromSurface":           true,  // 从渲染表面捕获（更清晰）
			"clip": map[string]interface{}{ // 截图区域（全屏）
				"x":      0,
				"y":      0,
				"width":  result.Width,
				"height": result.Height,
				"scale":  1.0, // 缩放比例
			},
		},
		"sessionId": chromeInstance.NowTabSession,
	}

	// 步骤3：发送截图请求
	err = chromeInstance.NowTabWSConn.WriteJSON(screenshotMsg)
	if err != nil {
		result.Error = fmt.Sprintf("发送截图请求失败: %v", err)
		return result, fmt.Errorf("发送截图请求失败: %w", err)
	}
	msgStr, _ := json.Marshal(screenshotMsg)
	utils.Debugf("发送截图请求: %s", string(msgStr))

	// 步骤4：等待并解析截图响应
	timeout := 10 * time.Second // 截图超时设为10秒（全屏截图可能耗时）
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				result.Error = "消息队列已关闭"
				return result, fmt.Errorf("消息队列已关闭")
			}

			//log.Println("收到截图响应 -> ", respMsg.Content)

			if chromeInstance.NextID == respMsg.ID {
				// 检查响应错误
				errObj, err := gt.JsonFind(respMsg.Content, "/error")
				if err == nil && errObj != nil {
					errMsg, _ := gt.JsonFind(respMsg.Content, "/error/message")
					result.Error = fmt.Sprintf("截图执行失败: %s", gt.Any2String(errMsg))
					return result, fmt.Errorf("截图执行失败: %s", gt.Any2String(errMsg))
				}

				// 解析base64截图数据
				base64Data, err := gt.JsonFind(respMsg.Content, "/result/data")
				if err != nil {
					result.Error = fmt.Sprintf("解析截图base64失败: %v", err)
					return result, fmt.Errorf("解析截图base64失败: %w", err)
				}
				//result.Base64Str = gt.Any2String(base64Data)
				result.Success = true

				// 步骤5：如果指定了保存路径，将base64保存为图片文件
				if outputPath != "" {

					outputPath = utils.SanitizeFileName(outputPath)

					err = saveBase64ToImage(gt.Any2String(base64Data), outputPath)
					if err != nil {
						result.Error = fmt.Sprintf("保存截图文件失败: %v", err)
						return result, fmt.Errorf("保存截图文件失败: %w", err)
					}
					result.FilePath = outputPath
					fmt.Printf("[Chrome]截图已保存到: %s \n", outputPath)
				}

				return result, nil
			} else {
				log.Println("忽略非当前截图请求的响应")
			}

		case <-timer.C:
			result.Error = "截图请求超时（10秒）"
			return result, fmt.Errorf("截图请求超时（10秒）")
		}
	}
}

// getPageFullSize 获取页面完整尺寸（宽高）
func getPageFullSize() (map[string]interface{}, error) {
	// 截图前获取尺寸等待500ms
	time.Sleep(500 * time.Millisecond)
	chromeInstance.NextID++
	// 执行JS获取页面完整宽高
	sizeMsg := map[string]interface{}{
		"id":     chromeInstance.NextID,
		"method": "Runtime.evaluate",
		"params": map[string]interface{}{
			"expression": `({
				width: Math.max(
					document.documentElement.scrollWidth,
					document.body.scrollWidth,
					document.documentElement.clientWidth,
					document.body.clientWidth
				),
				height: Math.max(
					document.documentElement.scrollHeight,
					document.body.scrollHeight,
					document.documentElement.clientHeight,
					document.body.clientHeight
				)
			})`,
			"returnByValue": true,
		},
		"sessionId": chromeInstance.NowTabSession,
	}

	msgStr, _ := json.Marshal(sizeMsg)
	utils.Debugf("发送: %s", string(msgStr))

	err := chromeInstance.NowTabWSConn.WriteJSON(sizeMsg)
	if err != nil {
		return nil, fmt.Errorf("发送获取尺寸请求失败: %w", err)
	}

	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return nil, fmt.Errorf("消息队列关闭")
			}
			if chromeInstance.NextID == respMsg.ID {

				utils.Debug("收到尺寸回复 : ", respMsg.Content)

				result, err := gt.JsonFind(respMsg.Content, "/result/result/value")
				if err != nil {
					return nil, fmt.Errorf("解析页面尺寸失败: %w", err)
				}
				return result.(map[string]interface{}), nil
			}

		case <-timer.C:
			return nil, fmt.Errorf("获取页面尺寸超时")
		}
	}
}

func saveBase64ToImage(base64Str, outputPath string) error {
	// 解码base64
	imgData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("base64解码失败: %w", err)
	}

	// 解析路径（处理相对路径→绝对路径，创建父目录）
	absPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("解析路径失败：%w", err)
	}

	// 获取父目录（如 "/tmp/data/test.txt" → "/tmp/data"）
	dir := filepath.Dir(absPath)
	// 创建父目录（不存在则创建，递归创建多级目录）
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建父目录失败：%w", err)
	}

	// 写入文件（覆盖写入，不存在则创建）
	if err := os.WriteFile(absPath, imgData, 0666); err != nil {
		return fmt.Errorf("写入文件失败：%w", err)
	}

	return nil
}
