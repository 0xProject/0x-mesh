package orderfilter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/xeipuuv/gojsonschema"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

const (
	pubsubTopicVersion = 3
)

var (
	// Built-in schemas
	addressSchemaLoader     = jsonschema.NewStringLoader(`{"id":"/address","type":"string","pattern":"^0x[0-9a-fA-F]{40}$"}`)
	wholeNumberSchemaLoader = jsonschema.NewStringLoader(`{"id":"/wholeNumber","anyOf":[{"type":"string","pattern":"^\\d+$"},{"type":"integer"}]}`)
	hexSchemaLoader         = jsonschema.NewStringLoader(`{"id":"/hex","type":"string","pattern":"^0x(([0-9a-fA-F][0-9a-fA-F])+)?$"}`)
	orderSchemaLoader       = jsonschema.NewStringLoader(`{"id":"/order","properties":{"makerAddress":{"$ref":"/address"},"takerAddress":{"$ref":"/address"},"makerFee":{"$ref":"/wholeNumber"},"takerFee":{"$ref":"/wholeNumber"},"senderAddress":{"$ref":"/address"},"makerAssetAmount":{"$ref":"/wholeNumber"},"takerAssetAmount":{"$ref":"/wholeNumber"},"makerAssetData":{"$ref":"/hex"},"takerAssetData":{"$ref":"/hex"},"salt":{"$ref":"/wholeNumber"},"exchangeAddress":{"$ref":"/exchangeAddress"},"feeRecipientAddress":{"$ref":"/address"},"expirationTimeSeconds":{"$ref":"/wholeNumber"}},"required":["makerAddress","takerAddress","makerFee","takerFee","senderAddress","makerAssetAmount","takerAssetAmount","makerAssetData","takerAssetData","salt","exchangeAddress","feeRecipientAddress","expirationTimeSeconds"],"type":"object"}`)
	signedOrderSchemaLoader = jsonschema.NewStringLoader(`{"id":"/signedOrder","allOf":[{"$ref":"/order"},{"properties":{"signature":{"$ref":"/hex"}},"required":["signature"]}]}`)

	// Root schemas
	rootOrderSchemaLoader = jsonschema.NewStringLoader(`{"id":"/rootOrder","allOf":[{"$ref":"/customOrder"},{"$ref":"/signedOrder"}]}`)
	// TODO(albrow): Add Topics as a required field for messages.
	rootMessageSchemaLoader = jsonschema.NewStringLoader(`{"id":"/rootMessage","properties":{"MessageType":{"type":"string"},"Order":{"$ref":"/rootOrder"}},"required":["MessageType","Order"]}`)

	// Default schema for /customOrder
	DefaultCustomOrderSchema = `{}`
)

var builtInSchemas = []jsonschema.JSONLoader{
	addressSchemaLoader,
	wholeNumberSchemaLoader,
	hexSchemaLoader,
	orderSchemaLoader,
	signedOrderSchemaLoader,
}

type Filter struct {
	version              int
	chainID              int
	rawCustomOrderSchema string
	orderSchema          *jsonschema.Schema
	messageSchema        *jsonschema.Schema
}

func New(chainID int, customOrderSchema string) (*Filter, error) {
	orderLoader, err := newLoader(chainID, customOrderSchema)
	rootOrderSchema, err := orderLoader.Compile(rootOrderSchemaLoader)
	if err != nil {
		return nil, err
	}

	messageLoader, err := newLoader(chainID, customOrderSchema)
	if err := messageLoader.AddSchemas(rootOrderSchemaLoader); err != nil {
		return nil, err
	}
	rootMessageSchema, err := messageLoader.Compile(rootMessageSchemaLoader)
	if err != nil {
		return nil, err
	}
	return &Filter{
		chainID:              chainID,
		rawCustomOrderSchema: customOrderSchema,
		orderSchema:          rootOrderSchema,
		messageSchema:        rootMessageSchema,
	}, nil
}

func loadExchangeAddress(loader *jsonschema.SchemaLoader, chainID int) error {
	contractAddresses, err := ethereum.GetContractAddressesForChainID(chainID)
	if err != nil {
		return err
	}
	// Note that exchangeAddressSchema accepts both checksummed and
	// non-checksummed (i.e. all lowercase) addresses.
	exchangeAddressSchema := fmt.Sprintf(`{"oneOf":[{"type":"string","pattern":%q},{"type":"string","pattern":%q}]}`, contractAddresses.Exchange.Hex(), strings.ToLower(contractAddresses.Exchange.Hex()))
	return loader.AddSchema("/exchangeAddress", jsonschema.NewStringLoader(exchangeAddressSchema))
}

func newLoader(chainID int, customOrderSchema string) (*jsonschema.SchemaLoader, error) {
	loader := jsonschema.NewSchemaLoader()
	if err := loadExchangeAddress(loader, chainID); err != nil {
		return nil, err
	}
	if err := loader.AddSchemas(builtInSchemas...); err != nil {
		return nil, err
	}
	if err := loader.AddSchema("/customOrder", jsonschema.NewStringLoader(customOrderSchema)); err != nil {
		return nil, err
	}
	return loader, nil
}

func NewFromTopic(topic string) (*Filter, error) {
	// TODO(albrow): Create a new Filter based on the topic.
	// TODO(albrow): Use a cache for topic -> filter
	return nil, errors.New("not yet implemented")
}

func (v *Filter) Topic() string {
	base64OrderSchema := base64.URLEncoding.EncodeToString([]byte(v.rawCustomOrderSchema))
	return fmt.Sprintf("/0x-orders/version/%d/chain/%d/schema/%s", pubsubTopicVersion, v.chainID, base64OrderSchema)
}

func (v *Filter) MatchMessageJSON(messageJSON []byte) (bool, error) {
	result, err := v.messageSchema.Validate(jsonschema.NewBytesLoader(messageJSON))
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (v *Filter) ValidateOrderJSON(orderJSON []byte) (*jsonschema.Result, error) {
	return v.orderSchema.Validate(gojsonschema.NewBytesLoader(orderJSON))
}

func (v *Filter) ValidateOrder(order *zeroex.SignedOrder) (*jsonschema.Result, error) {
	return v.orderSchema.Validate(gojsonschema.NewGoLoader(order))
}
