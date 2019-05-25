package zeroex

import (
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	signer "github.com/ethereum/go-ethereum/signer/core"
	"golang.org/x/crypto/sha3"
)

// OrderStatus represents the status of an order as returned from the 0x smart contracts
// as part of OrderInfo
type OrderStatus uint8

// OrderStatus values
const (
	Invalid OrderStatus = iota
	InvalidMakerAssetAmount
	InvalidTakerAssetAmount
	Fillable
	Expired
	FullyFilled
	Cancelled
	SignatureInvalid
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
	*Order
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

func (o *Order) ecSign(rpcClient *rpc.Client) ([]byte, error) {
	orderHash, err := o.ComputeOrderHash()
	if err != nil {
		return nil, err
	}

	var signer ethereum.Signer
	if rpcClient == nil {
		signer = ethereum.NewTestSigner()
	} else {
		signer = ethereum.NewEthRPCSigner(rpcClient)
	}
	ecSignature, err := signer.EthSign(orderHash.Bytes(), o.MakerAddress)
	if err != nil {
		return nil, err
	}

	// Generate 0x EthSign Signature (append the signature type byte)
	signature := make([]byte, 66)
	signature[0] = ecSignature.V
	copy(signature[1:33], ecSignature.R[:])
	copy(signature[33:65], ecSignature.S[:])
	signature[65] = byte(EthSignSignature)
	return signature, nil
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

// SignTestOrder converts a test 0x order to a signed 0x order
func SignTestOrder(order *Order) (*SignedOrder, error) {
	signature, err := order.ecSign(nil)
	if err != nil {
		return nil, err
	}
	signedOrder := &SignedOrder{
		Order:     order,
		Signature: signature,
	}
	return signedOrder, nil
}

// SignOrder converts a 0x order to a signed 0x order by fetching an order
// signature via an `eth_sign` call
func SignOrder(order *Order, rpcClient *rpc.Client) (*SignedOrder, error) {
	signature, err := order.ecSign(rpcClient)
	if err != nil {
		return nil, err
	}
	signedOrder := &SignedOrder{
		Order:     order,
		Signature: signature,
	}
	return signedOrder, nil
}

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}
