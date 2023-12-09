package components

import (
	"fmt"
	"time"
	"log"
	"bytes"
	"maps"
	"text/template"
	"encoding/xml"
	d "github.com/sixisgoood/matrix-ticker/data"

)


var (
	RegisteredViews = map[string]func(map[string]string) View{
	    "nhl-daily-games":  NHLDailyGamesViewCreate,
	    "nfl-daily-games":	NFLDailyGamesViewCreate,
	    "nfl-single-game":  NFLSingleGameViewCreate,
	    "sleeper-matchups": SleeperMatchupsViewCreate,
	    "text":				TextViewCreate,
	    "particle":			ParticlesViewCreate,
	    "image-player":		ImagePlayerViewCreate,
	} 
	GeneralConfig = ViewGeneralConfig{}
)

func SetViewGeneralConfig(config ViewGeneralConfig) {
	GeneralConfig = config
}

// timed refresh struct
type Refresher struct {
	interval	time.Duration
	ticker		*time.Ticker
	stopChan	chan struct{}
	onRefresh	func()
} 

func RefresherCreate(interval time.Duration, onRefresh func()) *Refresher {
	return &Refresher{
		interval: interval,
		stopChan: make(chan struct{}),
		onRefresh: onRefresh,
	}
}

func (r *Refresher) Start() {
	r.ticker = time.NewTicker(r.interval)
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.onRefresh()
			case <-r.stopChan:
				r.ticker.Stop()
				return
			}
		}
	}()
}

func (r *Refresher) Stop() {
	close(r.stopChan)
}


// define the View interface
type View interface {
	Init()
	Template()				*Template
	SetTemplate(*Template)
	TemplateString()		string
	TemplateData()			map[string]interface{}
	Stop()
}

// Define a base class for a view with basic implementations
type BaseView struct {
	template			*Template
	dataRefresh			*time.Ticker
	templateRefresh		*time.Ticker
	stopChan			chan struct{}
}

func (v *BaseView) Init() {}

func (v *BaseView) Template() *Template {
	return v.template
}

func TemplateRefresh(v View) {
	t := ExecuteViewTemplate(v)
	v.SetTemplate(t)
	t.Init()

}

func (v *BaseView) SetTemplate(t *Template) {
	v.template = t
}
func (v *BaseView) TemplateData() map[string]interface{} {
	return map[string]interface{}{}  
}

func (v *BaseView) TemplateString() string {
	return ""
}

func (v *BaseView) Stop() {}


func ExecuteViewTemplate(v View) *Template {
	// create the template object
	tmpl := template.New("view-template")

	// gather all the custom functions
	funcMap := template.FuncMap{
		"NilOrDefault": func() string { return "N/A" },
	}
	tmpl = tmpl.Funcs(funcMap)

	// construct the template string
	tmplString := `
		{{ $MatrixSizex := .Config.MatrixCols }}
		{{ $MatrixSizey := .Config.MatrixRows }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}

		%s
	`
	tmplString = fmt.Sprintf(tmplString, v.TemplateString())	

	// parse the template string from the view
	tmpl, err := tmpl.Parse(tmplString)
	if err != nil {
		log.Fatalf("Unable to parse view template")
		panic(err)
	}

	// merge data maps
	data := map[string]interface{}{
		"Config": GeneralConfig,
    }
    maps.Copy(data, v.TemplateData())

	// execute the template with the data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatalf("Unable to execute view template")
		panic(err)
	}	

	// convert to string
	tmplStr := buf.String()


	// unmarshall the string
	var t Template
	err = xml.Unmarshal([]byte(tmplStr), &t)
	if err != nil {
		log.Fatalf("Unable to unmarshal xml content: '%v'", err)
	}

	return &t
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
}

// ---------------------------
// NHL Daily Games
// ---------------------------
type NHLDailyGamesView struct {
	BaseView

	Date				string
	SportsFeedClient	*d.SportsFeed
	Games				d.DailyGamesNHLResponse
}

func NHLDailyGamesViewCreate(config map[string]string) View {
	client := d.SportsFeedClient()

	return &NHLDailyGamesView{
		Date: config["date"],
		SportsFeedClient: client,
	}
}

func (v *NHLDailyGamesView) Refresh() {
	// fetch the games
	v.Games = v.SportsFeedClient.FetchDailyNHLGamesInfo(v.Date)
}

func (v *NHLDailyGamesView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Games": v.Games,
	}
}

