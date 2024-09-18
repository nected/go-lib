package base64

import "encoding/base64"

func B64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func B64Decode(data string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decodedData), nil
}
