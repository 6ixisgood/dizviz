package types

import (
	"encoding/xml"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"github.com/fogleman/gg"
	"image"
	"image/color"
)

type Scroller struct {
	c.BaseComponent

	XMLName xml.Name    `xml:"scroller"`
	ScrollX int         `xml:"scrollX,attr"`
	ScrollY int         `xml:"scrollY,attr"`
	Slot    *c.Template `xml:"template"`
}

func (s *Scroller) Init() {
	s.Rr = 400 // render this one a bit faster than most
	s.BaseComponent.Init()
	s.Slot.SetParentSize(s.ComputedSizeX, s.ComputedSizeY)
	s.Slot.Init()
}

func (s *Scroller) Render() image.Image {
	// render the slot
	im := s.Slot.Render()

	if s.Ctx == nil {
		s.Ctx = gg.NewContext(s.ComputedSizeX, s.ComputedSizeY)
	}

	s.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	s.Ctx.Clear()

	s.Ctx.DrawImage(im, s.PosX, s.PosY)

	s.PosX = s.PosX + s.ScrollX
	s.PosY = s.PosY + s.ScrollY

	// wrap around
	if s.ScrollX < 0 {
		if s.PosX+s.Ctx.Width() < 0 {
			s.PosX = 0
		}
	}

	return s.Ctx.Image()
}

func init() {
	c.RegisterComponent("scroller", func() c.Component { return &Scroller{} })
}
