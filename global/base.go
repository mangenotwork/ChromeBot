package global

type Command string

var (
	Cron        Command = "cron"
	ChromeCheck Command = "chrome_check"
)

var globalSupport = map[Command]bool{
	Cron:        true,
	ChromeCheck: true,
}

func HasGlobalSupport(cmd Command) bool {
	_, ok := globalSupport[cmd]
	return ok
}
