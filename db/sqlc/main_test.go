package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"incidents_back/pkg/utils"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

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

	testQueries = New(conn)

	os.Exit(m.Run())
}
