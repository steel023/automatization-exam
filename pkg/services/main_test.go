package services

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	db "incidents_back/db/sqlc"
	"incidents_back/pkg/utils"
	"os"
	"testing"
)

var (
	repo  *db.Repo
	maker *utils.Maker
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("No .env file found")
	}

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
		log.Fatal("Cannot connect to DB", err)
	}

	repo = db.NewRepo(conn)

	maker, err = utils.NewMaker(utils.GetEnv("PASETO_KEY", ""), utils.GetEnv("PASETO_REFRESH_KEY", ""))
	if err != nil {
		log.Fatal("Cannot create maker", err)
	}

	os.Exit(m.Run())
}
