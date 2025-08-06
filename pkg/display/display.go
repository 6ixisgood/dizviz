package display

import (
	"image"
)

// Display interface for hardware-specific display implementations
type Display interface {
	Start(config interface{}) error
	Stop()
	Render(frame image.Image) error
	GetCapabilities() DisplayCapabilities
}

// DisplayCapabilities describes what a display can do
type DisplayCapabilities struct {
	MaxFrameRate int
	Width        int
	Height       int
	ColorDepth   int
}
