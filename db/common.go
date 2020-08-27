package db

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gibson042/canonicaljson-go"
	ds "github.com/ipfs/go-datastore"
)

const (
	// The amount of time to wait before timing out when connecting to the database for the first time.
	connectTimeout = 10 * time.Second
)

var (
	ErrDBFilledWithPinnedOrders = errors.New("the database is full of pinned orders; no orders can be removed in order to make space")
	ErrMetadataAlreadyExists    = errors.New("metadata already exists in the database (use UpdateMetadata instead?)")
	ErrNotFound                 = errors.New("could not find existing model or row in database")
	ErrClosed                   = errors.New("database is already closed")
)

type Database interface {
	AddOrders(orders []*types.OrderWithMetadata) (alreadyStored []common.Hash, added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error)
	GetOrder(hash common.Hash) (*types.OrderWithMetadata, error)
	GetOrderStatuses(hashes []common.Hash) (statuses []*StoredOrderStatus, err error)
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
	ResetMiniHeaders(miniHeaders []*types.MiniHeader) error
	GetMetadata() (*types.Metadata, error)
	SaveMetadata(metadata *types.Metadata) error
	UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) error
	PeerStore() ds.Batching
	DHTStore() ds.Batching
}

type Options struct {
	DriverName     string `json:"driverName"`
	DataSourceName string `json:"dataSourceName"`
	MaxOrders      int    `json:"maxOrders"`
	MaxMiniHeaders int    `json:"maxMiniHeaders"`
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
	OFLastValidatedBlockNumber OrderField = "lastValidatedBlockNumber"
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

type StoredOrderStatus struct {
	IsStored                 bool     `json:"isStored"`
	IsMarkedRemoved          bool     `json:"isMarkedRemoved"`
	FillableTakerAssetAmount *big.Int `json:"fillableTakerAssetAmount"`
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
	Filters []MiniHeaderFilter `json:"filters"`
	Sort    []MiniHeaderSort   `json:"sort"`
	Limit   uint               `json:"limit"`
	Offset  uint               `json:"offset"`
}

type MiniHeaderSort struct {
	Field     MiniHeaderField `json:"field"`
	Direction SortDirection   `json:"direction"`
}

type MiniHeaderFilter struct {
	Field MiniHeaderField `json:"field"`
	Kind  FilterKind      `json:"kind"`
	Value interface{}     `json:"value"`
}

// GetOldestMiniHeader is a helper method for getting the oldest MiniHeader.
// It returns ErrNotFound if there are no MiniHeaders in the database.
func (db *DB) GetOldestMiniHeader() (*types.MiniHeader, error) {
	oldestMiniHeaders, err := db.FindMiniHeaders(&MiniHeaderQuery{
		Sort: []MiniHeaderSort{
			{
				Field:     MFNumber,
				Direction: Ascending,
			},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(oldestMiniHeaders) == 0 {
		return nil, ErrNotFound
	}
	return oldestMiniHeaders[0], nil
}

// GetLatestMiniHeader is a helper method for getting the latest MiniHeader.
// It returns ErrNotFound if there are no MiniHeaders in the database.
func (db *DB) GetLatestMiniHeader() (*types.MiniHeader, error) {
	latestMiniHeaders, err := db.FindMiniHeaders(&MiniHeaderQuery{
		Sort: []MiniHeaderSort{
			{
				Field:     MFNumber,
				Direction: Descending,
			},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(latestMiniHeaders) == 0 {
		return nil, ErrNotFound
	}
	return latestMiniHeaders[0], nil
}

// GetCurrentMaxExpirationTime returns the maximum expiration time for non-pinned orders
// stored in the database. If there are no non-pinned orders in the database, it returns
// constants.UnlimitedExpirationTime.
func (db *DB) GetCurrentMaxExpirationTime() (*big.Int, error) {
	// Note(albrow): We don't include pinned orders because they are
	// never removed due to exceeding the max expiration time.
	ordersWithLongestExpirationTime, err := db.FindOrders(&OrderQuery{
		Filters: []OrderFilter{
			{
				Field: OFIsPinned,
				Kind:  Equal,
				Value: false,
			},
		},
		Sort: []OrderSort{
			{
				Field:     OFExpirationTimeSeconds,
				Direction: Descending,
			},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(ordersWithLongestExpirationTime) == 0 {
		return constants.UnlimitedExpirationTime, nil
	}
	return ordersWithLongestExpirationTime[0].ExpirationTimeSeconds, nil
}

func ParseContractAddressesAndTokenIdsFromAssetData(assetDataDecoder *zeroex.AssetDataDecoder, assetData []byte, contractAddresses ethereum.ContractAddresses) ([]*types.SingleAssetData, error) {
	if len(assetData) == 0 {
		return []*types.SingleAssetData{}, nil
	}
	singleAssetDatas := []*types.SingleAssetData{}

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
			as, err := ParseContractAddressesAndTokenIdsFromAssetData(assetDataDecoder, assetData, contractAddresses)
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
