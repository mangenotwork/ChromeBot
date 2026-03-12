package builtins

import (
	"ChromeBot/browser"
	"ChromeBot/dsl/interpreter"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"time"
)

var chromeFn = map[string]interpreter.Function{
	"ShowDemoTree":             chromeShowDemoTree,             // 显示当前demo树
	"MatchDemoContent":         chromeMatchDemoContent,         // 获取匹配到标签内容的xpath
	"MatchDemoContentOP":       chromeMatchDemoContentOP,       // 获取匹配到标签内容的xpath, 能用于操作的xpath
	"NowTabMatchDemoContentOP": chromeNowTabMatchDemoContentOP, // 获取当前操作的页面匹配到标签内容的xpath, 能用于操作的xpath
	"NowTabGetInputFirstXpath": chromeNowTabGetInputFirstXpath, // 获取当前操作的页面匹配到能输入的标签的xpath，返回匹配到的第一个
	"NowTabGetPointHTML":       chromeNowTabGetPointHTML,       // NowTabGetPointHTML(label, attr, val)  获取指定位置的HTML， 用标签， 标签属性， 属性值来定位
	"NowTabGetPointIDHTML":     chromeNowTabGetPointIDHTML,     // NowTabGetPointIDHTML(label, val) 获取指定位置的HTML， 用标签， 标签属性为id， 属性值来定位
	"NowTabGetPointClassHTML":  chromeNowTabGetPointClassHTML,  // NowTabGetPointClassHTML(label, val) 获取指定位置的HTML， 用标签， 标签属性为class， 属性值来定位
}

func chromeShowDemoTree(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ShowDemoTree(html_text) 需要一个参数")
	}

	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("ShowDemoTree(html_text) 参数要求是html字符串 ")
	}

	browser.ShowDemoTree(str)

	return time.Now().Unix(), nil
}

func chromeMatchDemoContent(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("MatchDemoContent(html_text, match_text) 需要两个参数")
	}

	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("MatchDemoContent(html_text, match_text) 参数要求是html字符串 ")
	}

	str1, ok1 := args[1].(string)
	if !ok1 {
		return nil, fmt.Errorf("MatchDemoContent(html_text, match_text) 参数要求是html字符串 ")
	}

	xpath := browser.MatchDemoContent(str, str1)

	return xpath, nil
}

func chromeMatchDemoContentOP(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("MatchDemoContent(html_text, match_text) 需要两个参数")
	}

	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("MatchDemoContent(html_text, match_text) 参数要求是html字符串 ")
	}

	str1, ok1 := args[1].(string)
	if !ok1 {
		return nil, fmt.Errorf("MatchDemoContent(html_text, match_text) 参数要求是html字符串 ")
	}

	xpath := browser.MatchDemoContentOP(str, str1)

	return xpath, nil
}

func chromeNowTabMatchDemoContentOP(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("NowTabMatchDemoContentOP(match_text) 需要一个参数")
	}

	matchText, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("NowTabMatchDemoContentOP(match_text) 参数要求是字符串 ")
	}

	// 获取当前页面
	htmlText, err := browser.GetHtml()
	if err != nil {
		fmt.Println("[Chrome]未获取到当前页面的html:", err.Error())
		return "", err
	}

	xpath := browser.MatchDemoContentOP(htmlText, matchText)
	return xpath, nil
}

func chromeNowTabGetInputFirstXpath(args []interpreter.Value) (interpreter.Value, error) {
	// 获取当前页面
	htmlText, err := browser.GetHtml()
	if err != nil {
		fmt.Println("[Chrome]未获取到当前页面的html:", err.Error())
		return "", err
	}

	xpath := browser.GetInputFirstXpath(htmlText)
	return xpath, nil
}

func chromeNowTabGetPointHTML(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("NowTabGetPointHTML(label, attr, val) 需要三个参数")
	}

	label, labelOK := args[0].(string)
	if !labelOK {
		return nil, fmt.Errorf("NowTabGetPointHTML(label, attr, val) 参数要求是字符串 ")
	}

	attr, attrOK := args[1].(string)
	if !attrOK {
		return nil, fmt.Errorf("NowTabGetPointHTML(label, attr, val) 参数要求是字符串 ")
	}

	val, valOK := args[2].(string)
	if !valOK {
		return nil, fmt.Errorf("NowTabGetPointHTML(label, attr, val) 参数要求是字符串 ")
	}

	htmlText, err := browser.GetHtml()
	if err != nil {
		fmt.Println("[Chrome]未获取到当前页面的html:", err.Error())
		return "", err
	}

	res, err := gt.GetPointHTML(htmlText, label, attr, val)
	if err != nil {
		fmt.Println("NowTabGetPointHTML 函数运行错误: ", err.Error())
	}

	fmt.Println("NowTabGetPointHTML 结果: ", res)
	return res, err
}

func chromeNowTabGetPointIDHTML(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("NowTabGetPointIDHTML(label, val) 需要两个参数")
	}

	label, labelOK := args[0].(string)
	if !labelOK {
		return nil, fmt.Errorf("NowTabGetPointIDHTML(label, val) 参数要求是字符串 ")
	}

	val, valOK := args[1].(string)
	if !valOK {
		return nil, fmt.Errorf("NowTabGetPointIDHTML(label, val) 参数要求是字符串 ")
	}

	htmlText, err := browser.GetHtml()
	if err != nil {
		fmt.Println("[Chrome]未获取到当前页面的html:", err.Error())
		return "", err
	}

	res, err := gt.GetPointIDHTML(htmlText, label, val)
	if err != nil {
		fmt.Println("NowTabGetPointIDHTML 函数运行错误: ", err.Error())
	}

	fmt.Println("NowTabGetPointIDHTML 结果: ", res)
	return res, err
}

func chromeNowTabGetPointClassHTML(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("NowTabGetPointClassHTML(label, val) 需要两个参数")
	}

	label, labelOK := args[0].(string)
	if !labelOK {
		return nil, fmt.Errorf("NowTabGetPointClassHTML(label, val) 参数要求是字符串 ")
	}

	val, valOK := args[1].(string)
	if !valOK {
		return nil, fmt.Errorf("NowTabGetPointClassHTML(label, val) 参数要求是字符串 ")
	}

	htmlText, err := browser.GetHtml()
	if err != nil {
		fmt.Println("[Chrome]未获取到当前页面的html:", err.Error())
		return "", err
	}

	res, err := gt.GetPointClassHTML(htmlText, label, val)
	if err != nil {
		fmt.Println("NowTabGetPointClassHTML 函数运行错误: ", err.Error())
	}

	fmt.Println("NowTabGetPointClassHTML 结果: ", res)
	return res, err
}
