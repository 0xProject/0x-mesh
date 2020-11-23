package zeroex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type V3Order struct {
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

// ResetHash resets the cached order hash. Usually only required for testing.
func (o *V3Order) resetHash() {
	o.hash = nil
}

type SignedV3Order struct {
	V3Order
	Signature []byte `json:"signature"`
}
