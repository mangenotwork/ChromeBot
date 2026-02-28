package builtins

import (
	"ChromeBot/dsl/ast"
	"ChromeBot/dsl/interpreter"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"regexp"
	"strings"
	"time"
)

// 时间相关的内置方法
var timeFn = map[string]interpreter.Function{
	"now":                 timeNow,                 // now 获取当前时间的时间戳
	"sleep":               timeSleep,               // sleep 休眠
	"Timestamp":           timeNow,                 // timestamp 时间戳
	"TimestampMilli":      timeNowMilli,            // timestamp 时间戳 milliseconds
	"date":                timeDate,                // 获取日期
	"TimestampToDate":     timeTimestampToDate,     // 时间戳转日期
	"TimestampToDateAT":   timeTimestampToDateAT,   // 指定时间格式  YYYYMMDD YYYY-MM-DD YYYYMMDDHHmmss YYYY-MM-DD HH:mm:ss MMdd HHmmss
	"BeginDayUnix":        timeBeginDayUnix,        // 获取当天0点的时间戳
	"EndDayUnix":          timeEndDayUnix,          // 获取当天24点的时间戳
	"MinuteAgo":           timeMinuteAgo,           // 获取多少分钟前的时间戳
	"HourAgo":             timeHourAgo,             // 获取多少小时前的时间戳
	"DayAgo":              timeDayAgo,              // 获取多少天前的时间戳
	"DayDiffAtUnix":       timeDayDiffAtUnix,       // 两个时间戳的插值
	"DayDiff":             timeDayDiff,             // 两个时间字符串的日期差, 返回的是天 格式是 YYYY-MM-DD HH:mm:ss
	"NowToEnd":            timeNowToEnd,            // 计算当前时间到这天结束还有多久,单位秒
	"IsToday":             timeIsToday,             // 判断时间戳是否是今天，返回今天的时分秒
	"Timestamp2Week":      timeTimestamp2Week,      // 传入的时间戳是周几
	"Timestamp2WeekXinQi": timeTimestamp2WeekXinQi, // 传入的时间戳是星期几
}

func timeNow(args []interpreter.Value) (interpreter.Value, error) {
	return time.Now().Unix(), nil
}

func timeSleep(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sleep(int) 需要一个参数")
	}
	ms, err := getInt64(args[0])
	if err != nil {
		return nil, fmt.Errorf("sleep(int) %s", err.Error())
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil, nil
}

// milliseconds
func timeNowMilli(args []interpreter.Value) (interpreter.Value, error) {
	return time.Now().UnixMilli(), nil
}

func timeDate(args []interpreter.Value) (interpreter.Value, error) {
	return time.Now().Format(gt.TimeTemplate), nil
}

func getInt(arg interpreter.Value) (int, error) {
	var value int
	switch v := arg.(type) {
	case int:
		value = v
	case int64:
		value = int(v)
	case float64:
		value = int(v)
	case ast.Integer:
		value = int(v.Value)
	default:
		return 0, fmt.Errorf("需要数字参数")
	}
	return value, nil
}

func getInt64(arg interpreter.Value) (int64, error) {
	var value int64
	switch v := arg.(type) {
	case int:
		value = int64(v)
	case int64:
		value = v
	case float64:
		value = int64(v)
	case ast.Integer:
		value = v.Value
	default:
		return 0, fmt.Errorf("需要数字参数")
	}
	return value, nil
}

func timeTimestampToDate(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("TimestampToDate(int) 需要一个参数")
	}
	ms, err := getInt64(args[0])
	if err != nil {
		return nil, fmt.Errorf("TimestampToDate(int) %s", err.Error())
	}
	return gt.Timestamp2Date(ms), nil
}

func timeTimestampToDateAT(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("TimestampToDateAT(timestamp,template) 需要两个参数")
	}
	ms, err := getInt64(args[0])
	if err != nil {
		return nil, fmt.Errorf("TimestampToDateAT(timestamp,template) %s", err.Error())
	}
	template, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("TimestampToDateAT(timestamp,template) 第二个参数要求是字符串 ")
	}
	format, err := getGoTimeFormat(template)
	if err != nil {
		return nil, fmt.Errorf("TimestampToDateAT(timestamp,template) 不支持%s格式, 正确格式如 YYYY-MM-DD, YYYYMMDDHHmmss ... ", template)
	}
	tm := time.Unix(ms, 0)
	return tm.Format(format), nil
}

// 通用时间格式符号 → Go时间格式模板 映射表
var timeFormatMap = map[string]string{
	// 基础日期格式
	"YYYYMMDD":    "20060102",
	"YYYY-MM-DD":  "2006-01-02",
	"YYYY/MM/DD":  "2006/01/02",
	"YYYY年MM月DD日": "2006年01月02日",
	"MMDD":        "0102",
	"MM-DD":       "01-02",
	"MM/DD":       "01/02",
	"YYYYMM":      "200601",
	"YYYY-MM":     "2006-01",

	// 基础时间格式
	"HHmmss":   "150405",
	"HH:mm:ss": "15:04:05",
	"HHmm":     "1504",
	"HH:mm":    "15:04",
	"mmss":     "0405",
	"mm:ss":    "04:05",

	// 组合格式（最常用）
	"YYYYMMDDHHmmss":      "20060102150405",
	"YYYY-MM-DD HH:mm:ss": "2006-01-02 15:04:05",
	"YYYY/MM/DD HH:mm:ss": "2006/01/02 15:04:05",
	"YYYY-MM-DD HH:mm":    "2006-01-02 15:04",
	"YYYYMMDDHHmm":        "200601021504",

	// 带毫秒格式
	"YYYY-MM-DD HH:mm:ss.SSS": "2006-01-02 15:04:05.000",
	"YYYYMMDDHHmmssSSS":       "20060102150405000",

	// 12小时制格式
	"YYYY-MM-DD hh:mm:ss": "2006-01-02 03:04:05", // 12小时制（需配合AM/PM）
}

