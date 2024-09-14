package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type ImagePlayerView struct {
	c.BaseView

	Src string
}

type ImagePlayerViewConfig struct {
	Src string `json:"src" spec:"required='true',label='Src (URL/Filepath)'"`
}

func ImagePlayerViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*ImagePlayerViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type ImagePlayerViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	return &ImagePlayerView{
		Src: config.Src,
	}, nil
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
	c.RegisterView("image", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &ImagePlayerViewConfig{} },
		NewView:   ImagePlayerViewCreate,
	})
}
