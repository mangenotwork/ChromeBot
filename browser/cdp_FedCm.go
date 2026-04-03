package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------  FedCm.clickDialogButton  -----------------------------------------------
// === 应用场景 ===
// 1. 对话框操作自动化: 自动化点击FedCM对话框中的按钮
// 2. 用户交互模拟: 模拟用户在FedCM对话框中的点击行为
// 3. 测试流程自动化: 自动化测试FedCM对话框的交互流程
// 4. 按钮功能验证: 验证对话框按钮的功能正确性
// 5. 用户体验测试: 测试对话框按钮的交互体验
// 6. 错误处理测试: 测试点击按钮后的错误处理流程

// CDPFedCMClickDialogButton 点击FedCM对话框按钮
// dialogID: 对话框ID
// buttonIndex: 按钮索引
func CDPFedCMClickDialogButton(dialogID string, buttonIndex int) (string, error) {
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
        "method": "FedCM.clickDialogButton",
        "params": {
            "dialogId": "%s",
            "buttonIndex": %d
        }
    }`, reqID, dialogID, buttonIndex)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FedCM.clickDialogButton 请求失败: %w", err)
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
			return "", fmt.Errorf("FedCM.clickDialogButton 请求超时")
		}
	}
}

// 辅助函数: CDPFedCMClickDialogConfirmButton 点击确认按钮
func CDPFedCMClickDialogConfirmButton(dialogID string) (string, error) {
	// 假设确认按钮通常是索引0
	return CDPFedCMClickDialogButton(dialogID, 0)
}

// 辅助函数: CDPFedCMClickDialogCancelButton 点击取消按钮
func CDPFedCMClickDialogCancelButton(dialogID string) (string, error) {
	// 假设取消按钮通常是索引1
	return CDPFedCMClickDialogButton(dialogID, 1)
}

// 辅助函数: CDPFedCMClickDialogButtonWithValidation 带验证的点击按钮
func CDPFedCMClickDialogButtonWithValidation(dialogID string, buttonIndex int, expectedResult string) (bool, error) {
	result, err := CDPFedCMClickDialogButton(dialogID, buttonIndex)
	if err != nil {
		return false, fmt.Errorf("点击按钮失败: %w", err)
	}

	// 简单验证结果是否包含预期内容
	if expectedResult != "" && !strings.Contains(result, expectedResult) {
		return false, fmt.Errorf("验证失败: 结果不包含预期内容。结果: %s", result)
	}

	return true, nil
}

/*

// 示例1: 完整的FedCM对话框交互流程
func CompleteFedCMDialogFlow() {
    log.Println("=== FedCM对话框交互测试 ===")

    // 1. 启用FedCM
    log.Println("步骤1: 启用FedCM...")
    if _, err := CDPFedCMEnable(); err != nil {
        log.Printf("启用FedCM失败: %v", err)
        return
    }
    defer CDPFedCMDisable()

    // 2. 假设FedCM对话框已显示，获取对话框ID
    dialogID := "fedcm_dialog_001"

    // 3. 点击对话框的确认按钮
    log.Println("步骤3: 点击确认按钮...")
    confirmResult, err := CDPFedCMClickDialogConfirmButton(dialogID)
    if err != nil {
        log.Printf("点击确认按钮失败: %v", err)

        // 4. 如果确认失败，尝试点击取消按钮
        log.Println("尝试点击取消按钮...")
        cancelResult, err := CDPFedCMClickDialogCancelButton(dialogID)
        if err != nil {
            log.Printf("点击取消按钮也失败: %v", err)
        } else {
            log.Printf("取消按钮点击成功: %s", cancelResult)
        }
    } else {
        log.Printf("确认按钮点击成功: %s", confirmResult)
    }

    // 5. 重置冷却期
    log.Println("步骤5: 重置冷却期...")
    if _, err := CDPFedCMResetCooldown(); err != nil {
        log.Printf("重置冷却期失败: %v", err)
    }
}

// 示例2: 测试不同按钮的交互
func TestAllDialogButtons() {
    dialogID := "test_dialog_buttons"

    log.Println("=== 测试FedCM对话框所有按钮 ===")

    if _, err := CDPFedCMEnable(); err != nil {
        log.Printf("启用FedCM失败: %v", err)
        return
    }
    defer CDPFedCMDisable()

    // 测试按钮索引0-3（假设最多4个按钮）
    for i := 0; i < 4; i++ {
        log.Printf("测试按钮索引 %d...", i)

        result, err := CDPFedCMClickDialogButton(dialogID, i)
        if err != nil {
            log.Printf("按钮 %d 点击失败: %v", i, err)

            // 如果点击失败，可能是按钮不存在，记录并继续
            if strings.Contains(err.Error(), "buttonIndex") {
                log.Printf("按钮索引 %d 可能不存在，继续测试其他按钮", i)
            }
        } else {
            log.Printf("按钮 %d 点击成功: %s", i, result)
        }

        // 重置冷却期以便下一次测试
        CDPFedCMResetCooldown()
        time.Sleep(500 * time.Millisecond) // 短暂延迟
    }
}

*/

// -----------------------------------------------  FedCm.disable  -----------------------------------------------
// === 应用场景 ===
// 1. 测试清理: 测试完成后禁用FedCM功能
// 2. 环境重置: 重置浏览器FedCM状态
// 3. 错误恢复: FedCM功能异常时禁用
// 4. 功能切换: 切换不同的认证测试场景
// 5. 资源释放: 释放FedCM相关资源
// 6. 测试隔离: 隔离不同测试的FedCM状态

// CDPFedCMDisable 禁用FedCM
func CDPFedCMDisable() (string, error) {
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
        "method": "FedCM.disable"
    }`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FedCM.disable 请求失败: %w", err)
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
			return "", fmt.Errorf("FedCM.disable 请求超时")
		}
	}
}

