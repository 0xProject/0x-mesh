package dbstack

import (
	"fmt"
	"sync"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/ethereum/simplestack"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum/common"
)

// DBStack is an in-memory stack that can be sync'd with the DB
type DBStack struct {
	meshDB      *meshdb.MeshDB
	mu          sync.Mutex
	simpleStack *simplestack.SimpleStack
}

// New instantiates a new DBStack. DBStack is go-routine safe.
func New(meshDB *meshdb.MeshDB, retentionLimit int) (*DBStack, error) {
	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	d := &DBStack{
		meshDB:      meshDB,
		simpleStack: simplestack.NewSimpleStack(retentionLimit, miniHeaders),
	}
	return d, nil
}

// Peek returns the top of the stack
func (d *DBStack) Peek() (*miniheader.MiniHeader, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.simpleStack.Peek()
}

// Pop returns the top of the stack and removes it from the stack and backing DB
func (d *DBStack) Pop() (*miniheader.MiniHeader, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.simpleStack.Pop()
}

// Push adds a miniheader.MiniHeader to the stack
func (d *DBStack) Push(miniHeader *miniheader.MiniHeader) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.simpleStack.Push(miniHeader)
}

// PeekAll returns all the miniHeaders currently in the stack
func (d *DBStack) PeekAll() ([]*miniheader.MiniHeader, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.simpleStack.PeekAll()
}

// Clear removes all items from the stack and the backing DB
func (d *DBStack) Clear() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.simpleStack.Clear(); err != nil {
		return err
	}
	return d.meshDB.ClearAllMiniHeaders()
}

// Reset resets the in-memory stack with the contents from the latest checkpoint
func (d *DBStack) Reset(checkpointID int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.simpleStack.Reset(checkpointID)
}

// Checkpoint checkpoints the changes to the stack such that a subsequent
// call to `Reset()` will reset any subsequent changes back to the state
// of the stack at the time of the latest checkpoint. The checkpointed state
// is also persisted to the DB.
func (d *DBStack) Checkpoint() (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	txn := d.meshDB.MiniHeaders.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	finalUpdates := map[common.Hash]*simplestack.Update{}
	for _, u := range d.simpleStack.GetUpdates() {
		if _, ok := finalUpdates[u.MiniHeader.Hash]; ok {
			delete(finalUpdates, u.MiniHeader.Hash)
		} else {
			finalUpdates[u.MiniHeader.Hash] = u
		}
	}

	for _, u := range finalUpdates {
		switch u.Type {
		case simplestack.Pop:
			if err := txn.Delete(u.MiniHeader.ID()); err != nil {
				return 0, err
			}
		case simplestack.Push:
			if err := txn.Insert(u.MiniHeader); err != nil {
				return 0, err
			}
		default:
			return 0, fmt.Errorf("Unrecognized update type encountered: %d", u.Type)
		}
	}

	if err := txn.Commit(); err != nil {
		return 0, err
	}
	return d.simpleStack.Checkpoint()
}
