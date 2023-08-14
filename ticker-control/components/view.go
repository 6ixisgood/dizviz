package components

import (
	"fmt"
	"time"
	"bytes"
	"text/template"
	d "github.com/sixisgoood/matrix-ticker/data"

)


var (
	RegisteredViews = map[string]func(map[string]string) View{
	    "nhl-daily-games":  NHLDailyGamesViewCreate,
	    "nfl-daily-games":	NFLDailyGamesViewCreate,
	} 
)

type View interface{
	Template()		string
	Refresh()
}

// ---------------------------
// NHL Daily Games
// ---------------------------
type NHLDailyGamesView struct {
	Date		string
	Games		d.DailyGamesNHLResponse
}


func NHLDailyGamesViewCreate(config map[string]string) View {
	return &NHLDailyGamesView{
		Date: config["date"],
	}
}

func (v *NHLDailyGamesView) Refresh() {
	// fetch the games
	v.Games = d.FetchDailyNHLGamesInfo(v.Date)
}

func (v *NHLDailyGamesView) Template() string {
	tmplStr := `
		{{ $MatrixSizex := 64 }}
		{{ $MatrixSizey := 64 }}
		{{ $DefaultImageSizex := 32 }}
		{{ $DefaultImageSizey := 32 }}
		{{ $DefaultFontSize := 24 }}
		{{ $DefaultFontType := "Ubuntu" }}
		{{ $DefaultFontStyle := "Regular" }}
		{{ $DefaultFontColor := "#ffffffff" }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<scroller scrollX="-1" scrollY="0">
				<template sizeX="10000" sizeY="{{ $MatrixSizey }}">
				    {{ range .Games.Games }}
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/data/images/nhl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/data/images/nhl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
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

	tmpl, err := template.New("temp").Funcs(template.FuncMap{
		"NilOrDefault": func() string { return "N/A" },
	}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Games": v.Games,
	} 


	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}

// ---------------------------
// NFL Daily Games
// ---------------------------
type NFLDailyGamesView struct {
	Date		string
	Games		d.DailyGamesNFLResponse
}


func NFLDailyGamesViewCreate(config map[string]string) View {
	return &NFLDailyGamesView{
		Date: config["date"],
	}
}

func (v *NFLDailyGamesView) Refresh() {
	// fetch the games
	v.Games = d.FetchDailyNFLGamesInfo(v.Date)
}

func (v *NFLDailyGamesView) Template() string {
	tmplStr := `
		{{ $MatrixSizex := 64 }}
		{{ $MatrixSizey := 64 }}
		{{ $DefaultImageSizex := 32 }}
		{{ $DefaultImageSizey := 32 }}
		{{ $DefaultFontSize := 24 }}
		{{ $DefaultFontType := "Ubuntu" }}
		{{ $DefaultFontStyle := "Regular" }}
		{{ $DefaultFontColor := "#ffffffff" }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<scroller scrollX="-1" scrollY="0">
				<template sizeX="10000" sizeY="{{ $MatrixSizey }}">
				    {{ range .Games.Games }}

				    {{ if eq .Schedule.PlayedStatus "UNPLAYED"}}
					<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/data/images/nfl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}"> @ </text>
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/data/images/nfl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Schedule.StartTime | FormatDate }} </text>
				    {{ else }}
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/data/images/nfl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/data/images/nfl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
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

	tmpl, err := template.New("temp").Funcs(template.FuncMap{
		"FormatDate": FormatDate,
	}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Games": v.Games,
	} 


	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}

func FormatDate(date time.Time) (string, error) {
	fmt.Printf("%v", date)
	destinationTimeZone := "America/New_York"
	destinationLocation, err := time.LoadLocation(destinationTimeZone)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	convertedTime := date.In(destinationLocation)
	return convertedTime.Format("01/01 01:01AM"), nil
}