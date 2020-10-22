package zeroex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/signer"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	gethsigner "github.com/ethereum/go-ethereum/signer/core"
	"golang.org/x/crypto/sha3"
)

// Order represents an unsigned 0x v3 order
type OrderV3 struct {
	ChainID               *big.Int       `json:"chainId"`
	ExchangeAddress       common.Address `json:"exchangeAddress"`
	MakerAddress          common.Address `json:"makerAddress"`
	MakerAssetData        []byte         `json:"makerAssetData"`
	MakerFeeAssetData     []byte         `json:"makerFeeAssetData"`
	MakerAssetAmount      *big.Int       `json:"makerAssetAmount"`
	MakerFee              *big.Int       `json:"makerFee"`
	TakerAddress          common.Address `json:"takerAddress"`
	TakerAssetData        []byte         `json:"takerAssetData"`
	TakerFeeAssetData     []byte         `json:"takerFeeAssetData"`
	TakerAssetAmount      *big.Int       `json:"takerAssetAmount"`
	TakerFee              *big.Int       `json:"takerFee"`
	SenderAddress         common.Address `json:"senderAddress"`
	FeeRecipientAddress   common.Address `json:"feeRecipientAddress"`
	ExpirationTimeSeconds *big.Int       `json:"expirationTimeSeconds"`
	Salt                  *big.Int       `json:"salt"`

	// Cache hash for performance
	hash *common.Hash
}

type OrderV4 struct {
	ChainID          *big.Int       `json:"chainId"`
	Exchange         common.Address `json:"exchangeAddress"`
	Maker            common.Address `json:"maker"`
	MakerToken       common.Address `json:"makerToken"`
	MakerAssetAmount *big.Int       `json:"makerAssetAmount"`
	MakerFee         *big.Int       `json:"makerFee"`
	Taker            common.Address `json:"taker"`
	TakerToken       common.Address `json:"takerToken"`
	TakerAssetAmount *big.Int       `json:"takerAssetAmount"`
	TakerFee         *big.Int       `json:"takerFee"`
	Sender           common.Address `json:"sender"`
	FeeRecipient     common.Address `json:"feeRecipient"`
	Expiry           *big.Int       `json:"expiry"`
	Salt             *big.Int       `json:"salt"`
	Pool             *big.Int       `json:"pool"`
	Origin           *big.Int       `json:"origin"`

	// Cache hash for performance
	hash *common.Hash
}

// FIXME(jalextowle) Version should be a string
type SignedOrder struct {
	Version   int64
	Order     interface{}
	Signature []byte
}

type signedOrderJSON struct {
	Version   int64
	Order     json.RawMessage
	Signature []byte
}

// SignatureType represents the type of 0x signature encountered
type SignatureType uint8

// SignatureType values
const (
	IllegalSignature SignatureType = iota
	InvalidSignature
	EIP712Signature
	EthSignSignature
	WalletSignature
	ValidatorSignature
	PreSignedSignature
	EIP1271WalletSignature
	NSignatureTypesSignature
)

// OrderStatus represents the status of an order as returned from the 0x smart contracts
// as part of OrderInfo
type OrderStatus uint8

// OrderStatus values
const (
	OSInvalid OrderStatus = iota
	OSInvalidMakerAssetAmount
	OSInvalidTakerAssetAmount
	OSFillable
	OSExpired
	OSFullyFilled
	OSCancelled
	OSSignatureInvalid
	OSInvalidMakerAssetData
	OSInvalidTakerAssetData
)

// ContractEvent is an event emitted by a smart contract
type ContractEvent struct {
	BlockHash  common.Hash
	TxHash     common.Hash
	TxIndex    uint
	LogIndex   uint
	IsRemoved  bool
	Address    common.Address
	Kind       string
	Parameters interface{}
}

type contractEventJSON struct {
	BlockHash  common.Hash
	TxHash     common.Hash
	TxIndex    uint
	LogIndex   uint
	IsRemoved  bool
	Address    common.Address
	Kind       string
	Parameters json.RawMessage
}

// MarshalJSON implements a custom JSON marshaller for the ContractEvent type
func (c ContractEvent) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"blockHash":  c.BlockHash.Hex(),
		"txHash":     c.TxHash.Hex(),
		"txIndex":    c.TxIndex,
		"logIndex":   c.LogIndex,
		"isRemoved":  c.IsRemoved,
		"address":    c.Address,
		"kind":       c.Kind,
		"parameters": c.Parameters,
	}
	return json.Marshal(m)
}

