# 0x Mesh Development and Contribution Guide

## Cloning the Repository and Opening PRs

Clone the repository with `git clone git@github.com:0xProject/0x-mesh.git`. Mesh
uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency
management, so depending on your settings you might need to clone the repository
_outside_ of your `GOPATH`.

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
-   [Go version 1.14.x](https://golang.org/dl/) (or use [the version manager called "g"](https://github.com/stefanmaric/g)).
-   [Node.js version >=11](https://nodejs.org/en/download/) (or use the [nvm version manager](https://github.com/creationix/nvm)).
-   [Yarn package manager](https://yarnpkg.com/en/).
-   [golangci-lint version 1.27.0](https://github.com/golangci/golangci-lint#install-golangci-lint).
-   [Python](https://www.python.org/downloads/). (Many OSes already have this).
-   [Google Chrome](https://www.google.com/chrome/). If you already have Google Chrome you typically don't need to do anything. On Ubuntu you can run `wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && dpkg -i google-chrome-stable_current_amd64.deb; apt-get -fy install`.
-   A C compiler such as [GCC](https://gcc.gnu.org/) or [Clang](https://clang.llvm.org/). Some OSes will already have this. On Ubuntu you can run `sudo apt-get install build-essential`.

## Installing Dependencies

```
make deps
```

## Building TypeScript packages

Mesh contains some TypeScript packages, all of which are contained in a small monorepo in the __packages/__ directory. Some
packages are published, and some are only used internally for development and testing.

To build all the TypeScript packages:

```
yarn build
```

You can also run `yarn build` inside of individual packages in order to just build that package and its dependencies. However,
changing any Go code will require running `yarn build` at the root of the project again.

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

### Potential Issues
The default maximum number of open files is too low in some operating systems
for the tests to be run successfully. If an error that reads like "Too many open files,"
it may be necessary to increase this limit. On Unix-like operating systems, the `ulimit`
command can be used as follows to accomplish this change:

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

## Running prettier

0x Mesh uses a tool called [prettier](https://prettier.io/), which is a tool that
enforces a consistent style across the Typescript codebase. The continuous integration
pipeline will fail if Typescript code is not "prettified," so `prettier` must be run
to ensure that the code is properly formatted.

The prettier tool can be run using the following command from the top-level of the
0x-Mesh repository (outside of a directory inside `packages/`):

```bash
yarn prettier
```

## Managing Dependencies

Mesh uses [Go Modules](https://github.com/golang/go/wiki/Modules) for managing
Go dependencies and [Yarn](https://yarnpkg.com/lang/en/) for managing
TypeScript/JavaScript dependencies.

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

### Prettier Configuration

Configurations for most popular text editors can be found
[here](https://prettier.io/docs/en/editors.html).
