package gqltypes

import (
	"fmt"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/ethereum/go-ethereum/common"
)

type LatestBlock struct {
	Number string
	Hash   string
}

func LatestBlockFromCommonType(latestBlock types.LatestBlock) LatestBlock {
	return LatestBlock{
		// TODO(albrow): Change this.
		Number: fmt.Sprint(latestBlock.Number),
		Hash:   latestBlock.Hash.Hex(),
	}
}

type Stats struct {
	Version                           string
	PubSubTopic                       string
	Rendezvous                        string
	PeerID                            string
	EthereumChainID                   int32
	LatestBlock                       LatestBlock
	NumPeers                          int32
	NumOrders                         int32
	NumOrdersIncludingRemoved         int32
	StartOfCurrentUTCDay              string
	EthRPCRequestsSentInCurrentUTCDay int32
	EthRPCRateLimitExpiredRequests    int32
	MaxExpirationTime                 string
}

func StatsFromCommonType(stats *types.Stats) *Stats {
	return &Stats{
		Version:                           stats.Version,
		PubSubTopic:                       stats.PubSubTopic,
		Rendezvous:                        stats.Rendezvous,
		PeerID:                            stats.PeerID,
		EthereumChainID:                   int32(stats.EthereumChainID),
		LatestBlock:                       LatestBlockFromCommonType(stats.LatestBlock),
		NumPeers:                          int32(stats.NumPeers),
		NumOrders:                         int32(stats.NumOrders),
		NumOrdersIncludingRemoved:         int32(stats.NumOrdersIncludingRemoved),
		StartOfCurrentUTCDay:              stats.StartOfCurrentUTCDay.Format(time.RFC3339),
		EthRPCRequestsSentInCurrentUTCDay: int32(stats.EthRPCRequestsSentInCurrentUTCDay),
		EthRPCRateLimitExpiredRequests:    int32(stats.EthRPCRateLimitExpiredRequests),
		MaxExpirationTime:                 stats.MaxExpirationTime.String(),
	}
}

type OrderWithMetadata struct {
	Hash                     string
	ChainID                  string
	ExchangeAddress          string
	MakerAddress             string
	MakerAssetData           string
	MakerFeeAssetData        string
	MakerAssetAmount         string
	MakerFee                 string
	TakerAddress             string
	TakerAssetData           string
	TakerFeeAssetData        string
	TakerAssetAmount         string
	TakerFee                 string
	SenderAddress            string
	FeeRecipientAddress      string
	ExpirationTimeSeconds    string
	Salt                     string
	Signature                string
	FillableTakerAssetAmount string
	LastUpdated              string
}

func OrderWithMetadataFromCommonType(order *types.OrderWithMetadata) *OrderWithMetadata {
	return &OrderWithMetadata{
		Hash:                     order.Hash.Hex(),
		ChainID:                  order.ChainID.String(),
		ExchangeAddress:          strings.ToLower(order.ExchangeAddress.Hex()),
		MakerAddress:             strings.ToLower(order.MakerAddress.Hex()),
		MakerAssetData:           common.Bytes2Hex(order.MakerAssetData),
		MakerFeeAssetData:        common.Bytes2Hex(order.MakerFeeAssetData),
		MakerAssetAmount:         order.MakerAssetAmount.String(),
		MakerFee:                 order.MakerFee.String(),
		TakerAddress:             strings.ToLower(order.TakerAddress.Hex()),
		TakerAssetData:           common.Bytes2Hex(order.TakerAssetData),
		TakerFeeAssetData:        common.Bytes2Hex(order.TakerFeeAssetData),
		TakerAssetAmount:         order.TakerAssetAmount.String(),
		TakerFee:                 order.TakerFee.String(),
		SenderAddress:            strings.ToLower(order.SenderAddress.Hex()),
		FeeRecipientAddress:      strings.ToLower(order.FeeRecipientAddress.Hex()),
		ExpirationTimeSeconds:    order.ExpirationTimeSeconds.String(),
		Salt:                     order.Salt.String(),
		Signature:                common.Bytes2Hex(order.Signature),
		FillableTakerAssetAmount: order.FillableTakerAssetAmount.String(),
		LastUpdated:              order.LastUpdated.Format(time.RFC3339),
	}
}

func OrdersWithMetadataFromCommonType(orders []*types.OrderWithMetadata) []*OrderWithMetadata {
	result := make([]*OrderWithMetadata, len(orders))
	for i, order := range orders {
		result[i] = OrderWithMetadataFromCommonType(order)
	}
	return result
}

func FilterKindToDBType(kind string) (db.FilterKind, error) {
	switch kind {
	case "EQUAL":
		return db.Equal, nil
	case "NOT_EQUAL":
		return db.NotEqual, nil
	case "GREATER":
		return db.Greater, nil
	case "GREATER_OR_EQUAL":
		return db.GreaterOrEqual, nil
	case "LESS":
		return db.Less, nil
	case "LESS_OR_EQUAL":
		return db.LessOrEqual, nil
	default:
		return "", fmt.Errorf("invalid filter kind: %q", kind)
	}
}

func SortDirectionToDBType(direction string) (db.SortDirection, error) {
	switch direction {
	case "ASC":
		return db.Ascending, nil
	case "DESC":
		return db.Descending, nil
	default:
		return "", fmt.Errorf("invalid sort direction: %q", direction)
	}
}