// OrderEvent is the order event emitted by the Mesh GraphQL API or the SubScribeToOrderEvents
// method in core.
type OrderEvent struct {
	// Timestamp is an order event timestamp that can be used for bookkeeping purposes.
	// If the OrderEvent represents a Mesh-specific event (e.g., ADDED, STOPPED_WATCHING),
	// the timestamp is when the event was generated. If the event was generated after
	// re-validating an order at the latest block height (e.g., FILLED, UNFUNDED, CANCELED),
	// then it is set to the latest block timestamp at which the order was re-validated.
	Timestamp time.Time `json:"timestamp"`
	// OrderHash is the EIP712 hash of the 0x order
	OrderHash common.Hash `json:"orderHash"`
	// SignedOrder is the 0x signed order interface
	SignedOrder *SignedOrder `json:"signedOrder"`
	// EndState is the end state of this order at the time this event was generated
	EndState OrderEventEndState `json:"endState"`
	// FillableTakerAssetAmount is the amount for which this order is still fillable
	FillableTakerAssetAmount *big.Int `json:"fillableTakerAssetAmount"`
	// ContractEvents contains all the contract events that triggered this orders re-evaluation.
	// They did not all necessarily cause the orders state change itself, only it's re-evaluation.
	// Since it's state _did_ change, at least one of them did cause the actual state change.
	ContractEvents []*ContractEvent `json:"contractEvents"`
}

type orderEventJSON struct {
	Timestamp                time.Time            `json:"timestamp"`
	OrderHash                string               `json:"orderHash"`
	SignedOrder              *SignedOrder         `json:"signedOrder"`
	EndState                 string               `json:"endState"`
	FillableTakerAssetAmount string               `json:"fillableTakerAssetAmount"`
	ContractEvents           []*contractEventJSON `json:"contractEvents"`
}

// MarshalJSON implements a custom JSON marshaller for the OrderEvent type
func (o OrderEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"timestamp":                o.Timestamp,
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder,
		"endState":                 o.EndState,
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
		"contractEvents":           o.ContractEvents,
	})
}

// FIXME(jalextowle): This will need to implement custom json unmarshalling in a
// similar way to how contract events are unmarshalled.
func (o *OrderEvent) UnmarshalJSON(data []byte) error {
	var orderEventJSON orderEventJSON
	err := json.Unmarshal(data, &orderEventJSON)
	if err != nil {
		return err
	}
	return o.fromOrderEventJSON(orderEventJSON)
}

func (o *OrderEvent) fromOrderEventJSON(orderEventJSON orderEventJSON) error {
	o.Timestamp = orderEventJSON.Timestamp
	o.OrderHash = common.HexToHash(orderEventJSON.OrderHash)
	o.SignedOrder = orderEventJSON.SignedOrder
	o.EndState = OrderEventEndState(orderEventJSON.EndState)
	var ok bool
	o.FillableTakerAssetAmount, ok = math.ParseBig256(orderEventJSON.FillableTakerAssetAmount)
	if !ok {
		return errors.New("Invalid uint256 number encountered for FillableTakerAssetAmount")
	}
	o.ContractEvents = make([]*ContractEvent, len(orderEventJSON.ContractEvents))
	for i, eventJSON := range orderEventJSON.ContractEvents {
		contractEvent, err := unmarshalContractEvent(eventJSON)
		if err != nil {
			return err
		}
		o.ContractEvents[i] = contractEvent
	}
	return nil
}

