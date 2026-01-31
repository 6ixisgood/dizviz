package types

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"

	c "github.com/6ixisgood/disco/pkg/component/common"
	"github.com/6ixisgood/disco/pkg/util"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type ComparisonBar struct {
	c.BaseComponent

	XMLName        xml.Name  `xml:"comparison-bar"`
	LeftValue      float64   `xml:"left-value,attr"`
	RightValue     float64   `xml:"right-value,attr"`
	StatLabel      string    `xml:"stat-label,attr"`
	CenterLabel    string    `xml:"center-label,attr"` // Label in the middle of the bar
	LeftLabel      string    `xml:"left-label,attr"`
	RightLabel     string    `xml:"right-label,attr"`
	LeftColor      util.RGBA `xml:"left-color,attr"`
	RightColor     util.RGBA `xml:"right-color,attr"`
	BarHeight      int       `xml:"bar-height,attr"`
	LabelFontSize  float64   `xml:"label-font-size,attr"`
	ValueFontSize  float64   `xml:"value-font-size,attr"`
	Font           string    `xml:"font,attr"`
	FontStyle      string    `xml:"style,attr"`
	TextColor      util.RGBA `xml:"text-color,attr"`
	ShowStatLabel  bool      `xml:"show-stat-label,attr"`
	ShowValues     bool      `xml:"show-values,attr"`
	ValuePlacement string    `xml:"value-placement,attr"` // "above", "below", "side"

	// calculated properties
	leftBarWidth  int
	rightBarWidth int
}

func (cb *ComparisonBar) Init() {
	cb.Rr = -1 // static component, no need to rerender
	cb.BaseComponent.Init()

	// Set defaults
	if cb.BarHeight == 0 {
		cb.BarHeight = 8
	}

	if cb.LabelFontSize == 0 {
		cb.LabelFontSize = 10
	}

	if cb.ValueFontSize == 0 {
		cb.ValueFontSize = 12
	}

	if cb.Font == "" {
		cb.Font = "Roboto"
	}

	if cb.FontStyle == "" {
		cb.FontStyle = "Regular"
	}

	if cb.LeftColor.RGBA == (color.RGBA{}) {
		cb.LeftColor = util.RGBA{RGBA: color.RGBA{R: 255, G: 68, B: 68, A: 255}} // red
	}

	if cb.RightColor.RGBA == (color.RGBA{}) {
		cb.RightColor = util.RGBA{RGBA: color.RGBA{R: 68, G: 68, B: 255, A: 255}} // blue
	}

	if cb.TextColor.RGBA == (color.RGBA{}) {
		cb.TextColor = util.RGBA{RGBA: color.RGBA{R: 255, G: 255, B: 255, A: 255}} // white
	}

	if cb.ValuePlacement == "" {
		cb.ValuePlacement = "below"
	}

	// Calculate bar widths
	cb.calculateBarWidths()

	// Create context
	cb.Ctx = gg.NewContext(cb.ComputedSizeX, cb.ComputedSizeY)
}

func (cb *ComparisonBar) calculateBarWidths() {
	// Start with full width
	maxBarSpace := cb.ComputedSizeX / 2

	// If we're showing values on the side, we need to reserve space for them
	if cb.ShowValues && cb.ValuePlacement == "side" {
		font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))
		face := truetype.NewFace(font, &truetype.Options{Size: cb.ValueFontSize})
		ctx := gg.NewContext(1, 1)
		ctx.SetFontFace(face)

		// Measure both values and use the larger width
		leftValueStr := formatValue(cb.LeftValue)
		rightValueStr := formatValue(cb.RightValue)
		leftValueWidth, _ := ctx.MeasureString(leftValueStr)
		rightValueWidth, _ := ctx.MeasureString(rightValueStr)
		maxValueWidth := math.Max(leftValueWidth, rightValueWidth)

		// Reserve space for value text plus padding (4px padding + 4px between value and bar)
		valueSpace := int(math.Max(maxValueWidth+16, 20))
		maxBarSpace = maxBarSpace - valueSpace
	}

	// If there's a center label, we need to account for its width
	if cb.CenterLabel != "" {
		font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))
		face := truetype.NewFace(font, &truetype.Options{Size: cb.LabelFontSize})
		ctx := gg.NewContext(1, 1)
		ctx.SetFontFace(face)
		labelWidth, _ := ctx.MeasureString(cb.CenterLabel)

		// Reduce available space for each bar by half the label width plus padding
		padding := 4.0 // 4px padding on each side of the label
		maxBarSpace = (cb.ComputedSizeX - int(labelWidth) - int(padding*2)) / 2

		// Also account for side values if present
		if cb.ShowValues && cb.ValuePlacement == "side" {
			font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))
			face := truetype.NewFace(font, &truetype.Options{Size: cb.ValueFontSize})
			ctx := gg.NewContext(1, 1)
			ctx.SetFontFace(face)

			leftValueStr := formatValue(cb.LeftValue)
			rightValueStr := formatValue(cb.RightValue)
			leftValueWidth, _ := ctx.MeasureString(leftValueStr)
			rightValueWidth, _ := ctx.MeasureString(rightValueStr)

			// Each side needs to fit: edge padding + value + spacing + bar
			maxBarSpace = maxBarSpace - int(math.Max(leftValueWidth, rightValueWidth)) - 8
		}
	} else if !cb.ShowValues || cb.ValuePlacement != "side" {
		// Reserve some space for edges when no side values
		maxBarSpace = maxBarSpace - 10
	}

	// Ensure we don't have negative space
	if maxBarSpace < 0 {
		maxBarSpace = 0
	}

	// Find max value for normalization
	maxValue := math.Max(cb.LeftValue, cb.RightValue)

	if maxValue > 0 && maxBarSpace > 0 {
		cb.leftBarWidth = int((cb.LeftValue / maxValue) * float64(maxBarSpace))
		cb.rightBarWidth = int((cb.RightValue / maxValue) * float64(maxBarSpace))
	} else {
		cb.leftBarWidth = 0
		cb.rightBarWidth = 0
	}
}

