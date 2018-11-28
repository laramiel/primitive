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
	Input       string
	Outputs     flagArray
	Background  string
	Configs     shapeConfigArray
	ColorPicker string
	Alpha       int
	InputSize   int
	OutputSize  int
	Mode        int
	Workers     int
	Nth         int
	Repeat      int
	V           bool
	VV          bool
	Seed        int64
	ZLevels     int // TODO
	Shapes      string
)

/*
Example configs.

{"SelectedShapes":{"Shapes":[
    {"Triangle":{"X1":0,"Y1":0,"X2":0,"Y2":0,"X3":0,"Y3":0,"MaxArea":60}},
    {"Ellipse":{"X":0,"Y":0,"Rx":3,"Ry":3,"EllipseType":3,"CX":0,"CY":0,"MaxRadius":0}},
    {"Polygon":{"Order":5,"Convex":true,"X":null,"Y":null}}
]}}

{"SelectedShapes":{"Shapes":[
    {"RadialLine":{"CX":0.5966666666666667,"CY":0.426,"Line":{"X1":0,"Y1":0,"X2":0,"Y2":0,"Width":0,"MaxLineWidth":0.5}}},{"RadialLine":{"CX":0.5580952380952381,"CY":0.2986666666666667,"Line":{"X1":0,"Y1":0,"X2":0,"Y2":0,"Width":0,"MaxLineWidth":0.5}}},{"Ellipse":{"X":0,"Y":0,"Rx":0,"Ry":0,"EllipseType":2,"CX":0.5966666666666667,"CY":0.426,"MaxRadius":0}},
    {"Triangle":{"X1":0,"Y1":0,"X2":0,"Y2":0,"X3":0,"Y3":0,"MaxArea":60}},
    {"Ellipse":{"X":0,"Y":0,"Rx":1,"Ry":1,"EllipseType":3,"CX":0,"CY":0,"MaxRadius":0}}
]}}

{"BasicShapes":{"T":0,"Mask":11}}

*/

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
	flag.IntVar(&Alpha, "a", 0, "alpha value")
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
	flag.StringVar(&Shapes, "shapes", "", "Shape JSON data")
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
		Configs[0].Shapes = Shapes
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
	}
	rand.Seed(Seed)
	plog.Log(1, "-seed %d\n", Seed)

	// determine worker count
	if Workers < 1 {
		Workers = runtime.NumCPU()
	}
	plog.Log(1, "-j %d\n", Workers)

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
		plog.Log(1, "Setting backgroud to average color\n")
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else if Background == "top" {
		plog.Log(1, "Setting backgroud to most frequent color\n")
		bg = primitive.MakeColor(primitive.MostFrequentImageColor(input))
	} else if Background == "center" {
		plog.Log(1, "Setting backgroud to center color\n")
		b := input.Bounds()
		bg = primitive.MakeColor(primitive.ColorAtPoint(input, (b.Max.X-b.Min.X)/2, (b.Max.Y-b.Min.Y)/2))
	} else {
		plog.Log(1, "Setting backgroud to %s\n", Background)
		bg = primitive.MakeHexColor(Background)
	}

	// run algorithm
	model := primitive.NewModel(input, bg, OutputSize, primitive.MakeColorPicker(ColorPicker))
	model.Init(Workers, Seed)
	plog.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	frame := 0
	for j, config := range Configs {
		plog.Log(1, "count=%d, mode=%d, alpha=%d, repeat=%d\n",
			config.Count, config.Mode, config.Alpha, config.Repeat)

		var factory shape.ShapeFactory = nil
		if config.Shapes != "" {
			factory = shape.UnmarshalShapeFactory(config.Shapes)
		} else {
			// "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=line 7=beziers 8=rotatedellipse 9=polygon"
			// TODO: Multiple Shapes for a BasicShapeFactory.
			factory = shape.NewBasicShapeFactory([]shape.ShapeType{shape.ShapeType(config.Mode)})
		}
		config.Shapes = shape.MarshalShapeFactory(factory)
		plog.Log(1, "%s\n", config.Shapes)

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
