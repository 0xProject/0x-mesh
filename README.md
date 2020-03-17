[![Version](https://img.shields.io/badge/version-9.2.0-orange.svg)](https://github.com/0xProject/0x-mesh/releases)
[![Docs](https://img.shields.io/badge/docs-website-yellow.svg)](https://0x-org.gitbook.io/mesh)
[![Chat with us on Discord](https://img.shields.io/badge/chat-Discord-blueViolet.svg)](https://discord.gg/HF7fHwk)
[![Circle CI](https://img.shields.io/circleci/project/0xProject/0x-mesh/master.svg)](https://circleci.com/gh/0xProject/0x-mesh/tree/master)

# 0x Mesh

0x Mesh is a peer-to-peer network for sharing orders that adhere to the
[0x order message format](https://0x.org/docs/guides/v3-specification#orders).

## Project status

We have reached the point where Mesh is being used by some teams in production. We feel that for many use cases, Mesh is stable enough for production. However, we caution that there are some issues and shortcomings in its current state, which generally fall into two categories:

- Order sharing: We have recently made significant improvements to our order sharing algorithm, including reducing bandwidth usage and CPU usage by at least an order of magnitude on average. See https://github.com/0xProject/0x-mesh/pull/692 and https://github.com/0xProject/0x-mesh/pull/732. However, we are still working on accurately testing and measuring the speed at which orders propagate through the network with different network sizes and topologies. In some circumstances, it may take longer than we would like for orders to reach the majority of nodes in the network. This is an area we will continue to focus on and improve.
- Browser usage: Mesh can run directly in the browser via the [@0x/mesh-browser](https://www.npmjs.com/package/@0x/mesh-browser) package. We have supported this for a while and have examples and integration tests in this repository. While we have made recent improvements to stability and performance (see https://github.com/0xProject/0x-mesh/pull/703, https://github.com/0xProject/0x-mesh/pull/697, and https://github.com/0xProject/0x-mesh/pull/694), there are still some [important missing features and issues to address](https://github.com/0xProject/0x-mesh/issues?q=is%3Aopen+is%3Aissue+label%3Abrowser) before `@0x/mesh-browser` is feasible for most production use cases.


## Overview

0x Mesh has a lot of different use cases for different categories of users:

- Relayers can use Mesh to share orders with one another and to receive orders
  from market makers. This allows them to increase the depth of their order
  books and provide a better user experience.
- Market makers can use Mesh to reach a broader audience. Their orders will be
  sent throughout the network and picked up by many trading venues and are therefore more likely to be filled.
- Mesh allows for a new type of relayer called a "serverless relayer". In the
  serverless relayer model, each user runs Mesh in their browser and there is
  no backend server or database. Instead, peers share orders directly with one
  another. (There are pros and cons to this approach and it is probably not
  suitable for all markets).

Both Relayers and Market makers can use Mesh to watch a set of 0x orders for changes in fillability (e.g., cancellations, fills, expirations, etc...).

0x Mesh is intended to be entirely automatic. It takes care of all the work of
receiving, sharing, and validating orders so that you can focus on building your
application. When you run a 0x Mesh node, it will automatically discover peers
in the network and begin receiving orders from and sending orders to them. You
do not need to know the identities (e.g., IP address or domain name) of any
peers in the network ahead of time and they do not need to know about you.

Developers can use the JSON-RPC API to interact with a Mesh node that they
control. The API allows you to send orders into the network, receive any new
orders, and get notified when the status of an existing order changes (e.g. when
it is filled, canceled, or expired). Under the hood, Mesh performs efficient
order validation and order book pruning, which takes out a lot of the hard work
for developers.

## Documentation

You can find documentation and guides for 0x Mesh at
https://0x-org.gitbook.io/mesh.

## Development

We love receiving contributions from the community :smile: If you are interested
in helping develop 0x Mesh, please read the
[Development Guide](CONTRIBUTING.md).
If you are looking for a place to start, take a look at the
[issues page](https://github.com/0xProject/0x-mesh/issues) and don't hesitate to
[reach out to us on Discord](https://discord.gg/HF7fHwk).

## Additional Background

-   [Announcement blog post](https://blog.0xproject.com/0x-roadmap-2019-part-3-networked-liquidity-0x-mesh-9a24026202b3)
-   [MVP architecture doc](https://drive.google.com/file/d/1dAVTEND7e1sISO9VZSOou0DN-igoUi9z/view)
-   [Video of talk at ETHNewYork](https://youtu.be/YUqe4fKBA2k?t=723)
