.PHONY: deps
deps: deps-go deps-js wasmbrowsertest


.PHONY: deps-go
deps-go:
	dep ensure


.PHONY: deps-js
deps-js:
	yarn install


# wasmbrowsertest is required for running WebAssembly tests in the browser.
.PHONY: wasmbrowsertest
wasmbrowsertest:
	go get -u github.com/agnivade/wasmbrowsertest


# Installs dependencies without updating Gopkg.lock or yarn.lock
.PHONY: deps-no-lockfile
deps-no-lockfile: deps-go-no-lockfile deps-js-no-lockfile wasmbrowsertest


.PHONY: deps-go-no-lockfile
deps-go-no-lockfile:
	dep ensure --vendor-only


.PHONY: deps-js-no-lockfile
deps-js-no-lockfile:
	yarn install --frozen-lockfile


.PHONY: test-all
test-all: test-go test-wasm-node test-wasm-browser


.PHONY: test-go
test-go:
	go test ./... -race -timeout 30s


.PHONY: test-integration
test-integration:
	go test ./integration-tests -timeout 185s --integration


.PHONY: test-wasm-node
test-wasm-node:
	export ZEROEX_MESH_ROOT_DIR=$$(pwd); GOOS=js GOARCH=wasm go test -exec="$$ZEROEX_MESH_ROOT_DIR/test-wasm/go_js_wasm_exec" ./...


.PHONY: test-wasm-browser
test-wasm-browser:
	GOOS=js GOARCH=wasm go test -tags=browser -exec="$$GOPATH/bin/wasmbrowsertest" ./...


.PHONY: lint
lint:
	golangci-lint run


.PHONY: mesh
mesh:
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
cut-release:
	go run ./cmd/cut-release/main.go


.PHONY: all
all: mesh mesh-keygen mesh-bootstrap db-integrity-check
