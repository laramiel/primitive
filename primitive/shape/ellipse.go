package shape

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type EllipseType int

const (
	EllipseAny EllipseType = iota
	EllipseCircle
	EllipseCenteredCircle
	EllipseFixedRadius
)

// An Ellipse shape.
type Ellipse struct {
	X, Y        int
	Rx, Ry      int
	EllipseType EllipseType
	CX, CY      float64
	MaxRadius   int
}

func NewEllipse() *Ellipse {
	c := &Ellipse{}
	c.EllipseType = EllipseAny
	return c
}

func NewCircle() *Ellipse {
	c := &Ellipse{}
	c.EllipseType = EllipseCircle
	return c
}

func NewFixedCircle(r int) *Ellipse {
	c := &Ellipse{}
	c.EllipseType = EllipseFixedRadius
	c.Rx = r
	c.Ry = r
	return c
}

func NewCenteredCircle(x, y float64) *Ellipse {
	c := &Ellipse{}
	c.EllipseType = EllipseCenteredCircle
	c.CX = x
	c.CY = y
	return c
}

func (c *Ellipse) Init(plane *Plane) {
	rnd := plane.Rnd
	maxr := 32
	if c.MaxRadius > 0 && maxr > c.MaxRadius {
		maxr = c.MaxRadius - 1
	}
	if c.EllipseType == EllipseCenteredCircle {
		c.X = int(c.CX * float64(plane.W))
		c.Y = int(c.CY * float64(plane.H))
	} else {
		c.X = rnd.Intn(plane.W)
		c.Y = rnd.Intn(plane.H)
	}
	if maxr > 1 {
		switch c.EllipseType {
		case EllipseAny:
			c.Rx = rnd.Intn(maxr) + 1
			c.Ry = rnd.Intn(maxr) + 1
		case EllipseCenteredCircle:
			fallthrough
		case EllipseCircle:
			c.Rx = rnd.Intn(maxr) + 1
			c.Ry = c.Rx
		case EllipseFixedRadius:
			// Don't adjust the radius
		}
	} else {
		switch c.EllipseType {
		case EllipseAny:
			fallthrough
		case EllipseCenteredCircle:
			fallthrough
		case EllipseCircle:
			c.Rx, c.Ry = maxr, maxr
		case EllipseFixedRadius:
			// Don't adjust the radius
		}
	}
}

func (c *Ellipse) Draw(dc *gg.Context, scale float64) {
	dc.DrawEllipse(float64(c.X), float64(c.Y), float64(c.Rx), float64(c.Ry))
	dc.Fill()
}

func (c *Ellipse) SVG(attrs string) string {
	return fmt.Sprintf(
		"<ellipse %s cx=\"%d\" cy=\"%d\" rx=\"%d\" ry=\"%d\" />",
		attrs, c.X, c.Y, c.Rx, c.Ry)
}

func (c *Ellipse) Copy() Shape {
	a := *c
	return &a
}

func (c *Ellipse) Mutate(plane *Plane, temp float64) {
	c.mutateImpl(plane, temp, ActionAny)
}

func (c *Ellipse) mutateImpl(plane *Plane, temp float64, actions ActionType) {
	if actions == ActionNone {
		return
	}
	w := plane.W
	h := plane.H
	rnd := plane.Rnd

	maxr := w - 1
	if c.MaxRadius > 0 {
		maxr = c.MaxRadius
	}
	scale := 16 * temp

	var id int
	if c.EllipseType == EllipseFixedRadius {
		id = 0
	} else if c.EllipseType == EllipseCenteredCircle {
		id = 1
	} else {
		id = rnd.Intn(3)
	}
	for {
		switch id {
		case 0: // Move
			if (actions & ActionTranslate) == 0 {
				continue
			}
			a := int(rnd.NormFloat64() * scale)
			b := int(rnd.NormFloat64() * scale)
			c.X = clampInt(c.X+a, 0, w-1)
			c.Y = clampInt(c.Y+b, 0, h-1)
		case 1: // Mutate X
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := int(rnd.NormFloat64() * temp * float64(maxr))
			c.Rx = clampInt(c.Rx+a, 1, maxr)
			if c.EllipseType == EllipseCircle || c.EllipseType == EllipseCenteredCircle {
				c.Ry = c.Rx
			}
		case 2:
			if (actions & ActionMutate) == 0 {
				continue
			}
			a := int(rnd.NormFloat64() * temp * float64(maxr))
			c.Ry = clampInt(c.Ry+a, 1, maxr)
			if c.EllipseType == EllipseCircle || c.EllipseType == EllipseCenteredCircle {
				c.Rx = c.Ry
			}
			// TODO: Rotate?
			// TODO: Scale?
		}
		break
	}
}