func (cb *ComparisonBar) Render() image.Image {
	// Clear the context
	cb.Ctx.SetColor(color.RGBA{0, 0, 0, 0})
	cb.Ctx.Clear()

	// Calculate layout
	centerX := float64(cb.ComputedSizeX) / 2
	currentY := 0.0

	// Load font
	font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))

	// Draw stat label if enabled
	if cb.ShowStatLabel && cb.StatLabel != "" {
		face := truetype.NewFace(font, &truetype.Options{Size: cb.LabelFontSize})
		cb.Ctx.SetFontFace(face)
		cb.Ctx.SetColor(cb.TextColor.RGBA)

		labelWidth, labelHeight := cb.Ctx.MeasureString(cb.StatLabel)
		cb.Ctx.DrawString(cb.StatLabel, centerX-labelWidth/2, currentY+labelHeight)
		currentY += labelHeight + 4
	}

	// Draw team labels if provided
	if cb.LeftLabel != "" || cb.RightLabel != "" {
		face := truetype.NewFace(font, &truetype.Options{Size: cb.LabelFontSize})
		cb.Ctx.SetFontFace(face)
		cb.Ctx.SetColor(cb.TextColor.RGBA)

		_, labelHeight := cb.Ctx.MeasureString("Ag") // Use Ag for consistent height

		if cb.LeftLabel != "" {
			cb.Ctx.DrawString(cb.LeftLabel, 5, currentY+labelHeight)
		}

		if cb.RightLabel != "" {
			rightLabelWidth, _ := cb.Ctx.MeasureString(cb.RightLabel)
			cb.Ctx.DrawString(cb.RightLabel, float64(cb.ComputedSizeX)-rightLabelWidth-5, currentY+labelHeight)
		}

		currentY += labelHeight + 4
	}

	// Calculate bar Y position (center the bar vertically in remaining space)
	remainingHeight := float64(cb.ComputedSizeY) - currentY
	if cb.ShowValues && (cb.ValuePlacement == "above" || cb.ValuePlacement == "below") {
		remainingHeight -= cb.ValueFontSize + 4
	}
	barY := currentY + (remainingHeight-float64(cb.BarHeight))/2

	// Draw values above bar if enabled
	if cb.ShowValues && cb.ValuePlacement == "above" {
		cb.drawValues(barY - cb.ValueFontSize - 2)
	}

	// Calculate center position and label width if there's a center label
	var labelWidth float64
	if cb.CenterLabel != "" {
		font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))
		face := truetype.NewFace(font, &truetype.Options{Size: cb.LabelFontSize})
		cb.Ctx.SetFontFace(face)
		labelWidth, _ = cb.Ctx.MeasureString(cb.CenterLabel)
	}

	// Adjust bar positions to account for center label
	padding := 4.0
	leftBarEnd := centerX - labelWidth/2 - padding
	rightBarStart := centerX + labelWidth/2 + padding

	// Draw left bar (growing leftward from label)
	if cb.leftBarWidth > 0 {
		cb.Ctx.SetColor(cb.LeftColor.RGBA)
		leftBarX := leftBarEnd - float64(cb.leftBarWidth)
		cb.Ctx.DrawRectangle(leftBarX, barY, float64(cb.leftBarWidth), float64(cb.BarHeight))
		cb.Ctx.Fill()
	}

	// Draw center label if provided
	if cb.CenterLabel != "" {
		cb.drawCenterLabel(barY, centerX, labelWidth)
	}

	// Draw right bar (growing rightward from label)
	if cb.rightBarWidth > 0 {
		cb.Ctx.SetColor(cb.RightColor.RGBA)
		cb.Ctx.DrawRectangle(rightBarStart, barY, float64(cb.rightBarWidth), float64(cb.BarHeight))
		cb.Ctx.Fill()
	}

	// Draw values based on placement
	if cb.ShowValues {
		switch cb.ValuePlacement {
		case "below":
			cb.drawValuesBelow(barY)
		case "side":
			cb.drawValuesSide(barY, centerX)
		}
	}

	return cb.Ctx.Image()
}

