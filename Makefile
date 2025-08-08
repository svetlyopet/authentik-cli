SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

BINARY_NAME=authentik-cli

OUT_DIR=out
DIST_DIR=dist

PLATFORMS=linux darwin
ARCHS=amd64 arm64


out:
	@mkdir -p $(OUT_DIR)

git-hooks: ## Setup githooks
	@git config --local core.hooksPath .githooks/

download: ## Downloads dependencies
	@go mod download

tidy: ## Cleans up go.mod and go.sum
	@go mod tidy

fmt: ## Formats all code with go fmt and goimports
	@go fmt ./...
	@go run golang.org/x/tools/cmd/goimports@latest -w .

.PHONY: mocks
mocks: ## Generates mocks
	@go run go.uber.org/mock/mockgen@latest -destination=mocks/ak/ak.go -package=mock_ak -source=internal/ak/ak.go

test-build: ## Tests whether the code compiles
	@go build -o /dev/null ./...

build: $(OUT_DIR)/bin ## Builds the binary

.PHONY: $(OUT_DIR)/bin
$(OUT_DIR)/bin:
	@mkdir -p "$(@)" && go build -ldflags="-w -s" -o "$(@)" ./...

dist:
	@mkdir -p $(DIST_DIR)

release: dist ## Build release binaries
	@if [ -z "$(VERSION)" ]; then \
		echo "Release version not set."; \
		echo "Set by passing VERSION environment variable to the make release command."; \
		echo "Exiting..."; \
		exit 1; \
	fi;
	@cd $(DIST_DIR); \
	for GOOS in $(PLATFORMS); do \
		for GOARCH in $(ARCHS); do \
			RELEASE_BIN=$(BINARY_NAME)-$$GOOS-$$GOARCH-$$VERSION; \
			RELEASE_TAR_GZ=$$RELEASE_BIN.tar.gz; \
			echo "Building $$RELEASE_BIN..."; \
			GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags="-s -w" -o $$RELEASE_BIN ../.; \
			echo "Compressing into $$RELEASE_TAR_GZ file..."; \
			tar -czf $$RELEASE_TAR_GZ $$RELEASE_BIN; \
			rm -f $$RELEASE_BIN; \
			echo "Calculating SHA256 checksum..."; \
			sha256sum $$RELEASE_TAR_GZ >> $(BINARY_NAME)-checksums.txt; \
		done \
	done

lint: fmt download ## Lints all code with golangci-lint
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run

lint-reports: $(OUT_DIR)/lint.xml

.PHONY: $(OUT_DIR)/lint.xml
$(OUT_DIR)/lint.xml: out download
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --out-format checkstyle | tee "$(@)"

govulncheck: ## Vulnerability detection using govulncheck
	@go run golang.org/x/vuln/cmd/govulncheck ./...

test: ## Runs unit tests
	@go test -v -covermode=atomic ./...

e2e-local: ## Run e2e tests locally
	@$(PWD)/ci/integration_tests/ak_bootstrap.sh create
	@$(PWD)/ci/integration_tests/ak_run_e2e.sh

e2e-ci: ## Run e2e tests in CI
	@AK_BOOTSTRAP_CI=true AK_BOOTSTRAP_WAIT=90 $(PWD)/ci/integration_tests/ak_bootstrap.sh create
	@$(PWD)/ci/integration_tests/ak_run_e2e.sh

e2e-cleanup: # Cleanup e2e environement
	@$(PWD)/ci/integration_tests/ak_bootstrap.sh destroy

coverage: $(OUT_DIR)/report.json ## Displays coverage per func on cli
	@go tool cover -func=$(OUT_DIR)/cover.out

html-coverage: $(OUT_DIR)/report.json ## Displays the coverage results in the browser
	@go tool cover -html=$(OUT_DIR)/cover.out

test-reports: $(OUT_DIR)/report.json

.PHONY: $(OUT_DIR)/report.json
$(OUT_DIR)/report.json: out
	@go test -count 1 ./... -coverprofile=$(OUT_DIR)/cover.out --json | tee "$(@)"

clean: ## Cleans up output and release files
	@rm -rf $(DIST_DIR) $(OUT_DIR)

define make-go-dependency
  # target template for go tools, can be referenced e.g. via /bin/<tool>
  bin/$(notdir $1):
	GOBIN=$(PWD)/bin go install $1
endef

# this creates a target for each go dependency to be referenced in other targets
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))
ci: lint-reports test-reports govulncheck ## Executes vulnerability scan, lint, test and generates reports

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''