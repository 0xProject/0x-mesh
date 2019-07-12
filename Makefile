.PHONY: deps
deps: deps-go deps-js


.PHONY: deps-go
deps-go:
	dep ensure


.PHONY: deps-js
deps-js:
	yarn install


# Installs dependencies without updating Gopkg.lock or yarn.lock
.PHONY: deps-no-lockfile
deps-no-lockfile: deps-go-no-lockfile deps-js-no-lockfile


.PHONY: deps-go-no-lockfile
deps-go-no-lockfile:
	dep ensure --vendor-only


.PHONY: deps-js-no-lockfile
deps-js-no-lockfile:
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


.PHONY: mesh-keygen
mesh-keygen:
	go install ./cmd/mesh-keygen


.PHONY: mesh-bootstrap
mesh-bootstrap:
	go install ./cmd/mesh-bootstrap


.PHONY: db-integrity-check
db-integrity-check:
	go install ./cmd/db-integrity-check


.PHONY: all
all: mesh mesh-keygen mesh-bootstrap db-integrity-check