func unmarshalContractEvent(eventJSON *contractEventJSON) (*ContractEvent, error) {
	event := &ContractEvent{
		BlockHash: eventJSON.BlockHash,
		TxHash:    eventJSON.TxHash,
		TxIndex:   eventJSON.TxIndex,
		LogIndex:  eventJSON.LogIndex,
		IsRemoved: eventJSON.IsRemoved,
		Address:   eventJSON.Address,
		Kind:      eventJSON.Kind,
	}

	switch eventJSON.Kind {
	case "ERC20TransferEvent":
		var parameters decoder.ERC20TransferEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ERC20ApprovalEvent":
		var parameters decoder.ERC20ApprovalEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ERC721TransferEvent":
		var parameters decoder.ERC721TransferEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ERC721ApprovalEvent":
		var parameters decoder.ERC721ApprovalEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ERC721ApprovalForAllEvent":
		var parameters decoder.ERC721ApprovalForAllEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ERC1155TransferSingleEvent":
		var parameters decoder.ERC1155TransferSingleEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ERC1155TransferBatchEvent":
		var parameters decoder.ERC1155TransferBatchEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ERC1155ApprovalForAllEvent":
		var parameters decoder.ERC1155ApprovalForAllEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "WethWithdrawalEvent":
		var parameters decoder.WethWithdrawalEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "WethDepositEvent":
		var parameters decoder.WethDepositEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ExchangeFillEvent":
		var parameters decoder.ExchangeFillEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ExchangeCancelEvent":
		var parameters decoder.ExchangeCancelEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	case "ExchangeCancelUpToEvent":
		var parameters decoder.ExchangeCancelUpToEvent
		if err := json.Unmarshal(eventJSON.Parameters, &parameters); err != nil {
			return nil, err
		}
		event.Parameters = parameters

	default:
		return nil, fmt.Errorf("unknown event kind: %s", eventJSON.Kind)
	}

	return event, nil
}

// OrderEventEndState enumerates all the possible order event types. An OrderEventEndState describes the
// end state of a 0x order after revalidation
type OrderEventEndState string

// OrderEventEndState values
const (
	// ESInvalid is an event that is never emitted. It is here to discern between a declared but uninitialized OrderEventEndState
	ESInvalid = OrderEventEndState("INVALID")
	// ESOrderAdded means an order was successfully added to the Mesh node
	ESOrderAdded = OrderEventEndState("ADDED")
	// ESOrderFilled means an order was filled for a partial amount
	ESOrderFilled = OrderEventEndState("FILLED")
	// ESOrderFullyFilled means an order was fully filled such that it's remaining fillableTakerAssetAmount is 0
	ESOrderFullyFilled = OrderEventEndState("FULLY_FILLED")
	// ESOrderCancelled means an order was cancelled on-chain
	ESOrderCancelled = OrderEventEndState("CANCELLED")
	// ESOrderExpired means an order expired according to the latest block timestamp
	ESOrderExpired = OrderEventEndState("EXPIRED")
	// ESOrderUnexpired means an order is no longer expired. This can happen if a block re-org causes the latest
	// block timestamp to decline below the order's expirationTimestamp (rare and usually short-lived)
	ESOrderUnexpired = OrderEventEndState("UNEXPIRED")
	// ESOrderBecameUnfunded means an order has become unfunded. This happens if the maker transfers the balance /
	// changes their allowance backing an order
	ESOrderBecameUnfunded = OrderEventEndState("UNFUNDED")
	// ESOrderFillabilityIncreased means the fillability of an order has increased. Fillability for an order can
	// increase if a previously processed fill event gets reverted, or if a maker tops up their balance/allowance
	// backing an order
	ESOrderFillabilityIncreased = OrderEventEndState("FILLABILITY_INCREASED")
	// ESStoppedWatching means an order is potentially still valid but was removed for a different reason (e.g.
	// the database is full or the peer that sent the order was misbehaving). The order will no longer be watched
	// and no further events for this order will be emitted. In some cases, the order may be re-added in the
	// future.
	ESStoppedWatching = OrderEventEndState("STOPPED_WATCHING")
)

var eip712OrderTypes = gethsigner.Types{
	"EIP712Domain": {
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "version",
			Type: "string",
		},
		{
			Name: "chainId",
			Type: "uint256",
		},
		{
			Name: "verifyingContract",
			Type: "address",
		},
	},
	"Order": {
		{
			Name: "makerAddress",
			Type: "address",
		},
		{
			Name: "takerAddress",
			Type: "address",
		},
		{
			Name: "feeRecipientAddress",
			Type: "address",
		},
		{
			Name: "senderAddress",
			Type: "address",
		},
		{
			Name: "makerAssetAmount",
			Type: "uint256",
		},
		{
			Name: "takerAssetAmount",
			Type: "uint256",
		},
		{
			Name: "makerFee",
			Type: "uint256",
		},
		{
			Name: "takerFee",
			Type: "uint256",
		},
		{
			Name: "expirationTimeSeconds",
			Type: "uint256",
		},
		{
			Name: "salt",
			Type: "uint256",
		},
		{
			Name: "makerAssetData",
			Type: "bytes",
		},
		{
			Name: "takerAssetData",
			Type: "bytes",
		},
		{
			Name: "makerFeeAssetData",
			Type: "bytes",
		},
		{
			Name: "takerFeeAssetData",
			Type: "bytes",
		},
	},
}

