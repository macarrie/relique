build: clean
	go mod tidy
	go build -o output/relique cmd/relique/main.go

clean:
	rm -f ./output/*

reset:
	rm -rf ~/.config/relique/db/relique.sqlite ~/.config/relique/storage/*

test:
	docker build -t relique_tests -f test/Dockerfile_tests  .
	docker run -it relique_tests go vet ./...
	docker run -it relique_tests go test -cover ./...

docker:
	docker build -t relique -f build/package/Dockerfile  .

.PHONY: build clean reset test docker