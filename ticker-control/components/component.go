package components

import (
	"log"
	"os"
	"image"
	"math"
	"time"
	"image/color"
	"image/gif"
	// "image/png"
	"bytes"
	"strings"
	"fmt"
	"math/rand"
	"encoding/xml"
	"strconv"
	"runtime"
	"io/ioutil"
	"path/filepath"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	fontpkg "golang.org/x/image/font"
	d "github.com/sixisgoood/matrix-ticker/data"
)

var (
	imageCache = make(map[string]image.Image)
)

type Component interface {
	Init()								// Ran before componenet is rendered
	Render()					image.Image		// Render the component to an image.Image representation
	Width()						int				// Return the width of the component. Used to help position components on the display
	Height()					int				// Return the height of the compoent. Used to help position components on the display
	PrevImg()					image.Image		// return the previously rendered image
	SetPrevImg(img image.Image)					// the the previously rendered iamge
	TickerChan()				<-chan time.Time // return the channel for the ticker
	Stop()										// Stop the ticker

}

type BaseComponent struct {
	SizeX		int			`xml:"sizeX,attr"`
	SizeY		int			`xml:"sizeY,attr"`
	PosX		int			`xml:"posX,attr"`
	PosY		int			`xml:"posY,attr"`
	ctx			*gg.Context
	prevImg		image.Image
	Ticker		*time.Ticker
	rr			int			// render rate in milliseconds
}

func (bc *BaseComponent) Init() {
	// create a context, if needed
	if (bc.ctx == nil && bc.SizeX > 0 && bc.SizeY > 0) {
		bc.ctx = gg.NewContext(bc.SizeX, bc.SizeY)
	}
	// set a render rate, if needed
	if bc.rr == 0 {
		bc.rr = 100
	}

	// set ticker for render rate
	if bc.rr > 0 {
		bc.Ticker = time.NewTicker(time.Duration(bc.rr) * time.Millisecond)
	}
}

func (bc *BaseComponent) Width() int {
	return bc.ctx.Width()
}

func (bc *BaseComponent) Height() int {
	return bc.ctx.Height()
}

func (bc *BaseComponent) Size() image.Point {
	return image.Point{bc.ctx.Width(), bc.ctx.Height()}
}

func (bc *BaseComponent) TickerChan() <-chan time.Time{
	if bc.Ticker != nil {
		return bc.Ticker.C
	} else {
		return nil
	}
}

func (bc *BaseComponent) Stop() {
	// pass
	if bc.Ticker != nil {
		bc.Ticker.Stop()
	}
}

func (bc *BaseComponent) PrevImg() image.Image {
	return bc.prevImg
}

func (bc *BaseComponent) SetPrevImg(img image.Image) {
	bc.prevImg = img
}


type Template struct {
	XMLName			xml.Name		`xml:"template"`
	SizeX			int				`xml:"sizeX,attr"`
	SizeY			int				`xml:"sizeY,attr"`
	Components		[]Component		`xml:",any"`

	ctx				*gg.Context

}

func (t *Template) Init() {
	ctxTmp := gg.NewContext(t.SizeX, t.SizeY)
	t.ctx = ctxTmp 

	for _, c  := range t.Components {
		c.Init()
	}
}

func (t *Template) Ready() bool {
	return t.ctx != nil
}

func (t *Template) ComponentWidth() int {
	sizeX := 0
	for _, c := range t.Components {
		sizeX += c.Width()
	}
	return sizeX
}

func (t *Template) Render() image.Image {
	if t.ctx == nil {
		t.ctx = gg.NewContext(t.SizeX, t.SizeY)
	}

	t.ctx.SetColor(color.RGBA{0,0,0,255})
	t.ctx.Clear()

	posX := 0
	var cIm image.Image
	for _, c  := range t.Components {
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
		// draw the image to the context
		t.ctx.DrawImage(cIm, posX, 0)
		posX += c.Width()
	}

	return t.ctx.Image()
}

func (tmpl *Template) Stop() {
	for _, c  := range tmpl.Components {
		c.Stop()
	}
}

