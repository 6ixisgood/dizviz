package types

import (
	"errors"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"strconv"
	"time"
)

type SleeperMatchupsView struct {
	c.BaseView

	League        string
	Week          int
	SleeperClient *d.Sleeper
	matchups      [][]d.SleeperTeamFormatted
	matchIndex    int
	Phase         int
	dataRefresh   *util.Refresher
	phaseRefresh  *util.Refresher
	league        d.SleeperLeagueFormatted
}

type SleeperMatchupsViewConfig struct {
	LeagueID string `json:"league_id" spec:"required='true',label='League ID'"`
	Week     int    `json:"week" spec:"required='true',min='1',max='18',label='Week'"`
}

func SleeperMatchupsViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*SleeperMatchupsViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type SleeperMatchupsViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	client := d.SleeperClient()

	return &SleeperMatchupsView{
		League:        config.LeagueID,
		Week:          config.Week,
		SleeperClient: client,
		Phase:         0,
	}, nil
}

func (v *SleeperMatchupsView) Init() {
	// init ticker and stop chan
	v.dataRefresh = util.RefresherCreate(60*time.Second, v.RefreshData)
	v.phaseRefresh = util.RefresherCreate(5*time.Second, v.RefreshPhase)
	v.RefreshData()
	v.dataRefresh.Start()
	v.phaseRefresh.Start()
}

func (v *SleeperMatchupsView) RefreshData() {
	v.matchups = v.SleeperClient.GetMatchupsFormatted(v.League, strconv.Itoa(v.Week))
	v.league = v.SleeperClient.GetLeagueFormatted(v.League)

}

func (v *SleeperMatchupsView) RefreshPhase() {
	// change the phase at a given rate
	if v.Phase < 2 {
		v.Phase += 1
	} else if v.Phase == 2 {
		v.Phase = 1
		v.matchIndex = (v.matchIndex + 1) % len(v.matchups)
	}
	c.TemplateRefresh(v)
}

func (v *SleeperMatchupsView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Team1":  v.matchups[v.matchIndex][0],
		"Team2":  v.matchups[v.matchIndex][1],
		"League": v.league,
		"Phase":  v.Phase,
		"Week":   v.Week,
	}
}

func (v *SleeperMatchupsView) Stop() {
	v.dataRefresh.Stop()
	v.phaseRefresh.Stop()
}

func (v *SleeperMatchupsView) TemplateString() string {
	return `
			{{ $BenchedColor := "#FF4542FF" }}
			{{ $ScoreColor := "#F2FF00FF" }}
			{{ $PlayingColor := "#FFFFFFFF"}}
			{{ $TeamNameColor := "#66CCFFFF"}}
			{{ $PositionColor := "#5FE512FF"}}
	
			{{ if eq .Phase 0 }}
	
			<template dir="col" justify="center" align="center" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">{{ .League.Name }}</rainbow-text>
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">Week {{ .Week }}</rainbow-text>
			</template>
	
			{{ else if gt .Phase 0 }}
			<template justify="space-between" align="center" dir="col" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">

				<!-- Team Headers -->
				<template size-x="100%" size-y="35%">
					<template justify="space-around" align="center" size-x="50%" size-y="100%" dir="col">
						<text size-x="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team1.Name }}</text>
						<image size-x="{{ $DefaultImageSizex }}" size-y="{{ $DefaultImageSizey }}" src="{{ .Team1.Avatar }}"></image>
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16"> {{ .Team1.Score }}</text>
					</template>
					<template justify="space-around" align="center" size-x="50%" size-y="100%" dir="col">
						<text size-x="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team2.Name }}</text>
						<image size-x="{{ $DefaultImageSizex }}" size-y="{{ $DefaultImageSizey }}" src="{{ .Team2.Avatar }}"></image>
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16"> {{ .Team2.Score }}</text>
					</template>
				</template>


				<!-- Player Info -->
				<template size-x="100%" size-y="65%">
					<template dir="col" size-x="45%" size-y="100%" justify="space-around">
						{{ if eq .Phase 1 }}	
							{{ range $index, $element := .Team1.Starters }}
							<template size-x="100%" size-y="10%" justify="space-between">
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
							</template>
							{{ end }}
						{{ else if eq .Phase 2 }}
							{{ range $index, $element := .Team1.Bench }}
							<template size-x="100%" size-y="10%" justify="space-between">
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $BenchedColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
							</template>
							{{ end }}
						{{ end }}
					</template>

					<template dir="col" size-x="10%" size-y="100%" justify="space-around">
						{{ if eq .Phase 1 }}
							{{ range $index, $element := .League.StartingPositions }}
							<template size-x="100%" size-y="10%" justify="center">
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $PositionColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element }}</text>
							</template>
							{{ end }}
						{{ end}}
					</template>

					<template dir="col" size-x="45%" size-y="100%" justify="space-around">
						{{ if eq .Phase 1 }}	
							{{ range $index, $element := .Team2.Starters }}
							<template size-x="100%" size-y="10%" justify="space-between">
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
							</template>
							{{ end }}
						{{ else if eq .Phase 2 }}
							{{ range $index, $element := .Team2.Bench }}
							<template size-x="100%" size-y="10%" justify="space-between">
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $BenchedColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
							</template>
							{{ end }}
						{{ end }}
					</template>
				</template>
			 </template>

			 {{ end }}
	`
}

func init() {
	c.RegisterView("sleeper-matchups", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &SleeperMatchupsViewConfig{} },
		NewView:   SleeperMatchupsViewCreate,
	})
}
