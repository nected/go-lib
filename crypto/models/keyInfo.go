package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
)

type EncryptStruct struct {
	AvailableKeys map[string]map[int]KeyInfo
}

var encryptKeysMap *EncryptStruct

type KeyInfo struct {
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
	Name    string
	Version int
}

func (k *KeyInfo) GetPrivKey() *rsa.PrivateKey {
	return k.PrivKey
}

func (k *KeyInfo) SetPrivKey(privKey rsa.PrivateKey) {
	k.PrivKey = &privKey
}

func (k *KeyInfo) GetPubKey() *rsa.PublicKey {
	return k.PubKey
}

func (k *KeyInfo) SetPubKey(pubKey rsa.PublicKey) {
	k.PubKey = &pubKey
}

func (k *KeyInfo) GetName() string {
	return k.Name
}

func (k *KeyInfo) SetName(name string) {
	k.Name = name
}

func (k *KeyInfo) GetVersion() int {
	return k.Version
}

func (k *KeyInfo) SetVersion(version int) {
	k.Version = version
}

func (k *KeyInfo) KeyNameVersion() string {
	return fmt.Sprintf("%s_%v", k.GetName(), k.GetVersion())
}

func (k *KeyInfo) Encrypt(data []byte) ([]byte, error) {
	msgLen := len(data)
	encryptHash := sha512.New()
	step := k.PubKey.Size() - 2*encryptHash.Size() - 2
	encryptedData := make([]byte, 0)
	for i := 0; i < msgLen; i += step {
		end := i + step
		if end > msgLen {
			end = msgLen
		}
		encrypted, err := rsa.EncryptOAEP(encryptHash, rand.Reader, k.GetPubKey(), data[i:end], []byte(k.KeyNameVersion()))
		if err != nil {
			return nil, err
		}
		encryptedData = append(encryptedData, encrypted...)
	}
	return encryptedData, nil
}

func (k *KeyInfo) Decrypt(data []byte) ([]byte, error) {
	msgLen := len(data)
	decryptHash := sha512.New()
	step := k.PubKey.Size()
	decryptedData := make([]byte, 0)
	for i := 0; i < msgLen; i += step {
		end := i + step
		if end > msgLen {
			end = msgLen
		}
		decrypted, err := rsa.DecryptOAEP(decryptHash, rand.Reader, k.GetPrivKey(), data[i:end], []byte(k.KeyNameVersion()))
		if err != nil {
			return nil, err
		}
		decryptedData = append(decryptedData, decrypted...)
	}
	return decryptedData, nil
}

func GetEncryptKeysMap() *EncryptStruct {
	return encryptKeysMap
}

func SetEncryptKeysMap(info *EncryptStruct) {
	encryptKeysMap = info
}

func GetEncryptionKey(keyName string, version int) *KeyInfo {
	info := GetEncryptKeysMap()
	if info == nil {
		return nil
	}
	keyInfoVersionMap, ok := info.AvailableKeys[keyName]
	if !ok {
		return nil
	}
	if version > 0 {
		keyInfo := keyInfoVersionMap[version]
		return &keyInfo
	}
	var latestKeyInfo *KeyInfo
	for _, keyInfo := range keyInfoVersionMap {
		if latestKeyInfo == nil {
			latestKeyInfo = &keyInfo
		}

		if keyInfo.GetVersion() > latestKeyInfo.GetVersion() {
			latestKeyInfo = &keyInfo
		}
	}
	return latestKeyInfo
}
