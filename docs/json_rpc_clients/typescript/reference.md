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

*Defined in [ws_client.ts:184](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L184)*

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

▸ **addOrdersAsync**(`signedOrders`: `SignedOrder`[]): *`Promise<ValidationResults>`*

*Defined in [ws_client.ts:209](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L209)*

Adds an array of 0x signed orders to the Mesh node.

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`signedOrders` | `SignedOrder`[] | signedOrders to add |

**Returns:** *`Promise<ValidationResults>`*

validation results

___

###  destroy

▸ **destroy**(): *void*

*Defined in [ws_client.ts:326](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L326)*

destroy unsubscribes all active subscriptions, closes the websocket connection
and stops the internal heartbeat connection liveness check.

**Returns:** *void*

___

###  getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): *`Promise<OrderInfo[]>`*

*Defined in [ws_client.ts:238](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L238)*

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

*Defined in [ws_client.ts:229](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L229)*

**Returns:** *`Promise<GetStatsResponse>`*

___

###  onClose

▸ **onClose**(`cb`: function): *void*

*Defined in [ws_client.ts:308](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L308)*

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

*Defined in [ws_client.ts:317](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L317)*

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: *function*

callback to call with the error when it occurs

▸ (): *void*

**Returns:** *void*

___

###  subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**(`cb`: function): *`Promise<string>`*

*Defined in [ws_client.ts:269](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L269)*

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

