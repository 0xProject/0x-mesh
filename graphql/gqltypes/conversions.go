package gqltypes

import (
	"fmt"
	"math/big"
	"strings"

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
		EthereumChainID: stats.EthereumChainID,
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
		Signature: newOrder.Signature,
	}
}

func NewOrdersToSignedOrders(newOrders []*NewOrder) []*zeroex.SignedOrder {
	result := make([]*zeroex.SignedOrder, len(newOrders))
	for i, newOrder := range newOrders {
		result[i] = NewOrderToSignedOrder(newOrder)
	}
	return result
}

func NewOrderFromSignedOrder(signedOrder *zeroex.SignedOrder) *NewOrder {
	return &NewOrder{
		ChainID:               BigNumber(*signedOrder.ChainID),
		ExchangeAddress:       Address(signedOrder.ExchangeAddress),
		MakerAddress:          Address(signedOrder.MakerAddress),
		MakerAssetData:        signedOrder.MakerAssetData,
		MakerFeeAssetData:     signedOrder.MakerFeeAssetData,
		MakerAssetAmount:      BigNumber(*signedOrder.MakerAssetAmount),
		MakerFee:              BigNumber(*signedOrder.MakerFee),
		TakerAddress:          Address(signedOrder.TakerAddress),
		TakerAssetData:        signedOrder.TakerAssetData,
		TakerFeeAssetData:     signedOrder.TakerFeeAssetData,
		TakerAssetAmount:      BigNumber(*signedOrder.TakerAssetAmount),
		TakerFee:              BigNumber(*signedOrder.TakerFee),
		SenderAddress:         Address(signedOrder.SenderAddress),
		FeeRecipientAddress:   Address(signedOrder.FeeRecipientAddress),
		ExpirationTimeSeconds: BigNumber(*signedOrder.ExpirationTimeSeconds),
		Salt:                  BigNumber(*signedOrder.Salt),
		Signature:             signedOrder.Signature,
	}
}

