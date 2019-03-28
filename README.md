# 0x-mesh

A peer-to-peer network for sharing orders

## Development

This project uses Go modules for managing dependencies. Check out the
[Go Modules Wiki page](https://github.com/golang/go/wiki/Modules#how-to-use-modules)
for more background on how this works.

### Prerequisites

- [Go version >= 1.12](https://golang.org/dl/) (or use [the version manager called "g"](https://github.com/stefanmaric/g))
- [Node.js version >= 10](https://nodejs.org/en/download/) (or use the [nvm version manager](https://github.com/creationix/nvm))
- [Yarn package manager](https://yarnpkg.com/en/)

### Installing dependencies

```
yarn install
```

### Running tests

Run tests in a vanilla Go environment:

```
yarn test:go
```

Run tests in a Node.js/WebAssembly environment:

```
yarn test:wasm
```

Run tests in both environments:

```
yarn test:all
```
