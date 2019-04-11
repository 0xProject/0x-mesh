package zeroex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// import (
// 	"fmt"
// 	"math/big"

// 	signer "github.com/ethereum/go-ethereum/signer/core"

// 	"github.com/ethereum/go-ethereum/common"
// 	"golang.org/x/crypto/sha3"
// )

// SignedOrder represents a signed 0x order
type SignedOrder struct {
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
	Signature             []byte         `json:"signature"`
}

// var eip712OrderTypes = signer.Types{
// 	"EIP712Domain": {
// 		{
// 			Name: "name",
// 			Type: "string",
// 		},
// 		{
// 			Name: "version",
// 			Type: "string",
// 		},
// 		{
// 			Name: "verifyingContract",
// 			Type: "address",
// 		},
// 	},
// 	"Order": {
// 		{
// 			Name: "makerAddress",
// 			Type: "address",
// 		},
// 		{
// 			Name: "takerAddress",
// 			Type: "address",
// 		},
// 		{
// 			Name: "feeRecipientAddress",
// 			Type: "address",
// 		},
// 		{
// 			Name: "senderAddress",
// 			Type: "address",
// 		},
// 		{
// 			Name: "makerAssetAmount",
// 			Type: "uint256",
// 		},
// 		{
// 			Name: "takerAssetAmount",
// 			Type: "uint256",
// 		},
// 		{
// 			Name: "makerFee",
// 			Type: "uint256",
// 		},
// 		{
// 			Name: "takerFee",
// 			Type: "uint256",
// 		},
// 		{
// 			Name: "expirationTimeSeconds",
// 			Type: "uint256",
// 		},
// 		{
// 			Name: "salt",
// 			Type: "uint256",
// 		},
// 		{
// 			Name: "makerAssetData",
// 			Type: "bytes",
// 		},
// 		{
// 			Name: "takerAssetData",
// 			Type: "bytes",
// 		},
// 	},
// }

// // ComputeOrderHash computes a 0x order hash
// func (s *SignedOrder) ComputeOrderHash() (common.Hash, error) {
// 	var domain = signer.TypedDataDomain{
// 		Name:              "0x Protocol",
// 		Version:           "2",
// 		VerifyingContract: s.ExchangeAddress.Hex(),
// 	}

// 	var message = map[string]interface{}{
// 		"makerAddress":          s.MakerAddress.Hex(),
// 		"takerAddress":          s.TakerAddress.Hex(),
// 		"senderAddress":         s.SenderAddress.Hex(),
// 		"feeRecipientAddress":   s.FeeRecipientAddress.Hex(),
// 		"makerAssetData":        s.MakerAssetData,
// 		"takerAssetData":        s.TakerAssetData,
// 		"salt":                  s.Salt,
// 		"makerFee":              s.MakerFee,
// 		"takerFee":              s.TakerFee,
// 		"makerAssetAmount":      s.MakerAssetAmount,
// 		"takerAssetAmount":      s.TakerAssetAmount,
// 		"expirationTimeSeconds": s.ExpirationTimeSeconds,
// 	}

// 	var typedData = signer.TypedData{
// 		Types:       eip712OrderTypes,
// 		PrimaryType: "Order",
// 		Domain:      domain,
// 		Message:     message,
// 	}

// 	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
// 	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
// 	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
// 	hashBytes := keccak256(rawData)
// 	hash := common.BytesToHash(hashBytes)
// 	return hash, nil
// }

// // keccak256 calculates and returns the Keccak256 hash of the input data.
// func keccak256(data ...[]byte) []byte {
// 	d := sha3.NewLegacyKeccak256()
// 	for _, b := range data {
// 		d.Write(b)
// 	}
// 	return d.Sum(nil)
// }
