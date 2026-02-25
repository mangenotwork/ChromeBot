package builtins

import (
	"ChromeBot/dsl/interpreter"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"strings"
)

// 字符串相关的内置方法
var strFn = map[string]interpreter.Function{
	"upper":           strUpper,           // upper 将参数转换为字符串并转为大写
	"repeat":          strRepeat,          // repeat 将字符串进行重复, 第二个参数必须是整数
	"lower":           strLower,           // lower 字符串转小写
	"trim":            strTrim,            // trim 取首字符
	"split":           strSplit,           // split 字符分割
	"CleanWhitespace": strCleanWhitespace, // CleanWhitespace 函数 清理字符串回车，换行符号，还有前后空格
	"StrDeleteSpace":  strDeleteSpace,     // StrDeleteSpace 函数 删除字符串前后的空格
	"UnicodeDecode":   strUnicodeDecode,   // UnicodeDec 函数 字符串进行unicode编码
	"UnescapeUnicode": strUnescapeUnicode, // UnescapeUnicode 函数 字符串进行unicode解码
	"Base64Encode":    strBase64Encode,    // Base64Encode 函数 字符串进行base64编码
	"Base64Decode":    strBase64Decode,    // Base64Decode 函数 字符串进行base64解码
	"UrlBase64Encode": strUrlBase64Encode, // UrlBase64Encode 函数 url进行base64编码
	"UrlBase64Decode": strUrlBase64Decode, // UrlBase64Decode 函数 url进行base64解码
	"MD5":             strMD5,             // MD5 函数 将字符串进行md5
	"MD516":           strMD516,           // MD516 函数 将字符串进行md5，返回16位
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

func strCleanWhitespace(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("CleanWhitespace(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("CleanWhitespace(str) 需要字符串参数")
	}
	return gt.CleaningStr(s), nil
}

func strDeleteSpace(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("StrDeleteSpace(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("StrDeleteSpace(str) 需要字符串参数")
	}
	return gt.StrDeleteSpace(s), nil
}

func strUnicodeDecode(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UnicodeDecode(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UnicodeDecode(str) 需要字符串参数")
	}
	return gt.UnicodeDec(s), nil
}

func strUnescapeUnicode(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UnicodeDec(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UnicodeDec(str) 需要字符串参数")
	}
	rse, err := gt.UnescapeUnicode([]byte(s))
	if err != nil {
		return nil, err
	}
	return string(rse), nil
}

func strBase64Encode(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Base64Encode(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("Base64Encode(str) 需要字符串参数")
	}
	return gt.Base64Encode(s), nil
}

func strBase64Decode(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Base64Decode(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("Base64Decode(str) 需要字符串参数")
	}
	rse, err := gt.Base64Decode(s)
	if err != nil {
		return nil, err
	}
	return rse, nil
}

func strUrlBase64Encode(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UrlBase64Encode(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UrlBase64Encode(str) 需要字符串参数")
	}
	return gt.Base64UrlEncode(s), nil
}

func strUrlBase64Decode(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UrlBase64Decode(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UrlBase64Decode(str) 需要字符串参数")
	}
	rse, err := gt.Base64UrlDecode(s)
	if err != nil {
		return nil, err
	}
	return rse, nil
}

func strMD5(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MD5(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("MD5(str) 需要字符串参数")
	}
	return gt.GetMD5Encode(s), nil
}

func strMD516(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MD516(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("MD516(str) 需要字符串参数")
	}
	return gt.Get16MD5Encode(s), nil
}