func (v *NHLDailyGamesView) TemplateString() string {
	return `
		{{ $MatrixSizex := .Config.MatrixCols }}
		{{ $MatrixSizey := .Config.MatrixRows }}
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
}

// ---------------------------
// NFL Daily Games
// ---------------------------
type NFLDailyGamesView struct {
	BaseView

	Date		string
	SportsFeedClient	*d.SportsFeed
	Games		d.DailyGamesNFLResponse
	Layout		string
}

func NFLDailyGamesViewCreate(config map[string]string) View {
	client := d.SportsFeedClient()

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

func (v *NFLDailyGamesView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Games": v.Games,
	} 
}

func (v *NFLDailyGamesView) TemplateString() string {
	var s string

	if v.Layout == "flat" {
		s = `
			{{ $MatrixSizex := .Config.MatrixCols }}
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
			{{ $MatrixSizex := .Config.MatrixCols }}
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

// ---------------------------
// NFL Single Game
// ---------------------------
type NFLSingleGameView struct {
	BaseView

	Matchup		string
	SportsFeedClient	*d.SportsFeed
	Game		d.NFLBoxScoreResponse
	Layout		string
}

func NFLSingleGameViewCreate(config map[string]string) View {
	client := d.SportsFeedClient()

	return &NFLSingleGameView{
		Matchup: config["matchup"],
		SportsFeedClient: client,
	}
}

func (v *NFLSingleGameView) Refresh() {
	// fetch the games
	v.Game = v.SportsFeedClient.FetchNFLBoxScore(v.Matchup)
}

func (v *NFLSingleGameView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Game": v.Game,
	} 
}

func (v *NFLSingleGameView) TemplateString() string {
	return `
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
		{{ $ScoreFontSize := 14 }}

		<template justify="space-around" align="center" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<container sizeX="40%" sizeY="100%">
				<template dir="col" justify="space-around" align="center"  sizeX="100%" sizeY="100%">
		    		<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Game.Game.AwayTeam.Abbreviation }}.png"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.Scoring.AwayScoreTotal }}</text>
				</template>
			</container>

			<container sizeX="20%" sizeY="100%">
				<template dir="col" justify="space-around" align="center" sizeX="100%" sizeY="100%">
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.CurrentQuarterSecondsRemaining }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.CurrentQuarter }}/4</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.CurrentDown }}&#38;{{ .Game.Scoring.CurrentYardsRemaining }}</text>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ .Game.Scoring.LineOfScrimmage }}</text>
				</template>
			</container>

			<container sizeX="40%" sizeY="100%">
				<template dir="col" justify="space-around"  align="center" sizeX="100%" sizeY="100%">
		    		<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ $ImageDir }}/nfl/{{ .Game.Game.HomeTeam.Abbreviation }}.png"></image>
					<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $ScoreFontSize }}">{{ .Game.Scoring.HomeScoreTotal }}</text>
				</template>
			</container>
		 </template>
		`
}

// ---------------------------
// Sleeper Matchups
// ---------------------------
type SleeperMatchupsView struct {
	BaseView

	League			string
	Week			string
	SleeperClient	*d.Sleeper
	CurrentMatchup	[]d.SleeperTeamFormatted
	matchIndex		int
	Phase			int
	dataRefresh		*Refresher
	phaseRefresh	*Refresher
}

func SleeperMatchupsViewCreate(config map[string]string) View {
	client := d.SleeperClient()

	view := &SleeperMatchupsView{
		League: config["league_id"],
		Week: config["week"],
		SleeperClient: client,
		Phase: 0,
	}

	return view
}

func (v *SleeperMatchupsView) Init() {
	// init ticker and stop chan
    v.dataRefresh = RefresherCreate(60 * time.Second, v.RefreshData)
    v.phaseRefresh = RefresherCreate(5 * time.Second, v.RefreshPhase)
    v.RefreshData()
    v.dataRefresh.Start()
    v.phaseRefresh.Start()
}

func (v *SleeperMatchupsView) RefreshData() {
	matchups := v.SleeperClient.GetMatchupsFormatted(v.League, v.Week)
	v.CurrentMatchup = matchups[v.matchIndex]
}

func (v *SleeperMatchupsView) RefreshPhase() {
	// change the phase at a given rate
	if v.Phase < 2 {
		v.Phase += 1
	} else if v.Phase == 2 {
		v.Phase = 1
	}
	TemplateRefresh(v)
}

func (v *SleeperMatchupsView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Team1": v.CurrentMatchup[0],
		"Team2": v.CurrentMatchup[1],
		"Phase": v.Phase,
	}
}

func (v *SleeperMatchupsView) Stop() {
	v.dataRefresh.Stop()
	v.phaseRefresh.Stop()
}

