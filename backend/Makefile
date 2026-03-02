.PHONY: build run test lint migrate-up migrate-down migrate-create

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v -race ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
