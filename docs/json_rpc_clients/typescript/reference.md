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

*Defined in [ws_client.ts:222](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L222)*

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

*Defined in [ws_client.ts:251](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L251)*

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

*Defined in [ws_client.ts:366](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L366)*

destroy unsubscribes all active subscriptions, closes the websocket connection
and stops the internal heartbeat connection liveness check.

**Returns:** *void*

___

###  getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): *`Promise<OrderInfo[]>`*

*Defined in [ws_client.ts:281](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L281)*

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

*Defined in [ws_client.ts:272](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L272)*

**Returns:** *`Promise<GetStatsResponse>`*

___

###  onClose

▸ **onClose**(`cb`: function): *void*

*Defined in [ws_client.ts:348](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L348)*

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

*Defined in [ws_client.ts:357](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L357)*

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: *function*

callback to call with the error when it occurs

▸ (): *void*

**Returns:** *void*

___

###  subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**(`cb`: function): *`Promise<string>`*

*Defined in [ws_client.ts:309](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L309)*

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

*Defined in [ws_client.ts:338](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/ws_client.ts#L338)*

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

*Defined in [types.ts:220](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L220)*

___

###  ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

*Defined in [types.ts:222](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L222)*

___

###  ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

*Defined in [types.ts:221](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L221)*

___

###  ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

*Defined in [types.ts:216](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L216)*

___

###  ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

*Defined in [types.ts:215](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L215)*

___

###  ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

*Defined in [types.ts:218](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L218)*

___

###  ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

*Defined in [types.ts:219](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L219)*

___

###  ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

*Defined in [types.ts:217](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L217)*

___

###  ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

*Defined in [types.ts:224](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L224)*

___

###  ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

*Defined in [types.ts:225](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L225)*

___

###  ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

*Defined in [types.ts:223](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L223)*

___

###  WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

*Defined in [types.ts:226](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L226)*

___

###  WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

*Defined in [types.ts:227](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L227)*

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
* [Unfunded](#unfunded)

## Enumeration members

###  Added

• **Added**: = "ADDED"

*Defined in [types.ts:284](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L284)*

___

###  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [types.ts:287](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L287)*

___

###  Expired

• **Expired**: = "EXPIRED"

*Defined in [types.ts:288](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L288)*

___

###  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [types.ts:290](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L290)*

___

###  Filled

• **Filled**: = "FILLED"

*Defined in [types.ts:285](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L285)*

___

###  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L286)*

___

###  Invalid

• **Invalid**: = "INVALID"

*Defined in [types.ts:283](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L283)*

___

###  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [types.ts:289](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L289)*

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

*Defined in [types.ts:352](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L352)*

___

###  MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MaxOrderSizeExceeded"

*Defined in [types.ts:353](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L353)*

___

###  NetworkRequestFailed

• **NetworkRequestFailed**: = "NetworkRequestFailed"

*Defined in [types.ts:356](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L356)*

___

###  OrderAlreadyStored

• **OrderAlreadyStored**: = "OrderAlreadyStored"

*Defined in [types.ts:354](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L354)*

___

###  OrderCancelled

• **OrderCancelled**: = "OrderCancelled"

*Defined in [types.ts:361](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L361)*

___

###  OrderExpired

• **OrderExpired**: = "OrderExpired"

*Defined in [types.ts:359](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L359)*

___

###  OrderForIncorrectChain

• **OrderForIncorrectChain**: = "OrderForIncorrectChain"

*Defined in [types.ts:355](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L355)*

___

###  OrderFullyFilled

• **OrderFullyFilled**: = "OrderFullyFilled"

*Defined in [types.ts:360](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L360)*

___

###  OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "OrderHasInvalidMakerAssetAmount"

*Defined in [types.ts:357](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L357)*

___

###  OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "OrderHasInvalidMakerAssetData"

*Defined in [types.ts:363](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L363)*

___

###  OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "OrderHasInvalidSignature"

*Defined in [types.ts:365](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L365)*

___

###  OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "OrderHasInvalidTakerAssetAmount"

*Defined in [types.ts:358](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L358)*

___

###  OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "OrderHasInvalidTakerAssetData"

*Defined in [types.ts:364](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L364)*

___

###  OrderUnfunded

• **OrderUnfunded**: = "OrderUnfunded"

*Defined in [types.ts:362](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L362)*

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

*Defined in [types.ts:347](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L347)*

___

###  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [types.ts:348](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L348)*

___

###  ZeroexValidation

• **ZeroexValidation**: = "ZEROEX_VALIDATION"

*Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L346)*

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

*Defined in [types.ts:329](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L329)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:330](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L330)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:327](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L327)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:328](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L328)*

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

