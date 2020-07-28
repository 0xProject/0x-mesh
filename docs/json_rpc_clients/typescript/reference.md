# Class: WSClient

This class includes all the functionality related to interacting with a Mesh JSON RPC
websocket endpoint.

## Hierarchy

* **WSClient**

## Index

### Constructors

* [constructor](#constructor)

### Methods

* [addOrdersAsync](#addordersasync)
* [destroy](#destroy)
* [getOrdersAsync](#getordersasync)
* [getOrdersForPageAsync](#getordersforpageasync)
* [getStatsAsync](#getstatsasync)
* [onClose](#onclose)
* [onReconnected](#onreconnected)
* [subscribeToOrdersAsync](#subscribetoordersasync)
* [unsubscribeAsync](#unsubscribeasync)

## Constructors

###  constructor

\+ **new WSClient**(`url`: string, `wsOpts?`: [WSOpts](#interface-wsopts)): *[WSClient](#class-wsclient)*

*Defined in [ws_client.ts:252](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L252)*

Instantiates a new WSClient instance

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`url` | string | WS server endpoint |
`wsOpts?` | [WSOpts](#interface-wsopts) | WebSocket options |

**Returns:** *[WSClient](#class-wsclient)*

An instance of WSClient

## Methods

###  addOrdersAsync

▸ **addOrdersAsync**(`signedOrders`: SignedOrder[], `pinned`: boolean): *Promise‹[ValidationResults](#interface-validationresults)›*

*Defined in [ws_client.ts:281](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L281)*

Adds an array of 0x signed orders to the Mesh node.

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`signedOrders` | SignedOrder[] | - | signedOrders to add |
`pinned` | boolean | true | Whether or not the orders should be pinned. Pinned orders will not be affected by any DDoS prevention or incentive mechanisms and will always stay in storage until they are no longer fillable. |

**Returns:** *Promise‹[ValidationResults](#interface-validationresults)›*

validation results

___

###  destroy

▸ **destroy**(): *void*

*Defined in [ws_client.ts:421](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L421)*

destroy unsubscribes all active subscriptions, closes the websocket connection
and stops the internal heartbeat connection liveness check.

**Returns:** *void*

___

###  getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

*Defined in [ws_client.ts:311](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L311)*

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`perPage` | number | 200 | number of signedOrders to fetch per paginated request |

**Returns:** *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts

___

###  getOrdersForPageAsync

▸ **getOrdersForPageAsync**(`page`: number, `perPage`: number, `snapshotID?`: undefined | string): *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

*Defined in [ws_client.ts:342](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L342)*

Get page of 0x signed orders stored on the Mesh node at the specified snapshot

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`page` | number | - | Page index at which to retrieve orders |
`perPage` | number | 200 | number of signedOrders to fetch per paginated request |
`snapshotID?` | undefined &#124; string | - | The DB snapshot at which to fetch orders. If omitted, a new snapshot is created |

**Returns:** *Promise‹[GetOrdersResponse](#interface-getordersresponse)›*

the snapshotID, snapshotTimestamp and all orders, their hashes and fillableTakerAssetAmounts

___

###  getStatsAsync

▸ **getStatsAsync**(): *Promise‹[GetStatsResponse](#interface-getstatsresponse)›*

*Defined in [ws_client.ts:302](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L302)*

**Returns:** *Promise‹[GetStatsResponse](#interface-getstatsresponse)›*

___

###  onClose

▸ **onClose**(`cb`: function): *void*

*Defined in [ws_client.ts:403](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L403)*

Get notified when the underlying WS connection closes normally. If it closes with an
error, WSClient automatically attempts to re-connect without emitting a `close` event.

**Parameters:**

▪ **cb**: *function*

callback to call when WS connection closes

▸ (): *void*

**Returns:** *void*

___

###  onReconnected

▸ **onReconnected**(`cb`: function): *void*

*Defined in [ws_client.ts:412](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L412)*

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: *function*

callback to call with the error when it occurs

▸ (): *void*

**Returns:** *void*

___

###  subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**(`cb`: function): *Promise‹string›*

*Defined in [ws_client.ts:363](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L363)*

Subscribe to the 'orders' topic and receive order events from Mesh. This method returns a
subscriptionId that can be used to `unsubscribe()` from this subscription.

**Parameters:**

▪ **cb**: *function*

callback function where you'd like to get notified about order events

▸ (`orderEvents`: [OrderEvent](#interface-orderevent)[]): *void*

**Parameters:**

Name | Type |
------ | ------ |
`orderEvents` | [OrderEvent](#interface-orderevent)[] |

**Returns:** *Promise‹string›*

subscriptionId

___

###  unsubscribeAsync

▸ **unsubscribeAsync**(`subscriptionId`: string): *Promise‹void›*

*Defined in [ws_client.ts:393](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/ws_client.ts#L393)*

Unsubscribe from a subscription

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`subscriptionId` | string | identifier of the subscription to cancel  |

**Returns:** *Promise‹void›*


<hr />

# Enumeration: ContractEventKind

## Index

### Enumeration members

* [ERC1155ApprovalForAllEvent](#erc1155approvalforallevent)
* [ERC1155TransferBatchEvent](#erc1155transferbatchevent)
* [ERC1155TransferSingleEvent](#erc1155transfersingleevent)
* [ERC20ApprovalEvent](#erc20approvalevent)
* [ERC20TransferEvent](#erc20transferevent)
* [ERC721ApprovalEvent](#erc721approvalevent)
* [ERC721ApprovalForAllEvent](#erc721approvalforallevent)
* [ERC721TransferEvent](#erc721transferevent)
* [ExchangeCancelEvent](#exchangecancelevent)
* [ExchangeCancelUpToEvent](#exchangecanceluptoevent)
* [ExchangeFillEvent](#exchangefillevent)
* [WethDepositEvent](#wethdepositevent)
* [WethWithdrawalEvent](#wethwithdrawalevent)

## Enumeration members

###  ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

*Defined in [types.ts:222](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L222)*

___

###  ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

*Defined in [types.ts:224](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L224)*

___

###  ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

*Defined in [types.ts:223](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L223)*

___

###  ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

*Defined in [types.ts:218](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L218)*

___

###  ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

*Defined in [types.ts:217](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L217)*

___

###  ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

*Defined in [types.ts:220](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L220)*

___

###  ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

*Defined in [types.ts:221](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L221)*

___

###  ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

*Defined in [types.ts:219](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L219)*

___

###  ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

*Defined in [types.ts:226](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L226)*

___

###  ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

*Defined in [types.ts:227](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L227)*

___

###  ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

*Defined in [types.ts:225](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L225)*

___

###  WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

*Defined in [types.ts:228](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L228)*

___

###  WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

*Defined in [types.ts:229](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L229)*


<hr />

# Enumeration: OrderEventEndState

## Index

### Enumeration members

* [Added](#added)
* [Cancelled](#cancelled)
* [Expired](#expired)
* [FillabilityIncreased](#fillabilityincreased)
* [Filled](#filled)
* [FullyFilled](#fullyfilled)
* [Invalid](#invalid)
* [StoppedWatching](#stoppedwatching)
* [Unexpired](#unexpired)
* [Unfunded](#unfunded)

## Enumeration members

###  Added

• **Added**: = "ADDED"

*Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L286)*

___

###  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [types.ts:289](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L289)*

___

###  Expired

• **Expired**: = "EXPIRED"

*Defined in [types.ts:290](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L290)*

___

###  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [types.ts:294](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L294)*

___

###  Filled

• **Filled**: = "FILLED"

*Defined in [types.ts:287](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L287)*

___

###  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [types.ts:288](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L288)*

___

###  Invalid

• **Invalid**: = "INVALID"

*Defined in [types.ts:285](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L285)*

___

###  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [types.ts:292](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L292)*

___

###  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [types.ts:291](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L291)*

___

###  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [types.ts:293](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L293)*


<hr />

# Enumeration: RejectedCode

## Index

### Enumeration members

* [InternalError](#internalerror)
* [MaxOrderSizeExceeded](#maxordersizeexceeded)
* [NetworkRequestFailed](#networkrequestfailed)
* [OrderAlreadyStored](#orderalreadystored)
* [OrderCancelled](#ordercancelled)
* [OrderExpired](#orderexpired)
* [OrderForIncorrectChain](#orderforincorrectchain)
* [OrderFullyFilled](#orderfullyfilled)
* [OrderHasInvalidMakerAssetAmount](#orderhasinvalidmakerassetamount)
* [OrderHasInvalidMakerAssetData](#orderhasinvalidmakerassetdata)
* [OrderHasInvalidSignature](#orderhasinvalidsignature)
* [OrderHasInvalidTakerAssetAmount](#orderhasinvalidtakerassetamount)
* [OrderHasInvalidTakerAssetData](#orderhasinvalidtakerassetdata)
* [OrderUnfunded](#orderunfunded)

## Enumeration members

###  InternalError

• **InternalError**: = "InternalError"

*Defined in [types.ts:358](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L358)*

___

###  MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MaxOrderSizeExceeded"

*Defined in [types.ts:359](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L359)*

___

###  NetworkRequestFailed

• **NetworkRequestFailed**: = "NetworkRequestFailed"

*Defined in [types.ts:362](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L362)*

___

###  OrderAlreadyStored

• **OrderAlreadyStored**: = "OrderAlreadyStored"

*Defined in [types.ts:360](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L360)*

___

###  OrderCancelled

• **OrderCancelled**: = "OrderCancelled"

*Defined in [types.ts:367](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L367)*

___

###  OrderExpired

• **OrderExpired**: = "OrderExpired"

*Defined in [types.ts:365](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L365)*

___

###  OrderForIncorrectChain

• **OrderForIncorrectChain**: = "OrderForIncorrectChain"

*Defined in [types.ts:361](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L361)*

___

###  OrderFullyFilled

• **OrderFullyFilled**: = "OrderFullyFilled"

*Defined in [types.ts:366](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L366)*

___

###  OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "OrderHasInvalidMakerAssetAmount"

*Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L363)*

___

###  OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "OrderHasInvalidMakerAssetData"

*Defined in [types.ts:369](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L369)*

___

###  OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "OrderHasInvalidSignature"

*Defined in [types.ts:371](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L371)*

___

###  OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "OrderHasInvalidTakerAssetAmount"

*Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L364)*

___

###  OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "OrderHasInvalidTakerAssetData"

*Defined in [types.ts:370](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L370)*

___

###  OrderUnfunded

• **OrderUnfunded**: = "OrderUnfunded"

*Defined in [types.ts:368](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L368)*


<hr />

# Enumeration: RejectedKind

## Index

### Enumeration members

* [MeshError](#mesherror)
* [MeshValidation](#meshvalidation)
* [ZeroexValidation](#zeroexvalidation)

## Enumeration members

###  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [types.ts:353](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L353)*

___

###  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [types.ts:354](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L354)*

___

###  ZeroexValidation

• **ZeroexValidation**: = "ZEROEX_VALIDATION"

*Defined in [types.ts:352](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L352)*


<hr />

# Interface: AcceptedOrderInfo

## Hierarchy

* **AcceptedOrderInfo**

## Index

### Properties

* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [isNew](#isnew)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)

## Properties

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:335](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L335)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:336](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L336)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:333](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L333)*

___

###  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:334](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L334)*


<hr />

# Interface: ClientConfig

WebSocketClient configs
Source: https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md#client-config-options

## Hierarchy

* **ClientConfig**

## Index

### Properties

* [assembleFragments](#optional-assemblefragments)
* [closeTimeout](#optional-closetimeout)
* [fragmentOutgoingMessages](#optional-fragmentoutgoingmessages)
* [fragmentationThreshold](#optional-fragmentationthreshold)
* [maxReceivedFrameSize](#optional-maxreceivedframesize)
* [maxReceivedMessageSize](#optional-maxreceivedmessagesize)
* [tlsOptions](#optional-tlsoptions)
* [webSocketVersion](#optional-websocketversion)

## Properties

### `Optional` assembleFragments

• **assembleFragments**? : *undefined | false | true*

*Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L16)*

___

### `Optional` closeTimeout

• **closeTimeout**? : *undefined | number*

*Defined in [types.ts:17](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L17)*

___

### `Optional` fragmentOutgoingMessages

• **fragmentOutgoingMessages**? : *undefined | false | true*

*Defined in [types.ts:14](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L14)*

___

### `Optional` fragmentationThreshold

• **fragmentationThreshold**? : *undefined | number*

*Defined in [types.ts:15](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L15)*

___

### `Optional` maxReceivedFrameSize

• **maxReceivedFrameSize**? : *undefined | number*

*Defined in [types.ts:12](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L12)*

___

### `Optional` maxReceivedMessageSize

• **maxReceivedMessageSize**? : *undefined | number*

*Defined in [types.ts:13](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L13)*

___

### `Optional` tlsOptions

• **tlsOptions**? : *any*

*Defined in [types.ts:18](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L18)*

___

### `Optional` webSocketVersion

• **webSocketVersion**? : *undefined | number*

*Defined in [types.ts:11](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L11)*


<hr />

# Interface: ContractEvent

## Hierarchy

* **ContractEvent**

## Index

### Properties

* [address](#address)
* [blockHash](#blockhash)
* [isRemoved](#isremoved)
* [kind](#kind)
* [logIndex](#logindex)
* [parameters](#parameters)
* [txHash](#txhash)
* [txIndex](#txindex)

## Properties

###  address

• **address**: *string*

*Defined in [types.ts:279](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L279)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:274](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L274)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:278](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L278)*

___

###  kind

• **kind**: *[ContractEventKind](#enumeration-contracteventkind)*

*Defined in [types.ts:280](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L280)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:277](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L277)*

___

###  parameters

• **parameters**: *[ContractEventParameters](#contracteventparameters)*

*Defined in [types.ts:281](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L281)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:275](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L275)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:276](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L276)*


<hr />

# Interface: ERC1155ApprovalForAllEvent

## Hierarchy

* **ERC1155ApprovalForAllEvent**

## Index

### Properties

* [approved](#approved)
* [operator](#operator)
* [owner](#owner)

## Properties

###  approved

• **approved**: *boolean*

*Defined in [types.ts:144](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L144)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:143](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L143)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:142](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L142)*


<hr />

# Interface: ERC1155TransferBatchEvent

## Hierarchy

* **ERC1155TransferBatchEvent**

## Index

### Properties

* [from](#from)
* [ids](#ids)
* [operator](#operator)
* [to](#to)
* [values](#values)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:127](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L127)*

___

###  ids

• **ids**: *BigNumber[]*

*Defined in [types.ts:129](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L129)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:126](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L126)*

___

###  to

• **to**: *string*

*Defined in [types.ts:128](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L128)*

___

###  values

• **values**: *BigNumber[]*

*Defined in [types.ts:130](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L130)*


<hr />

# Interface: ERC1155TransferSingleEvent

## Hierarchy

* **ERC1155TransferSingleEvent**

## Index

### Properties

* [from](#from)
* [id](#id)
* [operator](#operator)
* [to](#to)
* [value](#value)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:111](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L111)*

___

###  id

• **id**: *BigNumber*

*Defined in [types.ts:113](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L113)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:110](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L110)*

___

###  to

• **to**: *string*

*Defined in [types.ts:112](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L112)*

___

###  value

• **value**: *BigNumber*

*Defined in [types.ts:114](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L114)*


<hr />

# Interface: ERC20ApprovalEvent

## Hierarchy

* **ERC20ApprovalEvent**

## Index

### Properties

* [owner](#owner)
* [spender](#spender)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L68)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:69](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L69)*

___

###  value

• **value**: *BigNumber*

*Defined in [types.ts:70](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L70)*


<hr />

# Interface: ERC20TransferEvent

## Hierarchy

* **ERC20TransferEvent**

## Index

### Properties

* [from](#from)
* [to](#to)
* [value](#value)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L56)*

___

###  to

• **to**: *string*

*Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L57)*

___

###  value

• **value**: *BigNumber*

*Defined in [types.ts:58](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L58)*


<hr />

# Interface: ERC721ApprovalEvent

## Hierarchy

* **ERC721ApprovalEvent**

## Index

### Properties

* [approved](#approved)
* [owner](#owner)
* [tokenId](#tokenid)

## Properties

###  approved

• **approved**: *string*

*Defined in [types.ts:93](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L93)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:92](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L92)*

___

###  tokenId

• **tokenId**: *BigNumber*

*Defined in [types.ts:94](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L94)*


<hr />

# Interface: ERC721ApprovalForAllEvent

## Hierarchy

* **ERC721ApprovalForAllEvent**

## Index

### Properties

* [approved](#approved)
* [operator](#operator)
* [owner](#owner)

## Properties

###  approved

• **approved**: *boolean*

*Defined in [types.ts:106](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L106)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:105](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L105)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L104)*


<hr />

# Interface: ERC721TransferEvent

## Hierarchy

* **ERC721TransferEvent**

## Index

### Properties

* [from](#from)
* [to](#to)
* [tokenId](#tokenid)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L80)*

___

###  to

• **to**: *string*

*Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L81)*

___

###  tokenId

• **tokenId**: *BigNumber*

*Defined in [types.ts:82](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L82)*


<hr />

# Interface: ExchangeCancelEvent

## Hierarchy

* **ExchangeCancelEvent**

## Index

### Properties

* [feeRecipientAddress](#feerecipientaddress)
* [makerAddress](#makeraddress)
* [makerAssetData](#makerassetdata)
* [orderHash](#orderhash)
* [senderAddress](#senderaddress)
* [takerAssetData](#takerassetdata)

## Properties

###  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:178](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L178)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:176](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L176)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:180](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L180)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:179](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L179)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:177](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L177)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:181](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L181)*


<hr />

# Interface: ExchangeCancelUpToEvent

## Hierarchy

* **ExchangeCancelUpToEvent**

## Index

### Properties

* [makerAddress](#makeraddress)
* [orderEpoch](#orderepoch)
* [senderAddress](#senderaddress)

## Properties

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L185)*

___

###  orderEpoch

• **orderEpoch**: *BigNumber*

*Defined in [types.ts:187](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L187)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:186](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L186)*


<hr />

# Interface: ExchangeFillEvent

## Hierarchy

* **ExchangeFillEvent**

## Index

### Properties

* [feeRecipientAddress](#feerecipientaddress)
* [makerAddress](#makeraddress)
* [makerAssetData](#makerassetdata)
* [makerAssetFilledAmount](#makerassetfilledamount)
* [makerFeePaid](#makerfeepaid)
* [orderHash](#orderhash)
* [senderAddress](#senderaddress)
* [takerAddress](#takeraddress)
* [takerAssetData](#takerassetdata)
* [takerAssetFilledAmount](#takerassetfilledamount)
* [takerFeePaid](#takerfeepaid)

## Properties

###  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:151](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L151)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:148](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L148)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:157](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L157)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *BigNumber*

*Defined in [types.ts:152](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L152)*

___

###  makerFeePaid

• **makerFeePaid**: *BigNumber*

*Defined in [types.ts:154](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L154)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:156](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L156)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:150](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L150)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:149](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L149)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:158](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L158)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *BigNumber*

*Defined in [types.ts:153](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L153)*

___

###  takerFeePaid

• **takerFeePaid**: *BigNumber*

*Defined in [types.ts:155](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L155)*


<hr />

# Interface: GetOrdersResponse

## Hierarchy

* **GetOrdersResponse**

## Index

### Properties

* [ordersInfos](#ordersinfos)
* [snapshotID](#snapshotid)
* [snapshotTimestamp](#snapshottimestamp)

## Properties

###  ordersInfos

• **ordersInfos**: *[OrderInfo](#interface-orderinfo)[]*

*Defined in [types.ts:415](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L415)*

___

###  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:413](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L413)*

___

###  snapshotTimestamp

• **snapshotTimestamp**: *number*

*Defined in [types.ts:414](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L414)*


<hr />

# Interface: GetStatsResponse

## Hierarchy

* **GetStatsResponse**

## Index

### Properties

* [ethRPCRateLimitExpiredRequests](#ethrpcratelimitexpiredrequests)
* [ethRPCRequestsSentInCurrentUTCDay](#ethrpcrequestssentincurrentutcday)
* [ethereumChainID](#ethereumchainid)
* [latestBlock](#latestblock)
* [maxExpirationTime](#maxexpirationtime)
* [numOrders](#numorders)
* [numOrdersIncludingRemoved](#numordersincludingremoved)
* [numPeers](#numpeers)
* [numPinnedOrders](#numpinnedorders)
* [peerID](#peerid)
* [pubSubTopic](#pubsubtopic)
* [rendezvous](#rendezvous)
* [startOfCurrentUTCDay](#startofcurrentutcday)
* [version](#version)

## Properties

###  ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: *number*

*Defined in [types.ts:442](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L442)*

___

###  ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: *number*

*Defined in [types.ts:441](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L441)*

___

###  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:433](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L433)*

___

###  latestBlock

• **latestBlock**: *[LatestBlock](#interface-latestblock)*

*Defined in [types.ts:434](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L434)*

___

###  maxExpirationTime

• **maxExpirationTime**: *string*

*Defined in [types.ts:439](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L439)*

___

###  numOrders

• **numOrders**: *number*

*Defined in [types.ts:436](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L436)*

___

###  numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: *number*

*Defined in [types.ts:437](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L437)*

___

###  numPeers

• **numPeers**: *number*

*Defined in [types.ts:435](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L435)*

___

###  numPinnedOrders

• **numPinnedOrders**: *number*

*Defined in [types.ts:438](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L438)*

___

###  peerID

• **peerID**: *string*

*Defined in [types.ts:432](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L432)*

___

###  pubSubTopic

• **pubSubTopic**: *string*

*Defined in [types.ts:430](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L430)*

___

###  rendezvous

• **rendezvous**: *string*

*Defined in [types.ts:431](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L431)*

___

###  startOfCurrentUTCDay

• **startOfCurrentUTCDay**: *string*

*Defined in [types.ts:440](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L440)*

___

###  version

• **version**: *string*

*Defined in [types.ts:429](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L429)*


<hr />

# Interface: HeartbeatEventPayload

## Hierarchy

* **HeartbeatEventPayload**

## Index

### Properties

* [result](#result)
* [subscription](#subscription)

## Properties

###  result

• **result**: *string*

*Defined in [types.ts:304](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L304)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:303](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L303)*


<hr />

# Interface: LatestBlock

## Hierarchy

* **LatestBlock**

## Index

### Properties

* [hash](#hash)
* [number](#number)

## Properties

###  hash

• **hash**: *string*

*Defined in [types.ts:425](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L425)*

___

###  number

• **number**: *number*

*Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L424)*


<hr />

# Interface: OrderEvent

## Hierarchy

* **OrderEvent**

## Index

### Properties

* [contractEvents](#contractevents)
* [endState](#endstate)
* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)
* [timestampMs](#timestampms)

## Properties

###  contractEvents

• **contractEvents**: *[ContractEvent](#interface-contractevent)[]*

*Defined in [types.ts:322](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L322)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:320](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L320)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:321](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L321)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:318](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L318)*

___

###  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:319](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L319)*

___

###  timestampMs

• **timestampMs**: *number*

*Defined in [types.ts:317](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L317)*


<hr />

# Interface: OrderEventPayload

## Hierarchy

* **OrderEventPayload**

## Index

### Properties

* [result](#result)
* [subscription](#subscription)

## Properties

###  result

• **result**: *[RawOrderEvent](#interface-raworderevent)[]*

*Defined in [types.ts:299](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L299)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:298](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L298)*


<hr />

# Interface: OrderInfo

## Hierarchy

* **OrderInfo**

## Index

### Properties

* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)

## Properties

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [types.ts:348](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L348)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L346)*

___

###  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:347](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L347)*


<hr />

# Interface: RawAcceptedOrderInfo

## Hierarchy

* **RawAcceptedOrderInfo**

## Index

### Properties

* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [isNew](#isnew)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)

## Properties

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:328](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L328)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:329](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L329)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:326](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L326)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:327](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L327)*


<hr />

# Interface: RawGetOrdersResponse

## Hierarchy

* **RawGetOrdersResponse**

## Index

### Properties

* [ordersInfos](#ordersinfos)
* [snapshotID](#snapshotid)
* [snapshotTimestamp](#snapshottimestamp)

## Properties

###  ordersInfos

• **ordersInfos**: *[RawAcceptedOrderInfo](#interface-rawacceptedorderinfo)[]*

*Defined in [types.ts:406](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L406)*

___

###  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:404](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L404)*

___

###  snapshotTimestamp

• **snapshotTimestamp**: *string*

*Defined in [types.ts:405](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L405)*


<hr />

# Interface: RawOrderEvent

## Hierarchy

* **RawOrderEvent**

## Index

### Properties

* [contractEvents](#contractevents)
* [endState](#endstate)
* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)
* [timestamp](#timestamp)

## Properties

###  contractEvents

• **contractEvents**: *[StringifiedContractEvent](#interface-stringifiedcontractevent)[]*

*Defined in [types.ts:313](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L313)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:311](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L311)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:312](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L312)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:309](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L309)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:310](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L310)*

___

###  timestamp

• **timestamp**: *string*

*Defined in [types.ts:308](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L308)*


<hr />

# Interface: RawOrderInfo

## Hierarchy

* **RawOrderInfo**

## Index

### Properties

* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)

## Properties

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:342](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L342)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:340](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L340)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:341](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L341)*


<hr />

# Interface: RawRejectedOrderInfo

## Hierarchy

* **RawRejectedOrderInfo**

## Index

### Properties

* [kind](#kind)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)
* [status](#status)

## Properties

###  kind

• **kind**: *[RejectedKind](#enumeration-rejectedkind)*

*Defined in [types.ts:382](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L382)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:380](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L380)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:381](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L381)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:383](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L383)*


<hr />

# Interface: RawValidationResults

## Hierarchy

* **RawValidationResults**

## Index

### Properties

* [accepted](#accepted)
* [rejected](#rejected)

## Properties

###  accepted

• **accepted**: *[RawAcceptedOrderInfo](#interface-rawacceptedorderinfo)[]*

*Defined in [types.ts:394](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L394)*

___

###  rejected

• **rejected**: *[RawRejectedOrderInfo](#interface-rawrejectedorderinfo)[]*

*Defined in [types.ts:395](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L395)*


<hr />

# Interface: RejectedOrderInfo

## Hierarchy

* **RejectedOrderInfo**

## Index

### Properties

* [kind](#kind)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)
* [status](#status)

## Properties

###  kind

• **kind**: *[RejectedKind](#enumeration-rejectedkind)*

*Defined in [types.ts:389](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L389)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:387](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L387)*

___

###  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [types.ts:388](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L388)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:390](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L390)*


<hr />

# Interface: RejectedStatus

## Hierarchy

* **RejectedStatus**

## Index

### Properties

* [code](#code)
* [message](#message)

## Properties

###  code

• **code**: *[RejectedCode](#enumeration-rejectedcode)*

*Defined in [types.ts:375](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L375)*

___

###  message

• **message**: *string*

*Defined in [types.ts:376](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L376)*


<hr />

# Interface: StringifiedContractEvent

## Hierarchy

* **StringifiedContractEvent**

## Index

### Properties

* [address](#address)
* [blockHash](#blockhash)
* [isRemoved](#isremoved)
* [kind](#kind)
* [logIndex](#logindex)
* [parameters](#parameters)
* [txHash](#txhash)
* [txIndex](#txindex)

## Properties

###  address

• **address**: *string*

*Defined in [types.ts:253](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L253)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:248](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L248)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:252](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L252)*

___

###  kind

• **kind**: *string*

*Defined in [types.ts:254](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L254)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:251](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L251)*

___

###  parameters

• **parameters**: *[StringifiedContractEventParameters](#stringifiedcontracteventparameters)*

*Defined in [types.ts:255](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L255)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:249](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L249)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:250](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L250)*


<hr />

# Interface: StringifiedERC1155TransferBatchEvent

## Hierarchy

* **StringifiedERC1155TransferBatchEvent**

## Index

### Properties

* [from](#from)
* [ids](#ids)
* [operator](#operator)
* [to](#to)
* [values](#values)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:135](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L135)*

___

###  ids

• **ids**: *string[]*

*Defined in [types.ts:137](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L137)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:134](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L134)*

___

###  to

• **to**: *string*

*Defined in [types.ts:136](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L136)*

___

###  values

• **values**: *string[]*

*Defined in [types.ts:138](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L138)*


<hr />

# Interface: StringifiedERC1155TransferSingleEvent

## Hierarchy

* **StringifiedERC1155TransferSingleEvent**

## Index

### Properties

* [from](#from)
* [id](#id)
* [operator](#operator)
* [to](#to)
* [value](#value)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:119](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L119)*

___

###  id

• **id**: *string*

*Defined in [types.ts:121](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L121)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:118](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L118)*

___

###  to

• **to**: *string*

*Defined in [types.ts:120](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L120)*

___

###  value

• **value**: *string*

*Defined in [types.ts:122](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L122)*


<hr />

# Interface: StringifiedERC20ApprovalEvent

## Hierarchy

* **StringifiedERC20ApprovalEvent**

## Index

### Properties

* [owner](#owner)
* [spender](#spender)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:74](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L74)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:75](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L75)*

___

###  value

• **value**: *string*

*Defined in [types.ts:76](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L76)*


<hr />

# Interface: StringifiedERC20TransferEvent

## Hierarchy

* **StringifiedERC20TransferEvent**

## Index

### Properties

* [from](#from)
* [to](#to)
* [value](#value)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L62)*

___

###  to

• **to**: *string*

*Defined in [types.ts:63](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L63)*

___

###  value

• **value**: *string*

*Defined in [types.ts:64](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L64)*


<hr />

# Interface: StringifiedERC721ApprovalEvent

## Hierarchy

* **StringifiedERC721ApprovalEvent**

## Index

### Properties

* [approved](#approved)
* [owner](#owner)
* [tokenId](#tokenid)

## Properties

###  approved

• **approved**: *string*

*Defined in [types.ts:99](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L99)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L98)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:100](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L100)*


<hr />

# Interface: StringifiedERC721TransferEvent

## Hierarchy

* **StringifiedERC721TransferEvent**

## Index

### Properties

* [from](#from)
* [to](#to)
* [tokenId](#tokenid)

## Properties

###  from

• **from**: *string*

*Defined in [types.ts:86](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L86)*

___

###  to

• **to**: *string*

*Defined in [types.ts:87](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L87)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:88](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L88)*


<hr />

# Interface: StringifiedExchangeCancelUpToEvent

## Hierarchy

* **StringifiedExchangeCancelUpToEvent**

## Index

### Properties

* [makerAddress](#makeraddress)
* [orderEpoch](#orderepoch)
* [senderAddress](#senderaddress)

## Properties

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:191](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L191)*

___

###  orderEpoch

• **orderEpoch**: *string*

*Defined in [types.ts:193](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L193)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:192](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L192)*


<hr />

# Interface: StringifiedExchangeFillEvent

## Hierarchy

* **StringifiedExchangeFillEvent**

## Index

### Properties

* [feeRecipientAddress](#feerecipientaddress)
* [makerAddress](#makeraddress)
* [makerAssetData](#makerassetdata)
* [makerAssetFilledAmount](#makerassetfilledamount)
* [makerFeePaid](#makerfeepaid)
* [orderHash](#orderhash)
* [senderAddress](#senderaddress)
* [takerAddress](#takeraddress)
* [takerAssetData](#takerassetdata)
* [takerAssetFilledAmount](#takerassetfilledamount)
* [takerFeePaid](#takerfeepaid)

## Properties

###  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:165](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L165)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:162](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L162)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:171](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L171)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *string*

*Defined in [types.ts:166](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L166)*

___

###  makerFeePaid

• **makerFeePaid**: *string*

*Defined in [types.ts:168](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L168)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:170](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L170)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:164](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L164)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:163](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L163)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:172](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L172)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *string*

*Defined in [types.ts:167](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L167)*

___

###  takerFeePaid

• **takerFeePaid**: *string*

*Defined in [types.ts:169](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L169)*


<hr />

# Interface: StringifiedSignedOrder

## Hierarchy

* **StringifiedSignedOrder**

## Index

### Properties

* [exchangeAddress](#exchangeaddress)
* [expirationTimeSeconds](#expirationtimeseconds)
* [feeRecipientAddress](#feerecipientaddress)
* [makerAddress](#makeraddress)
* [makerAssetAmount](#makerassetamount)
* [makerAssetData](#makerassetdata)
* [makerFee](#makerfee)
* [salt](#salt)
* [senderAddress](#senderaddress)
* [signature](#signature)
* [takerAddress](#takeraddress)
* [takerAssetAmount](#takerassetamount)
* [takerAssetData](#takerassetdata)
* [takerFee](#takerfee)

## Properties

###  exchangeAddress

• **exchangeAddress**: *string*

*Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L49)*

___

###  expirationTimeSeconds

• **expirationTimeSeconds**: *string*

*Defined in [types.ts:51](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L51)*

___

###  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L50)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L40)*

___

###  makerAssetAmount

• **makerAssetAmount**: *string*

*Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L44)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L46)*

___

###  makerFee

• **makerFee**: *string*

*Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L42)*

___

###  salt

• **salt**: *string*

*Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L48)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L39)*

___

###  signature

• **signature**: *string*

*Defined in [types.ts:52](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L52)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L41)*

___

###  takerAssetAmount

• **takerAssetAmount**: *string*

*Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L45)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L47)*

___

###  takerFee

• **takerFee**: *string*

*Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L43)*


<hr />

# Interface: StringifiedWethDepositEvent

## Hierarchy

* **StringifiedWethDepositEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:212](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L212)*

___

###  value

• **value**: *string*

*Defined in [types.ts:213](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L213)*


<hr />

# Interface: StringifiedWethWithdrawalEvent

## Hierarchy

* **StringifiedWethWithdrawalEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:202](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L202)*

___

###  value

• **value**: *string*

*Defined in [types.ts:203](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L203)*


<hr />

# Interface: ValidationResults

## Hierarchy

* **ValidationResults**

## Index

### Properties

* [accepted](#accepted)
* [rejected](#rejected)

## Properties

###  accepted

• **accepted**: *[AcceptedOrderInfo](#interface-acceptedorderinfo)[]*

*Defined in [types.ts:399](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L399)*

___

###  rejected

• **rejected**: *[RejectedOrderInfo](#interface-rejectedorderinfo)[]*

*Defined in [types.ts:400](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L400)*


<hr />

# Interface: WethDepositEvent

## Hierarchy

* **WethDepositEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L207)*

___

###  value

• **value**: *BigNumber*

*Defined in [types.ts:208](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L208)*


<hr />

# Interface: WethWithdrawalEvent

## Hierarchy

* **WethWithdrawalEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L197)*

___

###  value

• **value**: *BigNumber*

*Defined in [types.ts:198](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L198)*


<hr />

# Interface: WSMessage

## Hierarchy

* **WSMessage**

## Index

### Properties

* [type](#type)
* [utf8Data](#utf8data)

## Properties

###  type

• **type**: *string*

*Defined in [types.ts:419](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L419)*

___

###  utf8Data

• **utf8Data**: *string*

*Defined in [types.ts:420](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L420)*


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

* **WSOpts**

## Index

### Properties

* [clientConfig](#optional-clientconfig)
* [headers](#optional-headers)
* [protocol](#optional-protocol)
* [reconnectDelay](#optional-reconnectdelay)
* [timeout](#optional-timeout)

## Properties

### `Optional` clientConfig

• **clientConfig**? : *[ClientConfig](#interface-clientconfig)*

*Defined in [types.ts:34](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L34)*

___

### `Optional` headers

• **headers**? : *undefined | __type*

*Defined in [types.ts:32](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L32)*

___

### `Optional` protocol

• **protocol**? : *undefined | string*

*Defined in [types.ts:33](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L33)*

___

### `Optional` reconnectDelay

• **reconnectDelay**? : *undefined | number*

*Defined in [types.ts:35](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L35)*

___

### `Optional` timeout

• **timeout**? : *undefined | number*

*Defined in [types.ts:31](https://github.com/0xProject/0x-mesh/blob/34c051d5/packages/rpc-client/src/types.ts#L31)*


<hr />

