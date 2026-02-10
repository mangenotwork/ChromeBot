package interpreter

import (
	"fmt"
	"strings"
)

// 注册内置函数
func (i *Interpreter) registerBuiltins() {

	builtinFnMap := map[string]Function{
		"print":   builtinPrint,  // print 打印函数
		"int":     builtinInt,    // int 类型转换 数值字符串转换数值类型
		"str":     builtinStr,    // str 类型转换 转换为字符串类型
		"len":     builtinLen,    // len 获取传入类型的长度，arg是任意类型，返回长度
		"keys":    builtinKeys,   // keys  获取字典的keys
		"values":  builtinValues, // values  获取字典的values
		"items":   builtinItems,  // items  获取所有键值对（每个键值对是一个包含两个元素的列表）
		"has":     builtinHas,    // has 字典或列表是否存在元素, arg第一个是字典或列表， 第二个是要找的元素
		"delete":  builtinDelete, // delete 删除字典或列表的指定元素, arg第一个是字典或列表， 第二个是要找的元素
		"type_of": builtinTypeOf, // type_of 获取变量类型
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
		// 检查参数类型，如果是字典类型，特殊处理
		if dict, ok := arg.(DictType); ok {
			fmt.Print(dictToString(dict))
		} else {
			fmt.Print(arg)
		}
	}
	fmt.Println()
	// 为了支持链式调用，返回最后一个参数
	// 如果没有参数，返回 nil
	if len(args) > 0 {
		return args[len(args)-1], nil
	}
	return nil, nil
}

func dictToString(d DictType) string {
	if d == nil {
		return "dict[]"
	}

	var sb strings.Builder
	sb.WriteString("dict[")

	first := true
	for key, value := range d {
		if !first {
			sb.WriteString(", ")
		}
		// 处理key
		sb.WriteString(fmt.Sprint(key))
		sb.WriteString(":")

		// 递归处理嵌套字典
		if nestedDict, ok := value.(DictType); ok {
			sb.WriteString(dictToString(nestedDict))
		} else {
			sb.WriteString(fmt.Sprint(value))
		}

		first = false
	}

	sb.WriteString("]")
	return sb.String()
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

	var exists Value

	switch args[0].(type) {
	case DictType:
		dict := args[0].(DictType)
		_, exists = dict[args[1]]
	case []Value:
		exists = false
		for _, v := range args[0].([]Value) {
			if v == args[1] {
				exists = true
			}
		}
	default:
		return nil, fmt.Errorf("has_key() 第一个参数必须是字典或者是列表，得到: %T", args[0])

	}

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

func builtinTypeOf(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type_of() 需要一个参数")
	}
	switch v := args[0].(type) {
	case int64:
		return "int", nil
	case float64:
		return "float", nil
	case string:
		return "string", nil
	case bool:
		return "bool", nil
	case []Value: // 优先匹配 []Value（如果切片元素类型是 Value）
		return "list", nil
	case []interface{}: // 匹配通用 []interface{} 切片
		return "list", nil
	case map[Value]Value: // 优先匹配 map[Value]Value（如果Map的键值都是 Value）
		return "dict", nil
	case map[interface{}]interface{}: // 匹配通用 map[interface{}]interface{}
		return "dict", nil
	case DictType:
		return "dict", nil
	default:
		switch fmt.Sprintf("%T", v) {
		case "[]int", "[]string", "[]bool": // 常见基础类型切片
			return "list", nil
		case "map[string]interface{}", "map[int]interface{}": // 常见基础类型键Map
			return "dict", nil
		default:
			return "unknown", nil
		}
	}
}
