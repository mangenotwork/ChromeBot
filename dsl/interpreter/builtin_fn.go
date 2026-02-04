package interpreter

import (
	"fmt"
	"strings"
)

// 注册内置函数
func (i *Interpreter) registerBuiltins() {

	builtinFnMap := map[string]Function{
		"print":  builtinPrint,  // print 打印函数
		"int":    builtinInt,    // int 类型转换 数值字符串转换数值类型
		"str":    builtinStr,    // str 类型转换 转换为字符串类型
		"len":    builtinLen,    // len 获取传入类型的长度，arg是任意类型，返回长度
		"keys":   builtinKeys,   // keys  获取字典的keys
		"values": builtinValues, // values  获取字典的values
		"items":  builtinItems,  // items  获取所有键值对（每个键值对是一个包含两个元素的列表）
		"has":    builtinHas,    // has 字典或列表是否存在元素, arg第一个是字典或列表， 第二个是要找的元素
		"delete": builtinDelete, // delete 删除字典或列表的指定元素, arg第一个是字典或列表， 第二个是要找的元素
		"upper":  builtinUpper,  // upper 将参数转换为字符串并转为大写
		"repeat": builtinRepeat, // repeat 将字符串进行重复, 第二个参数必须是整数
	}
	for name, fn := range builtinFnMap {
		i.global.SetFunc(name, fn)
	}
}

func builtinPrint(args []Value) (Value, error) {
	if len(args) == 0 {
		fmt.Println()
		return nil, nil
	}
	// 打印所有参数
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Println()
	// 为了支持链式调用，返回最后一个参数
	// 如果没有参数，返回 nil
	if len(args) > 0 {
		return args[len(args)-1], nil
	}
	return nil, nil
}

func builtinInt(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int() 需要一个参数")
	}
	switch v := args[0].(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		var result int64
		_, err := fmt.Sscanf(v, "%d", &result)
		if err != nil {
			return nil, fmt.Errorf("无法转换字符串为int: %s", v)
		}
		return result, nil
	case bool:
		if v {
			return int64(1), nil
		}
		return int64(0), nil
	default:
		return nil, fmt.Errorf("无法转换为int: %T", args[0])
	}
}

func builtinStr(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("str() 需要一个参数")
	}
	return fmt.Sprintf("%v", args[0]), nil
}

func builtinLen(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len() 需要一个参数")
	}
	switch v := args[0].(type) {
	case string:
		return int64(len(v)), nil
	case []Value: // 添加对列表的支持
		return int64(len(v)), nil
	case DictType: // 字典
		return int64(len(v)), nil
	default:
		return nil, fmt.Errorf("len() 不支持的类型: %T", args[0])
	}
}

func builtinKeys(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("keys() 需要一个参数")
	}
	dict, ok := args[0].(DictType)
	if !ok {
		return nil, fmt.Errorf("keys() 只支持字典，得到: %T", args[0])
	}
	keys := make([]Value, 0, len(dict))
	for key := range dict {
		keys = append(keys, key)
	}
	return keys, nil
}

func builtinValues(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("values() 需要一个参数")
	}
	dict, ok := args[0].(DictType)
	if !ok {
		return nil, fmt.Errorf("values() 只支持字典，得到: %T", args[0])
	}
	values := make([]Value, 0, len(dict))
	for _, value := range dict {
		values = append(values, value)
	}
	return values, nil
}

func builtinItems(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("items() 需要一个参数")
	}
	dict, ok := args[0].(DictType)
	if !ok {
		return nil, fmt.Errorf("items() 只支持字典，得到: %T", args[0])
	}
	items := make([]Value, 0, len(dict))
	for key, value := range dict {
		pair := []Value{key, value}
		items = append(items, pair)
	}
	return items, nil
}

func builtinHas(args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("has_key() 需要两个参数: dict, key")
	}
	dict, ok := args[0].(DictType)
	if !ok {
		return nil, fmt.Errorf("has_key() 第一个参数必须是字典，得到: %T", args[0])
	}

	// todo List类型也要支持

	_, exists := dict[args[1]]
	return exists, nil
}

func builtinDelete(args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("delete() 需要两个参数: dict, key")
	}
	dict, ok := args[0].(DictType)
	if !ok {
		return nil, fmt.Errorf("delete() 第一个参数必须是字典，得到: %T", args[0])
	}

	// todo List类型也要支持

	delete(dict, args[1])
	return nil, nil
}

func builtinUpper(args []Value) (Value, error) {
	if len(args) == 0 {
		return "", nil
	}
	str := fmt.Sprintf("%v", args[0])
	return strings.ToUpper(str), nil
}

func builtinRepeat(args []Value) (Value, error) {
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
