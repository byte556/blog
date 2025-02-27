package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// GenerateToken генерирует криптографически стойкий токен.
func GenerateToken() string {
	b := make([]byte, 32) // 32 байта = 256 бит
	if _, err := rand.Read(b); err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return ""
	}
	return hex.EncodeToString(b)
}

// ValidateIdempotencyToken сравнивает значение токена из сессии и из запроса.
func ValidateIdempotencyToken(sessionValue interface{}, requestToken string) bool {
	token, ok := sessionValue.(string)
	if !ok || token == "" || requestToken == "" || requestToken != token {
		return false
	}
	return true
}
