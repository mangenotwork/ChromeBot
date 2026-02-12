package builtins

import (
	"ChromeBot/dsl/ast"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/internal/httpclient"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
)

/*
参数说明
method ：请求方式 get post put delete options head patch
url : 请求的url,要求类型是str
body ： 请求的body,要求类型是字典或者是str
header ： 请求的header,要求类型是字典或者是json str
ctype ： 请求的 是 header key 为 Content-Type, 要求类型是str
cookie ：请求的cookie 是 header key 为 Cookie, 要求类型是字典或是str
timeout ：设置请求的超时时间, 要求类型是数值
proxy ：设置请求的代理，目前只支持 http/https代理, 要求类型是str
stress ：压力请求，并发请求设置的数量，要求类型是数值
save : 指定将响应内容存储，要求类型是str,本地文件路径
to : 将请求的返回存入到指定变量-如果变量未声明这里会自动声明变量
*/
func registerHttp(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("http", func(args []interpreter.Value) (interpreter.Value, error) {
		log.Println("执行 http 的操作，参数是 ", args, len(args))

		argMap := args[0].(map[string]interpreter.Value)

		// 校验参数类型
		req := &httpclient.HttpReq{
			Method: argMap["method"].(string),
		}

		if val, ok := argMap["url"]; ok {
			fmt.Printf("val : %T \n", val)
			switch val.(type) {
			case string:
				req.Url = val.(string)
			case *ast.String:
				req.Url = val.(*ast.String).Value
			default:
				interp.ErrorMessage("url参数要求类型是str")
				return nil, nil
			}
		}

		// todo
		if val, ok := argMap["body"]; ok {
			fmt.Printf("body : %T \n", val)
		}

		// todo
		if val, ok := argMap["header"]; ok {
			fmt.Printf("header : %T \n", val)
		}

		// todo
		if val, ok := argMap["ctype"]; ok {
			fmt.Printf("ctype : %T \n", val)
		}

		// todo
		if val, ok := argMap["cookie"]; ok {
			fmt.Printf("cookie : %T \n", val)
		}

		// todo
		if val, ok := argMap["timeout"]; ok {
			fmt.Printf("timeout : %T \n", val)
		}

		// todo
		if val, ok := argMap["proxy"]; ok {
			fmt.Printf("proxy : %T \n", val)
		}

		if val, ok := argMap["stress"]; ok {
			fmt.Printf("stress : %T \n", val)
			switch val.(type) {
			case string:
				req.Stress = gt.Any2Int(val)
			case *ast.Integer:
				req.Stress = val.(int)
			default:
				interp.ErrorMessage("stress参数要求类型是数值")
				return nil, nil
			}
		}

		if val, ok := argMap["save"]; ok {
			fmt.Printf("save : %T \n", val)
		}

		fmt.Println(" req = ", req)

		req.Do()

		if to, ok := argMap["to"]; ok {
			fmt.Println("http请求结果保存到变量: ", to)
			interp.Global().SetVar(to.(string), "http resp")
		}

		return nil, nil
	})
}
