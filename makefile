RABBIT_CONTAINER=rabbitmq
RABBIT_IMAGE=rabbitmq:3-management
CONDUCTOR=protheon-conductor
MIND=protheon-mind
VERSION:=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

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
	GOOS=linux GOARCH=amd64 go build -o bin/$(CONDUCTOR)-$(VERSION)-linux-amd64 cmd/conductor/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(MIND)-$(VERSION)-linux-amd64 cmd/mind/main.go

build-mac:
	@echo "Building Mac/arm64 binary..."
	GOOS=darwin GOARCH=arm64 go build -o bin/$(CONDUCTOR)-$(VERSION)-darwin-arm64 cmd/conductor/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/$(MIND)-$(VERSION)-darwin-arm64 cmd/mind/main.go

build-windows:
	@echo "Building Windows/amd64 binary..."
	GOOS=windows GOARCH=amd64 go build -o bin/$(CONDUCTOR)-$(VERSION)-windows.exe cmd/conductor/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(MIND)-$(VERSION)-windows.exe cmd/mind/main.go

release: build-linux build-mac build-windows
	@echo "Release build complete in ./bin"

run-server:
	./bin/protheon-darwin-arm64 --role=server

clean:
	rm -rf bin/*
