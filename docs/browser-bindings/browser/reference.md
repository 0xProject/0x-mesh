# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [mesh.ts:142](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L142)*

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

*Defined in [mesh.ts:292](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L292)*

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

*Defined in [mesh.ts:221](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L221)*

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

*Defined in [mesh.ts:263](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L263)*

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

*Defined in [mesh.ts:204](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L204)*

Returns various stats about Mesh, including the total number of orders
and the number of peers Mesh is connected to.

**Returns:** *Promise‹[Stats](#interface-stats)›*

___

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [mesh.ts:162](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L162)*

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

*Defined in [mesh.ts:177](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L177)*

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

*Defined in [mesh.ts:188](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/mesh.ts#L188)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*


<hr />

# Enumeration: ContractEventKind


### Enumeration members

##  ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

*Defined in [types.ts:468](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L468)*

___

##  ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

*Defined in [types.ts:470](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L470)*

___

##  ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

*Defined in [types.ts:469](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L469)*

___

##  ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

*Defined in [types.ts:464](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L464)*

___

##  ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

*Defined in [types.ts:463](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L463)*

___

##  ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

*Defined in [types.ts:466](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L466)*

___

##  ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

*Defined in [types.ts:467](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L467)*

___

##  ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

*Defined in [types.ts:465](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L465)*

___

##  ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

*Defined in [types.ts:472](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L472)*

___

##  ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

*Defined in [types.ts:473](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L473)*

___

##  ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

*Defined in [types.ts:471](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L471)*

___

##  WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

*Defined in [types.ts:474](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L474)*

___

##  WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

*Defined in [types.ts:475](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L475)*


<hr />

# Enumeration: OrderEventEndState


### Enumeration members

##  Added

• **Added**: = "ADDED"

*Defined in [types.ts:538](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L538)*

___

##  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [types.ts:541](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L541)*

___

##  Expired

• **Expired**: = "EXPIRED"

*Defined in [types.ts:542](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L542)*

___

##  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [types.ts:545](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L545)*

___

##  Filled

• **Filled**: = "FILLED"

*Defined in [types.ts:539](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L539)*

___

##  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [types.ts:540](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L540)*

___

##  Invalid

• **Invalid**: = "INVALID"

*Defined in [types.ts:537](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L537)*

___

##  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [types.ts:546](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L546)*

___

##  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [types.ts:543](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L543)*

___

##  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [types.ts:544](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L544)*


<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [types.ts:630](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L630)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [types.ts:628](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L628)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [types.ts:629](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L629)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [types.ts:627](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L627)*


<hr />

# Enumeration: Verbosity


### Enumeration members

##  Debug

• **Debug**: = 5

*Defined in [types.ts:211](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L211)*

___

##  Error

• **Error**: = 2

*Defined in [types.ts:208](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L208)*

___

##  Fatal

• **Fatal**: = 1

*Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L207)*

___

##  Info

• **Info**: = 4

*Defined in [types.ts:210](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L210)*

___

##  Panic

• **Panic**: = 0

*Defined in [types.ts:206](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L206)*

___

##  Trace

• **Trace**: = 6

*Defined in [types.ts:212](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L212)*

___

##  Warn

• **Warn**: = 3

*Defined in [types.ts:209](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L209)*


<hr />

# Interface: SchemaValidationResult

### Hierarchy

* **SchemaValidationResult**


### Properties

##  errors

• **errors**: *string[]*

*Defined in [schema_validator.ts:11](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/schema_validator.ts#L11)*

___

## `Optional` fatal

• **fatal**? : *undefined | string*

*Defined in [schema_validator.ts:12](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/schema_validator.ts#L12)*

___

##  success

• **success**: *boolean*

*Defined in [schema_validator.ts:10](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/schema_validator.ts#L10)*


<hr />

# Interface: SchemaValidator

### Hierarchy

* **SchemaValidator**


### Properties

##  messageValidator

• **messageValidator**: *function*

*Defined in [schema_validator.ts:17](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/schema_validator.ts#L17)*

#### Type declaration:

▸ (`input`: string): *[SchemaValidationResult](#class-schemavalidationresult)*

**Parameters:**

Name | Type |
------ | ------ |
`input` | string |

___

##  orderValidator

• **orderValidator**: *function*

*Defined in [schema_validator.ts:16](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/schema_validator.ts#L16)*

#### Type declaration:

▸ (`input`: string): *[SchemaValidationResult](#class-schemavalidationresult)*

**Parameters:**

Name | Type |
------ | ------ |
`input` | string |


<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:608](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L608)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [types.ts:609](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L609)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:606](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L606)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:607](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L607)*


<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [types.ts:118](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L118)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [types.ts:111](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L111)*

___

## `Optional` customContractAddresses

• **customContractAddresses**? : *[ContractAddresses](#interface-contractaddresses)*

*Defined in [types.ts:162](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L162)*

___

## `Optional` customOrderFilter

• **customOrderFilter**? : *[JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:187](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L187)*

___

## `Optional` enableEthereumRPCRateLimiting

• **enableEthereumRPCRateLimiting**? : *undefined | false | true*

*Defined in [types.ts:135](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L135)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:103](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L103)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [types.ts:127](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L127)*

___

## `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : *undefined | number*

*Defined in [types.ts:140](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L140)*

___

## `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : *undefined | number*

*Defined in [types.ts:146](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L146)*

___

## `Optional` ethereumRPCURL

• **ethereumRPCURL**? : *undefined | string*

*Defined in [types.ts:100](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L100)*

___

## `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : *undefined | number*

*Defined in [types.ts:167](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L167)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [types.ts:106](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L106)*

___

## `Optional` verbosity

• **verbosity**? : *[Verbosity](#enumeration-verbosity)*

*Defined in [types.ts:97](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L97)*

___

## `Optional` web3Provider

• **web3Provider**? : *SupportedProvider*

*Defined in [types.ts:190](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L190)*


<hr />

# Interface: ContractAddresses

### Hierarchy

* **ContractAddresses**


### Properties

## `Optional` coordinator

• **coordinator**? : *undefined | string*

*Defined in [types.ts:199](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L199)*

___

## `Optional` coordinatorRegistry

• **coordinatorRegistry**? : *undefined | string*

*Defined in [types.ts:200](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L200)*

___

##  devUtils

• **devUtils**: *string*

*Defined in [types.ts:195](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L195)*

___

##  erc1155Proxy

• **erc1155Proxy**: *string*

*Defined in [types.ts:198](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L198)*

___

##  erc20Proxy

• **erc20Proxy**: *string*

*Defined in [types.ts:196](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L196)*

___

##  erc721Proxy

• **erc721Proxy**: *string*

*Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L197)*

___

##  exchange

• **exchange**: *string*

*Defined in [types.ts:194](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L194)*

___

## `Optional` weth9

• **weth9**? : *undefined | string*

*Defined in [types.ts:201](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L201)*

___

## `Optional` zrxToken

• **zrxToken**? : *undefined | string*

*Defined in [types.ts:202](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L202)*


<hr />

# Interface: ContractEvent

### Hierarchy

* **ContractEvent**


### Properties

##  address

• **address**: *string*

*Defined in [types.ts:516](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L516)*

___

##  blockHash

• **blockHash**: *string*

*Defined in [types.ts:511](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L511)*

___

##  isRemoved

• **isRemoved**: *boolean*

*Defined in [types.ts:515](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L515)*

___

##  kind

• **kind**: *[ContractEventKind](#enumeration-contracteventkind)*

*Defined in [types.ts:517](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L517)*

___

##  logIndex

• **logIndex**: *number*

*Defined in [types.ts:514](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L514)*

___

##  parameters

• **parameters**: *ContractEventParameters*

*Defined in [types.ts:518](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L518)*

___

##  txHash

• **txHash**: *string*

*Defined in [types.ts:512](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L512)*

___

##  txIndex

• **txIndex**: *number*

*Defined in [types.ts:513](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L513)*


<hr />

# Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [types.ts:380](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L380)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:379](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L379)*

___

##  owner

• **owner**: *string*

*Defined in [types.ts:378](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L378)*


<hr />

# Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:362](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L362)*

___

##  ids

• **ids**: *BigNumber[]*

*Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L364)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:361](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L361)*

___

##  to

• **to**: *string*

*Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L363)*

___

##  values

• **values**: *BigNumber[]*

*Defined in [types.ts:365](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L365)*


<hr />

# Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:345](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L345)*

___

##  id

• **id**: *BigNumber*

*Defined in [types.ts:347](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L347)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:344](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L344)*

___

##  to

• **to**: *string*

*Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L346)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:348](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L348)*


