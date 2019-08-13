package blockwatch

// SimpleStack is a simple in-memory stack used in tests
type SimpleStack struct {
	limit       int
	miniHeaders []*MiniHeader
}

// NewSimpleStack instantiates a new SimpleStack
func NewSimpleStack(retentionLimit int) *SimpleStack {
	return &SimpleStack{
		limit:       retentionLimit,
		miniHeaders: []*MiniHeader{},
	}
}

// Peek returns the top of the stack
func (s *SimpleStack) Peek() (*MiniHeader, error) {
	if len(s.miniHeaders) == 0 {
		return nil, nil
	}
	return s.miniHeaders[len(s.miniHeaders)-1], nil
}

// Pop returns the top of the stack and removes it from the stack
func (s *SimpleStack) Pop() (*MiniHeader, error) {
	if len(s.miniHeaders) == 0 {
		return nil, nil
	}
	top := s.miniHeaders[len(s.miniHeaders)-1]
	s.miniHeaders = s.miniHeaders[:len(s.miniHeaders)-1]
	return top, nil
}

// Push adds a MiniHeader to the stack
func (s *SimpleStack) Push(miniHeader *MiniHeader) error {
	if len(s.miniHeaders) == s.limit {
		s.miniHeaders = s.miniHeaders[1:]
	}
	s.miniHeaders = append(s.miniHeaders, miniHeader)
	return nil
}

// Inspect returns the miniHeaders currently in the stack
func (s *SimpleStack) Inspect() ([]*MiniHeader, error) {
	return s.miniHeaders, nil
}
