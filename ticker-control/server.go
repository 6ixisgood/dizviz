package main

import (
	"bufio"
	"bytes"
	"os"
	"fmt"
	"log"
	"github.com/sixisgoood/matrix-ticker/sports_data"
	"github.com/sixisgoood/matrix-ticker/animations"
	"text/template"
)

type TemplateData struct {
	Matrix		animations.Matrix
	Games		sports_data.DailyGamesNHLResponse
}

var(
	content = `
<matrix sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}">
	<content sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}" posx="0" posy="0" scrollx="-1">
		{{ range .Games.Games }}
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">{{ .Schedule.AwayTeam.Abbreviation }}  </text>
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">{{ .Score.AwayScoreTotal }}  </text>
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">{{ .Schedule.HomeTeam.Abbreviation }}  </text>
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">{{ .Score.HomeScoreTotal }}  </text>
		{{ if eq .Score.CurrentPeriodSecondsRemaining nil }}
		{{ else }}
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">{{ .Score.CurrentPeriodSecondsRemaining }}  </text>
		{{ end }}
		{{ if eq .Schedule.PlayedStatus "COMPLETED" }}
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">FINAL  </text>
		{{ else if eq .Score.CurrentPeriod nil }}
		{{ else }}
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">{{ .Score.CurrentPeriod }}  </text>
		{{ end }}
		<text font="Ubuntu" fontstyle="Regular" color="#00ffffff" fontsize="10">• </text>
		{{ end }}
	</content>
</matrix>
`
)

func Serve() {

	for {
		// "listening"
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter text: ")
		scanner.Scan()
		text := scanner.Text()
		log.Printf(text)

		go HandleRequest()
	}
}


func HandleRequest() {

	data := TemplateData{
		Matrix: animations.Matrix{Sizex: 64, Sizey: 32},
		Games: sports_data.FetchDailyNHLGamesInfo("2022-2023-regular", "20221012"),
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

