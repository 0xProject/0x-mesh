# Class: BrowserLink

### Hierarchy

-   ApolloLink

    ↳ **BrowserLink**

### Constructors

## constructer

\+ **new BrowserLink**(`_mesh`: Mesh): _[BrowserLink](#class-browserlink)_

_Overrides void_

_Defined in [packages/mesh-graphql-client/src/browser_link.ts:8](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/browser_link.ts#L8)_

**Parameters:**

| Name    | Type |
| ------- | ---- |
| `_mesh` | Mesh |

**Returns:** _[BrowserLink](#class-browserlink)_

### Methods

## concat

▸ **concat**(`next`: ApolloLink | RequestHandler): _ApolloLink_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:12

**Parameters:**

| Name   | Type                             |
| ------ | -------------------------------- |
| `next` | ApolloLink &#124; RequestHandler |

**Returns:** _ApolloLink_

---

## request

▸ **request**(`operation`: Operation): _Observable‹FetchResult›_

_Overrides void_

_Defined in [packages/mesh-graphql-client/src/browser_link.ts:13](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/browser_link.ts#L13)_

**Parameters:**

| Name        | Type      |
| ----------- | --------- |
| `operation` | Operation |

**Returns:** _Observable‹FetchResult›_

---

## setOnError

▸ **setOnError**(`fn`: ApolloLink["onError"]): _this_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:15

**Parameters:**

| Name | Type                  |
| ---- | --------------------- |
| `fn` | ApolloLink["onError"] |

**Returns:** _this_

---

## split

▸ **split**(`test`: function, `left`: ApolloLink | RequestHandler, `right?`: ApolloLink | RequestHandler): _ApolloLink_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:11

**Parameters:**

▪ **test**: _function_

▸ (`op`: Operation): _boolean_

**Parameters:**

| Name | Type      |
| ---- | --------- |
| `op` | Operation |

▪ **left**: _ApolloLink | RequestHandler_

▪`Optional` **right**: _ApolloLink | RequestHandler_

**Returns:** _ApolloLink_

---

## `Static` concat

▸ **concat**(`first`: ApolloLink | RequestHandler, `second`: ApolloLink | RequestHandler): _ApolloLink_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:9

**Parameters:**

| Name     | Type                             |
| -------- | -------------------------------- |
| `first`  | ApolloLink &#124; RequestHandler |
| `second` | ApolloLink &#124; RequestHandler |

**Returns:** _ApolloLink_

---

## `Static` empty

▸ **empty**(): _ApolloLink_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:5

**Returns:** _ApolloLink_

---

## `Static` execute

▸ **execute**(`link`: ApolloLink, `operation`: GraphQLRequest): _Observable‹FetchResult›_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:8

**Parameters:**

| Name        | Type           |
| ----------- | -------------- |
| `link`      | ApolloLink     |
| `operation` | GraphQLRequest |

**Returns:** _Observable‹FetchResult›_

---

## `Static` from

▸ **from**(`links`: ApolloLink‹› | function[]): _ApolloLink_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:6

**Parameters:**

| Name    | Type                           |
| ------- | ------------------------------ |
| `links` | ApolloLink‹› &#124; function[] |

**Returns:** _ApolloLink_

---

## `Static` split

▸ **split**(`test`: function, `left`: ApolloLink | RequestHandler, `right?`: ApolloLink | RequestHandler): _ApolloLink_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:7

**Parameters:**

▪ **test**: _function_

▸ (`op`: Operation): _boolean_

**Parameters:**

| Name | Type      |
| ---- | --------- |
| `op` | Operation |

▪ **left**: _ApolloLink | RequestHandler_

▪`Optional` **right**: _ApolloLink | RequestHandler_

**Returns:** _ApolloLink_

<hr />

# Class: MeshGraphQLClient

### Hierarchy

-   **MeshGraphQLClient**

### Constructors

## constructer

\+ **new MeshGraphQLClient**(`linkConfig`: [LinkConfig](#interface-linkconfig)): _[MeshGraphQLClient](#class-meshgraphqlclient)_

_Defined in [packages/mesh-graphql-client/src/index.ts:96](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L96)_

**Parameters:**

| Name         | Type                                |
| ------------ | ----------------------------------- |
| `linkConfig` | [LinkConfig](#interface-linkconfig) |

**Returns:** _[MeshGraphQLClient](#class-meshgraphqlclient)_

### Methods

## addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[], `pinned`: boolean, `opts?`: [AddOrdersOpts](#interface-addordersopts)): _Promise‹[AddOrdersResults](#interface-addordersresults)‹[OrderWithMetadata](#interface-orderwithmetadata), SignedOrder››_

_Defined in [packages/mesh-graphql-client/src/index.ts:193](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L193)_

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

_Defined in [packages/mesh-graphql-client/src/index.ts:222](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L222)_

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

_Defined in [packages/mesh-graphql-client/src/index.ts:285](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L285)_

**Parameters:**

| Name    | Type                                | Default                                                  |
| ------- | ----------------------------------- | -------------------------------------------------------- |
| `query` | [OrderQuery](#interface-orderquery) | { sort: [], filters: [], limit: defaultOrderQueryLimit } |

**Returns:** _Promise‹[OrderWithMetadata](#interface-orderwithmetadata)[]›_

---

## findOrdersV4Async

▸ **findOrdersV4Async**(`query`: [OrderQuery](#interface-orderquery)): _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4)[]›_

_Defined in [packages/mesh-graphql-client/src/index.ts:302](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L302)_

**Parameters:**

| Name    | Type                                | Default                                                  |
| ------- | ----------------------------------- | -------------------------------------------------------- |
| `query` | [OrderQuery](#interface-orderquery) | { sort: [], filters: [], limit: defaultOrderQueryLimit } |

**Returns:** _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4)[]›_

---

## getOrderAsync

▸ **getOrderAsync**(`hash`: string): _Promise‹[OrderWithMetadata](#interface-orderwithmetadata) | null›_

_Defined in [packages/mesh-graphql-client/src/index.ts:253](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L253)_

**Parameters:**

| Name   | Type   |
| ------ | ------ |
| `hash` | string |

**Returns:** _Promise‹[OrderWithMetadata](#interface-orderwithmetadata) | null›_

---

## getOrderV4Async

▸ **getOrderV4Async**(`hash`: string): _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4) | null›_

_Defined in [packages/mesh-graphql-client/src/index.ts:269](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L269)_

**Parameters:**

| Name   | Type   |
| ------ | ------ |
| `hash` | string |

**Returns:** _Promise‹[OrderWithMetadataV4](#interface-orderwithmetadatav4) | null›_

---

## getStatsAsync

▸ **getStatsAsync**(): _Promise‹[Stats](#interface-stats)›_

_Defined in [packages/mesh-graphql-client/src/index.ts:182](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L182)_

**Returns:** _Promise‹[Stats](#interface-stats)›_

---

## onOrderEvents

▸ **onOrderEvents**(): _Observable‹[OrderEvent](#interface-orderevent)[]›_

_Defined in [packages/mesh-graphql-client/src/index.ts:323](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L323)_

**Returns:** _Observable‹[OrderEvent](#interface-orderevent)[]›_

---

## onReconnected

▸ **onReconnected**(`cb`: function): _void_

_Defined in [packages/mesh-graphql-client/src/index.ts:319](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L319)_

**Parameters:**

▪ **cb**: _function_

▸ (): _void_

**Returns:** _void_

---

## rawQueryAsync

▸ **rawQueryAsync**<**T**, **TVariables**>(`options`: QueryOptions‹TVariables›): _Promise‹ApolloQueryResult‹T››_

_Defined in [packages/mesh-graphql-client/src/index.ts:370](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L370)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:176](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L176)_

---

## ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:178](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L178)_

---

## ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:177](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L177)_

---

## ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:172](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L172)_

---

## ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:171](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L171)_

---

## ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:174](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L174)_

---

## ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:175](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L175)_

---

## ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:173](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L173)_

---

## ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:180](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L180)_

---

## ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:181](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L181)_

---

## ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:179](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L179)_

---

## WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:182](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L182)_

---

## WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:183](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L183)_

<hr />

# Enumeration: FilterKind

### Enumeration members

## Equal

• **Equal**: = "EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:222](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L222)_

---

## Greater

• **Greater**: = "GREATER"

_Defined in [packages/mesh-graphql-client/src/types.ts:224](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L224)_

---

## GreaterOrEqual

• **GreaterOrEqual**: = "GREATER_OR_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:225](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L225)_

---

## Less

• **Less**: = "LESS"

_Defined in [packages/mesh-graphql-client/src/types.ts:226](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L226)_

---

## LessOrEqual

• **LessOrEqual**: = "LESS_OR_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:227](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L227)_

---

## NotEqual

• **NotEqual**: = "NOT_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:223](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L223)_

<hr />

# Enumeration: OrderEventEndState

### Enumeration members

## Added

• **Added**: = "ADDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:189](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L189)_

---

## Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:195](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L195)_

---

## Expired

• **Expired**: = "EXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:197](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L197)_

---

## FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

_Defined in [packages/mesh-graphql-client/src/types.ts:206](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L206)_

---

## Filled

• **Filled**: = "FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:191](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L191)_

---

## FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:193](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L193)_

---

## Invalid

• **Invalid**: = "INVALID"

_Defined in [packages/mesh-graphql-client/src/types.ts:199](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L199)_

---

## StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

_Defined in [packages/mesh-graphql-client/src/types.ts:211](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L211)_

---

## Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:201](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L201)_

---

## Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:203](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L203)_

<hr />

# Enumeration: RejectedOrderCode

### Enumeration members

## DatabaseFullOfOrders

• **DatabaseFullOfOrders**: = "DATABASE_FULL_OF_ORDERS"

_Defined in [packages/mesh-graphql-client/src/types.ts:146](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L146)_

---

## EthRpcRequestFailed

• **EthRpcRequestFailed**: = "ETH_RPC_REQUEST_FAILED"

_Defined in [packages/mesh-graphql-client/src/types.ts:127](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L127)_

---

## IncorrectExchangeAddress

• **IncorrectExchangeAddress**: = "INCORRECT_EXCHANGE_ADDRESS"

_Defined in [packages/mesh-graphql-client/src/types.ts:144](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L144)_

---

## InternalError

• **InternalError**: = "INTERNAL_ERROR"

_Defined in [packages/mesh-graphql-client/src/types.ts:140](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L140)_

---

## MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MAX_ORDER_SIZE_EXCEEDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:141](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L141)_

---

## OrderAlreadyStoredAndUnfillable

• **OrderAlreadyStoredAndUnfillable**: = "ORDER_ALREADY_STORED_AND_UNFILLABLE"

_Defined in [packages/mesh-graphql-client/src/types.ts:142](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L142)_

---

## OrderCancelled

• **OrderCancelled**: = "ORDER_CANCELLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:132](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L132)_

---

## OrderExpired

• **OrderExpired**: = "ORDER_EXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:130](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L130)_

---

## OrderForIncorrectChain

• **OrderForIncorrectChain**: = "ORDER_FOR_INCORRECT_CHAIN"

_Defined in [packages/mesh-graphql-client/src/types.ts:143](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L143)_

---

## OrderFullyFilled

• **OrderFullyFilled**: = "ORDER_FULLY_FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:131](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L131)_

---

## OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "ORDER_HAS_INVALID_MAKER_ASSET_AMOUNT"

_Defined in [packages/mesh-graphql-client/src/types.ts:128](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L128)_

---

## OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "ORDER_HAS_INVALID_MAKER_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:134](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L134)_

---

## OrderHasInvalidMakerFeeAssetData

• **OrderHasInvalidMakerFeeAssetData**: = "ORDER_HAS_INVALID_MAKER_FEE_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:135](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L135)_

---

## OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "ORDER_HAS_INVALID_SIGNATURE"

_Defined in [packages/mesh-graphql-client/src/types.ts:138](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L138)_

---

## OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "ORDER_HAS_INVALID_TAKER_ASSET_AMOUNT"

_Defined in [packages/mesh-graphql-client/src/types.ts:129](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L129)_

---

## OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "ORDER_HAS_INVALID_TAKER_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:136](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L136)_

---

## OrderHasInvalidTakerFeeAssetData

• **OrderHasInvalidTakerFeeAssetData**: = "ORDER_HAS_INVALID_TAKER_FEE_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:137](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L137)_

---

## OrderMaxExpirationExceeded

• **OrderMaxExpirationExceeded**: = "ORDER_MAX_EXPIRATION_EXCEEDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:139](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L139)_

---

## OrderUnfunded

• **OrderUnfunded**: = "ORDER_UNFUNDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:133](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L133)_

---

## SenderAddressNotAllowed

• **SenderAddressNotAllowed**: = "SENDER_ADDRESS_NOT_ALLOWED"

_Defined in [packages/mesh-graphql-client/src/types.ts:145](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L145)_

---

## TakerAddressNotAllowed

• **TakerAddressNotAllowed**: = "TAKER_ADDRESS_NOT_ALLOWED"

_Defined in [packages/mesh-graphql-client/src/types.ts:147](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L147)_

<hr />

# Enumeration: SortDirection

### Enumeration members

## Asc

• **Asc**: = "ASC"

_Defined in [packages/mesh-graphql-client/src/types.ts:217](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L217)_

---

## Desc

• **Desc**: = "DESC"

_Defined in [packages/mesh-graphql-client/src/types.ts:218](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L218)_

<hr />

# Interface: LinkConfig

### Hierarchy

-   **LinkConfig**

### Properties

## `Optional` httpUrl

• **httpUrl**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/index.ts:87](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L87)_

---

## `Optional` mesh

• **mesh**? : _Mesh_

_Defined in [packages/mesh-graphql-client/src/index.ts:89](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L89)_

---

## `Optional` webSocketUrl

• **webSocketUrl**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/index.ts:88](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/index.ts#L88)_

<hr />

# Interface: AcceptedOrderResult <**T**>

### Type parameters

▪ **T**: _[GenericOrderWithMetadata](#genericorderwithmetadata)_

### Hierarchy

-   **AcceptedOrderResult**

### Properties

## isNew

• **isNew**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:109](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L109)_

---

## order

• **order**: _T_

_Defined in [packages/mesh-graphql-client/src/types.ts:106](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L106)_

<hr />

# Interface: AddOrdersOpts

### Hierarchy

-   **AddOrdersOpts**

### Properties

## `Optional` keepCancelled

• **keepCancelled**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:6](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L6)_

---

## `Optional` keepExpired

• **keepExpired**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:7](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L7)_

---

## `Optional` keepFullyFilled

• **keepFullyFilled**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:8](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L8)_

---

## `Optional` keepUnfunded

• **keepUnfunded**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:9](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L9)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:25](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L25)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:32](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L32)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:97](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L97)_

---

## rejected

• **rejected**: _[RejectedOrderResult](#interface-rejectedorderresult)‹K›[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:100](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L100)_

<hr />

# Interface: ContractEvent

### Hierarchy

-   **ContractEvent**

### Properties

## address

• **address**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:164](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L164)_

---

## blockHash

• **blockHash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:159](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L159)_

---

## isRemoved

• **isRemoved**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:163](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L163)_

---

## kind

• **kind**: _[ContractEventKind](#enumeration-contracteventkind)_

_Defined in [packages/mesh-graphql-client/src/types.ts:165](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L165)_

---

## logIndex

• **logIndex**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:162](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L162)_

---

## parameters

• **parameters**: _any_

_Defined in [packages/mesh-graphql-client/src/types.ts:167](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L167)_

---

## txHash

• **txHash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:160](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L160)_

---

## txIndex

• **txIndex**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:161](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L161)_

<hr />

# Interface: LatestBlock

### Hierarchy

-   **LatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:78](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L78)_

---

## number

• **number**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:77](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L77)_

<hr />

# Interface: OrderEvent

### Hierarchy

-   **OrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:155](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L155)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [packages/mesh-graphql-client/src/types.ts:154](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L154)_

---

## `Optional` order

• **order**? : _[OrderWithMetadata](#interface-orderwithmetadata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:152](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L152)_

---

## `Optional` orderv4

• **orderv4**? : _[OrderWithMetadataV4](#interface-orderwithmetadatav4)_

_Defined in [packages/mesh-graphql-client/src/types.ts:153](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L153)_

---

## timestampMs

• **timestampMs**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:151](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L151)_

<hr />

# Interface: OrderEventResponse

### Hierarchy

-   **OrderEventResponse**

### Properties

## orderEvents

• **orderEvents**: _[StringifiedOrderEvent](#interface-stringifiedorderevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:52](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L52)_

<hr />

# Interface: OrderFilter

### Hierarchy

-   **OrderFilter**

### Properties

## field

• **field**: _[OrderField](#orderfield)_

_Defined in [packages/mesh-graphql-client/src/types.ts:236](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L236)_

---

## kind

• **kind**: _[FilterKind](#enumeration-filterkind)_

_Defined in [packages/mesh-graphql-client/src/types.ts:237](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L237)_

---

## value

• **value**: _OrderWithMetadata[OrderField]_

_Defined in [packages/mesh-graphql-client/src/types.ts:238](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L238)_

<hr />

# Interface: OrderQuery

### Hierarchy

-   **OrderQuery**

### Properties

## `Optional` filters

• **filters**? : _[OrderFilter](#interface-orderfilter)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:242](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L242)_

---

## `Optional` limit

• **limit**? : _undefined | number_

_Defined in [packages/mesh-graphql-client/src/types.ts:244](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L244)_

---

## `Optional` sort

• **sort**? : _[OrderSort](#interface-ordersort)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:243](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L243)_

<hr />

# Interface: OrderResponse

### Hierarchy

-   **OrderResponse**

### Properties

## order

• **order**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:36](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L36)_

<hr />

# Interface: OrderResponseV4

### Hierarchy

-   **OrderResponseV4**

### Properties

## orderv4

• **orderv4**: _[StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:40](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L40)_

<hr />

# Interface: OrderSort

### Hierarchy

-   **OrderSort**

### Properties

## direction

• **direction**: _[SortDirection](#enumeration-sortdirection)_

_Defined in [packages/mesh-graphql-client/src/types.ts:232](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L232)_

---

## field

• **field**: _[OrderField](#orderfield)_

_Defined in [packages/mesh-graphql-client/src/types.ts:231](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L231)_

<hr />

# Interface: OrdersResponse

### Hierarchy

-   **OrdersResponse**

### Properties

## orders

• **orders**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:44](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L44)_

<hr />

# Interface: OrdersResponseV4

### Hierarchy

-   **OrdersResponseV4**

### Properties

## ordersv4

• **ordersv4**: _[StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:48](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L48)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:83](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L83)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:82](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L82)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:91](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L91)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:90](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L90)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:87](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L87)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:120](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L120)_

---

## `Optional` hash

• **hash**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/types.ts:114](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L114)_

---

## message

• **message**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:123](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L123)_

