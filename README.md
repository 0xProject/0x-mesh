[![Chat with us on Discord](https://img.shields.io/badge/chat-Discord-blueViolet.svg)](https://discord.gg/HF7fHwk)
[![HackMD](https://img.shields.io/badge/HackMD-Docs-orange.svg)](https://hackmd.io/s/By4pSjATE)
[![GoDoc](https://godoc.org/github.com/0xProject/0x-mesh/core?status.svg)](https://godoc.org/github.com/0xProject/0x-mesh/core)
[![Circle CI](https://img.shields.io/circleci/project/0xProject/0x-mesh/master.svg)](https://circleci.com/gh/0xProject/0x-mesh/tree/master)


# 0x-mesh

A peer-to-peer network for sharing orders

:no_entry: WARNING: This project is still under construction and is not ready for use :no_entry:

## Development

If you are working on 0x-mesh, the root directory for the project must be at
**\$GOPATH/src/github.com/0xProject/0x-mesh**. 0x Mesh uses [Dep](https://golang.github.io/dep/docs/installation.html) for dependency
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

### Installing the standalone Mesh node

To install a standalone version of 0x Mesh that can be run from the command
line, run:

```
make mesh
```

If you have `GOPATH/bin` in your `PATH`, you can now run 0x Mesh directly with:

```
mesh
```

### Running tests

Some of the tests depend on having a test Ethereum node running. Before running
the tests, make sure you have [Docker](https://docs.docker.com/install/)
installed locally and start
[0xorg/mesh-ganache-cli](https://cloud.docker.com/u/0xorg/repository/docker/0xorg/mesh-ganache-cli):

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
