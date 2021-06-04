OUTPUT := output
GO_BUILD_FLAGS ?= -tags musl


all: build
.PHONY: all

build:
	mkdir -p $(OUTPUT)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(OUTPUT)/web cmd/web/web.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(OUTPUT)/worker cmd/worker/worker.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(OUTPUT)/migration cmd/migration/migration.go
.PHONY: build