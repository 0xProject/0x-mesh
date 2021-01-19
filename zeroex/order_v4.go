package zeroex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

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
	Maker               common.Address `json:"makerAddress"`
	Taker               common.Address `json:"takerAddress"`
	Sender              common.Address `json:"sender"`
	FeeRecipient        common.Address `json:"feeRecipient"`
	Pool                Bytes32        `json:"pool"`   // bytes32
	Expiry              *big.Int       `json:"expiry"` // uint64
	Salt                *big.Int       `json:"salt"`   // uint256

	// Cache hash for performance
	hash *common.Hash
}

type SignatureFieldV4 struct {
	SignatureType SignatureTypeV4 `json:"signatureType"`
	V             uint8           `json:"v"`
	R             Bytes32         `json:"r"`
	S             Bytes32         `json:"s"`
}

// SignatureTypeV4 represents the type of 0x signature encountered
type SignatureTypeV4 uint8

// SignedOrderV4 represents a signed 0x order
// See <https://0xprotocol.readthedocs.io/en/latest/basics/orders.html#how-to-sign>
type SignedOrderV4 struct {
	OrderV4   `json:"order"`
	Signature SignatureFieldV4 `json:"signature"`
}

// SignedOrderJSONV4 is an unmodified JSON representation of a SignedOrder
type SignedOrderJSONV4 struct {
	ChainID             int64  `json:"chainId"`
	ExchangeAddress     string `json:"exchangeAddress"`
	MakerToken          string `json:"makerToken"`
	TakerToken          string `json:"takerToken"`
	MakerAmount         string `json:"makerAmount"`
	TakerAmount         string `json:"takerAmount"`
	TakerTokenFeeAmount string `json:"takerTokenFeeAmount"`
	Maker               string `json:"maker"`
	Taker               string `json:"taker"`
	Sender              string `json:"sender"`
	FeeRecipient        string `json:"feeRecipient"`
	Pool                string `json:"pool"`
	Expiry              string `json:"expiry"`
	Salt                string `json:"salt"`
	SignatureType       string `json:"signatureType"`
	SignatureR          string `json:"signatureR"`
	SignatureV          string `json:"signatureV"`
	SignatureS          string `json:"signatureS"`
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
		"pool":                o.Pool.Bytes(),
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
		OrderV4: *order,
		Signature: SignatureFieldV4{
			SignatureType: EthSignSignatureV4,
			V:             ecSignature.V,
			R:             HashToBytes32(ecSignature.R),
			S:             HashToBytes32(ecSignature.S),
		},
	}
	return signedOrder, nil
}

// UnmarshalJSON implements a custom JSON unmarshaller for the SignedOrderV4 type.
func (s *SignedOrderV4) UnmarshalJSON(data []byte) error {
	var signedOrderJSON SignedOrderJSONV4
	err := json.Unmarshal(data, &signedOrderJSON)
	if err != nil {
		return err
	}

	var ok bool
	s.ChainID = big.NewInt(signedOrderJSON.ChainID)
	s.ExchangeAddress = common.HexToAddress(signedOrderJSON.ExchangeAddress)
	s.MakerToken = common.HexToAddress(signedOrderJSON.MakerToken)
	s.TakerToken = common.HexToAddress(signedOrderJSON.TakerToken)
	s.MakerAmount, ok = math.ParseBig256(signedOrderJSON.MakerAmount)
	if !ok {
		s.MakerAmount = nil
	}
	s.TakerAmount, ok = math.ParseBig256(signedOrderJSON.TakerAmount)
	if !ok {
		s.TakerAmount = nil
	}
	s.TakerTokenFeeAmount, ok = math.ParseBig256(signedOrderJSON.TakerTokenFeeAmount)
	if !ok {
		s.TakerTokenFeeAmount = nil
	}
	s.Maker = common.HexToAddress(signedOrderJSON.Maker)
	s.Taker = common.HexToAddress(signedOrderJSON.Taker)
	s.Sender = common.HexToAddress(signedOrderJSON.Sender)
	s.FeeRecipient = common.HexToAddress(signedOrderJSON.FeeRecipient)
	s.Pool = HexToBytes32(signedOrderJSON.Pool)
	s.Expiry, ok = math.ParseBig256(signedOrderJSON.Expiry)
	if !ok {
		s.Expiry = nil
	}
	s.Salt, ok = math.ParseBig256(signedOrderJSON.Salt)
	if !ok {
		s.Expiry = nil
	}
	sigType, err := strconv.ParseUint(signedOrderJSON.SignatureType, 10, 8)
	if err != nil {
		return err
	}
	s.Signature.SignatureType = SignatureTypeV4(sigType)
	sigV, err := strconv.ParseUint(signedOrderJSON.SignatureV, 10, 8)
	if err != nil {
		return err
	}
	s.Signature.V = uint8(sigV)
	s.Signature.R = HexToBytes32(signedOrderJSON.SignatureR)
	s.Signature.S = HexToBytes32(signedOrderJSON.SignatureS)
	return nil
}

func (s *SignedOrderV4) MarshalJSON() ([]byte, error) {
	return json.Marshal(SignedOrderJSONV4{
		ChainID:             s.ChainID.Int64(),
		ExchangeAddress:     strings.ToLower(s.ExchangeAddress.Hex()),
		MakerToken:          strings.ToLower(s.MakerToken.Hex()),
		TakerToken:          strings.ToLower(s.TakerToken.Hex()),
		MakerAmount:         s.MakerAmount.String(),
		TakerAmount:         s.TakerAmount.String(),
		TakerTokenFeeAmount: s.TakerTokenFeeAmount.String(),
		Maker:               strings.ToLower(s.Maker.Hex()),
		Taker:               strings.ToLower(s.Taker.Hex()),
		Sender:              strings.ToLower(s.Sender.Hex()),
		FeeRecipient:        strings.ToLower(s.FeeRecipient.Hex()),
		Pool:                s.Pool.Hex(),
		Expiry:              s.Expiry.String(),
		Salt:                s.Salt.String(),
		SignatureType:       strconv.FormatUint(uint64(s.Signature.SignatureType), 10),
		SignatureR:          s.Signature.R.Hex(),
		SignatureV:          strconv.FormatUint(uint64(s.Signature.V), 10),
		SignatureS:          s.Signature.S.Hex(),
	})
}
