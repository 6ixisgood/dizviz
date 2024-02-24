package types

import (
	"encoding/xml"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	"image"
	"image/color"
)

type PaddleBallVisualizer struct {
	c.BaseComponent

	XMLName      xml.Name `xml:"pong"`
	BallRadius   float64  `xml:"ballRadius,attr"`
	BallSpeedX   float64
	BallSpeedY   float64
	BallX        float64
	BallY        float64
	PaddleHeight float64 `xml:"paddleHeight,attr"`
	PaddleWidth  float64 `xml:"paddleWidth,attr"`
	LeftPaddleY  float64
	RightPaddleY float64
	Amplitude    float64   // this will be updated by an external function based on the music beat
	Color        util.RGBA `xml:"color,attr"`
}

func (pbv *PaddleBallVisualizer) Init() {
	pbv.BaseComponent.Init()

	// Set initial values
	pbv.BallSpeedX = 1.0
	pbv.BallSpeedY = 0.0
	pbv.LeftPaddleY = float64(pbv.Height())/2 - pbv.PaddleHeight/2
	pbv.RightPaddleY = pbv.LeftPaddleY

	pbv.BallX = float64(pbv.Width()) / 2
	pbv.BallY = float64(pbv.Height()) / 2
}

func (pbv *PaddleBallVisualizer) Render() image.Image {
	pbv.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	pbv.Ctx.Clear()
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
	pbv.Ctx.SetColor(pbv.Color.RGBA)

	// Draw the ball
	pbv.Ctx.DrawCircle(pbv.BallX, pbv.BallY, pbv.BallRadius)
	pbv.Ctx.Fill()

	// Draw the paddles
	pbv.Ctx.DrawRectangle(0, pbv.LeftPaddleY, pbv.PaddleWidth, pbv.PaddleHeight)
	pbv.Ctx.DrawRectangle(float64(pbv.Width())-pbv.PaddleWidth, pbv.RightPaddleY, pbv.PaddleWidth, pbv.PaddleHeight)
	pbv.Ctx.Fill()

	return pbv.Ctx.Image()
}

func init() {
	c.RegisterComponent("paddleball", func() c.Component { return &PaddleBallVisualizer{} })
}
