# Custom Order Filters

Mesh supports the creation of separate sub-networks where 0x orders that adhere to a specific schema are shared. Each sub-network is built around a custom order filter. The custom filter defines which orders are allowed to be shared within a sub-network. For example:

- All orders for a specific asset pair (e.g., ETH/DAI)
- All orders for non-fungibles (i.e., ERC721, ERC1155)
- All orders used by a specific DApp

A custom filter may be passed into Mesh as a [JSON Schema](https://json-schema.org/) via the `CUSTOM_ORDER_FILTER` environment variable. Messages that contain orders that don't match this schema will be dropped. As a limitation, filtering is only possible by looking at the static fields of an order. So for example, it is not possible to filter orders by doing an on-chain check or sending an HTTP request to a third-party API. We don't expect that this limitation is going to be a problem in practice and it comes with the huge benefit of enabling cross-topic forwarding in the future (more on that later).

### New order and message schemas.

All orders must match the following JSON Schema:

```json
{
	"id": "/rootOrder",
	"allOf": [{
		"$ref": "/customOrder"
	}, {
		"$ref": "/signedOrder"
	}]
}
```

- `/signedOrder` is the JSON Schema that will match any valid 0x orders.
- `/customOrder` is the custom schema passed in through the `CUSTOM_ORDER_FILTER` environment variable.

Organizing the JSON Schema for orders like this means that `CUSTOM_ORDER_FILTER` can be relatively small. It doesn't need to contain all the required fields for a signed 0x order. It just needs to contain any _additional_ requirements on top of the default ones.

#### Example custom order schemas

The following `CUSTOM_ORDER_FILTER` doesn't add any additional requirements. All valid signed 0x orders will be accepted. This is the default value if no custom filter is passed in.

```json
{}
```

The following `CUSTOM_ORDER_FILTER` matches any valid signed 0x orders with a specific sender address:

```json
{
	"properties": {
		"senderAddress": {
			"pattern": "0x00000000000000000000000000000000ba5eba11",
			"type": "string"
		}
	}
}
```

This can easily be tweaked to filter orders by asset type, maker/taker address, or fee recipient. JSON Schema has support for [regular expressions](https://json-schema.org/understanding-json-schema/reference/regular_expressions.html) allowing for partial matching of any 0x order field.

### Limitations

Nodes that are spun up with a custom filter will share all their orders with nodes that are either using the same filter or the default "all" filter. They will not share orders with nodes using other custom filters since each filter results in a separate sub-network. Therefore, custom filters are most useful for applications where users care about a distinct sub-set of 0x orders. 

If you wanted to connect two sub-networks with overlapping valid orders, you could spin up a Mesh node for each and additionally run a [bridge script](https://github.com/0xProject/0x-mesh/blob/master/cmd/mesh-bridge/main.go) to sends orders from one sub-network to the other. Longer term, we hope to add support for cross-topic forwarding which will allow Mesh nodes to do this under-the-hood.
