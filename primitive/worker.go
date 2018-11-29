package primitive

import (
	"image"
	"math/rand"

	"github.com/golang/freetype/raster"
	"github.com/laramiel/primitive/primitive/shape"
)

type Worker struct {
	Plane       shape.Plane
	RC          shape.RasterContext
	Target      *image.RGBA
	Current     *image.RGBA
	Buffer      *image.RGBA
	Heatmap     *Heatmap
	Score       float64
	Counter     int
	ColorPicker ColorPicker // Picks the best color for the input scanlines
}

func NewWorker(target *image.RGBA, seed int64, picker ColorPicker) *Worker {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	worker := Worker{
		Plane: shape.Plane{
			W:   w,
			H:   h,
			Rnd: rand.New(rand.NewSource(seed)),
		},
		RC: shape.RasterContext{
			W:          w,
			H:          h,
			Lines:      make([]shape.Scanline, 0, 4096),
			Rasterizer: raster.NewRasterizer(w, h),
		},
	}
	worker.Target = target
	worker.Buffer = image.NewRGBA(target.Bounds())
	worker.Heatmap = NewHeatmap(w, h)
	worker.ColorPicker = picker
	vv("%+v\n", worker)
	return &worker
}

func (worker *Worker) Init(current *image.RGBA, score float64) {
	worker.Current = current
	worker.Score = score
	worker.Counter = 0
	worker.Heatmap.Clear()
}

func (worker *Worker) Energy(shape shape.Shape, alpha int) float64 {
	worker.Counter++
	lines := shape.Rasterize(&worker.RC)
	// worker.Heatmap.Add(lines)
	color := worker.ColorPicker.Select(worker.Target, worker.Current, lines, alpha)
	copyLines(worker.Buffer, worker.Current, lines)
	drawLines(worker.Buffer, color, lines)
	energy := differencePartial(worker.Target, worker.Current, worker.Buffer, worker.Score, lines)
	return energy
}

func (worker *Worker) BestHillClimbState(factory shape.ShapeFactory, a, n, age, m int) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < m; i++ {
		state := worker.BestRandomState(factory, a, n)
		before := state.Energy()
		state = HillClimb(state, age).(*State)
		energy := state.Energy()
		vv("%dx random: %.6f -> %dx hill climb: %.6f\n", n, before, age, energy)
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (worker *Worker) BestRandomState(factory shape.ShapeFactory, a, n int) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < n; i++ {
		state := NewState(worker, factory.MakeShape(&worker.Plane), a)
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

