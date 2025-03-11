.PHONY: tests coverage

run-server:
	@cd server && air

run-client:
	@cd example/custom-client && air

run-redis-client:
	@cd example/sample-go-redis && air

tests:
	@cd server && go test -v -coverpkg=./... -coverprofile=coverage.out ./...

coverage:
	@cd server && go tool cover -html=coverage.out