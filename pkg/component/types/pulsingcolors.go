package types

import (
	"encoding/xml"
	"image"
	"image/color"
	"math"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
)

type PulsingCircles struct {
	c.BaseComponent

	XMLName      xml.Name `xml:"pulsing-circles"`
	NumCircles   int      `xml:"numCircles,attr"`
	MaxRadius    float64  `xml:"maxRadius,attr"`
	Offset       float64
}

func (pc *PulsingCircles) Init() {
	pc.BaseComponent.Init()
	if pc.NumCircles == 0 {
		pc.NumCircles = 5 // Default number of circles
	}
	if pc.MaxRadius == 0 {
		pc.MaxRadius = 20 // Default maximum radius
	}
}

func (pc *PulsingCircles) Render() image.Image {
	pc.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	pc.Ctx.Clear()

	centerX, centerY := float64(pc.Width())/2, float64(pc.Height())/2
	for i := 0; i < pc.NumCircles; i++ {
		angle := 2 * math.Pi / float64(pc.NumCircles) * float64(i) + pc.Offset
		radius := pc.MaxRadius * (0.5 + 0.5*math.Sin(pc.Offset+float64(i)))

		r := uint8((math.Sin(angle) + 1) / 2 * 255)
		g := uint8((math.Cos(angle) + 1) / 2 * 255)
		b := uint8(255 - r)

		pc.Ctx.SetColor(color.RGBA{r, g, b, 255})
		pc.Ctx.DrawCircle(centerX+radius*math.Cos(angle), centerY+radius*math.Sin(angle), radius)
		pc.Ctx.Fill()
	}

	pc.Offset += 0.1 // Adjust for a pulsing effect
	return pc.Ctx.Image()
}

func init() {
	c.RegisterComponent("pulsing-circles", func() c.Component { return &PulsingCircles{} })
}
