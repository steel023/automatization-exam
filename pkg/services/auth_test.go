package services

import (
	"context"
	"database/sql"
	"github.com/guregu/null"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	db "incidents_back/db/sqlc"
	"incidents_back/pkg/utils"
	"testing"
	"time"
)

func createRandomUser() *db.User {
	pass, _ := utils.HashPassword("test123")
	params := db.CreateUserParams{
		Email:    utils.RandomString(8),
		Password: pass,
	}

	user, err := repo.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal("Unable to create a user")
	}

	return &user
}

func TestLogin(t *testing.T) {
	user := createRandomUser()

	require.NotEmpty(t, user)

	authResponse, err := Login(user.Email, "test123", repo, maker)

	require.NoError(t, err)
	require.NotEmpty(t, authResponse)
	require.NotEmpty(t, authResponse.ExpiresAt)
	require.NotEmpty(t, authResponse.AccessToken)
	require.NotEmpty(t, authResponse.RefreshToken)
	require.NotEmpty(t, authResponse.User)
	require.NotEqual(t, authResponse.AccessToken, authResponse.RefreshToken)

	authResponse, err = Login(user.Email, "test1232", repo, maker)

	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	require.Nil(t, authResponse)

	authResponse, err = Login(user.Email[:2], "test1232", repo, maker)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Nil(t, authResponse)

	authResponse, err = Login("", "test123", repo, maker)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestRegister(t *testing.T) {
	email := utils.RandomString(5)

	authResponse, err := Register(email, "test123", repo, maker)

	require.NoError(t, err)
	require.NotEmpty(t, authResponse)
	require.NotEmpty(t, authResponse.ExpiresAt)
	require.NotEmpty(t, authResponse.AccessToken)
	require.NotEmpty(t, authResponse.RefreshToken)
	require.NotEmpty(t, authResponse.User)
	require.NotEqual(t, authResponse.AccessToken, authResponse.RefreshToken)

	authResponse, err = Register(email, "test123", repo, maker)
	require.Error(t, err)

	if err, ok := err.(*pq.Error); ok {
		require.Equal(t, err.Code.Name(), "unique_violation")
	}

	authResponse, err = Register("", "", repo, maker)

	require.Error(t, err)
	if err, ok := err.(*pq.Error); ok {
		require.Equal(t, err.Code.Name(), "check_violation")
	}
}

func TestRefresh(t *testing.T) {
	user := createRandomUser()
	anotherUser := createRandomUser()

	authResponse, err := Login(user.Email, "test123", repo, maker)

	require.NoError(t, err)
	require.NotEmpty(t, authResponse)

	authResponse, err = Refresh(authResponse.RefreshToken, user.ID, repo, maker)

	require.NoError(t, err)
	require.NotEmpty(t, authResponse)
	require.NotEmpty(t, authResponse.ExpiresAt)
	require.NotEmpty(t, authResponse.AccessToken)
	require.NotEmpty(t, authResponse.RefreshToken)
	require.NotEmpty(t, authResponse.User)
	require.NotEqual(t, authResponse.AccessToken, authResponse.RefreshToken)

	authResponse, err = Refresh(authResponse.RefreshToken[1:], user.ID, repo, maker)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, authResponse)

	anotherAuthResponse, err := Login(anotherUser.Email, "test123", repo, maker)
	anotherAuthResponse, err = Refresh(anotherAuthResponse.RefreshToken, user.ID, repo, maker)

	require.Error(t, err)
	require.EqualError(t, err, utils.ErrInvalidToken.Error())
	require.Empty(t, anotherAuthResponse)

	expiredToken, _ := maker.GenerateToken(user.ID, -time.Hour, false)
	createTokenParams := db.CreateTokenParams{
		Token:  expiredToken,
		UserID: user.ID,
		ExpiresAt: null.Time{
			NullTime: sql.NullTime{
				Time:  time.Now().Add(-time.Hour),
				Valid: true,
			},
		},
	}
	token, _ := repo.CreateToken(context.Background(), createTokenParams)

	authResponse, err = Refresh(token.Token, user.ID, repo, maker)

	require.Error(t, err)
	require.EqualError(t, err, utils.ErrExpiredToken.Error())
	require.Empty(t, authResponse)
}
