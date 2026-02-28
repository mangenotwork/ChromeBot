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

// 字符串相关的内置方法，字符串的正则方法
var strFn = map[string]interpreter.Function{
	"upper":           strUpper,           // upper 将参数转换为字符串并转为大写
	"repeat":          strRepeat,          // repeat 将字符串进行重复, 第二个参数必须是整数
	"lower":           strLower,           // lower 字符串转小写
	"trim":            strTrim,            // trim 取首字符
	"split":           strSplit,           // split 字符分割
	"replace":         strReplace,         // replace 字符串替换
	"replaceN":        strReplaceN,        // replaceN 字符串替换 指定替换几个
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
	"Reg":             strReg,             // Reg 函数 字符串正则 第一个参数是字符串，第二个参数是正则串
	"RegHtml":         strRegHtml,         // RegHtml 函数 用正则提取html 第一个参数是html字符串，第二个是标签
	"RegHtmlText":     strRegHtmlText,     // RegHtmlText 函数 用正则提取html只匹配标签内的文本部分 第一个参数是html字符串，第二个是标签名
	"RegFn":           strRegFn,           // RegFn 函数 内置了很多用正则提取的常用场景方法 第一个参数是字符串，第二个是方法名
	"RegDel":          strRegDel,          // RegDel 函数 常见的删除方法支持html删除指定标签内容 第一个参数是字符串，第二个是方法名或标签名
	"RegHas":          strRegHas,          // RegHas 函数 使用正则判断是否存在某内容 第一个参数是字符串，第二个是方法名
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

func strReplace(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace() 需要3个参数")
	}
	s1, ok1 := args[0].(string)
	s2, ok2 := args[1].(string)
	s3, ok3 := args[2].(string)
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("replace() 需要字符串参数")
	}
	result := strings.ReplaceAll(s1, s2, s3)
	return result, nil
}

func strReplaceN(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("replaceN() 需要4个参数")
	}
	s1, ok1 := args[0].(string)
	s2, ok2 := args[1].(string)
	s3, ok3 := args[2].(string)
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("replaceN() 前三个参数需要字符串")
	}
	n, err := getInt(args[3])
	if err != nil {
		return nil, fmt.Errorf("replaceN()最后一个参数%s", err.Error())
	}
	result := strings.Replace(s1, s2, s3, n)
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

func strReg(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Reg(str,str) 需要2个参数 ")
	}
	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("Reg(str,str) 第一个参数要求是字符串 ")
	}
	regStr, ok2 := args[1].(string)
	if !ok2 {
		return nil, fmt.Errorf("Reg(str,str) 第二个参数要求是字符串 ")
	}
	res := make([]interpreter.Value, 0)
	list := gt.RegFindAll(regStr, str)

	for _, v := range list {
		for _, vItem := range v {
			res = append(res, vItem)
		}
	}
	return res, nil
}

var labelFunc = map[string]func(str string, property ...string) []string{
	"a":           gt.RegHtmlA,
	"title":       gt.RegHtmlTitle,
	"keywords":    gt.RegHtmlKeyword,
	"description": gt.RegHtmlDescription,
	"tr":          gt.RegHtmlTr,
	"input":       gt.RegHtmlInput,
	"td":          gt.RegHtmlTd,
	"p":           gt.RegHtmlP,
	"span":        gt.RegHtmlSpan,
	"src":         gt.RegHtmlSrc,
	"href":        gt.RegHtmlHref,
	"h1":          gt.RegHtmlH1,
	"h2":          gt.RegHtmlH2,
	"h3":          gt.RegHtmlH3,
	"h4":          gt.RegHtmlH4,
	"h5":          gt.RegHtmlH5,
	"h6":          gt.RegHtmlH6,
	"tbody":       gt.RegHtmlTbody,
	"video":       gt.RegHtmlVideo,
	"canvas":      gt.RegHtmlCanvas,
	"code":        gt.RegHtmlCode,
	"img":         gt.RegHtmlImg,
	"ul":          gt.RegHtmlUl,
	"li":          gt.RegHtmlLi,
	"meta":        gt.RegHtmlMeta,
	"select":      gt.RegHtmlSelect,
	"table":       gt.RegHtmlTable,
	"button":      gt.RegHtmlButton,
	"tableOnly":   gt.RegHtmlTableOnly,
	"div":         gt.RegHtmlDiv,
	"option":      gt.RegHtmlOption,
}

