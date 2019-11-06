package zeroex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/0xProject/0x-mesh/ethereum/signer"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	gethsigner "github.com/ethereum/go-ethereum/signer/core"
	"golang.org/x/crypto/sha3"
)

// Order represents an unsigned 0x order
type Order struct {
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

// SignedOrder represents a signed 0x order
type SignedOrder struct {
	Order
	Signature []byte `json:"signature"`
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

// OrderEvent is the order event emitted by Mesh nodes on the "orders" topic
// when calling JSON-RPC method `mesh_subscribe`
type OrderEvent struct {
	OrderHash                common.Hash        `json:"orderHash"`
	SignedOrder              *SignedOrder       `json:"signedOrder"`
	EndState                 OrderEventEndState `json:"endState"`
	FillableTakerAssetAmount *big.Int           `json:"fillableTakerAssetAmount"`
	// All the contract events that triggered this orders re-evaluation. They did not
	// all necessarily cause the orders state change itself, only it's re-evaluation.
	// Since it's state _did_ change, at least one of them did cause the actual state change.
	ContractEvents []*ContractEvent `json:"contractEvents"`
}

type orderEventJSON struct {
	OrderHash                string               `json:"orderHash"`
	SignedOrder              *SignedOrder         `json:"signedOrder"`
	EndState                 string               `json:"endState"`
	FillableTakerAssetAmount string               `json:"fillableTakerAssetAmount"`
	ContractEvents           []*contractEventJSON `json:"contractEvents"`
}

// MarshalJSON implements a custom JSON marshaller for the OrderEvent type
func (o OrderEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder,
		"endState":                 o.EndState,
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
		"contractEvents":           o.ContractEvents,
	})
}

// UnmarshalJSON implements a custom JSON unmarshaller for the OrderEvent type
func (o *OrderEvent) UnmarshalJSON(data []byte) error {
	var orderEventJSON orderEventJSON
	err := json.Unmarshal(data, &orderEventJSON)
	if err != nil {
		return err
	}
	return o.fromOrderEventJSON(orderEventJSON)
}

func (o *OrderEvent) fromOrderEventJSON(orderEventJSON orderEventJSON) error {
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
	ESInvalid          = OrderEventEndState("INVALID")
	ESOrderAdded       = OrderEventEndState("ADDED")
	ESOrderFilled      = OrderEventEndState("FILLED")
	ESOrderFullyFilled = OrderEventEndState("FULLY_FILLED")
	ESOrderCancelled   = OrderEventEndState("CANCELLED")
	ESOrderExpired     = OrderEventEndState("EXPIRED")
	ESOrderUnexpired   = OrderEventEndState("UNEXPIRED")
	// An order becomes unfunded if the maker transfers the balance / changes their
	// allowance backing an order
	ESOrderBecameUnfunded = OrderEventEndState("UNFUNDED")
	// Fillability for an order can increase if a previously processed fill event
	// gets reverted, or if a maker tops up their balance/allowance backing an order
	ESOrderFillabilityIncreased = OrderEventEndState("FILLABILITY_INCREASED")
	// Order is potentially still valid but was removed for a different reason
	// (e.g. the database is full or the peer that sent the order was
	// misbehaving). The order will no longer be watched and no further events for
	// this order will be emitted. In some cases, the order may be re-added in the
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
func (o *Order) ResetHash() {
	o.hash = nil
}

// ComputeOrderHash computes a 0x order hash
func (o *Order) ComputeOrderHash() (common.Hash, error) {
	if o.hash != nil {
		return *o.hash, nil
	}

	var domain = gethsigner.TypedDataDomain{
		Name:              "0x Protocol",
		Version:           "3.0.0",
		ChainId:           o.ChainID,
		VerifyingContract: o.ExchangeAddress.Hex(),
	}

	fmt.Println("domain", domain)

	var message = map[string]interface{}{
		"makerAddress":          o.MakerAddress.Hex(),
		"takerAddress":          o.TakerAddress.Hex(),
		"senderAddress":         o.SenderAddress.Hex(),
		"feeRecipientAddress":   o.FeeRecipientAddress.Hex(),
		"makerAssetData":        o.MakerAssetData,
		"makerFeeAssetData":     o.MakerFeeAssetData,
		"takerAssetData":        o.TakerAssetData,
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
func SignOrder(signer signer.Signer, order *Order) (*SignedOrder, error) {
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
func SignTestOrder(order *Order) (*SignedOrder, error) {
	testSigner := signer.NewTestSigner()
	signedOrder, err := SignOrder(testSigner, order)
	if err != nil {
		return nil, err
	}
	return signedOrder, nil
}

// Trim converts the order to a TrimmedOrder, which is the format expected by
// our smart contracts. It removes the ChainID and ExchangeAddress fields.
func (s *SignedOrder) Trim() wrappers.TrimmedOrder {
	return wrappers.TrimmedOrder{
		MakerAddress:          s.MakerAddress,
		TakerAddress:          s.TakerAddress,
		FeeRecipientAddress:   s.FeeRecipientAddress,
		SenderAddress:         s.SenderAddress,
		MakerAssetAmount:      s.MakerAssetAmount,
		TakerAssetAmount:      s.TakerAssetAmount,
		MakerFee:              s.MakerFee,
		TakerFee:              s.TakerFee,
		ExpirationTimeSeconds: s.ExpirationTimeSeconds,
		Salt:                  s.Salt,
		MakerAssetData:        s.MakerAssetData,
		MakerFeeAssetData:     s.MakerFeeAssetData,
		TakerAssetData:        s.TakerAssetData,
		TakerFeeAssetData:     s.TakerFeeAssetData,
	}
}

// SignedOrderJSON is an unmodified JSON representation of a SignedOrder
type SignedOrderJSON struct {
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

// MarshalJSON implements a custom JSON marshaller for the SignedOrder type
func (s SignedOrder) MarshalJSON() ([]byte, error) {
	makerAssetData := ""
	if len(s.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerAssetData))
	}
	makerFeeAssetData := "0x"
	if len(s.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerFeeAssetData))
	}
	takerAssetData := ""
	if len(s.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerAssetData))
	}
	takerFeeAssetData := "0x"
	if len(s.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerFeeAssetData))
	}
	signature := ""
	if len(s.Signature) != 0 {
		signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
	}

	signedOrderBytes, err := json.Marshal(SignedOrderJSON{
		ChainID:               s.ChainID.Int64(),
		ExchangeAddress:       strings.ToLower(s.ExchangeAddress.Hex()),
		MakerAddress:          strings.ToLower(s.MakerAddress.Hex()),
		MakerAssetData:        makerAssetData,
		MakerFeeAssetData:     makerFeeAssetData,
		MakerAssetAmount:      s.MakerAssetAmount.String(),
		MakerFee:              s.MakerFee.String(),
		TakerAddress:          strings.ToLower(s.TakerAddress.Hex()),
		TakerAssetData:        takerAssetData,
		TakerFeeAssetData:     takerFeeAssetData,
		TakerAssetAmount:      s.TakerAssetAmount.String(),
		TakerFee:              s.TakerFee.String(),
		SenderAddress:         strings.ToLower(s.SenderAddress.Hex()),
		FeeRecipientAddress:   strings.ToLower(s.FeeRecipientAddress.Hex()),
		ExpirationTimeSeconds: s.ExpirationTimeSeconds.String(),
		Salt:                  s.Salt.String(),
		Signature:             signature,
	})
	return signedOrderBytes, err
}

