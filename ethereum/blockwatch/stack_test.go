package blockwatch

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	miniHeaderOne = &types.MiniHeader{
		Number:    big.NewInt(1),
		Hash:      common.HexToHash("0x1"),
		Parent:    common.HexToHash("0x0"),
		Timestamp: time.Now().UTC(),
	}
)

func newTestStack(t *testing.T, ctx context.Context) *Stack {
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	return NewStack(database)
}

func TestStackPushPeekPop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := newTestStack(t, ctx)

	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)
	expectedMiniHeader := miniHeaderOne

	actualMiniHeaders, err := stack.PeekAll()
	require.NoError(t, err)
	require.Len(t, actualMiniHeaders, 1)
	assert.Equal(t, expectedMiniHeader, actualMiniHeaders[0])

	actualMiniHeader, err := stack.Peek()
	require.NoError(t, err)
	assert.Equal(t, expectedMiniHeader, actualMiniHeader)

	actualMiniHeaders, err = stack.PeekAll()
	require.NoError(t, err)
	assert.Len(t, actualMiniHeaders, 1)

	actualMiniHeader, err = stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, expectedMiniHeader, actualMiniHeader)

	actualMiniHeaders, err = stack.PeekAll()
	require.NoError(t, err)
	assert.Len(t, actualMiniHeaders, 0)
}

func TestStackErrorIfPushTwoHeadersWithSameNumber(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := newTestStack(t, ctx)
	// Push miniHeaderOne
	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Push miniHeaderOne again
	err = stack.Push(miniHeaderOne)
	assert.Error(t, err)
}

func TestStackClear(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := newTestStack(t, ctx)

	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)

	miniHeader, err := stack.Peek()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderOne, miniHeader)

	err = stack.Clear()
	require.NoError(t, err)

	miniHeader, err = stack.Peek()
	require.NoError(t, err)
	require.Nil(t, miniHeader)
}