func (tmpl *Template) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	log.Printf("Starting Unmarshalling of Template")
	tmpl.XMLName = start.Name

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "sizeX":
			tmpl.SizeX, _ = strconv.Atoi(attr.Value)
		case "sizeY":
			tmpl.SizeY, _ = strconv.Atoi(attr.Value)
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
			switch tt.Name.Local {
			case "text":
				i = new(Text)
			case "image":
				i = new(Image)
			case "scroller":
				i = new(Scroller)
			case "h-split":
				i = new(HorizonalSplit)
			case "rainbow-text":
				i = new(AnimatedRainbowText)
			case "spiral":
				i = new(SpiralGallery)
			case "pong":
				i = new(PaddleBallVisualizer)
			case "gravity-particles":
				i = new(GravityParticles)
			case "scenic-train":
				i = new(ScenicTrain)
			case "colorwave":
				i = new(ColorWave)
			case "matrix-rain":
				i = new(MatrixRain)
			default:
				log.Printf("Invalid component type %s", tt.Name.Local)
			}
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
	return nil
}

type Text struct {
	BaseComponent

	XMLName			xml.Name		`xml:"text"`
	Font			string			`xml:"font,attr"`		
	FontStyle		string			`xml:"style,attr"`	
	FontSize		float64			`xml:"size,attr"`	
	Color			RGBA			`xml:"color,attr"`
	Text			string			`xml:",chardata"`

	img         	*image.RGBA
	ftCtx			*freetype.Context
}

func (t *Text) Init() {
	t.rr = -1 // no need to rerender this once created
	t.BaseComponent.Init()

	t.ctx = gg.NewContext(0, 0)

	// init the font and style
	var font = loadFont(fmt.Sprintf("%s-%s", t.Font, t.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: t.FontSize})
	t.ctx.SetFontFace(face)

	// get the size of the string
	w, h := t.ctx.MeasureString(t.Text)
	w_i := int(math.Ceil(w))
	h_i := int(math.Ceil(h))
	if t.SizeX == 0 {
		t.SizeX = w_i
	}
	if t.SizeY == 0 {
		t.SizeY = h_i
	}


	// resize context
	t.ctx = gg.NewContext(t.SizeX, t.SizeY)
	// set up a blank image
	t.img = image.NewRGBA(image.Rect(0, 0, t.SizeX, t.SizeY))

	// Set up the freetype context
	t.ftCtx = freetype.NewContext()
	t.ftCtx.SetDPI(72)
	t.ftCtx.SetFont(font)
	t.ftCtx.SetFontSize(t.FontSize)
	t.ftCtx.SetClip(t.img.Bounds())
	t.ftCtx.SetDst(t.img)
	t.ftCtx.SetSrc(image.NewUniform(t.Color.RGBA)) // set the color
	t.ftCtx.SetHinting(fontpkg.HintingFull) 
}


func (t *Text) Render() image.Image {

	pt := freetype.Pt(0, int(t.FontSize))
	_, err := t.ftCtx.DrawString(t.Text, pt)
	if err != nil {
		log.Fatal(err)
	}

	return t.img
}

type Image struct {
	BaseComponent

	XMLName			xml.Name		`xml:"image"`
	Src				string			`xml:"src,attr"`
	Loop			bool			`xml:"loop,attr"`
	frames			[]image.Image
	currentFrame	int
}

func (i *Image) Init() {
	i.rr = -1
    i.BaseComponent.Init()


    data, filePath, err := d.FetchImage(i.Src, i.SizeX, i.SizeY)
    if err != nil {
    	log.Fatal(err)
    }

	extension := strings.ToLower(filepath.Ext(filePath))

    // Decode based on extension
    if extension == ".gif" {
        gifData, err := gif.DecodeAll(bytes.NewReader(data))
        if err != nil {
            log.Fatal(err)
        }

		for _, frame := range gifData.Image {
            i.frames = append(i.frames, frame)
        }

        // If the GIF should loop, reset the ticker accordingly
        if i.Loop {
            i.Ticker = time.NewTicker(time.Duration(gifData.Delay[0]) * 10 * time.Millisecond)
        }
    } else if extension == ".png" {
    	// do we have it saved in this size yet?
        img, _, err := image.Decode(bytes.NewReader(data))
        if err != nil {
            log.Fatal(err)
        }
        i.frames = append(i.frames, img)
    } else {
        log.Fatal("Unsupported file extension")
    }
}

func (i *Image) Width() int {
	return i.SizeX
}

