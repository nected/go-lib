package datatype

import (
	"reflect"

	"github.com/nected/go-lib/parser/date"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataType string

const (
	Unknown    DataType = ""
	Numeric    DataType = "numeric"
	Boolean    DataType = "boolean"
	String     DataType = "string"
	DateTime   DataType = "dateTime"
	Date       DataType = "date"
	JSON       DataType = "json"
	List       DataType = "list"
	ObjectID   DataType = "objectid"
	Decimal128 DataType = "decimal128"
)

func InspectDateType(option date.DateParseOption, v any) DataType {
	if v == nil {
		return Unknown
	}

	varType := reflect.TypeOf(v)
	switch varType.Kind() {
	case reflect.String:
		// check if v is of any of the supported Date formats
		_, ok := date.InspectDateFormat(v.(string), option.DateFormat)
		if ok {
			return Date
		}

		// check if v is of any of the supported DateTime formats
		_, ok = date.InspectDateTimeFormat(v.(string), option.DateFormat)
		if ok {
			return DateTime
		}
		return String

		/*
			epoch timestamps will also come as Numeric, and ideally it should work
			as mostly numeric operations would be done on that, any date manipulations
			would ideally be happening only when the user is converting the epoch in
			a date string
		*/
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return Numeric

	case reflect.Int64:
		switch varType.String() {
		case "primitive.DateTime", "primitive.Timestamp":
			return inferTypeNonStandard(v)
		}
		return Numeric

	case reflect.Bool:
		return Boolean

	case reflect.Slice, reflect.Array:
		inferredType := inferTypeNonStandard(v)
		if inferredType == Unknown {
			inferredType = List
		}
		return inferredType
	case reflect.Map:
		return JSON
	default:
		return inferTypeNonStandard(v)
	}
}

func inferTypeNonStandard(v interface{}) DataType {
	switch v.(type) {
	case primitive.ObjectID:
		return ObjectID
	case primitive.Decimal128:
		return Decimal128
	case primitive.DateTime, primitive.Timestamp:
		return DateTime
	case primitive.D:
		return List
	default:
		return Unknown
	}
}
