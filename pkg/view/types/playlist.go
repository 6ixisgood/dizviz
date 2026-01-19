package types

import (
	"context"
	"errors"
	"time"

	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type PlaylistView struct {
	c.BaseView

	views       []c.View
	activeIndex int
	timings     []time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
}

const (
	DefaultPlayTime = 60
)

type PlaylistViewConfigView struct {
	ViewId   string                     `json:"viewId" spec:"label='The View ID in the store'"`
	Settings PlaylistViewConfigSettings `json:"settings" spec:"label='Settings'"`
}

type PlaylistViewConfigSettings struct {
	Time time.Duration `json:"time" spec:"label='Duration(s)'"`
}

type PlaylistViewConfig struct {
	Views    []PlaylistViewConfigView   `json:"views" spec:"label='Views',min='1'"`
	Settings PlaylistViewConfigSettings `json:"settings" spec:"label='Global Settings'"`
}

func PlaylistViewConfigCreate() c.ViewConfig {
	return &PlaylistViewConfig{
		Views:    make([]PlaylistViewConfigView, 0),
		Settings: PlaylistViewConfigSettings{},
	}
}

func PlaylistViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*PlaylistViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type PlaylistViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	// parse global settings
	// var defaultTime time.Duration = DefaultPlayTime
	// if config.Settings.Time > 0 {
	// 	defaultTime = config.Settings.Time
	// }

	var views []c.View
	var timings []time.Duration

	// BROKEN: Playlist views are currently non-functional after control plane store refactor
	// This view type calls GetViewDefinition() which has been removed.
	// GetViewDefinition was deleted because agents no longer have store access.
	// The store is now owned exclusively by the control plane.
	//
	// To fix this, we need to redesign playlist views with one of these approaches:
	// 1. Control plane expands playlist definitions before sending to agent
	//    - When assigning playlist, control plane resolves all child view IDs
	//    - Sends a modified playlist config with full child view configs embedded
	// 2. Change playlist config to contain full child view definitions, not just IDs
	//    - Playlist config would be much larger (nested view definitions)
	//    - No need for store lookups
	// 3. Add a gRPC endpoint for agents to request view definitions
	//    - Agent calls control plane when it needs a child view definition
	//    - Adds network overhead but keeps playlist simple
	//
	// For now, attempting to create a playlist view will fail at this point.
	// TODO: Implement one of the above solutions

	// for _, v := range config.Views {
	// 	// This will fail - GetViewDefinition no longer exists
	// 	viewDef, err := c.GetViewDefinition(v.ViewId)
	// 	if err != nil {
	// 		return nil, errors.New(fmt.Sprintf("Error fetching view definition from store"))
	// 	}

	// 	regView, exists := c.RegisteredViews[viewDef.Type]
	// 	if !exists {
	// 		return nil, errors.New(fmt.Sprintf("View type %s does not exist", viewDef.Type))
	// 	}

	// 	// go from the ViewConfig (map[string]interface{}) to []byte
	// 	jsonConfig, err := json.Marshal(viewDef.Config)
	// 	if err != nil {
	// 		return nil, errors.New(fmt.Sprintf("Error marshaling generic ViewConfig to []byte"))
	// 	}

	// 	// go from []byte to specific ViewConfig type
	// 	configInstance := regView.NewConfig()
	// 	if err := json.Unmarshal(jsonConfig, &configInstance); err != nil {
	// 		return nil, errors.New(fmt.Sprintf("Config for view type %s is invalid", viewDef.Type))
	// 	}

	// 	newView, err := regView.NewView(configInstance)
	// 	if err != nil {
	// 		return nil, errors.New(fmt.Sprintf("Failed to create view of type %s with given config\nError: %s", viewDef.Type, err))
	// 	}

	// 	views = append(views, newView)
	// 	time := defaultTime
	// 	if v.Settings.Time > 0 {
	// 		time = v.Settings.Time
	// 	}
	// 	timings = append(timings, time)

	// }

	// if len(views) == 0 {
	// 	return nil, errors.New("No views supplied in playlist config")
	// }

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
	select {
	case <-v.ctx.Done():
		return
	default:
		prevIndex := v.activeIndex
		nextIndex := (v.activeIndex + 1) % len(v.views)

		// Propagate context to child view before Init
		childView := v.views[nextIndex]
		if ctx := v.GetContext(); ctx != nil {
			childView.SetContext(ctx)
		}
		childView.Init()

		// set next view as active
		v.SetTemplate(childView.Template())
		v.activeIndex = nextIndex

		c.TemplateRefresh(v)
		// stop active view
		if prevIndex >= 0 {
			v.views[prevIndex].Stop()
		}

		// wait for next view
		go func() {
			time.Sleep(v.timings[v.activeIndex] * time.Second)
			v.NextView()
		}()
	}
}

func (v *PlaylistView) Stop() {
	v.cancel()
}

func (v *PlaylistView) Init() {
	v.BaseView.Init()
	v.ctx, v.cancel = context.WithCancel(context.Background())

	v.NextView()
}

func init() {
	c.RegisterView("playlist", c.RegisteredView{
		NewConfig: PlaylistViewConfigCreate,
		NewView:   PlaylistViewCreate,
	})
}
