# Browser API documentation

## Class: Mesh

The main class for this package. Has methods for receiving order events and sending orders through the 0x Mesh network.

#### Hierarchy

* **Mesh**

#### Constructors

### constructer

+ **new Mesh**\(`config`: [Config](reference.md#interface-config)\): [_Mesh_](reference.md#class-mesh)

_Defined in_ [_index.ts:562_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L562)

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

_Defined in_ [_index.ts:636_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L636)

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

_Defined in_ [_index.ts:582_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L582)

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

_Defined in_ [_index.ts:597_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L597)

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

_Defined in_ [_index.ts:608_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L608)

Starts the Mesh node in the background. Mesh will automatically find peers in the network and begin receiving orders from them.

**Returns:** _Promise‹void›_

## Enumeration: OrderEventEndState

#### Enumeration members

### Added

• **Added**: = "ADDED"

_Defined in_ [_index.ts:427_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L427)

### Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in_ [_index.ts:430_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L430)

### Expired

• **Expired**: = "EXPIRED"

_Defined in_ [_index.ts:431_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L431)

### FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY\_INCREASED"

_Defined in_ [_index.ts:434_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L434)

### Filled

• **Filled**: = "FILLED"

_Defined in_ [_index.ts:428_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L428)

### FullyFilled

• **FullyFilled**: = "FULLY\_FILLED"

_Defined in_ [_index.ts:429_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L429)

### Invalid

• **Invalid**: = "INVALID"

_Defined in_ [_index.ts:426_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L426)

### StoppedWatching

• **StoppedWatching**: = "STOPPED\_WATCHING"

_Defined in_ [_index.ts:435_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L435)

### Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in_ [_index.ts:432_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L432)

### Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in_ [_index.ts:433_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L433)

## Enumeration: RejectedOrderKind

A set of categories for rejected orders.

#### Enumeration members

### CoordinatorError

• **CoordinatorError**: = "COORDINATOR\_ERROR"

_Defined in_ [_index.ts:516_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L516)

### MeshError

• **MeshError**: = "MESH\_ERROR"

_Defined in_ [_index.ts:514_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L514)

### MeshValidation

• **MeshValidation**: = "MESH\_VALIDATION"

_Defined in_ [_index.ts:515_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L515)

### ZeroExValidation

• **ZeroExValidation**: = "ZEROEX\_VALIDATION"

_Defined in_ [_index.ts:513_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L513)

## Enumeration: Verbosity

#### Enumeration members

### Debug

• **Debug**: = 5

_Defined in_ [_index.ts:141_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L141)

### Error

• **Error**: = 2

_Defined in_ [_index.ts:138_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L138)

### Fatal

• **Fatal**: = 1

_Defined in_ [_index.ts:137_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L137)

### Info

• **Info**: = 4

_Defined in_ [_index.ts:140_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L140)

### Panic

• **Panic**: = 0

_Defined in_ [_index.ts:136_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L136)

### Trace

• **Trace**: = 6

_Defined in_ [_index.ts:142_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L142)

### Warn

• **Warn**: = 3

_Defined in_ [_index.ts:139_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L139)

## Interface: AcceptedOrderInfo

Info for any orders that were accepted.

#### Hierarchy

* **AcceptedOrderInfo**

#### Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:494_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L494)

### isNew

• **isNew**: _boolean_

_Defined in_ [_index.ts:495_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L495)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:492_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L492)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:493_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L493)

## Interface: Config

A set of configuration options for Mesh.

#### Hierarchy

* **Config**

#### Properties

### `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : _undefined \| number_

_Defined in_ [_index.ts:79_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L79)

### `Optional` bootstrapList

• **bootstrapList**? : _string\[\]_

_Defined in_ [_index.ts:72_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L72)

### `Optional` customContractAddresses

