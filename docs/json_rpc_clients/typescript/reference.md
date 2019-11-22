> # Class: WSClient

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
* [getStatsAsync](#getstatsasync)
* [onClose](#onclose)
* [onReconnected](#onreconnected)
* [subscribeToOrdersAsync](#subscribetoordersasync)
* [unsubscribeAsync](#unsubscribeasync)

## Constructors

###  constructor

\+ **new WSClient**(`url`: string, `wsOpts?`: [WSOpts](#interface-wsopts)): *[WSClient](#class-wsclient)*

*Defined in [ws_client.ts:241](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L241)*

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

▸ **addOrdersAsync**(`signedOrders`: `SignedOrder`[], `pinned`: boolean): *`Promise<ValidationResults>`*

*Defined in [ws_client.ts:270](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L270)*

Adds an array of 0x signed orders to the Mesh node.

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`signedOrders` | `SignedOrder`[] | - | signedOrders to add |
`pinned` | boolean | true | Whether or not the orders should be pinned. Pinned orders will not be affected by any DDoS prevention or incentive mechanisms and will always stay in storage until they are no longer fillable. |

**Returns:** *`Promise<ValidationResults>`*

validation results

___

###  destroy

▸ **destroy**(): *void*

*Defined in [ws_client.ts:385](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L385)*

destroy unsubscribes all active subscriptions, closes the websocket connection
and stops the internal heartbeat connection liveness check.

**Returns:** *void*

___

###  getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): *`Promise<OrderInfo[]>`*

*Defined in [ws_client.ts:300](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L300)*

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`perPage` | number | 200 | number of signedOrders to fetch per paginated request |

**Returns:** *`Promise<OrderInfo[]>`*

all orders, their hash and their fillableTakerAssetAmount

___

###  getStatsAsync

▸ **getStatsAsync**(): *`Promise<GetStatsResponse>`*

*Defined in [ws_client.ts:291](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L291)*

**Returns:** *`Promise<GetStatsResponse>`*

___

###  onClose

▸ **onClose**(`cb`: function): *void*

*Defined in [ws_client.ts:367](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L367)*

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

*Defined in [ws_client.ts:376](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L376)*

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: *function*

callback to call with the error when it occurs

▸ (): *void*

**Returns:** *void*

___

###  subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**(`cb`: function): *`Promise<string>`*

*Defined in [ws_client.ts:328](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L328)*

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

**Returns:** *`Promise<string>`*

subscriptionId

___

###  unsubscribeAsync

▸ **unsubscribeAsync**(`subscriptionId`: string): *`Promise<void>`*

*Defined in [ws_client.ts:357](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/ws_client.ts#L357)*

Unsubscribe from a subscription

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`subscriptionId` | string | identifier of the subscription to cancel  |

**Returns:** *`Promise<void>`*

<hr />

> # Enumeration: ContractEventKind

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

*Defined in [types.ts:221](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L221)*

___

###  ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

*Defined in [types.ts:223](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L223)*

___

###  ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

*Defined in [types.ts:222](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L222)*

___

###  ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

*Defined in [types.ts:217](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L217)*

___

###  ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

*Defined in [types.ts:216](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L216)*

___

###  ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

*Defined in [types.ts:219](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L219)*

___

###  ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

*Defined in [types.ts:220](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L220)*

___

###  ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

*Defined in [types.ts:218](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L218)*

___

###  ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

*Defined in [types.ts:225](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L225)*

___

###  ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

*Defined in [types.ts:226](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L226)*

___

###  ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

*Defined in [types.ts:224](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L224)*

___

###  WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

*Defined in [types.ts:227](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L227)*

___

###  WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

*Defined in [types.ts:228](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L228)*

<hr />

> # Enumeration: OrderEventEndState

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

*Defined in [types.ts:285](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L285)*

___

###  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [types.ts:288](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L288)*

___

###  Expired

• **Expired**: = "EXPIRED"

*Defined in [types.ts:289](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L289)*

___

###  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [types.ts:293](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L293)*

___

###  Filled

• **Filled**: = "FILLED"

*Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L286)*

___

###  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [types.ts:287](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L287)*

___

###  Invalid

• **Invalid**: = "INVALID"

*Defined in [types.ts:284](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L284)*

___

###  StoppedWatching

• **StoppedWatching**: = "STOPPED_WATCHING"

*Defined in [types.ts:291](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L291)*

___

###  Unexpired

• **Unexpired**: = "UNEXPIRED"

*Defined in [types.ts:290](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L290)*

___

###  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [types.ts:292](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L292)*

<hr />

> # Enumeration: RejectedCode

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

*Defined in [types.ts:355](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L355)*

___

###  MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MaxOrderSizeExceeded"

*Defined in [types.ts:356](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L356)*

___

###  NetworkRequestFailed

• **NetworkRequestFailed**: = "NetworkRequestFailed"

*Defined in [types.ts:359](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L359)*

___

###  OrderAlreadyStored

• **OrderAlreadyStored**: = "OrderAlreadyStored"

*Defined in [types.ts:357](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L357)*

___

###  OrderCancelled

• **OrderCancelled**: = "OrderCancelled"

*Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L364)*

___

###  OrderExpired

• **OrderExpired**: = "OrderExpired"

*Defined in [types.ts:362](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L362)*

___

###  OrderForIncorrectChain

• **OrderForIncorrectChain**: = "OrderForIncorrectChain"

*Defined in [types.ts:358](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L358)*

___

###  OrderFullyFilled

• **OrderFullyFilled**: = "OrderFullyFilled"

*Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L363)*

___

###  OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "OrderHasInvalidMakerAssetAmount"

*Defined in [types.ts:360](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L360)*

___

###  OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "OrderHasInvalidMakerAssetData"

*Defined in [types.ts:366](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L366)*

___

###  OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "OrderHasInvalidSignature"

*Defined in [types.ts:368](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L368)*

___

###  OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "OrderHasInvalidTakerAssetAmount"

*Defined in [types.ts:361](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L361)*

___

###  OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "OrderHasInvalidTakerAssetData"

*Defined in [types.ts:367](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L367)*

___

###  OrderUnfunded

• **OrderUnfunded**: = "OrderUnfunded"

*Defined in [types.ts:365](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L365)*

<hr />

> # Enumeration: RejectedKind

## Index

### Enumeration members

* [MeshError](#mesherror)
* [MeshValidation](#meshvalidation)
* [ZeroexValidation](#zeroexvalidation)

## Enumeration members

###  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [types.ts:350](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L350)*

___

###  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [types.ts:351](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L351)*

___

###  ZeroexValidation

• **ZeroexValidation**: = "ZEROEX_VALIDATION"

*Defined in [types.ts:349](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L349)*

<hr />

> # Interface: AcceptedOrderInfo

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

• **fillableTakerAssetAmount**: *`BigNumber`*

*Defined in [types.ts:332](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L332)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:333](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L333)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:330](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L330)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:331](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L331)*

<hr />

> # Interface: ClientConfig

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

*Defined in [types.ts:15](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L15)*

___

### `Optional` closeTimeout

• **closeTimeout**? : *undefined | number*

*Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L16)*

___

### `Optional` fragmentOutgoingMessages

• **fragmentOutgoingMessages**? : *undefined | false | true*

*Defined in [types.ts:13](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L13)*

___

### `Optional` fragmentationThreshold

• **fragmentationThreshold**? : *undefined | number*

*Defined in [types.ts:14](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L14)*

___

### `Optional` maxReceivedFrameSize

• **maxReceivedFrameSize**? : *undefined | number*

*Defined in [types.ts:11](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L11)*

___

### `Optional` maxReceivedMessageSize

• **maxReceivedMessageSize**? : *undefined | number*

*Defined in [types.ts:12](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L12)*

___

### `Optional` tlsOptions

• **tlsOptions**? : *any*

*Defined in [types.ts:17](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L17)*

___

### `Optional` webSocketVersion

• **webSocketVersion**? : *undefined | number*

*Defined in [types.ts:10](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L10)*

<hr />

> # Interface: ContractEvent

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

*Defined in [types.ts:278](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L278)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:273](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L273)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:277](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L277)*

___

###  kind

• **kind**: *[ContractEventKind](#enumeration-contracteventkind)*

*Defined in [types.ts:279](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L279)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:276](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L276)*

___

###  parameters

• **parameters**: *[ContractEventParameters](#contracteventparameters)*

*Defined in [types.ts:280](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L280)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:274](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L274)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:275](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L275)*

<hr />

> # Interface: ERC1155ApprovalForAllEvent

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

*Defined in [types.ts:143](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L143)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:142](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L142)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:141](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L141)*

