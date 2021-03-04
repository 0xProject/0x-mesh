package gqltypes

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	log "github.com/sirupsen/logrus"
)

func AddOrderOptsToCommonType(opts *AddOrdersOpts) *types.AddOrdersOpts {
	commonTypeOpts := &types.AddOrdersOpts{}
	if opts.KeepCancelled != nil {
		commonTypeOpts.KeepCancelled = *opts.KeepCancelled
	}
	if opts.KeepExpired != nil {
		commonTypeOpts.KeepExpired = *opts.KeepExpired
	}
	if opts.KeepFullyFilled != nil {
		commonTypeOpts.KeepFullyFilled = *opts.KeepFullyFilled
	}
	if opts.KeepUnfunded != nil {
		commonTypeOpts.KeepUnfunded = *opts.KeepUnfunded
	}
	return commonTypeOpts
}

func StatsFromCommonType(stats *types.Stats) *Stats {
	return &Stats{
		Version:     stats.Version,
		PubSubTopic: stats.PubSubTopic,
		Rendezvous:  stats.Rendezvous,
		PeerID:      stats.PeerID,
		// TODO(albrow): This should be a big.Int in core package.
		EthereumChainID: stats.EthereumChainID,
		// TODO(albrow): LatestBlock should be a pointer in core package.
		LatestBlock:                       LatestBlockFromCommonType(stats.LatestBlock),
		NumPeers:                          stats.NumPeers,
		NumOrders:                         stats.NumOrders,
		NumOrdersIncludingRemoved:         stats.NumOrdersIncludingRemoved,
		StartOfCurrentUTCDay:              stats.StartOfCurrentUTCDay.Format(time.RFC3339),
		EthRPCRequestsSentInCurrentUTCDay: stats.EthRPCRequestsSentInCurrentUTCDay,
		EthRPCRateLimitExpiredRequests:    int(stats.EthRPCRateLimitExpiredRequests),
		SecondaryRendezvous:               stats.SecondaryRendezvous,
		MaxExpirationTime:                 stats.MaxExpirationTime.String(),
	}
}

func LatestBlockFromCommonType(latestBlock types.LatestBlock) *LatestBlock {
	return &LatestBlock{
		Number: latestBlock.Number.String(),
		Hash:   latestBlock.Hash.String(),
	}
}

func NewOrderToSignedOrder(newOrder *NewOrder) *zeroex.SignedOrder {
	return &zeroex.SignedOrder{
		Order: zeroex.Order{
			ChainID:               math.MustParseBig256(newOrder.ChainID),
			ExchangeAddress:       common.HexToAddress(newOrder.ExchangeAddress),
			MakerAddress:          common.HexToAddress(newOrder.MakerAddress),
			MakerAssetData:        types.HexToBytes(newOrder.MakerAssetData),
			MakerFeeAssetData:     types.HexToBytes(newOrder.MakerFeeAssetData),
			MakerAssetAmount:      math.MustParseBig256(newOrder.MakerAssetAmount),
			MakerFee:              math.MustParseBig256(newOrder.MakerFee),
			TakerAddress:          common.HexToAddress(newOrder.TakerAddress),
			TakerAssetData:        types.HexToBytes(newOrder.TakerAssetData),
			TakerFeeAssetData:     types.HexToBytes(newOrder.TakerFeeAssetData),
			TakerAssetAmount:      math.MustParseBig256(newOrder.TakerAssetAmount),
			TakerFee:              math.MustParseBig256(newOrder.TakerFee),
			SenderAddress:         common.HexToAddress(newOrder.SenderAddress),
			FeeRecipientAddress:   common.HexToAddress(newOrder.FeeRecipientAddress),
			ExpirationTimeSeconds: math.MustParseBig256(newOrder.ExpirationTimeSeconds),
			Salt:                  math.MustParseBig256(newOrder.Salt),
		},
		Signature: types.HexToBytes(newOrder.Signature),
	}
}

