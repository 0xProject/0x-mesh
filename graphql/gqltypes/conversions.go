package gqltypes

import (
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

func StatsFromCommonType(stats *types.Stats) *Stats {
	return &Stats{
		Version:     stats.Version,
		PubSubTopic: stats.PubSubTopic,
		Rendezvous:  stats.Rendezvous,
		PeerID:      stats.PeerID,
		// TODO(albrow): This should be a big.Int in core package.
		// EthereumChainID:                   stats.EthereumChainID,
		// TODO(albrow): LatestBlock should be a pointer in core package.
		LatestBlock:                       LatestBlockFromCommonType(stats.LatestBlock),
		NumPeers:                          stats.NumPeers,
		NumOrders:                         stats.NumOrders,
		NumOrdersIncludingRemoved:         stats.NumOrdersIncludingRemoved,
		StartOfCurrentUTCDay:              stats.StartOfCurrentUTCDay,
		EthRPCRequestsSentInCurrentUTCDay: stats.EthRPCRequestsSentInCurrentUTCDay,
		EthRPCRateLimitExpiredRequests:    int(stats.EthRPCRateLimitExpiredRequests),
		MaxExpirationTime:                 BigNumber(*stats.MaxExpirationTime),
	}
}

func LatestBlockFromCommonType(latestBlock types.LatestBlock) *LatestBlock {
	return &LatestBlock{
		Number: BigNumber(*latestBlock.Number),
		Hash:   Hash(latestBlock.Hash),
	}
}

func OrderWithMetadataFromCommonType(order *types.OrderWithMetadata) *OrderWithMetadata {
	return &OrderWithMetadata{
		Hash:                     Hash(order.Hash),
		ChainID:                  BigNumber(*order.ChainID),
		ExchangeAddress:          Address(order.ExchangeAddress),
		MakerAddress:             Address(order.MakerAddress),
		MakerAssetData:           Bytes(order.MakerAssetData),
		MakerFeeAssetData:        Bytes(order.MakerFeeAssetData),
		MakerAssetAmount:         BigNumber(*order.MakerAssetAmount),
		MakerFee:                 BigNumber(*order.MakerFee),
		TakerAddress:             Address(order.TakerAddress),
		TakerAssetData:           Bytes(order.TakerAssetData),
		TakerFeeAssetData:        Bytes(order.TakerFeeAssetData),
		TakerAssetAmount:         BigNumber(*order.TakerAssetAmount),
		TakerFee:                 BigNumber(*order.TakerFee),
		SenderAddress:            Address(order.SenderAddress),
		FeeRecipientAddress:      Address(order.FeeRecipientAddress),
		ExpirationTimeSeconds:    BigNumber(*order.ExpirationTimeSeconds),
		Salt:                     BigNumber(*order.Salt),
		Signature:                Bytes(order.Signature),
		FillableTakerAssetAmount: BigNumber(*order.FillableTakerAssetAmount),
	}
}

func OrdersWithMetadataFromCommonType(orders []*types.OrderWithMetadata) []*OrderWithMetadata {
	result := make([]*OrderWithMetadata, len(orders))
	for i, order := range orders {
		result[i] = OrderWithMetadataFromCommonType(order)
	}
	return result
}

func FilterKindToDBType(kind FilterKind) (db.FilterKind, error) {
	switch kind {
	case FilterKindEqual:
		return db.Equal, nil
	case FilterKindNotEqual:
		return db.NotEqual, nil
	case FilterKindGreater:
		return db.Greater, nil
	case FilterKindGreaterOrEqual:
		return db.GreaterOrEqual, nil
	case FilterKindLess:
		return db.Less, nil
	case FilterKindLessOrEqual:
		return db.LessOrEqual, nil
	default:
		return "", fmt.Errorf("invalid filter kind: %q", kind)
	}
}

func ConvertFilterValue(f *OrderFilter) (interface{}, error) {
	switch f.Field {
	case "chainID", "makerAssetAmount", "makerFee", "takerAssetAmount", "takerFee", "expirationTimeSeconds", "salt", "fillableTakerAssetAmount":
		return stringToBigInt(f.Value)
	case "hash":
		return stringToHash(f.Value)
	case "exchangeAddress", "makerAddress", "takerAddress", "senderAddress", "feeRecipientAddress":
		return stringToAddress(f.Value)
	case "makerAssetData", "makerFeeAssetData", "takerAssetData", "takerFeeAssetData":
		return stringToBytes(f.Value)
	default:
		return "", fmt.Errorf("invalid filter field: %q", f.Field)
	}
}

func filterValueToString(value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("invalid type for filter value (expected string but got %T)", value)
	}
	return valueString, nil
}

func stringToBigInt(value interface{}) (*big.Int, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return nil, err
	}
	result, valid := math.ParseBig256(valueString)
	if !valid {
		return nil, fmt.Errorf("could not convert %q to *big.Int", value)
	}
	return result, nil
}

func stringToHash(value interface{}) (common.Hash, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(valueString), nil
}

func stringToAddress(value interface{}) (common.Address, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(valueString), nil
}

func stringToBytes(value interface{}) ([]byte, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return nil, err
	}
	return common.FromHex(valueString), nil
}

func SortDirectionToDBType(direction SortDirection) (db.SortDirection, error) {
	switch direction {
	case SortDirectionAsc:
		return db.Ascending, nil
	case SortDirectionDesc:
		return db.Descending, nil
	default:
		return "", fmt.Errorf("invalid sort direction: %q", direction)
	}
}
