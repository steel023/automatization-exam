package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	db "incidents_back/db/sqlc"
)

type Config struct {
	Filter func(c *fiber.Ctx) bool // Required
	Repo   *db.Repo
	Role   int32
}

// New creates new admin middleware
func New(config *Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Filter != nil && config.Filter(c) {
			return c.Next()
		}

		userID := c.Locals("userId").(uuid.UUID)

		user, err := config.Repo.GetUserById(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User not found",
			})
		}

		if user.Role < config.Role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "You don't have such rights",
			})
		}

		return c.Next()
	}
}
