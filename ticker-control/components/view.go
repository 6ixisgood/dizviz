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
	    "nfl-single-game":  NFLSingleGameViewCreate,
	    "sleeper-matchups": SleeperMatchupsViewCreate,
	    "split-view":		SplitViewCreate,
	    "rainbow":			RainbowViewCreate,
	    "train":			TrainViewCreate,
	    "pong":				PongViewCreate,
	    "particle":			ParticlesViewCreate,
	    "colorwave":		ColorWaveViewCreate,
	    "image-player":		ImagePlayerViewCreate,
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
// NFL Single Game
// ---------------------------
type NFLSingleGameView struct {
	Matchup		string
	SportsFeedClient	d.SportsFeed
	Game		d.NFLBoxScoreResponse
	Layout		string
}

func NFLSingleGameViewCreate(config map[string]string) View {
	client := d.NewSportsFeed("",
		d.BasicAuthCredentials{
			Username: GeneralConfig.SportsFeedUsername,
			Password: GeneralConfig.SportsFeedPassword,
		},
	)

	return &NFLSingleGameView{
		Matchup: config["matchup"],
		SportsFeedClient: client,
	}
}

func (v *NFLSingleGameView) Refresh() {
	// fetch the games
	v.Game = v.SportsFeedClient.FetchNFLBoxScore(v.Matchup)
}

