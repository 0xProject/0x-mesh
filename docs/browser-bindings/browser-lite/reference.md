# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [mesh.ts:141](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L141)*

Instantiates a new Mesh instance.

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`config` | [Config](#interface-config) | Configuration options for Mesh |

**Returns:** *[Mesh](#class-mesh)*

An instance of Mesh

### Methods

##  addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[], `pinned`: boolean): *Promise‹[ValidationResults](#interface-validationresults)›*

*Defined in [mesh.ts:291](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L291)*

Validates and adds the given orders to Mesh. If an order is successfully
added, Mesh will share it with any peers in the network and start
watching it for changes (e.g. filled, canceled, expired). The returned
promise will only be rejected if there was an error validating or adding
the order; it will not be rejected for any invalid orders (check
results.rejected instead).

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`orders` | SignedOrder[] | - | An array of orders to add. |
`pinned` | boolean | true | Whether or not the orders should be pinned. Pinned orders will not be affected by any DDoS prevention or incentive mechanisms and will always stay in storage until they are no longer fillable. |

**Returns:** *Promise‹[ValidationResults](#interface-validationresults)›*

Validation results for the given orders, indicating which orders
were accepted and which were rejected.

___

##  getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

*Defined in [mesh.ts:220](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L220)*

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`perPage` | number | 200 | number of signedOrders to fetch per paginated request |

**Returns:** *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts

___

##  getOrdersForPageAsync

▸ **getOrdersForPageAsync**(`page`: number, `perPage`: number, `snapshotID?`: undefined | string): *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

*Defined in [mesh.ts:262](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L262)*

Get page of 0x signed orders stored on the Mesh node at the specified snapshot

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`page` | number | Page index at which to retrieve orders |
`perPage` | number | Number of signedOrders to fetch per paginated request |
`snapshotID?` | undefined &#124; string | The DB snapshot at which to fetch orders. If omitted, a new snapshot is created |

**Returns:** *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts

___

##  getStatsAsync

▸ **getStatsAsync**(): *Promise‹[Stats](#interface-stats)›*

*Defined in [mesh.ts:203](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L203)*

Returns various stats about Mesh, including the total number of orders
and the number of peers Mesh is connected to.

**Returns:** *Promise‹[Stats](#interface-stats)›*

___

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [mesh.ts:161](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L161)*

Registers a handler which will be called in the event of a critical
error. Note that the handler will not be called for non-critical errors.
In order to ensure no errors are missed, this should be called before
startAsync.

**Parameters:**

▪ **handler**: *function*

The handler to be called.

▸ (`err`: Error): *void*

**Parameters:**

Name | Type |
------ | ------ |
`err` | Error |

**Returns:** *void*

___

##  onOrderEvents

▸ **onOrderEvents**(`handler`: function): *void*

*Defined in [mesh.ts:176](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L176)*

Registers a handler which will be called for any incoming order events.
Order events are fired whenver an order is added, canceled, expired, or
filled. In order to ensure no events are missed, this should be called
before startAsync.

**Parameters:**

▪ **handler**: *function*

The handler to be called.

▸ (`events`: [OrderEvent](#interface-orderevent)[]): *void*

**Parameters:**

Name | Type |
------ | ------ |
`events` | [OrderEvent](#interface-orderevent)[] |

**Returns:** *void*

___

##  startAsync

▸ **startAsync**(): *Promise‹void›*

*Defined in [mesh.ts:187](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/mesh.ts#L187)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*


<hr />

# Enumeration: ContractEventKind


### Enumeration members

##  ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

*Defined in [types.ts:444](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L444)*

___

##  ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

*Defined in [types.ts:446](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L446)*

___

##  ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

*Defined in [types.ts:445](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L445)*

___

##  ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

*Defined in [types.ts:440](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L440)*

___

##  ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

*Defined in [types.ts:439](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L439)*

___

##  ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

*Defined in [types.ts:442](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L442)*

___

##  ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

*Defined in [types.ts:443](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L443)*

___

##  ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

*Defined in [types.ts:441](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L441)*

___

##  ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

*Defined in [types.ts:448](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L448)*

___

##  ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

*Defined in [types.ts:449](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L449)*

___

##  ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

*Defined in [types.ts:447](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L447)*

___

##  WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

*Defined in [types.ts:450](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L450)*

___

##  WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

*Defined in [types.ts:451](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L451)*


<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [types.ts:509](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L509)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [types.ts:512](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L512)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [types.ts:513](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L513)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [types.ts:516](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L516)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [types.ts:510](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L510)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [types.ts:511](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L511)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [types.ts:508](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L508)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [types.ts:517](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L517)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [types.ts:514](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L514)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [types.ts:515](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L515)*


<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [types.ts:600](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L600)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [types.ts:598](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L598)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [types.ts:599](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L599)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [types.ts:597](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L597)*


<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [types.ts:209](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L209)*

___

##  Error

• **Error**: = 2

*Defined in [types.ts:206](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L206)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [types.ts:205](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L205)*

___

##  Info

• **Info**: = 4

*Defined in [types.ts:208](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L208)*

___

##  Panic

• **Panic**: = 0

*Defined in [types.ts:204](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L204)*

___

##  Trace

• **Trace**: = 6

*Defined in [types.ts:210](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L210)*

___

##  Warn

• **Warn**: = 3

*Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L207)*


<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:578](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L578)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [types.ts:579](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L579)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:576](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L576)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:577](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L577)*


<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [types.ts:116](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L116)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [types.ts:109](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L109)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#interface-contractaddresses)*

*Defined in [types.ts:160](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L160)*

___

## `Optional` customOrderFilter

• **customOrderFilter**? : *[JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L185)*

___

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : *undefined | false | true*

*Defined in [types.ts:133](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L133)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:101](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L101)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [types.ts:125](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L125)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [types.ts:138](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L138)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [types.ts:144](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L144)*

___

## `Optional` ethereumRPCURL

• **ethereumRPCURL**? : *undefined | string*

*Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L98)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [types.ts:165](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L165)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L104)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [types.ts:95](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L95)*

___

## `Optional` web3Provider

• **web3Provider**? : *SupportedProvider*

*Defined in [types.ts:188](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L188)*


<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L197)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [types.ts:198](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L198)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [types.ts:193](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L193)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [types.ts:196](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L196)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [types.ts:194](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L194)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [types.ts:195](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L195)*

___

##  exchange

• **exchange**: *string*

*Defined in [types.ts:192](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L192)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [types.ts:199](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L199)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [types.ts:200](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L200)*


<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [types.ts:490](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L490)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [types.ts:485](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L485)*

___

##  isRemoved

• **isRemoved**: *boolean*

*Defined in [types.ts:489](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L489)*

___

##  kind

• **kind**: *[ContractEventKind](#enumeration-contracteventkind)*

*Defined in [types.ts:491](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L491)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [types.ts:488](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L488)*

___

##  parameters

• **parameters**: *[ContractEventParameters](#contracteventparameters)*

*Defined in [types.ts:492](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L492)*

___

##  txHash

• **txHash**: *string*

*Defined in [types.ts:486](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L486)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [types.ts:487](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L487)*


<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [types.ts:360](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L360)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:359](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L359)*

___

##  owner

• **owner**: *string*

*Defined in [types.ts:358](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L358)*


<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:343](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L343)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [types.ts:345](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L345)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:342](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L342)*

