package db

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
)

// DB instantiates the DB connection and creates all the collections used by the application
type DB struct {
	// database                 *db.DB
	// metadata                 *MetadataCollection
	// MiniHeaders              *MiniHeadersCollection
	// Orders                   *OrdersCollection
	// MiniHeaderRetentionLimit int
}

// New instantiates a new MeshDB instance
func New(path string, contractAddresses ethereum.ContractAddresses) (*DB, error) {
	return nil, errors.New("Not yet implemented")
}

func (m *DB) Close() error {
	return errors.New("Not yet implemented")
}

// FindAllMiniHeadersSortedByNumber returns all MiniHeaders sorted in ascending block number order
func (m *DB) FindAllMiniHeadersSortedByNumber() ([]*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// FindLatestMiniHeader returns the latest MiniHeader (i.e. the one with the
// largest block number). It returns nil, MiniHeaderCollectionEmptyError if there
// are no MiniHeaders in the database.
func (m *DB) FindLatestMiniHeader() (*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// FindMiniHeaderByBlockNumber returns the MiniHeader with the specified block number
func (m *DB) FindMiniHeaderByBlockNumber(blockNumber *big.Int) (*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// UpdateMiniHeaderRetentionLimit updates the MiniHeaderRetentionLimit. This is only used by tests in order
// to set the retention limit to a smaller size, making the tests shorter in length
func (m *DB) UpdateMiniHeaderRetentionLimit(limit int) error {
	return errors.New("Not yet implemented")
}

// PruneMiniHeadersAboveRetentionLimit prunes miniHeaders from the DB that are above the retention limit
func (m *DB) PruneMiniHeadersAboveRetentionLimit() error {
	return errors.New("Not yet implemented")
}

// ClearAllMiniHeaders removes all stored MiniHeaders from the database.
func (m *DB) ClearAllMiniHeaders() error {
	return errors.New("Not yet implemented")
}

// ClearOldMiniHeaders removes all stored MiniHeaders with a block number less then
// the given minBlockNumber.
func (m *DB) ClearOldMiniHeaders(minBlockNumber *big.Int) error {
	return errors.New("Not yet implemented")
}

// FindOrdersByMakerAddress finds all orders belonging to a particular maker address
func (m *DB) FindOrdersByMakerAddress(makerAddress common.Address) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressTokenAddressAndTokenID finds all orders belonging to a particular maker
// address where makerAssetData encodes for a particular token contract and optionally a token ID
func (m *DB) FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressMakerFeeAssetAddressTokenID finds all orders belonging to
// a particular maker address where makerFeeAssetData encodes for a particular
// token contract and optionally a token ID. To find orders without a maker fee,
// use constants.NullAddress for makerFeeAssetAddress.
func (m *DB) FindOrdersByMakerAddressMakerFeeAssetAddressAndTokenID(makerAddress, makerFeeAssetAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressAndMaxSalt finds all orders belonging to a particular maker address that
// also have a salt value less then or equal to X
func (m *DB) FindOrdersByMakerAddressAndMaxSalt(makerAddress common.Address, salt *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersLastUpdatedBefore finds all orders where the LastUpdated time is less
// than X
func (m *DB) FindOrdersLastUpdatedBefore(lastUpdated time.Time) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindRemovedOrders finds all orders that have been flagged for removal
func (m *DB) FindRemovedOrders() ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// GetMetadata returns the metadata (or a db.NotFoundError if no metadata has been found).
func (m *DB) GetMetadata() (*Metadata, error) {
	return nil, errors.New("Not yet implemented")
}

// SaveMetadata inserts the metadata into the database, overwriting any existing
// metadata.
func (m *DB) SaveMetadata(metadata *Metadata) error {
	return errors.New("Not yet implemented")
}

// UpdateMetadata updates the metadata in the database via a transaction. It
// accepts a callback function which will be provided with the old metadata and
// should return the new metadata to save.
func (m *DB) UpdateMetadata(updater func(oldmetadata Metadata) (newMetadata Metadata)) error {
	return errors.New("Not yet implemented")
}

type singleAssetData struct {
	Address common.Address
	TokenID *big.Int
}

func parseContractAddressesAndTokenIdsFromAssetData(assetData []byte, contractAddresses ethereum.ContractAddresses) ([]singleAssetData, error) {
	singleAssetDatas := []singleAssetData{}
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
		a := singleAssetData{
			Address: decodedAssetData.Address,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := singleAssetData{
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
			a := singleAssetData{
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
			as, err := parseContractAddressesAndTokenIdsFromAssetData(assetData, contractAddresses)
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
		// HACK(fabio): Despite Chai ERC20Bridge orders encoding the Dai address as
		// the tokenAddress, we actually want to react to the Chai token's contract
		// events, so we actually return it instead.
		if decodedAssetData.BridgeAddress == contractAddresses.ChaiBridge {
			tokenAddress = contractAddresses.ChaiToken
		}
		a := singleAssetData{
			Address: tokenAddress,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	default:
		return nil, fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
	}
	return singleAssetDatas, nil
}

func uint256ToConstantLengthBytes(v *big.Int) []byte {
	return []byte(fmt.Sprintf("%080s", v.String()))
}

// TrimOrdersByExpirationTime removes existing orders with the highest
// expiration time until the number of remaining orders is <= targetMaxOrders.
// It returns any orders that were removed and the new max expiration time that
// can be used to eliminate incoming orders that expire too far in the future.
func (m *DB) TrimOrdersByExpirationTime(targetMaxOrders int) (newMaxExpirationTime *big.Int, removedOrders []*Order, err error) {
	return nil, nil, errors.New("Not yet implemented")
}

// CountPinnedOrders returns the number of pinned orders.
func (m *DB) CountPinnedOrders() (int, error) {
	return 0, errors.New("Not yet implemented")
}
