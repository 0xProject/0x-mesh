[![Version](https://img.shields.io/badge/version-development-orange.svg)](https://github.com/0xProject/0x-mesh/releases)
[![Docs](https://img.shields.io/badge/docs-website-yellow.svg)](https://0x-org.gitbook.io/mesh)
[![Chat with us on Discord](https://img.shields.io/badge/chat-Discord-blueViolet.svg)](https://discord.gg/HF7fHwk)
[![Circle CI](https://img.shields.io/circleci/project/0xProject/0x-mesh/master.svg)](https://circleci.com/gh/0xProject/0x-mesh/tree/master)

# 0x Mesh

0x Mesh is a peer-to-peer network for sharing orders that adhere to the
[0x order message format](https://github.com/0xProject/0x-protocol-specification/blob/master/v2/v2-specification.md#order-message-format).

WARNING: This project is still under active development. Expect breaking changes before the official release.

## Overview

0x Mesh has a lot of different use cases for different categories of users:

- Relayers can use Mesh to share orders with one another and to receive orders
  from market makers. This allows them to increase the depth of their order
  books and provide a better user experience.
- Market makers can use Mesh to reach a broader audience. Their orders will be
  sent throughout the network and are more likely to be filled.
- Mesh allows for a new type of relayer called a "serverless relayer". In the
  serverless relayer model, each user runs Mesh in their browser and there is
  no backend server or database. Instead, peers share orders directly with one
  another. (There are pros and cons to this approach and it is probably not
  suitable for all markets).

Both Relayers and Market makers can use Mesh to watch a set of orders for changes in fillability (e.g., cancellations, fills, expirations, etc...).

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
[Development Guide](docs/development.md).
If you are looking for a place to start, take a look at the
[issues page](https://github.com/0xProject/0x-mesh/issues) and don't hesitate to
[reach out to us on Discord](https://discord.gg/HF7fHwk).

## Additional Background

-   [Announcement blog post](https://blog.0xproject.com/0x-roadmap-2019-part-3-networked-liquidity-0x-mesh-9a24026202b3)
-   [MVP architecture doc](https://drive.google.com/file/d/1dAVTEND7e1sISO9VZSOou0DN-igoUi9z/view)
-   [Video of talk at ETHNewYork](https://youtu.be/YUqe4fKBA2k?t=723)
