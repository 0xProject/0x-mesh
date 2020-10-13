# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

-   **Mesh**

### Constructors

## constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): _[Mesh](#class-mesh)_

_Defined in [mesh.ts:132](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L132)_

Instantiates a new Mesh instance.

**Parameters:**

| Name     | Type                        | Description                    |
| -------- | --------------------------- | ------------------------------ |
| `config` | [Config](#interface-config) | Configuration options for Mesh |

**Returns:** _[Mesh](#class-mesh)_

An instance of Mesh

### Properties

## `Optional` wrapper

• **wrapper**? : _MeshWrapper_

_Defined in [mesh.ts:129](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L129)_

### Methods

## addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[], `pinned`: boolean): _Promise‹[ValidationResults](#interface-validationresults)›_

_Defined in [mesh.ts:269](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L269)_

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

_Defined in [mesh.ts:207](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L207)_

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

_Defined in [mesh.ts:240](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L240)_

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

_Defined in [mesh.ts:190](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L190)_

Returns various stats about Mesh, including the total number of orders
and the number of peers Mesh is connected to.

**Returns:** _Promise‹[Stats](#interface-stats)›_

---

## onError

▸ **onError**(`handler`: function): _void_

_Defined in [mesh.ts:152](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L152)_

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

_Defined in [mesh.ts:165](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L165)_

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

_Defined in [mesh.ts:174](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/mesh.ts#L174)_

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** _Promise‹void›_

<hr />

# Enumeration: ContractEventKind

### Enumeration members

## ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

_Defined in [types.ts:505](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L505)_

---

## ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in [types.ts:507](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L507)_

---

## ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in [types.ts:506](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L506)_

---

## ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in [types.ts:501](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L501)_

---

## ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in [types.ts:500](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L500)_

---

## ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in [types.ts:503](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L503)_

---

## ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in [types.ts:504](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L504)_

---

## ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in [types.ts:502](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L502)_

---

## ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in [types.ts:509](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L509)_

---

## ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in [types.ts:510](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L510)_

---

## ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in [types.ts:508](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L508)_

---

## WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in [types.ts:511](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L511)_

---

## WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in [types.ts:512](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L512)_

<hr />

# Enumeration: OrderEventEndState

### Enumeration members

## Added

• **Added**: = "ADDED"

_Defined in [types.ts:575](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L575)_

---

## Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in [types.ts:578](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L578)_

---

## Expired

• **Expired**: = "EXPIRED"

_Defined in [types.ts:579](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L579)_

---

## FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

_Defined in [types.ts:582](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L582)_

---

## Filled

• **Filled**: = "FILLED"

_Defined in [types.ts:576](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L576)_

---

## FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

_Defined in [types.ts:577](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L577)_

---

## Invalid

• **Invalid**: = "INVALID"

_Defined in [types.ts:574](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L574)_

---

## StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

_Defined in [types.ts:583](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L583)_

---

## Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in [types.ts:580](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L580)_

---

## Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in [types.ts:581](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L581)_

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.

### Enumeration members

## MeshError

• **MeshError**: = "MESH_ERROR"

_Defined in [types.ts:713](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L713)_

---

## MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

_Defined in [types.ts:714](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L714)_

---

## ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

_Defined in [types.ts:712](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L712)_

<hr />

# Enumeration: Verbosity

### Enumeration members

## Debug

• **Debug**: = 5

_Defined in [types.ts:238](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L238)_

---

## Error

• **Error**: = 2

_Defined in [types.ts:235](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L235)_

---

## Fatal

• **Fatal**: = 1

_Defined in [types.ts:234](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L234)_

---

## Info

• **Info**: = 4

_Defined in [types.ts:237](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L237)_

---

## Panic

• **Panic**: = 0

_Defined in [types.ts:233](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L233)_

---

## Trace

• **Trace**: = 6

_Defined in [types.ts:239](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L239)_

---

## Warn

• **Warn**: = 3

_Defined in [types.ts:236](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L236)_

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

-   **AcceptedOrderInfo**

### Properties

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:693](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L693)_

---

## isNew

• **isNew**: _boolean_

_Defined in [types.ts:694](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L694)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:691](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L691)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:692](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L692)_

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

-   **Config**

### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : _undefined | number_

_Defined in [types.ts:144](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L144)_

---

## `Optional` bootstrapList

• **bootstrapList**? : _string[]_

_Defined in [types.ts:137](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L137)_

---

## `Optional` customContractAddresses

• **customContractAddresses**? : _[ContractAddresses](#interface-contractaddresses)_

_Defined in [types.ts:188](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L188)_

---

## `Optional` customOrderFilter

• **customOrderFilter**? : _[JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:213](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L213)_

---

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : _undefined | false | true_

_Defined in [types.ts:161](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L161)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [types.ts:129](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L129)_

---

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : _undefined | number_

_Defined in [types.ts:153](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L153)_

---

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : _undefined | number_

_Defined in [types.ts:166](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L166)_

---

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : _undefined | number_

_Defined in [types.ts:172](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L172)_

---

## `Optional` ethereumRPCURL

• **ethereumRPCURL**? : _undefined | string_

_Defined in [types.ts:126](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L126)_

---

## `Optional` maxBytesPerSecond

• **maxBytesPerSecond**? : _undefined | number_

_Defined in [types.ts:219](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L219)_

---

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : _undefined | number_

_Defined in [types.ts:193](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L193)_

---

## `Optional` useBootstrapList

• **useBootstrapList**? : _undefined | false | true_

_Defined in [types.ts:132](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L132)_

---

## `Optional` verbosity

• **verbosity**? : _[Verbosity](#enumeration-verbosity)_

_Defined in [types.ts:123](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L123)_

---

## `Optional` web3Provider

• **web3Provider**? : _SupportedProvider_

_Defined in [types.ts:216](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L216)_

<hr />

# Interface: ContractAddresses

### Hierarchy

-   **ContractAddresses**

### Properties

## devUtils

• **devUtils**: _string_

_Defined in [types.ts:224](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L224)_

---

## erc1155Proxy

• **erc1155Proxy**: _string_

_Defined in [types.ts:227](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L227)_

---

## erc20Proxy

• **erc20Proxy**: _string_

_Defined in [types.ts:225](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L225)_

---

## erc721Proxy

• **erc721Proxy**: _string_

_Defined in [types.ts:226](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L226)_

---

## exchange

• **exchange**: _string_

_Defined in [types.ts:223](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L223)_

---

## `Optional` weth9

• **weth9**? : _undefined | string_

_Defined in [types.ts:228](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L228)_

---

## `Optional` zrxToken

• **zrxToken**? : _undefined | string_

_Defined in [types.ts:229](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L229)_

<hr />

# Interface: ContractEvent

### Hierarchy

-   **ContractEvent**

### Properties

## address

• **address**: _string_

_Defined in [types.ts:553](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L553)_

---

## blockHash

• **blockHash**: _string_

_Defined in [types.ts:548](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L548)_

---

## isRemoved

• **isRemoved**: _boolean_

_Defined in [types.ts:552](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L552)_

---

## kind

• **kind**: _[ContractEventKind](#enumeration-contracteventkind)_

_Defined in [types.ts:554](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L554)_

---

## logIndex

• **logIndex**: _number_

_Defined in [types.ts:551](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L551)_

---

## parameters

• **parameters**: _ContractEventParameters_

_Defined in [types.ts:555](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L555)_

---

## txHash

• **txHash**: _string_

_Defined in [types.ts:549](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L549)_

---

## txIndex

• **txIndex**: _number_

_Defined in [types.ts:550](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L550)_

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

-   **ERC1155ApprovalForAllEvent**

### Properties

## approved

• **approved**: _boolean_

_Defined in [types.ts:417](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L417)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:416](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L416)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:415](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L415)_

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

-   **ERC1155TransferBatchEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:399](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L399)_

---

## ids

• **ids**: _BigNumber[]_

_Defined in [types.ts:401](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L401)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:398](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L398)_

