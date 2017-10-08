#!/usr/bin/make

SHELL := /bin/bash
PWD    = $(shell pwd)

PKG  = .
BIN  = darch
GO  := $(realpath ./go)
GITVERSION  := $(realpath ./gitversion)

FIND_STD_DEPS = $(GO) list std | sort | uniq
FIND_PKG_DEPS = $(GO) list -f '{{join .Deps "\n"}}' $(PKG) | sort | uniq | grep -v "^_"
DEPS          = $(shell comm -23 <($(FIND_PKG_DEPS)) <($(FIND_STD_DEPS)))

VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
SEMVAR_VERSION = $(shell $(GITVERSION) | jq '.SemVer' --raw-output)
GOARCH=$(shell go env GOARCH)

.PHONY: %

default: test

help:
	@echo 'Make commands for darch:'
	@echo
	@echo 'Usage:'
	@echo '    make build				Compile the project.'
	@echo '    make clean				Clean the directory tree.'
	@echo '    make deps				Install all the dependencies for the project.'
	@echo '    make test				Run test tests.'
	@echo '    make test-deps			Install all the dependencies for the tests.'
	@echo '    make run				Run the program. Use ARGS to pass in arguments.'
	@echo

all: build
build: deps
	$(GO) build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.Version=${SEMVAR_VERSION}" -o bin/${BIN}
package: build
	cp grub-mkconfig-script bin/
	tar cvzpf bin/darch-${GOARCH}.tar.gz -C bin darch grub-mkconfig-script
clean:
	$(GO) clean -i $(PKG)
	rm -r bin/
deps:
	$(GO) get -d $(PKG)
	$(GO) install $(DEPS)
test: test-deps
	$(GO) test $(PKG)
test-deps: deps
	$(GO) get -d -t $(PKG)
	$(GO) test -i $(PKG)
run: all
	./bin/$(BIN) ${ARGS}
