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
var chromeLock = utils.NewTimeoutLock(2 * time.Second)

var ChromeWait = 0

/*
前置说明: 当前的设计一个ChromeBot进程对应一个chrome子进程, 一行命令只支持一个操作

参数说明
init : 初始化打开浏览器，如果已经打开后续语句再出现init会忽略
close : 关闭浏览器
size : 设置浏览器窗口大小与init参数一起用,值为: 宽*高 （900*600） <值类型是字符串>
proxy : 设置浏览器代理与init参数一起用 <值类型是字符串>
userpath : 设置浏览器在本机的隔离目录与init参数一起用,对应浏览器的--user-data-dir，建议隔离 <值类型是字符串>
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

		if ChromeWait > 0 { // todo 以后更好的方案 ，  解决执行脚本命令之间操作太快
			time.Sleep(time.Duration(ChromeWait) * time.Second)
		}

		chromeLock.Lock()
		defer chromeLock.Unlock()

		argMap := make(map[string]string)
		for i, v := range args {
			utils.Debugf("参数 %d %v %T\n", i, v, v)
			switch v.(type) {
			case string:
				vList := strings.SplitN(v.(string), "=", 2)
				utils.Debug("vList = ", vList, len(vList))
				if len(vList) == 1 {
					argMap[vList[0]] = ""
				} else if len(vList) == 2 {
					argMap[vList[0]] = vList[1]
				}
			}

		}

		utils.Debug("argMap:", argMap)

		op := &chromeOperation{
			arg: make(map[string]interpreter.Value),
		}
		opNumber := 0

		if _, ok := argMap["init"]; ok {
			op = &chromeOperation{
				opType: opInit,
				//level:  executionPriority[opInit],
			}
			opNumber++
		}

		if _, ok := argMap["close"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opClose,
				//level:  executionPriority[opClose],
			}
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

		if val, ok := argMap["tab"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opTable,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opTable],
			}
			opNumber++
		}

		if val, ok := argMap["req"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opReq,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opReq],
			}
			opNumber++
		}

		if val, ok := argMap["click"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opClick,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opClick],
			}
			opNumber++
		}

		if val, ok := argMap["input"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opInput,
				arg:    map[string]interpreter.Value{"input": val},
				//level:  executionPriority[opInput],
			}
			opNumber++
		}

		if val, ok := argMap["xpath"]; ok {
			if op != nil && op.opType == opInput {
				op.arg["xpath"] = val
			}
		}

		if val, ok := argMap["check"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opCheck,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opCheck],
			}
			opNumber++
		}

		if val, ok := argMap["wait"]; ok && opNumber == 0 {
			op.arg["wait"] = val
		}

		if val, ok := argMap["pause"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opPause,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opScroll],
				extendType: 0,
			}
			opNumber++
		}

		if val, ok := argMap["scroll"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opScroll,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opScroll],
				extendType: 0,
			}
			opNumber++
		}

		if val, ok := argMap["scrollpixel"]; ok && opNumber == 0 {
			valList := strings.Split(val, ",")
			if len(valList) != 2 {
				fmt.Println("[Chrome] scrollpixel 参数错误，值为(x,y)如(2000, 500)")
			}
			op = &chromeOperation{
				opType: opScroll,
				arg:    map[string]interpreter.Value{"x": valList[0], "y": valList[1]},
				//level:  executionPriority[opScroll],
				extendType: 1,
			}
			opNumber++
		}

		if val, ok := argMap["scrollxpath"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opScroll,
				arg:    map[string]interpreter.Value{"xpath": val},
				//level:  executionPriority[opScroll],
				extendType: 2,
			}
			opNumber++
		}

		if val, ok := argMap["screenshot"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opScreenshot,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opScreenshot],
			}
			opNumber++
		}

		if val, ok := argMap["to"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opTo,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opTo],
			}
			opNumber++
		}

		if val, ok := argMap["save"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opSave,
				arg:    map[string]interpreter.Value{"arg": val},
				//level:  executionPriority[opSave],
			}
			opNumber++
		}

		utils.Debugf("执行 : %v", op)

		switch op.opType {
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
			err := browser.ChromeInit(windowSize, proxy, userPath)
			if err != nil {
				fmt.Println("[ERR]", err.Error())
			}

		case opClose:
			fmt.Println("[Chrome]关闭浏览器...")
			chromeObj := browser.GetChromeInstance()
			err := chromeObj.Close()
			if err != nil {
				fmt.Println("[ERR]", err.Error())
			}

		case opTable:
			fmt.Println("tab操作...")
			arg := op.arg["arg"].(string)
			if arg == "list" { // 获取所有tab
				chromeObj := browser.GetChromeInstance()
				_, err := chromeObj.GetAllTab()
				if err != nil {
					log.Println("[Chrome]获取tab错误: ", err.Error())
				}
				break
			}

			if arg == "new" { // 新建一个tab
				chromeObj := browser.GetChromeInstance()
				_, err := chromeObj.NewTab()
				if err != nil {
					log.Println("[Chrome]创建tab错误: ", err.Error())
				}
				break
			}

			if arg == "now" { // 当前tab
				chromeObj := browser.GetChromeInstance()
				chromeObj.NowTabInfo()
				break
			}

			if arg == "close" { // 关闭当前tab
				chromeObj := browser.GetChromeInstance()
				chromeObj.NowTabClose()
				break
			}

			chromeObj := browser.GetChromeInstance()
			chromeObj.SelectTab(arg)

		case opReq:
			fmt.Println("请求操作...")
			chromeObj := browser.GetChromeInstance()
			res, err := chromeObj.OpenUrl(op.arg["arg"].(string))
			if err != nil {
				fmt.Println("[Chrome]请求操作出现错误:", err.Error())
			}
			fmt.Println("[Chrome]请求操作结果:", res)

		case opClick:
			fmt.Println("点击操作...")

			if wait, waitOK := op.arg["wait"]; waitOK {
				waitInt := gt.Any2Int(wait)
				if waitInt > 0 {
					time.Sleep(time.Duration(waitInt) * time.Second)
				}
			}

			xPath := op.arg["arg"].(string)
			// 匹配一下判断arg是不是变量
			xPathVal, xPathValOK := interp.Global().GetVar(xPath)
			if xPathValOK {
				xPath = xPathVal.(string)
			}
			fmt.Println("[Chrome]点击的Xpath = ", xPath)

			chromeObj := browser.GetChromeInstance()
			err := chromeObj.Click(xPath)
			if err != nil {
				fmt.Println("[Chrome]点击操作出现错误:", err.Error())
			}

		case opInput:
			fmt.Println("输入操作...")

			if wait, waitOK := op.arg["wait"]; waitOK {
				waitInt := gt.Any2Int(wait)
				if waitInt > 0 {
					time.Sleep(time.Duration(waitInt) * time.Second)
				}
			}

			xPath := op.arg["xpath"].(string)
			// 匹配一下判断arg是不是变量
			xPathVal, xPathValOK := interp.Global().GetVar(xPath)
			if xPathValOK {
				xPath = xPathVal.(string)
			}
			fmt.Println("[Chrome]输入的Xpath = ", xPath)

			inputText := op.arg["input"].(string)
			chromeObj := browser.GetChromeInstance()
			err := chromeObj.Input(xPath, inputText)
			if err != nil {
				fmt.Println("[Chrome]输入操作出现错误:", err.Error())
			}

		case opCheck:
			fmt.Println("检查操作...")
			inputText := op.arg["arg"].(string)
			chromeObj := browser.GetChromeInstance()
			has, err := chromeObj.Check(inputText)
			if err != nil {
				fmt.Println("[Chrome]检查操作出现错误:", err.Error())
			}
			fmt.Printf("[Chrome]检查操作xPath: %s , %v", inputText, has)

		case opPause:
			fmt.Println("等待操作...")
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
			fmt.Println("滚动操作...")
			var err error
			chromeObj := browser.GetChromeInstance()
			switch op.extendType {
			case 0:
				high := gt.Any2Int(op.arg["arg"])
				log.Println("滚动的高度 high = ", high)
				err = chromeObj.ScrollByPixel(0, high)

			case 1:
				x := gt.Any2Int(op.arg["x"])
				y := gt.Any2Int(op.arg["y"])
				log.Println("滚动的高度 x = ", x, ", y = ", y)
				err = chromeObj.ScrollByPixel(gt.Any2Int(op.arg["x"]), gt.Any2Int(op.arg["y"]))

			case 2:
				xpath := op.arg["xpath"].(string)
				log.Println("滚动到Xpath = ", xpath)
				err = chromeObj.ScrollToElement(xpath)
			}

			if err != nil {
				fmt.Println("[Chrome]滚动操作出现错误:", err.Error())
			}

		case opScreenshot:
			fmt.Println("截图操作...")
			savePath := op.arg["arg"].(string)
			chromeObj := browser.GetChromeInstance()
			res, err := chromeObj.CaptureFullPageScreenshot(savePath)
			if err != nil {
				log.Println("[Chrome]截图操作错误: ", err.Error())
				return nil, fmt.Errorf("[Chrome]截图操作错误: %s", err.Error())
			}
			log.Printf("[Chrome]截图结果: %v", res)

		case opTo:
			fmt.Println("将当前页面的html赋值到变量操作...")
			chromeObj := browser.GetChromeInstance()
			htmlBody, err := chromeObj.GetHtml()
			if err != nil {
				log.Println("[Chrome]获取页面的html错误: ", err.Error())
				return nil, fmt.Errorf("[Chrome]获取页面的html错误: %s", err.Error())
			}
			to := op.arg["arg"].(string)
			interp.Global().SetVar(to, htmlBody)

		case opSave:
			fmt.Println("将当前页面的html保存到本地...")

			if wait, waitOK := op.arg["wait"]; waitOK {
				waitInt := gt.Any2Int(wait)
				if waitInt > 0 {
					time.Sleep(time.Duration(waitInt) * time.Second)
				}
			}

			chromeObj := browser.GetChromeInstance()
			htmlBody, err := chromeObj.GetHtml()
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
