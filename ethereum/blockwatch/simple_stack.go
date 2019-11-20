package blockwatch

import (
	"github.com/0xProject/0x-mesh/ethereum/miniheader"
)

// SimpleStack is a simple in-memory stack used in tests
type SimpleStack struct {
	limit       int
	miniHeaders []*miniheader.MiniHeader
}

// NewSimpleStack instantiates a new SimpleStack
func NewSimpleStack(retentionLimit int) *SimpleStack {
	return &SimpleStack{
		limit:       retentionLimit,
		miniHeaders: []*miniheader.MiniHeader{},
	}
}

// Peek returns the top of the stack
func (s *SimpleStack) Peek() (*miniheader.MiniHeader, error) {
	if len(s.miniHeaders) == 0 {
		return nil, nil
	}
	return s.miniHeaders[len(s.miniHeaders)-1], nil
}

// Pop returns the top of the stack and removes it from the stack
func (s *SimpleStack) Pop() (*miniheader.MiniHeader, error) {
	if len(s.miniHeaders) == 0 {
		return nil, nil
	}
	top := s.miniHeaders[len(s.miniHeaders)-1]
	s.miniHeaders = s.miniHeaders[:len(s.miniHeaders)-1]
	return top, nil
}

// Push adds a miniheader.MiniHeader to the stack
func (s *SimpleStack) Push(miniHeader *miniheader.MiniHeader) error {
	if len(s.miniHeaders) == s.limit {
		s.miniHeaders = s.miniHeaders[1:]
	}
	s.miniHeaders = append(s.miniHeaders, miniHeader)
	return nil
}

// PeekAll returns all the miniHeaders currently in the stack
func (s *SimpleStack) PeekAll() ([]*miniheader.MiniHeader, error) {
	return s.miniHeaders, nil
}

// Clear removes all items from the stack
func (s *SimpleStack) Clear() error {
	s.miniHeaders = []*miniheader.MiniHeader{}
	return nil
}
