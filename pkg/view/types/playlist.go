package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type PlaylistView struct {
	c.BaseView

	// Config fields (same pattern as other views like ImagePlayerView)
	Views            []PlaylistViewConfigView `json:"views" spec:"label='Views',min='1'"`
	GlobalDuration   int                      `json:"global_duration,omitempty" spec:"label='Global Duration (seconds)',min='1',default='60'"`
	TransitionEffect string                   `json:"transition_effect,omitempty" spec:"label='Transition Effect',default='none'"`

	// Runtime fields (not part of config)
	views       []c.View
	activeIndex int
	ctx         context.Context
	cancel      context.CancelFunc
	ticker      *time.Ticker
}

const (
	DefaultPlayTime = 60 // seconds
)

// PlaylistViewConfigView represents a single view in a playlist.
// This struct supports two formats:
// 1. UI format (for creation/storage): { "viewId": "abc123", "duration": 60 }
// 2. Agent format (after expansion): { "duration": 60, ...view fields... }
// The control plane expands format 1 to format 2 before sending to agents.
// In the expanded format, all fields from the referenced view definition are merged in.
type PlaylistViewConfigView struct {
	// ViewId is used by UI/control plane to reference stored view definitions (format 1)
	ViewId string `json:"viewId,omitempty" spec:"label='View Definition ID'"`

	// Duration controls how long this view is shown (overrides global duration)
	Duration int `json:"duration,omitempty" spec:"label='Duration (seconds)',default='0'"`
}

func (v *PlaylistView) Init(configJSON string, ctx c.ViewContext) error {
	// Initialize base view and unmarshal config directly into v (same as other views)
	if err := c.InitViewFromJSON(v, configJSON, ctx); err != nil {
		return err
	}

	// Validate we have at least one view
	if len(v.Views) == 0 {
		return errors.New("playlist must contain at least one view")
	}

	// Set default global duration if not specified
	if v.GlobalDuration == 0 {
		v.GlobalDuration = DefaultPlayTime
	}

	// Create context for lifecycle management
	v.ctx, v.cancel = context.WithCancel(context.Background())

	// We need to re-parse the config to get the raw view data with type information
	// The Views field in v only captures viewId and duration, but the expanded config
	// has type and all child view fields merged in
	var rawConfig struct {
		Views []map[string]interface{} `json:"views"`
	}
	if err := json.Unmarshal([]byte(configJSON), &rawConfig); err != nil {
		return fmt.Errorf("failed to parse raw playlist config: %w", err)
	}

	// Initialize all child views
	v.views = make([]c.View, 0, len(rawConfig.Views))
	for i, viewData := range rawConfig.Views {
		// Extract type from the expanded view data
		viewType, ok := viewData["type"].(string)
		if !ok || viewType == "" {
			return fmt.Errorf("view at index %d has no type - playlist config must be expanded by control plane before reaching agent", i)
		}

		// Get factory for this view type
		factory, exists := c.RegisteredViews[viewType]
		if !exists {
			return fmt.Errorf("unknown view type '%s' at index %d", viewType, i)
		}

		// Create view instance
		childView := factory()

		// Build child config by removing playlist-specific fields
		childConfigMap := make(map[string]interface{})
		for key, value := range viewData {
			// Skip playlist-specific fields
			if key != "viewId" && key != "type" && key != "duration" {
				childConfigMap[key] = value
			}
		}

		// Marshal the child config to JSON
		childConfigJSON, err := json.Marshal(childConfigMap)
		if err != nil {
			return fmt.Errorf("failed to marshal config for view %d: %w", i, err)
		}

		// Initialize child view with same context as parent
		if err := childView.Init(string(childConfigJSON), ctx); err != nil {
			return fmt.Errorf("failed to initialize view %d (%s): %w", i, viewType, err)
		}

		v.views = append(v.views, childView)
	}

	// Start with first view
	v.activeIndex = 0

	// Refresh the first view's template to ensure it's properly rendered
	c.TemplateRefresh(v.views[0])
	v.SetTemplate(v.views[0].Template())

	// Start rotation ticker
	duration := v.getDuration(0)
	v.ticker = time.NewTicker(time.Duration(duration) * time.Second)
	go v.rotationLoop()

	log.Printf("Playlist view initialized with %d views", len(v.views))
	return nil
}

// getDuration returns the duration for a specific view index
func (v *PlaylistView) getDuration(index int) int {
	viewDuration := v.Views[index].Duration
	if viewDuration > 0 {
		return viewDuration
	}
	return v.GlobalDuration
}

// rotationLoop handles automatic view rotation
func (v *PlaylistView) rotationLoop() {
	for {
		select {
		case <-v.ctx.Done():
			return
		case <-v.ticker.C:
			v.NextView()
		}
	}
}

func (v *PlaylistView) NextView() {
	prevIndex := v.activeIndex
	nextIndex := (v.activeIndex + 1) % len(v.views)

	// Update active index
	v.activeIndex = nextIndex

	// Refresh the next view's template to ensure it's properly rendered
	c.TemplateRefresh(v.views[nextIndex])

	// Update template to next view
	v.SetTemplate(v.views[nextIndex].Template())

	// Trigger template refresh on the playlist view itself
	c.TemplateRefresh(v)

	// Update ticker for next view's duration
	duration := v.getDuration(nextIndex)
	v.ticker.Reset(time.Duration(duration) * time.Second)

	log.Printf("Playlist: switched from view %d to view %d (duration: %ds)",
		prevIndex, nextIndex, duration)
}

func (v *PlaylistView) TemplateString() string {
	if v.activeIndex < 0 || v.activeIndex >= len(v.views) {
		return ""
	}
	return v.views[v.activeIndex].TemplateString()
}

func (v *PlaylistView) TemplateData() map[string]interface{} {
	if v.activeIndex < 0 || v.activeIndex >= len(v.views) {
		return map[string]interface{}{}
	}
	return v.views[v.activeIndex].TemplateData()
}

func (v *PlaylistView) Stop() {
	if v.cancel != nil {
		v.cancel()
	}
	if v.ticker != nil {
		v.ticker.Stop()
	}
	// Stop all child views
	for _, view := range v.views {
		view.Stop()
	}
	log.Printf("Playlist view stopped")
}

func init() {
	c.RegisterView("playlist", func() c.View { return &PlaylistView{} })
}
