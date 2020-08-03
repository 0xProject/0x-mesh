# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

-   **Mesh**

### Constructors

## constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): _[Mesh](#class-mesh)_

_Defined in [mesh.ts:132](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L132)_

Instantiates a new Mesh instance.

**Parameters:**

| Name     | Type                        | Description                    |
| -------- | --------------------------- | ------------------------------ |
| `config` | [Config](#interface-config) | Configuration options for Mesh |

**Returns:** _[Mesh](#class-mesh)_

An instance of Mesh

### Methods

## addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[], `pinned`: boolean): _Promise‹[ValidationResults](#interface-validationresults)›_

_Defined in [mesh.ts:275](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L275)_

Validates and adds the given orders to Mesh. If an order is successfully
added, Mesh will share it with any peers in the network and start
watching it for changes (e.g. filled, canceled, expired). The returned
promise will only be rejected if there was an error validating or adding
the order; it will not be rejected for any invalid orders (check
results.rejected instead).

**Parameters:**

| Name     | Type          | Default | Description                                                                                                                                                                                      |
| -------- | ------------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `orders` | SignedOrder[] | -       | An array of orders to add.                                                                                                                                                                       |
| `pinned` | boolean       | true    | Whether or not the orders should be pinned. Pinned orders will not be affected by any DDoS prevention or incentive mechanisms and will always stay in storage until they are no longer fillable. |

**Returns:** _Promise‹[ValidationResults](#interface-validationresults)›_

Validation results for the given orders, indicating which orders
were accepted and which were rejected.

---

## getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

_Defined in [mesh.ts:211](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L211)_

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

| Name      | Type   | Default | Description                                           |
| --------- | ------ | ------- | ----------------------------------------------------- |
| `perPage` | number | 200     | number of signedOrders to fetch per paginated request |

**Returns:** _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts

---

## getOrdersForPageAsync

▸ **getOrdersForPageAsync**(`perPage`: number, `minOrderHash?`: undefined | string): _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

_Defined in [mesh.ts:246](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L246)_

Get page of 0x signed orders stored on the Mesh node at the specified snapshot

**Parameters:**

| Name            | Type                    | Description                                                                                                      |
| --------------- | ----------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `perPage`       | number                  | Number of signedOrders to fetch per paginated request                                                            |
| `minOrderHash?` | undefined &#124; string | The minimum order hash for the returned orders. Should be set based on the last hash from the previous response. |

**Returns:** _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

Up to perPage orders with hash greater than minOrderHash, including order hashes and fillableTakerAssetAmounts

---

## getStatsAsync

▸ **getStatsAsync**(): _Promise‹[Stats](#interface-stats)›_

_Defined in [mesh.ts:194](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L194)_

Returns various stats about Mesh, including the total number of orders
and the number of peers Mesh is connected to.

**Returns:** _Promise‹[Stats](#interface-stats)›_

---

## onError

▸ **onError**(`handler`: function): _void_

_Defined in [mesh.ts:152](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L152)_

Registers a handler which will be called in the event of a critical
error. Note that the handler will not be called for non-critical errors.
In order to ensure no errors are missed, this should be called before
startAsync.

**Parameters:**

▪ **handler**: _function_

The handler to be called.

▸ (`err`: Error): _void_

**Parameters:**

| Name  | Type  |
| ----- | ----- |
| `err` | Error |

**Returns:** _void_

---

## onOrderEvents

▸ **onOrderEvents**(`handler`: function): _void_

_Defined in [mesh.ts:167](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L167)_

Registers a handler which will be called for any incoming order events.
Order events are fired whenver an order is added, canceled, expired, or
filled. In order to ensure no events are missed, this should be called
before startAsync.

**Parameters:**

▪ **handler**: _function_

The handler to be called.

▸ (`events`: [OrderEvent](#interface-orderevent)[]): _void_

**Parameters:**

| Name     | Type                                  |
| -------- | ------------------------------------- |
| `events` | [OrderEvent](#interface-orderevent)[] |

**Returns:** _void_

---

## startAsync

▸ **startAsync**(): _Promise‹void›_

_Defined in [mesh.ts:178](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/mesh.ts#L178)_

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** _Promise‹void›_

<hr />

# Enumeration: ContractEventKind

### Enumeration members

## ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

_Defined in [types.ts:470](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L470)_

---

## ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in [types.ts:472](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L472)_

---

## ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in [types.ts:471](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L471)_

---

## ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in [types.ts:466](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L466)_

---

## ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in [types.ts:465](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L465)_

---

## ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in [types.ts:468](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L468)_

---

## ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in [types.ts:469](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L469)_

---

## ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in [types.ts:467](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L467)_

---

## ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in [types.ts:474](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L474)_

---

## ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in [types.ts:475](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L475)_

---

## ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in [types.ts:473](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L473)_

---

## WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in [types.ts:476](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L476)_

---

## WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in [types.ts:477](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L477)_

<hr />

# Enumeration: OrderEventEndState

### Enumeration members

## Added

• **Added**: = "ADDED"

_Defined in [types.ts:540](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L540)_

---

## Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in [types.ts:543](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L543)_

---

## Expired

• **Expired**: = "EXPIRED"

_Defined in [types.ts:544](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L544)_

---

## FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

_Defined in [types.ts:547](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L547)_

---

## Filled

• **Filled**: = "FILLED"

_Defined in [types.ts:541](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L541)_

---

## FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

_Defined in [types.ts:542](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L542)_

---

## Invalid

• **Invalid**: = "INVALID"

_Defined in [types.ts:539](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L539)_

---

## StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

_Defined in [types.ts:548](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L548)_

---

## Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in [types.ts:545](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L545)_

---

## Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in [types.ts:546](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L546)_

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.

### Enumeration members

## CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

_Defined in [types.ts:632](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L632)_

---

## MeshError

• **MeshError**: = "MESH_ERROR"

_Defined in [types.ts:630](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L630)_

---

## MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

_Defined in [types.ts:631](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L631)_

---

## ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

_Defined in [types.ts:629](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L629)_

<hr />

# Enumeration: Verbosity

### Enumeration members

## Debug

• **Debug**: = 5

_Defined in [types.ts:212](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L212)_

---

## Error

• **Error**: = 2

_Defined in [types.ts:209](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L209)_

---

## Fatal

• **Fatal**: = 1

_Defined in [types.ts:208](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L208)_

---

## Info

• **Info**: = 4

_Defined in [types.ts:211](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L211)_

---

## Panic

• **Panic**: = 0

_Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L207)_

---

## Trace

• **Trace**: = 6

_Defined in [types.ts:213](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L213)_

---

## Warn

• **Warn**: = 3

_Defined in [types.ts:210](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L210)_

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

-   **AcceptedOrderInfo**

### Properties

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:610](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L610)_

---

## isNew

• **isNew**: _boolean_

_Defined in [types.ts:611](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L611)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:608](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L608)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:609](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L609)_

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

-   **Config**

### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : _undefined | number_

_Defined in [types.ts:116](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L116)_

---

## `Optional` bootstrapList

• **bootstrapList**? : _string[]_

_Defined in [types.ts:109](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L109)_

---

## `Optional` customContractAddresses

• **customContractAddresses**? : _[ContractAddresses](#interface-contractaddresses)_

_Defined in [types.ts:160](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L160)_

---

## `Optional` customOrderFilter

• **customOrderFilter**? : _[JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L185)_

---

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : _undefined | false | true_

_Defined in [types.ts:133](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L133)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [types.ts:101](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L101)_

---

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : _undefined | number_

_Defined in [types.ts:125](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L125)_

---

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : _undefined | number_

_Defined in [types.ts:138](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L138)_

---

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : _undefined | number_

_Defined in [types.ts:144](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L144)_

---

## `Optional` ethereumRPCURL

• **ethereumRPCURL**? : _undefined | string_

_Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L98)_

---

## `Optional` maxBytesPerSecond

• **maxBytesPerSecond**? : _undefined | number_

_Defined in [types.ts:191](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L191)_

---

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : _undefined | number_

_Defined in [types.ts:165](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L165)_

---

## `Optional` useBootstrapList

• **useBootstrapList**? : _undefined | false | true_

_Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L104)_

---

## `Optional` verbosity

• **verbosity**? : _[Verbosity](#enumeration-verbosity)_

_Defined in [types.ts:95](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L95)_

---

## `Optional` web3Provider

• **web3Provider**? : _SupportedProvider_

_Defined in [types.ts:188](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L188)_

<hr />

# Interface: ContractAddresses

### Hierarchy

-   **ContractAddresses**

### Properties

## `Optional` coordinator

• **coordinator**? : _undefined | string_

_Defined in [types.ts:200](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L200)_

---

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : _undefined | string_

_Defined in [types.ts:201](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L201)_

---

## devUtils

• **devUtils**: _string_

_Defined in [types.ts:196](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L196)_

---

## erc1155Proxy

• **erc1155Proxy**: _string_

_Defined in [types.ts:199](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L199)_

---

## erc20Proxy

• **erc20Proxy**: _string_

_Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L197)_

---

## erc721Proxy

• **erc721Proxy**: _string_

_Defined in [types.ts:198](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L198)_

---

## exchange

• **exchange**: _string_

_Defined in [types.ts:195](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L195)_

---

## `Optional` weth9

• **weth9**? : _undefined | string_

_Defined in [types.ts:202](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L202)_

---

## `Optional` zrxToken

• **zrxToken**? : _undefined | string_

_Defined in [types.ts:203](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L203)_

<hr />

# Interface: ContractEvent

### Hierarchy

-   **ContractEvent**

### Properties

## address

• **address**: _string_

_Defined in [types.ts:518](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L518)_

---

## blockHash

• **blockHash**: _string_

_Defined in [types.ts:513](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L513)_

---

## isRemoved

• **isRemoved**: _boolean_

_Defined in [types.ts:517](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L517)_

---

## kind

• **kind**: _[ContractEventKind](#enumeration-contracteventkind)_

_Defined in [types.ts:519](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L519)_

---

## logIndex

• **logIndex**: _number_

_Defined in [types.ts:516](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L516)_

---

## parameters

• **parameters**: _ContractEventParameters_

_Defined in [types.ts:520](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L520)_

---

## txHash

• **txHash**: _string_

_Defined in [types.ts:514](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L514)_

---

## txIndex

• **txIndex**: _number_

_Defined in [types.ts:515](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L515)_

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

-   **ERC1155ApprovalForAllEvent**

### Properties

## approved

• **approved**: _boolean_

_Defined in [types.ts:382](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L382)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:381](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L381)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:380](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L380)_

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

-   **ERC1155TransferBatchEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L364)_

---

## ids

• **ids**: _BigNumber[]_

_Defined in [types.ts:366](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L366)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L363)_

---

## to

• **to**: _string_

_Defined in [types.ts:365](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L365)_

---

## values

• **values**: _BigNumber[]_

_Defined in [types.ts:367](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L367)_

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

-   **ERC1155TransferSingleEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:347](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L347)_

