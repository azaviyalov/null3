.PHONY: prepare build lint clean

node_modules/.package-lock.json:
	npm ci

prepare: node_modules/.package-lock.json

build: prepare
	npm run build

lint: prepare
	npm run lint

format: prepare
	npm run prettier

clean:
	rm -rf dist .angular/cache node_modules
