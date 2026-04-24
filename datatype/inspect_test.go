package datatype

import (
	"fmt"
	"testing"

	"github.com/nected/go-lib/parser/date"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInspectDateType(t *testing.T) {
	tests := []struct {
		name   string
		option date.DateParseOption
		val    any
		want   DataType
	}{
		{name: "nil", val: nil, want: Unknown},

		{name: "int", val: 1, want: Numeric},
		{name: "int8", val: int8(1), want: Numeric},
		{name: "int16", val: int16(1), want: Numeric},
		{name: "int32", val: int32(1), want: Numeric},
		{name: "int64", val: int64(1), want: Numeric},
		{name: "uint", val: uint(1), want: Numeric},
		{name: "uint8", val: uint8(1), want: Numeric},
		{name: "uint16", val: uint16(1), want: Numeric},
		{name: "uint32", val: uint32(1), want: Numeric},
		{name: "uint64", val: uint64(1), want: Numeric},
		{name: "uintptr", val: uintptr(1), want: Numeric},
		{name: "float32", val: float32(1.1), want: Numeric},
		{name: "float64", val: 1.1, want: Numeric},

		{name: "bool true", val: true, want: Boolean},
		{name: "bool false", val: false, want: Boolean},

		{name: "string numeric", val: "1", want: String},
		{name: "string plain", val: "hello world", want: String},
		{name: "string empty", val: "", want: String},

		{name: "string date default", val: "1999/10/25", want: Date},
		{name: "string datetime default", val: "1999/10/25 10:10:10", want: DateTime},
		{name: "string date IN", val: "25/10/1999", want: Date, option: date.DateParseOption{DateFormat: date.IN_DATE_FORMAT}},
		{name: "string datetime IN", val: "25/10/1999 10:10:10", want: DateTime, option: date.DateParseOption{DateFormat: date.IN_DATE_FORMAT}},

		{name: "slice int", val: []int{1, 2}, want: List},
		{name: "slice string", val: []string{"a", "b"}, want: List},
		{name: "slice empty", val: []int{}, want: List},
		{name: "array", val: [3]int{1, 2, 3}, want: List},

		{name: "map string any", val: map[string]any{}, want: JSON},
		{name: "map string string", val: map[string]string{"k": "v"}, want: JSON},

		{name: "primitive ObjectID", val: primitive.NewObjectID(), want: ObjectID},
		{name: "primitive Decimal128", val: primitive.NewDecimal128(10, 10), want: Decimal128},
		{name: "primitive DateTime", val: primitive.DateTime(1700000000000), want: DateTime},
		{name: "primitive Timestamp", val: primitive.Timestamp{T: 1700000000, I: 1}, want: DateTime},
		{name: "primitive D", val: primitive.D{{Key: "k", Value: "v"}}, want: List},

		{name: "unsupported struct", val: struct{ X int }{X: 1}, want: Unknown},
		{name: "unsupported chan", val: make(chan int), want: Unknown},
	}

	for id, test := range tests {
		name := test.name
		if name == "" {
			name = fmt.Sprintf("%v", id)
		}
		t.Run(name, func(t *testing.T) {
			res := InspectDateType(test.option, test.val)
			assert.Equal(t, test.want, res)
		})
	}
}

func TestInferTypeNonStandard(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want DataType
	}{
		{name: "ObjectID", val: primitive.NewObjectID(), want: ObjectID},
		{name: "Decimal128", val: primitive.NewDecimal128(10, 10), want: Decimal128},
		{name: "DateTime", val: primitive.DateTime(1700000000000), want: DateTime},
		{name: "Timestamp", val: primitive.Timestamp{T: 1700000000, I: 1}, want: DateTime},
		{name: "D", val: primitive.D{{Key: "k", Value: "v"}}, want: List},
		{name: "unknown int", val: 1, want: Unknown},
		{name: "unknown string", val: "hello", want: Unknown},
		{name: "unknown nil", val: nil, want: Unknown},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := inferTypeNonStandard(test.val)
			assert.Equal(t, test.want, res)
		})
	}
}