<hr />

# Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:299](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L299)*

___

##  spender

• **spender**: *string*

*Defined in [types.ts:300](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L300)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:301](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L301)*


<hr />

# Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L286)*

___

##  to

• **to**: *string*

*Defined in [types.ts:287](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L287)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:288](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L288)*


<hr />

# Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**


### Properties

##  approved

• **approved**: *string*

*Defined in [types.ts:326](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L326)*

___

##  owner

• **owner**: *string*

*Defined in [types.ts:325](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L325)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [types.ts:327](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L327)*


<hr />

# Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**


### Properties

##  approved

• **approved**: *boolean*

*Defined in [types.ts:340](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L340)*

___

##  operator

• **operator**: *string*

*Defined in [types.ts:339](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L339)*

___

##  owner

• **owner**: *string*

*Defined in [types.ts:338](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L338)*


<hr />

# Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**


### Properties

##  from

• **from**: *string*

*Defined in [types.ts:312](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L312)*

___

##  to

• **to**: *string*

*Defined in [types.ts:313](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L313)*

___

##  tokenId

• **tokenId**: *BigNumber*

*Defined in [types.ts:314](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L314)*


<hr />

# Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:421](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L421)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:419](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L419)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:423](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L423)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:422](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L422)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:420](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L420)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L424)*


<hr />

# Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**


### Properties

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:428](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L428)*

___

##  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [types.ts:430](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L430)*

___

##  orderSenderAddress

• **orderSenderAddress**: *string*

*Defined in [types.ts:429](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L429)*


<hr />

# Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**


### Properties

##  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:387](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L387)*

___

##  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:384](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L384)*

___

##  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:394](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L394)*

___

##  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [types.ts:388](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L388)*

___

##  makerFeeAssetData

• **makerFeeAssetData**: *string*

*Defined in [types.ts:396](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L396)*

___

##  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [types.ts:390](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L390)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:393](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L393)*

___

##  protocolFeePaid

• **protocolFeePaid**: *BigNumber*

*Defined in [types.ts:392](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L392)*

___

##  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:386](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L386)*

___

##  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:385](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L385)*

___

##  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:395](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L395)*

___

##  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [types.ts:389](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L389)*

___

##  takerFeeAssetData

• **takerFeeAssetData**: *string*

*Defined in [types.ts:397](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L397)*

___

##  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [types.ts:391](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L391)*


<hr />

# Interface: GetOrdersResponse

### Hierarchy

* **GetOrdersResponse**


### Properties

##  ordersInfos

• **ordersInfos**: *[OrderInfo](#interface-orderinfo)[]*

*Defined in [types.ts:19](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L19)*

___

##  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:17](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L17)*

___

##  snapshotTimestamp

• **snapshotTimestamp**: *number*

*Defined in [types.ts:18](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L18)*


<hr />

# Interface: JsonSchema

An interface for JSON schema types, which are used for custom order filters.

### Hierarchy

* **JsonSchema**


### Properties

## `Optional` $ref

• **$ref**? : *undefined | string*

*Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L41)*

___

## `Optional` $schema

• **$schema**? : *undefined | string*

*Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L40)*

___

## `Optional` additionalItems

• **additionalItems**? : *boolean | [JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:52](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L52)*

___

## `Optional` additionalProperties

• **additionalProperties**? : *boolean | [JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:60](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L60)*

___

## `Optional` allOf

• **allOf**? : *[JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:82](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L82)*

___

## `Optional` anyOf

• **anyOf**? : *[JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:83](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L83)*

___

## `Optional` const

• **const**? : *any*

*Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L79)*

___

## `Optional` definitions

• **definitions**? : *undefined | object*

*Defined in [types.ts:61](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L61)*

___

## `Optional` dependencies

• **dependencies**? : *undefined | object*

*Defined in [types.ts:70](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L70)*

___

## `Optional` description

• **description**? : *undefined | string*

*Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L43)*

___

## `Optional` enum

• **enum**? : *any[]*

*Defined in [types.ts:73](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L73)*

___

## `Optional` exclusiveMaximum

• **exclusiveMaximum**? : *undefined | false | true*

*Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L46)*

___

## `Optional` exclusiveMinimum

• **exclusiveMinimum**? : *undefined | false | true*

*Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L48)*

___

## `Optional` format

• **format**? : *undefined | string*

*Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L81)*

___

## `Optional` id

• **id**? : *undefined | string*

*Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L39)*

___

## `Optional` items

• **items**? : *[JsonSchema](#interface-jsonschema) | [JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:53](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L53)*

___

## `Optional` maxItems

• **maxItems**? : *undefined | number*

*Defined in [types.ts:54](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L54)*

___

## `Optional` maxLength

• **maxLength**? : *undefined | number*

*Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L49)*

___

## `Optional` maxProperties

• **maxProperties**? : *undefined | number*

*Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L57)*

___

## `Optional` maximum

• **maximum**? : *undefined | number*

*Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L45)*

___

## `Optional` minItems

• **minItems**? : *undefined | number*

*Defined in [types.ts:55](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L55)*

___

## `Optional` minLength

• **minLength**? : *undefined | number*

*Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L50)*

___

## `Optional` minProperties

