package animations


import (
	"log"
	"image"
	"os"
	"image/color"
	"time"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"encoding/xml"
	"runtime"
	"reflect"
	"strconv"
	"net/http"

	"github.com/golang/freetype/truetype"
	"github.com/fogleman/gg"	
	"github.com/nfnt/resize"
)

var (
	imageCache = make(map[string]image.Image) // holds the images we've been using
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

	return image, time.After(time.Millisecond *1), nil
}

// Generate a new animations.Matrix type
func NewAnimation(content string) *Animation{
	// log.Printf("Creating a new animation with content %v", content)
	mtx := unmarshalContent(content)
	ctx := gg.NewContext(mtx.Sizex, mtx.Sizey)

	// clear the cache 
	imageCache = make(map[string]image.Image)

	return &Animation{
		Ctx:		ctx, 
		Mtx:		&mtx,
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

type Image struct {
	XMLName			xml.Name		`xml:"image"`
	Sizex			uint			`xml:"sizex,attr"`
	Sizey			uint			`xml:"sizey,attr"`
	Url				string			`xml:"url,attr"`
	FilePath		string			`xml:"filepath,attr"`
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
	log.Printf("Starting Unmarshalling of Matrix")
	mtx.XMLName = start.Name
	// grab any other attrs
	// grab other elements
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "sizex":
			mtx.Sizex, _ = strconv.Atoi(attr.Value)
		case "sizey":
			mtx.Sizey, _ = strconv.Atoi(attr.Value)
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
	Alignx			string			`xml:"alignx,attr"`
	Aligny			string			`xml:"aligny,attr"`
	Posx			int
	Posy			int
}

// Some custom unmarhsing to handle the XML config
func (cnt *Content) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	log.Printf("Starting Unmarshalling of Content")
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
		case "alignx":
			cnt.Alignx  = attr.Value
		case "aligny":
			cnt.Aligny = attr.Value
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
			case "image":
				i = new(Image)
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
	var font = loadFont(fmt.Sprintf("%s-%s", t.Font, t.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: t.FontSize})

	ctx.SetFontFace(face)
	// set the color for the text and draw
	ctx.SetColor(t.Color.RGBA)
	sizeX, sizeY := ctx.MeasureString(t.Text)
	ctx.DrawStringAnchored(t.Text, float64(pos.X), float64(pos.Y), 0, .35)

	return int(sizeX), int(sizeY)
 
}


func fetchImageFromURL(url string) image.Image {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching image from %v, %v", url, err)
	}
	defer res.Body.Close()


	img, _, err := image.Decode(res.Body)
	if err != nil {
		log.Printf("Error decoding image from %v, %v", url, err)
	}

	return img
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

func drawImage(ctx *gg.Context, i Image, pos image.Point) (int, int) {
	// fetch the image if it has a url
	var img image.Image
	if i.Url != "" {
		img = fetchImageFromURL(i.Url)
	} else if i.FilePath != "" {
		img = fetchImageFromPath(i.FilePath)
	}

	if i.Sizex != 0 {
		img = resizeImage(img, i.Sizex, i.Sizey)
	}

	ctx.DrawImageAnchored(img, pos.X, pos.Y, 0, .5)
	size := img.Bounds().Size()
	return size.X, size.Y
}

// Draw an animation.Content in its own gg.Context
// Return an image.Image representation of the Content
func drawContent(cnt *Content) image.Image {

	if cnt.Aligny == "top" {
		cnt.Posy = 0
	} else if cnt.Aligny == "center" {
		cnt.Posy = int(cnt.Sizey/2)
	}

	ctx := gg.NewContext(cnt.Sizex, cnt.Sizey)

	curx, cury := 0, 0
	for _, item := range cnt.Items {
		stepx, _ := 0, 0
		//y := 0
		switch i := item.(type) {
		case *Text:
			stepx, _ = drawText(ctx, *i, image.Point{cnt.Posx+curx, cnt.Posy+cury})
		case *Image:
			stepx, _ = drawImage(ctx, *i, image.Point{cnt.Posx+curx, cnt.Posy+cury})
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