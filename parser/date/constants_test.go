package date

import (
	"testing"
)

func TestHandleTimeString(t *testing.T) {
	tests := []struct {
		input          string
		format         string
		expectedInput  string
		expectedFormat string
	}{
		{"2023-10-10T10:10:10Z", "yyyy-mm-ddTHH:MM:ssZ", "2023-10-10 10:10:10 UTC", "2006-01-02 15:04:05 MST"},
		{"2023-10-10T10:10:10+0700", "yyyy-mm-ddTHH:MM:ssZ", "2023-10-10 10:10:10+0700", "2006-01-02 15:04:05 MST"},
		{"2023-10-10T10:10:10", "yyyy-mm-ddTHH:MM:ss", "2023-10-10 10:10:10", "2006-01-02 15:04:05"},
		{"", "yyyy-mm-dd", "", ""},
		{"2023-10-10", "", "2023-10-10", ""},
		{"2023-10-10T10:10:10Z", "yyyy-mm-ddTHH:MM:ss", "2023-10-10 10:10:10 UTC", "2006-01-02 15:04:05"},
		{"2023-10-10T10:10:10Z", "yyyy-mm-ddTHH:MM:ssZ", "2023-10-10 10:10:10 UTC", "2006-01-02 15:04:05 MST"},
	}

	for _, test := range tests {
		t.Run(test.input+"_"+test.format, func(t *testing.T) {
			gotInput, gotFormat := handleTimeString(test.input, test.format)
			if gotInput != test.expectedInput || gotFormat != test.expectedFormat {
				t.Errorf("handleTimeString(%q, %q) = (%q, %q); want (%q, %q)", test.input, test.format, gotInput, gotFormat, test.expectedInput, test.expectedFormat)
			}
		})
	}
}
