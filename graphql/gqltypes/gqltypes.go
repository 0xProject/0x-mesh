package gqltypes

import (
	"time"

	"github.com/0xProject/0x-mesh/common/types"
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
