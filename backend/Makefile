.DEFAULT_GOAL := build

.PHONY: build test-unit test coverage coverage-html lint staticcheck vet format clean clean-db

COVERAGE_PACKAGES=$(shell go list ./... | grep -v /internal/testutil)

build:
	mkdir -p bin
	go build -ldflags="-s -w" -o bin/server ./cmd/server

test-unit:
	go test -short ./...

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out $(COVERAGE_PACKAGES)
	go tool cover -func=coverage.out | tail -n 1

coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html

lint: staticcheck vet

staticcheck:
	go tool staticcheck ./...

vet:
	go vet ./...

format:
	go fmt ./...

clean:
	rm -rf bin/ internal/core/frontend/fs/*
	rm -f coverage.out coverage.html

clean-db:
	rm -f null3.db
