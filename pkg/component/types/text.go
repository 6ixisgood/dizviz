package types

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strings"

	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	fontpkg "golang.org/x/image/font"
)

type Text struct {
	c.BaseComponent

	XMLName   xml.Name  `xml:"text"`
	Font      string    `xml:"font,attr"`
	FontStyle string    `xml:"style,attr"`
	FontSize  float64   `xml:"size,attr"`
	Color     util.RGBA `xml:"color,attr"`
	WordWrap  bool      `xml:"word-wrap,attr"`
	Rainbow   bool      `xml:"rainbow,attr"`
	Text      string    `xml:",chardata"`

	img        *image.RGBA
	ftCtx      *freetype.Context
	lines      []string
	colorIndex int
}

func (t *Text) Init() {
	if t.Rainbow {
		t.Rr = 100
	} else {
		t.Rr = -1 // no need to rerender this once created
	}
	t.BaseComponent.Init()

	t.Ctx = gg.NewContext(0, 0)

	// init the font and style
	var font = util.LoadFont(fmt.Sprintf("%s-%s", t.Font, t.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: t.FontSize})
	t.Ctx.SetFontFace(face)

	// get the size of the string
	var w_i, h_i int
	if t.WordWrap && t.ComputedSizeX > 0 {
		t.lines = t.breakIntoLines(t.ComputedSizeX)

		// Find the width of the widest line
		maxLineWidth := 0
		for _, line := range t.lines {
			lineWidth, _ := t.Ctx.MeasureString(line)
			if int(math.Ceil(lineWidth)) > maxLineWidth {
				maxLineWidth = int(math.Ceil(lineWidth))
			}
		}

		w_i = int(math.Min(float64(t.ComputedSizeX), float64(maxLineWidth)))
		h_i = len(t.lines) * int(math.Ceil(t.FontSize*1.2))
	} else {
		w, h := t.Ctx.MeasureString(t.Text)
		w_i = int(math.Ceil(w))
		h_i = int(math.Ceil(h))
	}
	t.ComputedSizeX = w_i
	t.ComputedSizeY = h_i

	t.Ctx = gg.NewContext(t.ComputedSizeX, t.ComputedSizeY)
	// set up a blank image
	t.img = image.NewRGBA(image.Rect(0, 0, t.ComputedSizeX, t.ComputedSizeY))

	// Set up the freetype context
	t.ftCtx = freetype.NewContext()
	t.ftCtx.SetDPI(72)
	t.ftCtx.SetFont(font)
	t.ftCtx.SetFontSize(t.FontSize)
	t.ftCtx.SetClip(t.img.Bounds())
	t.ftCtx.SetDst(t.img)
	t.ftCtx.SetSrc(image.NewUniform(t.Color.RGBA)) // set the color
	t.ftCtx.SetHinting(fontpkg.HintingNone)
}

func (t *Text) breakIntoLines(maxWidth int) []string {
	words := strings.Fields(t.Text)
	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if currentLine != "" {
			testLine += " "
		}
		testLine += word

		w, _ := t.Ctx.MeasureString(testLine)
		if w > float64(maxWidth) {
			if currentLine == "" {
				lines = append(lines, t.breakWordByCharacter(word, maxWidth)...)
			} else {
				lines = append(lines, currentLine)
				wordWidth, _ := t.Ctx.MeasureString(word)
				if wordWidth > float64(maxWidth) {
					lines = append(lines, t.breakWordByCharacter(word, maxWidth)...)
				} else {
					currentLine = word
				}
			}
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (t *Text) breakWordByCharacter(word string, maxWidth int) []string {
	var lines []string
	var currentLine string

	for _, char := range word {
		testLine := currentLine + string(char)
		w, _ := t.Ctx.MeasureString(testLine)
		if w > float64(maxWidth) && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = string(char)
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (t *Text) renderRainbow() image.Image {
	t.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	t.Ctx.Clear()

	rainbowColors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{255, 127, 0, 255}, // Orange
		{255, 255, 0, 255}, // Yellow
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{75, 0, 130, 255},  // Indigo
		{148, 0, 211, 255}, // Violet
	}

	if t.WordWrap && len(t.lines) > 0 {
		lineHeight := int(t.FontSize * 1.2)
		for i, line := range t.lines {
			y := float64((i + 1) * lineHeight)
			startX := 0.0
			for _, char := range line {
				currentColor := rainbowColors[t.colorIndex]
				t.Ctx.SetColor(currentColor)
				charStr := string(char)
				t.Ctx.DrawString(charStr, startX, y)

				cw, _ := t.Ctx.MeasureString(charStr)
				startX += cw
				t.colorIndex = (t.colorIndex + 1) % len(rainbowColors)
			}
		}
	} else {
		_, h := t.Ctx.MeasureString(t.Text)
		startX := 0.0
		for _, char := range t.Text {
			currentColor := rainbowColors[t.colorIndex]
			t.Ctx.SetColor(currentColor)
			charStr := string(char)
			t.Ctx.DrawString(charStr, startX, h)

			cw, _ := t.Ctx.MeasureString(charStr)
			startX += cw
			t.colorIndex = (t.colorIndex + 1) % len(rainbowColors)
		}
	}

	return t.Ctx.Image()
}

func (t *Text) Render() image.Image {
	if t.Rainbow {
		return t.renderRainbow()
	}

	if t.WordWrap && len(t.lines) > 0 {
		lineHeight := int(t.FontSize * 1.2)

		for i, line := range t.lines {
			y := (i + 1) * lineHeight
			pt := freetype.Pt(0, y)
			_, err := t.ftCtx.DrawString(line, pt)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		pt := freetype.Pt(0, int(t.FontSize))
		_, err := t.ftCtx.DrawString(t.Text, pt)
		if err != nil {
			log.Fatal(err)
		}
	}

	return t.img
}

func init() {
	c.RegisterComponent("text", func() c.Component { return &Text{} })
	c.RegisterComponent("rainbow-text", func() c.Component {
		return &Text{Rainbow: true}
	})
}
