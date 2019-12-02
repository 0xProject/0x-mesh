# 0x Mesh Development and Contribution Guide

## Directory Location

If you are working on 0x-mesh, the root directory for the project must be at
**\$GOPATH/src/github.com/0xProject/0x-mesh**. 0x Mesh uses
[Dep](https://golang.github.io/dep/docs/installation.html) for dependency
management and does not support Go modules.

## Cloning the Repository and Opening PRs

0x Mesh uses two main branches:

1. The `development` branch contains the latest (possibly unreleased) changes
	and is not guaranteed to be stable.
2. The `master` branch contains the latest stable release.

If you intend to fork 0x Mesh and open a PR, you should work off of the
`development` branch. Make sure you check out the `development` branch and pull
the latest changes.

```
git checkout development
git pull
```

All PRs should use `development` as the base branch. When opening a new PR, use
the dropdown menu in the GitHub UI to select `development`.

![Selecting a branch](https://user-images.githubusercontent.com/800857/64901012-00272480-d64a-11e9-86f7-a2450ba8d0d9.png)

## Prerequisites

-   [GNU Make](https://www.gnu.org/software/make/) If you are using a Unix-like OS, you probably already have this.
-   [Go version 1.12.x](https://golang.org/dl/) (or use [the version manager called "g"](https://github.com/stefanmaric/g)). Go 1.13 is not supported yet (see https://github.com/0xProject/0x-mesh/issues/480).
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
[0xorg/ganache-cli](https://hub.docker.com/r/0xorg/ganache-cli). In these commands,
`$GANACHE_VERSION` should be set to the version of ganache-cli that is used in the mesh project's
CI found [here](https://github.com/0xProject/0x-mesh/blob/development/.circleci/config.yml#L10):

```
docker pull 0xorg/ganache-cli

# Run the $GANACHE_VERSION image of ganache-cli.
docker run -ti -p 8545:8545 -e VERSION=$GANACHE_VERSION 0xorg/ganache-cli
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

### Test Issues
Some of the tests will open a fairly large number of open files. The default allowance
for open files on most operating systems is 64, which will cause these tests to fail. The
allowance of open files can be configured using the following command:

```bash
# Increase number of open files that are tolerated to 2048 (a big number)
ulimit -S -n 2048
```

It may be convenient to add this line to the `.bashrc` (or `.bash_profile` for MacOs users)
file so that the change will go into effect whenever a new shell is created.

## Running the Linters

0x Mesh is configured to use linters for both Go and TypeScript code. To run all
available linters, run:

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

When working on code with the build tag `js,wasm`, you might need to add the
following to your editor config:

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
