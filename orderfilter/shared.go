package orderfilter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	canonicaljson "github.com/gibson042/canonicaljson-go"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

const (
	pubsubTopicVersionV3        = 3
	pubsubTopicVersionV4        = 4
	topicVersionFormat          = "/0x-orders/version/%d%s"
	topicChainIDAndSchemaFormat = "/chain/%d/schema/%s"
	fullTopicFormat             = "/0x-orders/version/%d/chain/%d/schema/%s"
	rendezvousVersionV3         = 1
	rendezvousVersionV4         = 2
	fullRendezvousFormat        = "/0x-custom-filter-rendezvous/version/%d/chain/%d/schema/%s"
)

type WrongTopicVersionError struct {
	expectedVersion int
	actualVersion   int
}

func (e WrongTopicVersionError) Error() string {
	return fmt.Sprintf("wrong topic version: expected %d but got %d", e.expectedVersion, e.actualVersion)
}

func GetDefaultFilter(chainID int, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	return New(chainID, DefaultCustomOrderSchema, DefaultCustomOrderSchema, contractAddresses)
}

func GetDefaultTopicV3(chainID int, contractAddresses ethereum.ContractAddresses) (string, error) {
	defaultFilter, err := GetDefaultFilter(chainID, contractAddresses)
	if err != nil {
		return "", err
	}
	return defaultFilter.TopicV3(), nil
}

func GetDefaultTopicV4(chainID int, contractAddresses ethereum.ContractAddresses) (string, error) {
	defaultFilter, err := GetDefaultFilter(chainID, contractAddresses)
	if err != nil {
		return "", err
	}
	return defaultFilter.TopicV4(), nil
}

// MatchOrder returns true if the order passes the filter. It only returns an
// error if there was a problem with validation. For details about
// orders that do not pass the filter, use ValidateOrder.
func (f *Filter) MatchOrder(order *zeroex.SignedOrder) (bool, error) {
	switch order.Order.(type) {
	case *zeroex.OrderV3:
		result, err := f.ValidateOrderV3(order)
		if err != nil {
			return false, err
		}
		return result.Valid(), nil
	case *zeroex.OrderV4:
		result, err := f.ValidateOrderV4(order)
		if err != nil {
			return false, err
		}
		return result.Valid(), nil
	default:
		return false, errors.New("Can't match unrecognized order type")
	}
}

func NewFromTopics(topicV3 string, topicV4 string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	// TODO(albrow): Use a cache for topic -> filter
	var version int
	var chainIDAndSchema string
	if _, err := fmt.Sscanf(topicV3, topicVersionFormat, &version, &chainIDAndSchema); err != nil {
		return nil, fmt.Errorf("could not parse topic version for topic: %q", topicV3)
	}
	if version != pubsubTopicVersionV3 {
		return nil, WrongTopicVersionError{
			expectedVersion: pubsubTopicVersionV3,
			actualVersion:   version,
		}
	}
	var chainID int
	var base64EncodedSchema string
	if _, err := fmt.Sscanf(chainIDAndSchema, topicChainIDAndSchemaFormat, &chainID, &base64EncodedSchema); err != nil {
		return nil, fmt.Errorf("could not parse chainID and schema from topic: %q", topicV3)
	}
	customOrderSchemaV3, err := base64.URLEncoding.DecodeString(base64EncodedSchema)
	if err != nil {
		return nil, fmt.Errorf("could not base64-decode order schema: %q", base64EncodedSchema)
	}
	if _, err := fmt.Sscanf(topicV4, topicVersionFormat, &version, &chainIDAndSchema); err != nil {
		return nil, fmt.Errorf("could not parse topic version for topic: %q", topicV4)
	}
	if version != pubsubTopicVersionV4 {
		return nil, WrongTopicVersionError{
			expectedVersion: pubsubTopicVersionV4,
			actualVersion:   version,
		}
	}
	if _, err := fmt.Sscanf(chainIDAndSchema, topicChainIDAndSchemaFormat, &chainID, &base64EncodedSchema); err != nil {
		return nil, fmt.Errorf("could not parse chainID and schema from topic: %q", topicV4)
	}
	customOrderSchemaV4, err := base64.URLEncoding.DecodeString(base64EncodedSchema)
	if err != nil {
		return nil, fmt.Errorf("could not base64-decode order schema: %q", base64EncodedSchema)
	}

	return New(chainID, string(customOrderSchemaV3), string(customOrderSchemaV4), contractAddresses)
}

func (f *Filter) RendezvousV3() string {
	if f.encodedSchemaV3 == "" {
		f.encodedSchemaV3 = f.generateEncodedSchemaV3()
	}
	return fmt.Sprintf(fullRendezvousFormat, rendezvousVersionV3, f.chainID, f.encodedSchemaV3)
}

