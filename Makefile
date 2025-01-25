.PHONY: build test clean lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=exo
BINARY_DIR=bin

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build

build: $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)

$(BINARY_DIR):
	mkdir -p $(BINARY_DIR)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

lint:
	golangci-lint run

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

install: build
	mv $(BINARY_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

uninstall:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Development helpers
fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet lint test
