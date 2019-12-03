package blockwatch

import (
	"fmt"
	"sync"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
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

// SimpleStack is a simple in-memory stack used in tests
type SimpleStack struct {
	limit       int
	miniHeaders []*miniheader.MiniHeader
	updates     []*update
	mu          sync.Mutex
}

// NewSimpleStack instantiates a new SimpleStack
func NewSimpleStack(retentionLimit int) *SimpleStack {
	return &SimpleStack{
		limit:       retentionLimit,
		miniHeaders: []*miniheader.MiniHeader{},
		updates:     []*update{},
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
	s.updates = append(s.updates, &update{
		Type:       pop,
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
	if len(s.miniHeaders) == s.limit {
		s.miniHeaders = s.miniHeaders[1:]
	}
	s.miniHeaders = append(s.miniHeaders, miniHeader)
	s.updates = append(s.updates, &update{
		Type: push,
		// Optimization: We don't need to store the MiniHeader for
		// pushes since reverting a push involves a `Pop()` and the
		// value to pop is already in the `miniHeaders` data structure
		MiniHeader: nil,
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

	s.updates = []*update{}
	return nil
}

// Reset resets the stack with the contents from the latest checkpoint
func (s *SimpleStack) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := len(s.updates) - 1; i >= 0; i-- {
		u := s.updates[i]
		switch u.Type {
		case pop:
			if err := s.push(u.MiniHeader); err != nil {
				return err
			}
		case push:
			if _, err := s.pop(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unrecognized update type encountered: %d", u.Type)
		}
	}
	return nil
}
