# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:576](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L576)*

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

*Defined in [index.ts:650](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L650)*

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

*Defined in [index.ts:596](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L596)*

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

*Defined in [index.ts:611](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L611)*

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

*Defined in [index.ts:622](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L622)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*


<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:439](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L439)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:442](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L442)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:443](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L443)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:446](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L446)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:440](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L440)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:441](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L441)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:438](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L438)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:447](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L447)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [index.ts:444](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L444)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:445](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L445)*


<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:530](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L530)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:528](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L528)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:529](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L529)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:527](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L527)*


<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:149](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L149)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:146](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L146)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:145](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L145)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:148](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L148)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:144](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L144)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:150](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L150)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:147](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L147)*


<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:508](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L508)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:509](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L509)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:506](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L506)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:507](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L507)*


<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:79](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L79)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:72](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L72)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:123](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L123)*

___

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : *undefined | false | true*

*Defined in [index.ts:96](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L96)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:64](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L64)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:88](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L88)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:101](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L101)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:107](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L107)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:61](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L61)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:128](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L128)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:67](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L67)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:58](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L58)*


<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:137](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L137)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:138](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L138)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:133](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L133)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:136](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L136)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:134](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L134)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:135](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L135)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:132](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L132)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:139](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L139)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:140](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L140)*


<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:420](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L420)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:415](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L415)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:419](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L419)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:421](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L421)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:418](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L418)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:422](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L422)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:416](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L416)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:417](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L417)*


<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:296](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L296)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:295](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L295)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:294](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L294)*


<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:279](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L279)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:281](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L281)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:278](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L278)*

___

##  to

• **to**: *string*

*Defined in [index.ts:280](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L280)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:282](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L282)*


<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:263](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L263)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:265](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L265)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:262](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L262)*

___

##  to

• **to**: *string*

*Defined in [index.ts:264](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L264)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:266](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L266)*


<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:220](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L220)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:221](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L221)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:222](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L222)*


<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:208](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L208)*

___

##  to

• **to**: *string*

*Defined in [index.ts:209](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L209)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:210](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L210)*


<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:245](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L245)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:244](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L244)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:246](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L246)*


<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:258](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L258)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:257](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L257)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:256](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L256)*


<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:232](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L232)*

___

##  to

• **to**: *string*

*Defined in [index.ts:233](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L233)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:234](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L234)*


<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:330](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L330)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:328](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L328)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:332](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L332)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:331](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L331)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:329](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L329)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:333](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L333)*


<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:337](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L337)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:339](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L339)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:338](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L338)*


<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:303](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L303)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:300](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L300)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:309](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L309)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:304](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L304)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:306](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L306)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:308](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L308)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:302](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L302)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:301](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L301)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:310](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L310)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:305](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L305)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:307](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L307)*


<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:469](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L469)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:467](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L467)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:468](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L468)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:465](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L465)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:466](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L466)*

___

##  timestampMs

• **timestampMs**: *number*

*Defined in [index.ts:464](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L464)*


<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:519](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L519)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:517](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L517)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:518](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L518)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:520](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L520)*


<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:537](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L537)*

___

##  message

• **message**: *string*

*Defined in [index.ts:538](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L538)*


<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:498](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L498)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:499](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L499)*


<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:359](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L359)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:360](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L360)*


<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:349](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L349)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:350](https://github.com/0xProject/0x-mesh/blob/1593c46/browser/ts/index.ts#L350)*


<hr />