___

##  to

• **to**: *string*

*Defined in [types.ts:344](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L344)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L346)*


<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:327](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L327)*

___

##  id

• **id**: *BigNumber*

*Defined in [types.ts:329](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L329)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:326](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L326)*

___

##  to

• **to**: *string*

*Defined in [types.ts:328](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L328)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:330](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L330)*


<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:284](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L284)*

___

##  spender

• **spender**: *string*

*Defined in [types.ts:285](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L285)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L286)*


<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:272](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L272)*

___

##  to

• **to**: *string*

*Defined in [types.ts:273](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L273)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:274](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L274)*


<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [types.ts:309](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L309)*

___

##  owner

• **owner**: *string*

*Defined in [types.ts:308](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L308)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [types.ts:310](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L310)*


<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [types.ts:322](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L322)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:321](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L321)*

___

##  owner

• **owner**: *string*

*Defined in [types.ts:320](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L320)*


<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:296](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L296)*

___

##  to

• **to**: *string*

*Defined in [types.ts:297](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L297)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [types.ts:298](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L298)*


<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:400](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L400)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:398](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L398)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:402](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L402)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:401](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L401)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:399](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L399)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:403](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L403)*


<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:407](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L407)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [types.ts:409](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L409)*

___

##  orderSenderAddress

• **orderSenderAddress**: *string*

*Defined in [types.ts:408](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L408)*


<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:367](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L367)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L364)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:374](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L374)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [types.ts:368](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L368)*

___

##  makerFeeAssetData

• **makerFeeAssetData**: *string*

*Defined in [types.ts:376](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L376)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [types.ts:370](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L370)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:373](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L373)*

___

##  protocolFeePaid

• **protocolFeePaid**: *BigNumber*

*Defined in [types.ts:372](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L372)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:366](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L366)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:365](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L365)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:375](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L375)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [types.ts:369](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L369)*

___

##  takerFeeAssetData

• **takerFeeAssetData**: *string*

*Defined in [types.ts:377](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L377)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [types.ts:371](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L371)*


<hr />

# Interface: GetOrdersResponse

### Hierarchy

* **GetOrdersResponse**


### Properties

##  ordersInfos

• **ordersInfos**: *[OrderInfo](#interface-orderinfo)[]*

*Defined in [types.ts:18](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L18)*

___

##  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L16)*

___

##  snapshotTimestamp

• **snapshotTimestamp**: *number*

*Defined in [types.ts:17](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L17)*


<hr />

# Interface: JsonSchema

An interface for JSON schema types, which are used for custom order filters.

### Hierarchy

* **JsonSchema**


### Properties

## `Optional` $ref

• **$ref**? : *undefined | string*

*Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L39)*

___

## `Optional` $schema

• **$schema**? : *undefined | string*

*Defined in [types.ts:38](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L38)*

___

## `Optional` additionalItems

• **additionalItems**? : *boolean | [JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L50)*

___

## `Optional` additionalProperties

• **additionalProperties**? : *boolean | [JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:58](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L58)*

___

## `Optional` allOf

• **allOf**? : *[JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L80)*

___

## `Optional` anyOf

• **anyOf**? : *[JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L81)*

___

## `Optional` const

• **const**? : *any*

*Defined in [types.ts:77](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L77)*

___

## `Optional` definitions

• **definitions**? : *undefined | object*

*Defined in [types.ts:59](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L59)*

___

## `Optional` dependencies

• **dependencies**? : *undefined | object*

*Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L68)*

___

## `Optional` description

• **description**? : *undefined | string*

*Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L41)*

___

## `Optional` enum

• **enum**? : *any[]*

*Defined in [types.ts:71](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L71)*

___

## `Optional` exclusiveMaximum

• **exclusiveMaximum**? : *undefined | false | true*

*Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L44)*

___

## `Optional` exclusiveMinimum

• **exclusiveMinimum**? : *undefined | false | true*

*Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L46)*

___

## `Optional` format

• **format**? : *undefined | string*

*Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L79)*

___

## `Optional` id

• **id**? : *undefined | string*

*Defined in [types.ts:37](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L37)*

___

## `Optional` items

• **items**? : *[JsonSchema](#interface-jsonschema) | [JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:51](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L51)*

___

## `Optional` maxItems

• **maxItems**? : *undefined | number*

*Defined in [types.ts:52](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L52)*

___

## `Optional` maxLength

• **maxLength**? : *undefined | number*

*Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L47)*

___

## `Optional` maxProperties

• **maxProperties**? : *undefined | number*

*Defined in [types.ts:55](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L55)*

___

## `Optional` maximum

• **maximum**? : *undefined | number*

*Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L43)*

___

## `Optional` minItems

• **minItems**? : *undefined | number*

*Defined in [types.ts:53](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L53)*

___

## `Optional` minLength

• **minLength**? : *undefined | number*

*Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L48)*

___

## `Optional` minProperties

• **minProperties**? : *undefined | number*

*Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L56)*

___

## `Optional` minimum

• **minimum**? : *undefined | number*

*Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L45)*

___

## `Optional` multipleOf

• **multipleOf**? : *undefined | number*

*Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L42)*

___

## `Optional` not

• **not**? : *[JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:83](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L83)*

___

## `Optional` oneOf

• **oneOf**? : *[JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:82](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L82)*

___

## `Optional` pattern

• **pattern**? : *string | RegExp*

*Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L49)*

___

## `Optional` patternProperties

• **patternProperties**? : *undefined | object*

*Defined in [types.ts:65](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L65)*

___

## `Optional` properties

• **properties**? : *undefined | object*

*Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L62)*

___

## `Optional` required

• **required**? : *string[]*

*Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L57)*

___

## `Optional` title

• **title**? : *undefined | string*

*Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L40)*

___

## `Optional` type

• **type**? : *string | string[]*

*Defined in [types.ts:78](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L78)*

___

## `Optional` uniqueItems

• **uniqueItems**? : *undefined | false | true*

*Defined in [types.ts:54](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L54)*


<hr />

# Interface: LatestBlock

### Hierarchy

* **LatestBlock**


### Properties

##  hash

• **hash**: *string*

*Defined in [types.ts:613](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L613)*

___

##  number

• **number**: *number*

*Defined in [types.ts:612](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L612)*


<hr />

# Interface: MeshWrapper

### Hierarchy

* **MeshWrapper**


### Methods

##  addOrdersAsync

▸ **addOrdersAsync**(`orders`: [WrapperSignedOrder](#interface-wrappersignedorder)[], `pinned`: boolean): *Promise‹[WrapperValidationResults](#interface-wrappervalidationresults)›*

*Defined in [types.ts:227](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L227)*

**Parameters:**

Name | Type |
------ | ------ |
`orders` | [WrapperSignedOrder](#interface-wrappersignedorder)[] |
`pinned` | boolean |

**Returns:** *Promise‹[WrapperValidationResults](#interface-wrappervalidationresults)›*

___

##  getOrdersForPageAsync

▸ **getOrdersForPageAsync**(`page`: number, `perPage`: number, `snapshotID?`: undefined | string): *Promise‹[WrapperGetOrdersResponse](#interface-wrappergetordersresponse)›*

*Defined in [types.ts:226](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L226)*

**Parameters:**

Name | Type |
------ | ------ |
`page` | number |
`perPage` | number |
`snapshotID?` | undefined &#124; string |

**Returns:** *Promise‹[WrapperGetOrdersResponse](#interface-wrappergetordersresponse)›*

___

##  getStatsAsync

▸ **getStatsAsync**(): *Promise‹[WrapperStats](#interface-wrapperstats)›*

*Defined in [types.ts:225](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L225)*

**Returns:** *Promise‹[WrapperStats](#interface-wrapperstats)›*

___

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [types.ts:223](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L223)*

**Parameters:**

▪ **handler**: *function*

▸ (`err`: Error): *void*

**Parameters:**

Name | Type |
------ | ------ |
`err` | Error |

**Returns:** *void*

___

##  onOrderEvents

▸ **onOrderEvents**(`handler`: function): *void*

*Defined in [types.ts:224](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L224)*

**Parameters:**

▪ **handler**: *function*

▸ (`events`: [WrapperOrderEvent](#interface-wrapperorderevent)[]): *void*

**Parameters:**

Name | Type |
------ | ------ |
`events` | [WrapperOrderEvent](#interface-wrapperorderevent)[] |

**Returns:** *void*

___

##  startAsync

▸ **startAsync**(): *Promise‹void›*

*Defined in [types.ts:222](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L222)*

**Returns:** *Promise‹void›*


<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#interface-contractevent)[]*

*Defined in [types.ts:539](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L539)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:537](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L537)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:538](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L538)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:535](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L535)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:536](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L536)*

___

##  timestampMs

• **timestampMs**: *number*

*Defined in [types.ts:534](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L534)*


<hr />

# Interface: OrderInfo

### Hierarchy

* **OrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L30)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:28](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L28)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:29](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L29)*


<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [types.ts:589](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L589)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:587](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L587)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:588](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L588)*

___

##  status

• **status**: *[RejectedOrderStatus](#interface-rejectedorderstatus)*

*Defined in [types.ts:590](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L590)*


<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [types.ts:607](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L607)*

___

##  message

• **message**: *string*

*Defined in [types.ts:608](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L608)*


<hr />

# Interface: Stats

### Hierarchy

* **Stats**


### Properties

##  ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: *number*

*Defined in [types.ts:649](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L649)*

___

##  ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: *number*

*Defined in [types.ts:648](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L648)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:640](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L640)*

___

##  latestBlock

• **latestBlock**: *[LatestBlock](#interface-latestblock)*

*Defined in [types.ts:641](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L641)*

___

##  maxExpirationTime

• **maxExpirationTime**: *BigNumber*

*Defined in [types.ts:646](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L646)*

___

##  numOrders

• **numOrders**: *number*

*Defined in [types.ts:643](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L643)*

___

##  numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: *number*

*Defined in [types.ts:644](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L644)*

___

##  numPeers

• **numPeers**: *number*

*Defined in [types.ts:642](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L642)*

___

##  numPinnedOrders

• **numPinnedOrders**: *number*

*Defined in [types.ts:645](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L645)*

___

##  peerID

• **peerID**: *string*

*Defined in [types.ts:639](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L639)*

___

##  pubSubTopic

• **pubSubTopic**: *string*

*Defined in [types.ts:636](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L636)*

___

##  rendezvous

• **rendezvous**: *string*

*Defined in [types.ts:637](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L637)*

___

##  secondaryRendezvous

• **secondaryRendezvous**: *string[]*

*Defined in [types.ts:638](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L638)*

___

##  startOfCurrentUTCDay

• **startOfCurrentUTCDay**: *Date*

*Defined in [types.ts:647](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L647)*

___

##  version

• **version**: *string*

*Defined in [types.ts:635](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L635)*


<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#interface-acceptedorderinfo)[]*

*Defined in [types.ts:568](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L568)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#interface-rejectedorderinfo)[]*

*Defined in [types.ts:569](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L569)*


<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:429](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L429)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:430](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L430)*


<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:419](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L419)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:420](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L420)*


