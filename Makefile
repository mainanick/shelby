.PHONY: all build test clean run build-linux

# Go parameters
GOCMD=go
LDFLAGS=-ldflags "-s -w"
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=shellby


# All target
all: test build

build: 
	$(GOBUILD) $(LDFLAGS) -o bin/ -v ./... 

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -rf bin/

run: 
	$(GOCMD) run cmd/shellby/main.go

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/unix/ -v ./...