package builtins

import (
	"ChromeBot/browser"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/utils"
	"encoding/json"
	"fmt"

	gt "github.com/mangenotwork/gathertool"
)

// GetMainWindowID 获取主窗口ID ex: chrome cdpfn=GetMainWindowID to=wid
// GetCurrentWindowInfo 获取当前活动窗口的信息  ex: chrome cdpfn=GetMainWindowID to=wid
// CDPBrowserSetContentsSize设置浏览器内容区域尺寸  ex: chrome cdpfn=CDPBrowserSetContentsSize params=`{"windowId":123, "width":900, "height":600, "keepPosition":false, "includeChrome":false}`
// 参数说明: keepPosition:是否保持当前位置(可选),includeChrome:是否包括浏览器边框(可选),windowState(可选):窗口状态(normal:正常窗口, minimized:最小化, maximized:最大化, fullscreen:全屏)left(可选):指定X坐标, top(可选):指定Y坐标
// NewTab 新建页签并返回sessionId，改方法会默认切换到这个新的页签  ex: chrome cdpfn=NewTab params=`{"url":""}` to=sid
// DOMStructureAnalysis DOM结构分析工具  ex: chrome cdpfn=DOMStructureAnalysis
// PageComparisonAtTime 指定时间页面对比分析   ex:  chrome cdpfn=PageComparisonAtTime  params=`{"second":5}`
// ClearLocalStorage 清除指定源的localStorage  ex: chrome cdpfn=ClearLocalStorage  params=`{"origin":"https://example.com"}`
// ClearSessionStorage 清除指定源的sessionStorage  ex: chrome cdpfn=ClearSessionStorage  params=`{"origin":"https://example.com"}`
// ComponentLibraryTest  UI组件库的交互状态测试测试按钮的各种状态  ex: chrome cdpfn=ComponentLibraryTest  params=`{"buttonNodeId": 1}`
// ResponsiveDesignDebugger 响应式设计调试 - 分析媒体查询  ex: chrome cdpfn=ResponsiveDesignDebugger
func runCDPFN(interp *interpreter.Interpreter, cdpfn string, params, to string) {

	fmt.Println("cdp = ", cdpfn)
	fmt.Println("params = ", params)

	paramsMap := make(map[string]any)
	err := json.Unmarshal([]byte(params), &paramsMap)
	if err != nil {
		fmt.Println("[Err] CDPFN params 解析失败 err:", err.Error())
		return
	}

	switch cdpfn {
	case "GetMainWindowID":
		windowId, err := browser.GetMainWindowID()
		if err != nil {
			fmt.Println("GetMainWindowID获取windowId失败, err: ", err.Error())
			break
		}

		fmt.Println("windowId = ", windowId)
		fmt.Println("to = ", to)

		interp.Global().SetVar(to, windowId)

	case "GetCurrentWindowInfo":
		info, err := browser.GetCurrentWindowInfo()
		if err != nil {
			fmt.Println("GetCurrentWindowInfo失败, err: ", err.Error())
			break
		}
		fmt.Println("窗口信息: ")
		utils.ShowJson(info)

	case "CDPBrowserSetContentsSize":
		windowId, ok := paramsMap["windowId"]
		if !ok {
			fmt.Println("未设置参数 windowId")
			break
		}
		width, ok := paramsMap["width"]
		if !ok {
			fmt.Println("未设置参数 width")
			break
		}
		height, ok := paramsMap["height"]
		if !ok {
			fmt.Println("未设置参数 height")
			break
		}

		// todo 可选参数如何优雅的传入, 或者在设计上如何避免这种不确定性，建议不到再一个方法上扩展太多可选参数，应该对应的场景进行方法拆解，
		// 个人认为这样具有确定性，AI在识别和按文档生成的时候也具备了确定性从而避免了不确定导致的问题

		windowIdInt := gt.Any2Int(windowId)
		widthInt := gt.Any2Int(width)
		heightInt := gt.Any2Int(height)
		browser.CDPBrowserSetContentsSizeFn(windowIdInt, widthInt, heightInt)

	case "NewTab":
		urlStr, ok := paramsMap["url"].(string)
		if !ok {
			fmt.Println("未设置参数 url")
			break
		}
		_, err := browser.NewTab()
		if err != nil {
			fmt.Println("打开页面失败,err:", err)
			break
		}
		browser.OpenUrl(urlStr)
		interp.Global().SetVar(to, browser.GetNowTabSession())

	case "DOMStructureAnalysis":
		browser.DOMStructureAnalysis()

	case "PageComparisonAtTime":
		second, ok := paramsMap["second"]
		if !ok {
			fmt.Println("未设置参数 second")
			break
		}
		browser.PageComparisonAtTime(gt.Any2Int(second))

	case "ClearLocalStorage":
		origin, ok := paramsMap["origin"].(string)
		if !ok {
			fmt.Println("未设置参数 origin")
			break
		}
		browser.ClearLocalStorage(origin)

	case "ClearSessionStorage":
		origin, ok := paramsMap["origin"].(string)
		if !ok {
			fmt.Println("未设置参数 origin")
			break
		}
		browser.ClearSessionStorage(origin)

	case "ComponentLibraryTest":
		browser.ComponentLibraryTest(paramsMap["buttonNodeId"].(int))

	case "ResponsiveDesignDebugger":
		browser.ResponsiveDesignDebugger()

	}

}
