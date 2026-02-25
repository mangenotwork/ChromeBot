package builtins

import (
	"ChromeBot/dsl/ast"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/internal/httpclient"
	"ChromeBot/utils"
	"encoding/json"
	gt "github.com/mangenotwork/gathertool"
	"strings"
)

/*
参数说明
method ：请求方式 get post put delete options head patch
url : 请求的url,要求类型是str
body ： 请求的body,要求类型者是str或是List和字典（根据ctype解析为from-data，json这些）
header ： 请求的header,要求类型是字典或者是json str
ctype ： 请求的 是 header key 为 Content-Type, 要求类型是str
cookie ：请求的cookie 是 header key 为 Cookie, 要求类型者是str(k=v;)或是List和字典（会解析为 k=v;） list是 ["k1=v1", "k2=v2"...]
timeout ：设置请求的超时时间单位为毫秒, 要求类型是数值
proxy ：设置请求的代理，目前只支持 http/https代理, 要求类型是str
stress ：压力请求，并发请求设置的数量，要求类型是数值
save : 指定将响应内容存储，要求类型是str,本地文件路径
to : 将请求的返回存入到指定变量-如果变量未声明这里会自动声明变量
*/
func registerHttp(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("http", func(args []interpreter.Value) (interpreter.Value, error) {
		utils.Debug("执行 http 的操作，参数是 ", args, len(args))

		argMap := args[0].(map[string]interpreter.Value)

		// 校验参数类型
		req := &httpclient.HttpReq{
			Method: argMap["method"].(string),
			Header: gt.Header{},
			Cookie: make(gt.Cookie),
		}

		if val, ok := argMap["url"]; ok {
			utils.Debugf("http val T : %T \n", val)
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

		if val, ok := argMap["body"]; ok {
			utils.Debugf("http body T : %T \n", val)
			switch val.(type) {
			case string:
				req.Body = []byte(val.(string))
			case *ast.String:
				req.Body = []byte(val.(*ast.String).Value)
			case interpreter.DictType, []interpreter.Value: // 需要处理为json // todo 根据 cType来决定val类型
				valJson, err := json.Marshal(val)
				if err != nil {
					utils.Debug("err :", err)
				}
				req.Body = valJson
				utils.Debug("valJson = ", string(valJson))
			default:
				interp.ErrorMessage("body参数要求类型是str")
				return nil, nil
			}
		}

		if val, ok := argMap["header"]; ok {
			utils.Debugf("http header T : %T  %s \n", val, val)
			switch val.(type) {
			case string:
				vMap := make(map[string]interface{})
				err := json.Unmarshal([]byte(val.(string)), &vMap)
				if err != nil {
					utils.Debug("解析header err :", err)
				}
				utils.Debug("vMap = ", vMap)
				for k, v := range vMap {
					req.Header[gt.Any2String(k)] = gt.Any2String(v)
				}
			case *ast.String:
				vMap := make(map[string]interface{})
				err := json.Unmarshal([]byte(val.(*ast.String).Value), &vMap)
				if err != nil {
					utils.Debug("解析header err :", err)
				}
				for k, v := range vMap {
					req.Header[gt.Any2String(k)] = gt.Any2String(v)
				}
			case interpreter.DictType:
				for k, v := range val.(interpreter.DictType) {
					req.Header[gt.Any2String(k)] = gt.Any2String(v)
				}
			default:
				interp.ErrorMessage("header参数要求类型者是str（json字符串）或是List和字典（根据ctype解析为from-data，json这些）")
				return nil, nil
			}
		}

		if val, ok := argMap["ctype"]; ok {
			utils.Debugf("ctype val T : %T \n", val)
			switch val.(type) {
			case string:
				req.Ctype = val.(string)
			case *ast.String:
				req.Ctype = val.(*ast.String).Value
			default:
				interp.ErrorMessage("ctype参数要求类型是str")
				return nil, nil
			}
		}

		if val, ok := argMap["cookie"]; ok {
			utils.Debugf("cookie val T : %T \n", val)
			switch val.(type) {
			case string:
				vList := strings.Split(val.(string), ";")
				utils.Debug("vList = ", vList)
				for _, items := range vList {
					item := strings.Split(items, "=")
					if len(item) == 2 {
						req.Cookie[gt.Any2String(item[0])] = gt.Any2String(item[1])
					}
				}
			case *ast.String:
				vList := strings.Split(val.(*ast.String).Value, ";")
				utils.Debug("vList = ", vList)
				for _, items := range vList {
					item := strings.Split(items, "=")
					if len(item) == 2 {
						req.Cookie[gt.Any2String(item[0])] = gt.Any2String(item[1])
					}
				}
			case interpreter.DictType:
				for k, v := range val.(interpreter.DictType) {
					req.Cookie[gt.Any2String(k)] = gt.Any2String(v)
				}
			case []interpreter.Value:
				for _, v := range val.([]interpreter.Value) {
					vStr := gt.Any2String(v)
					vItems := strings.Split(vStr, "=")
					if len(vItems) == 2 {
						req.Cookie[gt.Any2String(vItems[0])] = gt.Any2String(vItems[1])
					} else {
						interp.ErrorMessage("cookie参数类型者是List要求元素是[\"k1=v1\", \"k2=v2\"...]")
						return nil, nil
					}
				}
			default:
				interp.ErrorMessage("cookie参数要求类型者是str(k=v;)或是List和字典（会解析为 k=v;）")
				return nil, nil
			}
		}

		if val, ok := argMap["timeout"]; ok {
			utils.Debugf("timeout val T : %T \n", val)
			switch val.(type) {
			case string:
				req.Timeout = gt.ReqTimeOut(gt.Any2Int(val))
			case *ast.String:
				req.Timeout = gt.ReqTimeOut(gt.Any2Int(val.(*ast.String).Value))
			case int, int64, float64:
				req.Timeout = gt.ReqTimeOut(gt.Any2Int(val))
			case *ast.Integer:
				req.Timeout = gt.ReqTimeOut(val.(*ast.Integer).Value)
			default:
				interp.ErrorMessage("timeout参数要求类型是数值类型，单位为毫秒")
				return nil, nil
			}
		}

		if val, ok := argMap["proxy"]; ok {
			utils.Debugf("proxy val T : %T \n", val)
			switch val.(type) {
			case string:
				req.Proxy = gt.ProxyUrl(val.(string))
			case *ast.String:
				req.Proxy = gt.ProxyUrl(val.(*ast.String).Value)
			default:
				interp.ErrorMessage("proxy参数要求类型是str")
				return nil, nil
			}
		}

		if val, ok := argMap["stress"]; ok {
			utils.Debugf("stress val T : %T \n", val)
			switch val.(type) {
			case string:
				req.Stress = gt.Any2Int(val)
			case *ast.String:
				req.Stress = gt.Any2Int(val.(*ast.String).Value)
			case *ast.Integer:
				req.Stress = val.(int)
			default:
				interp.ErrorMessage("stress参数要求类型是数值")
				return nil, nil
			}
		}

		// todo 存储方式
		if val, ok := argMap["save"]; ok {
			utils.Debugf("save val T : %T \n", val)
			switch val.(type) {
			case string:
				req.Save = val.(string)
			case *ast.String:
				req.Save = val.(*ast.String).Value
			default:
				interp.ErrorMessage("save参数要求类型是字符串")
				return nil, nil
			}
		}

		utils.Debug(" req = ", req)

		rse := req.Do()
		utils.Debug(" rse = ", rse)

		if to, ok := argMap["to"]; ok {
			utils.Debug("http请求结果保存到变量: ", to)
			rseDict := make(interpreter.DictType)
			rseDict["code"] = rse.Code
			rseDict["body"] = string(rse.Body)
			rseDict["header"] = gt.Any2String(rse.Header)
			rseDict["req_time"] = rse.ReqTime
			if rse.Error != nil {
				rseDict["err"] = rse.Error.Error()
			}
			interp.Global().SetVar(to.(string), rseDict)
		}

		return nil, nil
	})
}