func strRegHtml(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("RegHtml(str,str) 需要2个参数 ")
	}
	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("RegHtml(str,str) 第一个参数要求是字符串 ")
	}
	label, ok2 := args[1].(string)
	if !ok2 {
		return nil, fmt.Errorf("RegHtml(str,str) 第二个参数要求是字符串 ")
	}
	fn, ok3 := labelFunc[label]
	if !ok3 {
		return nil, fmt.Errorf("RegHtml(str,str) 暂时不支持%s该标签的正则 ", label)
	}
	res := make([]interpreter.Value, 0)
	list := fn(str)
	for _, v := range list {
		res = append(res, v)
	}
	return res, nil
}

var labelTextFunc = map[string]func(str string, property ...string) []string{
	"a":           gt.RegHtmlATxt,
	"title":       gt.RegHtmlTitleTxt,
	"keywords":    gt.RegHtmlKeywordTxt,
	"description": gt.RegHtmlDescriptionTxt,
	"tr":          gt.RegHtmlTrTxt,
	"td":          gt.RegHtmlTdTxt,
	"p":           gt.RegHtmlPTxt,
	"span":        gt.RegHtmlSpanTxt,
	"src":         gt.RegHtmlSrcTxt,
	"href":        gt.RegHtmlHrefTxt,
	"h1":          gt.RegHtmlH1Txt,
	"h2":          gt.RegHtmlH2Txt,
	"h3":          gt.RegHtmlH3Txt,
	"h4":          gt.RegHtmlH4Txt,
	"h5":          gt.RegHtmlH5Txt,
	"h6":          gt.RegHtmlH6Txt,
	"code":        gt.RegHtmlCodeTxt,
	"ul":          gt.RegHtmlUlTxt,
	"li":          gt.RegHtmlLiTxt,
	"select":      gt.RegHtmlSelectTxt,
	"table":       gt.RegHtmlTableTxt,
	"button":      gt.RegHtmlButtonTxt,
	"div":         gt.RegHtmlDivTxt,
	"option":      gt.RegHtmlOptionTxt,
	"value":       gt.RegValue,
}

func strRegHtmlText(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 需要2个参数 ")
	}
	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 第一个参数要求是字符串 ")
	}
	label, ok2 := args[1].(string)
	if !ok2 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 第二个参数要求是字符串 ")
	}
	fn, ok3 := labelTextFunc[label]
	if !ok3 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 暂时不支持%s该标签的正则 ", label)
	}
	res := make([]interpreter.Value, 0)
	list := fn(str)
	for _, v := range list {
		res = append(res, v)
	}
	return res, nil
}

