# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:573](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L573)*

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

*Defined in [index.ts:647](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L647)*

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

*Defined in [index.ts:593](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L593)*

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

*Defined in [index.ts:608](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L608)*

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

*Defined in [index.ts:619](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L619)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:436](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L436)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:439](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L439)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:440](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L440)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:443](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L443)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:437](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L437)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:438](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L438)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:435](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L435)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:444](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L444)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [index.ts:441](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L441)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:442](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L442)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:527](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L527)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:525](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L525)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:526](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L526)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:524](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L524)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:149](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L149)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:146](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L146)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:145](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L145)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:148](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L148)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:144](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L144)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:150](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L150)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:147](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L147)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:505](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L505)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:506](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L506)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:503](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L503)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:504](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L504)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:79](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L79)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:72](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L72)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:123](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L123)*

___

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : *undefined | false | true*

*Defined in [index.ts:96](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L96)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:64](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L64)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:88](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L88)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:101](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L101)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:107](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L107)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:61](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L61)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:128](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L128)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:67](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L67)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:58](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L58)*

<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:137](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L137)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:138](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L138)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:133](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L133)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:136](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L136)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:134](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L134)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:135](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L135)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:132](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L132)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:139](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L139)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:140](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L140)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:417](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L417)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:412](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L412)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:416](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L416)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:418](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L418)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:415](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L415)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:419](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L419)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:413](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L413)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:414](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L414)*

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:293](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L293)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:292](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L292)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:291](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L291)*

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:276](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L276)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:278](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L278)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:275](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L275)*

___

##  to

• **to**: *string*

*Defined in [index.ts:277](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L277)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:279](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L279)*

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:260](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L260)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:262](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L262)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:259](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L259)*

___

##  to

• **to**: *string*

*Defined in [index.ts:261](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L261)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:263](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L263)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:217](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L217)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:218](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L218)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:219](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L219)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:205](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L205)*

___

##  to

• **to**: *string*

*Defined in [index.ts:206](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L206)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:207](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L207)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:242](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L242)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:241](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L241)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:243](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L243)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:255](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L255)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:254](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L254)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:253](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L253)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:229](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L229)*

___

##  to

• **to**: *string*

*Defined in [index.ts:230](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L230)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:231](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L231)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:327](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L327)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:325](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L325)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:329](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L329)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:328](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L328)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:326](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L326)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:330](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L330)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:334](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L334)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:336](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L336)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:335](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L335)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:300](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L300)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:297](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L297)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:306](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L306)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:301](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L301)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:303](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L303)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:305](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L305)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:299](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L299)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:298](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L298)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:307](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L307)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:302](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L302)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:304](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L304)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:466](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L466)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:464](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L464)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:465](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L465)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:462](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L462)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:463](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L463)*

___

##  timestampMs

• **timestampMs**: *number*

*Defined in [index.ts:461](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L461)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:516](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L516)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:514](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L514)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:515](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L515)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:517](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L517)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:534](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L534)*

___

##  message

• **message**: *string*

*Defined in [index.ts:535](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L535)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:495](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L495)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:496](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L496)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:356](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L356)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:357](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L357)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:346](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L346)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:347](https://github.com/0xProject/0x-mesh/blob/80bd6ae/browser/ts/index.ts#L347)*

<hr />

