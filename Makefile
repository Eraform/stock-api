GOCMD=go
GOBUILD=$(GOCMD) build

BINARY_NAME=stocks-api
APP_NAME = zhenik/$(BINARY_BASE_NAME)

build:
	${PREBUILD_FLAGS}$(GOBUILD) -o $(BINARY_NAME) -v -ldflags="-s -w"
run:
	./$(BINARY_NAME)

restart: build run
