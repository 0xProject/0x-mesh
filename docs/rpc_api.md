[![Version](https://img.shields.io/badge/version-10.0.0--unstable--3-orange.svg)](https://github.com/0xProject/0x-mesh/releases)

# 0x Mesh JSON-RPC API Documentation

The JSON-RPC API is intended to be a _private_ API. The API should only be
accessible to the developers running the Mesh node and should not be exposed to
the public. The API runs on a separate port from the peer-to-peer protocols and
access to it can be controlled via a firewall.

Peers in the network do not use the JSON-RPC API and instead use a peer-to-peer
PubSub mechanism (usually this is not something you need to worry about).

## Similarities to the Ethereum JSON-RPC API

Our JSON-RPC API is very similar to the
[Ethereum JSON-RPC API](https://github.com/ethereum/wiki/wiki/JSON-RPC); we even
use a lot of the same code from `go-ethereum`.

Some key differences:

-   It is only accessible via HTTP and WebSocket transports (IPC not supported)
-   uint256 amounts should not be hex encoded, but rather sent as numerical strings

Since the API adheres to the [JSON-RPC 2.0 spec](https://www.jsonrpc.org/specification),
you can use any JSON-RPC 2.0 compliant client in the language of your choice.
The clients made for Ethereum work even better since they extend the standard to
include [subscriptions](https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB).

### Recommended Clients:

-   Javascript/Typescript: We've published a [Typescript RPC client](json_rpc_clients/typescript/README.md).
-   Python: [Web3.py](https://github.com/ethereum/web3.py) has a [WebSocketProvider](https://web3py.readthedocs.io/en/stable/providers.html#web3.providers.websocket.WebsocketProvider) you can use.
-   Go: Mesh ships with a [Golang RPC client](https://godoc.org/github.com/0xProject/0x-mesh/rpc#Client)
    -   see the [examples](../examples/go/) directory for example usage.

## API

### `mesh_addOrders`

Adds an array of 0x signed orders to the Mesh node.

**Example payload:**

```json
{
    "jsonrpc": "2.0",
    "method": "mesh_addOrders",
    "params": [
        [
            {
                "makerAddress": "0x6440b8c5f5a3c725eb394c7c40994afaf50a0d39",
                "takerAddress": "0x0000000000000000000000000000000000000000",
                "feeRecipientAddress": "0xa258b39954cef5cb142fd567a46cddb31a670124",
                "senderAddress": "0x0000000000000000000000000000000000000000",
                "makerAssetAmount": "1233400000000000",
                "takerAssetAmount": "12334000000000000000000",
                "makerFee": "0",
                "takerFee": "0",
                "exchangeAddress": "0x080bf510fcbf18b91105470639e9561022937712",
                "chainId": 1,
                "expirationTimeSeconds": "1560917245",
                "signature": "0x1b6a49302774b0b0e14ef59e91fcf950dfb7db5705ae6929e06198518b1105301d4ef94b1b4760e550378bb5b7746b1a29c174290afe9448324cef4112dd03d7a103",
                "salt": "1545196045897",
                "makerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                "makerFeeAssetData": "0x",
                "takerAssetData": "0xf47261b00000000000000000000000000d8775f648430679a709e98d2b0cb6250d2887ef",
                "takerFeeAssetData": "0x"
            }
        ]
    ],
    "id": 1
}
```

**Example response:**

```json
{
    "jsonrpc": "2.0",
    "id": 2,
    "result": {
        "accepted": [
            {
                "orderHash": "0x4e7269386c8f2234305aafb421ba470f39064d79c4826006eaffe723b2066272",
                "signedOrder": {
                    "makerAddress": "0x8cff49b26d4d13e0601769f8a60fd697b713b9c6",
                    "makerAssetData": "0xf47261b0000000000000000000000000c778417e063141139fce010982780140aa0cd5ab",
                    "makerFeeAssetData": "0x",
                    "makerAssetAmount": "100000000000000000",
                    "makerFee": "0",
                    "takerAddress": "0x0000000000000000000000000000000000000000",
                    "takerAssetData": "0xf47261b0000000000000000000000000ff67881f8d12f372d91baae9752eb3631ff0ed00",
                    "takerFeeAssetData": "0x",
                    "takerAssetAmount": "1000000000000000000",
                    "takerFee": "0",
                    "senderAddress": "0x0000000000000000000000000000000000000000",
                    "exchangeAddress": "0x080bf510fcbf18b91105470639e9561022937712",
                    "chainId": 1,
                    "feeRecipientAddress": "0x0000000000000000000000000000000000000000",
                    "expirationTimeSeconds": "1559826927",
                    "salt": "48128453606684653105952683301312821720867493716494911784363103883716429240740",
                    "signature": "0x1cf5839d9a0025e684c3663151b1db14533cc8c9e495fb92543a37a7fffc0677a23f3b6d66a1f56d3fda46eb5277b4a91c7b7faad4fdaaa5aac9a1185dd545a8a002"
                },
                "fillableTakerAssetAmount": 1000000000000000000
            }
        ],
        "rejected": []
    }
}
```

Within the context of this endpoint:

-   _accepted_: means the order was found to be fillable for a non-zero amount and was therefore added to 0x Mesh (unless it already added of course)
-   _rejected_: means the order was not added to Mesh, however there could be many reasons for this. For example:
    -   The order could have been unfillable
    -   It could have failed some Mesh-specific validation (e.g., max order acceptable size in bytes)
    -   The network request to the Ethereum RPC endpoint used to validate the order failed

Some _rejected_ reasons warrant attempting to add the order again. Currently, the only reason we recommend re-trying adding the order is for the `NetworkRequestFailed` status code. Make sure to leave some time between attempts.

See the [AcceptedOrderInfo](https://godoc.org/github.com/0xProject/0x-mesh/zeroex/ordervalidator#AcceptedOrderInfo) and [RejectedOrderInfo](https://godoc.org/github.com/0xProject/0x-mesh/zeroex/ordervalidator#RejectedOrderInfo) type definitions as well as all the possible [RejectedOrderStatus](https://godoc.org/github.com/0xProject/0x-mesh/zeroex/ordervalidator#pkg-variables) types that could be returned.

**Note:** The `fillableTakerAssetAmount` takes into account the amount of the order that has already been filled AND the maker's balance/allowance. Thus, it represents the amount this order could _actually_ be filled for at this moment in time.

### `mesh_getOrders`

Gets orders already stored in a Mesh node at a particular snapshot of the DB state. This is a paginated endpoint with parameters (page, perPage and snapshotID).

**Example payload:**

```json
{
    "jsonrpc": "2.0",
    "method": "mesh_getOrders",
    "params": [1, 100, ""],
    "id": 1
}
```

This payload is requesting 100 orders from the 1st page (think: offset). The third parameter is the `snapshotID` which should be left empty for the first request. The response will include the snapshotID that can then be supplied in subsequent requests.

**Example response:**

```json
{
    "jsonrpc": "2.0",
    "result": {
        "snapshotID": "f47ac10b-58cc-0372-8567-0e02b2c3d479",
        "ordersInfos": [
            {
                "orderHash": "0xa0fcb54919f0b3823aa14b3f511146f6ac087ab333a70f9b24bbb1ba657a4250",
                "signedOrder": {
                    "makerAddress": "0xa3eCE5D5B6319Fa785EfC10D3112769a46C6E149",
                    "makerAssetData": "0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498",
                    "makerFeeAssetData": "0x",
                    "makerAssetAmount": "1000000000000000000",
                    "makerFee": "0",
                    "takerAddress": "0x0000000000000000000000000000000000000000",
                    "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                    "takerFeeAssetData": "0x",
                    "takerAssetAmount": "10000000000000000000000",
                    "takerFee": "0",
                    "senderAddress": "0x0000000000000000000000000000000000000000",
                    "exchangeAddress": "0x080bf510fcbf18b91105470639e9561022937712",
                    "chainId": 1,
                    "feeRecipientAddress": "0x0000000000000000000000000000000000000000",
                    "expirationTimeSeconds": "1586340602",
                    "salt": "41253767178111694375645046549067933145709740457131351457334397888365956743955",
                    "signature": "0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02"
                },
                "fillableTakerAssetAmount": "10000000000000000000000"
            }
        ]
    },
    "id": 1
}
```

### `mesh_getStats`

Gets certain configurations and stats about a Mesh node.

**Example payload:**

```json
{
    "jsonrpc": "2.0",
    "method": "mesh_getStats",
    "params": [],
    "id": 1
}
```

**Example response:**

```json
{
    "jsonrpc": "2.0",
    "result": {
        "version": "development",
        "pubSubTopic": "/0x-orders/network/1/version/1",
        "rendezvous": "/0x-mesh/network/1/version/1",
        "peerID": "16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF",
        "ethereumChainID": 1,
        "latestBlock": {
            "number": 8253150,
            "hash": "0x84aaae84147fc42fc77b33e2d3e05d86272663792d9cacaa8dc89f207b4d0642"
        },
        "numPeers": 18,
        "numOrders": 1095,
        "numOrdersIncludingRemoved": 1134,
        "startOfCurrentUTCDay": "1257811200",
        "ethRPCRequestsSentInCurrentUTCDay": 5039,
        "ethRPCRateLimitExpiredRequests": 0,
        "maxExpirationTime": "717784680"
    },
    "id": 1
}
```

### `mesh_subscribe` to `orders` topic

Allows the caller to subscribe to a stream of `OrderEvents`. An `OrderEvent` contains either newly discovered orders found by Mesh via the P2P network, or updates to the fillability of a previously discovered order (e.g., if an order gets filled, cancelled, expired, etc...). `OrderEvent`s _do not_ correspond 1-to-1 to smart contract events. Rather, an `OrderEvent` about an orders fillability change represents the aggregate change to it's fillability given _all_ the transactions included within the most recently mined/reverted blocks.

**Example:** If an order is both `filled` and `cancelled` within a single block, the `EndState`
of the `OrderEvent` will be `CANCELLED` (since this is the final state of the order after this block is
mined). The `OrderEvent` _will_ however list the contract events intercepted that could have impacted
this orders fillability. This list will include both the fill event and cancellation event.

Mesh has implemented subscriptions in the [same manner as Geth](https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB). In order to start a subscription, you must send the following payload:

```json
{
    "jsonrpc": "2.0",
    "method": "mesh_subscribe",
    "params": ["orders"],
    "id": 1
}
```

**Example response:**

```json
{
    "jsonrpc": "2.0",
    "result": "0xcd0c3e8af590364c09d0fa6a1210faf5",
    "id": 1
}
```

`result` contains the `subscriptionId` that uniquely identifies this subscription. The subscription is now active. You will now receive event payloads from Mesh of the following form:

**Example event:**

```json
{
    "jsonrpc": "2.0",
    "method": "mesh_subscription",
    "params": {
        "subscription": "0xcd0c3e8af590364c09d0fa6a1210faf5",
        "result": [
            {
                "orderHash": "0x96e6eb6174dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ecc13fd4",
                "signedOrder": {
                    "makerAddress": "0x50f84bbee6fb250d6f49e854fa280445369d64d9",
                    "makerAssetData": "0xf47261b00000000000000000000000000f5d2fb29fb7d3cfee444a200298f468908cc942",
                    "makerFeeAssetData": "0x",
                    "makerAssetAmount": "4424020538752105500000",
                    "makerFee": "0",
                    "takerAddress": "0x0000000000000000000000000000000000000000",
                    "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                    "takerFeeAssetData": "0x",
                    "takerAssetAmount": "1000000000000000061",
                    "takerFee": "0",
                    "senderAddress": "0x0000000000000000000000000000000000000000",
                    "exchangeAddress": "0x080bf510fcbf18b91105470639e9561022937712",
                    "chainId": 1,
                    "feeRecipientAddress": "0xa258b39954cef5cb142fd567a46cddb31a670124",
                    "expirationTimeSeconds": "1559422407",
                    "salt": "1559422141994",
                    "signature": "0x1cf16c2f3a210965b5e17f51b57b869ba4ddda33df92b0017b4d8da9dacd3152b122a73844eaf50ccde29a42950239ba36a525ed7f1698a8a5e1896cf7d651aed203"
                },
                "endState": "CANCELLED",
                "fillableTakerAssetAmount": 0,
                "contractEvents": [
                    {
                        "blockHash": "0x1be2eb6174dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ec11a4d2",
                        "txHash": "0xbcce172374dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ec232e3a",
                        "txIndex": 23,
                        "logIndex": 0,
                        "isRemoved": false,
                        "address": "0x4f833a24e1f95d70f028921e27040ca56e09ab0b",
                        "kind": "ExchangeCancelEvent",
                        "parameters": {
                            "makerAddress": "0x50f84bbee6fb250d6f49e854fa280445369d64d9",
                            "senderAddress": "0x0000000000000000000000000000000000000000",
                            "feeRecipientAddress": "0xa258b39954cef5cb142fd567a46cddb31a670124",
                            "orderHash": "0x96e6eb6174dbf0458686bdae44c9a330d9a9eb563962512a7be545c4ecc13fd4",
                            "makerAssetData": "0xf47261b00000000000000000000000000f5d2fb29fb7d3cfee444a200298f468908cc942",
                            "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
                        }
                    }
                ]
            }
        ]
    }
}
```

See the [OrderEvent](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#OrderEvent) type declaration as well as the [OrderEventEndState](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#pkg-constants) types for a complete list of the events that could be emitted.

To unsubscribe, send a `mesh_unsubscribe` request specifying the `subscriptionId`.

**Example unsubscription payload:**

```json
{
    "id": 1,
    "method": "mesh_unsubscribe",
    "params": ["0xcd0c3e8af590364c09d0fa6a1210faf5"]
}
```

### `mesh_subscribe` to `heartbeat` topic

After a sustained network disruption, it is possible that a WebSocket connection between client and server fails to reconnect. Both sides of the connection are unable to distinguish between network latency and a dropped connection and might continue to wait for new messages on the dropped connection. In order to avoid this, and promptly establish a new connection, clients can subscribe to a heartbeat from the server. The server will emit a heartbeat every 5 seconds. If the client hasn't received the expected heartbeat in a while, it can proactively close the connection and establish a new one. There are affordances for checking this edge-case in the [WebSocket specification](https://tools.ietf.org/html/rfc6455#section-5.5.2) however our research has found that [many WebSocket clients](https://github.com/0xProject/0x-mesh/issues/170#issuecomment-503391627) fail to provide this functionality. We therefore decided to support it at the application-level.

```json
{
    "jsonrpc": "2.0",
    "method": "mesh_subscribe",
    "params": ["heartbeat"],
    "id": 1
}
```

**Example response:**

```json
{
    "jsonrpc": "2.0",
    "result": "0xab1a3e8af590364c09d0fa6a12103ada",
    "id": 1
}
```

`result` contains the `subscriptionId` that uniquely identifies this subscription. The subscription is now active. You will now receive event payloads from Mesh of the following form:

**Example event:**

```json
{
    "jsonrpc": "2.0",
    "method": "mesh_subscription",
    "params": {
        "subscription": "0xab1a3e8af590364c09d0fa6a12103ada",
        "result": "tick"
    }
}
```

To unsubscribe, send a `mesh_unsubscribe` request specifying the `subscriptionId`.

**Example unsubscription payload:**

```json
{
    "id": 1,
    "method": "mesh_unsubscribe",
    "params": ["0xab1a3e8af590364c09d0fa6a12103ada"]
}
```
