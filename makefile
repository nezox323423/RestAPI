ifeq (,$(wildcard .env))
    $(error .env file is missing)
endif

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build
build:
	go build -v ./

.PHONY: run
run:
	go run ./

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build