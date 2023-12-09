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
	if a.view != nil {
		a.view.Stop()
	}

	// create a new view
	a.view = comp.RegisteredViews[view](viewConfig)
	a.view.Init()
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	t := a.view.Template()


	if t == nil {
		t = comp.ExecuteViewTemplate(a.view)
		a.view.SetTemplate(t)
		t.Init()
	}


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
