package primitive

import "github.com/fogleman/gg"

type Shape interface {
	Init(*Worker)
	Rasterize() []Scanline
	Copy() Shape
	Mutate()
	Draw(dc *gg.Context, scale float64)
	SVG(attrs string) string
}

type ShapeFactory interface {
	MakeShape(*Worker) Shape	
}

type ShapeType int

const (
	ShapeTypeAny ShapeType = iota
	ShapeTypeTriangle
	ShapeTypeRectangle
	ShapeTypeEllipse
	ShapeTypeCircle
	ShapeTypeRotatedRectangle
	ShapeTypeLine
	ShapeTypeQuadratic
	ShapeTypeRotatedEllipse
	ShapeTypePolygon
)
