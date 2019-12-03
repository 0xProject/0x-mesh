package simplestack

import (
	"fmt"
	"sync"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
)

type updateType int

const (
	Pop updateType = iota
	Push
)

// Update represents one update to the stack, either a pop or push of a miniHeader.
type Update struct {
	Type       updateType
	MiniHeader *miniheader.MiniHeader
}

// SimpleStack is a simple in-memory stack used in tests
type SimpleStack struct {
	limit       int
	miniHeaders []*miniheader.MiniHeader
	updates     []*Update
	mu          sync.Mutex
}

// NewSimpleStack instantiates a new SimpleStack
func NewSimpleStack(retentionLimit int, miniHeaders []*miniheader.MiniHeader) *SimpleStack {
	return &SimpleStack{
		limit:       retentionLimit,
		miniHeaders: miniHeaders,
		updates:     []*Update{},
	}
}

// Peek returns the top of the stack
func (s *SimpleStack) Peek() (*miniheader.MiniHeader, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.miniHeaders) == 0 {
		return nil, nil
	}
	return s.miniHeaders[len(s.miniHeaders)-1], nil
}

// Pop returns the top of the stack and removes it from the stack
func (s *SimpleStack) Pop() (*miniheader.MiniHeader, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.pop()
}

func (s *SimpleStack) pop() (*miniheader.MiniHeader, error) {
	if len(s.miniHeaders) == 0 {
		return nil, nil
	}
	top := s.miniHeaders[len(s.miniHeaders)-1]
	s.miniHeaders = s.miniHeaders[:len(s.miniHeaders)-1]
	s.updates = append(s.updates, &Update{
		Type:       Pop,
		MiniHeader: top,
	})
	return top, nil
}

// Push adds a miniheader.MiniHeader to the stack
func (s *SimpleStack) Push(miniHeader *miniheader.MiniHeader) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.push(miniHeader)
}

func (s *SimpleStack) push(miniHeader *miniheader.MiniHeader) error {
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

// PeekAll returns all the miniHeaders currently in the stack
func (s *SimpleStack) PeekAll() ([]*miniheader.MiniHeader, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.miniHeaders, nil
}

// Clear removes all items from the stack
func (s *SimpleStack) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.miniHeaders = []*miniheader.MiniHeader{}
	return nil
}

// Checkpoint checkpoints the changes to the stack such that a subsequent
// call to `Reset()` will reset any subsequent changes back to the state
// of the stack at the time of the latest checkpoint.
func (s *SimpleStack) Checkpoint() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.updates = []*Update{}
	return nil
}

// Reset resets the in-memory stack with the contents from the latest checkpoint
func (s *SimpleStack) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.reset()
}

func (s *SimpleStack) reset() error {
	for i := len(s.updates) - 1; i >= 0; i-- {
		u := s.updates[i]
		switch u.Type {
		case Pop:
			if err := s.push(u.MiniHeader); err != nil {
				return err
			}
		case Push:
			if _, err := s.pop(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unrecognized update type encountered: %d", u.Type)
		}
	}
	s.updates = []*Update{}
	return nil
}

// GetUpdates returns the updates applied since the last checkpoint
func (s *SimpleStack) GetUpdates() []*Update {
	return s.updates
}
