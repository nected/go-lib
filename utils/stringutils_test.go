package utils

import "testing"

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		name  string
		input string
		list  []string
		want  bool
	}{
		{"TestStringInSlice - String in slice", "apple", []string{"apple", "banana", "orange"}, true},
		{"TestStringInSlice - String not in slice", "grape", []string{"apple", "banana", "orange"}, false},
		{"TestStringInSlice - Empty slice", "apple", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.input, tt.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
