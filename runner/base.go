package runner

import (
	"ChromeBot/utils"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	VERSION    = "0.0.3"
	PROMPT     = ">>> "
	PROMPTCont = "... "
)

func Run() {

	utils.IsDebug = true
	gt.CloseLog()

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

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		fmt.Printf("程序已启动 主PID: %d\n", os.Getpid())

		//ctx, cancel := context.WithCancel(context.Background())
		//defer cancel()

		go runREPL(sigChan)

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

	} else { // 检查文件
		filename := os.Args[1]
		fileExt := filepath.Ext(filename)
		if fileExt != ".cbs" {
			fmt.Printf("无法读取文件 %s, 文件应为后缀是.cbs的脚本文件\n", filename)
			return
		}
		source, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("无法读取文件 %s: %v\n", filename, err)
			return
		}
		utils.Debug(source)

		runScript(string(source))
	}

}
