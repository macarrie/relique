build: clean
	go mod tidy
	go build -o output/relique cmd/relique/main.go

clean:
	rm -f ./output/*

reset:
	rm -rf ~/.config/relique/db/relique.sqlite ~/.config/relique/storage/* ~/.config/relique/catalog/*

test:
	docker build -t relique_tests -f test/Dockerfile_tests  .

docker:
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

.PHONY: build clean reset test docker release tag