package orderfilter

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/0xProject/0x-mesh/ethereum"
	canonicaljson "github.com/gibson042/canonicaljson-go"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

const (
	pubsubTopicVersion          = 3
	topicVersionFormat          = "/0x-orders/version/%d%s"
	topicChainIDAndSchemaFormat = "/chain/%d/schema/%s"
	fullTopicFormat             = "/0x-orders/version/%d/chain/%d/schema/%s"
	rendezvousVersion           = 1
	fullRendezvousFormat        = "/0x-custom-filter-rendezvous/version/%d/chain/%d/schema/%s"
)

type WrongTopicVersionError struct {
	expectedVersion int
	actualVersion   int
}

func (e WrongTopicVersionError) Error() string {
	return fmt.Sprintf("wrong topic version: expected %d but got %d", e.expectedVersion, e.actualVersion)
}

var (
	// Built-in schemas
	addressSchemaLoader     = jsonschema.NewStringLoader(`{"id":"/address","type":"string","pattern":"^0x[0-9a-fA-F]{40}$"}`)
	wholeNumberSchemaLoader = jsonschema.NewStringLoader(`{"id":"/wholeNumber","anyOf":[{"type":"string","pattern":"^\\d+$"},{"type":"integer"}]}`)
	hexSchemaLoader         = jsonschema.NewStringLoader(`{"id":"/hex","type":"string","pattern":"^0x(([0-9a-fA-F][0-9a-fA-F])+)?$"}`)
	orderSchemaLoader       = jsonschema.NewStringLoader(`{"id":"/order","properties":{"makerAddress":{"$ref":"/address"},"takerAddress":{"$ref":"/address"},"makerFee":{"$ref":"/wholeNumber"},"takerFee":{"$ref":"/wholeNumber"},"senderAddress":{"$ref":"/address"},"makerAssetAmount":{"$ref":"/wholeNumber"},"takerAssetAmount":{"$ref":"/wholeNumber"},"makerAssetData":{"$ref":"/hex"},"takerAssetData":{"$ref":"/hex"},"makerFeeAssetData":{"$ref":"/hex"},"takerFeeAssetData":{"$ref":"/hex"},"salt":{"$ref":"/wholeNumber"},"feeRecipientAddress":{"$ref":"/address"},"expirationTimeSeconds":{"$ref":"/wholeNumber"},"exchangeAddress":{"$ref":"/exchangeAddress"},"chainId":{"$ref":"/chainId"}},"required":["makerAddress","takerAddress","makerFee","takerFee","senderAddress","makerAssetAmount","takerAssetAmount","makerAssetData","takerAssetData","makerFeeAssetData","takerFeeAssetData","salt","feeRecipientAddress","expirationTimeSeconds","exchangeAddress","chainId"],"type":"object"}`)
	signedOrderSchemaLoader = jsonschema.NewStringLoader(`{"id":"/signedOrder","allOf":[{"$ref":"/order"},{"properties":{"signature":{"$ref":"/hex"}},"required":["signature"]}]}`)

	// Root schemas
	rootOrderSchemaLoader        = jsonschema.NewStringLoader(`{"id":"/rootOrder","allOf":[{"$ref":"/customOrder"},{"$ref":"/signedOrder"}]}`)
	rootOrderMessageSchemaLoader = jsonschema.NewStringLoader(`{"id":"/rootOrderMessage","properties":{"messageType":{"type":"string","pattern":"order"},"order":{"$ref":"/rootOrder"},"topics":{"type":"array","minItems":1,"items":{"type":"string"}}},"required":["messageType","order","topics"]}`)
)

const (
	// DefaultCustomOrderSchema is the default schema for /customOrder. It
	// includes all 0x orders and doesn't add any additional requirements.
	DefaultCustomOrderSchema = `{}`
)

var builtInSchemas = []jsonschema.JSONLoader{
	addressSchemaLoader,
	wholeNumberSchemaLoader,
	hexSchemaLoader,
	orderSchemaLoader,
	signedOrderSchemaLoader,
}

func GetDefaultFilter(chainID int, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	return New(chainID, DefaultCustomOrderSchema, contractAddresses)
}

func GetDefaultTopic(chainID int, contractAddresses ethereum.ContractAddresses) (string, error) {
	defaultFilter, err := GetDefaultFilter(chainID, contractAddresses)
	if err != nil {
		return "", err
	}
	return defaultFilter.Topic(), nil
}

type Filter struct {
	validatorLoaded      bool
	encodedSchema        string
	version              int
	chainID              int
	rawCustomOrderSchema string
	orderSchema          *jsonschema.Schema
	messageSchema        *jsonschema.Schema
}

