GOPATH := $(shell go env GOPATH)

default: dist

build:
	@echo ">> build <<"
	go build -o 'dist/kurse' ./...

clean:
	@echo ">> clean <<"
	rm -rf dist

ensuregox:
	@echo ">> ensure gox <<"
	@go install github.com/mitchellh/gox@latest

dist: clean ensuregox
	@echo ">> dist <<"
	@gox -osarch='darwin/amd64 linux/amd64 linux/arm linux/arm64 windows/amd64' -output 'dist/kurse_{{.OS}}-{{.Arch}}' .

install:
	@echo "installing to ${GOPATH}/kurse ..."
	@go install

edit:
	code '/Users/jens/Library/Application Support/kurse/portfolio.yml'