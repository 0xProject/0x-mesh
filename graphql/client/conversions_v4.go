package client

import (
	"strconv"

	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

func addOrdersResultsFromGQLTypeV4(results *gqltypes.AddOrdersResultsV4) *AddOrdersResultsV4 {
	return &AddOrdersResultsV4{
		Accepted: acceptedOrderResultsFromGQLTypeV4(results.Accepted),
		Rejected: rejectedOrderResultsFromGQLTypeV4(results.Rejected),
	}
}

func acceptedOrderResultsFromGQLTypeV4(results []*gqltypes.AcceptedOrderResultV4) []*AcceptedOrderResultV4 {
	result := make([]*AcceptedOrderResultV4, len(results))
	for i, r := range results {
		result[i] = acceptedOrderResultFromGQLTypeV4(r)
	}
	return result
}

func acceptedOrderResultFromGQLTypeV4(result *gqltypes.AcceptedOrderResultV4) *AcceptedOrderResultV4 {
	return &AcceptedOrderResultV4{
		Order: orderWithMetadataFromGQLTypeV4(result.Order),
		IsNew: result.IsNew,
	}
}

func orderWithMetadataFromGQLTypeV4(order *gqltypes.OrderV4WithMetadata) *OrderWithMetadataV4 {
	sigType, _ := zeroex.SignatureTypeV4FromString(order.SignatureType)
	vValue, _ := strconv.ParseUint(order.SignatureV, 10, 8)
	return &OrderWithMetadataV4{
		Hash:                common.HexToHash(order.Hash),
		ChainID:             math.MustParseBig256(order.ChainID),
		ExchangeAddress:     common.HexToAddress(order.ExchangeAddress),
		MakerToken:          common.HexToAddress(order.MakerToken),
		TakerToken:          common.HexToAddress(order.TakerToken),
		Maker:               common.HexToAddress(order.Maker),
		Taker:               common.HexToAddress(order.Taker),
		Sender:              common.HexToAddress(order.Sender),
		FeeRecipient:        common.HexToAddress(order.FeeRecipient),
		MakerAmount:         math.MustParseBig256(order.MakerAmount),
		TakerAmount:         math.MustParseBig256(order.TakerAmount),
		TakerTokenFeeAmount: math.MustParseBig256(order.TakerTokenFeeAmount),
		Salt:                math.MustParseBig256(order.Salt),
		Expiry:              math.MustParseBig256(order.Expiry),
		Pool:                zeroex.BigToBytes32(math.MustParseBig256(order.Pool)),
		Signature: zeroex.SignatureFieldV4{
			SignatureType: sigType,
			V:             uint8(vValue),
			R:             zeroex.HexToBytes32(order.SignatureR),
			S:             zeroex.HexToBytes32(order.SignatureS),
		},
	}
}

func rejectedOrderResultsFromGQLTypeV4(results []*gqltypes.RejectedOrderResultV4) []*RejectedOrderResultV4 {
	result := make([]*RejectedOrderResultV4, len(results))
	for i, r := range results {
		result[i] = rejectedOrderResultFromGQLTypeV4(r)
	}
	return result
}

func rejectedOrderResultFromGQLTypeV4(result *gqltypes.RejectedOrderResultV4) *RejectedOrderResultV4 {
	var hash *common.Hash
	if result.Hash != nil {
		h := common.HexToHash(*result.Hash)
		hash = &h
	}
	order := result.Order
	sigType, _ := zeroex.SignatureTypeV4FromString(order.SignatureType)
	vValue, _ := strconv.ParseUint(order.SignatureV, 10, 8)

	return &RejectedOrderResultV4{
		Hash: hash,
		Order: &OrderWithMetadataV4{
			// FIXME(oskar) - ??
			Hash:                common.HexToHash(order.Maker),
			ChainID:             math.MustParseBig256(order.ChainID),
			ExchangeAddress:     common.HexToAddress(order.ExchangeAddress),
			MakerToken:          common.HexToAddress(order.MakerToken),
			TakerToken:          common.HexToAddress(order.TakerToken),
			Maker:               common.HexToAddress(order.Maker),
			Taker:               common.HexToAddress(order.Taker),
			Sender:              common.HexToAddress(order.Sender),
			FeeRecipient:        common.HexToAddress(order.FeeRecipient),
			MakerAmount:         math.MustParseBig256(order.MakerAmount),
			TakerAmount:         math.MustParseBig256(order.TakerAmount),
			TakerTokenFeeAmount: math.MustParseBig256(order.TakerTokenFeeAmount),
			Salt:                math.MustParseBig256(order.Salt),
			Expiry:              math.MustParseBig256(order.Expiry),
			Pool:                zeroex.BigToBytes32(math.MustParseBig256(order.Pool)),
			Signature: zeroex.SignatureFieldV4{
				SignatureType: sigType,
				V:             uint8(vValue),
				R:             zeroex.HexToBytes32(order.SignatureR),
				S:             zeroex.HexToBytes32(order.SignatureS),
			},
		},
		Code:    result.Code,
		Message: result.Message,
	}
}

func ordersWithMetadataFromGQLTypeV4(orders []*gqltypes.OrderV4WithMetadata) []*OrderWithMetadataV4 {
	result := make([]*OrderWithMetadataV4, len(orders))
	for i, r := range orders {
		result[i] = orderWithMetadataFromGQLTypeV4(r)
	}
	return result
}
