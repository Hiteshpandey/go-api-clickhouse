build:
	@go build -o bin/velocityApi cmd/main.go

test:
	@go test -v ./..

run:
	@./bin/velocityApi