// -----------------------------------------------  FedCm.dismissDialog  -----------------------------------------------
// === 应用场景 ===
// 1. 对话框处理: 关闭FedCM认证对话框
// 2. 取消操作测试: 测试用户取消认证的流程
// 3. 错误处理测试: 测试认证取消的错误处理
// 4. 用户体验测试: 测试对话框关闭体验
// 5. 中断测试: 测试认证流程被中断的场景
// 6. 超时处理: 处理认证对话框超时

// CDPFedCMDismissDialog 关闭FedCM对话框
func CDPFedCMDismissDialog(dialogID string, trigger string) (string, error) {
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
        "method": "FedCM.dismissDialog",
        "params": {
            "dialogId": "%s",
            "trigger": "%s"
        }
    }`, reqID, dialogID, trigger)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FedCM.dismissDialog 请求失败: %w", err)
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
			return "", fmt.Errorf("FedCM.dismissDialog 请求超时")
		}
	}
}

// -----------------------------------------------  FedCm.enable  -----------------------------------------------
// === 应用场景 ===
// 1. 联邦认证测试: 启用FedCM功能进行联邦认证测试
// 2. 身份验证自动化: 自动化测试联邦身份验证流程
// 3. 兼容性测试: 测试浏览器对FedCM的支持
// 4. 调试支持: 调试联邦身份验证相关问题
// 5. 安全测试: 测试FedCM的安全特性
// 6. 性能测试: 测试FedCM的性能表现

// CDPFedCMEnable 启用FedCM
func CDPFedCMEnable() (string, error) {
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
        "method": "FedCM.enable"
    }`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FedCM.enable 请求失败: %w", err)
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
			return "", fmt.Errorf("FedCM.enable 请求超时")
		}
	}
}

// -----------------------------------------------  FedCm.openUrl  -----------------------------------------------
// === 应用场景 ===
// 1. 账户注册: 打开身份提供商的注册页面
// 2. 账户管理: 打开账户管理页面
// 3. 权限设置: 打开权限设置页面
// 4. 帮助文档: 打开身份提供商的帮助文档
// 5. 隐私政策: 打开隐私政策页面
// 6. 服务条款: 打开服务条款页面