func New(chainID int, customOrderSchema string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	orderLoader, err := newLoader(chainID, customOrderSchema, contractAddresses)
	if err != nil {
		return nil, err
	}
	rootOrderSchema, err := orderLoader.Compile(rootOrderSchemaLoader)
	if err != nil {
		return nil, err
	}

	messageLoader, err := newLoader(chainID, customOrderSchema, contractAddresses)
	if err := messageLoader.AddSchemas(rootOrderSchemaLoader); err != nil {
		return nil, err
	}
	rootOrderMessageSchema, err := messageLoader.Compile(rootOrderMessageSchemaLoader)
	if err != nil {
		return nil, err
	}
	return &Filter{
		chainID:              chainID,
		rawCustomOrderSchema: customOrderSchema,
		orderSchema:          rootOrderSchema,
		messageSchema:        rootOrderMessageSchema,
	}, nil
}

func loadExchangeAddress(loader *jsonschema.SchemaLoader, chainID int, contractAddresses ethereum.ContractAddresses) error {
	// Note that exchangeAddressSchema accepts both checksummed and
	// non-checksummed (i.e. all lowercase) addresses.
	exchangeAddressSchema := fmt.Sprintf(`{"enum":[%q,%q]}`, contractAddresses.Exchange.Hex(), strings.ToLower(contractAddresses.Exchange.Hex()))
	return loader.AddSchema("/exchangeAddress", jsonschema.NewStringLoader(exchangeAddressSchema))
}

func loadChainID(loader *jsonschema.SchemaLoader, chainID int) error {
	chainIDSchema := fmt.Sprintf(`{"const":%d}`, chainID)
	return loader.AddSchema("/chainId", jsonschema.NewStringLoader(chainIDSchema))
}

func newLoader(chainID int, customOrderSchema string, contractAddresses ethereum.ContractAddresses) (*jsonschema.SchemaLoader, error) {
	loader := jsonschema.NewSchemaLoader()
	if err := loadChainID(loader, chainID); err != nil {
		return nil, err
	}
	if err := loadExchangeAddress(loader, chainID, contractAddresses); err != nil {
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

func NewFromTopic(topic string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	// TODO(albrow): Use a cache for topic -> filter
	var version int
	var chainIDAndSchema string
	if _, err := fmt.Sscanf(topic, topicVersionFormat, &version, &chainIDAndSchema); err != nil {
		return nil, fmt.Errorf("could not parse topic version for topic: %q", topic)
	}
	if version != pubsubTopicVersion {
		return nil, WrongTopicVersionError{
			expectedVersion: pubsubTopicVersion,
			actualVersion:   version,
		}
	}
	var chainID int
	var base64EncodedSchema string
	if _, err := fmt.Sscanf(chainIDAndSchema, topicChainIDAndSchemaFormat, &chainID, &base64EncodedSchema); err != nil {
		return nil, fmt.Errorf("could not parse chainID and schema from topic: %q", topic)
	}
	customOrderSchema, err := base64.URLEncoding.DecodeString(base64EncodedSchema)
	if err != nil {
		return nil, fmt.Errorf("could not base64-decode order schema: %q", base64EncodedSchema)
	}
	return New(chainID, string(customOrderSchema), contractAddresses)
}

func (f *Filter) Topic() string {
	if f.encodedSchema == "" {
		f.encodedSchema = f.generateEncodedSchema()
	}
	return fmt.Sprintf(fullTopicFormat, pubsubTopicVersion, f.chainID, f.encodedSchema)
}

func (f *Filter) Rendezvous() string {
	if f.encodedSchema == "" {
		f.encodedSchema = f.generateEncodedSchema()
	}
	return fmt.Sprintf(fullRendezvousFormat, rendezvousVersion, f.chainID, f.encodedSchema)
}

func (f *Filter) generateEncodedSchema() string {
	// Note(albrow): We use canonicaljson to elminate any differences in spacing,
	// formatting, and the order of field names. This ensures that two filters
	// that are semantically the same JSON object always encode to exactly the
	// same canonical topic string.
	//
	// So for example:
	//
	//     {
	//         "foo": "bar",
	//         "biz": "baz"
	//     }
	//
	// Will encode to the same topic string as:
	//
	//     {
	//         "biz":"baz",
	//         "foo":"bar"
	//     }
	//
	var holder interface{} = struct{}{}
	_ = canonicaljson.Unmarshal([]byte(f.rawCustomOrderSchema), &holder)
	canonicalOrderSchemaJSON, _ := canonicaljson.Marshal(holder)
	return base64.URLEncoding.EncodeToString(canonicalOrderSchemaJSON)
}
