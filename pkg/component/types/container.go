package types

import (
	"encoding/xml"
	"image/color"
	"github.com/fogleman/gg"
	"image"
	compCommon "github.com/6ixisgood/matrix-ticker/pkg/component/common"
)

type Container struct {
	compCommon.BaseComponent

	XMLName			xml.Name		`xml:"container"`
	Slot			*compCommon.Template		`xml:"template"`
}

func (c *Container) Init() {
	c.Rr = 100 
	c.BaseComponent.Init()
	c.Slot.SetParentSize(c.ComputedSizeX, c.ComputedSizeY)
	c.Slot.Init()

}

func (c *Container) Render() image.Image {
	if c.Ctx == nil {
		c.Ctx = gg.NewContext(c.ComputedSizeX, c.ComputedSizeY)
	}

	c.Ctx.SetColor(color.RGBA{0,0,0,255})
	c.Ctx.Clear()

	// render and draw the slots
	im := c.Slot.Render()
	c.Ctx.DrawImage(im, 0, 0)

	return c.Ctx.Image()	
}

func init() {
	compCommon.RegisterComponent("container", func() compCommon.Component { return &Container{} })
}