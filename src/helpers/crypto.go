package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func GenRandString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	result := hex.EncodeToString(bytes)
	return result[0:length], nil
}

func GetHash(seed string) (string, error) {
	length := len(seed) / 2
	if length <= 0 {
		return "", errors.New("Invalid seed")
	}
	hasher := sha256.New()
	hasher.Write([]byte(seed))
	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum), nil
}
