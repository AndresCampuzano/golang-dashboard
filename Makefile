build:
	@go build -o bin/golang-gobank

run: build
	@./bin/golang-gobank

test:
	@go test -v ./...