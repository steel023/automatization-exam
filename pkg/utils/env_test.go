package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := GetEnv("PASETO_KEY", "aboba")

	require.NotEqual(t, key, "aboba")
	require.NotEmpty(t, key)
}
