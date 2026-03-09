package browser

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

// DOMNode 扩展后的DOM节点结构体，包含XPath和所有属性
type DOMNode struct {
	// 节点类型：元素节点、文本节点、注释节点等
	Type html.NodeType
	// 标签名（仅元素节点有效）
	TagName string
	// 文本内容（文本节点、注释节点等有效）
	Content string
	// 属性列表，完整保留所有属性键值对
	Attributes map[string]string
	// 节点的XPath定位路径
	XPath string
	// 子节点
	Children []*DOMNode
}

// ParseHTMLToDOM 将HTML内容解析为带XPath的DOM树
func ParseHTMLToDOM(htmlContent string) (*DOMNode, error) {
	// 将HTML字符串转为io.Reader
	reader := bytes.NewReader([]byte(htmlContent))
	// 解析HTML得到根节点
	rootNode, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	// 递归转换原生节点为自定义DOM节点（初始XPath为空）
	domRoot := convertToDOMNode(rootNode, "")
	return domRoot, nil
}

// isWhitespaceTextNode 判断是否是空白文本节点（过滤换行/空格等无意义文本）
func isWhitespaceTextNode(node *html.Node) bool {
	if node.Type != html.TextNode {
		return false
	}
	return strings.TrimSpace(node.Data) == ""
}

// getNodeIndex 计算当前节点在同类型兄弟节点中的索引（从1开始，符合XPath规范）
func getNodeIndex(node *html.Node) int {
	if node == nil {
		return 1
	}

	index := 0 // 先从0开始计数，最后+1
	// 遍历所有前置兄弟节点（NextSibling方向才能完整遍历所有同层级节点）
	for sibling := node.Parent.FirstChild; sibling != nil; sibling = sibling.NextSibling {
		// 遇到当前节点则终止遍历
		if sibling == node {
			break
		}

		// 过滤空白文本节点（避免干扰计数）
		if isWhitespaceTextNode(sibling) {
			continue
		}

		// 元素节点按标签名分组，非元素节点按类型分组
		isSameType := false
		if node.Type == html.ElementNode && sibling.Type == html.ElementNode {
			isSameType = (sibling.Data == node.Data)
		} else if node.Type != html.ElementNode && sibling.Type == node.Type {
			isSameType = true
		}

		if isSameType {
			index++
		}
	}

	// XPath索引从1开始
	return index + 1
}

// buildXPath 为当前节点构建XPath路径
func buildXPath(parentXPath string, node *html.Node) string {
	if node == nil || node.Type == html.DocumentNode {
		return "/" // 文档根节点的XPath
	}

	// 获取当前节点的索引
	index := getNodeIndex(node)
	indexStr := "[" + strconv.Itoa(index) + "]"

	// 根据节点类型构建XPath片段
	var nodePath string
	switch node.Type {
	case html.ElementNode:
		nodePath = node.Data + indexStr // 元素节点：div[1]、p[2]
	case html.TextNode:
		// 空白文本节点不参与XPath计数（保持逻辑一致）
		if isWhitespaceTextNode(node) {
			nodePath = "text()" + "[0]" // 标记为无效索引，后续可过滤
		} else {
			nodePath = "text()" + indexStr
		}
	case html.CommentNode:
		nodePath = "comment()" + indexStr // 注释节点：comment()[1]
	default:
		nodePath = "*" + indexStr // 其他节点：*[1]
	}

	// 拼接父节点路径和当前节点路径
	if parentXPath == "/" || parentXPath == "" {
		return "/" + nodePath
	}
	return parentXPath + "/" + nodePath
}

// convertToDOMNode 递归转换原生html.Node到带XPath的自定义DOMNode
func convertToDOMNode(node *html.Node, parentXPath string) *DOMNode {
	if node == nil {
		return nil
	}

	// 过滤空白文本节点（不加入DOM树）
	if isWhitespaceTextNode(node) {
		return nil
	}

	// 构建当前节点的XPath
	xpath := buildXPath(parentXPath, node)

	// 初始化自定义节点
	domNode := &DOMNode{
		Type:       node.Type,
		TagName:    node.Data,
		Attributes: make(map[string]string),
		XPath:      xpath,
	}

	// 处理不同类型的节点内容
	switch node.Type {
	case html.TextNode:
		domNode.Content = strings.TrimSpace(node.Data) // 清理文本空白
	case html.ElementNode:
		// 完整保留所有属性键值对
		for _, attr := range node.Attr {
			domNode.Attributes[attr.Key] = attr.Val
		}
	case html.CommentNode:
		domNode.Content = node.Data
	}

	// 递归处理子节点
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		domChild := convertToDOMNode(child, xpath)
		if domChild != nil {
			domNode.Children = append(domNode.Children, domChild)
		}
	}

	return domNode
}