var labelDelFunc = map[string]func(str string, property ...string) string{
	"html":   func(str string, property ...string) string { return gt.RegDelHtml(str) },
	"number": func(str string, property ...string) string { return gt.RegDelNumber(str) },
	"a":      func(str string, property ...string) string { return gt.RegDelHtmlA(str) },
	"title":  func(str string, property ...string) string { return gt.RegDelHtmlTitle(str) },
	"tr":     func(str string, property ...string) string { return gt.RegDelHtmlTr(str) },
	"input":  gt.RegDelHtmlInput,
	"td":     gt.RegDelHtmlTd,
	"p":      gt.RegDelHtmlP,
	"span":   gt.RegDelHtmlSpan,
	"src":    gt.RegDelHtmlSrc,
	"href":   gt.RegDelHtmlHref,
	"video":  gt.RegDelHtmlVideo,
	"canvas": gt.RegDelHtmlCanvas,
	"code":   gt.RegDelHtmlCode,
	"img":    gt.RegDelHtmlImg,
	"ul":     gt.RegDelHtmlUl,
	"li":     gt.RegDelHtmlLi,
	"meta":   gt.RegDelHtmlMeta,
	"select": gt.RegDelHtmlSelect,
	"table":  gt.RegDelHtmlTable,
	"button": gt.RegDelHtmlButton,
	"h1":     func(str string, property ...string) string { return gt.RegDelHtmlH(str, "1", property...) },
	"h2":     func(str string, property ...string) string { return gt.RegDelHtmlH(str, "2", property...) },
	"h3":     func(str string, property ...string) string { return gt.RegDelHtmlH(str, "3", property...) },
	"h4":     func(str string, property ...string) string { return gt.RegDelHtmlH(str, "4", property...) },
	"h5":     func(str string, property ...string) string { return gt.RegDelHtmlH(str, "5", property...) },
	"h6":     func(str string, property ...string) string { return gt.RegDelHtmlH(str, "6", property...) },
	"tbody":  gt.RegDelHtmlTbody,
}

func strRegDel(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("RegDel(str,str) 需要2个参数 ")
	}
	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("RegDel(str,str) 第一个参数要求是字符串 ")
	}
	label, ok2 := args[1].(string)
	if !ok2 {
		return nil, fmt.Errorf("RegDel(str,str) 第二个参数要求是字符串 ")
	}
	fn, ok3 := labelDelFunc[label]
	if !ok3 {
		return nil, fmt.Errorf("RegDel(str,str) 暂时不支持%s该标签的正则 ", label)
	}
	res := make([]interpreter.Value, 0)
	list := fn(str)
	for _, v := range list {
		res = append(res, v)
	}
	return res, nil
}

var labelHasFunc = map[string]func(str string, number int) bool{
	"IsNumber":        func(str string, number int) bool { return gt.IsNumber(str) },
	"IsNumber2Len":    gt.IsNumber2Len,
	"IsNumber2Heard":  gt.IsNumber2Heard,
	"IsFloat":         func(str string, number int) bool { return gt.IsFloat(str) },
	"IsFloat2Len":     gt.IsFloat2Len,
	"IsEngAll":        func(str string, number int) bool { return gt.IsEngAll(str) },
	"IsEngLen":        gt.IsEngLen,
	"IsEngNumber":     func(str string, number int) bool { return gt.IsEngNumber(str) },
	"IsLeastNumber":   gt.IsLeastNumber,
	"IsLeastCapital":  gt.IsLeastCapital,
	"IsLeastLower":    gt.IsLeastLower,
	"IsLeastSpecial":  gt.IsLeastSpecial,
	"HaveNumber":      func(str string, number int) bool { return gt.HaveNumber(str) },
	"HaveSpecial":     func(str string, number int) bool { return gt.HaveSpecial(str) },
	"IsEmail":         func(str string, number int) bool { return gt.IsEmail(str) },
	"IsDomain":        func(str string, number int) bool { return gt.IsDomain(str) },
	"IsURL":           func(str string, number int) bool { return gt.IsURL(str) },
	"IsPhone":         func(str string, number int) bool { return gt.IsPhone(str) },
	"IsLandline":      func(str string, number int) bool { return gt.IsLandline(str) },
	"AccountRational": func(str string, number int) bool { return gt.AccountRational(str) },
	"IsXMLFile":       func(str string, number int) bool { return gt.IsXMLFile(str) },
	"IsUUID3":         func(str string, number int) bool { return gt.IsUUID3(str) },
	"IsUUID4":         func(str string, number int) bool { return gt.IsUUID4(str) },
	"IsUUID5":         func(str string, number int) bool { return gt.IsUUID5(str) },
	"IsRGB":           func(str string, number int) bool { return gt.IsRGB(str) },
	"IsFullWidth":     func(str string, number int) bool { return gt.IsFullWidth(str) },
	"IsHalfWidth":     func(str string, number int) bool { return gt.IsHalfWidth(str) },
	"IsBase64":        func(str string, number int) bool { return gt.IsBase64(str) },
	"IsLatitude":      func(str string, number int) bool { return gt.IsLatitude(str) },
	"IsLongitude":     func(str string, number int) bool { return gt.IsLongitude(str) },
	"IsDNSName":       func(str string, number int) bool { return gt.IsDNSName(str) },
	"IsIPv4":          func(str string, number int) bool { return gt.IsWindowsPath(str) },
	"IsWindowsPath":   func(str string, number int) bool { return gt.IsWindowsPath(str) },
	"IsUnixPath":      func(str string, number int) bool { return gt.IsUnixPath(str) },
}

