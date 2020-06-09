package blockwatch

import (
	"fmt"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
)

type MiniHeaderAlreadyExistsError struct {
	miniHeader *types.MiniHeader
}

func (e MiniHeaderAlreadyExistsError) Error() string {
	return fmt.Sprintf("cannot add miniHeader with the same number (%s) or hash (%s) as an existing miniHeader", e.miniHeader.Number.String(), e.miniHeader.Hash.Hex())
}

// Stack is a simple in-memory stack used in tests
type Stack struct {
	db *db.DB
}

// New instantiates a new Stack
func NewStack(db *db.DB) *Stack {
	return &Stack{
		db: db,
	}
}

// Peek returns the top of the stack
func (s *Stack) Peek() (*types.MiniHeader, error) {
	latestMiniHeader, err := s.db.GetLatestMiniHeader()
	if err != nil {
		if err == db.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return latestMiniHeader, nil
}

// Pop returns the top of the stack and removes it from the stack
func (s *Stack) Pop() (*types.MiniHeader, error) {
	removed, err := s.db.DeleteMiniHeaders(&db.MiniHeaderQuery{
		Limit: 1,
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Descending,
			},
		},
	})
	if err != nil {
		return nil, err
	} else if len(removed) == 0 {
		return nil, nil
	}
	return removed[0], nil
}

// Push adds a db.MiniHeader to the stack. It returns an error if
// the stack already contains a miniHeader with the same number or
// hash.
func (s *Stack) Push(miniHeader *types.MiniHeader) error {
	added, _, err := s.db.AddMiniHeaders([]*types.MiniHeader{miniHeader})
	if len(added) == 0 {
		return MiniHeaderAlreadyExistsError{miniHeader: miniHeader}
	}
	return err
}

// PeekAll returns all the miniHeaders currently in the stack
func (s *Stack) PeekAll() ([]*types.MiniHeader, error) {
	return s.db.FindMiniHeaders(nil)
}

// Clear removes all items from the stack and clears any set checkpoint
func (s *Stack) Clear() error {
	_, err := s.db.DeleteMiniHeaders(nil)
	return err
}
