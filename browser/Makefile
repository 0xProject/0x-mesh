.PHONY: test
test: 
	export ZEROEX_MESH_ROOT_DIR=$$(dirname $$(pwd)); GOOS=js GOARCH=wasm go test -exec=$$ZEROEX_MESH_ROOT_DIR/test-wasm/go_js_wasm_exec ./go/
