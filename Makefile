APP_NAME ?= gitimpact-backend
VERSION ?= dev
GOFLAGS ?= -mod=vendor
IMAGE ?= gitimpact/backend:$(VERSION)
CONTAINER_NAME ?= gitimpact-backend
RUN_ARGS ?= -p 8080:8080
SERVICE_ARGS ?=

BACKEND_DIR := backend
SERVER_PKG := ./cmd/server
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)
LINUX_AMD64_BIN_PATH := $(BIN_DIR)/$(APP_NAME)-linux-amd64

.PHONY: build build-linux-amd64 test lint clean docker-build docker-run docker-push docker-build-run deploy

build:
	mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && GOFLAGS="$(GOFLAGS)" go build -trimpath -ldflags "-s -w" -o ../$(BIN_PATH) $(SERVER_PKG)

build-linux-amd64:
	mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS="$(GOFLAGS)" go build -trimpath -ldflags "-s -w" -o ../$(LINUX_AMD64_BIN_PATH) $(SERVER_PKG)

test:
	cd $(BACKEND_DIR) && GOFLAGS="$(GOFLAGS)" go test ./...

lint:
	cd $(BACKEND_DIR) && GOFLAGS="$(GOFLAGS)" golangci-lint run ./...

clean:
	rm -rf $(BIN_DIR) dist tmp

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm --name $(CONTAINER_NAME) $(RUN_ARGS) $(IMAGE) $(SERVICE_ARGS)

docker-push:
	docker push $(IMAGE)

docker-build-run: docker-build docker-run

deploy: docker-build docker-push