• **customContractAddresses**? : [_ContractAddresses_](reference.md#class-contractaddresses)

_Defined in_ [_index.ts:115_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L115)

### ethereumChainID

• **ethereumChainID**: _number_

_Defined in_ [_index.ts:64_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L64)

### `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : _undefined \| number_

_Defined in_ [_index.ts:88_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L88)

### `Optional` ethereumRPCMaxRequestsPer24HrUTC

• **ethereumRPCMaxRequestsPer24HrUTC**? : _undefined \| number_

_Defined in_ [_index.ts:93_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L93)

### `Optional` ethereumRPCMaxRequestsPerSecond

• **ethereumRPCMaxRequestsPerSecond**? : _undefined \| number_

_Defined in_ [_index.ts:99_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L99)

### ethereumRPCURL

• **ethereumRPCURL**: _string_

_Defined in_ [_index.ts:61_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L61)

### `Optional` maxOrdersInStorage

• **maxOrdersInStorage**? : _undefined \| number_

_Defined in_ [_index.ts:120_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L120)

### `Optional` useBootstrapList

• **useBootstrapList**? : _undefined \| false \| true_

_Defined in_ [_index.ts:67_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L67)

### `Optional` verbosity

• **verbosity**? : [_Verbosity_](reference.md#enumeration-verbosity)

_Defined in_ [_index.ts:58_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L58)

## Interface: ContractAddresses

#### Hierarchy

* **ContractAddresses**

#### Properties

### `Optional` coordinator

• **coordinator**? : _undefined \| string_

_Defined in_ [_index.ts:129_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L129)

### `Optional` coordinatorRegistry

• **coordinatorRegistry**? : _undefined \| string_

_Defined in_ [_index.ts:130_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L130)

### devUtils

• **devUtils**: _string_

_Defined in_ [_index.ts:125_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L125)

### erc1155Proxy

• **erc1155Proxy**: _string_

_Defined in_ [_index.ts:128_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L128)

### erc20Proxy

• **erc20Proxy**: _string_

_Defined in_ [_index.ts:126_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L126)

### erc721Proxy

• **erc721Proxy**: _string_

_Defined in_ [_index.ts:127_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L127)

### exchange

• **exchange**: _string_

_Defined in_ [_index.ts:124_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L124)

### `Optional` weth9

• **weth9**? : _undefined \| string_

_Defined in_ [_index.ts:131_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L131)

### `Optional` zrxToken

• **zrxToken**? : _undefined \| string_

_Defined in_ [_index.ts:132_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L132)

## Interface: ContractEvent

#### Hierarchy

* **ContractEvent**

#### Properties

### address

• **address**: _string_

_Defined in_ [_index.ts:408_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L408)

### blockHash

• **blockHash**: _string_

_Defined in_ [_index.ts:403_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L403)

### isRemoved

• **isRemoved**: _string_

_Defined in_ [_index.ts:407_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L407)

### kind

• **kind**: _ContractEventKind_

_Defined in_ [_index.ts:409_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L409)

### logIndex

• **logIndex**: _number_

_Defined in_ [_index.ts:406_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L406)

### parameters

• **parameters**: _ContractEventParameters_

_Defined in_ [_index.ts:410_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L410)

### txHash

• **txHash**: _string_

_Defined in_ [_index.ts:404_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L404)

### txIndex

• **txIndex**: _number_

_Defined in_ [_index.ts:405_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L405)

## Interface: ERC1155ApprovalForAllEvent

#### Hierarchy

* **ERC1155ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:284_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L284)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:283_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L283)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:282_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L282)

## Interface: ERC1155TransferBatchEvent

#### Hierarchy

* **ERC1155TransferBatchEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:267_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L267)

### ids

• **ids**: _BigNumber\[\]_

_Defined in_ [_index.ts:269_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L269)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:266_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L266)

### to

• **to**: _string_

_Defined in_ [_index.ts:268_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L268)

### values

• **values**: _BigNumber\[\]_

_Defined in_ [_index.ts:270_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L270)

## Interface: ERC1155TransferSingleEvent

#### Hierarchy

* **ERC1155TransferSingleEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:251_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L251)

### id

• **id**: _BigNumber_

_Defined in_ [_index.ts:253_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L253)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:250_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L250)

