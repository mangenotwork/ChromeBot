package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  Input.cancelDragging  -----------------------------------------------
// === 应用场景 ===
// 1. 取消拖拽操作: 取消正在进行的拖拽行为
// 2. 拖拽恢复: 拖拽过程中恢复原始状态
// 3. 自动化测试清理: 自动化测试中清理拖拽状态
// 4. 异常处理: 处理拖拽过程中出现的异常情况
// 5. 交互中断: 模拟用户中断拖拽操作
// 6. 测试场景重置: 测试完成后重置拖拽相关状态

// CDPInputCancelDragging 取消拖拽操作
func CDPInputCancelDragging() (string, error) {
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
		"method": "Input.cancelDragging"
	}`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 cancelDragging 请求失败: %w", err)
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
			return "", fmt.Errorf("cancelDragging 请求超时")
		}
	}
}

/*
// 场景描述：在自动化测试中，用户拖拽文件到上传区域后，需要取消拖拽操作
func TestCancelFileDrag() {
	// 1. 模拟拖拽开始
	// CDPInputDispatchMouseEvent(...) 发送拖拽开始事件

	// 2. 模拟取消拖拽
	result, err := CDPInputCancelDragging()
	if err != nil {
		log.Printf("取消拖拽失败: %v", err)
	} else {
		log.Printf("拖拽已取消: %s", result)
	}

	// 3. 验证拖拽被取消
	// 页面应恢复到拖拽前的状态
}

// 场景描述：拖拽过程中发生异常，需要恢复页面状态
func HandleDragException() error {
	// 模拟拖拽过程中的异常情况
	// 例如：页面卡住、元素位置错误等

	// 取消拖拽以恢复状态
	_, err := CDPInputCancelDragging()
	if err != nil {
		return fmt.Errorf("无法取消拖拽: %w", err)
	}

	log.Println("拖拽异常已恢复，页面状态已重置")
	return nil
}

// 场景描述：测试套件结束后清理拖拽相关状态
func TestTearDown() {
	// 确保所有拖拽操作都被正确结束
	result, err := CDPInputCancelDragging()
	if err != nil {
		log.Printf("警告：拖拽清理失败: %v", err)
	} else {
		log.Printf("拖拽状态已清理: %s", result)
	}

	// 其他清理操作...
}


*/

// -----------------------------------------------  Input.dispatchKeyEvent  -----------------------------------------------
// === 应用场景 ===
// 1. 键盘输入模拟: 模拟用户键盘按键输入
// 2. 快捷键测试: 测试应用程序的快捷键功能
// 3. 表单自动化: 自动化填写表单字段
// 4. 游戏控制: 模拟游戏中的键盘控制
// 5. 辅助功能测试: 测试键盘导航和无障碍功能
// 6. 组合键测试: 测试Ctrl、Alt、Shift等组合键功能

// CDPInputDispatchKeyEvent 发送键盘事件
func CDPInputDispatchKeyEvent(params map[string]interface{}) (string, error) {
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
		"method": "Input.dispatchKeyEvent",
		"params": %s
	}`, reqID, utils.MapToJson(params))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 dispatchKeyEvent 请求失败: %w", err)
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
			return "", fmt.Errorf("dispatchKeyEvent 请求超时")
		}
	}
}

/*

// 场景描述：在搜索框中输入搜索关键词
func SimulateSearchInput() {
	// 模拟输入 "hello world"
	params := map[string]interface{}{
		"type": "keyDown",
		"text": "h",
	}
	CDPInputDispatchKeyEvent(params)

	params["text"] = "e"
	CDPInputDispatchKeyEvent(params)

	params["text"] = "l"
	CDPInputDispatchKeyEvent(params)

	params["text"] = "l"
	CDPInputDispatchKeyEvent(params)

	params["text"] = "o"
	CDPInputDispatchKeyEvent(params)

	params["text"] = " "
	CDPInputDispatchKeyEvent(params)

	params["text"] = "w"
	CDPInputDispatchKeyEvent(params)

	// 继续输入其他字符...

	log.Println("搜索关键词输入完成")
}

// 场景描述：模拟Ctrl+C复制操作
func SimulateCopyShortcut() {
	// 按下Ctrl键
	params := map[string]interface{}{
		"type":         "keyDown",
		"key":          "Control",
		"code":         "ControlLeft",
		"windowsVirtualKeyCode": 17,
		"nativeVirtualKeyCode":  17,
	}
	CDPInputDispatchKeyEvent(params)

	// 按下C键
	params = map[string]interface{}{
		"type":         "keyDown",
		"key":          "c",
		"code":         "KeyC",
		"windowsVirtualKeyCode": 67,
		"nativeVirtualKeyCode":  67,
	}
	CDPInputDispatchKeyEvent(params)

	// 释放C键
	params["type"] = "keyUp"
	CDPInputDispatchKeyEvent(params)

	// 释放Ctrl键
	params = map[string]interface{}{
		"type":         "keyUp",
		"key":          "Control",
		"code":         "ControlLeft",
		"windowsVirtualKeyCode": 17,
		"nativeVirtualKeyCode":  17,
	}
	CDPInputDispatchKeyEvent(params)

	log.Println("Ctrl+C复制操作已执行")
}

// 场景描述：在输入框输入内容后按回车提交
func SimulateEnterSubmit() {
	// 先输入内容...

	// 按下回车键提交
	params := map[string]interface{}{
		"type":         "keyDown",
		"key":          "Enter",
		"code":         "Enter",
		"windowsVirtualKeyCode": 13,
		"nativeVirtualKeyCode":  13,
	}
	result, err := CDPInputDispatchKeyEvent(params)
	if err != nil {
		log.Printf("回车键按下失败: %v", err)
	} else {
		log.Printf("回车键已按下: %s", result)
	}

	// 释放回车键
	params["type"] = "keyUp"
	CDPInputDispatchKeyEvent(params)
}

// 场景描述：在游戏或列表中使用方向键导航
func SimulateArrowKeyNavigation() {
	// 按下向下箭头
	params := map[string]interface{}{
		"type":         "keyDown",
		"key":          "ArrowDown",
		"code":         "ArrowDown",
		"windowsVirtualKeyCode": 40,
		"nativeVirtualKeyCode":  40,
	}
	CDPInputDispatchKeyEvent(params)

	// 释放向下箭头
	params["type"] = "keyUp"
	CDPInputDispatchKeyEvent(params)

	// 可以继续模拟其他方向键...

	log.Println("方向键导航完成")
}


*/

// -----------------------------------------------  Input.dispatchMouseEvent  -----------------------------------------------
// === 应用场景 ===
// 1. 鼠标点击模拟: 模拟用户的鼠标点击操作
// 2. 拖拽操作: 模拟鼠标拖拽元素
// 3. 右键菜单: 模拟右键点击触发上下文菜单
// 4. 悬浮效果: 模拟鼠标悬停在元素上触发hover效果
// 5. 绘图应用: 模拟鼠标在画布上的绘制操作
// 6. 游戏控制: 模拟游戏中的鼠标控制

