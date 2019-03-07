package shape

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

// Triangle represents a triangular shape
type Triangle struct {
	X1, Y1  float64
	X2, Y2  float64
	X3, Y3  float64
	MaxArea int
}

func NewTriangle() *Triangle {
	return &Triangle{}
}

func NewMaxAreaTriangle(area int) *Triangle {
	t := &Triangle{}
	t.MaxArea = area
	return t
}

func (t *Triangle) Init(plane *Plane) {
	rnd := plane.Rnd
	t.X1 = randomW(plane)
	t.Y1 = randomH(plane)
	t.X2 = t.X1 + rnd.NormFloat64()*32
	t.Y2 = t.Y1 + rnd.NormFloat64()*32
	t.X3 = t.X1 + rnd.NormFloat64()*32
	t.Y3 = t.Y1 + rnd.NormFloat64()*32
	t.mutateImpl(plane, 1.0, 2, ActionAny)
}

func (t *Triangle) Draw(dc *gg.Context, scale float64) {
	dc.LineTo(t.X1, t.Y1)
	dc.LineTo(t.X2, t.Y2)
	dc.LineTo(t.X3, t.Y3)
	dc.ClosePath()
	dc.Fill()
}

func (t *Triangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<polygon %s points=\"%f,%f %f,%f %f,%f\" />",
		attrs, t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3)
}

func (t *Triangle) Copy() Shape {
	a := *t
	return &a
}

func (t *Triangle) Mutate(plane *Plane, temp float64) {
	t.mutateImpl(plane, temp, 100, ActionAny)
}

func (t *Triangle) mutateImpl(plane *Plane, temp float64, rollback int, actions ActionType) {
	if actions == ActionNone {
		return
	}

	const R = math.Pi / 4.0
	const m = 16
	w := float64(plane.W - 1 + m)
	h := float64(plane.H - 1 + m)
	rnd := plane.Rnd
	scale := 16 * temp
	save := *t

	for {
		switch rnd.Intn(5) {
		// Mutate.
		case 0:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			t.X1 = clamp(t.X1+a, -m, w)
			t.Y1 = clamp(t.Y1+b, -m, h)
		case 1:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			t.X2 = clamp(t.X2+a, -m, w)
			t.Y2 = clamp(t.Y2+b, -m, h)
		case 2:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			t.X3 = clamp(t.X3+a, -m, w)
			t.Y3 = clamp(t.Y3+b, -m, h)

		case 3: // Translate
			if (actions & ActionTranslate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			t.X1 = clamp(t.X1+a, -m, w)
			t.Y1 = clamp(t.Y1+b, -m, h)
			t.X2 = clamp(t.X2+a, -m, w)
			t.Y2 = clamp(t.Y2+b, -m, h)
			t.X3 = clamp(t.X3+a, -m, w)
			t.Y3 = clamp(t.Y3+b, -m, h)

		case 4: // Rotate
			if (actions & ActionRotate) == 0 {
				continue
			}
			cx := (t.X1 + t.X2 + t.X3) / 3
			cy := (t.Y1 + t.Y2 + t.Y3) / 3
			theta := rnd.NormFloat64() * temp * R
			cos := math.Cos(theta)
			sin := math.Sin(theta)

			var a, b float64
			a, b = rotateAbout(t.X1, t.Y1, cx, cy, cos, sin)
			t.X1 = clamp(a, -m, w)
			t.Y1 = clamp(b, -m, h)
			a, b = rotateAbout(t.X2, t.Y2, cx, cy, cos, sin)
			t.X2 = clamp(a, -m, w)
			t.Y2 = clamp(b, -m, h)
			a, b = rotateAbout(t.X2, t.Y2, cx, cy, cos, sin)
			t.X3 = clamp(a, -m, w-1+m)
			t.Y3 = clamp(b, -m, h-1+m)
		}
		if t.Valid() {
			break
		}
		if rollback > 0 {
			*t = save
			rollback -= 1
		}
	}
}

func (t *Triangle) Valid() bool {
	if t.MaxArea > 0 {
		// compute the area.
		a := (t.X1*(t.Y2-t.Y3) + t.X2*(t.Y3-t.Y1) + t.X3*(t.Y1-t.Y2)) / 2
		if a < 0 {
			a = -a
		}
		if a > float64(t.MaxArea) {
			return false
		}
	}

	const minDegrees = 15
	var a1, a2, a3 float64
	{
		x1 := t.X2 - t.X1
		y1 := t.Y2 - t.Y1
		x2 := t.X3 - t.X1
		y2 := t.Y3 - t.Y1
		d1 := math.Sqrt(x1*x1 + y1*y1)
		d2 := math.Sqrt(x2*x2 + y2*y2)
		x1 /= d1
		y1 /= d1
		x2 /= d2
		y2 /= d2
		a1 = degrees(math.Acos(x1*x2 + y1*y2))
	}
	{
		x1 := t.X1 - t.X2
		y1 := t.Y1 - t.Y2
		x2 := t.X3 - t.X2
		y2 := t.Y3 - t.Y2
		d1 := math.Sqrt(x1*x1 + y1*y1)
		d2 := math.Sqrt(x2*x2 + y2*y2)
		x1 /= d1
		y1 /= d1
		x2 /= d2
		y2 /= d2
		a2 = degrees(math.Acos(x1*x2 + y1*y2))
	}
	a3 = 180 - a1 - a2
	return a1 > minDegrees && a2 > minDegrees && a3 > minDegrees
}

func (t *Triangle) Rasterize(rc *RasterContext) []Scanline {
	var path raster.Path
	path.Start(fixp(t.X1, t.Y1))
	path.Add1(fixp(t.X2, t.Y2))
	path.Add1(fixp(t.X3, t.Y3))
	path.Add1(fixp(t.X1, t.Y1))
	return fillPath(rc, path)
}
