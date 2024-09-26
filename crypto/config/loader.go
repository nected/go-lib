package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nected/go-lib/crypto/models"
)

// key config format

func LoadKeysFromFile(keyName, keyPath string) error {
	info := models.GetEncryptKeysMap()
	if info == nil {
		info = &models.EncryptStruct{
			AvailableKeys: make(map[string]map[string]models.KeyInfo),
		}
		models.SetEncryptKeysMap(info)
	}
	privateKey, err := loadPrivateKeyFromFile(keyPath)
	if err != nil {
		return err
	}
	publicKey := generatePublicKey(privateKey)
	keyInfo := models.KeyInfo{
		PrivKey:   privateKey,
		PubKey:    publicKey,
		Name:      keyName,
		Version:   "1",
		CreatedAt: time.Now(),
	}

	info.AvailableKeys[keyName][keyInfo.Version] = keyInfo
	return nil
}

// env key format
//
// KEY_<key_name>_<key_version>
//
// Parameters:
//   - key_name: name of the key
//   - key_version(optional): version of the key, default is 1
//
// Example:
//   - KEY_TESTKEY_1_0
func LoadKeysFromEnv() error {
	info := models.GetEncryptKeysMap()

	if info == nil {
		info = &models.EncryptStruct{
			AvailableKeys: make(map[string]map[string]models.KeyInfo),
		}
		models.SetEncryptKeysMap(info)
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
			info.AvailableKeys[keyName] = make(map[string]models.KeyInfo)
		}

		info.AvailableKeys[keyName][keyVersion] = models.KeyInfo{
			Name:      keyName,
			Version:   keyVersion,
			CreatedAt: time.Now(),
			PrivKey:   privateKey,
			PubKey:    generatePublicKey(privateKey),
		}

	}
	return nil
}
