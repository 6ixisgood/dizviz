package common

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

type Template struct {
	BaseComponent

	XMLName xml.Name `xml:"template"`

	Align           string `xml:"align,attr"`
	Justify         string `xml:"justify,attr"`
	Direction       string `xml:"dir,attr"`
	BgColor         string `xml:"bg-color,attr"`
	Overflow        string `xml:"overflow,attr"`
	OverflowX       string `xml:"overflow-x,attr"`
	OverflowY       string `xml:"overflow-y,attr"`
	ScrollSpeed     int    `xml:"scroll-speed,attr"`
	BounceDelay     int    `xml:"bounce-delay,attr"`
	EasingMode      string `xml:"easing,attr"`
	CycleDurationMs int    `xml:"cycle-duration-ms,attr"`
	DwellPercent    int    `xml:"dwell-percent,attr"`
	PosX            int
	PosY            int
	// Internal scroll state
	bounceBackX       bool        // for horizontal bounce mode
	bounceBackY       bool        // for vertical bounce mode
	bounceDelayCountX int         // delay counter for horizontal bounce
	bounceDelayCountY int         // delay counter for vertical bounce
	easingStartTimeX  time.Time   // easing start time for X axis
	easingStartTimeY  time.Time   // easing start time for Y axis
	Components        []Component `xml:",any"`
}

func (t *Template) Init() {
	t.Rr = 30
	t.BaseComponent.Init()
	for _, c := range t.Components {
		c.SetParentSize(t.ComputedSizeX, t.ComputedSizeY) // Set parent size on each child component
		c.Init()
	}

	if t.BgColor == "" {
		t.BgColor = "#000000FF"
	}

	// Handle overflow attributes - separate X/Y take precedence over combined
	if t.OverflowX == "" && t.OverflowY == "" {
		// Use combined overflow for both axes if separate ones aren't specified
		if t.Overflow == "" {
			t.Overflow = "hidden"
		}
		t.OverflowX = t.Overflow
		t.OverflowY = t.Overflow
	} else {
		// Use separate overflow attributes, defaulting to hidden if not specified
		if t.OverflowX == "" {
			t.OverflowX = "hidden"
		}
		if t.OverflowY == "" {
			t.OverflowY = "hidden"
		}
	}

	// Initialize scroll settings
	if t.ScrollSpeed == 0 {
		t.ScrollSpeed = 1 // Default scroll speed
	}
	if t.BounceDelay == 0 {
		t.BounceDelay = 10 // Default bounce delay (frames to pause at each end)
	}
	if t.EasingMode == "" {
		t.EasingMode = "ease-in-out-sine" // Default easing mode
	}
	if t.DwellPercent == 0 {
		t.DwellPercent = 20 // Default 20% dwell time at each end
	}
	// CycleDurationMs = 0 means auto-calculate based on content and scroll speed

	// Initialize bounce state for separate axes
	t.bounceBackX = false
	t.bounceBackY = false
	t.bounceDelayCountX = 0
	t.bounceDelayCountY = 0

	// Initialize easing start times
	now := time.Now()
	t.easingStartTimeX = now
	t.easingStartTimeY = now

	// create context with sizes
	ctxTmp := gg.NewContext(t.ComputedSizeX, t.ComputedSizeY)
	t.Ctx = ctxTmp
}

func (t *Template) Ready() bool {
	return t.Ctx != nil
}

func (t *Template) computeSpace(availableSpace int, itemCount int, mode string) int {
	if mode == "space-between" && itemCount > 1 {
		return availableSpace / (itemCount - 1)
	} else if mode == "space-around" && itemCount > 0 {
		return availableSpace / itemCount
	} else if mode == "space-evenly" && itemCount > 0 {
		return availableSpace / (itemCount + 1)
	}
	return 0
}

func (t *Template) computePositionAndSpace(axis Axis, imListLen int, alignment string) (int, int) {
	position := 0
	space := 0
	switch alignment {
	case "center":
		position = (axis.TemplateSize - axis.Length) / 2
	case "end":
		position = axis.TemplateSize - axis.Length
	case "space-between":
		space = t.computeSpace(axis.TemplateSize-axis.Length, imListLen, "space-between")
	case "space-around":
		position += t.computeSpace(axis.TemplateSize-axis.Length, imListLen, "space-around") / 2
		space = t.computeSpace(axis.TemplateSize-axis.Length, imListLen, "space-around")
	case "space-evenly":
		space = t.computeSpace(axis.TemplateSize-axis.Length, imListLen, "space-evenly")
		position = space
	}

	return position, space
}

