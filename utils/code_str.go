package utils

import (
	"log"
	"regexp"
	"strings"
)

func EscapeQuotesInBackticks(input string) string {
	//re := regexp.MustCompile(`(` + "`[^`]*`" + `)`)
	//result := re.ReplaceAllStringFunc(input, func(match string) string {
	//	content := strings.Trim(match, "`")
	//	escapedContent := strings.ReplaceAll(content, `"`, `\"`)
	//	return "`" + escapedContent + "`"
	//})
	//log.Println("result = ", result)
	//return result

	lines := strings.Split(input, "\n")
	// 预编译正则（提升性能）
	reBacktick := regexp.MustCompile("`([^`]*)`")    // 匹配`包裹的内容
	reDoubleQuote := regexp.MustCompile(`"([^"]*)"`) // 匹配"包裹的内容

	// 步骤2：遍历每行处理
	processedLines := make([]string, len(lines))
	for i, line := range lines {
		trimmedLine := strings.TrimLeft(line, " \t")
		if trimmedLine == "" || !strings.HasPrefix(trimmedLine, "chrome") {
			log.Println("没 chrome 不需要处理 -> ", line)
			processedLines[i] = line
			continue
		}

		// 处理1：反引号` `内的"→\"、'→\"（新增单引号转义）
		line = reBacktick.ReplaceAllStringFunc(line, func(match string) string {
			content := match[1 : len(match)-1] // 去掉首尾`
			// 先替换双引号，再替换单引号（顺序不影响）
			escaped := strings.ReplaceAll(content, `"`, `\"`)
			escaped = strings.ReplaceAll(escaped, `'`, `\"`) // 新增：单引号转\"
			return "`" + escaped + "`"
		})

		// 处理2：双引号""内的'→\'
		line = reDoubleQuote.ReplaceAllStringFunc(line, func(match string) string {
			content := match[1 : len(match)-1] // 去掉首尾"
			escaped := strings.ReplaceAll(content, `'`, `\'`)
			return `"` + escaped + `"`
		})

		processedLines[i] = line
	}

	// 步骤3：还原换行结构
	result := strings.Join(processedLines, "\n")
	log.Println("处理完成，总行数：", len(processedLines))
	log.Println("result：", result)
	return result
}