---

## to

• **to**: _string_

_Defined in [types.ts:400](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L400)_

---

## values

• **values**: _BigNumber[]_

_Defined in [types.ts:402](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L402)_

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

-   **ERC1155TransferSingleEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:382](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L382)_

---

## id

• **id**: _BigNumber_

_Defined in [types.ts:384](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L384)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:381](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L381)_

---

## to

• **to**: _string_

_Defined in [types.ts:383](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L383)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:385](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L385)_

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

-   **ERC20ApprovalEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:336](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L336)_

---

## spender

• **spender**: _string_

_Defined in [types.ts:337](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L337)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:338](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L338)_

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

-   **ERC20TransferEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:323](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L323)_

---

## to

• **to**: _string_

_Defined in [types.ts:324](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L324)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:325](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L325)_

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

-   **ERC721ApprovalEvent**

### Properties

## approved

• **approved**: _string_

_Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L363)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:362](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L362)_

---

## tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L364)_

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

-   **ERC721ApprovalForAllEvent**

### Properties

## approved

• **approved**: _boolean_

_Defined in [types.ts:377](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L377)_

---

## operator

• **operator**: _string_

_Defined in [types.ts:376](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L376)_

---

## owner

• **owner**: _string_

_Defined in [types.ts:375](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L375)_

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