func NewOrdersFromSignedOrders(signedOrders []*zeroex.SignedOrder) []*NewOrder {
	result := make([]*NewOrder, len(signedOrders))
	for i, order := range signedOrders {
		result[i] = NewOrderFromSignedOrder(order)
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

func OrderEventFromZeroExType(event *zeroex.OrderEvent) *OrderEvent {
	return &OrderEvent{
		Order: &OrderWithMetadata{
			Hash:                     Hash(event.OrderHash),
			ChainID:                  BigNumber(*event.SignedOrder.ChainID),
			ExchangeAddress:          Address(event.SignedOrder.ExchangeAddress),
			MakerAddress:             Address(event.SignedOrder.MakerAddress),
			MakerAssetData:           Bytes(event.SignedOrder.MakerAssetData),
			MakerFeeAssetData:        Bytes(event.SignedOrder.MakerFeeAssetData),
			MakerAssetAmount:         BigNumber(*event.SignedOrder.MakerAssetAmount),
			MakerFee:                 BigNumber(*event.SignedOrder.MakerFee),
			TakerAddress:             Address(event.SignedOrder.TakerAddress),
			TakerAssetData:           Bytes(event.SignedOrder.TakerAssetData),
			TakerFeeAssetData:        Bytes(event.SignedOrder.TakerFeeAssetData),
			TakerAssetAmount:         BigNumber(*event.SignedOrder.TakerAssetAmount),
			TakerFee:                 BigNumber(*event.SignedOrder.TakerFee),
			SenderAddress:            Address(event.SignedOrder.SenderAddress),
			FeeRecipientAddress:      Address(event.SignedOrder.FeeRecipientAddress),
			ExpirationTimeSeconds:    BigNumber(*event.SignedOrder.ExpirationTimeSeconds),
			Salt:                     BigNumber(*event.SignedOrder.Salt),
			Signature:                Bytes(event.SignedOrder.Signature),
			FillableTakerAssetAmount: BigNumber(*event.FillableTakerAssetAmount),
		},
		EndState:       OrderEndState(event.EndState),
		Timestamp:      event.Timestamp,
		ContractEvents: ContractEventsFromZeroExType(event.ContractEvents),
	}
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
		BlockHash:  Hash(event.BlockHash),
		TxHash:     Hash(event.TxHash),
		TxIndex:    int(event.TxIndex),
		LogIndex:   int(event.LogIndex),
		IsRemoved:  event.IsRemoved,
		Address:    Address(event.Address),
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
	return common.ToHex(bytes), nil
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

// type orderWithMetadataJSON struct {
// 	ChainID                  int            `json:"chainId"`
// 	ExchangeAddress          common.Address `json:"exchangeAddress"`
// 	MakerAddress             common.Address `json:"makerAddress"`
// 	MakerAssetData           string         `json:"makerAssetData"`
// 	MakerAssetAmount         BigNumber      `json:"makerAssetAmount"`
// 	MakerFeeAssetData        string         `json:"makerFeeAssetData"`
// 	MakerFee                 BigNumber      `json:"makerFee"`
// 	TakerAddress             common.Address `json:"takerAddress"`
// 	TakerAssetData           string         `json:"takerAssetData"`
// 	TakerAssetAmount         BigNumber      `json:"takerAssetAmount"`
// 	TakerFeeAssetData        string         `json:"takerFeeAssetData"`
// 	TakerFee                 BigNumber      `json:"takerFee"`
// 	SenderAddress            common.Address `json:"senderAddress"`
// 	FeeRecipientAddress      common.Address `json:"feeRecipientAddress"`
// 	ExpirationTimeSeconds    BigNumber      `json:"expirationTimeSeconds"`
// 	Salt                     BigNumber      `json:"salt"`
// 	Signature                string         `json:"signature"`
// 	Hash                     common.Hash    `json:"hash"`
// 	FillableTakerAssetAmount BigNumber      `json:"fillableTakerAssetAmount"`
// }

// func (order *OrderWithMetadata) UnmarshalJSON(data []byte) error {
// 	var holder orderWithMetadataJSON
// 	if err := json.Unmarshal(data, &holder); err != nil {
// 		return err
// 	}
// 	order.ChainID = BigNumber(*big.NewInt(int64(holder.ChainID)))
// 	order.ExchangeAddress = Address(holder.ExchangeAddress)
// 	order.MakerAddress = Address(holder.MakerAddress)
// 	order.MakerAssetData = common.FromHex(holder.MakerAssetData)
// 	order.MakerAssetAmount = holder.MakerAssetAmount
// 	order.MakerFeeAssetData = common.FromHex(holder.MakerFeeAssetData)
// 	order.MakerFee = holder.MakerFee
// 	order.TakerAddress = Address(holder.TakerAddress)
// 	order.TakerAssetData = common.FromHex(holder.TakerAssetData)
// 	order.TakerAssetAmount = holder.TakerAssetAmount
// 	order.TakerFeeAssetData = common.FromHex(holder.TakerFeeAssetData)
// 	order.TakerFee = holder.TakerFee
// 	order.SenderAddress = Address(holder.SenderAddress)
// 	order.FeeRecipientAddress = Address(holder.FeeRecipientAddress)
// 	order.ExpirationTimeSeconds = holder.ExpirationTimeSeconds
// 	order.Salt = holder.Salt
// 	order.Signature = common.FromHex(holder.Signature)
// 	order.Hash = Hash(holder.Hash)
// 	order.FillableTakerAssetAmount = holder.FillableTakerAssetAmount
// 	return nil
// }

// type newOrderJSON struct {
// 	ChainID               string `json:"chainId"`
// 	ExchangeAddress       string `json:"exchangeAddress"`
// 	MakerAddress          string `json:"makerAddress"`
// 	MakerAssetData        string `json:"makerAssetData"`
// 	MakerAssetAmount      string `json:"makerAssetAmount"`
// 	MakerFeeAssetData     string `json:"makerFeeAssetData"`
// 	MakerFee              string `json:"makerFee"`
// 	TakerAddress          string `json:"takerAddress"`
// 	TakerAssetData        string `json:"takerAssetData"`
// 	TakerAssetAmount      string `json:"takerAssetAmount"`
// 	TakerFeeAssetData     string `json:"takerFeeAssetData"`
// 	TakerFee              string `json:"takerFee"`
// 	SenderAddress         string `json:"senderAddress"`
// 	FeeRecipientAddress   string `json:"feeRecipientAddress"`
// 	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
// 	Salt                  string `json:"salt"`
// 	Signature             string `json:"signature"`
// }

// func (order *NewOrder) MarshalJSON() ([]byte, error) {
// 	holder := newOrderJSON{
// 		ChainID:               BigNumberToBigInt(order.ChainID).String(),
// 		ExchangeAddress:       strings.ToLower(common.Address(order.ExchangeAddress).Hex()),
// 		MakerAddress:          strings.ToLower(common.Address(order.MakerAddress).Hex()),
// 		MakerAssetData:        common.ToHex(order.MakerAssetData),
// 		MakerAssetAmount:      BigNumberToBigInt(order.MakerAssetAmount).String(),
// 		MakerFeeAssetData:     common.ToHex(order.MakerFeeAssetData),
// 		MakerFee:              BigNumberToBigInt(order.MakerFee).String(),
// 		TakerAddress:          strings.ToLower(common.Address(order.TakerAddress).Hex()),
// 		TakerAssetData:        common.ToHex(order.TakerAssetData),
// 		TakerAssetAmount:      BigNumberToBigInt(order.TakerAssetAmount).String(),
// 		TakerFeeAssetData:     common.ToHex(order.TakerFeeAssetData),
// 		TakerFee:              BigNumberToBigInt(order.TakerFee).String(),
// 		SenderAddress:         strings.ToLower(common.Address(order.SenderAddress).Hex()),
// 		FeeRecipientAddress:   strings.ToLower(common.Address(order.FeeRecipientAddress).Hex()),
// 		ExpirationTimeSeconds: BigNumberToBigInt(order.ExpirationTimeSeconds).String(),
// 		Salt:                  BigNumberToBigInt(order.Salt).String(),
// 		Signature:             common.ToHex(order.Signature),
// 	}
// 	return json.Marshal(holder)
// }
