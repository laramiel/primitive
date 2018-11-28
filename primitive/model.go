package primitive

import (
	"fmt"
	"image"
	"math/rand"
	"strings"
	"time"
	// "sync/atomic"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
	"github.com/laramiel/primitive/primitive/shape"
)

type Model struct {
	Sw, Sh      int
	Scale       float64
	Background  Color
	Target      *image.RGBA
	Current     *image.RGBA
	Context     *gg.Context
	RC          shape.RasterContext // Rasterizes the shape into scanlines
	ColorPicker ColorPicker         // Picks the best color for the input scanlines
	Score       float64
	Shapes      []shape.Shape
	Colors      []Color
	Scores      []float64
	Workers     []*Worker
	counter     int64
}

func NewModel(target image.Image, background Color, size int, picker ColorPicker) *Model {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	aspect := float64(w) / float64(h)
	var sw, sh int
	var scale float64
	if aspect >= 1 {
		sw = size
		sh = int(float64(size) / aspect)
		scale = float64(size) / float64(w)
	} else {
		sw = int(float64(size) * aspect)
		sh = size
		scale = float64(size) / float64(h)
	}

	model := &Model{
		RC: shape.RasterContext{
			W:          w,
			H:          h,
			Lines:      make([]shape.Scanline, 0, 4096),
			Rasterizer: raster.NewRasterizer(w, h),
		},
		ColorPicker: picker,
	}

	model.Sw = sw
	model.Sh = sh
	model.Scale = scale
	model.Background = background
	model.Target = imageToRGBA(target)
	model.Current = uniformRGBA(target.Bounds(), background.NRGBA())
	model.Score = differenceFull(model.Target, model.Current)
	model.Context = model.newContext()
	vv("%+v\n", model)
	return model
}

func (model *Model) Init(numWorkers int, seed int64) {
	if seed == 0 {
		seed = time.Now().UnixNano()
		v("--seed=%d", seed)
	}
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(model.Target, rng.Int63(), model.ColorPicker)
		model.Workers = append(model.Workers, worker)
	}
}

func (model *Model) newContext() *gg.Context {
	dc := gg.NewContext(model.Sw, model.Sh)
	dc.Scale(model.Scale, model.Scale)
	dc.Translate(0.5, 0.5)
	dc.SetColor(model.Background.NRGBA())
	dc.Clear()
	return dc
}

func (model *Model) Frames(scoreDelta float64) []image.Image {
	var result []image.Image
	dc := model.newContext()
	result = append(result, imageToRGBA(dc.Image()))
	previous := 10.0
	for i, shape := range model.Shapes {
		c := model.Colors[i]
		dc.SetRGBA255(c.R, c.G, c.B, c.A)
		shape.Draw(dc, model.Scale)
		dc.Fill()
		score := model.Scores[i]
		delta := previous - score
		if delta >= scoreDelta {
			previous = score
			result = append(result, imageToRGBA(dc.Image()))
		}
	}
	return result
}

func (model *Model) SVG() string {
	bg := model.Background
	var lines []string
	// lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" width=\"%d\" height=\"%d\">", model.Sw, model.Sh))
	lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" width=\"100%%\" height=\"100%%\" preserveAspectRatio=\"none\" viewbox=\"0 0 %d %d\">", model.Sw, model.Sh))
	lines = append(lines, fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"#%02x%02x%02x\" />", model.Sw, model.Sh, bg.R, bg.G, bg.B))
	lines = append(lines, fmt.Sprintf("<g transform=\"scale(%f) translate(0.5 0.5)\" fill-opacity=\"%f\">", model.Scale, float64(model.Colors[0].A)/255))
	for i, shape := range model.Shapes {
		c := model.Colors[i]
		attrs := fmt.Sprintf("fill=\"#%02x%02x%02x\"", c.R, c.G, c.B)
		lines = append(lines, shape.SVG(attrs))
	}
	lines = append(lines, "</g>")
	lines = append(lines, "</svg>")
	return strings.Join(lines, "\n")
}

func (model *Model) Add(shape shape.Shape, alpha int) {
	before := copyRGBA(model.Current)
	lines := shape.Rasterize(&model.RC)
	color := model.ColorPicker.Select(model.Target, model.Current, lines, alpha)
	drawLines(model.Current, color, lines)
	score := differencePartial(model.Target, before, model.Current, model.Score, lines)

	model.Score = score
	model.Shapes = append(model.Shapes, shape)
	model.Colors = append(model.Colors, color)
	model.Scores = append(model.Scores, score)

	model.Context.SetRGBA255(color.R, color.G, color.B, color.A)
	shape.Draw(model.Context, model.Scale)
}

func (model *Model) Step(factory shape.ShapeFactory, alpha, repeat int) int {
	state := model.runWorkers(factory, alpha, 1000, 100, 16)
	// state = HillClimb(state, 1000).(*State)
	model.Add(state.Shape, state.Alpha)

	for i := 0; i < repeat; i++ {
		state.Worker.Init(model.Current, model.Score)
		a := state.Energy()
		state = HillClimb(state, 100).(*State)
		b := state.Energy()
		if a == b {
			break
		}
		model.Add(state.Shape, state.Alpha)
	}

	// for _, w := range model.Workers[1:] {
	// 	model.Workers[0].Heatmap.AddHeatmap(w.Heatmap)
	// }
	// SavePNG("heatmap.png", model.Workers[0].Heatmap.Image(0.5))

	counter := 0
	for _, worker := range model.Workers {
		counter += worker.Counter
	}
	return counter
}

func (model *Model) runWorkers(factory shape.ShapeFactory, a, n, age, m int) *State {
	wn := len(model.Workers)
	ch := make(chan *State, wn)
	wm := m / wn
	if m%wn != 0 {
		wm++
	}
	for i := 0; i < wn; i++ {
		worker := model.Workers[i]
		worker.Init(model.Current, model.Score)
		go model.runWorker(worker, factory, a, n, age, wm, ch)
	}
	var bestEnergy float64
	var bestState *State
	for i := 0; i < wn; i++ {
		state := <-ch
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) runWorker(worker *Worker, factory shape.ShapeFactory, a, n, age, m int, ch chan *State) {
	ch <- worker.BestHillClimbState(factory, a, n, age, m)
}
