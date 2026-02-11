package builtins

import (
	"ChromeBot/dsl/interpreter"
	"log"
)

func registerHttp(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("http", func(args []interpreter.Value) (interpreter.Value, error) {
		log.Println("执行 http 的操作，参数是 ", args, len(args))
		return nil, nil
	})
}
