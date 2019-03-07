package shape

import (
	"fmt"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Stamp struct {
	X1, Y1 float64
	X, Y   []float64
	// TODO: Angle, Scale
}

func NewStamp() *Stamp {
	return &Stamp{}
}

func (s *Stamp) Init(plane *Plane) {
	s.X1 = randomW(plane)
	s.Y1 = randomH(plane)
	s.mutateImpl(plane, 1.0, 1, ActionAny)
}

func (s *Stamp) Draw(dc *gg.Context, scale float64) {
	dc.NewSubPath()
	for i := 0; i < len(s.X); i++ {
		dc.LineTo(s.X1+s.X[i], s.Y1+s.Y[i])
	}
	dc.ClosePath()
	dc.Fill()
}

func (s *Stamp) SVG(attrs string) string {
	ret := fmt.Sprintf(
		"<polygon %s points=\"",
		attrs)
	points := make([]string, len(s.X))
	for i := 0; i < len(s.X); i++ {
		points[i] = fmt.Sprintf("%f,%f", s.X1+s.X[i], s.Y1+s.Y[i])
	}

	return ret + strings.Join(points, ",") + "\" />"
}

func (s *Stamp) Copy() Shape {
	a := *s
	a.X = make([]float64, len(s.X))
	a.Y = make([]float64, len(s.Y))
	copy(a.X, s.X)
	copy(a.Y, s.Y)
	return &a
}

func (s *Stamp) Mutate(plane *Plane, temp float64) {
	s.mutateImpl(plane, temp, 10, ActionAny)
}

func (s *Stamp) mutateImpl(plane *Plane, temp float64, rollback int, actions ActionType) {
	if actions == ActionNone {
		return
	}

	const R = math.Pi / 4.0
	const m = 16
	w := float64(plane.W - 1 + m)
	h := float64(plane.H - 1 + m)

	rnd := plane.Rnd
	scale := 16 * temp

	// Move center point
	a := rnd.NormFloat64() * scale
	b := rnd.NormFloat64() * scale
	s.X1 = clamp(s.X1+a, -m, w)
	s.Y1 = clamp(s.Y1+b, -m, h)
}

func (s *Stamp) Rasterize(rc *RasterContext) []Scanline {
	var path raster.Path
	path.Start(fixp(s.X1+s.X[0], s.Y1+s.Y[0]))
	for i := 1; i < len(s.X); i++ {
		path.Add1(fixp(s.X1+s.X[i], s.Y1+s.Y[i]))
	}
	path.Add1(fixp(s.X1+s.X[0], s.Y1+s.Y[0]))
	return fillPath(rc, path)
}