// CDPFedCMOpenUrl 打开FedCM相关URL
// configURL: 身份提供商的配置URL
// accountURL: 账户相关URL
// loginURL: 登录相关URL
func CDPFedCMOpenUrl(configURL, accountURL, loginURL string) (string, error) {
	if !DefaultBrowserWS() {
		return "", fmt.Errorf("CDP功能未启用")
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("浏览器WebSocket连接未建立")
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建参数
	var paramsJSON strings.Builder
	paramsJSON.WriteString("{")

	if configURL != "" {
		paramsJSON.WriteString(fmt.Sprintf(`"configUrl": "%s",`, configURL))
	}
	if accountURL != "" {
		paramsJSON.WriteString(fmt.Sprintf(`"accountUrl": "%s",`, accountURL))
	}
	if loginURL != "" {
		paramsJSON.WriteString(fmt.Sprintf(`"loginUrl": "%s",`, loginURL))
	}

	// 移除最后一个逗号
	paramsStr := paramsJSON.String()
	paramsStr = strings.TrimSuffix(paramsStr, ",")
	paramsStr += "}"

	// 构建消息
	message := fmt.Sprintf(`{
        "id": %d,
        "method": "FedCM.openUrl",
        "params": %s
    }`, reqID, paramsStr)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FedCM.openUrl 请求失败: %w", err)
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
			return "", fmt.Errorf("FedCM.openUrl 请求超时")
		}
	}
}

// 辅助函数: CDPFedCMOpenConfigUrl 打开身份提供商配置URL
func CDPFedCMOpenConfigUrl(configURL string) (string, error) {
	return CDPFedCMOpenUrl(configURL, "", "")
}

// 辅助函数: CDPFedCMOpenAccountUrl 打开账户管理URL
func CDPFedCMOpenAccountUrl(accountURL string) (string, error) {
	return CDPFedCMOpenUrl("", accountURL, "")
}

// 辅助函数: CDPFedCMOpenLoginUrl 打开登录URL
func CDPFedCMOpenLoginUrl(loginURL string) (string, error) {
	return CDPFedCMOpenUrl("", "", loginURL)
}

