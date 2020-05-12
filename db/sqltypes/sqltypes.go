package sqltypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// BigInt is a wrapper around *big.Int that implements the sql.Valuer
// and sql.Scanner interfaces.
type BigInt struct {
	*big.Int
}

func NewBigInt(v *big.Int) *BigInt {
	return &BigInt{
		Int: v,
	}
}

func BigIntFromString(v string) (*BigInt, error) {
	bigInt, ok := math.ParseBig256(v)
	if !ok {
		return nil, fmt.Errorf("sqltypes: could not convert %q to BigInt", v)
	}
	return NewBigInt(bigInt), nil
}

func BigIntFromInt64(v int64) *BigInt {
	return NewBigInt(big.NewInt(v))
}

func (i *BigInt) Value() (driver.Value, error) {
	if i == nil || i.Int == nil {
		return nil, nil
	}
	return i.String(), nil
}

func (i *BigInt) Scan(value interface{}) error {
	if value == nil {
		i = nil
		return nil
	}
	switch v := value.(type) {
	case int64:
		i.Int = big.NewInt(v)
	case string:
		parsed, ok := math.ParseBig256(v)
		if !ok {
			return fmt.Errorf("could not scan string value %q into Uint256", v)
		}
		i.Int = parsed
	default:
		return fmt.Errorf("could not scan type %T into Uint256", value)
	}

	return nil
}

func (i *BigInt) MarshalJSON() ([]byte, error) {
	if i == nil || i.Int == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(i.Int.String())
}

func (i *BigInt) UnmarshalJSON(data []byte) error {
	unqouted, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON data into Uint256: %s", string(data))
	}
	bigInt, ok := math.ParseBig256(unqouted)
	if !ok {
		return fmt.Errorf("could not unmarshal JSON data into Uint256: %s", string(data))
	}
	i.Int = bigInt
	return nil
}

type SingleAssetData struct {
	Address common.Address `json:"address"`
	TokenID *BigInt        `json:"tokenID"`
}

// ParsedAssetData is a wrapper around []*SingleAssetData that implements the
// sql.Valuer and sql.Scanner interfaces.
type ParsedAssetData []*SingleAssetData

func (s *ParsedAssetData) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *ParsedAssetData) Scan(value interface{}) error {
	if value == nil {
		s = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("could not scan type %T into EventLogs", value)
	}
}

// Order is the SQL database representation a 0x order along with some relevant metadata.
type Order struct {
	Hash                     common.Hash      `db:"hash"`
	ChainID                  *BigInt          `db:"chainID"`
	ExchangeAddress          common.Address   `db:"exchangeAddress"`
	MakerAddress             common.Address   `db:"makerAddress"`
	MakerAssetData           []byte           `db:"makerAssetData"`
	MakerFeeAssetData        []byte           `db:"makerFeeAssetData"`
	MakerAssetAmount         *BigInt          `db:"makerAssetAmount"`
	MakerFee                 *BigInt          `db:"makerFee"`
	TakerAddress             common.Address   `db:"takerAddress"`
	TakerAssetData           []byte           `db:"takerAssetData"`
	TakerFeeAssetData        []byte           `db:"takerFeeAssetData"`
	TakerAssetAmount         *BigInt          `db:"takerAssetAmount"`
	TakerFee                 *BigInt          `db:"takerFee"`
	SenderAddress            common.Address   `db:"senderAddress"`
	FeeRecipientAddress      common.Address   `db:"feeRecipientAddress"`
	ExpirationTimeSeconds    *BigInt          `db:"expirationTimeSeconds"`
	Salt                     *BigInt          `db:"salt"`
	Signature                []byte           `db:"signature"`
	LastUpdated              time.Time        `db:"lastUpdated"`
	FillableTakerAssetAmount *BigInt          `db:"fillableTakerAssetAmount"`
	IsRemoved                bool             `db:"isRemoved"`
	IsPinned                 bool             `db:"isPinned"`
	ParsedMakerAssetData     *ParsedAssetData `db:"parsedMakerAssetData"`
	ParsedMakerFeeAssetData  *ParsedAssetData `db:"parsedMakerFeeAssetData"`
}

// EventLogs is a wrapper around []*ethtypes.Log that implements the
// sql.Valuer and sql.Scanner interfaces.
type EventLogs struct {
	Logs []ethtypes.Log
}

func NewEventLogs(logs []ethtypes.Log) *EventLogs {
	return &EventLogs{
		Logs: logs,
	}
}

func (e *EventLogs) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	logsJSON, err := json.Marshal(e.Logs)
	if err != nil {
		return nil, err
	}
	return logsJSON, err
}

func (e *EventLogs) Scan(value interface{}) error {
	if value == nil {
		e = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &e.Logs)
	case string:
		return json.Unmarshal([]byte(v), &e.Logs)
	default:
		return fmt.Errorf("could not scan type %T into EventLogs", value)
	}
}

type MiniHeader struct {
	Hash      common.Hash `db:"hash"`
	Parent    common.Hash `db:"parent"`
	Number    *BigInt     `db:"number"`
	Timestamp time.Time   `db:"timestamp"`
	Logs      *EventLogs  `db:"logs"`
}

