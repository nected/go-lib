package aes

import (
	"reflect"
	"testing"

	"github.com/nected/go-lib/crypto/models"
)

func TestEncryptAES(t *testing.T) {
	type args struct {
		secret        string
		data          []byte
		encryptedData string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Payload
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				secret: "someRandomSecret",
				data:   []byte("data"),
			},
			want: &models.Payload{
				KeyType:       models.KeyTypeAES,
				Data:          "data",
				EncryptedData: "9zE+AFpfb3PhIfdaOlPxXZAVHb3oEiTxMYcIoDuaYVs=",
			},
			wantErr: false,
		},
		{
			name: "Test 2 : Invalid secret",
			args: args{
				secret: "1",
				data:   []byte("data"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.secret, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("EncryptAES() = %v, want data %v", got.Data, tt.want.Data)
			}

			if !reflect.DeepEqual(got.EncryptedData, tt.want.EncryptedData) {
				payload, err := Decrypt(tt.args.secret, got.EncryptedData)
				if err != nil {
					t.Errorf("DecryptAES() error = %v", err)
				}
				if !reflect.DeepEqual(payload.Data, tt.want.Data) {
					t.Errorf("DecryptAES() = %v, want data %v", payload.Data, tt.want.Data)
				}

			}
		})
	}
}

func TestDecryptAES(t *testing.T) {
	type args struct {
		secret string
		data   string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Payload
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				secret: "someRandomSecret",
				data:   "9zE+AFpfb3PhIfdaOlPxXZAVHb3oEiTxMYcIoDuaYVs=",
			},
			want: &models.Payload{
				KeyType: models.KeyTypeAES,
				Data:    "data",
			},
			wantErr: false,
		},
		{
			name: "Test 2 : Empty Data",
			args: args{
				secret: "someRandomSecret",
				data:   "",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test 3 : Invalid encoded data",
			args: args{
				secret: "someRandomSecret",
				data:   "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 4 : Invalid secret",
			args: args{
				secret: "1",
				data:   "9zE+AFpfb3PhIfdaOlPxXZAVHb3oEiTxMYcIoDuaYVs=",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.secret, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("DecryptAES() = %v, want %v", got, tt.want)
			}
		})
	}
}