-   **ERC721TransferEvent**

### Properties

## from

• **from**: _string_

_Defined in [types.ts:349](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L349)_

---

## to

• **to**: _string_

_Defined in [types.ts:350](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L350)_

---

## tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:351](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L351)_

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

-   **ExchangeCancelEvent**

### Properties

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:458](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L458)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:456](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L456)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:460](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L460)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:459](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L459)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:457](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L457)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:461](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L461)_

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

-   **ExchangeCancelUpToEvent**

### Properties

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:465](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L465)_

---

## orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in [types.ts:467](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L467)_

---

## orderSenderAddress

• **orderSenderAddress**: _string_

_Defined in [types.ts:466](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L466)_

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

-   **ExchangeFillEvent**

### Properties

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L424)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:421](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L421)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:431](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L431)_

---

## makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:425](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L425)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Defined in [types.ts:433](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L433)_

---

## makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in [types.ts:427](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L427)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:430](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L430)_

---

## protocolFeePaid

• **protocolFeePaid**: _BigNumber_

_Defined in [types.ts:429](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L429)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:423](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L423)_

---

## takerAddress

• **takerAddress**: _string_

_Defined in [types.ts:422](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L422)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:432](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L432)_

---

## takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:426](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L426)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Defined in [types.ts:434](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L434)_

---

## takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in [types.ts:428](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L428)_

<hr />

# Interface: GetOrdersResponse

### Hierarchy

-   **GetOrdersResponse**

### Properties

## ordersInfos

• **ordersInfos**: _[OrderInfo](#interface-orderinfo)[]_

_Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L45)_

---

## timestamp

• **timestamp**: _number_

_Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L44)_

<hr />

# Interface: JsonSchema

An interface for JSON schema types, which are used for custom order filters.

### Hierarchy

-   **JsonSchema**

### Properties

## `Optional` \$ref

• **\$ref**? : _undefined | string_

_Defined in [types.ts:67](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L67)_

---

## `Optional` \$schema

• **\$schema**? : _undefined | string_

_Defined in [types.ts:66](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L66)_

---

## `Optional` additionalItems

• **additionalItems**? : _boolean | [JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:78](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L78)_

---

## `Optional` additionalProperties

• **additionalProperties**? : _boolean | [JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:86](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L86)_

---

## `Optional` allOf

• **allOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:108](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L108)_

---

## `Optional` anyOf

• **anyOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:109](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L109)_

---

## `Optional` const

• **const**? : _any_

_Defined in [types.ts:105](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L105)_

---

## `Optional` definitions

• **definitions**? : _undefined | object_

_Defined in [types.ts:87](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L87)_

---

## `Optional` dependencies

• **dependencies**? : _undefined | object_

_Defined in [types.ts:96](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L96)_

---

## `Optional` description

• **description**? : _undefined | string_

_Defined in [types.ts:69](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L69)_

---

## `Optional` enum

• **enum**? : _any[]_

_Defined in [types.ts:99](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L99)_

---

## `Optional` exclusiveMaximum

• **exclusiveMaximum**? : _undefined | false | true_

_Defined in [types.ts:72](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L72)_

---

## `Optional` exclusiveMinimum

• **exclusiveMinimum**? : _undefined | false | true_

_Defined in [types.ts:74](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L74)_

---

## `Optional` format

• **format**? : _undefined | string_

_Defined in [types.ts:107](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L107)_

---

## `Optional` id

• **id**? : _undefined | string_

_Defined in [types.ts:65](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L65)_

---

## `Optional` items

• **items**? : _[JsonSchema](#interface-jsonschema) | [JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L79)_

---

## `Optional` maxItems

• **maxItems**? : _undefined | number_

_Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L80)_

---

## `Optional` maxLength

• **maxLength**? : _undefined | number_

_Defined in [types.ts:75](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L75)_

---

## `Optional` maxProperties

• **maxProperties**? : _undefined | number_

_Defined in [types.ts:83](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L83)_

---

## `Optional` maximum

• **maximum**? : _undefined | number_

_Defined in [types.ts:71](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L71)_

---

## `Optional` minItems

• **minItems**? : _undefined | number_

_Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L81)_

---

## `Optional` minLength

• **minLength**? : _undefined | number_

_Defined in [types.ts:76](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L76)_

---

## `Optional` minProperties

