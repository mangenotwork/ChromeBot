package builtins

import (
	"ChromeBot/dsl/interpreter"
	"fmt"
	"strings"
)

// 字符串相关的内置方法
var strFn = map[string]interpreter.Function{
	"upper":  strUpper,  // upper 将参数转换为字符串并转为大写
	"repeat": strRepeat, // repeat 将字符串进行重复, 第二个参数必须是整数
	"lower":  strLower,  // lower 字符串转小写
	"trim":   strTrim,   // trim 取首字符
	"split":  strSplit,  // split 字符分割
}

func strUpper(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) == 0 {
		return "", nil
	}
	str := fmt.Sprintf("%v", args[0])
	return strings.ToUpper(str), nil
}

func strRepeat(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("repeat 需要两个参数: 字符串和次数")
	}
	str := fmt.Sprintf("%v", args[0])
	count, ok := args[1].(int64)
	if !ok {
		return "", fmt.Errorf("repeat 的第二个参数必须是整数")
	}
	return strings.Repeat(str, int(count)), nil
}

func strLower(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("lower() 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("lower() 需要字符串参数")
	}
	return strings.ToLower(s), nil
}

func strTrim(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim() 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("trim() 需要字符串参数")
	}
	return strings.TrimSpace(s), nil
}

func strSplit(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split() 需要2个参数")
	}
	s, ok1 := args[0].(string)
	sep, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("split() 需要字符串参数")
	}
	parts := strings.Split(s, sep)
	result := make([]interpreter.Value, len(parts))
	for i, part := range parts {
		result[i] = part
	}
	return result, nil
}