*Defined in [ws_client.ts:298](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/ws_client.ts#L298)*

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

* [ERC20ApprovalEvent](#erc20approvalevent)
* [ERC20TransferEvent](#erc20transferevent)
* [ERC721ApprovalEvent](#erc721approvalevent)
* [ERC721TransferEvent](#erc721transferevent)
* [ExchangeCancelEvent](#exchangecancelevent)
* [ExchangeCancelUpToEvent](#exchangecanceluptoevent)
* [ExchangeFillEvent](#exchangefillevent)
* [WethDepositEvent](#wethdepositevent)
* [WethWithdrawalEvent](#wethwithdrawalevent)

## Enumeration members

###  ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

*Defined in [types.ts:178](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L178)*

___

###  ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

*Defined in [types.ts:177](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L177)*

___

###  ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

*Defined in [types.ts:180](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L180)*

___

###  ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

*Defined in [types.ts:179](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L179)*

___

###  ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

*Defined in [types.ts:182](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L182)*

___

###  ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

*Defined in [types.ts:183](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L183)*

___

###  ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

*Defined in [types.ts:181](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L181)*

___

###  WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

*Defined in [types.ts:184](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L184)*

___

###  WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

*Defined in [types.ts:185](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L185)*

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

*Defined in [types.ts:216](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L216)*

___

###  Cancelled

• **Cancelled**: = "CANCELLED"

*Defined in [types.ts:219](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L219)*

___

###  Expired

• **Expired**: = "EXPIRED"

*Defined in [types.ts:220](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L220)*

___

###  FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY_INCREASED"

*Defined in [types.ts:222](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L222)*

___

###  Filled

• **Filled**: = "FILLED"

*Defined in [types.ts:217](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L217)*

___

###  FullyFilled

• **FullyFilled**: = "FULLY_FILLED"

*Defined in [types.ts:218](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L218)*

___

###  Invalid

• **Invalid**: = "INVALID"

*Defined in [types.ts:215](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L215)*

___

###  Unfunded

• **Unfunded**: = "UNFUNDED"

*Defined in [types.ts:221](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L221)*

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
* [OrderForIncorrectNetwork](#orderforincorrectnetwork)
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

*Defined in [types.ts:284](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L284)*

___

###  MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MaxOrderSizeExceeded"

*Defined in [types.ts:285](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L285)*

___

###  NetworkRequestFailed

• **NetworkRequestFailed**: = "NetworkRequestFailed"

*Defined in [types.ts:288](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L288)*

___

###  OrderAlreadyStored

• **OrderAlreadyStored**: = "OrderAlreadyStored"

*Defined in [types.ts:286](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L286)*

___

###  OrderCancelled

• **OrderCancelled**: = "OrderCancelled"

*Defined in [types.ts:293](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L293)*

___

###  OrderExpired

• **OrderExpired**: = "OrderExpired"

*Defined in [types.ts:291](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L291)*

___

###  OrderForIncorrectNetwork

• **OrderForIncorrectNetwork**: = "OrderForIncorrectNetwork"

*Defined in [types.ts:287](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L287)*

___

###  OrderFullyFilled

• **OrderFullyFilled**: = "OrderFullyFilled"

*Defined in [types.ts:292](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L292)*

___

###  OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "OrderHasInvalidMakerAssetAmount"

*Defined in [types.ts:289](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L289)*

___

###  OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "OrderHasInvalidMakerAssetData"

*Defined in [types.ts:295](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L295)*

___

###  OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "OrderHasInvalidSignature"

*Defined in [types.ts:297](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L297)*

___

###  OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "OrderHasInvalidTakerAssetAmount"

*Defined in [types.ts:290](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L290)*

___

###  OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "OrderHasInvalidTakerAssetData"

*Defined in [types.ts:296](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L296)*

___

###  OrderUnfunded

• **OrderUnfunded**: = "OrderUnfunded"

*Defined in [types.ts:294](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L294)*

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

*Defined in [types.ts:279](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L279)*

___

###  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [types.ts:280](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L280)*

___

###  ZeroexValidation

• **ZeroexValidation**: = "ZEROEX_VALIDATION"

*Defined in [types.ts:278](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L278)*

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

*Defined in [types.ts:261](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L261)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:262](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L262)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:259](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L259)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:260](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L260)*

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

*Defined in [types.ts:14](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L14)*

___

### `Optional` closeTimeout

• **closeTimeout**? : *undefined | number*

*Defined in [types.ts:15](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L15)*

___

### `Optional` fragmentOutgoingMessages

• **fragmentOutgoingMessages**? : *undefined | false | true*

*Defined in [types.ts:12](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L12)*

___

### `Optional` fragmentationThreshold

• **fragmentationThreshold**? : *undefined | number*

*Defined in [types.ts:13](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L13)*

___

### `Optional` maxReceivedFrameSize

• **maxReceivedFrameSize**? : *undefined | number*

*Defined in [types.ts:10](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L10)*

___

### `Optional` maxReceivedMessageSize

• **maxReceivedMessageSize**? : *undefined | number*

*Defined in [types.ts:11](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L11)*

___

### `Optional` tlsOptions

• **tlsOptions**? : *any*

*Defined in [types.ts:16](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L16)*

___

### `Optional` webSocketVersion

• **webSocketVersion**? : *undefined | number*

*Defined in [types.ts:9](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L9)*

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

*Defined in [types.ts:209](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L209)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:204](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L204)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:208](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L208)*

___

###  kind

• **kind**: *[ContractEventKind](#enumeration-contracteventkind)*

*Defined in [types.ts:210](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L210)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:207](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L207)*

___

###  parameters

• **parameters**: *[ContractEventParameters](#contracteventparameters)*

*Defined in [types.ts:211](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L211)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:205](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L205)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:206](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L206)*

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

*Defined in [types.ts:66](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L66)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:67](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L67)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:68](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L68)*

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

*Defined in [types.ts:54](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L54)*

___

###  to

• **to**: *string*

*Defined in [types.ts:55](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L55)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:56](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L56)*

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

*Defined in [types.ts:91](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L91)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:90](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L90)*

___

###  tokenId

• **tokenId**: *`BigNumber`*

*Defined in [types.ts:92](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L92)*

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

*Defined in [types.ts:104](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L104)*

___

###  operator

• **operator**: *string*

*Defined in [types.ts:103](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L103)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:102](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L102)*

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

*Defined in [types.ts:78](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L78)*

___

###  to

• **to**: *string*

*Defined in [types.ts:79](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L79)*

___

###  tokenId

• **tokenId**: *`BigNumber`*

*Defined in [types.ts:80](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L80)*

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

*Defined in [types.ts:138](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L138)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:136](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L136)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:140](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L140)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:139](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L139)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:137](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L137)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:141](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L141)*

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

*Defined in [types.ts:145](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L145)*

___

###  orderEpoch

• **orderEpoch**: *`BigNumber`*

*Defined in [types.ts:147](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L147)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:146](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L146)*

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

*Defined in [types.ts:111](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L111)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:108](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L108)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:117](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L117)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *`BigNumber`*

*Defined in [types.ts:112](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L112)*

___

###  makerFeePaid

• **makerFeePaid**: *`BigNumber`*

*Defined in [types.ts:114](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L114)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:116](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L116)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:110](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L110)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:109](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L109)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:118](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L118)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *`BigNumber`*

*Defined in [types.ts:113](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L113)*

___

###  takerFeePaid

• **takerFeePaid**: *`BigNumber`*

*Defined in [types.ts:115](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L115)*

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

*Defined in [types.ts:331](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L331)*

___

###  snapshotID

• **snapshotID**: *string*

*Defined in [types.ts:330](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L330)*

<hr />