<hr />

> # Interface: ERC1155TransferBatchEvent

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

*Defined in [types.ts:126](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L126)*

___

###  ids

• **ids**: *`BigNumber`[]*

*Defined in [types.ts:128](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L128)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:125](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L125)*

___

###  to

• **to**: *string*

*Defined in [types.ts:127](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L127)*

___

###  values

• **values**: *`BigNumber`[]*

*Defined in [types.ts:129](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L129)*

<hr />

> # Interface: ERC1155TransferSingleEvent

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

*Defined in [types.ts:110](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L110)*

___

###  id

• **id**: *`BigNumber`*

*Defined in [types.ts:112](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L112)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:109](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L109)*

___

###  to

• **to**: *string*

*Defined in [types.ts:111](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L111)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:113](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L113)*

<hr />

> # Interface: ERC20ApprovalEvent

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

*Defined in [types.ts:67](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L67)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L68)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:69](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L69)*

<hr />

> # Interface: ERC20TransferEvent

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

*Defined in [types.ts:55](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L55)*

___

###  to

• **to**: *string*

*Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L56)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:57](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L57)*

<hr />

> # Interface: ERC721ApprovalEvent

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

*Defined in [types.ts:92](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L92)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:91](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L91)*

___

###  tokenId

• **tokenId**: *`BigNumber`*