// CDPInputDispatchMouseEvent 发送鼠标事件
func CDPInputDispatchMouseEvent(params map[string]interface{}) (string, error) {
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
		"method": "Input.dispatchMouseEvent",
		"params": %s
	}`, reqID, utils.MapToJson(params))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 dispatchMouseEvent 请求失败: %w", err)
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
			return "", fmt.Errorf("dispatchMouseEvent 请求超时")
		}
	}
}

/*

// 场景描述：点击网页上的提交按钮
func SimulateButtonClick() {
	// 移动到按钮位置
	params := map[string]interface{}{
		"type":       "mouseMoved",
		"x":          200,
		"y":          150,
		"button":     "none",
		"buttons":    0,
	}
	CDPInputDispatchMouseEvent(params)

	// 点击鼠标左键
	params = map[string]interface{}{
		"type":       "mousePressed",
		"x":          200,
		"y":          150,
		"button":     "left",
		"clickCount": 1,
		"buttons":    1,
	}
	CDPInputDispatchMouseEvent(params)

	// 释放鼠标左键
	params["type"] = "mouseReleased"
	CDPInputDispatchMouseEvent(params)

	log.Println("按钮点击模拟完成")
}

// 场景描述：在元素上右键点击打开上下文菜单
func SimulateRightClickMenu() {
	// 移动到目标元素
	params := map[string]interface{}{
		"type":       "mouseMoved",
		"x":          300,
		"y":          250,
		"button":     "none",
		"buttons":    0,
	}
	CDPInputDispatchMouseEvent(params)

	// 按下鼠标右键
	params = map[string]interface{}{
		"type":       "mousePressed",
		"x":          300,
		"y":          250,
		"button":     "right",
		"clickCount": 1,
		"buttons":    2,
	}
	result, err := CDPInputDispatchMouseEvent(params)
	if err != nil {
		log.Printf("右键按下失败: %v", err)
	} else {
		log.Printf("右键已按下: %s", result)
	}

	// 释放鼠标右键
	params["type"] = "mouseReleased"
	CDPInputDispatchMouseEvent(params)

	log.Println("右键菜单已触发")
}

// 场景描述：拖拽文件到上传区域
func SimulateDragAndDrop() {
	// 移动到拖拽起点
	params := map[string]interface{}{
		"type":       "mouseMoved",
		"x":          100,
		"y":          100,
		"button":     "none",
		"buttons":    0,
	}
	CDPInputDispatchMouseEvent(params)

	// 按下鼠标左键开始拖拽
	params = map[string]interface{}{
		"type":       "mousePressed",
		"x":          100,
		"y":          100,
		"button":     "left",
		"clickCount": 1,
		"buttons":    1,
	}
	CDPInputDispatchMouseEvent(params)

	// 移动到目标位置
	params = map[string]interface{}{
		"type":       "mouseMoved",
		"x":          300,
		"y":          300,
		"button":     "left",
		"buttons":    1,
	}
	CDPInputDispatchMouseEvent(params)

	// 释放鼠标左键完成拖拽
	params["type"] = "mouseReleased"
	CDPInputDispatchMouseEvent(params)

	log.Println("拖拽操作完成")
}

// 场景描述：悬停在元素上触发tooltip显示
func SimulateMouseHover() {
	// 移动到目标元素
	params := map[string]interface{}{
		"type":       "mouseMoved",
		"x":          150,
		"y":          200,
		"button":     "none",
		"buttons":    0,
	}
	result, err := CDPInputDispatchMouseEvent(params)
	if err != nil {
		log.Printf("鼠标移动失败: %v", err)
	} else {
		log.Printf("鼠标已移动到指定位置: %s", result)
	}

	// 等待一段时间让hover生效
	time.Sleep(500 * time.Millisecond)

	log.Println("鼠标悬停效果已触发")
}

// 场景描述：双击文件或图标
func SimulateDoubleClick() {
	// 移动到目标位置
	params := map[string]interface{}{
		"type":       "mouseMoved",
		"x":          250,
		"y":          180,
		"button":     "none",
		"buttons":    0,
	}
	CDPInputDispatchMouseEvent(params)

	// 第一次点击
	params = map[string]interface{}{
		"type":       "mousePressed",
		"x":          250,
		"y":          180,
		"button":     "left",
		"clickCount": 1,
		"buttons":    1,
	}
	CDPInputDispatchMouseEvent(params)

	params["type"] = "mouseReleased"
	CDPInputDispatchMouseEvent(params)

	// 短暂延迟
	time.Sleep(50 * time.Millisecond)

	// 第二次点击
	params = map[string]interface{}{
		"type":       "mousePressed",
		"x":          250,
		"y":          180,
		"button":     "left",
		"clickCount": 2,  // 注意clickCount设置为2表示双击
		"buttons":    1,
	}
	CDPInputDispatchMouseEvent(params)

	params["type"] = "mouseReleased"
	CDPInputDispatchMouseEvent(params)

	log.Println("双击操作模拟完成")
}

*/

// -----------------------------------------------  Input.dispatchTouchEvent  -----------------------------------------------
// === 应用场景 ===
// 1. 触摸屏测试: 模拟移动设备上的触摸操作
// 2. 手势操作: 模拟缩放、滑动等手势操作
// 3. 移动端自动化: 移动端应用的自动化测试
// 4. 多点触控: 测试多点触摸交互功能
// 5. 触屏游戏: 模拟触屏游戏中的控制操作
// 6. 响应式设计测试: 测试触屏设备上的UI响应

// CDPInputDispatchTouchEvent 发送触摸事件
func CDPInputDispatchTouchEvent(params map[string]interface{}) (string, error) {
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
		"method": "Input.dispatchTouchEvent",
		"params": %s
	}`, reqID, utils.MapToJson(params))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 dispatchTouchEvent 请求失败: %w", err)
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
			return "", fmt.Errorf("dispatchTouchEvent 请求超时")
		}
	}
}

/*

// 场景描述：在移动设备上点击按钮
func SimulateTouchTap() {
	// 触摸开始
	params := map[string]interface{}{
		"type": "touchStart",
		"touchPoints": []map[string]interface{}{
			{
				"x": 100,
				"y": 200,
				"radiusX": 5,
				"radiusY": 5,
			},
		},
	}
	result, err := CDPInputDispatchTouchEvent(params)
	if err != nil {
		log.Printf("触摸开始失败: %v", err)
	} else {
		log.Printf("触摸开始: %s", result)
	}

	// 短暂延迟模拟触摸时长
	time.Sleep(100 * time.Millisecond)

	// 触摸结束
	params["type"] = "touchEnd"
	params["touchPoints"] = []map[string]interface{}{}
	CDPInputDispatchTouchEvent(params)

	log.Println("触摸点击操作完成")
}

// 场景描述：在移动端页面上下滑动
func SimulateSwipeGesture() {
	// 开始触摸
	params := map[string]interface{}{
		"type": "touchStart",
		"touchPoints": []map[string]interface{}{
			{
				"x": 200,
				"y": 300,
			},
		},
	}
	CDPInputDispatchTouchEvent(params)

	// 移动触摸点（向下滑动）
	for i := 1; i <= 5; i++ {
		params["type"] = "touchMove"
		params["touchPoints"] = []map[string]interface{}{
			{
				"x": 200,
				"y": 300 + i*20,
			},
		}
		CDPInputDispatchTouchEvent(params)
		time.Sleep(50 * time.Millisecond)
	}

	// 结束触摸
	params["type"] = "touchEnd"
	params["touchPoints"] = []map[string]interface{}{}
	CDPInputDispatchTouchEvent(params)

	log.Println("滑动操作模拟完成")
}

// 场景描述：在图片或地图上进行双指缩放
func SimulatePinchZoom() {
	// 开始触摸 - 两个手指
	params := map[string]interface{}{
		"type": "touchStart",
		"touchPoints": []map[string]interface{}{
			{
				"x": 150,
				"y": 200,
			},
			{
				"x": 250,
				"y": 200,
			},
		},
	}
	CDPInputDispatchTouchEvent(params)
	time.Sleep(100 * time.Millisecond)

	// 缩放 - 两个手指向外移动
	params["type"] = "touchMove"
	params["touchPoints"] = []map[string]interface{}{
		{
			"x": 100,  // 向左移动
			"y": 200,
		},
		{
			"x": 300,  // 向右移动
			"y": 200,
		},
	}
	CDPInputDispatchTouchEvent(params)
	time.Sleep(200 * time.Millisecond)

	// 结束触摸
	params["type"] = "touchEnd"
	params["touchPoints"] = []map[string]interface{}{}
	CDPInputDispatchTouchEvent(params)

	log.Println("双指缩放手势模拟完成")
}

// 场景描述：长按元素触发上下文菜单
func SimulateLongPress() {
	// 开始触摸
	params := map[string]interface{}{
		"type": "touchStart",
		"touchPoints": []map[string]interface{}{
			{
				"x": 180,
				"y": 240,
			},
		},
	}
	result, err := CDPInputDispatchTouchEvent(params)
	if err != nil {
		log.Printf("长按开始失败: %v", err)
	} else {
		log.Printf("长按开始: %s", result)
	}

	// 保持触摸状态（模拟长按）
	time.Sleep(800 * time.Millisecond)

	// 结束触摸
	params["type"] = "touchEnd"
	params["touchPoints"] = []map[string]interface{}{}
	CDPInputDispatchTouchEvent(params)

	log.Println("长按操作完成")
}

// 场景描述：多手指拖拽元素
func SimulateMultiTouchDrag() {
	// 开始触摸 - 两个手指
	params := map[string]interface{}{
		"type": "touchStart",
		"touchPoints": []map[string]interface{}{
			{
				"x": 100,
				"y": 150,
			},
			{
				"x": 200,
				"y": 150,
			},
		},
	}
	CDPInputDispatchTouchEvent(params)

	// 同时移动两个手指
	for i := 1; i <= 3; i++ {
		params["type"] = "touchMove"
		params["touchPoints"] = []map[string]interface{}{
			{
				"x": 100 + i*20,
				"y": 150 + i*20,
			},
			{
				"x": 200 + i*20,
				"y": 150 + i*20,
			},
		}
		CDPInputDispatchTouchEvent(params)
		time.Sleep(100 * time.Millisecond)
	}

	// 结束触摸
	params["type"] = "touchEnd"
	params["touchPoints"] = []map[string]interface{}{}
	CDPInputDispatchTouchEvent(params)

	log.Println("多点触控拖拽模拟完成")
}


*/

