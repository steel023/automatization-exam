package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createMaker(t *testing.T) *Maker {
	maker, err := NewMaker(GetEnv("PASETO_KEY", ""), GetEnv("PASETO_REFRESH_KEY", ""))

	require.NoError(t, err)
	require.NotEmpty(t, maker)

	errMaker, err := NewMaker("sdfsdf", "sdfsdfsdf")

	require.Empty(t, errMaker)
	require.Error(t, err)

	return maker
}

func createToken(t *testing.T) (*Maker, string) {
	maker := createMaker(t)

	token, err := maker.GenerateToken(RandomUUID(), time.Hour, false)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	return maker, token
}

func TestNewPayload(t *testing.T) {
	payload, err := NewPayload(RandomUUID(), time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
}

func TestPayload_Valid(t *testing.T) {
	payload, err := NewPayload(RandomUUID(), time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	err = payload.Valid()
	require.NoError(t, err)

	expPayload, err := NewPayload(RandomUUID(), -time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, expPayload)

	err = expPayload.Valid()
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestNewMaker(t *testing.T) {
	createMaker(t)
}

func TestMaker_GenerateToken(t *testing.T) {
	createToken(t)
}

func TestMaker_VerifyToken(t *testing.T) {
	maker, token := createToken(t)

	payload, err := maker.VerifyToken(token, false)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	token = token[:10] + "e32423" + token[11:]

	payload, err = maker.VerifyToken(token, false)

	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())

	expToken, err := maker.GenerateToken(RandomUUID(), -time.Hour, false)

	require.NoError(t, err)
	require.NotEmpty(t, expToken)

	payload, err = maker.VerifyToken(expToken, false)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
