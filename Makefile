generate-mocks:
	mockgen -source=connection.go -destination=mock/connection.go -package=mock Connection
	mockgen -source=middleware.go -destination=mock/middleware.go -package=mock Middleware
	mockgen -destination=mock/mock_store.go -package=mock github.com/metinorak/varto Store
	mockgen -destination=mock/mock_topic.go -package=mock github.com/metinorak/varto Topic

test:
	go test -v -cover ./...
