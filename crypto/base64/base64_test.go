package base64

import (
	"testing"
)

func TestB64Encode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "aGVsbG8="},
		{[]byte("world"), "d29ybGQ="},
		{[]byte(""), ""},
		{[]byte("base64 encoding in Go"), "YmFzZTY0IGVuY29kaW5nIGluIEdv"},
	}

	for _, test := range tests {
		result := B64Encode(test.input)
		if result != test.expected {
			t.Errorf("B64Encode(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
func TestB64Decode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"aGVsbG8=", "hello", false},
		{"d29ybGQ=", "world", false},
		{"", "", false},
		{"YmFzZTY0IGVuY29kaW5nIGluIEdv", "base64 encoding in Go", false},
		{"invalid_base64", "", true},
	}

	for _, test := range tests {
		result, err := B64Decode(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("B64Decode(%q) error = %v; want error = %v", test.input, err != nil, test.hasError)
		}
		if result != test.expected {
			t.Errorf("B64Decode(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
