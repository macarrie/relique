PROJECTNAME="relique"
GO=go
VERSION=$(shell git describe --tags --always)
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)
PACKAGE_NAME=relique_$(VERSION)_$(GOOS)_$(GOARCH)

#MAKEFLAGS += --silent

BUILD_OUTPUT_DIR=output

## all: Build all relique components from scratch
all: clean build

## build: Build relique package distribution
build: clean $(BUILD_OUTPUT_DIR)
	./scripts/build.sh --output-dir "$(BUILD_OUTPUT_DIR)"

## server: Build relique-server
server:
	rm -f build/bin/relique-server
	$(MAKE) build

## client: Build relique-client
client:
	rm -f build/bin/relique-client
	$(MAKE) build

## cli: Build relique cli tool
cli:
	rm -f build/bin/relique
	$(MAKE) build

## check: Run code vet
check:
	go vet ./...

## test: Run tests
test: check
	go test ./... -cover

## certs: Generate self signed ssl certs to help start a quick relique configuration while getting real certs
certs:
	rm -rf build/certs/*
	mkdir -p build/certs
	echo  -e "[req]\ndistinguished_name=req\n[san]\nsubjectAltName=DNS.1:localhost,DNS.2:relique" > tmp.certs
	openssl req \
		-x509 \
		-newkey rsa:4096 \
		-sha256 \
		-days 3650 \
		-nodes \
		-keyout build/certs/key.pem \
		-out build/certs/cert.pem \
		-subj '/CN=relique' \
		-extensions san \
		-config tmp.certs
	rm tmp.certs

## install: Install
install:
	./scripts/install.sh --prefix "$(INSTALL_ROOT)" --src "$(INSTALL_SRC)" $(INSTALL_ARGS)

## clean: Clean all build artefacts
clean:
	rm -rf output

$(BUILD_OUTPUT_DIR):
	mkdir -p $@

## tar: Package sources to tar for rpm build
tar: $(BUILD_OUTPUT_DIR)
	tar --exclude "$(BUILD_OUTPUT_DIR)" -zcf $(BUILD_OUTPUT_DIR)/relique-$(VERSION).src.tar.gz .

~/rpmbuild:
	rpmdev-setuptree

build_single_rpm:
	sed "s/__VERSION__/$(VERSION)/" build/package/rpm/$(rpm).spec.tpl > ~/rpmbuild/SPECS/$(rpm).spec
	rpmlint ~/rpmbuild/SPECS/$(rpm).spec
	rpmbuild -ba ~/rpmbuild/SPECS/$(rpm).spec


## rpm: Build rpm packages
rpm: clean ~/rpmbuild tar
	cp $(BUILD_OUTPUT_DIR)/relique-$(VERSION).src.tar.gz ~/rpmbuild/SOURCES/
	$(MAKE) build_single_rpm rpm=relique-client
	$(MAKE) build_single_rpm rpm=relique-server

.PHONY: help clean server client cli test check certs install build_single_rpm rpm tar build
help: Makefile
	echo " Choose a command to run:"
	sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