// ResetHash resets the cached order hash. Usually only required for testing.
func (o *OrderV3) ResetHash() {
	o.hash = nil
}

// ComputeOrderHash computes a 0x order hash
func (o *OrderV3) ComputeOrderHash() (common.Hash, error) {
	if o.hash != nil {
		return *o.hash, nil
	}

	chainID := math.NewHexOrDecimal256(o.ChainID.Int64())
	var domain = gethsigner.TypedDataDomain{
		Name:              "0x Protocol",
		Version:           "3.0.0",
		ChainId:           chainID,
		VerifyingContract: o.ExchangeAddress.Hex(),
	}

	var message = map[string]interface{}{
		"makerAddress":          o.MakerAddress.Hex(),
		"takerAddress":          o.TakerAddress.Hex(),
		"senderAddress":         o.SenderAddress.Hex(),
		"feeRecipientAddress":   o.FeeRecipientAddress.Hex(),
		"makerAssetData":        o.MakerAssetData,
		"makerFeeAssetData":     o.MakerFeeAssetData,
		"takerAssetData":        o.TakerAssetData,
		"takerFeeAssetData":     o.TakerFeeAssetData,
		"salt":                  o.Salt.String(),
		"makerFee":              o.MakerFee.String(),
		"takerFee":              o.TakerFee.String(),
		"makerAssetAmount":      o.MakerAssetAmount.String(),
		"takerAssetAmount":      o.TakerAssetAmount.String(),
		"expirationTimeSeconds": o.ExpirationTimeSeconds.String(),
	}

	var typedData = gethsigner.TypedData{
		Types:       eip712OrderTypes,
		PrimaryType: "Order",
		Domain:      domain,
		Message:     message,
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return common.Hash{}, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Hash{}, err
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hashBytes := keccak256(rawData)
	hash := common.BytesToHash(hashBytes)
	o.hash = &hash
	return hash, nil
}

// SignOrder signs the 0x order with the supplied Signer
func SignOrder(signer signer.Signer, order *OrderV3) (*SignedOrder, error) {
	if order == nil {
		return nil, errors.New("cannot sign nil order")
	}
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		return nil, err
	}

	ecSignature, err := signer.EthSign(orderHash.Bytes(), order.MakerAddress)
	if err != nil {
		return nil, err
	}

	// Generate 0x EthSign Signature (append the signature type byte)
	signature := make([]byte, 66)
	signature[0] = ecSignature.V
	copy(signature[1:33], ecSignature.R[:])
	copy(signature[33:65], ecSignature.S[:])
	signature[65] = byte(EthSignSignature)
	signedOrder := &SignedOrder{
		Order:     *order,
		Signature: signature,
	}
	return signedOrder, nil
}

// SignTestOrder signs the 0x order with the local test signer
func SignTestOrder(order *OrderV3) (*SignedOrder, error) {
	testSigner := signer.NewTestSigner()
	signedOrder, err := SignOrder(testSigner, order)
	if err != nil {
		return nil, err
	}
	return signedOrder, nil
}

// FIXME(jalextowle): We'll need to implement trim on both order types. For now,
// we can get away with throwing an error if we can't cast to OrderV3.
//
// Trim converts the order to a LibOrderOrder, which is the format expected by
// our smart contracts. It removes the ChainID and ExchangeAddress fields.
func (s *SignedOrder) Trim() wrappers.LibOrderOrder {
	switch o := s.Order.(type) {
	case *OrderV3:
		return wrappers.LibOrderOrder{
			MakerAddress:          o.MakerAddress,
			TakerAddress:          o.TakerAddress,
			FeeRecipientAddress:   o.FeeRecipientAddress,
			SenderAddress:         o.SenderAddress,
			MakerAssetAmount:      o.MakerAssetAmount,
			TakerAssetAmount:      o.TakerAssetAmount,
			MakerFee:              o.MakerFee,
			TakerFee:              o.TakerFee,
			ExpirationTimeSeconds: o.ExpirationTimeSeconds,
			Salt:                  o.Salt,
			MakerAssetData:        o.MakerAssetData,
			MakerFeeAssetData:     o.MakerFeeAssetData,
			TakerAssetData:        o.TakerAssetData,
			TakerFeeAssetData:     o.TakerFeeAssetData,
		}
	default:
		panic("Can't trim non-v3 order")
	}
}

