package gqltypes

import (
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
)

type LatestBlock struct {
	Number string
	Hash   string
}

func LatestBlockFromCommonType(latestBlock types.LatestBlock) LatestBlock {
	return LatestBlock{
		Number: latestBlock.Number,
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
		MaxExpirationTime:                 stats.MaxExpirationTime,
	}
}

type OrderWithMetadata struct {
	Hash                              string
	ChainID                           string
	ExchangeAddress                   string
	MakerAddress                      string
	MakerAssetData                    string
	MakerFeeAssetData                 string
	MakerAssetAmount                  string
	MakerFee                          string
	TakerAddress                      string
	TakerAssetData                    string
	TakerFeeAssetData                 string
	TakerAssetAmount                  string
	TakerFee                          string
	SenderAddress                     string
	FeeRecipientAddress               string
	ExpirationTimeSeconds             string
	Salt                              string
	Signature                         string
	RemainingFillableTakerAssetAmount string
	LastUpdated                       string
}

func OrderWithMetadataFromCommonType(order *types.OrderWithMetadata) *OrderWithMetadata {
	return &OrderWithMetadata{
		Hash:                              order.Hash.Hex(),
		ChainID:                           order.ChainID.String(),
		ExchangeAddress:                   order.ExchangeAddress.Hex(),
		MakerAddress:                      order.MakerAddress.Hex(),
		MakerAssetData:                    common.Bytes2Hex(order.MakerAssetData),
		MakerFeeAssetData:                 common.Bytes2Hex(order.MakerFeeAssetData),
		MakerAssetAmount:                  order.MakerAssetAmount.String(),
		MakerFee:                          order.MakerFee.String(),
		TakerAddress:                      order.TakerAddress.Hex(),
		TakerAssetData:                    common.Bytes2Hex(order.TakerAssetData),
		TakerFeeAssetData:                 common.Bytes2Hex(order.TakerFeeAssetData),
		TakerAssetAmount:                  order.TakerAssetAmount.String(),
		TakerFee:                          order.TakerFee.String(),
		SenderAddress:                     order.SenderAddress.Hex(),
		FeeRecipientAddress:               order.FeeRecipientAddress.Hex(),
		ExpirationTimeSeconds:             order.ExpirationTimeSeconds.String(),
		Salt:                              order.Salt.String(),
		Signature:                         common.Bytes2Hex(order.Signature),
		RemainingFillableTakerAssetAmount: order.FillableTakerAssetAmount.String(),
		LastUpdated:                       order.LastUpdated.Format(time.RFC3339),
	}
}
