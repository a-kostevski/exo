.PHONY: build test clean install

BINARY_NAME=exo
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)
	go clean

# Development helpers
fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

check: fmt vet lint test
