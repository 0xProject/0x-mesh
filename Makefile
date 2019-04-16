.PHONY: deps
deps:
	dep ensure
	yarn install


.PHONY: test-all
test-all: test-go test-wasm


.PHONY: test-go
test-go:
	go test ./... -v -race


.PHONY: test-wasm
test-wasm:
	export ZEROEX_MESH_ROOT_DIR=$$(pwd); GOOS=js GOARCH=wasm go test -exec="$$ZEROEX_MESH_ROOT_DIR/test-wasm/go_js_wasm_exec" ./... -v
