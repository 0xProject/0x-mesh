package core

import (
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/xeipuuv/gojsonschema"
)

// JSON-schema schemas
var (
	addressSchemaLoader     = gojsonschema.NewStringLoader(`{"id":"/addressSchema","type":"string","pattern":"^0x[0-9a-fA-F]{40}$"}`)
	wholeNumberSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/wholeNumberSchema","anyOf":[{"type":"string","pattern":"^\\d+$"},{"type":"integer"}]}`)
	hexSchemaLoader         = gojsonschema.NewStringLoader(`{"id":"/hexSchema","type":"string","pattern":"^0x(([0-9a-fA-F])*)?$"}`)
	orderSchemaLoader       = gojsonschema.NewStringLoader(`{"id":"/orderSchema","properties":{"makerAddress":{"$ref":"/addressSchema"},"takerAddress":{"$ref":"/addressSchema"},"makerFee":{"$ref":"/wholeNumberSchema"},"takerFee":{"$ref":"/wholeNumberSchema"},"senderAddress":{"$ref":"/addressSchema"},"makerAssetAmount":{"$ref":"/wholeNumberSchema"},"takerAssetAmount":{"$ref":"/wholeNumberSchema"},"makerAssetData":{"$ref":"/hexSchema"},"takerAssetData":{"$ref":"/hexSchema"},"makerFeeAssetData":{"$ref":"/hexSchema"},"takerFeeAssetData":{"$ref":"/hexSchema"},"salt":{"$ref":"/wholeNumberSchema"},"feeRecipientAddress":{"$ref":"/addressSchema"},"expirationTimeSeconds":{"$ref":"/wholeNumberSchema"},"exchangeAddress":{"$ref":"/addressSchema"},"chainId": {"type": "number"}},"required":["makerAddress","takerAddress","makerFee","takerFee","senderAddress","makerAssetAmount","takerAssetAmount","makerAssetData","takerAssetData","makerFeeAssetData","takerFeeAssetData","salt","feeRecipientAddress","expirationTimeSeconds","exchangeAddress","chainId"],"type":"object"}`)
	signedOrderSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/signedOrderSchema","allOf":[{"$ref":"/orderSchema"},{"properties":{"signature":{"$ref":"/hexSchema"}},"required":["signature"]}]}`)
	meshMessageSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/meshMessageSchema","properties":{"MessageType":{"type":"string"},"Order":{"$ref":"/signedOrderSchema"}},"required":["MessageType","Order"]}`)
)

func setupMeshMessageSchemaValidator() (*gojsonschema.Schema, error) {
	sl := gojsonschema.NewSchemaLoader()
	if err := sl.AddSchemas(addressSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(wholeNumberSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(hexSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(orderSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(signedOrderSchemaLoader); err != nil {
		return nil, err
	}
	schema, err := sl.Compile(meshMessageSchemaLoader)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func setupOrderSchemaValidator() (*gojsonschema.Schema, error) {
	sl := gojsonschema.NewSchemaLoader()
	if err := sl.AddSchemas(addressSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(wholeNumberSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(hexSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(orderSchemaLoader); err != nil {
		return nil, err
	}
	schema, err := sl.Compile(signedOrderSchemaLoader)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func (app *App) schemaValidateOrder(o []byte) (*gojsonschema.Result, error) {
	orderLoader := gojsonschema.NewBytesLoader(o)

	result, err := app.orderJSONSchema.Validate(orderLoader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (app *App) schemaValidateMeshMessage(o []byte) (*gojsonschema.Result, error) {
	messageLoader := gojsonschema.NewBytesLoader(o)

	result, err := app.meshMessageJSONSchema.Validate(messageLoader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func validateMessageSize(message *p2p.Message) error {
	if len(message.Data) > constants.MaxOrderSizeInBytes {
		return constants.ErrMaxMessageSize
	}
	return nil
}
