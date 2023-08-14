package components

import (
	"bytes"
	"text/template"
	cd "github.com/sixisgoood/matrix-ticker/content_data"

)


var (
	RegisteredViews = map[string]func() View{
	    "myview":  MyViewCreate,
	} 
)

type View interface{
	Template()		string
	Refresh()
}

type MyView struct {
	Games		cd.DailyGamesNHLResponse
}


func MyViewCreate() View {
	return &MyView{}
}

func (m *MyView) Refresh() {
	// fetch the games
	m.Games = cd.FetchDailyNHLGamesInfo("2022-2023-regular", "20221112")
}

// func (m *MyView) Template() string {
// 	return `
// 	<template sizeX="64" sizeY="64">
// 		<scroller scrollX="-1" scrollY="0">
// 			<template sizeX="500" sizeY="64">
// 				<text font="Ubuntu" style="Regular" size="24" color="#ffffffff">Hi there! How are you?</text>
// 	    	</template>
// 	    </scroller>
// 	</template>
// 	`
// }

func (m *MyView) Template() string {
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
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/content_data/images/nhl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
				    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
				    <image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="/home/andrew/Lab/matrix-ticker/ticker-control/content_data/images/nhl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
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

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Games": m.Games,
	} 


	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content

}


