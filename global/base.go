package global

import (
	"ChromeBot/browser"
	"fmt"
	"os"
)

type Command string

var (
	Cron         Command = "cron"
	ConfJson     Command = "conf_json"
	ConfYaml     Command = "conf_yaml"
	ConfINI      Command = "conf_ini"
	ChromeCheck  Command = "chrome_check"
	NetworkCheck Command = "network_check"
)

var globalSupport = map[Command]bool{
	Cron:         true,
	ConfJson:     true,
	ConfYaml:     true,
	ConfINI:      true,
	ChromeCheck:  true,
	NetworkCheck: true,
}

func HasGlobalSupport(cmd Command) bool {
	_, ok := globalSupport[cmd]
	return ok
}

func ChromeCheckImplement() {
	chromePath, err := browser.FindChrome()
	if err != nil {
		fmt.Printf("[Err]本机未找到Chrome浏览器，请前往安装 https://www.google.cn/chrome/ ")
		os.Exit(0)
	}
	fmt.Printf("[@chrome_check]本机已找到Chrome浏览器, 路径为：%s \n", chromePath)
}

func NetworkCheckImplement(arg string) {
	fmt.Printf("[@network_check]正在检查网络连接... %s \n", arg)
	// err := browser.GetFirstTabWs()
	// if err != nil {
	// 	fmt.Printf("[Err]网络连接异常，请检查网络连接 \n")
	// 	os.Exit(0)
	// }
	fmt.Printf("[@network_check]网络连接正常 \n")
}
