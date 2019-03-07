package shape

import (
	"fmt"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

// Cubic represents a single cubic bezier curve
type Cubic struct {
	X1, Y1       float64
	X2, Y2       float64
	X3, Y3       float64
	X4, Y4       float64
	Width        float64
	MinLineWidth float64
	MaxLineWidth float64
	MinArcLength float64
}

func NewCubic() *Cubic {
	q := &Cubic{}
	q.MaxLineWidth = 1.0 / 2
	q.MinLineWidth = 0.2
	q.MinArcLength = 5
	return q
}

func (q *Cubic) Init(plane *Plane) {
	rnd := plane.Rnd
	q.X1 = randomW(plane)
	q.Y1 = randomH(plane)
	q.X2 = q.X1 + rnd.Float64()*40 - 20
	q.Y2 = q.Y1 + rnd.Float64()*40 - 20
	q.X3 = q.X2 + rnd.Float64()*40 - 20
	q.Y3 = q.Y2 + rnd.Float64()*40 - 20
	q.X4 = q.X3 + rnd.Float64()*40 - 20
	q.Y4 = q.Y3 + rnd.Float64()*40 - 20
	q.Width = 1.0 / 2
	q.mutateImpl(plane, 1.0, 2, ActionAny)
}

func (q *Cubic) Draw(dc *gg.Context, scale float64) {
	dc.MoveTo(q.X1, q.Y1)
	dc.CubicTo(q.X2, q.Y2, q.X3, q.Y3, q.X4, q.Y4)
	dc.SetLineWidth(q.Width * scale)
	dc.Stroke()
}

func (q *Cubic) SVG(attrs string) string {
	// TODO: this is a little silly
	attrs = strings.Replace(attrs, "fill", "stroke", -1)
	return fmt.Sprintf(
		"<path %s fill=\"none\" d=\"M %f %f Q %f %f, %f %f, %f %f\" stroke-width=\"%f\" />",
		attrs, q.X1, q.Y1, q.X2, q.Y2, q.X3, q.Y3, q.X4, q.Y4, q.Width)
}

func (q *Cubic) Copy() Shape {
	a := *q
	return &a
}

func (q *Cubic) Mutate(plane *Plane, temp float64) {
	q.mutateImpl(plane, temp, 10, ActionAny)
}

func (q *Cubic) mutateImpl(plane *Plane, temp float64, rollback int, actions ActionType) {
	if actions == ActionNone {
		return
	}
	const R = math.Pi / 4.0
	const m = 16
	w := plane.W
	h := plane.H
	rnd := plane.Rnd
	scale := temp * 16
	save := *q
	for {
		switch rnd.Intn(7) {
		case 0: // Mutate
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
		case 2:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			q.X3 = clamp(q.X3+a, -m, float64(w-1+m))
			q.Y3 = clamp(q.Y3+b, -m, float64(h-1+m))
		case 3:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			q.X4 = clamp(q.X4+a, -m, float64(w-1+m))
			q.Y4 = clamp(q.Y4+b, -m, float64(h-1+m))

		case 4: // Width
			q.Width = clamp(q.Width+rnd.NormFloat64()*temp, q.MinLineWidth, q.MaxLineWidth)

		case 5: // Translate
			if (actions & ActionTranslate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			q.X1 = clamp(q.X1+a, -m, float64(w-1+m))
			q.Y1 = clamp(q.Y1+b, -m, float64(h-1+m))
			q.X2 = clamp(q.X2+a, -m, float64(w-1+m))
			q.Y2 = clamp(q.Y2+b, -m, float64(h-1+m))
			q.X3 = clamp(q.X3+a, -m, float64(w-1+m))
			q.Y3 = clamp(q.Y3+b, -m, float64(h-1+m))
			q.X4 = clamp(q.X4+a, -m, float64(w-1+m))
			q.Y4 = clamp(q.Y4+b, -m, float64(h-1+m))

		case 6: // Rotate
			if (actions & ActionRotate) == 0 {
				continue
			}
			cx := (q.X1 + q.X2 + q.X3 + q.X4) / 4
			cy := (q.Y1 + q.Y2 + q.Y3 + q.Y4) / 4
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
			a, b = rotateAbout(q.X3, q.Y3, cx, cy, cos, sin)
			q.X3 = clamp(a, -m, float64(w-1+m))
			q.Y3 = clamp(b, -m, float64(h-1+m))
			a, b = rotateAbout(q.X4, q.Y4, cx, cy, cos, sin)
			q.X4 = clamp(a, -m, float64(w-1+m))
			q.Y4 = clamp(b, -m, float64(h-1+m))
			// TODO: Scale
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

func (q *Cubic) Valid() bool {
	dx12 := int(q.X1 - q.X2)
	dy12 := int(q.Y1 - q.Y2)
	d12 := dx12*dx12 + dy12*dy12

	dx23 := int(q.X2 - q.X3)
	dy23 := int(q.Y2 - q.Y3)
	d23 := dx23*dx23 + dy23*dy23

	dx34 := int(q.X3 - q.X4)
	dy34 := int(q.Y3 - q.Y4)
	d34 := dx34*dx34 + dy34*dy34

	return d12 > 1 && d23 > 1 && d34 > 1 && q.arcLength() > q.MinArcLength
}

func (q *Cubic) arcLength() float64 {
	d := 0.0
	// The actual answer requires that we do numerical integration
	// along the length of the spline; or we can approximate it by
	// sampling N points and adding the length of each segment.
	x := q.X1
	y := q.Y1
	const k = 1.0 / 48.0
	for t := k; t < 1.0; t += k {
		mt := 1.0 - t
		t2 := t * t
		mt2 := mt * mt
		a := mt2 * mt
		b := mt2 * t * 3
		c := mt * t2 * 3
		d := t * t2

		nx := a*q.X1 + b*q.X2 + c*q.X3 + d*q.X4
		ny := a*q.Y1 + b*q.Y2 + c*q.Y3 + d*q.Y4

		dx := nx - x
		dy := ny - y

		d += math.Sqrt(dx*dx + dy*dy)
		x = nx
		y = ny
	}
	return d
}

func (q *Cubic) Rasterize(rc *RasterContext) []Scanline {
	var path raster.Path
	p1 := fixp(q.X1, q.Y1)
	p2 := fixp(q.X2, q.Y2)
	p3 := fixp(q.X3, q.Y3)
	p4 := fixp(q.X4, q.Y4)
	path.Start(p1)
	path.Add3(p2, p3, p4)
	width := fix(q.Width)
	return strokePath(rc, path, width, raster.RoundCapper, raster.RoundJoiner)
}
