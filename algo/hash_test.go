package algo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetHash produces a non-empty sqids-encoded string that honors the
// MinLength=5 option. Different inputs must produce different hashes.
func TestGetHash(t *testing.T) {
	tests := []struct {
		name      string
		input     uint
		minLength int
	}{
		{name: "small number", input: 1, minLength: 5},
		{name: "zero", input: 0, minLength: 5},
		{name: "larger number", input: 12345, minLength: 5},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			h := GetHash(test.input)
			assert.NotEmpty(t, h)
			assert.GreaterOrEqual(t, len(h), test.minLength)
		})
	}

	// Uniqueness check lives outside the table — it's a relationship between
	// two GetHash calls, not a single-input assertion.
	t.Run("different inputs produce different hashes", func(t *testing.T) {
		assert.NotEqual(t, GetHash(1), GetHash(2))
	})
}

// DecodeHash reverses GetHash round-trips and returns nil on invalid input.
// Note: encode uses MinLength; decode uses default options. Round-trip still
// works because MinLength only pads output, it doesn't alter the decoded value.
func TestDecodeHash(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		want    *uint64
		wantNil bool
	}{
		// Empty string is the only input that returns nil without touching sqids.
		{name: "empty returns nil", hash: "", wantNil: true},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			got := DecodeHash(test.hash)
			if test.wantNil {
				assert.Nil(t, got)
				return
			}
			assert.Equal(t, test.want, got)
		})
	}

	// Round-trip verifies GetHash → DecodeHash recovers the original value.
	// Kept as a separate subtest so the hash isn't hard-coded (sqids output
	// format may vary by library version).
	t.Run("round trip", func(t *testing.T) {
		roundTrips := []uint{1, 42, 12345, 999999}
		for _, v := range roundTrips {
			got := DecodeHash(GetHash(v))
			if assert.NotNil(t, got) {
				assert.Equal(t, uint64(v), *got)
			}
		}
	})
}

// Encode wraps stdlib base64 encoding. Round-trip with Decode is the most
// useful check; an explicit fixture pins the encoding scheme (StdEncoding,
// not URLEncoding) so a swap would fail the test.
func TestEncode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty string", input: "", want: ""},
		{name: "ascii", input: "hello", want: "aGVsbG8="},
		// Bytes that differ between Std and URL encoding (`+/` vs `-_`) — pins
		// the StdEncoding choice in the source.
		{name: "std-vs-url discriminator bytes", input: "\xfb\xff", want: "+/8="},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, Encode(test.input))
		})
	}
}

// Decode inverts Encode and returns base64's error verbatim on malformed input.
func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty string", input: "", want: ""},
		{name: "valid ascii", input: "aGVsbG8=", want: "hello"},
		// '!' is not a valid base64 character — surfaces a CorruptInputError.
		{name: "invalid characters", input: "not_base64!", wantErr: true},
		// Length not a multiple of 4 and missing padding — also a CorruptInputError.
		{name: "bad padding", input: "abc", wantErr: true},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			got, err := Decode(test.input)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}

	// Round-trip across a few payloads, including non-ASCII bytes.
	t.Run("round trip", func(t *testing.T) {
		for _, s := range []string{"", "hello", "with spaces & symbols!", "日本語"} {
			got, err := Decode(Encode(s))
			assert.NoError(t, err)
			assert.Equal(t, s, got)
		}
	})
}