### to

• **to**: _string_

_Defined in_ [_index.ts:252_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L252)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:254_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L254)

## Interface: ERC20ApprovalEvent

#### Hierarchy

* **ERC20ApprovalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:208_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L208)

### spender

• **spender**: _string_

_Defined in_ [_index.ts:209_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L209)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:210_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L210)

## Interface: ERC20TransferEvent

#### Hierarchy

* **ERC20TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:196_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L196)

### to

• **to**: _string_

_Defined in_ [_index.ts:197_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L197)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:198_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L198)

## Interface: ERC721ApprovalEvent

#### Hierarchy

* **ERC721ApprovalEvent**

#### Properties

### approved

• **approved**: _string_

_Defined in_ [_index.ts:233_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L233)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:232_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L232)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:234_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L234)

## Interface: ERC721ApprovalForAllEvent

#### Hierarchy

* **ERC721ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:246_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L246)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:245_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L245)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:244_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L244)

## Interface: ERC721TransferEvent

#### Hierarchy

* **ERC721TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:220_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L220)

### to

• **to**: _string_

_Defined in_ [_index.ts:221_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L221)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:222_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L222)

## Interface: ExchangeCancelEvent

#### Hierarchy

* **ExchangeCancelEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:318_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L318)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:316_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L316)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:320_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L320)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:319_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L319)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:317_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L317)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:321_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L321)

## Interface: ExchangeCancelUpToEvent

#### Hierarchy

* **ExchangeCancelUpToEvent**

#### Properties

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:325_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L325)

### orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in_ [_index.ts:327_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L327)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:326_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L326)

## Interface: ExchangeFillEvent

#### Hierarchy

* **ExchangeFillEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:291_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L291)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:288_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L288)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:297_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L297)

### makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:292_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L292)

### makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:294_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L294)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:296_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L296)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:290_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L290)

### takerAddress

• **takerAddress**: _string_

_Defined in_ [_index.ts:289_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L289)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:298_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L298)

### takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:293_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L293)

### takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:295_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L295)

## Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired, or filled.

#### Hierarchy

* **OrderEvent**

#### Properties

### contractEvents

• **contractEvents**: [_ContractEvent_](reference.md#class-contractevent)_\[\]_

_Defined in_ [_index.ts:455_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L455)

### endState

• **endState**: [_OrderEventEndState_](reference.md#enumeration-ordereventendstate)

_Defined in_ [_index.ts:453_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L453)

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:454_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L454)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:451_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L451)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:452_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L452)

## Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were rejected.

#### Hierarchy

* **RejectedOrderInfo**

#### Properties

### kind

• **kind**: [_RejectedOrderKind_](reference.md#enumeration-rejectedorderkind)

_Defined in_ [_index.ts:505_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L505)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:503_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L503)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:504_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L504)

### status

• **status**: [_RejectedOrderStatus_](reference.md#class-rejectedorderstatus)

_Defined in_ [_index.ts:506_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L506)

## Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

#### Hierarchy

* **RejectedOrderStatus**

#### Properties

### code

• **code**: _string_

_Defined in_ [_index.ts:523_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L523)

### message

• **message**: _string_

_Defined in_ [_index.ts:524_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L524)

## Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

#### Hierarchy

* **ValidationResults**

#### Properties

### accepted

• **accepted**: [_AcceptedOrderInfo_](reference.md#class-acceptedorderinfo)_\[\]_

_Defined in_ [_index.ts:484_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L484)

### rejected

• **rejected**: [_RejectedOrderInfo_](reference.md#class-rejectedorderinfo)_\[\]_

_Defined in_ [_index.ts:485_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L485)

## Interface: WethDepositEvent

#### Hierarchy

* **WethDepositEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:347_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L347)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:348_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L348)

## Interface: WethWithdrawalEvent

#### Hierarchy

* **WethWithdrawalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:337_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L337)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:338_](https://github.com/0xProject/0x-mesh/blob/93b44a1/browser/ts/index.ts#L338)

