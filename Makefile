.DEFAULT_GOAL := build

.PHONY: run build lint format coverage release

SHELL := /bin/bash
FRONTEND_DEPS := frontend/node_modules/.package-lock.json

$(FRONTEND_DEPS): frontend/package.json frontend/package-lock.json
	npm --prefix frontend ci

run: $(FRONTEND_DEPS)
	@set -m; \
	(cd backend && exec go run ./cmd/server) & backend_pid=$$!; \
	(cd frontend && exec npm start) & frontend_pid=$$!; \
	cleanup() { \
		trap - EXIT INT TERM; \
		kill -- "-$$backend_pid" "-$$frontend_pid" 2>/dev/null || true; \
		wait "$$backend_pid" "$$frontend_pid" 2>/dev/null || true; \
	}; \
	trap cleanup EXIT INT TERM; \
	wait -n "$$backend_pid" "$$frontend_pid"

build: $(FRONTEND_DEPS)
	mkdir -p backend/bin
	go -C backend build -ldflags="-s -w" -o bin/server ./cmd/server
	npm --prefix frontend run build

lint: $(FRONTEND_DEPS)
	cd backend && go tool staticcheck ./...
	go -C backend vet ./...
	npm --prefix frontend run lint

format: $(FRONTEND_DEPS)
	go -C backend fmt ./...
	npm --prefix frontend run prettier

coverage:
	cd backend && go test -coverprofile=coverage.out $$(go list ./... | grep -v /internal/testutil)
	cd backend && go tool cover -func=coverage.out | tail -n 1

release: $(FRONTEND_DEPS)
	rm -rf backend/internal/core/frontend/fs/*
	npm --prefix frontend run build
	cp -r frontend/dist/frontend/browser/* backend/internal/core/frontend/fs/
	go -C backend build -ldflags="-s -w" -o ../null3-server ./cmd/server
