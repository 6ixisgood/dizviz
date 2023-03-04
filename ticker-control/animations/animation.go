package animations


import (
	"log"
	"image"
	"image/color"
	"time"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"encoding/xml"
	"runtime"
	"reflect"
	"strconv"

	"github.com/golang/freetype/truetype"
	"github.com/fogleman/gg"	
)

////////////////////////////////
// Structs
////////////////////////////////

// Animation Struct
// Wraps the Matrix and root gg.Context
type Animation struct {
	Ctx			*gg.Context
	Mtx			*Matrix

}

// Returns the next image in the animation and the delay between the next call to this function
func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	a.Ctx.SetColor(color.RGBA{0,0,0,255})
	a.Ctx.Clear()
	image := drawMatrix(a.Ctx, a.Mtx)
	return image, time.After(time.Millisecond *50), nil
}

// Generate a new animations.Matrix type
func NewAnimation(content string) *Animation{
	log.Printf("Creating a new animation with content %v", content)
	cnt := unmarshalContent(content)
	ctx := gg.NewContext(64, 32)
	return &Animation{
		Ctx:		ctx, 
		Mtx:		&cnt,
	}
}


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

type Text struct {
	XMLName			xml.Name		`xml:"text"`
	Font			string			`xml:"font,attr"`
	FontStyle		string			`xml:"fontstyle,attr"`
	FontSize		float64			`xml:"fontsize,attr"`
	// FontPosition	string			`xml:"fontposition,attr`
	Color			RGBA			`xml:"color,attr"`
	Text			string			`xml:",chardata"`
}


// Represents the overall Matrix you are drawing to. 
// A matrix can contain many Contents
type Matrix struct {
	XMLName		xml.Name		`xml:"matrix"`
	Items		[]interface{}	`xml:"content"`
	Sizex		int				`xml:"sizex,attr"`
	Sizey		int				`xml:"sizey,attr"`
}

// Custom unmashaling for Matrix type in XML
func (mtx *Matrix) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	log.Printf("Starting Unmarshalling")
	mtx.XMLName = start.Name
	// grab any other attrs

	// decode inner elements
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		var i interface{}
		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "content":
				i = new(Content)
			default:
				log.Printf("Nothing")
			}
			if i != nil {
				err = d.DecodeElement(i, &tt)
				if err != nil {
					return err
				}
				mtx.Items = append(mtx.Items, i)
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

// This struct holds the actual content that you are putting on the Matrix
// e.g. Text, Images, etc.
type Content struct {
	XMLName			xml.Name		`xml:"content"`
	Items			[]interface{}	`xml:",any"`
	Sizex			int				`xml:"sizex,attr"`
	Sizey			int				`xml:"sizey,attr"`	
	MtxPosx			int				`xml:"posx,attr"`
	MtxPosy			int				`xml:"posy,attr"`	
	Color			RGBA			`xml:"color,attr"`
	Bgcolor			RGBA			`xml:"bgcolor,attr"`
	Scroll			bool			`xml:"scroll,attr"`
	Scrollx			int				`xml:"scrollx,attr"`
	Scrolly			int				`xml:"scrolly,attr"`
	Posx			int
	Posy			int
}

// Some custom unmarhsing to handle the XML config
func (cnt *Content) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	log.Printf("Starting Unmarshalling")
	cnt.XMLName = start.Name
	// decode the start element
	// if err := d.DecodeElement(cnt, &start); err != nil {
	// 	return err
	// }

	// grab other elements
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "sizex":
			cnt.Sizex, _ = strconv.Atoi(attr.Value)
		case "sizey":
			cnt.Sizey, _ = strconv.Atoi(attr.Value)
		case "posx":
			cnt.MtxPosx, _ = strconv.Atoi(attr.Value)
		case "posy":
			cnt.MtxPosy, _ = strconv.Atoi(attr.Value)
		case "scroll":
			cnt.Scroll = attr.Value == "true" 
		case "scrollx":
			cnt.Scrollx, _ = strconv.Atoi(attr.Value)
		case "scrolly":
			cnt.Scrolly, _ = strconv.Atoi(attr.Value)
		}
	}

	// decode inner elements
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		var i interface{}
		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "text":
				i = new(Text)
			default:
				log.Printf("Nothing")
			}
			if i != nil {
				err = d.DecodeElement(i, &tt)
				if err != nil {
					return err
				}
				cnt.Items = append(cnt.Items, i)
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


// Take in a string representation of a Matrix XML and return an animation.Matrix
func unmarshalContent(content string) Matrix {
	var cnt Matrix
	err := xml.Unmarshal([]byte(content), &cnt)
	if err != nil {
		log.Fatalf("Unable to unmarshal xml content: '%v'", err)
	}
	return cnt

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

// Handle some errors for us
func fatal(err error) {
	if err != nil {
		panic(err)
	}
}


// Draw text to a screen
// Pass it the gg.Context to write to, the text represented as an animations.Text and the starting position
func drawText(ctx *gg.Context, t Text, pos image.Point) (int, int) {
	// ctx.DrawStringAnchored(t.Text, 10, 10, .5, .5)
	var font = loadFont(fmt.Sprintf("%s-%s", t.Font, t.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: t.FontSize})

	ctx.SetFontFace(face)
	// set the color for the text and draw
	ctx.SetColor(t.Color.RGBA)
	sizeX, sizeY := ctx.MeasureString(t.Text)
	ctx.DrawStringAnchored(t.Text, float64(pos.X), float64(pos.Y), 0, .9)

	return int(sizeX), int(sizeY)

}


func drawImage(ctx *gg.Context, i image.Image) {
	// do nothing
}

// Draw an animation.Content in its own gg.Context
// Return an image.Image representation of the Content
func drawContent(cnt *Content) image.Image {
	// log.Printf("Pos: %v, %v, %v, %v", cnt.Color, cnt.MtxPosx, cnt.MtxPosy, cnt)

	ctx := gg.NewContext(cnt.Sizex, cnt.Sizey)

	curx, cury := 0, 0
	for _, item := range cnt.Items {
		stepx, _ := 0, 0
		//y := 0
		switch i := item.(type) {
		case *Text:
			stepx, _ = drawText(ctx, *i, image.Point{cnt.Posx+curx, cnt.Posy+cury})
		case *image.Image:
			drawImage(ctx, *i)
		default:
			log.Printf("Type '%v'", reflect.TypeOf(i))
		}
		curx += stepx
		//cury += stepy
	}

	// scroll the content
	cnt.Posx += cnt.Scrollx
	cnt.Posy += cnt.Scrolly

	// wrap around
	if cnt.Scrollx < 0 {
		if cnt.Posx+curx < 0 {
			cnt.Posx = cnt.Sizex
		}
	}


	return ctx.Image()

}

// Creates the final image.Image to draw. Will delegate sub-tasks to draw Content
func drawMatrix(ctx *gg.Context, mtx *Matrix) image.Image {
	for _, item := range mtx.Items {
		switch i := item.(type) {
		case *Content:
			img := drawContent(i)  
			ctx.DrawImage(img, i.MtxPosx, i.MtxPosy)
		}

	}
	return ctx.Image()
}