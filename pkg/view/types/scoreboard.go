package types

import (
	"errors"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"time"
)

type NFLBoxView struct {
	c.BaseView

	Auto             bool
	Matchup          string
	Date             time.Time
	SportsFeedClient *d.SportsFeed
	Game             d.NFLBoxScoreResponseFormatted
	Games            []d.NFLBoxScoreResponseFormatted
	Duration         time.Duration
	gameIndex        int
	Layout           string
	refresh          *util.Refresher
}

type NFLBoxViewConfig struct {
	Auto     bool      `json:"auto" spec:"required='true',label='Auto'"`
	Matchup  string    `json:"matchup" spec:"required='false',label='Matchup'"`
	Date     util.Date `json:"date" spec:"required='false',label='Date'"`
	Duration int       `json:"duration" spec:"required='false',label='Duration'"`
}

func NFLBoxViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*NFLBoxViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type NFLBoxViewConfig")
	}

	if err := c.ValidateViewConfig(config); err != nil {
		return nil, err
	}

	client := d.SportsFeedClient()

	var d time.Time
	if config.Auto {
		d = time.Now()
	} else {
		d = config.Date.Time
	}

	if config.Duration == 0 {
		if config.Auto {
			config.Duration = 15
		} else {
			config.Duration = 60
		}
	}

	return &NFLBoxView{
		Auto:             config.Auto,
		Matchup:          config.Matchup,
		Date:             d,
		Duration:         time.Duration(config.Duration),
		SportsFeedClient: client,
	}, nil
}

func (v *NFLBoxView) Init() {
	var f func()
	if v.Auto {
		// we're grabbing the current games and looping through them
		week := v.SportsFeedClient.FetchNFLCurrentWeek()
		v.Games = v.SportsFeedClient.FetchNFLWeeklyGamesFormatted(week.SeasonSlug, week.Week)
		v.gameIndex = -1
		f = v.RefreshPhase
	} else {
		// a specific game to be on
		f = v.RefreshGame
	}
	v.refresh = util.RefresherCreate(v.Duration*time.Second, f)
	f()
	v.refresh.Start()
}

func (v *NFLBoxView) RefreshGame() {
	// fetch the games
	v.Game, _ = v.SportsFeedClient.FetchNFLBoxScore(v.Matchup, v.Date)
	c.TemplateRefresh(v)
}

func (v *NFLBoxView) RefreshPhase() {
	v.gameIndex = (v.gameIndex + 1) % len(v.Games)
	var err error
	v.Game, err = v.SportsFeedClient.FetchNFLBoxScore(v.Games[v.gameIndex].GameID, v.Date)
	if err != nil {
		// skip this game
		v.RefreshPhase()
	}
	c.TemplateRefresh(v)
}

func (v *NFLBoxView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Game": v.Game,
	}
}

func (v *NFLBoxView) TemplateString() string {
	return `
		{{ $ScoreFontSize := 32 }}
		{{ $DetailFontSize := 14}}
		{{ $LogoSize := 64 }}

		<template dir="col" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">

			<template justify="space-around" size-x="100%" size-y="50%">
				<template dir="col" justify="space-between" align="center"  size-x="40%" size-y="100%">
		    		<image size-x="{{ $LogoSize }}" size-y="{{ $LogoSize }}" src="{{ .Game.AwayLogo }}"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">({{ .Game.AwayWins }}-{{ .Game.AwayLosses }}-{{ .Game.AwayTies }})</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.AwayScore }}</text>
				</template>

				<template dir="col" justify="space-around" align="center" size-x="20%" size-y="100%">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ .Game.QuarterMinRemaining }}:{{ .Game.QuarterSecRemaining }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ CardinalToOrdinal .Game.Quarter}}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ CardinalToOrdinal .Game.Down}}{{ "&" }}{{ .Game.YardsRemaining }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ .Game.LineOfScrimmage }}</text>
				</template>

				<template dir="col" justify="space-between"  align="center" size-x="40%" size-y="100%">
		    		<image size-x="{{ $LogoSize }}" size-y="{{ $LogoSize }}" src="{{ .Game.HomeLogo }}"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">({{ .Game.HomeWins }}-{{ .Game.HomeLosses }}-{{ .Game.HomeTies }})</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.HomeScore }}</text>
				</template>
			</template>

			<template justify="space-between" size-x="100%" size-y="50%">
				<template size-x="45%" size-y="100%" dir="col" justify="space-around">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Passing: {{ .Game.AwayPassYards }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Rushing: {{ .Game.AwayRushYards }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Sacks: {{ .Game.AwaySacks }}</text>
				</template>

				<template size-x="45%" size-y="100%" dir="col" justify="space-around">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Passing: {{ .Game.HomePassYards }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Rushing: {{ .Game.HomeRushYards }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Sacks: {{ .Game.HomeSacks }}</text>
				</template>
			</template>

		 </template>
		`
}

func init() {
	c.RegisterView("nflbox", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &NFLBoxViewConfig{} },
		NewView:   NFLBoxViewCreate,
	})
}
