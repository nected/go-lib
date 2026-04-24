package algo

import (
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
