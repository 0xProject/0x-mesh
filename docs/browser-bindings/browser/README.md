# @0x/mesh-browser - v10.0.0-unstable-3

## @0x/mesh-browser

This package provides an easy way to run a browser-based Mesh node. Specifically, it
provides Typescript and Javascript bindings that can be used to interact with a Mesh
node that is running in the browser and handles the process of loading the mesh node
on the webpage. Because of the fact that this package handles Wasm loading, it is
considerably heavier-weight and may take longer to load than the @0x/mesh-browser-lite
package.

## Installation

```bash
yarn add @0x/mesh-browser
```

If your project is in [TypeScript](https://www.typescriptlang.org/), add the following to your `tsconfig.json`:

```json
"compilerOptions": {
    "typeRoots": ["node_modules/@0x/typescript-typings/types", "node_modules/@types"],
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

### Clean

```bash
yarn clean
```

### Lint

```bash
yarn lint
```
