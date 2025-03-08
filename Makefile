build:
	@go build -o bin/thoughtHub_Backend cmd/main.go

test:
	@go test -v ./...

run: build
	@exec ./bin/thoughthub_Backend

migration-create:
	@migrate create -ext sql -dir cmd/migrate/migrations -format "20060102150405" $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down