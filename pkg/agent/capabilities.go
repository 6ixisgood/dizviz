package agent

import (
	"github.com/6ixisgood/matrix-ticker/pkg/display"
)

// Capabilities describes what an agent can do
type Capabilities struct {
	DisplayType    string   // "led-matrix", "hdmi", "monitor", "file"
	Width          int      // Display width in pixels
	Height         int      // Display height in pixels
	MaxFPS         int      // Maximum frames per second
	ColorDepth     int      // Bits per pixel
	SupportedViews []string // List of view types this agent can render
}

// NewCapabilitiesFromDisplay creates capabilities from a display
func NewCapabilitiesFromDisplay(disp display.Display, displayType string, supportedViews []string) Capabilities {
	caps := disp.GetCapabilities()

	return Capabilities{
		DisplayType:    displayType,
		Width:          caps.Width,
		Height:         caps.Height,
		MaxFPS:         caps.MaxFrameRate,
		ColorDepth:     caps.ColorDepth,
		SupportedViews: supportedViews,
	}
}

// DefaultSupportedViews returns the default list of supported view types
func DefaultSupportedViews() []string {
	return []string{
		"text",
		"image",
		"nfl-scoreboard",
		"sleeper-matchups",
		"weather",
		"clock",
	}
}
