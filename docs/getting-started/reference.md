# Browser API documentation

## Class: Mesh

The main class for this package. Has methods for receiving order events and sending orders through the 0x Mesh network.

#### Hierarchy

* **Mesh**

#### Constructors

### constructer

+ **new Mesh**\(`config`: [Config](reference.md#interface-config)\): [_Mesh_](reference.md#class-mesh)

_Defined in_ [_index.ts:538_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L538)

Instantiates a new Mesh instance.

**Parameters:**

| Name | Type | Description |
| :--- | :--- | :--- |
| `config` | [Config](reference.md#interface-config) | Configuration options for Mesh |

**Returns:** [_Mesh_](reference.md#class-mesh)

An instance of Mesh

#### Methods

### addOrdersAsync

▸ **addOrdersAsync**\(`orders`: SignedOrder\[\], `pinned`: boolean\): _Promise‹_[_ValidationResults_](reference.md#interface-validationresults)_›_

_Defined in_ [_index.ts:612_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L612)

Validates and adds the given orders to Mesh. If an order is successfully added, Mesh will share it with any peers in the network and start watching it for changes \(e.g. filled, canceled, expired\). The returned promise will only be rejected if there was an error validating or adding the order; it will not be rejected for any invalid orders \(check results.rejected instead\).

**Parameters:**

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `orders` | SignedOrder\[\] | - | An array of orders to add. |
| `pinned` | boolean | true | Whether or not the orders should be pinned. Pinned orders will not be affected by any DDoS prevention or incentive mechanisms and will always stay in storage until they are no longer fillable. |

**Returns:** _Promise‹_[_ValidationResults_](reference.md#interface-validationresults)_›_

Validation results for the given orders, indicating which orders were accepted and which were rejected.

### onError

▸ **onError**\(`handler`: function\): _void_

_Defined in_ [_index.ts:558_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L558)

Registers a handler which will be called in the event of a critical error. Note that the handler will not be called for non-critical errors. In order to ensure no errors are missed, this should be called before startAsync.

**Parameters:**

▪ **handler**: _function_

The handler to be called.

▸ \(`err`: Error\): _void_

**Parameters:**

| Name | Type |
| :--- | :--- |
| `err` | Error |

**Returns:** _void_

### onOrderEvents

▸ **onOrderEvents**\(`handler`: function\): _void_

_Defined in_ [_index.ts:573_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L573)

Registers a handler which will be called for any incoming order events. Order events are fired whenver an order is added, canceled, expired, or filled. In order to ensure no events are missed, this should be called before startAsync.

**Parameters:**

▪ **handler**: _function_

The handler to be called.

▸ \(`events`: [OrderEvent](reference.md#interface-orderevent)\[\]\): _void_

**Parameters:**

| Name | Type |
| :--- | :--- |
| `events` | [OrderEvent](reference.md#interface-orderevent)\[\] |

**Returns:** _void_

### startAsync

▸ **startAsync**\(\): _Promise‹void›_

_Defined in_ [_index.ts:584_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L584)

Starts the Mesh node in the background. Mesh will automatically find peers in the network and begin receiving orders from them.

**Returns:** _Promise‹void›_

## Enumeration: OrderEventEndState

#### Enumeration members

### Added

• **Added**: = "ADDED"

_Defined in_ [_index.ts:404_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L404)

### Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in_ [_index.ts:407_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L407)

### Expired

• **Expired**: = "EXPIRED"

_Defined in_ [_index.ts:408_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L408)

### FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY\_INCREASED"

_Defined in_ [_index.ts:410_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L410)

### Filled

• **Filled**: = "FILLED"

_Defined in_ [_index.ts:405_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L405)

### FullyFilled

• **FullyFilled**: = "FULLY\_FILLED"

_Defined in_ [_index.ts:406_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L406)

### Invalid

• **Invalid**: = "INVALID"

_Defined in_ [_index.ts:403_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L403)

### StoppedWatching

• **StoppedWatching**: = "STOPPED\_WATCHING"

_Defined in_ [_index.ts:411_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L411)

### Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in_ [_index.ts:409_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L409)

## Enumeration: RejectedOrderKind

A set of categories for rejected orders.

#### Enumeration members

### CoordinatorError

• **CoordinatorError**: = "COORDINATOR\_ERROR"

_Defined in_ [_index.ts:492_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L492)

### MeshError

• **MeshError**: = "MESH\_ERROR"

_Defined in_ [_index.ts:490_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L490)

### MeshValidation

• **MeshValidation**: = "MESH\_VALIDATION"

_Defined in_ [_index.ts:491_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L491)

### ZeroExValidation

• **ZeroExValidation**: = "ZEROEX\_VALIDATION"

_Defined in_ [_index.ts:489_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L489)

## Enumeration: Verbosity

#### Enumeration members

### Debug

• **Debug**: = 5

_Defined in_ [_index.ts:118_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L118)

### Error

• **Error**: = 2

_Defined in_ [_index.ts:115_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L115)

### Fatal

• **Fatal**: = 1

_Defined in_ [_index.ts:114_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L114)

### Info

• **Info**: = 4

_Defined in_ [_index.ts:117_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L117)

### Panic

• **Panic**: = 0

_Defined in_ [_index.ts:113_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L113)

### Trace

• **Trace**: = 6

_Defined in_ [_index.ts:119_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L119)

### Warn

• **Warn**: = 3

_Defined in_ [_index.ts:116_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L116)

## Interface: AcceptedOrderInfo

Info for any orders that were accepted.

#### Hierarchy

* **AcceptedOrderInfo**

#### Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:470_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L470)

### isNew

• **isNew**: _boolean_

_Defined in_ [_index.ts:471_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L471)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:468_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L468)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:469_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L469)

## Interface: Config

A set of configuration options for Mesh.

#### Hierarchy

* **Config**

#### Properties

### `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : _undefined \| number_

_Defined in_ [_index.ts:56_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L56)

### `Optional` bootstrapList

• **bootstrapList**? : _string\[\]_

_Defined in_ [_index.ts:49_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L49)

### `Optional` customContractAddresses

• **customContractAddresses**? : [_ContractAddresses_](reference.md#class-contractaddresses)

_Defined in_ [_index.ts:92_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L92)

### ethereumChainID

• **ethereumChainID**: _number_

_Defined in_ [_index.ts:41_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L41)

### `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : _undefined \| number_

_Defined in_ [_index.ts:65_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L65)

### `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : _undefined \| number_

_Defined in_ [_index.ts:70_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L70)

### `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : _undefined \| number_

_Defined in_ [_index.ts:76_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L76)

### ethereumRPCURL

• **ethereumRPCURL**: _string_

_Defined in_ [_index.ts:38_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L38)

### `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : _undefined \| number_

_Defined in_ [_index.ts:97_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L97)

### `Optional` useBootstrapList

• **useBootstrapList**? : _undefined \| false \| true_

_Defined in_ [_index.ts:44_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L44)

### `Optional` verbosity

• **verbosity**? : [_Verbosity_](reference.md#enumeration-verbosity)

_Defined in_ [_index.ts:35_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L35)

## Interface: ContractAddresses

#### Hierarchy

* **ContractAddresses**

#### Properties

### `Optional` coordinator

• **coordinator**? : _undefined \| string_

_Defined in_ [_index.ts:106_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L106)

### `Optional` coordinatorRegistry

• **coordinatorRegistry**? : _undefined \| string_

_Defined in_ [_index.ts:107_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L107)

### devUtils

• **devUtils**: _string_

_Defined in_ [_index.ts:102_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L102)

### erc1155Proxy

• **erc1155Proxy**: _string_

_Defined in_ [_index.ts:105_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L105)

### erc20Proxy

• **erc20Proxy**: _string_

_Defined in_ [_index.ts:103_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L103)

### erc721Proxy

• **erc721Proxy**: _string_

_Defined in_ [_index.ts:104_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L104)

### exchange

• **exchange**: _string_

_Defined in_ [_index.ts:101_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L101)

### `Optional` weth9

• **weth9**? : _undefined \| string_

_Defined in_ [_index.ts:108_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L108)

### `Optional` zrxToken

• **zrxToken**? : _undefined \| string_

_Defined in_ [_index.ts:109_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L109)

## Interface: ContractEvent

#### Hierarchy

* **ContractEvent**

#### Properties

### address

• **address**: _string_

_Defined in_ [_index.ts:385_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L385)

### blockHash

• **blockHash**: _string_

_Defined in_ [_index.ts:380_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L380)

### isRemoved

• **isRemoved**: _string_

_Defined in_ [_index.ts:384_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L384)

### kind

• **kind**: _ContractEventKind_

_Defined in_ [_index.ts:386_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L386)

### logIndex

• **logIndex**: _number_

_Defined in_ [_index.ts:383_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L383)

### parameters

• **parameters**: _ContractEventParameters_

_Defined in_ [_index.ts:387_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L387)

### txHash

• **txHash**: _string_

_Defined in_ [_index.ts:381_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L381)

### txIndex

• **txIndex**: _number_

_Defined in_ [_index.ts:382_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L382)

## Interface: ERC1155ApprovalForAllEvent

#### Hierarchy

* **ERC1155ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:261_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L261)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:260_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L260)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:259_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L259)

## Interface: ERC1155TransferBatchEvent

#### Hierarchy

* **ERC1155TransferBatchEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:244_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L244)

### ids

• **ids**: _BigNumber\[\]_

_Defined in_ [_index.ts:246_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L246)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:243_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L243)

### to

• **to**: _string_

_Defined in_ [_index.ts:245_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L245)

### values

• **values**: _BigNumber\[\]_

_Defined in_ [_index.ts:247_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L247)

## Interface: ERC1155TransferSingleEvent

#### Hierarchy

* **ERC1155TransferSingleEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:228_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L228)

### id

• **id**: _BigNumber_

_Defined in_ [_index.ts:230_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L230)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:227_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L227)

### to

• **to**: _string_

_Defined in_ [_index.ts:229_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L229)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:231_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L231)

## Interface: ERC20ApprovalEvent

#### Hierarchy

* **ERC20ApprovalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:185_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L185)

