package builtins

import (
	"ChromeBot/dsl/interpreter"
)

// RegisterBuiltins 注册所有内置函数
func RegisterBuiltins(interp *interpreter.Interpreter) {

	// 注册 mathFn
	for name, fn := range mathFn {
		interp.Global().SetFunc(name, fn)
	}

	// 注册 strFn
	for name, fn := range strFn {
		interp.Global().SetFunc(name, fn)
	}

	// 注册 timeFn
	for name, fn := range timeFn {
		interp.Global().SetFunc(name, fn)
	}

	// 注册 chrome
	registerChrome(interp)
}
