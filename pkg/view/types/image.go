package types

import (
	c "github.com/6ixisgood/disco/pkg/view/common"
)

type ImagePlayerView struct {
	c.BaseView

	Src string `json:"src" spec:"required='true',min='1',label='Image Source (URL/Filepath)'"`
}

func (v *ImagePlayerView) Init(configJSON string, ctx c.ViewContext) error {
	return c.InitViewFromJSON(v, configJSON, ctx)
}

func (v *ImagePlayerView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Src": v.Src,
	}
}

func (v *ImagePlayerView) TemplateString() string {
	return `
		<template size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">
			<image size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}" src="{{ .Src }}" loop="true"></image>
		 </template>
	`
}

func init() {
	c.RegisterView("image", func() c.View { return &ImagePlayerView{} })
}
