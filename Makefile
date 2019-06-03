.PHONY: deps
deps:
	dep ensure
	yarn install


# Installs dependencies without updating Gopkg.lock or yarn.lock
.PHONY: deps-no-lockfile
deps-no-lockfile:
	dep ensure --vendor-only
	yarn install --frozen-lockfile


.PHONY: test-all
test-all: test-go test-wasm


.PHONY: test-go
test-go:
	go test ./... -race


.PHONY: test-wasm
test-wasm:
	export ZEROEX_MESH_ROOT_DIR=$$(pwd); GOOS=js GOARCH=wasm go test -exec="$$ZEROEX_MESH_ROOT_DIR/test-wasm/go_js_wasm_exec" ./...


.PHONY: lint
lint:
	golangci-lint run


.PHONY: mesh
mesh:
	go install ./cmd/mesh
