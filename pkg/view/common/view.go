package common

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	compCommon "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"html/template"
	"log"
	"maps"
	"reflect"
	"time"
)

// View a structure to describe a layout of components at a given time
type View interface {
	Init()
	Template() *compCommon.Template
	SetTemplate(*compCommon.Template)
	TemplateString() string
	TemplateData() map[string]interface{}
	Stop()
}

// ViewCommonConfig a set of global application configuration useful for rendering Views
type ViewCommonConfig struct {
	MatrixRows        int
	MatrixCols        int
	ImageDir          string
	CacheDir          string
	DefaultImageSizeX int
	DefaultImageSizeY int
	DefaultFontSize   int
	DefaultFontColor  string
	DefaultFontStyle  string
	DefaultFontType   string
}

// ViewConfig type alias to hold raw config definition for a View
type ViewConfig interface{}

// ViewDefinition what defines a View? The Type of View it is and the View's configuration
type ViewDefinition struct {
    Id string	`json:"id" spec:"label='View Definition ID',required='true'"`
	Name string	`json:"name" spec:"label='View Definition Name',required='true'"`
	Type   string     `json:"type" spec:"label='View Type',required='true'"`
	Config ViewConfig `json:"config" spec:"label='View Config',required='true'"`
}

// ViewDefinitionRaw similar to ViewDefinition, but Config is json.RawMessage ([]byte)
type ViewDefinitionRaw struct {
	 Id string `json:"id"`
	Name string `json:"name"`
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

// RegisteredView set of generic functions to a create a given View's config and the View itself
type RegisteredView struct {
	NewConfig func() ViewConfig
	NewView   func(ViewConfig) (View, error)
}

var (
	CommonConfig    = &ViewCommonConfig{}
	RegisteredViews = map[string]RegisteredView{}
)

func SetViewCommonConfig(config *ViewCommonConfig) {
	CommonConfig = config
}

func RegisterView(name string, creator RegisteredView) {
	RegisteredViews[name] = creator
}

// TemplateRefresh static function to generate a View's template
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

type ViewConfigFieldSpec struct {
	Field    string
	JsonKey  string
	Type     reflect.Type
	Value    interface{}
	Kind     reflect.Kind
	Required bool
	Min      int
	Max      int
	Label    string
}

// ValidateViewConfig takes a ViewConfig and using the rules defined by
// "spec" field tags, it validates the configuration
func ValidateViewConfig(v ViewConfig) error {
	validators := generateFieldSpecs(v, "")

	for _, validator := range validators {
		if validator.Required && isEmpty(validator.Value, validator.Kind) {
			return fmt.Errorf("field %s is required", validator.Field)
		}
		switch validator.Kind {
		case reflect.Int, reflect.Int64:
			switch validator.Type {
			case reflect.TypeOf(time.Duration(0)):
				dur, ok := validator.Value.(time.Duration)
				if !ok {
					return fmt.Errorf("error converting field %s of type %s to time.Duration", validator.Field, validator.Type)
				}
				if time.Duration(validator.Min) != 0 && dur < time.Duration(validator.Min) {
					return fmt.Errorf("field %s must be at least %d", validator.Field, validator.Min)
				}
				if time.Duration(validator.Max) != 0 && dur > time.Duration(validator.Max) {
					return fmt.Errorf("field %s must be no more than %d", validator.Field, validator.Max)
				}
			default:
				num, ok := validator.Value.(int)
				if !ok {
					return fmt.Errorf("error converting field %s of type %s to int", validator.Field, validator.Type)
				}
				if validator.Min != 0 && num < validator.Min {
					return fmt.Errorf("field %s must be at least %d", validator.Field, validator.Min)
				}
				if validator.Max != 0 && num > validator.Max {
					return fmt.Errorf("field %s must be no more than %d", validator.Field, validator.Max)
				}
			}

		case reflect.String:
			str, ok := validator.Value.(string)
			if !ok {
				return fmt.Errorf("error converting field %s of type %s to string", validator.Field, validator.Type)
			}
			if validator.Min != 0 && len(str) < validator.Min {
				return fmt.Errorf("field %s must have at least %d characters", validator.Field, validator.Min)
			}
			if validator.Max != 0 && len(str) > validator.Max {
				return fmt.Errorf("field %s must have no more than %d characters", validator.Field, validator.Max)
			}
		case reflect.Slice:
			sliceVal := reflect.ValueOf(validator.Value)
			if sliceVal.Kind() != reflect.Slice {
				return fmt.Errorf("field %s is not a slice", validator.Field)
			}
			if validator.Min != 0 && sliceVal.Len() < validator.Min {
				return fmt.Errorf("field %s must contain at least %d items", validator.Field, validator.Min)
			}
			if validator.Max != 0 && sliceVal.Len() > validator.Max {
				return fmt.Errorf("field %s must contain no more than %d items", validator.Field, validator.Max)
			}
		}
	}

	return nil
}