type Axis struct {
	Length       int
	Max          int
	TemplateSize int
	Position     int
	Space        int
}

func (t *Template) Render() image.Image {
	var r, g, b, a uint8
	fmt.Sscanf(t.BgColor, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	t.Ctx.SetColor(color.RGBA{r, g, b, a})
	t.Ctx.Clear()

	var componentLengthPrimary int
	var componentMaxSecondary int
	var cIm image.Image
	var imList []image.Image
	for _, c := range t.Components {

		// if display = inline then
		// end
		// elif display = block then
		// end

		// chan to check if we should re-render or just grab last image
		select {
		case <-c.TickerChan():
			// Ticker has ticked
			cIm = c.Render()
		default:
			// Ticker has not ticked
			cIm = c.PrevImg()
			// check for nil
			if cIm == nil {
				cIm = c.Render()
			}
		}

		// save the prev image for next time
		c.SetPrevImg(cIm)
		//  save the renderings to list and adjust dimensions
		imList = append(imList, cIm)

		// Calculate dimensions based on direction
		bounds := cIm.Bounds()
		if t.Direction == "col" {
			componentLengthPrimary += bounds.Dy()
			componentMaxSecondary = int(math.Max(float64(componentMaxSecondary), float64(bounds.Dx())))
		} else {
			componentLengthPrimary += bounds.Dx()
			componentMaxSecondary = int(math.Max(float64(componentMaxSecondary), float64(bounds.Dy())))
		}
	}

	var primary, secondary Axis
	if t.Direction == "col" {
		primary = Axis{TemplateSize: t.ComputedSizeY, Max: componentMaxSecondary, Length: componentLengthPrimary}
		secondary = Axis{TemplateSize: t.ComputedSizeX, Max: componentMaxSecondary, Length: 0} // Length not used for secondary in flex
	} else {
		primary = Axis{TemplateSize: t.ComputedSizeX, Max: componentMaxSecondary, Length: componentLengthPrimary}
		secondary = Axis{TemplateSize: t.ComputedSizeY, Max: componentMaxSecondary, Length: 0} // Length not used for secondary in flex
	}

	// Modularized positioning logic
	primary.Position, primary.Space = t.computePositionAndSpace(primary, len(imList), t.Justify)
	secondary.Position, _ = t.computePositionAndSpace(secondary, len(imList), t.Align)

	// Handle auto-expansion before calculating scroll offsets
	if t.OverflowX == "auto" || t.OverflowY == "auto" {
		t.handleAutoExpansion(componentLengthPrimary, componentMaxSecondary)
	}

	// Apply scroll offset to the starting position based on overflow mode
	scrollOffsetX, scrollOffsetY := t.calculateScrollOffset(componentLengthPrimary, componentMaxSecondary)

	for _, im := range imList {
		bounds := im.Bounds()
		var x, y int

		if t.Direction == "col" {
			// For column direction: secondary axis position can vary per item for individual alignment
			itemSecondaryPos := secondary.Position
			if t.Align == "start" || t.Align == "center" || t.Align == "end" {
				// For these alignments, each item can have individual cross-axis positioning
				itemSecondary := Axis{
					TemplateSize: t.ComputedSizeX,
					Length:       bounds.Dx(),
				}
				itemSecondaryPos, _ = t.computePositionAndSpace(itemSecondary, 1, t.Align)
			}
			x = itemSecondaryPos + scrollOffsetX
			y = primary.Position + scrollOffsetY
			primary.Position += bounds.Dy()
		} else {
			// For row direction
			itemSecondaryPos := secondary.Position
			if t.Align == "start" || t.Align == "center" || t.Align == "end" {
				itemSecondary := Axis{
					TemplateSize: t.ComputedSizeY,
					Length:       bounds.Dy(),
				}
				itemSecondaryPos, _ = t.computePositionAndSpace(itemSecondary, 1, t.Align)
			}
			x = primary.Position + scrollOffsetX
			y = itemSecondaryPos + scrollOffsetY
			primary.Position += bounds.Dx()
		}

		t.Ctx.DrawImage(im, x, y)
		primary.Position += primary.Space // Add spacing after each item except the last
	}

	im := t.Ctx.Image()

	// Apply clipping to keep content within template bounds
	// Scrolling happens during component positioning, not here
	// Handle both axes separately now
	needsClipping := false
	if t.OverflowX == "hidden" || t.OverflowY == "hidden" ||
		(t.OverflowX != "visible" && t.OverflowX != "auto" && t.OverflowY != "visible" && t.OverflowY != "auto") {
		needsClipping = true
	}

	if needsClipping {
		bounds := image.Rect(0, 0, t.ComputedSizeX, t.ComputedSizeY)
		clippedIm := image.NewRGBA(bounds)
		draw.Draw(clippedIm, bounds, im, image.Point{}, draw.Src)
		im = clippedIm
	}

	return im
}

// handleAutoExpansion expands the template size when overflow="auto" is set
func (t *Template) handleAutoExpansion(contentLength, contentMaxSecondary int) {
	// Handle horizontal auto-expansion
	if t.OverflowX == "auto" {
		var requiredWidth int
		if t.Direction == "row" || t.Direction == "" {
			// In row direction, contentLength is the total width needed
			requiredWidth = contentLength
		} else {
			// In column direction, contentMaxSecondary is the width of widest item
			requiredWidth = contentMaxSecondary
		}

		if requiredWidth > t.ComputedSizeX {
			t.ComputedSizeX = requiredWidth
			// Recreate context with new width
			t.Ctx = gg.NewContext(t.ComputedSizeX, t.ComputedSizeY)
		}
	}

	// Handle vertical auto-expansion
	if t.OverflowY == "auto" {
		var requiredHeight int
		if t.Direction == "col" {
			// In column direction, contentLength is the total height needed
			requiredHeight = contentLength
		} else {
			// In row direction, contentMaxSecondary is the height of tallest item
			requiredHeight = contentMaxSecondary
		}

		if requiredHeight > t.ComputedSizeY {
			t.ComputedSizeY = requiredHeight
			// Recreate context with new height
			t.Ctx = gg.NewContext(t.ComputedSizeX, t.ComputedSizeY)
		}
	}
}

// calculateScrollOffset calculates the scroll offset based on overflow mode and content size
func (t *Template) calculateScrollOffset(contentLength, contentMaxSecondary int) (int, int) {
	var offsetX, offsetY int

	// Calculate horizontal overflow offset
	offsetX = t.calculateHorizontalOffset(contentLength, contentMaxSecondary)

	// Calculate vertical overflow offset
	offsetY = t.calculateVerticalOffset(contentLength, contentMaxSecondary)

	return offsetX, offsetY
}

// easeInOutSine provides smooth sine-based easing
func easeInOutSine(t float64) float64 {
	return 0.5 * (1 - math.Cos(math.Pi*t))
}

// easeInOutQuad provides quadratic easing
func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - 2*(1-t)*(1-t)
}

