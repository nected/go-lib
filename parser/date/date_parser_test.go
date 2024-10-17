package date

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		format   string
		expected time.Time
		hasError bool
	}{
		{"01/Dec/2023 01:01:01:01", "dd/mmm/yyyy hh:MM:ss:ns", time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), false},
		{"2023-10-01", "2006-01-02", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"01/10/2023", "02/01/2006", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"2023-10-01T15:04:05Z", time.RFC3339, time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), false},
		{"invalid-date", "2006-01-02", time.Time{}, true},
		{"01/10/2023", "dd/mm/yyyy", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"01/Jan/2023", "dd/mmm/yyyy", time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Feb/2023", "dd/mmm/yyyy", time.Date(2023, 02, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Mar/2023", "dd/mmm/yyyy", time.Date(2023, 03, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Apr/2023", "dd/mmm/yyyy", time.Date(2023, 04, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/May/2023", "dd/mmm/yyyy", time.Date(2023, 05, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Jun/2023", "dd/mmm/yyyy", time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Jul/2023", "dd/mmm/yyyy", time.Date(2023, 07, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Aug/2023", "dd/mmm/yyyy", time.Date(2023, 8, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Sep/2023", "dd/mmm/yyyy", time.Date(2023, 9, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Oct/2023", "dd/mmm/yyyy", time.Date(2023, 10, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Nov/2023", "dd/mmm/yyyy", time.Date(2023, 11, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Dec/2023", "dd/mmm/yyyy", time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC), false},
		{"01/Jan/2023", "dd/MMM/yyyy", time.Time{}, true},
		{"01/Jan/2023", "dd/mmm/yyya", time.Time{}, true},
		{"01/Jab/2023", "dd/mmm/yyya", time.Time{}, true},
		{"", "dd/MMM/yyyy", time.Time{}, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := Parse(test.input, test.format)
			if (err != nil) != test.hasError {
				t.Errorf("expected error: %v, got: %v", test.hasError, err)
			}
			if !result.Equal(test.expected) {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
func TestFormat(t *testing.T) {
	tests := []struct {
		time     time.Time
		format   string
		expected string
		hasError bool
	}{
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy hh:MM:ss:ns", "01/Dec/2023 01:01:01:1", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmmm/yyyy hh:MM:ss:ns", "01/December/2023 01:01:01:1", false},
		{time.Date(2023, 12, 01, 13, 01, 01, 01, time.UTC), "dd/mmm/yyyy h:MM:ss:ns", "01/Dec/2023 1:01:01:1", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy h:MM:ss:ns", "01/Dec/2023 1:01:01:1", false},
		{time.Date(2023, 12, 01, 13, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:MM:ss:ns", "01/Dec/2023 13:01:01:1", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:MM:ss:ns", "01/Dec/2023 1:01:01:1", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy HH:MM:ss:ns", "01/Dec/2023 01:01:01:1", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:M:s:ns", "01/Dec/2023 1:1:1:1", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:M:s:na", "", true},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:M:sa:na", "", true},

		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "2006-01-02", "2023-10-01", false},
		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "02/01/2006", "01/10/2023", false},
		{time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), time.RFC3339, "2023-10-01T15:04:05Z", false},
		{time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), "", "2023-10-01T15:04:05Z", false},

		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "dd/mm/yyyy", "01/10/2023", false},
		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "d/mm/yyyy", "1/10/2023", false},
		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "ddd/mm/yyyy", "Sun/10/2023", false},
		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "dddd/mm/yyyy", "Sunday/10/2023", false},
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "dddd/m/yyyy", "Sunday/1/2023", false},
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "dddd/mm/yyyy", "Sunday/01/2023", false},
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "dddd/mm/yy", "Sunday/01/23", false},
		{time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC), "dddd/mm/yy", "Wednesday/01/03", false},
		{time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC), "dddd/mm/yyyy", "Wednesday/01/2003", false},
		{time.Date(0003, 1, 1, 0, 0, 0, 0, time.UTC), "dddd/mm/yyyy", "Wednesday/01/0003", false},

		{time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Jan/2023", false},
		{time.Date(2023, 02, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Feb/2023", false},
		{time.Date(2023, 03, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Mar/2023", false},
		{time.Date(2023, 04, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Apr/2023", false},
		{time.Date(2023, 05, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/May/2023", false},
		{time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Jun/2023", false},
		{time.Date(2023, 07, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Jul/2023", false},
		{time.Date(2023, 8, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Aug/2023", false},
		{time.Date(2023, 9, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Sep/2023", false},
		{time.Date(2023, 10, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Oct/2023", false},
		{time.Date(2023, 11, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Nov/2023", false},
		{time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC), "dd/mmm/yyyy", "01/Dec/2023", false},
		{time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC), "dd/ mmm/ yyyy", "01/ Dec/ 2023", false},
		{time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC), "dd mmm yyyy", "01 Dec 2023", false},
		{time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC), "dd mmm yyyy", "01 Dec 2023", false},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result, err := Format(test.time, test.format)
			if (err != nil) != test.hasError {
				t.Errorf("expected error: %v, got: %v", test.hasError, err)
				return
			}
			if result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