---

## order

• **order**: _K_

_Defined in [packages/mesh-graphql-client/src/types.ts:117](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L117)_

<hr />

# Interface: Stats

### Hierarchy

-   **Stats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:73](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L73)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:72](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L72)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:61](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L61)_

---

## latestBlock

• **latestBlock**: _[LatestBlock](#interface-latestblock)_

_Defined in [packages/mesh-graphql-client/src/types.ts:62](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L62)_

---

## maxExpirationTime

• **maxExpirationTime**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:70](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L70)_

---

## numOrders

• **numOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:64](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L64)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:66](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L66)_

---

## numOrdersIncludingRemovedV4

• **numOrdersIncludingRemovedV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:67](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L67)_

---

## numOrdersV4

• **numOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:65](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L65)_

---

## numPeers

• **numPeers**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:63](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L63)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:68](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L68)_

---

## numPinnedOrdersV4

• **numPinnedOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:69](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L69)_

---

## peerID

• **peerID**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:60](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L60)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:57](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L57)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:58](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L58)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:59](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L59)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _Date_

_Defined in [packages/mesh-graphql-client/src/types.ts:71](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L71)_

---

## version

• **version**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:56](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L56)_

<hr />

# Interface: StatsResponse

### Hierarchy

-   **StatsResponse**

