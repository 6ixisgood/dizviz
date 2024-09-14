package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// StructMetadata the tag metadata related to a specific struct and tag name
type StructMetadata struct {
	TypeName string
	Fields   []FieldMetadata
}

// FieldMetadata the individual key/value pairs defined on a specific field
type FieldMetadata struct {
	FieldName string
	Tags      []TagMetadata
	Type      reflect.Type
	Kind      reflect.Kind
	Value     interface{}
}

// TagMetadata key/value pair of a tag
type TagMetadata struct {
	Key   string
	Value string
}

const (
	viewDefinitionPrefix = "VD"
)

// parseStructMetadata return a StructMetadata containing a struct's tag metadata,
// given a struct and the name of the tag to extract info from
func parseStructMetadata(v interface{}, tagName string) *StructMetadata {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// for each field in the type of struct, create a list of the field's tags
	fieldMetadata := make([]FieldMetadata, 0)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		tagsStr := field.Tag.Get(tagName)
		if tagsStr == "" {
			continue
		}

		// grab each tag's key/value pair (split on "=")
		tagMetadata := make([]TagMetadata, 0)
		tagsStrSplit := strings.Split(tagsStr, ",")
		for _, tag := range tagsStrSplit {
			keyValue := strings.Split(tag, "=")
			if len(keyValue) != 2 {
				// theres no "=" sign
				// just assign the raw value as "root"
				tagMetadata = append(tagMetadata, TagMetadata{
					Key:   "root",
					Value: tag,
				})
				continue
			}

			// get this tag's key/value pair
			tagMetadata = append(tagMetadata, TagMetadata{
				Key:   keyValue[0],
				Value: strings.Trim(keyValue[1], "'"),
			})
		}

		fMeta := FieldMetadata{
			FieldName: typ.Field(i).Name,
			Tags:      tagMetadata,
			Value:     fieldValue.Interface(),
			Kind:      field.Type.Kind(),
			Type:      field.Type,
		}

		// capture the field's tags
		fieldMetadata = append(fieldMetadata, fMeta)
	}

	return &StructMetadata{
		TypeName: reflect.TypeOf(v).Name(),
		Fields:   fieldMetadata,
	}

}

// getBoolTag util function to take a string tag meant to be a bool and convert it
func getBoolTag(tags []TagMetadata, key string) bool {
	for _, tag := range tags {
		if tag.Key == key {
			return tag.Value == "true"
		}
	}
	return false
}

// getIntTag util function to take a string tag meant to be an int and convert it
func getIntTag(tags []TagMetadata, key string) int {
	for _, tag := range tags {
		if tag.Key == key {
			val, err := strconv.Atoi(tag.Value)
			if err == nil {
				return val
			}
		}
	}
	return 0
}

// getStrTag util function to get a string tag's value
func getStrTag(tags []TagMetadata, key string) string {
	for _, tag := range tags {
		if tag.Key == key {
			return tag.Value
		}
	}
	return ""
}

// isEmpty checks if a given reflect.Value is considered "empty" based on its type.
func isEmpty(value interface{}, kind reflect.Kind) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)

	switch kind {
	case reflect.String:
		return val.Len() == 0
	case reflect.Slice, reflect.Array, reflect.Map:
		return val.Len() == 0 || val.IsNil()
	case reflect.Int, reflect.Int64, reflect.Float32, reflect.Float64:
		return val.IsZero()
	case reflect.Bool:
		return !val.Bool()
	default:
		return val.IsZero()
	}
}

// generateFieldSpecs create ViewConfigFieldSpec for each field of a given interface{}
// if it comes across a struct or interface kind, it will recursivly create specs for
// those. If it comes across a slice, it creates specs for each element in the slice
func generateFieldSpecs(v interface{}, parentField string) []ViewConfigFieldSpec {
	if strings.HasPrefix(parentField, "Views.ViewDefinition") {
	}

	metadata := parseStructMetadata(v, "spec")
	jsonMetadata := parseStructMetadata(v, "json")
	specs := make([]ViewConfigFieldSpec, 0)

	if strings.HasPrefix(parentField, "Views.ViewDefinition") {
	}

	for i, field := range metadata.Fields {
		// init the spec for this field
		spec := ViewConfigFieldSpec{
			Field:    parentField + field.FieldName,
			JsonKey:  getStrTag(jsonMetadata.Fields[i].Tags, "root"),
			Type:     field.Type,
			Required: getBoolTag(field.Tags, "required"),
			Min:      getIntTag(field.Tags, "min"),
			Max:      getIntTag(field.Tags, "max"),
			Value:    field.Value,
			Kind:     field.Kind,
			Label:    getStrTag(field.Tags, "label"),
		}
		specs = append(specs, spec)

		if field.Kind == reflect.Struct {
			// Recursive call for nested structs
			nestedSpecs := generateFieldSpecs(field.Value, spec.Field+".")
			specs = append(specs, nestedSpecs...)
		} else if field.Kind == reflect.Slice {
			// Provide generic spec for slice type
			sliceElmTyp := reflect.TypeOf(field.Value).Elem()
			sliceElmInterface := reflect.New(sliceElmTyp).Elem().Interface()
			nestedSpec := generateFieldSpecs(sliceElmInterface, spec.Field+".")
			specs = append(specs, nestedSpec...)

			sliceVal := reflect.ValueOf(field.Value)
			for i := 0; i < sliceVal.Len(); i++ {
				element := sliceVal.Index(i).Interface()
				elementKind := reflect.TypeOf(element).Kind()

				// For struct elements within a slice, recursively determine spec
				if elementKind == reflect.Struct {
					nestedSpecs := generateFieldSpecs(element, fmt.Sprintf("%s[%d].", spec.Field, i))
					specs = append(specs, nestedSpecs...)
				}
			}
		}
	}
	return specs
}

