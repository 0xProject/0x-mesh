package decoder

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

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
