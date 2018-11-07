package primitive

import (
	"image"
  "strings"
  
	"github.com/laramiel/primitive/primitive/shape"
)

type ColorPicker interface {
	// Select provides a mechanism to pick the best color from a set of
	// scanlines and a target image.
	Select(target, current *image.RGBA, lines []shape.Scanline, alpha int) Color
}

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
	return Color{r, g, b, alpha}
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
	return Color{bw, bw, bw, alpha}
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
			rsum += int64((tr-cr))
			ssum += int64(cr)
			gsum += int64((tg-cg))
			hsum += int64(cg)
			bsum += int64((tb-cb))
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
	return Color{r, g, b, alpha}
}

type selectedColors struct {
	colors []Color	
	b BestColor
}

func (s *selectedColors) Select(target, current *image.RGBA, lines []shape.Scanline, alpha int) Color {
	if len(s.colors) == 1 {
		return s.colors[0]
	}
	best := s.b.Select(target, current, lines, alpha)

	selected := best	
	var delta int = 256 * 256 * 256;
	for _, c := range s.colors {
		d := c.Delta(&best)
		x := d.R + d.B + d.G 
		if x < delta {
			selected = c
			delta = x
		}
	}
	selected.A = alpha
	return selected
}

func MakeColorPicker(config string) ColorPicker {
	if config == "greyscale" {
		return &BestGreyscale{}
	}
	if config == "alpha" {
		return &BestAlpha{}		
	}

	var colors []Color
	for _, v := range strings.Split(config, "," ) {
		c := MakeHexColor(v)
		colors = append(colors, c)
	}
	if len(colors) == 0 {
		return &BestColor{}
	}
	return &selectedColors{colors, BestColor{}}	
}
