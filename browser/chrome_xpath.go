package browser

import (
	"fmt"
	"log"
	"regexp"
)

func ShowDemoTree(htmlText string) {
	domRoot, err := ParseHTMLToDOM(htmlText)
	if err != nil {
		fmt.Printf("html解析失败: %v\n", err)
		return
	}

	// 打印DOM树（含XPath和属性）
	fmt.Println("解析后的DOM树（含XPath和所有属性）：")
	PrintDOM(domRoot, 0)
}

func MatchDemoContent(htmlText, contentText string) string {
	domRoot, err := ParseHTMLToDOM(htmlText)
	if err != nil {
		fmt.Printf("html解析失败: %v\n", err)
		return ""
	}

	// 打印DOM树（含XPath和属性）
	fmt.Println("解析后的DOM树（含XPath和所有属性）：")
	return MatchContentPrintDOM(domRoot, contentText)
}

// MatchDemoContentOP 匹配标签内容获取到可交互的xpath
func MatchDemoContentOP(htmlText, contentText string) string {
	domRoot, err := ParseHTMLToDOM(htmlText)
	if err != nil {
		fmt.Printf("html解析失败: %v\n", err)
		return ""
	}

	// 打印DOM树（含XPath和属性）
	fmt.Println("MatchDemoContentOP  解析后的DOM树（含XPath和所有属性）：")
	xpath := MatchContentDOM(domRoot, contentText)
	fmt.Println("MatchDemoContentOP 匹配到的xpath = ", xpath)
	return RemoveNodeSuffix(xpath)
}

// RemoveNodeSuffix 移除XPath末尾的节点类型后缀（text()/comment()/*），保留元素节点路径
// 示例：
// 输入: /html[1]/body[1]/button[1]/text()[1]
// 输出: /html[1]/body[1]/button[1]
func RemoveNodeSuffix(xpath string) string {
	if xpath == "" || xpath == "/" {
		return xpath
	}

	// 循环移除末尾的非元素节点后缀，直到遇到元素节点路径
	regex := regexp.MustCompile(`(/(text\(\)|comment\(\)|\*)+\[\d+\])$`)
	cleanedXPath := xpath
	for {
		// 检查是否还有需要移除的后缀
		if !regex.MatchString(cleanedXPath) {
			break
		}
		cleanedXPath = regex.ReplaceAllString(cleanedXPath, "")
	}

	// 兜底处理空路径
	if cleanedXPath == "" {
		return "/"
	}
	return cleanedXPath
}

// GetInputFirstXpath  获取页面第一个能输入的标签的xpath， 只支持寻找 input,textarea 标签
// 排除条件，隐藏的和已经被赋予了值的会被排除
// todo: 支持 div contenteditable="true"
func GetInputFirstXpath(htmlText string) string {
	domRoot, err := ParseHTMLToDOM(htmlText)
	if err != nil {
		fmt.Printf("html解析失败: %v\n", err)
		return ""
	}
	xpathList := make([]string, 0)
	MatchInput(domRoot, &xpathList)
	log.Println("MatchInput xpathList len = ", len(xpathList))
	if len(xpathList) < 1 {
		return ""
	}
	return xpathList[0]
}
