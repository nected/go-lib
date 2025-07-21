package generators

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/iancoleman/strcase"
)

var tagName = "default"

func GenerateDefaults(t any) {
	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := val.Type().Field(i).Tag.Get(tagName)
		switch field.Kind() {
		case reflect.String:
			if tag != "" {
				field.SetString(tag)
			}
		case reflect.Int:
			if tag != "" {
				v, _ := strconv.ParseInt(tag, 10, 64)
				field.SetInt(v)
			}
		case reflect.Float64:
			if tag != "" {
				v, _ := strconv.ParseFloat(tag, 64)
				field.SetFloat(v)
			}
		case reflect.Int64:
			if tag != "" {
				field.SetInt(proccessInt64(tag))
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

// recursive function to process nested structs
func processStruct(val reflect.Value, dValues map[string]any) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := val.Type().Field(i)
		tag := structField.Tag.Get(tagName)
		dVal := dValues[strcase.ToLowerCamel(structField.Name)]
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
		case reflect.Float64:
			if tag != "" {
				v, _ := strconv.ParseFloat(tag, 64)
				field.SetFloat(v)
			}
			if dVal != nil {
				switch actualDVal := dVal.(type) {
				case float64:
					field.SetFloat(actualDVal)
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
