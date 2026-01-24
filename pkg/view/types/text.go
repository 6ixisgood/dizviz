package types

import (
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TextView struct {
	c.BaseView

	Text      string `json:"text" spec:"required='true',min='1',label='Text'"`
	Alignment string `json:"alignment" spec:"label='Alignment'"`
	Justify   string `json:"justify" spec:"label='Justify'"`
	Color     string `json:"color" spec:"label='Color'"`
	BgColor   string `json:"bg_color" spec:"label='Background Color'"`
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

	return nil
}

func (v *TextView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Text":      v.Text,
		"Justify":   v.Justify,
		"Alignment": v.Alignment,
		"Color":     v.Color,
		"BgColor":   v.BgColor,
	}
}

func (v *TextView) TemplateString() string {
	return `
		<template dir="col" justify="{{ .Justify }}" align="{{ .Alignment }}" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}" bg-color="{{ .BgColor }}" overflow-x="scroll-ease" scroll-speed="15">
			<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ .Color }}" size="{{ $DefaultFontSize }}">{{ .Text }}</text>
		</template>
	`
}

func init() {
	c.RegisterView("text", func() c.View { return &TextView{} })
}
