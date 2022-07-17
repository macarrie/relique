PROJECTNAME="relique"
GO=go
VERSION != cat .current_version
GOOS != $(GO) env GOOS
GOARCH != $(GO) env GOARCH
PACKAGE_NAME=relique_$(VERSION)_$(GOOS)_$(GOARCH)

UNAME != uname

MAKEFLAGS += --silent

BUILD_OUTPUT_DIR=output
INSTALL_SRC=output
INSTALL_ROOT=/

all: clean build ## Build all relique components from scratch

build: clean $(BUILD_OUTPUT_DIR) ## Build entire relique package distribution
	$(MAKE) build_client $(BUILD_OUTPUT_DIR)
	$(MAKE) build_server $(BUILD_OUTPUT_DIR)

build_server: $(BUILD_OUTPUT_DIR) ## Build relique server package distribution
	./scripts/build.sh --server --output-dir "$(BUILD_OUTPUT_DIR)"

build_client: $(BUILD_OUTPUT_DIR) ## Build relique client package distribution
	./scripts/build.sh --client --output-dir "$(BUILD_OUTPUT_DIR)"

server: ## Build relique-server
	rm -f output/bin/relique-server
	$(MAKE) build_server

client: ## Build relique-client
	rm -f output/bin/relique-client
	$(MAKE) build_client

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
	if [ "$(UNAME)" = "FreeBSD" ]; then make -C build/package/freebsd/relique-client clean; fi
	if [ "$(UNAME)" = "FreeBSD" ]; then make -C build/package/freebsd/relique-server clean; fi
	rm -f build/package/freebsd/relique-server/distinfo
	rm -f build/package/freebsd/relique-client/distinfo

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

prepare_release: ## Template files with current version and release info
	@echo "Templating freebsd port Makefiles"
	sed "s/__VERSION__/$(VERSION)/" build/package/freebsd/relique-client/Makefile.tpl > build/package/freebsd/relique-client/Makefile
	sed "s/__VERSION__/$(VERSION)/" build/package/freebsd/relique-server/Makefile.tpl > build/package/freebsd/relique-server/Makefile

release: clean ## Create relique release $TAG
	@if [ -z "$(TAG)" ]; then echo "Please provide tag with TAG=vx.y.z"; exit 1; fi
	@git diff --exit-code --quiet || (echo "Please commit pending changes before creating release commit"; exit 1)
	@echo "Writing current version file"
	echo "$(TAG)" > .current_version
	$(MAKE) prepare_release VERSION=$(TAG)
	git commit -am "Release v$(TAG)"
	$(MAKE) tag TAG=$(TAG)

tag:
	@if [ -z "$(TAG)" ]; then echo "Please provide tag with TAG=vx.y.z"; exit 1; fi
	@echo "Creating git tag v$(TAG)"
	git tag v$(TAG)

docker:
	@echo "Building Docker image"
	docker build --network host -t macarrie/relique-server:v$(VERSION) -f build/package/docker/server/Dockerfile .
	docker build --network host -t macarrie/relique-client:v$(VERSION) -f build/package/docker/client/Dockerfile .


.PHONY: help clean server client cli test check certs install build_single_rpm rpm tar build
help: Makefile ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
