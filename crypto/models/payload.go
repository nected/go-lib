package models

import (
	"fmt"

	"github.com/nected/go-lib/crypto/base64"
)

type KeyType string

const (
	KeyTypeRSA KeyType = "RSA"
	KeyTypeAES KeyType = "AES"
)

type Payload struct {
	// Payload is the data to be encrypted
	KeyName          string  `json:"keyName"`
	KeyVersion       string  `json:"keyVersion"`
	KeyType          KeyType `json:"keyType"`
	Data             string  `json:"data"`
	EncryptedData    string  `json:"encryptedData"`
	AlreadyEncrypted bool    `json:"alreadyEncrypted"`
}

func (p *Payload) String() string {
	if p.KeyName == "" {
		return p.EncryptedData
	}
	if p.KeyType == KeyTypeAES {
		return p.EncryptedData
	}
	data := fmt.Sprintf("$%s$%s$%s", p.KeyName, p.KeyVersion, p.EncryptedData)
	return base64.B64Encode([]byte(data))
}
