package builtins

import (
	"ChromeBot/dsl/interpreter"
	"fmt"
	"time"
)

// 时间相关的内置方法
var timeFn = map[string]interpreter.Function{
	"now":   timeNow,   // now 获取当前时间的时间戳
	"sleep": timeSleep, // sleep 休眠
}

func timeNow(args []interpreter.Value) (interpreter.Value, error) {
	return time.Now().Unix(), nil
}

func timeSleep(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sleep() 需要一个参数")
	}
	var ms int64
	switch v := args[0].(type) {
	case int64:
		ms = v
	case float64:
		ms = int64(v)
	default:
		return nil, fmt.Errorf("sleep() 需要数字参数")
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil, nil
}
