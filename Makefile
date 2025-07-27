.PHONY: all backend frontend prepare-backend prepare-frontend prepare build-backend build-frontend build lint-backend lint-frontend lint format-backend format clean-backend clean-frontend clean release

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

format: format-backend

clean-backend:
	$(MAKE) -C backend clean

clean-frontend:
	$(MAKE) -C frontend clean

clean: clean-backend clean-frontend
	rm null3-server || true

release:
	$(MAKE) -C frontend build
	mkdir -p backend/internal/frontend/dist
	cp -r frontend/dist/frontend/browser/* backend/internal/frontend/dist/browser/
	$(MAKE) -C backend build
	rm -rf backend/internal/frontend/dist/browser/*
	cp backend/bin/server ./null3-server
