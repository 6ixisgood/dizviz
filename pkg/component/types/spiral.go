package types

import (
	"encoding/xml"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"github.com/fogleman/gg"
	"image"
	"math"
)

type SpiralGallery struct {
	c.BaseComponent

	XMLName     xml.Name     `xml:"spiral"`
	Slots       []c.Template `xml:"template"`
	Angle       float64      // Rotation angle in degrees
	CurrentSlot int
}

func (sg *SpiralGallery) Init() {
	sg.BaseComponent.Init()

	// Initialize all slots
	for _, slot := range sg.Slots {
		slot.Init()
	}
}

func (sg *SpiralGallery) Render() image.Image {
	numSlots := len(sg.Slots)

	if numSlots == 0 {
		return sg.Ctx.Image()
	}

	// Calculate rotation angle step
	angleStep := 360.0 / float64(numSlots)

	// Render the slots in a spiral manner
	for i, slot := range sg.Slots {
		// Determine the angle for this slot
		currentAngle := sg.Angle + angleStep*float64(i)

		// Convert polar coordinates (r, theta) to Cartesian (x, y)
		// r is the distance from the center, theta is the angle from the positive x-axis
		r := float64(sg.Width()) / 3.0 // Let's position our images in a circle which is a third of our component's width
		x := r*math.Cos(gg.Radians(currentAngle)) + float64(sg.Width()/2)
		y := r*math.Sin(gg.Radians(currentAngle)) + float64(sg.Height()/2)

		img := slot.Render()
		sg.Ctx.DrawImageAnchored(img, int(x), int(y), 0.5, 0.5) // Anchored at center
	}

	// Update the angle for the next render
	sg.Angle += 5.0 // Rotate by 5 degrees. This can be adjusted for faster or slower rotation
	if sg.Angle >= 360 {
		sg.Angle = 0
		sg.CurrentSlot = (sg.CurrentSlot + 1) % numSlots
	}

	return sg.Ctx.Image()
}

func init() {
	c.RegisterComponent("spiral", func() c.Component { return &SpiralGallery{} })
}
