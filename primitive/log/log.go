package log

import "fmt"

var LogLevel int

func Log(level int, format string, a ...interface{}) {
	if LogLevel >= level {
		fmt.Printf(format, a...)
	}
}
