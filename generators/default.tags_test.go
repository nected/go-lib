package generators

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Verifies that `default` tag values override already-populated fields
// across a nested struct (string + int), including a grandchild struct
// reached through zero-JSON intermediate tags.
func TestGenerateDefaults(t *testing.T) {
	type SubAddress struct {
		Street string `default:"sub street"`
	}
	type Address struct {
		Street string `default:"street"`
		Number int    `default:"100"`
		Sub    SubAddress
	}
	type args struct {
		Name    string `default:"name"`
		Age     int    `default:"10"`
		Address Address
	}
	tests := []struct {
		name string
		args args
		want args
	}{
		{
			"Test case 1",
			args{Name: "hidden name", Age: 11, Address: Address{Street: "hidden street", Number: 101, Sub: SubAddress{Street: "hidden sub street"}}},
			args{Name: "name", Age: 10, Address: Address{Street: "street", Number: 100, Sub: SubAddress{Street: "sub street"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenerateDefaults(&tt.args)
			if tt.args != tt.want {
				t.Errorf("GenerateDefaults() = %v, want %v", tt.args, tt.want)
			}
		})
	}
}

// Verifies time.Duration parsing for reflect.Int64 fields and JSON-tag
// overrides on nested structs: the parent's JSON `default` tag supplies
// per-field overrides that win over the child's own scalar tags.
func TestGenerateDefaults2(t *testing.T) {
	type Timeout struct {
		Default time.Duration `default:"30s"`
		Max     time.Duration `default:"60s"`
	}

	type CustomRetryPolicy struct {
		InitialInterval    time.Duration `default:"1s"`
		MaximumInterval    time.Duration `default:"30s"`
		BackoffCoefficient float64       `default:"10"`
		MaximumAttempts    int           `default:"2"`
		Enabled            bool          `default:"true"`
	}

	type args struct {
		Duration    time.Duration     `default:"10s"`
		SR          Timeout           `default:"{\"max\":\"10s\"}"`
		ArrayInt    []int             `default:"[1,2,3]"`
		ArrayStruct []Timeout         `default:"[{\"max\":\"10s\"},{\"default\":\"10s\"}]"`
		Retry       CustomRetryPolicy `default:"{\"initialInterval\":\"2s\",\"maximumInterval\":\"1s\",\"backoffCoefficient\":4,\"maximumAttempts\":6}"`
	}

	var defaultArgs args
	GenerateDefaults(&defaultArgs)
	checkArgs := args{
		Duration:    10 * time.Second,
		SR:          Timeout{Default: 30 * time.Second, Max: 10 * time.Second},
		ArrayInt:    []int{1, 2, 3},
		ArrayStruct: []Timeout{{Default: 30 * time.Second, Max: 10 * time.Second}, {Default: 10 * time.Second, Max: 60 * time.Second}},
		Retry: CustomRetryPolicy{
			InitialInterval:    2 * time.Second,
			MaximumInterval:    time.Second,
			BackoffCoefficient: 4,
			MaximumAttempts:    6,
			Enabled:            true,
		},
	}

	assert.Equal(t, defaultArgs.Duration, checkArgs.Duration)
	assert.Equal(t, defaultArgs.SR.Max, checkArgs.SR.Max)
	assert.Equal(t, defaultArgs.SR.Default, checkArgs.SR.Default)
	assert.Equal(t, defaultArgs.Retry.InitialInterval, checkArgs.Retry.InitialInterval)
	assert.Equal(t, defaultArgs.Retry.MaximumInterval, checkArgs.Retry.MaximumInterval)
	assert.Equal(t, defaultArgs.Retry.BackoffCoefficient, checkArgs.Retry.BackoffCoefficient)
	assert.Equal(t, defaultArgs.Retry.MaximumAttempts, checkArgs.Retry.MaximumAttempts)
	assert.Equal(t, defaultArgs.Retry.Enabled, checkArgs.Retry.Enabled)
	assert.EqualValues(t, defaultArgs.ArrayInt, checkArgs.ArrayInt)
	assert.EqualValues(t, defaultArgs.ArrayStruct, checkArgs.ArrayStruct)
}

// Covers the dVal branches that are reachable via JSON unmarshalling
// of the parent struct's default tag: string, bool, float64.
func TestGenerateDefaults_NestedJSONOverrides(t *testing.T) {
	type Inner struct {
		Name     string  `default:"tag name"`
		Active   bool    `default:"false"`
		Count    int     `default:"0"`
		Distance float64 `default:"0"`
	}
	type Outer struct {
		// JSON keys correspond to lowerCamelCase field names.
		In Inner `default:"{\"name\":\"json name\",\"active\":true,\"count\":42,\"distance\":3.14}"`
	}

	var o Outer
	GenerateDefaults(&o)

	assert.Equal(t, "json name", o.In.Name)
	assert.Equal(t, true, o.In.Active)
	assert.Equal(t, 42, o.In.Count)
	assert.Equal(t, 3.14, o.In.Distance)
}

// Covers Int8/Int16/Int32 tag parsing (numeric + duration strings).
// NOTE: reflect.Uint is also routed through this case in the source but
// field.SetInt panics on uint kinds — see TestGenerateDefaults_UintPanicsBug.
func TestGenerateDefaults_SmallIntKindsFromTag(t *testing.T) {
	type args struct {
		I8  int8  `default:"12"`
		I16 int16 `default:"1200"`
		I32 int32 `default:"120000"`
		// Duration-style string should go through ParseDuration in proccessInt64.
		I32Dur int32 `default:"2ns"`
	}

	var a args
	GenerateDefaults(&a)

	assert.Equal(t, int8(12), a.I8)
	assert.Equal(t, int16(1200), a.I16)
	assert.Equal(t, int32(120000), a.I32)
	assert.Equal(t, int32(2), a.I32Dur)
}

// Documents a latent source-code bug: reflect.Uint is grouped with signed
// int kinds at default.tags.go:51, but field.SetInt panics on uint fields.
// Remove this test and update the test above once the source is fixed to
// route uint* kinds through field.SetUint.
func TestGenerateDefaults_UintPanicsBug(t *testing.T) {
	type args struct {
		U uint `default:"7"`
	}
	var a args
	assert.Panics(t, func() { GenerateDefaults(&a) })
}

// Float32 tag-path coverage.
func TestGenerateDefaults_Float32FromTag(t *testing.T) {
	type args struct {
		F32 float32 `default:"1.5"`
	}
	var a args
	GenerateDefaults(&a)
	assert.Equal(t, float32(1.5), a.F32)
}

// Invalid numeric / float tags should silently default to zero (errors swallowed).
func TestGenerateDefaults_InvalidTags(t *testing.T) {
	type args struct {
		N int     `default:"not-a-number"`
		F float64 `default:"not-a-float"`
	}
	var a args
	GenerateDefaults(&a)
	assert.Equal(t, 0, a.N)
	assert.Equal(t, float64(0), a.F)
}

// Invalid JSON in a struct-level tag should be swallowed; inner defaults still apply.
func TestGenerateDefaults_InvalidJSONStructTag(t *testing.T) {
	type Inner struct {
		Name string `default:"fallback"`
	}
	type Outer struct {
		In Inner `default:"{not-json"`
	}
	var o Outer
	GenerateDefaults(&o)
	assert.Equal(t, "fallback", o.In.Name)
}

// GenerateDefaults on a non-struct pointer should be a no-op.
func TestGenerateDefaults_NonStruct(t *testing.T) {
	n := 5
	GenerateDefaults(&n)
	assert.Equal(t, 5, n)
}

// Exercises dVal branches that JSON unmarshalling cannot produce (int8/16/32,
// uint, int, int64, float32). processStruct is package-private but reachable
// from same-package tests.
func TestProcessStruct_DirectDValTypes(t *testing.T) {
	// Uint intentionally omitted — the source bug at default.tags.go:51 routes
	// uint through SetInt which panics. See TestGenerateDefaults_UintPanicsBug.
	type S struct {
		I8  int8
		I16 int16
		I32 int32
		I   int
		I64 int64
		F32 float32
	}
	s := S{}
	val := reflect.ValueOf(&s).Elem()
	defaults := map[string]any{
		"i8":  int8(1),
		"i16": int16(2),
		"i32": int32(3),
		"i":   int(5),
		"i64": int64(6),
		"f32": float32(7.5),
	}
	processStruct(val, defaults)

	assert.Equal(t, int8(1), s.I8)
	assert.Equal(t, int16(2), s.I16)
	assert.Equal(t, int32(3), s.I32)
	assert.Equal(t, 5, s.I)
	assert.Equal(t, int64(6), s.I64)
	assert.Equal(t, float32(7.5), s.F32)
}

// reflect.Int64 + string dVal should be parsed as a time.Duration.
func TestProcessStruct_Int64StringDurationDVal(t *testing.T) {
	type S struct {
		D time.Duration
	}
	s := S{}
	val := reflect.ValueOf(&s).Elem()
	processStruct(val, map[string]any{"d": "5s"})
	assert.Equal(t, 5*time.Second, s.D)
}

// Covers the int and float64 dVal branches for reflect.Int64 — not reachable
// via the JSON-tag code path (JSON produces float64 only for numbers, and the
// float64-into-int64 branch shadows nothing when values are whole).
func TestProcessStruct_Int64IntAndFloat64DVal(t *testing.T) {
	type S struct {
		A int64
		B int64
	}
	s := S{}
	val := reflect.ValueOf(&s).Elem()
	processStruct(val, map[string]any{"a": int(9), "b": float64(10.9)})
	assert.Equal(t, int64(9), s.A)
	assert.Equal(t, int64(10), s.B)
}

// Unit-tests proccessInt64: duration strings return nanoseconds, plain
// integers parse as int64, and unparseable input falls back to zero
// (ParseDuration error → strconv.ParseInt error → 0).
func TestProcessInt64(t *testing.T) {
	assert.Equal(t, int64(time.Second), proccessInt64("1s"))
	assert.Equal(t, int64(42), proccessInt64("42"))
	assert.Equal(t, int64(0), proccessInt64("garbage"))
}
