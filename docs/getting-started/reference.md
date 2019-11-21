# Browser API documentation

## Class: Mesh

The main class for this package. Has methods for receiving order events and sending orders through the 0x Mesh network.

#### Hierarchy

* **Mesh**

#### Constructors

### constructer

+ **new Mesh**\(`config`: [Config](reference.md#interface-config)\): [_Mesh_](reference.md#class-mesh)

_Defined in_ [_index.ts:524_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L524)

Instantiates a new Mesh instance.

**Parameters:**

| Name | Type | Description |
| :--- | :--- | :--- |
| `config` | [Config](reference.md#interface-config) | Configuration options for Mesh |

**Returns:** [_Mesh_](reference.md#class-mesh)

An instance of Mesh

#### Methods

### addOrdersAsync

▸ **addOrdersAsync**\(`orders`: SignedOrder\[\]\): _Promise‹_[_ValidationResults_](reference.md#interface-validationresults)_›_

_Defined in_ [_index.ts:594_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L594)

Validates and adds the given orders to Mesh. If an order is successfully added, Mesh will share it with any peers in the network and start watching it for changes \(e.g. filled, canceled, expired\). The returned promise will only be rejected if there was an error validating or adding the order; it will not be rejected for any invalid orders \(check results.rejected instead\).

**Parameters:**

| Name | Type | Description |
| :--- | :--- | :--- |
| `orders` | SignedOrder\[\] | An array of orders to add. |

**Returns:** _Promise‹_[_ValidationResults_](reference.md#interface-validationresults)_›_

Validation results for the given orders, indicating which orders were accepted and which were rejected.

### onError

▸ **onError**\(`handler`: function\): _void_

_Defined in_ [_index.ts:544_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L544)

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

_Defined in_ [_index.ts:559_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L559)

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

_Defined in_ [_index.ts:570_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L570)

Starts the Mesh node in the background. Mesh will automatically find peers in the network and begin receiving orders from them.

**Returns:** _Promise‹void›_

## Enumeration: OrderEventEndState

#### Enumeration members

### Added

• **Added**: = "ADDED"

_Defined in_ [_index.ts:391_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L391)

### Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in_ [_index.ts:394_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L394)

### Expired

• **Expired**: = "EXPIRED"

_Defined in_ [_index.ts:395_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L395)

### FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY\_INCREASED"

_Defined in_ [_index.ts:397_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L397)

### Filled

• **Filled**: = "FILLED"

_Defined in_ [_index.ts:392_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L392)

### FullyFilled

• **FullyFilled**: = "FULLY\_FILLED"

_Defined in_ [_index.ts:393_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L393)

### Invalid

• **Invalid**: = "INVALID"

_Defined in_ [_index.ts:390_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L390)

### Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in_ [_index.ts:396_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L396)

## Enumeration: RejectedOrderKind

A set of categories for rejected orders.

#### Enumeration members

### CoordinatorError

• **CoordinatorError**: = "COORDINATOR\_ERROR"

_Defined in_ [_index.ts:478_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L478)

### MeshError

• **MeshError**: = "MESH\_ERROR"

_Defined in_ [_index.ts:476_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L476)

### MeshValidation

• **MeshValidation**: = "MESH\_VALIDATION"

_Defined in_ [_index.ts:477_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L477)

### ZeroExValidation

• **ZeroExValidation**: = "ZEROEX\_VALIDATION"

_Defined in_ [_index.ts:475_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L475)

## Enumeration: Verbosity

#### Enumeration members

### Debug

• **Debug**: = 5

_Defined in_ [_index.ts:106_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L106)

### Error

• **Error**: = 2

_Defined in_ [_index.ts:103_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L103)

### Fatal

• **Fatal**: = 1

_Defined in_ [_index.ts:102_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L102)

### Info

• **Info**: = 4

_Defined in_ [_index.ts:105_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L105)

### Panic

• **Panic**: = 0

_Defined in_ [_index.ts:101_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L101)

### Trace

• **Trace**: = 6

_Defined in_ [_index.ts:107_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L107)

### Warn

• **Warn**: = 3

_Defined in_ [_index.ts:104_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L104)

## Interface: AcceptedOrderInfo

Info for any orders that were accepted.

#### Hierarchy

* **AcceptedOrderInfo**

#### Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:456_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L456)

### isNew

• **isNew**: _boolean_

_Defined in_ [_index.ts:457_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L457)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:454_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L454)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:455_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L455)

## Interface: Config

A set of configuration options for Mesh.

#### Hierarchy

* **Config**

#### Properties

### `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : _undefined \| number_

_Defined in_ [_index.ts:59_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L59)

### `Optional` bootstrapList

• **bootstrapList**? : _string\[\]_

_Defined in_ [_index.ts:48_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L48)

### `Optional` customContractAddresses

• **customContractAddresses**? : [_ContractAddresses_](reference.md#class-contractaddresses)

_Defined in_ [_index.ts:85_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L85)

### ethereumNetworkID

• **ethereumNetworkID**: _number_

_Defined in_ [_index.ts:40_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L40)

### `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : _undefined \| number_

_Defined in_ [_index.ts:68_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L68)

### ethereumRPCURL

• **ethereumRPCURL**: _string_

_Defined in_ [_index.ts:38_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L38)

### `Optional` orderExpirationBufferSeconds

• **orderExpirationBufferSeconds**? : _undefined \| number_

_Defined in_ [_index.ts:52_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L52)

### `Optional` useBootstrapList

• **useBootstrapList**? : _undefined \| false \| true_

_Defined in_ [_index.ts:43_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L43)

### `Optional` verbosity

• **verbosity**? : [_Verbosity_](reference.md#enumeration-verbosity)

_Defined in_ [_index.ts:35_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L35)

## Interface: ContractAddresses

#### Hierarchy

* **ContractAddresses**

#### Properties

### `Optional` coordinator

• **coordinator**? : _undefined \| string_

_Defined in_ [_index.ts:94_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L94)

### `Optional` coordinatorRegistry

• **coordinatorRegistry**? : _undefined \| string_

_Defined in_ [_index.ts:95_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L95)

### devUtils

• **devUtils**: _string_

_Defined in_ [_index.ts:90_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L90)

### erc1155Proxy

• **erc1155Proxy**: _string_

_Defined in_ [_index.ts:93_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L93)

### erc20Proxy

• **erc20Proxy**: _string_

_Defined in_ [_index.ts:91_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L91)

### erc721Proxy

• **erc721Proxy**: _string_

_Defined in_ [_index.ts:92_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L92)

### exchange

• **exchange**: _string_

_Defined in_ [_index.ts:89_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L89)

### `Optional` weth9

• **weth9**? : _undefined \| string_

_Defined in_ [_index.ts:96_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L96)

### `Optional` zrxToken

• **zrxToken**? : _undefined \| string_

_Defined in_ [_index.ts:97_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L97)

## Interface: ContractEvent

#### Hierarchy

* **ContractEvent**

#### Properties

### address

• **address**: _string_

_Defined in_ [_index.ts:372_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L372)

### blockHash

• **blockHash**: _string_

_Defined in_ [_index.ts:367_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L367)

### isRemoved

• **isRemoved**: _string_

_Defined in_ [_index.ts:371_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L371)

### kind

• **kind**: _ContractEventKind_

_Defined in_ [_index.ts:373_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L373)

### logIndex

• **logIndex**: _number_

_Defined in_ [_index.ts:370_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L370)

### parameters

• **parameters**: _ContractEventParameters_

_Defined in_ [_index.ts:374_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L374)

### txHash

• **txHash**: _string_

_Defined in_ [_index.ts:368_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L368)

### txIndex

• **txIndex**: _number_

_Defined in_ [_index.ts:369_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L369)

## Interface: ERC1155ApprovalForAllEvent

#### Hierarchy

* **ERC1155ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:248_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L248)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:247_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L247)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:246_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L246)

## Interface: ERC1155TransferBatchEvent

#### Hierarchy

* **ERC1155TransferBatchEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:231_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L231)

### ids

• **ids**: _BigNumber\[\]_

_Defined in_ [_index.ts:233_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L233)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:230_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L230)

### to

• **to**: _string_

_Defined in_ [_index.ts:232_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L232)

### values

• **values**: _BigNumber\[\]_

_Defined in_ [_index.ts:234_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L234)

## Interface: ERC1155TransferSingleEvent

#### Hierarchy

* **ERC1155TransferSingleEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:215_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L215)

