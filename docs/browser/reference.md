# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:641](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L641)*

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

*Defined in [index.ts:791](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L791)*

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

*Defined in [index.ts:720](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L720)*

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

*Defined in [index.ts:762](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L762)*

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

*Defined in [index.ts:703](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L703)*

Returns various stats about Mesh, including the total number of orders
and the number of peers Mesh is connected to.

**Returns:** *Promise‹[Stats](#interface-stats)›*

___

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [index.ts:661](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L661)*

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

*Defined in [index.ts:676](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L676)*

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

*Defined in [index.ts:687](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L687)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*


<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:504](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L504)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:507](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L507)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:508](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L508)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:511](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L511)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:505](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L505)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:506](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L506)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:503](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L503)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:512](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L512)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [index.ts:509](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L509)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:510](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L510)*


<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:595](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L595)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:593](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L593)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:594](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L594)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:592](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L592)*


<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:212](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L212)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:209](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L209)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:208](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L208)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:211](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L211)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:207](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L207)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:213](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L213)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:210](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L210)*


<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:573](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L573)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:574](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L574)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:571](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L571)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:572](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L572)*


<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:79](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L79)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:72](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L72)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:123](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L123)*

___

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : *undefined | false | true*

*Defined in [index.ts:96](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L96)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:64](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L64)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:88](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L88)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:101](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L101)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:107](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L107)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:61](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L61)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:128](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L128)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:67](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L67)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:58](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L58)*


<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:137](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L137)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:138](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L138)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:133](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L133)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:136](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L136)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:134](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L134)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:135](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L135)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:132](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L132)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:139](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L139)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:140](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L140)*


<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:485](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L485)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:480](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L480)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:484](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L484)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:486](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L486)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:483](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L483)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:487](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L487)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:481](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L481)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:482](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L482)*


<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:361](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L361)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:360](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L360)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:359](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L359)*


<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:344](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L344)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:346](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L346)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:343](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L343)*

___

##  to

• **to**: *string*

*Defined in [index.ts:345](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L345)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:347](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L347)*


<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:328](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L328)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:330](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L330)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:327](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L327)*

___

##  to

• **to**: *string*

*Defined in [index.ts:329](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L329)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:331](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L331)*


<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:285](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L285)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:286](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L286)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:287](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L287)*


<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:273](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L273)*

___

##  to

• **to**: *string*

*Defined in [index.ts:274](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L274)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:275](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L275)*


<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:310](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L310)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:309](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L309)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:311](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L311)*


<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:323](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L323)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:322](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L322)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:321](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L321)*


<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:297](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L297)*

___

##  to

• **to**: *string*

*Defined in [index.ts:298](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L298)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:299](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L299)*


<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:395](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L395)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:393](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L393)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:397](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L397)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:396](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L396)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:394](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L394)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:398](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L398)*


<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:402](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L402)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:404](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L404)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:403](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L403)*


<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:368](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L368)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:365](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L365)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:374](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L374)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:369](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L369)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:371](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L371)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:373](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L373)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:367](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L367)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:366](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L366)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:375](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L375)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:370](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L370)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:372](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L372)*


<hr />

# Interface: GetOrdersResponse

### Hierarchy

* **GetOrdersResponse**


### Properties

##  ordersInfos

• **ordersInfos**: *[OrderInfo](#class-orderinfo)[]*

*Defined in [index.ts:203](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L203)*

___

##  snapshotID

• **snapshotID**: *string*

*Defined in [index.ts:201](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L201)*

___

##  snapshotTimestamp

• **snapshotTimestamp**: *number*

*Defined in [index.ts:202](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L202)*


<hr />

# Interface: LatestBlock

### Hierarchy

* **LatestBlock**


### Properties

##  hash

• **hash**: *string*

*Defined in [index.ts:145](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L145)*

___

##  number

• **number**: *number*

*Defined in [index.ts:144](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L144)*


<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:534](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L534)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:532](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L532)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:533](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L533)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:530](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L530)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:531](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L531)*

___

##  timestampMs

• **timestampMs**: *number*

*Defined in [index.ts:529](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L529)*


<hr />

# Interface: OrderInfo

### Hierarchy

* **OrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:191](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L191)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:189](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L189)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:190](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L190)*


<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:584](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L584)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:582](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L582)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:583](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L583)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:585](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L585)*


<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:602](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L602)*

___

##  message

• **message**: *string*

*Defined in [index.ts:603](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L603)*


<hr />

# Interface: Stats

### Hierarchy

* **Stats**


### Properties

##  ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: *number*

*Defined in [index.ts:179](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L179)*

___

##  ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: *number*

*Defined in [index.ts:178](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L178)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:170](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L170)*

___

##  latestBlock

• **latestBlock**: *[LatestBlock](#class-latestblock)*

*Defined in [index.ts:171](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L171)*

___

##  maxExpirationTime

• **maxExpirationTime**: *BigNumber*

*Defined in [index.ts:176](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L176)*

___

##  numOrders

• **numOrders**: *number*

*Defined in [index.ts:173](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L173)*

___

##  numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: *number*

*Defined in [index.ts:174](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L174)*

___

##  numPeers

• **numPeers**: *number*

*Defined in [index.ts:172](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L172)*

___

##  numPinnedOrders

• **numPinnedOrders**: *number*

*Defined in [index.ts:175](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L175)*

___

##  peerID

• **peerID**: *string*

*Defined in [index.ts:169](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L169)*

___

##  pubSubTopic

• **pubSubTopic**: *string*

*Defined in [index.ts:167](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L167)*

___

##  rendezvous

• **rendezvous**: *string*

*Defined in [index.ts:168](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L168)*

___

##  startOfCurrentUTCDay

• **startOfCurrentUTCDay**: *Date*

*Defined in [index.ts:177](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L177)*

___

##  version

• **version**: *string*

*Defined in [index.ts:166](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L166)*


<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:563](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L563)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:564](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L564)*


<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:424](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L424)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:425](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L425)*


<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:414](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L414)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:415](https://github.com/0xProject/0x-mesh/blob/ae6de374/browser/ts/index.ts#L415)*


<hr />