// -----------------------------------------------  Input.setIgnoreInputEvents  -----------------------------------------------
// === 应用场景 ===
// 1. 输入拦截: 临时禁用所有输入事件
// 2. 测试防抖: 测试应用在输入被禁用时的行为
// 3. 录制模式: 录制用户操作时临时忽略输入
// 4. 调试分析: 在调试过程中防止误操作干扰
// 5. 演示模式: 自动演示时禁用用户输入
// 6. 安全场景: 敏感操作时临时锁定用户输入

// CDPInputSetIgnoreInputEvents 设置是否忽略输入事件
func CDPInputSetIgnoreInputEvents(ignore bool) (string, error) {
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
		"method": "Input.setIgnoreInputEvents",
		"params": {
			"ignore": %t
		}
	}`, reqID, ignore)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setIgnoreInputEvents 请求失败: %w", err)
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
			return "", fmt.Errorf("setIgnoreInputEvents 请求超时")
		}
	}
}

/*

// 场景描述：录制用户操作流程时，防止误操作干扰
func StartRecordingMode() {
	// 开始录制前，禁用所有输入事件
	result, err := CDPInputSetIgnoreInputEvents(true)
	if err != nil {
		log.Printf("无法禁用输入事件: %v", err)
		return
	}
	log.Printf("输入事件已禁用: %s", result)

	// 开始录制操作
	log.Println("开始录制操作...")

	// 录制结束后，重新启用输入事件
	result, err = CDPInputSetIgnoreInputEvents(false)
	if err != nil {
		log.Printf("无法启用输入事件: %v", err)
	} else {
		log.Printf("输入事件已启用: %s", result)
	}

	log.Println("录制模式结束")
}

// 场景描述：执行敏感操作时临时锁定用户输入
func PerformSensitiveOperation() {
	// 锁定输入，防止误操作
	_, err := CDPInputSetIgnoreInputEvents(true)
	if err != nil {
		log.Printf("警告：无法锁定输入: %v", err)
	} else {
		log.Println("输入已锁定，开始执行敏感操作")
	}

	// 执行敏感操作
	// 例如：删除重要数据、修改配置等

	// 操作完成后，解锁输入
	result, err := CDPInputSetIgnoreInputEvents(false)
	if err != nil {
		log.Printf("警告：无法解锁输入: %v", err)
	} else {
		log.Printf("输入已解锁: %s", result)
	}
}

// 场景描述：测试应用在输入被快速禁用/启用时的行为
func TestInputThrottling() {
	log.Println("开始防抖测试...")

	// 快速切换输入状态
	for i := 0; i < 5; i++ {
		// 禁用输入
		_, err := CDPInputSetIgnoreInputEvents(true)
		if err != nil {
			log.Printf("第%d次禁用失败: %v", i+1, err)
		}

		// 短暂延迟
		time.Sleep(100 * time.Millisecond)

		// 启用输入
		result, err := CDPInputSetIgnoreInputEvents(false)
		if err != nil {
			log.Printf("第%d次启用失败: %v", i+1, err)
		} else {
			log.Printf("第%d次状态切换: %s", i+1, result)
		}

		time.Sleep(100 * time.Millisecond)
	}

	log.Println("防抖测试完成")
}

// 场景描述：在自动演示过程中防止用户输入干扰
func StartDemoMode() {
	// 开始演示前禁用输入
	_, err := CDPInputSetIgnoreInputEvents(true)
	if err != nil {
		log.Printf("无法进入演示模式: %v", err)
		return
	}

	log.Println("演示模式已启动，用户输入被禁用")

	// 执行演示操作序列
	// 例如：自动导航、点击、输入等

	// 模拟演示过程
	for i := 1; i <= 3; i++ {
		log.Printf("执行演示步骤 %d...", i)
		time.Sleep(1 * time.Second)
	}

	// 演示结束，恢复输入
	result, err := CDPInputSetIgnoreInputEvents(false)
	if err != nil {
		log.Printf("无法退出演示模式: %v", err)
	} else {
		log.Printf("演示模式已结束，用户输入已恢复: %s", result)
	}
}

// 场景描述：调试时隔离用户输入，专注于程序自动操作
func DebugWithInputIsolation() {
	log.Println("开始调试模式，输入事件将被隔离")

	// 保存当前输入状态
	originalIgnoreState := false // 假设初始状态为不忽略
	// 注意：实际应该先获取当前状态，这里简化处理

	// 启用输入隔离
	result, err := CDPInputSetIgnoreInputEvents(true)
	if err != nil {
		log.Printf("无法启用输入隔离: %v", err)
		return
	}
	log.Printf("输入隔离已启用: %s", result)

	// 执行调试操作
	// 这里可以运行自动化测试或调试代码
	log.Println("执行调试操作...")
	time.Sleep(2 * time.Second)

	// 恢复原始状态
	result, err = CDPInputSetIgnoreInputEvents(originalIgnoreState)
	if err != nil {
		log.Printf("无法恢复输入状态: %v", err)
	} else {
		log.Printf("输入状态已恢复: %s", result)
	}

	log.Println("调试模式结束")
}

*/

// -----------------------------------------------  Input.dispatchDragEvent  -----------------------------------------------
// === 应用场景 ===
// 1. 文件拖拽上传: 模拟文件拖拽到上传区域
// 2. 元素拖拽排序: 模拟列表或网格中元素的拖拽排序
// 3. 拖放操作测试: 测试拖放功能是否正常工作
// 4. 富文本编辑: 模拟在富文本编辑器中的拖拽操作
// 5. 图表操作: 拖拽图表元素进行交互
// 6. 跨窗口拖拽: 测试窗口间的拖拽功能

// CDPInputDispatchDragEvent 发送拖拽事件
func CDPInputDispatchDragEvent(params map[string]interface{}) (string, error) {
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
		"method": "Input.dispatchDragEvent",
		"params": %s
	}`, reqID, utils.MapToJson(params))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 dispatchDragEvent 请求失败: %w", err)
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
			return "", fmt.Errorf("dispatchDragEvent 请求超时")
		}
	}
}

