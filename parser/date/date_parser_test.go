package date

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	customTimeZone, _ := time.LoadLocation("Asia/Kolkata")
	// pstTimeZone, _ := time.LoadLocation("America/Los_Angeles")
	tests := []struct {
		input    string
		format   string
		expected time.Time
		hasError bool
	}{
		{"01/Dec/2023 01:01:01", "dd/mmm/yyyy hh:MM:ss", time.Date(2023, 12, 01, 01, 01, 01, 00, time.UTC), false},
		{"01/Dec/2023 a 01:01:01", "dd/mmm/yyyy \\a hh:MM:ss", time.Date(2023, 12, 01, 01, 01, 01, 00, time.UTC), false},
		{"2024-10-17 12:16:00", "yyyy-mm-dd hh:MM:ss", time.Date(2024, 10, 17, 12, 16, 00, 00, time.UTC), false},
		{"April 05, 2023 +05:30", "mmmm dd, yyyy p", time.Date(2023, 04, 05, 00, 00, 00, 00, customTimeZone), false},
		{"April 05, 2023", "mmmm dd, yyyy", time.Date(2023, 04, 05, 00, 00, 00, 00, time.UTC), false},
		{"April 05, 2023 +0530", "mmmm dd, yyyy o", time.Date(2023, 04, 05, 00, 00, 00, 00, customTimeZone), false},
		{"Apr 05, 2023 +05:30", "mmm dd, yyyy p", time.Date(2023, 04, 05, 00, 00, 00, 00, customTimeZone), false},
		{"2023-04-05 11:30 AM", "yyyy-mm-dd hh:MM AA", time.Date(2023, 04, 05, 11, 30, 00, 00, time.UTC), false},
		{"2023-04-05 18:30", "yyyy-mm-dd HH:MM", time.Date(2023, 04, 05, 18, 30, 00, 00, time.UTC), false},
		{"2023-04-05 06:30 am", "yyyy-mm-dd HH:MM aa", time.Date(2023, 04, 05, 06, 30, 00, 00, time.UTC), false},

		{"2023-04-05 12:30:45.1231", "yyyy-mm-dd HH:MM:ss", time.Date(2023, 04, 05, 12, 30, 45, 123100000, time.UTC), false},
		{"23-4-5 12:05:45.1231", "yy-m-d HH:MM:ss", time.Date(2023, 04, 05, 12, 5, 45, 123100000, time.UTC), false},
		{"2023-04-05 12:30:45 PST", "yyyy-mm-dd HH:MM:ss Z", time.Date(2023, 04, 05, 12, 30, 45, 0, time.UTC), false}, // to check

		{"2024-10-17 12:16:00", "yyyy-mm-dd hh:MM:ss", time.Date(2024, 10, 17, 12, 16, 0, 0, time.UTC), false},
		{"2023/14/32", "yyyy/mm/dd", time.Time{}, true},
		{"2023/01/15", "yyyy-mm-dd", time.Time{}, true},
		{"2023/01-15", "yyyy/mm-dd", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"00-00-0000", "mm/dd/yyyy", time.Time{}, true},
		{"00-00-0000", "mm-dd-yyyy", time.Time{}, true},

		{"2023-10-01", "2006-01-02", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"2023-10-01", "2006-01-02", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},

		{"01/10/2023", "02/01/2006", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"2023-10-01T15:04:05Z", "", time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), false},
		{"2023-10-01T15:04:05+05:30", "", time.Date(2023, 10, 1, 15, 4, 5, 0, customTimeZone), false},
		{"20230405", "yyyymmdd", time.Time{}, true},
		{"04/05", "dd/mm", time.Date(0, 5, 4, 0, 0, 0, 0, time.UTC), false},
		{"05-04-23", "dd-mm-yy", time.Date(2023, 4, 5, 0, 0, 0, 0, time.UTC), false},
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
		{"2023-10-01 PM 03:04:05 UTC", "yyyy-mm-dd A hh:MM:ss Z", time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), false},
		{"2023-10-01T03:04:05Z", "yyyy-mm-ddThh:MM:ssZ", time.Date(2023, 10, 1, 03, 4, 5, 0, time.UTC), false},
		{"2023-10-01T03:04:05 UTC", "yyyy-mm-ddThh:MM:ssZ", time.Date(2023, 10, 1, 03, 4, 5, 0, time.UTC), false},
		{"2023-10-01T03:04:05Z", "yyyy-mm-ddThh:MM:ssZ", time.Date(2023, 10, 1, 03, 4, 5, 0, time.UTC), false},
		{"2023-10-01 T 03:04:05Z", "yyyy-mm-dd T hh:MM:ssZ", time.Date(2023, 10, 1, 03, 4, 5, 0, time.UTC), false},
		{"01-10-2023T03:04:05Z", "dd-mm-yyyyThh:MM:ssZ", time.Date(2023, 10, 1, 03, 4, 5, 0, time.UTC), false},
		{"01-10-2023T03:04:05 +0000", "dd-mm-yyyyThh:MM:ss o", time.Date(2023, 10, 1, 03, 4, 5, 0, time.UTC), false},
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
	customTimeZone, _ := time.LoadLocation("Asia/Kolkata")
	pstTimeZone, _ := time.LoadLocation("America/Los_Angeles")

	tests := []struct {
		time     time.Time
		format   string
		expected string
		hasError bool
	}{
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy hh:MM:ss", "01/Dec/2023 01:01:01", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/y hh:MM:ss", "01/Dec/23 01:01:01", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 99, time.UTC), "dd/mmm/y hh:MM:ss.ns", "01/Dec/23 01:01:01.000000099", false},

		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmmm/yyyy hh:MM:ss", "01/December/2023 01:01:01", false},
		{time.Date(2023, 12, 01, 13, 01, 01, 01, time.UTC), "dd/mmm/yyyy h:MM:ss", "01/Dec/2023 1:01:01", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy h:MM:ss", "01/Dec/2023 1:01:01", false},
		{time.Date(2023, 12, 01, 13, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:MM:ss", "01/Dec/2023 13:01:01", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:MM:ss", "01/Dec/2023 1:01:01", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy HH:MM:ss", "01/Dec/2023 01:01:01", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:M:s", "01/Dec/2023 1:1:1", false},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:M:s:na", "", true},
		{time.Date(2023, 12, 01, 01, 01, 01, 01, time.UTC), "dd/mmm/yyyy H:M:sa:na", "", true},

		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "2006-01-02", "2023-10-01", false},
		{time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), "02/01/2006", "01/10/2023", false},
		{time.Date(2023, 10, 1, 15, 4, 5, 0, customTimeZone), time.RFC3339, "2023-10-01T15:04:05+05:30", false},
		{time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), "hh TT", "03 PM", false},
		{time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), "hh t", "03 p", false},
		{time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), "hh T", "03 P", false},

		{time.Date(2023, 10, 1, 5, 4, 5, 0, time.UTC), "hh TT", "05 AM", false},
		{time.Date(2023, 10, 1, 5, 4, 5, 0, time.UTC), "hh tt", "05 am", false},
		{time.Date(2023, 10, 1, 17, 4, 5, 0, time.UTC), "hh tt", "05 pm", false},
		{time.Date(2023, 10, 1, 17, 4, 5, 0, time.UTC), "hh tt Z", "05 pm UTC", false},
		{time.Date(2023, 10, 1, 17, 4, 5, 0, customTimeZone), "hh tt Z", "05 pm Asia/Kolkata", false},
		// {time.Date(2023, 10, 1, 17, 4, 5, 0, customTimeZone), "hh tt o", "05 pm 19800", false}, // to check
		{time.Date(2023, 10, 1, 17, 4, 5, 0, pstTimeZone), "hh tt o", "05 pm -0700", false},     // to check
		{time.Date(2023, 10, 1, 17, 4, 5, 0, customTimeZone), "hh tt o", "05 pm +0530", false},  // to check
		{time.Date(2023, 10, 1, 17, 4, 5, 0, customTimeZone), "hh tt p", "05 pm +05:30", false}, // to check
		{time.Date(2023, 10, 1, 17, 4, 5, 0, pstTimeZone), "hh tt p", "05 pm -07:00", false},    // to check

		// {time.Date(2023, 10, 1, 17, 4, 5, 0, customTimeZone), "hh tt p", "05 pm 5", false}, // to check

		{time.Date(2023, 10, 1, 15, 4, 5, 0, time.UTC), "HH TT", "15 PM", false},
		{time.Date(2023, 10, 1, 5, 4, 5, 0, time.UTC), "HH TT", "05 AM", false},

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
		{time.Date(2024, 10, 17, 0, 0, 0, 0, time.UTC), "DDD", "Thu", false},
		{time.Date(2024, 10, 17, 0, 0, 0, 0, time.UTC), "DDDD", "Thursday", false},
		{time.Date(2004, 10, 17, 0, 0, 0, 111111111111, time.UTC), "l", "111", false},
		{time.Date(2004, 10, 17, 0, 0, 0, 111111111111, time.UTC), "L", "11", false},
		// {time.Date(2004, 10, 17, 0, 0, 0, 111111111111, customTimeZone), "o", "11", false}, // to check
		{time.Date(2004, 10, 17, 0, 0, 0, 111111111111, time.UTC), "t", "a", false},
		{time.Date(2004, 10, 17, 0, 0, 0, 111111111111, time.UTC), "T", "A", false},
		{time.Date(2004, 10, 17, 06, 5, 0, 0, time.UTC), "hh:MM TT", "06:05 AM", false},
		{time.Date(2004, 10, 17, 18, 5, 0, 0, time.UTC), "hh:MM TT", "06:05 PM", false},
		{time.Date(2004, 10, 17, 8, 5, 0, 0, time.UTC), "g:MM TT", "8:05 AM", false},
		{time.Date(2007, 06, 9, 17, 46, 21, 0, time.UTC), "HH:MM:ss", "17:46:21", false},
		{time.Date(2007, 06, 9, 17, 46, 21, 0, time.UTC), time.RFC3339, "2007-06-09T17:46:21Z", false},
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
func TestFormatMap2Format(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"dd/mm/yyyy", "d/m/Y"},
		{"dd/mmm/yyyy", "d/M/Y"},
		{"yyyy-mm-dd", "Y-m-d"},
		{"hh:MM:ss", "h:i:s"},
		{"hh AA", "h A"},
		{"HH AA", "H A"},
		{"dd/MMM/yyyy", "d/MMM/Y"},
		{"dd/mmm/yyya", "d/M/yyya"},
		{"dd/mmm/yyyy hh:MM:ss.ns", "d/M/Y h:i:s.ns"},
		{"dd/mmmm/yyyy hh:MM:ss.ns", "d/F/Y h:i:s.ns"},
		{"yyyy-mm dd, A hh:MM:ss Z07:00", "Y-m d, A h:i:s Z07:00"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := formatMap2format(test.input)
			if result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
func TestFormat2Layout(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"dd/mm/yyyy", "02/01/2006"},
		{"dd/mmm/yyyy", "02/Jan/2006"},
		{"yyyy-mm-dd", "2006-01-02"},
		{"hh:MM:ss", "03:04:05"},
		{"hh AA", "03 PM"},
		{"HH AA", "15 PM"},
		{"dd/MMM/yyyy", "02/JanJanJan/2006"},
		{"dd/mmm/yyya", "02/Jan/060606pm"},
		{"dd/mmm/yyyy hh:MM:ss", "02/Jan/2006 03:04:05"},
		{"dd/mmmm/yyyy hh:MM:ss", "02/January/2006 03:04:05"},
		{"yyyy-mm dd, A hh:MM:ss -07:00", "2006-01 02, PM 03:04:05 -07:00"},
		{"mmmm dd, yyyy +05:30", "January 02, 2006 +05:30"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := convertToGoFormat(test.input)
			if result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