func (c *Ellipse) Rasterize(rc *RasterContext) []Scanline {
	w := rc.W
	h := rc.H
	lines := rc.Lines[:0]
	aspect := float64(c.Rx) / float64(c.Ry)
	for dy := 0; dy < c.Ry; dy++ {
		y1 := c.Y - dy
		y2 := c.Y + dy
		if (y1 < 0 || y1 >= h) && (y2 < 0 || y2 >= h) {
			continue
		}
		s := int(math.Sqrt(float64(c.Ry*c.Ry-dy*dy)) * aspect)
		x1 := c.X - s
		x2 := c.X + s
		if x1 < 0 {
			x1 = 0
		}
		if x2 >= w {
			x2 = w - 1
		}
		if y1 >= 0 && y1 < h {
			lines = append(lines, Scanline{y1, x1, x2, 0xffff})
		}
		if y2 >= 0 && y2 < h && dy > 0 {
			lines = append(lines, Scanline{y2, x1, x2, 0xffff})
		}
	}
	return lines
}

type RotatedEllipse struct {
	Plane     *Plane
	X, Y      float64
	Rx, Ry    float64
	Angle     float64
	MaxRadius int
}

func NewRotatedEllipse() *RotatedEllipse {
	return &RotatedEllipse{}
}

func (c *RotatedEllipse) Init(plane *Plane) {
	rnd := plane.Rnd
	c.Plane = plane
	c.X = rnd.Float64() * float64(plane.W)
	c.Y = rnd.Float64() * float64(plane.H)
	maxr := 32.0
	if c.MaxRadius > 0 && maxr > float64(c.MaxRadius) {
		maxr = float64(c.MaxRadius) - 1
	}
	c.Rx = rnd.Float64()*maxr + 1
	c.Ry = rnd.Float64()*maxr + 1
	c.Angle = rnd.Float64() * 360
}

func (c *RotatedEllipse) Draw(dc *gg.Context, scale float64) {
	dc.Push()
	dc.RotateAbout(radians(c.Angle), c.X, c.Y)
	dc.DrawEllipse(c.X, c.Y, c.Rx, c.Ry)
	dc.Fill()
	dc.Pop()
}

func (c *RotatedEllipse) SVG(attrs string) string {
	return fmt.Sprintf(
		"<g transform=\"translate(%f %f) rotate(%f) scale(%f %f)\"><ellipse %s cx=\"0\" cy=\"0\" rx=\"1\" ry=\"1\" /></g>",
		c.X, c.Y, c.Angle, c.Rx, c.Ry, attrs)
}

func (c *RotatedEllipse) Copy() Shape {
	a := *c
	return &a
}

func (c *RotatedEllipse) Mutate(plane *Plane, temp float64) {
	c.mutateImpl(plane, temp, ActionAny)
}

func (c *RotatedEllipse) mutateImpl(plane *Plane, temp float64, actions ActionType) {
	if actions == ActionNone {
		return
	}
	w := plane.W
	h := plane.H
	rnd := plane.Rnd

	maxr := w - 1
	if c.MaxRadius > 0 {
		maxr = c.MaxRadius
	}
	scale := 16 * temp
	for {
		switch rnd.Intn(3) {
		case 0:
			c.X = clamp(c.X+rnd.NormFloat64()*scale, 0, float64(w-1))
			c.Y = clamp(c.Y+rnd.NormFloat64()*scale, 0, float64(h-1))
		case 1:
			c.Rx = clamp(c.Rx+rnd.NormFloat64()*scale, 1, float64(maxr))
			c.Ry = clamp(c.Ry+rnd.NormFloat64()*scale, 1, float64(maxr))
		case 2:
			c.Angle = c.Angle + rnd.NormFloat64()*32*temp
		}
		break
	}
}

func (c *RotatedEllipse) Rasterize(rc *RasterContext) []Scanline {
	var path raster.Path
	const n = 16
	for i := 0; i < n; i++ {
		p1 := float64(i+0) / n
		p2 := float64(i+1) / n
		a1 := p1 * 2 * math.Pi
		a2 := p2 * 2 * math.Pi
		x0 := c.Rx * math.Cos(a1)
		y0 := c.Ry * math.Sin(a1)
		x1 := c.Rx * math.Cos(a1+(a2-a1)/2)
		y1 := c.Ry * math.Sin(a1+(a2-a1)/2)
		x2 := c.Rx * math.Cos(a2)
		y2 := c.Ry * math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		x0, y0 = rotate(x0, y0, radians(c.Angle))
		cx, cy = rotate(cx, cy, radians(c.Angle))
		x2, y2 = rotate(x2, y2, radians(c.Angle))
		if i == 0 {
			path.Start(fixp(x0+c.X, y0+c.Y))
		}
		path.Add2(fixp(cx+c.X, cy+c.Y), fixp(x2+c.X, y2+c.Y))
	}
	return fillPath(rc, path)
}