*Defined in [types.ts:93](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L93)*

<hr />

> # Interface: ERC721ApprovalForAllEvent

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

*Defined in [types.ts:105](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L105)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L104)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:103](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L103)*

<hr />

> # Interface: ERC721TransferEvent

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

*Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L79)*

___

###  to

• **to**: *string*

*Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L80)*

___

###  tokenId

• **tokenId**: *`BigNumber`*

*Defined in [types.ts:81](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L81)*

<hr />

> # Interface: ExchangeCancelEvent

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

*Defined in [types.ts:177](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L177)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:175](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L175)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:179](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L179)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:178](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L178)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:176](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L176)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:180](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L180)*

<hr />

> # Interface: ExchangeCancelUpToEvent

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

*Defined in [types.ts:184](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L184)*

___

###  orderEpoch

• **orderEpoch**: *`BigNumber`*

*Defined in [types.ts:186](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L186)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L185)*

<hr />

> # Interface: ExchangeFillEvent

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

*Defined in [types.ts:150](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L150)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:147](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L147)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:156](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L156)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *`BigNumber`*

*Defined in [types.ts:151](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L151)*

___

###  makerFeePaid

• **makerFeePaid**: *`BigNumber`*

*Defined in [types.ts:153](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L153)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:155](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L155)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:149](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L149)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:148](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L148)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:157](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L157)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *`BigNumber`*

*Defined in [types.ts:152](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L152)*

___

###  takerFeePaid

• **takerFeePaid**: *`BigNumber`*

*Defined in [types.ts:154](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L154)*

<hr />

> # Interface: GetOrdersResponse

## Hierarchy

* **GetOrdersResponse**

## Index

### Properties

