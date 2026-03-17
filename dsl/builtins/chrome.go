package builtins

import (
	"ChromeBot/browser"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/utils"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"strings"
	"time"
)

// 一个带timeout的锁
var chromeLock = utils.NewTimeoutLock(60 * time.Second) // 默认1分钟
var ChromeWait = 0

var chromeSupport = map[string]bool{
	"init":        true,
	"close":       true,
	"size":        true,
	"proxy":       true,
	"userpath":    true,
	"new":         true,
	"tab":         true,
	"req":         true,
	"click":       true,
	"xpath":       true,
	"input":       true,
	"check":       true,
	"wait":        true,
	"pause":       true,
	"scroll":      true,
	"scrollpixel": true,
	"scrollxpath": true,
	"screenshot":  true,
	"to":          true,
	"save":        true,
	"info":        true,
	"as":          true,
}

func hasChromeSupport(cmd string) bool {
	_, ok := chromeSupport[cmd]
	return ok
}

/*
前置说明: 当前的设计一个ChromeBot进程对应一个chrome子进程, 一行命令只支持一个操作

参数说明
init : 初始化打开浏览器，如果已经打开后续语句再出现init会忽略
close : 关闭浏览器
size : 设置浏览器窗口大小与init参数一起用,值为: 宽*高 （900*600） <值类型是字符串>
proxy : 设置浏览器代理与init参数一起用 <值类型是字符串>
userpath : 设置浏览器在本机的隔离目录与init参数一起用,对应浏览器的--user-data-dir，建议隔离 <值类型是字符串>
new : 设置浏览器新建一个隔离环境与init参数一起用；与userPath同时在时，优先使用userPath
tab : 页签, 值有get:获取；set:指定哪个标签切换到指定的页签; new：新建一个页签；1<number>:第一个页签；select：返回当前选中的页签; 注意: 如果是没有选中页签下文操作默认当前浏览器的页签进行操作 <值类型是指定的字符串>
req :  请求网址， 值为网址 <值类型是字符串>
（ dom : 获取当前页面html的dom树 - 改为函数 ）
click : 点击操作，值为xpath <值类型是字符串>
xpath : 当前选中的xpath, 输入的时候用
input : 输入操作，输入内容  <值类型是字符串>
check : 检查操作，检查页面是否存在指定xpath  <值类型是字符串>
wait : 默认会执行等待页面加载完成，这个参数给定操作时候设置等待的时间  <值类型是数值类型>
pause : 默认会执行等待页面加载完成，这个参数给定操作时候设置等待的时间  <值类型是数值类型>
scroll : 滚动操作，滚动页面  正数往下，负数往上 <值类型是数值类型>  注意: 该滚动存在局限性只针对根节点进行滚动，嵌套容器要想精确请使用 scrollxpath
scrollpixel : scroll by pixel 滚动操作,滚动到指定坐标， 值为(x,y)如(2000, 500)   注意: 该滚动存在局限性只针对根节点进行滚动, 嵌套容器要想精确请使用 scrollxpath
scrollxpath : 滚动操作,滚动到指定xpath <值类型是字符串>
screenshot : 截图操作，浏览器截图操作  值为保存位置  <值类型是字符串>
to : 将当前操作的页面html返回存入到指定变量-如果变量未声明这里会自动声明变量  <值类型是字符串>
save : 将将当前操作的页面html存入到指定文件  <值类型是字符串>
info : 获取chrome 的信息
as : 将指令的结果赋值给变量

语法

chrome click=`//*[@id="chat-submit-button"]` xpath=`//*[@id="chat-textarea"]` input=`aaaa`

例子1 ：简单访问百度进行查询操作
chrome init userpath="D:/chromeTest/"  // 打开浏览器
chrome req="www.baidu.com" // 访问
chrome inputxp=(`//*[@id="chat-textarea"]`,"mange") // 输入
chrome click=`//*[@id="chat-submit-button"]` // 点击确定
chrome to=res // 默认等待页面变换 将当前页面的html存储到变量 res
print(res) // 输出
chrome close  // 关闭浏览器
*/
func registerChrome(interp *interpreter.Interpreter) {

	interp.Global().SetFunc("chrome", func(args []interpreter.Value) (interpreter.Value, error) {
		utils.Debug("执行 chrome 的操作，参数是 ", args, len(args))

		argsStr := make([]string, 0)
		for i, arg := range args {
			utils.Debugf("参数 %d %v %T\n", i, arg, arg)
			argsStr = append(argsStr, arg.(string))
		}

		// 处理 函数类型的参数
		argsStr = processArgs(interp, argsStr)
		utils.Debug("执行 ProcessArgs 参数 处理  ", argsStr, len(args))

		if len(argsStr) == 0 {
			return nil, fmt.Errorf("未知命令，请参考文档")
		}

		chromeLock.Lock()
		defer chromeLock.Unlock()

		argMap := make(map[string]string)
		for _, v := range argsStr {
			vList := strings.SplitN(v, "=", 2)
			utils.Debug("vList = ", vList, len(vList))
			if !hasChromeSupport(vList[0]) {
				fmt.Println("[Chrome]未知命令 ", vList[0], ";请参考文档。")
				return nil, fmt.Errorf("[Chrome]未知命令 %s;请参考文档。", vList[0])
			}
			if len(vList) == 1 {
				argMap[vList[0]] = ""
			} else if len(vList) == 2 {
				argMap[vList[0]] = vList[1]
			}
		}

		utils.Debug("argMap:", argMap)

		var wait = 0

		op := &chromeOperation{
			arg: make(map[string]interpreter.Value),
		}
		opNumber := 0

		if _, ok := argMap["init"]; ok {
			op.opType = opInit
			opNumber++
		}

		if _, ok := argMap["info"]; ok {
			op.opType = opInfo
			opNumber++
		}

		if _, ok := argMap["close"]; ok && opNumber == 0 {
			op.opType = opClose
			opNumber++
		}

		if val, ok := argMap["size"]; ok {
			if op != nil && op.opType == opInit {
				op.arg["size"] = val
			}
		}

		if val, ok := argMap["proxy"]; ok {
			if op != nil && op.opType == opInit {
				op.arg["proxy"] = val
			}
		}

		if val, ok := argMap["userpath"]; ok {
			if op != nil && op.opType == opInit {
				op.arg["userpath"] = val
			}
		}

		if _, ok := argMap["new"]; ok {
			if op != nil && op.opType == opInit {
				op.arg["new"] = 1
			}
		}

		if val, ok := argMap["tab"]; ok && opNumber == 0 {
			op.opType = opTable
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["req"]; ok && opNumber == 0 {
			op.opType = opReq
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["click"]; ok && opNumber == 0 {
			op.opType = opClick
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["input"]; ok && opNumber == 0 {
			op.opType = opInput
			op.arg["input"] = val
			opNumber++
		}

		if val, ok := argMap["xpath"]; ok {
			if op != nil && op.opType == opInput {
				op.arg["xpath"] = val
			}
		}

		if val, ok := argMap["check"]; ok && opNumber == 0 {
			op.opType = opCheck
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["wait"]; ok && opNumber == 0 {
			wait = gt.Any2Int(val)
		}

		if val, ok := argMap["pause"]; ok && opNumber == 0 {
			op.opType = opPause
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["scroll"]; ok && opNumber == 0 {
			op.opType = opScroll
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["scrollpixel"]; ok && opNumber == 0 {
			valList := strings.Split(val, ",")
			if len(valList) != 2 {
				fmt.Println("[Chrome] scrollpixel 参数错误，值为(x,y)如(2000, 500)")
			}
			op.opType = opScroll
			op.arg["x"] = valList[0]
			op.arg["y"] = valList[1]
			op.extendType = 1
			opNumber++
		}

		if val, ok := argMap["scrollxpath"]; ok && opNumber == 0 {
			op.opType = opScroll
			op.arg["xpath"] = val
			op.extendType = 2
			opNumber++
		}

		if val, ok := argMap["screenshot"]; ok && opNumber == 0 {
			op.opType = opScreenshot
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["to"]; ok && opNumber == 0 {
			op.opType = opTo
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["save"]; ok && opNumber == 0 {
			op.opType = opSave
			op.arg["arg"] = val
			opNumber++
		}

		if val, ok := argMap["as"]; ok && opNumber == 1 {
			op.arg["as"] = val
		}

		utils.Debugf("执行 : %v", op)
		if ChromeWait > 0 { // todo 以后更好的方案 ，  解决执行脚本命令之间操作太快
			time.Sleep(time.Duration(ChromeWait) * time.Second)
		}

		if wait > 0 {
			fmt.Printf("[Chrome] wait %d s", wait)
			time.Sleep(time.Duration(wait) * time.Second)
		}

		switch op.opType {

		case opInfo:
			chromePath, err := browser.FindChrome()
			if err != nil {
				fmt.Println("获取chrome可执行文件路径失败")
			}
			fmt.Println("[Chrom] 路径 : ", chromePath)
			info, err := browser.GetChromeInfo(chromePath)
			if err != nil {
				fmt.Printf("获取Chrome信息失败: %v\n", err)
			} else {
				fmt.Println("Chrome 浏览器信息：")
				for k, v := range info {
					fmt.Printf("%-20s: %s\n", k, v)
				}
			}
			if asArg, ok := op.arg["as"]; ok {
				rse := make(interpreter.DictType)
				for k, v := range info {
					rse[k] = v
				}
				interp.Global().SetVar(asArg.(string), rse)
			}

		case opInit:
			fmt.Println("[Chrome]初始化浏览器...")
			windowSize, proxy, userPath := "", "", ""
			if val, ok := op.arg["size"]; ok {
				windowSize = val.(string)
			}
			if val, ok := op.arg["proxy"]; ok {
				proxy = val.(string)
			}
			if val, ok := op.arg["userpath"]; ok {
				userPath = val.(string)
			}
			isNew := false
			if _, ok := op.arg["new"]; ok {
				isNew = true
			}
			browser.ChromeInit(windowSize, proxy, userPath, isNew)

		case opClose:
			fmt.Println("[Chrome]关闭浏览器...")
			err := browser.Close()
			if err != nil {
				fmt.Println("[ERR]", err.Error())
			}

		case opTable:
			fmt.Println("[Chrome]tab操作...")
			arg := op.arg["arg"].(string)
			if arg == "list" { // 获取所有tab
				_, err := browser.GetAllTab()
				if err != nil {
					log.Println("[Chrome]获取tab错误: ", err.Error())
				}
				break
			}
			if arg == "new" { // 新建一个tab
				_, err := browser.NewTab()
				if err != nil {
					log.Println("[Chrome]创建tab错误: ", err.Error())
				}
				break
			}
			if arg == "now" { // 当前tab
				browser.NowTabInfo()
				break
			}
			if arg == "close" { // 关闭当前tab
				browser.NowTabClose()
				break
			}
			log.Println("arg = ", arg)
			browser.SelectTab(arg)

		case opReq:
			fmt.Println("[Chrome]请求操作...")
			reqUrl := op.arg["arg"].(string)
			utils.Debug("reqUrl = ", reqUrl)
			// 匹配一下判断arg是不是变量
			reqUrlVal, reqUrlValOK := interp.Global().GetVar(reqUrl)
			if reqUrlValOK {
				utils.Debug("存在变量 ", reqUrl, " | 值: ", reqUrlVal)
				reqUrl = reqUrlVal.(string)
			}
			fmt.Println("[Chrome]请求 url = ", reqUrl)
			rse, err := browser.OpenUrl(reqUrl)
			if err != nil {
				fmt.Println("[Chrome]请求操作出现错误:", err.Error())
			}
			if asArg, ok := op.arg["as"]; ok {
				interp.Global().SetVar(asArg.(string), interpreter.Value(rse))
			}

		case opClick:
			fmt.Println("[Chrome]点击操作...")

			xPath := op.arg["arg"].(string)
			// 匹配一下判断arg是不是变量
			xPathVal, xPathValOK := interp.Global().GetVar(xPath)
			if xPathValOK {
				xPath = xPathVal.(string)
			}

			_, err := utils.ValidateXPathPureNative(xPath)
			if err != nil {
				fmt.Println("[Chrome]点击操作警告: ", err.Error())
				break
			}

			fmt.Println("[Chrome]点击的Xpath = ", xPath)
			err = browser.Click(xPath)
			if err != nil {
				fmt.Println("[Chrome]点击操作出现错误:", err.Error())
			}

		case opInput:
			fmt.Println("[Chrome]输入操作...")

			xPath, ok := op.arg["xpath"].(string)
			if !ok {
				xPath = ""
			}
			// 匹配一下判断arg是不是变量
			xPathVal, xPathValOK := interp.Global().GetVar(xPath)
			if xPathValOK {
				xPath = xPathVal.(string)
			}
			fmt.Println("[Chrome]输入的Xpath = ", xPath)

			if xPath == "" {
				fmt.Println("[Chrome]输入操作警告: 未设置Xpath无法执行操作")
				break
			}

			_, err := utils.ValidateXPathPureNative(xPath)
			if err != nil {
				fmt.Println("[Chrome]输入操作警告: ", err.Error())
				break
			}

			inputText := op.arg["input"].(string)
			inputTextVal, inputTextValOK := interp.Global().GetVar(inputText)
			if inputTextValOK {
				inputText = inputTextVal.(string)
			}
			fmt.Println("[Chrome]输入内容 = ", inputText)

			err = browser.Input(xPath, inputText)
			if err != nil {
				fmt.Println("[Chrome]输入操作出现错误:", err.Error())
			}

		case opCheck:
			fmt.Println("[Chrome]检查操作...")
			inputText := op.arg["arg"].(string)
			has, err := browser.Check(inputText)
			if err != nil {
				fmt.Println("[Chrome]检查操作出现错误:", err.Error())
			}
			fmt.Printf("[Chrome]检查操作xPath: %s , %v", inputText, has)
			if asArg, ok := op.arg["as"]; ok {
				interp.Global().SetVar(asArg.(string), interpreter.Value(has))
			}

		case opPause:
			fmt.Println("[Chrome]等待操作...")
			if pause, pauseOK := op.arg["arg"]; pauseOK {
				pauseInt := gt.Any2Int(pause)
				if pauseInt > 0 {
					for i := 0; i < pauseInt; i++ {
						fmt.Printf("[Chrome] pause: %ds ...\n", pauseInt-i)
						time.Sleep(time.Duration(1) * time.Second)
					}
				}
			}

		case opScroll:
			fmt.Println("[Chrome]滚动操作...")
			var err error

			switch op.extendType {
			case 0:
				high := gt.Any2Int(op.arg["arg"])
				log.Println("滚动的高度 high = ", high)
				err = browser.ScrollByPixel(0, high)

			case 1:
				x := gt.Any2Int(op.arg["x"])
				y := gt.Any2Int(op.arg["y"])
				log.Println("滚动的高度 x = ", x, ", y = ", y)
				err = browser.ScrollByPixel(gt.Any2Int(op.arg["x"]), gt.Any2Int(op.arg["y"]))

			case 2:
				xpath := op.arg["xpath"].(string)
				_, err = utils.ValidateXPathPureNative(xpath)
				if err != nil {
					fmt.Println("[Chrome]滚动操作警告: ", err.Error())
					break
				}
				log.Println("滚动到Xpath = ", xpath)
				err = browser.ScrollToElement(xpath)
			}

			if err != nil {
				fmt.Println("[Chrome]滚动操作出现错误:", err.Error())
			}

		case opScreenshot:
			fmt.Println("[Chrome]截图操作...")
			savePath := op.arg["arg"].(string)
			// 匹配一下判断arg是不是变量
			savePathVal, savePathValOK := interp.Global().GetVar(savePath)
			if savePathValOK {
				savePath = savePathVal.(string)
			}
			fmt.Println("[Chrome]存储的 path = ", savePath)

			res, err := browser.CaptureFullPageScreenshot(savePath)
			if err != nil {
				log.Println("[Chrome]截图操作错误: ", err.Error())
				return nil, fmt.Errorf("[Chrome]截图操作错误: %s", err.Error())
			}
			utils.Debugf("[Chrome]截图结果: %v", res)

		case opTo:
			fmt.Println("[Chrome]将当前页面的html赋值到变量操作...")
			htmlBody, err := browser.GetHtml()
			if err != nil {
				log.Println("[Chrome]获取页面的html错误: ", err.Error())
				return nil, fmt.Errorf("[Chrome]获取页面的html错误: %s", err.Error())
			}
			to := op.arg["arg"].(string)
			interp.Global().SetVar(to, htmlBody)

		case opSave:
			fmt.Println("[Chrome]将当前页面的html保存到本地...")

			htmlBody, err := browser.GetHtml()
			if err != nil {
				log.Println("[Chrome]获取页面的html错误: ", err.Error())
				return nil, fmt.Errorf("[Chrome]获取页面的html错误: %s", err.Error())
			}
			savePath := op.arg["arg"].(string)
			err = utils.SaveDataToFile(savePath, htmlBody)
			if err != nil {
				fmt.Println("保存页面到文件出现了错误:", err.Error())
			}

		}

		return nil, nil
	})
}

type chromeOPType string

var (
	opInit       chromeOPType = "init"  // 初始化浏览器
	opInfo       chromeOPType = "info"  // 获取chrome info
	opClose      chromeOPType = "close" // 关闭浏览器
	opTable      chromeOPType = "tab"
	opReq        chromeOPType = "req"
	opClick      chromeOPType = "click"      // 点击操作
	opInput      chromeOPType = "input"      // 输入操作
	opCheck      chromeOPType = "check"      // 检查操作
	opPause      chromeOPType = "pause"      // 等待操作
	opScroll     chromeOPType = "scroll"     // 滚动操作
	opScreenshot chromeOPType = "screenshot" // 截图操作
	opTo         chromeOPType = "to"         // 将当前页面的html赋值到变量操作
	opSave       chromeOPType = "save"       // 将当前页面的html保存到本地
)

//type executionPriorityLevel int
//
//var (
//	executionPriorityLevel1 executionPriorityLevel = 1 // 优先级1：浏览器的直接操作
//	executionPriorityLevel2 executionPriorityLevel = 2 // 优先级2：Table页签操作，等待时间
//	executionPriorityLevel3 executionPriorityLevel = 3 // 优先级3：对Table页签进行输入网址
//	executionPriorityLevel4 executionPriorityLevel = 4 // 优先级4：检查页面是否存在某元素或Xpath
//	executionPriorityLevel5 executionPriorityLevel = 5 // 优先级5：对页面进行直接操作
//	executionPriorityLevel6 executionPriorityLevel = 6 // 优先级6: 将页面进行输出
//)
//
//// 执行优先级 顺序：小的在前
//var executionPriority = map[chromeOPType]executionPriorityLevel{
//	opInit:       executionPriorityLevel1,
//	opClose:      executionPriorityLevel1,
//	opTable:      executionPriorityLevel2,
//	opReq:        executionPriorityLevel3,
//	opClick:      executionPriorityLevel5,
//	opInput:      executionPriorityLevel5,
//	opCheck:      executionPriorityLevel4,
//	opWait:       executionPriorityLevel2,
//	opScroll:     executionPriorityLevel5,
//	opScreenshot: executionPriorityLevel6,
//	opTo:         executionPriorityLevel6,
//	opSave:       executionPriorityLevel6,
//}

type chromeOperation struct {
	opType chromeOPType                 // 操作的类型
	arg    map[string]interpreter.Value // 操作的参数
	//level  executionPriorityLevel // 操作等级
	extendType int // 扩展类型，用于同效果的多类型进行区分
}

func processArgs(interp *interpreter.Interpreter, args []string) []string {
	// 空数组直接返回
	if len(args) == 0 {
		return args
	}

	var result []string
	i := 0
	length := len(args)

	for i < length {
		// 找到左括号的位置
		if args[i] == "(" {
			// 检查左括号前是否有元素（避免数组越界）
			if i == 0 {
				// 左括号开头，直接保留（异常场景）
				result = append(result, args[i])
				i++
				continue
			}

			// 查找对应的右括号
			rightIdx := -1
			for j := i + 1; j < length; j++ {
				if args[j] == ")" {
					rightIdx = j
					break
				}
			}

			if rightIdx == -1 {
				// 未找到右括号，按原格式保留
				result = append(result, args[i-1], args[i])
				i++
				continue
			}

			// 提取括号内的元素并拼接
			innerElements := args[i+1 : rightIdx]
			//innerStr := strings.Join(innerElements, ",")
			prev := args[i-1]

			prevList := strings.SplitN(prev, "=", 2)

			if len(prevList) == 2 {

				if !hasChromeSupport(prevList[0]) {
					fmt.Println("[Chrome]未知命令 ", prevList[0], ";请参考文档。")
					return nil
				}

				funcName := prevList[1]
				utils.Debug("检查到函数 --> ", funcName)
				fn, has := interp.Global().GetFunc(funcName)
				utils.Debug("找到函数结果  --> ", has)
				utils.Debug("找到函数参数 --> ", innerElements, " | len:", len(innerElements))
				fnArgs := make([]interpreter.Value, 0)
				for _, arg := range innerElements {
					val, valHas := interp.Global().GetVar(arg)
					if valHas {
						arg = val.(string)
					}
					fnArgs = append(fnArgs, arg)
				}
				utils.Debug("整理出函数参数 --> ", fnArgs, " | len:", len(fnArgs))
				rse, rseErr := fn(fnArgs)
				if rseErr != nil {
					fmt.Println("[Chrome]命令行中的函数运行出现错误, err = ", rseErr.Error())
					return result
				}
				utils.Debug("函数运行结果 ", rse)
				// 拼接成完整字符串（前元素 + ( + 拼接内容 + )）
				fullStr := fmt.Sprintf("%s=%s", prevList[0], rse)
				utils.Debug("fullStr = ", fullStr)

				result = append(result, fullStr)

				// 跳过已处理的元素
				i = rightIdx + 1
			}

		} else {
			if i+1 < length && args[i+1] == "(" {
				i++
			} else {
				// 单元素/非括号后无左括号 → 直接添加
				result = append(result, args[i])
				i++
			}
		}
	}

	return result
}
