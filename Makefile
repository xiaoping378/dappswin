.PHONY: build  package start-docker run deps
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

build:
	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(GO) build -i $(GOFLAGS) $(GO_LINKER_FLAGS) ./dappswin.go

package:
	@echo Packaging dappswin
	rm -rf $(DIST_PATH)
	mkdir -p $(RELEASE_PATH)/bin
	mkdir -p $(RELEASE_PATH)/config
	cp $(GOPATH)/bin/dappswin $(RELEASE_PATH)/bin
	cp conf/*.json $(RELEASE_PATH)/config
	cp ./*.toml $(RELEASE_PATH)/
	tar -C $(DIST_PATH) -czf $(RELEASE_PATH)-$(BUILDER_GOOS_GOARCH).tar.gz release

deps:
	@echo make deps ... 
	@go get github.com/gravityblast/fresh
	@go get github.com/golang/dep/cmd/dep
	@dep ensure -v

run:
	@echo Running dappswin
	fresh
	



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

test: clean-docker start-docker test-server

test-server:
	docker run --rm  --name ginkgo-test --network my-net -v $(GOPATH)/src/dappswin:/go/src/dappswin  golang:latest /bin/sh -c  "/go/src/dappswin/ginkgo -r -trace -cover  -coverprofile=coverprofile.txt -outputdir=/go/src/dappswin   /go/src/dappswin"

start-docker: ## Starts the docker containers for local development.
	@echo Starting docker containers

	@if [ $(shell docker ps -a | grep -ci critic-mysql) -eq 0 ]; then \
		echo starting dappswin-mysql; \
		docker run --network my-net --name dappswin-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 \
		-e MYSQL_USER=dappswin -e MYSQL_PASSWORD=123456 -e MYSQL_DATABASE=dappswin_test -d mysql:5.7 > /dev/null; \
	elif [ $(shell docker ps | grep -ci critic-mysql) -eq 0 ]; then \
		echo restarting critic-mysql; \
		docker start critic-mysql > /dev/null; \
	fi

stop-docker: ## Stops the docker containers for local development.
	@echo Stopping docker containers

	@if [ $(shell docker ps -a | grep -ci dappswin-mysql) -eq 1 ]; then \
		echo stopping critic-mysql; \
		docker stop critic-mysql > /dev/null; \
	fi

clean-docker: ## Deletes the docker containers for local development.
	@echo Removing docker containers

	@if [ $(shell docker ps -a | grep -ci dappswin-mysql) -eq 1 ]; then \
		echo removing critic-mysql; \
		docker stop critic-mysql > /dev/null; \
		docker rm -v critic-mysql > /dev/null; \
	fi