func NewOrderToSignedOrderV4(newOrder *NewOrderV4) *zeroex.SignedOrderV4 {
	signatureType, err := zeroex.SignatureTypeV4FromString(newOrder.SignatureType)
	if err != nil {
		panic(err)
	}
	return &zeroex.SignedOrderV4{
		OrderV4: zeroex.OrderV4{
			ChainID:             math.MustParseBig256(newOrder.ChainID),
			VerifyingContract:   common.HexToAddress(newOrder.VerifyingContract),
			MakerToken:          common.HexToAddress(newOrder.MakerToken),
			TakerToken:          common.HexToAddress(newOrder.TakerToken),
			Maker:               common.HexToAddress(newOrder.Maker),
			Taker:               common.HexToAddress(newOrder.Taker),
			Sender:              common.HexToAddress(newOrder.Sender),
			FeeRecipient:        common.HexToAddress(newOrder.FeeRecipient),
			MakerAmount:         math.MustParseBig256(newOrder.MakerAmount),
			TakerAmount:         math.MustParseBig256(newOrder.TakerAmount),
			TakerTokenFeeAmount: math.MustParseBig256(newOrder.TakerTokenFeeAmount),
			Salt:                math.MustParseBig256(newOrder.Salt),
			Expiry:              math.MustParseBig256(newOrder.Expiry),
			Pool:                zeroex.BigToBytes32(math.MustParseBig256(newOrder.Pool)),
		},
		Signature: zeroex.SignatureFieldV4{
			SignatureType: signatureType,
			V:             parseUint8FromStringOrPanic(newOrder.SignatureV),
			R:             zeroex.BigToBytes32(math.MustParseBig256(newOrder.SignatureR)),
			S:             zeroex.BigToBytes32(math.MustParseBig256(newOrder.SignatureS)),
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

func NewOrdersToSignedOrdersV4(newOrders []*NewOrderV4) []*zeroex.SignedOrderV4 {
	result := make([]*zeroex.SignedOrderV4, len(newOrders))
	for i, newOrder := range newOrders {
		result[i] = NewOrderToSignedOrderV4(newOrder)
	}
	return result
}

func NewOrderFromSignedOrder(signedOrder *zeroex.SignedOrder) *NewOrder {
	return &NewOrder{
		ChainID:               signedOrder.ChainID.String(),
		ExchangeAddress:       strings.ToLower(signedOrder.ExchangeAddress.Hex()),
		MakerAddress:          strings.ToLower(signedOrder.MakerAddress.Hex()),
		MakerAssetData:        types.BytesToHex(signedOrder.MakerAssetData),
		MakerFeeAssetData:     types.BytesToHex(signedOrder.MakerFeeAssetData),
		MakerAssetAmount:      signedOrder.MakerAssetAmount.String(),
		MakerFee:              signedOrder.MakerFee.String(),
		TakerAddress:          strings.ToLower(signedOrder.TakerAddress.Hex()),
		TakerAssetData:        types.BytesToHex(signedOrder.TakerAssetData),
		TakerFeeAssetData:     types.BytesToHex(signedOrder.TakerFeeAssetData),
		TakerAssetAmount:      signedOrder.TakerAssetAmount.String(),
		TakerFee:              signedOrder.TakerFee.String(),
		SenderAddress:         strings.ToLower(signedOrder.SenderAddress.Hex()),
		FeeRecipientAddress:   strings.ToLower(signedOrder.FeeRecipientAddress.Hex()),
		ExpirationTimeSeconds: signedOrder.ExpirationTimeSeconds.String(),
		Salt:                  signedOrder.Salt.String(),
		Signature:             types.BytesToHex(signedOrder.Signature),
	}
}

func NewOrderFromSignedOrderV4(signedOrder *zeroex.SignedOrderV4) *NewOrderV4 {
	return &NewOrderV4{
		ChainID:             signedOrder.ChainID.String(),
		VerifyingContract:   strings.ToLower(signedOrder.VerifyingContract.Hex()),
		MakerToken:          strings.ToLower(signedOrder.MakerToken.Hex()),
		TakerToken:          strings.ToLower(signedOrder.TakerToken.Hex()),
		MakerAmount:         signedOrder.MakerAmount.String(),
		TakerAmount:         signedOrder.TakerAmount.String(),
		TakerTokenFeeAmount: signedOrder.TakerTokenFeeAmount.String(),
		Maker:               strings.ToLower(signedOrder.Maker.Hex()),
		Taker:               strings.ToLower(signedOrder.Taker.Hex()),
		Sender:              strings.ToLower(signedOrder.Sender.Hex()),
		FeeRecipient:        strings.ToLower(signedOrder.FeeRecipient.Hex()),
		Pool:                signedOrder.Pool.String(),
		Expiry:              signedOrder.Expiry.String(),
		Salt:                signedOrder.Salt.String(),
		SignatureType:       signedOrder.Signature.SignatureType.String(),
		SignatureV:          strconv.FormatUint(uint64(signedOrder.Signature.V), 10),
		SignatureR:          signedOrder.Signature.R.String(),
		SignatureS:          signedOrder.Signature.S.String(),
	}
}

func NewOrdersFromSignedOrders(signedOrders []*zeroex.SignedOrder) []*NewOrder {
	result := make([]*NewOrder, len(signedOrders))
	for i, order := range signedOrders {
		result[i] = NewOrderFromSignedOrder(order)
	}
	return result
}

func NewOrdersFromSignedOrdersV4(signedOrders []*zeroex.SignedOrderV4) []*NewOrderV4 {
	result := make([]*NewOrderV4, len(signedOrders))
	for i, order := range signedOrders {
		result[i] = NewOrderFromSignedOrderV4(order)
	}
	return result
}

func OrderWithMetadataFromCommonType(order *types.OrderWithMetadata) *OrderWithMetadata {
	return &OrderWithMetadata{
		Hash:                     order.Hash.Hex(),
		ChainID:                  order.OrderV3.ChainID.String(),
		ExchangeAddress:          strings.ToLower(order.OrderV3.ExchangeAddress.Hex()),
		MakerAddress:             strings.ToLower(order.OrderV3.MakerAddress.Hex()),
		MakerAssetData:           types.BytesToHex(order.OrderV3.MakerAssetData),
		MakerFeeAssetData:        types.BytesToHex(order.OrderV3.MakerFeeAssetData),
		MakerAssetAmount:         order.OrderV3.MakerAssetAmount.String(),
		MakerFee:                 order.OrderV3.MakerFee.String(),
		TakerAddress:             strings.ToLower(order.OrderV3.TakerAddress.Hex()),
		TakerAssetData:           types.BytesToHex(order.OrderV3.TakerAssetData),
		TakerFeeAssetData:        types.BytesToHex(order.OrderV3.TakerFeeAssetData),
		TakerAssetAmount:         order.OrderV3.TakerAssetAmount.String(),
		TakerFee:                 order.OrderV3.TakerFee.String(),
		SenderAddress:            strings.ToLower(order.OrderV3.SenderAddress.Hex()),
		FeeRecipientAddress:      strings.ToLower(order.OrderV3.FeeRecipientAddress.Hex()),
		ExpirationTimeSeconds:    order.OrderV3.ExpirationTimeSeconds.String(),
		Salt:                     order.OrderV3.Salt.String(),
		Signature:                types.BytesToHex(order.Signature),
		FillableTakerAssetAmount: order.FillableTakerAssetAmount.String(),
	}
}

func OrderWithMetadataFromCommonTypeV4(order *types.OrderWithMetadata) *OrderV4WithMetadata {
	return &OrderV4WithMetadata{
		Hash:                     order.Hash.Hex(),
		ChainID:                  order.OrderV4.ChainID.String(),
		VerifyingContract:        strings.ToLower(order.OrderV4.VerifyingContract.Hex()),
		Maker:                    strings.ToLower(order.OrderV4.Maker.Hex()),
		Taker:                    strings.ToLower(order.OrderV4.Taker.Hex()),
		Sender:                   strings.ToLower(order.OrderV4.Sender.Hex()),
		MakerAmount:              order.OrderV4.MakerAmount.String(),
		MakerToken:               strings.ToLower(order.OrderV4.MakerToken.Hex()),
		TakerAmount:              order.OrderV4.TakerAmount.String(),
		TakerToken:               strings.ToLower(order.OrderV4.TakerToken.Hex()),
		TakerTokenFeeAmount:      order.OrderV4.TakerTokenFeeAmount.String(),
		Pool:                     order.OrderV4.Pool.String(),
		Expiry:                   order.OrderV4.Expiry.String(),
		Salt:                     order.OrderV4.Salt.String(),
		SignatureType:            order.SignatureV4.SignatureType.String(),
		SignatureV:               strconv.FormatUint(uint64(order.SignatureV4.V), 10),
		SignatureR:               order.SignatureV4.R.String(),
		SignatureS:               order.SignatureV4.S.String(),
		FillableTakerAssetAmount: order.FillableTakerAssetAmount.String(),
	}
}

func OrdersWithMetadataFromCommonType(orders []*types.OrderWithMetadata) []*OrderWithMetadata {
	result := make([]*OrderWithMetadata, len(orders))
	for i, order := range orders {
		result[i] = OrderWithMetadataFromCommonType(order)
	}
	return result
}

func OrdersWithMetadataFromCommonTypeV4(orders []*types.OrderWithMetadata) []*OrderV4WithMetadata {
	result := make([]*OrderV4WithMetadata, len(orders))
	for i, order := range orders {
		result[i] = OrderWithMetadataFromCommonTypeV4(order)
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
			Hash:                     info.OrderHash.String(),
			ChainID:                  info.SignedOrder.ChainID.String(),
			ExchangeAddress:          strings.ToLower(info.SignedOrder.ExchangeAddress.Hex()),
			MakerAddress:             strings.ToLower(info.SignedOrder.MakerAddress.Hex()),
			MakerAssetData:           types.BytesToHex(info.SignedOrder.MakerAssetData),
			MakerFeeAssetData:        types.BytesToHex(info.SignedOrder.MakerFeeAssetData),
			MakerAssetAmount:         info.SignedOrder.MakerAssetAmount.String(),
			MakerFee:                 info.SignedOrder.MakerFee.String(),
			TakerAddress:             strings.ToLower(info.SignedOrder.TakerAddress.Hex()),
			TakerAssetData:           types.BytesToHex(info.SignedOrder.TakerAssetData),
			TakerFeeAssetData:        types.BytesToHex(info.SignedOrder.TakerFeeAssetData),
			TakerAssetAmount:         info.SignedOrder.TakerAssetAmount.String(),
			TakerFee:                 info.SignedOrder.TakerFee.String(),
			SenderAddress:            strings.ToLower(info.SignedOrder.SenderAddress.Hex()),
			FeeRecipientAddress:      strings.ToLower(info.SignedOrder.FeeRecipientAddress.Hex()),
			ExpirationTimeSeconds:    info.SignedOrder.ExpirationTimeSeconds.String(),
			Salt:                     info.SignedOrder.Salt.String(),
			Signature:                types.BytesToHex(info.SignedOrder.Signature),
			FillableTakerAssetAmount: info.FillableTakerAssetAmount.String(),
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
	var hash *string
	if hashString := info.OrderHash.Hex(); hashString != "0x" {
		hash = &hashString
	}
	code, err := RejectedCodeFromValidatorStatus(info.Status)
	if err != nil {
		return nil, err
	}
	return &RejectedOrderResult{
		Hash: hash,
		Order: &Order{
			ChainID:               info.SignedOrder.ChainID.String(),
			ExchangeAddress:       strings.ToLower(info.SignedOrder.ExchangeAddress.Hex()),
			MakerAddress:          strings.ToLower(info.SignedOrder.MakerAddress.Hex()),
			MakerAssetData:        types.BytesToHex(info.SignedOrder.MakerAssetData),
			MakerFeeAssetData:     types.BytesToHex(info.SignedOrder.MakerFeeAssetData),
			MakerAssetAmount:      info.SignedOrder.MakerAssetAmount.String(),
			MakerFee:              info.SignedOrder.MakerFee.String(),
			TakerAddress:          strings.ToLower(info.SignedOrder.TakerAddress.Hex()),
			TakerAssetData:        types.BytesToHex(info.SignedOrder.TakerAssetData),
			TakerFeeAssetData:     types.BytesToHex(info.SignedOrder.TakerFeeAssetData),
			TakerAssetAmount:      info.SignedOrder.TakerAssetAmount.String(),
			TakerFee:              info.SignedOrder.TakerFee.String(),
			SenderAddress:         strings.ToLower(info.SignedOrder.SenderAddress.Hex()),
			FeeRecipientAddress:   strings.ToLower(info.SignedOrder.FeeRecipientAddress.Hex()),
			ExpirationTimeSeconds: info.SignedOrder.ExpirationTimeSeconds.String(),
			Salt:                  info.SignedOrder.Salt.String(),
			Signature:             types.BytesToHex(info.SignedOrder.Signature),
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

func OrderEventFromZeroExType(event *zeroex.OrderEvent) *OrderEvent {
	baseEvent := &OrderEvent{
		EndState:       OrderEndState(event.EndState),
		Timestamp:      event.Timestamp.Format(time.RFC3339),
		ContractEvents: ContractEventsFromZeroExType(event.ContractEvents),
	}
	if event.SignedOrder != nil {
		baseEvent.Order = &OrderWithMetadata{
			Hash:                     event.OrderHash.String(),
			ChainID:                  event.SignedOrder.ChainID.String(),
			ExchangeAddress:          strings.ToLower(event.SignedOrder.ExchangeAddress.Hex()),
			MakerAddress:             strings.ToLower(event.SignedOrder.MakerAddress.Hex()),
			MakerAssetData:           types.BytesToHex(event.SignedOrder.MakerAssetData),
			MakerFeeAssetData:        types.BytesToHex(event.SignedOrder.MakerFeeAssetData),
			MakerAssetAmount:         event.SignedOrder.MakerAssetAmount.String(),
			MakerFee:                 event.SignedOrder.MakerFee.String(),
			TakerAddress:             strings.ToLower(event.SignedOrder.TakerAddress.Hex()),
			TakerAssetData:           types.BytesToHex(event.SignedOrder.TakerAssetData),
			TakerFeeAssetData:        types.BytesToHex(event.SignedOrder.TakerFeeAssetData),
			TakerAssetAmount:         event.SignedOrder.TakerAssetAmount.String(),
			TakerFee:                 event.SignedOrder.TakerFee.String(),
			SenderAddress:            strings.ToLower(event.SignedOrder.SenderAddress.Hex()),
			FeeRecipientAddress:      strings.ToLower(event.SignedOrder.FeeRecipientAddress.Hex()),
			ExpirationTimeSeconds:    event.SignedOrder.ExpirationTimeSeconds.String(),
			Salt:                     event.SignedOrder.Salt.String(),
			Signature:                types.BytesToHex(event.SignedOrder.Signature),
			FillableTakerAssetAmount: event.FillableTakerAssetAmount.String(),
		}
	} else {
		hash, err := event.SignedOrderV4.ComputeOrderHash()
		if err != nil {
			log.Error(err)
		}
		sigV := strconv.FormatUint(uint64(event.SignedOrderV4.Signature.V), 10)
		baseEvent.Orderv4 = &OrderV4WithMetadata{
			ChainID:                  event.SignedOrderV4.ChainID.String(),
			VerifyingContract:        event.SignedOrderV4.VerifyingContract.Hex(),
			MakerToken:               event.SignedOrderV4.MakerToken.Hex(),
			TakerToken:               event.SignedOrderV4.TakerToken.Hex(),
			MakerAmount:              event.SignedOrderV4.MakerAmount.String(),
			TakerAmount:              event.SignedOrderV4.TakerAmount.String(),
			TakerTokenFeeAmount:      event.SignedOrderV4.TakerTokenFeeAmount.String(),
			Maker:                    event.SignedOrderV4.Maker.Hex(),
			Taker:                    event.SignedOrderV4.Taker.Hex(),
			Sender:                   event.SignedOrderV4.Sender.Hex(),
			FeeRecipient:             event.SignedOrderV4.FeeRecipient.Hex(),
			Pool:                     event.SignedOrderV4.Pool.String(),
			Expiry:                   event.SignedOrderV4.Expiry.String(),
			Salt:                     event.SignedOrderV4.Salt.String(),
			SignatureType:            event.SignedOrderV4.Signature.SignatureType.String(),
			SignatureV:               sigV,
			SignatureR:               event.SignedOrderV4.Signature.R.String(),
			SignatureS:               event.SignedOrderV4.Signature.S.String(),
			Hash:                     hash.Hex(),
			FillableTakerAssetAmount: event.FillableTakerAssetAmount.String(),
		}
	}
	return baseEvent
}

func OrderEventsFromZeroExType(orderEvents []*zeroex.OrderEvent) []*OrderEvent {
	result := make([]*OrderEvent, len(orderEvents))
	for i, event := range orderEvents {
		result[i] = OrderEventFromZeroExType(event)
	}
	return result
}

func ContractEventFromZeroExType(event *zeroex.ContractEvent) *ContractEvent {
	return &ContractEvent{
		BlockHash:  event.BlockHash.Hex(),
		TxHash:     event.TxHash.Hex(),
		TxIndex:    int(event.TxIndex),
		LogIndex:   int(event.LogIndex),
		IsRemoved:  event.IsRemoved,
		Address:    strings.ToLower(event.Address.Hex()),
		Kind:       event.Kind,
		Parameters: event.Parameters,
	}
}

func ContractEventsFromZeroExType(events []*zeroex.ContractEvent) []*ContractEvent {
	result := make([]*ContractEvent, len(events))
	for i, event := range events {
		result[i] = ContractEventFromZeroExType(event)
	}
	return result
}

func RejectedCodeFromValidatorStatus(status ordervalidator.RejectedOrderStatus) (RejectedOrderCode, error) {
	switch status.Code {
	case ordervalidator.ROEthRPCRequestFailed.Code:
		return RejectedOrderCodeEthRPCRequestFailed, nil
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
	case ordervalidator.ROTakerAddressNotAllowed.Code:
		return RejectedOrderCodeTakerAddressNotAllowed, nil
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

// FilterValueFromJSON converts the filter value from the JSON type to the
// corresponding Go type. It returns an error if the JSON type does not match
// what was expected based on the filter field.
func FilterValueFromJSON(f OrderFilter) (interface{}, error) {
	switch f.Field {
	case OrderFieldChainID, OrderFieldMakerAssetAmount, OrderFieldMakerFee, OrderFieldTakerAssetAmount, OrderFieldTakerFee, OrderFieldExpirationTimeSeconds, OrderFieldSalt, OrderFieldFillableTakerAssetAmount:
		return stringToBigInt(f.Value)
	case OrderFieldHash:
		return stringToHash(f.Value)
	case OrderFieldExchangeAddress, OrderFieldMakerAddress, OrderFieldTakerAddress, OrderFieldSenderAddress, OrderFieldFeeRecipientAddress:
		return stringToAddress(f.Value)
	case OrderFieldMakerAssetData, OrderFieldMakerFeeAssetData, OrderFieldTakerAssetData, OrderFieldTakerFeeAssetData:
		return stringToBytes(f.Value)
	default:
		return "", fmt.Errorf("invalid filter field: %q", f.Field)
	}
}

// FilterValueFromJSONV4 converts the filter value from the JSON type to the
// corresponding Go type. It returns an error if the JSON type does not match
// what was expected based on the filter field.
func FilterValueFromJSONV4(f OrderFilterV4) (interface{}, error) {
	// TODO(oskar) add byte32 conversions here
	switch f.Field {
	case OrderFieldV4ChainID, OrderFieldV4MakerAmount, OrderFieldV4TakerAmount, OrderFieldV4TakerTokenFeeAmount, OrderFieldV4Expiry, OrderFieldV4Salt, OrderFieldV4FillableTakerAssetAmount:
		return stringToBigInt(f.Value)
	case OrderFieldV4Hash:
		return stringToHash(f.Value)
	case OrderFieldV4VerifyingContract, OrderFieldV4Maker, OrderFieldV4Taker, OrderFieldV4Sender, OrderFieldV4FeeRecipient:
		return stringToAddress(f.Value)
	default:
		return "", fmt.Errorf("invalid filter field: %q", f.Field)
	}
}

// FilterValueToJSON converts the filter value from a native Go type to the
// corresponding JSON value. It returns an error if the Go type does not match
// what was expected based on the filter field.
func FilterValueToJSON(f OrderFilter) (string, error) {
	switch f.Field {
	case OrderFieldChainID, OrderFieldMakerAssetAmount, OrderFieldMakerFee, OrderFieldTakerAssetAmount, OrderFieldTakerFee, OrderFieldExpirationTimeSeconds, OrderFieldSalt, OrderFieldFillableTakerAssetAmount:
		return bigIntToString(f.Value)
	case OrderFieldHash:
		return hashToString(f.Value)
	case OrderFieldExchangeAddress, OrderFieldMakerAddress, OrderFieldTakerAddress, OrderFieldSenderAddress, OrderFieldFeeRecipientAddress:
		return addressToString(f.Value)
	case OrderFieldMakerAssetData, OrderFieldMakerFeeAssetData, OrderFieldTakerAssetData, OrderFieldTakerFeeAssetData:
		return bytesToString(f.Value)
	default:
		return "", fmt.Errorf("invalid filter field: %q", f.Field)
	}
}

func bigIntToString(value interface{}) (string, error) {
	bigInt, ok := value.(*big.Int)
	if !ok {
		return "", fmt.Errorf("invalid type for filter value (expected *big.Int but got %T)", value)
	}
	return bigInt.String(), nil
}

func hashToString(value interface{}) (string, error) {
	hash, ok := value.(common.Hash)
	if !ok {
		return "", fmt.Errorf("invalid type for filter value (expected common.Hash but got %T)", value)
	}
	return hash.Hex(), nil
}

func addressToString(value interface{}) (string, error) {
	address, ok := value.(common.Address)
	if !ok {
		return "", fmt.Errorf("invalid type for filter value (expected common.Address but got %T)", value)
	}
	return strings.ToLower(address.Hex()), nil
}

func bytesToString(value interface{}) (string, error) {
	bytes, ok := value.([]byte)
	if !ok {
		return "", fmt.Errorf("invalid type for filter value (expected []byte but got %T)", value)
	}
	return types.BytesToHex(bytes), nil
}

func filterValueAsString(value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("invalid type for filter value (expected string but got %T)", value)
	}
	return valueString, nil
}

func stringToBigInt(value interface{}) (*big.Int, error) {
	valueString, err := filterValueAsString(value)
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
	valueString, err := filterValueAsString(value)
	if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(valueString), nil
}

func stringToAddress(value interface{}) (common.Address, error) {
	valueString, err := filterValueAsString(value)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(valueString), nil
}

func stringToBytes(value interface{}) ([]byte, error) {
	valueString, err := filterValueAsString(value)
	if err != nil {
		return nil, err
	}
	return types.HexToBytes(valueString), nil
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
