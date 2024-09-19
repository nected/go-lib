package rsa

import (
	"fmt"

	"github.com/nected/go-lib/crypto/base64"
	"github.com/nected/go-lib/crypto/models"
)

func Encrypt(keyName string, data []byte) (*models.Payload, error) {
	keyInfo := models.GetEncryptionKey(keyName)
	if keyInfo == nil {
		// if key not found return stringfied data
		return &models.Payload{
			Data:          string(data),
			EncryptedData: string(data),
		}, nil
	}

	encryptedData, err := keyInfo.Encrypt(data)
	if err != nil {
		return nil, err
	}

	encryptedDataString := base64.B64Encode(encryptedData)

	return &models.Payload{
		KeyName:       keyName,
		KeyVersion:    keyInfo.GetVersion(),
		KeyType:       models.KeyTypeRSA,
		Data:          string(data),
		EncryptedData: encryptedDataString,
	}, nil
}

func Decrypt(data string) (*models.Payload, error) {
	p := models.Payload{}
	if data == "" {
		return nil, nil
	}

	// decode data
	decodedData, err := base64.B64Decode(data)
	if err != nil {
		return nil, err
	}

	// if data is not encrypted return as is
	if decodedData[0] != '$' {
		p.Data = data
		return &p, nil
	}

	// split data into keyName, keyVersion and encryptedData
	// $keyName$keyVersion$encryptedData

	keyName, keyVersion, encryptedData := parseData(decodedData)

	keyInfo := models.GetEncryptionKey(keyName)

	if keyInfo == nil {
		return nil, fmt.Errorf("key %s not found", keyName)
	}

	encryptedData, err = base64.B64Decode(encryptedData)

	if err != nil {
		return nil, err
	}

	decryptedData, err := keyInfo.Decrypt([]byte(encryptedData))
	if err != nil {
		return nil, err
	}

	return &models.Payload{
		KeyName:       keyName,
		KeyVersion:    keyVersion,
		Data:          string(decryptedData),
		EncryptedData: encryptedData,
	}, nil
}

func parseData(data string) (string, string, string) {
	keyName := ""
	keyVersion := ""
	encryptedData := ""

	for i := 1; i < len(data); i++ {
		if data[i] == '$' {
			if keyName == "" {
				keyName = data[1:i]
				continue
			}
			if keyVersion == "" {
				keyVersion = data[len(keyName)+2 : i]
				encryptedData = data[i+1:]
				break
			}
		}
	}

	return keyName, keyVersion, encryptedData
}
