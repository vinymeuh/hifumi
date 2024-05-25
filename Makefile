SHELL := /usr/bin/env bash -o pipefail

EXE_NAME=hifumi

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help

## build: build the application
build:
	go build -o ${EXE_NAME} cmd/hifumi/main.go
.PHONY: build

## run: run the application
run: build
	@./hifumi
.PHONY: run

## test: run go tests
test:
	go test ./...
.PHONY: test

## test/debug: run movegen debug tests
test/debug:
	../perfttester/perfttester -d debug hifumi
.PHONY: test/debug

## test/perft: run perft tests
test/perft:
	../perfttester/perfttester hifumi
.PHONY: test/perft

