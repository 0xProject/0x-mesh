package db

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB instantiates the DB connection and creates all the collections used by the application
type DB struct {
	ctx                      context.Context
	sqldb                    *sqlx.DB
	MiniHeaderRetentionLimit int
}

// New creates a new connection to the database. The connection will be automatically closed
// when the given context is canceled.
func New(ctx context.Context, path string) (*DB, error) {
	connectCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if err := os.MkdirAll(path, os.ModePerm); err != nil && err != os.ErrExist {
		return nil, err
	}

	db, err := sqlx.ConnectContext(connectCtx, "sqlite3", filepath.Join(path, "db.sqlite"))
	if err != nil {
		return nil, err
	}

	// Automatically close the database connection when the context is canceled.
	go func() {
		select {
		case <-ctx.Done():
			_ = db.Close()
		}
	}()

	return &DB{
		ctx:   ctx,
		sqldb: db,
	}, nil
}

// TODO(albrow): Use a proper migration tool.
const schema = `
CREATE TABLE IF NOT EXISTS orders (
	hash                     TEXT UNIQUE NOT NULL,
	chainID                  NUMERIC(78, 0) NOT NULL,
	exchangeAddress          TEXT NOT NULL,
	makerAddress             TEXT NOT NULL,
	makerAssetData           TEXT NOT NULL,
	makerFeeAssetData        TEXT NOT NULL,
	makerAssetAmount         NUMERIC(78, 0) NOT NULL,
	makerFee                 NUMERIC(78, 0) NOT NULL,
	takerAddress             TEXT NOT NULL,
	takerAssetData           TEXT NOT NULL,
	takerFeeAssetData        TEXT NOT NULL,
	takerAssetAmount         NUMERIC(78, 0) NOT NULL,
	takerFee                 NUMERIC(78, 0) NOT NULL,
	senderAddress            TEXT NOT NULL,
	feeRecipientAddress      TEXT NOT NULL,
	expirationTimeSeconds    NUMERIC(78, 0) NOT NULL,
	salt                     NUMERIC(78, 0) NOT NULL,
	signature                TEXT NOT NULL,
	lastUpdated              DATETIME NOT NULL,
	fillableTakerAssetAmount NUMERIC(78, 0) NOT NULL,
	isRemoved                BOOLEAN NOT NULL,
	isPinned                 BOOLEAN NOT NULL
);
`

// TODO(albrow): Used prepared statement for inserts.
const insertQuery = `INSERT OR IGNORE INTO orders (
		hash,
		chainID,
		exchangeAddress,
		makerAddress,
		makerAssetData,
		makerFeeAssetData,
		makerAssetAmount,
		makerFee,
		takerAddress,
		takerAssetData,
		takerFeeAssetData,
		takerAssetAmount,
		takerFee,
		senderAddress,
		feeRecipientAddress,
		expirationTimeSeconds,
		salt,
		signature,
		lastUpdated,
		fillableTakerAssetAmount,
		isRemoved,
		isPinned
	) VALUES (
		:hash,
		:chainID,
		:exchangeAddress,
		:makerAddress,
		:makerAssetData,
		:makerFeeAssetData,
		:makerAssetAmount,
		:makerFee,
		:takerAddress,
		:takerAssetData,
		:takerFeeAssetData,
		:takerAssetAmount,
		:takerFee,
		:senderAddress,
		:feeRecipientAddress,
		:expirationTimeSeconds,
		:salt,
		:signature,
		:lastUpdated,
		:fillableTakerAssetAmount,
		:isRemoved,
		:isPinned
	)
`

func (db *DB) migrate() error {
	_, err := db.sqldb.ExecContext(db.ctx, schema)
	return err
}

func (db *DB) Close() error {
	panic(errors.New("Not implemented. Cancel the context instead."))
}

