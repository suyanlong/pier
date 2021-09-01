
SHELL := /bin/bash
CURRENT_PATH = $(shell pwd)
APP_NAME = sidercar

# build with verison infos
VERSION_DIR = github.com/link33/${APP_NAME}
BUILD_DATE = $(shell date +%FT%T)
GIT_COMMIT = $(shell git log --pretty=format:'%h' -n 1)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
ifeq (${GIT_BRANCH},HEAD)
  APP_VERSION = $(shell git describe --tags HEAD)
else
  APP_VERSION = dev
endif

GOLDFLAGS += -X "${VERSION_DIR}.BuildDate=${BUILD_DATE}"
GOLDFLAGS += -X "${VERSION_DIR}.CurrentCommit=${GIT_COMMIT}"
GOLDFLAGS += -X "${VERSION_DIR}.CurrentBranch=${GIT_BRANCH}"
GOLDFLAGS += -X "${VERSION_DIR}.CurrentVersion=${APP_VERSION}"

STATIC_LDFLAGS += ${GOLDFLAGS}
STATIC_LDFLAGS += -linkmode external -extldflags -static

GO = GO111MODULE=on go
TEST_PKGS := $(shell $(GO) list ./... | grep -v 'cmd' | grep -v 'mock_*' | grep -v 'proto' | grep -v 'imports' | grep -v 'internal/app' | grep -v 'api')

RED=\033[0;31m
GREEN=\033[0;32m
BLUE=\033[0;34m
NC=\033[0m

.PHONY: test

help: Makefile
	@echo "Choose a command run:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'

## make test: Run go unittest
test:
	go generate ./...
	@$(GO) test ${TEST_PKGS} -count=1

## make test-coverage: Test project with cover
test-coverage:
	go generate ./...
	@go test -short -coverprofile cover.out -covermode=atomic ${TEST_PKGS}
	@cat cover.out >> coverage.txt

packr:
	cd internal/repo && packr

prepare:
	cd scripts && bash prepare.sh

## make install: Go install the project (hpc)
install: packr
	rm -f imports/imports.go
	$(GO) install -ldflags '${GOLDFLAGS}' ./cmd/${APP_NAME}
	@printf "${GREEN}Build sidercar successfully${NC}\n"

build: packr
	@mkdir -p bin
	rm -f imports/imports.go
	$(GO) build -ldflags '${GOLDFLAGS}' ./cmd/${APP_NAME}
	@mv ./sidercar bin
	@printf "${GREEN}Build sidercar successfully!${NC}\n"

installent: packr
	cp imports/imports.go.template imports/imports.go
	@sed "s?)?$(MODS)@)?" go.mod  | tr '@' '\n' > goent.mod
	$(GO) install -tags ent -ldflags '${GOLDFLAGS}' -modfile goent.mod ./cmd/${APP_NAME}

buildent: packr
	@mkdir -p bin
	cp imports/imports.go.template imports/imports.go
	@sed "s?)?$(MODS)@)?" go.mod  | tr '@' '\n' > goent.mod
	$(GO) build -tags ent -ldflags '${GOLDFLAGS}' -modfile goent.mod ./cmd/${APP_NAME}
	@mv ./sidercar bin
	@printf "${GREEN}Build sidercar ent successfully!${NC}\n"

mod:
	sed "s?)?$(MODS)\n)?" go.mod

docker-build: packr
	$(GO) install -ldflags '${STATIC_LDFLAGS}' ./cmd/${APP_NAME}
	@echo "Build sidercar successfully"

## make build-linux: Go build linux executable file
build-linux:
	cd scripts && bash cross_compile.sh linux-amd64 ${CURRENT_PATH}

## make release: Build release before push
release-binary:
	@cd scripts && bash release_binary.sh

## make linter: Run golanci-lint
linter:
	golangci-lint run
	golangci-lint run -E goimports -E bodyclose --skip-dirs-use-default

fmt:
	go fmt ./...



all: pb grpc

pb:
	@cd model/pb && protoc -I=. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	--gogofaster_out=:. \
	ibtp.proto ibtpx.proto  basic.proto message.proto


#pb:
#	cd model/pb && protoc -I=. \
#	-I${GOPATH}/src \
#	--gogofaster_out=:. \
#	block.proto ibtp.proto ibtpx.proto network.proto receipt.proto bxh_transaction.proto chain.proto arg.proto interchain_meta.proto plugin.proto vp_info.proto basic.proto

grpc:
	cd model/pb && protoc -I=. \
	-I=${GOPATH}/src \
	-I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	--grpc-gateway_out=logtostderr=true:. \
	--swagger_out=logtostderr=true:. \
	--gogofaster_out=plugins=grpc:. \
	broker.proto plugin.proto


clean:
	rm pb/*.pb.go
	rm pb/*.json
	rm pb/*.gw.go

.PHONY: pb