> # Interface: GetStatsResponse

## Hierarchy

* **GetStatsResponse**

## Index

### Properties

* [EthereumNetworkID](#ethereumnetworkid)
* [LatestBlock](#latestblock)
* [NumOrders](#numorders)
* [NumPeers](#numpeers)
* [PeerID](#peerid)
* [PubSubTopic](#pubsubtopic)
* [Rendezvous](#rendezvous)
* [Version](#version)

## Properties

###  EthereumNetworkID

• **EthereumNetworkID**: *number*

*Defined in [types.ts:349](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L349)*

___

###  LatestBlock

• **LatestBlock**: *[LatestBlock](#interface-latestblock)*

*Defined in [types.ts:350](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L350)*

___

###  NumOrders

• **NumOrders**: *number*

*Defined in [types.ts:352](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L352)*

___

###  NumPeers

• **NumPeers**: *number*

*Defined in [types.ts:351](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L351)*

___

###  PeerID

• **PeerID**: *string*

*Defined in [types.ts:348](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L348)*

___

###  PubSubTopic

• **PubSubTopic**: *string*

*Defined in [types.ts:346](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L346)*

___

###  Rendezvous

• **Rendezvous**: *string*

*Defined in [types.ts:347](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L347)*

___

###  Version

• **Version**: *string*

*Defined in [types.ts:345](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L345)*

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

*Defined in [types.ts:232](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L232)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:231](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L231)*

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

*Defined in [types.ts:341](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L341)*

___

###  number

• **number**: *number*

*Defined in [types.ts:340](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L340)*

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

*Defined in [types.ts:248](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L248)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:246](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L246)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *`BigNumber`*

*Defined in [types.ts:247](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L247)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:244](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L244)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:245](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L245)*

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

*Defined in [types.ts:227](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L227)*

___

###  subscription

• **subscription**: *string*

*Defined in [types.ts:226](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L226)*

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

*Defined in [types.ts:274](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L274)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:272](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L272)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:273](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L273)*

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

*Defined in [types.ts:254](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L254)*

___

###  isNew

• **isNew**: *boolean*

*Defined in [types.ts:255](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L255)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:252](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L252)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:253](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L253)*

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

*Defined in [types.ts:240](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L240)*

___

###  endState

• **endState**: *[OrderEventEndState](#enumeration-ordereventendstate)*

*Defined in [types.ts:238](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L238)*

___

###  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *string*

*Defined in [types.ts:239](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L239)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:236](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L236)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:237](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L237)*

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

*Defined in [types.ts:268](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L268)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:266](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L266)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:267](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L267)*

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

*Defined in [types.ts:308](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L308)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:306](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L306)*

___

###  signedOrder

• **signedOrder**: *[StringifiedSignedOrder](#interface-stringifiedsignedorder)*

*Defined in [types.ts:307](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L307)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:309](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L309)*

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

*Defined in [types.ts:320](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L320)*

___

###  rejected

• **rejected**: *[RawRejectedOrderInfo](#interface-rawrejectedorderinfo)[]*

*Defined in [types.ts:321](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L321)*

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

*Defined in [types.ts:315](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L315)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:313](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L313)*

___

###  signedOrder

• **signedOrder**: *`SignedOrder`*

*Defined in [types.ts:314](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L314)*

___

###  status

• **status**: *[RejectedStatus](#interface-rejectedstatus)*

*Defined in [types.ts:316](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L316)*

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

*Defined in [types.ts:301](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L301)*

___

###  message

• **message**: *string*

*Defined in [types.ts:302](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L302)*

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

*Defined in [types.ts:196](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L196)*

___

###  blockHash

• **blockHash**: *string*

*Defined in [types.ts:191](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L191)*

___

###  isRemoved

• **isRemoved**: *string*

*Defined in [types.ts:195](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L195)*

___

###  kind

• **kind**: *string*

*Defined in [types.ts:197](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L197)*

___

###  logIndex

• **logIndex**: *number*

*Defined in [types.ts:194](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L194)*

___

###  parameters

• **parameters**: *[StringifiedContractEventParameters](#stringifiedcontracteventparameters)*

*Defined in [types.ts:198](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L198)*

___

###  txHash

• **txHash**: *string*

*Defined in [types.ts:192](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L192)*

___

###  txIndex

• **txIndex**: *number*

*Defined in [types.ts:193](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L193)*

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

*Defined in [types.ts:72](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L72)*

___

###  spender

• **spender**: *string*

*Defined in [types.ts:73](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L73)*

___

###  value

• **value**: *string*

*Defined in [types.ts:74](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L74)*

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

*Defined in [types.ts:60](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L60)*

___

###  to

• **to**: *string*

*Defined in [types.ts:61](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L61)*

___

###  value

• **value**: *string*

*Defined in [types.ts:62](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L62)*

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

*Defined in [types.ts:97](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L97)*

___

###  owner

• **owner**: *string*

*Defined in [types.ts:96](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L96)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:98](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L98)*

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

*Defined in [types.ts:84](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L84)*

___

###  to

• **to**: *string*

*Defined in [types.ts:85](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L85)*

___

###  tokenId

• **tokenId**: *string*

*Defined in [types.ts:86](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L86)*

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

*Defined in [types.ts:151](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L151)*

___

###  orderEpoch

• **orderEpoch**: *string*

*Defined in [types.ts:153](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L153)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:152](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L152)*

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

