package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"

	"github.com/nected/go-lib/crypto/config"
	"github.com/nected/go-lib/crypto/errors"
	"github.com/nected/go-lib/crypto/models"
)

var validEncryptedData, validEncryptedData2 string

func generatePrivateKey() string {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Cannot generate RSA key\n")
		os.Exit(1)
	}

	// encode private key to PEM format
	privatekeyBytes, err := x509.MarshalPKCS8PrivateKey(privatekey)
	if err != nil {
		fmt.Printf("Cannot marshal private key to bytes\n")
		return ""
	}

	privateKeyBlock := &pem.Block{
		Type:  config.PRIV_KEY_TYPE,
		Bytes: privatekeyBytes,
	}

	privKey := pem.EncodeToMemory(privateKeyBlock)

	return string(privKey)
}

func populateKeys() {
	payload, _ := Encrypt("TESTKEY", []byte("test data"))
	validEncryptedData = payload.String()
	payload, _ = Encrypt("TESTKEY1", []byte("test data"))
	validEncryptedData2 = payload.String()
}

func setupSuite(t *testing.T) func(t *testing.T) {
	var privateKey = generatePrivateKey()
	os.Setenv("KEY_TESTKEY_1_0", privateKey)
	os.Setenv("KEY_TESTKEY1_2_0", privateKey)
	os.Setenv("KEY_TESTKEYR_1_1726147578", privateKey)
	os.Setenv("KEY_TESTKEYINVALID_1_0", "lkajds")

	t.Log("setup suite")
	return func(t *testing.T) {
		defer os.Unsetenv("KEY_TESTKEY_1_0")
		defer os.Unsetenv("KEY_TESTKEY1_2_0")
		defer os.Unsetenv("KEY_TESTKEYR_1_1726147578")
		defer os.Unsetenv("KEY_TESTKEYINVALID_1_0")
		t.Log("teardown suite")
	}

}

func Test_parseData(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
		want2 string
	}{
		{
			name: "Test 1",
			args: args{
				data: "$keyName$1$encryptedData",
			},
			want:  "keyName",
			want1: 1,
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
func TestEncrypt(t *testing.T) {
	teardownSuite := setupSuite(t)
	config.LoadKeysFromEnv()
	defer teardownSuite(t)

	type args struct {
		keyName string
		data    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Payload
		wantErr bool
	}{
		{
			name: "Key not found",
			args: args{
				keyName: "nonexistentKey",
				data:    []byte("test data"),
			},
			want: &models.Payload{
				Data:          "test data",
				EncryptedData: "test data",
			},
			wantErr: false,
		},
		{
			name: "Valid key",
			args: args{
				keyName: "TESTKEY",
				data:    []byte("test data"),
			},
			want: &models.Payload{
				KeyName:    "TESTKEY",
				KeyVersion: 1,
				KeyType:    models.KeyTypeRSA,
				Data:       "test data",
				// EncryptedData will be checked separately
			},
			wantErr: false,
		},
		{
			name: "Invalid key",
			args: args{
				keyName: "TESTKEYINVALID",
				data:    []byte("test data"),
			},
			want: &models.Payload{
				Data:          "test data",
				EncryptedData: "test data",
			},
			wantErr: false,
		},
		{
			name: "Invalid Data",
			args: args{
				keyName: "TESTKEY",
				data:    []byte(""),
			},
			want: &models.Payload{
				KeyName:       "TESTKEY",
				KeyVersion:    1,
				KeyType:       models.KeyTypeRSA,
				Data:          "",
				EncryptedData: "",
			},
			wantErr: false,
		},
		{
			name: "Already Encrypted",
			args: args{
				keyName: "TESTKEY",
				data:    []byte("JFRFU1RLRVkkMSRhbHJlYWR5RW5jcnlwdGVk"),
			},
			want: &models.Payload{
				KeyName:       "",
				KeyVersion:    0,
				KeyType:       "",
				Data:          "JFRFU1RLRVkkMSRhbHJlYWR5RW5jcnlwdGVk",
				EncryptedData: "JFRFU1RLRVkkMSRhbHJlYWR5RW5jcnlwdGVk",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.keyName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.KeyName != tt.want.KeyName || got.KeyVersion != tt.want.KeyVersion || got.KeyType != tt.want.KeyType || got.Data != tt.want.Data {
					t.Errorf("Encrypt() got = %v, want %v", got, tt.want)
				}
				if tt.args.keyName == "TESTKEY" && tt.want.Data != "" && !got.AlreadyEncrypted && got.EncryptedData == tt.want.Data {
					t.Errorf("Encrypt() EncryptedData should not be equal to Data for valid key")
				}
			}

			if tt.args.keyName == "TESTKEYINVALID" {
				if got.EncryptedData != tt.want.Data {
					t.Errorf("Encrypt() EncryptedData should be equal to Data for invalid key")
				}
			}
		})
	}
}
func TestDecrypt(t *testing.T) {
	teardownSuite := setupSuite(t)
	config.LoadKeysFromEnv()
	populateKeys()
	defer teardownSuite(t)

	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Payload
		wantErr bool
		err     error
	}{
		{
			name: "Empty data",
			args: args{
				data: "",
			},
			want:    nil,
			wantErr: true,
			err:     errors.ErrEmptyData,
		},
		{
			name: "Not encrypted data",
			args: args{
				data: "dGVzdCBkYXRh", // base64 for "test data"
			},
			want: &models.Payload{
				Data: "dGVzdCBkYXRh",
			},
			wantErr: false,
		},
		{
			name: "Key not found: Invalid encoded data",
			args: args{
				data: "JFRFU1RLRVkkMSRlbmNyeXB0ZWRkYXRh", // base64 for "$TESTKEY$1$encrypteddata"
			},
			want:    nil,
			wantErr: true,
			err:     errors.ErrInvalidData,
		},
		{
			name: "Valid key version 1",
			args: args{
				data: validEncryptedData,
			},
			want: &models.Payload{
				KeyName:    "TESTKEY",
				KeyVersion: 1,
				Data:       "test data",
			},
			wantErr: false,
		},
		{
			name: "Valid key version 2",
			args: args{
				data: validEncryptedData2,
			},
			want: &models.Payload{
				KeyName:    "TESTKEY1",
				KeyVersion: 2,
				Data:       "test data",
			},
			wantErr: false,
		},
		{
			name: "Invalid base64 data",
			args: args{
				data: "invalidbase64data",
			},
			want: &models.Payload{
				Data: "invalidbase64data",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.err {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !tt.wantErr {
				if got.KeyName != tt.want.KeyName || got.KeyVersion != tt.want.KeyVersion || got.Data != tt.want.Data {
					t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
				}
				if got.KeyName == "TESTKEY" && got.Data != "" && got.Data != tt.want.Data {
					t.Errorf("Decrypt() EncryptedData should be equal to Data for valid key")
				}
			}
		})
	}
}
