package shape

import (
	"encoding/json"
)

type JsonShape struct {
	Ellipse          *Ellipse          `json:",omitempty"`
	RotatedEllipse   *RotatedEllipse   `json:",omitempty"`
	Line             *Line             `json:",omitempty"`
	RadialLine       *RadialLine       `json:",omitempty"`
	Polygon          *Polygon          `json:",omitempty"`
	Quadratic        *Quadratic        `json:",omitempty"`
	Cubic            *Cubic            `json:",omitempty"`
	Rectangle        *Rectangle        `json:",omitempty"`
	RotatedRectangle *RotatedRectangle `json:",omitempty"`
	Triangle         *Triangle         `json:",omitempty"`
}

func (s JsonShape) toShape() Shape {
	if s.Ellipse != nil {
		return s.Ellipse
	}
	if s.RotatedEllipse != nil {
		return s.RotatedEllipse
	}
	if s.Line != nil {
		return s.Line
	}
	if s.RadialLine != nil {
		return s.RadialLine
	}
	if s.Polygon != nil {
		return s.Polygon
	}
	if s.Quadratic != nil {
		return s.Quadratic
	}
	if s.Cubic != nil {
		return s.Cubic
	}
	if s.Rectangle != nil {
		return s.Rectangle
	}
	if s.RotatedRectangle != nil {
		return s.RotatedRectangle
	}
	if s.Triangle != nil {
		return s.Triangle
	}
	return nil
}

func makeJsonShape(input Shape) JsonShape {
	s := JsonShape{}

	switch v := input.(type) {
	case *Ellipse:
		s.Ellipse = v
	case *RotatedEllipse:
		s.RotatedEllipse = v
	case *Line:
		s.Line = v
	case *RadialLine:
		s.RadialLine = v
	case *Polygon:
		s.Polygon = v
	case *Quadratic:
		s.Quadratic = v
	case *Cubic:
		s.Cubic = v
	case *Rectangle:
		s.Rectangle = v
	case *RotatedRectangle:
		s.RotatedRectangle = v
	case *Triangle:
		s.Triangle = v
	default:
		panic("Unhandled shape")
	}

	return s
}

type SelectedShapesForJson struct {
	Shapes []JsonShape
}

func (s SelectedShapesForJson) toSelectedShapes() *SelectedShapes {
	r := &SelectedShapes{}
	for _, v := range s.Shapes {
		r.Shapes = append(r.Shapes, v.toShape())
	}
	return r
}

func makeSelectedShapesForJson(factory *SelectedShapes) *SelectedShapesForJson {
	r := &SelectedShapesForJson{}
	for _, v := range factory.Shapes {
		r.Shapes = append(r.Shapes, makeJsonShape(v))
	}
	return r
}

type JsonFactory struct {
	BasicShapes    *BasicShapes           `json:",omitempty"`
	SelectedShapes *SelectedShapesForJson `json:",omitempty"`
}

func makeJsonFactory(factory ShapeFactory) JsonFactory {
	s := JsonFactory{}
	switch v := factory.(type) {
	case *SelectedShapes:
		s.SelectedShapes = makeSelectedShapesForJson(v)
	case *BasicShapes:
		s.BasicShapes = v
	default:
		panic("Unhandled factory")
	}
	return s
}

func MarshalShapeFactory(input ShapeFactory) string {
	x := makeJsonFactory(input)
	data, err := json.Marshal(&x)
	if err != nil {
		panic("Marshal failed")
	}
	return string(data)
}

func UnmarshalShapeFactory(data string) ShapeFactory {
	mydata := []byte(data)
	x := JsonFactory{}
	if err := json.Unmarshal(mydata, &x); err != nil {
		panic("Unmarshal failed.")
	}

	if x.BasicShapes != nil {
		return x.BasicShapes
	}
	if x.SelectedShapes != nil {
		return x.SelectedShapes.toSelectedShapes()
	}
	panic("Unmarshal failed.")
}
