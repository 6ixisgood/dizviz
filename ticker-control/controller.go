package main

import (
	"log"
	"image"
	"time"
	comp "github.com/sixisgoood/matrix-ticker/components"
)

var (
	animation = &Animation{
	}
)

type Animation struct {
	view		comp.View
	template	comp.Template
}

func (a *Animation) Init(view string, viewConfig map[string]string) {
	log.Printf("Initializing view in controller")

	// init new view in background
	newView := comp.RegisteredViews[view](viewConfig)
	newView.Init()
	comp.TemplateRefresh(newView)

	// stop the old view and switch to new view
	if a.view != nil {
		a.view.Stop()
	}
	a.view = newView
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	t := a.view.Template()

	// render template to image.Image
	for {
		if !(t.Ready()) {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	im := t.Render()
	return im, time.After(time.Millisecond * 10), nil
}

func GetAnimation() *Animation {
	return animation
}
