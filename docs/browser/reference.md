# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:542](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L542)*

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

*Defined in [index.ts:616](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L616)*

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

*Defined in [index.ts:562](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L562)*

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

*Defined in [index.ts:577](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L577)*

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

*Defined in [index.ts:588](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L588)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:407](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L407)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:410](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L410)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:411](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L411)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:414](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L414)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:408](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L408)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:409](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L409)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:406](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L406)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [index.ts:415](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L415)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [index.ts:412](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L412)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:413](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L413)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:496](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L496)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:494](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L494)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:495](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L495)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:493](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L493)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:118](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L118)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:115](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L115)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:114](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L114)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:117](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L117)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:113](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L113)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:119](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L119)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:116](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L116)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:474](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L474)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:475](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L475)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:472](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L472)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:473](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L473)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:56](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L56)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:49](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L49)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:92](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L92)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [index.ts:41](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L41)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:65](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L65)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [index.ts:70](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L70)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [index.ts:76](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L76)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:38](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L38)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [index.ts:97](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L97)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:44](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L44)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:35](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L35)*

<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:106](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L106)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:107](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L107)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:102](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L102)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:105](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L105)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:103](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L103)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:104](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L104)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:101](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L101)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:108](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L108)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:109](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L109)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:388](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L388)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:383](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L383)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:387](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L387)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:389](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L389)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:386](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L386)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:390](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L390)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:384](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L384)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:385](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L385)*

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:264](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L264)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:263](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L263)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:262](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L262)*

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:247](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L247)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:249](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L249)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:246](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L246)*

___

##  to

• **to**: *string*

*Defined in [index.ts:248](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L248)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:250](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L250)*

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:231](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L231)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:233](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L233)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:230](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L230)*

___

##  to

• **to**: *string*

*Defined in [index.ts:232](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L232)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:234](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L234)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:188](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L188)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:189](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L189)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:190](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L190)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:176](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L176)*

___

##  to

• **to**: *string*

*Defined in [index.ts:177](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L177)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:178](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L178)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:213](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L213)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:212](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L212)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:214](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L214)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:226](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L226)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:225](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L225)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:224](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L224)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:200](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L200)*

___

##  to

• **to**: *string*

*Defined in [index.ts:201](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L201)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:202](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L202)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:298](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L298)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:296](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L296)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:300](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L300)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:299](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L299)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:297](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L297)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:301](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L301)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:305](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L305)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:307](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L307)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:306](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L306)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:271](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L271)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:268](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L268)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:277](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L277)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:272](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L272)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:274](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L274)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:276](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L276)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:270](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L270)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:269](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L269)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:278](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L278)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:273](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L273)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:275](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L275)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:435](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L435)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:433](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L433)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:434](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L434)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:431](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L431)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:432](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L432)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:485](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L485)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:483](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L483)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:484](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L484)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:486](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L486)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:503](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L503)*

___

##  message

• **message**: *string*

*Defined in [index.ts:504](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L504)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:464](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L464)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:465](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L465)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:327](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L327)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:328](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L328)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:317](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L317)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:318](https://github.com/0xProject/0x-mesh/blob/2053ba8/browser/ts/index.ts#L318)*

<hr />

