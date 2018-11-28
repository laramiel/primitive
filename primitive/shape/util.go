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
	cos, sin := math.Cos(theta), math.Sin(theta)
	rx = x*cos - y*sin
	ry = x*sin + y*cos
	return
}

// rotateAbout rotates the points x, y about x0, y0.
// cos, sin are math.Cos(theta), math.Sin(theta)
func rotateAbout(x, y int, x0, y0 int, cos, sin float64) (int, int) {
	xd := float64(x - x0)
	yd := float64(y - y0)
	return int(xd*cos - yd*sin + float64(x0)), int(xd*sin + yd*cos + float64(y0))
}