* [ordersInfos](#ordersinfos)
* [snapshotID](#snapshotid)

## Properties

###  ordersInfos

• **ordersInfos**: *[RawAcceptedOrderInfo](#interface-rawacceptedorderinfo)[]*

*Defined in [types.ts:402](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L402)*

___

###  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:401](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L401)*

<hr />

> # Interface: GetStatsResponse

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

*Defined in [types.ts:429](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L429)*

___

###  ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: *number*

*Defined in [types.ts:428](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L428)*

___

###  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:420](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L420)*

___

###  latestBlock

• **latestBlock**: *[LatestBlock](#interface-latestblock)*

*Defined in [types.ts:421](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L421)*

___

###  maxExpirationTime

• **maxExpirationTime**: *string*

*Defined in [types.ts:426](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L426)*

___

###  numOrders

• **numOrders**: *number*

*Defined in [types.ts:423](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L423)*

___

###  numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: *number*

*Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L424)*

___

###  numPeers

• **numPeers**: *number*

*Defined in [types.ts:422](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L422)*

___

###  numPinnedOrders

• **numPinnedOrders**: *number*

*Defined in [types.ts:425](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L425)*

___

###  peerID

• **peerID**: *string*

*Defined in [types.ts:419](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L419)*

___

###  pubSubTopic

• **pubSubTopic**: *string*

*Defined in [types.ts:417](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L417)*

___

###  rendezvous

• **rendezvous**: *string*

*Defined in [types.ts:418](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L418)*

___

###  startOfCurrentUTCDay

• **startOfCurrentUTCDay**: *string*

*Defined in [types.ts:427](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L427)*

___

###  version

• **version**: *string*

*Defined in [types.ts:416](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L416)*

<hr />

> # Interface: HeartbeatEventPayload

## Hierarchy

* **HeartbeatEventPayload**

## Index

### Properties

* [result](#result)
* [subscription](#subscription)

## Properties

###  result

• **result**: *string*

*Defined in [types.ts:303](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L303)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:302](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L302)*

<hr />

> # Interface: LatestBlock

## Hierarchy

* **LatestBlock**

## Index

### Properties

* [hash](#hash)
* [number](#number)

## Properties

###  hash

• **hash**: *string*

*Defined in [types.ts:412](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L412)*

___

###  number

• **number**: *number*

*Defined in [types.ts:411](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L411)*

<hr />

> # Interface: OrderEvent

## Hierarchy

* **OrderEvent**

## Index

### Properties

* [contractEvents](#contractevents)
* [endState](#endstate)
* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)

## Properties

###  contractEvents

• **contractEvents**: *[ContractEvent](#interface-contractevent)[]*

*Defined in [types.ts:319](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L319)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:317](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L317)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *`BigNumber`*

*Defined in [types.ts:318](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L318)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:315](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L315)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:316](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L316)*

<hr />

> # Interface: OrderEventPayload

## Hierarchy

* **OrderEventPayload**

## Index

### Properties

* [result](#result)
* [subscription](#subscription)

## Properties

###  result

• **result**: *[RawOrderEvent](#interface-raworderevent)[]*

*Defined in [types.ts:298](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L298)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:297](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L297)*

<hr />

> # Interface: OrderInfo

## Hierarchy

* **OrderInfo**

## Index

### Properties

* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)

## Properties

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *`BigNumber`*

*Defined in [types.ts:345](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L345)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:343](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L343)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:344](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L344)*

<hr />

> # Interface: RawAcceptedOrderInfo

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

*Defined in [types.ts:325](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L325)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:326](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L326)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:323](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L323)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:324](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L324)*

<hr />

> # Interface: RawOrderEvent

## Hierarchy

* **RawOrderEvent**

## Index

### Properties

* [contractEvents](#contractevents)
* [endState](#endstate)
* [fillableTakerAssetAmount](#fillabletakerassetamount)
* [orderHash](#orderhash)
* [signedOrder](#signedorder)

## Properties

###  contractEvents

• **contractEvents**: *[StringifiedContractEvent](#interface-stringifiedcontractevent)[]*

*Defined in [types.ts:311](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L311)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:309](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L309)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:310](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L310)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:307](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L307)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:308](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L308)*

<hr />

> # Interface: RawOrderInfo

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

*Defined in [types.ts:339](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L339)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:337](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L337)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:338](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L338)*

<hr />

> # Interface: RawRejectedOrderInfo

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

*Defined in [types.ts:379](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L379)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:377](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L377)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:378](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L378)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:380](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L380)*

<hr />

