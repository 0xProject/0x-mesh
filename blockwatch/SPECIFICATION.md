# Mesh Technical Design Doc

This doc outlines the high-level technical design specifications of 0x Mesh.

## Packages

Proposed package hierarchy and organization within the `0x-mesh` project.

- `db` -- Contains all the logic related to interfacing with IndexDB & Redis databases
  - `db.go` -- Contains the key-value database interface
  - `indexdb.go` -- Contains the IndexDB-backed implementation of the db interface
  - `redis.go` -- Contains the Redis-backed implementation of the db interface
- `core` -- Contains the core ETH-staked messaging logic
  - `core.go` -- Allows for a core node to be instantiated
  - `networking` -- Contains all p2p networking logic (TODO: Flesh this out more)
  - `ethwatch` -- Enables watching ETH balances for changes
  - `sharing` -- Contains the logic for sharing messages with peers
  - `reputation` -- Manages the reputations of neighbors
  - `messages` -- Manages the messages DB (`message` is defined as [payload, signature])
- `ethereum` -- Contains all logic that is Ethereum-specific.
  - `blockwatch` -- Enables watching Ethereum blocks in a block reorg friendly way
  - `eventwatch` -- Enables contract event watching in a block reorg friendly way
- `0x` -- Contains all 0x-specific logic
  - `orderwatch` -- Enables watching a set of orders for order-relevant state changes
    - `orderwatch.go` -- See above
    - `expirationwatch.go` -- Enables watching order expirations
    - `cleanup.go` -- Periodically updates all order-relevant state. Fixes any missed state updates
  - `node` -- A Mesh node that wraps a generic ETH-staked messaging Mesh node into a 0x-specific Mesh node.
    - `parser.go` -- Parses p2p messages into 0x orders and associated assetData
    - `validator.go` -- Validates incoming p2p messages to ensure that they are valid 0x orders
    - `ws` -- WebSocket server implementation
    - `http` -- Http server with order submission, diagnostics endpoints & potentially more

## DB

We will be using a key-value store database for 0x Mesh. Since the node runs both in the browser and server, we will use `IndexDB` in the browser, and `Redis` when running on a server. The differences between the interfaces of these two databases will be abstracted behind a `Database` interface.

### Database interface

```go
type Database interface {
    // TODO: Research Redis and see if a common interface is possible.
}
```

The rest of the node's logic will only interact with the database via a module that adheres to this interface.

## Core

### Networking

TODO(albrow)

### EthWatch

Watches for changes in ETH balances of interest and emits events whenever they change.

### Sharing

TODO(albrow)

### Reputation

TODO(albrow)

### Messages

TODO(fabio): Message schema / encoding / validation

### Core

Defines the core node and is the entry-point for interfacing with it.

```go
const (
    Stored EventType = iota
    Evicted
)

type Event type {
    Type EventType
    Message *Message
}

type Node interface {
    ValidationFn func(message *Message) bool
    Events chan*Event
    EvictionRequests chan*EvictionRequest
}
```

#### Database

##### Tables

###### Message

- messageHash (key) [messageHash == orderHash for 0x Mesh]
- message [most compact representation of message]
- stakingAddress (index) [stakingAddress == makerAddress for 0x Mesh]
- stakingSignature

```
    messageHash => {message: '...', stakingAddress: '...'}
```

###### Staker

- address (key)
- ethBalance

```
    address => ethBalance
```

##### DB Client

```go
type Message struct {
    encodedMessage []byte
    stakingAddress string
    stakingSignature string
}

type Staker struct {
    address string
    ethBalance *big.Int
}

type coreDb interface {
    // Initializes the database connection, and creates the tables if don't already exist
    init() *DB.conn

    // Message operations
    insertMessage(encodedMessage []byte, stakerAddress string, stakerSignature: string) (Message, error)
    removeMessage(messageHash string) error
    getRandomMessages(limit uint) ([]Message, error)
    removeMessagesWithLeastStake(limit uint) error

    // Staker operations
    upsertStaker(address string, ethBalance *big.Int) error
    upsertStakers(addresses []string, ethBalancees []*big.Int) (error)
    getAllStakerAddresses() ([]string, error)
    removeStaker(address string) bool
}
```

## Ethereum

##### BlockWatch

Watches for new Ethereum blocks and handles block-reorgs, network requests. Emits a steady stream of block additionals and removals.

##### EventWatch

Consumes the block feed from BlockWatch and emits the corresponding log events, filtered to those of interest to the caller.

## 0x

### OrderWatch

**Core assumptions**

1. The entity running the OrderWatcher already has a data-store containing orders & orderRelevantState
   (e.g Relayers need this to display the orderbook, as well as a users allowance/balance in a UI)
