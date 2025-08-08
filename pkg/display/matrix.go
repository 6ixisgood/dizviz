//go:build !noled
// +build !noled

package display

import (
	"fmt"
	"image"
	"image/draw"
	"log"

	rgbmatrix "github.com/sixisgoood/go-rpi-rgb-led-matrix"
)

// MatrixDisplay implements Display for RGB LED matrix hardware
type MatrixDisplay struct {
	matrix   rgbmatrix.Matrix
	canvas   *rgbmatrix.Canvas
	config   *rgbmatrix.HardwareConfig
	stopChan chan struct{}
	running  bool
}

// NewMatrixDisplay creates a new matrix display
func NewMatrixDisplay() *MatrixDisplay {
	return &MatrixDisplay{
		stopChan: make(chan struct{}),
	}
}

// Start initializes the matrix hardware with the given config
func (d *MatrixDisplay) Start(config interface{}) error {
	hardwareConfig, ok := config.(*rgbmatrix.HardwareConfig)
	if !ok {
		return fmt.Errorf("invalid config type for MatrixDisplay, expected *rgbmatrix.HardwareConfig")
	}

	// Setup matrix hardware
	matrix, err := rgbmatrix.NewRGBLedMatrix(hardwareConfig)
	if err != nil {
		return fmt.Errorf("failed to create matrix: %w", err)
	}

	d.matrix = matrix
	d.canvas = rgbmatrix.NewCanvas(matrix)
	d.config = hardwareConfig
	d.running = true

	log.Printf("Matrix display started successfully")
	return nil
}

// Stop shuts down the matrix display
func (d *MatrixDisplay) Stop() {
	if !d.running {
		return
	}

	close(d.stopChan)
	d.running = false

	if d.canvas != nil {
		d.canvas.Close()
	}

	if d.matrix != nil {
		d.matrix.Close()
	}

	log.Printf("Matrix display stopped")
}

// Render displays a single frame on the matrix
func (d *MatrixDisplay) Render(frame image.Image) error {
	if !d.running || d.matrix == nil || d.canvas == nil {
		return fmt.Errorf("matrix display not started")
	}

	if frame == nil {
		return fmt.Errorf("nil frame provided")
	}

	// Use draw.Draw to efficiently copy the frame to the canvas
	draw.Draw(d.canvas, d.canvas.Bounds(), frame, frame.Bounds().Min, draw.Src)

	// Render to the physical display
	return d.canvas.Render()
}

// GetCapabilities returns the display's capabilities
func (d *MatrixDisplay) GetCapabilities() DisplayCapabilities {
	width, height := 64, 32 // default matrix size
	if d.matrix != nil {
		width, height = d.matrix.Geometry()
	}

	return DisplayCapabilities{
		MaxFrameRate: 60,
		Width:        width,
		Height:       height,
		ColorDepth:   24, // RGB
	}
}
