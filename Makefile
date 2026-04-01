APP_NAME ?= gitimpact-backend
VERSION ?= dev
GOFLAGS ?= -mod=vendor
IMAGE ?= gitimpact/backend:$(VERSION)
CONTAINER_NAME ?= gitimpact-backend
RUN_ARGS ?= -p 8080:8080
SERVICE_ARGS ?=
POWERSHELL ?= powershell -ExecutionPolicy Bypass -File

BACKEND_DIR := backend
SERVER_PKG := ./cmd/server
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)
LINUX_AMD64_BIN_PATH := $(BIN_DIR)/$(APP_NAME)-linux-amd64
GO_TMP_DIR := .tmp-go-build

ifeq ($(OS),Windows_NT)
MKDIR_BIN = if not exist "$(BIN_DIR)" mkdir "$(BIN_DIR)"
MKDIR_GO_TMP = if not exist "$(GO_TMP_DIR)" mkdir "$(GO_TMP_DIR)"
SET_GOFLAGS = set "GOFLAGS=$(GOFLAGS)" && set "GOTMPDIR=../$(GO_TMP_DIR)" &&
SET_LINUX_ENV = set "CGO_ENABLED=0" && set "GOOS=linux" && set "GOARCH=amd64" && set "GOFLAGS=$(GOFLAGS)" && set "GOTMPDIR=../$(GO_TMP_DIR)" &&
CLEAN_DIRS = powershell -NoProfile -Command "foreach ($$p in @('$(BIN_DIR)','dist','tmp','$(GO_TMP_DIR)')) { if (Test-Path $$p) { Remove-Item -Recurse -Force $$p } }"
else
MKDIR_BIN = mkdir -p $(BIN_DIR)
MKDIR_GO_TMP = mkdir -p $(GO_TMP_DIR)
SET_GOFLAGS = GOFLAGS="$(GOFLAGS)"
SET_LINUX_ENV = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS="$(GOFLAGS)"
CLEAN_DIRS = rm -rf $(BIN_DIR) dist tmp $(GO_TMP_DIR)
endif

.PHONY: build build-linux-amd64 test lint clean docker-build docker-run docker-push docker-build-run deploy frontend-build release-build package-offline verify-offline

build:
	$(MKDIR_BIN)
	$(MKDIR_GO_TMP)
	cd $(BACKEND_DIR) && $(SET_GOFLAGS) go build -trimpath -ldflags "-s -w" -o ../$(BIN_PATH) $(SERVER_PKG)

build-linux-amd64:
	$(MKDIR_BIN)
	$(MKDIR_GO_TMP)
	cd $(BACKEND_DIR) && $(SET_LINUX_ENV) go build -trimpath -ldflags "-s -w" -o ../$(LINUX_AMD64_BIN_PATH) $(SERVER_PKG)

test:
	$(MKDIR_GO_TMP)
	cd $(BACKEND_DIR) && $(SET_GOFLAGS) go test ./...

lint:
	$(MKDIR_GO_TMP)
	cd $(BACKEND_DIR) && $(SET_GOFLAGS) golangci-lint run ./...

clean:
	$(CLEAN_DIRS)

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm --name $(CONTAINER_NAME) $(RUN_ARGS) $(IMAGE) $(SERVICE_ARGS)

docker-push:
	docker push $(IMAGE)

docker-build-run: docker-build docker-run

deploy: docker-build docker-push

frontend-build:
	$(POWERSHELL) scripts/build-frontend.ps1

release-build:
	$(POWERSHELL) scripts/build-release.ps1

package-offline:
	$(POWERSHELL) scripts/package-offline.ps1

verify-offline:
	$(POWERSHELL) scripts/verify-offline.ps1
