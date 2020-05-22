// +build js,wasm

package db

import (
	"context"
	"errors"
	"path/filepath"
	"syscall/js"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

var _ Database = (*DB)(nil)

func TestOptions() *Options {
	return &Options{
		DriverName:     "dexie",
		DataSourceName: filepath.Join("mesh_testing", uuid.New().String()),
		MaxOrders:      100,
		MaxMiniHeaders: 20,
	}
}

// New creates a new connection to the database. The connection will be automatically closed
// when the given context is canceled.
func New(ctx context.Context, opts *Options) (*DB, error) {
	return nil, errors.New("not yet implemented")
}

type DB struct {
	dexie *js.Value
}

func (db *DB) AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error) {
	return nil, nil, errors.New("not yet implemented")
}

func (db *DB) GetOrder(hash common.Hash) (*types.OrderWithMetadata, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) FindOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) CountOrders(opts *OrderQuery) (int, error) {
	return 0, errors.New("not yet implemented")
}

func (db *DB) DeleteOrder(hash common.Hash) error {
	return errors.New("not yet implemented")
}

func (db *DB) DeleteOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) error {
	return errors.New("not yet implemented")
}

func (db *DB) AddMiniHeaders(miniHeaders []*types.MiniHeader) (added []*types.MiniHeader, removed []*types.MiniHeader, err error) {
	return nil, nil, errors.New("not yet implemented")
}

func (db *DB) GetMiniHeader(hash common.Hash) (*types.MiniHeader, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) FindMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) DeleteMiniHeader(hash common.Hash) error {
	return errors.New("not yet implemented")
}

func (db *DB) DeleteMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) GetMetadata() (*types.Metadata, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) SaveMetadata(metadata *types.Metadata) error {
	return errors.New("not yet implemented")
}

func (db *DB) UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) error {
	return errors.New("not yet implemented")
}
