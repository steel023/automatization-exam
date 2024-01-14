package handlers

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"incidents_back/pkg/services"
	"strings"
)

type LoginParams struct {
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required"`
}

type RegisterParams struct {
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required"`
}

func (h *Handlers) login(c *fiber.Ctx) error {
	params := new(LoginParams)

	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request parameters",
		})
	}

	err := h.Validator.Struct(params)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validation failed",
		})
	}

	response, err := services.Login(params.Email, params.Password, h.Repo, h.Maker)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}
		log.WithError(err).Error("login error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handlers) register(c *fiber.Ctx) error {
	params := new(RegisterParams)

	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request parameters",
		})
	}

	err := h.Validator.Struct(params)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validation failed",
		})
	}

	response, err := services.Register(params.Email, params.Password, h.Repo, h.Maker)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validation failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handlers) refresh(c *fiber.Ctx) error {
	requestHeaders := c.GetReqHeaders()
	authHeader, ok := requestHeaders["Authorization"]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request, provide Auth header",
		})
	}

	token := strings.Split(authHeader, "Bearer ")[1]
	payload, err := h.Maker.VerifyToken(token, false)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	response, err := services.Refresh(token, payload.UserId, h.Repo, h.Maker)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
