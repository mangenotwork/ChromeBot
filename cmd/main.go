package main

import (
	"ChromeBot/runner"
	"flag"
	"fmt"
	"os"
)

// 定义全局命令行参数（也可以在 Run 函数内定义，按需选择）
var (
	// -v 查看版本，布尔类型，默认 false
	showVersion = flag.Bool("v", false, "查看 ChromeBot 版本信息")
	// -h 查看帮助，布尔类型，默认 false
	showHelp = flag.Bool("h", false, "查看 ChromeBot 帮助信息")
)

func main() {

	// ========== 第一步：解析命令行参数（必须在 os.Args 操作前执行） ==========
	// flag.Parse() 会自动解析 os.Args[1:] 中的参数，绑定到上面定义的变量
	flag.Parse()

	// ========== 第二步：处理 -v 和 -h 参数逻辑（优先执行，执行后退出） ==========
	// 处理 -h 查看帮助
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// 处理 -v 查看版本
	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	runner.Run()
}

// ========== 新增：打印版本信息 ==========
func printVersion() {
	fmt.Printf("v%s\n", runner.VERSION)
}

// ========== 新增：打印帮助信息 ==========
func printHelp() {
	fmt.Println("")
	fmt.Printf("ChromeBot v%s ( https://github.com/mangenotwork/ChromeBot )\n", runner.VERSION)
	fmt.Println("")
	fmt.Println("简介：")
	fmt.Println("  谷歌浏览器（Chrome）自动化平台，通过输入指令或脚本(chrome bot script)自动执行操作Chrome")
	fmt.Println("")
	fmt.Println("用法：")
	fmt.Println("  chromebot [选项] [文件名]")
	fmt.Println("")
	fmt.Println("选项：")
	// flag.PrintDefaults() 会自动打印所有定义的 flag 说明（无需手动写）
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("示例：")
	fmt.Println("  chromebot          # 启动交互式 REPL 环境")
	fmt.Println("  chromebot test.cbs # 执行 test.cbs 中的代码")
	fmt.Println("  chromebot -v       # 查看版本信息")
	fmt.Println("  chromebot -h       # 查看帮助信息")
}
