package builtins

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/utils"
	"encoding/json"
	"fmt"

	gt "github.com/mangenotwork/gathertool"
)

var jsonFn = map[string]interpreter.Function{
	"jsonDict": jsonToDict, // jsonDict(str) 将json字符串转换成字典
	"json":     jsonStr,    // json(arg) 将字典转换成json字符串
	"jsonFind": jsonFind,   // jsonFind(str, find) 查找json字符串  find是查询节点 如： {a:[{b:1},{b:2}]}  find=/a/[0]  =>   {b:1}   find=a/[0]/b  =>  1
	"jsonIS":   jsonIS,     // jsonIS(str) 判断是否是json字符串
	"jsonSave": jsonSave,   // jsonSave(arg, path) 将变量转存到本地文件，数据内容为json(格式化输出)
}

func jsonToDict(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json(str) 需要一个参数")
	}

	fmt.Printf("%T", args[0])

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("json(str) 参数要求是字符串 ")
	}
	dataMap, err := gt.Json2Map(str)
	if err != nil {
		fmt.Println("[Err]不是合法的json字符串")
		return nil, fmt.Errorf("[Err]不是合法的json字符串")
	}

	dict := convertToDictType(dataMap)

	return dict, nil
}

func convertToDictType(data any) interpreter.DictType {
	dict := make(interpreter.DictType)

	// 先尝试将输入转为 map[string]any（顶层/嵌套 map 都走这个逻辑）
	dataMap, ok := data.(map[string]any)
	if !ok {
		return dict // 非 map 类型直接返回空 DictType
	}

	// 遍历每个键值对，递归处理
	for k, v := range dataMap {
		// 键转换：string -> Value (any)
		key := interpreter.Value(k)

		// 值转换：递归处理不同类型
		switch val := v.(type) {
		case map[string]any:
			// 嵌套 map，递归转换为 DictType
			dict[key] = convertToDictType(val)
		case []any:
			// 处理切片：递归转换切片内的每个元素
			var slice []interpreter.Value
			for _, elem := range val {
				// 如果切片内是 map，转 DictType；否则直接转为 Value
				if elemMap, ok := elem.(map[string]any); ok {
					slice = append(slice, convertToDictType(elemMap))
				} else {
					slice = append(slice, interpreter.Value(elem))
				}
			}
			dict[key] = slice
		default:
			// 基础类型（string/int/float/bool/null 等）直接转为 Value
			dict[key] = interpreter.Value(val)
		}
	}

	return dict
}

func jsonStr(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("jsonStr(dict) 需要一个参数")
	}

	str, err := json.Marshal(args[0])
	if err != nil {
		fmt.Println("[Err]变量转json字符串错误")
		return nil, fmt.Errorf("[Err]变量转json字符串错误")
	}

	return string(str), nil
}

func jsonFind(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("jsonFind(str, find) 需要一个参数")
	}

	str, strOK := args[0].(string)
	if !strOK {
		return nil, fmt.Errorf("jsonFind(str, find) str 参数要求是字符串 ")
	}

	find, findOK := args[1].(string)
	if !findOK {
		return nil, fmt.Errorf("jsonFind(str, find) find 参数要求是字符串 ")
	}

	rse, err := gt.JsonFind(str, find)
	if err != nil {
		fmt.Println("[Err]jsonFind失败，", err.Error())
		return nil, fmt.Errorf("[Err]jsonFind失败，%s", err.Error())
	}

	if rse == nil {
		fmt.Println("[Wring]jsonFind未寻找指定节点")
		return "", nil
	}

	rseMap, rseMapOK := rse.(map[string]any)
	if rseMapOK {
		return convertToDictType(rseMap), nil
	}

	return rse, nil
}

func jsonIS(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("jsonIS(str) 需要一个参数")
	}

	str, strOK := args[0].(string)
	if !strOK {
		return nil, fmt.Errorf("jsonFind(str, find) str 参数要求是字符串 ")
	}

	return gt.IsJson(str), nil
}

func jsonSave(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("jsonSave(arg, path)  需要两个参数")
	}
	jsonStr := ""
	var strOK bool = false
	jsonStr, strOK = args[0].(string)
	if strOK {

		if !gt.IsJson(jsonStr) {
			fmt.Println("[Err]不是json字符串")
			return false, nil
		}

	} else {

		jsosB, err := json.Marshal(args[0])
		if err != nil {
			fmt.Println("[Err]变量转json字符串错误")
			return nil, fmt.Errorf("[Err]变量转json字符串错误")
		}

		jsonStr = string(jsosB)
	}

	path, pathOK := args[1].(string)
	if !pathOK {
		return nil, fmt.Errorf("jsonSave(arg, path) path 参数要求是字符串 ")
	}

	jsonStr, _ = prettyJSON(jsonStr, "    ")

	err := utils.SaveDataToFile(path, jsonStr)
	if err != nil {
		fmt.Println("保存http请求到文件出现了错误:", err.Error())
	}

	return true, nil
}

func prettyJSON(input string, indent string) (string, error) {
	var data any
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", fmt.Errorf("JSON解析失败: %w", err)
	}
	prettyBytes, err := json.MarshalIndent(data, "", indent)
	if err != nil {
		return "", fmt.Errorf("JSON格式化失败: %w", err)
	}
	return string(prettyBytes), nil
}