---

## id

• **id**: _BigNumber_

_Defined in [types.ts:349](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L349)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L346)_

---

## to

• **to**: _string_

_Defined in [types.ts:348](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L348)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:350](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L350)_

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

-   **ERC20ApprovalEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:301](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L301)_

---

## spender

• **spender**: _string_

_Defined in [types.ts:302](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L302)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:303](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L303)_

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

-   **ERC20TransferEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:288](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L288)_

---

## to

• **to**: _string_

_Defined in [types.ts:289](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L289)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:290](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L290)_

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

-   **ERC721ApprovalEvent**

### Properties

## approved

• **approved**: _string_

_Defined in [types.ts:328](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L328)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:327](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L327)_

---

## tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:329](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L329)_

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

-   **ERC721ApprovalForAllEvent**

### Properties

## approved

• **approved**: _boolean_

_Defined in [types.ts:342](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L342)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:341](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L341)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:340](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L340)_

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

-   **ERC721TransferEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:314](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L314)_

---

## to

• **to**: _string_

_Defined in [types.ts:315](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L315)_

---

## tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:316](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L316)_

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

-   **ExchangeCancelEvent**

### Properties

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:423](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L423)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:421](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L421)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:425](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L425)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L424)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:422](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L422)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:426](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L426)_

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

