package types

import (
	"encoding/xml"
	"fmt"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"math"
)

// Rainbow text
type AnimatedRainbowText struct {
	c.BaseComponent

	XMLName    xml.Name `xml:"rainbow-text"`
	Font       string   `xml:"font,attr"`
	FontStyle  string   `xml:"style,attr"`
	FontSize   float64  `xml:"size,attr"`
	Text       string   `xml:",chardata"`
	colorIndex int
}

func (art *AnimatedRainbowText) Init() {
	art.BaseComponent.Init()
	art.Ctx = gg.NewContext(0, 0)

	// get the size of the string
	w, h := art.Ctx.MeasureString(art.Text)
	w_i := int(math.Ceil(w))
	h_i := int(math.Ceil(h))
	if art.ComputedSizeX == 0 {
		art.ComputedSizeX = w_i
	}
	if art.ComputedSizeY == 0 {
		art.ComputedSizeY = h_i
	}

	// resize context
	art.Ctx = gg.NewContext(art.ComputedSizeX, art.ComputedSizeY)

	var font = util.LoadFont(fmt.Sprintf("%s-%s", art.Font, art.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: art.FontSize})
	art.Ctx.SetFontFace(face)
}

func (art *AnimatedRainbowText) Render() image.Image {
	art.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	art.Ctx.Clear()
	rainbowColors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{255, 127, 0, 255}, // Orange
		{255, 255, 0, 255}, // Yellow
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{75, 0, 130, 255},  // Indigo
		{148, 0, 211, 255}, // Violet
	}

	// Get the size of the text
	_, h := art.Ctx.MeasureString(art.Text)
	startX := 0.0

	for _, char := range art.Text {
		currentColor := rainbowColors[art.colorIndex]
		art.Ctx.SetColor(currentColor)
		charStr := string(char)
		art.Ctx.DrawString(charStr, startX, h) // Draw each character

		// Update starting x-coordinate for next character
		cw, _ := art.Ctx.MeasureString(charStr)
		startX += cw

		// Update color index for the next character
		art.colorIndex = (art.colorIndex + 1) % len(rainbowColors)
	}

	return art.Ctx.Image()
}

func init() {
	c.RegisterComponent("rainbow-text", func() c.Component { return &AnimatedRainbowText{} })
}
