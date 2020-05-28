package db

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db/sqltypes"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gibson042/canonicaljson-go"
)

const (
	// The default miniHeaderRetentionLimit used by Mesh. This default only gets overwritten in tests.
	defaultMiniHeaderRetentionLimit = 20
	// The maximum MiniHeaders to query per page when deleting MiniHeaders
	miniHeadersMaxPerPage = 1000
	// The amount of time to wait before timing out when connecting to the database for the first time.
	connectTimeout = 10 * time.Second
)

var (
	ErrDBFilledWithPinnedOrders = errors.New("the database is full of pinned orders; no orders can be removed in order to make space")
	ErrMetadataAlreadyExists    = errors.New("metadata already exists in the database (use UpdateMetadata instead?)")
	ErrNotFound                 = errors.New("could not find existing model or row in database")
)

type Database interface {
	AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error)
	GetOrder(hash common.Hash) (*types.OrderWithMetadata, error)
	FindOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
	CountOrders(opts *OrderQuery) (int, error)
	DeleteOrder(hash common.Hash) error
	DeleteOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
	UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) error
	AddMiniHeaders(miniHeaders []*types.MiniHeader) (added []*types.MiniHeader, removed []*types.MiniHeader, err error)
	GetMiniHeader(hash common.Hash) (*types.MiniHeader, error)
	FindMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
	DeleteMiniHeader(hash common.Hash) error
	DeleteMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
	GetMetadata() (*types.Metadata, error)
	SaveMetadata(metadata *types.Metadata) error
	UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) error
}

type Options struct {
	DriverName     string
	DataSourceName string
	MaxOrders      int
	MaxMiniHeaders int
}

func parseOptions(opts *Options) *Options {
	finalOpts := defaultOptions()
	if opts == nil {
		return finalOpts
	}
	if opts.DataSourceName != "" {
		finalOpts.DataSourceName = opts.DataSourceName
	}
	if opts.MaxOrders != 0 {
		finalOpts.MaxOrders = opts.MaxOrders
	}
	if opts.MaxMiniHeaders != 0 {
		finalOpts.MaxMiniHeaders = opts.MaxMiniHeaders
	}
	return finalOpts
}

type SortDirection string

const (
	Ascending  SortDirection = "ASC"
	Descending SortDirection = "DESC"
)

type FilterKind string

const (
	Equal          FilterKind = "="
	NotEqual       FilterKind = "!="
	Less           FilterKind = "<"
	Greater        FilterKind = ">"
	LessOrEqual    FilterKind = "<="
	GreaterOrEqual FilterKind = ">="
	Contains       FilterKind = "CONTAINS"
	// TODO(albrow): Starts with?
)

type OrderField string

const (
	OFHash                     OrderField = "hash"
	OFChainID                  OrderField = "chainID"
	OFExchangeAddress          OrderField = "exchangeAddress"
	OFMakerAddress             OrderField = "makerAddress"
	OFMakerAssetData           OrderField = "makerAssetData"
	OFMakerFeeAssetData        OrderField = "makerFeeAssetData"
	OFMakerAssetAmount         OrderField = "makerAssetAmount"
	OFMakerFee                 OrderField = "makerFee"
	OFTakerAddress             OrderField = "takerAddress"
	OFTakerAssetData           OrderField = "takerAssetData"
	OFTakerFeeAssetData        OrderField = "takerFeeAssetData"
	OFTakerAssetAmount         OrderField = "takerAssetAmount"
	OFTakerFee                 OrderField = "takerFee"
	OFSenderAddress            OrderField = "senderAddress"
	OFFeeRecipientAddress      OrderField = "feeRecipientAddress"
	OFExpirationTimeSeconds    OrderField = "expirationTimeSeconds"
	OFSalt                     OrderField = "salt"
	OFSignature                OrderField = "signature"
	OFLastUpdated              OrderField = "lastUpdated"
	OFFillableTakerAssetAmount OrderField = "fillableTakerAssetAmount"
	OFIsRemoved                OrderField = "isRemoved"
	OFIsPinned                 OrderField = "isPinned"
	OFParsedMakerAssetData     OrderField = "parsedMakerAssetData"
	OFParsedMakerFeeAssetData  OrderField = "parsedMakerFeeAssetData"
)

type OrderQuery struct {
	Filters []OrderFilter `json:"filters"`
	Sort    []OrderSort   `json:"sort"`
	Limit   uint          `json:"limit"`
	Offset  uint          `json:"offset"`
}

type OrderSort struct {
	Field     OrderField    `json:"field"`
	Direction SortDirection `json:"direction"`
}

type OrderFilter struct {
	Field OrderField  `json:"field"`
	Kind  FilterKind  `json:"kind"`
	Value interface{} `json:"value"`
}

// MakerAssetIncludesTokenAddressAndTokenID is a helper method which returns a filter that will match orders
// that include the token address and token ID in MakerAssetData.
func MakerAssetIncludesTokenAddressAndTokenID(tokenAddress common.Address, tokenID *big.Int) OrderFilter {
	return assetDataIncludesTokenAddressAndTokenID(OFParsedMakerAssetData, tokenAddress, tokenID)
}

// MakerFeeAssetIncludesTokenAddressAndTokenID is a helper method which returns a filter that will match orders
// that include the token address and token ID in MakerFeeAssetData.
func MakerFeeAssetIncludesTokenAddressAndTokenID(tokenAddress common.Address, tokenID *big.Int) OrderFilter {
	return assetDataIncludesTokenAddressAndTokenID(OFParsedMakerFeeAssetData, tokenAddress, tokenID)
}

