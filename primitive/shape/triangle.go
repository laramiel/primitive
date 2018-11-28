package shape

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

type Triangle struct {
	X1, Y1  int
	X2, Y2  int
	X3, Y3  int
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
	t.X1 = rnd.Intn(plane.W)
	t.Y1 = rnd.Intn(plane.H)
	t.X2 = t.X1 + rnd.Intn(31) - 15
	t.Y2 = t.Y1 + rnd.Intn(31) - 15
	t.X3 = t.X1 + rnd.Intn(31) - 15
	t.Y3 = t.Y1 + rnd.Intn(31) - 15
	t.mutateImpl(plane, 1.0, 2)
}

func (t *Triangle) Draw(dc *gg.Context, scale float64) {
	dc.LineTo(float64(t.X1), float64(t.Y1))
	dc.LineTo(float64(t.X2), float64(t.Y2))
	dc.LineTo(float64(t.X3), float64(t.Y3))
	dc.ClosePath()
	dc.Fill()
}

func (t *Triangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<polygon %s points=\"%d,%d %d,%d %d,%d\" />",
		attrs, t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3)
}

func (t *Triangle) Copy() Shape {
	a := *t
	return &a
}

func rotateAngle(x, y int, x0, y0 int, cos, sin float64) (int, int) {
	xd := float64(x - x0)
	yd := float64(y - y0)
	return int(xd*cos -yd*sin + float64(x0)), int(xd * sin + yd * cos + float64(y0))
}

func (t *Triangle) Mutate(plane *Plane, temp float64) {
	t.mutateImpl(plane, temp, 100)
}

func (t *Triangle) mutateImpl(plane *Plane, temp float64, rollback int) {
	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	const m = 16
	scale := 16 * temp
	save := *t
	for {
		switch rnd.Intn(5) {
			// Move.
		case 0:
			t.X1 = clampInt(t.X1+int(rnd.NormFloat64()*scale), -m, w-1+m)
			t.Y1 = clampInt(t.Y1+int(rnd.NormFloat64()*scale), -m, h-1+m)
		case 1:
			t.X2 = clampInt(t.X2+int(rnd.NormFloat64()*scale), -m, w-1+m)
			t.Y2 = clampInt(t.Y2+int(rnd.NormFloat64()*scale), -m, h-1+m)
		case 2:
			t.X3 = clampInt(t.X3+int(rnd.NormFloat64()*scale), -m, w-1+m)
			t.Y3 = clampInt(t.Y3+int(rnd.NormFloat64()*scale), -m, h-1+m)

		case 3: // Shift
			a := int(rnd.NormFloat64()*scale)
			b := int(rnd.NormFloat64()*scale)
			t.X1 = clampInt(t.X1+a, -m, w-1+m)
			t.Y1 = clampInt(t.Y1+b, -m, h-1+m)
			t.X2 = clampInt(t.X2+a, -m, w-1+m)
			t.Y2 = clampInt(t.Y2+b, -m, h-1+m)
			t.X3 = clampInt(t.X3+a, -m, w-1+m)
			t.Y3 = clampInt(t.Y3+b, -m, h-1+m)
		case 4: // Rotate 
			cx := (t.X1 + t.X2 + t.X3) / 3
			cy := (t.Y1 + t.Y2 + t.Y3) / 3
			angle := rnd.NormFloat64()*scale* math.Pi / 4
			cos := math.Cos(angle)
			sin := math.Sin(angle)

			var a, b int
			a, b = rotateAngle(t.X1, t.Y1, cx, cy, cos, sin) 
			t.X1 = clampInt(a, -m, w-1+m)
			t.Y1 = clampInt(b, -m, h-1+m)
			a, b = rotateAngle(t.X2, t.Y2, cx, cy, cos, sin) 
			t.X2 = clampInt(a, -m, w-1+m)
			t.Y2 = clampInt(b, -m, h-1+m)
			a, b = rotateAngle(t.X2, t.Y2, cx, cy, cos, sin) 
			t.X3 = clampInt(a, -m, w-1+m)
			t.Y3 = clampInt(b, -m, h-1+m)
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
		if a > t.MaxArea {
			return false
		}
	}

	const minDegrees = 15
	var a1, a2, a3 float64
	{
		x1 := float64(t.X2 - t.X1)
		y1 := float64(t.Y2 - t.Y1)
		x2 := float64(t.X3 - t.X1)
		y2 := float64(t.Y3 - t.Y1)
		d1 := math.Sqrt(x1*x1 + y1*y1)
		d2 := math.Sqrt(x2*x2 + y2*y2)
		x1 /= d1
		y1 /= d1
		x2 /= d2
		y2 /= d2
		a1 = degrees(math.Acos(x1*x2 + y1*y2))
	}
	{
		x1 := float64(t.X1 - t.X2)
		y1 := float64(t.Y1 - t.Y2)
		x2 := float64(t.X3 - t.X2)
		y2 := float64(t.Y3 - t.Y2)
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
	buf := rc.Lines[:0]
	lines := rasterizeTriangle(t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3, buf)
	return cropScanlines(lines, rc.W, rc.H)
}

func rasterizeTriangle(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	if y1 > y3 {
		x1, x3 = x3, x1
		y1, y3 = y3, y1
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y2 > y3 {
		x2, x3 = x3, x2
		y2, y3 = y3, y2
	}
	if y2 == y3 {
		return rasterizeTriangleBottom(x1, y1, x2, y2, x3, y3, buf)
	} else if y1 == y2 {
		return rasterizeTriangleTop(x1, y1, x2, y2, x3, y3, buf)
	} else {
		x4 := x1 + int((float64(y2-y1)/float64(y3-y1))*float64(x3-x1))
		y4 := y2
		buf = rasterizeTriangleBottom(x1, y1, x2, y2, x4, y4, buf)
		buf = rasterizeTriangleTop(x2, y2, x4, y4, x3, y3, buf)
		return buf
	}
}

func rasterizeTriangleBottom(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	s1 := float64(x2-x1) / float64(y2-y1)
	s2 := float64(x3-x1) / float64(y3-y1)
	ax := float64(x1)
	bx := float64(x1)
	for y := y1; y <= y2; y++ {
		a := int(ax)
		b := int(bx)
		ax += s1
		bx += s2
		if a > b {
			a, b = b, a
		}
		buf = append(buf, Scanline{y, a, b, 0xffff})
	}
	return buf
}

func rasterizeTriangleTop(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	s1 := float64(x3-x1) / float64(y3-y1)
	s2 := float64(x3-x2) / float64(y3-y2)
	ax := float64(x3)
	bx := float64(x3)
	for y := y3; y > y1; y-- {
		ax -= s1
		bx -= s2
		a := int(ax)
		b := int(bx)
		if a > b {
			a, b = b, a
		}
		buf = append(buf, Scanline{y, a, b, 0xffff})
	}
	return buf
}
