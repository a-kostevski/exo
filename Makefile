.PHONY: build install test clean

build:
	go build -v ./...

install:
	go install ./cmd/...

test:
	go test -v ./...

clean:
	go clean
	rm -f bin/*
