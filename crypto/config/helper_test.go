package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// pemEncodeRSA serializes a freshly generated RSA private key as PKCS8 PEM
// using the package's canonical block type. Returned bytes are what
// loadPrivateKey expects on the happy path.
func pemEncodeRSA(t *testing.T) (*rsa.PrivateKey, []byte) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatalf("MarshalPKCS8PrivateKey: %v", err)
	}
	block := &pem.Block{Type: PRIV_KEY_TYPE, Bytes: der}
	return priv, pem.EncodeToMemory(block)
}

// loadPrivateKey covers four branches:
//   - undecodable PEM → "Invalid private key"
//   - wrong block type → "Invalid private key type"
//   - decodable PEM but ParsePKCS8 fails (PKCS1-encoded body in a PKCS8 block)
//   - parsed key isn't *rsa.PrivateKey (EC key wrapped in PKCS8)
//   - happy path → key matches the source
func TestLoadPrivateKey(t *testing.T) {
	priv, validPEM := pemEncodeRSA(t)

	// Build an EC PKCS8 block to hit the "not *rsa.PrivateKey" branch.
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	ecDER, err := x509.MarshalPKCS8PrivateKey(ecKey)
	assert.NoError(t, err)
	ecPEM := pem.EncodeToMemory(&pem.Block{Type: PRIV_KEY_TYPE, Bytes: ecDER})

	// PKCS1-bodied PEM in a PKCS8-typed block — pem.Decode succeeds, type
	// matches, but ParsePKCS8PrivateKey rejects the PKCS1 body.
	pkcs1DER := x509.MarshalPKCS1PrivateKey(priv)
	pkcs1WrappedAsPKCS8 := pem.EncodeToMemory(&pem.Block{Type: PRIV_KEY_TYPE, Bytes: pkcs1DER})

	tests := []struct {
		name    string
		pemData []byte
		wantErr bool
	}{
		{name: "undecodable pem", pemData: []byte("not pem at all"), wantErr: true},
		{
			// Block decodes but Type is wrong ("CERTIFICATE" != "PRIVATE KEY").
			name:    "wrong block type",
			pemData: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("xx")}),
			wantErr: true,
		},
		{name: "pkcs1 body in pkcs8 block", pemData: pkcs1WrappedAsPKCS8, wantErr: true},
		{name: "ec key not rsa", pemData: ecPEM, wantErr: true},
		{name: "happy path", pemData: validPEM, wantErr: false},
	}

	for id, test := range tests {
		t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
			got, err := loadPrivateKey(test.pemData)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, priv.N.Cmp(got.N), 0)
		})
	}
}

// loadPrivateKeyFromFile delegates to loadPrivateKey after reading the file.
// Covers both branches: file-missing and successful file → key parse.
func TestLoadPrivateKeyFromFile(t *testing.T) {
	priv, validPEM := pemEncodeRSA(t)

	dir := t.TempDir()
	good := filepath.Join(dir, "key.pem")
	if err := os.WriteFile(good, validPEM, 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	t.Run("missing file returns error", func(t *testing.T) {
		_, err := loadPrivateKeyFromFile(filepath.Join(dir, "does-not-exist.pem"))
		assert.Error(t, err)
	})

	t.Run("valid pem file returns key", func(t *testing.T) {
		got, err := loadPrivateKeyFromFile(good)
		assert.NoError(t, err)
		if assert.NotNil(t, got) {
			assert.Equal(t, priv.N.Cmp(got.N), 0)
		}
	})
}

// generatePublicKey precomputes and returns &priv.PublicKey. The pointer must
// be non-nil and the modulus must match the source private key.
func TestGeneratePublicKey(t *testing.T) {
	priv, _ := pemEncodeRSA(t)
	pub := generatePublicKey(priv)
	if assert.NotNil(t, pub) {
		assert.Equal(t, priv.PublicKey.N.Cmp(pub.N), 0)
	}
}
