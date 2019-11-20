# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:562](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L562)*

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

*Defined in [index.ts:636](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L636)*

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

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [index.ts:582](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L582)*

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

*Defined in [index.ts:597](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L597)*

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

*Defined in [index.ts:608](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L608)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:427](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L427)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:430](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L430)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:431](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L431)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:434](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L434)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:428](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L428)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:429](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L429)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:426](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L426)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:435](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L435)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [index.ts:432](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L432)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:433](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L433)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:516](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L516)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:514](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L514)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:515](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L515)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:513](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L513)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:141](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L141)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:138](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L138)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:137](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L137)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:140](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L140)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:136](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L136)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:142](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L142)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:139](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L139)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:494](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L494)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:495](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L495)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:492](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L492)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:493](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L493)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:79](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L79)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:72](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L72)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:115](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L115)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:64](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L64)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:88](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L88)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:93](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L93)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:99](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L99)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:61](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L61)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:120](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L120)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:67](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L67)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:58](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L58)*

<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:129](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L129)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:130](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L130)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:125](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L125)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:128](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L128)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:126](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L126)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:127](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L127)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:124](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L124)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:131](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L131)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:132](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L132)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:408](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L408)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:403](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L403)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:407](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L407)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:409](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L409)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:406](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L406)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:410](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L410)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:404](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L404)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:405](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L405)*

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:284](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L284)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:283](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L283)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:282](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L282)*

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:267](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L267)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:269](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L269)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:266](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L266)*

___

##  to

• **to**: *string*

*Defined in [index.ts:268](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L268)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:270](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L270)*

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:251](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L251)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:253](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L253)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:250](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L250)*

___

##  to

• **to**: *string*

*Defined in [index.ts:252](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L252)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:254](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L254)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:208](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L208)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:209](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L209)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:210](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L210)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:196](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L196)*

___

##  to

• **to**: *string*

*Defined in [index.ts:197](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L197)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:198](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L198)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:233](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L233)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:232](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L232)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:234](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L234)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:246](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L246)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:245](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L245)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:244](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L244)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:220](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L220)*

___

##  to

• **to**: *string*

*Defined in [index.ts:221](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L221)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:222](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L222)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:318](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L318)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:316](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L316)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:320](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L320)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:319](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L319)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:317](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L317)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:321](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L321)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:325](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L325)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:327](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L327)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:326](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L326)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:291](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L291)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:288](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L288)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:297](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L297)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:292](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L292)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:294](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L294)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:296](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L296)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:290](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L290)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:289](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L289)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:298](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L298)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:293](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L293)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:295](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L295)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:455](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L455)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:453](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L453)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:454](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L454)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:451](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L451)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:452](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L452)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:505](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L505)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:503](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L503)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:504](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L504)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:506](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L506)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:523](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L523)*

___

##  message

• **message**: *string*

*Defined in [index.ts:524](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L524)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:484](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L484)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:485](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L485)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:347](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L347)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:348](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L348)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:337](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L337)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:338](https://github.com/0xProject/0x-mesh/blob/3a63262/browser/ts/index.ts#L338)*

<hr />

