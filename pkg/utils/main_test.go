package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("No .env file found")
	}

	os.Exit(m.Run())
}
