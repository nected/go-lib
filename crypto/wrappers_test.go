package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/nected/go-lib/crypto/config"
	"github.com/nected/go-lib/crypto/models"
	"github.com/stretchr/testify/assert"
)

func writeTempPEM(t *testing.T) string {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	assert.NoError(t, err)
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: config.PRIV_KEY_TYPE, Bytes: der})
	path := filepath.Join(t.TempDir(), "key.pem")
	assert.NoError(t, os.WriteFile(path, pemBytes, 0600))
	return path
}

// EncryptAES → DecryptAES with the same 16-byte secret must round-trip.
// Empty data short-circuits to (nil, nil) inside aes.Encrypt.
func TestEncryptAES_RoundTrip(t *testing.T) {
	secret := "0123456789abcdef" // 16 bytes → AES-128
	payload, err := EncryptAES(secret, []byte("hello world"))
	assert.NoError(t, err)
	if assert.NotNil(t, payload) {
		assert.Equal(t, models.KeyTypeAES, payload.KeyType)
		assert.NotEmpty(t, payload.EncryptedData)
	}

	dec, err := DecryptAES(secret, payload.EncryptedData)
	assert.NoError(t, err)
	if assert.NotNil(t, dec) {
		assert.Equal(t, "hello world", dec.Data)
	}
}

// Bad secret length (not 16/24/32) hits aes.NewCipher's KeySizeError on both
// Encrypt and Decrypt.
func TestEncryptAES_BadSecret(t *testing.T) {
	_, err := EncryptAES("short", []byte("data"))
	assert.Error(t, err)
	_, err = DecryptAES("short", "anything")
	assert.Error(t, err)
}

// DecryptAES short-circuits to (nil, nil) on empty ciphertext.
func TestDecryptAES_EmptyData(t *testing.T) {
	got, err := DecryptAES("0123456789abcdef", "")
	assert.NoError(t, err)
	assert.Nil(t, got)
}

// EncryptAES short-circuits to (nil, nil) on empty plaintext.
func TestEncryptAES_EmptyData(t *testing.T) {
	got, err := EncryptAES("0123456789abcdef", []byte(""))
	assert.NoError(t, err)
	assert.Nil(t, got)
}

// LoadKeysFromFile is a thin wrapper around config.LoadKeysFromFile and must
// insert the key into the models store. The pre-seed is required because
// the underlying loader has a latent nil-map bug (see crypto/config tests).
func TestLoadKeysFromFile_Wrapper(t *testing.T) {
	prev := models.GetEncryptKeysMap()
	t.Cleanup(func() { models.SetEncryptKeysMap(prev) })
	models.SetEncryptKeysMap(&models.EncryptStruct{
		AvailableKeys: map[string]map[int]models.KeyInfo{"FILEKEY": {}},
	})

	path := writeTempPEM(t)
	assert.NoError(t, LoadKeysFromFile("FILEKEY", path))
	assert.NotNil(t, models.GetEncryptionKey("FILEKEY", 1))
}

// ListKeys returns the package-level models store pointer. A nil store
// surfaces as nil; a populated store surfaces with the same address.
func TestListKeys(t *testing.T) {
	prev := models.GetEncryptKeysMap()
	t.Cleanup(func() { models.SetEncryptKeysMap(prev) })

	models.SetEncryptKeysMap(nil)
	assert.Nil(t, ListKeys())

	es := &models.EncryptStruct{AvailableKeys: map[string]map[int]models.KeyInfo{}}
	models.SetEncryptKeysMap(es)
	assert.Same(t, es, ListKeys())
}
