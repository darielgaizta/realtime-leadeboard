FROM golang:1.21-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY . .

ENTRYPOINT ["goose"]