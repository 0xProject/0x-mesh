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
)

// Stats is the return value for core.GetStats. Also used in the browser and RPC
// interface.
type Stats struct {
	Version                           string      `json:"version"`
	PubSubTopic                       string      `json:"pubSubTopic"`
	Rendezvous                        string      `json:"rendezvous"`
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
