package orderfilter

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

func GetDefaultFilter(chainID int, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	return New(chainID, DefaultCustomOrderSchema, contractAddresses)
}

// FIXME(jalextowle), is this doing the right thing?
func GetDefaultTopics(chainID int, contractAddresses ethereum.ContractAddresses) ([]string, error) {
	defaultFilter, err := GetDefaultFilter(chainID, contractAddresses)
	if err != nil {
		return []string{""}, err
	}
	return defaultFilter.Topics(), nil
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

func NewFromTopic(topics []string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	// TODO(albrow): Use a cache for topic -> filter
	var version int
	var chainIDAndSchema string
	// FIXME(jalextwole): This will need to be able to support multiple different order
	// versions
	if _, err := fmt.Sscanf(topics[0], topicVersionFormat, &version, &chainIDAndSchema); err != nil {
		return nil, fmt.Errorf("could not parse topic version for topic: %q", topics[0])
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
		return nil, fmt.Errorf("could not parse chainID and schema from topic: %q", topics[0])
	}
	customOrderSchema, err := base64.URLEncoding.DecodeString(base64EncodedSchema)
	if err != nil {
		return nil, fmt.Errorf("could not base64-decode order schema: %q", base64EncodedSchema)
	}
	return New(chainID, string(customOrderSchema), contractAddresses)
}

func (f *Filter) Rendezvous() string {
	if f.encodedSchema == "" {
		f.encodedSchema = f.generateEncodedSchema()
	}
	return fmt.Sprintf(fullRendezvousFormat, rendezvousVersion, f.chainID, f.encodedSchema)
}

// FIXME(jalextowle): We'll need to update the orderfilter implementation to accomodate
// v4 orders.
//
func (f *Filter) Topics() []string {
	if f.encodedSchema == "" {
		f.encodedSchema = f.generateEncodedSchema()
	}
	// FIXME(jalextowle): Add v4 topics when they are ready
	return []string{fmt.Sprintf(fullTopicFormat, pubsubTopicVersion, f.chainID, f.encodedSchema)}
}

// Dummy declaration to ensure that ValidatePubSubMessage matches the expected
// signature for pubsub.Validator.
var _ pubsub.Validator = (&Filter{}).ValidatePubSubMessage

// FIXME(jalextowle): We'll need to update the orderfilter implementation to accomodate
// v4 orders.
//
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

func (f *Filter) generateEncodedSchema() string {
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
	_ = canonicaljson.Unmarshal([]byte(f.rawCustomOrderSchema), &holder)
	canonicalOrderSchemaJSON, _ := canonicaljson.Marshal(holder)
	return base64.URLEncoding.EncodeToString(canonicalOrderSchemaJSON)
}

// NOTE(jalextowle): Due to the discrepancy between orderfilters used in browser
// nodes and those used in standalone nodes, we cannot simply encode orderfilter.Filter.
// Instead, we marshal a minimal representation of the filter, and then we recreate
// the filter in the node that unmarshals the filter. This ensures that any node
// that unmarshals the orderfilter will be capable of using the filter.
type jsonMarshallerForFilter struct {
	CustomOrderSchema string         `json:"customOrderSchema"`
	ChainID           int            `json:"chainID"`
	ExchangeAddress   common.Address `json:"exchangeAddress"`
}

func (f *Filter) MarshalJSON() ([]byte, error) {
	j := jsonMarshallerForFilter{
		CustomOrderSchema: f.rawCustomOrderSchema,
		ChainID:           f.chainID,
		ExchangeAddress:   f.exchangeAddress,
	}
	return json.Marshal(j)
}

func (f *Filter) UnmarshalJSON(data []byte) error {
	j := jsonMarshallerForFilter{}
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}
	filter, err := New(j.ChainID, j.CustomOrderSchema, ethereum.ContractAddresses{Exchange: j.ExchangeAddress})
	if err != nil {
		return err
	}
	*f = *filter
	return nil
}
