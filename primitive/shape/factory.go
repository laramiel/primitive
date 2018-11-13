package shape

import (
	"encoding/json"
)

type BasicShapes struct {
	t    ShapeType
	mask uint32
}

const biggest = int(ShapeTypePolygon)
const allShapes = uint32(1<<uint(biggest)) - 1

// NewBasicShapeFactory returns either the specific shape,
// or a randomly-selected shape.
func NewBasicShapeFactory(t []ShapeType) ShapeFactory {
	var mask uint32 = 0
	if len(t) == 0 {
		return &BasicShapes{ShapeTypeAny, allShapes}
	}
	if len(t) == 1 && t[0] != ShapeTypeAny {
		return &BasicShapes{t[0], 0}
	}
	for _, v := range t {
		if v == ShapeTypeAny {
			return &BasicShapes{ShapeTypeAny, allShapes}
		}
		mask |= 1 << (uint32(v) - 1)
	}
	return &BasicShapes{ShapeTypeAny, mask}
}

func (factory *BasicShapes) MakeShape(plane *Plane) Shape {

	t := factory.t
	for t == ShapeTypeAny {
		v := plane.Rnd.Intn(biggest - 1)
		if factory.mask&(1<<uint32(v)) != 0 {
			t = ShapeType(v + 1)
		}
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
		panic("Aah!")
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

func (factory *SelectedShapes) MarshalJSON() (b []byte, e error) {
	var data []json.RawMessage
	for _, v := range factory.Shapes {
		b, e := MarshalShape(v)
		if e == nil {
			data = append(data, json.RawMessage(b))
		}
	}
	return json.Marshal(data)
}

func (factory *SelectedShapes) Marshal() string {
	str, _ := json.Marshal(factory)
	return string(str)
}

func UnmarshalShapeFactory(data string) ShapeFactory {
	mydata := []byte(data)
	basic := BasicShapes{}
	if err := json.Unmarshal(mydata, &basic); err == nil {
		return &basic
	}

	selected := SelectedShapes{}
	if err := json.Unmarshal(mydata, &selected); err == nil {
		return &selected
	}
	panic("Unmarshal failed.")
}
