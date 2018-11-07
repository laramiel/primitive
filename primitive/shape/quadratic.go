package shape

import (
	"fmt"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Quadratic struct {
	X1, Y1       float64
	X2, Y2       float64
	X3, Y3       float64
	Width        float64
	MaxLineWidth float64
}

func NewQuadratic() *Quadratic {
	q := &Quadratic{}
	q.MaxLineWidth = 1.0 / 2
	return q
}

func (q *Quadratic) Init(plane *Plane) {
	rnd := plane.Rnd
	q.X1 = rnd.Float64() * float64(plane.W)
	q.Y1 = rnd.Float64() * float64(plane.H)
	q.X2 = q.X1 + rnd.Float64()*40 - 20
	q.Y2 = q.Y1 + rnd.Float64()*40 - 20
	q.X3 = q.X2 + rnd.Float64()*40 - 20
	q.Y3 = q.Y2 + rnd.Float64()*40 - 20
	q.Width = 1.0 / 2
	q.Mutate(plane)
}

func (q *Quadratic) Draw(dc *gg.Context, scale float64) {
	dc.MoveTo(q.X1, q.Y1)
	dc.QuadraticTo(q.X2, q.Y2, q.X3, q.Y3)
	dc.SetLineWidth(q.Width * scale)
	dc.Stroke()
}

func (q *Quadratic) SVG(attrs string) string {
	// TODO: this is a little silly
	attrs = strings.Replace(attrs, "fill", "stroke", -1)
	return fmt.Sprintf(
		"<path %s fill=\"none\" d=\"M %f %f Q %f %f, %f %f\" stroke-width=\"%f\" />",
		attrs, q.X1, q.Y1, q.X2, q.Y2, q.X3, q.Y3, q.Width)
}

func (q *Quadratic) Copy() Shape {
	a := *q
	return &a
}

func (q *Quadratic) Mutate(plane *Plane) {
	const m = 16
	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	for {
		switch rnd.Intn(4) {
		case 0:
			q.X1 = clamp(q.X1+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y1 = clamp(q.Y1+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 1:
			q.X2 = clamp(q.X2+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y2 = clamp(q.Y2+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 2:
			q.X3 = clamp(q.X3+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y3 = clamp(q.Y3+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 3:
			q.Width = clamp(q.Width+rnd.NormFloat64(), 0.25, q.MaxLineWidth)
		}
		if q.Valid() {
			break
		}
	}
}

func (q *Quadratic) Valid() bool {
	dx12 := int(q.X1 - q.X2)
	dy12 := int(q.Y1 - q.Y2)
	dx23 := int(q.X2 - q.X3)
	dy23 := int(q.Y2 - q.Y3)
	dx13 := int(q.X1 - q.X3)
	dy13 := int(q.Y1 - q.Y3)
	d12 := dx12*dx12 + dy12*dy12
	d23 := dx23*dx23 + dy23*dy23
	d13 := dx13*dx13 + dy13*dy13
	return d13 > d12 && d13 > d23
}

func (q *Quadratic) Rasterize(rc *RasterContext) []Scanline {
	var path raster.Path
	p1 := fixp(q.X1, q.Y1)
	p2 := fixp(q.X2, q.Y2)
	p3 := fixp(q.X3, q.Y3)
	path.Start(p1)
	path.Add2(p2, p3)
	width := fix(q.Width)
	return strokePath(rc, path, width, raster.RoundCapper, raster.RoundJoiner)
}