/*

// 场景描述：将文件拖拽到文件上传区域
func SimulateFileDragUpload() {
	// 定义拖拽数据
	dragData := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"mimeType": "text/plain",
				"data": "data:text/plain;base64,SSBhbSBhIGZpbGUgY29udGVudC4=", // 示例文件内容
			},
		},
		"dragOperationsMask": 1, // 1表示复制操作
	}

	// 拖拽开始
	params := map[string]interface{}{
		"type": "dragEnter",
		"x": 100,
		"y": 100,
		"data": dragData,
	}
	result, err := CDPInputDispatchDragEvent(params)
	if err != nil {
		log.Printf("拖拽进入失败: %v", err)
	} else {
		log.Printf("拖拽进入: %s", result)
	}

	// 拖拽到目标位置
	params = map[string]interface{}{
		"type": "dragOver",
		"x": 300,
		"y": 200,
		"data": dragData,
	}
	CDPInputDispatchDragEvent(params)

	// 释放拖拽
	params = map[string]interface{}{
		"type": "drop",
		"x": 300,
		"y": 200,
		"data": dragData,
	}
	CDPInputDispatchDragEvent(params)

	log.Println("文件拖拽上传模拟完成")
}

// 场景描述：在任务列表中拖拽任务项进行重新排序
func SimulateTaskDragSort() {
	// 定义拖拽的HTML元素数据
	dragData := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"mimeType": "text/html",
				"data": "data:text/html;charset=utf-8,%3Cdiv%20id%3D%22task-1%22%3E任务1%3C%2Fdiv%3E",
			},
		},
		"dragOperationsMask": 4, // 4表示移动操作
	}

	// 从原位置开始拖拽
	params := map[string]interface{}{
		"type": "dragEnter",
		"x": 150,
		"y": 100,
		"data": dragData,
	}
	result, err := CDPInputDispatchDragEvent(params)
	if err != nil {
		log.Printf("开始拖拽失败: %v", err)
		return
	}
	log.Printf("开始拖拽任务: %s", result)

	// 移动过程
	params["type"] = "dragOver"
	params["x"] = 150
	params["y"] = 150
	CDPInputDispatchDragEvent(params)

	// 移动到新位置
	params["x"] = 150
	params["y"] = 200
	CDPInputDispatchDragEvent(params)

	// 放置到新位置
	params["type"] = "drop"
	CDPInputDispatchDragEvent(params)

	// 拖拽结束
	params["type"] = "dragEnd"
	CDPInputDispatchDragEvent(params)

	log.Println("任务拖拽排序完成")
}

// 场景描述：在富文本编辑器中拖拽文本
func SimulateRichTextDrag() {
	// 文本拖拽数据
	dragData := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"mimeType": "text/plain",
				"data": "data:text/plain;base64,VGhpcyBpcyBkcmFnZ2VkIHRleHQ=", // "This is dragged text"
			},
		},
		"files": []string{}, // 没有文件
		"dragOperationsMask": 1,
	}

	// 开始拖拽文本
	params := map[string]interface{}{
		"type": "dragEnter",
		"x": 200,
		"y": 100,
		"data": dragData,
	}
	CDPInputDispatchDragEvent(params)

	// 拖拽到编辑器内部
	params["type"] = "dragOver"
	for y := 120; y <= 300; y += 20 {
		params["y"] = y
		CDPInputDispatchDragEvent(params)
		time.Sleep(20 * time.Millisecond)
	}

	// 放置文本
	params["type"] = "drop"
	params["x"] = 300
	params["y"] = 300
	result, err := CDPInputDispatchDragEvent(params)
	if err != nil {
		log.Printf("文本放置失败: %v", err)
	} else {
		log.Printf("文本已放置: %s", result)
	}

	log.Println("富文本拖拽完成")
}

// 场景描述：在图表应用中拖拽数据点
func SimulateChartElementDrag() {
	// 模拟拖拽图表中的数据点
	dragData := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"mimeType": "application/json",
				"data": "data:application/json;base64,eyJ2YWx1ZSI6IDEwMCwgImxhYmVsIjogIkRhdGEgUG9pbnQifQ==",
			},
		},
		"dragOperationsMask": 1,
	}

	// 开始拖拽数据点
	params := map[string]interface{}{
		"type": "dragEnter",
		"x": 250,
		"y": 180,
		"data": dragData,
	}
	CDPInputDispatchDragEvent(params)
	time.Sleep(100 * time.Millisecond)

	// 拖动到新位置
	params["type"] = "dragOver"
	params["x"] = 350
	params["y"] = 280
	CDPInputDispatchDragEvent(params)

	// 放置数据点
	params["type"] = "drop"
	result, err := CDPInputDispatchDragEvent(params)
	if err != nil {
		log.Printf("图表元素放置失败: %v", err)
	} else {
		log.Printf("图表元素已移动: %s", result)
	}

	log.Println("图表元素拖拽完成")
}

// 场景描述：开始拖拽但最终取消操作
func SimulateDragCancel() {
	dragData := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"mimeType": "text/plain",
				"data": "data:text/plain;base64,Q2FuY2VsZWQgZHJhZw==",
			},
		},
		"dragOperationsMask": 1,
	}

	// 开始拖拽
	params := map[string]interface{}{
		"type": "dragEnter",
		"x": 100,
		"y": 100,
		"data": dragData,
	}
	CDPInputDispatchDragEvent(params)

	// 拖动一段距离
	params["type"] = "dragOver"
	params["x"] = 200
	params["y"] = 150
	CDPInputDispatchDragEvent(params)

	// 取消拖拽（拖到无效区域）
	params["type"] = "dragEnd"
	params["x"] = 50
	params["y"] = 50
	result, err := CDPInputDispatchDragEvent(params)
	if err != nil {
		log.Printf("拖拽取消失败: %v", err)
	} else {
		log.Printf("拖拽已取消: %s", result)
	}

	log.Println("拖拽取消操作完成")
}




*/

// -----------------------------------------------  Input.emulateTouchFromMouseEvent  -----------------------------------------------
// === 应用场景 ===
// 1. 移动端兼容性测试: 在桌面浏览器中测试移动端触摸事件
// 2. 混合交互测试: 测试鼠标事件如何触发触摸事件
// 3. 响应式设计验证: 验证网站在触摸模拟下的响应
// 4. 游戏兼容性: 测试桌面鼠标控制对移动触摸的映射
// 5. 自动化测试转换: 将鼠标测试用例转换为触摸测试
// 6. 无障碍功能测试: 测试鼠标到触摸的转换是否正常工作