func (i *Image) Height() int {
	return i.SizeY
}

func (i *Image) Render() image.Image {
	// Return the current frame and update the frame index
	img := i.frames[i.currentFrame]
	i.currentFrame = (i.currentFrame + 1) % len(i.frames)
	return img
}


// Scroller Component

type Scroller struct {
	BaseComponent

	XMLName			xml.Name		`xml:"scroller"`
	ScrollX			int				`xml:"scrollX,attr"`
	ScrollY			int				`xml:"scrollY,attr`
	Slot			Template		`xml:"template"`
}

func (s *Scroller) Init() {
	s.rr = 10 // render this one a bit faster than most
	s.BaseComponent.Init()
	s.Slot.Init()
}

func (s *Scroller) Render() image.Image {
	// render the slot
	im := s.Slot.Render()

	if s.ctx == nil {
		s.ctx = gg.NewContext(s.Slot.ComponentWidth(), 64)
	}

	s.ctx.SetColor(color.RGBA{0,0,0,255})
	s.ctx.Clear()	

	s.ctx.DrawImage(im, s.PosX, s.PosY)

	s.PosX = s.PosX + s.ScrollX
	s.PosY = s.PosY + s.ScrollY

	// wrap around
	if s.ScrollX < 0 {
		if s.PosX+s.ctx.Width() < 0 {
			s.PosX = 0
		}
	}

	return s.ctx.Image()	
}

type HorizonalSplit struct {
	BaseComponent

	XMLName			xml.Name		`xml:"h-split"`
	Slots			[]Template		`xml:"template"`
}

func (s *HorizonalSplit) Init() {
	s.rr = -1 // Should only have to render once
	s.BaseComponent.Init()
	for _, s := range s.Slots {
		s.Init()
	}

}

func (s *HorizonalSplit) Render() image.Image {
	width := 0
	for _, slot := range s.Slots {
		slot.Render()
		width = int(math.Max(float64(slot.ComponentWidth()), float64(width)))
	}
	height := 64
	if s.ctx == nil {
		s.ctx = gg.NewContext(width, height)
	}

	s.ctx.SetColor(color.RGBA{0,0,0,255})
	s.ctx.Clear()

	// render and draw the slots
	var im image.Image
	var y int
	for _, slot := range s.Slots {
		im = slot.Render()
		s.ctx.DrawImage(im, s.PosX, y)
		y += height/len(s.Slots)
	}

	return s.ctx.Image()	
}



// Rainbow text
type AnimatedRainbowText struct {
	BaseComponent

	XMLName      xml.Name    `xml:"rainbow-text"`
	Font         string      `xml:"font,attr"`		
	FontStyle    string      `xml:"style,attr"`	
	FontSize     float64     `xml:"size,attr"`	
	Text         string      `xml:",chardata"`
	colorIndex   int
}

func (art *AnimatedRainbowText) Init() {
	art.BaseComponent.Init()
	art.ctx = gg.NewContext(200, 64)


	var font = loadFont(fmt.Sprintf("%s-%s", art.Font, art.FontStyle))
	var face = truetype.NewFace(font, &truetype.Options{Size: art.FontSize})
	art.ctx.SetFontFace(face)
}

func (art *AnimatedRainbowText) Render() image.Image {
	art.ctx.SetColor(color.RGBA{0,0,0,255})
	art.ctx.Clear()
	rainbowColors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{255, 127, 0, 255}, // Orange
		{255, 255, 0, 255}, // Yellow
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{75, 0, 130, 255},  // Indigo
		{148, 0, 211, 255}, // Violet
	}

	// Get the size of the text
	_, h := art.ctx.MeasureString(art.Text)
	startX := 0.0

	for _, char := range art.Text {
		currentColor := rainbowColors[art.colorIndex]
		art.ctx.SetColor(currentColor)
		charStr := string(char)
		art.ctx.DrawString(charStr, startX, h) // Draw each character

		// Update starting x-coordinate for next character
		cw, _ := art.ctx.MeasureString(charStr)
		startX += cw

		// Update color index for the next character
		art.colorIndex = (art.colorIndex + 1) % len(rainbowColors)
	}

	return art.ctx.Image()
}


