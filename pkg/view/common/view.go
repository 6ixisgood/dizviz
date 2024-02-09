package common

import (
	"log"
	"fmt"
	"maps"
	"bytes"
	"encoding/xml"
	"html/template"
	compCommon "github.com/6ixisgood/matrix-ticker/pkg/component/common"
)

// define the View interface
type View interface {
	Init()
	Template()				*compCommon.Template
	SetTemplate(*compCommon.Template)
	TemplateString()		string
	TemplateData()			map[string]interface{}
	Stop()
}

type ViewCommonConfig struct {
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

type ViewConfig interface {}

type RegisteredView struct {
	NewConfig		func() ViewConfig
	NewView			func(ViewConfig) (View, error)
}

var (
	CommonConfig = &ViewCommonConfig{}
	RegisteredViews = map[string]RegisteredView{}
)

func SetViewCommonConfig(config *ViewCommonConfig) {
	CommonConfig = config
}

func RegisterView(name string, creator RegisteredView) {
	RegisteredViews[name] = creator
}

func TemplateRefresh(v View) {
	// create the template object
	tmpl := template.New("view-template")

	// gather all the custom functions
	funcMap := template.FuncMap{
		"NilOrDefault": func() string { return "N/A" },
		"CardinalToOrdinal": func(card int) string {
			switch card {
			case 0:
				return "0th"
			case 1:
				return "1st"
			case 2:
				return "2nd"
			case 3:
				return "3rd"
			case 4:
				return "4th"
			default:
				return "N/A"
			}	
		},
	}
	tmpl = tmpl.Funcs(funcMap)

	// construct the template string
	tmplString := `
		{{ $MatrixSizex := .Ctx.MatrixCols }}
		{{ $MatrixSizey := .Ctx.MatrixRows }}
		{{ $DefaultImageSizex := .Ctx.DefaultImageSizeX }}
		{{ $DefaultImageSizey := .Ctx.DefaultImageSizeY }}
		{{ $DefaultFontSize := .Ctx.DefaultFontSize }}
		{{ $DefaultFontType := .Ctx.DefaultFontType }}
		{{ $DefaultFontStyle := .Ctx.DefaultFontStyle }}
		{{ $DefaultFontColor := .Ctx.DefaultFontColor }}
		{{ $ImageDir := .Ctx.ImageDir }}
		{{ $CacheDir := .Ctx.CacheDir }}

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
		"Ctx": CommonConfig,
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
	t := compCommon.Template{}
	err = xml.Unmarshal([]byte(tmplStr), &t)
	if err != nil {
		log.Fatalf("Unable to unmarshal xml content: '%v'", err)
	}

	// set new template and init
	v.SetTemplate(&t)
	t.Init()
}