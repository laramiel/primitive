package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/laramiel/primitive/primitive"
	plog "github.com/laramiel/primitive/primitive/log"
	"github.com/laramiel/primitive/primitive/shape"
	"github.com/nfnt/resize"
)

var (
	Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray
  ColorPicker string
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	Workers    int
	Nth        int
	Repeat     int
	V          bool
	VV         bool
	Seed       int64
	ZLevels    int // TODO
)

type flagArray []string

func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
	Shapes string
}

type shapeConfigArray []shapeConfig

func (i *shapeConfigArray) String() string {
	return ""
}

func (i *shapeConfigArray) Set(value string) error {
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, shapeConfig{int(n), Mode, Alpha, Repeat, ""})
	return nil
}

func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.Var(&Configs, "n", "number of primitives")
	flag.StringVar(&Background, "bg", "", "background color (hex)")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	flag.IntVar(&ZLevels, "z", 1, "Maximum z-index")
	flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.IntVar(&Mode, "m", 1, "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=line 7=beziers 8=rotatedellipse 9=polygon")
	flag.IntVar(&Workers, "j", 0, "number of parallel workers (default uses all cores)")
	flag.IntVar(&Nth, "nth", 1, "save every Nth frame (put \"%d\" in path)")
	flag.IntVar(&Repeat, "rep", 0, "add N extra shapes per iteration with reduced search")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
	flag.Int64Var(&Seed, "seed", 0, "RNG seed")
	flag.StringVar(&ColorPicker, "color", "", "Color picker to use")
}

func errorMessage(message string) bool {
	fmt.Fprintln(os.Stderr, message)
	return false
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// parse and validate arguments
	flag.Parse()
	ok := true
	if Input == "" {
		ok = errorMessage("ERROR: input argument required")
	}
	if len(Outputs) == 0 {
		ok = errorMessage("ERROR: output argument required")
	}
	if len(Configs) == 0 {
		ok = errorMessage("ERROR: number argument required")
	}
	if len(Configs) == 1 {
		Configs[0].Mode = Mode
		Configs[0].Alpha = Alpha
		Configs[0].Repeat = Repeat
		Configs[0].Shapes = ""
	}
	for _, config := range Configs {
		if config.Count < 1 {
			ok = errorMessage("ERROR: number argument must be > 0")
		}
	}
	if !ok {
		fmt.Println("Usage: primitive [OPTIONS] -i input -o output -n count")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// set log level
	if V {
		plog.LogLevel = 1
	}
	if VV {
		plog.LogLevel = 2
	}

	// seed random number generator
	if Seed == 0 {
		Seed = time.Now().UTC().UnixNano()
	} else {
		// Seed from command line, use only one worker by default
		if Workers < 1 {
			Workers = 1
		}
	}
	rand.Seed(Seed)

	// determine worker count
	if Workers < 1 {
		Workers = runtime.NumCPU()
	}

	// read input image
	plog.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input)
	check(err)

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}

	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else if Background == "top" {
		bg = primitive.MakeColor(primitive.MostFrequentImageColor(input))
	} else {
		bg = primitive.MakeHexColor(Background)
	}
	bg = primitive.MakeColor(primitive.ColorAtPoint(input, 5, 5))
	plog.Log(1, "%v\n", bg)

	// run algorithm
	model := primitive.NewModel(input, bg, OutputSize)

	// Change the color model.
	model.ColorPicker = primitive.MakeColorPicker(ColorPicker)

	model.Init(Workers, Seed)
	plog.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	frame := 0
	for j, config := range Configs {
		plog.Log(1, "count=%d, mode=%d, alpha=%d, repeat=%d\n",
			config.Count, config.Mode, config.Alpha, config.Repeat)

		{
			factory := shape.NewSelectedShapeFactory()

			// Radial line example
			const cX = (1051 + 202) / 2100.0
			const cY = (437 + 202) / 1500.0
			factory.AddShape(shape.NewRadialLine(cX, cY))

			// Radial line example
			const cX2 = 1172 / 2100.0
			const cY2 = 448 / 1500.0
			factory.AddShape(shape.NewRadialLine(cX2, cY2))

			// Centered circle example
			factory.AddShape(shape.NewCenteredCircle(cX, cY))
			factory.AddShape(shape.NewMaxAreaTriangle(60))
			factory.AddShape(shape.NewFixedCircle(1))

			plog.Log(1, "%s\n", factory.Marshal())
		}

		var factory shape.ShapeFactory
		if config.Shapes == "" {
			factory = shape.NewBasicShapeFactory(config.Mode)
			config.Shapes = factory.Marshal()
		} else {
			factory = shape.UnmarshalShapeFactory(config.Shapes)
		}

		for i := 0; i < config.Count; i++ {
			frame++

			// find optimal shape and add it to the model
			t := time.Now()
			n := model.Step(factory, config.Alpha, config.Repeat)
			nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
			elapsed := time.Since(start).Seconds()
			plog.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)

			// write output image(s)
			for _, output := range Outputs {
				ext := strings.ToLower(filepath.Ext(output))
				percent := strings.Contains(output, "%")
				saveFrames := percent && ext != ".gif"
				saveFrames = saveFrames && frame%Nth == 0
				last := j == len(Configs)-1 && i == config.Count-1
				if saveFrames || last {
					path := output
					if percent {
						path = fmt.Sprintf(output, frame)
					}
					plog.Log(1, "writing %s\n", path)
					switch ext {
					default:
						check(fmt.Errorf("unrecognized file extension: %s", ext))
					case ".png":
						check(primitive.SavePNG(path, model.Context.Image()))
					case ".jpg", ".jpeg":
						check(primitive.SaveJPG(path, model.Context.Image(), 95))
					case ".svg":
						check(primitive.SaveFile(path, model.SVG()))
					case ".gif":
						frames := model.Frames(0.001)
						check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
					}
				}
			}
		}
	}
}

/*

 */
