.PHONY: deps
deps: deps-ts wasmbrowsertest gqlgen


.PHONY: deps-ts
deps-ts:
	yarn install


# gobin allows us to install specific versions of binary tools written in Go.
.PHONY: gobin
gobin:
	GO111MODULE=off go get -u github.com/myitcv/gobin


# gqlgen is a tool for embedding files so that they are included in binaries.
# This installs the CLI for go-bindata.
.PHONY: gqlgen
gqlgen: gobin
	gobin github.com/99designs/gqlgen@v0.11.3


# wasmbrowsertest is required for running WebAssembly tests in the browser.
.PHONY: wasmbrowsertest
wasmbrowsertest: gobin
	gobin github.com/0xProject/wasmbrowsertest@mesh-fork


# Installs dependencies without updating Gopkg.lock or yarn.lock
.PHONY: deps-no-lockfile
deps-no-lockfile: deps-ts-no-lockfile wasmbrowsertest gqlgen


.PHONY: deps-ts-no-lockfile
deps-ts-no-lockfile:
	yarn install --frozen-lockfile


# Provides a pre-commit convenience command that runs all of the tests and the linters
.PHONY: check
check: test-all lint


.PHONY: test-all
test-all: test-go test-ts


.PHONY: test-go
test-go: generate test-go-parallel test-go-serial


.PHONY: test-go-parallel
test-go-parallel:
	go test ./... -race -timeout 30s


.PHONY: test-key-value-stores
test-key-value-stores: test-key-value-stores-go test-key-value-stores-wasm


.PHONY: test-key-value-stores-go
test-key-value-stores-go:
	ENABLE_KEY_VALUE_TESTS=true go test ./db


.PHONY: test-key-value-stores-wasm
test-key-value-stores-wasm:
	WASM_INIT_FILE="$$(pwd)/packages/mesh-browser-shim/dist/browser_shim.js" GOOS=js GOARCH=wasm ENABLE_KEY_VALUE_TESTS=true go test ./db -timeout 30m -tags=browser -exec="$$GOPATH/bin/wasmbrowsertest"


.PHONY: test-go-serial
test-go-serial:
	go test ./zeroex/ordervalidator ./zeroex/orderwatch ./core -race -timeout 300s -p=1 --serial

.PHONY: test-browser-integration
test-browser-integration: test-browser-legacy-integration test-browser-graphql-integration


.PHONY: test-browser-legacy-integration
test-browser-legacy-integration:
	go test ./integration-tests -timeout 60s --enable-browser-legacy-integration-tests -run BrowserLegacyIntegration


.PHONY: test-browser-graphql-integration
test-browser-graphql-integration:
	go test ./integration-tests -timeout 60s --enable-browser-graphql-integration-tests -run BrowserGraphQLIntegration


.PHONY: test-browser-conversion
test-browser-conversion:
	go test ./packages/mesh-browser/go/conversion-test -timeout 185s --enable-browser-conversion-tests -run BrowserConversions


.PHONY: test-wasm-browser
test-wasm-browser:
	WASM_INIT_FILE="$$(pwd)/packages/mesh-browser-shim/dist/browser_shim.js" GOOS=js GOARCH=wasm go test -tags=browser -exec="$$GOPATH/bin/wasmbrowsertest" ./...


.PHONY: test-ts
test-ts:
	yarn test


.PHONY: lint
lint: lint-go lint-ts lint-prettier


.PHONY: lint-go
lint-go:
	golangci-lint run --timeout 2m 


.PHONY: lint-ts
lint-ts:
	yarn lint

.PHONY: lint-prettier
lint-prettier:
	yarn prettier:ci


.PHONY: generate
generate:
	go generate ./...


.PHONY: mesh
mesh: generate
	go install ./cmd/mesh


.PHONY: mesh-keygen
mesh-keygen:
	go install ./cmd/mesh-keygen


.PHONY: mesh-bootstrap
mesh-bootstrap:
	go install ./cmd/mesh-bootstrap


.PHONY: db-integrity-check
db-integrity-check:
	go install ./cmd/db-integrity-check


.PHONY: cut-release
cut-release: generate
	go run ./cmd/cut-release/main.go


.PHONY: all
all: mesh mesh-keygen mesh-bootstrap db-integrity-check


# Docker images


.PHONY: docker-mesh
docker-mesh: generate
	docker build . -t 0xorg/mesh -f ./dockerfiles/mesh/Dockerfile


.PHONY: docker-mesh-bootstrap
docker-mesh-bootstrap:
	docker build . -t 0xorg/mesh-bootstrap -f ./dockerfiles/mesh-bootstrap/Dockerfile


.PHONY: docker-mesh-fluent-bit
docker-mesh-fluent-bit:
	docker build ./dockerfiles/mesh-fluent-bit -t 0xorg/mesh-fluent-bit -f ./dockerfiles/mesh-fluent-bit/Dockerfile


.PHONY: docker-mesh-bridge
docker-mesh-bridge: generate
	@echo 'WARN: mesh-bridge is currently disabled since it has not been updated to use the new GraphQL API' 
	# docker build . -t 0xorg/mesh-bridge -f ./dockerfiles/mesh-bridge/Dockerfile
