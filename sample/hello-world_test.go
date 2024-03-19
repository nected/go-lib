package sample

import "testing"

func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"TestHelloWorld", "\x1b[32mHello, world!\x1b[0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HelloWorld(); got != tt.want {
				t.Errorf("HelloWorld() = %v, want %v", got, tt.want)
			}
		})
	}
}