• **minProperties**? : _undefined | number_

_Defined in [types.ts:84](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L84)_

---

## `Optional` minimum

• **minimum**? : _undefined | number_

_Defined in [types.ts:73](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L73)_

---

## `Optional` multipleOf

• **multipleOf**? : _undefined | number_

_Defined in [types.ts:70](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L70)_

---

## `Optional` not

• **not**? : _[JsonSchema](#interface-jsonschema)_

_Defined in [types.ts:111](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L111)_

---

## `Optional` oneOf

• **oneOf**? : _[JsonSchema](#interface-jsonschema)[]_

_Defined in [types.ts:110](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L110)_

---

## `Optional` pattern

• **pattern**? : _string | RegExp_

_Defined in [types.ts:77](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L77)_

---

## `Optional` patternProperties

• **patternProperties**? : _undefined | object_

_Defined in [types.ts:93](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L93)_

---

## `Optional` properties

• **properties**? : _undefined | object_

_Defined in [types.ts:90](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L90)_

---

## `Optional` required

• **required**? : _string[]_

_Defined in [types.ts:85](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L85)_

---

## `Optional` title

• **title**? : _undefined | string_

_Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L68)_

---

## `Optional` type

• **type**? : _string | string[]_

_Defined in [types.ts:106](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L106)_

---

## `Optional` uniqueItems

• **uniqueItems**? : _undefined | false | true_

_Defined in [types.ts:82](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L82)_

<hr />

# Interface: LatestBlock

### Hierarchy

-   **LatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [types.ts:733](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L733)_

---

## number

• **number**: _BigNumber_

_Defined in [types.ts:732](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L732)_

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

-   **OrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [types.ts:606](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L606)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [types.ts:604](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L604)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:605](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L605)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:602](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L602)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:603](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L603)_

---

## timestampMs

• **timestampMs**: _number_

_Defined in [types.ts:601](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L601)_

<hr />

# Interface: OrderInfo

### Hierarchy

-   **OrderInfo**

### Properties

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:58](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L58)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L56)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L57)_

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

-   **RejectedOrderInfo**

### Properties

## kind

• **kind**: _[RejectedOrderKind](#enumeration-rejectedorderkind)_

_Defined in [types.ts:704](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L704)_

---

## orderHash

• **orderHash**: _string_

_Defined in [types.ts:702](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L702)_

---

## signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:703](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L703)_

---

## status

• **status**: _[RejectedOrderStatus](#interface-rejectedorderstatus)_

_Defined in [types.ts:705](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L705)_

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

-   **RejectedOrderStatus**

### Properties

## code

• **code**: _string_

_Defined in [types.ts:721](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L721)_

---

## message

• **message**: _string_

_Defined in [types.ts:722](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L722)_

<hr />

# Interface: Stats

### Hierarchy

-   **Stats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [types.ts:770](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L770)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [types.ts:769](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L769)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [types.ts:761](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L761)_

---

## `Optional` latestBlock

• **latestBlock**? : _[LatestBlock](#interface-latestblock)_

_Defined in [types.ts:762](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L762)_

---

## maxExpirationTime

• **maxExpirationTime**: _BigNumber_

_Defined in [types.ts:767](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L767)_

---

## numOrders

• **numOrders**: _number_

_Defined in [types.ts:764](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L764)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [types.ts:765](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L765)_

---

## numPeers

• **numPeers**: _number_

_Defined in [types.ts:763](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L763)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [types.ts:766](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L766)_

---

## peerID

• **peerID**: _string_

_Defined in [types.ts:760](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L760)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [types.ts:757](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L757)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [types.ts:758](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L758)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [types.ts:759](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L759)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _Date_

_Defined in [types.ts:768](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L768)_

---

## version

• **version**: _string_

_Defined in [types.ts:756](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L756)_

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

-   **ValidationResults**

### Properties

## accepted

• **accepted**: _[AcceptedOrderInfo](#interface-acceptedorderinfo)[]_

_Defined in [types.ts:683](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L683)_

---

## rejected

• **rejected**: _[RejectedOrderInfo](#interface-rejectedorderinfo)[]_

_Defined in [types.ts:684](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L684)_

<hr />

# Interface: WethDepositEvent

### Hierarchy

-   **WethDepositEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:489](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L489)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:490](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L490)_

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

-   **WethWithdrawalEvent**

### Properties

## owner

• **owner**: _string_

_Defined in [types.ts:478](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L478)_

---

## value

• **value**: _BigNumber_

_Defined in [types.ts:479](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-browser-lite/src/types.ts#L479)_

<hr />