*Defined in [types.ts:125](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L125)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:122](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L122)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:131](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L131)*

___

###  makerAssetFilledAmount

• **makerAssetFilledAmount**: *string*

*Defined in [types.ts:126](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L126)*

___

###  makerFeePaid

• **makerFeePaid**: *string*

*Defined in [types.ts:128](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L128)*

___

###  orderHash

• **orderHash**: *string*

*Defined in [types.ts:130](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L130)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:124](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L124)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:123](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L123)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:132](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L132)*

___

###  takerAssetFilledAmount

• **takerAssetFilledAmount**: *string*

*Defined in [types.ts:127](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L127)*

___

###  takerFeePaid

• **takerFeePaid**: *string*

*Defined in [types.ts:129](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L129)*

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

*Defined in [types.ts:47](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L47)*

___

###  expirationTimeSeconds

• **expirationTimeSeconds**: *string*

*Defined in [types.ts:49](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L49)*

___

###  feeRecipientAddress

• **feeRecipientAddress**: *string*

*Defined in [types.ts:48](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L48)*

___

###  makerAddress

• **makerAddress**: *string*

*Defined in [types.ts:38](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L38)*

___

###  makerAssetAmount

• **makerAssetAmount**: *string*

*Defined in [types.ts:42](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L42)*

___

###  makerAssetData

• **makerAssetData**: *string*

*Defined in [types.ts:44](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L44)*

___

###  makerFee

• **makerFee**: *string*

*Defined in [types.ts:40](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L40)*

___

###  salt

• **salt**: *string*

*Defined in [types.ts:46](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L46)*

___

###  senderAddress

• **senderAddress**: *string*

*Defined in [types.ts:37](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L37)*

___

###  signature

• **signature**: *string*

*Defined in [types.ts:50](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L50)*

___

###  takerAddress

• **takerAddress**: *string*

*Defined in [types.ts:39](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L39)*

___

###  takerAssetAmount

• **takerAssetAmount**: *string*

*Defined in [types.ts:43](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L43)*

___

###  takerAssetData

• **takerAssetData**: *string*

*Defined in [types.ts:45](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L45)*

___

###  takerFee

• **takerFee**: *string*

*Defined in [types.ts:41](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L41)*

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

*Defined in [types.ts:172](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L172)*

___

###  value

• **value**: *string*

*Defined in [types.ts:173](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L173)*

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

*Defined in [types.ts:162](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L162)*

___

###  value

• **value**: *string*

*Defined in [types.ts:163](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L163)*

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

*Defined in [types.ts:325](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L325)*

___

###  rejected

• **rejected**: *[RejectedOrderInfo](#interface-rejectedorderinfo)[]*

*Defined in [types.ts:326](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L326)*

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

*Defined in [types.ts:167](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L167)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:168](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L168)*

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

*Defined in [types.ts:157](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L157)*

___

###  value

• **value**: *`BigNumber`*

*Defined in [types.ts:158](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L158)*

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

*Defined in [types.ts:335](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L335)*

___

###  utf8Data

• **utf8Data**: *string*

*Defined in [types.ts:336](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L336)*

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

*Defined in [types.ts:32](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L32)*

___

### `Optional` headers

• **headers**? : *undefined | `__type`*

*Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L30)*

___

### `Optional` protocol

• **protocol**? : *undefined | string*

*Defined in [types.ts:31](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L31)*

___

### `Optional` reconnectAfter

• **reconnectAfter**? : *undefined | number*

*Defined in [types.ts:33](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L33)*

___

### `Optional` timeout

• **timeout**? : *undefined | number*

*Defined in [types.ts:29](https://github.com/0xProject/0x-mesh/blob/b470283/rpc/clients/typescript/src/types.ts#L29)*

<hr />

