package shape

import (
	"math"

	"github.com/laramiel/primitive/primitive/log"
)

func v(format string, a ...interface{}) {
	log.Log(1, format, a...)
}

func vv(format string, a ...interface{}) {
	log.Log(2, "  "+format, a...)
}

func vvv(format string, a ...interface{}) {
	log.Log(3, "    "+format, a...)
}

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func clampInt(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rotate(x, y, theta float64) (rx, ry float64) {
	rx = x*math.Cos(theta) - y*math.Sin(theta)
	ry = x*math.Sin(theta) + y*math.Cos(theta)
	return
}
