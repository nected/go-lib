package algo

import (
	"fmt"
	"reflect"
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
		// empty key returns source as-is
		{
			source: map[string]any{"name": "ram"},
			keyStr: "",
			want:   wantT{map[string]any{"name": "ram"}, nil},
		},
		// top-level map key lookup
		{
			source: map[string]any{"name": "ram"},
			keyStr: "name",
			want:   wantT{"ram", nil},
		},
		// nested map lookup via dotted path
		{
			source: map[string]any{"data": map[string]interface{}{"name": "ram"}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		// slice value inside a map, accessed by index
		{
			source: map[string]any{"data": []string{"name", "ram"}},
			keyStr: "data[1]",
			want:   wantT{"ram", nil},
		},
		// struct field lookup (exact case match)
		{
			source: map[string]any{"data": struct{ Name string }{Name: "ram"}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		// struct field lookup is case-insensitive
		{
			source: map[string]any{"data": struct{ NaMe string }{NaMe: "ram"}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		// pointer struct field is dereferenced to its underlying value
		{
			source: map[string]any{"data": struct{ NaMe *string }{NaMe: &ram}},
			keyStr: "data.name",
			want:   wantT{"ram", nil},
		},
		// nil pointer struct field returns the typed nil
		{
			source: map[string]any{"data": struct{ NaMe *string }{NaMe: nil}},
			keyStr: "data.name",
			want:   wantT{stringNil, nil},
		},
		// nil source with empty key returns source as-is
		{
			source: nil,
			keyStr: "",
			want:   wantT{nil, nil},
		},
		// nil source with key, default (no error option) returns nil, nil
		{
			source: nil,
			keyStr: "name",
			want:   wantT{nil, nil},
		},
		// nil source with key, error option returns error
		{
			source:  nil,
			keyStr:  "name",
			options: []string{ERROR_MISSING_KEY_VALUE},
			want:    wantT{nil, fmt.Errorf(errSourceNullForKey, "name")},
		},
		// missing map key, default swallows error and returns nil
		{
			source: map[string]any{"name": "ram"},
			keyStr: "missing",
			want:   wantT{nil, nil},
		},
		// missing map key, error option returns error
		{
			source:  map[string]any{"name": "ram"},
			keyStr:  "missing",
			options: []string{ERROR_MISSING_KEY_VALUE},
			want:    wantT{nil, fmt.Errorf(errKeyNotPresent, "missing")},
		},
		// missing struct field, error option returns error
		{
			source:  struct{ Name string }{Name: "ram"},
			keyStr:  "age",
			options: []string{ERROR_MISSING_KEY_VALUE},
			want:    wantT{nil, fmt.Errorf(errKeyNotPresent, "age")},
		},
		// direct slice source with index
		{
			source: []string{"foo", "bar"},
			keyStr: "[1]",
			want:   wantT{"bar", nil},
		},
		// direct array source with index
		{
			source: [2]int{10, 20},
			keyStr: "[0]",
			want:   wantT{10, nil},
		},
		// nested indexes on 2d slice
		{
			source: map[string]any{"data": [][]int{{1, 2}, {3, 4}}},
			keyStr: "data[1][0]",
			want:   wantT{3, nil},
		},
		// direct string source with index returns single char as string
		{
			source: "hello",
			keyStr: "[1]",
			want:   wantT{"e", nil},
		},
		// out of bound slice index, default returns nil
		{
			source: map[string]any{"data": []string{"a"}},
			keyStr: "data[5]",
			want:   wantT{nil, nil},
		},
		// out of bound slice index with error option
		{
			source:  map[string]any{"data": []string{"a"}},
			keyStr:  "data[5]",
			options: []string{ERROR_MISSING_KEY_VALUE},
			want:    wantT{nil, fmt.Errorf(errOutOfIndex)},
		},
		// out of bound string index, default returns nil
		{
			source: "hi",
			keyStr: "[5]",
			want:   wantT{nil, nil},
		},
		// out of bound string index with error option
		{
			source:  "hi",
			keyStr:  "[5]",
			options: []string{ERROR_MISSING_KEY_VALUE},
			want:    wantT{nil, fmt.Errorf(errIndexOutOfBound, 5)},
		},
		// key used on list item
		{
			source: []string{"foo"},
			keyStr: "name[0]",
			want:   wantT{nil, fmt.Errorf(errKeyUsedOnListItem, "name[0]")},
		},
		// non-map/list/string source with key
		{
			source: 42,
			keyStr: "name",
			want:   wantT{nil, fmt.Errorf(errInvalidKeyUsage, "name")},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			res, err := GetValFromSource(test.source, test.keyStr, test.options...)
			assert.Equal(t, test.want, wantT{res, err})
		})
	}
}

func TestGetMapOrStructKeyValue(t *testing.T) {
	var ram = "ram"
	tests := []struct {
		source any
		key    string
		want   wantT
	}{
		// nil source returns (nil, nil)
		{
			source: nil,
			key:    "name",
			want:   wantT{nil, nil},
		},
		// valid map key
		{
			source: map[string]any{"name": "ram"},
			key:    "name",
			want:   wantT{"ram", nil},
		},
		// missing map key returns error
		{
			source: map[string]any{"name": "ram"},
			key:    "missing",
			want:   wantT{nil, fmt.Errorf(errInvalidValue, reflect.Value{})},
		},
		// struct field lookup (case-insensitive)
		{
			source: struct{ Name string }{Name: "ram"},
			key:    "name",
			want:   wantT{"ram", nil},
		},
		// pointer struct field dereferenced
		{
			source: struct{ Name *string }{Name: &ram},
			key:    "name",
			want:   wantT{"ram", nil},
		},
		// missing struct field returns error
		{
			source: struct{ Name string }{Name: "ram"},
			key:    "age",
			want:   wantT{nil, fmt.Errorf(errInvalidValueForKey, "age")},
		},
		// non-map/struct source hits default case
		{
			source: 42,
			key:    "name",
			want:   wantT{nil, fmt.Errorf(errNotAMap, reflect.Int)},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			res, err := getMapOrStructKeyValue(test.source, test.key)
			assert.Equal(t, test.want, wantT{res, err})
		})
	}
}

func TestGetArrayIndexValue(t *testing.T) {
	tests := []struct {
		arr             any
		idx             int
		missingKeyError bool
		want            wantT
	}{
		// valid slice index
		{
			arr:  []int{10, 20, 30},
			idx:  1,
			want: wantT{20, nil},
		},
		// valid array index
		{
			arr:  [3]string{"a", "b", "c"},
			idx:  2,
			want: wantT{"c", nil},
		},
		// out of bound, default swallows to nil
		{
			arr:  []int{1},
			idx:  5,
			want: wantT{nil, nil},
		},
		// out of bound with missingKeyError returns error
		{
			arr:             []int{1},
			idx:             5,
			missingKeyError: true,
			want:            wantT{nil, fmt.Errorf(errOutOfIndex)},
		},
		// non-slice/array input returns error
		{
			arr:  "not-a-slice",
			idx:  0,
			want: wantT{nil, fmt.Errorf(errNotAnArray, reflect.String)},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			res, err := getArrayIndexValue(test.arr, test.idx, test.missingKeyError)
			assert.Equal(t, test.want, wantT{res, err})
		})
	}
}

func TestExtractKeyIndex(t *testing.T) {
	tests := []struct {
		key         string
		wantKey     string
		wantIndexes []int
		wantErr     error
	}{
		// plain key with no index
		{
			key:         "name",
			wantKey:     "name",
			wantIndexes: []int{},
		},
		// key with a single index
		{
			key:         "data[0]",
			wantKey:     "data",
			wantIndexes: []int{0},
		},
		// key with nested indexes
		{
			key:         "data[1][2]",
			wantKey:     "data",
			wantIndexes: []int{1, 2},
		},
		// index-only (no key prefix)
		{
			key:         "[3]",
			wantKey:     "",
			wantIndexes: []int{3},
		},
		// spaces are stripped before regex match
		{
			key:         "data [0]",
			wantKey:     "data",
			wantIndexes: []int{0},
		},
		// hyphenated key still matches
		{
			key:         "my-key[0]",
			wantKey:     "my-key",
			wantIndexes: []int{0},
		},
		// malformed brackets don't match the regex — treated as plain key
		{
			key:         "data[abc]",
			wantKey:     "data[abc]",
			wantIndexes: []int{},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			key, indexes, err := extractKeyIndex(test.key)
			assert.Equal(t, test.wantKey, key)
			assert.Equal(t, test.wantIndexes, indexes)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func TestGetListIndexValue(t *testing.T) {
	tests := []struct {
		source          any
		indexes         []int
		missingKeyError bool
		want            wantT
	}{
		// empty index list returns source unchanged
		{
			source: []int{1, 2, 3},
			want:   wantT{[]int{1, 2, 3}, nil},
		},
		// nil source with indexes returns error
		{
			source:  nil,
			indexes: []int{0},
			want:    wantT{nil, fmt.Errorf(errSourceMustNotBeNil)},
		},
		// single index on slice
		{
			source:  []string{"a", "b"},
			indexes: []int{1},
			want:    wantT{"b", nil},
		},
		// nested indexes on 2d slice
		{
			source:  [][]int{{1, 2}, {3, 4}},
			indexes: []int{1, 0},
			want:    wantT{3, nil},
		},
		// index into a string extracts single char
		{
			source:  "hello",
			indexes: []int{1},
			want:    wantT{"e", nil},
		},
		// out of bound string index, default returns nil
		{
			source:  "hi",
			indexes: []int{5},
			want:    wantT{nil, nil},
		},
		// out of bound string index with missingKeyError
		{
			source:          "hi",
			indexes:         []int{5},
			missingKeyError: true,
			want:            wantT{nil, fmt.Errorf(errIndexOutOfBound, 5)},
		},
		// non-slice/array/string source hits default case
		{
			source:  42,
			indexes: []int{0},
			want:    wantT{nil, fmt.Errorf(errInvalidIndexUsage)},
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v", id), func(t *testing.T) {
			res, err := getListIndexValue(test.source, test.indexes, test.missingKeyError)
			assert.Equal(t, test.want, wantT{res, err})
		})
	}
}

// UpdateValueToSource mutates a map[string]any by dotted path. Each case
// declares the starting source, the path+value, the expected error (nil if
// none), and the expected source state after the call (nil = skip check,
// e.g. when source itself is nil or not a map).
func TestUpdateValueToSource(t *testing.T) {
	tests := []struct {
		name       string
		source     any
		path       string
		value      any
		wantErr    error
		wantSource any
	}{
		{
			// empty path leaves source untouched and returns no error
			name:       "empty path is no-op",
			source:     map[string]interface{}{"a": 1},
			path:       "",
			value:      "new",
			wantSource: map[string]interface{}{"a": 1},
		},
		{
			// nil source is rejected before any traversal
			name:    "nil source returns error",
			source:  nil,
			path:    "a",
			value:   "new",
			wantErr: fmt.Errorf("source must not be null"),
		},
		{
			// only map[string]interface{} is supported as a source
			name:    "non-map source returns error",
			source:  []int{1, 2},
			path:    "a",
			value:   "new",
			wantErr: fmt.Errorf("source should be map"),
		},
		{
			// single-segment path updates the matching top-level key
			name:       "top-level key update",
			source:     map[string]interface{}{"name": "old", "age": 10},
			path:       "name",
			value:      "new",
			wantSource: map[string]interface{}{"name": "new", "age": 10},
		},
		{
			// multi-segment path recurses into a nested map and updates the leaf
			name: "nested key update",
			source: map[string]interface{}{
				"user": map[string]interface{}{"name": "old", "age": 10},
			},
			path:  "user.name",
			value: "new",
			wantSource: map[string]interface{}{
				"user": map[string]interface{}{"name": "new", "age": 10},
			},
		},
		{
			// missing top-level key is a silent no-op (no error, no insert)
			name:       "unknown top-level key is silent no-op",
			source:     map[string]interface{}{"a": 1},
			path:       "missing",
			value:      "new",
			wantSource: map[string]interface{}{"a": 1},
		},
		{
			// intermediate path segment must itself be a map; other types error
			name:    "nested intermediate is not a map returns error",
			source:  map[string]interface{}{"user": "plain-string"},
			path:    "user.name",
			value:   "new",
			wantErr: fmt.Errorf("unsupported value update user.name"),
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			err := UpdateValueToSource(test.source, test.path, test.value)
			assert.Equal(t, test.wantErr, err)
			if test.wantSource != nil {
				assert.Equal(t, test.wantSource, test.source)
			}
		})
	}
}
