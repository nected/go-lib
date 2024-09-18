package crypto

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nected/go-lib/crypto/common"
)

// key config format

const (
	KEY_ENV_PREFIX = "KEY_"
	PRIV_KEY_TYPE  = "PRIVATE KEY"
)

var encryptInfo *EncryptStruct

type EncryptStruct struct {
	AvailableKeys map[string]map[string]common.KeyInfo
}

func getEncryptInfo() *EncryptStruct {
	return encryptInfo
}

func GetEncryptionKey(keyName string) *common.KeyInfo {
	info := getEncryptInfo()
	if info == nil {
		return nil
	}
	keyInfoVersionMap, ok := info.AvailableKeys[keyName]
	if !ok {
		return nil
	}
	var latestKeyInfo *common.KeyInfo
	for _, keyInfo := range keyInfoVersionMap {
		if latestKeyInfo == nil {
			latestKeyInfo = &keyInfo
		}
		if keyInfo.GetCreatedAt().After(latestKeyInfo.GetCreatedAt()) {
			latestKeyInfo = &keyInfo
		}
	}
	return latestKeyInfo
}

func setEncryptInfo(info *EncryptStruct) {
	encryptInfo = info
}

func LoadKeysFromFile(keyName, keyPath string, rotateAt *time.Time) error {
	info := getEncryptInfo()
	if info == nil {
		info = &EncryptStruct{
			AvailableKeys: make(map[string]map[string]common.KeyInfo),
		}
		setEncryptInfo(info)
	}
	privateKey, err := loadPrivateKeyFromFile(keyPath)
	if err != nil {
		return err
	}
	publicKey := generatePublicKey(privateKey)
	keyInfo := common.KeyInfo{
		PrivKey:   privateKey,
		PubKey:    publicKey,
		Name:      keyName,
		Version:   "1",
		CreatedAt: time.Now(),
		RotateAt:  rotateAt,
	}

	info.AvailableKeys[keyName][keyInfo.Version] = keyInfo
	return nil
}

// env key format
//
// KEY_<key_name>_<key_version>_<rotate_time_milli>
//
// Parameters:
//   - key_name: name of the key
//   - key_version(optional): version of the key, default is 1
//   - rotate_time_milli(optional): time in milliseconds when the key should be rotated, default is 0
//
// Example:
//   - KEY_TESTKEY_1_0
//   - KEY_TESTKEY_2_0_1614556800000
func LoadKeysFromEnv() error {
	info := getEncryptInfo()

	if info == nil {
		info = &EncryptStruct{
			AvailableKeys: make(map[string]map[string]common.KeyInfo),
		}
		setEncryptInfo(info)
	}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, KEY_ENV_PREFIX) {
			continue
		}

		// vals := strings.Split(env, "=")
		key := ""
		// value := vals[1]
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				key = env[:i]
				break
			}
		}
		value, ok := os.LookupEnv(key)
		if !ok {
			// key not found
			continue
		}
		parts := strings.Split(key, "_")
		if len(parts) < 2 {
			// invalid key format
			// key name is missing
			continue
		}

		keyName := parts[1]

		keyVersion := "1"
		if len(parts) >= 3 {
			keyVersion = parts[2]
		}

		var rotateAt *time.Time
		if len(parts) >= 4 {
			rorateAtInt, err := strconv.ParseInt(parts[3], 10, 64)
			if err != nil {
				// invalid rotate time format
				// skip this key
				continue
			}
			if rorateAtInt <= 0 {
				rotateAt = nil
			} else {
				rotateAtTime := time.Unix(rorateAtInt, 0)
				rotateAt = &rotateAtTime
			}
		}

		privateKey, err := loadPrivateKey([]byte(value))
		if err != nil {
			fmt.Println("Invalid private key")
			continue
		}

		if err := privateKey.Validate(); err != nil {
			fmt.Println("Invalid private key")
			continue
		}

		if _, ok := info.AvailableKeys[keyName]; !ok {
			info.AvailableKeys[keyName] = make(map[string]common.KeyInfo)
		}

		info.AvailableKeys[keyName][keyVersion] = common.KeyInfo{
			Name:      keyName,
			Version:   keyVersion,
			CreatedAt: time.Now(),
			RotateAt:  rotateAt,
			PrivKey:   privateKey,
			PubKey:    generatePublicKey(privateKey),
		}

	}
	return nil
}
