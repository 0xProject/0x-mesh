[![Version](https://img.shields.io/badge/version-3.0.1--beta-orange.svg)](https://github.com/0xProject/0x-mesh/releases)

Welcome to the [0x Mesh](https://github.com/0xProject/0x-mesh) Development
Guide! This guide will help you get started with contributing to the 0x Mesh
codebase.

## Directory Location

If you are working on 0x-mesh, the root directory for the project must be at
**\$GOPATH/src/github.com/0xProject/0x-mesh**. 0x Mesh uses [Dep](https://golang.github.io/dep/docs/installation.html) for dependency
management and does not support Go modules.

## Prerequisites

-   [GNU Make](https://www.gnu.org/software/make/) If you are using a Unix-like OS, you probably already have this.
-   [Go version >= 1.12](https://golang.org/dl/) (or use [the version manager called "g"](https://github.com/stefanmaric/g))
-   [Dep package manager](https://golang.github.io/dep/docs/installation.html)
-   [Node.js version >=11](https://nodejs.org/en/download/) (or use the [nvm version manager](https://github.com/creationix/nvm))
-   [Yarn package manager](https://yarnpkg.com/en/)
-   [golangci-lint version 1.16.0](https://github.com/golangci/golangci-lint#install)

## Installing Dependencies

```
make deps
```

## Running Tests

Some of the tests depend on having a test Ethereum node running. Before running
the tests, make sure you have [Docker](https://docs.docker.com/install/)
installed locally and start
[0xorg/mesh-ganache-cli](https://hub.docker.com/r/0xorg/mesh-ganache-cli):

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

## Running the linter

```
make lint
```

## Managing Dependencies

See https://golang.github.io/dep/docs/daily-dep.html.

## Editor Configuration

### Visual Studio Code

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