### Properties

## stats

• **stats**: _[StringifiedStats](#interface-stringifiedstats)_

_Defined in [packages/mesh-graphql-client/src/types.ts:13](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L13)_

<hr />

# Interface: StringifiedAcceptedOrderResult <**T**>

### Type parameters

▪ **T**: _[GenericStringifiedOrderWithMetadata](#genericstringifiedorderwithmetadata)_

### Hierarchy

-   **StringifiedAcceptedOrderResult**

### Properties

## isNew

• **isNew**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:333](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L333)_

---

## order

• **order**: _T_

_Defined in [packages/mesh-graphql-client/src/types.ts:332](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L332)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:327](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L327)_

---

## rejected

• **rejected**: _[StringifiedRejectedOrderResult](#interface-stringifiedrejectedorderresult)‹K›[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:328](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L328)_

<hr />

# Interface: StringifiedLatestBlock

### Hierarchy

-   **StringifiedLatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:249](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L249)_

---

## number

• **number**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:248](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L248)_

<hr />

# Interface: StringifiedOrderEvent

### Hierarchy

-   **StringifiedOrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:349](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L349)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [packages/mesh-graphql-client/src/types.ts:347](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L347)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:348](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L348)_

---

## order

• **order**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:345](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L345)_

---

## orderv4

