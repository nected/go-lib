package generators

import (
	"testing"
)

type SampleSubStructA struct {
	Str string `default:"test"`
	Int int    `default:"1"`
}

type SampleStruct struct {
	DefaultTag
	SubStructA SampleSubStructA
	Str        string `default:"test"`
	Int        int    `default:"1"`
}

func TestDefaultTag_GenerateDefaults(t *testing.T) {
	tests := []struct {
		name string
		tr   *SampleStruct
	}{
		{"Test 1", &SampleStruct{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &SampleStruct{}
			tr.GenerateDefaults()
		})
	}
}