// CDPInputEmulateTouchFromMouseEvent 从鼠标事件模拟触摸事件
func CDPInputEmulateTouchFromMouseEvent(params map[string]interface{}) (string, error) {
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
		"method": "Input.emulateTouchFromMouseEvent",
		"params": %s
	}`, reqID, utils.MapToJson(params))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 emulateTouchFromMouseEvent 请求失败: %w", err)
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
			return "", fmt.Errorf("emulateTouchFromMouseEvent 请求超时")
		}
	}
}

/*

// 场景描述：在桌面浏览器中模拟移动端触摸点击
func SimulateTouchClickWithMouse() {
	// 模拟鼠标左键点击转换为触摸点击
	params := map[string]interface{}{
		"type": "mousePressed",
		"x": 200,
		"y": 150,
		"button": "left",
		"clickCount": 1,
		"modifiers": 0,
		"timestamp": float64(time.Now().UnixNano() / 1e6),
	}

	result, err := CDPInputEmulateTouchFromMouseEvent(params)
	if err != nil {
		log.Printf("触摸点击模拟失败: %v", err)
	} else {
		log.Printf("触摸点击已触发: %s", result)
	}

	// 模拟鼠标释放
	params["type"] = "mouseReleased"
	CDPInputEmulateTouchFromMouseEvent(params)

	log.Println("鼠标到触摸点击转换完成")
}

// 场景描述：用鼠标操作模拟触摸屏滑动
func SimulateTouchSwipeWithMouse() {
	startX := 100
	startY := 200
	endX := 300
	endY := 200

	// 模拟触摸开始（鼠标按下）
	params := map[string]interface{}{
		"type": "mousePressed",
		"x": startX,
		"y": startY,
		"button": "left",
		"clickCount": 1,
		"modifiers": 0,
		"timestamp": float64(time.Now().UnixNano() / 1e6),
	}
	CDPInputEmulateTouchFromMouseEvent(params)

	// 模拟滑动过程
	for i := 1; i <= 5; i++ {
		currentX := startX + (endX-startX)*i/5
		currentY := startY + (endY-startY)*i/5

		params["type"] = "mouseMoved"
		params["x"] = currentX
		params["y"] = currentY
		params["timestamp"] = float64(time.Now().UnixNano() / 1e6)

		CDPInputEmulateTouchFromMouseEvent(params)
		time.Sleep(50 * time.Millisecond)
	}

	// 模拟触摸结束（鼠标释放）
	params["type"] = "mouseReleased"
	params["x"] = endX
	params["y"] = endY
	params["timestamp"] = float64(time.Now().UnixNano() / 1e6)

	result, err := CDPInputEmulateTouchFromMouseEvent(params)
	if err != nil {
		log.Printf("滑动结束失败: %v", err)
	} else {
		log.Printf("触摸滑动完成: %s", result)
	}

	log.Println("鼠标到触摸滑动转换完成")
}

// 场景描述：用鼠标右键模拟移动端的长按操作
func SimulateTouchLongPressWithMouse() {
	// 鼠标右键按下模拟触摸开始
	params := map[string]interface{}{
		"type": "mousePressed",
		"x": 250,
		"y": 180,
		"button": "right",  // 使用右键模拟长按
		"clickCount": 1,
		"modifiers": 0,
		"timestamp": float64(time.Now().UnixNano() / 1e6),
	}

	result, err := CDPInputEmulateTouchFromMouseEvent(params)
	if err != nil {
		log.Printf("长按开始失败: %v", err)
	} else {
		log.Printf("长按开始: %s", result)
	}

	// 保持按下状态模拟长按
	time.Sleep(800 * time.Millisecond)

	// 鼠标右键释放模拟触摸结束
	params["type"] = "mouseReleased"
	params["timestamp"] = float64(time.Now().UnixNano() / 1e6)
	CDPInputEmulateTouchFromMouseEvent(params)

	log.Println("鼠标右键模拟长按完成")
}

// 场景描述：用鼠标拖拽模拟触摸屏的拖拽操作
func SimulateTouchDragWithMouse() {
	// 开始拖拽
	params := map[string]interface{}{
		"type": "mousePressed",
		"x": 150,
		"y": 120,
		"button": "left",
		"clickCount": 1,
		"modifiers": 0,
		"timestamp": float64(time.Now().UnixNano() / 1e6),
	}
	CDPInputEmulateTouchFromMouseEvent(params)

	// 拖拽移动
	for i := 1; i <= 4; i++ {
		params["type"] = "mouseMoved"
		params["x"] = 150 + i*50
		params["y"] = 120 + i*30
		params["timestamp"] = float64(time.Now().UnixNano() / 1e6)

		CDPInputEmulateTouchFromMouseEvent(params)
		time.Sleep(100 * time.Millisecond)
	}

	// 拖拽结束
	params["type"] = "mouseReleased"
	params["timestamp"] = float64(time.Now().UnixNano() / 1e6)
	result, err := CDPInputEmulateTouchFromMouseEvent(params)
	if err != nil {
		log.Printf("拖拽结束失败: %v", err)
	} else {
		log.Printf("触摸拖拽完成: %s", result)
	}

	log.Println("鼠标拖拽模拟触摸拖拽完成")
}

// 场景描述：测试网站在触摸模拟下的响应性
func TestResponsiveDesignWithTouchSimulation() {
	log.Println("开始响应式触摸模拟测试...")

	// 测试不同位置的触摸点击
	testPoints := []struct {
		x int
		y int
		desc string
	}{
		{100, 100, "左上角"},
		{300, 200, "中间区域"},
		{600, 400, "右下角"},
	}

	for _, point := range testPoints {
		log.Printf("测试%s区域触摸响应...", point.desc)

		// 鼠标按下模拟触摸开始
		params := map[string]interface{}{
			"type": "mousePressed",
			"x": point.x,
			"y": point.y,
			"button": "left",
			"clickCount": 1,
			"modifiers": 0,
			"timestamp": float64(time.Now().UnixNano() / 1e6),
		}

		result, err := CDPInputEmulateTouchFromMouseEvent(params)
		if err != nil {
			log.Printf("%s区域触摸失败: %v", point.desc, err)
		} else {
			log.Printf("%s区域触摸成功: %s", point.desc, result)
		}

		// 短暂延迟
		time.Sleep(200 * time.Millisecond)

		// 鼠标释放模拟触摸结束
		params["type"] = "mouseReleased"
		params["timestamp"] = float64(time.Now().UnixNano() / 1e6)
		CDPInputEmulateTouchFromMouseEvent(params)

		time.Sleep(300 * time.Millisecond)
	}

	log.Println("响应式触摸模拟测试完成")
}

// 场景描述：测试桌面鼠标控制转换为移动触摸的兼容性
func TestGameControlTouchSimulation() {
	log.Println("开始游戏控制触摸模拟测试...")

	// 模拟游戏中的触摸控制
	controlActions := []string{
		"点击攻击按钮",
		"滑动移动角色",
		"长按技能蓄力",
		"拖拽道具",
	}

	for i, action := range controlActions {
		log.Printf("测试动作: %s", action)

		switch i {
		case 0: // 点击攻击
			params := map[string]interface{}{
				"type": "mousePressed",
				"x": 400,
				"y": 500,
				"button": "left",
				"clickCount": 1,
				"timestamp": float64(time.Now().UnixNano() / 1e6),
			}
			CDPInputEmulateTouchFromMouseEvent(params)
			time.Sleep(50 * time.Millisecond)
			params["type"] = "mouseReleased"
			CDPInputEmulateTouchFromMouseEvent(params)

		case 1: // 滑动移动
			// 模拟滑动开始
			params := map[string]interface{}{
				"type": "mousePressed",
				"x": 200,
				"y": 300,
				"button": "left",
				"timestamp": float64(time.Now().UnixNano() / 1e6),
			}
			CDPInputEmulateTouchFromMouseEvent(params)

			// 模拟滑动过程
			for j := 1; j <= 3; j++ {
				params["type"] = "mouseMoved"
				params["x"] = 200 + j*50
				params["y"] = 300
				params["timestamp"] = float64(time.Now().UnixNano() / 1e6)
				CDPInputEmulateTouchFromMouseEvent(params)
				time.Sleep(30 * time.Millisecond)
			}

			// 滑动结束
			params["type"] = "mouseReleased"
			CDPInputEmulateTouchFromMouseEvent(params)

		case 2: // 长按蓄力
			params := map[string]interface{}{
				"type": "mousePressed",
				"x": 300,
				"y": 400,
				"button": "left",
				"timestamp": float64(time.Now().UnixNano() / 1e6),
			}
			CDPInputEmulateTouchFromMouseEvent(params)
			time.Sleep(1000 * time.Millisecond) // 长按1秒
			params["type"] = "mouseReleased"
			CDPInputEmulateTouchFromMouseEvent(params)

		case 3: // 拖拽道具
			params := map[string]interface{}{
				"type": "mousePressed",
				"x": 350,
				"y": 250,
				"button": "left",
				"timestamp": float64(time.Now().UnixNano() / 1e6),
			}
			CDPInputEmulateTouchFromMouseEvent(params)

			params["type"] = "mouseMoved"
			params["x"] = 450
			params["y"] = 350
			params["timestamp"] = float64(time.Now().UnixNano() / 1e6)
			CDPInputEmulateTouchFromMouseEvent(params)

			params["type"] = "mouseReleased"
			CDPInputEmulateTouchFromMouseEvent(params)
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Println("游戏控制触摸模拟测试完成")
}


*/

// -----------------------------------------------  Input.imeSetComposition  -----------------------------------------------
// === 应用场景 ===
// 1. 输入法模拟: 模拟输入法的组合输入过程
// 2. 国际化测试: 测试非拉丁字符的输入
// 3. 中文输入测试: 测试拼音、五笔等中文输入法
// 4. 日语输入测试: 测试假名、汉字转换
// 5. 韩文输入测试: 测试韩文组合输入
// 6. 输入法UI测试: 测试输入法候选框、组合框显示

// CDPInputImeSetComposition 设置IME（输入法）组合
func CDPInputImeSetComposition(params map[string]interface{}) (string, error) {
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
		"method": "Input.imeSetComposition",
		"params": %s
	}`, reqID, utils.MapToJson(params))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 imeSetComposition 请求失败: %w", err)
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
			return "", fmt.Errorf("imeSetComposition 请求超时")
		}
	}
}