### spender

• **spender**: _string_

_Defined in_ [_index.ts:186_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L186)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:187_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L187)

## Interface: ERC20TransferEvent

#### Hierarchy

* **ERC20TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:173_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L173)

### to

• **to**: _string_

_Defined in_ [_index.ts:174_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L174)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:175_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L175)

## Interface: ERC721ApprovalEvent

#### Hierarchy

* **ERC721ApprovalEvent**

#### Properties

### approved

• **approved**: _string_

_Defined in_ [_index.ts:210_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L210)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:209_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L209)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:211_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L211)

## Interface: ERC721ApprovalForAllEvent

#### Hierarchy

* **ERC721ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:223_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L223)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:222_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L222)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:221_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L221)

## Interface: ERC721TransferEvent

#### Hierarchy

* **ERC721TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:197_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L197)

### to

• **to**: _string_

_Defined in_ [_index.ts:198_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L198)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:199_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L199)

## Interface: ExchangeCancelEvent

#### Hierarchy

* **ExchangeCancelEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:295_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L295)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:293_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L293)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:297_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L297)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:296_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L296)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:294_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L294)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:298_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L298)

## Interface: ExchangeCancelUpToEvent

#### Hierarchy

* **ExchangeCancelUpToEvent**

#### Properties

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:302_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L302)

### orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in_ [_index.ts:304_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L304)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:303_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L303)

## Interface: ExchangeFillEvent

#### Hierarchy

* **ExchangeFillEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:268_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L268)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:265_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L265)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:274_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L274)

### makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:269_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L269)

### makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:271_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L271)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:273_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L273)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:267_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L267)

### takerAddress

• **takerAddress**: _string_

_Defined in_ [_index.ts:266_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L266)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:275_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L275)

### takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:270_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L270)

### takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:272_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L272)

## Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired, or filled.

#### Hierarchy

* **OrderEvent**

#### Properties

### contractEvents

• **contractEvents**: [_ContractEvent_](reference.md#class-contractevent)_\[\]_

_Defined in_ [_index.ts:431_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L431)

### endState

• **endState**: [_OrderEventEndState_](reference.md#enumeration-ordereventendstate)

_Defined in_ [_index.ts:429_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L429)

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:430_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L430)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:427_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L427)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:428_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L428)

## Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were rejected.

#### Hierarchy

* **RejectedOrderInfo**

#### Properties

### kind

• **kind**: [_RejectedOrderKind_](reference.md#enumeration-rejectedorderkind)

_Defined in_ [_index.ts:481_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L481)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:479_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L479)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:480_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L480)

### status

• **status**: [_RejectedOrderStatus_](reference.md#class-rejectedorderstatus)

_Defined in_ [_index.ts:482_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L482)

## Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

#### Hierarchy

* **RejectedOrderStatus**

#### Properties

### code

• **code**: _string_

_Defined in_ [_index.ts:499_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L499)

### message

• **message**: _string_

_Defined in_ [_index.ts:500_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L500)

## Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

#### Hierarchy

* **ValidationResults**

#### Properties

### accepted

• **accepted**: [_AcceptedOrderInfo_](reference.md#class-acceptedorderinfo)_\[\]_

_Defined in_ [_index.ts:460_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L460)

### rejected

• **rejected**: [_RejectedOrderInfo_](reference.md#class-rejectedorderinfo)_\[\]_

_Defined in_ [_index.ts:461_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L461)

## Interface: WethDepositEvent

#### Hierarchy

* **WethDepositEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:324_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L324)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:325_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L325)

## Interface: WethWithdrawalEvent

#### Hierarchy

* **WethWithdrawalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:314_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L314)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:315_](https://github.com/0xProject/0x-mesh/blob/0240f5d/browser/ts/index.ts#L315)

