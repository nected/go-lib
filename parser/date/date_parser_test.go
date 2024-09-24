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
		{"2023-10-01", "2006-01-02", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"01/10/2023", "02/01/2006", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"2023-10-01T15:04:05Z", time.RFC3339, time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), false},
		{"invalid-date", "2006-01-02", time.Time{}, true},
		{"01/10/2023", "dd/MM/yyyy", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"01/Jan/2023", "dd/mmm/yyyy", time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC), false},
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
