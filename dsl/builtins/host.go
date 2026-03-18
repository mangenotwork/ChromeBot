package builtins

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/internal/host"
	"ChromeBot/utils"
	"fmt"
	"strings"
)

var hostSupport = map[string]bool{
	"info": true,
	"name": true,
	"ip":   true,
	"to":   true,
	"disk": true,
}

func hasHostSupport(cmd string) bool {
	_, ok := hostSupport[cmd]
	return ok
}

/*
host 关键字，系统相关的操作与系统相关的命令; 一个命令只执行一个参数。

参数说明：

info : 获取系统的信息
name : 获取系统的名称
ip : 获取系统的ip
to : 将当前操作返回的值存入到指定变量-如果变量未声明这里会自动声明变量  <值类型是字符串>
disk : todo 系统的磁盘信息
*/
func registerHost(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("host", func(args []interpreter.Value) (interpreter.Value, error) {
		fmt.Println("执行 host 的操作，参数是 ", args, len(args))

		argsStr := make([]string, 0)
		for i, arg := range args {
			utils.Debugf("参数 %d %v %T\n", i, arg, arg)
			argsStr = append(argsStr, arg.(string))
		}

		// 处理 函数类型的参数
		argsStr = processArgs(interp, argsStr)
		utils.Debug("执行 ProcessArgs 参数 处理  ", argsStr, len(args))

		if len(argsStr) == 0 {
			return nil, fmt.Errorf("未知命令，请参考文档")
		}

		argMap := make(map[string]string)
		for _, v := range argsStr {
			vList := strings.SplitN(v, "=", 2)
			utils.Debug("vList = ", vList, len(vList))
			if !hasHostSupport(vList[0]) {
				fmt.Println("[Host]未知命令 ", vList[0], ";请参考文档。")
				return nil, fmt.Errorf("[Host]未知命令 %s;请参考文档。", vList[0])
			}
			if len(vList) == 1 {
				argMap[vList[0]] = ""
			} else if len(vList) == 2 {
				argMap[vList[0]] = vList[1]
			}
		}

		utils.Debug("argMap:", argMap)

		toArg, isTo := argMap["to"]

		if _, ok := argMap["info"]; ok {
			osInfo := host.GetOSInfo()
			if isTo {
				rse := make(interpreter.DictType)
				rse["HostName"] = osInfo.HostName
				rse["OSType"] = osInfo.OSType
				rse["OSArch"] = osInfo.OSArch
				rse["CpuCoreNumber"] = osInfo.CpuCoreNumber
				rse["InterfaceInfo"] = osInfo.InterfaceInfo
				interp.Global().SetVar(toArg, rse)
			}
			return nil, nil
		}

		if _, ok := argMap["name"]; ok {
			osName := host.GetOSName()
			if isTo {
				interp.Global().SetVar(toArg, osName)
			}
			return nil, nil
		}

		if _, ok := argMap["ip"]; ok {
			osIP := host.GetOSIP()
			if isTo {
				interp.Global().SetVar(toArg, osIP)
			}
			return nil, nil
		}

		return nil, nil
	})
}
