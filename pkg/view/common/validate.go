package common

import (
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

// ParseStructMetadata return a StructMetadata containing a struct's tag metadata,
// given a struct and the name of the tag to extract info from
func ParseStructMetadata(v interface{}, tagName string) StructMetadata {
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

	return StructMetadata{
		TypeName: reflect.TypeOf(v).Name(),
		Fields:   fieldMetadata,
	}

}

// ViewConfigFieldValidator holds info to make validating View Configs easier
type ViewConfigFieldValidator struct {
	Field    string
	Type     reflect.Type
	Required bool
	Min      int
	Max      int
	Value    interface{}
	Kind     reflect.Kind
}

// ValidateViewConfig takes a ViewConfig and using the rules defined by
// "spec" field tags, it validates the configuration
func ValidateViewConfig(v ViewConfig) error {
	validators := generateFieldValidators(v, "")

	for _, validator := range validators {
		if validator.Required && isEmpty(validator.Value, validator.Kind) {
			return fmt.Errorf("field %s is required", validator.Field)
		}

		switch validator.Kind {
		case reflect.Int, reflect.Int64:
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

// generateFieldValidators create ViewConfigValidators for each field of a given interface{}
// if it comes across a struct or interface kind, it will recursivly create validators for
// those. If it comes across a slice, it creates validators for each element in the slice
func generateFieldValidators(v interface{}, parentField string) []ViewConfigFieldValidator {
	metadata := ParseStructMetadata(v, "spec")
	validators := make([]ViewConfigFieldValidator, 0)

	for _, field := range metadata.Fields {
		// init the validator for this field
		validator := ViewConfigFieldValidator{
			Field:    parentField + field.FieldName,
			Type:     field.Type,
			Required: getBoolTag(field.Tags, "required"),
			Min:      getIntTag(field.Tags, "min"),
			Max:      getIntTag(field.Tags, "max"),
			Value:    field.Value,
			Kind:     field.Kind,
		}
		validators = append(validators, validator)

		if field.Kind == reflect.Struct || field.Kind == reflect.Interface {
			// Recursive call for nested structs
			nestedValidators := generateFieldValidators(field.Value, validator.Field+".")
			validators = append(validators, nestedValidators...)
		} else if field.Kind == reflect.Slice {
			sliceVal := reflect.ValueOf(field.Value)
			for i := 0; i < sliceVal.Len(); i++ {
				element := sliceVal.Index(i).Interface()
				elementKind := reflect.TypeOf(element).Kind()

				// For struct elements within a slice, recursively validate
				if elementKind == reflect.Struct {
					nestedValidators := generateFieldValidators(element, fmt.Sprintf("%s[%d].", validator.Field, i))
					validators = append(validators, nestedValidators...)
				}
			}
		}
	}
	return validators
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

// getBoolTag util function to take a string tag meant to be an int and convert it
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
