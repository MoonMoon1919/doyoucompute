# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVERSION=`cat go.mod | grep 'go\s\d.' | cut -d ' ' -f2`

GOENVCMD=goenv

# Check if required tools are installed
.PHONE: check-goenv
check-goenv:
	@which $(GOENVCMD) >/dev/null 2>&1 || \
		(echo "ERROR: goenv is not installed or not in PATH" && exit 1)

.PHONY: check-tools
check-tools:
	@which $(GOCMD) >/dev/null 2>&1 || \
		(echo "ERROR: Go is not installed or not in PATH" && exit 1)

# Format all go files
.PHONY: fmt
fmt: check-tools
	@$(GOFMT) ./...

# Run go vet
.PHONY: vet
vet: check-tools
	@$(GOVET) ./...

# Download dependencies
.PHONY: deps
deps: check-tools
	@$(GOMOD) download
	@$(GOMOD) verify

# Run tests
.PHONY: test/unit
test/unit: check-tools
	@$(GOTEST) -v ./...

.PHONY: test/unit/cover
test/unit/cover: check-tools
	@$(GOTEST) -v -cover ./...


.PHONY: init-shell
init-shell: check-goenv
	@$(GOENVCMD) local $(GOVERSION)

# docs
.PHONY: docs/readme
docs/readme:
	@$(GOCMD) run docs/main.go render --doc-name 'DOYOUCOMPUTE' --path README.md

.PHONY: validate/readme
validate/readme:
	@$(GOCMD) run docs/main.go compare --doc-name 'DOYOUCOMPUTE' --path README.md

.PHONY: docs/contrib
docs/contrib:
	@$(GOCMD) run docs/main.go render --doc-name 'Contributing' --path CONTRIBUTING.md

.PHONY: validate/contrib
validate/contrib:
	@$(GOCMD) run docs/main.go compare --doc-name 'Contributing' --path CONTRIBUTING.md

.PHONY: template/pullrequest
template/pullrequest:
	@$(GOCMD) run docs/main.go render --doc-name 'Pull Request' --path ./.github/PULL_REQUEST_TEMPLATE.md

.PHONY: validate/pullrequest
validate/pullrequest:
	@$(GOCMD) run docs/main.go compare --doc-name 'Pull Request' --path .github/PULL_REQUEST_TEMPLATE.md

.PHONY: template/bugreport
template/bugreport:
	@$(GOCMD) run docs/main.go render --doc-name 'Bug Report' --path ./.github/ISSUE_TEMPLATE/bug_report.md

.PHONY: validate/bugreport
validate/bugreport:
	@$(GOCMD) run docs/main.go compare --doc-name 'Bug Report' --path .github/ISSUE_TEMPLATE/bug_report.md

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  clean             - Removes build artifacts"
	@echo "  deps              - Downloads and verify dependencies"
	@echo "  fmt               - Formats Go source files"
	@echo "  help              - Shows this help message"
	@echo "  test/unit         - Runs unit tests"
	@echo "  vet               - Runs go vet"
	@echo "  init-shell        - Sets goversion using goenv"

# Default target
.DEFAULT_GOAL := help