// 将通用时间格式符号转换为Go的时间格式模板
func getGoTimeFormat(inputFormat string) (string, error) {
	// 去除首尾空格，统一大小写（兼容yyyyMMdd等写法）
	normalizedFormat := strings.TrimSpace(strings.ToUpper(inputFormat))

	if goFormat, ok := timeFormatMap[normalizedFormat]; ok {
		return goFormat, nil
	}

	// 动态替换, 规则：按字符长度从长到短替换，避免部分匹配
	replaceRules := []struct {
		old string
		new string
	}{
		{"YYYY", "2006"},
		{"MM", "01"},
		{"DD", "02"},
		{"HH", "15"}, // 24小时制
		{"hh", "03"}, // 12小时制
		{"mm", "04"},
		{"ss", "05"},
		{"SSS", "000"},
		{"YY", "06"}, // 2位年份
		{"M", "1"},   // 1位月份
		{"D", "2"},   // 1位日期
		{"H", "15"},  // 1位24小时
		{"h", "3"},   // 1位12小时
		{"m", "4"},   // 1位分钟
		{"s", "5"},   // 1位秒
	}

	goFormat := normalizedFormat
	for _, rule := range replaceRules {
		goFormat = strings.ReplaceAll(goFormat, rule.old, rule.new)
	}

	// 验证替换后的格式是否合法（避免无意义格式）
	if isInvalidFormat(goFormat) {
		return "", fmt.Errorf("不支持的时间格式：%s", inputFormat)
	}

	return goFormat, nil
}

// 简单验证Go时间格式是否合法（避免空/纯符号）
func isInvalidFormat(format string) bool {
	// 匹配仅包含符号的格式（无时间占位符）
	reg := regexp.MustCompile(`^[^\d]+$`)
	return format == "" || reg.MatchString(format)
}

// 通用时间格式化函数（封装常用逻辑）
func formatTime(t time.Time, inputFormat string) (string, error) {
	goFormat, err := getGoTimeFormat(inputFormat)
	if err != nil {
		return "", err
	}
	return t.Format(goFormat), nil
}

// 通用时间解析函数（将字符串转time.Time）
func parseTime(timeStr, inputFormat string) (time.Time, error) {
	goFormat, err := getGoTimeFormat(inputFormat)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(goFormat, timeStr)
}

func timeBeginDayUnix(args []interpreter.Value) (interpreter.Value, error) {
	return gt.BeginDayUnix(), nil
}

func timeEndDayUnix(args []interpreter.Value) (interpreter.Value, error) {
	return gt.EndDayUnix(), nil
}

func timeMinuteAgo(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MinuteAgo(int) 需要一个参数")
	}
	ms, err := getInt(args[0])
	if err != nil {
		return nil, fmt.Errorf("MinuteAgo(int) %s", err.Error())
	}
	return gt.MinuteAgo(ms), nil
}

func timeHourAgo(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("HourAgo(int) 需要一个参数")
	}
	ms, err := getInt(args[0])
	if err != nil {
		return nil, fmt.Errorf("HourAgo(int) %s", err.Error())
	}
	return gt.HourAgo(ms), nil
}

func timeDayAgo(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DayAgo(int) 需要一个参数")
	}
	ms, err := getInt(args[0])
	if err != nil {
		return nil, fmt.Errorf("DayAgo(int) %s", err.Error())
	}
	return gt.DayAgo(ms), nil
}

func timeDayDiffAtUnix(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("DayDiffAtUnix(int, int) 需要两个参数")
	}
	ms1, err := getInt64(args[0])
	if err != nil {
		return nil, fmt.Errorf("DayDiffAtUnix(int, int) %s", err.Error())
	}
	ms2, err := getInt64(args[1])
	if err != nil {
		return nil, fmt.Errorf("DayDiffAtUnix(int, int) %s", err.Error())
	}
	return gt.DayDiffAtUnix(ms1, ms2), nil
}

func timeDayDiff(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("DayDiff(str, str) 需要两个参数")
	}
	s1, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("DayDiff(str, str) 参数要求是字符串 ")
	}
	s2, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("DayDiff(str, str) 参数要求是字符串 ")
	}
	return gt.DayDiff(s1, s2), nil
}

func timeNowToEnd(args []interpreter.Value) (interpreter.Value, error) {
	res, err := gt.NowToEnd()
	if err != nil {
		return nil, fmt.Errorf("NowToEnd出现错误%s", err.Error())
	}
	return res, nil
}

func timeIsToday(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DayAgo(int) 需要一个参数")
	}
	ms, err := getInt64(args[0])
	if err != nil {
		return nil, fmt.Errorf("DayAgo(int) %s", err.Error())
	}
	return gt.IsToday(ms), nil
}

func timeTimestamp2Week(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DayAgo(int) 需要一个参数")
	}
	ms, err := getInt64(args[0])
	if err != nil {
		return nil, fmt.Errorf("DayAgo(int) %s", err.Error())
	}
	return gt.Timestamp2Week(ms), nil
}

func timeTimestamp2WeekXinQi(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DayAgo(int) 需要一个参数")
	}
	ms, err := getInt64(args[0])
	if err != nil {
		return nil, fmt.Errorf("DayAgo(int) %s", err.Error())
	}
	return gt.Timestamp2WeekXinQi(ms), nil
}
