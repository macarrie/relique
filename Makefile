PROJECTNAME="relique"
CURRENT_TAG != git describe --tags
BUILD_OUTPUT_DIR=output

#MAKEFLAGS += --silent

build: clean webui bin

bin: 
	go mod tidy
	go build -ldflags="-X 'github.com/macarrie/relique/api.ReliqueVersion=$(CURRENT_TAG)'" -o $(BUILD_OUTPUT_DIR)/relique cmd/relique/main.go

webui:
	cd webui && npm ci && npm run build
	cp -r ./webui/dist internal/server/

clean:
	rm -f ./output/*
	rm -rf ./webui/dist internal/server/dist

reset:
	rm -rf ~/.config/relique/db/relique.sqlite ~/.config/relique/storage/* ~/.config/relique/catalog/*

test:
	docker build -t relique_tests -f test/Dockerfile_tests  .

docker: clean
	docker build -t relique -f build/package/Dockerfile  .

release: clean
	@if [ -z "$(TAG)" ]; then echo "Please provide tag with TAG=x.y.z"; exit 1; fi
	@git diff --exit-code --quiet || (echo "Please commit pending changes before creating release commit"; exit 1)
	@echo "Writing current version file"
	echo "$(TAG)" > .current_version

	git commit -am "Release v$(TAG)"
	$(MAKE) tag TAG=$(TAG)

tag:
	@if [ -z "$(TAG)" ]; then echo "Please provide tag with TAG=x.y.z"; exit 1; fi
	@echo "Creating git tag v$(TAG)"
	git tag v$(TAG)

.PHONY: build clean reset test docker release tag bin webui