// FIXME(jalextowle): What low-level messaging format should Mesh use?
//
// SignedOrderV3JSON is an unmodified JSON representation of a SignedOrder
type SignedOrderV3JSON struct {
	ChainID               int64  `json:"chainId"`
	ExchangeAddress       string `json:"exchangeAddress"`
	MakerAddress          string `json:"makerAddress"`
	MakerAssetData        string `json:"makerAssetData"`
	MakerFeeAssetData     string `json:"makerFeeAssetData"`
	MakerAssetAmount      string `json:"makerAssetAmount"`
	MakerFee              string `json:"makerFee"`
	TakerAddress          string `json:"takerAddress"`
	TakerAssetData        string `json:"takerAssetData"`
	TakerFeeAssetData     string `json:"takerFeeAssetData"`
	TakerAssetAmount      string `json:"takerAssetAmount"`
	TakerFee              string `json:"takerFee"`
	SenderAddress         string `json:"senderAddress"`
	FeeRecipientAddress   string `json:"feeRecipientAddress"`
	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
	Salt                  string `json:"salt"`
	Signature             string `json:"signature"`
}

// FIXME(jalextowle): Should this have a Maker field?
// The answer to that question really depends on how this struct is interacted with.
// If it's typically used when messaging with other peers, we probably need to come up with a new
// communication mechanism.
//
// SignedOrderJSON is an unmodified JSON representation of a SignedOrder
type SignedOrderV4JSON struct {
	Version          int64  `json:"version"`
	ChainID          int64  `json:"chainId"`
	Exchange         string `json:"exchange"`
	Maker            string `json:"maker"`
	MakerToken       string `json:"makerToken"`
	MakerAssetAmount string `json:"makerAssetAmount"`
	MakerFee         string `json:"makerFee"`
	Taker            string `json:"taker"`
	TakerToken       string `json:"takerToken"`
	TakerAssetAmount string `json:"takerAssetAmount"`
	TakerFee         string `json:"takerFee"`
	Sender           string `json:"sender"`
	FeeRecipient     string `json:"feeRecipient"`
	Expiry           string `json:"expiry"`
	Salt             string `json:"salt"`
	Pool             string `json:"pool"`
	Origin           string `json:"origin"`
	Signature        string `json:"signature"`
}

