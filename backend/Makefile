.PHONY: prepare build lint clean clean-db build-all

FRONTEND_DIST=internal/frontend/dist

prepare:
	go mod tidy
	go mod download

build:
	go build -ldflags="-s -w" -o bin/server ./cmd/server/main.go

lint:
	staticcheck ./...

format:
	go fmt ./...

clean:
	rm -rf bin/ $(FRONTEND_DIST)/browser/*

clean-db:
	rm -f null3.db
