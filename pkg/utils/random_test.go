package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomInteger(t *testing.T) {
	i := RandomInteger(10)

	require.Less(t, i, 10)
	require.GreaterOrEqual(t, i, 0)
}

func TestRandomString(t *testing.T) {
	s := RandomString(10)

	require.NotEmpty(t, s)
	require.Len(t, s, 10)
}

func TestRandomPassword(t *testing.T) {
	p := RandomPassword()

	require.NotEmpty(t, p)
}

func TestRandomToken(t *testing.T) {
	token := RandomToken(RandomUUID())

	require.NotEmpty(t, token)
}

func TestRandomUUID(t *testing.T) {
	uuid := RandomUUID()

	require.NotEmpty(t, uuid)
}

func TestRandomEmail(t *testing.T) {
	email := RandomEmail()

	require.NotEmpty(t, email)
	require.Contains(t, email, "@")
}
