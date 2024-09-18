package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func loadPrivateKeyFromFile(keyPath string) (*rsa.PrivateKey, error) {
	fileData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return loadPrivateKey(fileData)
}

func loadPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	// Parse the key
	privKeyPem, _ := pem.Decode(pemData)
	if privKeyPem == nil {
		return nil, fmt.Errorf("Invalid private key")
	}
	if privKeyPem.Type != PRIV_KEY_TYPE {
		return nil, fmt.Errorf("Invalid private key type")
	}

	privKeyParsed, err := x509.ParsePKCS8PrivateKey(privKeyPem.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := privKeyParsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Invalid private key")
	}

	return privateKey, nil
}

func generatePublicKey(privateKey *rsa.PrivateKey) *rsa.PublicKey {
	privateKey.Precompute()
	return &privateKey.PublicKey
}
