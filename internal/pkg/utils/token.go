package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken generates a secure random token of given byte length
func GenerateRandomToken(nBytes int) string {
	bytes := make([]byte, nBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}