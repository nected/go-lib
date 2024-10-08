package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"

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
	os.Setenv(fmt.Sprintf("%s_%s", config.KEY_ENV_PREFIX, "TESTKEY_1"), privateKey)
	os.Setenv(fmt.Sprintf("%s_%s", config.KEY_ENV_PREFIX, "TESTKEY_2"), privateKey)
	os.Setenv(fmt.Sprintf("%s_%s", config.KEY_ENV_PREFIX, "TESTKEYINVALID"), "lkajds")

	t.Log("setup suite")
	return func(t *testing.T) {
		defer os.Unsetenv(fmt.Sprintf("%s_%s", config.KEY_ENV_PREFIX, "TESTKEY_1"))
		defer os.Unsetenv(fmt.Sprintf("%s_%s", config.KEY_ENV_PREFIX, "TESTKEY1_2"))
		defer os.Unsetenv(fmt.Sprintf("%s_%s", config.KEY_ENV_PREFIX, "TESTKEYINVALID"))
		t.Log("teardown suite")
	}

}

func TestLoadKeysFromEnv(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	tests := []struct {
		name        string
		wantErr     bool
		keyName     string
		keyExists   bool
		version     int
		privKeyNull bool
		pubKeyNull  bool
	}{
		{
			name:        "TestLoadKeysFromEnv version 1 - No errors",
			wantErr:     false,
			keyName:     "TESTKEY",
			keyExists:   true,
			version:     1,
			privKeyNull: false,
			pubKeyNull:  false,
		},
		{
			name:        "TestLoadKeysFromEnv version 2 - No errors",
			wantErr:     false,
			keyName:     "TESTKEY",
			keyExists:   true,
			version:     2,
			privKeyNull: false,
			pubKeyNull:  false,
		},
		{
			name:        "TestLoadKeysFromEnv - error",
			wantErr:     false,
			keyName:     "TESTKEYA",
			keyExists:   false,
			version:     1,
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
			info := models.GetEncryptionKey(tt.keyName, tt.version)
			if tt.keyExists {
				if info == nil {
					t.Errorf("Key %s not found in AvailableKeys", tt.keyName)
					return
				}
			} else {
				if info != nil {
					t.Errorf("Key %s found in AvailableKeys", tt.keyName)
				}
				return
			}

			if info.GetName() != tt.keyName {
				t.Errorf("GetName() = %v, want %v", info.GetName(), tt.keyName)
			}
			if info.GetVersion() != tt.version {
				t.Errorf("GetVersion() = %v, want %v", info.GetVersion(), tt.version)
			}

			if tt.privKeyNull {
				if info.GetPrivKey() != nil {
					t.Errorf("GetPrivKey() = %v, want nil", info.GetPrivKey())
				}
			} else {
				if info.GetPrivKey() == nil {
					t.Errorf("GetPrivKey() = %v, want not nil", info.GetPrivKey())
				}
			}

			if tt.pubKeyNull {
				if info.GetPubKey() != nil {
					t.Errorf("GetPubKey() = %v, want nil", info.GetPubKey())
				}
			} else {
				if info.GetPubKey() == nil {
					t.Errorf("GetPubKey() = %v, want not nil", info.GetPubKey())
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
		keyName    string
		keyVersion int
	}
	tests := []struct {
		name string
		args args
		want *models.KeyInfo
	}{
		{
			name: "TestGetEncryptionKey - No errors",
			args: args{
				keyName:    "TESTKEY",
				keyVersion: 1,
			},
			want: &models.KeyInfo{
				PrivKey: nil,
				PubKey:  nil,
				Name:    "TESTKEY",
				Version: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.GetEncryptionKey(tt.args.keyName, tt.args.keyVersion); got != nil {
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
				Data: "data",
			},
			wantErr: false,
		},
		{
			name: "TestEncryptRSA - similar pattern",
			args: args{
				keyName: "TESTKEY",
				data:    []byte("JGFiY2QkZGVmZyRsYWtzamxhc2RqYWxzZGo="), //$abcd$defg$laksjlasdjalsdj
			},
			want: &models.Payload{
				Data: "JGFiY2QkZGVmZyRsYWtzamxhc2RqYWxzZGo=",
			},
			wantErr: false,
		},
		{
			name: "TestEncryptRSA - Conflicting pattern",
			args: args{
				keyName: "TESTKEY",
				data:    []byte("QCRhYmNkJGRlZmckbGFrc2psYXNkamFsc2Rq"), //@$abcd$defg$laksjlasdjalsdj
			},
			want: &models.Payload{
				Data: "QCRhYmNkJGRlZmckbGFrc2psYXNkamFsc2Rq",
			},
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

			if !tt.wantErr && (got == nil || got.Data != tt.want.Data) {
				t.Errorf("EncryptRSA() = %v, want %v", got, tt.want)
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

				if got.KeyName != "" {
					if got.EncryptedData == "" {
						t.Errorf("Invalid encrypted data")
					}
				}
			}
		})
	}
}
