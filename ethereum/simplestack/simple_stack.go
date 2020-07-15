package simplestack

import (
	"fmt"
	"sync"

	"github.com/0xProject/0x-mesh/common/types"
)

// UpdateType is the type of update applied to the in-memory stack.
type UpdateType int

const (
	Pop UpdateType = iota
	Push
)

// Update represents one update to the stack, either a pop or push of a miniHeader.
type Update struct {
	Type       UpdateType
	MiniHeader *types.MiniHeader
}

// SimpleStack is a simple in-memory stack used in tests.
type SimpleStack struct {
	limit              int
	miniHeaders        []*types.MiniHeader
	updates            []*Update
	mu                 sync.RWMutex
	latestCheckpointID int
}

// New instantiates a new SimpleStack.
func New(retentionLimit int, miniHeaders []*types.MiniHeader) *SimpleStack {
	return &SimpleStack{
		limit:       retentionLimit,
		miniHeaders: miniHeaders,
		updates:     []*Update{},
	}
}

// Peek returns the top of the stack.
func (s *SimpleStack) Peek() *types.MiniHeader {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.miniHeaders) == 0 {
		return nil
	}
	return s.miniHeaders[len(s.miniHeaders)-1]
}

// Pop returns the top of the stack and removes it from the stack.
func (s *SimpleStack) Pop() *types.MiniHeader {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.pop()
}

// you MUST acquire a lock on the mutex `mu` before calling `pop()`.
func (s *SimpleStack) pop() *types.MiniHeader {
	if len(s.miniHeaders) == 0 {
		return nil
	}
	top := s.miniHeaders[len(s.miniHeaders)-1]
	s.miniHeaders = s.miniHeaders[:len(s.miniHeaders)-1]
	s.updates = append(s.updates, &Update{
		Type:       Pop,
		MiniHeader: top,
	})
	return top
}

// Push adds a types.MiniHeader to the stack.
func (s *SimpleStack) Push(miniHeader *types.MiniHeader) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.push(miniHeader)
}

// you MUST acquire a lock on the mutex `mu` before calling `push()`.
func (s *SimpleStack) push(miniHeader *types.MiniHeader) error {
	for _, h := range s.miniHeaders {
		if h.Number.Int64() == miniHeader.Number.Int64() {
			return fmt.Errorf("attempted to push multiple blocks with block number %d to the stack", miniHeader.Number.Int64())
		}
	}

	if len(s.miniHeaders) == s.limit {
		s.miniHeaders = s.miniHeaders[1:]
	}
	s.miniHeaders = append(s.miniHeaders, miniHeader)
	s.updates = append(s.updates, &Update{
		Type:       Push,
		MiniHeader: miniHeader,
	})
	return nil
}

// PeekAll returns all the miniHeaders currently in the stack.
func (s *SimpleStack) PeekAll() []*types.MiniHeader {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return copy of miniHeaders array
	m := make([]*types.MiniHeader, len(s.miniHeaders))
	copy(m, s.miniHeaders)

	return m
}

// Clear removes all items from the stack and clears any set checkpoint.
func (s *SimpleStack) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.miniHeaders = []*types.MiniHeader{}
	s.updates = []*Update{}
	s.latestCheckpointID = 0
}

// Checkpoint checkpoints the changes to the stack such that a subsequent
// call to `Reset(checkpointID)` with the checkpointID returned from this
// call will reset any subsequent changes back to the state of the stack
// at the time of this checkpoint.
func (s *SimpleStack) Checkpoint() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.updates = []*Update{}
	s.latestCheckpointID++
	return s.latestCheckpointID
}

// Reset resets the in-memory stack with the contents from the latest checkpoint.
func (s *SimpleStack) Reset(checkpointID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.latestCheckpointID == 0 {
		return fmt.Errorf("Checkpoint() must be called before Reset() since without it the checkpointID is unspecified")
	} else if checkpointID != s.latestCheckpointID {
		return fmt.Errorf("Attempted to reset the stack to checkpoint %d but the latest checkpoint has ID %d", checkpointID, s.latestCheckpointID)
	}

	for i := len(s.updates) - 1; i >= 0; i-- {
		u := s.updates[i]
		switch u.Type {
		case Pop:
			if err := s.push(u.MiniHeader); err != nil {
				return err
			}
		case Push:
			_ = s.pop()
		default:
			return fmt.Errorf("Unrecognized update type encountered: %d", u.Type)
		}
	}
	s.updates = []*Update{}
	return nil
}

// GetUpdates returns the updates applied since the last checkpoint.
func (s *SimpleStack) GetUpdates() []*Update {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return copy of updates array
	u := make([]*Update, len(s.updates))
	copy(u, s.updates)
	return u
}
