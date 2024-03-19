package generators

import "testing"

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
