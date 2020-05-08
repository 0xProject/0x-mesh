// Package types holds common types that are used across a variety of
// interfaces.
package types

import (
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

// Stats is the return value for core.GetStats. Also used in the browser and RPC
// interface.
type Stats struct {
	Version                           string      `json:"version"`
	PubSubTopic                       string      `json:"pubSubTopic"`
	Rendezvous                        string      `json:"rendezvous"`
	SecondaryRendezvous               []string    `json:"secondaryRendezvous"`
	PeerID                            string      `json:"peerID"`
	EthereumChainID                   int         `json:"ethereumChainID"`
	LatestBlock                       LatestBlock `json:"latestBlock"`
	NumPeers                          int         `json:"numPeers"`
	NumOrders                         int         `json:"numOrders"`
	NumOrdersIncludingRemoved         int         `json:"numOrdersIncludingRemoved"`
	NumPinnedOrders                   int         `json:"numPinnedOrders"`
	MaxExpirationTime                 string      `json:"maxExpirationTime"`
	StartOfCurrentUTCDay              time.Time   `json:"startOfCurrentUTCDay"`
	EthRPCRequestsSentInCurrentUTCDay int         `json:"ethRPCRequestsSentInCurrentUTCDay"`
	EthRPCRateLimitExpiredRequests    int64       `json:"ethRPCRateLimitExpiredRequests"`
}

// LatestBlock is the latest block processed by the Mesh node.
type LatestBlock struct {
	Number int         `json:"number"`
	Hash   common.Hash `json:"hash"`
}

// GetOrdersResponse is the return value for core.GetOrders. Also used in the
// browser and RPC interface.
type GetOrdersResponse struct {
	SnapshotID        string       `json:"snapshotID"`
	SnapshotTimestamp time.Time    `json:"snapshotTimestamp"`
	OrdersInfos       []*OrderInfo `json:"ordersInfos"`
}

// AddOrdersOpts is a set of options for core.AddOrders. Also used in the
// browser and RPC interface.
type AddOrdersOpts struct {
	// Pinned determines whether or not the added orders should be pinned. Pinned
	// orders will not be affected by any DDoS prevention or incentive mechanisms
	// and will always stay in storage until they are no longer fillable. Defaults
	// to true.
	Pinned bool `json:"pinned"`
}

// OrderInfo represents an fillable order and how much it could be filled for.
type OrderInfo struct {
	OrderHash                common.Hash         `json:"orderHash"`
	SignedOrder              *zeroex.SignedOrder `json:"signedOrder"`
	FillableTakerAssetAmount *big.Int            `json:"fillableTakerAssetAmount"`
}

type orderInfoJSON struct {
	OrderHash                string              `json:"orderHash"`
	SignedOrder              *zeroex.SignedOrder `json:"signedOrder"`
	FillableTakerAssetAmount string              `json:"fillableTakerAssetAmount"`
}

// MarshalJSON is a custom Marshaler for OrderInfo
func (o OrderInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder,
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
	})
}

// UnmarshalJSON implements a custom JSON unmarshaller for the OrderEvent type
func (o *OrderInfo) UnmarshalJSON(data []byte) error {
	var orderInfoJSON orderInfoJSON
	err := json.Unmarshal(data, &orderInfoJSON)
	if err != nil {
		return err
	}

	o.OrderHash = common.HexToHash(orderInfoJSON.OrderHash)
	o.SignedOrder = orderInfoJSON.SignedOrder
	var ok bool
	o.FillableTakerAssetAmount, ok = math.ParseBig256(orderInfoJSON.FillableTakerAssetAmount)
	if !ok {
		return errors.New("Invalid uint256 number encountered for FillableTakerAssetAmount")
	}
	return nil
}

type OrderWithMetadata struct {
	Hash                     common.Hash    `json:"hash"`
	ChainID                  *big.Int       `json:"chainID"`
	ExchangeAddress          common.Address `json:"exchangeAddress"`
	MakerAddress             common.Address `json:"makerAddress"`
	MakerAssetData           []byte         `json:"makerAssetData"`
	MakerFeeAssetData        []byte         `json:"makerFeeAssetData"`
	MakerAssetAmount         *big.Int       `json:"makerAssetAmount"`
	MakerFee                 *big.Int       `json:"makerFee"`
	TakerAddress             common.Address `json:"takerAddress"`
	TakerAssetData           []byte         `json:"takerAssetData"`
	TakerFeeAssetData        []byte         `json:"takerFeeAssetData"`
	TakerAssetAmount         *big.Int       `json:"takerAssetAmount"`
	TakerFee                 *big.Int       `json:"takerFee"`
	SenderAddress            common.Address `json:"senderAddress"`
	FeeRecipientAddress      common.Address `json:"feeRecipientAddress"`
	ExpirationTimeSeconds    *big.Int       `json:"expirationTimeSeconds"`
	Salt                     *big.Int       `json:"salt"`
	Signature                []byte         `json:"signature"`
	FillableTakerAssetAmount *big.Int       `json:"fillableTakerAssetAmount"`
	LastUpdated              time.Time      `json:"lastUpdated"`
	// Was this order flagged for removal? Due to the possibility of block-reorgs, instead
	// of immediately removing an order when FillableTakerAssetAmount becomes 0, we instead
	// flag it for removal. After this order isn't updated for X time and has IsRemoved = true,
	// the order can be permanently deleted.
	IsRemoved bool `json:"isRemoved"`
	// IsPinned indicates whether or not the order is pinned. Pinned orders are
	// not removed from the database unless they become unfillable.
	IsPinned bool `json:"isPinned"`
	// JSON-encoded list of assetdatas contained in MakerAssetData. For non-MAP
	// orders, the list contains only one element which is equal to MakerAssetData.
	// For MAP orders, it contains each component assetdata.
	ParsedMakerAssetData []*SingleAssetData `json:"parsedMakerAssetData"`
	// Same as ParsedMakerAssetData but for MakerFeeAssetData instead of MakerAssetData.
	ParsedMakerFeeAssetData []*SingleAssetData `json:"parsedMakerFeeAssetData"`
}

type SingleAssetData struct {
	Address common.Address `json:"address"`
	TokenID *big.Int       `json:"tokenID"`
}

type MiniHeader struct {
	Hash      common.Hash  `json:"hash"`
	Parent    common.Hash  `json:"parent"`
	Number    *big.Int     `json:"number"`
	Timestamp time.Time    `json:"timestamp"`
	Logs      []*types.Log `json:"logs"`
}
