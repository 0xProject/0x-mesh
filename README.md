# 0x-mesh

A peer-to-peer network for sharing orders

:no_entry: WARNING: This project is still under construction and is not ready for use :no_entry:

What's currently implemented? In it's current state, 0x Mesh is a P2P node that stores, validates, and shares 0x orders.- It uses FloodSub as a pubsub mechanism.

- It has a WebSocket endpoint for submitting 0x orders to the network.
- It performs 0x order validation and orderbook pruning.
- All relevant state is persisted to disk so that the node can trivially restart from where it left off after a crash.

What is still being built for the Beta release:

- Peer discovery (currently you must manually add new peers).
- WeijieSub (a more incentive-compatible and more efficient pubsub mechanism).
- A way to subscribe to order events via the WebSocket API.
- Tracking ETH reserves for orders stored (we have already implemented efficient ETH balance updating; it's just not currently used when storing orders).
- Full browser support (some components are already browser-ready, but others are not).

Open an issue and introduce yourself if you want to help build Mesh!

## Development

If you are working on 0x-mesh, the root directory for the project must be at
**\$GOPATH/src/github.com/0xProject/0x-mesh**. 0x Mesh uses Dep for dependency
management and does not support Go modules.

### Prerequisites

- [GNU Make](https://www.gnu.org/software/make/) If you are using a Unix-like OS, you probably already have this.
- [Go version >= 1.12](https://golang.org/dl/) (or use [the version manager called "g"](https://github.com/stefanmaric/g))
- [Dep package manager](https://golang.github.io/dep/docs/installation.html)
- [Node.js version >= 10](https://nodejs.org/en/download/) (or use the [nvm version manager](https://github.com/creationix/nvm))
- [Yarn package manager](https://yarnpkg.com/en/)
- [golangci-lint version 1.16.0](https://github.com/golangci/golangci-lint#install)

### Installing dependencies

```
make deps
```

### Running tests

Some of the tests depend on having a test Ethereum node running. Before running the tests, make sure you have [Docker](https://docs.docker.com/install/) installed locally and start [0xorg/mesh-ganache-cli](https://cloud.docker.com/u/0xorg/repository/docker/0xorg/mesh-ganache-cli):

```
docker pull 0xorg/mesh-ganache-cli
docker run -ti -p 8545:8545 0xorg/mesh-ganache-cli
```

Run tests in a vanilla Go environment:

```
make test-go
```

Run tests in a Node.js/WebAssembly environment:

```
make test-wasm
```

Run tests in both environments:

```
make test-all
```

### Running the linter

```
make lint
```

### Managing Dependencies

See https://golang.github.io/dep/docs/daily-dep.html.

### Editor Configuration

#### Visual Studio Code

For VS Code, the following editor configuration is recommended:

```javascript
{
  // ...

  "editor.formatOnSave": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.vetOnSave": "off"

  // ...
}
```
