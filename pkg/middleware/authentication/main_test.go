package authentication

import (
	"github.com/joho/godotenv"
	"incidents_back/pkg/utils"
	"log"
	"os"
	"testing"
)

var maker *utils.Maker

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("No .env file found")
	}

	maker, _ = utils.NewMaker(utils.GetEnv("PASETO_KEY", ""), utils.GetEnv("PASETO_REFRESH_KEY", ""))

	os.Exit(m.Run())
}
