package primitive

type BasicShapes struct {
	t ShapeType
}

func (factory *BasicShapes) MakeShape(worker *Worker) Shape {

	t := factory.t
	if (t == ShapeTypeAny) {
		rnd := worker.Rnd
		t = ShapeType(rnd.Intn(8)+1)
	}
	switch t {
	case ShapeTypeTriangle:
		return NewRandomTriangle(worker)
	case ShapeTypeRectangle:
		return NewRandomRectangle(worker)
	case ShapeTypeEllipse:
		return NewRandomEllipse(worker)
	case ShapeTypeCircle:
		return NewRandomCircle(worker)
	case ShapeTypeRotatedRectangle:
		return NewRandomRotatedRectangle(worker)
	case ShapeTypeLine:
		return NewRandomLine(worker)
	case ShapeTypeQuadratic:
		return NewRandomQuadratic(worker)
	case ShapeTypeRotatedEllipse:
		return NewRandomRotatedEllipse(worker)
	case ShapeTypePolygon:
		return NewRandomPolygon(worker, 4, false)
	}
	return nil
}

func NewBasicShapeFactory(t int) ShapeFactory {
	return &BasicShapes{ShapeType(t)};
}

/*

type FactoryImpl struct {
	Shapes []Shape
}

func (factory *FactoryImpl) MakeShape(worker *Worker) Shape {
	s := factory.Shapes[worker.Rnd.Intn(len(factory.Shapes))].Copy()
	s.Init(worker)
	return s
}

func (factory *FactoryImpl) AddShape(s Shape) {
	factory.Shapes = append(factory.Shapes, s)
}

func (factory *FactoryImpl) AddShapeType(t ShapeType, rnd *rand.Rand) {
	if t == ShapeTypeAny {
		u := ShapeTypePolygon
		for u != t {
			factory.AddShape(factory.newShape(u, rnd))
			u = ShapeType(int(u) - 1)
		}
	} else {
		factory.AddShape(factory.newShape(t, rnd))
	}
}

func (factory *FactoryImpl) NewRadialLine(worker *Worker) Shape {
	// Radial line example
	const centerX = (1051 + 202) / 2100.0
	const centerY = (437 + 202) / 1500.0
	return NewRadialLine(worker, centerX, centerY)
}

func (factory *FactoryImpl) NewFixedCircle(worker *Worker) Shape {
	// Random fixed circle.
	return NewRandomFixedCircle(worker, 4)
}

func (factory *FactoryImpl) newShape(t ShapeType, rnd *rand.Rand) Shape {
	switch t {
	case ShapeTypeTriangle:
		return NewRandomTriangle(rnd)
	case ShapeTypeRectangle:
		return NewRandomRectangle(rnd)
	case ShapeTypeEllipse:
		return NewRandomEllipse(rnd)
	case ShapeTypeCircle:
		return NewRandomCircle(rnd)
	case ShapeTypeRotatedRectangle:
		return NewRandomRotatedRectangle(rnd)
	case ShapeTypeLine:
		return NewRandomLine(rnd)
	case ShapeTypeQuadratic:
		return NewRandomQuadratic(rnd)
	case ShapeTypeRotatedEllipse:
		return NewRandomRotatedEllipse(rnd)
	case ShapeTypePolygon:
		return NewRandomPolygon(rnd, 4, false)
	}
}
*/