package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	SaltLength  = 32
	HashLength  = 32
	TokenLength = 32
)

// Argon2 parameters
const (
	Time    = 1
	Memory  = 64 * 1024 // 64 MB
	Threads = 4
)

// GenerateSalt creates a random salt for password hashing
func GenerateSalt() (string, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword creates a hash of the password using PBKDF2 with the provided salt
func HashPassword(password, salt string) (string, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("failed to decode salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), saltBytes, Time, Memory, Threads, HashLength)
	return hex.EncodeToString(hash), nil
}

// VerifyPassword checks if the provided password matches the stored hash and salt
func VerifyPassword(password, storedHash, salt string) (bool, error) {
	computedHash, err := HashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return computedHash == storedHash, nil
}

// GenerateAccessToken creates a secure random access token
func GenerateAccessToken() (string, error) {
	token := make([]byte, TokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
