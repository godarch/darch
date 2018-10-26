THIS_FILE := $(lastword $(MAKEFILE_LIST))
GIT_COMMIT=$(shell git rev-parse HEAD)
# v1.0.1
TAG=$(TRAVIS_TAG)
# 1.0.1, or NA if no tag
VERSION=$(shell test -n "$(TAG)" && echo $(TAG) | cut -d "v" -f 2 || echo "NA")
# Where the files will be installed
DESTDIR=/
CURRENTDIR=$(shell pwd)
GO_BUILD_FLAGS=

.PHONY: clean_build build test lint vendor clean_bundle bundle install release ci_deps ci

default: build

clean_build:
	@echo "cleaning bin/"
	@rm -rf bin/
build: clean_build
	@echo "bin/darch"
	@go build ${GO_BUILD_FLAGS} -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)" -o bin/darch pkg/cmd/darch/main.go
clean_runc:
	@echo "cleaning tmp/runc"
	@rm -rf tmp/runc
test:
	@echo "testing"
	@go test ./pkg/...
lint:
	@echo "linting"
	@gometalinter.v2 --config .gometalinter.json ./pkg/...
vendor:
	@echo "vendoring"
	@vndr
install:
	@echo "installing to $(DESTDIR)"
	@echo "installing /usr/bin/darch"
	@install -D -m 755 bin/darch $(DESTDIR)/usr/bin/darch
	@echo "installing /etc/darch/hooks/fstab/hook"
	@install -D -m 755 scripts/hooks/fstab $(DESTDIR)/etc/darch/hooks/fstab/hook
	@echo "installing /etc/darch/hooks/hostname/hook"
	@install -D -m 755 scripts/hooks/hostname $(DESTDIR)/etc/darch/hooks/hostname/hook
	@echo "installing /etc/darch/hooks/ssh/hook"
	@install -D -m 755 scripts/hooks/ssh $(DESTDIR)/etc/darch/hooks/ssh/hook
	@echo "installing /etc/grub.d/60_darch"
	@install -D -m 755 scripts/grub-mkconfig-script $(DESTDIR)/etc/grub.d/60_darch
ci_deps:
	@echo "fetching golint"
	@go get -u github.com/golang/lint/golint
	@echo "fetching gometalinter.v2"
	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@echo "fetching github-release"
	@go get -u github.com/aktau/github-release
	@echo "fetching vndr"
	@go get -u github.com/LK4D4/vndr
ci: ci_deps lint build