package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword Hash password with bcrypt library
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword Check if hashed password is equal to given password
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// HashWithSHA256 Hash string with SHA256
func HashWithSHA256(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))

	return hex.EncodeToString(hasher.Sum(nil))
}