// FIXME
//
// MarshalJSON implements a custom JSON marshaller for the SignedOrder type
func (s SignedOrder) MarshalJSON() ([]byte, error) {
	switch o := s.Order.(type) {
	case *OrderV3:
		makerAssetData := "0x"
		if len(o.MakerAssetData) != 0 {
			makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.MakerAssetData))
		}
		// Note(albrow): Because of how our smart contracts work, most fields of an
		// order cannot be null. However, makerAssetFeeData and takerAssetFeeData are
		// the exception. For these fields, "0x" is used to indicate a null value.
		makerFeeAssetData := "0x"
		if len(o.MakerFeeAssetData) != 0 {
			makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.MakerFeeAssetData))
		}
		takerAssetData := "0x"
		if len(o.TakerAssetData) != 0 {
			takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.TakerAssetData))
		}
		takerFeeAssetData := "0x"
		if len(o.TakerFeeAssetData) != 0 {
			takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.TakerFeeAssetData))
		}
		signature := "0x"
		if len(s.Signature) != 0 {
			signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
		}

		signedOrderBytes, err := json.Marshal(SignedOrderV3JSON{
			ChainID:               o.ChainID.Int64(),
			ExchangeAddress:       strings.ToLower(o.ExchangeAddress.Hex()),
			MakerAddress:          strings.ToLower(o.MakerAddress.Hex()),
			MakerAssetData:        makerAssetData,
			MakerFeeAssetData:     makerFeeAssetData,
			MakerAssetAmount:      o.MakerAssetAmount.String(),
			MakerFee:              o.MakerFee.String(),
			TakerAddress:          strings.ToLower(o.TakerAddress.Hex()),
			TakerAssetData:        takerAssetData,
			TakerFeeAssetData:     takerFeeAssetData,
			TakerAssetAmount:      o.TakerAssetAmount.String(),
			TakerFee:              o.TakerFee.String(),
			SenderAddress:         strings.ToLower(o.SenderAddress.Hex()),
			FeeRecipientAddress:   strings.ToLower(o.FeeRecipientAddress.Hex()),
			ExpirationTimeSeconds: o.ExpirationTimeSeconds.String(),
			Salt:                  o.Salt.String(),
			Signature:             signature,
		})
		return signedOrderBytes, err
	case *OrderV4:
		signature := "0x"
		if len(s.Signature) != 0 {
			signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
		}

		signedOrderBytes, err := json.Marshal(SignedOrderV4JSON{
			ChainID:          o.ChainID.Int64(),
			Exchange:         strings.ToLower(o.Exchange.Hex()),
			Maker:            strings.ToLower(o.Maker.Hex()),
			MakerToken:       strings.ToLower(o.MakerToken.Hex()),
			MakerAssetAmount: o.MakerAssetAmount.String(),
			MakerFee:         o.MakerFee.String(),
			Taker:            strings.ToLower(o.Taker.Hex()),
			TakerToken:       strings.ToLower(o.TakerToken.Hex()),
			TakerAssetAmount: o.TakerAssetAmount.String(),
			TakerFee:         o.TakerFee.String(),
			Sender:           strings.ToLower(o.Sender.Hex()),
			FeeRecipient:     strings.ToLower(o.FeeRecipient.Hex()),
			Expiry:           o.Expiry.String(),
			Salt:             o.Salt.String(),
			Signature:        signature,
		})
		return signedOrderBytes, err
	default:
		return nil, errors.New("Can't marshal unrecognized order type")
	}
}

type versionJSON struct {
	version int `json:"version"`
}

type signedOrderV3JSON struct {
	Version               int64  `json:"version"`
	ChainID               int64  `json:"chainId"`
	ExchangeAddress       string `json:"exchangeAddress"`
	MakerAddress          string `json:"makerAddress"`
	MakerAssetData        string `json:"makerAssetData"`
	MakerFeeAssetData     string `json:"makerFeeAssetData"`
	MakerAssetAmount      string `json:"makerAssetAmount"`
	MakerFee              string `json:"makerFee"`
	TakerAddress          string `json:"takerAddress"`
	TakerAssetData        string `json:"takerAssetData"`
	TakerFeeAssetData     string `json:"takerFeeAssetData"`
	TakerAssetAmount      string `json:"takerAssetAmount"`
	TakerFee              string `json:"takerFee"`
	SenderAddress         string `json:"senderAddress"`
	FeeRecipientAddress   string `json:"feeRecipientAddress"`
	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
	Salt                  string `json:"salt"`
	Signature             string `json:"signature"`
}

type signedOrderV4JSON struct {
	Version          int64  `json:"version"`
	ChainID          int64  `json:"chainId"`
	Exchange         string `json:"exchange"`
	Maker            string `json:"maker"`
	MakerToken       string `json:"makerToken"`
	MakerAssetAmount string `json:"makerAssetAmount"`
	MakerFee         string `json:"makerFee"`
	Taker            string `json:"taker"`
	TakerToken       string `json:"takerToken"`
	TakerAssetAmount string `json:"takerAssetAmount"`
	TakerFee         string `json:"takerFee"`
	Sender           string `json:"sender"`
	FeeRecipient     string `json:"feeRecipient"`
	Expiry           string `json:"expiry"`
	Salt             string `json:"salt"`
	Pool             string `json:"pool"`
	Origin           string `json:"origin"`
	Signature        string `json:"signature"`
}

