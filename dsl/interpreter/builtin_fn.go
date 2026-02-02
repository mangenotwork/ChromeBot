package interpreter

import (
	"fmt"
	"strings"
)

// 注册内置函数
func (i *Interpreter) registerBuiltins() {
	// 打印函数
	i.global.SetFunc("print", func(args []Value) (Value, error) {
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
	})

	i.global.SetFunc("println", func(args []Value) (Value, error) {
		return i.global.functions["print"](args)
	})

	// 类型转换
	i.global.SetFunc("int", func(args []Value) (Value, error) {
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
	})

	i.global.SetFunc("str", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("str() 需要一个参数")
		}
		return fmt.Sprintf("%v", args[0]), nil
	})

	// 数学函数
	i.global.SetFunc("len", func(args []Value) (Value, error) {
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
	})

	// 字典的 keys
	i.global.SetFunc("keys", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("keys() 需要一个参数")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("keys() 只支持字典，得到: %T", args[0])
		}

		// 获取所有键
		keys := make([]Value, 0, len(dict))
		for key := range dict {
			keys = append(keys, key)
		}
		return keys, nil
	})

	// 字典的 values
	i.global.SetFunc("values", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("values() 需要一个参数")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("values() 只支持字典，得到: %T", args[0])
		}

		// 获取所有值
		values := make([]Value, 0, len(dict))
		for _, value := range dict {
			values = append(values, value)
		}
		return values, nil
	})

	// 字典的 items
	i.global.SetFunc("items", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("items() 需要一个参数")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("items() 只支持字典，得到: %T", args[0])
		}

		// 获取所有键值对（每个键值对是一个包含两个元素的列表）
		items := make([]Value, 0, len(dict))
		for key, value := range dict {
			pair := []Value{key, value}
			items = append(items, pair)
		}
		return items, nil
	})

	// 字典的 has_key
	i.global.SetFunc("has_key", func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("has_key() 需要两个参数: dict, key")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("has_key() 第一个参数必须是字典，得到: %T", args[0])
		}

		// 检查键是否存在
		_, exists := dict[args[1]]
		return exists, nil
	})

	// 字典的 delete
	i.global.SetFunc("delete", func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("delete() 需要两个参数: dict, key")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("delete() 第一个参数必须是字典，得到: %T", args[0])
		}

		// 删除键
		delete(dict, args[1])
		return nil, nil
	})

	i.global.SetFunc("upper", func(args []Value) (Value, error) {
		if len(args) == 0 {
			return "", nil
		}

		// 将参数转换为字符串并转为大写
		str := fmt.Sprintf("%v", args[0])
		return strings.ToUpper(str), nil
	})

	i.global.SetFunc("repeat", func(args []Value) (Value, error) {
		if len(args) < 2 {
			return "", fmt.Errorf("repeat 需要两个参数: 字符串和次数")
		}

		str := fmt.Sprintf("%v", args[0])
		count, ok := args[1].(int64)
		if !ok {
			return "", fmt.Errorf("repeat 的第二个参数必须是整数")
		}

		return strings.Repeat(str, int(count)), nil
	})

}
