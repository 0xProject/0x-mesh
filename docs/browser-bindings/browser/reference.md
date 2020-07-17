# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

-   **Mesh**

### Constructors

## constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): _[Mesh](#class-mesh)_

_Defined in [mesh.ts:132](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L132)_

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

_Defined in [mesh.ts:275](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L275)_

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

_Defined in [mesh.ts:211](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L211)_

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

_Defined in [mesh.ts:246](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L246)_

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

_Defined in [mesh.ts:194](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L194)_

Returns various stats about Mesh, including the total number of orders
and the number of peers Mesh is connected to.

**Returns:** _Promise‹[Stats](#interface-stats)›_

---

## onError

▸ **onError**(`handler`: function): _void_

_Defined in [mesh.ts:152](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L152)_

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

_Defined in [mesh.ts:167](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L167)_

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

_Defined in [mesh.ts:178](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/mesh.ts#L178)_

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** _Promise‹void›_

<hr />

# Enumeration: ContractEventKind

### Enumeration members

## ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

_Defined in [types.ts:466](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L466)_

---

## ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in [types.ts:468](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L468)_

---

## ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in [types.ts:467](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L467)_

---

## ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in [types.ts:462](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L462)_

---

## ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in [types.ts:461](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L461)_

---

## ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in [types.ts:464](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L464)_

---

## ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in [types.ts:465](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L465)_

---

## ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in [types.ts:463](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L463)_

---

## ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in [types.ts:470](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L470)_

---

## ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in [types.ts:471](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L471)_

---

## ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in [types.ts:469](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L469)_

---

## WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in [types.ts:472](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L472)_

---

## WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in [types.ts:473](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L473)_

<hr />

# Enumeration: OrderEventEndState

### Enumeration members

## Added

• **Added**: = "ADDED"

_Defined in [types.ts:536](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L536)_

---

## Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in [types.ts:539](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L539)_

---

## Expired

• **Expired**: = "EXPIRED"

_Defined in [types.ts:540](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L540)_

---

## FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

_Defined in [types.ts:543](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L543)_

---

## Filled

• **Filled**: = "FILLED"

_Defined in [types.ts:537](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L537)_

---

## FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

_Defined in [types.ts:538](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L538)_

---

## Invalid

• **Invalid**: = "INVALID"

_Defined in [types.ts:535](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L535)_

---

## StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

_Defined in [types.ts:544](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L544)_

---

## Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in [types.ts:541](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L541)_

---

## Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in [types.ts:542](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L542)_

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.

### Enumeration members

## CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

_Defined in [types.ts:628](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L628)_

---

## MeshError

• **MeshError**: = "MESH_ERROR"

_Defined in [types.ts:626](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L626)_

---

## MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

_Defined in [types.ts:627](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L627)_

---

## ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

_Defined in [types.ts:625](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L625)_

<hr />

# Enumeration: Verbosity

### Enumeration members

## Debug

• **Debug**: = 5

_Defined in [types.ts:209](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L209)_

---

## Error

• **Error**: = 2

_Defined in [types.ts:206](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L206)_

---

## Fatal

• **Fatal**: = 1

_Defined in [types.ts:205](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L205)_

---

## Info

• **Info**: = 4

_Defined in [types.ts:208](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L208)_

---

## Panic

• **Panic**: = 0

_Defined in [types.ts:204](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L204)_

---

## Trace

• **Trace**: = 6

_Defined in [types.ts:210](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L210)_

---

## Warn

• **Warn**: = 3

_Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L207)_

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

-   **AcceptedOrderInfo**

### Properties

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:606](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L606)_

---

## isNew

• **isNew**: _boolean_

_Defined in [types.ts:607](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L607)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:604](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L604)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:605](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L605)_

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

-   **Config**

### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : _undefined | number_

_Defined in [types.ts:116](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L116)_

---

## `Optional` bootstrapList

• **bootstrapList**? : _string[]_

_Defined in [types.ts:109](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L109)_

---

## `Optional` customContractAddresses

