# Browser API documentation

## Class: Mesh

The main class for this package. Has methods for receiving order events and sending orders through the 0x Mesh network.

#### Hierarchy

* **Mesh**

#### Constructors

### constructer

+ **new Mesh**\(`config`: [Config](reference.md#interface-config)\): [_Mesh_](reference.md#class-mesh)

_Defined in_ [_index.ts:536_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L536)

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

_Defined in_ [_index.ts:610_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L610)

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

_Defined in_ [_index.ts:556_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L556)

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

_Defined in_ [_index.ts:571_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L571)

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

_Defined in_ [_index.ts:582_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L582)

Starts the Mesh node in the background. Mesh will automatically find peers in the network and begin receiving orders from them.

**Returns:** _Promise‹void›_

## Enumeration: OrderEventEndState

#### Enumeration members

### Added

• **Added**: = "ADDED"

_Defined in_ [_index.ts:402_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L402)

### Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in_ [_index.ts:405_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L405)

### Expired

• **Expired**: = "EXPIRED"

_Defined in_ [_index.ts:406_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L406)

### FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY\_INCREASED"

_Defined in_ [_index.ts:408_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L408)

### Filled

• **Filled**: = "FILLED"

_Defined in_ [_index.ts:403_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L403)

### FullyFilled

• **FullyFilled**: = "FULLY\_FILLED"

_Defined in_ [_index.ts:404_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L404)

### Invalid

• **Invalid**: = "INVALID"

_Defined in_ [_index.ts:401_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L401)

### StoppedWatching

• **StoppedWatching**: = "STOPPED\_WATCHING"

_Defined in_ [_index.ts:409_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L409)

### Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in_ [_index.ts:407_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L407)

## Enumeration: RejectedOrderKind

A set of categories for rejected orders.

#### Enumeration members

### CoordinatorError

• **CoordinatorError**: = "COORDINATOR\_ERROR"

_Defined in_ [_index.ts:490_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L490)

### MeshError

• **MeshError**: = "MESH\_ERROR"

_Defined in_ [_index.ts:488_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L488)

### MeshValidation

• **MeshValidation**: = "MESH\_VALIDATION"

_Defined in_ [_index.ts:489_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L489)

### ZeroExValidation

• **ZeroExValidation**: = "ZEROEX\_VALIDATION"

_Defined in_ [_index.ts:487_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L487)

## Enumeration: Verbosity

#### Enumeration members

### Debug

• **Debug**: = 5

_Defined in_ [_index.ts:116_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L116)

### Error

• **Error**: = 2

_Defined in_ [_index.ts:113_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L113)

### Fatal

• **Fatal**: = 1

_Defined in_ [_index.ts:112_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L112)

### Info

• **Info**: = 4

_Defined in_ [_index.ts:115_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L115)

### Panic

• **Panic**: = 0

_Defined in_ [_index.ts:111_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L111)

### Trace

• **Trace**: = 6

_Defined in_ [_index.ts:117_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L117)

### Warn

• **Warn**: = 3

_Defined in_ [_index.ts:114_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L114)

## Interface: AcceptedOrderInfo

Info for any orders that were accepted.

#### Hierarchy

* **AcceptedOrderInfo**

#### Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:468_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L468)

### isNew

• **isNew**: _boolean_

_Defined in_ [_index.ts:469_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L469)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:466_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L466)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:467_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L467)

## Interface: Config

A set of configuration options for Mesh.

#### Hierarchy

* **Config**

#### Properties

### `Optional` EthereumRPCMaxRequestsPer24HrUTC

• **EthereumRPCMaxRequestsPer24HrUTC**? : _undefined \| number_

_Defined in_ [_index.ts:69_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L69)

### `Optional` EthereumRPCMaxRequestsPerSecond

• **EthereumRPCMaxRequestsPerSecond**? : _undefined \| number_

_Defined in_ [_index.ts:74_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L74)

### `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : _undefined \| number_

_Defined in_ [_index.ts:56_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L56)

