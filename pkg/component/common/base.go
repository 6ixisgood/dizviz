package common

import (
	"github.com/fogleman/gg"
	"image"
	"strconv"
	"strings"
	"time"
)

type BaseComponent struct {
	SizeX         string `xml:"sizeX,attr"`
	SizeY         string `xml:"sizeY,attr"`
	ComputedSizeX int
	ComputedSizeY int
	ParentWidth   int
	ParentHeight  int
	PosX          int `xml:"posX,attr"`
	PosY          int `xml:"posY,attr"`

	Ctx     *gg.Context
	prevImg image.Image
	Ticker  *time.Ticker
	Rr      int // render rate in milliseconds
}

func (bc *BaseComponent) Init() {
	// determine sizing
	if strings.HasSuffix(bc.SizeX, "%") {
		percentage, _ := strconv.Atoi(bc.SizeX[:len(bc.SizeX)-1])
		bc.ComputedSizeX = int(bc.ParentWidth * percentage / 100)
	} else {
		bc.ComputedSizeX, _ = strconv.Atoi(bc.SizeX)
	}

	if strings.HasSuffix(bc.SizeY, "%") {
		percentage, _ := strconv.Atoi(bc.SizeY[:len(bc.SizeY)-1])
		bc.ComputedSizeY = int(bc.ParentHeight * percentage / 100)
	} else {
		bc.ComputedSizeY, _ = strconv.Atoi(bc.SizeY)
	}

	// create a context, if needed
	if bc.Ctx == nil && bc.ComputedSizeX > 0 && bc.ComputedSizeY > 0 {
		bc.Ctx = gg.NewContext(bc.ComputedSizeX, bc.ComputedSizeY)
	}
	// set a render rate, if needed
	if bc.Rr == 0 {
		bc.Rr = 5
	}

	// set ticker for render rate
	if bc.Rr > 0 {
		bc.Ticker = time.NewTicker(time.Duration(bc.Rr) * time.Millisecond)
	}
}

func (bc *BaseComponent) Width() int {
	return bc.ComputedSizeX
}

func (bc *BaseComponent) Height() int {
	return bc.ComputedSizeY
}

func (bc *BaseComponent) Size() image.Point {
	return image.Point{bc.Ctx.Width(), bc.Ctx.Height()}
}

func (bc *BaseComponent) SetParentSize(width int, height int) {
	bc.ParentWidth = width
	bc.ParentHeight = height
}

func (bc *BaseComponent) TickerChan() <-chan time.Time {
	if bc.Ticker != nil {
		return bc.Ticker.C
	} else {
		return nil
	}
}

func (bc *BaseComponent) Stop() {
	if bc.Ticker != nil {
		bc.Ticker.Stop()
	}
}

func (bc *BaseComponent) PrevImg() image.Image {
	return bc.prevImg
}

func (bc *BaseComponent) SetPrevImg(img image.Image) {
	bc.prevImg = img
}
