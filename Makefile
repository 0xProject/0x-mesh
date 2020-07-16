.PHONY: deps
deps: deps-ts wasmbrowsertest go-bindata


.PHONY: deps-ts
deps-ts:
	yarn install


# gobin allows us to install specific versions of binary tools written in Go.
.PHONY: gobin
gobin:
	GO111MODULE=off go get -u github.com/myitcv/gobin


# go-bindata is a tool for embedding files so that they are included in binaries.
# This installs the CLI for go-bindata.
.PHONY: go-bindata
go-bindata: gobin
	gobin github.com/go-bindata/go-bindata/v3/go-bindata@v3.1.3


# wasmbrowsertest is required for running WebAssembly tests in the browser.
.PHONY: wasmbrowsertest
wasmbrowsertest: gobin
	gobin github.com/0xProject/wasmbrowsertest@mesh-fork


# Installs dependencies without updating Gopkg.lock or yarn.lock
.PHONY: deps-no-lockfile
deps-no-lockfile: deps-ts-no-lockfile wasmbrowsertest


.PHONY: deps-ts-no-lockfile
deps-ts-no-lockfile:
	yarn install --frozen-lockfile


.PHONY: test-all
test-all: test-go test-wasm-browser test-ts test-browser-conversion test-browser-integration


.PHONY: test-go
test-go: generate test-go-parallel test-go-serial


.PHONY: test-go-parallel
test-go-parallel:
	go test ./... -race -timeout 30s


.PHONY: test-go-serial
test-go-serial:
	go test ./zeroex/ordervalidator ./zeroex/orderwatch ./core -race -timeout 90s -p=1 --serial


.PHONY: test-browser-integration
test-browser-integration:
	go test ./integration-tests -timeout 185s --enable-browser-integration-tests -run BrowserIntegration


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
	docker build . -t 0xorg/mesh-bridge -f ./dockerfiles/mesh-bridge/Dockerfile