// Spiral 
type SpiralGallery struct {
	BaseComponent

	XMLName     xml.Name     `xml:"spiral"`
	Slots       []Template   `xml:"template"`
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
		return sg.ctx.Image()
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
		sg.ctx.DrawImageAnchored(img, int(x), int(y), 0.5, 0.5) // Anchored at center
	}

	// Update the angle for the next render
	sg.Angle += 5.0 // Rotate by 5 degrees. This can be adjusted for faster or slower rotation
	if sg.Angle >= 360 {
		sg.Angle = 0
		sg.CurrentSlot = (sg.CurrentSlot + 1) % numSlots
	}

	return sg.ctx.Image()
}


// Pong game
type PaddleBallVisualizer struct {
	BaseComponent

	XMLName          xml.Name    `xml:"pong"`
	BallRadius       float64     `xml:"ballRadius,attr"`
	BallSpeedX       float64
	BallSpeedY       float64
	BallX			 float64
	BallY			 float64
	PaddleHeight     float64     `xml:"paddleHeight,attr"`
	PaddleWidth      float64     `xml:"paddleWidth,attr"`
	LeftPaddleY      float64
	RightPaddleY     float64
	Amplitude        float64     // this will be updated by an external function based on the music beat
	Color            RGBA        `xml:"color,attr"`
}

func (pbv *PaddleBallVisualizer) Init() {
	pbv.BaseComponent.Init()

	// Set initial values
	pbv.BallSpeedX = 1.0
	pbv.BallSpeedY = 0.0
	pbv.LeftPaddleY = float64(pbv.Height())/2 - pbv.PaddleHeight/2
	pbv.RightPaddleY = pbv.LeftPaddleY

	pbv.BallX = float64(pbv.Width())/2
	pbv.BallY = float64(pbv.Height())/2
}

func (pbv *PaddleBallVisualizer) Render() image.Image {
	pbv.ctx.SetColor(color.RGBA{0,0,0,255})
	pbv.ctx.Clear()
	// Move the ball
	pbv.BallX += pbv.BallSpeedX
	pbv.BallY += pbv.BallSpeedY

	// Ball collision with top and bottom
	if pbv.BallY-pbv.BallRadius < 0 || pbv.BallY+pbv.BallRadius > float64(pbv.Height()) {
		pbv.BallSpeedY = -pbv.BallSpeedY
	}

	// Ball collision with paddles
	if (pbv.BallX-pbv.BallRadius < pbv.PaddleWidth && pbv.BallY > pbv.LeftPaddleY && pbv.BallY < pbv.LeftPaddleY+pbv.PaddleHeight) || 
	   (pbv.BallX+pbv.BallRadius > float64(pbv.Width())-pbv.PaddleWidth && pbv.BallY > pbv.RightPaddleY && pbv.BallY < pbv.RightPaddleY+pbv.PaddleHeight) {
		pbv.BallSpeedX = -pbv.BallSpeedX
	}

	// Move the paddles based on amplitude (mimic sound beat)
	pbv.LeftPaddleY = float64(pbv.Height())/2 - pbv.PaddleHeight/2 + pbv.Amplitude
	pbv.RightPaddleY = float64(pbv.Height())/2 - pbv.PaddleHeight/2 - pbv.Amplitude

	// Draw everything
	pbv.ctx.SetColor(pbv.Color.RGBA)

	// Draw the ball
	pbv.ctx.DrawCircle(pbv.BallX, pbv.BallY, pbv.BallRadius)
	pbv.ctx.Fill()

	// Draw the paddles
	pbv.ctx.DrawRectangle(0, pbv.LeftPaddleY, pbv.PaddleWidth, pbv.PaddleHeight)
	pbv.ctx.DrawRectangle(float64(pbv.Width())-pbv.PaddleWidth, pbv.RightPaddleY, pbv.PaddleWidth, pbv.PaddleHeight)
	pbv.ctx.Fill()

	return pbv.ctx.Image()
}

// Gravity Particles

type Particle struct {
	X, Y     float64
	SpeedX   float64
	SpeedY   float64
	Color    color.RGBA
	Radius   float64
}

type GravityParticle struct {
	X, Y    float64
	Force   float64
	Color   color.RGBA
}

type GravityParticles struct {
	BaseComponent

	XMLName         xml.Name         `xml:"gravity-particles"`
	Particles       []Particle
	GravityPoints   []GravityParticle
}

