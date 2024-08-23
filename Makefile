build:
	@go build -o bin/dfss

run: build
	@./bin/dfss

test:
	@go test ./...