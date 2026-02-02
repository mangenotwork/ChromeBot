package runner

import (
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

	utils.IsDebug = true

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
	//interp := interpreter.NewInterpreter()

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
			//fullInput := strings.Join(inputLines, "\n")

			// 执行代码
			//executeCode(fullInput, interp)

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
		//executeCode(fullInput, interp)
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
