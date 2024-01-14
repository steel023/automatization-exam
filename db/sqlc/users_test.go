package db

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"incidents_back/pkg/utils"
	"testing"
)

func createRandomUser(t *testing.T) *User {
	params := CreateUserParams{
		Email:    utils.RandomString(8),
		Password: utils.RandomPassword(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	if err != nil {
		log.WithError(err).Fatal("Unable to create a user")
	}

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Role, int32(0))

	return &user
}

func TestQueries_GetUserByEmail(t *testing.T) {
	dbUser := createRandomUser(t)

	user, err := testQueries.GetUserByEmail(context.Background(), dbUser.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, dbUser.Email, user.Email)
	require.Equal(t, dbUser.ID, user.ID)

	user, err = testQueries.GetUserByEmail(context.Background(), "somenonexistingemail")

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)

	user, err = testQueries.GetUserByEmail(context.Background(), "")

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestQueries_CreateUser(t *testing.T) {
	password, _ := utils.HashPassword("test123")
	params := CreateUserParams{
		Email:    utils.RandomString(25),
		Password: password,
	}

	user, err := testQueries.CreateUser(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Role, int32(0))

	user, err = testQueries.CreateUser(context.Background(), params)

	require.Error(t, err)
	if err, ok := err.(*pq.Error); ok {
		require.Equal(t, err.Code.Name(), "unique_violation")
	}

	params.Email = ""

	user, err = testQueries.CreateUser(context.Background(), params)

	require.Error(t, err)
	if err, ok := err.(*pq.Error); ok {
		require.Equal(t, err.Code.Name(), "check_violation")
	}
}

func TestQueries_GetUserById(t *testing.T) {
	createdUser := createRandomUser(t)

	user, err := testQueries.GetUserById(context.Background(), createdUser.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Email, createdUser.Email)
}
