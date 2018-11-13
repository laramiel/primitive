package shape

import (
	"encoding/json"
)

type JsonShape struct {
	Type string
	Data json.RawMessage
}

func MarshalShape(input interface{}) ([]byte, error) {
	b, e := json.Marshal(input)
	if e != nil {
		return b, nil
	}
	s := JsonShape{"", json.RawMessage(b)}

	switch input.(type) {
	case *Ellipse:
		s.Type = "Ellipse"
	case *RotatedEllipse:
		s.Type = "RotatedEllipse"
	case *Line:
		s.Type = "Line"
	case *RadialLine:
		s.Type = "RadialLine"
	case *Polygon:
		s.Type = "Polygon"
	case *Quadratic:
		s.Type = "Quadratic"
	case *Rectangle:
		s.Type = "Rectangle"
	case *RotatedRectangle:
		s.Type = "RotatedRectangle"
	case *Triangle:
		s.Type = "Triangle"
	default:
		panic("Unhandled shape")
	}
	return json.Marshal(s)
}

func UnmarshalShape(b []byte) (Shape, error) {
	s := JsonShape{}
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}

	var out Shape
	switch s.Type {
	case "Ellipse":
		out = new(Ellipse)
	case "RotatedEllipse":
		out = new(RotatedEllipse)
	case "Line":
		out = new(Line)
	case "RadialLine":
		out = new(RadialLine)
	case "Polygon":
		out = new(Polygon)
	case "Quadratic":
		out = new(Quadratic)
	case "Rectangle":
		out = new(Rectangle)
	case "RotatedRectangle":
		out = new(RotatedRectangle)
	case "Triangle":
		out = new(Triangle)
	default:
		panic("Unhandled shape")
	}

	if err := json.Unmarshal(s.Data, out); err != nil {
		return nil, err
	}
	return out, nil
}
