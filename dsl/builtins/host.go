package builtins

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/internal/host"
	"ChromeBot/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"
)

var hostSupport = map[string]bool{
	"info": true,
	"name": true,
	"ip":   true,
	"to":   true,
	"disk": true,
	"ls":   true,
	"file": true,
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
  - s=<search word> : 搜索文件或目录
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
			ls(lsArg)

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

			case cOK:
				cArg, pathType = checkPath(cArg)
				fmt.Println("cArg = ", cArg, " | pathType = ", pathType)

			case dOK:
				dArg, pathType = checkPath(dArg)
				fmt.Println("dArg = ", dArg, " | pathType = ", pathType)

			case mOK:
				mArg, pathType = checkPath(mArg)
				fmt.Println("mArg = ", mArg, " | pathType = ", pathType)
				if !isGoto {
					fmt.Println("[Err]缺少goto")
					break
				}
				gotoArgPathType := 1
				gotoArg, gotoArgPathType = checkPath(gotoArg)
				fmt.Println("gotoArg = ", gotoArg, " | gotoArgPathType = ", gotoArgPathType)

			case cpOK:
				cpArg, pathType = checkPath(cpArg)
				fmt.Println("cpArg = ", cpArg, " | pathType = ", pathType)
				if !isGoto {
					fmt.Println("[Err]缺少goto")
					break
				}
				gotoArgPathType := 1
				gotoArg, gotoArgPathType = checkPath(gotoArg)
				fmt.Println("gotoArg = ", gotoArg, " | gotoArgPathType = ", gotoArgPathType)

			case rOK:
				rArg, pathType = checkPath(rArg)
				fmt.Println("rArg = ", rArg, " | pathType = ", pathType)
				if !isTo {
					fmt.Println("[Err]缺少to")
					break
				}
				toArgPathType := 1
				toArg, toArgPathType = checkPath(toArg)
				fmt.Println("toArg = ", toArg, " | toArgPathType = ", toArgPathType)

			case renmOK:
				renmArg, pathType = checkPath(renmArg)
				fmt.Println("renmArg = ", renmArg, " | pathType = ", pathType)
				if !isGoto {
					fmt.Println("[Err]缺少goto")
					break
				}
				gotoArgPathType := 1
				gotoArg, gotoArgPathType = checkPath(gotoArg)
				fmt.Println("gotoArg = ", gotoArg, " | gotoArgPathType = ", gotoArgPathType)

			case infoOK:
				infoArg, pathType = checkPath(infoArg)
				fmt.Println("infoArg = ", infoArg, " | pathType = ", pathType)
				rse := make(interpreter.DictType)
				var err error
				if pathType == 1 {
					rse, err = getFileInfo(infoArg)
				} else {
					rse, err = getDirInfo(infoArg)
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
				wArg, pathType = checkPath(wArg)
				fmt.Println("wArg = ", wArg, " | pathType = ", pathType)
				if !isFrom {
					fmt.Println("[Err]缺少from")
					break
				}
				fromArgPathType := 1
				fromArg, fromArgPathType = checkPath(fromArg)
				fmt.Println("fromArg = ", fromArg, " | fromArgPathType = ", fromArgPathType)

			case aOK:
				aArg, pathType = checkPath(aArg)
				fmt.Println("aArg = ", aArg, " | pathType = ", pathType)
				if !isFrom {
					fmt.Println("[Err]缺少from")
					break
				}
				fromArgPathType := 1
				fromArg, fromArgPathType = checkPath(fromArg)
				fmt.Println("fromArg = ", fromArg, " | fromArgPathType = ", fromArgPathType)

			}

		}

		return nil, nil
	})
}

// 返回 1是文件，2是目录
func checkPath(pathStr string) (string, int) {
	if pathStr == "" {
		return utils.ScriptDir, 2
	}
	if !filepath.IsAbs(pathStr) {
		pathStr = filepath.Join(utils.ScriptDir, pathStr)
	}
	pathStr = filepath.Clean(pathStr)
	stat, err := os.Stat(pathStr)
	if err != nil {
		return pathStr, 2
	}
	if stat.IsDir() {
		return pathStr, 2
	}
	return pathStr, 1
}

func ls(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "名称\t大小\t修改时间\t类型")
	fmt.Fprintln(w, "----\t----\t--------\t----")

	var totalSize int64

	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		var size int64

		if file.IsDir() {
			size = calcDirSizeParallel(fullPath) // 🔥 并行计算
		} else {
			size = file.Size()
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			file.Name(),
			formatSize(size),
			file.ModTime().Format("2006-01-02 15:04:05"),
			dirOrFile(file.IsDir()),
		)
		totalSize += size
	}

	w.Flush()
	fmt.Println("\n✅ 目录总大小：", formatSize(totalSize))
	return nil
}

var (
	lsSizeChan = make(chan int64, 1000)
	lsWg       sync.WaitGroup
	lsSema     = make(chan struct{}, 50) // 限制并发数，防止卡死
)

// 🔥 并行计算目录大小（超快）
func calcDirSizeParallel(root string) int64 {
	var total int64
	lsWg.Add(1)

	// 并发扫描
	go func() {
		defer lsWg.Done()
		scanDir(root)
	}()

	// 关闭通道并等待
	go func() {
		lsWg.Wait()
		close(lsSizeChan)
	}()

	// 汇总大小
	for s := range lsSizeChan {
		total += s
	}

	// 重置
	lsSizeChan = make(chan int64, 1000)
	return total
}

// 🔥 递归扫描（带并发限流，极快）
func scanDir(path string) {
	lsSema <- struct{}{}        //  acquire
	defer func() { <-lsSema }() // release

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			lsWg.Add(1)
			go func(p string) {
				defer lsWg.Done()
				scanDir(p)
			}(fullPath)
			continue
		}

		// 文件加入大小
		info, err := entry.Info()
		if err == nil {
			lsSizeChan <- info.Size()
		}
	}
}

func dirOrFile(isDir bool) string {
	if isDir {
		return "目录"
	}
	return "文件"
}

func formatSize(s int64) string {
	const unit = 1024
	if s < unit {
		return fmt.Sprintf("%d B", s)
	}
	div, exp := int64(unit), 0
	for n := s / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %sB", float64(s)/float64(div), []string{"K", "M", "G", "T"}[exp])
}

func getFileInfo(filePath string) (interpreter.DictType, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileInfo := interpreter.DictType{
		"name":    stat.Name(),
		"path":    filePath,
		"size":    stat.Size(),
		"sizeStr": formatSize(stat.Size()),
		"modtime": stat.ModTime().Format("2006-01-02 15:04:05"),
		"isdir":   stat.IsDir(),
		"mode":    stat.Mode().String(),
	}
	return fileInfo, nil
}

func getDirInfo(dirPath string) (interpreter.DictType, error) {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	fileCount := 0
	dirCount := 0
	entries, err := os.ReadDir(dirPath)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				dirCount++
			} else {
				fileCount++
			}
		}
	}

	dirInfo := interpreter.DictType{
		"name":      stat.Name(),
		"path":      dirPath,
		"modtime":   stat.ModTime().Format("2006-01-02 15:04:05"),
		"isdir":     stat.IsDir(),
		"mode":      stat.Mode().String(),
		"dircount":  dirCount,
		"filecount": fileCount,
		"count":     fileCount + dirCount,
	}
	return dirInfo, nil
}