-   **ExchangeCancelUpToEvent**

### Properties

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:430](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L430)_

---

## orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in [types.ts:432](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L432)_

---

## orderSenderAddress

• **orderSenderAddress**: _string_

_Defined in [types.ts:431](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L431)_

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

-   **ExchangeFillEvent**

### Properties

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:389](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L389)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:386](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L386)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:396](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L396)_

---

## makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:390](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L390)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Defined in [types.ts:398](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L398)_

---

## makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in [types.ts:392](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L392)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:395](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L395)_

---

## protocolFeePaid

• **protocolFeePaid**: _BigNumber_

_Defined in [types.ts:394](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L394)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:388](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L388)_

---

## takerAddress

• **takerAddress**: _string_

_Defined in [types.ts:387](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L387)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:397](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L397)_

---

## takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:391](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L391)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Defined in [types.ts:399](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L399)_

---

## takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in [types.ts:393](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L393)_

<hr />

# Interface: GetOrdersResponse

### Hierarchy

-   **GetOrdersResponse**

### Properties

## ordersInfos

• **ordersInfos**: _[OrderInfo](#interface-orderinfo)[]_

_Defined in [types.ts:17](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L17)_

---

## timestamp

• **timestamp**: _number_

_Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L16)_

<hr />

# Interface: JsonSchema

