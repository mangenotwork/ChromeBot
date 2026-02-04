package builtins

import (
	"ChromeBot/dsl/interpreter"
	"fmt"
)

// 数学相关的内置方法
var mathFn = map[string]interpreter.Function{
	"abs": mathAbs, // 计算绝对值
	"max": mathMax, // 计算最大值
	"min": mathMin, // 计算最小值
}

func mathAbs(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs() 需要一个参数")
	}
	switch v := args[0].(type) {
	case int64:
		if v < 0 {
			return -v, nil
		}
		return v, nil
	case float64:
		if v < 0 {
			return -v, nil
		}
		return v, nil
	default:
		return nil, fmt.Errorf("abs() 不支持的类型: %T", args[0])
	}
}

func mathMax(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("max() 需要至少2个参数")
	}

	// todo 实现max函数
	return nil, nil
}

func mathMin(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("min() 需要至少2个参数")
	}

	// todo 实现min函数
	return nil, nil
}