### `Optional` bootstrapList

• **bootstrapList**? : _string\[\]_

_Defined in_ [_index.ts:49_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L49)

### `Optional` customContractAddresses

• **customContractAddresses**? : [_ContractAddresses_](reference.md#class-contractaddresses)

_Defined in_ [_index.ts:90_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L90)

### ethereumChainID

• **ethereumChainID**: _number_

_Defined in_ [_index.ts:41_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L41)

### `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : _undefined \| number_

_Defined in_ [_index.ts:65_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L65)

### ethereumRPCURL

• **ethereumRPCURL**: _string_

_Defined in_ [_index.ts:38_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L38)

### `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : _undefined \| number_

_Defined in_ [_index.ts:95_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L95)

### `Optional` useBootstrapList

• **useBootstrapList**? : _undefined \| false \| true_

_Defined in_ [_index.ts:44_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L44)

### `Optional` verbosity

• **verbosity**? : [_Verbosity_](reference.md#enumeration-verbosity)

_Defined in_ [_index.ts:35_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L35)

## Interface: ContractAddresses

#### Hierarchy

* **ContractAddresses**

#### Properties

### `Optional` coordinator

• **coordinator**? : _undefined \| string_

_Defined in_ [_index.ts:104_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L104)

### `Optional` coordinatorRegistry

• **coordinatorRegistry**? : _undefined \| string_

_Defined in_ [_index.ts:105_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L105)

### devUtils

• **devUtils**: _string_

_Defined in_ [_index.ts:100_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L100)

### erc1155Proxy

• **erc1155Proxy**: _string_

_Defined in_ [_index.ts:103_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L103)

### erc20Proxy

• **erc20Proxy**: _string_

_Defined in_ [_index.ts:101_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L101)

### erc721Proxy

• **erc721Proxy**: _string_

_Defined in_ [_index.ts:102_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L102)

### exchange

• **exchange**: _string_

_Defined in_ [_index.ts:99_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L99)

### `Optional` weth9

• **weth9**? : _undefined \| string_

_Defined in_ [_index.ts:106_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L106)

### `Optional` zrxToken

• **zrxToken**? : _undefined \| string_

_Defined in_ [_index.ts:107_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L107)

## Interface: ContractEvent

#### Hierarchy

* **ContractEvent**

#### Properties

### address

• **address**: _string_

_Defined in_ [_index.ts:383_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L383)

### blockHash

• **blockHash**: _string_

_Defined in_ [_index.ts:378_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L378)

### isRemoved

• **isRemoved**: _string_

_Defined in_ [_index.ts:382_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L382)

### kind

• **kind**: _ContractEventKind_

_Defined in_ [_index.ts:384_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L384)

### logIndex

• **logIndex**: _number_

_Defined in_ [_index.ts:381_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L381)

### parameters

• **parameters**: _ContractEventParameters_

_Defined in_ [_index.ts:385_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L385)

### txHash

• **txHash**: _string_

_Defined in_ [_index.ts:379_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L379)

### txIndex

• **txIndex**: _number_

_Defined in_ [_index.ts:380_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L380)

## Interface: ERC1155ApprovalForAllEvent

#### Hierarchy

* **ERC1155ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:259_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L259)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:258_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L258)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:257_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L257)

## Interface: ERC1155TransferBatchEvent

#### Hierarchy

* **ERC1155TransferBatchEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:242_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L242)

### ids

• **ids**: _BigNumber\[\]_

_Defined in_ [_index.ts:244_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L244)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:241_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L241)

### to

• **to**: _string_

_Defined in_ [_index.ts:243_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L243)

### values

• **values**: _BigNumber\[\]_

_Defined in_ [_index.ts:245_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L245)

## Interface: ERC1155TransferSingleEvent

#### Hierarchy

* **ERC1155TransferSingleEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:226_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L226)

### id

• **id**: _BigNumber_

_Defined in_ [_index.ts:228_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L228)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:225_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L225)

### to