<hr />

# Interface: WrapperAcceptedOrderInfo

### Hierarchy

* **WrapperAcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:552](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L552)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [types.ts:553](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L553)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:550](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L550)*

___

##  signedOrder

• **signedOrder**: *[WrapperSignedOrder](#interface-wrappersignedorder)*

*Defined in [types.ts:551](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L551)*


<hr />

# Interface: WrapperConfig

### Hierarchy

* **WrapperConfig**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [types.ts:237](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L237)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *undefined | string*

*Defined in [types.ts:236](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L236)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *undefined | string*

*Defined in [types.ts:242](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L242)*

___

## `Optional` customOrderFilter

• **customOrderFilter**? : *undefined | string*

*Defined in [types.ts:244](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L244)*

___

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : *undefined | false | true*

*Defined in [types.ts:241](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L241)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:234](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L234)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [types.ts:238](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L238)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [types.ts:239](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L239)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [types.ts:240](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L240)*

___

## `Optional` ethereumRPCURL

• **ethereumRPCURL**? : *undefined | string*

*Defined in [types.ts:233](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L233)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [types.ts:243](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L243)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [types.ts:235](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L235)*

___

## `Optional` verbosity

• **verbosity**? : *undefined | number*

*Defined in [types.ts:232](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L232)*

___

## `Optional` web3Provider

• **web3Provider**? : *ZeroExProvider*

*Defined in [types.ts:245](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L245)*


<hr />

# Interface: WrapperContractEvent

### Hierarchy

* **WrapperContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [types.ts:502](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L502)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [types.ts:497](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L497)*

___

##  isRemoved

• **isRemoved**: *boolean*

*Defined in [types.ts:501](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L501)*

___

##  kind

• **kind**: *string*

*Defined in [types.ts:503](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L503)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [types.ts:500](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L500)*

___

##  parameters

• **parameters**: *[WrapperContractEventParameters](#wrappercontracteventparameters)*

*Defined in [types.ts:504](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L504)*

___

##  txHash

• **txHash**: *string*

*Defined in [types.ts:498](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L498)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [types.ts:499](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L499)*


<hr />

# Interface: WrapperERC1155TransferBatchEvent

### Hierarchy

* **WrapperERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:351](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L351)*

___

##  ids

• **ids**: *string[]*

*Defined in [types.ts:353](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L353)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:350](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L350)*

___

##  to

• **to**: *string*

*Defined in [types.ts:352](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L352)*

___

##  values

• **values**: *string[]*

*Defined in [types.ts:354](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L354)*


<hr />

# Interface: WrapperERC1155TransferSingleEvent

### Hierarchy

* **WrapperERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:335](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L335)*

___

##  id

• **id**: *string*

*Defined in [types.ts:337](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L337)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:334](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L334)*

___

##  to

• **to**: *string*

*Defined in [types.ts:336](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L336)*

___

##  value

• **value**: *string*

*Defined in [types.ts:338](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L338)*


<hr />

# Interface: WrapperERC20ApprovalEvent

### Hierarchy

* **WrapperERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:290](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L290)*

___

##  spender

• **spender**: *string*

*Defined in [types.ts:291](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L291)*

___

##  value

• **value**: *string*

