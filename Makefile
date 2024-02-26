build:
	@go build -o ./bin/shop ./*.go

run: build
	@./bin/*

test:
	@go test -v ./...