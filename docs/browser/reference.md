# Class: Mesh

The main class for this package. Has methods for receiving order events and
sending orders through the 0x Mesh network.

### Hierarchy

* **Mesh**


### Constructors

##  constructer

\+ **new Mesh**(`config`: [Config](#interface-config)): *[Mesh](#class-mesh)*

*Defined in [index.ts:253](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L253)*

Instantiates a new Mesh instance.

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`config` | [Config](#interface-config) | Configuration options for Mesh |

**Returns:** *[Mesh](#class-mesh)*

An instance of Mesh

### Methods

##  addOrdersAsync

▸ **addOrdersAsync**(`orders`: SignedOrder[]): *Promise‹[ValidationResults](#interface-validationresults)›*

*Defined in [index.ts:323](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L323)*

Validates and adds the given orders to Mesh. If an order is successfully
added, Mesh will share it with any peers in the network and start
watching it for changes (e.g. filled, canceled, expired). The returned
promise will only be rejected if there was an error validating or adding
the order; it will not be rejected for any invalid orders (check
results.rejected instead).

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
`orders` | SignedOrder[] | An array of orders to add. |

**Returns:** *Promise‹[ValidationResults](#interface-validationresults)›*

Validation results for the given orders, indicating which orders
were accepted and which were rejected.

___

##  onError

▸ **onError**(`handler`: function): *void*

*Defined in [index.ts:273](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L273)*

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

*Defined in [index.ts:288](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L288)*

Registers a handler which will be called for any incoming order events.
Order events are fired whenver an order is added, canceled, expired, or
filled. In order to ensure no events are missed, this should be called
before startAsync.

**Parameters:**

▪ **handler**: *function*

The handler to be called.

▸ (`events`: [OrderEvent](#interface-orderevent)[]): *void*

**Parameters:**

Name | Type |
------ | ------ |
`events` | [OrderEvent](#interface-orderevent)[] |

**Returns:** *void*

___

##  startAsync

▸ **startAsync**(): *Promise‹void›*

*Defined in [index.ts:299](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L299)*

Starts the Mesh node in the background. Mesh will automatically find
peers in the network and begin receiving orders from them.

**Returns:** *Promise‹void›*

<hr />

# Enumeration: RejectedOrderKind

A set of categories for rejected orders.


### Enumeration members

##  CoordinatorError

• **CoordinatorError**: = "COORDINATOR_ERROR"

*Defined in [index.ts:207](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L207)*

___

##  MeshError

• **MeshError**: = "MESH_ERROR"

*Defined in [index.ts:205](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L205)*

___

##  MeshValidation

• **MeshValidation**: = "MESH_VALIDATION"

*Defined in [index.ts:206](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L206)*

___

##  ZeroExValidation

• **ZeroExValidation**: = "ZEROEX_VALIDATION"

*Defined in [index.ts:204](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L204)*

<hr />

# Interface: AcceptedOrderInfo

Info for any orders that were accepted.

### Hierarchy

* **AcceptedOrderInfo**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:185](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L185)*

___

##  isNew

• **isNew**: *boolean*

*Defined in [index.ts:186](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L186)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:183](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L183)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:184](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L184)*

<hr />

# Interface: Config

A set of configuration options for Mesh.

### Hierarchy

* **Config**


### Properties

## `Optional` blockPollingIntervalSeconds

• **blockPollingIntervalSeconds**? : *undefined | number*

*Defined in [index.ts:59](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L59)*

___

## `Optional` bootstrapList

• **bootstrapList**? : *string[]*

*Defined in [index.ts:48](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L48)*

___

##  ethereumNetworkID

• **ethereumNetworkID**: *number*

*Defined in [index.ts:40](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L40)*

___

## `Optional` ethereumRPCMaxContentLength

• **ethereumRPCMaxContentLength**? : *undefined | number*

*Defined in [index.ts:68](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L68)*

___

##  ethereumRPCURL

• **ethereumRPCURL**: *string*

*Defined in [index.ts:38](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L38)*

___

## `Optional` orderExpirationBufferSeconds

• **orderExpirationBufferSeconds**? : *undefined | number*

*Defined in [index.ts:52](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L52)*

___

## `Optional` useBootstrapList

• **useBootstrapList**? : *undefined | false | true*

*Defined in [index.ts:43](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L43)*

___

##  verbosity

• **verbosity**: *Verbosity*

*Defined in [index.ts:35](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L35)*

<hr />

# Interface: OrderEvent

Order events are fired by Mesh whenever an order is added, canceled, expired,
or filled.

### Hierarchy

* **OrderEvent**


### Properties

##  fillableTakerAssetAmount

• **fillableTakerAssetAmount**: *BigNumber*

*Defined in [index.ts:145](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L145)*

___

##  kind

• **kind**: *string*

*Defined in [index.ts:144](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L144)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:142](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L142)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:143](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L143)*

___

##  txHashes

• **txHashes**: *string[]*

*Defined in [index.ts:146](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L146)*

<hr />

# Interface: RejectedOrderInfo

Info for any orders that were rejected, including the reason they were
rejected.

### Hierarchy

* **RejectedOrderInfo**


### Properties

##  kind

• **kind**: *[RejectedOrderKind](#enumeration-rejectedorderkind)*

*Defined in [index.ts:196](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L196)*

___

##  orderHash

• **orderHash**: *string*

*Defined in [index.ts:194](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L194)*

___

##  signedOrder

• **signedOrder**: *SignedOrder*

*Defined in [index.ts:195](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L195)*

___

##  status

• **status**: *[RejectedOrderStatus](#class-rejectedorderstatus)*

*Defined in [index.ts:197](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L197)*

<hr />

# Interface: RejectedOrderStatus

Provides more information about why an order was rejected.

### Hierarchy

* **RejectedOrderStatus**


### Properties

##  code

• **code**: *string*

*Defined in [index.ts:214](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L214)*

___

##  message

• **message**: *string*

*Defined in [index.ts:215](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L215)*

<hr />

# Interface: ValidationResults

Indicates which orders where accepted, which were rejected, and why.

### Hierarchy

* **ValidationResults**


### Properties

##  accepted

• **accepted**: *[AcceptedOrderInfo](#class-acceptedorderinfo)[]*

*Defined in [index.ts:175](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L175)*

___

##  rejected

• **rejected**: *[RejectedOrderInfo](#class-rejectedorderinfo)[]*

*Defined in [index.ts:176](https://github.com/0xProject/0x-mesh/blob/22c2a55/browser/ts/index.ts#L176)*

<hr />

