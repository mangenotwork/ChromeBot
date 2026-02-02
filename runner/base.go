package runner

import (
	"ChromeBot/dsl/builtins"
	"ChromeBot/dsl/interpreter"
	"ChromeBot/dsl/lexer"
	"ChromeBot/dsl/parser"
	"ChromeBot/utils"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	VERSION    = "0.0.1"
	PROMPT     = ">>> "
	PROMPTCont = "... "
)

func Run() {

	//utils.IsDebug = true

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("程序已启动 主PID: %d\n", os.Getpid())

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	if len(os.Args) < 2 {
		fmt.Printf("_________\n")
		fmt.Printf("|       |\n")
		fmt.Printf("|  o o  |\n")
		fmt.Printf("|   c   |\\\n")
		fmt.Printf("|_______| \\_chrome\n\n")
		fmt.Printf("欢迎使用 ChromeBot v%s\n", VERSION)
		fmt.Println("https://github.com/mangenotwork/ChromeBot")
		fmt.Println("输入代码并按回车执行。")
		fmt.Println("按Ctrl+C或者Ctrl+Z退出程序，也可以使用 'exit' 或 'quit' 命令退出程序。")
		fmt.Println("===================================================================")

		go runREPL(sigChan)

	} else { // 检查文件
		filename := os.Args[1]
		source, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("无法读取文件 %s: %v\n", filename, err)
			return
		}
		utils.Debug(source)
		//runScript(string(source))
	}

	sig := <-sigChan
	fmt.Printf("\n接收到信号: %v。开始执行清理工作...\n", sig)

	// 区分信号类型（可选）
	switch sig {
	case syscall.SIGINT:
		fmt.Println("这是 Ctrl+C 触发的中断信号")
	case syscall.SIGTERM:
		fmt.Println("这是 kill 命令触发的终止信号")
	}

	// 此处可以添加你的清理逻辑，例如关闭文件、断开网络连接等
	// ...
	//defer close(done)
	//cancel()
	//for _, v := range pidList {
	//	gt.Info("正在关闭进程 pid ====> ", v)
	//	time.Sleep(1 * time.Second)
	//	_ = SafeKillProcess(v)
	//}
	os.Exit(0)

}

func runREPL(sigChan chan os.Signal) {
	scanner := bufio.NewScanner(os.Stdin)
	// 创建解释器
	interp := interpreter.NewInterpreter()
	// 注册内置函数
	builtins.RegisterBuiltins(interp)

	fmt.Print(PROMPT)

	var inputLines []string
	braceCount := 0
	parenCount := 0
	bracketCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if shouldExit(line) {
			fmt.Println("BayBay.")
			sigChan <- syscall.SIGTERM
			return
		}

		// 统计括号数量以确定是否继续输入
		braceCount += countChars(line, '{', '}')
		parenCount += countChars(line, '(', ')')
		bracketCount += countChars(line, '[', ']')

		inputLines = append(inputLines, line)

		// 如果所有括号都匹配，执行代码
		if braceCount == 0 && parenCount == 0 && bracketCount == 0 {
			// 拼接所有行
			fullInput := strings.Join(inputLines, "\n")

			// 执行代码
			executeCode(fullInput, interp)

			// 重置状态
			inputLines = nil
			braceCount = 0
			parenCount = 0
			bracketCount = 0

			fmt.Print(PROMPT)
		} else {

			fmt.Print(PROMPTCont)
		}
	}

	if len(inputLines) > 0 {
		fullInput := strings.Join(inputLines, "\n")
		fmt.Print(fullInput)
		executeCode(fullInput, interp)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取输入错误: %v\n", err)
	}

	fmt.Println("BayBay.")
	sigChan <- syscall.SIGTERM
}

// 检查是否应该退出
func shouldExit(line string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(line))
	return trimmed == "exit" || trimmed == "quit" || trimmed == ":q"
}

// 统计括号数量
func countChars(line string, openChar, closeChar byte) int {
	count := 0
	for i := 0; i < len(line); i++ {
		if line[i] == openChar {
			count++
		} else if line[i] == closeChar {
			count--
		}
	}
	return count
}

// 执行代码并输出结果
func executeCode(input string, interp *interpreter.Interpreter) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// 处理特殊命令
	if handleSpecialCommands(input) {
		return
	}

	// 记录执行前的状态
	// 这里可以记录一些状态，如果需要的话

	// 解析和执行
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		printErrors("解析错误", p.Errors())
		return
	}

	result, err := interp.Interpret(program)
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		return
	}

	// 关键：只输出包含return或print的结果
	// 或者只输出表达式的值
	shouldOutput := false

	// 检查是否应该输出
	if result != nil {
		trimmedInput := strings.TrimSpace(input)

		// 如果输入以return开头，应该输出
		if strings.HasPrefix(trimmedInput, "return ") {
			shouldOutput = true
		} else if isExpression(trimmedInput) { // 如果输入是表达式（不以语句关键字开头），应该输出
			shouldOutput = true
		} else if strings.Contains(trimmedInput, "print(") { // 如果输入包含print，已经在print函数中输出了
			shouldOutput = false
		}
	}

	if shouldOutput {
		printResult(result)
	}
}

// 判断是否是表达式
func isExpression(input string) bool {
	trimmed := strings.TrimSpace(input)

	// 空行不是表达式
	if trimmed == "" {
		return false
	}

	// 检查是否是语句
	statements := []string{
		"var ", "if ", "while ", "for ",
		"func ", "break ", "continue ",
		"print(",
	}

	for _, stmt := range statements {
		if strings.HasPrefix(trimmed, stmt) {
			return false
		}
	}

	// 检查是否包含赋值（但不是比较）
	if strings.Contains(trimmed, "=") && !strings.Contains(trimmed, "==") {
		// 包含单个等号，可能是赋值语句
		return false
	}

	// 其他情况可能是表达式
	return true
}

// 格式化并打印结果
func printResult(result interpreter.Value) {
	switch v := result.(type) {
	case int64:
		fmt.Println(v)
	case float64:
		fmt.Println(v)
	case string:
		fmt.Printf("%q\n", v)
	case bool:
		fmt.Println(v)
	case []interpreter.Value:
		// 列表
		fmt.Print("[")
		for i, item := range v {
			if i > 0 {
				fmt.Print(", ")
			}
			printValue(item)
		}
		fmt.Println("]")
	case interpreter.DictType:
		// 字典
		fmt.Print("{")
		first := true
		for key, value := range v {
			if !first {
				fmt.Print(", ")
			}
			first = false
			fmt.Printf("%v: ", key)
			printValue(value)
		}
		fmt.Println("}")
	case nil:
		// 不输出nil
	default:
		fmt.Printf("%v\n", v)
	}
}

// 打印单个值
func printValue(value interpreter.Value) {
	switch v := value.(type) {
	case string:
		fmt.Printf("%q", v)
	default:
		fmt.Printf("%v", v)
	}
}

// 打印错误信息
func printErrors(prefix string, errors []string) {
	fmt.Printf("%s:\n", prefix)
	for _, err := range errors {
		fmt.Printf("  %s\n", err)
	}
}

// 处理特殊命令
func handleSpecialCommands(input string) bool {
	trimmed := strings.TrimSpace(input)

	switch strings.ToLower(trimmed) {
	case "clear", "cls":
		clearScreen()
		return true
	case "env", "variables":
		// todo 这里可以添加查看环境变量的功能
		fmt.Println("环境变量功能待实现")
		return true
	}

	return false
}

func clearScreen() {
	// 简单的清屏：打印多个空行
	for i := 0; i < 50; i++ {
		fmt.Println()
	}
}
