package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TextView struct {
	c.BaseView

	Text string
	Alignment string
	Justify string
	Color string
	BgColor string
}

type TextViewConfig struct {
	Text string `json:"text" spec:"required='true',min='1',label="Text"`
	Alignment string `json:"alignment" spec:"required='false',label="Alignment"`
	Justify string `json:"justify" spec:"required='false',label="Justify"`
	Color string `json:"color" spec:"required='false',label="Color"`
	BgColor string `json:"bg-color" spec:"required='false',label="Background Color"`
}

func TextViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*TextViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type TextViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	if config.Color == "" {
		config.Color = "#FFFFFFFF"
	}

	if config.BgColor == "" {
		config.BgColor = "#00000FF"
	}

	return &TextView{
		Text: config.Text,
		Justify: config.Justify,
		Alignment: config.Alignment,
		Color: config.Color,
		BgColor: config.BgColor,
	}, nil
}

func (v *TextView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Text": v.Text,
		"Justify": v.Justify,
		"Alignment": v.Alignment,
		"Color": v.Color,
		"BgColor": v.BgColor,
	}
}

func (v *TextView) TemplateString() string {
	return `
		<template dir="col" justify="{{ .Justify }}" align="{{ .Alignment }}" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}" bg-color="{{ .BgColor }}">
			<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ .Color }}" size="{{ $DefaultFontSize }}">{{ .Text }}</text>
		</template>
	`
}

func init() {
	c.RegisterView("text", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &TextViewConfig{} },
		NewView:   TextViewCreate,
	})
}
