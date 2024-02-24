package types

import (
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"time"
)

type PlaylistView struct {
	c.BaseView

	views       []c.View
	activeIndex int
	timings     []time.Duration
}

type PlaylistViewConfig struct {
	Views []struct {
		Type     string          `json:"type"`
		Config   json.RawMessage `json:"config"`
		Settings struct {
			Time time.Duration `json:"time"`
		} `json:"settings"`
	} `json:"views"`
	Settings struct {
		Time time.Duration `json:"time"`
	} `json:"settings"`
}

func (vc *PlaylistViewConfig) Validate() error {
	if vc.Views == nil {
		return errors.New("'Views' field is required")
	}

	for i, v := range vc.Views {
		if v.Type == "" {
			return errors.New(fmt.Sprintf("'Type' field is required for view in position %d", i))
		}
	}

	return nil
}

func PlaylistViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*PlaylistViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type PlaylistViewConfig")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// parse global settings
	var defaultTime time.Duration = 10
	if config.Settings.Time > 0 {
		defaultTime = config.Settings.Time
	}

	var views []c.View
	var timings []time.Duration
	// validate and create each view
	for _, v := range config.Views {
		regView, exists := c.RegisteredViews[v.Type]
		if !exists {
			return nil, errors.New(fmt.Sprintf("View type %s does not exist", v.Type))
		}

		configInstance := regView.NewConfig()
		if err := json.Unmarshal(v.Config, &configInstance); err != nil {
			return nil, errors.New(fmt.Sprintf("Config for view type %s is invalid", v.Type))
		}

		newView, err := regView.NewView(configInstance)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to create view of type %s with given config\nError: %s", v.Type, err))
		}

		views = append(views, newView)
		time := defaultTime
		if v.Settings.Time > 0 {
			time = v.Settings.Time
		}
		timings = append(timings, time)

	}

	if len(views) == 0 {
		return nil, errors.New("No views supplied in playlist config")
	}

	return &PlaylistView{
		views:       views,
		timings:     timings,
		activeIndex: -1,
	}, nil
}

func (v *PlaylistView) TemplateString() string {
	return v.views[v.activeIndex].TemplateString()
}

func (v *PlaylistView) TemplateData() map[string]interface{} {
	return v.views[v.activeIndex].TemplateData()
}

func (v *PlaylistView) NextView() {
	// init next view
	fmt.Println(v.activeIndex)
	nextIndex := (v.activeIndex + 1) % len(v.views)

	// stop active view
	if v.activeIndex > 0 {
		v.views[v.activeIndex].Stop()
	}
	v.views[nextIndex].Init()

	// set next view as active
	v.activeIndex = nextIndex
	c.TemplateRefresh(v)

	// wait for next view
	go func() {
		time.Sleep(v.timings[v.activeIndex] * time.Second)
		v.NextView()
	}()
}

func (v *PlaylistView) Init() {
	v.NextView()
}

func init() {
	c.RegisterView("playlist", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &PlaylistViewConfig{} },
		NewView:   PlaylistViewCreate,
	})
}
