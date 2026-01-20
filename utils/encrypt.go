package utils

import (
	"encoding/base64"
	"errors"
	"os"
)

func GetDecodedDBURL() (string, error) {
	encoded := os.Getenv("DATABASE_URL")
	if encoded == "" {
		return "", errors.New("DATABASE_URL is empty")
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
