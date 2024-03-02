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

	views       []c.View `ui:"label='Default Duration',type='number',min='1',hint='Enter default time in seconds'"`
	activeIndex int
	timings     []time.Duration
}

const (
	DefaultPlayTime = 60
)

type PlaylistViewConfig struct {
	Views []struct {
		c.ViewDefinition
		Settings struct {
			Time time.Duration `json:"time" spec:"label='Time (s)',type='number',required=true,min=10"`
		} `json:"settings" spec:"label='View Settings',type='section'"`
	} `json:"views" spec:"label='Views',type='list'"`
	Settings struct {
		Time time.Duration `json:"time" spec:"label='Default Duration',type='number',min='1',hint='Enter default time in seconds'"`
	} `json:"settings" spec:"label='Global Settings',type='section'"`
}

// type ViewConfigDesc struct {
// 	FieldName		string
// 	FieldType		string
// }

// {
//   "Views": {
//     "label": "Views",
//     "type": "list",
//     "fields": {
//       "Type": {
//         "label": "View Type",
//         "type": "text"
//       },
//       "Settings": {
//         "label": "View Settings",
//         "type": "section",
//         "fields": {
//           "Time": {
//             "label": "Time (s)",
//             "type": "number",
//             "min": "1",
//             "hint": "Duration in seconds"
//           }
//         }
//       }
//     }
//   },
//   "Settings": {
//     "label": "Global Settings",
//     "type": "section",
//     "fields": {
//       "Time": {
//         "label": "Default Duration",
//         "type": "number",
//         "min": "1",
//         "hint": "Enter default time in seconds"
//       }
//     }
//   }
// }

func PlaylistViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*PlaylistViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type PlaylistViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	// parse global settings
	var defaultTime time.Duration = DefaultPlayTime
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

		// go from the ViewConfig (map[string]interface{}) to []byte
		jsonConfig, err := json.Marshal(v.Config)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error marshaling generic ViewConfig to []byte"))
		}

		// go from []byte to specific ViewConfig type
		configInstance := regView.NewConfig()
		if err := json.Unmarshal(jsonConfig, &configInstance); err != nil {
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
