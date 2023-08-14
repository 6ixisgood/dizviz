package components

import (
	"log"
	"os"
	"image"
	"math"
	"image/color"
	"fmt"
	"encoding/xml"
	"strconv"
	"runtime"
	"io/ioutil"
	"path/filepath"
	"github.com/golang/freetype/truetype"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"

)

var (
	imageCache = make(map[string]image.Image)
)

type Component interface {
	Init()
	Render()			image.Image
	Size()				image.Point
	Width()				int
	Height()			int
}

type BaseComponent struct {
	SizeX		int			`xml:"sizeX,attr"`
	SizeY		int			`xml:"sizeY,attr"`
	PosX		int			`xml:"posX,attr"`
	PosY		int			`xml:"posY,attr"`
	ctx			*gg.Context
}

func (bc *BaseComponent) Init() {
	if (bc.ctx == nil) {
		bc.ctx = gg.NewContext(bc.SizeX, bc.SizeY)
	}
}

func (bc *BaseComponent) Width() int {
	return bc.ctx.Width()
}

func (bc *BaseComponent) Height() int {
	return bc.ctx.Height()
}

func (bc *BaseComponent) Size() image.Point {
	return image.Point{bc.ctx.Width(), bc.ctx.Height()}
}


type Template struct {
	XMLName			xml.Name		`xml:"template"`
	SizeX			int				`xml:"sizeX,attr"`
	SizeY			int				`xml:"sizeY,attr"`
	Slot			int				`xml:"slot,attr"`
	Components		[]Component		`xml:",any"`

	ctx				*gg.Context

}

func (t *Template) Init() {
	ctxTmp := gg.NewContext(t.SizeX, t.SizeY)
	t.ctx = ctxTmp 

	for _, c  := range t.Components {
		c.Init()
	}
}

func (t *Template) Ready() bool {
	return t.ctx != nil
}

func (t *Template) ComponentWidth() int {
	sizeX := 0
	for _, c := range t.Components {
		sizeX += c.Width()
	}
	return sizeX
}

func (t *Template) Render() image.Image {
	if t.ctx == nil {
		t.ctx = gg.NewContext(t.SizeX, t.SizeY)
	}

	t.ctx.SetColor(color.RGBA{0,0,0,255})
	t.ctx.Clear()

	posX := 0
	var cIm image.Image
	for _, c  := range t.Components {
		cIm = c.Render()
		t.ctx.SetColor(color.RGBA{222, 255, 255, 255})
		t.ctx.DrawImage(cIm, posX, 0)
		posX += c.Width()
	}

	return t.ctx.Image()
}

func (tmpl *Template) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	log.Printf("Starting Unmarshalling of Template")
	tmpl.XMLName = start.Name

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "sizeX":
			tmpl.SizeX, _ = strconv.Atoi(attr.Value)
		case "sizeY":
			tmpl.SizeY, _ = strconv.Atoi(attr.Value)
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
			switch tt.Name.Local {
			case "text":
				i = new(Text)
			case "image":
				i = new(Image)
			case "scroller":
				i = new(Scroller)
			case "h-split":
				i = new(HorizonalSplit)
			default:
				log.Printf("Invalid component type %s", tt.Name.Local)
			}
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

type Text struct {
	BaseComponent

	XMLName			xml.Name		`xml:"text"`
	Font			string			`xml:"font,attr"`		
	FontStyle		string			`xml:"style,attr"`	
	FontSize		float64			`xml:"size,attr"`	
	Color			RGBA			`xml:"color,attr"`
	Text			string			`xml:",chardata"`
}

func (t *Text) Init() {
	t.BaseComponent.Init()

	// init the font and style
	var font = loadFont(fmt.Sprintf("%s-%s", t.Font, t.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: t.FontSize})
	t.ctx.SetFontFace(face)

	// get the size of the string
	w, h := t.ctx.MeasureString(t.Text)
	w_i := int(math.Ceil(w))
	h_i := int(math.Ceil(h))

	// set up a new ctx
	t.ctx = gg.NewContext(w_i, h_i)

	// set these again
	t.ctx.SetFontFace(face)
	t.ctx.SetColor(t.Color.RGBA)
}


