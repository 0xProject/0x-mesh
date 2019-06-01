package zeroex

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/common"
	signer "github.com/ethereum/go-ethereum/signer/core"
	"golang.org/x/crypto/sha3"
)

// Order represents an unsigned 0x order
type Order struct {
	MakerAddress          common.Address `json:"makerAddress"`
	MakerAssetData        []byte         `json:"makerAssetData"`
	MakerAssetAmount      *big.Int       `json:"makerAssetAmount"`
	MakerFee              *big.Int       `json:"makerFee"`
	TakerAddress          common.Address `json:"takerAddress"`
	TakerAssetData        []byte         `json:"takerAssetData"`
	TakerAssetAmount      *big.Int       `json:"takerAssetAmount"`
	TakerFee              *big.Int       `json:"takerFee"`
	SenderAddress         common.Address `json:"senderAddress"`
	ExchangeAddress       common.Address `json:"exchangeAddress"`
	FeeRecipientAddress   common.Address `json:"feeRecipientAddress"`
	ExpirationTimeSeconds *big.Int       `json:"expirationTimeSeconds"`
	Salt                  *big.Int       `json:"salt"`
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

// OrderEvent is the order event emitted by Mesh nodes on the "orders" topic
// when calling JSON-RPC method `mesh_subscribe`
type OrderEvent struct {
	OrderHash                common.Hash    `json:"orderHash"`
	SignedOrder              *SignedOrder   `json:"signedOrder"`
	Kind                     OrderEventKind `json:"kind"`
	FillableTakerAssetAmount *big.Int       `json:"fillableTakerAssetAmount"`
	// The hash of the Ethereum transaction that caused the order status to change
	TxHash common.Hash `json:"txHash"`
}

// OrderEventKind enumerates all the possible order event types
type OrderEventKind string

// OrderEventKind values
const (
	EKInvalid          = OrderEventKind("INVALID")
	EKOrderAdded       = OrderEventKind("ADDED")
	EKOrderFilled      = OrderEventKind("FILLED")
	EKOrderFullyFilled = OrderEventKind("FULLY_FILLED")
	EKOrderCancelled   = OrderEventKind("CANCELLED")
	EKOrderExpired     = OrderEventKind("EXPIRED")
	// An order becomes unfunded if the maker transfers the balance / changes their
	// allowance backing an order
	EKOrderBecameUnfunded = OrderEventKind("UNFUNDED")
	// Fillability for an order can increase if a previously processed fill event
	// gets reverted, or if a maker tops up their balance/allowance backing an order
	EKOrderFillabilityIncreased = OrderEventKind("FILLABILITY_INCREASED")
)

var eip712OrderTypes = signer.Types{
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
	},
}

// ComputeOrderHash computes a 0x order hash
func (o *Order) ComputeOrderHash() (common.Hash, error) {
	var domain = signer.TypedDataDomain{
		Name:              "0x Protocol",
		Version:           "2",
		VerifyingContract: o.ExchangeAddress.Hex(),
	}

	var message = map[string]interface{}{
		"makerAddress":          o.MakerAddress.Hex(),
		"takerAddress":          o.TakerAddress.Hex(),
		"senderAddress":         o.SenderAddress.Hex(),
		"feeRecipientAddress":   o.FeeRecipientAddress.Hex(),
		"makerAssetData":        o.MakerAssetData,
		"takerAssetData":        o.TakerAssetData,
		"salt":                  o.Salt,
		"makerFee":              o.MakerFee,
		"takerFee":              o.TakerFee,
		"makerAssetAmount":      o.MakerAssetAmount,
		"takerAssetAmount":      o.TakerAssetAmount,
		"expirationTimeSeconds": o.ExpirationTimeSeconds,
	}

	var typedData = signer.TypedData{
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
	return hash, nil
}

// SignOrder signs the 0x order with the supplied Signer
func SignOrder(signer ethereum.Signer, order *Order) (*SignedOrder, error) {
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
	testSigner := ethereum.NewTestSigner()
	signedOrder, err := SignOrder(testSigner, order)
	if err != nil {
		return nil, err
	}
	return signedOrder, nil
}

// ConvertToOrderWithoutExchangeAddress re-formats a SignedOrder into the format expected by the 0x
// smart contracts.
func (s *SignedOrder) ConvertToOrderWithoutExchangeAddress() wrappers.OrderWithoutExchangeAddress {
	orderWithoutExchangeAddress := wrappers.OrderWithoutExchangeAddress{
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
		TakerAssetData:        s.TakerAssetData,
	}
	return orderWithoutExchangeAddress
}

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}
