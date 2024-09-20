package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nected/go-lib/crypto/config"
	"github.com/nected/go-lib/crypto/models"
)

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

func setupSuite(t *testing.T) func(t *testing.T) {
	var privateKey = generatePrivateKey()
	os.Setenv("KEY_TESTKEY_1_0", privateKey)
	os.Setenv("KEY_TESTKEYR_1_1726147578", privateKey)
	os.Setenv("KEY_TESTKEYINVALID_1_0", "lkajds")

	t.Log("setup suite")
	return func(t *testing.T) {
		defer os.Unsetenv("KEY_TESTKEY_1_0")
		defer os.Unsetenv("KEY_TESTKEYR_1_1726147578")
		defer os.Unsetenv("KEY_TESTKEYINVALID_1_0")
		t.Log("teardown suite")
	}
}

func TestLoadKeysFromEnv(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	rT := time.Unix(1726147578, 0)
	tests := []struct {
		name        string
		wantErr     bool
		keyName     string
		keyExists   bool
		version     string
		rotate      *time.Time
		privKeyNull bool
		pubKeyNull  bool
	}{
		{
			name:        "TestLoadKeysFromEnv - No errors",
			wantErr:     false,
			keyName:     "TESTKEY",
			keyExists:   true,
			version:     "1",
			rotate:      nil,
			privKeyNull: false,
			pubKeyNull:  false,
		},
		{
			name:        "TestLoadKeysFromEnv - No errors",
			wantErr:     false,
			keyName:     "TESTKEYR",
			keyExists:   true,
			version:     "1",
			rotate:      &rT,
			privKeyNull: false,
			pubKeyNull:  false,
		},
		{
			name:        "TestLoadKeysFromEnv - error",
			wantErr:     false,
			keyName:     "TESTKEYA",
			keyExists:   false,
			version:     "1",
			rotate:      nil,
			privKeyNull: true,
			pubKeyNull:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadKeysFromEnv()
			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadKeysFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("LoadKeysFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			info := models.GetEncryptInfo()
			if info == nil {
				t.Errorf("GetEncryptInfo() = %v, want %v", info, &models.EncryptStruct{
					AvailableKeys: make(map[string]map[string]models.KeyInfo),
				})
				return
			}

			if tt.keyExists {
				if _, ok := info.AvailableKeys[tt.keyName]; !ok {
					t.Errorf("Key %s not found in AvailableKeys", tt.keyName)
					return
				}
			} else {
				if _, ok := info.AvailableKeys[tt.keyName]; ok {
					t.Errorf("Key %s found in AvailableKeys", tt.keyName)
					return
				} else {
					return
				}
			}

			keyInfoVersionMap := info.AvailableKeys[tt.keyName]
			if keyInfo, ok := keyInfoVersionMap[tt.version]; !ok {
				if keyInfo.GetName() != tt.keyName {
					t.Errorf("GetName() = %v, want %v", keyInfo.GetName(), tt.keyName)
				}
				if keyInfo.GetVersion() != tt.version {
					t.Errorf("GetVersion() = %v, want %v", keyInfo.GetVersion(), tt.version)
				}

				if tt.privKeyNull {
					if keyInfo.GetPrivKey() != nil {
						t.Errorf("GetPrivKey() = %v, want nil", keyInfo.GetPrivKey())
					}
				} else {
					if keyInfo.GetPrivKey() == nil {
						t.Errorf("GetPrivKey() = %v, want not nil", keyInfo.GetPrivKey())
					}
				}

				if tt.pubKeyNull {
					if keyInfo.GetPubKey() != nil {
						t.Errorf("GetPubKey() = %v, want nil", keyInfo.GetPubKey())
					}
				} else {
					if keyInfo.GetPubKey() == nil {
						t.Errorf("GetPubKey() = %v, want not nil", keyInfo.GetPubKey())
					}
				}

				if tt.rotate != nil {
					if keyInfo.GetRotateAt() == nil {
						t.Errorf("GetRotateAt() = %v, want not nil", keyInfo.GetRotateAt())
					} else {
						if keyInfo.GetRotateAt().Unix() != tt.rotate.Unix() {
							t.Errorf("GetRotateAt() = %v, want %v", keyInfo.GetRotateAt().Unix(), tt.rotate.Unix())
						}
					}
				}

				if tt.rotate == nil {
					if keyInfo.GetRotateAt() != nil {
						t.Errorf("GetRotateAt() = %v, want nil", keyInfo.GetRotateAt())
					}
				}
			}

		})
	}
}

func TestGetEncryptionKey(t *testing.T) {
	teardownSuite := setupSuite(t)
	LoadKeysFromEnv()
	defer teardownSuite(t)
	type args struct {
		keyName string
	}
	tests := []struct {
		name string
		args args
		want *models.KeyInfo
	}{
		{
			name: "TestGetEncryptionKey - No errors",
			args: args{
				keyName: "TESTKEY",
			},
			want: &models.KeyInfo{
				PrivKey:   nil,
				PubKey:    nil,
				Name:      "TESTKEY",
				Version:   "1",
				CreatedAt: time.Time{},
				RotateAt:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.GetEncryptionKey(tt.args.keyName); got != nil {
				if got.GetName() != tt.want.GetName() {
					t.Errorf("GetName() = %v, want %v", got.GetName(), tt.want.GetName())
				}

				if got.GetVersion() != tt.want.GetVersion() {
					t.Errorf("GetVersion() = %v, want %v", got.GetVersion(), tt.want.GetVersion())
				}

				if got.GetPrivKey() == nil {
					t.Errorf("GetPrivKey() = %v, want not nil", got.GetPrivKey())
				}

				if got.GetPubKey() == nil {
					t.Errorf("GetPubKey() = %v, want not nil", got.GetPubKey())
				}
			}
		})
	}
}

func TestEncryptRSA(t *testing.T) {
	teardownSuite := setupSuite(t)
	LoadKeysFromEnv()
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
			name: "TestEncryptRSA - No errors",
			args: args{
				keyName: "TESTKEY",
				data:    []byte("data"),
			},
			want: &models.Payload{
				Data:          "data",
				EncryptedData: "9zE+AFpfb3PhIfdaOlPxXZAVHb3oEiTxMYcIoDuaYVs=",
			},
			wantErr: false,
		},
		{
			name: "TestEncryptRSA - No errors",
			args: args{
				keyName: "TESTKEYR",
				data:    []byte("data"),
			},
			want: &models.Payload{
				Data:          "data",
				EncryptedData: "9zE+AFpfb3PhIfdaOlPxXZAVHb3oEiTxMYcIoDuaYVs=",
			},
			wantErr: false,
		},
		{
			name: "TestEncryptRSA - error",
			args: args{
				keyName: "TESTKEYINVALID",
				data:    []byte("data"),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "TestEncryptRSA - plain text",
			args: args{
				keyName: "TESTKEYINVALID",
				data:    []byte("data"),
			},
			want:    &models.Payload{Data: "data"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptRSA(tt.args.keyName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptRSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				if got.Data != tt.want.Data {
					t.Errorf("EncryptRSA() = %v, want %v", got.Data, tt.want.Data)
				}

				payload, err := DecryptRSA(got.String())
				if err != nil {
					t.Errorf("DecryptRSA() error = %v", err)
					return
				}

				if payload.Data != tt.want.Data {
					t.Errorf("EncryptRSA() = %v, want %v", payload.Data, tt.want.Data)
				}
			}
		})
	}
}
