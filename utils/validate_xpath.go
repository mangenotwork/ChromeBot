package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ValidateXPathPureNative 校验XPath语法有效性
func ValidateXPathPureNative(xpathExpr string) (bool, error) {
	// 1. 空表达式校验
	trimmedExpr := strings.TrimSpace(xpathExpr)
	if trimmedExpr == "" {
		return false, errors.New("xpath表达式不能为空")
	}

	// 2. 执行核心语法校验
	if err := checkBracketMatching(trimmedExpr); err != nil {
		return false, fmt.Errorf("[xpath检查] 括号语法错误: %w", err)
	}
	if err := checkQuoteMatching(trimmedExpr); err != nil {
		return false, fmt.Errorf("[xpath检查] 引号语法错误: %w", err)
	}
	if err := checkAttributeSyntax(trimmedExpr); err != nil {
		return false, fmt.Errorf("[xpath检查] 属性语法错误: %w", err)
	}
	if err := checkInvalidPatterns(trimmedExpr); err != nil {
		return false, fmt.Errorf("[xpath检查] 无效语法模式: %w", err)
	}
	if err := checkBasicKeywords(trimmedExpr); err != nil {
		return false, fmt.Errorf("[xpath检查] 关键字语法错误: %w", err)
	}

	// 所有核心规则校验通过，判定为语法有效
	return true, nil
}

// checkBracketMatching 校验括号匹配（()、[]、{}）
func checkBracketMatching(expr string) error {
	bracketMap := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []rune{}

	for idx, char := range expr {
		// 跳过引号内的括号（引号内的括号不参与匹配）
		if inQuote(expr, idx) {
			continue
		}

		// 左括号入栈
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			// 右括号匹配校验
			if len(stack) == 0 {
				return fmt.Errorf("位置%d: 多余的闭合符号%c", idx, char)
			}
			last := stack[len(stack)-1]
			if last != bracketMap[char] {
				return fmt.Errorf("位置%d: 括号不匹配，期望%c但找到%c", idx, bracketMap[char], last)
			}
			stack = stack[:len(stack)-1]
		}
	}

	// 检查未闭合的括号
	if len(stack) > 0 {
		return fmt.Errorf("未闭合的括号: %c", stack[len(stack)-1])
	}
	return nil
}

// checkQuoteMatching 校验引号匹配（单/双引号）
func checkQuoteMatching(expr string) error {
	var quoteChar rune = 0 // 0表示未处于引号中
	for idx, char := range expr {
		switch char {
		case '\'', '"':
			if quoteChar == 0 {
				quoteChar = char // 开始引号
			} else if quoteChar == char {
				quoteChar = 0 // 结束引号
			}
		// 转义引号不中断匹配（如 'it\'s' 是合法的）
		case '\\':
			idx++ // 跳过转义字符后的字符
		}
	}

	if quoteChar != 0 {
		return fmt.Errorf("未闭合的%s引号", map[rune]string{'\'': "单", '"': "双"}[quoteChar])
	}
	return nil
}

// checkAttributeSyntax 校验@属性语法（如 @class 合法，@123/@ 不合法）
func checkAttributeSyntax(expr string) error {
	attrParts := strings.Split(expr, "@")
	for i := 1; i < len(attrParts); i++ {
		part := attrParts[i]
		// 跳过引号内的@（如 'a@b' 不是属性）
		if inQuote(expr, strings.Index(expr, "@"+part)) {
			continue
		}

		// 修剪空白字符
		trimmed := strings.TrimLeft(part, " \t\n\r")
		if len(trimmed) == 0 {
			return errors.New("@后未指定属性名")
		}

		// 属性名首字符必须是字母/下划线（XPath规范）
		firstChar := rune(trimmed[0])
		if !unicode.IsLetter(firstChar) && firstChar != '_' {
			return fmt.Errorf("属性名首字符非法: %c（必须是字母/下划线）", firstChar)
		}
	}
	return nil
}

// checkInvalidPatterns 校验无效语法模式（如空括号 []、()）
func checkInvalidPatterns(expr string) error {
	invalidPatterns := []string{"[]", "()", "{}", "[ ]", "( )", "{ }"}
	for _, pattern := range invalidPatterns {
		if strings.Contains(expr, pattern) {
			return fmt.Errorf("包含无效空括号模式: %s", pattern)
		}
	}
	return nil
}

// checkBasicKeywords 校验基础XPath关键字语法（可选，增强校验）
func checkBasicKeywords(expr string) error {
	// 常见XPath函数/轴关键字，校验格式（如 contains( 必须有参数）
	keywords := []string{"contains(", "text(", "node(", "ancestor::", "descendant::", "child::"}
	for _, kw := range keywords {
		idx := strings.Index(expr, kw)
		if idx == -1 {
			continue
		}
		// 跳过引号内的关键字
		if inQuote(expr, idx) {
			continue
		}

		// 检查关键字后是否有内容
		after := strings.TrimSpace(expr[idx+len(kw):])
		if len(after) == 0 || after[0] == ')' || after[0] == ']' {
			return fmt.Errorf("关键字%s后无有效参数", kw)
		}
	}
	return nil
}

// inQuote 辅助函数：判断指定位置的字符是否在引号内
func inQuote(expr string, pos int) bool {
	if pos < 0 || pos >= len(expr) {
		return false
	}

	var quoteChar rune = 0
	for i := 0; i < pos; i++ {
		char := rune(expr[i])
		switch char {
		case '\'', '"':
			if quoteChar == 0 {
				quoteChar = char
			} else if quoteChar == char {
				quoteChar = 0
			}
		case '\\':
			i++ // 跳过转义字符
		}
	}
	return quoteChar != 0
}
