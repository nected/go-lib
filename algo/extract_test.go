package algo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type wantT struct {
	res any
	err error
}

func TestGetValFromSource(t *testing.T) {
	var ram = "ram"
	var stringNil *string
	tests := []struct {
		source  any
		keyStr  string
		options []string
		want    wantT
	}{
		{
			source: map[string]any{"name": "ram"},
			keyStr: "",
			want:   wantT{map[string]any{"name": "ram"}, nil},
		},
		{
			source: map[string]any{"name": "ram"},
			keyStr: "name",
			want:   wantT{"ram", nil},
		},
		{
			source: map[string]any{"data": map[string]interface{}{"name": "ram"}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		{
			source: map[string]any{"data": []string{"name", "ram"}},
			keyStr: "data[1]",
			want:   wantT{"ram", nil},
		},
		{
			source: map[string]any{"data": struct{ Name string }{Name: "ram"}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		{
			source: map[string]any{"data": struct{ NaMe string }{NaMe: "ram"}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		{
			source: map[string]any{"data": struct{ NaMe *string }{NaMe: &ram}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		{
			source: map[string]any{"data": struct{ NaMe *string }{NaMe: nil}},
			keyStr: "data.name",
			want:   wantT{stringNil, nil},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			res, err := GetValFromSource(test.source, test.keyStr, test.options...)
			assert.Equal(t, test.want, wantT{res, err})
		})
	}
}
