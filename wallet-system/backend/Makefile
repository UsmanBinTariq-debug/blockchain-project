.PHONY: build run test clean dev

build:
	go build -o bin/wallet-server ./cmd/server

run: build
	./bin/wallet-server

dev:
	go run ./cmd/server/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	go clean

deps:
	go mod download
	go mod tidy

docker-build:
	docker build -t crypto-wallet-backend:latest .

docker-run:
	docker run -p 8080:8080 crypto-wallet-backend:latest
