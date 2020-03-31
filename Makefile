GOLANGCI_VERSION = 1.23.8

ci: clean lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w"

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

bin/go-acc:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/ory/go-acc

clean:
	@echo "--- clean all the things"
	@rm -rf dist
.PHONY: clean

lint: bin/golangci-lint
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint

test: bin/go-acc
	@echo "--- test all the things"
	@bin/go-acc --ignore mocks ./... -- -short -v -failfast
.PHONY: test