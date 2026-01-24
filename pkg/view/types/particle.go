package types

import (
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type ParticlesView struct {
	c.BaseView
}

func (v *ParticlesView) TemplateString() string {
	return `
		<template size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">
			<gravity-particles size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}"></gravity-particles>
		 </template>
	`
}

func init() {
	c.RegisterView("particle", func() c.View { return &ParticlesView{} })
}
