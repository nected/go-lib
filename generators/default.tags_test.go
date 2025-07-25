package generators

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
		Duration time.Duration     `default:"10s"`
		SR       Timeout           `default:"{\"max\":\"10s\"}"`
		Retry    CustomRetryPolicy `default:"{\"initialInterval\":\"2s\",\"maximumInterval\":\"1s\",\"backoffCoefficient\":4,\"maximumAttempts\":6}"`
	}

	var defaultArgs args
	GenerateDefaults(&defaultArgs)
	checkArgs := args{
		Duration: 10 * time.Second,
		SR:       Timeout{Default: 30 * time.Second, Max: 10 * time.Second},
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
}
