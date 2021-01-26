install:
	@go mod download
	@go mod tidy

build: install
	@go build -o app main.go

run: build
	@./app

test_ping:
	@go test -v -race ./test/