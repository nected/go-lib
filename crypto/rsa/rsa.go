package rsa

import (
	"fmt"

	"github.com/nected/go-lib/crypto"
	"github.com/nected/go-lib/crypto/base64"
)

func EncryptRSA(keyName string, data []byte) (*crypto.Payload, error) {
	keyInfo := crypto.GetEncryptionKey(keyName)
	if keyInfo == nil {
		// if key not found return stringfied data
		return &crypto.Payload{
			Data: string(data),
		}, nil
	}

	encryptedData, err := keyInfo.Encrypt(data)
	if err != nil {
		return nil, err
	}

	encryptedDataString := base64.B64Encode(encryptedData)

	return &crypto.Payload{
		KeyName:       keyName,
		KeyVersion:    keyInfo.GetVersion(),
		KeyType:       crypto.KeyTypeRSA,
		Data:          string(data),
		EncryptedData: encryptedDataString,
	}, nil
}

func DecryptRSA(data string) (*crypto.Payload, error) {
	p := crypto.Payload{}
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
		p.Data = decodedData
		return &p, nil
	}

	// split data into keyName, keyVersion and encryptedData
	// $keyName$keyVersion$encryptedData

	keyName, keyVersion, encryptedData := parseData(decodedData)

	keyInfo := crypto.GetEncryptionKey(keyName)

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

	return &crypto.Payload{
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
