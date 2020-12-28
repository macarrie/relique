PROJECTNAME="relique"
GO=go
VERSION=$(shell git describe --tags --always)
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)
PACKAGE_NAME=relique_$(VERSION)_$(GOOS)_$(GOARCH)

MAKEFLAGS += --silent

## all: Build all relique components from scratch
all: clean server client cli

build/bin:
	mkdir -p $@

build/bin/relique-server: build/bin
	go build -o $@ cmd/relique-server/main.go

build/bin/relique-client: build/bin
	go build -o $@ cmd/relique-client/main.go

build/bin/relique: build/bin
	go build -o $@ cmd/relique/main.go

## server: Build relique-server
server:
	rm -f build/bin/relique-server
	$(MAKE) build/bin/relique-server

## client: Build relique-client
client:
	rm -f build/bin/relique-client
	$(MAKE) build/bin/relique-client

## cli: Build relique cli tool
cli:
	rm -f build/bin/relique
	$(MAKE) build/bin/relique

## check: Run code vet
check:
	go vet ./...

## test: Run tests
test: check
	go test ./... -cover

## clean: Clean all build artefacts
clean:
	rm -rf build

.PHONY: clean server client cli test check
help: Makefile
	echo " Choose a command to run:"
	sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
