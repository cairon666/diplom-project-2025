package api_key

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// GenerateAPIKey создает случайный безопасный ключ длиной 32 байта и кодирует его в base64.
func GenerateAPIKey() (string, error) {
	key := make([]byte, 32) // 32 байта = 256 бит
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(key), nil // URL safe base64
}

// HashAPIKey возвращает hex-строку sha256 хеша ключа.
func HashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))

	return hex.EncodeToString(hash[:])
}