func OrderToCommonType(order *Order) *types.OrderWithMetadata {
	return &types.OrderWithMetadata{
		Hash:                     order.Hash,
		ChainID:                  order.ChainID.Int,
		ExchangeAddress:          order.ExchangeAddress,
		MakerAddress:             order.MakerAddress,
		MakerAssetData:           order.MakerAssetData,
		MakerFeeAssetData:        order.MakerFeeAssetData,
		MakerAssetAmount:         order.MakerAssetAmount.Int,
		MakerFee:                 order.MakerFee.Int,
		TakerAddress:             order.TakerAddress,
		TakerAssetData:           order.TakerAssetData,
		TakerFeeAssetData:        order.TakerFeeAssetData,
		TakerAssetAmount:         order.TakerAssetAmount.Int,
		TakerFee:                 order.TakerFee.Int,
		SenderAddress:            order.SenderAddress,
		FeeRecipientAddress:      order.FeeRecipientAddress,
		ExpirationTimeSeconds:    order.ExpirationTimeSeconds.Int,
		Salt:                     order.Salt.Int,
		Signature:                order.Signature,
		FillableTakerAssetAmount: order.FillableTakerAssetAmount.Int,
		LastUpdated:              order.LastUpdated,
		IsRemoved:                order.IsRemoved,
		IsPinned:                 order.IsPinned,
		ParsedMakerAssetData:     ParsedAssetDataToCommonType(order.ParsedMakerAssetData),
		ParsedMakerFeeAssetData:  ParsedAssetDataToCommonType(order.ParsedMakerFeeAssetData),
	}
}

func OrderFromCommonType(order *types.OrderWithMetadata) *Order {
	return &Order{
		Hash:                     order.Hash,
		ChainID:                  NewBigInt(order.ChainID),
		ExchangeAddress:          order.ExchangeAddress,
		MakerAddress:             order.MakerAddress,
		MakerAssetData:           order.MakerAssetData,
		MakerFeeAssetData:        order.MakerFeeAssetData,
		MakerAssetAmount:         NewBigInt(order.MakerAssetAmount),
		MakerFee:                 NewBigInt(order.MakerFee),
		TakerAddress:             order.TakerAddress,
		TakerAssetData:           order.TakerAssetData,
		TakerFeeAssetData:        order.TakerFeeAssetData,
		TakerAssetAmount:         NewBigInt(order.TakerAssetAmount),
		TakerFee:                 NewBigInt(order.TakerFee),
		SenderAddress:            order.SenderAddress,
		FeeRecipientAddress:      order.FeeRecipientAddress,
		ExpirationTimeSeconds:    NewBigInt(order.ExpirationTimeSeconds),
		Salt:                     NewBigInt(order.Salt),
		Signature:                order.Signature,
		LastUpdated:              order.LastUpdated,
		FillableTakerAssetAmount: NewBigInt(order.FillableTakerAssetAmount),
		IsRemoved:                order.IsRemoved,
		IsPinned:                 order.IsPinned,
		ParsedMakerAssetData:     ParsedAssetDataFromCommonType(order.ParsedMakerAssetData),
		ParsedMakerFeeAssetData:  ParsedAssetDataFromCommonType(order.ParsedMakerFeeAssetData),
	}
}

func OrdersToCommonType(orders []*Order) []*types.OrderWithMetadata {
	result := make([]*types.OrderWithMetadata, len(orders))
	for i, order := range orders {
		result[i] = OrderToCommonType(order)
	}
	return result
}

func ParsedAssetDataToCommonType(parsedAssetData *ParsedAssetData) []*types.SingleAssetData {
	if parsedAssetData == nil || len(*parsedAssetData) == 0 {
		return nil
	}
	assetDataSlice := []*SingleAssetData(*parsedAssetData)
	result := make([]*types.SingleAssetData, len(assetDataSlice))
	for i, singleAssetData := range assetDataSlice {
		result[i] = SingleAssetDataToCommonType(singleAssetData)
	}
	return result
}

func ParsedAssetDataFromCommonType(parsedAssetData []*types.SingleAssetData) *ParsedAssetData {
	if len(parsedAssetData) == 0 {
		return nil
	}
	result := ParsedAssetData(make([]*SingleAssetData, len(parsedAssetData)))
	for i, singleAssetData := range parsedAssetData {
		result[i] = SingleAssetDataFromCommonType(singleAssetData)
	}
	return &result
}

func SingleAssetDataToCommonType(singleAssetData *SingleAssetData) *types.SingleAssetData {
	var tokenID *big.Int
	if singleAssetData.TokenID != nil {
		tokenID = singleAssetData.TokenID.Int
	}
	return &types.SingleAssetData{
		Address: singleAssetData.Address,
		TokenID: tokenID,
	}
}

func SingleAssetDataFromCommonType(singleAssetData *types.SingleAssetData) *SingleAssetData {
	var tokenID *BigInt
	if singleAssetData.TokenID != nil {
		tokenID = NewBigInt(singleAssetData.TokenID)
	}
	return &SingleAssetData{
		Address: singleAssetData.Address,
		TokenID: tokenID,
	}
}

func MiniHeaderToCommonType(miniHeader *MiniHeader) *types.MiniHeader {
	return &types.MiniHeader{
		Hash:      miniHeader.Hash,
		Parent:    miniHeader.Parent,
		Number:    miniHeader.Number.Int,
		Timestamp: miniHeader.Timestamp,
		Logs:      miniHeader.Logs.Logs,
	}
}

func MiniHeaderFromCommonType(miniHeader *types.MiniHeader) *MiniHeader {
	return &MiniHeader{
		Hash:      miniHeader.Hash,
		Parent:    miniHeader.Parent,
		Number:    NewBigInt(miniHeader.Number),
		Timestamp: miniHeader.Timestamp,
		Logs:      NewEventLogs(miniHeader.Logs),
	}
}

func MiniHeadersToCommonType(miniHeaders []*MiniHeader) []*types.MiniHeader {
	result := make([]*types.MiniHeader, len(miniHeaders))
	for i, miniHeader := range miniHeaders {
		result[i] = MiniHeaderToCommonType(miniHeader)
	}
	return result
}
