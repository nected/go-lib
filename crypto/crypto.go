package crypto

import (
	"github.com/nected/go-lib/crypto/aes"
	"github.com/nected/go-lib/crypto/config"
	"github.com/nected/go-lib/crypto/models"
	"github.com/nected/go-lib/crypto/rsa"
)

func EncryptRSA(keyName string, data []byte) (*models.Payload, error) {
	return rsa.Encrypt(keyName, data)
}

func DecryptRSA(data string) (*models.Payload, error) {
	return rsa.Decrypt(data)
}

func EncryptAES(secret string, data []byte) (*models.Payload, error) {
	return aes.Encrypt(secret, data)
}

func DecryptAES(secret string, data string) (*models.Payload, error) {
	return aes.Decrypt(secret, data)
}

func LoadKeysFromEnv() error {
	return config.LoadKeysFromEnv()
}

func ListKeys() *models.EncryptStruct {
	return models.GetEncryptInfo()
}