func (gp *GravityParticles) Init() {
	gp.BaseComponent.Init()

	// Sample initialization for particles and gravity points
	// In a real-world scenario, these could be populated based on the music or other parameters
	for i := 0; i < 50; i++ {
		gp.Particles = append(gp.Particles, Particle{
			X:      float64(rand.Intn(gp.Width())),
			Y:      float64(rand.Intn(gp.Height())),
			SpeedX: float64(rand.Intn(5)-2), // Random speed between -2 and 2
			SpeedY: float64(rand.Intn(5)-2),
			Color:  color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
			Radius: 2,
		})
	}

	gp.GravityPoints = append(gp.GravityPoints, GravityParticle{
		X:     float64(gp.Width()) / 2,
		Y:     float64(gp.Height()) / 2,
		Force: 5,
		Color: color.RGBA{255, 0, 0, 255},
	})
}

func (gp *GravityParticles) Render() image.Image {
	gp.ctx.SetColor(color.RGBA{0, 0, 0, 255})
	gp.ctx.Clear()

	// Update and render particles
	for i, particle := range gp.Particles {
		for _, g := range gp.GravityPoints {
			// Calculate the force exerted by this gravity point
			distX := g.X - particle.X
			distY := g.Y - particle.Y
			distance := math.Sqrt(distX*distX + distY*distY)

			// The closer the particle is to the gravity point, the stronger the pull
			force := g.Force / (distance + 1) // +1 to avoid division by zero
			particle.SpeedX += force * distX / distance
			particle.SpeedY += force * distY / distance
		}

		particle.X += particle.SpeedX
		particle.Y += particle.SpeedY

		// Boundary checking
		if particle.X < 0 || particle.X > float64(gp.Width()) {
			particle.SpeedX = -particle.SpeedX
		}
		if particle.Y < 0 || particle.Y > float64(gp.Height()) {
			particle.SpeedY = -particle.SpeedY
		}

		// Update the particle in the slice
		gp.Particles[i] = particle

		// Draw the particle
		gp.ctx.SetColor(particle.Color)
		gp.ctx.DrawCircle(particle.X, particle.Y, particle.Radius)
		gp.ctx.Fill()
	}


	return gp.ctx.Image()
}


// train
type Mountain struct {
	X, Y      float64
	Width     float64
	Height    float64
	Color     color.RGBA
}

type River struct {
	X, Y      float64
	Width     float64
	Height    float64
	Color     color.RGBA
}

type Train struct {
	X         float64
	Y         float64
	SpeedX    float64
	Carriages int
	Color     color.RGBA
}

type ScenicTrain struct {
	BaseComponent

	XMLName   xml.Name   `xml:"scenic-train"`
	Mountains []Mountain
	River     River
	Train     Train
	skip	  int
}

func (st *ScenicTrain) Init() {
	st.BaseComponent.Init()

	// Sample initialization of mountains
	for i := 0; i < 4; i++ {
		st.Mountains = append(st.Mountains, Mountain{
			X:      float64(rand.Intn(st.Width())),
			Y:      float64(st.Height() - 20 - rand.Intn(20)),
			Width:  60 + float64(rand.Intn(40)),
			Height: 40 + float64(rand.Intn(20)),
			Color:  color.RGBA{uint8(120 + rand.Intn(50)), uint8(70 + rand.Intn(40)), uint8(30), 255},
		})
	}

	st.River = River{
		X:      0,
		Y:      float64(st.Height() - 20),
		Width:  float64(st.Width()),
		Height: 20,
		Color:  color.RGBA{0, 0, 255, 255},
	}

	st.Train = Train{
		X:         -300,
		Y:         float64(st.Height() - 45),
		SpeedX:    4,
		Carriages: 5,
		Color:     color.RGBA{255, 0, 0, 255},
	}
	st.skip = 100
}