// easeWithDwell provides easing with dwell time at start and end
// Pattern: 0000->1 1 1 1 1 1 1->0000
// This works for half-cycle (0->1), ping-pong conversion handles the return
func easeWithDwell(t float64, dwellPercent float64) float64 {
	if dwellPercent <= 0 || dwellPercent >= 0.5 {
		dwellPercent = 0.2 // Default 20% dwell time at each end
	}

	if t < dwellPercent {
		// Dwell at start (0)
		return 0.0
	} else if t > (1.0 - dwellPercent) {
		// Dwell at end (1)
		return 1.0
	} else {
		// Active movement phase - map to 0->1 using smooth easing
		activePhase := (t - dwellPercent) / (1.0 - 2*dwellPercent)
		return easeInOutSine(activePhase)
	}
}

// calculateEasedPosition calculates smooth eased scrolling position based on overflow distance
func (t *Template) calculateEasedPosition(startTime time.Time, overflow int) int {
	// Ensure overflow is positive - negative values mean no scrolling needed
	if overflow <= 0 {
		return 0
	}

	elapsed := time.Since(startTime)

	// Calculate speed-based duration: time to traverse full overflow distance
	// ScrollSpeed now represents pixels per second directly for better control
	// Lower values like 0.1, 0.5, 1.0 give much slower, more readable scrolling
	pixelsPerSecond := float64(t.ScrollSpeed)
	if pixelsPerSecond <= 0 {
		pixelsPerSecond = 10.0 // Default to 10 pixels per second if invalid
	}

	// Use absolute value to ensure positive calculation
	traversalTime := math.Abs(float64(overflow)) / pixelsPerSecond
	fullCycleDuration := time.Duration(traversalTime*2000) * time.Millisecond // *2 for round trip, *1000 for ms

	// Use CycleDurationMs if provided, otherwise use calculated duration
	var cycleDuration time.Duration
	if t.CycleDurationMs > 0 {
		cycleDuration = time.Duration(t.CycleDurationMs) * time.Millisecond
	} else {
		cycleDuration = fullCycleDuration
	}

	// Calculate progress in current cycle (0.0 to 1.0)
	progress := float64(elapsed%cycleDuration) / float64(cycleDuration)

	var easedProgress float64
	dwellPercent := float64(t.DwellPercent) / 100.0 // Convert percentage to decimal

	// Handle full cycle with ping-pong motion and dwell at both ends
	var halfCycleProgress float64

	if progress < 0.5 {
		// First half: 0 -> 1 (forward)
		halfCycleProgress = progress * 2.0
	} else {
		// Second half: 1 -> 0 (returning)
		halfCycleProgress = (1.0 - progress) * 2.0
	}

	switch t.EasingMode {
	case "ease-in-out-sine":
		easedProgress = easeWithDwell(halfCycleProgress, dwellPercent)
	case "ease-in-out-quad":
		// Apply dwell to quad easing too
		if halfCycleProgress < dwellPercent {
			easedProgress = 0.0
		} else if halfCycleProgress > (1.0 - dwellPercent) {
			easedProgress = 1.0
		} else {
			activePhase := (halfCycleProgress - dwellPercent) / (1.0 - 2*dwellPercent)
			easedProgress = easeInOutQuad(activePhase)
		}
	case "linear":
		easedProgress = halfCycleProgress
	default:
		easedProgress = easeWithDwell(halfCycleProgress, dwellPercent)
	}

	// No need for ping-pong conversion since we handled it above

	// Map to scroll position range - now covers full overflow distance
	return -int(easedProgress * math.Abs(float64(overflow)))
}

