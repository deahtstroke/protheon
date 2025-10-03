RABBIT_CONTAINER=rabbitmq
RABBIT_IMAGE=rabbitmq:3-management
BINARY=protheon

.PHONY: rabbit-up rabbit-down rabbit-logs

rabbit-up:
	docker run -d --rm\
		--hostname rabbit \
		--name $(RABBIT_CONTAINER) \
		-e RABBITMQ_DEFAULT_USER=protheon \
		-e RABBITMQ_DEFAULT_PASS=secretpassword \
		-p 5672:5672 \
		-p 15672:15672 \
		$(RABBIT_IMAGE)
rabbit-down:
	docker stop $(RABBIT_CONTAINER)

rabbit-logs:
	docker logs -f $(RABBIT_CONTAINER)

build-linux:
	@echo "Building Linux/amd64 binary..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY)-linux-amd64 cmd/protheon/main.go

build-mac:
	@echo "Building Mac/arm64 binary..."
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY)-darwin-arm64 cmd/protheon/main.go

release: build-linux build-mac
	@echo "Release build complete in ./bin"

run-server:
	./bin/protheon-darwin-arm64 --role=server
