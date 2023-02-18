package main

import (
	"image"
	"image/color"
	"encoding/json"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/sixisgoood/go-rpi-rgb-led-matrix"
	"github.com/sixisgoood/matrix-ticker/animations"
)

type Request struct {
	AnimationType	string				`json:"type"`
	AnimationConfig map[string]string	`json:"config"`
}


type Font struct {
	Size		float64	`json:"size"`
	Type		string	`json:"type"`	
}

type TextScrollConfigRequest struct {
	Size		[]int		`json:"size"`	
	TextColor	[]uint8		`json:"textColor"`	
	BgColor		[]uint8		`json:"bgColor"`
	Direction	[]int		`json:"direction"`	
	Font		Font		`json:"font"`	
	Text		string		`json:"text"`	
}

type AnimationRequest struct {
	Type		string					`json:"type"`
	Config		TextScrollConfigRequest	`json:"config"`
}

func HandleRequest(req []byte) {
	log.Printf("Handle Request: %v", string(req))
	var animation rgbmatrix.Animation
	var animation_req AnimationRequest
	json.Unmarshal(req, &animation_req)

	if animation_req.Type == "textscroll" {
		animation = createTextScrollAnimation(animation_req.Config)
	} else {
		// do nothing
	}

	SetLiveAnimation(animation)	
}


func createTextScrollAnimation(req TextScrollConfigRequest) (*animations.TextScrollAnimation) {
	var font, _ = truetype.Parse(goregular.TTF)
	var face = truetype.NewFace(font, &truetype.Options{Size: req.Font.Size})

	config := animations.TextScrollConfig{
		Size: 			image.Point{req.Size[0], req.Size[1]},
		BgColor:		color.RGBA{req.BgColor[0], req.BgColor[1], req.BgColor[2], req.BgColor[3]},
		TextColor:		color.RGBA{req.TextColor[0], req.TextColor[1], req.TextColor[2], req.TextColor[3]},
		TextFontFace: 	face,
		Dir:			image.Point{req.Direction[0], req.Direction[1]},
		Text:			req.Text,
	}

	return animations.NewTextScrollAnimation(config)
}