func strRegHas(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("RegHas(str,str) 需要2个参数 ")
	}
	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("RegHas(str,str) 第一个参数要求是字符串 ")
	}
	label, ok2 := args[1].(string)
	if !ok2 {
		return nil, fmt.Errorf("RegHas(str,str) 第二个参数要求是字符串 ")
	}
	fn, ok3 := labelHasFunc[label]
	if !ok3 {
		return nil, fmt.Errorf("RegHas(str,str) 暂时没有%s这个方法", label)
	}
	res := false
	if len(args) == 3 {
		number, ok4 := args[2].(int)
		if ok4 {
			res = fn(str, number)
			return res, nil
		}
	}

	res = fn(str, 1) // 默认1
	return res, nil
}

var labelFnFunc = map[string]func(str string, property ...string) []string{
	"RegTime":         gt.RegTime,
	"RegLink":         gt.RegLink,
	"RegEmail":        gt.RegEmail,
	"RegIPv4":         gt.RegIPv4,
	"RegIPv6":         gt.RegIPv6,
	"RegIP":           gt.RegIP,
	"RegMD5Hex":       gt.RegMD5Hex,
	"RegSHA1Hex":      gt.RegSHA1Hex,
	"RegSHA256Hex":    gt.RegSHA256Hex,
	"RegGUID":         gt.RegGUID,
	"RegMACAddress":   gt.RegMACAddress,
	"RegEmail2":       gt.RegEmail2,
	"RegUUID3":        gt.RegUUID3,
	"RegUUID4":        gt.RegUUID4,
	"RegUUID5":        gt.RegUUID5,
	"RegUUID":         gt.RegUUID,
	"RegInt":          gt.RegInt,
	"RegFloat":        gt.RegFloat,
	"RegRGBColor":     gt.RegRGBColor,
	"RegFullWidth":    gt.RegFullWidth,
	"RegHalfWidth":    gt.RegHalfWidth,
	"RegBase64":       gt.RegBase64,
	"RegLatitude":     gt.RegLatitude,
	"RegLongitude":    gt.RegLongitude,
	"RegDNSName":      gt.RegDNSName,
	"RegFullURL":      gt.RegFullURL,
	"RegURLSchema":    gt.RegURLSchema,
	"RegURLUsername":  gt.RegURLUsername,
	"RegURLPath":      gt.RegURLPath,
	"RegURLPort":      gt.RegURLPort,
	"RegURLIP":        gt.RegURLIP,
	"RegURLSubdomain": gt.RegURLSubdomain,
	"RegWinPath":      gt.RegWinPath,
	"RegUnixPath":     gt.RegUnixPath,
}

func strRegFn(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 需要2个参数 ")
	}
	str, ok1 := args[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 第一个参数要求是字符串 ")
	}
	label, ok2 := args[1].(string)
	if !ok2 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 第二个参数要求是字符串 ")
	}
	fn, ok3 := labelFnFunc[label]
	if !ok3 {
		return nil, fmt.Errorf("RegHtmlText(str,str) 暂时没有%s这个方法", label)
	}
	res := make([]interpreter.Value, 0)
	list := fn(str)
	for _, v := range list {
		res = append(res, v)
	}
	return res, nil
}
