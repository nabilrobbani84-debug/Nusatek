.PHONY: run build test clean docker-up docker-down

APP_NAME=nusatek-backend

run:
	go run cmd/api/main.go

build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
