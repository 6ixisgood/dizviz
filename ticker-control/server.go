package main

import (
	// "bufio"
	"bytes"
	// "os"
	// "fmt"
	// "time"
	// "log"
	"github.com/sixisgoood/matrix-ticker/sports_data"
	"github.com/sixisgoood/matrix-ticker/animations"	
	"text/template"
)

type Config struct {
	DefaultImageSizex		int
	DefaultImageSizey		int
	DefaultFontSize			int
	DefaultFontType			string
	DefaultFontStyle		string
	DefaultFontColor		string
}

type TemplateData struct {
	Matrix		animations.Matrix
	Games		sports_data.DailyGamesNHLResponse
	Config		Config	
}


var(
// 	content = `
// {{ $DefaultImageSizex := .Config.DefaultImageSizex }}
// {{ $DefaultImageSizey := .Config.DefaultImageSizey }}
// {{ $DefaultFontSize := .Config.DefaultFontSize }}
// {{ $DefaultFontType := .Config.DefaultFontType }}
// {{ $DefaultFontStyle := .Config.DefaultFontStyle }}
// {{ $DefaultFontColor := .Config.DefaultFontColor }}
// <matrix sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}">
// 	<content sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}" posx="0" posy="0" scrollx="-5">
// 		{{ range .Games.Games }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.HomeScoreTotal }}  </text>
// 		{{ if eq .Score.CurrentPeriodSecondsRemaining nil }}
// 		{{ else }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriodSecondsRemaining }}  </text>
// 		{{ end }}
// 		{{ if eq .Schedule.PlayedStatus "COMPLETED" }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">FINAL  </text>
// 		{{ else if eq .Score.CurrentPeriod nil }}
// 		{{ else }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriod }}  </text>
// 		{{ end }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">• </text>
// 		{{ end }}

// 	</content>
// </matrix>
// `
// 	content = `
// {{ $DefaultImageSizex := .Config.DefaultImageSizex }}
// {{ $DefaultImageSizey := .Config.DefaultImageSizey }}
// {{ $DefaultFontSize := .Config.DefaultFontSize }}
// {{ $DefaultFontType := .Config.DefaultFontType }}
// {{ $DefaultFontStyle := .Config.DefaultFontStyle }}
// {{ $DefaultFontColor := .Config.DefaultFontColor }}
// <matrix sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}">
// 	<content sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}" posx="0" posy="0" scrollx="-5">
// 		{{ range .Games.Games }}
// 		<image sizex="{{ $DefaultImageSizex }}" sizey="{{ $DefaultImageSizey }}" filepath="/home/andrew/Lab/matrix-ticker/ticker-control/sports_data/images/nhl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
// 		<image sizex="{{ $DefaultImageSizex }}" sizey="{{ $DefaultImageSizey }}" filepath="/home/andrew/Lab/matrix-ticker/ticker-control/sports_data/images/nhl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.HomeScoreTotal }}  </text>
// 		{{ if eq .Score.CurrentPeriodSecondsRemaining nil }}
// 		{{ else }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriodSecondsRemaining }}  </text>
// 		{{ end }}
// 		{{ if eq .Schedule.PlayedStatus "COMPLETED" }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">FINAL  </text>
// 		{{ else if eq .Score.CurrentPeriod nil }}
// 		{{ else }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriod }}  </text>
// 		{{ end }}
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">• </text>
// 		{{ end }}

// 	</content>
// </matrix>
// `
content = `
<matrix sizex="256" sizey="128">
	<content sizex="256" sizey="128" posx="0" posy="0" scrollx="-2">
 		<text font="Ubuntu" fontstyle="Regular" color="#FFFFFFFF" fontsize="10">I hope that a study of very long sentences will arm you with strategies that are almost as diverse as the sentences themselves, such as: starting each clause with the same word, tilting with dependent clauses toward a revelation at the end, padding with parentheticals, showing great latitude toward standard punctuation, rabbit-trailing away from the initial subject, encapsulating an entire life, and lastly, as this sentence is, celebrating the list.</text>
	</content>
</matrix>`

)

func Serve() {

	HandleRequest()

	// for {
	// 	// "listening"
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	fmt.Println("Enter text: ")
	// 	scanner.Scan()
	// 	text := scanner.Text()
	// 	log.Printf(text)

	// 	time.Sleep(5 * time.Second)

	// 	go HandleRequest()
	// }
}


func HandleRequest() {
	config := Config{
		DefaultImageSizex: 256,
		DefaultImageSizey: 128,
		DefaultFontSize: 84,
		DefaultFontType: "Ubuntu",
		DefaultFontStyle: "Regular",
		DefaultFontColor: "#ffffffff",
	}

	data := TemplateData{
		Matrix: animations.Matrix{Sizex: 256, Sizey: 128},
		Games:  GetGames(),
		Config: config,
	}

	tmpl, err := template.New("temp").Parse(content)
	if err != nil {
		panic(err)
	}


	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}	

	content := buf.String()
	animation := animations.NewAnimation(content)
	SetLiveAnimation(animation)	

}

