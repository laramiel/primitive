package shape

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

// Rectangle represents a rectangular shape.
type Rectangle struct {
	X1, Y1 int
	X2, Y2 int
}

func NewRectangle() *Rectangle {
	return &Rectangle{}
}

func (r *Rectangle) Init(plane *Plane) {
	rnd := plane.Rnd
	r.X1 = rnd.Intn(plane.W)
	r.Y1 = rnd.Intn(plane.H)
	r.X2 = clampInt(r.X1+rnd.Intn(32)+1, 0, plane.W-1)
	r.Y2 = clampInt(r.Y1+rnd.Intn(32)+1, 0, plane.H-1)
	r.mutateImpl(plane, 1.0, 2, ActionAny)
}

func (r *Rectangle) bounds() (x1, y1, x2, y2 int) {
	x1, y1 = r.X1, r.Y1
	x2, y2 = r.X2, r.Y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return
}

func (r *Rectangle) Draw(dc *gg.Context, temp float64) {
	x1, y1, x2, y2 := r.bounds()
	dc.DrawRectangle(float64(x1), float64(y1), float64(x2-x1+1), float64(y2-y1+1))
	dc.Fill()
}

func (r *Rectangle) SVG(attrs string) string {
	x1, y1, x2, y2 := r.bounds()
	w := x2 - x1 + 1
	h := y2 - y1 + 1
	return fmt.Sprintf(
		"<rect %s x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" />",
		attrs, x1, y1, w, h)
}

func (r *Rectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *Rectangle) Mutate(plane *Plane, temp float64) {
	r.mutateImpl(plane, temp, 10, ActionAny)
}

func (r *Rectangle) mutateImpl(plane *Plane, temp float64, rollback int, actions ActionType) {
	if actions == ActionNone {
		return
	}

	const R = math.Pi / 4.0
	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	scale := 16 * temp
	save := *r
	for {
		switch rnd.Intn(5) {
		case 0: // Mutate
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := int(rnd.NormFloat64() * scale)
			b := int(rnd.NormFloat64() * scale)
			r.X1 = clampInt(r.X1+a, 0, w-1)
			r.Y1 = clampInt(r.Y1+b, 0, h-1)
		case 1:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := int(rnd.NormFloat64() * scale)
			b := int(rnd.NormFloat64() * scale)
			r.X2 = clampInt(r.X2+a, 0, w-1)
			r.Y2 = clampInt(r.Y2+b, 0, h-1)
		case 2: // Translate
			if (actions & ActionTranslate) == 0 {
				continue
			}
			a := int(rnd.NormFloat64() * scale)
			b := int(rnd.NormFloat64() * scale)
			r.X1 = clampInt(r.X1+a, 0, w-1)
			r.Y1 = clampInt(r.Y1+b, 0, h-1)
			r.X2 = clampInt(r.X2+a, 0, w-1)
			r.Y2 = clampInt(r.Y2+b, 0, h-1)
		case 3: // Move
			if (actions & ActionTranslate) == 0 {
				continue
			}
			a := int(rnd.NormFloat64() * scale)
			r.X1 = clampInt(r.X1+a, 0, w-1)
			r.Y1 = clampInt(r.Y1+a, 0, h-1)
			r.X2 = clampInt(r.X2+a, 0, w-1)
			r.Y2 = clampInt(r.Y2+a, 0, h-1)
		}
		if r.Valid() {
			break
		}
		if rollback > 0 {
			*r = save
			rollback -= 1
		}
	}
}

func (r *Rectangle) Valid() bool {
	a, b := r.X1-r.X2, r.Y1-r.Y2
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	return a > 2 && b > 2
}

func (r *Rectangle) Rasterize(rc *RasterContext) []Scanline {
	x1, y1, x2, y2 := r.bounds()
	lines := rc.Lines[:0]
	for y := y1; y <= y2; y++ {
		lines = append(lines, Scanline{y, x1, x2, 0xffff})
	}
	return lines
}

// RotatedRectangle represents a rotated rectangular shape
type RotatedRectangle struct {
	X, Y   int
	Sx, Sy int
	// Angle of rotation of the rectangle.
	Angle int
}

func NewRotatedRectangle() *RotatedRectangle {
	return &RotatedRectangle{}
}

func (r *RotatedRectangle) Init(plane *Plane) {
	rnd := plane.Rnd
	r.X = rnd.Intn(plane.W)
	r.Y = rnd.Intn(plane.H)
	r.Sx = rnd.Intn(32) + 1
	r.Sy = rnd.Intn(32) + 1
	r.Angle = rnd.Intn(360)
	r.mutateImpl(plane, 1.0, 1, ActionAny)
}