• **orderv4**: _[StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:346](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L346)_

---

## timestamp

• **timestamp**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:344](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L344)_

<hr />

# Interface: StringifiedOrderWithMetadata

### Hierarchy

-   [StringifiedSignedOrder](#interface-stringifiedsignedorder)

    ↳ **StringifiedOrderWithMetadata**

### Properties

## chainId

• **chainId**: _string_

_Inherited from [StringifiedSignedOrder](#chainid)_

_Defined in [packages/mesh-graphql-client/src/types.ts:274](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L274)_

---

## exchangeAddress

• **exchangeAddress**: _string_

_Inherited from [StringifiedSignedOrder](#exchangeaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:275](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L275)_

---

## expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Inherited from [StringifiedSignedOrder](#expirationtimeseconds)_

_Defined in [packages/mesh-graphql-client/src/types.ts:284](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L284)_

---

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Inherited from [StringifiedSignedOrder](#feerecipientaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:278](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L278)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:320](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L320)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:319](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L319)_

---

## makerAddress

• **makerAddress**: _string_

_Inherited from [StringifiedSignedOrder](#makeraddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:276](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L276)_

---

## makerAssetAmount

• **makerAssetAmount**: _string_

_Inherited from [StringifiedSignedOrder](#makerassetamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:280](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L280)_

---

## makerAssetData

• **makerAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#makerassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:286](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L286)_

---

## makerFee

• **makerFee**: _string_

_Inherited from [StringifiedSignedOrder](#makerfee)_

_Defined in [packages/mesh-graphql-client/src/types.ts:282](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L282)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#makerfeeassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:288](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L288)_

---

## salt

• **salt**: _string_

_Inherited from [StringifiedSignedOrder](#salt)_

_Defined in [packages/mesh-graphql-client/src/types.ts:285](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L285)_

---

## senderAddress

• **senderAddress**: _string_

_Inherited from [StringifiedSignedOrder](#senderaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:279](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L279)_

---

## signature

• **signature**: _string_

_Inherited from [StringifiedSignedOrder](#signature)_

_Defined in [packages/mesh-graphql-client/src/types.ts:290](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L290)_

---

## takerAddress

• **takerAddress**: _string_

_Inherited from [StringifiedSignedOrder](#takeraddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:277](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L277)_

---

## takerAssetAmount

• **takerAssetAmount**: _string_

_Inherited from [StringifiedSignedOrder](#takerassetamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:281](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L281)_

---

## takerAssetData

• **takerAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#takerassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:287](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L287)_

---

## takerFee

• **takerFee**: _string_

_Inherited from [StringifiedSignedOrder](#takerfee)_

_Defined in [packages/mesh-graphql-client/src/types.ts:283](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L283)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#takerfeeassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:289](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L289)_

<hr />

# Interface: StringifiedOrderWithMetadataV4

### Hierarchy

-   [StringifiedSignedOrderV4](#interface-stringifiedsignedorderv4)

    ↳ **StringifiedOrderWithMetadataV4**

### Properties

## chainId

• **chainId**: _string_

_Inherited from [StringifiedSignedOrderV4](#chainid)_

_Defined in [packages/mesh-graphql-client/src/types.ts:294](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L294)_

---

## expiry

• **expiry**: _string_

_Inherited from [StringifiedSignedOrderV4](#expiry)_

_Defined in [packages/mesh-graphql-client/src/types.ts:306](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L306)_

---

## feeRecipient

• **feeRecipient**: _string_

_Inherited from [StringifiedSignedOrderV4](#feerecipient)_

_Defined in [packages/mesh-graphql-client/src/types.ts:304](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L304)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:315](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L315)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:314](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L314)_

---

## maker

• **maker**: _string_

_Inherited from [StringifiedSignedOrderV4](#maker)_

_Defined in [packages/mesh-graphql-client/src/types.ts:301](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L301)_

---

## makerAmount

• **makerAmount**: _string_

_Inherited from [StringifiedSignedOrderV4](#makeramount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:298](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L298)_

---

## makerToken

• **makerToken**: _string_

_Inherited from [StringifiedSignedOrderV4](#makertoken)_

_Defined in [packages/mesh-graphql-client/src/types.ts:296](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L296)_

---

## pool

• **pool**: _string_

_Inherited from [StringifiedSignedOrderV4](#pool)_

_Defined in [packages/mesh-graphql-client/src/types.ts:305](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L305)_

---

## salt

• **salt**: _string_

_Inherited from [StringifiedSignedOrderV4](#salt)_

_Defined in [packages/mesh-graphql-client/src/types.ts:307](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L307)_

---

## sender

• **sender**: _string_

_Inherited from [StringifiedSignedOrderV4](#sender)_

_Defined in [packages/mesh-graphql-client/src/types.ts:303](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L303)_

---

## signatureR

• **signatureR**: _string_

_Inherited from [StringifiedSignedOrderV4](#signaturer)_

_Defined in [packages/mesh-graphql-client/src/types.ts:309](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L309)_

---

## signatureS

• **signatureS**: _string_

_Inherited from [StringifiedSignedOrderV4](#signatures)_

_Defined in [packages/mesh-graphql-client/src/types.ts:310](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L310)_

---

## signatureType

• **signatureType**: _string_

_Inherited from [StringifiedSignedOrderV4](#signaturetype)_

_Defined in [packages/mesh-graphql-client/src/types.ts:308](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L308)_

---

## signatureV

• **signatureV**: _string_

_Inherited from [StringifiedSignedOrderV4](#signaturev)_

_Defined in [packages/mesh-graphql-client/src/types.ts:311](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L311)_

---

## taker

• **taker**: _string_

_Inherited from [StringifiedSignedOrderV4](#taker)_

_Defined in [packages/mesh-graphql-client/src/types.ts:302](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L302)_

---

## takerAmount

• **takerAmount**: _string_

_Inherited from [StringifiedSignedOrderV4](#takeramount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:299](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L299)_

---

## takerToken

• **takerToken**: _string_

_Inherited from [StringifiedSignedOrderV4](#takertoken)_

_Defined in [packages/mesh-graphql-client/src/types.ts:297](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L297)_

---

## takerTokenFeeAmount

• **takerTokenFeeAmount**: _string_

_Inherited from [StringifiedSignedOrderV4](#takertokenfeeamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:300](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L300)_

---

## verifyingContract

• **verifyingContract**: _string_

_Inherited from [StringifiedSignedOrderV4](#verifyingcontract)_

_Defined in [packages/mesh-graphql-client/src/types.ts:295](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L295)_

<hr />

# Interface: StringifiedRejectedOrderResult <**K**>

### Type parameters

▪ **K**: _[GenericStringifiedSignedOrders](#genericstringifiedsignedorders)_

### Hierarchy

-   **StringifiedRejectedOrderResult**

### Properties

## code

• **code**: _[RejectedOrderCode](#enumeration-rejectedordercode)_

_Defined in [packages/mesh-graphql-client/src/types.ts:339](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L339)_

---

## `Optional` hash

• **hash**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/types.ts:337](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L337)_

---

## message

• **message**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:340](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L340)_

---

## order

• **order**: _K_

_Defined in [packages/mesh-graphql-client/src/types.ts:338](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L338)_

<hr />

# Interface: StringifiedSignedOrder

### Hierarchy

-   **StringifiedSignedOrder**

    ↳ [StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)

### Properties

## chainId

• **chainId**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:274](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L274)_

---

## exchangeAddress

• **exchangeAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:275](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L275)_

---

## expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:284](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L284)_

---

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:278](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L278)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:276](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L276)_

---

## makerAssetAmount

• **makerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:280](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L280)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:286](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L286)_

---

## makerFee

• **makerFee**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:282](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L282)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:288](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L288)_

---

## salt

• **salt**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:285](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L285)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:279](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L279)_

---

## signature

• **signature**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:290](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L290)_

---

## takerAddress

• **takerAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:277](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L277)_

---

## takerAssetAmount

• **takerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:281](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L281)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:287](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L287)_

---

## takerFee

• **takerFee**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:283](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L283)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:289](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L289)_

<hr />

# Interface: StringifiedSignedOrderV4

### Hierarchy

-   **StringifiedSignedOrderV4**

    ↳ [StringifiedOrderWithMetadataV4](#interface-stringifiedorderwithmetadatav4)

### Properties

## chainId

• **chainId**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:294](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L294)_

---

## expiry

• **expiry**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:306](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L306)_

---

## feeRecipient

• **feeRecipient**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:304](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L304)_

---

## maker

• **maker**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:301](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L301)_

---

## makerAmount

• **makerAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:298](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L298)_

---

## makerToken

• **makerToken**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:296](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L296)_

---

## pool

• **pool**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:305](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L305)_

---

## salt

• **salt**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:307](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L307)_

---

## sender

• **sender**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:303](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L303)_

---

## signatureR

• **signatureR**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:309](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L309)_

---

## signatureS

• **signatureS**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:310](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L310)_

---

## signatureType

• **signatureType**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:308](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L308)_

---

## signatureV

• **signatureV**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:311](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L311)_

---

## taker

• **taker**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:302](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L302)_

---

## takerAmount

• **takerAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:299](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L299)_

