# make file for HTTP_SERVER

APP_NAME=HTTP_SERVER
BIN_DIR=bin

.PHONY: all run build tidy test fmt vet clean install

all: fmt vet tidy clean build run

run:
	go run ./...

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./...

tidy:
	go mod tidy

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(BIN_DIR)

install:
	go install ./...
