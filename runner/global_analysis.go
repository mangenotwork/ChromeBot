package runner

import (
	"ChromeBot/global"
	"ChromeBot/utils"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func globalAnalysisScript(input string) string {
	lines := strings.Split(input, "\n")
	processedLines := make([]string, len(lines))
	for i, line := range lines {
		trimmedLine := strings.TrimLeft(line, " \t")
		if trimmedLine == "" || !strings.HasPrefix(trimmedLine, "@") {
			utils.Debug("没 @ 不需要处理 -> ", line)
			processedLines[i] = line
			continue
		}

		// 匹配
		fmt.Println("@的语法句 : ", line)
		globalAnalysis(line)

		// 将第一个@替换为#
		line = strings.Replace(line, "@", "#", 1)

		processedLines[i] = line
	}

	// 步骤3：还原换行结构
	result := strings.Join(processedLines, "\n")
	//fmt.Println("处理完成，总行数：", len(processedLines))
	//fmt.Println("result：", result)
	//fmt.Println("global解析完成  -------- ")
	return result
}

func globalAnalysisLine(line string) string {
	trimmedLine := strings.TrimLeft(line, " \t")
	if trimmedLine == "" || !strings.HasPrefix(trimmedLine, "@") {
		utils.Debug("没 @ 不需要处理 -> ", line)
		return line
	}
	fmt.Println("[Wrong]REPL模式下不支持@语法")
	line = strings.Replace(line, "@", "#", 1)
	return line
}

func globalAnalysis(line string) {

	if strings.Contains(line, "//") {
		line = strings.SplitN(line, "//", 2)[0]
	}

	if strings.Contains(line, "#") {
		line = strings.SplitN(line, "#", 2)[0]
	}

	re := regexp.MustCompile(`\s*=\s*`) // 去掉 = 两边的空格
	normalized := re.ReplaceAllString(line, "=")
	lineList := strings.Fields(normalized)
	fmt.Println("lineList len = ", len(lineList), " v = ", lineList)
	if len(lineList) > 0 {
		command := lineList[0]
		command = strings.ReplaceAll(command, "@", "")
		globalCommand := global.Command(command)
		if has := global.HasGlobalSupport(globalCommand); !has {
			fmt.Printf("[Err]未知global指令%s, 请参考文档。\n", command)
			os.Exit(0)
		}

		fmt.Println("arg = ", lineList[1:])

		switch globalCommand {

		case global.Cron:
			if global.IsRegisterCron {
				fmt.Println("[Wrong]已设置过@cron,只能设置一次。")
			}
			arg := strings.Join(lineList[1:], " ")
			global.RegisterCron(arg)

		case global.ConfJson, global.ConfYaml, global.ConfINI:
			argMap := make(map[string]string)
			for _, argItem := range lineList[1:] {
				argItemList := strings.SplitN(argItem, "=", 2)
				argMap[argItemList[0]] = argItemList[1]
			}

			path, pathOK := argMap["path"]
			if !pathOK {
				fmt.Printf("[Err]@%s缺少参数path.正确语法@%s path=\"\" as=conf", globalCommand, globalCommand)
				os.Exit(0)
			}

			as, asOK := argMap["as"]
			if !asOK {
				fmt.Printf("[Err]@%s缺少参数as.正确语法@%s path=\"\" as=conf", globalCommand, globalCommand)
				os.Exit(0)
			}

			switch globalCommand {
			case global.ConfJson:
				global.ReadJsonToConf(path, as)
			case global.ConfYaml:
				global.ReadYamlToConf(path, as)
			case global.ConfINI:
				global.ReadINIToConf(path, as)
			}

		case global.ChromeCheck:
			global.ChromeCheckImplement()

		case global.NetworkCheck:
			if len(lineList) < 2 {
				fmt.Println(`[Wrong]@network_check缺少参数,参考语法 @network_check "254.254.254.254"  `)
				os.Exit(0)
			}
			arg := lineList[1]

			if !IsDomainOrIP(arg) {
				fmt.Println(`[Wrong]@network_check参数应该是ip或域名,参考语法 @network_check "254.254.254.254"  `)
				os.Exit(0)
			}

			global.NetworkCheckImplement(arg)

		}
	}

}

func IsDomainOrIP(s string) bool {
	if net.ParseIP(s) != nil {
		return true
	}

	// 移除协议前缀
	if strings.Contains(s, "://") {
		s = s[strings.Index(s, "://")+3:]
	}

	// 移除端口
	if idx := strings.Index(s, ":"); idx != -1 {
		s = s[:idx]
	}

	// 移除路径
	if idx := strings.Index(s, "/"); idx != -1 {
		s = s[:idx]
	}

	// 移除首尾空格和点
	s = strings.TrimSpace(s)
	s = strings.Trim(s, ".")

	if s == "" || len(s) > 253 {
		return false
	}

	// 必须包含点
	if !strings.Contains(s, ".") {
		return false
	}

	// 简单检查：不包含空格，有合理结构
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}

	// 最后一个部分至少2个字符
	if len(parts[len(parts)-1]) < 2 {
		return false
	}

	return true
}
