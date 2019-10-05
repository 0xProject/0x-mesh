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

_Defined in_ [_ws\_client.ts:71_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L71)

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

▸ **addOrdersAsync**\(`signedOrders`: `SignedOrder`\[\]\): _`Promise<ValidationResults>`_

_Defined in_ [_ws\_client.ts:96_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L96)

Adds an array of 0x signed orders to the Mesh node.

**Parameters:**

| Name | Type | Description |
| :--- | :--- | :--- |
| `signedOrders` | `SignedOrder`\[\] | signedOrders to add |

**Returns:** _`Promise<ValidationResults>`_

validation results

#### destroy

▸ **destroy**\(\): _void_

_Defined in_ [_ws\_client.ts:213_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L213)

destroy unsubscribes all active subscriptions, closes the websocket connection and stops the internal heartbeat connection liveness check.

**Returns:** _void_

#### getOrdersAsync

▸ **getOrdersAsync**\(`perPage`: number\): _`Promise<OrderInfo[]>`_

_Defined in_ [_ws\_client.ts:125_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L125)

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `perPage` | number | 200 | number of signedOrders to fetch per paginated request |

**Returns:** _`Promise<OrderInfo[]>`_

all orders, their hash and their fillableTakerAssetAmount

#### getStatsAsync

▸ **getStatsAsync**\(\): _`Promise<GetStatsResponse>`_

_Defined in_ [_ws\_client.ts:116_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L116)

**Returns:** _`Promise<GetStatsResponse>`_

#### onClose

▸ **onClose**\(`cb`: function\): _void_

_Defined in_ [_ws\_client.ts:195_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L195)

Get notified when the underlying WS connection closes normally. If it closes with an error, WSClient automatically attempts to re-connect without emitting a `close` event.

**Parameters:**

▪ **cb**: _function_

callback to call when WS connection closes

▸ \(\): _void_

**Returns:** _void_

#### onReconnected

▸ **onReconnected**\(`cb`: function\): _void_

_Defined in_ [_ws\_client.ts:204_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L204)

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: _function_

callback to call with the error when it occurs

▸ \(\): _void_

**Returns:** _void_

#### subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**\(`cb`: function\): _`Promise<string>`_

_Defined in_ [_ws\_client.ts:156_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L156)

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