/*

// 场景描述：模拟中文拼音输入法的输入过程
func SimulateChinesePinyinInput() {
	log.Println("开始模拟中文拼音输入...")

	// 步骤1: 输入拼音"nihao"
	params := map[string]interface{}{
		"text": "nihao",
		"selectionStart": 6,  // 光标在末尾
		"selectionEnd": 6,
		"replacementStart": 0,
		"replacementEnd": 0,
	}

	result, err := CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("拼音输入失败: %v", err)
	} else {
		log.Printf("拼音输入成功: %s", result)
	}

	// 步骤2: 选择候选词（这里模拟选择"你好"）
	time.Sleep(500 * time.Millisecond)
	params = map[string]interface{}{
		"text": "你好",
		"selectionStart": 2,  // 光标在"你好"后面
		"selectionEnd": 2,
		"replacementStart": 0,
		"replacementEnd": 6,  // 替换整个拼音
	}

	result, err = CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("候选词选择失败: %v", err)
	} else {
		log.Printf("候选词选择成功: %s", result)
	}

	log.Println("中文拼音输入模拟完成")
}

// 场景描述：模拟日语假名输入和汉字转换
func SimulateJapaneseKanaInput() {
	log.Println("开始模拟日语输入...")

	// 输入假名"こんにちは"
	params := map[string]interface{}{
		"text": "こんにちは",
		"selectionStart": 5,  // 光标在末尾
		"selectionEnd": 5,
		"replacementStart": 0,
		"replacementEnd": 0,
	}

	result, err := CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("假名输入失败: %v", err)
	} else {
		log.Printf("假名输入成功: %s", result)
	}

	// 转换为汉字"今日は"
	time.Sleep(500 * time.Millisecond)
	params = map[string]interface{}{
		"text": "今日は",
		"selectionStart": 3,  // 光标在"今日は"后面
		"selectionEnd": 3,
		"replacementStart": 0,
		"replacementEnd": 5,  // 替换整个假名
	}

	result, err = CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("汉字转换失败: %v", err)
	} else {
		log.Printf("汉字转换成功: %s", result)
	}

	log.Println("日语输入模拟完成")
}

// 场景描述：模拟韩文的组合字符输入
func SimulateKoreanInput() {
	log.Println("开始模拟韩文输入...")

	// 输入韩文字母组合
	compositionText := "안녕하세요"  // 你好

	params := map[string]interface{}{
		"text": compositionText,
		"selectionStart": len(compositionText),  // 光标在末尾
		"selectionEnd": len(compositionText),
		"replacementStart": 0,
		"replacementEnd": 0,
	}

	result, err := CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("韩文输入失败: %v", err)
	} else {
		log.Printf("韩文输入成功: %s", result)
	}

	// 确认输入
	time.Sleep(300 * time.Millisecond)
	params = map[string]interface{}{
		"text": compositionText,
		"selectionStart": len(compositionText),
		"selectionEnd": len(compositionText),
		"replacementStart": 0,
		"replacementEnd": len(compositionText),  // 替换整个组合
	}

	result, err = CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("韩文确认失败: %v", err)
	} else {
		log.Printf("韩文确认成功: %s", result)
	}

	log.Println("韩文输入模拟完成")
}

// 场景描述：模拟输入法候选词的选择过程
func SimulateInputMethodCandidateSelection() {
	log.Println("开始模拟输入法候选选择...")

	// 1. 输入拼音"zhong"
	params := map[string]interface{}{
		"text": "zhong",
		"selectionStart": 5,
		"selectionEnd": 5,
		"replacementStart": 0,
		"replacementEnd": 0,
	}
	CDPInputImeSetComposition(params)
	log.Println("已输入拼音: zhong")

	time.Sleep(300 * time.Millisecond)

	// 2. 输入拼音"guo"
	params["text"] = "zhongguo"
	params["selectionStart"] = 8
	params["selectionEnd"] = 8
	CDPInputImeSetComposition(params)
	log.Println("已输入拼音: zhongguo")

	time.Sleep(300 * time.Millisecond)

	// 3. 显示候选词"中国"
	params["text"] = "中国"
	params["selectionStart"] = 2
	params["selectionEnd"] = 2
	params["replacementStart"] = 0
	params["replacementEnd"] = 8
	result, err := CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("候选词显示失败: %v", err)
	} else {
		log.Printf("候选词显示成功: %s", result)
	}

	// 4. 最终确认输入
	time.Sleep(200 * time.Millisecond)
	params = map[string]interface{}{
		"text": "中国",
		"selectionStart": 2,
		"selectionEnd": 2,
		"replacementStart": 0,
		"replacementEnd": 0,  // 不替换，直接确认
	}

	result, err = CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("最终确认失败: %v", err)
	} else {
		log.Printf("最终确认成功: %s", result)
	}

	log.Println("输入法候选选择模拟完成")
}

// 场景描述：模拟输入法组合过程中的编辑操作
func SimulateInputMethodCompositionEditing() {
	log.Println("开始模拟输入法组合编辑...")

	// 初始输入
	params := map[string]interface{}{
		"text": "beijin",  // 故意拼错
		"selectionStart": 6,
		"selectionEnd": 6,
		"replacementStart": 0,
		"replacementEnd": 0,
	}
	CDPInputImeSetComposition(params)
	log.Println("初始输入: beijin")

	time.Sleep(300 * time.Millisecond)

	// 编辑：添加'g'变成正确的"beijing"
	params["text"] = "beijing"
	params["selectionStart"] = 7
	params["selectionEnd"] = 7
	result, err := CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("编辑失败: %v", err)
	} else {
		log.Printf("编辑成功: beijing, 结果: %s", result)
	}

	time.Sleep(300 * time.Millisecond)

	// 显示候选词"北京"
	params["text"] = "北京"
	params["selectionStart"] = 2
	params["selectionEnd"] = 2
	params["replacementStart"] = 0
	params["replacementEnd"] = 7
	CDPInputImeSetComposition(params)
	log.Println("显示候选词: 北京")

	time.Sleep(200 * time.Millisecond)

	// 最终确认
	params = map[string]interface{}{
		"text": "北京",
		"selectionStart": 2,
		"selectionEnd": 2,
		"replacementStart": 0,
		"replacementEnd": 0,
	}

	result, err = CDPInputImeSetComposition(params)
	if err != nil {
		log.Printf("确认失败: %v", err)
	} else {
		log.Printf("确认成功: %s", result)
	}

	log.Println("输入法组合编辑模拟完成")
}

// 场景描述：测试输入法候选框、组合框的UI显示
func TestInputMethodUIDisplay() {
	log.Println("开始测试输入法UI显示...")

	testCases := []struct {
		name     string
		text     string
		cursorPos int
		desc     string
	}{
		{"短词", "hello", 5, "测试短英文单词"},
		{"长词", "internationalization", 20, "测试长英文单词"},
		{"中文", "测试文本", 4, "测试中文字符"},
		{"混合", "hello世界", 7, "测试中英文混合"},
		{"符号", "hello@world.com", 15, "测试带符号文本"},
	}

	for _, tc := range testCases {
		log.Printf("测试用例: %s - %s", tc.name, tc.desc)

		params := map[string]interface{}{
			"text": tc.text,
			"selectionStart": tc.cursorPos,
			"selectionEnd": tc.cursorPos,
			"replacementStart": 0,
			"replacementEnd": 0,
		}

		result, err := CDPInputImeSetComposition(params)
		if err != nil {
			log.Printf("  %s 测试失败: %v", tc.name, err)
		} else {
			log.Printf("  %s 测试成功: %s", tc.name, result)

			// 验证输入法UI是否正确显示
			// 这里可以添加截图或其他验证逻辑
		}

		time.Sleep(200 * time.Millisecond)

		// 清除组合
		params = map[string]interface{}{
			"text": "",
			"selectionStart": 0,
			"selectionEnd": 0,
			"replacementStart": 0,
			"replacementEnd": len(tc.text),
		}
		CDPInputImeSetComposition(params)

		time.Sleep(100 * time.Millisecond)
	}

	log.Println("输入法UI显示测试完成")
}




*/

