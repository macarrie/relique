PROJECTNAME="relique"
GO=go
VERSION=$(shell git describe --tags --always)
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)
PACKAGE_NAME=relique_$(VERSION)_$(GOOS)_$(GOARCH)

MAKEFLAGS += --silent

BUILD_OUTPUT_DIR=output
INSTALL_SRC=output
INSTALL_ROOT=/

all: clean build ## Build all relique components from scratch

build: clean $(BUILD_OUTPUT_DIR) ## Build relique package distribution
	./scripts/build.sh --output-dir "$(BUILD_OUTPUT_DIR)"

server: ## Build relique-server
	rm -f build/bin/relique-server
	$(MAKE) build

client: ## Build relique-client
	rm -f build/bin/relique-client
	$(MAKE) build

check: ## Run code vet
	go vet ./...
	staticcheck ./...

test: check ## Run tests
	# Parallel db setup during unit tests can create errors (read only db).Use -p 1 to ensure tests are run sequentially
	sudo -u relique -g relique go test -p 1 ./... -cover

install: ## Install relique
	./scripts/install.sh --prefix "$(INSTALL_ROOT)" --src "$(INSTALL_SRC)" $(INSTALL_ARGS)

clean: ## Clean all build artefacts
	rm -rf output
	$(MAKE) -C build/package/freebsd/relique-client clean
	$(MAKE) -C build/package/freebsd/relique-server clean

$(BUILD_OUTPUT_DIR):
	mkdir -p $@

tar: $(BUILD_OUTPUT_DIR) ## Package sources to tar for rpm build
	tar --exclude "$(BUILD_OUTPUT_DIR)" --exclude "test" -zcf $(BUILD_OUTPUT_DIR)/relique-$(VERSION).src.tar.gz .

~/rpmbuild:
	rpmdev-setuptree

build_single_rpm:
	sed "s/__VERSION__/$(VERSION)/" build/package/rpm/$(rpm).spec.tpl > ~/rpmbuild/SPECS/$(rpm).spec
	rpmlint ~/rpmbuild/SPECS/$(rpm).spec
	rpmbuild -ba ~/rpmbuild/SPECS/$(rpm).spec


rpm: clean ~/rpmbuild tar ## Build rpm packages
	cp $(BUILD_OUTPUT_DIR)/relique-$(VERSION).src.tar.gz ~/rpmbuild/SOURCES/
	$(MAKE) build_single_rpm rpm=relique-client
	$(MAKE) build_single_rpm rpm=relique-server

port_makesum: ## Compute port sum
	$(MAKE) -C build/package/freebsd/relique-client clean makesum
	$(MAKE) -C build/package/freebsd/relique-server clean makesum

freebsd: port_makesum ## Build freebsd packages
	$(MAKE) -C build/package/freebsd/relique-client package
	$(MAKE) -C build/package/freebsd/relique-server package

.PHONY: help clean server client cli test check certs install build_single_rpm rpm tar build
help: Makefile ## Show this help
	echo " Choose a command run in "$(PROJECTNAME)":"
	@grep -E '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
