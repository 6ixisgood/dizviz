package types

import (
	"fmt"
	"time"

	c "github.com/6ixisgood/disco/pkg/view/common"
	"github.com/6ixisgood/disco/pkg/view/model"
)

type MatchupsScrollView struct {
	c.BaseView

	Layout string    `json:"layout" spec:"label='Layout'"`
	Date   time.Time `json:"date" spec:"required='true',label='Date'"`
	League string    `json:"league" spec:"required='true',label='League'"`

	Matchups []model.Matchup
}

func (v *MatchupsScrollView) Init(configJSON string, ctx c.ViewContext) error {
	if err := c.InitViewFromJSON(v, configJSON, ctx); err != nil {
		return err
	}

	v.Refresh()
	return nil
}

func (v *MatchupsScrollView) Refresh() {
	// fetch the games
	v.Matchups = model.FetchLeagueMatchupsByDate(v.League, v.Date)
	fmt.Println(v.Matchups)
}

func (v *MatchupsScrollView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Matchups": v.Matchups,
	}
}

func (v *MatchupsScrollView) TemplateString() string {
	var s string
	s = `
	<template size-x="{{ $Matrixsize-x }}" size-y="{{ $Matrixsize-y }}">
		<scroller scrollX="-1" scrollY="0" size-x="10000" size-y="100%">
			<template align="center" size-x="10000" size-y="{{ $Matrixsize-y }}">
			    {{ range .Matchups }}

			    {{ if eq .PlayedStatus "UNPLAYED"}}
				<image size-x="{{ $DefaultImagesize-x }}" size-y="{{ $DefaultImagesize-y }}" src="{{ (index .Teams 0).LogoSrc }}"></image>
				<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">@</text>
			    <image size-x="{{ $DefaultImagesize-x }}" size-y="{{ $DefaultImagesize-y }}" src="{{ (index .Teams 1).LogoSrc }}"></image>
				<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .StartTime }}</text>
			    {{ else }}
			    <image size-x="{{ $DefaultImagesize-x }}" size-y="{{ $DefaultImagesize-y }}" src="{{ (index .Teams 0).LogoSrc }}"></image>
			    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ (index .Teams 0).Score }}</text>
			    <image size-x="{{ $DefaultImagesize-x }}" size-y="{{ $DefaultImagesize-y }}" src="{{ (index .Teams 1).LogoSrc }}"></image>
			    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ (index .Teams 1).Score }}</text>
			    {{ end}}

			    {{ if eq .PlayedStatus "COMPLETED" }}
			    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}"> FINAL </text>
			    {{ else if eq .Period nil }}
			    {{ else }}
			    {{ if eq .PeriodMinRemaining nil }}
			    {{ else }}
			    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .PeriodMinRemaining }}:{{ .PeriodSecRemaining }}</text>
			    {{ end }}
			    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Period }}</text>
			    {{ end }}
			    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">  •  </text>
			    {{ end }}
			</template>
		</scroller>
	 </template>
	`
	return s
}

func init() {
	c.RegisterView("matchups-scroll", func() c.View { return &MatchupsScrollView{} })
}
