.PHONY: prepare build lint clean

node_modules/.package-lock.json:
	npm ci

prepare: node_modules/.package-lock.json

build: prepare
	ng build --configuration production

lint: prepare
	ng lint

format: prepare
	npm run prettier

clean:
	rm -rf dist .angular/cache node_modules
