all: clean bin/relique-client bin/relique-server

build/bin:
	mkdir -p $@

build/bin/relique-server: build/bin
	go build -o $@ cmd/relique-server/main.go

build/bin/relique-client: build/bin
	go build -o $@ cmd/relique-client/main.go

build/bin/relique: build/bin
	go build -o $@ cmd/relique/main.go

server: build/bin/relique-server
client: build/bin/relique-client
cli: build/bin/relique

clean:
	rm -rf build

.PHONY: clean server client