An interface for JSON schema types, which are used for custom order filters.

### Hierarchy

-   **JsonSchema**

### Properties

## `Optional` \$ref

• **\$ref**? : _undefined | string_

_Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L39)_

---

## `Optional` \$schema

• **\$schema**? : _undefined | string_

_Defined in [types.ts:38](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L38)_

---

## `Optional` additionalItems

• **additionalItems**? : _boolean | [JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L50)_

---

## `Optional` additionalProperties

• **additionalProperties**? : _boolean | [JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:58](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L58)_

---

## `Optional` allOf

• **allOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L80)_

---

## `Optional` anyOf

• **anyOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L81)_

---

## `Optional` const

• **const**? : _any_

_Defined in [types.ts:77](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L77)_

---

## `Optional` definitions

• **definitions**? : _undefined | object_

_Defined in [types.ts:59](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L59)_

---

## `Optional` dependencies

• **dependencies**? : _undefined | object_

_Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L68)_

---

## `Optional` description

• **description**? : _undefined | string_

_Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L41)_

---

## `Optional` enum

• **enum**? : _any[]_

_Defined in [types.ts:71](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L71)_

---

## `Optional` exclusiveMaximum

• **exclusiveMaximum**? : _undefined | false | true_

_Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L44)_

---

## `Optional` exclusiveMinimum

• **exclusiveMinimum**? : _undefined | false | true_

_Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L46)_

---

## `Optional` format

• **format**? : _undefined | string_

_Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L79)_

---

## `Optional` id

• **id**? : _undefined | string_

_Defined in [types.ts:37](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L37)_

---

## `Optional` items

• **items**? : _[JsonSchema](#interface-jsonschema) | [JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:51](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L51)_

---

## `Optional` maxItems

• **maxItems**? : _undefined | number_

_Defined in [types.ts:52](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L52)_

---

## `Optional` maxLength

• **maxLength**? : _undefined | number_

_Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L47)_

---

## `Optional` maxProperties

• **maxProperties**? : _undefined | number_

_Defined in [types.ts:55](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L55)_

---

## `Optional` maximum

• **maximum**? : _undefined | number_

_Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L43)_

---

## `Optional` minItems

• **minItems**? : _undefined | number_

_Defined in [types.ts:53](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L53)_

---

## `Optional` minLength

• **minLength**? : _undefined | number_

_Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L48)_

---

## `Optional` minProperties

• **minProperties**? : _undefined | number_

_Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L56)_

---

## `Optional` minimum

• **minimum**? : _undefined | number_

_Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L45)_

---

## `Optional` multipleOf

• **multipleOf**? : _undefined | number_

_Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L42)_

---

## `Optional` not

• **not**? : _[JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:83](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L83)_

---

## `Optional` oneOf

• **oneOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:82](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L82)_

---

## `Optional` pattern

• **pattern**? : _string | RegExp_

_Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L49)_

---

## `Optional` patternProperties

• **patternProperties**? : _undefined | object_

_Defined in [types.ts:65](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L65)_

---

## `Optional` properties

• **properties**? : _undefined | object_

_Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L62)_

---

## `Optional` required

• **required**? : _string[]_

_Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L57)_

---

## `Optional` title