func (v *SleeperMatchupsView) TemplateString() string {
	return `
			{{ $MatrixSizex :=  .Config.MatrixCols }}
			{{ $MatrixSizey := .Config.MatrixRows }}
			{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
			{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
			{{ $DefaultFontSize := 8 }}
			{{ $DefaultFontType := .Config.DefaultFontType }}
			{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
			{{ $DefaultFontColor := .Config.DefaultFontColor }}
			{{ $BenchedColor := "#FF4542FF" }}
			{{ $ScoreColor := "#F2FF00FF" }}
			{{ $PlayingColor := "#FFFFFFFF"}}
			{{ $TeamNameColor := "#66CCFFFF"}}
			{{ $ImageDir := .Config.ImageDir }}
			{{ $CacheDir := .Config.CacheDir }}
	
			{{ if eq .Phase 0 }}
	
			<template dir="col" justify="center" align="center" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">Bucket Hats</rainbow-text>
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">&#38;</rainbow-text>
				<rainbow-text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" size="12" color="{{ $DefaultFontColor }}">Trap Music</rainbow-text>
			</template>
	
			{{ else if eq .Phase 1 }}
	
			<template justify="space-between" align="center" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
				<container sizeX="50%" sizeY="100%">
					<template justify="space-around" align="center" sizeX="100%" sizeY="100%" dir="col">
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team1.Name }}</text>
						<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ .Team1.Avatar }}"></image>
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16"> {{ .Team1.Score }}</text>
					</template>
				</container>
				<container sizeX="50%" sizeY="100%">
					<template justify="space-around" align="center" sizeX="100%" sizeY="100%" dir="col">
						<text sizeX="90%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team2.Name }}</text>
						<image sizeX="{{ $DefaultImageSizex }}" sizeY="{{ $DefaultImageSizey }}" src="{{ .Team2.Avatar }}"></image>
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $ScoreColor }}" size="16"> {{ .Team2.Score }}</text>
					</template>
				</container>
			 </template>
	
			{{ else if eq .Phase 2 }}
	
			<template justify="space-between" align="center" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
				<container sizeX="50%" sizeY="100%">
					<template dir="col" sizeX="100%" sizeY="100%" justify="start" align="start">
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team1.Name }}</text>
	
						<scroller scrollX="0" scrollY="-1" sizeX="100%" sizeY="200">
							<template rr="10" dir="col" sizeX="100%" sizeY="100%">
								{{ range $index, $element := .Team1.Players }}
								<container sizeX="100%" sizeY="5%">
									<template sizeX="100%" sizeY="10">
										<text sizeX="80%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
										<text sizeX="20%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ if $element.Starter }}{{ $ScoreColor }}{{ else }}{{ $BenchedColor }}{{ end }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
									</template>
								</container>
								{{ end }}
							</template>
						</scroller>
					</template>
				</container>
				<container sizeX="50%" sizeY="100%">
					<template dir="col" sizeX="100%" sizeY="100%" justify="start" align="start">
						<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $TeamNameColor }}" size="{{ $DefaultFontSize }}">{{ .Team2.Name }}</text>
	
						<scroller scrollX="0" scrollY="-1" sizeX="100%" sizeY="200">
							<template rr="10" dir="col" sizeX="100%" sizeY="100%">
								{{ range $index, $element := .Team2.Players }}
								<container sizeX="100%" sizeY="5%">
									<template sizeX="100%" sizeY="10">
										<text sizeX="80%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">{{ printf "%s" $element.Name }}</text>
										<text sizeX="20%" font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ if $element.Starter }}{{ $ScoreColor }}{{ else }}{{ $BenchedColor }}{{ end }}" size="{{ $DefaultFontSize }}">{{ printf "%.2f" $element.Points }}</text>
									</template>
								</container>
								{{ end }}
							</template>
						</scroller>
	
					</template>
				</container>
			</template>
			{{ end }}
			`
}


// Particles view
type ParticlesView struct {
	BaseView
}

func ParticlesViewCreate(config map[string]string) View {
	return &ParticlesView{}
}

func (v *ParticlesView) TemplateString() string {
	return `
		{{ $MatrixSizex := .Config.MatrixCols }}
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
			<gravity-particles sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}"></gravity-particles>
		 </template>
	`
}

// Image Player view
type ImagePlayerView struct {
	BaseView

	Src		string
}


func ImagePlayerViewCreate(config map[string]string) View {
	return &ImagePlayerView{
		Src: config["src"],
	}
}

func (v *ImagePlayerView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Src": v.Src,
	}  
}

func (v *ImagePlayerView) TemplateString() string {
	return `
		{{ $MatrixSizex := .Config.MatrixCols }}
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
			<image sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}" src="{{ .Src }}" loop="true"></image>
		 </template>
	`
}


// Text view
type TextView struct {
	BaseView

	Text		string
}


func TextViewCreate(config map[string]string) View {
	return &TextView{
		Text: config["text"],
	}
}

func (v *TextView) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Text": v.Text,
	}  
}

func (v *TextView) TemplateString() string {
	return `
		{{ $MatrixSizex := .Config.MatrixCols }}
		{{ $MatrixSizey := .Config.MatrixRows }}
		{{ $DefaultImageSizex := .Config.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Config.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Config.DefaultFontSize }}
		{{ $DefaultFontType := .Config.DefaultFontType }}
		{{ $DefaultFontStyle := .Config.DefaultFontStyle }}
		{{ $DefaultFontColor := .Config.DefaultFontColor }}
		{{ $ImageDir := .Config.ImageDir }}
		{{ $CacheDir := .Config.CacheDir }}
		<template dir="row" justify="space-between" sizeX="{{ $MatrixSizex }}" sizeY="{{ $MatrixSizey }}">
			<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">One</text>
			<text font="{{ $DefaultFontType }}" style="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" size="{{ $DefaultFontSize }}">Two</text>
		</template>
	`
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