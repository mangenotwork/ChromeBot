package global

type Command string

var (
	Cron        Command = "cron"
	ConfJson    Command = "conf_json"
	ConfYaml    Command = "conf_yaml"
	ConfINI     Command = "conf_ini"
	ChromeCheck Command = "chrome_check"
)

var globalSupport = map[Command]bool{
	Cron:        true,
	ConfJson:    true,
	ConfYaml:    true,
	ConfINI:     true,
	ChromeCheck: true,
}

func HasGlobalSupport(cmd Command) bool {
	_, ok := globalSupport[cmd]
	return ok
}
