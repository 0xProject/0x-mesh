.PHONY: deps
deps: deps-go deps-ts wasmbrowsertest


.PHONY: deps-go
deps-go:
	dep ensure


.PHONY: deps-ts
deps-ts:
	yarn install
	cd rpc/clients/typescript && yarn install
	cd browser/ && yarn install


# gobin allows us to install specific versions of binary tools written in Go.
.PHONY: gobin
gobin:
	GO111MODULE=off go get -u github.com/myitcv/gobin


# wasmbrowsertest is required for running WebAssembly tests in the browser.
.PHONY: wasmbrowsertest
wasmbrowsertest: gobin
	gobin github.com/agnivade/wasmbrowsertest@v0.3.0


# Installs dependencies without updating Gopkg.lock or yarn.lock
.PHONY: deps-no-lockfile
deps-no-lockfile: deps-go-no-lockfile deps-ts-no-lockfile wasmbrowsertest


.PHONY: deps-go-no-lockfile
deps-go-no-lockfile:
	dep ensure --vendor-only


.PHONY: deps-ts-no-lockfile
deps-ts-no-lockfile:
	yarn install --frozen-lockfile
	cd rpc/clients/typescript && yarn install --frozen-lockfile
	cd browser/ && yarn install --frozen-lockfile


.PHONY: test-all
test-all: test-go test-wasm-node test-wasm-browser


.PHONY: test-go
test-go: test-go-parallel test-go-serial

.PHONY: test-go-parallel
test-go-parallel:
	go test ./... -race -timeout 30s

.PHONY: test-go-serial
test-go-serial:
	go test ./zeroex/ordervalidator ./zeroex/orderwatch -race -timeout 30s -p=1 --serial


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
lint: lint-go lint-ts


.PHONY: lint-go
lint-go:
	golangci-lint run


.PHONY: lint-ts
lint-ts:
	cd rpc/clients/typescript && yarn lint
	cd browser/ && yarn lint


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


# Docker images


.PHONY: docker-mesh
docker-mesh:
	docker build . -t 0xorg/mesh -f ./dockerfiles/mesh/Dockerfile


.PHONY: docker-mesh-bootstrap
docker-mesh-bootstrap:
	docker build . -t 0xorg/mesh-bootstrap -f ./dockerfiles/mesh-bootstrap/Dockerfile


.PHONY: docker-mesh-fluent-bit
docker-mesh-fluent-bit:
	docker build ./dockerfiles/mesh-fluent-bit -t 0xorg/mesh-fluent-bit -f ./dockerfiles/mesh-fluent-bit/Dockerfile
