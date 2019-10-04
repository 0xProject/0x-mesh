# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:425](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L425)*

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

*Defined in [index.ts:495](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L495)*

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

*Defined in [index.ts:445](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L445)*

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

*Defined in [index.ts:460](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L460)*

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

*Defined in [index.ts:471](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L471)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [index.ts:292](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L292)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [index.ts:295](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L295)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [index.ts:296](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L296)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [index.ts:298](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L298)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [index.ts:293](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L293)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [index.ts:294](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L294)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [index.ts:291](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L291)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [index.ts:297](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L297)*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:379](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L379)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:377](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L377)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:378](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L378)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:376](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L376)*

<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [index.ts:77](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L77)*

___

##  Error

• **Error**: = 2

*Defined in [index.ts:74](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L74)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [index.ts:73](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L73)*

___

##  Info

• **Info**: = 4

*Defined in [index.ts:76](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L76)*

___

##  Panic

• **Panic**: = 0

*Defined in [index.ts:72](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L72)*

___

##  Trace

• **Trace**: = 6

*Defined in [index.ts:78](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L78)*

___

##  Warn

• **Warn**: = 3

*Defined in [index.ts:75](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L75)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:357](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L357)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:358](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L358)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:355](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L355)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:356](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L356)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:59](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L59)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:48](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L48)*

___

##  ethereumNetworkID

• **ethereumNetworkID**: *number*

*Defined in [index.ts:40](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L40)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:68](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L68)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:38](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L38)*

___

## `Optional` orderExpirationBufferSeconds

• **orderExpirationBufferSeconds**? : *undefined | number*

*Defined in [index.ts:52](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L52)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:43](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L43)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [index.ts:35](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L35)*

<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [index.ts:273](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L273)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [index.ts:268](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L268)*

___

##  isRemoved

• **isRemoved**: *string*

*Defined in [index.ts:272](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L272)*

___

##  kind

• **kind**: *ContractEventKind*

*Defined in [index.ts:274](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L274)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [index.ts:271](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L271)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [index.ts:275](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L275)*

___

##  txHash

• **txHash**: *string*

*Defined in [index.ts:269](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L269)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [index.ts:270](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L270)*

<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:141](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L141)*

___

##  spender

• **spender**: *string*

*Defined in [index.ts:142](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L142)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:143](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L143)*

<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:129](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L129)*

___

##  to

• **to**: *string*

*Defined in [index.ts:130](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L130)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:131](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L131)*

<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [index.ts:166](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L166)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:165](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L165)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:167](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L167)*

<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [index.ts:179](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L179)*

___

##  operator

• **operator**: *string*

*Defined in [index.ts:178](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L178)*

___

##  owner

• **owner**: *string*

*Defined in [index.ts:177](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L177)*

<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [index.ts:153](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L153)*

___

##  to

• **to**: *string*

*Defined in [index.ts:154](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L154)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [index.ts:155](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L155)*

<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:213](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L213)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:211](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L211)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:215](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L215)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:214](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L214)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:212](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L212)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:216](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L216)*

<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:220](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L220)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [index.ts:222](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L222)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:221](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L221)*

<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [index.ts:186](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L186)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [index.ts:183](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L183)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [index.ts:192](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L192)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:187](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L187)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [index.ts:189](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L189)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:191](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L191)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [index.ts:185](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L185)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [index.ts:184](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L184)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [index.ts:193](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L193)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [index.ts:188](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L188)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [index.ts:190](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L190)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#class-contractevent)[]*

*Defined in [index.ts:318](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L318)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [index.ts:316](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L316)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:317](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L317)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:314](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L314)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:315](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L315)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:368](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L368)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:366](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L366)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:367](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L367)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:369](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L369)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:386](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L386)*

___

##  message

• **message**: *string*

*Defined in [index.ts:387](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L387)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:347](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L347)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:348](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L348)*

<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:242](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L242)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:243](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L243)*

<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [index.ts:232](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L232)*

___

##  value

• **value**: *BigNumber*

*Defined in [index.ts:233](https://github.com/0xProject/0x-mesh/blob/214e237/browser/ts/index.ts#L233)*

<hr />

