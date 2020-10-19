// +build js,wasm

package dexietypes

// Note(albrow): Could be optimized if needed by more directly converting between
// Go types and JavaScript types instead of using jsutil.IneffecientlyConvertX.
// The technique we used for MiniHeaders could be used in more places if needed.

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
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

func (i *SortedBigInt) FixedLengthString() string {
	// Note(albrow), strings in Dexie.js are sorted in alphanumerical order, not
	// numerical order. In order to sort by numerical order, we need to pad with
	// zeroes. The maximum length of an unsigned 256 bit integer is 80, so we
	// pad with zeroes such that the length of the number is always 80.
	return fmt.Sprintf("%080s", i.Int.String())
}

func (i *SortedBigInt) MarshalJSON() ([]byte, error) {
	if i == nil || i.Int == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(i.FixedLengthString())
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
	Version                  int            `json:"version"`
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
	IsRemoved                uint8          `json:"isRemoved"`
	IsPinned                 uint8          `json:"isPinned"`
	IsNotPinned              uint8          `json:"isNotPinned"` // Used in a compound index in queries related to max expiration time.
	IsUnfillable             uint8          `json:"isUnfillable"`
	IsExpired                uint8          `json:"isExpired"`
	ParsedMakerAssetData     string         `json:"parsedMakerAssetData"`
	ParsedMakerFeeAssetData  string         `json:"parsedMakerFeeAssetData"`
	LastValidatedBlockNumber *SortedBigInt  `json:"lastValidatedBlockNumber"`
	LastValidatedBlockHash   common.Hash    `json:"lastValidatedBlockHash"`
	KeepCancelled            uint8          `json:"keepCancelled"`
	KeepExpired              uint8          `json:"keepExpired"`
	KeepFullyFilled          uint8          `json:"keepFullyFilled"`
	KeepUnfunded             uint8          `json:"keepUnfunded"`
}

type Metadata struct {
	EthereumChainID                   int       `json:"ethereumChainID"`
	EthRPCRequestsSentInCurrentUTCDay int       `json:"ethRPCRequestsSentInCurrentUTCDay"`
	StartOfCurrentUTCDay              time.Time `json:"startOfCurrentUTCDay"`
}

func OrderToCommonType(order *Order) (*types.OrderWithMetadata, error) {
	if order == nil {
		return nil, nil
	}
	switch order.Version {
	case 3:
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
			IsRemoved:                order.IsRemoved == 1,
			IsPinned:                 order.IsPinned == 1,
			IsUnfillable:             order.IsUnfillable == 1,
			IsExpired:                order.IsExpired == 1,
			ParsedMakerAssetData:     ParsedAssetDataToCommonType(order.ParsedMakerAssetData),
			ParsedMakerFeeAssetData:  ParsedAssetDataToCommonType(order.ParsedMakerFeeAssetData),
			LastValidatedBlockNumber: order.LastValidatedBlockNumber.Int,
			LastValidatedBlockHash:   order.LastValidatedBlockHash,
			KeepCancelled:            order.KeepCancelled == 1,
			KeepExpired:              order.KeepExpired == 1,
			KeepFullyFilled:          order.KeepFullyFilled == 1,
			KeepUnfunded:             order.KeepUnfunded == 1,
		}, nil
	default:
		return nil, errors.New("Unknown order version stored in database")
	}
}

func OrderFromCommonType(order *types.OrderWithMetadata) *Order {
	if order == nil {
		return nil
	}
	return &Order{
		Hash:                     order.Hash,
		Version:                  3,
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
		IsRemoved:                BoolToUint8(order.IsRemoved),
		IsPinned:                 BoolToUint8(order.IsPinned),
		IsNotPinned:              BoolToUint8(!order.IsPinned),
		IsUnfillable:             BoolToUint8(order.IsUnfillable),
		IsExpired:                BoolToUint8(order.IsExpired),
		ParsedMakerAssetData:     ParsedAssetDataFromCommonType(order.ParsedMakerAssetData),
		ParsedMakerFeeAssetData:  ParsedAssetDataFromCommonType(order.ParsedMakerFeeAssetData),
		LastValidatedBlockNumber: NewSortedBigInt(order.LastValidatedBlockNumber),
		LastValidatedBlockHash:   order.LastValidatedBlockHash,
		KeepCancelled:            BoolToUint8(order.KeepCancelled),
		KeepExpired:              BoolToUint8(order.KeepExpired),
		KeepFullyFilled:          BoolToUint8(order.KeepFullyFilled),
		KeepUnfunded:             BoolToUint8(order.KeepUnfunded),
	}
}

