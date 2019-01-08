.PHONY: build  package start-docker run deps test pack
DIST_PATH = dist
RELEASE_PATH = $(DIST_PATH)/release
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)
GO = go
GINKGO = ginkgo
GO_LINKER_FLAGS ?= -ldflags \
	   "-X 'dappswin/common.BuildDate=$(BUILD_DATE)' \
	   -X dappswin/common.BuildHash=$(BUILD_HASH)"
BUILDER_GOOS_GOARCH=$(shell $(GO) env GOOS)_$(shell $(GO) env GOARCH)
PACKAGESLISTS=$(shell $(GO) list ./...)
TESTFLAGS ?= -short
PACKAGESLISTS_COMMA=$(shell echo $(PACKAGESLISTS) | tr ' ' ',')
ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

test:
	@echo Testing ...
	$(GO) test ./...

build: test
	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(GO) build -i $(GOFLAGS) $(GO_LINKER_FLAGS) ./dappswin.go

deps:
	@echo make deps ... 
	@go get github.com/gravityblast/fresh
	@go get github.com/golang/dep/cmd/dep
	@dep ensure -v

run:
	@echo Running dappswin
	fresh
	
pack: test build
	@echo tar builded
	tar zcvf dappswin.tar.gz dappswin dappswin.toml scripts install.sh

govet: ## Runs govet against all packages.
	@echo Running GOVET
	$(GO) vet $(GOFLAGS) $(PACKAGESLISTS) || exit 1

gofmt: ## Runs gofmt against all packages.
	@echo Running GOFMT
	@echo $(PACKAGESLISTS)
	@for package in $(PACKAGESLISTS);do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \

