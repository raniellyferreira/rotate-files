BINDIR      := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
DIST_DIRS   := find * -type d ! -name '[0-9]*' -exec
TARGETS     := darwin/amd64 darwin/arm64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
TARGET_OBJS ?= darwin-amd64.tar.gz darwin-amd64.tar.gz.sha256 darwin-amd64.tar.gz.sha256sum darwin-arm64.tar.gz darwin-arm64.tar.gz.sha256 darwin-arm64.tar.gz.sha256sum linux-amd64.tar.gz linux-amd64.tar.gz.sha256 linux-amd64.tar.gz.sha256sum linux-386.tar.gz linux-386.tar.gz.sha256 linux-386.tar.gz.sha256sum linux-arm.tar.gz linux-arm.tar.gz.sha256 linux-arm.tar.gz.sha256sum linux-arm64.tar.gz linux-arm64.tar.gz.sha256 linux-arm64.tar.gz.sha256sum linux-ppc64le.tar.gz linux-ppc64le.tar.gz.sha256 linux-ppc64le.tar.gz.sha256sum linux-s390x.tar.gz linux-s390x.tar.gz.sha256 linux-s390x.tar.gz.sha256sum windows-amd64.zip windows-amd64.zip.sha256 windows-amd64.zip.sha256sum
BINNAME     ?= rotate

GOBIN         = $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN         = $(shell go env GOPATH)/bin
endif
GOX           = $(GOBIN)/gox
GOIMPORTS     = $(GOBIN)/goimports
ARCH          = $(shell uname -p)

# go option
PKG         := ./...
TAGS        :=
TESTS       := .
TESTFLAGS   :=
LDFLAGS     := -w -s
GOFLAGS     :=
CGO_ENABLED ?= 0

# Rebuild the binary if any of these files change
SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match | sed 's/^v//' 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(shell echo ${VERSION} | sed 's/^v//')
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/raniellyferreira/rotate-files/internal/version.version=${BINARY_VERSION}
endif

VERSION_METADATA = unreleased
# Clear the "unreleased" string in BuildMetadata
ifneq ($(GIT_TAG),)
	VERSION_METADATA =
endif

LDFLAGS += -X github.com/raniellyferreira/rotate-files/internal/version.metadata=${VERSION_METADATA}
LDFLAGS += -X github.com/raniellyferreira/rotate-files/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/raniellyferreira/rotate-files/internal/version.gitTreeState=${GIT_DIRTY}
LDFLAGS += $(EXT_LDFLAGS)

.PHONY: all
all: build

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd/rotate

# ------------------------------------------------------------------------------
#  install

.PHONY: install
install: build
	@install "$(BINDIR)/$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test: build
ifeq ($(ARCH),s390x)
test: TESTFLAGS += -v
else
test: TESTFLAGS += -race -v
endif
test: test-style
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: test-coverage
test-coverage:
	@echo
	@echo "==> Running unit tests with coverage <=="
	@environment/scripts/coverage.sh

.PHONY: test-style
test-style:
	GO111MODULE=on golangci-lint run
	@environment/scripts/validate-license.sh

.PHONY: coverage
coverage:
	@environment/scripts/coverage.sh

.PHONY: format
format: $(GOIMPORTS)
	GO111MODULE=on go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w -local github.com/raniellyferreira/rotate-files

# Generate golden files used in unit tests
.PHONY: gen-test-golden
gen-test-golden:
gen-test-golden: PKG = ./cmd/rotate ./pkg/rotate
gen-test-golden: TESTFLAGS = -update
gen-test-golden: test-unit

# ------------------------------------------------------------------------------
#  dependencies

# If go install is run from inside the project directory it will add the
# dependencies to the go.mod file. To avoid that we change to a directory
# without a go.mod file when downloading the following dependencies

$(GOX):
	(cd /; GO111MODULE=on go install github.com/mitchellh/gox@latest)

$(GOIMPORTS):
	(cd /; GO111MODULE=on go install golang.org/x/tools/cmd/goimports@latest)

# ------------------------------------------------------------------------------
#  release

.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross: $(GOX)
	GOFLAGS="-trimpath" GO111MODULE=on CGO_ENABLED=0 $(GOX) -parallel=3 -output="_dist/{{.OS}}-{{.Arch}}/$(BINNAME)" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/rotate

.PHONY: dist
dist:
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) mkdir -p ${VERSION} \; && \
		$(DIST_DIRS) tar -zcf ${VERSION}/rotate-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r ${VERSION}/rotate-{}.zip {} \; \
	)

.PHONY: fetch-dist
fetch-dist:
	mkdir -p _dist/${VERSION}/;
	\
	for obj in ${TARGET_OBJS} ; do \
		curl -sSL --output-dir _dist/${VERSION}/ -o rotate-$${obj} https://awapi-rotate.s3.sa-east-1.amazonaws.com/${VERSION}/rotate-$${obj} ; \
	done

.PHONY: sign
sign:
	for f in $$(ls _dist/*.{gz,zip,sha256,sha256sum} 2>/dev/null) ; do \
		gpg --armor --detach-sign $${f} ; \
	done

.PHONY: checksum
checksum:
	for f in $$(ls _dist/${VERSION}/*.{gz,zip} 2>/dev/null) ; do \
		shasum -a 256 "$${f}" | sed 's/_dist\///' > "$${f}.sha256sum" ; \
		shasum -a 256 "$${f}" | awk '{print $$1}' > "$${f}.sha256" ; \
	done

# ------------------------------------------------------------------------------

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)' ./_dist

.PHONY: info
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"