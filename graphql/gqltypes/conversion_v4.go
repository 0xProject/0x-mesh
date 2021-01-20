package gqltypes

import (
	"strconv"
	"strings"

	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
)

func AddOrdersResultsFromValidationResultsV4(validationResults *ordervalidator.ValidationResults) (*AddOrdersResultsV4, error) {
	rejected, err := RejectedOrderResultsFromOrderInfosV4(validationResults.Rejected)
	if err != nil {
		return nil, err
	}
	return &AddOrdersResultsV4{
		Accepted: AcceptedOrderResultsFromOrderInfosV4(validationResults.Accepted),
		Rejected: rejected,
	}, nil
}

func RejectedOrderResultsFromOrderInfosV4(infos []*ordervalidator.RejectedOrderInfo) ([]*RejectedOrderResultV4, error) {
	result := make([]*RejectedOrderResultV4, len(infos))
	for i, info := range infos {
		rejectedResult, err := RejectedOrderResultFromOrderInfoV4(info)
		if err != nil {
			return nil, err
		}
		result[i] = rejectedResult
	}
	return result, nil
}

func RejectedOrderResultFromOrderInfoV4(info *ordervalidator.RejectedOrderInfo) (*RejectedOrderResultV4, error) {
	var hash *string
	if hashString := info.OrderHash.Hex(); hashString != "0x" {
		hash = &hashString
	}
	code, err := RejectedCodeFromValidatorStatus(info.Status)
	if err != nil {
		return nil, err
	}
	return &RejectedOrderResultV4{
		Hash: hash,
		Order: &OrderV4{
			ChainID:             info.SignedOrderV4.OrderV4.ChainID.String(),
			ExchangeAddress:     strings.ToLower(info.SignedOrderV4.OrderV4.ExchangeAddress.Hex()),
			Maker:               strings.ToLower(info.SignedOrderV4.OrderV4.Maker.Hex()),
			Taker:               strings.ToLower(info.SignedOrderV4.OrderV4.Taker.Hex()),
			Sender:              strings.ToLower(info.SignedOrderV4.OrderV4.Sender.Hex()),
			MakerAmount:         info.SignedOrderV4.OrderV4.MakerAmount.String(),
			MakerToken:          strings.ToLower(info.SignedOrderV4.OrderV4.MakerToken.Hex()),
			TakerAmount:         info.SignedOrderV4.OrderV4.TakerAmount.String(),
			TakerToken:          strings.ToLower(info.SignedOrderV4.OrderV4.TakerToken.Hex()),
			TakerTokenFeeAmount: info.SignedOrderV4.OrderV4.TakerTokenFeeAmount.String(),
			Pool:                info.SignedOrderV4.OrderV4.Pool.String(),
			Expiry:              info.SignedOrderV4.OrderV4.Expiry.String(),
			Salt:                info.SignedOrderV4.OrderV4.Salt.String(),
		},
		Code:    code,
		Message: info.Status.Message,
	}, nil
}

func AcceptedOrderResultFromOrderInfoV4(info *ordervalidator.AcceptedOrderInfo) *AcceptedOrderResultV4 {
	return &AcceptedOrderResultV4{
		Order: &OrderV4WithMetadata{
			Hash:                     info.OrderHash.Hex(),
			ChainID:                  info.SignedOrderV4.ChainID.String(),
			ExchangeAddress:          strings.ToLower(info.SignedOrderV4.ExchangeAddress.Hex()),
			Maker:                    strings.ToLower(info.SignedOrderV4.Maker.Hex()),
			Taker:                    strings.ToLower(info.SignedOrderV4.Taker.Hex()),
			Sender:                   strings.ToLower(info.SignedOrderV4.Sender.Hex()),
			MakerAmount:              info.SignedOrderV4.MakerAmount.String(),
			MakerToken:               strings.ToLower(info.SignedOrderV4.MakerToken.Hex()),
			TakerAmount:              info.SignedOrderV4.TakerAmount.String(),
			TakerToken:               strings.ToLower(info.SignedOrderV4.TakerToken.Hex()),
			TakerTokenFeeAmount:      info.SignedOrderV4.TakerTokenFeeAmount.String(),
			Pool:                     info.SignedOrderV4.Pool.String(),
			Expiry:                   info.SignedOrderV4.Expiry.String(),
			Salt:                     info.SignedOrderV4.Salt.String(),
			SignatureType:            info.SignedOrderV4.Signature.SignatureType.String(),
			SignatureV:               strconv.FormatUint(uint64(info.SignedOrderV4.Signature.V), 10),
			SignatureR:               info.SignedOrderV4.Signature.R.String(),
			SignatureS:               info.SignedOrderV4.Signature.S.String(),
			FillableTakerAssetAmount: info.FillableTakerAssetAmount.String(),
		},
		IsNew: info.IsNew,
	}
}

func AcceptedOrderResultsFromOrderInfosV4(infos []*ordervalidator.AcceptedOrderInfo) []*AcceptedOrderResultV4 {
	result := make([]*AcceptedOrderResultV4, len(infos))
	for i, info := range infos {
		added := AcceptedOrderResultFromOrderInfoV4(info)
		result[i] = added
	}
	return result
}
