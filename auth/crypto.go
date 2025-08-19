package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength  = 32
	hashLength  = 32
	tokenLength = 32
)

// Argon2 parameters
const (
	argon2_time    = 1
	argon2_memory  = 64 * 1024 // 64 MB
	argon2_threads = 4
)

// generateSalt creates a random salt for password hashing
func generateSalt() (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// hashPassword creates a hash of the password using PBKDF2 with the provided salt
func hashPassword(password, salt string) (string, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("failed to decode salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), saltBytes, argon2_time, argon2_memory, argon2_threads, hashLength)
	return hex.EncodeToString(hash), nil
}

// verifyPassword checks if the provided password matches the stored hash and salt
func verifyPassword(password, storedHash, salt string) (bool, error) {
	computedHash, err := hashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return computedHash == storedHash, nil
}

// generateAccessToken creates a secure random access token
func generateAccessToken() (string, error) {
	token := make([]byte, tokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
