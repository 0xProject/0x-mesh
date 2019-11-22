# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:565](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L565)*

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

*Defined in [index.ts:639](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L639)*

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

*Defined in [index.ts:585](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L585)*

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

*Defined in [index.ts:600](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L600)*

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

*Defined in [index.ts:611](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L611)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:430](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L430)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:433](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L433)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:434](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L434)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:437](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L437)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:431](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L431)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:432](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L432)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:429](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L429)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:438](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L438)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [index.ts:435](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L435)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:436](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L436)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:519](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L519)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:517](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L517)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:518](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L518)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:516](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L516)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:141](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L141)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:138](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L138)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:137](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L137)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:140](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L140)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:136](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L136)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:142](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L142)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:139](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L139)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:497](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L497)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:498](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L498)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:495](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L495)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:496](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L496)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:79](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L79)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:72](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L72)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:115](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L115)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:64](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L64)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:88](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L88)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:93](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L93)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:99](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L99)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:61](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L61)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:120](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L120)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:67](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L67)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:58](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L58)*

<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:129](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L129)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:130](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L130)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:125](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L125)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:128](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L128)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:126](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L126)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:127](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L127)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:124](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L124)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:131](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L131)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:132](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L132)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:411](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L411)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:406](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L406)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:410](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L410)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:412](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L412)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:409](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L409)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:413](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L413)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:407](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L407)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:408](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L408)*

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:287](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L287)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:286](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L286)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:285](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L285)*

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:270](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L270)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:272](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L272)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:269](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L269)*

___

##  to

• **to**: *string*

*Defined in [index.ts:271](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L271)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:273](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L273)*

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:254](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L254)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:256](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L256)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:253](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L253)*

___

##  to

• **to**: *string*

*Defined in [index.ts:255](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L255)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:257](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L257)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:211](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L211)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:212](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L212)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:213](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L213)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:199](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L199)*

___

##  to

• **to**: *string*

*Defined in [index.ts:200](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L200)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:201](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L201)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:236](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L236)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:235](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L235)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:237](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L237)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:249](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L249)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:248](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L248)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:247](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L247)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:223](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L223)*

___

##  to

• **to**: *string*

*Defined in [index.ts:224](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L224)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:225](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L225)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:321](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L321)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:319](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L319)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:323](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L323)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:322](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L322)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:320](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L320)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:324](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L324)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:328](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L328)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:330](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L330)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:329](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L329)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:294](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L294)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:291](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L291)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:300](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L300)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:295](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L295)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:297](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L297)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:299](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L299)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:293](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L293)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:292](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L292)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:301](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L301)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:296](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L296)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:298](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L298)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:458](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L458)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:456](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L456)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:457](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L457)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:454](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L454)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:455](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L455)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:508](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L508)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:506](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L506)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:507](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L507)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:509](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L509)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:526](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L526)*

___

##  message

• **message**: *string*

*Defined in [index.ts:527](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L527)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:487](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L487)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:488](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L488)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:350](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L350)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:351](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L351)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:340](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L340)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:341](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/browser/ts/index.ts#L341)*

<hr />

