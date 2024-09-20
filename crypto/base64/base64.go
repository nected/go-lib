package base64

import (
	"encoding/base64"

	"github.com/nected/go-lib/crypto/errors"
)

func B64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func B64Decode(data string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.ErrInvalidData
	}
	return string(decodedData), nil
}
