package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type ImagePlayerView struct {
	c.BaseView

	Src		string
}

type ImagePlayerViewConfig struct {
	Src			string		`json:"src"`
}

func (vc *ImagePlayerViewConfig) Validate() error {
	if vc.Src == "" {
		return errors.New("'src' field is required")
	}
	return nil
}

func ImagePlayerViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*ImagePlayerViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type ImagePlayerViewConfig")
	}

	if err := config.Validate(); err != nil {
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
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<image sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}" src="{{ .Src }}" loop="true"></image>
		 </template>
	`
}

func init() {
	c.RegisterView("image", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &ImagePlayerViewConfig{} },
		NewView: ImagePlayerViewCreate,
	})
}