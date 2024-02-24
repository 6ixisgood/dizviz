package common

import (
	"image"
	"time"
)

type Component interface {
	Init()                               // Ran before componenet is rendered
	Render() image.Image                 // Render the component to an image.Image representation
	Width() int                          // Return the width of the component. Used to help position components on the display
	Height() int                         // Return the height of the compoent. Used to help position components on the display
	PrevImg() image.Image                // return the previously rendered image
	SetPrevImg(img image.Image)          // the the previously rendered iamge
	TickerChan() <-chan time.Time        // return the channel for the ticker
	Stop()                               // Stop the ticker
	SetParentSize(width int, height int) // set the parent's size to use in calculations
}

type ComponentContext struct {
	CacheDir string
}

var (
	RegisteredComponents = map[string]func() Component{}
)

func RegisterComponent(name string, comp func() Component) {
	RegisteredComponents[name] = comp
}
