package types

import (
	"encoding/xml"
	"image"
	"image/color"
	"math"

	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
)

type ColorGrid struct {
	c.BaseComponent

	XMLName  xml.Name `xml:"color-grid"`
	GridSize int      `xml:"gridSize,attr"`
	Offset   float64
}

func (cg *ColorGrid) Init() {
	cg.BaseComponent.Init()
	if cg.GridSize == 0 {
		cg.GridSize = 10 // Default grid size
	}
}

func (cg *ColorGrid) Render() image.Image {
	cg.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	cg.Ctx.Clear()

	gridSize := float64(cg.GridSize)
	for x := 0.0; x < float64(cg.Width()); x += gridSize {
		for y := 0.0; y < float64(cg.Height()); y += gridSize {
			r := uint8((math.Sin(cg.Offset+x/10) + 1) / 2 * 255)
			g := uint8((math.Cos(cg.Offset+y/10) + 1) / 2 * 255)
			b := uint8(255 - r)

			cg.Ctx.SetColor(color.RGBA{r, g, b, 255})
			cg.Ctx.DrawRectangle(x, y, gridSize, gridSize)
			cg.Ctx.Fill()
		}
	}

	cg.Offset += 0.05 // Adjust for a smooth color transition
	return cg.Ctx.Image()
}

func init() {
	c.RegisterComponent("color-grid", func() c.Component { return &ColorGrid{} })
}
