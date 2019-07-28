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

* [clientConfig](_types_.wsopts.md#optional-clientconfig)
* [headers](_types_.wsopts.md#optional-headers)
* [protocol](_types_.wsopts.md#optional-protocol)
* [reconnectAfter](_types_.wsopts.md#optional-reconnectafter)
* [timeout](_types_.wsopts.md#optional-timeout)

## Properties

### `Optional` clientConfig

• **clientConfig**? : *[ClientConfig](_types_.clientconfig.md)*

*Defined in [types.ts:32](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/types.ts#L32)*

___

### `Optional` headers

• **headers**? : *undefined | `__type`*

*Defined in [types.ts:30](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/types.ts#L30)*

___

### `Optional` protocol

• **protocol**? : *undefined | string*

*Defined in [types.ts:31](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/types.ts#L31)*

___

### `Optional` reconnectAfter

• **reconnectAfter**? : *undefined | number*

*Defined in [types.ts:33](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/types.ts#L33)*

___

### `Optional` timeout

• **timeout**? : *undefined | number*

*Defined in [types.ts:29](https://github.com/0xProject/0x-mesh/blob/9ff2bf1/rpc/clients/typescript/src/types.ts#L29)*