---

## takerToken

• **takerToken**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:297](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L297)_

---

## takerTokenFeeAmount

• **takerTokenFeeAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:300](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L300)_

---

## verifyingContract

• **verifyingContract**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:295](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L295)_

<hr />

# Interface: StringifiedStats

### Hierarchy

-   **StringifiedStats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:270](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L270)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:269](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L269)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:258](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L258)_

---

## latestBlock

• **latestBlock**: _[StringifiedLatestBlock](#interface-stringifiedlatestblock)_

_Defined in [packages/mesh-graphql-client/src/types.ts:259](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L259)_

---

## maxExpirationTime

• **maxExpirationTime**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:267](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L267)_

---

## numOrders

• **numOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:261](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L261)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:263](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L263)_

---

## numOrdersIncludingRemovedV4

• **numOrdersIncludingRemovedV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:264](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L264)_

---

## numOrdersV4

• **numOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:262](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L262)_

---

## numPeers

• **numPeers**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:260](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L260)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:265](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L265)_

---

## numPinnedOrdersV4

• **numPinnedOrdersV4**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:266](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L266)_

---

## peerID

• **peerID**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:257](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L257)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:254](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L254)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:255](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L255)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:256](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L256)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:268](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L268)_

---

## version

• **version**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:253](https://github.com/0xProject/0x-mesh/blob/161799fc/packages/mesh-graphql-client/src/types.ts#L253)_

<hr />
