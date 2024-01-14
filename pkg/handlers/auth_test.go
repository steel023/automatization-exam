package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	db "incidents_back/db/sqlc"
	"incidents_back/pkg/services"
	"incidents_back/pkg/utils"
	"net/http"
	"testing"
)

func createRandomUser() *db.User {
	pass, _ := utils.HashPassword("test123")
	params := db.CreateUserParams{
		Email:    utils.RandomEmail(),
		Password: pass,
	}

	user, err := repo.CreateUser(context.Background(), params)
	if err != nil {
		log.WithError(err).Fatal("Unable to create a user")
	}

	return &user
}

func TestRegisterRoute(t *testing.T) {
	email := utils.RandomEmail()
	tests := []HTTPTest{
		{
			route:  "/api/v1/auth/register",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"email":    email,
				"password": "test123",
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBodyKeys:   []string{"access_token", "refresh_token", "expires_at", "user"},
		},
		{
			route:  "/api/v1/auth/register",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"email":    email,
				"password": "test123",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
		},
		{
			route:  "/api/v1/auth/register",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"emaiwl":    email,
				"passwword": "test123",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
		},
		{
			route:  "/api/v1/auth/register",
			method: http.MethodPost,
			requestHeaders: map[string]string{
				"Content-Type": "image/png",
			},
			requestBody: map[string]interface{}{
				"emaiwl":    email,
				"passwword": "test123",
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
	}

	HttpTest(&tests, t)
}

func TestLoginRoute(t *testing.T) {
	user := createRandomUser()
	tests := []HTTPTest{
		{
			route:  "/api/v1/auth/login",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"email":    user.Email,
				"password": "test123",
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBodyKeys:   []string{"access_token", "refresh_token", "expires_at", "user"},
		},
		{
			route:  "/api/v1/auth/login",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"email":    user.Email,
				"password": "test1232",
			},
			expectedStatusCode: fiber.StatusUnauthorized,
			expectedBodyKeys:   []string{"message"},
		},
		{
			route:  "/api/v1/auth/login",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"wefwe":  "test1232",
				"sdfsdf": "sdfsd",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBodyKeys:   []string{"message"},
		},
		{
			route:  "/api/v1/auth/login",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"email":    "",
				"password": "",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBodyKeys:   []string{"message"},
		},
		{
			route:  "/api/v1/auth/login",
			method: http.MethodPost,
			requestHeaders: map[string]string{
				"Content-Type": "image/png",
			},
			requestBody: map[string]interface{}{
				"email":    "",
				"password": "",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedBodyKeys:   []string{"message"},
		},
	}

	HttpTest(&tests, t)
}

func TestRefreshRoute(t *testing.T) {
	user := createRandomUser()

	authResponse, _ := services.Login(user.Email, "test123", repo, maker)

	tests := []HTTPTest{
		{
			route:  "/api/v1/auth/refresh",
			method: http.MethodPost,
			requestHeaders: map[string]string{
				"Authorization": "Bearer " + authResponse.RefreshToken,
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBodyKeys:   []string{"access_token", "refresh_token", "expires_at", "user"},
		},
		{
			route:  "/api/v1/auth/refresh",
			method: http.MethodPost,
			requestHeaders: map[string]string{
				"Authorization": "Bearer " + authResponse.RefreshToken,
			},
			expectedStatusCode: fiber.StatusUnauthorized,
		},
		{
			route:  "/api/v1/auth/refresh",
			method: http.MethodPost,
			requestHeaders: map[string]string{
				"Authorization": "Bearer sdfsdfsdfsdfsrortkortknv",
			},
			expectedStatusCode: fiber.StatusUnauthorized,
		},
		{
			route:              "/api/v1/auth/refresh",
			method:             http.MethodPost,
			expectedStatusCode: fiber.StatusBadRequest,
		},
	}

	HttpTest(&tests, t)
}
