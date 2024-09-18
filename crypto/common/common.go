package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"time"
)

type KeyInfo struct {
	PrivKey   *rsa.PrivateKey
	PubKey    *rsa.PublicKey
	Name      string
	Version   string
	CreatedAt time.Time
	RotateAt  *time.Time
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

func (k *KeyInfo) GetVersion() string {
	return k.Version
}

func (k *KeyInfo) SetVersion(version string) {
	k.Version = version
}

func (k *KeyInfo) KeyNameVersion() string {
	return fmt.Sprintf("%s_%s", k.GetName(), k.GetVersion())
}

func (k *KeyInfo) GetCreatedAt() time.Time {
	return k.CreatedAt
}

func (k *KeyInfo) SetCreatedAt(createdAt time.Time) {
	k.CreatedAt = createdAt
}

func (k *KeyInfo) GetRotateAt() *time.Time {
	return k.RotateAt
}

func (k *KeyInfo) SetRotateAt(rotateAt *time.Time) {
	k.RotateAt = rotateAt
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