// -----------------------------------------------  Input.insertText  -----------------------------------------------
// === 应用场景 ===
// 1. 文本插入: 在光标位置直接插入文本
// 2. 批量文本输入: 快速插入大量文本内容
// 3. 富文本编辑: 在富文本编辑器中插入格式化文本
// 4. 表单填充: 自动化填写表单字段
// 5. 代码编辑器: 在代码编辑器中插入代码片段
// 6. 内容替换: 替换选中的文本内容

// CDPInputInsertText 插入文本
func CDPInputInsertText(params map[string]interface{}) (string, error) {
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
		"method": "Input.insertText",
		"params": %s
	}`, reqID, utils.MapToJson(params))

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 insertText 请求失败: %w", err)
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
			return "", fmt.Errorf("insertText 请求超时")
		}
	}
}

/*

// 场景描述：在网页输入框中直接插入文本内容
func InsertTextIntoInputField() {
	// 在输入框中插入文本
	params := map[string]interface{}{
		"text": "这是一段插入的文本内容",
	}

	result, err := CDPInputInsertText(params)
	if err != nil {
		log.Printf("文本插入失败: %v", err)
	} else {
		log.Printf("文本插入成功: %s", result)
	}

	log.Println("输入框文本插入完成")
}

// 场景描述：自动化填写注册表单
func AutoFillRegistrationForm() {
	log.Println("开始自动填写注册表单...")

	// 填充用户名
	params := map[string]interface{}{
		"text": "testuser2023",
	}
	CDPInputInsertText(params)
	log.Println("已填写用户名")

	// 切换到下一个字段（这里需要先定位到下一个输入框）
	// 可以使用其他CDP方法切换焦点
	time.Sleep(200 * time.Millisecond)

	// 填充邮箱
	params["text"] = "test@example.com"
	CDPInputInsertText(params)
	log.Println("已填写邮箱")

	// 切换到密码字段
	time.Sleep(200 * time.Millisecond)

	// 填充密码
	params["text"] = "SecurePass123!"
	CDPInputInsertText(params)
	log.Println("已填写密码")

	// 切换到确认密码字段
	time.Sleep(200 * time.Millisecond)

	// 确认密码
	params["text"] = "SecurePass123!"
	CDPInputInsertText(params)
	log.Println("已确认密码")

	log.Println("注册表单自动填写完成")
}

// 场景描述：在富文本编辑器（如TinyMCE、CKEditor）中插入内容
func InsertFormattedContentIntoRichEditor() {
	log.Println("开始在富文本编辑器中插入格式化内容...")

	// 插入HTML格式的文本
	htmlContent := `<h1>标题</h1>
<p>这是一段<strong>加粗</strong>的文本。</p>
<ul>
	<li>列表项1</li>
	<li>列表项2</li>
</ul>`

	params := map[string]interface{}{
		"text": htmlContent,
	}

	result, err := CDPInputInsertText(params)
	if err != nil {
		log.Printf("富文本插入失败: %v", err)
	} else {
		log.Printf("富文本插入成功: %s", result)
	}

	log.Println("富文本编辑器内容插入完成")
}

// 场景描述：在在线代码编辑器（如CodePen、JSFiddle）中插入代码
func InsertCodeSnippetIntoCodeEditor() {
	log.Println("开始在代码编辑器中插入代码片段...")

	// 插入JavaScript代码
	javascriptCode := `function greet(name) {
	console.log('Hello, ' + name + '!');
}

// 调用函数
greet('World');`

	params := map[string]interface{}{
		"text": javascriptCode,
	}

	result, err := CDPInputInsertText(params)
	if err != nil {
		log.Printf("代码插入失败: %v", err)
	} else {
		log.Printf("代码插入成功: %s", result)

		// 验证代码是否正确插入
		// 可以添加代码高亮或语法检查验证
	}

	log.Println("代码编辑器内容插入完成")
}

// 场景描述：替换文档中选中的文本内容
func ReplaceSelectedText() {
	log.Println("开始替换选中文本...")

	// 注意：这个操作需要先有选中的文本
	// 可以使用其他CDP方法先选中文本

	// 替换选中的文本
	replacementText := "新的文本内容"
	params := map[string]interface{}{
		"text": replacementText,
	}

	result, err := CDPInputInsertText(params)
	if err != nil {
		log.Printf("文本替换失败: %v", err)
	} else {
		log.Printf("文本替换成功: %s", result)
	}

	log.Println("选中文本替换完成")
}

// 场景描述：在聊天应用或邮件客户端中插入常用短语
func InsertCommonPhrases() {
	log.Println("开始插入常用短语...")

	commonPhrases := []struct {
		name string
		text string
	}{
		{"问候语", "您好！\n祝您有愉快的一天！"},
		{"感谢语", "非常感谢您的帮助！"},
		{"询问语", "请问这个问题如何解决？"},
		{"确认语", "我已收到，会尽快处理。"},
		{"结束语", "祝好，\n[您的名字]"},
	}

	for _, phrase := range commonPhrases {
		log.Printf("插入短语: %s", phrase.name)

		params := map[string]interface{}{
			"text": phrase.text,
		}

		result, err := CDPInputInsertText(params)
		if err != nil {
			log.Printf("  %s 插入失败: %v", phrase.name, err)
		} else {
			log.Printf("  %s 插入成功: %s", phrase.name, result)
		}

		// 清空当前内容，为下一个插入做准备
		time.Sleep(500 * time.Millisecond)

		// 这里可以使用其他CDP方法清空输入框
		// 例如：全选然后删除
	}

	log.Println("常用短语插入完成")
}

// 场景描述：测试不同语言文本的插入功能
func TestMultilingualTextInsertion() {
	log.Println("开始多语言文本插入测试...")

	multilingualTexts := []struct {
		language string
		text     string
	}{
		{"English", "Hello, World!"},
		{"简体中文", "你好，世界！"},
		{"繁体中文", "你好，世界！"},
		{"日本語", "こんにちは、世界！"},
		{"한국어", "안녕하세요, 세계!"},
		{"Русский", "Привет, мир!"},
		{"Français", "Bonjour, le monde !"},
		{"Deutsch", "Hallo, Welt!"},
		{"Español", "¡Hola, mundo!"},
		{"العربية", "مرحبا بالعالم!"},
		{"Emoji", "👋🌍！测试表情符号 📱💻🔧"},
		{"混合文本", "Hello 世界! こんにちは! 안녕하세요!"},
	}

	for _, item := range multilingualTexts {
		log.Printf("测试语言: %s", item.language)

		params := map[string]interface{}{
			"text": item.text,
		}

		result, err := CDPInputInsertText(params)
		if err != nil {
			log.Printf("  %s 文本插入失败: %v", item.language, err)
		} else {
			log.Printf("  %s 文本插入成功: %s", item.language, result)

			// 验证文本是否正确插入
			// 可以添加字符编码和显示验证
		}

		// 清空输入，准备下一个测试
		time.Sleep(300 * time.Millisecond)

		// 清空输入框
		params["text"] = ""
		CDPInputInsertText(params)

		time.Sleep(200 * time.Millisecond)
	}

	log.Println("多语言文本插入测试完成")
}



*/

// -----------------------------------------------  Input.setInterceptDrags  -----------------------------------------------
// === 应用场景 ===
// 1. 拖拽拦截: 拦截并自定义处理拖拽事件
// 2. 拖拽分析: 分析拖拽事件的数据和行为
// 3. 安全防护: 防止恶意拖拽操作
// 4. 自定义拖拽UI: 拦截原生拖拽实现自定义UI
// 5. 拖拽调试: 调试拖拽相关的问题
// 6. 拖拽测试: 测试拖拽拦截功能

// CDPInputSetInterceptDrags 设置是否拦截拖拽事件
func CDPInputSetInterceptDrags(enabled bool) (string, error) {
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
		"method": "Input.setInterceptDrags",
		"params": {
			"enabled": %t
		}
	}`, reqID, enabled)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 setInterceptDrags 请求失败: %w", err)
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
			return "", fmt.Errorf("setInterceptDrags 请求超时")
		}
	}
}

