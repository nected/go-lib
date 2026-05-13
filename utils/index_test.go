package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleTarget struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// JsonToStruct marshals input to JSON, then unmarshals into output. It surfaces
// errors from either step. The interface{} input means non-marshalable values
// (chan/func) fail at Marshal; type-mismatched output fails at Unmarshal.
func TestJsonToStruct(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		output  any
		wantErr bool
		check   func(t *testing.T, out any)
	}{
		{
			name:   "map to struct",
			input:  map[string]any{"name": "alice", "count": 3},
			output: &sampleTarget{},
			check: func(t *testing.T, out any) {
				got := out.(*sampleTarget)
				assert.Equal(t, "alice", got.Name)
				assert.Equal(t, 3, got.Count)
			},
		},
		{
			name:   "struct to map",
			input:  sampleTarget{Name: "bob", Count: 7},
			output: &map[string]any{},
			check: func(t *testing.T, out any) {
				got := *out.(*map[string]any)
				assert.Equal(t, "bob", got["name"])
				// JSON numbers decode into float64 when target is map[string]any.
				assert.Equal(t, float64(7), got["count"])
			},
		},
		{
			// chan cannot be JSON-marshaled — fails at the Marshal step.
			name:    "unmarshalable input returns error",
			input:   make(chan int),
			output:  &sampleTarget{},
			wantErr: true,
		},
		{
			// Target field types don't match — fails at the Unmarshal step.
			name:    "type mismatch returns error",
			input:   map[string]any{"name": "x", "count": "not-an-int"},
			output:  &sampleTarget{},
			wantErr: true,
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			err := JsonToStruct(test.input, test.output)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if test.check != nil {
				test.check(t, test.output)
			}
		})
	}
}