func assetDataIncludesTokenAddressAndTokenID(field OrderField, tokenAddress common.Address, tokenID *big.Int) OrderFilter {
	filterValueJSON, err := canonicaljson.Marshal(sqltypes.SingleAssetData{
		Address: tokenAddress,
		TokenID: sqltypes.NewBigInt(tokenID),
	})
	if err != nil {
		// big.Int and common.Address types should never return an error when marshaling to JSON
		panic(err)
	}
	return OrderFilter{
		Field: field,
		Kind:  Contains,
		Value: string(filterValueJSON),
	}
}

// MakerAssetIncludesTokenAddress is a helper method which returns a filter that will match orders
// that include the token address (and any token id, including null) in MakerAssetData.
func MakerAssetIncludesTokenAddress(tokenAddress common.Address) OrderFilter {
	return assetDataIncludesTokenAddress(OFParsedMakerAssetData, tokenAddress)
}

// MakerFeeAssetIncludesTokenAddress is a helper method which returns a filter that will match orders
// that include the token address (and any token id, including null) in MakerFeeAssetData.
func MakerFeeAssetIncludesTokenAddress(tokenAddress common.Address) OrderFilter {
	return assetDataIncludesTokenAddress(OFParsedMakerFeeAssetData, tokenAddress)
}

func assetDataIncludesTokenAddress(field OrderField, tokenAddress common.Address) OrderFilter {
	tokenAddressJSON, err := canonicaljson.Marshal(tokenAddress)
	if err != nil {
		// big.Int and common.Address types should never return an error when marshaling to JSON
		panic(err)
	}
	filterValue := fmt.Sprintf(`"address":%s`, tokenAddressJSON)
	return OrderFilter{
		Field: field,
		Kind:  Contains,
		Value: filterValue,
	}
}

type MiniHeaderField string

const (
	MFHash      MiniHeaderField = "hash"
	MFParent    MiniHeaderField = "parent"
	MFNumber    MiniHeaderField = "number"
	MFTimestamp MiniHeaderField = "timestamp"
	MFLogs      MiniHeaderField = "logs"
)

type MiniHeaderQuery struct {
	Filters []MiniHeaderFilter
	Sort    []MiniHeaderSort
	Limit   uint
	Offset  uint
}

type MiniHeaderSort struct {
	Field     MiniHeaderField
	Direction SortDirection
}

type MiniHeaderFilter struct {
	Field MiniHeaderField
	Kind  FilterKind
	Value interface{}
}

func ParseContractAddressesAndTokenIdsFromAssetData(assetData []byte, contractAddresses ethereum.ContractAddresses) ([]*types.SingleAssetData, error) {
	if len(assetData) == 0 {
		return []*types.SingleAssetData{}, nil
	}
	singleAssetDatas := []*types.SingleAssetData{}
	assetDataDecoder := zeroex.NewAssetDataDecoder()

	assetDataName, err := assetDataDecoder.GetName(assetData)
	if err != nil {
		return nil, err
	}
	switch assetDataName {
	case "ERC20Token":
		var decodedAssetData zeroex.ERC20AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := &types.SingleAssetData{
			Address: decodedAssetData.Address,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := &types.SingleAssetData{
			Address: decodedAssetData.Address,
			TokenID: decodedAssetData.TokenId,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		for _, id := range decodedAssetData.Ids {
			a := &types.SingleAssetData{
				Address: decodedAssetData.Address,
				TokenID: id,
			}
			singleAssetDatas = append(singleAssetDatas, a)
		}
	case "StaticCall":
		var decodedAssetData zeroex.StaticCallAssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		// NOTE(jalextowle): As of right now, none of the supported staticcalls
		// have important information in the StaticCallData. We choose not to add
		// `singleAssetData` because it would not be used.
	case "MultiAsset":
		var decodedAssetData zeroex.MultiAssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		for _, assetData := range decodedAssetData.NestedAssetData {
			as, err := ParseContractAddressesAndTokenIdsFromAssetData(assetData, contractAddresses)
			if err != nil {
				return nil, err
			}
			singleAssetDatas = append(singleAssetDatas, as...)
		}
	case "ERC20Bridge":
		var decodedAssetData zeroex.ERC20BridgeAssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		tokenAddress := decodedAssetData.TokenAddress
		// TODO(albrow): Update orderwatcher to account for this instead of storing
		// it in the database. This would mean we can remove contractAddresses as an
		// argument and simplify the implementation. Maybe even have the db package
		// handle parsing asset data automatically.
		// HACK(fabio): Despite Chai ERC20Bridge orders encoding the Dai address as
		// the tokenAddress, we actually want to react to the Chai token's contract
		// events, so we actually return it instead.
		if decodedAssetData.BridgeAddress == contractAddresses.ChaiBridge {
			tokenAddress = contractAddresses.ChaiToken
		}
		a := &types.SingleAssetData{
			Address: tokenAddress,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	default:
		return nil, fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
	}
	return singleAssetDatas, nil
}

func checkOrderQuery(query *OrderQuery) error {
	if query == nil {
		return nil
	}
	if query.Offset != 0 && query.Limit == 0 {
		return errors.New("can't use Offset without Limit")
	}
	return nil
}
