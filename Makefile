generate-mocks:
	mockgen -source=connection.go -destination=mock/connection.go -package=mock Connection

test:
	go test -v -cover ./...