PROJECTNAME="nacosServiceExporter"

all: clean deps linux

clean:
	rm -rf dist/

deps:
	go mod vendor

linux:
	mkdir -p dist/linux/

	go build -o dist/linux/nacosServiceExporter cmd/main.go
