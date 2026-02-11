package runner

import (
	"ChromeBot/dsl/builtins"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/dsl/lexer"
	"ChromeBot/dsl/parser"
	"fmt"
)

func runScript(source string) {
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
	result, err := interp.Interpret(program)
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		return
	}

	if result != nil {
		fmt.Printf("程序返回值: %v\n", result)
	}
}
