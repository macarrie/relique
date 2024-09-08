build: clean
	goimports -w -local "github.com/macarrie/relique" cmd/**/*.go internal/**/*.go api/*.go
	go mod tidy
	go vet ./...
	go build -o relique cmd/relique/main.go

clean:
	rm -f ./relique

reset:
	rm -rf ~/.config/relique/db/relique.sqlite ~/.config/relique/storage/*