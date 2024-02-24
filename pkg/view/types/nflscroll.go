package types

import (
	"errors"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	c "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type NFLScrollView struct {
	c.BaseView

	Date             string
	SportsFeedClient *d.SportsFeed
	Games            d.DailyGamesNFLResponse
	Layout           string
}

type NFLScrollViewConfig struct {
	Layout string `json:"layout"`
	Date   string `json:"date"`
}

func (vc *NFLScrollViewConfig) Validate() error {
	if vc.Layout == "" {
		vc.Layout = "flat"
	}
	if vc.Date == "" {
		return errors.New("'date' field is required")
	}
	return nil
}

func NFLScrollViewCreate(viewConfig c.ViewConfig) (c.View, error) {
	config, ok := viewConfig.(*NFLScrollViewConfig)
	if !ok {
		return nil, errors.New("Error asserting type NFLScrollViewConfig")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := d.SportsFeedClient()

	return &NFLScrollView{
		Date:             config.Date,
		Layout:           config.Layout,
		SportsFeedClient: client,
	}, nil
}

func (v *NFLScrollView) Refresh() {
	// fetch the games
	v.Games = v.SportsFeedClient.FetchDailyNFLGamesInfo(v.Date)
}

func (v *NFLScrollView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Games": v.Games,
	}
}

func (v *NFLScrollView) TemplateString() string {
	var s string

	if v.Layout == "flat" {
		s = `
			<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
				<scroller scrollX="-1" scrollY="0" sizeX="10000" sizeY="100%">
					<template align="center" sizeX="10000" sizeY="{{ $MatrixSizey }}">
					    {{ range .Games.Games }}

					    {{ if eq .Schedule.PlayedStatus "UNPLAYED"}}
						<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}"> @ </text>
					    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Schedule.StartTime | FormatDate }} </text>
					    {{ else }}
					    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
					    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
					    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
					    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.HomeScoreTotal }}  </text>
					    {{ end}}

					    {{ if eq .Score.CurrentQuarterSecondsRemaining nil }}
					    {{ else }}
					    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentQuarterSecondsRemaining }}  </text>
					    {{ end }}
					    {{ if eq .Schedule.PlayedStatus "COMPLETED" }}
					    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">FINAL  </text>
					    {{ else if eq .Score.CurrentQuarter nil }}
					    {{ else }}
					    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentQuarter }}  </text>
					    {{ end }}
					    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">• </text>
					    {{ end }}
					</template>
				</scroller>
			 </template>
		`
	} else if v.Layout == "stack" {
		s = `
			<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">

				<scroller scrollX="-1" scrollY="0" sizeX="10000" sizeY="100%">
					<template sizeX="100%" sizeY="100%">
					    {{ range .Games.Games }}
					    <container sizeX="50" sizeY="100%">
					    	<template dir="col" justify="center" sizeX="100%" sizeY="100%">
					    		<container sizeX="100%" sizeY="50%"> 
					    			<template justify="space-between" align="center" sizeX="100%" sizeY="100%">
							    		<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
										<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}</text>
									</template>
								</container>
					    		<container sizeX="100%" sizeY="50%"> 
					    			<template justify="space-between" align="center" sizeX="100%" sizeY="100%">
							    		<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
										<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.HomeScoreTotal }}</text>
									</template>
								</container>
					    	</template>
					    </container>
					    <container sizeX="50" sizeY="100%">
					    	<template dir="col" justify="center" align="center" sizeX="100%" sizeY="100%">
								{{ if eq .Schedule.PlayedStatus "UNPLAYED"}}
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Schedule.StartTime | FormatDate }} </text>
							    {{ end}}
							    {{ if eq .Schedule.PlayedStatus "COMPLETED" }}
							    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">FINAL</text>
							    {{ else if eq .Score.CurrentQuarter nil }}
							    {{ else }}
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentQuarter }}</text>
							    {{ end }}

								{{ if eq .Score.CurrentQuarterSecondsRemaining nil }}
							    {{ else }}
							    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentQuarterSecondsRemaining }}</text>
							    {{ end }}
					    	</template>
					    </container>
					    {{ end }}
					</template>
				</scroller>
			 </template>
		`
	}

	return s
}

func init() {
	c.RegisterView("nflscroll", c.RegisteredView{
		NewConfig: func() c.ViewConfig { return &NFLScrollViewConfig{} },
		NewView:   NFLScrollViewCreate,
	})
}
