package client

import (
	"time"

	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

func addOrdersResultsFromGQLType(results *gqltypes.AddOrdersResults) *AddOrdersResults {
	return &AddOrdersResults{
		Accepted: acceptedOrderResultsFromGQLType(results.Accepted),
		Rejected: rejectedOrderResultsFromGQLType(results.Rejected),
	}
}

func acceptedOrderResultFromGQLType(result *gqltypes.AcceptedOrderResult) *AcceptedOrderResult {
	return &AcceptedOrderResult{
		Order: orderWithMetadataFromGQLType(result.Order),
		IsNew: result.IsNew,
	}
}

func acceptedOrderResultsFromGQLType(results []*gqltypes.AcceptedOrderResult) []*AcceptedOrderResult {
	result := make([]*AcceptedOrderResult, len(results))
	for i, r := range results {
		result[i] = acceptedOrderResultFromGQLType(r)
	}
	return result
}

func rejectedOrderResultFromGQLType(result *gqltypes.RejectedOrderResult) *RejectedOrderResult {
	var hash *common.Hash
	if result.Hash != nil {
		h := common.HexToHash(*result.Hash)
		hash = &h
	}
	return &RejectedOrderResult{
		Hash: hash,
		Order: &Order{
			ChainID:               math.MustParseBig256(result.Order.ChainID),
			ExchangeAddress:       common.HexToAddress(result.Order.ExchangeAddress),
			MakerAddress:          common.HexToAddress(result.Order.MakerAddress),
			MakerAssetData:        common.FromHex(result.Order.MakerAssetData),
			MakerFeeAssetData:     common.FromHex(result.Order.MakerFeeAssetData),
			MakerAssetAmount:      math.MustParseBig256(result.Order.MakerAssetAmount),
			MakerFee:              math.MustParseBig256(result.Order.MakerFee),
			TakerAddress:          common.HexToAddress(result.Order.TakerAddress),
			TakerAssetData:        common.FromHex(result.Order.TakerAssetData),
			TakerFeeAssetData:     common.FromHex(result.Order.TakerFeeAssetData),
			TakerAssetAmount:      math.MustParseBig256(result.Order.TakerAssetAmount),
			TakerFee:              math.MustParseBig256(result.Order.TakerFee),
			SenderAddress:         common.HexToAddress(result.Order.SenderAddress),
			FeeRecipientAddress:   common.HexToAddress(result.Order.FeeRecipientAddress),
			ExpirationTimeSeconds: math.MustParseBig256(result.Order.ExpirationTimeSeconds),
			Salt:                  math.MustParseBig256(result.Order.Salt),
			Signature:             common.FromHex(result.Order.Signature),
		},
		Code:    result.Code,
		Message: result.Message,
	}
}

func rejectedOrderResultsFromGQLType(results []*gqltypes.RejectedOrderResult) []*RejectedOrderResult {
	result := make([]*RejectedOrderResult, len(results))
	for i, r := range results {
		result[i] = rejectedOrderResultFromGQLType(r)
	}
	return result
}

func orderWithMetadataFromGQLType(order *gqltypes.OrderWithMetadataV3) *OrderWithMetadata {
	return &OrderWithMetadata{
		Hash:                     common.HexToHash(order.Hash),
		ChainID:                  math.MustParseBig256(order.ChainID),
		ExchangeAddress:          common.HexToAddress(order.ExchangeAddress),
		MakerAddress:             common.HexToAddress(order.MakerAddress),
		MakerAssetData:           common.FromHex(order.MakerAssetData),
		MakerFeeAssetData:        common.FromHex(order.MakerFeeAssetData),
		MakerAssetAmount:         math.MustParseBig256(order.MakerAssetAmount),
		MakerFee:                 math.MustParseBig256(order.MakerFee),
		TakerAddress:             common.HexToAddress(order.TakerAddress),
		TakerAssetData:           common.FromHex(order.TakerAssetData),
		TakerFeeAssetData:        common.FromHex(order.TakerFeeAssetData),
		TakerAssetAmount:         math.MustParseBig256(order.TakerAssetAmount),
		TakerFee:                 math.MustParseBig256(order.TakerFee),
		SenderAddress:            common.HexToAddress(order.SenderAddress),
		FeeRecipientAddress:      common.HexToAddress(order.FeeRecipientAddress),
		ExpirationTimeSeconds:    math.MustParseBig256(order.ExpirationTimeSeconds),
		Salt:                     math.MustParseBig256(order.Salt),
		Signature:                common.FromHex(order.Signature),
		FillableTakerAssetAmount: math.MustParseBig256(order.FillableTakerAssetAmount),
	}
}

func ordersWithMetadataFromGQLType(orders []*gqltypes.OrderWithMetadataV3) []*OrderWithMetadata {
	result := make([]*OrderWithMetadata, len(orders))
	for i, r := range orders {
		result[i] = orderWithMetadataFromGQLType(r)
	}
	return result
}

func statsFromGQLType(stats *gqltypes.Stats) (*Stats, error) {
	startOfCurrentUTCDay, err := time.Parse(time.RFC3339, stats.StartOfCurrentUTCDay)
	if err != nil {
		return nil, err
	}
	return &Stats{
		Version:                           stats.Version,
		PubSubTopic:                       stats.PubSubTopic,
		Rendezvous:                        stats.Rendezvous,
		PeerID:                            stats.PeerID,
		EthereumChainID:                   stats.EthereumChainID,
		LatestBlock:                       latestBlockFromGQLType(stats.LatestBlock),
		NumPeers:                          stats.NumPeers,
		NumOrders:                         stats.NumOrders,
		NumOrdersIncludingRemoved:         stats.NumOrdersIncludingRemoved,
		StartOfCurrentUTCDay:              startOfCurrentUTCDay,
		EthRPCRequestsSentInCurrentUTCDay: stats.EthRPCRequestsSentInCurrentUTCDay,
		EthRPCRateLimitExpiredRequests:    stats.EthRPCRateLimitExpiredRequests,
		MaxExpirationTime:                 math.MustParseBig256(stats.MaxExpirationTime),
	}, nil
}

func latestBlockFromGQLType(latestBlock *gqltypes.LatestBlock) *LatestBlock {
	return &LatestBlock{
		Number: math.MustParseBig256(latestBlock.Number),
		Hash:   common.HexToHash(latestBlock.Hash),
	}
}
