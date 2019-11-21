# Doc reference

> ## Class: WSClient

This class includes all the functionality related to interacting with a Mesh JSON RPC websocket endpoint.

### Hierarchy

* **WSClient**

### Index

#### Constructors

* [constructor](reference.md#constructor)

#### Methods

* [addOrdersAsync](reference.md#addordersasync)
* [destroy](reference.md#destroy)
* [getOrdersAsync](reference.md#getordersasync)
* [getStatsAsync](reference.md#getstatsasync)
* [onClose](reference.md#onclose)
* [onReconnected](reference.md#onreconnected)
* [subscribeToOrdersAsync](reference.md#subscribetoordersasync)
* [unsubscribeAsync](reference.md#unsubscribeasync)

### Constructors

#### constructor

+ **new WSClient**\(`url`: string, `wsOpts?`: [WSOpts](reference.md#interface-wsopts)\): [_WSClient_](reference.md#class-wsclient)

_Defined in_ [_ws\_client.ts:222_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L222)

Instantiates a new WSClient instance

**Parameters:**

| Name | Type | Description |
| :--- | :--- | :--- |
| `url` | string | WS server endpoint |
| `wsOpts?` | [WSOpts](reference.md#interface-wsopts) | WebSocket options |

**Returns:** [_WSClient_](reference.md#class-wsclient)

An instance of WSClient

### Methods

#### addOrdersAsync

▸ **addOrdersAsync**\(`signedOrders`: `SignedOrder`\[\], `pinned`: boolean\): _`Promise<ValidationResults>`_

_Defined in_ [_ws\_client.ts:251_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L251)

Adds an array of 0x signed orders to the Mesh node.

**Parameters:**

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `signedOrders` | `SignedOrder`\[\] | - | signedOrders to add |
| `pinned` | boolean | true | Whether or not the orders should be pinned. Pinned orders will not be affected by any DDoS prevention or incentive mechanisms and will always stay in storage until they are no longer fillable. |

**Returns:** _`Promise<ValidationResults>`_

validation results

#### destroy

▸ **destroy**\(\): _void_

_Defined in_ [_ws\_client.ts:366_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L366)

destroy unsubscribes all active subscriptions, closes the websocket connection and stops the internal heartbeat connection liveness check.

**Returns:** _void_

#### getOrdersAsync

▸ **getOrdersAsync**\(`perPage`: number\): _`Promise<OrderInfo[]>`_

_Defined in_ [_ws\_client.ts:281_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L281)

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `perPage` | number | 200 | number of signedOrders to fetch per paginated request |

**Returns:** _`Promise<OrderInfo[]>`_

all orders, their hash and their fillableTakerAssetAmount

#### getStatsAsync

▸ **getStatsAsync**\(\): _`Promise<GetStatsResponse>`_

_Defined in_ [_ws\_client.ts:272_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L272)

**Returns:** _`Promise<GetStatsResponse>`_

#### onClose

▸ **onClose**\(`cb`: function\): _void_

_Defined in_ [_ws\_client.ts:348_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L348)

Get notified when the underlying WS connection closes normally. If it closes with an error, WSClient automatically attempts to re-connect without emitting a `close` event.

**Parameters:**

▪ **cb**: _function_

callback to call when WS connection closes

▸ \(\): _void_

**Returns:** _void_

#### onReconnected

▸ **onReconnected**\(`cb`: function\): _void_

_Defined in_ [_ws\_client.ts:357_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L357)

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: _function_

callback to call with the error when it occurs

▸ \(\): _void_

**Returns:** _void_

#### subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**\(`cb`: function\): _`Promise<string>`_

_Defined in_ [_ws\_client.ts:309_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L309)

Subscribe to the 'orders' topic and receive order events from Mesh. This method returns a subscriptionId that can be used to `unsubscribe()` from this subscription.

**Parameters:**

▪ **cb**: _function_

callback function where you'd like to get notified about order events

▸ \(`orderEvents`: [OrderEvent](reference.md#interface-orderevent)\[\]\): _void_

**Parameters:**

| Name | Type |
| :--- | :--- |
| `orderEvents` | [OrderEvent](reference.md#interface-orderevent)\[\] |

**Returns:** _`Promise<string>`_

subscriptionId

#### unsubscribeAsync

▸ **unsubscribeAsync**\(`subscriptionId`: string\): _`Promise<void>`_

_Defined in_ [_ws\_client.ts:338_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/ws_client.ts#L338)

Unsubscribe from a subscription

**Parameters:**

| Name | Type | Description |
| :--- | :--- | :--- |
| `subscriptionId` | string | identifier of the subscription to cancel |

**Returns:** _`Promise<void>`_

> ## Enumeration: ContractEventKind

### Index

#### Enumeration members

* [ERC1155ApprovalForAllEvent](reference.md#erc1155approvalforallevent)
* [ERC1155TransferBatchEvent](reference.md#erc1155transferbatchevent)
* [ERC1155TransferSingleEvent](reference.md#erc1155transfersingleevent)
* [ERC20ApprovalEvent](reference.md#erc20approvalevent)
* [ERC20TransferEvent](reference.md#erc20transferevent)
* [ERC721ApprovalEvent](reference.md#erc721approvalevent)
* [ERC721ApprovalForAllEvent](reference.md#erc721approvalforallevent)
* [ERC721TransferEvent](reference.md#erc721transferevent)
* [ExchangeCancelEvent](reference.md#exchangecancelevent)
* [ExchangeCancelUpToEvent](reference.md#exchangecanceluptoevent)
* [ExchangeFillEvent](reference.md#exchangefillevent)
* [WethDepositEvent](reference.md#wethdepositevent)
* [WethWithdrawalEvent](reference.md#wethwithdrawalevent)

### Enumeration members

#### ERC1155ApprovalForAllEvent

• **ERC1155ApprovalForAllEvent**: = "ERC1155ApprovalForAllEvent"

_Defined in_ [_types.ts:220_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L220)

#### ERC1155TransferBatchEvent

• **ERC1155TransferBatchEvent**: = "ERC1155TransferBatchEvent"

_Defined in_ [_types.ts:222_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L222)

#### ERC1155TransferSingleEvent

• **ERC1155TransferSingleEvent**: = "ERC1155TransferSingleEvent"

_Defined in_ [_types.ts:221_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L221)

#### ERC20ApprovalEvent

• **ERC20ApprovalEvent**: = "ERC20ApprovalEvent"

_Defined in_ [_types.ts:216_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L216)

#### ERC20TransferEvent

• **ERC20TransferEvent**: = "ERC20TransferEvent"

_Defined in_ [_types.ts:215_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L215)

#### ERC721ApprovalEvent

• **ERC721ApprovalEvent**: = "ERC721ApprovalEvent"

_Defined in_ [_types.ts:218_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L218)

#### ERC721ApprovalForAllEvent

• **ERC721ApprovalForAllEvent**: = "ERC721ApprovalForAllEvent"

_Defined in_ [_types.ts:219_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L219)

#### ERC721TransferEvent

• **ERC721TransferEvent**: = "ERC721TransferEvent"

_Defined in_ [_types.ts:217_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L217)

#### ExchangeCancelEvent

• **ExchangeCancelEvent**: = "ExchangeCancelEvent"

_Defined in_ [_types.ts:224_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L224)

#### ExchangeCancelUpToEvent

• **ExchangeCancelUpToEvent**: = "ExchangeCancelUpToEvent"

_Defined in_ [_types.ts:225_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L225)

#### ExchangeFillEvent

• **ExchangeFillEvent**: = "ExchangeFillEvent"

_Defined in_ [_types.ts:223_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L223)

#### WethDepositEvent

• **WethDepositEvent**: = "WethDepositEvent"

_Defined in_ [_types.ts:226_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L226)

#### WethWithdrawalEvent

• **WethWithdrawalEvent**: = "WethWithdrawalEvent"

_Defined in_ [_types.ts:227_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L227)

> ## Enumeration: OrderEventEndState

### Index

#### Enumeration members

* [Added](reference.md#added)
* [Cancelled](reference.md#cancelled)
* [Expired](reference.md#expired)
* [FillabilityIncreased](reference.md#fillabilityincreased)
* [Filled](reference.md#filled)
* [FullyFilled](reference.md#fullyfilled)
* [Invalid](reference.md#invalid)
* [Unfunded](reference.md#unfunded)

### Enumeration members

#### Added

• **Added**: = "ADDED"

_Defined in_ [_types.ts:284_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L284)

#### Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in_ [_types.ts:287_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L287)

#### Expired

• **Expired**: = "EXPIRED"

_Defined in_ [_types.ts:288_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L288)

#### FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY\_INCREASED"

_Defined in_ [_types.ts:290_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L290)

#### Filled

• **Filled**: = "FILLED"

_Defined in_ [_types.ts:285_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L285)

#### FullyFilled

• **FullyFilled**: = "FULLY\_FILLED"

_Defined in_ [_types.ts:286_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L286)

#### Invalid

• **Invalid**: = "INVALID"

_Defined in_ [_types.ts:283_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L283)

#### Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in_ [_types.ts:289_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L289)

> ## Enumeration: RejectedCode

### Index

#### Enumeration members

* [InternalError](reference.md#internalerror)
* [MaxOrderSizeExceeded](reference.md#maxordersizeexceeded)
* [NetworkRequestFailed](reference.md#networkrequestfailed)
* [OrderAlreadyStored](reference.md#orderalreadystored)
* [OrderCancelled](reference.md#ordercancelled)
* [OrderExpired](reference.md#orderexpired)
* [OrderForIncorrectChain](reference.md#orderforincorrectchain)
* [OrderFullyFilled](reference.md#orderfullyfilled)
* [OrderHasInvalidMakerAssetAmount](reference.md#orderhasinvalidmakerassetamount)
* [OrderHasInvalidMakerAssetData](reference.md#orderhasinvalidmakerassetdata)
* [OrderHasInvalidSignature](reference.md#orderhasinvalidsignature)
* [OrderHasInvalidTakerAssetAmount](reference.md#orderhasinvalidtakerassetamount)
* [OrderHasInvalidTakerAssetData](reference.md#orderhasinvalidtakerassetdata)
* [OrderUnfunded](reference.md#orderunfunded)

### Enumeration members

#### InternalError

• **InternalError**: = "InternalError"

_Defined in_ [_types.ts:352_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L352)

#### MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MaxOrderSizeExceeded"

_Defined in_ [_types.ts:353_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L353)

#### NetworkRequestFailed

• **NetworkRequestFailed**: = "NetworkRequestFailed"

_Defined in_ [_types.ts:356_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L356)

#### OrderAlreadyStored

• **OrderAlreadyStored**: = "OrderAlreadyStored"

_Defined in_ [_types.ts:354_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L354)

#### OrderCancelled

• **OrderCancelled**: = "OrderCancelled"

_Defined in_ [_types.ts:361_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L361)

#### OrderExpired

• **OrderExpired**: = "OrderExpired"

_Defined in_ [_types.ts:359_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L359)

#### OrderForIncorrectChain

• **OrderForIncorrectChain**: = "OrderForIncorrectChain"

_Defined in_ [_types.ts:355_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L355)

#### OrderFullyFilled

• **OrderFullyFilled**: = "OrderFullyFilled"

_Defined in_ [_types.ts:360_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L360)

#### OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "OrderHasInvalidMakerAssetAmount"

_Defined in_ [_types.ts:357_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L357)

#### OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "OrderHasInvalidMakerAssetData"

_Defined in_ [_types.ts:363_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L363)

#### OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "OrderHasInvalidSignature"

_Defined in_ [_types.ts:365_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L365)

#### OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "OrderHasInvalidTakerAssetAmount"

_Defined in_ [_types.ts:358_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L358)

#### OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "OrderHasInvalidTakerAssetData"

_Defined in_ [_types.ts:364_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L364)

#### OrderUnfunded

• **OrderUnfunded**: = "OrderUnfunded"

_Defined in_ [_types.ts:362_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L362)

> ## Enumeration: RejectedKind

### Index

#### Enumeration members

* [MeshError](reference.md#mesherror)
* [MeshValidation](reference.md#meshvalidation)
* [ZeroexValidation](reference.md#zeroexvalidation)

### Enumeration members

#### MeshError

• **MeshError**: = "MESH\_ERROR"

_Defined in_ [_types.ts:347_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L347)

#### MeshValidation

• **MeshValidation**: = "MESH\_VALIDATION"

_Defined in_ [_types.ts:348_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L348)

#### ZeroexValidation

• **ZeroexValidation**: = "ZEROEX\_VALIDATION"

_Defined in_ [_types.ts:346_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L346)

> ## Interface: AcceptedOrderInfo

### Hierarchy

* **AcceptedOrderInfo**

### Index

#### Properties

* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [isNew](reference.md#isnew)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)

### Properties

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _`BigNumber`_

_Defined in_ [_types.ts:329_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L329)

#### isNew

• **isNew**: _boolean_

_Defined in_ [_types.ts:330_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L330)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:327_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L327)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:328_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L328)

> ## Interface: ClientConfig

WebSocketClient configs Source: [https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md\#client-config-options](https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md#client-config-options)

### Hierarchy

* **ClientConfig**

### Index

#### Properties

* [assembleFragments](reference.md#optional-assemblefragments)
* [closeTimeout](reference.md#optional-closetimeout)
* [fragmentOutgoingMessages](reference.md#optional-fragmentoutgoingmessages)
* [fragmentationThreshold](reference.md#optional-fragmentationthreshold)
* [maxReceivedFrameSize](reference.md#optional-maxreceivedframesize)
* [maxReceivedMessageSize](reference.md#optional-maxreceivedmessagesize)
* [tlsOptions](reference.md#optional-tlsoptions)
* [webSocketVersion](reference.md#optional-websocketversion)

### Properties

#### `Optional` assembleFragments

• **assembleFragments**? : _undefined \| false \| true_

_Defined in_ [_types.ts:14_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L14)

#### `Optional` closeTimeout

• **closeTimeout**? : _undefined \| number_

_Defined in_ [_types.ts:15_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L15)

#### `Optional` fragmentOutgoingMessages

• **fragmentOutgoingMessages**? : _undefined \| false \| true_

_Defined in_ [_types.ts:12_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L12)

#### `Optional` fragmentationThreshold

• **fragmentationThreshold**? : _undefined \| number_

_Defined in_ [_types.ts:13_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L13)

#### `Optional` maxReceivedFrameSize

• **maxReceivedFrameSize**? : _undefined \| number_

_Defined in_ [_types.ts:10_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L10)

#### `Optional` maxReceivedMessageSize

• **maxReceivedMessageSize**? : _undefined \| number_

_Defined in_ [_types.ts:11_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L11)

#### `Optional` tlsOptions

• **tlsOptions**? : _any_

_Defined in_ [_types.ts:16_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L16)

#### `Optional` webSocketVersion

• **webSocketVersion**? : _undefined \| number_

_Defined in_ [_types.ts:9_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L9)

> ## Interface: ContractEvent

### Hierarchy

* **ContractEvent**

### Index

#### Properties

* [address](reference.md#address)
* [blockHash](reference.md#blockhash)
* [isRemoved](reference.md#isremoved)
* [kind](reference.md#kind)
* [logIndex](reference.md#logindex)
* [parameters](reference.md#parameters)
* [txHash](reference.md#txhash)
* [txIndex](reference.md#txindex)

### Properties

#### address

• **address**: _string_

_Defined in_ [_types.ts:277_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L277)

#### blockHash

• **blockHash**: _string_

_Defined in_ [_types.ts:272_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L272)

#### isRemoved

• **isRemoved**: _string_

_Defined in_ [_types.ts:276_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L276)

#### kind

• **kind**: [_ContractEventKind_](reference.md#enumeration-contracteventkind)

_Defined in_ [_types.ts:278_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L278)

#### logIndex

• **logIndex**: _number_

_Defined in_ [_types.ts:275_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L275)

#### parameters

• **parameters**: [_ContractEventParameters_](reference.md#contracteventparameters)

_Defined in_ [_types.ts:279_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L279)

#### txHash

• **txHash**: _string_

_Defined in_ [_types.ts:273_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L273)

#### txIndex

• **txIndex**: _number_

_Defined in_ [_types.ts:274_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L274)

> ## Interface: ERC1155ApprovalForAllEvent

### Hierarchy

* **ERC1155ApprovalForAllEvent**

### Index

#### Properties

* [approved](reference.md#approved)
* [operator](reference.md#operator)
* [owner](reference.md#owner)

### Properties

#### approved

• **approved**: _boolean_

_Defined in_ [_types.ts:142_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L142)

#### operator

• **operator**: _string_

_Defined in_ [_types.ts:141_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L141)

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:140_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L140)

> ## Interface: ERC1155TransferBatchEvent

### Hierarchy

* **ERC1155TransferBatchEvent**

### Index

#### Properties

* [from](reference.md#from)
* [ids](reference.md#ids)
* [operator](reference.md#operator)
* [to](reference.md#to)
* [values](reference.md#values)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:125_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L125)

#### ids

• **ids**: _`BigNumber`\[\]_

_Defined in_ [_types.ts:127_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L127)

#### operator

• **operator**: _string_

_Defined in_ [_types.ts:124_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L124)

#### to

• **to**: _string_

_Defined in_ [_types.ts:126_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L126)

#### values

• **values**: _`BigNumber`\[\]_

_Defined in_ [_types.ts:128_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L128)

> ## Interface: ERC1155TransferSingleEvent

### Hierarchy

* **ERC1155TransferSingleEvent**

### Index

#### Properties

* [from](reference.md#from)
* [id](reference.md#id)
* [operator](reference.md#operator)
* [to](reference.md#to)
* [value](reference.md#value)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:109_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L109)

#### id

• **id**: _`BigNumber`_

_Defined in_ [_types.ts:111_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L111)

#### operator

• **operator**: _string_

_Defined in_ [_types.ts:108_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L108)

#### to

• **to**: _string_

_Defined in_ [_types.ts:110_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L110)

#### value

• **value**: _`BigNumber`_

_Defined in_ [_types.ts:112_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L112)

> ## Interface: ERC20ApprovalEvent

### Hierarchy

* **ERC20ApprovalEvent**

### Index

#### Properties

* [owner](reference.md#owner)
* [spender](reference.md#spender)
* [value](reference.md#value)

### Properties

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:66_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L66)

#### spender

• **spender**: _string_

_Defined in_ [_types.ts:67_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L67)

#### value

• **value**: _`BigNumber`_

_Defined in_ [_types.ts:68_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L68)

> ## Interface: ERC20TransferEvent

### Hierarchy

* **ERC20TransferEvent**

### Index

#### Properties

* [from](reference.md#from)
* [to](reference.md#to)
* [value](reference.md#value)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:54_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L54)

#### to

• **to**: _string_

_Defined in_ [_types.ts:55_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L55)

#### value

• **value**: _`BigNumber`_

_Defined in_ [_types.ts:56_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L56)

> ## Interface: ERC721ApprovalEvent

### Hierarchy

* **ERC721ApprovalEvent**

### Index

#### Properties

* [approved](reference.md#approved)
* [owner](reference.md#owner)
* [tokenId](reference.md#tokenid)

### Properties

#### approved

• **approved**: _string_

_Defined in_ [_types.ts:91_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L91)

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:90_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L90)

#### tokenId

• **tokenId**: _`BigNumber`_

_Defined in_ [_types.ts:92_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L92)

> ## Interface: ERC721ApprovalForAllEvent

### Hierarchy

* **ERC721ApprovalForAllEvent**

### Index

#### Properties

* [approved](reference.md#approved)
* [operator](reference.md#operator)
* [owner](reference.md#owner)

### Properties

#### approved

• **approved**: _boolean_

_Defined in_ [_types.ts:104_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L104)

#### operator

• **operator**: _string_

_Defined in_ [_types.ts:103_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L103)

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:102_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L102)

> ## Interface: ERC721TransferEvent

### Hierarchy

* **ERC721TransferEvent**

### Index

#### Properties

* [from](reference.md#from)
* [to](reference.md#to)
* [tokenId](reference.md#tokenid)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:78_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L78)

#### to

• **to**: _string_

_Defined in_ [_types.ts:79_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L79)

#### tokenId

• **tokenId**: _`BigNumber`_

_Defined in_ [_types.ts:80_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L80)

> ## Interface: ExchangeCancelEvent

### Hierarchy

* **ExchangeCancelEvent**

### Index

#### Properties

* [feeRecipientAddress](reference.md#feerecipientaddress)
* [makerAddress](reference.md#makeraddress)
* [makerAssetData](reference.md#makerassetdata)
* [orderHash](reference.md#orderhash)
* [senderAddress](reference.md#senderaddress)
* [takerAssetData](reference.md#takerassetdata)

### Properties

#### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_types.ts:176_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L176)

#### makerAddress

• **makerAddress**: _string_

_Defined in_ [_types.ts:174_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L174)

#### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_types.ts:178_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L178)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:177_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L177)

#### senderAddress

• **senderAddress**: _string_

_Defined in_ [_types.ts:175_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L175)

#### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_types.ts:179_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L179)

> ## Interface: ExchangeCancelUpToEvent

### Hierarchy

* **ExchangeCancelUpToEvent**

### Index

#### Properties

* [makerAddress](reference.md#makeraddress)
* [orderEpoch](reference.md#orderepoch)
* [senderAddress](reference.md#senderaddress)

### Properties

#### makerAddress

• **makerAddress**: _string_

_Defined in_ [_types.ts:183_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L183)

#### orderEpoch

• **orderEpoch**: _`BigNumber`_

_Defined in_ [_types.ts:185_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L185)

#### senderAddress

• **senderAddress**: _string_

_Defined in_ [_types.ts:184_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L184)

> ## Interface: ExchangeFillEvent

### Hierarchy

* **ExchangeFillEvent**

### Index

#### Properties

* [feeRecipientAddress](reference.md#feerecipientaddress)
* [makerAddress](reference.md#makeraddress)
* [makerAssetData](reference.md#makerassetdata)
* [makerAssetFilledAmount](reference.md#makerassetfilledamount)
* [makerFeePaid](reference.md#makerfeepaid)
* [orderHash](reference.md#orderhash)
* [senderAddress](reference.md#senderaddress)
* [takerAddress](reference.md#takeraddress)
* [takerAssetData](reference.md#takerassetdata)
* [takerAssetFilledAmount](reference.md#takerassetfilledamount)
* [takerFeePaid](reference.md#takerfeepaid)

### Properties

#### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_types.ts:149_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L149)

#### makerAddress

• **makerAddress**: _string_

_Defined in_ [_types.ts:146_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L146)

#### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_types.ts:155_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L155)

#### makerAssetFilledAmount

• **makerAssetFilledAmount**: _`BigNumber`_

_Defined in_ [_types.ts:150_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L150)

#### makerFeePaid

• **makerFeePaid**: _`BigNumber`_

_Defined in_ [_types.ts:152_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L152)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:154_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L154)

#### senderAddress

• **senderAddress**: _string_

_Defined in_ [_types.ts:148_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L148)

#### takerAddress

• **takerAddress**: _string_

_Defined in_ [_types.ts:147_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L147)

#### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_types.ts:156_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L156)

#### takerAssetFilledAmount

• **takerAssetFilledAmount**: _`BigNumber`_

_Defined in_ [_types.ts:151_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L151)

#### takerFeePaid

• **takerFeePaid**: _`BigNumber`_

_Defined in_ [_types.ts:153_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L153)

> ## Interface: GetOrdersResponse

### Hierarchy

* **GetOrdersResponse**

### Index

#### Properties

* [ordersInfos](reference.md#ordersinfos)
* [snapshotID](reference.md#snapshotid)

### Properties

#### ordersInfos

• **ordersInfos**: [_RawAcceptedOrderInfo_](reference.md#interface-rawacceptedorderinfo)_\[\]_

_Defined in_ [_types.ts:399_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L399)

#### snapshotID

• **snapshotID**: _string_

_Defined in_ [_types.ts:398_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L398)

> ## Interface: GetStatsResponse

### Hierarchy

* **GetStatsResponse**

### Index

#### Properties

* [ethRPCRateLimitExpiredRequests](reference.md#ethrpcratelimitexpiredrequests)
* [ethRPCRequestsSentInCurrentUTCDay](reference.md#ethrpcrequestssentincurrentutcday)
* [ethereumChainID](reference.md#ethereumchainid)
* [latestBlock](reference.md#latestblock)
* [maxExpirationTime](reference.md#maxexpirationtime)
* [numOrders](reference.md#numorders)
* [numOrdersIncludingRemoved](reference.md#numordersincludingremoved)
* [numPeers](reference.md#numpeers)
* [numPinnedOrders](reference.md#numpinnedorders)
* [peerID](reference.md#peerid)
* [pubSubTopic](reference.md#pubsubtopic)
* [rendezvous](reference.md#rendezvous)
* [startOfCurrentUTCDay](reference.md#startofcurrentutcday)
* [version](reference.md#version)

### Properties

#### ethRPCRateLimitExpiredRequests

• **ethRPCRateLimitExpiredRequests**: _number_

_Defined in_ [_types.ts:426_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L426)

#### ethRPCRequestsSentInCurrentUTCDay

• **ethRPCRequestsSentInCurrentUTCDay**: _number_

_Defined in_ [_types.ts:425_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L425)

#### ethereumChainID

• **ethereumChainID**: _number_

_Defined in_ [_types.ts:417_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L417)

#### latestBlock

• **latestBlock**: [_LatestBlock_](reference.md#interface-latestblock)

_Defined in_ [_types.ts:418_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L418)

#### maxExpirationTime

• **maxExpirationTime**: _string_

_Defined in_ [_types.ts:423_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L423)

#### numOrders

• **numOrders**: _number_

_Defined in_ [_types.ts:420_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L420)

#### numOrdersIncludingRemoved

• **numOrdersIncludingRemoved**: _number_

_Defined in_ [_types.ts:421_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L421)

#### numPeers

• **numPeers**: _number_

_Defined in_ [_types.ts:419_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L419)

#### numPinnedOrders

• **numPinnedOrders**: _number_

_Defined in_ [_types.ts:422_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L422)

#### peerID

• **peerID**: _string_

_Defined in_ [_types.ts:416_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L416)

#### pubSubTopic

• **pubSubTopic**: _string_

_Defined in_ [_types.ts:414_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L414)

#### rendezvous

• **rendezvous**: _string_

_Defined in_ [_types.ts:415_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L415)

#### startOfCurrentUTCDay

• **startOfCurrentUTCDay**: _string_

_Defined in_ [_types.ts:424_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L424)

#### version

• **version**: _string_

_Defined in_ [_types.ts:413_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L413)

> ## Interface: HeartbeatEventPayload

### Hierarchy

* **HeartbeatEventPayload**

### Index

#### Properties

* [result](reference.md#result)
* [subscription](reference.md#subscription)

### Properties

#### result

• **result**: _string_

_Defined in_ [_types.ts:300_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L300)

#### subscription

• **subscription**: _string_

_Defined in_ [_types.ts:299_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L299)

> ## Interface: LatestBlock

### Hierarchy

* **LatestBlock**

### Index

#### Properties

* [hash](reference.md#hash)
* [number](reference.md#number)

### Properties

#### hash

• **hash**: _string_

_Defined in_ [_types.ts:409_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L409)

#### number

• **number**: _number_

_Defined in_ [_types.ts:408_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L408)

> ## Interface: OrderEvent

### Hierarchy

* **OrderEvent**

### Index

#### Properties

* [contractEvents](reference.md#contractevents)
* [endState](reference.md#endstate)
* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)

### Properties

#### contractEvents

• **contractEvents**: [_ContractEvent_](reference.md#interface-contractevent)_\[\]_

_Defined in_ [_types.ts:316_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L316)

#### endState

• **endState**: [_OrderEventEndState_](reference.md#enumeration-ordereventendstate)

_Defined in_ [_types.ts:314_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L314)

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _`BigNumber`_

_Defined in_ [_types.ts:315_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L315)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:312_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L312)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:313_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L313)

> ## Interface: OrderEventPayload

### Hierarchy

* **OrderEventPayload**

### Index

#### Properties

* [result](reference.md#result)
* [subscription](reference.md#subscription)

### Properties

#### result

• **result**: [_RawOrderEvent_](reference.md#interface-raworderevent)_\[\]_

_Defined in_ [_types.ts:295_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L295)

#### subscription

• **subscription**: _string_

_Defined in_ [_types.ts:294_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L294)

> ## Interface: OrderInfo

### Hierarchy

* **OrderInfo**

### Index

#### Properties

* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)

### Properties

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _`BigNumber`_

_Defined in_ [_types.ts:342_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L342)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:340_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L340)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:341_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L341)

> ## Interface: RawAcceptedOrderInfo

### Hierarchy

* **RawAcceptedOrderInfo**

### Index

#### Properties

* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [isNew](reference.md#isnew)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)

### Properties

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in_ [_types.ts:322_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L322)

#### isNew

• **isNew**: _boolean_

_Defined in_ [_types.ts:323_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L323)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:320_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L320)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:321_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L321)

> ## Interface: RawOrderEvent

### Hierarchy

* **RawOrderEvent**

### Index

#### Properties

* [contractEvents](reference.md#contractevents)
* [endState](reference.md#endstate)
* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)

### Properties

#### contractEvents

• **contractEvents**: [_StringifiedContractEvent_](reference.md#interface-stringifiedcontractevent)_\[\]_

_Defined in_ [_types.ts:308_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L308)

#### endState

• **endState**: [_OrderEventEndState_](reference.md#enumeration-ordereventendstate)

_Defined in_ [_types.ts:306_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L306)

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in_ [_types.ts:307_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L307)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:304_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L304)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:305_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L305)

> ## Interface: RawOrderInfo

### Hierarchy

* **RawOrderInfo**

### Index

#### Properties

* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)

### Properties

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in_ [_types.ts:336_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L336)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:334_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L334)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:335_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L335)

> ## Interface: RawRejectedOrderInfo

### Hierarchy

* **RawRejectedOrderInfo**

### Index

#### Properties

* [kind](reference.md#kind)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)
* [status](reference.md#status)

### Properties

#### kind

• **kind**: [_RejectedKind_](reference.md#enumeration-rejectedkind)

_Defined in_ [_types.ts:376_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L376)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:374_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L374)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:375_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L375)

#### status

• **status**: [_RejectedStatus_](reference.md#interface-rejectedstatus)

_Defined in_ [_types.ts:377_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L377)

> ## Interface: RawValidationResults

### Hierarchy

* **RawValidationResults**

### Index

#### Properties

* [accepted](reference.md#accepted)
* [rejected](reference.md#rejected)

### Properties

#### accepted

• **accepted**: [_RawAcceptedOrderInfo_](reference.md#interface-rawacceptedorderinfo)_\[\]_

_Defined in_ [_types.ts:388_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L388)

#### rejected

• **rejected**: [_RawRejectedOrderInfo_](reference.md#interface-rawrejectedorderinfo)_\[\]_

_Defined in_ [_types.ts:389_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L389)

> ## Interface: RejectedOrderInfo

### Hierarchy

* **RejectedOrderInfo**

### Index

#### Properties

* [kind](reference.md#kind)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)
* [status](reference.md#status)

### Properties

#### kind

• **kind**: [_RejectedKind_](reference.md#enumeration-rejectedkind)

_Defined in_ [_types.ts:383_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L383)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:381_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L381)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:382_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L382)

#### status

• **status**: [_RejectedStatus_](reference.md#interface-rejectedstatus)

_Defined in_ [_types.ts:384_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L384)

> ## Interface: RejectedStatus

### Hierarchy

* **RejectedStatus**

### Index

#### Properties

* [code](reference.md#code)
* [message](reference.md#message)

### Properties

#### code

• **code**: [_RejectedCode_](reference.md#enumeration-rejectedcode)

_Defined in_ [_types.ts:369_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L369)

#### message

• **message**: _string_

_Defined in_ [_types.ts:370_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L370)

> ## Interface: StringifiedContractEvent

### Hierarchy

* **StringifiedContractEvent**

### Index

#### Properties

* [address](reference.md#address)
* [blockHash](reference.md#blockhash)
* [isRemoved](reference.md#isremoved)
* [kind](reference.md#kind)
* [logIndex](reference.md#logindex)
* [parameters](reference.md#parameters)
* [txHash](reference.md#txhash)
* [txIndex](reference.md#txindex)

### Properties

#### address

• **address**: _string_

_Defined in_ [_types.ts:251_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L251)

#### blockHash

• **blockHash**: _string_

_Defined in_ [_types.ts:246_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L246)

#### isRemoved

• **isRemoved**: _string_

_Defined in_ [_types.ts:250_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L250)

#### kind

• **kind**: _string_

_Defined in_ [_types.ts:252_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L252)

#### logIndex

• **logIndex**: _number_

_Defined in_ [_types.ts:249_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L249)

#### parameters

• **parameters**: [_StringifiedContractEventParameters_](reference.md#stringifiedcontracteventparameters)

_Defined in_ [_types.ts:253_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L253)

#### txHash

• **txHash**: _string_

_Defined in_ [_types.ts:247_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L247)

#### txIndex

• **txIndex**: _number_

_Defined in_ [_types.ts:248_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L248)

> ## Interface: StringifiedERC1155TransferBatchEvent

### Hierarchy

* **StringifiedERC1155TransferBatchEvent**

### Index

#### Properties

* [from](reference.md#from)
* [ids](reference.md#ids)
* [operator](reference.md#operator)
* [to](reference.md#to)
* [values](reference.md#values)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:133_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L133)

#### ids

• **ids**: _string\[\]_

_Defined in_ [_types.ts:135_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L135)

#### operator

• **operator**: _string_

_Defined in_ [_types.ts:132_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L132)

#### to

• **to**: _string_

_Defined in_ [_types.ts:134_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L134)

#### values

• **values**: _string\[\]_

_Defined in_ [_types.ts:136_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L136)

> ## Interface: StringifiedERC1155TransferSingleEvent

### Hierarchy

* **StringifiedERC1155TransferSingleEvent**

### Index

#### Properties

* [from](reference.md#from)
* [id](reference.md#id)
* [operator](reference.md#operator)
* [to](reference.md#to)
* [value](reference.md#value)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:117_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L117)

#### id

• **id**: _string_

_Defined in_ [_types.ts:119_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L119)

#### operator

• **operator**: _string_

_Defined in_ [_types.ts:116_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L116)

#### to

• **to**: _string_

_Defined in_ [_types.ts:118_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L118)

#### value

• **value**: _string_

_Defined in_ [_types.ts:120_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L120)

> ## Interface: StringifiedERC20ApprovalEvent

### Hierarchy

* **StringifiedERC20ApprovalEvent**

### Index

#### Properties

* [owner](reference.md#owner)
* [spender](reference.md#spender)
* [value](reference.md#value)

### Properties

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:72_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L72)

#### spender

• **spender**: _string_

_Defined in_ [_types.ts:73_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L73)

#### value

• **value**: _string_

_Defined in_ [_types.ts:74_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L74)

> ## Interface: StringifiedERC20TransferEvent

### Hierarchy

* **StringifiedERC20TransferEvent**

### Index

#### Properties

* [from](reference.md#from)
* [to](reference.md#to)
* [value](reference.md#value)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:60_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L60)

#### to

• **to**: _string_

_Defined in_ [_types.ts:61_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L61)

#### value

• **value**: _string_

_Defined in_ [_types.ts:62_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L62)

> ## Interface: StringifiedERC721ApprovalEvent

### Hierarchy

* **StringifiedERC721ApprovalEvent**

### Index

#### Properties

* [approved](reference.md#approved)
* [owner](reference.md#owner)
* [tokenId](reference.md#tokenid)

### Properties

#### approved

• **approved**: _string_

_Defined in_ [_types.ts:97_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L97)

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:96_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L96)

#### tokenId

• **tokenId**: _string_

_Defined in_ [_types.ts:98_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L98)

> ## Interface: StringifiedERC721TransferEvent

### Hierarchy

* **StringifiedERC721TransferEvent**

### Index

#### Properties

* [from](reference.md#from)
* [to](reference.md#to)
* [tokenId](reference.md#tokenid)

### Properties

#### from

• **from**: _string_

_Defined in_ [_types.ts:84_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L84)

#### to

• **to**: _string_

_Defined in_ [_types.ts:85_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L85)

#### tokenId

• **tokenId**: _string_

_Defined in_ [_types.ts:86_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L86)

> ## Interface: StringifiedExchangeCancelUpToEvent

### Hierarchy

* **StringifiedExchangeCancelUpToEvent**

### Index

#### Properties

* [makerAddress](reference.md#makeraddress)
* [orderEpoch](reference.md#orderepoch)
* [senderAddress](reference.md#senderaddress)

### Properties

#### makerAddress

• **makerAddress**: _string_

_Defined in_ [_types.ts:189_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L189)

#### orderEpoch

• **orderEpoch**: _string_

_Defined in_ [_types.ts:191_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L191)

#### senderAddress

• **senderAddress**: _string_

_Defined in_ [_types.ts:190_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L190)

> ## Interface: StringifiedExchangeFillEvent

### Hierarchy

* **StringifiedExchangeFillEvent**

### Index

#### Properties

* [feeRecipientAddress](reference.md#feerecipientaddress)
* [makerAddress](reference.md#makeraddress)
* [makerAssetData](reference.md#makerassetdata)
* [makerAssetFilledAmount](reference.md#makerassetfilledamount)
* [makerFeePaid](reference.md#makerfeepaid)
* [orderHash](reference.md#orderhash)
* [senderAddress](reference.md#senderaddress)
* [takerAddress](reference.md#takeraddress)
* [takerAssetData](reference.md#takerassetdata)
* [takerAssetFilledAmount](reference.md#takerassetfilledamount)
* [takerFeePaid](reference.md#takerfeepaid)

### Properties

#### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_types.ts:163_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L163)

#### makerAddress

• **makerAddress**: _string_

_Defined in_ [_types.ts:160_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L160)

#### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_types.ts:169_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L169)

#### makerAssetFilledAmount

• **makerAssetFilledAmount**: _string_

_Defined in_ [_types.ts:164_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L164)

#### makerFeePaid

• **makerFeePaid**: _string_

_Defined in_ [_types.ts:166_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L166)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:168_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L168)

#### senderAddress

• **senderAddress**: _string_

_Defined in_ [_types.ts:162_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L162)

#### takerAddress

• **takerAddress**: _string_

_Defined in_ [_types.ts:161_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L161)

#### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_types.ts:170_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L170)

#### takerAssetFilledAmount

• **takerAssetFilledAmount**: _string_

_Defined in_ [_types.ts:165_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L165)

#### takerFeePaid

• **takerFeePaid**: _string_

_Defined in_ [_types.ts:167_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L167)

> ## Interface: StringifiedSignedOrder

### Hierarchy

* **StringifiedSignedOrder**

### Index

#### Properties

* [exchangeAddress](reference.md#exchangeaddress)
* [expirationTimeSeconds](reference.md#expirationtimeseconds)
* [feeRecipientAddress](reference.md#feerecipientaddress)
* [makerAddress](reference.md#makeraddress)
* [makerAssetAmount](reference.md#makerassetamount)
* [makerAssetData](reference.md#makerassetdata)
* [makerFee](reference.md#makerfee)
* [salt](reference.md#salt)
* [senderAddress](reference.md#senderaddress)
* [signature](reference.md#signature)
* [takerAddress](reference.md#takeraddress)
* [takerAssetAmount](reference.md#takerassetamount)
* [takerAssetData](reference.md#takerassetdata)
* [takerFee](reference.md#takerfee)

### Properties

#### exchangeAddress

• **exchangeAddress**: _string_

_Defined in_ [_types.ts:47_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L47)

#### expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Defined in_ [_types.ts:49_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L49)

#### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_types.ts:48_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L48)

#### makerAddress

• **makerAddress**: _string_

_Defined in_ [_types.ts:38_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L38)

#### makerAssetAmount

• **makerAssetAmount**: _string_

_Defined in_ [_types.ts:42_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L42)

#### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_types.ts:44_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L44)

#### makerFee

• **makerFee**: _string_

_Defined in_ [_types.ts:40_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L40)

#### salt

• **salt**: _string_

_Defined in_ [_types.ts:46_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L46)

#### senderAddress

• **senderAddress**: _string_

_Defined in_ [_types.ts:37_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L37)

#### signature

• **signature**: _string_

_Defined in_ [_types.ts:50_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L50)

#### takerAddress

• **takerAddress**: _string_

_Defined in_ [_types.ts:39_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L39)

#### takerAssetAmount

• **takerAssetAmount**: _string_

_Defined in_ [_types.ts:43_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L43)

#### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_types.ts:45_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L45)

#### takerFee

• **takerFee**: _string_

_Defined in_ [_types.ts:41_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L41)

> ## Interface: StringifiedWethDepositEvent

### Hierarchy

* **StringifiedWethDepositEvent**

### Index

#### Properties

* [owner](reference.md#owner)
* [value](reference.md#value)

### Properties

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:210_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L210)

#### value

• **value**: _string_

_Defined in_ [_types.ts:211_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L211)

> ## Interface: StringifiedWethWithdrawalEvent

### Hierarchy

* **StringifiedWethWithdrawalEvent**

### Index

#### Properties

* [owner](reference.md#owner)
* [value](reference.md#value)

### Properties

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:200_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L200)

#### value

• **value**: _string_

_Defined in_ [_types.ts:201_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L201)

> ## Interface: ValidationResults

### Hierarchy

* **ValidationResults**

### Index

#### Properties

* [accepted](reference.md#accepted)
* [rejected](reference.md#rejected)

### Properties

#### accepted

• **accepted**: [_AcceptedOrderInfo_](reference.md#interface-acceptedorderinfo)_\[\]_

_Defined in_ [_types.ts:393_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L393)

#### rejected

• **rejected**: [_RejectedOrderInfo_](reference.md#interface-rejectedorderinfo)_\[\]_

_Defined in_ [_types.ts:394_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L394)

> ## Interface: WethDepositEvent

### Hierarchy

* **WethDepositEvent**

### Index

#### Properties

* [owner](reference.md#owner)
* [value](reference.md#value)

### Properties

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:205_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L205)

#### value

• **value**: _`BigNumber`_

_Defined in_ [_types.ts:206_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L206)

> ## Interface: WethWithdrawalEvent

### Hierarchy

* **WethWithdrawalEvent**

### Index

#### Properties

* [owner](reference.md#owner)
* [value](reference.md#value)

### Properties

#### owner

• **owner**: _string_

_Defined in_ [_types.ts:195_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L195)

#### value

• **value**: _`BigNumber`_

_Defined in_ [_types.ts:196_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L196)

> ## Interface: WSMessage

### Hierarchy

* **WSMessage**

### Index

#### Properties

* [type](reference.md#type)
* [utf8Data](reference.md#utf8data)

### Properties

#### type

• **type**: _string_

_Defined in_ [_types.ts:403_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L403)

#### utf8Data

• **utf8Data**: _string_

_Defined in_ [_types.ts:404_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L404)

> ## Interface: WSOpts

timeout: timeout in milliseconds to enforce on every WS request that expects a response headers: Request headers \(e.g., authorization\) protocol: requestOptions should be either null or an object specifying additional configuration options to be passed to http.request or https.request. This can be used to pass a custom agent to enable WebSocketClient usage from behind an HTTP or HTTPS proxy server using koichik/node-tunnel or similar. clientConfig: The client configs documented here: [https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md](https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md) reconnectAfter: time in milliseconds after which to attempt to reconnect to WS server after an error occurred \(default: 5000\)

### Hierarchy

* **WSOpts**

### Index

#### Properties

* [clientConfig](reference.md#optional-clientconfig)
* [headers](reference.md#optional-headers)
* [protocol](reference.md#optional-protocol)
* [reconnectAfter](reference.md#optional-reconnectafter)
* [timeout](reference.md#optional-timeout)

### Properties

#### `Optional` clientConfig

• **clientConfig**? : [_ClientConfig_](reference.md#interface-clientconfig)

_Defined in_ [_types.ts:32_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L32)

#### `Optional` headers

• **headers**? : _undefined \| `__type`_

_Defined in_ [_types.ts:30_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L30)

#### `Optional` protocol

• **protocol**? : _undefined \| string_

_Defined in_ [_types.ts:31_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L31)

#### `Optional` reconnectAfter

• **reconnectAfter**? : _undefined \| number_

_Defined in_ [_types.ts:33_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L33)

#### `Optional` timeout

• **timeout**? : _undefined \| number_

_Defined in_ [_types.ts:29_](https://github.com/0xProject/0x-mesh/blob/0240f5d/rpc/clients/typescript/src/types.ts#L29)

