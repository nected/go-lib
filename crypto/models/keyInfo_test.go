package models

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// freshKeyInfo returns a KeyInfo backed by a freshly generated 2048-bit RSA
// key. Tests need real keys so Encrypt/Decrypt can round-trip; the OAEP
// label is derived from KeyNameVersion() so matching name/version values
// between sender and receiver matter.
func freshKeyInfo(t *testing.T, name string, version int) *KeyInfo {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	return &KeyInfo{
		PrivKey: priv,
		PubKey:  &priv.PublicKey,
		Name:    name,
		Version: version,
	}
}

// Get/Set/KeyNameVersion are trivial field accessors. One table exercises
// each accessor pair; KeyNameVersion checks the "<name>_<version>" format.
func TestKeyInfoAccessors(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	k := &KeyInfo{}
	k.SetName("APP")
	k.SetVersion(3)
	k.SetPrivKey(*priv)
	k.SetPubKey(priv.PublicKey)

	assert.Equal(t, "APP", k.GetName())
	assert.Equal(t, 3, k.GetVersion())
	assert.NotNil(t, k.GetPrivKey())
	assert.NotNil(t, k.GetPubKey())
	assert.Equal(t, "APP_3", k.KeyNameVersion())
}

// Encrypt → Decrypt must round-trip across multiple payload sizes, including
// inputs larger than the RSA-OAEP step (key size - 2*hash size - 2) which
// forces multi-block encryption inside Encrypt.
func TestKeyInfoEncryptDecryptRoundTrip(t *testing.T) {
	k := freshKeyInfo(t, "ROUND", 1)

	tests := []struct {
		name string
		data []byte
	}{
		{name: "small", data: []byte("hello")},
		{name: "empty payload", data: []byte("")},
		// 2048-bit key + SHA-512 OAEP allows up to 2048/8 - 2*64 - 2 = 126 bytes per block.
		// 300 bytes guarantees multiple blocks.
		{name: "multi block", data: make([]byte, 300)},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			encrypted, err := k.Encrypt(test.data)
			assert.NoError(t, err)
			decrypted, err := k.Decrypt(encrypted)
			assert.NoError(t, err)
			assert.Equal(t, test.data, decrypted)
		})
	}
}

// Decrypt surfaces the underlying rsa.DecryptOAEP error when the ciphertext
// is malformed (wrong key, garbage bytes, etc.).
func TestKeyInfoDecrypt_InvalidData(t *testing.T) {
	k := freshKeyInfo(t, "BADDATA", 1)
	// Ciphertext block must be exactly PubKey.Size() bytes for OAEP — zeros
	// of that length parse as a block but fail integrity check.
	_, err := k.Decrypt(make([]byte, k.PubKey.Size()))
	assert.Error(t, err)
}

// GetEncryptKeysMap returns the package-level singleton (or nil if never set).
// SetEncryptKeysMap installs a new map. The store is process-wide, so this
// test resets it at the end to avoid cross-package contamination.
func TestEncryptKeysMapGetSet(t *testing.T) {
	prev := GetEncryptKeysMap()
	t.Cleanup(func() { SetEncryptKeysMap(prev) })

	SetEncryptKeysMap(nil)
	assert.Nil(t, GetEncryptKeysMap())

	fresh := &EncryptStruct{AvailableKeys: map[string]map[int]KeyInfo{}}
	SetEncryptKeysMap(fresh)
	assert.Same(t, fresh, GetEncryptKeysMap())
}

// GetEncryptionKey looks up a key by name and version. version > 0 returns
// that specific version; version <= 0 returns the highest known version.
// A nil store or unknown name returns nil.
func TestGetEncryptionKey(t *testing.T) {
	prev := GetEncryptKeysMap()
	t.Cleanup(func() { SetEncryptKeysMap(prev) })

	// Nil store path: must return nil for any lookup.
	SetEncryptKeysMap(nil)
	assert.Nil(t, GetEncryptionKey("ANY", 1))

	// Seed the store with two versions of a single key.
	k1 := freshKeyInfo(t, "APP", 1)
	k2 := freshKeyInfo(t, "APP", 2)
	SetEncryptKeysMap(&EncryptStruct{
		AvailableKeys: map[string]map[int]KeyInfo{
			"APP": {1: *k1, 2: *k2},
		},
	})

	// Exact version match.
	got := GetEncryptionKey("APP", 1)
	if assert.NotNil(t, got) {
		assert.Equal(t, 1, got.GetVersion())
	}

	// Missing version under a known key name returns nil rather than a
	// zero-value KeyInfo (otherwise downstream Encrypt would deref nil).
	assert.Nil(t, GetEncryptionKey("APP", 99))

	// version <= 0 selects the highest available version (map iteration
	// order is randomized, so the loop must compare and keep the max).
	latest := GetEncryptionKey("APP", 0)
	if assert.NotNil(t, latest) {
		assert.Equal(t, 2, latest.GetVersion())
	}

	// Unknown key name returns nil.
	assert.Nil(t, GetEncryptionKey("MISSING", 1))
}
