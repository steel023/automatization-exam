package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gopkg.in/go-playground/validator.v9"
	db "incidents_back/db/sqlc"
	"incidents_back/pkg/middleware/authentication"
	"incidents_back/pkg/utils"
)

type SetupConfig struct {
	App       *fiber.App
	Repo      *db.Repo
	Maker     *utils.Maker
	Validator *validator.Validate
}

// Handlers Struct to store utilities shared across all handlers
type Handlers struct {
	Repo      *db.Repo
	Maker     *utils.Maker
	Validator *validator.Validate
}

// NewHandlers Make new Handlers struct
func NewHandlers(repo *db.Repo, maker *utils.Maker, validator *validator.Validate) *Handlers {
	return &Handlers{Repo: repo, Maker: maker, Validator: validator}
}

func SetupRoutes(config *SetupConfig) {
	config.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	config.App.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetEnv("ALLOWED_ORIGINS", "*"),
		AllowMethods:     "POST,GET,PUT,DELETE,HEAD,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	config.App.Use(limiter.New(limiter.Config{
		Max: 240,
	}))

	config.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	authConfig := &authentication.Config{
		Filter: nil,
		Maker:  config.Maker,
	}

	_ = authentication.New(authConfig)

	handlers := NewHandlers(config.Repo, config.Maker, config.Validator)

	api := config.App.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("login", handlers.login)
	auth.Post("register", handlers.register)
	auth.Post("refresh", handlers.refresh)
}
