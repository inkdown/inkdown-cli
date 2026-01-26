package utils

import "encoding/base64"

func DecodeBase64(content string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func EncodeBase64(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}
