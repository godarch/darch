THIS_FILE := $(lastword $(MAKEFILE_LIST))
GIT_COMMIT=$(shell git rev-parse HEAD)
# v1.0.1
TAG=$(TRAVIS_TAG)
# 1.0.1, or NA if no tag
VERSION=$(shell test -n "$(TAG)" && echo $(TAG) | cut -d "v" -f 2 || echo "NA")
GOARCH=$(shell go env GOARCH)
BUNDLE_RUNTIME="no"
CONTAINERD_COMMIT="eed3b1c804bff194e2b53685a2cd95077e8aaaba"
RUNC_COMMIT="9f9c96235cc97674e935002fc3d78361b696a69e"
DESTINATION_BUNDLE_FILE_NAME="darch-$(GOARCH).tar.gz"
DESTINATION_BUNDLE_FILE_NAME_WITH_RUNTIME="darch-$(GOARCH)-with-runtime.tar.gz"
# Where the files will be installed
DESTDIR=/
CURRENTDIR=$(shell pwd)

.PHONY: clean_build build test lint vendor clean_bundle bundle install release ci_deps ci

default: build

clean_build:
	@echo "cleaning bin/"
	@rm -rf bin/
build: clean_build
	@echo "bin/darch"
	@go build -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)" -o bin/darch pkg/cmd/darch/main.go
clean_containerd:
	@echo "cleaning tmp/containerd"
	@rm -rf tmp/containerd
bundle_containerd: clean_containerd
	@echo "building containerd"
	@mkdir -p tmp/containerd/github.com/containerd
	@git clone https://github.com/containerd/containerd tmp/containerd/src/github.com/containerd/containerd
	@cd tmp/containerd/src/github.com/containerd/containerd && git checkout $(CONTAINERD_COMMIT)
	@GOPATH=$(CURRENTDIR)/tmp/containerd && cd $(CURRENTDIR)/tmp/containerd/src/github.com/containerd/containerd && make && make install DESTDIR=$(CURRENTDIR)/bundle/usr/
	@sed -i "s|/usr/local/bin/containerd|/usr/bin/containerd|" tmp/containerd/src/github.com/containerd/containerd/containerd.service
	@sed -i "s|/sbin/modprobe|/usr/bin/env modprobe|" tmp/containerd/src/github.com/containerd/containerd/containerd.service
	@install -Dm644 tmp/containerd/src/github.com/containerd/containerd/containerd.service bundle/usr/lib/systemd/system/containerd.service
clean_runc:
	@echo "cleaning tmp/runc"
	@rm -rf tmp/runc
bundle_runc: clean_runc
	@echo "building runc"
	@mkdir -p tmp/runc/github.com/opencontainers
	@git clone https://github.com/opencontainers/runc tmp/runc/src/github.com/opencontainers/runc
	@cd tmp/runc/src/github.com/opencontainers/runc && git checkout $(RUNC_COMMIT)
	@GOPATH=$(CURRENTDIR)/tmp/runc && cd $(CURRENTDIR)/tmp/runc/src/github.com/opencontainers/runc && make && make install BINDIR=$(CURRENTDIR)/bundle/usr/bin
test:
	@echo "testing"
	@go test ./pkg/...
lint:
	@echo "linting"
	@gometalinter.v2 --config .gometalinter.json ./pkg/...
vendor:
	@echo "vendoring"
	@vndr
clean_bundle:
	@echo "cleaning bundle/"
	@rm -rf bundle/
bundle: clean_bundle
	@echo "installing bundle/usr/bin/darch"
	@install -D -m 755 bin/darch bundle/usr/bin/darch
	@echo "installing bundle/etc/darch/hooks/fstab/hook"
	@install -D -m 755 scripts/hooks/fstab bundle/etc/darch/hooks/fstab/hook
	@echo "installing bundle/etc/darch/hooks/hostname/hook"
	@install -D -m 755 scripts/hooks/hostname bundle/etc/darch/hooks/hostname/hook
	@echo "installing bundle/etc/darch/hooks/ssh/hook"
	@install -D -m 755 scripts/hooks/ssh bundle/etc/darch/hooks/ssh/hook
	@echo "installing bundle/etc/grub.d/60_darch"
	@install -D -m 755 scripts/grub-mkconfig-script bundle/etc/grub.d/60_darch
ifeq ($(_BUNDLE_RUNTIME), yes)
	@echo "bundling containerd and runc"
	@$(MAKE) -f $(THIS_FILE) bundle_containerd
	@$(MAKE) -f $(THIS_FILE) bundle_runc
endif
ifeq ($(_BUNDLE_RUNTIME), yes)
	@echo "generating $(DESTINATION_BUNDLE_FILE_NAME_WITH_RUNTIME)"
	@mkdir -p output
	@rm -f output/$(DESTINATION_BUNDLE_FILE_NAME_WITH_RUNTIME)
	@tar -czpf output/$(DESTINATION_BUNDLE_FILE_NAME_WITH_RUNTIME) -C bundle/ .
else
	@echo "generating $(DESTINATION_BUNDLE_FILE_NAME)"
	@mkdir -p output
	@rm -f output/$(DESTINATION_BUNDLE_FILE_NAME)
	@tar -czpf output/$(DESTINATION_BUNDLE_FILE_NAME) -C bundle/ .
endif
ifeq ($(BUNDLE_RUNTIME), yes)
	@$(MAKE) -f $(THIS_FILE) bundle BUNDLE_RUNTIME="no" _BUNDLE_RUNTIME="yes"
endif
install: bundle
	@echo "installing to $(DESTDIR)"
	@tar xpzf output/$(DESTINATION_BUNDLE_FILE_NAME) -C $(DESTDIR)
release:
ifdef TRAVIS_TAG
	@echo "creating github release"
	@github-release release --user pauldotknopf --repo darch --tag $(TAG)
	@echo "uploading $(DESTINATION_BUNDLE_FILE_NAME)"
	@github-release upload --user pauldotknopf --repo darch --tag $(TAG) --name $(DESTINATION_BUNDLE_FILE_NAME) --file output/$(DESTINATION_BUNDLE_FILE_NAME)
ifeq ($(BUNDLE_RUNTIME), yes)
	@echo "uploading $(DESTINATION_BUNDLE_FILE_NAME_WITH_RUNTIME)"
	@github-release upload --user pauldotknopf --repo darch --tag $(TAG) --name $(DESTINATION_BUNDLE_FILE_NAME_WITH_RUNTIME) --file output/$(DESTINATION_BUNDLE_FILE_NAME_WITH_RUNTIME)
endif
	@echo "updating aur"
	scripts/aur/deploy-aur $(VERSION)
else
	@echo "not a tagged build, skipping release"
endif
ci_deps:
	@echo "fetching golint"
	@go get -u github.com/golang/lint/golint
	@echo "fetching gometalinter.v2"
	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@echo "fetching github-release"
	@go get -u github.com/aktau/github-release
	@echo "fetching vndr"
	@go get -u github.com/LK4D4/vndr
ci: ci_deps vendor lint build bundle release