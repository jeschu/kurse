default: dist

build:
	go build -o 'dist/kurse' ./...

test:
	 go test -v ./...

clean:
	rm -rf dist

dist: test clean
	GOOS='darwin'  GOARCH='amd64' go build -a -o 'dist/darwin-amd64-kurse'  .
	GOOS='linux'   GOARCH='amd64' go build -a -o 'dist/linux-amd64-kurse'   .
	GOOS='linux'   GOARCH='arm'   go build -a -o 'dist/linux-arm-kurse'     .
	GOOS='linux'   GOARCH='arm64' go build -a -o 'dist/linux-arm64-kurse'   .
	GOOS='windows' GOARCH='arm64' go build -a -o 'dist/windows-arm64-kurse'  .

install:
	@echo "installing to ${GOPATH}/kurse ..."
	@go install

edit:
	code '/Users/jens/Library/Application Support/kurse/depot.yml'