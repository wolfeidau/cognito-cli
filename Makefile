GOLANGCI_VERSION = 1.24.0

ci: clean awsmocks lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w"

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

bin/mockgen:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/golang/mock/mockgen

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov

mocks: bin/mockgen
	@echo "--- build all the mocks"
	@bin/mockgen -destination=mocks/cognito_service.go -package=mocks github.com/wolfeidau/cognito-cli/pkg/cognito Service
.PHONY: mocks

awsmocks: bin/mockgen
	@echo "--- build all the awsmocks"
	@bin/mockgen -destination=awsmocks/cognito.go -package=awsmocks github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface CognitoIdentityProviderAPI
.PHONY: awsmocks

clean:
	@echo "--- clean all the things"
	@rm -rf dist
.PHONY: clean

lint: bin/golangci-lint
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint

test: bin/gcov2lcov
	@echo "--- test all the things"
	@go test -v -covermode=count -coverprofile=coverage.txt ./ ./pkg/... ./internal/...
	@bin/gcov2lcov -infile=coverage.txt -outfile=coverage.lcov
.PHONY: test