> # Interface: RawValidationResults

## Hierarchy

* **RawValidationResults**

## Index

### Properties

* [accepted](#accepted)
* [rejected](#rejected)

## Properties

###  accepted

• **accepted**: *[RawAcceptedOrderInfo](#interface-rawacceptedorderinfo)[]*

*Defined in [types.ts:391](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L391)*

___

###  rejected

• **rejected**: *[RawRejectedOrderInfo](#interface-rawrejectedorderinfo)[]*

*Defined in [types.ts:392](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L392)*

<hr />

> # Interface: RejectedOrderInfo

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

*Defined in [types.ts:386](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L386)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:384](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L384)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:385](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L385)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:387](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L387)*

<hr />

> # Interface: RejectedStatus

## Hierarchy

* **RejectedStatus**

## Index

### Properties

* [code](#code)
* [message](#message)

## Properties

###  code

• **code**: *[RejectedCode](#enumeration-rejectedcode)*

*Defined in [types.ts:372](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L372)*

___

###  message

• **message**: *string*

*Defined in [types.ts:373](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L373)*

<hr />

> # Interface: StringifiedContractEvent

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

*Defined in [types.ts:252](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L252)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:247](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L247)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:251](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L251)*

___

###  kind

• **kind**: *string*

*Defined in [types.ts:253](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L253)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:250](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L250)*

___

###  parameters

• **parameters**: *[StringifiedContractEventParameters](#stringifiedcontracteventparameters)*

*Defined in [types.ts:254](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L254)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:248](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L248)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:249](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L249)*

<hr />

> # Interface: StringifiedERC1155TransferBatchEvent

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

*Defined in [types.ts:134](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L134)*

___

###  ids

• **ids**: *string[]*

*Defined in [types.ts:136](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L136)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:133](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L133)*

___

###  to

• **to**: *string*

*Defined in [types.ts:135](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L135)*

___

###  values

• **values**: *string[]*

*Defined in [types.ts:137](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L137)*

<hr />

> # Interface: StringifiedERC1155TransferSingleEvent

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

*Defined in [types.ts:118](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L118)*

___

###  id

• **id**: *string*

*Defined in [types.ts:120](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L120)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:117](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L117)*

___

###  to

• **to**: *string*

*Defined in [types.ts:119](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L119)*

___

###  value

• **value**: *string*

*Defined in [types.ts:121](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L121)*

<hr />

> # Interface: StringifiedERC20ApprovalEvent

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

*Defined in [types.ts:73](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L73)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:74](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L74)*

___

###  value

• **value**: *string*

*Defined in [types.ts:75](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L75)*

<hr />

> # Interface: StringifiedERC20TransferEvent

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

*Defined in [types.ts:61](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L61)*

___

###  to

• **to**: *string*

*Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L62)*

___

###  value

• **value**: *string*

*Defined in [types.ts:63](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L63)*

<hr />

> # Interface: StringifiedERC721ApprovalEvent

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

*Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L98)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:97](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L97)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:99](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L99)*

<hr />

> # Interface: StringifiedERC721TransferEvent

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

*Defined in [types.ts:85](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L85)*

___

###  to

• **to**: *string*

*Defined in [types.ts:86](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L86)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:87](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L87)*

<hr />

> # Interface: StringifiedExchangeCancelUpToEvent

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

*Defined in [types.ts:190](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L190)*

___

###  orderEpoch

• **orderEpoch**: *string*

*Defined in [types.ts:192](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L192)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:191](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L191)*

<hr />

> # Interface: StringifiedExchangeFillEvent

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

*Defined in [types.ts:164](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L164)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:161](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L161)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:170](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L170)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *string*

*Defined in [types.ts:165](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L165)*

___

###  makerFeePaid

• **makerFeePaid**: *string*

*Defined in [types.ts:167](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L167)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:169](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L169)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:163](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L163)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:162](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L162)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:171](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L171)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *string*

*Defined in [types.ts:166](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L166)*

___

###  takerFeePaid

• **takerFeePaid**: *string*

*Defined in [types.ts:168](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L168)*

