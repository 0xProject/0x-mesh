# Class: BrowserLink

### Hierarchy

-   ApolloLink

    ↳ **BrowserLink**

### Constructors

## constructer

\+ **new BrowserLink**(`_mesh`: Mesh): _[BrowserLink](#class-browserlink)_

_Overrides void_

_Defined in [packages/mesh-graphql-client/src/browser_link.ts:7](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/browser_link.ts#L7)_

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

_Defined in [packages/mesh-graphql-client/src/browser_link.ts:12](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/browser_link.ts#L12)_

**Parameters:**

| Name        | Type      |
| ----------- | --------- |
| `operation` | Operation |

**Returns:** _Observable‹FetchResult›_

---

## setOnError

▸ **setOnError**(`fn`: function): _this_

Defined in node_modules/@apollo/client/link/core/ApolloLink.d.ts:15

**Parameters:**

▪ **fn**: _function_

▸ (`reason`: any): _any_

**Parameters:**

| Name     | Type |
| -------- | ---- |
| `reason` | any  |

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

_Defined in [packages/mesh-graphql-client/src/index.ts:250](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L250)_

**Parameters:**

| Name         | Type                                |
| ------------ | ----------------------------------- |
| `linkConfig` | [LinkConfig](#interface-linkconfig) |

**Returns:** _[MeshGraphQLClient](#class-meshgraphqlclient)_

### Methods

## addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[], `pinned`: boolean, `opts?`: [AddOrdersOpts](#interface-addordersopts)): _Promise‹[AddOrdersResults](#interface-addordersresults)›_

_Defined in [packages/mesh-graphql-client/src/index.ts:335](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L335)_

**Parameters:**

| Name     | Type                                      | Default |
| -------- | ----------------------------------------- | ------- |
| `orders` | SignedOrder[]                             | -       |
| `pinned` | boolean                                   | true    |
| `opts?`  | [AddOrdersOpts](#interface-addordersopts) | -       |

**Returns:** _Promise‹[AddOrdersResults](#interface-addordersresults)›_

---

## findOrdersAsync

▸ **findOrdersAsync**(`query`: [OrderQuery](#interface-orderquery)): _Promise‹[OrderWithMetadata](#interface-orderwithmetadata)[]›_

_Defined in [packages/mesh-graphql-client/src/index.ts:377](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L377)_

**Parameters:**

| Name    | Type                                | Default                                                  |
| ------- | ----------------------------------- | -------------------------------------------------------- |
| `query` | [OrderQuery](#interface-orderquery) | { sort: [], filters: [], limit: defaultOrderQueryLimit } |

**Returns:** _Promise‹[OrderWithMetadata](#interface-orderwithmetadata)[]›_

---

## getOrderAsync

▸ **getOrderAsync**(`hash`: string): _Promise‹[OrderWithMetadata](#interface-orderwithmetadata) | null›_

_Defined in [packages/mesh-graphql-client/src/index.ts:361](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L361)_

**Parameters:**

| Name   | Type   |
| ------ | ------ |
| `hash` | string |

**Returns:** _Promise‹[OrderWithMetadata](#interface-orderwithmetadata) | null›_

---

## getStatsAsync

▸ **getStatsAsync**(): _Promise‹[Stats](#interface-stats)›_

_Defined in [packages/mesh-graphql-client/src/index.ts:324](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L324)_

**Returns:** _Promise‹[Stats](#interface-stats)›_

---

## onOrderEvents

▸ **onOrderEvents**(): _Observable‹[OrderEvent](#interface-orderevent)[]›_

_Defined in [packages/mesh-graphql-client/src/index.ts:394](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L394)_

**Returns:** _Observable‹[OrderEvent](#interface-orderevent)[]›_

---

## rawQueryAsync

▸ **rawQueryAsync**<**T**, **TVariables**>(`options`: QueryOptions‹TVariables›): _Promise‹ApolloQueryResult‹T››_

_Defined in [packages/mesh-graphql-client/src/index.ts:440](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L440)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:137](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L137)_

---

## ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:139](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L139)_

---

## ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:138](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L138)_

---

## ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:133](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L133)_

---

## ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:132](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L132)_

---

## ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:135](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L135)_

---

## ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:136](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L136)_

---

## ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:134](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L134)_

---

## ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:141](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L141)_

---

## ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:142](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L142)_

---

## ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:140](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L140)_

---

## WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:143](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L143)_

---

## WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in [packages/mesh-graphql-client/src/types.ts:144](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L144)_

<hr />

# Enumeration: FilterKind

### Enumeration members

## Equal

• **Equal**: = "EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:181](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L181)_

---

## Greater

• **Greater**: = "GREATER"

_Defined in [packages/mesh-graphql-client/src/types.ts:183](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L183)_

---

## GreaterOrEqual

• **GreaterOrEqual**: = "GREATER_OR_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:184](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L184)_

---

## Less

• **Less**: = "LESS"

_Defined in [packages/mesh-graphql-client/src/types.ts:185](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L185)_

---

## LessOrEqual

• **LessOrEqual**: = "LESS_OR_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:186](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L186)_

---

## NotEqual

• **NotEqual**: = "NOT_EQUAL"

_Defined in [packages/mesh-graphql-client/src/types.ts:182](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L182)_

<hr />

# Enumeration: OrderEventEndState

### Enumeration members

## Added

• **Added**: = "ADDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:150](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L150)_

---

## Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:156](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L156)_

---

## Expired

• **Expired**: = "EXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:158](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L158)_

---

## FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

_Defined in [packages/mesh-graphql-client/src/types.ts:165](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L165)_

---

## Filled

• **Filled**: = "FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:152](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L152)_

---

## FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:154](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L154)_

---

## StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

_Defined in [packages/mesh-graphql-client/src/types.ts:170](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L170)_

---

## Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:160](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L160)_

---

## Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:162](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L162)_

<hr />

# Enumeration: RejectedOrderCode

### Enumeration members

## DatabaseFullOfOrders

• **DatabaseFullOfOrders**: = "DATABASE_FULL_OF_ORDERS"

_Defined in [packages/mesh-graphql-client/src/types.ts:109](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L109)_

---

## EthRpcRequestFailed

• **EthRpcRequestFailed**: = "ETH_RPC_REQUEST_FAILED"

_Defined in [packages/mesh-graphql-client/src/types.ts:90](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L90)_

---

## IncorrectExchangeAddress

• **IncorrectExchangeAddress**: = "INCORRECT_EXCHANGE_ADDRESS"

_Defined in [packages/mesh-graphql-client/src/types.ts:107](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L107)_

---

## InternalError

• **InternalError**: = "INTERNAL_ERROR"

_Defined in [packages/mesh-graphql-client/src/types.ts:103](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L103)_

---

## MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MAX_ORDER_SIZE_EXCEEDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:104](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L104)_

---

## OrderAlreadyStoredAndUnfillable

• **OrderAlreadyStoredAndUnfillable**: = "ORDER_ALREADY_STORED_AND_UNFILLABLE"

_Defined in [packages/mesh-graphql-client/src/types.ts:105](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L105)_

---

## OrderCancelled

• **OrderCancelled**: = "ORDER_CANCELLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:95](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L95)_

---

## OrderExpired

• **OrderExpired**: = "ORDER_EXPIRED"

_Defined in [packages/mesh-graphql-client/src/types.ts:93](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L93)_

---

## OrderForIncorrectChain

• **OrderForIncorrectChain**: = "ORDER_FOR_INCORRECT_CHAIN"

_Defined in [packages/mesh-graphql-client/src/types.ts:106](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L106)_

---

## OrderFullyFilled

• **OrderFullyFilled**: = "ORDER_FULLY_FILLED"

_Defined in [packages/mesh-graphql-client/src/types.ts:94](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L94)_

---

## OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "ORDER_HAS_INVALID_MAKER_ASSET_AMOUNT"

_Defined in [packages/mesh-graphql-client/src/types.ts:91](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L91)_

---

## OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "ORDER_HAS_INVALID_MAKER_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:97](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L97)_

---

## OrderHasInvalidMakerFeeAssetData

• **OrderHasInvalidMakerFeeAssetData**: = "ORDER_HAS_INVALID_MAKER_FEE_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:98](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L98)_

---

## OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "ORDER_HAS_INVALID_SIGNATURE"

_Defined in [packages/mesh-graphql-client/src/types.ts:101](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L101)_

---

## OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "ORDER_HAS_INVALID_TAKER_ASSET_AMOUNT"

_Defined in [packages/mesh-graphql-client/src/types.ts:92](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L92)_

---

## OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "ORDER_HAS_INVALID_TAKER_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:99](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L99)_

---

## OrderHasInvalidTakerFeeAssetData

• **OrderHasInvalidTakerFeeAssetData**: = "ORDER_HAS_INVALID_TAKER_FEE_ASSET_DATA"

_Defined in [packages/mesh-graphql-client/src/types.ts:100](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L100)_

---

## OrderMaxExpirationExceeded

• **OrderMaxExpirationExceeded**: = "ORDER_MAX_EXPIRATION_EXCEEDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:102](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L102)_

---

## OrderUnfunded

• **OrderUnfunded**: = "ORDER_UNFUNDED"

_Defined in [packages/mesh-graphql-client/src/types.ts:96](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L96)_

---

## SenderAddressNotAllowed

• **SenderAddressNotAllowed**: = "SENDER_ADDRESS_NOT_ALLOWED"

_Defined in [packages/mesh-graphql-client/src/types.ts:108](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L108)_

<hr />

# Enumeration: SortDirection

### Enumeration members

## Asc

• **Asc**: = "ASC"

_Defined in [packages/mesh-graphql-client/src/types.ts:176](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L176)_

---

## Desc

• **Desc**: = "DESC"

_Defined in [packages/mesh-graphql-client/src/types.ts:177](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L177)_

<hr />

# Interface: LinkConfig

### Hierarchy

-   **LinkConfig**

### Properties

## `Optional` httpUrl

• **httpUrl**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/index.ts:242](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L242)_

---

## `Optional` mesh

• **mesh**? : _Mesh_

_Defined in [packages/mesh-graphql-client/src/index.ts:244](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L244)_

---

## `Optional` webSocketUrl

• **webSocketUrl**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/index.ts:243](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/index.ts#L243)_

<hr />

# Interface: AcceptedOrderResult

### Hierarchy

-   **AcceptedOrderResult**

### Properties

## isNew

• **isNew**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:73](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L73)_

---

## order

• **order**: _[OrderWithMetadata](#interface-orderwithmetadata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:70](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L70)_

<hr />

# Interface: AddOrdersOpts

### Hierarchy

-   **AddOrdersOpts**

### Properties

## `Optional` keepCancelled

• **keepCancelled**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:5](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L5)_

---

## `Optional` keepExpired

• **keepExpired**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:6](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L6)_

---

## `Optional` keepFullyFilled

• **keepFullyFilled**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:7](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L7)_

---

## `Optional` keepUnfunded

• **keepUnfunded**? : _undefined | false | true_

_Defined in [packages/mesh-graphql-client/src/types.ts:8](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L8)_

<hr />

# Interface: AddOrdersResponse

### Hierarchy

-   **AddOrdersResponse**

### Properties

## addOrders

• **addOrders**: _[StringifiedAddOrdersResults](#interface-stringifiedaddordersresults)_

_Defined in [packages/mesh-graphql-client/src/types.ts:16](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L16)_

<hr />

# Interface: AddOrdersResults

### Hierarchy

-   **AddOrdersResults**

### Properties

## accepted

• **accepted**: _[AcceptedOrderResult](#interface-acceptedorderresult)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:62](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L62)_

---

## rejected

• **rejected**: _[RejectedOrderResult](#interface-rejectedorderresult)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:65](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L65)_

<hr />

# Interface: ContractEvent

### Hierarchy

-   **ContractEvent**

### Properties

## address

• **address**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:125](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L125)_

---

## blockHash

• **blockHash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:120](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L120)_

---

## isRemoved

• **isRemoved**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:124](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L124)_

---

## kind

• **kind**: _[ContractEventKind](#enumeration-contracteventkind)_

_Defined in [packages/mesh-graphql-client/src/types.ts:126](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L126)_

---

## logIndex

• **logIndex**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:123](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L123)_

---

## parameters

• **parameters**: _any_

_Defined in [packages/mesh-graphql-client/src/types.ts:128](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L128)_

---

## txHash

• **txHash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:121](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L121)_

---

## txIndex

• **txIndex**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:122](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L122)_

<hr />

# Interface: LatestBlock

### Hierarchy

-   **LatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:51](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L51)_

---

## number

• **number**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:50](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L50)_

<hr />

# Interface: OrderEvent

### Hierarchy

-   **OrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:116](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L116)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [packages/mesh-graphql-client/src/types.ts:115](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L115)_

---

## order

• **order**: _[OrderWithMetadata](#interface-orderwithmetadata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:114](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L114)_

---

## timestampMs

• **timestampMs**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:113](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L113)_

<hr />

# Interface: OrderEventResponse

### Hierarchy

-   **OrderEventResponse**

### Properties

## orderEvents

• **orderEvents**: _[StringifiedOrderEvent](#interface-stringifiedorderevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:28](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L28)_

<hr />

# Interface: OrderFilter

### Hierarchy

-   **OrderFilter**

### Properties

## field

• **field**: _[OrderField](#orderfield)_

_Defined in [packages/mesh-graphql-client/src/types.ts:195](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L195)_

---

## kind

• **kind**: _[FilterKind](#enumeration-filterkind)_

_Defined in [packages/mesh-graphql-client/src/types.ts:196](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L196)_

---

## value

• **value**: _OrderWithMetadata[OrderField]_

_Defined in [packages/mesh-graphql-client/src/types.ts:197](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L197)_

<hr />

# Interface: OrderQuery

### Hierarchy

-   **OrderQuery**

### Properties

## `Optional` filters

• **filters**? : _[OrderFilter](#interface-orderfilter)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:201](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L201)_

---

## `Optional` limit

• **limit**? : _undefined | number_

_Defined in [packages/mesh-graphql-client/src/types.ts:203](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L203)_

---

## `Optional` sort

• **sort**? : _[OrderSort](#interface-ordersort)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:202](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L202)_

<hr />

# Interface: OrderResponse

### Hierarchy

-   **OrderResponse**

### Properties

## order

• **order**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata) | null_

_Defined in [packages/mesh-graphql-client/src/types.ts:20](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L20)_

<hr />

# Interface: OrderSort

### Hierarchy

-   **OrderSort**

### Properties

## direction

• **direction**: _[SortDirection](#enumeration-sortdirection)_

_Defined in [packages/mesh-graphql-client/src/types.ts:191](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L191)_

---

## field

• **field**: _[OrderField](#orderfield)_

_Defined in [packages/mesh-graphql-client/src/types.ts:190](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L190)_

<hr />

# Interface: OrdersResponse

### Hierarchy

-   **OrdersResponse**

### Properties

## orders

• **orders**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:24](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L24)_

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

_Defined in [packages/mesh-graphql-client/src/types.ts:56](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L56)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:55](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L55)_

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

# Interface: RejectedOrderResult

### Hierarchy

-   **RejectedOrderResult**

### Properties

## code

• **code**: _[RejectedOrderCode](#enumeration-rejectedordercode)_

_Defined in [packages/mesh-graphql-client/src/types.ts:83](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L83)_

---

## `Optional` hash

• **hash**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/types.ts:78](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L78)_

---

## message

• **message**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:86](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L86)_

---

## order

• **order**: _SignedOrder_

_Defined in [packages/mesh-graphql-client/src/types.ts:80](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L80)_

<hr />

# Interface: Stats

### Hierarchy

-   **Stats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:46](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L46)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:45](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L45)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:37](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L37)_

---

## latestBlock

• **latestBlock**: _[LatestBlock](#interface-latestblock)_

_Defined in [packages/mesh-graphql-client/src/types.ts:38](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L38)_

---

## maxExpirationTime

• **maxExpirationTime**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:43](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L43)_

---

## numOrders

• **numOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:40](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L40)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:41](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L41)_

---

## numPeers

• **numPeers**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:39](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L39)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:42](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L42)_

---

## peerID

• **peerID**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:36](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L36)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:33](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L33)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:34](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L34)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:35](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L35)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _Date_

