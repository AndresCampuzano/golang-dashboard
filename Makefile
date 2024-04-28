build:
	@go build -o bin/golang-dashboard

run: build
	@./bin/golang-dashboard

test:
	@go test -v ./...