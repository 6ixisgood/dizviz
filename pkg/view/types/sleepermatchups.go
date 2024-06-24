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
	
			<template dir="col" justify="center" align="center" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">{{ .League.Name }}</rainbow-text>
			</template>
	
			{{ else if gt .Phase 0 }}
			<template justify="space-between" align="center" dir="col" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">

				<!-- Team Headers -->
				<container sizeX="100%" sizeY="35%">
					<template sizeX="100%" sizeY="100%">
						<container sizeX="50%" sizeY="100%">
							<template justify="space-around" align="center" sizeX="100%" sizeY="100%" dir="col">
								<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team1.Name }}</text>
								<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ .Team1.Avatar }}"></image>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16"> {{ .Team1.Score }}</text>
							</template>
						</container>
						<container sizeX="50%" sizeY="100%">
							<template justify="space-around" align="center" sizeX="100%" sizeY="100%" dir="col">
								<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team2.Name }}</text>
								<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ .Team2.Avatar }}"></image>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16"> {{ .Team2.Score }}</text>
							</template>
						</container>
					</template>
				</container>


				<!-- Player Info -->
				<container sizeX="100%" sizeY="65%">
					<template sizeX="100%" sizeY="100%">
						<container sizeX="45%" sizeY="100%">
							<template dir="col" sizeX="100%" sizeY="100%" justify="space-around">
								{{ if eq .Phase 1 }}	
									{{ range $index, $element := .Team1.Starters }}
									<container sizeX="100%" sizeY="10%">
										<template sizeX="100%" sizeY="100%" justify="space-between">
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
										</template>
									</container>
									{{ end }}
								{{ else if eq .Phase 2 }}
									{{ range $index, $element := .Team1.Bench }}
									<container sizeX="100%" sizeY="10%">
										<template sizeX="100%" sizeY="100%" justify="space-between">
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $BenchedColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
										</template>
									</container>
									{{ end }}
								{{ end }}
							</template>
						</container>

						<container sizeX="10%" sizeY="100%">
							<template dir="col" sizeX="100%" sizeY="100%" justify="space-around">
								{{ if eq .Phase 1 }}
									{{ range $index, $element := .League.StartingPositions }}
									<container sizeX="100%" sizeY="10%">
										<template sizeX="100%" sizeY="100%" justify="center">
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $PositionColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element }}</text>
										</template>
									</container>
									{{ end }}
								{{ end}}
							</template>
						</container>

						<container sizeX="45%" sizeY="100%">
							<template dir="col" sizeX="100%" sizeY="100%" justify="space-around">
								{{ if eq .Phase 1 }}	
									{{ range $index, $element := .Team2.Starters }}
									<container sizeX="100%" sizeY="10%">
										<template sizeX="100%" sizeY="100%" justify="space-between">
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
										</template>
									</container>
									{{ end }}
								{{ else if eq .Phase 2 }}
									{{ range $index, $element := .Team2.Bench }}
									<container sizeX="100%" sizeY="10%">
										<template sizeX="100%" sizeY="100%" justify="space-between">
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $BenchedColor }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
											<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
										</template>
									</container>
									{{ end }}
								{{ end }}
							</template>
						</container>	
					</template>
				</container>
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