• **title**? : _undefined | string_

_Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L40)_

---

## `Optional` type

• **type**? : _string | string[]_

_Defined in [types.ts:78](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L78)_

---

## `Optional` uniqueItems

• **uniqueItems**? : _undefined | false | true_

_Defined in [types.ts:54](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L54)_

<hr />

# Interface: LatestBlock

### Hierarchy

-   **LatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [types.ts:645](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L645)_

---

## number

• **number**: _number_

_Defined in [types.ts:644](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L644)_

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

-   **OrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [types.ts:571](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L571)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [types.ts:569](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L569)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:570](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L570)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:567](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L567)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:568](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L568)_

---

## timestampMs

• **timestampMs**: _number_

_Defined in [types.ts:566](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L566)_

<hr />

# Interface: OrderInfo

### Hierarchy

-   **OrderInfo**

### Properties

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L30)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:28](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L28)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:29](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L29)_

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

-   **RejectedOrderInfo**

### Properties

## kind

• **kind**: _[RejectedOrderKind](#enumeration-rejectedorderkind)_

_Defined in [types.ts:621](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L621)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:619](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L619)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:620](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L620)_

---

## status

• **status**: _[RejectedOrderStatus](#interface-rejectedorderstatus)_

_Defined in [types.ts:622](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L622)_

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

-   **RejectedOrderStatus**

### Properties

## code

• **code**: _string_

_Defined in [types.ts:639](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L639)_

---

## message

• **message**: _string_

_Defined in [types.ts:640](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L640)_

<hr />

# Interface: Stats

### Hierarchy

-   **Stats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [types.ts:682](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L682)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [types.ts:681](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L681)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [types.ts:673](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L673)_

---

## latestBlock

• **latestBlock**: _[LatestBlock](#interface-latestblock)_

_Defined in [types.ts:674](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L674)_

---

## maxExpirationTime

• **maxExpirationTime**: _BigNumber_

_Defined in [types.ts:679](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L679)_

---

## numOrders

• **numOrders**: _number_

_Defined in [types.ts:676](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L676)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [types.ts:677](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L677)_

---

## numPeers

• **numPeers**: _number_

_Defined in [types.ts:675](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L675)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [types.ts:678](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L678)_

---

## peerID

• **peerID**: _string_

_Defined in [types.ts:672](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L672)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [types.ts:669](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L669)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [types.ts:670](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L670)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [types.ts:671](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L671)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _Date_

_Defined in [types.ts:680](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L680)_

---

## version

• **version**: _string_

_Defined in [types.ts:668](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L668)_

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

-   **ValidationResults**

### Properties

## accepted

• **accepted**: _[AcceptedOrderInfo](#interface-acceptedorderinfo)[]_

_Defined in [types.ts:600](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L600)_

---

## rejected

• **rejected**: _[RejectedOrderInfo](#interface-rejectedorderinfo)[]_

_Defined in [types.ts:601](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L601)_

<hr />

# Interface: WethDepositEvent

### Hierarchy

-   **WethDepositEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:454](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L454)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:455](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L455)_

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

-   **WethWithdrawalEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:443](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L443)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:444](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/types.ts#L444)_

<hr />

# Functions

## loadMeshStreamingForURLAsync

▸ **loadMeshStreamingWithURLAsync**(`url`: `string`): _Promise‹`void`›_

_Defined in [index.ts:7](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/index.ts#L7)_

Loads the Wasm module that is provided by fetching a url.

**Parameters:**

| Name  | Type     | Description                          |
| ----- | -------- | ------------------------------------ |
| `url` | `string` | The URL to query for the Wasm binary |

<hr />

## loadMeshStreamingAsync

▸ **loadMeshStreamingAsync**(`response`: `Response | Promise<Response>`): _Promise‹`void`›_

_Defined in [index.ts:15](https://github.com/0xProject/0x-mesh/blob/1f7ab983/packages/mesh-browser-lite/src/index.ts#L15)_

Loads the Wasm module that is provided by a response.

**Parameters:**

| Name       | Type                                | Description                                     |
| ---------- | ----------------------------------- | ----------------------------------------------- |
| `response` | `Response &#124; Promise<Response>` | The Wasm response that supplies the Wasm binary |

<hr />
