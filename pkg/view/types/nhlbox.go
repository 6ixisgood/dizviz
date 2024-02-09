package types

import (
	"errors"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
)

type NHLBoxView struct {
	c.BaseView

	Date				string
	SportsFeedClient	*d.SportsFeed
	Games				d.DailyGamesNHLResponse
}

type NHLBoxViewConfig struct {
	Date	string		`json:"date"`
}

func (vc *NHLBoxViewConfig) Validate() error {
	if vc.Date == "" {
		return errors.New("'Date' field is required")
	}
	return nil
}

func NHLBoxViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*NHLBoxViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type NHLBoxViewConfig")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := d.SportsFeedClient()

	return &NHLBoxView{
		Date: config.Date,
		SportsFeedClient: client,
	}, nil
}

func (v *NHLBoxView) Refresh() {
	// fetch the games
	v.Games = v.SportsFeedClient.FetchDailyNHLGamesInfo(v.Date)
}

func (v *NHLBoxView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Games": v.Games,
	}
}

func (v *NHLBoxView) TemplateString() string {
	return `
		<template sizeX="{{ $MatrixSizex  }}" sizeY="{{ $MatrixSizey }}">
			<scroller scrollX="-1" scrollY="0">
				<template sizeX="10000" sizeY="{{ $MatrixSizey}}">
				    {{ range .Games.Games }}
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nhl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nhl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.HomeScoreTotal }}  </text>
				    {{ if eq .Score.CurrentPeriodSecondsRemaining nil }}
				    {{ else }}
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriodSecondsRemaining }}  </text>
				    {{ end }}
				    {{ if eq .Schedule.PlayedStatus "COMPLETED" }}
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">FINAL  </text>
				    {{ else if eq .Score.CurrentPeriod nil }}
				    {{ else }}
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriod }}  </text>
				    {{ end }}
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">• </text>
				    {{ end }}
				</template>
			</scroller>
		 </template>
	`
}

func init() {
	c.RegisterView("nhlbox", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &NHLBoxViewConfig{} },
		NewView: NHLBoxViewCreate,
	})
}