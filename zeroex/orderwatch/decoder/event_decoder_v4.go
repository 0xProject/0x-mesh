package decoder

import (
	"github.com/ethereum/go-ethereum/common"
)

type ExchangeCancelEventV4 struct {
	Maker     common.Address
	OrderHash common.Hash
}
