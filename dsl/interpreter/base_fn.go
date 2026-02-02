package interpreter

import (
	"fmt"
	"reflect"
)

func (i *Interpreter) bool(v Value) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int64:
		return val != 0
	case float64:
		return val != 0
	case string:
		return val != ""
	default:
		return v != nil
	}
}

func (i *Interpreter) add(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l + r
		case float64:
			return float64(l) + r
		case string:
			return fmt.Sprintf("%d%s", l, r)
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l + float64(r)
		case float64:
			return l + r
		case string:
			return fmt.Sprintf("%f%s", l, r)
		}
	case string:
		return fmt.Sprintf("%s%v", l, right)

	case []Value: // 列表加法（连接）
		switch r := right.(type) {
		case []Value:
			// 连接两个列表
			result := make([]Value, len(l)+len(r))
			copy(result, l)
			copy(result[len(l):], r)
			return result
		default:
			// 将其他值添加到列表末尾
			result := make([]Value, len(l)+1)
			copy(result, l)
			result[len(l)] = r
			return result
		}

	case DictType: // 字典合并
		switch r := right.(type) {
		case DictType:
			// 合并两个字典
			result := make(DictType)
			// 先复制左边字典的所有元素
			for k, v := range l {
				result[k] = v
			}
			// 然后复制右边字典的所有元素（右边的会覆盖左边的）
			for k, v := range r {
				result[k] = v
			}
			return result
		}
	}

	i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T + %T", left, right))
	return nil
}

func (i *Interpreter) sub(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l - r
		case float64:
			return float64(l) - r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 - %T", right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l - float64(r)
		case float64:
			return l - r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: float64 - %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T - %T", left, right))
	}
	return nil
}

func (i *Interpreter) mul(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l * r
		case float64:
			return float64(l) * r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 * %T", right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l * float64(r)
		case float64:
			return l * r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: float64 * %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T * %T", left, right))
	}
	return nil
}

func (i *Interpreter) div(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return l / r
		case float64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return float64(l) / r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 / %T", right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return l / float64(r)
		case float64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return l / r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: float64 / %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T / %T", left, right))
	}
	return nil
}

func (i *Interpreter) mod(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("模零错误"))
				return nil
			}
			return l % r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 %% %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T %% %T", left, right))
	}
	return nil
}

func (i *Interpreter) equal(left, right Value) bool {
	// 如果是列表，需要深度比较
	if l, ok := left.([]Value); ok {
		if r, ok := right.([]Value); ok {
			if len(l) != len(r) {
				return false
			}
			for idx := 0; idx < len(l); idx++ {
				if !i.valuesEqual(l[idx], r[idx]) { // 使用辅助函数比较元素
					return false
				}
			}
			return true
		}
		return false
	}
	return reflect.DeepEqual(left, right)
}

// 添加辅助函数
func (i *Interpreter) valuesEqual(left, right Value) bool {
	// 递归处理嵌套列表和字典
	switch l := left.(type) {
	case []Value: // 列表
		if r, ok := right.([]Value); ok {
			if len(l) != len(r) {
				return false
			}
			for idx := 0; idx < len(l); idx++ {
				if !i.valuesEqual(l[idx], r[idx]) {
					return false
				}
			}
			return true
		}
		return false

	case DictType: // 字典
		if r, ok := right.(DictType); ok {
			if len(l) != len(r) {
				return false
			}
			for key, lVal := range l {
				if rVal, exists := r[key]; exists {
					if !i.valuesEqual(lVal, rVal) {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
		return false

	default:
		return reflect.DeepEqual(left, right)
	}
}
func (i *Interpreter) less(left, right Value) bool {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l < r
		case float64:
			return float64(l) < r
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l < float64(r)
		case float64:
			return l < r
		}
	case string:
		switch r := right.(type) {
		case string:
			return l < r
		}
	}

	i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T < %T", left, right))
	return false
}

func (i *Interpreter) greater(left, right Value) bool {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l > r
		case float64:
			return float64(l) > r
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l > float64(r)
		case float64:
			return l > r
		}
	case string:
		switch r := right.(type) {
		case string:
			return l > r
		}
	}

	i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T > %T", left, right))
	return false
}