• **customContractAddresses**? : _[ContractAddresses](#interface-contractaddresses)_

_Defined in [types.ts:160](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L160)_

---

## `Optional` customOrderFilter

• **customOrderFilter**? : _[JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L185)_

---

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : _undefined | false | true_

_Defined in [types.ts:133](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L133)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [types.ts:101](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L101)_

---

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : _undefined | number_

_Defined in [types.ts:125](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L125)_

---

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : _undefined | number_

_Defined in [types.ts:138](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L138)_

---

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : _undefined | number_

_Defined in [types.ts:144](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L144)_

---

## `Optional` ethereumRPCURL

• **ethereumRPCURL**? : _undefined | string_

_Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L98)_

---

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : _undefined | number_

_Defined in [types.ts:165](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L165)_

---

## `Optional` useBootstrapList

• **useBootstrapList**? : _undefined | false | true_

_Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L104)_

---

## `Optional` verbosity

• **verbosity**? : _[Verbosity](#enumeration-verbosity)_

_Defined in [types.ts:95](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L95)_

---

## `Optional` web3Provider

• **web3Provider**? : _SupportedProvider_

_Defined in [types.ts:188](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L188)_

<hr />

# Interface: ContractAddresses

### Hierarchy

-   **ContractAddresses**

### Properties

## `Optional` coordinator

• **coordinator**? : _undefined | string_

_Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L197)_

---

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : _undefined | string_

_Defined in [types.ts:198](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L198)_

---

## devUtils

• **devUtils**: _string_

_Defined in [types.ts:193](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L193)_

---

## erc1155Proxy

• **erc1155Proxy**: _string_

_Defined in [types.ts:196](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L196)_

---

## erc20Proxy

• **erc20Proxy**: _string_

_Defined in [types.ts:194](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L194)_

---

## erc721Proxy

• **erc721Proxy**: _string_

_Defined in [types.ts:195](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L195)_

---

## exchange

• **exchange**: _string_

_Defined in [types.ts:192](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L192)_

---

## `Optional` weth9

• **weth9**? : _undefined | string_

_Defined in [types.ts:199](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L199)_

---

## `Optional` zrxToken

• **zrxToken**? : _undefined | string_

_Defined in [types.ts:200](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L200)_

<hr />

# Interface: ContractEvent

### Hierarchy

-   **ContractEvent**

### Properties

## address

• **address**: _string_

_Defined in [types.ts:514](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L514)_

---

## blockHash

• **blockHash**: _string_

_Defined in [types.ts:509](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L509)_

---

## isRemoved

• **isRemoved**: _boolean_

_Defined in [types.ts:513](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L513)_

---

## kind

• **kind**: _[ContractEventKind](#enumeration-contracteventkind)_

_Defined in [types.ts:515](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L515)_

---

## logIndex

• **logIndex**: _number_

_Defined in [types.ts:512](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L512)_

---

## parameters

• **parameters**: _ContractEventParameters_

_Defined in [types.ts:516](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L516)_

---

## txHash

• **txHash**: _string_

_Defined in [types.ts:510](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L510)_

---

## txIndex

• **txIndex**: _number_

_Defined in [types.ts:511](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L511)_

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

-   **ERC1155ApprovalForAllEvent**

### Properties

## approved

• **approved**: _boolean_

_Defined in [types.ts:378](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L378)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:377](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L377)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:376](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L376)_

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

-   **ERC1155TransferBatchEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:360](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L360)_

---

## ids

• **ids**: _BigNumber[]_

_Defined in [types.ts:362](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L362)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:359](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L359)_

---

## to

• **to**: _string_

_Defined in [types.ts:361](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L361)_

---

## values

• **values**: _BigNumber[]_

_Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L363)_

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

-   **ERC1155TransferSingleEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:343](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L343)_

---

## id

• **id**: _BigNumber_

_Defined in [types.ts:345](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L345)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:342](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L342)_

---

## to

• **to**: _string_

_Defined in [types.ts:344](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L344)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L346)_

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

-   **ERC20ApprovalEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:297](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L297)_

---

## spender

• **spender**: _string_

_Defined in [types.ts:298](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L298)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:299](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L299)_

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

-   **ERC20TransferEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:284](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L284)_

---

## to

• **to**: _string_

_Defined in [types.ts:285](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L285)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L286)_

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

-   **ERC721ApprovalEvent**

### Properties

## approved

• **approved**: _string_

