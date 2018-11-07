package shape

import (
	"encoding/json"
)

type BasicShapes struct {
	t ShapeType
}

// NewBasicShapeFactory returns either the specific shape,
// or a randomly-selected shape.
func NewBasicShapeFactory(t int) ShapeFactory {
	return &BasicShapes{ShapeType(t)}
}

func (factory *BasicShapes) MakeShape(plane *Plane) Shape {
	t := factory.t
	if t == ShapeTypeAny {
		rnd := plane.Rnd
		t = ShapeType(rnd.Intn(9) + 1)
	}
	var s Shape
	switch t {
	case ShapeTypeTriangle:
		s = NewTriangle()
	case ShapeTypeRectangle:
		s = NewRectangle()
	case ShapeTypeEllipse:
		s = NewEllipse()
	case ShapeTypeCircle:
		s = NewCircle()
	case ShapeTypeRotatedRectangle:
		s = NewRotatedRectangle()
	case ShapeTypeLine:
		s = NewLine()
	case ShapeTypeQuadratic:
		s = NewQuadratic()
	case ShapeTypeRotatedEllipse:
		s = NewRotatedEllipse()
	case ShapeTypePolygon:
		s = NewPolygon(4, false)
	default:
		return nil
	}
	s.Init(plane)
	return s
}

func (factory *BasicShapes) Marshal() string {
	str, _ := json.Marshal(factory)
  return string(str)
}

// SelectedShapes allows the caller to add specific shapes.
// factory := NewSelectedShapeFactory()
// factory.AddShape(NewRadialLine(centerX, centerY))
// factory.AddShape(NewLine())
type SelectedShapes struct {
	Shapes []Shape
}

func NewSelectedShapeFactory() *SelectedShapes {
	return &SelectedShapes{}
}

func (factory *SelectedShapes) MakeShape(plane *Plane) Shape {
	s := factory.Shapes[plane.Rnd.Intn(len(factory.Shapes))].Copy()
	s.Init(plane)
	return s
}

func (factory *SelectedShapes) AddShape(shape Shape) {
	factory.Shapes = append(factory.Shapes, shape)
	vv("Shape: %v\n", shape)
}

func (factory *SelectedShapes) Marshal() string {
	str, _ := json.Marshal(factory)
  return string(str)
}

func UnmarshalShapeFactory(data string) ShapeFactory {
  mydata := []byte(data)
	basic := BasicShapes{ShapeType(0)}
	if err := json.Unmarshal(mydata, &basic); err == nil {
		return &basic
	}

	selected := SelectedShapes{}
	if err := json.Unmarshal(mydata, &selected); err == nil {
		return &selected
	}
	panic("Unmarshal failed.")
}
