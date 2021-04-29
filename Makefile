PROJECTNAME="nacosServiceExporter"

all: clean deps linux

clean:
	rm -rf dist/

deps:
	go mod vendor

linux:
	mkdir -p dist/linux/bin

	go build -o dist/linux/bin/nacosServiceExporter cmd/main.go
