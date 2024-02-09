package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TextView struct {
	c.BaseView

	Text		string
}

type TextViewConfig struct {
	Text	string		`json:"text"`
}

func (vc *TextViewConfig) Validate() error {
	if vc.Text == "" {
		return errors.New("'text' field is required")
	}
	return nil
}

func TextViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*TextViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type TextViewConfig")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &TextView{
		Text: config.Text,
	}, nil
}

func (v *TextView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Text": v.Text,
	}  
}

func (v *TextView) TemplateString() string {
	return `
		<template dir="row" justify="space-between" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">One</text>
			<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">Two</text>
		</template>
	`
}

func init() {
	c.RegisterView("text", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &TextViewConfig{} },
		NewView: TextViewCreate,
	})
}