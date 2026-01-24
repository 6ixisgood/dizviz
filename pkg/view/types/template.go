package types

import (
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TemplateView struct {
	c.BaseView

	RawTemplate string `json:"template" spec:"required='true',label='Raw Template'"`
}

func (v *TemplateView) TemplateString() string {
	return v.RawTemplate
}

func init() {
	c.RegisterView("template", func() c.View { return &TemplateView{} })
}
