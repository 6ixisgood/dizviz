package types

import (
	"errors"
	"strconv"
	"time"

	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
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
	phaseDuration time.Duration
	dataDuration  time.Duration
}

type SleeperMatchupsViewConfig struct {
	LeagueID      string `json:"league_id" spec:"required='true',label='League ID'"`
	Week          int    `json:"week" spec:"required='true',min='1',max='18',label='Week'"`
	PhaseDuration int    `json:"phase_duration" spec:"required='false',label='Phase Duration'"`
	DataDuration  int    `json:"data_duration" spec:"required='false',label='Data Duration'"`
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

	if config.PhaseDuration == 0 {
		config.PhaseDuration = 15
	}

	if config.DataDuration == 0 {
		config.DataDuration = 60
	}

	return &SleeperMatchupsView{
		League:        config.LeagueID,
		Week:          config.Week,
		SleeperClient: client,
		Phase:         1,
		phaseDuration: time.Duration(config.PhaseDuration),
		dataDuration:  time.Duration(config.DataDuration),
	}, nil
}

func (v *SleeperMatchupsView) Init() {
	v.BaseView.Init()
	// init ticker and stop chan
	v.dataRefresh = util.RefresherCreate(v.dataDuration*time.Second, v.RefreshData)
	v.phaseRefresh = util.RefresherCreate(v.phaseDuration*time.Second, v.RefreshPhase)
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
			{{ $PlayerFontSize := 10 }}
			{{ $PlayerScoreFontSize := 9 }}
			{{ $TeamNameFontSize := 10 }}
	
			{{ if eq .Phase 0 }}
	
			<template dir="col" justify="center" align="center" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">{{ .League.Name }}</rainbow-text>
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">Week {{ .Week }}</rainbow-text>
			</template>
	
			{{ else if gt .Phase 0 }}
			<template justify="space-between" align="center" dir="col" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">

				<!-- Team Headers -->
				<template dir="col" size-x="100%" size-y="40%" justify="space-around" overflow-y="auto">
					<!-- Team Names Row -->
					<template dir="row" size-x="100%" size-y="30%" justify="space-around" align="center" overflow-y="auto">
						<template size-x="45%" size-y="100%" justify="center" align="center" overflow-y="auto">
							<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $TeamNameFontSize }}">{{ .Team1.Name }}</text>
						</template>
						<template size-x="10%" size-y="100%">
						</template>
						<template size-x="45%" size-y="100%" justify="center" align="center" overflow-y="auto">
							<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $TeamNameFontSize }}">{{ .Team2.Name }}</text>
						</template>
					</template>
					
					<!-- Team Avatars Row -->
					<template dir="row" size-x="100%" size-y="40%" justify="space-around" align="center">
						<template size-x="45%" size-y="100%" justify="center" align="center">
							{{ if .Team1.Avatar }}
								<image size-x="{{ $DefaultImageSizex }}" size-y="{{ $DefaultImageSizey }}" src="{{ .Team1.Avatar }}"></image>
							{{ end}}
						</template>
						<template size-x="10%" size-y="100%">
						</template>
						<template size-x="45%" size-y="100%" justify="center" align="center">
							{{ if .Team2.Avatar }}
								<image size-x="{{ $DefaultImageSizex }}" size-y="{{ $DefaultImageSizey }}" src="{{ .Team2.Avatar }}"></image>
							{{ end}}
						</template>
					</template>
					
					<!-- Team Scores Row -->
					<template dir="row" size-x="100%" size-y="30%" justify="space-around" align="center">
						<template size-x="45%" size-y="100%" justify="center" align="center">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16">{{ .Team1.Score }}</text>
						</template>
						<template size-x="10%" size-y="100%">
						</template>
						<template size-x="45%" size-y="100%" justify="center" align="center">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16">{{ .Team2.Score }}</text>
						</template>
					</template>
				</template>


				<!-- Player Info -->
				<template dir="col" size-x="100%" size-y="60%" overflow-y="scroll-bounce" scroll-speed="5">
					{{ if eq .Phase 1 }}
						{{ range $index, $element := .Team1.Starters }}
						<!-- Player Row {{ $index }} -->
						<template dir="row" size-x="100%" size-y="12%" justify="space-between" overflow-y="auto">
							<!-- Team1 Player -->
							<template size-x="45%" size-y="100%" justify="space-between" overflow-y="auto">
								<template size-x="70%" size-y="100%" justify="start" overflow-y="auto">
									<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $PlayerFontSize }}">{{ printf "%s" $element.Name }}</text>
								</template>
								<text size-x="30%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="{{ $PlayerScoreFontSize }}">{{ printf "%.2f" $element.Points }}</text>
							</template>
							
							<!-- Position -->
							<template size-x="10%" size-y="100%" justify="center">
								{{ if lt $index (len $.League.StartingPositions) }}
									<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $PositionColor }}" size="{{ $PlayerScoreFontSize }}">{{ index $.League.StartingPositions $index }}</text>
								{{ else }}
									<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $PositionColor }}" size="{{ $PlayerScoreFontSize }}">--</text>
								{{ end }}
							</template>
							
							<!-- Team2 Player -->
							<template size-x="45%" size-y="100%" justify="space-between" overflow-y="auto">
								{{ if lt $index (len $.Team2.Starters) }}
									<text size-x="30%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="{{ $PlayerScoreFontSize }}">{{ printf "%.2f" (index $.Team2.Starters $index).Points }}</text>
									<template size-x="70%" size-y="100%" justify="end" overflow-y="auto">
										<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $PlayerFontSize }}">{{ printf "%s" (index $.Team2.Starters $index).Name }}</text>
									</template>
								{{ else }}
									<text size-x="30%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="{{ $PlayerScoreFontSize }}">--</text>
									<template size-x="70%" size-y="100%" justify="end" overflow-y="auto">
										<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $PlayerFontSize }}">No Player</text>
									</template>
								{{ end }}
							</template>
						</template>
						{{ end }}
					{{ else if eq .Phase 2 }}
						{{ range $index, $element := .Team1.Bench }}
						<!-- Bench Row {{ $index }} -->
						<template dir="row" size-x="100%" size-y="12%" justify="space-between" overflow-y="auto">
							<!-- Team1 Bench Player -->
							<template size-x="45%" size-y="100%" justify="space-between">
								<template size-x="70%" size-y="100%" justify="start" overflow-y="auto">
									<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $PlayerFontSize }}">{{ printf "%s" $element.Name }}</text>
								</template>
								<text size-x="30%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $BenchedColor }}" size="{{ $PlayerScoreFontSize }}">{{ printf "%.2f" $element.Points }}</text>
							</template>
							
							<!-- Empty Position Column for Bench -->
							<template size-x="10%" size-y="100%" justify="center">
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $PositionColor }}" size="{{ $PlayerScoreFontSize }}">BN</text>
							</template>
							
							<!-- Team2 Bench Player -->
							<template size-x="45%" size-y="100%" justify="space-between">
								{{ if lt $index (len $.Team2.Bench) }}
									<text size-x="30%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $BenchedColor }}" size="{{ $PlayerScoreFontSize }}">{{ printf "%.2f" (index $.Team2.Bench $index).Points }}</text>
									<template size-x="70%" size-y="100%" justify="end" overflow-y="auto">
										<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $PlayerFontSize }}">{{ printf "%s" (index $.Team2.Bench $index).Name }}</text>
									</template>
								{{ else }}
									<text size-x="30%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $BenchedColor }}" size="{{ $PlayerScoreFontSize }}">--</text>
									<template size-x="70%" size-y="100%" justify="end" overflow-y="auto">
										<text word-wrap="true" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $PlayerFontSize }}">No Player</text>
									</template>
								{{ end }}
							</template>
						</template>
						{{ end }}
					{{ end }}
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
