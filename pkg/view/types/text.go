package types

import (
	c "github.com/6ixisgood/disco/pkg/view/common"
)

type TextView struct {
	c.BaseView

	Text      string `json:"text" spec:"required='true',min='1',label='Text'"`
	Alignment string `json:"alignment" spec:"label='Alignment'"`
	Justify   string `json:"justify" spec:"label='Justify'"`
	Color     string `json:"color" spec:"label='Color'"`
	BgColor   string `json:"bg_color" spec:"label='Background Color'"`
	FontSize  int    `json:"font_size" spec:"label='Font Size'"`
	FontType  string `json:"font_type" spec:"label='Font Type'"`
	FontStyle string `json:"font_style" spec:"label='Font Style'"`
	Rainbow   bool   `json:"raindow" spec:"label='Rainbow?'"`
}

func (v *TextView) Init(configJSON string, ctx c.ViewContext) error {
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
		v.FontSize = v.BaseView.DefaultFontSize
	}

	if v.FontType == "" {
		v.FontType = v.BaseView.DefaultFontType
	}
	if v.FontStyle == "" {
		v.FontStyle = v.BaseView.DefaultFontStyle
	}

	return nil
}

func (v *TextView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Text":      v.Text,
		"Justify":   v.Justify,
		"Alignment": v.Alignment,
		"Color":     v.Color,
		"BgColor":   v.BgColor,
		"FontSize":  v.FontSize,
		"FontType":  v.FontType,
		"FontStyle": v.FontStyle,
		"Rainbow":   v.Rainbow,
	}
}

func (v *TextView) TemplateString() string {
	return `
		<template dir="col" justify="{{ .Justify }}" align="{{ .Alignment }}" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}" bg-color="{{ .BgColor }}" overflow-x="scroll-ease" scroll-speed="15">
			<text rainbow="{{ .Rainbow }}" font="{{ .FontType }}" style="{{ .FontStyle }}" color="{{ .Color }}" size="{{ .FontSize }}">{{ .Text }}</text>
		</template>
	`
}

func init() {
	c.RegisterView("text", func() c.View { return &TextView{} })
}
