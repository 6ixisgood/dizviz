package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type TemplateView struct {
	c.BaseView

	templateString string
}

type TemplateViewConfig struct {
	Template string `json:"template" spec:"required:'true'"`
}

func TemplateViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*TemplateViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type TemplateViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	return &TemplateView{
		templateString: config.Template,
	}, nil
}

func (v *TemplateView) TemplateString() string {
	return v.templateString
}

func init() {
	c.RegisterView("template", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &TemplateViewConfig{} },
		NewView:   TemplateViewCreate,
	})
}
