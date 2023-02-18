package animations


import (
	"log"
	"image"
	"image/color"
	"time"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font"
	"github.com/fogleman/gg"	
)


type TextScrollConfig struct {
	Position  		image.Point
	Size      		image.Point
	BgColor   		color.RGBA
	TextColor 		color.RGBA 
	TextFontFace  	font.Face
	TextFontSize	float64
	TextFontName	string
	Dir       		image.Point
	Text      		string
}

type TextScrollAnimation struct {
	ctx       *gg.Context
	position  image.Point
	stroke    int
	config    TextScrollConfig
}

func NewTextScrollAnimation(config TextScrollConfig) *TextScrollAnimation {
	var font = loadFont(config.TextFontName)
	var face = truetype.NewFace(font, &truetype.Options{Size: config.TextFontSize})
	config.TextFontFace = face

	return &TextScrollAnimation{
		ctx:    gg.NewContext(config.Size.X, config.Size.Y),
		stroke: 5,
		config: config,
	}
}

func QuickAnimation() *TextScrollAnimation {
	var font, _ = truetype.Parse(goregular.TTF)
	var face = truetype.NewFace(font, &truetype.Options{Size: 10})

	return &TextScrollAnimation{
		ctx:    gg.NewContext(64, 32),
		stroke: 5,
		config: TextScrollConfig{
			Position: 		image.Point{0, 0},
			Size:			image.Point{64, 32},
			BgColor:		color.RGBA{0, 0, 0, 255},	
			TextColor:		color.RGBA{255, 255, 255, 255},	
			TextFontFace:	face,
			Dir:			image.Point{1, 0},
			Text:			"Quick Animation",
		},
	} 
}

func (a *TextScrollAnimation) Next() (image.Image, <-chan time.Time, error) {
	// set initial postion
	var sizeX float64
	// set the text font
	a.ctx.SetFontFace(a.config.TextFontFace)
	// set the background color and clear
	a.ctx.SetColor(a.config.BgColor)
	a.ctx.Clear()

	// set the color for the text and draw
	a.ctx.SetColor(a.config.TextColor)
	a.ctx.DrawStringAnchored(a.config.Text, float64(a.position.X), float64(a.position.Y), 0, .9)
	// get the size of the text
	sizeX, _ = a.ctx.MeasureString(a.config.Text)
	// update the animation to scroll the text
	defer a.updatePosition(sizeX)
	return a.ctx.Image(), time.After(time.Millisecond * 50), nil
}

func (a *TextScrollAnimation) updatePosition(size float64) {
	a.position.X -= 1 * a.config.Dir.X

	if a.position.X+a.stroke+int(size) < 0 {
		a.position.X = a.ctx.Width()
	}
}


func loadFont(fontName string) *truetype.Font {
	// Read font file from disk
	_, callerFile, _, _ := runtime.Caller(0)
	callerDir := filepath.Dir(callerFile)
	filePath := filepath.Join(callerDir, fmt.Sprintf("./fonts/%s.ttf", fontName))
	fontBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read font file: %v", err)
	}

	// Parse font file into a truetype.Font
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Failed to parse font file: %v", err)
	}
	return font
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}