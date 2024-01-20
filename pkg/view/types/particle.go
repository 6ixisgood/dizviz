package types

import (
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type ParticlesView struct {
	c.BaseView
}

func ParticlesViewCreate(config map[string]string) c.View {
	return &ParticlesView{}
}

func (v *ParticlesView) TemplateString() string {
	return `
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<gravity-particles sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}"></gravity-particles>
		 </template>
	`
}

func init() {
	c.RegisterView("particle", ParticlesViewCreate)
}