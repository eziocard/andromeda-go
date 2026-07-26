package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

const RefreshTokenTTL = time.Hour * 24 * 7 // 7 días

// GenerateRefreshToken crea un token aleatorio seguro (se entrega al cliente tal cual).
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// HashToken se usa para guardar en BD solo el hash, nunca el token en texto plano.
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
