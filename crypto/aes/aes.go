package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/nected/go-lib/crypto/base64"
	"github.com/nected/go-lib/crypto/models"
)

func Encrypt(secret string, data []byte) (*models.Payload, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	// fill nonce with random data
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	// encrypt data
	encryptedData := gcm.Seal(nonce, nonce, data, nil)

	encryptedDataString := base64.B64Encode(encryptedData)

	return &models.Payload{
		KeyType:       models.KeyTypeAES,
		Data:          string(data),
		EncryptedData: encryptedDataString,
	}, nil
}

func Decrypt(secret string, data string) (*models.Payload, error) {
	p := models.Payload{}
	if data == "" {
		return nil, nil
	}

	// decode data
	decodedData, err := base64.B64Decode(data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(decodedData) < nonceSize {
		return nil, err
	}

	nonce, encryptedData := decodedData[:nonceSize], decodedData[nonceSize:]

	decryptedData, err := gcm.Open(nil, []byte(nonce), []byte(encryptedData), nil)
	if err != nil {
		return nil, err
	}

	p.Data = string(decryptedData)
	return &p, nil
}
