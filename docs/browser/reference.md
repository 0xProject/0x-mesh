# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:524](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L524)*

Instantiates a new Mesh instance.

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`config` | [Config](#interface-config) | Configuration options for Mesh |

**Returns:** *[Mesh](#class-mesh)*

An instance of Mesh

### Methods

##  addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[]): *Promise‹[ValidationResults](#interface-validationresults)›*

*Defined in [index.ts:594](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L594)*

Validates and adds the given orders to Mesh. If an order is successfully
added, Mesh will share it with any peers in the network and start
watching it for changes (e.g. filled, canceled, expired). The returned
promise will only be rejected if there was an error validating or adding
the order; it will not be rejected for any invalid orders (check
results.rejected instead).

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`orders` | SignedOrder[] | An array of orders to add. |

**Returns:** *Promise‹[ValidationResults](#interface-validationresults)›*

Validation results for the given orders, indicating which orders
were accepted and which were rejected.

___

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [index.ts:544](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L544)*

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

*Defined in [index.ts:559](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L559)*

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

*Defined in [index.ts:570](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L570)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:391](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L391)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:394](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L394)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:395](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L395)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:397](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L397)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:392](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L392)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:393](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L393)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:390](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L390)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:396](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L396)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:478](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L478)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:476](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L476)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:477](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L477)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:475](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L475)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:106](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L106)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:103](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L103)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:102](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L102)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:105](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L105)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:101](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L101)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:107](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L107)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:104](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L104)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:456](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L456)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:457](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L457)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:454](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L454)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:455](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L455)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:59](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L59)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:48](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L48)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#class-contractaddresses)*

*Defined in [index.ts:85](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L85)*

___

##  ethereumNetworkID

• **ethereumNetworkID**: *number*

*Defined in [index.ts:40](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L40)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:68](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L68)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:38](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L38)*

___

## `Optional` orderExpirationBufferSeconds

• **orderExpirationBufferSeconds**? : *undefined | number*

*Defined in [index.ts:52](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L52)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:43](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L43)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:35](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L35)*

<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [index.ts:94](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L94)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [index.ts:95](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L95)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [index.ts:90](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L90)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [index.ts:93](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L93)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [index.ts:91](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L91)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [index.ts:92](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L92)*

___

##  exchange

• **exchange**: *string*

*Defined in [index.ts:89](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L89)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [index.ts:96](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L96)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [index.ts:97](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L97)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:372](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L372)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:367](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L367)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:371](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L371)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:373](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L373)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:370](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L370)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:374](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L374)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:368](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L368)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:369](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L369)*

<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:248](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L248)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:247](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L247)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:246](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L246)*

<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:231](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L231)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [index.ts:233](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L233)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:230](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L230)*

___

##  to

• **to**: *string*

*Defined in [index.ts:232](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L232)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [index.ts:234](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L234)*

<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:215](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L215)*

___

##  id

• **id**: *BigNumber*

*Defined in [index.ts:217](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L217)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:214](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L214)*

___

##  to

• **to**: *string*

*Defined in [index.ts:216](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L216)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:218](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L218)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:171](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L171)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:172](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L172)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:173](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L173)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:159](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L159)*

___

##  to

• **to**: *string*

*Defined in [index.ts:160](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L160)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:161](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L161)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:196](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L196)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:195](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L195)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:197](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L197)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:209](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L209)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:208](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L208)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:207](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L207)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:183](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L183)*

___

##  to

• **to**: *string*

*Defined in [index.ts:184](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L184)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:185](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L185)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:282](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L282)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:280](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L280)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:284](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L284)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:283](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L283)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:281](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L281)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:285](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L285)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:289](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L289)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:291](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L291)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:290](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L290)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:255](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L255)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:252](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L252)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:261](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L261)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:256](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L256)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:258](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L258)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:260](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L260)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:254](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L254)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:253](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L253)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:262](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L262)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:257](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L257)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:259](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L259)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:417](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L417)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:415](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L415)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:416](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L416)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:413](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L413)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:414](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L414)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:467](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L467)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:465](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L465)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:466](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L466)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:468](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L468)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:485](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L485)*

___

##  message

• **message**: *string*

*Defined in [index.ts:486](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L486)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:446](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L446)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:447](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L447)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:311](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L311)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:312](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L312)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:301](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L301)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:302](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L302)*

<hr />

