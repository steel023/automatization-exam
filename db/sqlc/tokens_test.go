package db

import (
	"context"
	"database/sql"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"incidents_back/pkg/utils"
	"testing"
	"time"
)

func createRandomToken(t *testing.T) Token {
	user := createRandomUser(t)

	expiresAt := null.Time{
		sql.NullTime{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

	params := CreateTokenParams{
		Token:     utils.RandomToken(user.ID),
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	token, err := testQueries.CreateToken(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.Equal(t, params.Token, token.Token)
	require.Equal(t, params.UserID, token.UserID)
	require.WithinDuration(t, params.ExpiresAt.Time, token.ExpiresAt.Time, time.Second)

	require.NotZero(t, token.ID)

	return token
}

func TestQueries_CreateToken(t *testing.T) {
	createRandomToken(t)
}

func TestQueries_GetToken(t *testing.T) {
	randomToken := createRandomToken(t)
	token, err := testQueries.GetToken(context.Background(), randomToken.Token)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, randomToken.Token, token.Token)
	require.WithinDuration(t, randomToken.ExpiresAt.Time, token.ExpiresAt.Time, time.Second)
}

func TestQueries_DeleteToken(t *testing.T) {
	randomToken := createRandomToken(t)

	err := testQueries.DeleteToken(context.Background(), randomToken.ID)
	require.NoError(t, err)

	token, err := testQueries.GetToken(context.Background(), randomToken.Token)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, token)
}

func TestQueries_DeleteUsersTokens(t *testing.T) {
	token := createRandomToken(t)

	err := testQueries.DeleteUsersTokens(context.Background(), token.UserID)
	require.NoError(t, err)

	token, err = testQueries.GetToken(context.Background(), token.Token)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, token)
}
