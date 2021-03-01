[![Version](https://img.shields.io/badge/version-11.0.0-orange.svg)](https://github.com/0xProject/0x-mesh/releases)

# 0x Mesh GraphQL API Documentation

The GraphQL API is intended to be a _private_ API. The API should only be
accessible to the developers running the Mesh node and should not be exposed to
the public. The API runs on a separate port from the peer-to-peer protocols and
access to it can be controlled via a firewall.

Peers in the network do not use the GraphQL API and instead use a peer-to-peer
PubSub mechanism (usually this is not something you need to worry about).

## About GraphQL

[GraphQL](https://graphql.org/) is a structured query language for APIs. It:

-   Is transport layer agnostic (e.g. HTTP, WebSockets, or calling a function).
-   Is type-safe and uses well-structured schemas.
-   Has wide support across many programming languages.
-   Has great tooling (including automatic doc generation and playground environments).
-   Features built-in support for subscriptions.
-   Allows clients to only receive the data that they need in the format that they need it.

## Playground Environment

We have deployed a public playground environment for exploring the GraphQL API. You can access the playground
at [https://meshmock.spaceship.0x.org/](https://meshmock.spaceship.0x.org/). It supports auto-completion,
syntax highlighting, and subscriptions. In addition, interactive documentation for the API can be found
by clicking the "docs" button on the righthand side of the screen. See the Example Queries section below
for some queries to try.

## Recommended Clients:

-   [Typescript GraphQL client](graphql_clients/typescript/README.md).
-   [Go GraphQL client](https://godoc.org/github.com/0xProject/0x-mesh/graphql/client).
-   For any other languages see the [GraphQL page on recommended clients](https://graphql.org/code/#graphql-clients).

If you prefer, you can also [send requests directly over HTTP](https://graphql.org/learn/serving-over-http/)
without using a dedicated client.

## Example Queries

This section includes some example queries which you can copy and paste in the playground. Of course, you would
typically write queries programmatically with a GraphQL client, not by manually writing them. This is just for
illustrative purposes.

### Getting a Specific Order

You can get the details for any order by its hash:

```graphql
{
    order(hash: "0x38c1b56f95bf168b303e4b62d64f7f475f2ac34124e9678f0bd852f95a4ca377") {
        hash
        chainId
        exchangeAddress
        makerAddress
        makerAssetData
        makerAssetAmount
        makerFeeAssetData
        makerFee
        takerAddress
        takerAssetData
        takerAssetAmount
        takerFeeAssetData
        takerFee
        senderAddress
        feeRecipientAddress
        expirationTimeSeconds
        salt
        signature
        remainingFillableTakerAssetAmount
    }
}
```

### Querying and Filtering Orders

You can get all orders via the `orders` query. By default, it will return up to 100 orders at a time sorted by their hash. You can
also change the number of orders returned via the `limit` argument.

```graphql
{
    orders {
        hash
        chainId
        exchangeAddress
        makerAddress
        makerAssetData
        makerAssetAmount
        makerFeeAssetData
        makerFee
        takerAddress
        takerAssetData
        takerAssetAmount
        takerFeeAssetData
        takerFee
        senderAddress
        feeRecipientAddress
        expirationTimeSeconds
        salt
        signature
        remainingFillableTakerAssetAmount
    }
}
```

The `orders` query supports many different options. Here's an example of how to get orders with a minimum `expirationTimeSeconds`
and minimum `remainingFillableAssetAmount`. You can use this to exclude dust orders and orders which may expire too soon.

```graphql
{
    orders(
        filters: [
            { field: remainingFillableTakerAssetAmount, kind: GREATER_OR_EQUAL, value: "150000" }
            { field: expirationTimeSeconds, kind: GREATER_OR_EQUAL, value: "1598733429" }
        ]
    ) {
        hash
        makerAddress
        makerAssetData
        makerAssetAmount
        takerAddress
        takerAssetData
        takerAssetAmount
        expirationTimeSeconds
        remainingFillableTakerAssetAmount
    }
}
```

Here's an example of sorting orders by the `remainingFillableAssetAmount`. You can combine any number
of filters and sorts in an `orders` query.

```graphql
{
    orders(sort: [{ field: remainingFillableTakerAssetAmount, direction: DESC }]) {
        hash
        makerAddress
        makerAssetData
        makerAssetAmount
        takerAddress
        takerAssetData
        takerAssetAmount
        expirationTimeSeconds
        remainingFillableTakerAssetAmount
    }
}
```

### Adding Orders

You can add orders using a [`mutation`](https://graphql.org/learn/queries/#mutations).

```graphql
mutation AddOrders {
    addOrders(
        orders: [
            {
                signature: "0x1c91055b1ce93cdd341c423b889be703ce436e25fe62d94aabbae97528b4d247646c3cd3a20f0566540ac5668336d147d844cf1a7715d700f1a7c3e72f1c60e21502"
                senderAddress: "0x0000000000000000000000000000000000000000"
                makerAddress: "0xd965a4f8f5b49dd2f5ba83ef4e61880d0646fd00"
                takerAddress: "0x0000000000000000000000000000000000000000"
                makerFee: "1250000000000000"
                takerFee: "0"
                makerAssetAmount: "50000000000000000"
                takerAssetAmount: "10"
                makerAssetData: "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
                takerAssetData: "0xa7cb5fb7000000000000000000000000d4690a51044db77d91d7aa8f7a3a5ad5da331af0000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000e3a2a1f2146d86a604adc220b4967a898d7fe0700000000000000000000000009a379ef7218bcfd8913faa8b281ebc5a2e0bc040000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001360000000000000000000000000000000000000000000000000000000000000004"
                salt: "1584796917698"
                exchangeAddress: "0x61935cbdd02287b511119ddb11aeb42f1593b7ef"
                feeRecipientAddress: "0x0d056bb17ad4df5593b93a1efc29cb35ba4aa38d"
                expirationTimeSeconds: "1595164917"
                makerFeeAssetData: "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
                chainId: 1
                takerFeeAssetData: "0x"
            }
            {
                signature: "0x1bf931ab06551bbbd3a7e272ca4833503d768caca2cac564b157b46c906c7b41c57fd6146b500e1ad2dac729c351142764cb76efc975c6d7c64aef6cf7930c075d02"
                senderAddress: "0x0000000000000000000000000000000000000000"
                makerAddress: "0x0c5fa5fa51d84227bfacdc56b36329286b37d051"
                takerAddress: "0x0000000000000000000000000000000000000000"
                makerFee: "0"
                takerFee: "0"
                makerAssetAmount: "50911000000000000"
                takerAssetAmount: "10000000000000000000"
                makerAssetData: "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
                takerAssetData: "0xf47261b000000000000000000000000058b6a8a3302369daec383334672404ee733ab239"
                salt: "1590244702461"
                exchangeAddress: "0x61935cbdd02287b511119ddb11aeb42f1593b7ef"
                feeRecipientAddress: "0xa258b39954cef5cb142fd567a46cddb31a670124"
                expirationTimeSeconds: "1592663792"
                makerFeeAssetData: "0x"
                chainId: 1
                takerFeeAssetData: "0x"
            }
        ]
    ) {
        accepted {
            order {
                hash
            }
            isNew
        }
        rejected {
            hash
            code
            message
        }
    }
}
```

### Subscribing to Order Events

You can subscribe to order events via a [`subscription`](https://graphql.org/blog/subscriptions-in-graphql-and-relay/).

```graphql
subscription {
    orderEvents {
        timestamp
        endState
        order {
            hash
            remainingFillableTakerAssetAmount
        }
    }
}
```

### Getting Stats

You can get some stats about your Mesh node via the `stats` query.

```graphql
{
    stats {
        version
        pubSubTopic
        rendezvous
        peerID
        ethereumChainID
        latestBlock {
            number
            hash
        }
        numPeers
        numOrders
        numOrdersIncludingRemoved
        startOfCurrentUTCDay
        ethRPCRequestsSentInCurrentUTCDay
        ethRPCRateLimitExpiredRequests
        maxExpirationTime
    }
}
```

## Additional Tips

### Pagination

We recommend paginating through orders by using `filters` and `limit`. So for example, if you want to sort orders by their hash
(which is the default), you first send a query without any filters:

```graphql
{
    orders {
        hash
        makerAddress
        makerAssetData
        makerAssetAmount
        takerAddress
        takerAssetData
        takerAssetAmount
        expirationTimeSeconds
        remainingFillableTakerAssetAmount
    }
}
```

The orders in the response will be sorted by `hash` (which is the default). Look at the last order you received,
which in this case has a hash of `0x75d2b56b11f21235ec8faec8be9d081090678cf62f5c69fa118236d829424719`. Send the
next request by using the last hash you received in a filter:

```graphql
{
    orders(
        filters: [
            { field: hash, kind: GREATER, value: "0x75d2b56b11f21235ec8faec8be9d081090678cf62f5c69fa118236d829424719" }
        ]
    ) {
        hash
        makerAddress
        makerAssetData
        makerAssetAmount
        takerAddress
        takerAssetData
        takerAssetAmount
        expirationTimeSeconds
        remainingFillableTakerAssetAmount
    }
}
```

This will return any orders with a hash greater than `0x75d2b56b11f21235ec8faec8be9d081090678cf62f5c69fa118236d829424719`.
Repeat this process, changing the `hash` each time, until there are no orders left.

There may be orders added or removed while you are in the process of paginating. Following this method guarantees that:

1. No order will be included more than once.
2. Any order which was present at the start of pagination _and_ at the end of pagination will be included.
3. Any order which was added or removed after pagination started may or may not be included.

### Query Fragments

GraphQL requires you to specify all the fields that you want included in the response. However, you
can use [query fragments](https://graphql.org/learn/queries/#fragments) to avoid repeating the same
fields over and over. Here's an example of a query fragment that includes all the fields of an order.

```graphql
fragment AllOrderFields on OrderWithMetadata {
    hash
    chainId
    exchangeAddress
    makerAddress
    makerAssetData
    makerAssetAmount
    makerFeeAssetData
    makerFee
    takerAddress
    takerAssetData
    takerAssetAmount
    takerFeeAssetData
    takerFee
    senderAddress
    feeRecipientAddress
    expirationTimeSeconds
    salt
    signature
    remainingFillableTakerAssetAmount
}

{
    order(hash: "0x06d15403630b6d73fbacbf0864eb76c2db3d6e6fc8adec8a95fc536593f17c54") {
        ...AllOrderFields
    }
}
```