func (t *Text) Render() image.Image {
	
	t.ctx.DrawStringAnchored(t.Text, 0, 0, 0, 1)

	return t.ctx.Image()
}

type Image struct {
	BaseComponent

	XMLName			xml.Name		`xml:"image"`
	Src				string			`xml:"src,attr"`

	img				image.Image
}

func (i *Image) Init() {
	// maybe fetch the image if needed here?	
	i.img = fetchImageFromPath(i.Src)
	if i.SizeX != 0 {
		i.img = resizeImage(i.img, uint(i.SizeX), uint(i.SizeY))
	}
}

func (i *Image) Width() int {
	return i.SizeX
}

func (i *Image) Height() int {
	return i.SizeY
}

func (i *Image) Render() image.Image {
	return i.img
}



// Scroller Component

type Scroller struct {
	BaseComponent

	XMLName			xml.Name		`xml:"scroller"`
	ScrollX			int				`xml:"scrollX,attr"`
	ScrollY			int				`xml:"scrollY,attr`
	Slot			Template		`xml:"template"`
}

func (s *Scroller) Init() {
	s.Slot.Init()

}

func (s *Scroller) Render() image.Image {
	if s.ctx == nil {
		log.Printf("RENDER SCROLL")
		s.ctx = gg.NewContext(s.Slot.ComponentWidth(), 50)
	}

	s.ctx.SetColor(color.RGBA{0,0,0,255})
	s.ctx.Clear()

	// render the slot
	im := s.Slot.Render()

	s.ctx.DrawImage(im, s.PosX, s.PosY)

	s.PosX = s.PosX + s.ScrollX
	s.PosY = s.PosY + s.ScrollY

	// wrap around
	if s.ScrollX < 0 {
		if s.PosX+s.ctx.Width() < 0 {
			s.PosX = 0
		}
	}

	return s.ctx.Image()	
}

type HorizonalSplit struct {
	BaseComponent

	XMLName			xml.Name		`xml:"h-split"`
	Slots			[]Template		`xml:"template"`
}

func (s *HorizonalSplit) Init() {
	for _, s := range s.Slots {
		s.Init()
	}

}

func (s *HorizonalSplit) Render() image.Image {
	width := 0
	for _, slot := range s.Slots {
		slot.Render()
		width = int(math.Max(float64(slot.ComponentWidth()), float64(width)))
	}
	height := 64
	if s.ctx == nil {
		s.ctx = gg.NewContext(width, height)
	}

	s.ctx.SetColor(color.RGBA{0,0,0,255})
	s.ctx.Clear()

	// render and draw the slots
	var im image.Image
	var y int
	for _, slot := range s.Slots {
		im = slot.Render()
		s.ctx.DrawImage(im, s.PosX, y)
		y += height/len(s.Slots)
	}

	return s.ctx.Image()	
}



// HELPER FUNCTIONS

// RGBA Struct wraps color.RGBA for unmarshalling from XML
type RGBA struct {
	color.RGBA
}

func (c *RGBA) UnmarshalXMLAttr(attr xml.Attr) error {
	var r, g, b, a uint8
	_, err := fmt.Sscanf(attr.Value, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	if err != nil {
		return err
	}
	c.RGBA = color.RGBA{r, g, b, a}
	return nil
}

func loadFont(fontName string) *truetype.Font {
	// Read font file from disk
	_, callerFile, _, _ := runtime.Caller(0)
	callerDir := filepath.Dir(callerFile)
	filePath := filepath.Join(callerDir, fmt.Sprintf("./fonts/%s.ttf", fontName))
	fontBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read font file: %v", err)
	}
// Parse font file into a truetype.Font 
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Failed to parse font file: %v", err)
	}
	return font
}


func fetchImageFromPath(path string) image.Image {
	if contents, ok := imageCache[path]; ok {
		return contents
	} 
	// else fetch the file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open the image file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode the image: %v", err)
	}	

	imageCache[path] = img


	return img
}

func resizeImage(img image.Image, sizex uint, sizey uint) image.Image {
	return resize.Resize(sizex, sizey, img, resize.Lanczos3)
}