*Defined in [types.ts:292](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L292)*


<hr />

# Interface: WrapperERC20TransferEvent

### Hierarchy

* **WrapperERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:278](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L278)*

___

##  to

• **to**: *string*

*Defined in [types.ts:279](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L279)*

___

##  value

• **value**: *string*

*Defined in [types.ts:280](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L280)*


<hr />

# Interface: WrapperERC721ApprovalEvent

### Hierarchy

* **WrapperERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [types.ts:315](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L315)*

___

##  owner

• **owner**: *string*

*Defined in [types.ts:314](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L314)*

___

##  tokenId

• **tokenId**: *string*

*Defined in [types.ts:316](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L316)*


<hr />

# Interface: WrapperERC721TransferEvent

### Hierarchy

* **WrapperERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:302](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L302)*

___

##  to

• **to**: *string*

*Defined in [types.ts:303](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L303)*

___

##  tokenId

• **tokenId**: *string*

*Defined in [types.ts:304](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L304)*


<hr />

# Interface: WrapperExchangeCancelUpToEvent

### Hierarchy

* **WrapperExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:413](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L413)*

___

##  orderEpoch

• **orderEpoch**: *string*

*Defined in [types.ts:415](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L415)*

___

##  orderSenderAddress

• **orderSenderAddress**: *string*

*Defined in [types.ts:414](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L414)*


<hr />

# Interface: WrapperExchangeFillEvent

### Hierarchy

* **WrapperExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:384](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L384)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:381](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L381)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:391](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L391)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *string*

*Defined in [types.ts:385](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L385)*

___

##  makerFeeAssetData

• **makerFeeAssetData**: *string*

*Defined in [types.ts:393](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L393)*

___

##  makerFeePaid

• **makerFeePaid**: *string*

*Defined in [types.ts:387](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L387)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:390](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L390)*

___

##  protocolFeePaid

• **protocolFeePaid**: *string*

*Defined in [types.ts:389](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L389)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:383](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L383)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:382](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L382)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:392](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L392)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *string*

*Defined in [types.ts:386](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L386)*

___

##  takerFeeAssetData

• **takerFeeAssetData**: *string*

*Defined in [types.ts:394](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L394)*

___

##  takerFeePaid

• **takerFeePaid**: *string*

*Defined in [types.ts:388](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L388)*


<hr />

# Interface: WrapperGetOrdersResponse

### Hierarchy

* **WrapperGetOrdersResponse**


### Properties

##  ordersInfos

• **ordersInfos**: *[WrapperOrderInfo](#interface-wrapperorderinfo)[]*

*Defined in [types.ts:12](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L12)*

___

##  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:10](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L10)*

___

##  snapshotTimestamp

• **snapshotTimestamp**: *string*

*Defined in [types.ts:11](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L11)*


<hr />

# Interface: WrapperOrderEvent

### Hierarchy

* **WrapperOrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[WrapperContractEvent](#interface-wrappercontractevent)[]*

*Defined in [types.ts:526](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L526)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:524](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L524)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:525](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L525)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:522](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L522)*

___

##  signedOrder

• **signedOrder**: *[WrapperSignedOrder](#interface-wrappersignedorder)*

*Defined in [types.ts:523](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L523)*

___

##  timestamp

• **timestamp**: *string*

*Defined in [types.ts:521](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L521)*


<hr />

# Interface: WrapperOrderInfo

### Hierarchy

* **WrapperOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:24](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L24)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:22](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L22)*

___

##  signedOrder

• **signedOrder**: *[WrapperSignedOrder](#interface-wrappersignedorder)*

*Defined in [types.ts:23](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L23)*


<hr />

# Interface: WrapperRejectedOrderInfo

### Hierarchy

* **WrapperRejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [types.ts:560](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L560)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:558](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L558)*

___

##  signedOrder

• **signedOrder**: *[WrapperSignedOrder](#interface-wrappersignedorder)*

*Defined in [types.ts:559](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L559)*

___

##  status

• **status**: *[RejectedOrderStatus](#interface-rejectedorderstatus)*

*Defined in [types.ts:561](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L561)*


<hr />

# Interface: WrapperSignedOrder

### Hierarchy

* **WrapperSignedOrder**


### Properties

##  chainId

• **chainId**: *number*

*Defined in [types.ts:268](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L268)*

___

##  exchangeAddress

• **exchangeAddress**: *string*

*Defined in [types.ts:267](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L267)*

___

##  expirationTimeSeconds

• **expirationTimeSeconds**: *string*

*Defined in [types.ts:264](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L264)*

___

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:263](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L263)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:252](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L252)*

