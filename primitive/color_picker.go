package primitive

import (
	"image"
	"strings"
	"math"

	"github.com/laramiel/primitive/primitive/shape"
)

type ColorPicker interface {
	// Select provides a mechanism to pick the best color from a set of
	// scanlines and a target image.
	Select(target, current *image.RGBA, lines []shape.Scanline, alpha int) Color
}

const Palette1 = "#8b3336,#c83940,#e3acbb,#01afa4,#1ab7ad,#8cd1cc,#d1c0db,#dbc6de,#e6dbe6,#cae1a3,#cfe4a8,#e2edcb,#a59256,#94814c,#c8c7c1,#a5472e,#ac4d33,#e7ac8c,#543b34,#573c33,#aba190,#f0583d,#f26f4e,#f9bdb7,#b33835,#cd3434,#e47294,#f0553b,#f15b39,#f8a78a,#ee422e,#ef4a38,#f58d81,#faa720,#fbc017,#fdde39,#f5eb57,#f3eb5c,#f1efab,#fdbf30,#fdc01d,#fcdf76,#0875b8,#0375bb,#19c0ed,#325eab,#3363ae,#55b4e5,#e99c62,#ee9a5f,#eacebf,#603b4a,#793854,#d1a0c8,#383230,#383439,#8582bc,#e7c058,#ecd37d,#f1ebde,#135341,#076a42,#0eb69d,#039e4e,#0ea54f,#84c991,#35483b,#2e4636,#96c291,#343c3b,#494c46,#929698,#3dabc7,#43acc8,#99c9d9,#adcbea,#b6d0e9,#d2e0ef,#f8c8a3,#f8cea8,#f8e1c9,#363636,#323231,#666f74,#f284ae,#f288b1,#f7bcd4,#fcef9e,#f7ed9d,#f6f2ca,#b4babf,#b6bbc0,#dadddc,#393a3c,#373e4b,#a6b5c0,#31395c,#1869b0,#13b3e9,#0f5c4c,#0f5a48,#0db7a8,#1e4279,#2569b0,#26bcea,#d54344,#e34446,#f4aac3,#fdd220,#fcda22,#fbea8d,#805e9e,#b688bb,#ceb9d7,#a8335b,#d14b7e,#ea97c1,#a36f43,#b27647,#e1c5a3,#463d30,#5e4835,#b3ab9f,#b94a30,#c24b2f,#e99b90,#343836,#334536,#8ec0a1,#c1c9ca,#bfc4c3,#dce5e7,#e3e5e3,#e1e3e2,#e0e3e2,#2f3971,#34549e,#98bce3,#e1d5af,#e4d3ab,#eae2cd,#e2a530,#e9a636,#f0d6a8"

// BestColor calculates the color that should be used with the Scanlines to
// best approximate the target image.
type BestColor struct {
}

func (s *BestColor) Select(target, current *image.RGBA, lines []shape.Scanline, alpha int) Color {
	var rsum, gsum, bsum, count int64
	a := 0x101 * 255 / alpha
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			cr := int(current.Pix[i])
			cg := int(current.Pix[i+1])
			cb := int(current.Pix[i+2])
			i += 4
			rsum += int64((tr-cr)*a + cr*0x101)
			gsum += int64((tg-cg)*a + cg*0x101)
			bsum += int64((tb-cb)*a + cb*0x101)
			count++
		}
	}
	if count == 0 {
		return Color{}
	}
	r := clampInt(int(rsum/count)>>8, 0, 255)
	g := clampInt(int(gsum/count)>>8, 0, 255)
	b := clampInt(int(bsum/count)>>8, 0, 255)
	result := Color{r, g, b, alpha}
	return result
}

type BestGreyscale struct {
}

func (s *BestGreyscale) Select(target, current *image.RGBA, lines []shape.Scanline, alpha int) Color {
	var sum, count int64
	a := 0x101 * 255 / alpha
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			cr := int(current.Pix[i])
			cg := int(current.Pix[i+1])
			cb := int(current.Pix[i+2])
			i += 4
			sum += int64((tr-cr)*a + cr*0x101)
			sum += int64((tg-cg)*a + cg*0x101)
			sum += int64((tb-cb)*a + cb*0x101)
			count += 3
		}
	}
	if count == 0 {
		return Color{}
	}
	bw := clampInt(int(sum/count)>>8, 0, 255)
	result := Color{bw, bw, bw, alpha}
	return result
}

type BestAlpha struct {
}

func (s *BestAlpha) Select(target, current *image.RGBA, lines []shape.Scanline, alpha int) Color {
	var count int64
	var ssum, hsum, csum int64
	var rsum, gsum, bsum int64
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			cr := int(current.Pix[i])
			cg := int(current.Pix[i+1])
			cb := int(current.Pix[i+2])
			i += 4
			rsum += int64((tr - cr))
			ssum += int64(cr)
			gsum += int64((tg - cg))
			hsum += int64(cg)
			bsum += int64((tb - cb))
			csum += int64(cb)
			count++
		}
	}
	if count == 0 {
		return Color{}
	}
	r := clampInt(int(rsum/count)>>8, 0, 255)
	g := clampInt(int(gsum/count)>>8, 0, 255)
	b := clampInt(int(bsum/count)>>8, 0, 255)
	result := Color{r, g, b, alpha}
	return result
}

// ColorPalette allows you to restrict to a certain range of colors
type ColorPalette struct {
	hexStrings []string
	rgbColors  []Color
	b          BestColor
}

// closestColor returns the index of the closest color.
func (cp *ColorPalette) closestColorIdx(c Color) int {
	selected := 0
	score := 100000.0
	for i, cmp := range cp.rgbColors {
		d := c.Delta(&cmp)
		// minimize the euclidian distance
		x := math.Sqrt(float64(d.R*d.R + d.B*d.B + d.G*d.G))
		if x < score {
			selected = i
			score = x
		}
	}
	return selected
}

// ClosestColor returns the closest RGB color in the current palette
func (cp *ColorPalette) Select(target, current *image.RGBA, lines []shape.Scanline, alpha int) Color {
	if len(cp.rgbColors) == 1 {
		return cp.rgbColors[0]
	}
	best := cp.b.Select(target, current, lines, alpha)
	i := cp.closestColorIdx(best)
	return cp.rgbColors[i]
}

// NewColorPalette returns a new color palette
func NewColorPalette(hexes []string) (cp *ColorPalette) {
	cp = new(ColorPalette)
	cp.hexStrings = hexes
	cp.rgbColors = make([]Color, len(hexes))
	for i, hex := range hexes {
		cp.rgbColors[i] = MakeHexColor(hex)
	}
	return
}

func MakeColorPicker(config string) ColorPicker {
	if config == "" {
		return &BestColor{}
	}
	if config == "greyscale" {
		return &BestGreyscale{}
	}
	if config == "alpha" {
		return &BestAlpha{}
	}
	if config == "palette1" {
		config = Palette1
	}

	cp := NewColorPalette(strings.Split(config, ","))
	if len(cp.rgbColors) == 0 {
		return &BestColor{}
	}
	return cp
}

