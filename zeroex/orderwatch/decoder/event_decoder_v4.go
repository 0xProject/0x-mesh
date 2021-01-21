package decoder

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ExchangeFillEventV4 struct {
	OrderHash                 common.Hash
	Maker                     common.Address
	Taker                     common.Address
	FeeRecipient              common.Address
	MakerToken                common.Address
	TakerToken                common.Address
	TakerTokenFilledAmount    *big.Int
	MakerTokenFilledAmount    *big.Int
	TakerTokenFeeFilledAmount *big.Int
	ProtocolFeePaid           *big.Int
	Pool                      common.Hash // Decoder does not support bytes32
}

type ExchangeCancelEventV4 struct {
	Maker     common.Address
	OrderHash common.Hash
}

type ExchangePairCancelledLimitOrdersEventV4 struct {
	Maker        common.Address
	MakerToken   common.Address
	TakerToken   common.Address
	MinValidSalt *big.Int
}
