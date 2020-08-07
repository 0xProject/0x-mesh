# How to keep an external database in-sync with a Mesh node

This guide will walk you through syncing your own database with a Mesh node so that your database's
state mirrors that of the Mesh node (and vice-versa). Whenever new orders are discovered or added to
Mesh, they should be inserted into your database. If an order is filled, cancelled, or has its fillability
changed, you should update or remove it from your database. We are assuming that your database is storing
both the order itself and associated metadata (e.g., `fillableTakerAssetAmount`).

This guide fairly advanced, and is specifically written for developers who _need_ to sync their own
database with Mesh (e.g. in order to join Mesh orders with application-specific data). Note that with
the introduction of the [GraphQL API](graphql_api.md), most developers do not need to worry about database
syncing. Instead, you can think of Mesh itself as the database for orders and query it directly. If you
aren't sure if this guide is for you, try using the [GraphQL API](graphql_api.md) first.

## Standalone or Browser nodes

This guide is written with both standalone and browser nodes in mind and the process for syncing
your database is very similar.

If you are running a standalone node, you will interact with it via the [GraphQL API](graphql_api.md).

If you are running Mesh directly in the browser via the `@0x/mesh-browser` or `@0x/mesh-browser-lite`
packages, you will interact with it using the [TypeScript/JavaScript API](browser-bindings/browser/reference.md).

## Initial sync

When first connecting the DB and Mesh node, we need to make sure both have the same orders and
order-relevant state stored. We do this with the following steps:

#### 1. Subscribe to Mesh

Subscribe to the Mesh node's [`orderEvents` subscription](graphql_api.md#subscribing-to-order-events).

If you are using our TypeScript GraphQL client, you can use the [`onOrderEvents`](browser-bindings/browser/reference.md#onorderevents) method.

If you are using the `@0x/mesh-browser` or `@0x/mesh-browser-lite` packages, you use the method by
by the same name, [`onOrderEvents`](browser-bindings/browser/reference.md#onorderevents).

Whenever you receive an order event from this subscription, make the appropriate updates to your DB. Each
order event has an associated [OrderEventEndState](https://godoc.org/github.com/0xProject/0x-mesh/zeroex#pkg-constants).

| End state                                                    | DB operation     |
| ------------------------------------------------------------ | ---------------- |
| ADDED, FILLED, FILLABILITY_INCREASED, UNEXPIRED              | Insert or Update |
| FULLY_FILLED, EXPIRED, CANCELLED, UNFUNDED, STOPPED_WATCHING | Remove           |

**Note:** If you receive any event other than `ADDED`, `FILLABILITY_INCREASED`, or `UNEXPIRED`
for an order we do not find in our database, we ignore the event and noop.

#### 2. Get all orders currently stored in Mesh

There might have been orders stored in Mesh that your DB doesn't know about at this time. Because
of this, you should fetch all currently stored orders in the Mesh node and upsert them in the database.
This can be done using the [orders](graphql_api.md#querying-and-filtering-orders) GraphQL query.

If you are using our TypeScript GraphQL client, you can use the
[`getOrdersAsync`](graphql_clients/typescript/reference.md#getordersasync) method.

If you are using the `@0x/mesh-browser` or `@0x/mesh-browser-lite` packages, you can use
the method by the same name, [`getordersAsync`](browser-bindings/browser/reference.md#getordersasync).

Orders may be added or removed while you are getting existing orders from the Mesh DB. For this reason,
it is important to account for any order events received from step (1) while or after you get the existing
orders in this step.

**Note:** The [Mesh Typescript client](graphql_clients/typescript/README.md) has a convenience method
that does the multiple paginated requests for you under-the-hood. You can simply call the
[getOrders](graphql_clients/typescript/reference.md#getordersasync) method.

#### 3. Add all database orders to the Mesh node

Since there might also be orders in your database that Mesh doesn't know about, you should also
add those orders to Mesh. We can do this using the [addOrders](graphql_api.md#adding-orders)
GraphQL mutation.

If you are using our TypeScript GraphQL client, you can use the
[`addOrdersAsync`](graphql_clients/typescript/reference.md#addordersasync) method.

If you are using the `@0x/mesh-browser` or `@0x/mesh-browser-lite` packages, you can use
the method by the same name, [`addOrdersAsync`](browser-bindings/browser/reference.md#addordersasync).

This method accepts an array of signed 0x orders and returns which have been accepted
and rejected. The accepted orders are returned with their `fillableTakerAssetAmount` and so these
amounts should be updated in the database. Rejected orders are rejected with a specific
[RejectedOrderStatus](https://godoc.org/github.com/0xProject/0x-mesh/zeroex/ordervalidator), including
an identifying `code`. The following codes indicate temporary errors and you may try submitting the
order again (typically with exponential backoff): `INTERNAL_ERROR`, `ETH_RPC_REQUEST_FAILED`, or
`DATABASE_FULL_OF_ORDERS`. For any other code, the order has been rejected by Mesh and should not
be retried. If the order exists in your database is should be removed.

#### 4. Handle dropped connections

After performing the first 3 steps above, your database will be in-sync with the Mesh database, and continue to remain
in-sync thanks to the active order event subscription. If any new orders are added to the database, they will also need
to be added to Mesh of course. But what if the WebSocket connection to the Mesh node goes down? In that case, it
must be re-established and steps 1, 2 & 3 must be performed once again.

**Note:** With some WebSocket clients, we've noticed that the client is not always aware of when the connection has been
dropped. It can be hard for clients to discern between a network disruption and latency. Because of this, our GraphQL server
uses `GQL_CONNECTION_ACK` and `GQL_CONNECTION_KEEP_ALIVE` messages as described in
https://github.com/apollographql/subscriptions-transport-ws/blob/master/PROTOCOL.md. Our
[Typescript GraphQL client](graphql_clients/typescript) automatically handles these messages, and
most other GraphQL clients will too.

Happy database syncing!
