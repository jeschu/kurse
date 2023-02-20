default: dist

build:
	go build -v -o 'dist/kurse' ./...

test:
	 go test -v ./...

clean:
	rm -rf dist

dist: clean test showdist
	GOOS='darwin'  GOARCH='amd64' go build -x -a -o 'dist/darwin-amd64-kurse'  .
	GOOS='linux'   GOARCH='amd64' go build -x -a -o 'dist/linux-amd64-kurse'   .
	GOOS='linux'   GOARCH='arm'   go build -x -a -o 'dist/linux-arm-kurse'     .
	GOOS='linux'   GOARCH='arm64' go build -x -a -o 'dist/linux-arm64-kurse'   .
	GOOS='windows' GOARCH='arm64' go build -x -a -o 'dist/windows-arm64-kurse'  .

showdist:
	ls -hl dist/*