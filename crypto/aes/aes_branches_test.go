package aes

import (
	"fmt"
	"testing"

	"github.com/nected/go-lib/crypto/base64"
	cryptoerrors "github.com/nected/go-lib/crypto/errors"
	"github.com/stretchr/testify/assert"
)

// Decrypt with a non-AES-keysize secret must fail at aes.NewCipher. The
// existing TestDecryptAES "Invalid secret" case feeds data that contains
// '+' (StdEncoding), which makes B64DecodeURL fail first — so the
// aes.NewCipher branch was never reached. This case uses a b64url-valid
// ciphertext to exercise the cipher-construction failure.
func TestDecrypt_BadSecretLength(t *testing.T) {
	// Generate a valid ciphertext with a valid secret so the b64 path passes.
	good, err := Encrypt("0123456789abcdef", []byte("payload"))
	assert.NoError(t, err)

	// Now decrypt with a wrong-length secret — aes.NewCipher fails.
	got, err := Decrypt("1", good.EncryptedData)
	assert.Error(t, err)
	assert.Nil(t, got)
}

// Decrypt with a valid-length secret but the wrong key value must fail at
// gcm.Open (authentication tag mismatch). Covers the final error branch
// in Decrypt that the existing tests skip.
func TestDecrypt_WrongKey(t *testing.T) {
	good, err := Encrypt("0123456789abcdef", []byte("payload"))
	assert.NoError(t, err)

	// Same length, different bytes → gcm.Open returns "cipher: message
	// authentication failed".
	got, err := Decrypt("fedcba9876543210", good.EncryptedData)
	assert.Error(t, err)
	assert.Nil(t, got)
}

// Decrypt rejects ciphertext shorter than the GCM nonce size with
// cryptoerrors.ErrInvalidData. Before the fix at aes.go:70 this branch
// returned (nil, nil) because it reused the previously-cleared `err`.
func TestDecrypt_ShortNonceReturnsError(t *testing.T) {
	// 4 raw bytes → b64url encodes to 8 chars and decodes back to 4 bytes,
	// which is less than the GCM nonce size of 12.
	short := base64.B64EncodeURL([]byte{1, 2, 3, 4})
	got, err := Decrypt("0123456789abcdef", short)
	assert.Nil(t, got)
	assert.ErrorIs(t, err, cryptoerrors.ErrInvalidData)
}

// Encrypt → Decrypt round-trip across a few payload sizes, in the same
// package so the success path contributes to crypto/aes coverage rather
// than only to crypto/.
func TestEncryptDecrypt_RoundTripSamePackage(t *testing.T) {
	secret := "0123456789abcdef"
	tests := []struct {
		name string
		data []byte
	}{
		{name: "small", data: []byte("a")},
		{name: "ascii", data: []byte("hello world")},
		{name: "binary", data: []byte{0, 1, 2, 0xfe, 0xff}},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			enc, err := Encrypt(secret, test.data)
			assert.NoError(t, err)
			if assert.NotNil(t, enc) {
				assert.NotEmpty(t, enc.EncryptedData)
			}
			dec, err := Decrypt(secret, enc.EncryptedData)
			assert.NoError(t, err)
			if assert.NotNil(t, dec) {
				assert.Equal(t, string(test.data), dec.Data)
			}
		})
	}
}
