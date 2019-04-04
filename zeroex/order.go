package zeroex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SignedOrder struct {
	MakerAddress           common.Address
	MakerAssetData         []byte
	MakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerAddress           common.Address
	TakerAssetData         []byte
	TakerAssetFilledAmount *big.Int
	TakerFeePaid           *big.Int
	SenderAddress          common.Address
	ExchangeAddress        common.Address
	FeeRecipientAddress    common.Address
	ExpirationTimeSeconds  *big.Int
	Salt                   *big.Int
	Signature              []byte
}
