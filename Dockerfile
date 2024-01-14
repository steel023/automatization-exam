FROM golang:1.18-alpine

LABEL maintainer="Pasha Kubalov kubalov.pasha@gmail.com"

WORKDIR /app

RUN apk add --no-cache git
RUN apk add --no-cache make build-base

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install gotest.tools/gotestsum@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -buildvcs=false -o main ./cmd

EXPOSE 8888

CMD ["./main"]
