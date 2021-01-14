package zeroex

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/ethereum/signer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	gethsigner "github.com/ethereum/go-ethereum/signer/core"
)

// OrderV4 represents an unsigned 0x v4 limit order
// V4 Protocol also has RFQ orders, these are
// See <https://0xprotocol.readthedocs.io/en/latest/basics/orders.html#limit-orders>
type OrderV4 struct {
	// Domain information
	// TODO: These are constant within a chain context (mainnet/testnet/etc)
	// probably best to keep them out of the order struct, but this is how V3
	// does it.
	ChainID         *big.Int       `json:"chainId"`
	ExchangeAddress common.Address `json:"exchangeAddress"`

	// Limit order values
	MakerToken          common.Address `json:"makerToken"`
	TakerToken          common.Address `json:"takerToken"`
	MakerAmount         *big.Int       `json:"makerAmount"`         // uint128
	TakerAmount         *big.Int       `json:"takerAmount"`         // uint128
	TakerTokenFeeAmount *big.Int       `json:"takerTokenFeeAmount"` // uint128
	Maker               common.Address `json:"takerAddress"`
	Taker               common.Address `json:"makerAddress"`
	Sender              common.Address `json:"sender"`
	FeeRecipient        common.Address `json:"feeRecipient"`
	Pool                Bytes32        `json:"pool"`   // bytes32
	Expiry              *big.Int       `json:"expiry"` // uint64
	Salt                *big.Int       `json:"salt"`   // uint256

	// Cache hash for performance
	hash *common.Hash
}

// SignatureTypeV4 represents the type of 0x signature encountered
type SignatureTypeV4 uint8

// SignedOrderV4 represents a signed 0x order
// See <https://0xprotocol.readthedocs.io/en/latest/basics/orders.html#how-to-sign>
type SignedOrderV4 struct {
	OrderV4
	SignatureTypeV4 `json:"signatureType"` // uint8
	V               uint8                  `json:"v"` // uint8
	R               *big.Int               `json:"r"` // uint256
	S               *big.Int               `json:"s"` // uint256
}

// SignatureType values
const (
	IllegalSignatureV4 SignatureTypeV4 = iota
	InvalidSignatureV4
	EIP712SignatureV4
	EthSignSignatureV4
)

////////////////////////////////////////////////////////////////////////////////
//  O R D E R   H A S H I N G
////////////////////////////////////////////////////////////////////////////////

// See <https://0xprotocol.readthedocs.io/en/latest/basics/functions.html#getlimitorderhash>
// See <https://github.com/0xProject/protocol/blob/682c07cb73e572929547ee65afb0ce79885a7828/packages/protocol-utils/src/orders.ts#L127>
var eip712OrderTypesV4 = gethsigner.Types{
	"EIP712Domain": {
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"LimitOrder": {
		{Name: "makerToken", Type: "address"},
		{Name: "takerToken", Type: "address"},
		{Name: "makerAmount", Type: "uint128"},
		{Name: "takerAmount", Type: "uint128"},
		{Name: "takerTokenFeeAmount", Type: "uint128"},
		{Name: "maker", Type: "address"},
		{Name: "taker", Type: "address"},
		{Name: "sender", Type: "address"},
		{Name: "feeRecipient", Type: "address"},
		{Name: "pool", Type: "bytes32"},
		{Name: "expiry", Type: "uint64"},
		{Name: "salt", Type: "uint256"},
	},
}

// ResetHash resets the cached order hash. Usually only required for testing.
func (o *OrderV4) ResetHash() {
	o.hash = nil
}

// ComputeOrderHash computes a 0x order hash
func (o *OrderV4) ComputeOrderHash() (common.Hash, error) {
	if o.hash != nil {
		return *o.hash, nil
	}

	// TODO: This domain is constant for a given environment and should probably
	// not depend on the order.
	chainID := math.NewHexOrDecimal256(o.ChainID.Int64())
	var domain = gethsigner.TypedDataDomain{
		Name:              "ZeroEx",
		Version:           "1.0.0",
		ChainId:           chainID,
		VerifyingContract: o.ExchangeAddress.Hex(),
	}

	var message = map[string]interface{}{
		"makerToken":          o.MakerToken.Hex(),
		"takerToken":          o.TakerToken.Hex(),
		"makerAmount":         o.MakerAmount.String(),
		"takerAmount":         o.TakerAmount.String(),
		"takerTokenFeeAmount": o.TakerTokenFeeAmount.String(),
		"taker":               o.Taker.Hex(),
		"maker":               o.Maker.Hex(),
		"sender":              o.Sender.Hex(),
		"feeRecipient":        o.FeeRecipient.Hex(),
		"pool":                o.Pool.Hex(),
		"expiry":              o.Expiry.String(),
		"salt":                o.Salt.String(),
	}

	var typedData = gethsigner.TypedData{
		Types:       eip712OrderTypesV4,
		PrimaryType: "LimitOrder",
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

////////////////////////////////////////////////////////////////////////////////
//  O R D E R   S I G N I N G
////////////////////////////////////////////////////////////////////////////////

// SignOrderV4 signs the 0x order with the supplied Signer
func SignOrderV4(signer signer.Signer, order *OrderV4) (*SignedOrderV4, error) {
	if order == nil {
		return nil, errors.New("cannot sign nil order")
	}
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		return nil, err
	}

	ecSignature, err := signer.EthSign(orderHash.Bytes(), order.Maker)
	if err != nil {
		return nil, err
	}

	// Generate 0x V4 Signature
	signedOrder := &SignedOrderV4{
		OrderV4:         *order,
		SignatureTypeV4: EIP712SignatureV4,
		V:               ecSignature.V,
		R:               ecSignature.R.Big(),
		S:               ecSignature.S.Big(),
	}
	return signedOrder, nil
}

// SignTestOrderV4 signs the 0x order with the local test signer
func SignTestOrderV4(order *OrderV4) (*SignedOrderV4, error) {
	testSigner := signer.NewTestSigner()
	signedOrder, err := SignOrderV4(testSigner, order)
	if err != nil {
		return nil, err
	}
	return signedOrder, nil
}
