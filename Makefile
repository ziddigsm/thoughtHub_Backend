build:
	@go build -o bin/thoughtHub_Backend cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/thoughtHub_Backend