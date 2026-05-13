package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nected/go-lib/crypto/models"
	"github.com/stretchr/testify/assert"
)

// resetKeyStore swaps the global models map for a fresh one and restores
// the previous value on test teardown.
func resetKeyStore(t *testing.T) {
	t.Helper()
	prev := models.GetEncryptKeysMap()
	t.Cleanup(func() { models.SetEncryptKeysMap(prev) })
	models.SetEncryptKeysMap(nil)
}

// LoadKeysFromFile reads a key off disk and inserts it into the models
// global. The function only succeeds when AvailableKeys[keyName] already
// exists — see TestLoadKeysFromFile_NewKeyPanicsBug for the latent bug.
func TestLoadKeysFromFile(t *testing.T) {
	resetKeyStore(t)
	_, pemBytes := pemEncodeRSA(t)
	path := filepath.Join(t.TempDir(), "key.pem")
	assert.NoError(t, os.WriteFile(path, pemBytes, 0600))

	// Pre-seed the inner map for "APP" to avoid the nil-map assignment bug.
	models.SetEncryptKeysMap(&models.EncryptStruct{
		AvailableKeys: map[string]map[int]models.KeyInfo{
			"APP": {},
		},
	})

	assert.NoError(t, LoadKeysFromFile("APP", path))

	got := models.GetEncryptionKey("APP", 1)
	if assert.NotNil(t, got) {
		assert.Equal(t, "APP", got.GetName())
		assert.Equal(t, 1, got.GetVersion())
		assert.NotNil(t, got.GetPrivKey())
		assert.NotNil(t, got.GetPubKey())
	}
}

// File-read error: the path does not exist, ReadFile returns an error,
// the function returns it without touching the key store.
func TestLoadKeysFromFile_MissingFile(t *testing.T) {
	resetKeyStore(t)
	models.SetEncryptKeysMap(&models.EncryptStruct{
		AvailableKeys: map[string]map[int]models.KeyInfo{"APP": {}},
	})
	err := LoadKeysFromFile("APP", filepath.Join(t.TempDir(), "missing.pem"))
	assert.Error(t, err)
}

// Bad PEM: file exists but contains junk; loadPrivateKey returns "Invalid
// private key" and LoadKeysFromFile propagates it.
func TestLoadKeysFromFile_InvalidPEM(t *testing.T) {
	resetKeyStore(t)
	models.SetEncryptKeysMap(&models.EncryptStruct{
		AvailableKeys: map[string]map[int]models.KeyInfo{"APP": {}},
	})
	path := filepath.Join(t.TempDir(), "junk.pem")
	assert.NoError(t, os.WriteFile(path, []byte("not pem"), 0600))
	err := LoadKeysFromFile("APP", path)
	assert.Error(t, err)
}

// LoadKeysFromFile initializes the inner map on demand when the key name
// isn't already present. Before the fix at loader.go:34, this path wrote
// into a nil map and panicked.
func TestLoadKeysFromFile_NewKey(t *testing.T) {
	resetKeyStore(t)
	_, pemBytes := pemEncodeRSA(t)
	path := filepath.Join(t.TempDir(), "key.pem")
	assert.NoError(t, os.WriteFile(path, pemBytes, 0600))

	// Fresh global with no entry for the key name — the function must lazily
	// create the inner map and store the key under version 1.
	assert.NoError(t, LoadKeysFromFile("FRESH", path))
	got := models.GetEncryptionKey("FRESH", 1)
	if assert.NotNil(t, got) {
		assert.Equal(t, "FRESH", got.GetName())
		assert.Equal(t, 1, got.GetVersion())
	}
}

// LoadKeysFromEnv walks os.Environ() for entries prefixed with ENCRYPTKEY_.
// Key name comes from index 1 of the split parts; version (optional) from
// index 2. Bad PEMs, invalid versions, and malformed names are silently
// skipped. Valid entries land in models.AvailableKeys.
func TestLoadKeysFromEnv(t *testing.T) {
	resetKeyStore(t)
	_, pemBytes := pemEncodeRSA(t)
	pemStr := string(pemBytes)

	type envVar struct {
		key, value string
	}
	tests := []struct {
		name        string
		envs        []envVar
		wantName    string
		wantVersion int
		wantNil     bool // when true, no key should appear
	}{
		{
			name:        "valid versioned key",
			envs:        []envVar{{KEY_ENV_PREFIX + "_APPONE_3", pemStr}},
			wantName:    "APPONE",
			wantVersion: 3,
		},
		{
			// No version segment → defaults to 1.
			name:        "valid unversioned defaults to v1",
			envs:        []envVar{{KEY_ENV_PREFIX + "_APPTWO", pemStr}},
			wantName:    "APPTWO",
			wantVersion: 1,
		},
		{
			// Non-numeric version segment is rejected silently.
			name:        "invalid version is skipped",
			envs:        []envVar{{KEY_ENV_PREFIX + "_APPTHREE_abc", pemStr}},
			wantName:    "APPTHREE",
			wantNil:     true,
			wantVersion: 1,
		},
		{
			// Garbage PEM is rejected silently.
			name:        "invalid pem is skipped",
			envs:        []envVar{{KEY_ENV_PREFIX + "_APPFOUR_1", "not-pem"}},
			wantName:    "APPFOUR",
			wantNil:     true,
			wantVersion: 1,
		},
		{
			// Empty version segment ("APPFIVE_") keeps the default of 1.
			name:        "empty version segment defaults to 1",
			envs:        []envVar{{KEY_ENV_PREFIX + "_APPFIVE_", pemStr}},
			wantName:    "APPFIVE",
			wantVersion: 1,
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			resetKeyStore(t)
			for _, e := range test.envs {
				t.Setenv(e.key, e.value)
			}
			assert.NoError(t, LoadKeysFromEnv())
			got := models.GetEncryptionKey(test.wantName, test.wantVersion)
			if test.wantNil {
				assert.Nil(t, got)
				return
			}
			if assert.NotNil(t, got) {
				assert.Equal(t, test.wantName, got.GetName())
				assert.Equal(t, test.wantVersion, got.GetVersion())
			}
		})
	}
}

// The "no key name" branch: env like "ENCRYPTKEY=value" splits to ["ENCRYPTKEY"]
// (length 1), which fails the `len(parts) < 2` guard and is skipped silently.
func TestLoadKeysFromEnv_MissingKeyName(t *testing.T) {
	resetKeyStore(t)
	t.Setenv(KEY_ENV_PREFIX, "anything")
	assert.NoError(t, LoadKeysFromEnv())
	// Nothing should be inserted.
	m := models.GetEncryptKeysMap()
	if m != nil {
		assert.Empty(t, m.AvailableKeys)
	}
}