• **minProperties**? : *undefined | number*

*Defined in [types.ts:58](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L58)*

___

## `Optional` minimum

• **minimum**? : *undefined | number*

*Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L47)*

___

## `Optional` multipleOf

• **multipleOf**? : *undefined | number*

*Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L44)*

___

## `Optional` not

• **not**? : *[JsonSchema](#interface-jsonschema)*

*Defined in [types.ts:85](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L85)*

___

## `Optional` oneOf

• **oneOf**? : *[JsonSchema](#interface-jsonschema)[]*

*Defined in [types.ts:84](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L84)*

___

## `Optional` pattern

• **pattern**? : *string | RegExp*

*Defined in [types.ts:51](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L51)*

___

## `Optional` patternProperties

• **patternProperties**? : *undefined | object*

*Defined in [types.ts:67](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L67)*

___

## `Optional` properties

• **properties**? : *undefined | object*

*Defined in [types.ts:64](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L64)*

___

## `Optional` required

• **required**? : *string[]*

*Defined in [types.ts:59](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L59)*

___

## `Optional` title

• **title**? : *undefined | string*

*Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L42)*

___

## `Optional` type

• **type**? : *string | string[]*

*Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L80)*

___

## `Optional` uniqueItems

• **uniqueItems**? : *undefined | false | true*

*Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L56)*


<hr />

# Interface: LatestBlock

### Hierarchy

* **LatestBlock**


### Properties

##  hash

• **hash**: *string*

*Defined in [types.ts:643](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L643)*

___

##  number

• **number**: *number*

*Defined in [types.ts:642](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L642)*


<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  contractEvents

• **contractEvents**: *[ContractEvent](#interface-contractevent)[]*

*Defined in [types.ts:569](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L569)*

___

##  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:567](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L567)*

___

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:568](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L568)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:565](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L565)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:566](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L566)*

___

##  timestampMs

• **timestampMs**: *number*

*Defined in [types.ts:564](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L564)*


<hr />

# Interface: OrderInfo

### Hierarchy

* **OrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:32](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L32)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L30)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:31](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L31)*


<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [types.ts:619](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L619)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [types.ts:617](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L617)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:618](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L618)*

___

##  status

• **status**: *[RejectedOrderStatus](#interface-rejectedorderstatus)*

*Defined in [types.ts:620](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L620)*


<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [types.ts:637](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L637)*

___

##  message

• **message**: *string*

*Defined in [types.ts:638](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L638)*


<hr />

# Interface: Stats

### Hierarchy

* **Stats**


### Properties

##  ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: *number*

*Defined in [types.ts:680](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L680)*

___

##  ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: *number*

*Defined in [types.ts:679](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L679)*

___

##  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:671](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L671)*

___

##  latestBlock

• **latestBlock**: *[LatestBlock](#interface-latestblock)*

*Defined in [types.ts:672](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L672)*

___

##  maxExpirationTime

• **maxExpirationTime**: *BigNumber*

*Defined in [types.ts:677](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L677)*

___

##  numOrders

• **numOrders**: *number*

*Defined in [types.ts:674](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L674)*

___

##  numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: *number*

*Defined in [types.ts:675](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L675)*

___

##  numPeers

• **numPeers**: *number*

*Defined in [types.ts:673](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L673)*

___

##  numPinnedOrders

• **numPinnedOrders**: *number*

*Defined in [types.ts:676](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L676)*

___

##  peerID

• **peerID**: *string*

*Defined in [types.ts:670](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L670)*

___

##  pubSubTopic

• **pubSubTopic**: *string*

*Defined in [types.ts:667](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L667)*

___

##  rendezvous

• **rendezvous**: *string*

*Defined in [types.ts:668](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L668)*

___

##  secondaryRendezvous

• **secondaryRendezvous**: *string[]*

*Defined in [types.ts:669](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L669)*

___

##  startOfCurrentUTCDay

• **startOfCurrentUTCDay**: *Date*

*Defined in [types.ts:678](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L678)*

___

##  version

• **version**: *string*

*Defined in [types.ts:666](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L666)*


<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#interface-acceptedorderinfo)[]*

*Defined in [types.ts:598](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L598)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#interface-rejectedorderinfo)[]*

*Defined in [types.ts:599](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L599)*


<hr />

# Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:452](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L452)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:453](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L453)*


<hr />

# Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**


### Properties

##  owner

• **owner**: *string*

*Defined in [types.ts:441](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L441)*

___

##  value

• **value**: *BigNumber*

*Defined in [types.ts:442](https://github.com/0xProject/0x-mesh/blob/d606e1cf/packages/browser-lite/src/types.ts#L442)*


<hr />

