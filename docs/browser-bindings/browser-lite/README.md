# @0x/mesh-browser-lite - v10.0.0-unstable-3

## @0x/mesh-browser-lite

This packages provides a set of Typescript and Javascript bindings for running a
0x-mesh node in the browser. The browser node's Wasm binary is not bundled in
this package and is instead expected to be served by the consumer of the package.
This package has a smaller bundle size than the `@0x/mesh-browser` package and
may have faster load times.

## Installation

```bash
yarn add @0x/mesh-browser-lite
```

If your project is in [TypeScript](https://www.typescriptlang.org/), add the following to your `tsconfig.json`:

```json
"compilerOptions": {
    "typeRoots": ["node_modules/@types"],
}
```

## Contributing

If you would like to contribute bug fixes or new features to the client, checkout the [0xproject/0x-mesh](https://github.com/0xProject/0x-mesh) project and use the below commands to install the dependencies, build, lint and test your changes.

### Install dependencies

```bash
yarn install
```

### Build

```bash
yarn build
```

### Lint

```bash
yarn lint
```
