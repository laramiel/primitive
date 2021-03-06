package shape

import (
	"fmt"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

// Quadratic represents a single quadratic bezier
type Quadratic struct {
	X1, Y1       float64
	X2, Y2       float64
	X3, Y3       float64
	Width        float64
	MinLineWidth float64
	MaxLineWidth float64
	MinArcLength float64
}

func NewQuadratic() *Quadratic {
	q := &Quadratic{}
	q.MaxLineWidth = 1.0 / 2
	q.MinLineWidth = 0.2
	q.MinArcLength = 5
	return q
}

func (q *Quadratic) Init(plane *Plane) {
	rnd := plane.Rnd
	q.X1 = randomW(plane)
	q.Y1 = randomH(plane)
	q.X2 = q.X1 + rnd.Float64()*40 - 20
	q.Y2 = q.Y1 + rnd.Float64()*40 - 20
	q.X3 = q.X2 + rnd.Float64()*40 - 20
	q.Y3 = q.Y2 + rnd.Float64()*40 - 20
	q.Width = 1.0 / 2
	q.mutateImpl(plane, 1.0, 2, ActionAny)
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

func (q *Quadratic) Mutate(plane *Plane, temp float64) {
	q.mutateImpl(plane, temp, 10, ActionAny)
}

func (q *Quadratic) mutateImpl(plane *Plane, temp float64, rollback int, actions ActionType) {
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
		switch rnd.Intn(6) {
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
		case 2:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := rnd.NormFloat64() * scale
			b := rnd.NormFloat64() * scale
			q.X3 = clamp(q.X3+a, -m, float64(w-1+m))
			q.Y3 = clamp(q.Y3+b, -m, float64(h-1+m))
		case 3: // Width
			q.Width = clamp(q.Width+rnd.NormFloat64()*temp, q.MinLineWidth, q.MaxLineWidth)

		case 4: // Translate
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

		case 5: // Rotate
			if (actions & ActionRotate) == 0 {
				continue
			}
			cx := (q.X1 + q.X2 + q.X3) / 3
			cy := (q.Y1 + q.Y2 + q.Y3) / 3
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

func (q *Quadratic) Valid() bool {
	dx12 := int(q.X1 - q.X2)
	dy12 := int(q.Y1 - q.Y2)
	d12 := dx12*dx12 + dy12*dy12

	dx23 := int(q.X2 - q.X3)
	dy23 := int(q.Y2 - q.Y3)
	d23 := dx23*dx23 + dy23*dy23

	return d12 > 1 && d23 > 1 && q.arcLength() > q.MinArcLength
}

func (q *Quadratic) arcLength() float64 {
	// closed form solution to quadratic arcLength from stack overflow:
	// https://math.stackexchange.com/questions/12186/arc-length-of-bezier-curves
	xv := 2 * (q.X2 - q.X1)
	yv := 2 * (q.Y2 - q.Y1)
	xw := q.X3 - 2*q.X2 + q.X1
	yw := q.Y3 - 2*q.Y2 + q.Y1
	uu := 4 * (xw*xw + yw*yw)
	if uu < 0.0001 {
		return math.Sqrt((q.X3-q.X1)*(q.X3-q.X1) + (q.Y3-q.Y1)*(q.Y3-q.Y1))
	}
	vv := 4 * (xv*xw + yv*yw)
	ww := xv*xv + yv*yv
	t1 := 2.0 * math.Sqrt(uu*(uu+vv+ww))
	t2 := 2*uu + vv
	t3 := vv*vv - 4*uu*ww
	t4 := 2.0 * math.Sqrt(uu*ww)
	return ((t1*t2 - t3*math.Log(t2+t1) - (vv*t4 - t3*math.Log(vv+t4))) / (8 * math.Pow(uu, 1.5)))
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
