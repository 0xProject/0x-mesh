# 0x Mesh Browser Guide

This guide will walk you through how to run 0x Mesh directly in the browser and
discuss some of the advantages and drawbacks of doing so.

## Background

Mesh is written in Go, but can be compiled to WebAssembly and run directly in
the browser. This makes it possible for users to share orders and trade directly
with one another without relying on any third-party server or database. This
approach comes with a lot of advantages, but also has some drawbacks:

### Advantages

- Increased decentralization
- Little to no hosting costs
- Ability to trade experimental/niche assets

### Drawbacks

- Longer warm-up time
- Increased risk of trade collisions
- Consumes more end-user resources

## @0x/mesh-browser

### Installing

For your convenience, we've published an NPM package called `@0x/mesh-browser`
which includes the WebAssembly bytecode and a lightweight wrapper around it. To
find more information on optimizing the startup time of using Mesh in the browser,
refer to the next section on the `@0x/mesh-browser-lite` package. You install
this package in exactly the same way as any other NPM package and using it feels
exactly like using a native TypeScript/JavaScript library.

To install the NPM package, simply run:

```
npm install --save @0x/mesh-browser
```

Or if you are using `yarn`:

```
yarn add @0x/mesh-browser
```

We recommend using a tool like [Webpack](https://webpack.js.org/) to bundle the
0x Mesh package and all your other code into a single JavaScript file.

### Documentation

Documentation for the `@0x/mesh-browser` package is available at
[0x-org.gitbook.io/mesh/browser](https://0x-org.gitbook.io/mesh/browser).

### Example usage

[The examples/browser directory](../examples/browser) includes a bare-bones
example of how to use the `@0x/mesh-browser` and bundle everything together with
Webpack.

## @0x/mesh-browser-lite

### Installing

We've published a lightweight version of the `@0x/mesh-browser` package that provides
an identical abstraction around a browser-based Mesh node without requiring that
Wasm bytecode be bundled with the rest of the webpage's code. Additionally, this
package makes use of the WebAssembly's streaming functionality, which could provide
a speedup to load times for users of browser-based Mesh.

Using this package is a bit more complicated than using the `@0x/mesh-browser` package.
This package provides an API that will create a Mesh wrapper using a provided URL
or Response object that serves the WebAssembly bytecode that will be used to run a
Mesh node.

To install the NPM package, simply run:

```
npm install --save @0x/mesh-browser-lite
```

Or if you are using `yarn`:

```
yarn add @0x/mesh-browser-lite
```

### Documentation

Documentation for the `@0x/mesh-browser-lite` package is available at
[0x-org.gitbook.io/mesh/browser-lite](https://0x-org.gitbook.io/mesh/browser-lite).

### Example usage

[The examples/browser-lite directory](../examples/browser-lite) includes a bare-bones
example of how to use the `@0x/mesh-browser-lite` and bundle everything together with
Webpack.
