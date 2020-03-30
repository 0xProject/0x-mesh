# Custom Order Filters

Mesh supports the creation of separate sub-networks where 0x orders that adhere to a specific schema are shared. Each sub-network is built around a custom order filter. The custom filter defines which orders are allowed to be shared within a sub-network. For example:

- All orders for a specific asset pair (e.g., WETH/DAI)
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

Nodes that are spun up with a custom filter will share all their orders with nodes that are either using the exact same filter or the default "all" filter (i.e., "{}"). They will _not_ share orders with nodes using different custom filters (even if a given order matches both filters) because each filter results in a separate sub-network. Therefore, custom filters are most useful for applications where users care about a distinct subset of 0x orders.

If you wanted to connect two sub-networks with overlapping valid orders, you could spin up a Mesh node for each sub-network and additionally run a [bridge script](https://github.com/0xProject/0x-mesh/blob/master/cmd/mesh-bridge/main.go) to send orders from one sub-network to the other. Longer term, we hope to add support for cross-topic forwarding, which will allow Mesh nodes to do this under-the-hood.

### Examples

##### WETH <-> DAI orders:
```json
{
    "oneOf": [
        {
            "properties": {
                "makerAssetData": {
                    "pattern": "0xf47261b00000000000000000000000006b175474e89094c44da98b954eedeac495271d0f",
                    "type": "string"
                },
                "takerAssetData": {
                    "pattern": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                    "type": "string"
                }
            }
        },
        {
            "properties": {
                "makerAssetData": {
                    "pattern": "0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                    "type": "string"
                },
                "takerAssetData": {
                    "pattern": "0xf47261b00000000000000000000000006b175474e89094c44da98b954eedeac495271d0f",
                    "type": "string"
                }
            }
        }
    ]
}
```

##### Any ERC721 order:

```json
{
    "oneOf": [
        {
            "properties": {
                "makerAssetData": {
                    "pattern": "0x02571792.*",
                    "type": "string"
                }
            }
        },
        {
            "properties": {
                "takerAssetData": {
                    "pattern": "0x02571792.*",
                    "type": "string"
                }
            }
        }
    ]
}
```

##### Augur V2 orders:

```json
{
    "properties": {
        "makerAssetData": {
            "pattern": ".*${AUGUR_ERC1155_CONTRACT_ADDRESS}.*"
        }
    }
}
```

Where `${AUGUR_ERC1155_CONTRACT_ADDRESS}` needs to be replaced with the Augur ERC1155 token used to represent the outcomes of their various prediction markets.
