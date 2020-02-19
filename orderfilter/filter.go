package orderfilter

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	canonicaljson "github.com/gibson042/canonicaljson-go"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

const (
	pubsubTopicVersion          = 3
	topicVersionFormat          = "/0x-orders/version/%d%s"
	topicChainIDAndSchemaFormat = "/chain/%d/schema/%s"
	fullTopicFormat             = "/0x-orders/version/%d/chain/%d/schema/%s"
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

func GetDefaultFilter(chainID int) (*Filter, error) {
	return New(chainID, DefaultCustomOrderSchema)
}

func GetDefaultTopic(chainID int) (string, error) {
	defaultFilter, err := GetDefaultFilter(chainID)
	if err != nil {
		return "", err
	}
	return defaultFilter.Topic(), nil
}

type Filter struct {
	topic                string
	version              int
	chainID              int
	rawCustomOrderSchema string
	orderSchema          *jsonschema.Schema
	messageSchema        *jsonschema.Schema
}

func New(chainID int, customOrderSchema string) (*Filter, error) {
	orderLoader, err := newLoader(chainID, customOrderSchema)
	if err != nil {
		return nil, err
	}
	rootOrderSchema, err := orderLoader.Compile(rootOrderSchemaLoader)
	if err != nil {
		return nil, err
	}

	messageLoader, err := newLoader(chainID, customOrderSchema)
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

func loadExchangeAddress(loader *jsonschema.SchemaLoader, chainID int) error {
	contractAddresses, err := ethereum.GetContractAddressesForChainID(chainID)
	if err != nil {
		return err
	}
	// Note that exchangeAddressSchema accepts both checksummed and
	// non-checksummed (i.e. all lowercase) addresses.
	exchangeAddressSchema := fmt.Sprintf(`{"enum":[%q,%q]}`, contractAddresses.Exchange.Hex(), strings.ToLower(contractAddresses.Exchange.Hex()))
	return loader.AddSchema("/exchangeAddress", jsonschema.NewStringLoader(exchangeAddressSchema))
}

func loadChainID(loader *jsonschema.SchemaLoader, chainID int) error {
	chainIDSchema := fmt.Sprintf(`{"const":%d}`, chainID)
	return loader.AddSchema("/chainId", jsonschema.NewStringLoader(chainIDSchema))
}

func newLoader(chainID int, customOrderSchema string) (*jsonschema.SchemaLoader, error) {
	loader := jsonschema.NewSchemaLoader()
	if err := loadChainID(loader, chainID); err != nil {
		return nil, err
	}
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
	return New(chainID, string(customOrderSchema))
}

func (f *Filter) Topic() string {
	if f.topic == "" {
		f.topic = f.generateTopic()
	}
	return f.topic
}

func (v *Filter) generateTopic() string {
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
	_ = canonicaljson.Unmarshal([]byte(v.rawCustomOrderSchema), &holder)
	canonicalOrderSchemaJSON, _ := canonicaljson.Marshal(holder)
	base64EncodedSchema := base64.URLEncoding.EncodeToString(canonicalOrderSchemaJSON)
	return fmt.Sprintf(fullTopicFormat, pubsubTopicVersion, v.chainID, base64EncodedSchema)
}

// MatchOrder returns true if the order passes the filter. It only returns an
// error if there was a problem with validation. For details about
// orders that do not pass the filter, use ValidateOrder.
func (f *Filter) MatchOrder(order *zeroex.SignedOrder) (bool, error) {
	result, err := f.ValidateOrder(order)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (f *Filter) MatchOrderMessageJSON(messageJSON []byte) (bool, error) {
	result, err := f.messageSchema.Validate(jsonschema.NewBytesLoader(messageJSON))
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (f *Filter) ValidateOrderJSON(orderJSON []byte) (*jsonschema.Result, error) {
	return f.orderSchema.Validate(jsonschema.NewBytesLoader(orderJSON))
}

func (f *Filter) ValidateOrder(order *zeroex.SignedOrder) (*jsonschema.Result, error) {
	return f.orderSchema.Validate(jsonschema.NewGoLoader(order))
}

// Dummy declaration to ensure that ValidatePubSubMessage matches the expected
// signature for pubsub.Validator.
var _ pubsub.Validator = (&Filter{}).ValidatePubSubMessage

// ValidatePubSubMessage is an implementation of pubsub.Validator and will
// return true if the contents of the message pass the message JSON Schema.
func (f *Filter) ValidatePubSubMessage(ctx context.Context, sender peer.ID, msg *pubsub.Message) bool {
	isValid, err := f.MatchOrderMessageJSON(msg.Data)
	if err != nil {
		log.WithError(err).Error("MatchOrderMessageJSON returned an error")
		return false
	}
	return isValid
}