_Defined in_ [_ws\_client.ts:185_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/ws_client.ts#L185)

Unsubscribe from a subscription

**Parameters:**

| Name | Type | Description |
| :--- | :--- | :--- |
| `subscriptionId` | string | identifier of the subscription to cancel |

**Returns:** _`Promise<void>`_

> ## Enumeration: OrderEventKind

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

_Defined in_ [_types.ts:55_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L55)

#### Cancelled

• **Cancelled**: = "CANCELLED"

_Defined in_ [_types.ts:58_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L58)

#### Expired

• **Expired**: = "EXPIRED"

_Defined in_ [_types.ts:59_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L59)

#### FillabilityIncreased

• **FillabilityIncreased**: = "FILLABILITY\_INCREASED"

_Defined in_ [_types.ts:61_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L61)

#### Filled

• **Filled**: = "FILLED"

_Defined in_ [_types.ts:56_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L56)

#### FullyFilled

• **FullyFilled**: = "FULLY\_FILLED"

_Defined in_ [_types.ts:57_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L57)

#### Invalid

• **Invalid**: = "INVALID"

_Defined in_ [_types.ts:54_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L54)

#### Unfunded

• **Unfunded**: = "UNFUNDED"

_Defined in_ [_types.ts:60_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L60)

> ## Enumeration: RejectedCode

### Index

#### Enumeration members

* [InternalError](reference.md#internalerror)
* [MaxOrderSizeExceeded](reference.md#maxordersizeexceeded)
* [NetworkRequestFailed](reference.md#networkrequestfailed)
* [OrderAlreadyStored](reference.md#orderalreadystored)
* [OrderCancelled](reference.md#ordercancelled)
* [OrderExpired](reference.md#orderexpired)
* [OrderForIncorrectNetwork](reference.md#orderforincorrectnetwork)
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

_Defined in_ [_types.ts:123_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L123)

#### MaxOrderSizeExceeded

• **MaxOrderSizeExceeded**: = "MaxOrderSizeExceeded"

_Defined in_ [_types.ts:124_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L124)

#### NetworkRequestFailed

• **NetworkRequestFailed**: = "NetworkRequestFailed"

_Defined in_ [_types.ts:127_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L127)

#### OrderAlreadyStored

• **OrderAlreadyStored**: = "OrderAlreadyStored"

_Defined in_ [_types.ts:125_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L125)

#### OrderCancelled

• **OrderCancelled**: = "OrderCancelled"

_Defined in_ [_types.ts:132_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L132)

#### OrderExpired

• **OrderExpired**: = "OrderExpired"

_Defined in_ [_types.ts:130_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L130)

#### OrderForIncorrectNetwork

• **OrderForIncorrectNetwork**: = "OrderForIncorrectNetwork"

_Defined in_ [_types.ts:126_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L126)

#### OrderFullyFilled

• **OrderFullyFilled**: = "OrderFullyFilled"

_Defined in_ [_types.ts:131_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L131)

#### OrderHasInvalidMakerAssetAmount

• **OrderHasInvalidMakerAssetAmount**: = "OrderHasInvalidMakerAssetAmount"

_Defined in_ [_types.ts:128_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L128)

#### OrderHasInvalidMakerAssetData

• **OrderHasInvalidMakerAssetData**: = "OrderHasInvalidMakerAssetData"

_Defined in_ [_types.ts:134_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L134)

#### OrderHasInvalidSignature

• **OrderHasInvalidSignature**: = "OrderHasInvalidSignature"

_Defined in_ [_types.ts:136_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L136)

#### OrderHasInvalidTakerAssetAmount

• **OrderHasInvalidTakerAssetAmount**: = "OrderHasInvalidTakerAssetAmount"

_Defined in_ [_types.ts:129_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L129)

#### OrderHasInvalidTakerAssetData

• **OrderHasInvalidTakerAssetData**: = "OrderHasInvalidTakerAssetData"

_Defined in_ [_types.ts:135_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L135)

#### OrderUnfunded

• **OrderUnfunded**: = "OrderUnfunded"

_Defined in_ [_types.ts:133_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L133)

> ## Enumeration: RejectedKind

### Index

#### Enumeration members

* [MeshError](reference.md#mesherror)
* [MeshValidation](reference.md#meshvalidation)
* [ZeroexValidation](reference.md#zeroexvalidation)

### Enumeration members

#### MeshError

• **MeshError**: = "MESH\_ERROR"

_Defined in_ [_types.ts:118_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L118)

#### MeshValidation

• **MeshValidation**: = "MESH\_VALIDATION"

_Defined in_ [_types.ts:119_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L119)

#### ZeroexValidation

• **ZeroexValidation**: = "ZEROEX\_VALIDATION"

_Defined in_ [_types.ts:117_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L117)

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

_Defined in_ [_types.ts:100_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L100)

#### isNew

• **isNew**: _boolean_

_Defined in_ [_types.ts:101_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L101)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:98_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L98)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:99_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L99)

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

_Defined in_ [_types.ts:14_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L14)

#### `Optional` closeTimeout

• **closeTimeout**? : _undefined \| number_

_Defined in_ [_types.ts:15_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L15)

#### `Optional` fragmentOutgoingMessages

• **fragmentOutgoingMessages**? : _undefined \| false \| true_

_Defined in_ [_types.ts:12_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L12)

#### `Optional` fragmentationThreshold

• **fragmentationThreshold**? : _undefined \| number_

_Defined in_ [_types.ts:13_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L13)

#### `Optional` maxReceivedFrameSize

• **maxReceivedFrameSize**? : _undefined \| number_

_Defined in_ [_types.ts:10_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L10)

#### `Optional` maxReceivedMessageSize

• **maxReceivedMessageSize**? : _undefined \| number_

_Defined in_ [_types.ts:11_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L11)

#### `Optional` tlsOptions

• **tlsOptions**? : _any_

_Defined in_ [_types.ts:16_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L16)

#### `Optional` webSocketVersion

• **webSocketVersion**? : _undefined \| number_

_Defined in_ [_types.ts:9_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L9)

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

_Defined in_ [_types.ts:170_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L170)

#### snapshotID

• **snapshotID**: _string_

_Defined in_ [_types.ts:169_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L169)

> ## Interface: GetStatsResponse

### Hierarchy

* **GetStatsResponse**

### Index

#### Properties

* [EthereumNetworkID](reference.md#ethereumnetworkid)
* [LatestBlock](reference.md#latestblock)
* [NumOrders](reference.md#numorders)
* [NumPeers](reference.md#numpeers)
* [PeerID](reference.md#peerid)
* [PubSubTopic](reference.md#pubsubtopic)
* [Rendezvous](reference.md#rendezvous)
* [Version](reference.md#version)

### Properties

#### EthereumNetworkID

• **EthereumNetworkID**: _number_

_Defined in_ [_types.ts:188_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L188)

#### LatestBlock

• **LatestBlock**: [_LatestBlock_](reference.md#interface-latestblock)

_Defined in_ [_types.ts:189_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L189)

#### NumOrders

• **NumOrders**: _number_

_Defined in_ [_types.ts:191_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L191)

#### NumPeers

• **NumPeers**: _number_

_Defined in_ [_types.ts:190_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L190)

#### PeerID

• **PeerID**: _string_

_Defined in_ [_types.ts:187_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L187)

#### PubSubTopic

• **PubSubTopic**: _string_

_Defined in_ [_types.ts:185_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L185)

#### Rendezvous

• **Rendezvous**: _string_

_Defined in_ [_types.ts:186_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L186)

#### Version

• **Version**: _string_

_Defined in_ [_types.ts:184_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L184)

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

_Defined in_ [_types.ts:71_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L71)

#### subscription

• **subscription**: _string_

_Defined in_ [_types.ts:70_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L70)

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

_Defined in_ [_types.ts:180_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L180)

#### number

• **number**: _number_

_Defined in_ [_types.ts:179_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L179)

> ## Interface: OrderEvent

### Hierarchy

* **OrderEvent**

### Index

#### Properties

* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [kind](reference.md#kind)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)
* [txHashes](reference.md#txhashes)

### Properties

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _`BigNumber`_

_Defined in_ [_types.ts:86_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L86)

#### kind

• **kind**: [_OrderEventKind_](reference.md#enumeration-ordereventkind)

_Defined in_ [_types.ts:85_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L85)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:83_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L83)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:84_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L84)

#### txHashes

• **txHashes**: _string\[\]_

_Defined in_ [_types.ts:87_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L87)

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

_Defined in_ [_types.ts:66_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L66)

#### subscription

• **subscription**: _string_

_Defined in_ [_types.ts:65_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L65)

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

_Defined in_ [_types.ts:113_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L113)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:111_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L111)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:112_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L112)

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

_Defined in_ [_types.ts:93_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L93)

#### isNew

• **isNew**: _boolean_

_Defined in_ [_types.ts:94_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L94)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:91_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L91)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:92_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L92)

> ## Interface: RawOrderEvent

### Hierarchy

* **RawOrderEvent**

### Index

#### Properties

* [fillableTakerAssetAmount](reference.md#fillabletakerassetamount)
* [kind](reference.md#kind)
* [orderHash](reference.md#orderhash)
* [signedOrder](reference.md#signedorder)
* [txHashes](reference.md#txhashes)

### Properties

#### fillableTakerAssetAmount

• **fillableTakerAssetAmount**: _string_

_Defined in_ [_types.ts:78_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L78)

#### kind

• **kind**: [_OrderEventKind_](reference.md#enumeration-ordereventkind)

_Defined in_ [_types.ts:77_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L77)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:75_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L75)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:76_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L76)

#### txHashes

• **txHashes**: _string\[\]_

_Defined in_ [_types.ts:79_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L79)

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

_Defined in_ [_types.ts:107_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L107)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:105_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L105)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:106_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L106)

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

_Defined in_ [_types.ts:147_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L147)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:145_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L145)

#### signedOrder

• **signedOrder**: [_StringifiedSignedOrder_](reference.md#interface-stringifiedsignedorder)

_Defined in_ [_types.ts:146_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L146)

#### status

• **status**: [_RejectedStatus_](reference.md#interface-rejectedstatus)

_Defined in_ [_types.ts:148_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L148)

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

_Defined in_ [_types.ts:159_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L159)

#### rejected

• **rejected**: [_RawRejectedOrderInfo_](reference.md#interface-rawrejectedorderinfo)_\[\]_

_Defined in_ [_types.ts:160_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L160)

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

_Defined in_ [_types.ts:154_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L154)

#### orderHash

• **orderHash**: _string_

_Defined in_ [_types.ts:152_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L152)

#### signedOrder

• **signedOrder**: _`SignedOrder`_

_Defined in_ [_types.ts:153_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L153)

#### status

• **status**: [_RejectedStatus_](reference.md#interface-rejectedstatus)

_Defined in_ [_types.ts:155_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L155)

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

_Defined in_ [_types.ts:140_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L140)

#### message

• **message**: _string_

_Defined in_ [_types.ts:141_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L141)

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

_Defined in_ [_types.ts:47_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L47)

#### expirationTimeSeconds

• **expirationTimeSeconds**: _string_

_Defined in_ [_types.ts:49_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L49)

#### feeRecipientAddress

• **feeRecipientAddress**: _string_

_Defined in_ [_types.ts:48_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L48)

#### makerAddress

• **makerAddress**: _string_

_Defined in_ [_types.ts:38_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L38)

#### makerAssetAmount

• **makerAssetAmount**: _string_

_Defined in_ [_types.ts:42_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L42)

#### makerAssetData

• **makerAssetData**: _string_

_Defined in_ [_types.ts:44_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L44)

#### makerFee

• **makerFee**: _string_

_Defined in_ [_types.ts:40_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L40)

#### salt

• **salt**: _string_

_Defined in_ [_types.ts:46_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L46)

#### senderAddress

• **senderAddress**: _string_

_Defined in_ [_types.ts:37_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L37)

#### signature

• **signature**: _string_

_Defined in_ [_types.ts:50_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L50)

#### takerAddress

• **takerAddress**: _string_

_Defined in_ [_types.ts:39_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L39)

#### takerAssetAmount

• **takerAssetAmount**: _string_

_Defined in_ [_types.ts:43_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L43)

#### takerAssetData

• **takerAssetData**: _string_

_Defined in_ [_types.ts:45_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L45)

#### takerFee

• **takerFee**: _string_

_Defined in_ [_types.ts:41_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L41)

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

_Defined in_ [_types.ts:164_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L164)

#### rejected

• **rejected**: [_RejectedOrderInfo_](reference.md#interface-rejectedorderinfo)_\[\]_

_Defined in_ [_types.ts:165_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L165)

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

_Defined in_ [_types.ts:174_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L174)

#### utf8Data

• **utf8Data**: _string_

_Defined in_ [_types.ts:175_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L175)

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

_Defined in_ [_types.ts:32_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L32)

#### `Optional` headers

• **headers**? : _undefined \| `__type`_

_Defined in_ [_types.ts:30_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L30)

#### `Optional` protocol

• **protocol**? : _undefined \| string_

_Defined in_ [_types.ts:31_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L31)

#### `Optional` reconnectAfter

• **reconnectAfter**? : _undefined \| number_

_Defined in_ [_types.ts:33_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L33)

#### `Optional` timeout

• **timeout**? : _undefined \| number_

_Defined in_ [_types.ts:29_](https://github.com/0xProject/0x-mesh/blob/db314ef/rpc/clients/typescript/src/types.ts#L29)

