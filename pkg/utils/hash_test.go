package utils

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hashedPassword, err := HashPassword(RandomString(8))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestCheckPassword(t *testing.T) {
	password := RandomString(8)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	password = password[:4] + "s2332ew" + password[5:]

	err = CheckPassword(password, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	err = CheckPassword("test123", "$2a$10$jUOo/SnKN.kg2NgNFmZ7O.m2DPWmU9NczejYe3cfDL79ijvroum3q")
	require.NoError(t, err)
}