• **to**: _string_

_Defined in_ [_index.ts:227_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L227)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:229_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L229)

## Interface: ERC20ApprovalEvent

#### Hierarchy

* **ERC20ApprovalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:183_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L183)

### spender

• **spender**: _string_

_Defined in_ [_index.ts:184_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L184)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:185_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L185)

## Interface: ERC20TransferEvent

#### Hierarchy

* **ERC20TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:171_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L171)

### to

• **to**: _string_

_Defined in_ [_index.ts:172_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L172)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:173_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L173)

## Interface: ERC721ApprovalEvent

#### Hierarchy

* **ERC721ApprovalEvent**

#### Properties

### approved

• **approved**: _string_

_Defined in_ [_index.ts:208_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L208)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:207_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L207)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:209_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L209)

## Interface: ERC721ApprovalForAllEvent

#### Hierarchy

* **ERC721ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:221_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L221)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:220_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L220)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:219_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L219)

## Interface: ERC721TransferEvent

#### Hierarchy

* **ERC721TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:195_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L195)

### to

• **to**: _string_

_Defined in_ [_index.ts:196_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L196)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:197_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L197)

## Interface: ExchangeCancelEvent

#### Hierarchy

* **ExchangeCancelEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:293_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L293)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:291_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L291)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:295_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L295)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:294_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L294)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:292_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L292)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:296_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L296)

## Interface: ExchangeCancelUpToEvent

#### Hierarchy

* **ExchangeCancelUpToEvent**

#### Properties

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:300_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L300)

### orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in_ [_index.ts:302_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L302)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:301_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L301)

## Interface: ExchangeFillEvent

#### Hierarchy

* **ExchangeFillEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:266_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L266)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:263_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L263)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:272_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L272)

### makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:267_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L267)

### makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:269_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L269)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:271_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L271)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:265_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L265)

### takerAddress

• **takerAddress**: _string_

_Defined in_ [_index.ts:264_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L264)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:273_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L273)

### takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:268_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L268)

### takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:270_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L270)

## Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired, or filled.

#### Hierarchy

* **OrderEvent**

#### Properties

### contractEvents

• **contractEvents**: [_ContractEvent_](reference.md#class-contractevent)_\[\]_

_Defined in_ [_index.ts:429_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L429)

### endState

• **endState**: [_OrderEventEndState_](reference.md#enumeration-ordereventendstate)

_Defined in_ [_index.ts:427_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L427)

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:428_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L428)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:425_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L425)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:426_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L426)

## Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were rejected.

#### Hierarchy

* **RejectedOrderInfo**

#### Properties

### kind

• **kind**: [_RejectedOrderKind_](reference.md#enumeration-rejectedorderkind)

_Defined in_ [_index.ts:479_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L479)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:477_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L477)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:478_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L478)

### status

• **status**: [_RejectedOrderStatus_](reference.md#class-rejectedorderstatus)

_Defined in_ [_index.ts:480_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L480)

## Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

#### Hierarchy

* **RejectedOrderStatus**

#### Properties

### code

• **code**: _string_

_Defined in_ [_index.ts:497_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L497)

### message

• **message**: _string_

_Defined in_ [_index.ts:498_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L498)

## Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

#### Hierarchy

* **ValidationResults**

#### Properties

### accepted

• **accepted**: [_AcceptedOrderInfo_](reference.md#class-acceptedorderinfo)_\[\]_

_Defined in_ [_index.ts:458_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L458)

### rejected

• **rejected**: [_RejectedOrderInfo_](reference.md#class-rejectedorderinfo)_\[\]_

_Defined in_ [_index.ts:459_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L459)

## Interface: WethDepositEvent

#### Hierarchy

* **WethDepositEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:322_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L322)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:323_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L323)

## Interface: WethWithdrawalEvent

#### Hierarchy

* **WethWithdrawalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:312_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L312)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:313_](https://github.com/0xProject/0x-mesh/blob/7d5a102/browser/ts/index.ts#L313)

