BINARY_NAME=az-health-exporter
DOCKER_IMAGE=mobymatze/az-health-exporter
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
RELEASE_NAME?=az-health-monitor

PACKAGE_NAME=github.com/matzefriedrich/az-health-exporter
LDFLAGS=-s -w \
        -X $(PACKAGE_NAME)/internal.Version=$(VERSION) \
        -X $(PACKAGE_NAME)/internal.CommitSha=$(COMMIT) \
        -X $(PACKAGE_NAME)/internal.ReleaseDate=$(DATE) \
        -X $(PACKAGE_NAME)/internal.ReleaseName=$(RELEASE_NAME)

.PHONY: all build clean fmt vet test docker-build

all: build

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/az-health-exporter

clean:
	rm -f $(BINARY_NAME)

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

docker-build:
	docker build \
		--build-arg APP_VERSION=$(VERSION) \
		--build-arg CI_COMMIT_SHORT_SHA=$(COMMIT) \
		--build-arg APP_RELEASE_DATE=$(DATE) \
		--build-arg APP_RELEASE=$(RELEASE_NAME) \
		-t $(DOCKER_IMAGE):$(VERSION) .
