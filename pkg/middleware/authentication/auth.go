package authentication

import (
	"github.com/gofiber/fiber/v2"
	"incidents_back/pkg/utils"
	"strings"
)

// Config for paseto middleware, contains paseto token maker
type Config struct {
	Filter func(c *fiber.Ctx) bool // Required
	Maker  *utils.Maker
}

// New creates new paseto middleware
func New(config *Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Filter != nil && config.Filter(c) {
			return c.Next()
		}

		requestHeaders := c.GetReqHeaders()
		authHeader, ok := requestHeaders["Authorization"]
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request, provide Auth header",
			})
		}

		token := strings.Split(authHeader, "Bearer ")[1]

		payload, err := config.Maker.VerifyToken(token, true)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		c.Locals("userId", payload.UserId)

		return c.Next()
	}
}
