DB_URL=postgresql://postgres:pr_service_password@localhost:5432/pr_service?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

run:
	go run ./cmd

build:
	go build -o bin/pr-service ./cmd

tidy:
	go mod tidy

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

.PHONY: migrate-up migrate-down run build tidy docker-up docker-down