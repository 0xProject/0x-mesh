# Development and contribution guide

## Directory Location

If you are working on 0x-mesh, the root directory for the project must be at **$GOPATH/src/github.com/0xProject/0x-mesh**. 0x Mesh uses [Dep](https://golang.github.io/dep/docs/installation.html) for dependency management and does not support Go modules.

## Cloning the Repository and Opening PRs

0x Mesh uses two main branches:

1. The `development` branch contains the latest \(possibly unreleased\) changes

    and is not guaranteed to be stable.

2. The `master` branch contains the latest stable release.

If you intend to fork 0x Mesh and open a PR, you should work off of the `development` branch. Make sure you check out the `development` branch and pull the latest changes.

```text
git checkout development
git pull
```

All PRs should use `development` as the base branch. When opening a new PR, use the dropdown menu in the GitHub UI to select `development`.

![Selecting a branch](https://user-images.githubusercontent.com/800857/64901012-00272480-d64a-11e9-86f7-a2450ba8d0d9.png)

## Prerequisites

* [GNU Make](https://www.gnu.org/software/make/) If you are using a Unix-like OS, you probably already have this.
* [Go version &gt;= 1.12](https://golang.org/dl/) \(or use [the version manager called "g"](https://github.com/stefanmaric/g)\)
* [Dep package manager](https://golang.github.io/dep/docs/installation.html)
* [Node.js version &gt;=11](https://nodejs.org/en/download/) \(or use the [nvm version manager](https://github.com/creationix/nvm)\)
* [Yarn package manager](https://yarnpkg.com/en/)
* [golangci-lint version 1.16.0](https://github.com/golangci/golangci-lint#install)

## Installing Dependencies

```text
make deps
```

## Running Tests

Some of the tests depend on having a test Ethereum node running. Before running the tests, make sure you have [Docker](https://docs.docker.com/install/) installed locally and start [0xorg/ganache-cli](https://hub.docker.com/r/0xorg/ganache-cli):

```text
docker pull 0xorg/ganache-cli
docker run -ti -p 8545:8545 -e VERSION=4.3.0 0xorg/ganache-cli
```

There are various Make targets for running tests:

```bash
# Run tests in pure Go
make test-go

# Compile to WebAssembly and run tests in Node.js
make test-wasm-node

# Compile to WebAssembly and run tests in a headless Chrome browser
make test-wasm-browser

# Run tests in all available environments
make test-all
```

## Running the Linters

0x Mesh is configured to use linters for both Go and TypeScript code. To run all available linters, run:

```text
make lint
```

## Managing Dependencies

See [https://golang.github.io/dep/docs/daily-dep.html](https://golang.github.io/dep/docs/daily-dep.html).

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

When working on code with the build tag `js,wasm`, you might need to add the following to your editor config:

```javascript
{
    // ...

  "go.toolsEnvVars": {
    "GOARCH": "wasm",
    "GOOS": "js"
  },
  "go.testEnvVars": {
    "GOARCH": "wasm",
    "GOOS": "js"
    }

    // ...
}
```

