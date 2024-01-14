package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/validator.v9"
	db "incidents_back/db/sqlc"
	"incidents_back/pkg/utils"
	"io"
	"net/http"
	"os"
	"testing"
)

var (
	app   *fiber.App
	repo  *db.Repo
	maker *utils.Maker
)

type HTTPTest struct {
	route              string
	method             string
	requestBody        map[string]interface{}
	requestHeaders     map[string]string
	expectedStatusCode int
	expectedBody       string
	expectedBodyKeys   []string
	expectedMinLength  int
}

func HttpTest(tests *[]HTTPTest, t *testing.T) {
	for _, test := range *tests {
		requestBody, _ := json.Marshal(test.requestBody)
		request, _ := http.NewRequest(test.method, test.route, bytes.NewReader(requestBody))
		request.Header.Add(`Content-Type`, `application/json`)

		for k, v := range test.requestHeaders {
			request.Header.Add(k, v)
		}

		res, err := app.Test(request, -1)

		require.Equal(t, test.expectedStatusCode, res.StatusCode)

		body, err := io.ReadAll(res.Body)

		require.Nil(t, err)

		if test.expectedBody != "" {
			require.Equal(t, test.expectedBody, string(body))
		}

		if test.expectedMinLength != 0 {
			require.GreaterOrEqual(t, len(string(body)), test.expectedMinLength)
		}

		if test.expectedStatusCode != 200 || len(test.expectedBodyKeys) == 0 {
			continue
		}

		var mapBody map[string]interface{}
		_ = json.Unmarshal([]byte(string(body)), &mapBody)

		for _, key := range test.expectedBodyKeys {
			_, ok := mapBody[key]
			require.Equal(t, ok, true)
		}
	}
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("No .env file found")
	}

	conn, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=disable",
			utils.GetEnv("DB_HOST", "localhost"),
			utils.GetEnv("DB_PORT", "5433"),
			utils.GetEnv("DB_NAME", "watcher_test"),
			utils.GetEnv("DB_USER", "watcher"),
			utils.GetEnv("DB_PASSWORD", "secret"),
		),
	)

	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	maker, _ = utils.NewMaker(utils.GetEnv("PASETO_KEY", ""), utils.GetEnv("PASETO_REFRESH_KEY", ""))

	repo = db.NewRepo(conn)

	validate := validator.New()

	app = fiber.New()

	SetupRoutes(&SetupConfig{
		App:       app,
		Repo:      repo,
		Maker:     maker,
		Validator: validate,
	})

	os.Exit(m.Run())
}
