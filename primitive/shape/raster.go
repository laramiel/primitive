package shape

import (
	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(x * 64)
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{fix(x), fix(y)}
}

type painter struct {
	Lines []Scanline
}

func (p *painter) Paint(spans []raster.Span, done bool) {
	for _, span := range spans {
		p.Lines = append(p.Lines, Scanline{span.Y, span.X0, span.X1 - 1, span.Alpha})
	}
}

func fillPath(rc *RasterContext, path raster.Path) []Scanline {
	r := rc.Rasterizer
	r.Clear()
	r.UseNonZeroWinding = true
	r.AddPath(path)
	var p painter
	p.Lines = rc.Lines[:0]
	r.Rasterize(&p)
	return p.Lines
}

func strokePath(rc *RasterContext, path raster.Path, width fixed.Int26_6, cr raster.Capper, jr raster.Joiner) []Scanline {
	r := rc.Rasterizer
	r.Clear()
	r.UseNonZeroWinding = true
	r.AddStroke(path, width, cr, jr)
	var p painter
	p.Lines = rc.Lines[:0]
	r.Rasterize(&p)
	return p.Lines
}
