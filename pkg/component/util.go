package component

import (
	"encoding/xml"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	imageCache = make(map[string]image.Image)
)

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
