package types

import (
	"errors"
	"time"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
)

type NFLBoxView struct {
	c.BaseView

	Matchup		string
	SportsFeedClient	*d.SportsFeed
	Game		d.NFLBoxScoreResponseFormatted
	Layout		string
	dataRefresh *util.Refresher
}

type NFLBoxViewConfig struct {
	Matchup		string		`json:"matchup"`
}

func (vc *NFLBoxViewConfig) Validate() error {
	if vc.Matchup == "" {
		return errors.New("'matchup' field is required")
	}
	return nil
}

func NFLBoxViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*NFLBoxViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type NFLBoxViewConfig")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := d.SportsFeedClient()

	return &NFLBoxView{
		Matchup: config.Matchup,
		SportsFeedClient: client,
	}, nil
}

func (v *NFLBoxView) Init() {
	// init ticker and stop chan
    v.dataRefresh = util.RefresherCreate(60 * time.Second, v.RefreshData)
    v.RefreshData()
    v.dataRefresh.Start()
}

func (v *NFLBoxView) RefreshData() {
	// fetch the games
	v.Game = v.SportsFeedClient.FetchNFLBoxScore(v.Matchup)
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

		<template dir="col" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">

			<container sizeX="100%" sizeY="50%">
				<template justify="space-around" sizeX="100%" sizeY="100%">
					<container sizeX="40%" sizeY="100%">
						<template dir="col" justify="space-between" align="center"  sizeX="100%" sizeY="100%">
				    		<image sizeX="{{ $LogoSize }}" sizeY="{{ $LogoSize }}" src="{{ $ImageDir }}/nfl/{{ .Game.AwayAbbreviation }}.png"></image>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">({{ .Game.AwayWins }}-{{ .Game.AwayLosses }}-{{ .Game.AwayTies }})</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.AwayScore }}</text>
						</template>
					</container>

					<container sizeX="20%" sizeY="100%">
						<template dir="col" justify="space-around" align="center" sizeX="100%" sizeY="100%">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ .Game.QuarterMinRemaining }}:{{ .Game.QuarterSecRemaining }}</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ CardinalToOrdinal .Game.Quarter}}</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ CardinalToOrdinal .Game.Down}}{{ "&" }}{{ .Game.YardsRemaining }}</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">{{ .Game.LineOfScrimmage }}</text>
						</template>
					</container>

					<container sizeX="40%" sizeY="100%">
						<template dir="col" justify="space-between"  align="center" sizeX="100%" sizeY="100%">
				    		<image sizeX="{{ $LogoSize }}" sizeY="{{ $LogoSize }}" src="{{ $ImageDir }}/nfl/{{ .Game.HomeAbbreviation }}.png"></image>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">({{ .Game.HomeWins }}-{{ .Game.HomeLosses }}-{{ .Game.HomeTies }})</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.HomeScore }}</text>
						</template>
					</container>
				</template>
			</container>

			<container sizeX="100%" sizeY="50%">
				<template justify="space-between" sizeX="100%" sizeY="100%">
					<container sizeX="45%" sizeY="100%">
						<template sizeX="100%" sizeY="100%" dir="col" justify="space-around">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Passing: {{ .Game.AwayPassYards }}</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Rushing: {{ .Game.AwayRushYards }}</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Sacks: {{ .Game.AwaySacks }}</text>
						</template>
					</container>

					<container sizeX="45%" sizeY="100%">
						<template sizeX="100%" sizeY="100%" dir="col" justify="space-around">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Passing: {{ .Game.HomePassYards }}</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Rushing: {{ .Game.HomeRushYards }}</text>
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DetailFontSize }}">Sacks: {{ .Game.HomeSacks }}</text>
						</template>
					</container>
				</template>
			</container>

		 </template>
		`
}

func init() {
	c.RegisterView("nflbox", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &NFLBoxViewConfig{} },
		NewView: NFLBoxViewCreate,
	})
}