const addressHexLength = 42

// UnmarshalJSON implements a custom JSON unmarshaller for the SignedOrder type
func (s *SignedOrder) UnmarshalJSON(data []byte) error {
	var signedOrderJSON SignedOrderJSON
	err := json.Unmarshal(data, &signedOrderJSON)
	if err != nil {
		return err
	}
	var ok bool
	s.ChainID = big.NewInt(signedOrderJSON.ChainID)
	s.ExchangeAddress = common.HexToAddress(signedOrderJSON.ExchangeAddress)
	s.MakerAddress = common.HexToAddress(signedOrderJSON.MakerAddress)
	s.MakerAssetData = common.FromHex(signedOrderJSON.MakerAssetData)
	s.MakerFeeAssetData = common.FromHex(signedOrderJSON.MakerFeeAssetData)
	if signedOrderJSON.MakerAssetAmount != "" {
		s.MakerAssetAmount, ok = math.ParseBig256(signedOrderJSON.MakerAssetAmount)
		if !ok {
			s.MakerAssetAmount = nil
		}
	}
	if signedOrderJSON.MakerFee != "" {
		s.MakerFee, ok = math.ParseBig256(signedOrderJSON.MakerFee)
		if !ok {
			s.MakerFee = nil
		}
	}
	s.TakerAddress = common.HexToAddress(signedOrderJSON.TakerAddress)
	s.TakerAssetData = common.FromHex(signedOrderJSON.TakerAssetData)
	s.TakerFeeAssetData = common.FromHex(signedOrderJSON.TakerFeeAssetData)
	if signedOrderJSON.TakerAssetAmount != "" {
		s.TakerAssetAmount, ok = math.ParseBig256(signedOrderJSON.TakerAssetAmount)
		if !ok {
			s.TakerAssetAmount = nil
		}
	}
	if signedOrderJSON.TakerFee != "" {
		s.TakerFee, ok = math.ParseBig256(signedOrderJSON.TakerFee)
		if !ok {
			s.TakerFee = nil
		}
	}
	s.SenderAddress = common.HexToAddress(signedOrderJSON.SenderAddress)
	s.FeeRecipientAddress = common.HexToAddress(signedOrderJSON.FeeRecipientAddress)
	if signedOrderJSON.ExpirationTimeSeconds != "" {
		s.ExpirationTimeSeconds, ok = math.ParseBig256(signedOrderJSON.ExpirationTimeSeconds)
		if !ok {
			s.ExpirationTimeSeconds = nil
		}
	}
	if signedOrderJSON.Salt != "" {
		s.Salt, ok = math.ParseBig256(signedOrderJSON.Salt)
		if !ok {
			s.Salt = nil
		}
	}
	s.Signature = common.FromHex(signedOrderJSON.Signature)
	return nil
}

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}
