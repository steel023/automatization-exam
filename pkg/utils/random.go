package utils

import (
	"github.com/gofrs/uuid"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

// RandomString Generate random string of n letters
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomFloat32 Generate random float from 0 up to n (non-inclusive)
func RandomFloat32(n int) float32 {
	return float32(rand.Intn(n)) + rand.Float32()
}

// RandomFloat64 Generate random float from 0 up to n (non-inclusive)
func RandomFloat64(n int) float64 {
	return float64(rand.Intn(n)) + rand.Float64()
}

// RandomInteger Generate random integer from 0 up to n (non-inclusive)
func RandomInteger(n int) int {
	return rand.Intn(n)
}

// RandomUUID Generate random uuid.UUID (without error handling, for test purposes only)
func RandomUUID() uuid.UUID {
	id, _ := uuid.NewV6()
	return id
}

// RandomEmail Generates random email
func RandomEmail() string {
	return RandomString(12) + "@gmail.com"
}

// RandomPassword Generate random hashed password (without error handling, for test purposes only)
func RandomPassword() string {
	pass, _ := HashPassword(RandomString(8))
	return pass
}

// RandomToken Generate random paseto token (without error handling, for test purposes only)
func RandomToken(userId uuid.UUID) string {
	maker, _ := NewMaker(RandomString(32), RandomString(32))
	token, _ := maker.GenerateToken(userId, time.Hour, false)
	return token
}
