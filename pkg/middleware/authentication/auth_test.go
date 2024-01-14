package authentication

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	config := Config{
		Filter: nil,
		Maker:  maker,
	}

	middleware := New(&config)

	require.NotEmpty(t, middleware)
}
