package types

import (
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TextImageView struct {
	c.BaseView

	Text      string `json:"text" spec:"required='true',min='1',label='Text'"`
	FontSize  int    `json:"font_size" spec:"min='1',label='Font Size'"`
	Alignment string `json:"alignment" spec:"label='Alignment'"`
	Justify   string `json:"justify" spec:"label='Justify'"`
	Color     string `json:"color" spec:"label='Color'"`
	BgColor   string `json:"bg_color" spec:"label='Background Color'"`
	Src       string `json:"src" spec:"required='true',label='Image Source (URL/Filepath)'"`
}

func (v *TextImageView) Init(configJSON string, ctx c.ViewContext) error {
	if err := c.InitViewFromJSON(v, configJSON, ctx); err != nil {
		return err
	}

	// Set defaults
	if v.Color == "" {
		v.Color = "#FFFFFFFF"
	}
	if v.BgColor == "" {
		v.BgColor = "#00000FF"
	}
	if v.FontSize == 0 {
		v.FontSize = 12
	}

	return nil
}

func (v *TextImageView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Text":      v.Text,
		"FontSize":  v.FontSize,
		"Justify":   v.Justify,
		"Alignment": v.Alignment,
		"Color":     v.Color,
		"BgColor":   v.BgColor,
		"Src":       v.Src,
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
	c.RegisterView("text-image", func() c.View { return &TextImageView{} })
}
