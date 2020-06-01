fmt:
	@goimports -w .

lint:
	@golint ./...

run-server:
	@go run ./cmd/pinmonl-server/