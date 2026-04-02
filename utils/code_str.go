package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func EscapeQuotesInBackticks(input string) string {

	lines := strings.Split(input, "\n")
	reBacktick := regexp.MustCompile("`([^`]*)`")    // 匹配`包裹的内容
	reDoubleQuote := regexp.MustCompile(`"([^"]*)"`) // 匹配"包裹的内容
	processedLines := make([]string, len(lines))

	for i, line := range lines {

		// trimmedLine := strings.TrimLeft(line, " \t")
		// if trimmedLine == "" || !strings.HasPrefix(trimmedLine, "chrome") {
		// 	Debug("没 chrome 不需要处理 -> ", line)
		// 	processedLines[i] = line
		// 	continue
		// }

		// 处理反引号` `内的"→\"、'→\"（新增单引号转义）
		line = reBacktick.ReplaceAllStringFunc(line, func(match string) string {
			content := match[1 : len(match)-1] // 去掉首尾`
			// 先替换双引号，再替换单引号（顺序不影响）
			escaped := strings.ReplaceAll(content, `"`, `\"`)
			escaped = strings.ReplaceAll(escaped, `'`, `\"`) // 新增：单引号转\"
			return "`" + escaped + "`"
		})

		// 处理双引号""内的'→\'
		line = reDoubleQuote.ReplaceAllStringFunc(line, func(match string) string {
			content := match[1 : len(match)-1] // 去掉首尾"
			escaped := strings.ReplaceAll(content, `'`, `\'`)
			return `"` + escaped + `"`
		})

		processedLines[i] = line
	}

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

// UnescapeUnicode 将Unicode编码的字符串转为中文
func UnescapeUnicode(unicodeStr string) string {
	quotedStr := "\"" + unicodeStr + "\""
	res, err := strconv.Unquote(quotedStr)
	if err != nil {
		var jsonRes string
		err = json.Unmarshal([]byte(quotedStr), &jsonRes)
		if err != nil {
			Debugf("转义失败：%v", err)
			return unicodeStr
		}
		return jsonRes
	}
	return res
}

// ProcessCommandLine 处理命令行字符串，还原\ + 换行的长命令
// 参数: input - 带\换行的命令行字符串
// 返回: 处理后的字符串（\换行被替换为空格，其他换行保留）
func ProcessCommandLine(input string) string {
	// 按换行符拆分（兼容Windows(\r\n)和Linux(\n)）
	lines := strings.Split(input, "\n")
	var buf bytes.Buffer // 高效拼接字符串

	for i, line := range lines {
		// 剔除行末尾的空格/制表符/回车（只保留有效字符）
		trimmedLine := strings.TrimRight(line, " \t\r")
		lineLen := len(trimmedLine)

		// 判断是否以\结尾（且不是空行）
		if lineLen > 0 && trimmedLine[lineLen-1] == '\\' {
			// 统计末尾连续的\数量
			backslashCount := 0
			for j := lineLen - 1; j >= 0 && trimmedLine[j] == '\\'; j-- {
				backslashCount++
			}
			// 奇数个\：最后一个\用于换行，替换为空格；偶数个\：保留原\
			if backslashCount%2 == 1 {
				buf.WriteString(trimmedLine[:lineLen-backslashCount])
				buf.WriteString(strings.Repeat("\\", backslashCount-1))
				buf.WriteString(" ")
			} else {
				buf.WriteString(trimmedLine)
				buf.WriteString("\n")
			}
		} else {
			// 非\结尾的行，直接写入（最后一行不加多余换行）
			buf.WriteString(line)
			// 不是最后一行则保留换行符
			if i != len(lines)-1 {
				buf.WriteString("\n")
			}
		}
	}

	return buf.String()
}

// ProcessArgs 处理包含括号的参数数组，合并括号内的元素
func ProcessArgs(args []string) []string {
	// 空数组直接返回
	if len(args) == 0 {
		return args
	}

	var result []string
	i := 0
	length := len(args)

	for i < length {
		// 找到左括号的位置
		if args[i] == "(" {
			// 检查左括号前是否有元素（避免数组越界）
			if i == 0 {
				// 左括号开头，直接保留（异常场景）
				result = append(result, args[i])
				i++
				continue
			}

			// 查找对应的右括号
			rightIdx := -1
			for j := i + 1; j < length; j++ {
				if args[j] == ")" {
					rightIdx = j
					break
				}
			}

			if rightIdx == -1 {
				// 未找到右括号，按原格式保留
				result = append(result, args[i-1], args[i])
				i++
				continue
			}

			// 提取括号内的元素并拼接
			innerElements := args[i+1 : rightIdx]
			innerStr := strings.Join(innerElements, ",")

			// 拼接成完整字符串（前元素 + ( + 拼接内容 + )）
			fullStr := fmt.Sprintf("%s(%s)", args[i-1], innerStr)
			result = append(result, fullStr)

			// 跳过已处理的元素
			i = rightIdx + 1
		} else {
			if i+1 < length && args[i+1] == "(" {
				i++
			} else {
				// 单元素/非括号后无左括号 → 直接添加
				result = append(result, args[i])
				i++
			}
		}
	}

	return result
}

func SanitizeFileName(fullPath string) string {
	dir := filepath.Dir(fullPath)
	filename := filepath.Base(fullPath)
	// 正则匹配时间格式的冒号（HH:MM:SS 中的:），替换为-
	// 正则规则：匹配 数字:数字 或 数字:数字 的格式（避免误替换其他冒号）
	timeColonRegex := regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2})`)
	sanitizedFilename := timeColonRegex.ReplaceAllString(filename, "${1}-${2}-${3}")
	return filepath.Join(dir, sanitizedFilename)
}

func GetAbsolutePath(path string) (string, error) {
	// 清理路径
	clean := cleanPath(path)

	// 转换为绝对路径
	return filepath.Abs(clean)
}

func cleanPath(path string) string {
	// 移除首尾的引号
	path = strings.Trim(path, `"'`)

	// 如果是Windows路径，移除额外的转义
	if filepath.Separator == '\\' { // Windows
		path = strings.ReplaceAll(path, `\"`, `"`)
		path = strings.ReplaceAll(path, `\\`, `\`)
	}

	return path
}

func RemoveNewlinesInBackticks(input string) string {
	reBacktick := regexp.MustCompile(`(` + "`" + `)([\s\S]*?)(` + "`" + `)`)

	result := reBacktick.ReplaceAllStringFunc(input, func(match string) string {
		content := match[1 : len(match)-1]
		// 移除所有换行符（\n、\r\n、\r）
		noNewlines := strings.ReplaceAll(content, "\r\n", "")
		noNewlines = strings.ReplaceAll(noNewlines, "\n", "")
		noNewlines = strings.ReplaceAll(noNewlines, "\r", "")
		noNewlines = regexp.MustCompile(`\s+`).ReplaceAllString(noNewlines, " ")
		noNewlines = strings.TrimSpace(noNewlines)
		return "`" + noNewlines + "`"
	})

	return result
}

func ShowJson(data any) {
	if data == nil {
		fmt.Println("{}")
		return
	}

	// 格式化缩进输出，两个空格缩进
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("JSON 格式化失败: %v\n", err)
		return
	}

	fmt.Println(string(bytes))
}

// JsonPrettyFormat 输入原始 JSON 字符串，返回美化格式化后的字符串
func JsonPrettyFormat(jsonStr string) string {
	// 如果已经是格式化的 JSON，直接返回
	if strings.TrimSpace(jsonStr) == "" {
		return jsonStr
	}

	// 尝试解析
	var data interface{}

	// 创建自定义解码器
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.UseNumber() // 保留数字精度

	if err := decoder.Decode(&data); err != nil {
		// 如果是无效 JSON，返回原字符串
		fmt.Printf("[Wring] 无法解析为 JSON: %v", err)
		return jsonStr
	}

	// 创建自定义编码器
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // 关键：不转义 Unicode
	encoder.SetIndent("", "  ")  // 缩进2个空格

	if err := encoder.Encode(data); err != nil {
		fmt.Printf("[Wring] JSON 编码失败: %v", err)
		return jsonStr
	}

	// 移除末尾的换行符（Encode 会自动添加）
	result := buf.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result
}

// IsValidOrigin 验证 origin 格式
func IsValidOrigin(origin string) bool {
	// 基本格式验证
	if !strings.Contains(origin, "://") {
		return false
	}

	// 可以解析 URL 进行更严格的验证
	if parsed, err := url.Parse(origin); err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}

	return true
}