<hr />

> # Interface: StringifiedSignedOrder

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

*Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L48)*

___

###  expirationTimeSeconds

• **expirationTimeSeconds**: *string*

*Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L50)*

___

###  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L49)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L39)*

___

###  makerAssetAmount

• **makerAssetAmount**: *string*

*Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L43)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L45)*

___

###  makerFee

• **makerFee**: *string*

*Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L41)*

___

###  salt

• **salt**: *string*

*Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L47)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:38](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L38)*

___

###  signature

• **signature**: *string*

*Defined in [types.ts:51](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L51)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L40)*

___

###  takerAssetAmount

• **takerAssetAmount**: *string*

*Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L44)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L46)*

___

###  takerFee

• **takerFee**: *string*

*Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L42)*

<hr />

> # Interface: StringifiedWethDepositEvent

## Hierarchy

* **StringifiedWethDepositEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:211](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L211)*

___

###  value

• **value**: *string*

*Defined in [types.ts:212](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L212)*

<hr />

> # Interface: StringifiedWethWithdrawalEvent

## Hierarchy

* **StringifiedWethWithdrawalEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:201](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L201)*

___

###  value

• **value**: *string*

*Defined in [types.ts:202](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L202)*

<hr />

> # Interface: ValidationResults

## Hierarchy

* **ValidationResults**

## Index

### Properties

* [accepted](#accepted)
* [rejected](#rejected)

## Properties

###  accepted

• **accepted**: *[AcceptedOrderInfo](#interface-acceptedorderinfo)[]*

*Defined in [types.ts:396](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L396)*

___

###  rejected

• **rejected**: *[RejectedOrderInfo](#interface-rejectedorderinfo)[]*

*Defined in [types.ts:397](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L397)*

<hr />

> # Interface: WethDepositEvent

## Hierarchy

* **WethDepositEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:206](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L206)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L207)*

<hr />

> # Interface: WethWithdrawalEvent

## Hierarchy

* **WethWithdrawalEvent**

## Index

### Properties

* [owner](#owner)
* [value](#value)

## Properties

###  owner

• **owner**: *string*

*Defined in [types.ts:196](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L196)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L197)*

<hr />

> # Interface: WSMessage

## Hierarchy

* **WSMessage**

## Index

### Properties

* [type](#type)
* [utf8Data](#utf8data)

## Properties

###  type

• **type**: *string*

*Defined in [types.ts:406](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L406)*

___

###  utf8Data

• **utf8Data**: *string*

*Defined in [types.ts:407](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L407)*

<hr />

> # Interface: WSOpts

timeout: timeout in milliseconds to enforce on every WS request that expects a response
headers: Request headers (e.g., authorization)
protocol: requestOptions should be either null or an object specifying additional configuration options to be
passed to http.request or https.request. This can be used to pass a custom agent to enable WebSocketClient usage
from behind an HTTP or HTTPS proxy server using koichik/node-tunnel or similar.
clientConfig: The client configs documented here: https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md
reconnectAfter: time in milliseconds after which to attempt to reconnect to WS server after an error occurred (default: 5000)

## Hierarchy

* **WSOpts**

## Index

### Properties

* [clientConfig](#optional-clientconfig)
* [headers](#optional-headers)
* [protocol](#optional-protocol)
* [reconnectAfter](#optional-reconnectafter)
* [timeout](#optional-timeout)

## Properties

### `Optional` clientConfig

• **clientConfig**? : *[ClientConfig](#interface-clientconfig)*

*Defined in [types.ts:33](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L33)*

___

### `Optional` headers

• **headers**? : *undefined | `__type`*

*Defined in [types.ts:31](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L31)*

___

### `Optional` protocol

• **protocol**? : *undefined | string*

*Defined in [types.ts:32](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L32)*

___

### `Optional` reconnectAfter

• **reconnectAfter**? : *undefined | number*

*Defined in [types.ts:34](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L34)*

___

### `Optional` timeout

• **timeout**? : *undefined | number*

*Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/c1e1ff7/rpc/clients/typescript/src/types.ts#L30)*

<hr />

