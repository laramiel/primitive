package shape

import (
	"fmt"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Polygon struct {
	Order  int
	Convex bool
	X, Y   []float64
}

func NewPolygon(order int, convex bool) *Polygon {
	p := &Polygon{}
	p.Order = order
	p.Convex = convex
	return p
}

func (p *Polygon) Init(plane *Plane) {
	rnd := plane.Rnd
	p.X = make([]float64, p.Order)
	p.Y = make([]float64, p.Order)
	p.X[0] = rnd.Float64() * float64(plane.W)
	p.Y[0] = rnd.Float64() * float64(plane.H)
	for i := 1; i < p.Order; i++ {
		p.X[i] = p.X[0] + rnd.Float64()*40 - 20
		p.Y[i] = p.Y[0] + rnd.Float64()*40 - 20
	}
	p.mutateImpl(plane, 1.0, 1)
}

func (p *Polygon) Draw(dc *gg.Context, scale float64) {
	dc.NewSubPath()
	for i := 0; i < p.Order; i++ {
		dc.LineTo(p.X[i], p.Y[i])
	}
	dc.ClosePath()
	dc.Fill()
}

func (p *Polygon) SVG(attrs string) string {
	ret := fmt.Sprintf(
		"<polygon %s points=\"",
		attrs)
	points := make([]string, len(p.X))
	for i := 0; i < len(p.X); i++ {
		points[i] = fmt.Sprintf("%f,%f", p.X[i], p.Y[i])
	}

	return ret + strings.Join(points, ",") + "\" />"
}

func (p *Polygon) Copy() Shape {
	a := *p
	a.X = make([]float64, p.Order)
	a.Y = make([]float64, p.Order)
	copy(a.X, p.X)
	copy(a.Y, p.Y)
	return &a
}

func (p *Polygon) Mutate(plane *Plane, temp float64) {
	p.mutateImpl(plane, temp, 10)
}

func (p *Polygon) mutateImpl(plane *Plane, temp float64, rollback int) {
	const m = 16
	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	scale := 16 * temp
	for {
		if rnd.Float64() < 0.25 {
			// Swap a point
			i := rnd.Intn(p.Order)
			j := rnd.Intn(p.Order)
			p.X[i], p.Y[i], p.X[j], p.Y[j] = p.X[j], p.Y[j], p.X[i], p.Y[i]
			if p.Valid() {
				break
			}
			if rollback > 0 {
				p.X[i], p.Y[i], p.X[j], p.Y[j] = p.X[j], p.Y[j], p.X[i], p.Y[i]
				rollback -= 1
			}
		} else {
			// Move a point
			i := rnd.Intn(p.Order)
			xsave, ysave := p.X[i], p.Y[i]
			p.X[i] = clamp(p.X[i]+rnd.NormFloat64()*scale, -m, float64(w-1+m))
			p.Y[i] = clamp(p.Y[i]+rnd.NormFloat64()*scale, -m, float64(h-1+m))
			if p.Valid() {
				break
			}
			if rollback > 0 {
				p.X[i], p.Y[i] = xsave, ysave
				rollback -= 1
			}
		}
	}
}

func (p *Polygon) Valid() bool {
	if !p.Convex {
		return true
	}
	var sign bool
	for a := 0; a < p.Order; a++ {
		i := (a + 0) % p.Order
		j := (a + 1) % p.Order
		k := (a + 2) % p.Order
		c := cross3(p.X[i], p.Y[i], p.X[j], p.Y[j], p.X[k], p.Y[k])
		if a == 0 {
			sign = c > 0
		} else if c > 0 != sign {
			return false
		}
	}
	return true
}

func cross3(x1, y1, x2, y2, x3, y3 float64) float64 {
	dx1 := x2 - x1
	dy1 := y2 - y1
	dx2 := x3 - x2
	dy2 := y3 - y2
	return dx1*dy2 - dy1*dx2
}

func (p *Polygon) Rasterize(rc *RasterContext) []Scanline {
	var path raster.Path
	for i := 0; i <= p.Order; i++ {
		f := fixp(p.X[i%p.Order], p.Y[i%p.Order])
		if i == 0 {
			path.Start(f)
		} else {
			path.Add1(f)
		}
	}
	return fillPath(rc, path)
}
