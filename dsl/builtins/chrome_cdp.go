package builtins

import (
	"ChromeBot/browser"
	"ChromeBot/dsl/interpreter"
	"encoding/json"
	"fmt"

	gt "github.com/mangenotwork/gathertool"
)

// SystemInfo.getFeatureState 获取Feature状态 ex: chrome cdp=`SystemInfo.getFeatureState` params=`{"featureState":"webgl"}`  参数说明 feature：gpu_acceleration(GPU 加速),vulkan(Vulkan 渲染),direct3d11(D3D11),canvas_oop_rasterization(画布离屏渲染),video_acceleration(视频硬件加速),webgl,webgl2,webgpu
// SystemInfo.getInfo 获取系统信息信息 ex: chrome cdp=`SystemInfo.getInfo`
// SystemInfo.getProcessInfo 获取正在运行的进程的相关信息 ex: chrome cdp=`SystemInfo.getProcessInfo`
// Browser.close   关闭浏览器  ex: chrome cdp=`Browser.close`
// Browser.resetPermissions 重置权限 ex: chrome cdp=`Browser.resetPermissions` params=`{"origin": "https://example.com"}`
// Browser.getWindowForTarget  通过targetId获取对应的窗口ID ex: chrome cdp=`Browser.getWindowForTarget` params=`{"targetId": "..."}` to=wid
// Browser.setWindowBounds  设置浏览器窗口的大小。 ex: chrome cdp=`Browser.setWindowBounds` params=`{"windowId": "...", "left":100,"top":100,"width":800,"height":600,"windowState":"normal"}`    参数说明 windowState:窗口状态(normal:正常窗口, minimized:最小化, maximized:最大化, fullscreen:全屏)
// Browser.setContentsSize  设置浏览器窗口的位置和/或大小  ex: chrome cdp=`Browser.setContentsSize` params=`{"windowId": "...", "width":800,"height":600}`
// Target.createTarget  创建target  ex: chrome cdp=`Target.createTarget` params=`{"url":"https://example.com"}` to=tid
// Target.activateTarget 激活target 聚焦指定页面  ex: chrome cdp=`Target.activateTarget` params=`{"targetId":""}`
// Target.attachToTarget 聚焦目标页签返回sessionID ex: chrome cdp=`Target.attachToTarget` params=`{"targetId":""}` to=sid
// Target.closeTarget 关闭指定target,如果目标是页面，则页面也会被关闭。 ex: chrome cdp=`Target.closeTarget` params=`{"targetId":""}`
// Target.createBrowserContext 创建一个新的空浏览器上下文（它类似于浏览器的无痕模式） ex: chrome cdp=`Target.createBrowserContext`
// Target.detachFromTarget 分离掉指定sessionID  ex: chrome cdp=`Target.detachFromTarget` params=`{"sessionId":""}`
// Target.disposeBrowserContext  删除 BrowserContext  ex: chrome cdp=`Target.disposeBrowserContext` params=`{"browserContextId":""}`
// Target.getBrowserContexts 返回创建的所有浏览器上下文  ex: chrome cdp=`Target.getBrowserContexts`
// Target.getTargets  获取可用目标列表。  ex: chrome cdp=`Target.getTargets`
// Target.getTargetInfo  返回目标的相关信息   ex: chrome cdp=`Target.getTargetInfo` params=`{"targetId":""}`
// DOMSnapshot.captureSnapshot 返回文档快照，其中包含根节点的完整 DOM 树   ex: chrome cdp=`DOMSnapshot.captureSnapshot`
// DOMSnapshot.disable 禁用给定页面的 DOM 快照  ex: chrome cdp=`DOMSnapshot.disable`
// DOMSnapshot.enable  启用 DOM 快照   ex: chrome cdp=`DOMSnapshot.enable`
// DOMStorage.clear 清除指定存储区域的所有数据  ex: chrome cdp=`DOMStorage.clear` params=`{"securityOrigin":"https://example.com","isLocalStorage":true}`   securityOrigin:存储源 isLocalStorage(bool):是否是localStorage
// DOMStorage.disable 禁用存储跟踪    ex: chrome cdp=`DOMStorage.disable`
// DOMStorage.enable 启用存储跟踪功能  ex: chrome cdp=`DOMStorage.enable`
// DOMStorage.getDOMStorageItems  获取指定存储区域的所有项目   ex: chrome cdp=`DOMStorage.getDOMStorageItems` params=`{"securityOrigin":"https://example.com","isLocalStorage":true}`
// DOMStorage.removeDOMStorageItem  删除指定存储区域的特定项目  ex: chrome cdp=`DOMStorage.removeDOMStorageItem` params=`{"securityOrigin":"https://example.com","isLocalStorage":true, "key":""}`
// DOMStorage.setDOMStorageItem  在指定存储区域中设置项目  ex: chrome cdp=`DOMStorage.removeDOMStorageItem` params=`{"securityOrigin":"https://example.com","isLocalStorage":true, "key":"", value:""}`
// CSS.addRule  向样式表中添加新的CSS规则   ex: chrome cdp=`CSS.addRule` params=`{"styleSheetId": "1","rule": "div.test { background: blue; padding: 10px; }","index": 0}`
// CSS.collectClassNames 从指定样式表中收集所有类名  ex: chrome cdp=`CSS.collectClassNames` params=`{"styleSheetId": "1"}`
// CSS.enable 启用CSS域   ex: chrome cdp=`CSS.enable`
// CSS.Disable 禁用CSS域   ex: chrome cdp=`CSS.Disable`
// CSS.createStyleSheet 创建一个新的样式表   ex: chrome cdp=`CSS.createStyleSheet` params=`{"frameId": "YOUR_FRAME_ID", "force": false }`
// CSS.forcePseudoState 强制元素应用指定的伪类状态  ex: chrome cdp=`CSS.forcePseudoState` params=`{"nodeId": 123,"forcedPseudoClasses": ["hover", "focus"]}`
// CSS.forceStartingStyle 强制元素应用起始样式状态  ex: chrome cdp=`CSS.forceStartingStyle` params=`{"nodeId": 123,"focus": true}`
// CSS.getBackgroundColors 获取元素背后的背景颜色范围  ex: chrome cdp=`CSS.getBackgroundColors` params=`{"nodeId": 123}`
// CSS.getComputedStyleForNode 获取指定节点的计算样式  ex: chrome cdp=`CSS.getComputedStyleForNode` params=`{"nodeId": 123}`
// CSS.getInlineStylesForNode 获取指定节点的行内样式  ex: chrome cdp=`CSS.getInlineStylesForNode` params=`{"nodeId": 123}`
// CSS.getMatchedStylesForNode 获取指定节点的匹配样式  ex: chrome cdp=`CSS.getMatchedStylesForNode` params=`{"nodeId": 123}`
// CSS.getMediaQueries 获取所有媒体查询  ex: chrome cdp=`CSS.getMediaQueries`
// CSS.getPlatformFontsForNode  获取节点使用的平台字体信息  ex: chrome cdp=`CSS.getPlatformFontsForNode` params=`{"nodeId": 123}`
// CSS.getStyleSheetText 获取样式表文本  ex: chrome cdp=`CSS.getStyleSheetText` params=`{"styleSheetId": "123"}`
// CSS.setEffectivePropertyValueForNode 设置节点的属性值  ex: chrome cdp=`CSS.setEffectivePropertyValueForNode` params=`{"nodeId": 123, "propertyName": "color", "value": "red"}`
// CSS.setKeyframeKey 设置关键帧的键  ex: chrome cdp=`CSS.setKeyframeKey` params=`{"styleSheetId": "2:18","ruleIndex": 0,"keyIndex": 0,"key": "20%"}`
// CSS.setMediaText 设置媒体文本  ex: chrome cdp=`CSS.setMediaText` params=`{"styleSheetId": "2:0","ruleIndex": 0,"mediaText": "@media (max-width: 768px)"}`
// CSS.setPropertyRulePropertyName 设置属性规则属性名称 ex: chrome cdp=`CSS.setPropertyRulePropertyName` params=`{"styleSheetId": "2:0","ruleIndex": 0,"propertyIndex": 0,"name": "color"}`
// CSS.setRuleSelector 设置规则选择器 ex: chrome cdp=`CSS.setRuleSelector` params=`{"styleSheetId": "2:0","ruleIndex": 0,"selector": "body"}`
// CSS.setStyleSheetText 设置样式表的文本内容 ex: chrome cdp=`CSS.setStyleSheetText` params=`{"styleSheetId": "2:0","text": "body {color: red;}"}`
// CSS.setStyleTexts 设置样式文本 ex: chrome cdp=`CSS.setStyleTexts` params=`{"styleSheetId": "2:0","edits": [{"styleSheetId": "2:0","style": {"styleSheetId": "2:0","range": {"startLine": 0,"startColumn": 0,"endLine": 0,"endColumn": 0},"cssProperties": [{"name": "
// CSS.startRuleUsageTracking 开始规则使用跟踪 ex: chrome cdp=`CSS.startRuleUsageTracking`
// CSS.stopRuleUsageTracking 停止规则使用跟踪 ex: chrome cdp=`CSS.stopRuleUsageTracking`
// CSS.takeCoverageDelta 获取CSS规则使用跟踪结果 ex: chrome cdp=`CSS.takeCoverageDelta`
// CSS.getEnvironmentVariables  获取环境变量 ex: chrome cdp=`CSS.getEnvironmentVariables`
// CSS.setContainerQueryText 设置容器查询文本 ex: chrome cdp=`CSS.setContainerQueryText` params=`{"styleSheetId": "2:14","ruleIndex": 1,"containerQueryText": "(min-width: 400px)"}`
func runCDP(interp *interpreter.Interpreter, cdp string, params, to string) {

	fmt.Println("cdp = ", cdp)
	fmt.Println("params = ", params)

	paramsMap := make(map[string]any)
	err := json.Unmarshal([]byte(params), &paramsMap)
	if err != nil {
		fmt.Println("[Err] CDP params 解析失败 err:", err.Error())
		return
	}

	switch cdp {
	case "SystemInfo.getFeatureState":
		feature, ok := paramsMap["featureState"].(string)
		if !ok {
			fmt.Println("未设置参数 featureState")
			break
		}
		browser.CDPSystemInfoGetFeatureState(feature)

	case "SystemInfo.getInfo":
		browser.CDPSystemInfoGetInfo()

	case "SystemInfo.getProcessInfo":
		browser.CDPSystemInfoGetProcessInfo()

	case "Browser.close":
		browser.CDPBrowserClose()

	case "Browser.resetPermissions":
		origin, ok := paramsMap["origin"].(string)
		if !ok {
			fmt.Println("未设置参数 origin")
			break
		}
		browser.CDPBrowserResetPermissions(origin)

	case "Browser.getWindowForTarget":
		targetId, ok := paramsMap["targetId"].(string)
		if !ok {
			fmt.Println("未设置参数 targetId")
			break
		}
		windowId, err := browser.CDPBrowserGetWindowForTarget(targetId)
		if err != nil {
			fmt.Println("获取targetId对应的windowId失败, err: ", err.Error())
			break
		}
		interp.Global().SetVar(to, windowId)

	case "Browser.setWindowBounds":
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
		left, ok := paramsMap["left"]
		if !ok {
			fmt.Println("未设置参数 left")
			break
		}
		top, ok := paramsMap["top"]
		if !ok {
			fmt.Println("未设置参数 top")
			break
		}
		windowState, ok := paramsMap["top"].(string)
		if !ok {
			fmt.Println("未设置参数 windowState")
			break
		}
		windowIdInt := gt.Any2Int(windowId)
		widthInt := gt.Any2Int(width)
		heightInt := gt.Any2Int(height)
		leftInt := gt.Any2Int(left)
		topInt := gt.Any2Int(top)
		browser.CDPBrowserSetWindowBounds(windowIdInt, leftInt, topInt, widthInt, heightInt, windowState)

	case "Browser.setContentsSize":
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
		windowIdInt := gt.Any2Int(windowId)
		widthInt := gt.Any2Int(width)
		heightInt := gt.Any2Int(height)
		browser.CDPBrowserSetContentsSize(windowIdInt, widthInt, heightInt)

	case "Target.createTarget":
		urlStr, ok := paramsMap["url"].(string)
		if !ok {
			fmt.Println("未设置参数 url")
			break
		}
		tid, err := browser.CDPTargetCreateTarget(urlStr)
		if err != nil {
			fmt.Println("打开页面失败,err:", err)
			break
		}
		interp.Global().SetVar(to, tid)

	case "Target.activateTarget":
		targetId, ok := paramsMap["targetId"].(string)
		if !ok {
			fmt.Println("未设置参数 targetId")
			break
		}
		browser.CDPTargetActivateTarget(targetId)

	case "Target.attachToTarget":
		targetId, ok := paramsMap["targetId"].(string)
		if !ok {
			fmt.Println("未设置参数 targetId")
			break
		}
		sid, err := browser.CDPTargetAttachToTarget(targetId)
		if err != nil {
			fmt.Println("聚焦目标页签失败,err:", err)
			break
		}
		interp.Global().SetVar(to, sid)

	case "Target.closeTarget":
		targetId, ok := paramsMap["targetId"].(string)
		if !ok {
			fmt.Println("未设置参数 targetId")
			break
		}
		browser.CDPTargetCloseTarget(targetId)

	case "Target.createBrowserContext":
		browser.CDPTargetCreateBrowserContext()

	case "Target.detachFromTarget":
		sessionId, ok := paramsMap["sessionId"].(string)
		if !ok {
			fmt.Println("未设置参数 sessionId")
			break
		}
		browser.CDPTargetDetachFromTarget(sessionId)

	case "Target.disposeBrowserContext":
		browserContextId, ok := paramsMap["browserContextId"].(string)
		if !ok {
			fmt.Println("未设置参数 browserContextId")
			break
		}
		browser.CDPTargetDisposeBrowserContext(browserContextId)

	case "Target.getBrowserContexts":
		browser.CDPTargetGetBrowserContexts()

	case "Target.getTargets":
		browser.CDPTargetGetTargets()

	case "Target.getTargetInfo":
		targetId, ok := paramsMap["targetId"].(string)
		if !ok {
			fmt.Println("未设置参数 targetId")
			break
		}
		browser.CDPTargetGetTargetInfo(targetId)

	case "DOMSnapshot.captureSnapshot":
		browser.CDPDOMSnapshotCaptureSnapshot()

	case "DOMSnapshot.disable":
		browser.CDPDOMSnapshotDisable()

	case "DOMSnapshot.enable":
		browser.CDPDOMSnapshotEnable()

	case "DOMStorage.clear":
		securityOrigin, ok := paramsMap["sessionId"].(string)
		if !ok {
			fmt.Println("未设置参数 securityOrigin")
			break
		}
		isLocalStorage, ok := paramsMap["isLocalStorage"].(bool)
		if !ok {
			fmt.Println("未设置参数 isLocalStorage")
			break
		}
		browser.CDPDOMStorageClear(browser.DOMStorageId{
			SecurityOrigin: securityOrigin,
			IsLocalStorage: isLocalStorage,
		})

	case "DOMStorage.disable":
		browser.CDPDOMStorageDisable()

	case "DOMStorage.enable":
		browser.CDPDOMStorageEnable()

	case "DOMStorage.getDOMStorageItems":
		securityOrigin, ok := paramsMap["sessionId"].(string)
		if !ok {
			fmt.Println("未设置参数 securityOrigin")
			break
		}
		isLocalStorage, ok := paramsMap["isLocalStorage"].(bool)
		if !ok {
			fmt.Println("未设置参数 isLocalStorage")
			break
		}
		browser.CDPDOMStorageGetDOMStorageItems(browser.DOMStorageId{
			SecurityOrigin: securityOrigin,
			IsLocalStorage: isLocalStorage,
		})

	case "DOMStorage.removeDOMStorageItem":
		securityOrigin, ok := paramsMap["sessionId"].(string)
		if !ok {
			fmt.Println("未设置参数 securityOrigin")
			break
		}
		isLocalStorage, ok := paramsMap["isLocalStorage"].(bool)
		if !ok {
			fmt.Println("未设置参数 isLocalStorage")
			break
		}
		key, ok := paramsMap["key"].(string)
		if !ok {
			fmt.Println("未设置参数 key")
			break
		}
		browser.CDPDOMStorageRemoveDOMStorageItem(browser.DOMStorageId{
			SecurityOrigin: securityOrigin,
			IsLocalStorage: isLocalStorage,
		}, key)

	case "DOMStorage.setDOMStorageItem":
		securityOrigin, ok := paramsMap["sessionId"].(string)
		if !ok {
			fmt.Println("未设置参数 securityOrigin")
			break
		}
		isLocalStorage, ok := paramsMap["isLocalStorage"].(bool)
		if !ok {
			fmt.Println("未设置参数 isLocalStorage")
			break
		}
		key, ok := paramsMap["key"].(string)
		if !ok {
			fmt.Println("未设置参数 key")
			break
		}
		value, ok := paramsMap["value"].(string)
		if !ok {
			fmt.Println("未设置参数 value")
			break
		}
		browser.CDPDOMStorageSetDOMStorageItem(browser.DOMStorageId{
			SecurityOrigin: securityOrigin,
			IsLocalStorage: isLocalStorage,
		}, key, value)

	case "CSS.addRule":
		_, err := browser.CDPCSSAddRule(params)
		if err != nil {
			fmt.Println("添加CSS规则失败,err:", err)
		}

	case "CSS.collectClassNames":
		rse, err := browser.CDPCSSColectClassNames(paramsMap["styleSheetId"].(string))
		if err != nil {
			fmt.Println("从指定样式表中收集所有类名执行失败,err:", err)
			break
		}
		fmt.Println("收集的类名 : ", rse)

	case "CSS.enable":
		_, err := browser.CDPCSSEnable()
		if err != nil {
			fmt.Println("启动css域失败,err:", err)
			break
		}

	case "CSS.disable":
		_, err := browser.CDPCSSDisable()
		if err != nil {
			fmt.Println("禁用css域失败,err:", err)
			break
		}

	case "CSS.createStyleSheet":
		_, err := browser.CDPCSSCreateStyleSheet(params)
		if err != nil {
			fmt.Println("创建一个新的样式表失败,err:", err)
			break
		}

	case "CSS.forcePseudoState":
		nodeId := paramsMap["nodeId"].(int)
		forcedPseudoClasses := paramsMap["forcedPseudoClasses"].([]string)
		_, err := browser.CDPCSSForcePseudoState(nodeId, forcedPseudoClasses)
		if err != nil {
			fmt.Println("强制设置伪状态失败,err:", err)
			break
		}

	case "CSS.forceStartingStyle":
		nodeId := paramsMap["nodeId"].(int)
		forced := paramsMap["forced"].(bool)
		_, err := browser.CDPCSSForceStartingStyle(nodeId, forced)
		if err != nil {
			fmt.Println("强制设置起始样式失败,err:", err)
			break
		}

	case "CSS.getBackgroundColors":
		browser.CDPCSSGetBackgroundColors(paramsMap["nodeId"].(int))

	case "CSS.getComputedStyleForNode":
		browser.CDPCSSGetComputedStyleForNode(paramsMap["nodeId"].(int))

	case "CSS.getInlineStylesForNode":
		browser.CDPCSSGetInlineStylesForNode(paramsMap["nodeId"].(int))

	case "CSS.getMatchedStylesForNode":
		browser.CDPCSSGetMatchedStylesForNode(paramsMap["nodeId"].(int))

	case "CSS.getMediaQueries":
		browser.CDPCSSGetMediaQueries()

	case "CSS.getPlatformFontsForNode":
		browser.CDPCSSGetPlatformFontsForNode(paramsMap["nodeId"].(int))

	case "CSS.getStyleSheetText":
		browser.CDPCSSGetStyleSheetText(paramsMap["styleSheetId"].(string))

	case "CSS.setEffectivePropertyValueForNode":
		browser.CDPCSSSetEffectivePropertyValueForNode(params)

	case "CSS.setKeyframeKey":
		browser.CDPCSSSetKeyframeKey(params)

	case "CSS.setMediaText":
		browser.CDPCSSSetMediaText(params)

	case "CSS.setPropertyRulePropertyName":
		browser.CDPCSSSetPropertyRulePropertyName(params)

	case "CSS.setRuleSelector":
		browser.CDPCSSSetRuleSelector(params)

	case "CSS.setStyleSheetText":
		browser.CDPCSSSetStyleSheetText(paramsMap["styleSheetId"].(string), paramsMap["text"].(string))

	case "CSS.setStyleTexts":
		browser.CDPCSSSetStyleTexts(params)

	case "CSS.startRuleUsageTracking":
		browser.CDPCSSStartRuleUsageTracking()

	case "CSS.stopRuleUsageTracking":
		browser.CDPCSSStopRuleUsageTracking()

	case "CSS.takeCoverageDelta":
		browser.CDPCSSTakeCoverageDelta()

	case "CSS.getEnvironmentVariables":
		browser.CDPCSSGetEnvironmentVariables()

	case "CSS.setContainerQueryText":
		browser.CDPCSSSetContainerQueryText(params)

	}
}
