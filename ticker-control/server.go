package main

// import (
// 	"bufio"
// 	"bytes"
// 	"os"
// 	"fmt"
// 	"time"
// 	"log"
// 	"io/ioutil"
// 	cd "github.com/sixisgoood/matrix-ticker/content_data"
// 	"github.com/sixisgoood/matrix-ticker/animations"	
// 	"text/template"
// 	"gopkg.in/yaml.v2"	
// )

// type Config struct {
// 	DefaultImageSizex		int
// 	DefaultImageSizey		int
// 	DefaultFontSize			int
// 	DefaultFontType			string
// 	DefaultFontStyle		string
// 	DefaultFontColor		string
// 	MatrixSizex				int
// 	MatrixSizey				int
// }




// var(
// // 	content = `
// // {{ $DefaultImageSizex := .Config.DefaultImageSizex }}
// // {{ $DefaultImageSizey := .Config.DefaultImageSizey }}
// // {{ $DefaultFontSize := .Config.DefaultFontSize }}
// // {{ $DefaultFontType := .Config.DefaultFontType }}
// // {{ $DefaultFontStyle := .Config.DefaultFontStyle }}
// // {{ $DefaultFontColor := .Config.DefaultFontColor }}
// // <matrix sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}">
// // 	<content sizex="{{ .Matrix.Sizex }}" sizey="{{ .Matrix.Sizey }}" posx="0" posy="0" scrollx="-5">
// // 		{{ range .Games.Games }}
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.HomeScoreTotal }}  </text>
// // 		{{ if eq .Score.CurrentPeriodSecondsRemaining nil }}
// // 		{{ else }}
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriodSecondsRemaining }}  </text>
// // 		{{ end }}
// // 		{{ if eq .Schedule.PlayedStatus "COMPLETED" }}
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">FINAL  </text>
// // 		{{ else if eq .Score.CurrentPeriod nil }}
// // 		{{ else }}
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.CurrentPeriod }}  </text>
// // 		{{ end }}
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">• </text>
// // 		{{ end }}

// // 	</content>
// // </matrix>
// // `
// 	content = `
// {{ $MatrixSizex := .Config.MatrixSizex }}
// {{ $MatrixSizey := .Config.MatrixSizey }}
// {{ $DefaultImageSizex := .Config.DefaultImageSizex }}
// {{ $DefaultImageSizey := .Config.DefaultImageSizey }}
// {{ $DefaultFontSize := .Config.DefaultFontSize }}
// {{ $DefaultFontType := .Config.DefaultFontType }}
// {{ $DefaultFontStyle := .Config.DefaultFontStyle }}
// {{ $DefaultFontColor := .Config.DefaultFontColor }}
// <matrix sizex="{{ $MatrixSizex }}" sizey="{{ $MatrixSizey }}">
// 	<content sizex="{{ $MatrixSizex }}" sizey="{{ $MatrixSizey }}" posx="0" posy="0" scrollx="-15" aligny="center">
// 		{{ range .Games.Games }}
// 		<image sizex="{{ $DefaultImageSizex }}" sizey="{{ $DefaultImageSizey }}" filepath="/home/andrew/Lab/matrix-ticker/ticker-control/content_data/images/nhl/{{ .Schedule.AwayTeam.Abbreviation }}.png"></image>
// 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Score.AwayScoreTotal }}  </text>
// 		<image sizex="{{ $DefaultImageSizex }}" sizey="{{ $DefaultImageSizey }}" filepath="/home/andrew/Lab/matrix-ticker/ticker-control/content_data/images/nhl/{{ .Schedule.HomeTeam.Abbreviation }}.png"></image>
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
// // content = `
// // <matrix sizex="256" sizey="128">
// // 	<content sizex="256" sizey="128" posx="0" posy="0" scrollx="0">
// //  		<text font="Ubuntu" fontstyle="Regular" color="#FFFFFFFF" fontsize="20">FUCK</text>
// // 	</content>
// // </matrix>`
// // content = `
// // {{ $MatrixSizex := 256 }}
// // {{ $MatrixSizey := 128 }}
// // {{ $DefaultImageSizex := 128 }}
// // {{ $DefaultImageSizey := 64 }}
// // {{ $DefaultFontSize := 12 }}
// // {{ $DefaultFontType := "Roboto" }}
// // {{ $DefaultFontStyle := "Regular" }}
// // {{ $DefaultFontColor := "#ffffffff" }}
// // <matrix sizex="{{ $MatrixSizex }}" sizey="{{ $MatrixSizey }}">
// // 	<content sizex="{{ $MatrixSizex }}" sizey="32" posx="0" posy="0" scrollx="0" aligny="center">
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">{{ .Weather.Location.Name }}, {{ .Weather.Location.Region}}  {{ .Weather.Location.Localtime }} {{ .Weather.Current.TempF }}°</text>
// // 	</content>
// // 	<content sizex="{{ $MatrixSizex }}" sizey="32" posx="0" posy="32" scrollx="0" aligny="center">
// // 	{{ with index .Weather.Forecast.Forecastday 0 }}
// // 		<text font="{{ $DefaultFontType }}" fontstyle="{{ $DefaultFontStyle }}" color="{{ $DefaultFontColor }}" fontsize="{{ $DefaultFontSize }}">Low: {{ .Day.MintempF }}°  High: {{ .Day.MaxtempF }}° Rise: {{ .Astro.Sunrise }} Set: {{ .Astro.Sunset }} Phase: {{ .Astro.MoonPhase }}</text>
// // 	{{ end }}
// // 	</content>
// // </matrix>
// // `

