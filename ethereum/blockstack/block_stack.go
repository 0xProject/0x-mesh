package blockstack

import (
	"sync"

	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/meshdb"
)

// BlockStack allows performing basic stack operations on a stack of meshdb.MiniHeaders.
type BlockStack struct {
	// TODO(albrow): Use Transactions when db supports them instead of a mutex
	// here. There are cases where we need to make sure no modifications are made
	// to the database in between a read/write or read/delete.
	mut    sync.Mutex
	meshDB *meshdb.MeshDB
	limit  int
}

// New instantiates a new stack with the specified size limit. Once the size limit
// is reached, adding additional blocks will evict the deepest block.
func New(meshDB *meshdb.MeshDB, limit int) BlockStack {
	return BlockStack{
		meshDB: meshDB,
		limit:  limit,
	}
}

// Pop removes and returns the latest block header on the block stack. It
// returns nil if the stack is empty.
func (b BlockStack) Pop() (*blockwatch.MiniHeader, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	latestMiniHeader, err := b.meshDB.FindLatestMiniHeader()
	if err != nil {
		return nil, err
	}
	if latestMiniHeader == nil {
		return nil, nil
	}
	if err := b.meshDB.MiniHeaders.Delete(latestMiniHeader.ID()); err != nil {
		return nil, err
	}
	latest := blockwatch.MiniHeader(*latestMiniHeader)
	return &latest, nil
}

// Push pushes a block header onto the block stack. If the stack limit is
// reached, it will remove the oldest block header.
func (b BlockStack) Push(block *blockwatch.MiniHeader) error {
	b.mut.Lock()
	defer b.mut.Unlock()
	miniHeaders, err := b.meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return err
	}
	if len(miniHeaders) == b.limit {
		oldestMiniHeader := miniHeaders[0]
		if err := b.meshDB.MiniHeaders.Delete(oldestMiniHeader.ID()); err != nil {
			return err
		}
	}
	if err := b.meshDB.MiniHeaders.Insert(block); err != nil {
		return err
	}
	return nil
}

// Peek returns the latest block header from the block stack without removing
// it. It returns nil if the stack is empty.
func (b BlockStack) Peek() (*blockwatch.MiniHeader, error) {
	latestMiniHeader, err := b.meshDB.FindLatestMiniHeader()
	if err != nil {
		return nil, nil
	}
	latest := blockwatch.MiniHeader(*latestMiniHeader)
	return &latest, nil
}

// Inspect returns all the block headers currently on the stack
func (b BlockStack) Inspect() ([]*blockwatch.MiniHeader, error) {
	miniHeaders, err := b.meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	finalHeaders := []*blockwatch.MiniHeader{}
	for _, miniHeader := range miniHeaders {
		header := blockwatch.MiniHeader(*miniHeader)
		finalHeaders = append(finalHeaders, &header)
	}
	return finalHeaders, nil
}
