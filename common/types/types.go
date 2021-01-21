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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

// Stats is the return value for core.GetStats. Also used in the browser interface.
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
	MaxExpirationTime                 *big.Int    `json:"maxExpirationTime"`
	StartOfCurrentUTCDay              time.Time   `json:"startOfCurrentUTCDay"`
	EthRPCRequestsSentInCurrentUTCDay int         `json:"ethRPCRequestsSentInCurrentUTCDay"`
	EthRPCRateLimitExpiredRequests    int64       `json:"ethRPCRateLimitExpiredRequests"`
}

// LatestBlock is the latest block processed by the Mesh node.
type LatestBlock struct {
	Number *big.Int    `json:"number"`
	Hash   common.Hash `json:"hash"`
}

// GetOrdersResponse is the return value for core.GetOrders. Also used in the
// browser interface.
type GetOrdersResponse struct {
	Timestamp   time.Time    `json:"timestamp"`
	OrdersInfos []*OrderInfo `json:"ordersInfos"`
}

// AddOrdersOpts is a set of options for core.AddOrders. Also used in the
// browser interface.
type AddOrdersOpts struct {
	// Pinned determines whether or not the added orders should be pinned. Pinned
	// orders will not be affected by any DDoS prevention or incentive mechanisms
	// and will always stay in storage until they are no longer fillable. Defaults
	// to true.
	Pinned bool `json:"pinned"`
	// KeepCancelled signals that this order should not be deleted
	// even if it is cancelled.
	KeepCancelled bool `json:"keepCancelled"`
	// KeepExpired signals that this order should not be deleted
	// even if it becomes expired.
	KeepExpired bool `json:"keepExpired"`
	// KeepFullyFilled signals that this order should not be deleted
	// even if it is fully filled.
	KeepFullyFilled bool `json:"keepFullyFilled"`
	// KeepUnfunded signals that this order should not be deleted
	// even if it becomes unfunded.
	KeepUnfunded bool `json:"keepUnfunded"`
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
	OrderV3 *zeroex.Order
	OrderV4 *zeroex.OrderV4

	Hash                     common.Hash             `json:"hash"`
	Signature                []byte                  `json:"signature"`
	SignatureV4              zeroex.SignatureFieldV4 `json:"signaturev4"`
	FillableTakerAssetAmount *big.Int                `json:"fillableTakerAssetAmount"`
	LastUpdated              time.Time               `json:"lastUpdated"`
	// Was this order flagged for removal? Due to the possibility of block-reorgs, instead
	// of immediately removing an order when FillableTakerAssetAmount becomes 0, we instead
	// flag it for removal. After this order isn't updated for X time and has IsRemoved = true,
	// the order can be permanently deleted.
	IsRemoved bool `json:"isRemoved"`
	// IsPinned indicates whether or not the order is pinned. Pinned orders are
	// not removed from the database unless they become unfillable.
	IsPinned bool `json:"isPinned"`
	// IsUnfillable indicates whether or not the order has become unfillable.
	IsUnfillable bool `json:"isUnfillable"`
	// IsExpired indicates whether or not the order has become expired.
	IsExpired bool `json:"isExpired"`
	// JSON-encoded list of assetdatas contained in MakerAssetData. For non-MAP
	// orders, the list contains only one element which is equal to MakerAssetData.
	// For MAP orders, it contains each component assetdata.
	ParsedMakerAssetData []*SingleAssetData `json:"parsedMakerAssetData"`
	// Same as ParsedMakerAssetData but for MakerFeeAssetData instead of MakerAssetData.
	ParsedMakerFeeAssetData []*SingleAssetData `json:"parsedMakerFeeAssetData"`
	// LastValidatedBlockNumber is the block number at which the order was
	// last validated.
	LastValidatedBlockNumber *big.Int `json:"lastValidatedBlockNumber"`
	// LastValidatedBlockHash is the hash of the block at which the order was
	// last validated.
	LastValidatedBlockHash common.Hash `json:"lastValidatedBlockHash"`
	// KeepCancelled signals that this order should not be deleted
	// if it is cancelled.
	KeepCancelled bool `json:"keepCancelled"`
	// KeepExpired signals that this order should not be deleted
	// if it becomes expired.
	KeepExpired bool `json:"keepExpired"`
	// KeepFullyFilled signals that this order should not be deleted
	// if it is fully filled.
	KeepFullyFilled bool `json:"keepFullyFilled"`
	// KeepUnfunded signals that this order should not be deleted
	// if it becomes unfunded.
	KeepUnfunded bool `json:"keepUnfunded"`
}

func (order OrderWithMetadata) SignedOrder() *zeroex.SignedOrder {
	if order.OrderV3 == nil {
		return nil
	}
	return &zeroex.SignedOrder{
		Order: zeroex.Order{
			ChainID:               order.OrderV3.ChainID,
			ExchangeAddress:       order.OrderV3.ExchangeAddress,
			MakerAddress:          order.OrderV3.MakerAddress,
			MakerAssetData:        order.OrderV3.MakerAssetData,
			MakerFeeAssetData:     order.OrderV3.MakerFeeAssetData,
			MakerAssetAmount:      order.OrderV3.MakerAssetAmount,
			MakerFee:              order.OrderV3.MakerFee,
			TakerAddress:          order.OrderV3.TakerAddress,
			TakerAssetData:        order.OrderV3.TakerAssetData,
			TakerFeeAssetData:     order.OrderV3.TakerFeeAssetData,
			TakerAssetAmount:      order.OrderV3.TakerAssetAmount,
			TakerFee:              order.OrderV3.TakerFee,
			SenderAddress:         order.OrderV3.SenderAddress,
			FeeRecipientAddress:   order.OrderV3.FeeRecipientAddress,
			ExpirationTimeSeconds: order.OrderV3.ExpirationTimeSeconds,
			Salt:                  order.OrderV3.Salt,
		},
		Signature: order.Signature,
	}
}

func (order OrderWithMetadata) SignedOrderV4() *zeroex.SignedOrderV4 {
	if order.OrderV4 == nil {
		return nil
	}
	return &zeroex.SignedOrderV4{
		OrderV4:   *order.OrderV4, // QUESTION: Is there a reason we would explicitly copy every field like above instead of doing this?
		Signature: order.SignatureV4,
	}
}

type SingleAssetData struct {
	Address common.Address `json:"address"`
	TokenID *big.Int       `json:"tokenID"`
}

type MiniHeader struct {
	Hash      common.Hash `json:"hash"`
	Parent    common.Hash `json:"parent"`
	Number    *big.Int    `json:"number"`
	Timestamp time.Time   `json:"timestamp"`
	Logs      []types.Log `json:"logs"`
}

type Metadata struct {
	EthereumChainID                   int
	EthRPCRequestsSentInCurrentUTCDay int
	StartOfCurrentUTCDay              time.Time
}

// HexToBytes converts the the given hex string (with or without the "0x" prefix)
// to a slice of bytes. If the string is "0x" it returns nil.
func HexToBytes(s string) []byte {
	if s == "0x" {
		return nil
	}
	return common.FromHex(s)
}

// BytesToHex converts the given slice of bytes to a hex string with a "0x" prefix.
// If b is nil or has length 0, it returns "0x".
func BytesToHex(b []byte) string {
	if len(b) == 0 {
		return "0x"
	}
	return hexutil.Encode(b)
}