*Defined in [types.ts:14](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L14)*

___

### `Optional` closeTimeout

• **closeTimeout**? : *undefined | number*

*Defined in [types.ts:15](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L15)*

___

### `Optional` fragmentOutgoingMessages

• **fragmentOutgoingMessages**? : *undefined | false | true*

*Defined in [types.ts:12](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L12)*

___

### `Optional` fragmentationThreshold

• **fragmentationThreshold**? : *undefined | number*

*Defined in [types.ts:13](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L13)*

___

### `Optional` maxReceivedFrameSize

• **maxReceivedFrameSize**? : *undefined | number*

*Defined in [types.ts:10](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L10)*

___

### `Optional` maxReceivedMessageSize

• **maxReceivedMessageSize**? : *undefined | number*

*Defined in [types.ts:11](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L11)*

___

### `Optional` tlsOptions

• **tlsOptions**? : *any*

*Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L16)*

___

### `Optional` webSocketVersion

• **webSocketVersion**? : *undefined | number*

*Defined in [types.ts:9](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L9)*

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

*Defined in [types.ts:277](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L277)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:272](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L272)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:276](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L276)*

___

###  kind

• **kind**: *[ContractEventKind](#enumeration-contracteventkind)*

*Defined in [types.ts:278](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L278)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:275](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L275)*

___

###  parameters

• **parameters**: *[ContractEventParameters](#contracteventparameters)*

*Defined in [types.ts:279](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L279)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:273](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L273)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:274](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L274)*

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

*Defined in [types.ts:142](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L142)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:141](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L141)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:140](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L140)*

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

*Defined in [types.ts:125](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L125)*

___

###  ids

• **ids**: *`BigNumber`[]*

*Defined in [types.ts:127](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L127)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:124](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L124)*

___

###  to

• **to**: *string*

*Defined in [types.ts:126](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L126)*

___

###  values

• **values**: *`BigNumber`[]*

*Defined in [types.ts:128](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L128)*

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

*Defined in [types.ts:109](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L109)*

___

###  id

• **id**: *`BigNumber`*

*Defined in [types.ts:111](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L111)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:108](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L108)*

___

###  to

• **to**: *string*

*Defined in [types.ts:110](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L110)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:112](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L112)*

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

*Defined in [types.ts:66](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L66)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:67](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L67)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L68)*

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

*Defined in [types.ts:54](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L54)*

___

###  to

• **to**: *string*

*Defined in [types.ts:55](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L55)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L56)*

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

*Defined in [types.ts:91](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L91)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:90](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L90)*

___

###  tokenId

• **tokenId**: *`BigNumber`*

*Defined in [types.ts:92](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L92)*

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

*Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L104)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:103](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L103)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:102](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L102)*

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

*Defined in [types.ts:78](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L78)*

___

###  to

• **to**: *string*

*Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L79)*

___

###  tokenId

• **tokenId**: *`BigNumber`*

*Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L80)*

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

*Defined in [types.ts:176](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L176)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:174](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L174)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:178](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L178)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:177](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L177)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:175](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L175)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:179](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L179)*

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

*Defined in [types.ts:183](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L183)*

___

###  orderEpoch

• **orderEpoch**: *`BigNumber`*

*Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L185)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:184](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L184)*

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

*Defined in [types.ts:149](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L149)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:146](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L146)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:155](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L155)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *`BigNumber`*

