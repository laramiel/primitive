package primitive

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

type Color struct {
	R, G, B, A int
}

func (c *Color) NRGBA() color.NRGBA {
	return color.NRGBA{uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)}
}

func (c *Color) Delta(color *Color) Color {
	x := Color{c.R - color.R, c.G - color.G, c.B - color.B, c.A - color.A}
	if x.R < 0 {
		x.R = -x.R
	}
	x.R = clampInt(x.R, 0, 255)
	if x.G < 0 {
		x.G = -x.G
	}
	x.G = clampInt(x.G, 0, 255)
	if x.B < 0 {
		x.B = -x.B
	}
	x.B = clampInt(x.B, 0, 255)
	if x.A < 0 {
		x.A = -x.A
	}
	x.A = clampInt(x.A, 0, 255)
	return x
}

func MakeColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	result := Color{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
	vv("%v\n", result)
	return result
}

func MakeHexColor(x string) Color {
	x = strings.Trim(x, "#")
	var r, g, b, a int
	a = 255
	switch len(x) {
	case 3:
		fmt.Sscanf(x, "%1x%1x%1x", &r, &g, &b)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
	case 4:
		fmt.Sscanf(x, "%1x%1x%1x%1x", &r, &g, &b, &a)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
		a = (a << 4) | a
	case 6:
		fmt.Sscanf(x, "%02x%02x%02x", &r, &g, &b)
	case 8:
		fmt.Sscanf(x, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}
	result := Color{r, g, b, a}
	vv("%v\n", result)
	return result
}

// MostFrequentImageColor returns the average color in the image.
func AverageImageColor(im image.Image) color.NRGBA {
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y
	var r, g, b int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
		}
	}
	r /= w * h
	g /= w * h
	b /= w * h
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}

// MostFrequentImageColor returns the most-frequently used color in the image.
// NOTE: The low-order bits are masked off.
func MostFrequentImageColor(im image.Image) color.NRGBA {
	const mask = 0xff - 0x03
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y

	frequency := make(map[color.RGBA]int)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			c.A = 0
			// discard low bits.
			c.R &= mask
			c.G &= mask
			c.B &= mask
			frequency[c]++
		}
	}

	var best color.RGBA
	m := 0
	for k, v := range frequency {
		vv("%v = %d", k, v)
		if v > m {
			best = k
			m = v
		}
	}
	return color.NRGBA{best.R, best.G, best.B, 255}
}

// ColorAtPoint returns the color at a point in the image.
func ColorAtPoint(im image.Image, x, y int) color.NRGBA {
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	if x < 0 || x > size.X {
		x = 0
	}
	if y < 0 || y > size.Y {
		y = 0
	}
	c := rgba.RGBAAt(x, y)
	return color.NRGBA{c.R, c.G, c.B, 255}
}