### id

• **id**: _BigNumber_

_Defined in_ [_index.ts:217_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L217)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:214_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L214)

### to

• **to**: _string_

_Defined in_ [_index.ts:216_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L216)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:218_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L218)

## Interface: ERC20ApprovalEvent

#### Hierarchy

* **ERC20ApprovalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:171_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L171)

### spender

• **spender**: _string_

_Defined in_ [_index.ts:172_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L172)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:173_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L173)

## Interface: ERC20TransferEvent

#### Hierarchy

* **ERC20TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:159_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L159)

### to

• **to**: _string_

_Defined in_ [_index.ts:160_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L160)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:161_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L161)

## Interface: ERC721ApprovalEvent

#### Hierarchy

* **ERC721ApprovalEvent**

#### Properties

### approved

• **approved**: _string_

_Defined in_ [_index.ts:196_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L196)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:195_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L195)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:197_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L197)

## Interface: ERC721ApprovalForAllEvent

#### Hierarchy

* **ERC721ApprovalForAllEvent**

#### Properties

### approved

• **approved**: _boolean_

_Defined in_ [_index.ts:209_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L209)

### operator

• **operator**: _string_

_Defined in_ [_index.ts:208_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L208)

### owner

• **owner**: _string_

_Defined in_ [_index.ts:207_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L207)

