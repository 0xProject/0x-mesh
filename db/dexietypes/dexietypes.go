package dexietypes

// TODO(albrow): Can some of these types be de-duped with sqltypes without
// importing "database/sql/driver"?
// TODO(albrow): Could these be optimized by more directly converting between
// Go types and JavaScript types instead of using jsutil.IneffecientlyConvertX?

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gibson042/canonicaljson-go"
)

// BigInt is a wrapper around *big.Int that implements the json.Marshaler
// and json.Unmarshaler interfaces in a way that is compatible with Dexie.js
// but *does not* pad with zeroes and *does not* retain sort order.
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

// SortedBigInt is a wrapper around *big.Int that implements the json.Marshaler
// and json.Unmarshaler interfaces in a way that is compatible with Dexie.js and
// retains sort order by padding with zeroes.
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

// Order is the SQL database representation a 0x order along with some relevant metadata.
type Order struct {
	Hash                     common.Hash    `json:"hash"`
	ChainID                  *SortedBigInt  `json:"chainID"`
	ExchangeAddress          common.Address `json:"exchangeAddress"`
	MakerAddress             common.Address `json:"makerAddress"`
	MakerAssetData           []byte         `json:"makerAssetData"`
	MakerFeeAssetData        []byte         `json:"makerFeeAssetData"`
	MakerAssetAmount         *SortedBigInt  `json:"makerAssetAmount"`
	MakerFee                 *SortedBigInt  `json:"makerFee"`
	TakerAddress             common.Address `json:"takerAddress"`
	TakerAssetData           []byte         `json:"takerAssetData"`
	TakerFeeAssetData        []byte         `json:"takerFeeAssetData"`
	TakerAssetAmount         *SortedBigInt  `json:"takerAssetAmount"`
	TakerFee                 *SortedBigInt  `json:"takerFee"`
	SenderAddress            common.Address `json:"senderAddress"`
	FeeRecipientAddress      common.Address `json:"feeRecipientAddress"`
	ExpirationTimeSeconds    *SortedBigInt  `json:"expirationTimeSeconds"`
	Salt                     *SortedBigInt  `json:"salt"`
	Signature                []byte         `json:"signature"`
	LastUpdated              time.Time      `json:"lastUpdated"`
	FillableTakerAssetAmount *SortedBigInt  `json:"fillableTakerAssetAmount"`
	IsRemoved                bool           `json:"isRemoved"`
	IsPinned                 bool           `json:"isPinned"`
	ParsedMakerAssetData     string         `json:"parsedMakerAssetData"`
	ParsedMakerFeeAssetData  string         `json:"parsedMakerFeeAssetData"`
}

type MiniHeader struct {
	Hash      common.Hash   `json:"hash"`
	Parent    common.Hash   `json:"parent"`
	Number    *SortedBigInt `json:"number"`
	Timestamp time.Time     `json:"timestamp"`
	Logs      string        `json:"logs"`
}

type Metadata struct {
	EthereumChainID                   int           `json:"ethereumChainID"`
	MaxExpirationTime                 *SortedBigInt `json:"maxExpirationTime"`
	EthRPCRequestsSentInCurrentUTCDay int           `json:"ethRPCRequestsSentInCurrentUTCDay"`
	StartOfCurrentUTCDay              time.Time     `json:"startOfCurrentUTCDay"`
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

func OrdersFromCommonType(orders []*types.OrderWithMetadata) []*Order {
	result := make([]*Order, len(orders))
	for i, order := range orders {
		result[i] = OrderFromCommonType(order)
	}
	return result
}

func ParsedAssetDataToCommonType(parsedAssetData string) []*types.SingleAssetData {
	if parsedAssetData == "" {
		return nil
	}
	var dexieAssetDatas []*SingleAssetData
	_ = json.Unmarshal([]byte(parsedAssetData), &dexieAssetDatas)
	result := make([]*types.SingleAssetData, len(dexieAssetDatas))
	for i, singleAssetData := range dexieAssetDatas {
		result[i] = SingleAssetDataToCommonType(singleAssetData)
	}
	return result
}

func ParsedAssetDataFromCommonType(parsedAssetData []*types.SingleAssetData) string {
	dexieAssetDatas := ParsedAssetData(make([]*SingleAssetData, len(parsedAssetData)))
	for i, singleAssetData := range parsedAssetData {
		dexieAssetDatas[i] = SingleAssetDataFromCommonType(singleAssetData)
	}
	jsonAssetDatas, _ := canonicaljson.Marshal(dexieAssetDatas)
	return string(jsonAssetDatas)
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
		Logs:      EventLogsToCommonType(miniHeader.Logs),
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
		Logs:      EventLogsFromCommonType(miniHeader.Logs),
	}
}

func MiniHeadersToCommonType(miniHeaders []*MiniHeader) []*types.MiniHeader {
	result := make([]*types.MiniHeader, len(miniHeaders))
	for i, miniHeader := range miniHeaders {
		result[i] = MiniHeaderToCommonType(miniHeader)
	}
	return result
}

func MiniHeadersFromCommonType(miniHeaders []*types.MiniHeader) []*MiniHeader {
	result := make([]*MiniHeader, len(miniHeaders))
	for i, miniHeader := range miniHeaders {
		result[i] = MiniHeaderFromCommonType(miniHeader)
	}
	return result
}

func EventLogsToCommonType(eventLogs string) []ethtypes.Log {
	var result []ethtypes.Log
	_ = json.Unmarshal([]byte(eventLogs), &result)
	return result
}

func EventLogsFromCommonType(eventLogs []ethtypes.Log) string {
	result, _ := json.Marshal(eventLogs)
	return string(result)
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
