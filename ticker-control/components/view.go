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
	    "split-view":		SplitViewCreate,
	    "rainbow":			RainbowViewCreate,
	    "train":			TrainViewCreate,
	    "pong":				PongViewCreate,
	    "particle":			ParticlesViewCreate,
	    "colorwave":		ColorWaveViewCreate,
	} 
	GeneralConfig = ViewGeneralConfig{}
)

func SetViewGeneralConfig(config ViewGeneralConfig) {
	GeneralConfig = config
}

type View interface{
	Template()		string
	Refresh()
}

type ViewGeneralConfig struct {
	MatrixRows			int
	MatrixCols			int
	ImageDir			string
	CacheDir			string
	DefaultImageSizeX	int
	DefaultImageSizeY	int
	DefaultFontSize		int
	DefaultFontColor	string
	DefaultFontStyle	string
	DefaultFontType		string
	SportsFeedUsername	string
	SportsFeedPassword	string
}

// ---------------------------
// NHL Daily Games
// ---------------------------
type NHLDailyGamesView struct {
	Date				string
	SportsFeedClient	d.SportsFeed
	Games				d.DailyGamesNHLResponse
}

func NHLDailyGamesViewCreate(config map[string]string) View {
	client := d.NewSportsFeed("",
		d.BasicAuthCredentials{
			Username: GeneralConfig.SportsFeedUsername,
			Password: GeneralConfig.SportsFeedPassword,
		},
	)
	return &NHLDailyGamesView{
		Date: config["date"],
		SportsFeedClient: client,
	}
}

func (v *NHLDailyGamesView) Refresh() {
	// fetch the games
	v.Games = v.SportsFeedClient.FetchDailyNHLGamesInfo(v.Date)
}