/*

// 示例1: 打开身份提供商配置页面
func OpenIdentityProviderConfig() {
    log.Println("=== 打开身份提供商配置页面 ===")

    configURL := "https://accounts.google.com/.well-known/openid-configuration"

    result, err := CDPFedCMOpenConfigUrl(configURL)
    if err != nil {
        log.Printf("打开配置页面失败: %v", err)
        return
    }

    log.Printf("配置页面已打开: %s", result)

    // 解析响应
    var response struct {
        Result struct {
            Success bool   `json:"success"`
            Message string `json:"message"`
        } `json:"result"`
    }

    if err := json.Unmarshal([]byte(result), &response); err == nil {
        if response.Result.Success {
            log.Println("✓ 成功打开身份提供商配置")
        } else {
            log.Printf("⚠ 打开配置有问题: %s", response.Result.Message)
        }
    }
}

// 示例2: 打开账户管理页面
func OpenAccountManagementPage() {
    log.Println("=== 打开账户管理页面 ===")

    accountURL := "https://myaccount.google.com/"

    result, err := CDPFedCMOpenAccountUrl(accountURL)
    if err != nil {
        log.Printf("打开账户页面失败: %v", err)
        return
    }

    log.Printf("账户管理页面已打开: %s", result)
}

// 示例3: 完整的FedCM URL导航流程
func CompleteFedCMUrlNavigation() {
    log.Println("=== FedCM完整URL导航流程 ===")

    // 启用FedCM
    if _, err := CDPFedCMEnable(); err != nil {
        log.Printf("启用FedCM失败: %v", err)
        return
    }
    defer CDPFedCMDisable()

    // 1. 打开配置URL
    log.Println("步骤1: 打开身份提供商配置...")
    configResult, err := CDPFedCMOpenConfigUrl("https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration")
    if err != nil {
        log.Printf("打开配置URL失败: %v", err)
    } else {
        log.Printf("配置URL打开结果: %s", configResult)
    }

    // 2. 打开账户URL
    log.Println("步骤2: 打开账户管理页面...")
    accountResult, err := CDPFedCMOpenAccountUrl("https://account.microsoft.com/account")
    if err != nil {
        log.Printf("打开账户URL失败: %v", err)
    } else {
        log.Printf("账户URL打开结果: %s", accountResult)
    }

    // 3. 打开登录URL
    log.Println("步骤3: 打开登录页面...")
    loginResult, err := CDPFedCMOpenLoginUrl("https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
    if err != nil {
        log.Printf("打开登录URL失败: %v", err)
    } else {
        log.Printf("登录URL打开结果: %s", loginResult)
    }

    log.Println("=== FedCM URL导航流程完成 ===")
}

// 示例4: 批量打开多个身份提供商URL
func BatchOpenIdentityProviderURLs() {
    log.Println("=== 批量打开身份提供商URL ===")

    // 定义多个身份提供商的URL
    identityProviders := []struct {
        Name      string
        ConfigURL string
        AccountURL string
        LoginURL   string
    }{
        {
            Name:      "Google",
            ConfigURL: "https://accounts.google.com/.well-known/openid-configuration",
            AccountURL: "https://myaccount.google.com/",
            LoginURL:   "https://accounts.google.com/o/oauth2/auth",
        },
        {
            Name:      "Microsoft",
            ConfigURL: "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration",
            AccountURL: "https://account.microsoft.com/account",
            LoginURL:   "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
        },
        {
            Name:      "GitHub",
            ConfigURL: "https://github.com/.well-known/openid-configuration",
            AccountURL: "https://github.com/settings/profile",
            LoginURL:   "https://github.com/login/oauth/authorize",
        },
    }

    for _, provider := range identityProviders {
        log.Printf("打开 %s 身份提供商URL...", provider.Name)

        result, err := CDPFedCMOpenUrl(provider.ConfigURL, provider.AccountURL, provider.LoginURL)
        if err != nil {
            log.Printf("打开 %s URL失败: %v", provider.Name, err)
        } else {
            log.Printf("%s URL打开结果: %s", provider.Name, result)
        }

        // 短暂延迟
        time.Sleep(500 * time.Millisecond)
    }
}

// 示例5: FedCM URL测试工具
type FedCMUrlTester struct {
    TestURLs []TestURL
    Results  []UrlTestResult
    Timeout  time.Duration
}

type TestURL struct {
    ID        string
    URLType   string // "config", "account", "login"
    URL       string
    Expected  string
    Mandatory bool
}

type UrlTestResult struct {
    URLID     string
    Success   bool
    Error     string
    Duration  time.Duration
    Response  string
    Timestamp time.Time
}

func (tester *FedCMUrlTester) RunTests() {
    log.Printf("开始FedCM URL测试，共 %d 个URL", len(tester.TestURLs))

    // 启用FedCM
    if _, err := CDPFedCMEnable(); err != nil {
        log.Printf("启用FedCM失败: %v", err)
        return
    }
    defer CDPFedCMDisable()

    for _, testURL := range tester.TestURLs {
        log.Printf("测试URL: %s (%s)", testURL.ID, testURL.URLType)

        startTime := time.Now()

        var result string
        var err error

        switch testURL.URLType {
        case "config":
            result, err = CDPFedCMOpenConfigUrl(testURL.URL)
        case "account":
            result, err = CDPFedCMOpenAccountUrl(testURL.URL)
        case "login":
            result, err = CDPFedCMOpenLoginUrl(testURL.URL)
        default:
            err = fmt.Errorf("未知的URL类型: %s", testURL.URLType)
        }

        testResult := UrlTestResult{
            URLID:     testURL.ID,
            Timestamp: time.Now(),
            Duration:  time.Since(startTime),
            Response:  result,
        }

        if err != nil {
            testResult.Success = false
            testResult.Error = err.Error()
            if testURL.Mandatory {
                log.Printf("重要URL测试失败: %s - %v", testURL.ID, err)
            } else {
                log.Printf("URL测试失败: %s - %v", testURL.ID, err)
            }
        } else {
            testResult.Success = true
            log.Printf("URL测试成功: %s，耗时: %v", testURL.ID, testResult.Duration)

            // 验证响应
            if testURL.Expected != "" {
                if !strings.Contains(result, testURL.Expected) {
                    testResult.Success = false
                    testResult.Error = "响应不包含预期内容"
                    log.Printf("响应验证失败: 期望包含 '%s'", testURL.Expected)
                }
            }
        }

        tester.Results = append(tester.Results, testResult)

        // 短暂延迟
        time.Sleep(200 * time.Millisecond)
    }

    tester.GenerateReport()
}

func (tester *FedCMUrlTester) GenerateReport() {
    fmt.Println("=== FedCM URL测试报告 ===")
    fmt.Printf("测试时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Printf("测试URL数量: %d\n\n", len(tester.Results))

    total := len(tester.Results)
    passed := 0
    totalDuration := time.Duration(0)

    for _, result := range tester.Results {
        if result.Success {
            passed++
        }

        totalDuration += result.Duration

        status := "✓ 通过"
        if !result.Success {
            status = "✗ 失败"
        }

        fmt.Printf("%s %-20s 类型: %-8s 耗时: %v\n",
            status, result.URLID, getURLType(result.URLID), result.Duration)
        if result.Error != "" {
            fmt.Printf("   错误: %s\n", result.Error)
        }
    }

    avgDuration := totalDuration / time.Duration(total)
    successRate := float64(passed) / float64(total) * 100

    fmt.Printf("\n总计: %d/%d 通过 (%.1f%%)\n", passed, total, successRate)
    fmt.Printf("平均耗时: %v\n", avgDuration)
    fmt.Printf("总测试时间: %v\n", totalDuration)
}

func getURLType(urlID string) string {
    if strings.Contains(urlID, "config") {
        return "config"
    } else if strings.Contains(urlID, "account") {
        return "account"
    } else if strings.Contains(urlID, "login") {
        return "login"
    }
    return "unknown"
}

// 示例6: FedCM集成场景测试
func TestFedCMWithUrlNavigation() {
    log.Println("=== FedCM集成场景测试（包含URL导航）===")

    // 启用FedCM
    if _, err := CDPFedCMEnable(); err != nil {
        log.Printf("启用FedCM失败: %v", err)
        return
    }
    defer CDPFedCMDisable()

    // 场景1: 新用户注册流程
    log.Println("场景1: 新用户注册流程")
    testNewUserRegistration()

    // 场景2: 现有用户登录流程
    log.Println("\n场景2: 现有用户登录流程")
    testExistingUserLogin()

    // 场景3: 账户管理流程
    log.Println("\n场景3: 账户管理流程")
    testAccountManagement()

    log.Println("=== FedCM集成场景测试完成 ===")
}

func testNewUserRegistration() {
    // 1. 打开身份提供商配置
    log.Println("步骤1: 查看身份提供商配置...")
    result, err := CDPFedCMOpenConfigUrl("https://idp.example.com/.well-known/openid-configuration")
    if err != nil {
        log.Printf("查看配置失败: %v", err)
        return
    }
    log.Printf("配置信息: %s", result[:100]) // 只显示前100字符

    // 2. 打开注册页面
    log.Println("步骤2: 打开注册页面...")
    regResult, err := CDPFedCMOpenLoginUrl("https://idp.example.com/register")
    if err != nil {
        log.Printf("打开注册页面失败: %v", err)
    } else {
        log.Printf("注册页面打开结果: %s", regResult)
    }

    // 3. 打开服务条款
    log.Println("步骤3: 打开服务条款页面...")
    termsResult, err := CDPFedCMOpenUrl("", "https://idp.example.com/terms", "")
    if err != nil {
        log.Printf("打开服务条款失败: %v", err)
    } else {
        log.Printf("服务条款页面打开结果: %s", termsResult)
    }
}

func testExistingUserLogin() {
    // 1. 打开登录页面
    log.Println("步骤1: 打开登录页面...")
    result, err := CDPFedCMOpenLoginUrl("https://idp.example.com/login")
    if err != nil {
        log.Printf("打开登录页面失败: %v", err)
        return
    }
    log.Printf("登录页面打开结果: %s", result)

    // 2. 打开账户管理
    log.Println("步骤2: 打开账户管理页面...")
    accountResult, err := CDPFedCMOpenAccountUrl("https://idp.example.com/account")
    if err != nil {
        log.Printf("打开账户管理失败: %v", err)
    } else {
        log.Printf("账户管理页面打开结果: %s", accountResult)
    }

    // 3. 打开隐私政策
    log.Println("步骤3: 打开隐私政策页面...")
    privacyResult, err := CDPFedCMOpenUrl("", "https://idp.example.com/privacy", "")
    if err != nil {
        log.Printf("打开隐私政策失败: %v", err)
    } else {
        log.Printf("隐私政策页面打开结果: %s", privacyResult)
    }
}

func testAccountManagement() {
    // 1. 打开账户概览
    log.Println("步骤1: 打开账户概览页面...")
    result, err := CDPFedCMOpenAccountUrl("https://idp.example.com/account/overview")
    if err != nil {
        log.Printf("打开账户概览失败: %v", err)
        return
    }
    log.Printf("账户概览打开结果: %s", result)

    // 2. 打开安全设置
    log.Println("步骤2: 打开安全设置页面...")
    securityResult, err := CDPFedCMOpenUrl("", "https://idp.example.com/account/security", "")
    if err != nil {
        log.Printf("打开安全设置失败: %v", err)
    } else {
        log.Printf("安全设置页面打开结果: %s", securityResult)
    }

    // 3. 打开连接的应用程序
    log.Println("步骤3: 打开连接的应用程序页面...")
    appsResult, err := CDPFedCMOpenUrl("", "https://idp.example.com/account/apps", "")
    if err != nil {
        log.Printf("打开应用程序页面失败: %v", err)
    } else {
        log.Printf("应用程序页面打开结果: %s", appsResult)
    }
}


*/

