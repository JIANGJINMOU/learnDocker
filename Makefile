SHELL := /bin/bash

APP := cede
PKG := example.com/containeredu
BIN := bin/$(APP)
COVERAGE := coverage.out

.PHONY: all build test fmt vet lint cover check-coverage hooks docker-build

all: build

build:
	GOOS=linux GOARCH=amd64 go build -o $(BIN) $(PKG)/cmd/cede

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: vet

test:
	bash scripts/coverage.sh

cover: test
	go tool cover -func=$(COVERAGE)

check-coverage: cover
	@total=$$(go tool cover -func=$(COVERAGE) | awk '/total:/ {print $$3}' | sed 's/%//'); \
	req=80; \
	echo "Total coverage: $$total%"; \
	awk -v t=$$total -v r=$$req 'BEGIN { if (t+0 < r) { print "Coverage below " r "%"; exit 1 } else { print "Coverage OK" } }'

hooks:
	git config core.hooksPath .githooks

docker-build:
	docker build -t containeredu-builder .
