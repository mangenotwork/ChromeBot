package builtins

import (
	"ChromeBot/dsl/interpreter"
	"log"
)

/*
参数说明
method ：请求方式 get post put delete options head patch
url : 请求的url
body ： 请求的body
header ： 请求的header
ctype ： 请求的 是 header key 为 Content-Type
cookie ：请求的cookie 是 header key 为 Cookie
timeout ：设置请求的超时时间
proxy ：设置请求的代理，目前只支持 http/https代理
stress ：压力请求，并发请求设置的数量
*/
func registerHttp(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("http", func(args []interpreter.Value) (interpreter.Value, error) {
		log.Println("执行 http 的操作，参数是 ", args, len(args))
		return nil, nil
	})
}