func GenerateViewConfigSpecJson(v ViewConfig) map[string]interface{} {
	specs := generateFieldSpecs(v, "")
	result := make(map[string]interface{})

	for _, spec := range specs {
		path := strings.Split(spec.Field, ".")
		current := result

		for i, p := range path {
			if i == len(path)-1 { // Last item, set the spec
				current[p] = mapViewConfigSpec(spec)
			} else { // Intermediate, ensure map structure
				if _, exists := current[p]; !exists {
					current[p] = make(map[string]interface{})
				}
				// Move into the nested field for the next level, ensuring it exists
				if _, exists := current[p].(map[string]interface{})["nested"]; !exists {
					current[p].(map[string]interface{})["nested"] = make(map[string]interface{})
				}
				current = current[p].(map[string]interface{})["nested"].(map[string]interface{})
			}
		}
	}

	return result
}

func mapViewConfigSpec(spec ViewConfigFieldSpec) map[string]interface{} {
	specMap := map[string]interface{}{
		"label":    spec.Label,
		"key":      spec.JsonKey,
		"kind":     spec.Kind.String(),
		"type":     spec.Type.String(),
		"required": spec.Required,
		"min":      spec.Min,
		"max":      spec.Max,
	}
	return specMap
}

// SaveViewDefinition saves a view definition to the store
func SaveViewDefinition(definition ViewDefinition) error {
	data, err := json.Marshal(definition)
	if err != nil {
		return err
	}
	_, err = CommonConfig.Store.SaveItem(viewDefinitionPrefix, definition.Id, data)
	return err
}

// GetViewConfig retrieves a view definition from the store by name
func GetViewDefinition(id string) (ViewDefinition, error) {
	data, err := CommonConfig.Store.GetItem(viewDefinitionPrefix + "-" + id)

	if err != nil {
		return ViewDefinition{}, err
	}

	vd, err := unmarshalViewDefinition(data)
	return vd, err
}

// GetAllViewDefinitions retrieves all saved view definitions from the store
func GetAllViewDefinitions() ([]ViewDefinition, error) {
	var definitions []ViewDefinition
	viewDatas, err := CommonConfig.Store.GetPrefix(viewDefinitionPrefix + "-")
	if err != nil {
		return nil, err
	}

	for _, viewData := range viewDatas {
		definition, err := unmarshalViewDefinition(viewData)
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

// DeleteViewDefinition delete a given view def by id from the store
func DeleteViewDefinition(id string) error {
	return CommonConfig.Store.DeleteItem(viewDefinitionPrefix + "-" + id)
}

// unmarshalViewDefinition helper function to get a ViewDefinition from []byte
func unmarshalViewDefinition(data []byte) (ViewDefinition, error) {
	var definition ViewDefinition

	// Unmarshal the whole object
	err := json.Unmarshal(data, &definition)
	if err != nil {
		return definition, err
	}

	// Unmarshal into a map to get the raw config
	var rawDef map[string]json.RawMessage
	err = json.Unmarshal(data, &rawDef)
	if err != nil {
		fmt.Println("HERE 1")
		return definition, err
	}

	// Get the RegisteredView for this type
	regView, ok := RegisteredViews[definition.Type]
	if !ok {
		fmt.Println("HERE 2")
		return definition, errors.New(fmt.Sprintf("No registered view of type %s", definition.Type))
	}

	// Use the RegisteredView's NewConfig function to unmarshal the config
	definition.Config = regView.NewConfig()
	err = json.Unmarshal(rawDef["config"], &definition.Config)
	if err != nil {
		fmt.Println("HERE 3")
		return definition, err
	}

	return definition, nil
}
