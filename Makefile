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

## test: run all tests
test:
	go test ./...
.PHONY: test

## test/perft: run perft tests
test/perft:
	go test ./shogi ./tests_perft -covermode=count -coverpkg=github.com/vinymeuh/hifumi/shogi 
.PHONY: test/perft

## bench/perft: run perft benchmarks
bench/perft:
	go test ./tests_perft -bench=. -run=^# -benchmem -memprofile memprofile.out -cpuprofile profile.out
.PHONY: bench/perft

## bench/perft/cpu: run perft benchmarks then pprof cpu usage 
bench/perft/cpu: bench/perft
	go tool pprof profile.out
.PHONY: bench/perft/cpu

## bench/perft/mem: run perft benchmarks then pprof memory usage 
bench/perft/mem: bench/perft
	go tool pprof memprofile.out
.PHONY: bench/perft/mem
