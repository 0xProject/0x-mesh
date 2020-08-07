package gqltypes

import (
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
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

func BigNumberToBigInt(bigNumber BigNumber) *big.Int {
	bigInt := big.Int(bigNumber)
	return &bigInt
}

func NewOrderToSignedOrder(newOrder *NewOrder) *zeroex.SignedOrder {
	return &zeroex.SignedOrder{
		Order: zeroex.Order{
			ChainID:               BigNumberToBigInt(newOrder.ChainID),
			ExchangeAddress:       common.Address(newOrder.ExchangeAddress),
			MakerAddress:          common.Address(newOrder.MakerAddress),
			MakerAssetData:        newOrder.MakerAssetData,
			MakerFeeAssetData:     newOrder.MakerFeeAssetData,
			MakerAssetAmount:      BigNumberToBigInt(newOrder.MakerAssetAmount),
			MakerFee:              BigNumberToBigInt(newOrder.MakerFee),
			TakerAddress:          common.Address(newOrder.TakerAddress),
			TakerAssetData:        newOrder.TakerAssetData,
			TakerFeeAssetData:     newOrder.TakerFeeAssetData,
			TakerAssetAmount:      BigNumberToBigInt(newOrder.TakerAssetAmount),
			TakerFee:              BigNumberToBigInt(newOrder.TakerFee),
			SenderAddress:         common.Address(newOrder.SenderAddress),
			FeeRecipientAddress:   common.Address(newOrder.FeeRecipientAddress),
			ExpirationTimeSeconds: BigNumberToBigInt(newOrder.ExpirationTimeSeconds),
			Salt:                  BigNumberToBigInt(newOrder.Salt),
		},
	}
}

func NewOrdersToSignedOrders(newOrders []*NewOrder) []*zeroex.SignedOrder {
	result := make([]*zeroex.SignedOrder, len(newOrders))
	for i, newOrder := range newOrders {
		result[i] = NewOrderToSignedOrder(newOrder)
	}
	return result
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

func AddOrdersResultsFromValidationResults(validationResults *ordervalidator.ValidationResults) (*AddOrdersResults, error) {
	rejected, err := RejectedOrderResultsFromOrderInfos(validationResults.Rejected)
	if err != nil {
		return nil, err
	}
	return &AddOrdersResults{
		Accepted: AcceptedOrderResultsFromOrderInfos(validationResults.Accepted),
		Rejected: rejected,
	}, nil
}

func AcceptedOrderResultFromOrderInfo(info *ordervalidator.AcceptedOrderInfo) *AcceptedOrderResult {
	return &AcceptedOrderResult{
		Order: &OrderWithMetadata{
			Hash:                     Hash(info.OrderHash),
			ChainID:                  BigNumber(*info.SignedOrder.ChainID),
			ExchangeAddress:          Address(info.SignedOrder.ExchangeAddress),
			MakerAddress:             Address(info.SignedOrder.MakerAddress),
			MakerAssetData:           Bytes(info.SignedOrder.MakerAssetData),
			MakerFeeAssetData:        Bytes(info.SignedOrder.MakerFeeAssetData),
			MakerAssetAmount:         BigNumber(*info.SignedOrder.MakerAssetAmount),
			MakerFee:                 BigNumber(*info.SignedOrder.MakerFee),
			TakerAddress:             Address(info.SignedOrder.TakerAddress),
			TakerAssetData:           Bytes(info.SignedOrder.TakerAssetData),
			TakerFeeAssetData:        Bytes(info.SignedOrder.TakerFeeAssetData),
			TakerAssetAmount:         BigNumber(*info.SignedOrder.TakerAssetAmount),
			TakerFee:                 BigNumber(*info.SignedOrder.TakerFee),
			SenderAddress:            Address(info.SignedOrder.SenderAddress),
			FeeRecipientAddress:      Address(info.SignedOrder.FeeRecipientAddress),
			ExpirationTimeSeconds:    BigNumber(*info.SignedOrder.ExpirationTimeSeconds),
			Salt:                     BigNumber(*info.SignedOrder.Salt),
			Signature:                Bytes(info.SignedOrder.Signature),
			FillableTakerAssetAmount: BigNumber(*info.FillableTakerAssetAmount),
		},
		IsNew: info.IsNew,
	}
}

func AcceptedOrderResultsFromOrderInfos(infos []*ordervalidator.AcceptedOrderInfo) []*AcceptedOrderResult {
	result := make([]*AcceptedOrderResult, len(infos))
	for i, info := range infos {
		result[i] = AcceptedOrderResultFromOrderInfo(info)
	}
	return result
}

func RejectedOrderResultFromOrderInfo(info *ordervalidator.RejectedOrderInfo) (*RejectedOrderResult, error) {
	var hash *Hash
	if info.OrderHash.Hex() != "0x" {
		gqlHash := Hash(info.OrderHash)
		hash = &gqlHash
	}
	code, err := RejectedCodeFromValidatorStatus(info.Status)
	if err != nil {
		return nil, err
	}
	return &RejectedOrderResult{
		Hash: hash,
		Order: &Order{
			ChainID:               BigNumber(*info.SignedOrder.ChainID),
			ExchangeAddress:       Address(info.SignedOrder.ExchangeAddress),
			MakerAddress:          Address(info.SignedOrder.MakerAddress),
			MakerAssetData:        Bytes(info.SignedOrder.MakerAssetData),
			MakerFeeAssetData:     Bytes(info.SignedOrder.MakerFeeAssetData),
			MakerAssetAmount:      BigNumber(*info.SignedOrder.MakerAssetAmount),
			MakerFee:              BigNumber(*info.SignedOrder.MakerFee),
			TakerAddress:          Address(info.SignedOrder.TakerAddress),
			TakerAssetData:        Bytes(info.SignedOrder.TakerAssetData),
			TakerFeeAssetData:     Bytes(info.SignedOrder.TakerFeeAssetData),
			TakerAssetAmount:      BigNumber(*info.SignedOrder.TakerAssetAmount),
			TakerFee:              BigNumber(*info.SignedOrder.TakerFee),
			SenderAddress:         Address(info.SignedOrder.SenderAddress),
			FeeRecipientAddress:   Address(info.SignedOrder.FeeRecipientAddress),
			ExpirationTimeSeconds: BigNumber(*info.SignedOrder.ExpirationTimeSeconds),
			Salt:                  BigNumber(*info.SignedOrder.Salt),
			Signature:             Bytes(info.SignedOrder.Signature),
		},
		Code:    code,
		Message: info.Status.Message,
	}, nil
}

func RejectedOrderResultsFromOrderInfos(infos []*ordervalidator.RejectedOrderInfo) ([]*RejectedOrderResult, error) {
	result := make([]*RejectedOrderResult, len(infos))
	for i, info := range infos {
		rejectedResult, err := RejectedOrderResultFromOrderInfo(info)
		if err != nil {
			return nil, err
		}
		result[i] = rejectedResult
	}
	return result, nil
}

func RejectedCodeFromValidatorStatus(status ordervalidator.RejectedOrderStatus) (RejectedOrderCode, error) {
	switch status.Code {
	case ordervalidator.ROEthRPCRequestFailed.Code:
		return RejectedOrderCodeEthRPCRequestFailed, nil
	case ordervalidator.ROCoordinatorRequestFailed.Code:
		return RejectedOrderCodeCoordinatorRequestFailed, nil
	case ordervalidator.ROCoordinatorSoftCancelled.Code:
		return RejectedOrderCodeCoordinatorSoftCancelled, nil
	case ordervalidator.ROCoordinatorEndpointNotFound.Code:
		return RejectedOrderCodeCoordinatorEndpointNotFound, nil
	case ordervalidator.ROInvalidMakerAssetAmount.Code:
		return RejectedOrderCodeOrderHasInvalidMakerAssetAmount, nil
	case ordervalidator.ROInvalidTakerAssetAmount.Code:
		return RejectedOrderCodeOrderHasInvalidTakerAssetAmount, nil
	case ordervalidator.ROExpired.Code:
		return RejectedOrderCodeOrderExpired, nil
	case ordervalidator.ROFullyFilled.Code:
		return RejectedOrderCodeOrderFullyFilled, nil
	case ordervalidator.ROCancelled.Code:
		return RejectedOrderCodeOrderCancelled, nil
	case ordervalidator.ROUnfunded.Code:
		return RejectedOrderCodeOrderUnfunded, nil
	case ordervalidator.ROInvalidMakerAssetData.Code:
		return RejectedOrderCodeOrderHasInvalidMakerAssetData, nil
	case ordervalidator.ROInvalidMakerFeeAssetData.Code:
		return RejectedOrderCodeOrderHasInvalidMakerFeeAssetData, nil
	case ordervalidator.ROInvalidTakerAssetData.Code:
		return RejectedOrderCodeOrderHasInvalidTakerAssetData, nil
	case ordervalidator.ROInvalidTakerFeeAssetData.Code:
		return RejectedOrderCodeOrderHasInvalidTakerFeeAssetData, nil
	case ordervalidator.ROInvalidSignature.Code:
		return RejectedOrderCodeOrderHasInvalidSignature, nil
	case ordervalidator.ROMaxExpirationExceeded.Code:
		return RejectedOrderCodeOrderMaxExpirationExceeded, nil
	case ordervalidator.ROInternalError.Code:
		return RejectedOrderCodeInternalError, nil
	case ordervalidator.ROMaxOrderSizeExceeded.Code:
		return RejectedOrderCodeMaxOrderSizeExceeded, nil
	case ordervalidator.ROOrderAlreadyStoredAndUnfillable.Code:
		return RejectedOrderCodeOrderAlreadyStoredAndUnfillable, nil
	case ordervalidator.ROIncorrectChain.Code:
		return RejectedOrderCodeOrderForIncorrectChain, nil
	case ordervalidator.ROIncorrectExchangeAddress.Code:
		return RejectedOrderCodeIncorrectExchangeAddress, nil
	case ordervalidator.ROSenderAddressNotAllowed.Code:
		return RejectedOrderCodeSenderAddressNotAllowed, nil
	case ordervalidator.RODatabaseFullOfOrders.Code:
		return RejectedOrderCodeDatabaseFullOfOrders, nil
	default:
		return "", fmt.Errorf("unexpected RejectedOrderStatus.Code: %q", status.Code)
	}
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
