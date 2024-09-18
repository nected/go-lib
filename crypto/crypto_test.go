package crypto

import (
	"os"
	"testing"
	"time"

	"github.com/nected/go-lib/crypto/common"
)

func setupSuite(t *testing.T) func(t *testing.T) {
	var privateKey = "-----BEGIN PRIVATE KEY-----\nMIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEA1M+jjOGBlvxTfvQy\n9j+e7fW+viOqVFY7vTBjcllrPWxqpAhA4zGGEbCzEK90CKbNAXNMtxnKuqS5f5Sm\nKkesAQIDAQABAkEAqUXUAL5q1s80FvplAuxOHVdodlNmK5lAAVdY8t7fZ0W/ElhU\nRTQh0V2ZWFNMwkaBTxtWb/V2+Z3iuOcnjgUXAQIhAO5tkgpnPnHJ/m47F59ewZM0\nZZAYJfrv2+jTeROLrA6xAiEA5H7AyWGSNjO7RaAYdhtB16YbancRC+uz7Wv1Ea7k\n5lECIQC1RTS9GBWPqYT5BZBGKGJ/qlx1Gwb1K5tD/lOVGqGrYQIgf1dUweaqwaJb\nABaVC11teG2OYesxiN83S14bGlvKHcECIAXBMEUREDAEqOFrBFR3zi/m71in18d+\n5Gv+J6YWHl1B\n-----END PRIVATE KEY-----"
	os.Setenv("KEY_TESTKEY_1_0", privateKey)
	os.Setenv("KEY_TESTKEYR_1_1726147578", privateKey)
	os.Setenv("KEY_TESTKEYA_1_0", "lkajds")

	t.Log("setup suite")
	return func(t *testing.T) {
		defer os.Unsetenv("KEY_TESTKEY_1_0")
		defer os.Unsetenv("KEY_TESTKEYR_1_1726147578")
		defer os.Unsetenv("KEY_TESTKEYA_1_0")
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
			info := getEncryptInfo()
			if info == nil {
				t.Errorf("GetEncryptInfo() = %v, want %v", info, &EncryptStruct{
					AvailableKeys: make(map[string]map[string]common.KeyInfo),
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
		want *common.KeyInfo
	}{
		{
			name: "TestGetEncryptionKey - No errors",
			args: args{
				keyName: "TESTKEY",
			},
			want: &common.KeyInfo{
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
			if got := GetEncryptionKey(tt.args.keyName); got != nil {
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
