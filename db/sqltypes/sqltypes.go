// +build !js

package sqltypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gibson042/canonicaljson-go"
)

// BigInt is a wrapper around *big.Int that implements the sql.Valuer
// and sql.Scanner interfaces and *does not* retain sort order.
type BigInt struct {
	*big.Int
}

func NewBigInt(v *big.Int) *BigInt {
	return &BigInt{
		Int: v,
	}
}

func BigIntFromString(v string) (*BigInt, error) {
	bigInt, ok := ethmath.ParseBig256(v)
	if !ok {
		return nil, fmt.Errorf("dexietypes: could not convert %q to BigInt", v)
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
	return i.Int.String(), nil
}

func (i *BigInt) Scan(value interface{}) error {
	if value == nil {
		i = nil
		return nil
	}
	switch v := value.(type) {
	case int64:
		i.Int = big.NewInt(v)
	case float64:
		if math.Trunc(v) != v {
			// float64 may be used by the database driver to represent 0 or any other
			// whole number. This is okay as long as v is a whole number, i.e. does not
			// have anything after the decimal point. If this is not the case we return
			// an error.
			return fmt.Errorf("could not scan non-whole number float64 value %v into sqltypes.BigInt", value)
		}
		i.Int, _ = big.NewFloat(v).Int(big.NewInt(0))
	case string:
		parsed, ok := ethmath.ParseBig256(v)
		if !ok {
			return fmt.Errorf("could not scan string value %q into sqltypes.BigInt", v)
		}
		i.Int = parsed
	default:
		return fmt.Errorf("could not scan type %T into sqltypes.BigInt", value)
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
		return fmt.Errorf("could not unmarshal JSON data into dexietypes.BigInt: %s", string(data))
	}
	bigInt, ok := ethmath.ParseBig256(unqouted)
	if !ok {
		return fmt.Errorf("could not unmarshal JSON data into dexietypes.BigInt: %s", string(data))
	}
	i.Int = bigInt
	return nil
}

// SortedBigInt is a wrapper around *big.Int that implements the sql.Valuer
// and sql.Scanner interfaces and retains sort order by padding with zeroes.
type SortedBigInt struct {
	*big.Int
}

func NewSortedBigInt(v *big.Int) *SortedBigInt {
	return &SortedBigInt{
		Int: v,
	}
}

func SortedBigIntFromString(v string) (*SortedBigInt, error) {
	bigInt, ok := ethmath.ParseBig256(v)
	if !ok {
		return nil, fmt.Errorf("dexietypes: could not convert %q to BigInt", v)
	}
	return NewSortedBigInt(bigInt), nil
}

func SortedBigIntFromInt64(v int64) *SortedBigInt {
	return NewSortedBigInt(big.NewInt(v))
}

func (i *SortedBigInt) Value() (driver.Value, error) {
	if i == nil || i.Int == nil {
		return nil, nil
	}
	// Note(albrow), strings in SQL are sorted in alphanumerical order, not
	// numerical order. In order to sort by numerical order, we need to pad with
	// zeroes. The maximum length of an unsigned 256 bit integer is 80, so we
	// pad with zeroes such that the length of the number is always 80.
	return fmt.Sprintf("%080s", i.Int.String()), nil
}

func (i *SortedBigInt) Scan(value interface{}) error {
	if value == nil {
		i = nil
		return nil
	}
	switch v := value.(type) {
	case int64:
		i.Int = big.NewInt(v)
	case float64:
		if math.Trunc(v) != v {
			// float64 may be used by the database driver to represent 0 or any other
			// whole number. This is okay as long as v is a whole number, i.e. does not
			// have anything after the decimal point. If this is not the case we return
			// an error.
			return fmt.Errorf("could not scan non-whole number float64 value %v into sqltypes.BigInt", value)
		}
		i.Int, _ = big.NewFloat(v).Int(big.NewInt(0))
	case string:
		parsed, ok := ethmath.ParseBig256(v)
		if !ok {
			return fmt.Errorf("could not scan string value %q into sqltypes.BigInt", v)
		}
		i.Int = parsed
	default:
		return fmt.Errorf("could not scan type %T into sqltypes.BigInt", value)
	}

	return nil
}

func (i *SortedBigInt) MarshalJSON() ([]byte, error) {
	if i == nil || i.Int == nil {
		return json.Marshal(nil)
	}
	// Note(albrow), strings in Dexie.js are sorted in alphanumerical order, not
	// numerical order. In order to sort by numerical order, we need to pad with
	// zeroes. The maximum length of an unsigned 256 bit integer is 80, so we
	// pad with zeroes such that the length of the number is always 80.
	return json.Marshal(fmt.Sprintf("%080s", i.Int.String()))
}

func (i *SortedBigInt) UnmarshalJSON(data []byte) error {
	unqouted, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON data into dexietypes.BigInt: %s", string(data))
	}
	bigInt, ok := ethmath.ParseBig256(unqouted)
	if !ok {
		return fmt.Errorf("could not unmarshal JSON data into dexietypes.BigInt: %s", string(data))
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
	return canonicaljson.Marshal(s)
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
	ChainID                  *SortedBigInt    `db:"chainID"`
	ExchangeAddress          common.Address   `db:"exchangeAddress"`
	MakerAddress             common.Address   `db:"makerAddress"`
	MakerAssetData           []byte           `db:"makerAssetData"`
	MakerFeeAssetData        []byte           `db:"makerFeeAssetData"`
	MakerAssetAmount         *SortedBigInt    `db:"makerAssetAmount"`
	MakerFee                 *SortedBigInt    `db:"makerFee"`
	TakerAddress             common.Address   `db:"takerAddress"`
	TakerAssetData           []byte           `db:"takerAssetData"`
	TakerFeeAssetData        []byte           `db:"takerFeeAssetData"`
	TakerAssetAmount         *SortedBigInt    `db:"takerAssetAmount"`
	TakerFee                 *SortedBigInt    `db:"takerFee"`
	SenderAddress            common.Address   `db:"senderAddress"`
	FeeRecipientAddress      common.Address   `db:"feeRecipientAddress"`
	ExpirationTimeSeconds    *SortedBigInt    `db:"expirationTimeSeconds"`
	Salt                     *SortedBigInt    `db:"salt"`
	Signature                []byte           `db:"signature"`
	LastUpdated              time.Time        `db:"lastUpdated"`
	FillableTakerAssetAmount *SortedBigInt    `db:"fillableTakerAssetAmount"`
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
	logsJSON, err := canonicaljson.Marshal(e.Logs)
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
	Hash      common.Hash   `db:"hash"`
	Parent    common.Hash   `db:"parent"`
	Number    *SortedBigInt `db:"number"`
	Timestamp time.Time     `db:"timestamp"`
	Logs      *EventLogs    `db:"logs"`
}

type Metadata struct {
	EthereumChainID                   int           `db:"ethereumChainID"`
	MaxExpirationTime                 *SortedBigInt `db:"maxExpirationTime"`
	EthRPCRequestsSentInCurrentUTCDay int           `db:"ethRPCRequestsSentInCurrentUTCDay"`
	StartOfCurrentUTCDay              time.Time     `db:"startOfCurrentUTCDay"`
}

func OrderToCommonType(order *Order) *types.OrderWithMetadata {
	if order == nil {
		return nil
	}
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
	if order == nil {
		return nil
	}
	return &Order{
		Hash:                     order.Hash,
		ChainID:                  NewSortedBigInt(order.ChainID),
		ExchangeAddress:          order.ExchangeAddress,
		MakerAddress:             order.MakerAddress,
		MakerAssetData:           order.MakerAssetData,
		MakerFeeAssetData:        order.MakerFeeAssetData,
		MakerAssetAmount:         NewSortedBigInt(order.MakerAssetAmount),
		MakerFee:                 NewSortedBigInt(order.MakerFee),
		TakerAddress:             order.TakerAddress,
		TakerAssetData:           order.TakerAssetData,
		TakerFeeAssetData:        order.TakerFeeAssetData,
		TakerAssetAmount:         NewSortedBigInt(order.TakerAssetAmount),
		TakerFee:                 NewSortedBigInt(order.TakerFee),
		SenderAddress:            order.SenderAddress,
		FeeRecipientAddress:      order.FeeRecipientAddress,
		ExpirationTimeSeconds:    NewSortedBigInt(order.ExpirationTimeSeconds),
		Salt:                     NewSortedBigInt(order.Salt),
		Signature:                order.Signature,
		LastUpdated:              order.LastUpdated,
		FillableTakerAssetAmount: NewSortedBigInt(order.FillableTakerAssetAmount),
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
	result := ParsedAssetData(make([]*SingleAssetData, len(parsedAssetData)))
	for i, singleAssetData := range parsedAssetData {
		result[i] = SingleAssetDataFromCommonType(singleAssetData)
	}
	return &result
}

func SingleAssetDataToCommonType(singleAssetData *SingleAssetData) *types.SingleAssetData {
	if singleAssetData == nil {
		return nil
	}
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
	if singleAssetData == nil {
		return nil
	}
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
	if miniHeader == nil {
		return nil
	}
	return &types.MiniHeader{
		Hash:      miniHeader.Hash,
		Parent:    miniHeader.Parent,
		Number:    miniHeader.Number.Int,
		Timestamp: miniHeader.Timestamp,
		Logs:      miniHeader.Logs.Logs,
	}
}

func MiniHeaderFromCommonType(miniHeader *types.MiniHeader) *MiniHeader {
	if miniHeader == nil {
		return nil
	}
	return &MiniHeader{
		Hash:      miniHeader.Hash,
		Parent:    miniHeader.Parent,
		Number:    NewSortedBigInt(miniHeader.Number),
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

func MetadataToCommonType(metadata *Metadata) *types.Metadata {
	if metadata == nil {
		return nil
	}
	return &types.Metadata{
		EthereumChainID:                   metadata.EthereumChainID,
		MaxExpirationTime:                 metadata.MaxExpirationTime.Int,
		EthRPCRequestsSentInCurrentUTCDay: metadata.EthRPCRequestsSentInCurrentUTCDay,
		StartOfCurrentUTCDay:              metadata.StartOfCurrentUTCDay,
	}
}

func MetadataFromCommonType(metadata *types.Metadata) *Metadata {
	if metadata == nil {
		return nil
	}
	return &Metadata{
		EthereumChainID:                   metadata.EthereumChainID,
		MaxExpirationTime:                 NewSortedBigInt(metadata.MaxExpirationTime),
		EthRPCRequestsSentInCurrentUTCDay: metadata.EthRPCRequestsSentInCurrentUTCDay,
		StartOfCurrentUTCDay:              metadata.StartOfCurrentUTCDay,
	}
}
