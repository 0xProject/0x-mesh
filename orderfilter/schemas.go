package orderfilter

const (
	// Built-in schemas
	addressSchema     = `{"$id":"/address","type":"string","pattern":"^0x[0-9a-fA-F]{40}$"}`
	wholeNumberSchema = `{"$id":"/wholeNumber","anyOf":[{"type":"string","pattern":"^\\d+$"},{"type":"integer"}]}`
	hexSchema         = `{"$id":"/hex","type":"string","pattern":"^0x(([0-9a-fA-F][0-9a-fA-F])+)?$"}`

	// V3 Order schemas
	orderSchemaV3 = `{
  "$id": "/orderV3",
  "properties": {
    "makerAddress": {
      "$ref": "/address"
    },
    "takerAddress": {
      "$ref": "/address"
    },
    "makerFee": {
      "$ref": "/wholeNumber"
    },
    "takerFee": {
      "$ref": "/wholeNumber"
    },
    "senderAddress": {
      "$ref": "/address"
    },
    "makerAssetAmount": {
      "$ref": "/wholeNumber"
    },
    "takerAssetAmount": {
      "$ref": "/wholeNumber"
    },
    "makerAssetData": {
      "$ref": "/hex"
    },
    "takerAssetData": {
      "$ref": "/hex"
    },
    "makerFeeAssetData": {
      "$ref": "/hex"
    },
    "takerFeeAssetData": {
      "$ref": "/hex"
    },
    "salt": {
      "$ref": "/wholeNumber"
    },
    "feeRecipientAddress": {
      "$ref": "/address"
    },
    "expirationTimeSeconds": {
      "$ref": "/wholeNumber"
    },
    "exchangeAddress": {
      "$ref": "/exchangeAddress"
    },
    "chainId": {
      "$ref": "/chainId"
    }
  },
  "required": [
    "makerAddress",
    "takerAddress",
    "makerFee",
    "takerFee",
    "senderAddress",
    "makerAssetAmount",
    "takerAssetAmount",
    "makerAssetData",
    "takerAssetData",
    "makerFeeAssetData",
    "takerFeeAssetData",
    "salt",
    "feeRecipientAddress",
    "expirationTimeSeconds",
    "exchangeAddress",
    "chainId"
  ],
  "type": "object"
}`
	signedOrderSchemaV3 = `{"$id":"/signedOrderV3","allOf":[{"$ref":"/orderV3"},{"properties":{"signature":{"$ref":"/hex"}},"required":["signature"]}]}`

	// V4 Order schemas
	orderSchemaV4       = `{"$id":"/orderV4","properties":{"maker":{"$ref":"/address"},"taker":{"$ref":"/address"},"makerFee":{"$ref":"/wholeNumber"},"takerFee":{"$ref":"/wholeNumber"},"sender":{"$ref":"/address"},"makerAmount":{"$ref":"/wholeNumber"},"takerAmount":{"$ref":"/wholeNumber"},"makerToken":{"$ref":"/address"},"takerToken":{"$ref":"/address"},"salt":{"$ref":"/wholeNumber"},"pool":{"$ref":"/wholeNumber"},"origin":{"$ref":"/wholeNumber"},"feeRecipient":{"$ref":"/address"},"expiry":{"$ref":"/wholeNumber"},"exchange":{"$ref":"/exchange"},"chainId":{"$ref":"/chainId"}},"required":["maker","taker","makerFee","takerFee","sender","makerAmount","takerAmount","makerToken","takerToken","salt","pool","origin","feeRecipient","expiry","exchange","chainId"],"type":"object"}`
	signedOrderSchemaV4 = `{"$id":"/signedOrderV4","allOf":[{"$ref":"/orderV4"},{"properties":{"signature":{"$ref":"/hex"}},"required":["signature"]}]}`

	// V3 Root schemas
	rootOrderSchemaV3        = `{"$id":"/rootOrderV3","allOf":[{"$ref":"/customOrderV3"},{"$ref":"/signedOrderV3"}]}`
	rootOrderMessageSchemaV3 = `{"$id":"/rootOrderMessageV3","properties":{"messageType":{"type":"string","pattern":"order"},"order":{"$ref":"/rootOrderV3"},"topics":{"type":"array","minItems":1,"items":{"type":"string"}}},"required":["messageType","order","topics"]}`

	// V4 Root schemas
	rootOrderSchemaV4        = `{"$id":"/rootOrderV4","allOf":[{"$ref":"/customOrderV4"},{"$ref":"/signedOrderV4"}]}`
	rootOrderMessageSchemaV4 = `{"$id":"/rootOrderMessageV4","properties":{"messageType":{"type":"string","pattern":"order"},"order":{"$ref":"/rootOrderV4"},"topics":{"type":"array","minItems":1,"items":{"type":"string"}}},"required":["messageType","order","topics"]}`

	// DefaultCustomOrderSchema is the default schema for /customOrderV3 and
	// /customOrderV4. It includes all 0x orders and doesn't add any
	// additional requirements.
	DefaultCustomOrderSchema = `{}`
)
