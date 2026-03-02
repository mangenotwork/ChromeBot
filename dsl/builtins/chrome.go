package builtins

import (
	"ChromeBot/dsl/interpreter"
	"log"
	"strings"
)

/*
前置说明: 当前的设计一个ChromeBot进程对应一个chrome子进程, 一行命令只支持一个操作

参数说明
init : 初始化打开浏览器，如果已经打开后续语句再出现init会忽略
close : 关闭浏览器
size : 设置浏览器窗口大小与init参数一起用,值为: 宽*高 （900*600） <值类型是字符串>
proxy : 设置浏览器代理与init参数一起用 <值类型是字符串>
userpath : 设置浏览器在本机的隔离目录与init参数一起用,对应浏览器的--user-data-dir，建议隔离 <值类型是字符串>
table : 页签, 值有get:获取；set:指定哪个标签切换到指定的页签; new：新建一个页签；1<number>:第一个页签；selected：返回当前选中的页签; 注意: 如果是没有选中页签下文操作默认当前浏览器的页签进行操作 <值类型是指定的字符串>
req :  请求网址， 值为网址 <值类型是字符串>
（ dom : 获取当前页面html的dom树 - 改为函数 ）
click : 点击操作，值为xpath <值类型是字符串>
xpath : 当前选中的xpath, 输入的时候用
input : 输入操作，输入内容  <值类型是字符串>
check : 检查操作，检查页面是否存在指定xpath  <值类型是字符串>
wait : 默认会执行等待页面加载完成，这个参数给定操作时候设置等待的时间  <值类型是数值类型>
scroll : 滚动操作，滚动页面  正数往下，负数往上 <值类型是数值类型>
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
		log.Println("执行 chrome 的操作，参数是 ", args, len(args))

		argMap := make(map[string]string)
		for i, v := range args {
			log.Printf("参数 %d %v %T\n", i, v, v)
			switch v.(type) {
			case string:
				vList := strings.SplitN(v.(string), "=", 2)
				log.Println("vList = ", vList, len(vList))
				if len(vList) == 1 {
					argMap[vList[0]] = ""
				} else if len(vList) == 2 {
					argMap[vList[0]] = vList[1]
				}
			}

		}

		log.Println("argMap:", argMap)

		op := &chromeOperation{}
		opNumber := 0

		if _, ok := argMap["init"]; ok {
			op = &chromeOperation{
				opType: opInit,
				arg:    make([]interpreter.Value, 0),
				//level:  executionPriority[opInit],
			}
			opNumber++
		}

		if _, ok := argMap["close"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opClose,
				arg:    make([]interpreter.Value, 0),
				//level:  executionPriority[opClose],
			}
			opNumber++
		}

		if val, ok := argMap["size"]; ok {
			if op != nil && op.opType == opInit {
				op.arg = append(op.arg, "size=", val)
			}
		}

		if val, ok := argMap["proxy"]; ok {
			if op != nil && op.opType == opInit {
				op.arg = append(op.arg, "proxy=", val)
			}
		}

		if val, ok := argMap["userpath"]; ok {
			if op != nil && op.opType == opInit {
				op.arg = append(op.arg, "userpath="+val)
			}
		}

		if val, ok := argMap["table"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opTable,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opTable],
			}
			opNumber++
		}

		if val, ok := argMap["req"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opReq,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opReq],
			}
			opNumber++
		}

		if val, ok := argMap["click"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opClick,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opClick],
			}
			opNumber++
		}

		if val, ok := argMap["input"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opInput,
				arg:    []interpreter.Value{"input=" + val},
				//level:  executionPriority[opInput],
			}
			opNumber++
		}

		if val, ok := argMap["xpath"]; ok {
			if op != nil && op.opType == opInput {
				op.arg = append(op.arg, "xpath="+val)
			}
		}

		if val, ok := argMap["check"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opCheck,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opCheck],
			}
			opNumber++
		}

		if val, ok := argMap["wait"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opWait,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opWait],
			}
			opNumber++
		}

		if val, ok := argMap["scroll"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opScroll,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opScroll],
			}
			opNumber++
		}

		if val, ok := argMap["screenshot"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opScreenshot,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opScreenshot],
			}
			opNumber++
		}

		if val, ok := argMap["to"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opTo,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opTo],
			}
			opNumber++
		}

		if val, ok := argMap["save"]; ok && opNumber == 0 {
			op = &chromeOperation{
				opType: opSave,
				arg:    []interpreter.Value{val},
				//level:  executionPriority[opSave],
			}
			opNumber++
		}

		log.Printf("执行 : %v", op)

		return nil, nil
	})
}

type chromeOPType string

var (
	opInit       chromeOPType = "init"  // 初始化浏览器
	opClose      chromeOPType = "close" // 关闭浏览器
	opTable      chromeOPType = "table"
	opReq        chromeOPType = "req"
	opClick      chromeOPType = "click"      // 点击操作
	opInput      chromeOPType = "input"      // 输入操作
	opCheck      chromeOPType = "check"      // 检查操作
	opWait       chromeOPType = "wait"       // 等待操作
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
	opType chromeOPType        // 操作的类型
	arg    []interpreter.Value // 操作的参数
	//level  executionPriorityLevel // 操作等级
}
