package generators

import (
	"reflect"
	"strconv"
)

type DefaultTag struct{}

func (t *DefaultTag) GenerateDefaults() {
	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			if tag := val.Type().Field(i).Tag.Get("default"); tag != "" {
				field.SetString(tag)
			}
		}
		if field.Kind() == reflect.Int {
			if tag := val.Type().Field(i).Tag.Get("default"); tag != "" {
				v, _ := strconv.ParseInt(tag, 10, 64)
				field.SetInt(v)
			}
		}
		if field.Kind() == reflect.Struct {
			processStruct(field)
		}
	}
}

// recursive function to process nested structs
func processStruct(field reflect.Value) {
	for i := 0; i < field.NumField(); i++ {
		if field.Field(i).Kind() == reflect.Struct {
			processStruct(field.Field(i))
		}
		if field.Field(i).Kind() == reflect.String {
			if tag := field.Type().Field(i).Tag.Get("default"); tag != "" {
				field.Field(i).SetString(tag)
			}
		}
		if field.Field(i).Kind() == reflect.Int {
			if tag := field.Type().Field(i).Tag.Get("default"); tag != "" {
				v, _ := strconv.ParseInt(tag, 10, 64)
				field.Field(i).SetInt(v)
			}
		}

	}
}
