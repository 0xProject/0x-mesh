package blockwatch

import (
	"errors"

	"github.com/0xProject/0x-mesh/meshdb"
)

// TODO(albrow): Needs to be optimized and made goroutine-safe.

// Stack allows performing basic stack operations on a stack of meshdb.MiniHeaders.
type Stack struct {
	meshDB *meshdb.MeshDB
	limit  int
}

// NewStack instantiates a new stack with the specified size limit. Once the size limit
// is reached, adding additional blocks will evict the deepest block.
func NewStack(meshDB *meshdb.MeshDB, limit int) *Stack {
	return &Stack{
		meshDB: meshDB,
		limit:  limit,
	}
}

// Pop removes and returns the latest block header on the block stack.
func (s *Stack) Pop() (*meshdb.MiniHeader, error) {
	miniHeaders, err := s.meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	if len(miniHeaders) == 0 {
		return nil, errors.New("Cannot pop from empty stack")
	}
	latestMiniHeader := miniHeaders[len(miniHeaders)-1]
	if err := s.meshDB.MiniHeaders.Delete(latestMiniHeader.ID()); err != nil {
		return nil, err
	}
	return latestMiniHeader, nil
}

// Push pushes a block header onto the block stack. If the stack limit is reached,
// it will remove the oldest block header and return it.
func (s *Stack) Push(block *meshdb.MiniHeader) (*meshdb.MiniHeader, error) {
	miniHeaders, err := s.meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	var oldestMiniHeader *meshdb.MiniHeader
	if len(miniHeaders) == s.limit {
		oldestMiniHeader = miniHeaders[0]
		if err := s.meshDB.MiniHeaders.Delete(oldestMiniHeader.ID()); err != nil {
			return nil, err
		}
	}
	if err := s.meshDB.MiniHeaders.Insert(block); err != nil {
		return nil, err
	}
	return oldestMiniHeader, nil
}

// Peek returns the latest block header (if exists) from the block stack without removing it.
func (s *Stack) Peek() (*meshdb.MiniHeader, error) {
	miniHeaders, err := s.meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	if len(miniHeaders) == 0 {
		return nil, nil
	}
	return miniHeaders[len(miniHeaders)-1], nil
}

// Inspect returns all the block headers currently on the stack
func (s *Stack) Inspect() ([]*meshdb.MiniHeader, error) {
	miniHeaders, err := s.meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	return miniHeaders, nil
}
