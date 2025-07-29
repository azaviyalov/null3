.PHONY: prepare build lint staticcheck vet clean clean-db build-all

FRONTEND_DIST=internal/core/frontend/fs

prepare:
	go mod tidy
	go mod download

build:
	go build -ldflags="-s -w" -o bin/server ./cmd/server/main.go

lint: staticcheck vet

staticcheck:
	go tool staticcheck ./...

vet:
	go vet ./...

format:
	go fmt ./...

clean:
	rm -rf bin/ $(FRONTEND_DIST)/*

clean-db:
	rm -f null3.db
