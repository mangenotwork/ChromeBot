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
	"ls":   true,
	"file": true,
	"goto": true,
	"from": true,
	"s":    true,
	"c":    true,
	"d":    true,
	"m":    true,
	"cp":   true,
	"r":    true,
	"renm": true,
	"w":    true,
	"a":    true,
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
ls : 列出文件或目录
file : 操作系统文件
  - s=<search word> root=<path> : 搜索文件或目录
  - c=<path> : 创建文件或目录
  - d=<path> : 删除文件或目录
  - m=<path> goto=<path> : 移动文件或目录
  - cp=<path> goto=<path> : 复制文件或目录
  - r=<path> to=<arg> : 读文件
  - renm=<path> goto=<path> : 文件或目录改名, 路径不同则移动
  - info=<path> : 文件或目录信息
  - w=<path> from=<arg> : 将文件内容写入文件
  - a=<path> from=<arg> : 将文件内容追加写入文件
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
		gotoArg, isGoto := argMap["goto"]
		fromArg, isFrom := argMap["from"]

		_, hasInfo := argMap["info"]
		_, hasName := argMap["name"]
		_, hasIP := argMap["ip"]
		lsArg, hasLS := argMap["ls"]
		_, hasFile := argMap["file"]
		if hasFile { // 解决命令冲突
			hasInfo = false
		}

		switch {

		case hasInfo:
			osInfo := host.GetOSInfo()
			if isTo {
				interp.Global().SetVar(toArg, interpreter.DictType{
					"HostName":      osInfo.HostName,
					"OSType":        osInfo.OSType,
					"OSArch":        osInfo.OSArch,
					"CpuCoreNumber": osInfo.CpuCoreNumber,
					"InterfaceInfo": osInfo.InterfaceInfo,
				})
			}

		case hasName:
			osName := host.GetOSName()
			if isTo {
				interp.Global().SetVar(toArg, osName)
			}

		case hasIP:
			osIP := host.GetOSIP()
			if isTo {
				interp.Global().SetVar(toArg, osIP)
			}

		case hasLS:
			if lsArg == "" {
				lsArg = utils.ScriptDir
			}
			host.LS(lsArg)

		case hasFile:

			sArg, sOK := argMap["s"]
			cArg, cOK := argMap["c"]
			dArg, dOK := argMap["d"]
			mArg, mOK := argMap["m"]
			cpArg, cpOK := argMap["cp"]
			rArg, rOK := argMap["r"]
			renmArg, renmOK := argMap["renm"]
			infoArg, infoOK := argMap["info"]
			wArg, wOK := argMap["w"]
			aArg, aOK := argMap["a"]
			pathType := 1

			taskCount := 0
			for _, ok := range []bool{sOK, cOK, dOK, mOK, cpOK, rOK, renmOK, infoOK, wOK, aOK} {
				if ok {
					taskCount++
				}
			}

			if taskCount > 1 {
				fmt.Println("[Err]一行指令代码只做一个任务")
				break
			}

			switch {
			case sOK:
				fmt.Println("sArg = ", sArg)

				rootArg, rootOK := argMap["root"]
				if !rootOK {
					rootArg = utils.ScriptDir
				}

				res, err := host.SearchFilesDir(rootArg, sArg, true, true)
				if err != nil {
					fmt.Println("[Err]搜索失败err:", err.Error())
				}
				utils.ShowJson(res)

			case cOK:
				cArg, pathType = host.CheckPath(cArg)
				fmt.Println("cArg = ", cArg, " | pathType = ", pathType)
				var err error
				if pathType == 1 {
					err = host.CreateFile(infoArg)
				} else {
					err = host.CreateDir(infoArg)
				}
				if err != nil {
					fmt.Println("[Err]创建失败, err = ", err)
				}

			case dOK:
				dArg, pathType = host.CheckPath(dArg)
				fmt.Println("dArg = ", dArg, " | pathType = ", pathType)
				err := host.DeleteFile(dArg)
				if err != nil {
					fmt.Println("[Err]删除失败, err = ", err)
				}

			case mOK:
				mArg, pathType = host.CheckPath(mArg)
				fmt.Println("mArg = ", mArg, " | pathType = ", pathType)
				if !isGoto {
					fmt.Println("[Err]缺少goto")
					break
				}
				gotoArgPathType := 1
				gotoArg, gotoArgPathType = host.CheckPath(gotoArg)
				fmt.Println("gotoArg = ", gotoArg, " | gotoArgPathType = ", gotoArgPathType)

				err := host.MoveFile(mArg, gotoArg)
				if err != nil {
					fmt.Println("[Err]移动失败, err = ", err)
				}

			case cpOK:
				cpArg, pathType = host.CheckPath(cpArg)
				fmt.Println("cpArg = ", cpArg, " | pathType = ", pathType)
				if !isGoto {
					fmt.Println("[Err]缺少goto")
					break
				}
				gotoArgPathType := 1
				gotoArg, gotoArgPathType = host.CheckPath(gotoArg)
				fmt.Println("gotoArg = ", gotoArg, " | gotoArgPathType = ", gotoArgPathType)

				err := host.CopyFile(cpArg, gotoArg)
				if err != nil {
					fmt.Println("[Err]复制失败, err = ", err)
				}

			case rOK:
				rArg, pathType = host.CheckPath(rArg)
				fmt.Println("rArg = ", rArg, " | pathType = ", pathType)
				if !isTo {
					fmt.Println("[Err]缺少to")
					break
				}
				toArgPathType := 1
				toArg, toArgPathType = host.CheckPath(toArg)
				fmt.Println("toArg = ", toArg, " | toArgPathType = ", toArgPathType)

				str, err := host.ReadFileToString(rArg)
				if err != nil {
					fmt.Println("[Err]读取文件失败, err = ", err)
					break
				}

				if isTo {
					interp.Global().SetVar(toArg, str)
				}

			case renmOK:
				renmArg, pathType = host.CheckPath(renmArg)
				fmt.Println("renmArg = ", renmArg, " | pathType = ", pathType)
				if !isGoto {
					fmt.Println("[Err]缺少goto")
					break
				}
				gotoArgPathType := 1
				gotoArg, gotoArgPathType = host.CheckPath(gotoArg)
				fmt.Println("gotoArg = ", gotoArg, " | gotoArgPathType = ", gotoArgPathType)

				err := host.RenameOrMove(renmArg, gotoArg)
				if err != nil {
					fmt.Println("[Err]重命名失败, err = ", err)
				}

			case infoOK:
				infoArg, pathType = host.CheckPath(infoArg)
				fmt.Println("infoArg = ", infoArg, " | pathType = ", pathType)
				rse := make(interpreter.DictType)
				var err error
				if pathType == 1 {
					rse, err = host.GetFileInfo(infoArg)
				} else {
					rse, err = host.GetDirInfo(infoArg)
				}
				if err != nil {
					fmt.Println("[Err]获取文件或目录信息失败, err = ", err)
				}
				fmt.Println("Info : ")
				utils.ShowJson(rse)
				if isTo {
					interp.Global().SetVar(toArg, rse)
				}

			case wOK:
				wArg, pathType = host.CheckPath(wArg)
				fmt.Println("wArg = ", wArg, " | pathType = ", pathType)
				if !isFrom {
					fmt.Println("[Err]缺少from")
					break
				}
				fromArgPathType := 1
				fromArg, fromArgPathType = host.CheckPath(fromArg)
				fmt.Println("fromArg = ", fromArg, " | fromArgPathType = ", fromArgPathType)

				strVal, strValOK := interp.Global().GetVar(wArg)
				if !strValOK {
					fmt.Printf("[Err]%s变量不存在\n", fromArg)
					break
				}

				err := host.WriteFileOverwrite(wArg, strVal.(string))
				if err != nil {
					fmt.Println("[Err]写入文件失败, err = ", err)
				}

			case aOK:
				aArg, pathType = host.CheckPath(aArg)
				fmt.Println("aArg = ", aArg, " | pathType = ", pathType)
				if !isFrom {
					fmt.Println("[Err]缺少from")
					break
				}
				fromArgPathType := 1
				fromArg, fromArgPathType = host.CheckPath(fromArg)
				fmt.Println("fromArg = ", fromArg, " | fromArgPathType = ", fromArgPathType)

				strVal, strValOK := interp.Global().GetVar(wArg)
				if !strValOK {
					fmt.Printf("[Err]%s变量不存在\n", fromArg)
					break
				}

				err := host.AppendToFile(wArg, strVal.(string))
				if err != nil {
					fmt.Println("[Err]写入文件失败, err = ", err)
				}

			}

		}

		return nil, nil
	})
}
