package main

import (
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	db "incidents_back/db/sqlc"
	"incidents_back/pkg/handlers"
	"incidents_back/pkg/utils"
	"time"
)

func main() {
	conn, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=disable",
			utils.GetEnv("DB_HOST", "localhost"),
			utils.GetEnv("DB_PORT", "5432"),
			utils.GetEnv("DB_NAME", "incident"),
			utils.GetEnv("DB_USER", "incident"),
			utils.GetEnv("DB_PASSWORD", "incident"),
		),
	)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to DB")
	}

	repo := db.NewRepo(conn)

	app := fiber.New(fiber.Config{
		AppName:     "Incidents",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		IdleTimeout: time.Minute,
	})

	tokenMakerr, err := utils.NewMaker(utils.GetEnv("PASETO_KEY", ""), utils.GetEnv("PASETO_REFRESH_KEY", ""))
	if err != nil {
		log.Fatal("Invalid key size")
	}

	validate := validator.New()

	handlers.SetupRoutes(&handlers.SetupConfig{
		App:       app,
		Repo:      repo,
		Validator: validate,
		Maker:     tokenMakerr,
	})

	if err := app.Listen(":8888"); err != nil {
		log.WithError(err).Fatal("server failed to start!")
	}
}