func OrdersToCommonType(orders []*Order) ([]*types.OrderWithMetadata, error) {
	result := make([]*types.OrderWithMetadata, len(orders))
	for i, order := range orders {
		var err error
		result[i], err = OrderToCommonType(order)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
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

func MiniHeaderToCommonType(miniHeader js.Value) *types.MiniHeader {
	if jsutil.IsNullOrUndefined(miniHeader) {
		return nil
	}
	number, ok := ethmath.ParseBig256(miniHeader.Get("number").String())
	if !ok {
		panic(errors.New("could not convert number to uint64"))
	}
	timestamp, err := time.Parse(time.RFC3339Nano, miniHeader.Get("timestamp").String())
	if err != nil {
		panic(errors.New("could not convert timestamp: " + err.Error()))
	}
	return &types.MiniHeader{
		Hash:      common.HexToHash(miniHeader.Get("hash").String()),
		Parent:    common.HexToHash(miniHeader.Get("parent").String()),
		Number:    number,
		Timestamp: timestamp,
		Logs:      EventLogsToCommonType(miniHeader.Get("logs")),
	}
}

func MiniHeaderFromCommonType(miniHeader *types.MiniHeader) js.Value {
	if miniHeader == nil {
		return js.Null()
	}
	return js.ValueOf(
		map[string]interface{}{
			"hash":      miniHeader.Hash.Hex(),
			"parent":    miniHeader.Parent.Hex(),
			"number":    NewSortedBigInt(miniHeader.Number).FixedLengthString(),
			"timestamp": miniHeader.Timestamp.Format(time.RFC3339Nano),
			"logs":      EventLogsFromCommonType(miniHeader.Logs),
		},
	)
}

func MiniHeadersToCommonType(miniHeaders js.Value) []*types.MiniHeader {
	result := make([]*types.MiniHeader, miniHeaders.Length())
	for i := range result {
		result[i] = MiniHeaderToCommonType(miniHeaders.Index(i))
	}
	return result
}

func MiniHeadersFromCommonType(miniHeaders []*types.MiniHeader) js.Value {
	result := make([]interface{}, len(miniHeaders))
	for i, miniHeader := range miniHeaders {
		result[i] = MiniHeaderFromCommonType(miniHeader)
	}
	return js.ValueOf(result)
}

func EventLogsToCommonType(eventLogs js.Value) []ethtypes.Log {
	var result []ethtypes.Log
	buf := make([]byte, eventLogs.Get("length").Int())
	js.CopyBytesToGo(buf, eventLogs)
	if err := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&result); err != nil {
		panic(err)
	}
	return result
}

func EventLogsFromCommonType(eventLogs []ethtypes.Log) js.Value {
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(eventLogs); err != nil {
		panic(err)
	}
	result := js.Global().Get("Uint8Array").New(len(buf.Bytes()))
	js.CopyBytesToJS(result, buf.Bytes())
	return result
}

func MetadataToCommonType(metadata *Metadata) *types.Metadata {
	if metadata == nil {
		return nil
	}
	return &types.Metadata{
		EthereumChainID:                   metadata.EthereumChainID,
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
		EthRPCRequestsSentInCurrentUTCDay: metadata.EthRPCRequestsSentInCurrentUTCDay,
		StartOfCurrentUTCDay:              metadata.StartOfCurrentUTCDay,
	}
}

func BoolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
