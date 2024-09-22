package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TextImageView struct {
	c.BaseView

	Text      string
	FontSize  int
	Alignment string
	Justify   string
	Color     string
	BgColor   string
	Src string

}

type TextImageViewConfig struct {
	Text      string `json:"text" spec:"required='true',min='1',label='Text'"`
	FontSize  int `json:"font-size" spec:"required='false',min='1',label='Font Size'"`
	Alignment string `json:"alignment" spec:"required='false',label='Alignment'"`
	Justify   string `json:"justify" spec:"required='false',label='Justify'"`
	Color     string `json:"color" spec:"required='false',label='Color'"`
	BgColor   string `json:"bg-color" spec:"required='false',label='Background Color'"`
	Src string `json:"src" spec:"required='true',label='Src (URL/Filepath)'"`

}

func TextImageViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*TextImageViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type TextImageViewConfig")
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

	if config.FontSize == 0 {
		config.FontSize = 12
	}

	return &TextImageView{
		Text:      config.Text,
		FontSize:  config.FontSize,
		Justify:   config.Justify,
		Alignment: config.Alignment,
		Color:     config.Color,
		BgColor:   config.BgColor,
		Src:	   config.Src,
	}, nil
}

func (v *TextImageView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Text":      v.Text,
		"FontSize":	 v.FontSize,
		"Justify":   v.Justify,
		"Alignment": v.Alignment,
		"Color":     v.Color,
		"BgColor":   v.BgColor,
		"Src":   	 v.Src,
	}
}

func (v *TextImageView) TemplateString() string {
	return `
		<template dir="col" justify="{{ .Justify }}" align="{{ .Alignment }}" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}" bg-color="{{ .BgColor }}">
			<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" bg-color="{{ .BgColor }}" color="{{ .Color }}" size="24">{{ .Text }}</text>
			<image size-x="75%" size-y="75%" src="{{ .Src }}" loop="true"></image>
		</template>
	`
}

func init() {
	c.RegisterView("text-image", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &TextImageViewConfig{} },
		NewView:   TextImageViewCreate,
	})
}
