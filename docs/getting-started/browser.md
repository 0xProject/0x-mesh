# Browser guide

This guide will walk you through how to run 0x Mesh directly in the browser and discuss some of the advantages and drawbacks of doing so.

## Background

Mesh is written in Go, but can be compiled to WebAssembly and run directly in the browser. This makes it possible for users to share orders and trade directly with one another without relying on any third-party server or database. This approach comes with a lot of advantages, but also has some drawbacks:

### Advantages

* Increased decentralization
* Little to no hosting costs
* Ability to trade experimental/niche assets

### Drawbacks

* Longer warm-up time
* Increased risk of trade collisions
* Consumes more end-user resources 

## Installing

For your convenience, we've published an NPM package called `@0x/mesh-browser` which includes the WebAssembly bytecode and a lightweight wrapper around it. You install this package in exactly the same way as any other NPM package and using it feels exactly like using a native TypeScript/JavaScript library.

To install the NPM package, simply run:

```text
npm install --save @0x/mesh-browser
```

Or if you are using `yarn`:

```text
yarn add @0x/mesh-browser
```

We recommend using a tool like [Webpack](https://webpack.js.org/) to bundle the 0x Mesh package and all your other code into a single JavaScript file.

## Documentation

Documentation for the `@0x/mesh-browser` package is available at [0x-org.gitbook.io/mesh/browser](https://0x-org.gitbook.io/mesh/browser).

## Example usage

[The examples/browser directory](https://github.com/0xProject/0x-mesh/tree/9e86d4d3bf18fb19c9bbf04dcd3321490c819ba4/examples/browser/README.md) includes a bare-bones example of how to use the `@0x/mesh-browser` and bundle everything together with Webpack.

