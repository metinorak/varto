generate-mocks:
	mockgen -source=connection.go -destination=mock/connection.go -package=mock Connection
	mockgen -source=middleware.go -destination=mock/middleware.go -package=mock Middleware

test:
	go test -v -cover ./...