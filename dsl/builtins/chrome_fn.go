package builtins

import (
	"ChromeBot/browser"
	"ChromeBot/dsl/interpreter"
	"fmt"
	"time"

	gt "github.com/mangenotwork/gathertool"
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
	"HTMLGetPoint":             chromeHTMLGetPoint,             // HTMLGetPoint(html, label, attr, val)  获取指定位置的HTML， 用标签， 标签属性， 属性值来定位
	"HTMLGetPointID":           chromeHTMLGetPointID,           // HTMLGetPointID(html, label, val) 获取指定位置的HTML， 用标签， 标签属性为id， 属性值来定位
	"HTMLGetPointClass":        chromeHTMLGetPointClass,        // HTMLGetPointClass(html, label, val) 获取指定位置的HTML， 用标签， 标签属性为class， 属性值来定位
	"HtmlToTableSaveExcel":     chromeHtmlToTableSaveExcel,     // HtmlToTableSaveExcel(html, path, 可选参数sheetName) 提取html内的表格数据保存为Excel
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

	//fmt.Println("NowTabGetPointHTML 结果: ", res)

	resVal := make([]interpreter.Value, 0)
	for _, v := range res {
		resVal = append(resVal, interpreter.Value(v))
	}
	return resVal, err
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

	//fmt.Println("NowTabGetPointIDHTML 结果: ", res)

	resVal := make([]interpreter.Value, 0)
	for _, v := range res {
		resVal = append(resVal, v)
	}
	return resVal, err
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

	//fmt.Println("NowTabGetPointClassHTML 结果: ", res)

	resVal := make([]interpreter.Value, 0)
	for _, v := range res {
		resVal = append(resVal, interpreter.Value(v))
	}

	return resVal, err
}

func chromeHtmlToTableSaveExcel(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("HtmlToTableSaveExcel(html, path, 可选参数sheetName) 需要两个参数")
	}

	htmlStr, htmlStrOK := args[0].(string)
	if !htmlStrOK {
		return nil, fmt.Errorf("HtmlToTableSaveExcel(html, path, 可选参数sheetName) html 参数要求是字符串 ")
	}

	path, pathOK := args[1].(string)
	if !pathOK {
		return nil, fmt.Errorf("HtmlToTableSaveExcel(html, path, 可选参数sheetName) path 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	res := gt.RegHtmlTable(htmlStr)
	if len(res) == 0 {
		res = gt.RegHtmlTableOnly(htmlStr)
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("HtmlToTableSaveExcel(html, path, 可选参数sheetName) html 中没有table标签 ")
	}

	tableStr := res[0]

	inputData := make([]interpreter.Value, 0)
	trList := gt.RegHtmlTr(tableStr)
	for _, v := range trList {
		trList := gt.RegHtmlTr(v)
		for _, v := range trList {
			inputDataItem := make([]interpreter.Value, 0)
			tdList := gt.RegHtmlTdTxt(v)
			for _, item := range tdList {
				item := gt.RegDelHtml(item)
				inputDataItem = append(inputDataItem, item)
			}
			if len(inputDataItem) > 0 {
				inputData = append(inputData, inputDataItem)
			}
		}
	}

	// fmt.Println(" len(inputData) ", len(inputData))
	// fmt.Println(" len(inputData)[0] ", inputData[0])

	saveArgs := []interpreter.Value{path, inputData, sheetName}

	return excelSave(saveArgs)
}

func chromeHTMLGetPoint(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("HTMLGetPoint(html, label, attr, val) 需要四个参数")
	}

	htmlStr, htmlStrOK := args[0].(string)
	if !htmlStrOK {
		return nil, fmt.Errorf("HTMLGetPoint(html, label, attr, val) html 参数要求是字符串 ")
	}

	label, labelOK := args[1].(string)
	if !labelOK {
		return nil, fmt.Errorf("HTMLGetPoint(html, label, attr, val) label 参数要求是字符串 ")
	}

	attr, attrOK := args[2].(string)
	if !attrOK {
		return nil, fmt.Errorf("HTMLGetPoint(html, label, attr, val) attr 参数要求是字符串 ")
	}

	val, valOK := args[3].(string)
	if !valOK {
		return nil, fmt.Errorf("HTMLGetPoint(html, label, attr, val) val 参数要求是字符串 ")
	}

	res, err := gt.GetPointHTML(htmlStr, label, attr, val)
	if err != nil {
		fmt.Println("NowTabGetPointHTML 函数运行错误: ", err.Error())
	}

	resVal := make([]interpreter.Value, 0)
	for _, v := range res {
		resVal = append(resVal, interpreter.Value(v))
	}
	return resVal, err
}

func chromeHTMLGetPointID(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("HTMLGetPointID(html, label, val) 需要三个参数")
	}

	htmlStr, htmlStrOK := args[0].(string)
	if !htmlStrOK {
		return nil, fmt.Errorf("HTMLGetPointID(html, label, val) html 参数要求是字符串 ")
	}

	label, labelOK := args[1].(string)
	if !labelOK {
		return nil, fmt.Errorf("HTMLGetPointID(html, label, val) label 参数要求是字符串 ")
	}

	val, valOK := args[2].(string)
	if !valOK {
		return nil, fmt.Errorf("HTMLGetPointID(html, label, val) val 参数要求是字符串 ")
	}

	res, err := gt.GetPointIDHTML(htmlStr, label, val)
	if err != nil {
		fmt.Println("NowTabGetPointIDHTML 函数运行错误: ", err.Error())
	}

	resVal := make([]interpreter.Value, 0)
	for _, v := range res {
		resVal = append(resVal, v)
	}
	return resVal, err
}

func chromeHTMLGetPointClass(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("HTMLGetPointClass(html, label, val) 需要三个参数")
	}

	htmlStr, htmlStrOK := args[0].(string)
	if !htmlStrOK {
		return nil, fmt.Errorf("HTMLGetPointClass(html, label, val) html 参数要求是字符串 ")
	}

	label, labelOK := args[1].(string)
	if !labelOK {
		return nil, fmt.Errorf("HTMLGetPointClass(html, label, val) 参数要求是字符串 ")
	}

	val, valOK := args[2].(string)
	if !valOK {
		return nil, fmt.Errorf("HTMLGetPointClass(html, label, val) 参数要求是字符串 ")
	}

	res, err := gt.GetPointClassHTML(htmlStr, label, val)
	if err != nil {
		fmt.Println("NowTabGetPointClassHTML 函数运行错误: ", err.Error())
	}

	resVal := make([]interpreter.Value, 0)
	for _, v := range res {
		resVal = append(resVal, interpreter.Value(v))
	}
	return resVal, err
}