_Defined in [packages/mesh-graphql-client/src/types.ts:44](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L44)_

---

## version

• **version**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:32](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L32)_

<hr />

# Interface: StatsResponse

### Hierarchy

-   **StatsResponse**

### Properties

## stats

• **stats**: _[StringifiedStats](#interface-stringifiedstats)_

_Defined in [packages/mesh-graphql-client/src/types.ts:12](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L12)_

<hr />

# Interface: StringifiedAcceptedOrderResult

### Hierarchy

-   **StringifiedAcceptedOrderResult**

### Properties

## isNew

• **isNew**: _boolean_

_Defined in [packages/mesh-graphql-client/src/types.ts:261](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L261)_

---

## order

• **order**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:260](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L260)_

<hr />

# Interface: StringifiedAddOrdersResults

### Hierarchy

-   **StringifiedAddOrdersResults**

### Properties

## accepted

• **accepted**: _[StringifiedAcceptedOrderResult](#interface-stringifiedacceptedorderresult)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:255](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L255)_

---

## rejected

• **rejected**: _[StringifiedRejectedOrderResult](#interface-stringifiedrejectedorderresult)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:256](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L256)_

<hr />

# Interface: StringifiedLatestBlock

### Hierarchy

-   **StringifiedLatestBlock**

### Properties

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:208](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L208)_

---

## number

• **number**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:207](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L207)_

