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

// FIXME(jalextowle): We may have to fix these JSON marshallers.
//
// MarshalJSON is a custom Marshaler for OrderInfo
func (o OrderInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder,
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
	})
}

// FIXME(jalextowle): We may have to fix these JSON marshallers.
//
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

type OrderWithMetadataV4 struct {
	Hash                     common.Hash    `json:"hash"`
	ChainID                  *big.Int       `json:"chainID"`
	Exchange                 common.Address `json:"exchangeAddress"`
	Maker                    common.Address `json:"maker"`
	MakerToken               common.Address `json:"makerToken"`
	MakerAssetAmount         *big.Int       `json:"makerAssetAmount"`
	MakerFee                 *big.Int       `json:"makerFee"`
	Taker                    common.Address `json:"takerAddress"`
	TakerToken               common.Address `json:"takerToken"`
	TakerAssetAmount         *big.Int       `json:"takerAssetAmount"`
	TakerFee                 *big.Int       `json:"takerFee"`
	Sender                   common.Address `json:"sender"`
	FeeRecipient             common.Address `json:"feeRecipient"`
	Expiry                   *big.Int       `json:"expiry"`
	Salt                     *big.Int       `json:"salt"`
	Pool                     *big.Int       `json:"Pool"`
	Origin                   *big.Int       `json:"Origin"`
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
	// IsUnfillable indicates whether or not the order has become unfillable.
	IsUnfillable bool `json:"isUnfillable"`
	// IsExpired indicates whether or not the order has become expired.
	IsExpired bool `json:"isExpired"`
	// FIXME(jalextowle): Is this still needed?
	//
	// JSON-encoded list of assetdatas contained in MakerAssetData. For non-MAP
	// orders, the list contains only one element which is equal to MakerAssetData.
	// For MAP orders, it contains each component assetdata.
	ParsedMakerAssetData []*SingleAssetData `json:"parsedMakerAssetData"`
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

func (order OrderWithMetadataV4) SignedOrderV4() *zeroex.SignedOrder {
	return &zeroex.SignedOrder{
		Version: 4,
		Order: zeroex.OrderV4{
			ChainID:          order.ChainID,
			Exchange:         order.Exchange,
			Maker:            order.Maker,
			MakerToken:       order.MakerToken,
			MakerAssetAmount: order.MakerAssetAmount,
			MakerFee:         order.MakerFee,
			Taker:            order.Taker,
			TakerToken:       order.TakerToken,
			TakerAssetAmount: order.TakerAssetAmount,
			TakerFee:         order.TakerFee,
			Sender:           order.Sender,
			FeeRecipient:     order.FeeRecipient,
			Expiry:           order.Expiry,
			Salt:             order.Salt,
			Pool:             order.Pool,
			Origin:           order.Origin,
		},
		Signature: order.Signature,
	}
}

type OrderWithMetadataV3 struct {
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

func (order OrderWithMetadataV3) SignedOrderV3() *zeroex.SignedOrder {
	return &zeroex.SignedOrder{
		Version: 3,
		Order: zeroex.OrderV3{
			ChainID:               order.ChainID,
			ExchangeAddress:       order.ExchangeAddress,
			MakerAddress:          order.MakerAddress,
			MakerAssetData:        order.MakerAssetData,
			MakerFeeAssetData:     order.MakerFeeAssetData,
			MakerAssetAmount:      order.MakerAssetAmount,
			MakerFee:              order.MakerFee,
			TakerAddress:          order.TakerAddress,
			TakerAssetData:        order.TakerAssetData,
			TakerFeeAssetData:     order.TakerFeeAssetData,
			TakerAssetAmount:      order.TakerAssetAmount,
			TakerFee:              order.TakerFee,
			SenderAddress:         order.SenderAddress,
			FeeRecipientAddress:   order.FeeRecipientAddress,
			ExpirationTimeSeconds: order.ExpirationTimeSeconds,
			Salt:                  order.Salt,
		},
		Signature: order.Signature,
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
