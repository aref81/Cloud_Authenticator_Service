package utils

import (
	"encoding/base64"
	"errors"
)

func EncodeBase64(value string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))
	return encoded
}

func DecodeBase64(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("failed to decode Base64")
	}
	value := string(decoded)
	return value, nil
}