// )

// func Serve() {

// 	go HandleRequest()

// 	for {
// 		// "listening"
// 		scanner := bufio.NewScanner(os.Stdin)
// 		fmt.Println("Enter text: ")
// 		scanner.Scan()
// 		text := scanner.Text()
// 		log.Printf(text)

// 		time.Sleep(5 * time.Second)

// 		go HandleRequest()
// 	}
// }


// func HandleRequest() {
// 	appConfig := GetApplicationConfig()
// 	config := Config{
// 		DefaultImageSizex: 256,
// 		DefaultImageSizey: 128,
// 		DefaultFontSize: 84,
// 		DefaultFontType: "Ubuntu",
// 		DefaultFontStyle: "Regular",
// 		DefaultFontColor: "#ffffffff",
// 		MatrixSizex: appConfig.Matrix.Cols,
// 		MatrixSizey: appConfig.Matrix.Rows,

// 	}


// 	content := LoadTemplate()
// 	animation := animations.NewAnimation(content)
// 	SetLiveAnimation(animation)	

// }

// type TemplateConfig struct {
// 	Name		string		`yaml:"name"`	
// 	Template	string		`yaml:"template"`	
// }

// type TemplateData struct {
// 	Games		cd.DailyGamesNHLResponse
// 	Weather		cd.WeatherForecastResponse
// 	Config		Config	
// }

// func LoadView() {
// 	// create new view object
// 	view := &MyView{}

// 	// initalize the view
// 	view.Init()

// 	// load the view's template
// 	template := view.Template()

// 	// render the template to an image.Image
// 	renderImage(template)

// }

// // Take in a Template Config and convert into a full Content XML
// func LoadTemplate() {
// 	templateFileName := "nhl_ticker.yaml"


// 	// load the template from the directory
// 	appConfig := GetApplicationConfig()
// 	templateFile := fmt.Sprintf("%v/%v", appConfig.TemplateDir, templateFileName)
// 	data, err := ioutil.ReadFile(*templateFile)
// 	if err != nil {
// 		log.Fatalf("Error loading config: '%v'", err)
// 		return
// 	}

// 	var templateConfig TemplateConfig
// 	if err := yaml.Unmarshal(data, &templateConfig); err != nil {
// 		log.Fatalf("Error unmarshaling template file: '%v'", err)
// 		return
// 	}


// 	tmpl, err := template.New("temp").Parse(templateConfig.Template)
// 	if err != nil {
// 		panic(err)
// 	}

// 	templateConfig = interface {
// 		Games: GetGames(),
// 		Weather: GetWeather(),
// 		Config: appConfig,
// 	} 


// 	var buf bytes.Buffer
// 	err = tmpl.Execute(&buf, templateConfig)
// 	if err != nil {
// 		panic(err)
// 	}	

// 	content := buf.String()

// 	return content

// }