func (st *ScenicTrain) Render() image.Image {
	st.ctx.SetColor(color.RGBA{0, 0, 0, 255})
	st.ctx.Clear()

	// Render mountains
	for _, mtn := range st.Mountains {
		st.ctx.SetColor(mtn.Color)
		st.ctx.DrawRoundedRectangle(mtn.X, mtn.Y, mtn.Width, mtn.Height, 10)
		st.ctx.Fill()
	}

	// Render the river
	st.ctx.SetColor(st.River.Color)
	st.ctx.DrawRectangle(st.River.X, st.River.Y, st.River.Width, st.River.Height)
	st.ctx.Fill()

	// Render the train
	st.ctx.SetColor(st.Train.Color)
	for i := 0; i < st.Train.Carriages; i++ {
		st.ctx.DrawRoundedRectangle(st.Train.X+float64(i*70), st.Train.Y, 60, 25, 5)
		st.ctx.Fill()
	}

	// Update train position
	st.Train.X += st.Train.SpeedX
	if st.Train.X > float64(st.Width()) {
		st.Train.X = -300
	}

	return st.ctx.Image()
}


type Drop struct {
	x, y, speed float64
}

type MatrixRain struct {
	BaseComponent
	XMLName xml.Name `xml:"matrix-rain"`
	Drops   []Drop
	NumDrops int
}

func (mr *MatrixRain) Init() {
	mr.BaseComponent.Init()
	mr.NumDrops = 100
	for i := 0; i < mr.NumDrops; i++ {
		mr.Drops = append(mr.Drops, Drop{rand.Float64() * float64(mr.Width()), rand.Float64() * float64(mr.Height()), rand.Float64() * 5 + 1})
	}
}

func (mr *MatrixRain) Render() image.Image {
	mr.ctx.SetColor(color.RGBA{0, 0, 0, 255})
	mr.ctx.Clear()
	for _, drop := range mr.Drops {
		char := rune(33 + rand.Intn(94)) // Select a random ASCII character
		drop.y += drop.speed
		if drop.y > float64(mr.Height()) {
			drop.y = 0
		}
		mr.ctx.SetColor(color.RGBA{0, 255, 0, 255})
		mr.ctx.DrawStringAnchored(string(char), drop.x, drop.y, 0.5, 0.5)
	}
	return mr.ctx.Image()
}

// Color Wave Visual

type ColorWave struct {
	BaseComponent
	XMLName   xml.Name      `xml:"colorwave"`
	Colors    []color.RGBA  // Array of colors for the waves
	Phase     float64       // Current phase for the wave function
	WaveSpeed float64       // Speed at which the wave progresses
}

func (cw *ColorWave) Init() {
	cw.BaseComponent.Init()
	cw.Colors = []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	cw.Phase = 0
	cw.WaveSpeed = 0.05
}

func (cw *ColorWave) Render() image.Image {
	cw.ctx.SetColor(color.RGBA{0, 0, 0, 255})
	cw.ctx.Clear()
	for _, col := range cw.Colors {
		for x := 0; x < cw.Width(); x++ {
			y := math.Sin(float64(x)/float64(cw.Width())*2*math.Pi+cw.Phase) * float64(cw.Height()/4) + float64(cw.Height()/2)
			cw.ctx.SetColor(col)
			cw.ctx.DrawPoint(float64(x), y, 2)
		}
		cw.Phase += cw.WaveSpeed
	}
	cw.ctx.Fill()
	return cw.ctx.Image()
}



// HELPER FUNCTIONS

// RGBA Struct wraps color.RGBA for unmarshalling from XML
type RGBA struct {
	color.RGBA
}

func (c *RGBA) UnmarshalXMLAttr(attr xml.Attr) error {
	var r, g, b, a uint8
	_, err := fmt.Sscanf(attr.Value, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	if err != nil {
		return err
	}
	c.RGBA = color.RGBA{r, g, b, a}
	return nil
}

func loadFont(fontName string) *truetype.Font {
	// Read font file from disk
	_, callerFile, _, _ := runtime.Caller(0)
	callerDir := filepath.Dir(callerFile)
	filePath := filepath.Join(callerDir, fmt.Sprintf("./fonts/%s.ttf", fontName))
	fontBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read font file: %v", err)
	}
// Parse font file into a truetype.Font 
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Failed to parse font file: %v", err)
	}
	return font
}


func fetchImageFromPath(path string) image.Image {
	if contents, ok := imageCache[path]; ok {
		return contents
	} 
	// else fetch the file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open the image file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode the image: %v", err)
	}	

	imageCache[path] = img


	return img
}

func resizeImage(img image.Image, sizex uint, sizey uint) image.Image {
	return resize.Resize(sizex, sizey, img, resize.Lanczos3)
}
