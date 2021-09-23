# Needed SHELL since I'm using zsh
SHELL := /bin/bash
.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: format ## Run go lint in the current dir
	golint .

format:
	gofmt -d .

tidy:
	go mod tidy

compile: tidy lint     ## Compiles the program and leaves the binary in the current dir.
	go build .