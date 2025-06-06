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
		option date.DateParseOption
		val    any
		want   DataType
	}{
		{val: 1, want: Numeric},
		{val: 1.1, want: Numeric},
		{val: "1", want: String},
		{val: "1999/10/25", want: Date},
		{val: "1999/10/25 10:10:10", want: DateTime},
		{val: "25/10/1999", want: Date, option: date.DateParseOption{DateFormat: date.IN_DATE_FORMAT}},
		{val: "25/10/1999 10:10:10", want: DateTime, option: date.DateParseOption{DateFormat: date.IN_DATE_FORMAT}},
		{val: []int{1, 2}, want: List},
		{val: map[string]any{}, want: JSON},
		{val: true, want: Boolean},
		{val: primitive.NewObjectID(), want: ObjectID},
		{val: primitive.NewDecimal128(10, 10), want: Decimal128},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			res := InspectDateType(test.option, test.val)
			assert.Equal(t, test.want, res)
		})
	}
}
