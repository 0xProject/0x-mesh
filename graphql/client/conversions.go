package client

import (
	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/ethereum/go-ethereum/common"
)

func addOrdersResultsFromGQLType(results *gqltypes.AddOrdersResults) *AddOrdersResults {
	return &AddOrdersResults{
		Accepted: acceptedOrderResultsFromGQLType(results.Accepted),
		Rejected: rejectedOrderResultsFromGQLType(results.Rejected),
	}
}

func acceptedOrderResultFromGQLType(result *gqltypes.AcceptedOrderResult) *AcceptedOrderResult {
	return &AcceptedOrderResult{
		Order: &OrderWithMetadata{
			Hash:                     common.Hash(result.Order.Hash),
			ChainID:                  gqltypes.BigNumberToBigInt(result.Order.ChainID),
			ExchangeAddress:          common.Address(result.Order.ExchangeAddress),
			MakerAddress:             common.Address(result.Order.MakerAddress),
			MakerAssetData:           result.Order.MakerAssetData,
			MakerFeeAssetData:        result.Order.MakerFeeAssetData,
			MakerAssetAmount:         gqltypes.BigNumberToBigInt(result.Order.MakerAssetAmount),
			MakerFee:                 gqltypes.BigNumberToBigInt(result.Order.MakerFee),
			TakerAddress:             common.Address(result.Order.TakerAddress),
			TakerAssetData:           result.Order.TakerAssetData,
			TakerFeeAssetData:        result.Order.TakerFeeAssetData,
			TakerAssetAmount:         gqltypes.BigNumberToBigInt(result.Order.TakerAssetAmount),
			TakerFee:                 gqltypes.BigNumberToBigInt(result.Order.TakerFee),
			SenderAddress:            common.Address(result.Order.SenderAddress),
			FeeRecipientAddress:      common.Address(result.Order.FeeRecipientAddress),
			ExpirationTimeSeconds:    gqltypes.BigNumberToBigInt(result.Order.ExpirationTimeSeconds),
			Salt:                     gqltypes.BigNumberToBigInt(result.Order.Salt),
			Signature:                result.Order.Signature,
			FillableTakerAssetAmount: gqltypes.BigNumberToBigInt(result.Order.FillableTakerAssetAmount),
		},
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
		h := common.Hash(*result.Hash)
		hash = &h
	}
	return &RejectedOrderResult{
		Hash: hash,
		Order: &Order{
			ChainID:               gqltypes.BigNumberToBigInt(result.Order.ChainID),
			ExchangeAddress:       common.Address(result.Order.ExchangeAddress),
			MakerAddress:          common.Address(result.Order.MakerAddress),
			MakerAssetData:        result.Order.MakerAssetData,
			MakerFeeAssetData:     result.Order.MakerFeeAssetData,
			MakerAssetAmount:      gqltypes.BigNumberToBigInt(result.Order.MakerAssetAmount),
			MakerFee:              gqltypes.BigNumberToBigInt(result.Order.MakerFee),
			TakerAddress:          common.Address(result.Order.TakerAddress),
			TakerAssetData:        result.Order.TakerAssetData,
			TakerFeeAssetData:     result.Order.TakerFeeAssetData,
			TakerAssetAmount:      gqltypes.BigNumberToBigInt(result.Order.TakerAssetAmount),
			TakerFee:              gqltypes.BigNumberToBigInt(result.Order.TakerFee),
			SenderAddress:         common.Address(result.Order.SenderAddress),
			FeeRecipientAddress:   common.Address(result.Order.FeeRecipientAddress),
			ExpirationTimeSeconds: gqltypes.BigNumberToBigInt(result.Order.ExpirationTimeSeconds),
			Salt:                  gqltypes.BigNumberToBigInt(result.Order.Salt),
			Signature:             result.Order.Signature,
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
