build:
	@go build -o ./bin/shop ./*.go

run: build
	@./bin/*

drop:
	@go run ./dropdb/*.go
