package logger

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetZapFields(t *testing.T) {
	type abcd struct {
		a string
	}

	a := abcd{a: "a"}
	tests := []struct {
		name   string
		args   []interface{}
		fields []zapcore.Field
	}{
		{
			name: "TestGetZapFields - No errors",
			args: []interface{}{
				a, "value1",
				"key2", "value2",
				"key3", "value3",
			},
			fields: []zapcore.Field{
				zap.Any("{a}", "value1"),
				zap.Any("key2", "value2"),
				zap.Any("key3", "value3"),
			},
		},
		{
			name: "TestGetZapFields - With errors",
			args: []interface{}{
				"key1", "value1",
				"key2", "value2",
				"key3", "value3",
				errors.New("error1"),
				errors.New("error2"),
			},
			fields: []zapcore.Field{
				zap.Any("key1", "value1"),
				zap.Any("key2", "value2"),
				zap.Any("key3", "value3"),
				zap.Errors("errors", []error{errors.New("error1")}),
				zap.Error(errors.New("error2")),
			},
		},
		{
			name: "TestGetZapFields - Last argument is an error",
			args: []interface{}{
				"key1", "value1",
				"key2", "value2",
				"key3", "value3",
				errors.New("error1"),
			},
			fields: []zapcore.Field{
				zap.Any("key1", "value1"),
				zap.Any("key2", "value2"),
				zap.Any("key3", "value3"),
				zap.Errors("errors", []error{errors.New("error1")}),
				zap.Error(errors.New("error1")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getZapFields(tt.args...)
			if len(got) != len(tt.fields) {
				t.Errorf("getZapFields() returned %d fields, want %d fields", len(got), len(tt.fields))
			}
			for i := 0; i < len(got); i++ {
				if got[i].Key != tt.fields[i].Key {
					t.Errorf("getZapFields() returned field with key %s, want %s", got[i].Key, tt.fields[i].Key)
				}
				if got[i].Type != tt.fields[i].Type {
					t.Errorf("getZapFields() returned field with type %v, want %v", got[i].Type, tt.fields[i].Type)
				}
				// Compare field values based on their types
				switch got[i].Type {
				case zapcore.StringType:
					if got[i].String != tt.fields[i].String {
						t.Errorf("getZapFields() returned field with value %s, want %s", got[i].String, tt.fields[i].String)
					}
				case zapcore.ErrorType:
					assert.Equal(t, got[i].Interface, tt.fields[i].Interface)
					// if got[i].Interface != tt.fields[i].Interface {
					// 	t.Errorf("getZapFields() returned field with value %v, want %v", got[i].Interface, tt.fields[i].Interface)
					// }
				case zapcore.ArrayMarshalerType:
					g := fmt.Sprintf("%v", got[i].Interface)
					in := fmt.Sprintf("%v", tt.fields[i].Interface)
					if g != in {
						t.Errorf("getZapFields() returned field with value %v, want %v", got[i].Interface, tt.fields[i].Interface)
					}
				default:
					t.Errorf("getZapFields() returned field with unsupported type %v", got[i].Type)
				}
			}
		})
	}
}
