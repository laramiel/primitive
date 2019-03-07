package shape

import (
	"fmt"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

const LineMutate = 4 // 5 for line width

type Line struct {
	X1, Y1       float64
	X2, Y2       float64
	Width        float64
	MaxLineWidth float64
}

func NewLine() *Line {
	l := &Line{}
	l.MaxLineWidth = 1.0 / 2
	return l
}

func (q *Line) Init(plane *Plane) {
	q.X1 = randomW(plane)
	q.Y1 = randomH(plane)
	q.X2 = randomW(plane)
	q.Y2 = randomH(plane)
	q.Width = 1.0 / 2
	q.mutateImpl(plane, 1.0, 1, ActionAny)
}

func (q *Line) Draw(dc *gg.Context, scale float64) {
	dc.MoveTo(q.X1, q.Y1)
	dc.LineTo(q.X2, q.Y2)
	dc.SetLineWidth(q.Width * scale)
	dc.Stroke()
}

func (q *Line) SVG(attrs string) string {
	// TODO: this is a little silly
	attrs = strings.Replace(attrs, "fill", "stroke", -1)
	return fmt.Sprintf(
		"<path %s fill=\"none\" d=\"M %f %f L %f %f\" stroke-width=\"%f\" />",
		attrs, q.X1, q.Y1, q.X2, q.Y2, q.Width)
}

func (q *Line) Copy() Shape {
	a := *q
	return &a
}

func (q *Line) Mutate(plane *Plane, temp float64) {
	q.mutateImpl(plane, temp, 10, ActionAny)
}

func (q *Line) mutateImpl(plane *Plane, temp float64, rollback int, actions ActionType) {
	if actions == ActionNone {
		return
	}

	const R = math.Pi / 4.0
	const m = 16
	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	scale := 16 * temp
	save := *q
	for {
		switch rnd.Intn(LineMutate) {
		case 0: // Move
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			q.X1 = clamp(q.X1+a, -m, float64(w-1+m))
			q.Y1 = clamp(q.Y1+b, -m, float64(h-1+m))
		case 1:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			q.X2 = clamp(q.X2+a, -m, float64(w-1+m))
			q.Y2 = clamp(q.Y2+b, -m, float64(h-1+m))
		case 2: // Translate
			if (actions & ActionTranslate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			q.X1 = clamp(q.X1+a, -m, float64(w-1+m))
			q.Y1 = clamp(q.Y1+b, -m, float64(h-1+m))
			q.X2 = clamp(q.X2+a, -m, float64(w-1+m))
			q.Y2 = clamp(q.Y2+b, -m, float64(h-1+m))

		case 3: // Rotate
			if (actions & ActionRotate) == 0 {
				continue
			}
			cx := (q.X1 + q.X2) / 2
			cy := (q.Y1 + q.Y2) / 2
			theta := rnd.NormFloat64() * temp * R
			cos := math.Cos(theta)
			sin := math.Sin(theta)

			var a, b float64
			a, b = rotateAbout(q.X1, q.Y1, cx, cy, cos, sin)
			q.X1 = clamp(a, -m, float64(w-1+m))
			q.Y1 = clamp(b, -m, float64(h-1+m))
			a, b = rotateAbout(q.X2, q.Y2, cx, cy, cos, sin)
			q.X2 = clamp(a, -m, float64(w-1+m))
			q.Y2 = clamp(b, -m, float64(h-1+m))
		case 4: // Width
			if q.Width != q.MaxLineWidth {
				q.Width = clamp(q.Width+rnd.NormFloat64(), 0.1, q.MaxLineWidth)
			}
		}
		if q.Valid() {
			break
		}
		if rollback > 0 {
			*q = save
			rollback -= 1
		}
	}
}

func (q *Line) Valid() bool {
	x1, x2 := q.X1, q.X2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	y1, y2 := q.Y1, q.Y2
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return (y2-y1) > 1 || (x2-x1) > 1
}

func (q *Line) Rasterize(rc *RasterContext) []Scanline {
	var path raster.Path
	p1 := fixp(q.X1, q.Y1)
	p2 := fixp(q.X2, q.Y2)
	path.Start(p1)
	path.Add1(p2)
	width := fix(q.Width)
	return strokePath(rc, path, width, raster.RoundCapper, raster.RoundJoiner)
}

type RadialLine struct {
	CX, CY float64
	Line   Line
}

func NewRadialLine(cx, cy float64) *RadialLine {
	l := &RadialLine{}
	l.CX = cx
	l.CY = cy
	l.Line.MaxLineWidth = 1.0 / 2
	return l
}

func (l *RadialLine) Init(plane *Plane) {
	rnd := plane.Rnd
	l.Line.X1 = l.CX * float64(plane.W)
	l.Line.Y1 = l.CY * float64(plane.H)
	l.Line.X2 = rnd.Float64() * float64(plane.W)
	l.Line.Y2 = rnd.Float64() * float64(plane.H)
	l.Line.Width = 1.0 / 2
	l.mutateImpl(plane, 1.0, 1)
}

func (l *RadialLine) Draw(dc *gg.Context, scale float64) {
	l.Line.Draw(dc, scale)
}

func (l *RadialLine) SVG(attrs string) string {
	return l.Line.SVG(attrs)
}

func (l *RadialLine) Copy() Shape {
	a := *l
	return &a
}

func (l *RadialLine) Mutate(plane *Plane, temp float64) {
	l.mutateImpl(plane, temp, 10)
}

func (l *RadialLine) mutateImpl(plane *Plane, temp float64, rollback int) {
	const MaxLineWidth = 4
	const m = 16
	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	scale := 16 * temp
	save := *l
	for {
		switch rnd.Intn(LineMutate) {
		case 0:
			// Move along radial point
			xd := l.Line.X2 - l.Line.X1
			yd := l.Line.Y2 - l.Line.Y1
			v := rnd.NormFloat64() * temp
			l.Line.X1 = l.Line.X1 + v*xd
			l.Line.Y1 = l.Line.Y1 + v*yd

		case 1:
			// New radial point
			l.Line.X1 = l.CX * float64(w)
			l.Line.Y1 = l.CY * float64(h)
			l.Line.X2 = clamp(l.Line.X2+rnd.NormFloat64()*scale, -m, float64(w-1+m))
			l.Line.Y2 = clamp(l.Line.Y2+rnd.NormFloat64()*scale, -m, float64(h-1+m))
		case 2:
			l.Line.Width = clamp(l.Line.Width+rnd.NormFloat64(), 1, MaxLineWidth)
		}
		if l.Valid() {
			break
		}
		if rollback > 0 {
			*l = save
			rollback -= 1
		}
	}
}

func (l *RadialLine) Valid() bool {
	return l.Line.Valid()
}

func (l *RadialLine) Rasterize(rc *RasterContext) []Scanline {
	return l.Line.Rasterize(rc)
}