// calculateHorizontalOffset handles X-axis scrolling based on overflow-x
func (t *Template) calculateHorizontalOffset(contentLength, contentMaxSecondary int) int {
	var overflow int

	// For horizontal scrolling, we always want to compare the widest content
	// against the available horizontal space, regardless of layout direction
	if t.Direction == "row" || t.Direction == "" {
		// In row direction (or unset), contentLength is horizontal extent
		overflow = contentLength - t.ComputedSizeX
	} else {
		// In column direction, contentMaxSecondary is the width of widest item
		overflow = contentMaxSecondary - t.ComputedSizeX
	}

	// For scroll-ease, we want to scroll even if there's minimal overflow
	// This handles cases where text is just slightly smaller than container
	if t.OverflowX == "scroll-ease" && overflow >= -5 {
		// Force scrolling for text that's close to the container width
		if overflow < 5 {
			overflow = 5 // Minimum scroll distance
		}
	} else if overflow <= 0 {
		return 0 // No horizontal overflow
	}

	switch t.OverflowX {
	case "scroll", "scroll-left":
		t.PosX -= t.ScrollSpeed
		if t.PosX < -overflow {
			t.PosX = t.ComputedSizeX // Start from right side
		}
		return t.PosX

	case "scroll-right":
		t.PosX += t.ScrollSpeed
		if t.PosX > t.ComputedSizeX {
			t.PosX = -overflow // Start from left side
		}
		return t.PosX

	case "scroll-bounce":
		if t.EasingMode == "linear" {
			// Keep original bounce logic for linear mode
			if t.bounceBackX {
				if t.bounceDelayCountX > 0 {
					t.bounceDelayCountX--
					return -t.PosX
				}
				t.PosX -= t.ScrollSpeed
				if t.PosX <= 0 {
					t.PosX = 0
					t.bounceBackX = false
					t.bounceDelayCountX = t.BounceDelay
				}
			} else {
				if t.bounceDelayCountX > 0 {
					t.bounceDelayCountX--
					return -t.PosX
				}
				t.PosX += t.ScrollSpeed
				if t.PosX >= overflow {
					t.PosX = overflow
					t.bounceBackX = true
					t.bounceDelayCountX = t.BounceDelay
				}
			}
			return -t.PosX
		} else {
			// Use eased scrolling
			return t.calculateEasedPosition(t.easingStartTimeX, overflow)
		}

	case "scroll-ease":
		return t.calculateEasedPosition(t.easingStartTimeX, overflow)

	case "auto":
		// Auto-expansion mode - no scrolling needed as template expands
		return 0

	default: // "hidden", "visible", etc.
		return 0
	}
}

