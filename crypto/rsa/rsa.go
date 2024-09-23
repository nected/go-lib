package rsa

import (
	"fmt"

	"github.com/nected/go-lib/crypto/base64"
	"github.com/nected/go-lib/crypto/errors"
	"github.com/nected/go-lib/crypto/models"
)

// Encrypt encrypts the given data using the specified encryption key.
// If the data is already encrypted, it returns the data as is.
// If the key is not found, it returns the data as a string without encryption.
//
// Parameters:
//   - keyName: The name of the encryption key to use.
//   - data: The data to be encrypted.
//
// Returns:
//   - *models.Payload: A payload containing the original data and the encrypted data.
//   - error: An error if the encryption process fails.
func Encrypt(keyName string, data []byte) (*models.Payload, error) {
	if alreadyEncrypted(data) {
		return &models.Payload{
			Data:             string(data),
			EncryptedData:    string(data),
			AlreadyEncrypted: true,
		}, nil
	}
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

// Decrypt decrypts the given base64-encoded data string and returns a Payload object.
// If the data is not encrypted, it returns the data as is.
//
// Parameters:
//   - data: A base64-encoded string that may contain encrypted data.
//
// Returns:
//   - *models.Payload: A pointer to the Payload object containing the decrypted data.
//   - error: An error object if any error occurs during decryption.
//
// The function performs the following steps:
//  1. Checks if the input data string is empty and returns nil if true.
//  2. Decodes the base64-encoded input data.
//  3. Checks if the decoded data is not encrypted and returns it as is if true.
//  4. Parses the decoded data to extract keyName, keyVersion, and encryptedData.
//  5. Retrieves the encryption key information using the keyName.
//  6. Decodes the base64-encoded encryptedData.
//  7. Decrypts the encryptedData using the retrieved key information.
//  8. Constructs and returns a Payload object containing the decrypted data and other relevant information.
func Decrypt(data string) (*models.Payload, error) {
	p := models.Payload{}
	if data == "" {
		return nil, errors.ErrEmptyData
	}

	// decode data
	decodedData, err := base64.B64Decode(data)
	if err != nil {
		p.Data = data
		return &p, nil
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
		EncryptedData: base64.B64Encode([]byte(encryptedData)),
	}, nil
}

// parseData parses a string containing key name, key version, and encrypted data
// separated by '$' characters. It returns the key name, key version, and encrypted data.
//
// The input string is expected to be in the format: "$keyName$keyVersion$encryptedData".
//
// Parameters:
// - data: A string containing the key name, key version, and encrypted data.
//
// Returns:
// - keyName: The extracted key name.
// - keyVersion: The extracted key version.
// - encryptedData: The extracted encrypted data.
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

// alreadyEncrypted checks if the provided data is already encrypted.
// It attempts to decode the data from base64 and then parse it to extract a key name.
// If the key name is not empty, it returns true, indicating that the data is encrypted.
// Otherwise, it returns false.
//
// Parameters:
// - data: A byte slice containing the data to check.
//
// Returns:
// - bool: True if the data is already encrypted, false otherwise.
func alreadyEncrypted(data []byte) bool {
	decodedData, err := base64.B64Decode(string(data))
	if err != nil {
		return false
	}
	keyName, _, _ := parseData(decodedData)
	return keyName != ""
}
