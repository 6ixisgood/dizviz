package types

import (
    "log"
    "strings"
    "path/filepath"
    "bytes"
    "time"
    "image"
    "image/gif"
    "encoding/xml"
    c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
    "github.com/6ixisgood/matrix-ticker/pkg/util"
)

type Image struct {
	c.BaseComponent

	XMLName			xml.Name		`xml:"image"`
	Src				string			`xml:"src,attr"`
	Loop			bool			`xml:"loop,attr"`
	frames			[]image.Image
	currentFrame	int
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
    } else if extension == ".png" || extension == ".jpg" || extension == ".jpeg"  {
    	// do we have it saved in this size yet?
        img, _, err := image.Decode(bytes.NewReader(data))
        if err != nil {
            log.Fatal(err)
        }
        i.frames = append(i.frames, img)
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