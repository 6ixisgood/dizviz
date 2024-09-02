package types

import (
	"errors"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"time"
)

type SleeperDraftView struct {
	c.BaseView

	LeagueID        string
	league        d.SleeperLeagueFormatted
	DraftID			string
	Season string
	SleeperClient *d.Sleeper
	draft      d.SleeperDraftFormatted
	dataRefresh   *util.Refresher
}

type SleeperDraftViewConfig struct {
	LeagueID string `json:"league_id" spec:"required='true',label='League ID'"`
	DraftID string `json:"draft_id" spec:"required='true',label='Draft ID'"`
	Season     string    `json:"season" spec:"required='true',label='Season'"`
}

func SleeperDraftViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*SleeperDraftViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type SleeperDraftViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	client := d.SleeperClient()

	return &SleeperDraftView{
		LeagueID:        config.LeagueID,
		DraftID:          config.DraftID,
		Season:		config.Season,
		SleeperClient: client,
	}, nil
}

func (v *SleeperDraftView) Init() {
	// init data
	v.league = v.SleeperClient.GetLeagueFormatted(v.LeagueID)
	// init ticker and stop chan
	v.dataRefresh = util.RefresherCreate(60*time.Second, v.RefreshData)
	v.RefreshData()
	v.dataRefresh.Start()
}

func (v *SleeperDraftView) RefreshData() {
	v.draft = v.SleeperClient.GetDraftFormatted(v.LeagueID, v.DraftID)
}

func (v *SleeperDraftView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"League": v.league,
		"Season":	  v.Season,
		"Draft": v.draft,
	}
}

func (v *SleeperDraftView) Stop() {
	v.dataRefresh.Stop()
}

func (v *SleeperDraftView) TemplateString() string {
	return `
			{{ $BenchedColor := "#FF4542FF" }}
			{{ $ScoreColor := "#F2FF00FF" }}
			{{ $PlayingColor := "#FFFFFFFF"}}
			{{ $TeamNameColor := "#66CCFFFF"}}
			{{ $PositionColor := "#5FE512FF"}}

			<template justify="space-between" align="center" dir="col" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">

				<!-- Team Headers -->
				<template sizeX="100%" sizeY="100%">
					<template justify="space-around" align="center" sizeX="100%" sizeY="50%" dir="col">
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">On the Clock:</text>
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Draft.CurrentTeam.Name }}</text>
					</template>

					<template justify="space-around" align="center" sizeX="100%" sizeY="50%" dir="col">
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">Last Pick:</text>
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Draft.PrevTeam.Name }}</text>
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Draft.PrevPlayerPosition }}</text>
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Draft.PrevPlayerName }}</text>
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Draft.PrevPlayerTeam }}</text>
					</template>

				</template>

			</template>
	`
}

func init() {
	c.RegisterView("sleeper-draft", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &SleeperDraftViewConfig{} },
		NewView:   SleeperDraftViewCreate,
	})
}
