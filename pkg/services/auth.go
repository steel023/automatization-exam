package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	db "incidents_back/db/sqlc"
	"incidents_back/pkg/constants"
	"incidents_back/pkg/utils"
	"time"
)

// AuthResponse contains AccessToken, RefreshToken and ExpiresAt time for access token
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *db.User  `json:"user"`
}

// Login logs user in
func Login(email, password string, repo *db.Repo, maker *utils.Maker) (*AuthResponse, error) {
	user, err := repo.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	if pErr := utils.CheckPassword(password, user.Password); pErr != nil {
		return nil, pErr
	}

	accessToken, err := maker.GenerateToken(user.ID, constants.AccessTokenLifetime, true)
	if err != nil {
		return nil, err
	}

	refreshToken, err := maker.GenerateToken(user.ID, constants.RefreshTokenLifetime, false)
	if err != nil {
		return nil, err
	}

	expiresAt := null.Time{
		sql.NullTime{
			Time:  time.Now().Add(constants.RefreshTokenLifetime),
			Valid: true,
		},
	}

	params := db.CreateTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	_, err = repo.CreateToken(context.Background(), params)
	if err != nil {
		return nil, err
	}

	response := &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(constants.AccessTokenLifetime),
		User:         &user,
	}

	return response, nil
}

// Register creates new user, returns AuthResponse
func Register(email, password string, repo *db.Repo, maker *utils.Maker) (*AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.WithError(err).Error("Unable to hash a password")
		return nil, err
	}

	newUserParams := db.CreateUserParams{Email: email, Password: hashedPassword}
	user, err := repo.CreateUser(context.Background(), newUserParams)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "foreign_key_violation" {
				newUserParams = db.CreateUserParams{Email: email, Password: hashedPassword}
				user, err = repo.CreateUser(context.Background(), newUserParams)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	accessToken, err := maker.GenerateToken(user.ID, constants.AccessTokenLifetime, true)
	if err != nil {
		return nil, err
	}

	refreshToken, err := maker.GenerateToken(user.ID, constants.RefreshTokenLifetime, false)
	if err != nil {
		return nil, err
	}

	expiresAt := null.Time{
		sql.NullTime{
			Time:  time.Now().Add(constants.RefreshTokenLifetime),
			Valid: true,
		},
	}

	params := db.CreateTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	_, err = repo.CreateToken(context.Background(), params)
	if err != nil {
		return nil, err
	}

	response := &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(constants.AccessTokenLifetime),
		User:         &user,
	}

	return response, nil
}

// Refresh refreshes token, returns AuthResponse
func Refresh(oldRefreshToken string, userId uuid.UUID, repo *db.Repo, maker *utils.Maker) (*AuthResponse, error) {
	oldToken, err := repo.GetToken(context.Background(), oldRefreshToken)
	if err != nil {
		return nil, err
	}

	if oldToken.UserID != userId {
		return nil, utils.ErrInvalidToken
	}

	if !oldToken.ExpiresAt.Valid {
		return nil, utils.ErrInvalidToken
	}

	if time.Now().After(oldToken.ExpiresAt.Time) {
		return nil, utils.ErrExpiredToken
	}

	if err = repo.DeleteToken(context.Background(), oldToken.ID); err != nil {
		return nil, err
	}

	accessToken, err := maker.GenerateToken(userId, constants.AccessTokenLifetime, true)
	if err != nil {
		return nil, err
	}

	refreshToken, err := maker.GenerateToken(userId, constants.RefreshTokenLifetime, false)
	if err != nil {
		return nil, err
	}

	expiresAt := null.Time{
		sql.NullTime{
			Time:  time.Now().Add(constants.RefreshTokenLifetime),
			Valid: true,
		},
	}

	params := db.CreateTokenParams{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	_, err = repo.CreateToken(context.Background(), params)
	if err != nil {
		return nil, err
	}

	user, err := repo.GetUserById(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	response := &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(constants.AccessTokenLifetime),
		User:         &user,
	}

	return response, nil
}
