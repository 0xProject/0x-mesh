package orderfilter

const (
	// Built-in schemas
	addressSchema     = `{"$id":"/address","type":"string","pattern":"^0x[0-9a-fA-F]{40}$"}`
	wholeNumberSchema = `{"$id":"/wholeNumber","anyOf":[{"type":"string","pattern":"^\\d+$"},{"type":"integer"}]}`
	hexSchema         = `{"$id":"/hex","type":"string","pattern":"^0x(([0-9a-fA-F][0-9a-fA-F])+)?$"}`
	orderSchema       = `{"$id":"/order","properties":{"makerAddress":{"$ref":"/address"},"takerAddress":{"$ref":"/address"},"makerFee":{"$ref":"/wholeNumber"},"takerFee":{"$ref":"/wholeNumber"},"senderAddress":{"$ref":"/address"},"makerAssetAmount":{"$ref":"/wholeNumber"},"takerAssetAmount":{"$ref":"/wholeNumber"},"makerAssetData":{"$ref":"/hex"},"takerAssetData":{"$ref":"/hex"},"makerFeeAssetData":{"$ref":"/hex"},"takerFeeAssetData":{"$ref":"/hex"},"salt":{"$ref":"/wholeNumber"},"feeRecipientAddress":{"$ref":"/address"},"expirationTimeSeconds":{"$ref":"/wholeNumber"},"exchangeAddress":{"$ref":"/exchangeAddress"},"chainId":{"$ref":"/chainId"}},"required":["makerAddress","takerAddress","makerFee","takerFee","senderAddress","makerAssetAmount","takerAssetAmount","makerAssetData","takerAssetData","makerFeeAssetData","takerFeeAssetData","salt","feeRecipientAddress","expirationTimeSeconds","exchangeAddress","chainId"],"type":"object"}`
	orderV4Schema     = `
{
    "type": "object",
    "required": [
        "verifyingContract",
        "chainId",
        "makerToken",
        "takerToken",
        "makerAmount",
        "takerAmount",
        "takerTokenFeeAmount",
        "maker",
        "taker",
        "sender",
        "feeRecipient",
        "pool",
        "expiry",
        "salt"
    ],
    "properties": {
        "verifyingContract": {
            "$ref": "/hex"
        },
        "chainId": {
            "$ref": "/chainId"
        },
        "makerToken": {
            "$ref": "/hex"
        },
        "takerToken": {
            "$ref": "/hex"
        },
        "makerAmount": {
            "$ref": "/wholeNumber"
        },
        "takerAmount": {
            "$ref": "/wholeNumber"
        },
        "takerTokenFeeAmount": {
            "$ref": "/wholeNumber"
        },
        "maker": {
            "$ref": "/address"
        },
        "taker": {
            "$ref": "/address"
        },
        "sender": {
            "$ref": "/address"
        },
        "feeRecipient": {
            "$ref": "/address"
        },
        "pool": {
            "$ref": "/hex"
        },
        "expiry": {
            "$ref": "/wholeNumber"
        },
        "salt": {
            "$ref": "/wholeNumber"
        }
    },
    "$id": "/orderv4"
}
`

	signedOrderV4Schema = `
{
    "allOf": [
        {
            "$ref": "/orderv4"
        },
        {
            "properties": {
                "signatureType": {
                    "$ref": "/wholeNumber"
                },
                "signatureV": {
                    "$ref": "/wholeNumber"
                },
                "signatureR": {
                    "$ref": "/hex"
                },
                "signatureS": {
                    "$ref": "/hex"
                }
            },
            "required": [
                "signatureType",
                "signatureV",
                "signatureR",
                "signatureS"
            ]
        }
    ],
    "$id": "/signedOrderV4"
}
`
	signedOrderSchema = `{
    "allOf": [
        {
            "$ref": "/order"
        },
        {
            "required": [
                "signature"
            ],
            "properties": {
                "signature": {
                    "$ref": "/hex"
                }
            }
        }
    ],
    "$id": "/signedOrder"
}`

	// Root schemas
	rootOrderSchema        = `{"$id":"/rootOrder","allOf":[{"$ref":"/customOrder"},{"$ref":"/signedOrder"}]}`
	rootOrderV4Schema      = `{"$id":"/rootOrder","anyOf":[{"$ref":"/orderv4"}]}`
	rootOrderMessageSchema = `{"$id":"/rootOrderMessage","properties":{"messageType":{"type":"string","pattern":"order"},"order":{"$ref":"/rootOrder"},"topics":{"type":"array","minItems":1,"items":{"type":"string"}}},"required":["messageType","order","topics"]}`

	// DefaultCustomOrderSchema is the default schema for /customOrder. It
	// includes all 0x orders and doesn't add any additional requirements.
	DefaultCustomOrderSchema = `{}`
)
