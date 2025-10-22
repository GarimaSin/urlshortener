.PHONY: build run docker-compose

build:
	go build -o bin/shortener ./cmd/server

run:
	go run ./cmd/server

docker-compose:
	docker compose -f build/docker-compose.yml up --build
