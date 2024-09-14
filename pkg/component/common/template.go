package common

import (
	"encoding/xml"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"math"
)

type Template struct {
	BaseComponent

	XMLName xml.Name `xml:"template"`

	Align      string      `xml:"align,attr"`
	Justify    string      `xml:"justify,attr"`
	Direction  string      `xml:"dir,attr"`
	BgColor    string      `xml:"bg-color,attr"`
	Components []Component `xml:",any"`
}

func (t *Template) Init() {
	t.Rr = -1
	t.BaseComponent.Init()
	for _, c := range t.Components {
		c.SetParentSize(t.ComputedSizeX, t.ComputedSizeY) // Set parent size on each child component
		c.Init()
	}

	if t.BgColor == "" {
		t.BgColor = "#000000FF"
	}

	// create context with sizes
	ctxTmp := gg.NewContext(t.ComputedSizeX, t.ComputedSizeY)
	t.Ctx = ctxTmp
}

func (t *Template) Ready() bool {
	return t.Ctx != nil
}

func (t *Template) computeSpace(availableSpace int, itemCount int, mode string) int {
	if mode == "space-between" && itemCount > 1 {
		return availableSpace / (itemCount - 1)
	} else if mode == "space-around" && itemCount > 0 {
		return availableSpace / itemCount
	}
	return 0
}

func (t *Template) computePositionAndSpace(axis Axis, imListLen int, alignment string) (int, int) {
	position := 0
	space := 0
	switch alignment {
	case "center":
		position = (axis.TemplateSize - axis.Length) / 2
	case "end":
		position = axis.TemplateSize - axis.Length
	case "space-between":
		space = t.computeSpace(axis.TemplateSize-axis.Length, imListLen, "space-between")
	case "space-around":
		position += t.computeSpace(axis.TemplateSize-axis.Length, imListLen, "space-around") / 2
		space = t.computeSpace(axis.TemplateSize-axis.Length, imListLen, "space-around")
	}

	return position, space
}

type Axis struct {
	Length       int
	Max          int
	TemplateSize int
	Position     int
	Space        int
}

func (t *Template) Render() image.Image {
	var r, g, b, a uint8
	fmt.Sscanf(t.BgColor, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	t.Ctx.SetColor(color.RGBA{r, g, b, a})
	t.Ctx.Clear()

	var componentLengthX, componentLengthY int
	var componentMaxX, componentMaxY int
	var cIm image.Image
	var imList []image.Image
	for _, c := range t.Components {
		// chan to check if we should re-render or just grab last image
		select {
		case <-c.TickerChan():
			// Ticker has ticked
			cIm = c.Render()
		default:
			// Ticker has not ticked
			cIm = c.PrevImg()
			// check for nil
			if cIm == nil {
				cIm = c.Render()
			}
		}

		// save the prev image for next time
		c.SetPrevImg(cIm)
		//  save the renderings to list and adjust width
		imList = append(imList, cIm)
		componentLengthX += c.Width()
		componentLengthY += c.Height()
		componentMaxX = int(math.Max(float64(componentMaxX), float64(c.Width())))
		componentMaxY = int(math.Max(float64(componentMaxY), float64(c.Height())))
	}

	var primary, secondary Axis
	if t.Direction == "col" {
		primary = Axis{TemplateSize: t.ComputedSizeY, Max: componentMaxY, Length: componentLengthY}
		secondary = Axis{TemplateSize: t.ComputedSizeX, Max: componentMaxX, Length: componentLengthX}
	} else {
		primary = Axis{TemplateSize: t.ComputedSizeX, Max: componentMaxX, Length: componentLengthX}
		secondary = Axis{TemplateSize: t.ComputedSizeY, Max: componentMaxY, Length: componentLengthY}
	}

	// Modularized positioning logic
	primary.Position, primary.Space = t.computePositionAndSpace(primary, len(imList), t.Justify)
	secondary.Position, _ = t.computePositionAndSpace(secondary, len(imList), t.Align)
	for _, im := range imList {
		bounds := im.Bounds()
		if t.Direction == "col" {
			secondary.Length = bounds.Dx()
			secondary.Position, _ = t.computePositionAndSpace(secondary, len(imList), t.Align)
			t.Ctx.DrawImage(im, secondary.Position, primary.Position)
			primary.Position += bounds.Dy()
		} else {
			secondary.Length = bounds.Dy()
			secondary.Position, _ = t.computePositionAndSpace(secondary, len(imList), t.Align)
			t.Ctx.DrawImage(im, primary.Position, secondary.Position)
			primary.Position += bounds.Dx()
		}

		primary.Position += primary.Space // Increment only if it's space-between or space-around.
	}

	return t.Ctx.Image()
}

func (t *Template) Stop() {
	for _, c := range t.Components {
		c.Stop()
	}
	t.BaseComponent.Stop()
}

func (tmpl *Template) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tmpl.XMLName = start.Name

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "size-x":
			tmpl.SizeX = attr.Value
		case "size-y":
			tmpl.SizeY = attr.Value
		case "justify":
			tmpl.Justify = attr.Value
		case "align":
			tmpl.Align = attr.Value
		case "dir":
			tmpl.Direction = attr.Value
		case "bg-color":
			tmpl.BgColor = attr.Value
		}
	}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		var i Component
		switch tt := t.(type) {
		case xml.StartElement:
			i = RegisteredComponents[tt.Name.Local]()
			// log.Printf("Invalid component type %s", tt.Name.Local)
			if i != nil {
				err = d.DecodeElement(i, &tt)
				if err != nil {
					return err
				}
				tmpl.Components = append(tmpl.Components, i)
				i = nil
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}

	}
	return nil
}

func init() {
	RegisterComponent("template", func() Component { return &Template{} })
}