func (v *NFLSingleGameView) Template() string {
	var tmplStr string

	tmplStr = `
		{{ $MatrixSizex :=  .Config.MatrixCols }}
		{{ $MatrixSizey := .Config.MatrixRows }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}

		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
		    <h-split sizeX="50" sizeY="{{ $MatrixSizey }}">
				<template sizeX="50" sizeY="{{ $MatrixSizey }}">
		    		<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Game.Game.AwayTeam.Abbreviation }}.png"></image>
		    	</template>
				<template sizeX="50" sizeY="{{ $MatrixSizey }}">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">   {{ .Game.Scoring.AwayScoreTotal }}  </text>
				</template>
			</h-split>


			<h-split sizeX="28" sizeY="{{ $MatrixSizey }}">
				<template sizeX="28" sizeY="{{ $MatrixSizey }}">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.CurrentQuarterSecondsRemaining }}</text>
				</template>
				<template sizeX="28" sizeY="{{ $MatrixSizey }}">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.CurrentQuarter }}/4</text>
				</template>
				<template sizeX="28" sizeY="{{ $MatrixSizey }}">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.CurrentDown }}&#38;{{ .Game.Scoring.CurrentYardsRemaining }}</text>
				</template>
				<template sizeX="28" sizeY="{{ $MatrixSizey }}">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.LineOfScrimmage }}</text>
				</template>
			</h-split>


			<h-split sizeX="50" sizeY="{{ $MatrixSizey }}">
				<template sizeX="50" sizeY="{{ $MatrixSizey }}">
		    		<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Game.Game.HomeTeam.Abbreviation }}.png"></image>
				</template>
				<template sizeX="50" sizeY="{{ $MatrixSizey }}">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">   {{ .Game.Scoring.HomeScoreTotal }}  </text>
				</template>
			</h-split>
		 </template>
		`

	tmpl, err := template.New("temp").Funcs(template.FuncMap{
		"FormatDate": FormatDate,
	}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
		"Game": v.Game,
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
// Sleeper Matchups
// ---------------------------
type SleeperMatchupsView struct {
	League			string
	Week			string
	SleeperClient	d.Sleeper
	CurrentMatchup	[]d.SleeperTeamFormatted
	matchIndex		int
	Phase			int
}

func SleeperMatchupsViewCreate(config map[string]string) View {
	client := d.NewSleeper("")

	return &SleeperMatchupsView{
		League: config["league_id"],
		Week: config["week"],
		SleeperClient: client,
		Phase: -1,
	}
}

func (v *SleeperMatchupsView) Refresh() {
	// fetch the matchups
	matchups := v.SleeperClient.GetMatchupsFormatted(v.League, v.Week)

	if v.Phase < 0 {
		v.Phase = 0
	} else if v.Phase >= 0 && v.Phase < 2 {
		v.Phase += 1
	} else if v.Phase == 2 {
		v.Phase = 1
		v.matchIndex = (v.matchIndex + 1) % len(matchups)
	}

	v.CurrentMatchup = matchups[v.matchIndex]
}

func (v *SleeperMatchupsView) Template() string {
	var tmplStr string

	tmplStr = `
		{{ $MatrixSizex :=  .Config.MatrixCols }}
		{{ $MatrixSizey := .Config.MatrixRows }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := 8 }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}

		{{ if eq .Phase 0 }}

		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<h-split>
				<template sizeX="{{ $MatrixSizex }}" sizeY="28">
					<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12">   Bucket Hats</rainbow-text>
				</template>
				<template sizeX="{{ $MatrixSizex }}" sizeY="8">
					<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12">        &#38;</rainbow-text>
				</template>
				<template sizeX="{{ $MatrixSizex }}" sizeY="28">
					<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12">   Trap Music</rainbow-text>
				</template>
			</h-split>
		</template>

		{{ else if eq .Phase 1}}

		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<h-split sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
				<template sizeX="{{ $MatrixSizex }}" sizeY="32">
				    <h-split sizeX="{{ $MatrixSizex }}" sizeY="32">
						<template sizeX="{{ $MatrixSizex }}" sizeY="32">
							<text sizeX="55" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#FF4542FF" size="{{ $DefaultFontSize }}">{{ .Team1.Name }}</text>
							<text sizeX="5" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#FF4542FF" size="{{ $DefaultFontSize }}">  </text>
							<text sizeX="64" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#FF4542FF" size="{{ $DefaultFontSize }}">{{ .Team2.Name }}</text>
				    	</template>
					</h-split>
				</template>

				<template sizeX="{{ $MatrixSizex }}" sizeY="32">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#F2FF00FF" size="16"> {{ .Team1.Score }} - {{ .Team2.Score }}  </text>
				</template>
			</h-split>
		 </template>

		{{ else if eq .Phase 2}}

		<template sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
		    <h-split sizeX="64" sizeY="{{ $MatrixSizey }}">
				<template sizeX="64" sizeY="{{ $MatrixSizey }}">
					<scroller scrollX="-1" scrollY="0">
						<template sizeX="58" sizeY="{{ $MatrixSizey }}">
							<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#FF4542FF"  size="{{ $DefaultFontSize }}">{{ .Team1.Name }}</text>
							<text sizeX="5" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#FF4542FF" size="{{ $DefaultFontSize }}">  </text>
						</template>
					</scroller>
				</template>
				{{ range $index, $element := .Team1.Players }}
				{{ if lt $index 4 }}
				<template sizeX="60" sizeY="{{ $MatrixSizey }}">
					<text sizeX="45" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
					<text sizeX="17" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#F2FF00FF" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
				</template>
				{{ else }}
				{{ end }}
				{{ end }}
			</h-split>

			<h-split sizeX="64" sizeY="{{ $MatrixSizey }}">
				<template sizeX="64" sizeY="{{ $MatrixSizey }}">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#FF4542FF" size="{{ $DefaultFontSize }}">{{ .Team2.Name }}</text>
		    	</template>
				{{ range $index, $element := .Team2.Players }}
				{{ if lt $index 4 }}
				<template sizeX="64" sizeY="{{ $MatrixSizey }}">
					<text sizeX="45" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
					<text sizeX="17" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="#F2FF00FF" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
				</template>
				{{ else }}
				{{ end }}
				{{ end }}
			</h-split>
		 </template>
		 {{ end }}
		`

	tmpl, err := template.New("temp").Funcs(template.FuncMap{
		"FormatDate": FormatDate,
	}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
		"Team1": v.CurrentMatchup[0],
		"Team2": v.CurrentMatchup[1],
		"Phase": v.Phase,
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
// Split View Test
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




// Image Player view
type ImagePlayerView struct {
	Src		string
}


func ImagePlayerViewCreate(config map[string]string) View {
	return &ImagePlayerView{
		Src: config["src"],
	}
}

func (v *ImagePlayerView) Refresh() {
}

func (v *ImagePlayerView) Template() string {
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
			<image sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}" src="{{ .Src }}" loop="true"></image>
		 </template>
	`

	tmpl, err := template.New("temp").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"Config": GeneralConfig,
		"Src": v.Src,
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