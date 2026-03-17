package global

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// 定义语法规则  @cron 0 * * * * *
// 常用的cron
// */10 * * * * *	: 每 10 秒执行一次,实时监控、心跳检测
// 0 * * * * * 		: 每分钟执行一次（整点秒）,分钟级数据统计、日志清理
// 0 */5 * * * *	: 每 5 分钟执行一次,常规定时任务（如数据同步）
// 0 0 * * * *		: 每小时执行一次（整点）,小时级汇总、接口巡检
// 0 0 0 * * * 		: 每天凌晨 0 点执行,每日数据备份、日志归档
// 0 30 8 * * *		: 每天早上 8:30 执行, 每日早间通知、报表生成
// 0 0 22 * * * 	: 每天晚上 22:00 执行, 晚间定时任务、系统维护
// 0 0 9,18 * * * 	: 每天 9:00 和 18:00 各执行一次, 上下班时段的业务触发
// 0 0 9 * * 1		: 每周一早上 9:00 执行,每周一的工作任务、周报生成
// 0 0 18 * * 1-5	: 每周一至周五 18:00 执行, 工作日下班前的数据处理
// 0 0 10 * * 6,0	: 每周六、周日 10:00 执行, 周末定时任务（0 = 周日）
// 0 0 0 1 * * 		: 每月 1 号凌晨 0 点执行, 月度账单生成、数据对账
// 0 30 14 15 * *	: 每月 15 号 14:30 执行, 月度中期统计、提醒

var CronPerformance cronExecute
var cronOnce sync.Once
var IsRegisterCron = false

type cronExecute struct {
	Arg string
}

func NewCronPerformance(arg string) {
	cronOnce.Do(func() {
		CronPerformance = cronExecute{
			Arg: arg,
		}
	})
}

func RegisterCron(arg string) {
	fmt.Println("注册定时任务")
	IsRegisterCron = true
	NewCronPerformance(arg)
}

// CronField 定义cron字段的解析规则
type CronField struct {
	ChineseName string // 中文名称（秒/分/时/日/月/周）
	Min         int    // 最小值
	Max         int    // 最大值
}

// 初始化cron字段规则（6字段：秒、分、时、日、月、周）
var cronFields = []CronField{
	{"秒", 0, 59},
	{"分", 0, 59},
	{"时", 0, 23},
	{"日", 1, 31},
	{"月", 1, 12},
	{"周", 0, 6},
}

// 星期中文映射
var weekChineseMap = map[int]string{
	0: "周日",
	1: "周一",
	2: "周二",
	3: "周三",
	4: "周四",
	5: "周五",
	6: "周六",
}

// 解析步长表达式（如 */10 → 10，5/2 → 2）
func parseStep(fieldStr string) (int, bool) {
	stepRegex := regexp.MustCompile(`^(\*|\d+)-?\d*\/(\d+)$`)
	if !stepRegex.MatchString(fieldStr) {
		return 0, false
	}
	matches := stepRegex.FindStringSubmatch(fieldStr)
	step, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, false
	}
	return step, true
}

// 解析单个值/范围/枚举
func parseSingleValue(fieldStr string, field CronField) (string, error) {
	// 处理?（无指定值）
	if fieldStr == "?" {
		return "", nil
	}

	// 处理范围（n-m）
	rangeRegex := regexp.MustCompile(`^(\d+)-(\d+)$`)
	if rangeRegex.MatchString(fieldStr) {
		matches := rangeRegex.FindStringSubmatch(fieldStr)
		start, _ := strconv.Atoi(matches[1])
		end, _ := strconv.Atoi(matches[2])
		// 星期特殊处理
		if field.ChineseName == "周" {
			return fmt.Sprintf("%s至%s", weekChineseMap[start], weekChineseMap[end]), nil
		}
		return fmt.Sprintf("%d至%d%s", start, end, field.ChineseName), nil
	}

	// 处理枚举（n,m,k）
	if strings.Contains(fieldStr, ",") {
		parts := strings.Split(fieldStr, ",")
		var descParts []string
		for _, part := range parts {
			val, _ := strconv.Atoi(part)
			if field.ChineseName == "周" {
				descParts = append(descParts, weekChineseMap[val])
			} else {
				descParts = append(descParts, fmt.Sprintf("%d%s", val, field.ChineseName))
			}
		}
		return strings.Join(descParts, "、"), nil
	}

	// 处理单个值
	val, err := strconv.Atoi(fieldStr)
	if err != nil {
		return "", err
	}
	if field.ChineseName == "周" {
		return weekChineseMap[val], nil
	}
	return fmt.Sprintf("%d%s", val, field.ChineseName), nil
}

// CronToChinese 核心：将cron表达式转为简洁中文描述
func CronToChinese(cronExpr string) string {
	// 预处理：去除空格，分割字段
	fields := strings.Fields(strings.TrimSpace(cronExpr))

	// 格式校验：5位（标准cron）→ 补全6位（秒为*）；非5/6位直接报错
	if len(fields) == 5 {
		fields = append([]string{"*"}, fields...)
	} else if len(fields) != 6 {
		fmt.Printf("无效格式！仅支持5位（分 时 日 月 周）或6位（秒 分 时 日 月 周），当前：%d位", len(fields))
		return ""
	}

	// ========== 核心优化：优先处理「单字段步长，其余全*」的高频场景 ==========
	stepCount := 0
	var stepFieldIndex int
	var stepValue int
	for i, f := range fields {
		if f == "*" {
			continue
		}
		// 尝试解析步长
		step, ok := parseStep(f)
		if ok {
			stepCount++
			stepFieldIndex = i
			stepValue = step
		} else {
			// 非*且非步长 → 走通用逻辑
			goto generalLogic
		}
	}

	// 仅单个字段是步长，其余全* → 生成简洁描述
	if stepCount == 1 {
		fieldName := cronFields[stepFieldIndex].ChineseName
		if stepValue == 1 {
			return fmt.Sprintf("每%s执行一次", fieldName)
		}
		return fmt.Sprintf("每%d%s执行一次", stepValue, fieldName)
	}

	// 所有字段都是* → 每秒执行一次
	if stepCount == 0 {
		return "每秒执行一次"
	}

	// ========== 通用场景：非单字段步长 ==========
generalLogic:
	var descParts []string
	for i, f := range fields {
		field := cronFields[i]
		if f == "*" || f == "?" {
			continue
		}

		// 解析步长
		if step, ok := parseStep(f); ok {
			if step == 1 {
				descParts = append(descParts, fmt.Sprintf("每%s", field.ChineseName))
			} else {
				descParts = append(descParts, fmt.Sprintf("每%d%s", step, field.ChineseName))
			}
			continue
		}

		// 解析单个值/范围/枚举
		desc, err := parseSingleValue(f, field)
		if err != nil {
			fmt.Printf("解析%s字段「%s」失败：%v", field.ChineseName, f, err)
			return ""
		}
		if desc != "" {
			descParts = append(descParts, desc)
		}
	}

	// 拼接通用场景描述
	if len(descParts) == 0 {
		return "每秒执行一次"
	}
	return fmt.Sprintf("在%s执行一次", strings.Join(descParts, "、"))
}
