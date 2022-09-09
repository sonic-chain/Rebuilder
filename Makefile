OUT := rebuilder
PKG := github.com/Fogmeta/filecoin-ipfs-data-rebuilder
VERSION := $(shell git describe --always --long)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: build

deps:
	go mod tidy && go mod vendor

test:
	@go test -v ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

static: vet

build: deps static
	go build -v -o ${OUT}
clean:
	rm -rf vendor && rm -rf ${OUT}
	@go clean

.PHONY: run server static vet
