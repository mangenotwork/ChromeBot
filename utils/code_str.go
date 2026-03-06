package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func SaveDataToFile(path string, data interface{}) error {
	if path == "" {
		return errors.New("文件路径不能为空")
	}

	// 解析路径（处理相对路径→绝对路径，创建父目录）
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("解析路径失败：%w", err)
	}
	// 获取父目录（如 "/tmp/data/test.txt" → "/tmp/data"）
	dir := filepath.Dir(absPath)
	// 创建父目录（不存在则创建，递归创建多级目录）
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建父目录失败：%w", err)
	}

	// 转换数据为字节数组（适配不同数据类型）
	var content []byte
	switch v := data.(type) {
	case string:
		content = []byte(v)
	case []byte:
		content = v
	default:
		// 其他类型尝试JSON序列化（如结构体、map等）
		jsonContent, err := json.MarshalIndent(v, "", "  ") // 格式化JSON，易读
		if err != nil {
			return fmt.Errorf("数据JSON序列化失败：%w", err)
		}
		content = jsonContent
	}

	// 写入文件（覆盖写入，不存在则创建）
	if err := os.WriteFile(absPath, content, 0666); err != nil {
		return fmt.Errorf("写入文件失败：%w", err)
	}

	fmt.Printf("数据已成功保存到：%s\n", absPath)
	return nil
}
