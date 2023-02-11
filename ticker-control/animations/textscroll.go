package animations


import (
	"image"
	"image/color"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font"
	"github.com/fogleman/gg"	
)


type TextScrollConfig struct {
	Position  image.Point
	Size      image.Point
	BgColor   color.RGBA
	TextColor color.RGBA 
	TextFontFace  font.Face
	Dir       image.Point
	Text      string
}

type TextScrollAnimation struct {
	ctx       *gg.Context
	position  image.Point
	stroke    int
	config    TextScrollConfig
}

func NewTextScrollAnimation(config TextScrollConfig) *TextScrollAnimation {
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


func fatal(err error) {
	if err != nil {
		panic(err)
	}
}