# Class: MeshGraphQLClient

### Hierarchy

-   **MeshGraphQLClient**

### Constructors

## constructer

\+ **new MeshGraphQLClient**(`linkConfig`: [LinkConfig](#interface-linkconfig)): _[MeshGraphQLClient](#class-meshgraphqlclient)_

_Defined in [packages/mesh-graphql-client/src/index.ts:92](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L92)_

**Parameters:**

| Name         | Type                                |
| ------------ | ----------------------------------- |
| `linkConfig` | [LinkConfig](#interface-linkconfig) |

**Returns:** _[MeshGraphQLClient](#class-meshgraphqlclient)_

### Methods

## addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[], `pinned`: boolean, `opts?`: [AddOrdersOpts](#interface-addordersopts)): _Promise‹[AddOrdersResults](#interface-addordersresults)‹[OrderWithMetadata](#interface-orderwithmetadata), SignedOrder››_

_Defined in [packages/mesh-graphql-client/src/index.ts:172](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L172)_

**Parameters:**

| Name     | Type                                      | Default |
| -------- | ----------------------------------------- | ------- |
| `orders` | SignedOrder[]                             | -       |
| `pinned` | boolean                                   | true    |
| `opts?`  | [AddOrdersOpts](#interface-addordersopts) | -       |

**Returns:** _Promise‹[AddOrdersResults](#interface-addordersresults)‹[OrderWithMetadata](#interface-orderwithmetadata), SignedOrder››_

---

## addOrdersV4Async

▸ **addOrdersV4Async**(`orders`: [SignedOrderV4](#signedorderv4)››\*

_Defined in [packages/mesh-graphql-client/src/index.ts:201](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L201)_

**Parameters:**

| Name     | Type                                      | Default |
| -------- | ----------------------------------------- | ------- |
| `orders` | [SignedOrderV4](#signedorderv4)[]         | -       |
| `pinned` | boolean                                   | true    |
| `opts?`  | [AddOrdersOpts](#interface-addordersopts) | -       |

**Returns:** _Promise‹[AddOrdersResults](#signedorderv4)››_

---

## findOrdersAsync

▸ **findOrdersAsync**(`query`: [OrderQuery](#interface-orderquery)): _Promise‹[OrderWithMetadata](#interface-orderwithmetadata)[]›_

_Defined in [packages/mesh-graphql-client/src/index.ts:266](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L266)_

**Parameters:**

| Name    | Type                                | Default                                                  |
| ------- | ----------------------------------- | -------------------------------------------------------- |
| `query` | [OrderQuery](#interface-orderquery) | { sort: [], filters: [], limit: defaultOrderQueryLimit } |

**Returns:** _Promise‹[OrderWithMetadata](#interface-orderwithmetadata)[]›_

---

## findOrdersV4Async

▸ **findOrdersV4Async**(`query`: [OrderQuery](#interface-orderquery)): _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4)[]›_

_Defined in [packages/mesh-graphql-client/src/index.ts:284](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L284)_

**Parameters:**

| Name    | Type                                | Default                                                  |
| ------- | ----------------------------------- | -------------------------------------------------------- |
| `query` | [OrderQuery](#interface-orderquery) | { sort: [], filters: [], limit: defaultOrderQueryLimit } |

**Returns:** _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4)[]›_

---

## getOrderAsync

▸ **getOrderAsync**(`hash`: string): _Promise‹[OrderWithMetadata](#interface-orderwithmetadata) | null›_

_Defined in [packages/mesh-graphql-client/src/index.ts:232](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L232)_

**Parameters:**

| Name   | Type   |
| ------ | ------ |
| `hash` | string |

**Returns:** _Promise‹[OrderWithMetadata](#interface-orderwithmetadata) | null›_

---

## getOrderV4Async

▸ **getOrderV4Async**(`hash`: string): _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4) | null›_

_Defined in [packages/mesh-graphql-client/src/index.ts:249](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L249)_

**Parameters:**

| Name   | Type   |
| ------ | ------ |
| `hash` | string |

**Returns:** _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4) | null›_

---

## getStatsAsync

▸ **getStatsAsync**(): _Promise‹[Stats](#interface-stats)›_

_Defined in [packages/mesh-graphql-client/src/index.ts:160](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L160)_

**Returns:** _Promise‹[Stats](#interface-stats)›_

---

## onOrderEvents

▸ **onOrderEvents**(): _Observable‹[OrderEvent](#interface-orderevent)[]›_

_Defined in [packages/mesh-graphql-client/src/index.ts:306](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L306)_

**Returns:** _Observable‹[OrderEvent](#interface-orderevent)[]›_

---

## onReconnected

▸ **onReconnected**(`cb`: function): _void_

_Defined in [packages/mesh-graphql-client/src/index.ts:302](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L302)_

**Parameters:**

▪ **cb**: _function_

▸ (): _void_

**Returns:** _void_

---

## rawQueryAsync

▸ **rawQueryAsync**<**T**, **TVariables**>(`options`: QueryOptions‹TVariables›): _Promise‹ApolloQueryResult‹T››_

_Defined in [packages/mesh-graphql-client/src/index.ts:353](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L353)_

**Type parameters:**

▪ **T**

▪ **TVariables**

**Parameters:**

| Name      | Type                     |
| --------- | ------------------------ |
| `options` | QueryOptions‹TVariables› |

**Returns:** _Promise‹ApolloQueryResult‹T››_

<hr />

# Enumeration: ContractEventKind

### Enumeration members

## ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:197](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L197)_

---

## ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:199](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L199)_

---

## ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:198](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L198)_

---

## ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:193](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L193)_

---

## ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:192](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L192)_

---

## ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:195](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L195)_

---

## ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:196](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L196)_

---

## ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:194](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L194)_

---

## ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:201](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L201)_

---

## ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:202](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L202)_

---

## ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:200](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L200)_

---

## WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:203](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L203)_

---

## WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:204](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L204)_

<hr />

# Enumeration: FilterKind

### Enumeration members

## Equal

• **Equal**: = "EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:243](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L243)_

---

## Greater

• **Greater**: = "GREATER"

_Defined in [packages/mesh-graphql-client/src/types.ts:245](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L245)_

---

## GreaterOrEqual

• **GreaterOrEqual**: = "GREATER_OR_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:246](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L246)_

---

## Less

• **Less**: = "LESS"

_Defined in [packages/mesh-graphql-client/src/types.ts:247](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L247)_

---

## LessOrEqual

• **LessOrEqual**: = "LESS_OR_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:248](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L248)_

---

## NotEqual

• **NotEqual**: = "NOT_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:244](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L244)_

<hr />

# Enumeration: OrderEventEndState

### Enumeration members

## Added

• **Added**: = "ADDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:210](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L210)_

---

## Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:216](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L216)_

---

## Expired

• **Expired**: = "EXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:218](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L218)_

---

## FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

_Defined in [packages/mesh-graphql-client/src/types.ts:227](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L227)_

---

## Filled

• **Filled**: = "FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:212](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L212)_

---

## FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:214](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L214)_

---

## Invalid

• **Invalid**: = "INVALID"

_Defined in [packages/mesh-graphql-client/src/types.ts:220](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L220)_

---

## StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

_Defined in [packages/mesh-graphql-client/src/types.ts:232](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L232)_

---

## Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:222](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L222)_

---

## Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:224](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L224)_

<hr />

# Enumeration: RejectedOrderCode

### Enumeration members

## DatabaseFullOfOrders

• **DatabaseFullOfOrders**: = "DATABASE_FULL_OF_ORDERS"

_Defined in [packages/mesh-graphql-client/src/types.ts:167](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L167)_

---

## EthRpcRequestFailed

• **EthRpcRequestFailed**: = "ETH_RPC_REQUEST_FAILED"

_Defined in [packages/mesh-graphql-client/src/types.ts:148](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L148)_

---

## IncorrectExchangeAddress

• **IncorrectExchangeAddress**: = "INCORRECT_EXCHANGE_ADDRESS"

_Defined in [packages/mesh-graphql-client/src/types.ts:165](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L165)_

---

## InternalError

• **InternalError**: = "INTERNAL_ERROR"

_Defined in [packages/mesh-graphql-client/src/types.ts:161](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L161)_

---

## MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MAX_ORDER_SIZE_EXCEEDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:162](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L162)_

---

## OrderAlreadyStoredAndUnfillable

• **OrderAlreadyStoredAndUnfillable**: = "ORDER_ALREADY_STORED_AND_UNFILLABLE"

_Defined in [packages/mesh-graphql-client/src/types.ts:163](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L163)_

---

## OrderCancelled

• **OrderCancelled**: = "ORDER_CANCELLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:153](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L153)_

---

## OrderExpired

• **OrderExpired**: = "ORDER_EXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:151](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L151)_

---

## OrderForIncorrectChain

• **OrderForIncorrectChain**: = "ORDER_FOR_INCORRECT_CHAIN"

_Defined in [packages/mesh-graphql-client/src/types.ts:164](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L164)_

---

## OrderFullyFilled

• **OrderFullyFilled**: = "ORDER_FULLY_FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:152](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L152)_

---

## OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "ORDER_HAS_INVALID_MAKER_ASSET_AMOUNT"

_Defined in [packages/mesh-graphql-client/src/types.ts:149](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L149)_

---

## OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "ORDER_HAS_INVALID_MAKER_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:155](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L155)_

---

## OrderHasInvalidMakerFeeAssetData

• **OrderHasInvalidMakerFeeAssetData**: = "ORDER_HAS_INVALID_MAKER_FEE_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:156](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L156)_

---

## OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "ORDER_HAS_INVALID_SIGNATURE"

_Defined in [packages/mesh-graphql-client/src/types.ts:159](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L159)_

---

## OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "ORDER_HAS_INVALID_TAKER_ASSET_AMOUNT"

_Defined in [packages/mesh-graphql-client/src/types.ts:150](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L150)_

---

## OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "ORDER_HAS_INVALID_TAKER_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:157](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L157)_

---

## OrderHasInvalidTakerFeeAssetData

• **OrderHasInvalidTakerFeeAssetData**: = "ORDER_HAS_INVALID_TAKER_FEE_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:158](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L158)_

---

## OrderMaxExpirationExceeded

• **OrderMaxExpirationExceeded**: = "ORDER_MAX_EXPIRATION_EXCEEDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:160](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L160)_

---

## OrderUnfunded

• **OrderUnfunded**: = "ORDER_UNFUNDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:154](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L154)_

---

## SenderAddressNotAllowed

• **SenderAddressNotAllowed**: = "SENDER_ADDRESS_NOT_ALLOWED"

_Defined in [packages/mesh-graphql-client/src/types.ts:166](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L166)_

---

## TakerAddressNotAllowed

• **TakerAddressNotAllowed**: = "TAKER_ADDRESS_NOT_ALLOWED"

_Defined in [packages/mesh-graphql-client/src/types.ts:168](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L168)_

<hr />

# Enumeration: SortDirection

### Enumeration members

## Asc

• **Asc**: = "ASC"

_Defined in [packages/mesh-graphql-client/src/types.ts:238](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L238)_

---

## Desc

• **Desc**: = "DESC"

_Defined in [packages/mesh-graphql-client/src/types.ts:239](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L239)_

<hr />

# Interface: LinkConfig

### Hierarchy

-   **LinkConfig**

### Properties

## `Optional` httpUrl

• **httpUrl**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/index.ts:84](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L84)_

---

## `Optional` webSocketUrl

• **webSocketUrl**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/index.ts:85](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/index.ts#L85)_

<hr />

# Interface: AcceptedOrderResult <**T**>

### Type parameters

▪ **T**: _[GenericOrderWithMetadata](#genericorderwithmetadata)_

### Hierarchy

-   **AcceptedOrderResult**

### Properties

## isNew

• **isNew**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:130](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L130)_

---

## order

• **order**: _T_

_Defined in [packages/mesh-graphql-client/src/types.ts:127](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L127)_

<hr />

# Interface: AddOrdersOpts

### Hierarchy

-   **AddOrdersOpts**

### Properties

## `Optional` keepCancelled

• **keepCancelled**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:27](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L27)_

---

## `Optional` keepExpired

• **keepExpired**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:28](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L28)_

---

## `Optional` keepFullyFilled

• **keepFullyFilled**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:29](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L29)_

---

## `Optional` keepUnfunded

• **keepUnfunded**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:30](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L30)_

<hr />

# Interface: AddOrdersResponse <**T, K**>

### Type parameters

▪ **T**: _[GenericStringifiedOrderWithMetadata](#genericstringifiedorderwithmetadata)_

▪ **K**: _[GenericStringifiedSignedOrders](#genericstringifiedsignedorders)_

### Hierarchy

-   **AddOrdersResponse**

### Properties

## addOrders

• **addOrders**: _[StringifiedAddOrdersResults](#interface-stringifiedaddordersresults)‹T, K›_

_Defined in [packages/mesh-graphql-client/src/types.ts:46](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L46)_

<hr />

# Interface: AddOrdersResponseV4 <**T, K**>

### Type parameters

▪ **T**: _[GenericStringifiedOrderWithMetadata](#genericstringifiedorderwithmetadata)_

▪ **K**: _[GenericStringifiedSignedOrders](#genericstringifiedsignedorders)_

### Hierarchy

-   **AddOrdersResponseV4**

### Properties

## addOrdersV4

• **addOrdersV4**: _[StringifiedAddOrdersResults](#interface-stringifiedaddordersresults)‹T, K›_

_Defined in [packages/mesh-graphql-client/src/types.ts:53](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L53)_

<hr />

# Interface: AddOrdersResults <**T, K**>

### Type parameters

▪ **T**: _[GenericOrderWithMetadata](#genericorderwithmetadata)_

▪ **K**: _[GenericSignedOrder](#genericsignedorder)_

### Hierarchy

-   **AddOrdersResults**

### Properties

## accepted

• **accepted**: _[AcceptedOrderResult](#interface-acceptedorderresult)‹T›[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:118](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L118)_

---

## rejected

• **rejected**: _[RejectedOrderResult](#interface-rejectedorderresult)‹K›[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:121](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L121)_

<hr />

# Interface: ContractEvent

### Hierarchy

-   **ContractEvent**

### Properties

## address

• **address**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:185](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L185)_

---

## blockHash

• **blockHash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:180](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L180)_

---

## isRemoved

• **isRemoved**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:184](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L184)_

---

## kind

• **kind**: _[ContractEventKind](#enumeration-contracteventkind)_

_Defined in [packages/mesh-graphql-client/src/types.ts:186](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L186)_

---

## logIndex

• **logIndex**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:183](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L183)_

---

## parameters

• **parameters**: _any_

_Defined in [packages/mesh-graphql-client/src/types.ts:188](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L188)_

---

## txHash

• **txHash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:181](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L181)_

---

## txIndex

• **txIndex**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:182](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L182)_

<hr />

# Interface: LatestBlock

### Hierarchy

-   **LatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:99](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L99)_

---

## number

• **number**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:98](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L98)_

<hr />

# Interface: OrderEvent

### Hierarchy

-   **OrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:176](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L176)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [packages/mesh-graphql-client/src/types.ts:175](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L175)_

---

## `Optional` order

• **order**? : _[OrderWithMetadata](#interface-orderwithmetadata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:173](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L173)_

---

## `Optional` orderv4

• **orderv4**? : _[OrderWithMetadataV4](#interface-orderwithmetadatav4)_

_Defined in [packages/mesh-graphql-client/src/types.ts:174](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L174)_

---

## timestampMs

• **timestampMs**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:172](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L172)_

<hr />

# Interface: OrderEventResponse

### Hierarchy

-   **OrderEventResponse**

### Properties

## orderEvents

• **orderEvents**: _[StringifiedOrderEvent](#interface-stringifiedorderevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:73](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L73)_

<hr />

# Interface: OrderFilter

### Hierarchy

-   **OrderFilter**

### Properties

## field

• **field**: _[OrderField](#orderfield)_

_Defined in [packages/mesh-graphql-client/src/types.ts:257](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L257)_

---

## kind

• **kind**: _[FilterKind](#enumeration-filterkind)_

_Defined in [packages/mesh-graphql-client/src/types.ts:258](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L258)_

---

## value

• **value**: _OrderWithMetadata[OrderField]_

_Defined in [packages/mesh-graphql-client/src/types.ts:259](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L259)_

<hr />

# Interface: OrderQuery

### Hierarchy

-   **OrderQuery**

### Properties

## `Optional` filters

• **filters**? : _[OrderFilter](#interface-orderfilter)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:263](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L263)_

---

## `Optional` limit

• **limit**? : _undefined | number_

_Defined in [packages/mesh-graphql-client/src/types.ts:265](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L265)_

---

## `Optional` sort

• **sort**? : _[OrderSort](#interface-ordersort)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:264](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L264)_

<hr />

# Interface: OrderResponse

### Hierarchy

-   **OrderResponse**

### Properties

## order

• **order**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:57](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L57)_

<hr />

# Interface: OrderResponseV4

### Hierarchy

-   **OrderResponseV4**

### Properties

## orderv4

• **orderv4**: _[StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:61](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L61)_

<hr />

# Interface: OrderSort

### Hierarchy

-   **OrderSort**

### Properties

## direction

• **direction**: _[SortDirection](#enumeration-sortdirection)_

_Defined in [packages/mesh-graphql-client/src/types.ts:253](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L253)_

---

## field

• **field**: _[OrderField](#orderfield)_

_Defined in [packages/mesh-graphql-client/src/types.ts:252](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L252)_

<hr />

# Interface: OrdersResponse

### Hierarchy

-   **OrdersResponse**

### Properties

## orders

• **orders**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:65](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L65)_

<hr />

# Interface: OrdersResponseV4

### Hierarchy

-   **OrdersResponseV4**

### Properties

## ordersv4

• **ordersv4**: _[StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:69](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L69)_

<hr />

# Interface: OrderWithMetadata

### Hierarchy

-   SignedOrder

    ↳ **OrderWithMetadata**

### Properties

## chainId

• **chainId**: _number_

Defined in node_modules/@0x/types/lib/index.d.ts:4

---

## exchangeAddress

• **exchangeAddress**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:5

---

## expirationTimeSeconds

• **expirationTimeSeconds**: _BigNumber_

Defined in node_modules/@0x/types/lib/index.d.ts:14

---

## feeRecipientAddress

• **feeRecipientAddress**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:8

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:104](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L104)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:103](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L103)_

---

## makerAddress

• **makerAddress**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:6

---

## makerAssetAmount

• **makerAssetAmount**: _BigNumber_

Defined in node_modules/@0x/types/lib/index.d.ts:10

---

## makerAssetData

• **makerAssetData**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:16

---

## makerFee

• **makerFee**: _BigNumber_

Defined in node_modules/@0x/types/lib/index.d.ts:12

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:18

---

## salt

• **salt**: _BigNumber_

Defined in node_modules/@0x/types/lib/index.d.ts:15

---

## senderAddress

• **senderAddress**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:9

---

## signature

• **signature**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:22

---

## takerAddress

• **takerAddress**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:7

---

## takerAssetAmount

• **takerAssetAmount**: _BigNumber_

Defined in node_modules/@0x/types/lib/index.d.ts:11

---

## takerAssetData

• **takerAssetData**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:17

---

## takerFee

• **takerFee**: _BigNumber_

Defined in node_modules/@0x/types/lib/index.d.ts:13

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

Defined in node_modules/@0x/types/lib/index.d.ts:19

<hr />

# Interface: OrderWithMetadataV4

### Hierarchy

-   object & object

    ↳ **OrderWithMetadataV4**

### Properties

## chainId

• **chainId**: _number_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:31

---

## expiry

• **expiry**: _BigNumber_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:29

---

## feeRecipient

• **feeRecipient**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:21

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:112](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L112)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:111](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L111)_

---

## maker

• **maker**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:26

---

## makerAmount

• **makerAmount**: _BigNumber_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:24

---

## makerToken

• **makerToken**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:22

---

## pool

• **pool**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:28

---

## salt

• **salt**: _BigNumber_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:30

---

## sender

• **sender**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:20

---

## signature

• **signature**: _Signature_

_Defined in [packages/mesh-graphql-client/src/types.ts:108](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L108)_

---

## taker

• **taker**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:27

---

## takerAmount

• **takerAmount**: _BigNumber_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:25

---

## takerToken

• **takerToken**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:23

---

## takerTokenFeeAmount

• **takerTokenFeeAmount**: _BigNumber_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:19

---

## verifyingContract

• **verifyingContract**: _string_

Defined in node_modules/@0x/protocol-utils/lib/src/orders.d.ts:32

<hr />

# Interface: RejectedOrderResult <**K**>

### Type parameters

▪ **K**: _[GenericSignedOrder](#genericsignedorder)_

### Hierarchy

-   **RejectedOrderResult**

### Properties

## code

• **code**: _[RejectedOrderCode](#enumeration-rejectedordercode)_

_Defined in [packages/mesh-graphql-client/src/types.ts:141](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L141)_

---

## `Optional` hash

• **hash**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/types.ts:135](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L135)_

---

## message

• **message**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:144](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L144)_

---

## order

• **order**: _K_

_Defined in [packages/mesh-graphql-client/src/types.ts:138](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L138)_

<hr />

# Interface: Stats

### Hierarchy

-   **Stats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:94](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L94)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:93](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L93)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:82](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L82)_

---

## latestBlock

• **latestBlock**: _[LatestBlock](#interface-latestblock)_

_Defined in [packages/mesh-graphql-client/src/types.ts:83](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L83)_

---

## maxExpirationTime

• **maxExpirationTime**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:91](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L91)_

---

## numOrders

• **numOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:85](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L85)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:87](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L87)_

---

## numOrdersIncludingRemovedV4

• **numOrdersIncludingRemovedV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:88](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L88)_

---

## numOrdersV4

• **numOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:86](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L86)_

---

## numPeers

• **numPeers**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:84](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L84)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:89](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L89)_

---

## numPinnedOrdersV4

• **numPinnedOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:90](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L90)_

---

## peerID

• **peerID**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:81](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L81)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:78](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L78)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:79](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L79)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:80](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L80)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _Date_

_Defined in [packages/mesh-graphql-client/src/types.ts:92](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L92)_

---

## version

• **version**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:77](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L77)_

<hr />

# Interface: StatsResponse

### Hierarchy

-   **StatsResponse**

### Properties

## stats

• **stats**: _[StringifiedStats](#interface-stringifiedstats)_

_Defined in [packages/mesh-graphql-client/src/types.ts:34](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L34)_

<hr />

# Interface: StringifiedAcceptedOrderResult <**T**>

### Type parameters

▪ **T**: _[GenericStringifiedOrderWithMetadata](#genericstringifiedorderwithmetadata)_

### Hierarchy

-   **StringifiedAcceptedOrderResult**

### Properties

## isNew

• **isNew**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:354](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L354)_

---

## order

• **order**: _T_

_Defined in [packages/mesh-graphql-client/src/types.ts:353](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L353)_

<hr />

# Interface: StringifiedAddOrdersResults <**T, K**>

### Type parameters

▪ **T**: _[GenericStringifiedOrderWithMetadata](#genericstringifiedorderwithmetadata)_

▪ **K**: _[GenericStringifiedSignedOrders](#genericstringifiedsignedorders)_

### Hierarchy

-   **StringifiedAddOrdersResults**

### Properties

## accepted

• **accepted**: _[StringifiedAcceptedOrderResult](#interface-stringifiedacceptedorderresult)‹T›[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:348](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L348)_

---

## rejected

• **rejected**: _[StringifiedRejectedOrderResult](#interface-stringifiedrejectedorderresult)‹K›[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:349](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L349)_

<hr />

# Interface: StringifiedLatestBlock

### Hierarchy

-   **StringifiedLatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:270](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L270)_

---

## number

• **number**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:269](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L269)_

<hr />

# Interface: StringifiedOrderEvent

### Hierarchy

-   **StringifiedOrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:370](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L370)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [packages/mesh-graphql-client/src/types.ts:368](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L368)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:369](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L369)_

---

## order

• **order**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:366](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L366)_

---

## orderv4

• **orderv4**: _[StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:367](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L367)_

---

## timestamp

• **timestamp**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:365](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L365)_

<hr />

# Interface: StringifiedOrderWithMetadata

### Hierarchy

-   [StringifiedSignedOrder](#interface-stringifiedsignedorder)

    ↳ **StringifiedOrderWithMetadata**

### Properties

## chainId

• **chainId**: _string_

_Inherited from [StringifiedSignedOrder](#chainid)_

_Overrides [StringifiedSignedOrder](#chainid)_

_Defined in [packages/mesh-graphql-client/src/types.ts:22](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L22)_

---

## exchangeAddress

• **exchangeAddress**: _string_

_Inherited from [StringifiedSignedOrder](#exchangeaddress)_

_Overrides [StringifiedSignedOrder](#exchangeaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:21](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L21)_

---

## expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Inherited from [StringifiedSignedOrder](#expirationtimeseconds)_

_Overrides [StringifiedSignedOrder](#expirationtimeseconds)_

_Defined in [packages/mesh-graphql-client/src/types.ts:19](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L19)_

---

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Inherited from [StringifiedSignedOrder](#feerecipientaddress)_

_Overrides [StringifiedSignedOrder](#feerecipientaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:18](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L18)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:341](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L341)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:340](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L340)_

---

## makerAddress

• **makerAddress**: _string_

_Inherited from [StringifiedSignedOrder](#makeraddress)_

_Overrides [StringifiedSignedOrder](#makeraddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:7](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L7)_

---

## makerAssetAmount

• **makerAssetAmount**: _string_

_Inherited from [StringifiedSignedOrder](#makerassetamount)_

_Overrides [StringifiedSignedOrder](#makerassetamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:9](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L9)_

---

## makerAssetData

• **makerAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#makerassetdata)_

_Overrides [StringifiedSignedOrder](#makerassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:8](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L8)_

---

## makerFee

• **makerFee**: _string_

_Inherited from [StringifiedSignedOrder](#makerfee)_

_Overrides [StringifiedSignedOrder](#makerfee)_

_Defined in [packages/mesh-graphql-client/src/types.ts:10](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L10)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#makerfeeassetdata)_

_Overrides [StringifiedSignedOrder](#makerfeeassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:11](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L11)_

---

## salt

• **salt**: _string_

_Inherited from [StringifiedSignedOrder](#salt)_

_Overrides [StringifiedSignedOrder](#salt)_

_Defined in [packages/mesh-graphql-client/src/types.ts:20](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L20)_

---

## senderAddress

• **senderAddress**: _string_

_Inherited from [StringifiedSignedOrder](#senderaddress)_

_Overrides [StringifiedSignedOrder](#senderaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:17](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L17)_

---

## signature

• **signature**: _string_

_Inherited from [StringifiedSignedOrder](#signature)_

_Overrides [StringifiedSignedOrder](#signature)_

_Defined in [packages/mesh-graphql-client/src/types.ts:23](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L23)_

---

## takerAddress

• **takerAddress**: _string_

_Inherited from [StringifiedSignedOrder](#takeraddress)_

_Overrides [StringifiedSignedOrder](#takeraddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:12](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L12)_

---

## takerAssetAmount

• **takerAssetAmount**: _string_

_Inherited from [StringifiedSignedOrder](#takerassetamount)_

_Overrides [StringifiedSignedOrder](#takerassetamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:15](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L15)_

---

## takerAssetData

• **takerAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#takerassetdata)_

_Overrides [StringifiedSignedOrder](#takerassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:13](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L13)_

---

## takerFee

• **takerFee**: _string_

_Inherited from [StringifiedSignedOrder](#takerfee)_

_Overrides [StringifiedSignedOrder](#takerfee)_

_Defined in [packages/mesh-graphql-client/src/types.ts:16](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L16)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#takerfeeassetdata)_

_Overrides [StringifiedSignedOrder](#takerfeeassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:14](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L14)_

<hr />

# Interface: StringifiedOrderWithMetadataV4

### Hierarchy

-   [StringifiedSignedOrderV4](#interface-stringifiedsignedorderv4)

    ↳ **StringifiedOrderWithMetadataV4**

### Properties

## chainId

• **chainId**: _string_

_Inherited from [StringifiedSignedOrderV4](#chainid)_

_Defined in [packages/mesh-graphql-client/src/types.ts:315](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L315)_

---

## expiry

• **expiry**: _string_

_Inherited from [StringifiedSignedOrderV4](#expiry)_

_Defined in [packages/mesh-graphql-client/src/types.ts:327](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L327)_

---

## feeRecipient

• **feeRecipient**: _string_

_Inherited from [StringifiedSignedOrderV4](#feerecipient)_

_Defined in [packages/mesh-graphql-client/src/types.ts:325](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L325)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:336](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L336)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:335](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L335)_

---

## maker

• **maker**: _string_

_Inherited from [StringifiedSignedOrderV4](#maker)_

_Defined in [packages/mesh-graphql-client/src/types.ts:322](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L322)_

---

## makerAmount

• **makerAmount**: _string_

_Inherited from [StringifiedSignedOrderV4](#makeramount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:319](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L319)_

---

## makerToken

• **makerToken**: _string_

_Inherited from [StringifiedSignedOrderV4](#makertoken)_

_Defined in [packages/mesh-graphql-client/src/types.ts:317](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L317)_

---

## pool

• **pool**: _string_

_Inherited from [StringifiedSignedOrderV4](#pool)_

_Defined in [packages/mesh-graphql-client/src/types.ts:326](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L326)_

---

## salt

• **salt**: _string_

_Inherited from [StringifiedSignedOrderV4](#salt)_

_Defined in [packages/mesh-graphql-client/src/types.ts:328](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L328)_

---

## sender

• **sender**: _string_

_Inherited from [StringifiedSignedOrderV4](#sender)_

_Defined in [packages/mesh-graphql-client/src/types.ts:324](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L324)_

---

## signatureR

• **signatureR**: _string_

_Inherited from [StringifiedSignedOrderV4](#signaturer)_

_Defined in [packages/mesh-graphql-client/src/types.ts:330](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L330)_

---

## signatureS

• **signatureS**: _string_

_Inherited from [StringifiedSignedOrderV4](#signatures)_

_Defined in [packages/mesh-graphql-client/src/types.ts:331](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L331)_

---

## signatureType

• **signatureType**: _string_

_Inherited from [StringifiedSignedOrderV4](#signaturetype)_

_Defined in [packages/mesh-graphql-client/src/types.ts:329](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L329)_

---

## signatureV

• **signatureV**: _string_

_Inherited from [StringifiedSignedOrderV4](#signaturev)_

_Defined in [packages/mesh-graphql-client/src/types.ts:332](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L332)_

---

## taker

• **taker**: _string_

_Inherited from [StringifiedSignedOrderV4](#taker)_

_Defined in [packages/mesh-graphql-client/src/types.ts:323](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L323)_

---

## takerAmount

• **takerAmount**: _string_

_Inherited from [StringifiedSignedOrderV4](#takeramount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:320](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L320)_

---

## takerToken

• **takerToken**: _string_

_Inherited from [StringifiedSignedOrderV4](#takertoken)_

_Defined in [packages/mesh-graphql-client/src/types.ts:318](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L318)_

---

## takerTokenFeeAmount

• **takerTokenFeeAmount**: _string_

_Inherited from [StringifiedSignedOrderV4](#takertokenfeeamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:321](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L321)_

---

## verifyingContract

• **verifyingContract**: _string_

_Inherited from [StringifiedSignedOrderV4](#verifyingcontract)_

_Defined in [packages/mesh-graphql-client/src/types.ts:316](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L316)_

<hr />

# Interface: StringifiedRejectedOrderResult <**K**>

### Type parameters

▪ **K**: _[GenericStringifiedSignedOrders](#genericstringifiedsignedorders)_

### Hierarchy

-   **StringifiedRejectedOrderResult**

### Properties

## code

• **code**: _[RejectedOrderCode](#enumeration-rejectedordercode)_

_Defined in [packages/mesh-graphql-client/src/types.ts:360](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L360)_

---

## `Optional` hash

• **hash**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/types.ts:358](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L358)_

---

## message

• **message**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:361](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L361)_

---

## order

• **order**: _K_

_Defined in [packages/mesh-graphql-client/src/types.ts:359](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L359)_

<hr />

# Interface: StringifiedSignedOrder

### Hierarchy

-   **StringifiedSignedOrder**

    ↳ [StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)

### Properties

## chainId

• **chainId**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:22](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L22)_

_Defined in [packages/mesh-graphql-client/src/types.ts:295](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L295)_

---

## exchangeAddress

• **exchangeAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:21](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L21)_

_Defined in [packages/mesh-graphql-client/src/types.ts:296](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L296)_

---

## expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:19](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L19)_

_Defined in [packages/mesh-graphql-client/src/types.ts:305](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L305)_

---

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:18](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L18)_

_Defined in [packages/mesh-graphql-client/src/types.ts:299](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L299)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:7](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L7)_

_Defined in [packages/mesh-graphql-client/src/types.ts:297](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L297)_

---

## makerAssetAmount

• **makerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:9](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L9)_

_Defined in [packages/mesh-graphql-client/src/types.ts:301](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L301)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:8](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L8)_

_Defined in [packages/mesh-graphql-client/src/types.ts:307](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L307)_

---

## makerFee

• **makerFee**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:10](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L10)_

_Defined in [packages/mesh-graphql-client/src/types.ts:303](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L303)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:11](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L11)_

_Defined in [packages/mesh-graphql-client/src/types.ts:309](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L309)_

---

## salt

• **salt**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:20](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L20)_

_Defined in [packages/mesh-graphql-client/src/types.ts:306](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L306)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:17](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L17)_

_Defined in [packages/mesh-graphql-client/src/types.ts:300](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L300)_

---

## signature

• **signature**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:23](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L23)_

_Defined in [packages/mesh-graphql-client/src/types.ts:311](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L311)_

---

## takerAddress

• **takerAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:12](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L12)_

_Defined in [packages/mesh-graphql-client/src/types.ts:298](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L298)_

---

## takerAssetAmount

• **takerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:15](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L15)_

_Defined in [packages/mesh-graphql-client/src/types.ts:302](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L302)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:13](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L13)_

_Defined in [packages/mesh-graphql-client/src/types.ts:308](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L308)_

---

## takerFee

• **takerFee**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:16](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L16)_

_Defined in [packages/mesh-graphql-client/src/types.ts:304](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L304)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:14](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L14)_

_Defined in [packages/mesh-graphql-client/src/types.ts:310](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L310)_

<hr />

# Interface: StringifiedSignedOrderV4

### Hierarchy

-   **StringifiedSignedOrderV4**

    ↳ [StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4)

### Properties

## chainId

• **chainId**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:315](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L315)_

---

## expiry

• **expiry**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:327](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L327)_

---

## feeRecipient

• **feeRecipient**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:325](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L325)_

---

## maker

• **maker**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:322](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L322)_

---

## makerAmount

• **makerAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:319](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L319)_

---

## makerToken

• **makerToken**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:317](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L317)_

---

## pool

• **pool**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:326](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L326)_

---

## salt

• **salt**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:328](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L328)_

---

## sender

• **sender**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:324](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L324)_

---

## signatureR

• **signatureR**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:330](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L330)_

---

## signatureS

• **signatureS**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:331](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L331)_

---

## signatureType

• **signatureType**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:329](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L329)_

---

## signatureV

• **signatureV**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:332](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L332)_

---

## taker

• **taker**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:323](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L323)_

---

## takerAmount

• **takerAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:320](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L320)_

---

## takerToken

• **takerToken**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:318](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L318)_

---

## takerTokenFeeAmount

• **takerTokenFeeAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:321](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L321)_

---

## verifyingContract

• **verifyingContract**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:316](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L316)_

<hr />

# Interface: StringifiedStats

### Hierarchy

-   **StringifiedStats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:291](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L291)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:290](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L290)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:279](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L279)_

---

## latestBlock

• **latestBlock**: _[StringifiedLatestBlock](#interface-stringifiedlatestblock)_

_Defined in [packages/mesh-graphql-client/src/types.ts:280](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L280)_

---

## maxExpirationTime

• **maxExpirationTime**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:288](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L288)_

---

## numOrders

• **numOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:282](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L282)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:284](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L284)_

---

## numOrdersIncludingRemovedV4

• **numOrdersIncludingRemovedV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:285](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L285)_

---

## numOrdersV4

• **numOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:283](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L283)_

---

## numPeers

• **numPeers**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:281](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L281)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:286](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L286)_

---

## numPinnedOrdersV4

• **numPinnedOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:287](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L287)_

---

## peerID

• **peerID**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:278](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L278)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:275](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L275)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:276](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L276)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:277](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L277)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:289](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L289)_

---

## version

• **version**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:274](https://github.com/0xProject/0x-mesh/blob/a5807182/packages/mesh-graphql-client/src/types.ts#L274)_

<hr />
