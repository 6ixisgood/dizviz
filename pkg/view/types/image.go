package types

import (
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type ImagePlayerView struct {
	c.BaseView

	Src		string
}


func ImagePlayerViewCreate(config map[string]string) c.View {
	return &ImagePlayerView{
		Src: config["src"],
	}
}

func (v *ImagePlayerView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Src": v.Src,
	}  
}

func (v *ImagePlayerView) TemplateString() string {
	return `
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<image sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}" src="{{ .Src }}" loop="true"></image>
		 </template>
	`
}

func init() {
	c.RegisterView("image", ImagePlayerViewCreate)
}