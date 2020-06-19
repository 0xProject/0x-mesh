# 0x Mesh Browser Guide

This guide will walk you through how to run 0x Mesh directly in the browser and
discuss some of the advantages and drawbacks of doing so.

## Background

Mesh is written in Go, but can be compiled to WebAssembly and run directly in
the browser. This makes it possible for users to share orders and trade directly
with one another without relying on any third-party server or database. This
approach comes with a lot of advantages, but also has some drawbacks:

### Advantages

-   Increased decentralization
-   Little to no hosting costs
-   Ability to trade experimental/niche assets

### Drawbacks

-   Longer warm-up time
-   Increased risk of trade collisions
-   Consumes more end-user resources

## @0x/mesh-browser

For your convenience, we've published an NPM package called `@0x/mesh-browser`
which includes the WebAssembly bytecode and a lightweight wrapper around it. This
package is used exactly the same way as any other NPM package and using it feels
exactly like using a native TypeScript/JavaScript library. To find more information
on optimizing the startup time of using Mesh in the browser, refer to the next
section on the `@0x/mesh-browser-lite` package.

## @0x/mesh-browser-lite

We've published a lightweight version of the `@0x/mesh-browser` package -- the
`@0x/mesh-browser-lite` package -- that provides an identical abstraction around a
browser-based Mesh node without requiring that Wasm bytecode be bundled with the
rest of the webpage's code. Additionally, this package makes use of the
WebAssembly's streaming functionality, which could provide a speedup to load
times for users of browser-based Mesh.

Using this package is a bit more complicated than using the `@0x/mesh-browser` package.
WebAssembly binaries for each version of 0x-mesh that has an associated `@0x/mesh-browser-lite`
package can be found in the [0x-mesh release notes](https://github.com/0xProject/0x-mesh/releases).
The user will need to serve the appropriate binary on a server or CDN of their choice.
The package gives users the option of providing a URL to the `loadMeshStreamingWithURLAsync`
function or a `Response` object to the `loadMeshStreamingAsync` function in their
application. The URL or `Response` option should be chosen in such a way that they
load the Mesh Binary that is being served.

## Installation

To install the `@0x/mesh-browser` NPM package, simply run:

```
npm install --save @0x/mesh-browser
```

Or if you are using `yarn`:

```
yarn add @0x/mesh-browser
```

Similarly, the `@0x/mesh-browser-lite` NPM package can be installed by running:

```
npm install --save @0x/mesh-browser-lite
```

Or if you are using `yarn`:

```
yarn add @0x/mesh-browser-lite
```

## Documentation

-   Documentation for the `@0x/mesh-browser` package is available at
    [0x-org.gitbook.io/mesh/browser-bindings/browser](https://0x-org.gitbook.io/mesh/browser-bindings/browser).
-   Documentation for the `@0x/mesh-browser-lite` package is available at
    [0x-org.gitbook.io/mesh/browser-bindings/browser-lite](https://0x-org.gitbook.io/mesh/browser-bindings/browser-lite).

### Example usage

[The webpack-example directory](../packages/webpack-example) and
[the webpack-example-lite directory](../packages/webpack-example) include
examples of how to use the `@0x/mesh-browser` and the `@0x/mesh-browser-lite` packages,
respectively.