// calculateVerticalOffset handles Y-axis scrolling based on overflow-y
func (t *Template) calculateVerticalOffset(contentLength, contentMaxSecondary int) int {
	var overflow int
	// For vertical overflow, always measure vertical content vs vertical space
	// regardless of template direction
	if t.Direction == "col" {
		// In column direction, contentLength is the total height of all items
		overflow = contentLength - t.ComputedSizeY
	} else {
		// In row direction, contentMaxSecondary is the height of the tallest item
		overflow = contentMaxSecondary - t.ComputedSizeY
	}

	if overflow <= 0 {
		return 0 // No vertical overflow
	}

	switch t.OverflowY {
	case "scroll", "scroll-up":
		t.PosY -= t.ScrollSpeed
		if t.PosY < -overflow {
			t.PosY = t.ComputedSizeY // Start from bottom
		}
		return t.PosY

	case "scroll-down":
		t.PosY += t.ScrollSpeed
		if t.PosY > t.ComputedSizeY {
			t.PosY = -overflow // Start from top
		}
		return t.PosY

	case "scroll-bounce":
		if t.EasingMode == "linear" {
			// Keep original bounce logic for linear mode
			if t.bounceBackY {
				if t.bounceDelayCountY > 0 {
					t.bounceDelayCountY--
					return -t.PosY
				}
				t.PosY -= t.ScrollSpeed
				if t.PosY <= 0 {
					t.PosY = 0
					t.bounceBackY = false
					t.bounceDelayCountY = t.BounceDelay
				}
			} else {
				if t.bounceDelayCountY > 0 {
					t.bounceDelayCountY--
					return -t.PosY
				}
				t.PosY += t.ScrollSpeed
				if t.PosY >= overflow {
					t.PosY = overflow
					t.bounceBackY = true
					t.bounceDelayCountY = t.BounceDelay
				}
			}
			return -t.PosY
		} else {
			// Use eased scrolling
			return t.calculateEasedPosition(t.easingStartTimeY, overflow)
		}

	case "scroll-ease":
		return t.calculateEasedPosition(t.easingStartTimeY, overflow)

	case "auto":
		// Auto-expansion mode - no scrolling needed as template expands
		return 0

	default: // "hidden", "visible", etc.
		return 0
	}
}

func (t *Template) Stop() {
	for _, c := range t.Components {
		c.Stop()
	}
	t.BaseComponent.Stop()
}

func (tmpl *Template) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tmpl.XMLName = start.Name

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "size-x":
			tmpl.SizeX = attr.Value
		case "size-y":
			tmpl.SizeY = attr.Value
		case "justify":
			tmpl.Justify = attr.Value
		case "align":
			tmpl.Align = attr.Value
		case "dir":
			tmpl.Direction = attr.Value
		case "bg-color":
			tmpl.BgColor = attr.Value
		case "overflow":
			tmpl.Overflow = attr.Value
		case "overflow-x":
			tmpl.OverflowX = attr.Value
		case "overflow-y":
			tmpl.OverflowY = attr.Value
		case "scroll-speed":
			if speed, err := strconv.Atoi(attr.Value); err == nil {
				tmpl.ScrollSpeed = speed
			}
		case "bounce-delay":
			if delay, err := strconv.Atoi(attr.Value); err == nil {
				tmpl.BounceDelay = delay
			}
		case "easing":
			tmpl.EasingMode = attr.Value
		case "cycle-duration-ms":
			if duration, err := strconv.Atoi(attr.Value); err == nil {
				tmpl.CycleDurationMs = duration
			}
		case "dwell-percent":
			if dwell, err := strconv.Atoi(attr.Value); err == nil {
				tmpl.DwellPercent = dwell
			}
		}
	}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		var i Component
		switch tt := t.(type) {
		case xml.StartElement:
			i = RegisteredComponents[tt.Name.Local]()
			// log.Printf("Invalid component type %s", tt.Name.Local)
			if i != nil {
				err = d.DecodeElement(i, &tt)
				if err != nil {
					return err
				}
				tmpl.Components = append(tmpl.Components, i)
				i = nil
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}

	}
}

func init() {
	RegisterComponent("template", func() Component { return &Template{} })
}