_Defined in [types.ts:324](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L324)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:323](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L323)_

---

## tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:325](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L325)_

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

-   **ERC721ApprovalForAllEvent**

### Properties

## approved

• **approved**: _boolean_

_Defined in [types.ts:338](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L338)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:337](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L337)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:336](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L336)_

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

-   **ERC721TransferEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:310](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L310)_

---

## to

• **to**: _string_

_Defined in [types.ts:311](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L311)_

---

## tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:312](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L312)_

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

-   **ExchangeCancelEvent**

### Properties

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:419](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L419)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:417](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L417)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:421](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L421)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:420](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L420)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:418](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L418)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:422](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L422)_

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

-   **ExchangeCancelUpToEvent**

### Properties

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:426](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L426)_

---

## orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in [types.ts:428](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L428)_

---

## orderSenderAddress

• **orderSenderAddress**: _string_

_Defined in [types.ts:427](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L427)_

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

-   **ExchangeFillEvent**

### Properties

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:385](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L385)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:382](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L382)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:392](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L392)_

---

## makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:386](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L386)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Defined in [types.ts:394](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L394)_

---

## makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in [types.ts:388](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L388)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:391](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L391)_

---

## protocolFeePaid

• **protocolFeePaid**: _BigNumber_

_Defined in [types.ts:390](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L390)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:384](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L384)_

---

## takerAddress

• **takerAddress**: _string_

_Defined in [types.ts:383](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L383)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:393](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L393)_

---

## takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:387](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L387)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Defined in [types.ts:395](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L395)_

---

## takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in [types.ts:389](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L389)_

<hr />

# Interface: GetOrdersResponse

### Hierarchy

-   **GetOrdersResponse**

### Properties

## ordersInfos

• **ordersInfos**: _[OrderInfo](#interface-orderinfo)[]_

_Defined in [types.ts:17](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L17)_

---

## timestamp

• **timestamp**: _number_

_Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L16)_

<hr />

# Interface: JsonSchema

An interface for JSON schema types, which are used for custom order filters.

### Hierarchy

-   **JsonSchema**

### Properties

## `Optional` \$ref

• **\$ref**? : _undefined | string_

_Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L39)_

---

## `Optional` \$schema

• **\$schema**? : _undefined | string_

_Defined in [types.ts:38](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L38)_

---

## `Optional` additionalItems

• **additionalItems**? : _boolean | [JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L50)_

---

## `Optional` additionalProperties

• **additionalProperties**? : _boolean | [JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:58](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L58)_

---

## `Optional` allOf

• **allOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L80)_

---

## `Optional` anyOf

• **anyOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L81)_

---

## `Optional` const

• **const**? : _any_

_Defined in [types.ts:77](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L77)_

---

## `Optional` definitions

• **definitions**? : _undefined | object_

_Defined in [types.ts:59](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L59)_

---

## `Optional` dependencies

• **dependencies**? : _undefined | object_

_Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L68)_

---

## `Optional` description

• **description**? : _undefined | string_

_Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L41)_

---

## `Optional` enum

• **enum**? : _any[]_

_Defined in [types.ts:71](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L71)_

---

## `Optional` exclusiveMaximum

• **exclusiveMaximum**? : _undefined | false | true_

_Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L44)_

---

## `Optional` exclusiveMinimum

• **exclusiveMinimum**? : _undefined | false | true_

_Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L46)_

---

## `Optional` format

• **format**? : _undefined | string_

_Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L79)_

---

## `Optional` id

• **id**? : _undefined | string_

_Defined in [types.ts:37](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L37)_

---

## `Optional` items

• **items**? : _[JsonSchema](#interface-jsonschema) | [JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:51](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L51)_

---

## `Optional` maxItems

• **maxItems**? : _undefined | number_

_Defined in [types.ts:52](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L52)_

---

## `Optional` maxLength

• **maxLength**? : _undefined | number_

_Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L47)_

---

## `Optional` maxProperties

• **maxProperties**? : _undefined | number_

_Defined in [types.ts:55](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L55)_

---

## `Optional` maximum

• **maximum**? : _undefined | number_

_Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L43)_

---

## `Optional` minItems

• **minItems**? : _undefined | number_

_Defined in [types.ts:53](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L53)_

