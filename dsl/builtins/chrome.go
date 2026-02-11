package builtins

import (
	"ChromeBot/dsl/interpreter"
	"log"
)

func registerChrome(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("chrome", func(args []interpreter.Value) (interpreter.Value, error) {
		log.Println("执行 chrome 的操作，参数是 ", args, len(args))
		return nil, nil
	})
}
