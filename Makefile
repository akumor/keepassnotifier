export REPOSITORY=keepassnotifier
include boilerplate/docker_build/Makefile
include boilerplate/golang_test_targets/Makefile

GIT_VERSION := $(shell git describe --always --tags)
GIT_HASH := $(shell git rev-parse --short HEAD)
TIMESTAMP := $(shell date '+%Y-%m-%d')

LD_FLAGS="-s -w -X $(PACKAGE)/version.Version=$(GIT_VERSION) -X $(PACKAGE)/version.Build=$(GIT_HASH) -X $(PACKAGE)/version.BuildTime=$(TIMESTAMP)"

.PHONY: compile
compile:
	go build -o keepassnotifier -ldflags=$(LD_FLAGS) ./cmd/ && mv ./keepassnotifier ${GOPATH}/bin

.PHONY: linux_compile
linux_compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o /artifacts/keepassnotifier -ldflags=$(LD_FLAGS) ./cmd/

.PHONY: entries
server:
	go run cmd/main.go entries  --config ./test_config.yaml --database ./example.kbdx --credentials ./test_credentials.json

all: compile