func (f *Filter) RendezvousV4() string {
	if f.encodedSchemaV4 == "" {
		f.encodedSchemaV4 = f.generateEncodedSchemaV3()
	}
	return fmt.Sprintf(fullRendezvousFormat, rendezvousVersionV4, f.chainID, f.encodedSchemaV4)
}

func (f *Filter) TopicV3() string {
	if f.encodedSchemaV3 == "" {
		f.encodedSchemaV3 = f.generateEncodedSchemaV3()
	}
	return fmt.Sprintf(fullTopicFormat, pubsubTopicVersionV3, f.chainID, f.encodedSchemaV3)
}

func (f *Filter) TopicV4() string {
	if f.encodedSchemaV4 == "" {
		f.encodedSchemaV4 = f.generateEncodedSchemaV4()
	}
	return fmt.Sprintf(fullTopicFormat, pubsubTopicVersionV4, f.chainID, f.encodedSchemaV4)
}

// Dummy declarations to ensure that ValidatePubSubMessageV3 and match the expected
// signature for pubsub.Validator.
var _ pubsub.Validator = (&Filter{}).ValidatePubSubMessageV3
var _ pubsub.Validator = (&Filter{}).ValidatePubSubMessageV4

// ValidatePubSubMessageV3 is an implementation of pubsub.Validator and will
// return true if the contents of the message pass the message JSON Schema.
func (f *Filter) ValidatePubSubMessageV3(ctx context.Context, sender peer.ID, msg *pubsub.Message) bool {
	isValid, err := f.MatchOrderMessageV3JSON(msg.Data)
	if err != nil {
		log.WithError(err).Error("MatchOrderMessageV3JSON returned an error")
		return false
	}
	return isValid
}

// ValidatePubSubMessageV3 is an implementation of pubsub.Validator and will
// return true if the contents of the message pass the message JSON Schema.
func (f *Filter) ValidatePubSubMessageV4(ctx context.Context, sender peer.ID, msg *pubsub.Message) bool {
	isValid, err := f.MatchOrderMessageV4JSON(msg.Data)
	if err != nil {
		log.WithError(err).Error("MatchOrderMessageV4JSON returned an error")
		return false
	}
	return isValid
}

func (f *Filter) generateEncodedSchemaV3() string {
	// Note(albrow): We use canonicaljson to eliminate any differences in spacing,
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
	_ = canonicaljson.Unmarshal([]byte(f.rawCustomOrderSchemaV3), &holder)
	canonicalOrderSchemaV3JSON, _ := canonicaljson.Marshal(holder)
	return base64.URLEncoding.EncodeToString(canonicalOrderSchemaV3JSON)
}

func (f *Filter) generateEncodedSchemaV4() string {
	// Note(albrow): We use canonicaljson to eliminate any differences in spacing,
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
	_ = canonicaljson.Unmarshal([]byte(f.rawCustomOrderSchemaV3), &holder)
	canonicalOrderSchemaV3JSON, _ := canonicaljson.Marshal(holder)
	return base64.URLEncoding.EncodeToString(canonicalOrderSchemaV3JSON)
}

// NOTE(jalextowle): Due to the discrepancy between orderfilters used in browser
// nodes and those used in standalone nodes, we cannot simply encode orderfilter.Filter.
// Instead, we marshal a minimal representation of the filter, and then we recreate
// the filter in the node that unmarshals the filter. This ensures that any node
// that unmarshals the orderfilter will be capable of using the filter.
type jsonMarshallerForFilter struct {
	// NOTE(jalextowle): As of right now, we must keep the json field name of
	// "customOrderSchema" instead of switching to "customOrderSchemaV3" for
	// backwards compatability
	CustomOrderSchemaV3 string `json:"customOrderSchema"`
	CustomOrderSchemaV4 string `json:"customOrderSchemaV4"`
	ChainID             int    `json:"chainID"`
	// NOTE(jalextowle): As of right now, we must keep the json field name of
	// "exchangeAddress" instead of switching to "exchangeV3" for backwards
	// compatability
	ExchangeV3 common.Address `json:"exchangeAddress"`
	ExchangeV4 common.Address `json:"exchangeV4"`
}

func (f *Filter) MarshalJSON() ([]byte, error) {
	j := jsonMarshallerForFilter{
		CustomOrderSchemaV3: f.rawCustomOrderSchemaV3,
		CustomOrderSchemaV4: f.rawCustomOrderSchemaV4,
		ChainID:             f.chainID,
		ExchangeV3:          f.exchangeAddressV3,
		ExchangeV4:          f.exchangeAddressV4,
	}
	return json.Marshal(j)
}

func (f *Filter) UnmarshalJSON(data []byte) error {
	j := jsonMarshallerForFilter{}
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}
	filter, err := New(j.ChainID, j.CustomOrderSchemaV3, j.CustomOrderSchemaV4, ethereum.ContractAddresses{ExchangeV3: j.ExchangeV3, ExchangeV4: j.ExchangeV4})
	if err != nil {
		return err
	}
	*f = *filter
	return nil
}