---

## `Optional` minLength

• **minLength**? : _undefined | number_

_Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L48)_

---

## `Optional` minProperties

• **minProperties**? : _undefined | number_

_Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L56)_

---

## `Optional` minimum

• **minimum**? : _undefined | number_

_Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L45)_

---

## `Optional` multipleOf

• **multipleOf**? : _undefined | number_

_Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L42)_

---

## `Optional` not

• **not**? : _[JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:83](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L83)_

---

## `Optional` oneOf

• **oneOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:82](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L82)_

---

## `Optional` pattern

• **pattern**? : _string | RegExp_

_Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L49)_

---

## `Optional` patternProperties

• **patternProperties**? : _undefined | object_

_Defined in [types.ts:65](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L65)_

---

## `Optional` properties

• **properties**? : _undefined | object_

_Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L62)_

---

## `Optional` required

• **required**? : _string[]_

_Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L57)_

---

## `Optional` title

• **title**? : _undefined | string_

_Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L40)_

---

## `Optional` type

• **type**? : _string | string[]_

_Defined in [types.ts:78](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L78)_

---

## `Optional` uniqueItems

• **uniqueItems**? : _undefined | false | true_

_Defined in [types.ts:54](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L54)_

<hr />

# Interface: LatestBlock

### Hierarchy

-   **LatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [types.ts:641](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L641)_

---

## number

• **number**: _number_

_Defined in [types.ts:640](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L640)_

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

-   **OrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [types.ts:567](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L567)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [types.ts:565](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L565)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:566](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L566)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:563](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L563)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:564](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L564)_

---

## timestampMs

• **timestampMs**: _number_

_Defined in [types.ts:562](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L562)_

<hr />

# Interface: OrderInfo

### Hierarchy

-   **OrderInfo**

### Properties

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L30)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:28](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L28)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:29](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L29)_

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

-   **RejectedOrderInfo**

### Properties

## kind

• **kind**: _[RejectedOrderKind](#enumeration-rejectedorderkind)_

_Defined in [types.ts:617](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L617)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:615](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L615)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:616](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L616)_

---

## status

• **status**: _[RejectedOrderStatus](#interface-rejectedorderstatus)_

_Defined in [types.ts:618](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L618)_

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

-   **RejectedOrderStatus**

### Properties

## code

• **code**: _string_

_Defined in [types.ts:635](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L635)_

---

## message

• **message**: _string_

_Defined in [types.ts:636](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L636)_

<hr />

# Interface: Stats

### Hierarchy

-   **Stats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [types.ts:678](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L678)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [types.ts:677](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L677)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [types.ts:669](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L669)_

---

## latestBlock

• **latestBlock**: _[LatestBlock](#interface-latestblock)_

_Defined in [types.ts:670](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L670)_

---

## maxExpirationTime

• **maxExpirationTime**: _BigNumber_

_Defined in [types.ts:675](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L675)_

---

## numOrders

• **numOrders**: _number_

_Defined in [types.ts:672](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L672)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [types.ts:673](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L673)_

---

## numPeers

• **numPeers**: _number_

_Defined in [types.ts:671](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L671)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [types.ts:674](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L674)_

---

## peerID

• **peerID**: _string_

_Defined in [types.ts:668](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L668)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [types.ts:665](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L665)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [types.ts:666](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L666)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [types.ts:667](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L667)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _Date_

_Defined in [types.ts:676](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L676)_

---

## version

• **version**: _string_

_Defined in [types.ts:664](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L664)_

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

-   **ValidationResults**

### Properties

## accepted

• **accepted**: _[AcceptedOrderInfo](#interface-acceptedorderinfo)[]_

_Defined in [types.ts:596](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L596)_

---

## rejected

• **rejected**: _[RejectedOrderInfo](#interface-rejectedorderinfo)[]_

_Defined in [types.ts:597](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L597)_

<hr />

# Interface: WethDepositEvent

### Hierarchy

-   **WethDepositEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:450](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L450)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:451](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L451)_

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

-   **WethWithdrawalEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:439](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L439)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:440](https://github.com/0xProject/0x-mesh/blob/06584d8f/packages/mesh-browser-lite/src/types.ts#L440)_

<hr />
