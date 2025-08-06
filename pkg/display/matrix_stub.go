// +build noled

package display

import (
	"fmt"
	"image"
	"log"
)

// MatrixDisplay implements Display for testing without LED hardware
type MatrixDisplay struct {
	running bool
	width   int
	height  int
}

// NewMatrixDisplay creates a new matrix display stub
func NewMatrixDisplay() *MatrixDisplay {
	return &MatrixDisplay{
		width:  64,
		height: 32,
	}
}

// Start initializes the matrix display stub
func (d *MatrixDisplay) Start(config interface{}) error {
	log.Printf("Matrix display stub started")
	d.running = true
	return nil
}

// Stop shuts down the matrix display stub
func (d *MatrixDisplay) Stop() {
	if !d.running {
		return
	}
	d.running = false
	log.Printf("Matrix display stub stopped")
}

// Render displays a single frame (stub implementation)
func (d *MatrixDisplay) Render(frame image.Image) error {
	if !d.running {
		return fmt.Errorf("matrix display not started")
	}

	if frame == nil {
		return fmt.Errorf("nil frame provided")
	}

	// Stub implementation - just log that we would render
	bounds := frame.Bounds()
	log.Printf("Rendering frame %dx%d", bounds.Dx(), bounds.Dy())
	return nil
}

// GetCapabilities returns the display's capabilities
func (d *MatrixDisplay) GetCapabilities() DisplayCapabilities {
	return DisplayCapabilities{
		MaxFrameRate: 60,
		Width:        d.width,
		Height:       d.height,
		ColorDepth:   24, // RGB
	}
}
