package runner

import (
	"ChromeBot/global"
	"ChromeBot/utils"
	"fmt"
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
		switch globalCommand {
		case global.Cron:
			if global.IsRegisterCron {
				fmt.Println("[Wrong]已设置过@cron,只能设置一次。")
			}
			arg := strings.Join(lineList[1:len(lineList)], " ")
			global.RegisterCron(arg)
		}
	}

}
