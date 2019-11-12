# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:538](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L538)*

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

*Defined in [index.ts:612](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L612)*

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

*Defined in [index.ts:558](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L558)*

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

*Defined in [index.ts:573](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L573)*

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

*Defined in [index.ts:584](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L584)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:404](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L404)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:407](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L407)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:408](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L408)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:410](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L410)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:405](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L405)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:406](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L406)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:403](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L403)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:411](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L411)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:409](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L409)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:492](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L492)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:490](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L490)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:491](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L491)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:489](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L489)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:118](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L118)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:115](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L115)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:114](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L114)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:117](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L117)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:113](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L113)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:119](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L119)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:116](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L116)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:470](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L470)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:471](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L471)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:468](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L468)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:469](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L469)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:56](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L56)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:49](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L49)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:92](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L92)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:41](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L41)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:65](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L65)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:70](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L70)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:76](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L76)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:38](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L38)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:97](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L97)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:44](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L44)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:35](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L35)*

<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:106](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L106)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:107](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L107)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:102](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L102)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:105](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L105)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:103](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L103)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:104](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L104)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:101](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L101)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:108](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L108)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:109](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L109)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:385](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L385)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:380](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L380)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:384](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L384)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:386](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L386)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:383](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L383)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:387](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L387)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:381](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L381)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:382](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L382)*

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:261](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L261)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:260](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L260)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:259](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L259)*

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:244](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L244)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:246](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L246)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:243](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L243)*

___

##  to

• **to**: *string*

*Defined in [index.ts:245](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L245)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:247](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L247)*

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:228](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L228)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:230](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L230)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:227](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L227)*

___

##  to

• **to**: *string*

*Defined in [index.ts:229](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L229)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:231](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L231)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:185](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L185)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:186](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L186)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:187](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L187)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:173](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L173)*

___

##  to

• **to**: *string*

*Defined in [index.ts:174](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L174)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:175](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L175)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:210](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L210)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:209](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L209)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:211](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L211)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:223](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L223)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:222](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L222)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:221](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L221)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:197](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L197)*

___

##  to

• **to**: *string*

*Defined in [index.ts:198](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L198)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:199](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L199)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:295](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L295)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:293](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L293)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:297](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L297)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:296](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L296)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:294](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L294)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:298](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L298)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:302](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L302)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:304](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L304)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:303](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L303)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:268](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L268)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:265](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L265)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:274](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L274)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:269](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L269)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:271](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L271)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:273](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L273)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:267](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L267)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:266](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L266)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:275](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L275)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:270](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L270)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:272](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L272)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:431](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L431)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:429](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L429)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:430](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L430)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:427](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L427)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:428](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L428)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:481](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L481)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:479](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L479)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:480](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L480)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:482](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L482)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:499](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L499)*

___

##  message

• **message**: *string*

*Defined in [index.ts:500](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L500)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:460](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L460)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:461](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L461)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:324](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L324)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:325](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L325)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:314](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L314)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:315](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L315)*

<hr />

