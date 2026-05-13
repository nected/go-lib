package algo

import (
	"encoding/base64"

	"github.com/sqids/sqids-go"
)

func GetHash(number uint) string {
	s, err := sqids.New(sqids.Options{
		MinLength: 5,
	})
	if err != nil {
		return ""
	}

	id, err := s.Encode([]uint64{uint64(number)})
	if err != nil {
		return ""
	}
	return id
}

func DecodeHash(hash string) *uint64 {
	if hash == "" {
		return nil
	}
	s, err := sqids.New()
	if err != nil {
		return nil
	}
	return &s.Decode(hash)[0]
}

func Encode(planData string) string {
	return base64.StdEncoding.EncodeToString([]byte(planData))
}

func Decode(encodeData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodeData)
	return string(data), err
}