// -----------------------------------------------  FedCm.resetCooldown  -----------------------------------------------
// === 应用场景 ===
// 1. 冷却期重置: 重置FedCM认证的冷却期
// 2. 测试重复认证: 测试连续认证的场景
// 3. 自动化测试: 自动化测试中重置认证限制
// 4. 性能测试: 测试高频认证的性能
// 5. 压力测试: 对认证系统进行压力测试
// 6. 边界测试: 测试认证频率的边界情况

// CDPFedCMResetCooldown 重置FedCM冷却期
func CDPFedCMResetCooldown() (string, error) {
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
        "method": "FedCM.resetCooldown"
    }`, reqID)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FedCM.resetCooldown 请求失败: %w", err)
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
			return "", fmt.Errorf("FedCM.resetCooldown 请求超时")
		}
	}
}

// -----------------------------------------------  FedCm.selectAccount  -----------------------------------------------
// === 应用场景 ===
// 1. 账户选择自动化: 自动化联邦认证的账户选择流程
// 2. 测试用例模拟: 模拟用户选择账户的行为
// 3. 多账户测试: 测试多账户选择场景
// 4. 用户体验测试: 测试账户选择界面的用户体验
// 5. 安全性测试: 测试账户选择的安全性
// 6. 性能测试: 测试账户选择的性能

// CDPFedCMSelectAccount 选择FedCM账户
func CDPFedCMSelectAccount(dialogID string, accountIndex int) (string, error) {
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
        "method": "FedCM.selectAccount",
        "params": {
            "dialogId": "%s",
            "accountIndex": %d
        }
    }`, reqID, dialogID, accountIndex)

	// 发送请求
	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", fmt.Errorf("发送 FedCM.selectAccount 请求失败: %w", err)
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
			return "", fmt.Errorf("FedCM.selectAccount 请求超时")
		}
	}
}