/*
// 场景描述：拦截拖拽事件以实现自定义拖拽行为
func EnableDragInterceptionForCustomHandling() {
	log.Println("开始启用拖拽拦截...")
	
	// 启用拖拽拦截
	result, err := CDPInputSetInterceptDrags(true)
	if err != nil {
		log.Printf("启用拖拽拦截失败: %v", err)
		return
	}
	
	log.Printf("拖拽拦截已启用: %s", result)
	
	// 现在可以监听拖拽事件并进行自定义处理
	// 例如：修改拖拽数据、自定义拖拽UI等
	
	// 执行自定义拖拽操作
	log.Println("正在处理自定义拖拽逻辑...")
	
	// 处理完成后，如果需要，可以禁用拦截
	time.Sleep(2 * time.Second)
	
	// 禁用拖拽拦截
	result, err = CDPInputSetInterceptDrags(false)
	if err != nil {
		log.Printf("禁用拖拽拦截失败: %v", err)
	} else {
		log.Printf("拖拽拦截已禁用: %s", result)
	}
	
	log.Println("自定义拖拽处理完成")
}

// 场景描述：拦截拖拽事件以分析拖拽行为和调试问题
func AnalyzeDragEvents() {
	log.Println("开始拖拽事件分析...")
	
	// 启用拖拽拦截
	_, err := CDPInputSetInterceptDrags(true)
	if err != nil {
		log.Printf("无法启用拖拽拦截: %v", err)
		return
	}
	
	log.Println("拖拽拦截已启用，开始记录拖拽事件")
	
	// 模拟一些拖拽操作
	// 这里可以触发各种拖拽事件进行分析
	
	// 示例：监听特定时间内的拖拽事件
	analysisDuration := 10 * time.Second
	log.Printf("将在 %v 内分析拖拽事件...", analysisDuration)
	
	// 在这段时间内，可以手动或自动触发拖拽操作
	// 拖拽事件将被拦截，可以分析事件数据
	
	time.Sleep(analysisDuration)
	
	// 分析完成后禁用拦截
	result, err := CDPInputSetInterceptDrags(false)
	if err != nil {
		log.Printf("禁用拖拽拦截失败: %v", err)
	} else {
		log.Printf("拖拽拦截已禁用: %s", result)
	}
	
	// 输出分析结果
	log.Println("拖拽事件分析完成")
	log.Println("可以检查拦截到的拖拽事件数据进行分析和调试")
}

// 场景描述：启用拖拽拦截以防止恶意拖拽操作
func EnableDragSecurity() {
	log.Println("启用拖拽安全防护...")
	
	// 记录当前状态
	originalState := false // 假设初始状态为不拦截
	
	// 启用拖拽拦截
	result, err := CDPInputSetInterceptDrags(true)
	if err != nil {
		log.Printf("无法启用拖拽安全防护: %v", err)
		return
	}
	
	log.Printf("拖拽安全防护已启用: %s", result)
	
	// 执行敏感操作
	log.Println("正在执行敏感操作，拖拽被临时锁定...")
	
	// 这里可以执行需要保护的敏感操作
	// 例如：处理支付信息、敏感数据等
	
	time.Sleep(3 * time.Second)
	
	// 操作完成后恢复原始状态
	result, err = CDPInputSetInterceptDrags(originalState)
	if err != nil {
		log.Printf("无法恢复拖拽状态: %v", err)
	} else {
		log.Printf("拖拽状态已恢复: %s", result)
	}
	
	log.Println("拖拽安全防护完成")
}

// 场景描述：拦截原生拖拽以实现自定义拖拽UI
func ImplementCustomDragUI() {
	log.Println("开始实现自定义拖拽UI...")
	
	// 启用拖拽拦截
	_, err := CDPInputSetInterceptDrags(true)
	if err != nil {
		log.Printf("无法启用自定义拖拽UI: %v", err)
		return
	}
	
	log.Println("原生拖拽已被拦截，可以显示自定义拖拽UI")
	
	// 显示自定义拖拽UI
	log.Println("显示自定义拖拽效果...")
	
	// 这里可以实现：
	// 1. 自定义拖拽预览图像
	// 2. 自定义拖拽光标
	// 3. 自定义拖拽动画
	// 4. 自定义放置效果
	
	// 模拟自定义拖拽过程
	customDragSteps := []string{
		"显示自定义拖拽开始效果",
		"更新自定义拖拽预览",
		"显示自定义放置目标高亮",
		"执行自定义放置动画",
		"完成自定义拖拽",
	}
	
	for _, step := range customDragSteps {
		log.Printf("执行: %s", step)
		time.Sleep(500 * time.Millisecond)
	}
	
	// 完成后禁用拦截
	result, err := CDPInputSetInterceptDrags(false)
	if err != nil {
		log.Printf("禁用拖拽拦截失败: %v", err)
	} else {
		log.Printf("自定义拖拽UI已结束: %s", result)
	}
	
	log.Println("自定义拖拽UI实现完成")
}


// 场景描述：创建拖拽功能测试框架
func CreateDragTestFramework() {
	log.Println("初始化拖拽测试框架...")
	
	// 启用拖拽拦截以监控测试
	_, err := CDPInputSetInterceptDrags(true)
	if err != nil {
		log.Printf("无法初始化拖拽测试框架: %v", err)
		return
	}
	
	log.Println("拖拽测试框架已就绪，开始执行测试用例...")
	
	// 执行各种拖拽测试
	testCases := []struct {
		name        string
		description string
	}{
		{"文件拖拽上传", "测试文件拖拽到上传区域"},
		{"元素拖拽排序", "测试列表元素拖拽排序"},
		{"跨窗口拖拽", "测试窗口间的拖拽功能"},
		{"拖拽取消", "测试拖拽中途取消"},
		{"多点触控拖拽", "测试多点触控拖拽"},
	}
	
	for _, tc := range testCases {
		log.Printf("执行测试用例: %s - %s", tc.name, tc.description)
		
		// 这里可以执行具体的拖拽测试
		// 拖拽事件将被拦截，可以验证事件数据
		
		// 模拟测试执行
		time.Sleep(1 * time.Second)
		
		log.Printf("测试用例完成: %s", tc.name)
	}
	
	// 测试完成后禁用拦截
	result, err := CDPInputSetInterceptDrags(false)
	if err != nil {
		log.Printf("无法关闭测试框架: %v", err)
	} else {
		log.Printf("拖拽测试框架已关闭: %s", result)
	}
	
	// 生成测试报告
	log.Println("拖拽测试完成，生成测试报告...")
	log.Println("测试框架执行完毕")
}

// 场景描述：在系统维护期间临时禁用拖拽功能
func TemporaryDisableDragForMaintenance() {
	log.Println("开始系统维护，临时禁用拖拽功能...")
	
	// 记录维护开始时间
	maintenanceStart := time.Now()
	
	// 启用拖拽拦截
	result, err := CDPInputSetInterceptDrags(true)
	if err != nil {
		log.Printf("无法开始维护模式: %v", err)
		return
	}
	
	log.Printf("维护模式已启动: %s", result)
	log.Println("拖拽功能已被临时禁用")
	
	// 执行维护操作
	maintenanceTasks := []string{
		"清理临时文件",
		"更新拖拽配置",
		"优化拖拽性能",
		"修复拖拽相关bug",
	}
	
	for _, task := range maintenanceTasks {
		log.Printf("执行维护任务: %s", task)
		time.Sleep(1 * time.Second)
	}
	
	// 计算维护时长
	maintenanceDuration := time.Since(maintenanceStart)
	
	// 维护完成，恢复拖拽功能
	result, err = CDPInputSetInterceptDrags(false)
	if err != nil {
		log.Printf("无法结束维护模式: %v", err)
	} else {
		log.Printf("维护模式已结束: %s", result)
		log.Printf("维护总时长: %v", maintenanceDuration)
		log.Println("拖拽功能已恢复")
	}
	
	log.Println("系统维护完成")
}

 */
