run:
	go run cmd/main.go

sqlc:
	sqlc generate

migrate_up:
	migrate -path ./db/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

test:
	gotestsum --format pkgname -- -timeout=90s -coverprofile=cover.out ./...

.PHONY: run sqlc