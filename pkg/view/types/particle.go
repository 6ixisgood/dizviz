package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type ParticlesView struct {
	c.BaseView
}

type ParticlesViewConfig struct{}

func ParticlesViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*ParticlesViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type ParticlesViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	return &ParticlesView{}, nil
}

func (v *ParticlesView) TemplateString() string {
	return `
		<template size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">
			<gravity-particles size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}"></gravity-particles>
		 </template>
	`
}

func init() {
	c.RegisterView("particle", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &ParticlesViewConfig{} },
		NewView:   ParticlesViewCreate,
	})
}
