.PHONY: all backend frontend prepare-backend prepare-frontend prepare build-backend build-frontend build lint-backend lint-frontend lint format-backend format-frontend format clean-backend clean-frontend clean release

all: build-backend build-frontend

backend:
	$(MAKE) -C backend

frontend:
	$(MAKE) -C frontend

prepare-backend:
	$(MAKE) -C backend prepare

prepare-frontend:
	$(MAKE) -C frontend prepare

prepare: prepare-backend prepare-frontend

build-backend:
	$(MAKE) -C backend build

build-frontend:
	$(MAKE) -C frontend build

build: build-backend build-frontend

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
