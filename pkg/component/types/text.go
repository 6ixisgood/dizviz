package types

import (
	"math"
	"encoding/xml"
	"image"
	"fmt"
	"log"
	"github.com/fogleman/gg"
	fontpkg "golang.org/x/image/font"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
)

type Text struct {
	c.BaseComponent

	XMLName			xml.Name		`xml:"text"`
	Font			string			`xml:"font,attr"`		
	FontStyle		string			`xml:"style,attr"`	
	FontSize		float64			`xml:"size,attr"`	
	Color			util.RGBA		`xml:"color,attr"`
	Text			string			`xml:",chardata"`

	img         	*image.RGBA
	ftCtx			*freetype.Context
}

func (t *Text) Init() {
	t.Rr = -1 // no need to rerender this once created
	t.BaseComponent.Init()

	t.Ctx = gg.NewContext(0, 0)

	// init the font and style
	var font = util.LoadFont(fmt.Sprintf("%s-%s", t.Font, t.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: t.FontSize})
	t.Ctx.SetFontFace(face)

	// get the size of the string
	w, h := t.Ctx.MeasureString(t.Text)
	w_i := int(math.Ceil(w))
	h_i := int(math.Ceil(h))
	if t.ComputedSizeX == 0 {
		t.ComputedSizeX = w_i
	}
	if t.ComputedSizeY == 0 {
		t.ComputedSizeY = h_i
	}


	// resize context
	t.Ctx = gg.NewContext(t.ComputedSizeX, t.ComputedSizeY)
	// set up a blank image
	t.img = image.NewRGBA(image.Rect(0, 0, t.ComputedSizeX, t.ComputedSizeY))

	// Set up the freetype context
	t.ftCtx = freetype.NewContext()
	t.ftCtx.SetDPI(72)
	t.ftCtx.SetFont(font)
	t.ftCtx.SetFontSize(t.FontSize)
	t.ftCtx.SetClip(t.img.Bounds())
	t.ftCtx.SetDst(t.img)
	t.ftCtx.SetSrc(image.NewUniform(t.Color.RGBA)) // set the color
	t.ftCtx.SetHinting(fontpkg.HintingNone) 
}


func (t *Text) Render() image.Image {
	// Convert the point to fixed.Point26_6 format for freetype
	pt := freetype.Pt(0, int(t.FontSize))

	// draw to the image
	_, err := t.ftCtx.DrawString(t.Text, pt)
	if err != nil {
		log.Fatal(err)
	}

	return t.img
}

func init() {
	c.RegisterComponent("text", func() c.Component { return &Text{} })
}