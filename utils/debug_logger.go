package utils

import "log"

var IsDebug = false

func Debug(v ...interface{}) {
	if IsDebug {
		value := make([]interface{}, 0)
		value = append(value, "[DEBUG]")
		value = append(value, v...)
		log.Println(value...)
	}

}

func Debugf(format string, v ...interface{}) {
	if IsDebug {
		log.Printf("[DEBUG]"+format, v...)
	}
}
