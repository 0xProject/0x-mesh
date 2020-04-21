package meshdb

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Order is the database representation a 0x order along with some relevant metadata
type Order struct {
	Hash                  common.Hash
	ChainID               *big.Int
	ExchangeAddress       common.Address
	MakerAddress          common.Address
	MakerAssetData        []byte
	MakerFeeAssetData     []byte
	MakerAssetAmount      *big.Int
	MakerFee              *big.Int
	TakerAddress          common.Address
	TakerAssetData        []byte
	TakerFeeAssetData     []byte
	TakerAssetAmount      *big.Int
	TakerFee              *big.Int
	SenderAddress         common.Address
	FeeRecipientAddress   common.Address
	ExpirationTimeSeconds *big.Int
	Salt                  *big.Int
	Signature             []byte
	// When was this order last validated
	LastUpdated time.Time
	// How much of this order can still be filled
	FillableTakerAssetAmount *big.Int
	// Was this order flagged for removal? Due to the possibility of block-reorgs, instead
	// of immediately removing an order when FillableTakerAssetAmount becomes 0, we instead
	// flag it for removal. After this order isn't updated for X time and has IsRemoved = true,
	// the order can be permanently deleted.
	IsRemoved bool
	// IsPinned indicates whether or not the order is pinned. Pinned orders are
	// not removed from the database unless they become unfillable.
	IsPinned bool
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
