package types

import (
	"math"
	"encoding/xml"
	"image/color"
	"image"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
)

type ColorWave struct {
	c.BaseComponent
	XMLName   xml.Name      `xml:"colorwave"`
	Colors    []color.RGBA  // Array of colors for the waves
	Phase     float64       // Current phase for the wave function
	WaveSpeed float64       // Speed at which the wave progresses
}

func (cw *ColorWave) Init() {
	cw.BaseComponent.Init()
	cw.Colors = []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	cw.Phase = 0
	cw.WaveSpeed = 0.05
}

func (cw *ColorWave) Render() image.Image {
	cw.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	cw.Ctx.Clear()
	for _, col := range cw.Colors {
		for x := 0; x < cw.Width(); x++ {
			y := math.Sin(float64(x)/float64(cw.Width())*2*math.Pi+cw.Phase) * float64(cw.Height()/4) + float64(cw.Height()/2)
			cw.Ctx.SetColor(col)
			cw.Ctx.DrawPoint(float64(x), y, 2)
		}
		cw.Phase += cw.WaveSpeed
	}
	cw.Ctx.Fill()
	return cw.Ctx.Image()
}

func init() {
	c.RegisterComponent("colorwave", func() c.Component { return &ColorWave{} })
}