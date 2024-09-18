package rsa

import "testing"

func Test_parseData(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 string
	}{
		{
			name: "Test 1",
			args: args{
				data: "$keyName$keyVersion$encryptedData",
			},
			want:  "keyName",
			want1: "keyVersion",
			want2: "encryptedData",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := parseData(tt.args.data)
			if got != tt.want {
				t.Errorf("parseData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseData() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("parseData() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
