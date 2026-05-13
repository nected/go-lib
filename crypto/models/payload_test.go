package models

import (
	"fmt"
	"testing"

	"github.com/nected/go-lib/crypto/base64"
	"github.com/stretchr/testify/assert"
)

// Payload.String has three branches:
//   - empty KeyName: returns EncryptedData as-is.
//   - KeyType AES: returns EncryptedData as-is (key name is ignored).
//   - RSA (default): formats "$<name>$<version>$<encrypted>" and base64-encodes it.
func TestPayloadString(t *testing.T) {
	tests := []struct {
		name string
		p    Payload
		want string
	}{
		{
			name: "empty key name returns encrypted data verbatim",
			p:    Payload{EncryptedData: "abc"},
			want: "abc",
		},
		{
			name: "AES key type returns encrypted data verbatim",
			p:    Payload{KeyName: "APP", KeyVersion: 1, KeyType: KeyTypeAES, EncryptedData: "abc"},
			want: "abc",
		},
		{
			name: "RSA key type wraps and base64-encodes",
			p:    Payload{KeyName: "APP", KeyVersion: 2, KeyType: KeyTypeRSA, EncryptedData: "ciphertext"},
			want: base64.B64Encode([]byte("$APP$2$ciphertext")),
		},
		{
			// KeyType empty (zero value) is not KeyTypeAES, so the RSA branch runs.
			name: "default key type falls through to RSA formatting",
			p:    Payload{KeyName: "APP", KeyVersion: 1, EncryptedData: "x"},
			want: base64.B64Encode([]byte("$APP$1$x")),
		},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			assert.Equal(t, test.want, test.p.String())
		})
	}
}
