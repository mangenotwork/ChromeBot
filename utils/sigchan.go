package utils

import "os"

var SigChan = make(chan os.Signal, 1)
var RunMode = "REPL"
var ScriptDir = ""
