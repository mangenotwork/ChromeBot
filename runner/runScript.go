package runner

import (
	"ChromeBot/dsl/builtins"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/dsl/lexer"
	"ChromeBot/dsl/parser"
	"ChromeBot/global"
	"ChromeBot/utils"
	"fmt"
)

func runScript(source string) {

	source = utils.ProcessCommandLine(source)
	source = utils.EscapeQuotesInBackticks(source)
	source = globalAnalysisScript(source)

	builtins.ChromeWait = 2

	// 词法分析
	l := lexer.New(source)

	// 语法分析
	p := parser.New(l)
	program := p.ParseProgram()

	errs := p.CleanErrors()
	if len(errs) > 0 {
		fmt.Println("解析错误:")
		for _, err := range errs {
			fmt.Println("  " + err)
		}
		return
	}

	// 创建解释器
	interp := interpreter.NewInterpreter()

	// 注册内置函数
	builtins.RegisterBuiltins(interp)

	// 执行程序

	if global.IsRegisterCron {
		fmt.Println("开启了定时任务 ", global.CronPerformance.Arg, " -> ", global.CronToChinese(global.CronPerformance.Arg))
	}

	result, err := interp.Interpret(program)
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		return
	}

	if result != nil {
		fmt.Printf("程序返回值: %v\n", result)
	}

}