func (db *DB) AddOrders(orders []*Order) (added []*Order, removed []*Order, err error) {
	txn, err := db.sqldb.BeginTxx(db.ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	for _, order := range orders {
		result, err := txn.NamedExecContext(db.ctx, insertQuery, order)
		if err != nil {
			fmt.Printf("%T %s\n", err, err)
			spew.Dump(err)
			return nil, nil, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return nil, nil, err
		}
		if affected > 0 {
			added = append(added, order)
		}
	}
	if err := txn.Commit(); err != nil {
		return nil, nil, err
	}

	// TODO(albrow): Remove orders with longest expiration time.
	// TODO(albrow): Return appropriate values for added, removed.
	return added, nil, nil
}

func (db *DB) FindOrder(hash common.Hash) (*Order, error) {
	var order Order
	if err := db.sqldb.GetContext(db.ctx, &order, "SELECT * FROM orders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return nil, err
	}
	return &order, nil
}

// FindAllMiniHeadersSortedByNumber returns all MiniHeaders sorted in ascending block number order
func (db *DB) FindAllMiniHeadersSortedByNumber() ([]*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// FindLatestMiniHeader returns the latest MiniHeader (i.e. the one with the
// largest block number). It returns nil, MiniHeaderCollectionEmptyError if there
// are no MiniHeaders in the database.
func (db *DB) FindLatestMiniHeader() (*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// FindMiniHeaderByBlockNumber returns the MiniHeader with the specified block number
func (db *DB) FindMiniHeaderByBlockNumber(blockNumber *big.Int) (*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// UpdateMiniHeaderRetentionLimit updates the MiniHeaderRetentionLimit. This is only used by tests in order
// to set the retention limit to a smaller size, making the tests shorter in length
func (db *DB) UpdateMiniHeaderRetentionLimit(limit int) error {
	return errors.New("Not yet implemented")
}

// PruneMiniHeadersAboveRetentionLimit prunes miniHeaders from the DB that are above the retention limit
func (db *DB) PruneMiniHeadersAboveRetentionLimit() error {
	return errors.New("Not yet implemented")
}

// ClearAllMiniHeaders removes all stored MiniHeaders from the database.
func (db *DB) ClearAllMiniHeaders() error {
	return errors.New("Not yet implemented")
}

// ClearOldMiniHeaders removes all stored MiniHeaders with a block number less then
// the given minBlockNumber.
func (db *DB) ClearOldMiniHeaders(minBlockNumber *big.Int) error {
	return errors.New("Not yet implemented")
}

// FindOrdersByMakerAddress finds all orders belonging to a particular maker address
func (db *DB) FindOrdersByMakerAddress(makerAddress common.Address) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressTokenAddressAndTokenID finds all orders belonging to a particular maker
// address where makerAssetData encodes for a particular token contract and optionally a token ID
func (db *DB) FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressMakerFeeAssetAddressTokenID finds all orders belonging to
// a particular maker address where makerFeeAssetData encodes for a particular
// token contract and optionally a token ID. To find orders without a maker fee,
// use constants.NullAddress for makerFeeAssetAddress.
func (db *DB) FindOrdersByMakerAddressMakerFeeAssetAddressAndTokenID(makerAddress, makerFeeAssetAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressAndMaxSalt finds all orders belonging to a particular maker address that
// also have a salt value less then or equal to X
func (db *DB) FindOrdersByMakerAddressAndMaxSalt(makerAddress common.Address, salt *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersLastUpdatedBefore finds all orders where the LastUpdated time is less
// than X
func (db *DB) FindOrdersLastUpdatedBefore(lastUpdated time.Time) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindRemovedOrders finds all orders that have been flagged for removal
func (db *DB) FindRemovedOrders() ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// GetMetadata returns the metadata (or a db.NotFoundError if no metadata has been found).
func (db *DB) GetMetadata() (*Metadata, error) {
	return nil, errors.New("Not yet implemented")
}

// SaveMetadata inserts the metadata into the database, overwriting any existing
// metadata.
func (db *DB) SaveMetadata(metadata *Metadata) error {
	return errors.New("Not yet implemented")
}

// UpdateMetadata updates the metadata in the database via a transaction. It
// accepts a callback function which will be provided with the old metadata and
// should return the new metadata to save.
func (db *DB) UpdateMetadata(updater func(oldmetadata Metadata) (newMetadata Metadata)) error {
	return errors.New("Not yet implemented")
}

// type singleAssetData struct {
// 	Address common.Address
// 	TokenID *big.Int
// }

// func parseContractAddressesAndTokenIdsFromAssetData(assetData []byte, contractAddresses ethereum.ContractAddresses) ([]singleAssetData, error) {
// 	singleAssetDatas := []singleAssetData{}
// 	assetDataDecoder := zeroex.NewAssetDataDecoder()

// 	assetDataName, err := assetDataDecoder.GetName(assetData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	switch assetDataName {
// 	case "ERC20Token":
// 		var decodedAssetData zeroex.ERC20AssetData
// 		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
// 		if err != nil {
// 			return nil, err
// 		}
// 		a := singleAssetData{
// 			Address: decodedAssetData.Address,
// 		}
// 		singleAssetDatas = append(singleAssetDatas, a)
// 	case "ERC721Token":
// 		var decodedAssetData zeroex.ERC721AssetData
// 		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
// 		if err != nil {
// 			return nil, err
// 		}
// 		a := singleAssetData{
// 			Address: decodedAssetData.Address,
// 			TokenID: decodedAssetData.TokenId,
// 		}
// 		singleAssetDatas = append(singleAssetDatas, a)
// 	case "ERC1155Assets":
// 		var decodedAssetData zeroex.ERC1155AssetData
// 		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, id := range decodedAssetData.Ids {
// 			a := singleAssetData{
// 				Address: decodedAssetData.Address,
// 				TokenID: id,
// 			}
// 			singleAssetDatas = append(singleAssetDatas, a)
// 		}
// 	case "StaticCall":
// 		var decodedAssetData zeroex.StaticCallAssetData
// 		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// NOTE(jalextowle): As of right now, none of the supported staticcalls
// 		// have important information in the StaticCallData. We choose not to add
// 		// `singleAssetData` because it would not be used.
// 	case "MultiAsset":
// 		var decodedAssetData zeroex.MultiAssetData
// 		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, assetData := range decodedAssetData.NestedAssetData {
// 			as, err := parseContractAddressesAndTokenIdsFromAssetData(assetData, contractAddresses)
// 			if err != nil {
// 				return nil, err
// 			}
// 			singleAssetDatas = append(singleAssetDatas, as...)
// 		}
// 	case "ERC20Bridge":
// 		var decodedAssetData zeroex.ERC20BridgeAssetData
// 		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tokenAddress := decodedAssetData.TokenAddress
// 		// HACK(fabio): Despite Chai ERC20Bridge orders encoding the Dai address as
// 		// the tokenAddress, we actually want to react to the Chai token's contract
// 		// events, so we actually return it instead.
// 		if decodedAssetData.BridgeAddress == contractAddresses.ChaiBridge {
// 			tokenAddress = contractAddresses.ChaiToken
// 		}
// 		a := singleAssetData{
// 			Address: tokenAddress,
// 		}
// 		singleAssetDatas = append(singleAssetDatas, a)
// 	default:
// 		return nil, fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
// 	}
// 	return singleAssetDatas, nil
// }

func uint256ToConstantLengthBytes(v *big.Int) []byte {
	return []byte(fmt.Sprintf("%080s", v.String()))
}

// TrimOrdersByExpirationTime removes existing orders with the highest
// expiration time until the number of remaining orders is <= targetMaxOrders.
// It returns any orders that were removed and the new max expiration time that
// can be used to eliminate incoming orders that expire too far in the future.
func (db *DB) TrimOrdersByExpirationTime(targetMaxOrders int) (newMaxExpirationTime *big.Int, removedOrders []*Order, err error) {
	return nil, nil, errors.New("Not yet implemented")
}

// CountPinnedOrders returns the number of pinned orders.
func (db *DB) CountPinnedOrders() (int, error) {
	return 0, errors.New("Not yet implemented")
}
