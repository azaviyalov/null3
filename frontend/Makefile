.PHONY: prepare build lint clean

node_modules/.package-lock.json:
	npm install

prepare: node_modules/.package-lock.json

build: prepare
	ng build --configuration production

lint: prepare
	ng lint

clean:
	rm -rf dist .angular/cache node_modules
