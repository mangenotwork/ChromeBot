package builtins

import (
	"ChromeBot/browser"
	"ChromeBot/dsl/interpreter"
	"fmt"
	"time"
)

var chromeFn = map[string]interpreter.Function{
	"ShowDemoTree":             chromeShowDemoTree,             // 显示当前demo树
	"MatchDemoContent":         chromeMatchDemoContent,         // 获取匹配到标签内容的xpath
	"MatchDemoContentOP":       chromeMatchDemoContentOP,       // 获取匹配到标签内容的xpath, 能用于操作的xpath
	"NowTabMatchDemoContentOP": chromeNowTabMatchDemoContentOP, // 获取当前操作的页面匹配到标签内容的xpath, 能用于操作的xpath
	"NowTabGetInputFirstXpath": chromeNowTabGetInputFirstXpath, // 获取当前操作的页面匹配到能输入的标签的xpath，返回匹配到的第一个
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
	chromeObj := browser.GetChromeInstance()
	htmlText, err := chromeObj.GetHtml()
	if err != nil {
		fmt.Println("[Chrome]未获取到当前页面的html:", err.Error())
		return "", err
	}

	xpath := browser.MatchDemoContentOP(htmlText, matchText)
	return xpath, nil
}

func chromeNowTabGetInputFirstXpath(args []interpreter.Value) (interpreter.Value, error) {
	// 获取当前页面
	chromeObj := browser.GetChromeInstance()
	htmlText, err := chromeObj.GetHtml()
	if err != nil {
		fmt.Println("[Chrome]未获取到当前页面的html:", err.Error())
		return "", err
	}

	xpath := browser.GetInputFirstXpath(htmlText)
	return xpath, nil
}
