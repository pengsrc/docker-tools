SHELL := /bin/bash

DIRS_TO_CHECK=$(shell ls -d */ | grep -v "vendor")
PKGS_TO_CHECK=$(shell go list ./... | grep -vE "/vendor")
OS=$(shell uname | awk '{print tolower($$0)}')

ifneq (${PKG},)
	PKGS_TO_CHECK="github.com/pengsrc/docker-tools/${PKG}"
endif

PACKAGE_NAME:=github.com/pengsrc/docker-tools
BINARY_NAME:=docker-tools
VERSION=$(shell cat constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all          to check, build, test and release"
	@echo "  check        to vet and lint"
	@echo "  build        to build all"
	@echo "  test         to run test"
	@echo "  install      to install into ${GOPATH}/bin"
	@echo "  uninstall    to uninstall"
	@echo "  clean        to clean the test and built files"
	@echo "  release      to build and release"
	@echo "  builder      to build all builder images"

.PHONY: all
all: check build test

.PHONY: check
check: format vet lint

.PHONY: format
format:
	@echo "Formatting packages using gofmt..."
	@find . -path '*/vendor/*' -prune -o -name '*.go' -type f -exec gofmt -s -w {} \;
	@echo "Done"

.PHONY: vet
vet:
	@echo "Checking packages using go tool vet, skip vendor packages..."
	@go tool vet -all ${DIRS_TO_CHECK}
	@echo "Done"

.PHONY: lint
lint:
	@echo "Checking packages using golint, skip vendor packages..."
	@lint=$$(for pkg in ${PKGS_TO_CHECK}; do golint $${pkg}; done); \
	 lint=$$(echo "$${lint}"); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "Done"

.PHONY: build
build: check
	@mkdir -p ./bin
	@echo "Building ${BINARY_NAME}..."
	@GOOS=${OS} GOARCH=amd64 go build -o ./bin/${BINARY_NAME} .
	@echo "Done"

.PHONY: test
test:
	@echo "Running test..."
	@go test -v ${PKGS_TO_CHECK}
	@echo "Done"

.PHONY: install
install: build
	@if [[ -z "${GOPATH}" ]]; then echo "ERROR: $GOPATH not found."; exit 1; fi
	@echo "Installing into ${GOPATH}/bin/${BINARY_NAME}..."
	@cp ./bin/${BINARY_NAME} ${GOPATH}/bin/${BINARY_NAME}
	@echo "Done"

.PHONY: uninstall
uninstall:
	@if [[ -z "${GOPATH}" ]]; then echo "ERROR: $GOPATH not found."; exit 1; fi
	@echo "Uninstalling ${BINARY_NAME}..."
	rm -f ${GOPATH}/bin/${BINARY_NAME}
	@echo "Done"

.PHONY: clean
clean:
	@echo "Clean the test and built files"
	rm -rf ./bin
	rm -rf ./release
	rm -rf ./coverage
	@echo "Done"

.PHONY: release
release:
	@echo "Release ${BINARY_NAME}"
	mkdir -p ./release
	@echo "for Linux"
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux/${BINARY_NAME} .
	mkdir -p ./release
	tar -C ./bin/linux/ -czf ./release/${BINARY_NAME}-v${VERSION}-linux_amd64.tar.gz ${BINARY_NAME}
	@echo "for macOS"
	mkdir -p ./bin/linux
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/${BINARY_NAME} .
	tar -C ./bin/darwin/ -czf ./release/${BINARY_NAME}-v${VERSION}-darwin_amd64.tar.gz ${BINARY_NAME}
	@echo "for Windows"
	mkdir -p ./bin/windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows/${BINARY_NAME}.exe .
	cd ./bin/windows/ && zip ../../release/${BINARY_NAME}-v${VERSION}-windows_amd64.zip ${BINARY_NAME}.exe
	@echo "ok"
