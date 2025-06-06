package date

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInspectDateFormat(t *testing.T) {
	type expectedT = struct {
		layout DateLayout
		ok     bool
	}
	tests := []struct {
		val        string
		dateFormat DateFormatEnum
		expected   expectedT
	}{
		{
			val: "25/10/1999", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDate1, true},
		},
		{
			val: "25/10/1999", dateFormat: EMPTY_FORMAT,
			expected: expectedT{InDate1, true},
		},
		{
			val: "25-10-1999", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDate2, true},
		},
		{
			val: "25-10-1999", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Blank, false},
		},
		{
			val: "10/25/1999", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Date1, true},
		},
		{
			val: "10-25-1999", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Date2, true},
		},
		{
			val: "1999/10/25", dateFormat: US_DATE_FORMAT,
			expected: expectedT{CommonDate2, true},
		},
		{
			val: "1999-10-25", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{CommonDate1, true},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			layout, ok := InspectDateFormat(test.val, test.dateFormat)
			assert.Equal(t, test.expected, expectedT{layout, ok})
		})
	}
}

func TestInspectDateTimeFormat(t *testing.T) {
	type expectedT = struct {
		layout DateLayout
		ok     bool
	}
	tests := []struct {
		val        string
		dateFormat DateFormatEnum
		expected   expectedT
	}{
		{
			val: "25/10/1999 15:04:05", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDatetime1, true},
		},
		{
			val: "25/10/1999 15:04:05", dateFormat: EMPTY_FORMAT,
			expected: expectedT{InDatetime1, true},
		},
		{
			val: "25-10-1999 15:04:05", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDatetime2, true},
		},
		{
			val: "25-10-1999 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Blank, false},
		},
		{
			val: "10/25/1999 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Datetime1, true},
		},
		{
			val: "10-25-1999 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Datetime2, true},
		},
		{
			val: "1999/10/25 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{CommonDatetime2, true},
		},
		{
			val: "1999-10-25 15:04:05", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{CommonDatetime1, true},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			layout, ok := InspectDateTimeFormat(test.val, test.dateFormat)
			assert.Equal(t, test.expected, expectedT{layout, ok})
		})
	}
}

func TestInspectAllTimeFormat(t *testing.T) {
	type expectedT = struct {
		layout DateLayout
		ok     bool
	}
	tests := []struct {
		val        string
		dateFormat DateFormatEnum
		expected   expectedT
	}{
		{
			val: "25/10/1999", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDate1, true},
		},
		{
			val: "25/10/1999", dateFormat: EMPTY_FORMAT,
			expected: expectedT{InDate1, true},
		},
		{
			val: "25-10-1999", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDate2, true},
		},
		{
			val: "25-10-1999", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Blank, false},
		},
		{
			val: "10/25/1999", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Date1, true},
		},
		{
			val: "10-25-1999", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Date2, true},
		},
		{
			val: "1999/10/25", dateFormat: US_DATE_FORMAT,
			expected: expectedT{CommonDate2, true},
		},
		{
			val: "1999-10-25", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{CommonDate1, true},
		},
		{
			val: "25/10/1999 15:04:05", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDatetime1, true},
		},
		{
			val: "25/10/1999 15:04:05", dateFormat: EMPTY_FORMAT,
			expected: expectedT{InDatetime1, true},
		},
		{
			val: "25-10-1999 15:04:05", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{InDatetime2, true},
		},
		{
			val: "25-10-1999 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Blank, false},
		},
		{
			val: "10/25/1999 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Datetime1, true},
		},
		{
			val: "10-25-1999 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{Datetime2, true},
		},
		{
			val: "1999/10/25 15:04:05", dateFormat: US_DATE_FORMAT,
			expected: expectedT{CommonDatetime2, true},
		},
		{
			val: "1999-10-25 15:04:05", dateFormat: IN_DATE_FORMAT,
			expected: expectedT{CommonDatetime1, true},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			layout, ok := InspectAllTimeFormat(test.val, test.dateFormat)
			assert.Equal(t, test.expected, expectedT{layout, ok})
		})
	}
}

func TestGetAllTimeLayout(t *testing.T) {
	tests := []struct {
		dateFormat DateFormatEnum
		expected   []DateLayout
	}{
		{
			dateFormat: EMPTY_FORMAT, expected: append(inDateLayouts, inDatetimeLayouts...),
		},
		{
			dateFormat: IN_DATE_FORMAT, expected: append(inDateLayouts, inDatetimeLayouts...),
		},
		{
			dateFormat: US_DATE_FORMAT, expected: append(usDateLayouts, usDatetimeLayouts...),
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			layouts := GetAllTimeLayout(test.dateFormat)
			assert.Equal(t, test.expected, layouts)
		})
	}
}
