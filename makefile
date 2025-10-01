RABBIT_CONTAINER=rabbitmq
RABBIT_IMAGE=rabbitmq:3-management
BINARY=protheon

.PHONY: rabbit-up rabbit-down rabbit-logs

rabbit-up:
	docker run -d --rm \
		--hostname rabbit \
		--name $(RABBIT_CONTAINER) \
		-p 5672:5672 \
		-p 15672:15672 \
		$(RABBIT_IMAGE)
	docker exec -it rabbitmq rabbitmqctl add_user protheon secretpassword
	docker exec -it rabbitmq rabbitmqctl set_user_tags protheon administrator
	docker exec -it rabbitmq rabbitmqctl set_permissions -p / protheon ".*" ".*" ".*"

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