// PrintDOM 打印带XPath的DOM树（带缩进，便于查看）
func PrintDOM(node *DOMNode, indent int) {
	if node == nil {
		return
	}

	// 生成缩进
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	// 打印节点信息（包含XPath和属性）
	switch node.Type {
	case html.ElementNode:
		fmt.Printf("%s<%s> | XPath: %s | 属性: %v\n", indentStr, node.TagName, node.XPath, node.Attributes)
	case html.TextNode:
		// 过滤空文本（如换行、空格）
		if content := bytes.TrimSpace([]byte(node.Content)); len(content) > 0 {
			fmt.Printf("%s[文本] | XPath: %s | 内容: %s\n", indentStr, node.XPath, string(content))
		}
	case html.CommentNode:
		fmt.Printf("%s[注释] | XPath: %s | 内容: %s\n", indentStr, node.XPath, node.Content)
	default:
		fmt.Printf("%s[其他节点类型: %d] | XPath: %s | 名称: %s\n", indentStr, node.Type, node.XPath, node.TagName)
	}

	// 递归打印子节点
	for _, child := range node.Children {
		PrintDOM(child, indent+1)
	}
}

func MatchContentPrintDOM(node *DOMNode, text string) string {
	if node == nil {
		return ""
	}

	// 打印节点信息（包含XPath和属性）
	switch node.Type {
	case html.ElementNode:
		//fmt.Printf("%s<%s> | XPath: %s | 属性: %v\n", indentStr, node.TagName, node.XPath, node.Attributes)
	case html.TextNode:
		// 过滤空文本（如换行、空格）
		if content := bytes.TrimSpace([]byte(node.Content)); len(content) > 0 {
			if strings.Contains(string(content), text) {
				fmt.Printf("[文本] | XPath: %s | 内容: %s\n", node.XPath, string(content))
				return node.XPath
			}
		}
	case html.CommentNode:
		//fmt.Printf("%s[注释] | XPath: %s | 内容: %s\n", indentStr, node.XPath, node.Content)
	default:
		//fmt.Printf("%s[其他节点类型: %d] | XPath: %s | 名称: %s\n", indentStr, node.Type, node.XPath, node.TagName)
	}

	// 递归打印子节点
	for _, child := range node.Children {
		MatchContentPrintDOM(child, text)
	}

	return ""
}

func MatchContentDOM(node *DOMNode, text string) string {
	if node == nil {
		return ""
	}

	// 打印节点信息（包含XPath和属性）
	switch node.Type {
	case html.ElementNode:
		//fmt.Printf("%s<%s> | XPath: %s | 属性: %v\n", indentStr, node.TagName, node.XPath, node.Attributes)
	case html.TextNode:
		// 过滤空文本（如换行、空格）
		if content := bytes.TrimSpace([]byte(node.Content)); len(content) > 0 {
			if string(content) == text {
				fmt.Printf("[ MatchContentDOM 文本] | XPath: %s | 内容: %s\n", node.XPath, string(content))
				return node.XPath
			}
		}
	case html.CommentNode:
		//fmt.Printf("%s[注释] | XPath: %s | 内容: %s\n", indentStr, node.XPath, node.Content)
	default:
		//fmt.Printf("%s[其他节点类型: %d] | XPath: %s | 名称: %s\n", indentStr, node.Type, node.XPath, node.TagName)
	}

	// 递归打印子节点
	for _, child := range node.Children {
		xpath := MatchContentDOM(child, text)
		if xpath != "" {
			return xpath
		}
	}

	return ""
}

func MatchInput(node *DOMNode, res *[]string) {
	if node == nil {
		return
	}

	// 打印节点信息（包含XPath和属性）
	switch node.Type {
	case html.ElementNode:
		if node.TagName == "input" || node.TagName == "textarea" {
			// 忽略被隐藏的输入框
			if style, styleOK := node.Attributes["style"]; styleOK {
				if strings.Contains(style, "display:none") {
					return
				}
			}
			if inputType, inputTypeOK := node.Attributes["type"]; node.TagName == "input" && inputTypeOK {
				if strings.Contains(inputType, "hidden") {
					return
				}
			}
			// 忽略的一些不可输入的input type
			if inputType, inputTypeOK := node.Attributes["type"]; node.TagName == "input" && inputTypeOK {
				if inputType == "submit" || inputType == "reset" || inputType == "button" || inputType == "radio" ||
					inputType == "checkbox" || inputType == "file" {
					return
				}
			}

			*res = append(*res, node.XPath)

			fmt.Printf("<%s> | XPath: %s | 属性: %v\n", node.TagName, node.XPath, node.Attributes)
		}
	case html.TextNode:
		// 过滤空文本（如换行、空格）
		//if content := bytes.TrimSpace([]byte(node.Content)); len(content) > 0 {
		//	fmt.Printf("%s[文本] | XPath: %s | 内容: %s\n", indentStr, node.XPath, string(content))
		//}
	case html.CommentNode:
		//fmt.Printf("%s[注释] | XPath: %s | 内容: %s\n", indentStr, node.XPath, node.Content)
	default:
		//fmt.Printf("%s[其他节点类型: %d] | XPath: %s | 名称: %s\n", indentStr, node.Type, node.XPath, node.TagName)
	}

	// 递归打印子节点
	for _, child := range node.Children {
		MatchInput(child, res)
	}
}

// todo 循环匹配  可以通过内容  属性值  得出结果是 xpath
