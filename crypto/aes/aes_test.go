package aes

import (
	"reflect"
	"testing"

	"github.com/nected/go-lib/crypto/models"
)

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
			name: "Valid decryption",
			args: args{
				secret: "someRandomSecret",
				data:   "nVsk7LROacUqg5p70an8BBqBXwDYDGmhNz1CJiyfLak=",
			},
			want: &models.Payload{
				Data: "data",
			},
			wantErr: false,
		},
		{
			name: "Empty data",
			args: args{
				secret: "someRandomSecret",
				data:   "",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Invalid base64 data",
			args: args{
				secret: "someRandomSecret",
				data:   "invalid_base64",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid secret",
			args: args{
				secret: "1",
				data:   "9zE+AFpfb3PhIfdaOlPxXZAVHb3oEiTxMYcIoDuaYVs=",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Short decoded data",
			args: args{
				secret: "someRandomSecret",
				data:   "short",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.secret, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("Decrypt() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}
func TestEncryptAES(t *testing.T) {
	type args struct {
		secret string
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Payload
		wantErr bool
	}{
		{
			name: "Valid encryption",
			args: args{
				secret: "someRandomSecret",
				data:   []byte("data"),
			},
			want: &models.Payload{
				KeyType: models.KeyTypeAES,
				Data:    "data",
			},
			wantErr: false,
		},
		{
			name: "Empty data",
			args: args{
				secret: "someRandomSecret",
				data:   []byte(""),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Invalid secret",
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
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}
			if got.Data != tt.want.Data {
				t.Errorf("Encrypt() = %v, want %v", got.Data, tt.want.Data)
			}
			if got.KeyType != tt.want.KeyType {
				t.Errorf("Encrypt() KeyType = %v, want %v", got.KeyType, tt.want.KeyType)
			}
			if !tt.wantErr {
				// Decrypt to verify the encrypted data
				decryptedPayload, err := Decrypt(tt.args.secret, got.EncryptedData)
				if err != nil {
					t.Errorf("Decrypt() error = %v", err)
				}
				if decryptedPayload.Data != tt.want.Data {
					t.Errorf("Decrypt() = %v, want %v", decryptedPayload.Data, tt.want.Data)
				}
			}
		})
	}
}
