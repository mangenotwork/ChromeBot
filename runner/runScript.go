package runner

import (
	"ChromeBot/dsl/builtins"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/dsl/lexer"
	"ChromeBot/dsl/parser"
	"ChromeBot/global"
	"ChromeBot/utils"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func runScript(source string) {

	source = utils.RemoveNewlinesInBackticks(source)
	source = utils.ProcessCommandLine(source)
	source = utils.EscapeQuotesInBackticks(source)
	source = globalAnalysisScript(source)

	builtins.ChromeWait = 2

	// 词法分析
	l := lexer.New(source)

	// 语法分析
	p := parser.New(l)
	program := p.ParseProgram()

	errs := p.CleanErrors()
	if len(errs) > 0 {
		fmt.Println("解析错误:")
		for _, err := range errs {
			fmt.Println("  " + err)
		}
		return
	}

	// 创建解释器
	interp := interpreter.NewInterpreter()

	// 注册内置函数
	builtins.RegisterBuiltins(interp)

	// 执行程序

	if global.IsRegisterCron {
		fmt.Println("开启了定时任务 ", global.CronPerformance.Arg, " -> ", global.CronToChinese(global.CronPerformance.Arg))

		c := cron.New(
			cron.WithSeconds(),
			cron.WithParser(cron.NewParser(
				cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow,
			)),
		)

		taskID, err := c.AddFunc(global.CronPerformance.Arg, func() {

			entry := c.Entries()[0]
			currentTime := time.Now().Format(time.DateTime)
			nextTime := entry.Next.In(time.Local).Format(time.DateTime)
			fmt.Printf("[Cron] 定时任务执行脚本; 当前时间：%s | 任务下次执行时间：%s \n", currentTime, nextTime)

			result, err := interp.Interpret(program)
			if err != nil {
				fmt.Printf("执行错误: %v\n", err)
				return
			}

			if result != nil {
				fmt.Printf("程序返回值: %v\n", result)
			}
		})
		if err != nil {
			fmt.Printf("添加任务失败：%v\n", err)
			return
		}
		fmt.Printf("添加任务成功，ID：%d\n", taskID)

		// 启动cron
		c.Start()

		// 阻塞主线程
		select {}

	} else {
		result, err := interp.Interpret(program)
		if err != nil {
			fmt.Printf("执行错误: %v\n", err)
			return
		}

		if result != nil {
			fmt.Printf("程序返回值: %v\n", result)
		}

	}

}