<hr />

# Interface: StringifiedOrderEvent

### Hierarchy

-   **StringifiedOrderEvent**

### Properties

## contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:276](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L276)_

---

## endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [packages/mesh-graphql-client/src/types.ts:274](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L274)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [packages/mesh-graphql-client/src/types.ts:275](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L275)_

---

## order

• **order**: _[StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:273](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L273)_

---

## timestamp

• **timestamp**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:272](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L272)_

<hr />

# Interface: StringifiedOrderWithMetadata

### Hierarchy

-   [StringifiedSignedOrder](#interface-stringifiedsignedorder)

    ↳ **StringifiedOrderWithMetadata**

### Properties

## chainId

• **chainId**: _string_

_Inherited from [StringifiedSignedOrder](#chainid)_

_Defined in [packages/mesh-graphql-client/src/types.ts:230](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L230)_

---

## exchangeAddress

• **exchangeAddress**: _string_

_Inherited from [StringifiedSignedOrder](#exchangeaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:231](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L231)_

---

## expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Inherited from [StringifiedSignedOrder](#expirationtimeseconds)_

_Defined in [packages/mesh-graphql-client/src/types.ts:240](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L240)_

---

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Inherited from [StringifiedSignedOrder](#feerecipientaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:234](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L234)_

---

## fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:251](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L251)_

---

## hash

• **hash**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:250](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L250)_

---

## makerAddress

• **makerAddress**: _string_

_Inherited from [StringifiedSignedOrder](#makeraddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:232](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L232)_

---

## makerAssetAmount

• **makerAssetAmount**: _string_

_Inherited from [StringifiedSignedOrder](#makerassetamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:236](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L236)_

---

## makerAssetData

• **makerAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#makerassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:242](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L242)_

---

## makerFee

• **makerFee**: _string_

_Inherited from [StringifiedSignedOrder](#makerfee)_

_Defined in [packages/mesh-graphql-client/src/types.ts:238](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L238)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#makerfeeassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:244](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L244)_

---

## salt

• **salt**: _string_

_Inherited from [StringifiedSignedOrder](#salt)_

_Defined in [packages/mesh-graphql-client/src/types.ts:241](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L241)_

---

## senderAddress

• **senderAddress**: _string_

_Inherited from [StringifiedSignedOrder](#senderaddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:235](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L235)_

---

## signature

• **signature**: _string_

_Inherited from [StringifiedSignedOrder](#signature)_

_Defined in [packages/mesh-graphql-client/src/types.ts:246](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L246)_

---

## takerAddress

• **takerAddress**: _string_

_Inherited from [StringifiedSignedOrder](#takeraddress)_

_Defined in [packages/mesh-graphql-client/src/types.ts:233](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L233)_

---

## takerAssetAmount

• **takerAssetAmount**: _string_

_Inherited from [StringifiedSignedOrder](#takerassetamount)_

_Defined in [packages/mesh-graphql-client/src/types.ts:237](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L237)_

---

## takerAssetData

• **takerAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#takerassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:243](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L243)_

---

## takerFee

• **takerFee**: _string_

_Inherited from [StringifiedSignedOrder](#takerfee)_

_Defined in [packages/mesh-graphql-client/src/types.ts:239](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L239)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Inherited from [StringifiedSignedOrder](#takerfeeassetdata)_

_Defined in [packages/mesh-graphql-client/src/types.ts:245](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L245)_

<hr />

# Interface: StringifiedRejectedOrderResult

### Hierarchy

-   **StringifiedRejectedOrderResult**

### Properties

## code

• **code**: _[RejectedOrderCode](#enumeration-rejectedordercode)_

_Defined in [packages/mesh-graphql-client/src/types.ts:267](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L267)_

---

## `Optional` hash

• **hash**? : _undefined | string_

_Defined in [packages/mesh-graphql-client/src/types.ts:265](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L265)_

---

## message

• **message**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:268](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L268)_

---

## order

• **order**: _[StringifiedSignedOrder](#interface-stringifiedsignedorder)_

_Defined in [packages/mesh-graphql-client/src/types.ts:266](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L266)_

<hr />

# Interface: StringifiedSignedOrder

### Hierarchy

-   **StringifiedSignedOrder**

    ↳ [StringifiedOrderWithMetadata](#interface-stringifiedorderwithmetadata)

### Properties

## chainId

• **chainId**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:230](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L230)_

---

## exchangeAddress

• **exchangeAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:231](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L231)_

---

## expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:240](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L240)_

---

## feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:234](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L234)_

---

## makerAddress

• **makerAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:232](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L232)_

---

## makerAssetAmount

• **makerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:236](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L236)_

---

## makerAssetData

• **makerAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:242](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L242)_

---

## makerFee

• **makerFee**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:238](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L238)_

---

## makerFeeAssetData

• **makerFeeAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:244](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L244)_

---

## salt

• **salt**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:241](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L241)_

---

## senderAddress

• **senderAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:235](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L235)_

---

## signature

• **signature**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:246](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L246)_

---

## takerAddress

• **takerAddress**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:233](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L233)_

---

## takerAssetAmount

• **takerAssetAmount**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:237](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L237)_

---

## takerAssetData

• **takerAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:243](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L243)_

---

## takerFee

• **takerFee**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:239](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L239)_

---

## takerFeeAssetData

• **takerFeeAssetData**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:245](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L245)_

