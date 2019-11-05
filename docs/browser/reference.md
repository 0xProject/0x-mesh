# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:536](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L536)*

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

*Defined in [index.ts:610](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L610)*

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

*Defined in [index.ts:556](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L556)*

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

*Defined in [index.ts:571](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L571)*

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

*Defined in [index.ts:582](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L582)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:402](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L402)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:405](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L405)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:406](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L406)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:408](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L408)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:403](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L403)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:404](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L404)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:401](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L401)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:409](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L409)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:407](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L407)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:490](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L490)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:488](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L488)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:489](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L489)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:487](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L487)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:116](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L116)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:113](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L113)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:112](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L112)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:115](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L115)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:111](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L111)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:117](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L117)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:114](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L114)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:468](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L468)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:469](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L469)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:466](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L466)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:467](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L467)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` EthereumRPCMaxRequestsPer24HrUTC

• **EthereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:69](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L69)*

___

## `Optional` EthereumRPCMaxRequestsPerSecond

• **EthereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:74](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L74)*

___

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:56](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L56)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:49](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L49)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:90](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L90)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:41](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L41)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:65](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L65)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:38](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L38)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:95](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L95)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:44](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L44)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:35](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L35)*

<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:104](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L104)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:105](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L105)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:100](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L100)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:103](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L103)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:101](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L101)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:102](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L102)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:99](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L99)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:106](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L106)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:107](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L107)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:383](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L383)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:378](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L378)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:382](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L382)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:384](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L384)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:381](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L381)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:385](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L385)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:379](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L379)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:380](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L380)*

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:259](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L259)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:258](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L258)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:257](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L257)*

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:242](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L242)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:244](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L244)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:241](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L241)*

___

##  to

• **to**: *string*

*Defined in [index.ts:243](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L243)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:245](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L245)*

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:226](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L226)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:228](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L228)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:225](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L225)*

___

##  to

• **to**: *string*

*Defined in [index.ts:227](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L227)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:229](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L229)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:183](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L183)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:184](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L184)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:185](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L185)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:171](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L171)*

___

##  to

• **to**: *string*

*Defined in [index.ts:172](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L172)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:173](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L173)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:208](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L208)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:207](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L207)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:209](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L209)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:221](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L221)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:220](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L220)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:219](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L219)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:195](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L195)*

___

##  to

• **to**: *string*

*Defined in [index.ts:196](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L196)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:197](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L197)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:293](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L293)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:291](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L291)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:295](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L295)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:294](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L294)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:292](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L292)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:296](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L296)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:300](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L300)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:302](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L302)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:301](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L301)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:266](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L266)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:263](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L263)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:272](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L272)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:267](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L267)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:269](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L269)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:271](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L271)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:265](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L265)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:264](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L264)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:273](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L273)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:268](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L268)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:270](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L270)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:429](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L429)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:427](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L427)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:428](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L428)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:425](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L425)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:426](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L426)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:479](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L479)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:477](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L477)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:478](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L478)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:480](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L480)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:497](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L497)*

___

##  message

• **message**: *string*

*Defined in [index.ts:498](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L498)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:458](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L458)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:459](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L459)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:322](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L322)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:323](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L323)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:312](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L312)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:313](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L313)*

<hr />