*Defined in [types.ts:150](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L150)*

___

###  makerFeePaid

• **makerFeePaid**: *`BigNumber`*

*Defined in [types.ts:152](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L152)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:154](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L154)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:148](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L148)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:147](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L147)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:156](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L156)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *`BigNumber`*

*Defined in [types.ts:151](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L151)*

___

###  takerFeePaid

• **takerFeePaid**: *`BigNumber`*

*Defined in [types.ts:153](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L153)*

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

*Defined in [types.ts:399](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L399)*

___

###  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:398](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L398)*

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

*Defined in [types.ts:426](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L426)*

___

###  ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: *number*

*Defined in [types.ts:425](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L425)*

___

###  ethereumChainID

• **ethereumChainID**: *number*

*Defined in [types.ts:417](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L417)*

___

###  latestBlock

• **latestBlock**: *[LatestBlock](#interface-latestblock)*

*Defined in [types.ts:418](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L418)*

___

###  maxExpirationTime

• **maxExpirationTime**: *string*

*Defined in [types.ts:423](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L423)*

___

###  numOrders

• **numOrders**: *number*

*Defined in [types.ts:420](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L420)*

___

###  numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: *number*

*Defined in [types.ts:421](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L421)*

___

###  numPeers

• **numPeers**: *number*

*Defined in [types.ts:419](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L419)*

___

###  numPinnedOrders

• **numPinnedOrders**: *number*

*Defined in [types.ts:422](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L422)*

___

###  peerID

• **peerID**: *string*

*Defined in [types.ts:416](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L416)*

___

###  pubSubTopic

• **pubSubTopic**: *string*

*Defined in [types.ts:414](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L414)*

___

###  rendezvous

• **rendezvous**: *string*

*Defined in [types.ts:415](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L415)*

___

###  startOfCurrentUTCDay

• **startOfCurrentUTCDay**: *string*

*Defined in [types.ts:424](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L424)*

___

###  version

• **version**: *string*

*Defined in [types.ts:413](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L413)*

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

*Defined in [types.ts:300](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L300)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:299](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L299)*

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

*Defined in [types.ts:409](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L409)*

___

###  number

• **number**: *number*

*Defined in [types.ts:408](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L408)*

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

*Defined in [types.ts:316](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L316)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:314](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L314)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *`BigNumber`*

*Defined in [types.ts:315](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L315)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:312](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L312)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:313](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L313)*

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

*Defined in [types.ts:295](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L295)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:294](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L294)*

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

*Defined in [types.ts:342](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L342)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:340](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L340)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:341](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L341)*

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

*Defined in [types.ts:322](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L322)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:323](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L323)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:320](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L320)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:321](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L321)*

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

*Defined in [types.ts:308](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L308)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:306](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L306)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:307](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L307)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:304](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L304)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:305](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L305)*

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

*Defined in [types.ts:336](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L336)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:334](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L334)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:335](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L335)*

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

*Defined in [types.ts:376](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L376)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:374](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L374)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:375](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L375)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:377](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L377)*

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

*Defined in [types.ts:388](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L388)*

___

###  rejected

• **rejected**: *[RawRejectedOrderInfo](#interface-rawrejectedorderinfo)[]*

*Defined in [types.ts:389](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L389)*

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

*Defined in [types.ts:383](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L383)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:381](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L381)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:382](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L382)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:384](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L384)*

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

*Defined in [types.ts:369](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L369)*

___

###  message

• **message**: *string*

*Defined in [types.ts:370](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L370)*

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

*Defined in [types.ts:251](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L251)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:246](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L246)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:250](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L250)*

___

###  kind

• **kind**: *string*

*Defined in [types.ts:252](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L252)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:249](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L249)*

___

###  parameters

• **parameters**: *[StringifiedContractEventParameters](#stringifiedcontracteventparameters)*

*Defined in [types.ts:253](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L253)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:247](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L247)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:248](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L248)*

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

*Defined in [types.ts:133](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L133)*

___

###  ids