<hr />

# Interface: StringifiedStats

### Hierarchy

-   **StringifiedStats**

### Properties

## ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:226](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L226)_

---

## ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:225](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L225)_

---

## ethereumChainID

• **ethereumChainID**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:217](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L217)_

---

## latestBlock

• **latestBlock**: _[StringifiedLatestBlock](#interface-stringifiedlatestblock)_

_Defined in [packages/mesh-graphql-client/src/types.ts:218](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L218)_

---

## maxExpirationTime

• **maxExpirationTime**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:223](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L223)_

---

## numOrders

• **numOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:220](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L220)_

---

## numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:221](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L221)_

---

## numPeers

• **numPeers**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:219](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L219)_

---

## numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [packages/mesh-graphql-client/src/types.ts:222](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L222)_

---

## peerID

• **peerID**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:216](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L216)_

---

## pubSubTopic

• **pubSubTopic**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:213](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L213)_

---

## rendezvous

• **rendezvous**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:214](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L214)_

---

## secondaryRendezvous

• **secondaryRendezvous**: _string[]_

_Defined in [packages/mesh-graphql-client/src/types.ts:215](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L215)_

---

## startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:224](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L224)_

---

## version

• **version**: _string_

_Defined in [packages/mesh-graphql-client/src/types.ts:212](https://github.com/0xProject/0x-mesh/blob/b8104145/packages/mesh-graphql-client/src/types.ts#L212)_

<hr />
