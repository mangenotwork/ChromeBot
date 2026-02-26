package builtins

import (
	"ChromeBot/dsl/interpreter"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
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
	"GBKToUTF8":       strGBKToUTF8,       // GBKToUTF8 函数 将GBK编码的字符串转换为utf-8编码
	"UTF8ToGBK":       strUTF8ToGBK,       // UTF8ToGBK 函数 将utf-8编码的字符串转换为GBK编码
	"UTF8ToGB2312":    strUTF8ToGB2312,    // UTF8ToGB2312 函数 将UTF-8转换为GB2312
	"GB2312ToUTF8":    strGB2312ToUTF8,    // GB2312ToUTF8 函数 将GB2312转换为UTF-8
	"UTF8ToGB18030":   strUTF8ToGB18030,   // UTF8ToGB18030 函数 将UTF-8转换为GB18030
	"GB18030ToUTF8":   strGB18030ToUTF8,   // GB18030ToUTF8 函数 将GB18030转换为UTF-8
	"UTF8ToBIG5":      strUTF8ToBIG5,      // UTF8ToBIG5 函数 将UTF-8转换为BIG5
	"BIG5ToUTF8":      strBIG5ToUTF8,      // BIG5ToUTF8 函数 将BIG5转换为UTF-8
	"UTF8ToLatin1":    strUTF8ToLatin1,    // UTF8ToLatin1 函数 将UTF-8转换为ISO-8859-1（Latin1）
	"Latin1ToUTF8":    strLatin1ToUTF8,    // Latin1ToUTF8 函数 将ISO-8859-1转换为UTF-8
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

// encodeString 将字符串从UTF-8编码转换为指定目标编码
func encodeString(str string, enc encoding.Encoding) string {
	if str == "" {
		return ""
	}
	encoder := enc.NewEncoder()
	result, _, err := transform.String(encoder, str)
	if err != nil {
		return str // 转换失败返回原字符串
	}
	return result
}

// decodeString 将指定编码的字符串转换为UTF-8编码
func decodeString(str string, enc encoding.Encoding) string {
	if str == "" {
		return ""
	}
	decoder := enc.NewDecoder()
	result, _, err := transform.String(decoder, str)
	if err != nil {
		return str // 转换失败返回原字符串
	}
	return result
}

func strGBKToUTF8(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("GBKToUTF8(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("GBKToUTF8(str) 需要字符串参数")
	}
	return decodeString(s, simplifiedchinese.GBK), nil
}

func strUTF8ToGBK(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UTF8ToGBK(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UTF8ToGBK(str) 需要字符串参数")
	}
	return encodeString(s, simplifiedchinese.GBK), nil
}

func strUTF8ToGB2312(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UTF8ToGB2312(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UTF8ToGB2312(str) 需要字符串参数")
	}
	return encodeString(s, simplifiedchinese.HZGB2312), nil
}

func strGB2312ToUTF8(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("GB2312ToUTF8(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("GB2312ToUTF8(str) 需要字符串参数")
	}
	return decodeString(s, simplifiedchinese.HZGB2312), nil
}

func strUTF8ToGB18030(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UTF8ToGB18030(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UTF8ToGB18030(str) 需要字符串参数")
	}
	return encodeString(s, simplifiedchinese.GB18030), nil
}

func strGB18030ToUTF8(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("GB18030ToUTF8(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("GB18030ToUTF8(str) 需要字符串参数")
	}
	return decodeString(s, simplifiedchinese.GB18030), nil
}

func strUTF8ToBIG5(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UTF8ToBIG5(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UTF8ToBIG5(str) 需要字符串参数")
	}
	return encodeString(s, traditionalchinese.Big5), nil
}

func strBIG5ToUTF8(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("BIG5ToUTF8(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("BIG5ToUTF8(str) 需要字符串参数")
	}
	return decodeString(s, traditionalchinese.Big5), nil
}

func strUTF8ToLatin1(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UTF8ToLatin1(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UTF8ToLatin1(str) 需要字符串参数")
	}
	return encodeString(s, charmap.ISO8859_1), nil
}

func strLatin1ToUTF8(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Latin1ToUTF8(str) 需要一个参数")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("Latin1ToUTF8(str) 需要字符串参数")
	}
	return decodeString(s, charmap.ISO8859_1), nil
}
