# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:212](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L212)*

Instantiates a new Mesh instance.

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`config` | [Config](#interface-config) | Configuration options for Mesh |

**Returns:** *[Mesh](#class-mesh)*

An instance of Mesh

### Methods

##  addOrdersAsync

▸ **addOrdersAsync**(`orders`: Array‹SignedOrder›): *Promise‹[ValidationResults](#interface-validationresults)›*

*Defined in [index.ts:291](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L291)*

Validates and adds the given orders to Mesh. If an order is successfully
added, Mesh will share it with any peers in the network and start
watching it for changes (e.g. filled, canceled, expired). The returned
promise will only be rejected if there was an error validating or adding
the order; it will not be rejected for any invalid orders (check
results.rejected instead).

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`orders` | Array‹SignedOrder› | An array of orders to add. |

**Returns:** *Promise‹[ValidationResults](#interface-validationresults)›*

Validation results for the given orders, indicating which orders
were accepted and which were rejected.

___

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [index.ts:241](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L241)*

Registers a handler which will be called in the event of a critical
error. Note that the handler will not be called for non-critical errors.
In order to ensure no errors are missed, this should be called before
startAsync.

**Parameters:**

▪ **handler**: *function*

The handler to be called.

▸ (`err`: Error): *void*

**Parameters:**

Name | Type |
------ | ------ |
`err` | Error |

**Returns:** *void*

___

##  onOrderEvents

▸ **onOrderEvents**(`handler`: function): *void*

*Defined in [index.ts:256](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L256)*

Registers a handler which will be called for any incoming order events.
Order events are fired whenver an order is added, canceled, expired, or
filled. In order to ensure no events are missed, this should be called
before startAsync.

**Parameters:**

▪ **handler**: *function*

The handler to be called.

▸ (`events`: Array‹[OrderEvent](#interface-orderevent)›): *void*

**Parameters:**

Name | Type |
------ | ------ |
`events` | Array‹[OrderEvent](#interface-orderevent)› |

**Returns:** *void*

___

##  startAsync

▸ **startAsync**(): *Promise‹void›*

*Defined in [index.ts:267](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L267)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:173](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L173)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:171](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L171)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:172](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L172)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:170](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L170)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:151](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L151)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:152](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L152)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:149](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L149)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:150](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L150)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:47](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L47)*

___

##  ethereumNetworkID

• **ethereumNetworkID**: *number*

*Defined in [index.ts:33](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L33)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:56](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L56)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:31](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L31)*

___

## `Optional` orderExpirationBufferSeconds

• **orderExpirationBufferSeconds**? : *undefined | number*

*Defined in [index.ts:40](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L40)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:36](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L36)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:111](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L111)*

___

##  kind

• **kind**: *string*

*Defined in [index.ts:110](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L110)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:108](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L108)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:109](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L109)*

___

##  txHashes

• **txHashes**: *Array‹string›*

*Defined in [index.ts:112](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L112)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:162](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L162)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:160](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L160)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:161](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L161)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:163](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L163)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:180](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L180)*

___

##  message

• **message**: *string*

*Defined in [index.ts:181](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L181)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *Array‹[AcceptedOrderInfo](#class-acceptedorderinfo)›*

*Defined in [index.ts:141](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L141)*

___

##  rejected

• **rejected**: *Array‹[RejectedOrderInfo](#class-rejectedorderinfo)›*

*Defined in [index.ts:142](https://github.com/0xProject/0x-mesh/blob/ed8a11c/browser/ts/index.ts#L142)*

<hr />