## Interface: ERC721TransferEvent

#### Hierarchy

* **ERC721TransferEvent**

#### Properties

### from

• **from**: _string_

_Defined in_ [_index.ts:183_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L183)

### to

• **to**: _string_

_Defined in_ [_index.ts:184_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L184)

### tokenId

• **tokenId**: _BigNumber_

_Defined in_ [_index.ts:185_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L185)

## Interface: ExchangeCancelEvent

#### Hierarchy

* **ExchangeCancelEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:282_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L282)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:280_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L280)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:284_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L284)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:283_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L283)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:281_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L281)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:285_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L285)

## Interface: ExchangeCancelUpToEvent

#### Hierarchy

* **ExchangeCancelUpToEvent**

#### Properties

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:289_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L289)

### orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in_ [_index.ts:291_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L291)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:290_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L290)

## Interface: ExchangeFillEvent

#### Hierarchy

* **ExchangeFillEvent**

#### Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_index.ts:255_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L255)

### makerAddress

• **makerAddress**: _string_

_Defined in_ [_index.ts:252_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L252)

### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_index.ts:261_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L261)

### makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:256_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L256)

### makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:258_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L258)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:260_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L260)

### senderAddress

• **senderAddress**: _string_

_Defined in_ [_index.ts:254_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L254)

### takerAddress

• **takerAddress**: _string_

_Defined in_ [_index.ts:253_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L253)

### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_index.ts:262_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L262)

### takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in_ [_index.ts:257_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L257)

### takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in_ [_index.ts:259_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L259)

## Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired, or filled.

#### Hierarchy

* **OrderEvent**

#### Properties

### contractEvents

• **contractEvents**: [_ContractEvent_](reference.md#class-contractevent)_\[\]_

_Defined in_ [_index.ts:417_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L417)

### endState

• **endState**: [_OrderEventEndState_](reference.md#enumeration-ordereventendstate)

_Defined in_ [_index.ts:415_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L415)

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in_ [_index.ts:416_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L416)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:413_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L413)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:414_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L414)

## Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were rejected.

#### Hierarchy

* **RejectedOrderInfo**

#### Properties

### kind

• **kind**: [_RejectedOrderKind_](reference.md#enumeration-rejectedorderkind)

_Defined in_ [_index.ts:467_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L467)

### orderHash

• **orderHash**: _string_

_Defined in_ [_index.ts:465_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L465)

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in_ [_index.ts:466_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L466)

### status

• **status**: [_RejectedOrderStatus_](reference.md#class-rejectedorderstatus)

_Defined in_ [_index.ts:468_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L468)

## Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

#### Hierarchy

* **RejectedOrderStatus**

#### Properties

### code

• **code**: _string_

_Defined in_ [_index.ts:485_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L485)

### message

• **message**: _string_

_Defined in_ [_index.ts:486_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L486)

## Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

#### Hierarchy

* **ValidationResults**

#### Properties

### accepted

• **accepted**: [_AcceptedOrderInfo_](reference.md#class-acceptedorderinfo)_\[\]_

_Defined in_ [_index.ts:446_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L446)

### rejected

• **rejected**: [_RejectedOrderInfo_](reference.md#class-rejectedorderinfo)_\[\]_

_Defined in_ [_index.ts:447_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L447)

## Interface: WethDepositEvent

#### Hierarchy

* **WethDepositEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:311_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L311)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:312_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L312)

## Interface: WethWithdrawalEvent

#### Hierarchy

* **WethWithdrawalEvent**

#### Properties

### owner

• **owner**: _string_

_Defined in_ [_index.ts:301_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L301)

### value

• **value**: _BigNumber_

_Defined in_ [_index.ts:302_](https://github.com/0xProject/0x-mesh/blob/b6b8a86/browser/ts/index.ts#L302)

