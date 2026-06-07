VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS = -s -w \
	-X 'github.com/faridtriwicaksono/forgebe/internal/cli.version=$(VERSION)' \
	-X 'github.com/faridtriwicaksono/forgebe/internal/cli.commit=$(COMMIT)' \
	-X 'github.com/faridtriwicaksono/forgebe/internal/cli.buildDate=$(DATE)'

build:
	mkdir -p bin
	go build -ldflags="$(LDFLAGS)" -o bin/forgebe ./cmd/forgebe

.PHONY: build test clean install lint vet fmt

test:
	go test -count=1 ./...

vet:
	go vet ./...

fmt:
	gofmt -w ./cmd ./internal

clean:
	rm -rf bin/

install: build
	install -m 0755 bin/forgebe /usr/local/bin/forgebe

release: test vet build
	@echo "=== RELEASE BUILD OK ==="
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"
