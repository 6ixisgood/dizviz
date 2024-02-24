package types

import (
	"encoding/xml"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"image"
	"image/color"
	"math"
)

type ColorWave struct {
	c.BaseComponent

	XMLName   xml.Name `xml:"colorwave"`
	Frequency float64  `xml:"frequency,attr"`
	Amplitude float64  `xml:"amplitude,attr"`
	Speed     float64  `xml:"speed,attr"`
	Offset    float64
}

func (cw *ColorWave) Init() {
	cw.BaseComponent.Init()
	if cw.Frequency == 0 {
		cw.Frequency = 0.2 // More waves
	}
	if cw.Amplitude == 0 {
		cw.Amplitude = float64(cw.Height()) / 3 // Higher waves
	}
	if cw.Speed == 0 {
		cw.Speed = 2 // Faster speed
	}
}

func (cw *ColorWave) Render() image.Image {
	cw.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	cw.Ctx.Clear()

	for x := 0.0; x < float64(cw.Width()); x++ {
		y := float64(cw.Height())/2 + cw.Amplitude*math.Sin(cw.Frequency*x+cw.Offset)

		// Use a gradient of colors
		r := uint8((math.Sin(cw.Frequency*x) + 1) / 2 * 255)
		g := uint8((math.Cos(cw.Frequency*x) + 1) / 2 * 255)
		b := uint8((math.Sin(cw.Offset*2) + 1) / 2 * 255)

		// Draw circles for a cooler effect
		cw.Ctx.SetColor(color.RGBA{r, g, b, 255})
		cw.Ctx.DrawCircle(x, y, 3) // Larger circles
		cw.Ctx.Fill()
	}

	cw.Offset += cw.Speed * 0.05 // Adjust speed for a smoother animation

	return cw.Ctx.Image()
}

func init() {
	c.RegisterComponent("colorwave", func() c.Component { return &ColorWave{} })
}
