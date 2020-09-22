# Class: WSClient

This class includes all the functionality related to interacting with a Mesh JSON RPC
websocket endpoint.

## Hierarchy

-   **WSClient**

## Index

### Constructors

-   [constructor](#constructor)

### Methods

-   [addOrdersAsync](#addordersasync)
-   [destroy](#destroy)
-   [getOrdersAsync](#getordersasync)
-   [getOrdersForPageAsync](#getordersforpageasync)
-   [getStatsAsync](#getstatsasync)
-   [onClose](#onclose)
-   [onReconnected](#onreconnected)
-   [subscribeToOrdersAsync](#subscribetoordersasync)
-   [unsubscribeAsync](#unsubscribeasync)

## Constructors

### constructor

\+ **new WSClient**(`url`: string, `wsOpts?`: [WSOpts](#interface-wsopts)): _[WSClient](#class-wsclient)_

_Defined in [ws_client.ts:252](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L252)_

Instantiates a new WSClient instance

**Parameters:**

| Name      | Type                        | Description        |
| --------- | --------------------------- | ------------------ |
| `url`     | string                      | WS server endpoint |
| `wsOpts?` | [WSOpts](#interface-wsopts) | WebSocket options  |

**Returns:** _[WSClient](#class-wsclient)_

An instance of WSClient

## Methods

### addOrdersAsync

▸ **addOrdersAsync**(`signedOrders`: SignedOrder[], `pinned`: boolean): _Promise‹[ValidationResults](#interface-validationresults)›_

_Defined in [ws_client.ts:281](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L281)_

Adds an array of 0x signed orders to the Mesh node.

**Parameters:**

| Name           | Type          | Default | Description                                                                                                                                                                                      |
| -------------- | ------------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `signedOrders` | SignedOrder[] | -       | signedOrders to add                                                                                                                                                                              |
| `pinned`       | boolean       | true    | Whether or not the orders should be pinned. Pinned orders will not be affected by any DDoS prevention or incentive mechanisms and will always stay in storage until they are no longer fillable. |

**Returns:** _Promise‹[ValidationResults](#interface-validationresults)›_

validation results

---

### destroy

▸ **destroy**(): _void_

_Defined in [ws_client.ts:421](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L421)_

destroy unsubscribes all active subscriptions, closes the websocket connection
and stops the internal heartbeat connection liveness check.

**Returns:** _void_

---

### getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

_Defined in [ws_client.ts:311](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L311)_

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

| Name      | Type   | Default | Description                                           |
| --------- | ------ | ------- | ----------------------------------------------------- |
| `perPage` | number | 200     | number of signedOrders to fetch per paginated request |

**Returns:** _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts

---

### getOrdersForPageAsync

▸ **getOrdersForPageAsync**(`page`: number, `perPage`: number, `snapshotID?`: undefined | string): _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

_Defined in [ws_client.ts:342](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L342)_

Get page of 0x signed orders stored on the Mesh node at the specified snapshot

**Parameters:**

| Name          | Type                    | Default | Description                                                                     |
| ------------- | ----------------------- | ------- | ------------------------------------------------------------------------------- |
| `page`        | number                  | -       | Page index at which to retrieve orders                                          |
| `perPage`     | number                  | 200     | number of signedOrders to fetch per paginated request                           |
| `snapshotID?` | undefined &#124; string | -       | The DB snapshot at which to fetch orders. If omitted, a new snapshot is created |

**Returns:** _Promise‹[GetOrdersResponse](#interface-getordersresponse)›_

the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts

---

### getStatsAsync

▸ **getStatsAsync**(): _Promise‹[GetStatsResponse](#interface-getstatsresponse)›_

_Defined in [ws_client.ts:302](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L302)_

**Returns:** _Promise‹[GetStatsResponse](#interface-getstatsresponse)›_

---

### onClose

▸ **onClose**(`cb`: function): _void_

_Defined in [ws_client.ts:403](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L403)_

Get notified when the underlying WS connection closes normally. If it closes with an
error, WSClient automatically attempts to re-connect without emitting a `close` event.

**Parameters:**

▪ **cb**: _function_

callback to call when WS connection closes

▸ (): _void_

**Returns:** _void_

---

### onReconnected

▸ **onReconnected**(`cb`: function): _void_

_Defined in [ws_client.ts:412](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L412)_

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: _function_

callback to call with the error when it occurs

▸ (): _void_

**Returns:** _void_

---

### subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**(`cb`: function): _Promise‹string›_

_Defined in [ws_client.ts:363](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L363)_

Subscribe to the 'orders' topic and receive order events from Mesh. This method returns a
subscriptionId that can be used to `unsubscribe()` from this subscription.

**Parameters:**

▪ **cb**: _function_

callback function where you'd like to get notified about order events

▸ (`orderEvents`: [OrderEvent](#interface-orderevent)[]): _void_

**Parameters:**

| Name          | Type                                  |
| ------------- | ------------------------------------- |
| `orderEvents` | [OrderEvent](#interface-orderevent)[] |

**Returns:** _Promise‹string›_

subscriptionId

---

### unsubscribeAsync

▸ **unsubscribeAsync**(`subscriptionId`: string): _Promise‹void›_

_Defined in [ws_client.ts:393](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/ws_client.ts#L393)_

Unsubscribe from a subscription

**Parameters:**

| Name             | Type   | Description                              |
| ---------------- | ------ | ---------------------------------------- |
| `subscriptionId` | string | identifier of the subscription to cancel |

**Returns:** _Promise‹void›_

<hr />

# Enumeration: ContractEventKind

## Index

### Enumeration members

-   [ERC1155ApprovalForAllEvent](#erc1155approvalforallevent)
-   [ERC1155TransferBatchEvent](#erc1155transferbatchevent)
-   [ERC1155TransferSingleEvent](#erc1155transfersingleevent)
-   [ERC20ApprovalEvent](#erc20approvalevent)
-   [ERC20TransferEvent](#erc20transferevent)
-   [ERC721ApprovalEvent](#erc721approvalevent)
-   [ERC721ApprovalForAllEvent](#erc721approvalforallevent)
-   [ERC721TransferEvent](#erc721transferevent)
-   [ExchangeCancelEvent](#exchangecancelevent)
-   [ExchangeCancelUpToEvent](#exchangecanceluptoevent)
-   [ExchangeFillEvent](#exchangefillevent)
-   [WethDepositEvent](#wethdepositevent)
-   [WethWithdrawalEvent](#wethwithdrawalevent)

## Enumeration members

### ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

_Defined in [types.ts:222](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L222)_

---

### ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in [types.ts:224](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L224)_

---

### ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in [types.ts:223](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L223)_

---

### ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in [types.ts:218](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L218)_

---

### ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in [types.ts:217](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L217)_

---

### ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in [types.ts:220](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L220)_

---

### ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in [types.ts:221](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L221)_

---

### ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in [types.ts:219](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L219)_

---

### ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in [types.ts:226](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L226)_

---

### ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in [types.ts:227](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L227)_

---

### ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in [types.ts:225](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L225)_

---

### WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in [types.ts:228](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L228)_

---

### WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in [types.ts:229](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L229)_

<hr />

# Enumeration: OrderEventEndState

## Index

### Enumeration members

-   [Added](#added)
-   [Cancelled](#cancelled)
-   [Expired](#expired)
-   [FillabilityIncreased](#fillabilityincreased)
-   [Filled](#filled)
-   [FullyFilled](#fullyfilled)
-   [Invalid](#invalid)
-   [StoppedWatching](#stoppedwatching)
-   [Unexpired](#unexpired)
-   [Unfunded](#unfunded)

## Enumeration members

### Added

• **Added**: = "ADDED"

_Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L286)_

---

### Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in [types.ts:289](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L289)_

---

### Expired

• **Expired**: = "EXPIRED"

_Defined in [types.ts:290](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L290)_

---

### FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

_Defined in [types.ts:294](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L294)_

---

### Filled

• **Filled**: = "FILLED"

_Defined in [types.ts:287](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L287)_

---

### FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

_Defined in [types.ts:288](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L288)_

---

### Invalid

• **Invalid**: = "INVALID"

_Defined in [types.ts:285](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L285)_

---

### StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

_Defined in [types.ts:292](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L292)_

---

### Unexpired

• **Unexpired**: = "UNEXPIRED"

_Defined in [types.ts:291](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L291)_

---

### Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in [types.ts:293](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L293)_

<hr />

# Enumeration: RejectedCode

## Index

### Enumeration members

-   [InternalError](#internalerror)
-   [MaxOrderSizeExceeded](#maxordersizeexceeded)
-   [NetworkRequestFailed](#networkrequestfailed)
-   [OrderAlreadyStored](#orderalreadystored)
-   [OrderCancelled](#ordercancelled)
-   [OrderExpired](#orderexpired)
-   [OrderForIncorrectChain](#orderforincorrectchain)
-   [OrderFullyFilled](#orderfullyfilled)
-   [OrderHasInvalidMakerAssetAmount](#orderhasinvalidmakerassetamount)
-   [OrderHasInvalidMakerAssetData](#orderhasinvalidmakerassetdata)
-   [OrderHasInvalidSignature](#orderhasinvalidsignature)
-   [OrderHasInvalidTakerAssetAmount](#orderhasinvalidtakerassetamount)
-   [OrderHasInvalidTakerAssetData](#orderhasinvalidtakerassetdata)
-   [OrderUnfunded](#orderunfunded)

## Enumeration members

### InternalError

• **InternalError**: = "InternalError"

_Defined in [types.ts:358](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L358)_

---

### MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MaxOrderSizeExceeded"

_Defined in [types.ts:359](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L359)_

---

### NetworkRequestFailed

• **NetworkRequestFailed**: = "NetworkRequestFailed"

_Defined in [types.ts:362](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L362)_

---

### OrderAlreadyStored

• **OrderAlreadyStored**: = "OrderAlreadyStored"

_Defined in [types.ts:360](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L360)_

---

### OrderCancelled

• **OrderCancelled**: = "OrderCancelled"

_Defined in [types.ts:367](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L367)_

---

### OrderExpired

• **OrderExpired**: = "OrderExpired"

_Defined in [types.ts:365](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L365)_

---

### OrderForIncorrectChain

• **OrderForIncorrectChain**: = "OrderForIncorrectChain"

_Defined in [types.ts:361](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L361)_

---

### OrderFullyFilled

• **OrderFullyFilled**: = "OrderFullyFilled"

_Defined in [types.ts:366](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L366)_

---

### OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "OrderHasInvalidMakerAssetAmount"

_Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L363)_

---

### OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "OrderHasInvalidMakerAssetData"

_Defined in [types.ts:369](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L369)_

---

### OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "OrderHasInvalidSignature"

_Defined in [types.ts:371](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L371)_

---

### OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "OrderHasInvalidTakerAssetAmount"

_Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L364)_

---

### OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "OrderHasInvalidTakerAssetData"

_Defined in [types.ts:370](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L370)_

---

### OrderUnfunded

• **OrderUnfunded**: = "OrderUnfunded"

_Defined in [types.ts:368](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L368)_

<hr />

# Enumeration: RejectedKind

## Index

### Enumeration members

-   [MeshError](#mesherror)
-   [MeshValidation](#meshvalidation)
-   [ZeroexValidation](#zeroexvalidation)

## Enumeration members

### MeshError

• **MeshError**: = "MESH_ERROR"

_Defined in [types.ts:353](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L353)_

---

### MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

_Defined in [types.ts:354](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L354)_

---

### ZeroexValidation

• **ZeroexValidation**: = "ZEROEX_VALIDATION"

_Defined in [types.ts:352](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L352)_

<hr />

# Interface: AcceptedOrderInfo

## Hierarchy

-   **AcceptedOrderInfo**

## Index

### Properties

-   [fillableTakerAssetAmount](#fillabletakerassetamount)
-   [isNew](#isnew)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)

## Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:335](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L335)_

---

### isNew

• **isNew**: _boolean_

_Defined in [types.ts:336](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L336)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:333](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L333)_

---

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:334](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L334)_

<hr />

# Interface: ClientConfig

WebSocketClient configs
Source: https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md#client-config-options

## Hierarchy

-   **ClientConfig**

## Index

### Properties

-   [assembleFragments](#optional-assemblefragments)
-   [closeTimeout](#optional-closetimeout)
-   [fragmentOutgoingMessages](#optional-fragmentoutgoingmessages)
-   [fragmentationThreshold](#optional-fragmentationthreshold)
-   [maxReceivedFrameSize](#optional-maxreceivedframesize)
-   [maxReceivedMessageSize](#optional-maxreceivedmessagesize)
-   [tlsOptions](#optional-tlsoptions)
-   [webSocketVersion](#optional-websocketversion)

## Properties

### `Optional` assembleFragments

• **assembleFragments**? : _undefined | false | true_

_Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L16)_

---

### `Optional` closeTimeout

• **closeTimeout**? : _undefined | number_

_Defined in [types.ts:17](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L17)_

---

### `Optional` fragmentOutgoingMessages

• **fragmentOutgoingMessages**? : _undefined | false | true_

_Defined in [types.ts:14](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L14)_

---

### `Optional` fragmentationThreshold

• **fragmentationThreshold**? : _undefined | number_

_Defined in [types.ts:15](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L15)_

---

### `Optional` maxReceivedFrameSize

• **maxReceivedFrameSize**? : _undefined | number_

_Defined in [types.ts:12](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L12)_

---

### `Optional` maxReceivedMessageSize

• **maxReceivedMessageSize**? : _undefined | number_

_Defined in [types.ts:13](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L13)_

---

### `Optional` tlsOptions

• **tlsOptions**? : _any_

_Defined in [types.ts:18](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L18)_

---

### `Optional` webSocketVersion

• **webSocketVersion**? : _undefined | number_

_Defined in [types.ts:11](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L11)_

<hr />

# Interface: ContractEvent

## Hierarchy

-   **ContractEvent**

## Index

### Properties

-   [address](#address)
-   [blockHash](#blockhash)
-   [isRemoved](#isremoved)
-   [kind](#kind)
-   [logIndex](#logindex)
-   [parameters](#parameters)
-   [txHash](#txhash)
-   [txIndex](#txindex)

## Properties

### address

• **address**: _string_

_Defined in [types.ts:279](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L279)_

---

### blockHash

• **blockHash**: _string_

_Defined in [types.ts:274](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L274)_

---

### isRemoved

• **isRemoved**: _string_

_Defined in [types.ts:278](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L278)_

---

### kind

• **kind**: _[ContractEventKind](#enumeration-contracteventkind)_

_Defined in [types.ts:280](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L280)_

---

### logIndex

• **logIndex**: _number_

_Defined in [types.ts:277](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L277)_

---

### parameters

• **parameters**: _[ContractEventParameters](#contracteventparameters)_

_Defined in [types.ts:281](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L281)_

---

### txHash

• **txHash**: _string_

_Defined in [types.ts:275](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L275)_

---

### txIndex

• **txIndex**: _number_

_Defined in [types.ts:276](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L276)_

<hr />

# Interface: ERC1155ApprovalForAllEvent

## Hierarchy

-   **ERC1155ApprovalForAllEvent**

## Index

### Properties

-   [approved](#approved)
-   [operator](#operator)
-   [owner](#owner)

## Properties

### approved

• **approved**: _boolean_

_Defined in [types.ts:144](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L144)_

---

### operator

• **operator**: _string_

_Defined in [types.ts:143](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L143)_

---

### owner

• **owner**: _string_

_Defined in [types.ts:142](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L142)_

<hr />

# Interface: ERC1155TransferBatchEvent

## Hierarchy

-   **ERC1155TransferBatchEvent**

## Index

### Properties

-   [from](#from)
-   [ids](#ids)
-   [operator](#operator)
-   [to](#to)
-   [values](#values)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:127](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L127)_

---

### ids

• **ids**: _BigNumber[]_

_Defined in [types.ts:129](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L129)_

---

### operator

• **operator**: _string_

_Defined in [types.ts:126](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L126)_

---

### to

• **to**: _string_

_Defined in [types.ts:128](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L128)_

---

### values

• **values**: _BigNumber[]_

_Defined in [types.ts:130](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L130)_

<hr />

# Interface: ERC1155TransferSingleEvent

## Hierarchy

-   **ERC1155TransferSingleEvent**

## Index

### Properties

-   [from](#from)
-   [id](#id)
-   [operator](#operator)
-   [to](#to)
-   [value](#value)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:111](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L111)_

---

### id

• **id**: _BigNumber_

_Defined in [types.ts:113](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L113)_

---

### operator

• **operator**: _string_

_Defined in [types.ts:110](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L110)_

---

### to

• **to**: _string_

_Defined in [types.ts:112](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L112)_

---

### value

• **value**: _BigNumber_

_Defined in [types.ts:114](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L114)_

<hr />

# Interface: ERC20ApprovalEvent

## Hierarchy

-   **ERC20ApprovalEvent**

## Index

### Properties

-   [owner](#owner)
-   [spender](#spender)
-   [value](#value)

## Properties

### owner

• **owner**: _string_

_Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L68)_

---

### spender

• **spender**: _string_

_Defined in [types.ts:69](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L69)_

---

### value

• **value**: _BigNumber_

_Defined in [types.ts:70](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L70)_

<hr />

# Interface: ERC20TransferEvent

## Hierarchy

-   **ERC20TransferEvent**

## Index

### Properties

-   [from](#from)
-   [to](#to)
-   [value](#value)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L56)_

---

### to

• **to**: _string_

_Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L57)_

---

### value

• **value**: _BigNumber_

_Defined in [types.ts:58](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L58)_

<hr />

# Interface: ERC721ApprovalEvent

## Hierarchy

-   **ERC721ApprovalEvent**

## Index

### Properties

-   [approved](#approved)
-   [owner](#owner)
-   [tokenId](#tokenid)

## Properties

### approved

• **approved**: _string_

_Defined in [types.ts:93](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L93)_

---

### owner

• **owner**: _string_

_Defined in [types.ts:92](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L92)_

---

### tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:94](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L94)_

<hr />

# Interface: ERC721ApprovalForAllEvent

## Hierarchy

-   **ERC721ApprovalForAllEvent**

## Index

### Properties

-   [approved](#approved)
-   [operator](#operator)
-   [owner](#owner)

## Properties

### approved

• **approved**: _boolean_

_Defined in [types.ts:106](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L106)_

---

### operator

• **operator**: _string_

_Defined in [types.ts:105](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L105)_

---

### owner

• **owner**: _string_

_Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L104)_

<hr />

# Interface: ERC721TransferEvent

## Hierarchy

-   **ERC721TransferEvent**

## Index

### Properties

-   [from](#from)
-   [to](#to)
-   [tokenId](#tokenid)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L80)_

---

### to

• **to**: _string_

_Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L81)_

---

### tokenId

• **tokenId**: _BigNumber_

_Defined in [types.ts:82](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L82)_

<hr />

# Interface: ExchangeCancelEvent

## Hierarchy

-   **ExchangeCancelEvent**

## Index

### Properties

-   [feeRecipientAddress](#feerecipientaddress)
-   [makerAddress](#makeraddress)
-   [makerAssetData](#makerassetdata)
-   [orderHash](#orderhash)
-   [senderAddress](#senderaddress)
-   [takerAssetData](#takerassetdata)

## Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:178](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L178)_

---

### makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:176](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L176)_

---

### makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:180](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L180)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:179](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L179)_

---

### senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:177](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L177)_

---

### takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:181](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L181)_

<hr />

# Interface: ExchangeCancelUpToEvent

## Hierarchy

-   **ExchangeCancelUpToEvent**

## Index

### Properties

-   [makerAddress](#makeraddress)
-   [orderEpoch](#orderepoch)
-   [senderAddress](#senderaddress)

## Properties

### makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L185)_

---

### orderEpoch

• **orderEpoch**: _BigNumber_

_Defined in [types.ts:187](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L187)_

---

### senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:186](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L186)_

<hr />

# Interface: ExchangeFillEvent

## Hierarchy

-   **ExchangeFillEvent**

## Index

### Properties

-   [feeRecipientAddress](#feerecipientaddress)
-   [makerAddress](#makeraddress)
-   [makerAssetData](#makerassetdata)
-   [makerAssetFilledAmount](#makerassetfilledamount)
-   [makerFeePaid](#makerfeepaid)
-   [orderHash](#orderhash)
-   [senderAddress](#senderaddress)
-   [takerAddress](#takeraddress)
-   [takerAssetData](#takerassetdata)
-   [takerAssetFilledAmount](#takerassetfilledamount)
-   [takerFeePaid](#takerfeepaid)

## Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:151](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L151)_

---

### makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:148](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L148)_

---

### makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:157](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L157)_

---

### makerAssetFilledAmount

• **makerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:152](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L152)_

---

### makerFeePaid

• **makerFeePaid**: _BigNumber_

_Defined in [types.ts:154](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L154)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:156](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L156)_

---

### senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:150](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L150)_

---

### takerAddress

• **takerAddress**: _string_

_Defined in [types.ts:149](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L149)_

---

### takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:158](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L158)_

---

### takerAssetFilledAmount

• **takerAssetFilledAmount**: _BigNumber_

_Defined in [types.ts:153](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L153)_

---

### takerFeePaid

• **takerFeePaid**: _BigNumber_

_Defined in [types.ts:155](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L155)_

<hr />

# Interface: GetOrdersResponse

## Hierarchy

-   **GetOrdersResponse**

## Index

### Properties

-   [ordersInfos](#ordersinfos)
-   [snapshotID](#snapshotid)
-   [snapshotTimestamp](#snapshottimestamp)

## Properties

### ordersInfos

• **ordersInfos**: _[OrderInfo](#interface-orderinfo)[]_

_Defined in [types.ts:415](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L415)_

---

### snapshotID

• **snapshotID**: _string_

_Defined in [types.ts:413](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L413)_

---

### snapshotTimestamp

• **snapshotTimestamp**: _number_

_Defined in [types.ts:414](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L414)_

<hr />

# Interface: GetStatsResponse

## Hierarchy

-   **GetStatsResponse**

## Index

### Properties

-   [ethRPCRateLimitExpiredRequests](#ethrpcratelimitexpiredrequests)
-   [ethRPCRequestsSentInCurrentUTCDay](#ethrpcrequestssentincurrentutcday)
-   [ethereumChainID](#ethereumchainid)
-   [latestBlock](#latestblock)
-   [maxExpirationTime](#maxexpirationtime)
-   [numOrders](#numorders)
-   [numOrdersIncludingRemoved](#numordersincludingremoved)
-   [numPeers](#numpeers)
-   [numPinnedOrders](#numpinnedorders)
-   [peerID](#peerid)
-   [pubSubTopic](#pubsubtopic)
-   [rendezvous](#rendezvous)
-   [startOfCurrentUTCDay](#startofcurrentutcday)
-   [version](#version)

## Properties

### ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in [types.ts:442](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L442)_

---

### ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in [types.ts:441](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L441)_

---

### ethereumChainID

• **ethereumChainID**: _number_

_Defined in [types.ts:433](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L433)_

---

### latestBlock

• **latestBlock**: _[LatestBlock](#interface-latestblock)_

_Defined in [types.ts:434](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L434)_

---

### maxExpirationTime

• **maxExpirationTime**: _string_

_Defined in [types.ts:439](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L439)_

---

### numOrders

• **numOrders**: _number_

_Defined in [types.ts:436](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L436)_

---

### numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in [types.ts:437](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L437)_

---

### numPeers

• **numPeers**: _number_

_Defined in [types.ts:435](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L435)_

---

### numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in [types.ts:438](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L438)_

---

### peerID

• **peerID**: _string_

_Defined in [types.ts:432](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L432)_

---

### pubSubTopic

• **pubSubTopic**: _string_

_Defined in [types.ts:430](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L430)_

---

### rendezvous

• **rendezvous**: _string_

_Defined in [types.ts:431](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L431)_

---

### startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _string_

_Defined in [types.ts:440](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L440)_

---

### version

• **version**: _string_

_Defined in [types.ts:429](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L429)_

<hr />

# Interface: HeartbeatEventPayload

## Hierarchy

-   **HeartbeatEventPayload**

## Index

### Properties

-   [result](#result)
-   [subscription](#subscription)

## Properties

### result

• **result**: _string_

_Defined in [types.ts:304](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L304)_

---

### subscription

• **subscription**: _string_

_Defined in [types.ts:303](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L303)_

<hr />

# Interface: LatestBlock

## Hierarchy

-   **LatestBlock**

## Index

### Properties

-   [hash](#hash)
-   [number](#number)

## Properties

### hash

• **hash**: _string_

_Defined in [types.ts:425](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L425)_

---

### number

• **number**: _number_

_Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L424)_

<hr />

# Interface: OrderEvent

## Hierarchy

-   **OrderEvent**

## Index

### Properties

-   [contractEvents](#contractevents)
-   [endState](#endstate)
-   [fillableTakerAssetAmount](#fillabletakerassetamount)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)
-   [timestampMs](#timestampms)

## Properties

### contractEvents

• **contractEvents**: _[ContractEvent](#interface-contractevent)[]_

_Defined in [types.ts:322](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L322)_

---

### endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [types.ts:320](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L320)_

---

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:321](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L321)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:318](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L318)_

---

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:319](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L319)_

---

### timestampMs

• **timestampMs**: _number_

_Defined in [types.ts:317](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L317)_

<hr />

# Interface: OrderEventPayload

## Hierarchy

-   **OrderEventPayload**

## Index

### Properties

-   [result](#result)
-   [subscription](#subscription)

## Properties

### result

• **result**: _[RawOrderEvent](#interface-raworderevent)[]_

_Defined in [types.ts:299](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L299)_

---

### subscription

• **subscription**: _string_

_Defined in [types.ts:298](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L298)_

<hr />

# Interface: OrderInfo

## Hierarchy

-   **OrderInfo**

## Index

### Properties

-   [fillableTakerAssetAmount](#fillabletakerassetamount)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)

## Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _BigNumber_

_Defined in [types.ts:348](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L348)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L346)_

---

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:347](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L347)_

<hr />

# Interface: RawAcceptedOrderInfo

## Hierarchy

-   **RawAcceptedOrderInfo**

## Index

### Properties

-   [fillableTakerAssetAmount](#fillabletakerassetamount)
-   [isNew](#isnew)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)

## Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [types.ts:328](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L328)_

---

### isNew

• **isNew**: _boolean_

_Defined in [types.ts:329](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L329)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:326](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L326)_

---

### signedOrder

• **signedOrder**: _[StringifiedSignedOrder](#interface-stringifiedsignedorder)_

_Defined in [types.ts:327](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L327)_

<hr />

# Interface: RawGetOrdersResponse

## Hierarchy

-   **RawGetOrdersResponse**

## Index

### Properties

-   [ordersInfos](#ordersinfos)
-   [snapshotID](#snapshotid)
-   [snapshotTimestamp](#snapshottimestamp)

## Properties

### ordersInfos

• **ordersInfos**: _[RawAcceptedOrderInfo](#interface-rawacceptedorderinfo)[]_

_Defined in [types.ts:406](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L406)_

---

### snapshotID

• **snapshotID**: _string_

_Defined in [types.ts:404](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L404)_

---

### snapshotTimestamp

• **snapshotTimestamp**: _string_

_Defined in [types.ts:405](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L405)_

<hr />

# Interface: RawOrderEvent

## Hierarchy

-   **RawOrderEvent**

## Index

### Properties

-   [contractEvents](#contractevents)
-   [endState](#endstate)
-   [fillableTakerAssetAmount](#fillabletakerassetamount)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)
-   [timestamp](#timestamp)

## Properties

### contractEvents

• **contractEvents**: _[StringifiedContractEvent](#interface-stringifiedcontractevent)[]_

_Defined in [types.ts:313](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L313)_

---

### endState

• **endState**: _[OrderEventEndState](#enumeration-ordereventendstate)_

_Defined in [types.ts:311](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L311)_

---

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [types.ts:312](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L312)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:309](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L309)_

---

### signedOrder

• **signedOrder**: _[StringifiedSignedOrder](#interface-stringifiedsignedorder)_

_Defined in [types.ts:310](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L310)_

---

### timestamp

• **timestamp**: _string_

_Defined in [types.ts:308](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L308)_

<hr />

# Interface: RawOrderInfo

## Hierarchy

-   **RawOrderInfo**

## Index

### Properties

-   [fillableTakerAssetAmount](#fillabletakerassetamount)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)

## Properties

### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in [types.ts:342](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L342)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:340](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L340)_

---

### signedOrder

• **signedOrder**: _[StringifiedSignedOrder](#interface-stringifiedsignedorder)_

_Defined in [types.ts:341](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L341)_

<hr />

# Interface: RawRejectedOrderInfo

## Hierarchy

-   **RawRejectedOrderInfo**

## Index

### Properties

-   [kind](#kind)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)
-   [status](#status)

## Properties

### kind

• **kind**: _[RejectedKind](#enumeration-rejectedkind)_

_Defined in [types.ts:382](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L382)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:380](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L380)_

---

### signedOrder

• **signedOrder**: _[StringifiedSignedOrder](#interface-stringifiedsignedorder)_

_Defined in [types.ts:381](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L381)_

---

### status

• **status**: _[RejectedStatus](#interface-rejectedstatus)_

_Defined in [types.ts:383](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L383)_

<hr />

# Interface: RawValidationResults

## Hierarchy

-   **RawValidationResults**

## Index

### Properties

-   [accepted](#accepted)
-   [rejected](#rejected)

## Properties

### accepted

• **accepted**: _[RawAcceptedOrderInfo](#interface-rawacceptedorderinfo)[]_

_Defined in [types.ts:394](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L394)_

---

### rejected

• **rejected**: _[RawRejectedOrderInfo](#interface-rawrejectedorderinfo)[]_

_Defined in [types.ts:395](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L395)_

<hr />

# Interface: RejectedOrderInfo

## Hierarchy

-   **RejectedOrderInfo**

## Index

### Properties

-   [kind](#kind)
-   [orderHash](#orderhash)
-   [signedOrder](#signedorder)
-   [status](#status)

## Properties

### kind

• **kind**: _[RejectedKind](#enumeration-rejectedkind)_

_Defined in [types.ts:389](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L389)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:387](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L387)_

---

### signedOrder

• **signedOrder**: _SignedOrder_

_Defined in [types.ts:388](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L388)_

---

### status

• **status**: _[RejectedStatus](#interface-rejectedstatus)_

_Defined in [types.ts:390](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L390)_

<hr />

# Interface: RejectedStatus

## Hierarchy

-   **RejectedStatus**

## Index

### Properties

-   [code](#code)
-   [message](#message)

## Properties

### code

• **code**: _[RejectedCode](#enumeration-rejectedcode)_

_Defined in [types.ts:375](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L375)_

---

### message

• **message**: _string_

_Defined in [types.ts:376](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L376)_

<hr />

# Interface: StringifiedContractEvent

## Hierarchy

-   **StringifiedContractEvent**

## Index

### Properties

-   [address](#address)
-   [blockHash](#blockhash)
-   [isRemoved](#isremoved)
-   [kind](#kind)
-   [logIndex](#logindex)
-   [parameters](#parameters)
-   [txHash](#txhash)
-   [txIndex](#txindex)

## Properties

### address

• **address**: _string_

_Defined in [types.ts:253](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L253)_

---

### blockHash

• **blockHash**: _string_

_Defined in [types.ts:248](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L248)_

---

### isRemoved

• **isRemoved**: _string_

_Defined in [types.ts:252](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L252)_

---

### kind

• **kind**: _string_

_Defined in [types.ts:254](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L254)_

---

### logIndex

• **logIndex**: _number_

_Defined in [types.ts:251](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L251)_

---

### parameters

• **parameters**: _[StringifiedContractEventParameters](#stringifiedcontracteventparameters)_

_Defined in [types.ts:255](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L255)_

---

### txHash

• **txHash**: _string_

_Defined in [types.ts:249](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L249)_

---

### txIndex

• **txIndex**: _number_

_Defined in [types.ts:250](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L250)_

<hr />

# Interface: StringifiedERC1155TransferBatchEvent

## Hierarchy

-   **StringifiedERC1155TransferBatchEvent**

## Index

### Properties

-   [from](#from)
-   [ids](#ids)
-   [operator](#operator)
-   [to](#to)
-   [values](#values)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:135](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L135)_

---

### ids

• **ids**: _string[]_

_Defined in [types.ts:137](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L137)_

---

### operator

• **operator**: _string_

_Defined in [types.ts:134](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L134)_

---

### to

• **to**: _string_

_Defined in [types.ts:136](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L136)_

---

### values

• **values**: _string[]_

_Defined in [types.ts:138](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L138)_

<hr />

# Interface: StringifiedERC1155TransferSingleEvent

## Hierarchy

-   **StringifiedERC1155TransferSingleEvent**

## Index

### Properties

-   [from](#from)
-   [id](#id)
-   [operator](#operator)
-   [to](#to)
-   [value](#value)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:119](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L119)_

---

### id

• **id**: _string_

_Defined in [types.ts:121](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L121)_

---

### operator

• **operator**: _string_

_Defined in [types.ts:118](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L118)_

---

### to

• **to**: _string_

_Defined in [types.ts:120](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L120)_

---

### value

• **value**: _string_

_Defined in [types.ts:122](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L122)_

<hr />

# Interface: StringifiedERC20ApprovalEvent

## Hierarchy

-   **StringifiedERC20ApprovalEvent**

## Index

### Properties

-   [owner](#owner)
-   [spender](#spender)
-   [value](#value)

## Properties

### owner

• **owner**: _string_

_Defined in [types.ts:74](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L74)_

---

### spender

• **spender**: _string_

_Defined in [types.ts:75](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L75)_

---

### value

• **value**: _string_

_Defined in [types.ts:76](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L76)_

<hr />

# Interface: StringifiedERC20TransferEvent

## Hierarchy

-   **StringifiedERC20TransferEvent**

## Index

### Properties

-   [from](#from)
-   [to](#to)
-   [value](#value)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L62)_

---

### to

• **to**: _string_

_Defined in [types.ts:63](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L63)_

---

### value

• **value**: _string_

_Defined in [types.ts:64](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L64)_

<hr />

# Interface: StringifiedERC721ApprovalEvent

## Hierarchy

-   **StringifiedERC721ApprovalEvent**

## Index

### Properties

-   [approved](#approved)
-   [owner](#owner)
-   [tokenId](#tokenid)

## Properties

### approved

• **approved**: _string_

_Defined in [types.ts:99](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L99)_

---

### owner

• **owner**: _string_

_Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L98)_

---

### tokenId

• **tokenId**: _string_

_Defined in [types.ts:100](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L100)_

<hr />

# Interface: StringifiedERC721TransferEvent

## Hierarchy

-   **StringifiedERC721TransferEvent**

## Index

### Properties

-   [from](#from)
-   [to](#to)
-   [tokenId](#tokenid)

## Properties

### from

• **from**: _string_

_Defined in [types.ts:86](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L86)_

---

### to

• **to**: _string_

_Defined in [types.ts:87](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L87)_

---

### tokenId

• **tokenId**: _string_

_Defined in [types.ts:88](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L88)_

<hr />

# Interface: StringifiedExchangeCancelUpToEvent

## Hierarchy

-   **StringifiedExchangeCancelUpToEvent**

## Index

### Properties

-   [makerAddress](#makeraddress)
-   [orderEpoch](#orderepoch)
-   [senderAddress](#senderaddress)

## Properties

### makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:191](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L191)_

---

### orderEpoch

• **orderEpoch**: _string_

_Defined in [types.ts:193](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L193)_

---

### senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:192](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L192)_

<hr />

# Interface: StringifiedExchangeFillEvent

## Hierarchy

-   **StringifiedExchangeFillEvent**

## Index

### Properties

-   [feeRecipientAddress](#feerecipientaddress)
-   [makerAddress](#makeraddress)
-   [makerAssetData](#makerassetdata)
-   [makerAssetFilledAmount](#makerassetfilledamount)
-   [makerFeePaid](#makerfeepaid)
-   [orderHash](#orderhash)
-   [senderAddress](#senderaddress)
-   [takerAddress](#takeraddress)
-   [takerAssetData](#takerassetdata)
-   [takerAssetFilledAmount](#takerassetfilledamount)
-   [takerFeePaid](#takerfeepaid)

## Properties

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:165](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L165)_

---

### makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:162](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L162)_

---

### makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:171](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L171)_

---

### makerAssetFilledAmount

• **makerAssetFilledAmount**: _string_

_Defined in [types.ts:166](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L166)_

---

### makerFeePaid

• **makerFeePaid**: _string_

_Defined in [types.ts:168](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L168)_

---

### orderHash

• **orderHash**: _string_

_Defined in [types.ts:170](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L170)_

---

### senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:164](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L164)_

---

### takerAddress

• **takerAddress**: _string_

_Defined in [types.ts:163](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L163)_

---

### takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:172](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L172)_

---

### takerAssetFilledAmount

• **takerAssetFilledAmount**: _string_

_Defined in [types.ts:167](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L167)_

---

### takerFeePaid

• **takerFeePaid**: _string_

_Defined in [types.ts:169](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L169)_

<hr />

# Interface: StringifiedSignedOrder

## Hierarchy

-   **StringifiedSignedOrder**

## Index

### Properties

-   [exchangeAddress](#exchangeaddress)
-   [expirationTimeSeconds](#expirationtimeseconds)
-   [feeRecipientAddress](#feerecipientaddress)
-   [makerAddress](#makeraddress)
-   [makerAssetAmount](#makerassetamount)
-   [makerAssetData](#makerassetdata)
-   [makerFee](#makerfee)
-   [salt](#salt)
-   [senderAddress](#senderaddress)
-   [signature](#signature)
-   [takerAddress](#takeraddress)
-   [takerAssetAmount](#takerassetamount)
-   [takerAssetData](#takerassetdata)
-   [takerFee](#takerfee)

## Properties

### exchangeAddress

• **exchangeAddress**: _string_

_Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L49)_

---

### expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Defined in [types.ts:51](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L51)_

---

### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L50)_

---

### makerAddress

• **makerAddress**: _string_

_Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L40)_

---

### makerAssetAmount

• **makerAssetAmount**: _string_

_Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L44)_

---

### makerAssetData

• **makerAssetData**: _string_

_Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L46)_

---

### makerFee

• **makerFee**: _string_

_Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L42)_

---

### salt

• **salt**: _string_

_Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L48)_

---

### senderAddress

• **senderAddress**: _string_

_Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L39)_

---

### signature

• **signature**: _string_

_Defined in [types.ts:52](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L52)_

---

### takerAddress

• **takerAddress**: _string_

_Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L41)_

---

### takerAssetAmount

• **takerAssetAmount**: _string_

_Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L45)_

---

### takerAssetData

• **takerAssetData**: _string_

_Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L47)_

---

### takerFee

• **takerFee**: _string_

_Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L43)_

<hr />

# Interface: StringifiedWethDepositEvent

## Hierarchy

-   **StringifiedWethDepositEvent**

## Index

### Properties

-   [owner](#owner)
-   [value](#value)

## Properties

### owner

• **owner**: _string_

_Defined in [types.ts:212](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L212)_

---

### value

• **value**: _string_

_Defined in [types.ts:213](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L213)_

<hr />

# Interface: StringifiedWethWithdrawalEvent

## Hierarchy

-   **StringifiedWethWithdrawalEvent**

## Index

### Properties

-   [owner](#owner)
-   [value](#value)

## Properties

### owner

• **owner**: _string_

_Defined in [types.ts:202](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L202)_

---

### value

• **value**: _string_

_Defined in [types.ts:203](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L203)_

<hr />

# Interface: ValidationResults

## Hierarchy

-   **ValidationResults**

## Index

### Properties

-   [accepted](#accepted)
-   [rejected](#rejected)

## Properties

### accepted

• **accepted**: _[AcceptedOrderInfo](#interface-acceptedorderinfo)[]_

_Defined in [types.ts:399](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L399)_

---

### rejected

• **rejected**: _[RejectedOrderInfo](#interface-rejectedorderinfo)[]_

_Defined in [types.ts:400](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L400)_

<hr />

# Interface: WethDepositEvent

## Hierarchy

-   **WethDepositEvent**

## Index

### Properties

-   [owner](#owner)
-   [value](#value)

## Properties

### owner

• **owner**: _string_

_Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L207)_

---

### value

• **value**: _BigNumber_

_Defined in [types.ts:208](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L208)_

<hr />

# Interface: WethWithdrawalEvent

## Hierarchy

-   **WethWithdrawalEvent**

## Index

### Properties

-   [owner](#owner)
-   [value](#value)

## Properties

### owner

• **owner**: _string_

_Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L197)_

---

### value

• **value**: _BigNumber_

_Defined in [types.ts:198](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L198)_

<hr />

# Interface: WSMessage

## Hierarchy

-   **WSMessage**

## Index

### Properties

-   [type](#type)
-   [utf8Data](#utf8data)

## Properties

### type

• **type**: _string_

_Defined in [types.ts:419](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L419)_

---

### utf8Data

• **utf8Data**: _string_

_Defined in [types.ts:420](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L420)_

<hr />

# Interface: WSOpts

timeout: timeout in milliseconds to enforce on every WS request that expects a response
headers: Request headers (e.g., authorization)
protocol: requestOptions should be either null or an object specifying additional configuration options to be
passed to http.request or https.request. This can be used to pass a custom agent to enable WebSocketClient usage
from behind an HTTP or HTTPS proxy server using koichik/node-tunnel or similar.
clientConfig: The client configs documented here: https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md
reconnectDelay: time in milliseconds after which to attempt to reconnect to WS server after an error occurred (default: 5000)

## Hierarchy

-   **WSOpts**

## Index

### Properties

-   [clientConfig](#optional-clientconfig)
-   [headers](#optional-headers)
-   [protocol](#optional-protocol)
-   [reconnectDelay](#optional-reconnectdelay)
-   [timeout](#optional-timeout)

## Properties

### `Optional` clientConfig

• **clientConfig**? : _[ClientConfig](#interface-clientconfig)_

_Defined in [types.ts:34](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L34)_

---

### `Optional` headers

• **headers**? : _undefined | \_\_type_

_Defined in [types.ts:32](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L32)_

---

### `Optional` protocol

• **protocol**? : _undefined | string_

_Defined in [types.ts:33](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L33)_

---

### `Optional` reconnectDelay

• **reconnectDelay**? : _undefined | number_

_Defined in [types.ts:35](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L35)_

---

### `Optional` timeout

• **timeout**? : _undefined | number_

_Defined in [types.ts:31](https://github.com/0xProject/0x-mesh/blob/d7f70fc4/packages/rpc-client/src/types.ts#L31)_

<hr />