func (cb *ComparisonBar) drawValues(valueY float64) {
	font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))
	face := truetype.NewFace(font, &truetype.Options{Size: cb.ValueFontSize})
	cb.Ctx.SetFontFace(face)

	// Left value
	leftValueStr := formatValue(cb.LeftValue)
	cb.Ctx.SetColor(cb.LeftColor.RGBA)
	cb.Ctx.DrawString(leftValueStr, 5, valueY)

	// Right value
	rightValueStr := formatValue(cb.RightValue)
	rightValueWidth, _ := cb.Ctx.MeasureString(rightValueStr)
	cb.Ctx.SetColor(cb.RightColor.RGBA)
	cb.Ctx.DrawString(rightValueStr, float64(cb.ComputedSizeX)-rightValueWidth-5, valueY)
}

func (cb *ComparisonBar) drawValuesBelow(barY float64) {
	valueY := barY + float64(cb.BarHeight) + cb.ValueFontSize + 4
	cb.drawValues(valueY)
}

func (cb *ComparisonBar) drawValuesSide(barY float64, centerX float64) {
	font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))
	face := truetype.NewFace(font, &truetype.Options{Size: cb.ValueFontSize})
	cb.Ctx.SetFontFace(face)

	// Calculate vertical center of bar for text alignment
	_, textHeight := cb.Ctx.MeasureString("0")
	valueY := barY + float64(cb.BarHeight)/2 + textHeight/2

	// Measure the values
	leftValueStr := formatValue(cb.LeftValue)
	rightValueStr := formatValue(cb.RightValue)
	rightValueWidth, _ := cb.Ctx.MeasureString(rightValueStr)

	// Left value (left-aligned at edge)
	cb.Ctx.SetColor(cb.LeftColor.RGBA)
	cb.Ctx.DrawString(leftValueStr, 5, valueY)

	// Right value (right-aligned at edge)
	cb.Ctx.SetColor(cb.RightColor.RGBA)
	cb.Ctx.DrawString(rightValueStr, float64(cb.ComputedSizeX)-rightValueWidth-5, valueY)
}

func (cb *ComparisonBar) drawCenterLabel(barY float64, centerX float64, labelWidth float64) {
	font := util.LoadFont(fmt.Sprintf("%s-%s", cb.Font, cb.FontStyle))
	face := truetype.NewFace(font, &truetype.Options{Size: cb.LabelFontSize})
	cb.Ctx.SetFontFace(face)

	// Measure the label
	_, textHeight := cb.Ctx.MeasureString(cb.CenterLabel)

	// Calculate position: centered horizontally, vertically centered on the bar
	labelX := centerX - labelWidth/2
	labelY := barY + float64(cb.BarHeight)/2 + textHeight/2

	// Draw a background rectangle for better readability
	padding := 2.0
	cb.Ctx.SetColor(color.RGBA{0, 0, 0, 200}) // semi-transparent black background
	cb.Ctx.DrawRectangle(
		labelX-padding,
		barY+float64(cb.BarHeight)/2-textHeight/2-padding,
		labelWidth+padding*2,
		textHeight+padding*2,
	)
	cb.Ctx.Fill()

	// Draw the label text
	cb.Ctx.SetColor(cb.TextColor.RGBA)
	cb.Ctx.DrawString(cb.CenterLabel, labelX, labelY)
}

func formatValue(val float64) string {
	// If it's a whole number, don't show decimals
	if val == float64(int(val)) {
		return strconv.Itoa(int(val))
	}
	return fmt.Sprintf("%.1f", val)
}

func init() {
	c.RegisterComponent("comparison-bar", func() c.Component { return &ComparisonBar{} })
}
