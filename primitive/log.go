package primitive

import "github.com/laramiel/primitive/primitive/log"

func v(format string, a ...interface{}) {
	log.Log(1, format, a...)
}

func vv(format string, a ...interface{}) {
	log.Log(2, "  "+format, a...)
}

func vvv(format string, a ...interface{}) {
	log.Log(3, "    "+format, a...)
}