func (r *RotatedRectangle) Draw(dc *gg.Context, scale float64) {
	sx, sy := float64(r.Sx), float64(r.Sy)
	dc.Push()
	dc.Translate(float64(r.X), float64(r.Y))
	dc.Rotate(radians(float64(r.Angle)))
	dc.DrawRectangle(-sx/2, -sy/2, sx, sy)
	dc.Pop()
	dc.Fill()
}

func (r *RotatedRectangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<g transform=\"translate(%d %d) rotate(%d) scale(%d %d)\"><rect %s x=\"-0.5\" y=\"-0.5\" width=\"1\" height=\"1\" /></g>",
		r.X, r.Y, r.Angle, r.Sx, r.Sy, attrs)
}

func (r *RotatedRectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *RotatedRectangle) Mutate(plane *Plane, temp float64) {
	r.mutateImpl(plane, temp, 10, ActionAny)
}

func (r *RotatedRectangle) mutateImpl(plane *Plane, temp float64, rollback int, actions ActionType) {
	if actions == ActionNone {
		return
	}

	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	scale := 16 * temp
	save := *r
	for {
		a := int(rnd.NormFloat64() * scale)
		b := int(rnd.NormFloat64() * scale)
		switch rnd.Intn(3) {
		case 0: // Move Origin
			if (actions & ActionTranslate) == 0 {
				continue
			}
			r.X = clampInt(r.X+a, 0, w-1)
			r.Y = clampInt(r.Y+b, 0, h-1)
		case 1: // Resize
			if (actions & ActionScale) == 0 {
				continue
			}
			r.Sx = clampInt(r.Sx+a, 1, w-1)
			r.Sy = clampInt(r.Sy+b, 1, h-1)
		case 2: // Rotate
			if (actions & ActionRotate) == 0 {
				continue
			}
			r.Angle = r.Angle + a + a
		}
		if r.Valid() {
			break
		}
		if rollback > 0 {
			*r = save
			rollback -= 1
		}
	}
}

func (r *RotatedRectangle) Valid() bool {
	a, b := r.Sx, r.Sy
	if a < b {
		a, b = b, a
	}
	aspect := float64(a) / float64(b)
	return aspect <= 5
}

func (r *RotatedRectangle) Rasterize(rc *RasterContext) []Scanline {
	w := rc.W
	h := rc.H
	sx, sy := float64(r.Sx), float64(r.Sy)
	angle := radians(float64(r.Angle))
	rx1, ry1 := rotate(-sx/2, -sy/2, angle)
	rx2, ry2 := rotate(sx/2, -sy/2, angle)
	rx3, ry3 := rotate(sx/2, sy/2, angle)
	rx4, ry4 := rotate(-sx/2, sy/2, angle)
	x1, y1 := int(rx1)+r.X, int(ry1)+r.Y
	x2, y2 := int(rx2)+r.X, int(ry2)+r.Y
	x3, y3 := int(rx3)+r.X, int(ry3)+r.Y
	x4, y4 := int(rx4)+r.X, int(ry4)+r.Y
	miny := minInt(y1, minInt(y2, minInt(y3, y4)))
	maxy := maxInt(y1, maxInt(y2, maxInt(y3, y4)))
	n := maxy - miny + 1
	min := make([]int, n)
	max := make([]int, n)
	for i := range min {
		min[i] = w
	}
	xs := []int{x1, x2, x3, x4, x1}
	ys := []int{y1, y2, y3, y4, y1}
	// TODO: this could be better probably
	for i := 0; i < 4; i++ {
		x, y := float64(xs[i]), float64(ys[i])
		dx, dy := float64(xs[i+1]-xs[i]), float64(ys[i+1]-ys[i])
		count := int(math.Sqrt(dx*dx+dy*dy)) * 2
		for j := 0; j < count; j++ {
			t := float64(j) / float64(count-1)
			xi := int(x + dx*t)
			yi := int(y+dy*t) - miny
			min[yi] = minInt(min[yi], xi)
			max[yi] = maxInt(max[yi], xi)
		}
	}
	lines := rc.Lines[:0]
	for i := 0; i < n; i++ {
		y := miny + i
		if y < 0 || y >= h {
			continue
		}
		a := maxInt(min[i], 0)
		b := minInt(max[i], w-1)
		if b >= a {
			lines = append(lines, Scanline{y, a, b, 0xffff})
		}
	}
	return lines
}
