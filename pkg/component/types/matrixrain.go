package types

import (
	"encoding/xml"
	"image"
	"math/rand"
	"image/color"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
)

type Drop struct {
	x, y, speed float64
}

type MatrixRain struct {
	c.BaseComponent
	XMLName xml.Name `xml:"matrix-rain"`
	Drops   []Drop
	NumDrops int
}

func (mr *MatrixRain) Init() {
	mr.BaseComponent.Init()
	mr.NumDrops = 100
	for i := 0; i < mr.NumDrops; i++ {
		mr.Drops = append(mr.Drops, Drop{rand.Float64() * float64(mr.Width()), rand.Float64() * float64(mr.Height()), rand.Float64() * 5 + 1})
	}
}

func (mr *MatrixRain) Render() image.Image {
	mr.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	mr.Ctx.Clear()
	for _, drop := range mr.Drops {
		char := rune(33 + rand.Intn(94)) // Select a random ASCII character
		drop.y += drop.speed
		if drop.y > float64(mr.Height()) {
			drop.y = 0
		}
		mr.Ctx.SetColor(color.RGBA{0, 255, 0, 255})
		mr.Ctx.DrawStringAnchored(string(char), drop.x, drop.y, 0.5, 0.5)
	}
	return mr.Ctx.Image()
}

func init() {
	c.RegisterComponent("matrix-rain", func() c.Component { return &MatrixRain{} })
}