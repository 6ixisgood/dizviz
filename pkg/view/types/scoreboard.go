package types

import (
	"errors"
	"time"

	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
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
			config.Duration = 30
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
	v.BaseView.Init()
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

func (v *NFLBoxView) Stop() {
	v.refresh.Stop()
}

func (v *NFLBoxView) TemplateData() map[string]interface{} {
	// Create a list of stats for the comparison bars
	stats := []map[string]interface{}{
		// Offensive Stats
		{
			"Label":     "Passing Yards",
			"AwayValue": v.Game.AwayPassYards,
			"HomeValue": v.Game.HomePassYards,
		},
		{
			"Label":     "Rushing Yards",
			"AwayValue": v.Game.AwayRushYards,
			"HomeValue": v.Game.HomeRushYards,
		},
		{
			"Label":     "Passing TDs",
			"AwayValue": v.Game.AwayPassTD,
			"HomeValue": v.Game.HomePassTD,
		},
		{
			"Label":     "Rushing TDs",
			"AwayValue": v.Game.AwayRushTD,
			"HomeValue": v.Game.HomeRushTD,
		},
		{
			"Label":     "Completions",
			"AwayValue": v.Game.AwayPassCompletions,
			"HomeValue": v.Game.HomePassCompletions,
		},
		{
			"Label":     "Pass Attempts",
			"AwayValue": v.Game.AwayPassAttempts,
			"HomeValue": v.Game.HomePassAttempts,
		},
		{
			"Label":     "QB Rating",
			"AwayValue": v.Game.AwayQBRating,
			"HomeValue": v.Game.HomeQBRating,
		},
		{
			"Label":     "Total Yards",
			"AwayValue": v.Game.AwayTotalYards,
			"HomeValue": v.Game.HomeTotalYards,
		},
		// Defensive Stats
		{
			"Label":     "Sacks",
			"AwayValue": v.Game.AwaySacks,
			"HomeValue": v.Game.HomeSacks,
		},
		{
			"Label":     "Interceptions",
			"AwayValue": v.Game.AwayInterceptions,
			"HomeValue": v.Game.HomeInterceptions,
		},
		{
			"Label":     "Fumbles Lost",
			"AwayValue": v.Game.AwayFumblesLost,
			"HomeValue": v.Game.HomeFumblesLost,
		},
		{
			"Label":     "Turnovers",
			"AwayValue": v.Game.AwayTurnovers,
			"HomeValue": v.Game.HomeTurnovers,
		},
		{
			"Label":     "Tackles",
			"AwayValue": v.Game.AwayTackles,
			"HomeValue": v.Game.HomeTackles,
		},
		{
			"Label":     "Tackles for Loss",
			"AwayValue": v.Game.AwayTacklesForLoss,
			"HomeValue": v.Game.HomeTacklesForLoss,
		},
		{
			"Label":     "Passes Defended",
			"AwayValue": v.Game.AwayPassesDefended,
			"HomeValue": v.Game.HomePassesDefended,
		},
		// Efficiency Stats
		{
			"Label":     "3rd Down %",
			"AwayValue": v.Game.AwayThirdDownPct,
			"HomeValue": v.Game.HomeThirdDownPct,
		},
		{
			"Label":     "3rd Down Conv",
			"AwayValue": v.Game.AwayThirdDowns,
			"HomeValue": v.Game.HomeThirdDowns,
		},
		{
			"Label":     "First Downs",
			"AwayValue": v.Game.AwayFirstDowns,
			"HomeValue": v.Game.HomeFirstDowns,
		},
		{
			"Label":     "Penalty Yards",
			"AwayValue": v.Game.AwayPenaltyYards,
			"HomeValue": v.Game.HomePenaltyYards,
		},
		{
			"Label":     "Time of Possession",
			"AwayValue": v.Game.AwayTimeOfPossession,
			"HomeValue": v.Game.HomeTimeOfPossession,
		},
	}

	return map[string]interface{}{
		"Game":  v.Game,
		"Stats": stats,
	}
}

func (v *NFLBoxView) TemplateString() string {
	return `
		{{ $ScoreFontSize := 32 }}
		{{ $DetailFontSize := 14}}
		{{ $RecordFontSize := 10}}
		{{ $ScoreFontColor := "#F2FF00FF" }}
		{{ $LogoSize := 64 }}

		<template dir="col" size-x="{{ $MatrixSizex }}" size-y="{{ $MatrixSizey }}">

			<template justify="space-around" size-x="100%" size-y="50%">
				<template dir="col" justify="space-between" align="center"  size-x="40%" size-y="100%">
		    		<image size-x="{{ $LogoSize }}" size-y="{{ $LogoSize }}" src="{{ .Game.AwayLogo }}"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $RecordFontSize }}">({{ .Game.AwayWins }}-{{ .Game.AwayLosses }}-{{ .Game.AwayTies }})</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.AwayScore }}</text>
				</template>

				<template dir="col" justify="space-around" align="center" size-x="20%" size-y="100%">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ .Game.QuarterMinRemaining }}:{{ .Game.QuarterSecRemaining }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ CardinalToOrdinal .Game.Quarter}}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ CardinalToOrdinal .Game.Down}}{{ "&" }}{{ .Game.YardsRemaining }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ .Game.LineOfScrimmage }}</text>
				</template>

				<template dir="col" justify="space-between"  align="center" size-x="40%" size-y="100%">
		    		<image size-x="{{ $LogoSize }}" size-y="{{ $LogoSize }}" src="{{ .Game.HomeLogo }}"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $RecordFontSize }}">({{ .Game.HomeWins }}-{{ .Game.HomeLosses }}-{{ .Game.HomeTies }})</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.HomeScore }}</text>
				</template>
			</template>

			<template dir="col" size-x="100%" size-y="50%" overflow-y="scroll-bounce" scroll-speed="10">
				{{ range .Stats }}
				<comparison-bar 
					value-placement="side"
					left-value="{{ .AwayValue }}"
					right-value="{{ .HomeValue }}"
					stat-label="{{ .Label }}"
					show-stat-label="true"
					left-color="{{ $.Game.AwayColor }}"
					right-color="{{ $.Game.HomeColor }}"
					bar-height="6"
					show-values="true"
					size-x="100%"
					size-y="30%">
				</comparison-bar>
				{{ end }}
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
