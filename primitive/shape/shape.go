package shape

import (
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Plane struct {
	W, H int
	Rnd  *rand.Rand
}

type RasterContext struct {
	W, H       int
	Lines      []Scanline
	Rasterizer *raster.Rasterizer
}

// TODO: Shape should have an area method.

type Shape interface {
	Init(*Plane)
	Rasterize(*RasterContext) []Scanline
	Copy() Shape
	Mutate(*Plane, float64)
	Draw(dc *gg.Context, scale float64)
	SVG(attrs string) string
}

type ShapeFactory interface {
	MakeShape(*Plane) Shape
}

type ShapeType int

const (
	ShapeTypeAny ShapeType = iota
	ShapeTypeTriangle
	ShapeTypeRectangle
	ShapeTypeEllipse
	ShapeTypeCircle
	ShapeTypeRotatedRectangle
	ShapeTypeQuadratic
	ShapeTypeCubic
	ShapeTypeLine
	ShapeTypeRotatedEllipse
	ShapeTypePolygon // 10
)

type ActionType int

const (
	ActionNone      ActionType = 0
	ActionMutate    ActionType = 1
	ActionTranslate ActionType = 2
	ActionRotate    ActionType = 4
	ActionScale     ActionType = 8
	ActionAny       ActionType = 8 + 4 + 2 + 1
)