___

##  makerAssetAmount

• **makerAssetAmount**: *string*

*Defined in [types.ts:254](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L254)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:253](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L253)*

___

##  makerFee

• **makerFee**: *string*

*Defined in [types.ts:255](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L255)*

___

##  makerFeeAssetData

• **makerFeeAssetData**: *string*

*Defined in [types.ts:256](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L256)*

___

##  salt

• **salt**: *string*

*Defined in [types.ts:265](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L265)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:262](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L262)*

___

##  signature

• **signature**: *string*

*Defined in [types.ts:266](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L266)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:257](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L257)*

___

##  takerAssetAmount

• **takerAssetAmount**: *string*

*Defined in [types.ts:260](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L260)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:258](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L258)*

___

##  takerFee

• **takerFee**: *string*

*Defined in [types.ts:261](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L261)*

___

##  takerFeeAssetData

• **takerFeeAssetData**: *string*

*Defined in [types.ts:259](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L259)*


<hr />

# Interface: WrapperStats

### Hierarchy

* **WrapperStats**


### Properties

##  ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: *number*

*Defined in [types.ts:631](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L631)*

___

##  ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: *number*

*Defined in [types.ts:630](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L630)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:622](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L622)*

___

##  latestBlock

• **latestBlock**: *[LatestBlock](#interface-latestblock)*

*Defined in [types.ts:623](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L623)*

___

##  maxExpirationTime

• **maxExpirationTime**: *string*

*Defined in [types.ts:628](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L628)*

___

##  numOrders

• **numOrders**: *number*

*Defined in [types.ts:625](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L625)*

___

##  numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: *number*

*Defined in [types.ts:626](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L626)*

___

##  numPeers

• **numPeers**: *number*

*Defined in [types.ts:624](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L624)*

___

##  numPinnedOrders

• **numPinnedOrders**: *number*

*Defined in [types.ts:627](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L627)*

___

##  peerID

• **peerID**: *string*

*Defined in [types.ts:621](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L621)*

___

##  pubSubTopic

• **pubSubTopic**: *string*

*Defined in [types.ts:618](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L618)*

___

##  rendezvous

• **rendezvous**: *string*

*Defined in [types.ts:619](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L619)*

___

##  secondaryRendezvous

• **secondaryRendezvous**: *string[]*

*Defined in [types.ts:620](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L620)*

___

##  startOfCurrentUTCDay

• **startOfCurrentUTCDay**: *string*

*Defined in [types.ts:629](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L629)*

___

##  version

• **version**: *string*

*Defined in [types.ts:617](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L617)*


<hr />

# Interface: WrapperValidationResults

### Hierarchy

* **WrapperValidationResults**


### Properties

##  accepted

• **accepted**: *[WrapperAcceptedOrderInfo](#interface-wrapperacceptedorderinfo)[]*

*Defined in [types.ts:544](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L544)*

___

##  rejected

• **rejected**: *[WrapperRejectedOrderInfo](#interface-wrapperrejectedorderinfo)[]*

*Defined in [types.ts:545](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L545)*


<hr />

# Interface: WrapperWethDepositEvent

### Hierarchy

* **WrapperWethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:434](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L434)*

___

##  value

• **value**: *string*

*Defined in [types.ts:435](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L435)*


<hr />

# Interface: WrapperWethWithdrawalEvent

### Hierarchy

* **WrapperWethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L424)*

___

##  value

• **value**: *string*

*Defined in [types.ts:425](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L425)*


<hr />

# Interface: ZeroExMesh

### Hierarchy

* **ZeroExMesh**


### Methods

##  newWrapperAsync

▸ **newWrapperAsync**(`config`: [WrapperConfig](#interface-wrapperconfig)): *Promise‹[MeshWrapper](#interface-meshwrapper)›*

*Defined in [types.ts:215](https://github.com/0xProject/0x-mesh/blob/e3d1604/packages/browser-lite/src/types.ts#L215)*

**Parameters:**

Name | Type |
------ | ------ |
`config` | [WrapperConfig](#interface-wrapperconfig) |

**Returns:** *Promise‹[MeshWrapper](#interface-meshwrapper)›*


<hr />

