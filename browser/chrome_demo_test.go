package browser

import (
	"fmt"
	"os"
	"testing"
)

func TestDemo(t *testing.T) {
	// 测试用的HTML内容（包含重复标签、多属性）
	testHTML := `
	<!DOCTYPE html>
	<html lang="zh-CN" class="root">
	<head>
		<meta charset="UTF-8">
		<title>测试页面</title>
	</head>
	<body id="main" data-test="123" style="color: red; font-size: 14px;">
		<!-- 这是测试注释 -->
		<div class="container" id="content">
			<p class="text" data-id="456">Hello, Go!</p>
			<p class="text" data-id="789">XPath Test</p>
			<input type="text" value="默认值" disabled>
		</div>
	</body>
	</html>
	`

	// 解析HTML为带XPath的DOM树
	domRoot, err := ParseHTMLToDOM(testHTML)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
		os.Exit(1)
	}

	// 打印DOM树（含XPath和属性）
	fmt.Println("解析后的DOM树（含XPath和所有属性）：")
	PrintDOM(domRoot, 0)
}
