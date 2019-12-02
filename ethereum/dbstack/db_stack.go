package dbstack

import (
	"fmt"
	"sync"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/meshdb"
)

type updateType int

const (
	pop updateType = iota
	push
)

type update struct {
	Type       updateType
	MiniHeader *miniheader.MiniHeader
}

// DBStack is an in-memory stack that can be sync'd with the DB
type DBStack struct {
	meshDB      *meshdb.MeshDB
	limit       int
	miniHeaders []*miniheader.MiniHeader
	mu          sync.Mutex
	updates     []*update
}

// New instantiates a new DBStack
func New(meshDB *meshdb.MeshDB, retentionLimit int) (*DBStack, error) {
	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	d := &DBStack{
		meshDB:      meshDB,
		limit:       retentionLimit,
		miniHeaders: miniHeaders,
		updates:     []*update{},
	}
	return d, nil
}

// Peek returns the top of the stack
func (d *DBStack) Peek() (*miniheader.MiniHeader, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.miniHeaders) == 0 {
		return nil, nil
	}
	return d.miniHeaders[len(d.miniHeaders)-1], nil
}

// Pop returns the top of the stack and removes it from the stack
func (d *DBStack) Pop() (*miniheader.MiniHeader, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.miniHeaders) == 0 {
		return nil, nil
	}
	top := d.miniHeaders[len(d.miniHeaders)-1]
	d.miniHeaders = d.miniHeaders[:len(d.miniHeaders)-1]
	d.updates = append(d.updates, &update{
		Type: pop,
		// Optimization: We don't need to store the MiniHeader for
		// pops since checkpoints a pop involves removing the latest
		// block header from the DB which doesn't require the explicit
		// header value
		MiniHeader: nil,
	})
	return top, nil
}

// Push adds a miniheader.MiniHeader to the stack
func (d *DBStack) Push(miniHeader *miniheader.MiniHeader) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.miniHeaders) == d.limit {
		d.miniHeaders = d.miniHeaders[1:]
	}
	d.miniHeaders = append(d.miniHeaders, miniHeader)
	d.updates = append(d.updates, &update{
		Type:       push,
		MiniHeader: miniHeader,
	})
	return nil
}

// PeekAll returns all the miniHeaders currently in the stack
func (d *DBStack) PeekAll() ([]*miniheader.MiniHeader, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.miniHeaders, nil
}

// Clear removes all items from the stack
func (d *DBStack) Clear() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.miniHeaders = []*miniheader.MiniHeader{}
	return d.meshDB.ClearAllMiniHeaders()
}

// Checkpoint checkpoints the changes to the stack such that a subsequent
// call to `Reset()` will reset any subsequent changes back to the state
// of the stack at the time of the latest checkpoint. The checkpointed state
// is also persisted to the DB.
func (d *DBStack) Checkpoint() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	txn := d.meshDB.MiniHeaders.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	for _, u := range d.updates {
		switch u.Type {
		case pop:
			latestMiniHeader, err := d.meshDB.FindLatestMiniHeader()
			if err != nil {
				return err
			}
			if latestMiniHeader == nil {
				return nil
			}
			if err := txn.Delete(latestMiniHeader.ID()); err != nil {
				return err
			}
		case push:
			if err := txn.Insert(u.MiniHeader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unrecognized update type encountered: %d", u.Type)
		}
	}

	if err := txn.Commit(); err != nil {
		return err
	}
	d.updates = []*update{}
	return nil
}

// Reset resets the stack with the contents from the latest checkpoint. This
// reset is also persisted to the backing DB.
func (d *DBStack) Reset() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.miniHeaders = []*miniheader.MiniHeader{}
	d.updates = []*update{}
	storedHeaders, err := d.meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return err
	}
	d.miniHeaders = storedHeaders
	return nil
}
