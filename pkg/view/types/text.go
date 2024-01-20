package types

import (
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TextView struct {
	c.BaseView

	Text		string
}


func TextViewCreate(config map[string]string) c.View {
	return &TextView{
		Text: config["text"],
	}
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
	c.RegisterView("text", TextViewCreate)
}