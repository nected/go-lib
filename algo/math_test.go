package algo

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MAX returns the larger of two ints; ties return b (the second arg) because
// the `a > b` branch is false when equal.
func TestMAX(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{name: "a greater", a: 5, b: 3, want: 5},
		{name: "b greater", a: 3, b: 5, want: 5},
		{name: "equal returns b", a: 5, b: 5, want: 5},
		{name: "negatives", a: -5, b: -1, want: -1},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, MAX(test.a, test.b))
		})
	}
}

// MIN returns the smaller of two ints; ties return b (the second arg) because
// the `a < b` branch is false when equal.
func TestMIN(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{name: "a smaller", a: 3, b: 5, want: 3},
		{name: "b smaller", a: 5, b: 3, want: 3},
		{name: "equal returns b", a: 5, b: 5, want: 5},
		{name: "negatives", a: -5, b: -1, want: -5},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, MIN(test.a, test.b))
		})
	}
}

// In reports whether key is present in list. Generics are tested via separate
// subtests since a single table cannot mix type parameters.
func TestIn(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		tests := []struct {
			name string
			key  string
			list []string
			want bool
		}{
			{name: "present", key: "b", list: []string{"a", "b", "c"}, want: true},
			{name: "absent", key: "z", list: []string{"a", "b", "c"}, want: false},
			{name: "empty list", key: "a", list: []string{}, want: false},
			{name: "nil list", key: "a", list: nil, want: false},
		}
		for id, test := range tests {
			t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
				assert.Equal(t, test.want, In(test.key, test.list))
			})
		}
	})

	t.Run("ints", func(t *testing.T) {
		assert.True(t, In(2, []int{1, 2, 3}))
		assert.False(t, In(4, []int{1, 2, 3}))
	})

	t.Run("bools", func(t *testing.T) {
		assert.True(t, In(true, []bool{false, true}))
		assert.False(t, In(true, []bool{false, false}))
	})
}

// Intersect returns elements present in both lists, preserving list1's order
// (and its duplicates if they appear in list2). Empty intersection yields nil,
// not an empty slice, because `append` on a nil slice stays nil until first call.
func TestIntersect(t *testing.T) {
	tests := []struct {
		name  string
		list1 []int
		list2 []int
		want  []int
	}{
		{name: "partial overlap", list1: []int{1, 2, 3}, list2: []int{2, 3, 4}, want: []int{2, 3}},
		{name: "no overlap returns nil", list1: []int{1, 2}, list2: []int{3, 4}, want: nil},
		{name: "empty list1 returns nil", list1: []int{}, list2: []int{1, 2}, want: nil},
		{name: "duplicates preserved from list1", list1: []int{1, 1, 2}, list2: []int{1, 2}, want: []int{1, 1, 2}},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, Intersect(test.list1, test.list2))
		})
	}
}

// Union returns the distinct set across both lists. Order is non-deterministic
// (map iteration), so results are sorted before comparison.
func TestUnion(t *testing.T) {
	tests := []struct {
		name  string
		list1 []int
		list2 []int
		want  []int
	}{
		{name: "disjoint", list1: []int{1, 2, 3}, list2: []int{4, 5}, want: []int{1, 2, 3, 4, 5}},
		{name: "overlap collapsed", list1: []int{1, 2, 3}, list2: []int{3, 4, 5}, want: []int{1, 2, 3, 4, 5}},
		{name: "intra-list duplicates collapsed", list1: []int{1, 1, 2}, list2: []int{2, 2}, want: []int{1, 2}},
		// Both empty: the map has no keys, range produces nothing, finalList stays nil.
		{name: "both empty returns nil", list1: []int{}, list2: []int{}, want: nil},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			got := Union(test.list1, test.list2)
			sort.Ints(got)
			assert.Equal(t, test.want, got)
		})
	}
}

// AminusB returns elements in list1 not in list2, preserving list1 order.
// Full overlap returns nil (append-on-nil stays nil).
func TestAminusB(t *testing.T) {
	tests := []struct {
		name  string
		list1 []int
		list2 []int
		want  []int
	}{
		{name: "some removed", list1: []int{1, 2, 3}, list2: []int{2, 3, 4}, want: []int{1}},
		{name: "none removed", list1: []int{1, 2}, list2: []int{3, 4}, want: []int{1, 2}},
		{name: "all removed returns nil", list1: []int{1, 2}, list2: []int{1, 2, 3}, want: nil},
		{name: "empty list1 returns nil", list1: []int{}, list2: []int{1}, want: nil},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, AminusB(test.list1, test.list2))
		})
	}
}

// GenerateRandomString reads crypto/rand and URL-base64-encodes the bytes.
// The encoded length is the deterministic ceil(strLen/3)*4, and two successive
// calls with the same length must differ (sanity check on randomness).
func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name      string
		strLen    int
		wantLen   int
		wantEmpty bool
	}{
		{name: "len 0 returns empty", strLen: 0, wantLen: 0, wantEmpty: true},
		// 16 raw bytes -> ceil(16/3)*4 = 24 base64 chars (with padding).
		{name: "len 16", strLen: 16, wantLen: 24},
		// 32 raw bytes -> 44 base64 chars.
		{name: "len 32", strLen: 32, wantLen: 44},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			got, err := GenerateRandomString(test.strLen)
			assert.NoError(t, err)
			if test.wantEmpty {
				assert.Empty(t, got)
				return
			}
			assert.Equal(t, test.wantLen, len(got))
		})
	}

	// Two successive calls must produce different outputs — the chance of
	// collision on 16 random bytes is negligible.
	t.Run("two calls differ", func(t *testing.T) {
		a, err1 := GenerateRandomString(16)
		b, err2 := GenerateRandomString(16)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, a, b)
	})
}

// RemoveDuplicate preserves first-seen order. Unlike Intersect / AminusB it
// initializes `list` to an empty slice, so an empty input returns []int{} not nil.
func TestRemoveDuplicate(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "mixed duplicates", input: []int{1, 2, 2, 3, 1, 3}, want: []int{1, 2, 3}},
		{name: "all unique", input: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "all same", input: []int{5, 5, 5, 5}, want: []int{5}},
		{name: "empty returns empty slice", input: []int{}, want: []int{}},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, RemoveDuplicate(test.input))
		})
	}
}
