package agent

import (
	"github.com/6ixisgood/matrix-ticker/pkg/display"
)

// Capabilities describes what an agent can do
type Capabilities struct {
	Displays       []DisplayCapability // Multiple displays per agent
	SupportedViews []string            // List of view types this agent can render
}

// DisplayCapability describes a single display device
type DisplayCapability struct {
	DisplayID   string // Unique ID for this display
	DisplayType string // "led-matrix", "hdmi", "monitor", "file"
	Width       int    // Display width in pixels
	Height      int    // Display height in pixels
	MaxFPS      int    // Maximum frames per second
	ColorDepth  int    // Bits per pixel
}

// NewCapabilitiesFromDisplay creates capabilities from a single display
func NewCapabilitiesFromDisplay(disp display.Display, displayID string, displayType string, supportedViews []string) Capabilities {
	caps := disp.GetCapabilities()

	displayCap := DisplayCapability{
		DisplayID:   displayID,
		DisplayType: displayType,
		Width:       caps.Width,
		Height:      caps.Height,
		MaxFPS:      caps.MaxFrameRate,
		ColorDepth:  caps.ColorDepth,
	}

	return Capabilities{
		Displays:       []DisplayCapability{displayCap},
		SupportedViews: supportedViews,
	}
}

// NewCapabilitiesFromDisplays creates capabilities from multiple displays
func NewCapabilitiesFromDisplays(displays map[string]display.Display, displayTypes map[string]string, supportedViews []string) Capabilities {
	displayCaps := make([]DisplayCapability, 0, len(displays))

	for displayID, disp := range displays {
		caps := disp.GetCapabilities()
		displayType := displayTypes[displayID]
		if displayType == "" {
			displayType = "unknown"
		}

		displayCaps = append(displayCaps, DisplayCapability{
			DisplayID:   displayID,
			DisplayType: displayType,
			Width:       caps.Width,
			Height:      caps.Height,
			MaxFPS:      caps.MaxFrameRate,
			ColorDepth:  caps.ColorDepth,
		})
	}

	return Capabilities{
		Displays:       displayCaps,
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
