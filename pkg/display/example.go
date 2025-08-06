// +build !noled

package display

import (
	"github.com/sixisgoood/go-rpi-rgb-led-matrix"
)

// Example usage of the new simplified display system
func ExampleUsage() {
	// Create a matrix display
	matrixDisplay := NewMatrixDisplay()

	// Create a display manager (uses default settings: 30 FPS, buffer size 10)
	manager := NewManager(matrixDisplay)
	
	// Or create with custom settings
	// manager := NewManagerWithSettings(matrixDisplay, 60, 20) // 60 FPS, buffer size 20

	// Create matrix config
	config := &rgbmatrix.HardwareConfig{
		Rows:       32,
		Cols:       64,
		ChainLength:   1,
		Parallel:   1,
		PWMBits:    11,
		PWMLSBNanoseconds: 130,
		Brightness: 100,
	}

	// Create your view (this would be implemented elsewhere)
	// view := views.NewImageView(...)

	// Start the display pipeline
	// if err := manager.Start(view, config); err != nil {
	//     log.Fatal(err)
	// }

	// To change views later:
	// newView := views.NewTextView(...)
	// manager.ChangeView(newView)

	// When done:
	// manager.Stop()
	
	_ = manager // avoid unused variable error
	_ = config  // avoid unused variable error
}
