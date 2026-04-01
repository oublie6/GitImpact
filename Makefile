APP_NAME ?= gitimpact-backend
VERSION ?= dev
GOFLAGS ?= -mod=vendor
IMAGE ?= gitimpact/backend:$(VERSION)
CONTAINER_NAME ?= gitimpact-backend
RUN_ARGS ?= -p 8080:8080
SERVICE_ARGS ?=

BACKEND_DIR := backend
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

.PHONY: build test clean docker-build docker-run docker-push docker-build-run deploy

build:
	mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && GOFLAGS="$(GOFLAGS)" go build -o ../$(BIN_PATH) ./cmd/server

test:
	cd $(BACKEND_DIR) && GOFLAGS="$(GOFLAGS)" go test ./...

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
