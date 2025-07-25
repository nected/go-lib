package generators

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

var tagName = "default"

func GenerateDefaults(t any) {
	val := reflect.ValueOf(t).Elem()
	if val.IsValid() && val.Kind() == reflect.Struct {
		processStruct(val, make(map[string]any))
	}
}

// recursive function to process nested structs
func processStruct(val reflect.Value, structDefaultValues map[string]any) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := val.Type().Field(i)
		tag := structField.Tag.Get(tagName)
		dVal := structDefaultValues[strcase.ToLowerCamel(structField.Name)]
		switch field.Kind() {
		case reflect.String:
			if tag != "" {
				field.SetString(tag)
			}
			if dVal != nil {
				switch actualDVal := dVal.(type) {
				case string:
					field.SetString(actualDVal)
				}
			}
		case reflect.Bool:
			if tag != "" {
				field.SetBool(strings.ToLower(tag) == "true")
			}
			if dVal != nil {
				switch actualDVal := dVal.(type) {
				case bool:
					field.SetBool(actualDVal)
				}
			}

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint:
			if tag != "" {
				field.SetInt(proccessInt64(tag))
			}
			if dVal != nil {
				switch actualDVal := dVal.(type) {
				case int8:
					field.SetInt(int64(actualDVal))
				case int16:
					field.SetInt(int64(actualDVal))
				case int32:
					field.SetInt(int64(actualDVal))
				case uint:
					field.SetInt(int64(actualDVal))
				}
			}

		case reflect.Int:
			if tag != "" {
				field.SetInt(proccessInt64(tag))
			}
			if dVal != nil {
				switch actualDVal := dVal.(type) {
				case int:
					field.SetInt(int64(actualDVal))
				case float64:
					field.SetInt(int64(actualDVal))
				}
			}

		case reflect.Int64:
			if tag != "" {
				field.SetInt(proccessInt64(tag))
			}
			if dVal != nil {
				switch actualDVal := dVal.(type) {
				case int:
					field.SetInt(int64(actualDVal))
				case int64:
					field.SetInt(int64(actualDVal))
				case float64:
					field.SetInt(int64(actualDVal))
				case string:
					d, _ := time.ParseDuration(actualDVal)
					field.SetInt(d.Nanoseconds())
				}
			}
		case reflect.Float64, reflect.Float32:
			if tag != "" {
				v, _ := strconv.ParseFloat(tag, 64)
				field.SetFloat(v)
			}
			if dVal != nil {
				switch actualDVal := dVal.(type) {
				case float64:
					field.SetFloat(actualDVal)
				case float32:
					field.SetFloat(float64(actualDVal))
				}
			}
		case reflect.Struct:
			dValue := make(map[string]any)
			if tag != "" {
				_ = json.Unmarshal([]byte(tag), &dValue)
			}
			processStruct(field, dValue)
		}
	}
}

func proccessInt64(tag string) int64 {
	d, err := time.ParseDuration(tag)
	if err == nil {
		return d.Nanoseconds()
	}
	v, _ := strconv.ParseInt(tag, 10, 64)
	return v
}
