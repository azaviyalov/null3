.PHONY: all backend frontend prepare-frontend prepare build-backend build-frontend build test-backend-unit test-backend coverage-backend coverage-backend-html lint-backend lint-frontend lint format-backend format-frontend format clean-backend clean-frontend clean release

all: build-backend build-frontend

backend:
	$(MAKE) -C backend

frontend:
	$(MAKE) -C frontend

prepare-frontend:
	$(MAKE) -C frontend prepare

prepare: prepare-frontend

build-backend:
	$(MAKE) -C backend build

build-frontend:
	$(MAKE) -C frontend build

build: build-backend build-frontend

test-backend-unit:
	$(MAKE) -C backend test-unit

test-backend:
	$(MAKE) -C backend test

coverage-backend:
	$(MAKE) -C backend coverage

coverage-backend-html:
	$(MAKE) -C backend coverage-html

lint-backend:
	$(MAKE) -C backend lint

lint-frontend:
	$(MAKE) -C frontend lint

lint: lint-backend lint-frontend

format-backend:
	$(MAKE) -C backend format

format-frontend:
	$(MAKE) -C frontend format

format: format-backend format-frontend

clean-backend:
	$(MAKE) -C backend clean

clean-frontend:
	$(MAKE) -C frontend clean

clean: clean-backend clean-frontend
	rm null3-server || true

release:
	rm -rf backend/internal/core/frontend/fs/*
	$(MAKE) -C frontend build
	cp -r frontend/dist/frontend/browser/* backend/internal/core/frontend/fs/
	$(MAKE) -C backend build
	cp backend/bin/server ./null3-server
