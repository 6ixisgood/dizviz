package main

import (
	"log"
	"image"
	"time"
	"encoding/xml"
	comp "github.com/sixisgoood/matrix-ticker/components"
)

var (
	animation = &Animation{
		stopChan: make(chan struct{}),
	}
)

type Animation struct {
	view		comp.View
	template	comp.Template
	ticker		*time.Ticker
	stopChan	chan struct{}
}


func (a *Animation) ContinuousRefresh(notify <-chan time.Time) () {
	a.Refresh()
	for {
		select {
		// wait for the time period to pass to refresh
		case <-notify:
			a.Refresh()
		// wait for the stop signal
		case <-a.stopChan:
			a.ticker.Stop()
			return
		}
	}
}

func (a *Animation) Refresh() {
	log.Printf("Refreshing the currently active view")
	// refresh the view's data
	a.view.Refresh()

	// grab the view's updated template string
	template := a.view.Template()

	// unmarshall the updated template to comp.Template
	var t comp.Template
	err := xml.Unmarshal([]byte(template), &t)
	if err != nil {
		log.Fatalf("Unable to unmarshal xml content: '%v'", err)
	}

	// init the new template
	t.Init()

	// set this as the Animation's new comp.Template
	a.template = t
}

func (a *Animation) Init(view string, viewConfig map[string]string) {
	// create a new view
	a.view = comp.RegisteredViews[view](viewConfig)

	// close the stop channel, stopping the refresh cycle if it's running
	close(a.stopChan)

	// initialize the new channel
	a.stopChan = make(chan struct{})

	a.ticker = time.NewTicker(10000000 * time.Millisecond)
	// start refreshing at given rate
	go a.ContinuousRefresh(a.ticker.C)
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	// render template to image.Image
	for {
		if !(a.template.Ready()) {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	im := a.template.Render()
	return im, time.After(time.Millisecond * 1), nil
}

func GetAnimation() *Animation {
	return animation
}
