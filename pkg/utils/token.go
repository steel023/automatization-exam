package utils

import (
	"errors"
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/gofrs/uuid"
	"github.com/o1egl/paseto"
	"time"
)

var (
	// ErrInvalidToken Invalid token error
	ErrInvalidToken = errors.New("token is invalid")
	// ErrExpiredToken Expired token error
	ErrExpiredToken = errors.New("token has expired")
)

// Maker Paseto token maker
type Maker struct {
	paseto              *paseto.V2
	symmetricKey        []byte
	refreshSymmetricKey []byte
}

// Payload Paseto token payload
type Payload struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewPayload Create new Payload struct
func NewPayload(userId uuid.UUID, duration time.Duration) (*Payload, error) {
	tokenUUID, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}

	payload := Payload{
		Id:        tokenUUID,
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return &payload, nil
}

// Valid Check if token is expired
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

// NewMaker Create new Maker struct
func NewMaker(symmetricKey, refreshSymmetricKey string) (*Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize || len(refreshSymmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size")
	}

	maker := Maker{
		paseto:              paseto.NewV2(),
		symmetricKey:        []byte(symmetricKey),
		refreshSymmetricKey: []byte(refreshSymmetricKey),
	}

	return &maker, nil
}

// GenerateToken Generate new Paseto token
func (m *Maker) GenerateToken(userId uuid.UUID, duration time.Duration, access bool) (string, error) {
	payload, err := NewPayload(userId, duration)
	if err != nil {
		return "", err
	}

	if access {
		return m.paseto.Encrypt(m.symmetricKey, payload, nil)
	}

	return m.paseto.Encrypt(m.refreshSymmetricKey, payload, nil)
}

// VerifyToken Verify a Paseto token
func (m *Maker) VerifyToken(token string, access bool) (*Payload, error) {
	payload := &Payload{}

	var err error
	if access {
		err = m.paseto.Decrypt(token, m.symmetricKey, payload, nil)
	} else {
		err = m.paseto.Decrypt(token, m.refreshSymmetricKey, payload, nil)
	}

	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