• **ids**: *string[]*

*Defined in [types.ts:135](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L135)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:132](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L132)*

___

###  to

• **to**: *string*

*Defined in [types.ts:134](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L134)*

___

###  values

• **values**: *string[]*

*Defined in [types.ts:136](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L136)*

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

*Defined in [types.ts:117](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L117)*

___

###  id

• **id**: *string*

*Defined in [types.ts:119](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L119)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:116](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L116)*

___

###  to

• **to**: *string*

*Defined in [types.ts:118](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L118)*

___

###  value

• **value**: *string*

*Defined in [types.ts:120](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L120)*

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

*Defined in [types.ts:72](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L72)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:73](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L73)*

___

###  value

• **value**: *string*

*Defined in [types.ts:74](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L74)*

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

*Defined in [types.ts:60](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L60)*

___

###  to

• **to**: *string*

*Defined in [types.ts:61](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L61)*

___

###  value

• **value**: *string*

*Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L62)*

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

*Defined in [types.ts:97](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L97)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:96](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L96)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L98)*

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

*Defined in [types.ts:84](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L84)*

___

###  to

• **to**: *string*

*Defined in [types.ts:85](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L85)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:86](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L86)*

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

*Defined in [types.ts:189](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L189)*

___

###  orderEpoch

• **orderEpoch**: *string*

*Defined in [types.ts:191](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L191)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:190](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L190)*

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

*Defined in [types.ts:163](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L163)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:160](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L160)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:169](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L169)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *string*

*Defined in [types.ts:164](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L164)*

___

###  makerFeePaid

• **makerFeePaid**: *string*

*Defined in [types.ts:166](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L166)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:168](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L168)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:162](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L162)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:161](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L161)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:170](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L170)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *string*

*Defined in [types.ts:165](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L165)*

___

###  takerFeePaid

• **takerFeePaid**: *string*

*Defined in [types.ts:167](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L167)*

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

*Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L47)*

___

###  expirationTimeSeconds

• **expirationTimeSeconds**: *string*

*Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L49)*

___

###  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L48)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:38](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L38)*

___

###  makerAssetAmount

• **makerAssetAmount**: *string*

*Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L42)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L44)*

___

###  makerFee

• **makerFee**: *string*

*Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L40)*

___

###  salt

• **salt**: *string*

*Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L46)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:37](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L37)*

___

###  signature

• **signature**: *string*

*Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L50)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L39)*

___

###  takerAssetAmount

• **takerAssetAmount**: *string*

*Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L43)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L45)*

___

###  takerFee

• **takerFee**: *string*

*Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L41)*

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

*Defined in [types.ts:210](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L210)*

___

###  value

• **value**: *string*

*Defined in [types.ts:211](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L211)*

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

*Defined in [types.ts:200](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L200)*

___

###  value

• **value**: *string*

*Defined in [types.ts:201](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L201)*

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

*Defined in [types.ts:393](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L393)*

___

###  rejected

• **rejected**: *[RejectedOrderInfo](#interface-rejectedorderinfo)[]*

*Defined in [types.ts:394](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L394)*

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

*Defined in [types.ts:205](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L205)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:206](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L206)*

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

*Defined in [types.ts:195](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L195)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:196](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L196)*

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

*Defined in [types.ts:403](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L403)*

___

###  utf8Data

• **utf8Data**: *string*

*Defined in [types.ts:404](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L404)*

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

*Defined in [types.ts:32](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L32)*

___

### `Optional` headers

• **headers**? : *undefined | `__type`*

*Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L30)*

___

### `Optional` protocol

• **protocol**? : *undefined | string*

*Defined in [types.ts:31](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L31)*

___

### `Optional` reconnectAfter

• **reconnectAfter**? : *undefined | number*

*Defined in [types.ts:33](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L33)*

___

### `Optional` timeout

• **timeout**? : *undefined | number*

*Defined in [types.ts:29](https://github.com/0xProject/0x-mesh/blob/7d5a102/rpc/clients/typescript/src/types.ts#L29)*

<hr />

