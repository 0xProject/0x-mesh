# 0x Mesh Docs

Welcome to the [0x Mesh](https://github.com/0xProject/0x-mesh) documentation! 0x Mesh is a peer-to-peer network for sharing orders that adhere to the [0x order message format](https://github.com/0xProject/0x-protocol-specification/blob/master/v2/v2-specification.md#order-message-format).

Some resources:

- [Announcement blog post](https://blog.0xproject.com/0x-roadmap-2019-part-3-networked-liquidity-0x-mesh-9a24026202b3)
- [MVP architecture doc](https://drive.google.com/file/d/1dAVTEND7e1sISO9VZSOou0DN-igoUi9z/view)
- [Github repo](https://github.com/0xProject/0x-mesh/)

### Supported networks

- Mainnet
- Kovan
- Ropsten
- Rinkeby
- [Ganache snapshot](https://cloud.docker.com/u/0xorg/repository/docker/0xorg/mesh-ganache-cli)

## Running Mesh

Make sure you have Docker installed. Then run:

```bash
docker run \
-it \
--rm \
-p 60557:60557 \
-p 60558:60558 \
-e ETHEREUM_NETWORK_ID="1" \
-e ETHEREUM_RPC_URL="https://mainnet.infura.io/v3/a9a23d2566e542629179d6372ace13c9" \
-e VERBOSITY=5 \
-e PRIVATE_KEY_PATH="" \
-e USE_BOOTSTRAP_LIST=false \
-e BLOCK_POLLING_INTERVAL="5s" \
0xorg/mesh:latest
```

**Notes:**

- `60557` is the `RPC_PORT` and `60558` is the `P2P_LISTEN_PORT`
- In order to enable P2P order discovery and sharing, set `USE_BOOTSTRAP_LIST` to `true`.
- Running a VPN may interfere with Mesh. If you are having difficulty connecting to peers, disable your VPN.
- If you are running against a POA testnet (e.g., Kovan), you might want to shorten the `BLOCK_POLLING_INTERVAL` since blocks are mined more frequently then on mainnet.

### Persisting State

If you want the Mesh state to persist across Docker container re-starts, mount a local `0x_mesh` directory to your container. Add the following to the `docker run` command above:

```
-v {abs_local_path}/0x_mesh:/usr/mesh/0x_mesh
```

**Note:** Replace `{abs_local_path}` with the absolute path to the desired `0x_mesh` directory on the host machine.

The Mesh database will now be stored within `0x_mesh/db`.

## Possible environment variables

All possible env vars are detailed in the [Config](https://godoc.org/github.com/0xProject/0x-mesh/core#Config) struct. They are copied here for convenience, although the source code is authoritative.

```go
type Config struct {
    // Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
    Verbosity int `envvar:"VERBOSITY" default:"2"`
    // DataDir is the directory to use for persisting all data, including the
    // database and private key files.
    DataDir string `envvar:"DATA_DIR" default:"0x_mesh"`
    // P2PListenPort is the port on which to listen for new peer connections. By
    // default, 0x Mesh will let the OS select a randomly available port.
    P2PListenPort int `envvar:"P2P_LISTEN_PORT" default:"0"`
    // RPCPort is the port to use for the JSON RPC API over WebSockets. By
    // default, 0x Mesh will let the OS select a randomly available port.
    RPCPort int `envvar:"RPC_PORT" default:"0"`
    // EthereumRPCURL is the URL of an Etheruem node which supports the JSON RPC
    // API.
    EthereumRPCURL string `envvar:"ETHEREUM_RPC_URL"`
    // EthereumNetworkID is the network ID to use when communicating with
    // Ethereum.
    EthereumNetworkID int `envvar:"ETHEREUM_NETWORK_ID"`
    // UseBootstrapList is whether to use the predetermined list of peers to
    // bootstrap the DHT and peer discovery.
    UseBootstrapList bool `envvar:"USE_BOOTSTRAP_LIST" default:"false"`
    // OrderExpirationBuffer is the amount of time before the order's stipulated expiration time
    // that you'd want it pruned from the Mesh node.
    OrderExpirationBuffer time.Duration `envvar:"ORDER_EXPIRATION_BUFFER" default:"10s"`
    // BlockPollingInterval is the polling interval to wait before checking for a new Ethereum block
    // that might contain transactions that impact the fillability of orders stored by Mesh. Different
    // networks have different block producing intervals: POW networks are typically slower (e.g., Mainnet)
    // and POA networks faster (e.g., Kovan) so one should adjust the polling interval accordingly.
    BlockPollingInterval time.Duration `envvar:"BLOCK_POLLING_INTERVAL" default:"5s"`
}
```

## JSON-RPC API

Our JSON-RPC API is very similar to the [Ethereum JSON-RPC API](https://github.com/ethereum/wiki/wiki/JSON-RPC), we even use a bunch of `go-ethereum`'s code to generate it.

#### Some differences:

- It is **only accessible via a WebSocket connection**
- uint256 amounts should not be hex encoded, but rather sent as numerical strings

Since the API adheres to the [JSON-RPC 2.0 spec](https://www.jsonrpc.org/specification), you can use any JSON-RPC 2.0 compliant client in the language of your choice. The clients made for Ethereum work even better since they extend the standard to include [subscriptions](https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB).

#### Recommended clients:

- Javascript/Typescript: [Web3-providers](https://www.npmjs.com/package/web3-providers)
  - See our [example Mesh WS client](https://github.com/0xProject/0x-mesh-demo-client-javascript) built with it
- Python: [Web3.py](https://github.com/ethereum/web3.py) has a [WebSocketProvider](https://web3py.readthedocs.io/en/stable/providers.html#web3.providers.websocket.WebsocketProvider) you can use
- Go: Mesh ships with a [Mesh RPC client](https://godoc.org/github.com/0xProject/0x-mesh/rpc#Client)
  - see the [demos](https://github.com/0xProject/0x-mesh/tree/master/cmd/demo) for example usage

### Methods

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
        "exchangeAddress": "0x4f833a24e1f95d70f028921e27040ca56e09ab0b",
        "expirationTimeSeconds": "1560917245",
        "signature": "0x1b6a49302774b0b0e14ef59e91fcf950dfb7db5705ae6929e06198518b1105301d4ef94b1b4760e550378bb5b7746b1a29c174290afe9448324cef4112dd03d7a103",
        "salt": "1545196045897",
        "makerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
        "takerAssetData": "0xf47261b00000000000000000000000000d8775f648430679a709e98d2b0cb6250d2887ef"
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
          "makerAssetAmount": "100000000000000000",
          "makerFee": "0",
          "takerAddress": "0x0000000000000000000000000000000000000000",
          "takerAssetData": "0xf47261b0000000000000000000000000ff67881f8d12f372d91baae9752eb3631ff0ed00",
          "takerAssetAmount": "1000000000000000000",
          "takerFee": "0",
          "senderAddress": "0x0000000000000000000000000000000000000000",
          "exchangeAddress": "0x4530c0483a1633c7a1c97d2c53721caff2caaaaf",
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

- _accepted_: means the order was found to be fillable for a non-zero amount and was therefore added to 0x Mesh (unless it already added of course)
- _rejected_: means the order was not added to Mesh, however there could be many reasons for this. For example:
  - The order could have been unfillable
  - It could have failed some Mesh-specific validation (e.g., max order acceptable size in bytes)
  - The network request to the Ethereum RPC endpoint used to validate the order failed

Some _rejected_ reasons warrant attempting to add the order again. Currently, the only reason we recommend re-trying adding the order is for the `NetworkRequestFailed` status code. Make sure to leave some time between attempts.

See the [AcceptedOrderInfo](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#AcceptedOrderInfo) and [RejectedOrderInfo](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#RejectedOrderInfo) type definitions as well as all the possible [RejectedOrderStatus](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#pkg-variables) types that could be returned.

**Note:** The `fillableTakerAssetAmount` takes into account the amount of the order that has already been filled AND the maker's balance/allowance. Thus, it represents the amount this order could _actually_ be filled for at this moment in time.

### `mesh_subscribe` to `orders` topic

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
          "makerAssetAmount": "4424020538752105500000",
          "makerFee": "0",
          "takerAddress": "0x0000000000000000000000000000000000000000",
          "takerAssetData": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
          "takerAssetAmount": "1000000000000000061",
          "takerFee": "0",
          "senderAddress": "0x0000000000000000000000000000000000000000",
          "exchangeAddress": "0x4f833a24e1f95d70f028921e27040ca56e09ab0b",
          "feeRecipientAddress": "0xa258b39954cef5cb142fd567a46cddb31a670124",
          "expirationTimeSeconds": "1559422407",
          "salt": "1559422141994",
          "signature": "0x1cf16c2f3a210965b5e17f51b57b869ba4ddda33df92b0017b4d8da9dacd3152b122a73844eaf50ccde29a42950239ba36a525ed7f1698a8a5e1896cf7d651aed203"
        },
        "kind": "CANCELLED",
        "fillableTakerAssetAmount": 0,
        "txHash": "0x9e6830a7044b39e107f410e4f765995fd0d3d69d5c3b3582a1701b9d68167560"
      }
    ]
  }
}
```

See the [OrderEvent](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#OrderEvent) type declaration as well as the [OrderEventKind](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#pkg-constants) event types for a complete list of the events that could be emitted.

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
