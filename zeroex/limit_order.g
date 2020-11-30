package zeroex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ERC20LimitOrder struct {
	makerToken          common.Address
	takerToken          common.Address
	makerAmount         *big.Int
	takerAmount         *big.Int
	feeRecipient        common.Address
	takerTokenFeeAmount *big.Int
	maker               common.Address
	taker               common.Address
	sender              common.Address
	pool                []byte
	expiry              *big.Int
	salt                *big.Int
}

type SignedERC20LimitOrder struct {
	ERC20LimitOrder
	Signature []byte
}
