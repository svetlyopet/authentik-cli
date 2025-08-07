SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

out:
	@mkdir -p out

git-hooks:
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

build: out/bin ## Builds the binary

.PHONY: out/bin
out/bin:
	@mkdir -pv "$(@)" && go build -ldflags="-w -s" -o "$(@)" ./...

lint: fmt download ## Lints all code with golangci-lint
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run

lint-reports: out/lint.xml

.PHONY: out/lint.xml
out/lint.xml: out download
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

coverage: out/report.json ## Displays coverage per func on cli
	@go tool cover -func=out/cover.out

html-coverage: out/report.json ## Displays the coverage results in the browser
	@go tool cover -html=out/cover.out

test-reports: out/report.json

.PHONY: out/report.json
out/report.json: out
	@go test -count 1 ./... -coverprofile=out/cover.out --json | tee "$(@)"

clean: ## Cleans up output files
	@rm -rf bin out

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