all: clean bin/relique-client bin/relique-server

build/bin:
	mkdir -p $@

build/bin/relique-server: build/bin
	go build -o $@ cmd/relique-server/main.go

build/bin/relique-client: build/bin
	go build -o $@ cmd/relique-client/main.go

build/bin/relique: build/bin
	go build -o $@ cmd/relique/main.go

server:
	rm -f build/bin/relique-server
	$(MAKE) build/bin/relique-server
client:
	rm -f build/bin/relique-client
	$(MAKE) build/bin/relique-client
cli:
	rm -f build/bin/relique
	$(MAKE) build/bin/relique

clean:
	rm -rf build

.PHONY: clean server client