> # Class: WSClient

This class includes all the functionality related to interacting with a Mesh JSON RPC
websocket endpoint.

## Hierarchy

* **WSClient**

## Index

### Constructors

* [constructor](_ws_client_.wsclient.md#constructor)

### Methods

* [addOrdersAsync](_ws_client_.wsclient.md#addordersasync)
* [destroy](_ws_client_.wsclient.md#destroy)
* [getOrdersAsync](_ws_client_.wsclient.md#getordersasync)
* [onClose](_ws_client_.wsclient.md#onclose)
* [onReconnected](_ws_client_.wsclient.md#onreconnected)
* [subscribeToOrdersAsync](_ws_client_.wsclient.md#subscribetoordersasync)
* [unsubscribeAsync](_ws_client_.wsclient.md#unsubscribeasync)

## Constructors

###  constructor

\+ **new WSClient**(`url`: string, `wsOpts?`: [WSOpts](../interfaces/_types_.wsopts.md)): *[WSClient](_ws_client_.wsclient.md)*

*Defined in [ws_client.ts:55](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L55)*

Instantiates a new WSClient instance

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`url` | string | WS server endpoint |
`wsOpts?` | [WSOpts](../interfaces/_types_.wsopts.md) | WebSocket options |

**Returns:** *[WSClient](_ws_client_.wsclient.md)*

An instance of WSClient

## Methods

###  addOrdersAsync

▸ **addOrdersAsync**(`signedOrders`: `SignedOrder`[]): *`Promise<ValidationResults>`*

*Defined in [ws_client.ts:80](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L80)*

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

*Defined in [ws_client.ts:193](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L193)*

destroy unsubscribes all active subscriptions, closes the websocket connection
and stops the internal heartbeat connection liveness check.

**Returns:** *void*

___

###  getOrdersAsync

▸ **getOrdersAsync**(`perPage`: number): *`Promise<AcceptedOrderInfo[]>`*

*Defined in [ws_client.ts:105](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L105)*

Get all 0x signed orders currently stored in the Mesh node

**Parameters:**

Name | Type | Default | Description |
------ | ------ | ------ | ------ |
`perPage` | number | 200 | number of signedOrders to fetch per paginated request |

**Returns:** *`Promise<AcceptedOrderInfo[]>`*

all orders, their hash and their fillableTakerAssetAmount

___

###  onClose

▸ **onClose**(`cb`: function): *void*

*Defined in [ws_client.ts:175](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L175)*

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

*Defined in [ws_client.ts:184](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L184)*

Get notified when a connection to the underlying WS connection is re-established

**Parameters:**

▪ **cb**: *function*

callback to call with the error when it occurs

▸ (): *void*

**Returns:** *void*

___

###  subscribeToOrdersAsync

▸ **subscribeToOrdersAsync**(`cb`: function): *`Promise<string>`*

*Defined in [ws_client.ts:136](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L136)*

Subscribe to the 'orders' topic and receive order events from Mesh. This method returns a
subscriptionId that can be used to `unsubscribe()` from this subscription.

**Parameters:**

▪ **cb**: *function*

callback function where you'd like to get notified about order events

▸ (`orderEvents`: [OrderEvent](../interfaces/_types_.orderevent.md)[]): *void*

**Parameters:**

Name | Type |
------ | ------ |
`orderEvents` | [OrderEvent](../interfaces/_types_.orderevent.md)[] |

**Returns:** *`Promise<string>`*

subscriptionId

___

###  unsubscribeAsync

▸ **unsubscribeAsync**(`subscriptionId`: string): *`Promise<void>`*

*Defined in [ws_client.ts:165](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/ws_client.ts#L165)*

Unsubscribe from a subscription

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`subscriptionId` | string | identifier of the subscription to cancel  |

**Returns:** *`Promise<void>`*