func (v *NHLDailyGamesView) Template() string {
	tmplStr := `
		{{ $MatrixSizex :=  .Config.MatrixRows }}
		{{ $MatrixSizey := .Config.MatrixCols }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
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

	tmpl, err := template.New("temp").Funcs(template.FuncMap{
		"NilOrDefault": func() string { return "N/A" },
	}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
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
	SportsFeedClient	d.SportsFeed
	Games		d.DailyGamesNFLResponse
	Layout		string
}

func NFLDailyGamesViewCreate(config map[string]string) View {
	client := d.NewSportsFeed("",
		d.BasicAuthCredentials{
			Username: GeneralConfig.SportsFeedUsername,
			Password: GeneralConfig.SportsFeedPassword,
		},
	)
	if _, ok := config["layout"]; !ok {
		config["layout"] = "flat"
	}	
	return &NFLDailyGamesView{
		Date: config["date"],
		Layout: config["layout"],
		SportsFeedClient: client,
	}
}

func (v *NFLDailyGamesView) Refresh() {
	// fetch the games
	v.Games = v.SportsFeedClient.FetchDailyNFLGamesInfo(v.Date)
}

func (v *NFLDailyGamesView) Template() string {
	var tmplStr string

	if v.Layout == "flat" {
		tmplStr = `
			{{ $MatrixSizex :=  .Config.MatrixRows }}
			{{ $MatrixSizey := .Config.MatrixCols }}
			{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
			{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
			{{ $DefaultFontSize := .Config.DefaultFontSize }}
			{{ $DefaultFontType := .Config.DefaultFontType }}
			{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
			{{ $DefaultFontColor := .Config.DefaultFontColor }}
			{{ $ImageDir := .Config.ImageDir }}
			{{ $CacheDir := .Config.CacheDir }}
			<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
				<scroller scrollX="-1" scrollY="0">
					<template sizeX="10000" sizeY="{{ $MatrixSizey }}">
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
		tmplStr = `
			{{ $MatrixSizex :=  .Config.MatrixRows }}
			{{ $MatrixSizey := .Config.MatrixCols }}
			{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
			{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
			{{ $DefaultFontSize := .Config.DefaultFontSize }}
			{{ $DefaultFontType := .Config.DefaultFontType }}
			{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
			{{ $DefaultFontColor := .Config.DefaultFontColor }}
			{{ $ImageDir := .Config.ImageDir }}
			{{ $CacheDir := .Config.CacheDir }}

			<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">

				<scroller scrollX="-1" scrollY="0">
					<template sizeX="10000" sizeY="{{ $MatrixSizey }}">
					    {{ range .Games.Games }}
					    <h-split>
							<template slot="1" sizeX="10000" sizeY="{{ $MatrixSizey }}">
								<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>

								{{ if eq .Schedule.PlayedStatus "UNPLAYED"}}
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Schedule.StartTime | FormatDate }} </text>
							    {{ end}}
							    {{ if eq .Schedule.PlayedStatus "COMPLETED" }}
							    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">FINAL  </text>
							    {{ else if eq .Score.CurrentQuarter nil }}
							    {{ else }}
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentQuarter }}  </text>
							    {{ end }}

							</template>
							<template slot="2" sizeX="10000" sizeY="{{ $MatrixSizey }}">
					    		<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
								<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.HomeScoreTotal }}  </text>

								 {{ if eq .Score.CurrentQuarterSecondsRemaining nil }}
							    {{ else }}
							    <text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Score.CurrentQuarterSecondsRemaining }}  </text>
							    {{ end }}
							</template>
						</h-split>
						{{ end }}
					</template>
				</scroller>
			 </template>
		`
	}

	tmpl, err := template.New("temp").Funcs(template.FuncMap{
		"FormatDate": FormatDate,
	}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
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
type SplitView struct {
}


func SplitViewCreate(config map[string]string) View {
	return &SplitView{}
}

func (v *SplitView) Refresh() {
	// pass
}

func (v *SplitView) Template() string {
	tmplStr := `
		{{ $MatrixSizex :=  .Config.MatrixRows }}
		{{ $MatrixSizey := .Config.MatrixCols }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<h-split>
				<template slot="1" sizeX="10000" sizeY="{{ $MatrixSizey }}">
					<scroller scrollX="-1" scrollY="0">
						<template sizeX="100000" sizeY="{{ $MatrixSizey }}">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">This is the top line</text>
						</template>
					</scroller>
				</template>
				<template slot="2" sizeX="10000" sizeY="{{ $MatrixSizey }}">
					<scroller scrollX="-1" scrollY="0">
						<template sizeX="100000" sizeY="{{ $MatrixSizey }}">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">This is the bottom line</text>
						</template>
					</scroller>
				</template>
				<template slot="3" sizeX="10000" sizeY="{{ $MatrixSizey }}">
					<scroller scrollX="-1" scrollY="0">
						<template sizeX="100000" sizeY="{{ $MatrixSizey }}">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">This is the top line</text>
						</template>
					</scroller>
				</template>
				<template slot="4" sizeX="10000" sizeY="{{ $MatrixSizey }}">
					<scroller scrollX="-1" scrollY="0">
						<template sizeX="100000" sizeY="{{ $MatrixSizey }}">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">This is the bottom line</text>
						</template>
					</scroller>
				</template>
			</h-split>
		 </template>
	`

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
	} 

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}


// Rainbow text testing
type RainbowView struct {
}


func RainbowViewCreate(config map[string]string) View {
	return &RainbowView{}
}

func (v *RainbowView) Refresh() {
}

func (v *RainbowView) Template() string {
	tmplStr := `
		{{ $MatrixSizex :=  .Config.MatrixRows }}
		{{ $MatrixSizey := .Config.MatrixCols }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="{{ $DefaultFontSize }}">Rainbow!</rainbow-text>
		 </template>
	`

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
	} 

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}



// Train view
type TrainView struct {
}

func TrainViewCreate(config map[string]string) View {
	return &TrainView{}
}

func (v *TrainView) Refresh() {
}

func (v *TrainView) Template() string {
	tmplStr := `
		{{ $MatrixSizex :=  .Config.MatrixRows }}
		{{ $MatrixSizey := .Config.MatrixCols }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">

			<scenic-train sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}" color="{{ $DefaultFontColor }}"></scenic-train>

		 </template>
	`

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
	} 

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}


// Pong view
type PongView struct {
}


func PongViewCreate(config map[string]string) View {
	return &PongView{}
}

func (v *PongView) Refresh() {
}

func (v *PongView) Template() string {
	tmplStr := `
		{{ $MatrixSizex :=  .Config.MatrixRows }}
		{{ $MatrixSizey := .Config.MatrixCols }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<pong sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}" color="{{ $DefaultFontColor }}" ballRadius="2" paddleHeight="15" paddleWidth="5"></pong>
		 </template>
	`

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
	} 

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}

// Particles view
type ParticlesView struct {
}


func ParticlesViewCreate(config map[string]string) View {
	return &ParticlesView{}
}

func (v *ParticlesView) Refresh() {
}

func (v *ParticlesView) Template() string {
	tmplStr := `
		{{ $MatrixSizex :=  .Config.MatrixRows }}
		{{ $MatrixSizey := .Config.MatrixCols }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<gravity-particles sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}"></gravity-particles>
		 </template>
	`

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
	} 

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}

// ColorWave view
type ColorWaveView struct {
}


func ColorWaveViewCreate(config map[string]string) View {
	return &ColorWaveView{}
}

func (v *ColorWaveView) Refresh() {
}

func (v *ColorWaveView) Template() string {
	tmplStr := `
		{{ $MatrixSizex :=  .Config.MatrixRows }}
		{{ $MatrixSizey := .Config.MatrixCols }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<colorwave sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}"></colorwave>
		 </template>
	`

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
	}  

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}	

	content := buf.String()

	return content
}




// HELPERS

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