// FIXME(jalextowle): I think we'll need to unmarshal in stages. First, we'll unmarshal
// enough to see what the version is. After that, we'll unmarshal again, using that
// knowledge to guide our decoding.
//
// UnmarshalJSON implements a custom JSON unmarshaller for the SignedOrder type
func (s *SignedOrder) UnmarshalJSON(data []byte) error {
	var v versionJSON
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	switch v.version {
	case 3:
		var signedOrder signedOrderV3JSON
		err := json.Unmarshal(data, &signedOrder)
		if err != nil {
			return err
		}
		var ok bool
		o := &OrderV3{}
		o.ChainID = big.NewInt(signedOrder.ChainID)
		o.ExchangeAddress = common.HexToAddress(signedOrder.ExchangeAddress)
		o.MakerAddress = common.HexToAddress(signedOrder.MakerAddress)
		o.MakerAssetData = common.FromHex(signedOrder.MakerAssetData)
		o.MakerFeeAssetData = common.FromHex(signedOrder.MakerFeeAssetData)
		if signedOrder.MakerAssetAmount != "" {
			o.MakerAssetAmount, ok = math.ParseBig256(signedOrder.MakerAssetAmount)
			if !ok {
				o.MakerAssetAmount = nil
			}
		}
		if signedOrder.MakerFee != "" {
			o.MakerFee, ok = math.ParseBig256(signedOrder.MakerFee)
			if !ok {
				o.MakerFee = nil
			}
		}
		o.TakerAddress = common.HexToAddress(signedOrder.TakerAddress)
		o.TakerAssetData = common.FromHex(signedOrder.TakerAssetData)
		o.TakerFeeAssetData = common.FromHex(signedOrder.TakerFeeAssetData)
		if signedOrder.TakerAssetAmount != "" {
			o.TakerAssetAmount, ok = math.ParseBig256(signedOrder.TakerAssetAmount)
			if !ok {
				o.TakerAssetAmount = nil
			}
		}
		if signedOrder.TakerFee != "" {
			o.TakerFee, ok = math.ParseBig256(signedOrder.TakerFee)
			if !ok {
				o.TakerFee = nil
			}
		}
		o.SenderAddress = common.HexToAddress(signedOrder.SenderAddress)
		o.FeeRecipientAddress = common.HexToAddress(signedOrder.FeeRecipientAddress)
		if signedOrder.ExpirationTimeSeconds != "" {
			o.ExpirationTimeSeconds, ok = math.ParseBig256(signedOrder.ExpirationTimeSeconds)
			if !ok {
				o.ExpirationTimeSeconds = nil
			}
		}
		if signedOrder.Salt != "" {
			o.Salt, ok = math.ParseBig256(signedOrder.Salt)
			if !ok {
				o.Salt = nil
			}
		}
		s.Order = o
		s.Version = signedOrder.Version
		s.Signature = common.FromHex(signedOrder.Signature)
		return nil
	case 4:
		var signedOrder signedOrderV4JSON
		err := json.Unmarshal(data, &signedOrder)
		if err != nil {
			return err
		}
		var ok bool
		o := &OrderV4{}
		o.ChainID = big.NewInt(signedOrder.ChainID)
		o.Exchange = common.HexToAddress(signedOrder.Exchange)
		o.Maker = common.HexToAddress(signedOrder.Maker)
		o.MakerToken = common.HexToAddress(signedOrder.MakerToken)
		if signedOrder.MakerFee != "" {
			o.MakerFee, ok = math.ParseBig256(signedOrder.MakerFee)
			if !ok {
				o.MakerFee = nil
			}
		}
		o.Taker = common.HexToAddress(signedOrder.Taker)
		o.TakerToken = common.HexToAddress(signedOrder.TakerToken)
		if signedOrder.TakerFee != "" {
			o.TakerFee, ok = math.ParseBig256(signedOrder.TakerFee)
			if !ok {
				o.TakerFee = nil
			}
		}
		o.Sender = common.HexToAddress(signedOrder.Sender)
		o.FeeRecipient = common.HexToAddress(signedOrder.FeeRecipient)
		if signedOrder.Expiry != "" {
			o.Expiry, ok = math.ParseBig256(signedOrder.Expiry)
			if !ok {
				o.Expiry = nil
			}
		}
		if signedOrder.Salt != "" {
			o.Salt, ok = math.ParseBig256(signedOrder.Salt)
			if !ok {
				o.Salt = nil
			}
		}
		if signedOrder.Pool != "" {
			o.Salt, ok = math.ParseBig256(signedOrder.Pool)
			if !ok {
				o.Pool = nil
			}
		}
		if signedOrder.Origin != "" {
			o.Salt, ok = math.ParseBig256(signedOrder.Origin)
			if !ok {
				o.Pool = nil
			}
		}
		s.Order = o
		s.Version = signedOrder.Version
		s.Signature = common.FromHex(signedOrder.Signature)
		return nil
	default:
		return errors.New("")
	}
}

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}
