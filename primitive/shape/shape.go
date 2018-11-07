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

type Shape interface {
	Init(*Plane)
	Rasterize(*RasterContext) []Scanline
	Copy() Shape
	Mutate(*Plane)
	Draw(dc *gg.Context, scale float64)
	SVG(attrs string) string
}

type ShapeFactory interface {
	MakeShape(*Plane) Shape

	// Marshal marshals the factory into a JSON string which can be
	// unmarshalled by UnmarshalShapeFactory()
	Marshal() string
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
	ShapeTypePolygon // 9
)
