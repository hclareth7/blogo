.PHONY: all build dev test lint fmt vet clean fetch-content css docker help

BINARY_NAME := blogo
MODULE := github.com/hclareth7/blogo
GO := go
GOFLAGS := -trimpath
LDFLAGS := -s -w

CONTENT_URL ?= https://raw.githubusercontent.com/karanpratapsingh/system-design/main/README.md
CONTENT_DIR ?= ./content

all: lint test build

build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/blogo

dev:
	@command -v air > /dev/null 2>&1 || { echo "Install air: go install github.com/air-verse/air@latest"; exit 1; }
	air

run: build
	./bin/$(BINARY_NAME)

test:
	$(GO) test -race -count=1 ./...

test-cover:
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

lint:
	@command -v golangci-lint > /dev/null 2>&1 || { echo "Install golangci-lint: https://golangci-lint.run/welcome/install/"; exit 1; }
	golangci-lint run ./...

fmt:
	@command -v gofumpt > /dev/null 2>&1 || { echo "Install gofumpt: go install mvdan.cc/gofumpt@latest"; exit 1; }
	gofumpt -w .

vet:
	$(GO) vet ./...

fetch-content:
	@mkdir -p $(CONTENT_DIR)
	curl -sSL -o $(CONTENT_DIR)/README.md $(CONTENT_URL)
	@echo "Content fetched to $(CONTENT_DIR)/README.md"

css:
	@command -v tailwindcss > /dev/null 2>&1 || { echo "Install Tailwind CSS standalone CLI: https://tailwindcss.com/blog/standalone-cli"; exit 1; }
	tailwindcss -i web/assets/input.css -o web/static/css/styles.css --minify

css-watch:
	tailwindcss -i web/assets/input.css -o web/static/css/styles.css --watch

docker:
	docker build -t $(BINARY_NAME):latest -f deploy/docker/Dockerfile .

clean:
	rm -rf bin/ coverage.out coverage.html

help:
	@echo "Available targets:"
	@echo "  all            - Lint, test, and build"
	@echo "  build          - Build static binary"
	@echo "  dev            - Run with hot reload (requires air)"
	@echo "  run            - Build and run"
	@echo "  test           - Run tests with race detector"
	@echo "  test-cover     - Run tests with coverage report"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code with gofumpt"
	@echo "  vet            - Run go vet"
	@echo "  fetch-content  - Download content from source repo"
	@echo "  css            - Compile Tailwind CSS"
	@echo "  css-watch      - Watch and compile Tailwind CSS"
	@echo "  docker         - Build container image"
	@echo "  clean          - Remove build artifacts"
