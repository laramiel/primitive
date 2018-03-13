package primitive

import (
	"fmt"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

const LineMutate = 2 // 3 for line width

type Line struct {
	Worker *Worker
	X1, Y1 float64
	X2, Y2 float64
	Width  float64
    MaxLineWidth float64
}

func NewRandomLine(worker *Worker) *Line {
	l := &Line{}
	l.Init(worker)
	return l
}

func (q *Line) Init(worker *Worker) {	
	rnd := worker.Rnd
	q.Worker = worker
	q.X1 = rnd.Float64() * float64(worker.W)
	q.Y1 = rnd.Float64() * float64(worker.H)
	q.X2 = rnd.Float64() * float64(worker.W)
	q.Y2 = rnd.Float64() * float64(worker.H)
	q.Width = 1.0 / 2
	q.Mutate()
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

func (q *Line) Mutate() {
	const m = 16
	w := q.Worker.W
	h := q.Worker.H
	rnd := q.Worker.Rnd
	for {
		switch rnd.Intn(LineMutate) {
		case 0:
			q.X1 = clamp(q.X1+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y1 = clamp(q.Y1+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 1:
			q.X2 = clamp(q.X2+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y2 = clamp(q.Y2+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 2:
			if q.Width != q.MaxLineWidth {
			q.Width = clamp(q.Width+rnd.NormFloat64(), 0.1, q.MaxLineWidth)
			}
		}
		if q.Valid() {
			break
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
	return (y2-y1) > 2 || (x2-x1) > 2
}

func (q *Line) Rasterize() []Scanline {
	var path raster.Path
	p1 := fixp(q.X1, q.Y1)
	p2 := fixp(q.X2, q.Y2)
	path.Start(p1)
	path.Add1(p2)
	width := fix(q.Width)
	return strokePath(q.Worker, path, width, raster.RoundCapper, raster.RoundJoiner)
}

type RadialLine struct {
	Line   Line
	CX, CY float64
}

func NewRadialLine(worker *Worker, cx, cy float64) *RadialLine {
	l := &RadialLine{}
	l.CX = cx
	l.CY = cy
	l.Line.MaxLineWidth = 1.0 / 2
	l.Init(worker)
	return l
}

func (l *RadialLine) Init(worker *Worker) {	
	rnd := worker.Rnd
	l.Line.Worker = worker
	l.Line.X1 = l.CX * float64(worker.W)
	l.Line.Y1 = l.CY * float64(worker.H)
	l.Line.X2 = rnd.Float64() * float64(worker.W)
	l.Line.Y2 = rnd.Float64() * float64(worker.H)
	l.Line.Width = 1.0 / 2
	l.Line.Mutate()
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

func (l *RadialLine) Mutate() {
	const MaxLineWidth = 4
	const m = 16
	w := l.Line.Worker.W
	h := l.Line.Worker.H
	rnd := l.Line.Worker.Rnd
	for {
		switch rnd.Intn(LineMutate) {
		case 0:
			// Move along radial point
			xd := l.Line.X2 - l.Line.X1
			yd := l.Line.Y2 - l.Line.Y1
			v := rnd.NormFloat64()
			l.Line.X1 = l.Line.X1 + v*xd
			l.Line.Y1 = l.Line.Y1 + v*yd

		case 1:
			// New radial point
			l.Line.X1 = l.CX * float64(w)
			l.Line.Y1 = l.CY * float64(h)
			l.Line.X2 = clamp(l.Line.X2+rnd.NormFloat64()*16, -m, float64(w-1+m))
			l.Line.Y2 = clamp(l.Line.Y2+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 2:
			l.Line.Width = clamp(l.Line.Width+rnd.NormFloat64(), 1, MaxLineWidth)
		}
		if l.Valid() {
			break
		}
	}
}

func (l *RadialLine) Valid() bool {
	return l.Line.Valid()
}

func (l *RadialLine) Rasterize() []Scanline {
	return l.Line.Rasterize()
}
