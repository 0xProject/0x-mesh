package db

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

// Order is the database representation a 0x order along with some relevant metadata
type Order struct {
	Hash                  common.Hash    `db:"hash"`
	ChainID               *Uint256       `db:"chainID"`
	ExchangeAddress       common.Address `db:"exchangeAddress"`
	MakerAddress          common.Address `db:"makerAddress"`
	MakerAssetData        []byte         `db:"makerAssetData"`
	MakerFeeAssetData     []byte         `db:"makerFeeAssetData"`
	MakerAssetAmount      *Uint256       `db:"makerAssetAmount"`
	MakerFee              *Uint256       `db:"makerFee"`
	TakerAddress          common.Address `db:"takerAddress"`
	TakerAssetData        []byte         `db:"takerAssetData"`
	TakerFeeAssetData     []byte         `db:"takerFeeAssetData"`
	TakerAssetAmount      *Uint256       `db:"takerAssetAmount"`
	TakerFee              *Uint256       `db:"takerFee"`
	SenderAddress         common.Address `db:"senderAddress"`
	FeeRecipientAddress   common.Address `db:"feeRecipientAddress"`
	ExpirationTimeSeconds *Uint256       `db:"expirationTimeSeconds"`
	Salt                  *Uint256       `db:"salt"`
	Signature             []byte         `db:"signature"`
	// When was this order last validated
	LastUpdated time.Time `db:"lastUpdated"`
	// How much of this order can still be filled
	FillableTakerAssetAmount *Uint256 `db:"fillableTakerAssetAmount"`
	// Was this order flagged for removal? Due to the possibility of block-reorgs, instead
	// of immediately removing an order when FillableTakerAssetAmount becomes 0, we instead
	// flag it for removal. After this order isn't updated for X time and has IsRemoved = true,
	// the order can be permanently deleted.
	IsRemoved bool `db:"isRemoved"`
	// IsPinned indicates whether or not the order is pinned. Pinned orders are
	// not removed from the database unless they become unfillable.
	IsPinned bool `db:"isPinned"`
}

// Metadata is the database representation of MeshDB instance metadata
type Metadata struct {
	EthereumChainID                   int
	MaxExpirationTime                 *big.Int
	EthRPCRequestsSentInCurrentUTCDay int
	StartOfCurrentUTCDay              time.Time
}

// MiniHeader is a representation of a succinct Ethereum block headers
type MiniHeader struct {
	Hash      common.Hash
	Parent    common.Hash
	Number    *big.Int
	Timestamp time.Time
	Logs      []types.Log
}

type Uint256 struct {
	*big.Int
}

func NewUint256(v *big.Int) *Uint256 {
	return &Uint256{
		Int: v,
	}
}

func (u *Uint256) Value() (driver.Value, error) {
	if u == nil || u.Int == nil {
		return nil, nil
	}
	return u.String(), nil
}

func (u *Uint256) Scan(value interface{}) error {
	if value == nil {
		u = nil
		return nil
	}
	switch v := value.(type) {
	case int64:
		u.Int = big.NewInt(v)
	case string:
		parsed, ok := math.ParseBig256(v)
		if !ok {
			return fmt.Errorf("could not scan string value %q into Uint256", v)
		}
		u.Int = parsed
	default:
		return fmt.Errorf("could not scan type %T into Uint256", value)
	}

	return nil
}