1. We want to persist state caches & orders, resuming the OrderWatcher where it left off upon restart
1. We want to support as much of 0x's functionality as possible (e.g., all assetProxies, relayer models (i.e., open orderbook & matching), etc...)

#### High-level design

##### OrderWatch

Instantiates and orchestrates all the other components in this module. Acts as a single point of entry where 0x orders are fed in, and where order state changed events are emitted. Internally it instantiates an `eventwatch.Watcher` to watch the following events: ERC20 `Transfer`, ERC721 `Transfer`, ERC20 `Approve`, ERC721 `Approve`, ERC721 `ApproveForAll`, WETH `Deposit`, WETH `Withdraw`, Exchange `Fill`, Exchange `Cancel`, Exchange `CancelUpTo`. It also instantiates `ExpirationWatch` and `Cleanup`

##### ExpirationWatch

Maintains a red-black tree of order hashes sorted by expiration timestamp. Emits events whenever an order expires.

##### Cleanup

Unfortunately watching contract events is not 100% reliable. We therefore need to re-validate all orders on some interval, in order to prune any unfillable orders that experienced a state-change not caught by the event watcher.

##### Database

In order to persistently store orders and order-relevant state (e.g., trader balances, allowances, etc.) across re-starts, we require a database. This database will be constructed using the 0x Mesh key-value store interface that abstracts away the implementation differences between browser/server DB implementations.

##### Tables

###### Order

- Hash (key)
- expirationTimeSeconds
- makerAddress + makerAssetData (index)
- makerAssetAmount
- makerFee
- takerAddress + takerAssetData (index)
- takerAssetAmount
- takerFee

```
    hash => {...order}
```

###### AssetData

- singleAssetData (key)
- multiAssetDatas (array)
- tokenAddress
- tokenId?

```
    singleAssetData => {tokenAddress: '...', tokenId: '...', multiAssetDatas: [...multiAssetDatas]}
```

###### TraderInfo

- id [address + singleAssetData](key)
- balance
- allowance

```
    [address + signedAssetData] => {balance: '...', allowance: '...'}
```

##### DB Client

An additional layer of abstraction might be desired, one that does not lock OrderWatcher instantiators to even a key-value data store. We can abstract all DB queries into a DB client that is passed into OrderWatcher at instantiation.

```go
type Order struct {
    hash string
    expiration time.Time
    makerAddress string
    makerAssetData string
    makerAssetAmount *big.Int
    makerFee *big.Int
    takerAddress string
    takerAssetData string
    takerAssetAmount *big.Int
    takerFee *big.Int
    takerAssetFilledAmount *big.Int
    isCancelled bool
}

type TraderInfo struct {
    balance *big.Int
    allowance *big.Int
}

type DbClient interface {
    // Initializes the database connection, and creates the tables if don't already exist
    Init()

    // Order operations
    InsertOrder(order Order) error
    GetAllOrderHashesAndExpirations() ([]string, []time.Time, error)
    GetOrderByHash(hash string) (Order, error)
    GetOrdersByMakerInfo(makerAddress string, makerAssetData string) ([]Order, error)
    GetOrdersByTakerInfo(takerAddress string, takerAssetData string) ([]Order, error)
    UpdateOrderState(hash string, takerAssetFilledAmount *big.Int, isCancelled bool) error

    // AssetData operations
    InsertSingleERC721Asset(
        singleAssetData string,
        tokenAddress string,
        tokenId string,
        multiAssetDatas []string
    ) error
    InsertSingleERC20Asset(singleAssetData string, tokenAddress string, multiAssetDatas []string) error
    RemoveSingleAssetData(signleAssetData string) error
    AddMultiAssetData(singleAssetData string, multiAssetData string) error
    RemoveMultiAssetData(singleAssetData string, multiAssetData string) error
    GetMultiAssetDatas(singleAssetData string) ([]string, error)
    GetAssetDatasByTokenAddress(tokenId string) ([]string, error) // Returns both the single & multi assetDatas

    // TraderInfo operations
    UpsertTraderInfo(traderAddress string, singleAssetData string, traderInfo TraderInfo): error
    GetTraderInfo(traderAddress string, singleAssetData string) (TraderInfo, error)
    RemoveTraderInfo(traderAddress string) error
}
```

### Node

Creates a 0x-specific Mesh node by wrapping the core node and enforcing that messages contain orders and that messages including invalid orders get evicted. It also defines the public interface for the 0x Mesh node, including a WebSocket and HTTP server.

#### Validator

TODO(fabio)

#### Parser

TODO(fabio)

#### WS

TODO(fabio)

#### HTTP

TODO(fabio)
