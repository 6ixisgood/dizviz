package types

import (
	"bytes"
	"encoding/xml"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
	"image/gif"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type Image struct {
	c.BaseComponent

	XMLName      xml.Name `xml:"image"`
	Src          string   `xml:"src,attr"`
	Loop         bool     `xml:"loop,attr"`
	frames       []image.Image
	currentFrame int
}

func (i *Image) Init() {
	i.Rr = -1
	i.BaseComponent.Init()

	data, filePath, err := util.FetchImage(i.Src, i.ComputedSizeX, i.ComputedSizeY)
	if err != nil {
		log.Fatal(err)
	}

	extension := strings.ToLower(filepath.Ext(filePath))

	// Decode based on extension
	if extension == ".gif" {
		gifData, err := gif.DecodeAll(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}

		for _, frame := range gifData.Image {
			i.frames = append(i.frames, frame)
		}

		// If the GIF should loop, reset the ticker accordingly
		if i.Loop {
			i.Ticker = time.NewTicker(time.Duration(gifData.Delay[0]) * 10 * time.Millisecond)
		}
	} else if extension == ".png" || extension == ".jpg" || extension == ".jpeg" {
		// do we have it saved in this size yet?
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}
		i.frames = append(i.frames, img)
	} else if extension == ".svg" {
		// Handle SVG files
		svg, err := oksvg.ReadIconStream(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}

		// Set the desired size for the SVG (computed size)
		width := 32  //int(i.ComputedSizeX)
		height := 32 //int(i.ComputedSizeY)

		// Create an RGBA canvas to render the SVG
		rgba := image.NewRGBA(image.Rect(0, 0, width, height))

		// Set the target size for the SVG
		svg.SetTarget(0, 0, float64(width), float64(height))

		// Create a new Dasher for drawing the SVG onto the canvas
		scanner := rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())
		dasher := rasterx.NewDasher(width, height, scanner)

		// Draw the SVG onto the RGBA canvas
		svg.Draw(dasher, 1)

		i.frames = append(i.frames, rgba)
	} else {
		log.Fatal("Unsupported file extension")
	}
}

func (i *Image) Render() image.Image {
	// Return the current frame and update the frame index
	img := i.frames[i.currentFrame]
	i.currentFrame = (i.currentFrame + 1) % len(i.frames)
	return img
}

func init() {
	c.RegisterComponent("image", func() c.Component { return &Image{} })
}
