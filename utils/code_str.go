package utils

import (
	"os"
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
			Debug("没 chrome 不需要处理 -> ", line)
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
	Debug("处理完成，总行数：", len(processedLines))
	Debug("result：", result)
	return result
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// FixURLProtocol 自动补全URL的协议头（默认补https://，也可指定http）
// 参数：
//
//	url - 原始URL（如 "www.baidu.com"、"baidu.com:8080"、"https://xxx"）
//	defaultProto - 默认协议头（传 "" 则用 https://，传 "http://" 则补http）
//
// 返回：
//
//	补全后的合法URL
func FixURLProtocol(url string, defaultProto ...string) string {
	// 空URL直接返回
	if strings.TrimSpace(url) == "" {
		return url
	}

	// 确定默认协议头（优先用传入的，否则用https://）
	proto := "https://"
	if len(defaultProto) > 0 && defaultProto[0] != "" {
		proto = defaultProto[0]
		// 确保协议头以://结尾（容错：用户传"https"时自动补://）
		if !strings.HasSuffix(proto, "://") {
			proto += "://"
		}
	}

	// 检查URL是否已带协议头（http/https/ftp等），有则直接返回
	lowerURL := strings.ToLower(url)
	if strings.HasPrefix(lowerURL, "http://") ||
		strings.HasPrefix(lowerURL, "https://") ||
		strings.HasPrefix(lowerURL, "ftp://") {
		return url
	}

	// 补全协议头并返回
	return proto + url
}
